---
spec: 0074-git-spawn-economy
date: 2026-08-03
build: 4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0
status: closed
verdict: fail
rows_blocked_environment: 3
rows_blocked_finding: 0
surfaces: [backend, test]
---

# QA report — Git spawn economy

## Scope and environment

Full authored terminal QA gate for `task_07`. The PRD contains no explicit
user-story section; this matrix treats the maintainer, a Roundfix operator,
and a Baseline maintainer as the actors for the four goals, four core
features, and four Non-Goals. Pull Request journeys are planned but the QA
prompt states that no Pull Request is open.

- Build: `4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0`
- Run Worktree branch: `roundfix/run-run_20260803T171129Z_3d4f7965e9921af4`
- Spec target branch: `ma/spawn-economy`
- Host: macOS 26.5.2 (`25F84`), arm64
- Toolchain: Go 1.26.5; Git 2.55.0
- Application entry point: a fresh `roundfix` binary built from this checkout
- Test data: disposable local Git repositories and the committed canonical
  Baseline sources; no network or production data
- Parity gaps: no Open Pull Request; Pull Request observation is planned as a
  typed environment block
- Evidence root:
  `qa/evidence/2026-08-03-git-spawn-economy/`

## Project Constraint audit

Pass. The PRD and TechSpec each account for identifier strategy,
authentication/HTTP, active ADR obligations, and tooling authority with an
applicability decision, reason, and operative `docs/agents/` source. Accepted
ADR-0089 and ADR-0090 govern the implementation; ADR-0081 is not triggered.

`task_07` is the sole authored `qa` node, is terminal, and depends on
`task_06`, which transitively covers every other node. Tasks 01 through 06 are
`completed`; the Daemon-owned QA Task remains `pending` during this gate.
`git diff-tree` over all six Task commits and the independent
`dbdad8a..HEAD` path list found no protected repository-tooling change, so the
Spec's tooling-authority `not applicable` classification is accurate. The
current delta contains only this Spec-local QA directory.

## Static gate

Pass. `rtk make verify` exited 0: 3,125 Go tests passed across 24
packages, the focused four-test skills gate passed, `roundfix skills check`
passed, and the current `roundfix` binary built successfully.

## Planned execution

### QA-01 — Task Graph and Project Constraints

- Actor and entry point: maintainer reading the Spec and Git ancestry.
- Surface: governance/docs.
- Steps: validate `qa: task_07`; prove `task_07` is terminal and transitively
  depends on every non-QA leaf; verify all six dependencies are completed;
  audit identifier, authentication/HTTP, ADR, and tooling-authority rows in
  both active artifacts; resolve each Task commit with `git diff-tree`.
- Expected observable: a valid authored gate, complete applicable/not
  applicable constraint reasoning with operative `docs/agents/` sources, and
  no protected-tooling path in the six Task commits or current delta.
- Independent confirmation: compare commit path evidence with
  `git diff dbdad8a..HEAD --name-only` and a clean worktree status.
- Persistence check: rerun graph/path inspection from the current HEAD.
- Evidence: this report and `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-02 — Repository Verification

- Actor and entry point: maintainer running the repository gate.
- Surface: test/build.
- Steps: run `rtk make verify` once, read the exit status and every reported
  failure.
- Expected observable: formatting, test, Repository Skill Set, and build
  checks exit 0.
- Independent confirmation: inspect the final gate summary and a fresh built
  binary at the current commit.
- Persistence check: the later focused public-flow commands run against that
  binary without modifying the checkout.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-03 — Reproducible spawn census baseline

- Actor and entry point: Baseline maintainer following `baseline/README.md`.
- Surface: docs/test.
- Steps: syntax-check the shim and attribution program; verify executability;
  run the parser on an empty log; reconcile the recorded baseline and after
  tables with their named revisions and committed procedure.
- Expected observable: both procedures are runnable, zero input yields zero
  buckets, and every published number names its source run.
- Independent confirmation: rerun the focused syntax/parser commands and
  inspect both committed reports.
- Persistence check: repeat the empty-input parser check.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-04 — Skills restore batches object reads

- Actor and entry point: Roundfix operator using
  `roundfix baseline skills restore` on a disposable repository.
- Surface: CLI/backend.
- Steps: preview an offline restore of a multi-file skill; capture the plan
  digest; apply the exact digest; rerun for no changes; record Git invocations
  through the committed PATH shim.
- Expected observable: preview is confirmation-gated, confirmed apply succeeds,
  the repeat is idempotent, outputs remain stable, and one restore scope opens
  one `git cat-file --batch` process instead of one process per file.
- Independent confirmation: inspect restored files/lock through the public
  filesystem boundary and compare the second command's JSON.
- Persistence check: reopen the disposable repository and rerun restore.
- Probes: omit confirmation (safety), repeat confirmation (idempotence), and
  supply invalid input (recoverable error vocabulary).
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md` and
  `evidence/2026-08-03-git-spawn-economy/git-invocations.tsv`.

