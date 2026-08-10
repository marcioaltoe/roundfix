---
granted: 2026-08-10
action: correct-skill-adapter-guidance
paths:
  - .agents/skills/roundfix/SKILL.md
consuming: direct
---

# Tooling authorization — the skill stops promising a diagnosis it lost (2026-08-10)

On 2026-08-10 the maintainer directed the removal of the unofficial adapter
lineage:

> Remova o legado.

Asked afterwards what remained before the branch could close, the maintainer
approved the correction this record covers:

> sim

## What this covers

Removing `@zed-industries` package detection left four documents promising a
diagnosis the binary no longer performs. The measured behavior change is
narrow, and worth stating precisely so the correction does not overstate it
either:

- **What still works.** `internal/cli/setup.go:201` proposes an adapter
  migration whenever lineage proof fails, not when a specific package is
  recognized. An override resolving to `@zed-industries/codex-acp` still fails
  proof, still produces one migration offer to the official pinned package, and
  still asks before writing.
- **What was lost.** The error names the required package instead of the
  installed one. `did not prove required package lineage
  @agentclientprotocol/codex-acp` replaces `reported legacy package
  @zed-industries/codex-acp`, because an unrecognized package no longer
  resolves to a name.

So the guidance is not wrong about the offer; it is wrong about the naming. The
skill and the three user guides are corrected to describe lineage proof and its
migration offer without claiming that any particular superseded package is
recognized by name.

`TestProfilesDocumentationContractMatchesPublicGuidance` pins the removed
package strings in those documents. It changes in the same commit, because the
contract it enforces is that the documentation matches shipped behavior, and
the behavior is what moved.

## Authorized paths

- `.agents/skills/roundfix/SKILL.md`, limited to the adapter lineage and
  migration passages that name superseded packages.

The generated copy under `skills/roundfix/SKILL.md`, rewritten by
`make skills-sync`, is sanctioned fallout under ADR-0081 rather than a separate
target. The three files under `docs/user-guide/` are product documentation, and
`internal/cli/cli_test.go` is a test, so neither is repository tooling.

## Bounded by purpose

This grant covers making the skill's adapter guidance match the shipped
diagnosis. It does not authorize changing what Setup or Doctor do, restoring
the removed lineage, or editing any other passage of the skill.

## Consuming Spec

Applied directly rather than through a Spec: it corrects prose that a change
already committed on this branch made inaccurate.

## Commit choreography

This record lands as its own commit, before the commit that corrects the skill.
It is a consequent correction — the change that made the prose stale is already
in `b7b97c4b` — so it lands after that commit, never before it.
