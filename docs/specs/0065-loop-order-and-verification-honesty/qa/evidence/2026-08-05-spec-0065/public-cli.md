# Public CLI journeys

Built entry point:
`bin/roundfix spec check [<slug>] [--format text|json]`.

## Current assembled repository

- `rtk bin/roundfix spec check 0065-loop-order-and-verification-honesty` —
  exit 0; `No findings.`
- The same command with `--format json` — exit 0; schema
  `roundfix-speccheck/v1`, an empty `findings` array, and only the two declared
  skipped checks for the absent Vocabulary Contract and reference index.
- A fresh text process produced the same no-findings result.

## Spec 0060 replay

Scratch Git workspace:
`/private/tmp/roundfix-qa0065.wIsIi5`, branch `ma/qa-0065-replay`. The workspace
copies the shipped `replay-0060-task-03` fixture and the checker fixture's
agent, ADR, and authorization sources. No repository product file was edited.

The public JSON command exited 1 and emitted exactly these errors:

- `SC-VERIFY-WORK-INDEPENDENT` at `task_03.md:44`;
- `SC-REQUIREMENT-CONTRADICTORY` at `task_03.md:15` and `:13`, naming subject
  `commit`;
- `SC-REHEARSAL-UNDECLARED` at `task_03.md:9`.

This is the motivating Task observed through the built operator interface, not
only through a unit test.

## False-positive canaries and replay

The scratch Task was then changed only for the canaries:

- the contradictory phrase was made satisfiable;
- a complete `## Rehearsal Cases` entry was declared;
- Task Verification contained `make verify`, the clean-tree read, and a
  Task-specific `test -f .../qa/rehearsal-evidence.md` assertion.

Fresh JSON and text processes both exited 0 with no findings. Removing the
repository-wide gate and leaving only the Task-specific effect assertion also
exited 0 with no findings. This independently confirms both required
false-positive boundaries and process replay.

The initial scratch-copy command used relative sources from the scratch
directory and failed before copying anything. The corrected command used the
explicit workspace source paths; this setup mistake did not affect any CLI
verdict.
