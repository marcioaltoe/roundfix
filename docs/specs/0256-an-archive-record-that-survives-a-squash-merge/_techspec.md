---
spec: 0256-an-archive-record-that-survives-a-squash-merge
prd: _prd.md
created: 2026-10-12
---

# An Archive Record that survives a squash merge — Technical Spec

## Executive Summary

Three changes, one per reader class, plus the guidance. First, a single
predicate decides what an Archive Record is, and every reader of the archive
root uses it, so a prose `README.md` is ignored. Second, one recovery function
in `internal/spec` makes a missing `source_revision` readable by fetching
`refs/pull/<N>/head`. The four readers that open the folder at
`source_revision` call it before their existing availability check. Third, the
Archive Command gains `--pull-request <n>` and a squash warning, and the
Delivery Queue uses that flag to add a Pull Request Record Commit after it
opens the pull request. The primary trade-off is that a reader can now reach
the network. We accept it because the alternative, pushing an archive ref or
forcing merge commits, was rejected by the maintainer. We bound it: the fetch
is non-interactive, time-boxed, writes no ref and runs only when the commit is
missing and a number is known (ADR-0261). The Delivery Queue accepts the
record commit by proof instead of re-running the repository gate, because the
commit changes one scalar of one record.

## Project Constraints

- Identifier strategy: applicable — the record keeps `roundfix/archive-record/v1`
  and its existing optional `pull_request` field, a decimal string matching
  `^[1-9][0-9]{0,9}$`. An Archive Record is a file directly under the archive
  root whose name matches `^[0-9]{4}-[^/\\]+\.md$` and whose front matter
  declares the schema. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Source Revision Recovery runs
  `git fetch` against a configured remote with `GIT_TERMINAL_PROMPT=0`, no
  standard input and a 60-second time bound that cancels the process. It
  writes objects only and adds no credential handling. Tests use a local bare
  repository as the remote. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0261 is implemented here and
  extends ADR-0247. ADR-0260 (Spec 0255, merged first) makes the owner archive relink Live
  Documents and record them under `relinked:`; the Pull Request Record Commit
  changes only `pull_request`, never `relinked` or a relinked file, and the
  queue's exact-archive proof keeps accepting the relinked files it already
  proves. ADR-0223, ADR-0229 and ADR-0237 bind retries, which must
  adopt a Pull Request Record Commit instead of repeating it. ADR-0232 binds
  merge evidence, which reads the merged record. ADR-0165 and ADR-0169 bind
  the review's archived-Spec handling. ADR-0184 binds the Surface
  Transcripts. ADR-0187, ADR-0189 and ADR-0233 bind the owned skills'
  versions. ADR-0256 binds the Skill Coverage Map. ADR-0259 binds parallel
  tests. ADR-0252 and ADR-0253 bind the contract run. ADR-0244 binds the
  glossary terms. ADR-0179 binds the Governed Paths. The rest of the PRD's
  row applies unchanged. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 alone changes Governed Paths, and
  the grant bounds exactly those. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0256-an-archive-record-that-survives-a-squash-merge/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `.agents/skills/roundfix/references/runs.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## System Architecture

No new package. The Spec extends `internal/spec` (record identity, recovery,
pull request field), `internal/cli` (archive, review, runs causes, delivery
workflow), `internal/delivery` (one stage hook), `internal/runcause`,
`internal/worktree` and `internal/speccheck`.

### Reader inventory

Measured on `7a8b7f70` with a search for every caller of
`spec.ReadArchivedSpec`, `spec.ReadArchivedSpecAt`, `readArchivedSpecAt`,
`spec.ArchivedSpecSlugs` and every `cat-file -e <source_revision>^{commit}`.

Archive-root enumerations, which derive a slug from a file name:

| reader | place | today | after |
| --- | --- | --- | --- |
| review scope | `internal/cli/review.go`, archive-root loop in the scope builder | every changed `*.md` directly under an archive root is read as a record; a prose file blocks the review | only files that pass `spec.IsArchiveRecord` |
| review Spec context | `internal/cli/review.go`, changed-Spec loop | same | same predicate |
| archived Spec names | `spec.ArchivedSpecSlugs`, used by the Spec check's archived names | every `*.md` stem is a slug (`README`) | only files that pass `spec.IsArchiveRecord`, plus folders with `_prd.md` |

Every other archive reader takes a slug it already knows (the Delivery
Queue's archive stage, prerequisites, review correction, supersede, the
review convention, merge evidence, Task completion, the history inventory of
Legacy Archive Folders) and stays unchanged.

Readers that open the folder at `source_revision`:

| reader | place | today when the commit is missing | after |
| --- | --- | --- | --- |
| review Spec context | `internal/cli/review.go`, `cat-file -e` before reading PRD and TechSpec | Spec listed under skipped context | recovery first; skipped only when recovery fails |
| Run causes | `internal/runcause/report.go`, archived graph reader | `source_available: false` | recovery first |
| reconcile and merge evidence | `internal/worktree/archive_reader.go` `archivedSourceAtHead` (leftovers, newest QA Report, archived dirty paths) | "source_revision … is unavailable; keep the worktree" | recovery first |
| Spec check Task Context | `internal/speccheck/citations.go` `archivedTaskContextExists` | the path is unresolved | recovery first |

## Implementation Design

### Interfaces

```go
// internal/spec — record identity (task_02)
func IsArchiveRecordName(name string) bool            // ^[0-9]{4}-[^/\\]+\.md$
func IsArchiveRecord(name string, content []byte) bool // name passes and the front matter has schema: roundfix/archive-record/v1

