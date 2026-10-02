---
task: task_03
spec: 0214-measure-before-changing
status: completed
type: docs
complexity: medium
---

# Task 03: Measure who reads archived evidence and propose, removing nothing

## Overview

The history root holds most of the repository's tracked files, and nobody
has shown which archived evidence any agent, test, check or tool still
reads. This Task runs the archive measurement of `_techspec.md` at its
starting commit, writes the inventory, the readers, the ablation, the Agent
reads, the mirror and a proposal to one reference document, and writes a
Backlog Entry only when rule 3 finds a candidate. It deletes, moves and
rewrites nothing under the history root; the ablation happens in a disposable
clone outside the repository. It is verifiable on its own: the document's
counts equal what Git reports at the commit it names, and no archived file
changed.

## Requirements

1. MUST NOT delete, move, rename or rewrite any file under `docs/history/`. The maintainer decided on 2026-10-01 that this work measures and proposes only, and that any removal waits for the maintainer's explicit approval in a later change (ADR-0215).
2. MUST perform `_techspec.md` → The archive measurement, steps 1 to 5, at this Task's starting commit, and write `docs/references/archived-evidence-measurement.md` with, in this order near its top, the lines `Measured at: <40-hex commit>`, `History files: <n>`, `History bytes: <n>`, `QA evidence files: <n>`, `QA evidence bytes: <n>`, `Binary files: <n>` and `Binary bytes: <n>`, each count as that section defines it.
3. MUST give the document the sections `## Inventory`, `## Readers`, `## Ablation`, `## Agent reads`, `## Secondbrain mirror`, `## Proposal` and `## What removal would not reclaim`, and the line `Nothing under docs/history was deleted, moved or rewritten by this measurement.`
4. MUST run the ablation only in a clone created under a temporary directory outside the repository and outside Roundfix Home, record each command with its exit code and the names of failing tests, record the five `git grep` timings before and after the removal with their medians, and remove that clone afterwards.
5. MUST open the Run Database only as `file:<home>/.roundfix/roundfix.db?mode=ro&immutable=1`, read only `agent.tool_started` events, and record the counts by kind of path (QA evidence, QA Report, core Spec artifact, other), by tool kind and by distinct Run, with the query used. When the Run Database is absent, the section says so.
6. MUST record the Secondbrain mirror's file and byte counts for its copy of the history root, read-only, or state that the mirror is absent.
7. MUST apply `_techspec.md` → Decision rules, rule 3, and write the line `Candidates: none` or `Candidates: <kinds, comma-separated>` using the kinds `qa-evidence` and `binaries`. The Proposal states, for each candidate, the files and bytes it would remove, the Agent reads it would lose, the changes that must land first (each failing test of the ablation), and that removal waits for the maintainer's explicit approval. The last section states that history keeps every removed version, citing the Git book page of `_prd.md` → Acceptance evidence.
8. MUST write `docs/backlog/<date>-remove-archived-evidence-nobody-reads.md`, a `perf` Backlog Entry with status `open`, only when the candidate set is not empty, and its Target MUST say that the removal waits for the maintainer's explicit approval. Its file name carries the date the Task runs, so it is not declared in Context and the Daemon records it under `## Recorded paths`.

## Subtasks

- [ ] Inventory the history root at the starting commit.
- [ ] List the static readers outside the history root.
- [ ] Run the ablation in a disposable clone and time the searches.
- [ ] Count Agent reads and measure the mirror.
- [ ] Apply rule 3, write the proposal and, if it fires, the Backlog Entry.

## Acceptance Criteria

- [ ] The document's seven header lines equal what `git ls-tree -r -l` reports for the commit it names, and that commit is an ancestor of `HEAD`.
- [ ] No file under `docs/history/` was deleted or renamed between that commit and the working tree.
- [ ] The document carries every required section and the no-removal line, and its `Candidates:` line matches the Backlog Entry's presence.

## Context

- creates: `docs/references/archived-evidence-measurement.md`
- instruction: `docs/agents/docs-layout.md`
- instruction: `docs/agents/secondbrain.md`
- instruction: `internal/spec/archive.go`

## Verification

