# Public CLI journeys

Built entry point: `bin/roundfix spec check [<slug>] [--format text|json]`.

Build: `d603031e808e3c64a539c4875f00d62ed778f522`; `make verify` rebuilt the
binary before these journeys.

## Current assembled repository

- `rtk bin/roundfix spec check 0065-loop-order-and-verification-honesty`
  exited 0 and reported `No findings.`
- A fresh JSON process exited 0 with schema `roundfix-speccheck/v1`, an empty
  `findings` array, and only the two declared skipped checks for the absent
  Vocabulary Contract and references index.
- The authored Task Graph remains runnable before a Pull Request exists:
  `qa: task_06` is terminal, every non-QA dependency is completed, and the QA
  Task remains Daemon-owned `pending` during this gate.

## Spec 0060 replay

Scratch Git repository: `/private/tmp/roundfix-qa0065-01.rUiuCU`, branch
`ma/qa-0065-replay`. It copies the shipped `internal/speccheck/testdata/repo`
fixture and does not edit the product worktree.

The public JSON command exited 1 and emitted exactly:

- `SC-VERIFY-WORK-INDEPENDENT` at `task_03.md:44`;
- `SC-REQUIREMENT-CONTRADICTORY` at `task_03.md:15` and `:13`, naming subject
  `commit`;
- `SC-REHEARSAL-UNDECLARED` at `task_03.md:9`.

A fresh text process independently exited 1 with the same three codes,
locations, and fixes. This observes the motivating Task through the built
operator interface rather than crediting its unit test alone.

## False-positive canaries

The scratch Task was made satisfiable, given a complete
`## Rehearsal Cases` declaration, and assigned `make verify`, a clean-tree
read, and a Task-specific file assertion. Fresh JSON and text public processes
both exited 0 with no findings. This proves a repository-wide gate remains
accepted when a command can distinguish Task work from no work.

Removing the repository-wide gate and clean-tree read, leaving only the
Task-specific effect assertion, also exited 0 in a fresh JSON process. This
independently covers the effect-only acceptance boundary.

The scratch fixture intentionally lacks the three repository loop-order
carriers, which the CLI lists as skipped there; current-carrier behavior is
covered against the assembled repository by R04 and R10.
