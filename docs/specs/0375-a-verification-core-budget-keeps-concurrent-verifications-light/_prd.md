---
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: active
created: 2026-11-10
surfaces: [backend, cli, docs]
type: feat
---

# A Verification Core Budget keeps concurrent Verifications light

Requested by the maintainer through vortex on 2026-11-10: `go test` and most
test runners use every core of the machine, so a Verification that still runs
beside another one, or beside an Agent's own suite, slows the other enough to
fail its time-sensitive tests. The Machine Verification Capacity of Spec 0374
takes turns for the Verifications the Daemon and the Delivery Queue start, and
the maintainer asked for the budget exported to the Verifications that still
run together as well. After this Spec every Verification and Agent Session of
an Implement Run, and the Delivery Queue's repository gate, receive a
**Verification Core Budget**, half of the machine's logical CPUs divided by
the Machine Verification Capacity, as `GOMAXPROCS` and as
`ROUNDFIX_VERIFICATION_CORES`, and `roundfix doctor` reports the number
against the active Implement Runs.

## Sources

One request, owned by Spec 0374, which adopted it; this Spec links that copy
and adopts nothing:

1. [`../0374-a-machine-verification-capacity-bounds-verifications-across-repositories/references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md`](../0374-a-machine-verification-capacity-bounds-verifications-across-repositories/references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md),
   an ai-memory message from vortex of 2026-11-10 20:45 BRT, requested by the
   maintainer through vortex, with the operator's measured facts. Its second
   request names "um orçamento de núcleos por Run que o Roundfix exporta para
   os comandos de Verification e que o `doctor` confere contra o número de
   Runs ativas".

At the authoring base `docs/backlog/` holds four open entries, none about
cores or test parallelism, and `docs/findings/` is empty. This project's
ai-memory inbox, listed read-only on 2026-11-10 (BRT), held no pending
message; nothing was consumed.

## Prerequisites

`requires` names Spec 0374: this Spec divides by its `runs.max_verifying`,
places its doctor line after 0374's `verification-capacity:` line and its
configuration section after 0374's `## Machine Verification Capacity`.
Authored against `b233761e` (v0.137.9 plus the 0374 plan). Specs 0373 and
0374 deliver before this one and change, each in its own section, files this
Spec also changes:

- `CONTEXT.md`: 0373 its message terms, 0374 its capacity terms; this Spec
  adds one term after 0374's Machine Verification Capacity.
- The Roundfix owned skill's version fields in the `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json`: 0373 changes the `follow-up`
  and `modes` references, 0374 the `implement`, `runs`, `deliver` and `setup`
  references; this Spec changes the `implement`, `deliver` and `setup`
  references again, in sections of their own. Each Task raises the version
  with the record command from the tree it runs on and never writes a number.
- `docs/references/behavior-surfaces.json`: each Task re-records it after its
  own last guide edit.
- `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md` and
  `docs/user-guide/commands/deliver.md`, which 0374 also changes.
- The Daemon's dependencies and Verification request, the Implement
  Command's engine wiring, the Delivery Queue gate and the doctor command and
  its check names, which 0374 also changes in other functions or other lines.

## Project Constraints

- Identifier strategy: not applicable — no domain identifier is minted; the
  budget is a number computed per process and never stored. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the two variables carry a number, never a secret. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0385 (new) records the decision:
  "Roundfix computes one number and exports it under two names." ADR-0384
  applies unchanged: "`runs.max_verifying` is the Machine Verification
  Capacity", the bound the budget divides by. Not applying: ADR-0304, which
  ADR-0384 cites, sets the Run Slot Wait, which this Spec does not change;
  ADR-0365, which ADR-0384 cites, sets the pre-PR review's own capacity,
  which this Spec does not change; ADR-0373, which ADR-0384 cites, decides
  when a recorded owner is gone, which this Spec does not read; and
  ADR-0056, whose amendment by the capacity Spec cites ADR-0384, sets the
  per-Run Verification Capacity, which this Spec does not change. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answers of 2026-11-10 quoted in the authorization record,
  which grant the owned skills where the design needs them; the Governed Paths
  are the Roundfix skill files tasks 02 and 03 change. No Makefile, workflow,
  Go module, lint, formatter, test-runner or Baseline change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Go Verification on a machine that keeps the defaults uses half of the
  machine's cores, so an Agent's own suite or a second Verification still
  finds free cores.
- When the owner lets several Verifications run at once, their budgets
  together still fit in half of the machine.
- Any test runner, not only Go, can read the number from one variable.
- Whoever sets `GOMAXPROCS` explicitly keeps the number they chose.
- Anyone reading `roundfix doctor` sees the budget, where it came from and
  whether the bounds can oversubscribe the machine.

## User Stories

1. As the maintainer running Delivery Queues of several repositories on one
   machine, I want each Verification to use part of the cores, so that a
   Verification that still runs beside another one does not starve it.
2. As the maintainer, I want the QA Agent's own suites, which the Machine
   Verification Capacity cannot gate, to receive the same budget, so that they
   stay light too.
