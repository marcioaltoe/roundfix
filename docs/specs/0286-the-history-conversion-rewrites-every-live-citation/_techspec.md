---
spec: 0286-the-history-conversion-rewrites-every-live-citation
prd: _prd.md
created: 2026-10-20
---

# The history conversion rewrites every live citation — Technical Spec

## Executive Summary

Four slices correct the history conversion that `roundfix baseline update`
and `roundfix history sanitize` share, and extend the check that guards it.
No package, command, flag, configuration key or JSON field is new.

1. The plan finds History Citations in every Citation Source and rewrites
   every one it lists (`internal/spec`).
2. The lineage line names a retired ADR by the identifier its repository uses
   and dates it with the day it was superseded (`internal/spec`).
3. `live-links` and `SC-HISTORY-RETAINED` leave a link to a pending Task's
   `creates:` target, and resolve a missing link by its file name
   (`internal/spec`, `internal/speccheck`).
4. `SC-HISTORY-RETAINED` reports a citation of a Removed History Path that is
   not a Git Locator, and the guides, references and glossary describe the
   change (`internal/speccheck`, documentation).

The trade-off accepted: the conversion now reads and may rewrite any tracked
text file, a Makefile or a Go file included, instead of a fixed set of
Markdown families. Every rewrite is listed first and bound into the Plan
Digest, the skill trees and `testdata` stay out, and a citation outside
Markdown changes only its path text. A broader write set is the price of
leaving no citation behind, which is what the adopter's grep demanded.

## Project Constraints

- Identifier strategy: applicable — the lineage line identifies a retired
  ADR as `ADR-NNNN` for a numbered ADR and as its file name without `.md`
  otherwise (API Contract 3). No identifier, slug or record schema is minted.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the conversion and the check read the
  repository and its Git history only, and every test uses temporary
  repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Build Order follows recorded
  decisions. ADR-0299 governs Build Order 1–4: "The history conversion
  rewrites every citation it lists, and the check reports a removed history
  file cited without a Git Locator". ADR-0278, which ADR-0299 extends: "A
  retirement deletes, and the History Root holds only Archive Records".
  ADR-0282, which ADR-0299 amends: "The history cleanup names a superseded
  ADR in its successor, and a large candidate is reviewed in segments".
  ADR-0260: "An archive relinks live citations, and a Baseline update lists
  the citations it breaks". ADR-0254: "Baseline update sanitizes pending
  history in the change it plans". ADR-0295: "Reachability reads the fetched
  default branch, and origin names it when origin/HEAD is unset". ADR-0259,
  ADR-0271, ADR-0285, ADR-0189, ADR-0256, ADR-0280, ADR-0104, ADR-0240,
  ADR-0276, ADR-0265, ADR-0183 and ADR-0184, each named with its title in
  `_prd.md` → Project Constraints. Set no obligation here, because this Spec
  changes none of their rules: ADR-0092, ADR-0120, ADR-0165, ADR-0167,
  ADR-0169, ADR-0196, ADR-0197, ADR-0229, ADR-0233, ADR-0237, ADR-0247,
  ADR-0248, ADR-0250, ADR-0251, ADR-0258, ADR-0261, ADR-0262, ADR-0269,
  ADR-0272, ADR-0277, ADR-0279, ADR-0281, ADR-0286 and ADR-0290, each named
  with its title in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 ("Sim, mesmas regras do AJ (Recomendado)"),
  carrying the cycle AJ grant "Mexer em Governed Paths declarados", and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0286-the-history-conversion-rewrites-every-live-citation/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `9448aa64` with `bin/roundfix` built by `make
build`, on a temporary repository of the adopter's shape (the fixture of
Surface Transcript 1, before the change), `HOME` set to a temporary
directory.

