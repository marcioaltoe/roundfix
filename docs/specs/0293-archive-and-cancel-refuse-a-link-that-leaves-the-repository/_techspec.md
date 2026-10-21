---
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
prd: _prd.md
created: 2026-10-20
---

# Archive and cancel refuse a link that leaves the repository — Technical Spec

## Executive Summary

`RepositoryPath` in `internal/spec/repository_path.go` is the one resolver
`roundfix archive` and `roundfix cancel` use to turn a path into a
repository-relative one. It returns as soon as the path is lexically inside
the repository root (lines 25-27), so a path that crosses a symbolic link
whose target lies outside is accepted, and the commands write through it. The
archive's promotion destinations and the cancel's promotion and
returned-source destinations never pass through the resolver at all.

The resolver gains a containment test: a path is accepted only when
`config.CanonicalizeResolvedPath` of the path lies inside
`config.CanonicalizeResolvedPath` of the repository root. When the path is
lexically inside, it still returns the lexical relative path; otherwise it
returns the resolved one, as today. The three destination joins call the
resolver before anything is written. Every refusal keeps the text its step
already prints for a path outside the repository.

The trade-off accepted: a repository that links a destination directory
outside itself can no longer archive or cancel (ADR-0306 records the break and
narrows ADR-0269), in exchange for never writing a record, a promoted copy or
a returned source where no commit can hold it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0306,
  ADR-0268, ADR-0272, ADR-0247, ADR-0260, ADR-0278, ADR-0262,
  ADR-0154 and ADR-0299; Build Order 2 follows ADR-0306, ADR-0269,
  ADR-0262, ADR-0184 and ADR-0280; Build Order 3 follows ADR-0306, ADR-0269
  and ADR-0256; every step follows ADR-0259, ADR-0271, ADR-0285, ADR-0104,
  ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named with its title in
  `_prd.md` → Project Constraints. Set no obligation here, because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0165, ADR-0167, ADR-0169, ADR-0177, ADR-0196, ADR-0197, ADR-0229, ADR-0230, ADR-0233, ADR-0237, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0258, ADR-0261, ADR-0277, ADR-0279, ADR-0281, ADR-0282, ADR-0286, ADR-0290, ADR-0295, ADR-0298 and ADR-0301, each named in `_prd.md` → Project Constraints.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20, authorization "the same rules as cycle AL",
  carrying the cycle AJ grant "Mexer em Governed Paths declarados" and the
  2026-09-30 skills grant. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository/_authorization.md`;
  bounded files: `internal/spec/archive.go`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `70c78d15` with `bin/roundfix` built by `make
