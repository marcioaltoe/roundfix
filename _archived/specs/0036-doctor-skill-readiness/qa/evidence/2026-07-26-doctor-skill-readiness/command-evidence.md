# Command evidence — Doctor Skill Readiness

Build: `6e5618dba5059666c891f9078806631d30502d5b`

## Static verification

| Command | Exit | Observed result |
| --- | ---: | --- |
| `rtk make verify` | 0 | 2,394 tests passed across 22 packages; four focused skill-contract tests passed; all 14 shipped Roundfix Skills passed; `bin/roundfix` built. |
| `rtk go test -race ./...` | 0 | 2,394 tests passed across 22 packages under the race detector. |
| `rtk go test ./skills -run 'Test(CheckRepository\|SkillFolderHash)' -count=1` | 0 | 26 focused repository-readiness and folder-hash tests passed. |
| `rtk go test ./skills -count=1` | 0 | 99 skills-package tests passed, including the real repository match. |
| `rtk go test ./internal/cli -run 'Test(CommandUsage\|DocumentationContract)' -count=1` | 0 | Four public usage/documentation contract tests passed. |
| `rtk make skills-sync-check` | 0 | Four Roundfix Skill synchronization guards passed. |
| `rtk ./bin/roundfix skills check` | 0 | All 14 shipped Roundfix Skills passed. |
| `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` | 0 | No output; the canonical and generated files are byte-identical. |

## Public help

Command:

```text
rtk ./bin/roundfix doctor --help
```

Exit: `0`

```text
Usage:
  roundfix doctor

Diagnoses this machine's readiness for Runs. Checks Node.js, the minimum
supported acpx version, the effective adapter, and required
Agent Selection Profiles, the Repository Skill Set, and codex runtime hygiene.
The aggregate profiles: line exact-proves every distinct tuple. The skills:
line compares Roundfix-owned artifacts with the running binary and external
artifacts with skills-lock.json. Each failure reports its next action.
Doctor is offline, read-only, and mutates nothing.
```

## Complete repository flow

The fixture contains the current project configuration, `skills-lock.json`,
and all installed project skills. It is an independent Git repository under
`/tmp/roundfix-qa-0036.KrAsVB/clean`.

Running the built command from both the fixture root and
`nested/deeper` produced the same Repository Skill Set result:

```text
skills: ok (39 required: 14 Roundfix-owned, 25 external)
```

Both commands exited `1` because this host independently fails `profiles:` and
`codex:` readiness:

```text
/Users/marcio/.local/bin/codex has an invalid or missing code signature
```

The complete-repository file digest was identical before and after both runs:

```text
c1473d4bf5c378216f816a91e316d9d7c9782ef82a51eda512c14a5e1ea3df81
```

This proves the current build's complete `14/25/39` Repository Skill Set and
nested Git-root resolution. It does not prove the acceptance criterion's
otherwise-ready overall exit `0`; that remains blocked by the host Codex
prerequisite.

## Disposable repository matrix

Every row ran the built public command:

```text
/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260726T184338Z_c754fa2057c5fc07/bin/roundfix doctor
```

The unrelated host `profiles:` and `codex:` lines are omitted below. Every
failure fixture exited `1`. The `skills:` line appeared after `profiles:` and
before `codex:` in the complete output, proving independent later checks still
ran.

| Fixture action | Observed `skills:` line |
| --- | --- |
| Move owned `roundfix` directory out of `.agents/skills` | `skills: failed (missing: roundfix; next: roundfix skills install --target project)` |
| Move external `tech-writer` directory out of `.agents/skills` | `skills: failed (missing: tech-writer; next: bunx skills experimental_install && bunx skills update -p -y)` |
| Change owned `roundfix/SKILL.md` | `skills: failed (outdated: roundfix; next: roundfix skills install --target project)` |
| Add an unexpected owned file | `skills: failed (outdated: roundfix; next: roundfix skills install --target project)` |
| Remove owned `roundfix/SKILL.md` | `skills: failed (outdated: roundfix; next: roundfix skills install --target project)` |
| Change external `tech-writer/SKILL.md` | `skills: failed (outdated: tech-writer; next: bunx skills experimental_install && bunx skills update -p -y)` |
| Missing owned `roundfix`, missing external `tech-writer`, outdated owned `qa-gate`, outdated external `testing-boss` | `skills: failed (missing: roundfix, tech-writer; outdated: qa-gate, testing-boss; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)` |
| Replace owned `roundfix/SKILL.md` with an outside-tree symlink | `skills: failed (outdated: roundfix; next: roundfix skills install --target project)` |
| Make owned `roundfix/SKILL.md` unreadable | `skills: failed (read Roundfix-owned skill artifact ".../roundfix/SKILL.md": open SKILL.md: permission denied; next: roundfix skills install --target project)` |
| Remove `skills-lock.json` | `skills: failed (read skills lock ".../skills-lock.json": ... no such file or directory; next: bunx skills experimental_install && bunx skills update -p -y)` |
| Malform `skills-lock.json` JSON | `skills: failed (decode skills lock ".../skills-lock.json": invalid character ':' after array element; next: bunx skills experimental_install && bunx skills update -p -y)` |
| Set required `tech-writer.computedHash` to `not-a-valid-hash` | `skills: failed (validate required skills lock hash ".../skills-lock.json": external skill "tech-writer" computedHash must be 64 lowercase hexadecimal characters; next: bunx skills experimental_install && bunx skills update -p -y)` |
| Add unrelated installed skill plus unrelated unsafe lock entry | `skills: ok (39 required: 14 Roundfix-owned, 25 external)` |

