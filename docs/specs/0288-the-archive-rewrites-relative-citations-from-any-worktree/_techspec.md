---
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
prd: _prd.md
created: 2026-10-20
---

# The archive rewrites relative citations from any worktree — Technical Spec

## Executive Summary

The Archive Command and the Cancel Command take the repository root from
`git rev-parse --show-toplevel`, which prints the resolved path, and the Spec
Root from the configuration loaded at the work directory, which keeps the
spelling `os.Getwd` returned. When the work directory crosses a symbolic link
the two spellings share no prefix, and every `filepath.Rel` between them
produces a `../…` path that is not repository-relative. The citation rewrite
then fails, the Archive Relink plans nothing, and the cancel's citation
advice reads a path that does not exist.

One exported resolver in `internal/spec`, `RepositoryPath`, becomes the only
way the two commands turn an absolute path into a repository-relative one:
lexical first, then both sides resolved with `config.CanonicalizeResolvedPath`,
and a refusal when the path still lies outside. `spec.Archive`, `spec.Cancel`
and the command code in `internal/cli/archive.go` and
`internal/cli/cancel.go` call it wherever they build such a path. No command,
flag, output line, record field or configuration key is new.

The trade-off accepted: the resolver tries the lexical spelling first, so a
Spec Root reached through a link inside the repository keeps the bytes it has
today, at the cost of a second comparison only when the first one leaves the
root.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0247,
  ADR-0260, ADR-0278, ADR-0262, ADR-0154, ADR-0268, ADR-0272 and ADR-0299;
  Build Order 2 follows ADR-0247, ADR-0260, ADR-0262, ADR-0154, ADR-0269,
  ADR-0184, ADR-0280 and ADR-0256; every step follows ADR-0259, ADR-0271,
  ADR-0285, ADR-0104, ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named
  with its title in `_prd.md` → Project Constraints. Set no obligation here,
  because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0165, ADR-0167, ADR-0169, ADR-0177, ADR-0196, ADR-0197, ADR-0229, ADR-0230, ADR-0233, ADR-0237, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0258, ADR-0261, ADR-0277, ADR-0279, ADR-0281, ADR-0282, ADR-0286, ADR-0290, ADR-0295 and ADR-0298, each named in `_prd.md` → Project Constraints.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AL approval of 2026-10-20 ("Mesmas regras do AK (Recomendado)"), carrying
  the cycle AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0288-the-archive-rewrites-relative-citations-from-any-worktree/_authorization.md`;
  bounded files: `internal/spec/archive.go`.

## Current behavior

Characterized on 2026-10-20 at `7a678ddb` with the installed `roundfix
0.97.0`, in a temporary repository under a temporary `HOME` and
`GIT_CONFIG_GLOBAL`: a main checkout at `<tmp>/real/main`, a linked worktree
added with `git worktree add <tmp>/real/item`, outside the main checkout, and
a symbolic link `<tmp>/link -> <tmp>/real`. The Spec `0001-demo` had one
completed Task and a pending QA Task, and its `_prd.md` summary held
`[backlog of 2026-10-20](../../backlog/2026-10-20-follow-up.md)`; that Backlog
Entry linked `../specs/0001-demo/_prd.md`. A bare `origin` held `main`.

1. From `<tmp>/link/item`, `git rev-parse --show-toplevel` printed
   `<tmp>/real/item`, while the shell's work directory was `<tmp>/link/item`.
2. `roundfix archive 0001-demo --qa-override --approval probe --reason probe`
   from `<tmp>/link/item`, from `<tmp>/link/item/docs` and from
   `<tmp>/link/main` exited 2 with `rebase archive outcome: rebase citation
   "../../backlog/2026-10-20-follow-up.md": Rel: can't make
   docs/backlog/2026-10-20-follow-up.md relative to
   ../../link/item/docs/history/specs` (`../../link/main/…` for the main
   checkout), the error of the adopted Inbox Entry with a shorter climb.
3. The same command from `<tmp>/real/item` exited 0 with `archived 0001-demo
   with QA override -> docs/history/specs/0001-demo.md; removed 5 file(s)
   (1142 bytes) kept in Git at <revision>; relinked 1 file(s)`, a record
   whose summary keeps `../../backlog/2026-10-20-follow-up.md`, and the
   Backlog Entry relinked to `../history/specs/0001-demo.md`.
4. With no link in the PRD summary, the archive from `<tmp>/link/item`
   exited 0 but printed no `relinked` count and no warning, and the Backlog
   Entry still linked the removed `../specs/0001-demo/_prd.md`. The same run
   from `<tmp>/real/item` relinked it.
5. `roundfix cancel 0001-demo --reason probe` from `<tmp>/link/item` exited 2
   with `plan cancellation citation advice: list archive citation targets:
   lstat <tmp>/real/item<tmp>/link/item/docs/specs/0001-demo: no such file or
   directory`; from `<tmp>/real/item` it cancelled the Spec and relinked the
   Backlog Entry.
6. The record-only `roundfix archive 0001-demo --pull-request 7` worked from
   both spellings, because Git accepts the absolute linked path, but once the
   record was on `origin/main` its refusal from `<tmp>/link/item` named
   `<tmp>/link/item/docs/history/specs/0001-demo.md` instead of
   `docs/history/specs/0001-demo.md`.

Why the Delivery Queue's own archive works: `commandDeliveryWorkflow.Archive`
runs `roundfix archive <slug>` as a child with `Dir` set to the item worktree
and an environment that keeps the queue process's `PWD`. That `PWD` names
another directory, so the child's `os.Getwd` returns the resolved path, as Go
documents ("if the environment variable PWD provides an absolute name, and it
is a name of the current directory, it is returned"), and both spellings
agree by accident. An operator who `cd`s into the item worktree through
`~/.roundfix/worktrees`, a symbolic link to another volume on the
maintainer's machine since 2026-10-10 (ADR-0268), gets the linked spelling.

Where the spellings mix, at `7a678ddb`:

- `internal/spec/archive.go`: `rel` (`:204-207`) relativizes the Spec Root,
  the Archive Root, the Spec folder and the record path, built from the
  configured spelling, against `req.RepositoryRoot`, Git's spelling, for
  `PlanArchiveRelink`; `recordRel` (`:350`) does the same for the Legacy Link
  Rebase. `source` (`:338-346`) already resolves both sides.
- `internal/spec/cancel.go`: `relative` (`:96-105`) refuses the linked Spec
  folder as outside the repository; `planSourceReturn` (`:426`) relativizes
  each returned source the same way.
- `internal/cli/archive.go`: `folderRel` (`:175`) for the refusal and the
  removal recovery, with a resolved-path special case only when ignored files
  are the only finding; `recordRel` (`:285`) in the record-only
  `--pull-request` form; and, for an external Spec Root, the four relink
  paths of `archive --plan` (`:501-516`).
- `internal/cli/cancel.go`: `folderRel` (`:112`) and every path of
  `cancelCitations` (`:282-316`).

`filepathRelSlash` in `internal/cli/archive.go` returns the absolute target
when the relative path would leave the base, so a mixed comparison there
never fails; it hands an absolute path to Git or to the reader instead.

Checked and not in the class, because both sides come from one root:
`roundfix history sanitize` (`internal/cli/history.go`, repository and Spec
Root both from the loaded Git root); the history conversion of `roundfix
baseline update` (`internal/cli/baseline_update_history.go`, which reloads
the configuration at Git's top level); `RelinkContent` and the rest of
`internal/spec/live_citations.go`, `history_entries.go` and `live_links.go`,
which join every path they relativize onto the same root; the Delivery
Queue's `archivePaths` and `archiveRetirementIsExact`
(`internal/cli/deliver_workflow.go`), which read repository-relative paths
from Git trees; the display paths `filepathRelSlash(loaded.GitRoot, …)` in
both commands, whose sides share the loaded root; and
`resolveHistoricalSpecsRoot` (`internal/cli/archive_reader.go`). The Daemon
calls no archive code (`rg -n 'spec\.Archive\(' internal/daemon` finds
nothing).

The existing tests never cross a link: `newImplementWorkspace` resolves its
repository with `filepath.EvalSymlinks` before any command runs, and the
external-root fixtures that do not resolve their temporary directory reach
the mixed path only on macOS, where `/var` is a link, with no relative link
in a summary and no Live Document to relink.

## System Architecture

The fix stays inside the two packages that build archive paths.

- `internal/spec/repository_path.go` (new): `RepositoryPath`.
- `internal/spec/archive.go`: `Archive` builds its relink paths and its
  record path with `RepositoryPath`.
- `internal/spec/cancel.go`: `Cancel` and `planSourceReturn` build their
  paths with `RepositoryPath`.
- `internal/cli/archive.go` and `internal/cli/cancel.go`: the refusal,
  recovery, record-only and citation-advice paths use `RepositoryPath`
  against the root Git reported.
- `docs/user-guide/commands/archive.md` and `docs/user-guide/commands/cancel.md`:
  one sentence each, saying that the command reads every path relative to
  the checkout it runs in and naming the spellings that covers.

`internal/spec` already may import `internal/config` (`go list -deps
./internal/config` lists only `internal/openrouterkey`), so the resolver
reuses `config.CanonicalizeResolvedPath`, which resolves the longest existing
prefix of a path that does not exist yet, instead of a second copy.

## Implementation Design

### Interfaces

```go
// RepositoryPath returns path relative to repositoryRoot, slash-separated.
// It compares the two lexically first; when that leaves the root, it compares
// config.CanonicalizeResolvedPath of both; when that also leaves the root, it
// fails with: path "<path>" is outside the repository "<repositoryRoot>".
// The root itself yields ".".
func RepositoryPath(repositoryRoot, path string) (string, error)
```

### Data Models

None. The Archive Record, the cancelled record and the relink plan keep
their shapes; only the spelling of the paths fed into them changes.

### API Contracts

1. API Contract: `RepositoryPath(repositoryRoot, path)`. Both arguments are
   absolute. The result is `filepath.ToSlash` of a clean relative path that
   never starts with `../` and is never absolute. A path lexically inside the
   root returns its lexical relative path without touching the file system.
   Otherwise both sides go through `config.CanonicalizeResolvedPath` and the
   comparison repeats. A path inside neither way returns the error named in
   the Interfaces sketch, naming both arguments as given. A sibling sharing a
   name prefix (`<tmp>/repo2` beside `<tmp>/repo`) is outside.
2. API Contract: `spec.Archive`. The four paths handed to `PlanArchiveRelink`
   (Spec Root, Archive Root, Spec folder, record) and the Legacy Link
   Rebase's `RecordPath` come from `RepositoryPath(req.RepositoryRoot, …)`.
   `source`, which already compares both sides resolved, keeps its
   computation and bytes. An error from `RepositoryPath` is returned wrapped
   with the operation that failed, before anything is written. QA
   eligibility, the QA Archive Override, the refusals and their order, the
   record fields and the Archive Relink and Legacy Link Rebase rules do not
   change.
3. API Contract: `spec.Cancel`. Its `relative` helper and `planSourceReturn`
   use `RepositoryPath(req.RepositoryRoot, …)`. A path `RepositoryPath`
   refuses, or one `cancelRelativePath` rejects, keeps today's message
   `cancellation path "<path>" is outside the repository`.
4. API Contract: the Archive Command (`internal/cli/archive.go`). The Spec
   folder the uncommitted, ignored-file and removal-recovery messages name,
   the record path the record-only `--pull-request` form hands to `git
   ls-tree` and names in its refusal, and the four relink paths of `archive
   --plan` come from `RepositoryPath` against the root Git reported for that
   step. The resolved-path special case for ignored files is removed, because
   the resolver covers it. A path outside the repository is a Preflight
   failure naming it; it is never passed on as an absolute path.
5. API Contract: the Cancel Command (`internal/cli/cancel.go`). The Spec
   folder its refusals and recovery name and every path `cancelCitations`
   builds come from `RepositoryPath` against the root Git reported.
6. API Contract: everything else keeps its bytes from every spelling: the
   output lines, the record and relinked bytes, the warnings and their order,
   and every existing test name and assertion.

### Surface Transcripts

Both run in the temporary repository described under Current behavior, with
the Spec committed on `main`, `origin/main` at the same commit, the linked
worktree checked out on a new branch at that commit, and the command started
with the work directory `<tmp>/link/item`, the worktree root spelled through
the symbolic link.

1. Surface Transcript: the QA Archive Override archives a Spec whose PRD links
   outside its folder from the linked spelling of a worktree.

   ```transcript
   $ roundfix archive 0001-demo --qa-override --approval probe --reason probe
   stdout:
   archived 0001-demo with QA override -> docs/history/specs/0001-demo.md; removed <files> file(s) (<bytes> bytes) kept in Git at <revision>; relinked 1 file(s)
   stderr:
   exit: 0
   ```

2. Surface Transcript: the Cancel Command cancels the same Spec from the
   linked spelling.

   ```transcript
   $ roundfix cancel 0001-demo --reason probe
   stdout:
   cancelled 0001-demo -> docs/history/specs/0001-demo.md; 1 completed Task(s), 1 pending; removed <files> file(s) (<bytes> bytes) kept in Git at <revision>
   stderr:
   exit: 0
   ```

### Invariants

```text
1. A path RepositoryPath returns never starts with ../ and is never absolute.
2. A path lexically inside the repository root keeps the relative path it has today.
3. From every spelling of one checkout, the archive and the cancel print the same lines and write the same bytes as from the resolved spelling.
4. A path outside the repository is refused with a message naming it; it is never rewritten, skipped or passed on as an absolute path.
5. Every Archive Record already written stays byte-identical.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Features 1-2, Success Metric 1, Invariants 1-3 → `RepositoryPath` (API Contract 1), `spec.Archive` (API Contract 2), the Archive Command (API Contract 4), Surface Transcript 1
- Goal 2, Core Feature 2, Success Metric 2, Invariant 3 → `spec.Cancel` (API Contract 3), the Cancel Command (API Contract 5), Surface Transcript 2
- Goal 3, Success Metric 3 → the Archive Command's refusal and record-only paths (API Contract 4)
- Goal 4, Success Metric 4, Invariant 4 → `RepositoryPath`'s refusal (API Contract 1), `spec.Archive` and `spec.Cancel` (API Contracts 2-3)
- Core Feature 3 → the Archive Command run with both work-directory spellings a Delivery Queue child can see (API Contracts 4 and 6)
- Core Feature 4 → the `archive` and `cancel` user guide sentences, the Coverage Review of their guide surfaces and the re-recorded Behavior Surface Record
- Success Metric 5 → the PRD's Unreachable Acceptance
- Core Feature 2, Invariant 5, API Contract 6 → the existing archive, cancel and citation tests, the Coverage Review of `command archive` and `command cancel` in `docs/references/skill-coverage.json`, and the QA scope audit

