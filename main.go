package main

import (
	"bytes"
	"fmt"
	"github.com/sparrc/go-ping"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to read a config")
	}

	repo, err := LoadRecordRepository()
	if err != nil {
		log.Fatal("Failed to read records")
	}

	if !verifyConnection() {
		log.Fatal("Couldn't establish a connection")
	}

	ch := make(chan pingResult)
	for _, addr := range cfg.Destinations {
		go sendPing(addr, ch)
	}
	for range cfg.Destinations {
		res := <-ch
		addr := res.Address
		achieved := res.IsAvailable

		log.Println("Sent a ping to " + addr + ": " + strconv.FormatBool(achieved))

		switched := repo.Record(addr, achieved)
		rcd := repo.GetRecord(addr)

		if switched {
			percent := float32(rcd.CountSucceed) / float32(rcd.CountTrying) * 100.0

			var msg string
			if achieved {
				// down -> up
				msg = ":signal_strength: The server " + addr + " is currently up! | available: " + fmt.Sprintf("%.1f%%", percent)
			} else {
				// up -> down
				msg = ":warning: The server " + addr + " is currently down! | available: " + fmt.Sprintf("%.1f%%", percent)
			}

			err := report(cfg.Slack.WebHookURL, cfg.Slack.Mention, msg)
			if err != nil {
				log.Fatal("Failed to report")
			}
		}
	}

	err = repo.Flush()
	if err != nil {
		log.Fatal("Failed to write records a file")
	}
}

type pingResult struct {
	Address     string
	IsAvailable bool
}

func verifyConnection() bool {
	p, err := ping.NewPinger("8.8.8.8")
	if err != nil {
		return false
	}

	p.Count = 4
	p.Timeout = 10 * time.Second
	p.Run()

	return p.Statistics().PacketsRecv > 0
}

func sendPing(addr string, c chan pingResult) {
	p, err := ping.NewPinger(addr)
	if err != nil {
		log.Fatal("Failed to send a ping to " + addr)
	}

	p.Count = 5
	p.Timeout = 20 * time.Second
	p.OnFinish = func(s *ping.Statistics) { c <- pingResult{addr, s.PacketsRecv > 0} }
	p.Run()
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
