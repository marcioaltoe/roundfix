---
task: task_01
spec: 0131-cleanup-delivery-contracts
status: pending
type: test
complexity: medium
---

# Repair the current-clause fixture and CI history input

## Overview

Implement the two-file repair approved in `_authorization.md`. This Task closes
the final observed delivery defects of the existing cleanup branch.

## Requirements

1. MUST change only `internal/docscontract/corpus_test.go` and
   `.github/workflows/ci-verify.yml`, plus this Task's Result.
2. MUST derive the current mutation input and retain the three existing source
   cases, severity/diagnostic/source assertions and real divergent ordering.
   MUST fail on missing, empty or ambiguous source targets; no success by skipping.
3. MUST fetch full history in the existing CI checkout while preserving action
   versions, permissions, credentials policy, full verification and time budget.
4. MUST record criterion evidence and the unchanged Vocabulary Contract. Do not
   edit any other implementation, test, skill, Task or generated artifact.

## Subtasks

- [ ] Derive the bounded mutation input and retain strict disagreement coverage.
- [ ] Supply full Git history through the existing CI checkout.
- [ ] Record focused evidence, OE-1 connection and glossary confirmation.

## Acceptance Criteria

- [ ] Core Feature 1: all three sources yield the required divergence error and
  diagnostic assertions after mutation; the unmodified declarations agree.
- [ ] Core Feature 2: CI fetches the existing historical evidence without relaxing
  any verification or changing its credential policy.
- [ ] OE-1: the recorded CI Git-object failure is explained by the checkout depth
  and addressed by the full-history declaration, not a test skip.
- [ ] Exactly the two authorized implementation paths changed; no glossary term changed.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/speccheck/citations.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `.github/workflows/ci-verify.yml`

## Execution environment

Use Go 1.26.7 for focused commands. Format Go with
`/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`.
The system formatter has different output; do not reformat unrelated files.
The Supervisor runs the complete repository gates after this Task. Work only on
this narrow repair; no new architecture research or unrelated sweep is needed.

## Verification

- `grep -Eq '^[[:space:]]+fetch-depth: 0[[:space:]]*$' .github/workflows/ci-verify.yml && go test -count=1 -tags docscontract ./internal/docscontract -run '^TestCheckLoopOrderDivergent$'`

## References

- `_prd.md` → Goals 1–2, Core Features 1–2, Outside Evidence OE-1.
- `_techspec.md` → Implementation Design, Testing Approach, Vocabulary Contract.
