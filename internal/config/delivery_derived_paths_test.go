// Boundary: config proposal parsing, precedence, and derived-path validation.
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"strings"
	"testing"
)

func TestProjectConfigReadsDerivedPathDeclarations(t *testing.T) {
	user := []byte("delivery:\n  derived_paths:\n    - paths: [old]\n      regenerate: old-command\n")
	project := []byte("delivery:\n  derived_paths:\n    - paths: [generated.txt, cache/, 'assets/*.json']\n      regenerate: make generate\n")
	config, err := ResolveConfigProposal(user, project)
	if err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(t.TempDir(), "repo")
	home := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	if err := os.MkdirAll(filepath.Join(home, ".roundfix"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".roundfix", "config.yml"), user, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), project, 0600); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(LoadOptions{WorkDir: repo, HomeDir: home})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(loaded.Config.Delivery, config.Delivery) {
		t.Fatalf("loaded Project Config=%+v", loaded.Config.Delivery)
	}
	want := []DerivedPathDeclaration{{Paths: []string{"generated.txt", "cache/", "assets/*.json"}, Regenerate: "make generate"}}
	if !reflect.DeepEqual(config.Delivery.DerivedPaths, want) {
		t.Fatalf("declarations=%+v", config.Delivery.DerivedPaths)
	}
	for name, expected := range map[string]bool{"generated.txt": true, "cache/nested/file": true, "assets/a.json": true, "cacheish/a": false, "assets/nested/a.json": false} {
		if want[0].Matches(name) != expected {
			t.Errorf("match %q", name)
		}
	}
	without, err := ResolveConfigProposal(nil, nil)
	if err != nil || !reflect.DeepEqual(without, Builtin()) {
		t.Fatalf("absent key changed defaults: %v", err)
	}
	cleared, err := ResolveConfigProposal(user, []byte("delivery:\n  derived_paths: []\n"))
	if err != nil || len(cleared.Delivery.DerivedPaths) != 0 {
		t.Fatalf("explicit empty list did not override user: %v", err)
	}
}

func TestDerivedPathDeclarationsRefuseUnsafeEntries(t *testing.T) {
	for _, declaration := range []string{
		"paths: []\n      regenerate: make generate",
		"paths: [/tmp/file]\n      regenerate: make generate",
		"paths: [../file]\n      regenerate: make generate",
		"paths: [a/../file]\n      regenerate: make generate",
		"paths: [a/./file]\n      regenerate: make generate",
		"paths: ['']\n      regenerate: make generate",
		"paths: [file]\n      regenerate: ' '",
		"paths: ['[']\n      regenerate: make generate",
	} {
		t.Run(declaration, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte("delivery:\n  derived_paths:\n    - "+declaration+"\n"))
			if err == nil || !strings.Contains(err.Error(), "delivery.derived_paths") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
