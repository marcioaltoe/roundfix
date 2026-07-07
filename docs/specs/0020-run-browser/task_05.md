---
task: task_05
spec: 0020-run-browser
status: completed
type: docs
complexity: low
---

# Task 05: Docs and skill sync for the Run Browser surface

## Overview

Close the SKILL-matches-CLI gate: document the Run Browser, the enriched
`runs list` contract (`--state`, `--limit`, columns, notes), and the updated
attach wording everywhere the CLI surface is described, then re-sync the
embedded skill bundle. The CONTEXT.md Run Browser term already exists.

## Requirements

1. MUST update the README Commands and Command Boundaries for `runs`
   (browser), `runs list` (columns, flags, notes, Active default), and
   `attach` (browser entry, unknown-id wording), truthfully against
   implemented behavior.
2. MUST update the operational usage guide's monitoring flow: Run Browser
   for humans, bounded `runs list` for agents, including the Active-only
   default and the widening flags.
3. MUST update the canonical roundfix skill (SKILL.md) with the same
   contract — including that `runs list` defaults to Active Runs and agents
   widen with `--state`/`--limit` — and re-sync the embedded bundle.
4. MUST keep help text consistent with the docs.

## Subtasks

- [x] README Commands and Command Boundaries updates
- [x] Usage guide monitoring flow update
- [x] roundfix SKILL.md update and `make skills-sync`
- [x] Drift and skills checks pass

## Acceptance Criteria

- [x] README documents the browser keys, both flags, the Active default, the
      hidden-count notes, and the attach wording.
- [x] The usage guide routes humans to the Run Browser and agents to the
      bounded listing.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → API Contracts; Build
Order 5. CLAUDE.md SKILL.md-matches-CLI HARD RULE. CONTEXT.md → Run Browser.

## Result

Documented the Run Browser and the enriched `runs list`/`attach` contract
everywhere the CLI surface is described, and re-synced the embedded skill
bundle from the canonical skill.

What shipped:

- `README.md` Commands: the discovery section now leads with bare
  `roundfix runs` opening the Run Browser (header with repository and
  `ACTIVE`/`ALL` filter; `↑↓` move, `Enter` attach, `a` toggle,
  `q`/`Esc`/`Ctrl-C` quit exit `0`; cockpit-close returns to a refreshed
  browser; non-interactive exit `2` naming `runs list`), documents the full
  eight `runs list` columns, absolute UTC timestamps, `42m`/`1h12m` and
  `running <elapsed>` durations, `-` for missing fields, the 20-newest-
  Active default, `--state`/`--limit`/`--all`, both hidden-count note
  shapes, and the attach browser entry plus the unknown-id
  picker-numbers-are-not-stable wording. The stale `--active` flag and
  numbered-picker text are gone.
- `README.md` Command Boundaries: `runs list` boundary rewritten for the
  Active default, `--state`, `--limit`, the single trailing note, and the
  bare-`runs` browser/exit-2 split; `attach` boundary rewritten for the Run
  Browser loop (Active default, `a` widens, refreshed browser on return,
  cancel exits `0` with no side effects).
- `docs/usage.md`: the detached-Run monitoring flow now routes humans to
  the Run Browser (`roundfix runs` / `roundfix attach`) and scripts/agents
  to the bounded `runs list` with `--state all --limit 0` widening; the
  agent guidance bullet names the Active-only default and both widening
  flags and pins the Run Browser as the human surface; the command
  reference table gains the `runs` browser row.
- `.agents/skills/roundfix/SKILL.md` (canonical): "Run discovery and
  Attach" rewritten — 20-newest-Active default, eight columns, `--state`/
  `--limit` widening for agents, both note shapes, the Run Browser entry
  points with keys and empty-state text, non-interactive exit-2 contracts,
  and the unknown-id wording; the review-Run step 5, useful-commands list,
  spec-Run step 7, and detached monitor step 3 drop `--active`/picker
  language for the bounded listing and the Run Browser. `make skills-sync`
  regenerated `skills/roundfix/SKILL.md` byte-identical from the canonical
  file.
- Help text: unchanged in this Task — `runs`/`attach` usage already names
  the Run Browser, `--state`, and `--limit` (shipped with tasks 02/04);
  the docs above match it.

Acceptance evidence:

- README browser keys, flags, Active default, note shapes, attach wording:
  README.md Commands section (Run discovery block) and Command Boundaries
  bullets for `runs list` and `attach`; `rtk grep` over README.md,
  docs/usage.md, CONTEXT.md, and both SKILL.md copies finds zero remaining
  `--active`, numbered-picker, or Attach-picker references.
- Usage-guide routing: docs/usage.md monitoring flow ("At an interactive
  terminal, browse with the Run Browser; from a script or agent, use the
  bounded plain-text listing") plus the agent-guidance bullet ("The Run
  Browser is the human surface; agents stay on the plain-text listing").
- Drift and skills checks: `rtk make skills-sync-check` exits `0` after the
  re-sync, and `roundfix skills check` passes inside `rtk make verify`.

Verification:

- `rtk make verify` — pass: fmt-check, 978 tests in 19 packages,
  skills-sync-check, `roundfix skills check` (all 14 skills), and build.

Note: the first SKILL.md edit pass landed in the generated
`skills/roundfix/SKILL.md`; `make skills-sync` reverted it, exposing that
`.agents/skills/roundfix/SKILL.md` is the canonical source. The edits were
reapplied there and the bundle regenerated — both copies now match.
