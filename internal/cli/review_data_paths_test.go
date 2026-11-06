// Suite: Review Data Paths through the public review command.
// Invariant: data follows code in the same session and keeps its own anchors and record fields.
// Boundary IN: committed Project Config and candidates in temporary Git repositories.
// Boundary OUT: captured provider prompts, persisted records and verbatim answer files.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewSendsADataPathApart(t *testing.T) {
	t.Parallel()
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			t.Parallel()
			codeAnswer := "Findings:\n- code.go:1 Failure: the code loses a rule."
			dataAnswer := "Findings:\n- rules/icms.json:10 Failure: the catalog contradicts the code.\n- code.go:1 Failure: this anchor belongs to the code prompt."
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
				{result: agent.ExecuteResult{Message: codeAnswer, ACPSessionID: "data-session"}},
				{result: agent.ExecuteResult{Message: dataAnswer, ACPSessionID: "data-session"}},
			}}}
			fixture := reviewDataCommandFixture(t, provider, runner, "  data_paths: ['rules/*.json']\n")
			commitScopeFile(t, fixture, "code.go", "package code\n")
			fixture.headCommit = commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(9000))
			code, record, stderr := fixture.run(t)
			if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || len(runner.requests) != 2 || runner.prepareCalls != 1 || runner.endCalls != 0 {
				t.Fatalf("exit=%d outcome=%s prompts=%d prepares=%d ends=%d stderr=%q", code, record.Outcome, len(runner.requests), runner.prepareCalls, runner.endCalls, stderr)
			}
			codePrompt, dataPrompt := runner.requests[0].Prompt, runner.requests[1].Prompt
			wantOmissions := "Omitted from this diff (not reviewed): 0 path(s) of QA evidence, 0 of upstream-managed skills, 0 of archived Spec sources, 0 of retired history, 0 of generated files.\nReviewed apart, after this diff, in Data Review Segments: 1 path(s) of pre_pr_review.data_paths.\nrules/icms.json (data-path)\n"
			if strings.Contains(codePrompt, "\"rule_") || strings.Contains(codePrompt, "diff --git a/rules/") || !strings.Contains(codePrompt, "+package code") || !strings.HasSuffix(codePrompt, wantOmissions) {
				t.Fatalf("code prompt leaked data or lost the scope contract: %s", codePrompt)
			}
			assertDataReviewHeader(t, dataPrompt, 1, 1, fixture.baseCommit, fixture.headCommit)
			if !strings.Contains(dataPrompt, "+\"rule_00001\"") || strings.Contains(dataPrompt, "+package code") {
				t.Fatal("data prompt lacks the catalog or includes code")
			}
			for _, request := range runner.requests {
				if request.Session != runner.prepared[0].Session || request.Access != agent.SessionAccessReadOnly || request.Runtime.ID != provider {
					t.Fatalf("request left the prepared read-only session: %+v", request.Session)
				}
			}
			fields := reviewPromptBoundRecordFields(t, fixture)
			segments := dataReviewRecordSegments(t, fields)
			if len(segments) != 1 || !reflect.DeepEqual(segments[0].Paths, []string{"rules/icms.json"}) || len(record.OmittedPaths) != 0 {
				t.Fatalf("data segments=%+v omissions=%+v", segments, record.OmittedPaths)
			}
			if _, ok := fields["segments"]; ok {
				t.Fatal("a single code prompt acquired segments")
			}
			diff := dataReviewPromptDiff(t, dataPrompt)
			if segments[0].DiffBytes != len(diff) {
				t.Fatalf("data bytes=%d want=%d", segments[0].DiffBytes, len(diff))
			}
			if (segments[0].EstimatedPromptTokens != nil) != (provider == "claude") {
				t.Fatalf("wrong provider token field: %+v", segments[0])
			}
			if provider == "claude" && (*segments[0].EstimatedPromptTokens != estimateReviewPromptTokens(dataPrompt) || *segments[0].EstimatedPromptTokens > claudeReviewPromptBudget) {
				t.Fatal("Claude data token estimate differs from the sent prompt")
			}
			var findings []map[string]json.RawMessage
			if err := json.Unmarshal(fields["findingItems"], &findings); err != nil {
				t.Fatal(err)
			}
			if len(findings) != 3 {
				t.Fatalf("findings=%+v", record.FindingItems)
			}
			for i, finding := range record.FindingItems {
				if finding.ID != fmt.Sprintf("F%d", i+1) || finding.Segment != 0 || reviewFindingDismissedByValidation(finding) != (i == 2) {
					t.Fatalf("finding %d=%+v", i, finding)
				}
				if _, ok := findings[i]["segment"]; ok {
					t.Fatal("single code or data finding has a code-segment label")
				}
				got, ok := findings[i]["dataSegment"]
				if (i > 0 && string(got) != "1") || (i == 0 && ok) {
					t.Fatalf("data label on finding %d: %s", i, got)
				}
			}
			assertDataReviewAnswer(t, record, "--- Review ---\n"+codeAnswer+"\n--- Data Review Segment 1 of 1 ---\n"+dataAnswer+"\n")
			t.Run("code segments", func(t *testing.T) {
				assertDataReviewAfterCodeSegments(t, provider)
			})
			if provider == "claude" {
				t.Run("data uses the token budget", func(t *testing.T) {
					runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
						{result: agent.ExecuteResult{Message: "No findings."}},
						{result: agent.ExecuteResult{Message: "No findings."}},
					}}}
					fixture := reviewDataCommandFixture(t, provider, runner, "  data_paths: ['rules/']\n")
					commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(15000))
					code, record, stderr := fixture.run(t)
					if code != exitOK || len(record.DataSegments) != 1 || len(runner.requests) != 2 {
						t.Fatalf("exit=%d data segments=%d prompts=%d stderr=%q", code, len(record.DataSegments), len(runner.requests), stderr)
					}
					segment := record.DataSegments[0]
					if segment.DiffBytes <= reviewDiffBound || segment.EstimatedPromptTokens != estimateReviewPromptTokens(runner.requests[1].Prompt) || segment.EstimatedPromptTokens > claudeReviewPromptBudget {
						t.Fatalf("Claude data fit=%+v", segment)
					}
				})
			}
		})
	}
}

