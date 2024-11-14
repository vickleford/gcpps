package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listTopicsCmd = &cobra.Command{
	Use:  "topics [PROJECT]",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]

		client, err := gcpClient(endpoint, project)
		if err != nil {
			return err
		}

		topics, err := client.ListTopics(context.TODO())
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(topics, "\n"))

		return nil
	},
}
