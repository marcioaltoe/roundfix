# CI | Verify push runs on main, 2026-11-06 to 2026-11-10

Measured on 2026-11-10 by the Spec 0377 author, read-only, from the logs of
the `Verification gate` job of every completed push run of `CI | Verify` on
`main` that `gh run list --workflow "CI | Verify" --branch main --event push
--limit 40` listed between run 42227040380 and run 42905480439. Each log was
read with `gh run view <id> --log`; the suite time is the number in the line
`suite-time budget: <n>s of 240s` or `suite-time budget exceeded: the suite
took <n>s against a budget of 240s` that `make test-budget` prints, and the
`internal/cli` time is that package's `ok` line. Cancelled runs are left out.
Every run used the image `ubuntu-24.04`, version 20261101.360.1 or earlier,
on Azure.

| run | commit | created (UTC) | suite (s) | internal/cli (s) | verdict |
| --- | --- | --- | --- | --- | --- |
| 42227040380 | `9a285d77` | 2026-11-06T13:37:47Z | 231 | 191.107 | within the budget |
| 42249155780 | `c249c85d` | 2026-11-06T16:05:52Z | 306 | 255.023 | over the budget; every test passed |
| 42259665690 | `2aa3c360` | 2026-11-06T17:13:35Z | 258 | 207.569 | over the budget; every test passed |
| 42261554097 | `4856fec3` | 2026-11-06T17:25:49Z | - | - | test failure (TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow) |
| 42292208053 | `055df6f3` | 2026-11-06T20:53:56Z | 254 | 205.573 | over the budget; every test passed |
| 42293809705 | `c73b9f1a` | 2026-11-06T21:05:26Z | 263 | 214.443 | over the budget; every test passed |
| 42295657851 | `497224da` | 2026-11-06T21:18:42Z | 242 | 202.601 | over the budget; every test passed |
| 42297504894 | `cb3ed035` | 2026-11-06T21:32:04Z | 231 | 190.708 | within the budget |
| 42314923682 | `0c9edf9b` | 2026-11-06T23:54:40Z | 306 | 256.881 | over the budget; every test passed |
| 42316784104 | `530fb22b` | 2026-11-07T00:12:20Z | 315 | 263.750 | over the budget; every test passed |
| 42327170050 | `2dd887a2` | 2026-11-07T02:02:19Z | 250 | 215.275 | over the budget; every test passed |
| 42328191817 | `ba59c813` | 2026-11-07T02:14:02Z | 343 | 281.446 | over the budget; every test passed |
| 42367963162 | `2e59aebe` | 2026-11-07T11:23:11Z | 209 | 176.476 | within the budget |
| 42613899136 | `0d6b34aa` | 2026-11-09T12:03:52Z | 260 | 209.206 | over the budget; every test passed |
| 42614640398 | `af02cc71` | 2026-11-09T12:09:22Z | 305 | 255.437 | over the budget; every test passed |
| 42615868121 | `3c6d0203` | 2026-11-09T12:18:10Z | 320 | 266.820 | over the budget; every test passed |
| 42625481040 | `60a31037` | 2026-11-09T13:26:20Z | 306 | 252.653 | over the budget; every test passed |
| 42627234087 | `cfc08683` | 2026-11-09T13:38:20Z | 214 | 176.229 | within the budget |
| 42640144064 | `62bc941a` | 2026-11-09T14:59:54Z | 307 | 252.194 | over the budget; every test passed |
| 42641693593 | `a60928fa` | 2026-11-09T15:08:59Z | 311 | 257.763 | over the budget; every test passed |
| 42645144836 | `9ad01bbd` | 2026-11-09T15:29:15Z | 269 | 237.135 | over the budget; every test passed |
| 42648956079 | `3325e203` | 2026-11-09T15:52:04Z | 306 | 254.102 | over the budget; every test passed |
| 42650973877 | `520bf845` | 2026-11-09T16:04:09Z | 291 | 252.968 | over the budget; every test passed |
| 42676603832 | `13885902` | 2026-11-09T18:43:59Z | 215 | 180.566 | within the budget |
| 42679307494 | `162646a0` | 2026-11-09T19:01:15Z | 314 | 258.373 | over the budget; every test passed |
| 42700230482 | `15811f11` | 2026-11-09T21:12:06Z | 322 | 271.270 | over the budget; every test passed |
| 42702770607 | `aa1b53a2` | 2026-11-09T21:28:06Z | 217 | 182.326 | within the budget |
| 42724011143 | `9fd5fd7e` | 2026-11-09T23:53:08Z | 313 | 260.262 | over the budget; every test passed |
| 42726232433 | `4d04b3f0` | 2026-11-10T00:10:21Z | 312 | 257.883 | over the budget; every test passed |
| 42747838420 | `5be18746` | 2026-11-10T03:11:52Z | 215 | 178.019 | within the budget |
| 42749378382 | `9fda0b65` | 2026-11-10T03:25:43Z | 311 | 258.466 | over the budget; every test passed |
| 42885499119 | `4984bced` | 2026-11-10T19:27:48Z | 314 | 261.017 | over the budget; every test passed |
| 42901489270 | `c3a0c8c8` | 2026-11-10T21:04:52Z | 320 | 263.920 | over the budget; every test passed |
| 42905480439 | `1febafab` | 2026-11-10T21:29:07Z | 309 | 255.459 | over the budget; every test passed |

## Summary

- 34 completed runs: one failed a test, and 33 passed every test and
  printed a suite time.
- Of those 33, 7 finished within 240 s (209 s to 231 s) and 26 exceeded it
  (242 s to 343 s): 7 between 242 s and 269 s, and 19 between 291 s and
  343 s. When the `Verify` step failed, `Verify docs` and `Verify every
  contract` were skipped (run 42901489270, read with
  `gh run view 42901489270 --json jobs`).
- The two kinds of runner are visible on one tree: `9fda0b65` changes only
  the three version files of `5be18746`, and the suite took 215 s on run
  42747838420 and 311 s on run 42749378382. `internal/cli` alone went from
  178 s to 258 s.
- In the three passing runs read with `gh run view <id> --json jobs`
  (42747838420, 42702770607, 42676603832) the `Verify` step took 233 s to
  234 s, `Verify docs` 177 s to 192 s and `Verify every contract` 165 s to
  172 s.
