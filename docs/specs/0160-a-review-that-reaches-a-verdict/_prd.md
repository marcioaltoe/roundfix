---
spec: 0160-a-review-that-reaches-a-verdict
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# A review that reaches a verdict

`roundfix review`, delivered by Spec 0153, runs the configured reviewer over the
current candidate. In use it rarely reaches a verdict.

**Well-formed answers are refused.** The classifier accepts exactly
`No findings` or a message starting with `Findings:`. On 2026-09-24 the Codex
reviewer's answer on Spec 0157 was recorded `blocked: unclassifiable agent
output`, and because the raw answer is discarded nobody can tell whether it said
`No findings.`, wrapped the verdict in emphasis, or put a sentence before it.
Every such block sends the maintainer to a manual review outside the workflow.

**Only one provider runs.** `claude` is a valid policy value that the command
refuses, although the read-only session it needs is the same one Codex uses.

**The review cannot see the Spec.** The reviewer judges the diff for
correctness, but not whether the implementation obeys the Spec's decisions and
the alternatives it rejected — the half of Spec 0129 Core Feature 1 that waited
for a workflow-run reviewer.

This Spec is delivery 4 of the restructured queue: Spec 0126 Core Features 2
and 3 as they apply to agent providers, the review classifier debt recorded
after Spec 0153, and the semantic half of Spec 0129 Core Feature 1. A local
CodeRabbit surface stays refused with its reason: no supported local CLI is
installed or specified.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local ACP sessions only; no
  credential handling and no network call by Roundfix itself. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0153 keeps review evidence bound to the
  candidate it examined, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the standing grant of 2026-09-18 for keeping
  the shipped skill true to the CLI behavior a slice delivers, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A reviewer's answer that follows the contract in substance is classified, and
  every answer is kept as evidence.
- The `claude` policy runs a review.
- A review of a Spec's delivery judges it against the Spec's decisions.

## Core Features

1. **Tolerant classification, kept evidence.** The verdict line is found
   regardless of letter case, trailing punctuation, Markdown emphasis or a
   preamble before it. An answer that carries both verdicts or neither still
   blocks. The raw answer is written beside the review record on every outcome,
   and a blocked record names its path.
2. **The `claude` provider.** With `pre_pr_review.provider: claude`, the review
   runs through the `review` profile on a read-only Claude session, with the
   same evidence, fallback and blocking rules as Codex. `coderabbit` stays
   refused, naming the missing local surface.
3. **Spec-aware review.** When the candidate adds or changes a Spec folder under
   the Spec Root, the prompt carries that Spec's PRD Decisions and TechSpec and
   asks the reviewer to judge the implementation against the decisions and the
   alternatives they reject. The record names the Specs consulted. With `none`,
   no review of any kind is claimed.

## Non-Goals / Out of Scope

- A local CodeRabbit review surface.
- Finding dispositions and fresh-review enforcement across candidates, which the
  delivery loop of Spec 0156 owns.
- Changing the exit-code contract of `roundfix review`.

## Success Metrics

1. Answers such as `No findings.`, `**No findings**` and a preamble followed by
   `Findings:` are classified; an answer with both or neither verdict blocks;
   the raw answer file exists for every outcome.
2. A `claude` policy produces a review record with provider `claude`.
3. A candidate carrying a Spec folder produces a prompt with its Decisions and a
   record naming the Spec.

## Decisions

- **Tolerate form, not meaning.** Only presentation varies; an answer that does
  not commit to exactly one verdict still blocks, so tolerance never turns an
  unclear review into a pass.
- **Keep every answer.** Evidence that is discarded cannot be audited, and the
  block on Spec 0157 could not be diagnosed for that reason.
- **Read the Spec from the candidate.** The Spec a delivery carries is in its own
  diff, so no extra argument or configuration is needed to find it.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a classifier that read an
ambiguous answer as a pass would pass every happy path while approving
unreviewed work.

## Research basis

`classifyReviewCommandResult` in `internal/cli/review.go` compares the trimmed
message with `No findings` and the `Findings:` prefix and keeps nothing else of
it. The provider switch in the same file refuses `claude` and `coderabbit`. The
Spec 0157 review record of 2026-09-24 carries `reason: unclassifiable agent
output`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
