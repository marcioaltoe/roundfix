---
spec: 0012-npm-distribution
status: archived
created: 2026-07-06
archived: 2026-07-06
release: v0.1.0
qa_override: true
surfaces: [cli, infra, docs]
---

# npm Distribution and Skill Bundle

Roundfix has no public install path. A user who wants the CLI must clone the
repository and build it, and the Upgrade Command can only pull release assets
that no workflow produces yet. Node is already a hard prerequisite of the
agent layer, so the shortest install every target platform shares is npm. This
Spec ships Roundfix through npm as a launcher package with per-platform binary
packages, driven by a tag-triggered release workflow that cross-compiles the Go
binary for every supported platform and publishes both the npm packages and the
GitHub Release assets the Upgrade Command already consumes. The distributed
binary also carries the Roundfix skill bundle — the operational Roundfix Skill
plus the authorial workflow skills the repository owns — so an installed
Roundfix can teach an Agent runtime not just how to operate the CLI but the
whole spec-driven method it drives.

## Goals

- `npx roundfix`, `bunx roundfix`, and a global install run the current
  Roundfix on darwin, linux, and Windows across both architectures, from one
  npm channel. See ADR-0031.
- A `v*` tag triggers a workflow that cross-compiles, verifies version
  agreement, publishes the npm packages, and uploads the matching GitHub
  Release assets — no manual release steps.
- The Upgrade Command keeps working: the workflow feeds the same release-asset
  channel it already reads.
- The npm wrapper is transparent: it forwards arguments and exit codes
  verbatim, so Roundfix's stable exit-code contract survives the Node shim.
- The binary ships the Roundfix-owned skill bundle (the operational Roundfix
  Skill plus the authorial workflow skills), and `roundfix skills` manages the
  whole set — sync, check, and install — not the operational skill alone.

## User Stories

1. As a developer trying Roundfix, I want `npx roundfix` to run the current
   release without cloning or building, so that first use costs one command.
2. As a developer on any supported platform, I want the install to fetch the
   binary built for my OS and architecture, so that I never compile Go myself.
3. As a developer who installed Roundfix globally, I want `roundfix upgrade` to
   keep finding new releases, so that the npm install and the self-upgrade path
   stay consistent.
4. As a maintainer cutting a release, I want pushing a `v*` tag to build,
   verify, and publish every platform package and release asset, so that
   releasing is one tag push with no manual asset uploads.
5. As a maintainer, I want the release to fail fast when the tag, the npm
   package version, and the embedded app version disagree, so that a
   mismatched release never publishes.
6. As a developer scripting Roundfix behind `npx`, I want the launcher to exit
   with Roundfix's own exit code, so that my scripts read the real result
   through the Node wrapper.
7. As a developer who installed Roundfix, I want `roundfix skills install` to
   place the whole Roundfix-owned skill bundle into my Agent runtime, so that I
   get the spec-driven workflow skills the CLI drives, not just the operational
   skill.
8. As a developer, I want `roundfix skills check` and a skills listing to cover
   the whole owned bundle and point me at the recommended external skills, so
   that I can tell what Roundfix ships from what I install through my own skills
   tooling.

## Core Features

1. **Launcher package.** A root npm package named `roundfix` with a `bin` shim
   that resolves the correct per-platform binary package and executes it,
   forwarding every argument and the child's exit code and signals unchanged.
   Per-platform packages are declared as `optionalDependencies` so only the
   matching one installs. See ADR-0031.
2. **Per-platform binary packages.** One package per supported target
   (darwin-arm64, darwin-x64, linux-arm64, linux-x64, win32-x64), each carrying
   the prebuilt binary and constrained by npm `os`/`cpu` fields so the wrong
   platform is never selected.
3. **Tag-triggered release workflow.** A GitHub Actions workflow on `v*` tags
   that (a) checks the tag, the launcher package version, and the embedded app
   version all agree; (b) runs the full local verification gate; (c)
   cross-compiles the Go binary per platform with a `GOOS`/`GOARCH` matrix;
   (d) publishes the per-platform packages, then the launcher; (e) uploads the
   binaries as GitHub Release assets named so the Upgrade Command's platform
   selection matches them.
