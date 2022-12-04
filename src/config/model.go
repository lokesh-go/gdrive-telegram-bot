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
				ID        string
				Results   string
				Error     string
				Example   string
				NoResults string
			}
			SearchWaiting string
		}
		Admin struct {
			ChatId []string
		}
	}
}
