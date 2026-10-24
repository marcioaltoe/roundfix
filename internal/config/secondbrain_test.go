// Suite: Secondbrain configuration.
// Invariant: only User Config chooses the local root; projects may choose a namespace.
// Boundary IN: temporary layered YAML files.
// Boundary OUT: no host configuration or external services.
package config

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestSecondbrainRootIsReadFromUserConfig(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"absolute", "~/brain", "~"} {
		t.Run(value, func(t *testing.T) {
			home := t.TempDir()
			root := value
			want := filepath.Join(home, "brain")
			if value == "absolute" {
				root = want
			}
			if value == "~" {
				want = home
			}
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix/config.yml"), "secondbrain:\n  root: \""+root+"\"\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Secondbrain.Root != want {
				t.Fatalf("root=%q want=%q", loaded.Config.Secondbrain.Root, want)
			}
		})
	}
}

func TestProjectConfigCannotSetTheSecondbrainRoot(t *testing.T) {
	t.Parallel()
	home, repo := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(repo, ".git"))
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "secondbrain:\n  root: invalid-relative-root\n")
	var warnings bytes.Buffer
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, Stderr: &warnings})
	if err != nil {
		t.Fatal(err)
	}
	want := "config: secondbrain.root in Project Config is ignored; set secondbrain.root in User Config\n"
	if warnings.String() != want || loaded.Config.Secondbrain.Root != "" {
		t.Fatalf("warning=%q root=%q", warnings.String(), loaded.Config.Secondbrain.Root)
	}
}

func TestSecondbrainKeysRefuseInvalidValues(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ key, value, want string }{
		{"root", "relative", "secondbrain.root must be an absolute path or start with ~/"},
		{"inbox_namespace", "_triaged", "secondbrain.inbox_namespace must be one lowercase path segment of letters, digits, '.', '_' or '-', and not _triaged"},
		{"inbox_namespace", "Fluxus", "secondbrain.inbox_namespace must be one lowercase path segment of letters, digits, '.', '_' or '-', and not _triaged"},
		{"inbox_namespace", "a/b", "secondbrain.inbox_namespace must be one lowercase path segment of letters, digits, '.', '_' or '-', and not _triaged"},
	} {
		t.Run(tc.value, func(t *testing.T) {
			home := t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix/config.yml"), "secondbrain:\n  "+tc.key+": "+tc.value+"\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
			if err == nil || err.Error() != tc.want {
				t.Fatalf("error=%v want=%q", err, tc.want)
			}
		})
	}
}

func TestTheInboxNamespaceMayComeFromProjectConfig(t *testing.T) {
	t.Parallel()
	home, repo := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustMkdir(t, filepath.Join(repo, ".git"))
	root := filepath.Join(home, "brain")
	mustWrite(t, filepath.Join(home, ".roundfix/config.yml"), "secondbrain:\n  root: "+root+"\n  inbox_namespace: user-project\n")
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "secondbrain:\n  root: /ignored\n  inbox_namespace: project-choice\n")
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Secondbrain.Root != root || loaded.Config.Secondbrain.InboxNamespace != "project-choice" {
		t.Fatalf("config=%+v", loaded.Config.Secondbrain)
	}
}
