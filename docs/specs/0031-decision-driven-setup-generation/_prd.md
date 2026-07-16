---
spec: 0031-decision-driven-setup-generation
status: active
created: 2026-07-16
surfaces: [cli, docs]
---

# Decision-driven setup generation

## Problem statement

Manual variation QA for spec 0030 proved that the setup workflow stores its
durable decisions but does not make most answers affect generated guidance. It
can therefore report a clean audit while managed instructions contradict
confirmed choices. The same QA run also proved that the first-run response asks
the right questions but cannot show the promised managed-change preview before
the maintainer answers them.

## Goals

1. Make every decision in the initial setup catalog control an observable
   generated rule, artifact, value, or validation policy.
2. Show a complete read-only preview before all decisions are answered,
   distinguishing definite changes from decision-dependent changes.
3. Make preview, audit, and apply agree on one resolved setup state.
4. Reuse compatible answers from existing Setup Manifests and refresh only
   setup-owned content without repeating settled questions.
5. Preserve the atomicity, idempotency, portability, ownership, skill-safety,
   and Secondbrain-safety behavior that passed spec 0030 QA.

## Core features

1. The portable catalog declares the effects of every durable setup decision.
2. A false boolean omits the rules and guides controlled by that decision.
3. Enabling a capability requests only the dependent answers that capability
   needs; disabling it does not ask irrelevant follow-up questions.
4. Enum and string answers select or populate the managed guidance that owns
   them.
5. Preview names the selected profile and skill setup, definite managed
   operations, and conditional operations with their triggering decision and
   value.
6. Audit detects managed content that disagrees with the durable decisions,
   and apply repairs only proven setup-owned boundaries.
7. Existing compatible manifests migrate without confirmation; only new,
   missing, or invalid decisions require an answer.

## Success criteria

- The QA-03 first-run flow returns all unresolved decisions and a non-empty
  preview that covers every definite or conditional managed artifact without
  writing to the target repository.
- The QA-07 alternate-decision flow renders the selected domain layout,
  runtime values, and Verification command, while omitting local Spec and
  autonomous guidance when their decisions are false.
- Every initial decision has a macro test that proves its observable effect
  rather than only its manifest persistence; decisions with multiple accepted
  values cover at least one non-default value.
- A manifest produced by spec 0030 reuses every compatible answer, updates its
  managed inventory and content, and reaches a clean second audit without
  another confirmation.
- Repeated apply remains byte-for-byte idempotent and repository-authored bytes
  remain unchanged.

## Non-goals

- Adding new project profiles, setup presets, or user-facing decision IDs.
- Adding interactive prompts to the Python CLI; the skill continues to ask one
  unresolved question at a time.
- Changing skill installation, extra-skill classification, or removal policy.
- Validating or generating nested agent-instruction files.
- Reading from or writing to the Secondbrain.
- Replacing the Setup Manifest or ownership-marker formats.
