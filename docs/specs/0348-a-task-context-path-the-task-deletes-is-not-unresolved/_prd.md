---
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
status: active
created: 2026-11-04
surfaces: [backend, docs]
type: fix
---

# A Task Context path the Task deletes is not unresolved

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by fiscus on
2026-11-04: a Task that deleted the file its own `## Context` listed under
`interface:` failed its settlement, because the Spec Consistency settlement
check reported `SC-REF-UNRESOLVED` for the path the Task had been asked to
remove; the finding's fix, "Create <path> or update the Task Context entry",
then led the repair turn to swap in another file, which failed on
`SC-TOOLING-UNDECLARED`, and the Run ended Unresolved with the Task's own
Verification green. After this Spec a Task Context path that is missing from
the checked tree but still tracked at `HEAD` or at the Delivery Base is a
deletion made on the branch, not an unresolved reference; a path no such
revision tracks still reports the finding, whose fix now suggests removing
the entry or declaring the deletion; and `write-tasks` tells authors to list
a file the Task removes only under `deletes:`.

## Sources

One adopter report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md`](references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md),
   `type: fix`, Fiscus, ai-memory message `01a1a5c5`.

At the authoring base `6a0ed904` (v0.130.1) `docs/backlog/` holds one open
entry, about the waiting declaration of the backlog and findings contracts;
it shares no context with this report and stays apart. `docs/findings/`
holds no unresolved entry about Task Context paths.

## Prerequisites

None in `requires`. Authored against `6a0ed904` (v0.130.1), where no other
Spec is active. Specs 0345, 0346 and 0347 are authored beside this one and
deliver in the order 0345 → 0348 → 0347 → 0346. Files this Spec may share
with them, each edited narrowly so a later Spec rebases cleanly:

- the Roundfix skill version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md`, and `skills/testdata/owned-skill-versions.json`
  (0345 edits the roundfix `references/implement.md`, 0347 its
  `references/runs.md`; this Spec edits only `references/spec.md`, so the
  version lines and the record are re-recorded with their command after a
  rebase, never merged by hand);
- `docs/references/behavior-surfaces.json`, re-recorded with its command;
- `CONTEXT.md`, where this Spec changes only the **Delivery Base** entry.

This Spec owns `internal/speccheck/citations.go` and the `write-tasks` skill
for the cycle; it edits no file another cycle AX Spec owns.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads Git revisions of existing repository paths.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads the local repository with `git`. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0295: "One Default Branch Head
  serves every reachability and merge-base read"; this Spec's Task Context
  base reads that Default Branch Head and its merge base, the Delivery Base,
  through the same resolver and never fetches. No new ADR: the acceptance
  rule is a check fact stated in the guides and `CONTEXT.md`. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record, whose
  Rules answer keeps the cycle AW rules and the Governed Paths each Spec
  declares; the governed files are the canonical `write-tasks` skill and its
  Task template, the canonical Roundfix `spec` reference, and the `SKILL.md`
  pairs of the two skills whose version fields rise. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0348-a-task-context-path-the-task-deletes-is-not-unresolved/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A Task Context entry labeled `interface:` or `instruction:` whose path is
  missing from the checked tree reports no `SC-REF-UNRESOLVED` when `HEAD`
  or the Delivery Base tracks the path: the Task's own uncommitted removal
  at settlement, and a removal committed on the branch at the QA gate, in
  the Delivery Queue and in `roundfix spec check`.
- A path that neither revision tracks still reports exactly one
  `SC-REF-UNRESOLVED` with today's summary, so a mistyped path, a path the
  default branch removed before the branch began, and a check outside a Git
  work tree keep failing.
- The finding's fix suggests removing the entry or declaring the deletion
  and never starts with "Create".
- `write-tasks`, its Task template, the Roundfix `spec` reference, the spec
  command guide and `CONTEXT.md` state the rule.

## Core Features

