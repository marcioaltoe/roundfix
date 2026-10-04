// Boundary: item-binary declaration parsing, precedence, and path validation.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"roundfix/internal/gittest"
	"testing"
)

func itemBinaryLoadOptions(t *testing.T, user, project []byte) LoadOptions {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	home := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	if user != nil {
		if err := os.MkdirAll(filepath.Join(home, ".roundfix"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(home, ".roundfix", "config.yml"), user, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if project != nil {
		if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), project, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return LoadOptions{WorkDir: repo, HomeDir: home}
}

func checkItemBinaryConfig(t *testing.T, user, project []byte, want ItemBinaryDeclaration) {
	t.Helper()
	proposal, err := ResolveConfigProposal(user, project)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(itemBinaryLoadOptions(t, user, project))
	if err != nil {
		t.Fatal(err)
	}
	for name, config := range map[string]Config{"proposal": proposal, "load": loaded.Config} {
		wantDeclared := want != (ItemBinaryDeclaration{})
		if got := config.Delivery.ItemBinary; got != want || got.Declared() != wantDeclared {
			t.Errorf("%s declaration = %+v (declared %t), want %+v (declared %t)", name, got, got.Declared(), want, wantDeclared)
		}
	}
}

func TestDeliveryItemBinaryIsReadFromProjectConfig(t *testing.T) {
	project := []byte("delivery:\n  item_binary:\n    build: make build\n    path: bin/roundfix\n")
	checkItemBinaryConfig(t, nil, project, ItemBinaryDeclaration{Build: "make build", Path: "bin/roundfix"})
}

func TestDeliveryItemBinaryProjectReplacesUser(t *testing.T) {
	user := []byte("delivery:\n  item_binary:\n    build: user-build\n    path: user/roundfix\n")
	project := []byte("delivery:\n  item_binary:\n    build: make build\n    path: bin/roundfix\n")
	checkItemBinaryConfig(t, user, project, ItemBinaryDeclaration{Build: "make build", Path: "bin/roundfix"})
	t.Run("absent project declaration inherits user", func(t *testing.T) {
		checkItemBinaryConfig(t, user, []byte("delivery:\n  derived_paths: []\n"), ItemBinaryDeclaration{Build: "user-build", Path: "user/roundfix"})
	})
	t.Run("partial project declaration cannot inherit user fields", func(t *testing.T) {
		checkItemBinaryError(t, user, []byte("delivery:\n  item_binary:\n    build: make build\n"), "delivery.item_binary requires build and path")
	})
}

func TestDeliveryItemBinaryIsUndeclaredByDefault(t *testing.T) {
	if (ItemBinaryDeclaration{}).Declared() {
		t.Fatal("zero value is declared")
	}
	checkItemBinaryConfig(t, nil, nil, ItemBinaryDeclaration{})
}

func checkItemBinaryError(t *testing.T, user, project []byte, want string) {
	t.Helper()
	_, proposalErr := ResolveConfigProposal(user, project)
	_, loadErr := Load(itemBinaryLoadOptions(t, user, project))
	for name, err := range map[string]error{"proposal": proposalErr, "load": loadErr} {
		if err == nil {
			t.Errorf("%s accepted invalid declaration, want %q", name, want)
			continue
		}
		// Parsing adds scope context; compare the underlying contract error exactly.
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		if err.Error() != want {
			t.Errorf("%s error = %q, want %q", name, err, want)
		}
	}
}

func TestDeliveryItemBinaryRefusesAnIncompleteDeclaration(t *testing.T) {
	for name, declaration := range map[string]string{
		"missing build": "path: bin/roundfix",
		"missing path":  "build: make build",
		"empty build":   "build: ''\n    path: bin/roundfix",
		"empty path":    "build: make build\n    path: ''",
		"blank build":   "build: ' ' \n    path: bin/roundfix",
		"blank path":    "build: make build\n    path: ' '",
		"empty fields":  "build: ''\n    path: ''",
		"empty mapping": "{}",
	} {
		t.Run(name, func(t *testing.T) {
			project := []byte("delivery:\n  item_binary:\n    " + declaration + "\n")
			checkItemBinaryError(t, nil, project, "delivery.item_binary requires build and path")
		})
	}
}

func TestDeliveryItemBinaryRefusesAnUnsafePath(t *testing.T) {
	for name, entry := range map[string]string{
		"absolute":       "/abs/roundfix",
		"parent":         "../roundfix",
		"nested parent":  "bin/../roundfix",
		"backslash":      `bin\roundfix`,
		"dot":            ".",
		"dot segment":    "bin/./roundfix",
		"double slash":   "bin//roundfix",
		"trailing slash": "bin/roundfix/",
	} {
		t.Run(name, func(t *testing.T) {
			project := []byte("delivery:\n  item_binary:\n    build: make build\n    path: '" + entry + "'\n")
			checkItemBinaryError(t, nil, project, fmt.Sprintf("delivery.item_binary has unsafe path %q", entry))
		})
	}
}

func TestDeliveryItemBinaryRefusesAnUnknownKey(t *testing.T) {
	project := []byte("delivery:\n  item_binary:\n    build: make build\n    path: bin/roundfix\n    extra: true\n")
	checkItemBinaryError(t, nil, project, "delivery.item_binary.extra is not a supported config key")
	t.Run("user scope", func(t *testing.T) {
		checkItemBinaryError(t, project, nil, "delivery.item_binary.extra is not a supported config key")
	})
}
