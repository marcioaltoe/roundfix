---
spec: 0036-doctor-skill-readiness
status: active
created: 2026-07-17
surfaces: [cli, backend, docs]
---

# Doctor Skill Readiness

Roundfix depends on a repository-local set of Agent Skills, but the Doctor
Command currently validates only the runtime toolchain. A repository can pass
Doctor while a Roundfix-owned skill is stale relative to the running binary,
an externally managed skill is missing, or installed external content no
longer matches `skills-lock.json`. That creates a dangerous split: the CLI
executes one contract while the Supervisor or Agent follows another.

Doctor Skill Readiness makes the complete Repository Skill Set a first-class,
blocking readiness check. It compares the repository's installed skills with
their authoritative local sources, reports exactly which skills are missing or
outdated, and tells the user how to update them. The check is deterministic,
offline, read-only, and part of the normal `roundfix doctor` result.

Spec 0041 Agent Selection Runtime Readiness is a prerequisite. It replaces
Doctor's legacy runtime probe with the shared effective-adapter and exact
Agent Selection Profile proof. This Spec appends only the independent
Repository Skill Set result and must not recreate or bypass profile proof.

## Project Constraints

- Identifier strategy: not applicable — this feature creates no project-owned
  Internal Identifier or application identity; skill names remain sourced from
  the embedded bundle and `skills-lock.json`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Doctor reads only local repository
  state and must not add an authentication provider, HTTP route, or network
  call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 keep Agent
  Selection Profile proof authoritative, while ADR-0066 and ADR-0072 keep
  Baseline execution in the Go CLI and prohibit restoring the removed Python
  runtime. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the mechanically derived Baseline artifacts
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.normalized.json`, and
  `internal/baseline/testdata/catalog.digest`, and the mechanically derived
  parity artifacts
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Make `roundfix doctor` prove that every required repository skill is
  installed at the expected version before reporting the repository ready.
- Treat the Roundfix binary's embedded skill bundle as authoritative for the
  14 Roundfix-owned skills.
- Treat the repository's `skills-lock.json` names and `computedHash` values as
  authoritative for externally managed skills.
- Report one stable `skills:` line that distinguishes missing and outdated
  skills and provides the exact relevant update command or commands.
- Keep Doctor diagnosis-only: no downloads, network calls, installs, file
  writes, or lock-file updates.
- Keep `CONTEXT.md`, user guidance, and the canonical Roundfix Skill aligned
  with the shipped Doctor behavior.

## User Stories

1. As a developer preparing a repository, I want Doctor to fail when a
   required skill is missing, so an autonomous Run cannot begin with an
   incomplete workflow contract.
2. As a developer after upgrading Roundfix, I want Doctor to detect when a
   Roundfix-owned repository skill differs from the running binary, so I know
   to reinstall the shipped bundle.
3. As a maintainer updating external skills, I want Doctor to compare the
   installed content with `skills-lock.json`, so silent local drift or an
   incomplete update is visible immediately.
4. As an Agent or developer reading Doctor output, I want the failing skill
   names and exact remediation commands on one deterministic line, so I can
   act without researching repository setup.
5. As a security-conscious developer, I want the check to remain entirely
   local and read-only, so diagnosis cannot download or execute unreviewed
   skill updates.

## Core Features

1. **Repository Skill Set discovery.** Doctor resolves the current Git root,
   reads `<root>/.agents/skills/`, and requires the union of the Roundfix-owned
   embedded skill names and the binary's externally recommended names. Every
   required external name must have a matching declaration in
   `<root>/skills-lock.json`.
2. **Owned-skill version proof.** Every file in each Roundfix-owned installed
   skill must match the running binary's embedded artifact byte-for-byte.
   Missing files, changed bytes, or unexpected files make that skill outdated;
   a missing skill directory makes it missing.
3. **External-skill version proof.** Every external skill directory must hash
   to its lock entry's `computedHash`. Hashing follows the installed skills
   tool's deterministic local algorithm: sort files by slash-normalized path,
   then SHA-256 each relative path followed by its bytes, excluding `.git` and
   `node_modules` subtrees.
4. **Blocking Doctor result.** A current Repository Skill Set prints
   `skills: ok` with owned, external, and total counts. A missing, malformed,
   or mismatched declaration, missing skill, or outdated skill prints
   `skills: failed`, includes stable sorted names, and makes Doctor exit with
   its existing run-failed exit code.
5. **Actionable remediation.** Owned drift points to
   `roundfix skills install --target project`; external drift points to
   `bunx skills experimental_install && bunx skills update -p -y`. When both
   groups fail, the single Doctor line includes both commands in that order.
6. **Synchronized guidance.** Command help, user documentation, the canonical
   Roundfix Skill, its embedded copy, and the canonical glossary describe the
   new check and preserve the repository's skill-ownership boundary.

## User Experience

`roundfix doctor` keeps its existing one-line-per-check format and adds one
`skills:` line. Successful output identifies how many Roundfix-owned and
external skills were proven. Failure output uses sorted, comma-separated
`missing` and `outdated` groups and appends `next:` with only the commands
needed for the failing ownership groups. Doctor still runs and prints the
other readiness checks, then exits non-zero if any check failed.

The command never updates skills automatically. A developer or orchestrating
Agent can surface the failure, run the printed command after explicit workflow
authorization, and rerun Doctor to prove the result.

## Non-Goals / Out of Scope

- Downloading, installing, updating, deleting, or otherwise modifying skills.
- Contacting GitHub, a skills registry, or any other network service to decide
  whether a newer upstream version exists.
- Reporting unrelated extra skill directories or lock entries that are not
  members of the required Repository Skill Set, or offering to remove them.
- Changing the behavior of `roundfix skills check`, which validates the
  integrity of the bundle shipped by the binary rather than one repository's
  installation.
- Replacing `skills-lock.json` with a Roundfix-specific lock format.
- Validating global Agent Skill installations outside the current repository.
- Validating `AGENTS.md`, Context Documents, agent guides, or Baseline ADRs;
  those repository-document contracts belong to `setup-context-driven audit`
  and Spec 0040.
- Changing Baseline asset synchronization or restoring the removed Python
  setup runtime. Baseline asset behavior belongs to the Go CLI under ADR-0066
  and ADR-0072.

## Success Metrics

- A repository with all 14 owned skills byte-equal to the binary and all 25
  external skills matching `skills-lock.json` reports `skills: ok` for the
  current derived total of 39 required skills and does
  not fail Doctor on skill readiness.
- Removing one required skill makes Doctor name it as missing and exit
  non-zero without creating or changing any file.
- Changing one byte or adding one file to an owned or external skill makes
  Doctor name it as outdated and print the correct ownership-specific update
  command.
- A missing or malformed `skills-lock.json`, an empty or malformed
  `computedHash`, or an unreadable required artifact produces a deterministic
  failure with a useful next action rather than a panic.
- Focused tests cover clean, missing, outdated, mixed-ownership, malformed
  lock, and no-mutation behavior; the full repository verification gate and
  race suite pass.

## Decisions

- Missing or outdated required skills are blocking readiness failures.
- The current repository is the only validation scope: owned skills match the
  running binary, while external skills match the current repository's lock.
- Doctor emits one aggregate `skills:` line and remains read-only and offline.
- Unexpected skills outside the declared set are ignored; Doctor never offers
  removal.
- The local deterministic hash contract is implemented directly rather than
  invoking the external skills CLI.
- This change extends the existing Doctor Command and the skill ownership
  contract in `docs/agents/skill-dispatch.md`; it does not require a new
  architectural decision record.

## Dependencies

- Implement Spec 0041 Agent Selection Runtime Readiness first.
- Append `skills:` after the profile-aware readiness result delivered by Spec
  0041 while preserving all other deterministic Doctor results.
- Reuse Doctor's dependency-injection and aggregate-result boundaries; do not
  retain a parallel legacy model probe or add a second profile prover here.

## Open Questions

None.
