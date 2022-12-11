package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartConversation ...
func (b *Bot) StartConversation() {
	// Starts
	for update := range b.updatesChannel {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		// Starts process
		go b.startProcess(update)
	}
}

func (b *Bot) startProcess(update tgbotapi.Update) {
	// Sends
	quit, resMsgId := b.sendWaitingMsg(update)

	// Gets response
	var logErr interface{}
	res, cmd, statusCode, err := b.router(update.Message.Text)
	if err != nil {
		res, logErr = b.getCustomErrorMsg(statusCode, err)
	}

	// Sends
	resCount, errLog := b.sendResults(update, res, resMsgId, cmd, quit)
	if errLog != nil {
		logErr = errLog
	}

	// Sends
	b.sendLogs(update, resCount, logErr)
}
