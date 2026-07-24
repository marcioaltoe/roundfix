package skills

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var lockHashCompatibilityFixturePath = filepath.Join(
	"..",
	"internal",
	"baseline",
	"assets",
	"lock-hash-compatibility-v1.json",
)

type lockHashCompatibilityFixture struct {
	SchemaVersion string `json:"schemaVersion"`
	Version       int    `json:"version"`
	Files         []struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	} `json:"files"`
	ExpectedSHA256 string `json:"expectedSha256"`
}

func TestSkillFolderHashMatchesExternalCompatibilityFixture(t *testing.T) {
	data, err := os.ReadFile(lockHashCompatibilityFixturePath)
	if err != nil {
		t.Fatalf("read Go-owned lock compatibility fixture: %v", err)
	}
	var fixture lockHashCompatibilityFixture
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatalf("decode embedded lock compatibility fixture: %v", err)
	}
	if fixture.SchemaVersion != "setup-context-driven/external-lock-hash-compatibility-v1" || fixture.Version != 1 {
		t.Fatalf("unexpected lock compatibility fixture version: %#v", fixture)
	}

	root := t.TempDir()
	foundNestedPath := false
	for _, file := range fixture.Files {
		if strings.Contains(file.Path, "/") {
			foundNestedPath = true
		}
		path := filepath.Join(root, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create fixture directory for %q: %v", file.Path, err)
		}
		if err := os.WriteFile(path, []byte(file.Content), 0o644); err != nil {
			t.Fatalf("write fixture file %q: %v", file.Path, err)
		}
	}
	if !foundNestedPath {
		t.Fatal("lock compatibility fixture must include a slash-normalized nested path")
	}

	got, err := SkillFolderHash(root)
	if err != nil {
		t.Fatalf("hash compatibility fixture: %v", err)
	}
	if got != fixture.ExpectedSHA256 {
		t.Fatalf("SkillFolderHash() = %q, want pinned fixture digest %q", got, fixture.ExpectedSHA256)
	}
}

func TestSkillFolderHashExcludesMetadataAndDependencyDirectories(t *testing.T) {
	root := t.TempDir()
	writeSkillHashTestFile(t, root, "SKILL.md", "fixture\n")
	want, err := SkillFolderHash(root)
	if err != nil {
		t.Fatalf("hash baseline skill folder: %v", err)
	}

	for _, path := range []string{
		".git/config",
		"node_modules/package/index.js",
		"references/node_modules/package/index.js",
	} {
		writeSkillHashTestFile(t, root, path, "ignored\n")
	}

	got, err := SkillFolderHash(root)
	if err != nil {
		t.Fatalf("hash skill folder with excluded directories: %v", err)
	}
	if got != want {
		t.Fatalf("excluded directories changed hash: got %q, want %q", got, want)
	}
}

func TestSkillFolderHashRejectsUnsafeFilesystemShapes(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T) string
	}{
		{
			name: "root is regular file",
			prepare: func(t *testing.T) string {
				path := filepath.Join(t.TempDir(), "SKILL.md")
				writeSkillHashTestFile(t, filepath.Dir(path), filepath.Base(path), "fixture\n")
				return path
			},
		},
		{
			name: "file symlink",
			prepare: func(t *testing.T) string {
				root := t.TempDir()
				writeSkillHashTestFile(t, root, "SKILL.md", "fixture\n")
				if err := os.Symlink("SKILL.md", filepath.Join(root, "linked.md")); err != nil {
					t.Skipf("create file symlink: %v", err)
				}
				return root
			},
		},
		{
			name: "directory symlink",
			prepare: func(t *testing.T) string {
				root := t.TempDir()
				target := t.TempDir()
				writeSkillHashTestFile(t, target, "outside.md", "outside\n")
				if err := os.Symlink(target, filepath.Join(root, "linked")); err != nil {
					t.Skipf("create directory symlink: %v", err)
				}
				return root
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := SkillFolderHash(test.prepare(t)); err == nil {
				t.Fatal("expected unsafe filesystem shape to be rejected")
			}
		})
	}
}

func TestSkillFolderHashWrapsMissingRootError(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")
	_, err := SkillFolderHash(root)
	if err == nil {
		t.Fatal("expected missing root error")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected wrapped fs.ErrNotExist, got %v", err)
	}
	if !strings.Contains(err.Error(), root) {
		t.Fatalf("expected error to name root %q, got %v", root, err)
	}
}

func writeSkillHashTestFile(t *testing.T, root string, relative string, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create test directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file %q: %v", relative, err)
	}
}

