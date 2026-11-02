---
task: task_04
spec: 0338-the-release-decision-travels-with-the-spec
status: pending
type: docs
complexity: medium
---

# Task 04: The authoring skills fill the Delivery Type and ask the release question

## Overview

The Delivery Type and the Release Ceiling only help when the plan records
them. `write-prd` writes the PRD's `type` for a feature Spec, `write-techspec`
writes `fix` or `refactor` in the minimal PRD of its correction route, and
both ask the release question at planning time when the maintainer asks for a
release with the Spec, recording the answer in `_authorization.md`.
`write-tasks` keeps the release out of every Task and refuses a malformed
ceiling. This answers both adopted Backlog Entries of 2026-11-01
([references](references/2026-11-01-the-delivery-pull-request-type-follows-the-spec.md),
[references](references/2026-11-01-a-spec-authorization-carries-the-release-decision.md)),
whose Expected sections ask the authoring skills to fill the type and ask the
release question.

## Requirements

1. MUST add to the frontmatter block of
   `.agents/skills/write-prd/references/prd-template.md`, after `surfaces`,
   the line `type: feat # Delivery Type: feat | fix | perf | refactor; a feature Spec writes feat or perf`.
2. MUST add to `.agents/skills/write-prd/SKILL.md`, before `## Advisory
   judgment`, a section `## Delivery Type and release decision` that says:
   the PRD declares its Delivery Type in `type`, `feat` for a feature or
   `perf` for a performance change, and the Archive Record and the Delivery
   Queue's pull request title carry it; when the maintainer asks for a
   release with the Spec, ask the release question with the structured
   question tool (the Release Ceiling, `patch`, `minor` or `major`, and the
   version-only paths the release may change), then record `release` in
   `operations`, `release_ceiling` and `release_paths` in `_authorization.md`
   and quote the answer verbatim in its body; without that request, never
   add `release`. It contains the sentence "Ask the release question at
   planning time, never after the merge".
3. MUST state in `.agents/skills/write-techspec/SKILL.md` that the minimal
   PRD's frontmatter carries `type: fix` for a bug fix and `type: refactor`
   for a refactor (where the paragraph lists `spec`, `status: active`,
   `surfaces`), and add before `## Advisory judgment` the same `## Delivery
   Type and release decision` section, with the same sentence, worded for the
   correction route, and saying a `refactor` delivery that changes released
   code makes the release plan ask for a manual classification.
4. MUST add to `.agents/skills/write-tasks/SKILL.md`, before `## Process`, a
   section `## Release decision` that says: no Task performs the release,
   tags, edits `CHANGELOG.md`, a version file or a package manifest; a
   `release_ceiling` and its `release_paths` in `_authorization.md` grant the
   later release plan's answer, not a Task path, so they are never copied
   into a Task's `creates:` or `interface:`; and a record whose
   `release_ceiling` lacks `release` or `release_paths` is a decomposition
   refusal, because the reader refuses the whole record. It contains the
   sentence "A Release Ceiling answers the release plan, never a Task".
5. MUST NOT change the `### QA settlement` section or any other section of
   the three skills beyond the lines above, and MUST NOT touch the Roundfix
   skill, any guide, any Baseline module or any Go file.
6. MUST copy the canonical skills with `make skills-sync`, then raise both
   version fields of each of the three changed skills and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`. Never hand-edit a mirror or a version
   field.
7. MUST re-record the Behavior Surface Record as the LAST step with
   `go test -count=1 -tags docscontract ./internal/docscontract -run
   '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` when the
   coverage map moved, and commit `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Add the template line and the two `write-prd`/`write-techspec` sections.
- [ ] Add the `write-tasks` section.
- [ ] Sync the mirrors, record the three skill versions, re-record the fingerprints last.

## Acceptance Criteria

- [ ] The PRD template carries `type`, and `write-prd` and `write-techspec` fill it and ask the release question at planning time.
- [ ] `write-tasks` keeps the release out of every Task and refuses a malformed ceiling.
- [ ] The mirrors equal their canonical copies, and the skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`
- instruction: `.agents/skills/write-prd/references/clarification-record.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `write-prd`, `write-techspec` and `write-tasks` skill texts and the owned-skill version record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/write-prd/references/prd-template.md | grep -qF -- 'type: feat # Delivery Type: feat | fix | perf | refactor' || { echo 'template lacks the Delivery Type' >&2; exit 1; }; for f in .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md; do grep -qF -- '## Delivery Type and release decision' "$f" || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Ask the release question at planning time, never after the merge' || { printf 'missing release question in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/write-techspec/SKILL.md | grep -qF -- 'type: refactor' || { echo 'write-techspec lacks the refactor type' >&2; exit 1; }; for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Release Ceiling answers the release plan, never a Task' || { printf 'missing release rule in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/write-prd/references/prd-template.md skills/write-prd/references/prd-template.md && cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the template has no `type`, so the first check fails; afterwards the three skills carry the rules, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Sources; Goals; Core Features 1 and 5; Success Metric 8; Decisions
- `_techspec.md` → API Contract 12; Build Order 4
- ADR-0349
