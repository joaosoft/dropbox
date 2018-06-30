package godropbox

import (
	"fmt"

	golog "github.com/joaosoft/go-log/app"
	gomanager "github.com/joaosoft/go-manager/app"
)

type Dropbox struct {
	client        gomanager.IGateway
	config        *DropboxConfig
	pm            *gomanager.Manager
	isLogExternal bool

	// usage ...
	user   *user
	folder *folder
	file   *file
}

// NewDropbox ...
func NewDropbox(options ...dropboxOption) *Dropbox {
	pm := gomanager.NewManager(gomanager.WithRunInBackground(false))

	dropbox := &Dropbox{
		client: gomanager.NewSimpleGateway(),
		pm:     pm,
	}

	dropbox.Reconfigure(options...)

	if dropbox.isLogExternal {
		pm.Reconfigure(gomanager.WithLogger(log))
	}

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.GoDropbox.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(golog.WithLevel(level))
	}

	dropbox.config = &appConfig.GoDropbox

	return dropbox
}

// Api ...
func (d *Dropbox) User() *user {
	if d.user == nil {
		d.user = &user{
			client: d.client,
			config: d.config,
		}
	}
	return d.user
}

// Folder ...
func (d *Dropbox) Folder() *folder {
	if d.folder == nil {
		d.folder = &folder{
			client: d.client,
			config: d.config,
		}
	}
	return d.folder
}

// File ...
func (d *Dropbox) File() *file {
	if d.file == nil {
		d.file = &file{
			client: d.client,
			config: d.config,
		}
	}
	return d.file
}