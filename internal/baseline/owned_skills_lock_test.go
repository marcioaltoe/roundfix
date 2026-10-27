package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func ownedLockFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	folder := filepath.Join(root, ".agents", "skills", "roundfix")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "SKILL.md"), []byte("installed bundle\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, _, err := inspectRestoreTarget(root, ".agents/skills/roundfix")
	if err != nil {
		t.Fatal(err)
	}
	return root, externalSkillsLockDigest(files)
}

func TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill(t *testing.T) {
	t.Parallel()
	root, hash := ownedLockFixture(t)
	stale := strings.Repeat("0", 64)
	before := `{ "version":1,"skills":{"external":{"computedHash":"` + stale + `"},"roundfix":{"extra":{"number":1e2,"escaped":"\u0061"}, "ref":"keep", "computedHash":"` + stale + `"},"missing":{"computedHash":"` + stale + `"},"nonobject":false},"extra":true }` + "\n"
	filename := filepath.Join(root, "skills-lock.json")
	if err := os.WriteFile(filename, []byte(before), 0o640); err != nil {
		t.Fatal(err)
	}
	names, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"missing", "roundfix", "roundfix", "nonobject", "unlisted"})
	if err != nil || !reflect.DeepEqual(names, []string{"roundfix"}) {
		t.Fatalf("names=%v err=%v", names, err)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(before, `"computedHash":"`+stale+`"},"missing"`, `"computedHash":"`+hash+`"},"missing"`, 1)
	if string(after) != want {
		t.Fatalf("lock=%s want=%s", after, want)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode=%v", info.Mode())
	}
}

func TestRecordInstalledOwnedSkillsLeavesAMatchingLockUntouched(t *testing.T) {
	t.Parallel()
	root, hash := ownedLockFixture(t)
	filename := filepath.Join(root, "skills-lock.json")
	before := `{"version":1,"skills":{"roundfix":{"computedHash":"` + hash + `"}}}`
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		names, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
		if err != nil || len(names) != 0 {
			t.Fatalf("names=%v err=%v", names, err)
		}
		after, err := os.Stat(filename)
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		if !os.SameFile(info, after) || !info.ModTime().Equal(after.ModTime()) || string(data) != before {
			t.Fatal("matching lock was written")
		}
	}
}

func TestRecordInstalledOwnedSkillsAddsNoEntryAndNoLock(t *testing.T) {
	t.Parallel()
	root, _ := ownedLockFixture(t)
	filename := filepath.Join(root, "skills-lock.json")
	names, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
	if err != nil || len(names) != 0 {
		t.Fatalf("names=%v err=%v", names, err)
	}
	if _, err := os.Lstat(filename); !os.IsNotExist(err) {
		t.Fatalf("lock exists: %v", err)
	}
	before := `{"version":1,"skills":{"external":{"extra":true}}}`
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	names, err = RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
	if err != nil || len(names) != 0 {
		t.Fatalf("names=%v err=%v", names, err)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != before {
		t.Fatalf("lock changed: %s", after)
	}
}

func TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"symlink", "directory", "malformed", "schema"} {
		t.Run(kind, func(t *testing.T) {
			root, _ := ownedLockFixture(t)
			filename := filepath.Join(root, "skills-lock.json")
			original := []byte("{broken")
			switch kind {
			case "symlink":
				target := filepath.Join(root, "target")
				if err := os.WriteFile(target, original, 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("target", filename); err != nil {
					t.Fatal(err)
				}
			case "directory":
				if err := os.Mkdir(filename, 0o755); err != nil {
					t.Fatal(err)
				}
			case "schema":
				original = []byte(`{"version":2,"skills":{}}`)
				fallthrough
			default:
				if err := os.WriteFile(filename, original, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			info, err := os.Lstat(filename)
			if err != nil {
				t.Fatal(err)
			}
			names, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"roundfix"})
			if err == nil || !strings.Contains(err.Error(), "skills-lock.json") || len(names) != 0 {
				t.Fatalf("names=%v err=%v", names, err)
			}
			after, err := os.Lstat(filename)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(info, after) || info.Mode() != after.Mode() {
				t.Fatal("unsafe lock replaced")
			}
			if kind != "directory" {
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				if string(data) != string(original) {
					t.Fatal("unsafe lock bytes changed")
				}
			}
		})
	}
}

func TestRecordInstalledOwnedSkillsSetsMissingHashesInSortedOrder(t *testing.T) {
	t.Parallel()
	root, hash := ownedLockFixture(t)
	for _, name := range []string{"qa-gate", "write-tasks"} {
		folder := filepath.Join(root, ".agents", "skills", name)
		if err := os.MkdirAll(folder, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(folder, "SKILL.md"), []byte("installed bundle\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	filename := filepath.Join(root, "skills-lock.json")
	before := `{"version":1,"skills":{"roundfix":{ },"write-tasks":{"extra":true },"qa-gate":{"computedHash":null}}}`
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	names, err := RecordInstalledOwnedSkills(t.Context(), root, []string{"write-tasks", "roundfix", "qa-gate"})
	if err != nil || !reflect.DeepEqual(names, []string{"qa-gate", "roundfix", "write-tasks"}) {
		t.Fatalf("names=%v err=%v", names, err)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(before, `{ }`, `{ "computedHash":"`+hash+`"}`, 1)
	want = strings.Replace(want, `{"extra":true }`, `{"extra":true ,"computedHash":"`+hash+`"}`, 1)
	want = strings.Replace(want, `"computedHash":null`, `"computedHash":"`+hash+`"`, 1)
	if string(after) != want {
		t.Fatalf("lock=%s want=%s", after, want)
	}
}
