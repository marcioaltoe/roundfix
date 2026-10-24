---
spec: 0316-an-agent-may-start-the-setup-skill
status: active
created: 2026-10-24
surfaces: [docs]
---

# An Agent may start the setup skill

The `setup-context-driven` skill that Roundfix ships turns off model
invocation, so an Agent cannot start it and the Baseline's skill guide tells
the Agent to ask the person to run it, even for a managed refresh whose only
writes already wait for the person's approval of the Plan Digest. After this
Spec the skill no longer turns off model invocation: an Agent whose task
matches the skill's dispatch trigger may start it, and the skill and the
Baseline guide say that starting it writes nothing, because `roundfix baseline`
still writes only after the person approves the exact Plan Digest through
`--confirm-plan <digest>`, `--yes` or the interactive confirmation. The other
skills that keep the lock, and the skill guide's clause for them, are
unchanged.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-24-an-agent-may-start-the-setup-skill.md`](references/2026-10-24-an-agent-may-start-the-setup-skill.md),
   `type: feat`, from the maintainer's request relayed by Quaestor (ai-memory
   message `01a16dbe`, 2026-10-24): remove `disable-model-invocation: true`
   from the setup-context-driven skill and the matching block in skill
   dispatch, "preservando as confirmações do próprio comando para planos e
   operações quando aplicáveis".

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. This Spec is authored against `64c56d7f` (v0.117.0 plus
Spec 0311) and delivers after Spec 0315 and before Spec 0307, by the
maintainer's answer quoted in `_authorization.md`. Specs 0307, 0309, 0310
and 0314, delivered around it, also re-record
`skills/testdata/owned-skill-versions.json` or
`docs/references/behavior-surfaces.json`, and Spec 0309 also edits
`docs/user-guide/context-driven-development.md`; each Task re-records those
files on its own delivery base, so the order needs no `requires`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the skill keeps its name and gains a version. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — this Spec changes no Baseline
  module, clause or trigger; the skill guide's clause about skills that only a
  person can start keeps its text and keeps governing the skills that still
  turn off model invocation. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the
  `setup-context-driven` skill's canonical copy and its shipped mirror. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0316-an-agent-may-start-the-setup-skill/_authorization.md`;
  bounded files: `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`.

## Goals

- The canonical and shipped `setup-context-driven` skills carry no
  `disable-model-invocation` field, so an Agent can start the skill.
- Starting the skill authorizes no write: the skill tells the Agent to present
  the Plan and its Plan Digest and to confirm only a digest the person
  approved or that the repository's own instructions already authorize.
- The Baseline guide tells adopters the same.
- The skill's version rises and is recorded, so adopters receive it with the
  next release through `roundfix skills install`.
- The skill guide's person-only clause and every other skill that turns off
  model invocation are unchanged.

## Core Features

1. The skill metadata. `.agents/skills/setup-context-driven/SKILL.md` and its
   mirror `skills/setup-context-driven/SKILL.md` lose the
   `disable-model-invocation: true` line; both version fields rise from 0.0.5
   to 0.0.6, recorded in `skills/testdata/owned-skill-versions.json`.
2. The skill text. A new `## Who starts this skill` section, before `##
   Boundaries`, says an Agent may start the skill when its dispatch trigger
   matches, that starting it authorizes no write, and that the Agent passes
   `--confirm-plan <digest>` or `--yes` only after the person approves that
   Plan Digest or the repository's own instructions already authorize it.
3. The guide. `docs/user-guide/context-driven-development.md`, section "Adopt
   or update the Context-Driven Baseline", gains one paragraph saying the
   same, and the Behavior Surface record is re-recorded.
4. A contract test. `skills/setup_skill_invocation_test.go` asserts both
   copies carry no `disable-model-invocation` field and keep the Plan Digest
   approval, and that `implement-spec` still turns off model invocation.

## User Experience

In Claude Code, the skill's description now appears in the Agent's skill
listing. A request such as "atualize o baseline deste repositório" lets the
Agent start the skill, run `roundfix baseline update --repo . --format json`,
show the Plan and its Plan Digest, and wait. The person answers with the
approval; the Agent then runs `roundfix baseline update --repo .
--confirm-plan <digest> --format json`. Typing the skill's command by hand
still works.

## Non-Goals / Out of Scope

- Removing or rewording the skill guide's clause for skills only a person can
  start, or any Baseline module, trigger or Source Baseline row.
- Changing `roundfix baseline`'s confirmations: `--yes`, `--confirm-plan`,
  the interactive approval and the Plan Digest stay as they are.
