---
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
prd: _prd.md
created: 2026-10-01
---

# Gates and Run storage that let a correct delivery finish — Technical Spec

## Executive Summary

Four local causes stopped four correct deliveries. The pre-PR review sends
the whole `git diff` of the candidate as its prompt, with no filter and no
size limit, and Spec 0200's diff was 61% QA evidence. The failed-pass import
copies every added or modified file under the Spec's `qa/` directory,
whatever its extension. The Task Context parser knows three kinds, and the
unresolved-path detector requires every `interface` path to exist. And the
merged-head proof of reconcile compares every non-Task commit's paths with
the merged head, where the archive has moved and restamped the Spec's
directory, while a dirty Run Worktree is classified before any proof runs.

The fix narrows each one. The review omits QA evidence and upstream-managed
skill copies, records them, and blocks with a named cause above a measured
bound. The import skips compiled source and records the skip. A fourth
Context kind, `deletes`, is declared, audited and checked after completion.
Reconcile treats the archived Spec directory as representing the Run's
Spec-directory work, and declared or Spec-directory leftovers as superseded
(ADR-0212). The trade-off is that a review no longer reads QA evidence or
vendored skills, and that `--apply` may discard uncommitted changes it has
proved belong to a merged Spec's own scope.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; the
  review's provider call is unchanged and no test reaches a provider. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0212 (this Spec) governs the
  reconcile change and refines ADR-0053 and ADR-0161; the PRD accounts for
  ADR-0052, ADR-0080, ADR-0088, ADR-0091, ADR-0093, ADR-0096, ADR-0097,
  ADR-0104, ADR-0115, ADR-0117, ADR-0153, ADR-0155, ADR-0156, ADR-0166,
  ADR-0167, ADR-0168, ADR-0169, ADR-0170, ADR-0174, ADR-0176, ADR-0178,
  ADR-0183, ADR-0184, ADR-0187, ADR-0189, ADR-0194, ADR-0195, ADR-0196,
  ADR-0197 and ADR-0210, ADR-0165 and ADR-0192 hold unchanged, and ADR-0182
  and ADR-0211 do not apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Tasks edit the Roundfix, qa-gate and
  write-tasks skills, whose canonical files and `SKILL.md` mirrors are
  Governed Paths; express maintainer authorization: "considere autorizado a
  ajustar todas as skills se necessário"; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `skills/roundfix/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/write-tasks/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/_authorization.md`.

## System Architecture

No new package or command.

| Component | Where | Change |
| --- | --- | --- |
| Review diff | `reviewCandidateDiff` and the round-two delta in `internal/cli/review.go` | Omits evidence and upstream-managed skill paths; measures bytes |
| Review bound | `internal/cli/review.go` | Blocks above the bound before the provider call |
| Review record | `reviewRecord` in `internal/cli/review.go` | Gains `diffBytes`, `omittedPaths`, `runtimeStderrTail` |
| QA import | `findRecordedPriorQAPass` and `copyPriorQAPass` in `internal/daemon/qa_prior_pass.go` | Skips compiled source and records it |
| Context kinds | `internal/spec/spec.go`, `internal/spec/task.go` | Adds `deletes` |
| Context consumers | `internal/spec/recorded_paths.go`, `internal/spec/collision.go`, `internal/speccheck/citations.go`, `internal/speccheck/undeclared.go`, `internal/speccheck/surface.go`, `internal/daemon/task_context.go` | Treat `deletes` as declared |
| Deletion settlement | `verifyTask` in `internal/daemon/task_engine.go`, new `internal/daemon/task_deletes.go` | Refuses to settle while a `deletes` path exists |
| Merged-head proof | `inspectRunAtMergedHead` in `internal/worktree/merged_head.go` | Spec-directory paths represented by the archived Spec |
| Dirty classification | `inspectTerminalRunMerged` and `cleanupTerminalRun` in `internal/worktree/worktree.go` | Supersedes declared leftovers; forced removal only on that proof |
| Absent target | `classifyRunBranchSet` in `internal/worktree/worktree.go` | Falls back to the merged-head proof |
| Skills and guides | Roundfix `review` and `reconcile` references, qa-gate, write-tasks; `review` and `reconcile` command guides | Describe the behavior |

