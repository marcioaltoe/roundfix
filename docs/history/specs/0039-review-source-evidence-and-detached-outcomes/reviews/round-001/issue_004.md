---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/cli.go
line: 3907
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpb,comment:PRRC_kwDOS0qyts7aUVDn
review_hash: abef06f4a7b6af192365f0c725ef8b0df9feafd617cd335c411e1cb2c871567e
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Confirm the configured complexity limit and current findings.
fd -H -t f '.golangci' -x cat {}
```

Repository: marcioaltoe/roundfix

Length of output: 7620

---

🏁 Script executed:

```shell
python3 - <<'PY'
from pathlib import Path
import re

path = Path("internal/cli/cli.go")
text = path.read_text()
start = text.index("func inheritReviewArtifactEvidence")
brace = text.index("{", start)
depth = 0
end = None
for i in range(brace, len(text)):
    ch = text[i]
    if ch == "{":
        depth += 1
    elif ch == "}":
        depth -= 1
        if depth == 0:
            end = i + 1
            break

func = text[start:end]
print(func)

# crude cyclomatic complexity approximation for Go:
# 1 base + if/for/range/case + short-circuit operators in conditions
tokens = re.findall(r'\b(if|for|range|case|switch|&&|\|\|)\b', func)
complexity = 1
complexity += sum(1 for t in tokens if t in {"if", "for", "range", "case", "switch"})
complexity += sum(1 for t in tokens if t in {"&&", "||"})
print("\nAPPROX_COMPLEXITY:", complexity)
print("CONTROL_TOKENS:", tokens)
PY
```

Repository: marcioaltoe/roundfix

Length of output: 3511

---



</details>

**Split this function to get under the complexity limit.**

The current body is about 16 decision points, which exceeds the repo’s `gocyclo` threshold of 13. Extract the commit-shape checks and path-proof loop into helpers so the inheritance rules stay readable.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 3820 - 3907, Reduce the cyclomatic
complexity of inheritReviewArtifactEvidence by extracting the initial
commit/request validation into a helper and the diff path validation loop into
another helper. Keep all existing inheritance rules and return behavior
unchanged, and have inheritReviewArtifactEvidence delegate to these helpers
while retaining the Git verification and evidence construction flow.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:b5daafa6b44cc00f13b4c9fe -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Extracted request validation and review-root path validation into focused helpers, reducing `inheritReviewArtifactEvidence` complexity while preserving its fail-closed behavior. Focused artifact-inheritance tests passed.
