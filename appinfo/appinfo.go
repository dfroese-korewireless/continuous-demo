package appinfo

import (
	"encoding/json"
	"os"
)

// AppInfo contains the information to be conveyed
type AppInfo struct {
	AppVersion, Port string
	DBPath           string `json:"DatabasePath"`
}

const (
	configFileName = "appsettings.json"
)

func openConfigFile() (*os.File, error) {
	return os.Open(configFileName)
}

// GetAppInfo returns the application information from the config file
func GetAppInfo() (*AppInfo, error) {
	file, err := openConfigFile()
	info := &AppInfo{}

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&info)

	if err != nil {
		return nil, err
	}

	return info, nil
}