## Absence Evidence

1. `RepositoryPath` — command: `grep -rn 'CanonicalizeResolvedPath\|func RepositoryPath' internal/spec internal/cli/archive.go internal/cli/cancel.go`; result: no match at `7a678ddb`; no archive or cancel code compares a stored spelling with Git's resolved one, and `internal/spec/archive.go` resolves only `source` (`:338-346`) with `filepath.EvalSymlinks`.

## ADR Alternatives

- ADR-0247: no rejected alternative is taken; the archive still writes one small Archive Record and removes the folder, which stays in Git.
- ADR-0260: no rejected alternative is taken; the Archive Relink keeps its Live Documents, its promotion and its rewrite rules; it now receives repository-relative paths from every spelling.
- ADR-0278: no rejected alternative is taken; the Legacy Link Rebase keeps its rules and receives a repository-relative record path.
- ADR-0262: no rejected alternative is taken; the cancel keeps its record, its source return and its refusals.
- ADR-0154: no rejected alternative is taken; the override keeps its provenance and its refusals.
- ADR-0269: no rejected alternative is taken; both commands still refuse ignored and uncommitted Spec files and finish a failed removal forward; the paths they name become repository-relative.
- ADR-0268: no rejected alternative is taken; refusing every path that crosses a link is the alternative it rejected, and this Spec resolves both spellings instead, as it decided.
- ADR-0272: no rejected alternative is taken; like its review record, the archive paths are compared in Git's resolved spelling.
- ADR-0299: no rejected alternative is taken; Git Locator text and resolution are untouched.
- ADR-0184: no rejected alternative is taken; both changed command surfaces are stated as transcripts the QA gate runs through the built product.
- ADR-0280: no rejected alternative is taken; task_02 changes the archive and cancel sources and re-records their surface fingerprint itself.
- ADR-0256: no rejected alternative is taken; the covering Roundfix archive reference is reviewed and recorded unchanged, because no flag, line or exit code changes.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the vortex Inbox Entry and published documentation.
- ADR-0240: no rejected alternative is taken; the real item worktree run is declared under Unreachable Acceptance.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and record invariance compares with the Delivery Base.
- ADR-0265: no rejected alternative is taken; the new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0177: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0230: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0298: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