func TestReviewSplitsALargeDataFileAtLineBoundaries(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "No findings."}},
		{result: agent.ExecuteResult{Message: "Findings:\n- rules/icms.json:20000 Failure: belongs to the next segment."}},
		{result: agent.ExecuteResult{Message: "Findings:\n- rules/icms.json:20000 Failure: this rule has the wrong value."}},
	}}}
	fixture := reviewDataCommandFixture(t, "codex", runner, "  data_paths: ['rules/']\n")
	commitScopeFile(t, fixture, "code.go", "package code\n")
	fixture.headCommit = commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(24000))
	code, record, stderr := fixture.run(t)
	segments := dataReviewRecordSegments(t, reviewPromptBoundRecordFields(t, fixture))
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || len(segments) != 2 || len(runner.requests) != 3 || runner.prepareCalls != 1 {
		t.Fatalf("exit=%d outcome=%s segments=%+v prompts=%d stderr=%q", code, record.Outcome, segments, len(runner.requests), stderr)
	}
	lines := strings.Split(gittest.Run(t, fixture.repository, "show", fixture.headCommit+":rules/icms.json"), "\n")
	next := 1
	for i, request := range runner.requests[1:] {
		assertDataReviewHeader(t, request.Prompt, i+1, 2, fixture.baseCommit, fixture.headCommit)
		diff := dataReviewPromptDiff(t, request.Prompt)
		// The first cut retains its final newline; the full Git diff is trimmed.
		if i == 0 {
			diff += "\n"
		}
		if segments[i].DiffBytes != len(diff) || len(diff) > reviewDiffBound || !reflect.DeepEqual(segments[i].Paths, []string{"rules/icms.json"}) || request.Session != runner.prepared[0].Session {
			t.Fatalf("segment %d=%+v bytes=%d", i+1, segments[i], len(diff))
		}
		inHunk := false
		lineNumber := 0
		for _, line := range strings.Split(diff, "\n") {
			if strings.HasPrefix(line, "@@ ") {
				var oldStart, oldCount, count int
				if _, err := fmt.Sscanf(line, "@@ -%d,%d +%d,%d @@", &oldStart, &oldCount, &lineNumber, &count); err != nil {
					t.Fatalf("hunk %q: %v", line, err)
				}
				inHunk = true
				continue
			}
			if !inHunk || line == "" || line[0] == '\\' || line[0] == '-' {
				continue
			}
			if lineNumber != next || next > len(lines) || line[1:] != lines[next-1] {
				t.Fatalf("segment %d line %d (expected %d) differs from Git", i+1, lineNumber, next)
			}
			next++
			lineNumber++
		}
	}
	if next != 24003 || len(record.FindingItems) != 2 || !reviewFindingDismissedByValidation(record.FindingItems[0]) || record.FindingItems[1].Validation.Status != "stands" {
		t.Fatalf("next=%d findings=%+v", next, record.FindingItems)
	}
	var findings []struct {
		DataSegment int `json:"dataSegment"`
		Segment     int `json:"segment"`
	}
	if err := json.Unmarshal(reviewPromptBoundRecordFields(t, fixture)["findingItems"], &findings); err != nil {
		t.Fatal(err)
	}
	if findings[0].DataSegment != 1 || findings[1].DataSegment != 2 || findings[0].Segment != 0 || findings[1].Segment != 0 {
		t.Fatalf("finding segment labels=%+v", findings)
	}
	assertDataReviewAnswer(t, record, "--- Review ---\nNo findings.\n--- Data Review Segment 1 of 2 ---\n"+runner.results[1].result.Message+"\n--- Data Review Segment 2 of 2 ---\n"+runner.results[2].result.Message+"\n")
	for _, blockedAt := range []int{0, 1} {
		t.Run(fmt.Sprintf("blocked answer %d", blockedAt), func(t *testing.T) {
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
				{result: agent.ExecuteResult{Message: "No findings."}},
				{result: agent.ExecuteResult{Message: "No findings."}},
				{result: agent.ExecuteResult{Message: "No findings."}},
			}}}
			runner.results[blockedAt].result.Message = "unclassifiable answer"
			fixture := reviewDataCommandFixture(t, "codex", runner, "  data_paths: ['rules/']\n")
			commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(24000))
			code, record, stderr := fixture.run(t)
			if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || len(runner.requests) != blockedAt+1 || runner.endCalls != 1 {
				t.Fatalf("exit=%d outcome=%s prompts=%d ends=%d stderr=%q", code, record.Outcome, len(runner.requests), runner.endCalls, stderr)
			}
			want := "--- Review ---\n"
			if blockedAt == 1 {
				want += "No findings.\n--- Data Review Segment 1 of 2 ---\n"
			}
			assertDataReviewAnswer(t, record, want+"unclassifiable answer\n")
		})
	}
}

