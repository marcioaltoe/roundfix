package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/suiteguard"
)

var recordModuleVersions = flag.Bool("record-module-versions", false, "record module content under a version above every recorded version when needed")

const moduleVersionsSchema = "roundfix/baseline-module-versions/v1"
const moduleVersionsPath = "module-versions.json"
const recordModuleVersionsCommand = "go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1"
const moduleVersionLinePattern = `^  "version": [0-9]+,$`

var moduleVersionLine = regexp.MustCompile(`(?m)` + moduleVersionLinePattern)

type moduleVersionEntry struct {
	Version int64  `json:"version"`
	Digest  string `json:"digest"`
}

type moduleVersionRecord struct {
	SchemaVersion string                          `json:"schemaVersion"`
	Modules       map[string][]moduleVersionEntry `json:"modules"`
}

func moduleContent(data []byte) (moduleVersionEntry, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var object map[string]any
	if err := decoder.Decode(&object); err != nil {
		return moduleVersionEntry{}, fmt.Errorf("decode module: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return moduleVersionEntry{}, fmt.Errorf("module must contain exactly one JSON object")
	}
	number, ok := object["version"].(json.Number)
	if !ok {
		return moduleVersionEntry{}, fmt.Errorf("top-level version must be an integer of at least 1")
	}
	version, err := number.Int64()
	if err != nil || version < 1 {
		return moduleVersionEntry{}, fmt.Errorf("top-level version must be an integer of at least 1")
	}
	delete(object, "version")
	content, err := json.Marshal(object)
	if err != nil {
		return moduleVersionEntry{}, fmt.Errorf("encode module content: %w", err)
	}
	return moduleVersionEntry{Version: version, Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(content))}, nil
}

// Locate the decoded top-level field, rather than a similarly indented nested
// field. Call only after moduleContent has validated the JSON object.
func moduleVersionLineRange(data []byte, version int64) ([]int, error) {
	matches := moduleVersionLine.FindAllIndex(data, -1)
	if len(matches) != 1 {
		return nil, fmt.Errorf("top-level version must occupy exactly one line matching %s", moduleVersionLinePattern)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("read module object: %w", err)
	}
	count := 0
	matchesTopLevel := false
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("read module key: %w", err)
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, fmt.Errorf("read module value: %w", err)
		}
		if key == "version" {
			count++
			// InputOffset ends after the value; the matched line ends after its comma.
			matchesTopLevel = int(decoder.InputOffset()) == matches[0][1]-1 && string(value) == strconv.FormatInt(version, 10)
		}
	}
	if count != 1 || !matchesTopLevel {
		return nil, fmt.Errorf("top-level version must occupy exactly one line matching %s", moduleVersionLinePattern)
	}
	return matches[0], nil
}

func moduleVersionRefusal(format string, args ...any) error {
	return fmt.Errorf("%s; run %s", fmt.Sprintf(format, args...), recordModuleVersionsCommand)
}

// Clone history only after validation; neither a refusal nor an append can
// replace the caller's recorded entries through a shared map or backing array.
func checkModuleVersionRecord(record moduleVersionRecord, current map[string]moduleVersionEntry, recording bool) (moduleVersionRecord, error) {
	if record.SchemaVersion != moduleVersionsSchema {
		return record, moduleVersionRefusal("unsupported module version schema %q", record.SchemaVersion)
	}
	for _, name := range slices.Sorted(maps.Keys(record.Modules)) {
		previous := int64(0)
		for _, entry := range record.Modules[name] {
			if entry.Version <= previous {
				return record, moduleVersionRefusal("%s: recorded versions are not strictly ascending", name)
			}
			previous = entry.Version
		}
		if _, ok := current[name]; !ok && !recording {
			return record, moduleVersionRefusal("%s: recorded module is no longer in the catalog", name)
		}
	}
	updated := moduleVersionRecord{SchemaVersion: record.SchemaVersion, Modules: make(map[string][]moduleVersionEntry, len(current))}
	for _, name := range slices.Sorted(maps.Keys(current)) {
		entry := current[name]
		if entry.Version < 1 {
			return record, moduleVersionRefusal("%s: version must be at least 1", name)
		}
		entries := record.Modules[name]
		updated.Modules[name] = slices.Clone(entries)
		found := false
		for _, recorded := range entries {
			if recorded.Version != entry.Version {
				continue
			}
			if recorded.Digest == entry.Digest {
				found = true
				break
			}
			if !recording {
				return record, moduleVersionRefusal("%s: content changed under version %d", name, entry.Version)
			}
		}
		if found {
			continue
		}
		if !recording {
			return record, moduleVersionRefusal("%s: version %d is not recorded", name, entry.Version)
		}
		if len(entries) > 0 && entry.Version <= entries[len(entries)-1].Version {
			highest := entries[len(entries)-1].Version
			if highest == math.MaxInt64 {
				return record, moduleVersionRefusal("%s: cannot raise version %d: integer overflows", name, highest)
			}
			entry.Version = highest + 1
		}
		updated.Modules[name] = append(updated.Modules[name], entry)
	}
	return updated, nil
}

