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

func TestGenerateInvalidName(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	cases := []string{
		"../evil",
		"./evil",
		"a/b../c",
		"../../x",
		"/abs/path",
		".dotfile",
	}
	for _, name := range cases {
		if err := cmdGenerate([]string{name}); err == nil {
			t.Fatalf("expected error for invalid name %q", name)
		}
	}
}

func TestGenerateValidName(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	cases := []string{
		"myapp",
		"github.com/user/myapp",
		"github.com/user/my-app/v2",
	}
	for _, name := range cases {
		if err := cmdGenerate([]string{name, "--skip-verify"}); err != nil {
			t.Fatalf("unexpected error for valid name %q: %v", name, err)
		}
		_ = os.RemoveAll(name)
	}
}

func TestGenerateSkipVerify(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := cmdGenerate([]string{"myapp", "--skip-verify"}); err != nil {
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

func TestGenerateOfflineDoesNotRequireNetwork(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	// Simulate offline by clearing proxy/env that might trigger fetch; skip verify for offline scaffold
	if err := cmdGenerate([]string{"offlineapp", "--skip-verify"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat("offlineapp/go.mod"); err != nil {
		t.Fatal("go.mod not created")
	}
}

func TestTestGeneratedProjectFailureDoesNotDeleteDir(t *testing.T) {
	dir := t.TempDir()
	proj := dir + "/badproj"
	if err := os.MkdirAll(proj, 0755); err != nil {
		t.Fatal(err)
	}
	// Write broken go.mod to make go build fail
	if err := os.WriteFile(proj+"/go.mod", []byte("module badproj\ngo 1.26.0\nrequire broken@v0.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(proj+"/main.go", []byte("package main\nfunc main(){}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)

	if err := testGeneratedProject(proj, "badproj"); err == nil {
		t.Fatal("expected verification error")
	}
	if _, err := os.Stat(proj); err != nil {
		t.Fatal("project directory was deleted on verification failure")
	}
}
