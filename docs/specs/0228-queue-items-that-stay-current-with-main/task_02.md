---
task: task_02
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: backend
complexity: high
---

# Task 02: A derived merge resolves conflicts confined to declared version lines

## Overview

On 2026-10-04 the operator resolved two Pull Request conflicts by hand by
taking the default branch's version lines of the Roundfix Skill and raising
to the next patch (Backlog Entry
[queued Specs raise the same skill version](references/2026-10-04-queued-specs-raise-the-same-skill-version.md),
recorded 2026-10-04; operator log entries 157 and 159). The derived merge
treats any conflicted `SKILL.md` as source and parks. This Task adds
line-scoped paths to a derived-path declaration, resolves a conflict whose
hunks hold only matching lines by taking the default branch's side, limits
the regeneration to matching lines of those paths, and declares the version
fields of every `SKILL.md` for this repository's record command (ADR-0233).

## Requirements

1. MUST add the optional `lines` field to `DerivedPathDeclaration`, with
   `paths` validated like the declaration's `paths` and `match` a non-empty Go
   regular expression, refusing an incomplete or invalid `lines` with a config
   error naming `delivery.derived_paths[<n>].lines`, as "Line-scoped derived
   paths" step 1 of the TechSpec states.
2. MUST implement steps 2 to 4 of "Line-scoped derived paths" in
   `commandDeliveryWorkflow.ResolveConflict`: `paths` keeps its whole-file
   rule and precedence; a conflicted `lines.paths` file whose every conflict
   hunk holds only matching lines on both sides takes the default branch's
   side of each hunk; any other hunk parks `pull-request-conflict: <path>`;
   after the regeneration a `lines.paths` file may change only on matching
   lines with its line count kept, and any other change aborts with today's
   `regenerated <path> outside delivery.derived_paths` reason. The resolver
   runs its own merge with `merge.conflictStyle=merge`.
3. MUST add to the record declaration in `.roundfixrc.yml` the line-scoped
   paths `.agents/skills/*/SKILL.md` and `skills/*/SKILL.md` with the match
   `^ *version: `, leaving its `paths` and `regenerate` command unchanged, and
   update the expected declaration of
   `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` in
   `internal/config/verification_tools_test.go`, the one declared break.
4. MUST document `lines` in the "Delivery conflict recovery" section of
   `docs/user-guide/configuration.md`, including the phrase "takes the default
   branch's side of each conflict hunk".
5. MUST add `internal/config/delivery_derived_lines_test.go` with
   `TestDerivedLineDeclarationsAreReadAndValidated`, and
   `internal/cli/deliver_derived_lines_test.go` with
   `TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated`,
   `TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge` and
   `TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge`, over real Git
   repositories with a local default-branch remote and a shell regeneration
   command, as Testing Approach 2 describes. The tests of
   `internal/cli/deliver_conflict_test.go` pass unedited, and no test reaches
   GitHub or the network.
6. MUST change no Governed Path other than `.roundfixrc.yml`, the one path of
   `_authorization.md` this Task touches.

## Subtasks

- [ ] Read and validate line-scoped declarations.
- [ ] Resolve conflict hunks confined to matching lines and refuse every other hunk.
- [ ] Limit the regeneration to matching lines of line-scoped paths.
- [ ] Declare the repository's version lines and document the field.

## Acceptance Criteria

- [ ] Two branches that changed the same version line to different values and
      appended to the same derived record merge with the default branch's
      version line and a regenerated record.
- [ ] A conflict hunk with any non-matching line parks
      `pull-request-conflict: <path>`; a regeneration that changes a
      non-matching line aborts the merge.
- [ ] The repository declaration carries the line-scoped `SKILL.md` paths, and
      every existing derived-merge test passes unedited.
- [ ] The only Governed Path changed is `.roundfixrc.yml`.

## Context

- interface: `internal/config/config.go`
- interface: `internal/config/delivery.go`
- creates: `internal/config/delivery_derived_lines_test.go`
- interface: `internal/config/verification_tools_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_derived_lines_test.go`
- interface: `.roundfixrc.yml`
- interface: `docs/user-guide/configuration.md`
- instruction: `internal/cli/deliver_conflict_test.go`
- instruction: `internal/config/delivery_derived_paths_test.go`
- instruction: `docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "takes the default branch's side of each conflict hunk" || { printf 'missing phrase in configuration guide\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestDerivedLineDeclarationsAreReadAndValidated|TestThisRepositoryDeclaresItsToolsAndDerivedPaths|TestProjectConfigReadsDerivedPathDeclarations|TestDerivedPathDeclarationsRefuseUnsafeEntries)$" ./internal/config 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDerivedLineDeclarationsAreReadAndValidated TestThisRepositoryDeclaresItsToolsAndDerivedPaths TestProjectConfigReadsDerivedPathDeclarations TestDerivedPathDeclarationsRefuseUnsafeEntries; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide lacks the phrase and the new test does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated|TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge|TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge|TestADerivedConflictIsMergedAndRegenerated|TestASourceConflictAbortsTheMerge|TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge|TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch|TestConflictRecoveryRunsMatchedDeclarationsInOrder)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge TestADerivedConflictIsMergedAndRegenerated TestASourceConflictAbortsTheMerge TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch TestConflictRecoveryRunsMatchedDeclarationsInOrder; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three line-scoped tests do not exist, so the command fails.

## References

- `_prd.md` → Goals 1 and 4; Core Feature 2; Success Metric 2; Success Metric 4
- `_techspec.md` → Line-scoped derived paths; Interfaces; Data Models; API Contract 2; Testing Approach 2; Build Order 2
- ADR-0233; ADR-0192
