package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ErrorSettings struct {
	Errors map[string]string `json:"errors"`
}

var (
	errorSettings  ErrorSettings
	errorsFilePath = filepath.Join("resources", "errors.json")
)

func LoadErrors() error {
	file, err := os.Open(errorsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&errorSettings)
	if err != nil {
		return err
	}

	return nil
}

func GetErrorMessage(code string) string {
	if msg, exists := errorSettings.Errors[code]; exists {
		return msg
	}
	return "Unknown error occurred."
}
