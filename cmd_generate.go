package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const goModTmpl = `module %s

go 1.26.0

require github.com/lukman-ss/ballo v0.1.0
`

const mainTmpl = `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukman-ss/ballo"
)

func main() {
	e := ballo.New()
	e.Use(ballo.RecoveryMiddleware())

	e.GET("/health", func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusOK, map[string]string{"status": "ok"}), nil
	})

	e.NotFound(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusNotFound, map[string]string{"error": "not found"}), nil
	})

	e.MethodNotAllowed(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}), nil
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Println("server starting on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
	log.Println("server stopped")
}
`

const envExample = `APP_ENV=development
HTTP_ADDR=:8080
LOG_LEVEL=info
GO_PROJECT=%s
`

var modulePattern = regexp.MustCompile(`^([a-zA-Z0-9_.-]+/)*[a-zA-Z0-9_.-]+$`)

func isValidModulePath(path string) bool {
	if path == "" || len(path) > 256 {
		return false
	}
	if strings.HasPrefix(path, ".") || strings.Contains(path, "..") {
		return false
	}
	return modulePattern.MatchString(path)
}

// cmdGenerate scaffolds a new Ballo project.
func cmdGenerate(args []string) error {
	verify := false
	var name string
	for _, arg := range args {
		if arg == "--verify" {
			verify = true
			continue
		}
		if !strings.HasPrefix(arg, "-") && name == "" {
			name = arg
			continue
		}
	}

	if name == "" {
		return errors.New("missing project name")
	}

	if !isValidModulePath(name) {
		return errors.New("invalid project name: must be valid Go module path, no path traversal")
	}

	dest, err := filepath.Abs(name)
	if err != nil {
		return fmt.Errorf("absolute path: %w", err)
	}

	// Resolve to avoid path traversal via symlinks
	if resolved, err := filepath.EvalSymlinks(filepath.Dir(dest)); err == nil {
		dest = filepath.Join(resolved, filepath.Base(dest))
	}

	if err := ensureDirEmpty(dest); err != nil {
		return err
	}

	files := map[string]string{
		"go.mod":       fmt.Sprintf(goModTmpl, name),
		"main.go":      mainTmpl,
		".env.example": fmt.Sprintf(envExample, name),
	}

	for fname, content := range files {
		fpath := filepath.Join(dest, fname)
		if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
			os.RemoveAll(dest)
			return fmt.Errorf("write %s: %w", fname, err)
		}
	}

	if verify {
		if err := testGeneratedProject(dest, name); err != nil {
			return fmt.Errorf("verification failed for %s: %w", dest, err)
		}
	}

	fmt.Printf("generated: %s\n", dest)
	return nil
}

func ensureDirEmpty(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("directory already exists: %s", path)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	return nil
}

func testGeneratedProject(dir, module string) error {
	oldWd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(dir); err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	if err := cmd.Run(); err != nil {
		// ponytail: keep project; caller decides removal.
		return fmt.Errorf("go mod tidy: %w", err)
	}

	cmd = exec.Command("go", "build", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}
