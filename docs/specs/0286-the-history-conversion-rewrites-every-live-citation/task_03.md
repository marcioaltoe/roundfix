---
task: task_03
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: completed
type: backend
complexity: high
---

# Task 03: A pending Task's target survives, and a broken link is resolved by its file name

## Overview

The adopter's conversion turned four links to a file a pending Task will
create, and two links that were already broken, into text marked not
reachable (`_prd.md` → Sources 1, defects 3 and 4). This Task makes
`live-links` and `SC-HISTORY-RETAINED` leave a pending Task's `creates:`
target alone and resolve a missing link by its file name when exactly one
removed history file has that name. It closes the conversion with Surface
Transcript 1, which needs Tasks 01 and 02 as well.

## Requirements

1. MUST answer defects 3 and 4 of the Inbox Entry of 2026-10-19,
   [`references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt`](references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt)
   (`_prd.md` → Sources 1), by implementing API Contracts 4 and 5 of
   `_techspec.md` (ADR-0299): `PendingTaskCreates`, `RemovedHistoryPaths` and
   `ResolveMissingLinkByName` in a new `internal/spec/pending_targets.go`,
   used by
   `PlanLiveLinks` and by `SC-HISTORY-RETAINED`'s missing-link finding, as
   `_techspec.md` → Interfaces sketches them.
2. MUST state the hazard in the code's behavior, as API Contracts 4 and 5 do:
   - `PendingTaskCreates` reads only Task files of Spec directories whose PRD
     is `status: active` under the Spec Root, and a Task file that does not
     parse contributes nothing instead of failing the plan or the check;
   - `RemovedHistoryPaths` reads `HEAD`'s tree and history through Git only,
     limited to `docs/history/findings`, `docs/history/backlog` and
     `docs/history/adr`, with `--no-renames`, and returns none outside Git;
   - a file name matched by two or more candidates is never resolved; the
     link keeps today's not-reachable form.
3. MUST keep `PlanLiveLinks`' signature and its callers in `internal/cli`
   unchanged, and keep every existing `SC-HISTORY-RETAINED` finding and its
   wording; only the exempted links and the fix line of a resolved link
   change.
4. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and ending every child it starts in `t.Cleanup`:
   - in `internal/spec/live_links_test.go`:
     - `TestLiveLinksLeaveAPendingTaskCreatesTarget`: a link to a `creates:`
       target of a `pending` Task is neither listed nor rewritten; the same
       link becomes a Git Locator or the not-reachable form once the Task is
       `completed` or its Spec's PRD is not active;
     - `TestLiveLinksResolveABrokenLinkByTheRemovedFileName`: a link
       `../findings/<name>.md` whose file is removed by the plan from
       `docs/history/findings/<name>.md` gets that file's Git Locator, with
       the citation's `Target` the link's path; a file deleted by an earlier
       commit resolves the same way; two candidates with the same name keep
       the not-reachable form;
   - in `internal/spec/pending_targets_test.go`,
     `TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks`: a file
     deleted in the working tree, a file deleted by a commit and a file
     renamed by a commit are listed; a present file and a file outside the
     three families are not;
   - in `internal/speccheck/history_retained_test.go`:
     - `TestHistoryRetainedSkipsAPendingTaskCreatesTarget`: a Live Document's
       link to a pending Task's `creates:` target is not reported, and is
       reported again once the Task is `completed`;
     - `TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames`: a
       missing link resolved by its file name is still reported, and its fix
       line names the removed history file's Git Locator;
   - in `internal/cli/history_citations_transcript_test.go`,
     `TestHistorySanitizeRewritesEveryCitationTranscript`, asserting Surface
     Transcript 1 of `_techspec.md` whole: every standard output line in
     order, the empty standard error and exit `0`, with `<b>` and `<commit>`
     the only controlled variations.

## Subtasks

- [ ] Read pending `creates:` targets and list Removed History Paths.
- [ ] Exempt pending targets and resolve broken links by file name in `live-links`.
- [ ] Do the same in `SC-HISTORY-RETAINED`'s missing-link finding.
- [ ] Add the tests and the Surface Transcript 1 test.

## Acceptance Criteria

- [ ] A link to a pending Task's `creates:` target survives `live-links` and is not reported by the check.
- [ ] A broken link whose file name matches exactly one removed history file gets that file's Git Locator; none or several keep the not-reachable form.
- [ ] `roundfix history sanitize` reproduces Surface Transcript 1.

## Context

