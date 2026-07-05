---
task: task_05
spec: 0004-watch-merge-readiness
status: pending
type: docs
complexity: low
---

# Task 05: Sync docs and the Roundfix skill with the watch contract

## Overview

Document the shipped behavior: the merge-readiness Clean semantics, the
watch/resolve stdout report shape, and the `--no-agent-console` flag — in the
canonical Roundfix skill and command help, with the embedded copy regenerated.
Verifiable through the skills drift check inside the full gate.

## Requirements

1. MUST update the canonical Roundfix skill: watch's Clean now means the
   Review Source check succeeded on the final pushed commit (ADR-0019,
   including the `missing` note), the exact stdout report shapes for watch
   and resolve, and `--no-agent-console` on the operational commands;
   regenerate the embedded copy through the sync target.
2. MUST verify each documented flag and line shape against the built binary's
   actual output.
3. MUST verify every term against the glossary; call out gaps instead of
   inventing language (candidate gap to flag if felt: a term for
   merge-readiness).

## Subtasks

- [ ] Skill updates + `make skills-sync`
- [ ] Help-text cross-check against shipped output
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes inside
      the full gate.
- [ ] Every documented stdout line shape appears verbatim in a CLI test
      fixture.
- [ ] No new un-glossaried term, or the gap is called out in the Result.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Features 2–4; User Experience. `_techspec.md` → Build Order
5. ADR-0019. Repo hard rule (canonical skill ships with CLI behavior
changes).
