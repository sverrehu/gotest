package main

import (
	"log/slog"
	"os"
)

type Something struct {
	name string
	age  int
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	data := Something{
		name: "John Doe",
		age:  42,
	}
	slog.Debug("number one", "data", data)
	slog.Info("number two", "data", data)
}