## Implementation Design

### Interfaces

```go
// internal/cli/review.go
const reviewDiffBound = 917504 // bytes; below 919,745, the largest diff reviewed here

type reviewOmittedPath struct {
	Path   string `json:"path"`
	Reason string `json:"reason"` // "qa-evidence" or "upstream-skill"
}

type reviewRecord struct {
	// existing fields unchanged
	DiffBytes         int                 `json:"diffBytes"`
	OmittedPaths      []reviewOmittedPath `json:"omittedPaths"`
	RuntimeStderrTail string              `json:"runtimeStderrTail,omitempty"`
}

// internal/spec/spec.go
const ContextKindDeletes ContextKind = "deletes"

// internal/daemon/task_deletes.go
// undeletedTaskPaths returns the Task's deletes paths that still exist in
// workDir, sorted.
func undeletedTaskPaths(workDir string, task spec.Task) ([]string, error)
```

### The reviewed diff

Every diff the review sends, round one's candidate diff and round two's
delta, is computed with the same omission. The omitted set is:

1. **QA evidence**: each changed path under `<root>/<slug>/qa/evidence/`, for
   the resolved Specs Root and its archive root, `docs/specs` and
   `docs/history/specs` for the built-in root. Reason `qa-evidence`. The QA
   Report itself, `<root>/<slug>/qa/qa-report-*.md`, stays in the diff.
2. **Upstream-managed skill copies**: each changed path under
   `.agents/skills/<name>/` for every `<name>` in the `skills` map of
   `skills-lock.json` at the base commit or at the head commit. Reason
   `upstream-skill`. A repository without the file omits nothing for this
   class; an unreadable or malformed lock blocks the review with
   `review scope: read skills-lock.json: <error>`.

The review lists the changed paths with `git diff --name-only` over the same
range, classifies them in that order, and computes the diff text with
`git diff` and one `:(exclude)<path>` pathspec per omitted path, keeping
today's flags. `DiffBytes` is the length of the diff text sent.
`OmittedPaths` is sorted by path, always present and empty when nothing is
omitted. The prompt gains, after the candidate diff block, one line
`Omitted from this diff (not reviewed): <n> path(s) of QA evidence, <m> of
upstream-managed skills.` followed by the omitted paths, one per line, so the
reviewer does not raise their absence. A finding anchored in an omitted path
is out of the diff and validated as today (ADR-0196).

### The bound

When `DiffBytes` exceeds `reviewDiffBound`, the review records `blocked`
before any Agent Session is prepared, with the reason
`review diff too large: <bytes> bytes after omitting <n> path(s) exceeds the review bound of 917504 bytes`,
and exits as a blocked review does today. No provider is called, and the
configured fallback is not tried. The constant cites its measurement in a
comment: Spec 0194's candidate of 919,745 bytes was reviewed and Spec 0200's
of 1,295,055 bytes failed.

### The stderr tail

When the review fails with a runtime error that carries the runtime's
standard error, `agent.BatchFailureError`'s `Stderr` reached through
`errors.As`, the record's `RuntimeStderrTail` holds its last 10 lines,
capped at the last 1,024 bytes. The reason text is unchanged.

### The QA import

`copyPriorQAPass` skips each file whose base name ends in `.go`, `.rs`,
`.ts`, `.tsx`, `.mts`, `.cts`, `.js`, `.jsx`, `.mjs` or `.cjs`, the sources
the toolchains of the Baseline's built-in profiles compile or format. A
skipped file is neither written nor compared, so it never causes
`path differs`. The `prior_report` journal event gains `skipped`, a list of
`{path, reason}` with reason `compiled source`, and `files` lists only what
was imported. A pass whose only `qa-report-*.md` would be skipped cannot
occur, since reports are Markdown. A carried row whose cited evidence was
skipped is refused by the carry proof's existing evidence check and re-run.

### The `deletes` Context kind

