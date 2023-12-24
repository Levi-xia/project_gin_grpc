package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = &Configuration{}

type Config struct {
	viper *viper.Viper
}

type Configuration struct {
	Server ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

func InitConfig() *Config {
	config := &Config{
		viper: viper.New(),
	}
	workDir, _ := os.Getwd()
	config.viper.SetConfigName("app-dev")
	config.viper.SetConfigType("yaml")
	config.viper.AddConfigPath(workDir + "/config")

	if err := config.viper.ReadInConfig(); err != nil {
		log.Fatalln("Fatal error config file: ", err)
	}
	config.viper.WatchConfig()
    config.viper.OnConfigChange(func(in fsnotify.Event) {
        log.Println("config file changed:", in.Name)
        if err := config.viper.Unmarshal(Conf); err != nil {
            log.Println("Unmarshal config failed, err:", err)
        }
    })
    if err := config.viper.Unmarshal(Conf); err != nil {
        log.Println("Unmarshal config failed, err:", err)
    }
	return config
}