### QA-05 — Assets sync batches provenance reads

- Actor and entry point: Baseline maintainer using
  `roundfix baseline assets sync --check` on committed canonical sources.
- Surface: CLI/backend.
- Steps: run read-only JSON check twice, record Git invocations, compare
  output bytes/digest-bearing fields, then exercise invalid source input.
- Expected observable: check is read-only and stable; one provenance scope
  opens one `git cat-file --batch`; unsafe/non-repository input remains a
  specific actionable error.
- Independent confirmation: compare checkout status before/after and compare
  two public JSON outputs.
- Persistence check: fresh second process returns the same state.
- Probes: repeated check, invalid source path, and source outside a Git
  checkout.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md` and
  `evidence/2026-08-03-git-spawn-economy/git-invocations.tsv`.

### QA-06 — Repository resolution combines fresh facts

- Actor and entry point: operator invoking the Baseline CLI against a
  disposable Git repository.
- Surface: CLI/backend.
- Steps: capture the public result and shim log; mutate the repository with a
  new commit; rerun; exercise a non-repository path.
- Expected observable: related `rev-parse` facts share one invocation,
  repository identity refreshes after mutation, and non-repository/missing
  HEAD errors remain distinguishable.
- Independent confirmation: compare the public result to direct Git reads.
- Persistence check: restart the command after the mutation rather than reuse
  process state.
- Probes: mutation between reads, missing HEAD, and non-repository input.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md` and
  `evidence/2026-08-03-git-spawn-economy/git-invocations.tsv`.

### QA-07 — Agent runner takes its environment explicitly

- Actor and entry point: operator running the public Doctor Command; test
  maintainer running the race gate.
- Surface: CLI/backend/test.
- Steps: run Doctor through a controlled ACPX executable/environment and
  observe the child-visible marker; repeat in a fresh process with a changed
  marker; run `go test ./internal/agent -race -count=2 -parallel 16`.
- Expected observable: each child receives its invocation-scoped base and
  overrides, process defaults resolve at the boundary, public diagnostics stay
  stable, and parallel tests report no race.
- Independent confirmation: the controlled executable records its received
  environment while Doctor reports the same probe outcome.
- Persistence check: a second process sees the new marker and no stale value.
- Probes: changed environment between processes, missing executable, and
  parallel repetition.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-08 — Fresh suite completes under 60 seconds

- Actor and entry point: maintainer following the committed timing procedure.
- Surface: test/performance.
- Steps: create an empty task-local `GOCACHE`; run
  `/usr/bin/time -p go test ./... -count=1 -parallel 16`; record exit and
  wall/user/system time.
- Expected observable: exit 0 and `real` strictly below 60 seconds.
- Independent confirmation: compare with the committed baseline and Task 06
  after-measurement from the same host class.
- Persistence check: fresh cache and disabled test-result cache prevent a
  cached success.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-09 — Observable behavior remains byte-stable

- Actor and entry point: operator running released public command shapes on
  base and current binaries.
- Surface: CLI/backend.
- Steps: build the pre-Spec revision and current revision from clean exported
  trees; run representative success and failure Baseline commands on matched
  disposable inputs; compare stdout, stderr, exit code, and resulting digest
  or file state.
- Expected observable: same public outputs, errors, exit codes, and digests
  for identical inputs.
- Independent confirmation: compare filesystem/digest state after command
  restart.
- Persistence check: rerun both binaries on fresh matched inputs.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-10 — Tests were not deleted, skipped, or weakened

- Actor and entry point: maintainer reviewing the base-to-head Git diff and
  fresh suite inventory.
- Surface: test.
- Steps: inspect deleted/changed test paths, search the delta for skip or
  weakened-selection changes, compare top-level test inventory, and correlate
  with `make verify`.
- Expected observable: no test deletion, new tests only strengthen the changed
  seams, and no skip/disable pattern was introduced.
- Independent confirmation: base/head `go test -list` inventory comparison.
- Persistence check: rerun inventory from exported base/current trees.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-11 — CLI contract did not change

