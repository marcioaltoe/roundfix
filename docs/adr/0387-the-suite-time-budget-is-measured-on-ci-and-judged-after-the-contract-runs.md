---
status: accepted
created_at: 2026-11-10T23:00:00Z
updated_at: 2026-11-10T23:00:00Z
deprecated_at: null
superseded_by: null
---

# The suite-time budget is measured on CI and judged after the contract runs

ADR-0259 kept the CI suite budget at 240 s "until CI measures the converted
suite". CI has now measured it: the 33 pushes to main from 2026-11-06 to
2026-11-10 that passed every test took 209 s to 343 s, and 26 of them
exceeded 240 s. The spread is the runner, not the tree: commit `9fda0b65`
changes only the version files of `5be18746`, and the suite took 215 s on one
run and 311 s on the next. Published measurements agree that wall-clock time
on shared GitHub runners swings by up to three times between builds of one
revision, partly because each run lands on one of several CPU models. While
the step that times the suite failed, `Verify docs` and `Verify every
contract` were skipped, so most pushes to main lost the Full Contract Run that
ADR-0253 places on every push.

**The budget is one measured number, in the Makefile.** `SUITE_BUDGET_SECONDS`
defaults to 420 s and the CI workflow no longer overrides it. 420 s is 22 %
above the slowest passing run measured; on the fast runner, where that same
tree took 215 s, it fails a suite that slows by about 1.95 times. The Makefile comment
names the measurement; changing the number takes a new measurement.

**The verdict is recorded where the suite runs and judged last.** When
`SUITE_BUDGET_VERDICT` names a file, `make test-budget` writes its verdict line
there instead of failing on an overrun, and `make suite-budget-check` fails
when that line says the budget was exceeded, when it is missing and when it
cannot be read. On a push the `Verify` step records the verdict, `Verify docs`
and `Verify every contract` run as before, and the last step,
`Check the suite-time budget`, judges it whenever the `Verify` step passed,
even when the docs or contract steps failed. A failing test still fails the
`Verify` step and records no verdict. Without `SUITE_BUDGET_VERDICT`,
`make test-budget` fails at once on an overrun, with the same lines as before.

**Rejected.** Raising the number alone keeps the Full Contract Run hostage to
the next slow runner. A budget scaled by a calibration step of the same run
adds a second measurement with its own noise, and the packages slowed by
different factors on the slow runner (`internal/cli` 1.45 times,
`internal/daemon` 1.7 times), so no single factor fits; it would also need
real workflow runs to tune. Running the docs and contract steps before the
timed suite on a push needs no Makefile change, but it reports a failing test
only after more than five minutes of other steps and times the suite on a
build cache those steps warmed. A separate job for the Full Contract Run
spends a second runner, checkout and build on every push, and still leaves the
gate red on a slow runner. Parsing the `Verify` log to tell a budget overrun
from a test failure ties the workflow to the wording of a message.

## Consequences

- A push whose tests pass gets `Verify docs` and the Full Contract Run on every
  runner, and its gate is green unless the suite takes more than 420 s.
- A suite that slows by less than about 1.95 times on the fastest runner, or
  1.22 times on the slowest, is no longer caught. Smaller regressions are left
  to the per-package times every push log prints.
- The budget is judged after about six more minutes of CI on a push whose
  suite is over it, and the job still fails.
