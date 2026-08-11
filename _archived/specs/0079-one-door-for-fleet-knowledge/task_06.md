---
task: task_06
spec: 0079-one-door-for-fleet-knowledge
status: completed
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

## Result

### Implementation

- Made the pending Inbox Entry read the first process step in both authorial
  Skills. Each step reads the repository's brain-side
  `inbox/<repository>/` namespace before exploration, treats queued evidence
  and intent as decision inputs, and leaves Triage to the destination
  repository.
- Added the substantive-research capture step to `write-idea`: capture a
  sourced digest for the brain's namespace, run the advisory qmd duplicate
  check first through an authorized access path, verify returned paths, and
  review substantive overlap instead of deciding from score alone.
- Versioned the context-workflow module, its docs-layout guide, and its layout
  rule. The new fleet-flow clause routes every new observation through the
  Secondbrain inbox instead of loose checkout files and requires a closed
  Finding to mint each typed Backlog Entry its recorded actions call for.
- Versioned the secondbrain module, its guide, and its edited rule. The
  maintainer-owned session-end hook is explicitly outside this repository;
  the local contract binds only its output: every `capture: auto` draft is
  always pending triage and never self-triaged. The research clause carries
  the pilot's access, stale-path, substantive-review, and score-only friction
  into the binding workflow.
- Ran `rtk make skills-sync`; it exited 0 and regenerated only the two owned
  `skills/**` mirrors.
- Bootstrapped marker-delimited Source Baseline entries and manifest rows for
  `clause.context.inbox-02-fleet-flow`,
  `clause.secondbrain.inbox-auto-capture`, and
  `clause.secondbrain.research-capture`. Regeneration replaced every temporary
  span and digest with values calculated from the source bytes.
- Ran the required two-pass module chain. The first
  `rtk make baseline-digests` exited 0 with `changed:true`; the second exited 0
  with `changed:false`. Adopted both regenerated managed guide blocks without
  changing repository-authored content outside those blocks.

### Focused checks

- The pre-change focused search found none of `pending Inbox Entries`,
  `never self-triaged`, `substantive external research`, or
  `loose capture files` on the six canonical Skill/module/guide surfaces.
- Exact `cmp -s` checks passed for both canonical Skill/mirror pairs and for
  each adopted managed guide block against its regenerated formatter golden.
- Focused `jq -e` assertions passed for the module, guide, and rule version
  bumps; the fleet-flow semantics; the literal `capture: auto` and
  `never self-triaged` tokens; the hook ownership boundary; and the
  research-capture access, path-existence, substantive-review, and score
  boundaries.
- Three exact `cmp -s` checks passed between the new module guidance strings
  and their marker-delimited Source Baseline entries. A focused manifest
  assertion found no all-zero digest or invalid span after regeneration.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk go test
  ./internal/baseline -count=1 -run
  'TestAuthorialSkillSync|TestCatalogCompatibility|TestFormatterComposition'`
  exited 0 and reported two passing tests in one package.

### Acceptance-criterion evidence

1. Both Skills place `Read the pending inbox` at process step 1, before their
   exploration steps. Exact comparisons prove both regenerated mirrors match
   their canonical Skills byte-for-byte.
2. Both adopted guide blocks match their regenerated goldens. The layout guide
   carries the fleet-door and closed-Finding obligations; the Secondbrain
   guide carries `capture: auto`, `never self-triaged`, hook ownership, and
   sourced research capture with the advisory qmd safeguards measured by the
   pilot.
3. The second digest pass reported `changed:false`. The 24 changed paths are
   the two authorized canonical Skills and two mirrors, two authorized
   modules, two authorized guide postimages, three required Source Baseline
   bootstrap surfaces, twelve deterministic digest/formatter/catalog/plan
   artifacts, and this Task file. No path lies outside the recorded grant,
   its required regeneration choreography, or ADR-0081 fallout.

### Handoff boundary

- The Daemon-owned commands under `## Verification` were not run in this Agent
  turn. Task status remains Daemon-owned; no commit, push, or pull request was
  created.
