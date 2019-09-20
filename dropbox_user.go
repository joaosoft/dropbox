package dropbox

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/web"
)

type User struct {
	client manager.IGateway
	config *DropboxConfig
	logger logger.ILogger
}

type getUserResponse struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName       string `json:"given_name"`
		Surname         string `json:"surname"`
		FamiliarName    string `json:"familiar_name"`
		DisplayName     string `json:"display_name"`
		AbbreviatedName string `json:"abbreviated_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	Locale        string `json:"locale"`
	ReferralLink  string `json:"referral_link"`
	IsPaired      bool   `json:"is_paired"`
	AccountType   struct {
		Tag string `json:".tag"`
	} `json:"account_type"`
	RootInfo struct {
		Tag             string `json:".tag"`
		RootNamespaceID string `json:"root_namespace_id"`
		HomeNamespaceID string `json:"home_namespace_id"`
	} `json:"root_info"`
	Country string `json:"country"`
	Team    struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		SharingPolicies struct {
			SharedFolderMemberPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_member_policy"`
			SharedFolderJoinPolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_folder_join_policy"`
			SharedLinkCreatePolicy struct {
				Tag string `json:".tag"`
			} `json:"shared_link_create_policy"`
		} `json:"sharing_policies"`
		OfficeAddinPolicy struct {
			Tag string `json:".tag"`
		} `json:"office_addin_policy"`
	} `json:"team"`
	TeamMemberID string `json:"team_member_id"`
}

// Get ...
func (u *User) Get() (*getUserResponse, error) {
	dropboxResponse := &getUserResponse{}
	headers := manager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", u.config.Authorization.Access, u.config.Authorization.Token)},
	}

	if status, response, err := u.client.Request(http.MethodPost, u.config.Hosts.Api, "/users/get_current_account", string(web.ContentTypeEmpty), headers, nil); err != nil {
		err = u.logger.WithField("response", response).Error("error getting User account").ToError()
		return nil, err
	} else if status != http.StatusOK {
		err = u.logger.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError()
		return nil, err
	} else if response == nil {
		err = u.logger.Error("error getting User account").ToError()
		return nil, err
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			err = u.logger.Error("error converting Img User data").ToError()
			return nil, err
		}
		return dropboxResponse, nil
	}

	return nil, nil
}
