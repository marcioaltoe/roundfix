<!-- source-baseline-entry: clause.frontend.read-design-contract-before-ui-work -->
- MUST read the repository-owned design contract before writing, modifying, or reviewing UI code, and treat it as the selected design authority.
<!-- /source-baseline-entry: clause.frontend.read-design-contract-before-ui-work -->

<!-- source-baseline-entry: clause.frontend.organize-by-system -->
- MUST organize frontend feature code by domain system. Each system exposes one public boundary while its internal components, hooks, queries, routes, and state import each other directly.
<!-- /source-baseline-entry: clause.frontend.organize-by-system -->

<!-- source-baseline-entry: clause.frontend.public-system-boundary -->
- MUST import another system through that system's public boundary instead of reaching into its internal modules.
<!-- /source-baseline-entry: clause.frontend.public-system-boundary -->

<!-- source-baseline-entry: clause.frontend.test-user-visible-behavior -->
- MUST test user-visible roles, labels, text, state changes, loading, error, and empty states.
<!-- /source-baseline-entry: clause.frontend.test-user-visible-behavior -->

<!-- source-baseline-entry: clause.frontend.prohibit-incidental-ui-assertions -->
- MUST NOT make component internals or incidental CSS structure the source of correctness.
<!-- /source-baseline-entry: clause.frontend.prohibit-incidental-ui-assertions -->

<!-- source-baseline-entry: clause.frontend.inspect-runnable-ui -->
- MUST inspect significant local UI changes through the available browser when the target is runnable.
<!-- /source-baseline-entry: clause.frontend.inspect-runnable-ui -->

<!-- source-baseline-entry: recommendation.frontend.style-values -->
- RECOMMENDED: when no repository design contract selects style values, present portable defaults as a decision and adopt them only after confirmation.
<!-- /source-baseline-entry: recommendation.frontend.style-values -->