func TestCheckValidatesRoundfixSkillArtifacts(t *testing.T) {
	if diagnostics := Check(); len(diagnostics) > 0 {
		var messages []string
		for _, diagnostic := range diagnostics {
			messages = append(messages, diagnostic.Path+": "+diagnostic.Message)
		}
		t.Fatalf("expected no skill diagnostics, got %s", strings.Join(messages, "\n"))
	}
}

func TestOwnedSkillContractRejectsSetAndVersionDisagreement(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}

	tests := []struct {
		name      string
		mutate    func([]string, []File) ([]string, []File)
		wantMatch string
	}{
		{
			name: "missing owned skill",
			mutate: func(names []string, files []File) ([]string, []File) {
				return names[1:], files
			},
			wantMatch: "missing owned skill",
		},
		{
			name: "unexpected owned skill",
			mutate: func(names []string, files []File) ([]string, []File) {
				return append(names, "unexpected-owned"), append(files, File{
					Skill: "unexpected-owned",
					Path:  "unexpected-owned/SKILL.md",
					Data:  []byte("---\nname: unexpected-owned\nmetadata:\n  version: 0.0.1\n---\n"),
				})
			},
			wantMatch: "unexpected owned skill",
		},
		{
			name: "duplicate owned skill",
			mutate: func(names []string, files []File) ([]string, []File) {
				return append(names, names[0]), files
			},
			wantMatch: "duplicate owned skill",
		},
		{
			name: "owned version disagreement",
			mutate: func(names []string, files []File) ([]string, []File) {
				mutated := append([]File(nil), files...)
				for index := range mutated {
					if mutated[index].Path == "roundfix/SKILL.md" {
						mutated[index].Data = bytes.Replace(
							mutated[index].Data,
							[]byte("version: 0.0.1"),
							[]byte("version: 9.9.9"),
							1,
						)
						break
					}
				}
				return names, mutated
			},
			wantMatch: "metadata.version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			names, mutatedFiles := test.mutate(Names(), files)
			diagnostics := checkOwnedSkillBundle(names, mutatedFiles)
			messages := make([]string, 0, len(diagnostics))
			for _, diagnostic := range diagnostics {
				messages = append(messages, diagnostic.Message)
			}
			if joined := strings.Join(messages, "\n"); !strings.Contains(joined, test.wantMatch) {
				t.Fatalf("expected diagnostic containing %q, got %s", test.wantMatch, joined)
			}
		})
	}
}

func TestCheckOpenAIManifestRequiresEntrypointAndRuntimeCommand(t *testing.T) {
	diagnostics := checkOpenAIManifest("roundfix/agents/openai.yaml", []byte(`
name: roundfix
runtime_hints: {}
`))

	messages := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		messages = append(messages, diagnostic.Message)
	}
	joined := strings.Join(messages, "\n")
	for _, expected := range []string{
		"manifest field entrypoint is required",
		"manifest runtime command is required",
	} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("expected diagnostic %q, got %s", expected, joined)
		}
	}
}

func TestCheckOpenAIManifestAcceptsNestedRuntimeHints(t *testing.T) {
	diagnostics := checkOpenAIManifest("roundfix/agents/openai.yaml", []byte(`
name: roundfix
entrypoint: SKILL.md
runtime:
  hints:
    command: roundfix watch --source coderabbit --pr <number> --until-clean
`))

	if len(diagnostics) > 0 {
		t.Fatalf("expected nested runtime hints to pass, got %#v", diagnostics)
	}
}

func TestFilesIncludeEveryOwnedSkill(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}

	present := make(map[string]bool, len(files))
	for _, file := range files {
		present[file.Path] = true
		if len(file.Data) == 0 {
			t.Fatalf("expected embedded data for %s", file.Path)
		}
	}

	// The operational roundfix skill ships its SKILL.md and OpenAI manifest,
	// and every owned skill ships a SKILL.md.
	required := []string{"roundfix/SKILL.md", "roundfix/agents/openai.yaml"}
	for _, skill := range Names() {
		required = append(required, skill+"/SKILL.md")
	}
	for _, path := range required {
		if !present[path] {
			t.Fatalf("expected embedded file %q, got %d files without it", path, len(files))
		}
	}
}

func TestFrontmatterName(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{"valid", "---\nname: write-prd\ndescription: x\n---\nbody", "write-prd"},
		{"no frontmatter", "# Heading\nbody", ""},
		{"unterminated", "---\nname: write-prd\nbody", ""},
		{"missing name", "---\ndescription: x\n---\nbody", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := frontmatterName(tc.text); got != tc.want {
				t.Fatalf("frontmatterName(%q) = %q, want %q", tc.text, got, tc.want)
			}
		})
	}
}

