// Suite: Parallel Test Packages.
// Invariant: every top-level test is parallel first or names why it must be sequential, within its package ceiling.
// Boundary IN: test declarations and leading comments in listed package directories.
// Boundary OUT: contract-tagged files, subtests, helpers and unlisted packages.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var parallelTestPackages = map[string]int{
	"internal/cli":       8,
	"internal/daemon":    16,
	"internal/baseline":  12,
	"internal/store":     3,
	"internal/spec":      2,
	"internal/speccheck": 2,
	"internal/worktree":  2,
}

type sequentialTest struct {
	file, name, reason string
	line               int
}

func parallelTestViolations(root string, ceilings map[string]int) ([]string, []sequentialTest, error) {
	var violations []string
	var sequential []sequentialTest
	var packages []string
	for pkg := range ceilings {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	for _, pkg := range packages {
		entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(pkg)))
		if err != nil {
			return nil, nil, fmt.Errorf("read Parallel Test Package %s: %w", pkg, err)
		}
		start := len(sequential)
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}
			relative := filepath.ToSlash(filepath.Join(pkg, entry.Name()))
			path := filepath.Join(root, filepath.FromSlash(relative))
			source, err := os.ReadFile(path)
			if err != nil {
				return nil, nil, fmt.Errorf("read %s: %w", relative, err)
			}
			if parallelTestContractFile(string(source)) {
				continue
			}
			set := token.NewFileSet()
			file, err := parser.ParseFile(set, path, source, parser.ParseComments)
			if err != nil {
				return nil, nil, fmt.Errorf("parse %s: %w", relative, err)
			}
			testingName := ""
			for _, imp := range file.Imports {
				if imp.Path.Value == `"testing"` {
					testingName = "testing"
					if imp.Name != nil {
						testingName = imp.Name.Name
					}
				}
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || fn.Body == nil || fn.Name.Name == "TestMain" || !strings.HasPrefix(fn.Name.Name, "Test") || fn.Type.Params.NumFields() != 1 {
					continue
				}
				param := fn.Type.Params.List[0]
				pointer, ok := param.Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				selector, ok := pointer.X.(*ast.SelectorExpr)
				if !ok || selector.Sel.Name != "T" {
					continue
				}
				importName, ok := selector.X.(*ast.Ident)
				if !ok || testingName == "" || importName.Name != testingName {
					continue
				}
				paramName := ""
				if len(param.Names) == 1 {
					paramName = param.Names[0].Name
				}
				parallel := false
				first := fn.Body.Rbrace
				if len(fn.Body.List) > 0 {
					first = fn.Body.List[0].Pos()
					parallel = parallelTestFirstStatement(fn.Body.List[0], paramName)
				}
				var comments []*ast.CommentGroup
				if fn.Doc != nil {
					comments = append(comments, fn.Doc)
				}
				for _, group := range file.Comments {
					if group.Pos() > fn.Body.Lbrace && group.End() <= first {
						comments = append(comments, group)
					}
				}
				reason, marked := parallelTestSequentialReason(comments)
				line := set.Position(fn.Pos()).Line
				site := fmt.Sprintf("%s:%d: %s", relative, line, fn.Name.Name)
				switch {
				case marked && len(strings.Fields(reason)) < 2:
					violations = append(violations, site+": Sequential: reason must contain at least two words")
				case parallel && marked:
					violations = append(violations, site+": both first-statement Parallel() and Sequential: reason")
				case !parallel && !marked:
					violations = append(violations, site+": missing first-statement Parallel() or Sequential: reason")
				case !parallel:
					sequential = append(sequential, sequentialTest{file: relative, name: fn.Name.Name, reason: reason, line: line})
				}
			}
		}
		if count := len(sequential) - start; count > ceilings[pkg] {
			site := sequential[start]
			violations = append(violations, fmt.Sprintf("%s:%d: %s: %s has %d Sequential Tests, ceiling %d", site.file, site.line, site.name, pkg, count, ceilings[pkg]))
		}
	}
	sort.Strings(violations)
	return violations, sequential, nil
}

func parallelTestContractFile(source string) bool {
	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			break
		}
		if !constraint.IsGoBuild(line) && !constraint.IsPlusBuild(line) {
			continue
		}
		expr, err := constraint.Parse(line)
		if err != nil {
			continue // The parser below reports malformed source.
		}
		if parallelTestContractConstraint(expr) {
			return true
		}
	}
	return false
}

func parallelTestContractConstraint(expr constraint.Expr) bool {
	switch expr := expr.(type) {
	case *constraint.TagExpr:
		return expr.Tag == "docscontract" || expr.Tag == "repocontract"
	case *constraint.NotExpr:
		return parallelTestContractConstraint(expr.X)
	case *constraint.AndExpr:
		return parallelTestContractConstraint(expr.X) || parallelTestContractConstraint(expr.Y)
	case *constraint.OrExpr:
		return parallelTestContractConstraint(expr.X) || parallelTestContractConstraint(expr.Y)
	}
	return false
}

func parallelTestFirstStatement(statement ast.Stmt, param string) bool {
	expression, ok := statement.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := expression.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 0 {
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Parallel" {
		return false
	}
	receiver, ok := selector.X.(*ast.Ident)
	return ok && param != "" && receiver.Name == param
}

