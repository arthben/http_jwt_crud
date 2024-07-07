package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

const KeyLoggerObj = "objLogger"

func NewLogger(appMode string, serviceName string) *slog.Logger {
	var (
		writter io.Writer
		opts    slog.HandlerOptions
	)

	if appMode != "production" {
		writter = os.Stdout
	} else {
		writter = &lumberjack.Logger{
			Filename: "log/mylog.log",
			Compress: true,
		}

		replacer := func(groups []string, attr slog.Attr) slog.Attr {
			// when addSource is true, log will printed full qualified path name
			// replacer will remove the path name and remaining filename only
			if attr.Key == slog.SourceKey {
				source := attr.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return attr
		}

		opts = slog.HandlerOptions{
			AddSource:   true, // write log with filename and function name
			ReplaceAttr: replacer,
		}
	}

	logHandler := slog.NewJSONHandler(writter, &opts)
	logger := slog.New(logHandler).With(slog.String("service", "my-service-name"))
	slog.SetDefault(logger)

	return logger
}

// Attach slog.logger object to context
// It's usefull when slog.logger is attach to http.Request obect
func Attach(logger *slog.Logger, ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyLoggerObj, logger)
}

// FromContext is getting slog.logger from context
// It can be used when slog.logger attached to http.Request object
func FromContext(ctx context.Context) (*slog.Logger, error) {
	l := ctx.Value(KeyLoggerObj)

	logger, ok := l.(*slog.Logger)
	if (logger == nil) || !ok {
		return nil, fmt.Errorf("failed to get logger from context")
	}

	return logger, nil
}
