---
spec: 0315-a-task-that-deletes-a-tracked-file-settles
status: active
created: 2026-10-24
surfaces: [backend, cli]
---

# A Task that deletes a tracked file settles

Minimal PRD for a bug fix, minted by `write-techspec`. A Task that declares
`deletes:` removes the file from the working tree only, because the Daemon
makes the commit, so the file is still in the Git index when the settlement
checks run. The Spec Consistency settlement check runs the strict Spec check,
whose retired-history pass lists Citation Sources from the index and refuses
the whole check when a listed file is missing on disk. Every gated Task with a
`deletes:` path therefore fails with `inspect citing file <path>: lstat ...:
no such file or directory` in a repository that has a Removed History Path,
which is every adopter after the history conversion. After this Spec the
readers that list tracked files and then read them from the working tree skip
a tracked path the working tree no longer holds, so such a Task settles. The
deletion still answers to the deletes check, the scope record and the
authorization check, and a path behind a file instead of a directory is still
refused. The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md`](references/2026-10-22-task-com-deletes-nao-assenta-pela-consistencia.md),
   origin oraculum, type hint `fix`, 2026-10-22. With roundfix 0.108.0, Task
   08 of oraculum's Spec 0092 (`deletes:
   packages/backend/test/caixa-sem-venda-e-fluxus-meio-configurado.test.ts`)
   failed settlement twice with `Settlement check failed: settlement check:
   spec consistency: read refusing Spec Consistency findings: run "spec check
   --strict" for Spec 0092-as-vendas-do-erp-pelo-caixa: inspect citing file
   packages/backend/test/caixa-sem-venda-e-fluxus-meio-configurado.test.ts:
   lstat …: no such file or directory`. Oraculum worked around it by asking
   the Agent to run `git rm`. The entry asks to confirm the defect and fix it,
   for example by skipping tracked paths missing on disk or by reading the
   working tree's inventory.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. This Spec is authored against `64c56d7f` (v0.117.0 plus
Spec 0311) and delivers after Specs 0312, 0313 and 0314 and before Spec 0316,
by the maintainer's answer quoted in `_authorization.md`. No active Spec
declares any file this Spec changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change narrows which files three readers read. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and the existing fake Agent runner. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no accepted ADR decides whether a
  reader of tracked files skips a path the working tree lacks; the glossary's
  Citation Source is already a tracked regular text file, and a missing path
  is not a regular file, so this fix makes the readers agree with it without
  a new decision. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no Task creates, edits, renames, moves
  or deletes protected tooling configuration, scripts, ignore files, plugin
  declarations or version pins, and no Task changes a Governed Path (`paths:
  []` in the Spec-contained authorization record
  `docs/specs/0315-a-task-that-deletes-a-tracked-file-settles/_authorization.md`).
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A gated Task that declares `deletes:` and removes the file from the working
  tree only settles; the Spec Consistency settlement check no longer fails
  reading the deleted file.
- A deletion that breaks a link in a live document is still reported by the
  strict Spec check, as a finding instead of an error.
- A deletion of a Governed Path outside the authorization record is still
  refused by the authorization settlement check.
- `roundfix cancel` and the history conversion of `roundfix baseline update`
  read past a tracked Markdown file the working tree lacks instead of
  failing.
- A tracked path whose parent directory was replaced by a file, a symbolic
  link and a binary file keep today's handling.

## Core Features

1. The Citation Source inventory reads the working tree. `CitationSources`
   skips a tracked path that the working tree does not hold, as it already
   does under `docs/history/`.
2. The two Markdown inventories follow. `HistoryCitations` (used by
   `roundfix cancel` and the archive plan) and `StaleHistoryCitations` (used
   by the history conversion of `roundfix baseline update`) skip a tracked
   Markdown file that the working tree does not hold.
3. The settlement path is proved end to end. A gated Task Cycle whose Task
   declares `deletes:` and whose fake Agent removes the file from disk only
   completes with the real settlement checker.

## Non-Goals / Out of Scope

- Changing the deletes check (`deletes: <path> still exists`), the scope
  record of undeclared paths, the authorization settlement check or Governed
  Paths.
- Asking the Agent to stage its deletions, or staging them in the Daemon
  before the settlement checks.
- Changing readers that already skip a missing file or read it from Git: the
  superseded-record plan, the archive plan's Live Document reader, the
  Baseline relocation citation scan and the mechanical QA evidence reader.
- Changing any guide, skill or glossary entry.

## Success Metrics

1. Success Metric: through `TaskCycle` with the default settlement checker on
   a temporary repository with a Removed History Path and the recorded
   retirement rule, a gated Task that declares `deletes:
   packages/backend/test/gone.test.ts` and whose fake Agent removes that file
   from disk only completes.
2. Success Metric: on the same fixture, staging the deletion with `git rm`
   instead also completes, as it did before this Spec.
3. Success Metric: on the same fixture, deleting a tracked Markdown file that
   `docs/index.md` links to makes the real settlement checker return a
   `SC-HISTORY-RETAINED` finding naming the link, with no error.
4. Success Metric: the authorization settlement check refuses a prospective
   Task commit that deletes `.agents/skills/escaped/SKILL.md` from disk when
   the authorization record does not list it, with `QA-AUTH-PATHS`.
5. Success Metric: `CitationSources`, `HistoryCitations` and
   `StaleHistoryCitations` on a temporary repository whose tracked
   `docs/gone.md` and `pkg/gone.test.ts` were removed from disk only return
   no error and no entry for those paths, and still return the citations in
   the files that remain.
6. Success Metric: the three readers still fail, naming the path, when a
   tracked `notes/a.md` sits behind a file `notes` that replaced its
   directory.
7. Success Metric: after release, a Task of an oraculum Spec that declares
   `deletes:` settles without asking the Agent to run `git rm`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-24.

- The oraculum report in the adopted Inbox Entry, written before this Spec
  existed: the refusal `inspect citing file <path>: lstat …: no such file or
  directory` inside `settlement check: spec consistency`, which Success Metric
  1's fixture reproduces before the fix.
- Git's own behavior, Git 2.56.0: `git ls-files -h` describes `-c, --cached`
  as "show cached files in the output (default)" and `-d, --deleted` as "show
  deleted files in the output"; in a temporary repository a committed file
  removed from disk only is listed by both, and after `git rm` by neither.

## Unreachable Acceptance

- criterion: Success Metric 7, after release, a Task of an oraculum Spec that
  declares `deletes:` settles without asking the Agent to run `git rm`.
  reason: the Task runs in oraculum's checkout and Roundfix Home, which no
  Verification or QA row may read or write, and this Spec's authorization
  reads no other repository.
  satisfied-by: oraculum's next Run of a Task that declares `deletes:` with
  the released binary, reported back through the Secondbrain inbox.

## Validation

The decisions were checked on 2026-10-24 against these sources.

- Characterization at `64c56d7f`, in a scratch clone outside this checkout:
  - `bin/roundfix spec check 0312-the-ci-job-runs-the-documentation-contracts-once --strict`
    with `HOME` at a temporary directory exits 0; after `rm
    internal/agent/acp_stream.go` (tracked, not staged) it fails with
    `roundfix: spec check failed: inspect citing file
    internal/agent/acp_stream.go: lstat …: no such file or directory`.
  - A prototype `TaskCycle` test in `internal/daemon`, built as Success
    Metric 1 describes, failed with exactly oraculum's shape: `Settlement
    check failed: settlement check: spec consistency: read refusing Spec
    Consistency findings: run "spec check --strict" for Spec
    0001-sample-feature: inspect citing file
    packages/backend/test/gone.test.ts: lstat …: no such file or directory`.
    The same test with the readers fixed completed the Task.
  - Without a Removed History Path the retired-history pass does not run and
    the same Task settles today, which is why the defect shows only in
    repositories that retired history.
  - With the three readers fixed, `bin/roundfix spec check --strict` after
    deleting a tracked guide reported three `SC-HISTORY-RETAINED` findings
    for the links to it and no error.
  - The prototype passed `go test ./internal/spec ./internal/speccheck
    ./internal/baseline ./internal/daemon ./skills`.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "roundfix
  deletes settlement spec consistency citing file" --all --files --min-score
  0.3`, which returned the Roundfix mirror's settlement and spec guides and
  glossary; none records a decision about deleted tracked files. No Inbox
  Entry was written.
- Exa: none needed; Git's own help text and a temporary repository settle
  what the index lists, offline.

## Glossary

None.

## Decisions

- 2026-10-24 (about 11:15 BRT): the maintainer chose "Spec 0315 após a 0314
  (Recomendado)": a small fix, 1–2 Tasks + QA, delivered after the test-suite
  audit and before Spec 0307, and oraculum is told when it ships.

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 5 briefing for Spec 0315 of 2026-10-24 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- Skip a tracked path whose working-tree entry does not exist, rather than
  list the working tree instead of the index: the index keeps excluding
  untracked files unless asked, and only the missing entry changes.
- Only a missing entry is skipped; a path behind a file or any other read
  error keeps failing.
- Fix the two other readers of the same shape in the same Spec; leave the
  readers that already skip or read from Git.
- No guide, skill, glossary or ADR change: the glossary's Citation Source is
  already a tracked regular text file.

## Open Questions

None.
