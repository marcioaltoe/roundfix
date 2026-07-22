<!-- setup-context-driven:begin id=guide.agent-instructions version=3 -->

# Agent instructions

This setup-owned guide defines the portable baseline. Repository authors own
project-specific extensions outside setup markers and may add stricter rules.
The selected repository Verification is `python3 -B .formatter-fixture-verify.py`.

- **mandatory**: Keep root agent instructions as short mandatory pointers. Setup owns only marked baseline content; preserve repository-authored bytes outside setup markers and keep project-specific architecture and policy in repository-owned documents.

- **mandatory**: Fix root causes.

- **prohibited**: Do not suppress diagnostics, weaken assertions, swallow errors, add timing hacks, or bypass a required check to produce a passing result.

- **mandatory**: Keep follow-up work outside the current slice; record it for later instead of expanding the active change.

- **mandatory**: Record the commands and outcomes that provide fresh evidence for every acceptance criterion.

- **mandatory**: Use fresh evidence from the current worktree before claiming work complete, fixed, passing, ready, committed, or delivered. A narrower check supports only the behavior it exercised.

- **mandatory**: Run the selected repository Verification before completion claims. Treat every failure as blocking and report the command plus its actionable diagnostic.

- **stop-and-ask**: Stop and ask for explicit authority before intentionally changing lint, formatter, typecheck, test-runner, architecture, or Verification configuration.

- **prohibited**: Do not edit verification configuration, tests, fixtures, golden files, or generated expectations merely to make a failure disappear. Change them only when the repository contract intentionally changes, and prove the new contract.

- **mandatory**: Write generated repository guidance, identifiers, headings, and examples in English. Preserve repository-authored language outside setup-owned markers.

- **prohibited**: Do not use external research tools to discover or infer local repository code or behavior.

- **mandatory**: For external APIs and libraries, use current authoritative documentation through the profile's declared documentation skill.

- **mandatory**: Search repository files with local code-search tools.

- **mandatory**: When authoritative documentation cannot answer an external question, use the profile's declared external web-research fallback with varied searches and verify conclusions against primary sources.

- **mandatory**: Use the repository's declared package manager and lockfile workflow. Add or upgrade a dependency only for a named job the existing stack cannot perform, and keep manifest and lockfile changes together.

- **mandatory**: Review every new dependency for necessity, provenance, maintenance, and security before delivery.

- **stop-and-ask**: Stop and ask for explicit authority before committing, pushing, creating a branch, or opening a pull request when that authority is not already explicit.

- **stop-and-ask**: Stop and ask for explicit authority before destructive Git operations that discard, overwrite, or remove work.

- **mandatory**: Inspect repository status before staging or delivery and preserve unrelated work.

- **stop-and-ask**: Never guess a decision the user can answer cheaply. Ask through the available user-interaction tool, or ask plainly and stop when no such tool exists.

- **prohibited**: Never read, print, commit, or generate secrets. Keep credentials and environment-specific values in the repository's existing secure configuration boundary, and do not invent authentication, authorization, database, transport, or deployment policy.

<!-- setup-context-driven:end id=guide.agent-instructions -->