- Listing. `HistoryCitations` (`internal/spec/history_entries.go:475`) reads
  only `git ls-files -- '*.md'`. Its detector, `historyEntryCitesPath`
  (`:534`), matches a repository-relative path and resolves inline link
  destinations; a prose path such as `` `../history/findings/x.md` `` is
  never resolved. The plan listed the ADR, reference page and README
  citations, a skill file and a `testdata` file, and did not list the
  Makefile comment or the two relative paths between backticks.
- Rewriting. `PlanHistoryRewrite` (`:687`) rewrites a listed citation only
  when `scope` (`:761`) holds: a root instruction carrier, a Live Document, a
  `docs/user-guide/` page or an Archive Record. After `history sanitize
  --apply --batch 10`, `git grep -n -E 'history/(findings|backlog|adr)/'`
  without `@<commit>` still printed `Makefile:2`, `docs/adr/0002-live.md:6`,
  `docs/agents/deploy.md:1`, `docs/reference/notes.md:1` and
  `infra/README.md:1`. `TestRemovalRewriteLocatesAReducedEntryAtItsProvenance`
  (`internal/spec/history_removal_test.go:206`) pins that a `README.md`
  citation stays unchanged.
- Lineage. `PlanHistoryKind` builds `id := "ADR-" + strings.SplitN(name,
  "-", 2)[0]` (`:134`) and dates the line with `HEAD`'s committer date
  (`:136-148`). The successor of `docs/adr/2026-09-18-old-cluster.md`
  (`updated_at: 2026-10-01T00:00:00Z`, `superseded_by:
  2026-10-01-new-cluster`) received `Supersedes ADR-` followed by `2026 (2026-10-19);
  Git keeps its text at ...`, 2026-10-19 being the fixture commit's date.
  `historyADRReferenceNames` (`:292`) treats the leading `2026` of a date
  name as a number, so `superseded_by: 2026` would select the first ADR of
  2026.
- Pending Task targets. `PlanLiveLinks` (`internal/spec/live_links.go:15`)
  rewrote `[m](references/medicao.md)`, a file `task_01` (status `pending`)
  declares in `creates:`, to `` m (`docs/specs/0001-demo/references/medicao.md`
  (not reachable from any branch or tag)) ``. `SC-HISTORY-RETAINED`
  (`internal/speccheck/history_retained.go:194`) reports the same link as
  missing, through the same `MissingRelativeLinks`.
- Broken links. `[f](../findings/2026-01-01-x.md)` in
  `docs/specs/0001-demo/references/note.md`, whose file lived at
  `docs/history/findings/2026-01-01-x.md`, became `` `docs/specs/0001-demo/findings/2026-01-01-x.md`
  (not reachable from any branch or tag) ``: `ResolveGitLocator` looks only at
  the path as written.
- The check. With the retirement rule in the guide, `bin/roundfix spec check
  --strict` on the applied tree reported no `SC-HISTORY-RETAINED` finding for
  any of the five citations left behind: the detector walks only `docs` and
  the Spec Root, reads only Markdown, and checks only relative links of Live
  Documents, Archive Records and the user guide.
- Prototype. Widening `scope` alone in a scratch copy failed exactly these
  tests, each pinning the narrow set:
  `TestRemovalRewriteLocatesAReducedEntryAtItsProvenance`,
  `TestHistorySanitizePlanWritesNothing` (`internal/cli/history_test.go`),
  `TestBaselineUpdateHistoryCitationBoundaries`,
  `TestBaselineUpdateListsHistoryCitationsTranscript`,
  `TestBaselineUpdatePromotesOnlyOnRequestTranscript` and
  `TestBaselineUpdateReportsStaleHistoryCitationsTranscript`
  (`internal/cli/baseline_update_citations_test.go`). Every other
  `internal/spec`, `internal/speccheck` and `internal/baseline` test passed.
- The roundfix repository itself holds no citation of a removed history file
  outside `testdata` that is not already a Git Locator: of 33 distinct
  history-family paths cited, 4 were deleted from `HEAD`'s history, two of
  them written as Git Locators and two under `testdata` behind a `mirror/`
  prefix.

