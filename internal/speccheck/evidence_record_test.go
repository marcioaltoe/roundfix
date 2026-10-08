// Suite: QA evidence recording
// Invariant: only qualifying passing rows retain Daemon snapshots of Git blobs.
// Boundary IN: recorder, report reader, mechanical stage and temporary Git repositories
// Boundary OUT: Daemon settlement (internal/daemon/qa_evidence_snapshot_test.go)
package speccheck_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const recordReportPath = "docs/specs/mechanical/qa/qa-report-2026-08-11.md"

func recordFixture(t *testing.T) (string, string, string) {
	t.Helper()
	root := newMechanicalGitRepo(t)
	writeMechanicalFile(t, root, "evidence.txt", "stable evidence\n")
	head := commitMechanicalFiles(t, root, "establish evidence", "evidence.txt")
	report := mechanicalCarryReport(head, "stable evidence\n")
	start := strings.Index(report, "evidence_snapshots:")
	end := strings.Index(report[start:], "---\n") + start
	report = report[:start] + report[end:]
	return root, head, report
}

func recordRun(t *testing.T, root, head, report string) (speccheck.EvidenceRecord, string) {
	t.Helper()
	writeMechanicalFile(t, root, recordReportPath, report)
	record, err := speccheck.RecordEvidenceSnapshots(context.Background(), root, filepath.Join(root, recordReportPath), head)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(root, recordReportPath))
	if err != nil {
		t.Fatal(err)
	}
	return record, string(content)
}

func recordedRows(t *testing.T, content string) map[string]struct {
	Head   string `yaml:"head"`
	Inputs []struct {
		Ref    string `yaml:"ref"`
		Count  int    `yaml:"count"`
		SHA256 string `yaml:"sha256"`
	} `yaml:"inputs"`
} {
	t.Helper()
	var parsed struct {
		Rows map[string]struct {
			Head   string `yaml:"head"`
			Inputs []struct {
				Ref    string `yaml:"ref"`
				Count  int    `yaml:"count"`
				SHA256 string `yaml:"sha256"`
			} `yaml:"inputs"`
		} `yaml:"evidence_snapshots"`
	}
	front := strings.SplitN(strings.ReplaceAll(content, "\r\n", "\n"), "---", 3)[1]
	if err := yaml.Unmarshal([]byte(front), &parsed); err != nil {
		t.Fatal(err)
	}
	return parsed.Rows
}

func TestEvidenceRecordSnapshotsEveryQualifyingPassRow(t *testing.T) {
	t.Parallel()
	root, head, report := recordFixture(t)
	writeMechanicalFile(t, root, "src/z.txt", "z\n")
	writeMechanicalFile(t, root, "src/a.txt", "a\n")
	head = commitMechanicalFiles(t, root, "more inputs", "src")
	report = strings.Replace(report, "| R01 | Unchanged", "| R02 | Unchanged", 1) + "\n### R02 evidence\n\n```yaml\ninputs:\n  - kind: repository_path\n    ref: src/**\n  - kind: repository_path\n    ref: evidence.txt\n```\n"
	// Both passing rows qualify, in table order rather than sorted identifier order.
	report = strings.Replace(report, "| R02 | Unchanged evidence | maintainer / backend | pass | [evidence](../../../../evidence.txt) |", "| R02 | Unchanged evidence | maintainer / backend | pass | [evidence](../../../../evidence.txt) |\n| R01 | Stable | maintainer | pass | [evidence](../../../../evidence.txt) |", 1)
	writeMechanicalFile(t, root, "evidence.txt", "dirty worktree must not be hashed\n")
	record, content := recordRun(t, root, head, report)
	if record.Head != head || !reflect.DeepEqual(record.Rows, []string{"R02", "R01"}) {
		t.Fatalf("record = %+v", record)
	}
	rows := recordedRows(t, content)
	if len(rows) != 2 || rows["R02"].Head != head {
		t.Fatalf("snapshots = %+v", rows)
	}
	inputs := rows["R02"].Inputs
	if len(inputs) != 2 || inputs[0].Ref != "src/**" || inputs[1].Ref != "evidence.txt" {
		t.Fatalf("input order = %+v", inputs)
	}
	var summary strings.Builder
	for _, path := range []string{"src/a.txt", "src/z.txt"} {
		fmt.Fprintf(&summary, "%x  %s\n", sha256.Sum256([]byte(strings.TrimPrefix(strings.TrimSuffix(path, ".txt"), "src/")+"\n")), path)
	}
	if inputs[0].Count != 2 || inputs[0].SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(summary.String()))) {
		t.Fatalf("glob pair = %+v", inputs[0])
	}
	expected := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x  evidence.txt\n", sha256.Sum256([]byte("stable evidence\n"))))))
	if inputs[1].Count != 1 || inputs[1].SHA256 != expected {
		t.Fatalf("hashed worktree instead of Git: %+v", inputs[1])
	}
}

