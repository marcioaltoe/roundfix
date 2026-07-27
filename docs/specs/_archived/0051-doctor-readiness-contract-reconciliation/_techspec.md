---
spec: 0051-doctor-readiness-contract-reconciliation
prd: _prd.md
created: 2026-07-26
---

# Doctor Readiness Contract Reconciliation — Technical Spec

## Executive Summary

Reconcile four contracts left incomplete after Repository Skill Set
hardening: Go module ownership, cancellable filesystem inspection,
deterministic collation ties, and Doctor coordination. The existing package
boundaries stay intact: `skills` owns repository facts,
`internal/skillhash` owns pure digest ordering, and `internal/cli` owns Doctor
coordination and rendering. The primary trade-off is that byte-order
tie-breaking deliberately defines behavior where optionless
`String.localeCompare` reports equality and the external CLI leaves final
order dependent on traversal; this yields one reproducible local authority
instead of mirroring an undefined tie.

## Project Constraints

- Identifier strategy: not applicable — existing skill names and normalized
  repository-relative paths remain the only identifiers. No project-owned
  identity is added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the implementation remains local,
  offline, and filesystem-only and changes no authentication, authorization,
  transport, or HTTP contract. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 keep exact Agent
  Selection Profile proof independent and eager; ADR-0066 and ADR-0072 keep
  the shared Baseline hash consumer in Go without reopening asset
  synchronization; ADR-0077 requires this readable constraint snapshot.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization was given
  on 2026-07-26 to update `go.mod` and `go.sum` through Go tooling for the
  `golang.org/x/text` collation dependency, and to keep
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` aligned
  with shipped CLI behavior. After QA proved that this canonical edit
  invalidated the repository's authorial Baseline snapshot invariant, the
  maintainer also authorized the exact derived refresh of
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. These nine paths
  are the complete protected tooling boundary for this Spec; no other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- `skills` keeps Repository Skill Set ownership and filesystem policy. Its
  public blocking entry points become context-first, and cancellation flows
  through anchored repository reads, owned comparison, and external hashing.
- `internal/skillhash` keeps the shared pure digest algorithm. Its comparator
  becomes a total order by applying raw normalized path order only when the
  American English collator reports equality for distinct strings.
- `internal/baseline` remains a consumer of the shared hash. The only Baseline
  mutation is a mechanically derived refresh of the TypeScript/Bun setup
  snapshot and the catalog/parity metadata that authenticates that snapshot.
  Baseline behavior, source setup authority, lock files, and synchronization
  state do not change.
- `internal/cli` separates the profile-proof working directory from the
  repository readiness root, passes the command context into the injected
  skills seam, and owns canonical detail and fail-closed remediation text.
- User guidance and the Roundfix Skill pair mirror the shipped text contract.
  No new package, `HealthChecker` method, command, flag, or schema is added.

## Implementation Design

### Interfaces

The public filesystem APIs become context-first:

```go
func CheckRepository(
    ctx context.Context,
    root string,
) (RepositoryReadiness, error)

func SkillFolderHash(
    ctx context.Context,
    root string,
) (string, error)
```

The Doctor seam mirrors the public contract:

```go
checkSkills func(
    context.Context,
    string,
) (skills.RepositoryReadiness, error)
```

The pure hash operation keeps its existing shape:

```go
type File struct {
    Path    string
    Content []byte
}

