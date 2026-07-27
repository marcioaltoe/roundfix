# Command evidence — build 1736dbc

All commands ran from the repository root unless a section names another
working directory.

## PC-01 — Project Constraints and protected tooling scope

Both `_prd.md` and `_techspec.md` account for:

- identifier strategy: not applicable, with
  `docs/agents/domain.md` as the operative source;
- authentication and HTTP: not applicable, with
  `docs/agents/cli.md` as the operative source;
- active ADR obligations: applicable, citing ADR-0049, ADR-0055, ADR-0066,
  ADR-0072, and ADR-0077 through `docs/agents/domain.md`;
- Tooling Authority: applicable, with express 2026-07-26 maintainer
  authorization for the exact nine protected paths and
  `docs/agents/agent-instructions.md` as the operative source.

All six Task files have `status: completed`.

`rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r c27c8a1defa10405e0f159d1520ec33f0b28a7bb`
exited 0:

```text
docs/specs/0051-doctor-readiness-contract-reconciliation/task_01.md
go.mod
go.sum
```

`rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r fbc2ae4d2eabf232b5aeaf072c423667bce4b7b2`
exited 0:

```text
.agents/skills/roundfix/SKILL.md
docs/specs/0051-doctor-readiness-contract-reconciliation/task_05.md
skills/roundfix/SKILL.md
```

`rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r 1736dbc77683ab3704a6169a067eba09eb05d923`
exited 0:

```text
docs/specs/0051-doctor-readiness-contract-reconciliation/task_06.md
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
```

Each protected Task commit contains only its assigned Task file and its exact
authorized paths. The current worktree delta contains only this QA report and
its evidence.

## SG-01 — Full repository Verification

`rtk make verify` exited 0 on build `1736dbc`:

```text
rtk go test ./...
Go test: 2420 passed in 23 packages
rtk go test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'
Go test: 4 passed in 1 packages
rtk go run -buildvcs=false ./cmd/roundfix skills check
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
rtk go build -buildvcs=false -ldflags "-X 'roundfix/internal/app.BuildCommit=1736dbc-dirty' -X 'roundfix/internal/app.BuildTime=2026-07-26 21:05:00 -0300'" -o bin/roundfix ./cmd/roundfix
```

The post-gate
`rtk git -c core.fsmonitor=false status --porcelain=v1` output contains only
the current QA report and its new evidence directory.

## J-03 — Outside-Git working-directory separation

The built `bin/roundfix doctor` ran twice from the empty non-Git directory
`/private/tmp/roundfix-qa-0051-outside-git`. Both executions exited 1 with the
same ordered report. The environment's existing Codex signature problem made
`profiles:` and `codex:` fail, while the Repository Skill Set line was exactly:

```text
skills: failed (Repository Skill Set readiness requires a Git repository; next: run roundfix doctor from a Git repository)
```

The scratch directory was empty before and after both runs. The repeated public
output was byte-stable. Independent confirmation
`rtk go test ./internal/cli -run 'TestRunDoctorMissingRepositoryRoot' -count=1`
exited 0 with one passing test; that test proves the profile seam receives the
process working directory and the Repository Skill Set checker is not called.

## J-01 — Cooperative cancellation

A disposable Git repository under
`/private/tmp/roundfix-qa-0051-cancel-git-1736dbc` contained the current
Repository Skill Set plus a 4 GiB external-skill file to keep the public
filesystem walk active. The built `roundfix doctor` ran in a PTY. Ctrl-C during
that walk terminated the command immediately with exit 1 and no readiness
report after the interrupt marker. The scratch repository was removed after
the probe.

Independent confirmations:

- `rtk go test ./skills -run 'Test(CheckRepository|SkillFolderHash)' -count=1`
  exited 0: 40 tests passed.
- `rtk go test -race ./skills ./internal/cli -run
  'Test(CheckRepository|SkillFolderHash|RunDoctor)' -count=1` exited 0: 58
  tests passed across two packages.

Together the public interrupt and named package evidence cover prompt
termination, pre-cancelled and during-read cancellation, error-chain
inspection, traversal stop, context propagation through Doctor, context-first
caller compilation, and the affected race surface.

## J-02 — Total external hash order

