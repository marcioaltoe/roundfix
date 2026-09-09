// Suite: static Task Verification classification
// Invariant: Authored commands are classified from their declared shell form, independently of prose or status.
// Boundary IN: parsed Task Verification commands and the public detector
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const verificationSourceSlug = "committed-source"

func TestVerificationRefusesUntrustedSource(t *testing.T) {
	t.Parallel()

	t.Run("Spec Root outside repository Git tree", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		external := newExternalVerificationSource(t, "test -f expected.txt")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), speccheck.AuthoredCommandSourceRequest{
			RepoRoot:  fixture.repoRoot,
			SpecsRoot: external.specsRoot,
			SpecSlug:  verificationSourceSlug,
			Artifact:  external.artifact,
			Commands:  []string{"test -f expected.txt"},
		})

		assertSourceRefusal(t, err, speccheck.SourceConditionOutOfTree)
	})

	t.Run("untracked artifact", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		artifact := filepath.Join(verificationSourceSlug, "untracked.md")
		writeVerificationSourceFile(t, filepath.Join(fixture.specsRoot, artifact), "pending", "test -f expected.txt", "")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(artifact, "test -f expected.txt"))

		assertSourceRefusal(t, err, speccheck.SourceConditionUntrackedArtifact)
	})

	t.Run("modified authored requirement", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		writeVerificationSourceFile(t, fixture.artifactPath(), "pending", "test -f expected.txt", "2. MUST change the declared boundary.\n")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt"))

		assertSourceRefusal(t, err, speccheck.SourceConditionModifiedArtifact)
	})

	t.Run("altered command text", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		writeVerificationSourceFile(t, fixture.artifactPath(), "pending", "touch should-not-run", "")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "touch should-not-run"))

		assertSourceRefusal(t, err, speccheck.SourceConditionModifiedArtifact)
	})

	t.Run("requested command differs from carrying artifact", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "touch should-not-run"))

		assertSourceRefusal(t, err, speccheck.SourceConditionModifiedArtifact)
	})

	t.Run("modified source cannot append its own approval", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		changed := "touch should-not-run"
		approval := executionApprovalFixture{
			repository: fixture.repoRoot,
			revision:   fixture.artifactRevision(t),
			artifact:   fixture.repoArtifact(),
			command:    changed,
		}
		writeVerificationSourceFile(t, fixture.artifactPath(), "pending", changed, "")
		writeVerificationAuthorization(t, fixture.authorizationPath(), []executionApprovalFixture{approval})

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, changed))

		assertSourceRefusal(t, err, speccheck.SourceConditionModifiedArtifact)
	})

	t.Run("out-of-tree source cannot carry its own approval", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		external := newExternalVerificationSource(t, "test -f external.txt")
		writeVerificationAuthorization(t, filepath.Join(external.specsRoot, verificationSourceSlug, "_authorization.md"), []executionApprovalFixture{{
			repository: external.repoRoot,
			revision:   external.artifactRevision(t),
			artifact:   external.repoArtifact(),
			command:    "test -f external.txt",
		}})
		gittest.Run(t, external.repoRoot, "add", "-A")
		gittest.Run(t, external.repoRoot, "commit", "-m", "append self approval")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), sourceRequest(fixture.repoRoot, external, "test -f external.txt"))

		assertSourceRefusal(t, err, speccheck.SourceConditionOutOfTree)
	})

	t.Run("noncanonical approved artifact grants nothing", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		external := newExternalVerificationSource(t, "test -f external.txt")
		fixture.commitExecutionApprovals(t, []executionApprovalFixture{{
			repository: external.repoRoot,
			revision:   external.artifactRevision(t),
			artifact:   "specs/" + verificationSourceSlug + "/../" + verificationSourceSlug + "/task_01.md",
			command:    "test -f external.txt",
		}})

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), sourceRequest(fixture.repoRoot, external, "test -f external.txt"))

		assertSourceRefusal(t, err, speccheck.SourceConditionOutOfTree)
	})

	t.Run("unavailable revision", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		request := fixture.request(fixture.artifact, "test -f expected.txt")
		request.Revision = "revision-that-does-not-exist"

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), request)

		assertSourceRefusal(t, err, speccheck.SourceConditionUnresolved)
	})

	t.Run("unavailable approved source revision", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		changed := "touch should-not-run"
		fixture.commitExecutionApprovals(t, []executionApprovalFixture{{
			repository: fixture.repoRoot,
			revision:   "approved-revision-that-does-not-exist",
			artifact:   fixture.repoArtifact(),
			command:    changed,
		}})
		writeVerificationSourceFile(t, fixture.artifactPath(), "pending", changed, "")

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, changed))

		assertSourceRefusal(t, err, speccheck.SourceConditionUnresolved)
	})

	t.Run("unreadable Git object", func(t *testing.T) {
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		object := strings.TrimSpace(gittest.Run(t, fixture.repoRoot, "rev-parse", "HEAD:"+filepath.ToSlash(filepath.Join("docs", "specs", fixture.artifact))))
		objectPath := filepath.Join(fixture.repoRoot, ".git", "objects", object[:2], object[2:])
		if err := os.Remove(objectPath); err != nil {
			t.Fatalf("remove fixture blob %s: %v", object, err)
		}

		err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt"))

		assertSourceRefusal(t, err, speccheck.SourceConditionUnresolved)
	})
}

