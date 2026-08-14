---
source: coderabbit
pr: "19"
round: 1
round_created_at: "2026-07-07T22:44:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/run-browser
head_sha: f2726493dff5e63e604139d27d147973ff650cf5
file: internal/cli/cli.go
line: 55
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6PCMCi,comment:PRRC_kwDOS0qyts7S9VwD
review_hash: f1ff2dce621366965f6d9122b0809c91c4fa6891437903d28f315432a5f7238c
duplicate_of: ""
source_review_id: "4648487653"
source_review_submitted_at: "2026-07-07T19:59:41Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Add the bare `roundfix runs` entry to top-level help.**

`runRunsCommand` supports `roundfix runs` in interactive terminals, and `commandUsage("runs")` documents that path, but `roundfix --help` only advertises `roundfix runs list ...` here. That makes the top-level usage incomplete for a supported command. As per coding guidelines, help text must be concise, truthful, and backed by implemented behavior.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` at line 55, The top-level help text is missing the
supported bare roundfix runs command, so update the CLI usage string in cli.go
to advertise roundfix runs alongside the existing list form. Make the help entry
concise and truthful by reflecting the implemented behavior in runRunsCommand
and keeping it aligned with commandUsage("runs"), so roundfix --help shows both
the interactive runs entry and the list subcommand.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8461730f63bf59a848c10610 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `runRunsCommand` supports the bare interactive `roundfix runs` command and `commandUsage("runs")` documents it, but the top-level usage did not list it.

## Resolution

- Status: `resolved`
- Changes: Added `roundfix runs` to top-level help and extended `TestRunHelp`.
- Verification: `rtk make verify` passed in this session for the configured `make verify` gate.
