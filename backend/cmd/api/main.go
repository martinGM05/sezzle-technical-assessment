package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	transport "github.com/martinGM05/sezzle-technical-assessment/backend/internal/transport/http"
)

const shutdownTimeout = 5 * time.Second

func main() {
	os.Exit(run(newServer()))
}

func run(server *http.Server) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		slog.Error("server failed", "error", err)
		return 1
	}
	slog.Info("starting server", "address", listener.Addr().String())
	if err := serve(ctx, server, listener, shutdownTimeout); err != nil {
		slog.Error("server failed", "error", err)
		return 1
	}
	return 0
}

func newServer() *http.Server {
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	for index, origin := range origins {
		origins[index] = strings.TrimSpace(origin)
	}

	return &http.Server{
		Addr:              ":8080",
		Handler:           transport.NewHandler(origins...),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