Git, through the existing `git` executable, for `rev-parse --show-toplevel`,
`status`, `show` and `ls-tree` in the checkout the command runs in; the file
system, for resolving symbolic links. No test, Verification command or QA
row reads the real `~/.roundfix` or reaches the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, homes and Git configuration, and creates its own
symbolic link with `os.Symlink` under its temporary directory, so it crosses a
link on Linux as well as on macOS. Each fixture resolves its temporary
directory first, so the only link a test crosses is the one it made.

- `internal/spec/repository_path_test.go` (new):
  - `TestRepositoryPathComparesResolvedSpellings`, with subtests: one
    spelling (lexical result, no file-system access needed); root resolved
    and path through the link; root through the link and path resolved; a
    record path that does not exist yet under a linked parent; the root
    itself (`.`); a sibling sharing a name prefix (refused); a path outside
    (refused, message naming both).
  - `TestArchiveFromALinkedSpellingRewritesAndRelinks`, through
    `spec.Archive` with a QA override, a PRD summary linking
    `../../backlog/<file>.md` and a Backlog Entry linking the Spec's
    `_prd.md`: subtests with the Spec Root spelled through the link and with
    the repository root spelled through the link each succeed and write the
    same record bytes, relinked bytes and `Relinked` list as a control run
    with one spelling and the same archive date.
  - `TestCancelFromALinkedSpellingRelinksAndReturnsSources`, through
    `spec.Cancel` with an adopted Backlog Entry under `references/` and a
    Live Document citing the Spec: the same comparison, and every
    `Returned` path repository-relative.
