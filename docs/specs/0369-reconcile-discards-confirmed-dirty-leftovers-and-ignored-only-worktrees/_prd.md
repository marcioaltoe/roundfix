---
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
status: active
created: 2026-11-09
surfaces: [cli, docs]
type: fix
---

# Reconcile discards confirmed dirty leftovers and releases ignored-only worktrees

`roundfix reconcile` has no supported way to clean up two kinds of leftover that a delivered Spec leaves behind. A Run that ended Failed or Unresolved for a Spec another Run delivered and archived stays `dirty` for ever when its worktree holds uncommitted leftovers outside the archived Spec's declared and recorded paths, and a recorded Run Worktree directory that Git no longer registers stays `unknown` for ever even when it holds only files the repository ignores; empty repository directories under the Worktree Root also stay. On 2026-11-09 Fiscus had to clean eleven Run Worktrees and one directory by hand, which the Roundfix skill forbids. After this Spec the dry run lists each discardable Run's uncommitted paths with one digest, `roundfix reconcile --discard-dirty <digest>` records a Branch Disposition with those paths and then removes the Run Worktree and Run Branch, an unregistered directory that holds only ignored files is `safe`, and `--apply` on a full scan removes empty repository directories.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md`](references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md),
   ai-memory message 01a1c01f-7095, reported by fiscus on 2026-11-09, a
   measured report with a suggestion for each gap.

At the authoring base `docs/backlog/` holds one open entry, the disk-full
tests that Spec 0368 adopts, and `docs/findings/` holds no unresolved Finding;
neither shares this context.

## Prerequisites

None in `requires`. Authored against `cfc08683` (v0.137.2 and the 0366/0367
plans). The queue delivers Specs 0366, 0367 and 0368 first. They may also
change `CONTEXT.md` (other entries), `docs/references/behavior-surfaces.json`,
`skills/testdata/owned-skill-versions.json` and the Roundfix skill's
`SKILL.md` version fields. This Spec edits only the **Branch Disposition**,
**Reconcile Command** and **Run Worktree Reconciliation** entries of
`CONTEXT.md` and adds **Dirty Discard** after **Branch Disposition**; it edits
only the `reconcile` reference of the Roundfix skill, and re-records the
version record and the Behavior Surface Record by their commands, so a rebase
reruns the commands instead of merging bytes. It owns
`internal/cli/reconcile*.go`, `internal/daemon/reconcile.go` and the terminal
Run reconciliation in `internal/worktree/worktree.go`; it does not touch
`internal/cli/deliver_disk_full*.go`, `internal/cli/doctor_disk_space.go`,
`internal/spec/**`, `internal/speccheck/**` or `internal/baseline/**`.
`internal/cli/cli.go` and `internal/cli/cli_test.go` change only in the
`reconcile` usage text and the test row that pins it.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; the
  digest is a SHA-256 hex string computed from the report, and the Branch
  Disposition record keeps its `roundfix-branch-disposition/v1` schema with
  two optional fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the command reads the local Run
  Database and local Git only and makes no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0379 (new) records the decision,
  ADR-0379: "A Dirty Discard is explicit, confirmed and recorded." It amends
  ADR-0053, ADR-0053: "Ambiguous, dirty, and unintegrated work stays intact",
  ADR-0289, whose rejected option it answers, ADR-0289: "A stray edit outside
  the Spec's scope is work only its author can judge", and ADR-0290, whose
  Run Worktree row gains the Dirty Discard as a reclaimer, ADR-0290:
  "Reclaiming Run Worktrees by age was rejected."; task_03 writes the three
  notices. ADR-0115 makes a disposal a separate act that writes its evidence
  first, ADR-0115: "The act writes what the branch held", which the Dirty
  Discard follows. ADR-0212, ADR-0161 and ADR-0232 supply the commit proof
  the discard requires, ADR-0212: "When the merged head holds the Spec
  archived", ADR-0161: "A terminal Run of a merged Spec is released when
  every commit its Run Branch holds and the merged head lacks is represented
  there", ADR-0232: "Merge evidence releases the Runs and item branches of a
  merged Spec". Not applicable, and marked `not applicable` in the TechSpec:
  the related-ADR candidates that the listed ADRs draw in, each on a part
  this Spec leaves alone (Task Carry-Forward, Stale Run Owner Records, the
  history cleanup and conversion, the pre-PR review, the Delivery Queue's
  merge records and the authoring checks they cite), ADR-0092, ADR-0120, ADR-0165, ADR-0169, ADR-0170, ADR-0196, ADR-0197, ADR-0227, ADR-0233, ADR-0237, ADR-0247, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0256, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0282, ADR-0284, ADR-0286, ADR-0295, ADR-0299, ADR-0306, ADR-0308, ADR-0313, ADR-0317, ADR-0345, ADR-0346, ADR-0348, ADR-0349, ADR-0352, ADR-0366, ADR-0368 and ADR-0373.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA batch 2 answers of 2026-11-09 quoted in the authorization record
  grant this Spec, its owned skill edit and the governed help test row; no Makefile, workflow, Go module,
  lint or formatter change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- An operator can discard the uncommitted leftovers of an ended Run of an
  archived Spec through one supported, confirmed command that records what it
  discarded, and nothing else becomes discardable.
- A recorded Run Worktree directory that Git does not register and that holds
  only ignored files is released by `--apply`, and one that holds anything
  else stays `unknown`.
- `--apply` on a full scan leaves no empty repository directory under the
  Worktree Root.
- The glossary, the amended ADRs, the user guide and the Roundfix skill's
  `reconcile` reference describe the new switch and both releases.

## User Stories

1. As an operator cleaning up after a delivered Spec, I want the dry run to
   list the leftovers of each ended Run that a discard can remove and to print
   one digest, so that I confirm exactly what I saw.
2. As an operator, I want `--discard-dirty <digest>` to record a Branch
   Disposition with the discarded paths before it removes anything, and to
   refuse when the set changed, so that a discard is never broader than the
   one I confirmed.
3. As an operator, I want a Run Worktree directory that holds only ignored
   runtime files to be released by `--apply`, and empty repository
   directories to go with it, so that a full scan ends with nothing left.

## Core Features

1. Dirty Discard candidates: a terminal Run whose outcome is `Failed`,
   `Unresolved`, `Stopped`, `BudgetExceeded` or `TimedOut`, whose Spec's
   Archive Record is on the default branch head, whose commit proof is `safe`
   or `superseded`, whose only remaining obstacle is an uncommitted path
   outside the archived Spec's declared and recorded paths, and to which no
   live process tree is attributed. The report lists their paths
   (`discardPaths`) and one digest (`discardDirtyDigest`).
2. `--discard-dirty <digest>`, mutually exclusive with the other three
   switches: refuses without mutation when no candidate exists or the digest
   differs; otherwise, per candidate, writes the Branch Disposition record
   with the discarded paths, re-proves the Run head, the proof head and the
   path set, removes the Run Worktree with `--force` and deletes the Run
   Branch.
3. An unregistered recorded Run Worktree whose name is the Run ID, whose Run
   Branch is absent, with no `.git` entry at its top and only files that
   `git check-ignore` reports ignored (or none) is `safe`; `--apply` re-lists
   the files, refuses a changed list and removes the directory.
4. `--apply` on a full scan removes each empty directory directly under the
   configured Worktree Root whose name has the shape `<slug>-<8 hex>`; the
   dry run lists them (`emptyDirectoryCandidates`).
5. Documentation: ADR-0379 (written by this plan) and the amendment notices
   in ADR-0053, ADR-0289 and ADR-0290; the **Dirty Discard** term and the
   **Branch Disposition**, **Reconcile Command** and **Run Worktree
   Reconciliation** entries of `CONTEXT.md`; the `reconcile` usage text and
   user guide; the Roundfix skill's `reconcile` reference.

## User Experience

A dry run over a repository with one candidate keeps every existing line and
adds, under that Run, one `  discard-path: <path>` line per uncommitted path,
and after the `Apply with:` line:

```text
Discard dirty leftovers with: roundfix reconcile --discard-dirty <64 hex>
```

A confirmed discard prints the report in mode `discard-dirty`, the Run's
action `discarded`, and its evidence ending with
`; branch record: <artifact root>/runs/<run-id>/branch-disposition.json`.
A digest that no longer matches prints nothing on stdout, exits 2 and writes:

```text
roundfix: reconcile failed: the Dirty Discard set changed: confirmed digest <digest>, current digest <digest>; run roundfix reconcile and confirm the digest it prints
Run 'roundfix reconcile --help' for usage.
```

An unregistered ignored-only directory reads `classification: safe` with
evidence
`recorded Run Worktree is not registered in the Git root and holds only files the repository ignores`,
and an empty repository directory reads
`Empty repository directory: <path>` with `  action: would remove with --apply`.

## Non-Goals / Out of Scope

- Discarding a Run whose commit proof fails, whose Spec is not archived on
  the default branch, that is Clean or Integration Pending, or that still owns
  a live process tree.
- Any automatic discard: the Delivery Queue and `--apply` never discard
  dirty leftovers.
- Releasing an unregistered directory that holds a `.git` entry, a file the
  repository does not ignore or a tracked path, or whose Run Branch exists.
- Journaling the discard or the ignored-only release in the Run Database;
  the Branch Disposition record is the discard's record.
- Any change to the Makefile, the workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the Baseline, or the `### QA settlement` section
  of any skill.

## Success Metrics

1. Success Metric: in a temporary repository whose archived Spec's Stopped
   Run holds one committed, delivered Task and an untracked file outside the
   Spec's scope, the dry run lists that path and a digest, and
   `--discard-dirty <digest>` writes a Branch Disposition record naming it and
   leaves neither the Run Worktree nor the Run Branch; at the authoring base
   the Run is `dirty` with no supported exit.
2. Success Metric: a wrong digest, a Run with an unrepresented commit, a Run
   of an unarchived Spec, a Clean Run and a Run with a live process tree are
   never discarded, and their worktrees and branches stay.
3. Success Metric: an unregistered Run Worktree directory holding only
   `.storage/postgres/18/docker/base/1` under an ignored `.storage/` is
   `safe` and removed by `--apply`; the same directory with one file the
   repository does not ignore stays `unknown`, as at the authoring base.
4. Success Metric: after `--apply` on a full scan, no empty directory named
   `<slug>-<8 hex>` remains directly under the Worktree Root, and an empty
   directory with another name stays.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report quoted in the adopted reference, written by the Fiscus
  project before this Spec existed: "summary dirty=11 unknown=1" and "O
  diretório só tinha .storage/postgres/18/docker".
- The Git manual for `git check-ignore`
  (https://git-scm.com/docs/git-check-ignore, read through Exa on
  2026-11-09): "By default, tracked files are not shown at all since they are
  not subject to exclude rules" and, under EXIT STATUS, "1 : None of the
  provided paths are ignored." The ignored-only proof relies on both.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-09 against these sources.

- Reading at `cfc08683`: `inspectTerminalRunMerged` returns `unknown` with
  `recorded Run Worktree is not registered in the Git root` as soon as
  `recordedWorktree` reports the path unregistered, before it looks at the
  Run Branch or the directory; a dirty Run keeps no inspection evidence, so
  no apply can act on it; `revalidateTerminalRunApply` refuses every state
  but `safe`, `superseded` and `released`; `ReconcileIntegration` requires a
  Run head and the action `cleanup`, so an ignored-only release, which has no
  Run Branch, cannot be journaled there; `WriteBranchDispositionRecord`
  refuses a disposition that is not superseded.
- A probe with Git 2.56.0 in a temporary repository under this clone, with
  `.gitignore` holding `.storage/`, `*.log` and `!keep.log`: `git check-ignore
  -z --stdin` printed `.storage/postgres/18/docker/base/1`, `.storage/x.sock`
  and `a.log`, and not `keep.log`, the tracked `tracked.txt` (also after
  `tracked.txt` was added to `.gitignore`) or `src/new.go`; with only
  `src/new.go` on stdin it exited 1.
- The repository directory name: `repoSlug` writes
  `sanitizeSlugBase(base) + "-" + 8 hex`, and both directories Fiscus removed,
  `fiscus-2f7556b8` and `fiscus-0035-implementation-20260908-392a816f`, have
  that shape.
- Published sources (Exa): the `git check-ignore` and `git worktree` manual
  pages on git-scm.com; the second confirms that a working tree deleted
  without `git worktree remove` leaves administrative files for
  `git worktree prune`, which is how a recorded Run Worktree becomes
  unregistered.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix reconcile dirty Run Worktree discard leftovers archived Spec" --all --files --min-score 0.3`
  and
  `qmd query "git worktree not registered directory only ignored runtime artifacts removal" --all --files --min-score 0.3`;
  they returned only this repository's own mirrored pages (the `reconcile`
  user guide, `CONTEXT.md`, the archived Spec 0008 and two decision records
  on worktree paths and archive removal), which the decisions above already
  read from the repository itself; no page decides otherwise.
- Pending ai-memory messages: none read by this author; the operator relayed
  the Fiscus message.

## Glossary

- adds: **Dirty Discard** — the Reconcile Command's confirmed removal of an ended Run's dirty leftovers.
- changes: **Branch Disposition** — it also records the paths a Dirty Discard removed.
- changes: **Reconcile Command** — it gains `--discard-dirty <digest>`, and `--apply` removes empty repository directories.
- changes: **Run Worktree Reconciliation** — an unregistered Run Worktree that holds only ignored files is `safe`.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BA batch 2 brief for Spec 0369 of 2026-11-09 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-09 ~11:20 BRT), verbatim:

- Question "Fiscus relatou duas lacunas do `roundfix reconcile`: (1) sobra
  suja de Run terminal (Failed/Unresolved) de Spec já arquivada fica 'dirty'
  sem caminho suportado; (2) diretório de Run Worktree não registrado no Git,
  só com artefatos ignorados de runtime, fica 'unknown' para sempre (e
  diretórios de grupo vazios ficam). O que fazer?" → "Spec E agora
  (Recommended)" — option text: "Spec fix: modo explícito com confirmação
  (ex. --discard-dirty) só para Run terminal cuja Spec está arquivada no
  default branch, gravando Branch Disposition com os paths descartados;
  diretório não registrado só com artefatos ignorados vira liberável; grupos
  vazios removidos. Patch após o merge."
- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto fica
  para amanhã."

Taken from those and the evidence, with no further interview:

- The confirmation is a SHA-256 digest of the candidate set that the dry run
  prints and `--discard-dirty` must repeat, as `roundfix baseline` confirms a
  Plan Digest; a Run ID alone binds no path set. See ADR-0379.
- A discard requires the commit proof to hold, so it removes only
  uncommitted leftovers and never committed work the delivery lacks; paths
  inside the archived Spec's scope are already superseded by ADR-0289.
- The accepted outcomes are the ended-without-delivering ones the brief
  names as Failed, Unresolved and abandoned: `Failed`, `Unresolved`,
  `Stopped`, `BudgetExceeded` and `TimedOut`.
- "Ignored" is proven by `git check-ignore -z --stdin` in the repository,
  without `--no-index`, over every file below the directory; a `.git` entry
  at its top or an existing Run Branch keeps today's `unknown`.
- Empty repository directories are removed only on a full-scan `--apply`,
  after the Run releases, and only for the `<slug>-<8 hex>` name shape.
- Three implementation Tasks and the QA gate, in series: the Dirty Discard,
  the ignored-only release with empty directories, then the ADR notices and
  the skill reference.
- The release is the operator's patch right after the merge; no Task
  releases.

## Open Questions

None.
