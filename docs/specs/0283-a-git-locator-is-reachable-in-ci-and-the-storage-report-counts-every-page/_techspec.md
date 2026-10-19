---
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
prd: _prd.md
created: 2026-10-19
---

# A Git Locator is reachable in CI, and the storage report counts every page — Technical Spec

## Executive Summary

Two slices extend existing seams, and no package is new.

1. The Default Branch Head (ADR-0295).
   - `internal/preflight` gains the `origin-branch` detection step and one
     resolver, `ResolveDefaultBranchHead`, that reads the fetched
     remote-tracking branch first.
   - `spec.ResolveGitLocator` uses it, so `SC-HISTORY-RETAINED` and every
     locator rewrite answer the same in a CI checkout and a developer clone.
   - The QA Delivery Base, the Archive Command's default reference and the
     pre-PR review's default base use it too.
2. The Storage Snapshot (ADR-0296).
   - `roundfix storage report` picks its reader the way `roundfix gc
     --dry-run` already does and measures inside one read transaction.
   - Its page groups add pointer-map pages and the lock-byte page, which
     `dbstat` omits, and reconcile with the snapshot exactly. The main, `-wal`
     and `-shm` files are reported apart.

The trade-offs accepted:

- A commit that only an unpushed local default branch holds no longer proves
  reachability.
- A repository that keeps both `origin/main` and `origin/master` without
  origin/HEAD stays undetermined, as today.
- The storage report's text gains five lines, and its tolerance reads 0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one detection source value, `origin-branch`, and
  five report lines are added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every new read is a local Git ref, a file
  stat or a read-only SQLite query, and every test uses temporary repositories
  and stores. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — each Build Order step follows recorded
  decisions. ADR-0295 governs Build Orders 1 and 2: "Reachability reads the
  fetched default branch, and origin names it when origin/HEAD is unset".
  ADR-0296 governs Build Order 3: "The storage report reconciles every page of
  one Storage Snapshot". ADR-0278: "A retirement deletes, and the History Root
  holds only Archive Records", which ADR-0295 refines. ADR-0276: "A row the
  authorization keeps offline qualifies, and a byte invariant compares with
  the Delivery Base", which ADR-0295 refines. ADR-0255: "Run Retention removes
  terminal Runs whole and compacts incrementally", which ADR-0296 refines.
  ADR-0172: "Doctor reports reclaimable Run storage from cheap reads", which
  ADR-0296 keeps. ADR-0169: "The pre-PR review diffs the candidate from its merge
  base". ADR-0261: "An Archive Record names its pull request, and readers fetch
  the pull request head". ADR-0104: "A Spec accepts on evidence it did not
  author". ADR-0265: "A Spec proves what it claims, and every active Spec is
  held to it". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0285: "Every test file compiles for Windows, and a repository
  contract proves it". ADR-0189: "An owned skill's version names its content,
  and its minimum is the bundle". ADR-0280: "A guide Task records its own
  surface fingerprint, a PRD records who decided, and the glossary is held to
  its own rule". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0271: "A test bounds by real time only what it
  asserts". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0236 (a failed check is judged on a merge with the current default branch); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0266 (a signal goes only to a proven process, and worktrees share one build cache); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ grant of 2026-10-19 ("Mexer em Governed Paths declarados") and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-19 at `f382480b` (v0.91.0).

- Default-branch detection (`internal/preflight/default_branch.go`).
  - `DetectDefaultBranch` (`:49`) returns origin/HEAD's branch when
    `symbolic-ref refs/remotes/origin/HEAD` names one, else a current branch
    named `main` or `master` (source `name-match`), else undetermined (`:63`).
    It swallows every Git error.
  - `TestDetectDefaultBranch` (`internal/preflight/default_branch_test.go:9`)
    pins those cases with a fake runner.
