---
task: task_01
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
status: completed
type: backend
complexity: high
---

# Task 01: A confirmed Dirty Discard removes an ended Run's leftovers after recording them

## Overview

On 2026-11-09 Fiscus found eleven Runs that ended Failed or Unresolved for
Specs other Runs had delivered and archived, each kept `dirty` by uncommitted
leftovers outside the archived Spec's declared and recorded paths, with no
supported way out
([`references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md`](references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md),
reported by fiscus on 2026-11-09, gap 1). This Task delivers the Dirty
Discard of ADR-0379: the dry run lists each qualifying Run's uncommitted
paths with one digest, and `roundfix reconcile --discard-dirty <digest>`
writes a Branch Disposition record naming them, re-proves the Run and removes
its Run Worktree and Run Branch.

## Requirements

1. MUST, in `internal/worktree/worktree.go`, add
   `func DirtyDiscardOutcomes() []string` returning, in this order,
   `store.StateFailed`, `store.StateUnresolved`, `store.StateStopped`,
   `store.StateBudgetExceeded` and `store.StateTimedOut`, and the exported
   field `DiscardPaths []string` on `RunWorktreeReconciliation`.
2. MUST, in the dirty branch of `inspectTerminalRunMerged`, set
   `DiscardPaths` to the whole parsed uncommitted path set, in the order
   `terminalRunDirtyPaths` returns it, only when every one of these holds,
   and leave it nil otherwise:
   - the Run's outcome is one of `DirtyDiscardOutcomes()`;
   - the porcelain status parsed;
   - the Spec's Archive Record is on the default branch head
     (`resolveDefaultBranchHead` and `specArchivedAtMergedHead` at that
     head), whatever head the commit proof used;
   - the commit proof (`inspectRunAtMergedHead` or `inspectRunByDelivery`,
     as today) is `safe` or `superseded`;
   - `dirtyPathsInArchivedSpec` named an outside path (not the
     "could not be inspected" case).

   The Run stays `dirty` and its `Reason` stays byte-identical to what it is
   at the authoring base for every Run, candidate or not. A candidate MUST
   carry inspection evidence (via `newTerminalRunReconciliationEvidence`)
   holding the commit proof's `proofHead` and `proofRef`, the merged-head
   records it used, and the path set, so a later re-proof can compare them.
3. MUST add
   `func DiscardDirtyTerminalRun(ctx context.Context, result RunWorktreeReconciliation) error`,
   which refuses, changing nothing, a result whose evidence is missing or
   whose snapshot changed, or whose `DiscardPaths` is empty; re-inspects the
   Run with the same merged-head records; refuses, changing nothing, unless
   the fresh result is `dirty` with the same Run head, target head, proof
   head and `DiscardPaths`; and then runs `git worktree remove --force` on
   the recorded path and `git branch -D` on the Run Branch when it is
   present, reporting a failure through `terminalRunApplyFailure`. It MUST
   NOT write to the Run Database.
4. MUST, in `internal/daemon/reconcile.go`, add to `BranchDisposition` the
   fields `DiscardedPaths []string` with JSON name `discardedPaths` and
   `DiscardDigest string` with JSON name `discardDigest`, both `omitempty`,
   and
   `func WriteDirtyDiscardRecord(path string, disposition BranchDisposition, recordedAt time.Time) error`,
   which refuses a disposition that is `Superseded`, has no
   `DiscardedPaths` or no `DiscardDigest`, a path that is not absolute or
   that lies inside `disposition.Worktree`, and otherwise writes the same
   atomic `roundfix-branch-disposition/v1` record as
   `WriteBranchDispositionRecord`, sharing one unexported writer with it.
   `WriteBranchDispositionRecord` keeps its behavior and bytes.
