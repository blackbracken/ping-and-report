package repo

import (
	"github.com/spf13/viper"
	"log"
)

var configRepo *Config

type Config struct {
	Slack struct {
		WebHookURL string
		Mention    string
	}
	Destinations []string
	Message      struct {
		ServerUp    string
		ServerDown  string
		ServerStats string
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read a config")
	}

	configRepo = &Config{
		Slack: struct {
			WebHookURL string
			Mention    string
		}{
			viper.GetString("slack.webhookurl"),
			viper.GetString("slack.mention"),
		},
		Destinations: viper.GetStringSlice("destinations"),
		Message: struct {
			ServerUp    string
			ServerDown  string
			ServerStats string
		}{
			viper.GetString("message.server_up"),
			viper.GetString("message.server_down"),
			viper.GetString("message.server_stats"),
		},
	}
}

func GetConfigRepository() *Config {
	return configRepo
}
