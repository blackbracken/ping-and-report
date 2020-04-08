package repo

import "github.com/spf13/viper"

type Config struct {
	Slack struct {
		WebHookURL string
		Mention    string
	}
	Destinations []string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Slack: struct {
			WebHookURL string
			Mention    string
		}{
			viper.GetString("slack.webhookurl"),
			viper.GetString("slack.mention"),
		},
		Destinations: viper.GetStringSlice("destinations"),
	}, nil
}
