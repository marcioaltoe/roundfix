---
spec: 0183-run-storage-that-says-what-it-holds
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# Run storage that says what it holds

Roundfix keeps every Run's row, Run Event Journal rows and artifact directory
in the machine-wide Run Database and the Artifact Root. Four reports about that
storage are wrong or missing, and each misleads the operator who reads them:

- `roundfix gc` on 2026-09-29 reported `Runs pruned: 13` with
  `Journal rows removed: 0`, `Artifact bytes reclaimed: 0` and
  `Orphan artifact dirs removed: 0`. A `gc --dry-run` right after listed the
  same 13 Runs as eligible. Retention never deletes a `runs` row, so every Run
  it already emptied is counted again forever. The best-effort prune that
  `implement`, `resolve` and `watch` run at start has the same defect: 133 Run
  console logs under `~/.roundfix/artifacts` carry
  `pruned Run storage runs=12 journal_rows=0 artifact_bytes=0`. Each of those
  prunes also took the machine-wide write lock to delete nothing.
- Nothing tells the operator that storage is reclaimable. Only the GC Command
  reports it, and only when someone runs it. A 2026-09-17 capture recorded
  4.4 GB of Run storage with 619 MB reclaimable after one day.
- In a bare-repository layout, or a linked worktree of a `--separate-git-dir`
  repository, Spec 0162 moved the default Artifact Root from the checkout path
  to the common Git directory. `gc sanitize` still compares a recorded root
  only with the new default. It classifies a root created before the upgrade
  `overridden` and preserves it forever. The adopted entry calls the command
  `gc --sanitize`; it is the `sanitize` subcommand.
- On 2026-09-29 `roundfix runs list` in this repository printed
  `No Runs found.` and 21 stderr warnings, one per recorded checkout that no
  longer exists, such as
  `inspect retained terminal Runs in repository "/Users/marcio/dev/roundfix-wt-0179": inspect terminal Run: stat recorded Git root ...: no such file or directory`.
  The warnings repeat on every call, and no command clears them. Spec 0175 made
  these Runs visible from the main checkout through their repository key, but
  the retained-Run count still inspects each Run's recorded checkout.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; Runs keep their Run
  IDs, repository keys and recorded Git roots. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0171 (this Spec) makes a retention
  prune report only what it reclaimed, and ADR-0172 (this Spec) gives Doctor a
  read-only `storage` check. Both are recorded before implementation. ADR-0033
  prunes the journal and artifact directory of a terminal Run past the cutoff
  and never touches `runs` rows or Active Run locks; this Spec keeps that
  eligibility and changes only the count and the lock use. ADR-0023 runs
  Implement Runs in per-Run worktrees, whose recorded paths the retained-Run
  count still reads. ADR-0053 keeps reconciliation proof-based; the
  retained-Run count is a hint, not a reconciliation, and this Spec changes no
  reconciliation proof. ADR-0032 keeps Doctor's codex hygiene check and
  ADR-0107 keeps profile readiness over every configured category; the new
  check changes neither. ADR-0090 forbids reusing a repository fact across a
  mutation, and the prune rereads eligibility inside its own call. This Spec's
  gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155
  and ADR-0156. ADR-0093 checks Spec consistency by citation, and ADR-0094
  makes that check artifact-presence-aware. ADR-0097 cites ADR-0080 but carries
  a QA row forward, not storage. ADR-0141 cites ADR-0023 but runs a review Run
  in a clean tracked user checkout. ADR-0161 cites ADR-0053 but releases a
  merged Spec's Runs on the merged head. This Spec changes none of them, so they
  do not apply. ADR-0166 (Spec 0181) has the
  Daemon record the paths a Task changed without declaring them; every Task
  here declares the paths it edits, so it changes nothing here. ADR-0167
  (Spec 0181) keeps the pre-PR Pull Request row from deciding a qualifying
  partial; this Spec's gate aims at `pass`, so it does not apply. ADR-0168
  (Spec 0181) opens a related-ADR gap only for ADRs that predate the Spec,
  which narrows this Spec's own check and holds. ADR-0169 and ADR-0170
  (Spec 0182) govern the pre-PR review base and Task Carry-Forward, which this
  Spec does not touch. ADR-0173 (Spec 0184) reports the citations a Baseline
  history relocation breaks, which this Spec does not touch. All the others
  hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved the Onda 5 plan in
  chat on 2026-09-29 ("Duas filas": Specs 0183 and 0184 in a second
  `roundfix deliver start` after Onda 4); the Roundfix skill files ride the
  standing grant of 2026-09-18 for keeping the shipped skills true to the CLI,
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A retention prune, a GC report and a GC dry run count a Run only when it
  still holds journal rows or an artifact directory to reclaim.
- The operational sweep neither prints a prune line nor takes the write lock
  when there is nothing to prune.
- `roundfix doctor` tells the operator when Run storage is reclaimable, without
  writing anything and without failing.
- `gc sanitize` reclaims a default Artifact Root created before repository keys
  existed, in every layout.
- `runs list` stays silent about Runs whose recorded checkout no longer exists.

## Core Features

1. **A prune reports only what it reclaimed.** A terminal Run past the cutoff
   counts as reclaimable only while it holds Run Event Journal rows or an
   artifact directory under the Artifact Root. `gc`, `gc --dry-run` and the
   operational sweep report only those Runs. The prune takes the write lock
   only when a candidate still has events. A second `gc` after a complete one
   reports `Runs pruned: 0`, and a dry run then reports `Runs eligible: 0`.
