---
task: task_02
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
status: completed
type: docs
complexity: low
---

# Task 02: The guide, README, glossary and ADR-0370 describe the cache set with go env

## Overview

task_01 made the `Makefile` use the Go build cache set with
`go env -w GOCACHE` and fall back small with a warning, as gesttione-drive
asked on 2026-11-12
([`references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md`](references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md),
reported by gesttione-drive). This Task brings the four texts that describe
the Makefile's cache in line: the repository rule in
`docs/agents/specific-repository.md`, the README paragraph, the
**Build Cache Bound** entry of `CONTEXT.md`, and ADR-0370, which receives the
amendment notice of ADR-0394. The authorized tooling path is exactly
`docs/agents/specific-repository.md`.

## Requirements

1. MUST, in `docs/agents/specific-repository.md`, replace the rule line that
   starts ``- `make` shares `roundfix-go-build-<uid>` `` with exactly this
   line:

   ```markdown
   - `make` uses the Go build cache set with `go env -w GOCACHE` when it is an existing writable directory and never creates it; when that cache is set but unusable, `make` warns once and uses the shared cache below with a 2048 MiB default bound. Otherwise `make` shares `roundfix-go-build-<uid>` under `$TMPDIR` (or `/tmp`), falls back to `$(CURDIR)/.gocache` when it cannot create and write the shared cache, and exports an explicit `GOCACHE` unchanged; a sandboxed Agent sets nothing; before it builds or tests, the top-level `make` holds the cache it chose to the Build Cache Bound with `cmd/gocache-bound` and never clears it (ADR-0370, ADR-0394).
   ```

2. MUST, in `README.md`, replace the three lines that end the build-cache
   paragraph, from ``system disk, point `TMPDIR` or `GOCACHE` at another volume; Go's own``
   to ``User Config under `worktree.shared_caches`.``, with exactly these
   lines, keeping the paragraph's first four lines:

   ```markdown
   system disk, set it on another volume with `go env -w GOCACHE=<dir>`:
   `make` then uses that cache, under the same bound, whenever it is an
   existing writable directory, and `go` commands outside `make` share it.
   When that volume is not mounted, `make` never creates the directory; it
   warns once and falls back to `$TMPDIR/roundfix-go-build-<uid>` with a
   2048 MiB default bound. A `GOCACHE` in the environment still wins, and
   `make` bounds it only when `GOCACHE_BOUND_MB` is also set.
   ```

3. MUST, through the `domain-modeling` skill, change the definition line of
   the **Build Cache Bound** entry of `CONTEXT.md` to exactly this line,
   keeping its heading and its `_Avoid_` line:

   ```markdown
   The size, 10,240 MiB by default, to which this repository's Makefile holds the Go build cache it chose before each top-level make that builds or tests, removing with the Shared Build Cache rules only entries no build used for two hours; the cache stays above it while its entries are younger. The cache it chose is the one set with `go env -w GOCACHE` when that is an existing writable directory, and otherwise its temporary cache; when the cache set with `go env` is unusable, the Makefile warns and the default is 2,048 MiB ([ADR-0370](docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md), [ADR-0394](docs/adr/0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md)).
   ```

4. MUST, in
   `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`,
   set `updated_at: 2026-11-12T21:00:00Z`; add, after the `## Decision`
   paragraph that ends "The Free-Space Floor and the Disk-Full Park reach
   every repository.", one blank line and this paragraph, unchanged:

   ```markdown
   **Amended by [ADR-0394](0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md) (2026-11-12).** When the environment does not define `GOCACHE`, the cache the Makefile chooses is first the one set with `go env -w GOCACHE`, if it is an existing writable directory, under the same 10,240 MiB default bound; when that cache is set but unusable, the Makefile warns once and falls back to the shared `roundfix-go-build-<uid>` or the `.gocache` fallback under a 2,048 MiB default bound. An explicit `GOCACHE`, the Free-Space Floor and the Disk-Full Park are unchanged, and every other decision of this record still holds.
   ```

   and append to the Considered Option that starts "Moving the cache to
   another volume." these two lines, indented like it:

   ```markdown
     ADR-0394 takes this option for a cache set with `go env -w GOCACHE`: the
     Makefile now follows it when it is an existing writable directory.
   ```

