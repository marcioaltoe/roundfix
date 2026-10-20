---
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
prd: _prd.md
created: 2026-10-20
---

# GC sanitize reaches the default Artifact Roots — Technical Spec

## Executive Summary

Two code slices and one documentation slice, with no new package, command,
flag, configuration key or record field.

1. The storage report's resolver becomes the one exported resolver of a Run's
   Artifact Root, `store.ResolveRunArtifactRoot`, in
   `internal/store/journal.go`. The storage report keeps its output, and Run
   Retention (`internal/cli/run_retention.go`) drops its own copy and calls it.
2. `store.DiscoverArtifactRoots` groups Runs by the root the resolver returns,
   and `roundfix gc sanitize` (`internal/cli/gc.go`) classifies each resolved
   root with the proofs it already had. A Run whose root cannot be resolved
   keeps its whole group. A directory proven only by its Run's absence is
   checked again just before it is removed. The dry run names each resolved
   root, the recorded values behind it and what it would reclaim.
3. The `gc` guide, the Roundfix `storage` reference and the glossary describe
   the resolved root.

The trade-offs accepted:

- Sanitation resolves at read time, so it also reaches every Run already
  recorded with an empty root, at the cost of reading each Run's repository
  metadata and one `lstat` per Run.
- A blank recorded root resolves to the default root even for a Run that wrote
  under a root configured at the time. The default-root proof and the per-Run
  directory scan keep that safe: such a Run has no directory under the default
  root, and nothing outside `runs/` is touched.
- The recheck costs one Run Database read per directory proven by absence,
  only with `--apply`.

