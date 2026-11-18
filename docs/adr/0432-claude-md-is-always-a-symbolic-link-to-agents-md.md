---
status: accepted
created_at: 2026-11-18T09:50:00Z
updated_at: 2026-11-18T09:50:00Z
deprecated_at: null
superseded_by: null
---

# CLAUDE.md is always a symbolic link to AGENTS.md

In every repository the Baseline adopts, `AGENTS.md` is the one root
instruction file and `CLAUDE.md` is a relative symbolic link whose target is
exactly `AGENTS.md`. The repository has no option to keep a `CLAUDE.md` of its
own; its own rules go only in `docs/agents/specific-repository.md`.

## Context

On 2026-11-18 the maintainer set one standard for every repository that uses
the Baseline, with no option in the project: two separate files only if an
instruction had to differ per harness, otherwise `AGENTS.md` and a `CLAUDE.md`
that is always a link to it. No instruction differs: the Claude-only statement
of Spec 0405 opens with "In Claude Code,". Every adopted checkout on the
authoring host already keeps the link, but nothing in Roundfix creates it,
checks it or restores it. Adoption backs up a regular `CLAUDE.md` and leaves it
in place, a managed refresh never touches it, and a regular `CLAUDE.md` keeps
whatever Claude Code reads apart from what Codex and OpenCode read.

## Decision

**Adoption.** A `greenfield` or `preservation` plan writes `CLAUDE.md` as the
link whenever `AGENTS.md` is a regular file after the plan. It creates the link
when `CLAUDE.md` is absent and replaces a regular file or a safe link to
another file. The content-addressed root backup of ADR-0070 already keeps those
bytes, and a `preservation` plan already gives each of their Source Baseline
Entries a disposition.

**Managed refresh.** A `CLAUDE.md` that is already the link changes nothing.
An absent one is created. A regular file, or a safe link to another file, is
replaced by the link only when what Claude Code reads through it holds no
passage of its own: no Misplaced Repository Rule once a single `@AGENTS.md`
import line is set aside, or bytes identical to `AGENTS.md`. Otherwise the
plan withholds the conversion. It leaves `CLAUDE.md` byte-identical, names
each passage with `docs/agents/specific-repository.md` as its destination in
the warning `baseline.root.claude-link.withheld`, and refreshes everything
else as before. Once the passages are moved, the next update converts the
file. The conversion of a regular file writes the Baseline's content-addressed
backup `CLAUDE.<digest>.md`, as adoption does (ADR-0070), so no byte is lost
even though what it removes is only blank lines, the import line, Managed
Regions the Baseline renders into `AGENTS.md`, or a copy of `AGENTS.md`. A
link to another file needs no backup, because its target stays untouched. The
maintainer chose this on 2026-11-18: "Só depois de mover (Recommended)".

**Report.** `roundfix doctor` prints a `baseline-claude-link` line that is
`ok` for the link, `warn` with the next action for any other shape, and
`skipped` without a Setup Manifest. It never fails. The Baseline guide
`agent-instructions.md` states the rule.

**Unchanged.** An unsafe `CLAUDE.md` alias still blocks planning with
`baseline.inventory.unsafe-alias`. The link is not planned while `AGENTS.md`
is not a regular file, so the two files never point at each other.

## Considered Options

- **Convert a file with passages of its own and keep only a backup.**
  Rejected. The backup would hold the repository's own rules in a file no
  Agent reads, so the rules would silently stop applying. The only place for them is
  `docs/agents/specific-repository.md`, and a person moves them there.
- **Refuse the whole update until `CLAUDE.md` is clean.** Rejected. One file
  would block every catalog refresh of the adopter, the same reason ADR-0411
  rejected refusing an update over a Misplaced Repository Rule.
- **A decision that lets a repository keep a regular `CLAUDE.md`.** Rejected
  by the maintainer: "Sem opção ao usuário de mudar o baseline no projeto."
- **A copy of `AGENTS.md` where links do not work.** Rejected. A copy drifts
  from `AGENTS.md`, which is what the standard removes; ADR-0433 decides the
  checkout that does not create links.
- **Write the root blocks into a regular `CLAUDE.md` too.** Rejected. Two
  files that must stay equal are the shape the maintainer ruled out when no
  instruction differs per harness.

## Consequences

- This amends ADR-0100: a managed refresh may replace `CLAUDE.md`, and only
  that path, with the link when the preimage holds no passage of its own, and
  then writes the one root backup of the replaced regular file. The plan
  proves that condition on the bound preimage instead of the unchanged bytes
  outside the Managed Regions.
- This amends ADR-0411: a regular root `CLAUDE.md` whose only passage is the
  `@AGENTS.md` import line, or whose bytes are exactly `AGENTS.md`, reports no
  Misplaced Repository Rule, because the link replaces it.
- The file transaction replaces a regular file or a link at `CLAUDE.md` with
  the link and restores the preimage on rollback; every other path keeps the
  rule that a link postimage needs a missing preimage.
- A host that cannot create a symbolic link fails the apply before any file is
  replaced, and the transaction rolls back.
