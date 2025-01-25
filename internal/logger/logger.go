package logger

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// LogLevel is the log level for the logger. Default is Info level
const (
	LogLevel = 1
)

// Logger is a logger structure for the application
type Logger struct {
	log zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: zerolog.TimeFieldFormat,
	}

	log := zerolog.New(output).
		Level(zerolog.Level(LogLevel)).
		With().
		Timestamp().
		Logger()

	return &Logger{log: log}
}

func (l *Logger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.log.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.log.Error()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		l.Info().
			Str("method", r.Method).
			Str("content_type", r.Header.Get("Content-Type")).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Int64("request_size", r.ContentLength).
			Int("status", rw.statusCode).
			Int("response_size", rw.bytesWritten).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("request")
	})
}
