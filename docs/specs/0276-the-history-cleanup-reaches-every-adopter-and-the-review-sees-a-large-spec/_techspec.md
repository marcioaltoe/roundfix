---
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
prd: _prd.md
created: 2026-10-17
---

# The history cleanup reaches every adopter, and the review sees a large Spec — Technical Spec

## Executive Summary

Four slices extend existing seams; no package is new. The history planner in
`internal/spec` stops refusing the `adr` unit for a successor that does not
name its predecessor: it appends the missing lineage line to the successor and
refuses only an ADR without an accepted successor, by its path. A new
`live-links` planner beside `record-links` gives Live Documents and user guide
pages the Git Locators `SC-HISTORY-RETAINED` already prints. The pre-PR review
in `internal/cli` omits deleted history and declared generated files, and
reviews a `codex` round-1 candidate still above the bound as up to four Review
Segments in one session; Delivery Convention C2 exempts the QA Task. A
`internal/speccheck` detector holds an owned skill edit to its version record.
`roundfix archive --plan --offline` skips the judge, and the Delivery Queue
probes and runs a renamed copy of each item binary build.

The trade-off accepted is reach over generality. Segments split by commit
range, which keeps each finding anchored in a tree the author committed but
cannot split one oversized commit; omission uses only exact and glob derived
declarations of the base commit, which leaves a directory entry's files in the
diff; and the lineage line is written only where `superseded_by` already names
an accepted successor, so an ADR without one still needs a person (ADR-0282).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema is
  minted; one finding code joins the code set and the Pre-PR Review Record
  gains optional `segments` and per-finding `segment` fields. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; the offline flag removes one. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0282 governs every decision below:
  "The cleanup writes the lineage a successor lacks.", "Live Documents get
  their Git Locators.", "A QA Task settles on its QA Report." and "Spec check
  holds an owned skill edit to its version.". ADR-0278: "A retirement
  deletes." and "A link to a deleted file becomes a Git Locator.". ADR-0254:
  "Refused Units never block.". ADR-0260: "The update lists.". ADR-0281: "A
  test that rebuilds an executable replaces it.". ADR-0225: "A queue item runs
  the Roundfix binary its branch builds". ADR-0196: "A pre-PR review finding
  parks only after validation". ADR-0197: "A pre-PR reviewer lineage spans at
  most two rounds". ADR-0227: "The Claude review prompt is bounded in tokens,
  not bytes". ADR-0273: "The judge's keys never enter a Run.". ADR-0189: "A
  version names one content.". ADR-0280: "A guide Task records its own
  fingerprint.". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276: "A row
  the authorization keeps offline qualifies, and a byte invariant compares
  with the Delivery Base". ADR-0169: "The pre-PR review diffs the candidate from its merge base" is kept, since every Review Segment is a range of the candidate measured from that merge base. Set no obligation here, because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0251, ADR-0261, ADR-0262 and ADR-0269 (history families, Archive Records, cancellation and the History Full Tag); ADR-0233 and ADR-0250 (version regeneration at merge and module versions); ADR-0258 (reduced history); ADR-0277 (legacy folder conversion); ADR-0165, ADR-0272 and ADR-0279 (review parks, records and dismissals); ADR-0167, ADR-0229, ADR-0237 and ADR-0240 (QA rows, operator archive and delivery retry). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `docs/references/coverage-record.json`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Characterized on 2026-10-17 with `bin/roundfix` 0.85.0 built at `e5da154a`.

- ADR unit. `PlanHistoryKind` calls `historyADRSuccessor` for every
  superseded ADR and returns its error, so one ADR refuses the whole `adr`
  unit: `superseded ADR <name>: accepted successor <id> does not name it`, or
  `superseded ADR <name> has no accepted successor naming it` without front
  matter. `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest` pins the first
  text with an accepted successor that does not name its predecessor. The
  adopter dry runs of 2026-10-17, measured by the lead in scratch clones and
  recorded in the first Backlog Entry: vortex refused by ADR 0002 (successor
  0051), 18 retained ADRs and 3 broken live links; fluxus refused by ADR 0002
  (successor 0064), 22 retained ADRs; oraculum refused by ADR 0012 (no
  accepted successor), 38 retained ADRs and 2 broken live links; conexus
  clean.
