package baseline

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestThePythonModuleDeclaresItsStarterFiles(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	modules, ok := catalog.OrderedModules("python")
	if !ok {
		t.Fatal("python profile missing")
	}
	files, err := catalog.StarterFiles(modules)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{"Makefile", ".roundfixrc.yml", ".github/workflows/verify.yml"}
	ids := []string{"starter.python.makefile", "starter.python.project-config", "starter.python.verify-workflow"}
	lines := [][]string{
		{"SHELL := /bin/bash", ".PHONY: help sync lock-check fmt fmt-check lint typecheck test verify verify-incremental", "help:", "sync:\n\tuv sync --group dev", "lock-check:\n\tuv lock --check --offline", "fmt:\n\tuv run --frozen ruff format src tests", "fmt-check:\n\tuv run --frozen ruff format --check src tests", "lint:\n\tuv run --frozen ruff check src tests", "typecheck:\n\tuv run --frozen ty check --error-on-warning src tests", "test:\n\tuv run --frozen pytest -q", "verify: lock-check fmt-check lint typecheck test", "verify-incremental: verify"},
		{"verification: make verify", "bootstrap: uv sync --group dev --frozen", "bootstrap_timeout: 5m", "caches: [.venv, .pytest_cache, .ruff_cache]", "tools: [uv, make, git]"},
		{"name: Verify", "push:", "branches: [main]", "pull_request:", "contents: read", "name: make verify", "runs-on: ubuntu-latest", "uses: actions/checkout@v7", "persist-credentials: false", "uses: astral-sh/setup-uv@v10.3.0", `python-version: "3.12"`, "run: uv sync --group dev --locked", "run: make verify"},
	}
	if len(files) != len(paths) {
		t.Fatalf("StarterFiles() count = %d, want %d", len(files), len(paths))
	}
	for i, file := range files {
		if file.Path != paths[i] || file.ID != ids[i] || file.Module != "python" {
			t.Errorf("starter %d = %+v", i, file)
		}
		if want := fmt.Sprintf("sha256:%x", sha256.Sum256(file.Content)); file.ContentIdentity != want {
			t.Errorf("identity = %q, want %q", file.ContentIdentity, want)
		}
		if strings.Contains(string(file.Content), "{{") {
			t.Errorf("%s contains template tokens", file.Path)
		}
		for _, line := range lines[i] {
			if !strings.Contains(string(file.Content), line) {
				t.Errorf("%s missing %q", file.Path, line)
			}
		}
	}
	files[0].Content[0] = 'x'
	again, err := catalog.StarterFiles(modules)
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Content[0] != 'S' {
		t.Fatal("StarterFiles exposed mutable catalog bytes")
	}
	for _, profile := range catalog.ProfileIDs() {
		if profile == "python" {
			continue
		}
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			modules, ok := catalog.OrderedModules(profile)
			if !ok {
				t.Fatal("profile missing")
			}
			files, err := catalog.StarterFiles(modules)
			if err != nil || len(files) != 0 {
				t.Fatalf("StarterFiles(%s) = %v, %v", profile, files, err)
			}
		})
	}
	t.Run("module order", func(t *testing.T) {
		t.Parallel()
		// Repeating modules makes their order observable without inventing another starter-owning module.
		files, err := catalog.StarterFiles([]string{"python", "core", "python"})
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, file := range files {
			got = append(got, file.Path)
		}
		if want := append(slices.Clone(paths), paths...); !slices.Equal(got, want) {
			t.Fatalf("paths = %v, want %v", got, want)
		}
	})
	t.Run("unknown module", func(t *testing.T) {
		t.Parallel()
		if _, err := catalog.StarterFiles([]string{"python", "unknown"}); err == nil {
			t.Fatal("unknown module accepted")
		}
	})
}

func TestAStarterFileDeclarationIsValidated(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, field string
		value       any
		code        string
	}{
		{"unsafe path", "path", "../Makefile", "catalog.starter.file.path.invalid"},
		{"absolute path", "path", "/Makefile", "catalog.starter.file.path.invalid"},
		{"Windows path", "path", `C:\Makefile`, "catalog.starter.file.path.invalid"},
		{"empty path", "path", "", "catalog.starter.file.path.invalid"},
		{"duplicate path", "path", ".roundfixrc.yml", "catalog.starter.file.path.duplicate"},
		{"unknown template", "template", "template.unknown", "catalog.starter.file.template.unknown"},
		{"guide template", "template", "template.guide.python", "catalog.starter.file.template.kind"},
		{"missing id", "id", "", "catalog.starter.file.id.missing"},
		{"duplicate id", "id", "starter.python.project-config", "catalog.starter.file.id.duplicate"},
		{"not an array", "starterFiles", "invalid", "catalog.module.starterFiles.invalid"},
		{"null array", "starterFiles", nil, "catalog.module.starterFiles.invalid"},
		{"not an object", "starterFiles", []any{"invalid"}, "catalog.module.starterFiles.invalid"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assets := cloneEmbeddedAssets(t)
			var module map[string]any
			if err := json.Unmarshal(assets["modules/python.json"].Data, &module); err != nil {
				t.Fatal(err)
			}
			if test.field == "starterFiles" {
				module[test.field] = test.value
			} else {
				module["starterFiles"].([]any)[0].(map[string]any)[test.field] = test.value
			}
			data, err := json.Marshal(module)
			if err != nil {
				t.Fatal(err)
			}
			assets["modules/python.json"].Data = data
			_, err = LoadCatalog(assets)
			var invalid *ValidationError
			if !errors.As(err, &invalid) || !invalid.Has(test.code) {
				t.Fatalf("LoadCatalog() = %v, want %s", err, test.code)
			}
		})
	}
	t.Run("absent field", func(t *testing.T) {
		t.Parallel()
		assets := cloneEmbeddedAssets(t)
		var module map[string]any
		if err := json.Unmarshal(assets["modules/python.json"].Data, &module); err != nil {
			t.Fatal(err)
		}
		delete(module, "starterFiles")
		data, err := json.Marshal(module)
		if err != nil {
			t.Fatal(err)
		}
		assets["modules/python.json"].Data = data
		catalog, err := LoadCatalog(assets)
		if err != nil {
			t.Fatal(err)
		}
		files, err := catalog.StarterFiles([]string{"python"})
		if err != nil || len(files) != 0 {
			t.Fatalf("StarterFiles() = %v, %v", files, err)
		}
	})
}
