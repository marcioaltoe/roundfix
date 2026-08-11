---
source: coderabbit
pr: "110"
round: 1
round_created_at: "2026-08-04T22:55:35Z"
status: invalid
terminal_reason: "The approved TechSpec assigns criterion matching to qa-gate and intentionally carries only typed counts into Archive."
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0070-implementation
head_sha: a588c6ca3ab9d977284ba1f9e80a89b0e6336786
file: internal/spec/archive.go
line: 118
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WeYX5,comment:PRRC_kwDOS0qyts7dggqd
review_hash: 85c26d847a4fa9678601aa778c00f1aae77cf7e329d24d04629c7044b7c01bf8
duplicate_of: ""
source_review_id: "4859094834"
source_review_submitted_at: "2026-08-04T21:23:48Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Map the report and declaration contracts before inspecting matching behavior.
ast-grep outline internal/spec/qa.go --items all
ast-grep outline internal/spec/spec.go --items all

# Expect stable criterion or declaration identity to flow into declared-row handling.
rg -n -C 6 --glob '*.go' \
  'RowsBlockedDeclared|Unreachable\(|Unreachable Acceptance|criterion|satisfied-by|readQAReport|runQAGate' \
  internal

# Inspect QA policy and fixtures for an exact mismatch case.
rg -n -C 6 \
  'rows_blocked_declared|Unreachable Acceptance|criterion|satisfied-by|unmatched|wrongly-declared' \
  .agents/skills/qa-gate/SKILL.md internal/spec internal/cli
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Inspect the archive call path, report parser, and QA report format to determine
# whether declaration identity is enforced or only documented.
cat -n internal/spec/archive.go | sed -n '60,130p'
cat -n internal/spec/qa.go | sed -n '20,45p;183,235p'
cat -n internal/spec/spec.go | sed -n '120,135p;196,330p'
cat -n .agents/skills/qa-gate/SKILL.md | sed -n '155,172p;230,285p'

printf '\nQA report fixtures and tests:\n'
rg -n -C 8 --glob '*.md' --glob '*_test.go' \
  'rows_blocked_declared|blocked \(declared:|Unreachable Acceptance|qa-report-' \
  internal docs .agents | head -n 300

printf '\nArchive tests around declared coverage:\n'
cat -n internal/spec/archive_test.go | sed -n '1,180p;230,360p'
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

# Model the exact archive predicate with the implemented data shape. This
# read-only probe demonstrates whether equal counts can distinguish identities.
python3 - <<'PY'
from dataclasses import dataclass

`@dataclass`(frozen=True)
class Declaration:
    criterion: str
    satisfied_by: str

`@dataclass`(frozen=True)
class QAReport:
    rows_blocked_declared: int

declarations = [Declaration("A", "human action")]
report = QAReport(rows_blocked_declared=1)

archive_accepts = (
    report.rows_blocked_declared != 0
    and report.rows_blocked_declared <= len(declarations)
)
print({
    "report_blocked_criterion": "C",
    "spec_declaration_criteria": [d.criterion for d in declarations],
    "report_has_blocked_criterion_identity": False,
    "archive_predicate_accepts": archive_accepts,
    "recorded_unproven": [d.satisfied_by for d in declarations],
})
PY
```

Repository: marcioaltoe/roundfix

Length of output: 350

---



</details>

**Bind blocked rows to their declared criteria before archiving.** `QAReport` stores only `RowsBlockedDeclared`, so one blocked row for criterion `C` passes against an unrelated declaration `A`; archive then records `A` as unproven. Carry criterion or declaration IDs through the report and reject missing, extra, or mismatched declarations. Add a regression test for this case.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/archive.go` around lines 96 - 118, Update the archive
validation around Unreachable and QAReport so blocked rows retain and validate
their criterion or declaration IDs, rejecting missing, extra, or mismatched
declarations instead of comparing only counts. Ensure archived actions are
derived only from declarations proven by the report, and add a regression test
covering a blocked row for criterion C incorrectly satisfying declaration A.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4360b703780b22d1b9853d99 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The approved TechSpec explicitly assigns blocked-row criterion matching to `qa-gate`, carries only `rows_blocked_declared` into `internal/spec`, and records the accepted trade-off that Archive trusts the gate's typed count instead of parsing the human-authored matrix or carrying a second identity model. Adding IDs would contradict that decision and expand the persisted report contract beyond this review Batch.
- Evidence: `docs/specs/0070-declared-unreachable-acceptance/_techspec.md` sections Executive Summary, System Architecture, Data Models, and Decisions state this boundary; the canonical `qa-gate` skill requires `blocked (declared: <criterion>)` only for a matching pre-run declaration.