1. The Task Context base. `detectTaskContextReferences` in
   `internal/speccheck/citations.go` accepts a missing path that `HEAD` or
   the Delivery Base tracks, when the checked repository root is the top
   level of its Git work tree.
2. The fix text of the Task Context `SC-REF-UNRESOLVED` finding.
3. Settlement proof: the Daemon settles a gated Task whose turn deletes the
   file its Context lists under `interface:`, and still refuses a Context
   path that never existed.
4. The rule in writing: the `write-tasks` skill and Task template, the
   Roundfix `spec` reference, the spec command guide and the **Delivery
   Base** entry of `CONTEXT.md`.

## Non-Goals / Out of Scope

- `SC-TOOLING-UNDECLARED`, the mechanical authorization audit of a changed
  Governed Path, the scope audit and the `deletes:` rule that refuses
  settlement while a declared deletion still exists: all keep their
  behavior; a deleted Governed Path still needs its grant.
- The `creates:` and `deletes:` labels, which the check already skips, and
  the other `SC-REF-UNRESOLVED` findings (reference index rows, adopted
  sources left in place), which keep their text.
- A path created and removed on the same branch, which neither `HEAD` after
  the removal nor the Delivery Base tracks; it still reports the finding,
  whose fix names `deletes:` (see `_techspec.md` → Decisions).
- Reading a Task's settlement commit trailers, fetching a remote, or
  searching the whole history for a path.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`.

## Success Metrics

1. Success Metric: in a temporary Git repository whose default branch
   `main` (with `refs/remotes/origin/main`) tracks `src/gone.txt` and whose
   pending Task lists it, `roundfix spec check` on a work branch reports no
   `SC-REF-UNRESOLVED` when the file is removed from the working tree only,
   staged for removal, removed by a commit on the work branch, or listed
   under `instruction:` and removed by a commit, and when a file created by
   a commit on the work branch is removed from the working tree only; at
   `6a0ed904` each of these reports one.
2. Success Metric: the same check still reports exactly one
   `SC-REF-UNRESOLVED` whose summary is
   `<task> declares unresolved Task Context path <path>` byte for byte for a
   path no revision tracks, a path created and removed by commits on the
   work branch, a path removed on `main` before the work branch began, a
   repository root without Git, a repository root nested below another
   repository's top level whose `HEAD` tracks the path, and a committed
   removal when no Default Branch Head resolves.
3. Success Metric: that finding's fix reads
   `Remove the Task Context entry for <path> from <task>, or declare it under deletes: when the Task removes it; replace it only with a path the Task itself reads or changes.`
4. Success Metric: through the Daemon's Task cycle, a gated Task whose
   Context lists `interface:` for a tracked file that its turn deletes, on
   disk only or staged, completes; at `6a0ed904` it fails with
   `SC-REF-UNRESOLVED` (the Fiscus case). A Task whose turn adds a Context
   entry for a path no revision tracks still fails settlement, and its
   diagnostics carry the fix of Success Metric 3.
5. Success Metric: both copies of the `write-tasks` skill and of its Task
   template say that a file the Task removes is listed only under
   `deletes:`, never `interface:` and never `instruction:`; both copies of
   the Roundfix `spec` reference and the spec command guide state the base
   rule of Goal 1 and the fix; `CONTEXT.md` says the Spec Consistency Check
   reads a Task Context path at the Delivery Base.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report quoted in the adopted reference, written by the Fiscus
  project before this Spec existed: task_03 of its Spec 0061 deleted
  `.github/workflows/cd-frontend.yml`, listed under `interface:`, and failed
  on `SC-REF-UNRESOLVED`, then on `SC-TOOLING-UNDECLARED` after the repair
  turn swapped the line for another workflow.
- The Git manual, read through Exa on 2026-11-04 (git 2.56.0):
  `git cat-file` (https://git-scm.com/docs/git-cat-file), "-e Exit with zero
  status if <object> exists and is a valid object", and gitrevisions
  (https://git-scm.com/docs/gitrevisions), "`HEAD` names the commit on which
  you based the changes in the working tree" and "A suffix : followed by a
  path names the blob or tree at the given path in the tree-ish object named
  by the part before the colon. … The given path will be converted to be
  relative to the working tree's root directory." The last sentence is why
  the base is read only when the checked root is the work tree's top level.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Characterization at `6a0ed904`, through `speccheck.Check` in temporary
  repositories (an overlay test that wrote nothing to the tree): a pending
  Task listing `interface: src/gone.txt` reports no finding while the file
  exists, and exactly one `SC-REF-UNRESOLVED` with the fix
  `Create src/gone.txt or update the Task Context entry in <task>.` when the
  file is removed from the working tree only, staged for removal, removed by
  a commit, listed under `instruction:` and removed, created and removed on
  the branch, removed on `main` earlier, or never existed.
- Characterization at `6a0ed904` through the Daemon's Task cycle, with the
  fixture of `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles`: the
  Task completes with `deletes:` and fails with
  `Settlement check failed: settlement check: spec consistency: SC-REF-UNRESOLVED: … declares unresolved Task Context path packages/backend/test/gone.test.ts`
  when the same entry is labeled `interface:`. The Fiscus failure
  reproduces.
- Reading at `6a0ed904`: `detectTaskContextReferences` runs for every Task,
  completed ones included, so after a deleting Task settles the QA gate
  precondition (`specGatePrecondition`, which compares with no baseline)
  would report the same finding; the settlement check compares with the
  findings recorded before the Task (`specConsistencyBaseline`), keyed by
  code and summary.
- A prototype of the design in a scratch copy passed the table of Success
  Metrics 1 and 2 and the Daemon case of Success Metric 4; the existing
  `./internal/speccheck`, `./internal/daemon` and `./internal/spec` suites
  passed except `TestTaskContextResolvesThroughTheArchiveRecord`, whose
  fixture leaves the archived folder's removal uncommitted
  (`_techspec.md` → Declared breaks). With the skill wording prototyped,
  `./internal/cli` and `-tags docscontract ./internal/docscontract` passed,
  and `./skills` passed once the wording kept the phrase "never
  `instruction:`" that `TestWriteTasksSkillStatesTheDeclaredPathRules`
  requires.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "spec check task context deleted file unresolved reference settlement" --all --files --min-score 0.3`;
  see `_techspec.md` → Decisions for what it returned.