- `f=docs/references/archived-evidence-measurement.md; test -f "$f" || { printf 'missing %s\n' "$f" >&2; exit 1; }; sha="$(sed -n 's/^Measured at: \([0-9a-f]\{40\}\)$/\1/p' "$f")"; test -n "$sha" || { printf 'no Measured at line\n' >&2; exit 1; }; git merge-base --is-ancestor "$sha" HEAD || { printf 'measured commit is not an ancestor of HEAD\n' >&2; exit 1; }; tmp="$(mktemp -d)" || exit 1; git ls-tree -r -l "$sha" -- docs/history > "$tmp/tree" || exit 1; awk -F'\t' '{ split($1, meta, " "); size = meta[4] + 0; files++; bytes += size; if ($2 ~ /^docs\/history\/specs\/[^\/]+\/qa\/evidence\//) { ef++; eb += size } lower = tolower($2); if (lower ~ /\.(png|jpg|jpeg|gif|pdf|db|zip)$/) { bf++; bb += size } } END { printf "History files: %d\nHistory bytes: %d\nQA evidence files: %d\nQA evidence bytes: %d\nBinary files: %d\nBinary bytes: %d\n", files, bytes, ef, eb, bf, bb }' "$tmp/tree" > "$tmp/want" || exit 1; awk '$0 ~ "^(History files|History bytes|QA evidence files|QA evidence bytes|Binary files|Binary bytes): [0-9]+$"' "$f" > "$tmp/have" || exit 1; cmp "$tmp/want" "$tmp/have" || { printf 'counts differ from Git at %s\n' "$sha" >&2; exit 1; }; removed="$(git diff --name-only --diff-filter=DR "$sha" -- docs/history)" || exit 1; test -z "$removed" || { printf 'archived files removed or renamed:\n%s\n' "$removed" >&2; exit 1; }` — expected: exit 0; before this Task the document does not exist, so the command fails.
- `f=docs/references/archived-evidence-measurement.md; test -f "$f" || exit 1; for phrase in "## Inventory" "## Readers" "## Ablation" "## Agent reads" "## Secondbrain mirror" "## Proposal" "## What removal would not reclaim" "Nothing under docs/history was deleted, moved or rewritten by this measurement." "git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery"; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$f" "$phrase" >&2; exit 1; }; done; candidates="$(sed -n 's/^Candidates: //p' "$f")"; test -n "$candidates" || { printf 'no Candidates line\n' >&2; exit 1; }; count=0; for entry in docs/backlog/*-remove-archived-evidence-nobody-reads.md; do test -f "$entry" && count=$((count + 1)); done; if test "$candidates" = none; then test "$count" -eq 0 || { printf 'a Backlog Entry exists for no candidate\n' >&2; exit 1; }; else test "$count" -eq 1 || { printf 'candidates %s need exactly one Backlog Entry\n' "$candidates" >&2; exit 1; }; for entry in docs/backlog/*-remove-archived-evidence-nobody-reads.md; do tr -s '[:space:]' ' ' < "$entry" | grep -qF -- "explicit approval" || { printf 'the Backlog Entry does not state the approval\n' >&2; exit 1; }; done; fi` — expected: exit 0; before this Task the document does not exist, so the command fails.

## Result

Implemented the archive measurement in `docs/references/archived-evidence-measurement.md` at starting commit `513b22ebab2d62a23b4222073c937c2862907ec2`. It records the Git inventory, static readers, disposable-clone ablation, immutable Run Database counts, Secondbrain mirror, and rule 3 proposal. No file under `docs/history/` changed, and no Backlog Entry was created because `Candidates: none`.

Focused evidence: `git ls-tree -r -l` reported 5,271 files and 42,686,990 bytes, with 2,089 QA evidence files / 12,725,738 bytes and 18 binary files / 12,179,652 bytes; the measured commit is an ancestor and has no deleted or renamed history path. The disposable clone was removed after five pre-removal and five post-removal `git grep` timings, all exit 0, with medians 0.15s and 0.19s. The ablation's `make verify` and `make verify-docs` exited 2 with the failing tests recorded in the document. The immutable query found 876 matching `agent.tool_started` events across 127 Runs, and the read-only mirror contained 5,271 files / 42,686,990 bytes. The daemon must run the declared Verification commands and own Task status and settlement.

Verification feedback repair: the first attempt compared the report with the
Task contract's exact `docs/history/specs/<slug>/qa/evidence/` expression;
the recorded QA totals used a broader prefix match. The header and Result
evidence now use the exact Git-derived values: 2,089 files and 12,725,738
bytes. The diagnostic artifact was inspected without copying its body here.

## References

- [_prd.md](_prd.md) — Goals 4 and 5; User Stories 4 and 5; Core Features 9, 10 and 11; Success Metrics 4 and 5; Acceptance evidence
- [_techspec.md](_techspec.md) — The archive measurement; Decision rules; Integration Points; Testing Approach 6; Build Order 3
- [references/2026-09-30-archived-specs-keep-evidence-nobody-reads.md](references/2026-09-30-archived-specs-keep-evidence-nobody-reads.md)
- ADR-0215; ADR-0120; ADR-0121

## Carry-forward provenance

- Source Run: `run_20261002T165011Z_2a13efd81491b9d8`
- Source commit: `2b6e74ccab0b5aa850277b693dea7e0a8b3074a6`
