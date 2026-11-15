---
task: task_01
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
status: completed
type: backend
complexity: medium
---

# Task 01: The Verification probe prepares its checkout as the Daemon prepares a Task Worktree

## Overview

`roundfix spec check --run-verification` provisions its disposable checkout
with the repository's `worktree.copy` files and Worktree Bootstrap, through
the worktree package's provisioning function, once, right before the first
command it probes, so its verdict matches the Daemon's. A failed bootstrap
probes nothing and exits `1` with its own message, a report whose commands
ran after a bootstrap names it, and a repository without `worktree.copy` and
`worktree.bootstrap` is checked exactly as before. The `spec` guide, the
Roundfix skill's `spec` reference, the **Worktree Bootstrap** entry and the
notices in ADR-0148 and ADR-0034 state the rule. This answers the message
reported by indexus on 2026-11-14
([`references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md`](references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md)).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree,
   `spec check --run-verification` in a temporary repository whose Project
   Config sets a Worktree Bootstrap that creates an ignored file reports the
   command `test -f <that file>` as `honest`, and
   `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted`, written
   first, passes there.
2. MUST add `worktree.Provision` to `internal/worktree/worktree.go` and make
   `ProvisionItem` call it after validating its ref, as `_techspec.md` → API
   Contract 1 says; `ProvisionItem`, `Create` and `createTaskWithOptions`
   MUST NOT change behavior.
3. MUST add `verificationTreePreparation`, pass it from `runSpecCheck` with
   `loaded.Config.Worktree.Copy` and `worktreeBootstrapSpec(loaded.Config)`,
   and prepare the checkout lazily in `probeSpecVerifications`, exactly as
   API Contract 2 says: once per invocation, right before the first probed
   command, with `io.Discard` as the bootstrap output, and never for a
   completed Task, the QA gate, a Task without commands or a Spec without a
   Task Graph.
4. MUST add the `Bootstrap` report field, the text line and the JSON key of
   API Contract 3, and the failure path of API Contract 4 with its three
   standard-error parts word for word; every other probe error keeps
   `printSpecCheckFailure`.
5. MUST add `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes`,
   `TestSpecCheckStopsWhenTheWorktreeBootstrapFails`,
   `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand` and
   `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted`, with their
   subtests, in the new `internal/cli/spec_check_bootstrap_test.go`, with the
   fixtures and assertions of `_techspec.md` → Testing Approach. The first
   asserts Surface Transcript 1 and the second Surface Transcript 2, every
   line of both streams in order and the exit code, with the fixture's
   `Skipped:` block as the transcript's only controlled variation. No test
   reaches the network, reads the host's home or waits on elapsed time;
   every counter and marker file lives under `t.TempDir()`. MUST NOT change
   any existing test.
6. MUST add Text A of API Contract 5, word for word, at the anchor it names,
   to `docs/user-guide/commands/spec.md` and to the canonical
   `.agents/skills/roundfix/references/spec.md`, and MUST NOT add or change
   text inside any `### QA settlement` section.
7. MUST write Text B of API Contract 6 into the **Worktree Bootstrap** entry
   of `CONTEXT.md` through `domain-modeling`, keeping the entry's first
   sentence and its `_Avoid_` line, and MUST append Texts C and D of API
   Contract 7, word for word, as the last paragraphs of ADR-0148 and ADR-0034,
   setting each ADR's `updated_at` to the time of the edit and changing none
   of their other text.
