package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestServeRunsAndCleansTemp(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := os.WriteFile("go.mod", []byte("module example.com/app\ngo 1.26\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainSrc := `package main
import (
    // "log"
    "net/http"
    "os"
    "time"
)
func main() {
    addr := os.Getenv("HTTP_ADDR")
    if addr == "" { addr = ":8080" }
    srv := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })}
    go func() { <-time.After(300 * time.Millisecond); srv.Close() }()
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { panic(err) }
}
`
	if err := os.WriteFile("main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := cmdServe([]string{":9090"}); err != nil {
		t.Fatalf("cmdServe error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(os.TempDir(), "ballo_cli_user_app_*"))
	if len(matches) != 0 {
		t.Fatalf("temp dirs not cleaned: %v", matches)
	}
}

func TestServePassesAddrToApp(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := os.WriteFile("go.mod", []byte("module example.com/app\ngo 1.26\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainSrc := `package main
import (
    "net/http"
    "os"
    "time"
)
func main() {
    addr := os.Getenv("HTTP_ADDR")
    if addr == "" { addr = ":8080" }
    srv := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })}
    go func() { <-time.After(300 * time.Millisecond); srv.Close() }()
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
    _ = os.Args
}
`
	if err := os.WriteFile("main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := cmdServe([]string{":9091"}); err != nil {
		t.Fatalf("cmdServe error: %v", err)
	}
	_ = strings.TrimSpace
}

func TestServeProcessFailure(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := os.WriteFile("go.mod", []byte("module example.com/app\ngo 1.26\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainSrc := `package main
import "os"
func main() { os.Exit(1) }
`
	if err := os.WriteFile("main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	err := cmdServe([]string{})
	if err == nil {
		t.Fatal("expected error for failing process")
	}
	if !strings.Contains(err.Error(), "run:") {
		t.Fatalf("expected run error, got: %v", err)
	}
}

func TestServeCleanupOnBuildFailure(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := os.WriteFile("go.mod", []byte("module example.com/app\ngo 1.26\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	// broken main.go
	if err := os.WriteFile("main.go", []byte("package main\nfunc main() { this is broken\n"), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	err := cmdServe([]string{})
	if err == nil || !strings.Contains(err.Error(), "build failed") {
		t.Fatalf("expected build failed error, got: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(os.TempDir(), "ballo_cli_user_app_*"))
	if len(matches) != 0 {
		t.Fatalf("temp dirs not cleaned after build failure: %v", matches)
	}
}

func TestServeShutdownGraceful(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := os.WriteFile("go.mod", []byte("module example.com/app\ngo 1.26\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainSrc := `package main
import (
    "net/http"
    "os"
    "time"
)
func main() {
    addr := os.Getenv("HTTP_ADDR")
    if addr == "" { addr = ":8080" }
    srv := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })}
    go func() { <-time.After(500 * time.Millisecond); srv.Close() }()
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}
`
	if err := os.WriteFile("main.go", []byte(mainSrc), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmdServe([]string{":9092"}) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("cmdServe error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve timed out")
	}
}
