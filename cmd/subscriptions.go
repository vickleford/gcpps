package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listSubscriptionsCmd = &cobra.Command{
	Use:  "subscriptions [PROJECT] [TOPIC]",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		topic := args[1]

		client, err := gcpClient(endpoint, project)
		if err != nil {
			return err
		}

		subs, err := client.ListSubscriptions(context.TODO(), topic)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(subs, "\n"))

		return nil
	},
}
