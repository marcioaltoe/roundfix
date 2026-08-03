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

## Addendum — security review hardening

CodeRabbit raised two Major findings on the pin bump, both accepted:

- Every action is pinned to a full commit SHA rather than a mutable tag, across
  all three workflows. Each SHA was verified against its tag ref through the
  forge API before use — a wrong pin fails the workflow, so the bot's values
  were checked rather than trusted.
- `persist-credentials: false` on the `release.yml` and `ci-conventions.yml`
  checkouts. Neither workflow performs a Git write — confirmed by grep — and
  the release step already passes `GH_TOKEN` explicitly for `gh release`. The
  `secondbrain-sync.yml` checkout that pushes keeps its credentials, since
  disabling them there would break the sync.

A third finding followed: disable the implicit dependency caches. The release
job holds `id-token: write` and publishes, so a cache restored before the
verification, build, and publish steps is a supply-chain surface. `cache: false`
on `setup-go` and `package-manager-cache: false` on `setup-node`. Both input
names were confirmed against the pinned revisions' own `action.yml` before use —
`package-manager-cache` defaults to `true`, so leaving it unset was not neutral.

This extends the authorized paths to `.github/workflows/ci-conventions.yml` and
`.github/workflows/secondbrain-sync.yml`, which the pin change reaches.

## Not verified here

GitHub Actions cannot be exercised locally. These pins are proven only by the
next workflow run. The release path retains its bounded token fallback
(ADR-0084), so a failure in the publish stage is recoverable rather than
partial.


## Addendum — explicit dist-tag

The first real OIDC release failed on its first coordinate with:

```text
npm error Cannot implicitly apply the "latest" tag because previously published
version 0.3.0 is higher than the new version 0.0.3.
```

The platform packages still carry `0.1.0`, `0.2.0`, and `0.3.0` from before the
version reset, and npm 11 — the floor Trusted Publishing requires — refuses to
apply `latest` implicitly when a higher version exists. The 0.0.2 release
predates that npm version and did not hit it.

Publication now passes `--tag latest` on both the OIDC attempt and the bounded
fallback retry. The maintainer authorized the release; this is the fix that
makes it complete.

- `.github/workflows/release.yml` — explicit dist-tag on both publish paths.

Nothing was published by the failed run. The failure was correctly classified
`publish:` rather than `identity:`, so the bounded token fallback was not spent
on a non-authentication error — the behavior Spec 0057's QA-002 repair
introduced, exercised here for the first time in a real release.
