// IN: Check and CheckStage against disposable Git repositories.
// OUT: real repository history, provider calls, and Daemon settlement.
// Invariant: declarations, bold phrases and completed glossary-writing Tasks agree.
package speccheck

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/testfixture"
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

// glossaryDefinitionFormCases covers accepted definitions and nearby prose forms.
var glossaryDefinitionFormCases = []struct {
	name, line string
	defined    bool
}{
	{"colon", "**New Term**: d", true},
	{"dash list colon", "- **New Term**: d", true},
	{"star list colon", "* **New Term**: d", true},
	{"plus list colon", "+ **New Term**: d", true},
	{"em dash", "**New Term** — d", true},
	{"dash list em dash", "- **New Term** — d", true},
	{"em dash at line end", "- **New Term** —", true},
	{"heading", "### New Term", false},
	{"indented list item", "  - **New Term**: d", false},
	{"tab indented list item", "\t- **New Term** — d", false},
	{"ordered list item", "1. **New Term**: d", false},
	{"ascii hyphen", "- **New Term** - d", false},
	{"double hyphen", "- **New Term** -- d", false},
	{"en dash", "- **New Term** – d", false},
	{"unspaced em dash", "- **New Term**—d", false},
	{"unspaced list marker", "-**New Term**: d", false},
	{"bold prose", "- **New Term** is prose", false},
	{"bold sentence", "- **New Term.** Prose", false},
	{"inside sentence", "See **New Term**: d", false},
}

func TestGlossaryDefinitionFormsAreRead(t *testing.T) {
	t.Parallel()
	for _, tc := range glossaryDefinitionFormCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("completed binding", func(t *testing.T) {
				root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
				writeReceiptFixture(t, root, "CONTEXT.md", tc.line+"\n")
				glossaryTask(t, root, "completed", "- interface: CONTEXT.md", "**New Term**")
				result := glossaryCheck(t, root, StageAll)
				glossaryAssert(t, result, CodeGlossaryUnplanned)
				if tc.defined {
					glossaryAssert(t, result, CodeGlossaryMissing)
				} else {
					glossaryAssert(t, result, CodeGlossaryMissing, 5)
				}
			})
			t.Run("bold coverage", func(t *testing.T) {
				root := glossaryRepo(t, "## Glossary\nNone.\n## Story\n**New Term**\n")
				writeReceiptFixture(t, root, "CONTEXT.md", tc.line+"\n")
				result := glossaryCheck(t, root, StagePRD)
				if tc.defined {
					glossaryAssert(t, result, CodeGlossaryUndeclared)
				} else {
					glossaryAssert(t, result, CodeGlossaryUndeclared, 4)
				}
			})
		})
	}
}

func TestGlossaryVerificationCommandAgreesWithTheDefinitionReader(t *testing.T) {
	t.Parallel()
	grep := testfixture.RequireTool(t, "grep")
	type commandCase struct {
		name, term, line, path string
		defined                bool
	}
	var cases []commandCase
	for _, tc := range glossaryDefinitionFormCases {
		cases = append(cases, commandCase{tc.name, "New Term", tc.line, "CONTEXT.md", tc.defined})
	}
	cases = append(cases,
		commandCase{"metacharacter term", "New.Term", "**NewXTerm**: cites **New.Term**", "CONTEXT.md", false},
		commandCase{"apostrophe term", "Owner's Term", "- **Owner's Term** — d", "CONTEXT.md", true},
		commandCase{"later mention", "New Term", "**Other Term**: cites **New Term**: d", "CONTEXT.md", false},
		commandCase{"quoted glossary path", "New Term", "**New Term**: d", "owner's glossary/CONTEXT.md", true},
	)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeReceiptFixture(t, root, tc.path, tc.line+"\n")
			if tc.path != "CONTEXT.md" {
				writeReceiptFixture(t, root, "CONTEXT-MAP.md", "[Glossary](<"+tc.path+">)\n")
			}
			_, terms, err := glossaryFiles(root)
			if err != nil {
				t.Fatal(err)
			}
			if got := terms[glossaryNormalize(tc.term)]; got != tc.defined {
				t.Fatalf("reader defined = %v, want %v for %q", got, tc.defined, tc.line)
			}
			command := GlossaryVerificationCommand(tc.term, tc.path)
			cmd := exec.CommandContext(t.Context(), "sh", "-c", command)
			cmd.Dir = root
			cmd.Env = []string{"PATH=" + filepath.Dir(grep), "LC_ALL=C", "HOME=" + root, "TMPDIR=" + root}
			output, err := cmd.CombinedOutput()
			var exitErr *exec.ExitError
			if err != nil && (!errors.As(err, &exitErr) || exitErr.ExitCode() != 1) {
				t.Fatalf("command error: %v; output=%q; command=%q", err, output, command)
			}
			if got := err == nil; got != tc.defined {
				t.Fatalf("command accepted = %v, want %v; output=%q; command=%q", got, tc.defined, output, command)
			}
		})
	}
}

