---
spec: 0012-npm-distribution
prd: _prd.md
created: 2026-07-06
---

# npm Distribution — Technical Spec

## Executive Summary

Ship the existing Go binary through npm without touching the CLI's behavior. The
primary trade-off is adding a Node launcher shim in front of a native binary:
accepted because Node is already a hard prerequisite of the agent layer, Go
cross-compiles every target from one runner (no per-OS build machines, unlike a
Rust toolchain), and the shim stays a transparent pass-through so the stable
exit-code contract survives it. The release workflow's job is agreement and
fan-out: prove tag == launcher version == embedded app version, run the local
gate, cross-compile the `GOOS`/`GOARCH` matrix, publish per-platform packages
then the launcher, and upload GitHub Release assets whose names the Upgrade
Command's existing platform selection already matches.

## System Architecture

- New packaging tree (e.g. `dist/npm/`), not Go code:
  - `roundfix` launcher package — `package.json` (`type: commonjs`,
    `optionalDependencies` per platform), `bin/roundfix` shim.
  - Per-platform packages `@roundfix/cli-<os>-<arch>` (darwin-arm64, darwin-x64,
    linux-arm64, linux-x64, win32-x64) — each `package.json` sets `os`/`cpu` and
    carries the prebuilt binary.
- New `.github/workflows/release.yml` triggered on `v*` tags.
- `internal/app` — the embedded version (already read by the Upgrade Command via
  `app.NormalizeVersion`/`app.CompareVersions`) is the source of truth the
  version-agreement check compares against.
- Upgrade Command untouched: `selectPlatformAsset` already matches an asset whose
  name contains the running `runtime.GOOS` and `runtime.GOARCH` tokens; the
  workflow names assets to satisfy it.
- `skills/skills.go` grows from a single embedded `roundfix` skill to the
  Roundfix-owned bundle: an owned-skills manifest (the operational `roundfix`
  plus the 13 authorial workflow skills), embedded from a synced bundle
  directory, plus a recommended-external manifest (names only) derived from
  `skills-lock.json`. `Check`, `Install`, and a new `list` operate over the
  manifest; `make skills-sync`/`skills-sync-check` cover the whole bundle.

## Implementation Design

### Launcher shim

```js
// bin/roundfix — resolve the per-platform package's binary and exec it.
// 1. Map process.platform/process.arch → @roundfix/cli-<os>-<arch>.
// 2. require.resolve the binary inside the installed optional dependency.
// 3. spawnSync(binary, process.argv.slice(2), { stdio: 'inherit' }).
// 4. Propagate: if signal → re-raise/exit per signal; else process.exit(status).
// No stdout/stderr of its own on success; a missing platform package is the
// only shim-level error (clear message: unsupported platform / reinstall).
```

The shim adds nothing to stdout on success and forwards the child's exit code
and terminating signal verbatim, so `roundfix`'s exit-code contract holds through
`npx`/`bunx`/global installs.

### Platform tokens — npm vs Go

npm package names use the npm platform token (`win32`); the binary and the
GitHub Release asset use the Go `GOOS` token (`windows`). Keep them distinct:
`@roundfix/cli-win32-x64` wraps a `GOOS=windows GOARCH=amd64` binary, and its
Release asset name contains `windows` and `amd64` so `selectPlatformAsset`
(which matches against `runtime.GOOS`) resolves it. A mapping table is the single
source that drives package names, `os`/`cpu` fields, `GOOS`/`GOARCH`, and asset
names together.

### Release workflow (`v*` tag)

```
on: push: tags: ['v*']
jobs:
  guard:    tag == launcher package.json version == embedded app version, else fail
  verify:   make verify (the full local gate)
  build:    matrix GOOS/GOARCH → cross-compile, stage each into its @roundfix/cli-* package
  publish:  npm publish every @roundfix/cli-* (public), then npm publish roundfix
  release:  upload the same binaries as GitHub Release assets (names carry GOOS+GOARCH)
```

Per-platform packages publish before the launcher, so the launcher's
`optionalDependencies` are resolvable the moment it is live. The Release-asset
upload keeps the Upgrade Command's channel fed from the same artifacts.

### Version agreement

The guard job reads three versions — the pushed tag, the launcher
`package.json` version, and the embedded app version — normalizes them
(`app.NormalizeVersion` semantics), and fails the workflow before any build when
they disagree, so a mismatched release never publishes.

### Roundfix skill bundle

The embedded skill set expands from `roundfix` alone to the Roundfix-owned
bundle. The ownership boundary is fixed by CLAUDE.md: the operational `roundfix`
skill plus the 13 authorial workflow skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`, `setup-workflow`, `implement-task`,
`implement-spec`, `brainstorming`, `council`, `business-analyst`,
`archive-spec`, `qa-gate`, `evidence-gate`). Everything in `skills-lock.json`
stays external and is never embedded or modified.

```go
// Owned bundle manifest (14 entries) drives embed, sync, check, install, list.
var ownedSkills = []string{"roundfix", "write-idea", "write-prd", /* … */}

