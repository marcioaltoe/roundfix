# Skill governance

Who owns each skill in this repo, the roundfix skill-sync contract, and where workflow rules
live. `AGENTS.md` keeps only short mandatory pointers to these rules; the bodies below are
canonical.

## Rule placement

HARD RULEs and workflow orientation for the CONTEXT-driven spec workflow and for Roundfix live
in the skills themselves and in `docs/agents/*.md` files — this one,
`docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`, `docs/agents/issue-tracker.md`,
`docs/agents/domain.md`. `AGENTS.md` holds a one-line pointer that marks each rule as mandatory,
never the full rule body. When a rule outgrows its pointer, move the body here or into the
matching skill and leave the pointer behind.

## Roundfix skill sync (pre-PR gate)

Before opening any PR, confirm `.agents/skills/roundfix/SKILL.md` still matches the shipped CLI
behavior: commands, flags, output formats, exit codes, Batch contract semantics. If the PR
changes any of those, the skill update ships in the same PR.

- The project-local copy at `.agents/skills/roundfix/` is canonical.
- Its `metadata.version` tracks the released CLI version (the `v*` tag), not an independent
  skill version.
- The embedded `skills/roundfix/` is generated from it with `make skills-sync`; `make verify`
  fails on drift.
- Never keep a detached global copy (for example `~/.claude/skills/roundfix/`) — a stale copy
  shadows the canonical one and teaches agents a contract that no longer exists (observed
  2026-07-15, findings retrospective §2). If a machine-global install is wanted, symlink it to
  this repository's `skills/roundfix/`.

## Skill ownership

The authorial context-workflow skills are vendored in this repository and owned by it:
`write-idea`, `write-prd`, `write-techspec`, `write-tasks`, `setup-workflow`, `implement-task`,
`implement-spec`, `brainstorming`, `council`, `business-analyst`, `archive-spec`, `qa-gate`,
and `evidence-gate`.

- They live in `.agents/skills/` as plain repo files and are **never** listed in
  `skills-lock.json`.
- They may be adapted to Roundfix's needs — changes ship like any other repo change.
- Every other skill (including `grilling`, `grill-with-docs`, `the-fool`, `domain-modeling`,
  and all generic engineering skills) remains managed by the external `marcioaltoe/skills`
  origin through `skills-lock.json` and **must not** be modified locally — needed changes go
  upstream.
