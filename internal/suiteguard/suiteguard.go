// Package suiteguard detects repository writes and helper processes left by tests.
package suiteguard

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/suiteguardcontract"
)

// Violation is one path the suite touched inside the repository root.
type Violation struct {
	Path   string
	Change string
}

var sanctionedRegenerationDeclaration struct {
	sync.RWMutex
	command string
}

// DeclareSanctionedRegeneration identifies the sanctioned command running in
// this process. Regeneration tests call it before writing declared outputs.
func DeclareSanctionedRegeneration(command string) {
	sanctionedRegenerationDeclaration.Lock()
	defer sanctionedRegenerationDeclaration.Unlock()
	sanctionedRegenerationDeclaration.command = normalizeCommand(command)
}

// Main fingerprints repoRoot, runs the package's tests, fingerprints again,
// and fails the package naming every path the tests created, modified, or
// removed. Install it from the package's TestMain.
func Main(m *testing.M, repoRoot string) int {
	return run(m.Run, repoRoot, os.Stderr)
}

type fileState struct {
	mode   fs.FileMode
	digest string
}

func run(runTests func() int, repoRoot string, output io.Writer) int {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: resolve repository root: %v\n", err)
		return 1
	}

	beforeStarted := time.Now()
	before, err := fingerprint(root)
	beforeDuration := time.Since(beforeStarted)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: fingerprint repository before tests: %v\n", err)
		return 1
	}

	mark, markErr := markTestBinary()
	if markErr != nil {
		fmt.Fprintf(output, "suiteguard: helper residue check unavailable: %v\n", markErr)
	}
	testCode := runTests()
	if markErr == nil && checkHelperResidue(mark, output) {
		testCode = 1
	}

	afterStarted := time.Now()
	after, err := fingerprint(root)
	afterDuration := time.Since(afterStarted)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: fingerprint repository after tests: %v\n", err)
		return 1
	}

	fmt.Fprintf(
		output,
		"suiteguard: guard cost: measured %d paths; fingerprint took %s before and %s after\n",
		len(before), beforeDuration, afterDuration,
	)
	violations := compare(before, after)
	violations, err = removeSanctionedRegeneration(root, violations)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: read sanctioned-regeneration declaration: %v\n", err)
		return 1
	}
	if len(violations) == 0 {
		return testCode
	}

	fmt.Fprintln(output, "suiteguard: repository boundary violated:")
	for _, violation := range violations {
		fmt.Fprintf(output, "suiteguard: %s: %s\n", violation.Change, violation.Path)
	}
	return 1
}

func removeSanctionedRegeneration(root string, violations []Violation) ([]Violation, error) {
	if len(violations) == 0 {
		return violations, nil
	}
	declarations, err := suiteguardcontract.ReadSanctionedRegenerations(root)
	if err != nil {
		return nil, err
	}
	if len(declarations) == 0 {
		return violations, nil
	}

	command := declaredSanctionedRegenerationCommand()
	if command == "" {
		return violations, nil
	}
	sanctioned := make(map[string]bool)
	for _, declaration := range declarations {
		if normalizeCommand(declaration.Command) != command {
			continue
		}
		for _, output := range declaration.Outputs {
			sanctioned[output] = true
		}
	}
	if len(sanctioned) == 0 {
		return violations, nil
	}

	remaining := make([]Violation, 0, len(violations))
	for _, violation := range violations {
		if !sanctioned[violation.Path] {
			remaining = append(remaining, violation)
		}
	}
	return remaining, nil
}

func declaredSanctionedRegenerationCommand() string {
	sanctionedRegenerationDeclaration.RLock()
	defer sanctionedRegenerationDeclaration.RUnlock()
	return sanctionedRegenerationDeclaration.command
}

func normalizeCommand(command string) string {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) == 0 {
		return ""
	}
	fields[0] = filepath.Base(fields[0])
	return strings.Join(fields, " ")
}

func fingerprint(root string) (map[string]fileState, error) {
	files, err := fingerprintGit(root, "", nil, "ls-files", "-z", "--cached", "--others", "--exclude-per-directory=.gitignore")
	if err != nil {
		return nil, err
	}
	directories, err := fingerprintGit(root, "", nil, "ls-files", "-z", "--others", "--directory", "--exclude-per-directory=.gitignore")
	if err != nil {
		return nil, err
	}
	indexDirectory, err := os.MkdirTemp("", "roundfix-suiteguard-index-")
	if err != nil {
		return nil, fmt.Errorf("create temporary Git index directory: %w", err)
	}
	defer os.RemoveAll(indexDirectory)
	index := filepath.Join(indexDirectory, "index")
	if _, err := fingerprintGit(root, index, files, "update-index", "--add", "--remove", "--info-only", "-z", "--stdin"); err != nil {
		return nil, err
	}
	staged, err := fingerprintGit(root, index, nil, "ls-files", "-z", "--stage")
	if err != nil {
		return nil, err
	}

	result := make(map[string]fileState)
	addParents := func(relative string) {
		for parent := path.Dir(relative); parent != "."; parent = path.Dir(parent) {
			result[parent] = fileState{}
		}
	}
	for _, record := range strings.Split(string(staged), "\x00") {
		if record == "" {
			continue
		}
		metadata, relative, ok := strings.Cut(record, "\t")
		fields := strings.Fields(metadata)
		if !ok || len(fields) != 3 {
			return nil, fmt.Errorf("parse git ls-files --stage record %q", record)
		}
		result[relative] = fileState{digest: fields[0] + " " + fields[1]}
		addParents(relative)
	}
	for _, directory := range strings.Split(string(directories), "\x00") {
		if !strings.HasSuffix(directory, "/") {
			continue
		}
		relative := strings.TrimSuffix(directory, "/")
		result[relative] = fileState{}
		addParents(relative)
	}
	for relative, state := range result {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", relative, err)
		}
		state.mode = info.Mode()
		result[relative] = state
	}
	return result, nil
}

// Git reads repository contents in a subprocess; only the mode pass above
// touches repository paths in this process, outside m.Run's input-log window.
func fingerprintGit(root, index string, input []byte, args ...string) ([]byte, error) {
	options := []string{"-C", root,
		"-c", "core.fileMode=true", "-c", "core.symlinks=true",
		"-c", "core.autocrlf=false", "-c", "core.safecrlf=false", "-c", "core.fsmonitor=false"}
	command := exec.Command("git", append(options, args...)...)
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(entry), "GIT_") {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, "GIT_OPTIONAL_LOCKS=0")
	if index != "" {
		command.Env = append(command.Env, "GIT_INDEX_FILE="+index)
	}
	command.Stdin = bytes.NewReader(input)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return output, nil
}

func compare(before, after map[string]fileState) []Violation {
	violations := make([]Violation, 0)
	for filePath, beforeState := range before {
		afterState, exists := after[filePath]
		switch {
		case !exists:
			violations = append(violations, Violation{Path: filePath, Change: "removed"})
		case beforeState != afterState:
			violations = append(violations, Violation{Path: filePath, Change: "modified"})
		}
	}
	for filePath := range after {
		if _, exists := before[filePath]; !exists {
			violations = append(violations, Violation{Path: filePath, Change: "created"})
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].Path == violations[j].Path {
			return violations[i].Change < violations[j].Change
		}
		return violations[i].Path < violations[j].Path
	})
	return violations
}