- `internal/cli/archive_symlinked_checkout_test.go` (new), through the
  command entry point with the command environment's work directory set to
  the spelling under test, which is what `os.Getwd` returns when `PWD` names
  it:
  - `TestArchiveFromEverySpellingOfTheCheckoutTranscript`, with subtests:
    the worktree root resolved (control), the worktree root through the link
    (Surface Transcript 1 asserted whole), a worktree subdirectory through the
    link, and the main checkout through a link; each exits 0 with the same
    standard output and standard error as the control, and writes the same
    record and relinked Backlog Entry bytes.
  - `TestCancelFromALinkedSpellingTranscript`: the control and the worktree
    root through the link (Surface Transcript 2 asserted whole), with the
    same record and relinked bytes.
  - `TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling`: from the
    link spelling, the uncommitted-change refusal names `docs/specs/<slug>`
    and the record-only `--pull-request` refusal for a record already on the
    default branch names `docs/history/specs/<slug>.md`, each byte-identical
    to the control's refusal.
  - `TestArchivePlanWithALinkedExternalSpecRootListsTheRelink`: `archive
    --plan --offline` with an external Spec Root configured through a link
    prints `relinks docs/backlog/<file>.md: 1 citation(s)` for the external
    repository's Backlog Entry that cites the Spec.
