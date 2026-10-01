---
task: task_04
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: chore
complexity: high
---

# Task 04: The Jev Router, measured on four replayed docs and chore Tasks

## Overview

The adopted Backlog Entry's hypothesis is that `docs` and `chore` Tasks run
through the Jev Router cost less per settled Task without more corrective
work. This Task runs the TechSpec's measurement protocol on the maintainer's
machine: four archived Tasks, each on the router and on the current default,
in scratch clones of this repository, under the gate Task 03 added. It writes
the measurement record and, when the router holds up, a follow-up Backlog
Entry. It changes no profile and no default.

## Requirements

1. MUST stop with `jev_router_key_missing` in its result, writing none of the
   files below, when `ROUNDFIX_OPENROUTER_API_KEY` is not set. It MUST NOT
   print, echo or write the key, and MUST NOT read `OPENROUTER_API_KEY`.
2. MUST run with `NODE_OPTIONS` unset and with `bin/roundfix` rebuilt from
   the tree that holds tasks 01 to 03.
3. MUST follow the TechSpec "The measurement protocol" exactly: the four
   named Tasks or the named reserves, the pre-state construction, the routed
   replay first under a scratch `.roundfixrc.yml`, the default replay with the
   one-Run override, and a fresh scratch clone of this repository per replay.
4. MUST stop starting routed replays as soon as the gate refuses one, and
   record that refusal; it MUST NOT raise, bypass or work around the ceiling.
5. MUST write
   `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md`
   with `## Measured` (one row per replay: Task, category, selection, wall
   time, prompts, Verification repairs, outcome, fallback and reason, tokens,
   router cost) and `## Reading` (cost per settled Task on the router,
   whether tool calls worked through it, and what the sample can and cannot
   say).
6. MUST mint a dated Backlog Entry under `docs/backlog/` proposing the router
   for `docs` and `chore` only when the router settled every Task the default
   settled with no more Verification repairs; otherwise `## Reading` says why
   not.
7. MUST send prompts only from scratch clones of this repository under the
   system temporary directory, remove them after recording, and MUST NOT
   change `.roundfixrc.yml`, a User Config, or any profile of this
   repository.

## Subtasks

- [ ] Check the key, or stop with the named blocker.
- [ ] Build each replay's pre-state and run the routed and default replays.
- [ ] Stop at the first ceiling refusal.
- [ ] Write the measurement record and, when supported, the follow-up entry.

## Acceptance Criteria

- [ ] The measurement record has one row per replay that ran, each routed
      row with its cost from the Judge Log, and a reading.
- [ ] Without the key, the Task ends failed with `jev_router_key_missing`
      and no file written.

## Context

- creates: `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md`
- instruction: `docs/history/specs/0210-evidence-snapshots-that-stay-small/task_02.md`
- instruction: `docs/history/specs/0194-a-skill-and-a-command-guide-read-one-command-at-a-time/task_04.md`
- instruction: `docs/history/specs/0200-a-skill-snapshot-that-matches-its-upstream/task_04.md`
- instruction: `docs/history/specs/0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06.md`

## Verification

- `record='docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md'; test -f "$record" || { printf 'the measurement record is missing\n' >&2; exit 1; }; grep -q '^## Measured' "$record" && grep -q '^## Reading' "$record" && grep -q 'jev-router' "$record" || { printf 'the measurement record is incomplete\n' >&2; exit 1; }; if grep -nE 'sk-or-[A-Za-z0-9-]{8,}' "$record"; then printf 'the record holds a key\n' >&2; exit 1; fi` — expected: exit 0; before this Task the measurement record does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 4; Core Feature 6; Success Metric 5; Acceptance evidence; Open Questions
- [_techspec.md](_techspec.md) — The measurement protocol; Risks & Considerations; Testing Approach 5; Build Order 4
- ADR-0218; ADR-0211
