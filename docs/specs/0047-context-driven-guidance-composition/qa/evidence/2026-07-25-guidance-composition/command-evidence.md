# Command evidence — 2026-07-25 guidance composition

## Build identity and starting tree

```text
rtk git -c core.fsmonitor=false rev-parse HEAD
eaf9543938da78532b40310555328a46ef2bf918

rtk git -c core.fsmonitor=false status --short
?? docs/specs/0047-context-driven-guidance-composition/qa/qa-report-2026-07-25.md
```

The worktree had no pre-existing changes. The only reported path was the QA
report created before execution, as required by the resumable gate.

## Mandatory static gate

```text
rtk make verify
make: *** [test] Error 1
rtk go test ./...
Go test: 2199 passed, 2 failed, 1 skipped in 22 packages
baseline (317 passed, 2 failed)
  [FAIL] TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout
     testing.go:1464: TempDir RemoveAll cleanup: unlinkat .../.git: directory not empty
  [FAIL] TestAssetsSyncProvenanceAndPreMutationRefusals
```

The gate stopped in its test phase, before the full verification pipeline
could complete.

## Root-cause isolation

One isolated run passed, proving the failure was timing-dependent:

```text
rtk proxy go test -count=1 -v ./internal/baseline \
  -run 'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'
PASS
```

Fifty unmodified repetitions reproduced the cleanup race twice:

```text
rtk proxy go test -count=50 ./internal/baseline \
  -run 'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'
TempDir RemoveAll cleanup: unlinkat .../002/.git: directory not empty
TempDir RemoveAll cleanup: unlinkat .../002/.git/objects/info: directory not empty
FAIL
```

Both failed temporary repositories retained a file created after cleanup began:

```text
.../002/.git/objects/info/packs
```

The fixture already disables `core.fsmonitor`, but its Git invocations do not
disable automatic maintenance. The host has `core.fsmonitor=true`, and Git
background work can outlive the command that created the temporary
repository. A diagnostic-only single-variable run disabled automatic
maintenance and passed all 100 repetitions:

```text
rtk proxy env \
  GIT_CONFIG_COUNT=1 \
  GIT_CONFIG_KEY_0=maintenance.auto \
  GIT_CONFIG_VALUE_0=false \
  go test -count=100 ./internal/baseline \
  -run 'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'
ok roundfix/internal/baseline
```

This override was used only to identify the race. It was not used to replace or
reclassify the failed mandatory gate.

## Live-flow decision

The `qa-gate` static-gate contract requires flow QA to stop on a code-caused
failure. No Fluxus or Oraculum source checkout was copied or mutated, and no
live Baseline plan/apply journey was started.
