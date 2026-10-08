---
task: task_03
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: completed
type: backend
complexity: high
---

# Task 03: Baseline update plans and applies the Pending History in the change it plans

## Overview

Gives `roundfix baseline update` its history section (ADR-0254). The update
plans the Pending History with the History Sanitize Command's own per-unit
planning and lists every unit, every Refused Unit and the History Full Tag it
would create. It binds that section into its Plan Digest. On `--yes` or
`--confirm-plan` it applies the Baseline Plan, creates the tag when it is
absent and converts every selected unit, all in the working tree. It never
commits or pushes. `--no-history` leaves the section out. The Task is
verifiable on its own through command tests on synthetic adopted repositories.

## Requirements

1. MUST extract `planHistoryUnit` from `runHistoryCommand`, as
   `_techspec.md` → Interfaces states: the folder file inventory, the delivery
   lookup and `PlanLegacyConversion`, or `PlanHistoryKind`, for one
   inventoried unit. `runHistoryCommand` MUST call it and keep every output,
   refusal and exit code of `roundfix history sanitize` unchanged.
   `historyInventory` MUST keep its signature and behavior.
2. MUST implement the history section of `_techspec.md` → API Contract 1 in
   the new `internal/cli/baseline_update_history.go`:
   - Load the `--repo` Git root's configuration with the command
     environment's home and resolve its Spec Root and archive root.
   - Plan each inventoried unit at `HEAD` with `planHistoryUnit`.
   - Refuse a unit in this order: for its conversion refusal; for
     `has uncommitted changes under <path>; commit or restore them and rerun`
     when `git status --porcelain=v1 --untracked-files=all -- <paths>` over
     the paths it removes or rewrites is not empty; and for
     `history-full does not hold <path>` when the tag exists without that path.
   - Select every other unit, with no batch limit.
   - Set the status, the tag object and the `blocked` cases as the contract
     states.

   Every refusal reason and `blocked` message MUST go through task_02's
   `historyRefusalLine`. A `blocked` section MUST NOT change the Baseline
   part.
3. MUST add the always-present `history` field of `_techspec.md` → Data
   Models to `baselineUpdateResult`, with exactly the stated JSON keys. `units`
   and `refused` are empty arrays, never null.
4. MUST compute the update's `planDigest` as `_techspec.md` → API Contract 2
   states. It equals the Baseline Plan Digest when no unit is selected and is
   the domain-separated combined digest otherwise. The confirmation check MUST
   compare against it before any write. A mismatch reports `action_required`,
   category `approval`, exit 3, prints the current digest and writes nothing.
   `baseline.ApplyPlan` receives the Baseline Plan Digest.
5. MUST implement `_techspec.md` → API Contract 3: a selected unit makes the
   preview `plan_ready` with the stated message even when the Baseline Plan is
   empty. Every other history status leaves the state, category, message and
   exit code exactly as today.
6. MUST implement the apply order of `_techspec.md` → API Contract 4: first
   `ApplyPlan`; then the tag, created only for action `create` with the stated
   message at the planned revision and checked to be annotated; then each
   selected unit through `spec.ApplyLegacyConversion` or `spec.ApplyHistoryKind`;
   then the existing skills stage. A tag or unit failure exits 1 with state
   `failed`, category `history`, a message naming the unit (or the tag) and
   the stated next action, and keeps the Baseline fields of the result. The
   update MUST NOT commit, push, move, replace or delete a tag.
7. MUST implement the text block of `_techspec.md` → API Contract 5 exactly,
   placed directly before `Plan Digest:` (or before `Next action:` when no
   digest is printed). It is printed only for `pending`, `applied` or
   `blocked`, or when a Refused Unit exists.
8. MUST add the `--no-history` boolean flag of `_techspec.md` → API Contract 6.
   It runs no history configuration load, inventory or Git command and reports
   `history.status: skipped`. The `baseline update` usage text in
   `internal/cli/cli.go` lists `[--no-history]` on its own usage line and gains
   one paragraph naming the Pending History, its conversion on approval, the
   created `history-full` tag that must be pushed, and `--no-history`. Every
   snippet `internal/docscontract/publicdocs_test.go` pins for that usage stays
   present.