The plan assumes the code of Spec `reports-that-count-what-they-say` on the
default branch (`_prd.md` → Prerequisites); "Current behavior" names the
symbols it reuses.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  field is minted; Run IDs, repository IDs and Artifact Root paths keep their
  derivation. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; sanitation reads the Run Database, Git
  metadata and directories, and every test uses temporary homes and
  repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 and 2 follow ADR-0301
  ("GC sanitation classifies the Artifact Root each Run wrote to"), ADR-0255
  ("Run Retention removes terminal Runs whole and compacts incrementally"),
  ADR-0296 ("The storage report reconciles every page of one Storage
  Snapshot") and ADR-0033 ("The Run Event Journal is pruned by a retention window").
  Every step follows ADR-0259, ADR-0271, ADR-0285, ADR-0189, ADR-0256,
  ADR-0280, ADR-0104, ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named
  with its title in `_prd.md` → Project Constraints. Set no obligation here,
  because this Spec changes none of their rules: ADR-0171, ADR-0172,
  ADR-0290, ADR-0167, ADR-0184, ADR-0229, ADR-0233, ADR-0250, ADR-0266, ADR-0281, ADR-0282, ADR-0286, ADR-0295, ADR-0298, ADR-0197, ADR-0237, ADR-0272, ADR-0278, ADR-0299, ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0277, ADR-0279, ADR-0165, ADR-0251, ADR-0262, ADR-0269, ADR-0169, ADR-0196, each named in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AL approval of 2026-10-20 ("Mesmas regras do AK (Recomendado)"), carrying
  the cycle AJ grant "Mexer em Governed Paths declarados", and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0289-gc-sanitize-reaches-the-default-artifact-roots/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/storage.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `7a678ddb` (v0.97.0) and at `1e81b9d7`, the
head of the final Run branch of Spec `reports-that-count-what-they-say`, whose
Task 02 commit is `151468d6`. Line numbers are those of `7a678ddb` unless
marked.

- Sanitation (`internal/cli/gc.go`).
  - `runGCSanitation` (`:489`) opens the reader, calls
    `store.DiscoverArtifactRoots(ctx, runStore)` (`:510`), builds `allRunIDs`
    from every group's Runs, classifies each group with
    `classifyGCSanitationRoot` (`:561`) and lists candidates with
    `gcSanitationArtifactDirs` (`:685`) for `active` and `orphaned` roots
    only. With `--apply` it removes the candidates after the whole scan, with
    no second read of the Run Database.
  - `DiscoverArtifactRoots` (`internal/store/journal.go:1583`) groups Runs by
    the trimmed recorded `artifact_dir`, so every Run that recorded nothing
    falls into one group `""` across all repositories.
  - `classifyGCSanitationRoot`'s first proof refuses a path that is blank,
    relative or not clean: `preserved because recorded Artifact Root "" is not
    a clean absolute path`. The later proofs are the Run ID check, containment
    in Roundfix Home, not Roundfix Home itself, existence, physical root,
    physical `runs/`, and the default-root proof: the root must equal each
    Run's checkout default or repository-key default, otherwise `overridden`.
  - `gcSanitationArtifactDirs` keeps a directory whose Run is recorded under
    another group with `Run <id> is recorded under Artifact Root "<path>"`.
  - `printGCSanitationReport` (`:1153`) prints, per root, `Artifact Root:`,
    `Classification:`, `Evidence:`, `would remove:`/`removed:` and
    `preserved:` lines, then the totals `Directories reclaimable:` and
    `Artifact bytes reclaimable:` (or `removed:`/`reclaimed:`).
  - Implement Runs record no root (`internal/cli/implement.go:413` passes no
    `ArtifactDir`) and write under `ValidateArtifactDirectory(req.artifactDir,
    gitState.Root, home)` (`:249`); the Run row is created before the Run
    writes `runs/<run-id>`.
- The maintainer's dry run, read-only, with the installed `roundfix 0.97.0`
  on 2026-10-20 at 13:01 BRT (16:01 UTC), quoted whole except the
  trailing space the command prints after `Artifact Root:`:

  ```text
  GC sanitation dry-run
    Journal Retention: 336h
    Cutoff: 2026-10-06T16:01:59.862282Z
    Artifact Root:
      Classification: unsafe
      Evidence: preserved because recorded Artifact Root "" is not a clean absolute path
    Directories reclaimable: 0
    Artifact bytes reclaimable: 0
    No directories removed; rerun with --apply to mutate.
  ```

- Probe on a temporary Roundfix Home, seeded through a `go test -overlay`
  file that wrote nothing to the repository: one repository, three terminal
  Runs completed 30 days earlier (an Implement Run recording `""`, a `fetch`
  Run recording `relative-artifacts`, a `resolve` Run recording the
  repository's default root as an absolute path), each with bytes under the
  root it wrote to, and one orphaned `runs/run_orphan_without_metadata`
  directory. `roundfix gc sanitize` at `7a678ddb` and at `1e81b9d7` printed
  the same report: root `""` `unsafe`; the default root `orphaned` with two
  `would remove` lines (the `resolve` Run and the orphan) and the Implement
  Run's directory `preserved` as "recorded under Artifact Root """; root
  `relative-artifacts` `unsafe`; 2 directories, 46 bytes. The storage report
  at `1e81b9d7` listed the default root `present` with 2 Runs and
  `<repository>/relative-artifacts` `present`; at `7a678ddb` it listed `.`
  `missing`. The reader left `roundfix.db-wal` and `roundfix.db-shm` beside
  the database and changed no artifact file.
- The resolvers. Run Retention's `runRetentionArtifact`
  (`internal/cli/run_retention.go:186`) resolves only a blank root, inline:
  the key's default through `ResolveArtifactDirectory("", candidate.Repository,
  home)`, the checkout's default through
  `DefaultArtifactDirectoryForPath(candidate.GitRoot, home)`, the checkout's
  when the key's `runs/<run-id>` is absent (`os.Lstat`). A relative recorded
  root is used as is, so its existence is checked against the process's
  working directory. It then requires `classifyGCSanitationRoot` to say
  `orphaned`, resolves `runs/` with `filepath.EvalSymlinks` inside Roundfix
  Home, and requires a physical directory.
- The storage report at `1e81b9d7` (Spec `reports-that-count-what-they-say`,
  Task 02). This plan reuses these symbols:
  - `storageArtifactRoot(run storageRunRecord, homeDir string) (string,
    error)` in `internal/store/journal.go`. It rejects an unsafe Run ID with
    `measure Run artifacts: unsafe Run ID "<id>"`. It resolves a non-blank
    root with `config.ResolveArtifactDirectory(recorded, gitRoot, home)` and
    keeps the cleaned recorded value on error. It resolves a blank root to
    the key's default, or to the checkout's default when the key's
    `runs/<run-id>` is absent and the checkout's is present (`os.Stat`
    through `storagePathExists`), or when the key's default cannot be
    resolved, and to `.` when neither resolves;
  - `storageRunRecord` with `id`, `repository` (the Git root),
    `repositoryRoot` (`COALESCE(NULLIF(repository_root, ''), git_root)`),
    `state` and `artifactRoot`, read by `storageRunRecords`;
  - `(*Store).StorageReport(ctx, homeDir)` and
    `storageFilesystemGroups(runs, homeDir)`, called by
    `measureStorageReport` in `internal/cli/gc.go`;
  - the tests `TestStorageReportMeasuresTheArtifactRootARunWroteTo`
    (`internal/store/journal_test.go`) and
    `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`
    (`internal/cli/gc_test.go`), and the resolution sentences in
    `docs/user-guide/commands/storage-report.md` and the Roundfix `storage`
    reference.
- Other readers of a recorded root, swept with `grep -rn
  'ArtifactDir\b' internal --include='*.go'`: `internal/runcause/report.go`
  (`diagnosticTail`) reads it literally and skips a blank one, the open
  Backlog Entry this planning change writes; `internal/cli/attach.go` passes a
  review Run's recorded root, always absolute; `doctor_storage.go`,
  `reconcile.go`, `gc_item_binary.go` and `deliver_workflow.go` resolve the
  current checkout's configured root, not a recorded one.

## System Architecture

- `internal/store`: `journal.go` exports the resolver and its two errors,
  `StorageReport` calls it, and `DiscoverArtifactRoots` takes Roundfix Home
  and groups by its result. `ArtifactRootRun` gains the recorded value.
- `internal/cli`: `run_retention.go` calls the resolver; `gc.go` passes
  Roundfix Home to discovery, adds the unresolved proof in front of the
  existing ones, rechecks absence before removal and prints the new lines.

## Implementation Design

### Interfaces

```go
// internal/store
var ErrUnsafeRunID = errors.New("unsafe Run ID")
var ErrArtifactRootUnresolved = errors.New("Artifact Root cannot be resolved")

type ArtifactRootRun struct {
	ID, GitRoot, Repository, State string
	ArtifactDir string // the recorded value, trimmed
	CompletedAt *time.Time
}
type ArtifactRoot struct {
	Path       string // the resolved root, or the recorded value when Unresolved is set
	Runs       []ArtifactRootRun
	Unresolved string // why a Run of this group has no resolved root; "" when none
}

func ResolveRunArtifactRoot(run ArtifactRootRun, homeDir string) (string, error)
func DiscoverArtifactRoots(ctx context.Context, runStore *Store, homeDir string) ([]ArtifactRoot, error)
```

### Data Models

`ArtifactRootRun` gains `ArtifactDir`, and `ArtifactRoot` gains `Unresolved`.
The storage report's `storageRunRecord` and `storageArtifactRoot` give way to
the exported resolver. No file, table, record, flag or configuration key
changes shape.

### API Contracts

1. API Contract: the resolver, `ResolveRunArtifactRoot`, in
   `internal/store/journal.go`.
   - An unsafe Run ID (empty, `.`, `..` or not its own base name) returns an
     error wrapping `ErrUnsafeRunID` before any path is built.
   - A non-blank recorded root resolves with
     `config.ResolveArtifactDirectory(recorded, run.GitRoot, home)`: absolute
     stays (cleaned), `~/` expands, relative joins the Git root. A failure
     returns an error wrapping `ErrArtifactRootUnresolved`.
   - A blank recorded root resolves to the key's default,
     `config.ResolveArtifactDirectory("", run.Repository, home)`. The
     checkout's default, `config.DefaultArtifactDirectoryForPath(run.GitRoot,
     home)`, is used instead when it differs and its `runs/<run-id>` is
     present while the key's is absent, or when the key's default cannot be
     resolved. When neither resolves, the error wraps
     `ErrArtifactRootUnresolved`.
   - Presence is `os.Lstat` finding any entry, a symbolic link included, so a
     link at the key's `runs/<run-id>` keeps the key's root. Another `lstat`
     error returns `inspect Run artifact directory <path>: <error>`.
   - Hazard: data comes from the Run Database row and the recorded
     repositories' Git metadata; the result goes to the caller only. The
     resolver never calls `ValidateArtifactDirectory` and creates, writes or
     removes nothing.
2. API Contract: the storage report keeps its output. `StorageReport` calls
   the resolver for every Run. On `ErrArtifactRootUnresolved` it measures the
   value it measured before, the cleaned recorded root or `.` for a blank one.
   On `ErrUnsafeRunID` it fails with `measure Run artifacts: unsafe Run ID
   "<id>"`. Sections, columns, groups, totals and reconciliation are
   unchanged. The only declared difference is the probe: when the key's
   `runs/<run-id>` is a symbolic link whose target is missing, the key's root
   is now kept instead of the checkout's.
3. API Contract: Run Retention reads the shared root.
   `runRetentionArtifact` calls the resolver with the candidate's Run ID, Git
   root, repository root and recorded root. `ErrUnsafeRunID` and
   `ErrArtifactRootUnresolved` make the Run kept as `artifact root unproven`
   without any further filesystem read; another error fails the sweep as
   today. Every later step is unchanged: the Run path, the existence check,
   the `orphaned` proof, the `filepath.EvalSymlinks` containment in Roundfix
   Home and the physical-directory check. Declared changes: a relative or
   `~/` recorded root is now found under the Run's Git root or the home
   directory instead of the working directory, and a Run whose root cannot be
   resolved is kept instead of failing the sweep.
4. API Contract: discovery groups by resolved root.
   `DiscoverArtifactRoots(ctx, runStore, homeDir)` reads the same columns plus
   nothing new and calls the resolver per Run. Runs that resolve to the same
   root form one group. A Run whose error wraps `ErrUnsafeRunID` joins the
   group named by its trimmed recorded value without setting `Unresolved`, so
   the existing proofs report it with today's evidence. A Run with any other
   resolver error, `ErrArtifactRootUnresolved` or an `lstat` failure, joins
   the group named by its trimmed recorded value and sets that group's
   `Unresolved` to `Run "<id>" recorded as "<value>": <error>` (the first such
   Run in Run ID order); discovery itself fails only on a Run Database read
   error, as today. Groups are ordered by `Path`, Runs by ID.
5. API Contract: classification. `classifyGCSanitationRoot` first returns
   `unsafe` with the evidence `preserved because the Artifact Root of <Unresolved>
   cannot be resolved` when `Unresolved` is set. Otherwise every existing proof
   runs unchanged, in its order, with its evidence text, on the resolved
   `Path`, including the default-root proof against each Run of the group.
   `allRunIDs` maps each Run ID to its group's `Path`. Only `runs/<run-id>`
   entries are ever candidates.
6. API Contract: absence is proven again. With `--apply`, `runGCSanitation`
   opens one reader through the GC dependencies' reader opener for the removal
   phase. Immediately before removing a candidate proven only by absence (`Run
   <id> is absent from durable Run metadata`), it reads that Run ID with
   `Store.Run`. When the Run exists, the directory stays, it is reported
   `preserved: <path> (Run <id> appeared in durable Run metadata before
   removal)` and it is not counted. A read error fails the command before that
   removal with the existing `gc sanitize failed` diagnostic. Candidates of
   recorded terminal Runs are removed as today.
7. API Contract: the report names what it resolved.
   - After `  Artifact Root: <path>`, one line `    Recorded as: "<value>" by
     <n> Run(s)` for each distinct recorded value of the group's Runs that
     differs from `<path>`, sorted by value.
   - After the candidate and preserved lines, `    Directories reclaimable
     here: <n>` and `    Artifact bytes reclaimable here: <bytes>` in a dry
     run, or `    Directories removed here: <n>` and `    Artifact bytes
     reclaimed here: <bytes>` with `--apply`, for every root.
   - The totals, header, cutoff and closing lines keep their text, and the
     per-root labels differ from the total labels.
8. API Contract: everything else stays. Command names, flags, exit codes, the
   six classifications and their evidence texts, idempotence, the Journal
   Retention cutoff, `gc compact`, the per-repository `roundfix gc` sweeps,
   `storage report` output and every existing test name keep their behavior.

### Surface Transcripts

1. Surface Transcript: the dry run reaches the default root of a Run that
   recorded none. In a temporary Git repository with one commit, with `HOME`
   set to a temporary directory outside it and no User or Project Config, the
   Run Database holds one Implement Run for that repository with no Artifact
   Directory, terminal `clean`, completed 400 hours before now, and that Run's
   `runs/<run-id>/agent/batch-001.log` holds bytes under the repository's
   default root.

   ```transcript
   $ roundfix gc sanitize
   stdout:
   GC sanitation dry-run
     Journal Retention: 336h
     Cutoff: <cutoff>
     Artifact Root: <home>/.roundfix/artifacts/<repository id>
       Recorded as: "" by 1 Run(s)
       Classification: orphaned
       Evidence: no Active Run records this Artifact Root; terminal Runs recorded: 1; default derived from recorded checkout "<repository>" (Run "<run id>")
       would remove: <home>/.roundfix/artifacts/<repository id>/runs/<run id> (terminal Run <run id> completed before cutoff)
       Directories reclaimable here: 1
       Artifact bytes reclaimable here: <bytes>
     Directories reclaimable: 1
     Artifact bytes reclaimable: <bytes>
     No directories removed; rerun with --apply to mutate.
   stderr:
   exit: 0
   ```

### Invariants

```text
1. One function resolves a Run's Artifact Root; the storage report, GC sanitation and Run Retention call it and carry no resolution of their own.
2. The resolver creates, writes and removes nothing.
3. Sanitation classifies the resolved root with every existing proof, in order; a group with an unresolved Run is unsafe.
4. Sanitation removes only runs/<run-id> directories of retention-eligible terminal Runs recorded under the root and of Run IDs absent from the Run Database at the moment of removal.
5. The storage report prints the same Artifact Roots lines as before this Spec.
6. A second sanitation apply removes nothing.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 3, Core Feature 1, Success Metric 3, Invariants 1-2, 5 → `ResolveRunArtifactRoot`, `StorageReport`, `runRetentionArtifact` (API Contracts 1-3)
- Goal 1, Core Feature 2, Success Metric 1, Invariants 3-4, 6 → `DiscoverArtifactRoots`, `classifyGCSanitationRoot`, `runGCSanitation` (API Contracts 4-6)
- Goal 2, Success Metric 2, Invariants 3-4 → the existing proofs on the resolved root and the recheck (API Contracts 5-6)
- Core Feature 3, Success Metric 1 → `printGCSanitationReport` (API Contract 7), Surface Transcript 1
- Core Feature 4 → the `gc` guide, the Roundfix `storage` reference, **Artifact Root** in `CONTEXT.md`
- Success Metric 4 → the PRD's Unreachable Acceptance
- Goal 2, Core Features 1-2, API Contract 8 → the existing sanitation, retention and storage report suites and the QA scope audit

## Absence Evidence

1. `ResolveRunArtifactRoot`, `ErrUnsafeRunID` and `ErrArtifactRootUnresolved` — command: `grep -rn 'func Resolve.*ArtifactRoot\|ErrUnsafeRunID\|ErrArtifactRootUnresolved' internal`; result: no match at `7a678ddb` or `1e81b9d7`; the only resolvers are the unexported `storageArtifactRoot` and the inline blank-root branch of `runRetentionArtifact`.
2. Discovery by resolved root and the `Unresolved` group — command: `grep -n 'ResolveArtifactDirectory\|homeDir' internal/store/journal.go`; result: at `7a678ddb` no match; at `1e81b9d7` matches only inside the storage report; `DiscoverArtifactRoots` takes no Roundfix Home.
3. The absence recheck before removal — command: `grep -n 'openStoreReader\|\.Run(ctx' internal/cli/gc.go`; result: one `openStoreReader` call, the discovery reader in `runGCSanitation`, and no Run lookup; removal follows the scan with no read.
4. The per-root report lines — command: `grep -rn 'Recorded as\|reclaimable here\|reclaimed here' internal/cli docs .agents/skills`; result: no match at `7a678ddb`.

## ADR Alternatives

- ADR-0301: the design is the decision; none of its rejected alternatives (recording the root on new Runs, a third copy, a default without its proof) is taken.
- ADR-0255: no rejected alternative is taken; Run Retention keeps its kept reasons and proof, and only finds a relative root under the Git root.
- ADR-0296: no rejected alternative is taken; the Storage Snapshot and its tolerance are untouched, and the report's output stays the same.
- ADR-0033: no rejected alternative is taken; sanitation keeps Journal Retention's cutoff and never deletes a `runs` row.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering `storage` reference changes with the `gc` surface.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the maintainer's dry run, the existing suites and published collectors.
- ADR-0240: no rejected alternative is taken; the maintainer's real dry run and apply are one Unreachable Acceptance entry.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0171: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0172: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0184: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0266: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0298: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0299: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

The Run Database under Roundfix Home, read only, plus one Run lookup per
absence-proven candidate with `--apply`; Git metadata of each recorded
repository, read through `config.RepositoryRoot`; directories under each
resolved root, read, and removed only under `runs/` with `--apply`. No test,
Verification command or QA row reads or writes the real `~/.roundfix`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
and uses temporary homes and repositories; the CLI tests run the command in
process through `runCLIContext`, as the existing sanitation tests do.

- `internal/store/journal_test.go`:
  - new `TestResolveRunArtifactRootIsTheRootTheRunWroteTo`, with subtests:
    blank root (the key's default); blank root of a removed linked worktree
    (the main checkout's default); relative root (under the Git root); `~/`
    root; absolute root; blank root whose directory exists only under the
    checkout's default (that root); a symbolic link at the key's
    `runs/<run-id>` (the key's root); an unsafe Run ID (`ErrUnsafeRunID`);
    blank Git and repository roots (`ErrArtifactRootUnresolved`); and the
    home and repository trees unchanged after every call;
  - new `TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo`: a blank
    and an absolute record of the same default root form one group; two
    repositories' blank records form two; a relative record resolves under
    its Git root; a Run with blank Git and repository roots forms a group with
    `Unresolved` set;
  - declared change: `TestDiscoverArtifactRootsReturnsEveryRecordedRepository`
    passes the test's home and keeps its assertions.
- `internal/cli/gc_run_retention_test.go`, new test
  `TestRunRetentionReadsTheRootTheSharedResolverReturns`, through `gc`: a
  terminal Run past Run Retention whose relative root resolves outside
  Roundfix Home with a directory there is kept as `artifact root unproven`;
  one whose resolved relative root holds no directory is removed; a blank-root
  Run whose directory sits under the key's default is removed with it.
- `internal/cli/gc_sanitize_resolved_root_test.go`, new:
  - `TestGCSanitizeDryRunReachesTheDefaultRootTranscript`, asserting Surface
    Transcript 1 whole with its controlled values substituted;
  - `TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof`, with Implement
    Runs recording `""`: an eligible terminal Run and an orphaned directory are
    removed by `--apply`; an Active Run, a terminal Run inside Journal
    Retention, `reviews/` and `delivery/` files, a relative root resolved
    outside Roundfix Home, a symbolically linked default root of another
    repository, and a default root that another repository recorded as its
    own absolute root (`overridden`) all stay; a second `--apply` removes
    nothing;
  - `TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval`, which
    replaces the GC dependencies' reader opener, the seam the GC tests already
    replace, so that the opener's second call first gives the orphan's Run ID a
    row: the directory stays and is reported `preserved`;
  - `TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved`: a Run
    whose Git and repository roots were blanked in the Run Database makes its
    group `unsafe` with the new evidence, and nothing is removed.
- Unchanged and run as consumers: the sanitation tests in
  `internal/cli/gc_test.go` and `internal/cli/gc_sanitize_pre_key_root_test.go`,
  the Run Retention tests in `internal/cli/gc_run_retention_test.go`, and the
  storage report tests named under Current behavior.

## Build Order

1. One resolver: `ResolveRunArtifactRoot` and its errors, the storage report
   on it, Run Retention on it, their tests, and the `gc` guide's Run Retention
   sentence with the Behavior Surface Record (no dependency inside this Spec;
   the storage report code of Spec `reports-that-count-what-they-say` must be
   on the default branch, `_prd.md` → Prerequisites).
2. Sanitation over resolved roots: `DiscoverArtifactRoots(ctx, store, home)`,
   the unresolved proof, the recheck, the report lines, their tests, the `gc`
   guide's Sanitation section and the Behavior Surface Record (depends on: 1,
   for the resolver and the shared guide and record).
3. The Roundfix `storage` reference, the skill version and **Artifact Root**
   in `CONTEXT.md` (depends on: 2, because the reference describes both
   slices).
4. QA gate (depends on: 1, 2, 3).

## Risks & Considerations

- The first `--apply` after the release can remove directories that every
  earlier version kept. The dry run lists each one with its proof, and the
  maintainer applies only a clean dry run (Unreachable Acceptance).
- A Run that recorded `""` but wrote under a root configured at the time has
  no directory under the default root, so it contributes no candidate there;
  its real root stays unreached, as today.
- A Run that starts during the scan writes its row before its directory, so
  the recheck finds it. A Run Retention removal racing with sanitation removes
  the directory first, and a removal of an absent path succeeds.
- The reader opened by `gc sanitize` leaves `-wal` and `-shm` files beside
  the database, as it does today; tests compare artifact trees, not the
  database's sidecars.

## Glossary

- adds: **Artifact Root**

## Decisions

- Export the storage report's resolver rather than adding a sanitation copy;
  its fallbacks become typed errors that each caller maps, so the storage
  report keeps its output and the deleting callers fail closed. See ADR-0301.
- Probe with `os.Lstat`, the way Run Retention probed, so a link is never
  stepped over on the deleting paths.
- Group by resolved root and let an unresolved Run make its whole group
  `unsafe`, which is what an unsafe Run ID already did.
- Recheck only directories proven by absence: a recorded terminal Run cannot
  become live again.
- Name each resolved root and its reclaim in the dry run with labels distinct
  from the totals, so existing readers of the totals keep working.
