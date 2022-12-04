package drive

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	config "github.com/lokesh-go/google-services/src/config"
)

// Service ...
type Service struct {
	drive *drive.Service
}

// NewService ...
func NewService(client *http.Client) (*Service, error) {
	// Gets new drive service
	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Returns
	return &Service{
		drive: service,
	}, nil
}

// FileSearch ...
func (s *Service) FileSearch(searchKey string, searchExtend bool, removeDownloadQuotaLimitFile bool) (filesData []File, err error) {
	// Forms file search query
	searchQuery := s.formFileSearchQuery(searchKey)

	// Gets file list
	files, err := s.getFileList(searchQuery, searchExtend)
	if err != nil {
		return nil, err
	}

	// Forms file response data
	filesData = s.formFilesData(files, removeDownloadQuotaLimitFile)

	// Returns
	return filesData, nil
}

func (s *Service) FileDownload(fileId string) (res *http.Response, err error) {
	// Downloads
	res, err = s.drive.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}

	// Returns
	return res, nil
}

func (s *Service) formFileSearchQuery(searchKey string) (query string) {
	// Forms search query
	query = config.FileSearchQueryNotContainsFolder + " and " + config.FileSearchQueryNameContains + searchKey + "' and " + config.FileSearchQueryNotContainsTrash

	// Returns
	return query
}

func (s *Service) getFileList(searchQuery string, searchExtend bool) (files []*drive.File, err error) {
	var pageToken string
	files = []*drive.File{}

	for {
		// Gets file lists
		fileList, err := s.drive.Files.List().SupportsAllDrives(config.SupportsAllDrives).IncludeItemsFromAllDrives(config.IncludeItemsFromAllDrives).Corpora(config.Corpora).Spaces(config.Spaces).Fields(googleapi.Field(config.FileSearchFieldsIncluded)).Q(searchQuery).PageSize(config.FileSearchPageSize).PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}

		// Appends
		files = append(files, fileList.Files...)

		// Searches not extends
		if !searchExtend {
			break
		}

		// Checks
		pageToken = fileList.NextPageToken
		if pageToken == "" {
			break
		}
	}

	// Returns
	return files, nil
}

func (s *Service) formFilesData(files []*drive.File, removeDownloadQuotaLimitFile bool) (filesData []File) {
	// Removes download quota limit file
	if removeDownloadQuotaLimitFile {
		filesData = s.removeDownloadQuotaLimit(files)
		return filesData
	}

	// Ranges
	filesData = []File{}
	for _, file := range files {
		// Forms file response
		list := formFileResponse(file)

		// Appends
		filesData = append(filesData, list)
	}

	// Returns
	return filesData
}

func (s *Service) removeDownloadQuotaLimit(files []*drive.File) (fileList []File) {
	// Defines waitgroup
	var wg sync.WaitGroup

	// Ranges
	fileList = []File{}
	for k, file := range files {
		// Adds
		wg.Add(1)

		// Go routine
		go func(file *drive.File, index int) {
			// Done
			defer wg.Done()

			// Checks
			_, err := s.FileDownload(file.Id)
			if err == nil {
				fileData := formFileResponse(file)
				fileList = append(fileList, fileData)
			}
		}(file, k)
	}
	// Wait
	wg.Wait()

	// Returns
	return fileList
}

func formFileResponse(file *drive.File) (fileData File) {
	// Forms download link
	fileDownloadLink := strings.ReplaceAll(config.FileDownloadLink, config.FileDownloadLinkPlaceholder, file.Id)

	// Forms file response
	fileData = File{
		Id:            file.Id,
		DriveId:       file.DriveId,
		DownloadLink:  fileDownloadLink,
		Name:          file.Name,
		Size:          file.Size,
		FileExtension: file.FileExtension,
		MimeType:      file.MimeType,
		Md5Checksum:   file.Md5Checksum,
	}

	// Returns
	return fileData
}
