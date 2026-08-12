# Command evidence — Doctor Skill Readiness Hardening

Build: `14e03d61784c5cd9714ba89b25a22bcb7c223b09`

This file records fresh commands, exit codes, focused outputs, and state
comparisons for the 2026-07-26 QA gate. Sections are populated as their
corresponding matrix rows execute.

## Preconditions and Project Constraints

All four Task files declare `status: completed`. Before the QA artifacts were
created, `rtk git -c core.fsmonitor=false status --short --branch` reported
only the Run Branch header and no worktree changes.

The PRD and TechSpec each contain a Project Constraints section. Identifier
strategy and authentication/HTTP are explicitly not applicable with reasons;
the sections cite `docs/agents/domain.md` and `docs/agents/cli.md`. Active ADR
obligations are applicable and name ADR-0049, ADR-0055, ADR-0066, and ADR-0072.
Tooling authority records the maintainer's 2026-07-26 authorization for
`golang.org/x/text` and bounds protected changes to `go.mod` and `go.sum`.

The four implementation commits and fresh `diff-tree` paths were:

```text
6fd1d32
  docs/specs/0050-doctor-skill-readiness-hardening/task_01.md
  go.mod
  go.sum
93b3151
  docs/specs/0050-doctor-skill-readiness-hardening/task_02.md
  internal/baseline/skills_restore.go
  internal/skillhash/hash.go
  internal/skillhash/hash_test.go
  skills/skills.go
be91ebc
  docs/specs/0050-doctor-skill-readiness-hardening/task_03.md
  skills/repository.go
  skills/repository_test.go
  skills/skills.go
  skills/skills_test.go
14e03d6
  CONTEXT.md
  docs/specs/0050-doctor-skill-readiness-hardening/task_04.md
  internal/cli/cli_test.go
  internal/cli/doctor.go
```

Every `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only
-r <commit>` command exited `0`. Only Task 01 changed protected tooling files,
and its paths match the express authorization plus its own Task file.

## Task 01 authorized dependency

Commands:

```text
rtk go list -m -f '{{.Path}} {{.Version}}' golang.org/x/text
rtk go mod verify
rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r 6fd1d32
```

All exited `0`. The module query returned
`golang.org/x/text v0.40.0`; module verification returned
`all modules verified`; the commit changed only `go.mod`, `go.sum`, and
`task_01.md`.

## Root-anchored symlink rejection

The built CLI ran in disposable Git repositories derived from a known-ready
Repository Skill Set. Each command exited `1`, completed all independent
checks, and emitted the expected blocking skills result:

```text
.agents link:
skills: failed (inspect repository skill directory ".../.agents": symbolic links are not supported; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)

.agents/skills link:
skills: failed (inspect repository skill directory ".../.agents/skills": symbolic links are not supported; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)

skills-lock.json link:
skills: failed (inspect skills lock ".../skills-lock.json": symbolic links are not supported; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)
```

Nested link and FIFO entries also prevented readiness and produced
`skills: failed (outdated: roundfix; next: roundfix skills install --target
project)`. The owning regression suite gave the stronger internal
classification proof:

```text
rtk go test ./skills -run 'TestCheckRepository' -count=1
Go test: 26 passed in 1 packages
```

After all public commands, recursive `diff -qr` comparisons for the external
`.agents` and `.agents/skills` targets and `cmp` for the external lock all
exited `0`.

## Unicode hash compatibility

The exact pinned oracle ran in the focused suite and passed with the documented
skills CLI 1.5.19 digest
`2a46b6d704729eafc0148969028b9cc4030813059e1f7524def2f38b433011d4`.
The suite includes `_a` before `-a`, digits, case, composed and decomposed
Unicode, nested paths, and path/content negative companions.

On this macOS filesystem, case-only and canonical-equivalent names cannot
coexist. A representable ten-file subset was hashed independently with Node's
optionless `String.localeCompare`; its digest was
`57cb017b98f0f26a3782e0a4660885f61b74751d53633dd62f765e3ab301aa4b`.
The built Doctor accepted that locked external skill twice:

```text
skills: ok (39 required: 14 Roundfix-owned, 25 external)
```

Renaming `a.md` to `changed.md` in a copy while retaining the same lock changed
the public result to:

