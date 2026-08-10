<!-- setup-context-driven:begin id=guide.agent-instructions version=0.0.1 -->

# Agent instructions

This setup-owned guide defines the portable baseline. Repository authors own
project-specific extensions outside setup markers and may add stricter rules.
The selected repository Verification is `make verify`. Agent-created
branches **MUST** use the selected `ma/` prefix.

- **mandatory**: Keep root agent instructions as short mandatory pointers. Setup owns only marked baseline content; preserve repository-authored bytes outside setup markers and keep project-specific architecture and policy in repository-owned documents.

- **mandatory**: Fix root causes.

- **prohibited**: Do not suppress diagnostics, weaken assertions, swallow errors, add timing hacks, or bypass a required check to produce a passing result.

- **mandatory**: Keep follow-up work outside the current slice; record it for later instead of expanding the active change.

- **mandatory**: Record the commands and outcomes that provide fresh evidence for every acceptance criterion.

- **mandatory**: A hand-opened pull request asks for its own review when it changes code. Automatic review is off by configuration, so a pull request opened directly gets none, and its review check reports that automatic review is disabled — which reads like a pass. Put the review marker in the pull request description when the pull request is opened rather than adding it afterwards; a one-shot review comment is invalidated by any later push while the check still reads green. Before merging, read the review's own result against the head that will land, and treat an absent or stale result as a block rather than a pass.

- **mandatory**: Use fresh evidence from the current worktree before claiming work complete, fixed, passing, ready, committed, or delivered. A narrower check supports only the behavior it exercised.

- **mandatory**: Run the selected repository Verification before completion claims. Treat every failure as blocking and report the command plus its actionable diagnostic.

- **stop-and-ask**: Stop and ask for explicit authority before intentionally changing lint, formatter, typecheck, test-runner, architecture, or Verification configuration.

- **mandatory**: An assertion reads the constant it means. A test that copies a pinned version, digest, or identifier as a literal stops testing the day a legitimate change moves it, sometimes silently. Reference the exported or package constant; when a value must be duplicated, search for every occurrence and change them in the same commit.

- **mandatory**: Never let a pipe hide a gate's exit status. A pipeline exits with its last command's status, so piping Verification into a pager or filter reports that filter's success and lets `&&` proceed over a red gate. Run the gate on its own, capture its status, or redirect it to a file and read the file.

- **prohibited**: Do not edit verification configuration, tests, fixtures, golden files, or generated expectations merely to make a failure disappear. Change them only when the repository contract intentionally changes, and prove the new contract.

- **prohibited**: Do not create, edit, rename, move, or delete any linter, formatter, typechecker, test-runner, architecture-checker, build-tool, package-manager, code-generator, or other repository-tooling configuration, script, ignore file, plugin declaration, or version pin without express maintainer authorization. Setup completion, a Profile, a narrower guide, or a generic implementation request does not grant that authorization.

- **mandatory**: Land an authorized tooling change as its own commit, separate from the record that authorizes it. The express authorization record with its exact bounded paths, and any prerequisite fix repairing something already red before the change, are each their own commit landing before the authorized commit, in either relative order. A consequent fix, which only becomes necessary because the authorized change made something else stale, is its own commit landing after it; it cannot precede the cause that created it. Either kind folded into the authorized commit fails the tooling-authority gate. Prefer no consequent fix at all: a change's declared scope should include the tests its own edit invalidates.

- **mandatory**: Write generated repository guidance, identifiers, headings, and examples in English. Preserve repository-authored language outside setup-owned markers.

- **prohibited**: Do not use external research tools to discover or infer local repository code or behavior.

- **mandatory**: For external APIs and libraries, use current authoritative documentation through the profile's declared documentation skill.

- **mandatory**: Search repository files with local code-search tools.

- **mandatory**: When authoritative documentation cannot answer an external question, use the profile's declared external web-research fallback with varied searches and verify conclusions against primary sources.

- **mandatory**: Use the repository's declared package manager and lockfile workflow. Add or upgrade a dependency only for a named job the existing stack cannot perform, and keep manifest and lockfile changes together.

- **mandatory**: Review every new dependency for necessity, provenance, maintenance, and security before delivery.

- **stop-and-ask**: Stop and ask for explicit authority before committing, pushing, creating a branch, or opening a pull request when that authority is not already explicit.

- **stop-and-ask**: Stop and ask for explicit authority before destructive Git operations that discard, overwrite, or remove work.

- **mandatory**: Write commit subjects and pull request titles as Conventional Commits subjects, and follow the scope policy the repository's commit configuration declares. A squash merge often takes the pull request title as the final commit message, so the title is held to the same contract as a commit.

- **mandatory**: Release work starts with the read-only release plan before any changelog, version, tag, push, package, asset, or published-release mutation. A generic release request authorizes only a conclusive patch plan; minor, major, version-zero breaking, and manual classification outcomes require the maintainer decisions recorded in the repository's release runbook.

- **mandatory**: Inspect repository status before staging or delivery and preserve unrelated work.

- **mandatory**: When a pull request changes code and the repository's review automation does not review it on its own, request the review explicitly as part of opening the pull request rather than as a follow-up. Never read a skipped or absent review check as approval.

- **stop-and-ask**: Never guess a decision the user can answer cheaply. Ask through the runtime's structured user-interaction tool when it has one; when it has none, reproduce that shape by hand. Either way the question must be answerable without rereading the session: state what was found and why it forces a choice, then enumerate options by number or letter, each with the consequence of picking it, and name a recommendation. Never bundle several questions into one message, ask in open prose, or compress an option until only its author can tell what it means — an under-explained question is answered by guess, which costs more than not asking.

- **prohibited**: Never read, print, commit, or generate secrets. Keep credentials and environment-specific values in the repository's existing secure configuration boundary, and do not invent authentication, authorization, database, transport, or deployment policy.

<!-- setup-context-driven:end id=guide.agent-instructions -->
