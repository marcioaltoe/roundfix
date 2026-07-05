---
task: task_03
spec: 0006-acpx-run-robustness
status: pending
type: infra
complexity: low
---

# Task 03: Mitigate the acpx message-buffer limit, or evidence it upstream

## Overview

acpx 0.12.0's 10 MiB per-message buffer killed two finished turns in one day
(`-32603 Message buffer exceeded 10485760 bytes`, `acpxCode: RUNTIME`).
Investigate the pinned version for any configuration surface that raises or
disables the limit; apply it if real, and produce the upstream-ready report
either way. Evidence-driven: no invented flags. Verifiable by the recorded
investigation and, when a mitigation exists, a proven configuration.

## Requirements

1. MUST inspect the pinned acpx 0.12.0 for a message-buffer configuration
   surface — CLI help, `acpx config show` schema, and the installed
   package's source (`~/.asdf`/npm global tree) for the `10485760` constant
   and any override reading — and record exactly what was found with file or
   help references.
2. MUST, if a real surface exists: document the recommended setting in the
   README/skill guidance for orchestrator use and prove it with a local
   reproduction (a command that previously exceeded the limit, or the
   constant's behavior test) — no roundfix code changes unless an invocation
   flag is the surface, in which case the runner gains it behind the
   existing invocation builders with rig tests.
3. MUST, if no surface exists: write the upstream issue draft (title, body
   with the two incident reproductions, the buffer constant location, and
   the orchestrator use case) into this task's Result, and add the limit to
   the shipped docs as a known constraint with the "large docs-task
   payloads" trigger description.
4. MUST update `docs/_inbox/dogfood-findings-2.md` item 1 with the
   investigation outcome (mitigated with X / upstream-only).

## Subtasks

- [ ] Source and help inspection of the pinned acpx for the buffer constant
- [ ] Mitigation applied and proven, or upstream draft written
- [ ] Docs guidance and findings update

## Acceptance Criteria

- [ ] The Result names the exact finding: the buffer constant's location and
      whether any override exists, with references.
- [ ] Mitigation path: the setting is documented and demonstrated; or
      upstream path: the issue draft is complete enough to file verbatim.
- [ ] Findings-2 item 1 carries the dated outcome.

## Verification

- `rtk go test ./...` — expected: full suite passes (unchanged unless an
  invocation flag was added, in which case the agent-package rig covers it).
- `make verify` — expected: full gate passes.

## References

`_prd.md` → User Story 3; Core Feature 3; Decisions (evidence-driven).
`_techspec.md` → Build Order 3, Risks. Round-2 dogfood finding 1 (both
incident run ids and the exact error line).
