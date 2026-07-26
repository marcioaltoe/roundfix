package skills

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestCheckRepositoryMatchesRealRepository(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve real repository root: %v", err)
	}
	got, err := CheckRepository(t.Context(), root)
	if err != nil {
		t.Fatalf("check real repository: %v", err)
	}
	if !got.Ready() {
		t.Fatalf("real repository readiness = %#v, want ready", got)
	}
}

func TestCheckRepositoryHonorsPreCanceledContext(t *testing.T) {
	tests := []struct {
		name    string
		context func(*testing.T) context.Context
		want    error
	}{
		{
			name: "canceled",
			context: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			want: context.Canceled,
		},
		{
			name: "deadline exceeded",
			context: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithDeadline(t.Context(), time.Unix(1, 0))
				t.Cleanup(cancel)
				return ctx
			},
			want: context.DeadlineExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := test.context(t)
			root := filepath.Join(t.TempDir(), "must-not-be-inspected")

			_, err := CheckRepository(ctx, root)
			if !errors.Is(err, test.want) {
				t.Fatalf("CheckRepository() error = %v, want errors.Is(_, %v)", err, test.want)
			}

			_, err = checkRepository(ctx, root, nil, nil, nil)
			if !errors.Is(err, test.want) {
				t.Fatalf("checkRepository() error = %v, want errors.Is(_, %v)", err, test.want)
			}
			var repositoryErr *RepositoryReadinessError
			if !errors.As(err, &repositoryErr) {
				t.Fatalf("checkRepository() error type = %T, want RepositoryReadinessError", err)
			}
		})
	}
}

func TestCheckRepositoryReportsReadyRequiredSetWithoutMutation(t *testing.T) {
	root := writeReadyRepositoryFixture(t)
	before := snapshotRepositoryFixture(t, root)

	got, err := CheckRepository(t.Context(), root)
	if err != nil {
		t.Fatalf("check ready repository: %v", err)
	}
	if !got.Ready() {
		t.Fatalf("expected repository to be ready, got %#v", got)
	}
	if got.OwnedRequired != 14 || got.ExternalRequired != 25 {
		t.Fatalf("required counts = owned %d external %d, want 14 and 25", got.OwnedRequired, got.ExternalRequired)
	}
	if after := snapshotRepositoryFixture(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("CheckRepository mutated fixture:\nbefore=%#v\nafter=%#v", before, after)
	}
}