- Links. `PlanRecordLinks` reads only Archive Records under the archive root.
  `detectHistoryRetained` walks `docs/` and the Spec Root and reports a
  relative link to a missing file from an Archive Record, a Live Document
  (`spec.IsLiveDocument`) or a `docs/user-guide/` page, with the fix
  `Replace the link with <locator>.`; no unit writes that fix for a page that
  is not a record.
- Review diff. `reviewScopedDiff` omits `archived-spec-source`, `qa-evidence`
  and `upstream-skill` paths. Its summary line counts an archived Spec source
  as an upstream-managed skill. `checkReviewDiffBound` refuses a `codex`
  diff above 917504 bytes with `review diff too large: <bytes> bytes after
  omitting <n> path(s) exceeds the review bound of 917504 bytes`; Spec 0272's
  candidate measured 1180615 bytes after 174 omissions on 2026-10-17, and the
  operator reviewed it by hand in three segments. A `claude` prompt is bounded
  by its token budget instead.
- C2. `deliveryConventions()` in version `roundfix/delivery-conventions/v2`
  says the Daemon writes `## Result`, `## Recorded paths` and
  `## Carry-forward provenance` with no QA exemption; the review of Spec 0273
  on 2026-10-17 kept finding F1 on its QA Task from that text.
- Owned skills. `internal/speccheck` reads `skillcoverage.MapPath` and
  `skillcoverage.RecordPath`, never `skills/testdata/owned-skill-versions.json`;
  the record's `skills` keys are the thirteen owned skills.
- Archive plan. `runArchivePlan` always calls `judge.AdviseArchive`; without a
  key it prints
  `candidate <file> <bytes> bytes: no advice (ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set)`.
  `--offline` does not exist. The `qa-gate` skill's "Outcome for the Archive
  Record" asks for `roundfix archive <slug> --plan`; no Go source renders that
  instruction into a QA prompt.
- Item binary. `stepExecutable` builds, probes and returns
  `<item worktree>/<delivery.item_binary.path>` itself, so a later build that
  copies over that path rewrites a file an earlier step's process may still
  map. ADR-0281 fixed the same mechanism only in the delivery tests.

## System Architecture

- History (`internal/spec`, `internal/cli`). `PlanHistoryKind` for
  `ArchiveKindADR` collects a `HistoryLineage` for each superseded ADR whose
  accepted successor does not name it, and a `HistoryRefusal` for each without
  an accepted successor, which it leaves out of `Files`. `ApplyHistoryKind`
  appends each lineage line after removing the files. A new
  `internal/spec/live_links.go` holds `PlanLiveLinks`, sharing
  `MissingRelativeLinks`, `ResolveGitLocator` and `rewriteRecordLinks` with
  `record_links.go`. `historyInventory`, the selection loops of
  `internal/cli/history.go` and `internal/cli/baseline_update_history.go`, the
  rewrite plan and both printers gain the unit, the lineage citations and the
  per-file refusals.
- Review (`internal/cli`). `reviewScopedDiff` gains the two reasons and reads
  the base commit's derived declarations through
  `config.DeliveryConfigAtCommit`. `runConfiguredReviewSession` plans Review
  Segments when a `codex` round-1 diff exceeds the bound and runs each segment
  prompt with `RunPrepared` in the one prepared session; the command merges
  their answers. `deliveryConventions()` rewords C2 under
  `roundfix/delivery-conventions/v3`.
- Spec check (`internal/speccheck`). New `skill_version.go` holds
  `detectSkillVersion`, called from `Check` in `constraints.go` next to
  `detectSurfaceRecord` and registered at `StageTasks` in `coherence.go`.
- Archive and delivery (`internal/cli`). `parseArchiveCommand` accepts
  `--offline` with `--plan`; `runArchivePlan` skips `judge.Load` and
  `judge.AdviseArchive` under it. `stepExecutable` copies the built file
  before its probe.

## Implementation Design

### Interfaces

```go
// internal/spec
type HistoryLineage struct{ Successor, Removed, Line string } // Line has no newline
type HistoryRefusal struct{ Path, Reason string }
// HistoryKindPlan gains: Lineage []HistoryLineage; Refused []HistoryRefusal
func PlanLiveLinks(ctx context.Context, root, specRootRel, archiveRootRel string, skip map[string]bool) (RecordLinkPlan, error)

// internal/speccheck
const CodeSkillVersionUntasked = "SC-SKILL-VERSION-UNTASKED"
func detectSkillVersion(result *Result, repoRoot, specDir string, graph *spec.Graph) error

// internal/cli
type reviewSegment struct {
	Base      string `json:"base"`
	Head      string `json:"head"`
	DiffBytes int    `json:"diffBytes"`
}
// reviewRecord gains Segments []reviewSegment `json:"segments,omitempty"`;
// reviewFinding gains Segment int `json:"segment,omitempty"`.
const reviewSegmentLimit = 4
```

