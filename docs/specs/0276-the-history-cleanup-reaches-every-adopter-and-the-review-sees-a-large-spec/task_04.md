---
task: task_04
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: completed
type: backend
complexity: high
---

# Task 04: The archive plan has an offline form, and the Delivery Queue runs a fresh copy of each item binary

## Overview

Today `roundfix archive <slug> --plan` always calls the archive-value judge,
which a QA Task that forbids provider requests cannot allow, and the Delivery
Queue probes and runs the item binary at the path its build writes, so a later
build that writes in place can rewrite a file a running item binary maps
(`_techspec.md` → Current behavior). This Task adds `--plan --offline`, points
the `qa-gate` skill at it, and makes `stepExecutable` copy each build to a new
renamed file under the item's artifact directory and run that copy. It is
verifiable on its own through command tests with a counting judge transport
and through delivery tests that rebuild in place while a copy runs.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-16,
   [a QA gate's archive plan contacts the judge an offline QA forbids](references/2026-10-16-a-qa-gate-archive-plan-contacts-the-judge-an-offline-qa-forbids.md),
   by implementing API Contract 8 of `_techspec.md` in
   `internal/cli/archive.go`: `--offline` with `--plan` prints
   `judge: skipped (offline)` and each candidate as
   `candidate <file> <bytes> bytes: no advice (offline)`, calls neither
   `judge.Load` nor `judge.AdviseArchive`, and reads no Judge Key Variable;
   `--offline` without `--plan` is refused with `--offline requires --plan`;
   the usage gains `[--offline]` and its option line (ADR-0282, ADR-0273).
   `--plan` alone MUST keep today's output.
2. MUST change the "Outcome for the Archive Record" section of
   `.agents/skills/qa-gate/SKILL.md` to run
   `roundfix archive <slug> --plan --offline`, stating that a QA Task never
   contacts the judge, and MUST NOT change its `### QA settlement` section. No
   Go source renders that instruction into a QA prompt today; the Task MUST
   search the guides, skill copies and `internal/` for other instructions
   that tell a QA Agent to run the archive plan and change each.
3. MUST answer the adopted Backlog Entry of 2026-10-16,
   [an item binary built in place can be killed by code signing](references/2026-10-16-an-item-binary-built-in-place-can-be-killed-by-code-signing.md),
   by implementing API Contract 9 in `stepExecutable` in
   `internal/cli/deliver_workflow.go`: after a successful build, write the
   built regular file to
   `<artifact dir>/delivery/<slug>/item-binary/<first 16 hex of its SHA-256>/<base name>`
   through a temporary file created in that directory and a rename, with mode
   `0755`, then probe and run that copy and name it in the selection line. A
   failure to read or copy returns `copy item binary: <error>`, which parks the
   item `delivery-error`. The copy MUST never be written in place and the
   build output MUST never be executed (ADR-0281, ADR-0225).
4. MUST add the tests named in Verification, each calling `t.Parallel()`
   first: in the new `internal/cli/archive_plan_offline_test.go`,
   `TestArchivePlanOfflineTranscript` (asserts every line of Surface
   Transcript 1 of `_techspec.md`, in order, its empty standard error and exit
   `0`, for a Spec fixture whose only candidate is one `references/` file),
   `TestArchivePlanOfflineSkipsTheJudge` (with a Judge Key Variable set in the
   command environment and a judge transport that counts requests, the count
   stays zero) and `TestArchiveOfflineRequiresPlan`; in the new
   `internal/cli/deliver_item_binary_copy_test.go`,
   `TestDeliveryStepRunsACopyOfTheItemBinary` (the executed path is the copy
   under the artifact directory, its bytes equal the build output and the
   selection line names it) and
   `TestDeliveryRebuildInPlaceNeverRewritesARunningCopy` (a build command
   that copies over the build output in place runs while a held copy from the
   previous step still runs; the held copy's bytes and file identity are
   unchanged, the second step executes a path that differs from the build
   output, and the copy directory holds no temporary file).
5. MUST update `TestArchiveHelpNamesThePullRequestOption` in
   `internal/cli/archive_pull_request_test.go` to the new usage, and the
   helpers and assertions of `internal/cli/deliver_item_binary_test.go` that
   expect the build output path to expect the copy path; these are the
   declared breaks of their old expected text. MUST NOT rename or remove any
   test.
6. MUST document `--offline` in `docs/user-guide/commands/archive.md` and
   `.agents/skills/roundfix/references/archive.md`, and the copy, its path and
   the selection line in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`. No text inside any
   `### QA settlement` section may change.
7. MUST run `make skills-sync`, raise and record the Roundfix and qa-gate
   skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits user guide pages and command sources. No record or
   mirror is hand-edited.

## Subtasks

- [ ] Add `--offline` to the archive plan and its usage.
- [ ] Point the `qa-gate` skill at the offline plan.
- [ ] Copy, rename, probe and run the item binary copy.
- [ ] Add the five tests and move the help and item binary assertions.
- [ ] Document both changes in the guides and references.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] `roundfix archive <slug> --plan --offline` prints Surface Transcript 1 and sends no judge request with a key in its environment.
- [ ] `--offline` without `--plan` is refused, and `--plan` alone is unchanged.
- [ ] The `qa-gate` skill tells the QA Agent to use the offline plan.
- [ ] The Delivery Queue probes and runs a renamed copy, and an in-place rebuild never rewrites a running copy.