// internal/spec — Source Revision Recovery (task_02)
const SourceRecoveryFetchTimeout = 60 * time.Second
var ErrSourceRevisionUnavailable = errors.New("source_revision is unavailable")
type SourceRecovery struct{ Revision, PullRequest, From string } // From: local | pull-request
func PullRequestFromSubject(subject string) (string, bool)          // trailing "(#N)"
func EnsureSourceRevision(ctx context.Context, repo, remote string,
	record ArchiveRecord, head, recordPath string) (SourceRecovery, error)

// internal/spec — the pull request field (task_03)
func RecordArchivePullRequest(recordPath, number string) (changed bool, err error)

// internal/delivery — the record stage hook (task_04)
type PullRequestRecorder interface {
	RecordPullRequest(ctx context.Context, workDir, specSlug, head, number string) (PullRequestRecord, error)
}
type PullRequestRecord struct{ Head string; Exact bool }
```

`EnsureSourceRevision` runs Git itself with `exec.CommandContext`. It returns
`From: local` at once when `cat-file -e <rev>^{commit}` succeeds. Otherwise
it takes the number from `record.PullRequest`. When that is empty it uses
`git log -1 --diff-filter=A --format=%s <head> -- <recordPath>` and
`PullRequestFromSubject`. A number that does not match `^[1-9][0-9]{0,9}$`
is never used. It then runs
`git fetch --no-tags --no-write-fetch-head <remote> refs/pull/<N>/head` with
`GIT_TERMINAL_PROMPT=0` and no standard input, under a context limited to
`SourceRecoveryFetchTimeout`, and checks the commit again. Any failure wraps
`ErrSourceRevisionUnavailable` with the step that failed. The remote is the
caller's `watch.push_remote` when it holds the Project Config, otherwise
`origin`.

`RecordArchivePullRequest` parses the record, sets `PullRequest`, renders it
with `RenderArchiveRecord` and writes it. It returns `changed: false` when the
same number is already stored, and refuses a different stored number.

### Data Models

No schema change. `pull_request` is the existing optional key that
`RenderArchiveRecord` writes after `promoted` and `ParseArchiveRecord` reads.
Adding it may shorten the outcome paragraph to keep the 2,048-byte target.
That is the only other byte the change may move.

### API Contracts

1. API Contract: `roundfix archive <slug> --pull-request <n>` with an active
   Spec archives it as today and stores `<n>` in the record. With no active
   Spec and a record at `<archive-root>/<slug>.md` that the default branch tip
   does not hold, it stores `<n>` in that record, writes nothing else and
   prints `recorded pull request <n> in <record>`. When the same number is
   already stored it prints `<record> already records pull request <n>`.
   Exit 0 in all three cases. It refuses with exit 2 and no write when: `<n>`
   is not a positive integer; the record is on the default branch tip; the
   default branch is unknown; another number is stored; it is combined with
   `--plan`; or the record-only form is combined with `--promote`,
   `--qa-override`, `--approval` or `--reason`.
2. API Contract: after a successful archive, the Archive Command checks
   `git merge-base --is-ancestor <source_revision> <default-ref>`, where
   `<default-ref>` is `refs/remotes/origin/<name>` when `origin/HEAD` names
   the default branch, otherwise `refs/heads/<name>`. When the revision is not
   an ancestor, it writes one warning line to standard error (Surface
   Transcripts 1 and 2). When the default branch is unknown it writes
   `roundfix: warning: the default branch is unknown, so source_revision <12-hex> was not checked`.
   The exit code stays 0.
3. API Contract: Source Revision Recovery as in Interfaces. A reader whose
   recovery fails keeps exactly its current output, JSON and exit code.
4. API Contract: the Delivery Queue's Checking stage first calls
   `RecordPullRequest` with the candidate head and the item's pull request
   number. The command workflow resolves the record at the candidate head. It
   returns an empty result when the Spec has no record or the record already
   names the number. When the checkout `HEAD` is a child of the candidate
   head that changes only the record, it adopts that commit if it is exact.
   When the record in the working tree is the only change, it commits it
   without rewriting it. Otherwise it runs `roundfix archive <slug>
   --pull-request <n>` through the step executable and proves the result. It
   commits `docs: record pull request #<n> in the <slug> Archive Record`. The
   engine appends the new head to the item's candidate commits, logs
   `roundfix: Delivery Queue item <slug>: recorded pull request <n> in its Archive Record`
   and returns to Publishing, which pushes and finds the same pull request. An
   inexact commit parks `review-stale`. A failed recording returns an error,
   which parks `delivery-error`. Nothing merges with an Archive Record that
   lacks the number.
