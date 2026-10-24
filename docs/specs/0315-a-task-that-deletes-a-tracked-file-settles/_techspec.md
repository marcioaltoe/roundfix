---
spec: 0315-a-task-that-deletes-a-tracked-file-settles
prd: _prd.md
created: 2026-10-24
---

# A Task that deletes a tracked file settles — Technical Spec

## Executive Summary

`CitationSources` (`internal/spec/citation_sources.go`) lists Citation
Sources from the Git index and checks each path on disk with
`historyRewriteRegularFile`. A missing path is skipped only under
`docs/history`; anywhere else the whole call fails with `inspect citing file
<path>: lstat …: no such file or directory`. The strict Spec check calls it in
its retired-history pass, so a Task whose Agent deleted a tracked file without
staging the deletion fails the Spec Consistency settlement check. Two Markdown
inventories of the same shape, `HistoryCitations` and
`StaleHistoryCitations`, fail the same way for a tracked Markdown file the
working tree lacks. The fix skips, in those three readers, a tracked path
whose working-tree entry does not exist (`os.ErrNotExist`), wherever it is.
Every other read failure keeps failing, so a path behind a file, a symbolic
link and a binary file keep today's handling.

The trade-off accepted: a tracked file deleted from disk without staging is
no longer scanned for citations. Nothing is lost: the file will leave the
tree with the Task's commit, and a citation of the deleted file from a file
that remains is still reported.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; tests use temporary repositories and the
  existing fake Agent runner. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no accepted ADR decides whether a
  reader of tracked files skips a missing path, as `_prd.md` → Project
  Constraints states. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling file and no
  Governed Path changes; the Spec-contained authorization record
  `docs/specs/0315-a-task-that-deletes-a-tracked-file-settles/_authorization.md`
  records `paths: []`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Current behavior

Read and characterized on 2026-10-24 at `64c56d7f`. Line numbers move when
other Specs merge; the function names do not.

1. `CitationSources` (`internal/spec/citation_sources.go:39`) runs `git
   ls-files -z --stage --cached` (plus `--others --exclude-standard` when
   asked), keeps regular-file index entries, drops skill trees and testdata,
   skips a symbolic link, and calls `historyRewriteRegularFile`, which
   `Lstat`s the path and each parent. On an error it continues only when the
   error is `os.ErrNotExist` and the path is under `docs/history` (`:69`);
   otherwise it returns the error.
2. `detectHistoryRetained` (`internal/speccheck/history_retained.go`) calls
   `spec.CitationSources(ctx, root, true)` (`:260`) when Git is present and
   `spec.RemovedHistoryPaths` returns a candidate, that is, in every
   repository that removed a file under `docs/history/{findings,backlog,adr}`.
   The Daemon's Spec Consistency settlement check
   (`internal/daemon/settlement_spec_consistency.go`) reaches it through
   `SpecCheckSettlementChecker.RefusingFindings` → `specGatePrecondition` →
   `speccheck.Check`, and wraps the error as `read refusing Spec Consistency
   findings: run "spec check --strict" for Spec <slug>: …`.
3. The Daemon never stages before its settlement checks: the Agent's
   deletion stays an unstaged index entry until the Task commit. The deletes
   check (`undeletedTaskPaths`, `internal/daemon/task_deletes.go`) requires the
   path to be absent from disk; `prepareTaskCommit` stages the tracked
   removal and records undeclared paths; the authorization settlement check
   audits Governed Paths among the staged paths.
4. `HistoryCitations` (`internal/spec/history_entries.go:516`) lists tracked
   Markdown with `git ls-files -z -- *.md` and continues past a missing file
   only under `docs/history/` (`:546` and `:552`). Callers: `roundfix cancel`
   (`internal/cli/cancel.go:299`) and the archive plan
   (`PlanArchiveRelink`, `internal/spec/live_citations.go`, which already
   ignores an `os.ErrNotExist` from it).
5. `StaleHistoryCitations` (`internal/spec/stale_citations.go:16`) lists
   tracked Markdown the same way and fails on any read error (`:32`). Caller:
   the history conversion of `roundfix baseline update`
   (`internal/cli/baseline_update_history.go:294`).
6. Characterization (scratch clone of `64c56d7f`, `HOME` at a temporary
   directory): `bin/roundfix spec check
   0312-the-ci-job-runs-the-documentation-contracts-once --strict` exits 0;
   after `rm internal/agent/acp_stream.go` it fails with `roundfix: spec check
   failed: inspect citing file internal/agent/acp_stream.go: lstat …: no such
   file or directory`.