The unsafe-required-name branch is not a repository-controlled input in the
shipped application: the binary owns `Names()` and `Recommended()`. The current
artifact's sets passed `roundfix skills check`; the focused current-build test
`TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames` proves that an
injected unsafe required name is rejected before filesystem traversal.

## Determinism and no mutation

The mixed fixture produced the same sorted line on repeated runs. Each
remediation command appeared exactly once and owned remediation preceded
external remediation.

The mixed fixture digest before and after rerun was:

```text
2b5a1ee30fbe18f398b7e242f7e7bfb91dac836753ae436d6c32bfd827699038
```

The unrelated-extra fixture digest before and after rerun was:

```text
fabfc4128383c52d6d33e70f3d7cbb013c77ba5d97b695901624133f26386a77
```

The symlink fixture digest before and after rerun was:

```text
52cca50e7245e6a9f721d457d6abd9a43c3d68251376315a1ef20b2fab9412f9
```

The outside-tree sentinel remained:

```text
cf46cdea17f51f34cf6e21716f0d98148b59c456650adda64a72545760313040
```

Source-boundary inspection found only local observation operations in the
Repository Skill Set checker: `os.Lstat`, `os.ReadFile`, `fs.WalkDir`, and
`fs.ReadFile`. It found no `os.Write*`, `os.Create*`, `os.Remove*`,
`os.Rename`, `os.Mkdir*`, command execution, or network import/call. The two
remediation commands appear only as rendered next-action constants.

## Documentation evidence

The README, `docs/user-guide/commands.md`, `docs/user-guide/usage.md`,
`CONTEXT.md`, and both Roundfix Skill copies consistently state:

- Agent Selection Profile Readiness runs before the independent Repository
  Skill Set result.
- The running binary is authority for the 14 Roundfix-owned skills.
- `skills-lock.json` `computedHash` values are authority for 25 external
  skills.
- Missing, outdated, or invalid required state blocks Doctor with exit `1`.
- Owned remediation is `roundfix skills install --target project`.
- External remediation is
  `bunx skills experimental_install && bunx skills update -p -y`.
- Doctor is offline and read-only, never runs remediation, and ignores
  unrelated extras.

The built help and public command output matched these terms. The Roundfix
Skill additionally tells the Agent to surface the failed line and printed
next action before work continues and to require explicit workflow
authorization before remediation.

## Project Constraint and commit-scope evidence

The PRD and TechSpec both carry applicability and operative sources for
identifier strategy, authentication/HTTP, active ADR obligations, and Tooling
Authority.

- Identifier strategy is not applicable: no Internal Identifier was added.
- Authentication/HTTP is not applicable: the checker is local filesystem
  diagnosis.
- ADR-0049 and ADR-0055 remain authoritative for the preceding profile proof;
  the observed output kept `profiles:` before `skills:`.
- ADR-0066 and ADR-0072 remain intact: no Python setup runtime or new Baseline
  execution path was added.

Task 03 commit `6e5618d` changed exactly:

```text
.agents/skills/roundfix/SKILL.md
docs/specs/0036-doctor-skill-readiness/task_03.md
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
skills/roundfix/SKILL.md
```

This equals the exact bounded Tooling Authority recorded in both active Spec
artifacts plus the assigned Task file. Task 02 commit `ccfbd15` changed only:

```text
CONTEXT.md
README.md
docs/specs/0036-doctor-skill-readiness/task_02.md
docs/user-guide/commands.md
docs/user-guide/usage.md
```

The current worktree delta is confined to this QA report and its evidence
directory. No Task status or Task Graph file changed.
