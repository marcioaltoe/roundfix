# Glossary Declaration

Every PRD carries a `## Glossary` section. Use one entry per term or phrase:

- `adds: **<term>**` declares a domain term the Spec introduces.
- `changes: **<term>**` declares a domain term whose definition the Spec revises.
- `not a term: **<phrase>** — <reason>` declares a bolded phrase that is not a domain term; the reason is required.
- `None.` is valid only as the section's single entry.

Terms and phrases are bolded where they are introduced. A bold candidate is a
two-to-five-word span outside fenced blocks and outside the Glossary section,
with no backtick, digit or parenthesis, no terminal punctuation, and every word
starting with an uppercase letter except the permitted connectors after the
first word. A candidate is covered by an existing glossary definition, its
plural, or a declared term or phrase. An uncovered candidate reports
`SC-GLOSSARY-UNDECLARED`.

A declared `adds` or `changes` term is bound to a non-QA Task that declares the
glossary file under `interface:` or `creates:` and names the exact bolded term
in its Verification. A declared term without that binding reports
`SC-GLOSSARY-UNPLANNED`. A changed term that is not defined in the glossary
also reports `SC-GLOSSARY-UNPLANNED`; declare it as added to fix that case. When
all binding Tasks are completed and the glossary still lacks the term, the
check reports `SC-GLOSSARY-MISSING`.

The glossary files are the root `CONTEXT.md` and `GLOSSARY.md`, plus every
`CONTEXT.md` or `GLOSSARY.md` targeted by a relative Markdown link in the root
`CONTEXT-MAP.md` or `GLOSSARY-MAP.md`. The glossary horizon starts when this
guide is added: a PRD committed at or after that commit must carry a Glossary
Declaration. Before the horizon, the check skips the declaration requirement;
an uncommitted PRD or unreadable history follows the horizon contract.

A domain term the Spec introduces is written by one of the Spec's own Tasks, never left for a later Spec.
