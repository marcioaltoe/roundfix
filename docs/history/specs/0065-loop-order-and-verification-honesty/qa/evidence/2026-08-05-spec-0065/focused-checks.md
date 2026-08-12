# Focused assembled checks

All passing reruns used the repository-local Go cache already selected by the
Makefile.

- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/spec
  ./internal/speccheck -count=1 -run
  'TaskDocumentDeclarations|ContradictoryRequirements|UndeclaredRehearsal|Replay0060Task03RefusesContradictory'`
  — exit 0; both packages passed. This covers same-subject refusal,
  undecidable-pair silence, missing/incomplete rehearsal refusal, and complete
  declaration acceptance.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1
  -run '^TestCheckLoopOrder'` — exit 0; all current-agreement and per-source
  seeded-divergence cases passed.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -count=1
  -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$'` — exit 0.
- `rtk env GOCACHE=<worktree>/.gocache go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` — exit 0 in 0.942s.

The first focused declaration command used the sandbox-external default Go
cache and exited 1 before tests ran because access under
`/Users/marcio/Library/Caches/go-build` was denied. This is an environment-only
attempt. The documented repository-local cache rerun passed; the full
repository gate independently ran all packages and passed.