func TestVerificationExecutesOnCommittedProvenance(t *testing.T) {
	t.Parallel()

	t.Run("tracked unmodified artifact is authorized", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")

		if err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt")); err != nil {
			t.Fatalf("AuthorizeAuthoredCommands() error = %v", err)
		}
	})

	t.Run("committed approval authorizes modified in-tree artifact", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		approvedCommand := "test -f approved.txt"
		fixture.commitExecutionApprovals(t, []executionApprovalFixture{{
			repository: fixture.repoRoot,
			revision:   fixture.artifactRevision(t),
			artifact:   fixture.repoArtifact(),
			command:    approvedCommand,
		}})
		writeVerificationSourceFile(t, fixture.artifactPath(), "pending", approvedCommand, "")
		before := mustReadVerificationSource(t, fixture.authorizationPath())

		if err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, approvedCommand)); err != nil {
			t.Fatalf("AuthorizeAuthoredCommands() error = %v", err)
		}
		if after := mustReadVerificationSource(t, fixture.authorizationPath()); after != before {
			t.Fatal("authorization check changed the historical approval record")
		}
	})

	t.Run("committed approval authorizes out-of-tree artifact", func(t *testing.T) {
		t.Parallel()
		fixture := newVerificationSourceFixture(t, "test -f expected.txt")
		external := newExternalVerificationSource(t, "test -f external.txt")
		fixture.commitExecutionApprovals(t, []executionApprovalFixture{{
			repository: external.repoRoot,
			revision:   external.artifactRevision(t),
			artifact:   external.repoArtifact(),
			command:    "test -f external.txt",
		}})

		if err := speccheck.AuthorizeAuthoredCommands(context.Background(), speccheck.AuthoredCommandSourceRequest{
			RepoRoot:  fixture.repoRoot,
			SpecsRoot: external.specsRoot,
			SpecSlug:  verificationSourceSlug,
			Artifact:  external.artifact,
			Commands:  []string{"test -f external.txt"},
		}); err != nil {
			t.Fatalf("AuthorizeAuthoredCommands() error = %v", err)
		}
	})

	t.Run("approval expires with every bound identity field", func(t *testing.T) {
		approvedCommand := "test -f external.txt"
		tests := []struct {
			name   string
			mutate func(t *testing.T, fixture verificationSourceFixture, external verificationSourceFixture) speccheck.AuthoredCommandSourceRequest
		}{
			{
				name: "repository",
				mutate: func(t *testing.T, fixture verificationSourceFixture, _ verificationSourceFixture) speccheck.AuthoredCommandSourceRequest {
					other := newExternalVerificationSource(t, approvedCommand)
					return sourceRequest(fixture.repoRoot, other, approvedCommand)
				},
			},
			{
				name: "revision",
				mutate: func(t *testing.T, fixture verificationSourceFixture, external verificationSourceFixture) speccheck.AuthoredCommandSourceRequest {
					writeVerificationSourceFile(t, external.artifactPath(), "pending", approvedCommand, "2. MUST bind a different revision.\n")
					gittest.Run(t, external.repoRoot, "add", "-A")
					gittest.Run(t, external.repoRoot, "commit", "-m", "change source revision")
					return sourceRequest(fixture.repoRoot, external, approvedCommand)
				},
			},
			{
				name: "artifact",
				mutate: func(t *testing.T, fixture verificationSourceFixture, external verificationSourceFixture) speccheck.AuthoredCommandSourceRequest {
					external.artifact = filepath.Join(verificationSourceSlug, "task_02.md")
					writeVerificationSourceFile(t, external.artifactPath(), "pending", approvedCommand, "")
					gittest.Run(t, external.repoRoot, "add", "-A")
					gittest.Run(t, external.repoRoot, "commit", "-m", "add different artifact")
					return sourceRequest(fixture.repoRoot, external, approvedCommand)
				},
			},
			{
				name: "command digest",
				mutate: func(t *testing.T, fixture verificationSourceFixture, external verificationSourceFixture) speccheck.AuthoredCommandSourceRequest {
					changed := "touch should-not-run"
					writeVerificationSourceFile(t, external.artifactPath(), "pending", changed, "")
					return sourceRequest(fixture.repoRoot, external, changed)
				},
			},
		}

		for _, testCase := range tests {
			testCase := testCase
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()
				fixture := newVerificationSourceFixture(t, "test -f expected.txt")
				external := newExternalVerificationSource(t, approvedCommand)
				fixture.commitExecutionApprovals(t, []executionApprovalFixture{{
					repository: external.repoRoot,
					revision:   external.artifactRevision(t),
					artifact:   external.repoArtifact(),
					command:    approvedCommand,
				}})
				before := mustReadVerificationSource(t, fixture.authorizationPath())

				err := speccheck.AuthorizeAuthoredCommands(context.Background(), testCase.mutate(t, fixture, external))

				assertSourceRefusal(t, err, speccheck.SourceConditionOutOfTree)
				if after := mustReadVerificationSource(t, fixture.authorizationPath()); after != before {
					t.Fatal("expired approval changed the historical authorization record")
				}
			})
		}
	})
}

