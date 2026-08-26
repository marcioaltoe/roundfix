---
status: completed
type: docs
---

# Task: Document the vocabulary this Spec emitted

QA finding F-2 (`Trust-Damage`): the PRD declares `surfaces: [backend, docs]`
and makes the identifier strategy applicable, assigning the closing node the
check of whether the work introduced or changed a term. The check ran and found
three tokens entering the QA Report artifact with no documentation anywhere:
`rows_blocked_precondition`, `precondition_check`, and `precondition_reason`.
The whole documentation surface received nothing.

## Work

- `CONTEXT.md`: the **QA Report** entry names `rows_blocked_precondition` beside
  the counts it already carries, and the two metadata keys. Add the precondition
  refusal itself as a term if the glossary has no home for it, following the
  shape of the entries already there.
- `docs/user-guide/context-driven-development.md`: the report contract states
  what a precondition refusal leaves behind, and that the mechanical stage reads
  the newest report only. Both are behavior a maintainer whose gate refused will
  look for, and neither is written anywhere today.
- `_techspec.md`: add the `## Vocabulary Contract` section whose absence made
  the checker skip its vocabulary detector rather than decide this. The skip is
  what let three undocumented tokens reach a passing static gate.
- Claims are read from the delivered code, not from the TechSpec's draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected: section
  5 of the TechSpec states the shape detector reads any markdown table, which
  task_07 measured false — it bound the scan to the `## Results` section, and
  the blockers came from subsections inside that window.

## References

- User Story 3: Maintainer reads what caused the refusal
- Core Feature 2: Mechanical stage reads the newest report
- Core Feature 4: Only the Results table is read as results

## Verification
- `grep -q "rows_blocked_precondition" CONTEXT.md && grep -q "precondition_check" CONTEXT.md && grep -qi "precondition" docs/user-guide/context-driven-development.md && grep -q "newest" docs/user-guide/context-driven-development.md && grep -q "## Vocabulary Contract" docs/specs/0113-a-refused-gate-writes-its-refusal-once/_techspec.md`

## Result
Every token this Spec emitted is named where the repository names its
vocabulary, and the TechSpec agrees with the code that was delivered.

### The glossary carries the tokens and the term

`CONTEXT.md`'s **QA Report** entry now names `rows_blocked_precondition` beside
the `rows_blocked_environment` and `rows_blocked_finding` counts it already
carried, plus `precondition_check` and `precondition_reason`, and records that a
gate which reached its matrix writes none of those three.

The glossary had no home for the refusal itself — no entry mentioned
"precondition" at all — so **Precondition Refusal** is a new term in the shape
of its neighbours (bold name, definition, `_Avoid_` line), placed beside
**QA Report**. It defines the stop as the gate's entire report rather than a
missing one: verdict `fail`, the count, the two naming keys, and the terminal
row `0 | blocked | precondition`, with the reason an empty Results table is not
the alternative.

### The report contract is written where a refused maintainer looks

`docs/user-guide/context-driven-development.md`, under "QA report reachability
and verdicts", gains two paragraphs after the report-naming rules: what a
precondition refusal leaves behind (the four frontmatter facts and the single
terminal row, with row `0` explained as a refusal rather than a measurement),
and that the mechanical stage validates the newest report in the directory and
nothing else, ranked by embedded date then rerun sequence — the same key the
naming paragraph above it defines. Both were unwritten anywhere before this
Task.

### The TechSpec declares its Vocabulary Contract

`_techspec.md` gains `## Vocabulary Contract` between Project Constraints and
the Coverage Map, declaring three entries — `internal/spec/qa.go` (writes the
report), `.agents/skills/qa-gate/SKILL.md` (instructs the refusing gate), and
`internal/speccheck/mechanical.go` (reconciles the count against the Results
table) — each with pattern
`rows_blocked_precondition|precondition_check|precondition_reason` documented in
`CONTEXT.md`. The pattern names the three tokens literally instead of matching
`rows_blocked_*`, because the environment, finding, and declared counts predate
this Spec; the section says so.

**The detector now runs instead of skipping.** A Spec with no contract still
reports the skip, which is the contrast that proves the change:

| Spec | `roundfix spec check <slug> --strict` |
| --- | --- |
| 0098 (no contract) | `vocabulary documentation detector: missing docs/specs/0098-.../_techspec.md Vocabulary Contract` |
| 0113 (after this Task) | no vocabulary skip, no vocabulary finding |

**The pass is not vacuous.** Negative control: `precondition_reason` was
temporarily broken to `preconditionZZreason` in `CONTEXT.md` (a suffix alone is
not a control — the detector uses substring containment, so
`precondition_reason_NEGATIVE_CONTROL` still passed). The broken form produced
exactly the expected refusals, naming both emitting files:

```text
SC-VOCABULARY-UNDOCUMENTED: internal/spec/qa.go emits undocumented token "precondition_reason" absent from CONTEXT.md
  at internal/spec/qa.go:89
SC-VOCABULARY-UNDOCUMENTED: .agents/skills/qa-gate/SKILL.md emits undocumented token "precondition_reason" absent from CONTEXT.md
  at .agents/skills/qa-gate/SKILL.md:59
```

`CONTEXT.md` was restored immediately; `grep -n "NEGATIVE_CONTROL\|preconditionZZ" CONTEXT.md` finds nothing and the clean check was re-run afterwards.

### The TechSpec is corrected against the delivered code

Section 5 claimed `parseMechanicalReport` "collects rows from every markdown
table in the file". task_07 measured that false and the code is the fact: the
scan was already bound to the `## Results` section and stopped at the next
`## `, so Spec 0098's blockers came from the four `### Row detail` subsections
*inside* that window. The section now carries a dated correction stating the
narrower defect — a scan bound to the section rather than to the table — and its
**Change** paragraph describes what shipped:
`parseMechanicalResultsRows` collects the first table under `## Results` and
stops at the next heading of any depth or at the first line after the table that
is not a table row. Read from `internal/speccheck/mechanical.go:979-1021`, not
from the draft.

### Focused checks

- `bin/roundfix spec check 0113-a-refused-gate-writes-its-refusal-once --strict` — `No findings.`, and the vocabulary skip is gone; only the pre-existing `SC-REF-UNRESOLVED` reference-index skip remains.
- `make spec-check` — exit `0` across all 18 active Specs, no findings, no vocabulary finding anywhere; `grep -c "vocabulary documentation detector: missing docs/specs/0113"` on the output is `0`.
- `make docs-test` — `ok roundfix/internal/docscontract`, run after the final guide wording edits.
- `make spec-budget` — `ok`, so the added TechSpec section keeps the corpus sweep inside its budget.
- `make repo-test` — `ok roundfix/internal/baseline`, `ok roundfix/skills`; no derived artifact drifted.
- Changed files: `CONTEXT.md`, `docs/user-guide/context-driven-development.md`, `docs/specs/0113-.../_techspec.md`, and this Task file. No source file was touched.

### Follow-ups discovered

- `rows_blocked_declared` is emitted by `internal/spec/qa.go` and
  `.agents/skills/qa-gate/SKILL.md` but appears in neither `CONTEXT.md` nor the
  user guide. It predates this Spec and belongs to whichever Spec introduced it,
  so this Task's pattern deliberately excludes it rather than silently absorbing
  another Spec's debt. Worth a backlog note.
- The Spec's own vocabulary findings surface as `Repair inputs`, not as
  `Findings`, so an undocumented token this Spec declares does not by itself
  fail `spec check`. That routing is `declaredVocabularySpec` behavior in
  `internal/speccheck`, not this Task's slice.