## System Architecture

`internal/spec` owns the conversion; `internal/speccheck` owns the check;
`internal/cli` calls both unchanged.

- `internal/spec/citation_sources.go` (new): `CitationSources` and
  `SourceHistoryCitations`.
- `internal/spec/pending_targets.go` (new): `RemovedHistoryPaths`,
  `PendingTaskCreates` and `ResolveMissingLinkByName`.
- `internal/spec/history_entries.go`: `PlanHistoryRewrite` lists through
  `SourceHistoryCitations` and drops `scope`; `PlanHistoryKind` uses
  `historyADRIdentifier` and `historyADRSupersessionDate`.
- `internal/spec/git_locator.go`: `CitationRewrite` gains `Relative` and
  `ApplyText`.
- `internal/spec/live_links.go`: the creates exemption and the file-name
  resolution.
- `internal/speccheck/history_retained.go`: the same exemption and
  resolution for missing links, and the removed-citation finding.

`HistoryCitations`, the Archive Relink, `record-links`, Stale History
Citations and `PlanSupersededRecords` keep their sources and behavior.

## Implementation Design

### Interfaces

```go
// internal/spec
func CitationSources(ctx context.Context, root string, withUntracked bool) ([]string, error)
func SourceHistoryCitations(ctx context.Context, root string, removed []string) ([]HistoryCitation, error)
func RemovedHistoryPaths(ctx context.Context, root string) ([]string, error)
func PendingTaskCreates(specsRoot string) (map[string]bool, error)
func ResolveMissingLinkByName(target string, candidates []string) (string, bool)

type CitationRewrite struct {
	Moves    map[string]string
	Locate   func(path string) (GitLocator, bool, error)
	Relative bool // resolve ./ and ../ prose paths; only the history plan sets it
}
func (rw CitationRewrite) ApplyText(content []byte, path string) ([]byte, int, error)

func historyADRIdentifier(name string) string
func historyADRSupersessionDate(ctx context.Context, root, rel string, content []byte) (string, error)
```

### Data Models

None. No file, record, JSON field or configuration key changes shape.

### API Contracts

1. API Contract: Citation Sources and their citations.
   - `CitationSources` lists `git ls-files -z` (plus `--others
     --exclude-standard` when `withUntracked`), keeps a regular file that
     holds no NUL byte in its first 8,000 bytes, and drops a symbolic link and
     every path under `.agents/skills/`, `skills/` or a directory named
     `testdata`. Paths are repository-relative and sorted.
   - `SourceHistoryCitations` reads each Citation Source not removed by the
     same plan. It finds what `historyEntryCitesPath` and the archive scanner
     find today, plus a prose token that starts with `./` or `../`, between
     the boundaries `citationPathBoundary` uses, resolved against the citing
     file's directory. Citations come in source-path order, then line, then
     target, with the caller's target spelling.
   - Data source: the repository's own tracked files. Destination: the plan's
     text and JSON, and the Plan Digest. Nothing leaves the machine.
2. API Contract: every listed citation is rewritten.
   - `PlanHistoryRewrite` lists through `SourceHistoryCitations` and rewrites
     every citation it lists, except in a source the same plan removes. The
     `scope` condition is removed. `historyRootInstructionCarrier` goes with
     it, because no caller remains.
   - A Markdown source (`.md`, `.markdown`, `.mdx`) goes through `Apply` with
     `Relative: true`: a link becomes its label and the Git Locator in code,
     and a prose path becomes the locator, without a second pair of
     delimiters inside code. A moved path becomes its destination, as today.
   - Any other source goes through `ApplyText`: no link scanning; a bounded
     path or `./`/`../` token naming a removed path becomes `<path>@<commit>`,
     `<path> (not reachable from any branch or tag)` or `<path> (pull request
     #N)`, with no backticks. A moved path becomes its destination.
   - The citation's `Rewrite` field holds `GitLocator.Text(false)` for a
     Markdown source and `GitLocator.Text(true)` for any other source.
   - The plan never writes a file it did not list, and apply re-reads only
     the planned files, as today.
