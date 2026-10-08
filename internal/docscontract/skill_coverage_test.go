//go:build docscontract

//verify:always

package docscontract

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/cli"
	"roundfix/internal/skillcoverage"
)

var recordSkillCoverage = flag.Bool("record-skill-coverage", false, "record the current Behavior Surface fingerprints")

func TestTheSkillCoverageMapIsCurrent(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	computed := computeBehaviorSurfaces(t, root)
	if *recordSkillCoverage {
		record := skillcoverage.Record{SchemaVersion: skillcoverage.RecordSchema, Surfaces: computed}
		if err := os.WriteFile(filepath.Join(root, skillcoverage.RecordPath), skillcoverage.EncodeRecord(record), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	record, err := skillcoverage.ParseRecord([]byte(readContractDocument(t, filepath.Join(root, skillcoverage.RecordPath))))
	if err != nil {
		t.Fatal(err)
	}
	coverage, err := skillcoverage.ParseMap([]byte(readContractDocument(t, filepath.Join(root, skillcoverage.MapPath))))
	if err != nil {
		t.Fatal(err)
	}
	index := skillCoverageReferenceIndex(t, root)
	for _, problem := range skillCoverageProblems(computed, record, coverage, index) {
		t.Error(problem)
	}
	var owned struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal([]byte(readContractDocument(t, filepath.Join(root, "skills/testdata/owned-skill-versions.json"))), &owned); err != nil {
		t.Fatal(err)
	}
	var files []string
	if err := filepath.WalkDir(root, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == ".gocache" || entry.Name() == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		relative, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(relative))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, problem := range skillCoveragePathProblems(coverage, owned.Skills, files) {
		t.Error(problem)
	}
}

// cliHelp uses the nearest successful ancestor for a command that rejects --help.
func skillCoverageCLIHelp(t *testing.T, words []string) []byte {
	t.Helper()
	for {
		var stdout, stderr bytes.Buffer
		args := append(append([]string(nil), words...), "--help")
		if cli.Run(args, &stdout, &stderr) == 0 {
			return stdout.Bytes()
		}
		if len(words) == 0 {
			t.Fatalf("root help failed: %s", stderr.String())
		}
		words = words[:len(words)-1]
	}
}

func computeBehaviorSurfaces(t *testing.T, root string) map[string]string {
	t.Helper()
	surfaces := make(map[string]string)
	for _, command := range commandPaths(string(skillCoverageCLIHelp(t, nil))) {
		surfaces["command "+command] = skillcoverage.Fingerprint(skillCoverageCLIHelp(t, strings.Fields(command)))
	}
	tags := make(map[string]bool)
	entries, err := os.ReadDir(filepath.Join(root, "internal/config"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(root, "internal/config", entry.Name()), nil, 0)
		if err != nil {
			t.Fatalf("config keys: %v", err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			structure, ok := node.(*ast.StructType)
			if !ok {
				return true
			}
			for _, field := range structure.Fields.List {
				if field.Tag == nil {
					continue
				}
				literal, err := strconv.Unquote(field.Tag.Value)
				if err != nil {
					t.Fatalf("config keys: %v", err)
				}
				name := strings.Split(reflect.StructTag(literal).Get("yaml"), ",")[0]
				if name != "" && name != "-" {
					tags[name] = true
				}
			}
			return true
		})
	}
	names := make([]string, 0, len(tags))
	for name := range tags {
		names = append(names, name)
	}
	sort.Strings(names)
	surfaces["config keys"] = skillcoverage.Fingerprint([]byte(strings.Join(names, "\n") + "\n"))
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filepath.Join(root, "internal/cli/cli.go"), nil, 0)
	if err != nil {
		t.Fatalf("exit codes: %v", err)
	}
	// Evaluate constants so iota and expressions contribute their values, not syntax.
	var constants []ast.Decl
	for _, declaration := range file.Decls {
		group, ok := declaration.(*ast.GenDecl)
		if !ok || group.Tok != token.CONST {
			continue
		}
		hasExitCode := false
		for _, spec := range group.Specs {
			for _, name := range spec.(*ast.ValueSpec).Names {
				hasExitCode = hasExitCode || strings.HasPrefix(name.Name, "exit")
			}
		}
		if hasExitCode {
			constants = append(constants, declaration)
		}
	}
	info := &types.Info{Defs: make(map[*ast.Ident]types.Object)}
	var checker types.Config
	if _, err := checker.Check("cli", fileSet, []*ast.File{{Name: ast.NewIdent("cli"), Decls: constants}}, info); err != nil {
		t.Fatalf("exit codes: %v", err)
	}
	var codes []string
	for name, object := range info.Defs {
		if value, ok := object.(*types.Const); ok && strings.HasPrefix(name.Name, "exit") {
			codes = append(codes, name.Name+" "+value.Val().ExactString())
		}
	}
	sort.Strings(codes)
	surfaces["exit codes"] = skillcoverage.Fingerprint([]byte(strings.Join(codes, "\n") + "\n"))
	if err := filepath.WalkDir(filepath.Join(root, "docs/user-guide"), func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(file, ".md") {
			return nil
		}
		relative, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("guide %s: %w", relative, err)
		}
		surfaces["guide "+filepath.ToSlash(relative)] = skillcoverage.Fingerprint(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return surfaces
}

func skillCoverageReferenceIndex(t *testing.T, root string) map[string]string {
	t.Helper()
	text := readContractDocument(t, filepath.Join(root, ".agents/skills/roundfix/SKILL.md"))
	const begin = "<!-- roundfix:reference-index:begin -->"
	const end = "<!-- roundfix:reference-index:end -->"
	_, text, ok := strings.Cut(text, begin)
	if !ok {
		t.Fatal("Roundfix reference index is missing")
	}
	text, _, ok = strings.Cut(text, end)
	if !ok {
		t.Fatal("Roundfix reference index is unterminated")
	}
	links := regexp.MustCompile(`\]\((references/[^)]+)\)`)
	commands := regexp.MustCompile("`([^`]+)`")
	index := make(map[string]string)
	for _, line := range strings.Split(text, "\n") {
		columns := strings.Split(line, "|")
		if len(columns) < 4 {
			continue
		}
		link := links.FindStringSubmatch(columns[1])
		if len(link) == 0 || strings.HasPrefix(strings.TrimSpace(columns[2]), "—") {
			continue
		}
		reference := ".agents/skills/roundfix/" + link[1]
		for _, match := range commands.FindAllStringSubmatch(columns[2], -1) {
			family := strings.Fields(match[1])[0]
			if previous, exists := index[family]; exists && previous != reference {
				t.Fatalf("command %s: references %s and %s overlap", family, previous, reference)
			}
			index[family] = reference
		}
	}
	return index
}

func skillCoverageProblems(computed map[string]string, record skillcoverage.Record, coverage skillcoverage.Map, index map[string]string) []string {
	var problems []string
	for id, fingerprint := range computed {
		recorded, exists := record.Surfaces[id]
		if !exists {
			problems = append(problems, id+": missing from Behavior Surface Record")
		} else if recorded != fingerprint {
			problems = append(problems, id+": stale fingerprint")
		}
	}
	for id := range record.Surfaces {
		if _, exists := computed[id]; !exists {
			problems = append(problems, id+": extra Behavior Surface Record entry")
		}
	}
	entries := make(map[string]bool)
	for _, surface := range coverage.Surfaces {
		entries[surface.ID] = true
		if _, exists := computed[surface.ID]; !exists {
			problems = append(problems, surface.ID+": no computed Behavior Surface")
		}
		if guide, ok := strings.CutPrefix(surface.ID, "guide "); ok && !slices.Contains(surface.Sources, guide) {
			problems = append(problems, surface.ID+": sources omit own guide path "+guide)
		}
		if command, ok := strings.CutPrefix(surface.ID, "command "); ok {
			reference := index[strings.Fields(command)[0]]
			if reference == "" {
				problems = append(problems, surface.ID+": no Roundfix indexed reference")
			} else if !slices.Contains(surface.Skills, reference) {
				problems = append(problems, surface.ID+": missing indexed reference "+reference)
			}
			if slices.Contains(surface.Sources, "internal/cli/cli.go") {
				problems = append(problems, surface.ID+": sources must not include internal/cli/cli.go")
			}
		}
	}
	for id := range computed {
		if !entries[id] {
			problems = append(problems, id+": missing Skill Coverage Map entry")
		}
	}
	sort.Strings(problems)
	return problems
}

func skillCoveragePathProblems(coverage skillcoverage.Map, owned map[string]json.RawMessage, files []string) []string {
	existing := make(map[string]bool, len(files))
	for _, file := range files {
		existing[file] = true
	}
	var problems []string
	for _, surface := range coverage.Surfaces {
		for _, file := range surface.Skills {
			parts := strings.Split(file, "/")
			valid := len(parts) >= 4 && parts[0] == ".agents" && parts[1] == "skills"
			if valid {
				_, valid = owned[parts[2]]
			}
			if !valid {
				problems = append(problems, surface.ID+": covering file is outside an owned skill: "+file)
			}
			if !existing[file] {
				problems = append(problems, surface.ID+": covering file does not exist: "+file)
			}
		}
		for _, source := range surface.Sources {
			matches := false
			for _, file := range files {
				if strings.HasSuffix(source, "/") {
					matches = strings.HasPrefix(file, source)
				} else {
					matches, _ = path.Match(source, file)
				}
				if matches {
					break
				}
			}
			if !matches {
				problems = append(problems, surface.ID+": source matches no existing file: "+source)
			}
		}
	}
	sort.Strings(problems)
	return problems
}

func TestAStaleBehaviorSurfaceIsReported(t *testing.T) {
	id := "guide docs/user-guide/usage.md"
	fingerprint := skillcoverage.Fingerprint([]byte("current"))
	computed := map[string]string{id: fingerprint}
	coverage := skillcoverage.Map{Surfaces: []skillcoverage.Surface{{ID: id, Uncovered: "fixture", Sources: []string{"docs/user-guide/usage.md"}}}}
	for _, test := range []struct {
		name     string
		surfaces map[string]string
		want     string
	}{
		{"stale", map[string]string{id: skillcoverage.Fingerprint([]byte("old"))}, id + ": stale fingerprint"},
		{"missing", map[string]string{}, id + ": missing from Behavior Surface Record"},
		{"extra", map[string]string{id: fingerprint, "exit codes": fingerprint}, "exit codes: extra Behavior Surface Record entry"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := skillCoverageProblems(computed, skillcoverage.Record{Surfaces: test.surfaces}, coverage, nil)
			if !reflect.DeepEqual(got, []string{test.want}) {
				t.Fatalf("problems = %v, want %q", got, test.want)
			}
		})
	}
}

func TestAnUnmappedBehaviorSurfaceIsReported(t *testing.T) {
	id := "exit codes"
	fingerprint := skillcoverage.Fingerprint(nil)
	for _, test := range []struct {
		name     string
		computed map[string]string
		coverage skillcoverage.Map
		want     string
	}{
		{"missing", map[string]string{id: fingerprint}, skillcoverage.Map{}, id + ": missing Skill Coverage Map entry"},
		{"extra", map[string]string{}, skillcoverage.Map{Surfaces: []skillcoverage.Surface{{ID: id, Uncovered: "fixture"}}}, id + ": no computed Behavior Surface"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := skillCoverageProblems(test.computed, skillcoverage.Record{Surfaces: test.computed}, test.coverage, nil)
			if !reflect.DeepEqual(got, []string{test.want}) {
				t.Fatalf("problems = %v, want %q", got, test.want)
			}
		})
	}
}