```text
skills: failed (outdated: agentic-cli-design; next: bunx skills experimental_install && bunx skills update -p -y)
```

## Shared Doctor-Baseline hash authority

The first public preview without `--source-dir` failed because the sandbox
could not resolve `github.com`. The documented offline path was then used with
the clean local checkout at `/Users/marcio/dev/skills`, which contains the
profile's immutable commit.

In a disposable clone:

```text
roundfix baseline skills restore --profile go-cli-tui --skill agentic-cli-design --source-dir /Users/marcio/dev/skills --repo <clone>
```

Preview exited `3` with Plan Digest
`23a88969138a8072e167c3a5e3fe6689bede48010506ac12a13dbb38551e966f`
and one lock-entry update. Confirming that exact digest exited `0` with
`Baseline skills restore: applied`. A fresh unconfirmed rerun exited `0` with
`Baseline skills restore: no changes`. Doctor then reported the restored
repository's `skills: ok`.

Source inspection found the two production calls:

```text
skills/skills.go:190: return skillhash.Sum(files), nil
internal/baseline/skills_restore.go:1251: return skillhash.Sum(hashFiles)
```

No copied lowercase path comparator remains. Task 01's module-file object state
is unchanged in later implementation commits, and the focused shared suite
passed 27 tests across the hash, skills, and Baseline packages.

## Doctor root and remediation

From a directory with no Git root, the built command exited `1`, still rendered
all independent checks in Node, acpx, adapter, profiles, skills, codex order,
and printed:

```text
skills: failed (repository skill check requires a Git repository; next: run roundfix doctor from a Git repository)
```

The three unclassified authority errors also exited `1` and printed both
ownership remediations in owned-then-external order. The focused public-runner
suite passed 17 Doctor tests, including the empty-root call-count and exact
output assertions.

## Public Doctor no-mutation

A clean disposable clone had an empty `git status` before Doctor. The built
command ran twice; both runs produced the same ordered output and
`skills: ok (39 required: 14 Roundfix-owned, 25 external)`. A fresh
`rtk git -c core.fsmonitor=false status --short` remained empty after the
second run.

`rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` passed 17 tests. Its
real-checker public-runner fixture snapshots repository files, User Config,
user and repository `.roundfix`, lock, and skill trees and proves every
snapshot byte-identical after execution. The exact full `rtk make verify`
also passed.

## Repository classification and test ownership

The built Doctor reported `skills: ok` in the real repository and clean
disposable copies. Removing `autoresearch` and `agentic-cli-design` produced
the stable sorted result:

```text
skills: failed (missing: agentic-cli-design, autoresearch; next: bunx skills experimental_install && bunx skills update -p -y)
```

The Unicode path-change fixture produced `outdated: agentic-cli-design`.
`rtk go test ./skills -count=1` passed all 106 package tests.

Every `TestCheckRepository*` definition is in `skills/repository_test.go`;
every `TestSkillFolderHash*` definition is in `skills/skills_test.go`.

## Doctor CLI contract

Ready, missing-root, authority-error, and repeated Doctor executions retained
the `"; next: "` boundary and exit `1` for readiness failures. Invalid
arguments retained exit `2`:

```text
rtk ./bin/roundfix doctor unexpected
roundfix: doctor failed: unexpected argument "unexpected"
Run 'roundfix doctor --help' for usage.
```

`doctor --format json` also exited `2` with `flag provided but not defined:
-format`, proving no new machine schema. The 17 focused public CLI tests passed
their exact stdout, stderr, ordering, and exit-code assertions.

## Contract wording and bounded scope

`CONTEXT.md:370` states that Doctor "reports the detected acpx version against
the minimum". Task commit `14e03d6` changed only `CONTEXT.md`, `task_04.md`,
`internal/cli/cli_test.go`, and `internal/cli/doctor.go`.

The cumulative implementation diff from planning base `b71957d` contains only
the four Task files, `CONTEXT.md`, the authorized module files, and the
intended hash, Repository Skill Set, Baseline, and Doctor sources/tests. It
contains no archived Spec, upstream-managed skill, lock, or recommended-list
path.

## Non-Goals: history and archive

`rtk git merge-base --is-ancestor b71957d 14e03d6` exited `0`. The archived
Spec 0036 tree object is identical at base and build:

