<!-- setup-context-driven:begin id=guide.domain version=0.0.1 -->

# Domain docs

This repository uses a single `CONTEXT.md`.

- **mandatory**: The repository's selected domain context and its accepted ADRs are the base of the CONTEXT-driven workflow. Read both before naming domain concepts or changing behavior, and flag conflicts instead of silently overriding repository decisions.

- **mandatory**: Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.

- **mandatory**: A Spec adds each domain term it introduces to the domain context, and revises each term it changes, through `domain-modeling` in one of its own Tasks, whose Verification names the term. It lists those terms in a `## Glossary` section of its PRD or TechSpec, with any bolded phrase it declares not a domain term. `roundfix spec check` reports a bolded term the domain context lacks and the section does not cover, and the QA gate and `roundfix archive` refuse a Spec whose declared term is still missing. When `domain-modeling` names its glossary file differently, it writes to the selected domain context. Work outside a Spec checks at its close whether it introduced, changed, or retired a term and updates the domain context when it did. Neither the check nor the update waits for human interaction; reach for `grilling` only when a term is ambiguous enough to need sharpening before it is written down.

- **mandatory**: Follow the repository's declared single-context or multi-context layout. Setup can require that decision but cannot infer bounded contexts from directory names.

<!-- setup-context-driven:end id=guide.domain -->
