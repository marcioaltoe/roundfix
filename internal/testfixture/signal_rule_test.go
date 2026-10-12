// Suite: Proven Signal repository rule.
// Invariant: recorded PID signals use an identity proof, including in tests.
// Boundary IN: every Go source file and isolated seeded source trees.
// Boundary OUT: skipped artifact directories and child handles returned by Start.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Function-scoped exceptions: an adjacent function in the same file is checked.
func provenSignalAllowedFunction(relative, name string) bool {
	switch filepath.ToSlash(filepath.Dir(relative)) {
	case "internal/store":
		return name == "SignalProvenProcess" || name == "SignalProvenProcessGroup" || name == "signalOwnerProcess"
	case "internal/cli":
		return name == "interruptSelf"
	}
	return false
}

func signalRuleImportedCall(expression ast.Expr, imports map[string]string, pkg, name string) bool {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != name {
		return false
	}
	receiver, ok := selector.X.(*ast.Ident)
	return ok && imports[receiver.Name] == pkg
}

func signalRuleZero(expression ast.Expr, imports map[string]string) bool {
	if literal, ok := expression.(*ast.BasicLit); ok {
		return literal.Kind == token.INT && literal.Value == "0"
	}
	call, ok := expression.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 || !signalRuleImportedCall(call.Fun, imports, "syscall", "Signal") {
		return false
	}
	literal, ok := call.Args[0].(*ast.BasicLit)
	return ok && literal.Kind == token.INT && literal.Value == "0"
}

var shellSignalCommand = regexp.MustCompile(`(?m)(^|[;&|\n])[\t ]*(?:(?:if|while|then|exec)[\t ]+(?:![\t ]*)?)?(?:/[^\s;|&]+/)?(kill|pkill|killall)\b(?:[\t ]+([^;|&\n]*))?`)

func provenSignalViolations(root string) ([]string, error) {
	var violations []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".gocache", "bin", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		positions := token.NewFileSet()
		file, err := parser.ParseFile(positions, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse signal rule source %s: %w", relative, err)
		}
		imports := make(map[string]string)
		for _, imp := range file.Imports {
			imported, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return err
			}
			name := filepath.Base(imported)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			imports[name] = imported
		}
		report := func(position token.Pos, offset int, form string) {
			violations = append(violations, fmt.Sprintf("%s:%d: %s", relative, positions.Position(position).Line+offset, form))
		}
		for _, declaration := range file.Decls {
			if function, ok := declaration.(*ast.FuncDecl); ok && function.Recv == nil && provenSignalAllowedFunction(relative, function.Name.Name) {
				continue
			}
			// Diagnostic strings describe signals; they do not run a shell command.
			diagnostic := make(map[*ast.BasicLit]bool)
			commandLiterals := make(map[*ast.BasicLit]bool)
			ast.Inspect(declaration, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
					if selector.Sel.Name == "Errorf" || selector.Sel.Name == "Fatalf" || selector.Sel.Name == "Logf" || signalRuleImportedCall(call.Fun, imports, "errors", "New") {
						for _, argument := range call.Args {
							if literal, ok := argument.(*ast.BasicLit); ok {
								diagnostic[literal] = true
							}
						}
					}
				}
				// Also cover a command name and arguments supplied as separate literals,
				// including through repository exec runners rather than exec.Command.
				for index, argument := range call.Args {
					literal, ok := argument.(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						continue
					}
					value, _ := strconv.Unquote(literal.Value)
					command := filepath.Base(value)
					if command != "kill" && command != "pkill" && command != "killall" {
						continue
					}
					commandLiterals[literal] = true
					zero := false
					if index+1 < len(call.Args) {
						if next, ok := call.Args[index+1].(*ast.BasicLit); ok {
							flag, _ := strconv.Unquote(next.Value)
							zero = command == "kill" && flag == "-0"
						}
					}
					if !zero {
						report(literal.Pos(), 0, "shell "+command)
					}
				}
				return true
			})
			ast.Inspect(declaration, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.CallExpr:
					for _, pkg := range []string{"syscall", "golang.org/x/sys/unix"} {
						if signalRuleImportedCall(node.Fun, imports, pkg, "Kill") && (len(node.Args) != 2 || !signalRuleZero(node.Args[1], imports)) {
							report(node.Pos(), 0, filepath.Base(pkg)+".Kill without proof")
						}
					}
					if signalRuleImportedCall(node.Fun, imports, "os", "FindProcess") {
						report(node.Pos(), 0, "os.FindProcess without proof")
					}
				case *ast.BasicLit:
					if node.Kind != token.STRING || diagnostic[node] || commandLiterals[node] {
						break
					}
					value, err := strconv.Unquote(node.Value)
					if err != nil {
						break
					}
					for _, match := range shellSignalCommand.FindAllStringSubmatchIndex(value, -1) {
						command := value[match[4]:match[5]]
						args := ""
						if match[6] >= 0 {
							args = value[match[6]:match[7]]
						}
						fields := strings.Fields(args)
						if len(fields) == 0 {
							continue
						} // A command-name comparison does not run a command.
						if command == "kill" && len(fields) > 0 && fields[0] == "-0" {
							continue
						}
						offset := 0
						if strings.HasPrefix(node.Value, "`") {
							offset = strings.Count(value[:match[4]], "\n")
						}
						report(node.Pos(), offset, "shell "+command)
					}
				}
				return true
			})
		}
		return nil
	})
	sort.Strings(violations)
	return violations, err
}

