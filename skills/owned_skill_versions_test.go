package skills

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

var recordSkillVersions = flag.Bool("record-skill-versions", false, "record new, higher owned skill versions without replacing recorded digests")

const ownedSkillVersionsSchema = "roundfix/owned-skill-versions/v1"
const ownedSkillVersionsPath = "testdata/owned-skill-versions.json"
const recordSkillVersionsCommand = "go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions"

type ownedSkillVersionEntry struct {
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

type ownedSkillVersionRecord struct {
	SchemaVersion string                              `json:"schemaVersion"`
	Skills        map[string][]ownedSkillVersionEntry `json:"skills"`
}

// Return a new record only after every entry has been validated. A refusal
// cannot mutate the caller's history, even when another skill needs recording.
func checkOwnedSkillVersionRecord(record ownedSkillVersionRecord, current map[string]ownedSkillVersionEntry, recording bool) (ownedSkillVersionRecord, error) {
	if record.SchemaVersion != ownedSkillVersionsSchema {
		return record, fmt.Errorf("unsupported owned skill version schema %q", record.SchemaVersion)
	}
	for _, name := range slices.Sorted(maps.Keys(record.Skills)) {
		if _, ok := current[name]; !ok {
			return record, fmt.Errorf("recorded skill %q is no longer shipped by the bundle", name)
		}
		var previous [3]uint64
		for index, entry := range record.Skills[name] {
			version, ok := parseSkillVersion(entry.Version)
			if !ok {
				return record, fmt.Errorf("%s: invalid recorded version %q", name, entry.Version)
			}
			if index > 0 && slices.Compare(previous[:], version[:]) >= 0 {
				return record, fmt.Errorf("%s: recorded versions are not in ascending version order", name)
			}
			previous = version
		}
	}
	updated := ownedSkillVersionRecord{SchemaVersion: record.SchemaVersion, Skills: make(map[string][]ownedSkillVersionEntry, len(current))}
	for name, entries := range record.Skills {
		updated.Skills[name] = slices.Clone(entries)
	}
	for _, name := range slices.Sorted(maps.Keys(current)) {
		entry := current[name]
		version, ok := parseSkillVersion(entry.Version)
		if !ok {
			return record, fmt.Errorf("%s: invalid embedded version %q", name, entry.Version)
		}
		entries := record.Skills[name]
		found := false
		for _, recorded := range entries {
			parsed, _ := parseSkillVersion(recorded.Version)
			if parsed != version {
				continue
			}
			if recorded.Digest != entry.Digest {
				if recording {
					break
				}
				return record, fmt.Errorf("%s: content changed under version %s; raise the version", name, entry.Version)
			}
			found = true
			break
		}
		if found {
			continue
		}
		if !recording {
			return record, fmt.Errorf("%s: version %s is not recorded; run %s", name, entry.Version, recordSkillVersionsCommand)
		}
		if len(entries) > 0 {
			latest, _ := parseSkillVersion(entries[len(entries)-1].Version)
			if slices.Compare(version[:], latest[:]) <= 0 {
				if latest[2] == ^uint64(0) {
					return record, fmt.Errorf("%s: cannot raise version %s: patch overflows", name, entries[len(entries)-1].Version)
				}
				entry.Version = fmt.Sprintf("%d.%d.%d", latest[0], latest[1], latest[2]+1)
			}
		}
		updated.Skills[name] = append(updated.Skills[name], entry)
	}
	return updated, nil
}

var ownedSkillVersionLine = regexp.MustCompile(`(?m)^ *version: [^\n]*`)

// Only frontmatter is rewritten; matching examples in the body stay intact.
func mapOwnedSkillVersionLines(data []byte, rewrite func([]byte) []byte) []byte {
	lines := bytes.SplitAfter(data, []byte("\n"))
	if len(lines) == 0 || strings.TrimSpace(string(lines[0])) != "---" {
		return data
	}
	for index := 1; index < len(lines); index++ {
		if strings.TrimSpace(string(lines[index])) == "---" {
			break
		}
		line := lines[index]
		if ownedSkillVersionLine.Match(line) {
			lines[index] = rewrite(line)
		}
	}
	return bytes.Join(lines, nil)
}

func rewriteOwnedSkillVersion(data []byte, version string) []byte {
	return mapOwnedSkillVersionLines(data, func(line []byte) []byte {
		prefix := line[:bytes.Index(line, []byte("version: "))+len("version: ")]
		ending := ""
		if bytes.HasSuffix(line, []byte("\r\n")) {
			ending = "\r\n"
		} else if bytes.HasSuffix(line, []byte("\n")) {
			ending = "\n"
		}
		return []byte(string(prefix) + version + ending)
	})
}

func ownedSkillFolderContent(root string) (map[string]string, error) {
	content := make(map[string]string)
	tree := os.DirFS(root)
	err := fs.WalkDir(tree, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(tree, path)
		if err != nil {
			return err
		}
		if path == "SKILL.md" {
			data = mapOwnedSkillVersionLines(data, func([]byte) []byte {
				return []byte("version: <version>\n")
			})
		}
		content[path] = string(data)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read skill folder %s: %w", root, err)
	}
	return content, nil
}

// Preflight all mirrors and history before any version-line write. Read disk
// rather than the embedded bundle, which predates any version raise.
func recordOwnedSkillVersions(ctx context.Context, record ownedSkillVersionRecord, names []string, canonicalRoot, mirrorRoot string) (ownedSkillVersionRecord, error) {
	current := make(map[string]ownedSkillVersionEntry, len(names))
	for _, name := range names {
		canonical, err := ownedSkillFolderContent(filepath.Join(canonicalRoot, name))
		if err != nil {
			return record, err
		}
		mirror, err := ownedSkillFolderContent(filepath.Join(mirrorRoot, name))
		if err != nil {
			return record, err
		}
		if !reflect.DeepEqual(canonical, mirror) {
			return record, fmt.Errorf("%s: mirror differs from canonical skill outside version lines; run make skills-sync", name)
		}
		data, err := os.ReadFile(filepath.Join(mirrorRoot, name, "SKILL.md"))
		if err != nil {
			return record, fmt.Errorf("read %s skill: %w", name, err)
		}
		metadata, ok := parseSkillFrontmatter(string(data))
		if !ok || !ValidVersion(metadata.Version) {
			return record, fmt.Errorf("%s has no valid mirror version", name)
		}
		digest, err := skillFolderHash(ctx, os.DirFS(mirrorRoot), name, name)
		if err != nil {
			return record, err
		}
		current[name] = ownedSkillVersionEntry{Version: strings.TrimSpace(metadata.Version), Digest: digest}
	}
	updated, err := checkOwnedSkillVersionRecord(record, current, true)
	if err != nil {
		return record, err
	}
	for _, name := range names {
		entries := updated.Skills[name]
		if len(entries) == len(record.Skills[name]) {
			continue
		}
		last := len(entries) - 1
		if entries[last].Version == current[name].Version {
			continue
		}
		for _, root := range []string{canonicalRoot, mirrorRoot} {
			path := filepath.Join(root, name, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				return record, fmt.Errorf("read skill to raise %s: %w", path, err)
			}
			if err := os.WriteFile(path, rewriteOwnedSkillVersion(data, entries[last].Version), 0o644); err != nil {
				return record, fmt.Errorf("raise skill version %s: %w", path, err)
			}
		}
		digest, err := skillFolderHash(ctx, os.DirFS(mirrorRoot), name, name)
		if err != nil {
			return record, err
		}
		updated.Skills[name][last].Digest = digest
	}
	return updated, nil
}

func TestEveryOwnedSkillVersionIsRecorded(t *testing.T) {
	current := make(map[string]ownedSkillVersionEntry, len(Names()))
	for _, name := range Names() {
		data, err := embedded.ReadFile(name + "/SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		metadata, ok := parseSkillFrontmatter(string(data))
		if !ok || !ValidVersion(metadata.Version) {
			t.Fatalf("%s has no valid embedded version", name)
		}
		digest, err := skillFolderHash(t.Context(), embedded, name, name)
		if err != nil {
			t.Fatal(err)
		}
		current[name] = ownedSkillVersionEntry{Version: strings.TrimSpace(metadata.Version), Digest: digest}
	}
	data, err := os.ReadFile(ownedSkillVersionsPath)
	if err != nil && !(os.IsNotExist(err) && *recordSkillVersions) {
		t.Fatalf("read owned skill versions: %v; run %s", err, recordSkillVersionsCommand)
	}
	record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatalf("decode owned skill versions: %v", err)
		}
	}
	var updated ownedSkillVersionRecord
	if *recordSkillVersions {
		updated, err = recordOwnedSkillVersions(t.Context(), record, Names(), "../.agents/skills", ".")
	} else {
		updated, err = checkOwnedSkillVersionRecord(record, current, false)
	}
	if err != nil {
		t.Fatal(err)
	}
	if !*recordSkillVersions || reflect.DeepEqual(record, updated) {
		return
	}
	encoded, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ownedSkillVersionsPath, append(encoded, '\n'), 0o644); err != nil {
		t.Fatalf("write owned skill versions: %v", err)
	}
	t.Logf("recorded owned skill versions in %s", ownedSkillVersionsPath)
}

