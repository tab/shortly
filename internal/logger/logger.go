package logger

import (
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	LogLevel = 1
)

var (
	log  zerolog.Logger
	once sync.Once
)

func GetLogger() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFieldFormat,
		}

		log = zerolog.New(output).
			Level(zerolog.Level(LogLevel)).
			With().
			Timestamp().
			Logger()
	})

	return log
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := GetLogger()

		defer func() {
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Dur("duration", time.Since(start)).
				Int64("request_size", r.ContentLength).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Str("host", r.Host).
				Str("request_uri", r.RequestURI).
				Str("content_type", r.Header.Get("Content-Type")).
				Int64("content_length", r.ContentLength).
				Msg("request")
		}()

		next.ServeHTTP(w, r)
	})
}
