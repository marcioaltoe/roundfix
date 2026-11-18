package baseline

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestReadClaudeInstructionLinkClassifiesEveryShape(t *testing.T) {
	t.Parallel()

	const agents = "# Root instructions\nKeep repository policy.\n"
	const target = "docs/claude-notes.md"
	const body = "Managed guidance."
	block := managedRegionFixture("root.core", body)
	foreign := "<!-- BEGIN:tool-rules -->\nForeign guidance.\n<!-- END:tool-rules -->\n"
	rule := func(relative string, start, end int, first string) MisplacedRepositoryRule {
		return MisplacedRepositoryRule{Path: relative, StartLine: start, EndLine: end, FirstLine: first, MoveTo: SpecificRepositoryPath}
	}
	mismatch := rule(ClaudeInstructionLinkPath, 1, 5, "Local guidance.")
	mismatch.ManagedID = "root.core"
	tests := []struct {
		name        string
		content     string
		link        string
		absent      bool
		listed      bool
		gitMode     string
		agentsShape string
		readOnly    bool
		want        ClaudeInstructionLink
		convertible bool
		wantError   bool
	}{
		{name: "link", link: ClaudeInstructionLinkTarget, want: ClaudeInstructionLink{Shape: "link", Target: ClaudeInstructionLinkTarget}},
		{name: "dot_slash_target", link: "./AGENTS.md", want: ClaudeInstructionLink{Shape: "link", Target: "./AGENTS.md"}},
		{name: "absent", absent: true, want: ClaudeInstructionLink{Shape: "absent"}, convertible: true},
		{name: "empty_file", want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "import_line", content: " \n @AGENTS.md \n\n", want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "copy_of_agents", content: agents, want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "managed_regions_only", content: block, listed: true, want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "file_with_a_rule", content: "\n# Local rules\n- Keep this rule.\n", want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{rule(ClaudeInstructionLinkPath, 2, 3, "# Local rules")}}},
		{name: "listed_digest_mismatch", content: managedRegionFixture("root.core", "Local guidance."), listed: true, want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{mismatch}}},
		{name: "link_to_a_clean_file", link: target, content: block + foreign, want: ClaudeInstructionLink{Shape: "link-elsewhere", Target: target}, convertible: true},
		{name: "link_to_a_file_with_a_rule", link: target, content: foreign + "\n# Local rules\n- Keep this rule.\n", want: ClaudeInstructionLink{Shape: "link-elsewhere", Target: target, Passages: []MisplacedRepositoryRule{rule(target, 5, 6, "# Local rules")}}},
		{name: "unlinked_checkout", gitMode: "120000", want: ClaudeInstructionLink{Shape: "unlinked-checkout"}},
		{name: "agents_bytes_tracked_as_a_file", content: ClaudeInstructionLinkTarget, gitMode: "100644", want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "agents_missing", agentsShape: "missing", want: ClaudeInstructionLink{Shape: "agents-not-regular"}},
		{name: "agents_is_a_link", agentsShape: "link", want: ClaudeInstructionLink{Shape: "agents-not-regular"}},
		{name: "absolute_link_is_an_error", link: "absolute", wantError: true},
		{name: "reads_only", content: block + "\n# Local rules\n", listed: true, readOnly: true, want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{rule(ClaudeInstructionLinkPath, 7, 7, "# Local rules")}}},
		{name: "link_to_an_import_line", link: target, content: "@AGENTS.md\n", want: ClaudeInstructionLink{Shape: "link-elsewhere", Target: target}, convertible: true},
		{name: "link_with_missing_agents_is_an_error", link: ClaudeInstructionLinkTarget, agentsShape: "missing", wantError: true},
		{name: "link_with_linked_agents", link: ClaudeInstructionLinkTarget, agentsShape: "link", want: ClaudeInstructionLink{Shape: "link", Target: ClaudeInstructionLinkTarget}},
		{name: "agents_bytes_untracked", content: ClaudeInstructionLinkTarget, want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "agents_bytes_with_newline", content: ClaudeInstructionLinkTarget + "\n", want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{rule(ClaudeInstructionLinkPath, 1, 1, ClaudeInstructionLinkTarget)}}},
		{name: "link_to_a_copy_of_agents", link: target, content: agents, want: ClaudeInstructionLink{Shape: "link-elsewhere", Target: target}, convertible: true},
		{name: "import_beside_a_rule_is_a_passage", content: "@AGENTS.md\nKeep this rule.\n", want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{rule(ClaudeInstructionLinkPath, 1, 2, "@AGENTS.md")}}},
		{name: "foreign_block_splits_import_from_rule", content: "@AGENTS.md\n" + foreign + "Keep this rule.\n", want: ClaudeInstructionLink{Shape: "file", Passages: []MisplacedRepositoryRule{rule(ClaudeInstructionLinkPath, 5, 5, "Keep this rule.")}}},
		{name: "unlisted_managed_region", content: managedRegionFixture("root.core", "Local guidance."), want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "duplicate_managed_region", content: block + block, listed: true, want: ClaudeInstructionLink{Shape: "file"}, convertible: true},
		{name: "agents_is_a_directory", agentsShape: "directory", want: ClaudeInstructionLink{Shape: "agents-not-regular"}},
		{name: "claude_is_a_directory", agentsShape: "claude-directory", wantError: true},
		{name: "escaping_link_is_an_error", link: "../outside.md", wantError: true},
		{name: "missing_link_target_is_an_error", link: "missing.md", wantError: true},
		{name: "link_to_a_directory_is_an_error", link: "docs", wantError: true},
		{name: "indirect_escape_is_an_error", link: "docs/escape.md", wantError: true},
		{name: "cyclic_link_is_an_error", link: ClaudeInstructionLinkPath, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newInspectionRepository(t)
			writeInspectionFile(t, repo, ClaudeInstructionLinkTarget, agents)
			artifactPath := ClaudeInstructionLinkTarget
			if tt.listed {
				artifactPath = ClaudeInstructionLinkPath
			}
			writeManagedRefreshManifest(t, repo, ManifestArtifact{ID: "root.core", Path: artifactPath, Kind: "root-block", Module: "core", Template: "root.core", Version: ManifestVersion, Digest: managedRegionBodyDigest(body)})
			if tt.link == target {
				writeInspectionFile(t, repo, target, tt.content)
			}
			if tt.link == "docs" {
				if err := os.MkdirAll(filepath.Join(repo, "docs"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if tt.link == "docs/escape.md" {
				writeInspectionFile(t, repo, "docs/placeholder.md", "Untouched.\n")
				mustInspectionSymlink(t, "../../outside.md", filepath.Join(repo, "docs/escape.md"))
			}
			if !tt.absent {
				if tt.link != "" {
					link := tt.link
					if link == "absolute" {
						link = filepath.Join(repo, ClaudeInstructionLinkTarget)
					}
					mustInspectionSymlink(t, link, filepath.Join(repo, ClaudeInstructionLinkPath))
				} else if tt.gitMode == "120000" {
					mustInspectionSymlink(t, ClaudeInstructionLinkTarget, filepath.Join(repo, ClaudeInstructionLinkPath))
					commitInspectionRepository(t, repo, "Fixture instruction link")
					if err := os.Remove(filepath.Join(repo, ClaudeInstructionLinkPath)); err != nil {
						t.Fatal(err)
					}
					runInspectionCommand(t, repo, "git", "-c", "core.symlinks=false", "checkout", "--", ClaudeInstructionLinkPath)
				} else {
					writeInspectionFile(t, repo, ClaudeInstructionLinkPath, tt.content)
				}
			}
			if tt.gitMode == "100644" {
				commitInspectionRepository(t, repo, "Fixture regular instruction file")
			}
			if tt.agentsShape != "" {
				relative := ClaudeInstructionLinkTarget
				if tt.agentsShape == "claude-directory" {
					relative = ClaudeInstructionLinkPath
				}
				if err := os.Remove(filepath.Join(repo, relative)); err != nil {
					t.Fatal(err)
				}
				switch tt.agentsShape {
				case "link":
					writeInspectionFile(t, repo, "root-notes.md", agents)
					mustInspectionSymlink(t, "root-notes.md", filepath.Join(repo, relative))
				case "directory", "claude-directory":
					if err := os.Mkdir(filepath.Join(repo, relative), 0o755); err != nil {
						t.Fatal(err)
					}
				}
			}
			if tt.readOnly {
				for _, relative := range []string{ClaudeInstructionLinkPath, ClaudeInstructionLinkTarget, manifestPath} {
					if err := os.Chmod(filepath.Join(repo, relative), 0o440); err != nil {
						t.Fatal(err)
					}
				}
			}
			before := snapshotInspectionTree(t, repo)
			calls := 0
			ctx := t.Context()
			runner := inspectionGitRunnerFunc(func(gotCtx context.Context, dir string, args ...string) (string, error) {
				calls++
				if gotCtx != ctx || dir != repo || !reflect.DeepEqual(args, []string{"ls-files", "-s", "--", ClaudeInstructionLinkPath}) {
					t.Fatalf("unexpected Git request: context=%v dir=%q args=%q", gotCtx, dir, args)
				}
				return (ExecGitRunner{}).RunGit(gotCtx, dir, args...)
			})
			got, err := ReadClaudeInstructionLink(ctx, repo, runner)
			if after := snapshotInspectionTree(t, repo); !reflect.DeepEqual(before, after) {
				t.Fatal("reader changed fixture bytes, modes, links or paths")
			}
			wantCalls := 0
			if tt.gitMode == "120000" || tt.content == ClaudeInstructionLinkTarget {
				wantCalls = 1
			}
			if calls != wantCalls {
				t.Errorf("Git calls = %d, want %d", calls, wantCalls)
			}
			if tt.wantError {
				if err == nil || !strings.Contains(err.Error(), ClaudeInstructionLinkPath) {
					t.Fatalf("error = %v, want an error naming CLAUDE.md", err)
				}
			} else if err != nil {
				t.Fatalf("read instruction link: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("link = %#v, want %#v", got, tt.want)
			}
			if got.Convertible() != tt.convertible {
				t.Errorf("Convertible = %v, want %v", got.Convertible(), tt.convertible)
			}
		})
	}
}

func TestReadClaudeInstructionLinkPropagatesCancellationAndGitErrors(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"canceled", "git_error"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo := newInspectionRepository(t)
			writeInspectionFile(t, repo, ClaudeInstructionLinkTarget, "Root instructions.\n")
			writeInspectionFile(t, repo, ClaudeInstructionLinkPath, ClaudeInstructionLinkTarget)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			want := errors.New("index unreadable")
			if name == "canceled" {
				cancel()
				want = context.Canceled
			}
			runner := inspectionGitRunnerFunc(func(context.Context, string, ...string) (string, error) {
				if name == "canceled" {
					t.Fatal("canceled reader ran Git")
				}
				return "", want
			})
			_, err := ReadClaudeInstructionLink(ctx, repo, runner)
			if !errors.Is(err, want) {
				t.Fatalf("error = %v, want %v", err, want)
			}
		})
	}
}

