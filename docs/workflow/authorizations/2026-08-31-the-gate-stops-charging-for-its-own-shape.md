---
granted: 2026-08-12
action: Carry the characterization obligation and the Pull Request row's equivalent-evidence default into the authoring and gate skills for Spec 0105.
consuming: 0105-the-gates-own-economics
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/implement-task/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/implement-task/SKILL.md
---

# Tooling authorization — the gate stops charging for its own shape (2026-08-12)

The maintainer granted this on 2026-08-12, recorded at
`docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
which names Spec 0105 among its eight consuming Specs. This record re-bounds the
same grant to the three skills Spec 0105 actually edits, and adds the generated
copies the audit needs named.

## Why a second record rather than a citation

The 2026-08-12 record lists eight canonical skills under `.agents/skills/` and no
generated copy under `skills/`. The changed-path audit resolves a grant from
`paths:` frontmatter, so a Task that runs `make skills-sync` — which it must,
because a hand-edited generated copy is an unauthorized mutation — changes a path
no grant names. Spec 0118 was refused on 2026-08-27 for exactly that, with its
record describing the mirror as sanctioned fallout in prose the audit does not
read.

ADR-0149 says a grant should name only the regeneration command and let the audit
resolve its outputs from an `_ownership.yml` declaration, which is the shape
ADR-0081 prefers over enumerating consequences. That mechanism does not reach
here: `skills/` carries no such declaration, and the resolver reads records only
under `internal/baseline`. Until that is closed, enumerating is the only form the
audit can act on. The gap is recorded for Triage; this record works within it
rather than pretending it is absent.

## What this covers

Spec 0105 measured 123 of 201 failed Tasks across five repositories as the QA
gate returning a verdict rather than code breaking. Two of its five features are
guidance:

- the task-authoring skill gains the characterization obligation, so a Spec
  crossing an external boundary records what the real thing does before a premise
  the boundary does not support reaches the gate four Runs later;
- the QA gate skill gains the Pull Request row's equivalent-evidence default,
  because that row is unreachable by construction under ADR-0088 and one Spec
  paid six of its eight gate executions for it alone.

The task-execution skill is named because the characterization obligation reaches
the Agent that performs such a Task.

## Authorized paths

- `.agents/skills/write-tasks/SKILL.md`, limited to the characterization
  obligation and its ordering against the work that depends on it.
- `.agents/skills/qa-gate/SKILL.md`, limited to the Pull Request row's default
  and the evidence it must record.
- `.agents/skills/implement-task/SKILL.md`, limited to what a characterization
  Task requires of the Agent performing it.
- `skills/qa-gate/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/implement-task/SKILL.md` — the generated copies, rewritten by the
  declared `make skills-sync` and named here only because the audit cannot
  otherwise resolve them. They remain sanctioned fallout under ADR-0081; a
  hand-edited copy is still an unauthorized mutation.

## Bounded by purpose

This grant covers the two guidance changes above. It does not authorize changing
the gate's verdict rules, its row contract, the typed blocked-cause counts, which
checks it runs, or any other skill in the repository.

## Consuming Spec

This authorization is consumed by Spec `0105-the-gates-own-economics`.

## Commit choreography

This record lands as its own commit, before any commit that edits a skill.
