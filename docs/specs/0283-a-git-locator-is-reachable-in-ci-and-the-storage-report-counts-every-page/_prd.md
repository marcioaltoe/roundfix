---
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
status: active
created: 2026-10-19
surfaces: [backend, cli, docs]
---

# A Git Locator is reachable in CI, and the storage report counts every page

Two checks failed on 2026-10-18 only because of the environment they ran in.

- In GitHub Actions, `SC-HISTORY-RETAINED` called a valid Git Locator invalid.
  The CI checkout has no `refs/remotes/origin/HEAD` and no local `main`, so
  Roundfix could not name the default branch and only tags proved
  reachability. A scratch clone with a stale local `main` failed the same way
  for the opposite reason: the local branch was read before the fetched one.
- `roundfix storage report` failed while the Delivery Queue ran. Its
  reconciliation leaves out the pages `dbstat` never reports, pointer-map
  pages and the lock-byte page, and it reads the Run Database as an immutable
  file even while a writer appends to the write-ahead log.

This is a bug fix without a product interview: the route is `write-techspec`
with this minimal PRD, then `write-tasks`.

After this Spec:

- One **Default Branch Head** answers every reachability and merge-base read:
  the fetched remote-tracking branch, the local branch only when no
  remote-tracking branch exists, and `origin/main` or `origin/master` when
  origin/HEAD is unset. The Git Locator, the QA Delivery Base, the Archive
  Command's `source_revision` warning and the pre-PR review's default base read
  it, so they answer the same locally and in CI.
- The storage report measures one **Storage Snapshot**, groups every page and
  reports the main, `-wal` and `-shm` files apart, so it works on a store of
  any size and while a queue writes.

## Sources

Two Backlog Entries, both created on 2026-10-18 and adopted into this Spec's
`references/` (see `references/_index.md`):

1. `2026-10-18-a-git-locator-in-a-fresh-record-fails-the-check-in-ci.md`
   (references/2026-10-18-a-git-locator-in-a-fresh-record-fails-the-check-in-ci.md):
   pull request #580's Verification gate failed `TestCheckActiveCorpusHasNoErrors`
   on Spec 0279's locator `docs/specs/0279-…@b29c5430`, a commit on `main` newer
   than the last tag `v0.89.0`; the run passed once `v0.90.0` existed. A stale
   local `main` in a scratch clone made a valid locator look invalid.
2. `2026-10-18-storage-report-fails-while-the-run-database-wal-is-live.md`
   (references/2026-10-18-storage-report-fails-while-the-run-database-wal-is-live.md):
   `roundfix storage report` failed with `file bytes=1113714688 grouped
   bytes=1112350720 difference=1363968 tolerance=4096` while the queue ran.

The maintainer put both in one Spec on 2026-10-19. The only other open Backlog
Entry, `2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`, is
Spec 0282's source and shares no component with these.

## Prerequisites

