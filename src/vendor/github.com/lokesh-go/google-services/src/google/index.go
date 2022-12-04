package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	gdriveScopes "github.com/lokesh-go/google-services/src/google/services/drive/scopes"
	utils "github.com/lokesh-go/google-services/src/utils"
)

// Config ...
type Config struct {
	ClientSecretFilePath string
	TokenPath            string
	Scopes               Scopes
}

type Scopes struct {
	DriveScope bool
}

// New ...
func New(config *Config) *Config {
	// Returns
	return config
}

// GetClient ...
func (c *Config) GetClient() (client *http.Client, err error) {
	// Gets oauth config
	oauthConfig, err := c.getOAuthConfig()
	if err != nil {
		return nil, err
	}

	// Gets client
	client, err = c.getClient(oauthConfig)
	if err != nil {
		return nil, err
	}

	// Returns
	return client, nil
}

func (c *Config) getOAuthConfig() (oauthConfig *oauth2.Config, err error) {
	// Reads credential file
	bytes, err := ioutil.ReadFile(c.ClientSecretFilePath)
	if err != nil {
		return nil, err
	}

	// Gets scops
	scopes := c.getScopes()

	// Gets google oauth config
	oauthConfig, err = google.ConfigFromJSON(bytes, scopes)
	if err != nil {
		return nil, err
	}

	// Returns
	return oauthConfig, nil
}

func (c *Config) getClient(oauthConfig *oauth2.Config) (client *http.Client, err error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	oauthToken, err := c.tokenFromFile()
	if err != nil {
		// Request a token from the web, then returns the retrieved token.
		oauthToken, err = getTokenFromWeb(oauthConfig)
		if err != nil {
			return nil, err
		}

		// Saves token
		err = c.saveToken(oauthToken)
		if err != nil {
			return nil, err
		}
	}

	// Refresh token if token has expired
	if oauthToken.Expiry.Before(time.Now()) {
		// Gets new token
		oauthToken, err = oauthConfig.TokenSource(context.Background(), oauthToken).Token()
		if err != nil {
			return nil, err
		}

		// Saves token
		err = c.saveToken(oauthToken)
		if err != nil {
			return nil, err
		}
	}

	// Returns
	return oauthConfig.Client(context.Background(), oauthToken), nil
}

func (c *Config) tokenFromFile() (oauthToken *oauth2.Token, err error) {
	// Reads token file
	err = utils.ReadJSONFile(c.TokenPath, &oauthToken)
	if err != nil {
		return nil, err
	}

	// Returns
	return oauthToken, nil
}

func getTokenFromWeb(oauthConfig *oauth2.Config) (oauthToken *oauth2.Token, err error) {
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

func (c *Config) saveToken(oauthToken *oauth2.Token) (err error) {
	// Saves token to the path
	file, err := os.OpenFile(c.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	json.NewEncoder(file).Encode(oauthToken)

	// Returns
	return nil
}

func (c *Config) getScopes() (scope string) {
	// Gets drive scopes
	if c.Scopes.DriveScope {
		scope = gdriveScopes.Get()
	}

	// Returns
	return scope
}
