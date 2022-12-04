package telegram

import (
	"errors"
	"net/http"
)

const (
	SearchErrUserMsg     = "𝕀𝕟𝕧𝕒𝕝𝕚𝕕 𝕤𝕖𝕒𝕣𝕔𝕙 𝕢𝕦𝕖𝕣𝕪 😞"
	SearchExampleUserMsg = "/𝕤𝕣𝕔𝕙 𝕕𝕠𝕔𝕥𝕠𝕣 𝕤𝕥𝕣𝕒𝕟𝕘𝕖 𝟚𝟘𝟚𝟚"
	InvalidCmdErrUserMsg = "𝕀𝕟𝕧𝕒𝕝𝕚𝕕 𝕔𝕠𝕞𝕞𝕒𝕟𝕕 😞"
)

var (
	InvalidCmdErrMsg = "Invalid Command"
	BadRequestErrMsg = "Bad Request"
	ErrCustom        = errors.New("custom err")
)

func (b *Bot) getCustomErrorMsg(statusCode int, err error) (userErrMsg, logErrMsg interface{}) {
	switch statusCode {
	case http.StatusBadRequest:
		{
			logTemplate := b.config.Telegram.Template.Log
			msg := logTemplate.Tag.Open + "\n" + logTemplate.Error + SearchErrUserMsg + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Example + SearchExampleUserMsg + "\n" + logTemplate.Tag.Close

			// Forms msg
			userErrMsg = msg
			logErrMsg = BadRequestErrMsg
		}
	case http.StatusInternalServerError:
		{

		}
	case http.StatusNotFound:
		{
			logTemplate := b.config.Telegram.Template.Log
			msg := logTemplate.Tag.Open + "\n" + logTemplate.Error + InvalidCmdErrUserMsg + "\n" + logTemplate.Tag.Close

			// Forms msg
			userErrMsg = msg
			logErrMsg = InvalidCmdErrMsg
		}
	}

	// Returns
	return userErrMsg, logErrMsg
}
