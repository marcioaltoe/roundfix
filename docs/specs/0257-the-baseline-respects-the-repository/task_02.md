---
task: task_02
spec: 0257-the-baseline-respects-the-repository
status: completed
type: backend
complexity: high
---

# Task 02: A Baseline update records a declared decision

## Overview

`roundfix baseline update` never changes a recorded decision, so an adopter
that mounts Better Auth at `/auth` must restate every decision through
`baseline plan` and `baseline apply` to change one. This Task adds
`--decision <id=value>` to the update: the declared values are recorded before
the plan is built, held to the plan's normalization, reported, and covered by
the Plan Digest (ADR-0263). It is verifiable on its own through temporary
adopted repositories. It answers the Backlog Entry of 2026-10-12,
[the Baseline hardcodes the Better Auth base path](references/2026-10-12-the-baseline-hardcodes-the-better-auth-base-path.md),
whose expected adopter declares `/auth` at its next update.

This is an authorized tooling Task. It may change only the files in its
Context and this Task file.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12 on the hardcoded Better Auth
   base path by implementing API Contract 3 through
   `baseline.DeclareManifestDecisions` of `_techspec.md` → Interfaces, called
   after the Setup Manifest resolves and before the plan is built.
2. MUST refuse, with exit `2`, category `invalid` and no write, a malformed
   pair, an unknown decision, a decision the profile does not select, a
   profile-fixed decision and `preservation.mode`. A declared decision MUST
   remove itself from `newDecisions`, so it needs no `--adopt-suggested`.
3. MUST add the always-present `declaredDecisions` array to
   `roundfix/baseline-update-result/v1` and the text lines of API Contract 3.
   Without `--yes` or `--confirm-plan` nothing is written.
4. MUST implement API Contract 4 in the three synopsis lines and the update
   options, then re-record `docs/references/behavior-surfaces.json` with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
5. MUST create `internal/cli/baseline_update_decision_test.go` with
   `TestBaselineUpdateRecordsADeclaredDecision`,
   `TestBaselineUpdateDeclaresTheBetterAuthBasePath` and
   `TestBaselineUpdateRefusesAnUndeclarableDecision` of `_techspec.md` →
   Testing Approach. The Better Auth test MUST adopt a temporary Standard
   TypeScript repository, run the update with
   `--decision auth.base_path=/auth --yes`, assert Success Metric 1 on the
   written backend guide and Setup Manifest, and require a second update to
   report `current`. If a test in `internal/cli/baseline_update_test.go` pins
   the result's field set, it changes only to include `declaredDecisions`.
6. MUST add the task_02 text of `_techspec.md` → Skill and document texts to
   `.agents/skills/roundfix/references/baseline.md` and
   `docs/user-guide/commands/baseline.md`, then run `make skills-sync`,
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and `make skills-sync` again.
7. MUST prove the refusal can fail: the Result records one sabotage, for
   example accepting `preservation.mode`, and the test that failed.
8. MUST NOT let the update change the profile or the preservation mode, and
   MUST NOT change the catalog, any module or `baseline plan`.

## Subtasks

- [ ] Add the flag, its parsing and its refusals.
- [ ] Record declared decisions before planning and report them.
- [ ] Update the help and re-record the surface record.
- [ ] Add the three tests.
- [ ] Update the command reference, the skill and its version.

## Acceptance Criteria

- [ ] `--decision auth.base_path=/auth --yes` records `/auth` and moves the
      backend guide and the HTTP exception in a Standard TypeScript adopter.
- [ ] A preview with `--decision` writes nothing and lists the declared
      decision.
- [ ] Every undeclarable input exits `2` with no write.
- [ ] The help, the command reference and the skill name `--decision`.

## Context

- instruction: `docs/adr/0263-the-better-auth-base-path-is-a-repository-decision-a-baseline-update-records.md`
- instruction: `internal/cli/baseline_profile.go`
- interface: `internal/baseline/update.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/baseline_update_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/baseline_update_decision_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBaselineUpdateRecordsADeclaredDecision|TestBaselineUpdateDeclaresTheBetterAuthBasePath|TestBaselineUpdateRefusesAnUndeclarableDecision)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateRecordsADeclaredDecision TestBaselineUpdateDeclaresTheBetterAuthBasePath TestBaselineUpdateRefusesAnUndeclarableDecision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `help="$(go run -buildvcs=false ./cmd/roundfix baseline update --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -qF -- '[--decision <id=value> ...]' || { printf 'the update help lacks --decision\n' >&2; exit 1; }; for file in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "records or changes a decision the repository's profile selects before the plan is built" || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md || { printf 'the baseline reference mirror differs\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; go test -tags docscontract -count=1 -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0; before this Task the help has no --decision, so the command fails at its first check.

## References

