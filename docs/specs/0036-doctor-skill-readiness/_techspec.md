---
spec: 0036-doctor-skill-readiness
prd: _prd.md
created: 2026-07-17
---

# Doctor Skill Readiness — Technical Spec

## Executive Summary

Add a deterministic repository skill checker to the existing `skills` package
and inject it into the Doctor Command as one new read-only health result. The
checker compares the current repository's `.agents/skills` tree against two
local authorities: the binary's embedded files for Roundfix-owned skills and
`skills-lock.json` for externally managed skills. It returns a structured,
sorted readiness result; the CLI owns rendering that result as one `skills:`
line and applying Doctor's existing failure exit contract.

The central trade-off is external version proof. Calling `bunx skills update`
would answer a different question, require network access, and mutate state.
Instead, Roundfix reproduces the installed skills tool's small hash contract:
SHA-256 over sorted relative path and file-byte pairs. This proves that local
content matches the repository lock without executing third-party tooling.

This design assumes Spec 0041 lands first. Its profile-aware Doctor boundary
proves the effective adapter and exact required Agent Selection Profiles. The
work below injects one independent Repository Skill Set result into that
boundary and does not preserve or recreate the legacy single-model probe.

## Project Constraints

- Identifier strategy: not applicable — the checker introduces no project-owned
  Internal Identifier; it consumes existing embedded and lock-declared skill
  names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the design is filesystem-only and
  must not add credentials, HTTP behavior, or network access. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 bind the
  profile-aware Doctor boundary, and ADR-0066 plus ADR-0072 require any Baseline
  asset behavior to remain in the Go CLI rather than a Python skill runtime.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the mechanically derived Baseline artifacts
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.normalized.json`, and
  `internal/baseline/testdata/catalog.digest`. No other protected tooling
  mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- `skills` — add repository readiness types, `skills-lock.json` decoding,
  owned-tree comparison against `Files()`, external-tree hashing, stable
  classification, and remediation metadata. This package already owns the
  embedded bundle and externally recommended names, so it remains the single
  ownership boundary.
- `internal/cli` — inject the readiness function through `doctorDependencies`,
  translate the structured result into `CheckResult{Name: "skills"}`, and
  append it to the existing deterministic Doctor output. CLI code owns wording
  and exit behavior; the checker owns filesystem facts.
- `CONTEXT.md`, `docs/user-guide`, and README — add Repository Skill Set
  vocabulary and describe the new Doctor line and update commands.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` — one
  dedicated tooling-only Task updates the authorized canonical/generated pair
  without touching other protected tooling.
- No Run Store, configuration schema, Agent Session, TUI, network, or daemon
  changes. Agent adapter and profile proof belongs exclusively to Spec 0041.
- No Context-Driven artifact inspection: Doctor proves the installed skill
  version, while `setup-context-driven audit` validates generated repository
  instructions and Baseline ADRs under Spec 0040.

## Implementation Design

### Repository readiness contract

Add public result types in `skills`:

```go
type RepositoryReadiness struct {
    OwnedRequired    int
    ExternalRequired int
    MissingOwned     []string
    OutdatedOwned    []string
    MissingExternal  []string
    OutdatedExternal []string
}

func (r RepositoryReadiness) Ready() bool

func CheckRepository(root string) (RepositoryReadiness, error)
```

`CheckRepository` is a pure observation over the supplied root and the
binary's embedded filesystem. It never discovers a different repository,
executes a command, reads environment credentials, follows a network path, or
writes. The Doctor passes `roundconfig.Loaded.GitRoot`, keeping repository
resolution in the existing configuration boundary.

The result distinguishes ownership so the CLI can provide the correct update
command. All slices are lexically sorted before return. Errors are reserved for
facts that prevent a trustworthy classification — invalid root, missing or
malformed `skills-lock.json`, invalid lock schema/hash, or unreadable files.
Missing directories and content mismatches are normal readiness results rather
than Go errors.

### Owned-skill comparison

Build an expected map from `Files()`, grouped by `File.Skill` and keyed by the
path relative to that skill directory. For every name returned by `Names()`:

1. If `<root>/.agents/skills/<name>` does not exist, append `name` to
   `MissingOwned`.
2. Walk regular files under the installed directory in stable order.
3. Mark `name` outdated if a required file is absent, an installed byte slice
   differs, or the installed tree contains an unexpected regular file.
4. Propagate permission and other I/O failures with the operation and path
   wrapped using `%w`.

Directory metadata, permissions, modification times, and empty directories are
not version material. Symlinks are not accepted as regular versioned files;
their presence makes the skill outdated rather than causing Roundfix to read
outside the repository tree.

### External lock and hash comparison

Decode `<root>/skills-lock.json` strictly enough to require:

```go
type localSkillsLock struct {
    Version int                       `json:"version"`
    Skills  map[string]localLockSkill `json:"skills"`
}

type localLockSkill struct {
    ComputedHash string `json:"computedHash"`
}
```

Version must be `1`, names must be non-empty safe directory basenames, the map
must contain every external name shipped in `Recommended()`, and every
`computedHash` must be exactly 64 lowercase hexadecimal characters. Lock
entries required by `Recommended()` must not overlap `Names()` because that
would violate ownership. Other lock entries do not expand the binary's
required Repository Skill Set and are ignored, just like unrelated installed
skill directories. The check's external required set is therefore
`Recommended()`, with hashes resolved from matching lock entries.

For each external name, compute the installed tree hash exactly as skills CLI
1.5.19 does:

```text
files = recursive regular files excluding .git and node_modules directories
files = sort by slash-normalized path relative to the skill root
digest = sha256(concat(relative_path, file_bytes) for file in files)
```

