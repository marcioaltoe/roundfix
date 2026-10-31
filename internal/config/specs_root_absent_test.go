package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSpecsRootOrAbsentReadsAnAbsentRootInsideTheRepository(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, configured                                  string
		docsPresent, absolute, linkedAncestor, linkedRepo bool
		builtIn                                           bool
	}{
		{name: "built-in with docs missing", configured: "docs/specs", builtIn: true},
		{name: "built-in with docs present", configured: "docs/specs", docsPresent: true, builtIn: true},
		{name: "relative root", configured: "plans/active"},
		{name: "absolute root", configured: "plans/active", absolute: true},
		{name: "ancestor linked inside repository", configured: "docs/specs", linkedAncestor: true, builtIn: true},
		{name: "repository reached through link", configured: "docs/specs", linkedRepo: true, builtIn: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			if tt.docsPresent {
				mustMkdir(t, filepath.Join(repo, "docs"))
			}
			if tt.linkedAncestor {
				target := filepath.Join(repo, "actual-docs")
				mustMkdir(t, target)
				if err := os.Symlink(target, filepath.Join(repo, "docs")); err != nil {
					t.Fatal(err)
				}
			}
			if tt.linkedRepo {
				link := filepath.Join(t.TempDir(), "repo")
				if err := os.Symlink(repo, link); err != nil {
					t.Fatal(err)
				}
				repo = link
			}
			loaded := Loaded{Config: Builtin(), GitRoot: repo}
			loaded.Config.Specs.Root = tt.configured
			if tt.absolute {
				loaded.Config.Specs.Root = filepath.Join(repo, tt.configured)
			}
			wantPath := filepath.Join(repo, tt.configured)
			got, err := ResolveSpecsRootOrAbsent(loaded, "")
			want := SpecsRoot{Path: wantPath, Absent: true, BuiltInRoot: tt.builtIn}
			if err != nil || got != want {
				t.Fatalf("got %+v, %v; want %+v, nil", got, err, want)
			}
			if _, err := os.Lstat(wantPath); !os.IsNotExist(err) {
				t.Fatalf("resolver created the absent root: %v", err)
			}
			wantErr := fmt.Sprintf("specs.root resolved to %q, which does not exist; create the directory or update specs.root", wantPath)
			if err := got.RequirePresent(); err == nil || err.Error() != wantErr {
				t.Fatalf("RequirePresent() = %v; want %s", err, wantErr)
			}
			strict, err := ResolveSpecsRoot(loaded, "")
			if strict != (SpecsRoot{}) || err == nil || err.Error() != wantErr {
				t.Fatalf("strict = %+v, %v; want zero value, %s", strict, err, wantErr)
			}
			mustMkdir(t, wantPath)
			present, err := ResolveSpecsRootOrAbsent(loaded, repo)
			want.Absent = false
			if err != nil || present != want || present.RequirePresent() != nil {
				t.Fatalf("present = %+v, %v; want %+v", present, err, want)
			}
			strict, err = ResolveSpecsRoot(loaded, repo)
			if err != nil || strict != present {
				t.Fatalf("strict present = %+v, %v; want %+v", strict, err, present)
			}
		})
	}
}

func TestResolveSpecsRootOrAbsentStillRefusesAMissingExternalRoot(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"absolute external root", "external linked ancestor"} {
		t.Run(name, func(t *testing.T) {
			repo, external := t.TempDir(), t.TempDir()
			loaded := Loaded{Config: Builtin(), GitRoot: repo}
			root := filepath.Join(external, "missing", "specs")
			if name == "external linked ancestor" {
				if err := os.Symlink(external, filepath.Join(repo, "docs")); err != nil {
					t.Fatal(err)
				}
				root = filepath.Join(repo, "docs", "specs")
			}
			loaded.Config.Specs.Root = root
			assertAbsentSpecsRootRefusal(t, loaded, fmt.Sprintf("specs.root resolved to %q, which does not exist; create the directory or update specs.root", root))
		})
	}
}

func TestResolveSpecsRootOrAbsentStillRefusesADanglingLink(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"root link", "ancestor link"} {
		t.Run(name, func(t *testing.T) {
			repo := t.TempDir()
			loaded := Loaded{Config: Builtin(), GitRoot: repo}
			root := filepath.Join(repo, "docs", "specs")
			link := root
			if name == "root link" {
				mustMkdir(t, filepath.Dir(root))
			} else {
				link = filepath.Dir(root)
			}
			if err := os.Symlink(filepath.Join(repo, "missing"), link); err != nil {
				t.Fatal(err)
			}
			assertAbsentSpecsRootRefusal(t, loaded, fmt.Sprintf("specs.root resolved to %q, which does not exist; create the directory or update specs.root", root))
		})
	}
}

func TestResolveSpecsRootOrAbsentStillRefusesAFile(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	root := filepath.Join(repo, "docs", "specs")
	mustMkdir(t, filepath.Dir(root))
	mustWrite(t, root, "not a directory")
	loaded := Loaded{Config: Builtin(), GitRoot: repo}
	assertAbsentSpecsRootRefusal(t, loaded, fmt.Sprintf("specs.root resolved to %q, which is not a directory; update specs.root to a directory", root))
}

func assertAbsentSpecsRootRefusal(t *testing.T, loaded Loaded, want string) {
	t.Helper()
	for _, resolve := range []struct {
		name string
		call func(Loaded, string) (SpecsRoot, error)
	}{{"OrAbsent", ResolveSpecsRootOrAbsent}, {"strict", ResolveSpecsRoot}} {
		root, err := resolve.call(loaded, loaded.GitRoot)
		if root != (SpecsRoot{}) || err == nil || err.Error() != want {
			t.Errorf("%s = %+v, %v; want zero value, %s", resolve.name, root, err, want)
		}
	}
}
