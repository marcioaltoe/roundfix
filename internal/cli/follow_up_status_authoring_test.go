// Suite: the status view's Authoring Question Preference.
// Boundary IN: private homes and repositories, local refs and committed config.
// Boundary OUT: the injected Git runner permits reads only; no network or remote.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/preflight"
)

func TestTheStatusViewReportsTheAuthoringQuestionPreference(t *testing.T) {
	t.Parallel()
	const checkoutPart = "; this checkout's Project Config gives follow, which counts only once the default branch holds it"
	const invalidDetail = "the default branch's Project Config was not read: read conflict declarations: parse config \"Project Config proposal\": authoring.open_questions must be wait or follow"
	for _, tc := range []struct {
		name string
		want followUpAuthoringQuestions
		line string
	}{
		{"default_without_the_key", followUpAuthoringQuestions{Preference: "wait", Source: "default", Ref: "refs/heads/main"}, "Open authoring questions: wait (default)"},
		{"follow_committed_on_main", followUpAuthoringQuestions{Preference: "follow", Source: "project-config", Ref: "refs/heads/main"}, "Open authoring questions: follow (Project Config)"},
		{"follow_only_in_the_working_tree", followUpAuthoringQuestions{Preference: "wait", Source: "default", Ref: "refs/heads/main", Checkout: "follow"}, "Open authoring questions: wait (default)" + checkoutPart},
		{"follow_only_on_a_feature_branch", followUpAuthoringQuestions{Preference: "wait", Source: "default", Ref: "refs/heads/main", Checkout: "follow"}, "Open authoring questions: wait (default)" + checkoutPart},
		{"remote_tracking_ref_wins", followUpAuthoringQuestions{Preference: "follow", Source: "project-config", Ref: "refs/remotes/origin/main"}, "Open authoring questions: follow (Project Config)"},
		{"no_default_branch_head", followUpAuthoringQuestions{Preference: "wait", Source: "default", Detail: "no default branch head was found"}, "Open authoring questions: wait (default); no default branch head was found"},
		{"unreadable_project_config", followUpAuthoringQuestions{Preference: "wait", Source: "default", Ref: "refs/heads/main", Detail: invalidDetail}, "Open authoring questions: wait (default); " + invalidDetail},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			now := time.Date(2026, 11, 19, 20, 0, 0, 0, time.UTC)
			f, _ := newStatusQueueFixture(t, now)
			gitImplement(t, f.repo, "remote", "remove", "origin")
			base := strings.TrimSpace(gitImplementOutput(t, f.repo, "rev-parse", "HEAD"))
			writeConfig := func(value string) {
				followUpWrite(t, f.repo, ".roundfixrc.yml", "authoring:\n  open_questions: "+value+"\n")
			}
			commitConfig := func(value string) {
				writeConfig(value)
				gitImplement(t, f.repo, "add", ".roundfixrc.yml")
				gitImplement(t, f.repo, "commit", "-m", "docs: declare authoring preference")
			}
			switch tc.name {
			case "follow_committed_on_main":
				commitConfig("follow")
			case "follow_only_in_the_working_tree":
				writeConfig("follow")
			case "follow_only_on_a_feature_branch":
				// A local symbolic default hint names main without a remote or a
				// fetched branch; the resolver falls back to the local main ref.
				gitImplement(t, f.repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
				gitImplement(t, f.repo, "checkout", "-b", "work")
				commitConfig("follow")
			case "remote_tracking_ref_wins":
				gitImplement(t, f.repo, "checkout", "-b", "work")
				commitConfig("follow")
				gitImplement(t, f.repo, "update-ref", "refs/remotes/origin/main", "HEAD")
				gitImplement(t, f.repo, "checkout", "main")
				// The checkout agrees with the fetched value, but local main
				// remains at base. This case isolates precedence of the refs.
				writeConfig("follow")
				if got := strings.TrimSpace(gitImplementOutput(t, f.repo, "rev-parse", "main")); got != base {
					t.Fatalf("local main moved: %s", got)
				}
			case "no_default_branch_head":
				gitImplement(t, f.repo, "branch", "-m", "work")
			case "unreadable_project_config":
				gitImplement(t, f.repo, "branch", "work")
				commitConfig("maybe")
				gitImplement(t, f.repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
				gitImplement(t, f.repo, "checkout", "work")
			}
			runner := followUpReadGit(func(ctx context.Context, dir string, args ...string) (string, error) {
				switch args[0] {
				case "rev-parse", "symbolic-ref", "for-each-ref", "ls-tree", "show", "log", "rev-list", "diff-tree", "hash-object":
					return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
				default:
					t.Errorf("unexpected Git command: %v", args)
					return "", fmt.Errorf("unexpected Git command: %v", args)
				}
			})
			for _, format := range []string{"text", "json"} {
				env := commandEnvironmentForTest(t)
				var out, stderr bytes.Buffer
				loaded, err := loadCommandConfig(env, &stderr)
				if err != nil {
					t.Fatal(err)
				}
				code := runFollowUpStatusWithHistoryGit(t.Context(), followUpRequest{status: true, format: format}, &out, &stderr, loaded, env, runner)
				if code != 0 || stderr.Len() != 0 {
					t.Fatalf("%s exit=%d stderr=%q stdout=%s", format, code, stderr.String(), out.String())
				}
				if format == "text" {
					var lines []string
					for _, line := range strings.Split(out.String(), "\n") {
						if strings.HasPrefix(line, "Open authoring questions:") {
							lines = append(lines, line)
						}
					}
					if !reflect.DeepEqual(lines, []string{tc.line}) || !strings.Contains(out.String(), tc.line+"\nDelegated decisions: 0\nSources:\n") {
						t.Errorf("text=%s\nwant line=%s", out.String(), tc.line)
					}
				} else {
					var doc followUpStatusDocument
					if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
						t.Fatal(err)
					}
					if doc.AuthoringQuestions != tc.want {
						t.Errorf("authoring_questions=%+v, want %+v", doc.AuthoringQuestions, tc.want)
					}
					var fields struct {
						AuthoringQuestions map[string]json.RawMessage `json:"authoring_questions"`
					}
					if err := json.Unmarshal(out.Bytes(), &fields); err != nil {
						t.Fatal(err)
					}
					if len(fields.AuthoringQuestions) != 5 {
						t.Errorf("JSON fields=%v", fields.AuthoringQuestions)
					}
				}
			}
		})
	}
}
