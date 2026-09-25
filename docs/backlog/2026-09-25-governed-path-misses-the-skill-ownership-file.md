---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `GovernedPath` does not match `skills/_ownership.yml`, and the contract test that says so is outside the gates

## Symptom

`go test -tags repocontract -run TestEveryBoundedPathIsGoverned ./internal/speccheck` fails on `main` (595876fc): archived Spec 0163 bounds `skills/_ownership.yml`, which `GovernedPath` does not classify as governed. Neither `make verify` nor `make verify-docs` runs the `repocontract` tag, so the failure is silent.

## Where

`internal/speccheck/governed.go` (patterns and literal list); the Makefile targets that choose build tags.

## Expected

`skills/_ownership.yml` is governed (it decides skill ownership), and the `repocontract` contract tests run in a repository gate.

## Evidence

Found while authoring Spec 0170 on 2026-09-25.
