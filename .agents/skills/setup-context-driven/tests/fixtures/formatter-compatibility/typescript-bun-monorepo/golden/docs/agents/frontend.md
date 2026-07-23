<!-- setup-context-driven:begin id=guide.frontend version=0.0.1 -->

# Frontend

This setup-owned guide supplies portable frontend rules. Project-specific
visual language and interaction decisions belong to the
repository-owned `DESIGN.md`; setup does not invent architecture or product policy.

- MUST read the repository-owned design contract before writing, modifying, or reviewing UI code, and treat it as the selected design authority.

- MUST organize frontend feature code by domain system. Each system exposes one public boundary while its internal components, hooks, queries, routes, and state import each other directly.

- MUST import another system through that system's public boundary instead of reaching into its internal modules.

- MUST test user-visible roles, labels, text, state changes, loading, error, and empty states.

- MUST NOT make component internals or incidental CSS structure the source of correctness.

- MUST inspect significant local UI changes through the available browser when the target is runnable.

- RECOMMENDED: when no repository design contract selects style values, present portable defaults as a decision and adopt them only after confirmation.

<!-- setup-context-driven:end id=guide.frontend -->