func TestCommandCoverageFollowsTheRoundfixReferenceIndex(t *testing.T) {
	id := "command implement"
	reference := ".agents/skills/roundfix/references/implement.md"
	computed := map[string]string{id: skillcoverage.Fingerprint(nil)}
	coverage := skillcoverage.Map{Surfaces: []skillcoverage.Surface{{ID: id, Skills: []string{".agents/skills/roundfix/references/spec.md"}}}}
	index := map[string]string{"implement": reference}
	got := skillCoverageProblems(computed, skillcoverage.Record{Surfaces: computed}, coverage, index)
	want := []string{id + ": missing indexed reference " + reference}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("problems = %v, want %v", got, want)
	}
	coverage.Surfaces[0].Skills = []string{reference}
	if got := skillCoverageProblems(computed, skillcoverage.Record{Surfaces: computed}, coverage, index); len(got) != 0 {
		t.Fatalf("indexed reference rejected: %v", got)
	}
}

func TestSkillCoveragePathsRequireOwnedFilesAndMatchingSources(t *testing.T) {
	coverage := skillcoverage.Map{Surfaces: []skillcoverage.Surface{{
		ID:      "config keys",
		Skills:  []string{"skills/roundfix/SKILL.md", ".agents/skills/unowned/SKILL.md", ".agents/skills/roundfix/missing.md"},
		Sources: []string{"internal/config/*.go", "docs/user-guide/", "absent.go"},
	}}}
	owned := map[string]json.RawMessage{"roundfix": nil}
	files := []string{"skills/roundfix/SKILL.md", ".agents/skills/unowned/SKILL.md", "internal/config/config.go", "docs/user-guide/usage.md"}
	got := skillCoveragePathProblems(coverage, owned, files)
	want := []string{
		"config keys: covering file does not exist: .agents/skills/roundfix/missing.md",
		"config keys: covering file is outside an owned skill: .agents/skills/unowned/SKILL.md",
		"config keys: covering file is outside an owned skill: skills/roundfix/SKILL.md",
		"config keys: source matches no existing file: absent.go",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("problems = %v, want %v", got, want)
	}
}

func TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"internal/config/config.go":       "package config\ntype Config struct { B string `yaml:\"beta,omitempty\"`; A string `yaml:\"alpha\"`; Duplicate string `yaml:\"alpha\"`; Ignored string `yaml:\"-\"`; Untagged string }\n",
		"internal/config/config_test.go":  "this test source must not be parsed",
		"internal/config/nested/other.go": "this nested source must not be parsed",
		"internal/cli/cli.go":             "package cli\nconst (exitOK = iota; exitFailed; other = 9)\nconst exitPreflight = exitFailed + 1\n",
		"docs/user-guide/nested/guide.md": "# Guide\r\nExact bytes.\n",
	}
	for file, content := range files {
		target := filepath.Join(root, file)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got := computeBehaviorSurfaces(t, root)
	for id, content := range map[string]string{
		"config keys":                           "alpha\nbeta\n",
		"exit codes":                            "exitFailed 1\nexitOK 0\nexitPreflight 2\n",
		"guide docs/user-guide/nested/guide.md": files["docs/user-guide/nested/guide.md"],
	} {
		if want := skillcoverage.Fingerprint([]byte(content)); got[id] != want {
			t.Errorf("%s: fingerprint = %s, want %s", id, got[id], want)
		}
	}
	// This command rejects --help; its bytes must come from its successful parent.
	if got["command qa-report accept"] != skillcoverage.Fingerprint(skillCoverageCLIHelp(t, []string{"qa-report"})) {
		t.Error("command qa-report accept: fingerprint does not use nearest successful ancestor")
	}
}
