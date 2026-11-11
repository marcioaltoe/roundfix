---
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
status: active
created: 2026-11-10
surfaces: [infra, docs]
type: fix
---

# A push whose tests pass runs every contract under a measured budget

The `CI | Verify` gate of a push to `main` fails most of the time although every test passes. Its `Verify` step fails the suite-time budget of 240 s, which the workflow sets below the Makefile's own 360 s, and the steps after it, `Verify docs` and `Verify every contract`, are skipped, so the Full Contract Run that the repository rules place on every push to `main` did not run for those commits. Of the 33 pushes from 2026-11-06 to 2026-11-10 that passed every test, 26 exceeded 240 s, at 242 s to 343 s, because GitHub assigns runners of different speeds. After this Spec the budget is one measured number, 420 s, kept in the Makefile, and it is judged in a last step after the docs and contract steps, so a push whose tests pass gets both steps on every runner, while a suite that exceeds the budget still fails the gate.

## Sources

Two sources, both in `references/` (see `references/_index.md`):

1. [`references/2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md`](references/2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md),
   the open Backlog Entry of 2026-11-10, recorded by the operator, adopted
   under the adoption contract (`status: promoted`).
2. [`references/2026-11-10-ci-verify-push-suite-times.md`](references/2026-11-10-ci-verify-push-suite-times.md),
   the author's read-only measurement of the 34 completed push runs of the
   workflow on `main` from 2026-11-06 to 2026-11-10.

The other open Backlog Entries at the authoring base (a stale user-level skill
copy, the Catalog Skill Bump digest and a Portuguese marker in ADR
Alternatives) change other components and have their own Specs, 0376, 0378
and 0379; `docs/findings/` holds no entry.

## Prerequisites

None in `requires`. Authored against `60ac8142` and rebased on `b233761e`,
which carries the 0373 and 0374 plans. The queue delivers 0373, 0374 and 0375
before this Spec, then 0376, 0378 and 0379. This Spec owns
`.github/workflows/ci-verify.yml` and the `test-budget` lines of the
`Makefile`. It shares `docs/agents/setup-context.json` with Spec 0373, whose
task_02 regenerates the manifest after a Baseline change and merges first,
and the `Makefile` and the manifest with Spec 0378, which edits another
target after this Spec merges. task_01 regenerates the manifest with the
repository's command, never by hand, as its last step, so whatever an earlier
Spec merged into the manifest only moves the bytes that command writes. It
shares no file with 0374 or 0375 at their planned scope.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted or changed; the verdict file is a line of text in the
  runner's temporary directory. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no command reaches the network;
  the workflow keeps its permissions and its checkout. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0387 (new) records the decision
  and amends ADR-0259, whose last consequence this Spec closes. ADR-0259:
  "The CI suite budget stays at 240 s until CI measures the converted
  suite." ADR-0253 places the Full Contract Run on every push to main, which
  the change restores. ADR-0253: "CI runs the Full Contract Run on every push
  to main and in the release workflow before any build or publication."
  ADR-0089, ADR-0252 and ADR-0271 are not applicable: ADR-0259 and ADR-0253
  cite them or they cite ADR-0259, and the environment, contract-relevance
  and real-time rules they decide are unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization covers the
  CI workflow for this Spec ("CI na 0377") and, in the operator's cycle BB
  batch 2 brief, the `Makefile` when the design needs it; the design needs it
  because only `make test-budget` knows whether its suite failed or only took
  too long. `docs/agents/setup-context.json` is the derived manifest the
  repository's command regenerates after a `Makefile` change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget/_authorization.md`;
  bounded files: `.github/workflows/ci-verify.yml`, `Makefile`.

## Goals

- A push to `main` whose tests pass runs `Verify docs` and the Full Contract
  Run whatever its suite time.
- The suite-time budget is one measured number in the `Makefile`, and a push
  whose tests pass on any runner measured is green under it.
- A suite that exceeds the budget still fails the gate, after the contract
  runs, with the same verdict line as today.

## User Stories

1. As the maintainer, I want every push to `main` whose tests pass to run the
   Full Contract Run, so that a contract the selective gate left out is still
   checked before the next release.
2. As the maintainer, I want a red push gate to mean a real defect, so that I
   can trust it instead of rereading logs to find an over-budget suite.
