---
task: task_01
spec: 0257-the-baseline-respects-the-repository
status: completed
type: backend
complexity: high
---

# Task 01: A repository records its Better Auth base path

## Overview

The Standard TypeScript backend guide names `/api/auth/*` for every adopter,
and moving the provider's scope duplicates the `http.contract` exception
instead of moving it. This Task adds the optional `auth.base_path` decision,
derives the provider exception from it when it is recorded, renders the
recorded path, and leaves every adopter that records nothing byte-identical
(ADR-0263). It is verifiable on its own through plans built from the embedded
catalog. It answers the Backlog Entry of 2026-10-12,
[the Baseline hardcodes the Better Auth base path](references/2026-10-12-the-baseline-hardcodes-the-better-auth-base-path.md).

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12 on the hardcoded Better Auth
   base path by adding the `auth.base_path` entry of `_techspec.md` → Data
   Models to the decision catalog after `auth.provider`, listing it after
   `auth.provider` in the `entryDecisions` of the Standard TypeScript and Go
   CLI TypeScript profiles, raising `template.guide.backend` from 3 to 4 with
   the new token, and ending the backend template with
   `{{auth.provider}}{{auth.base_path}}`.
2. MUST implement API Contract 1 in project-decision normalization: a
   recorded valid base path sets the provider's scope to the path followed by
   `/*`, removes from `http.contract` the exception equal to the provider
   exception as recorded, and then derives `http.contract` as today. It MUST
   refuse an invalid value and a base path without `auth.provider` with the
   messages of API Contract 1, and MUST leave normalization unchanged when the
   base path is absent, including a recorded provider scope other than
   `/api/auth/*`.
3. MUST implement API Contract 2: the token renders the empty string while the
   decision is absent, and the exact sentence of API Contract 2 when it is
   recorded.
4. MUST create `internal/baseline/better_auth_base_path_test.go` with the four
   tests of `_techspec.md` → Testing Approach, and MUST change
   `TestDivergenceRendersProbe` in `internal/baseline/profile_alignment_test.go`
   only to expect `RemovedDecisions` `[auth.provider auth.base_path]`, a
   declared break.
5. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`;
   a second refresh MUST report `File changes: 0`. The Result MUST name every
   file the commands rewrote, which the prototype measured as the derived
   files listed under task_01 in `_techspec.md` → Derived files. The Standard
   TypeScript golden `backend.md` MUST stay byte-identical.
6. MUST append the task_01 sentence of `_techspec.md` → Skill and document
   texts to `.agents/skills/setup-context-driven/SKILL.md`, then run
   `make skills-sync`,
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and `make skills-sync` again. The existing proposal sentence stays, because
   a public-docs contract pins it.
7. MUST add the task_01 paragraph to
   `docs/user-guide/context-driven-development.md` and re-record
   `docs/references/behavior-surfaces.json` with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
8. MUST add **Better Auth Base Path** to `CONTEXT.md` through
   `domain-modeling` with the text of `_techspec.md` → Skill and document
   texts.
9. MUST prove each new test can fail: the Result records one sabotage of the
   scope derivation, for example leaving the stale exception, and the test
   that failed.
10. MUST NOT change `routeException`'s owner, methods or reason, the
    `http.contract` decision, any module, the Source Baseline assets or any
    file outside its Context and the derived files.
11. MUST keep the interactive Baseline prompt from recording
    `auth.base_path` unasked (addendum 2026-10-12, after Run
    run_20261012T191532Z_192c8d919680b754 failed `make verify-changed` twice
    in `internal/cli`: `TestProjectDecisionReuse`,
    `TestProjectDecisionParity`,
    `TestBetterAuthSuggestionReusesFullHTTPException` and
    `TestProjectDecisionJourney` ended their scripted input because
    `promptBaselineDecisions` in `internal/cli/baseline_human.go` asks every
    profile decision). The prompt MUST skip `auth.base_path` while the
    current decisions do not record it, so the answer stays absent and the
    recorded provider exception stands, as Requirement 2 and ADR-0263
    require; a repository records it with `baseline update --decision`
    (task_02). When it is recorded, the prompt asks it with the recorded
    value as the suggestion. The scripted answers in
    `internal/cli/baseline_human_test.go` MUST stay unchanged, and the four
    tests above MUST pass.

## Subtasks

- [ ] Add the decision, the profile entries and the template token.
- [ ] Derive the provider exception and render the recorded path.
- [ ] Add the tests and update the divergence probe.
- [ ] Regenerate, refresh twice and record the skill and surface records.
- [ ] Update the setup skill, the user guide and the glossary.

## Acceptance Criteria

- [ ] A plan with `auth.base_path` `/auth` renders `/auth/*` in both backend
      lines, the recorded-path sentence and one Better Auth exception.
- [ ] A plan without it renders the backend guide as before and records no
      `auth.base_path`.
- [ ] Each invalid value of the Testing Approach is refused.
- [ ] This repository's refresh is `current` after the Task.
- [ ] `CONTEXT.md` defines **Better Auth Base Path**.

## Context

- instruction: `docs/adr/0263-the-better-auth-base-path-is-a-repository-decision-a-baseline-update-records.md`
- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- instruction: `internal/baseline/plan_test.go`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`
- interface: `internal/baseline/assets/templates/guides/backend.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/project_decisions.go`
- interface: `internal/baseline/project_decision_render.go`
- interface: `internal/baseline/profile_alignment_test.go`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `internal/cli/baseline_human.go`
- instruction: `internal/cli/baseline_human_test.go`
- creates: `internal/baseline/better_auth_base_path_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheBetterAuthBasePathIsAnOptionalDecision|TestARecordedBetterAuthBasePathMovesTheProviderException|TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope|TestAnInvalidBetterAuthBasePathIsRefused|TestDivergenceRendersProbe|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBetterAuthBasePathIsAnOptionalDecision TestARecordedBetterAuthBasePathMovesTheProviderException TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope TestAnInvalidBetterAuthBasePathIsRefused TestDivergenceRendersProbe TestCatalogCompatibility TestFormatterComposition TestBaselinePlanCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four base-path tests do not exist, so the command fails.
- `tr -s '[:space:]' ' ' < .agents/skills/setup-context-driven/SKILL.md | grep -qF -- "is Better Auth's default base path, not a Baseline rule" || { printf 'the setup skill lacks the base path sentence\n' >&2; exit 1; }; cmp -s .agents/skills/setup-context-driven/SKILL.md skills/setup-context-driven/SKILL.md || { printf 'the setup skill mirror differs\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "The Better Auth base path is the optional" || { printf 'the user guide lacks the base path paragraph\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Better Auth Base Path**:" || { printf 'CONTEXT.md lacks **Better Auth Base Path**\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; go test -tags docscontract -count=1 -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the setup skill lacks the sentence, so the command fails at its first check.

## References

- `_prd.md` → Goal 1; User Stories 1-2; Core Features 1 and 5; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 1; API Contract 2; Data Models; Skill and document texts; Derived files; Invariants 1, 2, 7 and 8; Testing Approach; Build Order 1
- ADR-0263; ADR-0076; ADR-0144; ADR-0244

## Result

Implemented the task_01 slice for Daemon Verification. Task status and the
Subtasks/Acceptance Criteria checkboxes remain Daemon-owned; no declared
Verification command, repository Verification, commit, push or PR was run.
The only pre-existing worktree change was this Task's `pending` to
`in_progress` status edit, which was preserved.

### Implementation

- Added optional `auth.base_path` after `auth.provider` in the decision catalog
  and both TypeScript profiles. Backend template version 4 declares the token
  and ends with `{{auth.provider}}{{auth.base_path}}`.
- Project normalization validates the recorded path, requires the provider,
  removes only the HTTP exception equal to the recorded provider exception,
  changes its scope to `<path>/*`, and invokes the existing HTTP derivation.
  Owner, methods, reason, unrelated exceptions, HTTP mode and source are kept.
  An absent decision leaves the existing normalization unchanged.
- The token renders nothing while absent and the API Contract 2 sentence
  while recorded. Tests exercise the embedded catalog and temporary repository
  Plans, including Setup Manifest values and repeated normalization.
- `promptBaselineDecisions` skips an unrecorded base path; a recorded one uses
  the existing recorded-value prompt. `baseline_human_test.go` and its scripted
  answers remain byte-identical to HEAD.
- Added the four specified tests and changed only the declared
  `TestDivergenceRendersProbe` expectation to
  `[auth.provider auth.base_path]`.
- Added the exact setup-skill sentence, user-guide paragraph and **Better Auth
  Base Path** glossary entry with its ADR-0263 link. The existing setup-skill
  proposal sentence stays. The owned-skill record command raised both copies
  from `0.0.4` to `0.0.5`.

### Acceptance evidence

| Acceptance criterion | Fresh focused evidence |
| --- | --- |
| Recorded `/auth` renders both `/auth/*` lines, the sentence and one moved Better Auth exception | `TestARecordedBetterAuthBasePathMovesTheProviderException` passed for former `/api/auth/*` and `/v1/auth/*` scopes, including custom provider methods/reason, retained Operations exception/source, Plan and Manifest values, and replan stability. |
| Absent decision preserves the guide and records no answer | `TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope` passed for `/api/auth/*` and `/v1/auth/*`; its default case compares the rendered backend bytes with the existing golden. A separate byte comparison against HEAD confirmed the golden itself unchanged. |
| Invalid values are refused | `TestAnInvalidBetterAuthBasePathIsRefused` passed for all seven authored invalid strings, fragment, leading/Unicode whitespace, backtick, control, marker/token text, non-string input and missing provider. String-path cases exercise `BuildPlan`; missing-provider and non-string cases exercise project normalization directly because Plan collection first asks for required answers and validates scalar types. The authored normalization messages are asserted. |
| This repository's refresh is `current` | First approved refresh reported `verified` and one changed file; second approved refresh reported `File changes: 0` and verified idempotence. A subsequent text preview reported `Baseline update: current`, `File changes: 0`. |
| `CONTEXT.md` defines **Better Auth Base Path** | The exact glossary definition is present beside HTTP Contract Decision, with the accepted ADR link; inspected after the last glossary edit. |

### Focused checks and regeneration

- Before implementation, `go test ./internal/baseline -run 'BetterAuthBasePath'
  -count=1` exited 1: the catalog lacked the decision, recorded paths were
  unknown, and invalid-input tests did not receive their authored messages.
- After the final test edit,
  `go test ./internal/baseline -run 'BetterAuthBasePath|^TestHTTPContractConflict$|^TestProjectDecisionRendering$|^TestBetterAuthGuidance$|^TestDivergenceRendersProbe$' -count=1 -v`
  exited 0; all eight selected tests passed. Existing conflict detection and
  structured-render golden compatibility passed alongside the new behavior.
- `go test ./internal/cli -run '^(TestProjectDecisionReuse|TestProjectDecisionParity|TestBetterAuthSuggestionReusesFullHTTPException|TestProjectDecisionJourney)$' -count=1`
  exited 0 with all four required prompt flows selected.
- A temporary Go test overlay ran
  `TestATask01RecordedBasePathPromptProbe` in `internal/cli`, exit 0. It
  observed `Keep auth.base_path="/auth"`, retained the answer and derived one
  `/auth/*` exception with no former `/api/auth/*` exception. The overlay and
  log are outside the repository under
  `roundfix-task01-recorded-prompt-nk_0stn0` in the host temporary directory;
  no prompt test source was changed.
- Ran `make skills-sync`,
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  and `make skills-sync` again, all exit 0. They rewrote the canonical
  `.agents/skills/setup-context-driven/SKILL.md`, its mirror
  `skills/setup-context-driven/SKILL.md`, and
  `skills/testdata/owned-skill-versions.json`. The two skill copies match.
- Ran `make baseline-digests`, exit 0, including its sanctioned formatter,
  characterization and compatibility regeneration. It named exactly these
  rewritten files:
  - `internal/baseline/testdata/catalog.digest`
  - `internal/baseline/testdata/catalog.normalized.json`
  - `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
  - `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
  - `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
  - `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- Ran
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
  twice, both exit 0. Only the first invocation rewrote
  `docs/agents/setup-context.json`; the second reported no change. The
  Standard TypeScript golden backend and profile digest pins were unchanged.
- Ran
  `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
  exit 0; it rewrote only `docs/references/behavior-surfaces.json`.
- Ran the text preview
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format text`,
  exit 0, state `current` with no file changes.

### Sabotage evidence

Each mutation used `go test -overlay <temporary-overlay.json>
./internal/baseline -run '^<test-name>$' -count=1`; repository files were
unchanged. Each intended assertion failed with exit 1:

| Sabotage | Test that detected it |
| --- | --- |
| Set the catalog decision's `optional` field to false | `TestTheBetterAuthBasePathIsAnOptionalDecision`: expected optional string declaration. |
| Leave the former matching HTTP exception instead of removing it | `TestARecordedBetterAuthBasePathMovesTheProviderException`: HTTP contract retained both old and moved provider exceptions; backend also retained the stale scope. |
| Replace an unrecorded provider scope with `/wrong/*` | `TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope`: recorded provider scope changed. |
| Remove the leading-slash validation | `TestAnInvalidBetterAuthBasePathIsRefused/relative`: `auth` incorrectly produced a Plan without an error. |

Logs for the last three probes are in temporary directory
`roundfix-task01-sabotage-24qb_hov`; the corrected optional-field probe is in
`roundfix-task01-optional-sabotage-8ulsv264`. The initial optional probe used
colliding source/overlay filenames and failed catalog parsing; it is excluded
from the behavioral evidence above. All final focused baseline checks ran
without overlays and passed.

### Scope postflight

`git diff --check` exited 0. Comparing every tracked and untracked changed
path against this Task's Context found 24 paths and no out-of-scope path.
The Standard TypeScript golden `backend.md` and
`internal/cli/baseline_human_test.go` were compared byte-for-byte with HEAD
and are unchanged. No module, `http.contract` declaration, Source Baseline
asset, Task Graph manifest or other Task file was changed.

Follow-up: task_02 owns the public `baseline update --decision` input. This
Task adds no update flag or declaration API. Daemon Verification and settlement
remain pending.

### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T194228Z_1d1cb31ee8783296/verification/batch-001-attempt-1.log`.
The Daemon's `make verify-changed` attempt exited 2 in the core test stage;
its reported failure is
`internal/worktree.TestAddDetachedWorktreeWaitsForTheAdminLock`. This is a
Daemon-recorded Verification failure, not an Agent rerun or a Task verdict.

Investigation traced the test's 250 ms deadline through
`AddDetachedWorktree` (`internal/worktree/disposal.go`) to
`withWorktreeAdminLock` and `resolveGitCommonDir`
(`internal/worktree/adminlock.go`). The reported failure occurs during the
Git common-directory lookup, before waiting for the administration lock.
`execGitRunner.Run` (`internal/worktree/worktree.go`) wraps the subprocess
error without preserving `ctx.Err()` when a running command is terminated.
Consequently the test's expected context-deadline identity is not guaranteed
on that path. The cleanup failure follows test-context cancellation after
its fatal assertion. The original reason the lookup exceeded its deadline
is not established by the artifact; host scheduling/load is a hypothesis.

Ran exactly one focused diagnostic reproduction:
`go test ./internal/worktree -run '^TestAddDetachedWorktreeWaitsForTheAdminLock$' -count=1 -v`.
It exited 0 and reported the named test passed. This isolated result does not
resolve the observed Verification failure or establish that the full gate
passes; no repeated-until-green check was performed.

`git diff HEAD -- internal/worktree` is empty. The traced path calls no
Baseline decision/rendering code or interactive Baseline prompt code. No
repair in task_01's authorized files addresses this error. All affected
worktree source and test files are outside this Task's Context and the
sanctioned derived files, so they were left unchanged. Follow-up ownership:
a separately authorized worktree cancellation repair should preserve context
error identity during Git lookup and cover that boundary deterministically,
without increasing the timeout or weakening the existing assertion.

Only this Result evidence was updated during the feedback turn. Task status,
authored Verification and every other task file remain untouched. Neither
`make verify-changed` nor any declared Verification command was rerun by the
Agent. The Daemon owns the next full Verification attempt and settlement.
