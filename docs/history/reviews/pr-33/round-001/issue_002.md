---
source: coderabbit
pr: "33"
round: 1
round_created_at: "2026-07-21T17:02:02Z"
status: invalid
terminal_reason: "The cited rerun builder is for a same-runtime model probe; configured cross-runtime profile fallback uses agentSessionOwner and its regression test passes."
head_repository: marcioaltoe/roundfix
head_branch: ma/0041-agent-selection-runtime-readiness
head_sha: 6b48b67ab2154bb40d396befd673ad645a528214
file: internal/cli/selection.go
line: 215
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6SqIb6,comment:PRRC_kwDOS0qyts7YA3ST
review_hash: eef74b339a66922e9565e2540cb69abb0db0b58bf0544f9598ef8677131755fd
duplicate_of: ""
source_review_id: "4747041240"
source_review_submitted_at: "2026-07-21T17:01:00Z"
---

# Issue 002: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Use the fallback runtime in the rerun command.**

Line 212 reuses the original request runtime. A profile fallback may change runtime (for example, `claude` preferred to `codex` fallback), so this command can rerun the rejected selection. Populate `--agent` from the selected fallback runtime and add a cross-runtime fallback regression test.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/selection.go` around lines 211 - 215, Update the rerun command
argument construction in the selection flow to populate --agent from the
selected fallback runtime rather than req.agent, while preserving the fallback
model and reasoning-effort values. Add a regression test covering a profile
fallback that changes runtime, such as claude to codex, and verify the rerun
uses the fallback runtime.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:de74407347036b32a3c8e588 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `fallbackRerunCommand` belongs to the legacy dynamic model-probe path,
  whose `agent.FallbackSelection` contains only model and reasoning effort and
  therefore remains on the request runtime. Configured profile fallback is a
  separate daemon path: `agentSessionOwner.activate` passes each fallback's
  complete runtime/model/reasoning selection to the runtime factory. The
  existing cross-runtime regression proves a Codex preferred selection falls
  back to a Claude session.
- Focused check: `rtk go test ./internal/daemon -run '^TestCrossRuntimeFallbackUsesRuntimeFactory$'` — passed (1 test in 1 package).
