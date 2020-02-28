package main

import (
	"bytes"
	"github.com/sparrc/go-ping"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Slack struct {
		WebHookURL string
		Mention    string
	}
	Pinged []string
}

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get this executable")
	}
	dir := filepath.Dir(exe)

	buf, err := ioutil.ReadFile(dir + "/config.yml")
	if err != nil {
		log.Fatal("Failed to get config.yml")
	}

	cfg := Config{}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatal("Failed to parse config.yml")
	}

	for _, addr := range cfg.Pinged {
		pinger, err := ping.NewPinger(addr)
		if err != nil {
			log.Fatal("Failed to send a ping to " + addr)
		}

		pinger.Count = 3
		pinger.OnFinish = func(s *ping.Statistics) {
			var msg string
			if s.PacketsRecv == 0 {
				msg = "DOWN: " + s.Addr
			} else {
				msg = "UP  : " + s.Addr
			}

			err := report(cfg.Slack.WebHookURL, cfg.Slack.Mention, msg)
			if err != nil {
				log.Fatal("Failed to report")
			}
		}
		pinger.Timeout = 10 * time.Second
		pinger.Run()
	}
}

func report(url string, mention string, text string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{\"text\":\""+mention+" "+text+"\"}")))
	if err != nil {
		return err
	}

	_, err = (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	return nil
}
