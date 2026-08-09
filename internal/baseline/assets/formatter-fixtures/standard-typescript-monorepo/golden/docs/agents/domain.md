<!-- setup-context-driven:begin id=guide.domain version=0.0.1 -->

# Domain docs

This repository uses a single `CONTEXT.md`.

- **mandatory**: Read the repository's selected domain context and relevant accepted ADRs before naming domain concepts or changing behavior. Flag conflicts instead of silently overriding repository decisions.

- **mandatory**: Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.

- **mandatory**: At the close of a Spec, feature, refactor, or fix, check whether the work introduced, changed, or retired a term the glossary should carry, and update the domain context through `domain-modeling` when it did. The check is what is obliged; the update follows only from a check that found something, and neither waits for human interaction. Reach for `grilling` only when a term is ambiguous enough to need sharpening before it is written down.

- **mandatory**: Follow the repository's declared single-context or multi-context layout. Setup can require that decision but cannot infer bounded contexts from directory names.

## Identifier strategy

Use UUID version 7 for new project-owned Internal Identifiers only.

Preserve external provider identifiers, protocol identifiers, natural keys, and business codes according to their source contracts.

<!-- setup-context-driven:end id=guide.domain -->
