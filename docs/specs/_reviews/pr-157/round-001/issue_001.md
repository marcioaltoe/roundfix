---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: .agents/skills/archive-spec/SKILL.md
line: 3
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdC,comment:PRRC_kwDOS0qyts7gSdxe
review_hash: 20acf12c6e52a0e543316e9f2a35c9aa6ac0c6c1dec04fabe410e2d504820340
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Resolve the configured Spec Root before moving the Spec.**

These instructions hard-code `docs/specs/<slug>` as the source. The PR supports external Spec Roots, so this move fails when the active Spec is outside `docs/specs`. Resolve the source from `specs.root` and resolve the destination through the shared `ArchiveDir(ArchiveKindSpec)` contract.

<details>
<summary>Proposed documentation correction</summary>

```diff
- Move `docs/specs/<slug>/` to `_archived/specs/<slug>/`.
+ Move `<resolved-specs-root>/<slug>/` to the repository archive root `_archived/specs/<slug>/`.

- git mv docs/specs/<slug> _archived/specs/<slug>
+ Resolve the configured Spec Root and repository archive root before running the move.
```
</details>







Also applies to: 16-16, 220-221

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md at line 3, Update the archive workflow
instructions to resolve the active Spec source from the configured specs.root
instead of hard-coding docs/specs/<slug>. Resolve the archive destination using
the shared ArchiveDir(ArchiveKindSpec) contract, and apply this consistently to
the corresponding source and destination references throughout the document.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c0dbbd6831e2428d9e384030 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated `.agents/skills/archive-spec/SKILL.md` (and its `skills/` mirror via `make skills-sync`) so the archive workflow resolves the source from the configured Spec Root instead of hard-coding `docs/specs/<slug>`, and resolves the destination beside the active root (`<spec-root>/_archived/`) or the default `_archived/specs/`. The description, the intro paragraph, and the Step 2 move command now name both the built-in `docs/specs` root and a non-default/external root, matching `spec.ArchiveSpecRoot`. Verified: `make skills-sync-check` and `roundfix skills check` pass; the full `make verify` gate passes.
