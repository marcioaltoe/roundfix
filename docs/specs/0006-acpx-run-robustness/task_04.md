---
task: task_04
spec: 0006-acpx-run-robustness
status: completed
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

- [x] Skill updates for settle and ADR-0020 semantics + `make skills-sync`
- [x] Help-text and stdout-shape cross-check
- [x] Buffer guidance fold-in
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly; drift check passes inside
      the full gate.
- [x] Documented settle stdout lines appear verbatim in CLI test fixtures.
- [x] No new un-glossaried term, or the gap is called out in the Result.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Features 1–3; User Experience. `_techspec.md` → Build Order
4. ADR-0020. Repo hard rule (canonical skill ships with CLI behavior
changes).

## Result

Updated the shipped docs for the robustness changes:

- The canonical Roundfix skill now documents the Settle Command, including
  `--spec`, `--task`, exit codes `0`, `1`, and `2`, failed-only targets,
  deterministic stdout, the all-current-worktree-changes commit contract, no
  Run, no Run Event Journal entries, and no push.
- The embedded `skills/roundfix` copy was regenerated with `rtk make
  skills-sync`.
- The Roundfix skill now states the ADR-0020 classification: a Batch with a
  delivered `session/prompt` result and later nonzero acpx exit journals the
  anomaly and proceeds to Daemon Verification; without that parsed result, the
  nonzero exit remains a Batch failure.
- README and the Roundfix skill now carry task_03's finding that acpx `0.12.0`
  has a hard 10 MiB queue-owner per-message buffer with no discovered CLI,
  config, or environment override; recovery guidance points to ADR-0020
  classification and the Settle Command rather than an invented mitigation
  setting.
- `CONTEXT.md` now defines `Settle Command` and `Verification`; the new docs
  use glossary terms for the command, Spec, Task, Task Graph, Preflight
  Validation, Active Run, Run, Run Event Journal, Daemon, Batch, and Work Item
  language. No glossary gap remains.

Cross-check evidence:

- `rtk go run ./cmd/roundfix settle --help` printed the shipped command shape:
  `roundfix settle --spec <slug> --task <task_id>`, options `--spec` and
  `--task`, and exit codes `0`, `1`, and `2`.
- `rtk go run ./cmd/roundfix --help` listed
  `roundfix settle --spec <slug> --task <task_id>` and described `settle` as
  `Verify and commit all current worktree changes for one failed Task`.
- `rtk rg -n "verify test -f done.txt — ok|verify test -f missing.txt —
  failed|task_01 stays failed — verification failed|settled task_01 completed
  —" internal/cli/settle_test.go .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` found the documented stdout examples in both skill
  copies and in `internal/cli/settle_test.go` (`expectedStdout` assertions).

Verification:

- `rtk go run ./cmd/roundfix skills check`: passed with
  `Roundfix skill check passed: roundfix`.
- `rtk make verify`: passed. The gate ran `rtk go test ./...` (`605 passed in
  16 packages`), `rtk go run -buildvcs=false ./cmd/roundfix skills check`
  (`Roundfix skill check passed: roundfix`), and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.