- instruction: `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- instruction: `internal/cli/history.go`
- instruction: `internal/spec/declared_paths.go`
- creates: `internal/spec/pending_targets.go`
- creates: `internal/spec/pending_targets_test.go`
- interface: `internal/spec/live_links.go`
- interface: `internal/spec/live_links_test.go`
- interface: `internal/speccheck/history_retained.go`
- interface: `internal/speccheck/history_retained_test.go`
- creates: `internal/cli/history_citations_transcript_test.go`

## Consumers

- the `live-links` unit and the missing-link findings of `SC-HISTORY-RETAINED` — consumers: `./internal/cli`, `./internal/speccheck`

## Verification

- `grep -q 'func PendingTaskCreates' internal/spec/pending_targets.go || { printf 'missing PendingTaskCreates\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestLiveLinksLeaveAPendingTaskCreatesTarget|TestLiveLinksResolveABrokenLinkByTheRemovedFileName|TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLiveLinksLeaveAPendingTaskCreatesTarget TestLiveLinksResolveABrokenLinkByTheRemovedFileName TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestHistoryRetainedSkipsAPendingTaskCreatesTarget|TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryRetainedSkipsAPendingTaskCreatesTarget TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestHistorySanitizeRewritesEveryCitationTranscript$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistorySanitizeRewritesEveryCitationTranscript ' || { printf 'missing pass: TestHistorySanitizeRewritesEveryCitationTranscript\n' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/spec ./internal/speccheck ./internal/cli || exit 1; go test -count=1 ./internal/spec ./internal/speccheck || exit 1; go test -count=1 -run 'History|Citation|Baseline|Live|Sanitize' ./internal/cli` — expected: exit 0. Before this Task `PendingTaskCreates` does not exist, so the command fails at its first check. Afterwards the new `internal/spec`, `internal/speccheck` and transcript tests pass, the packages build for Windows, `internal/spec` and `internal/speccheck` pass whole, and the history, citation, baseline, live-links and sanitize tests of `internal/cli` pass.

## References

- `_prd.md` → Sources; Goals; Core Feature 3; Success Metric 3; Acceptance evidence
- `_techspec.md` → API Contract 4; API Contract 5; Surface Transcript 1; Invariants 4-5; Current behavior; Testing Approach; Build Order 3
- ADR-0299; ADR-0278; ADR-0282; ADR-0295; ADR-0259; ADR-0285; ADR-0184

## Result

Implemented this Task's pending-target exemption and unambiguous file-name
resolution. The initial consumers resolved only the missing link's exact path
and had no pending `creates:` exemption.

- `PendingTaskCreates` reads direct Spec directories with active PRDs and uses
  the existing Task parser. Non-completed Tasks contribute only `creates:`
  paths; malformed Tasks contribute nothing. Read errors propagate.
- `RemovedHistoryPaths` reads HEAD's tree and deletion history through Git,
  with `--no-renames` and only the findings, backlog and ADR history families.
  It excludes present working-tree paths, sorts and deduplicates results, and
  returns no paths outside Git or without HEAD.
- `PlanLiveLinks` keeps its signature and CLI callers. It exempts pending
  targets, then falls back from exact-path resolution to one distinct
  same-named candidate from planned removals and removed history. Duplicate
  occurrences of one candidate do not create ambiguity. Citation targets keep
  the original link path; the rewrite names the selected history file.
- The missing-link branch of `SC-HISTORY-RETAINED` shares the exemption and
  file-name lookup. Existing finding wording and Archive Record resolution
  remain unchanged; resolved missing links still produce findings whose fix
  line names the history file's locator.

Acceptance evidence from focused checks:

| Acceptance criterion | Evidence |
| --- | --- |
| A pending target survives conversion and is not reported | `TestLiveLinksLeaveAPendingTaskCreatesTarget` exercises pending/in-progress Tasks, completed Tasks and inactive PRDs. `TestHistoryRetainedSkipsAPendingTaskCreatesTarget` checks pending/completed Tasks both inside and outside Git. `TestPendingTaskCreatesUsesOnlyActiveSpecsAndValidCreates` covers malformed Tasks, inactive Specs and exclusion of interface paths. |
| One matching removed file resolves; zero or several preserve the original fallback | `TestLiveLinksResolveABrokenLinkByTheRemovedFileName` checks planned removals, earlier committed deletions, duplicate candidate paths, ambiguous names and unknown names, including applied bytes and original citation targets. `TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks` checks working-tree deletion, committed deletion, rename, present/outside-family exclusions, outside Git and unborn HEAD. `TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames` checks the retained finding and exact fix line. |
| Surface Transcript 1 is reproduced | `TestHistorySanitizeRewritesEveryCitationTranscript` runs CLI dispatch in an isolated temporary repository with the specified commit date and annotated tag. It asserts every stdout line in order, empty stderr and exit 0; only byte counts and the fixture commit vary. |

Focused commands and observed outcomes:

- `rtk proxy go test ./internal/spec ./internal/speccheck ./internal/cli -run 'Test(PendingTaskCreatesUsesOnly|RemovedHistoryPathsLists|LiveLinks|HistoryRetained|HistorySanitizeRewritesEveryCitationTranscript)' -count=1`
  — exit 0 for all three packages, including existing neighboring live-link
  and retained-history regressions.
- `GOOS=windows go test -c -buildvcs=false -o /private/tmp/roundfix-task03-spec.test.exe ./internal/spec`
  — exit 0.
- `GOOS=windows go test -c -buildvcs=false -o /private/tmp/roundfix-task03-speccheck.test.exe ./internal/speccheck`
  — exit 0.
- `GOOS=windows go test -c -buildvcs=false -o /private/tmp/roundfix-task03-cli.test.exe ./internal/cli`
  — exit 0.
- `rtk proxy git diff --check` — exit 0.

All new top-level tests begin with `t.Parallel()`. Newly authored child-process
helpers register cancellation and reaping in `t.Cleanup`; subprocess fixtures
use private homes. The Run's symlinked worktree required execution through its
resolved checkout with sandbox escalation; no product blocker was observed.

Task status remains Daemon-owned. Declared Verification and repository gates
were not run; the Daemon runs them after this handoff. No Task Graph or other
Task file was edited, and no commit, push or pull request was made. The new
removed-citation finding from API Contract 6 remains outside this Task's slice.
