---
spec: 0018-external-spec-root
prd: _prd.md
created: 2026-07-07
---

# External Spec Root — Technical Spec

## Executive Summary

The Spec Root becomes a resolved-once absolute path threaded through every
Spec consumer, replacing the assumption that Spec artifacts live under the
execution root's `docs/specs`. The primary trade-off is plumbing over
inference: threading an explicit root through implement, settle, attach, and
the task engine touches more call sites than following symlinks at the
staging boundary would, but it makes worktree behavior deterministic (a
relative root never re-resolves against a worktree) and makes the commit rule
("external artifacts are never staged") a property of configuration rather
than of git error handling. Field evidence (conexus, tax-poc) shows the
failure lands after verification passes — the daemon dies staging the task
file and the Run stays dead — so the commit boundary must degrade to settling
without the artifact, never to killing the Run.

## System Architecture

No new packages. The change threads one value — the resolved Spec Root —
through existing modules:

- `internal/config` — new `specs.root` key (Project > User > built-in
  `docs/specs`), validation, and resolution to an absolute path against the
  user checkout's repository root. Also owns the "external" predicate: the
  resolved root (after symlink evaluation) lies outside the repository
  working tree.
- `internal/spec` — `Load` and the active-Spec listing take the resolved
  Spec Root instead of deriving `<gitRoot>/docs/specs`. Task file paths in
  the Graph become resolvable against the Spec Root, not the execution
  `WorkDir`.
- `internal/cli` — implement, settle, archive, attach, and Interactive Input
  resolve the Spec Root once at command start and pass it down; Run startup
  reports it on stderr when it is not the default.
- `internal/daemon` — the task engine reads/writes task files and QA Reports
  through the Spec Root; the commit path drops external and symlink-crossing
  paths from staging with a journaled warning, and settles without a commit
  when nothing stageable changed.
- `internal/config` review artifact resolution — the Spec-associated review
  artifact paths follow the same root.

## Implementation Design

### Interfaces

Config resolution (in `internal/config`):

```go
// SpecsRoot is the resolved Spec Root for one command invocation.
type SpecsRoot struct {
    Path     string // absolute
    External bool   // outside the repository working tree after symlink evaluation
}

// ResolveSpecsRoot resolves specs.root (default docs/specs) against the
// user checkout's repository root and validates it is an existing directory.
func ResolveSpecsRoot(loaded Loaded, repoRoot string) (SpecsRoot, error)
```

Spec loading (in `internal/spec`, signature shift):

```go
// Load reads a Spec's Task Graph from the Spec Root (was: from gitRoot/docs/specs).
func Load(specsRoot string, slug string) (*Graph, error)

// Task.File stays relative to the Spec Root; consumers join it with the
// resolved root, never with the execution WorkDir.
```

Commit staging guard (in `internal/daemon`):

```go
// filterStageablePaths drops paths that are external to the repository or
// cross a symbolic link, returning kept paths and one journaled reason per
// dropped path.
func filterStageablePaths(workDir string, paths []string) (kept []string, dropped []droppedPath)
```

### Data Models

No schema change. The Run row is untouched: attach and settle re-resolve the
Spec Root from config in the same repository, which yields the same absolute
path the owning command used. Config gains:

```yaml
specs:
  # Directory holding Spec folders; relative values resolve against the
  # repository root. Default keeps today's layout.
  root: "docs/specs"
```

### API Contracts

- `specs.root` — string, default `docs/specs`. Unknown-key strict validation
  is unchanged; a configured root that does not exist or is not a directory
  fails Preflight Validation with the resolved path in the message.
- Run startup with a non-default root prints one stderr line naming the
  resolved Spec Root.
- Per-Task commits in external-root repositories contain only repository
  paths. The journal records one Daemon event per dropped artifact path with
  a reason (`external to repository` or `crosses a symbolic link`), and the
  progress stream prints a warning shaped like
  `roundfix: task file <path> kept outside the repository; committed without it`.
- A Task whose only changes are external artifacts settles `completed`
  without a commit; the report line and outcome contract are unchanged.
- Exit codes, report shapes, and default-layout behavior: byte-for-byte
  unchanged.

## Coverage Map

- Story 1 (symlinked layout completes and commits) → `filterStageablePaths`,
  task engine commit path
- Story 2 (explicit external root, worktree-stable) → `ResolveSpecsRoot`,
  `spec.Load` signature shift, resolve-once threading
- Story 3 (code commits carry only code) → external predicate + staging drop
- Core Feature 6 (validation) → `ResolveSpecsRoot` error path
- PRD User Experience (startup report, commit warnings) → CLI startup line,
  daemon progress warnings

## Integration Points

None external to the machine. The knowledge workspace repository is never
git-invoked; Roundfix only reads and writes files inside the resolved root.

## Testing Approach

Existing seams:

- `internal/config` table tests: default, relative, absolute, missing-dir
  failure, external predicate via a temp dir outside a temp repo root, and a
  symlinked root.
- `internal/spec` tests: `Load` against an external temp root; task file
  paths resolve against the root.
- `internal/daemon` tests: `filterStageablePaths` with symlink-crossing and
  external paths; task engine settles a Task with only-external changes
  without a commit; commit succeeds with mixed paths staging only internal
  ones.
- `internal/cli` buffer-captured tests: implement with a configured external
  root end-to-end in a temp git repo (fake agent), startup report line,
  settle and archive through the same root.

## Build Order

1. `specs.root` config key, `ResolveSpecsRoot`, external predicate, and
   config tests.
2. Spec Root threading: `spec.Load` signature shift and every consumer
   (implement, settle, archive, attach, Interactive Input listing, task
   engine reads/writes, QA Reports, review artifact resolution), with the
   resolve-once rule and startup report (depends on: 1).
3. Commit boundary: `filterStageablePaths`, journaled warnings,
   settle-without-commit when nothing stageable changed, daemon tests
   (depends on: 1, 2).
4. Docs and skill sync: README Config and Command Boundaries, usage guide,
   roundfix SKILL.md, `make skills-sync` (depends on: 2, 3).

## Risks & Considerations

- `spec.Load`'s signature shift touches every caller; the compiler enforces
  completeness, and the CLI tests pin behavior per command.
- Archive moves directories inside the external root; `os.Rename` stays
  within one filesystem there. Cross-device roots would need a copy fallback
  — out of scope, called out in the archive error message if hit.
- The settle-without-commit path must still publish the Task settled event;
  skipping the commit must not skip the journal.
- tax-poc and conexus both carry interim git shims; after this ships they
  must be removed or they will mask regressions.

## Decisions

- Resolve once, thread explicitly — never re-derive the root from an
  execution `WorkDir`. See ADR-0035.
- External artifacts never stage into code-repository commits; the symlink
  guard applies unconditionally. See ADR-0035.
- No per-Spec or multi-root support — one root per repository (PRD
  Non-Goals).
