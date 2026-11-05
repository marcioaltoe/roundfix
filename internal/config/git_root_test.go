// Suite: Git root discovery
// Invariant: only a .git entry Git accepts marks a Git root.
// Boundary IN: filesystem shapes, Git discovery, and configuration root consumers.
// Boundary OUT: CLI transcripts and skills project scope, owned by internal/cli.
package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

var gitRootShapes = []struct {
	name       string
	accepted   bool
	gitAnswers bool
}{
	{"only gk and info", false, true},
	{"empty directory", false, true},
	{"HEAD and objects without refs", false, true},
	{"HEAD that names nothing", false, true},
	{"file without gitdir", false, false},
	{"symbolic ref HEAD", true, true},
	{"detached HEAD", true, true},
	{"SHA-256 detached HEAD", true, true},
	{"common directory", true, true},
	{"gitdir file", true, true},
	{"gitdir file whose target is missing", true, false},
	{"git init repository", true, true},
	{"linked worktree", true, true},
	{"repository under a rejected directory", true, true},
	{"symlink HEAD", true, true},
	{"HEAD with whitespace after ref", true, true},
	{"short detached HEAD", false, true},
	{"HEAD naming a non-refs path", false, true},
	{"HEAD directory", false, true},
	{"missing objects", false, true},
	{"objects file", false, true},
	{"refs file", false, true},
	{"absolute common directory", true, true},
	{"symlink git directory", true, true},
}

func TestFindGitRootAcceptsOnlyAGitEntryGitAccepts(t *testing.T) {
	t.Parallel()
	for _, shape := range gitRootShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			_, root, work := makeGitRootShape(t, shape.name)
			want := ""
			if shape.accepted {
				want = root
			}
			if got := FindGitRoot(filepath.Join(work, ".")); got != want {
				t.Errorf("FindGitRoot = %q, want %q", got, want)
			}
			// A rejected inner entry must not hide an accepted ancestor.
			if !shape.accepted {
				outer := filepath.Dir(root)
				gittest.MarkRepository(t, outer)
				if got := FindGitRoot(work); got != outer {
					t.Errorf("FindGitRoot through rejected entry = %q, want %q", got, outer)
				}
			}
		})
	}
}

func TestFindGitRootAgreesWithGit(t *testing.T) {
	t.Parallel()
	for _, shape := range gitRootShapes {
		if !shape.gitAnswers {
			continue
		}
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			base, _, work := makeGitRootShape(t, shape.name)
			physicalBase, err := filepath.EvalSymlinks(base)
			if err != nil {
				t.Fatal(err)
			}
			command := exec.Command("git", "rev-parse", "--show-toplevel")
			command.Dir = work
			command.Env = append(gittest.IsolatedEnv(), "GIT_CEILING_DIRECTORIES="+filepath.Dir(physicalBase))
			output, err := command.CombinedOutput()
			want := ""
			if err == nil {
				want, err = filepath.EvalSymlinks(strings.TrimSpace(string(output)))
				if err != nil {
					t.Fatal(err)
				}
			} else if !strings.Contains(string(output), "not a git repository") {
				t.Fatalf("Git did not answer the discovery question: %v\n%s", err, output)
			}
			got := FindGitRoot(work)
			if got != "" {
				got, err = filepath.EvalSymlinks(got)
				if err != nil {
					t.Fatal(err)
				}
			}
			if got != want {
				t.Errorf("FindGitRoot = %q, Git = %q", got, want)
			}
		})
	}
}

