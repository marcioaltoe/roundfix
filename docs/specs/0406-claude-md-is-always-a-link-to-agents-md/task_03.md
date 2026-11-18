---
task: task_03
spec: 0406-claude-md-is-always-a-link-to-agents-md
status: completed
type: backend
complexity: high
---

# Task 03: Doctor and the update report the Claude Instruction Link, and the guides and skills describe it

## Overview

`roundfix doctor` gains one `baseline-claude-link` line that reads the root `CLAUDE.md` through the reader of Task 01 and never fails. `roundfix baseline update` already prints the plan's warnings. This Task proves the two new warnings through the command in text and JSON, and runs Surface Transcript 1 on this repository. The user guide pages for `baseline` and `doctor`, the Roundfix skill's `baseline`, `modes` and `setup` references and the `setup-context-driven` skill describe the rule, the warnings, the line, and the question the `--update` Skill Mode asks before it replaces `CLAUDE.md`.

## Requirements

1. MUST implement `_techspec.md` → API Contract 3 in `internal/cli`:
   - the constant `HealthCheckClaudeLink` (`baseline-claude-link`) in `health.go`;
   - `doctorClaudeLinkResult(ctx context.Context, repositoryRoot string) CheckResult` in `doctor.go`, with every line of the contract word for word and the codes `DR-BASELINE-CLAUDE-LINK`, `DR-BASELINE-CLAUDE-LINK-UNLINKED` and `DR-BASELINE-CLAUDE-LINK-UNREAD`;
   - `skipped` without a Setup Manifest, and otherwise the reader of Task 01 with `baseline.ExecGitRunner{}`;
   - passage locations as `<path>:<start>-<end>` separated by `, `, with `1 passage` or `<n> passages`.

   Doctor appends the line right after `baseline-rules:` and before the residue lines. The line never changes Doctor's exit code.
2. MUST add to the usage text of `doctor` in `internal/cli/cli.go` one sentence after the `baseline-rules:` sentence: that the `baseline-claude-link:` line reports whether `CLAUDE.md` is a symbolic link to `AGENTS.md` and never fails Doctor.
3. MUST add these tests in `internal/cli`, each top-level test calling `t.Parallel()` first, on temporary Git repositories and homes, with `runCLI` and stubbed command dependencies:
   - `TestDoctorReportsTheClaudeInstructionLink` in `internal/cli/doctor_claude_link_test.go`, with the subtests `ok`, `absent`, `regular_file`, `link_elsewhere`, `file_with_passages`, `link_elsewhere_with_passages`, `agents_not_regular`, `unlinked_checkout`, `unread`, `skipped` and `after_baseline_rules_and_exit_code`. Each asserts its whole line, its position and Doctor's exit code.
   - `TestBaselineUpdateWarnsAboutTheClaudeInstructionLink` in `internal/cli/baseline_update_claude_link_test.go`. Its fixtures are adopted, one with a `CLAUDE.md` holding a heading and a bullet, one with an unlinked checkout. The subtests are `withheld_text`, `withheld_json`, `unlinked_checkout_text`, `unlinked_checkout_json` and `exit_code_unchanged`. Each asserts the whole `Warning:` line or JSON warning of API Contract 2, that `CLAUDE.md` is byte-identical afterwards, and that the exit code equals that of the same fixture with `CLAUDE.md` the link.
   - `TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink` in the same file. It runs Surface Transcript 1 of `_techspec.md` whole through `runCLI` on this repository's root, plan only and without confirmation, with `<digest>` and `...` as the only variations. It also asserts that no output line names `baseline.root.claude-link` and that `CLAUDE.md` is still a symbolic link to `AGENTS.md`.
