package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

// OpenFile ...
func OpenFile(path string) (file *os.File, err error) {
	// Opens
	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// Return
	return file, nil
}
