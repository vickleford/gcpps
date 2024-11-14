package gcp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"google.golang.org/api/googleapi"
	pubsub "google.golang.org/api/pubsub/v1"
)

type Client struct {
	svc     *pubsub.Service
	project string
}

func New(project string, svc *pubsub.Service) *Client {
	return &Client{svc: svc, project: project}
}

func (c *Client) Publish(ctx context.Context, topic, data string) (string, error) {
	if err := c.createTopicIfNotExists(ctx, topic); err != nil {
		return "", err
	}

	topic = path.Join("projects", c.project, "topics", topic)

	encodedData := base64.StdEncoding.EncodeToString([]byte(data))

	call := c.svc.Projects.Topics.Publish(topic, &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			{
				Attributes: map[string]string{"Content-Type": "application/json"},
				Data:       encodedData,
			},
		},
	}).Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return "", err
	}

	if len(resp.MessageIds) != 1 {
		return "", fmt.Errorf("got %d message IDs", len(resp.MessageIds))
	}

	return resp.MessageIds[0], err
}

func (c *Client) createTopicIfNotExists(ctx context.Context, topic string) error {
	getCall := c.svc.Projects.Topics.Get(c.fqtp(topic)).Context(ctx)
	_, err := getCall.Do()
	var googleAPIError *googleapi.Error
	if errors.As(err, &googleAPIError) {
		if googleAPIError.Code != http.StatusNotFound {
			return err
		}
	} else if err != nil {
		return err
	}

	if err == nil {
		// Don't create it.
		return nil
	}

	call := c.svc.Projects.Topics.Create(c.fqtp(topic), &pubsub.Topic{}).Context(ctx)
	_, err = call.Do()
	return err
}

// fully qualified topic path
func (c *Client) fqtp(topic string) string {
	return path.Join("projects", c.project, "topics", topic)
}

type Message struct {
	ID         string
	Attributes map[string]string
	Data       string
}

type SubscribeEvent struct {
	Message Message
	Error   error
}

func (c *Client) Subscribe(ctx context.Context, topic, subscription string) (chan SubscribeEvent, error) {
	subscription = path.Join("projects", c.project, "subscriptions", subscription)

	// create a new subscription...
	call := c.svc.Projects.Subscriptions.Create(subscription, &pubsub.Subscription{
		Name:  subscription,
		Topic: c.fqtp(topic),
	}).Context(ctx)
	// call it...
	var googleAPIErr *googleapi.Error
	if _, err := call.Do(); errors.As(err, &googleAPIErr) {
		if googleAPIErr.Code != http.StatusConflict {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	events := make(chan SubscribeEvent, 1000)

	go func(ctx context.Context) {
		defer close(events)

		for {
			pull := c.svc.Projects.Subscriptions.Pull(subscription, &pubsub.PullRequest{
				MaxMessages:       100,
				ReturnImmediately: false,
			}).Context(ctx)

			resp, err := pull.Do()
			if err != nil {
				events <- SubscribeEvent{Error: err}
				return
			}

			ids := make([]string, 0)

			for _, msg := range resp.ReceivedMessages {
				decodedData, err := base64.StdEncoding.DecodeString(msg.Message.Data)
				if err != nil {
					events <- SubscribeEvent{Error: err}
					return
				}

				events <- SubscribeEvent{
					Message: Message{
						ID:         msg.Message.MessageId,
						Attributes: msg.Message.Attributes,
						Data:       string(decodedData),
					},
				}
				ids = append(ids, msg.AckId)
			}

			ack := c.svc.Projects.Subscriptions.Acknowledge(subscription, &pubsub.AcknowledgeRequest{
				AckIds: ids,
			}).Context(ctx)

			if _, err := ack.Do(); err != nil {
				events <- SubscribeEvent{Error: err}
				return
			}
		}
	}(ctx)

	return events, nil
}

func (c *Client) Drain(ctx context.Context, subscription string) error {
	subscription = path.Join("projects", c.project, "subscriptions", subscription)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	idsCh := make(chan []string)
	pullErr := make(chan error)

	go func() {
		var goroutineErr error

		for {
			if ctx.Err() != nil {
				// Technically this leaks the context error, but I'm not yet
				// sure if I need to preserve it.
				break
			}

			pull := c.svc.Projects.Subscriptions.Pull(subscription, &pubsub.PullRequest{
				MaxMessages:       100,
				ReturnImmediately: true,
			}).Context(ctx)

			resp, err := pull.Do()
			if err != nil {
				goroutineErr = err
				break
			}

			if len(resp.ReceivedMessages) == 0 {
				break
			}

			ids := make([]string, 0)
			for _, msg := range resp.ReceivedMessages {
				ids = append(ids, msg.AckId)
			}

			idsCh <- ids
		}

		close(idsCh)
		pullErr <- goroutineErr
	}()

	for ids := range idsCh {
		call := c.svc.Projects.Subscriptions.Acknowledge(subscription, &pubsub.AcknowledgeRequest{
			AckIds: ids,
		})

		_, err := call.Do()
		if err != nil {
			return fmt.Errorf("failed to ack messages (%s): %w", strings.Join(ids, ","), err)
		}
	}

	if err := <-pullErr; err != nil {
		return fmt.Errorf("failed pulling messages: %w", err)
	}

	return nil
}

func (c *Client) ListTopics(ctx context.Context) ([]string, error) {

	topics := make([]string, 0)

	var nextPage string
	for {
		call := c.svc.Projects.Topics.List(c.project).Context(ctx)
		call.PageToken(nextPage)

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		if resp.NextPageToken == "" {
			break
		}

		nextPage = resp.NextPageToken

		for _, topic := range resp.Topics {
			topics = append(topics, topic.Name)
		}
	}

	return topics, nil
}

func (c *Client) ListSubscriptions(ctx context.Context, topic string) ([]string, error) {
	call := c.svc.Projects.Topics.Subscriptions.List(topic).Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}
