---
source: coderabbit
pr: "61"
round: 1
round_created_at: "2026-08-01T13:42:52Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/npm-trusted-publishing-and-release-preflight
head_sha: b540a477ef11b1ddd09462656f6dab85bdd4affc
file: .github/workflows/release.yml
line: 16
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vogyi,comment:PRRC_kwDOS0qyts7cSDPj
review_hash: 04fffc800811e23d6403c2699f655a8256cd4e1dccb6b06521dbe8ef05746bc3
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_001.md
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

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `git show ma/npm-trusted-publishing-and-release-preflight:.github/workflows/release.yml` confirms both write permissions are declared at workflow scope on the reviewed head. The fix was not applied because this Batch checkout is on unrelated branch `ma/baseline-digest-regeneration-bootstrap` at `0cc8bf3`; switching branches is outside the Agent-owned Batch lifecycle, and editing the visible workflow would patch the wrong revision.