type symlinkModeFileInfo struct{ mode fs.FileMode }

func (info symlinkModeFileInfo) Name() string       { return ClaudeInstructionLinkPath }
func (info symlinkModeFileInfo) Size() int64        { return int64(len(ClaudeInstructionLinkTarget)) }
func (info symlinkModeFileInfo) Mode() fs.FileMode  { return info.mode }
func (info symlinkModeFileInfo) ModTime() time.Time { return time.Time{} }
func (info symlinkModeFileInfo) IsDir() bool        { return false }
func (info symlinkModeFileInfo) Sys() any           { return nil }

// A plan is portable (ADR-0071): lstat reports a symbolic link's permission
// bits as 0777 on Linux and from the umask on macOS, so the bound preimage
// must not carry them, or the same repository plans to different digests.
func TestSymbolicLinkPreimageModeIsPlatformIndependent(t *testing.T) {
	t.Parallel()

	want := uint32(fs.ModeSymlink)
	for _, perm := range []fs.FileMode{0o777, 0o755, 0o700} {
		info := symlinkModeFileInfo{mode: fs.ModeSymlink | perm}
		if got := preimageFromInfo(ClaudeInstructionLinkPath, info, ClaudeInstructionLinkTarget).Mode; got != want {
			t.Fatalf("preimage mode for link permissions %o = %d, want %d", perm, got, want)
		}
	}
	regular := symlinkModeFileInfo{mode: 0o644}
	if got := preimageMode(regular); got != uint32(0o644) {
		t.Fatalf("regular preimage mode = %o, want 644", got)
	}

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ClaudeInstructionLinkTarget), []byte("# Root instructions\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	if err := os.Symlink(ClaudeInstructionLinkTarget, filepath.Join(root, ClaudeInstructionLinkPath)); err != nil {
		t.Fatalf("create link: %v", err)
	}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		t.Fatalf("open root: %v", err)
	}
	defer anchored.Close()
	state, err := captureTransactionState(anchored, ClaudeInstructionLinkPath)
	if err != nil {
		t.Fatalf("capture link state: %v", err)
	}
	if state.Kind != PreimageSymlink || state.Mode != want {
		t.Fatalf("captured link state = kind %q mode %d, want symlink mode %d", state.Kind, state.Mode, want)
	}
	builder := inventoryBuilder{root: anchored, preimages: make(map[string]Preimage)}
	builder.recordPath(ClaudeInstructionLinkPath)
	preimage := builder.preimages[ClaudeInstructionLinkPath]
	if preimage.Kind != PreimageSymlink || preimage.Mode != want || !transactionStateMatchesPreimage(state, preimage) {
		t.Fatalf("inventoried link preimage = %+v, want symlink mode %d matching the transaction state %+v", preimage, want, state)
	}
}
