---
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
prd: _prd.md
created: 2026-10-20
---

# Runs causes reads the diagnostic where the Run wrote it — Technical Spec

## Executive Summary

One code slice and one documentation slice, with no new package, command,
flag, configuration key, record field or JSON field.

1. `runcause.Request` gains the user's home, which `roundfix runs causes`
   already loads, and `diagnosticTail` (`internal/runcause/report.go`) finds
   the Run's directory under the root `store.ResolveRunArtifactRoot` returns
   instead of under the recorded value taken literally. The containment check
   and the component walk then run unchanged on that directory.
2. The `runs` guide, the Roundfix `runs` reference and the glossary entry for
   **Artifact Root** say where the report reads diagnostics.

The trade-offs accepted:

- The report resolves at read time, so it also reaches every Run already
  recorded with an empty root, at the cost of reading each Run's repository
  metadata and, for a Run whose checkout and repository key differ, one
  `lstat` per failed attempt.
- A blank recorded root resolves to the default root even for a Run that
  wrote under a root configured at the time. Such a Run's diagnostic lies
  outside the resolved root, so the containment check gives it no evidence, as
  today; nothing wrong is read.
- The component walk keeps refusing any symbolic link on the path, the home
  directory's own components included, and keeps its check-then-open window
  (`_prd.md` → Acceptance evidence, Non-Goals).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  field is minted; Run IDs, repository IDs, Artifact Root paths and the
  `roundfix/runs-causes/v1` schema keep their form. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the report reads the Run Database, Git
  metadata and files under Roundfix Home, and every test uses temporary homes
  and repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0301 ("GC
  sanitation classifies the Artifact Root each Run wrote to"), ADR-0255 ("Run
  Retention removes terminal Runs whole and compacts incrementally") and
  ADR-0214 ("Why a Verification failed is classified by deterministic
  signatures"). Every step follows ADR-0259, ADR-0271, ADR-0285, ADR-0189,
  ADR-0256, ADR-0280, ADR-0104, ADR-0240, ADR-0276, ADR-0265 and ADR-0183,
  each named with its title in `_prd.md` → Project Constraints. Set no obligation here, because this Spec changes none of their rules: ADR-0167, ADR-0184, ADR-0229, ADR-0233, ADR-0250, ADR-0281, ADR-0286, ADR-0290, ADR-0295, ADR-0296, ADR-0298, ADR-0197, ADR-0237, ADR-0272, ADR-0282, ADR-0278, ADR-0279, ADR-0299, ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0277, ADR-0165, ADR-0251, ADR-0262, ADR-0269, ADR-0306, ADR-0169, ADR-0196, each named in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 (the same rules as cycle AL, "Mesmas regras do AK
  (Recomendado)"), carrying the cycle AJ grant "Mexer em Governed Paths
  declarados", and the 2026-09-30 skills grant ("considere autorizado a
  ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `70c78d15` (v0.100.0). Line numbers are those
of that commit.

- The read (`internal/runcause/report.go`).
  - `Build` (`:113`) lists the repository's terminal Runs, and for each failed
    Verification attempt calls `diagnosticTail(run, p.DiagnosticPath,
    t.DiagnosticTailBytes)` (`:217`), whose result is the evidence the
    `diagnostic` signatures read.
  - `diagnosticTail` (`:368`) returns `""` when the path is empty, the recorded
    `run.ArtifactDir` is empty, or the Run ID is not its own base name. It
    builds `filepath.Join(run.ArtifactDir, "runs", run.ID)` (`:372`), requires
    an absolute diagnostic path whose `filepath.Rel` to that directory stays
    inside it, walks every component from `/` with `os.Lstat` refusing a
    symbolic link or a missing component, requires a regular file, and reads
    at most the table's tail bytes. An open, stat, seek or read failure other
    than a missing file fails the report.
  - `runcause.Request` (`:30`) carries no home directory, and
    `runRunsCausesCommand` (`internal/cli/runs_causes.go`) fills it from
    `loaded` without one, although `loaded.HomeDir` opens the reader.
- The writer. Implement Runs record no Artifact Directory
  (`internal/cli/implement.go:413`) and write under
  `ValidateArtifactDirectory(req.artifactDir, gitState.Root, home)`
  (`:249`). The Daemon records each failed attempt's absolute
  `diagnostic_path` under `<that root>/runs/<run-id>/verification/`
  (`internal/daemon/daemon.go:196`, `internal/daemon/engine.go:588`).
- The resolver, on the default branch since Spec
  `gc-sanitize-reaches-the-default-artifact-roots`:
  `store.ResolveRunArtifactRoot(run ArtifactRootRun, homeDir string) (string,
  error)` in `internal/store/journal.go` (`:523`), with `ErrUnsafeRunID` and
  `ErrArtifactRootUnresolved`. A blank root resolves to the repository key's
  default or, for an older Run, the checkout's default when only it holds
  `runs/<run-id>`; a relative root resolves under the Run's Git root; an
  absolute root stays. It creates, writes and removes nothing. Its callers are
  the storage report (`:576`), `DiscoverArtifactRoots` (`:1690`) and Run
  Retention (`internal/cli/run_retention.go:187`).
- Probe through the command, on a temporary Roundfix Home, with a `go test
  -overlay` file in `internal/cli` that wrote nothing to the repository: the
  fixture of `TestRunsCausesPrintsAClassifiedWindow` (one terminal Implement
  Run of Spec `0300-example`, three failed attempts whose diagnostics read
  `database is locked`, `--- FAIL: TestCausesFixture` and `no matching
  words`), with the Run's directory and recorded `diagnostic_path` values
  moved to match three recorded roots. `roundfix runs causes --since
  2026-10-01` printed:
  - recorded root absolute (`<home>/artifacts`): task_01 attempt 1
    `environment database-locked`, task_02 attempt 1 `implementation-defect
    go-test-failure`, summary `implementation-defect 1, environment 1,
    unclassified 1; repository knowledge 1 of 3 classified`;
  - recorded root `""`, directory under the repository's default root
    `<home>/.roundfix/artifacts/<repository id>`: all three attempts
    `unclassified -`, summary `implementation-defect 0, environment 0,
    unclassified 3; repository knowledge 1 of 1 classified`;
  - recorded root `rel-artifacts`, directory under `<repository>/rel-artifacts`:
    the same three `unclassified` attempts.
  The same overlay with a prototype of API Contract 1 printed the absolute
  root's output for all three, and the existing `internal/runcause` suite and
  `runs causes` CLI tests passed with it.
- Readers of a recorded Artifact Directory, swept with `grep -rn
  '\.ArtifactDir\b' internal cmd --include='*.go'` and `grep -rn 'artifact_dir'
  internal --include='*.go'` over non-test files:
  1. `internal/runcause/report.go:369,372`, `diagnosticTail`: reads a Run's
     own `runs/<run-id>` directory from the recorded value literally. Fixed by
     this Spec.
  2. `internal/cli/attach.go:276`, `attachIssues`: passes a Run's recorded
     root to `rounds.SelectCompatibleIssues`, which reads the pull request's
     `reviews/pr-<N>/` artifacts, not `runs/<run-id>`. Review Runs record the
     root already resolved to an absolute path
     (`resolveArtifactDirectoryForPreflight`, `internal/cli/cli.go:1819`), and
     an Implement Run carries no pull request number, so the selection refuses
     before reading. Listed, unchanged.
  3. `internal/store/journal.go:487,527,1663,1689`,
     `internal/store/run_retention.go:60,76` with
     `internal/cli/run_retention.go:187`: the storage report, the resolver,
     discovery and Run Retention, which already resolve. Unchanged.
  4. `internal/cli/gc.go:595`: counts each distinct recorded value for the
     `Recorded as` lines of `gc sanitize`, a display of the record by design.
     Unchanged.
  5. `internal/store/store.go:641,2287,2475` and the other `artifact_dir`
     selects in `store.go`: the insert, the creation event payload and the row
     scan. Unchanged.
  Every other `.ArtifactDir` hit reads the current configuration
  (`loaded.Config.Defaults.ArtifactDir` in `gc.go`, `gc_item_binary.go`,
  `doctor_storage.go`, `reconcile.go`, `deliver*.go`, `review.go`,
  `settle.go`, `implement.go`, `cli.go`) or an in-process request already
  resolved by its writer (`internal/daemon`, `internal/agent`,
  `internal/rounds`, `internal/tui`, `detach.go`), not a recorded value.

## System Architecture

- `internal/runcause`: `Request` gains `HomeDir`; `diagnosticTail` takes it
  and calls `store.ResolveRunArtifactRoot`. No new file, type or package
  dependency (`internal/store` is already imported).
- `internal/cli`: `runRunsCausesCommand` sets `req.HomeDir =
  loaded.HomeDir`. Flags, help, exit codes, output and JSON are unchanged.

## Implementation Design

### Interfaces

```go
// internal/runcause
type Request struct {
	Reader                 Reader
	RepositoryRoot         string
	SpecsRoot, ArchiveRoot string
	PushRemote             string
	HomeDir                string // the user's home; resolves a Run's Artifact Root
	Since, Until           time.Time
}

func diagnosticTail(run store.Run, homeDir, path string, limit int64) (string, error)
```

### Data Models

`runcause.Request` gains `HomeDir`. No file, table, record, flag, JSON field
or configuration key changes shape.

### API Contracts

1. API Contract: the resolved read, `diagnosticTail` in
   `internal/runcause/report.go`.
   - An empty diagnostic path gives `""`, as today.
   - The root is `store.ResolveRunArtifactRoot(store.ArtifactRootRun{ID:
     run.ID, GitRoot: run.GitRoot, Repository: <run.RepositoryRoot, or
     run.GitRoot when blank>, ArtifactDir: <run.ArtifactDir trimmed>},
     homeDir)`; the repository fallback mirrors the `COALESCE(NULLIF(
     repository_root, ''), git_root)` the resolver's other callers read.
   - An error wrapping `store.ErrUnsafeRunID` or
     `store.ErrArtifactRootUnresolved` gives `""` without any filesystem
     read, as a blank root and an unsafe Run ID give today. Any other resolver
     error (an `lstat` failure while probing `runs/<run-id>`) fails the report
     with `resolve Artifact Root of Run "<id>": <error>`, as an open failure
     other than a missing file already does.
   - The Run's directory is `<resolved root>/runs/<run-id>`. The absolute-path
     requirement, the `filepath.Rel` containment, the `os.Lstat` walk of every
     component, the regular-file requirement, the tail bound and the error
     texts are unchanged.
   - Hazard: data comes from the Run Database row, the recorded repositories'
     Git metadata and the user's home that `runs causes` already loads; the
     resolved root and the diagnostic text go only to the classifier, and no
     diagnostic text, path or key reaches stdout or JSON, as today. The read
     creates, writes and removes nothing.
2. API Contract: the command passes the home. `runRunsCausesCommand` sets
   `Request.HomeDir` from `loaded.HomeDir`, the same value that opens the Run
   Database reader. A caller that leaves `HomeDir` empty still gets every
   absolute and relative recorded root, and no evidence for a blank one.
3. API Contract: everything else stays. Command name, flags, help, exit codes,
   text and JSON output shape, the signature table and its digest, the
   resolver, the storage report, GC sanitation, Run Retention, `attach`, and
   every existing test name and body keep their behavior. Declared change:
   an attempt of a Run that recorded a blank or relative root, whose
   diagnostic sits under the resolved root and passes the guards, is now
   classified from its diagnostic; the counts in the summary line move with
   it.

### Surface Transcripts

1. Surface Transcript: the report reads the diagnostic of a Run that recorded
   no Artifact Directory. In a temporary Git repository whose
   `docs/specs/0300-example/` holds the `_tasks.md` and `task_04.md` of the
   `TestRunsCausesPrintsAClassifiedWindow` fixture, with `HOME` set to a
   temporary directory outside it whose path has no symbolic link and no User
   or Project Config, the Run Database holds one terminal `clean` Implement
   Run of Spec `0300-example` created 2026-10-01T12:00:00Z with no Artifact
   Directory, the fixture's events, and its three diagnostics under
   `<home>/.roundfix/artifacts/<repository id>/runs/<run id>/verification/`.

   ```transcript
   $ roundfix runs causes --since 2026-10-01
   stdout:
   verification 0300-example/task_01 <run id> attempt-1 environment database-locked: go test ./internal/store
   verification 0300-example/task_02 <run id> attempt-1 implementation-defect go-test-failure: go test ./internal/cli
   verification 0300-example/task_02 <run id> attempt-2 unclassified -: make verify
   corrective 0300-example/task_04 <run id> - repository-convention record-or-golden: pre-pr-review
   Causes: 4 item(s) from 1 terminal Spec Run(s) in window 2026-10-01..*; scope-or-authorization 0, shared-section-contract 0, repository-convention 1, implementation-defect 1, environment 1, unclassified 1; repository knowledge 1 of 3 classified; signatures <digest>
   stderr:
   exit: 0
   ```

### Invariants

```text
1. runs causes finds a Run's directory under the root store.ResolveRunArtifactRoot returns and carries no resolution of its own.
2. A diagnostic outside <resolved root>/runs/<run-id>, behind a symbolic link, not a regular file or missing gives no evidence.
3. A Run whose root cannot be resolved, or whose ID is unsafe, gives no evidence and does not fail the report.
4. The report creates, writes and removes nothing under the home or the repository.
5. No diagnostic text, absolute path or key reaches stdout or JSON.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1-2, Invariant 1 → `diagnosticTail`, `Request.HomeDir`, `runRunsCausesCommand` (API Contracts 1-2), Surface Transcript 1
- Goal 2, Success Metric 4, Invariant 1 → `store.ResolveRunArtifactRoot` as the only resolution (API Contract 1)
- Goal 3, Core Feature 2, Success Metric 3, Invariants 2-5 → the unchanged guards on the resolved directory (API Contracts 1, 3) and the existing `runs causes` suites
- Core Feature 3 → the `runs` guide, the Roundfix `runs` reference, **Artifact Root** in `CONTEXT.md`

## Absence Evidence

1. `Request.HomeDir` and the resolved read in `diagnosticTail` — command: `grep -n 'HomeDir\|ResolveRunArtifactRoot\|ResolveArtifactDirectory\|DefaultArtifactDirectoryForPath' internal/runcause/report.go internal/cli/runs_causes.go`; result: at `70c78d15` the only match is `openRunsListReader(ctx, loaded.HomeDir)` in `internal/cli/runs_causes.go`; `internal/runcause` resolves nothing and joins the recorded value at `report.go:372`.
2. A test of a Run that recorded a blank or relative root — command: `grep -n 'ArtifactDir: ""\|rel-artifacts\|ResolveArtifactDirectory' internal/runcause/report_test.go internal/cli/runs_causes_test.go`; result: no match at `70c78d15`; every `runs causes` fixture records the absolute `<home>/artifacts`.

## ADR Alternatives

- ADR-0301: the design applies its one-resolver decision to a fourth, read-only caller; its rejected alternative, another copy of the resolution, is not taken, and recording the resolved root on new Runs is not taken either.
- ADR-0255: no rejected alternative is taken; Run Retention and its proofs are untouched.
- ADR-0214: no rejected alternative is taken; classification stays by signature from the same three evidence sources, and the report stays read-only.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering `runs` reference changes with the `runs causes` behavior.
- ADR-0280: no rejected alternative is taken; the guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the published measurement, the existing suites and the Go article.
- ADR-0240: no rejected alternative is taken; this Spec declares no unreachable acceptance and its web source is an offline-blocked row.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0184: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0296: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0298: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
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
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0306: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

The Run Database under Roundfix Home, read only through the reader that never
migrates; Git metadata of each recorded repository, read through
`config.RepositoryRoot` inside the resolver; diagnostics under each resolved
root, read only. No test, Verification command or QA row reads or writes the
real `~/.roundfix`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
starts no process in `internal/runcause` (which installs no suite guard), and
uses temporary homes and repositories with symbolic links resolved, because
the component walk refuses a linked temporary directory such as macOS `/var`.

- `internal/runcause/report_test.go`, new
  `TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot`, with subtests,
  each moving the fixture's Run directory and recorded `diagnostic_path`
  values to the root under test and setting `Request.HomeDir` to the fixture
  home:
  - `blank root`: recorded `""`, directory under the repository's default
    root; task_01 attempt 1 is `environment` with signature
    `database-locked`;
  - `relative root`: recorded `rel-artifacts`, directory under the fixture
    repository (with its symbolic links resolved); classified as above;
  - `absolute root`: the fixture as it is; classified as above;
  - `blank root without home`: recorded `""` with `Request.HomeDir` empty; no
    error, task_01 attempt 1 `unclassified` with nil evidence;
  - `leaf symlink`, `parent symlink`, `outside`, `other run`, `directory`,
    `missing`: recorded `""`, each guard of
    `TestBuildIgnoresALogOutsideTheArtifactDirectory` applied under the
    default root; task_01 attempt 1 `unclassified` with nil evidence;
  - `diagnostic under the recorded value`: recorded `""` with the diagnostic
    left under `<home>/artifacts` (a root the Run did not record); no evidence;
  - `inspection error`: the Run's recorded Git root set to a second plain
    directory (its repository root, which the listing filters on, stays), so
    the key's and the checkout's defaults differ, with a regular file at the
    key's `runs`, so the resolver's `lstat` of `runs/<run-id>` fails with an
    error that is not a missing entry on Unix; `Build` returns an error naming
    the Run;
  - for every subtest, the home's and the repository's file hashes are
    identical before and after `Build`, and the default root is not created
    when the subtest did not create it.
- `internal/cli/runs_causes_test.go`, new:
  - `TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript`,
    asserting Surface Transcript 1 whole (stdout line by line, empty stderr,
    exit 0) with `<run id>` and `<digest>` substituted;
  - `TestRunsCausesReadsTheDiagnosticUnderARelativeArtifactRoot`, through the
    command, recorded `rel-artifacts`: the same stdout as the transcript.
- Unchanged and run as consumers: `TestBuildIgnoresALogOutsideTheArtifactDirectory`,
  `TestBuildReportsAttemptsTasksAndCorrectives`,
  `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText`,
  `TestBuildWindowAndMultipleRuns` and `TestBuildLeavesTheRunDatabaseUnchanged`
  in `internal/runcause`, and `TestRunsCausesPrintsAClassifiedWindow`,
  `TestRunsCausesPrintsJSON` and `TestBuildLeavesTheRunDatabaseUnchanged` in
  `internal/cli`.

## Build Order

1. The resolved read: `Request.HomeDir`, `diagnosticTail` on
   `store.ResolveRunArtifactRoot`, the command passing the home, their tests,
   and the `runs` guide with the Behavior Surface Record (no dependency inside
   this Spec; the resolver is on the default branch).
2. The Roundfix `runs` reference with the skill version, and **Artifact
   Root** in `CONTEXT.md` (depends on: 1, because they describe its
   behavior).
3. QA gate (depends on: 1, 2).

## Risks & Considerations

- After the release, `runs causes` classifies attempts it used to leave
  `unclassified`, so its counts on an existing Run Database change. The
  published measurement under `docs/references/` is not re-run by this Spec.
- A Run that recorded `""` but wrote under a root configured at the time
  still gives no evidence: its diagnostic lies outside the resolved root.
- A Roundfix Home reached through a symbolic link still gives no evidence,
  because the component walk refuses every link; this Spec keeps that rule.
- A Run Retention or `gc sanitize` removal racing with the report makes the
  diagnostic missing, which gives no evidence.

## Glossary

- changes: **Artifact Root**

## Decisions

- Reuse `store.ResolveRunArtifactRoot` rather than resolve in
  `internal/runcause`, so four callers cannot disagree about where a Run wrote
  (ADR-0301).
- Map the two "no root" errors to no evidence, and fail on any other resolver
  error, matching how the read already treats a missing diagnostic and an
  unreadable one.
- Take the home as a `Request` field filled by the command, the value that
  already opens the reader, rather than reading the environment inside the
  package.
- No ADR: the decision is ADR-0301's; the glossary carries the caller list.
