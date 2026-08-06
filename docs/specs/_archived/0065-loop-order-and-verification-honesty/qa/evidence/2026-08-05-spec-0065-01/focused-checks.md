# Focused assembled checks

Build: `d603031e808e3c64a539c4875f00d62ed778f522`.

- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/spec
  ./internal/speccheck -count=1 -run
  'TaskDocumentDeclarations|ContradictoryRequirements|UndeclaredRehearsal|Replay0060Task03RefusesContradictory'`
  exited 0. Both package selectors passed in 0.295 and 0.584 seconds. The
  selected tables cover same-subject contradiction refusal, undecidable-pair
  silence, missing and incomplete rehearsal refusal, and complete
  `Case`/`Observation` acceptance.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1
  -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$'` exited 0
  in 2.764 seconds. Active and archived corpus characterization remains at its
  recorded finding counts with no active errors.
- `rtk env GOCACHE=<worktree>/.gocache go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` exited 0 in 0.941
  seconds. The unpiped full gate independently ran the same dedicated selector
  and passed.

Together with the public replay and false-positive canaries, these checks
cover R08, R09, and R11 without relying only on Task Results.
