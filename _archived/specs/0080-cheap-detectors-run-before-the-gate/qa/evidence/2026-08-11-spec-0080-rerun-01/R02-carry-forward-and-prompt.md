# R02 — carry-forward and real QA prompt

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

## Resolver evidence

`rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go test
./internal/speccheck -run '^(TestCarriable|TestMechanicalStageCarriable)'
-count=1 -v` exited 0.

The suite passed 14 independently named refusal cases and two predicate happy
paths. It refused failed, blocked, skipped, inputless, changed-path,
non-ancestor, changed-content, mixed elapsed-time, missing-snapshot,
missing-citation, uncovered-citation, glob-intersection, and expanded-path-set
changes. Real-Git resolver cases passed unchanged carry with establishing
report/head, changed-input refusal, non-ancestor refusal, and preservation of
the original citation across a later carry.

## Prompt builder evidence

`rtk env ... go test ./internal/agent -run
'^TestBuildQAPrompt(CarriesTheSpecContextBundle|CarriesThePreviousReportIdentity|StatesQAGateContract)$'
-count=1 -v` exited 0. The lower builder covers changed/no-changed paths,
present/absent previous report identity, and all three blocked counts.

The public TaskCycle integration remains red under F-003: its non-Git fixture
fails in the real mechanical stage before it can observe the prompt.

## Real prompt divergence

The Roundfix QA prompt that opened this run names the Spec, directories,
branches, user checkout, Pull Request fact, and QA contract. It does not contain
any of these current-code fields:

- `Spec Context Bundle:` or `Prior changed files:`;
- `Previous QA Report:` and `Previous QA Report head:` despite
  `qa/qa-report-2026-08-11.md` existing before this run;
- `Seeded QA Report:`; no current-build `-01` seed existed before this Agent
  created the required all-pending report.

This is the real Agent-facing surface, so lower passing builder tests do not
prove Story 4 or Feature 4 works in the running application. The evidence does
not establish whether the Daemon binary was stale or another assembly boundary
dropped the fields; that runtime cause remains unknown.

## Timing boundary

The assignment forbids the two disposable committed Runs needed to compare a
first gate with a corrective re-gate. The focused resolver proves attribution
to carried rows, but no authorized production timing comparison exists.
