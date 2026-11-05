package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// defaultAddr is the default HTTP listen address.
const defaultAddr = ":8080"

// cmdServe builds and starts an HTTP server on addr.
func cmdServe(args []string) error {
	addr := defaultAddr
	if len(args) > 0 {
		addr = args[0]
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"pass"}`)
	})

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		srv.Close()
	}()

	fmt.Printf("serving on %s\n", addr)
	return srv.ListenAndServe()
}

// newServeMux returns the default serve mux used by serve (exposed for tests).
func newServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"pass"}`)
	})
	return mux
}
