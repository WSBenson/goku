package kafku

import (
	"strings"
	"time"

	"github.com/WSBenson/goku/internal"

	//"github.com/WSBenson/goku/internal/kafka/utils"
	//"github.com/WSBenson/goku/internal/kafka/vars"

	//"bitbucket.di2e.net/dime/wisard-nlp/utils"
	//"bitbucket.di2e.net/dime/wisard-nlp/vars"
	"github.com/Shopify/sarama"
)

// ProduceZ ...
func ProduceZ(zfighter string) {

	//v := NewVars(DefaultVars())

	clientID := RandomIndexName(10)
	internal.Logger.Info().Msgf("Initialize config with clientID %s", clientID)
	config := NewDefaultConfig(clientID)
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// change to: "kafka:29092" to run in docker
	addrs := strings.Split("localhost:9092", ",")
	internal.Logger.Info().Msgf("Initializing producer with addresses %v", addrs)
	//t.Logf("Initializing producer with addresses %v", addrs)
	producer, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		internal.Logger.Fatal().Msg("Error initializing producer..." + err.Error())
	}
	defer producer.AsyncClose()
	go func() {
		for err := range producer.Errors() {
			internal.Logger.Error().Msg("AsyncProducer Error:" + err.Error())
		}
	}()

	total := 1
	topics := RandomKafkaTopics(1, 10)
	internal.Logger.Info().Msgf("Producing %d messages on topic: %s", total, topics[0])

	// publishes an event message when a fighter has been added to the ES index
	PublishFighterMessage(producer, total, topics[0], zfighter)

	internal.Logger.Info().Msgf("Handling topic %s until %d messages are consumed", topics[0], total)
	var count int
	timer := time.NewTimer(10 * time.Second)

	groupID := RandomIndexName(10)
	internal.Logger.Info().Msg("Initializing new consumer group")
	group, err := NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		internal.Logger.Error().Msgf("Error creating a new consumer group: %s", err.Error())
	}

	onConsume := func(_ *sarama.ConsumerMessage) error {
		count++
		if count >= total {
			timer.Reset(0)
		}
		return nil
	}

	internal.Logger.Info().Msg("Handling messages")
	go func() {
		err = HandleTopics(onConsume, group, topics)
		if err != nil {
			internal.Logger.Error().Msgf("Error handling topics: %s", err.Error())
		}
	}()

	<-timer.C
	internal.Logger.Info().Msgf("Consumed %d Z fighter", count)
}
