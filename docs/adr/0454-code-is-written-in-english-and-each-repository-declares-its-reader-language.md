---
status: accepted
created_at: 2026-11-27T17:08:15Z
updated_at: 2026-11-27T17:08:15Z
deprecated_at: null
superseded_by: null
knowledge: true
---

# Code is written in English and each repository declares its Reader Language

The Baseline carries the Language Rule so adopters share a convention for new
code and infrastructure names while declaring the language their users read.
Business terms without a faithful translation keep their original form.

## Decision

Four mandatory clauses live in `rule.core.code-and-reader-language` in `core`:
English names, untranslated terms, changes to new or touched names, and the
Reader Language for product text. Large renames require their own Spec; bulk
renames require a request. Specs, ADRs, agent guides, commits and pull requests
keep their existing language.

`language.reader` is a required string decision. Its suggestion is `English`,
which preserves ADR-0118's decision shape and produces an actionable required
decision stop. A suggestion is adopted only explicitly; it never answers for
a repository. Each built-in profile lists the decision.

The existing mandatory canonical-terms clause in `context-workflow` keeps its
identity and both original sentences. Two inserted sentences require a Name
in Code when a term is translated or kept without accents. Its Glossary
Definition records `_In code_:` followed by the name in backticks, after the
definition and before `_Avoid_`; the introducing or changing Spec writes it.

The approved design also provides a Code Language Exception: one optional
structured decision `language.code_exception`, with language and source and no
default. No source means no exception, and the exception covers only the
codebase that existed when its owner decided to keep its language; it cannot
start a new codebase in another language. English needs no exception. The
source is copied text, not a checked path: decisions may live outside a file,
and checking a path would neither prove its meaning nor keep validity stable.
A presence gate selects the exception clause instead of the three English,
untranslated-term and rename clauses; the Reader Language and Name in Code
still apply.

The new rule stays outside profile `requiredRules`, leaving Source Baseline
rows intact. The Pre-PR Review checks new names and the QA gate checks text
users read. These contextual judgments need no new deterministic checker.

## Considered options

- Put clauses in the existing English-guidance rule: rejected because it
  requires Source Baseline rows; one enlarged clause cannot be gated in parts.
- Add a neighbouring domain clause: rejected because it repeats the rule and
  needs another Source Baseline row. Insertion preserves the recorded words.
- Omit a Reader Language suggestion: rejected because it breaks the decision
  shape and yields invalid-manifest behaviour instead of an actionable stop.
- Use two exception strings or an optional enum: rejected because strings need
  a source default and an enum cannot carry both language and source.
- Check the source as a path: rejected because any file would pass while a
  valid decision outside a file would fail.
- Change the update's next-action sentence: rejected; existing guidance teaches
  explicit answers without changing that transcript.

## Consequences

A repository-owned profile composing `core` must add `language.reader` to its
entry decisions. Existing names change only when touched. The mandatory domain
clause retains its identity and force through updates, and the repository's
Reader Language appears in the opening lines of its agent guide.
