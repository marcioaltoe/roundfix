<!-- source-baseline-entry: clause.core.keep-root-compact -->
- MUST keep root agent instructions as short mandatory pointers. Setup owns only marked baseline content; preserve repository-authored bytes outside setup markers and keep repository-specific architecture and policy in repository-owned documents.
<!-- /source-baseline-entry: clause.core.keep-root-compact -->

<!-- source-baseline-entry: clause.core.fix-root-causes -->
- MUST fix root causes.
<!-- /source-baseline-entry: clause.core.fix-root-causes -->

<!-- source-baseline-entry: clause.core.prohibit-verification-workarounds -->
- MUST NOT suppress diagnostics, weaken assertions, swallow errors, add timing hacks, or bypass a required check to produce a passing result.
<!-- /source-baseline-entry: clause.core.prohibit-verification-workarounds -->

<!-- source-baseline-entry: clause.core.keep-follow-ups-outside-slice -->
- MUST keep follow-up work outside the current slice; record it for later instead of expanding the active change.
<!-- /source-baseline-entry: clause.core.keep-follow-ups-outside-slice -->

<!-- source-baseline-entry: clause.core.record-acceptance-evidence -->
- MUST record the commands and outcomes that provide fresh evidence for every acceptance criterion.
<!-- /source-baseline-entry: clause.core.record-acceptance-evidence -->

<!-- source-baseline-entry: clause.core.require-fresh-evidence -->
- MUST use fresh evidence from the current worktree before claiming work complete, fixed, passing, ready, committed, or delivered. A narrower check supports only the behavior it exercised.
<!-- /source-baseline-entry: clause.core.require-fresh-evidence -->

<!-- source-baseline-entry: clause.core.run-selected-verification -->
- MUST run the selected repository Verification before completion claims. Every failure blocks completion and must be reported with the command and actionable diagnostic.
<!-- /source-baseline-entry: clause.core.run-selected-verification -->

<!-- source-baseline-entry: clause.core.prohibit-verification-contract-bypass -->
- MUST NOT edit Verification configuration, tests, fixtures, golden files, or generated expectations merely to make a failure disappear. Change them only when the repository contract intentionally changes, and prove the new contract.
<!-- /source-baseline-entry: clause.core.prohibit-verification-contract-bypass -->

<!-- source-baseline-entry: clause.core.ask-before-verification-configuration-change -->
- MUST stop and ask for explicit authority before intentionally changing lint, formatter, typecheck, test-runner, architecture, or Verification configuration.
<!-- /source-baseline-entry: clause.core.ask-before-verification-configuration-change -->

<!-- source-baseline-entry: clause.core.activate-matching-skills -->
- MUST activate every matching required skill before governed work. When one skill has distinct active-surface triggers, retain and follow each trigger.
<!-- /source-baseline-entry: clause.core.activate-matching-skills -->

<!-- source-baseline-entry: clause.core.use-conventional-commits -->
- MUST use the governing Conventional Commits workflow before staging changes, writing commit messages, or preparing pull request titles.
<!-- /source-baseline-entry: clause.core.use-conventional-commits -->

<!-- source-baseline-entry: clause.core.use-github-pr-workflow -->
- MUST use the governing pull request workflow before preparing, opening, updating, or handing off a pull request.
<!-- /source-baseline-entry: clause.core.use-github-pr-workflow -->

<!-- source-baseline-entry: clause.core.write-generated-guidance-in-english -->
- MUST write generated repository guidance, identifiers, headings, and examples in English. Preserve repository-authored language outside setup-owned markers.
<!-- /source-baseline-entry: clause.core.write-generated-guidance-in-english -->

<!-- source-baseline-entry: contract.research.sequence -->
## External research procedure

1. Search repository code and local documentation with local code-search tools.
2. For external APIs and libraries, consult current authoritative documentation through the selected documentation capability.
3. If authoritative documentation cannot answer the external question, use the selected broad research capability with three to seven varied searches.
4. Verify conclusions against primary sources and cite the sources used.
5. Never use external research to discover or infer local repository code or behavior.
<!-- /source-baseline-entry: contract.research.sequence -->

<!-- source-baseline-entry: clause.core.follow-dependency-workflow -->
- MUST use the repository's declared package manager and lockfile workflow. Add or upgrade a dependency only for a named job the existing stack cannot perform, and keep manifest and lockfile changes together.
<!-- /source-baseline-entry: clause.core.follow-dependency-workflow -->

<!-- source-baseline-entry: clause.core.review-new-dependencies -->
- MUST review every new dependency for necessity, provenance, maintenance, and security before delivery.
<!-- /source-baseline-entry: clause.core.review-new-dependencies -->

<!-- source-baseline-entry: clause.core.preserve-git-scope -->
- MUST inspect repository status before staging or delivery and preserve unrelated work.
<!-- /source-baseline-entry: clause.core.preserve-git-scope -->

<!-- source-baseline-entry: clause.core.ask-before-destructive-git -->
- MUST stop and ask for explicit authority before destructive Git operations that discard, overwrite, or remove work.
<!-- /source-baseline-entry: clause.core.ask-before-destructive-git -->

<!-- source-baseline-entry: clause.core.ask-before-delivery -->
- MUST stop and ask for explicit authority before committing, pushing, creating a branch, or opening a pull request when that authority is not already explicit.
<!-- /source-baseline-entry: clause.core.ask-before-delivery -->

<!-- source-baseline-entry: clause.core.ask-user-answerable-decisions -->
- MUST NOT guess a decision the user can answer cheaply. Ask through the available user-interaction tool, or ask plainly and stop when no such tool exists.
<!-- /source-baseline-entry: clause.core.ask-user-answerable-decisions -->

<!-- source-baseline-entry: clause.core.prohibit-secret-exposure -->
- MUST NOT read, print, commit, or generate secrets. Keep credentials and environment-specific values in the repository's existing secure configuration boundary, and do not invent authentication, authorization, database, transport, or deployment policy.
<!-- /source-baseline-entry: clause.core.prohibit-secret-exposure -->

<!-- source-baseline-entry: recommendation.capability.firecrawl -->
- RECOMMENDED: provide a structured web-content extraction capability for research that authoritative documentation and broad search cannot complete.
<!-- /source-baseline-entry: recommendation.capability.firecrawl -->

<!-- source-baseline-entry: recommendation.tool.rtk -->
- RECOMMENDED: provide an output-filtering command wrapper that preserves exit status and keeps Agent diagnostics compact.
<!-- /source-baseline-entry: recommendation.tool.rtk -->

<!-- source-baseline-entry: recommendation.tool.rg -->
- RECOMMENDED: provide ripgrep for deterministic local text and file search.
<!-- /source-baseline-entry: recommendation.tool.rg -->
