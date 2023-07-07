package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alancesar/imgur-fetcher/pkg/imgur"
	"github.com/alancesar/imgur-fetcher/pkg/media"
	"github.com/alancesar/imgur-fetcher/pkg/status"
	"github.com/alancesar/imgur-fetcher/pkg/transport"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	defaultClient := &http.Client{
		Transport: transport.NewUserAgentRoundTripper("imgur-fetcher", http.DefaultTransport),
	}

	imgurAuthClient := &http.Client{
		Transport: transport.NewAuthorizationRoundTripper(func(_ context.Context) (string, error) {
			return "Client-ID " + os.Getenv("IMGUR_CLIENT_ID"), nil
		}, defaultClient.Transport),
	}

	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalln("failed to start amqp connection:", err)
	}

	defer func() {
		_ = amqpConnection.Close()
	}()

	subscriber, err := amqpConnection.Channel()
	if err != nil {
		log.Fatalln("failed to start amqp channel")
	}

	defer func() {
		_ = subscriber.Close()
	}()

	publisher, err := amqpConnection.Channel()
	if err != nil {
		log.Fatalln("failed to start amqp channel")
	}

	//downloadsPublisher, err := pubsub.NewRabbitMQPublisher(amqpConnection, "media", "downloads")
	//if err != nil {
	//	log.Fatalln("failed to start media.downloads publisher")
	//}
	//
	//defer func() {
	//	_ = downloadsPublisher.Close()
	//}()

	imgurClient := imgur.NewClient(imgurAuthClient)

	messages, err := subscriber.Consume(
		"fetcher.imgur",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("failed to start fetcher.imgur consumer:", err)
	}

	consumer := func(req media.Media) error {
		mediaList, err := imgurClient.GetMediaByURL(req.URL)
		if err != nil {
			if errors.Is(err, status.ErrNotFound) {
				return nil
			}

			return fmt.Errorf("failed to retrieve media: %w", err)
		}

		for _, m := range mediaList {
			body, err := json.Marshal(media.Media{
				URL:    m.HigherQualityURL(),
				Parent: req.Parent,
			})
			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}

			if err := publisher.PublishWithContext(
				ctx,
				"media",
				"downloads",
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				},
			); err != nil {
				return fmt.Errorf("failed to publish media: %w", err)
			}

			//if err := downloadsPublisher.Publish(ctx, m); err != nil {
			//	return fmt.Errorf("failed to publish media: %w", err)
			//}
		}

		return nil
	}

	go func() {
		for message := range messages {
			var m media.Media
			if err := json.Unmarshal(message.Body, &m); err != nil {
				fmt.Println("failed to unmarshal message")
				_ = message.Ack(false)
				continue
			}

			if err := consumer(m); err != nil {
				fmt.Println("failed to handle message")
				_ = message.Nack(false, true)
			} else {
				_ = message.Ack(false)
			}
		}
	}()

	fmt.Println("all systems go!")

	<-ctx.Done()
	stop()

	fmt.Println("shutting down...")
	fmt.Println("good bye")
}