- Git Locator (`internal/spec/git_locator.go`).
  - `ResolveGitLocator` (`:91`) detects the default (`:98`) and then reads
    `refs/heads/<default>` before `refs/remotes/origin/<default>` (`:104`). A
    preferred commit is reachable when it is an ancestor of that head or a tag
    contains it (`:134`). Without a preferred commit, it names the newest
    commit of the head that holds the path.
  - `TestResolveGitLocatorReadsTheRemoteDefaultBranch`
    (`internal/spec/git_locator_test.go:143`) asserts "Local head takes
    precedence over an older origin head".
  - Callers: `internal/speccheck/history_retained.go` (`:134`, and
    `invalidHistoryLocators` at `:251`, `:289`), `internal/spec/record_links.go`,
    `internal/spec/archive.go`, `internal/spec/live_links.go`,
    `internal/spec/history_entries.go`, `internal/spec/history_sanitize.go`,
    `internal/cli/history.go`, `internal/cli/baseline_update_history.go` and
    `internal/cli/deliver_workflow.go:2006`. All read the default through
    `ResolveGitLocator`, so they inherit its fix.
- Reproduction. A throwaway `go test -overlay` probe built two temporary
  clones.
  - A CI-like clone: `git init`, `git fetch origin
    +refs/heads/*:refs/remotes/origin/* +refs/tags/*:refs/tags/*`, detached at
    a pull request head, no origin/HEAD.
  - A stale-local clone: origin/HEAD set, local `main` reset to the older tag.

  In both, a locator whose commit `origin/main` holds and no tag contains got
  no commit.
- The CI checkout (`.github/workflows/ci-verify.yml:26-31`, read-only):
  `actions/checkout@v7` with `fetch-depth: 0` and, for a pull request, `ref:`
  the head SHA, so HEAD is detached. `release.yml:23-30` also uses
  `fetch-depth: 0`.
- Other default-branch callers, swept.
  - `internal/daemon/qa_delivery_base.go:16` reads
    `refs/remotes/origin/<default>` first (`:31`), local second (`:37`), but
    returns unresolved when detection is undetermined.
  - `internal/cli/archive.go:242` `archiveDefaultRef` maps origin/HEAD to
    `refs/remotes/origin/<name>` and `name-match` to `refs/heads/<name>`
    (`:251`). `warnArchiveSourceRevision` (`:322`) and
    `recordArchiveCommandPullRequest` (`:260`) use it, so an unset origin/HEAD
    on a work branch prints "the default branch is unknown".
  - `internal/cli/review.go:1296` `resolveReviewBaseCommit` uses the short
    name for `name-match` (`:1311`), which Git resolves to the local branch,
    and falls back to `refs/heads/main` then `refs/heads/master` (`:1301`)
    when detection is undetermined.
  - Unchanged, with reasons:
    - `internal/cli/implement.go:284` asks whether the current branch is the
      default (Preflight Validation), not what it reaches. A detached HEAD and
      a work branch still match no default.
    - `internal/worktree/worktree.go:1212` validates a local branch target and
      needs a local branch by design.
    - `internal/cli/deliver.go:277` and `internal/cli/deliver_workflow.go`
      (`:475`, `:535`, `:799`, `:843`, `:1546`, `:2216`, `:2277`) fetch the
      push remote and read `refs/remotes/<remote>/<default>`, or run without a
      remote. They never read a stale local branch, and they inherit the
      `origin-branch` step instead of refusing.
    - `internal/cli/deliver_linux_lane.go:82` and `:296` read the fetched
      remote the same way. That file belongs to Spec 0282 in this cycle.
- Storage report (`internal/store/journal.go`, `internal/cli/gc.go`).
  - `StorageReport` (`journal.go:228`) runs separate implicit reads of
    `page_size`, `freelist_count`, `dbstat` and the Run rows. It stats the main
    file and refuses when `|file − (dbstat + free)| > page_size` (`:263`),
    with the reason constant at `:161`.
  - `OpenStorageReader` (`store.go:423`) opens `mode=ro&immutable=1`
    (`storageReaderDSN`, `:490`), which skips locking and the WAL.
  - `measureStorageReport` (`gc.go:282`) always uses that reader, for
    `storage report` and for `gc compact`'s live-writer fallback (`:434`).
    `openGCStore` (`:880`) already picks the immutable reader only when the
    `-wal` file is absent or holds 32 bytes or fewer, and `OpenReader`
    otherwise.
  - Measured, read-only: this machine's idle store (WAL 0 bytes) fails with
    difference 1,380,352 = 337 pages, which is 336 pointer-map pages plus the
    lock-byte page. A temporary WAL database with 2,000 rows not yet
    checkpointed measured 3,013 pages immutable and 5,020 through a WAL
    reader, and 4 and 7 pointer-map pages that `dbstat` omitted. Any store
    past about 1,640 pages of 4096 bytes already exceeds the one-page
    tolerance.
  - `TestStorageReportReconcilesMeasuredTotalsWithoutMutation`
    (`internal/store/journal_test.go:736`) asserts the tolerance equals the
    page size and the reason holds "one SQLite page";
    `TestRunStorageReportOutsideGitRepository` (`internal/cli/gc_test.go:528`)
    asserts the output holds "one SQLite page".