None. This Spec is authored against `f382480b` (v0.91.0) and is delivered
after Specs 0281 and 0282 in the serial cycle order. It changes no file those
plans declare except the shared Roundfix skill version record, the Behavior
Surface Record and `CONTEXT.md`, which every Task regenerates from its own
edits.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one detection source value, `origin-branch`, and
  new report lines are added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every new read is a local Git ref, a local
  file stat or a read-only SQLite query, and no command contacts a remote.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the reachability rule and the storage
  snapshot follow recorded decisions. ADR-0295 (this Spec): "Reachability reads
  the fetched default branch, and origin names it when origin/HEAD is unset".
  ADR-0296 (this Spec): "The storage report reconciles every page of one
  Storage Snapshot". ADR-0278: "A retirement deletes, and the History Root
  holds only Archive Records", which ADR-0295 refines. ADR-0276: "A row the
  authorization keeps offline qualifies, and a byte invariant compares with
  the Delivery Base", which ADR-0295 refines. ADR-0255: "Run Retention removes
  terminal Runs whole and compacts incrementally", which ADR-0296 refines.
  ADR-0172: "Doctor reports reclaimable Run storage from cheap reads", which
  ADR-0296 keeps. ADR-0169: "The pre-PR review diffs the
  candidate from its merge base". ADR-0261: "An Archive Record names its pull
  request, and readers fetch the pull request head". ADR-0104: "A Spec accepts on evidence
  it did not author". ADR-0265: "A Spec proves what it claims, and every active
  Spec is held to it". ADR-0259: "A test runs in parallel unless it names why
  it cannot". ADR-0285: "Every test file compiles for Windows, and a
  repository contract proves it". ADR-0189: "An owned skill's version names its content, and its minimum is the bundle". ADR-0280: "A guide Task
  records its own surface fingerprint, a PRD records who decided, and the
  glossary is held to its own rule". ADR-0256: "A release waits for the skills that
  describe a changed surface". ADR-0271: "A test bounds by real time only what
  it asserts". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0236 (a failed check is judged on a merge with the current default branch); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0266 (a signal goes only to a proven process, and worktrees share one build cache); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ approval of 2026-10-19 through AskUserQuestion, option "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
  go.mod, lint, que pedem autorização específica)."), and the 2026-09-30 skills
  grant ("considere autorizado a ajustar todas as skills se necessário").
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter or Baseline module change.

## Goals

- `roundfix spec check` gives the same `SC-HISTORY-RETAINED` answer in a
  GitHub Actions checkout and in a developer clone of the same refs, including
  for a record archived between two releases.
- A stale local default branch never hides a commit the fetched default branch
  holds, in the Git Locator, the QA Delivery Base, the Archive Command's
  warning or the pre-PR review's default base.
- `roundfix storage report` exits 0 on a Run Database of any size, idle or
  with a live writer, and its page groups reconcile exactly.

## User Stories

1. As the maintainer, I want the CI Verification gate to accept a locator the
   default branch reaches, so that a planning pull request does not wait for a
   release tag.
2. As an operator with an old local `main`, I want Roundfix to read the fetched
   `origin/main`, so that a valid locator, Delivery Base or review base is not
   computed from stale history.
3. As an operator watching storage while the Delivery Queue runs, I want the
   storage report to succeed and to say how many bytes sit in the write-ahead
   log, so that I can decide whether to compact.

## Core Features

1. The Default Branch Head.
   - Detection adds one last step: when origin/HEAD is unset and the current
     branch is not `main` or `master`, exactly one fetched `origin/main` or
     `origin/master` names the default (source `origin-branch`).
   - The head is `refs/remotes/origin/<default>` when it exists, else
     `refs/heads/<default>`.
2. The Git Locator reads it, for both reachability and the newest commit that
   holds a removed path, so `SC-HISTORY-RETAINED`, the Archive Relink and the
   history and Baseline update rewrites inherit it.
3. The sweep. The QA Delivery Base, the Archive Command's `source_revision`
   warning and record-only pull request form, and the pre-PR review's default
   base read the Default Branch Head. Every other default-branch caller is
   named in the TechSpec with the reason it needs no change.
4. The Storage Snapshot. The storage report reads one read transaction, as an
   immutable file when the WAL holds no frame and as a WAL reader otherwise,
   and groups table, free, pointer-map and lock-byte pages; the grouped bytes
   equal the snapshot bytes exactly.
5. Documentation. The `archive`, `review` and `storage report` user guides, the
   Roundfix `archive`, `review` and `storage` references and the glossary
   describe each change.

## Non-Goals / Out of Scope

- Running `git remote set-head`, fetching, or contacting any remote from a
  check.
- Reading CI environment variables to find the default branch.
- Changing which branch Preflight Validation treats as the default for a named
  current branch, or the Linux Lane and Delivery Queue reads that already
  fetch the remote; their files belong to Spec 0282 in this cycle.
- Widening the storage report's tolerance, checkpointing, or writing anything
  while measuring.
- Changing the Doctor storage check, which stays a cheap read without
  `dbstat` (ADR-0172).