build` from that commit, in temporary repositories under a temporary `HOME`
and `GIT_CONFIG_GLOBAL`, each holding a Spec `0001-demo` with one completed
Task and a pending QA Task, archived with `--qa-override --approval probe
--reason probe` and cancelled with `--reason probe`. The temporary directory
was resolved first (`/private/var/…` on macOS), so the only links crossed
were the ones each case made.

1. Control: archive exited 0 with `archived 0001-demo with QA override ->
   docs/history/specs/0001-demo.md; removed 5 file(s) (1121 bytes) kept in Git
   at <revision>`.
2. `docs/history` committed as a link to `<tmp>/outside`: archive printed the
   same line and exited 0; the record was `<tmp>/outside/specs/0001-demo.md`,
   and `git status --porcelain` listed only the five deleted Spec files. Cancel
   did the same (`cancelled 0001-demo -> docs/history/specs/0001-demo.md; 1
   completed Task(s), 1 pending; …`, exit 0, record outside).
3. `docs/references` linked outside: `archive … --promote references/notes.md`
   exited 0 with `; promoted 1 file(s) to docs/references/` and wrote
   `notes.md` outside.
4. `docs/backlog` linked outside, the Spec adopting
   `2026-10-20-entry.md` as `backlog`: cancel exited 0, printed `returned
   docs/backlog/2026-10-20-entry.md as open`, and wrote the entry outside.
5. The Spec folder itself a committed link to `<tmp>/outside/0001-demo`:
   archive exited 0 with `removed 1 file(s) (96 bytes)`, building the record
   from files Git never held; cancel already refused with `unknown active Spec
   slug "0001-demo"` (`internal/cli/cancel.go:70-73` refuses a folder that is
   not a directory).
6. `docs/history` a link to `docs/hist-real` inside the repository: archive
   exited 0 and wrote `docs/hist-real/specs/0001-demo.md`.
7. The checkout reached through a link to it, and through the macOS `/var`
   spelling of `/private/var`: archive exited 0 with the control's line.
8. An external Spec Root configured as `specs.root:
   /var/…/knowledge/specs`, inside its own Git repository: archive exited 0
   and wrote `knowledge/specs/_archived/0001-demo.md`.
9. A tracked Live Document whose directory was replaced by a link outside:
   archive already refused with `citing file
   docs/backlog/2026-10-20-follow-up.md is not a regular repository file`
   (`internal/spec/history_entries.go:884`).
10. `docs/specs` itself a link into another Git repository: archive already
    refused with `read Spec folder status: fatal: … is outside repository at
    …`, because Git refuses the linked spelling; `docs/specs` linked to a
    directory outside any repository already refused with `fatal: not a git
    repository`.
11. `archive 0001-demo --plan --offline` with `docs/history` linked outside:
    exit 0, planning `writes docs/history/specs/0001-demo.md`.
12. `docs/references` linked outside and a committed open Backlog Entry
    citing the Spec's `references/notes.md`: cancel exited 0, promoted the
    cited file and wrote `notes.md` outside.

Where a path reaches a write without a containment test, at `70c78d15`:

- `internal/spec/repository_path.go:25-27`: the lexical branch returns before
  either side is resolved. Every caller below inherits it.
- `internal/spec/archive.go:308`: the promotion destination
  `docs/references/<basename>` is joined to the repository root and written.
- `internal/spec/cancel.go:174`: the cancel's promotion destination, the
  same way.
- `internal/spec/cancel.go:406`: the returned-source destination
  `docs/backlog/<file>` or `docs/findings/<file>`, the same way.

`RepositoryPath` has 21 call sites, all in the two commands: in
`internal/spec/archive.go` (`:204`, `:208`, `:212`, `:216`, `:346`),
`internal/spec/cancel.go` (`:97`, `:423`), `internal/cli/archive.go` (`:169`,
`:279`, `:496`, `:501`, `:506`, `:511`) and `internal/cli/cancel.go` (`:107`,
`:277`, `:289`, `:303`, `:307`, `:311`). `rg -n
'(^|[^A-Za-z_])RepositoryPath\(' -g '*.go'` finds no other caller outside
`internal/spec/repository_path_test.go`, so no other command changes. The
display paths `filepathRelSlash(loaded.GitRoot, …)` in both commands, which
only format output, keep their behavior; `source` in `spec.Archive`
(`:334-345`) keeps its computation, because the Spec folder is refused
earlier.

Already refused and not in the class: promotion sources and indexed
reference sources that cross a link (`promoted path … is not a regular file`,
`indexed cancellation path … is a symbolic link`), a Spec file that is a link
(`Spec contains symlink` in `archive --plan`), relinked Live Documents (case
9), and the cancel's Spec folder link (case 5). `roundfix history sanitize`,
the history conversion of `roundfix baseline update` and the Delivery Queue's
`archivePaths` never call the resolver.

## System Architecture

- `internal/spec/repository_path.go`: `RepositoryPath` adds the containment
  test.
- `internal/spec/archive.go`: the promotion destination passes through it.
- `internal/spec/cancel.go`: the promotion and returned-source destinations
  pass through it.
- `internal/cli/archive.go` and `internal/cli/cancel.go`: no change; every
  path they build already goes through the resolver.
- `docs/user-guide/commands/archive.md`, `docs/user-guide/commands/cancel.md`,
  `.agents/skills/roundfix/references/archive.md` (with its mirror): one
  paragraph each on the refusal. `docs/adr/0269-…`: a dated line naming
  ADR-0306.

## Implementation Design

### Interfaces

```go
// RepositoryPath returns path relative to repositoryRoot, slash-separated.
// It accepts path only when config.CanonicalizeResolvedPath(path) lies inside
// or equals config.CanonicalizeResolvedPath(repositoryRoot). It returns the
// lexical relative path when path is lexically inside, and the resolved one
// otherwise; when the resolved path lies outside it fails with:
// path "<path>" is outside the repository "<repositoryRoot>".
func RepositoryPath(repositoryRoot, path string) (string, error)
```

### Data Models

None. No record, plan or configuration shape changes.

### API Contracts

1. API Contract: `RepositoryPath(repositoryRoot, path)`. Both sides are
   always resolved with `config.CanonicalizeResolvedPath`, which resolves the
   longest existing prefix and keeps the missing remainder, so a record path
   that does not exist yet under a linked parent is judged by its parent's
   target. A resolved path inside or equal to the resolved root is accepted;
   the result is the lexical relative path when the lexical comparison stays
   inside, so every path accepted today with a target inside keeps its bytes,
   and the resolved relative path otherwise. A resolved path outside is
   refused with today's message, naming both arguments as given. The root is
   resolved as well as the path: a root and a path both spelled through a
   link to the checkout, such as macOS `/var` or a linked
   `worktree.location`, are inside.
2. API Contract: `spec.Archive`. Before the vacancy check of each promotion
   destination (`internal/spec/archive.go:308`), it calls
   `RepositoryPath(req.RepositoryRoot, destination)` and returns
   `resolve promotion destination: <resolver error>` before any write. Its
   other resolver calls keep their wrapping, for example `resolve Archive Root
   for archive relink: <resolver error>`.
3. API Contract: `spec.Cancel`. Before the vacancy check of each promotion
   destination (`:174`) and each returned-source destination (`:406`), it
   calls `RepositoryPath` and returns `cancellation path "<destination>" is
   outside the repository` before any write, the text its `relative` helper
   already uses.
4. API Contract: the Archive Command and the Cancel Command. Every resolver
   refusal, from the command code or from `spec.Archive` and `spec.Cancel`,
   is a Preflight failure: exit `2`, the refusal as the reason, nothing
   written, the Spec folder in place. The command code does not change. The
   first refusal follows today's order of steps: in `archive`, the Spec
   folder (`internal/cli/archive.go:169`), then `spec.Archive`'s Spec Root,
   Archive Root, Spec folder and record, then each promotion destination; in
   `archive --plan` and the record-only `--pull-request` form, their own
   resolver calls; in `cancel`, the Spec folder (`:107`), then the citation
   advice (`cancelCitations`), then `spec.Cancel`.
5. API Contract: everything else keeps its bytes: the output, records and
   relinked files of every path whose target stays inside, every existing
   refusal and its order, and every existing test name and assertion.

### Surface Transcripts

Both run in a temporary repository whose resolved root is `<repo>` and whose
one commit holds the Spec `0001-demo` of Current behavior. The command runs
with the work directory `<repo>`. `<repo>` is the only controlled variation.

1. Surface Transcript: the archive refuses an Archive Root reached through a
   `docs/history` link to a directory outside the repository.

   ```transcript
   $ roundfix archive 0001-demo --qa-override --approval probe --reason probe
   stdout:
   stderr:
   Preflight failed

   Reason:
     resolve Archive Root for archive relink: path "<repo>/docs/history/specs" is outside the repository "<repo>"

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix archive --help' for usage.
   exit: 2
   ```

2. Surface Transcript: the cancel refuses to return an adopted Backlog Entry
   into a `docs/backlog` link to a directory outside the repository.

   ```transcript
   $ roundfix cancel 0001-demo --reason probe
   stdout:
   stderr:
   Preflight failed

   Reason:
     cancellation path "<repo>/docs/backlog/2026-10-20-entry.md" is outside the repository

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix cancel --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. RepositoryPath accepts a path only when its resolved target lies inside the resolved repository root.
2. A path lexically inside the repository whose target lies inside keeps the relative path it has today.
3. A checkout or an external Spec Root reached through a symbolic link is inside its repository.
4. Archive and cancel write nothing, and leave the Spec folder in place, when they refuse a path.
5. Every Archive Record already written stays byte-identical.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Features 1-3, Success Metric 1, Invariants 1 and 4 → `RepositoryPath` (API Contract 1), `spec.Archive` (API Contract 2), the Archive Command (API Contract 4), Surface Transcript 1
- Goal 2, Core Features 2-3, Success Metric 2, Invariant 4 → `spec.Cancel` (API Contract 3), the Cancel Command (API Contract 4), Surface Transcript 2
- Goal 3, Success Metric 3, Invariants 2-3 → API Contracts 1 and 5, the inside and linked-checkout cases of the Testing Approach
- Success Metric 4 → API Contracts 1-3 and the package tests
- Goal 4, Core Feature 4 → the `archive` and `cancel` guides, the Roundfix `archive` reference and its version, ADR-0306 and the ADR-0269 line
- Core Feature 3, Invariant 5, API Contract 5 → the existing archive, cancel and citation tests and the QA scope audit