func TestReviewBlocksDataNeedingMoreThanFourDataSegments(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{}
	fixture := reviewDataCommandFixture(t, "codex", runner, "  data_paths: ['rules/']\n")
	commitScopeFile(t, fixture, "code.go", "package code\n")
	for i := range 5 {
		fixture.headCommit = commitScopeFile(t, fixture, fmt.Sprintf("rules/catalog%d.json", i), reviewDataCatalog(14000))
	}
	diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", fixture.baseCommit, fixture.headCommit, "--", ":(literal)rules/")
	if err != nil {
		t.Fatal(err)
	}
	code, record, stderr := fixture.run(t)
	want := fmt.Sprintf("review data too large: %d bytes of pre_pr_review.data_paths needs 5 data review segments; a review allows at most %d", len(diff), reviewSegmentLimit)
	want += reviewSizeRefusalTestSuffix(t, fixture, fixture.baseCommit, fixture.headCommit, []reviewOmittedPath{{Path: "code.go"}})
	assertBlockedReviewCommand(t, code, record, stderr, want)
	if record.Reason != want || stderr != "roundfix: review blocked: "+want+"\n" || record.AnswerPath != "" || len(record.OmittedPaths) != 0 {
		t.Fatalf("reason=%q omissions=%+v stderr=%q", record.Reason, record.OmittedPaths, stderr)
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 || runner.endCalls != 0 {
		t.Fatal("oversized data reached a provider")
	}
}

func TestReviewOmissionWinsOverADataPath(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}}
	fixture := reviewDataCommandFixture(t, "codex", runner, "  data_paths: ['rules/']\n  omit_paths: ['rules/icms.json']\n")
	commitScopeFile(t, fixture, "code.go", "package code\n")
	commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(9000))
	code, record, stderr := fixture.run(t)
	want := []reviewOmittedPath{{Path: "rules/icms.json", Reason: "omit-path"}}
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || !reflect.DeepEqual(record.OmittedPaths, want) || len(runner.requests) != 1 {
		t.Fatalf("exit=%d record=%+v prompts=%d stderr=%q", code, record, len(runner.requests), stderr)
	}
	if _, ok := reviewPromptBoundRecordFields(t, fixture)["dataSegments"]; ok {
		t.Fatal("omitted path acquired data segments")
	}
	if strings.Contains(runner.requests[0].Prompt, "\"rule_") || strings.Contains(runner.requests[0].Prompt, "Reviewed apart,") || !strings.HasSuffix(runner.requests[0].Prompt, ", 1 of pre_pr_review.omit_paths.\nrules/icms.json (omit-path)\n") {
		t.Fatal("omission prompt contract changed")
	}
	assertDataReviewAnswer(t, record, "No findings.")
}

