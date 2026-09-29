---
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
prd: _prd.md
created: 2026-09-29
---

# Baseline plans that show what a History Relocation breaks — Technical Spec

## Executive Summary

A new Relocation Citation scan runs inside Baseline planning, right where
`planHistoryMoves` has built the ordered History Relocation ledger, and only
when that ledger is non-empty. It lists the Git index once and reads each
tracked text file once. It reports every citation that resolves to a tracked
target before the plan and to nothing after it, with both sides resolved from
the citing file's own before and after location. It returns ordinary Baseline
`Finding` warnings. They join the plan's existing `warnings`, so every renderer
already prints them, and the Plan Digest binds them without a schema change.

The trade-off accepted is that a citing-file edit which changes the impact
changes the digest. That forces a fresh review even though apply would write
the same bytes. The maintainer chose this binding (ADR-0173). No preimage, file
change, History Relocation or apply stage changes.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; warnings use the
  existing Baseline finding shape keyed by code and repository-relative path.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and tracked files only;
  no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0173 (this Spec) is implemented
  here. ADR-0071, ADR-0073, ADR-0103, ADR-0068, ADR-0120, ADR-0064 and ADR-0070
  hold as the PRD states, and the gate ADRs ADR-0080, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156 bind the QA Task. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md) (standing grant of 2026-09-18 for the
  Roundfix skill, and the 2026-09-29 plan approval); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

The feature extends one existing seam, `planHistoryMoves` in
`internal/baseline/plan.go`. It already returns the History Relocation ledger
together with the history-layout warnings (retained reviews and collisions).
After it sorts the ledger and assigns ordinals, it calls the new scan with the
repository root and the ledger, and appends the returned findings to its
warnings. The scan lives in a new file, `internal/baseline/history_citations.go`,
beside `history_layout.go`. It needs no new package, type exported to other
packages, CLI flag or schema field.

Every Baseline Plan goes through `BuildPlan` → `buildPlanWithCatalog` →
`planHistoryMoves`, which covers `baseline plan`, `baseline update` and the
interactive command. Their renderers already print each warning as
`Warning: <code>: <path>: <message>` in text and append it to `warnings` in
JSON:

- `baseline_profile.go`
- `baseline_update.go`
- `baseline_human.go`
- `baseline_apply.go`

`computePlanDigest` hashes the whole document except `fileChanges`, so the new
warnings are digest-bound without any edit there. `baseline apply` validates a
portable document's digest from the document itself and never re-plans. It
therefore never scans, and it cannot be affected by a citing file.

## Implementation Design

### Interfaces

```go
// relocationCitationFindings returns one warning per tracked file whose
// citations resolve before the given History Relocations and not after them,
// followed by at most one omitted and one unscanned summary. It returns nil
// without touching Git or the filesystem when moves is empty.
func relocationCitationFindings(
	ctx context.Context, root string, moves []HistoryMove,
) ([]Finding, error)

// listTrackedPaths returns the Git index paths of root, NUL-separated and
// untrimmed, deduplicated and sorted.
func listTrackedPaths(ctx context.Context, root string) ([]string, error)
```

`listTrackedPaths` runs
`git -C <root> -c core.fsmonitor=false ls-files -z --cached --full-name` with
`GIT_OPTIONAL_LOCKS=0` and `GIT_TERMINAL_PROMPT=0`, the environment
`ExecGitRunner` sets. It does not trim the output: it splits on NUL, drops
empty entries and duplicate stage entries, and keeps only paths for which
`repositoryPathIsSafe` holds and that contain no control character
(U+0000–U+001F, U+007F–U+009F), so a dropped path is never opened or printed.
A Git failure is a planning error, as every other
Git failure in planning is.

### The resolution model

Every rule in this section is a consequence of one criterion: **a citation
counts when it resolves today to an existing target and would not resolve after
the plan**. Both sides are computed from sets built once:

- **Files before.** The tracked paths.
- **Files after.** The files before, minus every applied move's `From`, plus
  every applied move's `To`. A move whose `To` is already a file before is not
  applied: `baseline apply` refuses it and leaves its source in place
  (`HistoryMoveRefusal` in `internal/baseline/transaction.go`). Such a move
  changes neither set, and its source keeps its after-location.
- **Unit directories.** Every directory that contains a move's `From` and lies
  strictly below that move's source root, as `historyMoveSourceRoot` gives it.
  Examples are `docs/specs/_archived/0012-x` or a legacy Review Artifact
  folder. Family and legacy roots are never unit directories. They exist after
  the plan when any file after lies beneath them.
- **A target exists before** when it is a file before, or a unit directory with
  a file before beneath it. **It exists after** when it is a file after, or a
  unit directory with a file after beneath it.

A citing file `C` has an after-location `C′`: the move's `To` when `C` is
relocated, else `C` itself.

