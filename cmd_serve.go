package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukman-ss/ballo"
)

// defaultAddr is the default HTTP listen address.
const defaultAddr = ":8080"

// cmdServe builds and starts a Ballo application server on addr.
func cmdServe(args []string) error {
	addr := defaultAddr
	if len(args) > 0 {
		addr = args[0]
	}

	e := ballo.New()
	e.Use(ballo.RecoveryMiddleware())

	e.GET("/healthz", func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusOK, map[string]string{"status": "pass"}), nil
	})

	e.GET("/health", func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusOK, map[string]string{"status": "ok"}), nil
	})

	e.NotFound(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusNotFound, map[string]string{"error": "not found"}), nil
	})

	e.MethodNotAllowed(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}), nil
	})

	srv := &http.Server{Addr: addr, Handler: e}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("serving on %s\n", addr)
	go func() {
		<-sigCh
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "server shutdown error: %v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// newServeMux returns the default serve mux used by serve (exposed for tests).
func newServeMux() *http.ServeMux {
	e := ballo.New()
	e.Use(ballo.RecoveryMiddleware())
	e.GET("/healthz", func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusOK, map[string]string{"status": "pass"}), nil
	})
	e.GET("/health", func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusOK, map[string]string{"status": "ok"}), nil
	})
	e.NotFound(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusNotFound, map[string]string{"error": "not found"}), nil
	})
	e.MethodNotAllowed(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}), nil
	})
	mux := http.NewServeMux()
	mux.Handle("/", e)
	return mux
}
