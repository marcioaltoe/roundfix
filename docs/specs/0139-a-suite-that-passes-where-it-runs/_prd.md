---
spec: 0139-a-suite-that-passes-where-it-runs
status: active
created: 2026-09-15
surfaces: [backend]
---

# A suite that passes where it runs

The repository Verification passes in CI and fails on a maintainer's macOS machine
at the same tree. Spec 0138's QA gate failed on it, and so does every later Spec's
gate on this machine, whatever the Spec changes. On the unchanged base, the same
`make verify` failed in four packages no Spec had touched. Each failure has a
proven cause outside the code being verified:

- **ACPX adapter tests.** The fixtures create and delete hard links to the running
  test binary for every test. macOS sometimes kills an exec that lands during that
  churn.
- **Task-cycle tests.** They wait a fixed two seconds for work that takes longer
  under load. They fail although the Daemon returned the right result.
- **A historical audit subtest.** Its expectation went stale when the audit began
  reading a grant at its authorizing ancestor. It runs only where old Git objects
  survive locally, which is why CI never saw it.
- **Force-stop tests.** They cannot read the process table from inside the QA
  Agent's sandbox, which is where the gate runs the repository Verification today.

The outcome this Spec buys: the repository Verification means the same thing on
every machine that runs it, and a QA gate's result no longer depends on whether its
Agent chose to leave the sandbox.

## Project Constraints

- Identifier strategy: applicable — Verification event classifications and reasons, QA Report keys and Mechanical Refusal Codes keep their spelling and meaning; the QA gate's repository Verification publishes ordinary Verification events and a failed one reuses the existing Precondition Refusal record, so no token is coined. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, request or transport is created or changed; the network reach of the version freshness check in tests is recorded as a separate finding and left unchanged. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the repair touches Daemon-owned Verification, the QA gate's mechanical facts, Spec artifact placement and the test contract of a governed file. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs task verification and settles task status, so the QA gate's repository Verification is Daemon-run and the Agent only reads its result.
  ADR-0038 applies: the Daemon allows one Verification repair, and the QA gate's repository Verification adds no repair turn of its own.
  ADR-0056 applies: Spec Runs separate Task Capacity and Verification Capacity, so the QA gate's repository Verification acquires Verification Capacity like any other attempt.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and a repository Verification failure the Daemon observed stays a failure, never an environment block.
  ADR-0089 applies: code under test takes its environment explicitly, so no test may depend on local Git objects, machine load or churn of the running test binary.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, the repository verification gate is one of the facts it names, and a failed or unobserved result is a blocking fact that withholds the Agent Session.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the Daemon outside the Agent sandbox is the stage that can produce the repository Verification result.
  ADR-0130 applies: the audit judges governed paths, and history keeps the set honest, so only the one governed test file this Spec repairs needs its grant.
  ADR-0035 applies: Spec Root is configurable and external Spec artifacts stay uncommitted, so the retained Verification log follows the QA Report's placement and commit rule.
  ADR-0036 is not applicable to this change: review artifacts are committed in a separate docs commit, and this Spec changes no review Run or its artifact commit.
  ADR-0029 is not applicable to this change: review artifacts live with the Spec, and this Spec changes no review artifact location.
  ADR-0142 is not applicable to this change: head-bound Review Source Evidence decides the watch outcome, and this Spec changes no watch Run or Review Source evidence.
  ADR-0057 applies: the Daemon exclusively owns Implement Task status, and the QA gate's repository Verification settles the QA Task only through the existing refusal report.
  ADR-0091 applies: the QA gate is a Task node of its own type, and the repository Verification runs inside that node's step rather than as a separate command.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row replays the repository's own gate against Spec 0138's recorded failure.
  ADR-0020 is not applicable to this change: a parsed prompt result outranks the acpx exit code, and this Spec changes no Agent result handling; the ACPX repair touches only test fixtures.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule.
  ADR-0097 is not applicable to this change: a QA row carries forward only on declared, unmoved evidence, and this Spec changes no carry-forward rule.
  ADR-0127 is not applicable to this change: process residue is a readiness fact, not a Run record, and this Spec changes no residue inventory.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-15, recorded in `docs/specs/0139-a-suite-that-passes-where-it-runs/_authorization.md`; bounded files: `internal/speccheck/mechanical_test.go`. No sanctioned regeneration applies. The Daemon and ACPX files this Spec changes are ordinary source outside the governed set. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The repository Verification passes on a macOS maintainer machine at every tree
  where CI passes it. It has no failure that depends on unreachable Git objects,
  machine load or exec timing of the running test binary.