5. MUST, in `internal/cli/reconcile.go`:
   - parse `--discard-dirty <digest>` and `--discard-dirty=<digest>`; a
     missing value fails with `flag needs an argument: --discard-dirty`, a
     value that is not 64 lowercase hexadecimal characters fails with
     `--discard-dirty needs the 64-character digest that roundfix reconcile prints`,
     and combining it with another mutation switch fails with
     `--apply, --discard-superseded, --discard-dirty, and --carry-forward are mutually exclusive`
     (the existing three-switch message changes to this one); all three
     fail through `printReconcileValidationFailure` with exit 2 before the
     configuration is loaded;
   - copy `DiscardPaths` into a new `reconcileResult.DiscardPaths` (JSON
     `discardPaths`, `omitempty`) and, after the process candidates are
     inspected, call a new `withholdDirtyDiscardFromLiveRuns(report *reconcileReport)`
     that clears it on every Run whose ID appears with kind `process` in
     `processCandidates` or `preservedCandidates`, setting that Run's
     `refusalReason` to
     `a live process tree is attributed to this Run, so its dirty leftovers are not offered for discard`;
   - compute `dirtyDiscardDigest`: the lowercase hex SHA-256 of, for each
     Run with `DiscardPaths` in Run ID byte order, its Run ID, Run head,
     worktree and each discard path in byte order, each followed by one NUL
     byte, and one more NUL byte closing the Run; store it with the command
     `roundfix reconcile --discard-dirty <digest>` (or
     `roundfix reconcile <run-id> --discard-dirty <digest>` with a Run ID) in
     `reconcileReport` fields `DiscardDirtyDigest` and
     `DiscardDirtyCommand` (JSON `discardDirtyDigest` and
     `discardDirtyCommand`, `omitempty`), only when at least one Run has
     `DiscardPaths`;
   - in a dry run, set each such Run's action to
     `would discard <N> uncommitted path(s) with --discard-dirty`, where
     `<N>` is the path count;
   - with `--discard-dirty`, use the mode `discard-dirty`; when no Run has
     `DiscardPaths`, fail with
     `no ended Run has dirty leftovers a Dirty Discard can remove; run roundfix reconcile to see why each dirty Run is kept`;
     when the computed digest differs, fail with
     `the Dirty Discard set changed: confirmed digest <given>, current digest <current>; run roundfix reconcile and confirm the digest it prints`;
     both through `printReconcileValidationFailure`, exit 2, with nothing on
     stdout and nothing written or removed;
   - otherwise, for each such Run in report order, resolve the Artifact
     Root as `--discard-superseded` does, write
     `<artifact root>/runs/<run-id>/branch-disposition.json` with
     `WriteDirtyDiscardRecord` (branch, target, Run head, target head,
     worktree, `superseded: false`, empty `commits` and `changedFiles`,
     `reason` `dirty leftovers discarded by a confirmed Dirty Discard`, the
     Run's dirty reason as `refusalReason`, the paths and the digest), then
     call `runworktree.DiscardDirtyTerminalRun`; on success set the action
     to `discarded`, append `; branch record: <record path>` to the
     evidence, clear the refusal reason and count it in `summary.applied`;
     on a failure keep the Run with action `preserve`, the error as its
     refusal reason, and count an operational failure (exit 1). A record that
     cannot be written MUST stop that Run before anything is removed;
   - in the text report, print after a Run's `  refusal-reason:` line one
     `  discard-path: <path>` line per discard path, and in a dry run print
     `Discard dirty leftovers with: <command>` after the `Apply with:` line
     when a digest exists. A report without a candidate keeps every byte it
     prints at the authoring base, text and JSON.
6. MUST, in `internal/cli/cli.go`, change the `reconcile` synopsis in the
   global usage and in `commandUsage("reconcile")` to
   `roundfix reconcile [run-id] [--apply | --discard-superseded | --discard-dirty <digest> | --carry-forward] [--format <text|json>]`
   and add the option line
   `  --discard-dirty         Discard the dirty leftovers the dry run listed under the given digest`
   after the `--discard-superseded` line; update the `reconcile` row of
   `TestRunCommandHelp` in `internal/cli/cli_test.go` to the new synopsis,
   and every other test that pins the synopsis or the three-switch message on
   the default branch when the Task runs.
7. MUST create `internal/worktree/dirty_discard_test.go`, over the
   archived-leftover fixtures (`newArchivedLeftoverFixture`,
   `addDeclaredLeftovers`, `inspectMergedHeadTestRun`), with:
   - `TestAnEndedRunOfAnArchivedSpecWithAnOutsideLeftoverOffersADirtyDiscard`:
     the declared leftovers plus `unrelated.txt` give a `dirty` result whose
     `Reason` equals the outside-path reason and whose `DiscardPaths` holds
     all three paths;
   - `TestADirtyDiscardIsNeverOfferedOutsideItsBounds`, one subtest each for
     an unrepresented Task commit, a Spec not archived on the default branch,
     the outcome `Clean`, and leftovers all inside the scope (`superseded`):
     `DiscardPaths` is empty;
   - `TestDiscardDirtyTerminalRunRemovesTheRunWorktreeAndRunBranch`;
   - `TestDiscardDirtyTerminalRunRefusesAChangedPathSet`: a file added after
     inspection, and a result that is not a candidate, are refused and the
     Run Worktree and Run Branch stay.
8. MUST create `internal/daemon/dirty_discard_record_test.go` with
   `TestWriteDirtyDiscardRecordNamesTheDiscardedPaths`: the record in
   `t.TempDir()` decodes with schema `roundfix-branch-disposition/v1`,
   `superseded` false, `discardedPaths` and `discardDigest`; a superseded
   disposition, one without paths and a path inside the worktree are refused
   and write nothing.
9. MUST create `internal/cli/reconcile_discard_dirty_test.go`, over
   `newReconcileMergedFixture` with an archived, completed Spec and an
   untracked `unrelated.txt` in the Run Worktree, with temporary homes only:
   - `TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest`: JSON
     `discardPaths`, a 64-character `discardDirtyDigest`, the
     `discardDirtyCommand`, the dry-run action, and the text lines
     `  discard-path: unrelated.txt` and
     `Discard dirty leftovers with: roundfix reconcile <run-id> --discard-dirty <digest>`;
     a repeated dry run prints the same digest;
   - `TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns`: the
     confirmed digest exits 0, the record names `unrelated.txt`, the Run
     Worktree and Run Branch are gone, and the next dry run reads `released`;
   - `TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything`:
     a new untracked file after the dry run makes the old digest fail with
     exit 2, the exact stderr of Requirement 5, empty stdout, no record, and
     the Run Worktree and Run Branch kept;
   - `TestReconcileDiscardDirtyRefusesWhenNoRunIsOffered`: a dirty Run of an
     unarchived Spec gives exit 2 and the exact no-candidate stderr;
   - `TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree`: given a
     report whose Run carries `DiscardPaths` and a `process` candidate for
     the same Run ID (one case in `processCandidates`, one in
     `preservedCandidates`), `withholdDirtyDiscardFromLiveRuns` clears its
     `DiscardPaths`, sets the refusal reason of Requirement 5, and
     `dirtyDiscardDigest` then returns no digest; a `runOwnerRecord`
     candidate for the Run withholds nothing;
   - `TestReconcileDiscardDirtyNeedsItsDigestTranscript`,
     `TestReconcileDiscardDirtyRefusesAMalformedDigestTranscript` and
     `TestReconcileMutationSwitchesExcludeEachOtherTranscript`, each asserting
     its whole TechSpec Surface Transcript: empty stdout, both stderr lines in
     order, and exit 2.

   Every new top-level test calls `t.Parallel()` as its first statement, and
   every test file compiles for Windows.
10. MUST, in `docs/user-guide/commands/reconcile.md`: change the synopsis to
    the one of Requirement 6; change "Three switches mutate" to "Four
    switches mutate" and add the example line
    `roundfix reconcile --discard-dirty <digest>          # discard the dirty leftovers the dry run listed`;
    replace "`--apply` remains the only switch that releases Run Worktrees."
    with "`--apply` releases Run Worktrees on proof; `--discard-dirty`
    removes only the Runs the dry run listed under the confirmed digest.";
    and add, after the paragraph that begins "A kept dirty Run whose Spec is
    archived", a paragraph stating the bounds of Requirement 2, the dry-run
    lines and JSON fields of Requirement 5, the refusals with exit 2 and
    nothing on stdout, and that the switch writes
    `<artifact root>/runs/<run-id>/branch-disposition.json` with
    `discardedPaths` and `discardDigest` before it removes the Run Worktree
    with `--force` and deletes the Run Branch.
11. MUST, through the `domain-modeling` skill, make the three `CONTEXT.md`
    changes the block after this list gives, byte for byte: append the
    first sentence to the **Branch Disposition** definition line, after "A
    superseded Branch Disposition means every branch commit is already
    reachable from the target, or a later integrated Run covered the same
    Tasks." and a space; add the **Dirty Discard** entry, its three lines,
    after **Branch Disposition**'s `_Avoid_` line, separated by one blank
    line; and in **Reconcile Command** replace the sentence that begins "It
    is read-only by default and carries three named mutation switches" with
    the last sentence of the block. No other line of `CONTEXT.md` changes.
12. MUST, as the last edit of the Task, after every guide edit, re-record the
    Behavior Surface Record with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
13. MUST NOT change the classification, reason or apply behavior of any Run
    that is not a candidate, the Delivery Queue's release, any ADR, skill or
    other guide; task_02 and task_03 document the rest.


The `CONTEXT.md` text of Requirement 11, in order (Branch Disposition
sentence; Dirty Discard entry; Reconcile Command sentence):

```markdown
A Dirty Discard writes a Branch Disposition that is not superseded and names the uncommitted paths it discards and the confirmed digest, before anything is removed ([ADR-0379](docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md)).

**Dirty Discard**:
The Reconcile Command's confirmed removal of an ended Run's Run Worktree and Run Branch when only uncommitted leftovers outside its archived Spec's declared and recorded paths keep the Run `dirty`. The Run ended `Failed`, `Unresolved`, `Stopped`, `BudgetExceeded` or `TimedOut` and owns no live process tree, its Spec's Archive Record is on the default branch head, and its commit proof holds, so no committed work is lost; the dry run lists every uncommitted path with one digest of the whole set, and `--discard-dirty <digest>` acts only on that digest, writing the Branch Disposition before it re-proves each Run and removes it ([ADR-0379](docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md)).
_Avoid_: dirty release, forced reconcile, manual reset

It is read-only by default and carries four named mutation switches, none of which has a force bypass: `--apply` releases freshly revalidated `safe` or `superseded` work; `--discard-superseded` records a Branch Disposition before removing a branch it can prove superseded; `--discard-dirty <digest>` performs a Dirty Discard; and `--carry-forward` performs Task Carry-Forward.
```

## Subtasks

- [ ] Mark candidates and add the re-proving discard in `internal/worktree`.
- [ ] Add the Dirty Discard record writer in `internal/daemon`.
- [ ] Add the switch, digest, refusals and report fields in `internal/cli`.
- [ ] Update the usage text, the help test, the user guide and `CONTEXT.md`.
- [ ] Add the three test files and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] An ended Run of an archived Spec kept only by an outside leftover is listed with its paths and one stable digest, and the confirmed digest records then removes it.
- [ ] A changed set, a missing candidate, a malformed digest and combined switches exit 2 and change nothing; a Run with a failed commit proof, an unarchived Spec, a Clean outcome or a live process tree is never offered.
- [ ] Reports without a candidate print the same bytes as before, and the glossary defines **Dirty Discard**.

