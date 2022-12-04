package init

import (
	telegram "gdrive-telegram-bot/src/clients/telegram"
	configModule "gdrive-telegram-bot/src/config"
	utils "gdrive-telegram-bot/src/utils"
)

// Start ...
func Start() (err error) {
	// Initialize config
	config, err := initConfig()
	if err != nil {
		return err
	}

	// Gets telegram connection
	bot, err := telegram.New(config)
	if err != nil {
		return err
	}

	// Starts conversation
	bot.StartConversation()

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
