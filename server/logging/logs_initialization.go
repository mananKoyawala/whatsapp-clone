package logger

import (
	"log/slog"
	"os"
	"time"
)

var loggingPath = "../logging/logs/"

func InitUserLogger() *slog.Logger {

	// preparing users.log file // rw----rw-
	file, err := os.OpenFile(loggingPath+"users.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0606)
	if err != nil {
		loggerInitError("users", err)
	}

	// get user logger
	loggerHandler := getLogger(file)

	return slog.New(loggerHandler)
}

func InitMessageLogger() *slog.Logger {

	// preparing messages.log file
	file, err := os.OpenFile(loggingPath+"messages.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0606)
	if err != nil {
		loggerInitError("messages", err)
	}

	// get user logger
	loggerHandler := getLogger(file)

	return slog.New(loggerHandler)
}

func InitGroupLogger() *slog.Logger {

	// preparing groups.log file
	file, err := os.OpenFile(loggingPath+"groups.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0606)
	if err != nil {
		loggerInitError("groups", err)
	}

	// get user logger
	loggerHandler := getLogger(file)

	return slog.New(loggerHandler)
}

func InitContactLogger() *slog.Logger {

	// preparing contacts.log file
	file, err := os.OpenFile(loggingPath+"contacts.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0606)
	if err != nil {
		loggerInitError("contacts", err)
	}

	// get user logger
	loggerHandler := getLogger(file)

	return slog.New(loggerHandler)
}

func loggerInitError(filename string, err error) {
	slog.Error("Failed to open", filename, "  log file", slog.String("error", err.Error()))
	os.Exit(1)
}

func getLogger(file *os.File) *slog.JSONHandler {
	return slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			// it changes default key time to date
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
				a.Value = slog.Int64Value(time.Now().Unix())
			}
			return a
		},
	})
}