9. MUST add these tests to the new
   `internal/cli/baseline_update_history_test.go`. Each builds an adopted
   repository with `newBaselineUpdateRepository`, adds history under
   `docs/history/` and commits it. Each runs
   `runBaselineUpdateCommandWithSkillsStage` with an environment whose home is
   a temporary directory and a successful fake skills stage, so no test reads
   the real `~/.roundfix`. Each invocation passes `--no-skills`, as the
   Surface Transcripts do, so the skills readiness of the host never decides a
   state. The history holds `0001-maps-unproven`, with a
   list-of-maps `unproven`; `0002-plain`; `0003-broken`, with an unparsable
   `_prd.md`; and one unreduced retired Finding, unless a test says otherwise.
   - `TestBaselineUpdatePlansPendingHistoryWithoutWrites`: the text preview
     matches Surface Transcript 1 line for line, ignoring placeholders. It
     exits 3, and the JSON preview lists three units and one Refused Unit with
     tag action `create` at `HEAD`. The tree, the index and
     `git for-each-ref` are byte-identical afterwards, and `planDigest`
     differs from the Baseline Plan Digest.
   - `TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag`:
     `--confirm-plan` with the preview's digest exits 0 (Surface Transcript
     2). It writes `0001-maps-unproven.md` and `0002-plain.md`, removes their
     folders, reduces the Finding and leaves `0003-broken` byte-identical.
     `git cat-file -t history-full` is `tag`, the tag's commit is the
     preview's `HEAD`, the commit count is unchanged and no other ref changes.
     A second update without approval exits 0 with state `current`, history
     `current` and the Refused Unit still listed (Surface Transcript 3).
   - `TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest`: with no
     history, `planDigest` equals the Baseline Plan Digest from
     `baseline.BuildPlan`, history is `current`, and no text line starts with
     `History`.
   - `TestBaselineUpdateNoHistorySkipsTheSection`: `--no-history`, as preview
     and with `--yes`, reports `skipped`, leaves every history file and the
     refs unchanged and prints no `History` line (Surface Transcript 4).
     `baseline update --help` names `--no-history`.
   - `TestBaselineUpdateRefusesUnitsWithUncommittedChanges`: an uncommitted
     edit inside `0002-plain`, and an untracked file inside `0001-maps-unproven`,
     each make that unit a Refused Unit with the stated reason. The other units
     stay selected, and apply leaves the refused bytes as they were.
   - `TestBaselineUpdateHistoryHonorsAnExistingTag`: with an annotated
     `history-full` at a commit that lacks a later folder, the tag action is
     `present`, that folder is refused with
     `history-full does not hold <path>`, and apply leaves the tag object
     unchanged. A lightweight tag reports `blocked` with no unit selected and
     the Baseline state unchanged.
   - `TestBaselineUpdateRejectsAStaleHistoryDigest`: `--confirm-plan` with the
     Baseline Plan Digest alone, while a unit is selected, exits 3 with
     category `approval` and writes nothing.
   - `TestBaselineUpdateReportsATagFailureBeforeAnyConversion`: with the
     fixture's local Git configuration set so `git tag -a` fails
     (`tag.gpgSign=true` and `gpg.program=false`), `--yes` exits 1 with
     category `history`. No unit is converted, no tag exists, and the Baseline
     fields of the result are kept.
10. MUST NOT change any existing assertion in `internal/cli` tests. Every
    existing `baseline update` and `history` test passes unchanged
    (`_techspec.md` → Invariant 1).

## Subtasks

- [ ] Extract `planHistoryUnit` and keep `roundfix history sanitize` unchanged.
- [ ] Plan the history section with its refusals, tag state and digest.
- [ ] Apply in order with the tag and the conversions, and report text and JSON.
- [ ] Add `--no-history` and the usage text.
- [ ] Write the command tests.

## Acceptance Criteria

- [ ] A preview with Pending History lists units, Refused Units and the tag, exits 3 and writes nothing.
- [ ] The confirmed update converts every selected unit, creates the annotated tag at the planned `HEAD`, and commits and pushes nothing; a second update is current.
- [ ] A repository without Pending History gets today's digest and output, and `--no-history` touches no history.

## Context

