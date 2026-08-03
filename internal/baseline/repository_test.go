// Suite: Baseline repository inspection
// Invariant: repository identity and bounded carrier evidence are deterministic, root-confined, and read-only.
// Boundary IN: local Git lineage, instruction carriers, docs/agents files, aliases, preimages, and findings.
// Boundary OUT: Baseline decisions, plan postimages, apply transactions, and semantic classification.

package baseline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"sort"
	"strings"
	"syscall"
	"testing"
)

func TestInventoryWalkIgnoresTransientErrorsInsideExcludedTrees(t *testing.T) {
	t.Parallel()

	builder := &inventoryBuilder{}
	if err := builder.walk(".git/objects/maintenance.lock", nil, fs.ErrNotExist); err != nil {
		t.Fatalf("walk ignored transient Git path: %v", err)
	}
	if len(builder.blocking) != 0 {
		t.Fatalf("ignored transient Git path produced blocking findings: %+v", builder.blocking)
	}

	if err := builder.walk("docs/agents/AGENTS.md", nil, fs.ErrNotExist); err != nil {
		t.Fatalf("walk bounded transient path: %v", err)
	}
	if len(builder.blocking) != 1 ||
		builder.blocking[0].Code != "baseline.inventory.path-unreadable" {
		t.Fatalf("bounded transient path findings = %+v", builder.blocking)
	}
}

func TestRepositoryIdentityEquivalentClones(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "seed\n")
	commitInspectionRepository(t, repo, "seed")

	clone := filepath.Join(t.TempDir(), "clone")
	runInspectionCommand(t, "", "git", "clone", "--no-local", repo, clone)

	first, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect first clone: %v", err)
	}
	second, err := InspectRepository(context.Background(), clone, nil)
	if err != nil {
		t.Fatalf("inspect second clone: %v", err)
	}
	if first.Identity.Digest != second.Identity.Digest {
		t.Fatalf("repository identity differs by clone path: first=%+v second=%+v", first.Identity, second.Identity)
	}
	if first.Identity.ObjectFormat != "sha1" || len(first.Identity.RootCommits) != 1 {
		t.Fatalf("repository identity = %+v, want sha1 with one root commit", first.Identity)
	}
	if first.Root == second.Root {
		t.Fatalf("test clones unexpectedly share root %q", first.Root)
	}
}

func TestRepositoryIdentityAcceptsDetachedDirtyWithoutUpstream(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "seed\n")
	commitInspectionRepository(t, repo, "seed")
	runInspectionCommand(t, repo, "git", "checkout", "--detach")
	writeInspectionFile(t, repo, "scratch.txt", "dirty and unrelated\n")

	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect detached dirty repository without upstream: %v", err)
	}
	if inspection.Identity.Digest == "" {
		t.Fatal("repository identity digest is empty")
	}
	for _, preimage := range inspection.Snapshot.Preimages {
		if preimage.Path == "scratch.txt" {
			t.Fatalf("unrelated dirty path entered bounded preimages: %+v", inspection.Snapshot.Preimages)
		}
	}
}

func TestRepositoryIdentityRequiresCommittedGitWorktree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root func(*testing.T) string
		want string
	}{
		{
			name: "not a Git worktree",
			root: func(t *testing.T) string { return t.TempDir() },
			want: "detect Git worktree root",
		},
		{
			name: "Git worktree without commit",
			root: func(t *testing.T) string { return newInspectionRepository(t) },
			want: "requires at least one commit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InspectRepository(context.Background(), tt.root(t), nil)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("InspectRepository error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "README.md", "seed\n")
	writeInspectionFile(t, repo, "AGENTS.md", "root policy\n")
	writeInspectionFile(t, repo, "packages/api/AGENTS.md", "nested policy\n")
	writeInspectionFile(t, repo, "docs/agents/backend.md", "backend guide\n")
	writeInspectionFile(t, repo, "docs/agents/contracts/data.json", "{}\n")
	writeInspectionFile(t, repo, "node_modules/pkg/AGENTS.md", "ignored\n")
	writeInspectionFile(t, repo, "skills/example/CLAUDE.md", "ignored\n")
	commitInspectionRepository(t, repo, "seed")
	writeInspectionFile(t, repo, "scratch.txt", "unrelated dirty bytes\n")

	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect bounded repository: %v", err)
	}
	gotCarriers := make([]string, 0, len(inspection.Snapshot.Carriers))
	for _, carrier := range inspection.Snapshot.Carriers {
		gotCarriers = append(gotCarriers, carrier.Path)
	}
	wantCarriers := []string{
		"AGENTS.md",
		"docs/agents/backend.md",
		"docs/agents/contracts/data.json",
		"packages/api/AGENTS.md",
	}
	if !reflect.DeepEqual(gotCarriers, wantCarriers) {
		t.Fatalf("carrier paths = %v, want %v", gotCarriers, wantCarriers)
	}
	if !hasRepositoryFinding(inspection.Snapshot.Warnings, "baseline.inventory.nested-carrier-conflict", "packages/api/AGENTS.md") {
		t.Fatalf("nested carrier warning missing: %+v", inspection.Snapshot.Warnings)
	}
	if len(inspection.Snapshot.Blocking) != 0 {
		t.Fatalf("regular bounded repository unexpectedly blocked: %+v", inspection.Snapshot.Blocking)
	}
	for _, path := range []string{"AGENTS.md", "CLAUDE.md"} {
		if !hasRepositoryPreimage(inspection.Snapshot.Preimages, path) {
			t.Fatalf("mutable root path %q missing from preimages: %+v", path, inspection.Snapshot.Preimages)
		}
	}
	for _, excluded := range []string{"scratch.txt", "node_modules/pkg/AGENTS.md", "skills/example/CLAUDE.md"} {
		if hasRepositoryPreimage(inspection.Snapshot.Preimages, excluded) {
			t.Fatalf("unbounded path %q entered preimages", excluded)
		}
	}
}