### Data Models

No new file format. The lineage line is Markdown in the successor ADR; the
review record's new fields are optional and absent from an unsegmented or
older record.

### API Contracts

1. API Contract: lineage completion. For each superseded ADR the `adr` unit
   would remove, from `docs/adr/` or a legacy ADR root, whose front matter
   `superseded_by` is `ADR-NNNN` naming a file `docs/adr/NNNN-*.md` with
   `status: accepted` that does not already name the superseded ADR (today's
   `historyEntryNamesPath` test), the plan keeps the file in `Files` and adds
   one `HistoryLineage` whose line is exactly
   `Supersedes ADR-<number> (<date>); Git keeps its text at <locator>.`, where
   `<number>` is the superseded file's four-digit prefix, `<date>` is the
   committer date of `HEAD` in UTC as `YYYY-MM-DD`, and `<locator>` is
   `ResolveGitLocator(<removed path>, "", "").Text(false)`. Apply appends the
   line to the successor after one blank line, adding a final newline first
   when the file lacks one, and never twice; several lines for one successor
   follow in path order. The plan lists each line as a History Citation
   `{Path: successor, Line: <line it will occupy>, Target: <removed path>,
   Rewrite: <line>}`, so text prints `cites <successor>:<n> names <removed>;
   rewrites to <line>`, JSON lists it under `citations`, and the
   `baseline update` Plan Digest binds it. The successor joins the unit's
   uncommitted-change check but not its `paths`.
2. API Contract: refusal by path. A superseded ADR with no `superseded_by`, or
   one naming no `status: accepted` file, stays in the tree and becomes a
   `HistoryRefusal` with reason `superseded ADR has no accepted successor;
   name the accepted ADR that replaces it in superseded_by, or set its status
   to deprecated`. Both commands list it as a Refused Unit whose unit is the
   ADR's repository-relative path (`refused <path>: <reason>` in history text,
   `- refused <path>: <reason>` in baseline text, `{"unit": "<path>",
   "reason": "<reason>"}` in JSON); the `adr` unit still removes every other
   retired ADR, and is omitted only when nothing else remains.
3. API Contract: `live-links` unit. After `record-links`, both commands plan
   `live-links` over every regular Markdown file a walk of `docs/` and the
   Spec Root finds without following links that is not an Archive Record, is
   not removed by the selected units, and is a Live Document or lies under
   `docs/user-guide/`. Each `spec.MissingRelativeLinks` result that is not a
   self-link becomes `ResolveGitLocator(<target>, "", "")`, rewritten through
   the link-span rewrite `record-links` uses, so code spans, fences, URLs and
   existing Git Locators are untouched. Text: `live-links: rewrites <l> link(s)
   in <n> file(s)` (history) and `- live-links: rewrites <l> link(s) in <n>
   file(s)` followed by its paths (baseline); JSON unit `live-links`, action
   `rewrite`, `paths` the rewritten pages, `files`, `links`; its citations join
   the list and the digest as `record-links` ones do. A second plan after apply
   plans no `live-links` unit.
