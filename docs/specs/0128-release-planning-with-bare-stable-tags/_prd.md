---
spec: 0128-release-planning-with-bare-stable-tags
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# Release planning preserves stable tag identity

The Release Plan currently rejects repositories whose stable versions are tagged without v. Extend the read-only planner to accept both stable spellings, preserve the selected ref and its proposal convention, and make ambiguity explicit.

This Spec is in authoring. Its implementation intent is selected by the maintainer's
complete-queue request; the exact governed mutations and execution limits remain
proposed until their operative grant is approved.

## Project Constraints

- Identifier strategy: applicable — preserve Spec/Task identities and exact Git ref identity; any new recovery record follows the existing Run identity boundary. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential or HTTP API change is proposed; use existing read boundaries only. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0143 requires read-only ordinary/reset planning; ADR-0062 preserves operational history during identity-reset planning. ADR-0146 publication is outside this Spec. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `internal/docscontract/publicdocs_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Consumers using bare stable tags can use the existing planner.
- Existing v-prefixed results and release-classification rules remain unchanged.
- Every selected/reset ref preserves its exact identity and every operation remains read-only.

## Core Features

1. Accept strict vMAJOR.MINOR.PATCH and MAJOR.MINOR.PATCH; preserve current rejection of malformed, prerelease and build-metadata inputs.
2. Select the highest reachable semantic version across both spellings. If that highest version has more than one ref spelling, require the existing --from selector and report both exact names and commits. Explicit selection retains existing range/ancestry checks.
3. Preserve the exact selected tag and propose the next version in that same spelling, including approval output and reset planning.
4. Inventory local and remote refs without collapsing semantic aliases; match Releases by exact tag name and retain spelling/ref/commit in the reset digest.
5. Keep public states and exit codes: ambiguity is a typed preflight refusal, not release approval or a new readiness state.

## Non-Goals / Out of Scope

- Publishing or mutating releases, tags, remotes, operational config, changelogs or package versions.
- Changing this repository's v-prefixed release workflow.
- Adding another tag-convention configuration when the existing --from selector resolves ambiguity.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
later Task Graph. At least one acceptance row replays the pre-existing fleet
observation identified by the source ownership index; a fixture authored by this
Spec alone is insufficient to validate the original premise.

## Decisions and open authority

The maintainer selected the remaining source intent for the implementation queue.
The technical candidate is reviewable; proposed paths, experimental limits and
any change to an accepted ADR must receive an operative decision before execution.
The existing through-merge authority requires the configured review policy outcome and required
checks and grants no release or waiver.

## Research basis

Local source inspection and the archived fleet observations establish the
remaining behavior. Secondbrain's agent-workflow and harness syntheses explain
why durable state, external acceptance and explicit stop conditions matter.
Exa-read official Codex noninteractive and Claude headless documentation supports
independent execution interfaces, not the correctness of this Roundfix design.
The [historical research record](https://github.com/marcioaltoe/roundfix/blob/6b8ea48725cbca13974eee0b400b3482202874f6/docs/workflow/2026-09-08-pending-work-plan.md) preserves source URLs, their influence and remaining limitations at its cited Git revision.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
