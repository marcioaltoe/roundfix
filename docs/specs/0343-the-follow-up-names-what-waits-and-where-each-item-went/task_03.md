---
task: task_03
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
status: completed
type: docs
complexity: medium
---

# Task 03: The follow-up mode passes --all, writes the declaration and names each item's destination

## Overview

The Roundfix skill's `--follow-up` Skill Mode lists pending work with
`roundfix follow-up`, plans, authors and reports, but its report accounts only
for the handoff's items and the messages it handed in. After this Task the
mode passes `--all` when the request asks for it, writes a Waiting Declaration
on any Finding or Backlog Entry its plan leaves waiting, and its report names
the destination of every listed `finding` and `backlog` item and the
held-back count, so an item without a destination is a visible gap. This
answers the second request of the adopted Backlog Entry of 2026-11-02, "The
follow-up names what waits on someone else and where each item went"
([references](references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md)),
reported by fluxus.

## Requirements

1. MUST change the canonical `.agents/skills/roundfix/references/modes.md`
   as `_techspec.md` → API Contract 8 states, under its existing headings:
   - "How a mode is chosen": after `follow-up`, the word `all` or `--all`
     asks for Waiting Items, and step 1 then runs the command with `--all`;
     any other word keeps the existing rule;
   - `--follow-up` step 1: the command call is
     `roundfix follow-up --all --format json` only when the request asked for
     it, `roundfix follow-up --format json` otherwise, each with the messages
     file when that source exists, and the stop report for nothing to do
     names the held-back count the sources give;
   - `--follow-up` step 2: the plan leaves a listed Finding or Backlog Entry
     waiting only by writing its Waiting Declaration, `waiting_on: <who or
     what>` and `since: <YYYY-MM-DD>`, in the item's frontmatter (a Finding
     also updates `updated_at`), and the report names a declaration that is
     still uncommitted, as it names an uncommitted Backlog Entry;
   - `--follow-up` step 8: the report names, for every listed `finding` and
     `backlog` item, the Spec it became, the Spec it was grouped into, or
     `waiting on <waiting_on> since <since>`, with the sentence "A listed item
     without a destination is a gap in the report, as a handed-in message
     left pending is.", and it gives the number of Waiting Items the command
     held back.
   Change no other step, the `--status` section or the `--unattended`
   section.
2. MUST change the `--follow-up` line of the `## Skill Modes` summary in the
   canonical `.agents/skills/roundfix/SKILL.md` to also say that the mode
   names where each listed Finding and Backlog Entry went, changing no other
   line except the version fields the record step raises.
3. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`. Never hand-edit a mirror or a version
   field.
4. MUST re-record the Behavior Surface Record with `go test -count=1 -tags
   docscontract ./internal/docscontract -run
   '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` when the
   coverage map moved, and commit `docs/references/behavior-surfaces.json`.
5. MUST NOT touch any Go file, the follow-up, deliver or any other
   reference, any guide, `CONTEXT.md`, any `### QA settlement` section, any
   other skill, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Add the `all` mode word and the `--all` call to the mode's first step.
- [ ] Add the declaration rule to the plan step and the destination rule to the report step.
- [ ] Update the `SKILL.md` mode summary.
- [ ] Sync, record the skill version and, when needed, the fingerprints.

## Acceptance Criteria

- [ ] Both copies of the `modes` reference state the `all` word, the `--all` call, the declaration a waiting item gets and the destination of every listed Finding and Backlog Entry, with a missing destination named a gap.
- [ ] The `SKILL.md` summary names the destinations, the mirrors equal their canonical copies, and the skill version record passes.

## Context

- instruction: `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned-skill version record — consumers: `./skills`

## Verification

- `for f in .agents/skills/roundfix/references/modes.md skills/roundfix/references/modes.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A listed item without a destination is a gap in the report' || { printf 'missing destination rule in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'waiting on <waiting_on> since <since>' || { printf 'missing waiting destination in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'roundfix follow-up --all --format json' || { printf 'missing --all call in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'names where each listed Finding and Backlog Entry went' || { printf 'missing summary in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/modes.md skills/roundfix/references/modes.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the mode names no destination, so the first check fails; afterwards both copies of the reference and of `SKILL.md` carry the rules, the mirrors equal their canonical copies, and the skill version record passes.

## References

- `_prd.md` → Goals; Core Feature 6; Success Metric 8; Decisions
- `_techspec.md` → Current behavior 6; API Contract 8; Build Order 3

## Result

Skills loaded: implement-task, tech-writer, evidence-gate

Implementation evidence:

- The canonical `modes.md` now recognizes `all` and `--all`, selects the matching `roundfix follow-up --all --format json` call, records the Waiting Declaration frontmatter rule, reports uncommitted declarations, names every listed Finding and Backlog Entry destination, states the missing-destination gap, and reports the held-back Waiting Item count.
- The canonical `SKILL.md` summary now names where each listed Finding and Backlog Entry went. `make skills-sync` completed, and the canonical/mirror pairs have identical SHA-256 hashes.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed during the initial implementation; after the Verification Feedback repair, the final record is `0.1.176` in both skill version fields and the owned-skill record.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed; the Behavior Surface Record required no byte change.
- Verification Feedback repair: the diagnostic artifact showed that the existing modes contract test requires the original `--follow-up` summary sentence verbatim. The canonical summary now preserves that sentence and adds the destination sentence, then `make skills-sync` propagated it without changing any Go test.
- `go test ./skills -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' -count=1`: passed, 3 tests; the focused failure is repaired. The owned-skill recorder raised both copies and the record to `0.1.176`, and the coverage recorder passed again with no Behavior Surface Record change.

Acceptance evidence:

1. Both `modes.md` copies contain the `all`/`--all` call, declaration, destination, and gap rules; the pair hashes match.
2. Both `SKILL.md` copies contain the destination summary; the pair hashes match, and the owned-skill version recorder passed.

Not run: the Task's `## Verification` commands; the Daemon owns those checks and Task settlement.
