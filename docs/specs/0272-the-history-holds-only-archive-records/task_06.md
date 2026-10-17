---
task: task_06
spec: 0272-the-history-holds-only-archive-records
status: pending
type: backend
complexity: low
---

# Task 06: A tracked symbolic link is skipped by the history citation scan

## Overview

The first delivery of task_04 on 2026-10-17 found that both history commands
stop on this repository's tracked `CLAUDE.md`, a symbolic link to
`AGENTS.md`: `HistoryCitations` reads every tracked Markdown path and refuses
any that is not a regular file, so `roundfix baseline update` prints
`History result: citing file CLAUDE.md is not a regular repository file` and
plans nothing. Adopters keep the same `CLAUDE.md` link. Git tracks the link,
not the text; the target is tracked and scanned under its own path. This
corrective Task makes the citation scans skip a tracked symbolic link without
reading or writing through it, so task_04's Managed Refresh runs on the tree
as it is. It is verifiable on its own through a fixture with a tracked link.

## Requirements

1. MUST make `HistoryCitations` in `internal/spec/history_entries.go` skip a
   tracked Markdown path whose own entry is a symbolic link: it is not read,
   not reported and never rewritten, and its target is scanned only under the
   target's own tracked path. A regular file under a symbolically linked
   directory is still refused as today, and every other refusal keeps its
   text.
2. MUST make the citation scan of the `superseded` unit (API Contract 8) skip a
   tracked symbolic link the same way, so a link to a file outside the
   repository is never followed.
3. MUST NOT change the apply path's refusal of a symbolic link it is asked to
   rewrite, the plan output shape, or any other unit.
4. MUST add, each calling `t.Parallel()` first:
   `TestHistoryCitationsSkipATrackedSymbolicLink` in
   `internal/spec/history_symlink_test.go` (a fixture whose tracked
   `AGENTS.md` cites a removed path and whose tracked `CLAUDE.md` links to it:
   one citation, in `AGENTS.md`; `CLAUDE.md` is still a link to `AGENTS.md`
   after apply, and `AGENTS.md` carries the Git Locator), and
   `TestHistorySanitizeAppliesBesideATrackedSymbolicLink` in
   `internal/cli/history_symlink_test.go` (`roundfix history sanitize
   --apply` on the same kind of fixture exits 0 and leaves the link a link).
5. MUST leave every existing test passing unchanged.

## Subtasks

- [ ] Skip a tracked symbolic link in `HistoryCitations` and the `superseded` scan.
- [ ] Add the package and command tests on a fixture with a tracked link.

## Acceptance Criteria

- [ ] A tree with a tracked `CLAUDE.md` link plans and applies its history without a refusal.
- [ ] The link is never read through, reported or rewritten, and stays a link.

## Consumers

- History Citation planning for both history commands — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- interface: `internal/spec/history_entries.go`
- interface: `internal/spec/history_entries_test.go`
- creates: `internal/spec/history_symlink_test.go`
- creates: `internal/cli/history_symlink_test.go`
- interface: `internal/cli/history_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestHistoryCitationsSkipATrackedSymbolicLink$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistoryCitationsSkipATrackedSymbolicLink ' || { printf 'missing passing test TestHistoryCitationsSkipATrackedSymbolicLink\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestHistorySanitizeAppliesBesideATrackedSymbolicLink$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistorySanitizeAppliesBesideATrackedSymbolicLink ' || { printf 'missing passing test TestHistorySanitizeAppliesBesideATrackedSymbolicLink\n' >&2; exit 1; }; go test -count=1 -run 'History' ./internal/spec ./internal/cli || exit 1` — expected: exit 0. Before this Task the two tests do not exist, so the PASS checks fail.
- `dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$dir/roundfix" ./cmd/roundfix || exit 1; out="$(HOME="$dir" "$dir/roundfix" baseline update --repo . --no-skills --format text 2>&1)"; code=$?; if printf '%s\n' "$out" | grep -qF -- 'is not a regular repository file'; then printf '%s\n' "$out" >&2; exit 1; fi; case "$code" in 0|3) exit 0 ;; *) printf 'baseline update exited %s:\n%s\n' "$code" "$out" >&2; exit 1 ;; esac` — expected: exit 0. Before this Task the plan stops on `CLAUDE.md` with that refusal, so the command fails.

## References

- `_prd.md` → Core Feature 4; User Story 3
- `_techspec.md` → API Contract 8; API Contract 10; Invariant 5
- ADR-0278
