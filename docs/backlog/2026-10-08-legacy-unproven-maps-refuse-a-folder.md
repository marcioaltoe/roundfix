---
type: fix
status: open
created: 2026-10-08
spec: null
---

# A legacy `unproven` list of maps refuses its folder

## Problem

In Fluxus, two Legacy Archive Folders have `_prd.md` front matter whose
`unproven` field is a list of maps, the format older Roundfix versions wrote,
instead of a list of strings:

- 0032 has items with `row`, `goal`, `claim` and `satisfied-by`.
- 0039 has items with `row`, `claim` and `reason`.

Roundfix 0.55 refused the whole plan. Since 0.57, the plan works and lists both
folders as Refused Units, measured with 0.59 on 2026-10-08: 57 units, 23.9 MB,
2 refused. But neither folder can convert, and an archived Spec cannot be
edited in place. The refusal reason also prints the YAML error over several
lines (`yaml: unmarshal errors:` followed by the detail), so the plan line is
broken.

## Direction

- Under Lenient Legacy Reading, accept the legacy map form of `unproven`. Turn
  each map into one text line in the Archive Record, keeping its keys in a
  stable order, for example `row 03: <claim> (reason: ...)`. Active Specs keep
  today's string-only check.
- Print each Refused Unit's reason on one line.

## Sources

Secondbrain inbox, triaged 2026-10-08:
`inbox/roundfix/_triaged/2026-10-07-history-sanitize-recusa-unproven-em-formato-antigo.md`
(Fluxus).