- Changing `implement-spec` or any upstream skill that turns off model
  invocation (`grill-with-docs`, `handoff`); upstream skills are never edited.
- Changing the upstream skills repository's own copy of
  `setup-context-driven`; Roundfix installs its own.

## Success Metrics

1. Success Metric: `skills/setup_skill_invocation_test.go` passes: neither
   copy of the skill has a `disable-model-invocation` field, both carry the
   new section and the sentence "Starting it authorizes no write", the skill
   still holds "Ask the maintainer to approve the exact current
   `planDigest`", and `.agents/skills/implement-spec/SKILL.md` still holds
   `disable-model-invocation: true`.
2. Success Metric: `go test ./skills` passes with the skill at version 0.0.6
   recorded, and the canonical copy and its mirror are identical.
3. Success Metric: the documentation contracts pass with the Behavior
   Surface record re-recorded after the guide edit.
4. Success Metric: the rendered `docs/agents/skill-dispatch.md` and every
   Baseline module are byte-identical to the Delivery Base.
5. Success Metric: after release, in an adopter's Claude Code session, the
   Agent starts `setup-context-driven` without the person typing its
   command and stops at the Plan Digest.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-24.

- Claude Code documentation, "Extend Claude with skills"
  (https://code.claude.com/docs/en/skills), section "Control who invokes a
  skill": with `disable-model-invocation: true`, "Claude can't invoke the
  skill on its own" and the description is not in context; by default "both
  you and Claude can invoke any skill".
- The maintainer's request in the adopted Backlog Entry, written before this
  Spec existed, quoting "preservando as confirmações do próprio comando para
  planos e operações quando aplicáveis".

## Unreachable Acceptance

- criterion: Success Metric 5, after release, in an adopter's Claude Code
  session, the Agent starts `setup-context-driven` without the person typing
  its command and stops at the Plan Digest.
  reason: it needs an interactive Agent session in another repository with
  the released skill installed, which no Verification or QA row may start or
  read.
  satisfied-by: the operator's or an adopter's next managed refresh after the
  release, reported through the Secondbrain inbox.

## Validation

The decisions were checked on 2026-10-24 against these sources.

- Characterization at `64c56d7f`: `disable-model-invocation: true` appears in
  `.agents/skills/setup-context-driven/SKILL.md`, its mirror,
  `.agents/skills/implement-spec/SKILL.md` and its mirror, and the upstream
  `grill-with-docs` and `handoff`. No Go code reads the field: `rg -n -i
  "disable-model|model invocation"` over `internal` and `skills` finds only
  the Baseline's generic clause in `internal/baseline/assets/modules/core.json`
  and its two tests. The rendered `docs/agents/skill-dispatch.md` has no line
  about this skill beyond its trigger.
- A prototype in a scratch clone of `64c56d7f`: the line removed, the section
  added, both versions raised to 0.0.6, `make skills-sync`, the version
  record re-recorded; `go test ./skills` passed. The guide paragraph made
  `TestTheSkillCoverageMapIsCurrent` report a stale fingerprint until it was
  re-recorded; then the documentation contracts passed. `go test
  ./internal/baseline` passed with no golden change.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query
  "disable-model-invocation skill setup-context-driven baseline update agent"
  --all --files --min-score 0.3`, which returned the upstream skills mirror's
  copy of the skill (also locked), an AI Hero note on skills and memory, a
  Roundfix ADR on upstream setup names and Quaestor's skill guide. None
  records a decision to keep the lock. No Inbox Entry was written.
- Exa: the Claude Code skills documentation quoted under Acceptance evidence.

## Glossary

None.

## Decisions

- 2026-10-24 (about 11:15 BRT): the maintainer chose "Tirar a trava da skill
  (Recomendado)": remove the lock from the skill and the block in skill
  dispatch; the command's own confirmations (Plan Digest, `--confirm-plan`)
  stay; a small Spec with or right after Spec 0315.

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 5 briefing for Spec 0316 of 2026-10-24 and the characterization above.

Taken from those sources, with no further interview:

- The "block in skill dispatch" is the generic clause that tells an Agent to
  ask the person to run a skill whose metadata turns off model invocation. It
  names no skill; removing the field takes `setup-context-driven` out of its
  reach, so the clause stays for `implement-spec`, `grill-with-docs` and
  `handoff`. No Baseline module changes.
- The skill states the confirmation rule in its own section rather than
  changing its existing apply recipes.
- No new ADR: no Baseline clause or architecture decision changes.

## Open Questions

None.