5. API Contract: the exactness proof for a Pull Request Record Commit accepts
   it only when its diff against the candidate head touches exactly the
   record path. Both sides must parse as Archive Records whose fields are
   equal except `pull_request`, which must equal the number, and the outcome,
   which must equal the old outcome or a shortening that `RenderArchiveRecord`
   produced.

### Surface Transcripts

1. Surface Transcript: archiving on a delivery branch warns that a squash
   merge would drop the revision.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   archived <slug> -> docs/history/specs/<slug>.md; removed <count> file(s) (<bytes> bytes) kept in Git at <revision>
   stderr:
   roundfix: warning: source_revision <revision> is not an ancestor of origin/main; a squash merge keeps it only in the pull request's head ref, so record the pull request with --pull-request once it opens
   exit: 0
   ```

2. Surface Transcript: archiving with the pull request names the head ref
   readers use.

   ```transcript
   $ roundfix archive <slug> --pull-request 7
   stdout:
   archived <slug> -> docs/history/specs/<slug>.md; removed <count> file(s) (<bytes> bytes) kept in Git at <revision>; pull request 7
   stderr:
   roundfix: warning: source_revision <revision> is not an ancestor of origin/main; readers recover it from refs/pull/7/head
   exit: 0
   ```

3. Surface Transcript: recording the pull request in the record the branch
   already added.

   ```transcript
   $ roundfix archive <slug> --pull-request 7
   stdout:
   recorded pull request 7 in docs/history/specs/<slug>.md
   stderr:
   exit: 0
   ```

4. Surface Transcript: a record already on the default branch stays
   byte-identical.

   ```transcript
   $ roundfix archive <slug> --pull-request 7
   stdout:
   stderr:
   Preflight failed

   Reason:
     Archive Record docs/history/specs/<slug>.md is already on origin/main; an archived record stays byte-identical

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix archive --help' for usage.
   exit: 2
   ```

5. Surface Transcript: the help names the new form and option.

   ```transcript
   $ roundfix archive --help
   stdout:
   Usage:
     roundfix archive <slug>
     roundfix archive <slug> --qa-override --approval <source> --reason <text>
     roundfix archive <slug> --plan
     roundfix archive <slug> --pull-request <number>
   ...
     --pull-request Store the delivery pull request number in the Archive Record
   ...
   stderr:
   exit: 0
   ```

In Transcripts 1, 2 and 4, `origin/main` is the default ref of the test
repository, whose `origin/HEAD` names `main`. `<revision>` is the 12-hex
prefix of `source_revision`. The help keeps every other existing line.

## Coverage Map

- Goal 1 → Pull Request Record Commit (task_04), `RecordArchivePullRequest` and API Contract 1 (task_03)
- Goal 2 → `EnsureSourceRevision` and the four readers (task_02)
- Goal 3 → the archive warning, API Contract 2 (task_03)
- Goal 4 → `IsArchiveRecord` in the review scope, the review Spec context and `ArchivedSpecSlugs` (task_02)
- Goal 5 → skills, guides and glossary (task_01)
- Story 1 → task_02, task_04
- Story 2 → task_03
- Story 3 → task_04
- Story 4 → task_02
- Story 5 → task_02
- Core Feature 1 → task_03
- Core Feature 2 → task_04
- Core Feature 3 → task_02
- Core Feature 4 → task_03
- Core Feature 5 → task_02
- Core Feature 6 → task_01
- Success Metric 1 → task_04
- Success Metric 2 → task_02
- Success Metric 3 → task_03
- Success Metric 4 → task_02

## Integration Points

- Git and the forge: `git fetch` of `refs/pull/<N>/head` from the configured
  remote, the only new network boundary. Tests replace the forge with a bare
  repository on disk whose `refs/pull/7/head` is set with `git update-ref`.
- The Delivery Queue's existing pull request boundary is unchanged. The
  record stage reuses Publishing's push and find-or-create.

## Testing Approach

Every new top-level test in a Parallel Test Package calls `t.Parallel()`
first or names the process-wide state it changes (ADR-0259). Every test uses
temporary repositories and homes and opens no network connection.

- `internal/spec`: the predicate over names and contents; recovery against a
  bare remote with `refs/pull/7/head`. A clone made with `--no-local` from a
  repository that never held the pull request commit lacks it. Recovery is
  shown with the number in the record, with the number in the adding commit's
  subject, writing no ref (`git for-each-ref` before and after, no
  `FETCH_HEAD`), with an unreachable remote, and with a malformed number that
  fetches nothing. `RecordArchivePullRequest` round-trips every other field.
- `internal/cli`: a review whose diff changes a prose `README.md` reaches the
  fake provider, while a file declaring the schema but malformed still blocks.
  The review context and `runs causes` recover a squash-dropped revision. The
  Archive Command transcripts are each asserted whole in one named test. The
  delivery workflow's `RecordPullRequest` and proof run against real Git, and
  an end-to-end Delivery Queue run against the fake forge ends with the
  number in the merged record.
- `internal/worktree`: leftovers and the newest QA Report read through a
  recovered revision.
- `internal/speccheck`: an archived Task Context path resolves through a
  recovered revision.
- `internal/delivery`: the Checking stage with a fake recorder covers record
  and return to Publishing, adoption on retry, inexact parks `review-stale`,
  and error parks `delivery-error`.

The existing tests `TestRunCausesReportAnArchiveRecordWithoutItsRevision`
and `TestLeftoversKeepTheWorktreeWithoutTheSourceRevision` characterize
today's unavailable outcome. They must keep passing unchanged, because their
records carry no number and their fixtures have no reachable forge.

## Build Order

1. Guidance: skills, user guides, glossary terms and the Behavior Surface
   Record (task_01).
2. Record identity and Source Revision Recovery in every reader (task_02)
   (depends on: 1).
3. The Archive Command's `--pull-request` and warning (task_03) (depends on: 1).
4. The Delivery Queue's Pull Request Record Commit (task_04) (depends on: 3).
5. QA gate (task_05) (depends on: 2, 4).

No step depends on another active Spec's artifacts. Spec 0255 also changes
`internal/cli/archive.go` and the archive skill reference, so whichever lands
second rebases. The task_04 steps reuse the task_03 flag through the step
executable.

## Skill and document texts

task_01 writes these texts. Each phrase below must appear verbatim on
whitespace-normalized text.

- `CONTEXT.md`, revised **Archive Record**: keep the entry and its `_Avoid_`
  line, and replace its first two sentences with "The small `<slug>.md`
  record an archive leaves under the resolved archive root: a file named
  `NNNN-<slug>.md` whose front matter declares `roundfix/archive-record/v1`;
  any other Markdown file in that root is not a record. It carries the
  disposition, QA Report and verdict, promoted references, `source_revision`
  and, once the delivery pull request is open, `pull_request`; the removed
  Spec folder stays recoverable in Git at that revision, after a squash merge
  through Source Revision Recovery." Then append "(ADR-0261)" to the
  closing reference list.
- `CONTEXT.md`, revised **Archive Command**: after its first sentence add
  "It warns when `source_revision` is not an ancestor of the default branch,
  and `--pull-request <n>` stores the delivery pull request in the record,
  also in a record the current branch already added (ADR-0261)."
- `CONTEXT.md`, new **Source Revision Recovery**, placed after **Archive
  Command**: "A reader's fetch of `refs/pull/<N>/head` when an Archive
  Record's `source_revision` is missing locally, with the number taken from
  the record's `pull_request` or from the trailing `(#N)` of the commit that
  added the record. It writes no ref, never prompts, and leaves the reader's
  unavailable outcome when it fails (ADR-0261)." `_Avoid_`: Backfill, archive
  ref, history restore.
- `CONTEXT.md`, new **Pull Request Record Commit**, placed after **Source
  Revision Recovery**: "The Delivery Queue's commit, made after it opens a
  Spec's pull request, that stores the number in the Archive Record its
  archive commit added and changes nothing else; the queue accepts it by
  proof instead of a second repository gate (ADR-0261)." `_Avoid_`: Archive
  commit, record backfill, amend.
- The archive reference of the Roundfix skill and the archive user guide
  carry "`--pull-request <n>` stores the delivery pull request in the Archive
  Record" and "warns when `source_revision` is not an ancestor of the default
  branch".
- The review reference and the review user guide carry "Only a file named
  `NNNN-<slug>.md` that declares the Archive Record schema is read as a
  record" and "fetches `refs/pull/<N>/head` when the commit is absent".
- The runs reference and the runs user guide carry "fetches
  `refs/pull/<N>/head` when the commit is absent".
- The deliver reference and the deliver user guide carry "records the pull
  request number in the Archive Record in a Pull Request Record Commit".
- The spec-delivery reference carries "after the pull request opens, run
  `roundfix archive <slug> --pull-request <n>`, commit the record and push".
- The archive-spec skill carries, under a new heading `### Squash merges` and
  never inside `### QA settlement`, "A squash merge keeps `source_revision`
  only in the pull request's head ref; the record names that pull request in
  `pull_request`".

## Risks & Considerations

- Network from a reader. Mitigated by the trigger condition, the time bound
  (cancel, not abandon), `GIT_TERMINAL_PROMPT=0`, no standard input and no
  ref writes. A Run sandbox without network fails the fetch quickly, and the
  reader keeps today's outcome.
- Fetched objects without a ref may be pruned by a later `git gc`. Recovery
  then simply runs again.
- The subject fallback trusts the delivery commit's `(#N)`. A wrong number
  fetches a head that lacks the commit, and the second check reports it
  unavailable. Nothing is written.
- An extra CI run per delivery. Accepted (ADR-0261).
- A record commit on top of the archive commit: Delivery Retry resumes at
  Checking for an archived item with a pull request (ADR-0229, ADR-0237).
  Placing the record step at the start of Checking covers it. The
  archive-stage reconciliation is not re-entered for such an item.
- Spec 0255 changes the same command, so the operator orders the queue.
- Follow-up outside this Spec: the cancellation Backlog Entry gets its own
  Spec, which will add the `cancelled` disposition.

## Research record

Secondbrain, consulted on 2026-10-12 through `wiki/index.md` and
`qmd query "archive record source_revision squash merge pull request ref" --all --files --min-score 0.3`:

- `inbox/roundfix/_triaged/2026-10-09-archive-record-aponta-para-commit-de-branch-sob-squash.md`
  and `inbox/roundfix/_triaged/2026-10-09-review-le-o-readme-da-raiz-de-arquivo-como-archive-record.md`
  are the adopter observations behind both adopted entries. The first records
  that the adopter rejected archiving after the merge, because it inverts the
  Baseline order, and switched to merge commits. That shaped the rejected
  options of ADR-0261.
- `projects/pantheon/mirror/docs/agents/specific-repository.md`, section
  "Spec archive", is the adopter's merge-commit rule. It confirmed that the
  workaround is per repository, which is why this Spec does not prescribe it.
- `projects/fiscus/mirror/docs/adr/0081-a-entrega-usa-squash-com-evidencia-do-pr.md`
  shows another adopter that keeps squash merges and audits through pull
  request evidence. That supports recovery through the pull request rather
  than a merge-policy change.
- The pending `inbox/roundfix/` namespace was empty.

Exa, consulted the same day:

- GitHub Docs, "Checking out pull requests locally"
  (<https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/checking-out-pull-requests-locally>):
  pull request commits stay fetchable as `pull/ID/head`, and `refs/pull/` is
  read-only to pushes. This confirmed the recovery source and ruled out
  pushing an archive ref there.
- Git documentation, `git fetch` (<https://git-scm.com/docs/git-fetch>):
  `--no-write-fetch-head` and refspecs without a destination. This
  confirmed that the fetch can write no ref. The authoring measurement then
  proved it on this repository.
- Stack Overflow, "What is the difference between these two git fetch
  commands?" (<https://stackoverflow.com/questions/42467469/what-is-the-difference-between-these-two-git-fetch-commands>),
  is a secondary source. It describes how a squash merge's commit is not
  connected to the pull request head, which agrees with the measurement.

## Glossary

- adds: **Source Revision Recovery**
- adds: **Pull Request Record Commit**
- changes: **Archive Record**
- changes: **Archive Command**

## Decisions

- One recovery function owned by `internal/spec` that runs Git itself, so the
  non-interactive environment and time bound cannot be lost in a caller's
  runner; see ADR-0261.
- The record stage runs at the start of Checking, which covers a fresh
  publication, a merge-policy recheck and a Delivery Retry that resumes at
  Checking; see ADR-0261.
- The Pull Request Record Commit returns to Publishing without a second
  repository gate, because the proof shows a one-field record change; see
  ADR-0261.
- `RecordArchivePullRequest` lives beside the record code, not in the
  governed archive cut, so the cut is unchanged.
- Docs come first in the graph so every code Task depends on the Task that
  describes its surface, and the owned skills' versions rise once.
