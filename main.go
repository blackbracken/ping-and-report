package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	fmt.Println(filepath.Dir(exe))
	buf, err := ioutil.ReadFile(filepath.Dir(exe) + "/config.yml")
	if err != nil {
		log.Fatal("Failed to get config.yml")
	}

	cfg := Config{}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatal("Failed to parse config.yml")
	}

	slack := cfg.Slack
	webhook := slack.WebHookURL
	mention := slack.Mention

	fmt.Println("WH: " + webhook)
	fmt.Println("M : " + mention)

	p := cfg.Pinged
	for i := range p {
		s := p[i]
		fmt.Println(s)
	}

	pinger, err := ping.NewPinger("google.com")
	if err != nil {
		log.Fatal("Failed to send a ping google.com")
	}
	pinger.SetPrivileged(true)
	pinger.Count = 4
	pinger.Run()
}
