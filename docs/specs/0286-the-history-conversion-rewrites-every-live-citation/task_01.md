---
task: task_01
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: completed
type: backend
complexity: high
---

# Task 01: The history plan rewrites every citation in every Citation Source

## Overview

On 2026-10-19 an adopter's history conversion listed 45 citations, rewrote
30 and left 18 pointing at removed files, in ADRs, agent guides, READMEs and a
Makefile comment (`_prd.md` → Sources 1, defect 1). This Task makes the plan
find citations in every Citation Source, including prose paths that start
with `./` or `../`, and rewrite every citation it lists. It is verifiable on
its own through `PlanHistoryRewrite` and through `roundfix baseline update` on
temporary repositories.

## Requirements

1. MUST answer defect 1 of the Inbox Entry of 2026-10-19,
   [`references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt`](references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 and 2 of
   `_techspec.md` (ADR-0299): `CitationSources` and `SourceHistoryCitations`
   in a new `internal/spec/citation_sources.go`, `CitationRewrite.Relative`
   and `CitationRewrite.ApplyText` in `internal/spec/git_locator.go`, and
   `PlanHistoryRewrite` listing through `SourceHistoryCitations` with its
   `scope` condition and `historyRootInstructionCarrier` removed, as
   `_techspec.md` → Interfaces sketches them.
2. MUST state the hazard in the code's behavior, as API Contracts 1 and 2 do:
   - the data source is the repository's own tracked files, read only;
   - the plan writes only the files it listed, and apply re-reads only those;
   - a file under `.agents/skills/`, `skills/` or a `testdata` directory, a
     file holding a NUL byte in its first 8,000 bytes and a symbolic link are
     never read for citations and never written;
   - outside Markdown only bounded path text changes, and the Git Locator is
     written as plain `<path>@<commit>`, never with backticks.
3. MUST keep `HistoryCitations` reading only tracked Markdown, so the Archive
   Relink, its warnings and `PlanSupersededRecords` keep their behavior (API
   Contract 7), and keep `Relative` unset everywhere except
   `PlanHistoryRewrite`, so `RelinkContent`, `rewriteRecordLinks` and the
   Legacy Link Rebase keep their behavior.
4. MUST declare these breaks, each test keeping its name, as `_techspec.md`
   → Current behavior measured them:
   - `TestRemovalRewriteLocatesAReducedEntryAtItsProvenance`: the `README.md`
     citation is rewritten too, so the plan holds two files and count 2;
   - `TestHistorySanitizePlanWritesNothing`: the `README.md` citation line
     gains `; rewrites to <locator>`;
   - `TestBaselineUpdateHistoryCitationBoundaries`,
     `TestBaselineUpdateListsHistoryCitationsTranscript` and
     `TestBaselineUpdatePromotesOnlyOnRequestTranscript`: the handoff's
     citation of the converted folder's `_prd.md` carries `rewrite` to the
     Archive Record path;
   - `TestBaselineUpdateReportsStaleHistoryCitationsTranscript`: its stale
     citation is produced by deleting the cited file outside the conversion,
     so the case still reports one stale citation.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and ending every child it starts in `t.Cleanup`:
   - in `internal/spec/citation_sources_test.go`:
     - `TestHistoryRewriteReachesEveryCitationSource`: on the adopter's shape
       of Surface Transcript 1 (an ADR with a path between backticks and a
       relative link, an agent guide with a relative path between backticks,
       a reference page link, a README with a link and a relative path between
       backticks, a Makefile comment), the plan lists every citation, rewrites
       every one, and after `Apply` no Citation Source still names a removed
       path outside a Git Locator;
     - `TestHistoryRewriteResolvesARelativeProsePath`: `./` and `../` tokens,
       inside and outside code spans, resolve against the citing file's
       directory; an unprefixed `findings/x.md` is neither listed nor
       rewritten;
     - `TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown`: a Makefile and
       a `.go` file get `<path>@<commit>` with no backticks, the citation's
       `Rewrite` holds the same plain text, and no link syntax is rewritten
       in them;
     - `TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles`: files under
       `.agents/skills/`, `skills/` and a nested `testdata`, a file with a NUL
       byte and a tracked symbolic link are not listed, and with
       `withUntracked` an untracked file is listed while an ignored one is
       not;
   - in `internal/cli/baseline_update_citations_test.go`,
     `TestBaselineUpdateRewritesEveryListedCitation`: `roundfix baseline
     update --yes` on that shape applies, every listed citation carries
     `rewrite`, `History applied` counts every rewritten citation, and no file
     outside `docs/history` still names a removed path without a Git Locator.

## Subtasks

- [ ] List Citation Sources and their citations, including `./` and `../` prose paths.
- [ ] Rewrite every listed citation; add the plain form outside Markdown.
- [ ] Update the five declared tests.
- [ ] Add the four `internal/spec` tests and the `baseline update` test.

## Acceptance Criteria

- [ ] On the adopter's shape the plan lists and rewrites the ADR, agent guide, reference page, README and Makefile citations, relative paths between backticks included.
- [ ] Outside Markdown the rewrite writes a plain `<path>@<commit>`; skill trees, `testdata`, binary files and symbolic links stay untouched.
- [ ] The Archive Relink, `record-links` and the Legacy Link Rebase keep their tests green.

## Context

- instruction: `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `internal/cli/history.go`
- instruction: `internal/cli/baseline_update_history.go`
- instruction: `internal/spec/live_citations.go`
- creates: `internal/spec/citation_sources.go`
- creates: `internal/spec/citation_sources_test.go`
- interface: `internal/spec/history_entries.go`
- interface: `internal/spec/git_locator.go`
- interface: `internal/spec/history_removal_test.go`
- interface: `internal/cli/history_test.go`
- interface: `internal/cli/baseline_update_citations_test.go`

## Consumers

- the History Citations and rewrite plan `PlanHistoryRewrite` returns — consumers: `./internal/cli`

## Verification

- `grep -q 'func CitationSources' internal/spec/citation_sources.go || { printf 'missing CitationSources\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestHistoryRewriteReachesEveryCitationSource|TestHistoryRewriteResolvesARelativeProsePath|TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown|TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles|TestRemovalRewriteLocatesAReducedEntryAtItsProvenance)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryRewriteReachesEveryCitationSource TestHistoryRewriteResolvesARelativeProsePath TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles TestRemovalRewriteLocatesAReducedEntryAtItsProvenance; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestBaselineUpdateRewritesEveryListedCitation|TestHistorySanitizePlanWritesNothing|TestBaselineUpdateHistoryCitationBoundaries|TestBaselineUpdateListsHistoryCitationsTranscript|TestBaselineUpdatePromotesOnlyOnRequestTranscript|TestBaselineUpdateReportsStaleHistoryCitationsTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateRewritesEveryListedCitation TestHistorySanitizePlanWritesNothing TestBaselineUpdateHistoryCitationBoundaries TestBaselineUpdateListsHistoryCitationsTranscript TestBaselineUpdatePromotesOnlyOnRequestTranscript TestBaselineUpdateReportsStaleHistoryCitationsTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/spec ./internal/cli || exit 1; go test -count=1 ./internal/spec || exit 1; go test -count=1 -run 'History|Citation|Baseline|Archive|Sanitize|Relink' ./internal/cli` — expected: exit 0. Before this Task `internal/spec/citation_sources.go` does not exist, so the command fails at its first check. Afterwards the four new `internal/spec` tests, the new `baseline update` test and the five updated tests pass, the packages build for Windows, `internal/spec` passes whole, and the history, citation, baseline, archive and relink tests of `internal/cli` pass.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Success Metric 1; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 7; Invariants 1-2; Current behavior; Testing Approach; Build Order 1
- ADR-0299; ADR-0278; ADR-0260; ADR-0254; ADR-0259; ADR-0271; ADR-0285

## Result

Implemented the Task 01 slice; Task status and declared Verification remain
Daemon-owned. The only pre-existing changed path was this Task file (its
`in_progress` status); that status was preserved.

- Citation Sources now come from the repository's tracked regular text files,
  sorted by path, excluding both skill trees, every `testdata` directory,
  symbolic links (including index symlinks on Windows), and files with a NUL
  in their first 8,000 bytes. Untracked files are opt-in and ignored files
  remain excluded. Binary detection reads only that prefix.
- The history plan uses `SourceHistoryCitations`, removes the former source
  scope restriction, resolves explicit `./` and `../` prose paths, and marks
  every listed citation with its replacement. Markdown uses Markdown locators;
  other text uses plain locators and preserves surrounding link syntax.
  Apply re-reads only planned sources, preserves Baseline edits, and refuses
  a source that became a symlink or binary before writing.
- `HistoryCitations` retains its tracked-Markdown behavior. Only the history
  rewrite plan sets `CitationRewrite.Relative`. Empty target lists do not
  scan sources, preserving repeated-apply and stale-detector behavior.
- Updated every declared existing-test break and added the four spec tests
  and CLI test. New isolated test children have cancellation/reaping cleanup;
  CLI subprocess contexts also register cleanup.

Acceptance evidence:

1. The adopter-shaped spec and public CLI regressions list and rewrite seven
   citations across five files (ADR, agent guide, reference page, README and
   Makefile), then find no remaining citation of the removed path outside a
   Git Locator. Relative-prose coverage also checks relative moves and keeps
   an unprefixed `findings/x.md` unchanged.
2. Makefile and Go-text assertions check exact plain locator bytes and retained
   link syntax. Exclusion coverage checks a NUL at byte 7,999, the skill and
   fixture paths, a tracked symlink, untracked/ignored selection, and unchanged
   excluded bytes after Apply.
3. Focused Archive Relink, record-links and Legacy Link Rebase regressions
   passed; the shared legacy detector remains unchanged.

Focused checks (all commands below exited 0 unless explicitly labeled red):

- Initial red: `rtk proxy go test ./internal/spec -run '^TestHistoryRewriteResolvesARelativeProsePath$' -count=1`
  failed to compile because the two new Citation Source APIs did not exist.
- Final task cases: `rtk proxy go test ./internal/spec ./internal/cli -run '^(TestHistoryRewrite.*|TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles|TestRemovalRewriteLocatesAReducedEntryAtItsProvenance|TestBaselineUpdateRewritesEveryListedCitation|TestHistorySanitizePlanWritesNothing|TestBaselineUpdateHistoryCitationBoundaries|TestBaselineUpdateListsHistoryCitationsTranscript|TestBaselineUpdatePromotesOnlyOnRequestTranscript|TestBaselineUpdateReportsStaleHistoryCitationsTranscript|TestHistorySanitizeAppliesTheAdopterShapes|TestBaselineUpdateHistoryStaleDetectorFailureKeepsOutcome)$' -count=1`
  passed in both packages.
- Compatibility cases: `rtk proxy go test ./internal/spec ./internal/cli -run '^(TestRelinkContentRewritesOnlyDestinationsAndPaths|TestPlanArchiveRelink.*|TestCitationRewrite.*|TestLegacyLinkRebase.*|TestLegacyConversionRebasesTheOutcomeAndPromotedCopies|TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource|TestHistorySanitizeRewritesBrokenRecordLinks|TestBaselineUpdateRewritesBrokenRecordLinks)$' -count=1`
  passed in both packages.
- An additional focused history selection initially exposed unnecessary
  scans with no removal targets in the repeated-adopter and stale-detector
  tests. The empty-target return fixes the cause; both tests passed in the
  final focused task run above. No assertions in those tests were changed.
- `GOOS=windows rtk proxy go test -c -o /tmp/roundfix-task01-spec.test.exe ./internal/spec`
  and `GOOS=windows rtk proxy go test -c -o /tmp/roundfix-task01-cli.test.exe ./internal/cli`
  both cross-compiled the final packages and tests successfully.
- `rtk proxy git diff --check` passed.

Environment: the sandbox rejected the symlinked writable root before starting
any process. Authorized escalated execution used the resolved worktree under
`/Volumes/Brain-SSD/cache/.roundfix/worktrees/`.

Not run: the Task's declared Verification, selected repository Verification,
selected incremental Verification, or full test suites; these remain for the
Daemon. No other Task or Task Graph file was edited. No commit, push or pull
request was made. No follow-up outside this slice was implemented.
