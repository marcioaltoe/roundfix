# F6 sender diagnosis

suites: 10
recurrences: 0
sender: not reproduced
diagnostic report: not inspected; no suite reproduced F6

Base commit: `e5da154a7bf5d991d0b796f187e1abeb25c00312`. Task 01's fixture process IDs,
probe report and F6 failure-line changes were applied. Both delivery fixture
builds still copied in place; Task 02's replacement fix was not applied.

Each suite ran `go test -count=1 -parallel 16 ./...` with a separate temporary
home. Raw output and diagnosis metadata stayed outside the checkout. Paths
below are relative to the temporary root. The working-tree diff and untracked
source fingerprints matched before and after every suite; the driver wrote
nothing under the checkout while a suite ran.

| Suite | Start | Wall seconds | Exit status | Signal-related lines | F6 recurrences | Raw output |
| --- | --- | ---: | ---: | ---: | ---: | --- |
| 1 | 2026-10-17 09:14:35 BRT | 251.960 | 1 | 0 | 0 | `suite-01/raw.log` |
| 2 | 2026-10-17 09:18:47 BRT | 250.082 | 1 | 0 | 0 | `suite-02/raw.log` |
| 3 | 2026-10-17 09:22:57 BRT | 282.891 | 1 | 0 | 0 | `suite-03/raw.log` |
| 4 | 2026-10-17 09:27:40 BRT | 257.541 | 1 | 0 | 0 | `suite-04/raw.log` |
| 5 | 2026-10-17 09:31:58 BRT | 270.069 | 1 | 0 | 0 | `suite-05/raw.log` |
| 6 | 2026-10-17 09:36:28 BRT | 262.198 | 1 | 0 | 0 | `suite-06/raw.log` |
| 7 | 2026-10-17 09:40:50 BRT | 245.803 | 1 | 0 | 0 | `suite-07/raw.log` |
| 8 | 2026-10-17 09:44:56 BRT | 265.084 | 1 | 0 | 0 | `suite-08/raw.log` |
| 9 | 2026-10-17 09:49:21 BRT | 260.623 | 1 | 0 | 0 | `suite-09/raw.log` |
| 10 | 2026-10-17 09:53:42 BRT | 279.702 | 1 | 0 | 0 | `suite-10/raw.log` |

Signal-related lines count output lines matching `was killed by signal`,
`signal: killed` or `exited -1`. No failed suite was rerun to obtain a pass;
the diagnosis counted each unrelated failure and stopped at the ten-suite cap.

Every suite exited 1 in `TestGuidanceCompositionJourney/rust-cli`: the external
`rustfmt src/main.rs` command failed because rustup had no default toolchain
configured under the temporary home. This was the only failing test in each
suite. These are diagnosis attempts with an environment limitation, not ten
passing acceptance suites. F6 never produced a failure line or probe report,
so this run attributes no signal to a sender and did not read macOS diagnostic
reports. The raw logs remain in temporary storage, not in Git.
