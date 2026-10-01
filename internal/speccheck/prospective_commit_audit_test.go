package speccheck_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const prospectiveGrant = "docs/workflow/authorizations/mechanical.md"

func prospectiveFixture(t *testing.T) (speccheck.MechanicalRequest, speccheck.ProspectiveTaskCommit) {
	t.Helper()
	root := newMechanicalGitRepo(t)
	writeMechanicalFile(t, root, prospectiveGrant, mechanicalTypedAuthorization("mechanical", "Makefile"))
	parent := commitMechanicalFiles(t, root, "grant", prospectiveGrant)
	return speccheck.MechanicalRequest{RepoRoot: root, AuthorizationPath: prospectiveGrant, ConsumingSpec: "mechanical", DeliveryTargetRevision: parent}, speccheck.ProspectiveTaskCommit{TaskID: "task_02", TaskFile: "docs/specs/mechanical/task_02.md", Parent: parent}
}

func prospectiveFindings(t *testing.T, req speccheck.MechanicalRequest, commit speccheck.ProspectiveTaskCommit) []speccheck.MechanicalFinding {
	t.Helper()
	findings, err := speccheck.AuditProspectiveTaskCommit(context.Background(), req, commit)
	if err != nil {
		t.Fatal(err)
	}
	return findings
}

func TestProspectiveAuditRefusesAGovernedPathOutsideTheGrant(t *testing.T) {
	t.Parallel()
	req, commit := prospectiveFixture(t)
	commit.Changed = []string{commit.TaskFile, "Makefile", ".agents/skills/escaped/SKILL.md"}
	findings := prospectiveFindings(t, req, commit)
	want := "Task task_02's prospective commit changes .agents/skills/escaped/SKILL.md outside authorization grant " + prospectiveGrant + "'s exact bounded files"
	if len(findings) != 1 || findings[0].Detail != want {
		t.Fatalf("findings = %+v, want %s", findings, want)
	}
}

func TestProspectiveAuditAcceptsBoundedAndSanctionedPaths(t *testing.T) {
	t.Parallel()
	root := newMechanicalGitRepo(t)
	parent := commitMechanicalRegenerationFixture(t, root, prospectiveGrant, "command: "+mechanicalRegenerationCommand+"\noutputs:\n  - "+mechanicalEnumeratedOutput+"\n")
	req := speccheck.MechanicalRequest{RepoRoot: root, AuthorizationPath: prospectiveGrant, ConsumingSpec: "mechanical", DeliveryTargetRevision: parent}
	commit := speccheck.ProspectiveTaskCommit{TaskID: "task_02", TaskFile: "docs/specs/mechanical/task_02.md", Parent: parent, Changed: []string{"Makefile", mechanicalEnumeratedOutput, "src/ordinary.go"}}
	if findings := prospectiveFindings(t, req, commit); len(findings) != 0 {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestProspectiveAuditRefusesSelfApproval(t *testing.T) {
	t.Parallel()
	req, commit := prospectiveFixture(t)
	commit.Changed = []string{prospectiveGrant, "Makefile"}
	findings := prospectiveFindings(t, req, commit)
	if len(findings) != 1 || !strings.Contains(findings[0].Detail, "changes authorization grant "+prospectiveGrant+" in the commit that consumes it") {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestProspectiveAuditMatchesTheGateAuditOfTheCreatedCommit(t *testing.T) {
	t.Parallel()
	for _, paths := range [][]string{{"Makefile"}, {"Makefile", ".agents/skills/escaped/SKILL.md"}, {prospectiveGrant, "Makefile"}} {
		req, commit := prospectiveFixture(t)
		commit.Changed = paths
		prospective := prospectiveFindings(t, req, commit)
		for _, path := range paths {
			writeMechanicalFile(t, req.RepoRoot, path, "changed\n")
		}
		sha := commitMechanicalFiles(t, req.RepoRoot, "consumer", paths...)
		req.TaskCommits = []speccheck.MechanicalTaskCommit{{TaskID: commit.TaskID, TaskFile: commit.TaskFile, SHA: sha}}
		result := runMechanical(t, req)
		actual := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalAuthPaths)
		for i := range actual {
			actual[i].Detail = strings.ReplaceAll(actual[i].Detail, "Task "+commit.TaskID+" commit "+sha, "Task "+commit.TaskID+"'s prospective commit")
		}
		if !reflect.DeepEqual(prospective, actual) {
			t.Fatalf("paths %v: prospective %+v != gate %+v", paths, prospective, actual)
		}
	}
}

func TestProspectiveAuditMatchesTheGateWithoutAnAuthorizationReference(t *testing.T) {
	t.Parallel()
	req, commit := prospectiveFixture(t)
	req.AuthorizationPath = ""
	commit.Changed = []string{"Makefile"}
	prospective := prospectiveFindings(t, req, commit)
	writeMechanicalFile(t, req.RepoRoot, "Makefile", "changed\n")
	sha := commitMechanicalFiles(t, req.RepoRoot, "consumer", "Makefile")
	req.TaskCommits = []speccheck.MechanicalTaskCommit{{TaskID: commit.TaskID, TaskFile: commit.TaskFile, SHA: sha}}
	gate := runMechanical(t, req)
	if len(prospective) != 0 || len(mechanicalFindingsWithCode(gate, speccheck.CodeMechanicalAuthPaths)) != 0 {
		t.Fatalf("prospective %+v, gate %+v", prospective, gate.Findings)
	}
}
