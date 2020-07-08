package kafku

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	//"bitbucket.di2e.net/dime/wisard-nlp/utils"
	"github.com/Shopify/sarama"
	"github.com/WSBenson/goku/internal"
)

// NewConsumerGroup initializes a new sarama ConsumerGroup with the provided config values
func NewConsumerGroup(addrs []string, groupID string, config *sarama.Config) (group sarama.ConsumerGroup, err error) {
	//l := NewLogger()

	internal.Logger.Info().Msgf("Initializing new Kafka Consumer Group with Group ID %s and addresses %v", groupID, addrs)
	group, err = sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		internal.Logger.Fatal().Msgf("Error initializing a new Kafka Consumer Group: %s", err.Error())
		return
	}

	return
}

// HandleTopics consumes messages with a given function as part of a given consumer group
func HandleTopics(onConsume OnConsumeFunc, group sarama.ConsumerGroup, topics []string) (err error) {
	//l := NewLogger()

	internal.Logger.Info().Msgf("Handling Kafka Topics: %v", topics)
	handler := GroupHandler{
		ready:     make(chan bool, 0),
		onConsume: onConsume,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			if err := group.Consume(ctx, topics, &handler); err != nil {
				internal.Logger.Fatal().Msgf("Error from consumer: %v", err.Error())
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			handler.ready = make(chan bool, 0)
		}
	}()

	<-handler.ready
	internal.Logger.Info().Msg("Consumer is up and running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		internal.Logger.Info().Msg("Terminating: context cancelled")
	case <-sigterm:
		internal.Logger.Info().Msg("Terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = group.Close(); err != nil {
		internal.Logger.Error().Msgf("Error closing client: %v", err.Error())
		return
	}

	return
}
