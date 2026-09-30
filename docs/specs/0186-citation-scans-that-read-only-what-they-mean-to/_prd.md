---
spec: 0186-citation-scans-that-read-only-what-they-mean-to
status: active
created: 2026-09-30
surfaces: [backend]
---

# Citation scans that read only what they mean to

Two review findings from Wave 5 were dismissed at the review ceiling and kept
as Backlog Entries. Both are about a citation scan reading text it should not:

- **The citation projection strips sections from any file named like a Task.**
  Spec 0181's task_06 made the Spec citation walk behind `SC-ADR-UNLISTED` skip
  the Agent-owned `## Result` and the Daemon-owned `## Recorded paths` and
  `## Carry-forward provenance` sections (ADR-0176). It skips them in every
  file whose basename starts with `task_`, including an adopted source under
  `references/`. `references/` is authored and must be read in full.
- **The Relocation Citation scan has a parent-directory race.** Spec 0184's
  scan checks each tracked path's parent directories with `os.Lstat`, caches
  the result and opens the file by its joined path, with no-follow on the last
  component only. A directory swapped for a symbolic link after the check
  redirects the open, possibly outside the repository, and the later
  `os.SameFile` check still passes because both sides resolve through the
  link.

This Spec fixes both without changing any command output, flag, exit code or
schema.

## Project Constraints

- Identifier strategy: not applicable — no identifier, key or event payload is
  created or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0176 defines the authored projection
  the citation walk reads, and this Spec narrows it to the Spec folder's own
  top-level Task files, which is what ADR-0176 names. ADR-0093 checks Spec
  consistency by citation, and ADR-0094 skips a detector whose artifact is
  absent; both hold. ADR-0173 defines the Relocation Citation scan and its
  no-follow reads. This Spec adds ADR-0177, which makes the scan read through a
  repository root. ADR-0104 accepts on evidence a Spec did not author,
  ADR-0130 keeps a path governed once bounded, ADR-0155 makes the `qa` Task
  declare the matrix and ADR-0156 makes a declared promise name a consuming
  Task. This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097
  and ADR-0117. All hold.
  ADR-0168 gives `SC-ADR-RELATED` a commit-ancestry horizon, and it holds
  unchanged beside the narrowed projection. ADR-0166 (recorded paths) and
  ADR-0167 (the pre-PR Pull Request row) govern the Daemon's commit record and
  QA eligibility, which this Spec does not touch, so they do not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty. `internal/baseline/repository_test.go`,
  `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`,
  `docs/references/coverage-record.json`, `go.mod` and the Makefile stay
  untouched. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- Only the Spec folder's own Task files lose their Agent- and Daemon-owned
  sections to the citation walk; every other Markdown file under the Spec is
  read in full.
- A Relocation Citation scan never reads a file outside the repository, even
  when a tracked directory is replaced by a symbolic link during the scan.
- Neither change moves any command output, flag, exit code, schema or existing
  test.

## Core Features

1. **The projection applies to top-level Task files only.** `readSpecCitations`
   skips the three non-authorial sections only in files directly inside the
   Spec folder whose name matches `task_*.md`. A file under `references/`, or
   in any other subdirectory, is read in full. The `qa/` directory stays
   skipped.
2. **The scan reads through a repository root.** The Relocation Citation scan
   opens every tracked file through one `os.Root` opened on the repository root
   (ADR-0177). The per-component `os.Lstat`, the `os.SameFile` comparison and
   the non-blocking open on Unix stay, so a symbolic link inside the repository
   is still never followed on purpose.

## Non-Goals / Out of Scope

- Any change to what `SC-ADR-UNLISTED`, `SC-ADR-RELATED` or
  `SC-CITATION-UNSUPPORTED` report for authored text.
- Any change to which files the Relocation Citation scan reads, what it
  reports, its caps or its summaries.
- Defending against bind mounts or other constructs that need root privileges
  to create, which `os.Root` does not cover by design.
- A new dependency, or an edit to `go.mod`.

## Success Metrics

1. A Spec whose only unlisted ADR citation sits in the `## Result` of a file
   `references/task_example.md` reports `SC-ADR-UNLISTED` for it. The same
   citation in a top-level `task_01.md` `## Result` reports nothing.
2. In a disposable repository, a scan whose tracked directory is replaced by a
   symbolic link to a directory outside the repository, between the parent
   check and the open, reads nothing outside the repository. The file is
   skipped, and no citation from the outside file is reported.
3. Every existing speccheck and Baseline test passes unchanged, and the Windows
   build of the non-test code succeeds.

## Recorded limits

- `os.Root` follows a symbolic link that stays inside the repository. A link
  planted inside the repository mid-scan can still redirect a read to another
  repository file until the post-open checks see it. That read never leaves the
  repository, and its content only produces warnings.
- The projection keys on the file's location and name, not its content. A
  top-level file named `task_*.md` that is not a Task file of the graph would
  still be projected; the Task Graph rules forbid such a file.

## Decisions

- **Scope the projection by location.** ADR-0176 names "the Task files", and
  only the Spec folder's top-level `task_*.md` files are Task files. Parsing the
  Task Graph to find them would add a failure mode to a citation check for no
  gain.
- **Use `os.Root`, not a hand-written `openat` walk.** The standard library
  already implements traversal-resistant opening on every supported platform:
  `openat` with `O_NOFOLLOW` per component on Unix, and handle-relative opens
  on Windows. A hand-written walk would repeat it, and would need a separate
  Windows path. See ADR-0177.
- **Keep the existing no-follow checks.** `os.Root` alone guarantees
  containment, not "never follow a link". The scan's contract (ADR-0173) is
  the stronger one, so the checks stay alongside it.

## Acceptance evidence

Each Core Feature needs positive and negative evidence in the Task Graph. The
negative cases carry the weight:

- a Result section under `references/` that is still read;
- an authored section after a Result that is still read;
- a directory swapped for an outside symbolic link that yields no read;
- a normal tracked file that still yields its citations.

The outside-evidence row rests on published sources this Spec did not
produce: the Go blog's description of `os.Root` and its race model
(https://go.dev/blog/osroot), the `os.Root` documentation for Go 1.26
(https://pkg.go.dev/os#Root), and MITRE CWE-367
(https://cwe.mitre.org/data/definitions/367.html). The QA gate reproduces the
blog's parent-directory symlink attack against the built scan and records where
each source came from. When a source cannot be fetched, the row is recorded as
blocked with its reason.

## Research basis

- **Secondbrain.** `wiki/index.md` and `qmd query "TOCTOU symlink openat
  O_NOFOLLOW traversal-resistant file access"` returned no relevant knowledge.
  The closest hit, the Roundfix mirror of Spec 0150's task_02, concerns writing
  through a symlinked Task path, not reading. Nothing in the brain covers
  traversal-resistant reads.
- **Exa.**
  - The Go blog "Traversal-resistant file APIs" (https://go.dev/blog/osroot,
    2025-03-12) describes exactly this race: a symlink created after the
    program's check. It introduces `os.Root`, implemented with the `openat`
    family on Unix and held handles on Windows.
  - `src/os/root_unix.go` at go1.26.5 shows each intermediate component opened
    with `O_NOFOLLOW|O_DIRECTORY`.
  - golang/go#67002 records the design and the residual (links inside the
    root are followed).

  These sources changed the decision from a hand-written `openat` walk to
  `os.Root`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
