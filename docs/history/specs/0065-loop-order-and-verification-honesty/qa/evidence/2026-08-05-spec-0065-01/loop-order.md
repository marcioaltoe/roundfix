# Loop-order and evidence-count retest

Build: `d603031e808e3c64a539c4875f00d62ed778f522`.

The settled action order remains identical in the repository guide, shipped
formatter clause, and Baseline module:

`implement the graph including its authored gate, archive, open the Pull
Request, watch until Clean, and merge`.

Fresh
`rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1 -run '^TestCheckLoopOrder'`
exited 0 in 0.377 seconds. It covers current agreement plus an independent
seeded divergence in each of the three rule-owned sources.

## Archived source and corrected carriers

The archived Spec 0078 QA Report independently records:

- frontmatter `rows_blocked_environment: 11`;
- R06–R13 and R18, nine rows, blocked on no open Pull Request;
- Coverage of 7 pass and 11 environment-blocked rows out of 18.

The guide, both canonical Skills and mirrors, shipped formatter clause,
Baseline module, and TechSpec now state `eleven of eighteen` and attribute
nine of those eleven to the absent Pull Request. The bounded stale-count
search used by Task 07 exits 1 with no matches.

## F-002 repeated — active Task evidence still states the stale count

Impact: Trust-Damage.

Actor and step: a maintainer reads task_01's `## Result`, which is named
implementation evidence for the loop-order rationale.

Expected: the active Task Result matches the archived report and corrective
Task 07's requirement that no carrier state `nine of eighteen`.

Actual: `task_01.md:98` still says Spec 0078's gate proved the path with
`nine of eighteen rows blocked on no open Pull Request`. A repository-wide
search excluding the historical first QA report finds that active claim. The
Task 07 bounded search omitted every Task file, so its stated no-match evidence
could not detect this carrier.

The action order and ADR-0080/ADR-0091 conclusion remain correct. The failure
is the quantitative evidence carried by the active Task Result, the same
symptom prior F-002 identified.

Affected row: R10.
