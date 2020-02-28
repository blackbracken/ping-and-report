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
	"strconv"
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

	c := make(chan pingResult)
	for _, addr := range cfg.Pinged {
		go sendPing(addr, c)
	}
	for range cfg.Pinged {
		res := <-c
		addr := res.Address
		suc := res.IsAvailable

		log.Println("Sent a ping to " + addr + ": " + strconv.FormatBool(suc))

		addravb, ok := avb.AddressAvailables[addr]
		if !ok {
			addravb = AddressAvailable{0, 0, true}
		}

		if suc {
			addravb.CountSucceed++
		}
		addravb.CountTrying++
		// suc XOR last_suc
		if suc != addravb.LastAvailable {
			var percent float32
			if addravb.CountTrying == 0 {
				percent = 0
			} else {
				percent = float32(addravb.CountSucceed) / float32(addravb.CountTrying)
			}
			percent *= 100

			var msg string
			if suc {
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
		addravb.LastAvailable = suc

		avb.AddressAvailables[addr] = addravb
	}

	jsonBytes, err := json.Marshal(avb)
	if err != nil {
		log.Fatal("Failed to parse struct")
	}
	err = ioutil.WriteFile(avbp, jsonBytes, 0666)
	if err != nil {
		log.Fatal("Failed to write json")
	}
}

type pingResult struct {
	Address     string
	IsAvailable bool
}

func sendPing(addr string, c chan pingResult) {
	p, err := ping.NewPinger(addr)
	if err != nil {
		log.Fatal("Failed to send a ping to " + addr)
	}

	p.Count = 3
	p.Timeout = 10 * time.Second
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
