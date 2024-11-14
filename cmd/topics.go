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

var listTopicsCmd = &cobra.Command{
	Use:  "topics [PROJECT]",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]

		svc, err := pubsub.NewService(context.TODO(),
			option.WithEndpoint(endpoint),
			option.WithoutAuthentication(),
		)
		if err != nil {
			return err
		}

		client := gcs.New(project, svc)

		topics, err := client.ListTopics(context.TODO())
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(topics, "\n"))

		return nil
	},
}
