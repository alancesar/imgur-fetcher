package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alancesar/imgur-fetcher/internal/controller"
	"github.com/alancesar/imgur-fetcher/pkg/imgur"
	"github.com/alancesar/imgur-fetcher/pkg/pubsub"
	"github.com/alancesar/imgur-fetcher/pkg/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	publisher, err := pubsub.NewRabbitMQPublisher(amqpConnection, "fetcher", "imgur")
	if err != nil {
		log.Fatalln("failed to start fetcher publisher:", err)
	}

	imgurClient := imgur.NewClient(imgurAuthClient)
	imgurController := controller.New(defaultClient, imgurClient, publisher)

	mux := chi.NewMux()
	mux.Use(middleware.Logger, middleware.SetHeader("Content-Type", "application/json"))
	mux.Post("/", imgurController.GetMediaByURL)
	mux.Post("/publish", imgurController.PublishMedia)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + os.Getenv("PORT"),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Println("error on start http server:", err)
			panic(err)
		}
	}()

	fmt.Println("all systems go!")

	<-ctx.Done()
	stop()

	fmt.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	fmt.Println("good bye")
}
