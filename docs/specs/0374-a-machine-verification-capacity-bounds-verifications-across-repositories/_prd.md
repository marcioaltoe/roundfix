---
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: active
created: 2026-11-10
surfaces: [backend, cli, docs]
type: feat
---

# A Machine Verification Capacity bounds Verifications across repositories

Requested by the maintainer through vortex on 2026-11-10: Runs of different
repositories verify at the same time on one machine, their test suites
compete for the same cores, and time-sensitive tests that a Task never touched
fail, so a Run ends unresolved and an Agent spends a repair turn on a failure
it cannot fix. The maintainer wants the projects' Agents to keep working in
parallel and only the verification step to take turns. After this Spec a
**Machine Verification Capacity**, one User Config value that defaults to one,
bounds how many Verifications the Daemon and the Delivery Queue run at once
on the machine across every repository, and the status surfaces name who is
verifying and who is waiting.

## Sources

One request, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md`](references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md),
   an ai-memory message from vortex of 2026-11-10 20:45 BRT, requested by the
   maintainer through vortex, copied verbatim with the operator's measured
   facts. The same request also feeds the later core budget Spec, which links
   this copy and adopts nothing.

At the authoring base `docs/backlog/` holds no entry that shares this
context, and `docs/findings/` is empty. This project's ai-memory inbox, listed
read-only on 2026-11-10, holds only the vortex message above; it was not
consumed.

## Prerequisites

None in `requires`. Authored against `60ac8142` (v0.137.8 plus Spec 0371).
Spec 0373 delivers before this one and changes, each in its own section,
files this Spec also changes:

- `CONTEXT.md`: 0373 its message and follow-up terms; this Spec the capacity,
  waiting and Run Budget terms.
- The Roundfix owned skill's version fields in the `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json`: 0373 changes the `follow-up`
  and `modes` references; this Spec the `implement`, `runs`, `deliver` and
  `setup` references. Each Task raises the version with the record command
  from the tree it runs on and never writes a number.
- `docs/references/behavior-surfaces.json`: each Task re-records it after its
  own last guide or key edit.

The later core budget Spec builds on this one and is not authored here.

## Project Constraints

- Identifier strategy: not applicable — no domain identifier is minted; the
  capacity's tickets are a local sequence that never leaves the machine, and
  the Run and Delivery Queue identifiers are the existing ones. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the capacity lives in the Roundfix Home of the machine.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0384 (new) records the decision:
  "A Machine Verification Capacity bounds the Verifications of every
  repository". ADR-0056 is amended: it decides "The capacity is owned by one
  Task cycle, not by a machine-wide service, so it does not coordinate other
  Runs", and every Daemon Verification attempt now also takes a unit of the
  machine's capacity. ADR-0373 applies unchanged: "A proven start-identity
  mismatch at a recorded owner PID means the owner is gone", the only rule
  that reclaims a dead holder. ADR-0044 applies unchanged: "a provably dead
  owner makes the lock orphaned", the precedent for reclaiming by proof.
  ADR-0304 applies unchanged: "`runs.slot_wait` sits beside `runs.max_active`",
  the User Config placement the new key follows, and its rejection of a
  schema change, which this Spec keeps. ADR-0365 applies unchanged: "The turn
  is a file.", the precedent for state outside the Run Database. ADR-0342
  applies unchanged: the Integrated Repository Verification is gated like any
  other Daemon Verification. ADR-0313, cited by ADR-0373, applies unchanged:
  nothing that signals changes. Not applying: ADR-0025 and ADR-0038, which
  ADR-0056 refines, keep their meaning, and ADR-0034, which ADR-0056 cites,
  governs Worktree Bootstrap, which this Spec does not gate; ADR-0311, which
  cites ADR-0304, sets the Delivery Queue's check wait, which this Spec does
  not change; and ADR-0182, which cites ADR-0342, decides the facts a Task
  settles on, which this Spec does not change. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answers of 2026-11-10 quoted in the authorization record,
  which grant the owned skills where the design needs them; the Governed Paths
  are the Roundfix skill files tasks 02 and 03 change. No Makefile, workflow, Go
  module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0374-a-machine-verification-capacity-bounds-verifications-across-repositories/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runs.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Goals

- Two repositories' Verifications no longer run at the same time on a machine
  that keeps the default, while their Agents keep working in parallel.
- A Run that waits for another repository's Verification loses no Run Budget
  and no Task for it, and a Stop Request still ends it cleanly.
- A holder that dies never blocks the machine, and a holder that is still
  alive is never displaced.
- Anyone reading `runs list`, `deliver status` or `doctor` sees who holds the
  capacity, for which step, and who waits in which order.
- A machine owner can raise the bound or turn it off in one User Config key.

## User Stories

1. As the maintainer running Delivery Queues of several repositories on one
   machine, I want only one repository to verify at a time, so that
   time-sensitive tests stop failing because another suite took the cores.
2. As the maintainer, I want the Agents of every Run to keep working while one
   repository verifies, so that the machine still makes progress on every
   project.
3. As a Supervisor deciding whether to resume or stop a Run, I want `runs
   list`, `deliver status` and `doctor` to name the Run or queue item holding
   the capacity and the ones waiting, so that I know why a Run is not
   verifying.
4. As a machine owner with spare cores, I want to raise the bound or turn it
   off in User Config, so that a repository cannot decide it for me and I can
   return to today's behavior.
5. As the operator, I want a Run whose owner died to give its unit back
   without my intervention, so that one crash does not stop every repository's
   verification.

## Core Features

1. `runs.max_verifying` in User Config is the Machine Verification Capacity: a
   non-negative integer, `1` by default, `0` turns it off. A negative or
   non-integer value fails configuration validation. A Project Config value
   is ignored with the warning the other operator-only keys print. The User
   Config that `roundfix init --scope user` and `roundfix setup` create
   writes the key with its comment.
2. Every Verification the Daemon starts takes one unit of the Run's own
   Verification Capacity and then one unit of the Machine Verification
   Capacity, and gives the machine unit back first: the pre-work probe, the
   repository precondition, each Task attempt with the repository Verification
   appended at settlement, the exclusive retry after exit `75`, the Integrated
   Repository Verification and the QA Task's repository Verification run
   before the QA Agent. The exclusive retry takes the Run's whole capacity and
   one machine unit.
3. The Delivery Queue's repository gate takes one unit for its repository
   Verification and its Linux Lane together and gives it back when the gate
   ends, whatever its outcome.
4. Waiters are served in arrival order across every repository and process
   sharing the Roundfix Home: a waiter takes a unit only when the holders and
   the waiters ahead of it leave one free.
5. Before it counts, a waiter removes each holder or waiter whose recorded
   owner process is proven gone (absent, or a reused PID proven by its start
   identity). An owner that cannot be proven gone keeps its place, however
   long it has held it. A wait that is cancelled or stopped removes its own
   entry, and a release removes it even while the Run is being cancelled.
6. While any Verification of a Run waits for the Machine Verification
   Capacity, that Run's Run Budget does not run; the wait moves its deadline
   later by the wait's length. A Stop Request during the wait starts no
   command and leaves the Task resumable.
7. When a Verification begins to wait, the Run journals one Run Event naming
   the step, the capacity and each holder, and prints one progress line. The
   Delivery Queue owner prints one console line when its gate begins to wait.
8. `roundfix runs list`, `roundfix deliver status` and `roundfix doctor` name
   each holder and each waiter, with its Run or queue item, repository, Spec,
   step, place in line and start time; `runs list` and `deliver status` print
   nothing new when nobody holds or waits. They write nothing and leave out an
   entry whose owner is proven gone.
9. With `runs.max_verifying: 0` no entry is written and every Verification
   starts exactly as today.

## User Experience

When Vortex's Task waits while Roundfix verifies, `roundfix runs list --all
--format plain` ends with:

```text
Machine Verification Capacity: 1 of 1 held, 1 waiting
  holding: Run run_20261110T215130Z_fa9f763e0a43c500 of 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line in /Users/marcio/dev/roundfix, Task task_02 Verification attempt 1, since 2026-11-10T22:41:30Z
  waiting 1: Run run_20261110T221208Z_946b774b8b696f09 of 0179-example in /Users/marcio/dev/vortex, Task task_03 Verification attempt 2, since 2026-11-10T22:50:02Z