## Context

- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/dirty_discard_test.go`
- interface: `internal/daemon/reconcile.go`
- creates: `internal/daemon/dirty_discard_record_test.go`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- creates: `internal/cli/reconcile_discard_dirty_test.go`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `internal/worktree/merged_head.go`
- instruction: `internal/worktree/merged_spec_leftovers_test.go`
- instruction: `internal/cli/reconcile_merged_test.go`
- instruction: `docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- `runworktree.RunWorktreeReconciliation` and `InspectTerminalRunMerged` — consumers: `./internal/worktree`, `./internal/daemon`, `./internal/cli`
- The `daemon.BranchDisposition` record — consumers: `./internal/daemon`, `./internal/cli`
- The `roundfix-reconcile/v1` report and the `reconcile` usage — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnEndedRunOfAnArchivedSpecWithAnOutsideLeftoverOffersADirtyDiscard|TestADirtyDiscardIsNeverOfferedOutsideItsBounds|TestDiscardDirtyTerminalRunRemovesTheRunWorktreeAndRunBranch|TestDiscardDirtyTerminalRunRefusesAChangedPathSet)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnEndedRunOfAnArchivedSpecWithAnOutsideLeftoverOffersADirtyDiscard TestADirtyDiscardIsNeverOfferedOutsideItsBounds TestDiscardDirtyTerminalRunRemovesTheRunWorktreeAndRunBranch TestDiscardDirtyTerminalRunRefusesAChangedPathSet; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/worktree` — expected: exit 0. Before this Task none of the four tests exists; afterwards they pass with the whole package.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestWriteDirtyDiscardRecordNamesTheDiscardedPaths$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestWriteDirtyDiscardRecordNamesTheDiscardedPaths ' || { echo 'missing pass: TestWriteDirtyDiscardRecordNamesTheDiscardedPaths' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'Branch|Disposition|Superseded' ./internal/daemon` — expected: exit 0. Before this Task the test does not exist; afterwards it passes with the existing disposition tests.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest|TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns|TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything|TestReconcileDiscardDirtyRefusesWhenNoRunIsOffered|TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree|TestReconcileDiscardDirtyNeedsItsDigestTranscript|TestReconcileDiscardDirtyRefusesAMalformedDigestTranscript|TestReconcileMutationSwitchesExcludeEachOtherTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything TestReconcileDiscardDirtyRefusesWhenNoRunIsOffered TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree TestReconcileDiscardDirtyNeedsItsDigestTranscript TestReconcileDiscardDirtyRefusesAMalformedDigestTranscript TestReconcileMutationSwitchesExcludeEachOtherTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'Reconcile|ReleaseMergedRuns|AutomaticRelease|TestRunCommandHelp' ./internal/cli` — expected: exit 0. Before this Task none of the eight tests exists; afterwards they pass with every reconcile, merged-release and help test.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- 'roundfix reconcile [run-id] [--apply | --discard-superseded | --discard-dirty <digest> | --carry-forward] [--format <text|json>]' || { echo 'missing synopsis in reconcile guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- 'Four switches mutate' || { echo 'missing switch count in reconcile guide' >&2; exit 1; }; grep -F -- '**Dirty Discard**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Dirty Discard\*\*(:| —( |$))' || { echo 'missing glossary term Dirty Discard' >&2; exit 1; }; grep -F -- '**Branch Disposition**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Branch Disposition\*\*(:| —( |$))' || exit 1; grep -F -- '**Reconcile Command**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reconcile Command\*\*(:| —( |$))' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A Dirty Discard writes a Branch Disposition that is not superseded' || { echo 'missing Branch Disposition change' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'carries four named mutation switches' || { echo 'missing Reconcile Command change' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestGlossaryClosedListsEqualTheirConstants|TestGlossaryAvoidLinesArePlainLists|TestTheSkillCoverageMapIsCurrent)$'` — expected: exit 0. Before this Task the guide and the glossary lack the new text; afterwards they hold it, the glossary contracts pass and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 5; User Experience; Success Metric 1; Success Metric 2; Glossary
- `_techspec.md` → Interfaces; Data Models; API Contract 1; API Contract 2; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Testing Approach; Build Order 1
- ADR-0379; ADR-0115; ADR-0212; ADR-0289


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging.

