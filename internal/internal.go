package internal

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger ....
var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()
}