## Success Metrics

1. Success Metric: in a temporary clone that fetched every branch and tag into
   `refs/remotes/origin/*`, has no origin/HEAD and no local `main`, and is
   detached at a pull request head, a Git Locator whose commit `origin/main`
   holds and no tag contains resolves to that commit, and `SC-HISTORY-RETAINED`
   reports no finding for an Archive Record that cites it.
2. Success Metric: in a temporary clone whose local `main` is behind
   `origin/main`, the same locator resolves to its commit, and the newest
   commit that holds a removed path is read from `origin/main`.
3. Success Metric: in temporary repositories without origin/HEAD, the QA
   Delivery Base resolves against `origin/main`, the Archive Command's warning
   names `origin/main`, and the pre-PR review's default base is the
   `origin/main` tip, not a stale local `main`.
4. Success Metric: `roundfix storage report` exits 0 with `Reconciliation
   difference: 0` on a temporary store with at least two pointer-map pages, and
   on a temporary store whose writer holds WAL frames not yet checkpointed; in
   the second case the snapshot counts the WAL's pages, `WAL bytes` is
   positive, and the database and WAL bytes are unchanged by the report. An
   idle store gains no `-wal` or `-shm` file.
5. Success Metric: the pointer-map and lock-byte page counts give 333 pages
   for 271,903 pages of 4096 bytes, the difference recorded on 2026-10-18.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-19.

- The 2026-10-18 measurement in Source 2, taken by the Spec 0279 authoring
  session before this Spec existed: file bytes 1,113,714,688 and grouped bytes
  1,112,350,720. 1,113,714,688 / 4096 = 271,903 pages; the file format formula
  gives floor((271,903 − 2) / 820) + 1 = 332 pointer-map pages plus the
  lock-byte page 262,145, so 333 pages × 4096 = 1,363,968 bytes, the reported
  difference to the byte. The QA gate recomputes it offline from the adopted
  entry.
- This repository's history, which no Task of this Spec authors. Spec 0279's
  Archive Record cites the locator
  `docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage@b29c54309d36daf812c1cacbd33c114f3cb1d985`,
  whose commit lies between tags `v0.89.0` and `v0.90.0`. The QA gate rebuilds
  the CI checkout's ref layout from a local clone of the audited head, with
  every tag after `v0.89.0` deleted, and records the check's answer.
