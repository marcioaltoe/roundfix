---
spec: 0050-doctor-skill-readiness-hardening
status: active
surfaces: [cli, backend, docs]
---

# Doctor Skill Readiness Hardening

Spec 0036 made Repository Skill Set readiness a blocking Doctor Command
result, but review found three gaps in that implementation. Repository-level
symbolic links can redirect skill or lock reads beyond the intended Git root,
the external hash order only approximates the installed skills CLI's
`String.localeCompare` order, and the Doctor coordinator can substitute its
process working directory when configuration did not resolve a Git root.

This corrective Spec closes those gaps without reopening or rewriting the
archived Spec 0036. It also restores the missing executable evidence for
Doctor's no-mutation contract and removes duplicated hash-order logic from the
Repository Skill Set and Baseline restoration paths.

## Project Constraints

- Identifier strategy: not applicable — this correction creates no
  project-owned Internal Identifier or application identity. Skill names and
  paths remain governed by the existing Repository Skill Set contract. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change is local filesystem,
  hashing, CLI coordination, documentation, or test behavior. No credential,
  HTTP route, or network policy changes. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 keep Agent
  Selection Profile proof independent from Repository Skill Set readiness;
  ADR-0066 and ADR-0072 keep Baseline skill restoration and asset
  synchronization in the Go CLI. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26 the maintainer expressly
  authorized this Spec to add `golang.org/x/text` for Unicode collation.
  Bounded protected files: `go.mod` and `go.sum`. No other protected tooling
  mutation is authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- Keep every Repository Skill Set read anchored to the supplied Git root and
  reject symbolic links at repository-owned skill-tree and lock-authority
  boundaries.
- Match the installed skills CLI's path collation for punctuation, case, and
  Unicode instead of approximating it with lowercase byte ordering.
- Make the Repository Skill Set and Baseline restoration paths consume one
  shared external-skill hash implementation.
- Make the Doctor Command use the resolved Git root as its sole repository
  authority and print deterministic remediation for every readiness error.
- Prove through public CLI execution that Doctor does not mutate repository,
  User Config, `.roundfix`, lock, or skill state.
- Keep the canonical Doctor Command definition aligned with its still-shipped
  minimum-supported acpx version report.

## User Stories

1. As a security-conscious maintainer, I want Doctor to reject a symlinked
   `.agents`, `.agents/skills`, or `skills-lock.json`, so readiness cannot be
   controlled by files outside the repository trust boundary.
2. As a maintainer with valid punctuation or Unicode paths inside a skill, I
   want Doctor to accept the same folder hash produced by the installed skills
   CLI, so a valid lock does not become a false outdated result.
3. As a Baseline maintainer, I want restored external lock hashes to use the
   same ordering implementation as Doctor, so the two local authorities cannot
   drift.
4. As an Agent invoking Doctor, I want a missing Git root or unexpected skill
   checker failure to produce a useful next action without silently checking a
   different directory.
5. As a reviewer, I want package and CLI tests to exercise the real
   filesystem and public runner at their owning layers, so passing tests prove
   the security and no-mutation contracts.

## Core Features

1. **Root-anchored skill inspection.** Repository readiness opens one anchored
   repository root, validates each skill-tree ancestor and the lock authority
   without following links, and confines later reads and walks to that root.
2. **Compatible path collation.** One internal hash component uses explicit
   American English Unicode collation compatible with the installed skills
   CLI's resolved `String.localeCompare` behavior and preserves the existing
   slash-normalized path-plus-bytes digest shape.
3. **Shared hash authority.** `skills.SkillFolderHash` and Baseline external
   lock generation adapt their file inputs to the same internal hash
   component; neither carries a copied comparator.
4. **Deterministic Doctor failures.** Doctor never substitutes `os.Getwd()` for
   a missing `roundconfig.Loaded.GitRoot`. Missing-root and unclassified
   checker errors remain blocking, keep independent checks eager, and include
   deterministic remediation.
5. **Behavioral regression coverage.** Repository tests live with
   `repository.go`; hash compatibility tests include punctuation and Unicode;
   a public `Run(["doctor"])` test snapshots relevant paths before and after a
   real repository skill check.
6. **Contract reconciliation.** `CONTEXT.md` restores the Doctor Command's
   detected-acpx-version wording. The already implemented repository-owned
   Baseline digest preservation remains cited as pre-existing evidence rather
   than being rebuilt.

## Non-Goals

- Rebuilding, rebasing, force-pushing, or otherwise rewriting the current
  branch or its existing commits.
- Editing any artifact under
  `docs/specs/_archived/0036-doctor-skill-readiness/`.
- Reimplementing `sync-setups`; repository-owned digest preservation already
  exists in the Go Baseline asset synchronization path and its regression
  suite.
- Editing upstream-managed skills, `skills-lock.json`, or
  `skills/recommended.txt`.
- Changing the documented `"; next: "` Doctor output boundary or inventing a
  machine-readable Doctor schema.
- Changing required skill membership or current illustrative count examples.
- Invoking Node, Bun, or the external skills CLI from Doctor.

## Success Metrics

- Repository fixtures with symlinked `.agents`, `.agents/skills`, or
  `skills-lock.json` fail readiness without reading the external targets.
- A pinned punctuation, case, and Unicode corpus produces the same SHA-256
  digest as skills CLI 1.5.19's `String.localeCompare` algorithm.
- Doctor and Baseline external hashes are generated by one shared
  implementation and retain compatibility with all current lock entries.
- Doctor receives only `roundconfig.Loaded.GitRoot` as the repository
  readiness root; an empty value never calls the checker with the process
  working directory.
- A public Doctor execution leaves all snapshotted repository and user paths
  byte-identical.
- Focused tests, the race suite for affected packages, and `rtk make verify`
  pass before QA.

## Decisions

- Keep the archived Spec 0036 immutable and deliver review corrections through
  this follow-up Spec.
- Use `golang.org/x/text/collate` with
  `language.AmericanEnglish` because the supported Node and Bun environments
  resolve the installed skills CLI's optionless `localeCompare` to `en-US`.
- Add one internal hash package instead of exporting Baseline-specific file
  types or maintaining copied comparators.
- Keep Doctor read-only and offline; compatibility is implemented in Go rather
  than delegated to an external process.
- Add no ADR because the correction narrows existing trust, compatibility, and
  ownership contracts rather than introducing a new hard-to-reverse product
  boundary.