5. MUST NOT change any other line of these four files, the `Makefile`, any
   test, `docs/agents/setup-context.json`, any other guide, skill, ADR or
   user-guide page. The `./internal/docscontract` package stays green with
   these edits, so the Behavior Surface Record does not move and is not
   re-recorded.

## Subtasks

- [ ] Rewrite the Makefile cache rule in the repository guide and the README paragraph.
- [ ] Change the **Build Cache Bound** entry through `domain-modeling`.
- [ ] Write the ADR-0394 notice into ADR-0370.

## Acceptance Criteria

- [ ] The guide, README and glossary say that `make` uses the cache set with `go env` when usable, never creates it, and otherwise warns and falls back at 2048 MiB.
- [ ] ADR-0370 carries the dated notice, the note on its Considered Option and the new `updated_at`.
- [ ] Every documentation contract passes, including the glossary phrase test.

## Context

- interface: `docs/agents/specific-repository.md`
- interface: `README.md`
- interface: `CONTEXT.md`
- interface: `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- instruction: `docs/adr/0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md`
- instruction: `Makefile`
- instruction: `docs/agents/domain.md`
- instruction: `docs/agents/docs-layout.md`

## Consumers

- The **Build Cache Bound** entry and the repository rule — consumers: `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Build Cache Bound**: The size, 10,240 MiB by default' || { echo 'CONTEXT.md lost the Build Cache Bound entry' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'when the cache set with .go env. is unusable, the Makefile warns and the default is 2,048 MiB' || { echo 'CONTEXT.md lacks the cache set with go env' >&2; exit 1; }; tr -s '[:space:]' ' ' < README.md | grep -qF -- 'it warns once and falls back to' || { echo 'README.md lacks the fallback' >&2; exit 1; }` — expected: exit 0. Before this Task the **Build Cache Bound** entry names no cache set with `go env` and the README tells to point `TMPDIR` or `GOCACHE` at another volume, so the second check fails; afterwards the entry and the README describe the choice and the fallback.
- `tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qE -- 'when that cache is set but unusable, .make. warns once and uses the shared cache below with a 2048 MiB default bound' || { echo 'the repository rule lacks the fallback' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md | grep -qF -- '**Amended by [ADR-0394](0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md) (2026-11-12).**' || { echo 'ADR-0370 lacks the notice' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md | grep -qF -- 'ADR-0394 takes this option for a cache set with' || { echo 'ADR-0370 lacks the option note' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the rule and ADR-0370 say nothing of the cache set with `go env`, so the first check fails; afterwards the rule, the notice and the option note are present and every documentation contract, the glossary phrase test included, passes.

## References

- `_prd.md` → Core Feature 4; Glossary; Decisions
- `_techspec.md` → System Architecture; Coverage Map; Build Order 2
- ADR-0394; ADR-0370

## Result

Skills loaded: implement-task, domain-modeling, crafting-effective-readmes, tech-writer, knowledge-workspace, evidence-gate.

Rewrote the repository cache rule, README cache paragraph, Build Cache Bound glossary definition, and ADR-0370 amendment exactly within the requested documentation scope. Preserved the Daemon-owned `status: in_progress` and did not edit any other Task, `_tasks.md`, Makefile, test, manifest, guide, skill, ADR, or user-guide page.

Focused checks after the edit: inspected the four changed documentation blocks and the worktree diff; all required new cache-set, no-create, warning, 2048 MiB fallback, dated amendment, option note, and ADR link text is present. The Task's declared Verification and repository Verification were not run because this is a Daemon-assigned handoff; the Daemon runs them at settlement.

Acceptance evidence:

- Guide, README, and glossary: each now states that `make` follows a usable `go env -w GOCACHE` directory without creating it, and warns before the 2048 MiB fallback when unusable.
- ADR-0370: `updated_at` is `2026-11-12T21:00:00Z`, the dated ADR-0394 notice follows the Decision paragraph, and the Considered Option includes the requested two-line note.
- Documentation contracts: not run in this Daemon turn; the focused diff inspection confirms the authored contract text is present for the Daemon's declared Verification.
