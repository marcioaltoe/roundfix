---
type: fix
status: open
created: 2026-10-06
spec: null
---

# Settlement refuses a QA partial that archive and `qa-report accept` accept

## Problem

The `archive` reference says the Pull Request row "never decides a qualifying
partial" and that "Settlement and archive both apply this same policy".
Settlement applies a different rule:

- Fluxus, Spec 0100 (0.43/0.44): a `partial` with `rows_blocked_declared: 2`
  (both in `## Unreachable Acceptance`), `rows_blocked_environment: 1` (only
  the Pull Request row) and `rows_blocked_finding: 0`. `roundfix qa-report
  accept` exits 0, but the Daemon settles the QA Task `failed`, and `roundfix
  settle` says `rows_blocked_environment is 1; expected 0`.
- Oraculum, Spec 0069 (0.44, Run `run_20261006T030644Z_09840ba7986a5634`): a
  `partial` whose only blocked row is the Pull Request row, with no
  Unreachable Acceptance. The Daemon settles `failed`, the queue parks
  `qa-environment-partial`, and `roundfix settle` says `newest QA Report
  verdict is "partial"; expected "pass"`.

This repository meets it on nearly every delivery: the operator archived
0222–0234 with `--qa-override` because the Pull Request row (and sandbox-only
outside-evidence rows) left the report `partial` (operator log, 2026-10-04 to
2026-10-06). Each case costs a retry and an override.

## Expected

Settlement, archive and `qa-report accept` apply one policy. A `partial`
whose only unmet rows are the Pull Request row and declared Unreachable
Acceptance rows qualifies everywhere (or the QA gate writes `pass` for a
report whose only blocked row is the Pull Request row). Decide also whether
outside-evidence rows a Run sandbox cannot reach (network-denied) count as
environment rows that need the override.

## Sources

Secondbrain inbox, triaged 2026-10-06:
`inbox/roundfix/_triaged/2026-10-05-settle-recusa-partial-que-o-qa-report-accept-aceita.md`
(Fluxus) and
`inbox/roundfix/_triaged/2026-10-06-partial-so-com-a-linha-do-pr-e-recusado-no-assentamento.md`
(Oraculum).
