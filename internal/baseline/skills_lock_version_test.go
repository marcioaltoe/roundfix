// Suite: Baseline skills lock version compatibility
// Invariant: lock hashes follow the declared version and writers preserve that version.
// Boundary IN: temporary installed trees, ordered local locks and offline immutable sources.
// Boundary OUT: network acquisition, skills CLI execution and Doctor rendering.
package baseline

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/skillhash"
)

func TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion(t *testing.T) {
	t.Parallel()
	for _, version := range []int{1, 2} {
		t.Run(fmt.Sprintf("version %d", version), func(t *testing.T) {
			root, _ := ownedLockFixture(t)
			writeRetiredSkillTestFile(t, root, ".agents/skills/roundfix/SKILL.md", "# hi\r\n")
			writeRetiredSkillTestFile(t, root, ".agents/skills/roundfix/metadata.json", "meta\n")
			files := []skillhash.File{{Path: "SKILL.md", Content: []byte("# hi\r\n")}, {Path: "metadata.json", Content: []byte("meta\n")}}
			want, err := skillhash.SumForLockVersion(version, files)
			if err != nil {
				t.Fatal(err)
			}
			other, err := skillhash.SumForLockVersion(3-version, files)
			if err != nil {
				t.Fatal(err)
			}
			before := fmt.Sprintf(`{ "version": %d.0, "skills":{"external":{"extra":1e2},"roundfix":{"computedHash":%q,"extra":"\u0061"}}, "tail":true }`+"\n", version, other)
			filename := filepath.Join(root, skillsLockPath)
			writeRetiredSkillTestFile(t, root, skillsLockPath, before)
			document, err := loadSkillsLock(filename)
			if err != nil {
				t.Fatal(err)
			}
			cloned := document.clone()
			copied, err := cloned.digest([]restoreFile{{Path: "SKILL.md", Content: []byte("# hi\r\n")}, {Path: "metadata.json", Content: []byte("meta\n")}})
			if err != nil || cloned.version != version || copied != want {
				t.Fatalf("cloned version = %d, digest = %q, %v; want %d and %q", cloned.version, copied, err, version, want)
			}
			stale, err := StaleOwnedLockEntries(t.Context(), root, []string{"roundfix"})
			if err != nil || !reflect.DeepEqual(stale, []StaleOwnedLockEntry{{Skill: "roundfix", Hash: true}}) {
				t.Fatalf("other version's hash = %+v, %v", stale, err)
			}
			update, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
			if err != nil || !reflect.DeepEqual(update.Hashed, []string{"roundfix"}) || len(update.Unpinned) != 0 {
				t.Fatalf("record = %+v, %v", update, err)
			}
			after, err := os.ReadFile(filename)
			expected := strings.Replace(before, other, want, 1)
			if err != nil || string(after) != expected {
				t.Fatalf("record changed more than hash: %s, %v; want %s", after, err, expected)
			}
			stale, err = StaleOwnedLockEntries(t.Context(), root, []string{"roundfix"})
			if err != nil || len(stale) != 0 {
				t.Fatalf("current hash is stale: %+v, %v", stale, err)
			}
			update, err = RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
			if err != nil || len(update.Names()) != 0 {
				t.Fatalf("second record = %+v, %v", update, err)
			}
			if after, err := os.ReadFile(filename); err != nil || string(after) != expected {
				t.Fatalf("second record changed lock: %s, %v", after, err)
			}
		})
	}
}

func TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion(t *testing.T) {
	t.Parallel()
	for _, version := range []int{1, 2} {
		t.Run(fmt.Sprintf("version %d", version), func(t *testing.T) {
			root := t.TempDir()
			content := "# hi\r\n"
			writeRetiredSkillTestFile(t, root, ".agents/skills/tdd/SKILL.md", content)
			writeRetiredSkillTestFile(t, root, ".agents/skills/tdd/metadata.json", "meta\n")
			files := []skillhash.File{{Path: "SKILL.md", Content: []byte(content)}, {Path: "metadata.json", Content: []byte("meta\n")}}
			own, err := skillhash.SumForLockVersion(version, files)
			if err != nil {
				t.Fatal(err)
			}
			other, err := skillhash.SumForLockVersion(3-version, files)
			if err != nil {
				t.Fatal(err)
			}
			for _, test := range []struct{ name, hash, action, reason string }{
				{"other version", other, RetiredSkillKeep, RetiredSkillEdited},
				{"own version", own, RetiredSkillDelete, ""},
			} {
				t.Run(test.name, func(t *testing.T) {
					before := fmt.Sprintf(`{"version":%d,"skills":{"tdd":{"computedHash":%q}}}`, version, test.hash)
					writeRetiredSkillTestFile(t, root, skillsLockPath, before)
					got, err := InstalledRetiredSkills(root)
					if err != nil || len(got) != 1 || got[0].Skill != "tdd" || got[0].Action != test.action || got[0].Reason != test.reason {
						t.Fatalf("retired disposition = %+v, %v", got, err)
					}
					assertRestoreFile(t, root, skillsLockPath, before)
				})
			}
			writeRetiredSkillTestFile(t, root, ".agents/skills/tdd/SKILL.md", "# edited\n")
			got, err := InstalledRetiredSkills(root)
			if err != nil || len(got) != 1 || got[0].Action != RetiredSkillKeep || got[0].Reason != RetiredSkillEdited {
				t.Fatalf("edited copy = %+v, %v", got, err)
			}
		})
	}
}

