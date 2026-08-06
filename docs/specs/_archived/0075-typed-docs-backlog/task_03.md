---
task: task_03
spec: 0075-typed-docs-backlog
status: completed
type: docs
complexity: low
---

# Task 03: Adopt the contract in this repository

## Overview

The clauses are what adopting repositories receive; this repository is one of
them. This slice re-applies the baseline so the checked-in guide carries the new
contract, and seeds `docs/backlog/` with its first real entry using the template
verbatim — which is also the cheapest proof the template is usable.

## Requirements

1. MUST re-apply the baseline so `docs/agents/docs-layout.md` carries the
   backlog contract.
2. MUST create `docs/backlog/` with at least one real entry, using the template
   verbatim rather than an invented shape.
3. MUST make that first entry a genuine item, not a placeholder. The
   verification-performance contract the 2026-08-03 handoff describes as
   belonging in `docs/backlog/` as a `perf` entry is the natural candidate.
4. MUST NOT migrate any existing finding into the backlog. Deprecating findings
   into a backlog type is explicitly deferred.
5. MUST leave every other generated guide unchanged.

## Subtasks

- [ ] Re-apply the baseline and confirm the guide updated.
- [ ] Write the first backlog entry from the template.
- [ ] Confirm no finding was migrated and no other guide moved.

## Acceptance Criteria

- [ ] `docs/agents/docs-layout.md` documents `docs/backlog/` and its contract.
- [ ] `docs/backlog/` holds at least one entry whose frontmatter matches the
      contract exactly.
- [ ] That entry is a real item with real content, not a placeholder.
- [ ] No file under `docs/findings/` moved or changed.
- [ ] No generated guide other than the layout guide changed.

## Verification

- `grep -q "docs/backlog" docs/agents/docs-layout.md` — expected: exit 0.
- `ls docs/backlog/*.md | grep -q .` — expected: exit 0; the directory holds an
  entry.
- `head -8 docs/backlog/*.md | grep -qE "^type: (feat|fix|perf|refactor)( +#.*)?$"`
  — the optional trailing comment is not laxity: the template this Spec adopts
  documents the enum inline, as `type: perf # feat | fix | perf | refactor`, so
  an end-anchored pattern rejects the very shape the contract defines. The
  value is still pinned to the four members at the start of the line. Measured
  on 2026-08-05: the original pattern failed against a conforming file.
  — expected: exit 0; the entry's type is one of the four.
- `git diff --name-only HEAD -- docs/findings | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no finding moved.
- `make verify` — expected: exit 0.

## References

- `_prd.md` → Core Features; Non-Goals (no finding migration).
- `_techspec.md` → Build Order 3.
- `docs/handoffs/2026-08-03-after-the-0.3.1-release.md` → the `perf` entry it
  describes.

## Result

### Implementation

- Updated only the `guide.docs-layout` managed marker in
  `docs/agents/docs-layout.md` to its canonical regenerated postimage from
  task_02's formatter golden. The adopted guide now carries the backlog
  operational contract, all four body templates, the finding-versus-intent
  boundary, and the `docs/backlog/` directory job.
- Added `docs/backlog/2026-08-03-verification-performance-contract.md` from
  the generated copyable contract: the frontmatter keeps the template's field
  order, enum comments, and lifecycle fields, and the body uses the `perf`
  template's `Slow`, `Measured`, and `Target` sections without adding an
  invented shape.
- Grounded the entry in the 2026-08-03 handoff and the archived Spec 0071
  report. It records the measured 4.9s unchanged local gate, approximately 5s
  light-package edit, 48.5s `internal/cli` edit, 54.5s
  `internal/baseline` edit, and 88.9s fresh complete run, then carries forward
  the handoff's target of at most 10s unchanged, at most 60s after a typical
  change, and a fresh complete CI tier.
- Did not migrate, edit, rename, or delete any finding.

### Focused checks

- Read the Daemon diagnostic artifact for attempt 1; it contained no command
  output beyond the recorded exit status for the missing layout-guide text.
- `rtk rg -n -A55 "Backlog Operational Contract" internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
  printed the canonical frontmatter and `perf` body template used by the new
  entry.
- A marker-bounded `awk` extraction compared with `diff --unified=0` showed no
  difference between the adopted `guide.docs-layout` block and task_02's
  canonical formatter golden.
- `rtk read docs/backlog/2026-08-03-verification-performance-contract.md`
  showed the exact five-field frontmatter and real `Slow`, `Measured`, and
  `Target` content.
- `rtk git -c core.fsmonitor=false status --short` showed only this Task file
  and the new `docs/backlog/` path after the temporary diagnostic test was
  removed; no `docs/findings/` path or generated guide is changed.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  them after this handback.

### Acceptance-criterion evidence

1. **Satisfied — local layout guide carries the backlog contract:** the
   complete `guide.docs-layout` managed block is byte-identical to task_02's
   canonical regenerated formatter golden and documents `docs/backlog/`, its
   lifecycle, frontmatter, templates, and boundary with findings.
2. **Satisfied — exact backlog frontmatter:** the new entry uses `type: perf`,
   `status: open`, `created: 2026-08-03`, `spec: null`, and `reason: null` with
   the generated template's inline enum and lifecycle comments unchanged.
3. **Satisfied — genuine content:** the entry cites the archived measurement
   report, records its actual timings, and states the handoff's concrete
   two-tier budgets.
4. **Satisfied — findings unchanged:** the fresh changed-path inventory
   contains no path under `docs/findings/`.
5. **Satisfied — other generated guides unchanged:** the feedback repair
   copied only the canonical `guide.docs-layout` managed block. The fresh
   changed-path inventory contains no other file under `docs/agents/`.

### Follow-ups

- Repair the interactive Baseline workflow so its classification review uses
  the same catalog-aware carrier classifications as the final planner.
- Reconcile the active Baseline catalog's `guide.autonomous-work` clause with
  the current generated guide. The public full-repository plan tried to mutate
  that unrelated guide, so this feedback repair adopted the already-generated
  canonical layout postimage directly and kept the Task diff bounded. Both
  Baseline repairs are outside task_03's documentation-adoption slice.