func makeGitRootShape(t *testing.T, name string) (base, root, work string) {
	t.Helper()
	base = t.TempDir()
	root = filepath.Join(base, "shape")
	if name == "repository under a rejected directory" {
		mustMkdir(t, filepath.Join(base, ".git", "gk"))
		mustMkdir(t, filepath.Join(base, ".git", "info"))
		root = filepath.Join(base, "drive", "repo")
	}
	mustMkdir(t, root)
	gitDir := filepath.Join(root, ".git")
	switch name {
	case "only gk and info":
		mustMkdir(t, filepath.Join(gitDir, "gk"))
		mustMkdir(t, filepath.Join(gitDir, "info"))
	case "empty directory":
		mustMkdir(t, gitDir)
	case "file without gitdir":
		mustWrite(t, gitDir, "garbage\n")
	case "gitdir file", "symlink git directory":
		target := filepath.Join(base, "target")
		gittest.InitRepo(t, target)
		if name == "gitdir file" {
			mustWrite(t, gitDir, "gitdir: "+filepath.Join(target, ".git")+"\n")
		} else if err := os.Symlink(filepath.Join(target, ".git"), gitDir); err != nil {
			t.Fatal(err)
		}
	case "gitdir file whose target is missing":
		mustWrite(t, gitDir, "gitdir: "+filepath.Join(base, "missing")+"\n")
	case "git init repository", "repository under a rejected directory":
		gittest.InitRepo(t, root)
	case "linked worktree":
		source := filepath.Join(base, "source")
		gittest.InitRepo(t, source)
		gittest.Run(t, source, "commit", "--allow-empty", "-m", "initial")
		gittest.Run(t, source, "worktree", "add", "-b", "linked", root)
	default:
		mustMkdir(t, gitDir)
		head := "ref: refs/heads/main\n"
		switch name {
		case "HEAD that names nothing":
			head = "junk\n"
		case "detached HEAD":
			head = strings.Repeat("a", 40) + "\n"
		case "SHA-256 detached HEAD":
			head = strings.Repeat("B", 64) + "\n"
		case "short detached HEAD":
			head = strings.Repeat("a", 39) + "\n"
		case "HEAD naming a non-refs path":
			head = "ref: heads/main\n"
		case "HEAD with whitespace after ref":
			head = "ref:\t  refs/heads/main\n"
		}
		if name == "symlink HEAD" {
			if err := os.Symlink("refs/heads/main", filepath.Join(gitDir, "HEAD")); err != nil {
				t.Fatal(err)
			}
		} else if name == "HEAD directory" {
			mustMkdir(t, filepath.Join(gitDir, "HEAD"))
		} else {
			mustWrite(t, filepath.Join(gitDir, "HEAD"), head)
		}
		common := gitDir
		if name == "common directory" || name == "absolute common directory" {
			common = filepath.Join(base, "common")
			commonPath := common
			if name == "common directory" {
				commonPath = "../../common"
			}
			mustWrite(t, filepath.Join(gitDir, "commondir"), commonPath+"\n")
		}
		if name == "objects file" {
			mustWrite(t, filepath.Join(common, "objects"), "objects")
		} else if name != "missing objects" {
			mustMkdir(t, filepath.Join(common, "objects"))
		}
		if name == "refs file" {
			mustWrite(t, filepath.Join(common, "refs"), "refs")
		} else if name != "HEAD and objects without refs" {
			mustMkdir(t, filepath.Join(common, "refs"))
		}
	}
	work = filepath.Join(root, "sub")
	mustMkdir(t, work)
	return base, root, work
}

func TestLoadUnderARejectedGitDirectoryFindsNoGitRoot(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	work := filepath.Join(home, "drive")
	mustMkdir(t, work)
	mustMkdir(t, filepath.Join(home, ".git", "gk"))
	mustMkdir(t, filepath.Join(home, ".git", "info"))
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
	if err != nil {
		t.Fatalf("Load under rejected .git: %v", err)
	}
	if loaded.GitRoot != "" || loaded.ProjectConfigPath != "" {
		t.Errorf("GitRoot = %q, ProjectConfigPath = %q, want both empty", loaded.GitRoot, loaded.ProjectConfigPath)
	}
	if want := filepath.Join(home, ".roundfix", "worktrees"); loaded.Config.Worktree.Location != want {
		t.Errorf("worktree location = %q, want %q", loaded.Config.Worktree.Location, want)
	}
	_, err = Init(t.Context(), InitOptions{HomeDir: home, WorkDir: work, Scope: InitScopeProject})
	if err == nil || err.Error() != "project init requires a Git root; use --scope user outside a repository" {
		t.Errorf("project init error = %v", err)
	}
	_, _, _, err = profileConfigTarget(InitScopeProject, home, work)
	if err == nil || err.Error() != "project profiles configure requires a Git root; use --scope user outside a repository" {
		t.Errorf("project profiles configure error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".roundfixrc.yml")); !os.IsNotExist(err) {
		t.Errorf("project config at rejected root: %v, want not exist", err)
	}
}
