package telegram

import (
	"fmt"
	"os"
	"strconv"

	gdrive "gdrive-telegram-bot/src/gdrive"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendWaitingMsg(update tgbotapi.Update) (resMsgId int) {
	// Forms msg
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.config.Telegram.Template.SearchWaiting)
	msg.ReplyToMessageID = update.Message.MessageID

	// Sends
	sendRes, _ := b.bot.Send(msg)

	// Returns
	return sendRes.MessageID
}

func (b *Bot) sendLogs(update tgbotapi.Update, resCount int, logErr interface{}) {
	// Checks
	var errLog string
	if logErr != nil {
		errLog = logErr.(string)
	}

	// Gets telegram log template
	logTemplate := b.config.Telegram.Template.Log
	msg := logTemplate.Tag.Open + "\n" + logTemplate.Query + update.Message.Text + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.UserName + update.Message.Chat.UserName + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.ID + strconv.Itoa(int(update.Message.Chat.ID)) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Results + strconv.Itoa(int(resCount)) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Error + errLog + "\n" + logTemplate.Tag.Close

	// Sends
	for _, user := range b.config.Telegram.Admin.ChatId {
		chatId, _ := strconv.Atoi(os.Getenv(user))
		m := tgbotapi.NewMessage(int64(chatId), msg)
		b.bot.Send(m)
	}
}

func (b *Bot) sendResults(update tgbotapi.Update, res interface{}, resMsgId int, cmd string) {
	switch cmd {
	case SearchCmd:
		{
			// Checks
			var text string
			if res != nil {
				srchRes := res.([]gdrive.SearchResponse)
				fmt.Println(srchRes)
			} else {
				logTemplate := b.config.Telegram.Template.Log
				text = logTemplate.Tag.Open + "\n" + logTemplate.NoResults + "\n" + logTemplate.Tag.Close
			}

			// Sends
			msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
			b.bot.Send(msg)
		}
	case Err:
		{
			text := res.(string)
			msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
			b.bot.Send(msg)
		}
	}
}
