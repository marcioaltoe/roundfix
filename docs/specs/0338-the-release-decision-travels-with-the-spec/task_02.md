---
task: task_02
spec: 0338-the-release-decision-travels-with-the-spec
status: pending
type: backend
complexity: high
---

# Task 02: An authorization record's Release Ceiling answers the release plan's approval

## Overview

An approved `_authorization.md` that grants `release` may carry
`release_ceiling` and `release_paths`. The release plan, for a plan that would
ask, finds each Spec the range delivers, reads its record at the parent of the
delivering commit, and answers `approval_required` when every delivered Spec
grants a ceiling at or above the impact, nothing is breaking or classified by
hand, and every impact-setting change is such a delivery. Backlog Entries
become planning evidence, so plan commits stop forcing manual classification.
This answers the adopted Backlog Entry of 2026-11-01, "A Spec's authorization
carries the release decision"
([references](references/2026-11-01-a-spec-authorization-carries-the-release-decision.md)),
the Argus request of 2026-11-01.

## Requirements

1. MUST add `ReleaseCeiling` and `ReleasePaths` to `AuthorizationRecord` and
   read and classify `release_ceiling` and `release_paths` in
   `parseAuthorizationFrontmatter` and `classifyAuthorizationRecord`
   (`internal/authorization/authorization.go`) exactly as API Contract 5 of
   `_techspec.md` states, reusing the `executables` path rules for
   `release_paths`. A record without either key must classify as before.
2. MUST create `internal/authorization/release_ceiling_test.go` with
   `TestReleaseCeilingIsReadAndValidated` covering the granted pair, the four
   refusals with their code and field, and a record without the keys.
3. MUST create `internal/releaseplan/release_grant.go` with `ReleaseGrant`,
   `ReleaseGrantSource`, `ReleaseAuthorization` and `AnswerApproval`, add
   `ReleaseAuthorization` to `Plan` (`internal/releaseplan/model.go`), and
   make `Build` (`internal/releaseplan/build.go`) call `ReleaseGrants` only
   for an `approval_required` proposal whose source implements
   `ReleaseGrantSource` (API Contract 7, its five conditions, order and
   reason texts).
4. MUST add `docs/backlog/` to `isPlanningEvidencePath`
   (`internal/releaseplan/classifier.go`) and a
   `docs/backlog/2026-11-01-example.md` case with `want: true` to
   `TestMaintenanceOnly` (`internal/releaseplan/classifier_test.go`), a
   declared break.
5. MUST create `internal/releaseplan/release_grant_test.go` with
   `TestAnswerApprovalFollowsEveryCeiling`, one subtest per condition of API
   Contract 7 asserting its reason, the answered case asserting `ready`, the
   unchanged proposed version and an approval without a question, and the
   nil case for grants without a ceiling.
6. MUST implement `ReleaseGrants` on `releasePlanGitSource`
   (`internal/cli/releaseplan_git_source.go`) as API Contract 6 states,
   reading each record with `spec.ReadAuthorization` at `<sha>^`, and render
   `ReleaseAuthorization` in `printReleasePlanText` and
   `releasePlanJSONFromPlan` (`internal/cli/releaseplan_command.go`) as API
   Contract 8 states. A plan whose `ReleaseAuthorization` is nil keeps every
   byte of text and JSON.
7. MUST add the API Contract 10 paragraph to `commandUsage("release plan")`
   in `internal/cli/cli.go`, changing no other help line.
