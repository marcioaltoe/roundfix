---
spec: 0316-an-agent-may-start-the-setup-skill
prd: _prd.md
created: 2026-10-24
---

# An Agent may start the setup skill — Technical Spec

## Executive Summary

The lock is one frontmatter line in the Roundfix-owned skill
`setup-context-driven`. Removing it from the canonical copy and the shipped
mirror lets an Agent start the skill; the Baseline's generic person-only
clause in the skill guide stops matching it without any change, because it
names no skill and keys on the metadata. The skill gains a short section that
keeps the person's approval at the Plan Digest, the Baseline guide gains one
paragraph saying the same, the skill's version rises to 0.0.6, and a contract
test pins the result. No Go production code, Baseline module or rendered
guide changes.

The trade-off accepted: an Agent now reads the skill's description in every
session of an adopter that installs it, a few hundred characters of skill
listing, in exchange for not needing the person to start it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no Baseline module, clause or
  trigger changes, as `_prd.md` → Project Constraints states. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the
  `setup-context-driven` skill's canonical copy and its shipped mirror.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0316-an-agent-may-start-the-setup-skill/_authorization.md`;
  bounded files: `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`.

## Current behavior

Read on 2026-10-24 at `64c56d7f`.

1. `.agents/skills/setup-context-driven/SKILL.md` (canonical, version 0.0.5
   in `metadata.version` and `version`) and its mirror
   `skills/setup-context-driven/SKILL.md` (byte-identical, written by `make
   skills-sync`) carry `disable-model-invocation: true` on line 4. Its
   sections are `Boundaries`, `Interpret composed guidance`, `Refresh an
   adopted repository`, `Human adoption or profile change`, `Project
   decisions and Spec constraints`, `Non-interactive planning for adoption
   or profile change`, `Apply the approved plan`, `Baseline Profiles`,
   `Repository Skill Set restoration`, `Canonical asset synchronization`,
   `Recovery` and `Completion`. `Apply the approved plan` already says "Ask
   the maintainer to approve the exact current `planDigest`."
2. `skills/testdata/owned-skill-versions.json` records 0.0.2 to 0.0.5 for
   the skill; `TestEveryOwnedSkillVersionIsRecorded` fails when a version is
   not recorded.
3. The Baseline module `core` renders
   `clause.core.ask-the-person-to-start-a-person-only-skill` into
   `docs/agents/skill-dispatch.md`: "When a matching skill can be started
   only by a person, because its metadata turns off model invocation, ask the
   person to run it instead of activating it or carrying out its workflow
   yourself." The module `context-workflow` dispatches
   `trigger.context-workflow.setup-context-driven`. Neither names the field
   or a skill; `internal/baseline/person_only_skill_test.go` pins the clause.
4. Other skills with the field: `.agents/skills/implement-spec/SKILL.md`
   (owned, and its mirror) and the upstream `grill-with-docs` and `handoff`.
5. The Baseline's Setup Snapshots install `setup-context-driven` from source
   `repo` `roundfix` with `minimumVersion` 0.0.2, so adopters receive the
   shipped copy; `roundfix skills install` writes it.
6. No guide says an Agent cannot start the skill: `rg -n -i "model
   invocation|setup skill"` over `docs/user-guide`, `.agents/skills/roundfix`
   and `CONTEXT.md` finds no such sentence.
7. A guide edit makes `TestTheSkillCoverageMapIsCurrent` report a stale
   fingerprint for that guide until
   `docs/references/behavior-surfaces.json` is re-recorded (prototype,
   2026-10-24); a skill edit alone did not.

## System Architecture

- `.agents/skills/setup-context-driven/SKILL.md` and its mirror: metadata,
  version and one new section.
- `skills/testdata/owned-skill-versions.json`: the 0.0.6 record.
- `skills/setup_skill_invocation_test.go` (new): the contract test.
- `docs/user-guide/context-driven-development.md` and
  `docs/references/behavior-surfaces.json`: the guide paragraph and its
  fingerprint.

## Implementation Design

### Interfaces

The new section, placed between the skill's introduction and `##
Boundaries`:

```markdown
## Who starts this skill

An Agent may start this skill itself when its dispatch trigger matches; its
metadata does not turn off model invocation. Starting it authorizes no write.
Present the Plan and its Plan Digest first, and pass `--confirm-plan <digest>`
or `--yes` only after the person approves that Plan Digest or the
repository's own instructions already authorize it.
```

The guide paragraph, after the second paragraph of "Adopt or update the
Context-Driven Baseline":

```markdown
An Agent may start the `setup-context-driven` skill itself when its trigger
matches; the skill does not turn off model invocation. Starting it writes
nothing: the Baseline Command still writes only after the person approves the
exact Plan Digest, through `--confirm-plan <digest>`, `--yes` or the
interactive confirmation.
```

### Data Models

The skill frontmatter after the change, with every other line unchanged:

```yaml
name: setup-context-driven
description: <unchanged>
metadata:
  ...
  version: 0.0.6