4. MUST change, in `internal/cli/doctor_test.go`, only what the new line moves. That is the two table tests that count Doctor's lines, which gain one, the check-name list, which gains `HealthCheckClaudeLink` after `HealthCheckBaselineRules`, and the expected outputs that hold the `baseline-rules:` line, which gain the `baseline-claude-link: skipped (...)` or `ok` line their fixture prints. No other assertion changes.
5. MUST document the surfaces:
   - `docs/user-guide/commands/baseline.md`: the rule that `CLAUDE.md` is always a symbolic link to `AGENTS.md`, when a plan creates or replaces it, the two warnings with their messages, that a `CLAUDE.md` holding passages of its own is left until they move to `docs/agents/specific-repository.md`, and the unlinked checkout;
   - `docs/user-guide/commands/doctor.md`: the `baseline-claude-link:` line in the list after `baseline-rules:`, its forms, its three `DR-` codes and that it never fails;
   - `.agents/skills/roundfix/references/baseline.md`: the rule, the two warnings, and this sentence: "When a plan warns `baseline.root.claude-link.withheld`, report each passage to the person and move it into `docs/agents/specific-repository.md` only with the person's approval; the next update then replaces `CLAUDE.md` with the link.";
   - `.agents/skills/roundfix/references/modes.md`: in step 4 of the `--update` Skill Mode, ask first, in the same question as a deletion, when a file change replaces `CLAUDE.md`, that is a change with action `update` on that path: apply the whole plan, or stop;
   - `.agents/skills/roundfix/references/setup.md`: the `baseline-claude-link:` line, its place after `baseline-rules:` and its codes;
   - `.agents/skills/setup-context-driven/SKILL.md`: that `AGENTS.md` is the root instruction file and `CLAUDE.md` always a symbolic link to it, with no repository option, and the same sentence as the `baseline` reference.
6. MUST regenerate, never by hand:
   - the mirrors under `skills/`, with `make skills-sync`;
   - the owned skill versions, with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - the Behavior Surface Record, with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, any Baseline module, template, profile or golden, the Makefile, a workflow, `go.mod` or a lint configuration. Of the Governed Paths it changes only those it declares below, all listed in `_authorization.md`.

## Subtasks

- [ ] Add the Doctor line and its usage sentence.
- [ ] Add the Doctor, update and transcript tests, and adjust the Doctor tests the line moves.
- [ ] Document the line, the warnings and the rule in the user guide and the skills.
- [ ] Sync the mirrors and record the skill versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Doctor prints exactly one `baseline-claude-link` line in every form of the contract, right after `baseline-rules:`, and its exit code does not change.
- [ ] The update prints both warnings in text and JSON, leaves `CLAUDE.md` untouched, and keeps its exit code.
- [ ] This repository's update matches Surface Transcript 1 and plans nothing for `CLAUDE.md`.
- [ ] The guides and skills state the rule, the warnings, the line and the `--update` question, and the mirrors, versions and records are current.

## Context

