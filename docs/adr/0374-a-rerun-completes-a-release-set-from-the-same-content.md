---
status: accepted
created_at: 2026-11-06T23:30:00Z
updated_at: 2026-11-06T23:30:00Z
deprecated_at: null
superseded_by: null
---

# A rerun completes a Release Set that stopped part way, from the same content

Twice the release workflow's npm publication stopped after some packages of
the Release Set. v0.130.0 (run 41750342366) published
`@roundfix/cli-darwin-arm64`, and then `npm publish` for
`@roundfix/cli-darwin-x64` failed with `TLOG_CREATE_ENTRY_ERROR: (409) an
equivalent entry already exists in the transparency log`. v0.132.1 (run
42056741468) published the five platform packages and failed the same way for
`roundfix`. In both runs the workflow ran `npm publish` once per coordinate and
retried nothing: the failure was printed under `publish:` and the token
fallback, open only for identity errors, was closed. The evidence points to a
duplicate request inside one `npm publish`: the Sigstore client that signs the
provenance retries its own transparency-log request by default (two retries,
a five-second timeout, and `fetchOnConflict: false` in `sigstore-js`
`packages/client/src/config.ts`), so a request that landed but timed out is
resent and refused as a duplicate, and the package is never sent to the
registry. Both failures came about 46 seconds after the publish began. The npm
debug log that would prove the retry was not kept, so this cause is inferred.

The rerun of the failed job then refused in the Publication Preflight with
`registry: <package>@<version> is already used` for every package the first
attempt had published. ADR-0146 made publication all-or-nothing because a
partial release is irreversible under npm policy, but the preflight read every
published coordinate as a conflict, so a partial release could only be
abandoned: 0.130.1 and 0.133.0 republished the same changes under the next
version.

## Decision

**A coordinate already published at the target version counts as done for
the Release Set when its published content equals this run's build; the
rerun publishes only the coordinates that are missing.**

- The Publication Preflight classifies a coordinate whose packument lists the
  target version under `versions` as `published`, which does not stop it.
  A version that appears only under `time` (published, then unpublished)
  stays `used` and stops it, as do `cooldown`, `absent` and `undetermined`.
  When every coordinate is `published`, the preflight stops with
  `registry: every coordinate of the Release Set is already published at
  <version>`.
- On a tag run, after the build and before any publish, each `published`
  coordinate is compared with this run's build: `npm pack` of the staged
  package and `npm pack <coordinate>@<version>` from the registry, both
  decompressed, must be byte-identical tar streams. A difference stops the
  job with `registry: <coordinate>@<version> is already published with
  different content`; a tarball that cannot be packed, fetched or
  decompressed stops it under `undetermined:`.
- The publish step skips each coordinate that compared equal and publishes
  the others in the same order as before, platform packages first and the
  launcher last.
- The comparison is of the decompressed tar stream, not of the npm tarball
  integrity. On 2026-11-06 the v0.132.1 platform packages rebuilt with Go
  1.26.10 on another machine held byte-identical files and tar streams but
  different gzip bytes, so a different `dist.integrity`: the gzip output
  depends on the Node build. Two builds and packs on one machine were
  byte-identical. A build from another Go patch release differs, and is
  refused.

ADR-0146 still holds: no package publishes unless every coordinate is
eligible or already published with the same content, so the Release Set never
mixes two builds, and the bounded token fallback is unchanged.

## Considered Options

- Compare the tarball integrity (`dist.integrity` or the shasum): rejected,
  because the integrity changes with the gzip implementation of the Node
  build that packs it, even when every packed byte is the same.
- Record the first attempt's integrity in a workflow artifact and compare
  with it: rejected, because the registry already holds the published
  tarball, which is the truth to compare with, and an artifact can be missing
  or expired.
- Retry `npm publish` inside the job after a `TLOG_CREATE_ENTRY_ERROR`: not
  part of this decision. A new `npm publish` process signs with a new
  ephemeral key, so a retry would likely succeed, but it changes the publish
  step's failure contract; this decision makes the rerun of the failed job
  the recovery path.
- Skip a published coordinate without comparing it: rejected, because a
  rebuild from another commit or toolchain would then complete the Release
  Set with platform binaries that differ from the launcher's build.
