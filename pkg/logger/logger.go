package logger

import (
	"log"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() *zerolog.Logger {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(file).With().Timestamp().Logger()

	return &logger
}
