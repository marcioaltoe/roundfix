---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/daemon/task_engine.go
line: 215
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeNp,comment:PRRC_kwDOS0qyts7c2V0N
review_hash: 1910672539f7ad0fb05ccc7393d1fe88f7dcb4f96b363959cbaf285cc60cad47
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Find every remaining producer and consumer of TaskPlan.QA.
set -euo pipefail

# Structural search for assignments and reads of the QA field.
ast-grep run --pattern 'TaskPlan{$$$}' --lang go internal | rg -n 'QA' -C 2

# Textual cross-check for field selection on a plan value.
rg -nP --type=go '\b(plan|taskPlan)\.QA\b'
```

Repository: marcioaltoe/roundfix

Length of output: 471

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- TaskPlan definition and cycle logic ---'
sed -n '190,360p' internal/daemon/task_engine.go

printf '%s\n' '--- All QA references in Go files ---'
rg -n --type=go '\bQA\b|\.QA\b' internal

printf '%s\n' '--- TaskPlan literals and assignments ---'
rg -n --type=go 'TaskPlan\s*\{|\.QA\s*=' internal
```

Repository: marcioaltoe/roundfix

Length of output: 43461

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' '--- TaskPlan constructors and test cases ---'
sed -n '180,235p' internal/daemon/task_engine_test.go
sed -n '4845,4910p' internal/daemon/task_engine_test.go
sed -n '645,685p' internal/cli/implement.go

printf '%s\n' '--- TaskPlan field consumers outside its declaration ---'
python3 - <<'PY'
from pathlib import Path
import re

for path in Path("internal").rglob("*.go"):
    text = path.read_text()
    for number, line in enumerate(text.splitlines(), 1):
        if re.search(r'\b(?:plan|taskPlan)\.QA\b', line):
            print(f"{path}:{number}:{line}")
PY
```

Repository: marcioaltoe/roundfix

Length of output: 6258

---



</details>

**Remove `TaskPlan.QA` and the obsolete legacy-state test.** `TaskCycle` reads QA state only from `plan.Tasks`. Remove the assignment at `internal/daemon/task_engine_test.go:4876`; no production caller sets this field.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine.go` at line 215, Remove the obsolete TaskPlan.QA
field and its legacy-state test coverage. Update TaskCycle and related code to
rely solely on plan.Tasks for QA state, and delete the assignment in the
relevant task engine test near the existing QA setup.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:54b84f555334d60c16d9b008 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Removed the unused `TaskPlan.QA` field and the test's obsolete assignment; graph Tasks remain the sole QA-state source. `go test ./internal/daemon -count=1` passed with 171 tests.
