package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sparrc/go-ping"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ping-and-report/record"
	"strconv"
	"time"
)

type config struct {
	Slack struct {
		WebHookURL string
		Mention    string
	}
	Pinged []string
}

const configYml = "config.yml"
const recordJson = "record.json"

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get this executable")
	}
	dir := filepath.Dir(exe)

	cfgp := dir + "/" + configYml
	var cfg config
	{
		cfg = config{}

		buf, err := ioutil.ReadFile(cfgp)
		if err != nil {
			log.Fatal("Failed to read " + cfgp)
		}
		err = yaml.Unmarshal(buf, &cfg)
		if err != nil {
			log.Fatal("Failed to parse " + cfgp)
		}
	}

	arcdp := dir + "/" + recordJson
	var arcd record.AvailableRecord
	{
		arcd = record.AvailableRecord{}

		if fileExists(arcdp) {
			buf, err := ioutil.ReadFile(arcdp)
			if err != nil {
				log.Fatal("Failed to read " + arcdp)
			}
			err = json.Unmarshal(buf, &arcd)
			if err != nil {
				log.Fatal("Failed to parse " + arcdp)
			}
		}
	}

	if !verifyConnection() {
		log.Fatal("Couldn't establish a connection.")
		return
	}

	c := make(chan pingResult)
	for _, addr := range cfg.Pinged {
		go sendPing(addr, c)
	}
	for range cfg.Pinged {
		res := <-c
		addr := res.Address
		nowAvab := res.IsAvailable

		log.Println("Sent a ping to " + addr + ": " + strconv.FormatBool(nowAvab))

		switched := arcd.Write(addr, nowAvab)
		rcd := arcd.Record(addr)

		if switched {
			percent := float32(rcd.CountSucceed) / float32(rcd.CountTrying) * 100.0

			var msg string
			if nowAvab {
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

	jsonBytes, err := json.Marshal(arcd)
	if err != nil {
		log.Fatal("Failed to parse struct")
	}
	err = ioutil.WriteFile(arcdp, jsonBytes, 0666)
	if err != nil {
		log.Fatal("Failed to write json")
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
