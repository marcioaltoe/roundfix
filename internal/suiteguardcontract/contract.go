//go:build repocontract

// Package suiteguardcontract audits package-level suiteguard installation.
package suiteguardcontract

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var guardedSpawningPackages = []string{
	"internal/agent",
	"internal/baseline",
	"internal/cli",
	"internal/daemon",
	"internal/gittest",
	"internal/preflight",
	"internal/spec",
	"internal/specaudit",
	"internal/speccheck",
	"internal/store",
	"internal/suiteguard",
	"internal/worktree",
}

// GuardedSpawningPackages returns the packages the repository contract expects
// to spawn processes and install suiteguard.Main.
func GuardedSpawningPackages() []string {
	return append([]string(nil), guardedSpawningPackages...)
}

// CheckCurrentPackage verifies the calling package's part of the repository
// contract. Every test-bearing internal package calls it under repocontract so
// a newly introduced process spawn is named at its owning package.
func CheckCurrentPackage(t *testing.T) {
	t.Helper()
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolve package directory: %v", err)
	}
	repositoryRoot, err := findRepositoryRoot(workingDirectory)
	if err != nil {
		t.Fatal(err)
	}
	packageDirectory, err := filepath.Rel(repositoryRoot, workingDirectory)
	if err != nil {
		t.Fatalf("make package directory relative to repository root: %v", err)
	}
	packagePath := filepath.ToSlash(packageDirectory)

	facts, err := scanPackageTestFacts(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	findings := auditPackage(packagePath, facts[packagePath], guardedSpawningPackages)
	if len(findings) != 0 {
		t.Fatalf("suite guard installation contract failed:\n%s", strings.Join(findings, "\n"))
	}
}

// AuditInstallations compares the discovered spawning packages and guard calls
// with the repository's explicit guarded-package inventory.
func AuditInstallations(root string, enumerated []string) ([]string, error) {
	facts, err := scanPackageTestFacts(root)
	if err != nil {
		return nil, err
	}

	guarded := make(map[string]struct{}, len(enumerated))
	for _, packagePath := range enumerated {
		guarded[packagePath] = struct{}{}
	}

	var findings []string
	for packagePath, packageFacts := range facts {
		if !packageFacts.spawnsProcess {
			continue
		}
		findings = append(findings, auditPackage(packagePath, packageFacts, enumerated)...)
	}
	for packagePath := range guarded {
		packageFacts, ok := facts[packagePath]
		if !ok || !packageFacts.spawnsProcess {
			findings = append(findings, packagePath+" is enumerated as guarded but has no test subprocess call")
		}
	}
	sort.Strings(findings)
	return findings, nil
}

type packageTestFacts struct {
	spawnsProcess bool
	installsGuard bool
}

func auditPackage(packagePath string, facts packageTestFacts, enumerated []string) []string {
	if !facts.spawnsProcess {
		return nil
	}
	var findings []string
	if !containsString(enumerated, packagePath) {
		findings = append(findings, packagePath+" spawns a process but is not enumerated as guarded")
	}
	if !facts.installsGuard {
		findings = append(findings, packagePath+" spawns a process but does not install suiteguard.Main")
	}
	return findings
}

func scanPackageTestFacts(root string) (map[string]packageTestFacts, error) {
	internalRoot := filepath.Join(root, "internal")
	facts := make(map[string]packageTestFacts)
	err := filepath.WalkDir(internalRoot, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), filePath, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", filePath, err)
		}
		directory, err := filepath.Rel(root, filepath.Dir(filePath))
		if err != nil {
			return fmt.Errorf("make %s relative to %s: %w", filePath, root, err)
		}
		packagePath := filepath.ToSlash(directory)
		packageFacts := facts[packagePath]
		spawns, installs := inspectTestFile(parsed)
		packageFacts.spawnsProcess = packageFacts.spawnsProcess || spawns
		packageFacts.installsGuard = packageFacts.installsGuard || installs
		facts[packagePath] = packageFacts
		return nil
	})
	if err != nil {
		return nil, err
	}
	return facts, nil
}

func inspectTestFile(file *ast.File) (bool, bool) {
	imports := make(map[string]string)
	dotImports := make(map[string]bool)
	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			continue
		}
		name := pathpkg.Base(importPath)
		if importSpec.Name != nil {
			name = importSpec.Name.Name
		}
		if name == "." {
			dotImports[importPath] = true
			continue
		}
		imports[name] = importPath
	}

	var spawnsProcess bool
	var installsGuard bool
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch function := call.Fun.(type) {
		case *ast.SelectorExpr:
			qualifier, ok := function.X.(*ast.Ident)
			if !ok {
				return true
			}
			importPath := imports[qualifier.Name]
			spawnsProcess = spawnsProcess || isProcessSpawn(importPath, function.Sel.Name)
			installsGuard = installsGuard || isSuiteGuardMain(importPath, function.Sel.Name)
		case *ast.Ident:
			spawnsProcess = spawnsProcess || isDotImportedProcessSpawn(dotImports, function.Name)
			installsGuard = installsGuard || (function.Name == "Main" && hasSuiteGuardDotImport(dotImports))
		}
		return true
	})
	return spawnsProcess, installsGuard
}

func isProcessSpawn(importPath, function string) bool {
	switch importPath {
	case "os/exec":
		return function == "Command" || function == "CommandContext"
	case "os":
		return function == "StartProcess"
	case "syscall":
		return function == "ForkExec"
	default:
		return false
	}
}

func isDotImportedProcessSpawn(dotImports map[string]bool, function string) bool {
	return (dotImports["os/exec"] && (function == "Command" || function == "CommandContext")) ||
		(dotImports["os"] && function == "StartProcess") ||
		(dotImports["syscall"] && function == "ForkExec")
}

func isSuiteGuardMain(importPath, function string) bool {
	return function == "Main" && strings.HasSuffix(importPath, "/internal/suiteguard")
}

func hasSuiteGuardDotImport(dotImports map[string]bool) bool {
	for importPath := range dotImports {
		if strings.HasSuffix(importPath, "/internal/suiteguard") {
			return true
		}
	}
	return false
}

func findRepositoryRoot(start string) (string, error) {
	for directory := start; ; directory = filepath.Dir(directory) {
		_, err := os.Stat(filepath.Join(directory, "go.mod"))
		if err == nil {
			return directory, nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect repository root candidate %s: %w", directory, err)
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("find repository root from %s: go.mod not found", start)
		}
	}
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
