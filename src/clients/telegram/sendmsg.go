package telegram

import (
	"os"
	"strconv"
	"time"

	gdrive "gdrive-telegram-bot/src/gdrive"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendWaitingMsg(update tgbotapi.Update) (quitChannel chan bool, resMsgId int) {
	quitChannel = make(chan bool)

	// Forms msg
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.config.Telegram.Template.SearchWaiting[0])
	msg.ReplyToMessageID = update.Message.MessageID

	// Sends
	sendRes, _ := b.bot.Send(msg)

	// Updating msg
	searchMsgIndex := 0
	go func() {
		for {
			// Starts from initial msg
			if searchMsgIndex == len(b.config.Telegram.Template.SearchWaiting) {
				searchMsgIndex = 0
			}

			// Switch case
			select {
			case <-quitChannel:
				return
			default:
				m := tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, b.config.Telegram.Template.SearchWaiting[searchMsgIndex])
				b.bot.Send(m)
			}

			// Sleep
			searchMsgIndex++
			time.Sleep(time.Duration(200 * time.Millisecond))
		}
	}()

	// Returns
	return quitChannel, sendRes.MessageID
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

func (b *Bot) sendResults(update tgbotapi.Update, res interface{}, resMsgId int, cmd string, quitChannel chan bool) {
	// Define msg
	var msg tgbotapi.EditMessageTextConfig

	// Switch case
	switch cmd {
	case SearchCmd:
		{
			// Checks
			logTemplate := b.config.Telegram.Template.Log
			var text string
			if res != nil {
				srchRes := res.([]gdrive.SearchResponse)
				text = logTemplate.Tag.Open + "\n" + logTemplate.Results + strconv.Itoa(len(srchRes)) + "\n" + logTemplate.Tag.Close
			} else {
				text = logTemplate.Tag.Open + "\n" + logTemplate.NoResults + "\n" + logTemplate.Tag.Close
			}
			msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
		}
	case Err:
		{
			text := res.(string)
			msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
		}
	}

	// Sends
	quitChannel <- true
	b.bot.Send(msg)
}