func TestAuthoredProjectionIgnoresDaemonOwnedFields(t *testing.T) {
	t.Parallel()
	fixture := newVerificationSourceFixture(t, "test -f expected.txt")

	writeVerificationSourceFile(t, fixture.artifactPath(), "in_progress", "test -f expected.txt", "")
	if err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt")); err != nil {
		t.Fatalf("in_progress status was not excluded from authored projection: %v", err)
	}
	writeVerificationSourceFile(t, fixture.artifactPath(), "in_progress", "test -f expected.txt", "## Result\n\n- Agent implementation record.\n")
	if err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt")); err != nil {
		t.Fatalf("Result was not excluded from authored projection: %v", err)
	}
}

func TestExecutionEntryPointCharacterization(t *testing.T) {
	t.Parallel()
	for _, entryPoint := range []string{"probe", "post-Agent verification", "Settle"} {
		entryPoint := entryPoint
		t.Run(entryPoint, func(t *testing.T) {
			t.Parallel()
			fixture := newVerificationSourceFixture(t, "test -f expected.txt")
			if err := speccheck.AuthorizeAuthoredCommands(context.Background(), fixture.request(fixture.artifact, "test -f expected.txt")); err != nil {
				t.Fatalf("committed source at %s changed today's executable outcome: %v", entryPoint, err)
			}
		})
	}
}

type executionApprovalFixture struct {
	repository string
	revision   string
	artifact   string
	command    string
}

type verificationSourceFixture struct {
	repoRoot  string
	specsRoot string
	artifact  string
}

func newVerificationSourceFixture(t *testing.T, command string) verificationSourceFixture {
	t.Helper()
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	fixture := verificationSourceFixture{
		repoRoot:  repoRoot,
		specsRoot: specsRoot,
		artifact:  filepath.Join(verificationSourceSlug, "task_01.md"),
	}
	writeVerificationSourceFile(t, fixture.artifactPath(), "pending", command, "")
	writeVerificationAuthorization(t, fixture.authorizationPath(), nil)
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	gittest.Run(t, repoRoot, "add", "-A")
	gittest.Run(t, repoRoot, "commit", "-m", "seed committed Verification source")
	return fixture
}

func newExternalVerificationSource(t *testing.T, command string) verificationSourceFixture {
	t.Helper()
	repoRoot := t.TempDir()
	fixture := verificationSourceFixture{
		repoRoot:  repoRoot,
		specsRoot: filepath.Join(repoRoot, "specs"),
		artifact:  filepath.Join(verificationSourceSlug, "task_01.md"),
	}
	writeVerificationSourceFile(t, fixture.artifactPath(), "pending", command, "")
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	gittest.Run(t, repoRoot, "add", "-A")
	gittest.Run(t, repoRoot, "commit", "-m", "seed external Verification source")
	return fixture
}

