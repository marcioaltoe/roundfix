---
task: task_01
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
status: completed
type: backend
complexity: medium
---

# Task 01: The strict check names Unreachable Acceptance text the QA gate cannot read

## Overview

A PRD whose `## Unreachable Acceptance` section holds an entry the
declaration reader does not read passes `roundfix spec check --strict` and
fails only at `roundfix qa-report accept`, as Spec 0364's first QA gate did
on 2026-11-06. This Task answers the adopted Backlog Entry of 2026-11-07,
"The strict check names an Unreachable Acceptance entry it cannot read"
([`references/2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md`](references/2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md)):
the reader in `internal/spec` names every line it does not read, and the new
`SC-UNREACHABLE-UNREAD` error reports them, from the PRD stage, with the
fielded form in its fix.

## Requirements

1. MUST, in `internal/spec/spec.go`, make `parseUnreachableDeclarations`
   and a new exported function share one unexported walk of the section:

   ```go
   type UnreachableUnreadLine struct {
   	Line int    // 1-based line number in the PRD
   	Text string // the line, surrounding whitespace trimmed
   }

   func UnreachableUnread(content []byte, prdPath string) ([]UnreachableUnreadLine, error)
   ```

   `UnreachableUnread` returns, sorted by line, every non-empty line of the
   section that no declaration reads, and the same error `Unreachable`
   returns for that content (an `UnreachableDeclarationError` for a
   declaration without a field). `Unreachable` and
   `parseUnreachableDeclarations` MUST return, for every input, the same
   declarations and error as before this Task; this Task MUST NOT make
   either refuse unread text.
2. MUST apply this reading rule, and no other, inside the section (which
   still ends at the next `# ` or `## ` heading). A line is read when it is:
   blank; inside an HTML comment, from a line whose trimmed text starts with
   `<!--` through the line that holds `-->`; a top-level `- ` line, which
   opens a declaration; a `criterion:`, `reason:` or `satisfied-by:` line
   inside an open declaration; or a line indented by a space or a tab that
   continues the field the declaration read last. Before the first
   declaration, the paragraph (up to the next blank line) whose first line
   is the section's first non-blank, non-comment line and begins with
   `None.` is read only when the section holds no declaration at all. Every
   other non-empty line is unread, including an indented line after a
   top-level `- ` line that read no field.
3. MUST create `internal/speccheck/unreachable.go` with
   `const CodeUnreachableUnread = "SC-UNREACHABLE-UNREAD"` and
   `detectUnreachableDeclarations(result *Result, repoRoot, prdPath string) error`,
   which reads the PRD and adds, in this order:
   - one finding per run of consecutive unread line numbers, at the run's
     first line, `Severity: SeverityError`,
     `Where: []Location{{Path: <prd display path>, Line: <N>}}`, and the
     summary
     `<prd display path> line <N> under ## Unreachable Acceptance is not part of a fielded declaration, so the QA gate does not read it: <first line quoted with strconv.Quote>`;
   - when `UnreachableUnread` returns an `UnreachableDeclarationError`, one
     finding at its line with the summary
     `<prd display path> line <N> starts an Unreachable Acceptance declaration without its <field>: field`.

   Both carry exactly this fix:
   `Write each entry in the fielded form: a top-level "- criterion: <text the QA row repeats>" line, then indented "reason: <why no Verification or QA row can reach it>" and "satisfied-by: <evidence or follow-up that satisfies it>" lines; start the section with None. when no criterion is unreachable.`
   An absent section adds no finding and no skip.
4. MUST call `detectUnreachableDeclarations` from `Check` in
   `internal/speccheck/constraints.go` once the PRD is present, and from
   `checkAuthoringStage` in `internal/speccheck/coherence.go` once the PRD is
   present; where either adds its skips for a missing PRD, add
   `addSkip(..., CodeUnreachableUnread, <prd display path>)`. MUST add
   `{code: CodeUnreachableUnread, stage: StagePRD}` to `stagedDetectors`.
   Update every test that pins the staged code list or the skip list on the
   default branch when the Task runs, if any; at `2e59aebe` none does.
