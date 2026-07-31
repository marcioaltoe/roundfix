---
task: task_01
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: infra
complexity: low
---

# Task 01: Raise the publishing runtime and grant OIDC

## Overview

Trusted Publishing requires npm 11.5.1 or later on Node 22.14.0 or later, and
the release workflow pins Node 20 — whose bundled npm cannot perform an OIDC
exchange at all. This task raises the runtime, grants the workflow an OIDC
identity, and adds a guard that names a too-old toolchain instead of letting
publication silently fall back to token authentication. Verifiable on its own:
the workflow declares the new runtime and permission, and the guard rejects a
version below the floor.

## Requirements

1. MUST raise the release workflow's Node runtime to a version that satisfies
   the Trusted Publishing floor of Node 22.14.0 and ships npm 11.5.1 or later.
2. MUST grant the release job the `id-token: write` permission required to mint
   a GitHub Actions OIDC token.
3. MUST preserve the existing `contents: write` permission, which the GitHub
   Release stage depends on.
4. MUST add a guard step that resolves the running npm version and fails with a
   `runtime:` prefixed message when it is below 11.5.1, so an image change
   surfaces as a named failure.
5. MUST NOT change the registry URL, the tag validation, the Verification gate,
   the cross-compilation stage, publish ordering, or the GitHub Release stage.
6. MUST confine every change to the bounded authorized path
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Raise the `setup-node` version pin to the supported runtime.
- [ ] Add `id-token: write` beside the existing `contents: write` permission.
- [ ] Add the npm version guard step with its `runtime:` failure message.
- [ ] Confirm no other stage of the workflow changed.

## Acceptance Criteria

- [ ] The workflow requests a Node version of 24 and retains
      `registry-url: "https://registry.npmjs.org"`.
- [ ] The workflow's `permissions` block contains both `id-token: write` and
      `contents: write`.
- [ ] A guard step resolves `npm --version` and emits a `runtime:` prefixed
      error naming the 11.5.1 floor when the version is lower.
- [ ] The `Validate tag`, `Verify gate`, `Cross-compile and stage`, and
      `GitHub Release` stages are byte-identical to their previous content.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml` and this task file.

## Context

- interface: `.github/workflows/release.yml`

## Verification

- `grep -q 'node-version: "24"' .github/workflows/release.yml` — expected: exit
  0; the runtime pin is raised.
- `grep -q 'id-token: write' .github/workflows/release.yml` — expected: exit 0;
  OIDC token minting is permitted.
- `grep -q 'contents: write' .github/workflows/release.yml` — expected: exit 0;
  the pre-existing permission survived.
- `grep -q 'registry-url: "https://registry.npmjs.org"' .github/workflows/release.yml`
  — expected: exit 0; the registry target is unchanged.
- `grep -q '11.5.1' .github/workflows/release.yml` — expected: exit 0; the guard
  states the required npm floor.
- `grep -q 'runtime:' .github/workflows/release.yml` — expected: exit 0; the
  guard uses the spec's failure vocabulary.
- `grep -q 'make verify' .github/workflows/release.yml` — expected: exit 0; the
  Verification gate still runs.
- `grep -q 'gh release create' .github/workflows/release.yml` — expected: exit
  0; the GitHub Release stage is intact.

## References

- `_prd.md` → Core Features 1; Non-Goals (release contract preserved).
- `_techspec.md` → System Architecture: Runtime; Build Order 1.
- ADR-0084.
