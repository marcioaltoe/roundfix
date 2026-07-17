# Ubiquitous language in colocated `CONTEXT.md`

## Status

Accepted

## Context

A ubiquitous-language document loses authority when it is far from the code
and decisions it explains. A single central glossary also becomes ambiguous
when a repository contains multiple bounded contexts with distinct meanings
for similar terms. At the same time, forcing domain experts to use code
vocabulary can replace the language used by the business with an artificial
technical dialect.

Architecture Decision Records serve contributors and agents across projects.
Using one documentation language makes those decisions consistently searchable
and reusable without forcing the domain itself to speak that language.

## Decision

In a single-context repository, ubiquitous language lives in `CONTEXT.md` at
the repository root.

In a multi-context repository, `CONTEXT-MAP.md` at the repository root indexes
the bounded contexts and their relationships. Each bounded context keeps its
own `CONTEXT.md` colocated with the code it describes. System-wide ADRs remain
in root `docs/adr/`; a context may keep context-specific ADRs beside its own
code when that ownership boundary is explicit in the Context Map.

Canonical domain terms use the language spoken by the domain experts. When the
code uses a different English identifier, `CONTEXT.md` records that identifier
as an alias rather than replacing the domain term. Standard legal, regulatory,
or protocol vocabulary retains its authoritative form.

Every new or modified ADR filename, title, and prose must be written in
English. This language rule does not require canonical domain terms in
`CONTEXT.md` to be English.

Agents and contributors read the root Context Document or Context Map, the
relevant colocated Context Documents, and applicable ADRs before naming domain
concepts or changing architectural behavior.

## Alternatives Considered

### Keep every glossary under a central documentation directory

This makes glossaries easy to browse together but separates them from the code
whose language they govern, increasing the chance that they drift.

### Keep one root glossary for every repository

This is simple for a single context but mixes distinct languages as a system
grows and makes ownership unclear.

### Require domain terminology to follow English code identifiers

This aligns documentation with code but weakens collaboration with the people
who hold the domain knowledge when they use a different language.

## Consequences

- Changes to a bounded context keep its colocated `CONTEXT.md` and the root
  `CONTEXT-MAP.md` accurate.
- Root and local agent instructions point contributors to the applicable
  Context Documents and ADRs.
- Code can remain English while the glossary preserves the domain experts'
  canonical language through explicit aliases.
- ADRs have one portable documentation language even in projects whose domain
  language is not English.
