---
spec: 0147-a-planner-that-reads-both-tag-spellings
status: archived
created: 2026-09-18
surfaces: [backend, cli, docs]
archived: "2026-09-18"
source_slug: 0147-a-planner-that-reads-both-tag-spellings
---


# A planner that reads both tag spellings

The Release Plan Command reads a repository's stable tags to decide what the
next version is. It accepts exactly one spelling: a tag must begin with `v`, or
it is rejected as malformed with the advice to "use a tag like v1.2.3".

Plenty of repositories tag `1.2.3`. For them the planner does not report a wrong
version — it reports that their entire release history is malformed, and the
command that is supposed to precede every release refuses before it begins.

Accepting both spellings raises a question the current parser never had to
answer: what happens when `v1.2.3` and `1.2.3` both exist. They are the same
version under two refs, and picking one silently would make the planner's output
depend on which ref a repository happened to create first.

This Spec is the whole of Spec 0128, which is one family and already the size of
a slice.

## Project Constraints

- Identifier strategy: applicable — a tag's exact spelling is part of its identity, so the planner preserves the selected tag's text rather than normalizing it, and coins no identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable to range planning, which reads local Git refs only. Reset planning already reads GitHub Releases through the existing authenticated path, and this Spec adds no transport and changes no credential handling. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the command's read-only contract and its public refusal vocabulary are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is this repository's own tag history and the published guidance the command already carries.
- Tooling authority: applicable — the change is public CLI behavior, and the repository's hard rule ships the Roundfix skill update with it. Express maintainer authorization: standing grant of 2026-09-18 for the Roundfix skill and its mirror across this queue's remaining slices, consumed and bounded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`, and `make baseline-digests` for any derived pin it moves. The planner, the command surface and the user guide are ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A repository that tags `1.2.3` can plan a release without renaming its history.
- A repository that tags `v1.2.3` sees no change at all.
- When one version exists under both spellings, the command refuses and says how
  to resolve it, instead of choosing for the maintainer.
- The proposed version comes back in the spelling of the tag it was selected
  from.

## User Stories

1. As a maintainer of a repository that tags without the prefix, I want the
   planner to read my tags, so that I can plan a release without renaming
   history.
2. As a maintainer whose repository tags with the prefix, I want nothing to
   change, so that this costs me nothing.
3. As a maintainer whose repository carries both spellings of one version, I
   want the command to refuse and name the selector that resolves it, so that I
   decide which ref is the release.
4. As a maintainer reading a proposal, I want the next version spelled like the
   tag it follows, so that I can create it without translating.

## Core Features

1. **Both spellings are stable versions.** A tag of the form `MAJOR.MINOR.PATCH`
   is accepted beside `vMAJOR.MINOR.PATCH`. Malformed input, prerelease
   identifiers and build metadata keep being rejected exactly as they are today.
2. **Selection spans the spellings.** The highest reachable semantic version is
   selected across both, so a repository with `v1.2.3` and `1.3.0` selects the
   second.
3. **An ambiguous highest version refuses.** When the highest reachable version
   exists under more than one ref spelling, the command refuses in preflight,
   names both refs and names the existing selector that resolves it. It proposes
   no version and approves nothing.
4. **The proposal keeps the selected spelling.** The next version, its approval
   output and reset planning all use the spelling of the tag selection resolved
   to.
5. **The inventory keeps refs apart.** Local and remote refs are inventoried
   without collapsing semantic aliases, Releases are matched by exact tag name,
   and the reset digest retains each entry's spelling, ref and commit.
6. **Public states and exit codes are unchanged.** Ambiguity is a typed
   preflight refusal, not a new readiness state and not an approval.

## User Experience

A maintainer of a repository tagged `1.4.2` runs the command and reads a plan
whose base is `1.4.2` and whose proposal is `1.4.3`. A maintainer whose
repository carries both `v2.0.0` and `2.0.0` reads a refusal naming both refs
and the selector that picks one. Everyone else sees what they saw before.

## Non-Goals / Out of Scope

- Changing this repository's own release workflow, its `v*` tag trigger, or the
  spelling it publishes.
- Accepting prerelease identifiers, build metadata, or any tag shape the command
  rejects today.
- Normalizing, rewriting or deleting anyone's tags.
- Changing the states, exit codes, approval questions or digests the command
  already produces.
- Teaching the publication pipeline to publish a bare-tagged release.

## Declared intentional breaks

- A repository whose tags are bare stops being told its history is malformed and
  starts receiving plans. That is the feature.
- A repository carrying both spellings of its highest version stops receiving a
  plan and starts receiving a refusal. It was previously served by ignoring one
  of its refs.

## Regression locks

- A repository tagged only with the prefix produces the same plan, the same
  proposal, the same states and the same exit codes as today.
- Malformed, prerelease and build-metadata inputs keep their current rejections
  and their current messages.
- The command stays read-only: no Run, no config read, no mutation of refs,
  releases or files.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: this
repository's own tag history, which is entirely prefixed and includes a dead tag
beside a published one, and the command's published guidance, which states the
accepted shape. The replay plans against this repository's real tags and shows
an unchanged result, then plans against a fixture whose tags are bare and shows
the plan the current command refuses to produce. Where the repository's tags
cannot be read, the row records that reason and does not block.

## Success Metrics

1. A fixture repository tagged `1.4.2` receives a plan whose base and proposal
   are both bare, where the current command reports malformed input.
2. This repository's own tags produce the same plan, state and exit code as
   before the change.
3. A fixture carrying both `v2.0.0` and `2.0.0` receives a preflight refusal
   naming both refs and the selector, with no proposed version.

## Decisions

- **Refuse the ambiguity, do not resolve it.** Two refs for one version is a
  repository's decision to make; a planner that picks silently makes its output
  depend on ref creation order.
- **Preserve the selected spelling.** A proposal a maintainer has to translate
  before using is a proposal that invites a typo in a tag name.
- **Nothing changes for this repository.** Its tags are prefixed and its
  workflow triggers on `v*`; this Spec teaches the planner, not the pipeline.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
semantic version tag prefixes and on tooling that must read more than one
convention. The results were dominated by mirrors of this repository, which are
references rather than independent knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring; none is a
source for this Spec.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The parser's prefix requirement, its malformed-input
message, the prerelease and build-metadata rejections, the reset digest's
contents and this repository's own tag list were each read before this PRD was
written.

## Open Questions

None.
