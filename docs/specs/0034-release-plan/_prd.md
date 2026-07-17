---
spec: 0034-release-plan
status: active
created: 2026-07-16
surfaces: [backend, cli, docs]
---

# Release Plan

Roundfix validates release tags and publishes matching artifacts, but it does not determine which semantic-version increment the committed changes require or identify when that decision needs explicit human approval. The Release Plan gives maintainers and Agents one deterministic, read-only assessment before any changelog, tag, push, package, or GitHub Release mutation.

## Goals

- Classify every committed change in a release range by the highest consumer impact and propose the minimum valid next version.
- Distinguish a conclusive plan from a change set that needs manual impact classification instead of guessing from incomplete evidence.
- Make minor, major, and version-zero breaking decisions explicit approval gates while keeping patch and no-release outcomes low-friction.
- Give humans and Agents stable text and machine-readable results without changing repository or external state.
- Make the Release Plan the mandatory first step in the maintainer and Agent release instructions.

## User Stories

1. As a maintainer preparing a release, I want the committed delta classified against the latest release, so that I can propose a version supported by evidence.
2. As an Agent asked to cut the next release, I want a read-only command to identify the required increment and approval boundary, so that I do not make an unauthorized version decision.
3. As a maintainer reviewing ambiguous internal or consumer-facing changes, I want the plan to stop for manual classification, so that an uncertain change is never silently under-versioned.
4. As an automation caller, I want a stable machine-readable Release Plan and exit contract, so that I can distinguish ready, approval-required, manual-review, and invalid-input outcomes.
5. As a maintainer with a documentation-, test-, or CI-only delta, I want the plan to report that no release is required, so that repository maintenance does not force an empty package release.

## Core Features

1. The Release Plan Command analyzes a committed base release through a committed target revision and reports the base version, target identity, classified changes, highest impact, proposed version, and decision state.
2. Breaking changes outrank compatible features, compatible features outrank fixes, fixes outrank no-release changes, and commit order never changes that result.
3. Conventional Commit evidence provides deterministic minimum impact: breaking markers indicate breaking impact, compatible features indicate minor impact, and compatible fixes indicate patch impact. Other change types do not receive an automatic release meaning.
4. Changes without sufficient automatic evidence either qualify for no release under the documented maintenance-only boundary or require explicit manual impact classification. The command never guesses that an ambiguous consumer-facing change is safe.
5. While the project remains below `1.0.0`, a breaking change proposes the next minor version, labels the plan breaking, and requires explicit approval. At or above `1.0.0`, a breaking change proposes the next major version and requires explicit approval.
6. A generic request to cut a release authorizes a conclusive patch increment but not a minor, major, or breaking increment. A user-specified target version counts as approval only when it is at least the minimum version required by the plan.
7. The command is read-only on every outcome. It never edits the changelog or version-bearing files, creates or pushes tags, publishes packages, creates a GitHub Release, or contacts an external service.
8. Human-readable output summarizes the decision and next action. Machine-readable output uses a versioned schema and includes the evidence needed to reproduce the classification.
9. The release runbook and root Agent instructions require a Release Plan before release mutation and preserve the existing tag-driven publication workflow after approval.

## User Experience

The default flow starts from the latest reachable stable release and ends at the selected committed revision. A conclusive patch plan states the proposed version and that no additional version approval is required. A minor, major, or version-zero breaking plan states the exact approval question and remains read-only. An ambiguous plan names the commits that need manual impact classification and the explicit information required to continue. Invalid repositories, revisions, ranges, or manual classifications fail with one actionable diagnostic.

## Non-Goals / Out of Scope

- Editing `CHANGELOG.md`, package metadata, or any other release file.
- Creating or pushing a tag, publishing npm packages, uploading assets, or creating a GitHub Release.
- Monitoring the tag-triggered release workflow.
- Using an Agent, network service, or hidden heuristic to infer semantic meaning from source code.
- Replacing the existing tag validation, artifact version agreement, verification gate, or publication sequence.
- Deciding whether Roundfix is ready to move from version zero to `1.0.0`.
- Supporting pre-release version selection in the first release.

## Success Metrics

- From `v0.4.0`, a fix-only delta produces a patch Release Plan for `v0.4.1` without an additional version-approval requirement.
- From `v0.4.1`, a compatible feature produces an approval-required minor Release Plan for `v0.5.0`.
- From `v1.4.2`, a breaking delta produces an approval-required major Release Plan for `v2.0.0`.
- From a `v0.y.z` release, a breaking delta produces the next minor version, carries a breaking label, and requires approval.
- Mixed deltas always select the highest required increment regardless of commit order.
- Documentation-, test-, and CI-only deltas can produce a no-release plan.
- Ambiguous consumer-facing changes produce a manual-classification-required result and never a lower inferred increment.
- Every command outcome leaves repository bytes, refs, remotes, packages, and releases unchanged.

## Decisions

- Release planning is read-only and confirmation-gated. See ADR-0048.
- The first release includes both the Release Plan Command and mandatory release instructions; a documentation-only policy is insufficient for deterministic automation.
- Automatic classification is conservative: missing semantic evidence becomes an explicit manual decision, not a guessed patch or no-release result.

## Open Questions

None.