The parser accepts `- deletes: <path>` with the same path rules as the other
kinds. A path may appear once per Task across `interface`, `creates` and
`deletes`; a repeated `deletes` path or a path declared under two of those
kinds is a parse error naming both. Consumers:

- `detectTaskContextReferences` skips `deletes` as it skips `creates`.
- The Daemon refuses to settle a Task whose `deletes` path still exists:
  `verifyTask`, before it runs the Task's Verification commands, fails the
  attempt when `undeletedTaskPaths` is non-empty, with the diagnostic
  `deletes: <path> still exists` per path. The failure takes the existing
  failed-Verification route, so the bounded retry returns it to the same
  Agent Session. The Spec Consistency Check gains no finding code, because
  its stage table is a Governed Path this Spec does not touch.
- `UndeclaredTaskPaths` counts a `deletes` path as declared.
- `declaredGovernedTouches` counts it, so deleting a Governed Path still needs
  its grant (ADR-0178).
- The Wave collision rule counts it, so two Tasks touching the same path by
  edit and deletion run in series.
- `firstCLISurface` counts it; `taskNamesGuide` does not.
- The Daemon's Task context lists it with the interface paths, so the Agent
  is told which path to remove.

### The merged Spec's Runs

`inspectRunAtMergedHead` decides whether the merged head holds the Spec
archived: `<archive root>/<slug>/_prd.md` exists at the merged head. When it
does:

1. Paths of non-Task, non-QA commits under `<root>/<slug>/` or
   `<archive root>/<slug>/` are represented and are not compared. Every other
   path is compared as today.
2. The superseded reason counts them:
   `Run work is superseded at <label>: <t> Task commit(s) completed, <q> QA Report commit(s) superseded, <s> Spec-directory path(s) archived`.

`inspectTerminalRunMerged` no longer stops at a dirty worktree. It records
the dirty paths from `git status --porcelain=v1 -z --untracked-files=all`,
runs the merged-head proof, and classifies:

- `superseded` when the proof yields `safe` or `superseded`, the merged head
  holds the Spec archived, and each dirty path lies under the Spec's active
  or archive directory or is declared (`interface`, `creates`, `deletes`) or
  recorded (`## Recorded paths`) by a Task file of the archived Spec at the
  merged head. The reason appends `; <d> uncommitted path(s) superseded by the archived Spec`.
- `dirty` otherwise, with today's reason.

The evidence records the dirty path set. `cleanupTerminalRun` removes such a
worktree with `git worktree remove --force` only when the fresh revalidation
yields the same state, Run head and dirty path set; a worktree without dirty
paths is removed without `--force` as today.

`classifyRunBranchSet`, for an absent target branch and a terminal Run whose
superseding QA Report proof fails, applies `chooseMergedHead` and
`inspectRunAtMergedHead` before preserving, and releases the branch when that
proof yields `safe` or `superseded`. It also considers a Run whose working
root is a linked worktree of the repository, through
`roundconfig.RepositoryRoot`, as the Run Database's other listings do.

### Data Models

No Run Database change. The review record gains three fields; an older record
without them still reads, with `DiffBytes` zero and `OmittedPaths` empty. The
`prior_report` journal payload gains `skipped`.

### API Contracts

1. API Contract: the review record — `diffBytes` (integer, the bytes of diff
   sent, or measured when blocked by the bound), `omittedPaths` (array of
   `{path, reason}` with reason `qa-evidence` or `upstream-skill`, sorted by
   path, always present) and `runtimeStderrTail` (string, omitted when
   empty). Every existing field keeps its meaning.
2. API Contract: the bound refusal — outcome `blocked`, reason
   `review diff too large: <bytes> bytes after omitting <n> path(s) exceeds the review bound of 917504 bytes`,
   standard error `roundfix: review blocked: <reason>`, today's blocked exit
   code, no provider call.
3. API Contract: the `prior_report` journal event — `skipped` lists each
   compiled-source path with reason `compiled source`.
4. API Contract: Task Context — `- deletes: <path>` is a valid entry that
   the unresolved-path check skips and the scope audit, the Governed Path
   declaration audit and the Wave collision rule count; a Task whose
   `deletes` path exists at settlement fails its attempt with
   `deletes: <path> still exists`.
