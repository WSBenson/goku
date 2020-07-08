package kafku

import (
	"github.com/Shopify/sarama"
)

// NewDefaultConfig returns a Config pointer  to be used by a ConsumerGroup
func NewDefaultConfig(clientID string) (config *sarama.Config) {
	config = sarama.NewConfig()
	config.ClientID = clientID
	config.Version = sarama.V2_2_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return
}
