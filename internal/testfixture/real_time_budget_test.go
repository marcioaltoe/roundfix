// Suite: Real-Time Budgets.
// Invariant: each test file stays within its ceiling of short constant budgets.
// Boundary IN: timeout calls, named fields and assignments, and same-package budget parameters.
// Boundary OUT: production constants, dynamic durations, sleeps and excluded fixture/tool directories.
package testfixture

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

const realTimeBudgetThreshold = 10 * time.Second

type realTimeBudgetCeiling struct {
	max    int
	reason string
}

// A ceiling only falls. When a file loses a budget, lower its ceiling here.
var realTimeBudgetCeilings = map[string]realTimeBudgetCeiling{
	"internal/agent/proof_close_test.go":         {1, "The cleanup deadline is asserted against a blocked sessions close."},
	"internal/agent/acpx_runner_test.go":         {1, "StopGrace is unused by the successful prompt; it bounds no cancellation child."},
	"internal/agent/cursor_login_test.go":        {1, "The caller deadline is asserted against a login fixture that blocks."},
	"internal/cli/cli_test.go":                   {2, "Review and check deadlines are data for the fake watch clock and sleeper."},
	"internal/cli/detach_unix_test.go":           {2, "The asserted liveness and Run-creation deadlines race owner-watched blocking child modes."},
	"internal/cli/gc_run_retention_test.go":      {3, "Sweep budgets are measured by the injected retention clock."},
	"internal/cli/orphan_unix_test.go":           {1, "Teardown reaps an already killed or exited owner; it bounds no child startup."},
	"internal/cli/readiness_forge_test.go":       {1, "The probe deadline is asserted against a child that never answers before cancellation."},
	"internal/daemon/task_engine_test.go":        {4, "Bootstrap timeouts are data for fake worktrees, the disabled Run Budget bounds no work, and the BudgetExceeded deadline is asserted against blocked task_02 (its setup fix belongs to ADR-0259's delivery)."},
	"internal/delivery/check_rerun_test.go":      {1, "CheckTimeout is data for the fake delivery clock and sleeper."},
	"internal/delivery/conflict_test.go":         {2, "CheckTimeout is data for the fake delivery clock and sleeper."},
	"internal/delivery/engine_test.go":           {2, "Check timeouts are data for the fake delivery clock and sleeper."},
	"internal/delivery/merge_state_test.go":      {1, "CheckTimeout is data for the fake delivery clock and sleeper."},
	"internal/notify/notify_test.go":             {3, "Immediate fake runners bound no child or goroutine; the timeout case asserts expiry against a runner blocked on context cancellation."},
	"internal/runevent/fanout_test.go":           {1, "Publication's bound is asserted while the best-effort sink blocks until release."},
	"internal/store/journal_test.go":             {1, "The prune deadline is asserted against an independently held write lock."},
	"internal/store/process_signal_unix_test.go": {1, "The grace expiry is asserted before the identity re-proof; the graceful signal and liveness probe are fakes, so it bounds no child."},
	"internal/store/process_unix_test.go":        {18, "Stop grace applies to ready or fake owners, never startup; other paths reject identity or accept proven absence before waiting."},
	"internal/store/writetx_test.go":             {1, "The second writer's deadline is asserted after the first writer holds the lock and blocks on cancellation."},
	"internal/watch/watch_test.go":               {15, "Review, check and Run budgets are data for fake clocks and sleepers."},
	"internal/worktree/adminlock_test.go":        {1, "The helper deadline is asserted against the parent's held admin lock with an in-process fake Git runner."},
	"internal/worktree/worktree_test.go":         {2, "Bootstrap expiry is asserted against a shell blocked on a FIFO; the empty-command timeout bounds no work."},
}

var realTimeBudgetName = regexp.MustCompile(`(?i)timeout|deadline|grace|budget|liveness|runcreation|maxrunduration`)

type realTimeBudget struct {
	file     string
	line     int
	duration time.Duration
	form     string
}

type realTimeBudgetSource struct {
	file    *ast.File
	path    string
	imports map[string]string
}

type realTimeBudgetConstant struct {
	expression ast.Expr
	source     *realTimeBudgetSource
}

