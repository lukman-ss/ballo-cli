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

	// Build current project in a unique temp directory to avoid collisions
	tmpDir, err := os.MkdirTemp("", "ballo_cli_user_app_*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	binPath := filepath.Join(tmpDir, "ballo_cli_user_app")
	// Ensure cleanup on return
	defer os.RemoveAll(tmpDir)

	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Dir = cwd
	buildCmd.Env = append(os.Environ(), "GO111MODULE=on")

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return errors.New("binary not built")
	}

	srvCmd := exec.Command(binPath)
	srvCmd.Stdout = os.Stdout
	srvCmd.Stderr = os.Stderr
	srvCmd.Dir = cwd
	// Pass address to app via documented contract (HTTP_ADDR)
	srvCmd.Env = append(os.Environ(), "HTTP_ADDR="+addr)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	defer close(sigCh)

	fmt.Printf("serving project from %s on %s\n", cwd, addr)
	if err := srvCmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- srvCmd.Wait()
	}()

	// Wait for signal or process exit
	select {
	case <-sigCh:
		// Gracefully terminate child
		if srvCmd.Process != nil {
			_ = srvCmd.Process.Signal(syscall.SIGTERM)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		select {
		case <-done:
			cancel()
		case <-ctx.Done():
			cancel()
			if srvCmd.Process != nil {
				_ = srvCmd.Process.Kill()
			}
			// Drain done
			<-done
		}
		return nil
	case err := <-done:
		if err != nil {
			if strings.Contains(err.Error(), "signal: ") {
				return nil
			}
			return fmt.Errorf("run: %w", err)
		}
		return nil
	}
}