- Documentation.
  - `docs/user-guide/commands/archive.md:24` says the locator prefers
    `source_revision` "when the local default branch or a tag reaches it".
  - `docs/user-guide/commands/review.md:95` and the Roundfix `review`
    reference (`:85`) say the default base is "the repository's default
    branch".
  - `docs/user-guide/commands/storage-report.md` and the Roundfix `storage`
    reference describe the report without its page groups.

## System Architecture

- `internal/preflight`: `default_branch.go` gains `DefaultBranchFromOriginBranch`
  and the `origin-branch` step; a new `default_branch_head.go` holds
  `DefaultBranchHead` and `ResolveDefaultBranchHead`.
- `internal/spec`: `git_locator.go` reads the Default Branch Head.
- `internal/daemon`: `qa_delivery_base.go` reads it.
- `internal/cli`: `archive.go` and `review.go` read it; `gc.go` shares one
  snapshot-reader choice between `openGCStore` and `measureStorageReport` and
  prints the new lines.
- `internal/store`: `journal.go` measures inside one read transaction and adds
  the pointer-map and lock-byte groups.

## Implementation Design

### Interfaces

```go
// internal/preflight
const DefaultBranchFromOriginBranch DefaultBranchSource = "origin-branch"
type DefaultBranchHead struct {
	Branch DefaultBranch
	Ref    string // refs/remotes/origin/<name> or refs/heads/<name>
	Commit string
}
func ResolveDefaultBranchHead(ctx context.Context, workDir, currentBranch string, runner GitRunner) (DefaultBranchHead, bool, error)

// internal/store: StorageReport gains
//   DatabaseWALBytes, DatabaseSharedMemoryBytes, DatabaseSnapshotBytes,
//   DatabasePointerMapBytes, DatabaseLockBytePageBytes int64
```

### Data Models

No file format, schema or record changes. `StorageReport` is an in-process
value; its new fields are printed only by `storage report`.

### API Contracts

1. API Contract: the `origin-branch` detection step.
   - `DetectDefaultBranch` keeps its first two steps. When both fail, it runs
     one `git for-each-ref --format=%(refname) refs/remotes/origin/main
     refs/remotes/origin/master` and keeps only exact ref names.
   - Exactly one match returns that branch with source `origin-branch`. Zero
     or two matches, or any Git error, return undetermined. The step never
     runs when an earlier step answered.
2. API Contract: `ResolveDefaultBranchHead`.
   - It reads the current branch the caller passes, detects the default, and
     runs one `git for-each-ref --format=%(refname)%00%(objectname)
     refs/remotes/origin/<name> refs/heads/<name>`.
   - The exact `refs/remotes/origin/<name>` wins; `refs/heads/<name>` answers
     only when the remote-tracking ref is absent.
   - It returns `found=false` with no error when detection is undetermined or
     neither ref exists. A Git error is returned wrapped with `%w`.
3. API Contract: the Git Locator.
   - `ResolveGitLocator` reads the head from `ResolveDefaultBranchHead` with
     its existing runner. That head proves reachability, with tags as today,
     and supplies the newest commit that holds a removed path.
   - Its signature, error wrapping and `Text` forms stay.
   - `SC-HISTORY-RETAINED` accepts a locator whose commit the head reaches,
     with no other change.
4. API Contract: the QA Delivery Base. `qaDeliveryBase` merges `HEAD` with the
   head's `Ref`. `found=false` still leaves the base unresolved, and a Git
   error is still returned.
5. API Contract: the Archive Command's default reference.
   - `archiveDefaultRef` returns the head's `Ref` and displays `origin/<name>`
     for a remote-tracking ref or `<name>` for a local one.
   - `found=false` keeps today's texts: "the default branch is unknown, so
     source_revision … was not checked" and "the default branch is unknown;
     cannot record a pull request".