3. API Contract: the lineage line.
   - `historyADRIdentifier(name)`: when `name` is `<digits>-<slug>.md` and
     does not start `YYYY-MM-DD-`, it returns `ADR-<digits>`; otherwise the
     file name without `.md`.
   - `historyADRSupersessionDate`: the retired ADR's front matter
     `updated_at`, a YAML date or timestamp converted to UTC, formatted
     `2006-01-02`; when absent or unparsable, the committer date in UTC of
     `git log -1 --format=%cI HEAD -- <path>`; when that is empty, `HEAD`'s
     committer date in UTC.
   - The line stays `Supersedes <identifier> (<date>); Git keeps its text at
     <locator>.`.
   - `historyADRReferenceNames` accepts `ADR-NNNN` and `NNNN` only for a
     numbered name; the stem, the file name and a path ending in the file
     name match any name. `historyADRNamesPredecessor` looks for the
     identifier and the file name.
4. API Contract: a pending Task's target is not missing.
   - `PendingTaskCreates(specsRoot)` reads each Spec directory whose `_prd.md`
     has `status: active` and each `task_*.md` in it, through the parser
     `DeclaredTaskPaths` uses. A Task whose status is not `completed`
     contributes its `creates:` paths. A Task file that does not parse
     contributes nothing; a directory that cannot be read returns an error.
   - `PlanLiveLinks` computes it from `root` and `specRootRel` and skips a
     missing link whose target it holds. `SC-HISTORY-RETAINED` skips the same
     link. The function's signature and its callers do not change.
5. API Contract: a broken link is resolved by its file name.
   - `RemovedHistoryPaths` returns the union of `git ls-tree -r --name-only
     -z HEAD` and `git log --format= --name-only -z --diff-filter=D
     --no-renames HEAD`, both limited to `docs/history/findings`,
     `docs/history/backlog` and `docs/history/adr`, minus paths present in
     the working tree, sorted and deduplicated. Outside Git, or with no
     `HEAD`, it returns none.
   - When `ResolveGitLocator` gives no commit for a missing link's target,
     `PlanLiveLinks` calls `ResolveMissingLinkByName(target, candidates)`.
     The candidates are the plan's removed paths (its `skip` keys) and
     `RemovedHistoryPaths`. Exactly one candidate with the target's base name
     gives `ResolveGitLocator(candidate, "history-full", "")`. None or
     several keep the path as written. The citation's `Target` stays the
     link's path, and its `Rewrite` is the chosen locator.
   - The missing-link diagnostic of `SC-HISTORY-RETAINED` prints the same
     locator in its fix line.
6. API Contract: `SC-HISTORY-RETAINED` reports a removed citation.
   - With the retirement rule present and inside Git, the detector reads
     `CitationSources(ctx, root, true)` outside `docs/history/` and looks
     for citations of `RemovedHistoryPaths` with `SourceHistoryCitations`'
     detector. Text already written as a Git Locator, or followed by `(not
     reachable from any branch or tag)` or `(pull request #`, is not a
     citation.
   - Each remaining citation is an error: summary `<path> cites a removed
     history file without a Git Locator: <target>`, location `<path>:<line>`,
     fix `Replace the citation with <locator>.`. The locator is
     `ResolveGitLocator(target, "", "")`, written `Text(false)` in Markdown and
     `Text(true)` elsewhere.
   - Exempt: every file under `docs/history/`, Archive Records included; a
     target `PendingTaskCreates` holds; a citation on a line where the same
     detector run already reports a missing relative link to the same target.
   - Outside Git the finding is skipped under the existing skip reason;
     without the retirement rule the whole detector stays skipped.
