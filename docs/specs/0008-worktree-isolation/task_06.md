---
task: task_06
spec: 0008-worktree-isolation
status: pending
type: docs
complexity: low
---

# Task 06: Sync docs and the Roundfix skill with worktree isolation

## Overview

Document the shipped model: Runs execute in Run Worktrees, outcomes gained
Integration Pending, settle reads kept worktrees, and two config surfaces
appeared — in the canonical Roundfix skill and README, cross-checked
against the binary. Verifiable through the skills drift check inside the
full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: the Run Worktree
   lifecycle (where it lives, when it is kept, the printed path), the
   Integration Pending outcome with its exit code and the exact integration
   command shape, the demoted implement dirty-tree note, settle's
   worktree scope, and the `worktree.copy` / artifact-default changes;
   regenerate the embedded copy through the sync target.
2. MUST update the README's behavior notes: the one gate difference (cold
   worktrees lack untracked files; `worktree.copy` is the remedy) called
   out explicitly.
3. MUST cross-check every documented line shape against CLI test fixtures
   and the built binary.
4. MUST verify glossary coverage (Run Worktree, Run Branch, Integration
   Pending already exist) and call out any further gap instead of
   inventing language.

## Subtasks

- [ ] Skill updates + `make skills-sync`
- [ ] README worktree behavior notes
- [ ] Binary and fixture cross-check
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [ ] The documented Integration Pending outcome line and worktree-path
      report lines appear verbatim in CLI test fixtures.
- [ ] README carries the cold-worktree note with the config remedy.
- [ ] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–6. `_techspec.md` → API
Contracts, Risks, Build Order 6. ADR-0023, ADR-0024. Repo hard rule
(canonical skill ships with CLI behavior changes).