## Absence Evidence

1. Containment test in `RepositoryPath` — command: `sed -n '17,32p' internal/spec/repository_path.go`; result: at `70c78d15` the lexical branch (`:25-27`) returns before `config.CanonicalizeResolvedPath` is called, and the resolved comparison runs only when the lexical one leaves the root.
2. Destination check for promotions and returned sources — command: `grep -n 'RepositoryPath' internal/spec/archive.go internal/spec/cancel.go`; result: no call near `internal/spec/archive.go:308`, `internal/spec/cancel.go:174` or `:406`; the destinations are joined to the repository root and only checked for vacancy.

## ADR Alternatives

- ADR-0306: no rejected alternative is taken; containment is tested before any write in the resolver, not by refusing every link, not by returning resolved paths, not by writing through `os.Root`.
- ADR-0269: no rejected alternative is taken; its ignored-file refusal and forward recovery stand, and only its symbolic-link bullet is narrowed by ADR-0306.
- ADR-0268: no rejected alternative is taken; accepting any symbolic link is the alternative it rejected, and this Spec stops accepting a link that leaves the repository while still resolving the root.
- ADR-0272: no rejected alternative is taken; paths are still compared in resolved spellings.
- ADR-0247: no rejected alternative is taken; the record stays one small file that a commit can now always hold.
- ADR-0260: no rejected alternative is taken; promotion and relink rules are unchanged, and a promotion outside the repository is refused before writing.
- ADR-0278: no rejected alternative is taken; Git Locators and the Legacy Link Rebase are untouched.
- ADR-0262: no rejected alternative is taken; the cancel keeps its record and Source Return, refused before writing when a destination leaves the repository.
- ADR-0154: no rejected alternative is taken; the override keeps its provenance and refusals.
- ADR-0299: no rejected alternative is taken; Git Locator text and resolution are untouched.
- ADR-0184: no rejected alternative is taken; both changed command surfaces are stated as transcripts the QA gate runs through the built product.
- ADR-0280: no rejected alternative is taken; task_02 changes the two guides and re-records their surface fingerprint itself.
- ADR-0256: no rejected alternative is taken; the Roundfix `archive` reference changes with the surfaces it describes, with its version raised.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the Go documentation and the Go blog.
- ADR-0240: no rejected alternative is taken; no criterion needs a Run sandbox it lacks.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and byte invariance compares with the Delivery Base.
- ADR-0265: no rejected alternative is taken; each new component has its absence proof.
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
- ADR-0301: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

