// Suite: Surface Transcripts
// Invariant: written command surfaces are shaped and traced without executing them.
// Boundary IN: temporary artifacts, disposable Git histories and public checks.
// Boundary OUT: command execution, QA reproduction and Daemon settlement.
package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

const transcriptTech = "docs/specs/0200-example/_techspec.md"
const transcriptBlock = "```transcript\n$ never-execute-this-command\nstdout:\nhello\nstderr:\nexit: 0\n```\n"

func transcriptFixture(t *testing.T, held bool, section string) string {
	t.Helper()
	root := t.TempDir()
	writeReceiptFixture(t, root, receiptFixturePRD, "---\nspec: 0200-example\nstatus: active\n---\n# Example\n")
	if held {
		writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	}
	writeReceiptFixture(t, root, transcriptTech, "# TechSpec\n"+section)
	return root
}
func transcriptSection(block string) string {
	return "## Surface Transcripts\n\n1. Surface Transcript: example\n\n" + block
}
func assertTranscriptCount(t *testing.T, result Result, code string, count int) []Finding {
	t.Helper()
	found := receiptCodeFindings(result, code)
	if len(found) != count {
		t.Fatalf("%s = %#v, want %d", code, found, count)
	}
	return found
}
func TestAHeldTechSpecWithoutSurfaceTranscriptsIsAGap(t *testing.T) {
	root := transcriptFixture(t, true, "")
	found := assertTranscriptCount(t, receiptStage(t, root, StageTechSpec), CodeTranscriptUndeclared, 1)
	if found[0].Severity != SeverityGap {
		t.Fatal(found)
	}
}
func TestSurfaceTranscriptsNoneWithAReasonIsAccepted(t *testing.T) {
	for _, heading := range []string{"##", "###"} {
		root := transcriptFixture(t, true, heading+" Surface Transcripts\n\nNone. No command surface changes.\n")
		assertTranscriptCount(t, receiptStage(t, root, StageTechSpec), CodeTranscriptUndeclared, 0)
	}
}
func TestAWellFormedSurfaceTranscriptReportsNothing(t *testing.T) {
	root := transcriptFixture(t, true, transcriptSection(transcriptBlock))
	result := receiptStage(t, root, StageTechSpec)
	assertTranscriptCount(t, result, CodeTranscriptUndeclared, 0)
	assertTranscriptCount(t, result, CodeTranscriptMalformed, 0)
	parsed := SurfaceTranscripts([]byte(transcriptSection(transcriptBlock)))
	if len(parsed) != 1 || parsed[0].Command != "never-execute-this-command" || strings.Join(parsed[0].Stdout, "\n") != "hello" || len(parsed[0].Stderr) != 0 || parsed[0].ExitCode != 0 || parsed[0].Line != 3 || parsed[0].Title != "example" {
		t.Fatalf("parsed = %#v", parsed)
	}
}
func TestEachMalformedSurfaceTranscriptNamesItsReason(t *testing.T) {
	cases := []struct{ reason, block string }{
		{"without a transcript block", ""},
		{"with more than one transcript block", transcriptBlock + transcriptBlock},
		{"without a command line", strings.Replace(transcriptBlock, "$ never-execute-this-command", "command", 1)},
		{"without a stdout line", strings.Replace(transcriptBlock, "stdout:", "out:", 1)},
		{"without a stderr line after its stdout line", strings.Replace(transcriptBlock, "stderr:", "err:", 1)},
		{"without an exit line", strings.Replace(transcriptBlock, "exit: 0\n", "", 1)},
		{"with an exit code outside 0-255", strings.Replace(transcriptBlock, "exit: 0", "exit: 256", 1)},
	}
	for _, c := range cases {
		for _, held := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/held=%v", c.reason, held), func(t *testing.T) {
				root := transcriptFixture(t, held, transcriptSection(c.block))
				findings := assertTranscriptCount(t, receiptStage(t, root, StageTechSpec), CodeTranscriptMalformed, 1)
				if findings[0].Severity != SeverityError || !strings.HasSuffix(findings[0].Summary, c.reason) {
					t.Fatal(findings)
				}
			})
		}
	}

}
func TestNumberedOutputInsideATranscriptIsNotAnItem(t *testing.T) {
	block := strings.Replace(transcriptBlock, "hello", "2. Surface Transcript: output\n## Surface Transcripts", 1)
	parsed := SurfaceTranscripts([]byte(transcriptSection(block)))
	if len(parsed) != 1 || parsed[0].Malformed != "" {
		t.Fatal(parsed)
	}
}
func transcriptGraph(t *testing.T, root, workRefs, qaRefs, qaReq, status string, declined bool) {
	t.Helper()
	gate := "qa: task_02\n"
	nodes := "    - id: task_01\n      file: task_01.md\n      needs: []\n"
	if declined {
		gate = "qa: declined\nqa_reason: no gate in this fixture\n"
	} else {
		nodes += "    - id: task_02\n      file: task_02.md\n      needs: [task_01]\n"
	}
	writeReceiptFixture(t, root, "docs/specs/0200-example/_tasks.md", "---\nschema: spec-tasks/v1\nspec: 0200-example\n"+gate+"graph:\n  nodes:\n"+nodes+"---\n# Tasks\n")
	for _, task := range []struct{ id, kind, refs, req, status string }{{"task_01", "backend", workRefs, "Implement.", "pending"}, {"task_02", "qa", qaRefs, qaReq, status}} {
		if declined && task.kind == "qa" {
			continue
		}
		writeReceiptFixture(t, root, "docs/specs/0200-example/"+task.id+".md", fmt.Sprintf("---\ntask: %s\nspec: 0200-example\nstatus: %s\ntype: %s\ncomplexity: low\n---\n# Task\n\n## Requirements\n\n1. %s\n\n## Verification\n\n- `true`\n\n## References\n\n%s\n", task.id, task.status, task.kind, task.req, task.refs))
	}
}
func TestASurfaceTranscriptNamedOnlyByTheQATaskIsUntasked(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "", "Surface Transcript 1", "Reproduce Surface Transcript 1.", "pending", false)
	findings := assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeCoverageUntasked, 1)
	if !strings.Contains(findings[0].Summary, "no Task other than the QA gate") {
		t.Fatal(findings)
	}
}
func TestASurfaceTranscriptNamedByANonQATaskIsTasked(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcripts 1-3", "", "Reproduce Surface Transcript 1.", "pending", false)
	assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeCoverageUntasked, 0)
}
func TestASurfaceTranscriptTheQATaskDoesNotNameIsUngated(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcript 1", "Surface Transcript 1", "Reproduce output.", "pending", false)
	assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeTranscriptUngated, 1)
}
func TestASurfaceTranscriptNamedInAQARequirementIsGated(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcript 1", "", "Reproduce Surface Transcripts 1-3.", "pending", false)
	assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeTranscriptUngated, 0)
}
func TestADeclinedGateRaisesNoUngatedTranscript(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcript 1", "", "", "pending", true)
	assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeTranscriptUngated, 0)
}
func TestASpecThatIsNotHeldSkipsTheTranscriptDeclarationGap(t *testing.T) {
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, receiptFixturePRD, "# Example\n")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, transcriptTech, "# TechSpec\n")
	result := receiptStage(t, root, StageTechSpec)
	assertTranscriptCount(t, result, CodeTranscriptUndeclared, 0)
	if !strings.Contains(RenderText(result, VerificationCoverage{}), CodeTranscriptUndeclared+": missing a PRD committed at or after "+ConcreteContractGuidePath) {
		t.Fatal(result.Skipped)
	}
}
func TestAMalformedTranscriptIsReportedInASpecThatIsNotHeld(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(""))
	assertTranscriptCount(t, receiptStage(t, root, StageTechSpec), CodeTranscriptMalformed, 1)
}
func TestTranscriptFindingsRenderSurfaceTranscriptsThreeAndFour(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(strings.Replace(transcriptBlock, "exit: 0\n", "", 1))+"\n2. Surface Transcript: second\n\n"+transcriptBlock)
	transcriptGraph(t, root, "Surface Transcripts 1-2", "", "Reproduce Surface Transcript 1.", "pending", false)
	text := RenderText(receiptStage(t, root, StageTasks), VerificationCoverage{})
	for _, want := range []string{
		"[error] SC-TRANSCRIPT-MALFORMED: " + transcriptTech + " declares Surface Transcript 1 without an exit line\n  at " + transcriptTech + ":4\n  fix: Write the block as a command line starting with \"$ \", then \"stdout:\", \"stderr:\" and \"exit: <code>\", in that order.",
		"[error] SC-TRANSCRIPT-UNGATED: " + transcriptTech + " declares Surface Transcript 2, but QA Task docs/specs/0200-example/task_02.md names it in no Requirement\n  at " + transcriptTech + ":13\n  at docs/specs/0200-example/task_02.md:12\n  fix: Name Surface Transcript 2 in a Requirement of docs/specs/0200-example/task_02.md so the gate reproduces it.",
		"fix: Name Surface Transcript 2 in a Requirement of docs/specs/0200-example/task_02.md so the gate reproduces it.",
		"SC-TRANSCRIPT-UNDECLARED: missing " + ConcreteContractGuidePath,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}
}
func TestCompletedQATaskIsHistoricalTranscriptEvidence(t *testing.T) {
	root := transcriptFixture(t, false, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcript 1", "", "Historical.", "pending", false)
	// Use the detector seam: full graph loading also validates completed QA reports.
	result := Result{}
	graph, present, err := loadOptionalTaskGraph(filepath.Join(root, "docs/specs"), "0200-example", filepath.Join(root, "docs/specs/0200-example"))
	if err != nil {
		t.Fatal(err)
	}
	if !present {
		t.Fatal("missing graph")
	}
	for index := range graph.Tasks {
		if graph.Tasks[index].ID == graph.QATaskID {
			graph.Tasks[index].Status = spec.StatusCompleted
		}
	}
	detectTranscriptGate(&result, root, graph, coverageUnitsDeclaredIn([]coverageUnit{{Kind: coverageTranscript, Number: 1, Line: 5}}, transcriptTech))
	assertTranscriptCount(t, result, CodeTranscriptUngated, 0)
}
func TestThisSpecsSurfaceTranscriptsAreWellFormed(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(thisSpecDir(t), "_techspec.md"))
	if err != nil {
		t.Fatal(err)
	}
	parsed := SurfaceTranscripts(content)
	if len(parsed) != 4 {
		t.Fatal(parsed)
	}
	for _, item := range parsed {
		if item.Malformed != "" {
			t.Fatal(item)
		}
	}
}