3. As the maintainer of a repository whose tests do not run on Go, I want a
   Roundfix variable that names the budget, so that my test runner can bound
   its workers by it.
4. As an operator who exports `GOMAXPROCS` for a reason of their own, I want
   Roundfix to keep my number, so that its default does not override my
   choice.
5. As a Supervisor deciding whether the machine is overloaded, I want `doctor`
   to show the budget against the active Implement Runs, so that I can see
   whether the bounds let Verifications take more cores than the machine has.

## Core Features

1. The Verification Core Budget is half of the machine's logical CPUs,
   rounded down and at least 1, divided by `runs.max_verifying` when it is
   above `0`, otherwise by `runs.max_active` when that is above `0`, otherwise
   not divided, rounded down and never below 1.
2. A `GOMAXPROCS` that the environment of the Roundfix process holds as a
   positive whole number is the budget instead; any other value is ignored.
3. Every Verification command the Daemon starts for an Implement Run — the
   pre-work probe, the repository precondition, each Task attempt with the
   repository Verification appended at settlement, the exclusive retry, the
   Integrated Repository Verification and the QA Task's repository
   Verification — runs with `GOMAXPROCS` and `ROUNDFIX_VERIFICATION_CORES`
   both set to the budget, replacing any inherited copy.
4. Every Agent Session of an Implement Run receives both variables with the
   budget, beside its Run Marker.
5. The Delivery Queue's repository Verification receives both variables, and
   its Linux Lane container receives both as container environment; `roundfix
   deliver lane` gives its container the same.
6. `roundfix doctor` prints one `verification-cores:` line after the
   `verification-capacity:` line: the budget, its source, and the number of
   active Implement Runs. It warns, without changing the exit code, when the
   Verifications the bounds allow at once, each with the budget, would need
   more logical CPUs than the machine has, and when the Run Database cannot be
   read.
7. No User Config or Project Config key is added, and no other command,
   Verification or child process receives the variables from Roundfix.

## User Experience

On the maintainer's twelve-core machine with the defaults, `roundfix doctor`
prints:

```text
verification-cores: ok (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; 2 active Implement Runs)
```

A Go Verification then runs `go test` with six package binaries at most, and a
Vitest configuration that reads `ROUNDFIX_VERIFICATION_CORES` into its
`maxWorkers` runs six workers. With `runs.max_verifying` and
`runs.max_active` both at `0` and three active Implement Runs, the line warns
that three Verifications at once may take eighteen of twelve logical CPUs.

## Non-Goals / Out of Scope

- A Baseline clause, a Setup Snapshot or a generated test-runner
  configuration that makes an adopter's test runner read the variable: each
  repository opts in through its own authorized tooling change, guided by the
  user guide.
- Changing this repository's Makefile, its `-parallel` setting or any test
  runner configuration.
- A User Config key for the share, a per-repository budget, or a budget that
  changes while a Run runs.
- Exporting the variables to `roundfix spec check --run-verification`,
  `roundfix settle`, the item binary build, derived-file regeneration,
  Worktree Bootstrap, review Runs or any command the user runs by hand.
- Limiting cores with operating-system controls such as CPU affinity, nice
  levels or container CPU quotas.
- Telling adopters by message: they read the change in the release notes, the
  user guide and the Update Check.

## Success Metrics

1. Success Metric: on twelve logical CPUs the budget is 6 with
   `runs.max_verifying` at `1`, 3 at `2`, 2 with the capacity off and
   `runs.max_active` at `3`, and 6 with both off; it is never below 1, and an
   explicit `GOMAXPROCS=4` makes it 4 while `0`, `-2`, `+4`, `four` and an
   empty value are ignored. At `b233761e` no such number exists.
2. Success Metric: a Verification command the Daemon starts and an Agent
   Session of the same Implement Run each see `GOMAXPROCS` and
   `ROUNDFIX_VERIFICATION_CORES` equal to the budget, with exactly one
   `GOMAXPROCS` entry; with no budget the command sees the inherited
   environment unchanged.
3. Success Metric: the Delivery Queue's repository gate command sees both
   variables equal to the budget, and the Linux Lane's container is started
   with both as container environment.
4. Success Metric: `roundfix doctor` prints the prescribed
   `verification-cores:` line for each source and warns exactly when the
   bounded number of Verifications times the budget exceeds the logical CPUs;
   every existing doctor test passes unchanged.
5. Success Metric: the configuration guide, the doctor and deliver command
   guides, both copies of the three Roundfix references and `CONTEXT.md`
   describe the budget, and the documentation, skill and glossary contracts
   pass.

## Acceptance evidence

The outside-evidence rows rest on a measurement this Spec did not design, on
the Go toolchain's own documentation and on published sources read through
Exa on 2026-11-10.

- The vortex measurement of 2026-11-10 in Spec 0374's adopted copy: "A suíte
  do Vortex já se limita a 4 processos (vitest `maxWorkers`). O `go test` usa
  todos os núcleos por padrão (`-p` e GOMAXPROCS iguais ao número de CPUs)."
  This is why the budget goes to `GOMAXPROCS` and to a variable a non-Go
  runner can read.
