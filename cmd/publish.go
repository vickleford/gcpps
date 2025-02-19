package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const flagContentType = "content-type"

var publishCmd = &cobra.Command{
	Use:  "publish [PROJECT] [TOPIC] [DATA]",
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		topic := args[1]
		data := args[2]

		client, err := gcpClient(endpoint, project)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "publishing to project %s on topic %s: %s\n", project, topic, data)

		contentType, err := cmd.Flags().GetString(flagContentType)
		if err != nil {
			return err
		}

		id, err := client.Publish(context.TODO(), topic, contentType, data)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "published message successfully with id %s\n", id)

		return nil
	},
}

func init() {
	publishCmd.Flags().String(flagContentType, "text/plain", "Set the content type of the data to publish")
}
