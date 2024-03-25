package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path"
	"sync"
)

//var (
//	defaultConfigPath = "_config.yaml"
//)

func createPath() (string, error) {
	u, err := user.Current()
	fmt.Println("Home dir:", u.HomeDir)
	return u.HomeDir, err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createFile(path string) error {
	// 检查是否存在文件
	exists, existsErr := pathExists(path)
	if exists {
		return existsErr
	}

	f, createErr := os.Create(path)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)

	return createErr
}

type config struct {
}

var instance *config
var once sync.Once

func GetInstance() *config {
	if instance == nil {
		once.Do(func() {
			err := InitConfig()
			if err != nil {
				fmt.Println("init config error: ", err.Error())
				return
			}
		})
	}
	return instance
}

// InitConfig init config
func InitConfig() error {
	// 生成家目录
	defaultConfigPath, err := createPath()
	if err != nil {
		return err
	}
	defaultConfigPath = path.Join(defaultConfigPath, "spark_ai_cli_config.yaml")
	// 创建配置文件
	if err := createFile(defaultConfigPath); err != nil {
		return err
	}

	// 设置配置文件
	viper.SetConfigFile(defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 初始化基础配置
	viper.Set("SPARK.BASE", "wss://spark-api.xf-yun.com/v3.5/chat")
	viper.Set("SPARK.DOMAIN", "generalv3.5")

	instance = &config{}

	return nil
}

// Set config with key
func Set(key, value string) {
	_ = GetInstance()
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		return
	}
}

// Get config value with key
func Get(key string) any {
	_ = GetInstance()
	return viper.Get(key)
}
