package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/skills"
)

func TestDoctorNamesAStaleOwnedSkillLockEntry(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		stale bool
	}{
		{name: "pinned after a bundle install", stale: true},
		{name: "unpinned with the installed hash"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := doctorOwnedLockFixture(t, test.stale)
			dependencies := doctorOwnedLockDependencies()
			counts := fmt.Sprintf("%d required: %d Roundfix-owned, 0 external", len(skills.Names()), len(skills.Names()))
			wantOK := "skills: ok (" + counts + ")\n"
			want := wantOK
			if test.stale {
				want = "skills: warn (" + counts + "; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: archive-spec, roundfix, write-prd; next: roundfix skills install --target project)\n"
			}
			before := snapshotDoctorPath(t, root)
			var output bytes.Buffer
			printDoctorResult(&output, repositorySkillsCheck(t.Context(), dependencies, root))
			if output.String() != want {
				t.Fatalf("skills line = %q, want %q", output.String(), want)
			}
			if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) {
				t.Fatal("Doctor changed repository files")
			}
			if _, err := baseline.RecordInstalledOwnedSkills(t.Context(), root, skills.Names()); err != nil {
				t.Fatal(err)
			}
			output.Reset()
			printDoctorResult(&output, repositorySkillsCheck(t.Context(), dependencies, root))
			if output.String() != wantOK {
				t.Fatalf("skills line after recording the bundle = %q, want %q", output.String(), wantOK)
			}
		})
	}
}

func TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry(t *testing.T) {
	t.Parallel()
	root := doctorOwnedLockFixture(t, false)
	filename := filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md")
	content := regexp.MustCompile(`(?m)^([ \t]*)version:.*$`).ReplaceAllString(string(mustReadBytes(t, filename)), "${1}version: 0.0.1")
	mustWrite(t, filename, content)
	mustWrite(t, filepath.Join(root, "skills-lock.json"), `{"version":1,"skills":{"roundfix":{"ref":"old","computedHash":"`+strings.Repeat("0", 64)+`"},"testing-boss":{"ref":"keep"}}}`)
	var output bytes.Buffer
	printDoctorResult(&output, repositorySkillsCheck(t.Context(), doctorOwnedLockDependencies(), root))
	line := output.String()
	if !strings.HasPrefix(line, "skills: failed (") || !strings.Contains(line, "; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: roundfix; next: ") || strings.Count(line, "roundfix skills install --target project") != 1 || strings.Contains(line, "testing-boss") {
		t.Fatalf("skills line = %q", line)
	}
}

func doctorOwnedLockDependencies() doctorDependencies {
	return doctorDependencies{
		resolveExternal: func(string) ([]string, bool, error) { return nil, true, nil },
		checkSkills:     skills.CheckRepositoryWithExternal,
	}
}

func doctorOwnedLockFixture(t *testing.T, stale bool) string {
	t.Helper()
	root := t.TempDir()
	if _, err := skills.Install(t.Context(), skills.InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(root, ".agents", "skills", "testing-boss")
	mustMkdir(t, external)
	mustWrite(t, filepath.Join(external, "SKILL.md"), "external skill\n")
	entries := map[string]map[string]any{"testing-boss": {"ref": "keep", "computedHash": strings.Repeat("0", 64)}}
	for _, name := range []string{"archive-spec", "roundfix", "write-prd"} {
		hash, err := skills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", name))
		if err != nil {
			t.Fatal(err)
		}
		entry := map[string]any{"computedHash": hash}
		if stale {
			if name == "write-prd" {
				entry["computedHash"] = strings.Repeat("0", 64)
			} else {
				entry["ref"] = "old"
			}
		}
		entries[name] = entry
	}
	data, err := json.Marshal(map[string]any{"version": 1, "skills": entries})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "skills-lock.json"), string(data))
	return root
}

func TestStaleOwnedLockEntriesReadsOnlyInstalledNamedTrees(t *testing.T) {
	t.Parallel()
	root := doctorOwnedLockFixture(t, true)
	filename := filepath.Join(root, "skills-lock.json")
	before := string(mustReadBytes(t, filename))
	before = strings.Replace(before, `"skills":{`, `"skills":{"absent":{"ref":"keep"},"nonobject":false,`, 1)
	mustWrite(t, filename, before)
	snapshot := snapshotDoctorPath(t, root)
	names := []string{"write-prd", "roundfix", "absent", "nonobject", "unlisted", "archive-spec", "roundfix"}
	got, err := baseline.StaleOwnedLockEntries(t.Context(), root, names)
	want := []baseline.StaleOwnedLockEntry{
		{Skill: "archive-spec", Ref: "old"},
		{Skill: "roundfix", Ref: "old"},
		{Skill: "write-prd", Hash: true},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("stale entries = %+v, error = %v; want %+v", got, err, want)
	}
	if !reflect.DeepEqual(snapshot, snapshotDoctorPath(t, root)) {
		t.Fatal("lock scan wrote repository files")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := baseline.StaleOwnedLockEntries(ctx, root, names); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled scan error = %v", err)
	}
}

