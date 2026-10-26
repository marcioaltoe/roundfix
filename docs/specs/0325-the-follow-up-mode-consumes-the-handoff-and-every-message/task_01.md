---
task: task_01
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: pending
type: docs
complexity: medium
---

# Task 01: The follow-up mode takes the handoff and settles every handed-in message

## Overview

The Roundfix skill's `--follow-up` Skill Mode consumes ai-memory messages
only in its fourth step and never claims the project's Session Handoff. This
Task rewrites the mode's `Optional sources`, `--status`, step 1, step 4 and
step 8 text and the two `## Modes` bullets of the skill entry, syncs the
mirror, records the owned-skill version, and adds a contract test that pins
the text in both layouts.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-follow-up-mode-leaves-the-handoff-and-handed-in-messages-pending.md`](references/2026-10-26-the-follow-up-mode-leaves-the-handoff-and-handed-in-messages-pending.md)
   of 2026-10-26 (`_prd.md` → Sources 1) with `_techspec.md` → Interfaces
   and API Contracts 1-3.
2. MUST replace, in the canonical `.agents/skills/roundfix/references/modes.md`,
   the `## Optional sources` paragraph, the `--status` sentence, step 1, step
   4 and the closing clause of step 8 with the texts of `_techspec.md` →
   Interfaces, byte for byte, keeping steps 2, 3 and 5-7,
   `## How a mode is chosen` and `## --unattended` unchanged.
3. MUST replace the `--follow-up` and `--status` bullets of the `## Modes`
   section of `.agents/skills/roundfix/SKILL.md` with the two bullets of
   `_techspec.md` → Interfaces, keeping the section's first paragraph and its
   `--unattended` bullet unchanged.
4. MUST add `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` to
   `skills/roundfix_modes_test.go` as `_techspec.md` → Testing Approach
   describes, for both `roundfixLayoutSources()` layouts, comparing
   whitespace-normalized text, and keep every existing test of the file
   unchanged.
5. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the `roundfix` skill and records them),
   then `make skills-sync` again so the mirror matches.
6. MUST NOT change the `roundfix follow-up` command, any Go file other than
   `skills/roundfix_modes_test.go`, any guide under `docs/user-guide/`, any
   Baseline module or `docs/agents/` file, any other skill, or the
   `### QA settlement` section of any skill.

## Subtasks

- [ ] Rewrite the five passages of the canonical `modes` reference and the two `## Modes` bullets.
- [ ] Add the contract test.
- [ ] Sync the mirror, record the owned-skill version and sync again.

## Acceptance Criteria

- [ ] Step 1 takes the Session Handoff before the command runs and hands its items in as `handoff:<handoff id>:<n>` entries, in both layouts.
- [ ] Step 4 leaves no handed-in message pending, also on an early stop, and step 8 reports the handoff and the consumed messages.
- [ ] `--status` and `--unattended` accept no handoff and consume no message, and the `## Modes` bullets say the same.
- [ ] The new test and every existing `skills` test pass, and the owned-skill version is recorded.

## Consumers

- The Roundfix skill's `modes` reference and `## Modes` section — consumers: `./skills`

## Context

- instruction: `skills/roundfix_layout_test.go`
- instruction: `skills/owned_skill_versions_test.go`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/roundfix_modes_test.go`

## Verification

- `grep -q 'func TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage' skills/roundfix_modes_test.go && go test -timeout 10m -count=1 ./skills && make skills-sync-check` — expected: exit 0. Before this Task the new test does not exist, so the first check fails; afterwards the canonical skill and its mirror carry the handoff and settlement text the new test pins, every existing `skills` test passes, the owned-skill version is recorded and the mirror matches the canonical copy.

## References

- `_prd.md` → Sources 1; Goals 1-5; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 1-6; Interfaces; Data Models; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 1
