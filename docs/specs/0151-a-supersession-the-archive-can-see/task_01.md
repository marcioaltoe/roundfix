---
task: task_01
spec: 0151-a-supersession-the-archive-can-see
status: pending
type: backend
complexity: medium
---

# Task 01: The supersession record and the command that writes it

## Overview

A Spec whose content another Spec delivered has no way to say so that a command
can read. This Task adds the record and the command that writes it.

## Requirements

1. MUST write `_supersession.md` in the superseded Spec, with frontmatter
   carrying the superseding slug, the date and the reason, and a body carrying
   the explanation given.
2. MUST leave every other file in the Spec byte-identical.
3. MUST refuse with exit 2 when the superseded Spec is unknown, when the
   superseding slug is neither active nor archived, when the two slugs are the
   same, or when the Spec already carries a supersession; and MUST name which
   condition failed.
4. MUST write nothing when it refuses.
5. MUST refuse unknown flags, matching the surrounding commands.
6. MUST NOT create a Run, write a Run Event Journal entry, commit or push.

## Subtasks

- [ ] Add the record type and its reader.
- [ ] Add the command, its refusals and its writer.
- [ ] Register it on the public command surface.
- [ ] Add tests for the accepted path and every refusal.

## Acceptance Criteria

- [ ] A fixture Spec gains `_supersession.md` naming an existing Spec, and every
      other file in it is byte-identical.
- [ ] Each of the four refusals is exercised and leaves the directory
      byte-identical.
- [ ] The command appears on the public surface.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/archive.go`
- interface: `internal/spec/spec.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSupersede" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSupersede"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix supersede --spec"` — expected: exit 0; the command appears on the public surface. Before this Task the usage has no supersede line, so the command fails.

## References

- [_techspec.md](_techspec.md) — The command