8. MUST create `internal/cli/releaseplan_release_ceiling_test.go` with
   `TestReleasePlanAnswersApprovalFromTheReleaseCeiling` and the six subtests
   `_techspec.md` → Testing Approach lists ("ceiling minor", "ceiling patch",
   "no ceiling", "a direct feat commit", "a breaking footer", "classified by
   hand"), each building its temporary repository with the existing release
   plan test helpers, running the command in text and JSON, and asserting the
   exit code, the `Decision:` line, the `Release authorization:` line and the
   JSON `releaseAuthorization`; "no ceiling" compares the full text and JSON
   byte for byte with the same repository whose record omits `release`. Every
   top-level test and subtest calls `t.Parallel()` first.
9. MUST update `docs/user-guide/release-runbook.md`,
   `docs/user-guide/usage.md`, the canonical
   `.agents/skills/roundfix/references/release.md` and the `## Release
   planning` section of `.agents/skills/roundfix/SKILL.md` so each keeps its
   existing sentences, the phrase "generic release request authorizes only a
   conclusive patch plan" included, and adds that a Release Ceiling recorded
   in the authorization record of every Spec the range delivers answers
   `approval_required` within the ceiling, with the sentence "A Release
   Ceiling never answers a breaking plan or a manual classification", the
   `Release authorization:` lines and `releaseAuthorization`, Backlog Entries
   as planning evidence, and that a `refactor:` delivery is a commit to
   classify by hand. The runbook and the reference also name
   `release_ceiling` and `release_paths`.
10. MUST add **Release Ceiling** to `CONTEXT.md` through `domain-modeling`:
    the highest impact (`patch`, `minor` or `major`) an approved Spec
    authorization record that grants `release` answers in the release plan
    without a second approval, with its release paths, linked to ADR-0349;
    and add one sentence to **Release Plan** that the records of the Specs
    a range delivers can answer its approval within their Release Ceilings.
11. MUST copy the canonical skill with `make skills-sync`, raise both
    version fields of the Roundfix `SKILL.md` pair with
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
    -record-skill-versions -count=1`, and re-record the Behavior Surface
    Record as the LAST step, after the last edit of any guide or help text,
    with `go test -count=1 -tags docscontract ./internal/docscontract -run
    '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
12. MUST keep ADR-0349
    (`docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`,
    committed by the planning change) as written; change it only if the
    implementation proves a sentence wrong, and then say so in the Result.
13. MUST NOT touch any Baseline module, `docs/agents/`, the `write-*`
    skills, the release plan's reset mode, `follow_up_status.go`, any
    `### QA settlement` section, the Makefile, CI workflows, `go.mod` or
    `go.sum`.

## Subtasks

- [ ] Characterize: a minor range that delivers one Spec reports `approval_required` today, and a plan commit that moves a Backlog Entry is a blocking commit.
- [ ] Add the record keys and their test.
- [ ] Add the grants, the answer, the planning path and their tests.
- [ ] Add the Git lookup, the rendering, the help and the CLI test.
- [ ] Update the guides, references, `SKILL.md` and `CONTEXT.md`; sync, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] A range delivering one Spec whose record grants ceiling `minor` exits 0 with `Decision: ready` and the record named; every other case of `_techspec.md` → Testing Approach exits 3 with its reason.
- [ ] A range whose records carry no ceiling keeps its text and JSON bytes.
- [ ] The guides, references, `SKILL.md` and `CONTEXT.md` state the rule, the mirrors equal their canonical copies, and the skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0143-release-planning-is-read-only-including-its-reset-mode.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `docs/adr/0340-an-untracked-executable-rides-a-task-commit-only-when-the-spec-authorization-names-it.md`
- instruction: `docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/releaseplan_command_test.go`
- instruction: `internal/cli/releaseplan_git_source_test.go`
- instruction: `internal/cli/follow_up_status.go`
- interface: `internal/authorization/authorization.go`
- creates: `internal/authorization/release_ceiling_test.go`
- interface: `internal/releaseplan/model.go`
- interface: `internal/releaseplan/build.go`
- interface: `internal/releaseplan/classifier.go`
- interface: `internal/releaseplan/classifier_test.go`
- creates: `internal/releaseplan/release_grant.go`
- creates: `internal/releaseplan/release_grant_test.go`
- interface: `internal/cli/releaseplan_git_source.go`
- interface: `internal/cli/releaseplan_command.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/releaseplan_release_ceiling_test.go`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `AuthorizationRecord.ReleaseCeiling` and `ReleasePaths` — consumers: `./internal/authorization`, `./internal/spec`, `./internal/cli`
- `releaseplan.Build`, `AnswerApproval` and `MaintenanceOnly` — consumers: `./internal/releaseplan`, `./internal/cli`
- The release plan text, JSON and help — consumers: `./internal/cli`, `./internal/docscontract`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestReleasePlanAnswersApprovalFromTheReleaseCeiling$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in 'ceiling_minor' 'ceiling_patch' 'no_ceiling' 'a_direct_feat_commit' 'a_breaking_footer' 'classified_by_hand'; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestReleasePlanAnswersApprovalFromTheReleaseCeiling/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestAnswerApprovalFollowsEveryCeiling|TestMaintenanceOnly)$' ./internal/releaseplan 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAnswerApprovalFollowsEveryCeiling ' || { echo 'missing pass: TestAnswerApprovalFollowsEveryCeiling' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^TestReleaseCeilingIsReadAndValidated$' ./internal/authorization 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestReleaseCeilingIsReadAndValidated ' || { echo 'missing pass: TestReleaseCeilingIsReadAndValidated' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/releaseplan ./internal/authorization ./internal/spec && go test -count=1 -timeout 20m -run 'ReleasePlan|FollowUp' ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the first test run reports no pass line; afterwards a sufficient ceiling answers the plan, every other case asks with its reason, and the release plan, follow-up, authorization and Spec suites pass.
- `for f in docs/user-guide/release-runbook.md docs/user-guide/usage.md .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Release Ceiling never answers a breaking plan or a manual classification' || { printf 'missing ceiling rule in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- 'release_ceiling' docs/user-guide/release-runbook.md || { echo 'runbook does not name release_ceiling' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Release Ceiling**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; cmp .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide names a Release Ceiling, so the first check fails; afterwards the guides, references, `SKILL.md` pair and glossary carry the rule, the mirrors equal their canonical copies, and the skill version and documentation contracts, the release-plan documentation snippets among them, pass.

## References

- `_prd.md` → Sources 2; Goals; Core Features 3-5; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 8; Glossary; Decisions
- `_techspec.md` → Current behavior 4-6; Interfaces; Data Models; API Contracts 5-10 and 12; Invariants 3-6; Declared breaks; Testing Approach; Build Order 2
- ADR-0349; ADR-0143; ADR-0256; ADR-0340
