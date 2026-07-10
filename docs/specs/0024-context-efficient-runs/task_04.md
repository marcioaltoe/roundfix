---
task: task_04
spec: 0024-context-efficient-runs
status: pending
type: backend
complexity: medium
---

# Task 04: Compact Agent reads and edits

## Overview

Render repetitive ACP file reads and edits as bounded one-line Console Log
summaries while preserving the original Agent Run Event payload. The slice is
verifiable with measured-shape fixtures containing 31 edits and 330 reads and
with byte-for-byte journal payload assertions.

## Requirements

1. MUST retain ACP Tool Kind, locations, and diff old/new line information in the structured stream model.
2. MUST render reads as `read <path> (N lines)` and edits as `edit <path> (+N/-M)`.
3. MUST omit raw input, raw output, file bodies, ACP JSON, and unified diffs from Console Log and Live Run View text.
4. MUST degrade incomplete structured tool metadata to a bounded marker rather than raw payload output.
5. MUST keep original Agent Run Event payload bytes lossless and must not enable acpx read suppression.
6. MUST use one compact formatter for Detached Run output, journal summaries, and the Live Run View timeline.

## Subtasks

- [ ] Extend ACP parsing with structured read/edit metadata.
- [ ] Calculate deterministic read and diff line counts.
- [ ] Replace raw tool detail rendering with compact summaries.
- [ ] Add metadata-only fallback behavior.
- [ ] Preserve raw Run Event payload capture.
- [ ] Add measured-shape and Live Run View regression fixtures.

## Acceptance Criteria

- [ ] A 31-edit/330-read fixture produces exactly 31 edit lines and 330 read lines.
- [ ] No compact output contains a file body, old/new diff text, raw tool output, or serialized ACP object.
- [ ] Edit counts reflect added and removed lines from structured diff content.
- [ ] A structured read reports the correct path and line count.
- [ ] Missing optional metadata produces a bounded tool marker without exposing raw content.
- [ ] The corresponding journal events retain their original payload bytes exactly.
- [ ] Console Log and Live Run View use the same summary wording.

## Verification

- `rtk go test ./internal/agent ./internal/tui` - expected: ACP parsing, compact formatting, measured fixture, timeline, and lossless-payload tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/agent/acp_stream.go`
- interface: `internal/agent/stream.go`
- interface: `internal/agent/event.go`
- interface: `internal/tui/timeline.go`

## References

`_prd.md` -> User Story 4; Core Features 8-9; Success Metrics. `_techspec.md` -> API Contracts: Console rendering; Integration Points: acpx/ACP; Build Order 4. ADR-0008.
