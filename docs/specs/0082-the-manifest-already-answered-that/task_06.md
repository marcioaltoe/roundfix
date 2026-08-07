---
task: task_06
spec: 0082-the-manifest-already-answered-that
status: pending
type: backend
complexity: medium
---

# Task 06: Ask only what the manifest does not already answer

## Overview

The interactive command announces `update` mode and then asks all twelve
decisions anyway, each offering the stored value as a default, before blocking
on a spawned ACP runtime that re-segments the root instruction corpus. This task
short-circuits that path: when the manifest resolves, the interactive workflow
skips the settled prompts and the classification step and goes straight to the
plan, prompting only for decisions the manifest does not carry. First adoption
is untouched.

## Requirements

1. MUST skip the preservation-mode prompt, the profile prompt, and every decision
   prompt whose answer the resolved manifest carries.
2. MUST prompt for exactly those decisions the current catalog requires and the
   manifest does not carry, and for nothing else.
3. MUST NOT invoke the semantic analyzer, and therefore MUST NOT spawn an ACP
   runtime, when the manifest resolves.
4. MUST keep the plan confirmation prompt and its revision flow, so mutation
   still happens only against a reviewed Plan Digest.
5. MUST leave first adoption — a repository with no manifest — behaviorally
   identical, including its preservation prompt and supervised classification.
6. MUST fall back to the current full interactive path when the manifest is
   present but unreadable or incompatible, rather than refusing.
7. MUST keep the profile-change route reachable, so a maintainer can still move
   a repository to a different Baseline Profile interactively.

## Subtasks

- [ ] Short-circuit the settled prompts when the manifest resolves.
- [ ] Prompt only for newly required decisions.
- [ ] Skip classification on the resolved path.
- [ ] Keep the profile-change route reachable from the short-circuited path.
- [ ] Prove first adoption and the incompatible-manifest fallback are unchanged.

## Acceptance Criteria

- [ ] On a repository with a complete current manifest, the interactive workflow
      reaches the plan confirmation having issued zero decision prompts.
- [ ] On a repository whose manifest lacks a catalog-required decision, exactly
      that decision is prompted.
- [ ] No semantic analyzer call occurs on the resolved path, proven by a test
      whose injected analyzer fails the test if it is called.
- [ ] On a repository with no manifest, the prompt sequence matches the task_01
      characterization corpus exactly.
- [ ] On a repository with an unreadable manifest, the full interactive path runs
      rather than the command refusing.
- [ ] On a repository whose stored profile digest no longer matches the catalog
      but whose profile resolves and whose decisions validate, the workflow
      announces update rather than adoption and issues zero decision prompts.
- [ ] A maintainer can still choose to change the Baseline Profile.

## Context

- interface: `internal/cli/baseline_human.go`
- interface: `internal/baselineacp/analyzer.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineHuman' -v 2>&1 | grep -q '^--- PASS: .*BaselineHuman'` — expected: exits 0.
- `go test ./internal/cli/ -run 'BaselineHuman' -v 2>&1 | grep -q -i 'analyzer'` — expected: exits 0, proving the never-called-analyzer case ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the task_01 corpus proving first adoption is unchanged.

## References

- `_techspec.md` → Build Order 8; System Architecture.
- `_prd.md` → Core Features 7 and 8; User Story 2; Goal 1; Non-Goals: first adoption, profile changes.
- ADR-0068, ADR-0069, ADR-0099.