- Actor and entry point: operator invoking base and current binaries.
- Surface: CLI.
- Steps: compare root, Baseline, skills-restore, and assets-sync help; compare
  representative invalid-command exit/stdout/stderr behavior.
- Expected observable: command names, flags, stdout/stderr placement, and exit
  codes are byte-identical.
- Independent confirmation: the base-to-head diff contains no CLI production
  file.
- Persistence check: repeat comparison in fresh processes.
- Probes: `--help`, unknown subcommand, and invalid format.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-12 — Batching stays bounded; runners stay package-local

- Actor and entry point: maintainer inspecting the current implementation and
  mutation probe results.
- Surface: backend/architecture.
- Steps: verify no cross-mutation repository cache exists, the two batch
  scopes close their child, and no shared Git client extraction landed.
- Expected observable: reads batch only inside immutable scopes; subsequent
  commands observe repository mutation; package-local error vocabularies and
  runners remain distinct.
- Independent confirmation: QA-04 through QA-06 public restart/mutation
  probes plus base-to-head symbol/path inspection.
- Persistence check: restart after mutation and compare direct Git state.
- Evidence: `evidence/2026-08-03-git-spawn-economy/command-evidence.md`.

### QA-13 — Pull Request review readiness on exact head

- Actor and entry point: maintainer/reviewer observing the Spec target Pull
  Request read-only.
- Surface: Pull Request/review.
- Steps: resolve the Open Pull Request for `ma/spawn-economy`, then inspect
  review decision, checks, unresolved threads, and exact-head ancestry.
- Expected observable: reviewer decision and checks apply to the exact Spec
  target head with no unresolved review work.
- Independent confirmation: compare Pull Request head SHA with local target
  branch ancestry.
- Persistence check: reload Pull Request state.
- Evidence: QA prompt environment fact and this report's blocked-row record.

### QA-14 — QA Report contract

- Actor and entry point: Daemon reading the authored gate result.
- Surface: docs/workflow.
- Steps: confirm this is the first collision-safe report for the date; close
  every row; validate frontmatter verdict and exact typed blocked counts;
  resolve every evidence path.
- Expected observable: a dated closed report with zero pending rows and a
  mechanically derived verdict.
- Independent confirmation: grep the report frontmatter and result table after
  closure.
- Persistence check: reread from disk after the final patch.
- Evidence: this report.

## Results

| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| - | --- | --- | --- | --- |
| QA-01 | Task Graph and Project Constraints | Maintainer — governance/docs | pass | `command-evidence.md` → QA-01; exact graph, constraints, six commit path sets, and current delta audited |
| QA-02 | Repository Verification | Maintainer — test/build | pass | `rtk make verify` exit 0; 3,125 + 4 tests, skill check, and build passed (`command-evidence.md` → QA-02) |
| QA-03 | Reproducible spawn census baseline | Baseline maintainer — docs/test | pass | Shim syntax/executable checks and empty-log parser passed; both published revisions resolve (`command-evidence.md` → QA-03) |
| QA-04 | Skills restore batches object reads | Operator — CLI/backend | pass | Public preview exit 3, confirmed apply exit 0, restart idempotence, invalid-input exit 2, and one `cat-file` per five-file scope (`command-evidence.md` → QA-04; `git-invocations.tsv`) |
| QA-05 | Assets sync batches provenance reads | Baseline maintainer — CLI/backend | pass | Three exact-source public checks exited 0 with stable JSON and clean source; stale/non-Git sources exited 2 actionably; one batch per multi-file scope confirmed (`command-evidence.md` → QA-05; `assets-sync-invocations.tsv`) |
| QA-06 | Repository resolution combines fresh facts | Operator — CLI/backend | pass | Public shim logged one combined root/format/HEAD `rev-parse`; fresh process reread after checkout mutation; missing-HEAD and non-repository errors stayed distinct (`command-evidence.md` → QA-06; `git-args-before.tsv`, `git-args-after.tsv`) |
| QA-07 | Agent runner takes its environment explicitly | Operator/test maintainer — CLI/backend/test | pass | Public Doctor children independently received alpha then beta environments; actionable machine-readiness diagnostics remained; exact race/count gate passed 530 executions (`command-evidence.md` → QA-07; `acpx-alpha.tsv`, `acpx-beta.tsv`) |
| QA-08 | Fresh suite completes under 60 seconds | Maintainer — test/performance | fail | Exact empty-cache command exited 0 but took `real 86.36`, 26.36s over target; F-001 (`command-evidence.md` → QA-08) |
| QA-09 | Observable behavior remains byte-stable | Operator — CLI/backend | blocked (environment: exported base build inherited sandboxed `~/Library/Caches/go-build`) | Equivalent evidence: full current characterization and live flows passed; CLI/user-guide source is byte-identical base-to-head. Full-access base build comparison remains (`command-evidence.md` → QA-09) |
| QA-10 | Tests were not deleted, skipped, or weakened | Maintainer — test | pass | No deleted test path or added skip/forbidden dependency; top-level inventory grew 1,515→1,529; fresh full/race gates passed (`command-evidence.md` → QA-10) |
| QA-11 | CLI contract did not change | Operator — CLI | blocked (environment: exported base build inherited sandboxed `~/Library/Caches/go-build`) | Equivalent evidence: no CLI/user-guide source diff; current four help surfaces exit 0; invalid format exits 2 with stable stream placement/action (`command-evidence.md` → QA-11) |
| QA-12 | Batching stays bounded; runners stay package-local | Maintainer — backend/architecture | pass | No cache/shared client symbol or package; readers close on success/error; seven focused error cases and fresh mutation read passed (`command-evidence.md` → QA-12) |
| QA-13 | Pull Request review readiness on exact head | Maintainer/reviewer — Pull Request | blocked (environment: no Open Pull Request) | Prompt confirms absence; local target→Run ancestry and current gate observed, but no reviewer/check/thread state exists (`command-evidence.md` → QA-13) |
| QA-14 | QA Report contract | Daemon — docs/workflow | pass | First collision-safe report for 2026-08-03; closed frontmatter, exact typed counts, zero pending/skipped rows, and resolving evidence paths verified |

