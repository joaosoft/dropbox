package dropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joaosoft/manager"
)

type writeMode string

const (
	writeModeAdd       writeMode = "add"
	writeModeOverwrite           = "overwrite"
)

type File struct {
	client manager.IGateway
	config *DropboxConfig
}

type uploadFileRequest struct {
	Path       string    `json:"path"`
	Mode       writeMode `json:"mode"`
	AutoRename bool      `json:"auto_rename"`
	Mute       bool      `json:"mute"`
}

type uploadFileResponse struct {
	Name           string    `json:"name"`
	ID             string    `json:"id"`
	ClientModified time.Time `json:"client_modified"`
	ServerModified time.Time `json:"server_modified"`
	Rev            string    `json:"rev"`
	Size           int       `json:"size"`
	PathLower      string    `json:"path_lower"`
	PathDisplay    string    `json:"path_display"`
	SharingInfo    struct {
		ReadOnly             bool   `json:"read_only"`
		ParentSharedFolderID string `json:"parent_shared_folder_id"`
		ModifiedBy           string `json:"modified_by"`
	} `json:"sharing_info"`
	PropertyGroups []struct {
		TemplateID string `json:"template_id"`
		Fields     []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	} `json:"property_groups"`
	HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
	ContentHash              string `json:"content_hash"`
}

func (f *File) Upload(path string, file []byte) (*uploadFileResponse, error) {
	var err error
	var bodyArgs []byte
	args := uploadFileRequest{
		Path:       path,
		Mode:       writeModeOverwrite,
		AutoRename: true,
		Mute:       false,
	}

	if bodyArgs, err = json.Marshal(args); err != nil {
		err = log.Error("errors converting upload input arguments").ToError()
		return nil, err
	}

	headers := manager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":    {"application/octet-stream"},
		"Dropbox-API-Arg": {string(bodyArgs)},
	}

	dropboxResponse := &uploadFileResponse{}
	if err != nil {
		err = log.Error("errors marshal arguments").ToError()
		return nil, err
	}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/upload", headers, file); err != nil {
		log.WithField("response", response).Errorf("error uploading file to %s", path)
		return nil, err
	} else if status != http.StatusOK {
		var err error
		err = log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError()
		return nil, err
	} else if response == nil {
		err = log.WithField("response", response).Errorf("error uploading file to %s", path).ToError()
		return nil, err
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			err = log.Error("errors converting Dropbox response data").ToError()
			return nil, err
		}
		return dropboxResponse, nil
	}

	return nil, nil
}

type downloadFileRequest struct {
	Path string `json:"path"`
}

func (f *File) Download(path string) ([]byte, error) {
	var err error
	var bodyArgs []byte
	args := downloadFileRequest{
		Path: path,
	}

	if bodyArgs, err = json.Marshal(args); err != nil {
		err = log.Error("errors converting download input arguments").ToError()
		return nil, err
	}

	headers := manager.Headers{
		"Authorization":   {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Dropbox-API-Arg": {string(bodyArgs)},
	}

	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Content, "/files/download", headers, []byte("")); err != nil {
		err = log.WithField("response", response).Error("errors downloading File").ToError()
		return nil, err
	} else if status != http.StatusOK {
		err = log.WithField("response", response).WithFields(map[string]interface{}{"response": string(response)}).Errorf("response status %d instead of %d", status, http.StatusOK).ToError()
		return nil, err
	} else if response == nil {
		err = log.Error("errors downloading File").ToError()
		return nil, err
	} else {
		return response, nil
	}

	return nil, nil
}

type deleteFileRequest struct {
	Path string `json:"path"`
}

type deleteFileResponse struct {
	Metadata struct {
		Tag            string    `json:".tag"`
		Name           string    `json:"name"`
		ID             string    `json:"id"`
		ClientModified time.Time `json:"client_modified"`
		ServerModified time.Time `json:"server_modified"`
		Rev            string    `json:"rev"`
		Size           int       `json:"size"`
		PathLower      string    `json:"path_lower"`
		PathDisplay    string    `json:"path_display"`
		SharingInfo    struct {
			ReadOnly             bool   `json:"read_only"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			ModifiedBy           string `json:"modified_by"`
		} `json:"sharing_info"`
		PropertyGroups []struct {
			TemplateID string `json:"template_id"`
			Fields     []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"property_groups"`
		HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
		ContentHash              string `json:"content_hash"`
	} `json:"metadata"`
}

func (f *File) Delete(path string) (*deleteFileResponse, error) {
	if path == "/" {
		path = ""
	}
	body, err := json.Marshal(deleteFileRequest{
		Path: path,
	})
	if err != nil {
		err = log.Error("errors marshal arguments").ToError()
		return nil, err
	}

	headers := manager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":  {"application/json"},
	}

	dropboxResponse := &deleteFileResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Api, "/files/delete_v2", headers, body); err != nil {
		err = log.WithField("response", response).Error("errors deleting File").ToError()
		return nil, err
	} else if status != http.StatusOK {
		err = log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError()
		return nil, err
	} else if response == nil {
		err = log.Error("errors deleting File").ToError()
		return nil, err
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			err = log.Error("errors converting Dropbox response data").ToError()
			return nil, err
		}
		return dropboxResponse, nil
	}

	return nil, nil
}