5. MUST create `internal/spec/unreachable_unread_test.go` with
   `TestUnreachableUnreadReadsTheSectionAsUnreachableDoes`: over every
   fixture under `internal/spec/testdata/unreachable/`, `UnreachableUnread`
   returns the error `Unreachable` returns and no unread line, and
   `parseUnreachableDeclarations` returns the declarations `Unreachable`
   returns; on a section holding a numbered two-line entry, a fielded
   declaration and an unindented prose line after it, the unread lines are
   exactly the numbered entry's two lines and the prose line.
6. MUST create `internal/speccheck/unreachable_test.go`, whose fixtures are
   PRDs written into `t.TempDir()`, with:
   - `TestAnUnreachableSectionReadAsWrittenReportsNothing`: a fielded section
     (one declaration with one criterion, one with two `criterion:` lines and
     a wrapped criterion), a `None.` paragraph of two lines, an empty
     section, an absent section, and the `write-prd` template's HTML comment
     left in place each report no `SC-UNREACHABLE-UNREAD`;
   - `TestANumberedUnreachableEntryIsReportedAtItsFirstLine`: Spec 0364's
     shape, `1. A real rerun of a failed release job on GitHub Actions against the npm`
     and two indented lines, reports exactly one error, at the `1.` line,
     whose summary holds that line quoted and whose fix holds
     `"- criterion: <text the QA row repeats>"`, both through `Check` and
     through `CheckStage` at `StagePRD`;
   - `TestProseAfterADeclarationAndAMissingFieldAreReported`: an unindented
     prose line after a declaration's `reason:` is reported at its line, and
     a declaration without `satisfied-by:` is reported with
     `without its satisfied-by: field`;
   - `TestNoneBesideADeclarationAndAPlainBulletAreReported`: a `None.`
     paragraph followed by a fielded declaration reports the `None.` line,
     and `- A real rerun of a failed release job.` with an indented second
     line reports the indented line as unread and the bullet with
     `without its criterion: field`.

   Every new top-level test calls `t.Parallel()` as its first statement.
7. MUST, in `CONTEXT.md`, append to the definition line of
   **Unreachable Acceptance**, after "the declaration count remains the
   fallback for reports accepted before this rule (ADR-0343).", a space and
   this sentence, exactly:

   ```markdown
   Each declaration is a top-level `- criterion:` line with its `reason:` and `satisfied-by:` lines; the Spec Consistency Check reports any other text in the section as `SC-UNREACHABLE-UNREAD`, except a first paragraph that begins with `None.` in a section that holds no declaration.
   ```

   Write it through the `domain-modeling` skill; no other line of
   `CONTEXT.md` changes.
8. MUST NOT change `internal/spec/qa*.go`, `internal/spec/archive.go`, any
   guide, skill, ADR or record; task_03 documents the code.

## Subtasks

- [ ] Share one walk between the declaration reader and `UnreachableUnread`.
- [ ] Add the detector, its two call sites, its skips and its stage.
- [ ] Add the two test files.
- [ ] Add the glossary sentence.

## Acceptance Criteria

- [ ] Spec 0364's numbered entry is reported once, at its first line, with the fielded form in the fix.
- [ ] A fielded, `None.`, empty, absent or commented section reports nothing, and `Unreachable` returns what it returned before.
- [ ] The glossary states the reading rule.

## Context

- interface: `internal/spec/spec.go`
- creates: `internal/spec/unreachable_unread_test.go`
- creates: `internal/speccheck/unreachable.go`
- creates: `internal/speccheck/unreachable_test.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- interface: `CONTEXT.md`

## Consumers

- `spec.Unreachable` and the new `spec.UnreachableUnread` — consumers: `./internal/spec`, `./internal/speccheck`
- The `SC-UNREACHABLE-UNREAD` finding of `speccheck.Check` and `speccheck.CheckStage` — consumers: `./internal/speccheck`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestUnreachableUnreadReadsTheSectionAsUnreachableDoes$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestUnreachableUnreadReadsTheSectionAsUnreachableDoes ' || { echo 'missing pass: TestUnreachableUnreadReadsTheSectionAsUnreachableDoes' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/spec` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes with the whole package, every existing `Unreachable` test included.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnUnreachableSectionReadAsWrittenReportsNothing|TestANumberedUnreachableEntryIsReportedAtItsFirstLine|TestProseAfterADeclarationAndAMissingFieldAreReported|TestNoneBesideADeclarationAndAPlainBulletAreReported)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnUnreachableSectionReadAsWrittenReportsNothing TestANumberedUnreachableEntryIsReportedAtItsFirstLine TestProseAfterADeclarationAndAMissingFieldAreReported TestNoneBesideADeclarationAndAPlainBulletAreReported; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/speccheck` — expected: exit 0. Before this Task none of the four tests exists; afterwards they pass with the whole Spec Check suite.
- `grep -F -- '**Unreachable Acceptance**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Unreachable Acceptance\*\*(:| —( |$))' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'the Spec Consistency Check reports any other text in the section as' || exit 1; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce)$'` — expected: exit 0. Before this Task the glossary lacks the sentence, so the second check fails; afterwards the term keeps its definition line, the sentence is in place and the glossary contracts pass.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; API Contract 1; Testing Approach; Build Order 1
- ADR-0094; ADR-0343

