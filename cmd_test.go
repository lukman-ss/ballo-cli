package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateDefaultOffline(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := cmdGenerate([]string{"defaultapp"}); err != nil {
		t.Fatalf("generate default offline failed: %v", err)
	}

	if _, err := os.Stat("defaultapp/go.mod"); err != nil {
		t.Fatalf("go.mod missing: %v", err)
	}
}

func TestVerificationSuccess(t *testing.T) {
	dir := t.TempDir()
	modPath := filepath.Join(dir, "go.mod")
	err := os.WriteFile(modPath, []byte("module example.com/app\ngo 1.26\n"), 0644)
	if err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainPath := filepath.Join(dir, "main.go")
	err = os.WriteFile(mainPath, []byte("package main\nfunc main(){}\n"), 0644)
	if err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := testGeneratedProject(dir, "example.com/app"); err != nil {
		t.Fatalf("verification should succeed: %v", err)
	}
}

func TestServeNoGoMod(t *testing.T) {
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	if err := cmdServe([]string{}); err == nil || !strings.Contains(err.Error(), "no go.mod") {
		t.Fatalf("expected no go.mod error, got %v", err)
	}
}
