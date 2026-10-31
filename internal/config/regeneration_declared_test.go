//go:build docscontract

//verify:relevant .roundfixrc.yml internal/baseline/ .agents/skills/ skills/ docs/agents/ docs/references/coverage-record.json internal/spec/ internal/cli/baseline_* cmd/roundfix/ internal/suiteguard/ internal/suiteguardcontract/ docs/references/behavior-surfaces.json docs/references/skill-coverage.json

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRegenerationIsDeclared(t *testing.T) {
	root := configTestRepositoryRoot(t)
	declarations := repositoryDerivedPaths(t, filepath.Join(root, ".roundfixrc.yml"))
	if len(declarations) < 4 {
		t.Fatalf("derived declarations = %d, want the module, profile, formatter, and managed guidance declarations", len(declarations))
	}

	changed := runRegenerationFixture(t, root, -1)
	assertChangedPathsDeclared(t, changed, declarations)

	// Each declaration introduced by this Task must protect an output that the
	// fixture actually regenerates. Removing a declaration does not change what
	// the commands write, so one regeneration answers every ablation: without
	// that declaration, some changed path must be left undeclared.
	for _, index := range []int{1, 2, 3} {
		t.Run(fmt.Sprintf("declaration_%d_is_required", index), func(t *testing.T) {
			if uncovered := undeclaredChangedPaths(changed, removeDeclaration(declarations, index)); len(uncovered) == 0 {
				t.Fatalf("removing declaration %d still covered every changed path: %v", index, changed)
			}
		})
	}

	t.Run("an_unsanctioned_coverage_write_is_refused", func(t *testing.T) {
		fixture := filepath.Join(t.TempDir(), "repository")
		copyRepository(t, root, fixture)
		writeArchivedItemRecord(t, fixture, false)
		editBaselineModule(t, fixture)
		staleCoverageRecord(t, fixture)
		gitFixture(t, fixture, "record the unsanctioned coverage fixture")

		for _, declaration := range declarations {
			if !declaration.Matches("docs/references/coverage-record.json") {
				continue
			}
			command := exec.Command("sh", "-c", declaration.Regenerate)
			command.Dir = fixture
			command.Env = os.Environ()
			output, err := command.CombinedOutput()
			if err == nil {
				t.Fatalf("unsanctioned coverage regeneration succeeded:\n%s", output)
			}
			if !bytes.Contains(output, []byte("suiteguard: modified: docs/references/coverage-record.json")) {
				t.Fatalf("unsanctioned coverage regeneration did not name the refused record: %v\n%s", err, output)
			}
			return
		}
		t.Fatal("coverage regeneration declaration not found")
	})
}

func runRegenerationFixture(t *testing.T, sourceRoot string, removeIndex int) []changedFile {
	t.Helper()
	fixture := filepath.Join(t.TempDir(), "repository")
	copyRepository(t, sourceRoot, fixture)
	if removeIndex >= 0 {
		removeDeclarationFromFile(t, filepath.Join(fixture, ".roundfixrc.yml"), removeIndex)
	}
	writeArchivedItemRecord(t, fixture, true)
	staleCoverageRecord(t, fixture)
	editBaselineModule(t, fixture)
	gitFixture(t, fixture, "record the module edit")

	declarations := repositoryDerivedPaths(t, filepath.Join(fixture, ".roundfixrc.yml"))
	commands := declarations
	if removeIndex >= 0 {
		commands = repositoryDerivedPaths(t, filepath.Join(sourceRoot, ".roundfixrc.yml"))
	}
	for index, declaration := range commands {
		command := exec.Command("sh", "-c", declaration.Regenerate)
		command.Dir = fixture
		command.Env = os.Environ()
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("regeneration declaration %d (%q): %v\n%s", index, declaration.Regenerate, err, output)
		}
	}
	changed := changedFiles(t, fixture)
	coverageChanged := false
	for _, file := range changed {
		if file.path == "docs/references/coverage-record.json" {
			coverageChanged = true
		}
	}
	if !coverageChanged {
		t.Fatal("regeneration did not change the stale coverage record")
	}
	coverage, err := os.ReadFile(filepath.Join(fixture, "docs/references/coverage-record.json"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(coverage, []byte("TestRegenerationFixtureStaleEntry")) {
		t.Fatal("regenerated coverage record still contains the stale entry")
	}
	return changed
}

type changedFile struct {
	path  string
	lines []string
}

func assertChangedPathsDeclared(t *testing.T, changed []changedFile, declarations []DerivedPathDeclaration) {
	t.Helper()
	if uncovered := undeclaredChangedPaths(changed, declarations); len(uncovered) != 0 {
		t.Fatalf("regeneration produced undeclared changes: %v", uncovered)
	}
	if len(changed) == 0 {
		t.Fatal("regeneration produced no changed paths")
	}
}

func undeclaredChangedPaths(changed []changedFile, declarations []DerivedPathDeclaration) []string {
	var uncovered []string
	for _, file := range changed {
		declared := false
		for _, declaration := range declarations {
			if declaration.MatchesLines(file.path) {
				declared = changedLinesMatch(file.lines, declaration.Lines.Match)
			} else if declaration.Matches(file.path) {
				declared = true
			}
			if declared {
				break
			}
		}
		if !declared {
			uncovered = append(uncovered, file.path)
		}
	}
	sort.Strings(uncovered)
	return uncovered
}

func changedLinesMatch(lines []string, expression string) bool {
	pattern := regexp.MustCompile(expression)
	if len(lines) == 0 {
		return false
	}
	for _, line := range lines {
		if !pattern.MatchString(line) {
			return false
		}
	}
	return true
}

func changedFiles(t *testing.T, repository string) []changedFile {
	t.Helper()
	output := gitOutput(t, repository, "diff", "--unified=0", "--no-color", "HEAD")
	byPath := make(map[string][]string)
	var current string
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "+++ b/") {
			current = strings.TrimPrefix(line, "+++ b/")
			continue
		}
		if current == "" || line == "" || strings.HasPrefix(line, "@@") || strings.HasPrefix(line, "diff ") || strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "\\") {
			continue
		}
		if line[0] == '+' || line[0] == '-' {
			byPath[current] = append(byPath[current], line[1:])
		}
	}
	paths := make([]string, 0, len(byPath))
	for path := range byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	files := make([]changedFile, 0, len(paths))
	for _, path := range paths {
		files = append(files, changedFile{path: path, lines: byPath[path]})
	}
	return files
}