func TestGlossaryFindingsNameTheDefinitionForms(t *testing.T) {
	t.Parallel()
	t.Run("missing", func(t *testing.T) {
		root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
		glossaryTask(t, root, "completed", "- interface: CONTEXT.md", "**New Term**")
		result := glossaryCheck(t, root, StageAll)
		glossaryAssert(t, result, CodeGlossaryMissing, 5)
		findings := receiptCodeFindings(result, CodeGlossaryMissing)
		if want := "Define the declared term in the glossary as " + GlossaryDefinitionForms + "."; findings[0].Fix != want {
			t.Fatalf("missing fix = %q, want %q", findings[0].Fix, want)
		}
	})
	for _, path := range []string{"CONTEXT.md", "GLOSSARY.md", "a/CONTEXT.md"} {
		t.Run(path, func(t *testing.T) {
			root := glossaryRepo(t, "## Glossary\n- adds: **New Term**\n")
			writeReceiptFixture(t, root, "GLOSSARY.md", "# Glossary\n")
			writeReceiptFixture(t, root, "CONTEXT-MAP.md", "[A](a/CONTEXT.md)\n[Z](z/GLOSSARY.md)\n")
			writeReceiptFixture(t, root, "a/CONTEXT.md", "# Glossary\n")
			writeReceiptFixture(t, root, "z/GLOSSARY.md", "# Glossary\n")
			if path != "CONTEXT.md" {
				if err := os.Remove(filepath.Join(root, "CONTEXT.md")); err != nil {
					t.Fatal(err)
				}
			}
			if path == "a/CONTEXT.md" {
				if err := os.Remove(filepath.Join(root, "GLOSSARY.md")); err != nil {
					t.Fatal(err)
				}
			}
			glossaryTask(t, root, "pending", "- interface: "+path, "no term binding")
			result := glossaryCheck(t, root, StageAll)
			glossaryAssert(t, result, CodeGlossaryUnplanned, 5)
			findings := receiptCodeFindings(result, CodeGlossaryUnplanned)
			command := GlossaryVerificationCommand("New Term", path)
			if want := "Add a non-QA Task declaring the glossary under interface: or creates: and naming **New Term** in ## Verification, such as: " + command; findings[0].Fix != want {
				t.Fatalf("unplanned fix = %q, want %q", findings[0].Fix, want)
			}
			if findings := WrapFragileVerification(spec.Task{File: "task_01.md", Verification: []string{command}}); len(findings) != 0 {
				t.Fatalf("rendered command is wrap-fragile: %#v", findings)
			}
			glossaryTask(t, root, "pending", "- interface: "+path, command)
			glossaryAssert(t, glossaryCheck(t, root, StageAll), CodeGlossaryUnplanned)
		})
	}
}