- instruction: `docs/adr/0432-claude-md-is-always-a-symbolic-link-to-agents-md.md`
- instruction: `docs/adr/0433-a-checkout-without-symbolic-links-keeps-the-link-git-records-for-claude-md.md`
- instruction: `internal/cli/doctor_baseline_rules_test.go`
- instruction: `internal/cli/baseline_update_misplaced_rules_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/doctor_test.go`
- creates: `internal/cli/doctor_claude_link_test.go`
- creates: `internal/cli/baseline_update_claude_link_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The Doctor output and the update's warnings — consumers: `./internal/cli`
- The user guide pages, the skills and the Behavior Surface Record — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDoctorReportsTheClaudeInstructionLink|TestBaselineUpdateWarnsAboutTheClaudeInstructionLink|TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorReportsTheClaudeInstructionLink/ok TestDoctorReportsTheClaudeInstructionLink/absent TestDoctorReportsTheClaudeInstructionLink/regular_file TestDoctorReportsTheClaudeInstructionLink/link_elsewhere TestDoctorReportsTheClaudeInstructionLink/file_with_passages TestDoctorReportsTheClaudeInstructionLink/link_elsewhere_with_passages TestDoctorReportsTheClaudeInstructionLink/agents_not_regular TestDoctorReportsTheClaudeInstructionLink/unlinked_checkout TestDoctorReportsTheClaudeInstructionLink/unread TestDoctorReportsTheClaudeInstructionLink/skipped TestDoctorReportsTheClaudeInstructionLink/after_baseline_rules_and_exit_code TestBaselineUpdateWarnsAboutTheClaudeInstructionLink/withheld_text TestBaselineUpdateWarnsAboutTheClaudeInstructionLink/withheld_json TestBaselineUpdateWarnsAboutTheClaudeInstructionLink/unlinked_checkout_text TestBaselineUpdateWarnsAboutTheClaudeInstructionLink/unlinked_checkout_json TestBaselineUpdateWarnsAboutTheClaudeInstructionLink/exit_code_unchanged TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the Doctor line and its tests do not exist, so their pass lines are missing; afterwards every form of the line, both warnings and the transcript pass, and the whole CLI suite passes with the adjusted Doctor tests.

- `grep -q 'baseline-claude-link' docs/user-guide/commands/doctor.md || { echo 'doctor guide lacks the line' >&2; exit 1; }; grep -q 'baseline.root.claude-link.withheld' docs/user-guide/commands/baseline.md || { echo 'baseline guide lacks the warning' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- "the next update then replaces" || { echo 'baseline reference lacks the sentence' >&2; exit 1; }; grep -q 'CLAUDE.md' .agents/skills/roundfix/references/modes.md || { echo 'modes reference lacks the question' >&2; exit 1; }; grep -q 'baseline-claude-link' .agents/skills/roundfix/references/setup.md || { echo 'setup reference lacks the line' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/setup-context-driven/SKILL.md | grep -qF -- "the next update then replaces" || { echo 'setup-context-driven lacks the sentence' >&2; exit 1; }; make skills-sync-check && go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides and skills do not mention the line, the warnings or the question; afterwards they do, the mirrors match, the owned skill versions are recorded, and the documentation contracts pass.

## References

- `_prd.md` → Goals (third, fourth and fifth bullets); Core Features 4 and 5; Success Metric 4; User Experience
- `_techspec.md` → Command surfaces; invariant 9; API Contracts 2, 3 and 5; Surface Transcript 1; Testing Approach (Task 3); Risks & Considerations (Doctor tests that count lines); Build Order 3
- ADR-0432, ADR-0433, ADR-0220, ADR-0411

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-testing, testing-boss, tech-writer, evidence-gate, roundfix, systematic-debugging.

### Implementation and acceptance evidence

- Doctor now reads the Task 01 reader with `baseline.ExecGitRunner{}` and prints exactly one advisory `baseline-claude-link:` line immediately after `baseline-rules:`. The line covers all API Contract 3 shapes, singular and plural passage counts, comma-separated locations, absent manifest and reader failure; it never returns `failed`. The Doctor usage describes it. `TestDoctorReportsTheClaudeInstructionLink` exercises all eleven required subtests through `runCLI`, checks whole lines and adjacent checks, preserves repository snapshots, and proves an independent failed check still controls the exit code.
- `TestBaselineUpdateWarnsAboutTheClaudeInstructionLink` covers both API Contract 2 warnings in text and JSON on adopted temporary Git repositories. The fixtures preserve `CLAUDE.md` byte-for-byte and the repository snapshot during preview, compare exits against the same fixture with the link, and exercise both exit 0 and exit 3 with an unrelated pending guide change. The unlinked fixture uses `core.symlinks=false` and a Git checkout of the tracked link, rather than a fake reader.
- `TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink` asserts Surface Transcript 1's ordered output prefix and empty stderr, exit 0, the declared remaining-output variation, no `CLAUDE.md` plan or `baseline.root.claude-link` warning, and the `AGENTS.md` symbolic-link target before and after the read-only preview of this repository.
- The two user guide pages and the canonical skill references describe the rule, both warning messages, Doctor's forms and codes, the passage-move approval requirement, and the single whole-plan-or-stop question before deletion or a `fileChanges` entry that updates `CLAUDE.md`. Preservation wording includes the narrow root-link exception. The mirrors were regenerated with `make skills-sync`, versions with the owned-skill recorder, and `docs/references/behavior-surfaces.json` with the coverage recorder. Final recorded skill versions: Roundfix `0.1.273`; setup-context-driven `0.0.10`.

### Focused checks

- Initial reproduction: `go test -count=1 -timeout 5m ./internal/cli -run '^TestDoctorReportsTheClaudeInstructionLink$/^ok$'` exited 1 because the required whole line was absent before the implementation.
- `go test -count=1 -timeout 5m -v ./internal/cli -run '^(TestDoctorReportsTheClaudeInstructionLink|TestBaselineUpdateWarnsAboutTheClaudeInstructionLink|TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink|TestDoctorReportsMisplacedRepositoryRules)$'` exited 0; every required new subtest and the existing misplaced-rule cases reported PASS.
- `go test -count=1 -timeout 5m ./internal/cli -run '^TestRunDoctor(ProfileReadinessProvesEffectiveCategoriesAndReportsCounts|AdapterReadinessReportsRequiredProfileRuntimes|RepositorySkillReadiness|MissingRepositoryRoot|RealRepositoryCheckDoesNotMutateState)$'` exited 0 for the existing outputs and line counts affected by insertion.
- `go test -count=1 -timeout 5m ./skills -run '^(TestTheUpdateModeAppliesOnlyThePlanItShowed|TestBaselineSkillContract|TestThinSetupSkill|TestRoundfixSkillIndexNamesEveryReference|TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak)$'` exited 0 against the final skill text.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^(TestBaselineDocumentationContract|TestUserGuideLinksResolve|TestGlossaryPhrasesAreDefinedTerms|TestCommandCoverageFollowsTheRoundfixReferenceIndex|TestSkillCoveragePathsRequireOwnedFilesAndMatchingSources)$'` exited 0.
- Required generation commands exited 0: `make skills-sync`; `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; and, after the last guide edit, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`. Mirrors were synchronized again after recording versions.
- `gofmt -l` reported no changed Go file; `git diff --check` reported no whitespace error. A changed-path audit found only the declared paths, this Task file and the companion assertion below; all five canonical/mirror pairs matched byte-for-byte. `CLAUDE.md` remained a symbolic link to `AGENTS.md`.

### Companion assertion and execution boundary

`internal/cli/doctor_baseline_rules_test.go` previously required `residue:` immediately after `baseline-rules:`. A focused run reproduced its failure with the required new line. Its sole changed assertion now requires `baseline-claude-link:` between them. This is an ordinary, non-Governed companion path within the same output-ordering slice, permitted by `docs/agents/spec-routing.md`; the Daemon can record it under Recorded paths. No assertion about passages or exit codes was weakened. The existing `doctor_test.go` fixture without root instruction files correctly gains the contract's `agents-not-regular` warning; its other expected additions are skipped lines.

The sandbox rejected its symlinked writable root before process creation; commands ran through the physical checkout with escalation. No other Task, Task Graph, prohibited path, Baseline source or golden was edited. Task status remains Daemon-owned. The Agent did not run declared Verification or repository Verification; no terminal Task verdict, commit, push or pull request is claimed.

### Verification Feedback repair — attempt 1

Inspected the Daemon's `batch-003-attempt-1.log` and the entry-size assertion in `skills/roundfix_layout_test.go`. The added mode-summary wording put the Roundfix skill entry over its existing 20,000-byte budget. Shortened only that summary, preserving the Update Check, release and Baseline preview, and approval before deleting files or replacing `CLAUDE.md`. The detailed approval question remains in `references/modes.md`; the budget and its test are unchanged.

Regenerated the mirror with `make skills-sync`, recorded Roundfix `0.1.273` with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, synchronized again, and ran the sanctioned coverage recorder after the final edit. All commands exited 0. Canonical and embedded entries are byte-identical at 19,990 bytes each.

Focused repair check: `go test -count=1 -timeout 5m ./skills -run '^(TestRoundfixSkillEntryFileStaysWithinItsBudget|TestTheUpdateModeAppliesOnlyThePlanItShowed|TestTheRoundfixSkillReadsItsModeFromTheRequest)$'` exited 0. `git diff --check` reported no whitespace errors. The full declared Verification sequence remains for the Daemon's next attempt; Task status and all other implementation files were left unchanged in this repair turn.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/doctor_baseline_rules_test.go`
