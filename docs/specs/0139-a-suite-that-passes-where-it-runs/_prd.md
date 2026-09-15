---
spec: 0139-a-suite-that-passes-where-it-runs
status: active
created: 2026-09-15
surfaces: [backend]
---

# A suite that passes where it runs

The repository Verification passes in CI and fails on a maintainer's macOS
machine at the same tree. Spec 0138's QA gate failed on it, and so does every
later Spec's gate on this machine, whatever the Spec itself changes. The same
`make verify`, run on the unchanged base, failed in four packages that no
Spec had touched. Each failure has a proven cause outside the code being
verified:

- **ACPX adapter tests.** The fixtures create and delete hard links to the
  running test binary for every test. macOS sometimes kills an exec that lands
  during that churn.
- **Task-cycle tests.** They wait a fixed two seconds for work that takes longer
  under load, so they fail although the Daemon returned the right result.
- **A historical audit subtest.** Its expectation went stale when the audit began
  reading a grant at its authorizing ancestor. It runs only where old Git objects
  survive locally, which is why CI never saw it.
- **Force-stop tests.** They cannot read the process table from inside the QA
  Agent's sandbox, where the gate runs the repository Verification today.

The outcome this Spec buys: the repository Verification means the same thing on
every machine that runs it, and a QA gate's result no longer depends on whether
its Agent chose to leave the sandbox.

## Project Constraints

- Identifier strategy: applicable — Verification event classifications and reasons, QA Report keys and Mechanical Refusal Codes keep their spelling and meaning; the QA gate's repository Verification publishes ordinary Verification events and coins no token. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, request or transport is created or changed; the network reach of the version freshness check in tests is recorded as a separate finding and left unchanged. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the repair touches Daemon-owned Verification, the QA gate's mechanical facts and the test contract of governed files. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs task verification and settles task status, so the QA gate's repository Verification is Daemon-run and the Agent only reads its result.
  ADR-0038 applies: the Daemon allows one Verification repair, and the QA gate's repository Verification adds no repair turn of its own.
  ADR-0056 applies: Spec Runs separate Task Capacity and Verification Capacity, so the QA gate's repository Verification acquires Verification Capacity like any other attempt.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and a repository Verification failure the Daemon observed stays a failure, never an environment block.
  ADR-0089 applies: code under test takes its environment explicitly, so no test may depend on local Git objects, machine load or churn of the running test binary.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, and the repository verification gate is one of the facts it names; this Spec runs it there without making any verdict more permissive.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the Daemon outside the Agent sandbox is the stage that can produce the repository Verification result.
  ADR-0130 applies: the audit judges governed paths, and history keeps the set honest, so the three governed paths this Spec repairs need its own grant.
  ADR-0057 applies: the Daemon exclusively owns Implement Task status, and the QA gate's repository Verification settles no Task status of its own.
  ADR-0091 applies: the QA gate is a Task node of its own type, and the repository Verification runs inside that node's step rather than as a separate command.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row compares the repository's own gate before and after this Spec on the machine where it failed.
  ADR-0020 is not applicable to this change: a parsed prompt result outranks the acpx exit code, and this Spec changes no Agent result handling; the ACPX repair touches only test fixtures.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule.
  ADR-0097 is not applicable to this change: a QA row carries forward only on declared, unmoved evidence, and this Spec changes no carry-forward rule.
  ADR-0127 is not applicable to this change: process residue is a readiness fact, not a Run record, and this Spec changes no residue inventory.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-15, recorded in `docs/specs/0139-a-suite-that-passes-where-it-runs/_authorization.md`; bounded files: `internal/daemon/task_engine.go`, `internal/daemon/task_engine_test.go`, `internal/speccheck/mechanical_test.go`. No sanctioned regeneration applies. The ACPX fixture file carries no historical grant and needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The repository Verification passes on a macOS maintainer machine at every tree
  where CI passes it. It has no failure that depends on unreachable Git objects,
  machine load, or exec timing of the running test binary.
- A QA gate's repository Verification result is the same whether or not the QA
  Agent runs inside a sandbox.
- A genuine hang in a Task-cycle test still fails. Slowness under load no longer
  does.
- No verdict, blocked-cause rule or Verification contract becomes more
  permissive.

## User Stories

1. As a Supervisor on a maintainer machine, I want `make verify` to report the
   same result CI reports for the same tree, so that a QA gate fails only on
   defects in the work.
2. As a Supervisor reading a QA Report, I want the repository Verification to
   have run outside the Agent's sandbox, with its command, verdict and log
   recorded, so that the static gate is a machine fact and not an Agent's
   judgement.
3. As a maintainer, I want tests to wait for the Daemon's result up to the test
   run's own deadline, so that load cannot fail a correct result and a real hang
   still fails with the test binary's timeout report.
4. As a maintainer, I want the historical authorization audit proven by a fixture
   the test controls, so that its result does not depend on which Git objects a
   clone happens to hold.

## Core Features

1. **ACPX fixtures stop churning links to the test binary.** No test creates and
   then deletes a hard link to the running test binary while other tests exec it.
   A link that must resolve to a real path under a package directory is created
   once before the package's tests run and removed after them.
2. **Task-cycle waits end at the test deadline.** Every wait in the Task-cycle
   test helpers derives its bound from the running test's deadline, falling back
   to a generous bound when none is set. No wait keeps a guessed fixed wall-clock
   duration. A genuine hang fails at the deadline with the test binary's timeout
   report.
3. **The historical audit is proven by a controlled fixture.**
   - The stale historical expectation is replaced by a repository the test builds.
     In it, a grant widened after its consuming commit is refused for the widened
     path, and the same change is authorized when the grant already bounded it at
     the authorizing ancestor.
   - The historical-commit skip guard requires that a commit is reachable from the
     audited history, not merely that its object exists.
   - The other caller of that guard keeps its meaning.