func TestInstructionAliasRetainsOneSourceEvidence(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "policy/shared.md", "shared policy\n")
	if err := os.Symlink("policy/shared.md", filepath.Join(repo, "AGENTS.md")); err != nil {
		t.Fatalf("create AGENTS alias: %v", err)
	}
	if err := os.Symlink("policy/shared.md", filepath.Join(repo, "CLAUDE.md")); err != nil {
		t.Fatalf("create CLAUDE alias: %v", err)
	}
	commitInspectionRepository(t, repo, "seed")

	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect safe aliases: %v", err)
	}
	if len(inspection.Snapshot.Blocking) != 0 {
		t.Fatalf("safe aliases blocked: %+v", inspection.Snapshot.Blocking)
	}
	if len(inspection.Snapshot.Sources) != 1 {
		t.Fatalf("source evidence count = %d, want one: %+v", len(inspection.Snapshot.Sources), inspection.Snapshot.Sources)
	}
	source := inspection.Snapshot.Sources[0]
	if source.Path != "policy/shared.md" || source.ContentIdentity == "" {
		t.Fatalf("source evidence = %+v", source)
	}
	if len(inspection.Snapshot.Carriers) != 2 {
		t.Fatalf("alias carriers = %+v, want two", inspection.Snapshot.Carriers)
	}
	for _, carrier := range inspection.Snapshot.Carriers {
		if carrier.Kind != CarrierAlias || carrier.TargetPath != "policy/shared.md" ||
			carrier.ContentIdentity != source.ContentIdentity {
			t.Fatalf("alias carrier = %+v, source = %+v", carrier, source)
		}
	}
}

func TestInstructionAliasRetainsDirectSourcePreimageIdentity(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "shared policy\n")
	if err := os.Symlink("AGENTS.md", filepath.Join(repo, "CLAUDE.md")); err != nil {
		t.Fatalf("create CLAUDE alias: %v", err)
	}
	commitInspectionRepository(t, repo, "seed")

	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect direct source alias: %v", err)
	}
	for _, preimage := range inspection.Snapshot.Preimages {
		if preimage.Path != "AGENTS.md" {
			continue
		}
		if preimage.Kind != PreimageRegular || preimage.ContentIdentity == "" {
			t.Fatalf("AGENTS preimage lost its trusted identity: %+v", preimage)
		}
		return
	}
	t.Fatal("AGENTS preimage is missing")
}

func TestInstructionAliasUnsafeTargetsBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(*testing.T, string)
	}{
		{
			name: "absolute external",
			setup: func(t *testing.T, repo string) {
				outside := filepath.Join(t.TempDir(), "outside.md")
				writeInspectionFile(t, filepath.Dir(outside), filepath.Base(outside), "outside\n")
				mustInspectionSymlink(t, outside, filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "relative escape",
			setup: func(t *testing.T, repo string) {
				mustInspectionSymlink(t, "../outside.md", filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "cycle",
			setup: func(t *testing.T, repo string) {
				mustInspectionSymlink(t, "CLAUDE.md", filepath.Join(repo, "AGENTS.md"))
				mustInspectionSymlink(t, "AGENTS.md", filepath.Join(repo, "CLAUDE.md"))
			},
		},
		{
			name: "directory",
			setup: func(t *testing.T, repo string) {
				if err := os.Mkdir(filepath.Join(repo, "policy"), 0o755); err != nil {
					t.Fatalf("create policy directory: %v", err)
				}
				mustInspectionSymlink(t, "policy", filepath.Join(repo, "AGENTS.md"))
			},
		},
		{
			name: "special file",
			setup: func(t *testing.T, repo string) {
				special := filepath.Join(repo, "policy.pipe")
				if err := syscall.Mkfifo(special, 0o600); err != nil {
					t.Fatalf("create named pipe: %v", err)
				}
				mustInspectionSymlink(t, "policy.pipe", filepath.Join(repo, "AGENTS.md"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInspectionRepository(t)
			tt.setup(t, repo)
			commitInspectionRepository(t, repo, "seed")

			inspection, err := InspectRepository(context.Background(), repo, nil)
			if err != nil {
				t.Fatalf("unsafe carrier must be reported in snapshot, got error: %v", err)
			}
			if !hasRepositoryFinding(inspection.Snapshot.Blocking, "baseline.inventory.unsafe-alias", "AGENTS.md") {
				t.Fatalf("unsafe alias was not apply-blocking: %+v", inspection.Snapshot.Blocking)
			}
		})
	}
}

func TestInstructionAliasUnreadableTargetBlocks(t *testing.T) {
	t.Parallel()

	if os.Geteuid() == 0 {
		t.Skip("root can read mode-000 fixtures")
	}
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "policy.md", "private policy\n")
	mustInspectionSymlink(t, "policy.md", filepath.Join(repo, "AGENTS.md"))
	commitInspectionRepository(t, repo, "seed")
	target := filepath.Join(repo, "policy.md")
	if err := os.Chmod(target, 0); err != nil {
		t.Fatalf("make alias target unreadable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(target, 0o644) })

	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("unreadable carrier must be reported in snapshot, got error: %v", err)
	}
	if !hasRepositoryFinding(inspection.Snapshot.Blocking, "baseline.inventory.unsafe-alias", "AGENTS.md") {
		t.Fatalf("unreadable alias was not apply-blocking: %+v", inspection.Snapshot.Blocking)
	}
}

func TestRepositoryInspectionNoMutation(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "root policy\n")
	writeInspectionFile(t, repo, "nested/CLAUDE.md", "nested policy\n")
	commitInspectionRepository(t, repo, "seed")
	writeInspectionFile(t, repo, "untracked.txt", "must remain byte-identical\n")

	before := snapshotInspectionTree(t, repo)
	if _, err := InspectRepository(context.Background(), repo, nil); err != nil {
		t.Fatalf("inspect repository: %v", err)
	}
	after := snapshotInspectionTree(t, repo)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("repository changed during inspection:\nbefore=%v\nafter=%v", before, after)
	}
}

func hasRepositoryFinding(findings []Finding, code, path string) bool {
	for _, finding := range findings {
		if finding.Code == code && finding.Path == path {
			return true
		}
	}
	return false
}

func hasRepositoryPreimage(preimages []Preimage, path string) bool {
	for _, preimage := range preimages {
		if preimage.Path == path {
			return true
		}
	}
	return false
}

func newInspectionRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runInspectionCommand(t, repo, "git", "init", "--quiet")
	gittest.AppendConfig(t, repo, "[user]\n\tname = Roundfix Test\n\temail = roundfix@example.test\n[commit]\n\tgpgsign = false\n")
	return repo
}

func commitInspectionRepository(t *testing.T, repo, message string) {
	t.Helper()
	runInspectionCommand(t, repo, "git", "add", "--all")
	runInspectionCommand(t, repo, "git", "commit", "--quiet", "-m", message)
}

func writeInspectionFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}

func mustInspectionSymlink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.Symlink(target, path); err != nil {
		t.Fatalf("create symlink %s -> %s: %v", path, target, err)
	}
}

func runInspectionCommand(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	if name == "git" {
		args = append([]string{
			"-c", "core.fsmonitor=false",
			"-c", "maintenance.auto=false",
		}, args...)
	}
	command := exec.Command(name, args...)
	command.Dir = dir
	command.Env = append(
		os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_OPTIONAL_LOCKS=0",
		"GIT_AUTHOR_DATE=2000-01-01T00:00:00Z",
		"GIT_COMMITTER_DATE=2000-01-01T00:00:00Z",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
}

type inspectionTreeEntry struct {
	Path string
	Mode fs.FileMode
	Link string
	Data []byte
}

func snapshotInspectionTree(t *testing.T, root string) []inspectionTreeEntry {
	t.Helper()
	var entries []inspectionTreeEntry
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := inspectionTreeEntry{Path: filepath.ToSlash(relative), Mode: info.Mode()}
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			item.Link, err = os.Readlink(path)
		case info.Mode().IsRegular():
			item.Data, err = os.ReadFile(path)
		}
		if err != nil {
			return err
		}
		entries = append(entries, item)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot repository tree: %v", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries
}

func TestRepositoryInspectionHonorsCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := InspectRepository(ctx, t.TempDir(), nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("InspectRepository error = %v, want context.Canceled", err)
	}
}

