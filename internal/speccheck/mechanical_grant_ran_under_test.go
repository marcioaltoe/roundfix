// Suite: mechanical authorization grant provenance
// Invariant: a Task commit uses its parent's grant only while the delivery target carries the same record.
// Boundary IN: public speccheck mechanical API and real temporary Git histories
// Boundary OUT: Daemon scheduling and QA verdict computation
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const (
	grantRanUnderSlug              = "grant-ran-under"
	grantRanUnderAuthorizationPath = "docs/specs/grant-ran-under/_authorization.md"
	grantRanUnderTaskFile          = "docs/specs/grant-ran-under/task_01.md"
	grantRanUnderNewPath           = ".golangci.yml"
)

func TestAuthPathsAcceptTheGrantTheTaskRanUnder(t *testing.T) {
	t.Parallel()

	repoRoot, fork := newGrantRanUnderRepository(t)
	runMechanicalGit(t, repoRoot, "switch", "main")
	writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
		mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile", grantRanUnderNewPath))
	widening := commitMechanicalFiles(t, repoRoot, "widen authorization on main", grantRanUnderAuthorizationPath)
	writeMechanicalFile(t, repoRoot, "README.md", "delivery target advanced\n")
	target := commitMechanicalFiles(t, repoRoot, "advance main after authorization", "README.md")

	runMechanicalGit(t, repoRoot, "switch", "item")
	runMechanicalGit(t, repoRoot, "cherry-pick", "--quiet", widening)
	consumer := commitGrantRanUnderPath(t, repoRoot, "change newly bounded path")
	result := auditGrantRanUnder(t, repoRoot, target, consumer)

	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	read := singleGrantRanUnderRead(t, result)
	if read.Outcome != spec.AuthorizationGranted {
		t.Fatalf("authorization outcome = %q, want %q", read.Outcome, spec.AuthorizationGranted)
	}
	if read.Source.Revision != widening {
		t.Fatalf("authorization revision = %q, want main widening commit %q (fork %q)", read.Source.Revision, widening, fork)
	}
}

func TestAuthPathsRefuseAParentGrantMainNeverHeld(t *testing.T) {
	t.Parallel()

	t.Run("parent has no record", func(t *testing.T) {
		t.Parallel()

		repoRoot, fork := newGrantRanUnderRepository(t)
		runMechanicalGit(t, repoRoot, "rm", "--quiet", grantRanUnderAuthorizationPath)
		runMechanicalGit(t, repoRoot, "commit", "--quiet", "-m", "remove authorization on item")
		consumer := commitGrantRanUnderPath(t, repoRoot, "change path without parent grant")

		result := auditGrantRanUnder(t, repoRoot, fork, consumer)

		assertMechanicalPathEscapedGrant(t, result, grantRanUnderNewPath, grantRanUnderAuthorizationPath)
		read := singleGrantRanUnderRead(t, result)
		if read.Source.Revision != fork {
			t.Fatalf("authorization revision = %q, want unchanged fork point %q", read.Source.Revision, fork)
		}
	})

	t.Run("record exists only on item branch", func(t *testing.T) {
		t.Parallel()

		repoRoot, fork := newGrantRanUnderRepository(t)
		writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
			mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile", grantRanUnderNewPath))
		commitMechanicalFiles(t, repoRoot, "widen authorization only on item", grantRanUnderAuthorizationPath)
		consumer := commitGrantRanUnderPath(t, repoRoot, "change path from item-only grant")

		result := auditGrantRanUnder(t, repoRoot, fork, consumer)

		assertMechanicalPathEscapedGrant(t, result, grantRanUnderNewPath, grantRanUnderAuthorizationPath)
		read := singleGrantRanUnderRead(t, result)
		if read.Source.Revision != fork {
			t.Fatalf("authorization revision = %q, want unchanged fork point %q", read.Source.Revision, fork)
		}
	})

	t.Run("older main grant was narrowed", func(t *testing.T) {
		t.Parallel()

		repoRoot, fork := newGrantRanUnderRepository(t)
		runMechanicalGit(t, repoRoot, "switch", "main")
		writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
			mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile", grantRanUnderNewPath))
		widening := commitMechanicalFiles(t, repoRoot, "widen authorization on main", grantRanUnderAuthorizationPath)

		runMechanicalGit(t, repoRoot, "switch", "item")
		runMechanicalGit(t, repoRoot, "cherry-pick", "--quiet", widening)
		runMechanicalGit(t, repoRoot, "switch", "main")
		writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
			mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile"))
		narrowing := commitMechanicalFiles(t, repoRoot, "narrow authorization on main", grantRanUnderAuthorizationPath)

		runMechanicalGit(t, repoRoot, "switch", "item")
		consumer := commitGrantRanUnderPath(t, repoRoot, "change path from revoked grant")
		result := auditGrantRanUnder(t, repoRoot, narrowing, consumer)

		assertMechanicalPathEscapedGrant(t, result, grantRanUnderNewPath, grantRanUnderAuthorizationPath)
		read := singleGrantRanUnderRead(t, result)
		if read.Source.Revision != fork {
			t.Fatalf("authorization revision = %q, want unchanged fork point %q after main narrowed the grant", read.Source.Revision, fork)
		}
	})
}