func readModuleVersionFiles(modulesDir string) (map[string]moduleVersionEntry, map[string][]byte, error) {
	paths, err := filepath.Glob(filepath.Join(modulesDir, "*.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("list modules: %w", err)
	}
	if len(paths) == 0 {
		return nil, nil, fmt.Errorf("no module files in %s", modulesDir)
	}
	current := make(map[string]moduleVersionEntry, len(paths))
	contents := make(map[string][]byte, len(paths))
	for _, path := range paths {
		name := strings.TrimSuffix(filepath.Base(path), ".json")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("read module %s: %w", name, err)
		}
		entry, err := moduleContent(data)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", name, err)
		}
		if _, err := moduleVersionLineRange(data, entry.Version); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", name, err)
		}
		current[name], contents[name] = entry, data
	}
	return current, contents, nil
}

// Preflight every module and the entire history before writing any version.
func recordModuleVersionFiles(record moduleVersionRecord, modulesDir string) (moduleVersionRecord, error) {
	current, contents, err := readModuleVersionFiles(modulesDir)
	if err != nil {
		return record, err
	}
	updated, err := checkModuleVersionRecord(record, current, true)
	if err != nil {
		return record, err
	}
	for _, name := range slices.Sorted(maps.Keys(current)) {
		entries := updated.Modules[name]
		version := current[name].Version
		if len(entries) == len(record.Modules[name]) || entries[len(entries)-1].Version == version {
			continue
		}
		data := contents[name]
		line, err := moduleVersionLineRange(data, version)
		if err != nil {
			return record, fmt.Errorf("%s: %w", name, err)
		}
		replacement := []byte(fmt.Sprintf("  \"version\": %d,", entries[len(entries)-1].Version))
		rewritten := append(bytes.Clone(data[:line[0]]), replacement...)
		rewritten = append(rewritten, data[line[1]:]...)
		if err := os.WriteFile(filepath.Join(modulesDir, name+".json"), rewritten, 0o644); err != nil {
			return record, fmt.Errorf("write module %s: %w", name, err)
		}
	}
	return updated, nil
}

func updateModuleVersionRecordFile(recordPath, modulesDir string, recording bool) error {
	data, err := os.ReadFile(recordPath)
	if err != nil && !(os.IsNotExist(err) && recording) {
		return fmt.Errorf("read module versions: %w; run %s", err, recordModuleVersionsCommand)
	}
	record := moduleVersionRecord{SchemaVersion: moduleVersionsSchema}
	if err == nil {
		if err := json.Unmarshal(data, &record); err != nil {
			return fmt.Errorf("decode module versions: %w; run %s", err, recordModuleVersionsCommand)
		}
	}
	var updated moduleVersionRecord
	if recording {
		updated, err = recordModuleVersionFiles(record, modulesDir)
	} else {
		var current map[string]moduleVersionEntry
		current, _, err = readModuleVersionFiles(modulesDir)
		if err == nil {
			updated, err = checkModuleVersionRecord(record, current, false)
		}
	}
	if err != nil {
		return err
	}
	if !recording || reflect.DeepEqual(record, updated) {
		return nil
	}
	encoded, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return fmt.Errorf("encode module versions: %w", err)
	}
	if err := os.WriteFile(recordPath, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("write module versions: %w", err)
	}
	return nil
}