func TestDoctorKeepsReadinessWhenOwnedLockComparisonIsUnavailable(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"absent", "malformed", "schema", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			filename := filepath.Join(root, "skills-lock.json")
			switch kind {
			case "malformed":
				mustWrite(t, filename, "{broken")
			case "schema":
				mustWrite(t, filename, `{"version":3,"skills":{}}`)
			case "symlink":
				mustWrite(t, filepath.Join(root, "target"), `{"version":1,"skills":{}}`)
				if err := os.Symlink("target", filename); err != nil {
					t.Fatal(err)
				}
			case "directory":
				mustMkdir(t, filename)
			}
			for _, failed := range []bool{false, true} {
				readiness := skills.RepositoryReadiness{}
				if failed {
					readiness.MissingOwned = []string{"roundfix"}
				}
				dependencies := doctorDependencies{
					resolveExternal: func(string) ([]string, bool, error) { return nil, true, nil },
					checkSkills:     func(context.Context, string, []string) (skills.RepositoryReadiness, error) { return readiness, nil },
				}
				before := snapshotDoctorPath(t, root)
				got := repositorySkillsCheck(t.Context(), dependencies, root)
				want := doctorSkillReadinessResult(readiness, nil)
				if kind != "absent" {
					// The comparison error is the loader's diagnostic, with no readiness change.
					entries, err := baseline.StaleOwnedLockEntries(t.Context(), root, skills.Names())
					if err == nil || len(entries) != 0 {
						t.Fatalf("invalid lock scan = %+v, error = %v", entries, err)
					}
					want.Detail += "; skills-lock.json comparison unavailable: " + err.Error()
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("result = %+v, want %+v", got, want)
				}
				if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) {
					t.Fatal("Doctor changed repository files")
				}
			}
		})
	}
}

func TestDoctorOwnedLockWarningKeepsReadinessDetailAndNextActions(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name        string
		unversioned bool
	}{
		{name: "external failure"},
		{name: "unversioned", unversioned: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := doctorOwnedLockFixture(t, true)
			dependencies := doctorOwnedLockDependencies()
			if test.unversioned {
				filename := filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md")
				content := regexp.MustCompile(`(?m)^[ \t]*version:.*$`).ReplaceAllString(string(mustReadBytes(t, filename)), "")
				mustWrite(t, filename, content)
			} else {
				dependencies.resolveExternal = func(string) ([]string, bool, error) { return []string{"testing-boss", "agentic-cli-design"}, true, nil }
			}
			external, _, err := dependencies.resolveExternal(root)
			if err != nil {
				t.Fatal(err)
			}
			readiness, readinessErr := skills.CheckRepositoryWithExternal(t.Context(), root, external)
			want := doctorSkillReadinessResult(readiness, readinessErr)
			got := repositorySkillsCheck(t.Context(), dependencies, root)
			if got.Status != want.Status || (test.unversioned && got.Status != CheckStatusUnversioned) || (!test.unversioned && got.Status != CheckStatusFailed) {
				t.Fatalf("status = %q, readiness status = %q", got.Status, want.Status)
			}
			if !strings.Contains(got.Detail, "; DR-SKILL-LOCK-STALE:") {
				t.Fatalf("detail suffix order = %q", got.Detail)
			}
			if want.NextAction != "" {
				want.NextAction += " && "
			}
			want.NextAction += ownedSkillsNextAction
			if got.NextAction != want.NextAction {
				t.Fatalf("next action = %q, want %q", got.NextAction, want.NextAction)
			}
		})
	}
}

func TestStaleOwnedLockEntriesDetectsAnyRefAndInvalidHashes(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		members string
		hash    bool
	}{
		{name: "empty ref", members: `"ref":"","computedHash":"<hash>"`},
		{name: "null ref", members: `"ref":null,"computedHash":"<hash>"`},
		{name: "nonstring ref", members: `"ref":false,"computedHash":"<hash>"`},
		{name: "missing hash", hash: true},
		{name: "nonstring hash", members: `"computedHash":false`, hash: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := doctorOwnedLockFixture(t, false)
			hash, err := skills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", "roundfix"))
			if err != nil {
				t.Fatal(err)
			}
			members := strings.ReplaceAll(test.members, "<hash>", hash)
			mustWrite(t, filepath.Join(root, "skills-lock.json"), `{"version":1,"skills":{"roundfix":{`+members+`}}}`)
			got, err := baseline.StaleOwnedLockEntries(t.Context(), root, []string{"roundfix"})
			if err != nil || len(got) != 1 || got[0].Skill != "roundfix" || got[0].Hash != test.hash {
				t.Fatalf("stale entries = %+v, error = %v; want roundfix with Hash=%t", got, err, test.hash)
			}
		})
	}
}
