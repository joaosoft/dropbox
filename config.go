package dropbox

import (
	"fmt"
	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Dropbox DropboxConfig `json:"dropbox"`
}

// DropboxConfig ...
type DropboxConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Authorization struct {
		Access string `json:"access"`
		Token  string `json:"token"`
	} `json:"authorization"`
	Hosts struct {
		Api     string `json:"api"`
		Content string `json:"content"`
	} `json:"hosts"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	if err != nil {
		log.Error(err.Error())
	}

	return appConfig, simpleConfig, err
}