5. API Contract: `roundfix reconcile` — a released merged-Spec Run is
   reported `superseded` with the reason of "The merged Spec's Runs", and
   `--apply` removes it.

### Surface Transcripts

None. No command gains a flag, and the changed outputs are the review's
one-line JSON record, whose bytes carry commit identifiers and the checkout
path of the machine, a new settlement diagnostic in the Daemon's existing
Verification failure, and new reason text in reconcile's existing report. API Contracts 1,
2, 4 and 5 fix the new bytes, the Tasks' tests assert them field by field,
and the QA gate observes them through the built binary.

## Vocabulary Contract

No new glossary term. The emitted words are `qa-evidence`, `upstream-skill`,
`review diff too large`, `compiled source`, `deletes`,
`deletes: <path> still exists`, `Spec-directory path(s) archived` and `uncommitted path(s) superseded by the
archived Spec`. task_01 documents the review words in the `review` command
guide and the Roundfix Skill's `review` reference, task_02 the import rule in
the qa-gate skill, task_03 the `deletes` kind and its settlement
check in the write-tasks skill and template, and task_04 the reconcile words in the
`reconcile` command guide and reference.

## Coverage Map

- Goal 1 → The reviewed diff; The bound; API Contract 1.
- Goal 2 → The QA import; API Contract 3.
- Goal 3 → The `deletes` Context kind; API Contract 4.
- Goal 4 → The merged Spec's Runs; API Contract 5.
- User Story 1 → The reviewed diff.
- User Story 2 → The bound; The stderr tail.
- User Story 3 → The QA import.
- User Story 4 → The `deletes` Context kind.
- User Story 5 → The merged Spec's Runs.
- Core Feature 1 → The reviewed diff.
- Core Feature 2 → The bound; API Contract 2.
- Core Feature 3 → The reviewed diff; The stderr tail; API Contract 1.
- Core Feature 4 → The QA import; API Contract 3.
- Core Feature 5 → Build Order 2.
- Core Feature 6 → The `deletes` Context kind; API Contract 4.
- Core Feature 7 → The merged Spec's Runs; API Contract 5.
- Core Feature 8 → Build Order 1; Build Order 4.
- Success Metric 1 → Testing Approach 1; Measured outside evidence.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- **Spec 0210** owns the Evidence Snapshot; the import's skip changes only
  which files arrive, and the carry proof refuses a row whose cited evidence
  is absent.
- **Spec 0211** raises the Roundfix Skill's version first; this Spec raises it
  again from the value on its base.
- **The delivery owner's release after a merge** (ADR-0161) calls the same
  classification, so it releases these Runs too.

### Measured outside evidence

Measured on 2026-10-01 with `git diff --no-ext-diff --no-textconv --no-color`
over commit pairs in this repository's history, matching the review records
under this repository's Artifact Directory:

| Candidate | Outcome | Diff bytes | QA evidence | `.agents/skills/` | Rest |
| --- | --- | --- | --- | --- | --- |
| Spec 0194, `c49a0128`..`e8fad6f7` | `reviewed` | 919,745 | — | — | — |
| Spec 0200, `64aff3f7`..`954ad599` | `blocked`, `agent/protocol error` | 1,295,055 | 786,524 | 48,124 | 460,407 |

Every changed path under `.agents/skills/` in Spec 0200's candidate belongs
to a skill the skills lock names at its base or head. This refutes the
adopted Backlog Entry's guess that vendored skills dominated the diff, and
puts the bound between the two observed sizes.

## Testing Approach

1. **Review**, in `internal/cli`, new file `internal/cli/review_scope_test.go`,
   over `reviewCandidateFixture`-style temporary repositories and the
   `reviewCommandRunner` fake: evidence and lock-named skill paths are absent
   from the prompt's diff and listed in the prompt and the record, sorted,
   with their reasons; a skill removed from the lock between base and head is
   still omitted; the QA Report stays in the diff; a malformed lock blocks
   with the scope reason; a diff above the bound is `blocked` with API
   Contract 2's reason and the fake records no prompt; a runtime failure
   carrying stderr records its tail of at most 10 lines and 1,024 bytes; an
   older record without the new fields still reads. The existing
   `TestReviewPromptCarriesTheCandidateDiff` and the review command tests
   pass unedited.
