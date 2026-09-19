---
task: task_03
spec: 0148-a-profile-that-declares-both-tiers
status: pending
type: chore
complexity: low
---

# Task 03: Regenerate the derived guides and pins

## Overview

The repository's own guides are generated from the Baseline sources this Spec
changed. This slice runs the sanctioned regeneration and commits its output, so
the generated copies agree with what owns them.

## Requirements

1. MUST regenerate the repository's guides from their Baseline sources through
   the sanctioned command.
2. MUST rewrite any derived pin the asset change moves, only through the
   sanctioned command.
3. MUST NOT hand-edit any line inside a setup-context marker.
4. MUST leave the clause texts unchanged: this Spec makes the declared value
   exist, it does not rewrite what the clauses ask for.
5. MUST NOT change any path outside this Spec's bounded files and its own Task
   file.

## Subtasks

- [ ] Run the sanctioned regeneration.
- [ ] Confirm the generated guides publish both tiers.
- [ ] Confirm the regeneration contract passes.

## Acceptance Criteria

- [ ] The generated guides name the incremental command.
- [ ] The regeneration contract test passes.
- [ ] No clause text changed.
- [ ] No path outside the bounded list changed.

## Context

- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`

## Verification

- `grep -rq "verify-incremental" docs/agents/` — expected: exit 0; the generated guidance names the repository's fast command. Before this Task it does not.
- `grep -rq "verification.incremental" docs/agents/ && go test -count=1 -tags repocontract -run "^TestDeclaredStepRegenerationAndFrozenBoundaries$" ./internal/baseline > /tmp/0148-regen.txt 2>&1 && grep -q "^ok" /tmp/0148-regen.txt` — expected: exit 0; the regenerated guidance carries the declared tier and the regeneration contract holds over it. Before this Task the guidance does not name it, so the command fails.

## References

`_prd.md` → Core Feature 3; Declared intentional breaks; Regression locks;
`_techspec.md` → Implementation Design: Regeneration, not editing;
Build Order 3; `_authorization.md`.
