---
granted: 2026-08-10
action: correct-review-exclusion-rationale
paths:
  - .coderabbit.yaml
consuming: direct
---

# Tooling authorization — the archive exclusion stops claiming the wrong mechanism (2026-08-10)

On 2026-08-10 the maintainer directed:

> _archived deve sempre ficar de fora da revisão do coderabbit também. valide
> @.coderabbit.yaml

## What the validation established

The exclusion works; its written justification does not.

- **Patterns cover every archived tree.** `!**/_archived/**` reaches
  `docs/specs/_archived` and `docs/findings/_archived`, `!**/_reviews/**`
  reaches `docs/specs/_reviews`, at any depth and for every extension found
  inside them.
- **The file is schema-valid.** Validated against CodeRabbit's published
  `schema.v2.json`: zero errors.
- **Every archived-path review comment in history predates the exclusion.**
  PR #136 looked like a counterexample — its head carried the exclusion and
  still drew five comments on `_archived` paths — until the timestamps
  settled it: the comments landed at 17:04 UTC and the commit adding the
  exclusion reached the branch at 20:33 UTC. After the exclusion, no reviewed
  PR has an archived-path comment.
- **The comment's mechanism is wrong.** The block says the exclusions sit
  last "para vencer qualquer inclusão acima", assuming gitignore-style
  last-match-wins. CodeRabbit's glossary documents the opposite model: "Files
  not matching any exclude pattern are included by default" — an exclusion
  wins regardless of position. And position cannot be relied on at all under
  `inheritance: true`, because inherited arrays append the organization's
  unique entries *after* this file's, so nothing this file writes can be
  guaranteed last.
- **`knowledge_base.filePatterns` documents no `!` negation**, so it cannot
  be hardened that way; its protection is that no pattern matches an archived
  tree, which holds today for every rule-named file (checked: none exists
  under any archived tree).

## Authorized paths

- `.coderabbit.yaml`, limited to the comment block explaining the archived
  exclusion and a matching note in the knowledge-base section. No pattern
  changes.

## Bounded by purpose

A future maintainer who trusts the current comment could add an include after
the block believing order governs, or reorder the list believing position
protects the archive. The correction states the documented semantics so the
exclusion's safety stops resting on a mechanism CodeRabbit does not promise.

## Consuming Spec

Applied directly: it corrects two comments in one configuration file.

## Commit choreography

This record lands as its own commit, before the comment correction.
