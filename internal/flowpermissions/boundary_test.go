// Suite: Flow Command List import boundary
// Invariant: only the permission command and the two reports consume the list.
// Boundary IN: imports of all production Go files, including platform files.
// Boundary OUT: tests, generated binaries and dependencies outside this module.
package flowpermissions

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func flowListImportViolations(root string) ([]string, error) {
	allowed := map[string]bool{"internal/cli/permissions.go": true, "internal/cli/doctor_flow_permissions.go": true, "internal/cli/update_check.go": true}
	var violations []string
	for _, directory := range []string{"internal", "cmd"} {
		err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, item := range file.Imports {
				imported, err := strconv.Unquote(item.Path.Value)
				if err != nil {
					return err
				}
				if imported == "roundfix/internal/flowpermissions" && !strings.HasPrefix(relative, "internal/flowpermissions/") && !allowed[relative] {
					violations = append(violations, relative)
					break
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	slices.Sort(violations)
	return violations, nil
}

func TestOnlyTheCommandAndTheReportsUseTheList(t *testing.T) {
	t.Parallel()
	t.Run("this_repository", func(t *testing.T) {
		t.Parallel()
		violations, err := flowListImportViolations(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		if len(violations) > 0 {
			t.Fatalf("unauthorized Flow Command List import in %s", strings.Join(violations, ", "))
		}
	})
	t.Run("planted_baseline_import_is_reported", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		files := map[string]string{
			"internal/baseline/planted.go":            "package baseline\nimport _ \"roundfix/internal/flowpermissions\"\n",
			"internal/cli/permissions.go":             "package cli\nimport _ \"roundfix/internal/flowpermissions\"\n",
			"internal/cli/doctor_flow_permissions.go": "package cli\nimport _ \"roundfix/internal/flowpermissions\"\n",
			"internal/cli/update_check.go":            "package cli\nimport _ \"roundfix/internal/flowpermissions\"\n",
			"internal/baseline/allowed_test.go":       "package baseline\nimport _ \"roundfix/internal/flowpermissions\"\n",
			"cmd/roundfix/main.go":                    "package main\n",
		}
		for name, content := range files {
			path := filepath.Join(root, name)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
		violations, err := flowListImportViolations(root)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(violations, []string{"internal/baseline/planted.go"}) {
			t.Fatalf("violations=%q", violations)
		}
	})
}
