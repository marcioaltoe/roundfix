---
type: fix
status: open
created: 2026-10-06
spec: null
---

# An unformatted QA report breaks the next Run's precondition

## Problem

In Pantheon, `make verify` runs `bun x oxfmt@latest --check .`, and it is the
QA gate's precondition. With Roundfix 0.46.0, on Spec 0047, on 2026-10-06:

1. **First QA Run.** `run_20261006T132943Z_87bfc40ea92ff2ff` failed and wrote
   `qa/qa-report-2026-10-06.md` without formatting it.
2. **Next Run.** `run_20261006T134004Z_f96397038ab85269` started from the
   branch head, which did not have that report. The Run carried the
   invalidated report into its worktree anyway. Its `make verify`
   precondition exited 2 on that file's formatting, and the Run wrote a
   `fail` report with `rows_blocked_precondition: 1` without running QA
   (`Tokens: no prompts recorded`).
3. **Workaround.** The operator brought the previous Run's `qa/` onto the
   branch, formatted it and committed it. The third Run closed with a
   qualifying `partial`.

Every passing QA in Pantheon also needed a separate formatting commit (Specs
0042, 0046, 0047).

A second report from the same entry: `roundfix reopen` does not reopen a
completed QA Task when a new corrective Task, already `completed`, becomes one
of its dependencies. The operator had to set `status: pending` by hand
(Spec 0046).

## Expected

1. **Formatting.** A Run formats the QA report and its evidence with the
   repository's own formatter before committing them, or the precondition
   ignores the files the Run itself carries. Decide which, and how the
   formatter is discovered (Baseline stack setup, `make fmt`, a configured
   command).
2. **Reopen.** When the Task Graph gains a dependency of a completed QA Task,
   `roundfix reopen` (or the settlement of the new Task) returns the QA Task to
   `pending`.

## Sources

Secondbrain inbox, triaged 2026-10-06:
`inbox/roundfix/_triaged/2026-10-06-relatorio-de-qa-sem-formatacao-derruba-a-proxima-run.md`
(Pantheon).
