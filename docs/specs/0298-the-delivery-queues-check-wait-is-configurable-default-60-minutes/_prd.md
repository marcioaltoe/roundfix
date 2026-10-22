---
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# The Delivery Queue's check wait is configurable, default 60 minutes

Minimal PRD for a bug fix, minted by `write-techspec`. The Delivery Queue
waits at most 30 minutes, a constant, for a Pull Request's current-head checks
and then parks the item `checks-timeout`. On 2026-10-21 this repository's CI
Verification gate needed 30 to 33 minutes, so an item parked while CI was
still running and waited for a manual retry after its checks turned green.
After this Spec the wait is the User Config key `delivery.check_timeout`, a Go
duration greater than zero and at most a day, 60 minutes by default, and the
park names the key and the value it waited. The TechSpec holds the design.

## Sources

One Backlog Entry, written by this Spec's planning change and adopted at
birth into `references/` (see `references/_index.md`):

1. [`references/2026-10-22-the-delivery-queue-parks-checks-timeout-while-ci-is-still-running.md`](references/2026-10-22-the-delivery-queue-parks-checks-timeout-while-ci-is-still-running.md),
   type `fix`. Spec 0292's item parked `checks-timeout` twice on 2026-10-21
   while its CI ran 30.4 and 32.4 minutes; the wait is
   `defaultCheckTimeout = 30 * time.Minute` in `internal/delivery/engine.go`,
   overridden only by tests; the CI job's `timeout-minutes` is now 60.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` held only the reconcile entry Spec 0300 adopts,
and `docs/findings/` held no entry.

## Prerequisites

None in `requires`. This Spec is authored against `34b4e246` (v0.107.0), with
no active Spec under `docs/specs/`. It delivers after Spec 0297 and before
Specs 0299 and 0300 of the same cycle. Spec 0299 also changes
`internal/delivery`, for an item whose Spec is already archived on its
branch; this Spec changes only the Check Wait in `Engine.createPullRequest`,
`Engine.checkCandidate` and the engine's default, and it names functions
rather than line numbers. Its documentation Task raises the Roundfix skill's
version over whatever the tree holds when it starts, so it rebases over the
other Specs' version bumps.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new config key is a YAML key, and the park text extends an
  existing blocker. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the queue reads the same forge checks as
  before, and every test uses fake clocks, fake forges and temporary homes.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0236: "A failed check is judged on
  a merge with the current default branch"; ADR-0304: "A queue item waits for
  a Run slot up to runs.slot_wait"; ADR-0311: "The Delivery Queue's check wait
  is a User Config duration, sixty minutes by default". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 through AskUserQuestion, authorization "the same
  rules as cycle AN", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- The operator sets the Delivery Queue's check wait in User Config without a
  release, and nothing else can set it.
- With nothing set, the wait is 60 minutes, the CI Verification gate's job
  timeout, instead of 30.
- A value that cannot bound a wait is refused when the configuration loads.
- Each wait still starts its own deadline, and a check re-run still restarts
  it once.
- A `checks-timeout` park says which key bounded the wait and for how long.
- The configuration guide, the `deliver` guide, the Roundfix `deliver`
  reference and the glossary describe the wait.

## Core Features

1. The key. `delivery.check_timeout` is a Go duration in User Config, default
   `60m`. `0`, a negative value or a value above `24h` is refused with
   `delivery.check_timeout must be greater than 0 and at most 24h`; an
   unparsable value is refused as an invalid duration.
2. User Config only. A `delivery.check_timeout` in Project Config is ignored
   with the warning `config: delivery.check_timeout in Project Config is
   ignored; set delivery.check_timeout in User Config`, as for
   `runs.slot_wait`.
3. One bound for every Check Wait. The Delivery Queue's engine receives the
   configured value and uses it for the wait in `checking` and for the wait of
   publication on an existing Pull Request that still reports an earlier
   candidate head. Each wait starts its own deadline, and a re-run of a failed
   check restarts it once, as today.
4. The park text. At the deadline the item parks with the blocker
   `checks-timeout: current-head checks not settled within
   delivery.check_timeout <duration>`, the duration in Go's form such as
   `1h0m0s`; its Park Class stays `environment`.
5. Documentation. `docs/user-guide/configuration.md`, the `deliver` guide and
   the Roundfix `deliver` reference state the key, its default, its bounds,
   its scope and the park text, and the glossary defines **Check Wait**.

## Non-Goals / Out of Scope

- The check poll interval (15 seconds), the merge policy refusal, the stale
  and outside-change re-runs and their single timeout restart.
- The queue deadline, the retry limit, `runs.slot_wait` and the Run Slot Wait.
- Reading a workflow's `timeout-minutes` or any CI configuration.
- Resuming an item whose Spec is already archived on its branch, which Spec
  0299 owns.
- Rewriting a `checks-timeout` blocker recorded before this change.

## Success Metrics

1. Success Metric: with no `delivery.check_timeout` set, a Check Wait whose
   checks stay pending parks after exactly 60 minutes of fake-clock time,
   where the Delivery Base parks after 30, with the blocker `checks-timeout:
   current-head checks not settled within delivery.check_timeout 1h0m0s` in
   the `environment` Park Class.
2. Success Metric: a User Config `delivery.check_timeout: 45m` loads as 45
   minutes, the engine the `deliver` command builds receives it, a pending
   Check Wait parks after exactly 45 minutes with `... delivery.check_timeout
   45m0s`, and the wait for an existing Pull Request's head ends after
   exactly 45 minutes.