## Findings

### F-001 — Fresh suite remains above the 60-second target

- Impact: Friction.
- Actor and journey step: maintainer, QA-08 timing step.
- Expected: `go test ./... -count=1 -parallel 16` with an empty task-local
  `GOCACHE` exits 0 in strictly less than 60 seconds.
- Actual: all packages passed, but `/usr/bin/time -p` reported `real 86.36`,
  `user 137.32`, and `sys 243.25`; the wall time missed by 26.36 seconds.
- Reproduction: from repository root, create an empty task-local `GOCACHE`
  and run `rtk proxy env GOCACHE=<empty> /usr/bin/time -p go test ./...
  -count=1 -parallel 16`.
- Independent evidence: Task 06's same-host after run reported 83.38 seconds;
  this fresh QA run reported 86.36 seconds. Current package floors were
  `internal/cli` 72.501s and `internal/baseline` 71.872s.
- Evidence: `command-evidence.md` → QA-08.
- Affected rows: QA-08 only. The behavior-preservation and batching rows
  remain independently runnable.

## Blocked and skipped

QA-09 and QA-11 are environment-blocked by the same first base-export build
attempt: Go inherited `~/Library/Caches/go-build`, which this sandbox cannot
write. The denial happened before compilation and was not retried. Equivalent
observed evidence is recorded in each row: the full current characterization,
live current public flows, and a byte-identical base-to-head CLI/user-guide
source comparison. A full-access rerun must build `dbdad8a` with a writable
fresh Go cache and perform the planned byte-for-byte base/current comparisons.

QA-13 is environment-blocked because no Pull Request is open. The prompt is
supervised evidence of that cause, and local ancestry proves the Run head
descends from `ma/spawn-economy`, but no equivalent reviewer/check/thread
decision exists. Integrate the Run commits, open the Spec target Pull Request,
obtain review on its exact head, and rerun this row read-only.

## Coverage

Planned: four PRD goals, four core features, four Non-Goals, the two declared
surfaces, the authored QA Task criteria, three high-risk public journeys, and
their safety, invalid-input, mutation, restart, idempotence, and concurrency
Executed: 14 of 14 rows. Passed: 10. Failed: 1. Environment-blocked: 3.
Finding-blocked: 0. Skipped: 0. Pending: 0. Findings: one Friction. All four
PRD goals, four core features, four Non-Goals, both declared surfaces, both
QA Task criteria, and the planned safety, invalid-input, mutation, restart,
idempotence, and concurrency probes have terminal coverage.

## Final verdict

Fail. Ten rows passed, one failed, and three are environment-blocked; there
are zero finding-blocked or skipped rows and one Friction finding. Reduce the
empty-cache fresh suite below 60 seconds before rerunning the gate; a
full-access base build and an Open Pull Request are also required to replace
the three environment blocks with direct evidence.
