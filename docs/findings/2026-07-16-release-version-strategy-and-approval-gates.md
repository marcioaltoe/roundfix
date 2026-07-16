---
status: done
created_at: 2026-07-16
updated_at: 2026-07-16
---

# Release instructions — semantic version selection and approval gates (2026-07-16)

The current release path validates that a pushed tag has the shape
`vMAJOR.MINOR.PATCH`, but it does not tell a maintainer or Agent which component
to increment. Implement one shared instruction strategy that classifies the
release from consumer impact and requires explicit human approval before any
major or minor increment.

## 1. Release instructions do not define the next version or its authorization boundary

- **Symptom / evidence**: `docs/user-guide/release-runbook.md` asks the
  maintainer to finalize `## [<version>]` and push `v<version>`. The workflow
  rejects malformed tags, but neither surface defines when to increment major,
  minor, or revision/patch. They also do not stop an Agent from choosing a
  major or minor version and mutating the changelog or tag without approval.
- **Impact**: a release can understate a compatibility break, overstate a fix as
  a feature release, or make a high-impact version decision that the user did
  not authorize. The tag remains internally consistent while communicating the
  wrong upgrade risk to npm, Go, and CLI consumers.
- **Community guidance**: Semantic Versioning assigns incompatible public API
  changes to major, backward-compatible functionality to minor, and
  backward-compatible bug fixes to patch. It also requires minor and patch to
  reset after a major increment, and patch to reset after a minor increment.
  [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html) defines the base contract.
  [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
  maps `fix` to patch, `feat` to minor, and `BREAKING CHANGE` to major, while
  leaving other commit types without an automatic release meaning. Go's
  [module version guidance](https://go.dev/doc/modules/version-numbers) and
  npm's [semantic versioning guidance](https://docs.npmjs.com/about-semantic-versioning/)
  use the same compatibility signals. The
  [semantic-release commit analyzer](https://github.com/semantic-release/commit-analyzer)
  applies the highest required release type when multiple changes are present.
- **Action / request**: add the strategy below to the authoritative release
  instructions and expose it through the Agent instruction path used to cut a
  release. Keep `AGENTS.md` as a short pointer if the detailed rule belongs in
  the release runbook or a repo-owned `cut-release` skill.

## Requested strategy

Before proposing a version, inspect every consumer-visible change since the
latest release tag. Roundfix's public contract includes CLI commands and flags,
stdout/stderr, exit codes, machine-readable fields, configuration keys, release
assets, npm package behavior, and documented compatibility guarantees. Classify
the release by the highest-impact change, not by commit count or perceived
effort.

| Increment | Use when | Version effect | Permission |
| --- | --- | --- | --- |
| Revision / patch | The release contains only backward-compatible fixes or internal changes that preserve the public contract. `fix` and compatible `perf` changes normally belong here. | `X.Y.Z` → `X.Y.(Z+1)` | A generic user request to cut a release authorizes this increment; no second version prompt is required. |
| Minor | The release adds backward-compatible functionality, deprecates public behavior, or makes a substantial compatible capability available to consumers. `feat` normally belongs here. | `X.Y.Z` → `X.(Y+1).0` | Show the proposed version and evidence, then wait for explicit user approval before editing release files, tagging, or publishing. |
| Major | The release removes, renames, or incompatibly changes any public contract and requires consumers to change usage, configuration, parsing, automation, imports, or installation behavior. | `X.Y.Z` → `(X+1).0.0` | Show the compatibility break, migration impact, and proposed version, then wait for explicit user approval before editing release files, tagging, or publishing. |

Commits such as `docs`, `test`, `ci`, `chore`, and compatible `refactor` have no
automatic SemVer effect. If they do not change shipped behavior or artifacts,
the strategy must allow “no release.” If they must be published to deliver a
consumer-visible fix, classify that impact instead of relying only on the
commit type.

For mixed changes, the order is `major > minor > revision/patch`. A fix bundled
with a backward-compatible feature is minor; a feature bundled with a breaking
change is major.

## Version-zero rule

Roundfix currently releases under `0.y.z`. SemVer and Go both treat version zero
as initial development without a stability guarantee. The instruction must not
use that allowance to hide compatibility impact:

- compatible fixes increment revision/patch;
- new compatible capabilities increment minor;
- incompatible public-contract changes also require a minor increment while
  the project remains on major zero, but the approval prompt must label the
  release as breaking;
- promotion from `0.y.z` to `1.0.0` is a major increment and requires explicit
  approval.

The future spec must decide whether Roundfix is ready to declare `1.0.0`; this
finding does not make that product decision.

## Approval interaction

For a computed major or minor increment, the Agent must stop before any write or
publish action and present:

1. current version and proposed version;
2. required increment type;
3. the specific changes that require it;
4. compatibility and migration impact;
5. the exact question: `Approve the <major|minor> increment to <version>?`

Only an explicit approval continues the release. A user request that already
names the target version, such as `release v0.5.0`, counts as approval for that
increment after the Agent verifies the version is not smaller than the impact
requires. A generic request such as `cut the next release` does not authorize a
computed major or minor increment.

Read-only analysis is allowed before approval. Editing `CHANGELOG.md`, changing
version-bearing files, creating or pushing a tag, publishing packages, and
creating the GitHub Release are blocked until approval. Existing authorization
rules for commits, pushes, and external publication still apply independently.

## Acceptance criteria for the future implementation

1. Starting from `v0.4.0`, a fix-only delta proposes `v0.4.1` without an extra
   version prompt after the user asks to cut a release.
2. Starting from `v0.4.1`, a compatible feature proposes `v0.5.0` and performs
   no release mutation until the user approves the minor increment.
3. Starting from `v1.4.2`, a breaking public-contract change proposes `v2.0.0`
   and performs no release mutation until the user approves the major increment.
4. Starting from a `v0.y.z` tag, a breaking contract change proposes the next
   minor version, labels it breaking, and requires minor-increment approval.
5. A mixed delta selects the highest required increment regardless of commit
   order.
6. A docs/test/CI-only delta can report that no release is required.
7. An explicitly requested target version counts as approval only when it is at
   least as large as the classified impact requires; otherwise the Agent stops
   and explains the mismatch.
8. The existing tag-driven version agreement, verification, package publishing,
   asset naming, and GitHub Release workflow remain unchanged.

## Follow-up spec

Routed to [Spec 0034 — Release Plan](../specs/0034-release-plan/_prd.md), with the read-only and confirmation-gated authority boundary recorded in [ADR-0048](../adr/0048-release-planning-is-read-only-and-confirmation-gated.md).
