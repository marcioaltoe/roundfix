---
task: task_01
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: completed
type: backend
complexity: medium
---

# Task 01: An unselected provider receives no request

## Overview

The Roundfix skill's section "A review only happens when it is asked for"
tells every agent that a pull request gets no review unless someone requests
one, and how to request CodeRabbit's, "across these repositories", whatever
Pre-PR Review Policy the repository selected. The guides call `coderabbit` "the
supported Review Source" without saying it is the legacy PR-feedback source,
and `defaultConfigYAML` in `internal/config/config.go` writes a
`review_source` block with no such label. The guidance is read by every agent
that loads the skill, and the generated config by every repository that runs
`roundfix init`. In code, the only CodeRabbit request Roundfix publishes is
the legacy `watch`/`resolve` review request, and `roundfix review` and the
Delivery Queue never read `.coderabbit.yaml`. This Task scopes the guidance
and the generated config, and pins the code paths with tests so a later
change cannot quietly add a request.

## Requirements

1. MUST rescope the Roundfix skill's section "A review only happens when it is
   asked for" in `.agents/skills/roundfix/SKILL.md` to a repository whose
   Pre-PR Review Policy selects `coderabbit` and to the legacy PR-feedback
   commands `fetch`, `watch` and `resolve`. The section MUST state that a
   repository whose policy selects `codex`, `claude` or `none` is never asked
   for a CodeRabbit review (the phrase `never asked for a CodeRabbit review`),
   MUST no longer carry the universal claim "across these repositories", and
   MUST keep every request instruction for the scoped case. Regenerate
   `skills/roundfix/SKILL.md` with `make skills-sync`.
2. MUST extend the sentence "The supported Review Source is `coderabbit`" in
   `docs/user-guide/commands.md`, the matching sentence in
   `docs/user-guide/usage.md`, and the `review_source.name` row of
   `docs/user-guide/configuration.md` to say that the Review Source is the
   legacy PR-feedback source, read only by `fetch`, `watch` and `resolve`, that
   never selects or requests a pre-PR reviewer (the phrase `never selects or
   requests a pre-PR reviewer` in each of the three guides).
3. MUST add, in `defaultConfigYAML` in `internal/config/config.go`, a comment
   line directly above `review_source:` stating that it is the legacy
   PR-feedback Review Source read only by fetch, watch and resolve, and never
   selects or requests a pre-PR reviewer (the phrase `read only by fetch, watch
   and resolve`). Mirror it in the YAML example of
   `docs/user-guide/configuration.md`. The generated config MUST keep
   `request_review: false` and MUST write no `pre_pr_review` key, so it opts the
   repository into neither a service nor `none`.
4. MUST add `internal/config/review_source_scope_test.go` with:
   - `TestDefaultConfigScopesReviewSourceToPullRequestFeedback`: both
     `DefaultConfigYAML()` and `DefaultProjectConfigYAML()` contain `read only
     by fetch, watch and resolve` and `request_review: false`.
   - `TestDefaultConfigSelectsNoPrePRReviewProvider`: neither contains
     `pre_pr_review`, and loading the generated project YAML as Project Config
     through `Load` resolves `PrePRReview.Provider` to `codex` with source
     `default`.
5. MUST add `internal/cli/review_provider_scope_test.go` with
   `TestReviewUnderCodexIgnoresCodeRabbitConfiguration` and
   `TestReviewUnderNoneIgnoresCodeRabbitConfiguration`. Each writes a
   `.coderabbit.yaml` with `reviews.auto_review.enabled: false` and a Project
   Config with `review_source.request_review: false`, the pair that
   `validateReviewRequestCoherence` refuses for `watch`. It then runs
   `roundfix review` through the existing review command fixture. Under
   `codex` with a `No findings` answer it asserts exit `0` and outcome
   `reviewed`; under `none` it asserts exit `0`, outcome `omitted` and zero
   Agent runner calls.
6. MUST add `internal/cli/deliver_publication_scope_test.go` with
   `TestDeliveryPublicationRequestsNoCodeRabbitReview`: for each policy
   `codex`, `claude`, `coderabbit` and `none`, the title and body returned by
   `commandDeliveryWorkflow.Publication` contain neither `@coderabbitai` nor
   `coderabbit:review`, compared case-insensitively.
7. MUST keep `TestRunEnforcesReviewRequestCoherence` in
   `internal/preflight/preflight_test.go` green and unchanged, as the mirror:
   when the legacy Review Source is exercised, CodeRabbit keeps its coherence
   refusal. MUST NOT change `validateReviewRequestCoherence`, the legacy review
   request publication, or any legacy command behavior.

## Subtasks

