---
task: task_01
spec: 0151-a-supersession-the-archive-can-see
status: completed
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

## Result

Implemented a typed supersession amendment with a validating reader and an
exclusive writer for `_supersession.md`. Added the `supersede` command with
required `--spec`, `--by`, and `--reason` flags, active-or-archived superseding
Spec validation, the four required exit-2 refusals, unknown-flag refusal, and
public help registration. The command does not enter the Run lifecycle.

Focused evidence:

- Acceptance criterion 1: `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/cli -run '^TestSupersede$/accepted_path_writes_only_the_amendment$'` passed. The test reads the written amendment through `spec.ReadSupersession`, checks its superseding slug, date, reason, and explanation, and compares every pre-existing file byte-for-byte after removing only the new amendment from the snapshot. It also proves HEAD is unchanged and no Run Database exists.
- Acceptance criterion 2: `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/cli -run '^TestSupersede$/refusals_leave_the_Spec_root_byte-identical$'` passed for an unknown superseded Spec, an unknown superseding Spec, self-supersession, and an existing supersession. Each case checks exit 2, the condition-specific diagnostic, an empty stdout, byte-identical Spec-root contents, and no Run Database. `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/cli -run '^TestSupersede$/unknown_flag_refuses_without_writing$'` also passed.
- Acceptance criterion 3: the binary built by the incremental gate returned exit 0 from both `rtk ./bin/roundfix --help` and `rtk ./bin/roundfix supersede --help`; the top-level output includes `roundfix supersede --spec <slug> --by <slug> --reason <text>` and the command list entry.
- Repository incremental check: `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache make verify-incremental` passed outside the filesystem sandbox, including vet, all Go packages, skill checks, and build. The first sandboxed attempt reached the test sweep but two unrelated force-stop integration tests could not read the macOS process table (`operation not permitted`); the permitted rerun passed those tests.

The Daemon-owned Verification commands were not run in this Agent turn.
