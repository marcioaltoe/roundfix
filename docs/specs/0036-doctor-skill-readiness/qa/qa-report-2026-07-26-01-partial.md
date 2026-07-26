---
spec: 0036-doctor-skill-readiness
date: 2026-07-26
build: 6e5618dba5059666c891f9078806631d30502d5b
status: closed
verdict: partial
surfaces: [cli, backend, docs]
---

# QA report — Doctor Skill Readiness

## Scope and environment

Full QA gate for every PRD user story and every Task acceptance criterion.
The primary actor is a developer preparing a local Git repository for a
Roundfix Run. The public entry points are the built `roundfix doctor`,
`roundfix doctor --help`, and `roundfix skills check` commands plus the
supported README, user guide, glossary, and Roundfix Skill.

- Build: `6e5618dba5059666c891f9078806631d30502d5b`
- Branch: `roundfix/run-run_20260726T184338Z_c754fa2057c5fc07`
- Host: macOS Darwin 25.5.0 arm64
- Go: `go1.26.5 darwin/arm64`
- Runtime tools found: `bun`, `acpx`, and `codex-acp`
- Test data: current repository plus disposable Git repository copies under a
  QA-owned temporary directory
- Network: not required and not authorized by the feature contract
- Evidence directory:
  `qa/evidence/2026-07-26-doctor-skill-readiness/`

Behavior probes are selected for their risk to the readiness boundary:
repeat the read-only command, invoke it from a nested directory, remove a
required skill, change/add/remove versioned content, introduce malformed and
unsafe lock declarations, replace content with a symlink, combine ownership
failures, and add unrelated extras.

## Project Constraint audit

| # | Constraint | Applicability and operative source | Status | Evidence |
| - | --- | --- | --- | --- |
| PC-01 | Identifier strategy | Not applicable: no project-owned Internal Identifier is created; names come from the embedded bundle and `skills-lock.json`. Operative source: `docs/agents/domain.md`. | pass | PRD/TechSpec carry the required applicability and source; Task commit inspection found no identity model or generator. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |
| PC-02 | Authentication and HTTP | Not applicable: the feature is local filesystem diagnosis and adds no authentication, HTTP, or network path. Operative source: `docs/agents/cli.md`. | pass | PRD/TechSpec carry the required applicability and source; public flows and checker inspection found local reads only. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| PC-03 | Active ADR obligations | Applicable: ADR-0049 and ADR-0055 preserve Agent Selection Profile authority; ADR-0066 and ADR-0072 keep Baseline execution in Go. Operative source: `docs/agents/domain.md`. | pass | Accepted ADRs were inspected; real output kept `profiles:` authoritative before independent `skills:`, and changed paths added no Python setup runtime or competing Baseline execution. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |
| PC-04 | Tooling authority and bounded scope | Applicable: the PRD and TechSpec expressly authorize the protected Roundfix Skill pair and five mechanically derived Baseline/parity artifacts. Operative source: `docs/agents/agent-instructions.md`. | pass | `git diff-tree` for Task 03 commit `6e5618d` matched the seven exact protected/derived paths plus `task_03.md`; Task 02 touched no protected path and the current delta is QA-only. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |

## Static gate

| Command | Status | Evidence |
| --- | --- | --- |
| `rtk make verify` | pass | Exit `0`: 2,394 tests passed across 22 packages; four focused skills checks passed; all 14 shipped Roundfix Skills passed `skills check`; `bin/roundfix` built from `6e5618d-dirty`. |
| `rtk go test -race ./...` | pass | Exit `0`: 2,394 tests passed across 22 packages under the race detector. |

## Results

| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| - | --- | --- | --- | --- |
| US-01 | Missing required skill blocks readiness and names the missing skill. | Developer; CLI from a disposable repository. | pass | Owned `roundfix` and external `tech-writer` removals each produced exit `1`, the missing name, and the correct ownership-specific remediation. [Fixture matrix](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#disposable-repository-matrix). |
| US-02 | Roundfix-owned drift from the running binary is detected. | Developer after upgrade; CLI. | pass | Changed, added, and removed owned artifacts each reported `outdated: roundfix` and the owned install command. [Fixture matrix](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#disposable-repository-matrix). |
| US-03 | External installed content is compared with `skills-lock.json`. | Maintainer; CLI/backend through Doctor. | pass | The complete real tree matched all 25 hashes; external byte drift, missing lock, malformed JSON, and invalid required hash all failed with bounded diagnostics. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#complete-repository-flow). |
| US-04 | One deterministic line names failures and exact remediation. | Agent or developer; CLI. | pass | The mixed fixture repeatedly printed one sorted line with both commands exactly once in owned-then-external order. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| US-05 | Diagnosis remains entirely local and read-only. | Security-conscious developer; CLI. | pass | Complete, mixed, extra, and symlink fixture digests remained unchanged across repeated public commands; the outside sentinel was preserved and checker inspection found no write, command-execution, or network path. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| T01-AC01 | Complete Repository Skill Set reports `skills: ok` with derived `14/25/39` counts and does not fail an otherwise ready Doctor. | Developer; CLI. | blocked | The built CLI printed `skills: ok (39 required: 14 Roundfix-owned, 25 external)` identically from the fixture root and `nested/deeper`, but this host cannot provide the criterion's otherwise-ready precondition: independent `profiles:` and `codex:` checks fail because `/Users/marcio/.local/bin/codex` has an invalid code signature. Exit-zero confirmation requires reinstalling the official signed Codex binary and rerunning the unchanged command. |
| T01-AC02 | Removing owned or external skills reports `missing`, prints the ownership-specific command, and exits `1`. | Developer; CLI. | pass | Separate owned and external public flows matched the exact line and exit contract. [Fixture matrix](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#disposable-repository-matrix). |
| T01-AC03 | Changed, added, or removed versioned files report `outdated`; unrelated extras are ignored. | Developer; CLI. | pass | Owned change/add/remove and external change reported `outdated`; an unrelated installed skill plus unrelated unsafe lock entry preserved `skills: ok`. [Fixture matrix](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#disposable-repository-matrix). |
| T01-AC04 | External folder hashing matches the installed skills CLI algorithm, including ordering and exclusions. | Maintainer; CLI/backend. | pass | The complete public flow matched all 25 current lock hashes; 26 focused compatibility checks, including the pinned locale-order fixture and exclusions, passed. [Static verification](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#static-verification). |
| T01-AC05 | Malformed lock data, invalid hashes, unsafe names, unreadable artifacts, and symlinks fail deterministically without panic or traversal. | Developer; CLI. | pass | Public missing/malformed/invalid-lock, unreadable, and symlink variants failed deterministically; the outside sentinel remained unchanged. Unsafe required names are binary-owned, not repository input, and the focused current-build guard passed before traversal. [Fixture matrix](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#disposable-repository-matrix). |
| T01-AC06 | Mixed failures produce one sorted `skills:` line and each remediation exactly once in owned-then-external order. | Agent; CLI. | pass | Repeated mixed flow matched the exact deterministic line and command order. [Determinism evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| T01-AC07 | Doctor runs all existing checks despite skill failure and changes no repository, config, Run, lock, or skill path. | Developer; CLI. | pass | Full output continued through the later `codex:` check after skill failure; multiple file manifests and the outside sentinel were byte-stable. Focused current-build no-mutation tests passed. [Determinism evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| T01-AC08 | Doctor help names Agent Selection Profiles and the appended Repository Skill Set check. | Developer; CLI/docs. | pass | Built `doctor --help` exited `0` and named both terms contiguously with the offline/read-only contract. [Public help](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#public-help). |
| T02-AC01 | Supported user docs identify authorities, outcomes, remediation, and the offline/read-only guarantee. | Documentation reader; docs. | pass | README and both user-guide entry points state both authorities, `ok/failed`, exit `1`, both exact commands, ignored extras, and no automatic mutation/network. [Documentation evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#documentation-evidence). |
| T02-AC02 | All required external installed trees match current `skills-lock.json`. | Maintainer; CLI. | pass | The public complete-repository flow printed 25 external required and `skills: ok`; the full skills package and focused compatibility suite passed. [Complete flow](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#complete-repository-flow). |
| T02-AC03 | No external authorial edit is introduced; unrelated extras are neither flagged nor removed. | Maintainer; Git/CLI. | pass | Task 02 commit scope contains only glossary/user docs/Task file; the extra fixture remained byte-stable and printed `skills: ok`. [Commit scope](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |
| T02-AC04 | Task 02 changes no protected tooling path. | Maintainer; Git. | pass | `git diff-tree` for `ccfbd15` listed only `CONTEXT.md`, README, two user guides, and `task_02.md`. [Commit scope](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |
| T02-AC05 | Glossary, docs, and CLI help use Repository Skill Set, Doctor Command, and Roundfix Skill consistently. | Documentation reader; docs/CLI. | pass | Canonical cross-read matched the built help/output; four documentation-contract tests passed. [Documentation evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#documentation-evidence). |
| T03-AC01 | Roundfix Skill tells an Agent to surface failed readiness and printed remediation before work continues without claiming Doctor updates. | Agent; docs. | pass | Canonical Skill instructs surfacing the failed line/next action, says Doctor is diagnosis-only, and requires explicit workflow authorization before remediation. [Documentation evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#documentation-evidence). |
| T03-AC02 | Canonical and generated Roundfix Skill files are byte-identical. | Maintainer; docs/tooling. | pass | `cmp` exited `0`; `make skills-sync-check` passed four guards. [Static verification](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#static-verification). |
| T03-AC03 | Task 03 changed paths contain only the authorized Skill pair, five derived artifacts, and its Task file. | Maintainer; Git. | pass | `git diff-tree` for `6e5618d` matched the eight allowed paths exactly. [Commit scope](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |
| T03-AC04 | No upstream-managed or other Roundfix-owned skill changed. | Maintainer; Git/tooling. | pass | Task 03 path audit found only the Roundfix Skill pair; complete synchronization and shipped checks passed. [Static verification](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#static-verification). |
| T03-AC05 | Shipped Skill validation and complete repository gate pass. | Maintainer; CLI/static gate. | pass | Built `roundfix skills check` passed all 14 owned skills and `rtk make verify` exited `0`. [Static verification](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#static-verification). |
| CLI-01 | Public CLI contract preserves stdout/stderr placement, exit codes, deterministic order, nested-directory resolution, and independent later checks. | Developer; CLI sweep. | pass | Help exit `0`; failures exit `1`; requested one-line checks stayed on stdout in stable order; nested invocation resolved the same Git root; `codex:` remained after skill failures. [Command evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md). |
| DOC-01 | Documentation can be followed without undocumented knowledge and matches executable help/output. | New developer; docs sweep. | pass | Documented `doctor`, `doctor --help`, and `skills check` commands were run; their output and read-only/remediation boundaries matched the guides. [Documentation evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#documentation-evidence). |
| NG-01 | Doctor does not download, install, update, delete, write, or contact a network service. | Security-conscious developer; CLI/scope. | pass | Repeated public-flow digests were unchanged; checker inspection found only local read APIs and rendered remediation constants. [Determinism evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| NG-02 | Unrelated extra skills and lock entries stay ignored and are not removed. | Maintainer; CLI. | pass | Extra fixture printed `skills: ok`; its digest was identical before and after. [Determinism evidence](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#determinism-and-no-mutation). |
| NG-03 | `roundfix skills check` retains its shipped-bundle integrity role; global installations remain outside repository readiness scope. | Maintainer; CLI. | pass | `roundfix skills check` independently passed all 14 shipped Skills; isolated project fixtures determined Repository Skill Set output from their Git root. [Static verification](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#static-verification). |
| NG-04 | Doctor does not validate Context Documents/Baseline ADRs or restore the removed Python runtime. | Maintainer; CLI/Git. | pass | Public output added only `skills:`; changed-path and static guards found no Python setup runtime or Context/Baseline audit path. [Commit scope](evidence/2026-07-26-doctor-skill-readiness/command-evidence.md#project-constraint-and-commit-scope-evidence). |

## Findings

No product finding was observed. The only incomplete criterion is an external
host-readiness prerequisite, recorded below.

## Blocked and skipped

- `T01-AC01` is blocked only on the existing host Codex readiness prerequisite.
  The Repository Skill Set itself passed from both public entry locations and
  the fixture digest remained
  `c1473d4bf5c378216f816a91e316d9d7c9782ef82a51eda512c14a5e1ea3df81`
  across both runs. Unblock by reinstalling the official signed Codex binary,
  then rerun the built `roundfix doctor` in the complete fixture and require
  exit `0`.

## Coverage

- User stories: `5/5` passed.
- Task acceptance criteria: `17/18` passed; `1` blocked.
- Constraint audits: `4/4` passed.
- Surface and Non-Goal sweeps: `6/6` passed.
- Static gates: `2/2` passed.
- Behavior probes: nested invocation, repeat, owned/external missing,
  change/add/remove drift, malformed and invalid lock, unreadable file,
  symlink traversal guard, mixed ownership, and unrelated extras all
  attempted.
- Pending rows: `0`.

## Final verdict

Partial. All feature-specific stories, 17 of 18 Task criteria, constraint
audits, surface sweeps, and static gates passed; rerun the complete fixture
after restoring a valid officially signed Codex binary to close `T01-AC01`
with an overall Doctor exit `0`.