3. Success Metric: a Project Config `delivery.check_timeout` of `0`, `5m`,
   `-1h` or `invalid` leaves the User Config value in force and writes exactly
   the ignored-setting warning once.
4. Success Metric: a User Config value of `0`, `0s`, `-1m` or `24h1s` is
   refused with `delivery.check_timeout must be greater than 0 and at most
   24h`, `soon` is refused as an invalid duration, and `24h` is accepted.
5. Success Metric: a second Check Wait that starts after a first one lasted
   40 of its 45 minutes still lasts 45 minutes, and the existing re-run tests
   still restart the wait once and only once.
6. Success Metric: the configuration guide, the `deliver` guide and both
   copies of the Roundfix `deliver` reference give the key, `60m`, the bounds,
   the User Config scope and the park text, and `CONTEXT.md` defines **Check
   Wait**.
7. Success Metric: after release, the first Delivery Queue item whose CI runs
   longer than 30 minutes but finishes within the configured wait reaches
   `merging` without a `checks-timeout` park.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- This repository's history: `1f06abf1`, `ci: raise the verification gate
  timeout to 30 minutes (#612)`, and `c7bc381d`, `ci: raise the verification
  gate timeout to 60 minutes (#620)`, which leaves the `verify` job of
  `.github/workflows/ci-verify.yml` at `timeout-minutes: 60`.
- GitHub's documentation of Actions limits
  (https://docs.github.com/en/actions/reference/limits): for all
  GitHub-hosted runners, "Each job in a workflow can run for up to 6 hours of
  execution time. If a job reaches this limit, the job is terminated and
  fails."; for self-hosted runners, "A job can be in the queue for 24 hours
  before it is automatically cancelled."

## Unreachable Acceptance

- criterion: Success Metric 7, after release, the first Delivery Queue item
  whose CI runs longer than 30 minutes but finishes within the configured wait
  reaches `merging` without a `checks-timeout` park.
  reason: only a real forge runs CI on a Pull Request, and no Verification or
  QA row may reach the network or run `roundfix deliver` against a forge.
  satisfied-by: the operator's `roundfix deliver status` line for that item
  after the first such delivery following this Spec's release.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `34b4e246` with the engine's fake clock: an engine
  built without `CheckTimeout`, whose checks stay pending, parked
  `checks-timeout` after exactly 30 minutes and 121 polls; a User Config or
  Project Config holding `delivery:\n  check_timeout: 60m` failed to load with
  `yaml: unmarshal errors`. `newCommandDeliveryEngine` passes no
  `CheckTimeout`, so every queue used the constant. The TechSpec records the
  readers of the blocker.
- The CI run durations of 2026-10-21, read with `gh run list --workflow
  ci-verify.yml` on 2026-10-22: Spec 0292's Pull Request runs took 30.4
  minutes (cancelled at the job's then 30-minute timeout) and 32.4 minutes
  (passed); the v0.103.0 release Pull Request took 29.7 minutes.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "delivery
  queue checks timeout configurable CI duration roundfix" --all --files
  --min-score 0.3`, which returned the Roundfix mirror's decision record on
  the profile proof setup deadline (a configurable duration "greater than 0
  and at most 10m", the bounded-duration precedent), its `configuration.md`
  and `deliver.md`, and ADR-0236 (the single timeout restart after a
  re-run); no prior decision
  on the check wait exists, and no Inbox Entry under `inbox/roundfix/`
  mentions it. No Inbox Entry was written, because the briefing forbids
  writing there.
- Exa: GitHub's workflow syntax documents `jobs.<job_id>.timeout-minutes` as
  "The maximum number of minutes to let a job run before GitHub
  automatically cancels it. Default: 360", and GitHub's Actions limits page
  gives the 6-hour hosted job limit and the 24-hour self-hosted queue limit
  quoted above. Challenging: a job's limit does not include its queue time,
  and a workflow can chain several jobs, so the queue's wait cannot be derived
  from one job's limit; a day bounds every hosted job plus a full queue.

## Glossary

- adds: **Check Wait**

## Decisions

- 2026-10-22 (about 00:55 BRT): the cycle AO scope, all four in order ("Z:
  portão e corretiva", "Espera de checks", "Retomar item arquivado", "Y:
  reconcile PID"), one Spec and one minor each; this Spec's design "Config
  com padrão 60m (Recomendado)" (option text "Nova chave de User Config (ex.:
  `delivery.check_timeout`) com padrão 60 min, igual ao timeout atual do job
  ci-verify. O operador ajusta sem release."); and the authorization "the same
  rules as cycle AN".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AO briefing for Spec 0298 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The key is `delivery.check_timeout`, the maintainer's suggestion, in User
  Config only. The existing `delivery.*` keys that Project Config sets
  describe the repository's build and are read from the fetched default
  branch; none is a wait, while the machine-wide waits and ceilings
  (`runs.slot_wait`, `runs.max_active`) are User Config only. No contract
  forbids the User Config choice.
- `0` is refused, as for the profile proof setup deadline: a zero wait parks
  every item at its first poll, and an unbounded wait contradicts a bounded
  queue.
- The upper bound is `24h`: it covers GitHub's 6-hour hosted job limit and
  its 24-hour self-hosted queue limit, and a longer wait needs a person.
- The same value bounds the wait for an existing Pull Request's head, which
  used the same constant; splitting them would add a key nobody asked for.
- The park text extends the existing `checks-timeout` blocker after `: `, so
  the Park Class and every prefix reader keep working.
- One new ADR, ADR-0311, records the key, the default change and the bounds;
  the glossary gains **Check Wait**.

## Open Questions

None.
