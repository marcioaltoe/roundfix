---
task: task_04
spec: 0034-release-plan
status: pending
type: backend
complexity: high
---

# Task 04: Expose the read-only Release Plan CLI contracts

## Overview

Deliver the public `roundfix release plan` command over the completed domain and Git source. The slice includes parsing, help, deterministic text and JSON, exit mapping, decision guidance, the PRD acceptance matrix, and a command-level mutation audit.

## Requirements

1. MUST dispatch `roundfix release plan` through the repository's stdlib `flag.FlagSet` and `Run(args, stdout, stderr) int` conventions without introducing a CLI dependency.
2. MUST implement `--from`, `--to`, paired `--impact` and `--reason`, and `--format text|json` exactly as specified.
3. MUST render decision output only on stdout and diagnostics or corrective actions only on stderr; JSON output MUST be exactly one `roundfix.release-plan/v1` object.
4. MUST return exit `0` for `ready` and `no_release`, exit `3` for `approval_required` and `manual_classification_required`, and exit `2` with no partial result for invalid input or repository failures.
5. MUST print the exact approval question for minor, version-zero breaking, and major decisions and an explicit rerun shape for manual classification.
6. MUST add truthful root and command help without changing any existing command, flag, JSON, or exit contract.
7. MUST create no Run, read no Roundfix configuration, contact no external service, and mutate no release or repository state on any outcome.
8. MUST cover every PRD success-metric range, mixed commit order, manual classification, dirty state, and text/JSON stdout-stderr isolation with real CLI tests.

## Subtasks

- [ ] Add nested `release plan` dispatch and flag parsing.
- [ ] Connect the Git source and domain plan builder.
- [ ] Render deterministic text decisions and next actions.
- [ ] Render the versioned JSON schema without diagnostic leakage.
- [ ] Map all decision and invalid-input exit codes.
- [ ] Update root and command help contracts.
- [ ] Add command-level acceptance and read-only mutation tests.

## Acceptance Criteria

- [ ] The fix-only, compatible-feature, major-breaking, version-zero-breaking, mixed-order, no-release, and ambiguous scenarios match the PRD outcomes.
- [ ] Text output leads with the decision and prints only determining or blocking commits; JSON includes evidence for every commit.
- [ ] Approval-required and manual-classification plans remain valid stdout results with exit `3`.
- [ ] Invalid requests emit no partial text or JSON plan and return exit `2` with one actionable stderr diagnostic.
- [ ] Help describes every supported flag, default, state, and read-only boundary truthfully.
- [ ] Files, refs, tags, remotes, configuration, packages, and releases remain unchanged for every tested command outcome.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `cmd/roundfix/main.go`

## Verification

- `go test ./internal/cli -run 'TestReleasePlan(Command|Help|Text|JSON|ExitCodes|ReadOnly|DirtyTree)' -count=1` — expected: the public command, all decision states, invalid inputs, and mutation audit pass.
- `go test ./internal/releaseplan ./internal/cli -count=1` — expected: domain, Git integration, and CLI suites pass together.
- `go run -buildvcs=false ./cmd/roundfix release plan --help` — expected: concise help lists `--from`, `--to`, `--impact`, `--reason`, and `--format` and describes a read-only plan.

## References

- `_prd.md` → Goals 3-4; User Stories 2-4; Core Features 1 and 6-8; User Experience; Success Metrics.
- `_techspec.md` → API Contracts; exit codes; System Architecture: `internal/cli`; Coverage Map; Testing Approach; Build Order 4 and 6.
- ADR-0048 → Release planning is read-only and confirmation-gated.
