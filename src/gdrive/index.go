package gdrive

import (
	configModule "gdrive-telegram-bot/src/config"
	utils "gdrive-telegram-bot/src/utils"
	"net/url"

	googleService "github.com/lokesh-go/google-services/src/google"
	driveService "github.com/lokesh-go/google-services/src/google/services/drive"
)

// Module ...
type Module struct {
	config       *configModule.Config
	driveService *driveService.Service
}

// New ...
func New(config *configModule.Config) (*Module, error) {
	// Gets config
	scopes := googleService.Scopes{
		DriveScope: true,
	}
	gconfig := &googleService.Config{
		ClientSecretFilePath: config.GDrive.Credential.ClientSecret,
		TokenPath:            config.GDrive.Credential.Token,
		Scopes:               scopes,
	}

	// Gets google client
	googleModule := googleService.New(gconfig)
	client, err := googleModule.GetClient()
	if err != nil {
		return nil, err
	}

	// Gets drive service
	driveService, err := driveService.NewService(client)
	if err != nil {
		return nil, err
	}

	// Returns
	return &Module{config, driveService}, nil
}

// FileSearch ...
func (m *Module) FileSearch(text string) (response []SearchResponse, err error) {
	// Gets file search list
	searchResults, err := m.driveService.FileSearch(text, false, true)
	if err != nil {
		return nil, err
	}

	// Forms results
	response = []SearchResponse{}
	for _, file := range searchResults {
		// Forms link
		endpoint := url.QueryEscape(file.Name) + "?id=" + file.Id
		link1 := m.config.Telegram.Results.Worker.Download1 + endpoint
		link2 := m.config.Telegram.Results.Worker.Download2 + endpoint
		link3 := m.config.Telegram.Results.Worker.MediaLink + endpoint

		res := SearchResponse{
			FileId:             file.Id,
			FileName:           file.Name,
			FileSize:           utils.ConvertBytesToHumanReadableForm(uint64(file.Size)),
			MimeType:           file.MimeType,
			GDriveDownloadLink: file.DownloadLink,
			DownloadLink1:      link1,
			DownloadLink2:      link2,
			MediaLink:          link3,
		}
		response = append(response, res)
	}

	// Returns
	return response, nil
}

// FileUpload ...
func (m *Module) FileUpload(fileName string, filePath string) (id string, err error) {
	// File upload
	uploadFileId, err := m.driveService.FileCreate(fileName, "text/html", m.config.Telegram.Results.UploadFolderId, filePath)
	if err != nil {
		return id, err
	}

	// Returns
	return uploadFileId, nil
}
