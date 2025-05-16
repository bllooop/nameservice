package log

import (
	"io"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger(output io.Writer, level zerolog.Level) {
	Logger = zerolog.New(output).Level(level).With().Timestamp().Logger()
}
