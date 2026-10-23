---
task: task_01
spec: 0301-the-light-tier-yields-to-the-category-profile
status: completed
type: docs
complexity: low
---

# Task 01: The guides, the skill references and the glossary describe the skip, the yield and the Doctor line

## Overview

The configuration guide and the Roundfix `runtime` reference list three
Light Tier skip reasons and no pre-work yield, the Doctor guide and the
Roundfix `setup` reference describe no `light-tier:` line, and the glossary's
**Light Tier** entry says nothing about access. This Task describes the
`access_mode_unavailable` skip, the `light_tier_yielded` yield and the
`light-tier:` Doctor line in each, as `_techspec.md` states them for task_02
and task_03 to build, raises the Roundfix skill's version, and
re-records the Behavior Surface Record last (`_prd.md` → Core Feature 4).

## Requirements

1. MUST, in `docs/user-guide/configuration.md`, section "Light implementation
   tier":
   - replace "names one of the three reasons: the key variable is not set,"
     with "names one of the four reasons: the light runtime offers no mode for
     the requested access policy, the key variable is not set,";
   - add, right after the paragraph that ends "or this month's light spend
     reached the ceiling.", this paragraph: "When `defaults.agent_full_access`
     or `--agent-full-access` requests full access and no light runtime offers
     a full-access mode, as OpenCode offers none, the light tier is skipped
     before any other check with the reason code `access_mode_unavailable`,
     and the Task runs on its category profile with the same access request.
     A light session refused for any other reason before Agent work, such as
     a failed setup or a prompt that fails before any Agent output, yields:
     Roundfix prints `roundfix: warning: light tier yielded for Task <id>
     before Agent work: <cause>; it runs on its <category> profile`, publishes
     the Task event phase `light_tier_yielded`, and opens the category's
     Preferred Selection without trying the remaining light models.
     `roundfix doctor` warns about the first case on its `light-tier:` line
     (ADR-0314)." The ADR reference takes the guide's link form, a relative
     link to `../adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`,
     as the guide's other ADR references do.
2. MUST, in `docs/user-guide/commands/doctor.md`:
   - add, right after the `recommendations:` bullet that ends "This line
     opens no Agent Session and never fails Doctor.", this bullet: "-
     `light-tier:` — printed after `recommendations:` only when
     `openrouter.light_models` is not empty and `defaults.agent_full_access`
     is `true`. It builds each light model's runtime as a Run would and
     reports `DR-LIGHT-ACCESS-MODE` (`warn`) when none offers a full-access
     mode, because light Tasks then run on their category profile; otherwise
     it reports `ok`. It opens no Agent Session, reads no key, and never
     changes the exit code.";
   - add to the finding table, right after the `DR-LANE-IMAGE-UNKNOWN` row,
     the row "| `DR-LIGHT-ACCESS-MODE` | warn | No light runtime offers a
     full-access mode; set `openrouter.light_models: []` in User Config, or
     `defaults.agent_full_access: false`. |".
3. MUST, in the canonical Roundfix reference
   `.agents/skills/roundfix/references/runtime.md`, section "Light
   implementation tier", replace "A skipped light Task emits phase
   `light_tier_skipped` with reason code `key_missing`, `spend_unreadable`,
   or `ceiling_reached`." with "A skipped light Task emits phase
   `light_tier_skipped` with reason code `access_mode_unavailable`,
   `key_missing`, `spend_unreadable`, or `ceiling_reached`. The
   `access_mode_unavailable` check comes first: under a full-access request,
   a light model whose runtime offers no full-access mode is not chosen, and
   when none remains the Task runs on its category profile. Any other refusal
   of a light session before Agent work emits phase `light_tier_yielded`,
   naming the light model, a reason code and the cause, and opens the
   category's Preferred Selection without trying the remaining light models."
4. MUST, in the canonical Roundfix reference
   `.agents/skills/roundfix/references/setup.md`:
   - replace "the conditional `linux-lane` check" with "the conditional
     `light-tier` and `linux-lane` checks";
   - add, right after the `DR-LANE-IMAGE-UNKNOWN` row, the row "|
     `DR-LIGHT-ACCESS-MODE` | warn | No light runtime offers a full-access
     mode; set `openrouter.light_models: []` in User Config, or
     `defaults.agent_full_access: false`. |";
   - add, right before the `linux-lane:` bullet, the bullet "- `light-tier:`
     — printed after `recommendations:` only when `openrouter.light_models` is
     not empty and `defaults.agent_full_access` is `true`; `warn` with
     `DR-LIGHT-ACCESS-MODE` when no light runtime offers a full-access mode,
     otherwise `ok`. It opens no Agent Session and never changes the exit
     code."
   No text inside any `### QA settlement` section of any skill may change.
