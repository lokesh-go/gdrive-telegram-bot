package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	configModule "gdrive-telegram-bot/src/config"
	utils "gdrive-telegram-bot/src/utils"
)

// Model ...
type Model struct {
	Config *configModule.Config
}

// New ...
func New(config *configModule.Config) *Model {
	// Returns
	return &Model{
		Config: config,
	}
}

// Connect ...
func (m *Model) Connect() (service *drive.Service, err error) {
	// Gets oauth config
	oauthConfig, err := m.getOAuthConfig()
	if err != nil {
		return nil, err
	}

	// Gets client
	client, err := m.getClient(oauthConfig)
	if err != nil {
		return nil, err
	}

	// Gets drive service
	service, err = drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Returns
	return service, nil
}

func (m *Model) getOAuthConfig() (oauthConfig *oauth2.Config, err error) {
	// Reads credential file
	bytes, err := ioutil.ReadFile(m.Config.GDrive.Credential.ClientSecret)
	if err != nil {
		return nil, err
	}

	// Gets google oauth config
	oauthConfig, err = google.ConfigFromJSON(bytes, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, err
	}

	// Returns
	return oauthConfig, nil
}

func (m *Model) getClient(oauthConfig *oauth2.Config) (client *http.Client, err error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	oauthToken, err := m.tokenFromFile()
	if err != nil {
		// Request a token from the web, then returns the retrieved token.
		oauthToken, err = m.getTokenFromWeb(oauthConfig)
		if err != nil {
			return nil, err
		}

		// Saves token
		err = m.saveToken(oauthToken)
		if err != nil {
			return nil, err
		}
	}

	// Returns
	return oauthConfig.Client(context.Background(), oauthToken), nil
}

func (m *Model) tokenFromFile() (oauthToken *oauth2.Token, err error) {
	// Reads token file
	err = utils.ReadJSONFile(m.Config.GDrive.Credential.Token, &oauthToken)
	if err != nil {
		return nil, err
	}

	// Returns
	return oauthToken, nil
}

func (m *Model) getTokenFromWeb(oauthConfig *oauth2.Config) (oauthToken *oauth2.Token, err error) {
	// Tokens from web
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	oauthToken, err = oauthConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}

	// Returns
	return oauthToken, nil
}

func (m *Model) saveToken(oauthToken *oauth2.Token) (err error) {
	// Saves token to the path
	file, err := os.OpenFile(m.Config.GDrive.Credential.Token, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(oauthToken)

	// Returns
	return nil
}
