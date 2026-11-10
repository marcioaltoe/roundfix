---
task: task_03
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
status: completed
type: docs
complexity: low
---

# Task 03: The qa-gate skill, the guides, the references and the glossary say where Results rows go

## Overview

The qa-gate skill tells the gate to add every row directly below the last row
of the Results table, with no blank line and the same three columns, which
closes the cause of the defect reported by pantheon on 2026-11-09
([references](references/2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md)).
The qa-report, settle and implement guides, the context-driven development
guide, the Roundfix settle and implement references and `CONTEXT.md` state
which lines are Results rows, the refusal task_01 produces and the repair
turn task_02 gives.

## Requirements

1. MUST add to the canonical `.agents/skills/qa-gate/SKILL.md`, after the
   section `### Evidence belongs to this Run` and before `### Declared rows
   name their criterion`, this section verbatim (line breaks free), and change
   no `### QA settlement` text:

   > ### Results rows stay inside the table
   >
   > Add every row to the `## Results` table directly below its last row, with no blank line before it and the same three columns, `#`, `Status` and `Provenance`; in a Run that is the table the Daemon seeded. Put each row's detail in a block below the table. A row after a blank line is not part of the table: Markdown ends a table at its first blank line. `roundfix qa-report accept`, the Settle Command and `roundfix archive` refuse a report with such a row, naming its line and the blank line to remove, and in a Run the Daemon returns the QA turn once to repair it.

2. MUST add, in `docs/user-guide/commands/qa-report.md` after the Borrowed
   Evidence refusal block, and in the `## QA Report acceptance` section of the
   canonical `.agents/skills/roundfix/references/settle.md` after its Borrowed
   Evidence refusal block, this paragraph verbatim (line breaks free) and its
   `text` block:

   > Every Results row must belong to a table under `## Results`. A table starts at a header line followed by a separator line, and its rows run to the first line that is not a table row; the rows of every such table with a `Status` column are read, under a deeper heading too. A table line that no header leads, such as a row after a blank line, refuses the report before any other rule, naming the row's line and the repair: the blank line to remove, the line to move the row below, or the separator line a header lacks. The command exits `1` with no stdout and this stderr line:

   ```text
   roundfix: QA Report eligibility refused: QA Report <path>: Results row at line <n> is outside the Results table that starts at line <h>; remove the blank line at <path>:<b> so the row continues the table
   ```

3. MUST add to `docs/user-guide/commands/settle.md`, after the Borrowed
   Evidence refusal block, the sentence "For a `qa` Task, the Settle Command
   refuses a report with a Results row outside every table before any other
   rule, leaves the Task file unchanged, and names the report relative to the
   selected surface:" and the `text` block
   `roundfix: settle QA Report is ineligible: QA Report <path>: Results row at line <n> is outside the Results table that starts at line <h>; remove the blank line at <path>:<b> so the row continues the table (report <path>)`.
4. MUST add to `docs/user-guide/commands/implement.md` and the canonical
   `.agents/skills/roundfix/references/implement.md`, after the paragraph that
   ends with the `evidence_binding` event, this paragraph verbatim (line breaks
   free; in the reference, indented as its neighbours):

   > The QA prompt tells the gate to add each row directly below the last row of the seeded Results table, with no blank line before it. When a Results row of the newest QA Report lies outside every table, the Daemon returns the QA turn once, in the same turn as an unresolved link or Borrowed Evidence, naming the row's line and the repair, and publishes the `results_table` event. A `pass` or `partial` report with a row still outside settles the QA Task `failed` with `QA verdict <verdict> not accepted: QA Report <path>: Results row at line <n> is outside the Results table that starts at line <h>; <repair>`.

5. MUST replace, in `docs/user-guide/context-driven-development.md`, the
   clause "only the one table under the `## Results` heading is read as
   results, so every other table a report carries is evidence a reader weighs,
   never a row a later run can be blocked by" with "under the `## Results`
   heading only tables whose header has a `Status` column are read as
   results, so every other table a report carries is evidence a reader weighs,
   never a row a later run can be blocked by, and a table line there that no
   header leads, such as a row after a blank line, refuses the report, naming
   its line".
6. MUST append to the **QA Report** entry of `CONTEXT.md`, through
   `domain-modeling`, after its last sentence, the sentence "Its Results rows
   are the rows of every table under `## Results` whose header has a `Status`
   column; a table line there that no header leads, such as a row after a
   blank line, refuses the report, naming its line", followed by the citation
   of ADR-0381 in the entry's existing form: the label `ADR-0381` linked to
   `docs/adr/0381-a-results-row-outside-every-results-table-refuses-the-qa-report.md`
   in parentheses, then a period. Change no other glossary entry.
