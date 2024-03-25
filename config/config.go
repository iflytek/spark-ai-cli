package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	defaultConfigPath = ".sparkai/config.yaml"
)

type config struct {
}

var instance *config
var once sync.Once

func GetInstance() *config {
	if instance == nil {
		once.Do(func() {
			instance = &config{}
		})
	}
	return instance
}

func fileFullPath(relativeFilePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %s", err)
	}
	configFilePath := filepath.Join(homeDir, relativeFilePath)
	return configFilePath, nil
}

// InitConfig init config
func InitConfig() error {

	return nil
}

// Set config with key
func Set(key, value string) {

}

// Get config value with key
func Get(key string) error {

	return nil
}
