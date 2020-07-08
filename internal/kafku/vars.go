package kafku

import (
	//"github.com/WSBenson/goku/internal/kafka/utils"
	//"bitbucket.di2e.net/dime/wisard-nlp/utils"

	"github.com/WSBenson/goku/internal"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Vars is a list of key value pairs
type Vars map[string]string

// Keys
const (
	Port                        = "PORT"
	DocsFile                    = "DOCS_FILE"
	KafkaTopics                 = "KAFKA_TOPICS"
	KafkaBrokers                = "KAFKA_BROKERS"
	KafkaClientID               = "KAFKA_CLIENT_ID"
	KafkaGroupID                = "KAFKA_GROUP_ID"
	RevealVars                  = "REVEAL_VARS"
	GMDataAddr                  = "GM_DATA_ADDR"
	GmDataDerivativeAddr        = "GM_DATA_DERIVATIVE_ADDR"
	MimeTopics                  = "MIME_TOPICS"
	FallbackTopic               = "FALLBACK_TOPIC"
	SpacyURL                    = "SPACY_URL"
	SpacyBodySize               = "SPACY_BODY_SIZE"
	CogitoURL                   = "COGITO_URL"
	CogitoBodySize              = "COGITO_BODY_SIZE"
	InternalServerName          = "INTERNAL_SERVER_NAME"
	InternalServerCrtPath       = "INTERNAL_SERVER_CRT_PATH"
	InternalServerKeyPath       = "INTERNAL_SERVER_KEY_PATH"
	InternalCaCrtPath           = "INTERNAL_CA_CRT_PATH"
	InternalTLSOn               = "INTERNAL_TLS_ON"
	ExternalSourceServerName    = "EXTERNAL_SOURCE_SERVER_NAME"
	ExternalSourceServerCrtPath = "EXTERNAL_SOURCE_SERVER_CRT_PATH"
	ExternalSourceServerKeyPath = "EXTERNAL_SOURCE_SERVER_KEY_PATH"
	ExternalSourceCaCrtPath     = "EXTERNAL_SOURCE_CA_CRT_PATH"
	ChangelogPath               = "CHANGELOG_PATH"
	IndexName                   = "INDEX_NAME"
	DocType                     = "DOCTYPE"
	Derivtype                   = "DERIVTYPE"
	DType                       = "DTYPE"
	SpacyEndpoint               = "SPACY_ENDPOINT"
	RowDtype                    = "ROW_DTYPE"
)

// DefaultVars returns a map of default config vars to be used by viper
func DefaultVars() (defaultVars Vars) {
	return DefaultConfigVars
}

// NewVars return new config variables
func NewVars(defaultVars map[string]string) (v Vars) {
	// initialize logger
	//l := utils.NewLogger()

	// Set var value lookups
	viper.SetConfigName("config")
	//viper.AddConfigPath(".")
	//viper.AddConfigPath("../")
	viper.AddConfigPath("/etc/goku/")
	viper.AutomaticEnv()

	// Attempt to set vars
	if err := viper.ReadInConfig(); err != nil {
		internal.Logger.Fatal().Err(err).Msgf("Error loading in config file: ")
		errors.Wrap(err, "viper.ReadInConfig")
	} else {
		for k := range defaultVars {
			if viper.IsSet(k) {
				defaultVars[k] = viper.GetString(k)
			}
		}
	}
	return defaultVars
}
