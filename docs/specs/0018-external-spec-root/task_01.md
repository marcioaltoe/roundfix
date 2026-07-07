---
task: task_01
spec: 0018-external-spec-root
status: pending
type: backend
complexity: low
---

# Task 01: specs.root config with resolution and the external predicate

## Overview

Add the `specs.root` configuration key and the resolution that turns it into
one absolute Spec Root per command invocation, including the predicate that
classifies the resolved root as external to the repository working tree.
Verifiable on its own through config package tests.

## Requirements

1. MUST add `specs.root` (string, built-in default `docs/specs`) with the
   existing precedence: Project Config over User Config over the built-in
   default. Unknown-key strict validation behavior is unchanged.
2. MUST resolve relative values against the repository root of the user's
   checkout and accept absolute values as-is, producing one absolute Spec
   Root plus an external flag.
3. MUST classify the root as external when, after symbolic link evaluation,
   it lies outside the repository working tree.
4. MUST fail validation with an actionable message naming the resolved path
   when the root does not exist or is not a directory.

## Subtasks

- [ ] Config key, default, and precedence wiring
- [ ] Resolution to an absolute Spec Root with the external predicate
- [ ] Validation failure for a missing or non-directory root
- [ ] Config tests: default, relative, absolute, symlinked, external, invalid

## Acceptance Criteria

- [ ] With no configuration, the resolved Spec Root is the repository's
      `docs/specs` and is classified internal.
- [ ] A relative configured root resolves against the repository root; an
      absolute one is used as-is.
- [ ] A root that is a symlink to a directory outside the repository — or an
      absolute path outside it — is classified external.
- [ ] A nonexistent configured root fails validation naming the resolved
      path.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass, including the
  new Spec Root resolution tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 6. `_techspec.md` → Interfaces: ResolveSpecsRoot;
API Contracts: specs.root; Build Order 1. ADR-0035.