func TestAChangedOwnedSkillUnderARecordedVersionIsRefused(t *testing.T) {
	t.Parallel()
	record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{
		"example": {{Version: "1.0.0", Digest: "original"}},
	}}
	current := map[string]ownedSkillVersionEntry{"example": {Version: "1.0.0", Digest: "changed"}}
	_, err := checkOwnedSkillVersionRecord(record, current, false)
	if err == nil || !strings.Contains(err.Error(), "content changed under version 1.0.0; raise the version") {
		t.Fatalf("changed content refusal = %v", err)
	}
}

func TestAnUnrecordedOwnedSkillVersionIsRefused(t *testing.T) {
	t.Parallel()
	record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{
		"example": {{Version: "1.0.0", Digest: "original"}},
	}}
	current := map[string]ownedSkillVersionEntry{"example": {Version: "1.0.1", Digest: "new"}}
	_, err := checkOwnedSkillVersionRecord(record, current, false)
	if err == nil || !strings.Contains(err.Error(), "version 1.0.1 is not recorded") || !strings.Contains(err.Error(), recordSkillVersionsCommand) {
		t.Fatalf("unrecorded version refusal = %v", err)
	}
}

func TestRecordingNeverReplacesARecordedVersion(t *testing.T) {
	t.Parallel()
	t.Run("changed digest", func(t *testing.T) {
		assertRecordingRaisesVersion(t, []ownedSkillVersionEntry{{"1.0.0", "original"}}, ownedSkillVersionEntry{"1.0.0", "changed"}, "1.0.1")
	})
	t.Run("lower unrecorded version", func(t *testing.T) {
		assertRecordingRaisesVersion(t, []ownedSkillVersionEntry{{"1.0.1", "original"}}, ownedSkillVersionEntry{"1.0.0", "new"}, "1.0.2")
	})
	for _, test := range []struct {
		name    string
		entries []ownedSkillVersionEntry
		current ownedSkillVersionEntry
		wantErr string
	}{
		{"descending versions", []ownedSkillVersionEntry{{"1.0.1", "new"}, {"1.0.0", "original"}}, ownedSkillVersionEntry{"1.0.2", "next"}, "ascending version order"},
		{"duplicate version", []ownedSkillVersionEntry{{"1.0.0", "original"}, {"1.0.0", "original"}}, ownedSkillVersionEntry{"1.0.1", "new"}, "ascending version order"},
		{"unchanged digest", []ownedSkillVersionEntry{{"1.0.0", "original"}}, ownedSkillVersionEntry{"1.0.0", "original"}, ""},
		{"higher version", []ownedSkillVersionEntry{{"1.0.9", "original"}}, ownedSkillVersionEntry{"1.0.10", "new"}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{"example": test.entries}}
			before, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			updated, err := checkOwnedSkillVersionRecord(record, map[string]ownedSkillVersionEntry{"example": test.current}, true)
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("recording refusal = %v, want %q", err, test.wantErr)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				want := slices.Clone(test.entries)
				if test.current.Version != want[len(want)-1].Version {
					want = append(want, test.current)
				}
				if !reflect.DeepEqual(updated.Skills["example"], want) {
					t.Fatalf("recorded history = %v, want %v", updated.Skills["example"], want)
				}
			}
			after, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatal("recording mutated its input history")
			}
		})
	}
	t.Run("no longer shipped", func(t *testing.T) {
		record := ownedSkillVersionRecord{SchemaVersion: ownedSkillVersionsSchema, Skills: map[string][]ownedSkillVersionEntry{"removed": {{Version: "1.0.0", Digest: "original"}}}}
		for _, recording := range []bool{false, true} {
			_, err := checkOwnedSkillVersionRecord(record, map[string]ownedSkillVersionEntry{}, recording)
			if err == nil || !strings.Contains(err.Error(), "no longer shipped") {
				t.Fatalf("removed skill refusal (recording=%v) = %v", recording, err)
			}
		}
	})
}
