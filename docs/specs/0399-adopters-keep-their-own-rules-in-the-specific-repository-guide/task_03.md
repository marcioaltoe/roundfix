---
task: task_03
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
status: pending
type: backend
complexity: high
---

# Task 03: Baseline update and Doctor name each Misplaced Repository Rule, and the guides and skills describe them

## Overview

`roundfix baseline update` reports every **Misplaced Repository Rule** in
each result it gives after reading a valid Setup Manifest, in text and in
the JSON field `misplacedRepositoryRules`, outside the Plan Document and
without changing its exit code. `roundfix doctor` prints one
`baseline-rules` line that is `ok`, `warn` or `skipped` and never fails.
The usage text, the user guides of both commands, the Roundfix skill's
`baseline` and `setup` references and the `setup-context-driven` skill
describe them, and say that a person or a supervised Agent moves a passage,
never Roundfix.

## Requirements

1. MUST implement `_techspec.md` → API Contract 1 in
   `internal/cli/baseline_update.go`: the field
   `MisplacedRepositoryRules *[]baseline.MisplacedRepositoryRule` with the
   tag `json:"misplacedRepositoryRules,omitempty"`, set from
   `baseline.MisplacedRepositoryRules` after the plan, current or apply
   step, as `[]` when empty, whenever the command read a valid Setup
   Manifest, and nil on every earlier failure. An applied result computes it
   on the bytes after the apply. The text section is printed after
   `Unrecorded managed regions` and before `Retention evidence`, with the
   lines of the contract, the `Next for misplaced rules:` line only when
   the extension decision is not `true` and an entry exists. A detector
   error is no new failure: the result keeps its state and exit code and
   reports the error as one warning line of the existing `Warning:` form.
   The Plan Document, the Plan Digest and every exit code stay as they are
   for the same repository.
2. MUST implement `_techspec.md` → API Contract 2 in `internal/cli`:
   `HealthCheckBaselineRules` (`baseline-rules`) and
   `doctorBaselineRulesResult(repositoryRoot string) CheckResult`, with the
   four lines of the contract, word for word, and the codes
   `DR-BASELINE-RULE-MISPLACED` and `DR-BASELINE-RULES-UNREAD`. The `warn`
   line lists every entry as `<path>:<start>-<end>`, separated by `, `, and
   says `1 passage belongs` or `<n> passages belong`. Doctor appends the
   line after `skills:` and after `user-skills:` and `authoring-models:`
   when those are printed, before the residue lines. The line never changes
   Doctor's exit code.
3. MUST add to the usage text of `baseline update` and of `doctor` in
   `internal/cli/cli.go` one sentence each: that the result names each
   passage outside a Managed Region that belongs in
   `docs/agents/specific-repository.md` and moves none; and that the
   `baseline-rules:` line names those passages and never fails Doctor.
4. MUST add these tests in `internal/cli`, each top-level test calling
   `t.Parallel()` first, on temporary repositories and homes, with
   `runCLI` and stubbed command dependencies:
   - `TestBaselineUpdateNamesMisplacedRepositoryRules` in
     `internal/cli/baseline_update_misplaced_rules_test.go`, on a fixture
     with one passage after a guide's Managed Region and one
     `digest-mismatch` region, with the subtests `plan_text`, `plan_json`,
     `current_text`, `current_json`, `applied_json` (the line numbers are
     those of the bytes after the apply), `plan_digest_unchanged` (the
     result's `planDigest` equals the one `baseline` planning computes for
     the same fixture), `exit_code_unchanged` (the exit code equals that
     of the same fixture without the passage), `extension_decision_off`
     (the `Next for misplaced rules:` line) and `no_manifest_omits_the_field`
     (no `misplacedRepositoryRules` key and no section). Each asserts whole
     lines and the whole JSON array.
   - `TestDoctorReportsMisplacedRepositoryRules` in
     `internal/cli/doctor_baseline_rules_test.go`, with the subtests `ok`,
     `warn_one_passage`, `warn_two_passages`, `warn_extension_off`,
     `skipped`, `unread_invalid_manifest` and `after_skills_and_exit_code`,
     each asserting its whole line, its position and Doctor's exit code.
   - `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` in
     `internal/cli/baseline_update_misplaced_rules_test.go`: Surface
     Transcript 1 of `_techspec.md` run whole through `runCLI` on this
     repository's root, with plan only and no confirmation, asserting the
     first eight lines in order with `<digest>` as the only variation, that
     more lines follow, an empty standard error and exit `0`.
5. MUST change, in `internal/cli/doctor_test.go`, only the line counts of
   the two table tests that count fifteen lines, to the count with the
   `baseline-rules: skipped (...)` line their fake repository root now
   prints, and no other assertion.
