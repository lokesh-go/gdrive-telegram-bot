package config

// Config ...
type Config struct {
	GDrive struct {
		Credential struct {
			ClientSecret string
			Token        string
		}
	}
	Telegram struct {
		Token         string
		Debug         bool
		UpdateTimeout int
		Template      struct {
			Log struct {
				Tag struct {
					Open  string
					Mid   string
					Close string
				}
				Query     string
				UserName  string
				Name      string
				ID        string
				Results   string
				Error     string
				Example   string
				NoResults string
			}
			SearchWaiting []string
		}
		Admin struct {
			ChatId []string
		}
		Results struct {
			Worker struct {
				Download1 string
				Download2 string
				MediaLink string
				SearchRes string
			}
			Template struct {
				Placeholder struct {
					Query              string
					Count              string
					Results            string
					HexColor           string
					GDriveDownloadLink string
					FileName           string
					FileSize           string
					DownloadLink1      string
					DownloadLink2      string
					MediaLink          string
					MXPlayerLink       string
				}
				SearchResults string
				MXPlayer      string
				OtherPlayer   string
			}
			FilePath       string
			UploadFolderId string
			ButtonText     string
		}
	}
}
