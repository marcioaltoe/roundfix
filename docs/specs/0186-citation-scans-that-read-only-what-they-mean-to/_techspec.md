---
spec: 0186-citation-scans-that-read-only-what-they-mean-to
prd: _prd.md
created: 2026-09-30
---

# Citation scans that read only what they mean to — Technical Spec

## Executive Summary

Two independent fixes.

- **The Spec citation walk** in `internal/speccheck/citations.go` applies its
  Task-section projection only to the Spec folder's top-level `task_*.md`
  files.
- **The Relocation Citation scan** in `internal/baseline/history_citations.go`
  opens tracked files through one `os.Root` on the repository root, instead of
  a joined absolute path.

The trade-off is accepting `os.Root`'s rule that a link staying inside the root
may be followed, instead of hand-writing a per-component `openat` walk with a
separate Windows path. Containment comes from the standard library on every
platform. The scan's existing checks keep the stronger no-follow promise,
except for a race confined to the repository.

## Project Constraints

- Identifier strategy: not applicable — no identifier, key or payload changes.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093, ADR-0094, ADR-0104, ADR-0117,
  ADR-0130, ADR-0155, ADR-0156, ADR-0168, ADR-0173 and ADR-0176 hold;
  ADR-0166 and ADR-0167 do not apply; the gate is bound
  by ADR-0080, ADR-0091, ADR-0096 and ADR-0097; this Spec adds ADR-0177.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`;
  no governed test, `go.mod` or Makefile edit. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package, command or file type.

- `internal/speccheck/citations.go`: `readSpecCitations` walks the Spec folder.
  Its projection flag becomes "the file is directly in `specDir` and its name
  matches `task_*.md`", instead of "the basename starts with `task_`". The `qa/`
  skip is unchanged.
- `internal/baseline/history_citations.go`: `relocationCitationFindings` opens
  `os.OpenRoot(root)` once per scan, before reading any file, and closes it
  when the scan returns. `readTrackedCitationFile` receives the `*os.Root` and
  passes it, with the repository-relative path, to the platform open.
- `internal/baseline/history_citations_open_unix.go` and
  `internal/baseline/history_citations_open_windows.go`:
  `openCitationFileNoFollow` becomes a function of `(*os.Root, relative string)`.
  - On Unix it calls
    `root.OpenFile(filepath.FromSlash(relative), os.O_RDONLY|syscall.O_NONBLOCK, 0)`.
    `os.Root` already adds `O_NOFOLLOW|O_CLOEXEC` and opens every intermediate
    component with `O_NOFOLLOW|O_DIRECTORY`.
  - On Windows it calls `root.Open(filepath.FromSlash(relative))`.
- The pre-open `lstatTrackedCitationPath`, the post-open uncached re-check and
  the `os.SameFile` comparisons stay exactly as they are.

## Implementation Design

### Interfaces

```go
// internal/speccheck/citations.go
// projectedTaskFile reports whether path is one of the Spec's own Task
// files: directly inside specDir and named task_*.md.
func projectedTaskFile(specDir, path string) bool

// internal/baseline (both platform files)
func openCitationFileNoFollow(root *os.Root, relative string) (*os.File, error)
```

`readTrackedCitationFile` gains a `root *os.Root` parameter. It is
unexported, so no exported signature changes.

### Data Models

None. No schema, record or payload changes.

### API Contracts

None. `roundfix spec check` and the Baseline plan commands keep every output,
flag and exit code. The two fixes change which bytes an internal scan reads,
not what any command prints for authored or tracked content.

## Coverage Map

- Goal 1, Core Feature 1 → `projectedTaskFile`, `readSpecCitations`; Testing
  Approach 1.
- Goal 2, Core Feature 2 → `os.Root` in `relocationCitationFindings`,
  `openCitationFileNoFollow`; Testing Approach 2.
- Goal 3 → Testing Approach 1–3.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.

## Integration Points

None external. `os.Root` is in the standard library of the module's Go
version (`go 1.26` in `go.mod`), so no dependency changes.

## Testing Approach

1. **Projection scope.** New `internal/speccheck/citation_projection_scope_test.go`
   over plain temporary Spec folders:
   - an unlisted ADR cited only in `references/task_example.md`'s
     `## Result` opens `SC-ADR-UNLISTED`;
   - the same citation in a top-level `task_01.md` `## Result` opens none;
   - a citation in a nested directory's `task_*.md` Daemon section still opens
     the finding.

   The existing `TestASpecCitationInATaskResultIsNotAnObligation` and
   `TestAnAuthoredSpecCitationStillMustBeListed` stay green unchanged.
2. **Root containment.** New `internal/baseline/history_citations_root_test.go`
   over a real temporary Git repository. A package-level test hook,
   `citationBeforeOpen func(relative string)`, is nil in production and called
   after the pre-open check and before the open.
   - The race test sets the hook to replace a tracked directory with a
     symbolic link to an outside directory that holds a file citing a
     relocated path under the same relative name. It asserts that no finding
     names that citation and that the outside file was never opened.
   - A second test proves that an unraced repository still reports its
     citations.

   The hook is declared in `history_citations.go`.
3. **Unchanged suite and Windows build.** The existing
   `TestRelocationCitations*` tests stay green unchanged, and
   `GOOS=windows GOARCH=amd64 go build -buildvcs=false -o /dev/null ./cmd/roundfix`
   succeeds. Only non-test code is built for Windows, because an existing
   Unix-only FIFO test in `internal/baseline/repository_test.go` does not
   compile there.

## Build Order

1. The projection applies to top-level Task files only, task_01 (depends on:
   none).
2. The scan reads through a repository root, task_02 (depends on: none).
3. Terminal QA, task_03 (depends on: 1, 2).

## Risks & Considerations

- **Unix non-blocking open through `os.Root`.** `Root.OpenFile` passes caller
  flags through, and `os` records a non-blocking flag
  (`unix.HasNonblockFlag`), so a FIFO swapped in still never blocks. The
  existing FIFO test must stay green.
- **Windows.** `Root.Open` opens components relative to held handles and
  refuses to leave the root. `FILE_FLAG_OPEN_REPARSE_POINT` is no longer
  passed explicitly. The post-open `os.SameFile` check still refuses a final
  component that was a link when checked.
- **Hook discipline.** The test hook is a single package variable, nil in
  production, and never assigned outside tests.

## Decisions

- Scope the projection by location, not by parsing the Task Graph.
- Read through `os.Root`, keeping the existing checks. See ADR-0177.
