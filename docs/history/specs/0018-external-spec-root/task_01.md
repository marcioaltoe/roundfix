---
task: task_01
spec: 0018-external-spec-root
status: completed
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

- [x] Config key, default, and precedence wiring
- [x] Resolution to an absolute Spec Root with the external predicate
- [x] Validation failure for a missing or non-directory root
- [x] Config tests: default, relative, absolute, symlinked, external, invalid

## Acceptance Criteria

- [x] With no configuration, the resolved Spec Root is the repository's
      `docs/specs` and is classified internal.
- [x] A relative configured root resolves against the repository root; an
      absolute one is used as-is.
- [x] A root that is a symlink to a directory outside the repository — or an
      absolute path outside it — is classified external.
- [x] A nonexistent configured root fails validation naming the resolved
      path.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass, including the
  new Spec Root resolution tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 6. `_techspec.md` → Interfaces: ResolveSpecsRoot;
API Contracts: specs.root; Build Order 1. ADR-0035.

## Result

Implemented `specs.root` in config with built-in default `docs/specs`, strict
`specs` key parsing, User Config and Project Config precedence, and generated
default-config output. Added `ResolveSpecsRoot`, which resolves relative roots
against the repository root, keeps absolute roots as configured, validates the
resolved path is an existing directory, and sets `External` after symlink
evaluation against the repository working tree.

Acceptance evidence:

- Default root: `TestResolveSpecsRootUsesLoadedDefault` passes and verifies no
  config resolves to `<repo>/docs/specs` with `External == false`.
- Relative and absolute roots: `TestResolveSpecsRoot` passes and verifies a
  relative `configured-specs` root resolves under the repository root while an
  absolute root is returned unchanged.
- External predicate: `TestResolveSpecsRoot` passes for an absolute external
  root, and `TestResolveSpecsRootClassifiesExternalSymlink` passes for
  `docs/specs` symlinked to a directory outside the repository.
- Invalid roots: `TestResolveSpecsRootRejectsInvalidRoots` passes and verifies
  missing and non-directory roots fail with messages naming the resolved path.

Verification:

- `rtk go test ./internal/config/`: passed, 64 config package tests.
- `rtk make verify`: passed; `go test ./...` reported 850 passing tests,
  `roundfix skills check` passed, and `go build` completed.