## Consumers

- The archive plan output and the Delivery Queue's item binary selection — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- instruction: `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`
- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/archive_pull_request_test.go`
- creates: `internal/cli/archive_plan_offline_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- creates: `internal/cli/deliver_item_binary_copy_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchivePlanOfflineTranscript|TestArchivePlanOfflineSkipsTheJudge|TestArchiveOfflineRequiresPlan|TestDeliveryStepRunsACopyOfTheItemBinary|TestDeliveryRebuildInPlaceNeverRewritesARunningCopy|TestArchiveHelpNamesThePullRequestOption)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchivePlanOfflineTranscript TestArchivePlanOfflineSkipsTheJudge TestArchiveOfflineRequiresPlan TestDeliveryStepRunsACopyOfTheItemBinary TestDeliveryRebuildInPlaceNeverRewritesARunningCopy TestArchiveHelpNamesThePullRequestOption; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in .agents/skills/qa-gate/SKILL.md docs/user-guide/commands/archive.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '--plan --offline' || { printf 'missing --plan --offline in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'item-binary/' || { printf 'missing the item binary copy path in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Archive|DeliveryStep|DeliveryItem|ItemBinary' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the five new tests do not exist, so the command fails at its first pass check; afterwards the archive, delivery, documentation, skill version and Behavior Surface Record checks pass.

## References

- `_prd.md` → Goals; User Story 6; User Story 7; Core Feature 7; Core Feature 8; Core Feature 9; Success Metric 5; Success Metric 6
- `_techspec.md` → API Contract 8; API Contract 9; API Contract 10; Surface Transcript 1; Invariants 7 and 8; Current behavior; Testing Approach; Build Order 4
- ADR-0282; ADR-0281; ADR-0225; ADR-0273

## Result

Implemented the Task 04 slice for Daemon Verification. Task status and the
Subtasks/Acceptance Criteria checkboxes remain Daemon-owned. No declared
Verification command, selected repository Verification, commit, push or Pull
Request operation was run.

### Acceptance evidence

- Offline archive plan: `runArchivePlan` bypasses `judge.Load`, Judge Key
  Variable inspection and `judge.AdviseArchive` when `--offline` is set. It
  prints the skip line after evidence and every candidate in sorted plan file
  order. `TestArchivePlanOfflineTranscript` passed with the complete ordered
  Surface Transcript 1, empty stderr, exit 0 and unchanged Spec files.
  `TestArchivePlanOfflineSkipsTheJudge` passed with every Judge Key Variable
  populated in the isolated command environment and zero transport calls.
- Flag validation and the existing plan: `TestArchiveOfflineRequiresPlan`
  passed with the exact `--offline requires --plan` preflight transcript and
  exit 2. `TestArchiveHelpNamesThePullRequestOption` passed with the new usage
  and option. `TestArchivePlanWithoutAKeyPrintsNoAdvice` passed for the plain
  `--plan` form. Its judge loading, key handling and advice formatting remain
  inside the original online path.
- QA instruction: the canonical qa-gate Outcome section now runs
  `roundfix archive <slug> --plan --offline` and states that a QA Task never
  contacts the judge. Searched guides, canonical/distributed skill copies and
  `internal/` for archive-plan instructions. The qa-gate pair was the only
  QA-specific instruction; the generic archive operator instructions still
  describe optional advice and now document the offline form. No Go source
  renders that QA instruction. Direct comparisons against HEAD and
  `TestSettlementGuidanceIsOneTable` proved the protected settlement sections
  unchanged.
- Delivery copy: successful builds are read only as regular files, hashed,
  written through a sibling temporary file with mode 0755, closed and renamed
  into `delivery/<slug>/item-binary/<first 16 SHA-256 hex>/<base name>` under
  the artifact directory. Only the published copy is probed and selected;
  its absolute path appears in the selection line. Copy failures return
  `copy item binary: <error>` before any probe. Existing owner fallback is
  preserved.
- Running-copy isolation: `TestDeliveryStepRunsACopyOfTheItemBinary` passed
  for the exact hash path, matching bytes, mode, selection transcript and
  recorded executable paths of both probe and step.
  `TestDeliveryRebuildInPlaceNeverRewritesARunningCopy` passed while a pipe
  held the first executable copy running: an in-place build preserved the
  build-output inode, left the held copy's bytes and identity unchanged,
  probed a different artifact copy and left no temporary files. A repeated
  identical build also replaced its published copy with a new inode.
  `TestDeliveryItemBinaryCopyFailureNeverRunsTheBuildOutput` passed for a
  non-regular output and blocked publication; neither ran a probe, and failed
  publication cleaned its temporary file. Existing argument, fallback,
  build-error, missing-output, rebuild, probe-report and signal tests passed.

### Focused checks and generators

- Before implementation,
  `go test ./internal/cli -run '^(TestArchivePlanOfflineTranscript|TestArchiveOfflineRequiresPlan|TestDeliveryStepRunsACopyOfTheItemBinary)$' -count=1`
  failed as expected: undefined `--offline` and selection of the build output.
- After implementation,
  `go test ./internal/cli -run '^(TestArchivePlanOffline.*|TestArchiveOfflineRequiresPlan|TestDelivery(Step.*|StepsStart.*|RebuildInPlace.*|ItemBinary.*|ProbeReport.*)|TestArchiveHelpNamesThePullRequestOption|TestArchivePlanWithoutAKeyPrintsNoAdvice)$' -count=1 -v`
  passed, including all five required new tests and the intentional existing
  contract changes.
- `go test ./internal/cli -run '^TestDeliveryItemBinaryCopyFailureNeverRunsTheBuildOutput$' -count=1 -v`
  passed both copy-error cases.
- `make skills-sync` succeeded. The owning record command
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  succeeded and raised Roundfix to 0.1.89 and qa-gate to 0.0.16 in both copies.
- The owning Behavior Surface Record command
  `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  succeeded, updating only the archive command and archive/deliver guide
  fingerprints. No mirror or record was hand-edited.
- The repository-required `make baseline-digests` succeeded with
  `changed: false`; no derived Baseline artifact changed.
- `go test ./skills -run '^(TestSettlementGuidanceIsOneTable|TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical)$' -count=1`
  passed. Direct file comparisons confirmed all four changed canonical/mirror
  pairs match, and the changed skills' QA settlement sections equal HEAD.
- `GOOS=linux go build ./internal/cli` passed for the other platform named
  by the updated delivery signal test's build constraint.
- `git diff --check` passed.

### Additional affected test consumers

Three existing test files also consume the intentionally changed executable
selection or its shared selection-line helper:
`internal/cli/deliver_item_binary_rebuild_test.go`,
`internal/cli/deliver_item_probe_report_test.go`, and
`internal/cli/deliver_item_binary_signal_unix_test.go`.
Their assertions/call sites now expect the copy; their tests were neither
renamed nor removed. Build-output inode assertions still inspect the build
output. No other Task or Task Graph file was edited.

Daemon-owned Verification and settlement remain pending.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/deliver_item_binary_rebuild_test.go`
- `internal/cli/deliver_item_binary_signal_unix_test.go`
- `internal/cli/deliver_item_probe_report_test.go`
