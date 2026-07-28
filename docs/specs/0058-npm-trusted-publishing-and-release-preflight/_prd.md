---
spec: 0058-npm-trusted-publishing-and-release-preflight
status: active
created: 2026-07-28
surfaces: [infra, docs]
---

# npm Trusted Publishing and release preflight

The `v0.0.1` release reset exposed two independent publication risks that
remain open. Roundfix's release workflow authenticates six sequential
`npm publish` commands with a long-lived repository secret, and it starts
publishing platform packages before proving that every package in the
release set can accept the target version — during the reset, the launcher
name sat inside npm's 24-hour post-unpublish block while the five platform
packages were publishable, a combination that without cancellation would
have produced a partial release: installable platform packages with no
launcher. Evidence:
[token authentication and registry state can produce a partial release](../../findings/2026-07-25-npm-trusted-publishing-and-release-preflight.md).

## Project Constraints

- Identifier strategy: not applicable — package names, tags, and versions
  keep their existing identities; no project-owned Internal Identifier is
  created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — publication authentication moves
  from a long-lived npm token to GitHub Actions OIDC Trusted Publishing;
  no Roundfix-runtime authentication or HTTP surface changes, and registry
  access stays inside the release workflow. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0031 keeps Roundfix shipping
  through npm platform packages with the existing layout; ADR-0048 keeps
  release planning read-only and confirmation-gated (the Release Plan
  Command is unchanged; the preflight lives in the workflow); ADR-0082
  makes publication all-or-nothing across the release set. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.github/workflows/release.yml`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- npm publication authenticates through short-lived OIDC identity bound to
  this repository and workflow, with the long-lived token removed after one
  proven release.
- No package publishes unless every coordinate in the release set is
  eligible for the exact target version.
- A publication failure names whether identity or registry state blocked
  it, and which coordinate.
- The existing release contract survives unchanged: tag as version
  authority, package names, asset names, ordering, changelog extraction,
  GitHub Release last.

## User Stories

1. As the maintainer cutting a release, I want the workflow to publish via
   npm Trusted Publishing with OIDC, so that publication no longer depends
   on a repository-wide secret that can leak or outlive its purpose.
2. As the maintainer, I want a read-only publication preflight after
   Verification and before cross-compilation that evaluates the launcher
   and every platform package as one release set, so that an ineligible
   coordinate stops the release before the first irreversible publish.
3. As the maintainer diagnosing a failed release, I want the failure to
   name the blocked coordinate and whether the cause was identity or
   registry state, so that recovery starts at the right place.
4. As a user installing Roundfix during a failed release window, I want no
   partial release set on the registry, so that the launcher and platform
   packages never disagree.

## Core Features

1. The release workflow authenticates npm publication through Trusted
   Publishing: OIDC id-token permission, per-package trusted-publisher
   configuration for the launcher and every platform package, and only the
   publish grant.
2. A publication preflight runs after Verification and before
   cross-compilation: for the launcher and all platform packages it
   refuses an already-used exact version, detects a post-unpublish
   cooldown, and detects missing ownership or trusted-publisher
   configuration, emitting the exact blocked coordinate; no package
   publishes unless every coordinate is eligible.
3. Identity migration and release-set preflight are separable failures: a
   failed release names which one blocked it.
4. The long-lived `NPM_TOKEN` remains through one complete OIDC-proven
   release as rollback, then is removed and token publication disallowed
   for the owned packages; the runbook records the rollback window.
5. The release contract is preserved: the Git tag stays the version
   authority, package names and platform asset names are unchanged,
   platform packages still publish before the launcher, changelog
   extraction and GitHub Release creation stay last.

## User Experience

- A blocked preflight prints one line per blocked coordinate with the
  reason (version used, cooldown until when, ownership or publisher
  configuration missing) and stops before any publish.
- The release runbook documents the OIDC setup, the rollback window, and
  the preflight's failure vocabulary.

## Non-Goals / Out of Scope

- Renaming packages, reusing unpublished versions, or changing the
  `v0.0.1` tag.
- Changing the Release Plan Command, which stays read-only per ADR-0048.
- A new registry, scope, or distribution channel.
- Retrying publication automatically after a cooldown.

## Success Metrics

- One complete release publishes all six packages through OIDC with
  `NPM_TOKEN` unused; the follow-up release runs with the secret removed.
- A simulated ineligible coordinate (used version or cooldown) stops the
  workflow before cross-compilation with the coordinate named, and the
  registry receives zero packages from that run.
- The published artifact set, names, and asset contract are byte-compatible
  with the previous release process, and the Upgrade Command consumes the
  release unchanged.

## Decisions

- Publication eligibility is all-or-nothing across the release set: five
  installable platform packages without a launcher is a worse outcome than
  a stopped release. See
  [ADR-0082](../../adr/0082-release-publication-is-all-or-nothing-across-the-package-set.md).
- Preflight lives in the workflow, not in the Release Plan Command —
  planning stays local and read-only per ADR-0048; registry truth is
  checked where publication happens.
- Migration and preflight are separate workflow steps so failures
  attribute cleanly.

## Open Questions

None.
