package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vickleford/gcpps/internal/gcp"
	"google.golang.org/api/option"
	pubsub "google.golang.org/api/pubsub/v1"
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

func gcpClient(endpoint, project string) (*gcp.Client, error) {
	svc, err := pubsub.NewService(context.TODO(),
		option.WithEndpoint(endpoint),
		option.WithoutAuthentication(),
	)
	if err != nil {
		return nil, err
	}

	client := gcp.New(project, svc)

	return client, nil
}
