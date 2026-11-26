package verifyselect

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func CLITestNames(root string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(root, "internal/cli/*_test.go"))
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, name := range files {
		data, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(data), "//go:build docscontract") || strings.Contains(string(data), "//go:build repocontract") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), name, data, 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") && fn.Name.Name != "TestMain" {
				seen[fn.Name.Name] = true
			}
		}
	}
	var names []string
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}
func CLIShard(root string, shard, of int) (string, string, error) {
	if of < 1 || shard < 1 || shard > of {
		return "", "", fmt.Errorf("invalid shard %d of %d", shard, of)
	}
	names, err := CLITestNames(root)
	if err != nil {
		return "", "", err
	}
	var chosen []string
	for _, name := range names {
		h := fnv.New32a()
		_, _ = h.Write([]byte(name))
		bucket := int(h.Sum32()%uint32(of)) + 1
		if (shard == of && bucket != of) || (shard != of && bucket == shard) {
			chosen = append(chosen, regexp.QuoteMeta(name))
		}
	}
	flag := "-run"
	if shard == of {
		flag = "-skip"
	}
	pattern := "^$"
	if len(chosen) > 0 {
		pattern = "^(" + strings.Join(chosen, "|") + ")$"
	}
	return flag, pattern, nil
}
func JobArguments(root string, plan Plan, job string) ([]string, error) {
	packages := plan.Packages
	if len(packages) == 1 && packages[0] == "./..." {
		listed, abs, err := listGoPackages(context.Background(), root)
		if err != nil {
			return nil, err
		}
		packages = nil
		for _, p := range listed {
			rel, err := filepath.Rel(abs, p.Dir)
			if err != nil {
				return nil, err
			}
			packages = append(packages, "./"+filepath.ToSlash(rel))
		}
		sort.Strings(packages)
	}
	var result []string
	for _, p := range packages {
		if job == "packages" && p != "./internal/cli" {
			result = append(result, p)
		}
		if job != "packages" && p == "./internal/cli" {
			n := 1
			if job == "cli-2" {
				n = 2
			}
			flag, pattern, err := CLIShard(root, n, 2)
			if err != nil {
				return nil, err
			}
			result = append(result, flag, pattern, p)
		}
	}
	if len(result) > 0 && plan.Fresh {
		result = append([]string{"-count=1"}, result...)
	}
	return result, nil
}