- `_prd.md` → Goal 2; User Stories 1 and 3; Core Feature 2; Success Metric 1
- `_techspec.md` → API Contract 3; API Contract 4; Interfaces; Skill and document texts; Invariant 3; Testing Approach; Build Order 2
- ADR-0263; ADR-0118; ADR-0256


## Result

Implemented the declared-decision slice for Daemon Verification. Task status,
Subtasks and Acceptance Criteria remain for Daemon settlement.

`baseline.DeclareManifestDecisions` now overlays declared values after Setup
Manifest resolution and before planning, using the resolved built-in or
repository-owned Profile and the plan's existing normalization. It preserves
the prior manifest for recorded-value evidence, recomputes unresolved required
decisions, and rejects profile-fixed decisions and `preservation.mode` before
planning. The update accepts repeatable `--decision` values through the existing
inline value parser, reports an always-present `declaredDecisions` array and
compact JSON text lines, and uses its existing confirmation and Plan Digest
contract for the resulting managed refresh.

### Acceptance evidence

| Criterion | Implementation and focused-check evidence |
| --- | --- |
| `/auth` is recorded and moves the backend guide and HTTP exception | `TestBaselineUpdateDeclaresTheBetterAuthBasePath` adopts a temporary Standard TypeScript repository, previews without writes, applies `--decision auth.base_path=/auth --yes`, checks the recorded path and provider scope, exactly one matching HTTP exception, both `/auth/*` guide occurrences and the base-path sentence, no `/api/auth/*`, and a second update reporting `current` with zero file changes. |
| A declaration preview writes nothing and lists the declaration | `TestBaselineUpdateRecordsADeclaredDecision` checks JSON values and previous recorded values, exact text lines, repeated declarations, empty arrays without declarations, unchanged repository bytes, distinct digests for changed values, rejection of an earlier digest for a different declaration, and application of the confirmed preview digest. It also covers a new required decision without `--adopt-suggested`, remaining suggestions, and a repository-owned Profile. |
| Every undeclarable input exits `2` without writes | `TestBaselineUpdateRefusesAnUndeclarableDecision` checks malformed pairs, empty identifiers and invalid empty values, unknown and unselected decisions, changed and identical profile-fixed values, `preservation.mode`, invalid typed values, duplicates, a valid declaration followed by an invalid declaration, and seven invalid Better Auth paths. Each checks category `invalid` and unchanged repository bytes. |
| Help, command reference and skill name `--decision` | The new help subtest checks all three synopsis surfaces and the exact option description. Both references contain the TechSpec's authored paragraph. The Behavior Surface Record was regenerated, and skill sync plus the generated version record keep the canonical and distributed Roundfix skill at `0.1.64`. |

### Focused checks and regeneration

- Before implementation,
  `go test ./internal/cli -run '^TestBaselineUpdateRecordsADeclaredDecision$/changed_recorded_decision$' -count=1`
  exited `1`: the update refused the undefined `--decision` flag with CLI exit
  `2`. This established the missing behavior.
- After implementation,
  `go test ./internal/cli -run '^TestBaselineUpdate(RecordsADeclaredDecision|DeclaresTheBetterAuthBasePath|RefusesAnUndeclarableDecision)/' -count=1`
  exited `0`, including the new-decision, confirmation, Better Auth and refusal
  cases.
- Final focused regression run:
  `go test ./internal/cli -run '^TestBaselineUpdate' -count=1` exited `0`.
  An earlier run passed its assertions but suiteguard rejected a concurrently
  regenerated behavior record; the final run followed regeneration with a
  stable repository and passed the guard.
- `go test ./internal/baseline -run '^TestResolveManifestInput' -count=1`
  exited `0`.
- Ran `make skills-sync`,
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  then `make skills-sync`; all exited `0`. Repeated after the final paragraph
  formatting so the recorded version and mirror describe the final content.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited `0` after the final help and guide edits.
- The repository-required `make baseline-digests` exited `0` with
  `changed: false`; no catalog, module or Baseline derived-artifact diff remains.
- `git diff --check` exited `0`; `gofmt -l` listed no changed Go file.

### Sabotage

A temporary Go overlay replaced the `preservation.mode` refusal with
`continue`, causing the update to silently accept that prohibited input. Ran
`go test -overlay <temporary-overlay.json> ./internal/cli -run '^TestBaselineUpdateRefusesAnUndeclarableDecision$/preservation_mode$' -count=1`.
The test exited `1` and failed at `refusal exit=0`: the sabotaged update returned
`ok: true`, state `verified`, and empty stderr. The overlay was removed; no
repository source was changed by the sabotage. The final focused regression
run above exercised the original refusal again.

### Scope and handoff

The initial worktree change was the Daemon's Task 02 status. Final changed
paths are confined to this Task's declared Context and this Task file. No
other Task or Task Graph change, catalog/module change, `baseline plan`
behavior change, commit, push or Pull Request was made. The Task's authored
Verification commands and the repository Verification were not run; both
remain Daemon-owned. No follow-up implementation was added to this slice.