- interface: `internal/cli/history.go`
- creates: `internal/cli/baseline_update_history.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/baseline_update_history_test.go`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/cli/history_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBaselineUpdatePlansPendingHistoryWithoutWrites|TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag|TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest|TestBaselineUpdateNoHistorySkipsTheSection|TestBaselineUpdateRefusesUnitsWithUncommittedChanges|TestBaselineUpdateHistoryHonorsAnExistingTag|TestBaselineUpdateRejectsAStaleHistoryDigest|TestBaselineUpdateReportsATagFailureBeforeAnyConversion)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestBaselineUpdatePlansPendingHistoryWithoutWrites TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest TestBaselineUpdateNoHistorySkipsTheSection TestBaselineUpdateRefusesUnitsWithUncommittedChanges TestBaselineUpdateHistoryHonorsAnExistingTag TestBaselineUpdateRejectsAStaleHistoryDigest TestBaselineUpdateReportsATagFailureBeforeAnyConversion; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix baseline update --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -q -- '--no-history' || { printf 'baseline update help does not name --no-history\n%s\n' "$help" >&2; exit 1; }; go test -count=1 -run '^(TestBaselineUpdate|TestHistorySanitize|TestHistorySnapshot)' ./internal/cli` — expected: exit 0. Before this Task the eight new tests do not exist and the usage has no `--no-history`, so the first check fails. After it, the new tests pass, the built help names the flag, and every existing `baseline update` and History Sanitize Command test passes unchanged.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; Goals; User Story 1; User Story 2; User Story 3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-9; Invariant 12; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Build Order 3
- ADR-0254; ADR-0248; ADR-0100

## Result

Implemented the Pending History section of `baseline update`. The update reuses
`planHistoryUnit` with the History Sanitize Command, loads the target Git root's
configuration through the command environment, and selects all convertible,
committed units covered by an existing annotated ancestor tag or the planned
new tag. Conversion refusals, dirty paths and tag coverage refusals stay in
place; a blocked history section leaves the Baseline outcome intact. History
status, units and refusals are always present in JSON with non-null arrays.

Selected units join the Baseline Plan Digest through the domain-separated
history digest. Confirmation is checked before any write; `ApplyPlan` receives
its own Baseline Plan Digest. Apply then creates the annotated tag at the
planned revision, converts the selected units and reaches the existing skills
stage. History failures retain the Baseline result and name recovery commands.
The text history block and `--no-history` usage follow the TechSpec contracts.
Approval refusals keep the existing Baseline wording.
Status uses Git's `--no-optional-locks` to keep preview index bytes unchanged.

### Focused evidence by acceptance criterion

- Preview: `TestBaselineUpdatePlansPendingHistoryWithoutWrites` checks the
  Surface Transcript 1 history block line for line, the three selected units,
  one Refused Unit and planned tag, exit 3, the combined digest, and identical
  working-tree bytes, actual Git index bytes and refs. It deliberately changes
  a tracked file's timestamp without changing its content to exercise index
  refresh safety.
- Confirmed update: `TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag`
  checks both records, the reduced Finding, unchanged refused bytes, the tag
  object type and planned commit, unchanged commit count and other refs, and a
  second `current` update retaining its refusal. The `apply text` boundary
  subtest checks Surface Transcript 2's history block. Dirty-unit, existing-tag,
  stale-digest and tag-failure tests check the refusal and failure contracts.
- No pending history and opt-out:
  `TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest` compares
  `baseline.BuildPlan`'s digest and checks silent history text and non-null JSON
  arrays and unchanged approval-refusal wording.
  `TestBaselineUpdateNoHistorySkipsTheSection` checks preview and apply,
  unchanged history and refs, skipped status and help. The configuration-failure
  boundary subtest also proves that opt-out succeeds with an unavailable home.

Focused commands and outcomes:

- Before implementation,
  `rtk proxy go test ./internal/cli -run '^TestBaselineUpdatePlansPendingHistoryWithoutWrites$' -count=1`
  failed to compile because `baselineUpdateResult.History` was absent.
- After the final implementation edits,
  `GOCACHE="$PWD/.gocache" rtk proxy go test ./internal/cli -run 'TestBaselineUpdate(PlansPending|AppliesPending|WithoutPending|NoHistory|RefusesUnits|HistoryHonors|RejectsAStaleHistory|ReportsATagFailure|HistoryPlanningBoundaries|RejectsDigestOtherThanCurrentPlanWithoutWrites)' -count=1`
  exited 0. All eight required tests and the additional tag-name,
  non-ancestor-tag, configuration-failure/opt-out, dirty-kind and apply-text
  boundary cases passed. The existing stale-digest regression test also passed
  unchanged.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestBaselineDocumentationContract$' -count=1`
  exited 0. The initial invocation without the scoped cache was denied access
  to the host Go cache; no repository configuration was changed. After the
  final edits, the same documentation test also exited 0 using
  `GOCACHE="$PWD/.gocache"`.
- The first `rtk make verify-incremental` invocation failed only in
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion` and
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner`: the sandbox
  denied reading the host process table. Other reported packages passed. A
  rerun with the required process-table permission exited 0, including the
  complete CLI suite, vet, skills checks and build. The final rerun after the
  approval-message compatibility refinement also exited 0, with the complete
  CLI suite passing in 225.096 seconds.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The pre-existing Task status change was preserved. No existing test assertion,
Task status, other Task file or Task Graph manifest was edited. The authored
Verification commands were not run; their execution and settlement remain with
the Daemon. No commit, push or Pull Request was made. No follow-up scope was
added.
