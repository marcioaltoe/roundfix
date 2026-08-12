# R02 — carry-forward and defeat probes

The real first-Run/corrective-re-Run comparison is environment-blocked by the
same no-commit boundary as R01. Equivalent fresh evidence passed:

- `rtk ... go test ./internal/speccheck ./internal/daemon -run '^(TestCarriable|TestMechanicalStageCarriable|TestPriorChangedFilesUseCurrentWorktreeHeadAndIgnoreSiblingBranch|TestBuildQAPromptCarriesTheSpecContextBundle|TestBuildQAPromptCarriesThePreviousReportIdentity)$' -count=1 -v`
  passed 18 cases across two packages.
- The refusal cases cover prior non-pass, no inputs, moved input, non-ancestor,
  changed bytes, missing/malformed snapshots, changed glob membership, and a
  mixed `repository_path` plus `elapsed_time` input.
- The happy paths retain the original establishing report and head across a
  later carry. Prompt cases carry changed paths and prior-report identity.

No production re-gate timing was observed, so the material cost reduction is
not treated as passed.
