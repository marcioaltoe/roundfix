# Release runbook

Roundfix releases are cut by pushing a `v*` git tag. The tag-triggered workflow
at `.github/workflows/release.yml` cross-compiles every target, publishes the
npm launcher and per-platform packages, and uploads GitHub Release assets in one
run. This runbook is for maintainers cutting a release; end users install with
`npx`/`bunx`/global npm (see the README **Install** section). See ADR-0031 for
why distribution goes through npm platform packages.

Every release starts with the read-only Release Plan Command:

```bash
roundfix release plan
```

Run it before changelog edits, version-file edits, tags, pushes, package
publication, asset uploads, or GitHub Release creation. The command creates no
Run, reads no Roundfix configuration, contacts no external service, and
mutates no repository or release state.

The plan's state controls what a generic release request authorizes:

- `ready` with a patch proposal can proceed without another version approval.
- `approval_required` for a minor, major, or version-zero breaking proposal
  requires the human to answer the printed approval question before any release
  mutation.
- `manual_classification_required` requires a rerun with
  `--impact <none|patch|minor|major> --reason <text>`. Manual classification
  records the release impact and reason, but it does not approve a resulting
  minor, major, or version-zero breaking version.
- `no_release` means there is no release to cut for the committed range.

For ordinary range planning, after the plan's required decision is satisfied,
keep using the tag-triggered publication workflow below.

## Planning the 0.0.1 release history reset

The exceptional `0.0.1` reset starts with a clean checkout whose `HEAD` is
committed, then runs:

```bash
roundfix release plan --reset-to v0.0.1
```

This mode is read-only. It inventories every local and remote stable tag and
every GitHub Release through complete pagination, sorts the
inventory deterministically, and binds the reset target, committed target
revision, tags, and releases to `planDigest`. Text output lists every immutable
identity and target commit. Use `--format json` for the
`roundfix.release-plan/0.0.1` structured result.

`--reset-to` cannot be combined with `--from`, `--to`, `--impact`, or
`--reason`. A complete reset plan returns state `approval_required`, prints the
approval question, and exits `3`. Invalid flags, a dirty tree, an invalid
target, incomplete local or remote tag access, or incomplete GitHub Release
inventory exits `2` without a partial plan. The command creates no Run, reads
no Roundfix configuration, and never edits files, refs, tags, remotes,
packages, releases, or configuration.

Exit `3` is the planning boundary, not authorization to delete anything. The
Release Plan Command exposes no tag or GitHub Release deletion action. Review
the complete inventory and digest, then stop. Implementation and QA may inspect
this plan, but they must not mutate release history.

After implementation and QA pass, rerun the same command to obtain a fresh
inventory and `planDigest`. If the inventory or target changed, the new digest
replaces the earlier one. Any later tag or GitHub Release deletion is a
separate destructive release operation and requires explicit human approval
for that fresh plan. Approval of setup, implementation, QA, a prior plan, or
the public CLI's printed question does not grant deletion authority.

## Agent Selection CLI contract correction

Release notes for this change must identify a partial Agent Selection override
as rejected usage. `resolve`, `watch`, and `implement` now accept only these two
forms:

- omit all three selection flags to use Agent Selection Profiles;
- provide `--agent`, `--model`, and `--reasoning-effort` together for one complete override.

For example:

```bash
roundfix implement --spec <slug>
roundfix implement --spec <slug> --agent codex --model gpt-5.6-sol --reasoning-effort high
```

A bare `--agent`, or any other proper subset, is rejected with exit `2` before configuration, proof, or Run mutation. This is an intentional public CLI correction, not a
runtime regression.

## Version agreement

The pushed tag is the single source of truth for the release version. The
workflow strips the leading `v` and:

- injects it into the binary via `-ldflags "-X roundfix/internal/app.Version=<tag>"`;
- sets it on every per-platform package and on the launcher with
  `npm version <tag> --no-git-tag-version`;
- rewrites the launcher's `optionalDependencies` to pin the same version.

Nothing can disagree because every artifact is stamped from the one tag —
**after** the run is allowed to start. `Validate tag` gates that, and it makes
**two** checks before any build:

1. the tag is a semver version (`vMAJOR.MINOR.PATCH`, with an optional
   pre-release suffix);
2. the tag matches the **checked-in launcher version** in
   `dist/npm/roundfix/package.json`.

