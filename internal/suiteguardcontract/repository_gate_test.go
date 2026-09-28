// Suite: repository-contract gate completeness.
// Invariant: every test requiring the repocontract build tag runs in the repository gate.
// Boundary IN: Go test declarations and repository-gate wiring in the module Makefile.
// Boundary OUT: behavior inside each repository-contract test.
package suiteguardcontract

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestEveryRepositoryContractTestRunsInTheRepositoryGate(t *testing.T) {
	t.Parallel()

	repositoryRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	findings, err := auditRepositoryGate(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("repository gate audit failed:\n%s", strings.Join(findings, "\n"))
	}
}

func TestRepositoryGateAuditNamesAMissingContractTest(t *testing.T) {
	t.Parallel()

	repositoryRoot := newRepositoryGateFixture(t, "TestSomeOtherContract", "./...")
	findings, err := auditRepositoryGate(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	assertRepositoryGateFindingContains(t, findings, "TestFixtureRepositoryContract")
}

func TestRepositoryGateAuditNamesADroppedRepositoryWideRun(t *testing.T) {
	t.Parallel()

	repositoryRoot := newRepositoryGateFixture(t, "TestFixtureRepositoryContract", "./internal/example")
	findings, err := auditRepositoryGate(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	assertRepositoryGateFindingContains(t, findings, "./...")
}

func auditRepositoryGate(repositoryRoot string) ([]string, error) {
	contractTests, err := repositoryContractTests(repositoryRoot)
	if err != nil {
		return nil, err
	}
	makefile, err := os.ReadFile(filepath.Join(repositoryRoot, "Makefile"))
	if err != nil {
		return nil, fmt.Errorf("read repository Makefile: %w", err)
	}

	var findings []string
	listedTests, found := makeVariable(makefile, "REPO_CONTRACT_TESTS")
	if !found {
		findings = append(findings, "Makefile does not define REPO_CONTRACT_TESTS")
	}
	listed := make(map[string]struct{})
	for _, name := range strings.Split(listedTests, "|") {
		if name = strings.TrimSpace(name); name != "" {
			listed[name] = struct{}{}
		}
	}
	for _, name := range contractTests {
		if _, ok := listed[name]; !ok {
			findings = append(findings, fmt.Sprintf("repository-contract test %s is missing from REPO_CONTRACT_TESTS", name))
		}
	}

	repoTestRecipe, found := makeTargetRecipe(makefile, "repo-test")
	if !found {
		findings = append(findings, "Makefile has no repo-test recipe")
	} else {
		if !recipePassesBuildTag(repoTestRecipe, "repocontract") {
			findings = append(findings, "repo-test must pass -tags repocontract")
		}
		if !containsMakeToken(repoTestRecipe, "./...") {
			findings = append(findings, "repo-test must pass ./... to run repository-contract tests across the module")
		}
	}

	verifyDocsDependencies, found := makeTargetDependencies(makefile, "verify-docs")
	if !found || !containsGateString(verifyDocsDependencies, "repo-test") {
		findings = append(findings, "verify-docs must depend on repo-test")
	}

	sort.Strings(findings)
	return findings, nil
}

func repositoryContractTests(repositoryRoot string) ([]string, error) {
	tests := make(map[string]struct{})
	err := filepath.WalkDir(repositoryRoot, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if filePath == repositoryRoot {
				return nil
			}
			name := entry.Name()
			if name == "testdata" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}

		requires, err := fileRequiresBuildTag(filePath, "repocontract")
		if err != nil {
			return err
		}
		if !requires {
			return nil
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), filePath, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", filePath, err)
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if ok && function.Recv == nil && isGoTestName(function.Name.Name) {
				tests[function.Name.Name] = struct{}{}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(tests))
	for name := range tests {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func fileRequiresBuildTag(filePath, tag string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("open %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "package ") {
			break
		}
		if !strings.HasPrefix(line, "//go:build ") {
			continue
		}
		expression, err := constraint.Parse(line)
		if err != nil {
			return false, fmt.Errorf("parse build constraint in %s: %w", filePath, err)
		}
		return buildConstraintRequiresTag(expression, tag), nil
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("read %s: %w", filePath, err)
	}
	return false, nil
}

func buildConstraintRequiresTag(expression constraint.Expr, requiredTag string) bool {
	tags := make(map[string]struct{})
	collectBuildConstraintTags(expression, tags)
	if _, found := tags[requiredTag]; !found {
		return false
	}
	delete(tags, requiredTag)
	otherTags := make([]string, 0, len(tags))
	for tag := range tags {
		otherTags = append(otherTags, tag)
	}

	values := map[string]bool{requiredTag: false}
	var canMatchWithoutRequiredTag func(int) bool
	canMatchWithoutRequiredTag = func(index int) bool {
		if index == len(otherTags) {
			return expression.Eval(func(tag string) bool { return values[tag] })
		}
		tag := otherTags[index]
		values[tag] = false
		if canMatchWithoutRequiredTag(index + 1) {
			return true
		}
		values[tag] = true
		return canMatchWithoutRequiredTag(index + 1)
	}
	return !canMatchWithoutRequiredTag(0)
}

func collectBuildConstraintTags(expression constraint.Expr, tags map[string]struct{}) {
	switch expression := expression.(type) {
	case *constraint.TagExpr:
		tags[expression.Tag] = struct{}{}
	case *constraint.NotExpr:
		collectBuildConstraintTags(expression.X, tags)
	case *constraint.AndExpr:
		collectBuildConstraintTags(expression.X, tags)
		collectBuildConstraintTags(expression.Y, tags)
	case *constraint.OrExpr:
		collectBuildConstraintTags(expression.X, tags)
		collectBuildConstraintTags(expression.Y, tags)
	}
}

func isGoTestName(name string) bool {
	if !strings.HasPrefix(name, "Test") || len(name) == len("Test") {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name[len("Test"):])
	return !unicode.IsLower(r)
}

func makeVariable(makefile []byte, name string) (string, bool) {
	lines := strings.Split(string(makefile), "\n")
	prefix := name + " :="
	for index, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		parts := []string{strings.TrimSpace(strings.TrimPrefix(line, prefix))}
		for strings.HasSuffix(parts[len(parts)-1], "\\") && index+1 < len(lines) {
			parts[len(parts)-1] = strings.TrimSpace(strings.TrimSuffix(parts[len(parts)-1], "\\"))
			index++
			parts = append(parts, strings.TrimSpace(lines[index]))
		}
		return strings.Join(parts, " "), true
	}
	return "", false
}

func makeTargetRecipe(makefile []byte, target string) (string, bool) {
	lines := strings.Split(string(makefile), "\n")
	for index, line := range lines {
		if !strings.HasPrefix(line, target+":") {
			continue
		}
		var recipe []string
		for index++; index < len(lines) && strings.HasPrefix(lines[index], "\t"); index++ {
			recipe = append(recipe, strings.TrimSpace(lines[index]))
		}
		return strings.Join(recipe, " "), true
	}
	return "", false
}

func makeTargetDependencies(makefile []byte, target string) ([]string, bool) {
	prefix := target + ":"
	for _, line := range strings.Split(string(makefile), "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.Fields(strings.SplitN(strings.TrimPrefix(line, prefix), "##", 2)[0]), true
		}
	}
	return nil, false
}

func recipePassesBuildTag(recipe, tag string) bool {
	fields := strings.Fields(recipe)
	for index, field := range fields {
		if field == "-tags" && index+1 < len(fields) && fields[index+1] == tag {
			return true
		}
		if field == "-tags="+tag {
			return true
		}
	}
	return false
}

func containsMakeToken(recipe, wanted string) bool {
	return containsGateString(strings.Fields(recipe), wanted)
}

func containsGateString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func newRepositoryGateFixture(t *testing.T, listedTest, packagePattern string) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	writeRepositoryGateFixtureFile(t, filepath.Join(repositoryRoot, "go.mod"), "module fixture\n\ngo 1.26\n")
	writeRepositoryGateFixtureFile(t, filepath.Join(repositoryRoot, "internal", "example", "contract_repocontract_test.go"), `//go:build repocontract

package example

import "testing"

func TestFixtureRepositoryContract(t *testing.T) {}
`)
	makefile := fmt.Sprintf(`REPO_CONTRACT_TESTS := %s

verify-docs: repo-test

repo-test:
	go test -count=1 -tags repocontract -run '^($(REPO_CONTRACT_TESTS))$$' %s
`, listedTest, packagePattern)
	writeRepositoryGateFixtureFile(t, filepath.Join(repositoryRoot, "Makefile"), makefile)
	return repositoryRoot
}

func writeRepositoryGateFixtureFile(t *testing.T, filePath, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertRepositoryGateFindingContains(t *testing.T, findings []string, wanted string) {
	t.Helper()
	for _, finding := range findings {
		if strings.Contains(finding, wanted) {
			return
		}
	}
	t.Fatalf("findings = %q, want one naming %q", findings, wanted)
}