func Sum(files []File) string
```

### Context propagation and errors

`CheckRepository` checks `ctx.Err()` before starting and passes the same
context to internal repository inspection. Walk callbacks and file-read loops
check cancellation before processing the next entry and before reading file
content. `SkillFolderHash` applies the same checks around collection and
delegates only the completed file slice to the pure hash operation.

Cancellation is wrapped with the current operation and path through
`RepositoryReadinessError` where repository ownership is known. Wrapping uses
`%w` semantics so `errors.Is(err, context.Canceled)` and
`errors.Is(err, context.DeadlineExceeded)` remain true. No timeout, background
context, goroutine, polling loop, or test-only production hook is introduced.

Embedded bundle enumeration remains synchronous package data and does not
receive a separate context API. The cancellation boundary covers the supplied
repository and external filesystem walks that can grow with user-controlled
content.

### Deterministic hash order

Every path is slash-normalized before comparison. Ordering uses:

1. `collate.New(language.AmericanEnglish).CompareString`.
2. When the comparison is non-zero, its sign decides order.
3. When the comparison is zero and the normalized strings differ, ordinary Go
   string comparison decides order.

The second comparison turns the stable partial order into a total order. It
preserves `_a` before `-a` and the existing punctuation/case/Unicode
compatibility corpus while making precomposed and decomposed canonical
equivalents independent of caller order. The test oracle pins the resulting
digest and separately permutes the tie corpus; it does not invoke Node or Bun
from production or tests.

### Doctor working directories

Doctor resolves two values:

- `profileWorkDir`: trimmed Git root when present; otherwise the process
  working directory obtained with `os.Getwd()`.
- `repositoryRoot`: the trimmed Git root only.

Profile proof always receives `profileWorkDir`. When `repositoryRoot` is
empty, Doctor does not call `checkSkills` and emits a failed skills result with
detail `Repository Skill Set readiness requires a Git repository` plus the
existing run-from-Git next action. A `Getwd` error follows the pre-existing
command failure behavior and must not be disguised as a Repository Skill Set
result.

Independent machine checks, adapter proof, profile proof, skills readiness,
and Codex hygiene remain eager and keep their current output order.

### Ownership and remediation

Every `skills-lock.json` inspection error, including a symbolic link, carries
`RepositoryOwnershipExternal`. Shared ancestor errors remain unclassified
because they block both ownership groups. Owned and external skill-root
symlinks keep their current distinct classifications.

Doctor derives required repair groups from readiness fields plus typed error
ownership. An unclassified error with no classified group conservatively
selects both actions. Multiple actions are joined with ` && ` in
owned-then-external order:

```text
roundfix skills install --target project && bunx skills experimental_install && bunx skills update -p -y
```

The existing `"; next: "` separator between detail and the `NextAction` field
does not change. Only composition inside `NextAction` changes, making the
printed shell expression fail closed.

### Module and naming hygiene

The tooling Task runs `go mod tidy` rather than editing dependency metadata by
hand. The expected result moves `golang.org/x/text` to the direct requirement
block, removes stale module sums, and adds the sums selected by the current Go
module graph without changing dependency versions.

Repository helper parameters that currently shadow the imported `path`
package are renamed to describe their role, such as `relativePath` and
`lockPath`. This is a bounded readability correction inside the same
Repository Skill Set slice.

### Documentation and skill synchronization

Update the Doctor user guide for:

- context-aware, read-only Repository Skill Set readiness;
- canonical missing-Git terminology;
- the `&&`-joined mixed remediation chain.

After the CLI and user guide settle, update only the canonical Roundfix Skill
and its shipped copy in a dedicated tooling Task so
`.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` remain
byte-identical. Do not combine that protected mutation with code or user-guide
changes.

The canonical skill edit changes the repository-owned content digest embedded
in the TypeScript/Bun Baseline setup snapshot. A subsequent isolated Task
updates that one digest and regenerates only the directly derived setup digest,
normalized catalog/digest, parity asset-sync fixture, and parity manifest.
The Task does not change the Baseline source setup, implementation, schemas, or
other assets. Do not edit upstream-managed skills, recommendation or lock
authorities, archived Specs, `.coderabbit.yaml`, or `.roundfixrc.yml`.

## Data Models

No persisted model or schema changes. `RepositoryReadiness`,
`RepositoryReadinessError`, and `CheckResult` keep their existing fields and
semantics.

## API Contracts

The Go signatures for blocking Repository Skill Set inspection add
`context.Context` as the first parameter. Doctor's command-line arguments,
line ordering, status vocabulary, exit codes, stdout/stderr split, and
`"; next: "` detail boundary remain unchanged.

## Coverage Map

- Goal: cancellable filesystem inspection; User Story 1 → context-first
  `skills` APIs, walk/read cancellation checks, and error-chain tests.
- Goal: deterministic collation ties; User Story 2 →
  `internal/skillhash` total comparator and permutation corpus.
- Goal: separated working directories; User Story 3 → Doctor coordinator and
  public command tests outside Git.
- Goal: canonical fail-closed failures; User Story 4 → external lock ownership,
  canonical detail, and exact mixed-remediation output tests.
- Goal: tool-produced metadata; User Story 5 → bounded module Task and tidy
  postflight.
- Goal: synchronized guidance → Doctor user guide and Roundfix Skill sync.
- Goal: repository gate integrity → derived TypeScript/Bun snapshot,
  catalog, and parity metadata agree with the synchronized Roundfix Skill.

## Integration Points

- `golang.org/x/text/collate` and `language` remain the in-process collation
  authority.
- `os.Root` remains the repository confinement boundary.
- `os.Getwd` restores only the Agent Selection Profile proof fallback.
- Baseline restoration continues to consume `internal/skillhash.Sum` without
  changing its asset or transaction contracts.

## Testing Approach

The test invariants and owners are:

- `internal/skillhash/hash_test.go`: the same normalized file set produces one
  digest for every input permutation, including distinct canonically
  equivalent paths. Boundary IN is pure ordering and hashing; filesystem
  collection is outside.
- `skills/skills_test.go`: `SkillFolderHash` honors pre-cancellation and
  cancellation observed while walking or reading. Boundary IN is the public
  folder-hash filesystem contract; Doctor rendering is outside.
- `skills/repository_test.go`: `CheckRepository` preserves cancellation in the
  error chain and assigns external ownership to a symlinked lock. Boundary IN
  is anchored Repository Skill Set inspection; CLI text is outside.
- `internal/cli/doctor_test.go`: outside Git, profile proof receives the real
  process working directory while `checkSkills` is not called; mixed failures
  print one exact fail-closed chain and independent checks remain eager.
  Boundary IN is public Doctor coordination and output; package hashing is
  covered below.

Focused tests use real temporary filesystems and the existing injected Doctor
collaborators only at machine/runtime boundaries. No assertion repeats a value
written into a mock as its sole oracle. Negative companions cover
pre-cancellation, cancellation during work, collator equality, symlinked lock
ownership, and mixed failure.

Affected-package race tests and the public Doctor command provide the macro
proof. Final QA runs the repository gate, executes Doctor through its public
entry point inside and outside Git, and confirms repository state remains
unchanged.

## Build Order

1. Tidy the authorized Go module metadata with Go tooling.
2. Make external-skill hash ordering total and permutation-independent
   (depends on: 1).
3. Make Repository Skill Set filesystem inspection context-aware and reconcile
   typed ownership and helper naming (depends on: 2).
4. Reconcile Doctor working directories, terminology, remediation, public
   tests, and user guidance (depends on: 3).
5. Synchronize only the authorized canonical and shipped Roundfix Skill pair
   with the settled Doctor contract (depends on: 4).
6. Refresh only the derived TypeScript/Bun Baseline snapshot and catalog/parity
   metadata required by the canonical Roundfix Skill digest (depends on: 5).

## Risks & Considerations

- The external CLI does not define ordering when `localeCompare` returns zero.
  The raw normalized-path tie-breaker is intentionally stricter and
  deterministic; it cannot promise to reproduce a nondeterministic historical
  lock generated from the opposite order of canonically equivalent paths.
- Cancellation is cooperative between filesystem operations. A single
  in-flight operating-system read may finish before the next context check;
  the design adds no goroutine solely to interrupt one file read.
- Adding context to exported functions is a compile-time API change. Current
  repository callers are migrated in the same Task, and no compatibility
  wrapper preserves a blocking contextless API that violates repository Go
  guidance.
- Restoring `os.Getwd()` for profile proof must not leak that value into
  Repository Skill Set inspection. Tests assert both values independently.
- Broad skill synchronization can touch unrelated owned skills. The dedicated
  tooling Task must mutate only the authorized Roundfix Skill pair and its own
  Task file.
- The TypeScript/Bun setup snapshot participates in nested digest authorities.
  The derived refresh must update the complete catalog/parity chain and prove
  it with the existing authorial and compatibility tests; changing only the
  visible `contentDigest` is incomplete.

## Decisions

- Use one context from Doctor through all repository-controlled blocking work.
- Keep the current package boundaries and injected Doctor seam; do not add a
  `HealthChecker` method.
- Define raw normalized path order as the final collation tie-breaker.
- Restore working-directory fallback only for Agent Selection Profile proof.
- Join multiple remediation actions with `&&`, preserving owned-first order.
- Produce module metadata only through `go mod tidy`.
- Preserve all prior commits and every archived Spec byte-for-byte.
- Reconcile the failed QA digest by regenerating the exact derived Baseline
  chain, without changing Baseline runtime behavior or external authorities.
