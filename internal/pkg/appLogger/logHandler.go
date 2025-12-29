package appLogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

type customHandler struct {
	handler slog.Handler
	writer  io.Writer
	level   slog.Level
	colors  bool
}

func newCustomHandler(w io.Writer, level slog.Level, colors bool) slog.Handler { // TODO доделать хэндлер
	return &customHandler{
		writer: w,
		level:  level,
		colors: colors,
	}
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	levelStr := r.Level.String()
	if h.colors {
		levelStr = h.colorizeLevel(levelStr, r.Level)
	}

	// Собираем атрибуты
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	attrsJSON, _ := json.Marshal(attrs)

	msg := fmt.Sprintf("%s [%s] %s %s\n",
		timeStr,
		levelStr,
		r.Message,
		string(attrsJSON))

	_, err := h.writer.Write([]byte(msg))
	return err
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // упрощенная реализация
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return h // упрощенная реализация
}

func (h *customHandler) colorizeLevel(level string, l slog.Level) string {
	colorMap := map[slog.Level]string{
		slog.LevelDebug: "\033[36m", // голубой
		slog.LevelInfo:  "\033[32m", // зеленый
		slog.LevelError: "\033[31m", // красный
	}

	reset := "\033[0m"
	if color, ok := colorMap[l]; ok {
		return fmt.Sprintf("%s%s%s", color, level, reset)
	}
	return level
}
