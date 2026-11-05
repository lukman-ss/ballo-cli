package main

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRunNoArgs(t *testing.T) {
	if got := run([]string{}); got != exitUsage {
		t.Fatalf("expected %d, got %d", exitUsage, got)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	if got := run([]string{"foobar"}); got != exitUsage {
		t.Fatalf("expected %d, got %d", exitUsage, got)
	}
}

func TestRunHelp(t *testing.T) {
	if got := run([]string{"help"}); got != exitOK {
		t.Fatalf("expected %d, got %d", exitOK, got)
	}
}

func TestServeHealthz(t *testing.T) {
	mux := newServeMux()
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"pass"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestServe404(t *testing.T) {
	mux := newServeMux()
	req := httptest.NewRequest("GET", "/nope", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGenerateMissingName(t *testing.T) {
	if err := cmdGenerate([]string{}); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestGenerateSucceeds(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := cmdGenerate([]string{"myapp"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat("myapp/go.mod"); err != nil {
		t.Fatal("go.mod not created")
	}
	if _, err := os.Stat("myapp/main.go"); err != nil {
		t.Fatal("main.go not created")
	}
	if _, err := os.Stat("myapp/.env.example"); err != nil {
		t.Fatal(".env.example not created")
	}
}
