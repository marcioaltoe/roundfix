---
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: active
created: 2026-10-26
surfaces: [backend, cli]
---

# Runs trim the shared build cache, and doctor warns on low disk

Minimal PRD for a bug fix, minted by `write-techspec`. Spec 0314 made
`roundfix gc` trim a declared **Shared Build Cache**, but only when someone
runs it, so this repository's Go build cache grew past its 20,480 MiB limit to
43.7 GB on 2026-10-22 and to 51 GB on 2026-10-25 and stopped Runs with no
space left on the device, while `roundfix doctor` never mentioned the full
disk. After this Spec every Implement Run and the Delivery Queue's repository
gate run the same trim by themselves, under a thirty-second budget, without
ever changing their outcome, and record what it did; and `roundfix doctor`
prints a `disk-space:` warning when the filesystem holding the worktree
location or a declared Shared Build Cache has less than 10 GiB free.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md`](references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md),
   `type: fix`, from the operator's queue records of 2026-10-22 (Spec 0300's
   Run Unresolved, 399 MB free, cache 43.7 GB), 2026-10-25 (157 MiB free,
   cache 51 GB, Specs 0315 and 0307 failed for lack of space) and 2026-10-26
   (cache back at 17 GB by 15:30 BRT).

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `d08f8557`, where Specs 0319 and 0320
are active and deliver first. Neither touches the GC Command, `doctor`, the
Implement Run start, the repository gate or `internal/sharedcache`; both
change `CONTEXT.md`, the Roundfix skill's version fields,
`skills/testdata/owned-skill-versions.json` and
`docs/references/behavior-surfaces.json`, which this Spec's documentation
Task re-records with their commands on its own Delivery Base.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or schema
  is minted; the Run Event reuses the `daemon.status` kind with a new payload
  and adds no kind. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the trim and the free-space read touch
  local files only. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0327 (GC trims a declared shared
  build cache) keeps every trim rule; ADR-0334, recorded by this planning
  change, adds the automatic triggers, the budget and the best-effort rule.
  ADR-0327: "A trim acts only when the cache exceeds `max_size_mb`, removes
  idle entries oldest first until the cache is within its limit" Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix
  skill's `deliver`, `implement`, `setup` and `storage` references with its
  version fields, and the command test harness in `internal/cli/cli_test.go`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A declared Shared Build Cache is trimmed at the start of every Implement
  Run, direct or started by the Delivery Queue, and before every repository
  gate of the queue, with exactly ADR-0327's rules.
- The trim never fails, parks or blocks the Run or the gate: a refusal, an
  error or the end of its thirty-second budget is recorded and the work goes
  on.
- Each automatic trim leaves a record: a Run Event per declared cache in the
  Run, a line per cache in the Queue Owner's log, and a console line whenever
  it removed something, failed or stopped.
- `roundfix doctor` warns, without failing, when the filesystem that holds the
  worktree location or a declared Shared Build Cache has less than 10 GiB
  free, and names `roundfix gc` as the next action.
- A repository with no declaration, and a machine with enough free space, see
  every output unchanged.

## Core Features

1. A bounded trim engine. `internal/sharedcache` gains `MeasureContext` and
   `ApplyContext`, which stop between entries when their context ends;
   `Measure` and `Apply` keep their behavior.
2. The automatic trim. One function in `internal/cli` trims every declared
   cache under one budget and reports an outcome per cache (`trimmed`,
   `unchanged`, `kept`, `failed`, `stopped`). The Implement Run calls it right
   after the Run is recorded and appends a `daemon.status` Run Event per
   cache; the repository gate calls it before its Verification and writes to
   the Queue Owner's log. `roundfix gc` uses the same function and keeps its
   report.
3. The `disk-space:` doctor line, printed after `storage:` only when a
   location is below 10 GiB free (`DR-DISK-LOW`) or its free space cannot be
   read (`DR-DISK-UNREAD`), through an exported reader that reuses the one
   `gc compact` already uses.
4. Documentation. The GC, doctor, implement, deliver and configuration
   guides, the Roundfix skill's `storage`, `setup`, `implement` and `deliver`
   references, and the revised **Shared Build Cache** entry in `CONTEXT.md`. ADR-0334 is recorded by this planning
   change.

## User Experience

The operator's Delivery Queue starts Spec 0322's Implement Run while the
shared cache holds 26 GB against its 20,480 MiB limit. Before the Run
Worktree exists the Run removes the idle entries, oldest first, prints
`roundfix: Shared Build Cache <path>: removed <n> idle entries (<bytes> bytes)
at Run start`, and records the outcome in the Run's events; `roundfix events`
shows it. A later `deliver retry` that resumes at the gate trims again first.
When the disk is nevertheless low, `roundfix doctor` prints
`disk-space: warn (DR-DISK-LOW: Shared Build Cache <path> has <n> MiB free,
below 10240 MiB; next: roundfix gc)`.

## Non-Goals / Out of Scope

- Refusing to start a Run, or parking a queue item, below a free-space floor
  (the brief's optional third decision). The measured failures happened when
  space ran out during a Run that started with room, which a start floor does
  not prevent; the automatic trim removes the growth that filled the disk,
  and the doctor warning covers what is left. It can be proposed later with
  its own park class.
- Trimming at the start of `resolve` and `watch` Runs, or before every Task
  (ADR-0334 → Considered Options).
- A configuration key for the doctor threshold or the trim budget. Both are
  fixed and named in the guides; the warning never fails Doctor, so a machine
  with a smaller disk sees a warning, not a failure.
- Any change to ADR-0327's trim rules, to `worktree.shared_caches`, to this
  repository's Makefile or `.roundfixrc.yml`, or to Worktree Build Caches.
- Grouping the warning by filesystem: each location is named on its own, so a
  shared filesystem can be named twice with the same number.

## Success Metrics

1. Success Metric: an Implement Run whose declared cache is above its limit
   removes exactly the idle entries ADR-0327 selects before its first Task
   starts, keeps the fresh entry and the README, ends as it would without the
   cache, and records one `daemon.status` Run Event with outcome `trimmed`
   and the removed count and bytes.
2. Success Metric: a refused cache (no README), and a trim whose budget ends,
   leave the Run's exit code and outcome unchanged and record `kept` or
   `stopped` with the cause; a Run with no declaration records no trim event
   and prints no trim line.
3. Success Metric: the repository gate removes the idle entries before its
   Verification runs (a Verification that fails while an idle entry exists
   passes), writes the trim line to the Queue Owner's log, and its pass or
   fail still comes only from the Verification.
4. Success Metric: `MeasureContext` and `ApplyContext` stop between entries
   when their context ends, and every existing `internal/sharedcache` and GC
   Shared Build Cache test passes unchanged.
5. Success Metric: `roundfix doctor` prints the `disk-space:` line between
   `storage:` and the next line when a location has 10,239 MiB free or less,
   prints nothing new at exactly 10,240 MiB, reports `DR-DISK-UNREAD` when the
   read fails, and keeps its exit code; no existing doctor test reads the
   host's disk.
6. Success Metric: the guides, the Roundfix skill references and
   `CONTEXT.md` describe the triggers, the budget, the records and the
   doctor line, the owned-skill version is recorded, and the Behavior Surface
   record is re-recorded last.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- `docs/references/2026-10-24-test-suite-audit.md`, committed by Spec 0311's
  planning change before this Spec existed: "The shared build cache
  `roundfix-go-build-501` under `$TMPDIR` held 31 GB on 2026-10-24 (43.7 GB
  on 2026-10-22, when it filled the disk)", and "`roundfix gc` does not
  reclaim it" — the growth the automatic trim answers.
- The Go command documentation, `go help cache`
  (https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching), read on
  2026-10-26: "The cache is safe for concurrent invocations of the go
  command." and "The go command periodically deletes cached data that has
  not been used recently." The concurrent trims rely on the first; the
  second is the slow trim that did not keep up.
- The Bazel remote caching documentation (https://bazel.build/remote/caching),
  read on 2026-10-26: "Bazel will automatically garbage collect the disk cache
  while idling between builds" and "As an alternative to automatic garbage
  collection, we also provide a tool to run a garbage collection on demand."
  Another build tool bounds a shared disk cache automatically between builds
  and keeps an on-demand command, the shape ADR-0334 takes.

## Unreachable Acceptance

- criterion: after release, on the operator's delivery machine, the shared
  cache stays near its declared 20,480 MiB across a day of queued Runs, each
  Run's events show its trim outcome, and `roundfix doctor` prints no
  `disk-space:` warning while more than 10 GiB is free.
  reason: it needs the released binary, the real Roundfix Home and the
  machine's real disk, which no Verification or QA row may read or fill.
  satisfied-by: the operator's queue records after the release that carries
  this Spec, recorded session evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Code read at `d08f8557`: `sweepSharedCaches`
  (`internal/cli/gc_shared_cache.go`) is called only by `runGC`;
  `sharedcache.Measure` and `sharedcache.Apply` take no context. Every
  Implement Run, direct or the queue's `implement --spec` child
  (`commandDeliveryWorkflow.RunSpec`), passes `runImplementCommand`, which
  records the Run with `CreateRun` and then calls `createRunWorktree`; the
  queue builds outside a Run only in `commandDeliveryWorkflow.Gate`.
  `availableTemporaryCapacity` in `internal/store` already reads a
  filesystem's free bytes on Unix and Windows for `gc compact`.
- Measurement on this repository's machine, read-only, at 16:40 BRT:
  `sharedcache.Measure` read the real 21.6 GB cache (about 38,000 files) in
  1.2 s; on a synthetic cache of 40,000 idle entries in a scratch directory,
  `Measure` took 1.4 s and `Apply` removed all of them in 5.6 s. Thirty
  seconds covers twice the worst measured cache.
- Threshold: the failures happened at 399 MB and 157 MiB free, and the
  operator's earlier disk interventions (2026-10-07 and 2026-10-09) at
  255 MB, 2.1 GB and 5.3 GB were already emergencies. ADR-0327 measured
  33 GB of growth in fifteen hours with three concurrent Runs, about 2.2 GB
  an hour, so 10 GiB leaves several hours of warning.
- Secondbrain, read only: `qmd query "roundfix shared build cache gc disk full
  GOCACHE" --all --files --min-score 0.3` returned the Roundfix mirror's
  ADR-0327 (score 1.00), its Worktree Build Cache ADR, the 0314 Archive Record and a digest of
  2026-10-12 on the Go cache in the sandbox; `qmd query "doctor low disk free
  space warning"` returned only the mirror's doctor guide (0.93). None
  records a decision against an automatic trim or a disk warning.
- Exa: the Go command and Bazel documentation quoted under Acceptance
  evidence.

## Glossary

- changes: **Shared Build Cache**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 2 brief for Spec 0321 of 2026-10-26, the maintainer's cycle AQ answers quoted in `_authorization.md`, and the measurements above.

Taken from those sources, with no further interview:

- Triggers: the start of every Implement Run (after `CreateRun`, before the
  Run Worktree) and the queue's repository gate (before its Verification);
  not `resolve`, `watch` or each Task. See ADR-0334.
- Budget: thirty seconds for all declared caches at one trigger; the trim
  stops between entries and the next trigger continues.
- Records: a `daemon.status` Run Event per declared cache in an Implement
  Run, a line per cache in the Queue Owner's log at the gate, and a console
  line only for `trimmed`, `failed` and `stopped`.
- Doctor: a warning line printed only below 10 GiB free or when the read
  fails, after `storage:`, not configurable, never `failed`.
- No start floor and no new park class (Non-Goals).
- ADR-0334, a new accepted ADR, records the automatic triggers; ADR-0327 is
  unchanged.

## Open Questions

None.
