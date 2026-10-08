package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRetentionDaysDefaultsToThirty(t *testing.T) {
	cfg, err := ResolveConfigProposal(nil, nil)
	if err != nil || cfg.Store.RunRetentionDays != 30 {
		t.Fatalf("config=%+v error=%v", cfg.Store, err)
	}
}

func TestRunRetentionDaysAcceptsOnlySevenFifteenOrThirty(t *testing.T) {
	for _, value := range []string{"7", "15", "30", "0", "10", "31", "-7", "30d", `"30"`, "null", "30.0"} {
		t.Run(value, func(t *testing.T) {
			cfg, err := ResolveConfigProposal([]byte("store:\n  run_retention_days: "+value+"\n"), nil)
			if value == "7" || value == "15" || value == "30" {
				if err != nil || fmt.Sprint(cfg.Store.RunRetentionDays) != value {
					t.Fatalf("config=%+v error=%v", cfg.Store, err)
				}
			} else if err == nil || err.Error() != `parse config "User Config proposal": store.run_retention_days must be 7, 15 or 30` {
				t.Fatalf("error=%v", err)
			}
		})
	}
}

func TestRunRetentionDaysIsUserConfigOnly(t *testing.T) {
	home, repo := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".roundfix"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".roundfix", "config.yml"), []byte("store:\n  run_retention_days: 15\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".roundfixrc.yml"), []byte("store:\n  run_retention_days: invalid\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var warnings bytes.Buffer
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo, Stderr: &warnings})
	if err != nil || loaded.Config.Store.RunRetentionDays != 15 {
		t.Fatalf("config=%+v error=%v", loaded.Config.Store, err)
	}
	if warnings.String() != "config: store.run_retention_days in Project Config is ignored; set store.run_retention_days in User Config\n" {
		t.Fatalf("warnings=%q", warnings.String())
	}
}

func TestInitUserConfigCarriesRunRetentionDays(t *testing.T) {
	home, repo := t.TempDir(), t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, scope := range []string{InitScopeUser, InitScopeProject} {
		result, err := Init(context.Background(), InitOptions{HomeDir: home, WorkDir: repo, Scope: scope})
		if err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(result.Path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "run_retention_days: 30") != (scope == InitScopeUser) {
			t.Fatalf("unexpected %s template", scope)
		}
		if !strings.Contains(string(content), "journal_retention: 336h") {
			t.Fatal("Journal Retention missing")
		}
	}
}
