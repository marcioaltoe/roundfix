// Suite: pre-PR review records and reviewer input.
// Invariant: each review is bound to one candidate whose diff is supplied to a read-only Agent Session.
// Boundary IN: review-record encoding, candidate diff collection, prompt construction, and runner requests.
// Boundary OUT: policy resolution, command routing, output classification, and record storage.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
)

func TestReviewRecordRoundTripsEachOutcome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		outcome  reviewOutcome
		findings string
		reason   string
	}{
		{
			name:    "reviewed",
			outcome: reviewOutcomeReviewed,
		},
		{
			name:     "findings",
			outcome:  reviewOutcomeFindings,
			findings: "internal/cli/review.go: head validation is missing",
		},
		{
			name:    "blocked",
			outcome: reviewOutcomeBlocked,
			reason:  "reviewer output was empty",
		},
		{
			name:    "omitted",
			outcome: reviewOutcomeOmitted,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			policy := roundconfig.PrePRReview{
				Provider: "codex",
				Source:   "project",
			}
			want := newReviewRecord(
				"github.com/example/project",
				"1111111111111111111111111111111111111111",
				"2222222222222222222222222222222222222222",
				policy,
				test.outcome,
			)
			want.Findings = test.findings
			want.Reason = test.reason

			var output bytes.Buffer
			if err := writeReviewRecord(&output, want); err != nil {
				t.Fatalf("write review record: %v", err)
			}

			var got reviewRecord
			if err := json.Unmarshal(output.Bytes(), &got); err != nil {
				t.Fatalf("decode review record: %v", err)
			}
			if got != want {
				t.Fatalf("review record = %+v, want %+v", got, want)
			}
		})
	}
}

func TestReviewRecordWithoutHeadIsRefusedBeforeWriting(t *testing.T) {
	t.Parallel()

	record := newReviewRecord(
		"github.com/example/project",
		"1111111111111111111111111111111111111111",
		"",
		roundconfig.PrePRReview{Provider: "codex", Source: "project"},
		reviewOutcomeReviewed,
	)
	var output bytes.Buffer

	err := writeReviewRecord(&output, record)

	if err == nil || !strings.Contains(err.Error(), "head commit is required") {
		t.Fatalf("write review record error = %v, want missing-head error", err)
	}
	if output.Len() != 0 {
		t.Fatalf("invalid review record wrote %q", output.String())
	}
}

func TestReviewRecordRefusesInvalidOutcomeDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		outcome  reviewOutcome
		findings string
		reason   string
		want     string
	}{
		{
			name:    "unknown outcome",
			outcome: "unknown",
			want:    "outcome \"unknown\" is invalid",
		},
		{
			name:    "findings outcome without findings",
			outcome: reviewOutcomeFindings,
			want:    "findings are required",
		},
		{
			name:    "blocked outcome without reason",
			outcome: reviewOutcomeBlocked,
			want:    "reason is required",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			record := newReviewRecord(
				"github.com/example/project",
				"1111111111111111111111111111111111111111",
				"2222222222222222222222222222222222222222",
				roundconfig.PrePRReview{Provider: "codex", Source: "project"},
				test.outcome,
			)
			record.Findings = test.findings
			record.Reason = test.reason

			var output bytes.Buffer
			err := writeReviewRecord(&output, record)

			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("write review record error = %v, want %q", err, test.want)
			}
			if output.Len() != 0 {
				t.Fatalf("invalid review record wrote %q", output.String())
			}
		})
	}
}

func TestReviewPromptCarriesTheCandidateDiff(t *testing.T) {
	t.Parallel()

	repository, baseCommit, headCommit := reviewCandidateFixture(t)
	runner := &recordingReviewRunner{}

	if _, err := runReviewSession(
		t.Context(),
		agent.ExecuteRequest{GitRoot: repository},
		baseCommit,
		headCommit,
		preflight.ExecGitRunner{},
		runner,
		runevent.Discard,
	); err != nil {
		t.Fatalf("run review session: %v", err)
	}

	for _, want := range []string{
		baseCommit,
		headCommit,
		"diff --git a/review.txt b/review.txt",
		"+the reviewer receives this changed line",
	} {
		if !strings.Contains(runner.request.Prompt, want) {
			t.Fatalf("review prompt does not contain %q:\n%s", want, runner.request.Prompt)
		}
	}
}

func TestReviewSessionReadsWithoutWriting(t *testing.T) {
	t.Parallel()

	repository, baseCommit, headCommit := reviewCandidateFixture(t)
	runner := &recordingReviewRunner{}

	if _, err := runReviewSession(
		t.Context(),
		agent.ExecuteRequest{GitRoot: repository},
		baseCommit,
		headCommit,
		preflight.ExecGitRunner{},
		runner,
		runevent.Discard,
	); err != nil {
		t.Fatalf("run review session: %v", err)
	}

	if !runner.request.Access.CanRead() {
		t.Fatal("review session cannot read files")
	}
	if runner.request.Access.CanWrite() {
		t.Fatal("review session can write")
	}
}

func reviewCandidateFixture(t *testing.T) (repository string, baseCommit string, headCommit string) {
	t.Helper()

	repository = t.TempDir()
	gittest.InitRepo(t, repository, "--initial-branch=main")
	gittest.PersistIdentity(t, repository)
	path := filepath.Join(repository, "review.txt")
	if err := os.WriteFile(path, []byte("before review\n"), 0o644); err != nil {
		t.Fatalf("write base candidate: %v", err)
	}
	gittest.Run(t, repository, "add", "review.txt")
	gittest.Run(t, repository, "commit", "-m", "base")
	baseCommit = strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))

	if err := os.WriteFile(path, []byte("before review\nthe reviewer receives this changed line\n"), 0o644); err != nil {
		t.Fatalf("write head candidate: %v", err)
	}
	gittest.Run(t, repository, "add", "review.txt")
	gittest.Run(t, repository, "commit", "-m", "candidate")
	headCommit = strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
	return repository, baseCommit, headCommit
}

type recordingReviewRunner struct {
	request agent.ExecuteRequest
}

func (*recordingReviewRunner) Probe(context.Context, agent.ProbeRequest) error {
	return nil
}

func (runner *recordingReviewRunner) Run(
	_ context.Context,
	request agent.ExecuteRequest,
	_ runevent.Sink,
) (agent.ExecuteResult, error) {
	runner.request = request
	return agent.ExecuteResult{Output: "No findings"}, nil
}

func (*recordingReviewRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}
