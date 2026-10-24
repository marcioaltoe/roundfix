---
spec: 0314-gc-trims-the-shared-build-cache
status: active
created: 2026-10-24
surfaces: [backend, cli]
---

# GC trims the shared build cache

The Go build cache that this repository's Makefile shares across every
worktree grew to 43.7 GB and filled the disk on 2026-10-22, and held 33 GB
again fifteen hours after it was emptied, because each Run worktree compiles
under its own path and the go command deletes an entry only after five days
without use. The GC Command reclaims Worktree Build Caches inside ended
worktrees but leaves a shared cache alone. After this Spec a repository or a
user declares a **Shared Build Cache** in `worktree.shared_caches`, and
`roundfix gc` trims it above its declared size by removing only cache
entries no build used for the declared idle time, oldest first, as ADR-0327
decides; `roundfix gc --dry-run` reports the same without writing. This
repository declares its Makefile cache.

## Sources

1. The maintainer's answer of 2026-10-24, about 01:50 BRT, recorded in this
   Spec's authorization record: the cache reclamation becomes a change to the
   GC Command.
2. The test suite audit
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md),
   committed by Spec 0311's planning change: 33 GB in 69,988 files at 01:20
   BRT on 2026-10-24, none older than fifteen hours, 31.4 GB untouched for
   two hours and 18.4 GB for four; `du -sk` read it in 2.0 s.

## Prerequisites

None in `requires`. This Spec is authored against `56dd1a76` (v0.116.0) and
delivers after Specs 0311, 0312 and 0313 and before 0307, 0309 and 0310, by
the maintainer's answer quoted in `_authorization.md`. Its sources cite the
audit reference, which Spec 0311's planning change adds, so the operator
merges the plans together. Spec 0306, delivered before it, does not touch
the GC Command.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted; the new configuration key holds paths, sizes and
  durations. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; GC reads and removes local files only.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0327 (GC trims a declared shared
  build cache), recorded by this planning change, is the decision the
  implementation Tasks carry out. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix
  skill's storage reference with its version fields and this repository's
  Project Config declaration of its cache. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0314-gc-trims-the-shared-build-cache/_authorization.md`;
  bounded files: `.agents/skills/roundfix/references/storage.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.roundfixrc.yml`.

## Goals

- A declared Shared Build Cache never stays above its limit after
  `roundfix gc`, as long as enough of it is idle.
- No file that is not an idle Go build cache entry is ever removed: not a
  fresh entry, not the README or `trim.txt`, not a file in a directory
  without the go command's README, not through a symbolic link.
- `roundfix gc --dry-run` reports the same plan without writing.
- A repository without a declaration sees GC's output unchanged.

## Core Features

1. The trim engine. A new package `internal/sharedcache` measures a cache
   directory and plans the removal of idle entries, oldest first, until the
   cache is within the limit, and applies a plan by rereading each entry's
   modification time just before removing it.
2. The configuration key. `worktree.shared_caches` lists entries with
   `path`, `kind` (`go-build`), `max_size_mb` and `idle` (at least `2h`), in
   User Config or Project Config; the path is absolute or starts with
   `{tmpdir}` or `{home}`, and `{uid}` stands for the user id.
3. The GC section. After the Item Binary section, `roundfix gc` and
   `roundfix gc --dry-run` print one block per declared cache: its resolved
   path, its size against its limit, and the entries removed or
   reclaimable, or why it was kept.
4. This repository's declaration and the documents: `.roundfixrc.yml`
   declares `{tmpdir}/roundfix-go-build-{uid}` with 20,480 MB and `4h`; the
   GC guide, the configuration guide, the Roundfix skill's storage reference
   and the glossary describe the key and the section.

## User Experience

```text
  Shared Build Cache: /var/folders/…/T/roundfix-go-build-501
    Size: 35433480192 bytes of 21474836480 allowed
    Entries reclaimable: 31218 (14982201344 bytes)
```

With `roundfix gc` the last line reads `Entries removed`. A cache the trim
refuses prints `Kept: <reason>` instead of the size lines.

## Non-Goals / Out of Scope

- Kinds other than `go-build`.
- Running the trim automatically at the start of a Run or a Delivery Queue.
- Changing Worktree Build Caches, Item Binaries, Run Retention or
  `gc sanitize`.
- The Makefile: it keeps creating the cache exactly as today.

## Success Metrics

1. Success Metric: on a fixture cache of four 1,000-byte entries aged 10 h,
   5 h, 3 h and 30 min plus a non-entry file aged 48 h, a 3,500-byte limit
   and a 2 h idle time, the plan removes the 10 h and 5 h entries and the
   apply removes exactly those two.
2. Success Metric: under its limit a cache loses nothing; over a zero limit
   only the three entries older than two hours are planned; an entry touched
   after planning survives the apply.
3. Success Metric: a directory without the README, a symbolic link and an
   idle time below two hours are refused.
4. Success Metric: `roundfix gc --dry-run` on a fixture cache of three 1 MiB
   entries (10 h, 3 h, fresh) with a 1 MB limit reports
   `Entries reclaimable: 2 (2097152 bytes)` and changes nothing; `roundfix
   gc` removes those two and keeps the fresh entry and the README; with a
   9 MB limit it removes nothing.
5. Success Metric: an unknown kind, a relative path, a `..` segment, an
   unknown placeholder, a missing size or an idle time below two hours fails
   configuration loading naming `worktree.shared_caches`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The go command's source in the Go 1.26 toolchain, read under
  `$(go env GOROOT)/src/cmd/go/internal/cache/` offline: its comment says the
  go command sets the modification time of a cache file on each use at most
  once per `mtimeInterval`, one hour, and its constants are
  `mtimeInterval = 1 * time.Hour` and `trimLimit = 5 * 24 * time.Hour`.
- `go help cache` in the audited checkout: "The go command periodically
  deletes cached data that has not been used recently."

## Validation

The decisions were checked on 2026-10-24 against these sources.

- A prototype in a scratch worktree of `56dd1a76`: the engine, the key and
  the GC section; the four engine tests, the two configuration tests and the
  three GC tests passed; `roundfix gc --dry-run` with `HOME` and `TMPDIR` at
  temporary directories printed the section for a fixture cache; the
  Behavior Surface fingerprint needed re-recording after the key, and the
  documentation contracts and the skills tests passed after the skill's
  version was recorded. Each sabotage the QA lists failed its named test.
- Secondbrain, read only: `qmd query "roundfix test suite flaky slow
  performance" --all --files --min-score 0.3` returned nothing on build
  caches. No Inbox Entry was written.
- Exa: none needed; the go command's own source and documentation are read
  offline.

## Glossary

- adds: **Shared Build Cache**

## Decisions

- 2026-10-24 (about 01:50 BRT): the maintainer's build cache answer,
  recorded in this Spec's authorization record.

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 4 briefing, the maintainer's answers of 2026-10-24 and the measurements recorded in the audit.

Taken from those sources, with no further interview:

- Declared caches only, found by placeholders rather than by reading build
  tooling (ADR-0327).
- Trim idle entries oldest first down to the limit, never empty the cache.
- Minimum idle time of two hours; this repository uses four, from the
  measured ages.
- Prove "never in use" by the go command's documented refresh interval and
  a reread before each removal.

## Open Questions

None.
