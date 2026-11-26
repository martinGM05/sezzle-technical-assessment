package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// serve owns the listener and waits for both request draining and Serve to stop.
// The shutdown deadline is independent of the already-canceled signal context.
func serve(ctx context.Context, server *http.Server, listener net.Listener, timeout time.Duration) error {
	served := make(chan error, 1)
	go func() {
		served <- server.Serve(listener)
	}()

	select {
	case err := <-served:
		return serverError(err)
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		shutdownError := server.Shutdown(shutdownContext)
		var closeError error
		if shutdownError != nil {
			// Shutdown does not close active connections when its deadline expires.
			closeError = server.Close()
		}
		return errors.Join(shutdownError, closeError, serverError(<-served))
	}
}

func serverError(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