func TestEveryBaselineModuleVersionIsRecorded(t *testing.T) {
	// Sequential: can rewrite shared module version records and assets when the record flag is enabled.
	if *recordModuleVersions {
		suiteguard.DeclareSanctionedRegeneration(recordModuleVersionsCommand)
	}
	if err := updateModuleVersionRecordFile(moduleVersionsPath, "assets/modules", *recordModuleVersions); err != nil {
		t.Fatal(err)
	}
}

func moduleVersionFixture(t *testing.T, version int64) (string, string, []byte, moduleVersionRecord) {
	t.Helper()
	root := t.TempDir()
	modulesDir := filepath.Join(root, "modules")
	if err := os.Mkdir(modulesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data := []byte(fmt.Sprintf("{\n  \"id\": \"example\",\n  \"version\": %d,\n  \"rules\": [{\"version\": 3, \"text\": \"original\"}],\n  \"guides\": [{\"version\": 4}]\n}\n", version))
	writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), data)
	entry, err := moduleContent(data)
	if err != nil {
		t.Fatal(err)
	}
	record := moduleVersionRecord{SchemaVersion: moduleVersionsSchema, Modules: map[string][]moduleVersionEntry{"example": {entry}}}
	recordPath := filepath.Join(root, moduleVersionsPath)
	writeModuleVersionFixtureRecord(t, recordPath, record)
	return modulesDir, recordPath, data, record
}

func writeModuleVersionFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeModuleVersionFixtureRecord(t *testing.T, path string, record moduleVersionRecord) {
	t.Helper()
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeModuleVersionFixture(t, path, append(data, '\n'))
}

func readModuleVersionFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func requireModuleVersionRefusal(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), message) || !strings.HasSuffix(err.Error(), "run "+recordModuleVersionsCommand) {
		t.Fatalf("refusal = %v, want %q ending with record command", err, message)
	}
}

func TestModuleContentDigestIgnoresVersionAndLayout(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, data, _ := moduleVersionFixture(t, 1)
	original, err := moduleContent(data)
	if err != nil {
		t.Fatal(err)
	}
	variants := []string{
		`{"guides":[{"version":4}],"rules":[{"text":"original","version":3}],"version":92,"id":"example"}`,
		string(bytes.Replace(data, []byte(`"version": 1,`), []byte(`"version": 9,`), 1)),
	}
	for _, variant := range variants {
		got, err := moduleContent([]byte(variant))
		if err != nil {
			t.Fatal(err)
		}
		if got.Digest != original.Digest {
			t.Fatalf("digest changed with version/layout: %s != %s", got.Digest, original.Digest)
		}
	}
	changed := bytes.Replace(data, []byte(`"version": 3`), []byte(`"version": 5`), 1)
	got, err := moduleContent(changed)
	if err != nil {
		t.Fatal(err)
	}
	if got.Digest == original.Digest {
		t.Fatal("nested content change did not change digest")
	}
	// Large numbers retain precision under UseNumber.
	a, err := moduleContent([]byte(`{"version":1,"number":9007199254740992}`))
	if err != nil {
		t.Fatal(err)
	}
	b, err := moduleContent([]byte(`{"version":1,"number":9007199254740993}`))
	if err != nil {
		t.Fatal(err)
	}
	if a.Digest == b.Digest {
		t.Fatal("large integer content lost precision")
	}
	if err := updateModuleVersionRecordFile(recordPath, modulesDir, false); err != nil {
		t.Fatal(err)
	}
}

func TestAModuleChangedUnderARecordedVersionIsRefused(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, data, _ := moduleVersionFixture(t, 1)
	writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), bytes.Replace(data, []byte("original"), []byte("changed"), 1))
	requireModuleVersionRefusal(t, updateModuleVersionRecordFile(recordPath, modulesDir, false), "example: content changed under version 1")
}

func TestAnUnrecordedModuleVersionIsRefused(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, data, _ := moduleVersionFixture(t, 1)
	writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), bytes.Replace(data, []byte(`"version": 1,`), []byte(`"version": 2,`), 1))
	requireModuleVersionRefusal(t, updateModuleVersionRecordFile(recordPath, modulesDir, false), "example: version 2 is not recorded")
}

