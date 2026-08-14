# Focused boundary checks

Build: `9d7f834e`

- `rtk make verify` exited 0 on the unmodified build: 3,347 Go tests passed,
  the isolated corpus budget test passed, four Skill tests passed, Repository
  Skill Set check passed, the binary built, and `bin/roundfix spec check`
  reported no findings for Spec 0067.
- `rtk go test ./internal/baseline -count=1 -run
  '^(TestDerivedOwnership.*|TestDeclaredStepRegenerationAndFrozenBoundaries)$'
  -v` exited 0 with 24 tests passing.
- In a fresh scratch clone, a stale plan-characterization golden made the
  strict test exit 1 and print the exact record path
  `testdata/plan-characterization/_ownership.yml` and its declared update
  command. The printed command exited 0, and the fresh strict test then exited
  0.
- In the same scratch clone, a stale parity fixture made
  `TestBaselineCompatibilityCorpus` exit 1. Its diagnostic said nothing
  regenerates the artifact and printed the parity ownership record path plus
  the full 2026-07-30 tried-and-reverted reason. Running the sanctioned command
  left that directly perturbed fixture unchanged and reported `changed: false`.
- `rtk go test ./internal/baseline -count=1 -run
  '^TestDerivedOwnershipRemediationDiagnostics$' -v` exited 0 with four tests
  passing, including the exact sanctioned remediation text.
- `rtk git diff --exit-code a188c987..HEAD -- internal/baseline/assets
  internal/baseline/testdata ':(exclude)**/_ownership.yml'
  ':(exclude)**/*_ownership.yml'` exited 0. The implementation changed no
  pre-existing artifact or digest under the derived roots.
- `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  c7ad3f62` listed only `Makefile` and `task_03.md`. Authorization commit
  `2e560cea` is an ancestor and records Spec 0067 with bounded path `Makefile`.

The catalog-diagnostic ownership record says the artifact is not in
`BASELINE_DIGEST_STEPS`, while `Makefile:110` places its exact test and update
flag in that list. The parity ownership record says its corpus is regenerated
by nothing, while the owned-Skill journey above caused the sanctioned command
to rewrite `fixtures/asset-sync.json` and `manifest.json` in that corpus.