4. API Contract: review omissions. Besides today's reasons, `reviewScopedDiff`
   omits a path deleted between base and head under `docs/history/` as
   `retired-history`, and a path matched by an entry of a
   `delivery.derived_paths[].paths` list that does not end in `/`, through
   `path.Match`, as `generated`. Declarations are read with
   `config.DeliveryConfigAtCommit` at the base commit, never the head; a base
   without a Project Config omits nothing as generated, and an unreadable one
   blocks with `review scope: read derived paths: <error>`. Precedence:
   `archived-spec-source`, `qa-evidence`, `upstream-skill`, `retired-history`,
   `generated`. The prompt's summary becomes `Omitted from this diff (not
   reviewed): <q> path(s) of QA evidence, <u> of upstream-managed skills, <a>
   of archived Spec sources, <r> of retired history, <g> of generated files.`
   The rules apply to round 2's delta too.
5. API Contract: Review Segments. When the provider is `codex`, the round is 1
   and the scoped diff exceeds `reviewDiffBound`, the command lists the
   candidate's commits with `git rev-list --first-parent --reverse
   <base>..<head>` and packs consecutive commits greedily: a segment grows
   while `reviewScopedDiff(<segment base>, <next commit>)` stays within the
   bound. A single commit whose own scoped diff exceeds the bound blocks with
   today's reason. More than `reviewSegmentLimit` segments blocks with
   `review diff too large: <bytes> bytes after omitting <n> path(s) needs <k>
   review segments; a review allows at most 4`, before any provider call.
   Otherwise one session is prepared as today and each segment prompt,
   `buildReviewPrompt(<segment base>, <segment head>, <segment diff>,
   <segment omissions>)` preceded by the line `Review Segment <i> of <k>:
   commits <segment base>..<segment head> of the candidate <base>..<head>;
   answer for this segment only.`, is sent in order with `RunPrepared`; the
   Spec context is appended to the first prompt only. Each answer is classified
   as today; the first answer that blocks or fails decides the record and no
   later segment is sent. Otherwise findings from every segment are numbered
   `F1`… in segment order, each with its `segment`, validated against its own
   segment's diff, and the record carries `segments` and `diffBytes` equal to
   their sum. The persisted answer joins the answers under lines `--- Review
   Segment <i> of <k> ---`. Verdict reuse, dispositions, the convention
   validator and round accounting stay per candidate head: all segments are
   one round. Round 2 and `claude` are unchanged.
6. API Contract: C2. `deliveryConventionsVersion` becomes
   `roundfix/delivery-conventions/v3`, and the C2 prompt text is exactly the
   line in the block below, where the section names and `completed` are code
   spans as today. The other conventions keep their bytes.

   ```text
   The Daemon writes a Task file's status after the Task's Verification passes, its `## Result` and `## Carry-forward provenance` sections for non-QA Tasks only, and `## Recorded paths` only when a Task changed a path it did not declare; a QA Task settles on its committed QA Report and carries none of these sections. A Result that calls status Daemon-owned agrees with a `completed` status.
   ```
7. API Contract: `SC-SKILL-VERSION-UNTASKED`. When
   `skills/testdata/owned-skill-versions.json` is a regular file, its `skills`
   keys are the owned skills; for each Task of the graph that is not
   completed, not the manifest's `qa:` Task and not of type `qa`, and each
   owned skill `<name>` for which the Task declares a path
   `.agents/skills/<name>/<rel>` under `interface:`, `creates:` or `deletes:`,
   the detector reports one gap unless the Task declares the mirror
   `skills/<name>/<rel>` of each such path under the same kind (`deletes:` for
   a deletion, `interface:` or `creates:` otherwise),
   `.agents/skills/<name>/SKILL.md`, `skills/<name>/SKILL.md` and
   `skills/testdata/owned-skill-versions.json` under `interface:` or
   `creates:`, and names `-record-skill-versions` in a requirement. Summary:
   `Task <id> edits owned skill <name> without declaring its mirror, both
   SKILL.md files and skills/testdata/owned-skill-versions.json and naming its
   -record-skill-versions step`, located at the Context line of the skill's
   first declared path. Fix: `In this Task, add <missing>, following the
   write-tasks rule "Owned skill versions".`, where `<missing>` lists, joined
   by `; `, `interface: <path>` for each missing declaration and `a requirement
   naming -record-skill-versions` when the step is missing. A repository
   without the record skips the code with that path, a Spec without a graph
   with its `_tasks.md` path; an unreadable record is an error. A gap with no
   horizon, from the Tasks stage.
8. API Contract: offline archive plan. `roundfix archive <slug> --plan
   --offline` prints the plan of today with `judge: skipped (offline)` after
   the `evidence` line and every candidate, in plan file order, as
   `candidate <file> <bytes> bytes: no advice (offline)`; it calls neither
   `judge.Load` nor `judge.AdviseArchive` and reads no Judge Key Variable.
   `--offline` without `--plan` is a validation error `--offline requires
   --plan`. The usage line becomes `roundfix archive <slug> --plan [--offline]`
   with the option `--offline      With --plan, skip Archive Advice and contact
   no judge`. `--plan` alone is unchanged.
9. API Contract: item binary copy. After a successful build, `stepExecutable`
   computes the SHA-256 of the regular file at the declared path and writes a
   copy to `<artifact dir>/delivery/<slug>/item-binary/<first 16 hex>/<base
   name of path>`: it creates the directory, writes the bytes to a temporary
   file created in that directory, sets mode `0755`, closes it and renames it
   over the final name. The probe and the step run that copy, and the
   selection line names it. A declared path that is not a regular file, or any
   copy failure, returns `copy item binary: <error>` and parks the item
   `delivery-error` as a failed build does. The fallback to the owner's binary
   is unchanged.
10. API Contract: everything else stays. `SC-HISTORY-RETAINED`, the review
    bound value, the Claude budget, the other conventions, every other judge
    call, existing finding counts and the build command keep their behavior.

### Surface Transcripts

1. Surface Transcript: the offline archive plan of a Spec whose only
   candidate is one file under `references/` and which no Live Document cites.

   ```transcript
   $ roundfix archive <slug> --plan --offline
   stdout:
   archive plan for <slug>: removes <count> file(s) (<bytes> bytes) and writes docs/history/specs/<slug>.md
   core <count> file(s): <files>
   evidence 0 file(s) (0 bytes) under qa/evidence/
   judge: skipped (offline)
   candidate references/<file> <bytes> bytes: no advice (offline)
   promote with: roundfix archive <slug> --promote <path>
   stderr:
   exit: 0
   ```

The history, baseline and review output lines are fixed by API Contracts 1 to
5 and asserted by the named tests; their full transcripts depend on fixture
paths and commits.

### Invariants

```text
1. A plan never removes a superseded ADR whose successor will not name it after apply.
2. Every other retired ADR leaves even when one ADR is refused.
3. After apply, SC-HISTORY-RETAINED reports no missing link in a page live-links planned.
4. A review omission never reads the candidate's own Project Config.
5. Every candidate commit is in exactly one Review Segment, and each segment's diff is within the bound.
6. A segmented review is one round with one record.
7. --offline sends no judge request whatever the environment holds.
8. No item binary build writes into a file a probe or step executed.
9. Every earlier finding count in the corpus golden stays.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → `PlanHistoryKind` lineage and refusals, `ApplyHistoryKind` (API Contracts 1, 2)
- Goal 1, User Story 2, Core Feature 2, Success Metric 1 → `PlanLiveLinks` and the `live-links` unit (API Contract 3)
- Goal 2, User Story 3, Core Features 3 and 4, Success Metric 2 → `reviewScopedDiff`, Review Segments in `runConfiguredReviewSession` (API Contracts 4, 5)
- Goal 3, User Story 4, Core Feature 5, Success Metric 3 → `deliveryConventions()` C2 (API Contract 6)
- Goal 4, User Story 5, Core Feature 6, Success Metric 4 → `detectSkillVersion` (API Contract 7)
- Goal 5, User Story 6, Core Feature 7, Success Metric 5 → `runArchivePlan` offline, the `qa-gate` skill (API Contract 8, Surface Transcript 1)
- Goal 5, User Story 7, Core Feature 8, Success Metric 6 → `stepExecutable` copy (API Contract 9)
- Core Feature 9 → the six user guides, the Roundfix references, `write-tasks`, `qa-gate`, `CONTEXT.md`
- Success Metric 4 (unchanged corpus counts), API Contract 10, Invariant 9 → corpus pins and the scope audit

