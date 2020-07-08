package kafku

import (
	"github.com/Shopify/sarama"
	"github.com/WSBenson/goku/internal"
)

// OnConsumeFunc is a func type that handles consumer messages received from kafka
type OnConsumeFunc func(*sarama.ConsumerMessage) error

// GroupHandler implements the ConsumerGroupHandler interface in order to consume Kafka messages
type GroupHandler struct {
	onConsume OnConsumeFunc
	ready     chan bool
}

// Setup runs when the GroupHandler starts consuming messages
func (h *GroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	//l := utils.NewLogger()
	internal.Logger.Info().Msg("Setting up Group Handler...")
	internal.Logger.Debug().Msgf("Member ID:\t%s", session.MemberID())
	close(h.ready)
	return nil
}

// Cleanup runs when the GroupHandler is done
func (*GroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	//l := utils.NewLogger()
	internal.Logger.Info().Msg("Cleaning up Group Handler...")
	internal.Logger.Debug().Msgf("Member ID:\t%s", session.MemberID())
	return nil
}

// ConsumeClaim runs a consume func that is passed through the GroupHandler every time a message is pulled from claim
func (h *GroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	//l := utils.NewLogger()
	for msg := range claim.Messages() {
		sess.MarkMessage(msg, "")
		err = h.onConsume(msg)
		if err != nil {
			internal.Logger.Error().Msgf("Error consuming Kafka message: %s", err.Error())
			return
		}
	}
	return
}
