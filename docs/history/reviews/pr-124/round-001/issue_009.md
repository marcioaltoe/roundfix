---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: invalid
terminal_reason: "Spec 0077 Task 01 requires the previously verified empty-title CodeRabbit success shape to remain verified; non-empty unknown titles still resolve pending."
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/reviewsource/coderabbit/coderabbit.go
line: 849
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wua0h,comment:PRRC_kwDOS0qyts7d35uB
review_hash: 046dc43236fc10c7706a772f8662ca32868d8f74aec8870f8f691d89a3744636
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 009: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Locate recorded CodeRabbit check-run payloads and spec statements about review-completed evidence.
set -euo pipefail

rg -n -C 5 'review completed' --glob '!**/vendor/**'
rg -n -C 5 'reviewCompleted|OutputTitle' --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- target implementation ---'
sed -n '820,950p' internal/reviewsource/coderabbit/coderabbit.go

printf '%s\n' '--- related tests and call sites ---'
rg -n -C 8 'reviewCompleted|successful check verifies|OutputTitle|reviewStatusCompleted' internal/reviewsource/coderabbit --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 32756

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- check-run parsing ---'
sed -n '440,515p' internal/reviewsource/coderabbit/coderabbit.go

printf '%s\n' '--- evidence classification ---'
sed -n '630,725p' internal/reviewsource/coderabbit/coderabbit.go

printf '%s\n' '--- contract references ---'
rg -n -i -C 4 'recognised review|review-completed|review completed|positive evidence|green check|without.*title|output title' . \
  --glob '!**/.git/**' --glob '!**/vendor/**' --glob '*.go' --glob '*.md' --glob '*.yaml' --glob '*.yml' || true

printf '%s\n' '--- signal helper ---'
rg -n -C 6 'func isCodeRabbitSignal|func normalized' internal/reviewsource/coderabbit/coderabbit.go
```

Repository: marcioaltoe/roundfix

Length of output: 50379

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

python3 - <<'PY'
def normalized(value):
    return " ".join(value.strip().lower().split())

def review_completed(run):
    if normalized(run["status"]) != "completed" or normalized(run["conclusion"]) != "success":
        return False
    title = normalized(run["output_title"])
    if title:
        return title == "review completed"
    name = normalized(run["name"])
    return name in {"coderabbit", "review completed"}

cases = [
    {"name": "CodeRabbit", "status": "completed", "conclusion": "success", "output_title": ""},
    {"name": "CodeRabbit", "status": "completed", "conclusion": "success", "output_title": "Unknown outcome"},
    {"name": "CodeRabbit", "status": "completed", "conclusion": "success", "output_title": "Review completed"},
]
for case in cases:
    print(case, "=>", review_completed(case))
PY
```

Repository: marcioaltoe/roundfix

Length of output: 487

---



</details>

**Do not treat an empty output title as completed-review evidence.**

When `OutputTitle` is empty, a completed successful check named `CodeRabbit` returns `true` and reaches `verified` at `classifyEvidence`. This violates the documented rule that an unrecognised green check remains `pending` because it does not prove that a review ran. Remove the `name == "coderabbit"` fallback, or document a separate positive signal that justifies it, and update the test at `coderabbit_test.go:680`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/reviewsource/coderabbit/coderabbit.go` around lines 840 - 849,
Update reviewCompleted to return false when OutputTitle is empty; remove the
name == "coderabbit" fallback so an unrecognised successful check remains
pending. Preserve the existing “review completed” title and name handling, and
update the corresponding coderabbit test to assert the empty-title CodeRabbit
case is not completed.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:622e3dea1cfcc40897776582 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The accepted Spec deliberately preserves this compatibility shape.
  Task 01 requirement 4 requires every payload that previously resolved
  `verified` to keep doing so, and its Result calls the unchanged empty-title
  “successful check” cases the canaries for `verified` and `reviewed`. The
  fallback applies only to a current-head, completed, successful check whose
  recognised name is `CodeRabbit` or `Review completed`; any non-empty unknown
  output title already returns `pending`. Removing it would contradict the
  accepted non-regression contract.
- Focused evidence: inspected archived Spec 0077 `_prd.md`, `_techspec.md`, and
  `task_01.md`; `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/reviewsource/coderabbit -count=1` passed unchanged.
- Daemon Verification: `make verify` not run; Daemon-owned.