func TestReviewReadsADataPathApartInRoundTwo(t *testing.T) {
	t.Parallel()
	runner := newLineageRunner("Findings:\n- code.go:1 Failure: code loses a rule.", "No findings.", "No findings.", "No findings.")
	fixture := reviewDataCommandFixture(t, "codex", runner, "  data_paths: ['rules/']\n")
	commitScopeFile(t, fixture, "code.go", "package code\n")
	catalog := reviewDataCatalog(9000)
	commitScopeFile(t, fixture, "rules/icms.json", catalog)
	// Head declarations cannot clear the immutable base's policy, even next round.
	config, err := os.ReadFile(filepath.Join(fixture.repository, ".roundfixrc.yml"))
	if err != nil {
		t.Fatal(err)
	}
	commitScopeFile(t, fixture, ".roundfixrc.yml", strings.Replace(string(config), "data_paths: ['rules/']", "data_paths: []", 1))
	code, first, stderr := fixture.run(t)
	if code != exitRunFailed || len(runner.prompts) != 2 || first.FindingItems[0].Segment != 0 {
		t.Fatalf("first exit=%d prompts=%d stderr=%q", code, len(runner.prompts), stderr)
	}
	changed := strings.Replace(catalog, "\"rule_04500\": {\"rate\": 18, \"description\": \"authored value\"}", "\"rule_04500\": {\"rate\": 19, \"description\": \"corrected value\"}", 1)
	if changed == catalog {
		t.Fatal("catalog change did not match")
	}
	fixture.headCommit = commitScopeFile(t, fixture, "rules/icms.json", changed)
	code, second, stderr := fixture.run(t)
	if code != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || second.Lineage.PreviousHead != first.HeadCommit || len(runner.prompts) != 4 {
		t.Fatalf("second exit=%d lineage=%+v prompts=%d stderr=%q", code, second.Lineage, len(runner.prompts), stderr)
	}
	codePrompt, dataPrompt := runner.prompts[2], runner.prompts[3]
	if !strings.Contains(codePrompt, "--- BEGIN ROUND 2 DELTA DIFF ---") || strings.Contains(codePrompt, "\"rule_") || strings.Contains(codePrompt, "diff --git a/rules/") {
		t.Fatal("round two code prompt lacks the delta envelope or leaks the catalog")
	}
	assertDataReviewHeader(t, dataPrompt, 1, 1, first.HeadCommit, fixture.headCommit)
	diff := dataReviewPromptDiff(t, dataPrompt)
	wantDiff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", first.HeadCommit, fixture.headCommit, "--", ":(literal)rules/icms.json")
	if err != nil {
		t.Fatal(err)
	}
	if diff != wantDiff || !strings.Contains(diff, "\"rate\": 19") || strings.Contains(diff, "\"rule_00001\"") {
		t.Fatal("data prompt is not exactly the changed entry's Git delta")
	}
	if len(second.Segments) != 0 || len(second.OmittedPaths) != 0 || len(dataReviewRecordSegments(t, reviewPromptBoundRecordFields(t, fixture))) != 1 || runner.request.Session.Name != first.Lineage.Session {
		t.Fatalf("round two record=%+v session=%s", second, runner.request.Session.Name)
	}
	assertDataReviewAnswer(t, second, "--- Review ---\nNo findings.\n--- Data Review Segment 1 of 1 ---\nNo findings.\n")
}

func reviewDataCommandFixture(t *testing.T, provider string, runner agent.Runner, declarations string) reviewCommandFixture {
	t.Helper()
	fixture := newReviewCommandFixture(t, provider, runner)
	writeReviewCommandProfileConfig(t, fixture.repository, provider, fixture.artifactDir, provider, provider)
	file := filepath.Join(fixture.repository, ".roundfixrc.yml")
	config, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Replace(string(config), "  provider: "+provider+"\n", "  provider: "+provider+"\n"+declarations, 1)
	fixture.baseCommit = commitScopeFile(t, fixture, ".roundfixrc.yml", body)
	gittest.Run(t, fixture.repository, "branch", "-f", "main", fixture.baseCommit)
	return fixture
}

func reviewDataCatalog(count int) string {
	var body strings.Builder
	body.WriteString("{\n")
	for i := 1; i <= count; i++ {
		fmt.Fprintf(&body, "\"rule_%05d\": {\"rate\": 18, \"description\": \"authored value\"}", i)
		if i < count {
			body.WriteByte(',')
		}
		body.WriteByte('\n')
	}
	body.WriteString("}\n")
	return body.String()
}