No separator is inserted between path and bytes. A missing skill directory is
`MissingExternal`; a digest mismatch is `OutdatedExternal`. Unrelated skill
directories under `.agents/skills` are ignored.

### Doctor integration and output

Extend `doctorDependencies` with:

```go
checkSkills func(root string) (skills.RepositoryReadiness, error)
```

The default calls `skills.CheckRepository`; tests inject fixed results. Doctor
runs the skill check independently of Agent Selection Profile proof outcomes
and appends its result after the profile-readiness result introduced by Spec
0041, preserving every other line and its order.

Success detail is deterministic:

```text
skills: ok (39 required: 14 Roundfix-owned, 25 external)
```

Failure detail uses only non-empty stable groups:

```text
skills: failed (missing: handoff; outdated: roundfix; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)
```

Names from both ownership groups may share the aggregate `missing` and
`outdated` labels, but remediation is selected from the structured ownership
fields. Owned command comes first, external command second, and each appears at
most once. A lock/read error becomes a failed `skills` result whose detail
names the failed operation and whose next action is the external update
command. Any failed result uses the existing `exitRunFailed` (`1`); argument or
usage errors remain `exitPreflight` (`2`). Stdout remains requested Doctor
output and stderr remains reserved for command-level load/usage failures.

### Documentation and skill synchronization

- Add **Repository Skill Set** to `CONTEXT.md` and extend **Doctor Command** to
  include it.
- Update README and Doctor command documentation with the `skills:` line,
  ownership rules, failure behavior, and read-only update guidance.
- In a separate tooling-only Task, update
  `.agents/skills/roundfix/SKILL.md` and its
  `skills/roundfix/SKILL.md` generated copy byte-identically. Do not run the
  broad `make skills-sync` mutation target because it rewrites every owned
  skill directory; use `make skills-sync-check` as read-only verification.
- Keep externally managed skill files, `skills/recommended.txt`, and
  `skills-lock.json` unchanged; the required external count remains derived
  from the shipped recommendation list.

## Coverage Map

- User stories 1–3 → repository checker, owned exact comparison, external lock
  hash comparison.
- User story 4 → stable Doctor formatter, sorted classifications, exact next
  actions, non-zero exit.
- User story 5 → local filesystem-only implementation and no-mutation tests.
- Core feature 6 → docs, canonical vocabulary, Roundfix Skill sync, and
  ownership contract verification.

## Testing Approach

- `skills` package table tests create temporary repository trees from embedded
  fixtures and cover: all current; missing owned; missing external; changed,
  added, and missing owned files; changed and added external files; mixed
  failures; malformed JSON; wrong lock version; missing/invalid hash; missing
  recommended entry; ownership overlap; symlink; unreadable path where the
  platform permits it; and deterministic ordering.
- A hash-compatibility fixture asserts a known directory produces the same
  digest as the external skills CLI algorithm. At least one current repository
  skill hash is verified in a read-only integration test or QA flow.
- Doctor buffer tests inject readiness results and assert exact success,
  owned-only, external-only, mixed, and lock-error output; existing Doctor
  output remains byte-stable except for the appended line.
- No-mutation tests snapshot relevant paths before and after Doctor and assert
  no repository, user config, `.roundfix`, or skill file is created or changed.
- QA runs the built binary against a clean disposable repository, then copies
  that fixture and exercises missing and outdated variants without changing
  the working repository.
- Run `make verify`; because filesystem behavior and Doctor dependency
  injection change, also run `go test -race ./...`.

## Build Order

Cross-Spec prerequisite: implement Spec 0041 through its profile-aware Doctor
integration before starting this graph.

1. Repository skill readiness types, lock validation, owned comparison,
   external hash compatibility, and focused package tests.
2. Doctor dependency injection, rendering after profile readiness, failure
   semantics, command tests, and CLI help (depends on 1 and Spec 0041).
3. Canonical vocabulary and user documentation (depends on 2).
4. Dedicated tooling-only update of the exact authorized canonical/generated
   Roundfix Skill files and read-only sync verification (depends on 3).

## Risks & Considerations

- The hash algorithm is a compatibility contract with the installed skills
  tool. A focused fixture must guard path normalization, sort order, excluded
  directories, and the absence of separators.
- Comparing the installed owned tree against the running binary intentionally
  means a newly built dirty binary may report local skill drift until
  the authorized canonical/generated pair is synchronized and the binary is
  rebuilt; that is the desired pre-PR signal.
- A repository lock can contain stale or malicious paths. Skill names are
  validated as basenames, and the checker never follows symlinked files, so no
  declared value can escape `.agents/skills`.
- Doctor must report runtime checks even when skill readiness fails; one failed
  local subsystem must not suppress evidence from the others.
- Spec 0041 and this Spec both touch Doctor dependencies and output. Implement
  Spec 0041 first and extend its coordinator; do not maintain competing legacy
  and profile-aware runtime checks.
- The project currently declares 39 required skills. Counts are derived at
  runtime and must not be hard-coded in production behavior.

## Decisions

- Missing and outdated skills are blocking Doctor failures.
- Roundfix-owned authority is the running binary; external authority is the
  current repository's `skills-lock.json` constrained to `Recommended()`.
- The check reproduces the local hash contract and never invokes `bunx`.
- One aggregate `skills:` line preserves Doctor's compact public surface.
- No ADR is added because this is an extension of the accepted Doctor and the
  skill ownership boundary in
  `docs/agents/skill-dispatch.md`, not a new hard-to-reverse architecture
  choice.