7. Characterization through the settlement path: a prototype of task_02's
   `TaskCycle` test (gated plan, default settlement checker, retirement rule
   recorded, one Removed History Path, Task declaring `deletes:
   packages/backend/test/gone.test.ts`, fake Agent removing it from disk)
   fails the Task with `Settlement check failed: settlement check: spec
   consistency: read refusing Spec Consistency findings: run "spec check
   --strict" for Spec 0001-sample-feature: inspect citing file
   packages/backend/test/gone.test.ts: lstat …: no such file or directory`,
   oraculum's message. Staging the deletion with `git rm` settles, as
   oraculum observed. Without a Removed History Path the pass does not run
   and the Task settles.

Sweep of every reader that lists tracked files and then reads them from the
working tree (`rg -n 'ls-files' internal -g '!*_test.go'`):

- Changed: `CitationSources`, `HistoryCitations`, `StaleHistoryCitations`.
- `SourceHistoryCitations` and `PlanHistoryRewrite` (`roundfix history
  sanitize`) read through `CitationSources`; they follow it.
- Unchanged, already skips a missing file: the superseded-record plan in
  `internal/spec/history_entries.go` (`ls-files … *.md`, `os.Lstat` with
  `os.ErrNotExist` → continue).
- Unchanged, reads Git for a missing file on purpose: the archive plan's Live
  Document reader (`internal/spec/live_citations.go`), which takes the HEAD
  blob when the working-tree file is gone.
- Unchanged, already skips: the Baseline relocation citation scan
  (`internal/baseline/history_citations.go` `listTrackedPaths` →
  `readTrackedCitationFile`, which returns `trackedReadSkipped` when the
  `Lstat` fails).
- Unchanged, reads blobs from the index: the mechanical QA evidence reader
  (`internal/speccheck/mechanical.go` `mechanicalTrackedEvidence`).
- Unchanged, not on a Task path: `internal/cli/deliver_workflow.go` lists
  tracked files to snapshot derived lines during a Delivery Queue merge
  resolution, after every Task commit.
- Unchanged, list paths only: `indexTracksExactPath`
  (`internal/daemon/task_engine.go`) and `verifyselect`'s untracked listing.

## System Architecture

- `internal/spec/citation_sources.go`, `internal/spec/history_entries.go`,
  `internal/spec/stale_citations.go`: the missing-entry skip.
- `internal/spec/citation_sources_working_tree_test.go` (new) and
  `internal/daemon/task_deletes_settlement_test.go` (new): the proofs.

No package, import edge, exported symbol, message or schema changes.

## Implementation Design

### Interfaces

No signature changes:

```go
func CitationSources(ctx context.Context, root string, withUntracked bool) ([]string, error)
func HistoryCitations(ctx context.Context, repositoryRoot string, removed []string) ([]HistoryCitation, error)
func StaleHistoryCitations(ctx context.Context, repositoryRoot, historyRoot string, held map[string]bool) ([]HistoryCitation, error)
```

### Data Models

None changes.

### API Contracts

1. API Contract: `CitationSources` continues past a listed path when
   `historyRewriteRegularFile` fails with an error that `errors.Is(err,
   os.ErrNotExist)`, wherever the path is; every other error from it, from
   opening the file or from reading its prefix is returned as today.
2. API Contract: `HistoryCitations` continues past a listed Markdown path
   when `historyRewriteRegularFile` or `os.ReadFile` fails with
   `os.ErrNotExist`, wherever the path is; every other error is returned as
   today.
3. API Contract: `StaleHistoryCitations` continues past a listed Markdown
   path when `os.ReadFile` fails with `os.ErrNotExist`; every other error is
   returned as today, with today's message.
4. API Contract: unchanged consumers and checks. The strict Spec check's
   findings and messages, the Spec Consistency settlement check, the deletes
   check, `prepareTaskCommit`, the authorization settlement check, `roundfix
   cancel`, `roundfix history sanitize`, `roundfix baseline update` and the
   readers listed as unchanged in Current behavior keep their code and bytes.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes; the
strict check, the settlement, `cancel`, `history sanitize` and `baseline
update` stop failing for a tracked file the working tree lacks and otherwise
print what they print today.

### Invariants

