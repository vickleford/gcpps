package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var drainCmd = &cobra.Command{
	Use:  "drain [PROJECT] [SUBSCRIPTION]",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		subscription := args[1]

		client, err := gcpClient(endpoint, project)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "draining from project %s on subscription %s\n", project, subscription)

		if err := client.Drain(context.TODO(), subscription); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "drained messages successfully")

		return nil
	},
}
