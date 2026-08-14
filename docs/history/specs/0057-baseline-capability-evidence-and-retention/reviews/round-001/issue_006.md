---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "No active repository complexity rule requires this refactor; the cited gocyclo configuration is a Skill asset, not the repository Verification configuration."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/classification.go
line: 1035
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Yma,comment:PRRC_kwDOS0qyts7cjgE2
review_hash: 5fa6a5b05f45e05379897f7dc06d523adc920a634cb218718a230d0db4e6c3ee
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:29Z"
---

# Issue 006: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Report the configured complexity limits and estimate function complexity.
set -euo pipefail

fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec rg -n -C 3 'cyclop|gocyclo|gocognit|funlen|complexity' {}

ast-grep outline internal/baseline/classification.go --items all --type function
```

Repository: marcioaltoe/roundfix

Length of output: 3886

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- configured complexity rules ---'
fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec sed -n '105,130p' {}

printf '%s\n' '--- target function ---'
sed -n '971,1035p' internal/baseline/classification.go

printf '%s\n' '--- available complexity tooling ---'
command -v golangci-lint || true
command -v gocyclo || true
go version 2>/dev/null || true
```

Repository: marcioaltoe/roundfix

Length of output: 2701

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path

path = Path("internal/baseline/classification.go")
text = path.read_text()
start = text.index("func classifyMarkedCarrier(")
end = text.index("\n}\n", start) + 2
fn = text[start:end]

# Approximate gocyclo's decision-point model for this isolated function:
# base 1, plus if/for/range/switch/select/case/comm clauses, and boolean
# operators in conditions.
import re
decisions = []
for m in re.finditer(r'\b(if|for|range|switch|select|case|default)\b|&&|\|\|', fn):
    decisions.append(m.group())
print(fn)
print("decision tokens:", decisions)
print("estimated cyclomatic complexity:", 1 + len(decisions))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 2191

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
from pathlib import Path
import re

text = Path("internal/baseline/classification.go").read_text()
start = text.index("func classifyMarkedCarrier(")
end = text.index("\n}\n", start) + 2
fn = text[start:end]

decisions = re.findall(r'\b(if|for|range|switch|select|case|default)\b|&&|\|\|', fn)
print(fn)
print("decision tokens:", decisions)
print("estimated cyclomatic complexity:", 1 + len(decisions))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 2183

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

golangci-lint run \
  --disable-all \
  --enable=gocyclo \
  --config=.golangci.yml \
  ./internal/baseline
```

Repository: marcioaltoe/roundfix

Length of output: 264

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

golangci-lint version
golangci-lint run --help | sed -n '1,120p'
```

Repository: marcioaltoe/roundfix

Length of output: 5751

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

golangci-lint run \
  --config=.golangci.yml \
  ./internal/baseline
```

Repository: marcioaltoe/roundfix

Length of output: 388

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

config="$(fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' | head -n 1)"
test -n "$config"
printf 'config=%s\n' "$config"

golangci-lint run \
  --config="$config" \
  --enable-only=gocyclo \
  ./internal/baseline
```

Repository: marcioaltoe/roundfix

Length of output: 26319

---



</details>

**Reduce `classifyMarkedCarrier` complexity below 13.**

`gocyclo` reports complexity 19. Extract per-entry ownership evaluation and managed-set completeness checks into helpers.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/classification.go` around lines 971 - 1035, Reduce
classifyMarkedCarrier’s cyclomatic complexity below 13 by extracting the
per-entry ownership and staleness logic into a helper, and extracting the
expected-versus-seen managed-ID completeness check into another helper. Preserve
all existing classification outcomes, including setup-context markers, unmanaged
entries, manifest validation, current-byte matching, and stale managed
detection.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:30438b3dd9b8e79cd4acb947 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The active gate is `make verify`, whose test/build flow does not run gocyclo, and the repository defines no complexity threshold of 13. Refactoring this direct classification flow solely for an inactive metric would violate the bounded KISS and surgical-change rules.