- Two fresh executions of
  `rtk go test ./internal/skillhash -run 'TestSum' -count=1` each exited 0
  with eight passing tests.
- `rtk go test ./internal/baseline -run 'TestSkillsRestore' -count=1` exited 0
  with 14 passing tests.
- `rtk go test -race ./internal/skillhash ./internal/baseline -run
  'Test(Sum|SkillsRestore)' -count=1` exited 0 with 22 passing tests across two
  packages.

The repeated fresh runs independently confirm the pinned tie-corpus digest.
The named group also covers American English primary collation, normalized raw
path tie-breaking, punctuation and Unicode compatibility, path/content
negative cases, caller-order preservation, the shared Baseline consumer, and
the affected race surface.

## J-04 — Fail-closed ownership remediation

A disposable Git repository contained the current Repository Skill Set plus
one extra file in the Roundfix-owned `roundfix` skill and one in the externally
managed `coding-guidelines` skill. Two executions of the built Doctor command
exited 1 with byte-stable ordered output. The exact skills line was:

```text
skills: failed (outdated: coding-guidelines, roundfix; next: roundfix skills install --target project && bunx skills experimental_install && bunx skills update -p -y)
```

Both public runs continued to the independent `codex:` check. A tar-stream
SHA-256 of `.agents` and `skills-lock.json` was
`a5f7207bd5acf7ff6156d2bb440a94191828cddb2a86e4e3de19545886efe7b8`
before and after both runs, proving diagnosis did not mutate the scratch
Repository Skill Set. The scratch repository was removed after capture.

Independent confirmations:

- `rtk go test ./internal/cli -run
  'TestRunDoctorRepositorySkillReadiness' -count=1` exited 0 with seven passing
  coordinator cases.
- `rtk go test ./skills -run
  'TestCheckRepositoryRejectsSymlinkedAuthoritiesBeforeReadingTargets'
  -count=1` exited 0 with four passing ownership/confinement cases.

Those cases cover external-only remediation for a symlinked lock, target
non-read, shared/owned/external classifications, conservative unclassified
errors, the exact mixed chain, and eager line order.

## J-05 — Tool-produced module metadata

Fresh evidence that completed:

- `rtk go mod edit -json` exited 0 and reports
  `golang.org/x/text v0.40.0` without `Indirect`.
- `rtk go mod verify` exited 0 with `all modules verified`.
- `rtk git -c core.fsmonitor=false diff --exit-code -- go.mod go.sum` exited 0,
  proving the current build has no module-file delta.
- `rtk git -c core.fsmonitor=false show
  c27c8a1defa10405e0f159d1520ec33f0b28a7bb -- go.mod` shows only
  `golang.org/x/text v0.40.0` moving from the indirect block to the direct
  block; the selected version is unchanged.

The daemon sandbox could not complete `rtk go mod tidy -diff` because it
denied reads under `/Users/marcio/Library/Caches/go-build`. After the Run
closed and its commits were integrated, the supervisor executed the same
literal command from the integrated branch in a full-access session. It exited
0 with no output. No alternate cache, environment override, dependency change,
or module-file write was used. This closes the environment-only block with the
exact required evidence.

## AC-01 — Repository Skill Set inspection

The fresh J-01 commands are the complete named AC-01 group: 40 focused
`CheckRepository`/`SkillFolderHash` tests and 58 affected race-tested cases
passed. J-04 adds four explicit symlink-authority and target-non-read cases.
The public interrupt supplies the actor-level cancellation evidence. Together
they cover all six Task 03 acceptance criteria.

## AC-02 — Doctor coordinator and public surface

The built Doctor command ran inside the real Git repository twice during this
gate. Both runs kept the ordered six-line interface. Repository Skill Set
readiness passed:

```text
skills: ok (39 required: 14 Roundfix-owned, 25 external)
```

The machine's pre-existing invalid/missing Codex code signature made
`profiles:` and `codex:` fail with the documented actions, and Doctor exited
1. This parity deviation does not prevent validating the Spec's independent,
eager Repository Skill Set contract: `skills:` still ran after `profiles:`
failed and before `codex:`.

The invalid-input probe
`rtk ./bin/roundfix doctor --unexpected` exited 2 with:

```text
roundfix: doctor failed: flag provided but not defined: -unexpected
Run 'roundfix doctor --help' for usage.
```

