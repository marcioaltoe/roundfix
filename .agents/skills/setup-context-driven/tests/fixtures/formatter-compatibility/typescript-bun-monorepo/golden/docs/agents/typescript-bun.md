<!-- setup-context-driven:begin id=guide.typescript-bun version=3 -->

# TypeScript and Bun

- **mandatory**: Keep type errors visible and preserve the repository's package and lockfile workflow.

- **mandatory**: Use current authoritative documentation before changing TypeScript library APIs, framework configuration, runtime behavior, or dependencies.

- **mandatory**: Read the dependent interfaces and their current contracts before writing or changing tests.

- **prohibited**: Do not make mocks, snapshots of incidental output, private structure, or type assertions the source of correctness.

- **mandatory**: Test observable behavior and explicit failure modes.

<!-- setup-context-driven:end id=guide.typescript-bun -->

<!-- setup-context-driven:begin id=guide.bun version=2 -->

# Bun

- **mandatory**: Run `bun add` from the workspace package that owns the dependency.

- **prohibited**: Do not substitute another package manager or hand-edit the lockfile.

- **mandatory**: Use Bun-owned commands for dependency installation, scripts, tests, and lockfile updates.

- **mandatory**: Verify that a dependency exists and inspect its current version before adding it.

- **mandatory**: When the selected TypeScript/Bun profile treats warnings as errors, every warning reported by Verification blocks completion.

<!-- setup-context-driven:end id=guide.bun -->
