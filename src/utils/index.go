package utils

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"

	"github.com/dustin/go-humanize"
)

// ReadJSONFile ...
func ReadJSONFile(path string, instance interface{}) (err error) {
	// Forms bytes
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal
	err = json.Unmarshal(bytes, instance)
	if err != nil {
		return err
	}

	// Returns
	return nil
}

// ConvertBytesToHumanReadableForm ...
func ConvertBytesToHumanReadableForm(s uint64) (res string) {
	// Returns
	return humanize.Bytes(s)
}

// GetTemplateString ...
func GetTemplateString() (template string, err error) {
	// Forms bytes
	bytes, err := ioutil.ReadFile("utils/template.html")
	if err != nil {
		return template, err
	}

	// Returns
	return string(bytes), nil
}

// GetRandomHexValue ...
func GetRandomHexValue(length int) (hexVal string, err error) {
	// Gets random integer value
	randInt := rand.Intn(5000)

	bytes := make([]byte, randInt)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Gets hex string
	hexString := hex.EncodeToString(bytes)

	// Trim
	hexVal = hexString[len(hexString)-length:]

	// Returns
	return hexVal, nil
}
