---
task: task_03
spec: 0354-an-external-skill-is-checked-only-for-presence
status: pending
type: docs
complexity: medium
---

# Task 03: The Roundfix skill, the glossary and four ADRs state that an external skill is checked only for presence

## Overview

The Roundfix `baseline`, `setup` and `runtime` references still say that an external skill must hash
to its lock entry, that Doctor names `DR-SKILL-DIFFERS-FROM-SNAPSHOT` and that
the managed refresh lists `skills.differing`; the glossary defines a Missing
Required Skill by its lock hash; and ADR-0221, ADR-0224, ADR-0241 and
ADR-0264 still prescribe or name the comparisons tasks 01 and 02 removed. This Task writes the rule
of ADR-0364 into those documents, word for word from `_techspec.md` → API
Contract 7, and the four amendment notices. It answers the fluxus report of
2026-11-05, adopted as
[references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md](references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md).

## Requirements

1. MUST make in the canonical `.agents/skills/roundfix/references/baseline.md`,
   `.agents/skills/roundfix/references/setup.md` and
   `.agents/skills/roundfix/references/runtime.md` every change API Contract
   7 names for each, word for word.
2. MUST revise, through `domain-modeling`, the `CONTEXT.md` entries
   **Missing Required Skill** and **Setup Snapshot Difference** to the text
   API Contract 7 gives, with each `ADR-NNNN` written as the glossary's
   Markdown link to `docs/adr/<file>.md`, keeping their `_Avoid_` lines, and
   add no other capitalized phrase to `CONTEXT.md`.
3. MUST write the four amendment notices of API Contract 7 into
   `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`,
   `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`,
   `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
   and
   `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`,
   each after the paragraph API Contract 7 names, with `ADR-0364` in the bold
   opening as a Markdown link to
   `0364-a-required-external-skill-is-checked-only-for-presence.md`, and set
   each file's `updated_at` to `2026-11-05T12:00:00Z`. No other line of the
   four ADRs changes.
4. MUST run `make skills-sync` to copy the canonical references to
   `skills/roundfix/references/`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and never hand-edit a mirror, a version field or the record.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last skill edit, and never hand-edit
   it.

## Subtasks

- [ ] Write the rule into the Roundfix references, sync the mirrors and record the skill version.
- [ ] Revise the two glossary entries and write the four amendment notices.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the Roundfix `baseline` and `setup` references carry the rule sentence, and no copy of the three references names `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or `skills.differing`.
- [ ] Both copies of the Roundfix `runtime` reference say each external skill is checked only for presence.
- [ ] The two glossary entries say that a present copy is never a Missing Required Skill and that Roundfix no longer looks for a Setup Snapshot Difference.
- [ ] ADR-0221, ADR-0224, ADR-0241 and ADR-0264 carry an amendment notice naming ADR-0364 and an `updated_at` of 2026-11-05.
- [ ] The owned-skill record, the mirror check and the documentation contracts pass.

## Context

- instruction: `docs/adr/0364-a-required-external-skill-is-checked-only-for-presence.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`
- interface: `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`
- interface: `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- interface: `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill references and version record — consumers: `./skills`
- The glossary, the ADR notices and the skill fingerprints in the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `rule='checks a required external skill only for presence: an entry in .skills-lock\.json. and a folder under .\.agents/skills.'; for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- "$rule" || { echo "missing rule: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'checks each remaining skill only for presence' || { echo "missing runtime rule: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/setup.md .agents/skills/roundfix/references/runtime.md; do ! grep -q -- 'DR-SKILL-DIFFERS-FROM-SNAPSHOT\|skills\.differing' "$f" || { echo "old comparison still named: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task no Roundfix reference carries the rule sentence, so the first check fails; afterwards each carries it, none names the retired code or field, and the owned-skill record and the mirror check pass with the raised version.

- `t="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; printf '%s' "$t" | grep -qF -- 'a present copy is never one, whatever its content' || { echo 'missing **Missing Required Skill** revision' >&2; exit 1; }; printf '%s' "$t" | grep -qF -- 'so Roundfix no longer looks for one' || { echo 'missing **Setup Snapshot Difference** revision' >&2; exit 1; }; for f in docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md; do grep -qF -- '**Amended by [ADR-0364](0364-a-required-external-skill-is-checked-only-for-presence.md) (2026-11-05).**' "$f" || { echo "missing amendment: $f" >&2; exit 1; }; grep -q -- '^updated_at: 2026-11-05T12:00:00Z$' "$f" || { echo "missing updated_at: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither glossary entry carries its revision, so the first check fails; afterwards both do, the four ADRs carry their notices, and the documentation contracts, the glossary phrase test and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 7; Glossary
- `_techspec.md` → API Contract 7; Build Order 3; Glossary
