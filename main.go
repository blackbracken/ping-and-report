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

type anyMap = map[interface{}]interface{}

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

	m := make(anyMap)
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		log.Fatal("Failed to parse config.yml")
	}

	slack := m["Slack"].(anyMap)
	webhook := slack["WebHookURL"].(string)
	mention := slack["Mention"].(string)

	fmt.Println(webhook)
	fmt.Println(mention)

	p := m["Pinged"].([]interface{})
	for i := range p {
		s := p[i].(string)
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