5. MUST revise **Light Tier** in `CONTEXT.md` through `domain-modeling` by
   appending to its definition, after "escalates its one feedback turn to its
   category's Preferred Selection.", the sentence "A light runtime that
   offers no mode for the requested access policy is not chosen, and a light
   session refused before Agent work yields to the category's Preferred
   Selection (ADR-0314)." The ADR reference takes the glossary's link form,
   to `docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`.
   Its `_Avoid_` line keeps its bytes.
6. MUST run these commands, because this Task edits an owned skill and two
   guides. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `SKILL.md` of the Roundfix skill by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited. Re-record after
   the last guide edit: the Behavior Surface Record command runs last, after
   every edit to `docs/user-guide/**` and the references, and no guide or
   reference is edited after it.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     configuration and Doctor guides' and the references' fingerprints
     change.
7. MUST search the guides, the agent guides and the skill copies for wording
   that says a light Task is skipped for only three reasons, that a light
   session's refusal fails the Task, or that lists every Doctor line without
   `light-tier:`, and report what it found and what it changed.
8. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page.

## Subtasks

- [ ] Describe the skip, the yield and the Doctor line in the configuration and Doctor guides.
- [ ] Describe them in the Roundfix `runtime` and `setup` references.
- [ ] Revise **Light Tier** in the glossary.
- [ ] Raise the skill version, sync the mirrors, record the version, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both guides and both copies of each reference name `access_mode_unavailable`, `light_tier_yielded` or the `light-tier:` line as their requirement states.
- [ ] `CONTEXT.md` revises **Light Tier** and links ADR-0314.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirrors, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/configuration.md .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'access_mode_unavailable' || { printf 'missing skip reason in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'light_tier_yielded' || { printf 'missing yield phase in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'names one of the four reasons: the light runtime offers no mode for the requested access policy' || { echo 'missing four reasons in configuration guide' >&2; exit 1; }; for f in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'light-tier:' || { printf 'missing light-tier line in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'DR-LIGHT-ACCESS-MODE' || { printf 'missing finding code in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'the conditional .light-tier. and .linux-lane. checks' || { printf 'missing check list in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'is not chosen, and a light session refused before Agent work yields to the category' || { echo 'missing Light Tier revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document names `access_mode_unavailable` and the glossary's **Light Tier** says nothing about access, so the command fails at its first phrase check. Afterwards the configuration guide and both `runtime` reference copies name the skip reason and the yield phase, the Doctor guide and both `setup` reference copies describe the `light-tier:` line and its finding code, the glossary revises **Light Tier**, the mirrors are in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 4; Goals; Glossary
- `_techspec.md` → API Contract 6; Build Order 1
- ADR-0314

## Result

- Updated the configuration and Doctor guides with the `access_mode_unavailable`
  skip, the `light_tier_yielded` pre-work yield, the `light-tier:` Doctor line,
  and `DR-LIGHT-ACCESS-MODE`.
- Updated the canonical Roundfix `runtime` and `setup` references, then ran
  `make skills-sync`; the canonical and embedded Roundfix skill copies are
  byte-identical.
- Revised the **Light Tier** glossary definition with the ADR-0314 link while
  preserving its `_Avoid_` line.
- Raised both Roundfix skill version fields from `0.1.130` to `0.1.131` and
  recorded that version with `go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed,
  `1 passed in 1 packages`.
- Re-recorded `docs/references/behavior-surfaces.json` last with `go test
  -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` —
  passed, `1 passed in 1 packages`.
- Focused checks after the final edit: required-term search found the new terms
  in all requested guide/reference copies; mirror comparison passed; `git diff
  --check` passed. The wording sweep found no remaining targeted Light Tier
  statement limited to three skip reasons, treating a refusal as a Task failure,
  or omitting `light-tier:` from the updated Doctor descriptions. Remaining
  generic refusal/failure matches are in unrelated history, delivery,
  profiles, implementation, and workflow documents.
- The Daemon-owned `## Verification` commands and repository Verification were
  not run.
