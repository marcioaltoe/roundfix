---
task: task_04
spec: 0068-spec-close-audit
status: pending
type: backend
complexity: medium
---

# Task 04: Expose the audit as `roundfix spec audit`

## Overview

The user surface: a read-only sibling of `spec check`, taking only a slug
because the operator should not need to know which Runs or branches the Spec
produced. Demoable the moment it lands — running it against this repository
prints a real audit.

## Requirements

1. MUST add `roundfix spec audit <slug> [--format text|json]`.
2. MUST exit `0` when nothing needs attention, `1` when residue or undelivered
   work is found, and `2` on usage error or an unknown slug.
3. MUST print each survivor with its kind, its evidence, and — for residue —
   the exact reclaim command.
4. MUST emit one `roundfix-specaudit/v1` object on stdout under
   `--format json`, with diagnostics on stderr.
5. MUST leave every existing command's behaviour unchanged, including
   `spec check`.
6. MUST appear in the usage block and the command list.

## Subtasks

- [ ] Add the command dispatch, flags, and help text.
- [ ] Implement both renderers and the exit-code contract.
- [ ] Add the usage and command-list entries.
- [ ] Add command tests for every exit code and both formats.

## Acceptance Criteria

- [ ] A clean fixture exits `0` and prints a clean report.
- [ ] A fixture with residue exits `1` and prints the kind, the evidence, and
      the reclaim command.
- [ ] A fixture with an undelivered artifact exits `1` and names the holding
      branch.
- [ ] `--format json` emits one parseable `roundfix-specaudit/v1` object.
- [ ] An unknown slug exits `2` and names it on stderr.
- [ ] `roundfix --help` lists the command.
- [ ] A test asserts `spec check` behaviour is unchanged.

## Context

- instruction: `docs/agents/cli.md`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1 -run 'SpecAudit' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix --help | grep -q "spec audit"`
  — expected: exit 0; the command is discoverable.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; the Skill is task_06's bounded scope.

## References

- `_prd.md` → Core Feature 6; Goals.
- `_techspec.md` → API Contracts; Build Order 4.
