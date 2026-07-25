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

## Current-build rerun after F-06 repair

Revision:
`5a97f2d8b85630dc413f0fed738cd32b24d5a228`.

The QA binary was built with that exact identity and had SHA-256
`813fe65580f4d1638b727d59ad4ed34248236601d53648e954701674596b3bae`.
The disposable run root was
`/private/tmp/roundfix-qa-0047.ECKGEG`. Clean local clones were made from:

- Fluxus `1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`
- Oraculum `ad74f46197500de63dc0d9ff0d3e09f61a6a43ce`

The original source checkouts were clean before and after the run.

### Mandatory static gate and repaired race

```text
rtk make verify
Go test: 2201 passed in 22 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed
Roundfix build passed
exit 0

rtk go test -count=100 ./internal/baseline \
  -run 'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'
exit 0
```

The exact F-06 regression passed all 100 consecutive repetitions without an
environment override.

### LIVE-01 — Fluxus greenfield

The public automation Plan exited `0`:

```text
roundfix baseline plan \
  --repo <fluxus-greenfield> \
  --profile standard-typescript-monorepo \
  --decision-file fluxus-greenfield-decisions.json \
  --format json

Plan Digest:
sha256:c52337dce84d5153752f70bebb671f1907bb1d0b8f12cdd29a8709c73952d121
File changes: 13
Managed entries: 26
Upgrade Retention Contract entries: 0
Warnings: 13
```

Applying with an all-zero digest exited `3`, named the approved digest, and
left the disposable repository unchanged. Applying the exact digest exited
`0`, reported `state: verified`, and verified 15 postimages. `rtk bun run fmt`
and the disposable Fluxus `rtk make verify` both passed.

Fresh generated-state checks confirmed the ordered Instruction Hierarchy, the
strengthen-only rule, no generic or repository-specific carrier, the complete
ADR lifecycle overlay, and the complete Findings template.

The required first fresh Plan then exited `2`:

```text
roundfix: baseline plan failed: validate assembled Baseline Plan:
root carrier source "AGENTS.md" has no immutable backup
```

Focused artifacts:
`rerun/fluxus-greenfield-plan.json`,
`rerun/fluxus-greenfield-wrong-digest.json`,
`rerun/fluxus-greenfield-apply.json`, and
`rerun/fluxus-greenfield-replan-stderr.txt`.

### LIVE-02 — Fluxus update

The public interactive workflow was driven through a real PTY:

1. selected Preservation and the existing
   `standard-typescript-monorepo` Profile;
2. retained every current repository decision;
3. reached ready Profile alignment;
4. reviewed 19 exact `AGENTS.md` source segments, each proposed as
   non-governed with an individual reasoned rejection;
5. accepted the complete proposal;
6. reviewed 13 file changes, 26 managed entries, and all 19 Upgrade Retention
   Contract entries;
7. confirmed Plan Digest
   `sha256:434be69ff208a5c34da005dd4d18f9094f436baeacba4c8c5847162563eb5e82`.

Apply exited `0`, reported `Baseline apply: verified`, and verified 15
postimages. `rtk bun run fmt` and the disposable Fluxus `rtk make verify`
passed.

On a fresh public interactive run, the maintainer again retained the current
Profile and decisions and accepted the complete 19-entry proposal. Before a
fresh Change Plan appeared, the command exited `1`:

```text
roundfix: baseline failed: build human Baseline Plan:
validate assembled Baseline Plan:
root carrier source "AGENTS.md" has no immutable backup
```

### LIVE-03 — Oraculum divergence and adaptation

The built-in `standard-typescript-monorepo` automation path exited `3` and
named ten profile-specific blockers individually. The public interactive flow
then:

1. selected Greenfield and the existing built-in Profile;
2. reached `action_required` alignment;
3. chose repository-owned Profile adaptation;
4. reviewed and accepted removal of `frontend`, `autonomous-work`, and the ten
   profile-specific capabilities;
5. supplied `oraculum-backend`;
6. re-audited the draft to ready;
7. reviewed the Profile postimage, 14 file changes, 23 managed entries, and
   the empty retention ledger;
8. confirmed Plan Digest
   `sha256:c5e892ce37890bdd76a4e46376a65643d6b9760f48d92df39b3efaebef710911`.

Apply exited `0`, verified 14 postimages, and wrote the Profile only through
the approved Plan. Public `baseline profile validate` reported the resulting
repository-owned Profile `valid`; the disposable Oraculum `rtk make verify`
passed. Supplying `--profile` with `--profile-file` exited `2` without
mutation.

A fresh automation Plan with the applied Profile and its selected decisions
then exited `2` with the same missing immutable backup error as both Fluxus
journeys.

Focused artifacts:
`rerun/oraculum-built-in.json`,
`rerun/oraculum-backend-profile.json`,
`rerun/oraculum-mutual-exclusion.json`, and
`rerun/oraculum-replan-stderr.txt`.

## Current-build conclusion

F-06 is fixed and the mandatory gate is green. Initial public adoption,
Readoption, and Profile-adaptation planning and apply boundaries now work, but
all three live journeys fail the required persistence/empty-reapply step on
the same assembled-Plan validation error.
