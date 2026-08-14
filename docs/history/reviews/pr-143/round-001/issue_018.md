---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: .agents/skills/roundfix/SKILL.md
line: 466
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKK6,comment:PRRC_kwDOS0qyts7e9kiu
review_hash: d2b972853794a73640a74d53d6ed36c5f76905ed3f0b8effb8bf84cd1cf1ffde
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---


# Issue 018: _ Functional Correctness_ _ Trivial_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

rg -n -C 6 \
  'opus\[1m\]|bracketed|context suffix|canonical prefix|reasoning_effort' \
  internal/agent internal/config .roundfixrc.yml .agents/skills/roundfix/SKILL.md
```

Repository: marcioaltoe/roundfix

Length of output: 50379

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- relevant documentation ---'
sed -n '437,458p' .agents/skills/roundfix/SKILL.md

printf '%s\n' '--- capability parser and selection logic ---'
rg -n -C 8 \
  'type ModelCapability|parseModelCapability|CanonicalModel|AdapterValue|PlanSelectionAssignment|SelectionEncodingIndependent|ModelNotAdvertised' \
  internal/agent/selection_capabilities.go internal/agent/acpx_runner.go

printf '%s\n' '--- profile validation and model lookup ---'
rg -n -C 8 \
  'ModelCatalog|claudeModelCatalog|CanonicalModel|AdapterModel|model.*catalog|catalog.*model|Validate.*Profile|validate.*model' \
  internal/agent internal/config \
  -g '*.go' | head -n 500
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- assignment planner ---'
sed -n '240,320p' internal/agent/selection_assignment.go
sed -n '380,490p' internal/agent/selection_assignment.go

printf '%s\n' '--- profile/model validation call sites ---'
rg -n -C 6 \
  'ModelCatalog\(|catalog|ModelChoice|RuntimeSpec.*Model|Model.*allow|allowlist|profiles.*validate|validate.*profiles|Validate.*Runtime' \
  internal -g '*.go' | grep -v '_test.go' | head -n 400

printf '%s\n' '--- relevant config and runtime types ---'
rg -n -C 5 \
  'type RuntimeSpec|type Profile|ReasoningEffort|Model string|model:' \
  internal/config internal/agent -g '*.go' | head -n 400
```

Repository: marcioaltoe/roundfix

Length of output: 30446

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- profile validation implementation ---'
rg -n 'func validateAgentSelectionProfile|func validate.*Selection|validateAgentSelectionProfile' internal/config -g '*.go'
line=$(rg -n 'func validateAgentSelectionProfile' internal/config -g '*.go' | cut -d: -f1 | head -n1)
file=$(rg -l 'func validateAgentSelectionProfile' internal/config -g '*.go' | head -n1)
start=${line:-1}
end=$((start + 100))
sed -n "${start},${end}p" "$file"

printf '%s\n' '--- runtime selection construction ---'
rg -n -C 10 \
  'RuntimeSpec\{|Model:.*profile|profile.*Model|ReasoningEffort:.*profile|resolve.*profile|effective.*profile' \
  internal -g '*.go' | head -n 300

printf '%s\n' '--- focused tests for canonical and advertised spellings ---'
sed -n '270,345p' internal/agent/selection_capabilities_test.go
sed -n '400,440p' internal/agent/selection_capabilities_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 980

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- profile validation implementation ---'
sed -n '360,430p' internal/config/profiles.go

printf '%s\n' '--- runtime selection construction ---'
rg -n -C 10 \
  'RuntimeSpec\{|Model:.*profile|profile.*Model|ReasoningEffort:.*profile|resolve.*profile|effective.*profile' \
  internal -g '*.go' | head -n 300 || true

printf '%s\n' '--- focused tests for canonical and advertised spellings ---'
sed -n '270,345p' internal/agent/selection_capabilities_test.go
sed -n '400,440p' internal/agent/selection_capabilities_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 27508

---



</details>

**Document the canonical and adapter model namespaces.**

State that `model: opus` is the canonical profile value. The planner maps it to the advertised adapter value, such as `opus[1m]`, while `model: opus[1m]` selects that adapter value unchanged. Clarify that exact advertised values take precedence when both forms are advertised.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 25: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[error] 1677: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1677: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 1153: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 441 - 444, Update the model
namespace documentation near the listed Codex and Claude identifiers to state
that opus is the canonical profile value, the planner maps it to the advertised
adapter value such as opus[1m], and opus[1m] is passed through unchanged; also
state that exact advertised values take precedence when both canonical and
adapter forms are available.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:2c6c1df263b03d9bce5ab283 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: SKILL.md has strict byte-identical contract tests (TestAuthorialSkillSync, TestBaselineSkillContract). Edit reverted — the existing documentation at lines 460-479 already covers canonical adapter namespace mapping, opus-1m bracketed suffix, and catalog value semantics. Adding new text would require syncing embedded distributed copies via make skills-sync.
