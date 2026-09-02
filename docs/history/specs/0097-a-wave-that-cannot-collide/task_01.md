---
status: completed
type: backend
---

# Task: The collision rule over a Task Graph

Two Tasks a graph declares independent that edit the same file die at
integration, after every byte of Agent work is done. The rule that would have
caught it does not exist: concurrency is configured as a number and verified as
nothing.

## Work

- One function over a Task Graph returns every pair the graph permits in the
  same wave that is known to touch one path, naming both Tasks, each shared
  path, and how the path was learned.
- Learn a Task's paths from three sources, none authoritative alone: the
  repository paths its Verification commands name, its declared `## Context`
  entries, and the paths a prior Run's settlement commit for that Task changed
  when such a Run exists. The measured collision declared no Context at all, so
  a rule reading only that source would report nothing.
- A path counts only when it resolves to a file in the repository. A package
  selector such as `./internal/cli`, a flag, or a test name is not a path, and
  the rule never infers what a Task means to edit.
- Tasks the graph already orders never collide. Follow `needs` transitively, not
  just directly: a chain of three serializes its ends as surely as an edge does.
- Read the graph and the repository only. No command runs, no Agent Session
  opens, and nothing is written.
- Cover the measured shape first: two Tasks with no declared Context whose
  Verifications name one file. Then a transitive `needs` chain, a package
  selector that must not become a path, and a path learned from each of the
  three sources.

## References

- `_prd.md` → Goal 1, User Story 1, Core Feature 1
- `_techspec.md` → Build Order 1; Interfaces: `Collisions`, `TaskTouchSet`,
  `TouchSource`
- ADR-0093 bounds the rule to what artifacts say rather than to inference

## Verification
- `grep -q "func Collisions" internal/spec/collision.go || exit 1; grep -q "TouchFromVerification" internal/spec/collision.go || exit 1; grep -q "TestCollisionsFindsTheMeasuredShape" internal/spec/collision_test.go || exit 1; go test -count=1 ./internal/spec`

## Result

Implementation:

- `Collisions` returns every unordered independent Task pair with its shared
  repository files and touch sources.
- Touch sets combine literal file paths from Verification commands, declared
  Context entries, and the newest reachable settlement commit for each Task.
  Git objects, including packed and delta-compressed objects, are read directly
  from the repository without starting a process.
- The dependency closure follows `needs` transitively in both directions before
  comparing a pair. Candidate paths must resolve to regular files inside the
  repository; directories, package selectors, flags, test names, expansions,
  and paths outside the repository are ignored.

Acceptance evidence:

- Same-Wave pair and source reporting: `TestCollisionsFindsTheMeasuredShape`
  exercises two Tasks with no Context whose Verification commands name
  `internal/speccheck/mechanical.go` and observes one collision attributed to
  `TouchFromVerification`; `TestCollisionsReturnsEveryPairAndSharedPath`
  observes all three pairs among three independent Tasks and both shared paths
  on the pair that names both.
- Three-source union: the measured-shape test covers Verification,
  `TestCollisionsLearnsPathFromDeclaredContext` covers Context, and
  `TestCollisionsLearnsPathFromPackedPriorRunSettlementCommitsWithoutCommands`
  covers settlement commits.
- File-only filtering:
  `TestCollisionsRejectsPackageSelectorsFlagsAndTestNames` exercises
  `./internal/cli`, `-run`, and `TestCommand` and observes no collision.
- Transitive ordering: `TestCollisionsExcludesTransitivelyOrderedTasks` uses
  `task_01 <- task_02 <- task_03` with one common file and observes no pair,
  including the chain's ends.
- Read-only execution: the packed-prior-Run test clears `PATH` before calling
  `Collisions`; the call still reads the settlement commits and returns the
  expected collision.
- Invalid-input behavior: `TestCollisionsRequiresGraphAndRepositoryRoot`
  observes errors for a missing graph and a missing repository root.

Focused checks:

- Pre-change signal: `rtk rg -n "func Collisions|TestCollisionsFindsTheMeasuredShape" internal/spec`
  exited 1 with no matches.
- `rtk go test -run '^TestCollisions' ./internal/spec` initially exposed packed
  symbolic-reference handling, then passed all 9 tests after the production
  reader fix.
- `rtk go vet ./internal/spec` exited 0 with no diagnostics.

Daemon Verification was not run in this Agent turn; the Daemon owns the
declared command and terminal Task settlement.
