---
source: coderabbit
pr: "19"
round: 1
round_created_at: "2026-07-07T22:44:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/run-browser
head_sha: f2726493dff5e63e604139d27d147973ff650cf5
file: internal/tui/styles.go
line: 96
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6PCMCs,comment:PRRC_kwDOS0qyts7S9VwR
review_hash: ade3067d0cd799aa35cf1350a7586aef346d7bbaf06a8a53fe676219b7d4f969
duplicate_of: ""
source_review_id: "4648487653"
source_review_submitted_at: "2026-07-07T19:59:41Z"
---

# Issue 004: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'FILES\n'
git ls-files 'internal/tui/styles.go' 'internal/tui/cockpit.go' 'internal/tui/runbrowser.go' 'SKILL.md'

printf '\nOUTLINE styles.go\n'
ast-grep outline internal/tui/styles.go --view expanded || true

printf '\nOUTLINE cockpit.go\n'
ast-grep outline internal/tui/cockpit.go --view expanded || true

printf '\nOUTLINE runbrowser.go\n'
ast-grep outline internal/tui/runbrowser.go --view expanded || true

printf '\nRELEVANT SEARCHES\n'
rg -n --context 2 'cockpitTokens|ResolveTokens\(|ColorEnabled|NO_COLOR|ROUNDFIX_COLOR|color-mode|color mode|Live Run View|Run Browser' internal/tui SKILL.md || true
```

Repository: marcioaltoe/roundfix

Length of output: 27435

---



</details>

**Run Browser still bypasses color-mode tokens** `internal/tui/runbrowser.go:137-186` still renders header/footer/rows through the package-level `cockpitTokens`, so it ignores `ColorEnabled` and will keep emitting the styled palette when color is disabled. Thread per-instance tokens into `RunBrowser` before removing the legacy package var.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 89-89: var styleBright is unused

(unused)

---

[error] 91-91: var styleBar is unused

(unused)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/styles.go` around lines 76 - 96, RunBrowser is still using the
package-level cockpitTokens, so it bypasses the instance ColorEnabled setting
and keeps styling enabled when color should be off. Update RunBrowser to accept
and use per-instance tokens instead of the global cockpitTokens, then thread
those tokens in from the caller so header/footer/row rendering respects the
current color mode. Keep the legacy style aliases in styles.go only as fallback
wiring until the per-instance path is fully in place.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d7b25a7447e7bd917cef9e61 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Run Browser header, footer, and Active-row rendering used the package-level styled `cockpitTokens`, bypassing the caller's effective color mode.

## Resolution

- Status: `resolved`
- Changes: Added per-instance Run Browser tokens, threaded them from the CLI with `ResolveTokens(colorEnabled(stdout))`, and added a no-color rendering regression test.
- Verification: `rtk make verify` passed in this session for the configured `make verify` gate.
