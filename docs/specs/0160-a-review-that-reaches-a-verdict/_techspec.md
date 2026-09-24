---
spec: 0160-a-review-that-reaches-a-verdict
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# A review that reaches a verdict

## Executive Summary

Classify the reviewer's verdict by substance and keep its raw answer, run the
`claude` provider through the existing read-only review path, and give the
reviewer the decisions of any Spec the candidate carries.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local ACP sessions only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0153, ADR-0155 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — standing grant recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Classification

`classifyReviewCommandResult` normalizes each line of the answer — trimmed,
case-folded, with surrounding Markdown emphasis and trailing punctuation removed
— and looks for a line equal to `no findings` and a line starting with
`findings:`. Exactly one kind present decides the outcome; findings text is
everything after the `Findings:` line. Both or neither blocks with a reason
naming the ambiguity. The raw message is written to `pre-pr-review-answer.txt`
beside `pre-pr-review.json` for every outcome that reached the reviewer, and the
record gains `answerPath`.

## The claude provider

The provider switch sends `claude` down the same path as `codex`:
`validateReviewProfileProvider` already requires the `review` profile's runtime
to match the provider, and the read-only access applies to any ACP runtime.
`coderabbit` keeps its refusal, reworded to name the missing local surface.

## Spec-aware prompt

From the candidate's changed paths, the command collects the Spec folders under
the configured Spec Root. For each, it reads `_prd.md`'s `## Decisions` section
and `_techspec.md`, appends them to the prompt under a labelled block, and asks
the reviewer to report any implementation choice that contradicts a decision or
adopts an alternative the Spec rejected. The record gains `specs`, the slugs
consulted. A candidate with no Spec folder gets today's prompt.

## API Contracts

1. `roundfix review` classifies a verdict by substance and writes the raw answer
   beside its record, whose `answerPath` names it.
2. `pre_pr_review.provider: claude` runs a read-only review.
3. The review record's `specs` names every Spec whose decisions the prompt
   carried.

## Coverage Map

- Goal 1 → Classification; API Contract 1.
- Goal 2 → The claude provider; API Contract 2.
- Goal 3 → Spec-aware prompt; API Contract 3.
- Core Feature 1 → Classification.
- Core Feature 2 → The claude provider.
- Core Feature 3 → Spec-aware prompt.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → Classification, The claude provider, Spec-aware prompt.

## Integration Points

- **Spec 0156.** The delivery loop's reviewing stage calls this command and
  benefits from every change here without depending on its prompt text.
- **Spec 0153.** Owns the record shape this Spec extends with `answerPath` and
  `specs`.

## Testing Approach

1. **Classification.** Table-driven over exact, punctuated, emphasized,
   preamble, lowercase, both-verdict and no-verdict answers through a fake
   runner; the answer file exists for each reached outcome.
2. **Provider.** A `claude` policy with a matching `review` profile produces a
   record with provider `claude`; `coderabbit` is refused naming the missing
   surface.
3. **Spec awareness.** A candidate adding a Spec folder yields a prompt with its
   Decisions and a record naming it; a candidate without one yields the prior
   prompt and an empty `specs`.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Classification and kept answers (depends on: none).
2. The claude provider (depends on: 1).
3. Spec-aware prompt (depends on: 2).
4. Shipped skill and guide (depends on: 1, 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4, 6, 7).
6. No findings header escapes, no answer is invented (depends on: 4).
7. Spec context that never blocks and stays bounded (depends on: 6).

## Risks & Considerations

- **Tolerance as a loophole.** The classifier accepts presentation variants only;
  an answer carrying both verdicts blocks.
- **Prompt size.** A TechSpec adds context; the Decisions section and TechSpec
  are bounded per Spec, and the diff stays the primary input.
