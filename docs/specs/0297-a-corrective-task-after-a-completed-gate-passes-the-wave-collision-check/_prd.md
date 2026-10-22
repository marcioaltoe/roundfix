---
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# A corrective Task after a completed gate passes the Wave collision check

Minimal PRD for a bug fix, minted by `write-techspec`. Oraculum reported that
its QA gate refused a corrective Task four times with `SC-WAVE-COLLISION`
although the Task's `needs` edge was committed and `roundfix spec check
--strict` passed on the same head. Characterization on oraculum's own
history and in a temporary repository shows the gate refused once, correctly,
before the edge existed; the three later gates passed that check and failed
on two documentation findings, and the first refusal report, carried into each
later Run Branch beside the newer report, was read as current. The gate and
the strict check already judge the same working-tree Task Graph and settlement
history. What the refusal lacked was the remedy: the finding the gate records
names no edge, while the Daemon's pre-dispatch refusal does. After this Spec
every `SC-WAVE-COLLISION` message names the one `needs` edge that orders the
pair, built by one helper, and tests pin the corrective shape so the gate, the
strict check, the settlement checks and the pre-dispatch refusal keep agreeing.
The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-21-sc-wave-collision-no-portao-com-needs-e-check-estrito-limpo.md`](references/2026-10-21-sc-wave-collision-no-portao-com-needs-e-check-estrito-limpo.md),
   origin oraculum, type hint `fix`, 2026-10-21. Spec 0087 there gained a
   corrective `task_05` after a pre-PR review sent an archived candidate to
   `corrective-spec-required`; the entry reports four gate refusals with
   `SC-WAVE-COLLISION: Task task_02 and Task task_05 can run in the same Wave
   and share repository files: packages/backend/test/rateio-repartido-versoes.test.ts
   (declared context; prior Run settlement commit)` and asks that the gate
   read the current graph or that the message name the edge that resolves it.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry about the
Task Graph, Wave collisions or the QA gate precondition.

## Prerequisites

None in `requires`. This Spec is authored against `34b4e246` (v0.107.0), with
no Spec active, and delivers first of cycle AO, before Specs 0298, 0299 and
0300. It changes `internal/spec/collision.go`, `internal/speccheck/coherence.go`
and `refuseTaskPlanWaveCollisions` in `internal/daemon/task_engine.go`; none
of the other cycle AO Specs is planned to touch them. Its documentation Task
raises the Roundfix skill version over whatever the tree holds when it starts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted; the message names Task IDs the graph already holds.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and Git configuration. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0117: "A defect is checked by the
  stage that can produce it", which this Spec keeps by holding the authoring
  check and the gate precondition to one judgement; ADR-0182: "A Task settles
  on the facts its gate will check", which this Spec keeps by proving that
  the settlement checks report the same Wave collision text the gate records.
  ADR-0096, which ADR-0117 and ADR-0182 cite, is not changed: the mechanical
  stage and its refusal record keep their shape. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 through AskUserQuestion, authorization "the same
  rules as cycle AN", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check/_authorization.md`;
  bounded files: `internal/speccheck/coherence.go`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. No
  Makefile, workflow, Go module, lint, formatter, Baseline module or Project
  Config file change.

## Goals

- Every `SC-WAVE-COLLISION` message, the Spec Consistency Check finding with
  the refusal reason the QA gate records and the Daemon's pre-dispatch
  refusal, names the one `needs` edge that orders the pair, in the same words.
- A corrective Task whose `needs` orders it after the Task it shares a file
  with passes the gate precondition, the strict check, the settlement checks
  and the pre-dispatch refusal, in the shape oraculum met: a prior Run Branch
  with settlement commits, an archive and its revert, and the corrective Task
  added afterwards.
- Two Tasks the current graph lets run in one Wave that share a file are
  still refused by all four.
- The `spec` user guide and the Roundfix `spec` reference state the finding,
  its evidence sources and its remedy.

## Core Features

