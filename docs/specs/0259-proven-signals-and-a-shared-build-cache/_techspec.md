---
spec: 0259-proven-signals-and-a-shared-build-cache
prd: _prd.md
created: 2026-10-12
---

# Proven signals and a shared build cache — Technical Spec

## Executive Summary

F6's sender is unknown, so this design removes every way the repository can
send one: each signal to a recorded PID or process group goes through a start
identity proof in `internal/store`, a `go/ast` repository rule forbids every
other signal path, waiting helpers get a lifeline to their test binary, and the
suite guard reports and kills any helper still alive after its package. The
build cache half moves `make` to one Go cache under the temporary directory,
which a measurement shows the Run sandbox can write, and teaches `doctor` and
`gc` the repository-declared caches left in ended worktrees.

The trade-off accepted is the window between the proof and the signal. A
start-identity proof followed by `kill` is not atomic, and macOS offers no
process descriptor to close the gap; psutil, LWN and this repository's Force
Stop accept the same window. With sequential PIDs the window would have to span
a full wrap of the PID counter, so the design closes the realistic hazard,
which is a PID held across minutes, and leaves the theoretical one documented.
The second trade-off is that the cache lives where the operating system may
clean it: a cleaned entry is a cache miss, never a wrong build.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; tests use temporary homes, repositories and fake runners. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0266 governs every decision below.
  ADR-0213: "must end when the test binary that started it ends". ADR-0172:
  "It never deletes, compacts or writes anything", and ADR-0266 amends only
  its byte rule. ADR-0127: "Process residue is a readiness fact, not a Run
  record", unchanged. ADR-0044: "a provably dead owner makes the lock
  orphaned", and the legacy PID-only proof stays. ADR-0126: "The guard
  therefore runs per package". ADR-0125: "A spawned fixture is a compiled
  binary, never a written script". ADR-0259: "a ceiling only falls". ADR-0171:
  "A retention prune reports only what it reclaimed". ADR-0256: "A release
  waits for the skills that describe a changed surface". ADR-0244: "A Spec
  declares the domain terms it introduces". ADR-0257: "Inside a Run the Daemon
  is the only full gate". ADR-0258: "A Spec is authored against the QA rerun
  classes". Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AE grant of 2026-10-11 ("Concedo" for the Governed Paths each Spec declares;
  the Makefile for this Spec only); bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/storage.md`, `.roundfixrc.yml`,
  `Makefile`, `docs/agents/specific-repository.md`,
  `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0259-proven-signals-and-a-shared-build-cache/_authorization.md`.

## Current behavior

Measured on 2026-10-12 at `7a8b7f70` on the maintainer's arm64 Mac.

**Signal senders.** A search for `syscall.Kill`, `unix.Kill`, `.Signal(`,
`.Kill()` and `os.FindProcess` in every Go file found 31 sites. Zero-signal
probes and signals through the `*os.Process` of an unreaped own child are
safe. These are not:

| Site | What it signals | Why it is unproven |
| --- | --- | --- |
| `internal/cli/detach_test.go`, `terminate` | a PID read from a file, through `os.FindProcess` | the survivor may have exited; nothing compares identity |
| `internal/store/process_unix_test.go`, tree helper cleanup | `-pid` with `SIGKILL`, always | runs after the test reaped the root and its group emptied |
| `internal/cli/implement_detach_teardown_test.go`, `killDetachFixtureGroup` | `-pid` with `SIGKILL` | keyed on a child PID, with no leader proof |
| `internal/cli/implement_test.go`, detach caller kill | `-pid` with `SIGKILL` | same |
| `internal/cli/detach_unix.go`, `killDetachedProcessGroup` | `-pid` with `SIGKILL` | a concurrent `Wait` may already have reaped the leader |
| `internal/daemon/qa_format_unix.go`, the `Cancel` hook | `-pid` with `SIGKILL` | the hook can run after the leader was reaped |
| `internal/store/process.go`, `TerminateAndWait` | the owner with `SIGKILL` after the grace | identity proven only before the graceful signal |

Shell fixtures contain only `kill -0` probes. Production Force Stop, reconcile
and Delivery Queue stop already prove the owner before the first signal.

**PIDs.** Six consecutive spawns received 91867 to 91872: macOS allocates in
sequence and wraps below 100,000, and the suite's at least 83,660 Git
processes make about one wrap per full run.

**Helpers that wait.** `TestOwnerProcessHelper` (`graceful`, `ignore`,
`tree-root`, `tree-grandchild`) and `TestCLIForceStopOwnerProcessHelper` wait
for a signal and watch nothing; `ignore` survives `SIGTERM`. The detach
children watch the test binary with `store.ProcessAlive`, a PID-only probe.
The worktree staging helper already ends on standard-input end of file, the
pattern this Spec adopts. Six `store.test` helpers were found orphaned for 45
hours on 2026-10-08.

**F6.** Spec 0254's QA saw the stdout subtest of
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` report the item
binary's `migrate --check` exit as -1 in one of two head suites, and its
Verification once more; the test was then made sequential. One hypothesis was
ruled out: 300 runs of copying a Go binary over an executed copy of itself and
running it, and 40 runs alternating two different binaries, produced no
`SIGKILL` from code signing.

**Cache.** `GOCACHE ?= $(CURDIR)/.gocache` dates from `fe345ab0` (2026-07-29),
whose QA recorded `open /Users/marcio/Library/Caches/go-build/...: operation
not permitted` inside the Run sandbox. `.roundfixrc.yml` keeps
`agent_full_access: false`, and the Codex user configuration selects
`workspace-write` with no extra writable roots. A probe script run under
`codex sandbox -P :workspace` reported: working directory writable, `$TMPDIR`
writable, `/tmp` writable, `~/Library/Caches/go-build` denied, another
`~/Library/Caches` directory denied; the same script outside the sandbox wrote
everywhere. Spec 0037's QA had already built with
`GOCACHE=/private/tmp/roundfix-qa-0037-gocache` inside the sandbox. Today the
user checkout's `.gocache` holds 1.9 GB and the Go default cache 416 MB.
`roundfix gc` touches no worktree, and Doctor's `storage` check never measures
directory bytes.

## System Architecture

- **Proven Signal operations** (`internal/store`, Unix): two exported
  functions beside the existing owner-identity reader, used by production
  detach and QA-format teardown and by the tests. `TerminateAndWait` reuses
  `proveOwner` before its forced kill.
- **Signal rule** (`internal/testfixture`): a repository test beside the
  Parallel Test Package rule; it parses every Go file of the module.
- **Helper Lifeline**: test-only changes to the helper families listed above,
  plus an owner identity carried in the detach children's environment.
- **Helper Residue check** (`internal/suiteguard`): `Main` marks the
  environment before the package's tests and scans the process table after
  them. It reads processes through `golang.org/x/sys/unix`, already required,
  and does not import `internal/store`, whose tests import the guard.
- **Shared build cache** (`Makefile`): the `GOCACHE` default.
- **Worktree Build Caches** (`internal/config`, `internal/cli`): a
  `worktree.caches` key, one discovery used by the Doctor `storage` check and
  the GC Command, and their report lines.

## Implementation Design

### Interfaces

```go
// internal/store/process_signal_unix.go
// sent is false when the process is gone or its identity changed.
func SignalProvenProcess(ctx context.Context, pid int, identity string, signal syscall.Signal) (sent bool, err error)
func SignalProvenProcessGroup(ctx context.Context, pgid int, leaderIdentity string, signal syscall.Signal) (sent bool, err error)

// internal/suiteguard/residue.go
const testBinaryMarkEnv = "SUITEGUARD_TEST_BINARY"
type residueProcess struct{ PID int; Command string }
func helperResidue(mark string, self int) ([]residueProcess, error) // per-platform readers
func killResidue(process residueProcess, mark string) error          // rereads the mark first

// internal/config/config.go
type Worktree struct { /* existing fields */ Caches []string }

// internal/cli/worktree_build_cache.go
type worktreeBuildCache struct{ Worktree, Path string; Bytes int64 }
type worktreeBuildCacheScan struct{ Reclaimable []worktreeBuildCache; Kept int; Bytes int64; Bounded bool }
func scanWorktreeBuildCaches(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, budget time.Duration) (worktreeBuildCacheScan, error)
```

### Data Models

No schema change. `worktree.caches` is a list of strings in Project or User
Config, empty by default. This repository declares `.gocache`.

### API Contracts

1. API Contract: `SignalProvenProcess`. It refuses, with
   `ErrOwnerProcessIdentityUnproven` and no signal, a PID at or below zero, the
   caller's own PID, an empty identity and the zero signal. It reads the live
   start identity with the same reader as `OwnerProcessIdentity`. Absent
   process or different identity: `sent` false, no error. Unreadable identity
   of a live process: no signal and the read error. Otherwise it sends and
   reports `sent`; `ESRCH` from the send is `sent` false, no error.
2. API Contract: `SignalProvenProcessGroup`. When a process with PID `pgid`
   exists, its identity must equal `leaderIdentity`, or nothing is sent. When
   none exists, it sends to `-pgid` only if the platform lists at least one
   member of group `pgid`; darwin reads `kern.proc.pgrp`, Linux reads
   `/proc/<pid>/stat`, other Unix systems return `ErrOwnerProcessUnsupported`.
   `EPERM` from a group whose only members are exiting is `sent` false.
3. API Contract: Force Stop escalation. After the graceful signal and the
   grace, `TerminateAndWait` runs `proveOwner` with the same recorded identity
   before the forced signal. A mismatch is proof of exit: it returns nil and
   sends nothing. An empty recorded identity keeps today's PID-only behavior.
4. API Contract: the Signal rule. `TestEverySignalGoesToAProvenProcess` walks
   every `.go` file under the repository root except `.git`, `.gocache`, `bin`
   and `testdata`, and reports `<file>:<line>: <form>` for: a call to
   `syscall.Kill` or `unix.Kill` whose signal argument is not the literal `0`
   or `syscall.Signal(0)`; a call to `os.FindProcess`; a string literal that
   runs `kill`, `pkill` or `killall` with anything other than `-0`. The
   allowlist names functions: the store's `signalOwnerProcess` on each
   platform, `SignalProvenProcess`, `SignalProvenProcessGroup`, the suite
   guard's `killResidue`, and the self-interrupt in `internal/cli/runui.go`.
   It also runs on seeded trees: `direct-kill`, `bare-pid-lookup`,
   `group-kill` and `shell-kill` each fail naming the seeded line, and
   `zero-probe` and `allowlisted` pass.
5. API Contract: Helper Lifeline. A store or CLI helper started by a test
   receives a pipe as standard input and the test keeps the write end until
   cleanup; the helper ends with status 0 when standard input reaches end of
   file, in addition to its existing behavior. `tree-middle` passes its
   standard input to `tree-grandchild`. A detach child receives
   `ROUNDFIX_DETACH_TEST_OWNER_IDENTITY` and ends when the owner's live
   identity differs or the owner is absent.
6. API Contract: Helper Residue. Before running the tests, `suiteguard.Main`
   sets `SUITEGUARD_TEST_BINARY` to `<pid>.<16 random hex digits>`. After the
   tests it lists live processes of the same user whose environment holds
   exactly that assignment, other than itself, polling for up to 10 s until
   none remains. For each one left it prints
   `suiteguard: helper residue: pid <pid>: <command>`, kills it with
   `killResidue`, which rereads the process environment and sends `SIGKILL`
   only while the mark is still there, and returns exit 1. When the process
   table cannot be read it prints
   `suiteguard: helper residue check unavailable: <reason>` and does not fail.
   Windows and Unix systems other than darwin and Linux print the
   unavailable line with `unsupported on <GOOS>`.
7. API Contract: the shared build cache. With `GOCACHE` neither in the
   environment nor on the command line, the Makefile computes once
   `<root>/roundfix-go-build-<uid>`, where `<root>` is `$TMPDIR` without a
   trailing slash, or `/tmp` when `TMPDIR` is empty, and `<uid>` is `id -u`.
   It uses that path when it can create the directory and create and remove a
   file in it, and `$(CURDIR)/.gocache` otherwise, and exports the result. An
   explicit `GOCACHE` is exported unchanged.
8. API Contract: `worktree.caches`. Each entry is a clean, relative,
   slash-separated path with no `..` segment and not `.git` or under it;
   anything else is a configuration error naming the key and the entry.
   Duplicates collapse.
9. API Contract: ended worktrees. For the loaded repository a worktree is
   ended when it is the Run Worktree recorded on a terminal Run, a kept Task
   Worktree of a terminal Run (a sibling directory named for that Run's
   Tasks), or the worktree of a Delivery Queue item in stage `merged`. A
   worktree of an Active Run or of an item in any other stage is kept and
   counted. A worktree must lie under the configured worktree location, and
   each `<worktree>/<entry>` must be a real directory, not a symbolic link,
   whose resolved path stays inside the worktree; otherwise it is skipped.
10. API Contract: Doctor `storage`. With no declared cache the line is
    byte-identical to today. Otherwise the detail gains, before `; next:`,
    `; Worktree Build Caches reclaimable: <n> (<bytes> bytes)`, with
    `at least <bytes> bytes` when the walk stopped at its two-second budget.
    `<n>` greater than zero makes the status `found` and adds `roundfix gc` to
    the next action. Doctor still writes nothing and never reports `failed`.
11. API Contract: GC Command. With no declared cache the output is
    byte-identical to today. Otherwise, after the Run Retention lines, the
    dry run prints `Worktree Build Caches declared`, `reclaimable`,
    `bytes reclaimable` and `kept` lines and the list
    `Reclaimable Worktree Build Caches`; it writes nothing. The live run
    rereads each worktree's state immediately before removing its cache,
    removes only the declared directory, and prints `removed`,
    `bytes reclaimed` and `kept` lines and the list
    `Removed Worktree Build Caches`. A removal failure is reported on stderr
    and exits 1 after the other removals.

### Surface Transcripts

1. Surface Transcript: `roundfix gc --dry-run` in a repository declaring
   `.gocache`, with one ended Run Worktree holding a cache and one Active Run.

   ```transcript
   $ roundfix gc --dry-run
   stdout:
   GC dry-run
   ...
     Database bytes reclaimable (estimated): <bytes>
     Worktree Build Caches declared: .gocache
     Worktree Build Caches reclaimable: 1
     Worktree Build Cache bytes reclaimable: <bytes>
     Worktree Build Caches kept: 1
     Reclaimable Worktree Build Caches:
       <worktree>/.gocache
   stderr:
   exit: 0
   ```

2. Surface Transcript: the live `roundfix gc` on the same fixture.

   ```transcript
   $ roundfix gc
   stdout:
   GC complete
   ...
     Database bytes after: <bytes>
     Worktree Build Caches declared: .gocache
     Worktree Build Caches removed: 1
     Worktree Build Cache bytes reclaimed: <bytes>
     Worktree Build Caches kept: 1
     Removed Worktree Build Caches:
       <worktree>/.gocache
   stderr:
   exit: 0
   ```

3. Surface Transcript: `roundfix doctor` on the same fixture before `gc`.

   ```transcript
   $ roundfix doctor
   stdout:
   ...
   storage: found (Runs reclaimable: 0; Run Database free bytes: <bytes>; Worktree Build Caches reclaimable: 1 (<bytes> bytes); next: roundfix gc)
   ...
   stderr:
   ...
   exit: 0
   ```

### Invariants

1. No signal other than the zero probe leaves Roundfix code or tests except
   through an unreaped own child's `*os.Process` or an allowlisted Proven
   Signal operation.
2. A Proven Signal sends nothing when the live identity differs from the
   recorded one.
3. Every waiting helper ends within the test bound after its test binary dies.
4. The residue check never signals a process that does not carry the binary's
   mark at the moment of the signal.
5. With no declared Worktree Build Cache, `doctor` and `gc` output is
   byte-identical to today's.
6. `gc` never removes anything outside a declared directory of an ended
   worktree under the worktree location, and `doctor` removes nothing.
7. An explicit `GOCACHE` always wins, and a sandboxed `make` still builds with
   no variable set by the Agent.

## Coverage Map

- Goal 1, Story 2, Core Features 1–2, Success Metric 2 → Proven Signal
  operations, Force Stop re-proof, Signal rule (task_01).
- Goal 2, Story 3, Core Features 3–4, Success Metric 3 → Helper Lifeline,
  Helper Residue check (task_02).
- Goal 3, Story 1, Core Feature 5, Success Metric 1 → F6 test returned to
  parallel (task_02) and the ten-suite run (task_05).
- Goal 4, Story 4, Core Feature 6, Success Metric 4 → shared build cache
  (task_03).
- Goal 5, Story 5, Core Feature 7, Success Metric 5 → `worktree.caches`,
  discovery, Doctor and GC lines (task_04).
- Core Feature 8 → `CONTEXT.md` in task_01, task_02 and task_04; guides and
  skill in task_04.

## Integration Points

None outside the repository. The Codex sandbox is an environment the shared
cache must work inside, measured above; nothing calls it.

## Testing Approach

- `internal/store/process_signal_unix_test.go` (new): a real helper whose
  recorded identity matches is signalled; a mismatched identity sends nothing;
  a group whose leader identity matches is signalled; a group whose leader was
  reaped but whose member lives is signalled; an empty identity and the
  caller's own PID are refused. `TestOwnerProcessControllerReprovesBeforeForceKill`
  uses the controller's injected readers to change the identity during the
  grace and asserts no forced signal.
- `internal/testfixture/signal_rule_test.go` (new): API Contract 4 on the
  repository and on seeded trees in `t.TempDir()`.
- Helper Lifeline death tests: `TestOwnerProcessHelpersEndWithTheirTestBinary`
  in `internal/store` and `TestCLITestHelpersEndWithTheirTestBinary` in
  `internal/cli` start an inner test binary that starts every helper family
  and records each PID and identity, kill the inner binary with `SIGKILL`, and
  poll until every recorded identity is gone, within `testwait`'s bound.
- `internal/suiteguard/suiteguard_test.go`: the existing fixture-binary pattern
  runs `run` with a test function that leaves a marked child alive and one
  whose child ends; the first fails, prints the residue line and leaves the
  child dead, the second passes. A third asserts the mark's form.
- `internal/verifyselect/build_cache_makefile_test.go` (new): runs the
  repository Makefile's `test` target with a stub `GO` that records `GOCACHE`,
  under a writable `TMPDIR`, an unwritable one and an explicit `GOCACHE`.
- `internal/config/worktree_caches_test.go` and
  `internal/cli/worktree_build_cache_test.go` (new): API Contracts 8–11 on
  temporary homes with seeded Runs, Delivery Queue items and worktrees,
  including a symbolic-link escape and an Active Run. One test per Surface
  Transcript asserts every line.
- Each Task that converts or adds tests in a Parallel Test Package runs that
  package once with `-race` (ADR-0259).

## Build Order

1. Proven Signal operations, Force Stop re-proof, production and test senders
   converted, the Signal rule, ADR-0266 and the **Proven Signal** term
   (task_01).
2. Helper Lifeline, the Helper Residue check, the F6 test returned to parallel
   and the terms **Helper Lifeline** and **Helper Residue** (task_02)
   (depends on: 1).
3. The shared build cache in the Makefile and the repository guide (task_03)
   (depends on: 2).
4. Worktree Build Caches in config, Doctor and GC, the guides, the skill and the
   term **Worktree Build Cache** (task_04) (depends on: 2).
5. QA gate with ten full suites (task_05) (depends on: 2, 3, 4).

task_02 edits the same helper and test files task_01 converts, and task_04
shares `CONTEXT.md` with both, so the three run in series. task_03 shares only
the repository guide with task_02, so it follows task_02 and runs beside
task_04.

## Risks & Considerations

- **F6 may have another cause.** Memory pressure while the disk was full could
  kill a child with `SIGKILL` too. The shared cache removes the duplicated
  gigabytes that filled the disk, and the QA gate's ten suites decide; a
  recurrence fails the gate and its report records the signal.
- **The residue check may find real leaks.** A package that leaves a process
  today will fail once task_02 lands; task_02 must fix the leak it finds, not
  weaken the check, and reports each one.
- **Environment assertions.** A test that compares a child's whole environment
  will now see `SUITEGUARD_TEST_BINARY`. The mark has no `ROUNDFIX_` prefix so
  that production filters of Roundfix variables neither strip nor report it.
- **Sandboxed residue check.** The sandbox may deny reading other processes'
  arguments; the check then reports itself unavailable. The Daemon's
  Verification runs outside the sandbox, so the gate still exercises it.
- **Temporary-directory cleaning.** macOS may delete cache files nobody read
  for three days; Go treats a missing entry as a miss.
- **Doctor time.** The walk is capped at two seconds and only in repositories
  that declare caches.

## Glossary

- adds: **Proven Signal**
- adds: **Helper Lifeline**
- adds: **Helper Residue**
- adds: **Worktree Build Cache**
- changes: **GC Command**
- not a term: **Signal rule** — the name of one repository rule test
- not a term: **Shared build cache** — a Makefile default of this repository
- not a term: **Worktree Build Caches** — plural of the declared term
- not a term: **Proven Signal operations** — the store functions that send a Proven Signal
- not a term: **Helper Residue check** — the suite guard step that finds Helper Residue
- not a term: **Current Behavior** — a section heading
- not a term: **Ended Worktrees** — the API Contract's description of eligibility

## Decisions

- Proven Signals compare start identities, group signals prove the leader or
  live membership, and a `go/ast` rule forbids every other path; see ADR-0266.
- Helpers end on a lifeline pipe or an identity-proven owner watch, and the
  suite guard kills marked residue; see ADR-0266.
- `make` defaults to `roundfix-go-build-<uid>` under the temporary directory
  with the worktree cache as fallback; see ADR-0266.
- `worktree.caches` declares Worktree Build Caches; Doctor measures them under
  a two-second budget and GC removes them from ended worktrees; see ADR-0266.
- Ten consecutive clean full suites accept F6; see ADR-0266.
