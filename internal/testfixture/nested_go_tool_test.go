// Suite: Nested Go Tool repository rule.
// Invariant: code tests do not start nested Go test or vet runs.
// Boundary IN: module test sources, exec calls and contract build constraints.
// Boundary OUT: testdata, production sources and the frozen migration inventory.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Frozen inventory: entries may only shrink.
var nestedGoToolInventory = map[string]string{
	"internal/verifyselect/verifyselect_test.go": "TestPartitionFollowsTheMakefileRecipes changes in task_03",
	"internal/spec/coverage_test.go":             "only the platform docscontract calls the toolchain listing",
}

func TestNoCodeTestStartsAGoTestOrGoVetRun(t *testing.T) {
	t.Parallel()
	reports, err := goToolReports("../..", false, nestedGoToolInventory)
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range reports {
		t.Error(report)
	}
}

func TestTheNestedGoToolRuleReportsSeededViolations(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, header, call string
		bad                bool
	}{
		{"test", "", `exec.Command("go", "test", "./other")`, true},
		{"vet", "", `exec.CommandContext(nil, "go", "vet", "./...")`, true},
		{"slice", "", `args := []string{"test", "./..."}; exec.Command("go", args...)`, true},
		{"variable", "", `tool := "go"; exec.Command(tool, "test")`, true},
		{"docs", "//go:build docscontract\n\n", `exec.Command("go", "test")`, false},
		{"repo", "//go:build repocontract\n\n", `exec.Command("go", "vet")`, false},
		{"optional contract", "//go:build docscontract || linux\n\n", `exec.Command("go", "test")`, true},
		{"negated platform", "//go:build !linux\n\n", `exec.Command("go", "test")`, true},
		{"list", "", `exec.Command("go", "list", "./...")`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "seed_test.go")
			source := tc.header + "package seeded\nimport \"os/exec\"\nfunc fixture(){" + tc.call + "}\n"
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			reports, err := goToolReports(root, false, nil)
			if err != nil {
				t.Fatal(err)
			}
			if (len(reports) > 0) != tc.bad {
				t.Fatalf("reports = %q, violation = %v", reports, tc.bad)
			}
			if tc.bad && !strings.Contains(reports[0], "seed_test.go:") {
				t.Fatalf("missing filename: %q", reports)
			}
		})
	}
}

func goToolReports(root string, builds bool, inventory map[string]string) ([]string, error) {
	var reports []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && (entry.Name() == "testdata" || entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "_")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if reason := inventory[rel]; reason != "" {
			return nil
		}
		if builds && strings.HasPrefix(rel, "internal/testfixture/") {
			return nil
		}
		set := token.NewFileSet()
		file, err := parser.ParseFile(set, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", rel, err)
		}
		if !builds && requiresContractTag(file) {
			return nil
		}
		imports := map[string]bool{}
		for _, imp := range file.Imports {
			if imp.Path.Value == `"os/exec"` {
				name := "exec"
				if imp.Name != nil {
					name = imp.Name.Name
				}
				imports[name] = true
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			name := ""
			switch fun := call.Fun.(type) {
			case *ast.SelectorExpr:
				pkg, ok := fun.X.(*ast.Ident)
				if ok && pkg.Obj == nil && imports[pkg.Name] {
					name = fun.Sel.Name
				}
			case *ast.Ident:
				if fun.Obj == nil && imports["."] {
					name = fun.Name
				}
			}
			start := 0
			if name == "CommandContext" {
				start = 1
			} else if name != "Command" {
				return true
			}
			var args []string
			for _, arg := range call.Args[start:] {
				args = append(args, goToolStrings(arg, map[*ast.Object]bool{})...)
			}
			if len(args) < 2 || (filepath.Base(args[0]) != "go" && filepath.Base(args[0]) != "go.exe") {
				return true
			}
			verb := args[1]
			compile := verb == "test" && containsGoToolArg(args[2:], "-c")
			bad := verb == "test" || verb == "vet"
			if builds {
				bad = verb == "build" || verb == "run" || compile
			}
			if bad {
				reports = append(reports, fmt.Sprintf("%s:%d: starts go %s outside the frozen inventory", rel, set.Position(call.Pos()).Line, verb))
			}
			return true
		})
		return nil
	})
	return reports, err
}

func containsGoToolArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want || strings.HasPrefix(arg, want+"=") {
			return true
		}
	}
	return false
}

func requiresContractTag(file *ast.File) bool {
	for _, group := range file.Comments {
		if group.Pos() > file.Package {
			break
		}
		for _, comment := range group.List {
			if !strings.HasPrefix(comment.Text, "//go:build ") {
				continue
			}
			expr, err := constraint.Parse(comment.Text)
			// If disabling both tags can satisfy the expression, it is still a code test.
			if err != nil {
				return false
			}
			tags := map[string]bool{}
			var collect func(constraint.Expr)
			collect = func(e constraint.Expr) {
				switch e := e.(type) {
				case *constraint.TagExpr:
					if e.Tag != "docscontract" && e.Tag != "repocontract" {
						tags[e.Tag] = false
					}
				case *constraint.NotExpr:
					collect(e.X)
				case *constraint.AndExpr:
					collect(e.X)
					collect(e.Y)
				case *constraint.OrExpr:
					collect(e.X)
					collect(e.Y)
				}
			}
			collect(expr)
			var names []string
			for tag := range tags {
				names = append(names, tag)
			}
			var codeBuildPossible func(int) bool
			codeBuildPossible = func(i int) bool {
				if i == len(names) {
					return expr.Eval(func(tag string) bool { return tags[tag] })
				}
				tags[names[i]] = false
				if codeBuildPossible(i + 1) {
					return true
				}
				tags[names[i]] = true
				return codeBuildPossible(i + 1)
			}
			return !codeBuildPossible(0)
		}
	}
	return false
}

// Resolve local literal command arguments, including variadic string slices.
func goToolStrings(expr ast.Expr, seen map[*ast.Object]bool) []string {
	switch value := expr.(type) {
	case *ast.BasicLit:
		if value.Kind == token.STRING {
			text, err := strconv.Unquote(value.Value)
			if err == nil {
				return []string{text}
			}
		}
	case *ast.CompositeLit:
		var result []string
		for _, element := range value.Elts {
			result = append(result, goToolStrings(element, seen)...)
		}
		return result
	case *ast.Ident:
		if value.Obj == nil || seen[value.Obj] {
			return nil
		}
		seen[value.Obj] = true
		switch decl := value.Obj.Decl.(type) {
		case *ast.AssignStmt:
			for i, lhs := range decl.Lhs {
				if id, ok := lhs.(*ast.Ident); ok && id.Name == value.Name && i < len(decl.Rhs) {
					return goToolStrings(decl.Rhs[i], seen)
				}
			}
		case *ast.ValueSpec:
			for i, id := range decl.Names {
				if id.Name == value.Name && i < len(decl.Values) {
					return goToolStrings(decl.Values[i], seen)
				}
			}
		}
	}
	return nil
}
