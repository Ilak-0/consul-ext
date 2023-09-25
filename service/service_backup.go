package service

import (
	"consul-ext/pkg/config"
	"consul-ext/pkg/consul"
	"consul-ext/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
	"sync"
	"time"
)

func DailySyncConsulSvcs() error {
	//sync consul svcs to db storage
	for _, consulAddress := range consul.AddressList {
		err := PullServicesToDB(consulAddress)
		if err != nil {
			return err
		}
	}
	//delete expired records
	if err := consul.ConsulSvcsControl.Delete("DATEDIFF(NOW(), backup_time) > " + config.Data.BackupTime); err != nil {
		log.Println("delete expired records failed :", err)
		return err
	}
	return nil
}

func PullServicesToDB(address string) error {
	now := time.Now().Format("2006-01-02 00:00:00")
	svcList, err := GetAllSvcsFromConsul(address)
	if err != nil {
		return err
	}
	for _, svc := range svcList {
		s, err := json.Marshal(svc)
		if err != nil {
			return err
		}
		svc := consul.ConsulSvcs{
			SvcId:          svc.ServiceID,
			SvcName:        svc.ServiceName,
			ConsulAddress:  address,
			SvcCatalogJson: string(s),
			BackupTime:     now,
		}
		if err := consul.ConsulSvcsControl.CreateOnDuplicate(&svc); err != nil {
			return err
		}
	}
	return nil
}

func GetAllSvcsFromConsul(address string) ([]*api.CatalogService, error) {
	svcList := make([]*api.CatalogService, 0)
	consulConfig, err := consul.ClientMap.GetConfig(address)
	if err != nil {
		return nil, err
	}
	consulClient := consulConfig.Client
	svcMap, _, err := consulClient.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}
	for n, _ := range svcMap {
		if consulConfig.SvcName == nil || len(consulConfig.SvcName) == 0 || tool.ContainStr(consulConfig.SvcName, n) {
			svcs, _, err := consulClient.Catalog().Service(n, "", nil)
			if err != nil {
				return nil, err
			}
			svcList = append(svcList, svcs...)
		}
	}
	return svcList, nil
}

// ConsulSvcsRestore avoid using unless necessary
func ConsulSvcsRestore(params *consul.SvcsRestoreParams) error {
	var (
		err         error
		svcNameList = make([]string, 0)
		svcs        = make([]consul.ConsulSvcs, 0)
	)
	//change date to time.time

	if params.SvcName == "all" {
		svcs, err = consul.ConsulSvcsControl.ListWithoutPage(consul.CONSUL_SVC_FIELDS,
			fmt.Sprintf("backup_time='%s' and consul_address='%s'", params.BackupTime.Format("2006-01-02 00:00:00"), params.ReadConsulAddress))
		if err != nil {
			return err
		}
	} else {
		svcNameList = strings.Split(params.SvcName, ",")
		for _, svc := range svcNameList {
			singleSvcs, err := consul.ConsulSvcsControl.ListWithoutPage(consul.CONSUL_SVC_FIELDS,
				fmt.Sprintf("backup_time='%s' and consul_address='%s' and svc_name='%s'", params.BackupTime.Format("2006-01-02 00:00:00"), params.ReadConsulAddress, svc))
			if err != nil {
				return err
			}
			svcs = append(svcs, singleSvcs...)
		}
	}

	if len(svcs) == 0 {
		return nil
	}
	consulClient, err := consul.GetConsulClientByAddress(params.WriteConsulAddress)
	if err != nil {
		return err
	}
	//decide whether to remove the existing service
	errChan := make(chan error, 100)
	done := make(chan struct{}, 1)
	defer func() {
		close(errChan)
		close(done)
	}()
	if params.DeleteCurrentSvcs {
		go func() {
			defer func() {
				recover()
			}()
			registerSvcsWithDeregister(svcNameList, svcs, consulClient, errChan)
			done <- struct{}{}
		}()
	} else {
		go func() {
			defer func() {
				recover()
			}()
			registerSvcs(svcs, consulClient, errChan)
			done <- struct{}{}
		}()
	}
	select {
	case <-done:
		return nil
	case er := <-errChan:
		return er
	}
}

func registerSvcsWithDeregister(svcNameList []string, svcs []consul.ConsulSvcs, consulClient *api.Client, errChan chan error) {
	defer func() {
		recover()
	}()
	wg := sync.WaitGroup{}
	limitChan := make(chan bool, 100)
	defer close(limitChan)
	for _, svc := range svcNameList {
		ss, _, err := consulClient.Catalog().Service(svc, "", nil)
		if err != nil {
			errChan <- err
			return
		}
		for _, v := range ss {
			limitChan <- true
			wg.Add(1)
			go func(v *api.CatalogService) {
				defer func() {
					<-limitChan
					wg.Done()
				}()
				if v.ServiceID == "" {
					errChan <- fmt.Errorf("service id is empty")
					return
				}
				log.Println("deregister service:", v.ServiceID)
				if err := consulClient.Agent().ServiceDeregister(v.ServiceID); err != nil {
					errChan <- err
					return
				}
			}(v)
		}
	}
	wg.Wait()
	registerSvcs(svcs, consulClient, errChan)
	log.Println("registerSvcsWithDeregister finish")
	return
}

func registerSvcs(svcs []consul.ConsulSvcs, consulClient *api.Client, errChan chan error) {
	defer func() {
		recover()
	}()
	wg := sync.WaitGroup{}
	limitChan := make(chan bool, 100)
	defer close(limitChan)
	for _, svc := range svcs {
		log.Println("register service:", svc)
		limitChan <- true
		wg.Add(1)
		go func(svc consul.ConsulSvcs) {
			defer func() {
				<-limitChan
				wg.Done()
			}()
			s := new(api.CatalogService)
			if err := json.Unmarshal([]byte(svc.SvcCatalogJson), s); err != nil {
				log.Println(err)
				errChan <- err
				return
			}
			agentSvc := consul.ConvertCatalogServiceAddScrape(s)
			if err := consulClient.Agent().ServiceRegister(agentSvc); err != nil {

			}
		}(svc)
	}
	wg.Wait()
	log.Println("registerSvcs finish")
	return
}
