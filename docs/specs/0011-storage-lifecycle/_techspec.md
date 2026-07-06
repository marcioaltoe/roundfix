---
spec: 0011-storage-lifecycle
prd: _prd.md
created: 2026-07-06
---

# Review Artifacts, Run Logs, and Spec Archiving — Technical Spec

## Executive Summary

Three storage/lifecycle fixes plus two output-shaping cleanups, all local in
blast radius. The consequential decision is where review artifacts default: an
in-repo review root resolved by hierarchy (explicit config → spec-associated →
default), accepted because artifacts that version with the feature beat a loose
user-scoped pile, and Run Worktree isolation already keeps them out of Batch
commits (ADR-0029). Opt-in agent logs are a gate in front of an existing writer,
justified because the Run Event Journal already stores every payload (ADR-0030).
Spec archiving becomes a first-class Archive Command so the precondition
(all-completed + QA-passed) is enforced deterministically rather than trusted to
a skill. Title hygiene and poll dedup are pure output shaping.

## System Architecture

- `internal/config` — a review-artifact-root resolver and a new opt-in
  `logs.agent` (User/Project) key. The existing `artifact_dir` stays the
  explicit override at the top of the hierarchy.
- `internal/cli` — the review-artifact root today is derived as
  `<artifact_dir>/reviews/pr-<n>/round-*`; it moves behind a resolver that
  chooses among explicit config, `docs/specs/<slug>/reviews/`, and
  `docs/specs/_reviews/pr-<n>/`. Spec association reads the newest
  `Roundfix-Spec` PR-head commit trailer or an explicit spec selector. A new
  `archive` command reads the Spec's Task Graph and QA verdict, then moves the
  folder.
- `internal/agent` (or wherever per-Batch logs are written) — the log-file
  writer becomes conditional on the `logs.agent` key; the Run Event Journal and
  the Detached Run console log are untouched.
- `internal/reviewsource/coderabbit` — the Review Issue title derivation strips
  markup/emoji.
- The watch status-poll loop in `internal/cli` — dedups identical poll lines.
- No store schema, journal, or TUI changes.

## Implementation Design

### Interfaces

```go
// internal/config (or internal/cli) — review artifact root resolution.
type ReviewArtifactContext struct {
    ExplicitArtifactDir string // set → wins (ADR-0003's survivor)
    RepoRoot            string
    SpecSlug            string // from --spec or the newest Roundfix-Spec trailer; "" if none
    PRNumber           int
}

// Resolve returns the directory under which reviews/pr-<n>/round-* is written.
func ResolveReviewRoot(ReviewArtifactContext) (string, error)
// explicit  → ExplicitArtifactDir
// spec      → <RepoRoot>/docs/specs/<slug>/reviews
// default   → <RepoRoot>/docs/specs/_reviews/pr-<n>
```

```go
// Archive Command contract.
type ArchiveRequest struct{ Slug string }
// Preconditions checked before any move:
//  - every task_NN.md in _tasks.md has status: completed
//  - a QA report under qa/ carries a passing verdict
// On pass: stamp archive metadata (archived date, source slug) and move
// docs/specs/<slug>/ → docs/specs/_archived/<slug>/. On failure: return an
// error naming the unmet condition; move nothing.
```

### Review artifact location (ADR-0029)

The current join at the resolve/watch call sites is replaced by
`ResolveReviewRoot`. Spec association: an explicit `--spec <slug>` selector wins
over trailer discovery; otherwise the newest `Roundfix-Spec: <slug>` trailer on
the PR head commit supplies the slug. When neither the explicit `artifact_dir`
nor a spec is present, the root is `docs/specs/_reviews/pr-<n>/`. Roundfix never
`git add`s or `.gitignore`s the tree — placement only. Run Worktree isolation
already excludes these paths from Batch commits, so no clean-tree Preflight
regression is possible.

### Opt-in agent logs (ADR-0030)

The per-Batch agent log-file writer is wrapped in a guard reading the
`logs.agent` config key (default off). When off, no `*.log` files are created
under the Run's artifact area; the Run Event Journal still records every payload
(ADR-0008). The Detached Run console log (ADR-0028) is a separate path and stays
unconditional.

### Archive Command

`roundfix archive <slug>` (support command, non-interactive): parse
`_tasks.md`, read each task file's `status`, read the QA verdict from the newest
report under `qa/`; if all Tasks are `completed` and the verdict passes, write
archive metadata and move the folder to `docs/specs/_archived/<slug>/`; else
exit non-zero naming the first unmet condition. It creates no Run and never
pushes — folder move plus metadata only.

### Title hygiene and poll dedup

- Title derivation in the CodeRabbit source strips markdown table pipes,
  heading/emoji markup, and surrounding whitespace before the title is stored,
  so Work Queue rows and issue files read as plain text.
- The watch status-poll writer keeps the last emitted poll line and prints only
  when the new line differs, collapsing identical intervals.
- The merge-readiness `missing` branch appends the documentation expectation
  that explains the state to its stderr note.

## Coverage Map

- Stories 1-3 → `ResolveReviewRoot` + spec association (ADR-0029)
- Stories 4-5 → `logs.agent` guard (ADR-0030)
- Story 6 → Archive Command
- Story 7 → CodeRabbit title derivation
- Story 8 → status-poll dedup

## Integration Points

Git only, read-only: the `Roundfix-Spec` trailer is read from the PR head commit
(existing commit-reading path). No new external systems. The Archive Command
touches the filesystem within the repository (`docs/specs/`).

## Testing Approach

- `ResolveReviewRoot`: table tests over the three branches (explicit, spec via
  flag, spec via trailer, spec-less default) asserting the resolved path.
- Agent logs: a test asserts no log files are written with the key off and files
  appear with it on, while the journal is populated in both cases.
- Archive Command: buffer-captured CLI tests over a temp spec folder —
  all-completed + QA pass moves and stamps; an incomplete task or failing/absent
  QA verdict refuses with the reason; the moved folder lands under `_archived/`.
- Title hygiene: table tests over CodeRabbit table-fragment inputs → plain
  titles.
- Poll dedup: drive the poll writer with repeated identical statuses; assert one
  line.

## Build Order

1. `ResolveReviewRoot` resolver + spec association, wired into resolve/watch
   artifact paths (ADR-0029) (no deps)
2. Opt-in `logs.agent` gate around the per-Batch log writer (ADR-0030) (no deps)
3. Archive Command with the all-completed + QA-passed precondition (no deps)
4. Title hygiene + status-poll dedup + merge-readiness docs note (no deps)
5. Docs and skill sync (depends on: 1, 2, 3, 4)

## Risks & Considerations

- Spec association must be conservative: a wrong slug would write artifacts into
  an unrelated Spec's folder — an explicit `--spec` and a validated trailer slug
  (folder must exist) are the only sources.
- The Archive Command move must be atomic enough that a failure mid-move leaves
  the Spec either wholly in place or wholly archived; validate preconditions
  before touching the filesystem.
- The `logs.agent` default flip is a behavior change for anyone who parsed those
  files — the deprecation/None path is documented; the journal remains the
  supported record.

## Decisions

- Review artifact root resolves explicit → spec-associated → in-repo default;
  never committed or ignored by Roundfix. See ADR-0029 (supersedes ADR-0003's
  default).
- Per-Batch agent logs gate behind `logs.agent`; Detached Run console log stays
  unconditional. See ADR-0030.
- Archiving is the Archive Command enforcing all-completed + QA-passed; glossary
  gains Archive Command.
