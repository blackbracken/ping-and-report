package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"ping-and-report/repo"
	"ping-and-report/slack"
	"time"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show stats to the given channel of slack",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := repo.GetConfigRepository()
		recordRepo := repo.GetRecordRepository()

		for _, r := range recordRepo.Records {
			addr := r.Address

			if contains(cfg.Destinations, addr) {
				err := slack.ReportStats(addr)
				if err != nil {
					log.Fatal("Failed to report stats")
				}

				time.Sleep(1 * time.Second)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func contains(s []string, e string) bool {
	for _, x := range s {
		if x == e {
			return true
		}
	}

	return false
}