func TestTranscriptDeclarationRejectsMissingReasonsAndWrongShapes(t *testing.T) {
	for _, section := range []string{
		"## Surface Transcripts\n\nNone.\n",
		"## Surface Transcripts\n\nNone. Reason.\n\n1. Other entry.\n",
		"# Surface Transcripts\n\n" + transcriptSection(transcriptBlock)[len("## Surface Transcripts\n\n"):],
		"## Surface Transcripts\n\n1. Other entry.\n" + transcriptBlock,
		"```markdown\n" + transcriptSection(transcriptBlock) + "```\n",
	} {
		root := transcriptFixture(t, true, section)
		assertTranscriptCount(t, receiptStage(t, root, StageTechSpec), CodeTranscriptUndeclared, 1)
	}
}

func TestTranscriptBlockBoundariesAndExactStreams(t *testing.T) {
	cases := []struct{ name, body, reason string }{
		{"empty streams", "$ command\nstdout:\nstderr:\nexit: 255", ""},
		{"first stderr ends stdout", "$ command\nstdout:\n...\n<line>\nstderr:\nstderr:\nexit: 1", ""},
		{"empty command", "$ \nstdout:\nstderr:\nexit: 0", "without a command line"},
		{"blank before command", "\n$ command\nstdout:\nstderr:\nexit: 0", "without a command line"},
		{"stdout must be second", "$ command\nextra\nstdout:\nstderr:\nexit: 0", "without a stdout line"},
		{"stderr must follow stdout", "$ command\nstderr:\nstdout:\nexit: 0", "without a stdout line"},
		{"exit must be last", "$ command\nstdout:\nstderr:\nexit: 0\nextra", "without an exit line"},
		{"noninteger exit", "$ command\nstdout:\nstderr:\nexit: nope", "with an exit code outside 0-255"},
		{"negative exit", "$ command\nstdout:\nstderr:\nexit: -1", "with an exit code outside 0-255"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			section := transcriptSection("   ~~~transcript\n   " + strings.ReplaceAll(c.body, "\n", "\n   ") + "\n   ~~~~\n")
			parsed := SurfaceTranscripts([]byte(section))
			if len(parsed) != 1 || parsed[0].Malformed != c.reason {
				t.Fatalf("parsed = %#v, want %q", parsed, c.reason)
			}
			if c.name == "first stderr ends stdout" && (strings.Join(parsed[0].Stdout, "\n") != "...\n<line>" || strings.Join(parsed[0].Stderr, "\n") != "stderr:") {
				t.Fatal(parsed)
			}
		})
	}
	for _, info := range []string{"text", "transcript extra", "Transcript"} {
		parsed := SurfaceTranscripts([]byte(transcriptSection(strings.Replace(transcriptBlock, "```transcript", "```"+info, 1))))
		if len(parsed) != 1 || parsed[0].Malformed != "without a transcript block" {
			t.Fatal(parsed)
		}
	}
}

