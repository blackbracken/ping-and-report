package slack

import (
	"bytes"
	"net/http"
	"ping-and-report/repo"
)

func SendMessage(text string) error {
	cfg, err := repo.LoadConfig()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", cfg.Slack.WebHookURL, bytes.NewBuffer([]byte("{\"text\":\""+text+"\"}")))
	if err != nil {
		return err
	}

	_, err = (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	return nil
}