1. One resolving edge. `spec.WaveCollision` gains `ResolvingEdge`, the text
   ``add a `needs` edge from Task <Second> to Task <First>``, and
   `SharedPathList`, the sorted `<path> (<source>)` list; the Spec
   Consistency Check and the Daemon's pre-dispatch refusal both build their
   text from them.
2. The finding names the edge. The `SC-WAVE-COLLISION` summary ends with `;
   ` and the resolving edge, so the refusal reason a QA Report records names
   it, and the fix line names the same two Tasks.
3. The corrective shape is pinned. A Daemon-package test builds oraculum's
   shape and proves that the gate precondition, the strict check in the item
   checkout, the settlement check and the pre-dispatch refusal agree, both
   before the edge (all refuse with the edge named) and after it (none
   refuses).
4. Documentation. The `spec` user guide and the Roundfix `spec` reference
   describe `SC-WAVE-COLLISION`.

## Non-Goals / Out of Scope

- Changing which pairs collide: the evidence sources, the dependency closure,
  the newest-settlement rule and the serial-capacity exemption stay as they
  are.
- Changing the QA Report's precondition refusal shape, the Prior QA Report
  import, or the names of QA Reports.
- Resuming a Delivery Queue item whose Spec is archived on its branch, which
  Spec 0299 owns.
- Rewriting a QA Report or a commit already made.

## Success Metrics

1. Success Metric: in a temporary repository shaped like oraculum's Spec 0087
   with `task_05 needs: [task_01]`, the gate precondition run in a linked Run
   Worktree, the strict Spec check run in the item checkout and the
   settlement check each report exactly one `SC-WAVE-COLLISION`, whose
   refusal reason is `SC-WAVE-COLLISION: Task task_02 and Task task_05 can
   run in the same Wave and share repository files: pkg/versions.test.ts
   (declared context; prior Run settlement commit); add a `needs` edge from
   Task task_05 to Task task_02`, and the pre-dispatch refusal at Task
   Capacity 2 names the same path, sources and edge.
2. Success Metric: in the same shape with `task_05 needs: [task_02]` and a
   second Run Branch settling `task_05` on the shared file, none of the four
   reports a Wave collision.
3. Success Metric: for a pair in either manifest order, adding the edge
   `ResolvingEdge` names removes that collision and forms no cycle.
4. Success Metric: `TestTaskCycleRefusesAWaveCollision` and the existing Wave
   collision tests of `internal/spec` and `internal/speccheck` pass unchanged.
5. Success Metric: after release, the next corrective Task an adopting
   repository adds after a completed gate, with the edge the finding names,
   passes the gate's precondition without a QA Archive Override.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- This repository's history: Spec 0296's Archive Record
  (`docs/history/specs/0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters.md`)
  names `source_revision: 3dc0b8dd17eeb4cebc861e71ef4cb16d8a079438`. There a
  corrective `task_04` declared `internal/cli/deliver_queue_subject_test.go`
  with `needs: [task_02]`, which created that file, and the QA Task `task_03`
  needs `task_04`; `roundfix spec check` on that tree reports no finding, and
  the gate reached its matrix and recorded `qa_verdict: partial`.
