package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIStandardLibraryImports(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "mikoto")
	build := exec.Command("go", "build", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	files := map[string]string{
		"go.mod":    "module example.com/stdlib\n\ngo 1.26.0\n",
		"stdlib.go": "package stdlib\nimport \"os\"\nfunc ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }\n",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(binary, "-max-lines=80", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("analyze standard library import: %v\n%s", err, output)
	}
}
