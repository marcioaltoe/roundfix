---
task: task_04
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: pending
type: docs
complexity: medium
---

# Task 04: The baseline guide and the skills describe the gap refusal and the guide rule

## Overview

Tasks 02 and 03 changed what `baseline profile validate`, `baseline plan`
and `baseline update` do with a repository profile that lacks a template
decision, and which guides raise a nested-carrier warning. This Task writes
both rules into the baseline command guide, the Roundfix skill's baseline
reference and the setup skill, whose promise "A single-gate Setup Manifest
exits `3`, names the missing incremental decision" becomes true for a
repository profile too. Both owned skills record their next versions, and
the Behavior Surface Record is re-recorded last.

## Requirements

1. MUST add to `docs/user-guide/commands/baseline.md`, right after the
   paragraph that ends "cannot compose profiles or load remote executable
   content.", these two paragraphs, verbatim:
   - "A repository profile must list every decision a template it renders requires. When it does not, `baseline profile validate` exits `2` with `custom.profile.decision.template-required`, and `baseline plan` and `baseline update` exit `3` with category `decision`. Each names the decision, the template, the entry to add to the profile's `decisions` and a `--decision` flag with the suggested value, and none writes a byte or edits the profile. Once the profile lists the decision, `baseline update` asks for it as a new decision."
   - "A guide under `docs/agents/` that has no `setup-context-driven` marker line and sits at a path the Baseline does not manage belongs to the repository and raises no `baseline.inventory.nested-carrier-conflict`; naming the marker in its prose changes nothing. A nested `AGENTS.md` or `CLAUDE.md`, an unmarked guide at a managed path and a guide with an unknown marker still warn."
2. MUST add the same two paragraphs to the canonical Roundfix skill
   reference `.agents/skills/roundfix/references/baseline.md`, right after
   the fenced block of maintenance commands that lists
   `roundfix baseline profile validate <profile-id> --format text`.
3. MUST add to the canonical setup skill
   `.agents/skills/setup-context-driven/SKILL.md`, right after the paragraph
   that ends "later updates use the recorded decision without that flag.",
   this paragraph, verbatim: "A repository-owned Baseline Profile that does not list a decision a template it renders requires is refused first: `baseline update` exits `3` and names the decision, the entry to add to the profile's `decisions` and a `--decision` flag, and `baseline profile validate` reports the same gap with exit `2`. Add the entry to the profile, then rerun the update with `--adopt-suggested` or that flag."
4. MUST raise both version fields (`metadata.version` and `version`) of the
   Roundfix skill and of the setup skill by one patch step over the values
   the tree holds when the Task runs, run `make skills-sync` to refresh the
   mirrors under `skills/`, and re-record the version record with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   No mirror, version record or digest is hand-edited.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's LAST step, after the last edit of the user guide and of the
   skills.
6. MUST NOT change any other paragraph of the three files, `CONTEXT.md`, any
   Go file, any Baseline asset or guide under `docs/agents/`, or any other
   skill.

## Subtasks

- [ ] Add the two paragraphs to the baseline command guide and the Roundfix reference.
- [ ] Add the paragraph to the setup skill.
- [ ] Raise both skills' versions, sync the mirrors and re-record the version record.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The guide and the Roundfix reference state the gap refusal with its exits and the repository-authored guide rule.
- [ ] The setup skill tells a repository-profile maintainer what to add and how to rerun.
- [ ] Both owned skills record their next versions, and the skill and documentation contracts pass.

## Context

- instruction: `docs/specs/0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings/_techspec.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned skill versions and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'A repository profile must list every decision a template it renders requires.' || { printf 'missing phrase in %s\n' docs/user-guide/commands/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'belongs to the repository and raises no' || { printf 'missing phrase in %s\n' docs/user-guide/commands/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'A repository profile must list every decision a template it renders requires.' || { printf 'missing phrase in %s\n' .agents/skills/roundfix/references/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'belongs to the repository and raises no' || { printf 'missing phrase in %s\n' .agents/skills/roundfix/references/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/roundfix/references/baseline.md | grep -qF -- 'A repository profile must list every decision a template it renders requires.' || { printf 'missing phrase in %s\n' skills/roundfix/references/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/roundfix/references/baseline.md | grep -qF -- 'belongs to the repository and raises no' || { printf 'missing phrase in %s\n' skills/roundfix/references/baseline.md >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/setup-context-driven/SKILL.md | grep -qF -- 'A repository-owned Baseline Profile that does not list a decision a template it renders requires is refused first' || { printf 'missing paragraph in %s\n' .agents/skills/setup-context-driven/SKILL.md >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/setup-context-driven/SKILL.md | grep -qF -- 'A repository-owned Baseline Profile that does not list a decision a template it renders requires is refused first' || { printf 'missing paragraph in %s\n' skills/setup-context-driven/SKILL.md >&2; exit 1; }; go test -timeout 20m -count=1 ./skills && go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms)$'` — expected: exit 0. Before this Task none of the files carries the paragraphs, so the first check fails; afterwards the mirrors match, the version record is current and the Behavior Surface Record matches the guide.

## References

- `_prd.md` → Goals; User Story 2; User Story 3; Core Feature 4
- `_techspec.md` → API Contracts 4-6; Build Order 4
- ADR-0372
