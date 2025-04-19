package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

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

func NewLogger(logLevel string) *slog.Logger {
	var output io.Writer = os.Stdout
	var level = slog.LevelInfo

	if logLevel == "debug" {
		level = slog.LevelDebug
		output = &prettyJSONWriter{writer: os.Stdout}
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
