---
task: task_02
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: pending
type: backend
complexity: medium
---

# Task 02: The failed-pass QA import leaves compiled source behind

## Overview

The next QA pass imports a failed pass's files byte for byte. In Spec 0203's
delivery that carried a Go test file a QA Agent wrote as evidence, which
failed the repository gate before any row ran and would have looped every
retry. This Task skips each compiled-source file in the import, records the
skip, and adds the qa-gate skill's rule that runnable evidence uses a
non-compiled extension.

## Requirements

1. MUST make `copyPriorQAPass` skip each file whose base name ends in `.go`,
   `.rs`, `.ts`, `.tsx`, `.mts`, `.cts`, `.js`, `.jsx`, `.mjs` or `.cjs`,
   neither writing nor comparing it, as "The QA import" states.
2. MUST add `skipped`, a list of `{path, reason}` with reason
   `compiled source`, to the `prior_report` journal event, and MUST keep
   `files` to what was imported (API Contract 3).
3. MUST NOT change which pass is chosen, the order of the existing refusals,
   the report shape check or the Evidence Snapshot.
4. MUST add the tests named in Verification to the new file
   `internal/daemon/qa_prior_pass_skip_test.go`, over the `priorFixture`
   helpers, and MUST leave every `TestPriorQAPass…` test unedited and
   passing.
5. MUST add to the qa-gate skill, under its own new heading
   `### Runnable evidence` outside `### QA settlement`, that runnable evidence
   is stored under a non-compiled extension such as `.go.txt`, never as a
   source file the repository gate compiles or formats, and that a later
   pass's import skips compiled source; MUST raise the skill's version by one
   patch level in both front-matter fields, run `make skills-sync`, and
   re-record the version.
6. MUST leave the `### QA settlement` section byte-identical.

## Subtasks

- [ ] Skip compiled-source files in the import.
- [ ] Record the skip in the journal event.
- [ ] Add the import tests.
- [ ] Add the qa-gate skill's runnable-evidence rule and raise its version.

## Acceptance Criteria

- [ ] A failed pass holding `qa/evidence/replay_test.go` and
      `qa/evidence/ledger.txt` imports the report and the text file, skips
      the Go file, and the event lists it under `skipped` with
      `compiled source`.
- [ ] Each listed extension is skipped, and a different file already at a
      skipped path never yields `path differs`.
- [ ] The existing import tests pass unedited.
- [ ] The qa-gate skill carries `### Runnable evidence`, its mirror equals it,
      its version moved past `0.0.6` and is recorded, and `### QA settlement`
      is unchanged.

## Context

- interface: `internal/daemon/qa_prior_pass.go`
- creates: `internal/daemon/qa_prior_pass_skip_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/daemon/qa_prior_pass_test.go`
- instruction: `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/references/2026-10-01-an-imported-qa-pass-can-carry-a-file-that-breaks-the-gate.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPriorQAPassSkipsCompiledSourceEvidence|TestPriorQAPassSkipsEveryCompiledExtension|TestPriorQAPassSkipNeverReportsPathDiffers)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPriorQAPassSkipsCompiledSourceEvidence TestPriorQAPassSkipsEveryCompiledExtension TestPriorQAPassSkipNeverReportsPathDiffers; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte|TestPriorQAPassRefusesADifferingPath|TestPriorQAPassKeepsIdenticalExistingEvidence|TestPriorQAPassSkipsCompiledSourceEvidence)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte TestPriorQAPassRefusesADifferingPath TestPriorQAPassKeepsIdenticalExistingEvidence TestPriorQAPassSkipsCompiledSourceEvidence; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing import tests run unedited beside the skip test, which does not exist before this Task.
- `tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "### Runnable evidence" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md "### Runnable evidence" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "non-compiled extension" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md "non-compiled extension" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- ".go.txt" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md ".go.txt" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "compiled source" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md "compiled source" >&2; exit 1; }; cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && ! grep -q '^version: 0.0.6$' .agents/skills/qa-gate/SKILL.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the qa-gate skill has no `### Runnable evidence` heading, so the command fails.

## References

- `_prd.md` → User Story 3; Core Features 4-5; Success Metric 3; Acceptance evidence
- `_techspec.md` → The QA import; API Contract 3; Testing Approach 2; Build Order 2
- ADR-0194; ADR-0210; ADR-0097
