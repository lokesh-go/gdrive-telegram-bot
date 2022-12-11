package telegram

import (
	"fmt"
	"gdrive-telegram-bot/src/utils"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/showwin/speedtest-go/speedtest"
)

// Router ...
func (b *Bot) router(reqMsg string) (res interface{}, cmd string, statusCode int, err error) {
	// Gets key details
	cmd, text := getSearchKeyDetails(reqMsg)

	// Commands cases
	switch cmd {
	case SearchCmd:
		{
			// Checks
			if text == "" {
				return nil, Err, http.StatusBadRequest, ErrCustom
			}

			// Search file
			res, err = b.gdriveModule.FileSearch(text)
			if err != nil {
				return nil, Err, http.StatusInternalServerError, err
			}
		}
	case DriveCmd:
		{
			// Gets drive counts
		}
	case InfoCmd:
		{
			// Gets system information
			res = b.getSystemInformation()
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
	req = strings.ToLower(req)

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

	// Checks for info
	if strings.HasPrefix(req, InfoCmd) {
		// Assigns search cmd
		searchCmd = InfoCmd
	}

	// Returns
	return searchCmd, searchText
}

func (b *Bot) getSystemInformation() (res interface{}) {
	// Gets cpu
	cpu := runtime.NumCPU()

	// Gets memory
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	// Gets current no of goroutines
	numGoRoutine := runtime.NumGoroutine()

	// Gets speed test
	latency, downSpeed, uploadSpeed := speedTest()

	// Gets template
	logTemplate := b.config.Telegram.Template.Log
	resText := logTemplate.Tag.Open + "\n" + logTemplate.BotName + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.CPU + strconv.Itoa(cpu) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.RAM + utils.ConvertBytesToHumanReadableForm(mem.TotalAlloc) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Goroutine + strconv.Itoa(numGoRoutine) + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.DownSpeed + downSpeed + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.UploadSpeed + uploadSpeed + "\n" + logTemplate.Tag.Mid + "\n" + logTemplate.Latency + latency + "\n" + logTemplate.Tag.Close

	// Assigns
	res = resText

	// Returns
	return res
}

func speedTest() (latency, download, upload string) {
	user, _ := speedtest.FetchUserInfo()
	serverList, _ := speedtest.FetchServers(user)
	targets, _ := serverList.FindServer([]int{})
	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)
		latency = fmt.Sprintf("%d 𝕞𝕤", s.Latency.Milliseconds())
		download = fmt.Sprintf("%0.2f 𝕄𝕓𝕡𝕤", s.DLSpeed)
		upload = fmt.Sprintf("%0.2f 𝕄𝕓𝕡𝕤", s.ULSpeed)
	}

	// Returns
	return latency, download, upload
}