func (fixture verificationSourceFixture) request(artifact string, commands ...string) speccheck.AuthoredCommandSourceRequest {
	return speccheck.AuthoredCommandSourceRequest{
		RepoRoot:  fixture.repoRoot,
		SpecsRoot: fixture.specsRoot,
		SpecSlug:  verificationSourceSlug,
		Artifact:  artifact,
		Commands:  commands,
	}
}

func sourceRequest(repoRoot string, source verificationSourceFixture, command string) speccheck.AuthoredCommandSourceRequest {
	return speccheck.AuthoredCommandSourceRequest{
		RepoRoot:  repoRoot,
		SpecsRoot: source.specsRoot,
		SpecSlug:  verificationSourceSlug,
		Artifact:  source.artifact,
		Commands:  []string{command},
	}
}

func (fixture verificationSourceFixture) artifactPath() string {
	return filepath.Join(fixture.specsRoot, fixture.artifact)
}

func (fixture verificationSourceFixture) repoArtifact() string {
	relative, err := filepath.Rel(fixture.repoRoot, fixture.artifactPath())
	if err != nil {
		panic(err)
	}
	return filepath.ToSlash(relative)
}

func (fixture verificationSourceFixture) artifactRevision(t *testing.T) string {
	t.Helper()
	return strings.TrimSpace(gittest.Run(t, fixture.repoRoot, "log", "-1", "--format=%H", "HEAD", "--", fixture.repoArtifact()))
}

func (fixture verificationSourceFixture) authorizationPath() string {
	return filepath.Join(fixture.repoRoot, "docs", "specs", verificationSourceSlug, "_authorization.md")
}

func (fixture verificationSourceFixture) commitExecutionApprovals(t *testing.T, approvals []executionApprovalFixture) {
	t.Helper()
	writeVerificationAuthorization(t, fixture.authorizationPath(), approvals)
	gittest.Run(t, fixture.repoRoot, "add", filepath.ToSlash(filepath.Join("docs", "specs", verificationSourceSlug, "_authorization.md")))
	gittest.Run(t, fixture.repoRoot, "commit", "-m", "approve exact Verification source")
}

func writeVerificationSourceFile(t *testing.T, path string, status string, command string, extra string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create Verification source directory: %v", err)
	}
	requirementExtra := extra
	resultExtra := ""
	if strings.HasPrefix(extra, "## Result") {
		requirementExtra = ""
		resultExtra = extra
	}
	content := "---\n" +
		"task: task_01\n" +
		"spec: " + verificationSourceSlug + "\n" +
		"status: " + status + "\n" +
		"type: backend\n" +
		"complexity: low\n" +
		"---\n\n" +
		"# Task 01: Exercise committed source\n\n" +
		"## Requirements\n\n" +
		"1. MUST keep the authored command bounded.\n" +
		requirementExtra + "\n" +
		"## Verification\n\n" +
		"- `" + command + "`\n"
	if resultExtra != "" {
		content = strings.TrimSuffix(content, "\n") + "\n\n" + resultExtra
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write Verification source: %v", err)
	}
}

func writeVerificationAuthorization(t *testing.T, path string, approvals []executionApprovalFixture) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create authorization directory: %v", err)
	}
	var content strings.Builder
	content.WriteString("---\nstatus: approved\ngranted: 2026-09-09\naction: approve exact Verification source\nconsuming: " + verificationSourceSlug + "\npaths:\n  - docs/specs/" + verificationSourceSlug + "/task_01.md\noperations:\n  - implement\n")
	if len(approvals) > 0 {
		content.WriteString("execution_approvals:\n")
		for _, approval := range approvals {
			content.WriteString("  - repository: " + approval.repository + "\n")
			content.WriteString("    revision: " + approval.revision + "\n")
			content.WriteString("    artifact: " + approval.artifact + "\n")
			content.WriteString("    command_digest: " + speccheck.AuthoredCommandDigest(approval.command) + "\n")
		}
	}
	content.WriteString("---\n\n# Approved Verification source\n")
	if err := os.WriteFile(path, []byte(content.String()), 0o644); err != nil {
		t.Fatalf("write Verification authorization: %v", err)
	}
}