- The GCC internals manual, "Guidelines for Diagnostics", section 29.1.2
  (https://gcc.gnu.org/onlinedocs/gccint/Guidelines-for-Diagnostics.html):
  "A good diagnostic is actionable: it should assist the user in taking
  action." and "Given an error, an end user will think: “How do I fix
  this?”".
- Barik, Ford, Murphy-Hill and Parnin, "How Should Compilers Explain Problems
  to Developers?", ESEC/FSE 2018
  (https://static.barik.net/barik/publications/fse2018/barik_fse18.pdf):
  developers "will accept a deficient argument structure if it provides a
  resolution to the problem".

## Unreachable Acceptance

- criterion: Success Metric 5, after release, the next corrective Task an
  adopting repository adds after a completed gate, with the edge the finding
  names, passes the gate's precondition without a QA Archive Override.
  reason: only an adopter's own Run on its own Spec shows it, and no
  Verification or QA row may read another repository, the real `~/.roundfix`
  or a provider.
  satisfied-by: the adopter operator's newest QA Report of that Spec after
  this Spec's release, whose frontmatter carries no `precondition_reason`
  naming `SC-WAVE-COLLISION`.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `34b4e246`, recorded in the TechSpec's Current
  behavior. Oraculum's history was read without writing to it, through a
  mirror clone in the operator's scratch space: its item branch shows the
  `SC-WAVE-COLLISION` refusal only in `qa-report-2026-10-21-01.md`, the Run
  at `57a7c55` whose `task_05` still needed `task_01`; reports `-02`, `-03`
  and `-04`, built at `8736693`, `7472ad34` and `75a9dd0`, all with
  `task_05 needs: [task_02]`, reached their matrices, recorded the strict
  precondition passing and failed on two documentation findings. The current
  binary run on those four heads reports the collision on `57a7c55` only. A
  temporary repository built in oraculum's shape reproduces the exact refusal
  text through `roundfix spec check` and the Daemon's
  `specGatePrecondition`, and both stop reporting it once the edge is
  committed.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "wave
  collision needs edge corrective task QA gate precondition" --all --files
  --min-score 0.3`, which returned the `write-tasks` skill mirror, the fleet's
  `autonomous-work.md` copies ("When a corrective Task is added after the gate
  settled, add it to the QA Task's `needs`") and Roundfix's archived Spec 0097
  record, which introduced the collision rule; and `qmd query "QA report
  precondition refusal imported prior report newest sequence misread" --all
  --files --min-score 0.3`, which returned the Roundfix mirror's
  `qa-report.md` and `context-driven-development.md` guides ("only the newest
  report in the directory is read") and the ADRs on refused gates and the
  prior-pass import. No prior decision on the collision message exists. No
  Inbox Entry was written, because the briefing forbids writing there.
- Exa: the GCC guideline and the Barik et al. study above, and the Rust
  compiler development guide
  (https://rustc-dev-guide.rust-lang.org/diagnostics): "The error or warning
  portion should not suggest how to fix the problem, only the “help”
  sub-diagnostic should." Challenging: a Roundfix finding already has a
  separate fix line, the analogue of Rust's help, but the QA Report's
  `precondition_reason` records only the code and the summary, and the
  Daemon's refusal already carries the edge in its sentence; so the summary
  carries the edge as well and the fix line repeats it.

## Glossary

None.

## Decisions

- 2026-10-22 (about 00:55 BRT): the cycle AO scope, multi-select, "Z: portão
  e corretiva", "Espera de checks", "Retomar item arquivado", "Y: reconcile
  PID", one Spec and one minor each, in that order; the authorization "the
  same rules as cycle AN".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AO briefing for Spec 0297 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The briefing's decision has two parts: the gate and the strict check judge
  the same current graph and evidence, and every message names the edge. The
  first already holds; this Spec pins it with tests in oraculum's shape
  instead of changing the collision path, and fixes the second.
- The briefing's hypothesis that a gate reader loads the Task Graph or the
  settlement evidence from an archived record's `source_revision` was ruled
  out: the gate loads `_tasks.md` from the Run Worktree and reads settlement
  commits with `git log --all`, exactly as `roundfix spec check` does.
- The named edge makes the later Task in graph order need the earlier one,
  the direction the Daemon's refusal already names; neither Task reaches the
  other, so it cannot form a cycle.
- The Daemon's pre-dispatch refusal keeps its bytes and reads its text from
  the shared helper.
- No new ADR is written: no decision changes, and only the message gains
  its remedy.
- The precondition refusal report carries no audited head, which is what let
  an imported refusal pass for a current one; changing that report shape is a
  product decision outside the briefing and is left out of this Spec.

## Open Questions

None.
