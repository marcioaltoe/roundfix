---
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
prd: _prd.md
created: 2026-11-13
---

# A nested carrier Git does not list raises no conflict — Technical Spec

## Executive Summary

The Baseline inventory walk (`inventoryBuilder.walk` in
`internal/baseline/repository.go`) asks Git, on the first nested `AGENTS.md`
or `CLAUDE.md` it meets, for the carriers that belong to the repository:
`git ls-files -z --cached --others --exclude-standard --
':(glob)**/AGENTS.md' ':(glob)**/CLAUDE.md'` at the repository root. A
nested carrier that the output lists is inventoried exactly as today; any
other one is left out, with no preimage, carrier, warning or counted bytes.
One Git rule covers ignored paths, linked worktrees, nested repositories and
submodules. The trade-off is one more read-only Git command per inspection
that meets a nested carrier, and one more failure mode (a failed listing
stops the inspection), in exchange for never second-guessing Git's own
ignore and checkout boundaries (ADR-0399).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0399 (new) and ADR-0070, which it
  amends, apply; ADR-0372 applies unchanged; ADR-0351, ADR-0330, ADR-0074,
  ADR-0078 and ADR-0100 do not (see ADR Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BE batch 2 answer of 2026-11-13 quoted in the authorization record;
  the Governed Paths are the Roundfix skill files task_02 changes and
  `internal/baseline/repository_test.go`, which task_01 changes. No Makefile,
  workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0390-a-nested-carrier-git-does-not-list-raises-no-conflict/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `internal/baseline/repository_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `2e7a2a7a`:

1. `inspectRepositorySnapshot(root, request)` walks the working tree with
   `fs.WalkDir`. `inventoryPathIgnored` skips any path with a `.git`, `.hg`,
   `.svn`, `.venv`, `__pycache__`, `node_modules`, `vendor` or `venv`
   component, and `skills/` and `.agents/skills/`. Nothing else is skipped
   and no Git command is run: the walk does not read `.gitignore`.
2. `isBoundedCarrierPath` selects a path whose base name is `AGENTS.md` or
   `CLAUDE.md`, or any path under `docs/agents/`. A path with a `/` gets scope
   `nested`; every nested carrier outside the repository-rule carriers raises
   `baseline.inventory.nested-carrier-conflict`, later filtered for
   repository-authored guides by ADR-0372's classification.
3. A linked worktree's `.git` is a file; `inventoryPathIgnored` skips the file
   but not its directory, so the walk enters the worktree and inventories
   its carriers. The same holds for a nested repository and a submodule.
4. `InspectRepository` runs exactly two Git commands (`rev-parse ...` and
   `rev-list --max-parents=0 HEAD`); `TestRepositoryInspectionUsesNarrowReadOnlyGitCommands`
   pins them for a repository without nested carriers.
5. `ExecGitRunner.RunGit` returns `strings.TrimSpace` of standard output.
6. In a temporary repository, `git ls-files -z --cached --others
   --exclude-standard -- ':(glob)**/AGENTS.md' ':(glob)**/CLAUDE.md'` printed
   the untracked entries first, then the tracked ones; it omitted the
   ignored `.storage/` tree, the files of linked worktrees (ignored or not),
   a nested repository and a submodule, and listed a tracked file under an
   ignored directory, a symbolic link named `CLAUDE.md` and the files of a
   directory whose `.git` file pointed to a missing repository. `git
   check-ignore` exited 128 for a path inside the submodule.

## System Architecture

The change stays in the inventory of `internal/baseline/repository.go`. No
package, file or command is added; every caller of
`inspectRepositorySnapshot` (`InspectRepository` with its runner,
`buildPlanWithCatalog`, `validatePlanRepositoryState`,
`inspectApplyPreimage`, and `verifyAppliedPlanState`, which the two
transaction constructors in `internal/baseline/transaction.go` wrap in a
closure over their context) passes its context, so the
plan, the apply re-inspection and the post-apply verification all apply the
same rule.

## Implementation Design

### Interfaces

```go
// ExecGitRunner keeps RunGit (trimmed) and gains:
func (ExecGitRunner) RunGitRaw(ctx context.Context, workDir string, args ...string) ([]byte, error)

type rawGitRunner interface {
	RunGitRaw(context.Context, string, ...string) ([]byte, error)
}

func listRepositoryCarriers(ctx context.Context, root string, runner GitRunner) (map[string]struct{}, error)
func inspectRepositorySnapshot(ctx context.Context, root string, request InventoryRequest, runner GitRunner) (RepositorySnapshot, error)
func verifyAppliedPlanState(ctx context.Context, root string, document PlanDocument) error
func (builder *inventoryBuilder) inventoryCarrier(normalized string, entry fs.DirEntry)
```

### Data Models

`inventoryBuilder` gains `listCarriers func() (map[string]struct{}, error)`,
set by `inspectRepositorySnapshot` to call `listRepositoryCarriers` with its
context, root and runner, and `listed map[string]struct{}`, nil until the
first nested carrier. The snapshot schema and its digest domain are
unchanged.

### API Contracts

1. API Contract: `ExecGitRunner.RunGitRaw` runs Git exactly as `RunGit` does
   today (same `-C`, `-c core.fsmonitor=false`, environment and error text)
   and returns standard output byte for byte. `RunGit` returns
   `strings.TrimSpace` of `RunGitRaw`'s output; its behavior is unchanged.
2. API Contract: `listRepositoryCarriers` runs, with `runner` (or
   `ExecGitRunner{}` when nil) at `root`, the arguments `ls-files`, `-z`,
   `--cached`, `--others`, `--exclude-standard`, `--`,
   `:(glob)**/AGENTS.md`, `:(glob)**/CLAUDE.md`, in that order. It reads the
   output through `RunGitRaw` when the runner implements `rawGitRunner`,
   otherwise through `RunGit`, splits it on NUL and keeps every non-empty
   entry without trimming. A failure returns
   `list repository instruction carriers: %w`.
3. API Contract: in `inventoryBuilder.walk`, after the existing skips and the
   `isBoundedCarrierPath` test, a path that contains `/` and does not start
   with `docs/agents/` is a nested candidate. On the first candidate the walk
   calls `listCarriers` once and keeps the result; a listing error is
   returned from the walk, so `inspectRepositorySnapshot` fails with `walk
   bounded carriers: list repository instruction carriers: <git error>` and
   `InspectRepository` with `inspect bounded repository state: walk bounded
   carriers: list repository instruction carriers: <git error>`. A candidate
   the result does not list is left out (`fs.SkipDir` for a directory,
   otherwise `nil`). Every other path, and a listed candidate, goes to
   `inventoryCarrier`, which holds today's body from `entry.Info()` to the
   warning, unchanged, so the walk order and every finding stay as they are.
4. API Contract: an inspection whose walk meets no nested candidate runs no
   new Git command.
5. API Contract: Text A, a new paragraph placed right after the paragraph
   that ends "A nested `AGENTS.md` or `CLAUDE.md`, an unmarked guide at a
   managed path and a guide with an unknown marker still warn." in
   `docs/user-guide/commands/baseline.md` and in both copies of the Roundfix
   skill's `references/baseline.md`, word for word:

   ```markdown
   A nested `AGENTS.md` or `CLAUDE.md` enters the inventory only when `git ls-files --cached --others --exclude-standard` lists it: a tracked one, or an untracked one Git does not ignore. A copy in a path Git ignores, and a carrier inside a linked worktree, a nested repository or a submodule, is not a file of the repository: it raises no `baseline.inventory.nested-carrier-conflict` and does not enter the Plan Digest.
   ```

6. API Contract: Text B, the amendment notice appended as the last paragraph
   of ADR-0070, whose `updated_at` becomes the time the Task writes it, word
   for word:

   ```markdown
   **Amended by [ADR-0399](0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md) (2026-11-13).** A nested `AGENTS.md` or `CLAUDE.md` is inventoried only when `git ls-files --cached --others --exclude-standard` lists it; a carrier in a path Git ignores, or inside a linked worktree, a nested repository or a submodule, is not a carrier of this repository. Root carriers, their backups and the warning-only handling of the nested carriers Git lists still hold.
   ```

### Surface Transcripts

None. No command's text changes: `baseline update`, `baseline plan` and
`baseline apply` print fewer `baseline.inventory.nested-carrier-conflict`
lines when Git does not list a nested carrier, in the existing format, and a
failed listing goes through the existing Git-failure path (exit 2).

### Invariants

1. A repository whose nested `AGENTS.md` and `CLAUDE.md` Git all lists has
   the same snapshot, warnings and Plan Digest as at `2e7a2a7a`.
2. Adding an ignored copy, a linked worktree, a nested repository or a
   submodule changes neither the snapshot digest nor the warnings.
3. Root carriers and every path under `docs/agents/` never depend on the
   listing.

### Declared breaks

- An untracked nested carrier that `.gitignore`, `info/exclude` or the
  user's `core.excludesFile` ignores no longer warns.
- A carrier inside a linked worktree, a nested repository or a submodule no
  longer warns.
- An untracked directory or special file named `AGENTS.md` or `CLAUDE.md`
  below the root is no longer an opaque nested carrier with a warning, since
  Git lists files only.
- A failed `git ls-files` fails an inspection that meets a nested carrier;
  before, no Git command ran there.

## Coverage Map

- Goal 1, Goal 2, Core Feature 1, Success Metric 1 → API Contracts 2 and 3,
  Invariant 2, `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists`.
- Goal 3, Success Metric 1 → API Contract 3, Invariants 1 and 3,
  `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists`,
  `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths`.
- Core Feature 2, Success Metric 2 → API Contracts 1, 2 and 4,
  `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments`.
- Success Metric 3 → API Contract 3 through `baseline update`, the QA gate's
  clone row.
- Goal 4, Core Feature 3, Success Metric 4 → API Contracts 5 and 6.

Deviations from the PRD: None.

## Absence Evidence

1. `listRepositoryCarriers` and `RunGitRaw` — command: `grep -rn 'listRepositoryCarriers\|RunGitRaw\|exclude-standard' internal/baseline`; result: no match at `2e7a2a7a`.

## ADR Alternatives

- ADR-0399: the design is the decision; it takes none of the five options
  the record rejects.
- ADR-0070: it limited automatic mutation to root carriers and kept nested
  conflicts as warnings, rejecting remediation across the repository; both
  hold, and the audit's scope narrows to the carriers Git lists, recorded by
  the amendment notice task_02 writes.
- ADR-0372: it rejected silencing `nested-carrier-conflict` for every guide
  under `docs/agents/`; this Spec never applies the listing there, so that
  classification stays as it is.
- ADR-0351: not applicable — it governs QA Report links; ADR-0399 only
  follows its precedent of asking `git ls-files --cached --others
  --exclude-standard` instead of reading ignore files, and this Spec changes
  nothing it decides.
- ADR-0330: not applicable — it moved cross-project lessons to ai-memory, a
  sentence ADR-0372 cites and this Spec does not touch.
- ADR-0074: not applicable — it splits repository rules between root and
  semantic owners, which no nested carrier left out of the inventory changes.
- ADR-0078: not applicable — it moves confirmed root rules to semantic
  owners; its statement that the warning-only handling of nested carriers
  remains active still holds for every carrier Git lists, and no root rule
  changes.
- ADR-0100: not applicable — it makes a managed refresh prove preservation of
  root carriers through preimages; root carriers never depend on the
  listing.

## Integration Points

- The installed `git`, the same dependency `InspectRepository` already has;
  a Linux Lane test runs it as uid 65534 over temporary repositories it
  creates.
- `baseline update`, `baseline plan`, `baseline apply`, the interactive
  `baseline` command and every other caller of `InspectRepository` see the
  rule through the inventory.

## Testing Approach

New tests in `internal/baseline/repository_test.go`, each calling
`t.Parallel()` first:

- `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists`: a repository from
  `newInspectionRepository` whose local config sets `core.excludesFile` to an
  empty file under the test's temporary directory (so the host's global
  excludes play no part), with `.gitignore` holding `.storage/` and
  `forced/`, a committed root `AGENTS.md` and `packages/api/AGENTS.md`, a
  committed `forced/AGENTS.md` added with `git add --force`, and untracked
  `packages/web/CLAUDE.md` and ` lead/AGENTS.md` (a leading space, so Git
  lists it first). The first `InspectRepository` is taken there. Then: a
  linked worktree at `.storage/simulation/wt-ignored`, an ignored
  `.storage/copy/CLAUDE.md`, a linked worktree at `trees/wt-visible`, a
  nested repository `nested/` with an `AGENTS.md`, and a submodule at
  `modules/sub` added with `git -c protocol.file.allow=always submodule add`
  from a second local repository holding an `AGENTS.md`. Both inspections
  must list exactly the nested carriers and nested-carrier-conflict paths
  ` lead/AGENTS.md`, `forced/AGENTS.md`, `packages/api/AGENTS.md`,
  `packages/web/CLAUDE.md`, in order, and equal snapshot digests.
- `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments`: a
  `t.TempDir()` root with `packages/api/AGENTS.md` and
  `packages/web/CLAUDE.md`, an `inspectionGitRunnerFunc` answering the two
  identity commands and the listing of API Contract 2 with
  `packages/api/AGENTS.md\x00`: exactly three Git calls and the carriers
  `[packages/api/AGENTS.md]`; then a runner whose listing fails with
  `listing failed`: the error of API Contract 3, word for word.
- The existing `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths`,
  `TestRepositoryInspectionUsesNarrowReadOnlyGitCommands`, the inventory and
  classification tests of `./internal/baseline` and the `Baseline` tests of
  `./internal/cli` pass unchanged; the only existing test edit is the
  `inspectRepositorySnapshot` call of the over-budget subtest, which gains
  `context.Background()` and `nil`.

Sabotages, each on a temporary copy of the audited head made with `git
archive` and committed in a fresh Git repository, each observed by the named
test (all run by the author on 2026-11-13 against the prototype):

1. Drop `--exclude-standard` →
   `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists` (`.storage/copy/CLAUDE.md`
   appears) and `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments`.
2. Read through `RunGit` even when `RunGitRaw` exists →
   `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists` (` lead/AGENTS.md`
   disappears).
3. Drop `--cached` → `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths`
   and `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists` (the tracked
   carriers disappear).
4. Inventory a candidate the listing omits →
   `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists` (the worktree,
   nested and submodule carriers appear) and
   `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments`.
5. Apply the listing to `docs/agents/` paths too →
   `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths`
   (`docs/agents/backend.md` disappears).
6. List again for every candidate →
   `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments` (four calls).

## Build Order

1. The inventory rule, `RunGitRaw`, the context threading and the two tests
   (task_01).
2. The `baseline` guide, both copies of the Roundfix `baseline` reference,
   the skill version, the Behavior Surface Record and the ADR-0070 notice
   (task_02) (depends on: 1, so the text describes merged behavior).
3. The QA gate (task_03) (depends on: 2).

No other Spec is required. Spec 0389 delivers first and changes the
Roundfix skill's `deliver` reference; the shared version fields and records
are named in `_prd.md` → Prerequisites. Each Task works on the tree it runs
on, cites functions rather than lines, and task_02 raises the skill version
with the record command.

## Risks & Considerations

- A user's `core.excludesFile` can leave out an untracked nested carrier on
  one machine and not on another; the Plan Digest already depends on
  untracked files, and `git status` shows the same difference.
- The listing walks untracked, not ignored directories like `git status`;
  on this repository it took about 50 ms.
- A Git that ignores the `:(glob)` pathspec magic would list nothing and
  leave every nested carrier out; Git has supported it since 1.8.3, and the
  tests run the installed Git.

## Glossary

None.

## Decisions

- Git lists the repository's nested carriers in one call, read untrimmed.
  See ADR-0399.
- The listing runs lazily inside the walk so the walk order and the byte
  budget order stay as they are. See ADR-0399.
- A failed listing stops the inspection. See ADR-0399.
- Nothing reports a carrier left out. See ADR-0399.
- ADR-0070 receives an amendment notice. See ADR-0399 and ADR-0070.