- **Repository-path tokens**, in every scanned file. A token is a maximal run
  of letters, digits and `._~/@+-` that contains `/`. It is normalized by
  removing one leading `./` and trailing `.,:;/` characters. A normalized token
  `T` that exists before is a citation. It resolves to `T` both before and
  after, so it counts when `T` does not exist after.
- **Markdown link destinations**, in files whose name ends in `.md` or
  `.markdown` (case-insensitive). The scan reads inline links and images
  (`](dest)`) and reference definitions (a line starting with up to three
  spaces, `[label]:`, then the destination). It ignores lines inside fenced
  code blocks opened by ```` ``` ```` or `~~~`.
  - The destination is unwrapped from `<…>`, loses any title, and loses its
    `#fragment` and `?query`. It is then percent-decoded, and kept raw when
    decoding fails.
  - A destination with a URL scheme, one starting with `//`, an empty or
    fragment-only destination, and one whose raw or decoded form contains a
    control character are skipped.
  - A destination starting with `/` resolves from the repository root.
    Otherwise it resolves from `path.Dir(C)` before and from `path.Dir(C′)`
    after. A resolution that leaves the root is skipped.
  - It counts when the before-resolution exists before and the
    after-resolution does not exist after.
- The same `(line, before-resolution)` pair is reported once, even when both
  forms match it.

A co-moved pair keeps resolving, because the after-resolution from `C′` lands
on the partner's `To`. That is why no special case is needed for documents that
move together, or for relocated files that link outward.

### What the scan reads

The scan walks the tracked paths in sorted order.

- **Symbolic links.** It checks each path component with `os.Lstat`, caching
  directories, and skips a path whose final component or any parent is a
  symbolic link. It also skips a path that is not a regular file. It then
  opens the path with `O_RDONLY|O_NOFOLLOW` and reads it only when
  `os.SameFile` holds between the opened file's `Stat` and that `os.Lstat`, so
  a path swapped for a link between the check and the open is never read.
- **Binary files.** A file whose first 8,000 bytes hold a NUL byte is skipped
  as binary.
- **Unscanned files.** A regular file larger than 4 MiB, or one that fails to
  open or read, is not scanned. It is collected for the unscanned summary.
- **Untracked, ignored or outside files.** None of these is ever stat-ed or
  opened.
- **No content echo.** A finding carries the citing path, line numbers, the
  cited text, and resolved repository paths. By construction that is a
  path-shaped token that equals a tracked path, or a link destination.

### Data Models

No schema changes. The scan adds warnings with three new codes to the existing
`Finding{Code, Path, Message}` list:

- **`baseline.history.citation`.** Path is the citing file's current path. The
  message joins up to three citations, ordered by line and cited text, with
  `; `. Each citation reads
  `line <L> cites <T>, which resolves to <P>; after this plan <after>`, where
  `<after>` is one of two phrases:
  - `nothing exists there`, when the after-resolution equals `<P>`;
  - `it resolves to <P′>, where nothing exists`, otherwise.

  When `<P>` is a relocated file or lies in a relocated unit directory, the
  citation appends ` (it moves to <destination>)`. With more than three
  citations, the message ends in `; and <N> more`.
- **`baseline.history.citation.omitted`.** It is emitted once, when more than
  200 files cite. Path is `.`, and the message reads
  `<N> more tracked files cite paths this plan relocates and are not listed`.
- **`baseline.history.citation.unscanned`.** It is emitted once, when any file
  was unscanned. Path is `.`, and the message reads
  `<N> tracked files were not scanned for citations (larger than 4 MiB or unreadable): <first three paths>`.

The citation findings are ordered by path, then comes the omitted summary, then
the unscanned summary. They follow the existing history-layout warnings in
`planHistoryMoves`'s result.

### API Contracts

1. API Contract: `roundfix baseline plan` and `roundfix baseline update`, in
   text and JSON, and the interactive plan review. For a plan with History
   Relocations, the output lists the Relocation Citation warnings in the
   existing warnings form and JSON array. The plan's `planDigest` and the
   update result's `planDigest` cover them. Exit codes, flags and the
   `roundfix/baseline-plan/v1` and `roundfix/baseline-update-result/v1` schema
   versions do not change. A plan without History Relocations produces exactly
   the output it produced before.
2. API Contract: `roundfix baseline update --confirm-plan <digest>`, when a
   citing file changed after the preview in a way that changes the reported
   impact, refuses the stale digest through the existing digest-mismatch path.
   `roundfix baseline apply --plan <file> --confirm-plan <digest>` applies a
   portable plan exactly as before, without scanning.

## Coverage Map

