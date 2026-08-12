# Pull Request environment and equivalent evidence

The QA prompt states that no Pull Request is open and Pull Request journeys
are environment-blocked. The Run Worktree branch is per-Run, is never pushed,
and has no Pull Request of its own. R16 is therefore
`blocked (environment: no open Pull Request)` without attempting a lookup from
the Run Worktree branch.

Unblock the hosted journey by opening a Pull Request from the Spec target
branch `ma/0065-loop-order-and-verification-honesty`, then rerun the journey
read-only against that Pull Request.

Equivalent observed and supervised evidence exists in
`docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`
and its linked `equivalent-pr-evidence.md`. That passing report records
`rows_blocked_environment: 11`; nine of those eleven rows (R06–R13 and R18)
were blocked on no open Pull Request and link public CLI, request-boundary,
command-runner, persistence/replay, and Run Event evidence. Its Coverage
independently records 7 pass and 11 environment-blocked rows out of 18.

ADR-0080 is accepted and unchanged across this Spec's commit range. It permits
an environment-blocked row to avoid capping the verdict only when equivalent
observed or supervised evidence exists; R16 uses that exact path.