## Absence Evidence

1. Lineage line and per-file refusal — command: `grep -rn 'Supersedes ADR\|lineage' internal/spec --include='*.go'`; result: no production match at `e5da154a`; `historyADRSuccessor` only returns an error.
2. `PlanLiveLinks` and the `live-links` unit — command: `grep -rn 'live-links\|PlanLiveLinks' internal --include='*.go'`; result: no match at `e5da154a`.
3. Review omissions — command: `grep -rn 'retired-history\|"generated"\|DerivedPaths' internal/cli/review.go internal/cli/review_conventions.go`; result: no match at `e5da154a`.
4. Review Segments — command: `grep -rni 'segment' internal/cli/review.go internal/cli/review_validation.go`; result: no match at `e5da154a`.
5. `detectSkillVersion` — command: `grep -rln 'owned-skill-versions\|record-skill-versions' internal/speccheck`; result: no match at `e5da154a`.
6. `--offline` — command: `grep -n 'offline' internal/cli/archive.go`; result: no match at `e5da154a`.
7. Item binary copy — command: `grep -n 'item-binary/\|Rename' internal/cli/deliver_workflow.go`; result: no match at `e5da154a`; the build log is the only file written under the item's artifact directory.

## ADR Alternatives

- ADR-0282: this Spec's decision; it rejects refusing the whole unit, holding back the guide, raising the bound, splitting by file or bytes, reading the candidate's declarations, skipping the judge only without a key, and removing the built file first.
- ADR-0278: no rejected alternative is taken; superseded ADRs still leave the tree, and the lineage the decision keeps through the successor's line is now written when missing.
- ADR-0254: no rejected alternative is taken; one plan and one digest stay, and a kept ADR is a Refused Unit that never blocks.
- ADR-0260: no rejected alternative is taken; the update still lists each citation it writes.
- ADR-0281: no rejected alternative is taken; the copy is renamed into place, never written over a used file, and nothing waits or retries.
- ADR-0225: no rejected alternative is taken; the item still runs its own build, and the owner never re-executes as the item binary.
- ADR-0196: no rejected alternative is taken; conventions stay a closed numbered list and validation is unchanged.
- ADR-0197: no rejected alternative is taken; there is still no third round, and segments live inside round 1.
- ADR-0227: no rejected alternative is taken; the Claude prompt keeps its token budget and is never segmented.
- ADR-0273: no rejected alternative is taken; keys stay out of Runs, and the offline plan adds an explicit refusal to call the judge.
- ADR-0189: no rejected alternative is taken; the version record stays the only record, and the check reads it.
- ADR-0280: no rejected alternative is taken; the Task that edits the skill records its own version, as a guide Task records its fingerprint.
- ADR-0256: no rejected alternative is taken; each Task re-records the Behavior Surface Record it changes.
- ADR-0104: no rejected alternative is taken; the outside-evidence row quotes published sources.
- ADR-0276: no rejected alternative is taken; the network-free rows qualify and the corpus row compares with the Delivery Base.
- ADR-0092: no rejected alternative is taken; the backlog stays typed intent and retirement still deletes.
- ADR-0120: no rejected alternative is taken; the History Root keeps its one job.
- ADR-0247: no rejected alternative is taken; an Archive Record and its folder in Git are unchanged.
- ADR-0248: no rejected alternative is taken; the History Full Tag still precedes every removal.
- ADR-0261: no rejected alternative is taken; an Archive Record still names its pull request.
- ADR-0233: no rejected alternative is taken; a version raise is still recorded by its command.
- ADR-0250: no rejected alternative is taken; module versions are untouched.
- ADR-0258: no rejected alternative is taken; no Reduced History Entry is written.
- ADR-0272: no rejected alternative is taken; a review record still survives a fix at its head.
- ADR-0277: no rejected alternative is taken; legacy folder conversion is unchanged.
- ADR-0169: no rejected alternative is taken; each segment is diffed inside the candidate range from its merge base.
- ADR-0251: no rejected alternative is taken; lenient legacy reading is unchanged.
- ADR-0262: no rejected alternative is taken; cancellation is unchanged.
- ADR-0269: no rejected alternative is taken; archive and cancel refusals are unchanged.
- ADR-0165: no rejected alternative is taken; a blocking review still parks publication.
- ADR-0279: no rejected alternative is taken; a dismissal still survives a fix that leaves its anchor.
- ADR-0167: no rejected alternative is taken; the Pull Request row stays blocked before a Pull Request exists.
- ADR-0229: no rejected alternative is taken; the operator archive is unchanged.
- ADR-0237: no rejected alternative is taken; the Delivery Retry is unchanged.
- ADR-0240: no rejected alternative is taken; the QA partial policy is unchanged.