6. API Contract: the pre-PR review's default base.
   - Without `--base`, `resolveReviewBaseCommit` uses the head's `Ref`.
   - `found=false` keeps today's `refs/heads/main`, then `refs/heads/master`,
     fallback and its error text.
   - `--base <ref>` is unchanged.
7. API Contract: the Storage Snapshot reader.
   - One helper in `internal/cli/gc.go` returns the immutable storage reader
     when `<database>-wal` is absent or 32 bytes or fewer, and the read-only
     WAL reader (`store.OpenReader`) otherwise. A stat error other than
     not-exist is returned.
   - Both `openGCStore`'s dry run and `measureStorageReport` use it, so
     `storage report` and `gc compact`'s live-writer fallback read the same
     snapshot.
8. API Contract: the page groups.
   - `StorageReport` runs every query on one connection inside one read
     transaction: `page_size`, `page_count`, `freelist_count`, `auto_vacuum`,
     `dbstat` and the Run rows.
   - The usable size is the page size minus the reserved-bytes field at offset
     20 of the main file's header.
   - Pointer-map pages exist when `auto_vacuum` is not 0. For `J = usable/5`,
     candidate position `k` is `2 + k·(J+1)`, moved to the next page when it
     equals the lock-byte page. Their count is the number of positions that
     do not exceed `page_count`.
   - The lock-byte page is `1073741824/page_size + 1`, counted when
     `page_count` reaches it.
   - Snapshot bytes are `page_count·page_size`. Grouped bytes are table, free,
     pointer-map and lock-byte bytes. A difference other than 0 refuses with
     `reconcile Run Database storage groups: snapshot bytes=<n> grouped
     bytes=<n> difference=<n>`.
   - `ReconciliationToleranceBytes` is 0, with the reason `one read
     transaction measures every page group, pointer-map and lock-byte pages
     included`.
   - `DatabaseBytes` stays the main file's size. The `-wal` and `-shm` sizes
     are stat'ed, and an absent file is 0.
9. API Contract: the report text. The Database section prints, in order: the
   path, the file bytes, `WAL bytes`, `Shared-memory bytes`, `Snapshot bytes`,
   the table bytes, the free bytes, `Pointer-map bytes`, `Lock-byte page bytes`,
   the grouped bytes, the difference between snapshot and grouped bytes, and
   the tolerance line. Every other section is unchanged.
10. API Contract: everything else stays. These keep their behavior: Preflight
    Validation's default-branch veto, `DefaultBranch.IsDefault`, the Delivery
    Queue and Linux Lane reads, `roundfix gc compact --apply`, the Doctor
    storage check, the compaction preview's own tolerance, every existing
    error text not named above, and every existing test name.

### Surface Transcripts

1. Surface Transcript: storage report without a Run Database.

   Run with `HOME` set to a temporary directory that holds no Run Database;
   `<home>` is that directory. Tabs separate the column headings.

   ```transcript
   $ roundfix storage report
   stdout:
   Storage report
   Database
     Path: <home>/.roundfix/roundfix.db
     File bytes: 0
     WAL bytes: 0
     Shared-memory bytes: 0
     Snapshot bytes: 0
     Table bytes: 0
     Free bytes: 0
     Pointer-map bytes: 0
     Lock-byte page bytes: 0
     Grouped bytes: 0
     Reconciliation difference: 0
     Reconciliation tolerance: 0 bytes (no Run Database exists)
   Tables
     TABLE	ROWS	BYTES
   Repositories
     REPOSITORY	STATUS	RUN_ROWS	RUN_ARTIFACT_BYTES
   States
     STATE	RUN_ROWS	RUN_ARTIFACT_BYTES
   Artifact Roots
     ARTIFACT_ROOT	STATUS	RUN_ROWS	BYTES
     Measured bytes: 0
     Grouped bytes: 0
     Reconciliation difference: 0
   stderr:
   exit: 0
   ```

### Invariants