7. MUST copy the canonical skills with `make skills-sync`, then raise both
   version fields of `.agents/skills/qa-gate/SKILL.md`,
   `.agents/skills/roundfix/SKILL.md` and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`,
   recording `skills/testdata/owned-skill-versions.json`. Never hand-edit a
   mirror or a version field.
8. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of any guide, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and keep `docs/references/behavior-surfaces.json`.
9. MUST NOT touch any Go file, any other skill or guide, any Baseline module,
   `docs/agents/`, any ADR, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Add the qa-gate section.
- [ ] Add the acceptance paragraph to the qa-report guide and the settle reference, and the settle refusal to the settle guide.
- [ ] Add the repair-turn paragraph to the implement guide and reference, and correct the context-driven development guide.
- [ ] Append the **QA Report** sentence to `CONTEXT.md`.
- [ ] Sync the mirrors, record the skill versions, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both copies of the qa-gate skill carry `### Results rows stay inside the table`, and its `### QA settlement` section is unchanged.
- [ ] The guides and both reference copies carry the rule, the refusal lines and the repair turn, and the context-driven development guide no longer says only one table is read.
- [ ] `CONTEXT.md` carries the new **QA Report** sentence; the mirrors equal their canonical copies; the owned-skill version, glossary and Behavior Surface contracts pass.

## Context

- instruction: `docs/adr/0381-a-results-row-outside-every-results-table-refuses-the-qa-report.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/settle.md`
- interface: `docs/user-guide/commands/implement.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- The owned-skill version record, the glossary and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do grep -q '^### Results rows stay inside the table$' "$f" || { printf 'missing qa-gate section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'directly below its last row, with no blank line before it and the same three columns' || { printf 'missing qa-gate rule in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/qa-report.md .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Every Results row must belong to a table under' || { printf 'missing acceptance rule in %s\n' "$f" >&2; exit 1; }; grep -qF -- 'is outside the Results table that starts at line <h>; remove the blank line at <path>:<b> so the row continues the table' "$f" || { printf 'missing refusal line in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/settle.md | grep -qF -- 'roundfix: settle QA Report is ineligible: QA Report <path>: Results row at line <n>' || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/settle.md 'roundfix: settle QA Report is ineligible: QA Report <path>: Results row at line <n>' >&2; exit 1; }; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'When a Results row of the newest QA Report lies outside every table, the Daemon returns the QA turn once' || { printf 'missing repair turn in %s\n' "$f" >&2; exit 1; }; done; ! tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'only the one table under the' || { echo 'the guide still reads one table' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'only tables whose header has a' || { echo 'missing corrected guide clause' >&2; exit 1; }; grep -F -- '**QA Report**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*QA Report\*\*(:| —( |$))' || { echo 'missing QA Report term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'Its Results rows are the rows of every table under' || { echo 'missing QA Report sentence' >&2; exit 1; }; cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md && cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither qa-gate copy has the section, so the first check fails; afterwards both skills, the guides, both reference copies and the glossary carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint and the glossary phrase test among them, pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 5; Glossary
- `_techspec.md` → Current behavior 8; API Contract 11; Build Order 3
- ADR-0381

## Result

Skills loaded: implement-task, tech-writer, knowledge-workspace, domain-modeling, evidence-gate

- Added the Results-table placement rule to the canonical qa-gate skill; added the QA Report acceptance refusal to the qa-report guide and settle reference; added the Settle refusal to its guide; added the repair-turn rule to both implement surfaces; corrected the context-driven development guide; and appended the QA Report glossary sentence with ADR-0381.
- `make skills-sync` passed. The owned-skill version recorder passed with `GOCACHE=/tmp/roundfix-skill-cache-task03 go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`; the focused non-recording version test also passed.
- `cmp` passed for both skill files and both Roundfix reference mirrors. A byte comparison of the `### QA settlement` section against `HEAD` passed. `git diff --check` passed.
- The final Behavior Surface regeneration passed with `GOCACHE=/tmp/roundfix-skill-cache-task03 go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
- The Daemon-owned `status: in_progress` field and the declared Verification section remain unchanged; no Go file, other task, Task Graph manifest, ADR, Baseline module, Makefile, CI workflow, `go.mod` or `go.sum` was changed.

## Carry-forward provenance

- Source Run: `run_20261110T215130Z_fa9f763e0a43c500`
- Source commit: `4d1ed9ebe9781958abafcee3c475cb6ef2b3fac9`