- `docs/references/skill-coverage.json`: task_02 records the Coverage Review
  of the four surfaces `_prd.md` → Skills lists unchanged (`command archive`,
  `command cancel` and the two guides), as a dated `review` on each, and
  re-records `docs/references/behavior-surfaces.json` with the repository's
  record command, because the two guide fingerprints change with their new
  sentence; the command fingerprints come from help text, which does not
  change. `TestTheSkillCoverageMapIsCurrent` accepts the result.
- Declared break: none. No existing test, golden or fixture changes; the
  removed ignored-file special case leaves `TestArchiveRefusesIgnoredSpecFiles`,
  `TestCancelRefusesIgnoredSpecFiles` and
  `TestSpecFolderStatusReadsTheExternalSpecsRootRepository` passing with
  their bytes.

## Build Order

1. The resolver and the package code: `RepositoryPath`, `spec.Archive`,
   `spec.Cancel` and their tests (no dependency).
2. The command code: the Archive Command's refusal, record-only and plan
   paths and the Cancel Command's refusal, recovery and citation-advice
   paths, their transcript tests, the `archive` and `cancel` guide
   sentences, the Coverage Review of the four surfaces and the Behavior
   Surface Record (depends on: 1, which exports the
   resolver the commands call and fixes the package code their transcripts
   pass through).
3. QA gate (depends on: 1, 2).

No other Spec is a prerequisite. Spec 0287 merges first and touches
`internal/spec/archive_record.go`, which this Spec does not change; the plan
is rebased and checked again before its Run.

## Risks & Considerations

- A Spec Root that is itself a link inside the repository resolves lexically
  first, so its paths keep today's bytes; only a comparison that leaves the
  root lexically reaches the resolved spelling.
- A link whose target is outside the repository resolves outside it and is
  refused with a named path, where today it would be rewritten as `../…` or
  passed on as an absolute path.
- Resolving costs one `EvalSymlinks` walk per path, and only when the lexical
  comparison leaves the root; an archive builds fewer than ten such paths.
- Live Documents that an earlier archive from a linked spelling left citing a
  removed folder are not repaired here; `roundfix spec check` already reports
  them as broken citations.

## Glossary

None.

## Decisions

- Fix the callers, not `CitationRewrite.Apply`: the rewrite's inputs are
  documented as repository-relative ("All paths and promotion keys are
  repository-relative and slash-separated", `internal/spec/legacy_link_rebase.go`),
  and the callers broke that contract.
- One resolver, exported from `internal/spec`, reusing
  `config.CanonicalizeResolvedPath`; no third copy of the resolution.
- Lexical first, resolved second, refusal third, so every spelling that works
  today keeps its bytes and a real escape is named.
- Keep the Delivery Queue's archive step as it is: it runs the same command,
  so the command's fix covers it.
- No new ADR: ADR-0268 and ADR-0272 already decided that a stored spelling and
  Git's spelling are compared resolved.
