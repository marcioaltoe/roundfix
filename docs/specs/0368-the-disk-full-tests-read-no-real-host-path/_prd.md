---
spec: 0368-the-disk-full-tests-read-no-real-host-path
status: active
created: 2026-11-09
surfaces: [cli]
type: fix
---

# The disk-full tests read no real host path

Two `internal/cli` tests of the Disk-Full Park write the literal path `/tmp/x/_pkg_.a` into their fake build output, and the code under test measures free space at that path's nearest existing directory before it calls the stubbed capacity reader. On 2026-11-09 an empty file `/tmp/x` on the maintainer's MacBook made `os.Stat` return `ENOTDIR`, the ancestor lookup failed, the reported free space fell to `0 MiB`, and Spec 0365's first Run parked `run-unresolved` for a defect it did not touch. After this Spec every disk-full test that reaches the filesystem uses a path under its own temporary directory, so no file on the host can change its outcome, and the nearest-directory lookup walks past a file that stands where a directory is expected, as the `doctor` guide already says it does, so a real disk-full park in the same situation names the free space instead of `0 MiB`.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-09-the-disk-full-tests-read-a-real-host-path.md`](references/2026-11-09-the-disk-full-tests-read-a-real-host-path.md),
   the open Backlog Entry recorded by the operator on 2026-11-09, `type: fix`;
   this Spec adopts it.

At the authoring base `docs/backlog/` holds no other open entry and
`docs/findings/` no unresolved Finding, so no other source shares this
context.

## Prerequisites

