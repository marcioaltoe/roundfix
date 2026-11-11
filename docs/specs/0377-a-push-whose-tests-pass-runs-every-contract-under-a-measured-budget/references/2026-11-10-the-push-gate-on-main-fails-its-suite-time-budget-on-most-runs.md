---
type: fix
status: promoted
created: 2026-11-10
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
reason: null
---

# The push gate on main fails its suite-time budget on most runs

## Symptom

`CI | Verify` failed on 9 of the 12 completed pushes to `main` between 2026-11-09 15:52 UTC and 2026-11-10 21:04 UTC. The three failures the operator read (`9fda0b65`, `4984bced` and `c3a0c8c8`) passed every test and failed only the suite-time budget: 311, 314 and 320 seconds against 240. The three successes recorded 215, 217 and 215 seconds of 240. When the `Verify` step fails, the steps after it are skipped, so `Verify docs` and `Verify every contract`, the Full Contract Run that the repository rules place on every push to `main`, did not run for those commits.

## Where

The `Verify` step of `.github/workflows/ci-verify.yml` on a push, which sets `VERIFY_TEST_TARGET: test-budget` and `SUITE_BUDGET_SECONDS: "240"`, and the `test-budget` target of the `Makefile`.

## Expected

A push to `main` whose tests pass gets a green gate and its Full Contract Run on either kind of runner, and a real wall-clock regression still fails. How the budget tells a slow runner from a slow tree is a choice for the Spec: a larger margin, a budget measured against a calibration step of the same run, or a Full Contract Run that does not depend on the budget step.

## Evidence

Commit `9fda0b65` changes only the three version files relative to `5be18746`. Run 42747838420 on `5be18746` took 215 seconds and run 42749378382 on `9fda0b65` took 311, and every package slowed by about the same factor: `internal/cli` 178 to 258 seconds, `internal/daemon` 62 to 107, `internal/spec` 106 to 149, `internal/store` 80 to 113. The difference is the runner, not the tree. Runs read with `gh run view <id> --log` and `gh run view <id> --json jobs`: 42901489270 (`c3a0c8c8`, 320 seconds, later steps skipped), 42885499119 (`4984bced`, 314 seconds), 42749378382, 42747838420, 42702770607 (217 seconds) and 42676603832 (215 seconds). The other six failures were listed by `gh run list --workflow "CI | Verify" --branch main --event push` and not read. Recorded by the operator on 2026-11-10.
