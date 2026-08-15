// Package suiteguard detects repository writes made by a package's tests.
package suiteguard

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
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
	digest [sha256.Size]byte
}

type ignoreRule struct {
	base      string
	pattern   string
	negated   bool
	directory bool
	anchored  bool
}

func run(runTests func() int, repoRoot string, output io.Writer) int {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: resolve repository root: %v\n", err)
		return 1
	}

	beforeStarted := time.Now()
	rules, err := loadIgnoreRules(root)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: load repository ignore rules: %v\n", err)
		return 1
	}
	before, err := fingerprint(root, rules)
	beforeDuration := time.Since(beforeStarted)
	if err != nil {
		fmt.Fprintf(output, "suiteguard: fingerprint repository before tests: %v\n", err)
		return 1
	}

	testCode := runTests()

	afterStarted := time.Now()
	after, err := fingerprint(root, rules)
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

func fingerprint(root string, rules []ignoreRule) (map[string]fileState, error) {
	result := make(map[string]fileState)
	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if current == root {
			return nil
		}

		relative, err := filepath.Rel(root, current)
		if err != nil {
			return fmt.Errorf("make %s relative to %s: %w", current, root, err)
		}
		relative = filepath.ToSlash(relative)
		if relative == ".git" || ignored(rules, relative, entry.IsDir()) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat %s: %w", relative, err)
		}
		state := fileState{mode: info.Mode()}
		switch {
		case info.Mode().IsRegular():
			state.digest, err = digestFile(current)
		case info.Mode()&fs.ModeSymlink != 0:
			var target string
			target, err = os.Readlink(current)
			state.digest = sha256.Sum256([]byte(target))
		}
		if err != nil {
			return fmt.Errorf("fingerprint %s: %w", relative, err)
		}
		result[relative] = state
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func digestFile(filePath string) ([sha256.Size]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return [sha256.Size]byte{}, err
	}

	hash := sha256.New()
	if _, copyErr := io.Copy(hash, file); copyErr != nil {
		return [sha256.Size]byte{}, errors.Join(copyErr, file.Close())
	}
	if err := file.Close(); err != nil {
		return [sha256.Size]byte{}, err
	}
	var digest [sha256.Size]byte
	copy(digest[:], hash.Sum(nil))
	return digest, nil
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

func loadIgnoreRules(root string) ([]ignoreRule, error) {
	var rules []ignoreRule
	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == ".git" || (relative != "." && ignored(rules, relative, entry.IsDir())) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() {
			return nil
		}

		ignoreFile := filepath.Join(current, ".gitignore")
		file, err := os.Open(ignoreFile)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("open %s: %w", ignoreFile, err)
		}
		base := relative
		if base == "." {
			base = ""
		}
		parsed, parseErr := parseIgnoreRules(file, base)
		closeErr := file.Close()
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", ignoreFile, parseErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close %s: %w", ignoreFile, closeErr)
		}
		rules = append(rules, parsed...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func parseIgnoreRules(reader io.Reader, base string) ([]ignoreRule, error) {
	var rules []ignoreRule
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		line = trimUnescapedTrailingSpaces(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, `\#`) || strings.HasPrefix(line, `\!`) {
			line = line[1:]
		}
		rule := ignoreRule{base: base}
		if strings.HasPrefix(line, "!") {
			rule.negated = true
			line = line[1:]
		}
		if line == "" {
			continue
		}
		rule.directory = strings.HasSuffix(line, "/")
		line = strings.TrimSuffix(line, "/")
		rule.anchored = strings.HasPrefix(line, "/") || strings.Contains(line, "/")
		rule.pattern = strings.TrimPrefix(line, "/")
		if rule.pattern != "" {
			rules = append(rules, rule)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func trimUnescapedTrailingSpaces(value string) string {
	for strings.HasSuffix(value, " ") {
		backslashes := 0
		for index := len(value) - 2; index >= 0 && value[index] == '\\'; index-- {
			backslashes++
		}
		if backslashes%2 == 1 {
			return value[:len(value)-2] + " "
		}
		value = value[:len(value)-1]
	}
	return value
}

func ignored(rules []ignoreRule, relative string, directory bool) bool {
	isIgnored := false
	for _, rule := range rules {
		candidate, applies := relativeToBase(relative, rule.base)
		if !applies || (rule.directory && !directory) {
			continue
		}
		matched := false
		if rule.anchored {
			matched = matchPath(rule.pattern, candidate)
		} else {
			matched = matchSegment(rule.pattern, path.Base(candidate))
		}
		if matched {
			isIgnored = !rule.negated
		}
	}
	return isIgnored
}

func relativeToBase(relative, base string) (string, bool) {
	if base == "" {
		return relative, true
	}
	prefix := base + "/"
	if !strings.HasPrefix(relative, prefix) {
		return "", false
	}
	return strings.TrimPrefix(relative, prefix), true
}

func matchPath(pattern, candidate string) bool {
	patternParts := strings.Split(pattern, "/")
	candidateParts := strings.Split(candidate, "/")
	return matchParts(patternParts, candidateParts)
}

func matchParts(pattern, candidate []string) bool {
	if len(pattern) == 0 {
		return len(candidate) == 0
	}
	if pattern[0] == "**" {
		return matchParts(pattern[1:], candidate) ||
			(len(candidate) > 0 && matchParts(pattern, candidate[1:]))
	}
	return len(candidate) > 0 && matchSegment(pattern[0], candidate[0]) &&
		matchParts(pattern[1:], candidate[1:])
}

func matchSegment(pattern, candidate string) bool {
	matched, err := path.Match(pattern, candidate)
	if err != nil {
		return pattern == candidate
	}
	return matched
}