func TestRecordingRaisesAModuleAboveTheHighestRecordedVersion(t *testing.T) {
	t.Parallel()
	for _, version := range []int64{1, 2} {
		t.Run(fmt.Sprintf("current-%d", version), func(t *testing.T) {
			modulesDir, recordPath, data, record := moduleVersionFixture(t, version)
			// Version 1 is recorded with different content; version 2 is unrecorded.
			record.Modules["example"] = []moduleVersionEntry{{Version: 1, Digest: "old"}, {Version: 7, Digest: "latest"}}
			writeModuleVersionFixtureRecord(t, recordPath, record)
			if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
				t.Fatal(err)
			}
			expected := bytes.Replace(data, []byte(fmt.Sprintf(`"version": %d,`, version)), []byte(`"version": 8,`), 1)
			if got := readModuleVersionFixture(t, filepath.Join(modulesDir, "example.json")); !bytes.Equal(got, expected) {
				t.Fatalf("rewritten module = %s, want %s", got, expected)
			}
			if err := updateModuleVersionRecordFile(recordPath, modulesDir, false); err != nil {
				t.Fatal(err)
			}
			before := readModuleVersionFixture(t, recordPath)
			if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, readModuleVersionFixture(t, recordPath)) || !bytes.Equal(expected, readModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"))) {
				t.Fatal("second recording changed files")
			}
		})
	}
}

func TestRecordingKeepsAHigherModuleVersion(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, data, record := moduleVersionFixture(t, 9)
	record.Modules["example"] = []moduleVersionEntry{{Version: 7, Digest: "old"}}
	writeModuleVersionFixtureRecord(t, recordPath, record)
	if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, readModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"))) {
		t.Fatal("higher version was rewritten")
	}
	var updated moduleVersionRecord
	if err := json.Unmarshal(readModuleVersionFixture(t, recordPath), &updated); err != nil {
		t.Fatal(err)
	}
	if entries := updated.Modules["example"]; len(entries) != 2 || entries[1].Version != 9 {
		t.Fatalf("recorded history = %#v", entries)
	}
	if err := updateModuleVersionRecordFile(recordPath, modulesDir, false); err != nil {
		t.Fatal(err)
	}
}

func TestRecordingNeverRewritesARecordedModuleEntry(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, data, record := moduleVersionFixture(t, 1)
	original := readModuleVersionFixture(t, recordPath)
	if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, readModuleVersionFixture(t, recordPath)) {
		t.Fatal("recorded content was rewritten")
	}
	writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), bytes.Replace(data, []byte("original"), []byte("changed"), 1))
	updated, err := recordModuleVersionFiles(record, modulesDir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updated.Modules["example"][:1], record.Modules["example"]) {
		t.Fatal("existing recorded entry changed")
	}
	if len(record.Modules["example"]) != 1 {
		t.Fatal("caller's record mutated")
	}
	if updated.Modules["example"][1].Version != 2 {
		t.Fatal("new content did not append a higher version")
	}
}

func TestAModuleVersionNotOnItsOwnLineIsRefused(t *testing.T) {
	t.Parallel()
	for name, invalid := range map[string]string{
		"inline":       `{"id":"example", "version":1,"rules":[]}`,
		"wrong-indent": "{\n \"version\": 1,\n \"id\": \"example\"\n}",
		"nested-decoy": "{\n  \"version\":1,\n  \"rules\": [{\n  \"version\": 1,\n  \"text\": \"x\"}]\n}",
		"duplicate":    "{\n  \"version\": 1,\n  \"version\": 1,\n  \"id\": \"example\"\n}",
	} {
		t.Run(name, func(t *testing.T) {
			modulesDir, recordPath, _, _ := moduleVersionFixture(t, 1)
			writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), []byte(invalid))
			for _, recording := range []bool{false, true} {
				err := updateModuleVersionRecordFile(recordPath, modulesDir, recording)
				if err == nil || !strings.Contains(err.Error(), "example") || !strings.Contains(err.Error(), moduleVersionLinePattern) {
					t.Fatalf("layout refusal = %v", err)
				}
			}
		})
	}
}

