// Suite: Project Config review data paths.
// Invariant: data paths replace or clear inherited values and match apart with the omit-path grammar.
// Boundary IN: User and Project Config proposals through ResolveConfigProposal.
// Boundary OUT: review prompts, segmentation, provider calls, and delivery regeneration.
package config

import (
	"reflect"
	"testing"
)

func TestPrePRReviewDataPathsLoadsAndValidates(t *testing.T) {
	t.Parallel()

	user := []byte("pre_pr_review:\n  provider: claude\n  data_paths: [user.json]\n")
	for _, tc := range []struct {
		name     string
		user     []byte
		project  []byte
		paths    []string
		provider string
		source   string
	}{
		{name: "absent", provider: "codex", source: "default"},
		{name: "inherits user", user: user, project: []byte("{}"), paths: []string{"user.json"}, provider: "claude", source: "user"},
		{name: "project replaces user", user: user, project: []byte("pre_pr_review:\n  data_paths: [rules/a.json, catalogs/, 'assets/*.json', 'literal/[a].json']\n"), paths: []string{"rules/a.json", "catalogs/", "assets/*.json", "literal/[a].json"}, provider: "claude", source: "user"},
		{name: "empty clears user", user: user, project: []byte("pre_pr_review:\n  data_paths: []\n"), paths: []string{}, provider: "claude", source: "user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := ResolveConfigProposal(tc.user, tc.project)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg.PrePRReview.DataPaths, tc.paths) || cfg.PrePRReview.Provider != tc.provider || cfg.PrePRReview.Source != tc.source {
				t.Fatalf("review config=%+v, want paths=%v provider=%q source=%q", cfg.PrePRReview, tc.paths, tc.provider, tc.source)
			}
			if len(tc.paths) == 0 && cfg.PrePRReview.ReviewsApart("user.json") {
				t.Fatal("absent or cleared data paths review a file apart")
			}
		})
	}

	for _, tc := range []struct {
		name  string
		value string
		want  string
	}{
		{name: "a mapping", value: "{}", want: "pre_pr_review.data_paths must be a list of paths"},
		{name: "a number entry", value: "[42]", want: `pre_pr_review.data_paths has non-string path "42"`},
		{name: "an absolute path", value: "[/tmp/file]", want: `pre_pr_review.data_paths has unsafe path "/tmp/file"`},
		{name: "a parent segment", value: "[a/../file]", want: `pre_pr_review.data_paths has unsafe path "a/../file"`},
		{name: "a bad pattern", value: "['[']", want: `pre_pr_review.data_paths has invalid pattern "[": syntax error in pattern`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			content := []byte("pre_pr_review:\n  data_paths: " + tc.value + "\n")
			for _, scope := range []struct {
				name    string
				user    []byte
				project []byte
			}{
				{name: "User Config proposal", user: content},
				{name: "Project Config proposal", project: content},
			} {
				t.Run(scope.name, func(t *testing.T) {
					t.Parallel()
					_, err := ResolveConfigProposal(scope.user, scope.project)
					want := `parse config "` + scope.name + `": ` + tc.want
					if err == nil || err.Error() != want {
						t.Fatalf("error=%v, want %q", err, want)
					}
				})
			}
		})
	}

	t.Run("matches whole paths", func(t *testing.T) {
		t.Parallel()
		cfg, err := ResolveConfigProposal(nil, []byte("pre_pr_review:\n  data_paths: [rules/a.json, catalogs/, 'assets/*.json', 'literal/[a].json']\n"))
		if err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string]bool{
			"rules/a.json": true, "rules/b.json": false,
			"catalogs/nested/file": true, "catalogsish/file": false, "catalogs": false,
			"assets/a.json": true, "assets/nested/a.json": false,
			"literal/[a].json": true, "literal/a.json": true,
			"src/a.go": false,
		} {
			if got := cfg.PrePRReview.ReviewsApart(name); got != want {
				t.Errorf("ReviewsApart(%q)=%v, want %v", name, got, want)
			}
			if cfg.PrePRReview.Omits(name) {
				t.Errorf("data_paths alone omits %q", name)
			}
		}
	})
}