func TestSkillsRestoreWritesTheHashOfTheLockVersion(t *testing.T) {
	t.Parallel()
	repo, source, dependencies := newSkillsRestoreFixture(t, map[string]map[string]string{
		"agentic-cli-design": {"SKILL.md": "# hi\r\n", "metadata.json": "meta\n"},
	})
	writeInspectionFile(t, repo, skillsLockPath, "{\n  \"version\": 2.0,\n  \"skills\": {\"unrelated\": {\"extra\": true}}\n}\n")
	before := snapshotVisibleTree(t, repo)
	request := SkillsRestoreRequest{Repository: repo, ProfileID: dependencies.profile.ID, Skills: []string{"agentic-cli-design"}, SourceDir: source}
	preview, err := restoreSkills(t.Context(), request, dependencies)
	var restoreErr *SkillsRestoreError
	if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != "plan.confirmation.required" || preview.PlanDigest == nil {
		t.Fatalf("preview = %+v, %v", preview, err)
	}
	assertVisibleTree(t, repo, before)
	request.Confirmation = *preview.PlanDigest
	applied, err := restoreSkills(t.Context(), request, dependencies)
	if err != nil || !applied.OK || !applied.Applied {
		t.Fatalf("apply = %+v, %v", applied, err)
	}
	assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/SKILL.md", "# hi\r\n")
	assertRestoreFile(t, repo, ".agents/skills/agentic-cli-design/metadata.json", "meta\n")
	lock := readSkillsRestoreLock(t, repo)
	entry := lock["skills"].(map[string]any)["agentic-cli-design"].(map[string]any)
	const want = "93c36952d8843b33e68ae970ebabcc03e1687c9bb0dd1ecaf357060318952f55"
	if lock["version"] != float64(2) || entry["computedHash"] != want || entry["ref"] != dependencies.profile.Skills["agentic-cli-design"].Source.Ref {
		t.Fatalf("restored lock = %+v", lock)
	}
	if !reflect.DeepEqual(lock["skills"].(map[string]any)["unrelated"], map[string]any{"extra": true}) {
		t.Fatalf("unrelated entry changed: %+v", lock)
	}
	data, err := os.ReadFile(filepath.Join(repo, skillsLockPath))
	if err != nil || !strings.Contains(string(data), `"version": 2.0`) {
		t.Fatalf("version representation changed: %s, %v", data, err)
	}
	request.Confirmation = ""
	again, err := restoreSkills(t.Context(), request, dependencies)
	if err != nil || !again.OK || again.Applied || len(again.PlannedChanges) != 0 {
		t.Fatalf("second restore = %+v, %v", again, err)
	}
}

func TestEveryLockReaderRefusesAnUnsupportedVersionByName(t *testing.T) {
	t.Parallel()
	for _, text := range []string{"3", "0", "1.5", "-1", "3e0"} {
		t.Run(text, func(t *testing.T) {
			root := t.TempDir()
			before := fmt.Sprintf(`{"version":%s,"skills":{}}`, text)
			writeRetiredSkillTestFile(t, root, skillsLockPath, before)
			readers := []struct {
				name string
				read func() error
			}{
				{"document", func() error { _, err := loadSkillsLock(filepath.Join(root, skillsLockPath)); return err }},
				{"owned recorder", func() error {
					_, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
					return err
				}},
				{"owned comparison", func() error { _, err := StaleOwnedLockEntries(t.Context(), root, []string{"roundfix"}); return err }},
				{"retired comparison", func() error { _, err := InstalledRetiredSkills(root); return err }},
			}
			for _, reader := range readers {
				t.Run(reader.name, func(t *testing.T) {
					err := reader.read()
					want := fmt.Sprintf("skills-lock.json version %s is not supported; this Roundfix reads versions 1 and 2.", text)
					var restoreErr *SkillsRestoreError
					if !errors.Is(err, skillhash.ErrUnsupportedLockVersion) || !errors.As(err, &restoreErr) || restoreErr.Category != SkillsRestoreInvalid || restoreErr.Finding.Code != "lock.unsupported-version" || restoreErr.Finding.Message != want || restoreErr.Finding.Action != "Run roundfix upgrade; a newer Roundfix may read this version." {
						t.Fatalf("refusal = %+v, %v; want %q", restoreErr, err, want)
					}
					assertRestoreFile(t, root, skillsLockPath, before)
				})
			}
		})
	}
	for _, text := range []string{`[]`, `null`, `{"skills":{}}`, `{"version":"2","skills":{}}`, `{"version":true,"skills":{}}`, `{"version":null,"skills":{}}`, `{"version":2}`, `{"version":2,"skills":[]}`, `{"version":2,"skills":null}`} {
		t.Run("invalid "+text, func(t *testing.T) {
			root := t.TempDir()
			writeRetiredSkillTestFile(t, root, skillsLockPath, text)
			_, err := loadSkillsLock(filepath.Join(root, skillsLockPath))
			var restoreErr *SkillsRestoreError
			if !errors.As(err, &restoreErr) || restoreErr.Finding.Code != "lock.invalid" || restoreErr.Finding.Message != "skills-lock.json must hold a version number and a skills object." || errors.Is(err, skillhash.ErrUnsupportedLockVersion) {
				t.Fatalf("malformed refusal = %+v, %v", restoreErr, err)
			}
		})
	}
	t.Run("absent lock starts in version 1", func(t *testing.T) {
		document, err := loadSkillsLock(filepath.Join(t.TempDir(), skillsLockPath))
		if err != nil {
			t.Fatal(err)
		}
		files := []restoreFile{{Path: "SKILL.md", Content: []byte("# hi\r\n")}}
		for _, copy := range []skillsLockDocument{document, document.clone()} {
			got, err := copy.digest(files)
			if err != nil || copy.version != 1 || got != externalSkillsLockDigest(files) {
				t.Fatalf("absent lock digest = %q, %v; version %d", got, err, copy.version)
			}
		}
	})
}
