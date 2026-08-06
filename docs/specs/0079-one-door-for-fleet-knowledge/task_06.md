---
task: task_06
spec: 0079-one-door-for-fleet-knowledge
status: pending
type: chore
complexity: medium
---

# Task 06: Teach the authorial skills to start at the door

## Overview

Tooling task two of two, under the same authorization as task_02, and the
moment obligation becomes legitimate: the pilot (task_05) proved the door,
so the skills and guides may now require it. `write-idea` and `write-prd`
gain the inbox-first step, and the two modules gain the obligating clauses —
the layout guide's inbox-flow rule (a closed finding spawns its typed
Backlog Entries; new fleet knowledge enters through the door) and the
`capture: auto` contract for the maintainer's session-end hook (drafts are
always pending triage and never self-triaged). After Spec 0073, the
owned-skill content edits need only the mirror sync; the module edits carry
the digest chain.

## Requirements

1. MUST add the inbox-first step to `.agents/skills/write-idea/SKILL.md` and
   `.agents/skills/write-prd/SKILL.md`: before exploration, read the
   repository's pending inbox (the brain-side namespace for this
   repository) so queued knowledge feeds the decision; regenerate the
   `skills/**` mirrors with `make skills-sync`.
2. MUST author the obligating flow clauses in
   `internal/baseline/assets/modules/context-workflow.json`: fleet
   observations enter through the Secondbrain inbox — never as loose files
   in project checkouts — and a Finding whose lifecycle closes spawns the
   typed Backlog Entries its actions call for.
3. MUST author the session-end capture contract in
   `internal/baseline/assets/modules/secondbrain.json` with these literal
   tokens: `capture: auto` drafts are always pending triage and
   `never self-triaged`; the hook itself is the maintainer's, outside this
   repository — this clause contracts only what the hook writes.
4. MUST author the research-capture obligation: a session that performed
   substantive external research captures its digest with sources for the
   brain's namespace, advisory qmd duplicate check first.
5. MUST fold the pilot's recorded friction observations into the clause
   wording where they demand it; the pilot report is the evidence the
   obligations stand on.
6. MUST run the module chain per the Spec 0075 choreography (manifest-row
   bootstrap, two-pass `make baseline-digests`) and adopt both regenerated
   guide postimages.
7. MUST change only the authorization's bounded files, their sanctioned
   digest fallout (ADR-0081), and this task file.

## Subtasks

- [ ] Author the inbox-first step in both skills; sync mirrors.
- [ ] Author the obligating clauses in both modules; run the chain.
- [ ] Adopt both guide postimages.

## Acceptance Criteria

- [ ] Both skills open their exploration with the pending-inbox read, and
      their mirrors match.
- [ ] Both guides carry the obligating clauses with the literal tokens from
      Requirements 3.
- [ ] The digest chain is converged and the diff stays inside the bounded
      files, fallout, and this task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-fleet-knowledge-door.md
- interface: .agents/skills/write-idea/SKILL.md
- interface: .agents/skills/write-prd/SKILL.md
- interface: internal/baseline/assets/modules/context-workflow.json
- interface: internal/baseline/assets/modules/secondbrain.json

## Verification

- `grep -qi 'inbox' .agents/skills/write-idea/SKILL.md && grep -qi 'inbox' .agents/skills/write-prd/SKILL.md`
  — expected: exit 0; both canonical skills carry the inbox-first step.
- `make skills-sync-check`
  — expected: exit 0; the embedded mirrors match the canonical skills.
- `grep -q 'capture: auto' docs/agents/secondbrain.md && grep -q 'never self-triaged' docs/agents/secondbrain.md`
  — expected: exit 0; the session-end capture contract binds with its
  literal tokens.
- `grep -qi 'inbox' docs/agents/docs-layout.md`
  — expected: exit 0; the layout guide carries the inbox-flow rule.
- `s1="$(git status --porcelain | sort)"; make baseline-digests >/dev/null 2>&1; s2="$(git status --porcelain | sort)"; [ "$s1" = "$s2" ]`
  — expected: exit 0; the digest chain is converged.
- `output="$(git status --porcelain | awk '{print $NF}' | grep -vE '^(\.agents/skills/write-idea/|\.agents/skills/write-prd/|skills/|internal/baseline/assets/|internal/baseline/testdata/|docs/agents/docs-layout\.md$|docs/agents/secondbrain\.md$|docs/agents/setup-context\.json$|docs/specs/0079-one-door-for-fleet-knowledge/task_06\.md$)')"; [ -z "$output" ]`
  — expected: exit 0; nothing outside the bounded files, mirrors, and
  sanctioned fallout changed.

## References

- `_prd.md` → User Stories 4–6; Core Features 7–9; Decisions (session-end
  capture over recorded dissent; research capture).
- `_techspec.md` → Build Order 5; Decisions (permission precedes
  obligation; owned-skill edits need only the mirror sync).
- ADR-0095, ADR-0081.
