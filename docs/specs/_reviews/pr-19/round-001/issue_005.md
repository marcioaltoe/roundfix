---
source: coderabbit
pr: "19"
round: 1
round_created_at: "2026-07-07T22:44:40Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/run-browser
head_sha: f2726493dff5e63e604139d27d147973ff650cf5
file: internal/tui/styles.go
line: 96
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6PCMCy,comment:PRRC_kwDOS0qyts7S9VwZ
review_hash: c50f35a5ff92ef33c63f112b83d2f0b9fbeeb039c212c4fe69ea47c70ab2839b
duplicate_of: ""
source_review_id: "4648487653"
source_review_submitted_at: "2026-07-07T19:59:41Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Unused vars `styleBright` and `styleBar` flagged by golangci-lint.**

Static analysis reports `styleBright` (Line 89) and `styleBar` (Line 91) as unused. Since they're not referenced by any renderer yet, this will fail lint in CI.




<details>
<summary>🧹 Remove or gate unused legacy styles</summary>

```diff
 var (
 	styleAccent       = cockpitTokens.SectionLabel
 	styleMuted        = cockpitTokens.Muted
 	styleError        = cockpitTokens.Failed
 	styleActiveBorder = cockpitTokens.ActiveBorder
-	styleBright       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
 	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
-	styleBar          = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Background(lipgloss.Color("238"))
 	styleBarFill      = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("27")).Bold(true)
 	styleBarRest      = lipgloss.NewStyle().Foreground(lipgloss.Color("17")).Background(lipgloss.Color("153"))
 	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
 	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
 )
```

If these are intentionally kept for near-term use in an upcoming task, a `//nolint:unused` with a tracking note would silence the lint failure without deleting the styles.
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
var (
	styleAccent       = cockpitTokens.SectionLabel
	styleMuted        = cockpitTokens.Muted
	styleError        = cockpitTokens.Failed
	styleActiveBorder = cockpitTokens.ActiveBorder
	styleTool         = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	styleBarFill      = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("27")).Bold(true)
	styleBarRest      = lipgloss.NewStyle().Foreground(lipgloss.Color("17")).Background(lipgloss.Color("153"))
	styleFooter       = lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Background(lipgloss.Color("234"))
	styleBorder       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("238")).Padding(0, 1)
)
```

</details>

<!-- suggestion_end -->

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

In `@internal/tui/styles.go` around lines 84 - 96, The style declarations in the
styles block include legacy variables that are not referenced anywhere, causing
golangci-lint unused warnings. Remove the unused style definitions for
styleBright and styleBar, or if they must be kept for imminent use, gate them
with an explicit nolint comment and a tracking note in internal/tui/styles.go so
the lint check passes.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c56bd574132db1c44aaea555 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `invalid`
- Notes: The finding is stale in this worktree. `styleBright` is referenced in `internal/tui/agent_live.go` by `renderAgentHeader` and `colorTimelineLines`; `styleBar` is referenced by `renderPipelineBar`. Removing them would break current code, and no code change is needed for this issue.

## Resolution

- Status: `invalid`
- Verification: `rtk make verify` passed in this session for the configured `make verify` gate.
