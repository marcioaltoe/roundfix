// Suite: written executable residue.
// Invariant: literal executable writes in internal tests match the frozen inventory exactly.
// Boundary IN: os.WriteFile and os.OpenFile calls with integer literal modes in internal/**/*_test.go.
// Boundary OUT: modes held in constants, computed modes and os.Chmod calls.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type writtenExecutable struct {
	file string
	line int
}

var writtenExecutableResidue = map[string]int{
	"internal/cli/cli_test.go":                        1,
	"internal/cli/upgrade_test.go":                    1,
	"internal/baseline/plan_characterization_test.go": 1,
	"internal/baseline/profile_alignment_test.go":     1,
	"internal/baseline/skills_lock_read_test.go":      1,
	"internal/daemon/commit_hook_test.go":             1,
	"internal/daemon/task_engine_test.go":             2,
	"internal/preflight/preflight_test.go":            1,
	"internal/store/process_unix_test.go":             1,
}

func writtenExecutables(root string) ([]writtenExecutable, error) {
	var sites []writtenExecutable
	err := filepath.WalkDir(filepath.Join(root, "internal"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		set := token.NewFileSet()
		file, err := parser.ParseFile(set, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("resolve %s: %w", path, err)
		}
		// Resolve imported names, including aliases, so a local variable named os is not mistaken for the package.
		names := map[string]bool{}
		for _, imp := range file.Imports {
			if imp.Path.Value == `"os"` {
				name := "os"
				if imp.Name != nil {
					name = imp.Name.Name
				}
				names[name] = true
			}
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) != 3 {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := selector.X.(*ast.Ident)
			if !ok || !names[pkg.Name] || pkg.Obj != nil {
				return true
			}
			if selector.Sel.Name != "WriteFile" && selector.Sel.Name != "OpenFile" {
				return true
			}
			literal, ok := call.Args[2].(*ast.BasicLit)
			if !ok || literal.Kind != token.INT {
				return true
			}
			mode, err := strconv.ParseUint(literal.Value, 0, 32)
			if err == nil && mode&0o111 != 0 {
				sites = append(sites, writtenExecutable{file: filepath.ToSlash(relative), line: set.Position(call.Pos()).Line})
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan written executables: %w", err)
	}
	return sites, nil
}

func checkWrittenExecutableResidue(root string, residue map[string]int) error {
	sites, err := writtenExecutables(root)
	if err != nil {
		return err
	}
	counts := map[string]int{}
	var findings []string
	for _, site := range sites {
		counts[site.file]++
		if _, ok := residue[site.file]; !ok {
			findings = append(findings, fmt.Sprintf("%s:%d: written executable outside the residue", site.file, site.line))
		}
	}
	for file, want := range residue {
		if got := counts[file]; got != want {
			var locations []string
			for _, site := range sites {
				if site.file == file {
					locations = append(locations, fmt.Sprintf("%s:%d", file, site.line))
				}
			}
			findings = append(findings, fmt.Sprintf("%s: executable write count = %d, want %d; sites: %s", file, got, want, strings.Join(locations, ", ")))
		}
	}
	if len(findings) != 0 {
		sort.Strings(findings)
		return fmt.Errorf("written executable inventory differs:\n%s", strings.Join(findings, "\n"))
	}
	return nil
}

func TestNoTestWritesAnExecutableOutsideTheResidue(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		if err := checkWrittenExecutableResidue(root, writtenExecutableResidue); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("seeded executable names its site", func(t *testing.T) {
		root := writtenExecutableTree(t, "package seeded\nimport \"os\"\nfunc fixture() { os.WriteFile(\"fixture\", nil, 0o755) }\n")
		err := checkWrittenExecutableResidue(root, nil)
		if err == nil || !strings.Contains(err.Error(), "internal/seeded/fixture_test.go:3") {
			t.Fatalf("guard error = %v, want seeded file:line", err)
		}
	})
	t.Run("inventory count differs by one", func(t *testing.T) {
		root := writtenExecutableTree(t, "package seeded\nimport \"os\"\nfunc fixture() { os.OpenFile(\"fixture\", os.O_CREATE, 0o100) }\n")
		for _, want := range []int{0, 1, 2} {
			err := checkWrittenExecutableResidue(root, map[string]int{"internal/seeded/fixture_test.go": want})
			if want == 1 {
				if err != nil {
					t.Fatalf("exact count rejected: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("count = 1, want %d", want)) {
				t.Fatalf("guard error = %v, want count mismatch for %d", err, want)
			}
		}
	})
	t.Run("literal modes and boundary", func(t *testing.T) {
		root := writtenExecutableTree(t, `package seeded
import filesystem "os"
const executable = 0o755
func fixture() {
 filesystem.WriteFile("a", nil, 0o001)
 filesystem.OpenFile("b", filesystem.O_CREATE, 0110)
 filesystem.WriteFile("c", nil, 73)
 filesystem.WriteFile("data", nil, 0o600)
 filesystem.WriteFile("constant", nil, executable)
 filesystem.Chmod("data", 0o755)
}
`)
		sites, err := writtenExecutables(root)
		if err != nil {
			t.Fatal(err)
		}
		if len(sites) != 3 {
			t.Fatalf("sites = %+v, want three literal executable modes", sites)
		}
		for index, site := range sites {
			if site.line != index+5 {
				t.Fatalf("site = %+v, want line %d", site, index+5)
			}
		}
	})
}

func writtenExecutableTree(t *testing.T, source string) string {
	t.Helper()
	root := t.TempDir()
	directory := filepath.Join(root, "internal", "seeded")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture_test.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}