func TestEvidenceRecordSkipsRowsThatCannotCarry(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, from, to string }{
		{"failed", "| pass |", "| fail |"},
		{"blocked", "| pass |", "| blocked (environment: offline) |"},
		{"skipped", "| pass |", "| skipped |"},
		{"no inputs", "inputs:", "other:"},
		{"external repository", "kind: repository_path", "kind: external_repository"},
		{"live service", "kind: repository_path", "kind: live_service"},
		{"elapsed time", "kind: repository_path", "kind: elapsed_time"},
		{"commit range", "kind: repository_path", "kind: " + string(speccheck.EvidenceCommitRange)},
		{"repository verification", "| pass |", "| pass | repository Verification |"},
		{"pull request", "| pass |", "| pass | " + spec.QAPullRequestRowSource + " |"},
		{"absent input", "ref: evidence.txt", "ref: missing.txt"},
		{"uncovered evidence", "ref: evidence.txt", "ref: other.txt"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, head, report := recordFixture(t)
			writeMechanicalFile(t, root, "other.txt", "other\n")
			head = commitMechanicalFiles(t, root, "other", "other.txt")
			if tc.name == "repository verification" || tc.name == "pull request" {
				report = strings.Replace(report, "Status | Evidence", "Status | Provenance | Evidence", 1)
			}
			record, content := recordRun(t, root, head, strings.Replace(report, tc.from, tc.to, 1))
			if len(record.Rows) != 0 || strings.Contains(content, "evidence_snapshots:") {
				t.Fatalf("ineligible row recorded: %+v\n%s", record, content)
			}
		})
	}
}

func TestEvidenceRecordReplacesAnAgentWrittenKey(t *testing.T) {
	t.Parallel()
	for _, qualifies := range []bool{true, false} {
		t.Run(fmt.Sprint(qualifies), func(t *testing.T) {
			root, head, report := recordFixture(t)
			report = strings.Replace(report, "verdict: pass", "evidence_snapshots: fabricated\n  fake: value\nverdict: pass", 1)
			if !qualifies {
				report = strings.Replace(report, "| pass |", "| fail |", 1)
			}
			record, content := recordRun(t, root, head, report)
			if strings.Contains(content, "fabricated") || strings.Contains(content, "fake:") {
				t.Fatal(content)
			}
			if qualifies && len(recordedRows(t, content)) != 1 || !qualifies && len(record.Rows) != 0 {
				t.Fatal(content)
			}
		})
	}
}

func TestEvidenceRecordKeepsEveryOtherByte(t *testing.T) {
	t.Parallel()
	for _, newline := range []string{"\n", "\r\n"} {
		t.Run(fmt.Sprintf("newline %q", newline), func(t *testing.T) {
			root, head, report := recordFixture(t)
			report = strings.Replace(report, "status: closed", "# keep this comment\nstatus: closed  # exact spaces", 1)
			report = strings.ReplaceAll(report, "\n", newline)
			writeMechanicalFile(t, root, recordReportPath, report)
			before, err := spec.ReadQAReportFile(filepath.Join(root, recordReportPath))
			if err != nil && newline == "\n" {
				t.Fatal(err)
			}
			_, content := recordRun(t, root, head, report)
			start := strings.Index(content, "evidence_snapshots:")
			end := start + strings.Index(content[start:], "---"+newline)
			if content[:start]+content[end:] != report {
				t.Fatalf("other bytes changed:\n%s", content)
			}
			after, err := spec.ReadQAReportFile(filepath.Join(root, recordReportPath))
			if err != nil && newline == "\n" {
				t.Fatal(err)
			}
			if newline == "\n" && !reflect.DeepEqual(before, after) {
				t.Fatalf("report semantics changed: %+v => %+v", before, after)
			}
		})
	}
}

func TestEvidenceRecordRoundTripsThroughTheMechanicalStage(t *testing.T) {
	t.Parallel()
	root, head, report := recordFixture(t)
	recordRun(t, root, head, report)
	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: root, ReportPath: recordReportPath})
	if len(result.Carried) != 1 || result.Carried[0].ID != "R01" || result.Carried[0].EstablishedHead != head {
		t.Fatalf("carry = %+v", result.Carried)
	}
}

func TestAlwaysObservedNamesItsReason(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		provenance string
		inputs     []speccheck.EvidenceInput
		reason     string
	}{
		{"criterion; repository Verification, another", nil, speccheck.CarryReasonRepositoryVerification},
		{"criterion, " + spec.QAPullRequestRowSource, nil, speccheck.CarryReasonPullRequestRow},
		{"criterion", []speccheck.EvidenceInput{{Kind: speccheck.EvidenceCommitRange}}, speccheck.CarryReasonCommitRangeInput},
		{"not repository Verification; Pull Request row extra", nil, ""},
		{"criterion", []speccheck.EvidenceInput{{Kind: speccheck.EvidenceRepositoryPath, Ref: "evidence.txt"}}, ""},
	} {
		reason, observed := speccheck.AlwaysObserved(tc.provenance, tc.inputs)
		if reason != tc.reason || observed != (tc.reason != "") {
			t.Fatalf("reason=%q observed=%v for %+v", reason, observed, tc)
		}
	}
}

func TestEvidenceRecordStripsTheKeyOnGitReadError(t *testing.T) {
	t.Parallel()
	root, _, report := recordFixture(t)
	report = strings.Replace(report, "verdict: pass", "evidence_snapshots: fabricated\nverdict: pass", 1)
	writeMechanicalFile(t, root, recordReportPath, report)
	record, err := speccheck.RecordEvidenceSnapshots(context.Background(), root, filepath.Join(root, recordReportPath), "missing-head")
	if err == nil || len(record.Rows) != 0 {
		t.Fatalf("record=%+v err=%v", record, err)
	}
	content, err := os.ReadFile(filepath.Join(root, recordReportPath))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "evidence_snapshots:") {
		t.Fatal(string(content))
	}
}
