---
spec: 0242-an-archive-that-leaves-an-archive-record
prd: _prd.md
created: 2026-10-06
---

# An archive that leaves an Archive Record — Technical Spec

## Executive Summary

The record and the cut live in `internal/spec`, beside today's archive.
`spec.Archive` keeps every eligibility rule and refusal. Instead of moving the
folder, it builds an `ArchiveRecord` from the Spec's files, writes it to
`<archive-root>/<slug>.md`, copies any promoted files to `docs/references/`
and removes the folder. A single resolver, `spec.ReadArchivedSpec`, answers
"is this Spec archived, and how" from a record or from a legacy folder. Every
runtime reader moves to it.

Readers that need more than the record read the pre-archive tree from Git at
the record's `source_revision`, which the Run's or the candidate's own branch
holds. The Delivery Queue's exact-archive proof becomes an exact retirement:
the commit removes the folder and adds a record that agrees with it.
Repository tests that characterize the archived corpus read it from Git at the
pinned commit `40a7893d872c8a6705f6d7745e6efe430ae9deeb`.

`roundfix archive <slug> --plan` lists what the cut removes, and Jev advises
on candidate files through the Spec judge's transport. `--promote <path>`
copies a file upstream in the archive change.

The primary trade-off is provenance by revision instead of bytes. A removed
file is recoverable only from Git, and the record is what makes it findable
there (ADR-0247).

## Project Constraints

