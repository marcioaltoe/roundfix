---
spec: 0050-doctor-skill-readiness-hardening
prd: _prd.md
created: 2026-07-26
---

# Doctor Skill Readiness Hardening — Technical Spec

## Executive Summary

Harden Repository Skill Set readiness at the two boundaries that review proved
unsafe: filesystem anchoring and external hash ordering. One new internal hash
package will own slash-normalized path-plus-content hashing and use the
American English collation tables from `golang.org/x/text`, matching the
installed skills CLI's resolved `String.localeCompare` behavior for
punctuation, case, and Unicode. The repository checker will perform every lock
and skill read through an `os.Root` anchored at the supplied Git root after
non-following metadata validation. The trade-off is one official Go
supplementary dependency and its generated Unicode tables in exchange for
removing an incomplete handwritten collation approximation.

## Project Constraints

- Identifier strategy: not applicable — existing skill names and
  repository-relative paths remain the only identifiers. This Spec adds no
  project-owned identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the implementation is offline and
  filesystem-only and changes no authentication, authorization, transport, or
  HTTP contract. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 require Doctor's
  profile proof to remain independent and eager; ADR-0066 and ADR-0072 require
  Baseline restoration to remain a Go CLI operation and preserve the existing
  Go `sync-setups` destination. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization:
  authorization given on 2026-07-26 to add `golang.org/x/text` for Unicode
  collation; bounded protected files: `go.mod`, `go.sum`. No other protected
  tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- `internal/skillhash` becomes the single internal compatibility component for
  external skill digests. It accepts already collected relative paths and
  bytes, orders them with one collator, and emits the existing lowercase
  SHA-256 value.
- `skills` remains the Repository Skill Set ownership boundary. Its filesystem
  walker collects regular files and rejects special entries; its repository
  checker anchors all repository-relative operations and classifies missing,
  outdated, or unreadable state.
- `internal/baseline` keeps ownership of Baseline skill restoration and lock
  rendering. It converts restoration files to the shared hash input rather
  than duplicating collation.
- `internal/cli` keeps Doctor output, check ordering, exit status, and
  remediation ownership. It passes only the resolved Git root to Repository
  Skill Set readiness.
- `CONTEXT.md` remains the canonical vocabulary contract. The archived Spec
  0036 and existing Baseline asset synchronization code remain unchanged.

The new package is deliberately narrow. It does not walk filesystems, decode
locks, know about Doctor output, or expose Baseline types. This keeps
filesystem security in `skills`, presentation in `internal/cli`, and only the
shared compatibility rule in one place.

## Implementation Design

### Shared external-skill hash

The internal component exposes one small value type and one pure operation:

```go
type File struct {
    Path    string
    Content []byte
}

func Sum(files []File) string
```

`Sum` copies or indexes the supplied slice without changing caller-owned order,
sorts stably by `collate.New(language.AmericanEnglish).CompareString`, and
feeds each slash-normalized relative path immediately followed by its bytes
into SHA-256. It inserts no separator or length prefix because the installed
skills CLI does not. Callers remain responsible for rejecting unsafe paths,
links, excluded directories, and unreadable files.

`skills.SkillFolderHash` retains its public signature and filesystem failure
behavior. Its walker collects `internal/skillhash.File` values and delegates
only ordering and digest production. Baseline's
`externalSkillsLockDigest` performs the same adaptation from `restoreFile`.
The copied lowercase-byte comparators are deleted.

The compatibility test uses a pinned corpus containing `_a`, `-a`, digits,
case variants, precomposed accented characters, and punctuation at multiple
path depths. Its expected digest is produced once with the locally installed
skills CLI 1.5.19 algorithm and then treated as a product compatibility
contract. Existing real-lock integration tests remain the macro proof.

### Root-anchored repository reads

`checkRepository` validates required names before opening the supplied root,
then calls `os.OpenRoot`. Every repository-relative authority is inspected
through that handle:

1. `skills-lock.json` is checked with `Root.Lstat`; a symbolic link or
   non-regular file is rejected before `Root.ReadFile`.
2. `.agents` and `.agents/skills` are checked component by component with
   `Root.Lstat`. A missing component produces the existing missing-skill
   classifications; a symbolic link or non-directory produces a typed
   readiness error.
3. Each required skill root is inspected without following its final path.
   Missing roots remain normal readiness results. Present roots must be real
   directories before an anchored sub-root or filesystem view is walked.
4. Walkers reject symbolic links and special files and skip `.git` and
   `node_modules` exactly as before.

The `os.Root` handle confines all later operations to the repository even if a
path changes between metadata inspection and use. Component `Lstat` checks
enforce the static no-symlink contract; root confinement prevents an ancestor
replacement race from redirecting reads outside the supplied Git root.
Errors retain the failed operation and repository-relative authority while
wrapping their underlying cause for `errors.Is` and `errors.As`.

### Doctor coordination and remediation

The Doctor coordinator removes the `os.Getwd()` fallback. A non-empty
`roundconfig.Loaded.GitRoot` remains the work directory supplied to profile
proof and the root passed to `checkSkills`. When the loaded Git root is empty,
Doctor does not call the repository checker and instead emits a failed
`skills:` result whose next action tells the caller to run Doctor from a Git
repository.

`doctorSkillReadinessResult` retains owned and external remediation selection
for typed errors. If an injected or future checker returns an error without
ownership and no classified readiness fields, the conservative fallback
prints both existing remediation commands in owned-then-external order. The
output boundary stays unchanged:

```text
skills: failed (repository skill check failed; next: roundfix skills install --target project; bunx skills experimental_install && bunx skills update -p -y)
```

