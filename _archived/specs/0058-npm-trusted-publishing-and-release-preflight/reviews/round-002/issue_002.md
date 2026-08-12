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
line: 250
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vogyj,comment:PRRC_kwDOS0qyts7cSDPl
review_hash: d053b39ebf5c9547ef74df2947da43e9ce4036e4f00324beb62878eed3ba422b
duplicate_of: ""
source_review_id: "4834755579"
source_review_submitted_at: "2026-08-01T13:42:06Z"
---

# Issue 002: _ Security & Privacy_ _ Major_ _ Quick win_

## Review Comment

_🔒 Security & Privacy_ | _🟠 Major_ | _⚡ Quick win_

**Pass `NPM_TOKEN` through `env:` instead of interpolating the secret directly into the script.**

`NODE_AUTH_TOKEN="${{ secrets.NPM_TOKEN }}"` expands the secret directly into the generated shell script at template time. zizmor flags this as a template-injection risk. Even though this specific value is a maintainer-controlled secret rather than attacker-controlled input, GitHub's own hardening guidance is to always route `${{ }}` expressions, including secrets, through `env:` and reference them as shell variables, so the literal value never appears as inline script text.

<details>
<summary>🔒 Proposed fix</summary>

```diff
   - name: Publish to npm
     if: ${{ github.event_name == 'push' }}
     env:
       FALLBACK_WINDOW: ${{ vars.NPM_TRUSTED_PUBLISHING_FALLBACK }}
+      NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
     run: |
@@
           echo "::warning::identity: $coordinate failed OIDC publish; retrying with the bounded token fallback"
-          ( cd "$pkgdir" && NODE_AUTH_TOKEN="${{ secrets.NPM_TOKEN }}" npm publish --access public ) || return 1
+          ( cd "$pkgdir" && NODE_AUTH_TOKEN="$NPM_TOKEN" npm publish --access public ) || return 1
```
</details>

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 zizmor (1.28.0)</summary>

[warning] 250-250: code injection via template expansion (template-injection): may expand into attacker-controllable code

(template-injection)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.github/workflows/release.yml at line 250, Update the release workflow step
containing the npm publish command to pass the NPM_TOKEN secret through the
step’s env configuration, then reference the resulting shell environment
variable in the command instead of interpolating the GitHub expression inline.
Preserve the existing package-directory, public-access, and failure behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4e77b5cd93f74c0d94bfc0b5 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The reviewed head interpolated `secrets.NPM_TOKEN` inside the generated shell script. The publish step now maps the secret through `env.NPM_TOKEN`, and the fallback command reads `$NPM_TOKEN` while preserving its package directory, public access, and failure behavior. A focused Ruby YAML assertion confirmed the environment mapping, shell-variable reference, and absence of the secret expression from the `run` block; `rtk git diff --check` exited 0. The Daemon owns authoritative Verification.
