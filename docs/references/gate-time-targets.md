# Gate time targets

Approved on 2026-11-25. The JSON document beside this reference is the source
read by `gate-time`. No individual report fails for elapsed time.

| Gate | Machine | Target | Measurement and judgement |
| --- | --- | --- | --- |
| Task settlement | the mini | median 120 s, p95 240 s | 2026-11-25 field baseline: 607 s median, 882 s p95; first 30 passing checks after delivery |
| Task settlement lab | the mini | 45 s, 150 s with regeneration | 2026-11-25 lab: unchanged gate 17 s, changed gate 251 s, contracts 81 s; median of three runs after changes |
| Complete suite run | the mini | 360 s | 2026-11-25 whole-suite lab baseline: 376 s |
| Before pull request | the mini | 90 s | 2026-11-25 documentation gate lab baseline: 77 s |
| Pull request | standard macOS runner | median 480 s, p95 600 s | 2026-11-25 cold-cache run 45529968311; parallel-part estimates 350–380 s |
| Test binary | standard macOS runner | 300 s | 2026-11-25 CLI baseline: 564 s plus 69 s Baseline tests before sharing and sharding |

The approved targets are objectives; the baseline measurements above precede
the changes. The allowance is 25%, rounded up to ten seconds and capped by the
approved target. A measured result above its target gets no number.

## Lab protocol

Use the mini, its normal Go build cache, `GOMAXPROCS=5` and the same checkout. Record the
revision and machine. Run each case three times and report the median wall time:
`make verify-changed` on a Run branch with its existing Run Marker, then with a
change selecting regeneration contracts; `make verify-changed` on the candidate
branch after its archive commit; and `make verify-docs` before publication.
Do not set a Run Marker for the experiment. Record package-test invocations to
prove the settlement runs none and the candidate runs all packages fresh.
Keep the commands' exit status and distinguish failed checks from measurements.

## Windows

Both windows are open. They started on 2026-11-26 at 10:58:57 UTC, when Spec
0421's delivering pull request (#962) merged; no closing date is recorded yet.
The start is recorded in the JSON document. `make gate-times`, or
`go run -buildvcs=false ./cmd/gate-time windows`, reads both windows: the first
30 passing Task settlement checks from the Run Database and the first 20 passing
pull request runs of `CI | Verify` from GitHub. It pairs each verification start
with its passed command for the same Run, Task and attempt. Pull request time is
the longest job duration, excluding runner queue time. `-latest` refreshes the
available samples. Failed runs are not passing samples; a timeout misses its
target and must be investigated separately. The test-binary report names the
slowest binary for each part.

The rule is not recorded yet. After both windows close, the closing change records
the measured medians, p95s, numbers and rule. Before a release, read the windows;
a closed window over its recorded number holds the release until time returns
within it or the maintainer decides on the measured evidence.