All independent health checks still execute and print in the current order.
Missing Git root or skill failure remains exit `1`; argument failures remain
exit `2`.

### Test ownership and no-mutation proof

Tests whose primary subject is `CheckRepository` and their repository fixture
helpers move from `skills/skills_test.go` to
`skills/repository_test.go`. `SkillFolderHash` tests remain beside
`skills.go`. The move changes no assertions merely to make a failure pass.

Focused negative cases cover symlinked `.agents`, symlinked
`.agents/skills`, and symlinked `skills-lock.json`, with targets containing
otherwise valid ready state. Each case must fail without classifying the
repository ready.

`internal/cli/doctor_test.go` adds public-runner coverage with fake machine
health and profile proof but the real `skills.CheckRepository`. The test
creates a complete temporary Repository Skill Set, snapshots repository files,
the configured User Config path, and `.roundfix`, runs
`Run([]string{"doctor"}, ...)`, then proves every snapshot is unchanged. A
separate case proves an empty loaded Git root does not invoke `checkSkills` and
prints repository-specific remediation.

### Contract reconciliation

`CONTEXT.md` restores that the Doctor Command reports the detected acpx version
against the minimum. This is a vocabulary correction for behavior that never
stopped shipping.

The ownership-safe `sync-setups` behavior removed from Spec 0036's text is not
rebuilt here. Its Go implementation already preserves repository-owned
digests in `internal/baseline/assets_sync.go`, its regression test already
proves the behavior, and ADR-0072 assigns that responsibility to the Go
Baseline operation. This Spec records that evidence while keeping the archived
Spec byte-identical.

## Coverage Map

- Goal: root-anchored reads and User Story 1 → anchored repository checker,
  ancestor and lock `Lstat` cases.
- Goals: exact ordering and shared implementation; User Stories 2–3 →
  `internal/skillhash`, punctuation/Unicode oracle, Doctor and Baseline
  adapters.
- Goal: resolved Git root and deterministic remediation; User Story 4 →
  Doctor coordinator and exact-output command tests.
- Goal: executable no-mutation evidence; User Story 5 →
  `repository_test.go` and real-checker public Doctor tests.
- Goal: canonical terminology → `CONTEXT.md` reconciliation plus existing
  Baseline asset synchronization evidence.

## Integration Points

- Installed skills CLI 1.5.19 and 1.5.20 are compatibility references only;
  production code never invokes them.
- `golang.org/x/text/collate` and `golang.org/x/text/language` provide pinned,
  in-process Unicode collation tables.
- `os.Root` is the standard-library repository confinement boundary.
- Baseline restoration consumes the shared hash package without changing its
  lock JSON or transaction contracts.

## Testing Approach

- Unit: a pure `internal/skillhash` oracle test compares the pinned expected
  digest for punctuation, case, digits, and Unicode; a negative companion
  changes one path or byte and requires a different digest.
- Package integration: `skills/repository_test.go` builds complete disposable
  repositories and proves normal readiness, missing/outdated classification,
  malformed authorities, ancestor and lock symlink rejection, and no
  mutation.
- Baseline integration: existing restoration tests plus a focused shared-hash
  assertion prove generated `computedHash` values still use the compatibility
  contract.
- CLI integration: public `Run` tests prove exact output, missing-Git-root
  handling, independent check execution, conservative error remediation, and
  real-checker no mutation.
- Repository integration: the existing real Repository Skill Set contract,
  affected-package race suite, and `rtk make verify` remain mandatory.

No test introduces production-only hooks. Machine health and profile proof use
the pre-existing Doctor dependency seam; the repository checker itself remains
real in the no-mutation test.

## Build Order

1. Add the authorized `golang.org/x/text` module requirement in the bounded
   `go.mod` and `go.sum` tooling slice.
2. Add the shared external-skill hash component, its CLI-oracle corpus, and
   migrate both existing consumers (depends on: 1).
3. Anchor Repository Skill Set reads, reject skill-tree and lock symlinks, and
   relocate repository tests to their owning file (depends on: 2).
4. Remove Doctor's working-directory fallback, complete remediation and
   public no-mutation coverage, and reconcile `CONTEXT.md` (depends on: 3).

## Risks & Considerations

- The upstream skills CLI calls optionless `String.localeCompare`, so its
  default locale is environment-dependent. Roundfix pins American English
  because supported Node and Bun currently resolve `en-US`; the oracle makes
  that assumption visible and fails loudly if the compatibility contract
  changes.
- `os.Root` confines symbolic-link traversal to the root but still permits
  in-root links by design. Explicit component `Lstat` validation is therefore
  required in addition to anchoring.
- Unicode collation tables increase module and binary size. The dependency is
  official Go supplementary code, already present in the local module cache,
  and replaces a handwritten algorithm that cannot satisfy the contract.
- Moving tests can create noisy diffs. The move must preserve existing
  assertions and separate only repository-owned helpers; behavior changes
  belong in newly named cases.
- A lock or shared ancestor security error cannot safely identify one damaged
  ownership group. Conservative unclassified remediation may print both
  commands, but it must never run them.

## Decisions

- Pin `golang.org/x/text` through Go tooling; never hand-edit module checksums.
- Pin `language.AmericanEnglish` rather than inherit the host locale.
- Share only the pure hash rule; keep file collection and policy with each
  owning package.
- Combine non-following component validation with `os.Root` confinement.
- Preserve Doctor's output delimiter, line ordering, read-only contract, and
  exit codes.
- Preserve branch history and every archived Spec artifact.
