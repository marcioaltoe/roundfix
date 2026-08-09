package store

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

func TestWaitBudgetIsNeverRestatedAtACallSite(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "agent", "acpx_runner_test.go"),
		"process_unix_test.go",
	}
	for _, path := range targets {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			fileSet := token.NewFileSet()
			file, err := parser.ParseFile(fileSet, path, nil, 0)
			if err != nil {
				t.Fatalf("parse wait-budget target %s: %v", path, err)
			}
			if violations := restatedWaitBudgets(fileSet, file); len(violations) > 0 {
				t.Fatalf("literal wait budget at call site in %s:\n%s", path, strings.Join(violations, "\n"))
			}
		})
	}
}

func TestWaitBudgetGuardRejectsRestatedLiteral(t *testing.T) {
	t.Parallel()

	const source = `package sample

import "time"

func wait() {
	<-time.After(5 * time.Second)
	_ = newOwnerProcessController(sharedGrace, 2 * time.Second, sharedPollInterval)
}
`
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "restated_wait_budget.go", source, 0)
	if err != nil {
		t.Fatalf("parse negative control: %v", err)
	}
	violations := restatedWaitBudgets(fileSet, file)
	if len(violations) != 2 {
		t.Fatalf("restated wait-budget violations = %v, want time.After and controller stop-window violations", violations)
	}
}

func restatedWaitBudgets(fileSet *token.FileSet, file *ast.File) []string {
	var violations []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		arguments := waitBudgetArguments(call)
		for _, argument := range arguments {
			if !isLiteralDuration(argument) {
				continue
			}
			position := fileSet.Position(argument.Pos())
			violations = append(violations, fmt.Sprintf("%s:%d", position.Filename, position.Line))
		}
		return true
	})
	return violations
}

func waitBudgetArguments(call *ast.CallExpr) []ast.Expr {
	switch function := call.Fun.(type) {
	case *ast.Ident:
		if function.Name == "newOwnerProcessController" {
			return call.Args
		}
	case *ast.SelectorExpr:
		packageName, ok := function.X.(*ast.Ident)
		if ok && packageName.Name == "time" && (function.Sel.Name == "After" || function.Sel.Name == "NewTimer") {
			return call.Args
		}
	}
	return nil
}

func isLiteralDuration(expression ast.Expr) bool {
	switch expression := expression.(type) {
	case *ast.ParenExpr:
		return isLiteralDuration(expression.X)
	case *ast.BinaryExpr:
		if expression.Op == token.MUL {
			return isIntegerLiteral(expression.X) && isTimeDurationUnit(expression.Y) ||
				isTimeDurationUnit(expression.X) && isIntegerLiteral(expression.Y)
		}
		return isLiteralDuration(expression.X) || isLiteralDuration(expression.Y)
	default:
		return isTimeDurationUnit(expression)
	}
}

func isIntegerLiteral(expression ast.Expr) bool {
	literal, ok := expression.(*ast.BasicLit)
	return ok && literal.Kind == token.INT
}

func isTimeDurationUnit(expression ast.Expr) bool {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	packageName, ok := selector.X.(*ast.Ident)
	if !ok || packageName.Name != "time" {
		return false
	}
	switch selector.Sel.Name {
	case "Nanosecond", "Microsecond", "Millisecond", "Second", "Minute", "Hour":
		return true
	default:
		return false
	}
}
