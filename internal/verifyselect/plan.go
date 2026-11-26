package verifyselect

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Mode string

const (
	Settlement Mode = "settlement"
	Complete   Mode = "complete"
	Reach      Mode = "reach"
)
const runBranchPrefix = "roundfix/run-"

type Facts struct {
	RunMarker, Branch string
	Changed, Added    []string
}
type Plan struct {
	Mode     Mode
	Reason   string
	Packages []string
	Fresh    bool
}

func uniquePaths(paths []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, p := range paths {
		p = strings.TrimPrefix(filepath.ToSlash(filepath.Clean(p)), "./")
		if p != "." && p != "" && !seen[p] {
			seen[p] = true
			result = append(result, p)
		}
	}
	sort.Strings(result)
	return result
}

func PlanPackages(facts Facts, module []goListPackage, root string) (Plan, error) {
	prefix := runBranchPrefix + facts.RunMarker
	if facts.RunMarker != "" && (facts.Branch == prefix || strings.HasPrefix(facts.Branch, prefix+"-")) {
		return Plan{Mode: Settlement, Reason: "none at a settlement; the Complete Suite Run runs them once per Spec"}, nil
	}
	for _, p := range facts.Added {
		if path.Dir(p) == "docs/history/specs" && strings.HasSuffix(p, ".md") {
			return Plan{Mode: Complete, Reason: "every package with fresh results for the Complete Suite Run", Packages: []string{"./..."}, Fresh: true}, nil
		}
	}
	dirs := map[string]string{}
	imports := map[string][]string{}
	for _, pkg := range module {
		rel, err := filepath.Rel(root, pkg.Dir)
		if err != nil || rel == ".." || strings.HasPrefix(rel, "../") {
			return Plan{}, fmt.Errorf("package %q is outside repository", pkg.Dir)
		}
		rel = filepath.ToSlash(rel)
		dirs[rel] = pkg.ImportPath
		imports[pkg.ImportPath] = append(append(append([]string{}, pkg.Imports...), pkg.TestImports...), pkg.XTestImports...)
	}
	selected := map[string]bool{}
	all := false
	for _, p := range facts.Changed {
		ignored := false
		for _, tree := range []string{"docs/history", "docs/adr", "docs/specs", "docs/user-guide"} {
			if inDirectory(p, tree) {
				ignored = true
			}
		}
		if ignored {
			continue
		}
		if inDirectory(p, ".agents/skills") {
			for dir, name := range dirs {
				if inDirectory(dir, "internal/baseline") || inDirectory(dir, "internal/baselineacp") || inDirectory(dir, "skills") {
					selected[name] = true
				}
			}
			continue
		}
		owner := path.Dir(p)
		for {
			if name, ok := dirs[owner]; ok {
				selected[name] = true
				break
			}
			if owner == "." {
				all = true
				break
			}
			owner = path.Dir(owner)
		}
	}
	for changed := true; changed; {
		changed = false
		for name, deps := range imports {
			if selected[name] {
				continue
			}
			for _, dep := range deps {
				if selected[dep] {
					selected[name] = true
					changed = true
					break
				}
			}
		}
	}
	result := Plan{Mode: Reach, Reason: "packages reached by the changes"}
	for dir, name := range dirs {
		if all || selected[name] {
			arg := "./" + dir
			if dir == "." {
				arg = "."
			}
			result.Packages = append(result.Packages, arg)
		}
	}
	sort.Strings(result.Packages)
	return result, nil
}

// AddedPaths preserves the distinction between additions and modified archive records.
func AddedPaths(ctx context.Context, root, base string) ([]string, error) {
	merge, err := runGit(ctx, root, "merge-base", "HEAD", base)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, args := range [][]string{{"diff", "--name-only", "--diff-filter=A", "--no-renames", "-z", strings.TrimSpace(string(merge))}, {"ls-files", "--others", "--exclude-standard", "-z"}} {
		out, err := runGit(ctx, root, args...)
		if err != nil {
			return nil, err
		}
		result = append(result, splitNUL(out)...)
	}
	return uniquePaths(result), nil
}
func Select(ctx context.Context, root, base string) (Plan, error) {
	branch, _ := runGit(ctx, root, "symbolic-ref", "--quiet", "--short", "HEAD")
	facts := Facts{RunMarker: os.Getenv("ROUNDFIX_RUN_ID"), Branch: strings.TrimSpace(string(branch))}
	quick, _ := PlanPackages(facts, nil, root)
	if quick.Mode == Settlement {
		return quick, nil
	}
	module, abs, err := listGoPackages(ctx, root)
	if err != nil {
		return Plan{}, err
	}
	facts.Changed, err = ChangedPaths(ctx, root, base)
	if err == nil {
		facts.Added, err = AddedPaths(ctx, root, base)
	}
	if err != nil {
		return Plan{Mode: Reach, Reason: fmt.Sprintf("Git could not list changes (%s); selecting every package", strings.Join(strings.Fields(err.Error()), " ")), Packages: []string{"./..."}}, nil
	}
	return PlanPackages(facts, module, abs)
}