- [ ] Scope the skill section and the three guides.
- [ ] Label `review_source` in the generated config.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] The skill never tells an agent to request a CodeRabbit review for a
      repository whose policy selects `codex`, `claude` or `none`.
- [ ] `roundfix review` succeeds under `codex` and `none` beside an incoherent
      CodeRabbit configuration, and the delivery publication carries no request.
- [ ] The generated config labels `review_source` and selects no provider.
- [ ] The legacy coherence refusal still fires for `watch`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `internal/config/config.go`
- creates: `internal/config/review_source_scope_test.go`
- creates: `internal/cli/review_provider_scope_test.go`
- creates: `internal/cli/deliver_publication_scope_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDefaultConfigScopesReviewSourceToPullRequestFeedback|TestDefaultConfigSelectsNoPrePRReviewProvider|TestReviewUnderCodexIgnoresCodeRabbitConfiguration|TestReviewUnderNoneIgnoresCodeRabbitConfiguration|TestDeliveryPublicationRequestsNoCodeRabbitReview|TestRunEnforcesReviewRequestCoherence)$" ./internal/config ./internal/cli ./internal/preflight 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDefaultConfigScopesReviewSourceToPullRequestFeedback TestDefaultConfigSelectsNoPrePRReviewProvider TestReviewUnderCodexIgnoresCodeRabbitConfiguration TestReviewUnderNoneIgnoresCodeRabbitConfiguration TestDeliveryPublicationRequestsNoCodeRabbitReview TestRunEnforcesReviewRequestCoherence; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "never asked for a CodeRabbit review" && ! { tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "across these repositories"; } && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "never selects or requests a pre-PR reviewer" && tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF -- "never selects or requests a pre-PR reviewer" && tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "never selects or requests a pre-PR reviewer" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the five new named tests exists and neither phrase is documented, so the command fails.

## References

- [_techspec.md](_techspec.md) — Unselected providers receive no request
- `_prd.md` → Goal 1; Core Feature 1; Success Metric 1
- `_techspec.md` → API Contract 6; Testing Approach 1

## Result

### Implementation

- Scoped the CodeRabbit request instructions to a `coderabbit` Pre-PR Review
  Policy or the legacy `fetch`, `watch` and `resolve` workflow. The scoped
  section preserves every request instruction, states that `codex`, `claude`
  and `none` are never asked for a CodeRabbit review by the pre-PR workflow,
  and no longer contains `across these repositories`.
- Labeled `review_source` in the command, usage and configuration guides and
  in generated User and Project Config as the legacy PR-feedback source read
  only by `fetch`, `watch` and `resolve`. Generated config still writes
  `request_review: false` and no `pre_pr_review` key.
- Added package-level contract tests for generated config scope and default
  provider resolution, `roundfix review` under `codex` and `none` beside an
  incoherent legacy CodeRabbit configuration, and CodeRabbit-free Delivery
  Queue publication under every policy.
- Regenerated `skills/roundfix` from the canonical skill with
  `rtk make skills-sync`. `rtk make baseline-digests` reported that derived
  artifacts already matched their canonical sources and changed no files.

### Focused checks and acceptance evidence

1. Skill scope: `rtk rg -n 'never asked for a CodeRabbit review'
   .agents/skills/roundfix/SKILL.md` found the scoped statement;
   `rtk rg -n 'across these repositories'
   .agents/skills/roundfix/SKILL.md` returned no match (expected exit 1); and
   `rtk diff -r .agents/skills/roundfix skills/roundfix` exited 0.
2. Review and publication scope: `rtk env
   GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run
   '^(TestReviewUnderCodexIgnoresCodeRabbitConfiguration|TestReviewUnderNoneIgnoresCodeRabbitConfiguration|TestDeliveryPublicationRequestsNoCodeRabbitReview)$'
   ./internal/cli` exited 0. The `none` case also asserted zero Agent runner
   probes, prepares, prompts and session closes.
3. Generated config: `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test
   -count=1 -run
   '^(TestDefaultConfigScopesReviewSourceToPullRequestFeedback|TestDefaultConfigSelectsNoPrePRReviewProvider)$'
   ./internal/config` exited 0. The loading assertion resolved provider
   `codex` from source `default`.
4. Legacy mirror: `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test
   -count=1 -run '^TestRunEnforcesReviewRequestCoherence$'
   ./internal/preflight` exited 0 without changing the test or legacy
   coherence code.
5. Repository incremental check: the first sandboxed `rtk make
   verify-incremental` exposed the preserved-skill wording contract and two
   host process-table permission refusals. After restoring the scoped exact
   request instruction, `TestReviewRequestContract` passed; a fresh
   `rtk make verify-incremental` with host process-table access exited 0.

The Task's authored `## Verification` command was not run; the Daemon owns
that command and terminal settlement.