```text
758792d870b4ffd93e3b6984b71cbb185f7d1dd9
758792d870b4ffd93e3b6984b71cbb185f7d1dd9
```

## Non-Goals: Baseline and managed files

The base and build object IDs are identical for `skills-lock.json`,
`skills/recommended.txt`, and `.agents/skills`; the scoped `git diff --quiet`
exited `0`.

No Baseline asset-sync production path changed. The existing preservation and
idempotence regression passed:

```text
rtk go test ./internal/baseline -run 'TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent' -count=1
Go test: 1 passed in 1 packages
```

## Non-Goals: output and membership

`doctor --help` exposes only `roundfix doctor`; it describes an offline,
read-only command and no format option. `doctor --format json` exited `2`.
Every readiness failure retained the human `"; next: "` boundary.

The base/build object IDs for `skills-lock.json` and
`skills/recommended.txt` match exactly. Current public output remains
`39 required: 14 Roundfix-owned, 25 external`.

## Non-Goal: offline Doctor

Every public Doctor run completed without a network diagnostic. Source
inspection found `doctorDeps.checkSkills = skills.CheckRepository` and no
process launch in `skills/repository.go` or `internal/skillhash/hash.go`.
The existing Node and acpx machine-health checks remain independent; no Node,
Bun, or external skills CLI process was added for hash compatibility.

## CLI surface sweep

The built binary covered:

- `doctor --help`: exit `0`, current offline/read-only contract.
- `doctor`: ready, missing, outdated, three symlink authorities, nested link,
  special entry, missing Git root, and repeat execution.
- `doctor unexpected` and `doctor --format json`: exit `2`.
- `skills check`: exit `0`, all 14 shipped Roundfix skills passed.
- `baseline skills restore`: offline preview exit `3`, confirmed apply exit
  `0`, idempotent rerun exit `0`.

Fresh reruns and Git/directory comparisons confirmed stable state.

## Backend surface sweep

All fresh suites passed:

```text
rtk go test ./internal/skillhash ./skills ./internal/baseline -run 'Test(Sum|SkillFolderHash|CheckRepositoryMatchesRealRepository|SkillsRestore)' -count=1
Go test: 27 passed in 3 packages

rtk go test ./skills -run 'TestCheckRepository' -count=1
Go test: 26 passed in 1 packages

rtk go test ./internal/cli -run 'TestRunDoctor' -count=1
Go test: 17 passed in 1 packages

rtk go test -race ./internal/skillhash ./skills ./internal/baseline ./internal/cli -run 'Test(Sum|SkillFolderHash|CheckRepository|SkillsRestore|RunDoctor)' -count=1
Go test: 69 passed in 4 packages

rtk go test ./skills -count=1
Go test: 106 passed in 1 packages
```

## Documentation surface sweep

`CONTEXT.md` defines Doctor with detected-versus-minimum acpx wording,
independent Repository Skill Set evaluation, next actions, and no mutation.
`roundfix doctor --help` matched that contract, and real output included the
detected `acpx: ok (0.12.1 >= 0.12.0)` result plus deterministic remediation.

The full Verification's Repository Skill Set check passed the canonical
Roundfix Skill synchronization for all 14 shipped skills.

## Static gate

Command:

```text
rtk make verify
```

Exit: `0`.

Focused output:

```text
rtk go test ./...
Go test: 2409 passed in 23 packages
rtk go test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'
Go test: 4 passed in 1 packages
rtk go run -buildvcs=false ./cmd/roundfix skills check
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
rtk go build -buildvcs=false ... -o bin/roundfix ./cmd/roundfix
```

The exact repository Verification passed on the QA build.

Final artifact postflight:

```text
rtk rg -n '[[:blank:]]+$' docs/specs/0050-doctor-skill-readiness-hardening/qa
exit 1, no trailing-whitespace matches

rtk jq empty docs/specs/0050-doctor-skill-readiness-hardening/qa/evidence/2026-07-26-doctor-skill-readiness/unicode-skills-lock.json
exit 0

filesystem-compatible oracle recheck
10 files, 57cb017b98f0f26a3782e0a4660885f61b74751d53633dd62f765e3ab301aa4b

rtk git -c core.fsmonitor=false diff --check
exit 0
```

`git status --short --untracked-files=all` listed only the dated QA report and
its evidence directory.