- SQLite, "The DBSTAT Virtual Table" (https://www.sqlite.org/dbstat.html):
  "The DBSTAT table only reports on the content of btrees within the database
  file. Freelist pages, pointer-map pages, and the lock page are omitted from
  the analysis."
- SQLite, "Database File Format" (https://sqlite.org/fileformat.html), section
  1.8: "In a database with ptrmap pages, the first ptrmap page is page 2. ...
  The second pointer map page will be on page J+3", with J = U/5.
- SQLite, "Uniform Resource Identifiers" (https://sqlite.org/uri.html):
  immutable "skips all file locking and change detection on immutable database
  files. If this query parameter ... asserts that a database file is immutable
  and that file changes anyhow, then SQLite might return incorrect query
  results and/or SQLITE_CORRUPT errors."
- actions/checkout, README (https://github.com/actions/checkout): "Set
  `fetch-depth: 0` to fetch all history for all branches and tags", and "In a
  pull request trigger, `ref` is required as GitHub Actions checks out in
  detached HEAD mode"; issue #910, "`origin/HEAD` Not Present"
  (https://github.com/actions/checkout/issues/910): "I find the absence of
  `origin/HEAD` when passing `fetch-depth: 0` counterintuitive."

## Validation

The decisions were checked on 2026-10-19 against these sources.

- Secondbrain.
  - Read `wiki/index.md`, then ran `qmd query` from the Secondbrain root:
    - "git default branch origin/HEAD detached HEAD CI checkout reachability"
    - "SQLite WAL storage report run database size reconciliation"
  - The first returned the mirror of ADR-0236
    (`projects/roundfix/mirror/docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`),
    which judges a failed check against the current default-branch tip, the
    mirror of Source 1 and the mirror of
    `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`,
    which keeps squash-merged provenance through the pull request ref. None
    holds a decision on reading the local branch before the fetched one, which
    supports making the fetched branch authoritative.
  - The second returned the mirror of Source 2, of the `storage report` guide
    and of ADR-0172, and
    `projects/roundfix/mirror/docs/user-guide/run-database-lifecycle.md`; none
    mentions pointer-map pages or the WAL, so nothing prior conflicts.
  - A third query on pointer-map pages did not finish within five minutes and
    is not counted.
- Exa.
  - Supporting: the SQLite `dbstat`, file format and URI pages and the
    actions/checkout README and issue #910, quoted above; actions/checkout
    issue #2219 reports that Git 2.50 stopped setting origin/HEAD even after a
    full fetch, so origin/HEAD cannot be assumed on any runner.
  - Challenging: issue #910 and #2219 recommend `git remote set-head origin
    --auto`. That contacts the remote and writes a ref; a read-only offline
    check cannot rely on a workflow step, so Roundfix reads the fetched
    branches instead.
  - SQLite's WAL page states that a read-only WAL database can be read when
    the `-shm` and `-wal` files already exist, which is why the WAL reader is
    chosen only when those files exist.
  - No research Inbox Entry was captured; the authoring brief keeps the
    Secondbrain read-only for this session.
- Measurement, 2026-10-19, read-only.
  - At `f382480b`, a throwaway `go test -overlay` probe built the CI-like clone
    and the stale-local clone described in Success Metrics 1 and 2:
    `ResolveGitLocator` returned no commit in both.
  - `roundfix storage report` (v0.91.0) on this machine's idle store failed
    with `file bytes=1126957056 grouped bytes=1125576704 difference=1380352`;
    its `-wal` held 0 bytes. A read-only immutable query gave 275,136 pages,
    0 free pages and 274,799 `dbstat` pages: 337 missing, the formula's 336
    pointer-map pages plus the lock-byte page. No byte of that store was
    copied.
  - A throwaway Python probe on a temporary WAL database with incremental
    auto-vacuum measured 3,013 pages through the immutable reader and 5,020
    through a WAL reader while 2,000 rows sat in the WAL, with both readers
    leaving the database and WAL bytes unchanged.

## Glossary

- adds: **Default Branch Head**
- adds: **Storage Snapshot**
- changes: **Git Locator**
- changes: **Delivery Base**

## Decisions

- 2026-10-19 (about 16:20 BRT): decomposition — "Três Specs (Recomendado)":
  J (0281) → K (0282) → L (0283), one minor each.
- 2026-10-19: authorizations — "Mexer em Governed Paths declarados", option
  text "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica)."; and "Merge squash + release + push".
- 2026-10-19: costs "Regras do ciclo AH/AI (Recomendado)"; failures "Regras
  AE–AH e parar na proposta (Recomendado)"; authoring "Opus + strict check,
  enfileirar sem esperar você (Recomendado)".
- 2026-10-19: scope of Spec L, from the lead's brief of that decision — both
  sources in one Spec; reachability is proved the same way locally and in a
  GitHub Actions checkout and never trusts a stale local branch over the
  fetched remote; the storage report accounts for the WAL or reads a
  consistent snapshot, and widening the tolerance is not the fix.
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the lead's cycle AJ
brief of 2026-10-19 and from the code and measurements characterized in the
TechSpec.

Taken from those sources, with no further interview:

- origin names the default only when exactly one of `origin/main` and
  `origin/master` exists; ambiguity stays undetermined. See ADR-0295.
- The fetched branch answers first and the local branch only without one. See
  ADR-0295.
- The storage report's failure is caused by pointer-map and lock-byte pages
  even on an idle store; the WAL is handled by choosing the reader, and both
  are fixed. See ADR-0296.

## Open Questions

None.