There are **two** checked-in version sources, and `TestVersionMatchesTheReleaseManifest`
keeps them equal:

- `dist/npm/roundfix/package.json` — what `Validate tag` compares against;
- `internal/app/version.go` (`var Version`) — the binary's fallback when no
  ldflags stamp is applied.

The second check is why the tag alone is not enough to cut a release. Bump both
files and commit them **before** creating the tag, or the run stops with

```text
tag v<version> does not match the checked-in Roundfix version <checked-in>
```

The five per-platform packages under `dist/npm/packages/` are **not** checked:
the workflow stamps them from the tag, so they stay at their placeholder
version in the repository. Only the launcher is compared.

Recovering from a tag pushed against a stale launcher version means deleting
and recreating that tag, which is a destructive release operation needing its
own explicit approval. Bumping first avoids it. Measured on 2026-08-06: the
v0.4.0 tag was pushed against a checked-in 0.3.1 and the run stopped at
`Validate tag` before publishing anything.

Keep `CHANGELOG.md` ahead of the tag: the GitHub Release notes are extracted
from the `## [<version>]` section, falling back to a bare `Release <tag>` line
when no section matches.

## Configuring npm Trusted Publishing

Configure GitHub Actions as a trusted publisher for every package before the
first OIDC release. npm stores this binding per package, so configuring the
launcher does not configure any platform package.

1. Sign in to npmjs.com with maintainer access to the `roundfix` and
   `@roundfix` packages.
2. Open each package's settings, add a GitHub Actions trusted publisher, and
   enter the exact values in this table:

   | Coordinate | GitHub owner | Repository | Workflow filename |
   | --- | --- | --- | --- |
   | `roundfix` | `marcioaltoe` | `roundfix` | `release.yml` |
   | `@roundfix/cli-darwin-arm64` | `marcioaltoe` | `roundfix` | `release.yml` |
   | `@roundfix/cli-darwin-x64` | `marcioaltoe` | `roundfix` | `release.yml` |
   | `@roundfix/cli-linux-arm64` | `marcioaltoe` | `roundfix` | `release.yml` |
   | `@roundfix/cli-linux-x64` | `marcioaltoe` | `roundfix` | `release.yml` |
   | `@roundfix/cli-win32-x64` | `marcioaltoe` | `roundfix` | `release.yml` |

3. Review all six entries for exact spelling before cutting the first OIDC
   release. Each entry must bind to the `marcioaltoe/roundfix` repository and
   `.github/workflows/release.yml`.

npm validates a trusted-publisher configuration only when `npm publish`
attempts the OIDC exchange. Saving the configuration does not test it; an
owner, repository, or workflow typo first appears during publication as an
`identity:` authentication failure.

## Rehearsing the Publication Preflight

Use the release workflow's manual `workflow_dispatch` trigger to exercise the
Publication Preflight without cutting a tag:

1. Open the `release` workflow in GitHub Actions and choose **Run workflow**.
2. Enter `v<version>` in the required `version` input. The version must match
   the checked-in launcher version.
3. Inspect the Publication Preflight table and any failure prefix in the run
   summary.

This rehearsal is publish-free. It validates the version, runs `make verify`,
and checks registry eligibility for the full Release Set, then stops before
cross-compilation, npm publication, or GitHub Release creation.

## Closing the fallback window

The fallback window protects the first OIDC release from a trusted-publisher
configuration error that npm cannot reveal before publish time. While the
GitHub Actions repository variable `NPM_TRUSTED_PUBLISHING_FALLBACK` is exactly
`1`, an OIDC authentication failure retries that coordinate with `NPM_TOKEN`
and records the coordinate under **npm Trusted Publishing fallback** in the
job summary. This keeps one bad package binding from leaving a partial Release
Set.

Only one result closes the window: one complete tagged release whose fallback
record is empty. The job summary states `No coordinates required the bounded
token fallback.` A successful release with any coordinate in that record does
not close the window; correct that coordinate's trusted-publisher binding and
prove another complete release.

After the empty record exists, remove the repository variable (or set it to a
value other than `1`) so an OIDC failure stops instead of retrying. Remove the
`NPM_TOKEN` repository secret at the same time. Do not disallow token
publication before that empty record exists: while any coordinate still needs
the token fallback, disabling token publication would break the release path.

As part of the same window-closing procedure, use npm package settings to
disallow token publication with the per-package setting for each Release Set
coordinate:

