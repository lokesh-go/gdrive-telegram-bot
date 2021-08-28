package init

import (
	configModule "gdrive-telegram-bot/src/config"
	gdrive "gdrive-telegram-bot/src/storage/cloud/gdrive"
	utils "gdrive-telegram-bot/src/utils"
)

// Start ...
func Start() (err error) {
	// Initialize config
	config, err := initConfig()
	if err != nil {
		return err
	}

	// Initialize gdrive connection
	gdriveModule := gdrive.New(config)
	driveService, err := gdriveModule.Connect()
	if err != nil {
		return err
	}

	// Returns
	return nil
}

func initConfig() (config *configModule.Config, err error) {
	// Binds config models
	err = utils.ReadJSONFile("config/config.json", &config)
	if err != nil {
		return nil, err
	}

	// Returns
	return config, nil
}
