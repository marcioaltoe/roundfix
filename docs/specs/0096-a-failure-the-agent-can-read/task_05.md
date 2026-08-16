---
task: task_05
spec: 0096-a-failure-the-agent-can-read
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 05: Say what the run budget bounds where it is set

## Overview

A maintainer setting a maximum Run duration cannot tell from the setting whether
it bounds wall-clock time from the Run's start or is evaluated at Work Item
boundaries. The single measured overrun has no established cause, so the
behaviour stays and its contract is stated where it is configured.

## Requirements

1. MUST state, in the configuration template the tool renders, what the maximum
   Run duration bounds and when it is evaluated.
2. MUST leave the evaluation behaviour unchanged; this Task documents, it does not
   move the check.
3. MUST leave the rendered template valid input to the loader that reads it.

## Subtasks

- [ ] State the contract in the rendered template.
- [ ] Prove the template still loads.

## Acceptance Criteria

- [ ] The rendered configuration states what the budget bounds and when it is
      evaluated.
- [ ] The rendered template round-trips through the loader.
- [ ] No budget evaluation code changed.

## Verification

- `go test -count=1 ./internal/config -run 'TestRenderedConfig' -v > /tmp/0096-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestRenderedConfig' /tmp/0096-t05.log || { cat /tmp/0096-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing template tests; fails today if no such test exists, in which case this Task adds it.
- `grep -A 3 'max_run_duration' internal/config/config.go | grep -qiE 'bounds|evaluated|wall.clock|work item' || { echo 'the rendered template still does not say what the budget bounds'; grep -B 2 -A 3 'max_run_duration' internal/config/config.go; exit 1; }` — expected: exits 0, proving the contract is stated beside the setting. Fails today.
- `go build -buildvcs=false -o /tmp/0096-t05-roundfix ./cmd/roundfix && /tmp/0096-t05-roundfix profiles show --json > /dev/null 2>&1; git diff --name-only HEAD -- internal/daemon internal/rounds > /tmp/0096-t05-behaviour.txt; test ! -s /tmp/0096-t05-behaviour.txt || { echo 'budget evaluation code changed, which this Task forbids:'; cat /tmp/0096-t05-behaviour.txt; exit 1; }; grep -A 3 'max_run_duration' internal/config/config.go | grep -qiE 'bounds|evaluated' || { echo 'nothing changed at all'; exit 1; }` — expected: exits 0, proving the built command still runs, no evaluation code moved, and the statement landed. Fails today on the last clause.

## Context

- interface: `internal/config/config.go`

## References

`_techspec.md` → Build Order 5; System Architecture, the configuration surface.
`_prd.md` → Core Feature 6; Open Questions. ADR-0137.

## Result

The rendered YAML now states that `max_run_duration` bounds wall-clock time
from Run start, is evaluated before review Rounds and during Review Source
waits, and does not interrupt active Work Item resolution. A focused config
test asserts that contract against `DefaultConfigYAML`, writes the rendered
template as User Config, and loads it through `Load`.

Focused checks:

- Before the template edit,
  `rtk env GOCACHE=/tmp/roundfix-task05-gocache go test ./internal/config -run '^TestRenderedConfig$' -count=1`
  failed in `states_Run_Budget_contract` because the rendered template lacked
  the contract text. After the edit, the same focused check passed.
- `rtk env GOCACHE=/tmp/roundfix-task05-gocache go test ./internal/config -run '^Test(DefaultConfigYAMLVerificationCapacity|RenderedConfig)$' -count=1`
  passed, covering the contract assertion, the real loader round-trip, and the
  adjacent rendered-template behavior.
- `rtk git diff --check` passed.
- `rtk git diff --name-only HEAD -- internal/watch internal/daemon internal/rounds`
  produced no paths; no budget evaluation code changed.
- `rtk make verify-incremental GOCACHE=/tmp/roundfix-task05-gocache` first reached
  two `internal/cli` force-stop integration failures because the sandbox denied
  process-table inspection. The approved rerun with process-table access exited
  0; the Go suite, skill checks, and build passed.

Acceptance evidence:

- The rendered-configuration contract is asserted by
  `TestRenderedConfig/states_Run_Budget_contract` and the focused test passed.
- `TestRenderedConfig/round-trips_through_Load` writes the rendered YAML and
  loads it through the package loader; the focused test passed.
- The changed evaluation-path query produced no paths; the implementation diff
  is confined to the rendered config text and its config-package test.

Daemon Verification was not run in this Agent turn, as required by the
Daemon-assigned execution contract.
