// Suite: bounded evidence snapshots
// Boundary IN: recorder, mechanical carry stage, real temporary Git repositories
// Boundary OUT: Daemon report settlement and prior-pass import
package speccheck_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

// This oracle builds the summary from content, independently of production.
func inputSummaryDigest(contents map[string]string) string {
	paths := make([]string, 0, len(contents))
	for path := range contents {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var summary strings.Builder
	for _, path := range paths {
		fmt.Fprintf(&summary, "%x  %s\n", sha256.Sum256([]byte(contents[path])), path)
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(summary.String())))
}

func digestFixture(t *testing.T, contents map[string]string) (string, string, string) {
	t.Helper()
	root, _, report := recordFixture(t)
	for path, content := range contents {
		writeMechanicalFile(t, root, path, content)
	}
	head := commitMechanicalFiles(t, root, "establish glob", "src")
	report = strings.ReplaceAll(report, "evidence.txt", "src/a.txt")
	report = strings.Replace(report, "ref: src/a.txt", "ref: src/**", 1)
	return root, head, report
}

func assertDigestStage(t *testing.T, root, reason string) {
	t.Helper()
	result := rowCarryResult(t, root)
	if len(result.Findings) != 0 {
		t.Fatalf("findings = %+v", result.Findings)
	}
	if reason != "" {
		rowCarryRefusal(t, result, reason)
		return
	}
	if len(result.Carried) != 1 || len(result.Dispositions) != 1 || !result.Dispositions[0].Carried {
		t.Fatalf("carry = %+v", result)
	}
}

func TestEvidenceRecordWritesOneLinePerInputWhateverItMatches(t *testing.T) {
	contents := map[string]string{"src/a.txt": "a\n"}
	for i := 1; i < 2500; i++ {
		contents[fmt.Sprintf("src/file-%04d.txt", i)] = fmt.Sprint(i)
	}
	root, head, report := digestFixture(t, contents)
	record, content := recordRun(t, root, head, report)
	start := strings.Index(content, "evidence_snapshots:")
	block := strings.SplitN(content[start:], "---\n", 2)[0]
	if got, want := strings.Count(block, "\n"), 1+3*len(record.Rows)+1; got != want {
		t.Fatalf("block lines = %d, want %d", got, want)
	}
	input := recordedRows(t, content)["R01"].Inputs[0]
	if !strings.Contains(block, "count: 2500") || input.Count != 2500 || input.SHA256 != inputSummaryDigest(contents) {
		t.Fatalf("input = %+v", input)
	}
}

func TestEvidenceInputDigestFollowsTheSortedSummary(t *testing.T) {
	contents := map[string]string{"src/z.txt": "z\n", "src/a.txt": "a\n", "src/m.txt": "m\n"}
	root, head, report := recordFixture(t)
	for _, path := range []string{"src/z.txt", "src/a.txt", "src/m.txt"} {
		writeMechanicalFile(t, root, path, contents[path])
		head = commitMechanicalFiles(t, root, "commit "+path, path)
	}
	report = strings.ReplaceAll(report, "evidence.txt", "src/a.txt")
	report = strings.Replace(report, "ref: src/a.txt", "ref: src/**", 1)
	_, content := recordRun(t, root, head, report)
	input := recordedRows(t, content)["R01"].Inputs[0]
	if input.Count != 3 || input.SHA256 != inputSummaryDigest(contents) {
		t.Fatalf("input = %+v", input)
	}
	t.Run("newline path is unresolved", func(t *testing.T) {
		writeMechanicalFile(t, root, "src/new\nline.txt", "ambiguous")
		head = commitMechanicalFiles(t, root, "newline path", "src")
		record, _ := recordRun(t, root, head, report)
		if len(record.Rows) != 0 {
			t.Fatalf("newline input recorded: %+v", record)
		}
	})
}

func TestCarryComparesTheRecordedDigestPerInput(t *testing.T) {
	contents := map[string]string{"src/a.txt": "a\n", "src/z.txt": "z\n"}
	for _, name := range []string{"equal", "count", "digest"} {
		t.Run(name, func(t *testing.T) {
			root, head, report := digestFixture(t, contents)
			_, content := recordRun(t, root, head, report)
			reason := ""
			if name == "count" {
				content = strings.Replace(content, "count: 2", "count: 3", 1)
				reason = speccheck.CarryReasonInputMoved + "src/**"
			}
			if name == "digest" {
				content = strings.Replace(content, inputSummaryDigest(contents), strings.Repeat("0", 64), 1)
				reason = speccheck.CarryReasonInputMoved + "src/**"
			}
			writeMechanicalFile(t, root, recordReportPath, content)
			assertDigestStage(t, root, reason)
		})
	}
}