func mustReadVerificationSource(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertSourceRefusal(t *testing.T, err error, condition speccheck.SourceCondition) {
	t.Helper()
	var refusal *speccheck.SourceUntrustedError
	if !errors.As(err, &refusal) {
		t.Fatalf("error = %v, want SourceUntrustedError", err)
	}
	if refusal.Code != speccheck.CodeSourceUntrusted || refusal.Condition != condition {
		t.Fatalf("refusal = %#v, want %s/%s", refusal, speccheck.CodeSourceUntrusted, condition)
	}
	for _, token := range []string{speccheck.CodeSourceUntrusted, string(condition)} {
		if !strings.Contains(err.Error(), token) {
			t.Errorf("error %q does not name %q", err, token)
		}
	}
}

func TestWorkIndependentVerificationRefusesOnlyWorkIndependentCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       spec.Status
		verification []string
		wantFinding  bool
	}{
		{
			name:         "repository gate and clean tree",
			status:       spec.StatusPending,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
		{
			name:         "repository wide Go gates and diff cleanliness",
			status:       spec.StatusInProgress,
			verification: []string{"go build -buildvcs=false ./...", "go test -parallel 16 ./...", "git diff --check"},
			wantFinding:  true,
		},
		{
			name:         "repository gate plus focused effect assertion",
			status:       spec.StatusPending,
			verification: []string{"make verify", "go test ./internal/speccheck -run '^TestTaskEffect$'"},
		},
		{
			name:         "repository wide command with effect selector",
			status:       spec.StatusPending,
			verification: []string{"go test ./... -run '^TestTaskEffect$'"},
		},
		{
			name:         "only effect assertions",
			status:       spec.StatusPending,
			verification: []string{"grep -q 'expected effect' internal/speccheck/verification.go", "go test ./internal/speccheck"},
		},
		{
			name:         "status does not change the declared command property",
			status:       spec.StatusCompleted,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			finding, got := speccheck.WorkIndependentVerification(spec.Task{
				File:         "fixture/task_01.md",
				Status:       tt.status,
				Verification: tt.verification,
			})
			if got != tt.wantFinding {
				t.Fatalf("WorkIndependentVerification() finding = %v, want %v; finding: %#v", got, tt.wantFinding, finding)
			}
			if !got {
				return
			}
			if finding.Code != speccheck.CodeVerifyWorkIndependent || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyWorkIndependent, speccheck.SeverityError)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" || finding.Where[0].Line < 1 {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
			if !strings.Contains(finding.Fix, "Task's own effect") {
				t.Errorf("finding fix = %q, want an effect-asserting repair", finding.Fix)
			}
		})
	}

	t.Run("task creates temporary artifact in an earlier command", func(t *testing.T) {
		t.Parallel()

		findings := speccheck.NonHermeticVerification(spec.Task{
			File: "fixture/task_01.md",
			Verification: []string{
				`go test ./internal/store > /tmp/task.log 2>&1`,
				`grep -q '^ok' /tmp/task.log`,
			},
		})
		if len(findings) != 0 {
			t.Fatalf("NonHermeticVerification() = %#v, want Task-created path permitted", findings)
		}
	})
}

