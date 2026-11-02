---
task: task_02
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
status: pending
type: docs
complexity: medium
---

# Task 02: The glossary, the ADRs and the Roundfix skill say what an approved update deletes

## Overview

ADR-0350, written by the planning change, supersedes ADR-0246 and narrows
ADR-0191. This Task retires ADR-0246 by deletion, amends ADR-0191, revises
**Retired Skill** in `CONTEXT.md`, and rewrites the retired-skills paragraph of
the Roundfix skill's baseline reference to describe the deletion task_01
ships, as task_01 did for the baseline guide, with the skill version and the Behavior Surface
Record re-recorded. It completes the adopted Backlog Entry of 2026-11-01, "An
approved Baseline update deletes the Retired Skills"
([references](references/2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md)),
the Argus report of 2026-11-01.

## Requirements

1. MUST delete
   `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
   under the retirement clause of `docs/agents/docs-layout.md`: ADR-0350
   names it with its date. Name the reason in the commit message.
2. MUST append to `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
   an amendment beginning "Amendment (2026-11-01): ADR-0350 lets an approved
   Baseline update delete a Retired Skill it proves unedited." that states
   the never-delete rule still holds for every other skill, renamed skills
   included, and set its `updated_at` to `2026-11-01T00:00:00Z`. Change no
   other line.
3. MUST revise **Retired Skill** in `CONTEXT.md` through `domain-modeling`:
   keep its first two sentences' meaning, replace the clause on listing a
   copy "without deleting it" with: an approved `roundfix baseline update`
   deletes a copy it proves unedited against `skills-lock.json` and keeps an
   edited or unrecorded copy for the adopter, and link
   `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`
   instead of ADR-0246. The phrase "deletes a copy it proves unedited against"
   appears verbatim. Change no other entry.
4. MUST rewrite the retired-skills paragraph of the canonical
   `.agents/skills/roundfix/references/baseline.md` as `_techspec.md` → API
   Contract 7 states: the sentence "An approved update deletes each Retired
   Skill it proves unedited" appears verbatim; the paragraph names the
   folder, the `.claude/skills` entry and the lock entry, the Plan Digest,
   `Skills deleted` (`skills.deleted` in JSON), the three keep reasons and
   that nothing outside the Retired Skills set is deleted, and no longer
   says the update never deletes the copy. Keep the list of Retired Skills
   and every other paragraph unchanged.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. Never hand-edit a mirror or a
   version field.
6. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of the skill, with `go test -count=1 -tags docscontract
   ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
   -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
7. MUST NOT edit `CHANGELOG.md` (the operator corrects the v0.127.0 line in
   the release), ADR-0350, ADR-0206, the history-sanitize ADR,
   `docs/user-guide/commands/baseline.md`,
   `docs/references/2026-10-08-baseline-audit.md`, any `### QA settlement`
   section, any Baseline module, `docs/agents/`,
   `docs/references/skill-coverage.json`, the Makefile, CI workflows,
   `go.mod` or `go.sum`.

## Subtasks

- [ ] Delete ADR-0246 and amend ADR-0191.
- [ ] Revise the glossary entry.
- [ ] Rewrite the reference paragraph, sync the mirror and record the skill version.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] ADR-0246 is gone, ADR-0191 carries the amendment, and `CONTEXT.md` defines **Retired Skill** with the deletion and links ADR-0350.
- [ ] Both copies of the reference state the deletion, its limits and its report, and neither says the update never deletes a retired copy.
- [ ] The mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/agents/docs-layout.md`
- interface: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- interface: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill's baseline reference, its version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `test ! -e docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md || { echo 'ADR-0246 is not retired' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md | grep -qF -- 'Amendment (2026-11-01): ADR-0350 lets an approved Baseline update delete a Retired Skill it proves unedited.' || { echo 'missing ADR-0191 amendment' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Retired Skill**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'deletes a copy it proves unedited against' || { echo 'missing deletion in the Retired Skill entry' >&2; exit 1; }; grep -qF -- '0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md' CONTEXT.md || { echo 'CONTEXT.md does not link ADR-0350' >&2; exit 1; }; if grep -qF -- '0246-a-retired-skill' CONTEXT.md; then echo 'CONTEXT.md still links ADR-0246' >&2; exit 1; fi` — expected: exit 0. Before this Task ADR-0246 exists, so the first check fails; afterwards it is retired, ADR-0191 carries its amendment, and the glossary entry states the deletion and links ADR-0350.
- `for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'An approved update deletes each Retired Skill it proves unedited' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Skills deleted' || { printf 'missing report in %s\n' "$f" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$f" | grep -qiE -- 'never deletes (an installed|the) copy'; then printf 'stale never-deletes wording in %s\n' "$f" >&2; exit 1; fi; done; cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the reference lacks the rule, so the first check fails; afterwards both reference copies state it with its report and drop the never-deletes wording, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 6; Glossary
- `_techspec.md` → API Contract 7; Declared breaks; Build Order 2
