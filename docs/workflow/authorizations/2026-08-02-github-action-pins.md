# Tooling authorization — GitHub Action version pins (2026-08-02)

A rehearsal run of the release workflow reported:

```text
Node.js 20 is deprecated. The following actions target Node.js 20 but are being
forced to run on Node.js 24: actions/checkout@v4, actions/setup-go@v5,
actions/setup-node@v4
```

The maintainer directed the pins be updated. Action versions are protected
tooling — the rules name version pins explicitly — so the grant is recorded
rather than assumed.

## Authorized paths

- `.github/workflows/release.yml` — `actions/checkout` v4→v7,
  `actions/setup-go` v5→v7, `actions/setup-node` v4→v7.
- `.github/workflows/ci-conventions.yml` — `actions/checkout` v6→v7.

`.github/workflows/secondbrain-sync.yml` already pinned `actions/checkout@v7`
and was not touched.

## Why these versions are safe for this repository

Each major's release notes were read before bumping. Every breaking change
across the range is the Node runtime upgrade the deprecation notice asks for:

- `setup-go` v6.0.0 — "Upgrade Nodejs runtime from node20 to node 24".
- `setup-node` v5.0.0 — "Upgrade action to use node24".
- `checkout` v5.0.0 — "Update actions checkout to use node 24", with a minimum
  runner of v2.327.1 that GitHub-hosted runners exceed.

No input contract this repository uses changed: `node-version`,
`registry-url`, `go-version`, and `check-latest` are unaffected.

`checkout` v7.0.0 adds fork-pull-request blocking for `pull_request_target` and
`workflow_run`. Neither workflow uses those triggers — `ci-conventions.yml`
runs on `pull_request` and `release.yml` on tag push and `workflow_dispatch` —
so the change does not reach this repository.

## Not verified here

GitHub Actions cannot be exercised locally. These pins are proven only by the
next workflow run. The release path retains its bounded token fallback
(ADR-0084), so a failure in the publish stage is recoverable rather than
partial.
