// Suite: production Git worktree administration routing.
// Invariant: production packages outside internal/worktree never administer Git worktrees directly.
// Boundary IN: Go call expressions in non-test files under internal and cmd.
// Boundary OUT: commands represented dynamically rather than by consecutive string literals.
package worktree

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

func TestNoProductionWorktreeAdministrationBypassesTheAdminLock(t *testing.T) {
	t.Parallel()
	repository := repositoryRootForAdminLockGuard(t)
	var findings []string
	for _, directory := range []string{"internal", "cmd"} {
		root := filepath.Join(repository, directory)
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				if path == filepath.Join(repository, "internal", "worktree") {
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			fileFindings, err := directWorktreeAdministrationCalls(path)
			if err != nil {
				return err
			}
			for _, finding := range fileFindings {
				relative, err := filepath.Rel(repository, finding.Filename)
				if err != nil {
					return err
				}
				findings = append(findings, fmt.Sprintf("%s:%d:%d", filepath.ToSlash(relative), finding.Line, finding.Column))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("scan production Go files under %s: %v", directory, err)
		}
	}
	if len(findings) != 0 {
		sort.Strings(findings)
		t.Fatalf("production Git worktree administration bypasses internal/worktree:\n%s", strings.Join(findings, "\n"))
	}
}

func TestAdminLockGuardReportsADirectWorktreeCall(t *testing.T) {
	t.Parallel()
	source := `package synthetic
func bypass(run func(...string)) {
	run("-C", ".", "worktree", "remove", "path")
}`
	path := filepath.Join(t.TempDir(), "direct.go")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write synthetic production file: %v", err)
	}

	findings, err := directWorktreeAdministrationCalls(path)
	if err != nil {
		t.Fatalf("inspect synthetic direct call: %v", err)
	}
	if len(findings) != 1 || findings[0].Line != 3 {
		t.Fatalf("direct-call findings = %+v, want one finding on line 3", findings)
	}
}

func directWorktreeAdministrationCalls(path string) ([]token.Position, error) {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	var findings []token.Position
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		for index := 0; index+1 < len(call.Args); index++ {
			if stringLiteral(call.Args[index]) != "worktree" {
				continue
			}
			switch stringLiteral(call.Args[index+1]) {
			case "add", "remove", "prune", "move":
				findings = append(findings, set.Position(call.Args[index].Pos()))
			}
		}
		return true
	})
	return findings, nil
}

func stringLiteral(expression ast.Expr) string {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return ""
	}
	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		return ""
	}
	return value
}

func repositoryRootForAdminLockGuard(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	return root
}
