package skills

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCheckRepositoryChecksAnExternalSkillOnlyForPresence(t *testing.T) {
	t.Parallel()
	const name = "agentic-cli-design"
	for _, test := range []struct {
		name         string
		mutate       func(*testing.T, string, string)
		missing      bool
		missingEntry bool
		rootLink     bool
	}{
		{name: "changed file", mutate: func(t *testing.T, root, folder string) {
			writeSkillHashTestFile(t, folder, "SKILL.md", "changed\n")
		}},
		{name: "added file", mutate: func(t *testing.T, root, folder string) {
			writeSkillHashTestFile(t, folder, "extra.txt", "added\n")
		}},
		{name: "nested symbolic link", mutate: func(t *testing.T, root, folder string) {
			if err := os.Symlink(filepath.Join(root, "skills-lock.json"), filepath.Join(folder, "linked")); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "lock entry without computedHash", mutate: func(t *testing.T, root, folder string) {
			if err := os.WriteFile(filepath.Join(root, "skills-lock.json"), []byte(`{"version":1,"skills":{"agentic-cli-design":{}}}`), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "missing directory", missing: true, mutate: func(t *testing.T, root, folder string) {
			if err := os.RemoveAll(folder); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "no lock entry", missingEntry: true, mutate: func(t *testing.T, root, folder string) {
			lock := readRepositoryLockFixture(t, root)
			delete(lock.Skills, name)
			writeRepositoryLockFixture(t, root, lock)
		}},
		{name: "skill root is a symbolic link", rootLink: true, mutate: func(t *testing.T, root, folder string) {
			target := filepath.Join(root, "external-skill")
			if err := os.Rename(folder, target); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, folder); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			folder := filepath.Join(root, ".agents", "skills", name)
			test.mutate(t, root, folder)
			got, err := CheckRepositoryWithExternal(t.Context(), root, []string{name})
			if test.missingEntry || test.rootLink {
				var readinessErr *RepositoryReadinessError
				if !errors.As(err, &readinessErr) {
					t.Fatalf("error = %v, want RepositoryReadinessError", err)
				}
				if readinessErr.Ownership != RepositoryOwnershipExternal {
					t.Fatalf("ownership = %q, want external", readinessErr.Ownership)
				}
				if test.missingEntry && !reflect.DeepEqual(readinessErr.MissingExternal, []string{name}) {
					t.Fatalf("missing lock entries = %v, want [%s]", readinessErr.MissingExternal, name)
				}
				if test.rootLink && (readinessErr.Path != folder || !strings.Contains(err.Error(), "symbolic links are not supported")) {
					t.Fatalf("root symlink error = %+v, want path %q and symbolic link refusal", readinessErr, folder)
				}
				return
			}
			if err != nil {
				t.Fatalf("check repository: %v", err)
			}
			var wantMissing []string
			if test.missing {
				wantMissing = []string{name}
			}
			if !reflect.DeepEqual(got.MissingExternal, wantMissing) || got.Ready() == test.missing {
				t.Fatalf("readiness = %+v, want missing=%v ready=%t", got, wantMissing, !test.missing)
			}
		})
	}
}