func TestAuthPathsStillRefuseSelfApproval(t *testing.T) {
	t.Parallel()

	repoRoot, fork := newGrantRanUnderRepository(t)
	runMechanicalGit(t, repoRoot, "switch", "main")
	writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
		mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile", grantRanUnderNewPath))
	widening := commitMechanicalFiles(t, repoRoot, "widen authorization on main", grantRanUnderAuthorizationPath)
	runMechanicalGit(t, repoRoot, "switch", "item")
	runMechanicalGit(t, repoRoot, "cherry-pick", "--quiet", widening)

	writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
		mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile", grantRanUnderNewPath)+"\nself approval\n")
	writeMechanicalFile(t, repoRoot, grantRanUnderNewPath, "linters: {}\n")
	consumer := commitMechanicalFiles(t, repoRoot, "edit grant while consuming it", grantRanUnderAuthorizationPath, grantRanUnderNewPath)
	result := auditGrantRanUnder(t, repoRoot, widening, consumer)

	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
	if len(findings) == 0 {
		t.Fatalf("%s findings = none, want self-approval refusal", speccheck.CodeMechanicalAuthPaths)
	}
	want := "changes authorization grant " + grantRanUnderAuthorizationPath + " in the commit that consumes it"
	found := false
	for _, finding := range findings {
		if strings.Contains(finding.Detail, want) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("%s findings = %#v, want existing self-approval text %q", speccheck.CodeMechanicalAuthPaths, findings, want)
	}
	read := singleGrantRanUnderRead(t, result)
	if read.Source.Revision != fork {
		t.Fatalf("authorization revision = %q, want self-approval check to keep fork point %q", read.Source.Revision, fork)
	}
}

func TestAuthPathsKeepTheForkPointRevisionWhenItCovers(t *testing.T) {
	t.Parallel()

	repoRoot, fork := newGrantRanUnderRepository(t)
	writeMechanicalFile(t, repoRoot, "Makefile", "verify:\n\t@true\n")
	consumer := commitMechanicalFiles(t, repoRoot, "change path covered at fork", "Makefile")
	result := auditGrantRanUnder(t, repoRoot, fork, consumer)

	assertNoMechanicalCode(t, result, speccheck.CodeMechanicalAuthPaths)
	read := singleGrantRanUnderRead(t, result)
	if read.Source.Revision != fork {
		t.Fatalf("authorization revision = %q, want fork point %q", read.Source.Revision, fork)
	}
}

func newGrantRanUnderRepository(t *testing.T) (string, string) {
	t.Helper()
	repoRoot := newMechanicalGitRepo(t)
	writeMechanicalFile(t, repoRoot, grantRanUnderAuthorizationPath,
		mechanicalTypedAuthorization(grantRanUnderSlug, "Makefile"))
	fork := commitMechanicalFiles(t, repoRoot, "record narrow authorization", grantRanUnderAuthorizationPath)
	runMechanicalGit(t, repoRoot, "switch", "-c", "item")
	writeMechanicalFile(t, repoRoot, grantRanUnderTaskFile, "---\nstatus: in_progress\n---\n")
	commitMechanicalFiles(t, repoRoot, "start item branch", grantRanUnderTaskFile)
	return repoRoot, fork
}

func commitGrantRanUnderPath(t *testing.T, repoRoot, message string) string {
	t.Helper()
	writeMechanicalFile(t, repoRoot, grantRanUnderNewPath, "linters: {}\n")
	return commitMechanicalFiles(t, repoRoot, message, grantRanUnderNewPath)
}

func auditGrantRanUnder(t *testing.T, repoRoot, target, consumer string) speccheck.MechanicalResult {
	t.Helper()
	return runMechanical(t, speccheck.MechanicalRequest{
		RepoRoot:               repoRoot,
		AuthorizationPath:      grantRanUnderAuthorizationPath,
		ConsumingSpec:          grantRanUnderSlug,
		DeliveryTargetRevision: target,
		TaskCommits: []speccheck.MechanicalTaskCommit{{
			TaskID:   "task_01",
			SHA:      consumer,
			TaskFile: grantRanUnderTaskFile,
		}},
	})
}

func singleGrantRanUnderRead(t *testing.T, result speccheck.MechanicalResult) speccheck.MechanicalAuthorizationRead {
	t.Helper()
	if len(result.AuthorizationReads) != 1 {
		t.Fatalf("AuthorizationReads = %#v, want one read", result.AuthorizationReads)
	}
	return result.AuthorizationReads[0]
}
