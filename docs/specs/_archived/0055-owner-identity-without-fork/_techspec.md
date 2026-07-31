---
spec: 0055-owner-identity-without-fork
prd: _prd.md
created: 2026-07-30
---

# Owner identity without fork — Technical Spec

## Executive Summary

One function changes shape: `processStartIdentity` in
`internal/store/process_unix.go` shells out to `ps -o lstart=` on every read.
It becomes two build-tagged implementations that read the kernel directly —
procfs on Linux, `sysctl` on macOS — with no subprocess and no cgo.
`golang.org/x/sys` is already a direct dependency, so the macOS path needs no
new module and the Linux path needs only `os.ReadFile`.

The primary trade-off is the token's stability across platforms. The current
token is a `ps`-formatted wall-clock string; the replacement is a
platform-native start time. Both are opaque and equality-compared, so the
change is safe for new Runs but **not** comparable with tokens already recorded
by the `ps` implementation. We treat a token the current platform cannot have
produced as *unreadable*, not as a mismatch, which is exactly the distinction
this Spec introduces — so the migration falls out of the feature rather than
needing a schema step.

The second trade-off is where the unreadable/mismatch split lives. Both
conditions today collapse into `ErrOwnerProcessIdentityUnproven` at
`internal/store/process.go:99` and `:102`. We add one sibling sentinel rather
than an error-kind enum, because every caller already branches with
`errors.Is` and only the diagnostic text differs.

## Project Constraints

- Identifier strategy: not applicable — the owner identity stays an opaque
  equality-compared token; no project-owned Internal Identifier is created.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local process inspection and CLI
  argument parsing only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0044 requires orphaned Run locks to
  be reclaimed only on proven owner death, which the kernel read must keep
  proving; ADR-0052 protects compare-and-set terminal completion; Spec 0037's
  refuse-on-proven-mismatch semantics are preserved. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the PRD authorizes exactly
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, and the
  Skill-digest fallout in exactly
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.
- Dependency discipline: applicable — `golang.org/x/sys` is already required at
  `v0.45.0`; no new module is introduced and the stdlib covers the Linux path.
  Source: `docs/agents/specific-repository.md`.

## System Architecture

- `internal/store/process_unix.go` keeps `processAbsent` and
  `signalOwnerProcess` and loses `processStartIdentity`.
- `internal/store/process_linux.go` (new, `//go:build linux`) reads
  `/proc/<pid>/stat`.
- `internal/store/process_darwin.go` (new, `//go:build darwin`) reads
  `KERN_PROC_PID` through `unix.SysctlRaw`.
- `internal/store/process_unix_other.go` (new,
  `//go:build unix && !linux && !darwin`) returns the unreadable sentinel, so
  an unsupported Unix degrades exactly like the documented non-Unix stub
  instead of failing to build.
- `internal/store/process.go` gains the `ErrOwnerIdentityUnreadable` sentinel
  and the classification at the two existing branches.
- `internal/store/store.go` records the capture-failure marker; the Run
  inspection surface renders it.
- `internal/cli/cli.go` fixes Stop Command argument ordering, reusing whatever
  Spec 0042 introduced for Attach.
- `internal/cli/implement.go` emits the one startup warning.

No new packages. `process_windows.go` and `process_other.go` are untouched.

## Implementation Design

### Interfaces

```go
// internal/store/process.go
var (
    // ErrOwnerProcessIdentityUnproven — unchanged; a proven different owner.
    ErrOwnerIdentityUnreadable = errors.New("owner process identity is unreadable")
)
```

Linux, stdlib only. Field 22 of `/proc/<pid>/stat` is the start time in clock
ticks since boot. The field is parsed from the last `)` forward, because the
comm field at index 2 may itself contain spaces and parentheses:

```go
func processStartIdentity(_ context.Context, pid int) (string, error) {
    // returns "linux:<starttime-ticks>"; ENOENT means the process is gone,
    // every other error is unreadable.
}
```

macOS, through the existing `golang.org/x/sys/unix`:

```go
func processStartIdentity(_ context.Context, pid int) (string, error) {
    // unix.SysctlRaw("kern.proc.pid", pid) -> unix.KinfoProc
    // returns "darwin:<p_starttime.Sec>.<p_starttime.Usec>";
    // ESRCH means the process is gone, every other error is unreadable.
}
```

The platform prefix is what makes the migration self-describing: a token
recorded by the `ps` implementation carries no prefix, so the comparison sees a
token this platform cannot produce.

### Classification at the proof

`proveOwner` keeps its shape and structure; only the two failure branches
change:

```go
// read failed and the process is still present  -> ErrOwnerIdentityUnreadable
//   wrapping the host error, with the next action.
// recorded token is not comparable on this platform (missing prefix, or a
//   prefix from another platform)                -> ErrOwnerIdentityUnreadable
// live token differs from a comparable recorded token
//                                               -> ErrOwnerProcessIdentityUnproven
```

