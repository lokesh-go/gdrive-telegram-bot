package gdrive

import (
	configModule "gdrive-telegram-bot/src/config"
	utils "gdrive-telegram-bot/src/utils"

	googleService "github.com/lokesh-go/google-services/src/google"
	driveService "github.com/lokesh-go/google-services/src/google/services/drive"
)

// Module ...
type Module struct {
	config *configModule.Config
}

// New ...
func New(config *configModule.Config) *Module {
	// Returns
	return &Module{config}
}

// FileSearch ...
func (m *Module) FileSearch(text string) (response []SearchResponse, err error) {
	// Gets file search list
	searchResults, err := m.fileSearch(text)
	if err != nil {
		return nil, err
	}

	// Forms results
	response = []SearchResponse{}
	for _, r := range searchResults {
		res := SearchResponse{
			FileId:             r.Id,
			FileName:           r.Name,
			FileSize:           utils.ConvertBytesToHumanReadableForm(uint64(r.Size)),
			GDriveDownloadLink: r.DownloadLink,
			DownloadLink:       "",
		}
		response = append(response, res)
	}

	// Returns
	return response, nil
}

func (m *Module) fileSearch(text string) (searchResults []driveService.File, err error) {
	// Gets config
	scopes := googleService.Scopes{
		DriveScope: true,
	}
	config := &googleService.Config{
		ClientSecretFilePath: m.config.GDrive.Credential.ClientSecret,
		TokenPath:            m.config.GDrive.Credential.Token,
		Scopes:               scopes,
	}

	// Gets google client
	googleModule := googleService.New(config)
	client, err := googleModule.GetClient()
	if err != nil {
		return nil, err
	}

	// Gets drive service
	driveService, err := driveService.NewService(client)
	if err != nil {
		return nil, err
	}

	// Gets list
	searchResults, err = driveService.FileSearch(text, false, true)
	if err != nil {
		return nil, err
	}

	// Returns
	return searchResults, nil
}
