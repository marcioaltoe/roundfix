---
task: task_04
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: pending
type: chore
complexity: high
---

# Task 04: Grok through Cursor, measured on two replayed Tasks

## Overview

The adopted Backlog Entry asked which Grok model IDs Cursor offers, whether
an unattended Cursor session needs a person, how the plan bills a long Run,
and how Grok compares with the current default on real Tasks. None of it
could be measured during authoring, because `cursor-agent` was not logged in.
This Task runs the TechSpec's measurement protocol on the maintainer's
machine and existing Cursor login, commits a recorded session and its parsing
test, and writes the measurement record. It changes no profile and no
default.

## Requirements

1. MUST first run `cursor-agent status`. Without a login, it MUST record the
   blocker `cursor_login_required` in its result and stop, writing none of
   the files below. It MUST NOT run `cursor-agent login`, pass a key or token,
   or read `CURSOR_API_KEY` or `CURSOR_AUTH_TOKEN`.
2. MUST run with `NODE_OPTIONS` unset and with `bin/roundfix` rebuilt from
   the tree that holds tasks 01 and 02.
3. MUST follow the TechSpec "The measurement protocol" steps 2 to 6 exactly:
   capture one real session's `configOptions` with every account field
   removed; prove one Grok selection; replay
   `0210-evidence-snapshots-that-stay-small/task_02` and
   `0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06`, or
   the named reserves, on the Grok selection and on the current default, in a
   fresh scratch clone of this repository per replay.
4. MUST write `internal/agent/testdata/cursor-session-recorded.json` and
   `internal/agent/cursor_recorded_session_test.go` with
   `TestCursorRecordedSessionAdvertisesGrok`, which parses the recording under
   `cursor`, asserts at least one advertised value starting with `grok`, and
   assigns that value as `model_managed`.
5. MUST write
   `docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md`
   with `## Measured` (one row per replay: selection, wall time, prompts,
   Verification repairs, outcome, person needed, tokens; plus the advertised
   Grok values, whether acpx opened the session alone, whether `session/load`
   worked, and the billing mode the value states) and `## Reading`.
6. MUST mint a dated Backlog Entry under `docs/backlog/` proposing a Grok
   fallback only when Grok settled both replays with no person needed;
   otherwise `## Reading` says why not.
7. MUST send prompts only from scratch clones of this repository under the
   system temporary directory, remove them after recording, and MUST NOT
   change `.roundfixrc.yml`, a User Config, or any profile.

## Subtasks

- [ ] Check the login, or stop with the named blocker.
- [ ] Capture and commit one real session, with its parsing test.
- [ ] Prove a Grok selection and run the four replays.
- [ ] Write the measurement record and, when supported, the follow-up entry.

## Acceptance Criteria

- [ ] The recording parses under `cursor`, advertises a `grok` value, and
      that value assigns.
- [ ] The measurement record has one row per replay and a reading.
- [ ] Without a login, the Task ends failed with `cursor_login_required` and
      no file written.

## Context

- creates: `internal/agent/testdata/cursor-session-recorded.json`
- creates: `internal/agent/cursor_recorded_session_test.go`
- creates: `docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md`
- instruction: `internal/agent/selection_capabilities.go`
- instruction: `docs/history/specs/0210-evidence-snapshots-that-stay-small/task_02.md`
- instruction: `docs/history/specs/0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06.md`

## Verification

- `out="$(go test -count=1 -v -run '^TestCursorRecordedSessionAdvertisesGrok$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestCursorRecordedSessionAdvertisesGrok' || { printf '%s\n' "$out"; exit 1; }; record='docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md'; grep -q '^## Measured' "$record" && grep -q '^## Reading' "$record" && grep -q 'grok' "$record" || { printf 'the measurement record is incomplete\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the recorded session nor the record exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 4; Core Feature 5; Success Metric 5; Acceptance evidence; Open Questions
- [_techspec.md](_techspec.md) — The measurement protocol; Testing Approach 5; Risks & Considerations; Build Order 4
- ADR-0217; ADR-0211