7. API Contract: everything else stays. The text and JSON line formats, the
   JSON field set, the Plan Digest's fields, `record-links`, Stale History
   Citations, `superseded` records, the Archive Relink and its warnings,
   `HistoryCitations` and the existing `SC-HISTORY-RETAINED` findings keep
   their behavior. Only the set of rewritten citations, the lineage
   identifier and date, the `live-links` targets and the new finding change.

### Surface Transcripts

The transcript runs in a temporary repository on branch `main`, `HOME` a
temporary directory, committed at `2026-10-19T12:00:00Z` with an annotated
`history-full` tag. It holds:

- `docs/history/findings/2026-01-01-x.md` (`status: done`),
  `docs/history/backlog/2026-01-02-b.md` (`status: done`) and
  `docs/history/adr/0001-old.md` (`status: deprecated`);
- `docs/adr/0002-live.md`, citing `` `docs/history/findings/2026-01-01-x.md` ``
  and `[old](../history/adr/0001-old.md)` on line 6;
- `docs/agents/deploy.md`: `` `../history/findings/2026-01-01-x.md` `` on
  line 1;
- `docs/reference/notes.md`: `[b](../history/backlog/2026-01-02-b.md)` on
  line 1;
- `infra/README.md`: `[b](../docs/history/backlog/2026-01-02-b.md)` and
  `` `../docs/history/findings/2026-01-01-x.md` `` on line 1;
- `Makefile`: `# see docs/history/findings/2026-01-01-x.md` on line 2;
- `.agents/skills/demo/SKILL.md` and `internal/testdata/case.md`, each citing
  `docs/history/findings/2026-01-01-x.md`;
- `docs/adr/2026-09-18-old-cluster.md` (`status: superseded`, `updated_at:
  2026-10-01T00:00:00Z`, `superseded_by: 2026-10-01-new-cluster`) and
  `docs/adr/2026-10-01-new-cluster.md` (`status: accepted`, eight lines);
- the active Spec `docs/specs/0001-demo`: `_prd.md` linking
  `[m](references/medicao.md)` on line 9, `task_01.md` (`status: pending`,
  `type: docs`) declaring `creates:
  docs/specs/0001-demo/references/medicao.md`, and `references/note.md`
  linking `[f](../findings/2026-01-01-x.md)` on line 1.

`<b>` stands for a byte count and `<commit>` for the fixture commit, the only
controlled variations.

1. Surface Transcript: the history plan lists and rewrites every citation.

   ```transcript
   $ roundfix history sanitize
   stdout:
   history sanitize plan: 4 unit(s) pending; 4 file(s) (<b> bytes) leave docs/history
   findings: removes 1 file(s) (<b> bytes)
   backlog: removes 1 file(s) (<b> bytes)
   adr: removes 2 file(s) (<b> bytes)
   live-links: rewrites 1 link(s) in 1 file(s)
   cites Makefile:2 names docs/history/findings/2026-01-01-x.md; rewrites to docs/history/findings/2026-01-01-x.md@<commit>
   cites docs/adr/0002-live.md:6 names docs/history/adr/0001-old.md; rewrites to `docs/history/adr/0001-old.md@<commit>`
   cites docs/adr/0002-live.md:6 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`
   cites docs/agents/deploy.md:1 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`
   cites docs/reference/notes.md:1 names docs/history/backlog/2026-01-02-b.md; rewrites to `docs/history/backlog/2026-01-02-b.md@<commit>`
   cites infra/README.md:1 names docs/history/backlog/2026-01-02-b.md; rewrites to `docs/history/backlog/2026-01-02-b.md@<commit>`
   cites infra/README.md:1 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`
   cites docs/adr/2026-10-01-new-cluster.md:10 names docs/adr/2026-09-18-old-cluster.md; rewrites to Supersedes 2026-09-18-old-cluster (2026-10-01); Git keeps its text at `docs/adr/2026-09-18-old-cluster.md@<commit>`.
   cites docs/specs/0001-demo/references/note.md:1 names docs/specs/0001-demo/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

