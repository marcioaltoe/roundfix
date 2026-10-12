---
task: task_01
spec: 0257-the-baseline-respects-the-repository
status: pending
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
