# Feature contract evidence

All commands ran from the audited Run Worktree at `c6a5ea15`.

- `rtk proxy go test -count=1 -v ./internal/suiteguardcontract -run
  '^TestCleanupRegenerationDiscovery$'` exited 0. The approved Spec grant,
  compatible legacy declaration, symlink refusal, missing optional roots and
  invalid-root error cases all passed.
- The ownership group first encountered a nested Go-cache sandbox denial:
  `open /Users/marcio/Library/Caches/go-build/...: operation not permitted`.
  The unchanged focused reproduction passed with cache access. The complete
  command `rtk proxy go test -count=1 -v -tags repocontract
  ./internal/baseline -run
  '^(TestCleanupRegenerationOwnershipParity|TestOutputsForCommand|TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries)$'`
  then exited 0. It exercised reader parity, exact ownership slices, invalid and
  conflicting declarations, the dedicated command, wrong commands, unchanged
  output, sanctioned rewriting, frozen refusal and idempotence.
- `rtk proxy go test -count=1 -v -tags repocontract ./internal/speccheck
  -run '^(TestCleanupHistoricalGrantEvidence|TestAuditJudgesTheGrant|TestEveryBoundedPathIsGoverned)$'`
  exited 0. The two unavailable historical Task commits were explicit skipped
  subtests; the real historical corpus, unmatched-path refusal, governed-path
  refusal, hand-edited-derived refusal, wrong-Spec refusal and folded
  authorization refusal all executed.