func TestVerifyInvertedExit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		command         string
		wantForm        string
		wantReplacement string
	}{
		{
			name:            "grep count uses match status instead of the printed count",
			command:         "grep -c 'PASS' /tmp/task.log",
			wantForm:        "grep -c count-and-exit",
			wantReplacement: "grep -q",
		},
		{
			name:            "filtered count uses wc pipeline status",
			command:         "grep -v '^PASS' /tmp/task.log | wc -l",
			wantForm:        "grep -v ... | wc -l filtered count",
			wantReplacement: `test "$(grep -v ... | wc -l)" -eq 0`,
		},
		{
			name:            "test command substitution has no comparison",
			command:         `test "$(grep -v '^PASS' /tmp/task.log)"`,
			wantForm:        "test $(...) without a comparison",
			wantReplacement: `test -z "$(cmd)"`,
		},
		{
			name:    "grep q uses match status intentionally",
			command: "grep -q 'PASS' /tmp/task.log",
		},
		{
			name:    "test command substitution compares the count",
			command: `test "$(grep -c 'PASS' /tmp/task.log)" -ge 4`,
		},
		{
			name:    "redirected grep count is compared by the terminal test",
			command: `grep -c 'PASS' /tmp/task.log > /tmp/pass-count; test "$(cat /tmp/pass-count)" -ge 4`,
		},
		{
			name:    "filtered count is compared by the enclosing test",
			command: `test "$(grep -v '^PASS' /tmp/task.log | wc -l)" -eq 0`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			findings := speccheck.InvertedExitVerification(spec.Task{
				File:         "fixture/task_01.md",
				Verification: []string{tt.command},
			})
			if tt.wantForm == "" {
				if len(findings) != 0 {
					t.Fatalf("InvertedExitVerification(%q) = %#v, want no finding", tt.command, findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("InvertedExitVerification(%q) = %#v, want one finding", tt.command, findings)
			}
			finding := findings[0]
			if finding.Code != speccheck.CodeVerifyInvertedExit || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyInvertedExit, speccheck.SeverityError)
			}
			if !strings.Contains(finding.Summary, tt.wantForm) {
				t.Errorf("finding summary = %q, want form %q", finding.Summary, tt.wantForm)
			}
			if !strings.Contains(finding.Fix, tt.wantReplacement) {
				t.Errorf("finding fix = %q, want replacement %q", finding.Fix, tt.wantReplacement)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
		})
	}
}

func TestVerifyInvertedExitSkipsWithoutTaskGraph(t *testing.T) {
	t.Parallel()

	result, err := speccheck.Check(fixtureSpecRoot, "testdata/repo", "no-taskgraph")
	if err != nil {
		t.Fatalf("Check(no-taskgraph): %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeVerifyInvertedExit); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want detector skipped", speccheck.CodeVerifyInvertedExit, findings)
	}
	if !hasSkip(result, speccheck.CodeVerifyInvertedExit, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeVerifyInvertedExit)
	}
}

func TestVerifyNonHermetic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		command  string
		wantForm string
	}{
		{
			name:     "undeclared environment variable",
			command:  `go test ./internal/store -run "$TASK_04_BASELINE"`,
			wantForm: "undeclared environment variable TASK_04_BASELINE",
		},
		{
			name:     "environment presence guard",
			command:  `test -n "$VORTEX_POSTGRES_INTEGRATION_URL" && go test ./internal/store`,
			wantForm: "environment-presence guard VORTEX_POSTGRES_INTEGRATION_URL",
		},
		{
			name:     "external temporary tree snapshot",
			command:  `diff /tmp/task-before.tree current.tree`,
			wantForm: "external path /tmp/task-before.tree",
		},
		{
			name:    "redirect creates temporary artifact before reading it",
			command: `go test ./internal/store -run '^TestStore$' -v > /tmp/task.log 2>&1; s=$?; grep -q '^--- PASS: TestStore' /tmp/task.log || { cat /tmp/task.log; exit 1; }; exit $s`,
		},
		{
			name:    "command declares its own shell variable",
			command: `expected=PASS; test "$expected" = PASS`,
		},
		{
			name:    "single quoted variable text is literal",
			command: `grep -q '$DOCUMENTED_VARIABLE' docs/contract.md`,
		},
		{
			name:    "read declares a command local variable",
			command: `printf 4 | { read n; test "$n" -ge 4; }`,
		},
		{
			name:    "for declares its iteration variable",
			command: `for form in alpha beta; do test -n "$form"; done`,
		},
		{
			name:    "command creates external directory before reading it",
			command: `mkdir -p /tmp/task-fixture; test -d /tmp/task-fixture`,
		},
		{
			name:    "output flag creates executable before running it",
			command: `go build -o /tmp/task-binary ./cmd/roundfix && /tmp/task-binary version`,
		},
		{
			name:     "external path is read before a later write",
			command:  `cat /tmp/task.log; go test ./internal/store > /tmp/task.log`,
			wantForm: "external path /tmp/task.log",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			findings := speccheck.NonHermeticVerification(spec.Task{
				File:         "fixture/task_01.md",
				Verification: []string{tt.command},
			})
			if tt.wantForm == "" {
				if len(findings) != 0 {
					t.Fatalf("NonHermeticVerification(%q) = %#v, want no finding", tt.command, findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("NonHermeticVerification(%q) = %#v, want one finding", tt.command, findings)
			}
			finding := findings[0]
			if finding.Code != speccheck.CodeVerifyNonHermetic || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyNonHermetic, speccheck.SeverityError)
			}
			if !strings.Contains(finding.Summary, tt.wantForm) {
				t.Errorf("finding summary = %q, want form %q", finding.Summary, tt.wantForm)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
		})
	}
}

