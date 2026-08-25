---
spec: 0113-a-refused-gate-writes-its-refusal-once
status: active
created: 2026-08-25
---

# TechSpec: Gate Refusal Report Shape

## Context

When the QA gate refuses at a precondition check (e.g., a strict Spec consistency
check), it writes a QA Report. Currently, the Results table in that report is
empty — correctly, because the gate stopped before building the matrix.

However, the mechanical stage of every later run reads that report and refuses
with `QA-REPORT-SHAPE`:

```text
Results table has no report rows
fix: Materialize every planned QA row with one terminal status.
```

This creates an unsolvable loop:
1. Gate refuses at precondition → writes empty report
2. Next run reads empty report → refuses on shape error
3. The only exit is manual deletion of the report (evidence removal, forbidden)

Measured on Spec 0103 (2026-08-14) and Spec 0094 (2026-08-12/13). The gate is
doing the right thing by checking facts before spending an Agent turn (ADR-0096).
The problem is the report it leaves behind.

## Root Cause Analysis

1. **Precondition refusal produces invalid report shape**: The gate builds a
   matrix only after passing preconditions. A precondition refusal stops before
   that, leaving an empty Results table.

2. **Mechanical stage rejects empty reports**: The report validation requires
   every planned row to be materialized with a terminal status. An empty table
   violates this contract.

3. **Report reading is not filtered**: The mechanical stage reads every report
   in the QA directory, not just the current one. A superseded refusal blocks
   the run that supersedes it.

4. **Report is written despite stopping**: When the gate stops at a precondition,
   it could write nothing, or it could write a valid report shape. Currently it
   writes an invalid one.

## Solution Design

### 1. Gate Refusal Produces a Terminal Row

When the QA gate refuses at a precondition:

**Before**:
```markdown
## Results

| # | Status | Provenance |
| - | --- | --- |
```

**After**:
```markdown
## Results

| # | Status | Provenance |
| - | --- | --- |
| 0 | blocked | precondition |
```

The terminal row:
- **Status**: `blocked` (gate did not run, precondition prevented it)
- **Provenance**: `precondition` (not a matrix result, a refusal record)
- **Row 0**: Special terminal row that records the refusal itself

### 2. Precondition Refusal Information

The terminal row includes metadata about what precondition was checked:

**Implementation**:
- Store in QA Report frontmatter: `rows_blocked_precondition: <check_name>`
- Check name: the name of the Spec check that refused (e.g., `strict`, `vocabulary`)
- Reason: the reason it failed (e.g., "undocumented term", "contradictory requirement")

**Report shape**:
```yaml
---
verdict: fail
rows_blocked_precondition: 1
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_unproven: 0
precondition_check: "strict"
precondition_reason: "SC-VOCABULARY-UNDOCUMENTED"
---
```

### 3. Mechanical Stage Reads Newest Report Only

Update the stage that reads QA Reports before running the gate:

**Before**: Read every report in the QA directory, fail if any is invalid

**After**: Read only the newest report by timestamp/name. If that report is
newer than the gate's last run, accept it. Superseded reports are ignored.

**Logic**:
```
reports = ls docs/specs/SLUG/qa/qa-report-*.md | sort -r
if len(reports) == 0:
  proceed (no prior report)
elif len(reports) >= 1:
  latest = reports[0]
  if validate(latest) == ok:
    proceed
  else:
    refuse SC-REPORT-SHAPE on latest only
```

This way:
- A refusal from a prior run does not block a new run
- The current run's own newest report is the authority
- Superseded refusals are naturally ignored

### 4. Terminal Row Semantics

The terminal row that records a precondition refusal:

| # | Status | Provenance | Requirement |
| - | --- | --- | --- |
| 0 | blocked | precondition | (none — this row records the refusal, not a requirement) |

**Interpretation**:
- The gate tried to run and was prevented by a precondition check
- No requirements were measured
- The refusal is recorded as metadata (precondition_check, precondition_reason)
- The report is valid per the QA Report contract

## Specification

### QA Gate Refusal Path

1. Gate runs precondition checks (e.g., `spec check --strict`)
2. If precondition check fails:
   - Do not build matrix (matrix requires passed preconditions)
   - Create QA Report with verdict `fail`
   - Frontmatter: `rows_blocked_precondition: 1`
   - Store check name and reason
   - Create Results table with one terminal row (status: blocked, provenance: precondition)
   - Write report and stop

3. If precondition check passes:
   - Build matrix (existing behavior)
   - Run verification and QA tasks
   - Materialize every planned row
   - Verdict: pass/partial/fail based on results

### Mechanical Stage Report Reading

1. Load all QA reports from `docs/specs/SLUG/qa/qa-report-*.md`
2. Sort by filename (timestamp-based)
3. Select the newest report only
4. Validate report shape:
   - If valid: proceed with gate
   - If invalid (empty Results, missing frontmatter): refuse SC-REPORT-SHAPE
5. Ignore reports that are not the newest (superseded reports)

### Acceptance Rows

Every acceptance row must list what was blocked and why:

| # | Requirement | Blocked By | Evidence |
| - | --- | --- | --- |
| 1 | Term must be in glossary | precondition (strict) | SC-VOCABULARY-UNDOCUMENTED |
| 2 | Contradictions must be resolved | precondition (strict) | SC-REQUIREMENT-CONTRADICTORY |
| etc | - | - | - |

When a requirement is blocked by precondition, it is listed in `rows_blocked_precondition`
in the frontmatter.

## Acceptance Criteria

### Functional

1. **Gate writes valid report on precondition refusal**
   - Results table has one terminal row (not empty)
   - Row status is `blocked`, provenance is `precondition`
   - Report frontmatter has `rows_blocked_precondition: 1`
   - Precondition check name and reason recorded

2. **Mechanical stage reads newest report only**
   - Multiple reports in QA directory → newest is read
   - Superseded report does not block next run
   - Report ordering is deterministic (by timestamp/name)

3. **SC-REPORT-SHAPE never blocks a run that supersedes its cause**
   - Spec 0103: precondition refusal on 2026-08-14 → can retry on 2026-08-15
   - No manual intervention needed
   - Report is valid per its own contract

4. **Precondition check information is recoverable**
   - Report includes what check refused
   - Report includes why it refused
   - Supervisor can read the report and understand the block

### Evidence

5. **Three consecutive Specs confirm the pattern is stable**
   - Spec 0078 (historical)
   - Spec 0094 (2026-08-12/13)
   - Spec 0103 (2026-08-14)
   - Each shows the same shape, fixed by same solution

6. **No external acceptance row required**
   - Gate refusal is internal to QA contract
   - Report shape is Roundfix-defined
   - Evidence is the measured pattern

## Related Specs

**Spec 0098** (Hook Strictness) depends on this: After 0113 gates correctly on
precondition refusal, 0098 can proceed without deadlock. The two Specs are
independent implementations but sequenced for dependency reasons.

**Spec 0103** (Vocabulary Contract) is unblocked by this: Once 0113 allows
precondition refusal reports to not block the next run, Spec 0103 can retry
and pass.