Focused and independent checks:

- `rtk go test ./internal/cli -run 'Test(RunDoctor|Doctor)' -count=1`:
  18 passed.
- `rtk go test -race ./internal/cli -run 'Test(RunDoctor|Doctor)' -count=1`:
  18 passed.
- `rtk grep -n 'func TestRunDoctor' internal/cli/doctor_test.go` found all nine
  Doctor behavior entry points.
- The same grep in `internal/cli/cli_test.go` exited 1 with no match.
- A fresh Git status contains only the QA report and this evidence directory.

The public command, invalid-input probe, focused suite, race suite, source
placement, and status confirmation cover all seven Task 04 acceptance
criteria.

## AC-03 — Synchronized Roundfix guidance

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited 0.
- `rtk make skills-sync-check` exited 0 with four passing tests.
- Exact-string greps found the canonical missing-Git result at line 52 and the
  mixed fail-closed chain at line 58 in both Skill copies.
- `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  fbc2ae4d2eabf232b5aeaf072c423667bce4b7b2` lists only the canonical Skill,
  shipped copy, and `task_05.md`.

The current bytes, repository synchronization guard, exact public wording, and
Daemon-owned Task commit prove all five Task 05 acceptance criteria.

## AC-04 — Reconciled derived snapshot

- `rtk go test ./skills -run
  'TestAuthorialSkillSync/typescript-bun.json' -count=1` exited 0 with two
  passing results.
- `rtk go test ./internal/baseline -run
  'Test(CatalogCompatibility|AssetsSyncCompatibilityMatchesMaintainedPythonContract|BaselineCompatibilityCorpus)'
  -count=1` exited 0 with three passing owners.
- The full `rtk make verify` result in SG-01 passed after the last derived
  artifact edit.
- `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  1736dbc77683ab3704a6169a067eba09eb05d923` lists only the five authorized
  derived artifacts and `task_06.md`.
- A current `git diff --exit-code` over those five artifacts exited 0.

The authorial, catalog, asset-sync, parity, full-gate, commit-scope, and
unchanged-current-state evidence proves all five Task 06 acceptance criteria
and Core Feature 7.

## DOC-01 — Executable and exact user guidance

`rtk ./bin/roundfix doctor --help` exited 0 and exposes only
`roundfix doctor`. It states that Doctor checks profiles and Repository Skill
Set independently and is offline, read-only, and mutates nothing.

The intended-reader walkthrough in `docs/user-guide/commands.md` matches the
real commands exercised in AC-02, J-03, and J-04:

- lines 54–56 specify one stdout line per check and nonzero on failure;
- lines 70–80 describe Repository Skill Set and Codex checks;
- lines 83–90 specify independent order and canonical missing-Git behavior;
- line 106 contains the exact owned-then-external `&&` chain;
- line 114 states the Repository Skill Set check is offline and read-only.

`rtk go test ./internal/cli -run
'Test(RunCommandHelp|ProfilesDocumentationContractMatchesPublicGuidance)'
-count=1` exited 0 with 11 passing tests. The help, guide, Skill pair, and real
inside/outside/mixed commands agree without undocumented steps.

## NG-01 — Non-Goals and no-mutation boundary

`rtk git -c core.fsmonitor=false diff --name-only
53a15221530e73c910e32ae665c29e9e86aef0cd^ HEAD` lists only this active Spec,
the authorized implementation/test/docs/module/Skill paths, and the five
derived Baseline artifacts.

`rtk git -c core.fsmonitor=false diff --exit-code
53a15221530e73c910e32ae665c29e9e86aef0cd^ HEAD --
docs/specs/_archived skills-lock.json skills/recommended.txt .coderabbit.yaml
.roundfixrc.yml` exited 0 with no output.

The Baseline-only path projection contains exactly:

```text
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
```

`internal/cli/doctor.go` still injects
`checkSkills func(context.Context, string)` and calls it only when the
repository root is non-empty; no `HealthChecker` ownership move occurred. The
invalid-argument and help probes prove no Doctor flag or schema was added.
Public Doctor runs executed no network repair and preserved repository and
scratch state. The final worktree status contains only this QA report and its
evidence. All QA scratch repositories were removed.
