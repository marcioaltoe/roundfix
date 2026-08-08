<!-- source-baseline-entry: clause.context.read-domain-contract -->
- MUST read the repository's selected domain context and relevant accepted ADRs before naming domain concepts or changing behavior. Flag conflicts instead of silently overriding repository decisions.
<!-- /source-baseline-entry: clause.context.read-domain-contract -->

<!-- source-baseline-entry: clause.domain.canonical-language -->
- MUST use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.
<!-- /source-baseline-entry: clause.domain.canonical-language -->

<!-- source-baseline-entry: clause.domain.layout-decision -->
- MUST follow the repository's declared single-context or multi-context layout. Setup can require that decision but cannot infer bounded contexts from directory names.
<!-- /source-baseline-entry: clause.domain.layout-decision -->

<!-- source-baseline-entry: clause.domain.glossary-currency -->
- MUST check, at the close of a Spec, feature, refactor, or fix, whether the work introduced, changed, or retired a term the glossary should carry, and update the domain context when it did. The check is obliged; the update follows only from a check that found something, and neither waits for human interaction.
<!-- /source-baseline-entry: clause.domain.glossary-currency -->