3. As the maintainer, I want a large wall-clock regression to fail the gate,
   so that a change that serializes the suite again is caught on main.

## Core Features

1. `make test-budget` records its verdict line in the file
   `SUITE_BUDGET_VERDICT` names and does not fail on an overrun when that
   variable is set; a failing suite records no verdict and removes a stale
   one; without the variable it fails at once, as today.
2. `make suite-budget-check` prints the recorded verdict and fails when it says
   the budget was exceeded, when no verdict was recorded and when the verdict
   cannot be read.
3. The push path of `.github/workflows/ci-verify.yml` records the verdict in
   the runner's temporary directory during `Verify` and judges it in a last
   step, `Check the suite-time budget`, which runs whenever `Verify` passed,
   even when the docs or contract steps failed; the workflow no longer sets
   `SUITE_BUDGET_SECONDS`.
4. `SUITE_BUDGET_SECONDS` defaults to 420, with a comment naming the
   measurement and ADR-0387; ADR-0259 carries the "Amended by" notice.

## User Experience

Before this Spec, on a slow runner, the push log ends the `Verify` step with
`suite-time budget exceeded: the suite took 320s against a budget of 240s` and
`the budget is SUITE_BUDGET_SECONDS in the Makefile; raise it deliberately or
find the regression`, and `Verify docs` and `Verify every contract` show as
skipped. After it, the `Verify` step prints `suite-time budget: 320s of 420s`
and `suite-time verdict recorded in <runner temp>/suite-budget-verdict; make
suite-budget-check judges it`, the docs and contract steps run, and the last
step prints `suite-time budget: 320s of 420s` and passes. A suite of 450 s
prints `suite-time budget exceeded: the suite took 450s against a budget of
420s` in both steps, and the last step fails with the same hint line as today.

## Non-Goals / Out of Scope

- Making the suite faster, or splitting `internal/cli`.
- A calibrated or per-runner budget, a CPU-model report, or a second job for
  the Full Contract Run (ADR-0387 records why).
- Running `Verify docs` or the Full Contract Run after a failing test: a
  failing test still stops the push at `Verify`, as today.
- The pull request path (`Verify changed`, `Verify docs` with the skipped
  documentation contracts) and `.github/workflows/release.yml`, which never
  run `test-budget`.
- The flaky test that failed run 42261554097
  (`TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow`), which
  is not a budget failure.
- Any Go production code, guide, skill, `CONTEXT.md`, `go.mod`, `go.sum`,
  lint or formatter configuration.

## Success Metrics

1. Success Metric: `make test-budget` with `SUITE_BUDGET_VERDICT` set and a
   budget below the suite time exits 0 and records the exceeded verdict;
   `make suite-budget-check` then exits non-zero with that line and the hint.
   On the authoring base the first command exits non-zero and records nothing.
2. Success Metric: the parsed workflow has exactly one timed step, `Verify`
   with id `verify` and the verdict file, sets `SUITE_BUDGET_SECONDS` nowhere,
   and ends with `Check the suite-time budget`, which runs
   `make suite-budget-check` under
   `${{ !cancelled() && github.event_name == 'push' && steps.verify.outcome == 'success' }}`
   after `Verify every contract`; `actionlint`, where installed, reports
   nothing.
3. Success Metric: the default budget is 420 s, which every one of the 33
   measured suite times (at most 343 s) is under.
4. Success Metric: without `SUITE_BUDGET_VERDICT`, an overrun still fails
   `make test-budget` at once with the same two lines as before.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry, recorded by the operator on 2026-11-10 from
  runs this Spec did not design: "`CI | Verify` failed on 9 of the 12
  completed pushes to `main` between 2026-11-09 15:52 UTC and 2026-11-10
  21:04 UTC" and "Run 42747838420 on `5be18746` took 215 seconds and run
  42749378382 on `9fda0b65` took 311".