## Integration Points

None outside the repository. The lineage date and locator read `git` in the
repository being planned; the review reads the base commit's Project Config
through `git show`; the item binary copy writes only under the configured
artifact directory.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and builds temporary
repositories with the fixtures already used beside it.

- `internal/spec/history_lineage_test.go`: a successor that does not name its
  predecessor yields one lineage line with the exact bytes and date, applied
  once; a missing or non-accepted successor yields a refusal by path while a
  rejected ADR in the same unit is removed.
- `internal/spec/live_links_test.go`: a reference page and a user guide page
  with missing links are rewritten to Git Locators; Archive Records, code
  spans, URLs, existing locators, self-links and files the plan removes are
  untouched.
- `internal/cli/history_adopter_test.go`: the adopter shapes through `history
  sanitize` plan and `--apply`, then `speccheck.Check` reports
  `SC-HISTORY-RETAINED` only for the kept ADR; `baseline update --format json`
  lists the `live-links` unit, the lineage citation and the refusal, and its
  Plan Digest changes when the lineage line would change.
- `internal/cli/history_refusal_test.go`:
  `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest` keeps its name; its
  fixture's successor becomes `status: proposed`, so it asserts the refusal by
  path (declared break of its expected text).
- `internal/cli/review_scope_test.go`: deleted history and a generated golden
  are omitted with their reasons; a declaration the candidate adds and a
  directory entry omit nothing; the summary line of
  `TestReviewOmitsQAEvidenceAndUpstreamSkills` moves to API Contract 4's form
  (declared break).
