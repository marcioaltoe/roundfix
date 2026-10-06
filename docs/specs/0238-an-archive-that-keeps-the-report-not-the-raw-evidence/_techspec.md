---
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
prd: _prd.md
created: 2026-10-06
---

# An archive that keeps the report, not the raw evidence — Technical Spec

## Executive Summary

The cut lives in `internal/spec`, beside the archive's existing link pass
(ADR-0230), and runs inside `spec.Archive` after eligibility and before the
stamp. It walks `qa/evidence/`, hashes every file, writes
`qa/evidence-manifest.md`, points the Spec's evidence links at the manifest,
moves the Spec, and only then removes the evidence directory. The same
planner serves `roundfix archive <slug> --drop-evidence`, which cuts one
archived Spec on request and is a dry run unless `--apply` is passed. The
Delivery Queue's exact-move check learns to accept a cut whose manifest
lists exactly the dropped blobs. Two readers that made the repository's
gates depend on archived evidence are fixed first. The skills, the archive
user guide and two Baseline clauses then describe the cut. The primary
trade-off is provenance by digest and revision instead of bytes: a dropped
file is recoverable only from Git history, and the manifest is what makes it
findable and checkable there (ADR-0243).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  new names are the `--drop-evidence` and `--apply` flags, the manifest file
  name and schema string, and Go functions. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the cut reads and writes files
  and reads Git objects through the local `git` binary; no credential or
  network call is added, and every test uses temporary repositories, homes
  and fake runners. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0243 (this Spec): "The Archive
  Command now cuts the raw evidence as part of a normal archive", and
  "Already-archived Specs are cut only on request, one Spec at a time".
  It applies ADR-0215: "A removal is its own later change, and the
  maintainer's explicit approval of that removal is recorded in it". It
  extends ADR-0230, whose queue "accepts an archive commit whose only
  content changes besides the archive stamp are these rewrites". ADR-0154
  stands, ADR-0154: "The exception can waive the terminal QA Task's archive
  completion/evidence requirement without rewriting its status, Result or
  report", so an override archive keeps its QA files. ADR-0223's refusal for
  a file that "names an active Spec's directory" gains a counterpart for the
  archived evidence directory. ADR-0120 places the History Root, ADR-0120:
  "retired documentation is documentation". ADR-0210's Evidence Snapshots
  hash repository inputs, not evidence files, ADR-0210: "The record now
  holds one entry per declared input". ADR-0232's merge evidence supersedes
  a merged Spec's other commits, ADR-0232: "without a content comparison",
  so neither reads a dropped file. ADR-0229's operator archive and
  ADR-0237's Delivery Retry run the same Archive Command and exact-move
  check, so they inherit the cut unchanged, and ADR-0165's park of a
  blocking review after archive reads the review, not the evidence. ADR-0169's pre-PR review diffs the candidate, from which QA evidence is already omitted, and ADR-0240's partial policy is applied before the cut and reads only the report and declarations, so neither changes. ADR-0104, ADR-0167 and ADR-0196 bind the QA gate's rows, its outside evidence and the pre-PR review, which this Spec uses unchanged.
  ADR-0182 binds each Task's Verification to the facts its gate checks. ADR-0179 bounds the Governed Paths, ADR-0187, ADR-0189 and
  ADR-0233 bind the owned skills' versions, and ADR-0184 binds the Surface
  Transcripts. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183
  check consistency by citation and receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 changes three Governed Paths
  (the coverage record, its test and the repository-copy helper's file),
  task_02 changes `internal/spec/archive.go`, task_03 changes none, and
  task_04 changes the skill and Baseline ones. Express maintainer
  authorization: "Sim, como Spec" of 2026-10-06, beside the standing grants
  of 2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário"; "Autorizar os dois"); bounded files:
  `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`,
  `docs/agents/spec-routing.md`, `docs/references/coverage-record.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/spec/archive.go`, `internal/spec/coverage_test.go`,
  `skills/archive-spec/SKILL.md`, `skills/baseline_skill_contract_test.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0238-an-archive-that-keeps-the-report-not-the-raw-evidence/_authorization.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Coverage record | `collectCoverageRecord` in `internal/spec/coverage_test.go`; `docs/references/coverage-record.json` | Counts no package under `roundfix/docs/`; the three archived-evidence packages leave the record |
| Repository copy | `copyTrackedRepository` in `skills/baseline_skill_contract_test.go` | Skips a tracked path missing from the working tree |
| Evidence cut | new `internal/spec/archive_evidence.go` | Plans, renders and applies the cut; parses the manifest |
| Archived cut | new `internal/spec/archive_evidence_cut.go` | Cuts one archived Spec on request |
| Archive | `Archive`, `ArchiveRequest`, `ArchiveResult`, `ArchiveLinksMatch` in `internal/spec/archive.go` | Calls the cut for a normal archive; counts manifest links with the rewrites; accepts an evidence link that now reaches the manifest |
| Archive Command | `internal/cli/archive.go` | Passes the revision and source; reports the cut; `--drop-evidence` and `--apply` |
| Archived-evidence pins | new `internal/speccheck/archived_evidence_paths.go` | Files other than Markdown that name an archived Spec's evidence directory |
| Delivery exact move | `archiveCommitIsExact` in `internal/cli/deliver_workflow.go` | Accepts dropped evidence that the archived manifest lists exactly |
| Guidance | archive-spec, qa-gate and Roundfix skills; `docs/user-guide/commands/archive.md`; two clauses of `spec-workflow.json` | Describe the cut, the manifest and the request-only cut |

Checked and unchanged: the review scope (`internal/cli/review.go`) already
omits any path under a Spec's `qa/evidence/`, and `review_scope_test.go`
writes its evidence paths into a disposable repository, so it reads no
archived file. Merged-Run reconciliation (`internal/worktree/merged_head.go`)
counts a changed path under an archived Spec's directory as superseded
without reading it. QA eligibility, settlement and Evidence Snapshots read
the QA Report and declarations only. `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive`
reads archived QA Reports, which the cut keeps.

## Implementation Design

### Interfaces

```go
// internal/spec/archive_evidence.go
const (
	EvidenceManifestName   = "evidence-manifest.md" // in the Spec's qa/ directory
	EvidenceManifestSchema = "roundfix/qa-evidence-manifest/v1"
)
type EvidenceManifestEntry struct { Path string; Bytes int64; SHA256 string }
type EvidenceManifest struct {
	DroppedOn, Revision, Source string
	Files                       int
	Bytes                       int64
	Entries                     []EvidenceManifestEntry
}
func ParseEvidenceManifest(content []byte) (EvidenceManifest, error)
func RenderEvidenceManifest(manifest EvidenceManifest) []byte

type ArchivedEvidenceCutRequest struct {
	SpecsRoot   string; BuiltInRoot bool; Slug string
	Revision    string; Source string; CutOn time.Time; Apply bool
}
type ArchivedEvidenceDisposition string // "cut", "would-cut", "nothing", "kept-override", "kept-superseded"
type ArchivedEvidenceCutResult struct {
	Disposition ArchivedEvidenceDisposition
	SpecDir, ManifestPath string; Files int; Bytes int64; Links int
}
// internal/spec/archive_evidence_cut.go
func CutArchivedEvidence(req ArchivedEvidenceCutRequest) (ArchivedEvidenceCutResult, error)
```

```go
// internal/spec/archive.go — added fields only
type ArchiveRequest struct { /* existing */ EvidenceRevision, EvidenceSource string }
type ArchiveResult struct { /* existing */ DroppedEvidenceFiles int; DroppedEvidenceBytes int64 }

// internal/speccheck/archived_evidence_paths.go
func ArchivedEvidencePathPins(repoRoot, specsRoot, slug string) ([]SpecPathPin, error)
```

```text
1. The cut set is every entry under <spec>/qa/evidence/ other than a directory, walked without following links. An entry that is not a regular file, or a path containing a newline, carriage return, tab, "|" or a backtick, refuses before any file changes and names the path relative to the Spec.
2. A manifest is written only when the cut set is non-empty. Its rows are sorted by path in byte order; each path is relative to the Spec directory with "/" separators, bytes is decimal and sha256 is lowercase hex of the file's content.
3. The manifest's front matter holds, in order: schema, dropped_on (YYYY-MM-DD), revision, source (the repository-relative Spec directory at that revision), files and bytes. Its body is "# QA evidence manifest", one paragraph saying the files were dropped at archive and can be read with git show at the revision, and the table "| path | bytes | sha256 |" with each path in a code span.
4. In every regular Markdown file of the Spec outside qa/evidence/, a relative link destination that the archive's link scanner reads and whose target resolves to qa/evidence or below becomes the relative path from that file's directory to qa/evidence-manifest.md, without query or fragment. Link text, title and angle-bracket form are kept. These rewrites count in RewrittenLinks with ADR-0230's.
5. A normal archive of a Spec with a Task Graph and no QA override cuts. A QA Archive Override archive and a superseded Spec without a Task Graph cut nothing and move exactly as today. A Spec without files under qa/evidence/ archives byte-for-byte as today.
6. Every archive refusal of today happens first; a cut refusal (Invariant 1, or a missing EvidenceRevision when the cut set is non-empty) happens before the stamp. The evidence directory is removed only after the move succeeds. A removal failure after the move returns an error naming the leftover directory; the archive itself stands.
7. ArchiveLinksMatch also accepts a destination that changed from a target under the active Spec's qa/evidence/ to the archived Spec's qa/evidence-manifest.md. Its signature is unchanged.
8. archiveCommitIsExact accepts source-tree entries under qa/evidence/ that are absent from the destination only when the destination adds qa/evidence-manifest.md, the manifest parses, its source equals the active Spec path, its rows list exactly those entries, and each row's bytes and sha256 equal the parent commit's blob. Every other rule of the exact move is unchanged.
9. CutArchivedEvidence resolves <archive root>/<slug>. It refuses when the Spec is active or not archived, and when the evidence directory holds files and the manifest already exists. It returns kept-override for a PRD with qa_override true, kept-superseded for a Spec without _tasks.md that has _supersession.md, and nothing for a Spec without files under qa/evidence/. Otherwise it returns would-cut without writing, or cut after writing the rewrites and the manifest and removing the directory when Apply is set. Dry run and apply plan the same set.
10. --drop-evidence refuses, before planning, when ArchivedEvidencePathPins finds a file other than Markdown outside the Spec roots and the History Root that names <archive-relative>/<slug>/qa/evidence, and when git status reports any change under that directory. --apply records the HEAD revision.
11. No path under docs/ is a package of the coverage record. copyTrackedRepository skips a path that git ls-files lists but that is absent from the working tree; every other copy error still fails the test.
```

### Data Models

`qa/evidence-manifest.md` is the only new artifact. An example of its form:

```markdown
---
schema: roundfix/qa-evidence-manifest/v1
dropped_on: 2026-10-07
revision: 0123456789abcdef0123456789abcdef01234567
source: docs/specs/0238-an-archive-that-keeps-the-report-not-the-raw-evidence
files: 2
bytes: 1536
---

# QA evidence manifest
...
| path | bytes | sha256 |
| --- | ---: | --- |
| `qa/evidence/run/R01.txt` | 512 | <64 hex> |
```

No Run Database, Task file, QA Report front matter or PRD stamp field
changes.

### API Contracts

1. API Contract: `roundfix archive <slug>` on an eligible Spec with files
   under `qa/evidence/` drops them, writes the manifest, points the Spec's
   evidence links at it, and appends
   `; dropped <f> QA evidence file(s) (<b> bytes) listed in qa/evidence-manifest.md`
   to its confirmation, after the existing rewrite suffix. Without evidence
   its output is unchanged.
2. API Contract: `roundfix archive <slug> --drop-evidence` on an archived
   Spec prints what the cut would do and changes no file; with `--apply` it
   cuts and prints what it did; both exit 0. A kept or already-cut Spec
   prints one line and exits 0.
3. API Contract: `--drop-evidence` exits 2 through Preflight Validation for
   an active Spec, a missing archived Spec, a pinned evidence directory,
   uncommitted changes under it, a manifest beside remaining evidence, or an
   entry Invariant 1 refuses. `--apply` without `--drop-evidence`, and
   `--drop-evidence` with `--qa-override`, `--approval` or `--reason`, are
   refused as usage errors with exit 2.
4. API Contract: the Delivery Queue treats an archive commit with a cut as
   an exact Spec move under Invariant 8, and as not exact otherwise.

### Surface Transcripts

1. Surface Transcript: a normal archive drops the evidence.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   archived <slug> -> docs/history/specs/<slug>; rewrote <n> relative link(s); dropped <f> QA evidence file(s) (<b> bytes) listed in qa/evidence-manifest.md
   stderr:
   exit: 0
   ```

2. Surface Transcript: a dry run on an archived Spec.

   ```transcript
   $ roundfix archive <slug> --drop-evidence
   stdout:
   would drop <f> QA evidence file(s) (<b> bytes) from docs/history/specs/<slug>/qa/evidence, rewrite <n> relative link(s) and write docs/history/specs/<slug>/qa/evidence-manifest.md
   stderr:
   exit: 0
   ```

3. Surface Transcript: the cut applied.

   ```transcript
   $ roundfix archive <slug> --drop-evidence --apply
   stdout:
   dropped <f> QA evidence file(s) (<b> bytes) from docs/history/specs/<slug>/qa/evidence, rewrote <n> relative link(s) and wrote docs/history/specs/<slug>/qa/evidence-manifest.md
   stderr:
   exit: 0
   ```

4. Surface Transcript: a second apply.

   ```transcript
   $ roundfix archive <slug> --drop-evidence --apply
   stdout:
   <slug> has no QA evidence to drop
   stderr:
   exit: 0
   ```

5. Surface Transcript: an override archive keeps its evidence.

   ```transcript
   $ roundfix archive <slug> --drop-evidence
   stdout:
   <slug> keeps its QA evidence: archived with a QA Archive Override
   stderr:
   exit: 0
   ```

6. Surface Transcript: an active Spec is refused.

   ```transcript
   $ roundfix archive <slug> --drop-evidence
   stdout:
   stderr:
   Preflight failed
   ...
   Spec "<slug>" is active; a normal archive drops its QA evidence
   ...
   exit: 2
   ```

A superseded Spec prints
`<slug> keeps its QA evidence: superseded without a Task Graph`.

## Coverage Map

- Goal 1 → evidence cut and manifest (Invariants 1-3, 5, 6; API Contract 1)
- Goal 2 → link rewrite (Invariant 4)
- Goal 3 → exact move (Invariants 7, 8; API Contract 4)
- Goal 4 → `CutArchivedEvidence`, `--drop-evidence` (Invariants 9, 10; API Contracts 2, 3)
- Goal 5 → coverage record and repository copy (Invariant 11)
- User Story 1 → Archive, Archive Command
- User Story 2 → Evidence Manifest, link rewrite
- User Story 3 → `archiveCommitIsExact`, `ArchiveLinksMatch`
- User Story 4 → `CutArchivedEvidence`, Archive Command flags
- User Story 5 → coverage record, `copyTrackedRepository`
- Core Feature 1 → Invariants 1, 5, 6
- Core Feature 2 → Invariants 2, 3
- Core Feature 3 → Invariant 4; API Contract 1; Surface Transcript 1
- Core Feature 4 → Invariants 1, 5
- Core Feature 5 → Invariants 7, 8; API Contract 4
- Core Feature 6 → Invariants 9, 10; API Contracts 2, 3; Surface Transcripts 2-6
- Core Feature 7 → Invariant 11
- Core Feature 8 → skills, user guide, Baseline clauses (Exact texts)
- Success Metric 1 → task_02 tests
- Success Metric 2 → task_02 delivery tests
- Success Metric 3 → task_03 tests
- Success Metric 4 → task_01 tests; QA row in a disposable clone
- Success Metric 5 → QA dry run over the History Root

## Integration Points

Git only, through the local `git` binary already used by the Archive Command
and the Delivery Queue: `rev-parse HEAD`, `status --porcelain`, `ls-tree`
and reading blobs. No external system.

## Testing Approach

- `internal/spec`: a new `archive_evidence_test.go` archives temporary Specs
  through `Archive` and `CutArchivedEvidence`: evidence with nested
  directories and a binary file; a QA Report, a Task file and the PRD
  linking into evidence by inline, reference and angle-bracket forms, with a
  fragment; a link to a sibling outside the Spec (ADR-0230's rewrite still
  applies); no evidence; an override and a superseded Spec; a symbolic link
  and a path with "|". It parses the written manifest with
  `ParseEvidenceManifest` and recomputes sizes and digests. It also covers
  `ArchiveLinksMatch` with an evidence link that now reaches the manifest,
  and one that reaches anything else.
- `internal/cli`: `archive_evidence_test.go` runs the Archive Command in a
  temporary repository and asserts Surface Transcript 1;
  `deliver_archive_evidence_test.go` builds parent and archive commits and
  asserts `archiveCommitIsExact` for an exact cut, a manifest that omits a
  file, a wrong digest and an extra dropped path;
  `archive_drop_evidence_test.go` covers Surface Transcripts 2-6, the
  superseded line, the pin refusal, uncommitted evidence, the usage errors,
  and that a dry run leaves every file byte-identical.
- `internal/spec/coverage_test.go`: a unit test of the package filter, and
  `TestCoverageEquivalence` against the re-recorded record.
- `skills`: a new `copy_tracked_repository_test.go` builds a temporary
  repository, deletes one tracked file from its working tree and copies it.
- `internal/baseline`: a new `evidence_cut_clause_test.go` asserts both
  added sentences in the embedded clauses, the two formatter goldens and
  this repository's two generated guides.
- Existing archive, delivery, review-scope and Baseline tests keep passing
  unchanged: no existing fixture writes a file under `qa/evidence/` before an
  archive.

## Build Order

1. Readers: the coverage record counts no package under `docs/`, the record
   is re-recorded, and the repository copy skips a missing tracked file.
2. The cut at archive: `archive_evidence.go`, `Archive`, the link rewrite,
   `ArchiveLinksMatch`, the Archive Command's confirmation and the Delivery
   Queue's exact move. Independent of step 1; no shared file.
3. The request-only cut of an archived Spec: `CutArchivedEvidence`, the pin
   check, `--drop-evidence` and `--apply`, and the archive user guide
   (depends on: 2).
4. Guidance: the archive-spec, qa-gate and Roundfix skills, the two Baseline
   clauses, `make baseline-digests` and the Managed Refresh (depends on: 3).
5. Final QA gate (depends on: 1, 4).

Step 4 follows step 3 because it documents step 3's command. Step 1 is
needed before any cut of this repository's History Root, which this Spec
does not perform; the QA gate's dry run over the History Root depends on it.

## Risks & Considerations

- A conclusion that lived only in an evidence file disappears from the
  archive. The qa-gate skill now asks for each row's observed result in the
  report, and the bytes stay in Git history at the recorded revision.
- Evidence that was never committed cannot be recovered. The Delivery Queue
  commits QA evidence with the QA Task before archive, and `--apply` refuses
  uncommitted evidence.
- The `### QA settlement` table, identical in three skills, says a passing
  archive keeps "The Spec and its QA report and evidence". This Spec does
  not edit that section. Each skill gains its own section saying that what
  archive keeps of the evidence is its manifest.
- Specs 0235, 0236 and 0237 are delivered first and touch
  `spec-workflow.json`, its derived files, `internal/cli/deliver_workflow.go`,
  `internal/cli/archive_test.go` and the Roundfix skill. This Spec declares
  new test files instead of editing `archive_test.go`, regenerates derived
  files on top of theirs, and re-records skill versions (ADR-0233).
- Reading every dropped blob in the exact-move check costs one batched
  `git cat-file` per archive, bounded by the evidence of one Spec.

## Exact texts

Each Baseline sentence is appended after the clause's existing last
sentence; every existing sentence stays byte-identical:

- `clause.spec.project-constraints-05-legacy-and-ownership`: "The one
  exception is the runtime's QA evidence cut: an archive, or an explicit
  request to cut an archived Spec, may drop the Spec's `qa/evidence/`
  directory, write the evidence manifest that lists each dropped file, and
  point the Spec's links into that directory at the manifest."
- `clause.spec.keep-artifacts-in-spec-folder`: "A normal archive through a
  runtime with an evidence cut, such as `roundfix archive <slug>`, drops
  `qa/evidence/` and keeps every QA Report beside `qa/evidence-manifest.md`,
  which names each dropped file's path, size and SHA-256 and the revision
  that holds it; record each row's observed result in the QA Report, and
  keep material a later reader needs in `references/` or
  `docs/references/`, never only in `qa/evidence/`."

The skills gain sections outside `### QA settlement`: archive-spec a
`## Raw QA evidence at archive` section before `## Unarchive`; qa-gate a
`## Raw evidence does not outlive the archive` section before
`## Anti-patterns`; the Roundfix archive reference a `### Raw QA evidence`
subsection at the end of `## Archive Command`. Each names
`qa/evidence-manifest.md`; the archive-spec and Roundfix ones also name
`roundfix archive <slug> --drop-evidence --apply`.

## Vocabulary Contract

- emits: `internal/cli/archive.go`
  pattern: `evidence-manifest.md`
  documented-in: `docs/user-guide/commands/archive.md`

`CONTEXT.md` is not edited here. Its **Archive Command** entry says the
command "moves the whole Spec", and **Evidence Manifest** is a candidate
term; the QA gate records both for the glossary.

## Decisions

- The cut covers `qa/evidence/` only and runs inside the normal archive.
  See ADR-0243.
- Markdown manifest in `qa/`, not JSON, so links can reach it and readers
  can read it.
- `roundfix archive --drop-evidence` rather than a new history command:
  the same cut, the same command, one help text.
- Override and superseded archives keep their files.
- The coverage record excludes packages under `docs/` instead of keeping
  archived evidence packages listed.
