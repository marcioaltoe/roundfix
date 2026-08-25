---
schema: spec-tasks/v1
spec: 0113-a-refused-gate-writes-its-refusal-once
qa: task_06
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Gate Refusal Report Shape

| id      | title                                   | type    | complexity | needs    |
| ------- | --------------------------------------- | ------- | ---------- | -------- |
| task_01 | Write terminal row on precondition     | backend | high       | —        |
| task_02 | Detect precondition failure            | backend | high       | task_01  |
| task_03 | Store precondition metadata            | backend | medium     | task_02  |
| task_04 | Update mechanical stage validation     | backend | high       | task_03  |
| task_05 | Read newest report only                | backend | high       | task_04  |
| task_06 | QA gate                                | qa      | medium     | task_05  |

Waves: 1 → task_01 · 2 → task_02 · task_03 · task_04 · task_05 · 3 → task_06

## Task: task_01 — Write Terminal Row on Precondition Refusal

**Type**: backend **Status**: pending

Implement gate refusal path to write valid QA Report with terminal row instead of empty Results table.

**Acceptance**:
- Precondition check failure → single terminal row
- Row: `| 0 | blocked | precondition |`
- Frontmatter: `rows_blocked_precondition: 1`
- Store check name and reason
- Set verdict: `fail`

**Verification**:
```bash
grep -q "| 0 | blocked | precondition |" /tmp/qa-report*.md && \
grep -q "rows_blocked_precondition: 1" /tmp/qa-report*.md
```

---

## Task: task_02 — Detect Precondition Failure

**Type**: backend **Status**: pending **Needs**: task_01

Implement detection of which precondition failed and why.

**Acceptance**:
- Call `spec check --strict` and capture output
- Parse output for error codes (SC-VOCABULARY-UNDOCUMENTED, etc.)
- Extract check name and failure reason
- Store for report writing

**Verification**:
```bash
make test -k TestPreconditionDetection | grep -q "ok"
```

---

## Task: task_03 — Store Precondition Metadata

**Type**: backend **Status**: pending **Needs**: task_02

Update QA Report structure to include precondition refusal metadata.

**Acceptance**:
- Frontmatter fields: `precondition_check`, `precondition_reason`
- Fields are optional (not required for passing gates)
- Preserved when report is read/written

**Verification**:
```bash
make test -k TestQAReportMetadata | grep -q "ok"
```

---

## Task: task_04 — Update Mechanical Stage Validation

**Type**: backend **Status**: pending **Needs**: task_03

Modify mechanical stage to accept terminal blocked row from precondition.

**Acceptance**:
- Empty Results table → refuse SC-REPORT-SHAPE
- Terminal blocked row → accept
- Status and provenance validated

**Verification**:
```bash
make test -k TestReportValidation | grep -q "ok.*2 subtests"
```

---

## Task: task_05 — Read Newest Report Only

**Type**: backend **Status**: pending **Needs**: task_04

Update report discovery to read only the newest report, ignoring superseded ones.

**Acceptance**:
- List all `qa-report-*.md` files
- Sort by filename (newest first)
- Select only first (newest) report
- Ignore all older reports

**Verification**:
```bash
make test -k TestNewestReportOnly | grep -q "ok"
```

---

## Task: task_06 — QA Gate

**Type**: qa **Status**: pending **Needs**: task_05

Verify all deliverables and document acceptance evidence.

**Acceptance Rows**:

| # | Requirement | Evidence |
| - | --- | --- |
| 1 | Terminal row on precondition refusal | Test: blocked row written |
| 2 | Precondition metadata captured | Test: check_name and reason stored |
| 3 | Mechanical stage validates new shape | Test: blocked row passes |
| 4 | Newest report only is read | Test: older reports ignored |
| 5 | Three specs confirm pattern | Historical: 0078, 0094, 0103 |

**Verification**:
```bash
roundfix spec check 0113-a-refused-gate-writes-its-refusal-once --strict && \
make test -k TestGateRefusal | grep -q "ok"
```