func repositoryDerivedPaths(t *testing.T, path string) []DerivedPathDeclaration {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	config, err := ResolveConfigProposal(nil, data)
	if err != nil {
		t.Fatal(err)
	}
	return config.Delivery.DerivedPaths
}

func removeDeclaration(declarations []DerivedPathDeclaration, index int) []DerivedPathDeclaration {
	result := append([]DerivedPathDeclaration(nil), declarations...)
	return append(result[:index:index], result[index+1:]...)
}

func removeDeclarationFromFile(t *testing.T, path string, index int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	root := document.Content[0]
	derived := yamlMappingValue(root, "delivery", "derived_paths")
	if derived == nil || index >= len(derived.Content) {
		t.Fatalf("derived declaration %d is absent", index)
	}
	derived.Content = append(derived.Content[:index], derived.Content[index+1:]...)
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&document); err != nil {
		t.Fatal(err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, output.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func yamlMappingValue(node *yaml.Node, keys ...string) *yaml.Node {
	for _, key := range keys {
		if node == nil || node.Kind != yaml.MappingNode {
			return nil
		}
		var next *yaml.Node
		for index := 0; index+1 < len(node.Content); index += 2 {
			if node.Content[index].Value == key {
				next = node.Content[index+1]
				break
			}
		}
		node = next
	}
	return node
}

func editBaselineModule(t *testing.T, repository string) {
	t.Helper()
	path := filepath.Join(repository, "internal/baseline/assets/modules/context-workflow.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	old := []byte("Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.")
	if !bytes.Contains(data, old) {
		t.Fatal("baseline module guidance clause was not found")
	}
	data = bytes.Replace(data, old, []byte("Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing phrase."), 1)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeArchivedItemRecord(t *testing.T, repository string, sanctionCoverage bool) {
	t.Helper()
	for _, path := range []string{"docs/specs", "docs/workflow/authorizations"} {
		if err := os.RemoveAll(filepath.Join(repository, path)); err != nil {
			t.Fatal(err)
		}
	}
	records, err := filepath.Glob(filepath.Join(repository, "docs/history/specs/*.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		if err := os.Remove(record); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(repository, "docs/history/specs/0001-the-archived-item.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `---
schema: roundfix/archive-record/v1
spec: 0001-the-archived-item
title: The archived item
status: archived
created: "2026-10-30"
archived: "2026-10-30"
source: docs/specs/0001-the-archived-item
source_revision: 0123456789012345678901234567890123456789
disposition: pass
qa_task: task_01
qa_report: qa/2026-10-30.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
promoted: []
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
`
	if sanctionCoverage {
		content += `  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
`
	}
	content += "---\n\n# The archived item\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func staleCoverageRecord(t *testing.T, repository string) {
	t.Helper()
	path := filepath.Join(repository, "docs/references/coverage-record.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]json.RawMessage
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	var packages map[string][]string
	if err := json.Unmarshal(record["packages"], &packages); err != nil {
		t.Fatal(err)
	}
	names := packages["roundfix/internal/spec"]
	if len(names) == 0 {
		t.Fatal("coverage record has no internal/spec tests to replace")
	}
	names[0] = "TestRegenerationFixtureStaleEntry"
	record["packages"], err = json.Marshal(packages)
	if err != nil {
		t.Fatal(err)
	}
	data, err = json.MarshalIndent(record, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func copyRepository(t *testing.T, source, destination string) {
	t.Helper()
	if err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == ".git" || strings.HasPrefix(relative, ".git"+string(filepath.Separator)) || relative == ".gocache" || strings.HasPrefix(relative, ".gocache"+string(filepath.Separator)) || relative == "bin" || strings.HasPrefix(relative, "bin"+string(filepath.Separator)) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, info.Mode().Perm())
	}); err != nil {
		t.Fatal(err)
	}
}

func gitFixture(t *testing.T, repository, message string) {
	t.Helper()
	gitOutput(t, repository, "init", "-q")
	gitOutput(t, repository, "config", "user.email", "test@example.invalid")
	gitOutput(t, repository, "config", "user.name", "Regeneration Test")
	gitOutput(t, repository, "config", "commit.gpgSign", "false")
	gitOutput(t, repository, "add", ".")
	gitOutput(t, repository, "commit", "-q", "-m", message)
}

func gitOutput(t *testing.T, repository string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repository}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func configTestRepositoryRoot(t *testing.T) string {
	t.Helper()
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := workingDirectory
	for {
		if _, err := os.Stat(filepath.Join(root, ".roundfixrc.yml")); err == nil {
			return root
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("repository root not found")
		}
		root = parent
	}
}
