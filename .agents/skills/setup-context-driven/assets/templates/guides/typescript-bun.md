# TypeScript and Bun

Use current documentation before changing library APIs, framework
configuration, or runtime behavior. Prefer repository-defined Bun commands and
keep the lockfile consistent with dependency changes.

Tests must assert observable behavior and failure modes. Avoid tests that only
exercise mocks, implementation details, or snapshots of incidental output.
