---
task: task_03
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: completed
type: docs
complexity: medium
---

# Task 03: The Supervisor announces the Run and queue events and the guide lists every owner

## Overview

No Agent inside a Run can announce, so the Roundfix skill's implement
reference makes the Supervisor session that started the Run or queued the
Delivery Queue item announce `started`, `blocked` and, for a Spec the queue
archives, `completed` when the item merges. The context-driven development
guide then tells users the step exists, that their own instructions define it,
and which skill announces each Spec Event. This answers the adopted Backlog
Entry of 2026-11-03
([references](references/2026-11-03-the-spec-skills-announce-each-spec-event.md)),
reported by jarvis, whose events "iniciada" and "bloqueada" these are.

## Requirements

1. MUST add to `.agents/skills/roundfix/references/implement.md`, right after
   the `### Token usage` section and before `## Fixture values`, a section
   headed `## Spec Event Announcement` that holds the paragraph of
   `_techspec.md` → API Contract 1 and then the text of API Contract 6, its
   list included, word for word (line wrapping is free).
2. MUST NOT edit the Roundfix skill's `references/modes.md` or
   `references/follow-up.md`, which Spec 0343 owns, or any other Roundfix
   reference, and MUST NOT change the `### QA settlement` section of the
   Roundfix `SKILL.md`.
3. MUST add to `docs/user-guide/context-driven-development.md`, right before
   `### QA report reachability and verdicts`, the section
   `### Spec Event Announcements` exactly as API Contract 8 gives it, its
   ADR-0354 link and table included.
4. MUST NOT name any personal assistant, tracker, channel or message block in
   the reference or the guide.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of the Roundfix skill and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   Never hand-edit a mirror or a version field.
6. MUST re-record the Behavior Surface Record as the LAST step, after the
   last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and commit `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Add the implement reference section.
- [ ] Add the guide section.
- [ ] Sync the mirror, record the Roundfix skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] The implement reference, canonical and mirror, carries the generic paragraph and the Supervisor's `started`, `blocked` and `completed` events.
- [ ] The context-driven development guide describes the step and the owner of each of the seven Spec Events.
- [ ] The mirrors equal their canonical copies, and the skill version, reference layout and Behavior Surface contracts pass.

## Context

- instruction: `docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix implement reference, the context-driven development guide, the owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "## Spec Event Announcement" || { printf 'missing phrase in %s: Spec Event Announcement\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "this skill never names one" || { printf 'missing phrase in %s: this skill never names one\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "An Agent inside a Daemon-assigned Run never announces." || { printf 'missing phrase in %s: An Agent inside a Daemon-assigned Run ne\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "The Supervisor session that started the Run or queued the Spec's Delivery Queue item announces each event of that Run or item when it observes it" || { printf 'missing phrase in %s: The Supervisor session that started the\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "Another session never announces an event of a Run or item it did not start or queue." || { printf 'missing phrase in %s: Another session never announces an event\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/context-driven-development.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "### Spec Event Announcements" || { printf 'missing phrase in %s: Spec Event Announcements\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "the user's own instructions, global or for the repository, define the destination and the format" || { printf 'missing phrase in %s: the users own instructions, global or fo\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the implement reference has no heading, so the first check fails; afterwards the reference, its mirror and the guide carry the step, and the skill version, reference layout and Behavior Surface contracts pass.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1, 2 and 4; Success Metric 1; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 6; API Contract 8; Invariants 1-5; Build Order 3
- ADR-0354

## Result

Skills loaded: implement-task, tech-writer, evidence-gate

Implementation evidence:

- Added the Roundfix implement reference's `## Spec Event Announcement` section between `### Token usage` and `## Fixture values`, carrying the generic API Contract 1 paragraph plus the Supervisor ownership and `started`, `blocked`, and `completed` event rules from API Contract 6.
- Added the `### Spec Event Announcements` guide section immediately before `### QA report reachability and verdicts`, including the ADR-0354 link and the seven-event owner table from API Contract 8.
- Ran `rtk make skills-sync` with no output or error; the canonical and mirror implement references now have the same SHA-256 (`f60bad94d8d8925684fae6a459ffc0f1b6944744418ae0890223ed9617c4b0a5`), and the Roundfix `SKILL.md` pair has the same SHA-256 (`aeaf494abe3652ea56b9113cc44450db10ec2be1a2b23c334d66d87d889521f8`).
- Ran `rtk go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed (1 test package); both Roundfix skill version fields recorded `0.1.177`.
- Ran the final `rtk go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed (1 test package); `docs/references/behavior-surfaces.json` was re-recorded after the guide edit.
- Focused checks passed: `rtk git diff --check`; section placement inspection; canonical/mirror hash comparison; and prohibited-name scan for the reference and guide.

Acceptance evidence: the canonical and synced mirror implement references contain the generic announcement paragraph and Supervisor event rules; the guide contains all seven Spec Event owners; the sanctioned sync/version/coverage recorders passed and produced matching skill copies. The task's declared Verification commands were not run because this is a Daemon-assigned turn; the Daemon owns Verification and settlement.
