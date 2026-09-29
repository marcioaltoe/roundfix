// Suite: mechanical authorization commit identity
// Invariant: every authorization audit read and rendered row names the exact Task commit it audited.
// Boundary IN: public speccheck mechanical API, real temporary Git histories, and report rendering
// Boundary OUT: Daemon commit discovery and QA Agent command execution
package speccheck_test

import (
	"bytes"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestMechanicalAuthorizationReadNamesTheAuditedCommit(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const authorizationPath = "docs/specs/commit-column/_authorization.md"
	writeMechanicalFile(t, repoRoot, authorizationPath, mechanicalTypedAuthorization("commit-column", "Makefile"))
	target := commitMechanicalFiles(t, repoRoot, "record authorization", authorizationPath)
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
	taskCommit := commitMechanicalFiles(t, repoRoot, "change tooling", "Makefile")

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:               repoRoot,
		AuthorizationPath:      authorizationPath,
		ConsumingSpec:          "commit-column",
		DeliveryTargetRevision: target,
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID: "task_02",
			SHA:    taskCommit,
		}},
	})

	if len(result.AuthorizationReads) != 1 {
		t.Fatalf("AuthorizationReads = %#v, want one granted read", result.AuthorizationReads)
	}
	read := result.AuthorizationReads[0]
	if read.Outcome != spec.AuthorizationGranted {
		t.Fatalf("AuthorizationReads[0].Outcome = %q, want %q", read.Outcome, spec.AuthorizationGranted)
	}
	if read.Commit != taskCommit {
		t.Fatalf("AuthorizationReads[0].Commit = %q, want %q", read.Commit, taskCommit)
	}
}

func TestMechanicalUnresolvedAuthorizationReadNamesTheAuditedCommit(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const authorizationPath = "docs/specs/commit-column/_authorization.md"
	writeMechanicalFile(t, repoRoot, authorizationPath, mechanicalTypedAuthorization("commit-column", "Makefile"))
	target := commitMechanicalFiles(t, repoRoot, "record authorization", authorizationPath)
	missingTaskCommit := strings.Repeat("f", 40)

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:               repoRoot,
		AuthorizationPath:      authorizationPath,
		ConsumingSpec:          "commit-column",
		DeliveryTargetRevision: target,
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID: "task_02",
			SHA:    missingTaskCommit,
		}},
	})

	if len(result.AuthorizationReads) != 1 {
		t.Fatalf("AuthorizationReads = %#v, want one unresolved read", result.AuthorizationReads)
	}
	read := result.AuthorizationReads[0]
	if read.Outcome != spec.AuthorizationUnresolved {
		t.Fatalf("AuthorizationReads[0].Outcome = %q, want %q", read.Outcome, spec.AuthorizationUnresolved)
	}
	if read.Commit != missingTaskCommit {
		t.Fatalf("AuthorizationReads[0].Commit = %q, want %q", read.Commit, missingTaskCommit)
	}
}

func TestMechanicalReportAuditTableHasACommitColumn(t *testing.T) {
	t.Parallel()

	const taskCommit = "0123456789abcdef0123456789abcdef01234567"
	result := speccheck.MechanicalResult{
		AuthorizationReads: []speccheck.MechanicalAuthorizationRead{{
			TaskID:  "task_02",
			Commit:  taskCommit,
			Outcome: spec.AuthorizationGranted,
			Source: spec.AuthorizationSource{
				Path:     "docs/specs/commit-column/_authorization.md",
				Revision: "fedcba9876543210fedcba9876543210fedcba98",
			},
		}},
	}

	var report bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&report, result); err != nil {
		t.Fatalf("WriteMechanicalResult() error = %v", err)
	}
	if !strings.Contains(report.String(), "| Task | Commit | Outcome | Record | Revision | Detail |") {
		t.Fatalf("authorization audit table lacks Commit header:\n%s", report.String())
	}
	if !strings.Contains(report.String(), "| task_02 | "+taskCommit+" | granted |") {
		t.Fatalf("authorization audit row lacks full Task commit %q:\n%s", taskCommit, report.String())
	}
}

func TestMechanicalReportAuditTableListsEachCommitOfATask(t *testing.T) {
	t.Parallel()

	repoRoot := newMechanicalGitRepo(t)
	const authorizationPath = "docs/specs/commit-column/_authorization.md"
	writeMechanicalFile(t, repoRoot, authorizationPath, mechanicalTypedAuthorization("commit-column", "Makefile"))
	target := commitMechanicalFiles(t, repoRoot, "record authorization", authorizationPath)
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
	firstCommit := commitMechanicalFiles(t, repoRoot, "first tooling change", "Makefile")
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n\ncheck:\n\t@true\n")
	secondCommit := commitMechanicalFiles(t, repoRoot, "second tooling change", "Makefile")

	result := runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:               repoRoot,
		AuthorizationPath:      authorizationPath,
		ConsumingSpec:          "commit-column",
		DeliveryTargetRevision: target,
		TaskCommits: []speccheck.MechanicalTaskCommit{
			{TaskID: "task_02", SHA: firstCommit},
			{TaskID: "task_02", SHA: secondCommit},
		},
	})

	var report bytes.Buffer
	if err := speccheck.WriteMechanicalResult(&report, result); err != nil {
		t.Fatalf("WriteMechanicalResult() error = %v", err)
	}
	for _, taskCommit := range []string{firstCommit, secondCommit} {
		if count := strings.Count(report.String(), "| task_02 | "+taskCommit+" | granted |"); count != 1 {
			t.Errorf("authorization audit row count for commit %s = %d, want 1:\n%s", taskCommit, count, report.String())
		}
	}
}
