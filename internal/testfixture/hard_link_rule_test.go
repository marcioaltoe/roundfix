// Suite: Hard Link repository rule.
// Invariant: no internal test creates a hard link, preserving Single-Name Test Binaries.
// Boundary IN: imported os.Link, syscall.Link/Linkat and unix.Link/Linkat calls in internal/**/*_test.go.
// Boundary OUT: other receivers, local variables, symbolic links and link reads.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func hardLinkViolations(root string) ([]string, error) {
	var violations []string
	err := filepath.WalkDir(filepath.Join(root, "internal"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		positions := token.NewFileSet()
		file, err := parser.ParseFile(positions, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse hard link rule source %s: %w", path, err)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("resolve %s: %w", path, err)
		}
		imports := make(map[string]string)
		for _, imp := range file.Imports {
			imported, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return fmt.Errorf("read import in %s: %w", path, err)
			}
			switch imported {
			case "os", "syscall", "golang.org/x/sys/unix":
				name := filepath.Base(imported)
				if imp.Name != nil {
					name = imp.Name.Name
				}
				imports[name] = imported
			}
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			receiver, ok := selector.X.(*ast.Ident)
			if !ok || receiver.Obj != nil {
				return true
			}
			imported, ok := imports[receiver.Name]
			if !ok || (selector.Sel.Name != "Link" && (selector.Sel.Name != "Linkat" || imported == "os")) {
				return true
			}
			violations = append(violations, fmt.Sprintf("%s:%d: %s.%s creates a hard link; copy or use a symbolic link", filepath.ToSlash(relative), positions.Position(call.Pos()).Line, filepath.Base(imported), selector.Sel.Name))
			return true
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan hard links: %w", err)
	}
	return violations, nil
}

func TestNoTestCreatesAHardLink(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		violations, err := hardLinkViolations(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, violation := range violations {
			t.Error(violation)
		}
	})
	for _, tc := range []struct {
		name, source string
		want         []string
	}{
		{"os Link names its site", "package seeded\nimport \"os\"\nfunc fixture() { os.Link(\"source\", \"target\") }\n", []string{"internal/seeded/fixture_test.go:3: os.Link"}},
		{"aliased os import", "package seeded\nimport filesystem \"os\"\nfunc fixture() { filesystem.Link(\"source\", \"target\") }\n", []string{"internal/seeded/fixture_test.go:3: os.Link"}},
		{"syscall Link", "package seeded\nimport \"syscall\"\nfunc fixture() { syscall.Link(\"source\", \"target\") }\n", []string{"internal/seeded/fixture_test.go:3: syscall.Link"}},
		{"aliased syscall Linkat", "package seeded\nimport system \"syscall\"\nfunc fixture() { system.Linkat(0, \"source\", 0, \"target\", 0) }\n", []string{"internal/seeded/fixture_test.go:3: syscall.Linkat"}},
		{"unix Link", "package seeded\nimport \"golang.org/x/sys/unix\"\nfunc fixture() { unix.Link(\"source\", \"target\") }\n", []string{"internal/seeded/fixture_test.go:3: unix.Link"}},
		{"unix Linkat", "package seeded\nimport \"golang.org/x/sys/unix\"\nfunc fixture() { unix.Linkat(0, \"source\", 0, \"target\", 0) }\n", []string{"internal/seeded/fixture_test.go:3: unix.Linkat"}},
		{"aliased unix import", "package seeded\nimport system \"golang.org/x/sys/unix\"\nfunc fixture() { system.Link(\"source\", \"target\"); system.Linkat(0, \"source\", 0, \"target\", 0) }\n", []string{"internal/seeded/fixture_test.go:3: unix.Link", "internal/seeded/fixture_test.go:3: unix.Linkat"}},
		{"other receiver and local Link", `package seeded
import "net/url"
type filesystem struct{}
func (filesystem) Link(string, string) {}
func fixture() {
 local := filesystem{}
 local.Link("source", "target")
 Link := func(string, string) {}
 Link("source", "target")
 url.Link("source", "target")
}
`, nil},
		{"shadowed os import", `package seeded
import "os"
type filesystem struct{}
func (filesystem) Link(string, string) {}
func fixture() {
 os := filesystem{}
 os.Link("source", "target")
}
`, nil},
		{"symbolic links and reads", "package seeded\nimport \"os\"\nfunc fixture() { os.Symlink(\"source\", \"target\"); os.Readlink(\"target\") }\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			directory := filepath.Join(root, "internal", "seeded")
			if err := os.MkdirAll(directory, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(directory, "fixture_test.go"), []byte(tc.source), 0o600); err != nil {
				t.Fatal(err)
			}
			violations, err := hardLinkViolations(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(violations) != len(tc.want) {
				t.Fatalf("violations = %q, want %q", violations, tc.want)
			}
			for index, want := range tc.want {
				if !strings.HasPrefix(violations[index], want) {
					t.Errorf("violation = %q, want prefix %q", violations[index], want)
				}
			}
		})
	}
}
