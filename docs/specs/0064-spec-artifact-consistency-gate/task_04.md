---
task: task_04
spec: 0064-spec-artifact-consistency-gate
status: pending
type: backend
complexity: medium
---

# Task 04: Expose the check as `roundfix spec check`

## Overview

Give the checker its user surface: a read-only support command in the family
`doctor`, `archive`, and `release plan` already occupy. With no slug it checks
every active Spec; with slugs it checks those. It reports to stdout, sends
diagnostics to stderr, offers a machine-readable format, and carries the
exit-code contract that lets a gate depend on it.

Demoable on its own the moment it lands: running it against this repository
prints a real report.

## Requirements

1. MUST add `roundfix spec check [<slug> ...] [--format text|json] [--strict]`,
   defaulting to every active Spec when no slug is given.
2. MUST exit `0` when there are no findings or only gaps, `1` when at least one
   finding has severity `error`, and `2` on usage error or an unreadable Spec
   Root.
3. MUST promote gaps to errors under `--strict`, changing the exit code
   accordingly.
4. MUST print each finding with its code, severity, one-line summary, every
   location as `path:line`, and its fix.
5. MUST emit one `roundfix-speccheck/v1` JSON object per Spec under
   `--format json`, on stdout, with diagnostics on stderr.
6. MUST report an unknown slug as a usage error naming the slug, not as a
   clean result.
7. MUST leave every existing command's behavior unchanged, and MUST NOT add a
   consistency precondition to `implement`.
8. MUST appear in the usage block and the command list with the other support
   commands.

## Subtasks

- [ ] Add the command dispatch, its flag set, and its help text.
- [ ] Wire single-slug, multi-slug, and all-active-Specs selection.
- [ ] Implement both renderers and the exit-code contract.
- [ ] Add the usage and command-list entries.
- [ ] Add command tests over fixture Spec Roots for every exit code and both
      formats.

## Acceptance Criteria

- [ ] `spec check` against a fixture Spec Root with a clean Spec exits `0` and
      prints a clean report.
- [ ] `spec check` against a fixture Spec with one `error` exits `1` and prints
      the code, both locations, and the fix.
- [ ] `spec check` against a fixture Spec whose only finding is a gap exits `0`,
      and exits `1` under `--strict`.
- [ ] `spec check --format json` emits one parseable
      `roundfix-speccheck/v1` object per Spec on stdout.
- [ ] `spec check no-such-slug` exits `2` and names the slug on stderr.
- [ ] With no slug, the command checks every active Spec in the Spec Root.
- [ ] `roundfix --help` lists the command in the usage block and the command
      list.
- [ ] A test asserts `implement` gained no consistency precondition.

## Context

- instruction: `docs/agents/cli.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/archive.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1 -run 'SpecCheck' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix --help | grep -q "spec check"`
  — expected: exit 0; the command appears in the usage block.
- `go run -buildvcs=false ./cmd/roundfix spec check 0064-spec-artifact-consistency-gate --format json | grep -q "roundfix-speccheck/v1"`
  — expected: exit 0; the command runs against this Spec and emits its schema.

## References

- `_prd.md` → Core Feature 1; Decisions.
- `_techspec.md` → API Contracts; Build Order 5.
- ADR-0093.
