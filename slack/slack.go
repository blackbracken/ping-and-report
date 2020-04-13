package slack

import (
	"bytes"
	"fmt"
	"net/http"
	"ping-and-report/repo"
	"strings"
)

var cfg = repo.GetConfigRepository()

func ReportServerUp(addr string) error {
	return sendMessage(convertVariables(cfg.Message.ServerUp, addr))
}

func ReportServerDown(addr string) error {
	return sendMessage(convertVariables(cfg.Message.ServerDown, addr))
}

func ReportStats(addr string) error {
	return sendMessage(convertVariables(cfg.Message.ServerStats, addr))
}

func sendMessage(s string) error {
	req, err := http.NewRequest("POST", cfg.Slack.WebHookURL, bytes.NewBuffer([]byte("{\"text\":\""+s+"\"}")))
	if err != nil {
		return err
	}

	removethis, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("LOG: %d\n", removethis.StatusCode)

	return nil
}

func convertVariables(s, addr string) string {
	return strings.NewReplacer(
		"$address$", addr,
		"$available_percent$", "DUMMY-100.0%",
		"$up_time$", "DUMMY-10:20",
		"$total_running_time$", "DUMMY-20:40",
		"\r\n", "\n",
		"\r", "\n",
	).Replace(s)
}
