---
task: task_04
spec: 0068-spec-close-audit
status: completed
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

## Result

### Implementation

- Added `spec audit` as a sibling of `spec check`, with one required active or
  archived Spec slug, `--format text|json`, command help, top-level usage, and
  command-list discovery.
- Added text and `roundfix-specaudit/v1` JSON renderers. Text reports every
  survivor's kind and evidence, includes the exact reclaim command only for
  residue, and names the holding branch for undelivered artifacts.
- Mapped clean reports to exit `0`, residue or undelivered artifacts to exit
  `1`, and invalid input or an unknown slug to exit `2`. Diagnostics remain on
  stderr and JSON stdout contains one object.
- Kept `spec check` on its existing parser and renderer, with a command-level
  non-regression test for its clean text output and exit code.

### Focused checks

- Before implementation, `rtk go run -buildvcs=false ./cmd/roundfix spec audit
  0068-spec-close-audit` reported `unknown spec command "audit"` and the
  command's internal usage exit `2`.
- After the final code and test edits, `rtk go test ./internal/cli -run
  '^TestRunSpecAudit' -count=1` exited `0`: 7 tests passed.
- `rtk git -c core.fsmonitor=false diff --check` exited `0` after the final
  code and test edits.

### Acceptance evidence

- Clean report: `TestRunSpecAuditCleanText` audits an archived fixture and
  asserts exit `0`, the exact clean text report, and empty stderr.
- Residue: `TestRunSpecAuditResidueText` asserts exit `1`, kind `residue`, the
  classifier evidence, and `git branch -d -- '<branch>'`.
- Undelivered artifact: `TestRunSpecAuditUndeliveredTextNamesHoldingBranch`
  asserts exit `1`, the archived artifact path, and its holding branch.
- JSON: `TestRunSpecAuditJSONWritesOneObject` decodes exactly one object and
  asserts schema `roundfix-specaudit/v1` and the requested slug.
- Unknown slug: `TestRunSpecAuditUnknownSlugIsUsageError` asserts exit `2`, no
  partial JSON on stdout, and the unknown slug on stderr.
- Discovery: `TestRunSpecAuditHelpAppearsInUsageAndCommandList` asserts the
  top-level usage and command-list entries plus command-specific help.
- `spec check` non-regression: `TestRunSpecAuditPreservesSpecCheckBehavior`
  asserts its existing clean exit, report prefix, finding absence, and stderr
  behavior.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
