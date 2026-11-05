package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name           string
		gitFile        bool
		locationIsRoot bool
	}{
		{name: "git directory"},
		{name: "git file", gitFile: true},
		{name: "location is the root", locationIsRoot: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			homeDir, gitRoot := t.TempDir(), t.TempDir()
			workDir := filepath.Join(gitRoot, "nested", "directory")
			mustMkdir(t, workDir)
			if test.gitFile {
				mustWrite(t, filepath.Join(gitRoot, ".git"), "gitdir: "+t.TempDir()+"\n")
			} else {
				gittest.MarkRepository(t, gitRoot)
			}
			location := filepath.Join(gitRoot, ".roundfix", "worktrees")
			if test.locationIsRoot {
				location = gitRoot
			}
			mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
			mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"),
				fmt.Sprintf("worktree:\n  location: %q\n", location))

			_, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
			if err == nil {
				t.Fatal("expected in-repository worktree.location to fail")
			}
			oldText := fmt.Sprintf("worktree.location must not be inside the repository tree: %q", location)
			want := fmt.Sprintf("%s is inside the Git root %q (found %q)", oldText, gitRoot, filepath.Join(gitRoot, ".git"))
			if err.Error() != want {
				t.Errorf("refusal = %q, want %q", err.Error(), want)
			}
			if !strings.HasPrefix(err.Error(), oldText) {
				t.Errorf("refusal = %q, want old-text prefix %q", err.Error(), oldText)
			}
		})
	}
}

func TestLoadWithOptionalWorktreeLocationReturnsTheRefusal(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		stray   bool
		gitFile bool
	}{
		{name: "stray git directory", stray: true},
		{name: "stray git file", stray: true, gitFile: true},
		{name: "repository"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			homeDir, gitRoot := t.TempDir(), t.TempDir()
			if test.stray {
				gitRoot = homeDir
			}
			workDir := filepath.Join(gitRoot, "drive")
			mustMkdir(t, workDir)
			if test.gitFile {
				mustWrite(t, filepath.Join(gitRoot, ".git"), "gitdir: "+t.TempDir()+"\n")
			} else {
				gittest.MarkRepository(t, gitRoot)
			}
			userConfigPath := filepath.Join(homeDir, ".roundfix", "config.yml")
			mustMkdir(t, filepath.Dir(userConfigPath))
			mustWrite(t, userConfigPath, "watch:\n  max_rounds: 7\n")

			opts := LoadOptions{HomeDir: homeDir, WorkDir: workDir}
			strict, strictErr := Load(opts)
			opts.WorktreeLocationOptional = true
			loaded, err := Load(opts)
			if err != nil {
				t.Fatalf("optional Load: %v", err)
			}
			if loaded.GitRoot != gitRoot || loaded.HomeDir != homeDir ||
				loaded.UserConfigPath != userConfigPath || loaded.ProjectConfigPath != filepath.Join(gitRoot, ".roundfixrc.yml") {
				t.Fatalf("optional Load lost its root or config paths: %+v", loaded)
			}
			if loaded.Config.Watch.MaxRounds != 7 {
				t.Errorf("User Config max_rounds = %d, want 7", loaded.Config.Watch.MaxRounds)
			}
			if test.stray {
				if strictErr == nil {
					t.Fatal("strict Load must refuse the location under the stray Git root")
				}
				if loaded.WorktreeLocationErr == nil || loaded.WorktreeLocationErr.Error() != strictErr.Error() {
					t.Errorf("optional refusal = %v, want %v", loaded.WorktreeLocationErr, strictErr)
				}
				wantConfig := Builtin()
				wantConfig.Watch.MaxRounds = 7
				wantConfig.Worktree.Location = ""
				if !reflect.DeepEqual(loaded.Config, wantConfig) {
					t.Errorf("optional config = %+v, want %+v", loaded.Config, wantConfig)
				}
				return
			}
			if strictErr != nil {
				t.Fatalf("strict Load in repository: %v", strictErr)
			}
			if loaded.WorktreeLocationErr != nil {
				t.Errorf("optional refusal in repository = %v, want nil", loaded.WorktreeLocationErr)
			}
			if want := filepath.Join(homeDir, ".roundfix", "worktrees"); loaded.Config.Worktree.Location != want {
				t.Errorf("resolved location = %q, want %q", loaded.Config.Worktree.Location, want)
			}
			if !reflect.DeepEqual(loaded, strict) {
				t.Errorf("optional Load = %+v, want strict Load %+v", loaded, strict)
			}
		})
	}
}