func TestRecordingRewritesOnlyTheTopLevelVersionLine(t *testing.T) {
	t.Parallel()
	modulesDir, recordPath, _, _ := moduleVersionFixture(t, 1)
	data := []byte("{\n  \"id\": \"example\",\n  \"version\": 1,\n  \"rules\": [{\n    \"version\": 3,\n    \"text\": \"changed\"\n  }],\n  \"guides\": [{\n    \"version\": 4,\n    \"text\": \"guide\"\n  }]\n}\n")
	writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), data)
	if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
		t.Fatal(err)
	}
	expected := bytes.Replace(data, []byte("  \"version\": 1,"), []byte("  \"version\": 2,"), 1)
	if got := readModuleVersionFixture(t, filepath.Join(modulesDir, "example.json")); !bytes.Equal(got, expected) {
		t.Fatalf("bytes outside top-level version changed: %s", got)
	}
}

func TestRecordingWritesNothingWhenAModuleIsRefused(t *testing.T) {
	t.Parallel()
	for _, refusal := range []string{"layout", "history", "schema", "overflow", "invalid-version", "malformed-json"} {
		t.Run(refusal, func(t *testing.T) {
			modulesDir, recordPath, data, record := moduleVersionFixture(t, 1)
			// example would be raised before z-refused if recording wrote incrementally.
			writeModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"), bytes.Replace(data, []byte("original"), []byte("changed"), 1))
			other := bytes.Replace(data, []byte("example"), []byte("z-refused"), 1)
			switch refusal {
			case "layout":
				other = bytes.Replace(other, []byte("  \"version\": 1,"), []byte(" \"version\": 1,"), 1)
			case "history":
				record.Modules["z-refused"] = []moduleVersionEntry{{Version: 3}, {Version: 2}}
			case "schema":
				record.SchemaVersion = "unsupported"
			case "overflow":
				record.Modules["z-refused"] = []moduleVersionEntry{{Version: math.MaxInt64, Digest: "old"}}
			case "invalid-version":
				other = bytes.Replace(other, []byte("  \"version\": 1,"), []byte("  \"version\": 0,"), 1)
			case "malformed-json":
				other = []byte("{")
			}
			writeModuleVersionFixture(t, filepath.Join(modulesDir, "z-refused.json"), other)
			writeModuleVersionFixtureRecord(t, recordPath, record)
			beforeRecord := readModuleVersionFixture(t, recordPath)
			beforeModule := readModuleVersionFixture(t, filepath.Join(modulesDir, "example.json"))
			if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err == nil {
				t.Fatal("expected refusal")
			}
			for path, expected := range map[string][]byte{recordPath: beforeRecord, filepath.Join(modulesDir, "example.json"): beforeModule, filepath.Join(modulesDir, "z-refused.json"): other} {
				if !bytes.Equal(expected, readModuleVersionFixture(t, path)) {
					t.Fatalf("refusal changed %s", path)
				}
			}
		})
	}
}

func TestModuleVersionRecordHistoryRefusals(t *testing.T) {
	t.Parallel()
	for _, refusal := range []string{"missing", "schema", "ascending", "removed"} {
		t.Run(refusal, func(t *testing.T) {
			modulesDir, recordPath, _, record := moduleVersionFixture(t, 1)
			message := ""
			switch refusal {
			case "missing":
				if err := os.Remove(recordPath); err != nil {
					t.Fatal(err)
				}
				message = "read module versions"
			case "schema":
				record.SchemaVersion = "other"
				message = "unsupported module version schema"
			case "ascending":
				record.Modules["example"] = []moduleVersionEntry{{Version: 1}, {Version: 1}}
				message = "example: recorded versions are not strictly ascending"
			case "removed":
				record.Modules["removed"] = []moduleVersionEntry{{Version: 1}}
				message = "removed: recorded module is no longer in the catalog"
			}
			if refusal != "missing" {
				writeModuleVersionFixtureRecord(t, recordPath, record)
			}
			requireModuleVersionRefusal(t, updateModuleVersionRecordFile(recordPath, modulesDir, false), message)
			if refusal == "removed" || refusal == "missing" {
				if err := updateModuleVersionRecordFile(recordPath, modulesDir, true); err != nil {
					t.Fatal(err)
				}
				if err := updateModuleVersionRecordFile(recordPath, modulesDir, false); err != nil {
					t.Fatal(err)
				}
				if bytes.Contains(readModuleVersionFixture(t, recordPath), []byte(`"removed"`)) {
					t.Fatal("removed module remains in record")
				}
			}
		})
	}
}
