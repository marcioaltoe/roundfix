# Task 03 corpus sweep evidence

Recorded: 2026-09-17, before the corpus sweep.
Boundary: the active Specs returned by `roundfix spec check` from this worktree.

## Prediction

| Finding | Predicted count | Predicted severity |
| --- | ---: | --- |
| `SC-METRIC-UNDECLARED` | 10 | `gap` |
| `SC-CONTRACT-UNDECLARED` | 10 | `gap` |
| New `error` findings caused by these detectors | 0 | `error` |

The prediction comes from ten pre-existing Specs in authoring: seven PRDs have no Success Metrics section, three express it as bullets or a table, and all ten TechSpecs express API Contracts as prose. Spec 0140 is excluded from those ten because it declares numbered Success Metrics and API Contracts under the new rule.

## Observation

The first sweep after recording the prediction ran:

`rtk go run ./cmd/roundfix spec check --format json`

It returned 11 JSON documents with exit status 0. The comparison matched:

| Finding | Predicted | Observed | Difference |
| --- | ---: | ---: | ---: |
| `SC-METRIC-UNDECLARED` | 10 gaps | 10 gaps | 0 |
| `SC-CONTRACT-UNDECLARED` | 10 gaps | 10 gaps | 0 |
| New `error` findings caused by these detectors | 0 | 0 | 0 |

Both gap codes named the same ten pre-existing Specs: 0120 through 0129. Spec 0140 emitted neither code. The sweep reported no other finding code, and it did not edit any Spec in authoring.
