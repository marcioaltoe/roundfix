# Loop-order and evidence-count retest

Build commit: `9252430f9e6c63332775a90ee9dcb08f7bbccef7`.

The settled action order remains identical in the repository guide, shipped
formatter clause, and Baseline module:

`implement the graph including its authored gate, archive, open the Pull
Request, watch until Clean, and merge`.

The first direct focused selector inherited a sandbox-external user Go cache
and failed before compilation with `operation not permitted`. A single rerun
using the repository-local `.gocache` documented by the project exited 0:

`rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1 -run '^TestCheckLoopOrder'`

The selector covers current agreement plus an independent seeded divergence in
each of the shipped clause, repository guide, and Baseline module sources. The
full unpiped gate independently passed all assembled tests and the public
repository Spec check.

## Archived source and active carriers

The archived Spec 0078 QA Report independently records:

- frontmatter `rows_blocked_environment: 11`;
- R06–R13 and R18, nine rows, blocked on no open Pull Request;
- Coverage of 7 pass and 11 environment-blocked rows out of 18.

A fresh bounded `rg` sweep shows the guide, both canonical Skills and mirrors,
shipped formatter clause, Baseline module, TechSpec, and active task_01 Result
state `eleven of eighteen`; each attributes nine of those eleven to the absent
Pull Request. Occurrences of `nine of eighteen` in task_01 and task_07 quote
the stale wording while documenting its correction or express the bounded
absence assertion; they do not assert it as the current measurement.

The prior F-002 is closed. Result: R10 passes.