- Identifier strategy: applicable — new names are the record file
  `<slug>.md`, its schema string `roundfix/archive-record/v1`, the judgment
  kind `archive-value`, the flags `--plan` and `--promote`, and Go
  identifiers. Spec, Task and ADR identifiers keep their schemes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — only `--plan` can send a request,
  through `internal/judge`'s existing client, key variables
  (`openrouterkey.StageJudge`), monthly ceiling and judge log. With no key or
  a reached ceiling it sends nothing and prints the skip reason. Every test
  injects an `http.RoundTripper` and a temporary home. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0247 (this Spec) decides that an
  archive leaves an Archive Record and the folder stays in Git, that Jev
  advises and never gates, and that repository tests read the pre-cut corpus
  at a pinned commit. ADR-0247 records, for archives from now on, the approval
  ADR-0215 requires, ADR-0215: "A removal is its own later change, and the
  maintainer's explicit approval of that removal is recorded in it". ADR-0120
  places the History Root, ADR-0120: "retired documentation is documentation",
  which holds the records. ADR-0154: "The exception can waive the terminal QA
  Task's archive completion/evidence requirement without rewriting its status,
  Result or report", so the record carries every override field. ADR-0230's
  link pass keeps serving folders archived earlier. ADR-0121 records the
  Baseline history layout's relocations of those folders, ADR-0121: "Archive
  relocations are therefore recorded as their own ordered ledger of source,
  destination, and content identity". ADR-0223's refusal for a file that
  "names an active Spec's directory" runs before the cut, unchanged. ADR-0229
  and ADR-0237 reach the cut through the same Archive Command. ADR-0232 grants
  merge evidence when the default branch "holds its archived `_prd.md`", and
  the record now stands for that file. ADR-0165 parks a blocking review after
  archive, and the review must still find the archived Spec; ADR-0169's
  merge-base diff is the diff from which the review omits the removed folder.
  ADR-0193 counts an archived prerequisite. ADR-0210 hashes declared inputs
  and reads no archived folder. ADR-0179 bounds the Governed Paths, and
  ADR-0187, ADR-0189 and ADR-0233 bind the owned skills' versions. ADR-0184
  binds the Surface Transcripts, ADR-0240 the QA partial, and ADR-0182 each
  Task's Verification. The gate is bound by ADR-0080, ADR-0091, ADR-0096,
  ADR-0097, ADR-0104 and ADR-0167, and ADR-0196 validates a pre-PR finding
  before it parks. ADR-0194 and ADR-0195 cite ADR-0097 but decide what a QA
  row records and when it is observed again; this Spec changes neither, so
  neither applies. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and
  ADR-0183 check consistency by citation and receipt. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — 24 of the declared files are Governed Paths.
  They were measured with `GovernedPath` through a `go test -overlay` probe
  that wrote nothing, and each is bounded in `_authorization.md`. The express
  authorization comes from the maintainer: "Concedo" on 2026-10-06, beside
  the standing grants of 2026-09-30 ("considere autorizado a ajustar todas as
  skills se necessário"; "Autorizar os dois"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0242-an-archive-that-leaves-an-archive-record/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/derived_ownership_test.go`, `internal/baseline/derived_regeneration_repocontract_test.go`, `internal/spec/archive.go`, `internal/spec/archive_layout_characterization_test.go`, `internal/spec/archive_test.go`, `internal/spec/coverage_test.go`, `internal/speccheck/backlog.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `skills/archive-spec/SKILL.md`, `skills/owned_skill_edit_repocontract_test.go`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Readers

The inventory was measured on 2026-10-06 at `40a7893d`:

- `git grep -l -I -e docs/history -- ':!docs/history/**'` named 97 tracked
  files: 44 Go and 53 others.
- `git grep -l -e ArchiveDir -e ArchiveSpecRoot -e history/specs -- '*.go'`
  named 59 Go files.
- The ablation (`_prd.md` → Acceptance evidence) found the readers that
  break on the real repository.

| Kind | Count | Readers | Disposition |
| --- | ---: | --- | --- |
| Archive itself | 2 files | `spec.Archive`, Archive Command | task_01: writes the record, removes the folder |
| Delivery Queue archive stage | 1 file, 5 functions | `Archive`, `archivePaths`, `archiveDiffIsExact`, `archiveCommitIsExact`, `reconcileArchiveCommit` in `deliver_workflow.go` | task_01: exact retirement, legacy proof kept |
| Supersede | 1 | `knownDelivererSpec` | task_01, because its tests run the real archive |
| Other queue readers | 4 | `InspectItem`, `UnmetPrerequisites`, `validateDeliveryPrerequisites`, `ProveReviewCorrection` | task_02 |
| Reconcile | 7 functions | `chooseMergedHead`, `provenDeliveryEvidence`, `specArchivedAtMergedHead`, `taskCompletedAtMergedHead`, `taskCompletedInRoots`, `dirtyPathsInArchivedSpec` in `merged_head.go`; `qaReportDirectories` / `archivedQAReportDirectory` in `worktree.go` | task_02 |
| Pre-PR review | 3 | `reviewCandidateSpecContexts`, `reviewSpecQAOverride`, `reviewScopedDiff` | task_02 |
| Run causes | 2 | `runs_causes.go`, `runcause.findGraph` | task_02 |
| Spec check | 3 | `detectFindingsConsistency` and `repositoryDirectoryNames`, `detectBacklogPromotion`, `detectTaskContextReferences` | task_03 |
| Spec audit | 3 | `claimedArtifacts`, `archivedSpecArtifactPath`, `validateSpecAuditSlug` | task_03 |
| Suite guard authority | 1 | `readSanctionedRegenerations`: 85 archived declarations | task_03: record `regeneration`, legacy folders |
| Repository tests failing in the ablation | 22 tests in 15 files | `make verify`: 17, `make verify-docs`: 5 (list below) | task_03 |
| Coverage record | 1 | `docs/references/coverage-record.json` lists packages under `docs/history` | task_03 |
| Checked user-guide links | 7 links in 5 files | `TestUserGuideLinksResolve` | task_03 |
| Guidance naming the archive move | 9 | archive-spec, qa-gate and Roundfix skills; archive and context-driven guides; two Baseline clauses and their generated guide; `CONTEXT.md` | task_04 |
| Unaffected, checked | 31 | see below | none |

These 22 tests failed in the ablation:

- `TestCurrentRecordPermitsItsDeclaredOperations`
- `TestQuestionFileLoads`, `TestGroupingQuestionIsTheMeasuredOne`
- `TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily`
- `TestRepositorySpecCorpusStillLoads`,
  `TestQAGateLegacyArchivedManifestsLoadUnchanged`
- `TestArchivedPassCorpusRemainsArchiveEligible`,
  `TestArchivedQAOverrideCorpusIncludesFailedSpec`
- `TestSpec0058ReplayArchivesDeclaredUnreachableRelease`,
  `TestSpec0058ReplayRefusesUnmatchedBlockedRow`,
  `TestSpec0058ReplayReportsWronglyDeclaredReachableRow`
- `TestArchivedQAReportCorpusRemainsReadable`, `TestCoverageEquivalence`
- `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive`
- `TestReceiptCharacterizationKeepsEveryParsedClaim`
- `TestThisSpecsSurfaceTranscriptsAreWellFormed`,
  `TestThisSpecsTranscriptsHaveImplementationAndGateReferences`
- `TestMeasuredSanctionedOwnershipMatchesRecords`,
  `TestDeclaredStepRegenerationAndFrozenBoundaries`
- `TestEveryBoundedPathIsGoverned`
- `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`
- `TestUserGuideLinksResolve`

These readers were checked and stay unchanged:

- Path exclusions: `ActiveSpecPathPins`, `.coderabbit.yaml`, the review scope
  roots.
- Revision reads at a commit before the archive: the delivery authorization
  in `deliver_workflow.go` `Authorization`, and
  `discoverMechanicalAuthorizationPaths`.
- `internal/baseline/history_layout.go`. It relocates any file under a legacy
  `_archived/` tree, records included.
- `TestCleanupHistoricalGrantEvidence`, which already reads an ancestor
  commit.
- `.secondbrain-export`, which Spec 0243 revisits.
- Prose that names `docs/history` as a rule and stays true for records: 16
  ADRs, `docs/agents/docs-layout.md`, `docs/agents/specific-repository.md`,
  and the brainstorming, write-idea, write-prd, write-tasks and
  write-techspec skills with their mirrors. A record also lives under
  `docs/history/specs/` and starts with `NNNN`.
- Prose paths that resolve in Git: `CHANGELOG.md` and four
  `docs/references/` files.
- Tests that read only temporary fixtures, in `internal/cli`,
  `internal/worktree`, `internal/speccheck` and `internal/spec`. They build
  legacy folders, which stay readable.
- `TestTaskAcceptance` measurement roots, which run only behind a flag.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Archive Record | new `internal/spec/archive_record.go` | Format, build, render, parse, resolver |
| Archive | `internal/spec/archive.go` | Writes the record, copies promotions, removes the folder; the link pass no longer runs for a cut |
| Archive Command | `internal/cli/archive.go` | Passes the revision, refuses a folder that differs from `HEAD`, prints the cut; `--plan`, `--promote` (task_04) |
| Exact retirement | `internal/cli/deliver_workflow.go` | Archive stage and reconcile accept a record commit; the legacy proof stays |
| Runtime readers | `deliver*.go`, `review*.go`, `supersede.go`, `runs_causes.go`, `runcause`, `worktree` | Read `ReadArchivedSpec`, then Git at `source_revision` |
| Check readers | `speccheck`, `specaudit`, `spec_check.go`, `suiteguardcontract` | Record stems, record regenerations |
| Pinned history | new `internal/gittest/pinned_history.go` | Materializes paths of the pinned commit for tests |
| Archive Advice | new `internal/judge/archive_advice.go`, `questions.json` | `archive-value` judgment over candidate files |
| Guidance | skills, guides, Baseline clauses, `CONTEXT.md` | Describe the record, the cut and the advice |

## Implementation Design

### Interfaces

```go
// internal/spec/archive_record.go
const ArchiveRecordSchema = "roundfix/archive-record/v1"
const ArchiveRecordTargetBytes = 2048

type ArchiveDisposition string // "pass", "partial", "qa-override", "superseded"

type ArchiveRegeneration struct{ Command string; Outputs []string }

type ArchiveRecord struct {
	Spec, Title, Created, Archived string
	Disposition                     ArchiveDisposition
	Source, SourceRevision          string // repository-relative active path; 40-hex HEAD at archive
	QATask, QAReport, QAVerdict     string
	Unproven                        []string
	QAOverride                      *QAArchiveOverrideRecord // nil unless disposition is qa-override
	SupersededBy                    string
	ADRs, Sources, Promoted         []string
	Regeneration                    []ArchiveRegeneration
	PullRequest, DeliveryCommit     string // optional; empty at archive, filled by migration when known
	Outcome                         string
}
type QAArchiveOverrideRecord struct{ Approval, Reason, QAOutcome, QATaskStatus, Revision string }

func ArchiveRecordPath(archiveRoot, slug string) string // <archiveRoot>/<slug>.md
func BuildArchiveRecord(in ArchiveRecordInput) (ArchiveRecord, error)
func RenderArchiveRecord(record ArchiveRecord) ([]byte, error)
func ParseArchiveRecord(content []byte) (ArchiveRecord, error)

type ArchivedForm string // "record", "folder"
type ArchivedSpec struct {
	Slug   string
	Form   ArchivedForm
	Path   string        // record file or legacy folder
	Record ArchiveRecord // for a folder: built from its _prd.md, _tasks.md and newest QA Report
}
// ReadArchivedSpec returns ErrNotArchived when neither form exists, and an
// error when both exist.
func ReadArchivedSpec(archiveRoot, slug string) (ArchivedSpec, error)
func ArchivedSpecSlugs(archiveRoot string) ([]string, error) // folder names and record stems, sorted, unique
var ErrNotArchived = errors.New("Spec is not archived")

// internal/spec/archive.go — added fields
type ArchiveRequest struct { /* existing */ SourceRevision string; Promote []string; RepositoryRoot string }
type ArchiveResult struct { /* existing */ RecordPath string; RemovedFiles int; RemovedBytes int64; Promoted []string }

// internal/gittest/pinned_history.go
const PinnedHistoryRevision = "40a7893d872c8a6705f6d7745e6efe430ae9deeb"
func PinnedHistory(t testing.TB, repoRoot string, pathspecs ...string) string // temp dir holding those paths at the pin; t.Skip when the commit is absent

// internal/judge/archive_advice.go
type ArchiveAdviceRequest struct {
	RepoRoot, SpecDir, Spec string
	Files                   []string // Spec-relative candidate files
	Keys                    map[string]string
	HomeDir                 string
	Transport               http.RoundTripper
	Now                     func() time.Time
}
type ArchiveAdvice struct{ File, Outcome string; Choice *string; Probabilities map[string]float64; Confidence *float64; Reason *string; Model *string }
type ArchiveAdviceReport struct{ Model string; Skipped, Stopped *string; Advice []ArchiveAdvice; Calls int; CostUSD, MonthCostUSD, MonthCeilingUSD float64 }
func AdviseArchive(ctx context.Context, q Questions, req ArchiveAdviceRequest) (ArchiveAdviceReport, error)
```

### Invariants

1. The record path is `ArchiveRecordPath(ArchiveSpecRoot(specsRoot, builtInRoot), slug)`:
   `docs/history/specs/<slug>.md` for the built-in root and
   `<spec-root>/_archived/<slug>.md` for any other.
2. `BuildArchiveRecord` derives each field from the Spec as follows. Nothing
   is invented.
   - `title` is the PRD's first H1, and `created` is its `created`.
   - `disposition` is `superseded` for a Spec without a Task Graph,
     `qa-override` with an override, and otherwise the QA verdict, `pass` or
     `partial`.
   - `qa_task` is the `_tasks.md` `qa` id. `qa_report` is the newest QA
     Report's file name and `qa_verdict` its verdict. Both are empty when no
     report exists.
   - `unproven` holds today's stamp values, and the override fields are
     today's stamp values.
   - `superseded_by` comes from `_supersession.md`.
   - `adrs` are the `ADR-NNNN` tokens in the PRD's `## Decisions` section,
     unique, in first-appearance order.
   - `sources` are the `path` column of `references/_index.md`, in table
     order.
   - `regeneration` is `suiteguardcontract.ParseSanctionedRegenerations` of
     `_authorization.md`, with the outputs it enumerates.
   - `promoted` holds the destinations of the promoted files.
   - `outcome` is the first paragraph of the newest QA Report's `## Outcome`
     section when one exists, else the first paragraph after the PRD's H1.
     For a superseded Spec it is the supersession reason. Newlines collapse
     to single spaces.
3. `RenderArchiveRecord` writes YAML front matter in this key order:
   `schema`, `spec`, `title`, `status: archived`, `created`, `archived`,
   `disposition`, `source`, `source_revision`, `qa_task`, `qa_report`,
   `qa_verdict`, `unproven`, the six `qa_override*` keys only for
   `qa-override`, `superseded_by` only for `superseded`, `adrs`, `sources`,
   `regeneration`, `promoted`, and `pull_request` and `delivery_commit` only
   when set. The body is `# <title>`, a blank line and the outcome paragraph.
   When the rendering exceeds `ArchiveRecordTargetBytes`, only the outcome
   is shortened: to the last sentence end that fits, else to a rune boundary
   followed by `…`. Every other field is written whole. `ParseArchiveRecord`
   of a rendering returns the same record. It refuses an unknown schema, a
   missing required key and a `spec` that differs from the file stem when
   read through `ReadArchivedSpec`.
4. `Archive` keeps every refusal of today, in today's order, and adds these
   before any file changes:
   - an empty or non-hex `SourceRevision`;
   - a record or legacy folder already present for the slug;
   - a promoted path that is outside the Spec, is not a regular file, is a
     core artifact (`_*.md` at the Spec root, `task_*.md`, a QA Report) or
     has an occupied destination under `docs/references/`.
5. A cut then writes the record, copies each promoted file byte-identically
   to `<RepositoryRoot>/docs/references/<basename>`, and only then removes
   the Spec folder. A failure before the removal leaves the folder whole and
   removes the record and the copies it wrote. A removal failure returns an
   error naming the leftover folder, and the record stays. The PRD is not
   stamped and no link is rewritten. The cut applies to normal, override and
   superseded archives alike.
6. The Archive Command passes `HEAD` as `SourceRevision`, from the Spec
   Root's repository when it is external, as the override revision does
   today. Before `spec.Archive` runs, it refuses through Preflight Validation
   when `git status --porcelain --untracked-files=all -- <spec dir>` reports
   anything.
7. `ReadArchivedSpec` reads the record when `<slug>.md` exists, else the
   legacy folder's `_prd.md` (status `archived`), `_tasks.md`, newest QA
   Report and `_supersession.md`, through the same builder. It refuses when
   both forms exist. `ArchivedSpecSlugs` lists both forms, and a missing
   archive root yields an empty list.
8. The exact-retirement proof is in `archiveCommitIsExact` and the archive
   stage's `archiveDiffIsExact`. It accepts a commit whose parent holds
   `<source>/_prd.md` when all of the following hold:
   - the commit deletes every path under `<source>/`;
   - it adds exactly the record path and the record's `promoted` paths;
   - it changes nothing else;
   - the record parses, with `spec` equal to the slug, `source` equal to
     `<source>`, `source_revision` equal to the parent, and `title` and
     `created` equal to the parent's `_prd.md`;
   - each promoted blob equals a deleted blob under `<source>/`.

   A commit that adds `<destination>/` instead keeps today's proof
   (ADR-0230). Anything else is not exact.
9. Readers that need a Task file, a TechSpec, a QA Report or a Task Context
   path read them with `git show <source_revision>:<source>/<path>` when
   `git cat-file -e <source_revision>^{commit}` succeeds in the repository at
   hand. Otherwise they fall back to the record and decide conservatively:
   a worktree is kept, never released, on missing knowledge.
10. Task completion of an archived Spec follows from the record. Every Task
    other than the QA Task is completed. The QA Task is completed unless the
    disposition is `qa-override` with a `qa_override_qa_task_status`.
11. The pre-PR review's diff omits the deletions under an archived Spec's
    `source` and keeps the record. A record added under the archive root
    marks its slug archived. The review's Spec context reads the PRD and
    TechSpec at `source_revision`.
12. `suiteguardcontract` reads `regeneration` from each record's front matter
    with its own YAML read, and does not import `internal/spec`. It also
    reads legacy archived authorizations and active ones as today.
13. `PinnedHistory` uses `git archive` of the pinned commit, limited to the
    pathspecs, and extracts the result with `archive/tar` into `t.TempDir()`.
    It writes nothing in the repository. It skips with the commit named when
    the commit is absent, the same way `TestCleanupHistoricalGrantEvidence`
    does today.
14. `AdviseArchive` judges only candidate files, never a core artifact or a
    path under `qa/evidence/`. It reads at most the first 6,000 bytes of a
    text file, skips a binary file with reason `binary`, and uses the
    `archive-value` question of `questions.json`. That question is a choice
    of `reusable_knowledge`, `repository_record` and `transient_evidence`. It
    reuses the Spec judge's key variables, transport order, ceiling, retry,
    redaction and log. It fails open: a missing key, a reached ceiling or a
    service error yields `skipped` advice with the reason and a nil error.
15. `--plan` reads the Spec and writes nothing under the repository. Its only
    write is the judge log under the Roundfix Home. It does not require
    archive eligibility. It refuses, exit 2, for a Spec that is not active.

### Data Models

An Archive Record, as the archive writes it:

```markdown
---
schema: roundfix/archive-record/v1
spec: 0242-an-archive-that-leaves-an-archive-record
title: An archive that leaves an Archive Record
status: archived
created: 2026-10-06
archived: "2026-10-08"
disposition: pass
source: docs/specs/0242-an-archive-that-leaves-an-archive-record
source_revision: 0123456789abcdef0123456789abcdef01234567
qa_task: task_05
qa_report: qa-report-2026-10-08.md
qa_verdict: pass
unproven: []
adrs: [ADR-0247]
sources: [2026-10-06-history-keeps-only-what-the-secondbrain-needs.md]
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
---

# An archive that leaves an Archive Record

<one paragraph>
```

No Run Database, Task file, QA Report or PRD field changes. An archive no
longer stamps the PRD.

### API Contracts

1. API Contract: `roundfix archive <slug>` on an eligible, committed Spec
   writes the record, removes the folder and prints one confirmation line.
   Exit 0.
2. API Contract: `roundfix archive <slug>` exits 2 through Preflight
   Validation in two new cases: the Spec folder differs from `HEAD`, or a
   record or folder for the slug already exists. Every refusal of today
   stays.
3. API Contract: `roundfix archive <slug> --promote <path>`, repeatable, also
   copies each file to `docs/references/<basename>`. It exits 2 for a path
   Invariant 4 refuses. It combines with `--qa-override`.
4. API Contract: `roundfix archive <slug> --plan` prints the plan to stdout
   and exits 0. It refuses, exit 2, for a Spec that is not active, and as a
   usage error with `--promote`, `--qa-override`, `--approval` or
   `--reason`.
5. API Contract: the Delivery Queue treats an archive commit as exact under
   Invariant 8 and parks any other archive commit as today.

### Surface Transcripts

1. Surface Transcript: a normal archive.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   archived <slug> -> docs/history/specs/<slug>.md; removed <n> file(s) (<b> bytes) kept in Git at <12-hex>
   stderr:
   exit: 0
   ```

2. Surface Transcript: an override archive with a promotion.

   ```transcript
   $ roundfix archive <slug> --qa-override --approval <source> --reason <text> --promote references/<file>.md
   stdout:
   archived <slug> with QA override -> docs/history/specs/<slug>.md; removed <n> file(s) (<b> bytes) kept in Git at <12-hex>; promoted 1 file(s) to docs/references/
   stderr:
   exit: 0
   ```

3. Surface Transcript: the plan with a fake judge.

   ```transcript
   $ roundfix archive <slug> --plan
   stdout:
   archive plan for <slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md
   core <k> file(s): _prd.md, _techspec.md, _tasks.md, _authorization.md, <task files>, <QA Reports>, references/_index.md and adopted sources
   evidence <e> file(s) (<eb> bytes) under qa/evidence/
   candidate references/<file>.md <bytes> bytes: reusable_knowledge (p=0.94)
   candidate measurement/<file>.md <bytes> bytes: transient_evidence (p=0.81)
   promote with: roundfix archive <slug> --promote <path>
   stderr:
   exit: 0
   ```

4. Surface Transcript: the plan without a key.

   ```transcript
   $ roundfix archive <slug> --plan
   stdout:
   archive plan for <slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md
   ...
   candidate references/<file>.md <bytes> bytes: no advice (<KEY_VARIABLE> is not set)
   promote with: roundfix archive <slug> --promote <path>
   stderr:
   exit: 0
   ```

5. Surface Transcript: a folder that differs from `HEAD`.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   stderr:
   Preflight failed
   ...
   Reason:
     Spec "<slug>" has changes not committed at HEAD under docs/specs/<slug>; commit them before archive so the Archive Record's source_revision holds the Spec
   ...
   exit: 2
   ```

## Coverage Map

- Goal 1 → Archive Record and cut (Invariants 1-6; API Contracts 1, 2; Surface Transcripts 1, 2, 5)
- Goal 2 → `source_revision` (Invariants 5, 6, 9)
- Goal 3 → Archive Advice and promotion (Invariants 4, 14, 15; API Contracts 3, 4; Surface Transcripts 3, 4)
- Goal 4 → readers (Invariants 7-13; API Contract 5; Readers)
- Goal 5 → guidance (Exact texts; Vocabulary Contract)
- User Story 1 → Archive Record
- User Story 2 → Archive Advice, `--promote`
- User Story 3 → exact retirement
- User Story 4 → readers, pinned history
- Core Feature 1 → Invariants 2, 3
- Core Feature 2 → Invariants 4-6
- Core Feature 3 → Invariants 14, 15
- Core Feature 4 → Invariants 7-12
- Core Feature 5 → Invariant 13; Readers
- Core Feature 6 → Exact texts
- Success Metric 1 → task_01 tests
- Success Metric 2 → task_01 and task_02 tests
- Success Metric 3 → task_04 tests
- Success Metric 4 → task_03 tests; QA ablation row
- Success Metric 5 → Readers; QA reader audit

## Integration Points

- Git, through the local `git` binary the archive and the queue already use:
  `rev-parse`, `status --porcelain`, `ls-tree`, `cat-file`, `show`,
  `archive`, `log`.
- The Spec judge's existing transports (OpenRouter `systemone`, TypeSafe),
  for `--plan` only.

## Testing Approach

- `internal/spec/archive_record_test.go` archives temporary Specs: a
  passing Spec with adopted sources, ADRs and a regeneration; a partial with
  `unproven`; an override with a 300-byte reason; a superseded Spec. It
  covers the size bound with a long outcome and a long reason, the
  round-trip, every refusal of Invariant 4, a failure before removal, and
  `ReadArchivedSpec` over a record, a legacy folder, both and neither.
- `internal/cli/archive_record_test.go` runs the Archive Command in a
  temporary repository and asserts Surface Transcripts 1, 2 and 5. It also
  checks that the removed bytes equal `git show <source_revision>:...`.
- `internal/cli/deliver_archive_record_test.go` builds parent and archive
  commits and covers the exact retirement: exact, a kept file, an extra
  path, a wrong revision, a promoted blob that differs, and a legacy move.
  It also runs the queue's archive stage end to end with the real command.
- `internal/cli/archive_record_readers_test.go` and
  `internal/worktree/archive_record_test.go` cover the task_02 readers.
  Each reader gets a record fixture that reaches `main` through a squash
  merge, and a legacy-folder fixture. Existing reader tests stay unchanged
  and keep proving the legacy form.
- task_03 adds record tests to `speccheck`, `specaudit` and
  `suiteguardcontract`. It moves the 22 repository tests to
  `gittest.PinnedHistory` or to fixtures they own, and proves them on a copy
  without Spec folders.
- `internal/judge/archive_advice_test.go` and
  `internal/cli/archive_plan_test.go` use a fake transport and a temporary
  home, and assert Surface Transcripts 3 and 4.
- `internal/baseline/archive_record_clause_test.go` asserts both appended
  sentences in the embedded clauses, the formatter golden and this
  repository's generated guide.

## Build Order

1. The record, the cut, the Archive Command and the queue's archive stage
   and exact retirement, with supersede. They ship together because the
   queue's and supersede's existing tests run the real Archive Command.
2. The other runtime readers: queue inspection and prerequisites, review
   correction, reconcile, the pre-PR review, and Run causes (depends on: 1).
3. The Spec check, the Spec audit, the suite guard's regenerations, the
   pinned history, the 22 repository tests, the coverage record and the
   user-guide links (depends on: 1). This step shares no file with step 2.
4. Archive Advice, `--plan`, `--promote`, and the skills, guides, Baseline
   clauses with `make baseline-digests` and the Managed Refresh, and the
   glossary (depends on: 2, 3).
5. Final QA gate (depends on: 4).

## Risks & Considerations

- A conclusion that lived only in a removed file leaves the tree. The plan
  and the advice make the candidates visible, the qa-gate skill asks for an
  `## Outcome`, and the bytes stay at `source_revision`.
- After a squash merge, `source_revision` is reachable only through the
  delivery branch and the Pull Request's `pull/<n>/head` ref. A fresh clone
  of `main` does not fetch it. No reader on `main` needs it, and the readers
  that do run against the Run's or candidate's own branch (Invariant 9).
- The exact retirement proves agreement, not byte identity. The removed bytes
  are the parent's tree, which the proof names.
- Reconcile becomes conservative when `source_revision` is absent locally. A
  worktree is kept rather than released.
- The `### QA settlement` tables in the archive-spec, qa-gate and Roundfix
  skills say a passing archive keeps "The Spec and its QA report and
  evidence", which the cut makes false. The operator lifted the briefing's
  rule against editing that section for this Spec only (2026-10-06), so
  task_04 rewrites the three tables identically (Exact texts), and
  `TestSettlementGuidanceIsOneTable` keeps them byte-identical.
- Specs 0239, 0240 and 0241 land first and may touch `CONTEXT.md`, owned
  skills and Baseline-derived files. task_04 regenerates on top of theirs.
- Jev's monthly ceiling was breached on 2026-10-02. `--plan` reports the
  ceiling as its skip reason and never raises it.
- What Spec 0243 must do:
  - tag `history-full-2026-10` at the last full commit;
  - add `roundfix history sanitize`, a dry run unless `--apply`, working in
    batches through Pull Requests;
  - convert each legacy folder through `BuildArchiveRecord`, filling
    `pull_request` and `delivery_commit` from Git;
  - reduce reviews, handoffs and terminal backlog and findings to records or
    delete them;
  - promote the files the 2026-10-06 judgment found reusable;
  - repoint the `absorbed_by` and `CHANGELOG.md` paths that name folders;
  - restore `.secondbrain-export` to the whole of `docs/`;
  - then retire the legacy-folder branches of `ReadArchivedSpec`, the
    exact-move proof, ADR-0230's link pass and the pinned-history fallback
    where nothing reads them.

## Exact texts

Each Baseline sentence is appended after the clause's last sentence. Every
existing sentence stays byte-identical.

- `clause.spec.keep-artifacts-in-spec-folder`: "A runtime whose archive
  leaves an Archive Record, such as `roundfix archive <slug>`, writes
  `<archive-root>/<slug>.md` and removes the Spec folder in the same change;
  the folder stays in Git at the record's `source_revision`, so copy every
  file a later reader needs to `docs/references/`, an accepted ADR, the
  glossary, an agent guide or the Secondbrain inbox before or at archive."
- `clause.context.docs-one-job-per-directory`: "When the runtime's archive
  leaves an Archive Record, the record names each Spec-owned adopted
  reference, and the reference leaves the tree with the Spec folder."

The `### QA settlement` section changes identically in the archive-spec,
qa-gate and Roundfix skills. The outcome names and the Settles column stay
byte-identical, and so does every other row's Archives cell:

- The intro sentence ends "determines what the archive leaves:" instead of
  "determines what archive may move:".
- `pass`, Archives: "The Archive Record, which names the QA Report and
  verdict; the Spec, its QA report and evidence stay in Git at the record's
  `source_revision`."
- qualifying declared `partial`, Archives: "The Archive Record, which names
  the QA Report and verdict and carries the declarations' `satisfied-by`
  record as `unproven`; the Spec, its QA report and evidence stay in Git at
  the record's `source_revision`."
- `override`, Archives: "The Archive Record with `qa_override`,
  `qa_override_approval`, `qa_override_reason`, `qa_override_qa_outcome`,
  `qa_override_qa_task_status` when the QA Task is incomplete, and
  `qa_override_revision`; the QA Task and Reports stay unchanged in Git at
  the record's `source_revision`."

Each skill also adds a section outside `### QA settlement`:

- archive-spec: `## Archive Record`, before `## Unarchive`.
- qa-gate: `## Outcome for the Archive Record`, before `## Anti-patterns`.
  It asks for one `## Outcome` paragraph in the QA Report. In a Delivery
  Queue Run, it also asks the gate to copy any file `roundfix archive <slug>
  --plan` shows as reusable to `docs/references/` in the QA commit.
- Roundfix archive reference: `### Archive Record`, at the end of
  `## Archive Command`.

Each names `<slug>.md`, `source_revision`, `--plan` and `--promote`.

## Vocabulary Contract

- emits: `internal/cli/archive.go`
  pattern: `kept in Git at`
  documented-in: `docs/user-guide/commands/archive.md`
- emits: `internal/cli/archive.go`
  pattern: `archive plan for`
  documented-in: `docs/user-guide/commands/archive.md`

`CONTEXT.md` gains **Archive Record** and **Archive Advice**, and the
**Archive Command** entry stops saying the command "moves the whole Spec".

## Glossary

- adds: **Archive Record**
- adds: **Archive Advice**
- changes: **Archive Command**

## Research basis

- Secondbrain: `wiki/index.md` was read, then
  `qmd query "archive record spec history retention provenance" --all --files --min-score 0.3`.
  The query returned, among others,
  `projects/roundfix/mirror/docs/history/specs/0085-what-an-agent-reads-before-it-decides/references/2026-08-09-the-archive-belongs-outside-the-read-path.md`,
  the origin of the single History Root that the records keep. It also
  returned `wiki/sources/digest-roundfix-checkout-completo-para-testes-de-evidencia-historica-2026-09-09.md`,
  which records that CI needs `fetch-depth: 0` to read a historical commit.
  `ci-verify.yml` already sets that depth, so the pinned-history reads work in
  CI.
- Exa: GitHub's "Checking out pull requests locally"
  (<https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/checking-out-pull-requests-locally>)
  confirms that `git fetch origin pull/ID/head` reaches a pull request's
  commits after its branch is gone. A Stack Overflow answer and a TruffleHog
  write-up found with it agree, and add that a squash merge leaves those
  commits off the base branch. That is why a squash-merged Spec's
  `source_revision` stays reachable but is not on `main`, and why no reader on
  `main` depends on it.
- The parked draft Spec 0238 (branch
  `docs/archive-keeps-the-report-not-the-raw-evidence`) supplied the reader
  analysis for the coverage record and the owned-skill test. Its manifest and
  link rewrite are not built, by the maintainer's decision "Incorporar ao
  saneamento".

## Decisions

- Archive Record at `<archive-root>/<slug>.md`, removing the folder, with
  bytes in Git at `source_revision`. See ADR-0247.
- One resolver reads both the record and the legacy folder until Spec 0243
  migrates every folder.
- The exact retirement checks agreement with the parent, and the legacy
  exact move stays.
- Jev advises through the Spec judge's plumbing, and a person or Agent
  promotes. See ADR-0247.
- Repository tests read the pre-cut corpus at
  `40a7893d872c8a6705f6d7745e6efe430ae9deeb` instead of copying 50 MB into
  test data. See ADR-0247.