Git, through the existing `git` executable, for `rev-parse --show-toplevel`,
`status`, `show` and `ls-tree` in the checkout the command runs in; the file
system, for resolving symbolic links. No test, Verification command or QA row
reads the real `~/.roundfix` or reaches the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, homes and Git configuration, resolves its
temporary directory with `filepath.EvalSymlinks` first, and creates its own
links with `os.Symlink`, so it crosses a link on Linux as on macOS. Each
refusal case asserts that the outside directory holds no new file, that the
Spec folder bytes and `git status --porcelain --untracked-files=all` are
unchanged, and the exact refusal text.

- `internal/spec/repository_path_containment_test.go` (new):
  - `TestRepositoryPathRefusesALinkThatLeavesTheRepository`, with subtests:
    an existing directory under an in-repository link to an outside
    directory; a file that does not exist yet under that link; a link to an
    outside directory spelled through a root that is itself a link to the
    checkout. Each fails with `path "<path>" is outside the repository
    "<root>"`, naming both arguments as given.
  - `TestRepositoryPathKeepsLinksThatStayInside`, with subtests: a path under
    an in-repository link to another directory inside returns the lexical
    relative path (`docs/history/specs/demo.md`, not the target's path); the
    root and the path both spelled through a link to the checkout return the
    lexical relative path (the trap: resolving only the path would refuse
    it); the root through the link and the path resolved, and the reverse,
    return the relative path.
  - `TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting`, through
    `spec.Archive` with a QA override in a committed fixture: `docs/history`
    linked outside fails with `resolve Archive Root for archive relink: path
    "<repo>/docs/history/specs" is outside the repository "<repo>"`;
    `docs/references` linked outside with `Promote:
    []string{"references/notes.md"}` fails with `resolve promotion
    destination: path "<repo>/docs/references/notes.md" is outside the
    repository "<repo>"`; neither writes anything.
  - `TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting`, through
    `spec.Cancel`: `docs/backlog` linked outside with an adopted Backlog Entry
    fails with `cancellation path "<repo>/docs/backlog/<file>" is outside
    the repository`; `docs/references` linked outside with a committed Live
    Document citing the Spec's `references/notes.md` (so the relink promotes
    it) fails with `cancellation path "<repo>/docs/references/notes.md" is
    outside the repository`; `docs/history` linked outside fails with
    `cancellation path "<repo>/docs/history/specs" is outside the
    repository`; none writes anything.
  - `TestArchiveAndCancelKeepLinksThatStayInside`: with `docs/history` a link
    to `docs/hist-real` inside, and with the repository root and Spec Root
    both spelled through a link to the checkout, `spec.Archive` and
    `spec.Cancel` succeed and write the same record and relinked bytes as a
    control without links, with the same archive date.
