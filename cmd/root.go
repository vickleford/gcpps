package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gcsps",
	Short: "GCS PubSub Utility",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	endpoint string
)

func init() {
	rootCmd.AddCommand(
		publishCmd,
		drainCmd,
		listTopicsCmd,
		listSubscriptionsCmd,
		subscribeCmd,
	)

	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "http://localhost:8085", "set the pubsub endpoint")
}
