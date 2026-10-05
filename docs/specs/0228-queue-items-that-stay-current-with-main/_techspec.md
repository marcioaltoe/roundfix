---
spec: 0228-queue-items-that-stay-current-with-main
prd: _prd.md
created: 2026-10-04
---

# Queue items that stay current with main — Technical Spec

## Executive Summary

Three seams carry the two defects the PRD describes. The owned-skill record
command (`TestEveryOwnedSkillVersionIsRecorded` with `-record-skill-versions`
in `skills/owned_skill_versions_test.go`) refuses content recorded under a
colliding version, so neither a Task nor the derived merge can recover from a
collision. `commandDeliveryWorkflow.ResolveConflict` in
`internal/cli/deliver_workflow.go` treats any conflicted `SKILL.md` as source,
so a version-line conflict parks. `Engine.Retry` in
`internal/delivery/engine.go` refuses every moved head of a
`corrective-spec-required` item. This Spec makes the record command choose
the next free version, lets a derived declaration name line-scoped paths
whose version-line conflicts take the default branch's side before the
regeneration, and lets the retry accept a correction it can prove answers
only the review. The trade-off accepted is a version sequence that may skip
a number: an item's own unmerged entry survives when no conflict regenerates
the record, in exchange for a record that never rewrites a digest
(ADR-0233).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; queue
  items, Runs, blockers, review findings and skill names keep their names.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the record command, the derived
  merge and the retry read and write local files, local Git and the Run
  Database; the merge keeps the existing fetch of the default branch, and no
  credential, forge read or network call is added. No test reaches GitHub.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0233 (this Spec) decides the three
  rules, ADR-0233: "may then change a line-scoped path only on matching
  lines", and ADR-0233: "The next review is round 2 of the same Reviewer
  Lineage". ADR-0189 ties an owned skill's version to its content, ADR-0189: "A
  version names one content". ADR-0192 keeps the whole-file resolution of
  declared derived paths, ADR-0192: "Any other conflicted path aborts the
  merge", now with line-scoped paths beside it. ADR-0165 still parks a
  finding about the delivered change, ADR-0165: "an archived Spec is never
  edited to absorb a finding". ADR-0197 bounds the Reviewer Lineage,
  ADR-0197: "Round 2 reviews only the diff from the round-1 head to the
  current head". ADR-0223 is narrowed for a review-only correction,
  ADR-0223: "A park for a corrective Spec still refuses a moved head".
  ADR-0229 keeps the operator-archive anchor unchanged, ADR-0229: "The
  Run start head is a weak anchor". ADR-0187 and
  ADR-0189 govern both skill edits. ADR-0184 has this TechSpec state its
  changed command surfaces, ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts". The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency. ADR-0178,
  ADR-0179 and ADR-0166 decide each Task commit's grant and record undeclared
  paths, and ADR-0182 settles a Task on the facts its gate checks. ADR-0096
  and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row
  carry, ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA
  row records, when it is observed again and its evidence snapshot, ADR-0169
  cites ADR-0165 but decides the review's merge-base diff, and ADR-0196 cites
  ADR-0169 but decides when a review finding parks; this Spec changes none of
  them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization for owned
  skills ("considere autorizado a ajustar todas as skills se necessário"),
  for guides and `.roundfixrc.yml` ("Autorizar os dois") and for this cycle
  ("Aprovado", 2026-10-04). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0228-queue-items-that-stay-current-with-main/_authorization.md`;
  bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`, `.roundfixrc.yml`,
  `docs/agents/specific-repository.md`,
  `skills/implement-task/SKILL.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Three existing components change; nothing new is added beside one engine
dependency.

- **Owned-skill record.** `TestEveryOwnedSkillVersionIsRecorded` reads each
  embedded skill's declared version and folder digest and checks them against
  `skills/testdata/owned-skill-versions.json` through
  `checkOwnedSkillVersionRecord`. With `-record-skill-versions` it appends
  new, higher versions. The repository already declares this command as the
  regeneration of the record in `.roundfixrc.yml`.
- **Derived merge.** At a `CONFLICTING` Pull Request, `Engine.checkCandidate`
  calls `ConflictResolver.ResolveConflict`. The workflow merges the fetched
  default branch with `--no-ff --no-commit`, takes the default branch's
  bytes for each conflicted path a `DerivedPathDeclaration` matches, runs each
  matched command, and commits only when the regeneration changed no path
  outside the declarations.
- **Delivery Retry.** `Engine.Retry` returns a `corrective-spec-required`
  item to `reviewing` only when the item head equals the parked candidate.
  The review the queue runs next is `roundfix review` in the item worktree,
  whose Reviewer Lineage makes a descendant head round 2 and closes a third
  round at the ceiling (ADR-0197). Findings and their operator dispositions
  live in the review record and the disposition ledger under the artifact
  directory.

## Implementation Design

### Interfaces

```go
// internal/config: a derived declaration may name line-scoped paths.
type DerivedPathDeclaration struct {
	Paths      []string               `yaml:"paths"`
	Lines      *DerivedLineDeclaration `yaml:"lines"`
	Regenerate string                 `yaml:"regenerate"`
}
type DerivedLineDeclaration struct {
	Paths []string `yaml:"paths"` // same path rules as Paths
	Match string   `yaml:"match"` // RE2 pattern a whole line must match
}
func (declaration DerivedPathDeclaration) MatchesLines(name string) bool