## Result

Skills loaded: implement-task, coding-guidelines, domain-modeling, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, no-workarounds, systematic-debugging, tech-writer, evidence-gate.

The shared section walk now returns unread lines beside the existing declarations. It retains the first declaration error while continuing to collect unread lines. Declaration parsing preserves its historical comment and two-space continuation behavior; the unread-line state applies the specified comment, one-space/tab continuation and initial `None.` paragraph rules without changing QA consumers' declarations.

`SC-UNREACHABLE-UNREAD` runs immediately after PRD presence is established in both check entry points, is registered at `StagePRD`, and records a missing-PRD skip. It reports consecutive unread lines once per run, followed by the first missing-field error, with the exact summaries and fix required above.

Acceptance evidence from focused checks:

- Numbered entry: `TestANumberedUnreachableEntryIsReportedAtItsFirstLine` asserts exactly one error at the numbered line, the quoted first line, the complete fix and its fielded-form example through `Check`, `StagePRD` and `StageTechSpec`. `TestUnreadRunsKeepTheirFirstLines` checks separate unread runs; the prose/missing-field and plain-bullet tests check finding order and locations.
- Readable sections and compatibility: all five existing unreachable fixtures are checked for declaration/error parity and zero unread lines. The readable-section test covers fielded declarations, repeated and wrapped criteria, a two-line `None.` paragraph, empty/absent sections and the literal write-prd template comment. Boundary tests cover comments, one-space/tab continuation, orphan fields/continuations, the paragraph's blank-line boundary, section headings, and unread text after a malformed declaration. The missing-PRD test checks exactly one path-specific skip at all three stages and no finding; present PRDs get no such skip.
- Glossary: the exact required sentence was appended through `domain-modeling`. A Python comparison against `git show HEAD:CONTEXT.md` confirmed that removing that single suffix reproduces the original bytes; no other glossary content changed.

Focused commands and outcomes:

- Before implementation, `rtk proxy go test ./internal/spec ./internal/speccheck -run 'Test(UnreachableUnread|ANumberedUnreachable)' -count=1 -timeout=2m` failed to compile because `UnreachableUnread`, `UnreachableUnreadLine` and `CodeUnreachableUnread` did not exist. The first sandbox attempt could not access Go's existing build cache; the authorized retry reached these expected compiler diagnostics.
- Final focused run: `rtk proxy go test -overlay /private/tmp/roundfix-unreachable-wm9784gt/overlay.json ./internal/spec ./internal/speccheck -run 'Test(Unreachable|AnUnreachable|ANumberedUnreachable|ProseAfterADeclaration|NoneBesideADeclaration|UnreadRuns|StageScope)' -count=1 -timeout=3m` exited 0 for both packages. The temporary overlay adds a compatibility probe without changing repository files: it compares the original parser and builder from `HEAD` against the shared walk on 10,000 deterministic inputs and checks identical declarations and errors, including `UnreachableUnread` errors. The regular focused tests also passed before adding the overlay.
- `gofmt` was applied to the changed Go files; `rtk proxy git diff --check` exited 0.

Declared Verification commands, repository Verification and the full package suites were not run; they remain Daemon-owned. The Task status, Task Graph, other Task files, guides, skills, ADRs, records and QA/archive consumers were not edited. No commit, push or pull request was made. No follow-up outside this slice was needed.
