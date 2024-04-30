package logger

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
)

type Detail map[string]interface{}

type StandardLogger struct {
	*slog.Logger
}

type Config struct {
	LogLevel      slog.Level
	GrayLogActive bool
	GrayLogServer string
	GrayLogStream string
}

type APILogStruct struct {
	URL    string
	File   string
	Line   string
	Detail Detail
}

type ErrorSource struct {
	File string
	Line string
}

// CreateLogger Creates a StandardLogger by given config.
// According to config it can create a Gelf writer or a simple text
func CreateLogger(cfg Config) *StandardLogger {
	if _, exists := sloggraylog.LogLevels[cfg.LogLevel]; !exists {
		panic(fmt.Errorf("invalid Graylog level %d", cfg.LogLevel))
	}

	var logger *slog.Logger

	loggerOptions := &slog.HandlerOptions{
		Level:     cfg.LogLevel,
		AddSource: true,
	}

	textHandler := slog.New(slog.NewTextHandler(os.Stdout, loggerOptions)).Handler()

	if !cfg.GrayLogActive {
		return &StandardLogger{
			slog.New(slogmulti.Fanout(textHandler)),
		}
	}

	// Setup Gelf writer.
	gelfWriter, err := gelf.NewWriter(cfg.GrayLogServer)

	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}

	graylogHandler := slog.New(
		sloggraylog.Option{Level: cfg.LogLevel, Writer: gelfWriter}.NewGraylogHandler(),
	).With("stream", cfg.GrayLogStream)

	logger = slog.New(
		slogmulti.Fanout(
			textHandler,
			graylogHandler.Handler(),
		))

	return &StandardLogger{
		logger,
	}
}

func (logger *StandardLogger) LogAPIError(path string, err error, source *ErrorSource, detail Detail) {
	data := APILogStruct{
		File:   source.File,
		Line:   source.Line,
		URL:    path,
		Detail: detail,
	}

	logger.With(slog.Any("source", data)).Error(err.Error())
}

func (logger *StandardLogger) InfoF(format string, args ...any) {
	logger.Info(fmt.Sprintf(format, args...))
}

func (logger *StandardLogger) FatalF(format string, args ...any) {
	logger.Error(fmt.Sprintf(format, args...))

	os.Exit(1)
}

func (logger *StandardLogger) WarnF(format string, args ...any) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func GetErrorSource() *ErrorSource {
	_, file, line, _ := runtime.Caller(1)

	return &ErrorSource{
		File: file,
		Line: strconv.Itoa(line),
	}
}
