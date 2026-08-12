---
task: task_07
spec: 0088-a-third-runtime-that-can-run
status: completed
type: docs
complexity: low
---

# Task 07: Move the measured OpenCode facts upstream

## Overview

The adopted measurement lives inside this Spec, and an archived Spec may be
deleted at any time. The durable half of what it recorded — which models the
subscription grants, how the effort vocabulary varies by model, and why Roundfix
treats OpenCode reasoning as model-managed — belongs to the repository's model
selection reference, which `.roundfixrc.yml` already cites.

## Requirements

1. MUST record, in the repository's model selection reference, the OpenCode model
   identifiers Roundfix can select, distinguishing the subscribed `opencode-go`
   tier from the pay-per-use and aggregator tiers.
2. MUST record that the reasoning-effort vocabulary is advertised per model and
   varies, with the measured examples, and that Roundfix therefore treats
   OpenCode reasoning as model-managed.
3. MUST date the snapshot and name the runtime and adapter versions it was taken
   against, matching the existing snapshot convention in that document.
4. MUST reference the decision rather than restate its reasoning, so the two
   documents cannot drift into disagreement.
5. MUST NOT reference any path under the Spec directory, because durable
   knowledge never depends on a Spec that may be deleted.
6. MUST NOT change `CONTEXT.md`; a glossary gap, if any, is raised by the closing
   gate.

## Subtasks

- [x] Add the OpenCode section to the model selection reference.
- [x] Record the subscribed tier and the measured per-model effort variation.
- [x] Date the snapshot with its runtime and adapter versions.
- [x] Confirm no Spec path is referenced.

## Acceptance Criteria

- [x] The model selection reference names the `opencode-go` subscribed models.
- [x] It records that the effort vocabulary is per model, with the measured
      examples.
- [x] It carries a dated snapshot line naming the runtime and adapter versions.
- [x] It cites the model-managed reasoning decision by its ADR identifier.
- [x] It contains no path under the Spec directory.

## Bounded scope

This Task may create or modify only:

- `docs/references/model-selection.md`
- `docs/specs/0088-a-third-runtime-that-can-run/task_07.md`

## Verification

- `grep -q 'opencode-go' docs/references/model-selection.md` — expected: exits 0.
- `grep -q 'ADR-0106' docs/references/model-selection.md` — expected: exits 0, proving the decision is cited rather than restated.
- `grep -q '2026-08-08' docs/references/model-selection.md` — expected: exits 0, proving the snapshot is dated.
- `grep 'docs/specs/' docs/references/model-selection.md` — expected: exits non-zero with no output, proving durable knowledge does not depend on a Spec.
- `go run -buildvcs=false ./cmd/roundfix spec check 0088-a-third-runtime-that-can-run` — expected: exits 0.

## References

- `_prd.md` → Core Features.
- `_techspec.md` → Build Order 7.
- `docs/agents/specific-repository.md` → the durable-knowledge-flows-upstream
  HARD RULE.
- ADR-0106.

## Result

The durable half of the adopted measurement now lives in the repository's model
selection reference, where `.roundfixrc.yml` already points.

**What changed.** The existing `The opencode runtime reaches everything else`
section was remeasured and extended. Its identifier count moved from the
2026-08-07 reading of 431 to the 2026-08-08 reading of 417, with the split
recorded — 339 `openrouter/`, 60 `opencode/`, 18 `opencode-go/` — and a note that
the number moves and should be read rather than quoted. Two subsections were
added. The first names the three tiers hiding behind one runtime, lists the
eighteen subscribed `opencode-go` models, and records that `gpt-5.6-luna` and
`deepseek-v4-flash` bill at 2x while `deepseek-v4-pro` does not. The second
records the per-model effort table, the fact that a session ensured without
`--model` advertises no `effort` option at all, and the mechanical reason
Roundfix declines to set effort on this runtime. The reasoning-effort sentence
above them was corrected: that vocabulary is the codex and claude adapters', not
OpenCode's. The document header now carries the remeasurement date with the
runtime and adapter versions.

**Commands and outcomes.**

- `grep -q 'opencode-go' docs/references/model-selection.md` — exit 0.
- `grep -q 'ADR-0106' docs/references/model-selection.md` — exit 0.
- `grep -q '2026-08-08' docs/references/model-selection.md` — exit 0.
- `grep 'docs/specs/' docs/references/model-selection.md` — exit 1, no output.
- `roundfix spec check 0088-a-third-runtime-that-can-run` — exit 0, no findings.
- `make verify` — exit 0 on a genuinely cold cache, zero `(cached)` lines.
- `git status --porcelain` — `docs/references/model-selection.md` and nothing
  else.

**Evidence per acceptance criterion.**

- Subscribed models named: the eighteen `opencode-go/` identifiers are listed
  verbatim.
- Per-model effort variation recorded: a seven-row table with the measured
  advertised values and defaults, plus the no-model case.
- Snapshot dated with versions: both the section heading and the document header
  carry 2026-08-08, opencode 1.18.15, acpx 0.13.0.
- Decision cited rather than restated: the reasoning paragraph points at
  ADR-0106 for the decision and keeps only the mechanical reason here.
- No Spec path: the grep found none, so this reference survives the Spec's
  deletion.

**Deliberately left out.** The measurement's Roundfix-side observations — the
failing `profiles validate` output and the doctor counts — stayed in the Spec.
They record what was broken on one day, not a durable fact about the runtime.
