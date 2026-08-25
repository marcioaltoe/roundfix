---
spec: 0113-a-refused-gate-writes-its-refusal-once
qa: task_06
---

# Task Graph: Gate Refusal Report Shape

## Dependency Graph

```
task_01 (precondition terminal row) →
  task_02 (gate refusal detection) →
    task_03 (report metadata) →
      task_04 (mechanical stage update) →
        task_05 (newest-only reading) →
  task_06 (qa gate) [terminal, depends on all]
```

## Tasks

### task_01 — Write Terminal Row on Precondition Refusal
**Type**: backend  
**Status**: pending

Implement the gate refusal path to produce a valid QA Report with a terminal row
instead of an empty Results table.

**Work**:
- Modify QA gate precondition handling (internal/gate/gate.go or similar)
- When precondition check fails:
  - Do not attempt matrix building
  - Create Results table with exactly one row
  - Row: `| 0 | blocked | precondition |`
  - Set frontmatter: `rows_blocked_precondition: 1`
  - Store check name: `precondition_check: <name>`
  - Store reason: `precondition_reason: <reason>`
  - Set verdict: `fail`
  - Write report

**Verification**:
```bash
# Run spec check with strict precondition failure
# Verify report has terminal row and valid shape
grep -q "| 0 | blocked | precondition |" qa-report*.md && \
grep -q "rows_blocked_precondition: 1" qa-report*.md
```

