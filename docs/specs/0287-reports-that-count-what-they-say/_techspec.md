---
spec: 0287-reports-that-count-what-they-say
prd: _prd.md
created: 2026-10-20
---

# Reports that count what they say — Technical Spec

## Executive Summary

Three independent slices, each inside the function that already builds its
report. No package, command, flag, configuration key or record field is new.

1. The Linux Lane script prints `roundfix-linux-lane: packages: <n>`, the
   number of packages `go list` selected, before it runs `go test`; the host
   reads the first such line and prints it on the pass line
   (`internal/cli/deliver_linux_lane.go`).
2. `Store.StorageReport` resolves each Run's recorded Artifact Root the way
   the Run resolved it when it wrote its artifacts, and measures that root
   (`internal/store/journal.go`).
3. `archiveParagraph`, which reads every Archive Record summary, carries the
   list or paragraph that follows a paragraph ending with a colon
   (`internal/spec/archive_record.go`).

The trade-offs accepted:

- The lane's count comes from the container's output, so a lane whose script
  printed no count passes with a line that says the count is unknown, rather
  than with a number the host would have to guess.
- The report resolves the empty value at read time instead of recording the
  resolved directory on new Runs, so it also covers the 508 Runs already
  recorded, at the cost of reading each Run's repository metadata.