type dataReviewRecordedSegment struct {
	Paths                 []string `json:"paths"`
	DiffBytes             int      `json:"diffBytes"`
	EstimatedPromptTokens *int     `json:"estimatedPromptTokens"`
}

func dataReviewRecordSegments(t *testing.T, fields map[string]json.RawMessage) []dataReviewRecordedSegment {
	t.Helper()
	var segments []dataReviewRecordedSegment
	if err := json.Unmarshal(fields["dataSegments"], &segments); err != nil {
		t.Fatalf("decode dataSegments %q: %v", fields["dataSegments"], err)
	}
	return segments
}

func assertDataReviewHeader(t *testing.T, prompt string, index, total int, from, head string) {
	t.Helper()
	want := fmt.Sprintf("Data Review Segment %d of %d: pre_pr_review.data_paths changes %s..%s of the candidate; answer for this segment only.\n", index, total, from, head) +
		"These files are data the project writes by hand and declares in pre_pr_review.data_paths, such as rule catalogs. Check each changed value against the Spec contexts, the code that reads it and the rest of the file; report a wrong, missing, duplicated or contradictory value.\n"
	if !strings.HasPrefix(prompt, want) {
		t.Fatalf("data prompt lacks %q", want)
	}
}

func dataReviewPromptDiff(t *testing.T, prompt string) string {
	t.Helper()
	_, body, ok := strings.Cut(prompt, "--- BEGIN CANDIDATE DIFF ---\n")
	if !ok {
		t.Fatal("data prompt lacks a diff")
	}
	diff, _, ok := strings.Cut(body, "--- END CANDIDATE DIFF ---")
	if !ok {
		t.Fatal("data prompt lacks the diff closing marker")
	}
	return strings.TrimSuffix(diff, "\n")
}

func assertDataReviewAnswer(t *testing.T, record reviewRecord, want string) {
	t.Helper()
	answer, err := os.ReadFile(record.AnswerPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(answer) != want {
		t.Fatalf("answer=%q, want %q", answer, want)
	}
}

func assertDataReviewAfterCodeSegments(t *testing.T, provider string) {
	t.Helper()
	answers := []string{
		"Findings:\n- part1.txt:1 Failure: the first code segment breaks.",
		"Findings:\n- part2.txt:1 Failure: the second code segment breaks.",
		"Findings:\n- rules/icms.json:10 Failure: the data segment breaks.",
	}
	runner := &segmentReviewRunner{}
	for _, answer := range answers {
		runner.results = append(runner.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: answer}})
	}
	fixture := reviewDataCommandFixture(t, provider, runner, "  data_paths: ['rules/']\n")
	heads := []string{fixture.baseCommit}
	for i := 1; i <= 2; i++ {
		heads = append(heads, commitScopeFile(t, fixture, fmt.Sprintf("part%d.txt", i), strings.Repeat("x", 700000)+"\n"))
	}
	fixture.headCommit = commitScopeFile(t, fixture, "rules/icms.json", reviewDataCatalog(9000))
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || len(record.Segments) != 2 || len(record.DataSegments) != 1 || len(runner.requests) != 3 || runner.prepareCalls != 1 || len(record.FindingItems) != 3 {
		t.Fatalf("exit=%d code segments=%d data segments=%d prompts=%d stderr=%q", code, len(record.Segments), len(record.DataSegments), len(runner.requests), stderr)
	}
	for i, request := range runner.requests {
		if request.Session != runner.prepared[0].Session {
			t.Fatal("mixed segments left their prepared session")
		}
		finding := record.FindingItems[i]
		if finding.ID != fmt.Sprintf("F%d", i+1) || finding.Validation.Status != "stands" {
			t.Fatalf("finding=%+v", finding)
		}
		if i < 2 {
			header := fmt.Sprintf("Review Segment %d of 2: commits %s..", i+1, heads[i])
			if !strings.HasPrefix(request.Prompt, header) || strings.Contains(request.Prompt, "\"rule_") || finding.Segment != i+1 || finding.DataSegment != 0 {
				t.Fatalf("code segment %d header or label changed: %+v", i+1, finding)
			}
		} else {
			assertDataReviewHeader(t, request.Prompt, 1, 1, fixture.baseCommit, fixture.headCommit)
			if finding.Segment != 0 || finding.DataSegment != 1 {
				t.Fatalf("data label=%+v", finding)
			}
		}
	}
	assertDataReviewAnswer(t, record, "--- Review Segment 1 of 2 ---\n"+answers[0]+"\n--- Review Segment 2 of 2 ---\n"+answers[1]+"\n--- Data Review Segment 1 of 1 ---\n"+answers[2]+"\n")
}
