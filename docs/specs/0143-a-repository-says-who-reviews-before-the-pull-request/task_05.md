---
task: task_05
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: completed
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

## Result

### Implementation

- Added the required `prePRReviewProviderValue` decoder. It accepts only scalar
  supported providers, preserves key absence through its pointer, and renders
  non-scalar YAML in the existing named invalid-provider error.
- Added an explicit parsed-node refusal for null because `yaml.v3` otherwise
  skips a value unmarshaler for null nodes. This distinguishes a written empty
  or null provider from an absent `pre_pr_review` section.
- Added configuration tests for a valueless provider, explicit null, an empty
  string, sequence and mapping values. Existing focused cases continue to
  cover every supported provider and inheritance from an absent layer.
- Documented the Doctor `pre-pr-review:` check in the canonical Roundfix skill:
  it reports the resolved provider and source layer, calls explicit `none`
  disabled by configuration, invokes no provider and mutates nothing. Regenerated
  the distributed skill mirror from the canonical copy.

### Focused checks

- Initial `GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1
  ./internal/config` exposed that `yaml.v3` bypassed the value unmarshaler for
  null nodes; the empty-value and explicit-null cases failed while non-scalar
  cases were already refused. The parsed-node refusal fixes that root cause.
- `GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1 -v -run
  '^TestPrePRReview(PolicyResolution|ProviderRefusesNullValue)$'
  ./internal/config`: passed all 12 subtests. Evidence includes valueless,
  explicit-null, empty-string, sequence and mapping refusals; all four supported
  providers; default inheritance; and User Config inheritance when Project
  Config has no section.
- `make skills-sync`: exited 0 and regenerated `skills/roundfix/SKILL.md`.
- `make baseline-digests`: exited 0 and reported that derived artifacts already
  matched their canonical sources, with no additional changed path.
- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`: exited 0;
  the canonical and distributed skill copies are byte-identical.
- `git diff --check`: exited 0.

### Acceptance evidence

1. `TestPrePRReviewProviderRefusesNullValue/empty_value` and `/explicit_null`
   observe load refusal with `pre_pr_review.provider`, the written value and all
   four supported values in the error.
2. The same test's `/sequence` and `/mapping` cases observe the same named error
   and include the rendered offending value.
3. `TestPrePRReviewPolicyResolution/user_selection_overrides_default` proves an
   absent Project Config section inherits User Config; `/both_layers_silent_use_default`
   proves complete silence resolves the built-in default.
4. The skill text states the resolved provider, source layer, disabled meaning
   of `none`, and absence of provider calls; `cmp -s` proves the mirror is
   identical.
