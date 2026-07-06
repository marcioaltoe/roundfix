---
task: task_02
spec: 0011-storage-lifecycle
status: pending
type: backend
complexity: low
---

# Task 02: Opt-in per-Batch agent logs

## Overview

Stop writing per-Batch agent log files by default and gate them behind a config
key. The Run Event Journal already stores every raw agent payload, so the log
files are a redundant on-disk copy; this task makes them development-only while
leaving the journal and the Detached Run console log untouched.

## Requirements

1. MUST add a `logs.agent` config key (User and Project Config), default off,
   that enables per-Batch agent log-file writing.
2. MUST NOT write any per-Batch agent log file when the key is off; the Run
   Event Journal MUST still record every payload in both states.
3. MUST leave the Detached Run console log (ADR-0028) unconditional and
   unaffected by this key.
4. MUST keep every other config default and generated output byte-stable and
   route the deprecation/None case through the config warning path if a removed
   key is involved.

## Subtasks

- [ ] `logs.agent` config key with default-off semantics
- [ ] Guard the per-Batch log writer on the key
- [ ] Confirm the journal path is independent of the guard
- [ ] Tests: no log files with key off, files with key on, journal populated in both

## Acceptance Criteria

- [ ] A production Run with the default config writes zero per-Batch agent log files.
- [ ] The same Run with `logs.agent` enabled writes the log files.
- [ ] The Run Event Journal contains the agent payloads in both cases.
- [ ] A Detached Run's console log is written regardless of `logs.agent`.

## Verification

- `rtk go test ./internal/config/ ./internal/agent/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 4-5; Core Feature 2. `_techspec.md` → Opt-in agent
logs, Build Order 2. ADR-0030 (extends ADR-0008). Work-plan finding R1-8.