Implemented this Task's Dirty Discard slice. The pre-edit inspection found no
Dirty Discard outcome list, discard paths, record writer, digest or switch.
The only pre-existing changed path was this Task file (Daemon-owned status);
its status and authored content remain untouched by this implementation.

- Acceptance criterion 1: eligible ended Runs stay `dirty` with the original
  reason, while carrying all parsed paths and immutable inspection evidence.
  The CLI withholds live-process Runs, prints a stable digest, validates the
  confirmation before mutations, writes the Branch Disposition, and re-proves
  each Run before forced worktree removal and Run Branch deletion. The success
  test also compares Run Database bytes before and after and checks that the
  next scan reports `released`. Focused checks passed:
  `go test -count=1 -timeout 4m ./internal/cli -run 'DirtyDiscard|DiscardDirty'`
  and
  `go test -count=1 -timeout 4m ./internal/worktree -run 'DirtyDiscard|DiscardDirty|MergedSpecRunWithAnUndeclaredLeftover'`.
- Acceptance criterion 2: tests cover changed paths, changed Run/target/proof
  heads, changed inspection metadata and merged records, a missing candidate,
  malformed and absent digests, all mutation-switch pairs, failed commit proof,
  absent default-branch archive, Clean outcomes, scoped leftovers, and both
  live-process collections. Record-write failure preserves the Git surfaces.
  The two initial focused-check failures were fixture errors (the CLI fixture
  lacked archived `_prd.md`; the target-head fixture tried to move checked-out
  `main`); those fixtures were corrected and their checks passed. Additional
  focused checks passed:
  `go test -count=1 -timeout 4m ./internal/worktree -run '^TestDirtyDiscard(KeepsTheDeliveryProofWhenDefaultHasAdvanced|ReportsALockedWorktreeRemovalFailure)$'`
  and
  `go test -count=1 -timeout 4m ./internal/daemon -run '^TestWrite(DirtyDiscardRecordNamesTheDiscardedPaths|BranchDispositionRecordSurvivesDiscard)$'`.
  These prove a delivery proof remains pinned when the default branch advances,
  a locked-worktree removal failure names both remaining surfaces, invalid
  records write nothing, and the existing record-and-discard path still works.
