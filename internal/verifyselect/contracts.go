package verifyselect

import (
	"context"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ContractClass describes when a Repository Contract Test is relevant.
type ContractClass string

const (
	ContractAlways   ContractClass = "always"
	ContractPackage  ContractClass = "package"
	ContractRelevant ContractClass = "relevant"
	ContractBoundary ContractClass = "boundary"
)

// ContractTest is a discovered test and its repository-relative relevance.
type ContractTest struct {
	Name, Package, Tag, File string
	Class                    ContractClass
	Paths                    []string
	Reason                   string
}

// DiscoverContracts reads tagged test declarations without invoking the Go toolchain.
func DiscoverContracts(repoRoot string) ([]ContractTest, error) {
	contracts := make([]ContractTest, 0)
	err := filepath.WalkDir(repoRoot, func(file string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("discover contracts: walk %s: %w", file, walkErr)
		}
		if entry.IsDir() {
			if file != repoRoot && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "testdata" || entry.Name() == "vendor" || entry.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		relative, err := filepath.Rel(repoRoot, file)
		if err != nil {
			return fmt.Errorf("discover contracts: relative path: %w", err)
		}
		relative = filepath.ToSlash(relative)
		source, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("%s: read contract source: %w", relative, err)
		}
		// Parse only the header first: an unrelated test's body is not our input.
		header, err := parser.ParseFile(token.NewFileSet(), relative, source, parser.PackageClauseOnly|parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s: parse contract header: %w", relative, err)
		}
		tag, err := contractBuildTag(header, relative)
		if err != nil {
			return err
		}
		if tag == "" {
			return nil
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), relative, source, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s: parse contract tests: %w", relative, err)
		}
		var headerComments []*ast.Comment
		for _, group := range parsed.Comments {
			for _, comment := range group.List {
				if comment.Pos() < parsed.Package {
					headerComments = append(headerComments, comment)
				}
			}
		}
		defaults, err := contractDirective(headerComments, relative)
		if err != nil {
			return err
		}
		for _, declaration := range parsed.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || !isContractTest(function, parsed) {
				continue
			}
			relevance := defaults
			if function.Doc != nil {
				own, err := contractDirective(function.Doc.List, relative+": "+function.Name.Name)
				if err != nil {
					return err
				}
				if own.Class != ContractPackage {
					relevance = own
				}
			}
			relevance.Name, relevance.Package, relevance.Tag, relevance.File = function.Name.Name, path.Dir(relative), tag, relative
			contracts = append(contracts, relevance)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(contracts, func(i, j int) bool {
		if contracts[i].Package != contracts[j].Package {
			return contracts[i].Package < contracts[j].Package
		}
		return contracts[i].Name < contracts[j].Name
	})
	return contracts, nil
}

func contractBuildTag(header *ast.File, file string) (string, error) {
	var modern, legacy constraint.Expr
	for _, group := range header.Comments {
		for _, comment := range group.List {
			if comment.Pos() >= header.Package {
				continue
			}
			line := comment.Text
			if !constraint.IsGoBuild(line) && !constraint.IsPlusBuild(line) {
				continue
			}
			expression, err := constraint.Parse(line)
			if err != nil {
				return "", fmt.Errorf("%s: parse build constraint: %w", file, err)
			}
			if constraint.IsGoBuild(line) {
				if modern != nil {
					return "", fmt.Errorf("%s: more than one go:build constraint", file)
				}
				modern = expression
			} else if legacy == nil {
				legacy = expression
			} else {
				legacy = &constraint.AndExpr{X: legacy, Y: expression}
			}
		}
	}
	expression := modern
	if expression == nil {
		expression = legacy
	}
	if expression == nil {
		return "", nil
	}
	if !namesContractTag(expression) {
		return "", nil
	}
	neither := expression.Eval(func(string) bool { return false })
	docs := expression.Eval(func(tag string) bool { return tag == "docscontract" })
	repo := expression.Eval(func(tag string) bool { return tag == "repocontract" })
	if neither || docs == repo {
		return "", fmt.Errorf("%s: build constraint must select exactly one contract tag", file)
	}
	if docs {
		return "docscontract", nil
	}
	return "repocontract", nil
}

func namesContractTag(expression constraint.Expr) bool {
	switch expression := expression.(type) {
	case *constraint.TagExpr:
		return expression.Tag == "docscontract" || expression.Tag == "repocontract"
	case *constraint.NotExpr:
		return namesContractTag(expression.X)
	case *constraint.AndExpr:
		return namesContractTag(expression.X) || namesContractTag(expression.Y)
	case *constraint.OrExpr:
		return namesContractTag(expression.X) || namesContractTag(expression.Y)
	}
	return false
}

func isContractTest(function *ast.FuncDecl, file *ast.File) bool {
	name := function.Name.Name
	if function.Recv != nil || !strings.HasPrefix(name, "Test") || len(name) == 4 {
		return false
	}
	first, _ := utf8.DecodeRuneInString(name[4:])
	if unicode.IsLower(first) || function.Type.TypeParams != nil || function.Type.Results != nil || function.Type.Params == nil || len(function.Type.Params.List) != 1 {
		return false
	}
	parameter := function.Type.Params.List[0]
	if len(parameter.Names) > 1 {
		return false
	}
	pointer, ok := parameter.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	testingName := ""
	for _, imported := range file.Imports {
		importedPath, err := strconv.Unquote(imported.Path.Value)
		if err == nil && importedPath == "testing" {
			testingName = "testing"
			if imported.Name != nil {
				testingName = imported.Name.Name
			}
		}
	}
	if testingName == "." {
		ident, ok := pointer.X.(*ast.Ident)
		return ok && ident.Name == "T"
	}
	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "T" {
		return false
	}
	qualifier, ok := selector.X.(*ast.Ident)
	return ok && qualifier.Name == testingName
}

func contractDirective(comments []*ast.Comment, scope string) (ContractTest, error) {
	result := ContractTest{Class: ContractPackage}
	var directives []string
	for _, comment := range comments {
		if strings.HasPrefix(comment.Text, "//verify:") {
			directives = append(directives, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//verify:")))
		}
	}
	if len(directives) > 1 {
		return result, fmt.Errorf("%s: more than one Contract Relevance directive", scope)
	}
	if len(directives) == 0 {
		return result, nil
	}
	fields := strings.Fields(directives[0])
	class := ""
	if len(fields) > 0 {
		class = fields[0]
	}
	result.Class = ContractClass(class)
	switch result.Class {
	case ContractAlways:
	case ContractRelevant:
		if len(fields) < 2 {
			return result, fmt.Errorf("%s: relevant needs at least one path", scope)
		}
		for _, pattern := range fields[1:] {
			_, err := path.Match(pattern, "")
			invalid := strings.HasPrefix(pattern, "/") || strings.HasPrefix(pattern, "./") || strings.Contains(pattern, "**") || strings.Contains(pattern, "\\") || err != nil
			for _, segment := range strings.Split(pattern, "/") {
				invalid = invalid || segment == ".."
			}
			if invalid {
				return result, fmt.Errorf("%s: invalid relevant path %q", scope, pattern)
			}
		}
		result.Paths = fields[1:]
	case ContractBoundary:
		if len(fields) < 2 {
			return result, fmt.Errorf("%s: boundary needs a reason", scope)
		}
		result.Reason = strings.TrimSpace(strings.TrimPrefix(directives[0], class))
	default:
		return result, fmt.Errorf("%s: unknown Contract Relevance class %q", scope, class)
	}
	return result, nil
}

// SelectContractPaths selects contracts from a known change list.
func SelectContractPaths(contracts []ContractTest, changed []string) []ContractTest {
	moduleWide := false
	for _, file := range changed {
		moduleWide = moduleWide || file == "go.mod" || file == "go.sum" || file == "Makefile"
	}
	var selected []ContractTest
	for _, contract := range contracts {
		if contract.Class == ContractBoundary {
			continue
		}
		matches := moduleWide || contract.Class == ContractAlways
		for _, file := range changed {
			matches = matches || path.Dir(file) == contract.Package
			if contract.Class == ContractRelevant {
				for _, pattern := range contract.Paths {
					if strings.HasSuffix(pattern, "/") {
						matches = matches || strings.HasPrefix(file, pattern)
					} else {
						match, _ := path.Match(pattern, file)
						matches = matches || match
					}
				}
			}
		}
		if matches {
			selected = append(selected, contract)
		}
	}
	return selected
}

// SelectContracts discovers contracts and uses Git changes, selecting every
// non-boundary contract together with a diagnostic if Git cannot list changes.
func SelectContracts(ctx context.Context, repoRoot, baseRef string) ([]ContractTest, []ContractTest, error) {
	all, err := DiscoverContracts(repoRoot)
	if err != nil {
		return nil, nil, err
	}
	changed, err := ChangedPaths(ctx, repoRoot, baseRef)
	if err != nil {
		var selected []ContractTest
		for _, contract := range all {
			if contract.Class != ContractBoundary {
				selected = append(selected, contract)
			}
		}
		return all, selected, fmt.Errorf("select contracts: %w", err)
	}
	return all, SelectContractPaths(all, changed), nil
}

// ContractInvocations formats sorted tag, test pattern, and package arguments.
func ContractInvocations(selected []ContractTest) []string {
	if len(selected) == 0 {
		return nil
	}
	tags, names, packages := make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{})
	for _, contract := range selected {
		tags[contract.Tag] = struct{}{}
		names[regexp.QuoteMeta(contract.Name)] = struct{}{}
		packages["./"+contract.Package] = struct{}{}
	}
	return []string{strings.Join(sortedContractKeys(tags), ",") + " ^(" + strings.Join(sortedContractKeys(names), "|") + ")$ " + strings.Join(sortedContractKeys(packages), " ")}
}

func sortedContractKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