// Per-skill validation policy — roundfix is strict, the rest are structural.
//   roundfix:   existing required wording + Roundfix branding + openai manifest
//   authorial:  SKILL.md present, frontmatter parses, name == dir, no banned
//               "reference project" branding; NO Roundfix-branding requirement
//               and NO version-tracks-tag requirement (they keep own versions)
func Check() []Diagnostic
```

- **Embed shape.** `make skills-sync` copies each owned skill from the canonical
  `.agents/skills/<name>/` into a synced bundle directory under `skills/` that
  `skills.go` embeds as a unit (isolated from the package's `.go` files), and
  writes a recommended-external manifest (names read from `skills-lock.json`)
  for `skills list`. `skills-sync-check` diffs the whole synced bundle and fails
  `make verify` on any drift — the same guarantee the single skill has today,
  widened to 14.
- **`skills install`.** Already generalizes over the embedded set (it walks
  `Files()` and writes each), so it writes all 14 owned skills to the chosen
  target with no per-skill special-casing; the `--target`/`--dir` semantics are
  unchanged.
- **`skills check`.** Applies the per-skill policy above: `roundfix` keeps its
  strict contract-wording and openai-manifest checks; the authorial skills get
  structural validation only, so their independent versions and generic
  authorial language never trip the check.
- **`skills list`.** New read-only subcommand printing the owned bundle and,
  separately, the recommended external skills (names from the embedded
  recommended manifest) with a one-line note that external skills install
  through the user's own skills tooling — never by Roundfix.
- **Version rule.** The `roundfix` skill's `metadata.version` continues to track
  the released CLI version (`v*` tag); the authorial skills keep their own
  versions, and the check does not couple them.

## Coverage Map

- Stories 1-2 → launcher shim + per-platform packages (`os`/`cpu` selection)
- Story 3, Feature 5 → asset naming compatible with `selectPlatformAsset`
- Stories 4-5 → release workflow (matrix build, publish, version-agreement guard)
- Story 6 → shim exit-code/signal propagation
- Story 7 → owned-bundle embed + `skills install` over the whole set
- Story 8 → per-skill `skills check` policy + `skills list` (owned vs external)

## Integration Points

npm registry (publish) and GitHub Releases (asset upload) — both from the
tag-triggered workflow with repository-scoped tokens. The Go toolchain
cross-compiles all targets on one runner. No runtime external systems: the shim
only execs a local binary.

## Testing Approach

- Shim: a Node-level test (or a shell harness invoked from the workflow) asserts
  the shim forwards args and returns the child's exit code for a success and a
  known non-zero failure, and errors clearly when no platform package is
  present. Kept out of the Go suite; runs in the release/CI lane.
- Version agreement: a unit-testable helper (or a workflow step with a golden
  fixture) covering agree/disagree cases.
- Asset-name compatibility: a Go test asserting `selectPlatformAsset` resolves
  the workflow's asset-name scheme for every matrix entry — this ties the
  workflow's naming to the Upgrade Command's selection so they cannot drift.
- The matrix build itself is validated by the workflow producing all five
  binaries; a dry-run `npm pack` per package checks file inclusion.

## Build Order

1. Platform mapping table + per-platform `@roundfix/cli-*` package scaffolding
   with `os`/`cpu` fields (no deps)
2. Launcher `roundfix` package + `bin` shim with pass-through exit/signal
   semantics (depends on: 1)
3. `selectPlatformAsset` compatibility test pinning the asset-name scheme
   (depends on: 1)
4. Roundfix skill bundle: owned-skills manifest, widened embed + `make
   skills-sync`/`skills-sync-check`, per-skill `Check` policy, `skills install`
   over the set, `skills list`, recommended-external manifest (no deps)
5. Release workflow: version-agreement guard, `make verify` (now covering the
   widened bundle sync/check), matrix cross-compile, publish order,
   Release-asset upload (depends on: 1, 2, 3, 4)
6. Docs (install + release runbook + skill bundle) and skill sync (depends on:
   2, 4, 5)

## Risks & Considerations

- The npm-vs-Go platform token split is the sharpest footgun; the single mapping
  table and the `selectPlatformAsset` compatibility test are the guardrails.
- Publish ordering: a failed launcher publish after platform packages succeed
  leaves platform packages orphaned for that version — acceptable (they are
  inert without the launcher) but the workflow should not retry into a partial
  state; publish platform packages, then launcher, and let a failed run be
  re-cut on a new tag.
- The shim must never buffer or reinterpret I/O — `stdio: 'inherit'` and verbatim
  exit/signal propagation keep the CLI contract intact.
- Binaries are unsigned/un-notarized in this Spec; notarization is out of scope
  and tracked with the agent-runtime hygiene work.
- Widening the embedded skill set widens the drift surface: `skills-sync-check`
  must diff every owned skill, and `Check`'s per-skill policy must not impose
  the operational skill's strict wording on the authorial ones (which carry
  generic language and independent versions) — otherwise `make verify` breaks
  on legitimate authorial-skill edits. The authorial skills stay editable as
  ordinary repo files (CLAUDE.md ownership rule); external skills are never
  embedded or modified.
- The recommended-external manifest is names only, derived from
  `skills-lock.json` at sync time — no external skill content is copied, so
  there is no upstream-drift or licensing coupling from the external set.

## Decisions

- npm launcher + per-platform packages under `optionalDependencies`; workflow
  uploads GitHub Release assets to feed the Upgrade Command. See ADR-0031.
- Asset names carry `GOOS`+`GOARCH` tokens to match `selectPlatformAsset`; a Go
  test pins the scheme.
- Homebrew deferred. See ADR-0031.
- The binary embeds the Roundfix-owned skill bundle (operational + 13 authorial)
  behind an owned-skills manifest; `skills sync/check/install/list` operate over
  it; `Check` is strict for `roundfix`, structural for the authorial skills.
- External `skills-lock.json` skills are listed, never embedded or modified
  (CLAUDE.md skill-ownership rule).
