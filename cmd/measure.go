package cmd

import (
	"github.com/sparrc/go-ping"
	"github.com/spf13/cobra"
	"log"
	"ping-and-report/repo"
	"ping-and-report/slack"
	"strconv"
	"time"
)

var measureCmd = &cobra.Command{
	Use:   "measure",
	Short: "Send ping and report the result the given channel of slack",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := repo.GetConfigRepository()
		recordRepo := repo.GetRecordRepository()

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

			switched := recordRepo.Record(addr, achieved)
			if switched {
				if achieved {
					// down -> up
					err := slack.ReportServerUp(addr)
					if err != nil {
						log.Fatal("Failed to send a message")
					}
				} else {
					// up -> down
					err := slack.ReportServerDown(addr)
					if err != nil {
						log.Fatal("Failed to send a message")
					}
				}
			}
		}

		err := recordRepo.Flush()
		if err != nil {
			log.Fatal("Failed to write records a file")
		}
	},
}

func init() {
	rootCmd.AddCommand(measureCmd)
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
