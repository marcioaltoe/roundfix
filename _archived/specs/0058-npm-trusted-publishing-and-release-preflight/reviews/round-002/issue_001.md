---
source: coderabbit
pr: "61"
round: 2
round_created_at: "2026-08-01T13:51:04Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/npm-trusted-publishing-and-release-preflight
head_sha: b540a477ef11b1ddd09462656f6dab85bdd4affc
file: .github/workflows/release.yml
line: 16
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vogyi,comment:PRRC_kwDOS0qyts7cSDPj
review_hash: c8b9ad65f3e54018d38713a74f6248be8c7a080204b655070d25f50e332cfdba
duplicate_of: ""
source_review_id: "4834755579"
source_review_submitted_at: "2026-08-01T13:42:06Z"
---

# Issue 001: _ Security & Privacy_ _ Trivial_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🔵 Trivial_ | _⚡ Quick win_

**Scope `id-token`/`contents` permissions to the job, not the workflow.**

zizmor flags both `id-token: write` and `contents: write` as overly broad at the workflow level. The workflow currently has one job, so this has no functional effect today, but workflow-level permissions apply to every job added in the future without an explicit opt-in. Move these into `permissions:` under the `release` job instead.

<details>
<summary>♻️ Proposed fix</summary>

```diff
-permissions:
-  id-token: write # mint the npm Trusted Publishing OIDC token
-  contents: write # create the GitHub Release and upload binary assets
-
 jobs:
   release:
     runs-on: ubuntu-latest
+    permissions:
+      id-token: write # mint the npm Trusted Publishing OIDC token
+      contents: write # create the GitHub Release and upload binary assets
     env:
```
</details>

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 zizmor (1.28.0)</summary>

[error] 15-15: overly broad permissions (excessive-permissions): id-token: write is overly broad at the workflow level

(excessive-permissions)

---

[error] 16-16: overly broad permissions (excessive-permissions): contents: write is overly broad at the workflow level

(excessive-permissions)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.github/workflows/release.yml around lines 14 - 16, Move the id-token and
contents write permissions from the workflow-level permissions block into a
permissions block within the release job, preserving both required scopes and
leaving other workflow behavior unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:9ebdb9be38c337674f0ae807 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The reviewed head declared both write permissions at workflow scope. The release job now owns the unchanged `id-token: write` and `contents: write` grants, so future jobs do not inherit them. A focused Ruby YAML assertion confirmed the workflow has no top-level permission map and the release job has exactly those two grants; `rtk git diff --check` exited 0. The Daemon owns authoritative Verification.
