---
type: fix
status: promoted
created: 2026-11-02
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
reason: null
---

# A glossary definition written as a list item is not read

## Symptom

`roundfix spec check` reads a glossary term as defined only when its line
starts with `**Term**:`. A glossary that writes its vocabulary as list items,
`- **Term** — definition`, never defines the term: `SC-GLOSSARY-MISSING`
fails the QA gate's precondition after the Task that wrote the term passed its
own Verification, the Run ends Unresolved, and nothing names the form the
check wanted. Reported by pantheon on roundfix 0.128.0.

## Where

`internal/speccheck/glossary.go`, `glossaryDefinitionPattern`, read by
`glossaryFiles` for every glossary code; the `SC-GLOSSARY-MISSING` fix text;
the `write-prd` glossary guide and the `write-tasks` Glossary requirement.

## Expected

A list-item or em-dash definition counts, the finding names the accepted
forms, and the Verification `write-tasks` asks for passes only on a form the
check accepts.

## Evidence

The Pantheon report, as the operator reproduced it from the ai-memory message
(verbatim below):

> # Pantheon report — SC-GLOSSARY-MISSING only accepts "**Term**:" on its own line
>
> From: Pantheon (gesttione-solutions-pantheon), ai-memory message 01a19a04, 2026-11-02 02:09 UTC. Type: fix. Nature: measured.
>
> Observation: the SC-GLOSSARY-MISSING detector (and therefore glossary candidate coverage) only accepts a term definition written as its own line starting with "**Term**:". Pantheon's CONTEXT.md, like many glossaries, lists vocabulary as list items "- **Term** — definition". In that format the term never counts as defined.
>
> Evidence (roundfix 0.128.0, `spec check --strict` in a worktree of the Run branch of Pantheon Spec 0055):
> - "**Acervo**: definition" (own line, no marker) → 0 findings;
> - "- **Acervo**: definition" → 2 SC-GLOSSARY-MISSING findings (one in the PRD, one in the TechSpec);
> - "### Acervo" → 2 findings;
> - "**Acervo** — definition" → 2 findings.
>
> Effect: the Task that adds the term passes its own Verification (which looks for "**Acervo** —"; write-tasks says "bolded term followed by a colon", but the Pantheon guide uses an em dash), and the QA gate refuses the precondition with verdict fail. The Run is Unresolved and the queue parks in run-unresolved without saying which format is missing. The formats had to be tested by hand.
>
> Suggestions: accept a list item ("- **Term**" followed by ":" or " —"), the most common Markdown glossary format; or document the exact format in write-prd references/glossary.md and in the finding's message. write-tasks says "bolded term followed by a colon", which does not make clear that a list marker breaks detection.
>
> Workaround applied in Pantheon: the Acervo definition became its own line "**Acervo**: ..." before the list, and the Task's Verification was adjusted (commit a12bee8 on the 0055 delivery branch).
>
> Run: run_20261102T014441Z_1d007372e1efa0a7 (Spec 0055), QA report qa-report-2026-11-01.md verdict fail, precondition "SC-GLOSSARY-MISSING: Glossary term \"Acervo\" is missing after its binding Tasks completed".
