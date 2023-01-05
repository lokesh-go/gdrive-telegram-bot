package telegram

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	gdrive "gdrive-telegram-bot/src/gdrive"
	utils "gdrive-telegram-bot/src/utils"

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
	msg := logTemplate.Tag.Open + "\n" + logTemplate.Query + update.Message.Text + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.UserName + update.Message.Chat.UserName + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Name + update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.ID + strconv.Itoa(int(update.Message.Chat.ID)) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Results + strconv.Itoa(int(resCount)) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Error + errLog + "\n" + logTemplate.Tag.Close

	// Sends
	for _, user := range b.config.Telegram.Admin.ChatId {
		chatId, _ := strconv.Atoi(os.Getenv(user))
		m := tgbotapi.NewMessage(int64(chatId), msg)
		b.bot.Send(m)
	}
}

func (b *Bot) sendResults(update tgbotapi.Update, res interface{}, resMsgId int, cmd string, quitChannel chan bool) (resCount int, logErrMsg interface{}) {
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
				if len(srchRes) == 0 {
					text = logTemplate.Tag.Open + "\n" + logTemplate.NoResults + "\n" + logTemplate.Tag.Close
					msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
				} else {
					resCount = len(srchRes)

					// Gets search text
					_, searchText := getSearchKeyDetails(update.Message.Text)

					// Forms template result string
					templateString := b.formSearchResultTemplate(srchRes, searchText)

					// Forms random file name
					randomText, _ := utils.GetRandomHexValue(10)
					fileName := searchText + "_" + randomText + ".html"

					// Creates file
					filePath := os.Getenv(b.config.Telegram.Results.FilePath) + fileName
					f, err := os.Create(filePath)
					if err != nil {
						userErrMsg, logErrMsg := b.getCustomErrorMsg(http.StatusInternalServerError, err)
						text := userErrMsg.(string)
						msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
						quitChannel <- true
						b.bot.Send(msg)
						return resCount, logErrMsg
					}

					// Writes file string
					f.WriteString(templateString)
					f.Close()

					// Upload file into gdrive
					fileUploadedId, err := b.gdriveModule.FileUpload(fileName, filePath)
					if err != nil {
						userErrMsg, logErrMsg := b.getCustomErrorMsg(http.StatusInternalServerError, err)
						text := userErrMsg.(string)
						msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
						quitChannel <- true
						b.bot.Send(msg)
						return resCount, logErrMsg
					}

					// Forms inline keyborad button
					url := b.config.Telegram.Results.Worker.SearchRes + fileUploadedId
					inlineButton := tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text: b.config.Telegram.Results.ButtonText,
							URL:  &url,
						},
					)
					inlineMarkUp := tgbotapi.NewInlineKeyboardMarkup(
						inlineButton,
					)

					// Forms reply text
					replyText := logTemplate.Tag.Open + "\n" + logTemplate.Results + strconv.Itoa(resCount) + " 🎬🤩\n" + logTemplate.Tag.Close

					// Sends msg
					m := tgbotapi.NewEditMessageTextAndMarkup(update.Message.Chat.ID, resMsgId, replyText, inlineMarkUp)
					quitChannel <- true
					b.bot.Send(m)

					// Delete file
					os.Remove(filePath)

					// Returns
					return resCount, nil
				}
			} else {
				text = logTemplate.Tag.Open + "\n" + logTemplate.NoResults + "\n" + logTemplate.Tag.Close
				msg = tgbotapi.NewEditMessageText(update.Message.Chat.ID, resMsgId, text)
			}
		}
	case InfoCmd:
		{
			text := res.(string)
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

	// Returns
	return resCount, nil
}

func (b *Bot) formSearchResultTemplate(res []gdrive.SearchResponse, searchQuery string) (resTemplate string) {
	searchCount := strconv.Itoa(len(res))
	placeholderConfig := b.config.Telegram.Results.Template.Placeholder
	searchResults := ""

	// Ranges
	for _, r := range res {
		// Adds file name & buttons
		randHexColor, _ := utils.GetRandomHexValue(6)
		searchRes := strings.ReplaceAll(b.config.Telegram.Results.Template.SearchResults, placeholderConfig.HexColor, randHexColor)
		searchRes = strings.ReplaceAll(searchRes, placeholderConfig.GDriveDownloadLink, r.GDriveDownloadLink)
		searchRes = strings.ReplaceAll(searchRes, placeholderConfig.FileName, r.FileName)
		searchRes = strings.ReplaceAll(searchRes, placeholderConfig.FileSize, r.FileSize)
		searchRes = strings.ReplaceAll(searchRes, placeholderConfig.DownloadLink1, r.DownloadLink1)
		searchRes = strings.ReplaceAll(searchRes, placeholderConfig.DownloadLink2, r.DownloadLink2)

		// Adds media player links
		if strings.Contains(strings.ToLower(r.MimeType), "video") {
			mxPlayerLink := strings.ReplaceAll(b.config.Telegram.Results.Template.MXPlayer, placeholderConfig.MediaLink, r.MediaLink)
			mxPlayerLink = strings.ReplaceAll(mxPlayerLink, placeholderConfig.FileName, r.FileName)
			playersLink := strings.ReplaceAll(b.config.Telegram.Results.Template.OtherPlayer, placeholderConfig.MediaLink, r.MediaLink)
			playersLink = strings.ReplaceAll(playersLink, placeholderConfig.MXPlayerLink, mxPlayerLink)
			searchRes += playersLink
		}

		// Adds close tag
		searchRes += "</div></div>"
		searchResults += searchRes
	}

	// Gets template string
	responseTemplate, _ := utils.GetTemplateString()

	// Forms response
	responseTemplate = strings.ReplaceAll(responseTemplate, placeholderConfig.Query, searchQuery)
	responseTemplate = strings.ReplaceAll(responseTemplate, placeholderConfig.Count, searchCount)
	responseTemplate = strings.ReplaceAll(responseTemplate, placeholderConfig.Results, searchResults)

	// Returns
	return responseTemplate
}
