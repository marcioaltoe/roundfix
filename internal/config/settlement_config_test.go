package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryAtSettlementDefaultsToTrue(t *testing.T) {
	t.Parallel()
	home, work := t.TempDir(), t.TempDir()
	// Both config files exist but omit the new key.
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustMkdir(t, filepath.Join(work, ".git"))
	mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "verification:\n  concurrency: 2\n")
	mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "verification:\n  concurrency: 3\n")
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Config.Verification.RepositoryAtSettlement {
		t.Fatal("absent key must default to true")
	}
}

func TestRepositoryAtSettlementReadsFalse(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"user", "project", "project overrides user"} {
		t.Run(scope, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustMkdir(t, filepath.Join(work, ".git"))
			path := filepath.Join(home, ".roundfix", "config.yml")
			if scope != "user" {
				path = filepath.Join(work, ".roundfixrc.yml")
			}
			if scope == "project overrides user" {
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "verification:\n  repository_at_settlement: true\n")
			}
			mustWrite(t, path, "verification:\n  repository_at_settlement: false\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Verification.RepositoryAtSettlement {
				t.Fatal("explicit false was ignored")
			}
		})
	}
}

func TestRepositoryAtSettlementRefusesANonBoolean(t *testing.T) {
	t.Parallel()
	for _, value := range []string{`"false"`, "1", "[]", "{}", "perhaps", "null"} {
		t.Run(value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "verification:\n  repository_at_settlement: "+value+"\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err == nil || !strings.Contains(err.Error(), "verification.repository_at_settlement") {
				t.Fatalf("invalid %s: error = %v", value, err)
			}
		})
	}
}

func TestConfigTemplateCarriesRepositoryAtSettlement(t *testing.T) {
	t.Parallel()
	for _, template := range []string{DefaultConfigYAML(), DefaultProjectConfigYAML()} {
		if !strings.Contains(template, "repository_at_settlement: true") {
			t.Fatal("template omitted default")
		}
		config, err := ResolveConfigProposal([]byte(template), nil)
		if err != nil {
			t.Fatal(err)
		}
		if !config.Verification.RepositoryAtSettlement {
			t.Fatal("template must enable repository verification at settlement")
		}
	}
}