```

`roundfix deliver status` prints the same block after the queue's items,
`roundfix doctor` prints a `verification-capacity:` line, and the waiting
Run's progress shows that its Verification waits for the Machine
Verification Capacity and who holds it.

## Non-Goals / Out of Scope

- Bounding the commands an Agent runs inside its session, including the QA
  Agent's own suites: the Daemon never sees them.
- Gating `roundfix spec check --run-verification`, `roundfix settle`, the
  Delivery Queue's item binary build and derived-file regeneration, Worktree
  Bootstrap, or the Verification of review Runs (`fetch`, `resolve`,
  `watch`).
- A core budget exported to Verification commands (`GOMAXPROCS` or a Roundfix
  variable): the next Spec, built on this one.
- Reclaiming a holder by elapsed time, a heartbeat or a timeout, and adding a
  time limit to a Verification command.
- A Run Database schema change, a new Run state or a new Delivery Queue stage.
- Changing `runs.max_active`, `runs.slot_wait`, `verification.concurrency`,
  `worktree.concurrency` or the Run Slot Wait.
- Coordinating CI jobs, other machines, other Roundfix Homes or processes
  outside Roundfix.
- Telling adopters by message: they read the change in the release notes, the
  user guide and the Update Check.

## Success Metrics

1. Success Metric: with a capacity of `1`, two holders requested at once from
   two processes' worth of entries never hold a unit together, and a third
   waiter is granted only after both, in ticket order; at `60ac8142` no such
   capacity exists and `runs.max_verifying` is refused as an unknown key.
2. Success Metric: an entry whose owner PID is reused (a live PID with a
   different start identity) or absent is removed by the next waiter, which
   then takes the unit; an entry with an unreadable, foreign or empty
   identity on a live PID keeps it.
3. Success Metric: in a Task cycle with a fake gate, every Daemon Verification
   step (pre-work probe, repository precondition, Task attempt with the
   settlement command, exclusive retry, Integrated Repository Verification,
   QA repository Verification) asks for the machine unit after the Run's own
   unit and starts no command before it is granted.
4. Success Metric: a Stop Request during a machine wait starts no command,
   and a one-minute Run Budget paused for a two-minute machine wait expires
   only once one minute of time outside the wait has passed.
5. Success Metric: the Delivery Queue's repository gate command observes its
   own entry in the holding state while it runs, and no entry remains after
   the gate passes or fails.
6. Success Metric: `runs list`, `deliver status` and `doctor` print the
   prescribed lines for a holder and a waiter in a temporary Roundfix Home and
   print no capacity block when it is empty; every existing test of their
   output passes unchanged.
7. Success Metric: the configuration guide, the command guides, both copies
   of the four Roundfix references, `CONTEXT.md` and ADR-0056's amendment
   notice describe the capacity, and the documentation, skill and glossary
   contracts pass.

## Acceptance evidence

The outside-evidence rows rest on a measurement this Spec did not design and
on published literature read through Exa on 2026-11-10.

- The vortex measurement of 2026-11-10, committed with this Spec's planning
  change in `references/` and copied verbatim from vortex's message: three
  `run-unresolved` stops "sempre na verificação do repositório no assentamento
  da Task", one test that "estourou 5000 ms. Sozinho, o mesmo teste leva 1,3
  s", and "Com a máquina sem outra Run, a mesma árvore do Vortex passou
  inteira: 771 arquivos, 7.140 testes, 285 s". This is why the bound is on the
  Verification step and not on the Runs.
- Michał Górny, "One jobserver to rule them all"
  (https://blogs.gentoo.org/mgorny/2025/11/30/one-jobserver-to-rule-them-all/):
  "if you run multiple make processes in parallel, and they run multiple jobs
  simultaneously, you end up overscheduling", and "clients are responsible for
  acquiring a token whenever they are about to start a job, and returning it
  once the job finishes. The availability of job tokens therefore limits the
  total number of processes started." This is why independent processes
  share one fixed pool of units that each takes before it starts and returns
  when it ends.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, a Run of this
  repository and a Run of another repository on the maintainer's machine
  never verify at the same time with the default capacity, and `roundfix runs
  list` names the waiting one.
  reason: it needs two live Runs of two real repositories with real Agents on
  the maintainer's machine and its real Roundfix Home, which no Verification
  or QA row may reach.
  satisfied-by: the operator's runs list and doctor evidence from the first
  cycle after the release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Characterization at `60ac8142`:
  - every Daemon Verification goes through one function that takes the Run's
    own Verification Capacity, `acquireVerificationCapacity`, reached by the
    pre-work probe and by every attempt request (Task attempt, exclusive
    retry, repository precondition, Integrated Repository Verification, QA
    repository Verification); the per-Run gate is an in-process object, so two
    Runs never see each other's attempts;
  - the Delivery Queue's repository gate runs the configured Verification and
    then the Linux Lane in the Queue Owner, outside any Run;
  - a User Config with `runs.max_verifying` is refused as an unknown key;
  - `runs.max_active` is User Config only and ignored with a warning in
    Project Config; the user template writes it with a comment;
  - the Run Budget is a deadline that renews at each Task settlement and
    cancels the Task cycle when it passes; no Verification command has its
    own time limit;
  - `runs list` and `deliver status` print nothing about other repositories'
    Verifications, and Vortex read `Verifying` for 50 minutes without knowing
    whether that Run was running a suite.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for shared
  verification capacity and for reused owner PIDs. The relevant results were
  the conexus mirror `projects/conexus/mirror/docs/references/test-execution-capacity.md`
  (Vitest `maxWorkers: 2` after joint runs timed out: "o resultado sustenta
  sensibilidade a concorrência/recursos"), the vortex mirror
  `projects/vortex/mirror/docs/references/2026-11-10-autorizacao-limite-de-tempo-dos-testes-do-backend.md`
  (the 30 s `testTimeout` vortex adopted the same day), this repository's
  ADR-0056 mirror and `wiki/sources/digest-roundfix-pid-reutilizado-e-cache-go-no-sandbox-2026-10-12.md`
  ("um kill por PID gravado precisa comparar a identidade de início"). They
  confirm that each project already limited itself and still collided, and
  that a recorded PID needs its start identity.
- Exa: the jobserver article above and the GNU make jobserver paper
  (https://make.mad-scientist.net/papers/jobserver-implementation/), which
  share a fixed number of tokens among independent processes; and a commit of
  an unrelated project's SQLite semaphore
  (https://github.com/browser-cli/browser-cli/commit/e8e106d11fa3d86c794352751360a0f6b2367b4d)
  whose PID-only liveness check kept a slot "held" forever after PID reuse and
  whose fix reclaims by a 180 s heartbeat age. It challenged the design; this
  Spec keeps proof of death over age, because a Verification can legitimately
  run for an hour, and ADR-0373 already proves reuse by start identity.

## Glossary

- adds: **Machine Verification Capacity**
- changes: **Verification Capacity**
- changes: **Waiting for Verification**
- changes: **Run Budget**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0374 of 2026-11-10, the maintainer's answers below and the adopted request.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~20:55–21:00 BRT), verbatim:

- Question "Pedido do Vortex (Spec 0374): verificações de projetos diferentes
  rodando juntas nesta máquina derrubam testes sensíveis a tempo. Hoje o
  Roundfix só limita Runs inteiras (`runs.max_active`), e você não quer isso.
  Que mecanismo o Roundfix ganha?" → "Os dois" — option text: "Capacidade da
  máquina para serializar a verificação e, além dela, o orçamento exportado
  para as verificações que ainda rodarem juntas. Spec maior, provavelmente
  dividida em duas." The two options it combines: "Capacidade da máquina" —
  "Uma chave na config de usuário limita quantas verificações rodam ao mesmo
  tempo na máquina, somando todos os repositórios. Os agentes seguem em
  paralelo; só a etapa de verificação espera a vez, e `deliver status` e
  `runs list` mostram quem espera. Vale para qualquer stack (Go, Vitest,
  Bun)."; "Orçamento de núcleos" — "O Roundfix divide os núcleos entre as Runs
  ativas e exporta o número (GOMAXPROCS e uma variável própria) para os
  comandos de verificação. O Go obedece sozinho; Vitest e Bun só se cada
  projeto ler a variável. As verificações continuam simultâneas, só mais
  leves."
- Question "A Capacidade da máquina (Spec 0374) é uma chave da config de
  usuário. Que valor ela tem por padrão para quem atualizar o Roundfix sem
  configurá-la?" → "1 por padrão (Recommended)" — option text: "Toda máquina
  passa a verificar um repositório por vez, e os agentes seguem em paralelo.
  O `setup` grava o valor na config de usuário. Muda o comportamento de todos
  os adotantes logo que atualizarem."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those answers, the request and the characterization, with no
further interview (see ADR-0384):

- The key is `runs.max_verifying`, beside `runs.max_active` and
  `runs.slot_wait`, User Config only, and `0` turns it off as `max_active: 0`
  does.
- The capacity wraps the Run's own Verification Capacity: the Run's unit
  first, the machine's second, released in reverse, so no Run holds a machine
  unit while it waits for its own.
- The entries are files in the Roundfix Home under one advisory lock, not a
  Run Database table, so no older binary refuses the shared Run Database.
- Waiters are served in ticket order; a holder is removed only by ADR-0373's
  proof that its owner is gone.
- The Run Budget pauses during the wait; no other time limit changes.
- The release notes and the configuration guide say that every adopter's
  Verifications become one at a time after the update; the release that
  carries this Spec carries the core budget Spec too.
- ADR-0384 records the decision and amends ADR-0056. Four implementation
  Tasks in series and a QA gate.

## Open Questions

None.