2. **Import**, in `internal/daemon`, new file
   `internal/daemon/qa_prior_pass_skip_test.go`, over the `priorFixture`
   helpers: a failed pass holding `qa/evidence/replay_test.go` and
   `qa/evidence/ledger.txt` imports the report and the text file, skips the
   Go file, publishes `skipped` with `compiled source`, and an existing file
   at the skipped path never yields `path differs`; each listed extension is
   skipped. Every `TestPriorQAPass…` test passes unedited.
3. **Deletes**, in `internal/spec`, `internal/speccheck` and
   `internal/daemon`, new files `internal/spec/context_deletes_test.go`,
   `internal/speccheck/context_deletes_test.go` and
   `internal/daemon/task_deletes_test.go`: the parser accepts the kind and
   refuses a path under two kinds; a Task whose `deletes` path is gone raises
   no `SC-REF-UNRESOLVED`; `UndeclaredTaskPaths`, `declaredGovernedTouches`
   and the Wave collision rule count it; over the task-cycle fixture, a Task
   whose Agent leaves a `deletes` path in place fails its attempt with
   `deletes: <path> still exists` and settles `completed` once the path is
   removed.
   `TestCheckReferenceUnresolved`, `TestContextDeclaredOutput` and
   `TestInstructionContextPathIsNotAudited` pass unedited.
4. **Reconcile**, in `internal/worktree`, new file
   `internal/worktree/merged_spec_leftovers_test.go`, over
   `newMergedHeadTestFixture` and `terminalRunFixture`: a Run whose non-Task
   commits changed the Spec directory before an archive move is `superseded`
   with the Spec-directory count; a dirty worktree whose changes are a
   declared path and a Spec-directory file is `superseded`, and `--apply`
   removes it; a dirty undeclared file outside the Spec keeps it `dirty`; an
   unrepresented Task commit keeps it `unintegrated`; a dirty path set that
   changes between inspection and apply makes apply refuse as stale; an
   absent target branch whose Run the merged head supersedes is released.
   `TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved`,
   `TestMergedHeadRefusesAnUnrepresentedCommit` and
   `TestInspectTerminalRunRequiresArchivedEvidence` pass unedited.

## Build Order

1. The review omission, bound, record fields and stderr tail, with the
   Roundfix Skill `review` reference, the `review` command guide and the
   skill version, task_01 (depends on: none).
2. The QA import skip and the qa-gate skill's evidence rule, task_02 (depends
   on: 1, for the shared skill version record).
3. The `deletes` Context kind, its settlement check and the write-tasks skill
   and template, task_03 (depends on: 2, for the shared skill version record).
4. The merged Spec's Runs, with the Roundfix Skill `reconcile` reference and
   the `reconcile` command guide, task_04 (depends on: 3, for the shared
   skill version record and the Context kinds it reads).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Unreviewed content.** QA evidence and upstream-managed skills merge
  without the reviewer reading them. Evidence is checked by the QA gate and
  the repository gate; vendored skills are third-party content this
  repository never edits. The record lists both.
- **The bound is one repository's measurement.** A provider with a larger
  window would review more; the constant is conservative and named.
- **Forced removal.** `--apply` may now discard uncommitted changes. The
  proof limits them to the archived Spec's own scope, re-proved at apply time.
- **Extension list.** A toolchain outside the built-in profiles may compile
  another extension; the list follows the Baseline's profiles.

## Decisions

- Omit QA evidence and upstream-managed skill copies; do not split the diff.
- Block above a bound below the largest diff reviewed here.
- Skip compiled source in the import rather than refusing the pass.
- Add `deletes` as a fourth Context kind, checked by the Daemon at settlement
  rather than by a new Spec Consistency Check finding.
- A merged Spec supersedes its Runs' Spec-directory work and declared
  leftovers. See ADR-0212.
