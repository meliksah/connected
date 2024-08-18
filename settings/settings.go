package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	ServerPort int    `json:"server_port"`
	ServerIP   string `json:"server_ip"`
	Password   string `json:"password"`
	LastIP     string `json:"last_ip"`
	LastPort   int    `json:"last_port"`
}

var (
	settingsPath = filepath.Join(os.Getenv("HOME"), ".connected", "settings.json")
	settings     Settings
)

const MagicWord = "PASSWORD_TEXT"

func LoadSettings() error {
	file, err := os.Open(settingsPath)
	if os.IsNotExist(err) {
		settings = Settings{ServerPort: 59599, Password: "", LastIP: "192.168.50.50", LastPort: 59599}
		return SaveSettings()
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	return err
}

func SaveSettings() error {
	err := os.MkdirAll(filepath.Dir(settingsPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(settingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(&settings)
}

func GetPort() int {
	return settings.ServerPort
}

func GetPassword() string {
	return settings.Password
}

func GetLastIP() string {
	return settings.LastIP
}

func GetLastPort() int {
	return settings.LastPort
}

func GetLocalIP() string {
	return settings.ServerIP
}

func SetPort(port int) {
	settings.ServerPort = port
}

func SetPassword(password string) {
	settings.Password = password
}

func SetLastIP(lastIP string) {
	settings.LastIP = lastIP
}

func SetLastPort(lastPort int) {
	settings.LastPort = lastPort
}
