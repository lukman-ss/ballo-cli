package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const goModTmpl = `module %s

go 1.26.0

require github.com/lukman-ss/ballo v0.1.0
`

const mainTmpl = `package main

import (
	"log"
	"net/http"
)

func main() {
	engine := NewApp()
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, engine); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
`

const envExample = `APP_ENV=development
HTTP_ADDR=:8080
LOG_LEVEL=info
GO_PROJECT=github.com/user/myapp
`

// cmdGenerate scaffolds a new Ballo project.
func cmdGenerate(args []string) error {
	if len(args) == 0 {
		return errors.New("missing template name")
	}

	name := args[0]
	dir := filepath.Join(".", name)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("directory already exists: %s", dir)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	files := map[string]string{
		"go.mod":       fmt.Sprintf(goModTmpl, name),
		"main.go":      mainTmpl,
		".env.example": envExample,
	}
	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(dir, fname), []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s: %w", fname, err)
		}
	}

	fmt.Printf("generated: %s\n", dir)
	return nil
}
