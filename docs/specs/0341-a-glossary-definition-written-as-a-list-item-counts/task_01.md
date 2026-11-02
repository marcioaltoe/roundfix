---
task: task_01
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
status: pending
type: backend
complexity: medium
---

# Task 01: The glossary reader accepts list-item and em dash definitions and the findings name the forms

## Overview

The Spec Consistency Check reads a Glossary Definition only from a line that
starts with `**Term**:`, so a glossary written as `- **Term** — definition`
never defines its terms and `SC-GLOSSARY-MISSING` fails the QA gate without
saying which form it wanted. This Task widens the one reader to the accepted
forms, makes `SC-GLOSSARY-MISSING` name them and `SC-GLOSSARY-UNPLANNED`
print the line-anchored Verification command, proves the command and the
reader agree, and writes the rule into the spec command guide and
`CONTEXT.md`. It answers the Pantheon report of 2026-11-02 (ai-memory
message `01a19a04`, adopted as
[references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md](references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md)).

## Requirements

1. MUST characterize first: add the package-level table
   `glossaryDefinitionFormCases` and `TestGlossaryDefinitionFormsAreRead`
   (`_techspec.md` → Testing Approach) with today's expectations (only the
   `colon` row defined), run it green on the unchanged reader, and record that
   run in the Task Result; then set the expectations of API Contract 1 and
   change the reader.
2. MUST replace `glossaryDefinitionPattern` in
   `internal/speccheck/glossary.go` with
   `^(?:[-*+] )?\*\*([^*]+)\*\*(?::| —(?: |$))` and keep `glossaryFiles`'
   loop, normalization and file set (`_techspec.md` → API Contract 1).
3. MUST add the exported `GlossaryDefinitionForms` constant with the exact
   text of `_techspec.md` → Interfaces, and `GlossaryVerificationCommand(term,
   glossary string) string` rendering API Contract 2 (literal `grep -F` of the
   bolded term reading the glossary on standard input (`< <glossary>`),
   piped to a line-anchored `grep -qE` whose term is escaped with
   `regexp.QuoteMeta`, each argument single-quoted with an apostrophe written
   as `'\''`).
4. MUST change the `SC-GLOSSARY-MISSING` fix to
   `Define the declared term in the glossary as ` + `GlossaryDefinitionForms`
   + `.`, and append `, such as: ` + `GlossaryVerificationCommand(term, path)`
   to the `SC-GLOSSARY-UNPLANNED` fix for a term without a binding Task, with
   `path` chosen as API Contract 3 states; every other summary, code,
   location and fix keeps its bytes.
5. MUST add `TestGlossaryVerificationCommandAgreesWithTheDefinitionReader`
   and `TestGlossaryFindingsNameTheDefinitionForms` as `_techspec.md` →
   Testing Approach describes, including the assertion that
   `WrapFragileVerification` reports nothing for the rendered command. The
   agreement test resolves `grep` with
   `testfixture.RequireTool`, runs only `sh` in a `t.TempDir()`, and starts no
   process that outlives the test. Every new top-level test calls
   `t.Parallel()` first.
6. MUST state the forms in `docs/user-guide/commands/spec.md`: replace
   "Definitions begin a line with `**Term**:`." with a sentence that starts
   "A Glossary Definition is a line that starts with" and names the four
   forms and the `*` and `+` markers, and add to the `SC-GLOSSARY-MISSING` row
   that the glossary must hold a Glossary Definition of the term in one of
   those forms.
7. MUST add **Glossary Definition** to `CONTEXT.md` through domain-modeling,
   next to **Glossary Declaration**, written as a `**Glossary Definition**:`
   line: a line of the glossary that defines a term, the bolded term at the
   start of the line followed by a colon or by a spaced em dash, bare or as a
   top-level list item marked `-`, `*` or `+`; a heading, an indented or
   ordered list item and bold text inside a sentence define nothing. It
   keeps this repository's glossary contract (no undefined capitalized
   phrase, a plain `_Avoid_` line).
8. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of the spec command guide,
   and never hand-edit it.
9. MUST NOT touch any skill, the Glossary Declaration grammar, the bold
   candidate rule, the glossary horizon, the binding rule, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize today's reader with the form table.
- [ ] Widen the reader and add the exported form text and command.
- [ ] Change the two fix texts and add the agreement and message tests.
- [ ] Update the spec command guide and `CONTEXT.md`, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Every accepted form defines the term for `SC-GLOSSARY-MISSING` and the bold coverage; every rejected row still does not.
- [ ] The rendered command exits 0 on exactly the rows the reader accepts.
- [ ] The two fix texts name the forms and the command; the guide and `CONTEXT.md` state the forms.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/speccheck/glossary.go`
- interface: `internal/speccheck/glossary_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `speccheck.GlossaryDefinitionForms`, `speccheck.GlossaryVerificationCommand` and the glossary findings — consumers: `./internal/speccheck`, `./internal/docscontract`
- The Behavior Surface Record and the glossary contract over `CONTEXT.md` — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestGlossaryDefinitionFormsAreRead|TestGlossaryVerificationCommandAgreesWithTheDefinitionReader|TestGlossaryFindingsNameTheDefinitionForms)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGlossaryDefinitionFormsAreRead TestGlossaryVerificationCommandAgreesWithTheDefinitionReader TestGlossaryFindingsNameTheDefinitionForms; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/speccheck` — expected: exit 0. Before this Task none of the three tests exists, so the first pass line is missing; afterwards they pass with the whole Spec Check suite.
- `grep -F -- '**Glossary Definition**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Glossary Definition\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'A Glossary Definition is a line that starts with' || { echo 'missing forms in the spec guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary has no such definition, so the first check fails; afterwards the glossary and the guide carry the rule and the documentation contracts, the glossary contract and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Glossary
- `_techspec.md` → Current behavior 1-5; API Contracts 1-4; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1