8. MUST run `make skills-sync` to copy the reference to
   `skills/roundfix/references/spec.md`, then raise the Roundfix skill version
   to the next version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   on the tree the Task runs on, after Specs 0392 and 0393 raised it; never
   hand-edit a mirror, a version field or the record. Of the Governed Paths
   the Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/spec.md` and
   `skills/roundfix/SKILL.md`.
9. MUST, as its last step and after the last edit of the guide and of
   `internal/cli/spec_check.go`, re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Record the characterization.
- [ ] Add `worktree.Provision` and the lazy preparation of the probe.
- [ ] Add the report line, the JSON key and the bootstrap failure message.
- [ ] Add the four tests.
- [ ] Write Texts A to D, sync the mirror, record the skill version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A command that passes only on the tree prepared by `worktree.copy` and the Worktree Bootstrap is reported `vacuous`, and the check exits `1` (Surface Transcript 1).
- [ ] A failed bootstrap runs no command, prints nothing on standard output and the three-part message on standard error, and exits `1`, in text and JSON (Surface Transcript 2).
- [ ] The bootstrap runs zero times when nothing is probed and once for two Specs.
- [ ] Without `worktree.copy` and `worktree.bootstrap` the report and exit code are those of the starting tree.
- [ ] The guide, both reference copies with the same bytes, the glossary entry and both ADR notices carry their texts; the skills and documentation contracts pass.

## Context

- instruction: `docs/adr/0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md`
- instruction: `internal/cli/spec_check_test.go`
- instruction: `internal/cli/spec_check_completed_tasks_test.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/cli/spec_check.go`
- creates: `internal/cli/spec_check_bootstrap_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
- interface: `docs/adr/0034-run-worktrees-can-bootstrap-their-environment.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The text and JSON report and the exit code of `roundfix spec check --run-verification` — consumers: `./internal/cli`
- `worktree.Provision` and `worktree.ProvisionItem` — consumers: `./internal/worktree`
- The `spec` reference and the skill version record — consumers: `./skills`
- The `spec` guide, the glossary and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestSpecCheck(PreparesTheVerificationTreeAsTheDaemonDoes|StopsWhenTheWorktreeBootstrapFails|RunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand|WithoutWorktreePreparationPrintsWhatItPrinted)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes/surface_transcript_1 TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes/json_names_the_bootstrap TestSpecCheckStopsWhenTheWorktreeBootstrapFails/surface_transcript_2_text TestSpecCheckStopsWhenTheWorktreeBootstrapFails/surface_transcript_2_json TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand/nothing_to_probe TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand/two_Specs_and_two_Tasks TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted/text TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted/json; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'SpecCheck' ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/cli ./internal/worktree` — expected: exit 0. Before this Task the four tests do not exist, so their pass lines are missing; afterwards both Surface Transcripts hold, the bootstrap runs once and only when a command is probed, a repository without preparation prints what it printed, every `SpecCheck` test of the CLI passes, and both packages build for Windows.

- `grep -qF -- 'func Provision(ctx context.Context, userRoot, worktreePath string, opts ItemProvisionOptions) error' internal/worktree/worktree.go || { echo 'worktree.Provision missing' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/worktree` — expected: exit 0. Before this Task the function does not exist; afterwards it exists with the signature of `_techspec.md` → Interfaces and every test of the worktree package, the item provisioning tests included, passes.