func TestCarryNamesTheMovedRefNotItsFiles(t *testing.T) {
	contents := map[string]string{"src/a.txt": "a\n", "src/b.txt": "b\n", "src/z.txt": "z\n"}
	root, head, report := digestFixture(t, contents)
	recordRun(t, root, head, report)
	assertDigestStage(t, root, "")
	writeMechanicalFile(t, root, "src/a.txt", "a changed\n")
	writeMechanicalFile(t, root, "src/b.txt", "b changed\n")
	writeMechanicalFile(t, root, "src/new.txt", "added\n")
	if err := os.Remove(filepath.Join(root, "src/z.txt")); err != nil {
		t.Fatal(err)
	}
	commitMechanicalFiles(t, root, "move matched files", "src")
	assertDigestStage(t, root, speccheck.CarryReasonInputMoved+"src/**")
	t.Run("declaration order", func(t *testing.T) {
		root, head, report := digestFixture(t, contents)
		report = strings.Replace(report, "    ref: src/**", "    ref: src/z.txt\n  - kind: repository_path\n    ref: src/**", 1)
		recordRun(t, root, head, report)
		writeMechanicalFile(t, root, "src/z.txt", "changed")
		commitMechanicalFiles(t, root, "move two declared inputs", "src")
		assertDigestStage(t, root, speccheck.CarryReasonInputMoved+"src/z.txt, src/**")
	})
}

func TestCarryReadsAPerFileSnapshotAsItsDigest(t *testing.T) {
	contents := map[string]string{"src/a.txt": "a\n", "src/z.txt": "z\n"}
	// The independent summary also pins the expected pair used by the new form.
	digest := inputSummaryDigest(contents)
	for _, name := range []string{"valid", "missing", "unsorted", "mixed", "mixed count", "mixed zero count", "mixed empty digest"} {
		t.Run(name, func(t *testing.T) {
			root, head, report := digestFixture(t, contents)
			_, content := recordRun(t, root, head, report)
			first := fmt.Sprintf("          - path: src/a.txt\n            sha256: %x\n", sha256.Sum256([]byte(contents["src/a.txt"])))
			last := fmt.Sprintf("          - path: src/z.txt\n            sha256: %x\n", sha256.Sum256([]byte(contents["src/z.txt"])))
			files, reason, extra := first+last, "", ""
			switch name {
			case "missing":
				files = first
				reason = speccheck.CarryReasonInputMoved + "src/**"
			case "unsorted":
				files = last + first
				reason = speccheck.CarryReasonNoEvidenceSnapshot
			case "mixed":
				extra = "        sha256: " + digest + "\n"
				reason = speccheck.CarryReasonNoEvidenceSnapshot
			case "mixed count":
				extra = "        count: 2\n"
				reason = speccheck.CarryReasonNoEvidenceSnapshot
			case "mixed zero count":
				extra = "        count: 0\n"
				reason = speccheck.CarryReasonNoEvidenceSnapshot
			case "mixed empty digest":
				extra = "        sha256: ''\n"
				reason = speccheck.CarryReasonNoEvidenceSnapshot
			}
			line := fmt.Sprintf("      - {ref: src/**, count: 2, sha256: %s}\n", digest)
			content = strings.Replace(content, line, "      - ref: src/**\n"+extra+"        files:\n"+files, 1)
			if strings.Contains(content, line) || !strings.Contains(content, "files:") {
				t.Fatal("legacy fixture replacement failed")
			}
			writeMechanicalFile(t, root, recordReportPath, content)
			assertDigestStage(t, root, reason)
		})
	}
}

func TestCarriableAcceptsARecordedDigestAgainstTheCurrentFiles(t *testing.T) {
	contents := map[string]string{"src/a.txt": "a\n"}
	prior := speccheck.ReportRow{Status: "pass", EstablishedBy: recordReportPath, EstablishedHead: "prior", AncestryVerified: true, Inputs: []speccheck.EvidenceInput{{Kind: speccheck.EvidenceRepositoryPath, Ref: "src/**"}}, EvidencePaths: []string{"src/a.txt"}}
	established := []speccheck.EvidenceSnapshot{{Ref: "src/**", Count: 1, SHA256: inputSummaryDigest(contents)}}
	current := []speccheck.EvidenceSnapshot{{Ref: "src/**", Files: []speccheck.EvidenceFile{{Path: "src/a.txt", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(contents["src/a.txt"])))}}}}
	if !speccheck.Carriable(prior, "current", nil, established, current) {
		t.Fatal("equal pair refused")
	}
	prior.EvidencePaths = []string{"src/uncited.txt"}
	if speccheck.Carriable(prior, "current", nil, established, current) {
		t.Fatal("uncovered citation carried")
	}
	prior.EvidencePaths = []string{"src/a.txt"}
	established[0].SHA256 = strings.Repeat("0", 64)
	if speccheck.Carriable(prior, "current", nil, established, current) {
		t.Fatal("different pair carried")
	}
}
