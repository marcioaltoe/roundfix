---
task: task_05
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: pending
type: backend
complexity: medium
---

# Task 05: Refuse a null provider and tell the skill what Doctor reports

## Overview

Independent review found that `pre_pr_review.provider:` written with no value
leaves the decoded pointer nil, so validation is skipped and the layer below is
inherited in silence — the opposite of a declared decision. It also found the
shipped Roundfix skill silent about the new Doctor check, which the repository's
hard rule does not allow for a change in CLI behavior. This slice closes both,
and it is the first of the two corrective Tasks the contract allows.

## Requirements

1. MUST refuse `pre_pr_review.provider` written as null, as an empty value, or
   as a non-scalar node, with the error that names the key, the value and the
   four supported values.
2. MUST keep an absent `pre_pr_review` section inheriting, since silence and an
   empty written value are different statements.
3. MUST describe, in the Roundfix skill and its regenerated mirror, that the
   Doctor Command reports the resolved provider and the configuration layer that
   supplied it, that `none` reads as review disabled by configuration, and that
   the check invokes no provider.
4. MUST change no path outside the two bounded in this Spec's approved authority
   record, plus its own Task file and ordinary source.
5. MUST name the decoding type `prePRReviewProviderValue`, so the Task's own
   check and later readers agree.

## Subtasks

- [ ] Decode the provider through a value type that refuses null and non-scalar
      nodes, as the repository already does for another boolean key.
- [ ] Cover null, empty, non-scalar, absent-section and each supported value.
- [ ] Write the Doctor check into the Roundfix skill and regenerate the mirror.

## Acceptance Criteria

- [ ] `provider:` with no value fails the load with the named error.
- [ ] A non-scalar value fails the same way.
- [ ] An absent section still inherits and resolves to the layer below.
- [ ] The Roundfix skill and its mirror describe the check, and the two copies
      are identical.

## Context

- interface: `internal/config/config.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "prePRReviewProviderValue" internal/config/config.go` — expected: exit 0; the provider decodes through a value type that can refuse. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestPrePRReviewProviderRefusesNullValue$" ./internal/config 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `grep -q "pre-pr-review" .agents/skills/roundfix/SKILL.md` — expected: exit 0; the canonical skill describes the check. Before this Task it does not.
- `grep -q "pre-pr-review" skills/roundfix/SKILL.md && diff -q .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; the mirror carries the same text and stays identical to its canonical copy.

## References

`_prd.md` → Core Features 3-4; User Story 5; Goal 3;
Project Constraints: Tooling authority;
`_techspec.md` → Implementation Design: Refusal at load, The report;
API Contracts 1-2; `_authorization.md`.
