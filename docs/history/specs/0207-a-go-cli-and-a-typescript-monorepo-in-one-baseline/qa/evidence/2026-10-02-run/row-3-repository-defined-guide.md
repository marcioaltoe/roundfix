<!-- setup-context-driven:begin id=guide.frontend version=0.0.1 -->

# Frontend

This setup-owned guide supplies portable frontend rules. Project-specific
visual language and interaction decisions belong to the
repository-owned `DESIGN.md`; setup does not invent architecture or product policy.

These rules govern the repository's web frontend workspace. A terminal
interface follows its own guide.

The repository records its own frontend layout, stated in its repository-owned rules.

- **mandatory**: Inspect significant local UI changes through the available browser when the target is runnable.

- **prohibited**: Do not make component internals or incidental CSS structure the source of correctness.

- **mandatory**: Read the repository-owned `DESIGN.md` before writing, modifying, or reviewing UI code, and treat it as the selected design contract.

- **mandatory**: Test user-visible roles, labels, text, state changes, loading, error, and empty states.

- **mandatory**: Organize frontend feature code by the layout the repository's own rules state, and update those rules in the same change that alters the layout.

<!-- setup-context-driven:end id=guide.frontend -->
