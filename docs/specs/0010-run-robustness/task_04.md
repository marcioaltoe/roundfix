---
task: task_04
spec: 0010-run-robustness
status: completed
type: docs
complexity: low
---

# Task 04: Sync docs and the Roundfix skill with the robustness changes

## Overview

Document the shipped surface: the config deprecation policy and warning
shape, the session-close reaping reports, and Detached Runs with their
four-line report and follow/stop surfaces — in the canonical Roundfix
skill and README, cross-checked against the binary. Verifiable through the
skills drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: `--detach` on the three
   operational commands (report shape, attach/stop follow-up, console-log
   location, Preflight relay), the deprecation warning behavior, and the
   force-stop/sweep session-close reports; regenerate the embedded copy
   through the sync target.
2. MUST update the README: the config-migration promise (removed keys warn,
   never break) and the Detached Run usage pattern for scripts/CI.
3. MUST cross-check every documented line shape against CLI test fixtures
   and the built binary.
4. MUST verify glossary coverage — Detached Run exists; call out any
   further gap instead of inventing language.

## Subtasks

- [x] Skill updates + `make skills-sync`
- [x] README migration promise and detach usage
- [x] Fixture and binary cross-check
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [x] The four detach lines and the deprecation warning appear verbatim in
      CLI test fixtures.
- [x] README carries both notes.
- [x] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–3. `_techspec.md` → Build
Order 4. ADR-0027, ADR-0028. Repo hard rule (canonical skill ships with CLI
behavior changes).

## Result

Updated the canonical Roundfix skill and regenerated the embedded
`skills/roundfix` copy with `rtk make skills-sync`. The skill now documents the
`resolve` / `watch` / `implement --detach` report, follow/stop commands,
console-log location, Preflight Validation relay, config deprecation warning,
and force-stop plus implement-sweep session-close reports. Updated the README
with the config-migration promise and Detached Run usage for scripts or CI.

Evidence:

- Skill drift: `rtk diff -r .agents/skills/roundfix skills/roundfix` produced
  no output after `rtk make skills-sync`.
- Fixture cross-check: `rtk rg -n -e "Run detached: %s|Console log: %s|Follow:
  roundfix attach %s|Stop: roundfix stop %s|config: resolve.concurrent is
  deprecated and ignored; use worktree.concurrency" internal/cli
  internal/config` found the four detach lines in
  `internal/cli/detach.go` / `internal/cli/implement_test.go` and the
  deprecation warning in `internal/cli/cli_test.go` /
  `internal/config/config_test.go`.
- Binary cross-check: a temporary binary built with
  `rtk go build -buildvcs=false -o /private/tmp/roundfix-task04/roundfix
  ./cmd/roundfix`; `resolve --help`, `watch --help`, and `implement --help`
  each listed `--detach`; `rtk strings /private/tmp/roundfix-task04/roundfix |
  rtk rg ...` found the detach report formats and config warning text in the
  built binary.
- Glossary: `CONTEXT.md` already defines `Detached Run`, `Attach`, `Stop
  Command`, `Preflight Validation`, `Artifact Directory`, and `Agent Session`.
  `Console log` is used only as the literal shipped CLI report label.
- `rtk go run ./cmd/roundfix skills check` passed:
  `Roundfix skill check passed: roundfix`.
- `rtk make verify` passed: Go tests reported 734 passing tests across 17
  packages, the Roundfix skill check passed, and the build completed.

Follow-ups: none.
