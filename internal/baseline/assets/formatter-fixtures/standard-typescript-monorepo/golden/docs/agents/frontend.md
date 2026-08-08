<!-- setup-context-driven:begin id=guide.frontend version=0.0.1 -->

# Frontend

This setup-owned guide supplies portable frontend rules. Project-specific
visual language and interaction decisions belong to the
repository-owned `DESIGN.md`; setup does not invent architecture or product policy.

- **mandatory**: Inspect significant local UI changes through the available browser when the target is runnable.

- **mandatory**: Organize frontend feature code by domain system. Each system exposes one public boundary while its internal components, hooks, queries, routes, and state import each other directly.

- **prohibited**: Do not make component internals or incidental CSS structure the source of correctness.

- **mandatory**: Import another system through that system's public boundary instead of reaching into its internal modules.

- **mandatory**: Read the repository-owned `DESIGN.md` before writing, modifying, or reviewing UI code, and treat it as the selected design contract.

- **mandatory**: Test user-visible roles, labels, text, state changes, loading, error, and empty states.

<!-- setup-context-driven:end id=guide.frontend -->
