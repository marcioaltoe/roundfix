# Static gate evidence

Build: `c6a5ea156c70d33d59e9e5e1d98fcc2fe24cb318`.

- `rtk roundfix spec check 0130-documentation-cleanup-compatibility --strict`
  exited 0 with `No findings. Authored Verification commands were not executed.`
  It reported the missing `docs/backlog`, Vocabulary Contract and references index
  inputs as skipped detectors; QA retained the applicable manual checks.
- `rtk make verify` exited 2 at `fmt-check`. It named
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` as needing formatting.
- `rtk gofmt -d internal/cli/baseline_skills_restore_test.go
  internal/cli/baseline_assets_sync_test.go` exited 1 and showed indentation-only
  differences in the two files.
- `rtk git diff --name-status 54614682..c6a5ea15 --
  internal/cli/baseline_skills_restore_test.go
  internal/cli/baseline_assets_sync_test.go` returned no paths. The formatter
  failure therefore predates the authorized repair and is outside its scope.
- `rtk make verify-docs` exited 0. It built the CLI, passed `docscontract`, passed
  the tagged Baseline and Skill regeneration contracts, passed the corpus budget,
  and reported no findings for active Spec consistency checks.

