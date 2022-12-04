package drive

// File ...
type File struct {
	Id            string `json:"id"`
	DriveId       string `json:"driveId"`
	DownloadLink  string `json:"downloadLink"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	FileExtension string `json:"fileExtension"`
	MimeType      string `json:"mimeType"`
	Md5Checksum   string `json:"md5Checksum"`
}