4. **The QA gate step runs the repository Verification in the Daemon.** This
   applies when a repository Verification command is configured and the
   mechanical stage does not withhold the Agent.
   - The Daemon runs that command once, before the Agent turn and outside the
     Agent's sandbox, through the existing Verification machinery and Verification
     Capacity.
   - The seeded QA Report records the command, verdict, exit status and a link to
     the captured log. The log is stored under the Spec's QA evidence.
   - The gate prompt tells the Agent that the repository Verification already
     ran, gives its verdict and evidence, and tells the Agent not to run it again.
5. **Unconfigured and withheld cases keep today's behavior.** Without a
   configured command, the Daemon runs nothing, and the prompt says the gate runs
   the repository Verification itself. When the mechanical stage withholds the
   Agent, the Daemon does not run it either.
6. **The Daemon settles nothing from that result.** A failed repository
   Verification remains a failure the gate records under its existing rules. The
   Daemon does not settle the QA Task from it, and it gets no Verification
   Feedback repair.

## Non-Goals / Out of Scope

- The version freshness check that operational command tests run against GitHub.
  That is recorded as a finding for a later Spec.
- The 31 lock-copy analyzer diagnostics in `internal/agent`, which Spec 0123
  owns.
- The qa-gate skill text. Spec 0138 owns that text and aligns its static gate
  wording once this lands.
- CI workflow changes, a pinned Go toolchain, or any Verification configuration
  change.
- Raising a timeout, adding a retry or skip, or removing `t.Parallel` as a
  repair.
- New Mechanical Refusal Codes, Verification classifications, or a stage that
  withholds the Agent on a failed repository Verification.
- Deciding which revision authorizes a Task commit from an earlier Run. That is
  recorded in the 2026-09-14 finding on the mechanical stage's commit range.

## Declared intentional breaks

- With a repository Verification configured, the QA Agent no longer runs it. It
  reads the Daemon's recorded result.
- The historical audit subtest that read unreachable commits from the local
  object store is removed. A controlled fixture replaces it.

## Regression locks

- Every existing assertion of the ACPX, Task-cycle and audit tests keeps
  passing. Only the stale historical expectation is replaced.
- A repository Verification failure still makes the gate's static gate row a
  `fail`.
- A mechanical stage that withholds the Agent still withholds it, and runs no
  repository Verification.

## Acceptance evidence

The outside-evidence row compares measurements this Spec did not design: the
repository's own gate and the package test runs, on this machine, before and
after the change.

- **Before.** At the unchanged base, `make verify` exits 2. It fails in
  `internal/agent` (two ACPX tests), `internal/daemon` (a Task-cycle wait) and
  `internal/speccheck` (the historical audit subtest). Run in isolation, the same
  ACPX and Task-cycle tests pass.
- **After.** At the assembled tree, the same gate exits 0, and the stressed
  package runs pass:
  - the ACPX package at `-count=5 -parallel 16`;
  - the Daemon package under concurrent load.
- **Fresh clone.** The new audit fixture runs, rather than skips, in a clone that
  holds no stale objects.

If the base cannot be reproduced, the row is recorded as blocked with its reason.

## Success Metrics

- `make verify` exits 0 on this machine at the delivered tree.
- Spec 0138's rerun QA gate records the repository Verification from the Daemon
  and does not fail on any of the four families.
- The next two Spec QA gates on this machine record no repository Verification
  failure that CI does not also report for the same tree.

## Decisions

- **Suite first.** On 2026-09-15 the maintainer chose to fix the suite before
  retrying Spec 0138's gate.
- **Verification in the Daemon.** On 2026-09-15 the maintainer chose that the
  Daemon runs the QA gate's repository Verification in the mechanical step,
  rather than deferring to Spec 0122 or giving the QA Agent full access.
- **Symlinks over hard links.** The ACPX fixtures use symlinks, and one-time
  links for the package-directory case. Per-test hard links reproduce the kill,
  and neither symlinks nor links that are never deleted did.
- **Deadline-bound waits.** Task-cycle waits bind to the test deadline instead
  of a larger constant. A constant is still a guess that heavier load exceeds.
- **Refuse late widening.** The historical audit keeps refusing a grant widened
  after its consuming commit, as Spec 0119 requires. Legacy prose sanctions stay
  unparsed.

## Research basis

**Secondbrain, consulted before authoring.**

- The index-first query
  `qmd query "flaky go tests parallel hard link exec killed macOS test timeout deadline" --all --files --min-score 0.3`
  returned `projects/roundfix/mirror/docs/history/specs/0103-a-suite-that-leaks-nothing/qa/qa-report-2026-08-15.md`
  and `projects/roundfix/mirror/docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`.
  These framed the defects as environment the tests take implicitly, and ruled
  out serializing the suite as the fix.
- `inbox/roundfix/_triaged/2026-08-08-atrito-medido-do-loop-autonomo-em-tres-specs.md`
  records gate reruns caused by contract and environment rather than product
  defects. It supports fixing the environment instead of re-running gates.

**Exa MCP.** Consultation was attempted, but no Exa MCP tool was available in
this session. No external source was read, and no external validation is
claimed. That includes macOS exec behavior, which is established here only by
local experiment.

**Local diagnosis.** Each cause was reproduced before it was accepted:

- **ACPX kill.** A standalone experiment, plus an A/B run of the package with and
  without symlinks.
- **Task-cycle slowness.** An instrumented copy whose waits recorded late
  results, plus goroutine dumps and a race-detector run.
- **Audit failure.** A first-parent bisect showing the audit starts reading
  grants at the authorizing ancestor in the change that also stales the test.

## Open Questions

None.