- `for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Before the first command it probes, the check prepares the disposable checkout as the Daemon prepares a Task Worktree' || { echo "Text A missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A repository without' || { echo "Text A incomplete: $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md || { echo 'spec reference mirror differs' >&2; exit 1; }; grep -F -- '**Worktree Bootstrap**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Worktree Bootstrap\*\*(:| —( |$))' || { echo 'Worktree Bootstrap definition missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'runs it once in its disposable checkout, after the same copy and before the first command it probes' || { echo 'Text B missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md | grep -qF -- 'The authoring caller no longer supplies a bare working directory' || { echo 'ADR-0148 notice missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0034-run-worktrees-can-bootstrap-their-environment.md | grep -qF -- 'The Worktree Bootstrap also runs once in the disposable checkout' || { echo 'ADR-0034 notice missing' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the texts exists; afterwards the guide, both reference copies with the same bytes, the **Worktree Bootstrap** definition and both ADR notices carry them, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record and the glossary phrase contract pass.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Glossary
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contracts 1-7; Surface Transcript 1; Surface Transcript 2; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1
- ADR-0406, ADR-0148, ADR-0034
- [The adopted indexus message](references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md)


## Result

Skills loaded: implement-task, agentic-cli-design, coding-guidelines, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds, roundfix.

### Starting-tree characterization

Before changing production code, wrote `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted` first. `go test ./internal/cli -run '^TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted$' -count=1 -v` passed both `text` and `json` on the starting implementation. The text assertion covers all lines except the fixture's controlled `Skipped:` block; JSON asserts the original verification object suffix and exit `1`.

Then ran `go test ./internal/cli -run '^TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes$/surface_transcript_1$' -count=1 -v` before the implementation edit. In its temporary committed repository, `.gitignore` excludes `.deps-installed` and `local.env`, Project Config copies `local.env` and sets bootstrap to `test -f local.env && touch .deps-installed`. The CLI exited `0`, omitted the bootstrap line, and printed:

```text
Verification tree: HEAD
- task_01: honest — "test -f .deps-installed" (exited non-zero before work)
- task_01: honest — "test -f local.env" (exited non-zero before work)
- task_01: honest — "test -f output-the-task-writes" (exited non-zero before work)
```

The new transcript test failed on exit `0` versus required `1`, establishing the preparation mismatch before the production edit. No pre-existing tests were edited.


### Implementation and acceptance evidence

1. Prepared-tree verdicts: `worktree.Provision` copies the configured files and then runs bootstrap, returning the first error unchanged. `ProvisionItem` validates its ref before delegating; `Create` and `createTaskWithOptions` are unchanged. The CLI passes the loaded copy list and bootstrap spec, calls provisioning once immediately before the first probed command, and discards bootstrap output. The new `surface_transcript_1` subtest asserts the full ordered text transcript, empty stderr and exit `1` (only the fixture's `Skipped:` block varies). Its JSON companion asserts `HEAD`, the bootstrap command and the vacuous/vacuous/honest verdicts.
2. Failed bootstrap: both `surface_transcript_2_text` and `surface_transcript_2_json` assert empty stdout, the exact three-part stderr transcript, exit `1`, and absence of the command marker. The CLI identifies the wrapped `BootstrapError` with `errors.As`; other probe failures retain `printSpecCheckFailure`.
3. Laziness and frequency: `nothing_to_probe` proves zero bootstrap runs for a completed Task plus pending QA gate, a Task without commands, and a Spec without a Task Graph. The empty-command graph is already refused during loading with exit `2`; the test preserves that existing refusal and proves no bootstrap ran. `two_Specs_and_two_Tasks` proves four honest command reports, two bootstrap report lines and exactly one counter line. The second Spec has its own copied fixture authorization so consistency errors cannot determine the result. Counters and failure markers use `t.TempDir()`; all repositories are disposable and all bootstrap commands are local shell operations.
4. No preparation configured: the original text and JSON characterizations pass again after the implementation, preserving the original ordered report, absent bootstrap JSON key, empty stderr and exit `1`. Existing resumed-Task, QA omission, malformed-command and source-provenance tests pass unchanged in focused checks.
5. Documentation and generated records: exact Text A is inserted after the parser paragraph in the guide and canonical reference; `make skills-sync` generated the mirror. Through domain-modeling, Text B extends only the Worktree Bootstrap definition, preserving its first sentence and Avoid line. Texts C and D are each ADR's final paragraph; only `updated_at` changes elsewhere. The owned-skill recorder raised both Roundfix copies from `0.1.239` to `0.1.240` and appended the version record. Byte-level inspection confirmed exact Texts A-D, anchors, mirrors, unchanged QA settlement sections, and original ADR/glossary bytes outside the specified edits. The Behavior Surface Record recorder passed and changed only the spec guide fingerprint; it is rerun after this Result update as the final repository mutation.

### Focused checks

All commands below exited `0` on the implementation tree:

- `go test ./internal/cli -run '^TestSpecCheck(PreparesTheVerificationTreeAsTheDaemonDoes|StopsWhenTheWorktreeBootstrapFails|RunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand|WithoutWorktreePreparationPrintsWhatItPrinted)$' -count=1 -v` — all four tests and their required subtests passed, including the additional no-commands and no-Task-Graph cases.
- `go test ./internal/cli -run '^Test(SpecCheckRunVerification|SpecCheckLeavesCompletedTasksAlone|SpecCheckNeverProbesTheQAGate|SpecCheckReportsAMalformedCommand|SpecCheckNamesNoSourceWhenCommitted)$' -count=1` — focused existing CLI regressions passed.
- `go test ./internal/worktree -run '^Test(CreateRunsBootstrapAfterCopyInRunWorktreeRoot|CreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot|RunBootstrapReturnsBootstrapErrorOnNonZeroExit|RunBootstrapSkipsEmptyCommand|CleanupKeepsARegisteredItemWorktreeOnAnotherBranch|CleanupKeepsADetachedItemWorktree|CleanupRecoversAHalfRemovedItemWorktree)$' -count=1` — focused Run, Task and item worktree behavior passed.
- `make skills-sync` — synchronized the owned-skill mirror.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — recorded the next Roundfix skill version.
- `go test ./skills -run '^Test(CheckValidatesRoundfixSkillArtifacts|SpecReferenceLifecycleSkillContracts|RoundfixSkillIndexNamesEveryReference|AuthorialSkillSync|OwnedSkillEditLeavesDerivedArtifactsByteIdentical)$' -count=1` — the normal-build artifact, lifecycle, index and sync checks passed; the repocontract-tagged test is not part of this invocation.
- `go test -tags docscontract ./internal/docscontract -run '^Test(GlossaryPhrasesAreDefinedTerms|GlossaryLinksNameAcceptedADRs|GlossaryAvoidLinesArePlainLists|GlossaryStatesEachRuleOnce|EveryCommandIsNamedInTheRoundfixSkill)$' -count=1` — focused glossary and command-documentation contracts passed.
- `make baseline-digests` — repository-required regeneration returned `ok: true`, `changed: false`; no Baseline or frozen derived artifact changed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — regenerated the Behavior Surface Record after the last guide and CLI edits.
- `git diff --check` and byte-level Python inspections — clean whitespace; exact text, anchor, mirror and preservation contracts held. The Task's authored projection and existing tests are unchanged; its original Daemon-owned `in_progress` status is preserved.

The only new or changed paths belong to this Task's declared interfaces, new test file and Result. No follow-up work identified. Declared Verification commands, the full package suites, the Windows build and repository Verification remain for the Daemon. No Task status, other Task file or Task Graph manifest was edited; no commit, push, pull request or external message was created.
