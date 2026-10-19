// Suite: External Test Tools.
// Invariant: tests resolve tools outside the Base Test Tools through RequireTool.
// Boundary IN: literal os/exec tool lookups in internal, cmd and skills test files.
// Boundary OUT: paths, variable programs and non-test files.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// Sequential: sets PATH
func TestRequireToolReturnsAToolOnPath(t *testing.T) {
	name := "stand-in"
	if runtime.GOOS == "windows" {
		name += ".exe"
		t.Setenv("PATHEXT", ".EXE")
	}
	binary := FixtureBinary(t, name, "package main\nfunc main() {}\n")
	t.Setenv("PATH", filepath.Dir(binary))
	if got := RequireTool(t, name); got != binary {
		t.Fatalf("RequireTool = %q, want %q", got, binary)
	}
}

// Sequential: sets PATH
func TestRequireToolSkipsNamingAnAbsentTool(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	var skipped *testing.T
	var message string
	continued := false
	t.Run("absent", func(t *testing.T) {
		skipped = t
		RequireTool(&toolSkipRecorder{TB: t, message: &message}, "rustfmt")
		continued = true
	})
	if !skipped.Skipped() || continued || message != "requires rustfmt, which is not on PATH" {
		t.Fatalf("skipped = %v, continued = %v, message = %q", skipped.Skipped(), continued, message)
	}
}

type toolSkipRecorder struct {
	testing.TB
	message *string
}

func (t *toolSkipRecorder) Skipf(format string, args ...any) {
	t.Helper()
	*t.message = fmt.Sprintf(format, args...)
	t.TB.Skipf(format, args...)
}

func TestEveryExternalTestToolGoesThroughRequireTool(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	reports, err := externalTestToolReports(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 0 {
		t.Fatalf("External Test Tool rule:\n%s", strings.Join(reports, "\n"))
	}
}

func TestTheExternalTestToolRuleReportsALiteralLookup(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	sources := map[string]string{
		"internal/seeded/fixture_test.go": `package seeded
import "os/exec"
func fixture() {
 exec.LookPath("node")
 exec.Command("rustfmt")
 exec.Command("git")
 exec.Command("/bin/sh")
 program := "node"
 exec.Command(program)
 exec.Command("C:\\tools\\node.exe")
}
`,
		"cmd/seeded/fixture_test.go": `package seeded
import process "os/exec"
func fixture() {
 process.CommandContext(nil, "node")
 process.CommandContext(nil, "go")
}
`,
		"skills/seeded/fixture_test.go": `package seeded
import "os/exec"
func fixture() {
 exec.CommandContext(nil, "rustfmt")
}
`,
		"internal/seeded/production.go": `package seeded
import "os/exec"
func fixture() { exec.LookPath("node") }
`,
	}
	for path, source := range sources {
		path = filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := externalTestToolReports(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"internal/seeded/fixture_test.go:4: node is not a Base Test Tool; resolve it with testfixture.RequireTool",
		"internal/seeded/fixture_test.go:5: rustfmt is not a Base Test Tool; resolve it with testfixture.RequireTool",
		"cmd/seeded/fixture_test.go:4: node is not a Base Test Tool; resolve it with testfixture.RequireTool",
		"skills/seeded/fixture_test.go:4: rustfmt is not a Base Test Tool; resolve it with testfixture.RequireTool",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("reports = %q, want %q", got, want)
	}
}

func externalTestToolReports(root string) ([]string, error) {
	base := map[string]bool{
		"go": true, "gofmt": true, "git": true, "make": true,
		"sh": true, "cat": true, "kill": true, "ps": true,
		"id": true, "mkfifo": true, "cmd": true,
	}
	var reports []string
	for _, directory := range []string{"internal", "cmd", "skills"} {
		err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, walkErr error) error {
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
			names := map[string]bool{}
			for _, imp := range file.Imports {
				if imp.Path.Value == `"os/exec"` {
					name := "exec"
					if imp.Name != nil {
						name = imp.Name.Name
					}
					names[name] = true
				}
			}
			ast.Inspect(file, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				var function string
				switch fun := call.Fun.(type) {
				case *ast.SelectorExpr:
					pkg, ok := fun.X.(*ast.Ident)
					if !ok || !names[pkg.Name] || pkg.Obj != nil {
						return true
					}
					function = fun.Sel.Name
				case *ast.Ident:
					if !names["."] || fun.Obj != nil {
						return true
					}
					function = fun.Name
				default:
					return true
				}
				index := 0
				switch function {
				case "Command", "LookPath":
				case "CommandContext":
					index = 1
				default:
					return true
				}
				if len(call.Args) <= index {
					return true
				}
				literal, ok := call.Args[index].(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING {
					return true
				}
				tool, err := strconv.Unquote(literal.Value)
				if err == nil && !base[tool] && !strings.ContainsAny(tool, "/\\") {
					reports = append(reports, fmt.Sprintf("%s:%d: %s is not a Base Test Tool; resolve it with testfixture.RequireTool", filepath.ToSlash(relative), set.Position(call.Pos()).Line, tool))
				}
				return true
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("scan External Test Tools: %w", err)
		}
	}
	return reports, nil
}
