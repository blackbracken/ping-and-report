package main

import (
	"github.com/sparrc/go-ping"
)

func main() {
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		panic(err)
	}

	pinger.Run()
}