- Andrey Akinshin, "Performance stability of GitHub Actions" (2023-03-21,
  https://aakinshin.net/posts/github-actions-perf-stability/), read through
  Exa on 2026-11-10: "CPU-bound benchmarks are much more stable than
  Memory/Disk-bound benchmarks, but the “average” performance levels still can
  be up to three times different across builds."
- "Rust Benchmarking: How to Detect Performance Regressions"
  (https://hotpath.rs/blog/rust-performance-benchmarking), read through Exa on
  2026-11-10: "100 runs landed on 6 different CPU models."
- GitHub Docs, "Contexts reference"
  (https://docs.github.com/en/actions/reference/workflows-and-actions/contexts),
  read on 2026-11-10: `jobs.<job_id>.steps.env` may use the `runner` context,
  `jobs.<job_id>.env` may not; `jobs.<job_id>.steps.if` may use `steps` and
  `cancelled`; `steps.<step_id>.outcome` is "The result of a completed step
  before continue-on-error is applied. Possible values are success, failure,
  cancelled, or skipped."; and `runner.temp` is "The path to a temporary
  directory on the runner. This directory is emptied at the beginning and end
  of each job."

## Unreachable Acceptance

- criterion: the first push to `main` after this Spec merges runs `Verify`,
  `Verify docs`, `Verify every contract` and `Check the suite-time budget`,
  and its gate is green when every test passes.
  reason: it needs a push to `main` and a GitHub-hosted run of the workflow,
  which no Verification or QA row may start; the QA gate runs before any Pull
  Request exists.
  satisfied-by: the operator's reading of the first push run of
  `CI | Verify` on `main` after the merge, with `gh run view <id> --json jobs`,
  recorded session evidence.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Reading at `60ac8142`: the push `Verify` step of
  `.github/workflows/ci-verify.yml` sets `VERIFY_TEST_TARGET: test-budget` and
  `SUITE_BUDGET_SECONDS: "240"`; the `test-budget` target of the `Makefile`
  defaults the budget to 360 and exits 1 on an overrun after printing the
  verdict and the hint; `TestMainAndReleaseRunEveryContract` and
  `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts`
  (`internal/verifyselect`) parse the workflow, and no test names
  `test-budget`; no ADR or Archive Record sets the budget other than ADR-0259.
- The measurement in `references/`: 33 suite times from 209 s to 343 s, and
  the later steps skipped on run 42901489270.
- A prototype of the design on a `git archive` of `60ac8142`: the two new tests
  fail on the authoring base (seven of eight subtests and the workflow test)
  and pass with the change; each sabotage listed in the TechSpec fails its
  named subtest; the whole `./internal/verifyselect` package passes;
  `actionlint` 1.7 reports nothing on the changed workflow; and
  `TestTheSetupManifestIsCurrent` fails after the `Makefile` edit until the
  repository's manifest command runs, whose second run changes nothing.
- Published sources (Exa and GitHub Docs): the passages quoted above.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix CI suite-time budget GitHub runner variance push to main full contract run" --all --files --min-score 0.3`,
  which returned this repository's mirrored ADR-0253, release runbook,
  `CONTEXT.md`, CHANGELOG and Spec 0247 and 0282 records; none sets or
  discusses the suite-time budget, and none contradicts this design.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0377 of 2026-11-10, the adopted Backlog Entry, the measurement in `references/` and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~20:55–21:00 BRT), verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among others,
  "Portão do CI (Recommended)" — option text: "O push no main falhou 9 de 12
  vezes só pelo orçamento de 240 s da suíte (os testes passam) e pula o Full
  Contract Run. Spec pequena no workflow de CI."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those and the evidence, with no further interview:

- Both a new number and a deferred verdict: the number alone leaves the Full
  Contract Run hostage to the next slow runner, and the deferral alone leaves
  the gate red on 26 of 33 passing pushes.
- 420 s, in the `Makefile`, with no override in the workflow: 22 % above the
  slowest passing suite measured, one number instead of 240 in the workflow
  and 360 in the `Makefile`.
- The verdict is recorded by `make test-budget` and judged by a new
  `make suite-budget-check`, because only the target that runs the suite
  knows whether it failed or only took too long. The operator's brief grants
  the `Makefile` for 0377 when the design needs it, and lessons 42 and 55
  bring `docs/agents/setup-context.json` with it.
- The check runs whenever `Verify` passed, even after a failing docs or
  contract step, so the budget verdict is always reported when it exists.
- ADR-0387 records the decision and amends ADR-0259; a non-QA Task writes the
  notice into ADR-0259.
- One implementation Task and the QA gate; `type: fix`, Release Ceiling
  `patch`. The release is the operator's, right after the merge.

## Open Questions

None.