func TestRepositoryInspectionUsesNarrowReadOnlyGitCommands(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	runner := &recordingInspectionGitRunner{
		root: root,
		head: strings.Repeat("a", 40),
	}
	if _, err := InspectRepository(context.Background(), root, runner); err != nil {
		t.Fatalf("inspect with recording Git runner: %v", err)
	}
	want := [][]string{
		{"rev-parse", "--show-toplevel", "--show-object-format", "--verify", "HEAD^{commit}"},
		{"rev-list", "--max-parents=0", "HEAD"},
	}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("Git calls = %v, want narrow read-only commands %v", runner.calls, want)
	}
}

func TestRepositoryInspectionParsesCombinedResolutionPositionally(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	head := strings.Repeat("a", 40)
	runner := inspectionGitRunnerFunc(func(_ context.Context, _ string, args ...string) (string, error) {
		switch strings.Join(args, "\x00") {
		case "rev-parse\x00--show-toplevel\x00--show-object-format\x00--verify\x00HEAD^{commit}":
			return strings.Join([]string{root, head, "sha1"}, "\n"), nil
		default:
			return "", errors.New("unexpected Git command")
		}
	})

	_, err := InspectRepository(context.Background(), root, runner)
	want := fmt.Sprintf("detect Git object format: unsupported format %q", head)
	if err == nil || err.Error() != want {
		t.Fatalf("InspectRepository error = %v, want %q", err, want)
	}
}

func TestRepositoryInspectionCombinedResolutionErrorsRemainDistinct(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	head := strings.Repeat("a", 40)
	combinedCommand := "rev-parse\x00--show-toplevel\x00--show-object-format\x00--verify\x00HEAD^{commit}"
	tests := []struct {
		name   string
		runner GitRunner
		want   string
	}{
		{
			name: "non-repository directory",
			runner: inspectionGitRunnerFunc(func(_ context.Context, _ string, args ...string) (string, error) {
				if strings.Join(args, "\x00") == "rev-parse\x00--show-toplevel" {
					return "", errors.New("not a repository")
				}
				return "", errors.New("combined resolution failed")
			}),
			want: "detect Git worktree root: not a repository",
		},
		{
			name: "missing HEAD",
			runner: inspectionGitRunnerFunc(func(_ context.Context, _ string, args ...string) (string, error) {
				if strings.Join(args, "\x00") == "rev-parse\x00--show-toplevel" {
					return root, nil
				}
				return "", errors.New("HEAD is missing")
			}),
			want: "Baseline repository requires at least one commit: HEAD is missing",
		},
		{
			name: "unknown object format",
			runner: inspectionGitRunnerFunc(func(_ context.Context, _ string, args ...string) (string, error) {
				if strings.Join(args, "\x00") != combinedCommand {
					return "", errors.New("unexpected Git command")
				}
				return strings.Join([]string{root, "unknown", head}, "\n"), nil
			}),
			want: "detect Git object format: unsupported format \"unknown\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := InspectRepository(context.Background(), root, tt.runner)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("InspectRepository error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestRepositorySnapshotDigestChangesWithBoundedBytesOnly(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "one\n")
	commitInspectionRepository(t, repo, "seed")
	first, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect first snapshot: %v", err)
	}
	writeInspectionFile(t, repo, "scratch.txt", "unrelated\n")
	second, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect unrelated dirty state: %v", err)
	}
	if first.Snapshot.Digest != second.Snapshot.Digest {
		t.Fatalf("unrelated dirty file changed bounded digest: %s != %s", first.Snapshot.Digest, second.Snapshot.Digest)
	}
	writeInspectionFile(t, repo, "AGENTS.md", "two\n")
	third, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect changed carrier: %v", err)
	}
	if bytes.Equal([]byte(first.Snapshot.Digest), []byte(third.Snapshot.Digest)) {
		t.Fatalf("bounded carrier change did not change snapshot digest %q", first.Snapshot.Digest)
	}
}

type recordingInspectionGitRunner struct {
	root  string
	head  string
	calls [][]string
}

type inspectionGitRunnerFunc func(context.Context, string, ...string) (string, error)

func (run inspectionGitRunnerFunc) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	return run(ctx, workDir, args...)
}

func (runner *recordingInspectionGitRunner) RunGit(_ context.Context, _ string, args ...string) (string, error) {
	runner.calls = append(runner.calls, append([]string(nil), args...))
	switch strings.Join(args, "\x00") {
	case "rev-parse\x00--show-toplevel\x00--show-object-format\x00--verify\x00HEAD^{commit}":
		return strings.Join([]string{runner.root, "sha1", runner.head}, "\n"), nil
	case "rev-list\x00--max-parents=0\x00HEAD":
		return runner.head, nil
	default:
		return "", errors.New("unexpected Git command")
	}
}
