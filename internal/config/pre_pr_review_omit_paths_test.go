// Boundary: review omission config parsing, precedence, validation, and matching.
package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPrePRReviewOmitPathsLoadsAndValidates(t *testing.T) {
	t.Parallel()

	user := []byte("pre_pr_review:\n  provider: claude\n  omit_paths: [user.json]\n")
	for _, tc := range []struct {
		name     string
		project  []byte
		paths    []string
		provider string
		source   string
	}{
		{name: "absent", provider: "codex", source: "default"},
		{name: "inherits user", project: []byte("{}"), paths: []string{"user.json"}, provider: "claude", source: "user"},
		{name: "project replaces user", project: []byte("pre_pr_review:\n  omit_paths: [gen/a.json, cache/, 'assets/*.json', 'literal/[a].json']\n"), paths: []string{"gen/a.json", "cache/", "assets/*.json", "literal/[a].json"}, provider: "claude", source: "user"},
		{name: "empty clears user", project: []byte("pre_pr_review:\n  omit_paths: []\n"), paths: []string{}, provider: "claude", source: "user"},
		{name: "provider still overrides", project: []byte("pre_pr_review:\n  provider: none\n"), paths: []string{"user.json"}, provider: "none", source: "project"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			userContent := user
			if tc.name == "absent" {
				userContent = nil
			}
			cfg, err := ResolveConfigProposal(userContent, tc.project)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg.PrePRReview.OmitPaths, tc.paths) || cfg.PrePRReview.Provider != tc.provider || cfg.PrePRReview.Source != tc.source {
				t.Fatalf("review config=%+v, want paths=%v provider=%q source=%q", cfg.PrePRReview, tc.paths, tc.provider, tc.source)
			}

			home, repo := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(repo, ".git"))
			if userContent != nil {
				mustMkdir(t, filepath.Join(home, ".roundfix"))
				mustWrite(t, filepath.Join(home, userConfigRelPath), string(userContent))
			}
			if tc.project != nil {
				mustWrite(t, filepath.Join(repo, projectConfigName), string(tc.project))
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil || !reflect.DeepEqual(loaded.Config.PrePRReview, cfg.PrePRReview) {
				t.Fatalf("loaded=%+v error=%v, want %+v", loaded.Config.PrePRReview, err, cfg.PrePRReview)
			}
		})
	}

	t.Run("matches whole paths", func(t *testing.T) {
		t.Parallel()
		cfg, err := ResolveConfigProposal(nil, []byte("pre_pr_review:\n  omit_paths: [gen/a.json, cache/, 'assets/*.json', 'literal/[a].json']\n"))
		if err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string]bool{
			"gen/a.json": true, "gen/b.json": false,
			"cache/nested/file": true, "cacheish/file": false, "cache": false,
			"assets/a.json": true, "assets/nested/a.json": false,
			"literal/[a].json": true, "literal/a.json": true,
			"src/a.go": false,
		} {
			if got := cfg.PrePRReview.Omits(name); got != want {
				t.Errorf("Omits(%q)=%v, want %v", name, got, want)
			}
		}
		if (PrePRReview{}).Omits("src/a.go") {
			t.Fatal("zero value omits authored source")
		}
	})

	for _, entry := range []string{"", ".", "./", "/tmp/file", "../file", "a/../file", "a/./file", "a//file", "a//", `a\file`, "["} {
		t.Run("refuses "+fmt.Sprintf("%q", entry), func(t *testing.T) {
			t.Parallel()
			_, err := ResolveConfigProposal(nil, []byte(fmt.Sprintf("pre_pr_review:\n  omit_paths: [%q]\n", entry)))
			kind := "unsafe path"
			if entry == "[" {
				kind = "invalid pattern"
			}
			want := fmt.Sprintf("pre_pr_review.omit_paths has %s %q", kind, entry)
			if err == nil || !strings.Contains(err.Error(), `parse config "Project Config proposal": `+want) {
				t.Fatalf("error=%v, want parse config label and %q", err, want)
			}
			cfg := Builtin()
			cfg.PrePRReview.OmitPaths = []string{entry}
			if err := Validate(cfg); err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("Validate error=%v, want %q", err, want)
			}
		})
	}
	for _, value := range []string{"[42]", "[true]", "[null]", "[{}]", "[[]]"} {
		t.Run("refuses non-string entry "+value, func(t *testing.T) {
			t.Parallel()
			_, err := ResolveConfigProposal(nil, []byte("pre_pr_review:\n  omit_paths: "+value+"\n"))
			if err == nil || !strings.Contains(err.Error(), `parse config "Project Config proposal": pre_pr_review.omit_paths has non-string path`) {
				t.Fatalf("error=%v, want labeled non-string path error", err)
			}
		})
	}
	for _, value := range []string{"gen/a.json", "{}", "null", "true", "42"} {
		t.Run("refuses non-list "+value, func(t *testing.T) {
			t.Parallel()
			_, err := ResolveConfigProposal(nil, []byte("pre_pr_review:\n  omit_paths: "+value+"\n"))
			want := `parse config "Project Config proposal": pre_pr_review.omit_paths must be a list of paths`
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v, want %q", err, want)
			}
		})
	}
}