- `roundfix`
- `@roundfix/cli-darwin-arm64`
- `@roundfix/cli-darwin-x64`
- `@roundfix/cli-linux-arm64`
- `@roundfix/cli-linux-x64`
- `@roundfix/cli-win32-x64`

Reopen each package's settings and confirm that its per-package setting still
reports token publication as disallowed. The window is not closed until all
six confirmations pass. Then remove the bounded fallback branch from
`release.yml` in a dedicated workflow change.

## Release workflow failure vocabulary

The prefix identifies the failing stage and the maintainer's first recovery
step:

| Prefix | Emitting stage | Meaning | First recovery step |
| --- | --- | --- | --- |
| `registry:` | Publication Preflight | The named coordinate cannot accept the target version because it is already used, remains in the printed post-unpublish cooldown, or is absent. | Follow the named reason: choose an unused version, wait until the printed cooldown time, or restore the absent coordinate; then rerun the rehearsal. |
| `undetermined:` | Publication Preflight | Registry eligibility could not be proven because the Release Set, transport, HTTP response, registry body, or cooldown timestamp could not be read. Nothing has published. | Inspect the named read or classification failure and rerun the rehearsal after registry access is reliable. |
| `identity:` | Publish to npm | OIDC authentication failed for the named coordinate. During the open fallback window this can be a warning before the token retry; after the window closes it is fatal. | Compare that package's trusted-publisher owner, repository, and workflow filename with the table above and correct the mismatch. |
| `publish:` | Publish to npm | The named coordinate failed for a reason that is **not** authentication — a network error, a validation rejection, or a registry fault. The underlying npm error is printed beneath the prefix, and the token fallback is never attempted. | Read the printed npm error and act on it directly. A `publish:` failure never means the trusted publisher is misconfigured, so do not change publisher settings in response to one. |
| `runtime:` | Guard Trusted Publishing runtime | The resolved npm runtime is below the Trusted Publishing floor of `11.5.1`. | Restore the workflow's declared Node 24/npm runtime and rerun the workflow. |

## Required secrets

- `NPM_TOKEN` — a temporary npm **Automation** token for the `roundfix` org,
  read only by the bounded fallback retry while the fallback window is open.
  Store it as a repository secret and remove it when the empty fallback record
  closes the window.
- `GITHUB_TOKEN` — the workflow's built-in token, used with `contents: write`
  permission to create the GitHub Release and upload assets. No manual setup.

## Cutting a release

1. Run `roundfix release plan` from a clean checkout and satisfy the decision
   boundary described above. Use the proposed version as `<version>` only after
   the plan is `ready` or the required human approval has been given.
2. Land all release content on `main` and confirm `make verify` is green
   locally. The workflow re-runs the full gate and refuses to publish if it
   fails. No Spec merges without a passing QA Report in its `qa/` directory —
   `roundfix archive` enforces the same rule later, but the gap must be caught
   at merge time, not at archive time.
3. Add or finalize the `## [<version>]` section in `CHANGELOG.md`.
4. Tag and push:

   ```bash
   git tag v<version>
   git push origin v<version>
   ```

5. Watch the `release` workflow. In order it: validates the tag, runs
   `make verify`, runs the Publication Preflight, cross-compiles and stages
   every target, publishes the per-platform packages first and then the
   launcher, and creates the GitHub Release with the matching binary assets.

The per-platform packages publish before the launcher so the launcher's
`optionalDependencies` resolve the moment it goes live. A failed `make verify`,
a disagreeing tag, or a blocked Publication Preflight stops the run before
anything is published.

## How assets feed the Upgrade Command

Each target's binary is uploaded as a GitHub Release asset named with its
`GOOS`+`GOARCH` tokens (for example `roundfix-darwin-arm64`,
`roundfix-linux-amd64`, `roundfix-windows-amd64.exe`). The Upgrade Command
(`roundfix upgrade`) resolves the latest release through the GitHub CLI and
selects the asset for the current platform using the same naming scheme, which
is pinned by `internal/cli/release_asset_compat_test.go`. Changing an asset name
without updating that fixture fails `make verify`, keeping the release channel
and the Upgrade Command in agreement.

## Verifying a published release

```bash
npx roundfix@<version> --version
```

The printed version must match the tag. Newly published per-platform packages
can take a moment to propagate across npm; a transient `E404` for one platform
usually resolves on retry.
