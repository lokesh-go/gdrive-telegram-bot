package utils

import (
	"encoding/json"
	"io/ioutil"

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
