---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/cli/implement_test.go
line: 5715
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWjK,comment:PRRC_kwDOS0qyts7fswOz
review_hash: d3b3544fef4d6c659796d242c68c85c9ccecd6c8dfca6b07e15f50662be07d03
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 011: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: List every model named by CLI test profiles and compare with the fake catalogues.
set -euo pipefail

echo '== models named in macroProfilesYAML and related fixtures =='
rg -nP --type=go -C25 'func macroProfilesYAML\(' internal/cli

echo '== runtimes configured in tests that use fakeACPXCommand =='
rg -nP --type=go -C3 'fakeACPXCommand\(' internal/cli

echo '== every model string in the CLI test fixtures =='
rg -nP --type=go -o '\bmodel:\s*"[^"]+"|"macro-[a-z-]+"|"gpt-5\.[0-9a-z.-]+"|"claude-[a-z0-9-]+"' internal/cli | sort -u
```

Repository: marcioaltoe/roundfix

Length of output: 10495

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== catalogue declarations and agent selection =='
sed -n '5688,5722p' internal/cli/implement_test.go
sed -n '5860,5905p' internal/cli/implement_test.go

echo '== macro profile models =='
sed -n '5948,6045p' internal/cli/implement_test.go

echo '== smaller fake catalogue and its callers =='
sed -n '150,225p' internal/cli/implement_test.go
sed -n '3360,3400p' internal/cli/implement_test.go

echo '== preflight model validation implementation =='
rg -n -C8 'Preflight|preflight|advertis|catalog|model.*(available|supported|known)|available.*model' internal/cli --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 50377

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== requested source ranges =='
sed -n '5688,5722p' internal/cli/implement_test.go
sed -n '5868,5898p' internal/cli/implement_test.go
sed -n '5948,6045p' internal/cli/implement_test.go
sed -n '154,220p' internal/cli/implement_test.go

echo '== model catalogue and agent-selection symbols =='
rg -n -C5 'CODEX_MODELS|CLAUDE_MODELS|AGENTS|func .*Models|agent.*model|models.*agent|advertis' internal/cli/implement_test.go

echo '== profile model literals in implement_test.go =='
python3 - <<'PY'
import re
from pathlib import Path
p = Path("internal/cli/implement_test.go")
text = p.read_text()
for name, body in re.findall(r'func (macroProfilesYAML)\(\) string \{(.*?)\n\}', text, re.S):
    print("macroProfilesYAML models:")
    for model in re.findall(r'^\s+model:\s*([^\s]+)', body, re.M):
        print(model)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 12178

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

echo '== all fakeACPXCommand call sites and nearby model selections =='
python3 - <<'PY'
from pathlib import Path
p = Path("internal/cli/implement_test.go")
lines = p.read_text().splitlines()
for i, line in enumerate(lines, 1):
    if "fakeACPXCommand(" in line:
        print(f"\n--- call at {i} ---")
        for n in range(max(1, i-8), min(len(lines), i+14)+1):
            print(f"{n}:{lines[n-1]}")
PY

echo '== macro fake catalogue response logic =='
sed -n '5880,5950p' internal/cli/implement_test.go

echo '== production advertised-model membership checks =='
rg -n -C6 'not advertised|advertised.*model|current_model_id|config_options|options.*model|Model.*advertis|advertis.*Model' --glob '*.go' .
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---



</details>

**Select the model catalogue by agent explicitly.**

The configured macro models are already present in the correct catalogues. However, `AGENTS` includes `opencode`, while `sessions show` returns `CODEX_MODELS` for every non-`claude` agent. Use an explicit agent-to-catalogue map and fail loudly for unmapped agents.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement_test.go` around lines 5704 - 5715, Update the
sessions-show model selection logic to map each agent explicitly to its
catalogue, including the existing CODEX_MODELS, CLAUDE_MODELS, and opencode
catalogue. Remove the fallback that returns CODEX_MODELS for every non-claude
agent, and fail loudly when an agent has no mapping. Use the existing AGENTS and
catalogue symbols to keep configured model assignments unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:48cb00343e0825eec7c609cd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: The macro fake ACPX `sessions show` handler in internal/cli/implement_test.go now selects the model catalogue via an explicit `CATALOGUES` map (codex/claude/opencode) instead of returning `CODEX_MODELS` for every non-claude agent, and refuses loudly (stderr + exit 2) when an agent has no mapping. Added `OPENCODE_MODELS` so opencode is not silently served the codex catalogue. Focused: `go test ./internal/cli -run 'TestImplement|TestMacro|TestProfile'` ok.
