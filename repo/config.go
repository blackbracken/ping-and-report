package repo

import "github.com/spf13/viper"

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
		Message: struct {
			ServerUp    string
			ServerDown  string
			ServerStats string
		}{
			viper.GetString("message.server_up"),
			viper.GetString("message.server_down"),
			viper.GetString("message.server_stats"),
		},
	}, nil
}