func TestEverySignalGoesToAProvenProcess(t *testing.T) {
	t.Parallel()
	t.Run("repository", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		violations, err := provenSignalViolations(root)
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
		{"direct-kill", "package store\nimport s \"syscall\"\nfunc forbidden() {\n s.Kill(123, s.SIGKILL)\n}\nfunc SignalProvenProcess() { s.Kill(123, s.SIGKILL) }\n", []string{"internal/store/seed.go:4: syscall.Kill"}},
		{"bare-pid-lookup", "package store\nimport \"os\"\nfunc forbidden() {\n os.FindProcess(123)\n}\n", []string{"internal/store/seed.go:4: os.FindProcess"}},
		{"group-kill", "package store\nimport u \"golang.org/x/sys/unix\"\nfunc forbidden() {\n u.Kill(-123, u.SIGTERM)\n}\n", []string{"internal/store/seed.go:4: unix.Kill"}},
		{"shell-kill", "package store\nfunc forbidden() {\n _ = `#!/bin/sh\n" + "ki" + "ll -9 123\n" + "pk" + "ill child\n" + "kill" + "all child\n`\n}\n", []string{"internal/store/seed.go:4: shell kill", "internal/store/seed.go:5: shell pkill", "internal/store/seed.go:6: shell killall"}},
		{"zero-probe", "package store\nimport (\"syscall\"; \"os/exec\")\nfunc probe() { syscall.Kill(123, 0); syscall.Kill(-123, syscall.Signal(0)); exec.Command(\"kill\", \"-0\", \"123\"); _ = `" + "ki" + "ll -0 123` }\n", nil},
		{"allowlisted", "package store\nimport (\"syscall\"; \"os\")\nfunc SignalProvenProcess() { syscall.Kill(123, syscall.SIGKILL) }\nfunc SignalProvenProcessGroup() { syscall.Kill(-123, syscall.SIGKILL) }\nfunc signalOwnerProcess() { os.FindProcess(123) }\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			dir := filepath.Join(root, "internal", "store")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "seed.go"), []byte(tc.source), 0o600); err != nil {
				t.Fatal(err)
			}
			violations, err := provenSignalViolations(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(violations) != len(tc.want) {
				t.Fatalf("violations=%q, want %q", violations, tc.want)
			}
			for index, want := range tc.want {
				if !strings.HasPrefix(violations[index], want) {
					t.Errorf("violation=%q, want %q", violations[index], want)
				}
			}
		})
	}
}