type realTimeBudgetPackage struct {
	constants  map[string]*realTimeBudgetConstant
	locals     map[*ast.Object]*realTimeBudgetConstant
	parameters map[string][]int
	sources    []*realTimeBudgetSource
}

func realTimeBudgets(root string) ([]realTimeBudget, error) {
	set := token.NewFileSet()
	packages := map[string]*realTimeBudgetPackage{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".gocache", "bin", "testdata", "node_modules", ".agents", "skills":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(set, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("resolve %s: %w", path, err)
		}
		key := filepath.Dir(path) + "/" + file.Name.Name
		pkg := packages[key]
		if pkg == nil {
			pkg = &realTimeBudgetPackage{
				constants:  map[string]*realTimeBudgetConstant{},
				locals:     map[*ast.Object]*realTimeBudgetConstant{},
				parameters: map[string][]int{},
			}
			packages[key] = pkg
		}
		// Production functions supply parameter names, but production constants do not supply budgets.
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv != nil {
				continue
			}
			index := 0
			for _, field := range function.Type.Params.List {
				for _, name := range field.Names {
					if realTimeBudgetName.MatchString(name.Name) {
						pkg.parameters[function.Name.Name] = append(pkg.parameters[function.Name.Name], index)
					}
					index++
				}
				if len(field.Names) == 0 {
					index++
				}
			}
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		source := &realTimeBudgetSource{file: file, path: filepath.ToSlash(relative), imports: map[string]string{}}
		for _, imp := range file.Imports {
			name := strings.Trim(imp.Path.Value, `"`)
			if name != "time" && name != "context" {
				continue
			}
			alias := name
			if imp.Name != nil {
				alias = imp.Name.Name
			}
			source.imports[alias] = name
		}
		pkg.sources = append(pkg.sources, source)
		ast.Inspect(file, func(node ast.Node) bool {
			declaration, ok := node.(*ast.GenDecl)
			if !ok || declaration.Tok != token.CONST {
				return true
			}
			var previous []ast.Expr
			for _, spec := range declaration.Specs {
				value := spec.(*ast.ValueSpec)
				if len(value.Values) != 0 {
					previous = value.Values
				}
				for index, name := range value.Names {
					if index >= len(previous) {
						continue
					}
					binding := &realTimeBudgetConstant{expression: previous[index], source: source}
					pkg.locals[name.Obj] = binding
					if file.Scope.Lookup(name.Name) == name.Obj {
						pkg.constants[name.Name] = binding
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan Real-Time Budgets: %w", err)
	}
	var budgets []realTimeBudget
	for _, pkg := range packages {
		for _, source := range pkg.sources {
			record := func(expression ast.Expr, form string) {
				value := realTimeBudgetValue(expression, source, pkg, map[*realTimeBudgetConstant]bool{})
				nanos, exact := constant.Int64Val(constant.ToInt(value))
				if exact && nanos > 0 && nanos < int64(realTimeBudgetThreshold) {
					budgets = append(budgets, realTimeBudget{source.path, set.Position(expression.Pos()).Line, time.Duration(nanos), form})
				}
			}
			ast.Inspect(source.file, func(node ast.Node) bool {
				switch node := node.(type) {
				case *ast.CallExpr:
					if selector, ok := node.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := selector.X.(*ast.Ident); ok && ident.Obj == nil {
							library := source.imports[ident.Name]
							index := -1
							switch library + "." + selector.Sel.Name {
							case "context.WithTimeout":
								index = 1
							case "time.After", "time.NewTimer", "time.AfterFunc":
								index = 0
							}
							if index >= 0 && index < len(node.Args) {
								record(node.Args[index], library+"."+selector.Sel.Name)
							}
						}
					}
					if function, ok := node.Fun.(*ast.Ident); ok && (function.Obj == nil || function.Obj.Kind == ast.Fun) {
						// Alternate build-tagged declarations can contribute the same parameter index.
						seen := map[int]bool{}
						for _, index := range pkg.parameters[function.Name] {
							if index < len(node.Args) && !seen[index] {
								record(node.Args[index], "parameter of "+function.Name)
								seen[index] = true
							}
						}
					}
				case *ast.KeyValueExpr:
					if name, ok := node.Key.(*ast.Ident); ok && realTimeBudgetName.MatchString(name.Name) {
						record(node.Value, "field "+name.Name)
					}
				case *ast.AssignStmt:
					for index, left := range node.Lhs {
						if index < len(node.Rhs) {
							if name := realTimeBudgetAssignmentName(left); realTimeBudgetName.MatchString(name) {
								record(node.Rhs[index], "assignment "+name)
							}
						}
					}
				case *ast.GenDecl:
					if node.Tok == token.VAR {
						for _, spec := range node.Specs {
							value := spec.(*ast.ValueSpec)
							for index, name := range value.Names {
								if index < len(value.Values) && realTimeBudgetName.MatchString(name.Name) {
									record(value.Values[index], "assignment "+name.Name)
								}
							}
						}
					}
				}
				return true
			})
		}
	}
	sort.Slice(budgets, func(i, j int) bool {
		if budgets[i].file != budgets[j].file {
			return budgets[i].file < budgets[j].file
		}
		if budgets[i].line != budgets[j].line {
			return budgets[i].line < budgets[j].line
		}
		if budgets[i].form != budgets[j].form {
			return budgets[i].form < budgets[j].form
		}
		return budgets[i].duration < budgets[j].duration
	})
	return budgets, nil
}

func realTimeBudgetAssignmentName(expression ast.Expr) string {
	switch expression := expression.(type) {
	case *ast.Ident:
		return expression.Name
	case *ast.SelectorExpr:
		return expression.Sel.Name
	}
	return ""
}

func realTimeBudgetValue(expression ast.Expr, source *realTimeBudgetSource, pkg *realTimeBudgetPackage, resolving map[*realTimeBudgetConstant]bool) constant.Value {
	unknown := constant.MakeUnknown()
	switch expression := expression.(type) {
	case *ast.BasicLit:
		if expression.Kind == token.INT || expression.Kind == token.FLOAT {
			return constant.MakeFromLiteral(expression.Value, expression.Kind, 0)
		}
	case *ast.ParenExpr:
		return realTimeBudgetValue(expression.X, source, pkg, resolving)
	case *ast.UnaryExpr:
		if expression.Op == token.ADD || expression.Op == token.SUB {
			return constant.UnaryOp(expression.Op, realTimeBudgetValue(expression.X, source, pkg, resolving), 0)
		}
	case *ast.BinaryExpr:
		switch expression.Op {
		case token.ADD, token.SUB, token.MUL, token.QUO:
			left := realTimeBudgetValue(expression.X, source, pkg, resolving)
			right := realTimeBudgetValue(expression.Y, source, pkg, resolving)
			if left.Kind() == constant.Unknown || right.Kind() == constant.Unknown {
				return unknown
			}
			if expression.Op == token.QUO && constant.Sign(right) == 0 {
				return unknown
			}
			op := expression.Op
			if op == token.QUO && left.Kind() == constant.Int && right.Kind() == constant.Int {
				op = token.QUO_ASSIGN
			}
			return constant.BinaryOp(left, op, right)
		}
	case *ast.SelectorExpr:
		if name, ok := expression.X.(*ast.Ident); ok && name.Obj == nil && source.imports[name.Name] == "time" {
			units := map[string]time.Duration{"Nanosecond": time.Nanosecond, "Microsecond": time.Microsecond, "Millisecond": time.Millisecond, "Second": time.Second, "Minute": time.Minute, "Hour": time.Hour}
			if value, ok := units[expression.Sel.Name]; ok {
				return constant.MakeInt64(int64(value))
			}
		}
	case *ast.CallExpr:
		if selector, ok := expression.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Duration" && len(expression.Args) == 1 {
			if name, ok := selector.X.(*ast.Ident); ok && name.Obj == nil && source.imports[name.Name] == "time" {
				return constant.ToInt(realTimeBudgetValue(expression.Args[0], source, pkg, resolving))
			}
		}
	case *ast.Ident:
		var binding *realTimeBudgetConstant
		if expression.Obj != nil {
			binding = pkg.locals[expression.Obj]
		} else {
			binding = pkg.constants[expression.Name]
		}
		if binding != nil && !resolving[binding] {
			resolving[binding] = true
			value := realTimeBudgetValue(binding.expression, binding.source, pkg, resolving)
			delete(resolving, binding)
			return value
		}
	}
	return unknown
}

func realTimeBudgetViolations(root string, ceilings map[string]realTimeBudgetCeiling) ([]string, error) {
	budgets, err := realTimeBudgets(root)
	if err != nil {
		return nil, err
	}
	return realTimeBudgetFindings(budgets, ceilings), nil
}

func realTimeBudgetFindings(budgets []realTimeBudget, ceilings map[string]realTimeBudgetCeiling) []string {
	byFile := map[string][]realTimeBudget{}
	for _, budget := range budgets {
		byFile[budget.file] = append(byFile[budget.file], budget)
	}
	var files []string
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)
	var violations []string
	for _, file := range files {
		sites := byFile[file]
		ceiling := ceilings[file]
		if len(sites) <= ceiling.max {
			continue
		}
		violations = append(violations, fmt.Sprintf("%s: %d Real-Time Budgets, ceiling %d", file, len(sites), ceiling.max))
		for _, site := range sites {
			violations = append(violations, fmt.Sprintf("%s:%d: Real-Time Budget %s in %s", file, site.line, site.duration, site.form))
		}
	}
	return violations
}

func TestEveryRealTimeBudgetIsWithinItsFileCeiling(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	budgets, err := realTimeBudgets(root)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for _, budget := range budgets {
		counts[budget.file]++
	}
	for file, ceiling := range realTimeBudgetCeilings {
		if strings.TrimSpace(ceiling.reason) == "" {
			t.Errorf("%s: Real-Time Budget ceiling needs a reason", file)
		}
		if count := counts[file]; count < ceiling.max {
			t.Logf("%s: %d Real-Time Budgets below ceiling %d; lower the ceiling on its next change", file, count, ceiling.max)
		}
	}
	if violations := realTimeBudgetFindings(budgets, realTimeBudgetCeilings); len(violations) != 0 {
		t.Fatal(strings.Join(violations, "\n"))
	}
}

func TestRealTimeBudgetRuleRejectsSeededViolations(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		declarations string
		companion    string
		production   bool
		ceiling      *realTimeBudgetCeiling
		want         []string
	}{
		{name: "unlisted-timeout", declarations: "func fixture() { context.WithTimeout(context.Background(), 50 * time.Millisecond) }", want: []string{"internal/seeded/fixture_test.go: 1 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:3: Real-Time Budget 50ms in context.WithTimeout"}},
		{name: "over-ceiling", declarations: "func fixture() {\n time.After(time.Second)\n time.NewTimer(2 * time.Second)\n}", ceiling: &realTimeBudgetCeiling{1, "asserted against blocking work"}, want: []string{"internal/seeded/fixture_test.go: 2 Real-Time Budgets, ceiling 1", "internal/seeded/fixture_test.go:4: Real-Time Budget 1s in time.After", "internal/seeded/fixture_test.go:5: Real-Time Budget 2s in time.NewTimer"}},
		{name: "parameter-budget", declarations: "func fixture() { wait(50 * time.Millisecond) }", companion: "package seeded\nimport \"time\"\nfunc wait(grace time.Duration) {}\n", want: []string{"internal/seeded/fixture_test.go: 1 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:3: Real-Time Budget 50ms in parameter of wait"}},
		{name: "named-constant", declarations: "func fixture() { time.After(shortBudget) }", companion: "package seeded\nimport clock \"time\"\nconst shortBudget = clock.Duration(100) * clock.Millisecond / 2\n", want: []string{"internal/seeded/fixture_test.go: 1 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:3: Real-Time Budget 50ms in time.After"}},
		{name: "long-budget", declarations: "func fixture() { context.WithTimeout(context.Background(), 10 * time.Second); time.After(time.Hour) }"},
		{name: "sleep-only", declarations: "func fixture() { time.Sleep(time.Millisecond) }"},
		{name: "listed", declarations: "func fixture() { time.After(time.Second) }", ceiling: &realTimeBudgetCeiling{1, "asserted against blocking work"}},
		{name: "fields-and-assignments", declarations: "var deadline = time.Second\nfunc fixture() {\n cfg := struct{ Budget time.Duration }{Budget: 2 * time.Second}\n cfg.Budget = time.Duration(3) * time.Second\n}", want: []string{"internal/seeded/fixture_test.go: 3 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:3: Real-Time Budget 1s in assignment deadline", "internal/seeded/fixture_test.go:5: Real-Time Budget 2s in field Budget", "internal/seeded/fixture_test.go:6: Real-Time Budget 3s in assignment Budget"}},
		{name: "constant-arithmetic", declarations: "func fixture() {\n const (\n  allowance = +(2 * time.Second + time.Second - 500 * time.Millisecond) / 2\n  inherited\n )\n time.AfterFunc(inherited, func(){})\n}", want: []string{"internal/seeded/fixture_test.go: 1 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:8: Real-Time Budget 1.25s in time.AfterFunc"}},
		{name: "production-constant", declarations: "func fixture() { time.After(shortBudget) }", companion: "package seeded\nimport \"time\"\nconst shortBudget = time.Millisecond\n", production: true},
		{name: "production-parameter", declarations: "func fixture() { wait(50 * time.Millisecond) }", companion: "package seeded\nimport \"time\"\nfunc wait(grace time.Duration) {}\n", production: true, want: []string{"internal/seeded/fixture_test.go: 1 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:3: Real-Time Budget 50ms in parameter of wait"}},
		{name: "dynamic-duration", declarations: "func fixture(duration time.Duration) { time.After(duration) }"},
		{name: "shadowed-constant", declarations: "const allowance = time.Second\nfunc fixture(allowance time.Duration) { time.After(allowance) }"},
		{name: "below-ceiling", declarations: "func fixture() {}", ceiling: &realTimeBudgetCeiling{1, "a concurrent change removed an asserted budget"}},
		{name: "numeric-expressions", declarations: "func fixture() {\n time.After(0.5 * time.Second)\n time.NewTimer(time.Duration(3 / 2) * time.Second)\n}", want: []string{"internal/seeded/fixture_test.go: 2 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:4: Real-Time Budget 500ms in time.After", "internal/seeded/fixture_test.go:5: Real-Time Budget 1s in time.NewTimer"}},
		{name: "grouped-parameters", declarations: "func wait(unused int, timeout, grace time.Duration) {}\nfunc fixture() { wait(0, time.Second, 2 * time.Second) }", want: []string{"internal/seeded/fixture_test.go: 2 Real-Time Budgets, ceiling 0", "internal/seeded/fixture_test.go:4: Real-Time Budget 1s in parameter of wait", "internal/seeded/fixture_test.go:4: Real-Time Budget 2s in parameter of wait"}},
		{name: "other-package-constant", declarations: "func fixture() { time.After(shortBudget) }", companion: "package seeded_test\nimport \"time\"\nconst shortBudget = time.Millisecond\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := writtenExecutableTree(t, "package seeded\nimport (\"context\"; \"time\")\n"+tc.declarations+"\n")
			if tc.companion != "" {
				companion := "companion_test.go"
				if tc.production {
					companion = "companion.go"
				}
				if err := os.WriteFile(filepath.Join(root, "internal", "seeded", companion), []byte(tc.companion), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			ceilings := map[string]realTimeBudgetCeiling{}
			if tc.ceiling != nil {
				ceilings["internal/seeded/fixture_test.go"] = *tc.ceiling
			}
			violations, err := realTimeBudgetViolations(root, ceilings)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Join(violations, "\n") != strings.Join(tc.want, "\n") {
				t.Fatalf("violations = %q, want %q", violations, tc.want)
			}
		})
	}
	for _, directory := range []string{".git", ".gocache", "bin", "testdata", "node_modules", ".agents", "skills"} {
		t.Run("excluded-"+directory, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			dir := filepath.Join(root, "nested", directory)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "fixture_test.go"), []byte("package seeded\nimport \"time\"\nfunc fixture() { time.After(time.Second) }\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			violations, err := realTimeBudgetViolations(root, nil)
			if err != nil || len(violations) != 0 {
				t.Fatalf("excluded tree: violations = %q, error = %v", violations, err)
			}
		})
	}
}
