// IN: Check and CheckStage against disposable Git repositories.
// OUT: real repository history, provider calls, and Daemon settlement.
// Invariant: declarations, bold phrases and completed glossary-writing Tasks agree.
package speccheck

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

const glossaryFixtureSlug = "0200-example"
const glossaryFixtureDir = "docs/specs/" + glossaryFixtureSlug

func glossaryRepo(t *testing.T, prd string) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeReceiptFixture(t, root, "CONTEXT.md", "**Existing Term**: definition\n")
	writeReceiptFixture(t, root, GlossaryGuidePath, "guide\n")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_prd.md", prd)
	return root
}

func glossaryCheck(t *testing.T, root string, stage Stage) Result {
	t.Helper()
	result, err := CheckStage(filepath.Join(root, "docs/specs"), root, glossaryFixtureSlug, stage)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func glossaryAssert(t *testing.T, result Result, code string, lines ...int) {
	t.Helper()
	findings := receiptCodeFindings(result, code)
	if len(findings) != len(lines) {
		t.Fatalf("%s findings = %#v, want lines %v", code, findings, lines)
	}
	for i, line := range lines {
		if findings[i].Severity != SeverityError || len(findings[i].Where) == 0 || findings[i].Where[0].Line != line || findings[i].Where[0].Path != glossaryFixtureDir+"/_prd.md" {
			t.Fatalf("finding = %#v, want PRD line %d", findings[i], line)
		}
	}
}

func glossaryTask(t *testing.T, root, status, context, verification string) {
	t.Helper()
	prdPath := filepath.Join(root, glossaryFixtureDir, "_prd.md")
	content, err := os.ReadFile(prdPath)
	if err != nil {
		t.Fatal(err)
	}
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_prd.md", "---\nstatus: active\n---\n"+string(content))
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_tasks.md", "---\nschema: spec-tasks/v1\nqa: declined\nqa_reason: detector fixture\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n---\n")
	writeReceiptFixture(t, root, glossaryFixtureDir+"/task_01.md", "---\ntask: task_01\nspec: "+glossaryFixtureSlug+"\nstatus: "+status+"\ntype: backend\n---\n# Task 01: Write the term\n\n## Context\n\n"+context+"\n\n## Verification\n\n- `true` — "+verification+"\n")
}

func TestGlossaryDeclarationIsRequiredAfterTheHorizon(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "# Example\n")
	glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared, 1)
	receiptCommit(t, root)
	glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared, 1)
}
func TestGlossaryMalformedEntryIsReported(t *testing.T) {
	t.Parallel()
	for _, entry := range []string{"- adds: New Term", "- not a term: **New Term**", "None.\n- adds: **New Term**", "- changes: ****", "- adds: **New Term** — "} {
		t.Run(entry, func(t *testing.T) {
			root := glossaryRepo(t, "## Glossary\n\n"+entry+"\n")
			glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared, 3)
		})
	}
}
func TestGlossaryUndeclaredBoldTermIsReported(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "## Glossary\nNone.\n## Story\n**New Term** and **New Term**\n**Existing Terms** and **Existing Termes**\n")
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_techspec.md", "## Design\n**New Term**\n")
	result := glossaryCheck(t, root, StageTechSpec)
	found := receiptCodeFindings(result, CodeGlossaryUndeclared)
	if len(found) != 2 || found[0].Where[0] != (Location{Path: glossaryFixtureDir + "/_prd.md", Line: 4}) || found[1].Where[0] != (Location{Path: glossaryFixtureDir + "/_techspec.md", Line: 2}) {
		t.Fatalf("findings: %#v", found)
	}
}
func TestGlossaryBoldRuleIgnoresLabelsDigitsAndFences(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "## Glossary\nNone.\n## Story\n**Label** **Phase 2** **Some (Term)** **Some Term.**\n```markdown\n**Fenced Term**\n```\n~~~\n**Other Fence**\n~~~\n**New Term**\n")
	glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared, 11)
}
func TestGlossaryNotATermCoversAPhrase(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "## Glossary\n- not a term: **New Term** — a label\n## Story\n**New Term**\n")
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryUndeclared)
}
func TestGlossaryAddedTermNeedsATaskThatWritesIt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, context, verification string
		want                        bool
	}{
		{"bound", "- interface: CONTEXT.md", "**New Term**", false},
		{"created", "- creates: GLOSSARY.md", "**New Term**", false},
		{"no path", "- instruction: CONTEXT.md", "**New Term**", true},
		{"no verification", "- interface: CONTEXT.md", "New Term", true},
		{"different case", "- interface: CONTEXT.md", "**new term**", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
			writeReceiptFixture(t, root, "GLOSSARY.md", "# Glossary\n")
			glossaryTask(t, root, "pending", tc.context, tc.verification)
			result := glossaryCheck(t, root, StageAll)
			if tc.want {
				glossaryAssert(t, result, CodeGlossaryUnplanned, 5)
			} else {
				glossaryAssert(t, result, CodeGlossaryUnplanned)
				glossaryAssert(t, result, CodeGlossaryMissing)
			}
		})
	}
	t.Run("QA Task does not bind", func(t *testing.T) {
		root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
		glossaryTask(t, root, "pending", "- interface: CONTEXT.md", "**New Term**")
		manifest, err := os.ReadFile(filepath.Join(root, glossaryFixtureDir, "_tasks.md"))
		if err != nil {
			t.Fatal(err)
		}
		writeReceiptFixture(t, root, glossaryFixtureDir+"/_tasks.md", strings.Replace(strings.Replace(string(manifest), "qa: declined", "qa: task_01", 1), "qa_reason: detector fixture\n", "", 1))
		task, err := os.ReadFile(filepath.Join(root, glossaryFixtureDir, "task_01.md"))
		if err != nil {
			t.Fatal(err)
		}
		writeReceiptFixture(t, root, glossaryFixtureDir+"/task_01.md", strings.Replace(string(task), "type: backend", "type: qa", 1))
		glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryUnplanned, 5)
	})
}
func TestGlossaryChangedTermMustExist(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "## Glossary\n- changes: **New Term**\n")
	glossaryTask(t, root, "pending", "- interface: CONTEXT.md", "**New Term**")
	result := glossaryCheck(t, root, StageAll)
	glossaryAssert(t, result, CodeGlossaryUnplanned, 5)
	if !strings.Contains(receiptCodeFindings(result, CodeGlossaryUnplanned)[0].Fix, "added") {
		t.Fatal("fix does not recommend adds")
	}
	writeReceiptFixture(t, root, "CONTEXT.md", "**new   TERM**: definition\n")
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryUnplanned)
}
func TestGlossaryCompletedTaskMustLeaveTheTermDefined(t *testing.T) {
	t.Parallel()
	root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
	glossaryTask(t, root, "completed", "- interface: CONTEXT.md", "**New Term**")
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryMissing, 5)

	manifestPath := filepath.Join(root, glossaryFixtureDir, "_tasks.md")
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_tasks.md", strings.Replace(string(manifest), "      needs: []\n", "      needs: []\n    - id: task_02\n      file: task_02.md\n      needs: []\n", 1))
	task, err := os.ReadFile(filepath.Join(root, glossaryFixtureDir, "task_01.md"))
	if err != nil {
		t.Fatal(err)
	}
	second := strings.Replace(string(task), "task: task_01", "task: task_02", 1)
	writeReceiptFixture(t, root, glossaryFixtureDir+"/task_02.md", strings.Replace(second, "status: completed", "status: pending", 1))
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryMissing)
	writeReceiptFixture(t, root, glossaryFixtureDir+"/task_02.md", second)
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryMissing, 5)
	writeReceiptFixture(t, root, "CONTEXT.md", "**New Term**: definition\n")
	glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryMissing)
}
func TestGlossaryReadsMappedAndRenamedGlossaries(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"CONTEXT.md", "GLOSSARY.md"} {
		t.Run(name, func(t *testing.T) {
			root := glossaryRepo(t, "## Glossary\n- adds: **Mapped Term**\n## Story\n**Root Term** **Mapped Terms**\n")
			writeReceiptFixture(t, root, "CONTEXT.md", "# no definitions\n")
			writeReceiptFixture(t, root, "GLOSSARY.md", "**Root Term**: definition\n")
			mapped := "domains/" + name
			writeReceiptFixture(t, root, "GLOSSARY-MAP.md", "[Domain](<"+mapped+">)\n[Reference][domain]\n[domain]: "+mapped+"\n[external](https://invalid/CONTEXT.md)\n[outside](../../CONTEXT.md)\n")
			writeReceiptFixture(t, root, mapped, "**Mapped Term**: definition\n")
			glossaryTask(t, root, "completed", "- interface: "+mapped, "**Mapped Term**")
			result := glossaryCheck(t, root, StageAll)
			for _, code := range []string{CodeGlossaryUndeclared, CodeGlossaryUnplanned, CodeGlossaryMissing} {
				glossaryAssert(t, result, code)
			}
		})
	}
}
func TestGlossaryStagesReportOnlyTheirCodes(t *testing.T) {
	t.Parallel()
	t.Run("missing inputs are recorded", func(t *testing.T) {
		root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
		result := glossaryCheck(t, root, StageAll)
		for _, code := range []string{CodeGlossaryUnplanned, CodeGlossaryMissing} {
			found := false
			for _, skip := range result.Skipped {
				if skip.Code == code && skip.Missing == glossaryFixtureDir+"/_tasks.md" {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing graph skip: %#v", result.Skipped)
			}
		}
		if err := os.Remove(filepath.Join(root, "CONTEXT.md")); err != nil {
			t.Fatal(err)
		}
		result = glossaryCheck(t, root, StageAll)
		for _, code := range []string{CodeGlossaryUndeclared, CodeGlossaryUnplanned, CodeGlossaryMissing} {
			glossaryAssert(t, result, code)
			found := false
			for _, skip := range result.Skipped {
				if skip.Code == code && skip.Missing == "CONTEXT.md or GLOSSARY.md" {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing glossary skip: %#v", result.Skipped)
			}
		}
	})
	root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_techspec.md", "## Story\n**Uncovered Term**\n")
	glossaryTask(t, root, "completed", "- interface: CONTEXT.md", "**New Term**")
	glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared)
	for _, stage := range []Stage{StagePRD, StageTechSpec} {
		glossaryAssert(t, glossaryCheck(t, root, stage), CodeGlossaryMissing)
	}
	tech := receiptCodeFindings(glossaryCheck(t, root, StageTechSpec), CodeGlossaryUndeclared)
	if len(tech) != 1 || tech[0].Where[0].Path != glossaryFixtureDir+"/_techspec.md" {
		t.Fatalf("tech findings: %#v", tech)
	}
	full := glossaryCheck(t, root, StageAll)
	glossaryAssert(t, full, CodeGlossaryMissing, 5)
	direct, err := Check(filepath.Join(root, "docs/specs"), root, glossaryFixtureSlug)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(full, direct) || !reflect.DeepEqual(full, glossaryCheck(t, root, StageTasks)) {
		t.Fatal("full and staged checks disagree")
	}
}
func TestGlossaryLegacySpecIsSkipped(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeReceiptFixture(t, root, "CONTEXT.md", "**Existing Term**: definition\n")
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_prd.md", "**Legacy Term**\n")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, GlossaryGuidePath, "guide\n")
	receiptCommit(t, root)
	result := glossaryCheck(t, root, StageAll)
	for _, code := range []string{CodeGlossaryUndeclared, CodeGlossaryUnplanned, CodeGlossaryMissing} {
		glossaryAssert(t, result, code)
		found := false
		for _, skip := range result.Skipped {
			if skip.Code == code && skip.Missing == "a PRD committed at or after "+GlossaryGuidePath {
				found = true
			}
		}
		if !found {
			t.Fatalf("skip absent: %#v", result.Skipped)
		}
	}
	writeReceiptFixture(t, root, "docs/specs/0201-new/_prd.md", "# New Spec\n")
	receiptCommit(t, root)
	newer, err := Check(filepath.Join(root, "docs/specs"), root, "0201-new")
	if err != nil {
		t.Fatal(err)
	}
	found := receiptCodeFindings(newer, CodeGlossaryUndeclared)
	if len(found) != 1 || found[0].Where[0] != (Location{Path: "docs/specs/0201-new/_prd.md", Line: 1}) {
		t.Fatalf("newer finding: %#v", found)
	}
	// A declaration opts even a legacy Spec in.
	writeReceiptFixture(t, root, glossaryFixtureDir+"/_prd.md", "## Glossary\nNone.\n## Story\n**Legacy Term**\n")
	glossaryAssert(t, glossaryCheck(t, root, StagePRD), CodeGlossaryUndeclared, 4)
}
