package dropbox

import (
	"fmt"

	gomanager "github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Dropbox *DropboxConfig `json:"dropbox"`
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
func NewConfig() (*DropboxConfig, error) {
	appConfig := &AppConfig{}
	if _, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())

		return &DropboxConfig{}, nil
	}

	return appConfig.Dropbox, nil
}
