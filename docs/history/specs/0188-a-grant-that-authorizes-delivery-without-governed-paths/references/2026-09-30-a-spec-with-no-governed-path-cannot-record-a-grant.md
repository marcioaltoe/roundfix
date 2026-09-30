---
type: fix
status: promoted
created: 2026-09-30
spec: 0188-a-grant-that-authorizes-delivery-without-governed-paths
reason: null
---

# A Spec that changes no Governed Path cannot record a deliverable grant

## Symptom

Spec 0186 changes only ordinary files. `deliver plan` blocked it with `authorization refused: paths`, because the authorization reader requires at least one exact path. Listing the ordinary files it changes then failed CI: `TestEveryBoundedPathIsGoverned/repository_records_are_governed` requires every listed path to be governed. No record is valid for a Spec with delivery operations and no Governed Path, so such a Spec can only be delivered by hand. PR #287 was closed for this reason.

## Where

- The `paths` rule in `internal/authorization/authorization.go` (about line 601).
- `TestEveryBoundedPathIsGoverned` (repository contract test, `make repo-test`).
- The `deliver start` and `deliver plan` authority check.

## Expected

A grant whose consuming Spec changes no Governed Path can authorize its delivery operations with an empty `paths` list. Any governed mutation still needs its exact path in the list, and a listed path must still be governed.

## Evidence

`bin/roundfix deliver plan 0186-…` on 2026-09-30, and the failed Verification gate of PR #287 (run 36701160969).