4. **Exit-code and signal fidelity.** The launcher shim is a pass-through: it
   never rewrites stdout/stderr, never adds output of its own on success, and
   propagates the child's exit code and terminating signal so the CLI contract
   holds through the wrapper.
5. **Upgrade-channel continuity.** Release assets keep the naming the Upgrade
   Command's platform selection already expects (asset name contains the
   platform's OS and architecture tokens), so `roundfix upgrade` resolves them
   without changes.
6. **Roundfix skill bundle.** The binary embeds the Roundfix-owned skill set —
   the operational `roundfix` skill plus the 13 authorial workflow skills the
   repository owns — and `roundfix skills` manages all of them: `check`
   validates the whole set, `install` writes the whole set to the chosen
   target, and a listing distinguishes the owned bundle from the recommended
   external skills a user installs through their own skills tooling. The
   embedded set stays in sync with the canonical `.agents/skills/` copies and
   `make verify` fails on drift. Externally-managed skills (the
   `skills-lock.json` set) are never embedded or modified — they are only
   listed as recommendations.

## User Experience

A first-time user runs `npx roundfix <command>` and gets the current release
behaving exactly as the built binary, including exit codes. A global installer
runs `npm i -g roundfix` and then `roundfix` directly. Nothing about the CLI's
own output, flags, or exit codes changes — the wrapper is invisible in normal
use. A maintainer pushes a `v*` tag and watches one workflow build, verify,
and publish everything.

## Non-Goals / Out of Scope

- Homebrew, a tap, or a formula — deferred: it excludes Windows and needs a
  maintained tap for reach npm already covers. See ADR-0031.
- Embedding, vendoring, or modifying the externally-managed skills (the
  `skills-lock.json` set from `marcioaltoe/skills`) — Roundfix ships only the
  skills it owns; the external ones are listed as recommendations and installed
  through the user's own skills tooling.
- A skill dependency resolver that pulls the external skills a workflow skill
  references — the listing names them; fetching them is the user's tooling job.
- Linux distro packages (apt, dnf, AUR) and Windows installers (winget, MSI).
- Changing the Upgrade Command's download or platform-selection logic — the
  workflow adapts to it, not the reverse.
- Signing or notarizing the binaries in this Spec (tracked separately for the
  agent-runtime hygiene work).
- Publishing pre-release or nightly channels — released `v*` tags only.

## Success Metrics

- On each supported platform, `npx roundfix --version` reports the tagged
  version and runs the platform's own binary.
- A `v*` tag with agreeing versions publishes all per-platform packages plus
  the launcher and uploads matching Release assets in one workflow run; a tag
  whose versions disagree fails the workflow before publishing.
- After a global npm install, `roundfix upgrade --check` resolves the latest
  release asset for the running platform.
- The launcher returns Roundfix's exit code unchanged for a success and a
  known failure exit code (asserted).
- `roundfix skills install --target <t>` writes all 14 owned skills to the
  target; `roundfix skills check` validates the whole owned set; the embedded
  set matches `.agents/skills/` (no `make verify` drift).

## Decisions

- Distribution is an npm launcher plus per-platform binary packages under
  `optionalDependencies`; the tag-triggered workflow uploads GitHub Release
  assets to keep the Upgrade Command fed. See ADR-0031.
- Homebrew is deferred; npm is the single first channel. See ADR-0031.
- Release-asset names must remain compatible with the Upgrade Command's
  existing OS/architecture substring selection.
- The distributed binary carries the Roundfix-owned skill bundle (operational
  `roundfix` + 13 authorial workflow skills); `roundfix skills` manages the
  whole set. External `skills-lock.json` skills are listed, never embedded.
- The operational `roundfix` skill keeps its version tracking the released CLI
  version (`v*` tag); the authorial skills keep their own independent versions.

## Open Questions

None.
