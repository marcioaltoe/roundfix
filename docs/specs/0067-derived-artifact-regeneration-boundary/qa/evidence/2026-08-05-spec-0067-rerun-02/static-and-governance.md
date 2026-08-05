# Static and governance evidence

Build: `c6cf8033b7fe75d6f97404735e798cd42427ba89`.

## Project Constraints and tooling authority

- The PRD and TechSpec each account for identifier strategy, authentication
  and HTTP, active ADRs, and tooling authority with operative `docs/agents/`
  sources.
- Authorization commit
  `2e560cea708006024286881c5948702e1e4599c2` is an ancestor of tooling Task 03
  commit `c7ad3f62`; `git merge-base --is-ancestor` exited 0.
- Task 03's changed paths from
  `git diff-tree --no-commit-id --name-only -r c7ad3f62` are exactly
  `Makefile` and its assigned `task_03.md`.
- Corrective Task 06 commit `c6cf8033` is a separate descendant of Task 03;
  `git merge-base --is-ancestor c7ad3f62 c6cf8033` exited 0. Its changed paths
  are the authorized `Makefile`, assigned `task_06.md`, the ownership test,
  and three ownership records named by that corrective Task. No other tooling
  configuration changed, no prerequisite repair was folded into Task 03, and
  the consequent repair lands after the change that caused it.
- The authorization record names Spec 0067 and bounds its tooling mutation to
  `Makefile` for the regeneration step list and derived path scan.

## Repository Verification

Command, unpiped from the Run Worktree root:

```text
rtk make verify
exit 0
Go test: 3349 passed in 26 packages
isolated corpus budget: 1 passed
Skill tests: 4 passed
Roundfix skill check passed
go build -buildvcs=false: passed
Spec 0067 consistency check: No findings
```

The two Spec 0067 consistency skips concern a missing Vocabulary Contract and
`references/_index.md`; neither is required by this Spec's declared scope.

## Derived-content range audit

`rtk git -c core.fsmonitor=false diff --exit-code a188c987^ c6cf8033 --
internal/baseline/assets internal/baseline/testdata
':(exclude)**/_ownership.yml' ':(exclude)**/*_ownership.yml'` exited 0 with no
diff. Across the implementation range the only changes below the derived roots
are ten newly added ownership metadata files; no digest value or derived
artifact content moved.