2. **Doctor reports reclaimable storage.** `roundfix doctor` prints one
   `storage:` line. It reports `found` with the reclaimable Run count, the Run
   Database free bytes and the next command (`roundfix gc`,
   `roundfix gc compact` or both) when a Run past the cutoff still holds
   something, or when free bytes reach 64 MiB. Otherwise it reports `ok`; a
   missing database reports `ok (no Run Database)`, and a database it cannot
   read reports `partial`. It never reports `failed` and never changes
   Doctor's exit code.
3. **`gc sanitize` recognizes a pre-key default root.** A recorded Artifact
   Root that equals the default derived from the Run's recorded checkout, as
   it was before repository keys, is a default root like one derived from the
   repository key. Every existing safety guard still applies.
4. **`runs list` ignores a vanished checkout.** When a terminal Run's recorded
   Git root no longer exists, the retained-Run count lists the Run's branches
   through its repository key when that repository still exists, and otherwise
   counts only a recorded Run Worktree that still exists. Neither case is a
   warning. A recorded root that exists but fails validation is still a
   warning.

## Non-Goals / Out of Scope

- Deleting `runs` rows, changing Journal Retention eligibility or its default,
  or pruning anything automatically beyond today's operational sweep.
- A storage line in `runs list`, a byte walk in Doctor, or any Doctor mutation.
- Changing `roundfix storage`, `gc compact` or the `gc sanitize` safety guards.
- Making `reconcile` inspect a Run whose recorded checkout is gone; it keeps
  reporting such a Run `unknown` with its retry action.
- A Run Database schema change or a change to the Doctor or GC help text.

## Success Metrics

1. In a disposable Roundfix Home, `gc` over a terminal Run past the cutoff with
   events and an artifact directory reports `Runs pruned: 1`; a second `gc`
   reports `Runs pruned: 0` and a dry run then reports `Runs eligible: 0`. A
   Run whose events are gone but whose artifact directory remains is still
   reclaimed and reported.
2. The operational sweep over a store holding only emptied Runs prints no
   `pruned Run storage` line and takes no write transaction.
3. Doctor over a disposable home with a reclaimable Run prints
   `storage: found` naming `roundfix gc`, and exits `0` when every other check
   passes. Over a home with nothing to reclaim it prints `storage: ok`. The Run
   Database bytes are unchanged afterwards.
4. In a bare clone with `git worktree add`, a Run recorded with a root derived
   from its checkout path is classified `orphaned`, not `overridden`, and
   `gc sanitize --apply` removes its retention-eligible directory. A root equal
   to neither default is still `overridden`.
5. `runs list` over Runs whose recorded checkout was deleted prints no warning
   and counts a Run Branch that still exists in the repository key. A recorded
   root that is a symlink still warns.

## Recorded limits

- A Run whose repository was deleted entirely has no branches left to count;
  only a recorded Run Worktree that still exists is counted as retained.
- Doctor inspects artifact directories only under the current repository's
  Artifact Root. Outside Git it counts journal rows only and says so.
- The free-space threshold is a fixed 64 MiB; it is not configurable.

## Decisions

- **Reclaimable means something is left.** The count follows what the prune can
  still remove, not the eligibility window alone. See ADR-0171.
- **Doctor, not `runs list`.** Doctor is diagnostic and rarely run; `runs list`
  is the Agents' hot path and already prints at most one note. See ADR-0172.
- **One predicate for Doctor and the dry run.** Both count the same Runs, so
  they never disagree.
- **Accept the recorded checkout's default without re-resolving.** Re-resolving
  the checkout maps it back to the common Git directory, which is exactly the
  comparison that fails.
- **Fall back to the repository key, not to silence.** A Run Branch of a
  removed linked checkout still lives in the shared repository, so the count
  looks there before concluding nothing is retained.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight. A prune that still counts
an emptied Run, a sweep that still prints or locks, a Run whose leftover
directory is no longer removed, a Doctor line that fails or writes, a
sanitize that accepts any root, and a `runs list` that hides a symlinked root
would each pass a happy-path test.

The outside-evidence rows rest on sources this Spec did not produce:

- The live Run console logs under `~/.roundfix/artifacts`, written by earlier
  Runs. A read-only count of the files carrying
  `pruned Run storage runs=<n> journal_rows=0 artifact_bytes=0` records the
  false report this Spec removes (133 files on 2026-09-29).
- SQLite's own documentation. `PRAGMA freelist_count` "Return[s] the number of
  unused pages in the database file" (<https://sqlite.org/pragma.html>), while
  the DBSTAT virtual table returns one row per btree page
  (<https://sqlite.org/dbstat.html>). That is why Doctor reads the pragma and
  never `dbstat`.
- Spec 0162's recorded limit, written by another session on 2026-09-24, names
  the bare-layout reproduction and the carried fix this Spec implements.

## Research basis

The four adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). The Secondbrain was consulted
through `wiki/index.md` and
`qmd query "roundfix run storage reclaimable gc doctor retention notice"` and
`qmd query "runs list warning checkout removed worktree recorded git root"`.
It holds nothing beyond this repository's own mirrors: Spec 0014's retention
Tasks, Spec 0059's sanitation, Spec 0162's recorded limit and Spec 0175's
legacy repository key reference. Exa located the two SQLite primary sources
cited above, which settle the cheap-read design of ADR-0172.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
