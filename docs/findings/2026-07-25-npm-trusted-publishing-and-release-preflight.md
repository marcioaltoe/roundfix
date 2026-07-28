---
status: done
created_at: 2026-07-25
updated_at: 2026-07-28
---

# npm publishing — token authentication and registry state can produce a partial release (2026-07-25)

The `v0.0.1` release reset exposed two independent publication risks. Roundfix
still authenticates npm publication with a long-lived repository secret, and
the release workflow does not prove that every package can accept the target
version before it begins publishing the platform packages.

Session evidence:

- Release Run
  [`30166957746`](https://github.com/marcioaltoe/roundfix/actions/runs/30166957746)
  failed in the Verification gate because
  `TestBaselineApplyCommandRealCLI` used the macOS-only
  `/private/tmp/roundfix-task09-go-cache` path on a Linux runner.
- Commit `f98a12fcb99015d37be1fb32c26a0e2e57fd0c6e` replaced that path with
  a test-owned temporary directory, and the focused test and local
  `make verify` passed.
- Release Run
  [`30167217166`](https://github.com/marcioaltoe/roundfix/actions/runs/30167217166)
  was canceled during its Verification gate before cross-compilation or
  publication.
- At cancellation time, the unscoped `roundfix` package reported that the
  whole package had been unpublished on `2026-07-25T17:14:13.763Z`.
  The five `@roundfix/cli-*` packages remained published at `0.3.0`.
- npm blocks every new version of a fully unpublished package for 24 hours.
  A previously published exact `package@version` can never be reused. See the
  [npm unpublish policy](https://docs.npmjs.com/policies/unpublish/).

## 1. The release workflow depends on a long-lived npm credential

- **Symptom / evidence**:
  [`.github/workflows/release.yml`](../../.github/workflows/release.yml)
  exposes `secrets.NPM_TOKEN` as `NODE_AUTH_TOKEN` to six sequential
  `npm publish --access public` commands. The
  [release runbook](../user-guide/release-runbook.md) describes the secret as
  an npm Automation token.
- **Root cause**: The publication contract predates npm Trusted Publishing.
  Authentication is repository-wide secret state rather than a short-lived
  identity bound to the GitHub repository, workflow file, and hosted runner.
- **Action / suggestion**: Migrate the release workflow to
  [npm Trusted Publishing](https://docs.npmjs.com/trusted-publishers/) with
  GitHub Actions OIDC. Configure `marcioaltoe/roundfix` and `release.yml` as
  the trusted publisher for the launcher and every platform package, grant
  only the `npm publish` action, add `id-token: write`, and use supported Node
  and npm versions. Keep `NPM_TOKEN` until one complete release proves the
  OIDC path; then remove the workflow secret and disallow token publication
  for the owned packages.

## 2. Publication starts before registry eligibility is established

- **Symptom / evidence**: The workflow publishes every `@roundfix/cli-*`
  package before publishing `roundfix`. During the reset, the platform
  packages could accept `0.0.1`, but the launcher name was inside npm's
  24-hour post-unpublish block. Without cancellation, publication could have
  stopped after creating five installable platform packages but before
  creating the launcher or GitHub Release.
- **Root cause**: `Validate tag` checks SemVer and the checked-in owned
  version, but no read-only step checks all registry coordinates, the exact
  target version, an unpublish cooldown, or package ownership before the first
  irreversible publish.
- **Action / suggestion**: Add a publication preflight after Verification and
  before cross-compilation. It must evaluate the launcher and all platform
  packages as one release set, refuse an already-used target version, detect a
  package-name cooldown or missing ownership, and emit the exact blocked
  coordinate. No package may publish unless every coordinate is eligible.

## 3. Migration must preserve the existing release contract

- **Symptom / evidence**: The release workflow stamps every binary and npm
  package from the Git tag, publishes platform packages before the launcher,
  and creates the GitHub Release last. The ordering protects launcher
  installation only when the full npm set is eligible.
- **Root cause**: Authentication migration and release-set validation affect
  the same irreversible boundary, but neither requires a new package layout
  or version authority.
- **Action / suggestion**: Preserve the Git tag as the release version
  authority, the existing package names, platform asset names, changelog
  extraction, and GitHub Release assets. Treat release-set preflight and OIDC
  authentication as separate Tasks so a failure identifies whether registry
  state or identity blocked publication.

## Routing

A future Spec must define the npm preflight contract, OIDC migration,
credential rollback window, package-by-package trusted-publisher setup, and
real release validation. This finding does not authorize renaming packages,
reusing unpublished versions, or changing the current `v0.0.1` tag.

## Addendum — 2026-07-28 — Routed to Spec 0058

Owned by
[Spec 0058 — npm Trusted Publishing and release preflight](../specs/0058-npm-trusted-publishing-and-release-preflight/_prd.md);
the maintainer expressly authorized the `.github/workflows/release.yml`
mutation in that Spec's Tooling authority entry.
