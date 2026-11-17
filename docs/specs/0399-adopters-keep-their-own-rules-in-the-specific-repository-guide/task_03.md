---
task: task_03
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, skill-creator, evidence-gate, knowledge-workspace, roundfix, systematic-debugging.

### Verification feedback repair

The Daemon's attempt-1 diagnostic artifact was inspected at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261117T004735Z_ff97830eca75451a/verification/batch-003-attempt-1.log`.
It identifies four existing Doctor tests whose whole-output or ordered-name
expectations omit the new line. The subsequent Verification Feedback directs
repair of the affected tests in this Task's slice, resolving the earlier
request for authority to update those expectations beyond the two counts.

Only the corresponding expectations in the declared
`internal/cli/doctor_test.go` interface changed during this repair:

- The three profile-readiness cases and the missing-repository-root case now
  expect the skipped line after skills and before residue.
- The adapter test's ordered-name list includes `HealthCheckBaselineRules`
  after skills, preserving the complete ordering assertion.
- The real-repository state-preservation test expects the unread warning:
  its existing fixture has a partial manifest accepted by the skill reader
  but invalid for Baseline parsing. The fixture itself is unchanged.

Exit-code, adapter, runtime, profile-proof and state-preservation assertions
remain intact. No production behavior, fixture, skill, record or guide was
changed in this feedback repair. No terminal Task verdict is claimed.

### Implementation and acceptance evidence

- **Plan/current/applied reporting:** the result gains the required pointer
  field and computes the report at result emission, only after manifest
  parsing succeeded. This also covers decision refusals and later operation
  failures. Applied results read the bytes after apply. Text emits exact
  ranges, managed identities, first lines and destinations between the
  existing unrecorded-region and retention sections. Extension advice depends
  on the recorded extension decision. A detector error adds one warning and
  an empty array without changing the command's existing state or exit.
  All nine required `TestBaselineUpdateNamesMisplacedRepositoryRules`
  subtests pass. They check whole section lines, whole JSON arrays, the
  direct planner's digest, unchanged preview bytes, shifted applied ranges
  and unchanged exit behavior with the outside passage removed.
- **No manifest and earlier failures:** `no_manifest_omits_the_field` checks
  text and JSON omission. The additional
  `TestBaselineUpdateMisplacedRuleReportAcrossFailures` passes for an empty
  array, invalid-manifest omission, a declaration refusal after a valid
  manifest, and a missing carrier. The missing-carrier case checks the whole
  detector warning line exactly once and preserves the existing
  `action_required` exit 3; the new test's initial expectation of exit 2 was
  corrected after reading the actual planner outcome.
- **Doctor:** the new check produces the four required forms and both DR
  codes, including singular/plural wording and extension-off advice. It is
  appended after skill, user-skill and authoring-model reports, before
  residue, and only returns ok/warn/skipped. All seven required
  `TestDoctorReportsMisplacedRepositoryRules` subtests pass, asserting whole
  lines, position, repository byte preservation and both advisory exit 0
  and an independent failure's exit 1. The existing repository-skill table
  test passes with its count change. After the feedback repair, the adapter
  table, profile-readiness table, missing-root case and real-repository
  state-preservation case also pass their focused checks.
- **Surface Transcript 1:**
  `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` passes through
  `runCLI` against this repository with a temporary home and no confirmation.
  Its first eight stdout lines match the transcript with the embedded
  catalog digest, more lines follow, stderr is empty and exit is 0.
- **Documentation and skills:** both usage sentences, both command guides,
  the Roundfix baseline/setup references and the managed-refresh paragraph
  of setup-context-driven describe the public surfaces and approval before
  an Agent moves a passage. The required move sentence matches all four
  canonical/mirror files. Generated versions are Roundfix `0.1.248` and
  setup-context-driven `0.0.9`; both copies match byte for byte. The owned
  version record and Behavior Surface Record were generated, never edited
  by hand.

### Focused checks and generation

- Feedback repair check: `go test ./internal/cli -count=1 -run '^(TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes|TestRunDoctorMissingRepositoryRoot|TestRunDoctorRealRepositoryCheckDoesNotMutateState|TestRunDoctorRepositorySkillReadiness|TestDoctorReportsMisplacedRepositoryRules)$' -v`
  exited 0 after the repair. All six selected tests and their subtests
  passed, including all four test groups named in the Daemon diagnostic.
  `git diff --check` also exited 0.
- Red starting point: `go test ./internal/cli -run 'TestBaselineUpdateNamesMisplacedRepositoryRules/plan_json$' -count=1`
  rejected the absent `misplacedRepositoryRules` array before the implementation.
- Final focused CLI check: `go test ./internal/cli -run '^(TestBaselineUpdateNamesMisplacedRepositoryRules|TestDoctorReportsMisplacedRepositoryRules|TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent|TestBaselineUpdateMisplacedRuleReportAcrossFailures|TestRunDoctorRepositorySkillReadiness)$' -count=1 -v`
  exited 0; every selected test and subtest passed.
- The earlier focused check also selected
  `TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes`; its three
  subtests exposed the expected-name-list mismatch subsequently repaired
  under the Verification Feedback instruction.
- `make skills-sync` regenerated the mirrors after source edits and version
  recording.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and generated both skill version fields and their record.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 after the final skill-reference edit.
- `make baseline-digests` exited 0 with `changed: false`; no Baseline asset or
  derived Baseline path changed.
- `go test ./skills -run '^(TestAuthorialSkillSync|TestEveryOwnedSkillVersionIsRecorded)$' -count=1`
  exited 0.
- `go test -tags docscontract ./internal/docscontract -run '^(TestBaselineDocumentationContract|TestTheSkillCoverageMapIsCurrent|TestEveryCommandIsNamedInTheRoundfixSkill|TestCommandCoverageFollowsTheRoundfixReferenceIndex|TestSkillCoveragePathsRequireOwnedFilesAndMatchingSources)$' -count=1`
  exited 0 after the final reference edit.
- Scope inspection found exactly 20 changed paths, all declared by this Task;
  four canonical/mirror pairs match byte for byte. `git diff --check` is clean.

The initial worktree had only the Daemon-owned status change in this Task.
Task status, the Task Graph, other Task files, the forbidden paths and QA
settlement sections are unchanged by this implementation. No declared
Verification command, repository Verification, full suite, commit, push,
pull request or cross-project message was run. The Daemon owns Verification
and settlement; the configured sequence has not been rerun by this Agent.

## Carry-forward provenance

- Source Run: `run_20261117T004735Z_ff97830eca75451a`
- Source commit: `f9376c79472de9ab33eaa1021cd03373f16424d1`
