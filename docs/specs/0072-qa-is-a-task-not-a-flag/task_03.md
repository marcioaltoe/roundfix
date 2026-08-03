---
task: task_03
spec: 0072-qa-is-a-task-not-a-flag
status: pending
type: backend
complexity: medium
---

# Task 03: Delete the QA parameter from the Implement Command

## Overview

`--qa` stops existing. What runs is what the graph says. The flag is removed
from the Implement Command's flag set and help text; passing it becomes an
unknown-flag error whose remediation names the authored-gate contract. The
"all Tasks completed" no-op rule inverts its special case: a graph whose
gate node is still unsettled is not a no-op — the Run starts and executes
the gate, which is the behavior the flag used to force. Agent selection
profile categories derive the `qa` category from the graph, not from the
request.

## Requirements

1. MUST remove the `qa` flag from the Implement Command's flag definitions,
   its help text, and the interactive-input QA choice merge; the request
   struct loses its QA field.
2. MUST make `--qa` an unknown-flag failure whose message tells the author
   the gate is declared in the Spec's Task Graph.
3. MUST treat an all-completed graph with an unsettled gate node as
   runnable: the Run starts and the cycle ends with the gate.
4. MUST treat an all-completed graph with a settled gate, a declined
   declaration, or a legacy shape as the no-op it is today.
5. MUST derive `implementProfileCategories`' `qa` category from the graph's
   gate declaration.
6. MUST update the CLI tests that pass `--qa` to graph-driven setups,
   keeping their assertions; add the unknown-flag remediation test.

## Subtasks

- [ ] Remove the flag, the request field, and the interactive merge path.
- [ ] Add the unknown-flag remediation message and its test.
- [ ] Invert the no-op rule around the unsettled gate node.
- [ ] Derive the profile category from the graph.
- [ ] Migrate the flag-driven CLI tests.

## Acceptance Criteria

- [ ] `roundfix implement --spec <slug> --qa` fails as an unknown flag and
      the message names the Task Graph as where the gate lives.
- [ ] The help text contains no QA parameter.
- [ ] An all-completed graph with an unsettled gate starts a Run that ends
      in the gate; with a settled gate or no gate it is a no-op.
- [ ] The `qa` selection profile category resolves for a graph with a gate
      and does not resolve for one without.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and
      this task file.

## Verification

- `go test ./internal/cli -count=1 -run 'Implement.*QA|QA.*Implement|UnknownFlag' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix implement --help 2>&1 | grep -q -- "--qa"`
  — expected: exit 1; the parameter is gone from help. (A `grep -qv` here
  would pass vacuously on any line without the flag.)
- `go test ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Success Metric 2.
- `_techspec.md` → API Contracts; Integration Points (`internal/cli`).
