---
status: pending
type: docs
---

# Task: The skill ships with the CLI change

This repository carries a HARD RULE: a pull request that changes CLI behavior
ships the roundfix skill update with it. The skill documents `reconcile` with
`--apply` alone and mentions neither disposition flag Spec 0092 shipped, so an
Agent following it cannot reach the command this Spec exists to make reachable.

## Tooling authorization

Express maintainer authorization granted 2026-08-27, recorded at
`docs/workflow/authorizations/2026-08-27-carry-forward-reaches-an-unresolved-run.md`.

Bounded files, and nothing else:

- `.agents/skills/roundfix/SKILL.md`

The generated copy under `skills/roundfix/SKILL.md` is rewritten by the
declared `make skills-sync` and is sanctioned fallout of the authorized source
edit, not a separate target. A hand-edited generated copy is an unauthorized
mutation. Changing any other path fails this Task; stop before mutating one.

## Work

- The reconcile guidance names `--carry-forward` and `--discard-superseded`
  beside `--apply`, and records which Run outcomes carry-forward accepts.
- The implement guidance describes the Preflight refusal and the command that
  clears it, so an Agent that meets the refusal can act on it.
- Run the declared skill sync so the generated copy matches the source.
- Stay inside the grant. It does not cover the skill's resolve, watch, settle,
  archive, baseline, or release guidance, nor its agent bundles.
- Claims are read from the delivered code, not from the TechSpec draft.

## References

- `_prd.md` → Goal 5, User Story 5, Core Feature 7; Project Constraints,
  Tooling authority
- `_techspec.md` → Build Order 6
- `docs/agents/specific-repository.md` carries the HARD RULE this Task settles

## Verification
- `grep -q -- "--carry-forward" .agents/skills/roundfix/SKILL.md && grep -q -- "--discard-superseded" .agents/skills/roundfix/SKILL.md && grep -q -- "--carry-forward" skills/roundfix/SKILL.md && go test -count=1 -tags docscontract ./internal/docscontract 2>&1 | grep -q "^ok"`
