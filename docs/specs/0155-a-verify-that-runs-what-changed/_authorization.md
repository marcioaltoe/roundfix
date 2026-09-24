---
status: approved
granted: 2026-09-24
action: add a selective verification gate and point Runs and pull request CI at it
consuming: 0155-a-verify-that-runs-what-changed
paths:
  - Makefile
  - .roundfixrc.yml
  - .github/workflows/ci-verify.yml
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0155

Two grants are consumed here.

The maintainer decided on 2026-09-24, when approving the restructured queue, to
keep one repository and give it selective verification gates. A selective gate
is a build-tool target and a Verification setting, so that decision is the
express authority for the three files below.

The standing grant of 2026-09-21 covers any governed path a slice genuinely
needs, on the condition that the bounded set is measured rather than predicted
and recorded here with its reason. The set was measured with `GovernedPath`
itself over every path this Spec changes.

## Why each governed path is unavoidable

- `Makefile` — the selective target `verify-changed` is a build-tool target,
  and `make verify` must keep running the complete gate.
- `.roundfixrc.yml` — `defaults.verification` is what the Daemon runs as the
  repository Verification; pointing it at `make verify-changed` is the point of
  the Spec.
- `.github/workflows/ci-verify.yml` — pull request CI runs the selective gate
  while pushes to `main` keep the complete one.

## What is not governed

`cmd/verify-select`, `internal/verifyselect` and their tests are ordinary
source.

## Approved bounded mutation

Add `verify-changed` beside `verify` in the `Makefile` without changing what
`verify` runs; set `defaults.verification: make verify-changed` in
`.roundfixrc.yml`; run `make verify-changed` for pull requests and keep
`make verify` for pushes to `main` in `.github/workflows/ci-verify.yml`.

## Limits

- No action, operation or path beyond those above.
- `make verify` keeps its meaning: the complete gate.
- No change to the generated guides, Baseline assets or skill files.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
