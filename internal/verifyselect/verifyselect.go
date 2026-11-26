// Package verifyselect maps repository changes to the packages and contracts they can affect.
package verifyselect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// ChangedPaths lists committed changes since the merge base with baseRef,
// staged changes, unstaged changes, and untracked files. Returned paths are
// repository-relative, unique, and sorted.
func ChangedPaths(ctx context.Context, repoRoot, baseRef string) ([]string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	baseRef = strings.TrimSpace(baseRef)
	if repoRoot == "" {
		return nil, errors.New("list changed paths: repository root is required")
	}
	if baseRef == "" {
		return nil, errors.New("list changed paths: base ref is required")
	}

	mergeBaseOutput, err := runGit(ctx, repoRoot, "merge-base", "HEAD", baseRef)
	if err != nil {
		return nil, fmt.Errorf("list changed paths: resolve merge base: %w", err)
	}
	mergeBase := strings.TrimSpace(string(mergeBaseOutput))
	if mergeBase == "" {
		return nil, errors.New("list changed paths: resolve merge base: git returned an empty revision")
	}

	commands := [][]string{
		{"diff", "--name-only", "-z", "--no-renames", mergeBase, "HEAD"},
		{"diff", "--cached", "--name-only", "-z", "--no-renames"},
		{"diff", "--name-only", "-z", "--no-renames"},
		{"ls-files", "--others", "--exclude-standard", "-z"},
	}
	paths := make([]string, 0)
	for _, args := range commands {
		output, runErr := runGit(ctx, repoRoot, args...)
		if runErr != nil {
			return nil, fmt.Errorf("list changed paths: %w", runErr)
		}
		paths = append(paths, splitNUL(output)...)
	}
	return uniquePaths(paths), nil
}

// Run executes the verify-select command and returns its process exit code.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("verify-select", flag.ContinueOnError)
	flags.SetOutput(stderr)
	baseRef := flags.String("base", "main", "base ref used to select changed tests")
	repoRoot := flags.String("repo", ".", "repository root")
	planFlag := flags.Bool("plan", false, "print the package plan")
	job := flags.String("job", "", "print arguments for packages, cli-1 or cli-2")
	tag := flags.String("tag", "", "restrict contracts to docscontract or repocontract")
	contracts := flags.Bool("contracts", false, "print relevant Repository Contract Test invocations")
	allContracts := flags.Bool("all", false, "print every Repository Contract Test invocation (requires -contracts)")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "verify-select: positional arguments are not supported")
		return 2
	}
	explicitBase, explicitPlan, explicitJob := false, false, false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "base":
			explicitBase = true
		case "plan":
			explicitPlan = true
		case "job":
			explicitJob = true
		}
	})
	if *contracts && (explicitPlan || explicitJob) {
		fmt.Fprintln(stderr, "verify-select: incompatible flags")
		return 2
	}
	if *allContracts && !*contracts {
		fmt.Fprintln(stderr, "verify-select: -all requires -contracts")
		return 2
	}
	if *allContracts && explicitBase {
		fmt.Fprintln(stderr, "verify-select: -all cannot be combined with -base")
		return 2
	}
	if (*tag != "" && *tag != "docscontract" && *tag != "repocontract") || (*tag != "" && !*contracts) || (*allContracts && !*contracts) || (*contracts && (*planFlag || *job != "")) {
		fmt.Fprintln(stderr, "verify-select: incompatible flags")
		return 2
	}
	filter := func(items []ContractTest) []ContractTest {
		if *tag == "" {
			return items
		}
		var result []ContractTest
		for _, item := range items {
			if item.Tag == *tag {
				result = append(result, item)
			}
		}
		return result
	}
	if *contracts {
		if *allContracts {
			all, err := DiscoverContracts(*repoRoot)
			if err != nil {
				fmt.Fprintf(stderr, "verify-select: %v\n", err)
				return 1
			}
			for _, invocation := range ContractInvocations(filter(all)) {
				fmt.Fprintln(stdout, invocation)
			}
			fmt.Fprintf(stderr, "verify-select: contracts: %d selected (every Repository Contract Test), 0 not selected\n", len(all))
			return 0
		}
		all, selected, err := SelectContracts(ctx, *repoRoot, *baseRef)
		if err != nil {
			if all == nil {
				fmt.Fprintf(stderr, "verify-select: %v\n", err)
				return 1
			}
			fmt.Fprintf(stderr, "verify-select: %v; selecting every contract that is not boundary\n", err)
		}
		for _, invocation := range ContractInvocations(filter(selected)) {
			fmt.Fprintln(stdout, invocation)
		}
		always := 0
		type contractKey struct{ name, packagePath, tag string }
		selectedContracts := make(map[contractKey]bool, len(selected))
		for _, contract := range selected {
			selectedContracts[contractKey{contract.Name, contract.Package, contract.Tag}] = true
			if contract.Class == ContractAlways {
				always++
			}
		}
		var relevantNotSelected, boundaries []string
		for _, contract := range all {
			if contract.Class == ContractRelevant && !selectedContracts[contractKey{contract.Name, contract.Package, contract.Tag}] {
				relevantNotSelected = append(relevantNotSelected, fmt.Sprintf("%s (%s)", contract.Name, contract.Package))
			}
			if contract.Class == ContractBoundary {
				boundaries = append(boundaries, fmt.Sprintf("%s (%s)", contract.Name, contract.Package))
			}
		}
		fmt.Fprintf(stderr, "verify-select: contracts: %d selected (%d always, %d by change), %d not selected, %d boundary", len(selected), always, len(selected)-always, len(all)-len(selected), len(boundaries))
		if len(relevantNotSelected) != 0 {
			fmt.Fprintf(stderr, "; relevant not selected: %s", strings.Join(relevantNotSelected, ", "))
		}
		if len(boundaries) != 0 {
			fmt.Fprintf(stderr, "; boundary: %s", strings.Join(boundaries, ", "))
		}
		fmt.Fprintln(stderr)
		return 0
	}

	if *job != "" && *job != "packages" && *job != "cli-1" && *job != "cli-2" {
		fmt.Fprintln(stderr, "verify-select: unknown job")
		return 2
	}
	plan, err := Select(ctx, *repoRoot, *baseRef)
	if err != nil {
		fmt.Fprintf(stderr, "verify-select: %v\n", err)
		return 1
	}
	fmt.Fprintln(stderr, "verify-select: packages: "+plan.Reason)
	arguments := plan.Packages
	if *job != "" {
		arguments, err = JobArguments(*repoRoot, plan, *job)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	for _, argument := range arguments {
		fmt.Fprintln(stdout, argument)
	}
	return 0
}