- A QA gate's repository Verification result is the same whether or not the QA
  Agent runs inside a sandbox.
- A genuine hang in a Task-cycle test still fails, and slowness under load no
  longer does.
- No verdict, blocked-cause rule or Verification contract becomes more permissive.

## User Stories

1. As a Supervisor on a maintainer machine, I want `make verify` to report the same
   result CI reports for the same tree, so that a QA gate fails only on defects in
   the work.
2. As a Supervisor reading a QA Report, I want the repository Verification to have
   run outside the Agent's sandbox, with its command, verdict and log recorded, so
   that the static gate is a machine fact and not an Agent's judgement.
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
   once, before the package's tests run, and removed after them.
2. **Task-cycle waits end at the test deadline.** Every wait in the Task-cycle test
   helpers derives its bound from the running test's deadline, falling back to a
   generous bound when none is set. No wait keeps a guessed fixed wall-clock
   duration. A genuine hang fails at the deadline with the test binary's timeout
   report.
3. **The historical audit is proven by a controlled fixture.**
   - The stale historical expectation is replaced by a repository the test builds.
   - In that repository, a grant widened after its consuming commit is refused for
     the widened path.
   - The same change is authorized when the grant already bounded it at the
     authorizing ancestor.
   - The historical-commit skip guard requires the commit to be reachable from the
     audited history, not merely that its object exists.
   - The guard's other caller keeps its meaning.
4. **The QA gate step runs the repository Verification in the Daemon.** When a
   command is configured and the mechanical stage has not withheld the Agent, the
   Daemon runs it once:
   - before the Agent turn and outside the Agent's sandbox;
   - through the existing Verification machinery and Verification Capacity;
   - with its log retained whether it passes or fails.
5. **A failed or unobserved result stops the gate before the Agent.**
   - A non-zero exit, or an outcome the runner could not observe, is recorded as a
     Precondition Refusal. The refusal names the command, its exit status or the
     unobserved cause, and the retained log.
   - The report's verdict is `fail` and the Agent Session is withheld, as ADR-0096
     requires for a blocking machine fact.
   - No retry and no Verification Feedback apply.
   - A stopping Run stops as it does today.
6. **A passing result reaches the Agent as a fact.**
   - The seeded QA Report records the command, verdict, exit status and a link to
     the retained log.
   - The gate prompt tells the Agent that the repository Verification already ran,
     gives its result and evidence, and tells the Agent not to run it again.
   - Without a configured command the Daemon runs nothing, and the prompt says the
     gate runs the repository Verification itself.
7. **The log follows the QA Report.** The retained log lives beside the QA Report
   under the Spec's QA evidence. Under an external Spec Root it stays uncommitted,
   exactly as the report does.

## Non-Goals / Out of Scope

- The version freshness check that operational command tests run against GitHub.
  It is recorded as a finding for a later Spec.
- The 31 lock-copy analyzer diagnostics in `internal/agent`, which Spec 0123 owns.
- The qa-gate skill text. Spec 0138 owns it and aligns its static gate wording once
  this lands.
- CI workflow changes, a pinned Go toolchain, or any Verification configuration
  change.
- Raising a timeout, adding a retry or skip, or removing `t.Parallel` as a repair.
- New Mechanical Refusal Codes, Verification classifications or report frontmatter
  keys.
- The seeded report's initial verdict for gates that reach the Agent.
- Deciding which revision authorizes a Task commit from an earlier Run. The
  2026-09-14 finding on the mechanical stage's commit range records that decision.

## Declared intentional breaks

- With a repository Verification configured, the QA Agent no longer runs it. It
  reads the Daemon's recorded result.
- A failed or unobserved repository Verification now ends the gate as a
  Precondition Refusal before the Agent turn, so no flow row is observed until it
  passes.
- The historical audit subtest that read unreachable commits from the local object
  store is removed and replaced by a controlled fixture.

## Regression locks

- Every existing assertion of the ACPX, Task-cycle, verifier and audit tests keeps
  passing. Only the stale historical expectation is replaced.
