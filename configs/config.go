package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config는 JSON 구성 파일의 구조를 정의합니다.
type Config struct {
	JWTKey       string        `json:"jwt_key"`
	JWTTokenName string        `json:"jwt_token_name"`
	Ansible      AnsibleConfig `json:"ansible"`
	MongoDB      MongoDBConfig `json:"mongodb"`
	Erp          ErpConfig     `json:"erp"`
	Smtp         SmtpConfig    `json:"smtp"`
}

type AnsibleConfig struct {
	Playbook           string          `json:"playbook"`
	Options            AnsibleOptions  `json:"options"`
	DefaultExtraVars   string          `json:"default_extra_vars"`
	Patterns           AnsiblePatterns `json:"patterns"`
	PathStaticPlaybook string          `json:"path_static_playbook"`
}

type AnsibleOptions struct {
	ExtraVars string `json:"extra_vars"`
	Inventory string `json:"inventory"`
}

type AnsiblePatterns struct {
	InventoryINI string `json:"inventory_ini"`
	AnsibleYML   string `json:"ansible_yml"`
}

type MongoDBConfig struct {
	URL         string             `json:"url"`
	Database    string             `json:"database"`
	Collections MongoDBCollections `json:"collections"`
}

type MongoDBCollections struct {
	AnsibleProcessStatus string `json:"ansible_process_status"`
	Auth                 string `json:"auth"`
}

type ErpConfig struct {
	BaseURL string         `json:"base_url"`
	Login   ErpLoginConfig `json:"login"`
}

type ErpLoginConfig struct {
	Url         string `json:"url"`
	AdminId     string `json:"admin_id"`
	AdminPasswd string `json:"admin_passwd"`
	AllowType   string `json:"allow_type"`
	LoginBtn    string `json:"login_btn"`
}

type SmtpConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	From string `json:"from"`
}

// 전역 변수로 설정을 저장합니다.
var GlobalConfig *Config

// init 함수에서 전역 변수를 초기화합니다.
func init() {
	var err error
	GlobalConfig, err = LoadConfig("/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

// LoadConfig는 구성 파일을 읽어와 Config 구조체로 반환합니다.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