version: 0.0.6
```

### API Contracts

1. API Contract: neither copy of the skill has a frontmatter line whose key
   is `disable-model-invocation`; both copies are byte-identical; both
   version fields read `0.0.6`, recorded with a digest in
   `skills/testdata/owned-skill-versions.json`.
2. API Contract: the skill holds the `## Who starts this skill` section of
   Interfaces before `## Boundaries`, and keeps every other byte apart from
   the removed line and the two version values, so "Ask the maintainer to
   approve the exact current `planDigest`" and the `--confirm-plan` recipes
   remain.
3. API Contract: the guide holds the paragraph of Interfaces, and
   `docs/references/behavior-surfaces.json` is re-recorded after the guide's
   last edit.
4. API Contract: unchanged Baseline. Every file under
   `internal/baseline/assets/`, `docs/agents/skill-dispatch.md` and
   `docs/agents/setup-context.json` keeps its bytes; the person-only clause
   still governs `implement-spec`, `grill-with-docs` and `handoff`.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes; the
skill's text and metadata and one guide paragraph change.

### Invariants

```text
1. An Agent may start setup-context-driven; no copy of it turns off model invocation.
2. Starting the skill writes nothing; the Baseline Command writes only after the person approves the exact Plan Digest or the repository's instructions authorize it.
3. The person-only clause and every other skill that turns off model invocation keep their bytes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariant 1, task_01
- Goal 2, Core Feature 2, Success Metric 1 → API Contract 2, Invariant 2, task_01
- Goal 3, Core Feature 3, Success Metric 3 → API Contract 3, task_01
- Goal 4, Core Feature 1, Success Metric 2 → API Contract 1, task_01
- Goal 5, Core Feature 4, Success Metrics 1 and 4 → API Contract 4, Invariant 3, task_01 and task_02
- Success Metric 5 → the PRD's Unreachable Acceptance

## Absence Evidence

1. A test of the skill's invocation metadata — command: `grep -rn -e "disable-model-invocation" -e "Who starts this skill" skills internal --include=*.go`; result: no match at `64c56d7f`; nothing in Go reads or pins the field.

## ADR Alternatives

None. No ADR is listed: no Baseline clause or architecture decision changes.

## Integration Points

None. No external system, network or provider is reached; the tests read
the repository's own skill files.

## Testing Approach

- `skills/setup_skill_invocation_test.go` (new), `package skills`, with a
  `Suite:`/`Invariant:`/`Boundary` header, reading files relative to the
  package directory (`../.agents/skills/...` and `setup-context-driven/...`,
  as the package's other tests do), one top-level test calling `t.Parallel()`
  first:
  - `TestAnAgentMayStartTheSetupSkill`, subtests `canonical` and `shipped`:
    the frontmatter (between the first two `---` lines) has no line whose key
    is `disable-model-invocation`; the body holds `## Who starts this skill`
    before `## Boundaries`, the sentence "Starting it authorizes no write.",
    and "Ask the maintainer to approve the exact current `planDigest`";
    subtest `person-only-skills-keep-the-lock`: the frontmatter of
    `.agents/skills/implement-spec/SKILL.md` still has
    `disable-model-invocation: true`.
- Unchanged and run: `go test ./skills` (mirror, version record, version
  raise, owned-skill edit and baseline skill contract tests),
  `TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill` and
  `TestTheSkillGuidesStateThePersonOnlySkillClause` in
  `./internal/baseline`, and the documentation contracts with
  `-tags docscontract`.
- A prototype passed these at authoring time; restoring the field in the
  canonical copy fails the `canonical` subtest, and restoring it in the
  mirror only fails the `shipped` subtest and the mirror test.

## Build Order

1. The skill change, its sync, version record, contract test, guide
   paragraph and the Behavior Surface record re-recorded last (task_01; no
   dependency).
2. QA gate (task_02) (depends on: 1).

No Spec is listed in `requires`.

## Risks & Considerations

- An Agent might pass `--yes` on its own; the skill and the guide forbid it
  without the person's approval or a standing authorization in the
  repository's instructions, and the command's confirmation is unchanged.
- The upstream skills repository keeps its own locked copy; Roundfix's
  Setup Snapshots install the Roundfix copy, so adopters get the unlocked
  one. The upstream copy is the skills project's concern.
- Harnesses that ignore the field (Codex, OpenCode) behave as before.

## Glossary

None.

## Decisions

- Keep the generic person-only clause: it names no skill and still governs
  the skills that keep the lock.
- State the confirmation rule in a new section of the skill and one guide
  paragraph; change no recipe.
- No ADR and no Baseline module edit.