- The Go command's own help, `go help build`, which the QA gate reads offline:
  "The default is GOMAXPROCS, normally the number of CPUs available." This is
  why setting `GOMAXPROCS` bounds both the package binaries `go test` starts
  and the threads each one runs.
- The Go runtime package documentation
  (https://pkg.go.dev/runtime): "If the GOMAXPROCS environment variable is set
  to a positive whole number, GOMAXPROCS defaults to that value." This is why
  only a positive whole number counts as an explicit choice.
- Jest's worker default (https://github.com/facebook/jest/blob/master/packages/jest-config/src/getMaxWorkers.ts):
  "In watch mode, Jest should be unobtrusive and not use all available CPUs.",
  followed by `Math.floor(numCpus / 2)`; and Vitest's `maxWorkers`
  (https://vitest.dev/config/maxworkers): "if `watch` is enabled, uses half of
  all available parallelism". This is why half of the logical CPUs is the
  share a run beside other work takes.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, `roundfix doctor` on
  the maintainer's twelve-core machine reports a Verification Core Budget of 6
  divided by `runs.max_verifying` 1, and a Roundfix Run's Verification on that
  machine runs with that budget.
  reason: it needs the released binary, the maintainer's machine and its real
  Roundfix Home and User Config, which no Verification or QA row may reach.
  satisfied-by: the operator's doctor line and a Verification diagnostic from
  the first cycle after the release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Characterization at `b233761e`:
  - the Daemon starts every Verification command through one verifier that
    passes the parent environment, without the judge's key variables and with
    the Run Marker, to `sh -c`; nothing sets `GOMAXPROCS`;
  - the Agent Sessions of a Run receive the Run Marker through the acpx
    command environment, and nothing else of the Run's;
  - the Delivery Queue's repository gate runs the same verifier in the Queue
    Owner, then the Linux Lane, whose container environment is fixed by its
    argument list;
  - this repository's Makefile runs `go test -parallel 16` without `-p`, so
    `go test` runs as many package binaries as `GOMAXPROCS`, twelve on the
    maintainer's machine;
  - `roundfix doctor` says nothing about cores.
- A prototype of tasks 01 to 03 on top of a prototype of Spec 0374 passed the
  `internal/corebudget`, `internal/daemon` and `internal/agent` packages and
  the Implement, Doctor, Deliver and Linux Lane tests of `internal/cli`;
  only the Linux Lane test that pins the container's arguments changed.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for a core budget
  shared by test verifications and for test-runner workers under machine
  load. See the TechSpec's Validation for the files it returned and how they
  informed the design.
- Exa: the Go runtime documentation, the Go blog's container-aware
  `GOMAXPROCS` post (https://go.dev/blog/container-aware-gomaxprocs), Jest's
  worker default, Vitest's `maxWorkers` page and a Jest pull request that
  reads `JEST_MAX_WORKERS` from the environment
  (https://github.com/jestjs/jest/pull/15944). They confirm that Go obeys the
  variable alone, that the common non-Go runners take half the CPUs when they
  share the machine and read a count from configuration, and that runners can
  read a count from an environment variable.

## Glossary

- adds: **Verification Core Budget**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0375 of 2026-11-10, the maintainer's answers below, Spec 0374's plan and the adopted request.

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
- Question "Até a 0374 existir, a verificação da 0373 do Roundfix vai rodar
  esta noite em paralelo com a fila do Vortex. O que faço com a fila do
  Roundfix?" → custom answer "GOMAXPROCS=4 até o final das implementações e
  depois usar 6", confirmed by "Isso (Recommended)" to "GOMAXPROCS=4 em toda
  Run do Roundfix até a última Spec deste ciclo (0373 a 0379) ser entregue, e
  depois 6 como valor fixo da fila do Roundfix nesta máquina, até o orçamento
  de núcleos da 0375 assumir."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those answers, the request, Spec 0374's plan and the
characterization, with no further interview (see ADR-0385):

- The share is half of the logical CPUs, at least 1: the maintainer's six of
  twelve, and the default the common test runners take beside other work.
- The share divides by `runs.max_verifying`, or by `runs.max_active` while
  the capacity is off, so the budgets that may run at once fit in the share;
  the divisor never changes while a Run runs.
- A positive whole-number `GOMAXPROCS` in the environment of the Roundfix
  process wins, so the operator's interim measure and any deliberate choice
  keep working; both variables then carry it.
- The Agent Sessions of an Implement Run receive the budget, because the QA
  Agent's suites are the load the Machine Verification Capacity cannot gate.
- `roundfix doctor` reports the budget against the active Implement Runs and
  warns on oversubscription only.
- The Baseline does not change: telling adopters' test runners to read the
  variable would ask their agents for a test-runner change that each
  maintainer must authorize; the user guide shows the opt-in.
- No new key. The release that carries this Spec carries Spec 0374 too.
- ADR-0385 records the decision. Four implementation Tasks in series and a QA
  gate.

## Open Questions

None.