func parallelTestSequentialReason(comments []*ast.CommentGroup) (string, bool) {
	for _, group := range comments {
		for _, line := range strings.Split(group.Text(), "\n") {
			if reason, ok := strings.CutPrefix(strings.TrimSpace(line), "Sequential:"); ok {
				return strings.TrimSpace(reason), true
			}
		}
	}
	return "", false
}

func TestEveryTestInAParallelTestPackageRunsInParallel(t *testing.T) {
	t.Parallel()
	t.Run("repository", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatal(err)
		}
		violations, sequential, err := parallelTestViolations(root, parallelTestPackages)
		if err != nil {
			t.Fatal(err)
		}
		var packages []string
		for pkg := range parallelTestPackages {
			packages = append(packages, pkg)
		}
		sort.Strings(packages)
		for _, pkg := range packages {
			count := 0
			for _, test := range sequential {
				if filepath.ToSlash(filepath.Dir(test.file)) == pkg {
					count++
				}
			}
			t.Logf("%s: %d Sequential Tests, ceiling %d", pkg, count, parallelTestPackages[pkg])
		}
		for _, violation := range violations {
			t.Error(violation)
		}
	})

	cases := []struct {
		name, prefix, declarations string
		ceiling                    int
		want                       []string
		sequential                 int
	}{
		{name: "parallel", declarations: "func TestSeed(subject *testing.T) { subject.Parallel() }", ceiling: 1},
		{name: "neither form", declarations: "func TestSeed(t *testing.T) {}", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:3: TestSeed: missing first-statement"}},
		{name: "one word reason", declarations: "// Sequential: environment\nfunc TestSeed(t *testing.T) {}", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:4: TestSeed: Sequential: reason must contain at least two words"}},
		{name: "both forms", declarations: "// Sequential: changes environment\nfunc TestSeed(t *testing.T) { t.Parallel() }", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:4: TestSeed: both"}},
		{name: "doc comment", declarations: "// Sequential: changes environment\nfunc TestSeed(t *testing.T) {}", ceiling: 1, sequential: 1},
		{name: "body comment", declarations: "func TestSeed(t *testing.T) {\n// Sequential: changes environment\nt.Helper()\n}", ceiling: 1, sequential: 1},
		{name: "later Parallel", declarations: "func TestSeed(t *testing.T) { t.Helper(); t.Parallel() }", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:3: TestSeed: missing first-statement"}},
		{name: "repocontract skipped", prefix: "//go:build repocontract\n\n", declarations: "func TestSeed(t *testing.T) {}", ceiling: 1},
		{name: "docscontract skipped", prefix: "//go:build linux && docscontract\n\n", declarations: "func TestSeed(t *testing.T) {}", ceiling: 1},
		{name: "legacy contract skipped", prefix: "// +build repocontract\n\n", declarations: "func TestSeed(t *testing.T) {}", ceiling: 1},
		{name: "over ceiling", declarations: "// Sequential: changes environment\nfunc TestSeed(t *testing.T) {}", ceiling: 0, sequential: 1, want: []string{"internal/seeded/fixture_test.go:4: TestSeed: internal/seeded has 1 Sequential Tests, ceiling 0"}},
		{name: "all violations", declarations: "func TestFirst(t *testing.T) {}\nfunc TestSecond(t *testing.T) {}", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:3: TestFirst:", "internal/seeded/fixture_test.go:4: TestSecond:"}},
		{name: "wrong receiver", declarations: "func TestSeed(subject *testing.T) { t.Parallel() }", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:3: TestSeed: missing first-statement"}},
		{name: "late reason", declarations: "func TestSeed(t *testing.T) { t.Helper()\n// Sequential: changes environment\n}", ceiling: 1, want: []string{"internal/seeded/fixture_test.go:3: TestSeed: missing first-statement"}},
		{name: "helper and TestMain ignored", declarations: "func helper(t *testing.T) {}\nfunc TestMain(m *testing.M) {}", ceiling: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			dir := filepath.Join(root, "internal", "seeded")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			source := tc.prefix + "package seeded\nimport \"testing\"\n" + tc.declarations + "\n"
			if err := os.WriteFile(filepath.Join(dir, "fixture_test.go"), []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			violations, sequential, err := parallelTestViolations(root, map[string]int{"internal/seeded": tc.ceiling})
			if err != nil {
				t.Fatal(err)
			}
			if len(violations) != len(tc.want) || len(sequential) != tc.sequential {
				t.Fatalf("violations = %q, Sequential Tests = %d; want %q, %d", violations, len(sequential), tc.want, tc.sequential)
			}
			for i, want := range tc.want {
				if !strings.Contains(violations[i], want) {
					t.Errorf("violation = %q, want %q", violations[i], want)
				}
			}
			for _, test := range sequential {
				if test.reason != "changes environment" || test.name != "TestSeed" || test.file != "internal/seeded/fixture_test.go" {
					t.Errorf("Sequential Test = %+v", test)
				}
			}
		})
	}
	t.Run("unlisted package ignored", func(t *testing.T) {
		t.Parallel()
		root := writtenExecutableTree(t, "package seeded\nimport \"testing\"\nfunc TestSeed(t *testing.T) {}\n")
		violations, sequential, err := parallelTestViolations(root, nil)
		if err != nil || len(violations) != 0 || len(sequential) != 0 {
			t.Fatalf("unlisted package: violations = %q, Sequential Tests = %+v, error = %v", violations, sequential, err)
		}
	})
}