6. MUST document both surfaces:
   - `docs/user-guide/commands/baseline.md`: the section
     `Misplaced repository rules: <n>`, the entry lines, the
     `Next for misplaced rules:` line, the JSON field
     `misplacedRepositoryRules` with its keys, that the field is outside the
     Plan Document and changes no exit code, that a block another tool
     delimits with `BEGIN:`/`END:` is not reported, and that the command
     moves nothing;
   - `docs/user-guide/commands/doctor.md`: the `baseline-rules:` line, its
     four forms, both `DR-` codes and that it never fails;
   - `.agents/skills/roundfix/references/baseline.md`: the field and the
     section, and this sentence: "When `roundfix baseline update` or
     `roundfix doctor` names a Misplaced Repository Rule, report each
     passage to the person and move it into
     `docs/agents/specific-repository.md` only with the person's approval.";
   - `.agents/skills/roundfix/references/setup.md`: the `baseline-rules:`
     line and its two codes;
   - `.agents/skills/setup-context-driven/SKILL.md`: in its managed-refresh
     paragraph, the same sentence as the `baseline` reference.

   MUST synchronize the mirrors under `skills/` with `make skills-sync`,
   record the owned skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   never by hand.
7. MUST NOT change `docs/agents/specific-repository.md`,
   `internal/cli/cli_test.go`, `docs/references/coverage-record.json`, the
   `### QA settlement` section of any skill, any Baseline asset or any
   other Governed Path than the five skill files it declares.

## Subtasks

- [ ] Add the result field and the text section to `baseline update`.
- [ ] Add the `baseline-rules` Doctor line after the skill lines.
- [ ] Add the usage sentences and the tests, with the Surface Transcript.
- [ ] Update the two Doctor line counts.
- [ ] Document both surfaces in the guides and skills, then sync the mirrors and record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] In plan, current and applied results, `baseline update` names every Misplaced Repository Rule by path, line range and first line with its destination, in text and JSON, and the Plan Digest and exit code equal those of the same repository without the passage.
- [ ] Without a Setup Manifest the JSON carries no `misplacedRepositoryRules` key and the text no section.
- [ ] `roundfix doctor` prints `baseline-rules: ok`, `warn` or `skipped` after the skill lines and exits as before.
- [ ] Surface Transcript 1 reproduces on this repository.
- [ ] The guides and both copies of the skills describe the field, the line and the move only with the person's approval; the mirrors, the version record and the Behavior Surface Record are current.

## Context

- instruction: `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md`
- instruction: `docs/adr/0220-doctor-checks-what-delivery-needs-and-only-a-definite-answer-fails.md`
- instruction: `internal/cli/baseline_update_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/doctor_test.go`
- creates: `internal/cli/baseline_update_misplaced_rules_test.go`
- creates: `internal/cli/doctor_baseline_rules_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `roundfix/baseline-update-result/v1` JSON result and the `baseline update` text — consumers: `./internal/cli`, `./internal/docscontract`
- The `roundfix doctor` output — consumers: `./internal/cli`, `./internal/docscontract`
- The Roundfix and `setup-context-driven` skills — consumers: `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestBaselineUpdateNamesMisplacedRepositoryRules|TestDoctorReportsMisplacedRepositoryRules|TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateNamesMisplacedRepositoryRules/plan_text TestBaselineUpdateNamesMisplacedRepositoryRules/plan_json TestBaselineUpdateNamesMisplacedRepositoryRules/current_text TestBaselineUpdateNamesMisplacedRepositoryRules/current_json TestBaselineUpdateNamesMisplacedRepositoryRules/applied_json TestBaselineUpdateNamesMisplacedRepositoryRules/plan_digest_unchanged TestBaselineUpdateNamesMisplacedRepositoryRules/exit_code_unchanged TestBaselineUpdateNamesMisplacedRepositoryRules/extension_decision_off TestBaselineUpdateNamesMisplacedRepositoryRules/no_manifest_omits_the_field TestDoctorReportsMisplacedRepositoryRules/ok TestDoctorReportsMisplacedRepositoryRules/warn_one_passage TestDoctorReportsMisplacedRepositoryRules/warn_two_passages TestDoctorReportsMisplacedRepositoryRules/warn_extension_off TestDoctorReportsMisplacedRepositoryRules/skipped TestDoctorReportsMisplacedRepositoryRules/unread_invalid_manifest TestDoctorReportsMisplacedRepositoryRules/after_skills_and_exit_code TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the tests do not exist, so the command fails; afterwards both surfaces, their states and the Surface Transcript pass by name, and the whole CLI suite, the Doctor line counts included, passes.

- `for f in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do grep -qF -- 'misplacedRepositoryRules' "$f" || { echo "missing field in $f" >&2; exit 1; }; done; for f in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do grep -qF -- 'DR-BASELINE-RULE-MISPLACED' "$f" && grep -qF -- 'DR-BASELINE-RULES-UNREAD' "$f" || { echo "missing Doctor codes in $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md .agents/skills/setup-context-driven/SKILL.md skills/setup-context-driven/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "names a Misplaced Repository Rule, report each passage to the person and move it into" || { echo "missing the move sentence in $f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "Misplaced repository rules:" || { echo 'baseline guide lacks the section' >&2; exit 1; }; go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide or skill names the field, the codes or the move; afterwards both guides and both copies of the skills do, the mirrors and the owned skill versions are current, and the documentation contracts and the Behavior Surface Record pass.

## References

- `_prd.md` → Goals (second and third bullets); User Experience; Core Features 3, 4 and 5; Success Metric 3; Success Metric 4; Non-Goals
- `_techspec.md` → Command surfaces; API Contract 1; API Contract 2; Surface Transcript 1; Data Models; Interfaces (invariants 6 and 8); Testing Approach (Task 3); Build Order 3
- ADR-0411, ADR-0220, ADR-0145