// internal/delivery: the retry asks the workflow to prove a correction.
type ReviewCorrection struct {
	Accepted bool
	Reason   string // why the proof refused; empty when accepted
}
type ReviewCorrectionProver interface {
	ProveReviewCorrection(ctx context.Context, workDir string,
		archivedSpecs []string, candidate, head string) (ReviewCorrection, error)
}
// EngineDependencies gains Corrections ReviewCorrectionProver.
```

### The record command chooses the version

1. Without `-record-skill-versions` the check is unchanged: content that
   differs from the digest recorded under its version, or a version that is
   not recorded, fails with today's messages.
2. With the flag, for each owned skill whose content is not recorded under
   its declared version: when the declared version is higher than every
   recorded version, record it as today; otherwise the version becomes one
   patch above the highest recorded version.
3. A raise rewrites every front-matter line matching `^ *version: ` in the
   canonical `../.agents/skills/<name>/SKILL.md` and in the mirror
   `<name>/SKILL.md`, then computes the digest from the mirror on disk, since
   the embedded copy predates the rewrite, and records that version and
   digest.
4. A recorded entry is never replaced or removed, and versions stay
   ascending. A skill whose mirror differs from its canonical copy outside
   the version lines refuses with a message naming `make skills-sync`.

### Line-scoped derived paths

1. `lines` is optional. When present it needs a non-empty `paths` list,
   validated like `paths`, and a non-empty `match` that compiles as a Go
   regular expression; otherwise the config error names
   `delivery.derived_paths[<n>].lines`.
2. A conflicted path that matches a declaration's `paths` is resolved as
   today; `paths` takes precedence over `lines`.
3. A conflicted path that matches a declaration's `lines.paths` is resolved
   only when every conflict hunk in its working file holds, on both sides,
   lines that each match `match`. The resolver writes the default branch's
   side of each hunk, stages the file and marks the declaration matched; a
   hunk with any other line makes the path a source path, and the merge
   aborts and parks `pull-request-conflict: <path>` as today. The resolver
   runs its own merge with `merge.conflictStyle=merge`, so the hunks have the
   two-sided form Git documents.
4. After the regeneration, a changed path is accepted when it matches a
   declaration's `paths`, or when it matches `lines.paths` and its diff from
   the pre-regeneration snapshot keeps the line count and changes only lines
   that match `match` before and after. Any other change aborts the merge
   with today's `regenerated <path> outside delivery.derived_paths` reason.
5. This repository's record declaration in `.roundfixrc.yml` gains
   `lines: {paths: [.agents/skills/*/SKILL.md, skills/*/SKILL.md], match: '^ *version: '}`,
   and the regeneration command is unchanged.

### A review-only correction returns to round 2

1. The `corrective-spec-required` branch of `Engine.Retry` keeps an unchanged
   head returning to `reviewing`.
2. When the head moved, the engine has a `ReviewCorrectionProver`, and the
   head is not empty, it calls `ProveReviewCorrection` with the item
   worktree, the slugs the blocker names, the parked candidate and the head.
   On acceptance the retry appends the head to the candidate commits, sets
   `reviewing` and hands the item to the owner.
3. The workflow's proof accepts only when all of these hold:
   - `git merge-base --is-ancestor <candidate> <head>` succeeds;
   - the review record persisted for the item worktree names the candidate
     as its head and has the outcome `findings`;
   - every finding of that record not dismissed by validation has exactly
     one disposition in the ledger at that head: dismissed with non-empty
     evidence, or fixed by a commit that descends from the candidate and is
     an ancestor of the head, the same test the review's ceiling applies, now
     shared through one helper in `internal/cli/review_lineage.go`;
   - every path `git diff --name-only --no-renames <candidate> <head>` lists
     lies under the archive directory of a slug the blocker names.
4. Any failed condition returns `Accepted: false` with a reason naming it; a
   Git or file error returns an error, which the retry logs and treats as a
   refusal. A refusal keeps today's text and inserts the reason in
   parentheses after the item head, as API Contract 3 shows; without a
   prover the text is byte-identical to today's.
5. The next review runs `roundfix review` at the head. The Reviewer Lineage
   makes it round 2 with the round-1 findings and dispositions in its prompt;
   a later correction is closed at the ceiling as ADR-0197 already decides.

### Data Models

`DerivedPathDeclaration` gains the optional `lines` field. The record file,
the review record, the disposition ledger and the queue item keep their
fields.

### API Contracts

1. API Contract: owned-skill record. Input: the record command on a tree
   whose skill content is not recorded under its declared version, which is
   not above every recorded version. Output: both version fields of the skill
   and its mirror hold one patch above the highest recorded version, and the
   record gains that version and digest. Failure: a mirror that differs from
   its canonical copy outside the version lines fails, naming
   `make skills-sync`, and writes nothing.
2. API Contract: line-scoped derived merge. Input: a `CONFLICTING` Pull
   Request whose conflicted paths each match a declaration's `paths`, or its
   `lines.paths` with every hunk confined to matching lines. Output: the
   merge commits with the `Roundfix-Delivery: derived-merge` trailer and the
   item returns to `gating`. Failure: any other conflicted path parks
   `pull-request-conflict: <path>`; a regeneration that changes another path
   or another line aborts and parks as today.
3. API Contract: review-only correction. Input: a parked
   `corrective-spec-required` item whose head moved. Output, when the proof
   holds: stage `reviewing`, the head appended to the candidate commits,
   blocker cleared. Failure: refusal with exit `2`, the item unchanged, and
   `retry Delivery Queue item "<slug>": archived Specs <slugs> were reviewed at parked candidate head "<candidate>", but the item head is "<head>" (<reason>); author a corrective Spec with its own authorization and QA gate`.

### Surface Transcripts

None. `roundfix deliver retry` prints its existing result line on success
and its existing refusal, with the proof's reason inserted, on failure;
reproducing a successful retry through the built binary would start a
detached queue owner that runs a review, and a derived merge needs a Pull
Request. API Contracts 1 to 3 are asserted through the record test, the
workflow and the engine with real Git, a real queue store and disposable
artifact directories.

## Vocabulary Contract

No new glossary term. The Spec uses **Delivery Retry**, **Reviewer
Lineage**, **Park Class** and the derived-path declaration as the `deliver`
and configuration guides already do; a line-scoped path is a field of that
declaration, documented in the configuration guide by task_02. task_04
documents the retry and record rules in the `deliver` and `review` guides and
the Roundfix Skill's references.

## Coverage Map

- Goal 1 → The record command chooses the version; Line-scoped derived
  paths; API Contracts 1 and 2.
- Goal 2 → The record command chooses the version, step 2; Build Order 4
  (the `implement-task` skill and the repository rule).
- Goal 3 → A review-only correction returns to round 2; API Contract 3.
- Goal 4 → The record command chooses the version, steps 1 and 4;
  Line-scoped derived paths, steps 2 to 4; A review-only correction returns
  to round 2, steps 1 and 4.
- Core Feature 1 → The record command chooses the version; API Contract 1.
- Core Feature 2 → Line-scoped derived paths; API Contract 2.
- Core Feature 3 → A review-only correction returns to round 2; API
  Contract 3.
- Core Feature 4 → Build Order 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1, 2 and 3.

## Integration Points

- Local Git: the existing default-branch fetch and merge of the derived
  merge, `git merge-base --is-ancestor` and `git diff --name-only`.
- The artifact directory: the review record and disposition ledger the
  Pre-PR Review Command already writes.
- No forge, provider or network boundary is added.

## Testing Approach

1. **Record**, in the new file `skills/owned_skill_version_raise_test.go`:
   `TestRecordingRaisesAVersionRecordedWithOtherContent` and
   `TestRecordingRaisesAVersionBelowTheHighestRecorded` over
   `checkOwnedSkillVersionRecord`, and
   `TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror` over
   temporary canonical and mirror roots, asserting both front-matter fields,
   the recorded digest of the rewritten mirror, and a refusal that writes
   nothing for a drifted mirror. The record-mode cases `changed digest` and
   `lower unrecorded version` of `TestRecordingNeverReplacesARecordedVersion`
   change to expect the raise, the one declared break there; every
   non-recording refusal test passes unedited.
2. **Derived merge**, in `internal/config/delivery_derived_lines_test.go`
   (`TestDerivedLineDeclarationsAreReadAndValidated`) and
   `internal/cli/deliver_derived_lines_test.go` with real Git repositories
   and a shell regeneration command:
   `TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated`,
   `TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge` and
   `TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge`. The tests
   of `internal/cli/deliver_conflict_test.go` pass unedited.
   `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` changes its expected
   record declaration, the one declared break there.
3. **Retry**, in `internal/delivery/review_correction_retry_test.go` with a
   real queue store and a fake prover:
   `TestRetryReturnsAReviewOnlyCorrectionToReview` and
   `TestRetryRefusesACorrectionTheProofRejects`; and in
   `internal/cli/deliver_review_correction_test.go` with a real repository
   and a disposable artifact directory:
   `TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord`,
   `TestACorrectionOutsideTheArchivedSpecIsRefused`,
   `TestACorrectionWithAnUndisposedFindingIsRefused` and
   `TestACorrectionThatDoesNotDescendIsRefused`. The tests of
   `internal/delivery/corrective_spec_test.go`,
   `internal/delivery/archived_head_retry_test.go` and
   `internal/cli/review_lineage_test.go` pass unedited.

## Build Order

1. The record command chooses the version, with Testing Approach 1, task_01
   (depends on: none).
2. Line-scoped derived paths, the repository declaration and the
   configuration guide, with Testing Approach 2, task_02 (depends on: none).
3. A review-only correction returns to round 2, with Testing Approach 3,
   task_03 (depends on: 2, which changes the same workflow file).
4. The Roundfix Skill's `deliver` and `review` references, their user guides,
   the `implement-task` skill and the repository's version rule, with both
   skill versions recorded by the record command, task_04 (depends on: 1).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A pattern that matches a body line.** A conflicting body line that
  matches `^ *version: ` resolves from the default branch. The repository
  gate after the merge is the net, as for any derived path.
- **A skipped version.** A Task that recorded too early and records again
  leaves its first entry; ADR-0233 accepts it, and a conflict regenerates the
  record from the default branch's bytes.
- **A forged disposition.** Dispositions come from `roundfix review dispose`,
  which already checks the fixing commit's ancestry; the proof checks it
  again against the head, and round 2 still reviews the correction.
- **Concurrent skill edits until merge.** This Spec's own raise follows the
  rule it replaces; the operator orders the queue against 0227 and 0229.

## Decisions

- The record command chooses the version and never replaces a digest. See
  ADR-0233.
- Line-scoped derived paths take the default branch's side of hunks confined
  to matching lines; `paths` keeps its whole-file rule. See ADR-0233.
- A review-only correction after archive returns to round 2. See ADR-0233.
- Refreshing an item from the default branch before its Run was rejected,
  and so was declaring every `SKILL.md` a derived path. See ADR-0233.