```text
1. A tracked path whose working-tree entry does not exist is never read and never fails the three readers.
2. A path behind a file, a symbolic link, a binary file and any read error other than a missing entry keep today's result.
3. A gated Task with a declared deletion still passes the deletes check, and a deletion that breaks a live link or touches an unauthorized Governed Path is still reported or refused.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Features 1 and 3, Success Metrics 1-2 → API Contract 1, Invariant 1, task_01 and task_02
- Goal 2, Success Metric 3, Invariant 3 → API Contract 4, task_02
- Goal 3, Success Metric 4, Invariant 3 → API Contract 4, task_02
- Goal 4, Core Feature 2, Success Metric 5 → API Contracts 2-3, Invariant 1, task_01
- Goal 5, Success Metric 6, Invariant 2 → API Contracts 1-3, task_01
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. A test of a tracked file missing on disk for these readers — command: `grep -rn -e TrackedSourceReaders -e DeclaredDeletionIsOnlyOnDisk internal`; result: no match at `64c56d7f`; the only missing-file skip is the `docs/history` one in `CitationSources` and `HistoryCitations`.

## ADR Alternatives

None. No ADR is listed: the fix narrows three readers to the glossary's
Citation Source and records no new decision.

## Integration Points

None outside the repository: Git on temporary repositories and the existing
fake Agent runner and verifier. No test, Verification command or QA row reads
the real `~/.roundfix`, another repository or the network.

## Testing Approach

Every new top-level test and subtest calls `t.Parallel()` first; fixtures are
synthetic. A callback that runs on the Agent's turn reports with `t.Error`,
never `t.Fatal`, and tolerates a file already gone, because a repair turn runs
it again.

- `internal/spec/citation_sources_working_tree_test.go` (new), one fixture: a
  temporary repository tracking `README.md`, `docs/gone.md`,
  `pkg/gone.test.ts` and `notes/a.md`, each citing
  `docs/history/findings/x.md`, committed. Subtests `citation-sources`,
  `history-citations` and `stale-history-citations` per test:
  - `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks`: after removing
    `docs/gone.md` and `pkg/gone.test.ts` from disk only, each reader returns
    no error and exactly `README.md` and `notes/a.md` (the sources, or the
    citing paths).
  - `TestTrackedSourceReadersStillRefuseAPathBehindAFile`: after replacing
    the directory `notes` with a file, each reader returns an error naming
    `notes/a.md`.
- `internal/daemon/task_deletes_settlement_test.go` (new), one fixture built
  on `settlementRepositoryFixture`: the recorded retirement rule
  (`docs/agents/setup-context.json` and a `guide.docs-layout` region holding
  `spec.RetirementRuleSentinel`), `docs/history/findings/2026-01-01-old.md`
  committed then removed with `git rm` and committed (one Removed History
  Path), `docs/user-guide/index.md` linking `guide.md`,
  `docs/user-guide/guide.md`, `packages/backend/test/gone.test.ts`, and
  task_01 declaring `deletes: packages/backend/test/gone.test.ts`:
  - `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles`, subtests
    `disk-only` (the fake Agent removes the file from disk) and `staged` (it
    runs `git rm --ignore-unmatch`): `TaskCycle` completes the Task.
  - `TestSettlementStillReportsALinkToADeletedFile`: after removing
    `docs/user-guide/guide.md` from disk, `SpecCheckSettlementChecker`
    `RefusingFindings` returns no error and an `SC-HISTORY-RETAINED` finding
    whose summary names `docs/user-guide/index.md` and `guide.md`.
  - `TestSettlementRefusesAnUnauthorizedGovernedDeletion`: a repository whose
    authorization record lists only `Makefile` and which tracks
    `.agents/skills/escaped/SKILL.md`; after removing that file from disk,
    `AuditCommit` for a prospective commit changing it returns exactly one
    `QA-AUTH-PATHS` finding naming it.
- Unchanged and run by name: `TestATaskThatLeavesItsDeletesPathDoesNotSettle`,
  `TestATaskThatRemovesItsDeletesPathSettles`,
  `TestDefaultSettlementCheckerAuditsAProspectiveCommit`,
  `TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings` and
  `TestStaleHistoryCitationsReportsFilesTheTagHolds`.
- A prototype passed these tests at authoring time; before the fix the
  `disk-only` subtest, `TestSettlementStillReportsALinkToADeletedFile` and
  every `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks` subtest failed,
  while `staged`, the authorization test and the path-behind-a-file test
  passed. Three sabotages, each widening one reader's skip to any error, each
  failed exactly its reader's subtest of
  `TestTrackedSourceReadersStillRefuseAPathBehindAFile`.

## Build Order

1. The missing-entry skip in the three readers and their tests in
   `internal/spec` (task_01; no dependency).
2. The settlement-path tests in `internal/daemon` (task_02) (depends on: 1,
   because its `disk-only` subtest passes only with the readers fixed).
3. QA gate (task_03) (depends on: 1, 2).

No Spec is listed in `requires`.

## Risks & Considerations

- A tracked file deleted from disk but kept in the index on purpose is no
  longer scanned for citations; its text leaves with the next commit, and
  Git keeps it.
- A Task that deletes a file cited by a live document now gets the citation
  finding at settlement instead of an opaque error, so its repair turn sees
  what to fix.
- The Delivery Queue's revalidation and the QA precondition run the same
  strict check and stop failing for the same reason.

## Glossary

None.

## Decisions

- Skip only a missing working-tree entry, in the three readers that list
  tracked files and then read them; keep listing from the index.
- Leave the readers that already skip a missing file or read it from Git.
- No guide, skill, glossary or ADR change; the glossary's Citation Source
  is already a tracked regular text file.
