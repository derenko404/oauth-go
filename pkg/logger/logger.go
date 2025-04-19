package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

const LogFilePath = "logs/application.log"

type prettyJSONWriter struct {
	writer io.Writer
}

func (p *prettyJSONWriter) Write(data []byte) (int, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return p.writer.Write(data)
	}
	return p.writer.Write(append(buf.Bytes(), '\n'))
}

func getLogger(logLevel string, prettyJSON bool) *slog.Logger {
	var output io.Writer = os.Stdout
	var level = slog.LevelInfo

	if logLevel == "debug" {
		level = slog.LevelDebug
		if prettyJSON {
			output = &prettyJSONWriter{writer: os.Stdout}
		} else {
			output = os.Stdout
		}
	} else {
		level = slog.LevelInfo

		file, err := os.OpenFile(LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("failed to open log file", "error", err)
			os.Exit(1)
		}

		output = file
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	return logger
}

func NewLogger(logLevel string) *slog.Logger {
	return getLogger(logLevel, true)
}
