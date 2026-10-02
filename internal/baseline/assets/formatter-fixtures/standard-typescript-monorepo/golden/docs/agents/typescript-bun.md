<!-- setup-context-driven:begin id=guide.typescript-bun version=0.0.1 -->

# TypeScript and Bun

These rules govern the repository's TypeScript sources and tests. Code in
another language follows its own guide.

- **mandatory**: Keep TypeScript type errors visible; never hide one to make Verification pass.

- **mandatory**: Use current authoritative documentation before changing TypeScript library APIs, framework configuration, runtime behavior, or dependencies.

- **mandatory**: Read the dependent interfaces and their current contracts before writing or changing tests.

- **prohibited**: Do not make mocks, snapshots of incidental output, private structure, or type assertions the source of correctness.

- **mandatory**: Test observable behavior and explicit failure modes.

- **mandatory**: Type a test fixture that stands for a stored row from the schema's inferred row type, with a complete base row and partial overrides, never from an untyped record.

<!-- setup-context-driven:end id=guide.typescript-bun -->
<!-- setup-context-driven:begin id=guide.bun version=0.0.1 -->

# Bun

These rules govern the Bun workspace: the packages under the root
`package.json`. A toolchain for another language in the same repository keeps
its own commands.

- **mandatory**: Run `bun add` from the workspace package that owns the dependency.

- **prohibited**: Inside the Bun workspace, do not substitute another JavaScript package manager or runner (npm, pnpm, yarn, npx) or hand-edit the lockfile. A toolchain for another language in the same repository keeps its own package manager.

- **mandatory**: Inside the Bun workspace, use Bun-owned commands for dependency installation, scripts, and lockfile updates, and run tests through the package's `test` script (`bun run test`), never through a bare runner such as `bun test`.

- **mandatory**: Verify that a dependency exists and inspect its current version before adding it.

<!-- setup-context-driven:end id=guide.bun -->