- Acceptance criterion 3: a new wire-contract test pins every text and JSON
  byte of a report without Dirty Discard candidates. Existing report, help
  and superseded-discard checks also passed:
  `go test -count=1 -timeout 4m ./internal/cli -run '^TestRunCommandHelp$/^reconcile$|^TestRunReconcileDryRunReadOnly$|^TestRunReconcileJSONMatchesTextFields$|^TestReconcileDiscardWritesTheRecordBeforeRemoving$|^TestReconcileDiscardRefusesAnUnreachableCommit$'`.
  Updated both usage synopses, the help option and test row, and only the
  requested user-guide content. `CONTEXT.md` was checked against the Task's
  own Markdown block and the base revision: its three requested replacements
  match byte for byte, with no other line changed. Focused glossary checks
  passed with
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossary(ClosedListsEqualTheirConstants|StatesEachRuleOnce|AvoidLinesArePlainLists)$'`
  and
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs)$'`.

Windows compilation passed without running Windows executables:
`GOOS=windows go test -run '^$' -exec /usr/bin/true ./internal/worktree ./internal/daemon ./internal/cli`;
the worktree package was compiled again after its last two additional tests.
Every new top-level test starts with `t.Parallel()`. `git diff --check` passed.
Go checks used `rtk proxy`; the first sandbox attempts could not write the
host Go build cache, so subsequent checks used the approved cache access.

Final derived-record edit, executed after this Result and every guide edit:
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
Its exit status is in the session execution evidence; the generated record is
not edited by hand.

Declared Task Verification, repository Verification, status settlement and
committing remain Daemon-owned and were not run here. No commit, push, PR,
Task Graph, other Task, ADR, skill, Delivery Queue release or ignored-only
cleanup change is part of this diff. The latter documentation and cleanup
slices remain with task_02 and task_03.

## Carry-forward provenance

- Source Run: `run_20261109T214437Z_6567933336913b3f`
- Source commit: `84ea7bea3d958484aa73f603a540685909cac1c9`