None in `requires`. Authored against `cfc08683` (v0.137.2 plus the 0366 and
0367 plans). The queue delivers Specs 0366 and 0367 first, and Specs 0369
and 0370 follow; none of them changes `internal/cli/doctor_disk_space.go`,
`internal/cli/doctor_disk_space_test.go` or
`internal/cli/deliver_disk_full_test.go`, the only files this Spec's
implementation Task changes. This Spec does not touch `CONTEXT.md`, any
guide, skill, ADR or record.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change reads the local
  filesystem only and makes no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0370 defines the Disk-Full Park
  and the Free-Space Floor, ADR-0370: "`<N>` is its free space at the park",
  and this Spec makes that value hold when a file stands in the failing path;
  the decision itself is unchanged. ADR-0327 and ADR-0334 carry ADR-0370's
  amendment notices on the Shared Build Cache trim; not applicable, because
  this Spec does not change the trim. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: not applicable — no protected tooling configuration,
  script, ignore file, plugin declaration or version pin is created or
  changed, and no Governed Path is touched (`paths: []`, measured with
  `speccheck.GovernedPath`). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0368-the-disk-full-tests-read-no-real-host-path/_authorization.md`.

## Goals

- No disk-full test reads a path outside its own temporary directory, so a
  file left anywhere on the host cannot change its result.
- The nearest-directory lookup of free space treats a file where a directory
  is expected like a missing path and walks on to the nearest existing
  directory, while any other read error still reports as today.

## User Stories

1. As the operator, I want the disk-full tests to depend on nothing outside
   their temporary directories, so that a stray file on the machine never
   parks a Spec that did not touch the code.
2. As a user whose item parks `disk-full`, I want the park to name the free
   space of the filesystem that holds the failing path even when a file
   stands in that path, so that the reason is not `0 MiB free`.

## Core Features

1. Every test in `internal/cli/deliver_disk_full_test.go` that hands a path
   to the disk-full code builds that path under its own `t.TempDir()`; the
   parsing-only test keeps literal paths, none of them under `/tmp/`.
2. `doctorDiskSpaceAncestor` walks past `ENOTDIR` as it walks past a missing
   path; the Disk-Full Park, the queue's low-disk warning and `doctor`'s
   `disk-space:` line, which share it, measure the nearest existing
   directory.
3. Two new tests prove it with a file created inside a temporary root, never
   in a shared system location: the disk-full path below a file is measured
   at the root, and the lookup still fails on an unreadable directory.

## User Experience

Before this Spec, an item build that fails with
`write <root>/x/_pkg_.a: no space left on device` while `<root>/x` is a file
parks as `disk-full: <root>/x/_pkg_.a: 0 MiB free`; after it, the park reads
`disk-full: <root>/x/_pkg_.a: <N> MiB free`, `<N>` read at `<root>`. A
`doctor` run whose `worktree.location` lies below a file reports the free
space of the nearest existing directory instead of a `DR-DISK-UNREAD`
warning, as `docs/user-guide/commands/doctor.md` already describes
("read at its nearest existing directory"). No text changes.

## Non-Goals / Out of Scope

- Any change to the Disk-Full Park's text, Park Class, next action or retry
  counting, the Free-Space Floor, or `doctor`'s finding codes.
- Following symbolic-link loops, permission errors or other read errors past
  the failing component; they keep failing the lookup.
- Tests outside the disk-free-space code that name `/tmp/` paths only as
  strings, such as `internal/delivery/disk_full_park_test.go`, which never
  touches the filesystem.
- Any change to the Makefile, the workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the Baseline, any guide or skill.

## Success Metrics

1. Success Metric: `internal/cli/deliver_disk_full_test.go` names no `/tmp/`
   path, and its four existing tests pass unchanged in what they assert.
2. Success Metric: with a file created at `<root>/x` inside a temporary root,
   `diskFullFailure` reports `<root>/x/_pkg_.a` with the stubbed free space
   read at `<root>`; on the authoring base it reports 0 bytes and reads
   nothing.
3. Success Metric: `doctorDiskSpaceAncestor` returns the temporary root for
   a path below a file and still returns an error that is not
   `os.ErrNotExist` for a path inside an unreadable directory.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry, recorded by the operator from Spec 0365's Run
  `run_20261109T120541Z_75e62679b1abfda6` before this Spec existed: the
  empty `/tmp/x` "made the reported free space `0 MiB` instead of the
  stubbed value", and removing it "made both tests pass 3 of 3 on main".
  The QA reads the entry, never the Run's diagnostics under `~/.roundfix`.
- The POSIX `stat()` errors (The Open Group Base Specifications, read
  through Exa on 2026-11-09, mirrored at
  https://linux.die.net/man/3/stat): "ENOTDIR: A component of the path
  prefix is not a directory." and "ENOENT: A component of path does not name
  an existing file". `ENOTDIR` therefore means the path cannot exist as
  written, the same conclusion as `ENOENT` for a free-space lookup.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-09 against these sources.

- Reading at `cfc08683`: `diskFullFailure` (`internal/cli/deliver_disk_full.go`)
  calls `doctorDiskSpaceAncestor` before the stubbed `availableBytes`;
  `doctorDiskSpaceAncestor` (`internal/cli/doctor_disk_space.go`) walks to
  the parent only when `os.Stat` reports `os.ErrNotExist` or a non-directory
  that exists; `inspectFreeSpace` (doctor and the queue's low-disk warning in
  `deliver_free_space.go`) uses the same function.
- A prototype of the design in a scratch worktree of `cfc08683`: the two new
  tests fail on the authoring base (`stat ...: not a directory`, 0 bytes) and
  pass with the change; the four existing disk-full tests and every test
  matching `DiskFull|DiskSpace|LowOnDisk` in `./internal/cli` pass; walking
  past every error makes the unreadable-directory case fail;
  `GOOS=windows go vet ./internal/cli` passes. No file was created in a
  shared system location.
- Published sources (Exa): the POSIX `stat()` error list quoted above.
- Secondbrain, read only:
  `qmd query "disk-full test host path tmp ENOTDIR roundfix" --all --files --min-score 0.3`
  returned this repository's mirrored ADR-0370, the Spec 0360 Archive Record,
  the `doctor` guide and a digest on symlinked worktree roots; no page
  decides otherwise.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BA batch 2 brief for Spec 0368 of 2026-11-09, the adopted Backlog Entry, and the maintainer's answer below.

The maintainer's answer (AskUserQuestion via `/roundfix --follow-up`,
2026-11-09 ~11:20 BRT), verbatim:

- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto
  fica para amanhã."

Taken from those and the evidence, with no further interview:

- The lookup walks past `ENOTDIR`. A file where a directory is expected
  means the path cannot exist as written, exactly as a missing component
  does; the guide already promises the nearest existing directory, and the
  Disk-Full Park otherwise reports `0 MiB free`. Every other `os.Stat` error
  still fails the lookup, so a read that cannot be proven stays
  `DR-DISK-UNREAD`.
- The tests build their paths under `t.TempDir()`; the proof creates the
  file inside a temporary root and never in `/tmp` or `$HOME`.
- No ADR: the change restores the behavior ADR-0370 and the `doctor` guide
  already state.
- One implementation Task and the QA gate. The release is the operator's
  patch right after the merge; no Task releases.

## Open Questions

None.