func TestRecommendedListsExternallyManagedSkills(t *testing.T) {
	recommended := Recommended()
	if len(recommended) == 0 {
		t.Fatal("expected recommended skills from skills-lock.json, got none")
	}
	owned := make(map[string]bool, len(Names()))
	for _, skill := range Names() {
		owned[skill] = true
	}
	for _, skill := range recommended {
		if owned[skill] {
			t.Fatalf("recommended skill %q must not be an owned bundle skill", skill)
		}
	}
}

func TestRecommendedSkillsMatchLock(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "skills-lock.json"))
	if err != nil {
		t.Fatalf("read skills lock: %v", err)
	}
	var lock struct {
		Version int                        `json:"version"`
		Skills  map[string]json.RawMessage `json:"skills"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	if lock.Version != 1 {
		t.Fatalf("skills lock version = %d, want 1", lock.Version)
	}
	want := make([]string, 0, len(lock.Skills))
	for name := range lock.Skills {
		want = append(want, name)
	}
	sort.Strings(want)
	if got := Recommended(); !reflect.DeepEqual(got, want) {
		t.Fatalf("recommended skills = %v, want lock entries %v", got, want)
	}
}

func TestCheckRejectsExecutableSetupEngineArtifacts(t *testing.T) {
	diagnostics := checkThinSetupSkill([]File{
		{
			Skill: "setup-context-driven",
			Path:  "setup-context-driven/SKILL.md",
		},
		{
			Skill: "setup-context-driven",
			Path:  "setup-context-driven/scripts/" + "context_setup.py",
		},
	})
	if len(diagnostics) != 1 {
		t.Fatalf("thin setup diagnostics = %#v", diagnostics)
	}
	if diagnostics[0].Path != "setup-context-driven/scripts/context_setup.py" ||
		!strings.Contains(diagnostics[0].Message, "must not ship runtime") {
		t.Fatalf("thin setup diagnostic = %#v", diagnostics[0])
	}
}

func TestInstallCopiesSkillsToSupportedTargetDirectories(t *testing.T) {
	root := t.TempDir()
	targetDirs := map[string]string{
		"codex":    filepath.Join(root, "codex"),
		"claude":   filepath.Join(root, "claude"),
		"opencode": filepath.Join(root, "opencode"),
	}

	result, err := Install(context.Background(), InstallRequest{
		Target:     "all",
		TargetDirs: targetDirs,
	})
	if err != nil {
		t.Fatalf("install skills: %v", err)
	}
	if len(result.Targets) != 3 {
		t.Fatalf("expected three install targets, got %#v", result.Targets)
	}
	wantFiles := embeddedFileCount(t)
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded files: %v", err)
	}
	for _, target := range result.Targets {
		if target.Files != wantFiles {
			t.Fatalf("expected %d files for %s, got %d", wantFiles, target.Target, target.Files)
		}
		for _, file := range embeddedFiles {
			installed, err := os.ReadFile(filepath.Join(target.Dir, file.Path))
			if err != nil {
				t.Fatalf("read installed file %s for %s: %v", file.Path, target.Target, err)
			}
			if !bytes.Equal(installed, file.Data) {
				t.Fatalf("installed file %s for %s differs from trusted embedded bytes", file.Path, target.Target)
			}
		}
	}
}

func embeddedFileCount(t *testing.T) int {
	t.Helper()
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}
	return len(files)
}

func TestInstallCopiesSkillsToProjectDirectoryByDefault(t *testing.T) {
	projectDir := t.TempDir()

	result, err := Install(context.Background(), InstallRequest{ProjectDir: projectDir})
	if err != nil {
		t.Fatalf("install project skill: %v", err)
	}
	if len(result.Targets) != 1 {
		t.Fatalf("expected one install target, got %#v", result.Targets)
	}
	target := result.Targets[0]
	if target.Target != "project" {
		t.Fatalf("expected project target, got %q", target.Target)
	}
	expectedDir := filepath.Join(projectDir, ".agents", "skills")
	if target.Dir != expectedDir {
		t.Fatalf("expected project target dir %q, got %q", expectedDir, target.Dir)
	}
	if want := embeddedFileCount(t); target.Files != want {
		t.Fatalf("expected %d installed files, got %d", want, target.Files)
	}
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
		"write-prd/SKILL.md",
	} {
		if _, err := os.Stat(filepath.Join(expectedDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
}

func TestInstallRejectsUnsupportedTarget(t *testing.T) {
	_, err := Install(context.Background(), InstallRequest{Target: "other"})
	if err == nil {
		t.Fatal("expected unsupported target error")
	}
	if !strings.Contains(err.Error(), "unsupported skill install target") {
		t.Fatalf("expected unsupported target error, got %v", err)
	}
}
