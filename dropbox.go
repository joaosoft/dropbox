package dropbox

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Dropbox struct {
	client        manager.IGateway
	config        *DropboxConfig
	pm            *manager.Manager
	logger logger.ILogger
	isLogExternal bool

	// usage ...
	user   *User
	folder *Folder
	file   *File
}

// NewDropbox ...
func NewDropbox(options ...DropboxOption) *Dropbox {
	config, simpleConfig, err := NewConfig()
	pm := manager.NewManager(manager.WithRunInBackground(false))
	log := logger.NewLogDefault("dropbox", logger.DebugLevel)

	service := &Dropbox{
		client: manager.NewSimpleGateway(),
		pm:     pm,
		config: &config.Dropbox,
		logger: log,
	}

	if service.isLogExternal {
		pm.Reconfigure(manager.WithLogger(log))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Dropbox.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	return service
}

// Api ...
func (d *Dropbox) User() *User {
	if d.user == nil {
		d.user = &User{
			client: d.client,
			config: d.config,
			logger: d.logger,
		}
	}
	return d.user
}

// Folder ...
func (d *Dropbox) Folder() *Folder {
	if d.folder == nil {
		d.folder = &Folder{
			client: d.client,
			config: d.config,
			logger: d.logger,
		}
	}
	return d.folder
}

// File ...
func (d *Dropbox) File() *File {
	if d.file == nil {
		d.file = &File{
			client: d.client,
			config: d.config,
			logger: d.logger,
		}
	}
	return d.file
}
