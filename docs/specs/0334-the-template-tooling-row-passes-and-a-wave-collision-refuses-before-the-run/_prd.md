---
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
status: active
created: 2026-10-30
surfaces: [backend, cli, docs]
---

# The template's tooling row passes the Spec Check, and a Wave collision is refused before the Run

Minimal PRD for a bug fix, minted by `write-techspec`. Quaestor followed the
`write-prd` and `write-techspec` templates for a Spec that changes no
Governed Path and `roundfix spec check --strict` refused the Tooling
authority row the templates prescribe with `SC-TOOLING-UNBOUNDED`. In another
Spec Quaestor added a corrective Task that shares a file with a Task already
completed; `roundfix implement` created a Run and only then refused the Wave,
although the completed Task can never run again. After this Spec the
template's no-mutation row passes when its record grants no Governed Path, a
completed Task never collides, and a real collision is refused before any Run
exists.

## Sources

Two Backlog Entries, adopted by this Spec's planning change (see
`references/_index.md`). Both were written on the operator's unpushed branch
`docs/backlog-ar-batch-1` (commit `34434d59`) from Quaestor's ai-memory
message `01a17e01-45fd-7a81-9052-54c660114746` (2026-10-27 12:36 BRT) and are
created here with the same body and `status: promoted`:

1. [`references/2026-10-27-the-authoring-templates-and-the-tooling-check-disagree-on-no-tooling-change.md`](references/2026-10-27-the-authoring-templates-and-the-tooling-check-disagree-on-no-tooling-change.md),
   `type: fix`, item 1 of the message.
2. [`references/2026-10-27-a-wave-collision-with-a-completed-task-refuses-after-the-run-starts.md`](references/2026-10-27-a-wave-collision-with-a-completed-task-refuses-after-the-run-starts.md),
   `type: fix`, item 2 of the message.

At the authoring base `docs/backlog/` and `docs/findings/` hold no other
entry on the Spec Check's tooling rows or on Wave collisions.

## Prerequisites

None in `requires`. Authored against `35de7f74` (v0.126.0) with no active
Spec under `docs/specs/`. Specs 0332 and 0333 deliver before this one and
Spec 0335 after it. Spec 0333 edits authoring skills and the Roundfix skill:
if it changes `write-prd`, `write-techspec` or the Roundfix `SKILL.md`, the
later Task here edits those files after that merge, raises the version above
whatever 0333 recorded through the record command, and pins none of their
bytes. Spec 0335 also raises the Roundfix skill version after this one.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the refusal and the template wording are text. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  network access; the commands read a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the new ADR-0347 records both
  decisions. ADR-0131 (the tooling row states applicability, not mutation)
  is the rule the checker returns to; ADR-0179 (an explicitly empty paths
  list grants operations only) defines the record the row cites, and its
  rule that every bounded path stays governed is unchanged. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AS answers of 2026-10-30 quoted in the authorization record, which keep the
  cycle AJ grant for the Governed Paths each Spec declares; the governed files
  are the Spec Check's constraint reader and the owned skills `write-prd`,
  `write-techspec` and `roundfix`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `internal/speccheck/constraints.go`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A PRD or TechSpec whose Tooling authority row is the templates' no-mutation
  row, citing a record with `paths: []`, passes `roundfix spec check
  --strict` with no tooling finding.
- A row that cites a record granting Governed Paths, or that declares a
  mutation, without bounded files is still refused with
  `SC-TOOLING-UNBOUNDED`.
- The templates say how that row ends, and a contract test keeps them and
  the checker in agreement.
- A completed Task never forms a Wave collision in `roundfix spec check`, the
  QA gate's precondition, the settlement checks or the Daemon's
  pre-dispatch refusal; a Task reopened to `pending` counts again.
- `roundfix implement` refuses a real Wave collision before it records a
  Run, when the Task Capacity is above 1.
- The Roundfix skill, the user guides and `CONTEXT.md` say both rules.

## Core Features

1. The no-mutation row with an empty grant. A Tooling authority row whose
   reason says `no protected tooling mutation` declares no mutation when it
   cites no record, as today, and when the one record it cites is an
   operative grant whose `paths` is `[]`. Every other row keeps today's
   reading.
2. The templates name the row's end. Both templates say to cite the record
   and, when its `paths` is `[]`, to end the row with `bounded files: none`.
   A documentation contract test composes the row from each template and
   requires the checker to accept it and to refuse the same row when the
   record lists a Governed Path.
3. A completed Task is outside every Wave. `spec.Collisions` reports no pair
   in which either Task's status is `completed`.
4. Refusal before the Run. The Implement Command's Preflight Validation runs
   the collision rule on the committed graph when the Task Capacity is above
   1 and exits `2`, naming each pair, its shared files and the edge that
   orders it, with nothing written to the Run Database. The Daemon keeps its
   pre-dispatch refusal.
5. The rules in writing. The `implement` and `spec` guides, the Roundfix
   `implement` and `spec` references and `CONTEXT.md` (**Wave Collision**)
   say that a completed Task never collides and that `implement` refuses a
   collision before the Run.

## Non-Goals / Out of Scope

- Changing which files make a Task touch a path (Verification operands,
  declared Context, prior Run settlement commits).
- A `not applicable` Tooling authority row: it keeps today's reading.
- Tightening a bounded-files value such as `none (...)`, which today counts
  as bounded; changing it would refuse rows adopters already write.
- Changing the Delivery Queue: it starts `implement`, which now refuses
  earlier.
- Any text inside a `### QA settlement` section; the Makefile, the CI
  workflows, `go.mod`, `go.sum`; any Baseline module or guide under
  `docs/agents/`.

## Success Metrics