- `internal/cli/review_segment_test.go`: a three-commit candidate above the
  bound sends three prompts in one prepared session and writes one record
  with `segments` and per-finding `segment`; five segments block before a
  provider call; `TestReviewBlocksAboveTheBoundWithoutAProviderCall` keeps
  passing for one commit and for round 2.
- `internal/cli/review_convention_qa_task_test.go`: the prompt carries v3 and
  the new C2 text, and a candidate holding a completed QA Task file shaped like
  Spec 0274's `task_05.md` gives a prompt whose C2 exempts it; no provider is
  called.
- `internal/speccheck/skill_version_test.go`: the gap, its acceptance, its
  skips and its stage, as `surface_record_test.go` does.
- `internal/cli/archive_plan_offline_test.go`: Surface Transcript 1 as one
  test; a judge transport that counts requests sees none with a judge key in
  the environment; `--offline` without `--plan` is refused.
  `TestArchiveHelpNamesThePullRequestOption` takes the new usage (declared
  break).
- `internal/cli/deliver_item_binary_copy_test.go`: the step runs the copy; a
  rebuild that copies in place over the build output while a held copy runs
  leaves that copy's bytes and file identity unchanged, executes a path that
  differs from the build output, and leaves no temporary file. The helpers of
  `deliver_item_binary_test.go` that expect the build output path move to the
  copy path (declared break).
- No existing test is renamed or removed, so the coverage record stays.

## Build Order

1. History lineage, per-file refusals and `live-links`, with their tests, the
   `history` and `baseline` guides, the Roundfix `archive`, `baseline` and
   `setup` references, versions and the Behavior Surface Record (no
   dependency).
2. Review omissions, Review Segments and C2, with their tests, the `review`
   guide and reference, the glossary entry, versions and the record (depends
   on: 1, which raises the same Roundfix version and record).
3. `SC-SKILL-VERSION-UNTASKED`, its tests, corpus pins, the `spec` guide and
   reference, the `write-tasks` rule, versions and the record (depends on: 2,
   same version and record).
4. The offline archive plan, the item binary copy, their tests, the `archive`
   and `deliver` guides and references, the `qa-gate` skill, versions and the
   record (depends on: 3, same version and record, and the 0275 prerequisite).
5. QA gate (depends on: 4).

Spec 0275 is a prerequisite: Step 4 changes `deliver_item_binary_test.go`,
which 0275 rewrites, and extends ADR-0281, so the manifest's `requires` names
it. Rebase on main after it merges and run the strict check again.

## Risks & Considerations

- Every Task raises the Roundfix skill version and re-records two shared
  files, so the graph is a chain; a failure in one Task holds the rest.
- A Live Document in an active Spec of another branch may gain a locator
  rewrite in the same update; the plan lists every page, and the operator
  reviews it.
- A candidate whose single commit exceeds the bound still needs a person; the
  reason says so and the guide names the remedy, splitting the commit.
- An item binary copy per distinct build accumulates under the artifact
  directory; ADR-0282 records that nothing reclaims it yet.
- The adopter facts are the lead's measurements; the QA gate rebuilds each
  shape in a temporary repository instead of reading an adopter.

## Glossary

- adds: **Review Segment**

## Decisions

- The lineage line follows the house `Supersedes` shape and HEAD's committer date. See ADR-0282.
- A kept ADR is a Refused Unit named by its path. See ADR-0282.
- `live-links` reuses the link scan and fix of `SC-HISTORY-RETAINED`. See ADR-0282.
- Derived declarations come from the base commit, without directory entries. See ADR-0282.
- Segments are first-parent commit ranges, at most four, for `codex` round 1. See ADR-0282.
- The archived Spec source omission gets its own count in the summary line.
- The item binary copy is content-addressed under the item's artifact directory. See ADR-0282.
