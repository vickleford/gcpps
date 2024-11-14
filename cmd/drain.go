package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vickleford/gcsps/internal/gcs"
	"google.golang.org/api/option"
	pubsub "google.golang.org/api/pubsub/v1"
)

var drainCmd = &cobra.Command{
	Use:  "drain [PROJECT] [SUBSCRIPTION]",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		subscription := args[1]

		svc, err := pubsub.NewService(context.TODO(),
			option.WithEndpoint(endpoint),
			option.WithoutAuthentication(),
		)
		if err != nil {
			return err
		}

		client := gcs.New(project, svc)

		fmt.Fprintf(cmd.OutOrStdout(), "draining from project %s on subscription %s\n", project, subscription)

		if err := client.Drain(context.TODO(), subscription); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "drained messages successfully")

		return nil
	},
}