1. Success Metric: a PRD and a TechSpec whose Tooling authority row is
   `applicable — no protected tooling mutation proposed or authorized`,
   citing a record with `paths: []` and ending `bounded files: none`, give no
   `SC-TOOLING-UNBOUNDED`; at `35de7f74` both report it, as the v0.124.0
   binary does for the same row.
2. Success Metric: the same row citing a record whose `paths` lists
   `Makefile` still reports `SC-TOOLING-UNBOUNDED`, and so does the row
   `applicable — protected tooling mutation proposed and authorized.` with
   no bounded files.
3. Success Metric: `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck`
   passes on the changed templates and checker and fails when the checker
   reverts to `35de7f74` or when a template loses its `bounded files: none`
   wording.
4. Success Metric: a corrective Task that shares a file with a completed
   Task and has no edge to it gives no `SC-WAVE-COLLISION` in the Spec
   Check, the QA gate's precondition, the settlement checks or the Daemon's
   pre-dispatch refusal; two Tasks that are not completed and share a file
   still give the same refusal in all four.
5. Success Metric: `roundfix implement` with a Task Capacity of 2 and two
   pending Tasks sharing `shared/notes.txt` exits `2` with `Task Capacity 2
   would run colliding Tasks in one Wave, so no Run was created:`, records no
   Run and starts no Agent session; at `35de7f74` it exits `1` with
   `implement failed after Run start: task cycle: refuse Wave before
   dispatch because Tasks collide`. At a Task Capacity of 1 it creates the
   Run.
6. Success Metric: the Quaestor shape, a completed `task_01` and a pending
   `task_02` sharing a file, creates the Run and dispatches `task_02` with no
   collision text.
7. Success Metric: the two guides and the two Roundfix references named in
   Core Feature 5 carry both rules, and `CONTEXT.md` adds **Wave Collision**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Quaestor's report, quoted verbatim in both adopted Backlog Entries,
  written by another project before this Spec existed.
- The POSIX `make` specification (IEEE Std 1003.1-2024, `make`,
  https://pubs.opengroup.org/onlinepubs/9799919799.2024edition/utilities/make.html,
  read through Exa on 2026-10-30), which schedules only out-of-date targets
  in parallel: a target is brought up to date in parallel only when "The
  target's prerequisites (if any) are up-to-date" and make "is not currently
  bringing maxjobs targets up-to-date in parallel", among its other
  conditions.

## Unreachable Acceptance

- criterion: Quaestor, on the release that carries this Spec, authors a Spec
  with the templates' Tooling authority row unchanged and `paths: []`, and
  runs a corrective Task that shares a file with a completed Task, without
  either refusal.
  reason: it needs another repository's Specs and Runs after the release,
  which no Verification or QA row may read.
  satisfied-by: Quaestor's reply through ai-memory to the release message the
  operator sends, recorded session evidence.

## Validation

The decisions were checked on 2026-10-30 against these sources.

- Characterization at `35de7f74` in a scratch repository with the
  `bin/roundfix` of that commit: the templates' row citing a record with
  `paths: []` and `bounded files: none` reports `SC-TOOLING-UNBOUNDED`; the
  same row as `not applicable`, or without the record citation, passes. A
  binary built from `v0.124.0` reports the same finding, and
  `declaresProtectedToolingMutation` has not changed since v0.85.0, so no
  Roundfix release moved the verdict. A row whose bounded-files value is
  `none` followed by a parenthesized `paths: []`, the form Fiscus writes in
  its Specs 0052 to 0057, counts as bounded and passes in both versions; that
  is the likely reason Quaestor's Spec 0022 passed.
- Characterization through CLI dispatch on the implement test fixtures: two
  pending Tasks sharing a Verification operand at a Task Capacity of 2 exit
  `1` with `implement failed after Run start: task cycle: refuse Wave before
  dispatch because Tasks collide`, and a completed Task sharing the file
  with a pending one is refused the same way.
- A prototype of both Tasks' code in a second worktree passed those cases
  under the new rules. Across `internal/spec`, `internal/speccheck` and
  `internal/daemon` it moved only the `unordered corrective Task` case of
  `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph`, which `_techspec.md` →
  Declared breaks names.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix tooling authority row no protected tooling mutation SC-TOOLING-UNBOUNDED" --all --files --min-score 0.3`
  (the write-prd and write-techspec skill mirrors and two Fluxus tooling
  authorization references; the Fiscus PRDs 0052 and 0057, found by a direct
  search of the mirrors, carry that parenthesized form) and
  `qmd query "roundfix wave collision completed task implement refuses after run" --all --files --min-score 0.3`
  (the Roundfix `spec-delivery` and `deliver` reference mirrors). None
  records a decision on reading the row or on completed Tasks in a Wave.
- Exa: the POSIX `make` specification quoted above.

## Glossary

- adds: **Wave Collision**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AS brief for Spec 0334 of 2026-10-30 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-10-30), verbatim:

- Plan: "Aprovar como está (Recommended)" — option text: "Quatro Specs na
  ordem validada pelo Jev; a 0332 entra primeiro para que a fila resolva
  sozinha conflitos de registros gerados durante a 0333."
- Authorization: "Mesmas regras do AR (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."
- Release: "Uma no fim do ciclo (Recommended)" — option text: "v0.127.0
  depois das quatro Specs; os adotantes recebem a revisão de skills e as
  correções juntas."

Taken from those answers, ADR-0131 and the characterization, with no further
interview (ADR-0347):

- The checker moves, not the template: the row states applicability
  (ADR-0131), and a cited record with `paths: []` agrees with "no mutation".
  The templates gain only the row's ending, so the contract test can compose
  it.
- A completed Task never collides; a reopened one counts again.
- `implement` refuses a collision in Preflight Validation at a Task Capacity
  above 1; the Daemon's refusal stays as the second line.

## Open Questions

None.
