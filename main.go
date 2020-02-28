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
	"time"
)

type Config struct {
	Slack struct {
		WebHookURL string
		Mention    string
	}
	Pinged []string
}

type Available struct {
	AddressAvailables map[string]AddressAvailable
}

type AddressAvailable struct {
	CountTrying   uint64
	CountSucceed  uint64
	LastAvailable bool
}

const CONFIG = "config.yml"
const AVAILABLE = "available.json"

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get this executable")
	}
	dir := filepath.Dir(exe)

	var cfg Config
	cfgp := dir + "/" + CONFIG
	{
		cfg = Config{}

		buf, err := ioutil.ReadFile(cfgp)
		if err != nil {
			log.Fatal("Failed to read " + cfgp)
		}
		err = yaml.Unmarshal(buf, &cfg)
		if err != nil {
			log.Fatal("Failed to parse " + cfgp)
		}
	}

	var avb Available
	avbp := dir + "/" + AVAILABLE
	{
		avb = Available{map[string]AddressAvailable{}}

		if exists(avbp) {
			buf, err := ioutil.ReadFile(avbp)
			if err != nil {
				log.Fatal("Failed to read " + avbp)
			}
			err = json.Unmarshal(buf, &avb)
			if err != nil {
				log.Fatal("Failed to parse " + avbp)
			}
		}
	}

	for _, addr := range cfg.Pinged {
		pinger, err := ping.NewPinger(addr)
		if err != nil {
			log.Fatal("Failed to send a ping to " + addr)
		}

		pinger.Count = 3
		pinger.OnFinish = func(s *ping.Statistics) {
			addravb, ok := avb.AddressAvailables[s.Addr]
			if !ok {
				addravb = AddressAvailable{0, 0, true}
			}

			suc := s.PacketsRecv > 0

			addravb.CountTrying++
			addravb.LastAvailable = suc

			var msg string
			if suc {
				addravb.CountSucceed++
				msg = "UP  : " + s.Addr
			} else {
				msg = "DOWN: " + s.Addr
			}

			err := report(cfg.Slack.WebHookURL, cfg.Slack.Mention, msg)
			if err != nil {
				log.Fatal("Failed to report")
			}

			avb.AddressAvailables[s.Addr] = addravb
		}
		pinger.Timeout = 10 * time.Second
		pinger.Run()
	}

	fmt.Println(avb)

	jsonBytes, err := json.Marshal(avb)
	if err != nil {
		log.Fatal("Failed to parse struct")
	}
	err = ioutil.WriteFile(avbp, jsonBytes, 0666)
	if err != nil {
		log.Fatal("Failed to write json")
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
