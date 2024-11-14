package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vickleford/gcsps/internal/gcs"
	"google.golang.org/api/option"
	pubsub "google.golang.org/api/pubsub/v1"
)

var listSubscriptionsCmd = &cobra.Command{
	Use:  "subscriptions [PROJECT] [TOPIC]",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		topic := args[1]

		svc, err := pubsub.NewService(context.TODO(),
			option.WithEndpoint(endpoint),
			option.WithoutAuthentication(),
		)
		if err != nil {
			return err
		}

		client := gcs.New(project, svc)

		subs, err := client.ListSubscriptions(context.TODO(), topic)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(subs, "\n"))

		return nil
	},
}
