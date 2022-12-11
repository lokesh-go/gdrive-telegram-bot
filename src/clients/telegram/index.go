package telegram

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	configModule "gdrive-telegram-bot/src/config"
	gdrive "gdrive-telegram-bot/src/gdrive"
)

// Bot ...
type Bot struct {
	bot            *tgbotapi.BotAPI
	updatesChannel tgbotapi.UpdatesChannel
	config         *configModule.Config
	gdriveModule   *gdrive.Module
}

// New ...
func New(config *configModule.Config) (*Bot, error) {
	// Connects
	bot, updatesChannel, err := connect(config)
	if err != nil {
		return nil, err
	}

	// Gets the gdrive module
	gdriveModule, err := gdrive.New(config)
	if err != nil {
		return nil, err
	}

	// Returns
	return &Bot{
		bot:            bot,
		updatesChannel: updatesChannel,
		config:         config,
		gdriveModule:   gdriveModule,
	}, nil
}

func connect(config *configModule.Config) (bot *tgbotapi.BotAPI, updatesChannel tgbotapi.UpdatesChannel, err error) {
	// Connects
	bot, err = tgbotapi.NewBotAPI(os.Getenv(config.Telegram.Token))
	if err != nil {
		return nil, nil, err
	}

	// Sets debugger
	bot.Debug = config.Telegram.Debug

	// Update config
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = config.Telegram.UpdateTimeout

	// Start polling Telegram for updates.
	updatesChannel = bot.GetUpdatesChan(updateConfig)

	// Returns
	return bot, updatesChannel, nil
}
