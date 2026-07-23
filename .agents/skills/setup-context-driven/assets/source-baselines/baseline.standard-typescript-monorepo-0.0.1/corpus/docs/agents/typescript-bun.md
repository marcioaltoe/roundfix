<!-- source-baseline-entry: clause.typescript.read-current-authoritative-docs -->
- MUST use current authoritative documentation before changing TypeScript library APIs, framework configuration, runtime behavior, or dependencies.
<!-- /source-baseline-entry: clause.typescript.read-current-authoritative-docs -->

<!-- source-baseline-entry: clause.typescript.keep-type-errors-visible -->
- MUST keep type errors visible and preserve the repository's package and lockfile workflow.
<!-- /source-baseline-entry: clause.typescript.keep-type-errors-visible -->

<!-- source-baseline-entry: clause.typescript.inspect-dependent-interfaces-before-tests -->
- MUST read dependent interfaces and their current contracts before writing or changing tests.
<!-- /source-baseline-entry: clause.typescript.inspect-dependent-interfaces-before-tests -->

<!-- source-baseline-entry: clause.typescript.test-observable-failure-modes -->
- MUST test observable behavior and explicit failure modes.
<!-- /source-baseline-entry: clause.typescript.test-observable-failure-modes -->

<!-- source-baseline-entry: clause.typescript.prohibit-incidental-test-oracles -->
- MUST NOT make mocks, snapshots of incidental output, private structure, or type assertions the source of correctness.
<!-- /source-baseline-entry: clause.typescript.prohibit-incidental-test-oracles -->

<!-- source-baseline-entry: clause.bun.use-bun-owned-commands -->
- MUST use Bun-owned commands for dependency installation, scripts, tests, and lockfile updates.
<!-- /source-baseline-entry: clause.bun.use-bun-owned-commands -->

<!-- source-baseline-entry: clause.bun.verify-dependency-before-add -->
- MUST verify that a dependency exists and inspect its current version before adding it.
<!-- /source-baseline-entry: clause.bun.verify-dependency-before-add -->

<!-- source-baseline-entry: clause.bun.add-from-owning-workspace -->
- MUST run `bun add` from the workspace package that owns the dependency.
<!-- /source-baseline-entry: clause.bun.add-from-owning-workspace -->

<!-- source-baseline-entry: clause.bun.prohibit-other-package-managers -->
- MUST NOT substitute another package manager or hand-edit the lockfile.
<!-- /source-baseline-entry: clause.bun.prohibit-other-package-managers -->

<!-- source-baseline-entry: clause.bun.block-warnings-when-profile-treats-them-as-errors -->
- MUST treat every warning reported by Verification as blocking when the selected profile declares warnings as errors.
<!-- /source-baseline-entry: clause.bun.block-warnings-when-profile-treats-them-as-errors -->
