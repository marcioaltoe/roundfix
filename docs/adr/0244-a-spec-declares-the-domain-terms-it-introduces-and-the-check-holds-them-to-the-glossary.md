---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec declares the domain terms it introduces, and the check holds them to the glossary

Roundfix was built with the glossary and the ADRs as the base of the
CONTEXT-driven process. By 2026-10-06 the glossary had not changed for four
days while twenty-one Specs archived, and reports had named terms that nobody
wrote down. The domain guide already asked for a glossary check at the close of
a Spec, but nothing read the answer, the autonomous authoring route never
reached `domain-modeling`, and the operator's own briefing had told authors not
to edit the glossary.

A Spec now carries a Glossary Declaration: a `## Glossary` section in its PRD
or TechSpec that lists each domain term the Spec adds or changes and each
bolded phrase it declares not a domain term. The Spec Consistency Check holds
the declaration to the glossary deterministically. A bolded phrase of two to
five capitalized words that the glossary does not define and the declaration
does not cover is a gap. A term the Spec adds or changes needs a non-QA Task
that declares the glossary file and names the term in its Verification. Once
every such Task has completed, the glossary must define the term. Each of these
is an error, so the QA gate's precondition and `roundfix archive` refuse a Spec
that leaves one. The glossary is the repository's `CONTEXT.md`, or
`GLOSSARY.md` under the name the upstream skills adopted on 2026-10-05, plus
any context its map file links.

The section is required only for a Spec whose PRD was committed at or after
the commit that added the glossary guide to the `write-prd` skill, so no Spec
authored earlier is blocked by a rule it could not know. A Spec that already
carries the section is checked whatever its age.

A model judging "is this a domain term" was rejected as the gate: the check
must give the same answer on every run, and the Jev judge stays advisory. Bold is
the convention this repository already uses for glossary terms. A replay of the
bold rule over the PRDs and TechSpecs of Specs 0200 to 0237 found eight phrases
in five Specs, six distinct concept names, none of them in the glossary; the
same replay counting one-word bold labels added twenty-four more, nearly all
labels such as "Review" or "Configuration", so one-word phrases are left to the
declaration. A lowercase term is caught only when its author
declares it, which is why the authoring skills now call `domain-modeling` when
they name a concept.