func TestVerifyNonHermeticRegistersAtTasksStage(t *testing.T) {
	t.Parallel()

	const slug = "non-hermetic-stage"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: non-hermetic-stage
status: active
created: 2026-08-14
surfaces: [backend]
---

# Non-hermetic stage fixture
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: non-hermetic-stage
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/task_01.md", `---
task: task_01
spec: non-hermetic-stage
status: pending
type: backend
complexity: low
---

# Task 01: Refuse external state

## Verification

- `+"`test -n \"$UNDECLARED_FIXTURE\" && go test ./internal/store`"+`
`)

	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeVerifyNonHermetic)
	if len(findings) != 1 {
		t.Fatalf("StageTasks %s findings = %#v, want one", speccheck.CodeVerifyNonHermetic, findings)
	}
	if !strings.Contains(findings[0].Summary, "environment-presence guard UNDECLARED_FIXTURE") {
		t.Errorf("finding summary = %q, want named guard", findings[0].Summary)
	}
}

func TestNonHermeticVerificationSkipsRoundfixOwnedQAVerification(t *testing.T) {
	t.Parallel()

	task := spec.Task{
		File:         "derived-qa-verification/task_01.md",
		Type:         spec.TaskTypeQA,
		Verification: spec.DerivedQAVerification("derived-qa-verification"),
	}
	if findings := speccheck.NonHermeticVerification(task); len(findings) != 0 {
		t.Fatalf("NonHermeticVerification() findings = %#v, want Roundfix-owned Verification skipped", findings)
	}
}

func TestVerifyNonHermeticSkipsWithoutTaskGraph(t *testing.T) {
	t.Parallel()

	result, err := speccheck.Check(fixtureSpecRoot, "testdata/repo", "no-taskgraph")
	if err != nil {
		t.Fatalf("Check(no-taskgraph): %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeVerifyNonHermetic); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want detector skipped", speccheck.CodeVerifyNonHermetic, findings)
	}
	if !hasSkip(result, speccheck.CodeVerifyNonHermetic, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeVerifyNonHermetic)
	}
}