### Invariants

```text
1. After an applied conversion, no Citation Source cites a path the conversion removed except through a Git Locator, the surviving Archive Record or a promoted copy.
2. A file under a skill tree or a testdata directory, a binary file and a symbolic link are never read for citations nor written by the conversion.
3. ADR-NNNN and NNNN name only an ADR whose file name is a number followed by a slug.
4. A link to a creates: target of a Task that is not completed, in an active Spec, is neither rewritten nor reported.
5. A missing link is resolved by file name only when exactly one candidate has that name.
6. With the retirement rule present, SC-HISTORY-RETAINED reports exactly the citations of Removed History Paths outside docs/history that are not Git Locators and not exempt.
```

## Coverage Map

- Goal 1, Goal 2, Core Feature 1, Success Metric 1, Invariants 1-2 → `CitationSources`, `SourceHistoryCitations`, `PlanHistoryRewrite` without `scope`, `CitationRewrite.Relative` and `ApplyText` (API Contracts 1-2), Surface Transcript 1
- Goal 3, Core Feature 2, Success Metric 2, Invariant 3 → `historyADRIdentifier`, `historyADRSupersessionDate`, `historyADRReferenceNames` (API Contract 3)
- Goal 4, Goal 5, Core Feature 3, Success Metric 3, Invariants 4-5 → `PendingTaskCreates`, `RemovedHistoryPaths`, `ResolveMissingLinkByName` in `PlanLiveLinks` and `SC-HISTORY-RETAINED` (API Contracts 4-5)
- Goal 6, Core Feature 4, Success Metric 4, Invariant 6 → the removed-citation finding of `SC-HISTORY-RETAINED` (API Contract 6)
- Core Feature 5 → the `baseline`, `history` and `spec` guides, the Roundfix skill and its `baseline`, `archive` and `spec` references, and `CONTEXT.md`
- API Contract 7 → the existing history, baseline update, archive and check tests, and the QA scope audit

Deviations from the PRD: None.

## Absence Evidence

1. `CitationSources` and the binary test — command: `grep -rn '8000' internal/spec --include='*.go' --exclude='*_test.go'`; result: no match at `9448aa64`; `HistoryCitations` lists only `git ls-files -- '*.md'`.
2. Resolution of `./` and `../` prose paths — command: `grep -n 'HasPrefix(token, "\.\./")\|HasPrefix(path, "\./")' internal/spec/git_locator.go internal/spec/history_entries.go internal/spec/live_citations.go`; result: no match at `9448aa64`; the only `../` tests reject such a path rather than resolve it.
3. `historyADRSupersessionDate` — command: `grep -rn 'updated_at' internal/spec --include='*.go' --exclude='*_test.go'`; result: only `internal/spec/cancel.go`, which writes a Finding's date; nothing reads an ADR's `updated_at`.
4. `PendingTaskCreates` — command: `grep -n 'ContextKindCreates' internal/spec/live_links.go internal/spec/record_links.go internal/speccheck/history_retained.go`; result: no match at `9448aa64`.
5. `RemovedHistoryPaths` — command: `grep -rn 'diff-filter=D' internal/spec --include='*.go' --exclude='*_test.go'`; result: only `internal/spec/history_sanitize.go`, which finds the commit that deleted a legacy Spec folder's `_prd.md`; nothing lists the deleted history-family files.
6. `ResolveMissingLinkByName` — command: `grep -n 'filepath.Base' internal/spec/live_links.go internal/spec/record_links.go internal/speccheck/history_retained.go`; result: no match at `9448aa64`.
7. The removed-citation finding — command: `grep -n 'ls-files\|without a Git Locator' internal/speccheck/history_retained.go`; result: no match at `9448aa64`; the detector walks `docs` and the Spec Root and reads only Markdown.

