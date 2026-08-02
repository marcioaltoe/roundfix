---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: invalid
terminal_reason: "No active repository gocyclo, funlen, or nestif threshold requires extracting this already bounded retention branch."
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/plan.go
line: 1145
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymh,comment:PRRC_kwDOS0qyts7cjgE9
review_hash: 8d6145f6c406b8982bd37af83cc104f641c9e6388d9d6ede64bd0a94c481daa5
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 011: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Report configured complexity limits and function sizes.
set -euo pipefail

fd -H -t f '^\.golangci\.(yml|yaml|toml|json)$' --exec rg -n -C 3 'cyclop|gocyclo|gocognit|funlen|nestif' {}
ast-grep outline internal/baseline/plan.go --items all --type function --match 'resolvePlanRetention'
```

Repository: marcioaltoe/roundfix

Length of output: 2034

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

sed -n '1040,1175p' internal/baseline/plan.go
printf '\n-- configured linter availability --\n'
command -v gocyclo || true
command -v golangci-lint || true
printf '\n-- relevant function metrics --\n'
python3 - <<'PY'
from pathlib import Path
p = Path("internal/baseline/plan.go")
lines = p.read_text().splitlines()
start = next(i for i, line in enumerate(lines) if line.startswith("func resolvePlanRetention("))
# Find the next top-level function declaration.
end = next((i for i in range(start + 1, len(lines))
            if lines[i].startswith("func ")), len(lines))
body = lines[start:end]
print(f"start_line={start+1}")
print(f"end_line={end}")
print(f"line_count={len(body)}")
print(f"nonblank_statement_like_lines={sum(bool(x.strip()) for x in body)}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 5013

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

sed -n '1040,1175p' internal/baseline/plan.go
printf '\n-- configured linter availability --\n'
command -v gocyclo || true
command -v golangci-lint || true
printf '\n-- relevant function metrics --\n'
python3 - <<'PY'
from pathlib import Path
p = Path("internal/baseline/plan.go")
lines = p.read_text().splitlines()
start = next(i for i, line in enumerate(lines) if line.startswith("func resolvePlanRetention("))
end = next((i for i in range(start + 1, len(lines))
            if lines[i].startswith("func ")), len(lines))
body = lines[start:end]
print(f"start_line={start+1}")
print(f"end_line={end}")
print(f"line_count={len(body)}")
print(f"nonblank_statement_like_lines={sum(bool(x.strip()) for x in body)}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 5013

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

golangci-lint run ./internal/baseline \
  --disable-all \
  --enable=gocyclo \
  --enable=funlen \
  --enable=nestif \
  --out-format=tab 2>&1
status=$?
printf '\nexit_status=%s\n' "$status"
exit 0
```

Repository: marcioaltoe/roundfix

Length of output: 279

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

golangci-lint run --help 2>&1 | sed -n '1,180p'
printf '\n-- version --\n'
golangci-lint version 2>&1
```

Repository: marcioaltoe/roundfix

Length of output: 5766

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

golangci-lint run ./internal/baseline \
  --no-config \
  --enable-only=gocyclo,funlen,nestif \
  --output.tab.colors=false \
  --output.tab.print-linter-name=true 2>&1
status=$?
printf '\nexit_status=%s\n' "$status"
exit 0
```

Repository: marcioaltoe/roundfix

Length of output: 14163

---

🏁 Script executed:

```shell
#!/bin/bash
set -u

golangci-lint run ./internal/baseline \
  --enable-only=gocyclo,funlen,nestif \
  --output.tab.colors=false \
  --output.tab.print-linter-name=true 2>&1
status=$?
printf '\nexit_status=%s\n' "$status"
exit 0
```

Repository: marcioaltoe/roundfix

Length of output: 14163

---



</details>

**Extract the same-baseline retention branch from `resolvePlanRetention`.** The manifest branch at line 1089 has nested-block complexity 12. Extracting it will keep the retention stages readable and independently testable.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan.go` around lines 1089 - 1145, Extract the
same-baseline retention logic from resolvePlanRetention into a focused helper,
including strict manifest parsing, equality checks, retention-source lookup,
classifySourceClauseTransition, and delta handling. Have resolvePlanRetention
delegate to the helper while preserving all existing return values and fallback
behavior, so the main retention stages remain readable and the extracted logic
is independently testable.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:6d2a4949a90011bb758af0a5 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The reviewer ran opt-in linters outside `make verify`; the repository has no matching complexity contract. The same-baseline branch is one cohesive retention decision flow, so extraction would add indirection without a behavioral or gate requirement.