```text
1. In one repository state, a Git Locator is valid in a CI checkout exactly when it is valid in a clone whose local default branch is stale.
2. No reachability or merge-base read answers from a local default branch while the fetched remote-tracking branch exists.
3. A storage report that exits 0 has grouped bytes equal to its snapshot bytes.
4. The storage report changes no byte of the database or WAL files and creates no sidecar file.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Features 1-2, Success Metrics 1-2 → `DetectDefaultBranch`, `ResolveDefaultBranchHead`, `ResolveGitLocator` (API Contracts 1-3)
- Goal 2, User Story 2, Core Feature 3, Success Metric 3 → the QA Delivery Base, the Archive Command and the review base (API Contracts 4-6)
- Goal 3, User Story 3, Core Feature 4, Success Metrics 4-5 → the snapshot reader, the page groups and the report text (API Contracts 7-9, Surface Transcript 1)
- Core Feature 5 → the `archive`, `review` and `storage report` guides, the Roundfix references, `CONTEXT.md`
- Success Metrics 1-5, API Contract 10 and Invariants 1-4 → the named tests and the scope audit

## Absence Evidence

1. The `origin-branch` detection source — command: `grep -rn 'origin-branch' internal/preflight`; result: no match at `f382480b`.
2. The Default Branch Head resolver — command: `grep -rn 'func ResolveDefaultBranchHead' internal`; result: no match at `f382480b`.
3. The pointer-map and lock-byte page groups — command: `grep -rni 'pointer-map\|ptrmap\|lock-byte' internal/store internal/cli`; result: no match at `f382480b`.
4. The shared snapshot-reader choice for the storage report — command: `grep -n 'openStorageReader' internal/cli/gc.go`; result: `measureStorageReport` calls `openStorageReader` unconditionally at `:292`; only `openGCStore` (`:880`) checks the `-wal` size.

## ADR Alternatives

- ADR-0295: this Spec's decision; it rejects `git remote set-head`, any remote-tracking branch as proof, the remote ref as a second candidate, preferring `main` when both exist, and CI environment variables.
- ADR-0296: this Spec's decision; it rejects a wider tolerance, counting the omitted pages by subtraction, always opening a WAL reader, and checkpointing before measuring.
- ADR-0278: refined by ADR-0295, which defines the default branch head its Git Locator reads; its rejected alternatives stay rejected, and live work branches still never prove reachability.
- ADR-0276: refined by ADR-0295; the Delivery Base keeps its definition and now reads the fetched head; no rejected alternative is taken.
- ADR-0255: refined by ADR-0296, which accounts for the pointer-map pages its incremental auto-vacuum creates; no rejected alternative is taken.
- ADR-0172: kept; the Doctor storage check still runs no `dbstat`, so its rejected full measurement stays out of Doctor.
- ADR-0169: no rejected alternative is taken; the review still diffs from the merge base, and only the default base ref's source changes.
- ADR-0261: no rejected alternative is taken; the archive warning still names the pull request head ref, against the fetched default.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on a measurement this Spec did not design, repository history and published sources.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0285: no rejected alternative is taken; every new test compiles for Windows.
- ADR-0189: no rejected alternative is taken; each skill edit raises its recorded version.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0256: no rejected alternative is taken; each changed surface's skill text changes with it.
- ADR-0271: no rejected alternative is taken; no new test bounds a wait by real time.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0236: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0240: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0266: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

None outside the repository. The resolver reads local refs, and the storage
report reads the local Run Database read-only.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, builds temporary
repositories or stores, compiles for Windows, and reaches no network.
Repositories that need a private Git environment reuse the package's existing
child-process isolation (`gitLocatorTestHome`, `retainedTestHome`).

- `internal/preflight/default_branch_test.go`: new cases in
  `TestDetectDefaultBranch` for the `origin-branch` step (one, two and zero
  remote branches, and a Git error), and the new
  `TestResolveDefaultBranchHeadPrefersTheFetchedBranch` (fake runner: remote
  wins, local answers only without a remote, undetermined and absent refs
  return `found=false`, a Git error is returned).
- `internal/spec/git_locator_test.go`:
  - `TestResolveGitLocatorReachesTheFetchedDefaultInACIClone`: the CI-like
    clone of Current behavior;
  - `TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch`: the stale-local
    clone;
  - `TestResolveGitLocatorReadsTheRemoteDefaultBranch` changes its last
    assertion: an unpushed local commit on the default branch no longer wins
    over the fetched head (declared break, ADR-0295).
- `internal/speccheck/history_retained_test.go`:
  `TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone`, an
  Archive Record whose locator commit lies between two tags.
- `internal/daemon/qa_every_run_audit_test.go`:
  `TestQADeliveryBaseReadsTheFetchedDefaultWithoutOriginHead`.
- `internal/cli/archive_pull_request_test.go`:
  `TestArchiveChecksTheRevisionAgainstOriginWithoutOriginHead`, a full
  transcript whose warning names `origin/main`.
- `internal/cli/review_merge_base_test.go`:
  `TestReviewDefaultBaseIsTheFetchedDefaultBranch`, a work branch without
  origin/HEAD whose local `main` is behind `origin/main`.
- `internal/store/journal_test.go`:
  - `TestStoragePageStructureMatchesTheFileFormat`, table cases: 271,903 and
    275,136 pages of 4096 bytes give 333 and 337 pages; 3,013 and 5,020 give 4
    and 7; `auto_vacuum` 0 gives only the lock-byte page; page size 1024 moves
    the pointer-map page that falls on page 1,048,577.
  - `TestStorageReportGroupsPointerMapPages`, a store of more than 1,700 pages;
  - `TestStorageReportReadsALiveWALSnapshot`, a writer that keeps WAL frames
    open while the report reads through the WAL reader, with database and WAL
    bytes unchanged;
  - `TestStorageReportReconcilesMeasuredTotalsWithoutMutation` changes its
    tolerance assertions to 0 and the new reason (declared break, ADR-0296).
- `internal/cli/gc_test.go`:
  - `TestRunStorageReportWithoutDatabaseTranscript` asserts the Surface
    Transcript whole;
  - `TestRunStorageReportWhileAWriterHoldsTheWAL` asserts exit 0, a positive
    `WAL bytes` line and `Reconciliation difference: 0`;
  - `TestRunStorageReportOutsideGitRepository` replaces "one SQLite page" with
    the new reason (declared break).

## Build Order

1. The Default Branch Head and the Git Locator: the detection step, the
   resolver, `ResolveGitLocator`, their tests and the changed locator test,
   the `archive` guide's locator sentence and the Roundfix `archive`
   reference, the glossary, the skill version and the Behavior Surface Record
   (no dependency).
2. The sweep: the QA Delivery Base, the Archive Command's default reference
   and the review base, their tests, the `archive` and `review` guides and
   references, the glossary, the skill version and the record (depends on: 1,
   which adds the resolver and edits the same `archive` guide and reference).
3. The Storage Snapshot: the reader choice, the page groups, the report text,
   their tests and the two changed tests, the `storage report` guide and the
   Roundfix `storage` reference, the glossary, the skill version and the
   record (depends on: 2, which raises the same Roundfix skill version and
   re-records the same files).
4. QA gate (depends on: 3).

No other Spec is a prerequisite. Specs 0281 and 0282 are delivered first in
this cycle. They change the Roundfix skill version record, the Behavior
Surface Record and possibly `CONTEXT.md`, which each Task here regenerates
from its own edits, and no source file this Spec declares.

## Risks & Considerations

- A developer clone whose `origin/main` was never fetched answers from the
  local branch, as before.
- A repository whose remote is not called `origin` keeps today's detection;
  the Delivery Queue's push-remote reads are unchanged.
- A WAL that a writer resets between the size check and the open is read as a
  WAL reader, whose sidecar files already exist.
- Steps 1 to 3 each raise the Roundfix skill version and re-record shared
  files, so they run in series.

## Glossary

- adds: **Default Branch Head**
- adds: **Storage Snapshot**
- changes: **Git Locator**
- changes: **Delivery Base**

## Decisions

- origin names the default only when exactly one of `origin/main` and
  `origin/master` exists. See ADR-0295.
- One resolver serves every reachability and merge-base read, and the fetched
  branch answers first. See ADR-0295.
- The storage report reuses `gc --dry-run`'s reader choice, measures in one
  read transaction, and groups pointer-map and lock-byte pages with a 0
  tolerance. See ADR-0296.
- `for-each-ref` reads refs because it exits 0 when a ref is absent, so absence
  never needs an exit-status test through the `GitRunner` interface.