func TestCheckRepositoryClassifiesMissingAndOutdatedSkills(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
		want   RepositoryReadiness
	}{
		{
			name: "missing owned",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "archive-spec")); err != nil {
					t.Fatalf("remove owned skill: %v", err)
				}
			},
			want: RepositoryReadiness{MissingOwned: []string{"archive-spec"}},
		},
		{
			name: "missing external",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "agentic-cli-design")); err != nil {
					t.Fatalf("remove external skill: %v", err)
				}
			},
			want: RepositoryReadiness{MissingExternal: []string{"agentic-cli-design"}},
		},
		{
			name: "changed owned file",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				writeSkillHashTestFile(t, filepath.Join(root, ".agents", "skills", "roundfix"), "SKILL.md", "changed\n")
			},
			want: RepositoryReadiness{OutdatedOwned: []string{"roundfix"}},
		},
		{
			name: "added owned file",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				writeSkillHashTestFile(t, filepath.Join(root, ".agents", "skills", "roundfix"), "unexpected.txt", "extra\n")
			},
			want: RepositoryReadiness{OutdatedOwned: []string{"roundfix"}},
		},
		{
			name: "removed owned file",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, ".agents", "skills", "roundfix", "agents", "openai.yaml")
				if err := os.Remove(path); err != nil {
					t.Fatalf("remove owned artifact: %v", err)
				}
			},
			want: RepositoryReadiness{OutdatedOwned: []string{"roundfix"}},
		},
		{
			name: "changed external file",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				writeSkillHashTestFile(t, filepath.Join(root, ".agents", "skills", "agentic-cli-design"), "SKILL.md", "changed\n")
			},
			want: RepositoryReadiness{OutdatedExternal: []string{"agentic-cli-design"}},
		},
		{
			name: "added external file",
			mutate: func(t *testing.T, root string) {
				t.Helper()
				writeSkillHashTestFile(t, filepath.Join(root, ".agents", "skills", "agentic-cli-design"), "unexpected.txt", "extra\n")
			},
			want: RepositoryReadiness{OutdatedExternal: []string{"agentic-cli-design"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			test.mutate(t, root)

			got, err := CheckRepository(t.Context(), root)
			if err != nil {
				t.Fatalf("check repository: %v", err)
			}
			got.OwnedRequired = 0
			got.ExternalRequired = 0
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("readiness = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestCheckRepositoryClassifiesMissingSharedSkillDirectories(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "missing .agents", path: ".agents"},
		{name: "missing .agents/skills", path: filepath.Join(".agents", "skills")},
	}

	wantOwned := Names()
	sort.Strings(wantOwned)
	wantExternal := Recommended()
	sort.Strings(wantExternal)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			if err := os.RemoveAll(filepath.Join(root, test.path)); err != nil {
				t.Fatalf("remove shared skill directory %q: %v", test.path, err)
			}

			got, err := CheckRepository(t.Context(), root)
			if err != nil {
				t.Fatalf("check repository: %v", err)
			}
			if !reflect.DeepEqual(got.MissingOwned, wantOwned) {
				t.Fatalf("missing owned = %v, want %v", got.MissingOwned, wantOwned)
			}
			if !reflect.DeepEqual(got.MissingExternal, wantExternal) {
				t.Fatalf("missing external = %v, want %v", got.MissingExternal, wantExternal)
			}
		})
	}
}

func TestCheckRepositoryIgnoresUnrelatedSkillsAndLockEntries(t *testing.T) {
	root := writeReadyRepositoryFixture(t)
	writeSkillHashTestFile(t, filepath.Join(root, ".agents", "skills", "unrelated"), "SKILL.md", "unrelated\n")
	lock := readRepositoryLockFixture(t, root)
	lock.Skills["../ignored-unsafe-entry"] = repositoryLockSkillFixture{ComputedHash: "not-a-hash"}
	writeRepositoryLockFixture(t, root, lock)

	got, err := CheckRepository(t.Context(), root)
	if err != nil {
		t.Fatalf("check repository with unrelated entries: %v", err)
	}
	if !got.Ready() {
		t.Fatalf("unrelated entries changed readiness: %#v", got)
	}
}

func TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames(t *testing.T) {
	tests := []struct {
		name  string
		write func(*testing.T, string)
	}{
		{
			name: "missing lock",
			write: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "skills-lock.json")); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "malformed lock",
			write: func(t *testing.T, root string) {
				if err := os.WriteFile(filepath.Join(root, "skills-lock.json"), []byte("{"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "wrong version",
			write: func(t *testing.T, root string) {
				lock := readRepositoryLockFixture(t, root)
				lock.Version = 2
				writeRepositoryLockFixture(t, root, lock)
			},
		},
		{
			name: "missing required entry",
			write: func(t *testing.T, root string) {
				lock := readRepositoryLockFixture(t, root)
				delete(lock.Skills, Recommended()[0])
				writeRepositoryLockFixture(t, root, lock)
			},
		},
		{
			name: "invalid required hash",
			write: func(t *testing.T, root string) {
				lock := readRepositoryLockFixture(t, root)
				lock.Skills[Recommended()[0]] = repositoryLockSkillFixture{ComputedHash: strings.Repeat("A", 64)}
				writeRepositoryLockFixture(t, root, lock)
			},
		},
		{
			name: "unreadable lock artifact shape",
			write: func(t *testing.T, root string) {
				path := filepath.Join(root, "skills-lock.json")
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			test.write(t, root)

			_, err := CheckRepository(t.Context(), root)
			if err == nil {
				t.Fatal("expected repository lock error")
			}
			var repositoryErr *RepositoryReadinessError
			if !errors.As(err, &repositoryErr) {
				t.Fatalf("expected RepositoryReadinessError, got %T: %v", err, err)
			}
			if repositoryErr.Ownership != RepositoryOwnershipExternal {
				t.Fatalf("lock error ownership = %q, want external", repositoryErr.Ownership)
			}
			if !strings.Contains(err.Error(), filepath.Join(root, "skills-lock.json")) {
				t.Fatalf("lock error does not name path: %v", err)
			}
		})
	}

	root := t.TempDir()
	outside := filepath.Join(filepath.Dir(root), "outside")
	writeSkillHashTestFile(t, outside, "marker", "must not be read\n")
	_, err := checkRepository(t.Context(), root, []string{"../outside"}, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "unsafe") {
		t.Fatalf("expected unsafe required name error, got %v", err)
	}
}

func TestCheckRepositoryWrapsFilesystemCauses(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")

	_, err := CheckRepository(t.Context(), root)
	if err == nil {
		t.Fatal("expected missing repository root error")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected wrapped fs.ErrNotExist, got %v", err)
	}
	var repositoryErr *RepositoryReadinessError
	if !errors.As(err, &repositoryErr) {
		t.Fatalf("expected RepositoryReadinessError, got %T: %v", err, err)
	}
	if repositoryErr.Path != root || repositoryErr.Ownership != "" {
		t.Fatalf("repository root error = %#v, want unclassified path %q", repositoryErr, root)
	}
}

func TestCheckRepositoryRejectsSymlinkedAuthoritiesBeforeReadingTargets(t *testing.T) {
	tests := []struct {
		name      string
		authority string
		target    func(*testing.T, string) string
	}{
		{
			name:      "symlinked .agents",
			authority: ".agents",
			target: func(_ *testing.T, outside string) string {
				return filepath.Join(outside, ".agents")
			},
		},
		{
			name:      "symlinked .agents/skills",
			authority: filepath.Join(".agents", "skills"),
			target: func(_ *testing.T, outside string) string {
				return filepath.Join(outside, ".agents", "skills")
			},
		},
		{
			name:      "symlinked skills-lock.json",
			authority: "skills-lock.json",
			target: func(t *testing.T, outside string) string {
				t.Helper()
				target := filepath.Join(outside, "malformed-lock.json")
				if err := os.WriteFile(target, []byte("{"), 0o644); err != nil {
					t.Fatalf("write malformed lock target: %v", err)
				}
				return target
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeReadyRepositoryFixture(t)
			outside := writeReadyRepositoryFixture(t)
			authority := filepath.Join(root, test.authority)
			if err := os.RemoveAll(authority); err != nil {
				t.Fatalf("remove repository authority %q: %v", authority, err)
			}
			if err := os.Symlink(test.target(t, outside), authority); err != nil {
				t.Skipf("create authority symlink: %v", err)
			}

			_, err := CheckRepository(t.Context(), root)
			if err == nil {
				t.Fatal("expected symlinked authority to fail readiness")
			}
			var repositoryErr *RepositoryReadinessError
			if !errors.As(err, &repositoryErr) {
				t.Fatalf("expected RepositoryReadinessError, got %T: %v", err, err)
			}
			if repositoryErr.Path != authority {
				t.Fatalf("authority error path = %q, want %q", repositoryErr.Path, authority)
			}
			wantOwnership := RepositoryOwnership("")
			if test.authority == "skills-lock.json" {
				wantOwnership = RepositoryOwnershipExternal
			}
			if repositoryErr.Ownership != wantOwnership {
				t.Fatalf("authority error ownership = %q, want %q", repositoryErr.Ownership, wantOwnership)
			}
			if strings.Contains(repositoryErr.Operation, "decode") {
				t.Fatalf("authority target was decoded before rejection: %v", err)
			}
		})
	}
}

func TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering(t *testing.T) {
	root := writeReadyRepositoryFixture(t)
	ownedPath := filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md")
	if err := os.Remove(ownedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "skills-lock.json"), ownedPath); err != nil {
		t.Skipf("create owned artifact symlink: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "agentic-cli-design")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "archive-spec")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "autoresearch")); err != nil {
		t.Fatal(err)
	}

	got, err := CheckRepository(t.Context(), root)
	if err != nil {
		t.Fatalf("check repository with owned symlink: %v", err)
	}
	if !reflect.DeepEqual(got.MissingOwned, []string{"archive-spec"}) ||
		!reflect.DeepEqual(got.MissingExternal, []string{"agentic-cli-design", "autoresearch"}) ||
		!reflect.DeepEqual(got.OutdatedOwned, []string{"roundfix"}) {
		t.Fatalf("unstable classifications: %#v", got)
	}

	externalRoot := writeReadyRepositoryFixture(t)
	externalPath := filepath.Join(externalRoot, ".agents", "skills", "agentic-cli-design", "linked")
	if err := os.Symlink(filepath.Join(externalRoot, "skills-lock.json"), externalPath); err != nil {
		t.Skipf("create external symlink: %v", err)
	}
	_, err = CheckRepository(t.Context(), externalRoot)
	if err == nil || !strings.Contains(err.Error(), externalPath) {
		t.Fatalf("expected external symlink error naming %q, got %v", externalPath, err)
	}

	specialRoot := writeReadyRepositoryFixture(t)
	specialPath := filepath.Join(specialRoot, ".agents", "skills", "agentic-cli-design", "socket")
	listener, err := net.Listen("unix", specialPath)
	if err != nil {
		t.Skipf("create special filesystem entry: %v", err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Errorf("close special-entry listener: %v", err)
		}
	})
	_, err = CheckRepository(t.Context(), specialRoot)
	if err == nil || !strings.Contains(err.Error(), specialPath) {
		t.Fatalf("expected external special-entry error naming %q, got %v", specialPath, err)
	}
}