- `internal/cli/archive_link_escape_test.go` (new), through the command entry
  point:
  - `TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript`: the
    `docs/history` case asserts Surface Transcript 1 whole; subtests assert
    the whole refusal transcript for the Spec folder link (`path
    "<repo>/docs/specs/0001-demo" is outside the repository "<repo>"`), the
    promotion into a linked `docs/references`, and `archive 0001-demo --plan
    --offline` with a linked `docs/history` (`path
    "<repo>/docs/history/specs" is outside the repository "<repo>"`). Each
    leaves the outside directory, the Spec folder and `git status` unchanged.
  - `TestCancelRefusesALinkThatLeavesTheRepositoryTranscript`: the
    `docs/backlog` case asserts Surface Transcript 2 whole; a subtest asserts
    the `docs/history` case, refused by the citation advice with `path
    "<repo>/docs/history/specs" is outside the repository "<repo>"`.
  - `TestArchiveKeepsLinksThatStayInsideTranscript`: subtests for a
    `docs/history` link inside the repository, the checkout reached through a
    link to it (work directory spelled through the link), and an external
    Spec Root configured in `.roundfixrc.yml` through a link spelling of its
    repository; each exits 0 with the same standard output and standard error
    as its control and writes the control's record bytes.
- Declared break: none in tests. A prototype of API Contracts 1-3 on a
  scratch copy left every existing test in `internal/spec` and the
  `Archive|Cancel|SpecFolderStatus` tests of `internal/cli` passing, including
  `TestRepositoryPathComparesResolvedSpellings`,
  `TestArchiveRollsBackTheRecordWhenPromotionCannotBeWritten` and
  `TestCancelWriteFailureUsesExitOne`, whose dangling links resolve as missing
  components inside the repository. The behavior break is the one `_prd.md`
  → Declared break states.

## Build Order

1. The resolver and the package code: `RepositoryPath`, the destination
   checks in `spec.Archive` and `spec.Cancel`, and their tests (no
   dependency).
2. The command transcripts, the `archive` and `cancel` guide paragraphs and
   the Behavior Surface Record (depends on: 1, whose resolver the commands
   reach through every path they build).
3. The Roundfix skill's `archive` reference, its version and the ADR-0269
   line (depends on: 2, which settles the wording the reference repeats).
4. QA gate (depends on: 1, 2, 3).

No other Spec is a prerequisite. Specs 0290, 0291 and 0292 deliver before
this one and declare none of its files; the plan is rebased and checked again
before its Run.

## Risks & Considerations

- A repository that keeps a destination directory behind a link outside
  itself loses archive and cancel until it moves the link inside; the
  refusal names the path so the operator finds it.
- A dangling link is resolved as a missing component, so the containment
  test sees its lexical spelling. Nothing is written through it:
  `os.MkdirAll` refuses to create a directory where a link exists, and an
  existing destination link fails the vacancy check
  (`TestCancelWriteFailureUsesExitOne` holds the first today).
- The test happens before the writes, not at them; a link swapped in between
  by another process is outside the threat model, which is a repository
  layout, not an attacker racing the command.
- Resolving costs one `EvalSymlinks` walk per path; an archive resolves fewer
  than twenty.

## Glossary

None.

## Decisions

- Test containment in the resolver, so every path the two commands build is
  covered at once, and route the three destination joins through it.
- Resolve the root as well as the path, so a checkout or external Spec Root
  reached through a link is inside.
- Keep the lexical relative path when it is inside, so no accepted path
  changes its bytes.
- Reuse each step's existing refusal text; add no message.
- Record ADR-0306, which narrows ADR-0269's symbolic-link bullet.