The existing exited-between-checks recovery at `process.go:96` is preserved
ahead of both, so proven absence still wins over any read failure.

### Supervised path for an unreadable identity

Force Stop keeps failing closed. The supervised exit is an explicit flag on the
Stop Command — `--owner-identity-unreadable` — that is accepted **only** when
the proof returned `ErrOwnerIdentityUnreadable`. Passing it when the identity
is readable, or when the proof returned a mismatch, exits `2` and signals
nothing. This keeps "explicit, never the default" literal: there is no
configuration key, no environment variable, and no timeout that reaches it.

### Data Models

- One additive column on the Run row: `owner_identity_unproven` (integer, `0`
  default), set at creation when capture failed. A legacy row predating the
  identity column keeps NULL identity with the marker unset, preserving
  ADR-0044's PID-only degradation exactly.
- No change to Stop Request semantics, signaling order, or the Agent Session
  registry.

### API Contracts

- `roundfix stop <run-id> --force` and `roundfix stop --force <run-id>` parse
  identically, matching the Attach Command.
- A Force Stop refusal names its condition: proven mismatch, or unreadable
  identity with the host error and the documented next action.
- Run inspection output renders the marker for a Run created without reuse
  protection.

## Coverage Map

- Goal 1 / Story 1 → the two kernel implementations plus the no-subprocess
  proof (Feature 1).
- Goal 2 / Stories 2, 4 → the sentinel split, its diagnostics, and the
  supervised flag (Features 2, 4).
- Goal 3 / Story 3 → the startup warning and the durable marker (Feature 3).
- Goal 4 / Story 5 → Stop Command argument ordering (Feature 5).
- Feature 6 → the preserved-behavior cases in the existing proof table.

## Integration Points

- Kernel only: procfs reads and `sysctl`. No network, no external process, no
  new file written outside the Run Database.

## Testing Approach

`internal/store/process_unix_test.go` holds the existing proof table and gains:
a case asserting the identity read spawns no subprocess (assert via a
`PATH`-shadowing `ps` that fails the test if executed, so the assertion holds
without depending on process accounting); the unreadable-versus-mismatch split
with a synthetic recorded token; a non-comparable legacy token classifying
unreadable; and stability of the token across repeated reads of a live process.
`internal/store/store_test.go` covers the marker at creation and legacy NULL
degradation. `internal/cli/cli_test.go` covers both Stop argument orders, the
supervised flag's acceptance only under an unreadable identity, and its
refusal otherwise. `internal/cli/orphan_unix_test.go` keeps ADR-0044 reclaim
behavior green.

## Build Order

1. Platform identity reads: `process_linux.go`, `process_darwin.go`,
   `process_unix_other.go`, `processStartIdentity` removed from
   `process_unix.go`, with the no-subprocess and token-stability tests.
2. Classification split: `ErrOwnerIdentityUnreadable`, the non-comparable-token
   rule, and the two diagnostics (depends on: 1).
3. Creation marker and startup warning, including legacy NULL degradation
   (no dependency on 1–2).
4. Stop Command argument ordering and the supervised
   `--owner-identity-unreadable` flag (depends on: 2 for the condition it
   keys on).
5. Skill pair and user-guide documentation for the diagnostics, the marker, and
   the supervised path, plus the authorized derived digest fallout
   (depends on: 1–4).

## Risks & Considerations

- **Token format change is a silent no-op risk.** If the platform prefix were
  omitted, a legacy `ps` token would compare unequal to a live kernel token and
  Force Stop would report a *mismatch* on a correctly-owned Run — failing closed
  for the wrong reason. The prefix is what makes that case classify unreadable;
  a test pins it directly rather than through the happy path.
- **procfs field parsing.** The comm field can contain spaces and
  parentheses, so field 22 must be counted from the last `)`. Parsing from the
  left is the classic defect here.
- **`unix.KinfoProc` is a platform struct.** It is only referenced under
  `//go:build darwin`, so a Linux or Windows build never resolves it.
- **The supervised flag is a real hazard.** It signals a process whose identity
  could not be proven. It is gated on the unreadable classification, refuses
  under mismatch, and is documented as an operator action of last resort — but
  it is the one place this Spec widens what Force Stop may do.

## Decisions

- The token gains a platform prefix and stays opaque; equality is still the
  only supported comparison, and the prefix is never parsed for meaning beyond
  "this platform could have produced it".
- One sibling sentinel rather than an error-kind enum, because callers already
  branch with `errors.Is`.
- An unsupported Unix gets its own stub file so the kernel reads never need a
  portable fallback that would reintroduce a subprocess.