func inDirectory(path, directory string) bool {
	return path == directory || strings.HasPrefix(path, directory+"/")
}

type goListPackage struct {
	Dir          string
	ImportPath   string
	Imports      []string
	TestImports  []string
	XTestImports []string
}

func decodePackages(output []byte) ([]goListPackage, error) {
	decoder := json.NewDecoder(bytes.NewReader(output))
	packages := make([]goListPackage, 0)
	for {
		var listedPackage goListPackage
		if err := decoder.Decode(&listedPackage); err != nil {
			if errors.Is(err, io.EOF) {
				return packages, nil
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		packages = append(packages, listedPackage)
	}
}

func listGoPackages(ctx context.Context, repoRoot string) ([]goListPackage, string, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, "", fmt.Errorf("resolve repository root: %w", err)
	}
	output, err := runCommand(ctx, root, "go", "list", "-json", "./...")
	if err != nil {
		return nil, "", err
	}
	listed, err := decodePackages(output)
	if err != nil {
		return nil, "", err
	}
	return listed, root, nil
}

func splitNUL(output []byte) []string {
	parts := bytes.Split(output, []byte{0})
	paths := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) != 0 {
			paths = append(paths, string(part))
		}
	}
	return paths
}

func runGit(ctx context.Context, repoRoot string, args ...string) ([]byte, error) {
	gitArgs := []string{"--no-optional-locks", "-C", repoRoot, "-c", "core.fsmonitor=false"}
	gitArgs = append(gitArgs, args...)
	return runCommand(ctx, "", "git", gitArgs...)
}

func runCommand(ctx context.Context, workDir, name string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = workDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("%s %s: %s: %w", name, strings.Join(args, " "), detail, err)
	}
	return stdout.Bytes(), nil
}