**Context**:
- Gate entry point: internal/gate/gate.go
- Precondition checks: spec check, vocabulary, requirements, etc.
- Report writing: use existing report serialization (don't create new paths)

### task_02 — Detect Precondition Failure
**Type**: backend  
**Status**: pending  
**Depends**: task_01

Implement detection of which precondition failed and why.

**Work**:
- Call `spec check --strict` and capture output
- Parse output for Spec check error codes (SC-VOCABULARY-UNDOCUMENTED, etc.)
- Extract check name and failure reason
- Store in gate context for report writing
- Each precondition check carries a stable name (vocabulary, requirements, etc.)

**Verification**:
```bash
# Verify precondition check names are captured correctly
# Example: SC-VOCABULARY-UNDOCUMENTED → check_name="vocabulary"
make test -k "TestPreconditionDetection" | grep -q "ok"
```

**Context**:
- Spec check output format: documented in spec check README
- Error codes: SC-VOCABULARY-*, SC-REQUIREMENT-*, etc. (stable tokens)
- Parse first error found (precondition stops at first failure)

### task_03 — Store Precondition Metadata
**Type**: backend  
**Status**: pending  
**Depends**: task_02

Update QA Report structure to include precondition refusal metadata.

**Work**:
- Add frontmatter fields to QA Report:
  - `precondition_check: <check_name>`
  - `precondition_reason: <error_code_or_message>`
- Ensure these fields are optional (not required for passing gates)
- Preserve when report is read and re-written
- Include in report validation (valid reports may have these fields)

**Verification**:
```bash
# Verify frontmatter includes precondition fields
# Read report and confirm unmarshal succeeds
make test -k "TestQAReportMetadata" | grep -q "ok"
```

**Context**:
- QA Report: docs/specs/SLUG/qa/qa-report-YYYY-MM-DD.md
- Frontmatter: YAML at top of file
- Do not include in Results table (metadata only)

### task_04 — Update Mechanical Stage Report Validation
**Type**: backend  
**Status**: pending  
**Depends**: task_03

Modify the mechanical stage that validates reports before gate runs.

**Work**:
- Load QA Report from disk
- Validate report shape:
  - **If Results table is empty** (zero rows): refuse SC-REPORT-SHAPE
    (This should not happen after task_01, but keep guard)
  - **If Results table has rows**: parse each row
    - Check status is valid (completed, failed, blocked, unproven)
    - Check provenance is valid (task, gate, precondition, etc.)
  - **Accept**: terminal row with status=blocked, provenance=precondition
  - **Refuse**: empty table (no rows at all)

**Verification**:
```bash
# Valid report with precondition row should pass
# Empty report should fail with SC-REPORT-SHAPE
make test -k "TestReportValidation" | grep -q "ok.*2 subtests"
```

**Context**:
- File: internal/gate/mechanical.go (or report validation module)
- Current check: "Results table must have at least one row"
- New check: "Accept terminal blocked row from precondition"

### task_05 — Read Newest Report Only
**Type**: backend  
**Status**: pending  
**Depends**: task_04

Update report discovery to read only the newest report, ignoring superseded ones.

**Work**:
- When loading QA report for validation:
  - List all `qa-report-*.md` files in QA directory
  - Sort by filename (YYYY-MM-DD based): newest first
  - Select only the first (newest) report
  - Validate only that report
  - Ignore all older reports (do not loop through them)

- When gate runs and writes new report:
  - Use timestamp-based filename (existing practice)
  - New report automatically becomes the newest
  - Older reports remain but are ignored by next run

**Verification**:
```bash
# Create multiple reports
# Verify only newest is read
# Verify older report does not block new run
make test -k "TestNewestReportOnly" | grep -q "ok"
```

**Context**:
- Filename format: `qa-report-2026-08-14.md` (timestamp-based)
- Sort order: lexicographic sorts correctly by date
- Fallback: if no reports exist, proceed (allow first-run gate)

### task_06 — QA Gate
**Type**: qa  
**Status**: pending  
**Depends**: task_01, task_02, task_03, task_04, task_05

Verify all deliverables and document acceptance evidence.

**Acceptance Rows**:

| # | Requirement | Evidence | Status |
| - | --- | --- | --- |
| 1 | Terminal row on precondition refusal | Test: gate writes blocked row | TBD |
| 2 | Precondition metadata captured | Test: check_name and reason stored | TBD |
| 3 | Mechanical stage validates new shape | Test: blocked row passes validation | TBD |
| 4 | Newest report only is read | Test: older reports ignored | TBD |
| 5 | Spec 0103 can retry without deadlock | Integration: 0103 runs after 0113 merge | TBD |
| 6 | Three specs confirm pattern | Historical: 0078, 0094, 0103 evidence | External |

**External Acceptance Evidence**:
- **Three consecutive Specs** in the codebase have measured the same defect
  - Spec 0078: documented precondition refusal loop
  - Spec 0094 (2026-08-12/13): reproduced the empty-table path
  - Spec 0103 (2026-08-14): deadlock from precondition refusal
  - All three show: refusal → empty report → next run refuses on shape
  - All three would be solved by: terminal row on refusal + newest-only read

**Verify**:
```bash
roundfix spec check 0113 --strict && \
make test -k "TestGateRefusal" && \
# Confirm 0103 can execute after 0113 is merged
git log --oneline | grep -q "0113"
```

## Implementation Notes

### Report Filename Ordering

The report filenames use dates: `qa-report-2026-08-14.md`, `qa-report-2026-08-15.md`, etc.

Lexicographic sort on these filenames naturally produces chronological order:
```
qa-report-2026-08-13.md  <- oldest
qa-report-2026-08-14.md
qa-report-2026-08-15.md  <- newest (sorted first in reverse)
```

This is reliable and requires no timestamp parsing.

### Interaction with Spec 0103

Spec 0103 failed with:
1. Gate ran successfully
2. Gate refused at precondition (vocabulary term missing)
3. Gate wrote empty report
4. Next run read empty report
5. Mechanical stage refused on SC-REPORT-SHAPE
6. Run deadlocked

After this Spec 0113 is merged:
1. Gate runs successfully
2. Gate encounters same precondition (vocabulary term missing)
3. Gate writes report with terminal blocked row
4. Next run reads newest report (now valid)
5. Mechanical stage accepts report
6. Gate can run again with the precondition fixed
7. No deadlock, no manual intervention

### Irreversible Design Choice

Once this Spec is merged, reports with precondition refusals will always have
a terminal row. This is backwards-compatible: older reports (if any) with empty
Results will still be refused on shape, but new refusals will be valid.

Never go back to empty-table reports after this change.
