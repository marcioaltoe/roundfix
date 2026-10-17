---
type: fix
status: promoted
created: 2026-10-17
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
reason: null
---

# One unnamed superseded ADR keeps an adopter's whole history and turns its check red

## Symptom

On 2026-10-17, `roundfix baseline update` with v0.84.0 refused the `adr`
unit in three of four adopters. In each case a single superseded ADR had a
successor that does not name it:

- vortex: ADR 0002, successor 0051.
- fluxus: ADR 0002, successor 0064.
- oraculum: ADR 0012, which has no accepted successor that names it.

The rest of the plan applies with `--yes`, and the Managed Refresh writes the
retirement rule into the docs-layout guide. That turns `SC-HISTORY-RETAINED`
on while `docs/history/adr/` still holds every retired ADR:

- vortex: 18 errors;
- oraculum: 38 errors;
- fluxus: 22 files, reported as soon as a Spec is active.

The check also reports links from `docs/references/` to files of Specs that
an earlier layout removed:

- vortex: 3 links;
- oraculum: 2 links.

No unit rewrites those links, because `record-links` covers Archive Records
only. A repository that applies the update cannot pass `roundfix spec check`
until someone edits by hand.

## Where

- `internal/spec/history_entries.go`, `PlanHistoryKind` for the `adr` kind:
  one refused ADR refuses the unit.
- The `record-links` unit, which covers Archive Records only.
- The `SC-HISTORY-RETAINED` gate, which opens on the guide sentinel and not
  on a history the update left current.

## Expected

- An ADR the plan cannot retire is reported by name, and the unit still
  removes the others.
- A Live Document's link to a missing file gets its Git Locator through the
  same rewrite as an Archive Record's.
- Alternatively, the guide that turns the check on is held back until the
  history applies without a refusal.

## Evidence

The adopter dry runs of 2026-10-17. Scratch clones of fluxus and oraculum
were updated with `--yes` and checked with `roundfix spec check --strict`.
vortex was updated on a local branch that was held unpushed. conexus had no
refusal and passed.
