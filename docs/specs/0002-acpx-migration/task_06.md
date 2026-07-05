---
task: task_06
spec: 0002-acpx-migration
status: pending
type: docs
complexity: low
---

# Task 06: Update docs and the Roundfix skill for the acpx dependency

## Overview

Bring the user- and agent-facing docs in line with the new agent layer: the canonical Roundfix skill and the README name the acpx pin and the Node prerequisite, the latency recommendation is recorded, and handoff work-plan item 2 is marked done. Verifiable through the skills drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill (and regenerate the embedded copy through the sync target): the agent layer runs through acpx at the pinned version, Node is a prerequisite, the install command, and that Runs drive one Agent Session per Run — without changing any documented command, flag, output, or exit-code text (none changed).
2. MUST update the README (or the equivalent install/usage doc) with the acpx prerequisite and the recommendation to configure direct adapter binaries in acpx config for latency-sensitive setups instead of default npx launches.
3. MUST record handoff work-plan item 2 as done in the handoff document's work plan (a status note, not a rewrite — history stays).
4. MUST verify every term used comes from the glossary (Agent Session, ACP Runtime, Run, Work Item, Stop Request); call out any gap in the Result instead of inventing language.

## Subtasks

- [ ] Roundfix skill: acpx pin, Node prerequisite, Agent Session note; embedded copy regenerated
- [ ] README/install docs: prerequisite plus latency recommendation
- [ ] Handoff work-plan item 2 status note
- [ ] Glossary coverage pass

## Acceptance Criteria

- [ ] The canonical skill names the exact pinned version shipped by task_03's constant and the install command; the drift check passes inside the full gate.
- [ ] The README prerequisite section matches what Preflight Validation actually demands (same version, same command).
- [ ] The handoff document shows item 2 closed with a pointer to this spec.
- [ ] No new un-glossaried term appears in the updated text.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts validate.
- `make verify` — expected: full gate passes, including the skills drift check.

## References

`_prd.md` → Core Feature 1; Non-Goals (skill semantics unchanged beyond the dependency); User Experience. `_techspec.md` → Integration Points, Build Order 6, Risks (first-run latency). Repo hard rule (canonical skill ships with behavior changes). ADR-0017.