- A repository Verification failure still ends the gate with verdict `fail`.
- A Task's successful Verification output is still removed. Only the QA gate step
  asks the runner to retain it.
- A mechanical stage that withholds the Agent still withholds it, and runs no
  repository Verification.

## Acceptance evidence

The outside-evidence row replays artifacts and revisions this Spec did not write:

- **Spec 0138's QA Report of 2026-09-15.** Its finding F-001 recorded the
  repository Verification failing in `internal/agent`, `internal/cli`,
  `internal/daemon` and `internal/speccheck`. It is at commit `990bb65b` on the
  kept Run Branch `roundfix/run-run_20260915T105843Z_41c059d11301702f`.
- **Delivery target revision `7a15e970`.** On this machine, `make verify` exited 2
  there. It failed in `internal/agent`, `internal/daemon` and `internal/speccheck`,
  while `internal/cli` passed outside the sandbox.

The terminal QA does four things:

1. reads that report at its commit;
2. runs `make verify` at the base revision;
3. runs `make verify` at the assembled tree;
4. runs the new audit fixture in a fresh clone that holds no stale objects.

When the commit or revision no longer resolves, the row is recorded as blocked with
that reason.

## Success Metrics

- `make verify` exits 0 on this machine at the delivered tree.
- Spec 0138's rerun QA gate records the repository Verification from the Daemon and
  does not fail on any of the four families.
- The next two Spec QA gates on this machine record no repository Verification
  failure that CI does not also report for the same tree.

## Decisions

- **Suite before gate.** On 2026-09-15 the maintainer chose to fix the suite before
  retrying Spec 0138's gate.
- **Daemon runs the Verification.** On 2026-09-15 the maintainer chose that the
  Daemon runs the QA gate's repository Verification in its mechanical step, rather
  than deferring to Spec 0122 or giving the QA Agent full access.
- **Failure withholds the Agent.** A failed or unobserved result becomes a
  Precondition Refusal. ADR-0096's blocking-fact rule and the existing refusal
  contract then decide the verdict, rather than the Agent.
- **Symlinks, not hard links.** The ACPX fixtures use symlinks, and one-time links
  for the package-directory case. Per-test hard links reproduced the kill; symlinks
  and never-deleted links did not.
- **Deadline-bound waits.** Task-cycle waits bind to the test deadline rather than
  to a larger constant, because a constant is still a guess that heavier load
  exceeds.
- **Late widening stays refused.** The historical audit keeps refusing a grant
  widened after its consuming commit, as Spec 0119 requires. Legacy prose sanctions
  stay unparsed.
- **One governed path.** After review, the grant bounds only the governed test
  file. The Daemon files are ordinary source.

## Research basis

Secondbrain was consulted before authoring:

- The index-first query
  `qmd query "flaky go tests parallel hard link exec killed macOS test timeout deadline" --all --files --min-score 0.3`
  returned
  `projects/roundfix/mirror/docs/history/specs/0103-a-suite-that-leaks-nothing/qa/qa-report-2026-08-15.md`
  and
  `projects/roundfix/mirror/docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`.
  They framed the defects as environment the tests take implicitly, and ruled out
  serializing the suite as the fix.
- `inbox/roundfix/_triaged/2026-08-08-atrito-medido-do-loop-autonomo-em-tres-specs.md`
  records gate reruns caused by contract and environment rather than product
  defects. It supports fixing the environment instead of re-running gates.

Exa MCP: consultation was attempted, but no Exa MCP tool was available in this
session. No external source was read and no external validation is claimed. That
includes macOS exec behavior, which is established here only by local experiment.

Local diagnosis reproduced each cause before it was accepted:

- the ACPX kill, by a standalone experiment and an A/B run of the package with and
  without symlinks;
- the Task-cycle slowness, by an instrumented copy whose waits recorded late
  results, plus goroutine dumps and a race-detector run;
- the audit failure, by a first-parent bisect.

The pre-PR Codex review of this Spec found seven gaps, each confirmed in the code
before it was accepted:

- the ungoverned grant paths;
- the deleted success log;
- the settlement authority;
- the ADR-0096 withholding rule;
- the unobserved outcome;
- the external Spec Root;
- the unnamed evidence source.

## Open Questions

None.
