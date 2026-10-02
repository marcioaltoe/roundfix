---
task: task_01
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: completed
type: backend
complexity: high
---

# Task 01: The pre-PR review omits QA evidence and upstream-managed skills and blocks above a measured bound

## Overview

The pre-PR review sends the candidate's whole `git diff` as its prompt, and
Spec 0200's diff of 1,295,055 bytes, 61% of it QA evidence, failed with
`agent/protocol error` three times. This Task omits QA evidence and
upstream-managed skill copies from every diff the review sends, lists them in
the prompt and the record, blocks with a named cause before any provider call
when the remainder exceeds the measured bound, records the runtime's stderr
tail on a runtime failure, and describes the change in the Roundfix Skill's
`review` reference and the `review` command guide.

## Requirements

1. MUST compute every diff the review sends, round one's and round two's,
   with the omission of "The reviewed diff": paths under any Spec's
   `qa/evidence/` in the resolved Specs Root and its archive root, reason
   `qa-evidence`; paths under `.agents/skills/<name>/` for each name in the
   `skills` map of `skills-lock.json` at the base or the head, reason
   `upstream-skill`. The QA Report MUST stay in the diff. A missing lock MUST
   omit nothing for that class; an unreadable or malformed lock MUST block
   with `review scope: read skills-lock.json: <error>`.
2. MUST list the omitted paths in the prompt after the candidate diff block
   with the line of "The reviewed diff", and MUST record them in
   `omittedPaths`, sorted, with `diffBytes`, as API Contract 1 states.
3. MUST add `reviewDiffBound` at 917504 bytes, citing its measurement in a
   comment, and MUST record `blocked` with API Contract 2's reason before any
   Agent Session is prepared when `diffBytes` exceeds it, without trying the
   configured fallback.
4. MUST record `runtimeStderrTail`, the last 10 lines capped at 1,024 bytes of
   an `agent.BatchFailureError`'s `Stderr` reached through `errors.As`, on a
   runtime failure, leaving the reason text unchanged.
5. MUST keep reading a record without the new fields.
6. MUST add the tests named in Verification to the new file
   `internal/cli/review_scope_test.go`, over temporary repositories and the
   `reviewCommandRunner` fake, and MUST leave every existing review test
   unedited and passing, except `TestReviewBoundsSpecContext` in
   `internal/cli/review_test.go`: its oversized fixture writes about 18 times
   the per-Spec context limit into each of two files, so its diff (1,180,372
   bytes) exceeds the new bound. Shrink only that fixture to
   `strings.Repeat("candidate context ", reviewSpecContextPerSpecLimit/8)`, which
   still exceeds the per-Spec limit and keeps every assertion of the test.
7. MUST describe the omission, its reasons, the bound and its reason text,
   and the three record fields in `.agents/skills/roundfix/references/review.md`
   and `docs/user-guide/commands/review.md`; MUST raise the Roundfix Skill's
   version by one patch level from the value on this Task's base in both
   front-matter fields, run `make skills-sync`, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.

## Subtasks

- [ ] Classify and omit evidence and upstream-managed skill paths in both review diffs.
- [ ] List omitted paths in the prompt and the record, with the diff size.
- [ ] Block above the bound before any provider call.
- [ ] Record the runtime's stderr tail.
- [ ] Add the scope tests.
- [ ] Describe the change in the skill reference and the command guide, and raise the version.

## Acceptance Criteria

- [ ] A candidate holding QA evidence, a lock-named skill file, a skill the
      lock dropped at the head, the QA Report and a Go source sends a diff
      without the first three, lists them sorted with their reasons in the
      prompt and the record, and keeps the QA Report and the Go source.
- [ ] A malformed lock blocks with the scope reason.
- [ ] A diff above the bound is recorded `blocked` with API Contract 2's
      reason and the fake runner records no prompt.
- [ ] A runtime failure carrying 30 lines of stderr records at most the last
      10 lines and 1,024 bytes.
- [ ] A record written before this Task still reads.
- [ ] The review reference and guide carry the new words, the mirrors equal
      their canonical files, and the raised version is recorded.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_scope_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/review.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/review.md`
- interface: `internal/cli/review_test.go`
- instruction: `internal/cli/review_session_test.go`
- instruction: `internal/agent/acpx_runner.go`
- instruction: `skills-lock.json`
- instruction: `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/references/2026-10-01-the-pre-pr-review-fails-on-a-large-vendored-diff.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewOmitsQAEvidenceAndUpstreamSkills|TestReviewOmitsASkillTheLockDroppedAtTheHead|TestReviewKeepsTheQAReportInTheDiff|TestReviewBlocksOnAMalformedSkillsLock|TestReviewBlocksAboveTheBoundWithoutAProviderCall|TestReviewRecordsTheRuntimeStderrTail|TestReviewReadsARecordWithoutTheNewFields)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewOmitsQAEvidenceAndUpstreamSkills TestReviewOmitsASkillTheLockDroppedAtTheHead TestReviewKeepsTheQAReportInTheDiff TestReviewBlocksOnAMalformedSkillsLock TestReviewBlocksAboveTheBoundWithoutAProviderCall TestReviewRecordsTheRuntimeStderrTail TestReviewReadsARecordWithoutTheNewFields; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestReviewBoundsSpecContext|TestReviewPromptCarriesTheCandidateDiff|TestReviewCommandExitsZeroOnExplicitClean|TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt|TestReviewOmitsQAEvidenceAndUpstreamSkills)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewBoundsSpecContext TestReviewPromptCarriesTheCandidateDiff TestReviewCommandExitsZeroOnExplicitClean TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt TestReviewOmitsQAEvidenceAndUpstreamSkills; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing review tests run unedited beside the scope test, which does not exist before this Task.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "qa-evidence" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/review.md "qa-evidence" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "upstream-skill" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/review.md "upstream-skill" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "review diff too large" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/review.md "review diff too large" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "omittedPaths" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/review.md "omittedPaths" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- "omittedPaths" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/review.md "omittedPaths" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- "review diff too large" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/review.md "review diff too large" >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the review reference and guide carry none of these words, so the command fails.

