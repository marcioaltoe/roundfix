# Focused assembled checks

Build commit: `9252430f9e6c63332775a90ee9dcb08f7bbccef7`.

- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/spec
  ./internal/speccheck -count=1 -run
  'TaskDocumentDeclarations|ContradictoryRequirements|UndeclaredRehearsal|Replay0060Task03RefusesContradictory'`
  exited 0. The `internal/spec` and `internal/speccheck` selectors passed in
  0.755 and 0.500 seconds. The selected tables cover same-subject
  contradiction refusal, undecidable-pair silence, missing and incomplete
  rehearsal refusal, and complete `Case`/`Observation` acceptance.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1
  -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$'` exited 0 in
  1.704 seconds. Active and archived corpus characterization remains at its
  recorded finding counts with no active errors.
- `rtk env GOCACHE=<worktree>/.gocache go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` exited 0 in 1.231
  seconds. The unpiped full gate independently ran the same dedicated selector
  and passed.

Together with the public replay and false-positive canaries, these checks make
R08, R09, and R11 pass without relying only on Task Results.