## ADR Alternatives

- ADR-0299: this Spec's decision; it rejects reporting without rewriting, Markdown-only sources, rewriting skill trees and `testdata`, resolving unprefixed relative words, a new rule ID, choosing among same-named files, and keeping `HEAD`'s date.
- ADR-0278: extended by ADR-0299 with a dated note; none of its rejected alternatives is taken. Citations still become Git Locators, nothing is reduced or relocated, and the check still prints the locator to use.
- ADR-0282: amended by ADR-0299 with a dated note; its rejected "Refusing the whole `adr` unit until a person edits the successor" stays rejected, and the plan still writes the lineage line itself. Only the identifier and the date change.
- ADR-0260: no rejected alternative is taken; the update still promotes only on `--promote` and never promotes automatically. ADR-0278 already superseded its "never rewrites a citing file", and the Archive Relink is unchanged.
- ADR-0254: no rejected alternative is taken; the history still converts in the change the Baseline Plan plans, bound into one Plan Digest.
- ADR-0295: no rejected alternative is taken; Git Locators still read the Default Branch Head and tags.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows and use no unix-only call.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering skill changes with the `baseline`, `history` and `spec` surfaces.
- ADR-0280: no rejected alternative is taken; the guide Task records its own fingerprint, and the PRD records who decided.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the adopter's run and published sources.
- ADR-0240: no rejected alternative is taken; network-blocked rows are recorded as the policy says.
- ADR-0276: no rejected alternative is taken; offline rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0184: no rejected alternative is taken; the changed `history sanitize` output is stated as Surface Transcript 1.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

Git only, on the host, through the existing read helpers (`ls-files`,
`ls-tree`, `log`, `cat-file`). No network, provider, judge or reviewer.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
runs on temporary repositories with the package's existing Git helpers
(`gitLocatorRepo`, `legacyCommit`, `historyEntryWrite`, `gittest`) and ends any
child it starts in `t.Cleanup`.

- Build Order 1, `internal/spec/citation_sources_test.go` (new):
  `TestHistoryRewriteReachesEveryCitationSource`,
  `TestHistoryRewriteResolvesARelativeProsePath`,
  `TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown`,
  `TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles`. In
  `internal/cli/baseline_update_citations_test.go`, new
  `TestBaselineUpdateRewritesEveryListedCitation`. Declared breaks, each
  keeping its name: `TestRemovalRewriteLocatesAReducedEntryAtItsProvenance`
  (the `README.md` citation is rewritten: two files, count 2);
  `TestHistorySanitizePlanWritesNothing` (the `README.md` citation gains
  `; rewrites to <locator>`); `TestBaselineUpdateHistoryCitationBoundaries`,
  `TestBaselineUpdateListsHistoryCitationsTranscript` and
  `TestBaselineUpdatePromotesOnlyOnRequestTranscript` (the handoff citation
  carries `rewrite` to the Archive Record);
  `TestBaselineUpdateReportsStaleHistoryCitationsTranscript` (its stale
  citation is made by deleting the file outside the conversion).
- Build Order 2, `internal/spec/history_adr_identity_test.go` (new):
  `TestHistoryADRLineageNamesADateNamedADRByItsFileName`,
  `TestHistoryADRLineageDatesTheSupersession`,
  `TestHistoryADRNumericReferencesSelectOnlyNumberedADRs`.
  `TestHistoryADRPlanWritesTheLineageASuccessorLacks` stays unchanged: its
  retired ADRs have no `updated_at`, and their last commit is the dated
  fixture commit.
