---
task: task_06
spec: 0008-worktree-isolation
status: completed
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

- [x] Skill updates + `make skills-sync`
- [x] README worktree behavior notes
- [x] Binary and fixture cross-check
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [x] The documented Integration Pending outcome line and worktree-path
      report lines appear verbatim in CLI test fixtures.
- [x] README carries the cold-worktree note with the config remedy.
- [x] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–6. `_techspec.md` → API
Contracts, Risks, Build Order 6. ADR-0023, ADR-0024. Repo hard rule
(canonical skill ships with CLI behavior changes).

## Result

Status: completed.

Acceptance evidence:

- Skill text was updated in `.agents/skills/roundfix/SKILL.md` for the Run
  Worktree lifecycle, Integration Pending exit/command shape, the demoted
  implement dirty-tree note, settle's Run Worktree scope, `worktree.copy`, and
  the Roundfix Home Artifact Directory default. `rtk make skills-sync`
  regenerated `skills/roundfix/SKILL.md`.
- The documented Integration Pending and Run Worktree line shapes were
  cross-checked against fixtures: `internal/cli/implement_test.go` asserts
  `IntegrationPending: 1 completed, 0 failed, 0 skipped, 0 pending; integrate with ...`,
  `Integration command: git merge --ff-only ...`, the dirty-tree note, and
  `Run Worktree kept: ...`; `internal/cli/cli_test.go` and
  `internal/tui/tui_test.go` assert `Run Worktree: ...` report lines. The
  built binary's `settle --help` reports the kept Run Worktree scope.
- `README.md` now explains that a new Run Worktree starts from committed Git
  state, so untracked files are absent unless repository-relative paths are
  listed under `worktree.copy`; it also documents the Roundfix Home Artifact
  Directory default.
- Glossary coverage was verified in `CONTEXT.md`: Run Worktree, Run Branch,
  Integration Pending, Roundfix Home, and Artifact Directory already exist; no
  new glossary gap was introduced.

Verification:

- `rtk go run ./cmd/roundfix skills check` passed:
  `Roundfix skill check passed: roundfix`.
- `rtk make verify` passed: `go test ./...` reported 682 passed tests in 17
  packages, `skills check` passed, and the binary built successfully.