- Goal 1 → relocationCitationFindings; API Contract 1.
- Goal 2 → The resolution model.
- Goal 3 → Data Models (digest-bound warnings); API Contract 2.
- Goal 4 → System Architecture (scan only when moves exist; apply untouched).
- Story 1 → API Contract 1; Data Models.
- Story 2 → API Contract 1 (JSON warnings).
- Story 3 → API Contract 2.
- Core Feature 1 → relocationCitationFindings; System Architecture.
- Core Feature 2 → The resolution model.
- Core Feature 3 → Data Models.
- Core Feature 4 → System Architecture; API Contract 2.
- Core Feature 5 → What the scan reads; listTrackedPaths.
- Core Feature 6 → System Architecture; API Contract 1.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 4.

## Integration Points

- **History layout discovery.** The scan consumes the ledger that
  `planHistoryMoves` builds, including occupied-destination collisions, which
  the ledger lists. It uses `historyMoveSourceRoot` to find unit directories.
  Discovery itself is unchanged.
- **Plan digest and validation.** `computePlanDigest` and `ValidatePlanDocument`
  are unchanged. Warnings are already part of the hashed payload, and shape
  validation does not constrain warning codes.
- **The Roundfix skill and user guide.** They describe the new warnings.

## Testing Approach

1. **Scanner unit tests.** A new `internal/baseline/history_citations_test.go`
   calls `relocationCitationFindings` over real temporary Git repositories. It
   covers:
   - a repository-path citation, an inline relative link, an image, a
     reference definition and a root-relative link to a relocated ADR, each
     reported once with its line;
   - a relocated ADR's relative link to an ADR that stays, reported with its
     after-resolution;
   - a relative link between two co-relocated ADRs, not reported;
   - a citation already broken before the plan, not reported;
   - a citation inside a fenced code block and a URL with a scheme, not
     reported as links;
   - an untracked file and an ignored file holding citations, neither opened
     nor reported, proven by making them unreadable;
   - a tracked symbolic link and a file under a symlinked directory, not
     followed;
   - a binary file, skipped;
   - an oversized file, reported only in the unscanned summary;
   - the three-citation and 200-file caps;
   - a unit directory citation reported, while a family-root citation is not;
   - no Git call and no finding when moves is empty;
   - byte-identical results across two runs.
2. **Plan-level tests.** A new
   `internal/baseline/history_citation_plan_test.go` uses the package's
   existing plan fixtures (`newPlanRepository`, `writeInspectionFile`,
   `commitInspectionRepository`, `buildTestPlan`) without editing them. It
   covers:
   - a plan whose relocated ADR is cited carries the warning, while a current
     layout plan has none;
   - two repositories differing only in a citing file have different Plan
     Digests and identical preimages, postimages, History Relocations and file
     changes;
   - applying both plans yields identical trees apart from the citing file.
3. **CLI tests.** A new `internal/cli/baseline_history_citation_test.go` reuses
   the existing Baseline Plan test helpers without editing them. It covers:
   - `baseline plan --format=text` prints the warning line;
   - `--format=json` carries it in `warnings`;
   - `baseline update --confirm-plan` with a digest taken before a citing
     file changed is refused and leaves the tree unchanged.
4. **QA.**
   - The terminal QA Task replays Fluxus commit `986f9dc` read-only. It
     restores the sixteen relocated ADRs to `docs/adr/` in a disposable
     repository and runs the built `baseline plan`. It compares the warnings
     with a link recount at that commit that is independent of Roundfix.
   - It times `baseline update` on a clone of this repository with one added
     relocation, against the v0.20.0 binary on the same clone.
   - It records the repository Verification as a fact.

## Build Order

1. Relocation Citation scan (depends on: none).
2. Plan wiring and public-contract tests (depends on: 1).
3. Guide, skill and glossary (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **Noise in large histories.** A repository with many relocations and many
  citing files could flood the review. The per-file three-citation cap and the
  200-file cap bound the output. The omitted summary keeps the rest countable.
- **Cost.** The scan reads every tracked text file once, but only when the
  plan relocates something. The reads use map lookups, not per-target string
  searches. The QA measurement bounds the added time on this repository's size.
- **False negatives by design.** Prose relative paths, HTML anchors and wiki
  links are not recognized. The PRD records this limit, and the warning never
  claims completeness.
- **A collided relocation.** A move that apply later refuses is still counted
  as moving, because the plan lists it. The PRD records this limit.
- **Secrets.** Only index paths are read, and only path-shaped text is echoed.
  A secret committed to a tracked file never reaches a message, because a
  message quotes only a token that equals a tracked path or a link destination.

## Decisions

- Report inside the same plan and digest; no deferral flag and no second
  digest. See ADR-0173.
- One resolution criterion, applied from the citing file's before and after
  locations, instead of per-case rules.
- A new file in the existing package, called from the existing seam. No new
  package or exported API.
- Read only from the Git index, never follow symbolic links, and echo only
  path text.
- Bound the output with per-file and per-plan caps, and summarize what is not
  listed or not scanned rather than dropping it silently.
