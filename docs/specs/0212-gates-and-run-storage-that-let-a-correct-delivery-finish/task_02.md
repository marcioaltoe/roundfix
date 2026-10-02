---
task: task_02
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: completed
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


## Result

Implemented the compiled-source exclusion before the import reads blobs or
compares destination paths. The existing pass selection, date/newest-report
refusals, report shape check and rollback remain in place; Evidence Snapshot
code is unchanged. The `prior_report` event now always has a `skipped` list of
`{path, reason}` entries and lists only imported paths under `files` (`null`
when no files were imported, preserving the no-prior-pass event contract).

Acceptance evidence from focused implementation checks:

1. `TestPriorQAPassSkipsCompiledSourceEvidence` exercises a recorded failed
   pass through real Git and the Run journal. It confirms the report, text
   evidence and `.go.txt` evidence are byte-identical after import, the
   `replay_test.go` file is absent, and the event names it with reason
   `compiled source` while excluding it from `files`.
2. `TestPriorQAPassSkipsEveryCompiledExtension` covers all ten extensions.
   `TestPriorQAPassSkipNeverReportsPathDiffers` repeats all ten with different
   local bytes, confirms those bytes survive and checks the imported event
   rather than a `path differs` refusal.
3. The focused prior-pass suite includes every existing `TestPriorQAPass…`
   test, recorded-commit selection, imported-head prompting and failed-pass
   row carry. A byte comparison against `HEAD` confirms
   `internal/daemon/qa_prior_pass_test.go` is unedited.
4. Added `### Runnable evidence` before `### QA settlement`, raised both
   qa-gate front-matter versions from `0.0.6` to `0.0.7`, synchronized the
   mirror and recorded the version. A Python byte comparison confirms the
   mirrors match and the settlement section equals its pre-edit snapshot.
   Focused skill checks confirm the version record and settlement guidance.

Commands and outcomes:

- `rtk proxy go test -count=1 ./internal/daemon -run '^TestPriorQAPassSkipsCompiledSourceEvidence$'`
  before the implementation: exit 1, reproducing `compiled source imported`.
- `rtk proxy go test -count=1 ./internal/daemon -run 'TestPriorQAPass|TestAPriorPass|TestQAPromptNamesTheImportedPassHead|TestQAGateCarriesARowFromAnUnintegratedFailedPass'`:
  final exit 0 (`ok roundfix/internal/daemon`, 9.001s). The first attempt
  encountered a sandbox Go-cache denial; an escalated attempt reached test
  PASS but the repository guard saw concurrent digest regeneration. After
  regeneration finished, the recorded final run passed without concurrent
  mutations.
- `rtk make skills-sync`: exit 0.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  exit 0; appended the qa-gate `0.0.7` version/digest record.
- `rtk make baseline-digests`: exit 0, `changed: false`; required derived
  artifacts already matched the canonical source after regeneration.
- `rtk proxy go test -count=1 ./skills -run 'TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations|TestAuthorialSkillSync'`:
  exit 0 (`ok roundfix/skills`, 0.177s).
- Python artifact assertions: matching skill mirrors, both versions `0.0.7`,
  byte-identical settlement section and unchanged existing import tests.
- `rtk proxy git -c core.fsmonitor=false diff --check` and focused
  `gofmt -l`: exit 0, no whitespace or formatting diagnostics.

Only this Task's declared paths were changed; the pre-existing Daemon status
change is preserved. No follow-ups were found. No declared Verification
command, full repository gate, commit, push or Pull Request operation was
run. Task status and settlement remain Daemon-owned.


### Verification feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T100707Z_cbb39973a4b07b9b/verification/batch-002-attempt-1.log`.
The configured gate exposed an event compatibility regression: the first
pass's `prior_report.files` must remain JSON `null` when there is no prior
pass. Initializing the new filtered event list as empty had changed it to
`[]`, causing both existing two-pass carry tests to refuse that event.

Preserved the established `null` representation when nothing is imported;
`skipped` still serializes as a list, and successful imports still list only
non-compiled files. Added
`TestPriorQAPassWithoutRecordedPassPreservesNullFiles` to this Task's new test
file to assert both fields at the no-prior-pass boundary. Existing prior-pass
and two-pass carry tests remain unedited.

Focused repair evidence:

- `rtk proxy go test -count=1 ./internal/daemon -run 'TestPriorQAPass|TestAPriorPass|TestQAPromptNamesTheImportedPassHead|TestQAGateCarriesARowFromAnUnintegratedFailedPass|TestTwoGatePasses'`:
  exit 0 (`ok roundfix/internal/daemon`, 10.578s), covering the new null-field
  regression, all compiled-source cases and both two-pass carry journeys.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.
- Byte comparisons against `HEAD`: both existing import and two-pass carry
  test files unchanged.

The repair changes only Task 02's implementation, new tests and Result.
Task status remains untouched; the configured gate and declared Verification
were not rerun. No commit, push or Pull Request operation was performed.

## Carry-forward provenance

- Source Run: `run_20261002T100707Z_cbb39973a4b07b9b`
- Source commit: `3294ea6bf3bba0e902e6a6015fb3e976653cf301`
