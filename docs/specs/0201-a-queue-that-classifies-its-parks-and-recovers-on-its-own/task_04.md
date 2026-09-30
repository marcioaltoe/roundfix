---
task: task_04
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: pending
type: backend
complexity: high
---

# Task 04: A conflicting Pull Request parks at once and a derived conflict is resolved by regeneration

## Overview

On 2026-09-30 PR #298 was `CONFLICTING`; GitHub never started its checks, and the owner parked it `checks-timeout` thirty minutes later. Every conflicted file but one was derived. This Task reads GitHub's `mergeable` state while the owner waits for checks, resolves a conflict confined to paths declared in `delivery.derived_paths` by merging the default branch and regenerating, and parks any other conflict as `pull-request-conflict` with its paths (ADR-0192). It also brings the Roundfix skill's delivery reference up to every behavior this Spec adds.

## Requirements

1. MUST add the Project Config key `delivery.derived_paths`, a list of `paths` and `regenerate` pairs, as the TechSpec's Data Models section states. An empty `paths` list, an absolute or `..` entry, or an empty `regenerate` MUST be a config error that names the key. A config without the key MUST load exactly as today.
2. MUST add `mergeable` to the Pull Request fields `GitHubCLI` reads and return it as `CheckReport.Mergeable`.
3. MUST make `checkCandidate` send a `CONFLICTING` report to the optional `ConflictResolver` on the first read, and keep an `UNKNOWN` report pending. A nil `ConflictResolver` MUST park `pull-request-conflict` with no paths.
4. MUST implement `ResolveConflict` in `internal/cli/deliver_workflow.go` in the five steps of the TechSpec's Conflicts section: a clean worktree at the candidate head; a fetch and `git merge --no-ff --no-commit` of the remote default branch; `git merge --abort` and the undeclared paths when any conflicted path matches no declaration; `--theirs` for each conflicted path and each matched `regenerate` command once, in declaration order, through the Verification executor; and an abort when the regeneration changed an undeclared path, otherwise a merge commit with the trailer `Roundfix-Delivery: derived-merge`.
5. MUST make the engine park `pull-request-conflict: <path>, …` on source paths, classified `conflict` with the next command the TechSpec's Park Classes table states, and on a new head append it to the candidate commits and move to `gating`.
6. MUST make `createPullRequest` re-read an existing Pull Request whose head is an earlier candidate commit of the item, every check interval up to the check timeout, before it fails.
7. MUST make `Retry` of `pull-request-conflict` accept a head that descends from the candidate, through task_03's acceptance, and target `reviewing`.
8. MUST add a test that reproduces Surface Transcript 1 through `roundfix deliver status`.
9. MUST describe the key in `docs/user-guide/configuration.md`, and the conflict park and derived merge in `docs/user-guide/commands/deliver.md`.
10. MUST describe the prerequisites, the Park Classes and the `Park:` line, the one re-run, the operator-archived retry and the conflict handling in the Delivery queue section of `.agents/skills/roundfix/references/deliver.md`, the file Spec 0194 creates, with the phrases `prerequisite-unmerged`, `pull-request-conflict`, `qa-environment-partial`, `flaky-check` and `delivery.derived_paths`. It MUST raise both version fields of `.agents/skills/roundfix/SKILL.md`, run `make skills-sync`, and record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`. It MUST NOT change any text inside the `### QA settlement` section.
11. MUST NOT add this repository's own declaration to `.roundfixrc.yml`, rebase, force-push, or reach GitHub or a network remote in any test.

## Subtasks

- [ ] Add and validate `delivery.derived_paths`.
- [ ] Read `mergeable` and stop the check wait on a conflict.
- [ ] Merge, regenerate and commit a derived conflict; abort and park a source conflict.
- [ ] Tolerate a lagging Pull Request head and accept the conflict retry.
- [ ] Update the configuration and delivery guides and the Roundfix skill, with a test for each criterion.

## Acceptance Criteria

- [ ] A `CONFLICTING` report stops the wait on the first read; an `UNKNOWN` report keeps waiting.
- [ ] In a repository with a bare origin, a conflict on one declared derived file produces a merge commit whose file equals the regeneration command's output, carries the trailer, and returns the item to `gating`.
- [ ] A conflict that also touches an undeclared file aborts the merge, leaves the worktree at the candidate head and parks `pull-request-conflict` naming only that file; a regeneration that writes an undeclared file aborts the same way.
- [ ] A Pull Request head that still reports the previous candidate is read again, and `deliver status` reproduces Surface Transcript 1.
- [ ] The skill's delivery reference and its mirror carry the five phrases, the recorded skill version test passes, and `TestSettlementGuidanceIsOneTable` passes.

## Context

- interface: `internal/config/config.go`
- creates: `internal/config/delivery_derived_paths_test.go`
- interface: `internal/delivery/github.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/park_class.go`
- creates: `internal/delivery/conflict_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_conflict_test.go`
- creates: `internal/cli/deliver_park_status_test.go`
- interface: `docs/user-guide/configuration.md`
- creates: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/deliver.md`
- creates: `skills/roundfix/references/deliver.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProjectConfigReadsDerivedPathDeclarations|TestDerivedPathDeclarationsRefuseUnsafeEntries|TestAConflictingPullRequestStopsTheCheckWaitAtOnce|TestAnUnknownMergeableStateKeepsWaiting|TestADerivedConflictReturnsTheItemToTheGate|TestASourceConflictParksWithItsPaths|TestALaggingPullRequestHeadIsReadAgain|TestRetryResumesAResolvedConflictAtReview|TestGitHubCLIReadsTheMergeableState)$" ./internal/config ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProjectConfigReadsDerivedPathDeclarations TestDerivedPathDeclarationsRefuseUnsafeEntries TestAConflictingPullRequestStopsTheCheckWaitAtOnce TestAnUnknownMergeableStateKeepsWaiting TestADerivedConflictReturnsTheItemToTheGate TestASourceConflictParksWithItsPaths TestALaggingPullRequestHeadIsReadAgain TestRetryResumesAResolvedConflictAtReview TestGitHubCLIReadsTheMergeableState; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestADerivedConflictIsMergedAndRegenerated|TestASourceConflictAbortsTheMerge|TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge|TestDeliverStatusReproducesSurfaceTranscriptOne)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestADerivedConflictIsMergedAndRegenerated TestASourceConflictAbortsTheMerge TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge TestDeliverStatusReproducesSurfaceTranscriptOne; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `for pair in "docs/user-guide/configuration.md|derived_paths" "docs/user-guide/commands/deliver.md|pull-request-conflict" ".agents/skills/roundfix/references/deliver.md|prerequisite-unmerged" ".agents/skills/roundfix/references/deliver.md|pull-request-conflict" ".agents/skills/roundfix/references/deliver.md|qa-environment-partial" ".agents/skills/roundfix/references/deliver.md|flaky-check" ".agents/skills/roundfix/references/deliver.md|delivery.derived_paths"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task none of the phrases is in the configuration guide, the delivery guide or the skill's delivery reference, so the command fails.

## References

- `_prd.md` → Goal 2; User Stories 2-3; Core Feature 2; Success Metric 2; Recorded limits
- `_techspec.md` → Conflicts; Data Models; API Contracts 1, 3, 4; Surface Transcript 1; Testing Approach 4-5; Build Order 4
- ADR-0192; ADR-0149; ADR-0178; ADR-0081; ADR-0160; ADR-0169; ADR-0174; ADR-0196; ADR-0197
