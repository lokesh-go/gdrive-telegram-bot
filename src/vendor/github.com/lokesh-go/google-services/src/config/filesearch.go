package config

const (
	FileDownloadLink                 = "https://drive.google.com/uc?id=<fileIdPlaceholder>&export=download"
	FileDownloadLinkPlaceholder      = "<fileIdPlaceholder>"
	FileSearchPageSize               = 100
	SupportsAllDrives                = true
	IncludeItemsFromAllDrives        = true
	Corpora                          = "allDrives"
	Spaces                           = "drive"
	FileSearchFieldsIncluded         = "nextPageToken, files(id, driveId, name, mimeType, size, fileExtension, md5Checksum)"
	FileCreateFieldsIncluded         = "id"
	FileSearchQueryNameContains      = "name contains '"
	FileSearchQueryNotContainsFolder = "mimeType != 'application/vnd.google-apps.folder'"
	FileSearchQueryNotContainsTrash  = "trashed=false"
)
