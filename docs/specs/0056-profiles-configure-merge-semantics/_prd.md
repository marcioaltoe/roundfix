---
spec: 0056-profiles-configure-merge-semantics
status: active
created: 2026-07-28
surfaces: [backend, cli, docs]
---

# Profiles configure merge semantics

`roundfix profiles configure --scope project --file <fragment>` replaces the
entire `profiles:` map with the fragment, so configuring one category
silently deletes every other configured profile — a fragment naming only
`frontend` destroyed `general`, `backend`, `qa`, and `review` with their
Fallback Chains, reported `changed: true` naming only the requested
category, and rewrote the file's indentation so the 103-line diff read as
formatting. This contradicts the command's documented contract of preserving
unrelated config. A second defect compounds it: a declined confirmation in a
non-interactive context writes nothing and exits `0`, so automation reads a
refusal as success. Evidence:
[profiles configure replaces the whole profiles map](../../findings/2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md).

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier
  is created; profile categories and Agent Selection tuples keep their
  existing identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all behavior is local
  configuration reading, merging, confirmation, and writing; no
  authentication or HTTP surface. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 keeps each present profile
  atomic (a named category still replaces that category as one object; the
  merge is at map level, never field level); ADR-0037/0039 selection-proof
  obligations are unchanged — proof-before-write is preserved. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Configuring one category leaves every other configured category
  byte-identical.
- A destructive write can no longer be mistaken for a formatting diff: the
  effective change is summarized per category before confirmation.
- A single-value change produces a single-value diff.
- A declined non-interactive write is distinguishable from a successful one
  by exit code alone.

## User Stories

1. As a user configuring the `frontend` profile from a fragment, I want the
   fragment's named categories merged into the existing map, so that my
   other configured profiles survive untouched.
2. As a user reviewing the confirmation, I want a per-category summary —
   added, replaced, removed — so that a removal cannot hide inside
   reformatting churn.
3. As a maintainer reviewing the resulting config diff, I want the file's
   existing indentation, key order, and comments preserved with only the
   changed values touched, so that the diff is reviewable.
4. As an Agent driving the command non-interactively, I want a declined
   confirmation to exit non-zero, so that my automation cannot read a
   silent no-op as success.
5. As a user who genuinely wants to remove a category, I want removal to
   require naming it explicitly, so that deletion is always intentional.
6. As a maintainer running the command against my existing config, I want
   every path that works today to keep working, so that fixing the merge
   never costs me a working one.

## Core Features

1. Fragment semantics are merge-by-category: each category named in the
   fragment replaces that category atomically per ADR-0049; categories
   absent from the fragment are preserved byte-identically.
2. Removing a category requires an explicit removal declaration; a fragment
   can never delete by omission.
3. The confirmation renders a per-category effective-change summary —
   added, replaced, removed — before any write, in both interactive and
   `--json` flows.
4. The writer preserves existing indentation, key order, and comments,
   touching only the lines the change requires.
5. A declined or unconfirmable write exits non-zero with a distinct
   refusal, while `--dry-run` keeps its current successful-no-write
   contract; the `roundfix/profiles-configure/v1` schema evolves additively.
6. The proof-before-write behavior is unchanged: every distinct tuple in
   the resulting effective map is exact-proven before confirmation.
7. A characterization corpus captured from the current implementation
   before any behavior change bounds the new writer: every config the
   current writer handles round-trips through the new one with only the
   intended values differing, and none fails to parse afterward.

## User Experience

- The confirmation reads like a change plan: one line per affected
  category, nothing about untouched ones, and an explicit `removed:` line
  when removal was requested.
- A single-value edit shows a one-line summary and produces a minimal diff.
- Non-interactive decline prints the refusal on stderr and exits non-zero;
  `--yes` remains the explicit consent path.

## Non-Goals / Out of Scope

- Field-level merging inside a category — a configured category remains one
  atomic object per ADR-0049.
- Changing profile resolution precedence, fallback semantics, or proof
  behavior.
- A general YAML formatting or linting pass over the config file.
- Changing `profiles show` or `profiles validate`.

## Success Metrics

- Configuring one category leaves every other configured category
  byte-identical, proven against a five-category config.
- A fragment that would remove a category without declaring removal is
  refused; a declared removal appears in the confirmation summary.
- A single-value change produces a diff touching only that value's lines.
- A declined non-interactive invocation exits non-zero; automation
  distinguishes it from an applied one by exit code alone.
- Every config in the characterization corpus round-trips through the new
  writer with only the intended change, and none fails to parse after a
  write.
- The only observable behavior changes are the two breaks declared below;
  any other change to output, exit code, or file content is a defect.

## Decisions

- Merge is by category, replacement is within category — ADR-0049's atomic
  profile object is the merge unit.
- Removal is a first-class declared operation, never an omission side
  effect.
- The write path becomes format-preserving rather than re-serializing the
  whole document.
- A Spec evolves Roundfix and never regresses it. Behavior that works
  today is preserved unless this PRD names the break. Exactly two breaks
  are declared: a fragment no longer removes categories by omission —
  removal must be declared — and a declined non-interactive write exits
  non-zero instead of zero.
- The current writer's behavior is captured as a characterization corpus
  before the format-preserving writer replaces it; that corpus is the
  regression gate, not a test written after the fact.

## Open Questions

None.
