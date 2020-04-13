package slack

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"ping-and-report/repo"
	"strings"
	"time"
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

	_, err = (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	return nil
}

var recordRepo = repo.GetRecordRepository()

func convertVariables(s, addr string) string {
	r := recordRepo.GetOrNewRecord(addr)

	avabPercent := fmt.Sprintf("%.1f%%",
		float32(r.CountSucceed)/float32(math.Max(1, float64(r.CountTrying)))*100.0)

	var uptimeText string
	if r.LastAchieved {
		uptimeText = transcribeDuration(time.Since(r.LastBootAt))
	} else {
		uptimeText = "NOT_RUNNING_NOW"
	}

	return strings.NewReplacer(
		"$address$", addr,
		"$available_percent$", avabPercent,
		"$up_time$", uptimeText,
		"$total_measured_time$", transcribeDuration(time.Since(r.FirstBootAt)),
		"\r\n", "\n",
		"\r", "\n",
	).Replace(s)
}

func transcribeDuration(d time.Duration) string {
	return fmt.Sprintf("%.0fH %dM", math.Floor(d.Hours()), int64(math.RoundToEven(d.Minutes()))%60)
}