- Exa: the Git manual passages quoted above.

## Glossary

- changes: **Delivery Base** — the Spec Consistency Check also reads a Task Context path at it

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AX brief for Spec 0348 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04), verbatim:

- 0348: "Sim, as três sugestões (Recommended)" — option text: "Spec 0348:
  SC-REF-UNRESOLVED aceita caminho de Context que a própria Task apaga no
  diff dela ou que existe no tree base; o Fix: sugere remover a entrada em
  vez de criar o arquivo; write-tasks orienta não listar em Context um
  arquivo que a Task apaga."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text:
  "Uma minor (a mudança da Baseline e o runs list fora de repo são feat);
  notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e o
  contorno que podem largar."

Taken from those answers and the characterization, with no further
interview:

- "The Task's own diff" is read as the Task's uncommitted change against
  `HEAD`, which is what the settlement check sees; "the base tree" is the
  Delivery Base, which `roundfix spec check`, the QA gate and the Delivery
  Queue all resolve the same way, unlike a Run start that only a Run knows.
- An "authorized deletion" needs no judgment of its own in this check: the
  entry is the Task's own declaration, and whether the Task may remove a
  Governed Path stays with `SC-TOOLING-UNDECLARED` and the mechanical
  authorization audit, unchanged.
- The finding keeps its summary, so recorded settlement baselines and
  refusal reasons keep matching; only the fix changes.
- `deletes:` stays the one authored form for a removal; the guides say so.
- No new ADR. Two implementation Tasks that share no file, and a QA gate.

## Open Questions

None.