- A carried list is flattened into the summary's single line, because every
  record reader takes the summary as one paragraph.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the count line and the no-count wording are fixed strings.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the lane keeps `--network none` and
  `--pull never`, and every test uses temporary repositories, homes and
  compiled stand-ins. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Build Order follows recorded
  decisions. Build Order 1 follows ADR-0294 ("The Linux Lane runs like CI"),
  ADR-0297 ("The Linux Lane reads the repository history read-only") and
  ADR-0288 ("The Delivery Queue gate runs the changed Go packages in a Linux
  Lane"). Build Order 2 follows ADR-0296 ("The storage report reconciles
  every page of one Storage Snapshot") and ADR-0255 ("Run Retention removes
  terminal Runs whole and compacts incrementally"). Build Order 3 follows
  ADR-0247 ("An archive leaves an Archive Record, and the Spec folder stays
  in Git"), ADR-0277 ("A Legacy Archive Folder conversion keeps the summary,
  links and dates its folder held") and ADR-0262 ("A cancelled Spec leaves an
  Archive Record and returns its sources"). Every step follows ADR-0125,
  ADR-0259, ADR-0271, ADR-0285, ADR-0189, ADR-0256, ADR-0280, ADR-0104,
  ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named with its title in
  `_prd.md` → Project Constraints. Set no obligation here, because this Spec
  changes none of their rules: ADR-0092, ADR-0120, ADR-0165, ADR-0167,
  ADR-0169, ADR-0184, ADR-0196, ADR-0197, ADR-0229, ADR-0233, ADR-0237,
  ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0258, ADR-0260, ADR-0299, ADR-0261,
  ADR-0269, ADR-0272, ADR-0278, ADR-0279, ADR-0281, ADR-0282, ADR-0286,
  ADR-0290 and ADR-0295, each named in `_prd.md` → Project Constraints.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 ("Sim, mesmas regras do AJ (Recomendado)"),
  carrying the cycle AJ grant "Mexer em Governed Paths declarados", and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0287-reports-that-count-what-they-say/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `9448aa64`.

- The lane (`internal/cli/deliver_linux_lane.go`).
  - `runLinuxLane` sets `result.Packages` to `./...` with `--all` (`:92-93`)
    or to `linuxLanePackages` (`:115`), which returns `./...` for a `go.mod`
    or `go.sum` change and the nearest package directories otherwise.
  - `linuxLaneScript` (`:26`) expands those patterns with `go list -e -f
    '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}'`,
    prints `no package builds on linux` and exits `0` when nothing is
    selected, and otherwise runs `exec go test -count=1 $pkgs`.
  - On exit `0` the lane log gets `Linux Lane passed %d package(s) on %s`
    (`:239`) and the console `roundfix: %sLinux Lane passed %d package(s) on
    %s` (`:243`), both with `len(result.Packages)`. The gate (`Gate`,
    `internal/cli/deliver_workflow.go`) and `runDeliverLane`
    (`internal/cli/deliver_lane.go`) both reach this one builder.
  - `docs/user-guide/commands/deliver.md` and
    `.agents/skills/roundfix/references/deliver.md` state "`<n>` counts
    selected package arguments".
  - The tests' Docker stand-in (`linuxLaneToolSource`,
    `internal/cli/deliver_linux_lane_test.go:28`) writes `\n fixture
    output\nsecond line\n` to standard output and `LANE_OUTPUT` to standard
    error on `run`; the host gives Docker one writer for both streams, so
    their order is kept.
  - The operator's `roundfix deliver lane --all` on `main` printed `passed 1
    package(s)` while 43 packages ran (Source 1).
- The storage report (`internal/store/journal.go`).
  - `storageRunRecords` (`:481`) reads `id, git_root, state, artifact_dir`
    and cleans `artifact_dir` with `filepath.Clean`, so `""` becomes `.`.
  - `storageFilesystemGroups` (`:501`) groups Artifact Roots by that value;
    `storageRunArtifactBytes` and `storageDirectoryBytes` treat `""` and `.`
    as absent, and a relative value is measured against the process's
    working directory.
  - Implement Runs record no `artifact_dir` (`internal/cli/implement.go:413`
    passes none; `validateCreateRunRequest` requires it only for `fetch`,
    `resolve` and `watch`), while `runImplement` writes under
    `ValidateArtifactDirectory(req.artifactDir, gitState.Root, home)`
    (`:249`), which for an empty value is
    `<home>/.roundfix/artifacts/<repoID(RepositoryRoot(gitRoot))>`.
  - Run Retention resolves an empty recorded root
    (`internal/cli/run_retention.go:187`, `runRetentionArtifact`): the
    repository key's default through `ResolveArtifactDirectory("",
    repository, home)` and the recorded checkout's default through
    `DefaultArtifactDirectoryForPath(gitRoot, home)`.
  - `internal/cli/gc.go` calls `StorageReport` from `measureStorageReport`
    (`:282`), used by `storage report` (`:273`) and by `gc compact`'s
    live-writer fallback (`:439`).
  - On the maintainer's machine (2026-10-20) the report listed `.` `missing`
    with 508 Run rows and 0 bytes and every state group at 0 bytes, while
    `~/.roundfix/artifacts` held 37 repository directories, about 605 MiB.
    `roundfix gc sanitize` showed the recorded value is `""`.
  - Probe: a temporary store whose Implement Run recorded no root, with 10
    bytes under the Run's default root, reported `root="." missing=true
    rows=1 bytes=0`; a `fetch` Run recording `relative-artifacts` reported
    that root `missing` with 0 bytes.
- The summary (`internal/spec/archive_record.go`).
  - `archiveParagraph` (`:489`) collects one paragraph's trimmed lines and
    stops at the first blank line after text or at a line starting with `#`.
  - `BuildArchiveRecord` reads the PRD's title paragraph with it (`:186`),
    the first section's paragraph for a Legacy Archive Folder (`:188`) and
    the newest QA Report's `## Outcome` paragraph (`:263`).
    `legacyArchiveParagraph` (`:516`) strips block quote markers first. The
    Cancel Command (`internal/spec/cancel.go:207`) and `ParseArchiveRecord`
    (`:676`) use it too.
  - `RenderArchiveRecord` (`:608`) writes the summary as one line after the
    title, cut with `…` past `ArchiveRecordTargetBytes`.
    `history_sanitize.go` (`:336-349`) and the pull request record update
    (`internal/cli/deliver_workflow.go:1412-1425`) parse a rendered record
    and render it again, expecting the same bytes.
  - Probe: a PRD whose title paragraph was "Two gaps failed it:" followed by
    a two-item list gave the summary "Two gaps failed it:"; a QA Outcome "The
    gate held three facts:" followed by an ordered list gave "The gate held
    three facts:".
  - On 2026-10-20, 38 of 256 records under `docs/history/specs/` ended with a
    colon.

## System Architecture

Three packages, one slice each, no new seam.

- `internal/cli`: `deliver_linux_lane.go` gains the script's count line and
  the host's reading of it; `deliver_workflow.go` and `deliver_lane.go` keep
  calling `runLinuxLane` unchanged.
- `internal/store`: `journal.go` resolves each Run's root inside
  `StorageReport`, which gains the Roundfix Home it resolves against;
  `internal/cli/gc.go` passes it.
- `internal/spec`: `archive_record.go`'s `archiveParagraph` carries the
  following block; every summary reader inherits it.

## Implementation Design

### Interfaces

```go
// internal/cli
// linuxLanePackageCount reads the first line that starts with
// "roundfix-linux-lane: packages: " and holds only decimal digits after it.
func linuxLanePackageCount(output []byte) (count int, reported bool)

// internal/store
func (store *Store) StorageReport(ctx context.Context, homeDir string) (StorageReport, error)

// storageArtifactRoot returns the root a Run wrote its artifacts to.
func storageArtifactRoot(run storageRunRecord, homeDir string) (string, error)

// internal/spec — unchanged signature, extended behavior.
func archiveParagraph(body, after string) string
```

### Data Models

`storageRunRecord` gains the Run's Git root and repository root
(`COALESCE(NULLIF(repository_root, ''), git_root)`); its `artifactRoot`
becomes the resolved root. No file, table, record or configuration key
changes shape.

### API Contracts

1. API Contract: the lane script becomes exactly:
   `fail() { printf 'roundfix-linux-lane: container: %s\n' "$1" >&2; exit 3; }; for tool in tar git go; do command -v "$tool" >/dev/null 2>&1 || fail "missing tool $tool"; done; [ -d /repo-objects ] && [ -r /repo-objects ] && [ -x /repo-objects ] || fail "cannot read the repository history at /repo-objects"; mkdir -p /tmp/work && cd /tmp/work || fail "cannot prepare /tmp/work"; tar -xf - || fail "cannot unpack the candidate"; { git init -q && echo /repo-objects > .git/objects/info/alternates && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane"; } || fail "cannot record the candidate"; pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; echo "roundfix-linux-lane: packages: 0"; exit 0; fi && set -- $pkgs && echo "roundfix-linux-lane: packages: $#" && exec go test -count=1 "$@"`
   - The count line goes to standard output before `go test` starts, so the
     first line of the combined output that starts with `roundfix-linux-lane:
     packages: ` is the script's; a test that prints one later is ignored.
2. API Contract: the pass lines.
   - On exit `0`, `linuxLanePackageCount` reads the combined output. With a
     count `<n>`, the lane log ends with `Linux Lane passed <n> package(s) on
     <image>` and the console line is `roundfix: <label>Linux Lane passed <n>
     package(s) on <image>`.
   - Without one, the lane log ends with `Linux Lane passed on <image>; the
     lane script reported no package count` and the console line is
     `roundfix: <label>Linux Lane passed on <image>; the lane script reported
     no package count`. The outcome stays a pass.
   - `<label>` is `Delivery Queue item <slug>: ` in the gate and empty in
     `roundfix deliver lane`. `result.Packages` keeps the patterns; nothing
     else in the lane changes: failure lines, the classification, blockers,
     exit codes, mounts, image and package selection.
   - Data source: the container's combined output, which the candidate's
     tests also write to. Destination: the console and the lane log only.
3. API Contract: the resolved Artifact Root, in `StorageReport`.
   - A recorded root that is not blank resolves with
     `config.ResolveArtifactDirectory(recorded, gitRoot, homeDir)`: an
     absolute path stays as recorded (cleaned), `~` expands, and a relative
     path joins the Run's Git root. A resolution error keeps the cleaned
     recorded value.
   - A blank recorded root resolves to the repository key's default,
     `config.ResolveArtifactDirectory("", repositoryRoot, homeDir)`. When the
     recorded checkout's default, `config.DefaultArtifactDirectoryForPath(gitRoot,
     homeDir)`, differs and holds `runs/<run-id>` while the key's default does
     not, it is used instead. When the key's default cannot be resolved, the
     checkout's default is used; when neither can, the value stays `.` and is
     reported `missing`, as today.
   - The Run ID is checked as `storageRunArtifactBytes` checks it before any
     path is built; an unsafe ID fails the report as today. A stat error
     other than not-exist fails it with `inspect Run artifact directory <path>:
     <error>`.
   - The Artifact Root groups, the repository groups and the state groups
     all measure `<resolved root>/runs/<run-id>` and the resolved root; the
     sections, columns, totals and reconciliation keep their shape, and the
     `ARTIFACT_ROOT` column prints the resolved path.
   - Hazard: data comes from the Run Database and the repositories' Git
     metadata; the destination is standard output only. The report never
     calls `ValidateArtifactDirectory`, so it creates no directory and writes
     no byte, in the database or under any root.
   - `measureStorageReport` passes its `homeDir`; `storage report` and `gc
     compact`'s fallback keep their behavior otherwise.
4. API Contract: the summary that carries its list, in `archiveParagraph`.
   - When the collected paragraph, joined, ends with `:`, reading continues
     past blank lines to the next block, and exactly one block is carried:
     - a list: its first line starts, after up to three spaces, with `-`,
       `*` or `+` followed by a space, or with digits followed by `.` or `)`
       and a space. Every line that starts that way, at any indentation,
       starts an item; the other non-blank lines continue the item. A blank
       line followed by an item line continues the list; a blank line
       followed by anything else, or a line starting with `#`, ends it. Each
       item loses its marker and is whitespace-normalized; an item that does
       not end with `.`, `!`, `?`, `:` or `;` gains `;`, or `.` when it is the
       last item. The items follow the paragraph, separated by single spaces;
     - a paragraph: its lines up to the next blank line or `#` line,
       whitespace-normalized, follow after one space;
     - a line starting with `#`, a fence (`` ``` `` or `~~~`), a table row
       (`|`) or the end of the text: nothing is carried, and the paragraph
       stays as today.
   - Every caller inherits it: the PRD title paragraph, a Legacy Archive
     Folder's section fallback after its block quote markers are removed,
     the QA Report's Outcome, the Cancel Command and `ParseArchiveRecord`.
   - A written record holds one summary line, so reading it back carries
     nothing: every existing record, including those ending with a colon,
     parses and renders to the same bytes, and no record is rewritten.
5. API Contract: everything else stays. Command names, flags, exit codes,
   the storage report's sections, the Archive Record schema and fields, the
   size cut, the outcome precedence (supersession, QA Outcome, PRD), Legacy
   Link Rebase, and every existing test name keep their behavior.

### Surface Transcripts

Both run in a temporary repository whose `go.mod` declares `go 1.26.1` and
whose `.roundfixrc.yml` sets `delivery.linux_lane: false`, with `HOME` set to a
temporary directory and compiled Docker and Go stand-ins first on `PATH`. The
Docker stand-in answers `info` and `image inspect` with exit `0` and, on `run`,
reads its input and exits `0`.

1. Surface Transcript: the on-demand lane reports the packages the container
   selected. The Docker stand-in writes `roundfix-linux-lane: packages: 43`
   to standard output on `run`.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: Linux Lane passed 43 package(s) on golang:1.26.1
   stderr:
   exit: 0
   ```

2. Surface Transcript: the on-demand lane passes without a package count. The
   Docker stand-in writes no count line.

   ```transcript
   $ roundfix deliver lane --all
   stdout:
   roundfix: Linux Lane passed on golang:1.26.1; the lane script reported no package count
   stderr:
   exit: 0
   ```

### Invariants

```text
1. The Linux Lane pass line prints a number only when the lane script reported it, and that number is the count of packages go list selected in the container.
2. The script's count line precedes go test, and only the first count line in the output is read.
3. The storage report measures each Run's artifacts under the root the Run wrote them to: as recorded when absolute, under the Git root when relative, under the repository default when blank.
4. The storage report creates no directory and writes no byte.
5. A summary paragraph that ends with a colon carries the one list or paragraph that follows it.
6. Every Archive Record already written parses and renders to the same bytes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariants 1-2 → the lane script's count line and `linuxLanePackageCount` (API Contracts 1-2), Surface Transcripts 1-2
- Goal 2, Core Feature 2, Success Metric 2, Invariants 3-4 → `storageArtifactRoot` and `StorageReport` (API Contract 3)
- Goal 3, Core Feature 3, Success Metric 3, Invariant 5 → `archiveParagraph` (API Contract 4)
- Goal 4, Success Metric 3, Invariant 6 → `ParseArchiveRecord` and `RenderArchiveRecord` round trip (API Contract 4), the QA scope audit
- Core Feature 4 → the `deliver`, `storage report`, `archive` and `history` guides, the Roundfix `deliver`, `storage` and `archive` references, **Archive Record** in `CONTEXT.md`
- Success Metric 4 → the QA gate's real lane row when Docker is reachable, or the PRD's Unreachable Acceptance
- Core Features 1-3 and API Contract 5 → the existing lane, storage, archive, cancel and history tests and the QA scope audit

## Absence Evidence

1. The lane's count line and `linuxLanePackageCount` — command: `grep -rn 'packages: ' internal/cli/deliver_linux_lane.go docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md`; result: no match at `9448aa64`; no code or guide reads a package count from the container.
2. `storageArtifactRoot` — command: `grep -n 'ResolveArtifactDirectory\|DefaultArtifactDirectoryForPath\|repository_root' internal/store/journal.go`; result: one match, `DiscoverArtifactRoots`'s query at `:1585` selecting `repository_root`; nothing in the storage report resolves a recorded root.
3. The carried list in `archiveParagraph` — command: `grep -n 'HasSuffix' internal/spec/archive_record.go`; result: no match at `9448aa64`; no summary reader looks at how its paragraph ends.

## ADR Alternatives

- ADR-0294: no rejected alternative is taken; the lane's image, init, classification and skips stay, and only the pass line's number changes.
- ADR-0297: no rejected alternative is taken; the object store mount and its skips are untouched, and the script keeps its history check and alternates line.
- ADR-0288: no rejected alternative is taken; the lane still never pulls, mounts no checkout and stays opt-in, and the gate keeps its single builder.
- ADR-0296: no rejected alternative is taken; the tolerance stays 0 and the Storage Snapshot reader is unchanged; only the artifact groups' paths are resolved.
- ADR-0255: no rejected alternative is taken; the report reuses Run Retention's two default roots for a blank recorded root and removes nothing.
- ADR-0247: no rejected alternative is taken; the record stays small, its summary stays one line, and the Spec folder stays in Git.
- ADR-0277: no rejected alternative is taken; records already written stay byte-identical, and a Legacy Archive Folder's summary keeps its block quote reading before the list is carried.
- ADR-0262: no rejected alternative is taken; a cancelled Spec's record keeps its fields and gains only the carried list in its summary.
- ADR-0125: no rejected alternative is taken; the Docker and Go stand-ins stay compiled binaries.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0189: no rejected alternative is taken; each Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering skill changes with the `deliver`, `storage` and `archive` surfaces.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the repository's records, the operator's runs and published documentation.
- ADR-0240: no rejected alternative is taken; the real lane and the maintainer's storage report are declared under Unreachable Acceptance.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and the record invariance compares with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0184: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0299: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

Docker, only through the existing `docker` executable on the operator's or
Queue Owner's `PATH`; the Run Database and recorded Artifact Roots under
Roundfix Home, read only; Git metadata of each recorded repository, read
through `config.RepositoryRoot`. No test, Verification command or QA row needs
Docker, except the QA row that runs the real lane only when it can, and none
reads the real `~/.roundfix`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories and homes, and ends every child it starts in
`t.Cleanup`.

- `internal/cli/deliver_linux_lane_test.go`:
  - declared change: `linuxLaneToolSource` writes `roundfix-linux-lane:
    packages: <LANE_PACKAGES>` to standard output before its fixture output
    on `run` when `LANE_PACKAGES` is not empty, and `newLinuxLaneFixture`
    sets `LANE_PACKAGES=1`, so the existing pass-line assertions keep their
    text;
  - declared change: `TestLinuxLaneRunsTheChangedPackagesInAContainer`'s
    independent script copy becomes API Contract 1's, keeping its name;
  - new: `TestLinuxLaneReportsThePackagesTheContainerSelected`, through the
    gate with a `go.mod` change (pattern `./...`): `LANE_PACKAGES=43` gives
    the 43-package console and log lines; `LANE_PACKAGES=0` gives `passed 0
    package(s)`; an empty `LANE_PACKAGES` gives the no-count lines; a later
    `roundfix-linux-lane: packages: 7` in `LANE_OUTPUT` leaves 43; and the
    script constant holds `echo "roundfix-linux-lane: packages: $#"` before
    `exec go test`.
- `internal/cli/deliver_lane_test.go`: new
  `TestDeliverLanePackageCountTranscript` and
  `TestDeliverLaneNoPackageCountTranscript`, asserting Surface Transcripts 1
  and 2 whole.
- `internal/store/journal_test.go`:
  - declared change: the five `StorageReport(ctx)` calls pass the test's
    home;
  - new: `TestStorageReportMeasuresTheArtifactRootARunWroteTo`, with
    subtests: an Implement Run with no recorded root whose default root holds
    bytes (present, bytes in the root, repository and state groups); a Run
    recorded from a removed linked worktree whose repository root is the
    main checkout (the main checkout's default); a relative recorded root (under the
    Git root); a Run whose directory exists only under its checkout's default
    (that root); and a blank root whose default does not exist (reported
    `missing`, and the directory still absent after the report).
- `internal/cli/gc_test.go`: new
  `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`, through `storage
  report`: the Artifact Roots line is `<home>/.roundfix/artifacts/<id>` tab
  `present` tab `1` tab the bytes, no `.` root is listed, and the database
  bytes are unchanged.
- `internal/spec/archive_record_test.go`: new
  `TestArchiveSummaryCarriesWhatItsColonIntroduces`, through
  `BuildArchiveRecord`, with subtests: a tight list, a loose list, an ordered
  list, items without closing punctuation, a following paragraph, a fence, a
  heading, end of text, a QA Outcome; and
  `TestAnArchiveRecordEndingWithAColonRoundTrips`, a record whose summary is
  "Two environment gaps failed it:" parsing and rendering to the same bytes.
- `internal/spec/archive_record_cancelled_test.go`: new
  `TestACancelledRecordSummaryCarriesItsList`.
- `internal/spec/legacy_record_summary_test.go`: new
  `TestLegacySummaryCarriesTheListItsColonIntroduces`, a block-quoted
  paragraph and list in a Legacy Archive Folder.

## Build Order

1. The lane's package count: the script, `linuxLanePackageCount`, the pass
   lines, their tests, the `deliver` guide and reference, the skill version
   and the Behavior Surface Record (no dependency).
2. The storage report's resolved roots: `StorageReport(ctx, homeDir)`,
   `storageArtifactRoot`, the caller, their tests, the `storage report` guide
   and the Roundfix `storage` reference, the skill version and the Behavior
   Surface Record (depends on: 1, because both raise the Roundfix skill
   version and re-record the same files).
3. The summary that carries its list: `archiveParagraph`, its tests, the
   `archive` and `history` guides, the Roundfix `archive` reference,
   **Archive Record** in `CONTEXT.md`, the skill version and the Behavior
   Surface Record (depends on: 2, for the same shared files).
4. QA gate (depends on: 1, 2, 3).

No other Spec is a prerequisite. Specs 0285 and 0286 may merge first and move
the Roundfix skill version and the Behavior Surface Record, so this plan is
rebased on the default branch and checked again before its Run.

## Risks & Considerations

- A candidate's test that prints a count line before the script does cannot
  exist: the script prints its line before `go test` starts. A test can only
  print a later one, which is ignored.
- A pass with the no-count line means the container exited `0` without
  running the script's last steps the usual way, for example an image whose
  entrypoint swallowed the command; the line names it instead of hiding it.
- Resolving a blank root reads each Run's repository metadata once per Run;
  508 Runs read at most 508 `.git` entries, well inside the report's cost.
- A summary paragraph that ends with a colon before a fence or a table still
  ends with the colon; lesson 36 of the authoring briefing keeps PRD title
  paragraphs self-contained meanwhile.
- A long carried list is cut with `…` at `ArchiveRecordTargetBytes`, as any
  long summary is today.

## Glossary

- changes: **Archive Record**

## Decisions

- Count the packages `go list` selected, reported by the script, rather than
  calling the number "patterns": the patterns (`./...`) say nothing about
  how much ran, and only the container knows what its `go list` selected.
- Read the first count line only, so a test's output cannot replace it.
- Pass without a number when the count is missing; the outcome
  classification is not this Spec's to change.
- Resolve a blank recorded root at read time with Run Retention's two
  defaults, preferring the repository key's, and use the checkout's only when
  the Run's directory is there and not under the key's.
- Pass Roundfix Home to `StorageReport` explicitly instead of deriving it from
  the database path.
- Carry one following block, flattened into the summary line; never rewrite a
  written record.
- No new ADR: each fix restores what its report already promised, and no
  recorded decision changes.