func TestVacuousVerificationCommandIsCaughtBesideHonestSiblings(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name    string
		command string
		vacuous bool
	}{
		{
			name:    "an assertion over which paths changed",
			command: "git diff --name-only HEAD | grep -v -E '^(kept)$' | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			name:    "a porcelain cleanliness check",
			command: "git status --porcelain | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			// Verified on a clean tree on 2026-08-10: git writes nothing, grep
			// matches nothing, and the command exits 1. It fails on an
			// unchanged tree rather than passing, so it is not vacuous.
			name:    "a nonempty-diff assertion fails on an unchanged tree",
			command: "git diff --name-status HEAD | grep -q .",
			vacuous: false,
		},
		{
			name:    "the same paths with a predicate that passes when empty",
			command: "git diff --name-status HEAD | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			// The equality predicate compares two snapshots captured by earlier
			// chain segments. An intervening command (`make regen`) can change
			// the second snapshot, so `[ "$s1" = "$s2" ]` is not a guaranteed
			// success: its operands are variables, not statically proven
			// values. The matcher must not report it vacuous.
			name:    "a two-snapshot comparison with an intervening command",
			command: `s1="$(git status --porcelain)"; make regen >/dev/null 2>&1; s2="$(git status --porcelain)"; [ "$s1" = "$s2" ]`,
			vacuous: false,
		},
		{
			name:    "a named test that does not exist yet",
			command: "go test ./internal/agent -run '^TestX$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestX'",
			vacuous: false,
		},
		{
			name:    "a grep over declared source",
			command: "grep -c 'Declared break' internal/agent/selection_test.go",
			vacuous: false,
		},
		{
			// git --exit-code does not make the terminal `grep -q .` pass on an
			// unchanged tree; grep still receives no input and exits 1.
			name:    "an exit-code diff piped into grep still fails when empty",
			command: "git diff --exit-code HEAD | grep -q .",
			vacuous: false,
		},
		{
			// A terminal that consumes empty output successfully: cat exits 0
			// on an unchanged tree, so the command passes before any work.
			name:    "a diff piped into cat passes when empty",
			command: "git diff --name-only HEAD | cat",
			vacuous: true,
		},
		{
			// An inverted empty-output check: it asserts changes exist, so on
			// an unchanged tree the terminal `exit 1` fires and fails honestly.
			name:    "an inverted empty-output check fails when empty",
			command: "git diff --name-status HEAD | grep -q . && exit 1",
			vacuous: false,
		},
		{
			// A success form inside another command's quoted argument is not a
			// terminal predicate: grep finds no lines in an unchanged tree's
			// empty output and exits 1, whatever text it was asked to match.
			name:    "a quoted success form inside grep still fails when empty",
			command: "git diff --name-only HEAD | grep -q 'exit 0'",
			vacuous: false,
		},
		{
			// `&&` short-circuits: grep -q . fails on an unchanged tree, so
			// cat is skipped and never decides the exit status, which stays 1.
			name:    "an and-chain whose tail is skipped when empty fails honestly",
			command: "git diff --name-only HEAD | grep -q . && cat",
			vacuous: false,
		},
		{
			// `||` runs its right side when the left fails, so the terminal cat
			// does execute on an unchanged tree and is genuinely vacuous.
			name:    "an or-chain whose tail runs when the left fails passes when empty",
			command: "git diff --name-only HEAD | grep -q . || cat",
			vacuous: true,
		},
		{
			// A quoted test -z form is likewise a grep of ordinary text, not an
			// empty-string test, so it fails on an unchanged tree.
			name:    "a quoted empty-test inside grep still fails when empty",
			command: `git diff --name-only HEAD | grep -q 'test -z'`,
			vacuous: false,
		},
		{
			// `test -z "$(printf x)"` reads a substitution of fixed non-empty
			// output, not the working tree, so it fails on an unchanged tree.
			name:    "an empty-test over non-empty command substitution fails",
			command: `git diff --name-only HEAD | test -z "$(printf x)"`,
			vacuous: false,
		},
		{
			// Single quotes make the operand literal text, not a substitution;
			// `'$(printf x)'` is a fixed non-empty string that fails the test.
			name:    "a single-quoted command substitution is literal and fails",
			command: `git diff --name-only HEAD | test -z '$(printf x)'`,
			vacuous: false,
		},
		{
			// A working-tree substitution is genuinely vacuous: an unchanged
			// tree yields an empty value from it and the test passes.
			name:    "an empty-test over a working-tree substitution passes when empty",
			command: `git diff --name-only HEAD | test -z "$(git diff --name-only HEAD)"`,
			vacuous: true,
		},
		{
			// The bracket form behaves like the test form.
			name:    "a bracket empty-test over a working-tree substitution passes when empty",
			command: `git diff --name-only HEAD | [ -z "$(git diff --name-only HEAD)" ]`,
			vacuous: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			task := spec.Task{Verification: []string{testCase.command}}
			vacuous := speccheck.VacuousVerificationCommands(task)
			if (len(vacuous) != 0) != testCase.vacuous {
				t.Fatalf("VacuousVerificationCommands(%q) = %#v, want vacuous %v", testCase.command, vacuous, testCase.vacuous)
			}
		})
	}
}

// TestOneHonestCommandDoesNotAbsolveAVacuousSibling pins the unit of judgement:
// the Daemon's pre-work probe refuses a Task for any single vacuous command, so
// the static check must not stay silent because a sibling is honest.
func TestOneHonestCommandDoesNotAbsolveAVacuousSibling(t *testing.T) {
	t.Parallel()

	task := spec.Task{Verification: []string{
		"go test ./internal/agent -run '^TestY$' -count=1 -v 2>&1 | grep -q '^--- PASS: TestY'",
		"test -z \"$(git diff --name-only HEAD)\"",
	}}
	vacuous := speccheck.VacuousVerificationCommands(task)
	if len(vacuous) != 1 || !strings.Contains(vacuous[0], "--name-only") {
		t.Fatalf("VacuousVerificationCommands() = %#v, want only the working-tree assertion", vacuous)
	}
	if _, reported := speccheck.WorkIndependentVerification(task); reported {
		t.Fatal("WorkIndependentVerification judges the whole Task and must stay silent here")
	}
}