## References

- `_prd.md` → User Stories 1-2; Core Features 1-3; Core Feature 8; Success Metrics 1-2; Acceptance evidence
- `_techspec.md` → The reviewed diff; The bound; The stderr tail; Interfaces; API Contract 1; API Contract 2; Measured outside evidence; Testing Approach 1; Build Order 1
- ADR-0153; ADR-0169; ADR-0174; ADR-0196; ADR-0197; ADR-0187; ADR-0189

## Result

Implemented the Task 01 review slice. Task status and declared Verification
remain Daemon-owned; no commit, push or Pull Request was performed.

### Implementation and acceptance evidence

1. Both review ranges use the same changed-path classification and literal
   Git exclusions. QA evidence under active and archived resolved Spec roots
   is omitted as `qa-evidence`; skill names are unioned from the base and
   head locks and omitted as `upstream-skill`. Reports and source stay in
   scope. The prompt lists sorted omissions immediately after its diff, and
   records always write `diffBytes` and an array `omittedPaths`.
   Focused evidence: `TestReviewOmitsQAEvidenceAndUpstreamSkills` includes
   omitted payloads larger than the bound and retained Go source;
   `TestReviewOmitsASkillTheLockDroppedAtTheHead` proves base-lock ownership;
   `TestReviewKeepsTheQAReportInTheDiff` proves report retention and missing
   lock behavior. `TestReviewScopesRoundTwoDelta` and
   `TestReviewScopesConfiguredSpecRoots` cover the second range and custom
   active/archive roots. These tests assert prompt and persisted record
   scope, size and ordering.
2. Malformed and unreadable locks block with the exact scope prefix before
   provider probes, preparation or prompts. Focused evidence:
   `TestReviewBlocksOnAMalformedSkillsLock`,
   `TestReviewBlocksOnAMalformedBaseLock` and
   `TestReviewBlocksOnAnUnreadableSkillsLock`.
3. The measurement-cited 917504-byte bound runs before readiness probes and
   session preparation. Oversized round-one and round-two diffs retain their
   measured size and block with the required reason without fallback.
   Focused evidence: `TestReviewBlocksAboveTheBoundWithoutAProviderCall`
   asserts unchanged provider call counts with a configured fallback;
   `TestReviewAdmitsADiffExactlyAtTheBound` proves the inclusive boundary.
4. Wrapped `agent.BatchFailureError` stderr is reached through `errors.As`
   and recorded as at most the last 10 lines and 1,024 bytes. Existing
   runtime reason text is preserved, including the readiness reason prefix.
   Focused evidence: `TestReviewRecordsTheRuntimeStderrTail` covers both
   line-limited and byte-limited 30-line failures and persisted diagnostics.
5. Older records remain readable and normalize an absent omission list to
   an empty array on writing. Focused evidence:
   `TestReviewReadsARecordWithoutTheNewFields` uses pre-change JSON.
6. The review skill reference and command guide describe omission reasons,
   the bound and exact reason, and all three new record fields. Both skill
   version fields increased from the Task base's 0.1.10 to 0.1.11; mirrors
   were regenerated and the new version was recorded. A Python inspection
   confirmed the required documentation phrases, canonical/mirror byte
   equality, base-relative patch increase and ledger entry. Skills-package
   checks passed. Existing review tests were left unchanged except the
   expressly required `/8` fixture reduction, confirmed by Git diff.

### Focused checks

- `GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/cli -run '^TestReview' -count=1`
  — exit 0 after the final implementation and test edits (`29.572s`).
- `GOCACHE=/private/tmp/roundfix-task01-gocache go test ./skills -count=1`
  — exit 0 (`1.299s`).
- `make skills-sync` — exit 0; skill mirrors regenerated.
- `GOCACHE=/private/tmp/roundfix-task01-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; 0.1.11 recorded.
- `GOCACHE=/private/tmp/roundfix-task01-gocache make baseline-digests`
  — exit 0; derived artifacts already matched, no changes.
- Documentation/mirror/version inspection and `git -c core.fsmonitor=false diff --check`
  — exit 0.
- An isolated Go overlay of the original `HEAD` implementation with the
  omission regression test exited 1 at `omitted payload reached reviewer`.
  The overlay changed no repository files and reached no real provider.

The declared `## Verification` commands were not run. No follow-up outside
this Task's slice was implemented.

## Carry-forward provenance

- Source Run: `run_20261002T100707Z_cbb39973a4b07b9b`
- Source commit: `e41af9618b3c1e4769b2b6135f5fb82e4d0b224e`
