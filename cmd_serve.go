package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// defaultAddr is the default HTTP listen address.
const defaultAddr = ":8080"

// cmdServe runs the Go application in the current directory.
func cmdServe(args []string) error {
	addr := defaultAddr
	if len(args) > 0 {
		addr = args[0]
	}

	// Ensure current directory has a go.mod
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("current directory: %w", err)
	}

	goModPath := filepath.Join(cwd, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return errors.New("no go.mod found in current directory; use 'generate' to create a project")
	}

	// Build current project
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(os.TempDir(), "ballo_cli_user_app"), ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Dir = cwd
	buildCmd.Env = append(os.Environ(), "GO111MODULE=on")

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Run built binary
	binPath := filepath.Join(os.TempDir(), "ballo_cli_user_app")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return errors.New("binary not built")
	}

	srvCmd := exec.Command(binPath)
	srvCmd.Stdout = os.Stdout
	srvCmd.Stderr = os.Stderr
	srvCmd.Dir = cwd

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("serving project from %s on %s\n", cwd, addr)
	if err := srvCmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	go func() {
		<-sigCh
		srvCmd.Process.Signal(syscall.SIGTERM)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		<-ctx.Done()
		srvCmd.Process.Kill()
	}()

	if err := srvCmd.Wait(); err != nil {
		if strings.Contains(err.Error(), "signal: ") {
			return nil
		}
		return fmt.Errorf("run: %w", err)
	}
	return nil
}