func TestTranscriptDetectorStagesAndMissingArtifacts(t *testing.T) {
	root := transcriptFixture(t, true, transcriptSection(""))
	result := receiptStage(t, root, StagePRD)
	assertTranscriptCount(t, result, CodeTranscriptMalformed, 0)
	for _, code := range []string{CodeTranscriptUndeclared, CodeTranscriptMalformed, CodeTranscriptUngated} {
		if !strings.Contains(RenderText(result, VerificationCoverage{}), code+": missing stage ") {
			t.Fatal(result.Skipped)
		}
	}
	if err := os.Remove(filepath.Join(root, transcriptTech)); err != nil {
		t.Fatal(err)
	}
	result = receiptStage(t, root, StageTechSpec)
	assertTranscriptCount(t, result, CodeTranscriptUndeclared, 0)
	if !strings.Contains(RenderText(result, VerificationCoverage{}), CodeTranscriptUndeclared+": missing "+transcriptTech) {
		t.Fatal(result.Skipped)
	}
	if err := os.Remove(filepath.Join(root, receiptFixturePRD)); err != nil {
		t.Fatal(err)
	}
	result = receiptStage(t, root, StageAll)
	for _, code := range []string{CodeTranscriptUndeclared, CodeTranscriptMalformed, CodeTranscriptUngated} {
		if !strings.Contains(RenderText(result, VerificationCoverage{}), code+": missing "+receiptFixturePRD) {
			t.Fatal(result.Skipped)
		}
	}
}