- Build Order 3, `internal/spec/live_links_test.go`:
  `TestLiveLinksLeaveAPendingTaskCreatesTarget`,
  `TestLiveLinksResolveABrokenLinkByTheRemovedFileName`.
  `internal/spec/pending_targets_test.go` (new):
  `TestRemovedHistoryPathsListsWhatHEADHeldAndTheTreeLacks`.
  `internal/speccheck/history_retained_test.go`:
  `TestHistoryRetainedSkipsAPendingTaskCreatesTarget`,
  `TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames`.
  `internal/cli/history_citations_transcript_test.go` (new):
  `TestHistorySanitizeRewritesEveryCitationTranscript`, asserting Surface
  Transcript 1 whole.
- Build Order 4, `internal/speccheck/history_retained_test.go`:
  `TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator`, on the
  shape an earlier conversion left; each of the five citations is one error
  with its fix, and a Git Locator, an Archive Record, a skill tree file, a
  `testdata` file and a pending `creates:` target are not reported.

## Build Order

1. Citation Sources and the full rewrite: `CitationSources`,
   `SourceHistoryCitations`, `CitationRewrite.Relative` and `ApplyText`,
   `PlanHistoryRewrite` without `scope`, their tests and the declared breaks
   (no dependency).
2. Lineage identity and date: `historyADRIdentifier`,
   `historyADRSupersessionDate`, the reference forms and their tests
   (depends on: 1, because both change `internal/spec/history_entries.go`).
3. Pending Task targets and broken links: `PendingTaskCreates`,
   `RemovedHistoryPaths`, `ResolveMissingLinkByName` in `PlanLiveLinks` and in
   `SC-HISTORY-RETAINED`'s missing-link finding, their tests, and the Surface
   Transcript 1 test, which needs Build Orders 1 and 2 (depends on: 1, 2).
4. The removed-citation finding and the documentation: the new
   `SC-HISTORY-RETAINED` finding, the `baseline`, `history` and `spec`
   guides, the Roundfix skill and references, `CONTEXT.md`, the skill
   version and the Behavior Surface Record (depends on: 1, 3, because it
   reads `SourceHistoryCitations`, `RemovedHistoryPaths` and
   `PendingTaskCreates`).
5. QA gate (depends on: 4).

No other Spec is a prerequisite. Specs 0285 and 0287 may merge first and move
the Roundfix skill version and the Behavior Surface Record, so this plan is
rebased on the default branch and checked again before its Run.

## Risks & Considerations

- Reading every tracked text file costs one read per file per plan or check.
  Roundfix holds about 3,000 tracked files, so the cost is milliseconds, but a
  repository that tracks very large text files reads them in full. Measure in
  the QA gate and record it.
- A Go or shell file that names a removed history file in code changes the
  program text. It is listed first and bound into the digest, and only a path
  the plan removes can match, so a synthetic test path never does.
- A citation inside a setup-marker region is rewritten like any other; the
  Baseline never cites an adopter's history files, so the next Baseline
  update has nothing to restore.
- `updated_at` edited after the supersession dates the line with the later
  edit; the docs-layout guide asks for `updated_at` to change with the
  status, which makes it the best recorded value.

## Glossary

- adds: **Citation Source**
- adds: **Removed History Path**
- changes: **History Citation**
- changes: **Git Locator**

## Decisions

- Every listed citation in every Citation Source is rewritten; skill trees,
  `testdata`, binary files and symbolic links are not sources. See ADR-0299.
- Only `./` and `../` prose paths are resolved as relative. See ADR-0299.
- Outside Markdown a Git Locator is written without backticks. See ADR-0299.
- The lineage identifier comes from the file name and the date from
  `updated_at`, with commit dates as fallbacks. See ADR-0299.
- A pending `creates:` target is exempt in `live-links` and the check. See
  ADR-0299.
- File-name resolution needs exactly one candidate. See ADR-0299.
- The check extends `SC-HISTORY-RETAINED` and exempts `docs/history/`. See
  ADR-0299.
- `HistoryCitations` stays Markdown-only, so the Archive Relink and the
  `superseded` record check keep their behavior.
