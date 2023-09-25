package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var Data = new(Config)

func init() {
	configPath := "./etc/consul-ext.yaml"
	var err error
	Data, err = LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	s, _ := yaml.Marshal(Data)
	log.Println(string(s))
	if Data.Consul == nil {
		log.Fatalf("Failed to load consul data: %v", err)
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

type Config struct {
	Database *struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBname   string `yaml:"dbname"`
	} `yaml:"database"`
	Gitea *struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Url    string `yaml:"url"`
		Owner  string `yaml:"owner"`
		Repo   string `yaml:"repo"`
		Branch string `yaml:"branch"`
		Token  string `yaml:"token"`
	} `yaml:"gitea"`
	Gitlab *struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Url    string `yaml:"url"`
		Repo   string `yaml:"repo"`
		Branch string `yaml:"branch"`
		Token  string `yaml:"token"`
	} `yaml:"gitlab"`
	Consul []*struct {
		Host    string   `yaml:"host"`
		Port    int      `yaml:"port"`
		SvcName []string `yaml:"svc_name"`
	} `yaml:"consul"`
	BackupTime string `yaml:"backup_time"` // backup store time,only days
	KVPefix    string `yaml:"kv_prefix"`
}
