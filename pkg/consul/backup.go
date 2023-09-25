package consul

import (
	"consul-ext/db"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"

	"fmt"
	"time"
)

const (
	CONSUL_SVC_FIELDS = "svc_id,svc_name,consul_address,svc_catalog_json,backup_time"
)

var ConsulSvcsControl = &ConsulSvcsDao{
	sourceDB:  db.Gpool,
	replicaDB: []*gorm.DB{db.Gpool},
	m:         new(ConsulSvcs),
}

type ConsulSvcs struct {
	SvcId          string `json:"svc_id" gorm:"column:svc_id" db:"svc_id"`
	SvcName        string `json:"svc_name" gorm:"column:svc_name" db:"svc_name"`
	ConsulAddress  string `json:"consul_address" gorm:"column:consul_address" db:"consul_address"`       // consul 注册地址
	SvcCatalogJson string `json:"svc_catalog_json" gorm:"column:svc_catalog_json" db:"svc_catalog_json"` // 服务信息详情
	BackupTime     string `json:"backup_time" gorm:"column:backup_time" db:"backup_time"`                // 备份时间
}

func (m *ConsulSvcs) TableName() string {
	return "consul_svcs"
}

type ConsulSvcsDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *ConsulSvcs
}

func NewConsulSvcsDao(dbs ...*gorm.DB) *ConsulSvcsDao {
	dao := new(ConsulSvcsDao)
	switch len(dbs) {
	case 0:
		panic("database connection required")
	case 1:
		dao.sourceDB = dbs[0]
		dao.replicaDB = []*gorm.DB{dbs[0]}
	default:
		dao.sourceDB = dbs[0]
		dao.replicaDB = dbs[1:]
	}
	return dao
}

func (d *ConsulSvcsDao) Create(obj *ConsulSvcs) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("ConsulSvcsDao: %w", err)
	}
	return nil
}

func (d *ConsulSvcsDao) CreateOnDuplicate(obj *ConsulSvcs) error {
	err := d.sourceDB.Model(d.m).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "service_id"}, {Name: "backup_time"}},
			DoUpdates: clause.AssignmentColumns([]string{"svc_name", "consul_address", "svc_catalog_json"}),
		}).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("ConsulSvcsDao: %w", err)
	}
	return nil
}

func (d *ConsulSvcsDao) Get(fields, where string) (*ConsulSvcs, error) {
	items, err := d.List(fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("ConsulSvcsDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *ConsulSvcsDao) List(fields, where string, offset, limit int) ([]ConsulSvcs, error) {
	var results []ConsulSvcs
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("ConsulSvcsDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func (d *ConsulSvcsDao) Update(where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("ConsulSvcsDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *ConsulSvcsDao) Delete(where string, args ...interface{}) error {
	if len(where) == 0 {
		return gorm.ErrRecordNotFound
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("ConsulSvcsDao: Delete where=%s: %w", where, err)
	}
	return nil
}

func (d *ConsulSvcsDao) ListWithoutPage(fields, where string) ([]ConsulSvcs, error) {
	var results []ConsulSvcs
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("ConsoleProjectPlatTeamDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func ConvertCatalogServiceAddScrape(catalogService *api.CatalogService) *api.AgentServiceRegistration {
	agentServiceRegistration := &api.AgentServiceRegistration{
		ID:      catalogService.ServiceID,
		Name:    catalogService.ServiceName,
		Port:    catalogService.ServicePort,
		Address: catalogService.ServiceAddress,
		Tags:    catalogService.ServiceTags,
		Meta:    catalogService.ServiceMeta,
	}
	return agentServiceRegistration
}

type SvcsRestoreParams struct {
	SvcName            string    `json:"svc_name" bind:"required"`
	BackupTime         time.Time `json:"backup_time" `
	ReadConsulAddress  string    `json:"read_consul_address" bind:"required"`
	WriteConsulAddress string    `json:"write_consul_address"`
	DeleteCurrentSvcs  bool      `json:"delete_current_svcs" bind:"required"`
}

type WatchPayload struct {
	Key string `json:"Key"`
	// Add more fields here if needed
}
