---
task: task_04
spec: 0006-acpx-run-robustness
status: pending
type: docs
complexity: low
---

# Task 04: Sync docs and the Roundfix skill with the robustness changes

## Overview

Document the shipped behavior: the Settle Command joins the canonical
Roundfix skill's command surface, the ADR-0020 classification is stated
where the skill describes Batch failure semantics, and any buffer guidance
from task_03 lands in the shipped docs. Verifiable through the skills drift
check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: the `settle` command
   (flags, exit codes, the stage-everything commit contract, failed-only
   targets), and the ADR-0020 rule that a delivered prompt result with a
   dirty transport exit proceeds to verification with the anomaly journaled;
   regenerate the embedded copy through the sync target.
2. MUST cross-check every documented flag and line shape against the built
   binary's output.
3. MUST fold task_03's outcome into the README/skill guidance (mitigation
   setting or known-constraint note).
4. MUST verify every term against the glossary; call out gaps instead of
   inventing language.

## Subtasks

- [ ] Skill updates for settle and ADR-0020 semantics + `make skills-sync`
- [ ] Help-text and stdout-shape cross-check
- [ ] Buffer guidance fold-in
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes inside
      the full gate.
- [ ] Documented settle stdout lines appear verbatim in CLI test fixtures.
- [ ] No new un-glossaried term, or the gap is called out in the Result.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Features 1–3; User Experience. `_techspec.md` → Build Order
4. ADR-0020. Repo hard rule (canonical skill ships with CLI behavior
changes).
