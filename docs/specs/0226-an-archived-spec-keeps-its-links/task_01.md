---
task: task_01
spec: 0226-an-archived-spec-keeps-its-links
status: completed
type: docs
complexity: low
---

# Task 01: The Roundfix Skill and the archive guide describe the link rewrite

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task describes, in the Roundfix Skill's
archive reference and the archive command guide, the behavior task_02 and
task_03 implement, as the TechSpec states it: the rewrite of relative links
that leave the Spec, the keep rule for an already archived target, the
refusal, the count on the confirmation line, and the limits.

## Requirements

1. MUST describe in `.agents/skills/roundfix/references/archive.md`, after the
   confirmation line example of the Archive Command section, and in
   `docs/user-guide/commands/archive.md`, after the paragraph that names the
   archive destinations: that before the move the archive rewrites each
   relative Markdown link that leaves the Spec (inline links, images and
   reference definitions outside code blocks and code spans) to reach the same
   path from the archived location, keeping its fragment, query and
   angle-bracket form; that it keeps a link whose target was already archived
   when the unchanged destination reaches it from the archived location; that
   when any other link leaves the Spec and reaches nothing it exits `2` before
   changing any file, naming the relative links that leave the Spec and do not
   resolve with file, line and destination; that a successful archive that
   rewrote links appends `; rewrote <n> relative link(s)`; and that
   destinations inside the Spec, absolute destinations, URLs, HTML anchors and
   non-Markdown files, including evidence scripts that climb a fixed number of
   directories, are not rewritten. This answers the Backlog Entry of
   2026-10-03, "Archiving a Spec breaks its relative links that leave the
   Spec".
2. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
3. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, `CONTEXT.md` or `.roundfixrc.yml`.

## Subtasks

- [ ] Describe the rewrite, the keep rule, the refusal, the count and the limits.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The archive reference and guide carry "rewrites each relative Markdown
      link that leaves the Spec", "relative links that leave the Spec and do
      not resolve", "whose target was already archived" and the count suffix.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/archive.md`
- instruction: `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- "rewrites each relative Markdown link that leaves the Spec" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/archive.md "rewrites each relative Markdown link that leaves the Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- "relative links that leave the Spec and do not resolve" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/archive.md "relative links that leave the Spec and do not resolve" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- "whose target was already archived" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/archive.md "whose target was already archived" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- "rewrote <n> relative link(s)" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/archive.md "rewrote <n> relative link(s)" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- "rewrites each relative Markdown link that leaves the Spec" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md "rewrites each relative Markdown link that leaves the Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- "relative links that leave the Spec and do not resolve" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md "relative links that leave the Spec and do not resolve" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- "whose target was already archived" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md "whose target was already archived" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- "rewrote <n> relative link(s)" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md "rewrote <n> relative link(s)" >&2; exit 1; }` — expected: exit 0; before this Task none of these phrases is in these files, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/archive.md | grep -qF -- "rewrites each relative Markdown link that leaves the Spec" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirror does not carry the new phrase, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → Core Feature 7; User Stories 1-2
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Surface Transcript 2; Build Order 1
- ADR-0187; ADR-0189; ADR-0230

## Result

Implemented the documentation slice for the archive link rewrite contract.
The Roundfix Skill archive reference and archive command guide now describe
rewriting outward relative Markdown links, preserving fragments, queries and
angle-bracket form; keeping links whose target was already archived; refusing
unresolved outward links before changing files with file, line and destination;
reporting `; rewrote <n> relative link(s)`; and excluding in-Spec, absolute,
URL, HTML-anchor and non-Markdown destinations, including fixed-depth evidence
scripts.

Raised both Roundfix Skill front-matter version fields from `0.1.25` to
`0.1.26`, synchronized the embedded mirror, recorded the version digest, and
ran the required baseline digest regeneration (no derived changes).

Focused checks completed:

- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- `make baseline-digests` — passed; reported no derived changes.
- `make skills-sync-check` — passed.
- `cmp` checks for the canonical and embedded Roundfix Skill files — passed.
- `git diff --check` — passed.
- The recorded version is `0.1.26` with digest `0bb47417be37f1a5c84d7c019f5331f7f611170d45a46e775105e2b0fe83460b`.

The task status remains daemon-owned as `in_progress`; the daemon must run the
authored Verification commands and settle the Task.
