package telegram

import (
	"net/http"
	"strings"

	gdrive "gdrive-telegram-bot/src/gdrive"
)

// Router ...
func (b *Bot) router(reqMsg string) (res interface{}, cmd string, statusCode int, err error) {
	// Converts into lower
	req := strings.ToLower(reqMsg)

	// Gets key details
	cmd, text := getSearchKeyDetails(req)

	// Gets the gdrive module
	gdriveModule := gdrive.New(b.config)

	// Commands cases
	switch cmd {
	case SearchCmd:
		{
			// Checks
			if text == "" {
				return nil, Err, http.StatusBadRequest, ErrCustom
			}

			// Search file
			res, err = gdriveModule.FileSearch(text)
			if err != nil {
				return nil, Err, http.StatusInternalServerError, err
			}
		}
	case DriveCmd:
		{
			// Gets drive counts
		}
	default:
		{
			// Returns
			return nil, Err, http.StatusNotFound, ErrCustom
		}
	}

	// Returns
	return res, cmd, http.StatusOK, nil
}

func getSearchKeyDetails(req string) (searchCmd, searchText string) {
	// Checks for search
	if strings.HasPrefix(req, SearchCmd) {
		// Splits
		splitedStrings := strings.SplitN(req, " ", 2)

		// Assigns search text
		if len(splitedStrings) == 2 {
			searchText = splitedStrings[1]
		}

		// Assigns search cmd
		searchCmd = SearchCmd
	}

	// Returns
	return searchCmd, searchText
}
