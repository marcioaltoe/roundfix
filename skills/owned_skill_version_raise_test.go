package skills

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func assertRecordingRaisesVersion(t *testing.T, entries []ownedSkillVersionEntry, current ownedSkillVersionEntry, wantVersion string) {
	t.Helper()
	record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{"example": entries}}
	before := append([]ownedSkillVersionEntry(nil), entries...)
	updated, err := checkOwnedSkillVersionRecord(record, map[string]ownedSkillVersionEntry{"example": current}, true)
	if err != nil {
		t.Fatal(err)
	}
	want := append(append([]ownedSkillVersionEntry(nil), before...), ownedSkillVersionEntry{Version: wantVersion, Digest: current.Digest})
	if !reflect.DeepEqual(updated.Skills["example"], want) {
		t.Fatalf("recorded history = %v, want %v", updated.Skills["example"], want)
	}
	if !reflect.DeepEqual(record.Skills["example"], before) {
		t.Fatal("recording mutated its input history")
	}
}

func TestRecordingRaisesAVersionRecordedWithOtherContent(t *testing.T) {
	t.Parallel()
	assertRecordingRaisesVersion(t, []ownedSkillVersionEntry{{"0.1.26", "original"}}, ownedSkillVersionEntry{"0.1.26", "changed"}, "0.1.27")
}

func TestRecordingRaisesAVersionBelowTheHighestRecorded(t *testing.T) {
	t.Parallel()
	assertRecordingRaisesVersion(t, []ownedSkillVersionEntry{{"0.1.25", "older"}, {"0.1.26", "original"}}, ownedSkillVersionEntry{"0.1.24", "changed"}, "0.1.27")
}

func TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		declared string
		drift    string
	}{
		{"colliding version", "0.1.26", ""},
		{"stale version", "0.1.24", ""},
		{"different version lines", "0.1.26", "versions"},
		{"drifted body", "0.1.26", "body"},
		{"drifted reference", "0.1.26", "reference"},
		{"extra mirror file", "0.1.26", "extra"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			canonicalRoot, mirrorRoot := filepath.Join(root, "canonical"), filepath.Join(root, "mirror")
			original := "---\nname: example\nmetadata:\n  version: " + test.declared + "\nversion: " + test.declared + "\n---\n# Example\nversion: body-example\n"
			for _, dir := range []string{canonicalRoot, mirrorRoot} {
				if err := os.MkdirAll(filepath.Join(dir, "example", "references"), 0o755); err != nil {
					t.Fatal(err)
				}
				for path, data := range map[string]string{"SKILL.md": original, "references/guide.md": "reference content\n"} {
					if err := os.WriteFile(filepath.Join(dir, "example", path), []byte(data), 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			mirrorPath := filepath.Join(mirrorRoot, "example", "SKILL.md")
			switch test.drift {
			case "versions":
				if err := os.WriteFile(filepath.Join(canonicalRoot, "example", "SKILL.md"), []byte(strings.ReplaceAll(original, test.declared, "0.1.25")), 0o644); err != nil {
					t.Fatal(err)
				}
			case "body":
				if err := os.WriteFile(mirrorPath, []byte(original+"drift\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			case "reference", "extra":
				path := "references/guide.md"
				if test.drift == "extra" {
					path = "extra.md"
				}
				if err := os.WriteFile(filepath.Join(mirrorRoot, "example", path), []byte("drift\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{"example": {{"0.1.26", "original"}}}}
			names := []string{"example"}
			if test.drift != "" && test.drift != "versions" {
				// A preceding skill also needs a raise: refusal must leave it intact.
				for _, dir := range []string{canonicalRoot, mirrorRoot} {
					if err := os.MkdirAll(filepath.Join(dir, "pending"), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(dir, "pending", "SKILL.md"), []byte(original), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				record.Skills["pending"] = []ownedSkillVersionEntry{{"0.1.26", "original"}}
				names = []string{"pending", "example"}
			}
			beforeCanonical := snapshotSkillTree(t, canonicalRoot)
			beforeMirror := snapshotSkillTree(t, mirrorRoot)
			beforeRecord, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			updated, err := recordOwnedSkillVersions(t.Context(), record, names, canonicalRoot, mirrorRoot)
			if test.drift != "" && test.drift != "versions" {
				if err == nil || !strings.Contains(err.Error(), "make skills-sync") {
					t.Fatalf("drift refusal = %v", err)
				}
				if !reflect.DeepEqual(beforeCanonical, snapshotSkillTree(t, canonicalRoot)) || !reflect.DeepEqual(beforeMirror, snapshotSkillTree(t, mirrorRoot)) || !reflect.DeepEqual(updated, record) {
					t.Fatal("drift refusal wrote skill content or history")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := strings.ReplaceAll(original, test.declared, "0.1.27")
			for _, dir := range []string{canonicalRoot, mirrorRoot} {
				got := snapshotSkillTree(t, dir)
				if got["example/SKILL.md"] != want || got["example/references/guide.md"] != "reference content\n" {
					t.Fatalf("rewritten tree = %v", got)
				}
			}
			digest, err := skillFolderHash(t.Context(), os.DirFS(mirrorRoot), "example", "example")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(updated.Skills["example"], []ownedSkillVersionEntry{{"0.1.26", "original"}, {"0.1.27", digest}}) {
				t.Fatalf("recorded history = %v", updated.Skills["example"])
			}
			afterRecord, err := json.Marshal(record)
			if err != nil || !bytes.Equal(beforeRecord, afterRecord) {
				t.Fatalf("input history changed: %v", err)
			}
			repeated, err := recordOwnedSkillVersions(t.Context(), updated, []string{"example"}, canonicalRoot, mirrorRoot)
			if err != nil || !reflect.DeepEqual(repeated, updated) {
				t.Fatalf("repeated recording = %v, %v", repeated, err)
			}
		})
	}
}

func snapshotSkillTree(t *testing.T, root string) map[string]string {
	t.Helper()
	content := make(map[string]string)
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content[filepath.ToSlash(relative)] = string(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return content
}