func TestSurfaceTranscriptsNeedNoCoverageMap(t *testing.T) {
	root := transcriptFixture(t, true, transcriptSection(transcriptBlock))
	transcriptGraph(t, root, "Surface Transcript 1", "", "Reproduce Surface Transcript 1.", "pending", false)
	assertTranscriptCount(t, receiptStage(t, root, StageTasks), CodeCoverageUnmapped, 0)
}

func TestThisSpecsTranscriptsHaveImplementationAndGateReferences(t *testing.T) {
	root := t.TempDir()
	slug := "0191-claims-with-receipts-and-contracts-as-they-ship"
	for _, name := range []string{"_prd.md", "_techspec.md", "_tasks.md", "task_01.md", "task_02.md", "task_03.md", "task_04.md"} {
		path := filepath.Join("docs/specs", slug, name)
		content, err := os.ReadFile(filepath.Join(thisSpecDir(t), name))
		if err != nil {
			t.Fatal(err)
		}
		text := string(content)
		if name == "_prd.md" {
			// The fixture replays the Spec as active, even after its archive.
			text = strings.Replace(text, "status: archived", "status: active", 1)
		}
		writeReceiptFixture(t, root, path, text)
	}
	result, err := CheckStage(filepath.Join(root, "docs/specs"), root, slug, StageTasks)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{CodeTranscriptUndeclared, CodeTranscriptMalformed, CodeTranscriptUngated} {
		assertTranscriptCount(t, result, code, 0)
	}
	for _, finding := range receiptCodeFindings(result, CodeCoverageUntasked) {
		if strings.Contains(finding.Summary, "Surface Transcript") {
			t.Fatal(finding)
		}
	}
}

// thisSpecDir finds Spec 0191 where it lives: active under docs/specs, or
// archived under docs/history/specs once its delivery archives it.
func thisSpecDir(t *testing.T) string {
	t.Helper()
	const slug = "0191-claims-with-receipts-and-contracts-as-they-ship"
	for _, root := range []string{"../../docs/specs", "../../docs/history/specs"} {
		dir := filepath.Join(root, slug)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	t.Fatalf("Spec %s is neither active nor archived", slug)
	return ""
}
