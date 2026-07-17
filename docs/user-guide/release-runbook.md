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

After the plan's required decision is satisfied, keep using the tag-triggered
publication workflow below.

## Version agreement

The pushed tag is the single source of truth for the release version. The
workflow strips the leading `v` and:

- injects it into the binary via `-ldflags "-X roundfix/internal/app.Version=<tag>"`;
- sets it on every per-platform package and on the launcher with
  `npm version <tag> --no-git-tag-version`;
- rewrites the launcher's `optionalDependencies` to pin the same version.

Nothing can disagree because every artifact is stamped from the one tag. The
`Validate tag` step fails the run before any build when the tag is not a semver
version (`vMAJOR.MINOR.PATCH`, with an optional pre-release suffix). Keep
`CHANGELOG.md` ahead of the tag: the GitHub Release notes are extracted from the
`## [<version>]` section, falling back to a bare `Release <tag>` line when no
section matches.

## Required secrets

- `NPM_TOKEN` — an npm **Automation** token for the `roundfix` org, exposed to
  the publish step as `NODE_AUTH_TOKEN`. An automation token is required because
  interactive 2FA cannot be satisfied in CI. Store it as a repository secret.
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
   `make verify`, cross-compiles and stages every target, publishes the
   per-platform packages first and then the launcher, and creates the GitHub
   Release with the matching binary assets.

The per-platform packages publish before the launcher so the launcher's
`optionalDependencies` resolve the moment it goes live. A failed `make verify`
or a disagreeing tag stops the run before anything is published.

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
