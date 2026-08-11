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

<!-- source-baseline-entry: clause.core.verification-two-tiers -->
- MUST use the incremental verification command declared by the active Baseline Profile for fast local checks; it answers whether the current change remains valid while reusing safe local state. CI MUST run the complete verification command declared by that Profile from a fresh run; it answers whether the complete tree satisfies the repository contract. If the Profile declares no incremental command, this clause is unmet rather than satisfied by omission.
<!-- /source-baseline-entry: clause.core.verification-two-tiers -->

<!-- source-baseline-entry: clause.core.prohibit-verification-contract-bypass -->
- MUST NOT edit Verification configuration, tests, fixtures, golden files, or generated expectations merely to make a failure disappear. Change them only when the repository contract intentionally changes, and prove the new contract.
<!-- /source-baseline-entry: clause.core.prohibit-verification-contract-bypass -->

<!-- source-baseline-entry: clause.core.ask-before-verification-configuration-change -->
- MUST stop and ask for explicit authority before intentionally changing lint, formatter, typecheck, test-runner, architecture, or Verification configuration.
<!-- /source-baseline-entry: clause.core.ask-before-verification-configuration-change -->

<!-- source-baseline-entry: clause.core.never-let-a-pipe-hide-a-gate -->
- MUST NOT let a pipe hide a gate's exit status. A pipeline exits with its last command's status, so piping Verification into a pager or filter reports that filter's success and lets `&&` proceed over a red gate. Run the gate on its own, capture its status, or redirect it to a file and read the file.
<!-- /source-baseline-entry: clause.core.never-let-a-pipe-hide-a-gate -->

<!-- source-baseline-entry: clause.core.assertion-reads-the-constant -->
- MUST make an assertion read the constant it means. A test that copies a pinned version, digest, or identifier as a literal stops testing the day a legitimate change moves it, sometimes silently. Reference the exported or package constant; when a value must be duplicated, search for every occurrence and change them in the same commit.
<!-- /source-baseline-entry: clause.core.assertion-reads-the-constant -->

<!-- source-baseline-entry: clause.core.require-tooling-authorization -->
- MUST NOT create, edit, rename, move, or delete any linter, formatter, typechecker, test-runner, architecture-checker, build-tool, package-manager, code-generator, or other repository-tooling configuration, script, ignore file, plugin declaration, or version pin without express maintainer authorization. Setup completion, a Profile, a narrower guide, or a generic implementation request does not grant that authorization.
<!-- /source-baseline-entry: clause.core.require-tooling-authorization -->

<!-- source-baseline-entry: clause.core.tooling-commit-choreography -->
- MUST land an authorized tooling change as its own commit, separate from the record that authorizes it. The express authorization record with its exact bounded paths, and any prerequisite fix repairing something already red before the change, are each their own commit landing before the authorized commit, in either relative order. A consequent fix, which only becomes necessary because the authorized change made something else stale, is its own commit landing after it; it cannot precede the cause that created it. Either kind folded into the authorized commit fails the tooling-authority gate. Prefer no consequent fix at all: a change's declared scope should include the tests its own edit invalidates.
<!-- /source-baseline-entry: clause.core.tooling-commit-choreography -->

<!-- source-baseline-entry: clause.core.activate-matching-skills -->
- MUST activate every matching required skill before governed work. When one skill has distinct active-surface triggers, retain and follow each trigger.
<!-- /source-baseline-entry: clause.core.activate-matching-skills -->

<!-- source-baseline-entry: clause.core.use-conventional-commits -->
- MUST use the governing Conventional Commits workflow before staging changes, writing commit messages, or preparing pull request titles.
<!-- /source-baseline-entry: clause.core.use-conventional-commits -->

<!-- source-baseline-entry: clause.core.use-github-pr-workflow -->
- MUST use the governing pull request workflow before preparing, opening, updating, or handing off a pull request.
<!-- /source-baseline-entry: clause.core.use-github-pr-workflow -->

<!-- source-baseline-entry: clause.core.request-pull-request-review -->
- MUST ask a hand-opened pull request for its own review when it changes code: automatic review is off by configuration, so a pull request opened directly gets none and its review check reports that automatic review is disabled, which reads like a pass. Put the review marker in the pull request description when the pull request is opened rather than adding it afterwards; a one-shot review comment is invalidated by any later push while the check still reads green. Before merging, read the review's own result against the head that will land and treat an absent or stale result as a block.
<!-- /source-baseline-entry: clause.core.request-pull-request-review -->

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
<!-- source-baseline-entry: clause.core.request-review-explicitly -->
- MUST request the review explicitly as part of opening a pull request that changes code, when the repository's review automation does not review it on its own, and MUST NOT read a skipped or absent review check as approval.
<!-- /source-baseline-entry: clause.core.request-review-explicitly -->

<!-- source-baseline-entry: clause.core.conventional-commit-titles -->
- MUST write commit subjects and pull request titles as Conventional Commits subjects, following the scope policy the repository's commit configuration declares. A squash merge often takes the pull request title as the final commit message, so the title is held to the same contract as a commit.
<!-- /source-baseline-entry: clause.core.conventional-commit-titles -->

<!-- source-baseline-entry: clause.core.plan-the-release-first -->
- MUST start release work with the read-only release plan before any changelog, version, tag, push, package, asset, or published-release mutation. A generic release request authorizes only a conclusive patch plan; minor, major, version-zero breaking, and manual classification outcomes require the maintainer decisions recorded in the repository's release runbook.
<!-- /source-baseline-entry: clause.core.plan-the-release-first -->

<!-- source-baseline-entry: clause.core.ask-user-answerable-decisions -->
- MUST NOT guess a decision the user can answer cheaply, and MUST ask it so it is answerable without rereading the session: what was found and why it forces a choice, options enumerated by number or letter with the consequence of each, and a named recommendation. Never bundle questions, ask in prose, or compress an option until only its author can read it.
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

<!-- source-baseline-entry: clause.core.research-authoritative-external-sources -->
For external APIs and libraries, use current authoritative documentation through the profile's declared documentation skill.
<!-- /source-baseline-entry: clause.core.research-authoritative-external-sources -->

<!-- source-baseline-entry: clause.core.research-local-code-locally -->
Search repository files with local code-search tools.
<!-- /source-baseline-entry: clause.core.research-local-code-locally -->

<!-- source-baseline-entry: clause.core.use-declared-external-research-fallback -->
When authoritative documentation cannot answer an external question, use the profile's declared external web-research fallback with varied searches and verify conclusions against primary sources.
<!-- /source-baseline-entry: clause.core.use-declared-external-research-fallback -->

<!-- source-baseline-entry: clause.core.prohibit-external-research-for-local-code -->
Do not use external research tools to discover or infer local repository code or behavior.
<!-- /source-baseline-entry: clause.core.prohibit-external-research-for-local-code -->