type repositoryLockFixture struct {
	Version int                                   `json:"version"`
	Skills  map[string]repositoryLockSkillFixture `json:"skills"`
}

type repositoryLockSkillFixture struct {
	ComputedHash string `json:"computedHash"`
}

func writeReadyRepositoryFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	skillsRoot := filepath.Join(root, ".agents", "skills")
	files, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	for _, file := range files {
		writeSkillHashTestFile(t, skillsRoot, file.Path, string(file.Data))
	}

	lock := repositoryLockFixture{
		Version: 1,
		Skills:  make(map[string]repositoryLockSkillFixture, len(Recommended())),
	}
	for _, name := range Recommended() {
		skillRoot := filepath.Join(skillsRoot, name)
		writeSkillHashTestFile(t, skillRoot, "SKILL.md", name+"\n")
		hash, err := SkillFolderHash(t.Context(), skillRoot)
		if err != nil {
			t.Fatalf("hash external fixture %q: %v", name, err)
		}
		lock.Skills[name] = repositoryLockSkillFixture{ComputedHash: hash}
	}
	writeRepositoryLockFixture(t, root, lock)
	return root
}

func readRepositoryLockFixture(t *testing.T, root string) repositoryLockFixture {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "skills-lock.json"))
	if err != nil {
		t.Fatalf("read repository lock fixture: %v", err)
	}
	var lock repositoryLockFixture
	if err := json.Unmarshal(data, &lock); err != nil {
		t.Fatalf("decode repository lock fixture: %v", err)
	}
	return lock
}

func writeRepositoryLockFixture(t *testing.T, root string, lock repositoryLockFixture) {
	t.Helper()
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		t.Fatalf("encode repository lock fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills-lock.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write repository lock fixture: %v", err)
	}
}

func snapshotRepositoryFixture(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relative)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot repository fixture: %v", err)
	}
	return snapshot
}
