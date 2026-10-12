---
type: fix
status: open
created: 2026-10-12
spec: null
---

# `baseline update` sanitizes history without reporting the citations it breaks

## Symptom

On 2026-10-09 the vortex ran `roundfix baseline update --format json` with
v0.65.0. The plan carried 79 Pending History units, `warnings: []`, and no
citation field. After the digest was applied, 2,518 files (50 MB) left
`docs/history`, and 46 links in 23 tracked Markdown files outside
`docs/history` stopped resolving. The broken links came from open Specs,
live backlog and findings, handoffs and `docs/references/`. Their targets were
files under `references/`, `evidence/` and `qa/` of the converted folders, and
those Specs' `_prd.md`. The vortex repaired them by hand.

Confirmed in code on 2026-10-12: `roundfix history sanitize` calls
`spec.HistoryCitations` (`internal/cli/history.go`) and prints
`cites <path>:<line> names <target>`. The detector already resolves relative
links and backticked paths (`internal/spec/history_entries.go`). The Pending
History path of `baseline update` (`internal/cli/baseline_update_history.go`)
never calls it, so the operator approves a digest that names no citation and
never receives the `--promote` offer. ADR-0254 says the update plans these
units the same way the History Sanitize Command does.

## Where

`roundfix baseline update` (Pending History plan and apply): `internal/cli/baseline_update_history.go`; detector `spec.HistoryCitations` in `internal/spec/history_entries.go`.

## Expected

- The `baseline update` plan reports, in text and JSON, every citation that the
  sanitization breaks, and the Plan Digest covers them.
- A citation from an active Spec, an open Backlog Entry or Finding, or
  `docs/references/` leads to a promotion offer, or at least to a warning that
  names the file and line before approval.
- The adopter release notes say how to find links an earlier update already
  broke.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-09-o-baseline-update-sanitiza-o-historico-sem-relatar-os-links-que-quebra.md`
  (vortex, 2026-10-09).
