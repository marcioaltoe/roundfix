---
task: task_01
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
status: completed
type: backend
complexity: high
---

# Task 01: The Default Branch Head, and a Git Locator is reachable in a CI checkout

## Overview

On 2026-10-18 the CI Verification gate of pull request #580 called a valid Git
Locator invalid (`_prd.md` → Sources 1). The GitHub Actions checkout has no
origin/HEAD and no local `main`, so the default branch was undetermined and
only tags proved reachability. A stale local `main` failed the same way,
because the locator read the local branch before the fetched one.

This Task adds the `origin-branch` detection step and one resolver, the
**Default Branch Head**, and makes the Git Locator read it. It is verifiable on
its own through preflight unit tests and temporary clones that reproduce both
cases.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-a-git-locator-in-a-fresh-record-fails-the-check-in-ci.md`
   (`_prd.md` → Sources 1), by implementing API Contracts 1 to 3 of
   `_techspec.md` (ADR-0295).
2. MUST add the `origin-branch` step to `DetectDefaultBranch` in
   `internal/preflight/default_branch.go` (API Contract 1).
   - Add `DefaultBranchFromOriginBranch DefaultBranchSource = "origin-branch"`.
   - The step runs only after origin/HEAD and the `main`/`master` name match
     both fail. It reads one `git for-each-ref --format=%(refname)
     refs/remotes/origin/main refs/remotes/origin/master` and keeps exact ref
     names only.
   - Exactly one match names the default; zero, two, or a Git error leave it
     undetermined.
   - `IsDefault` and every existing case of `TestDetectDefaultBranch` keep
     their answers.
3. MUST add `DefaultBranchHead` and `ResolveDefaultBranchHead` in the new
   `internal/preflight/default_branch_head.go` (API Contract 2).
   - One `git for-each-ref --format=%(refname)%00%(objectname)
     refs/remotes/origin/<name> refs/heads/<name>` gives both candidates. The
     exact remote-tracking ref wins, and the local ref answers only when it is
     absent.
   - `found=false` with a nil error when detection is undetermined or neither
     ref exists. A Git error is returned wrapped with `%w`, never treated as
     absence.
4. MUST make `ResolveGitLocator` in `internal/spec/git_locator.go` read the
   head from `ResolveDefaultBranchHead` with its existing
   `gitLocatorBranchRunner` (API Contract 3).
   - That head proves reachability, with tags as today, and supplies the
     newest commit that holds a removed path.
   - The local-then-remote loop is removed, and the signature, the error
     wrapping and `GitLocator.Text` stay.
   - `internal/speccheck/history_retained.go` needs no change: it inherits the
     fix through `ResolveGitLocator`.
5. MUST add these tests. Each calls `t.Parallel()` first, compiles for Windows
   and reaches no network. Git repositories are temporary, and the `spec` and
   `speccheck` tests reuse their package's child-process isolation
   (`gitLocatorTestHome`, `retainedTestHome`).
   - `internal/preflight/default_branch_test.go`:
     - new `TestDetectDefaultBranch` cases for one fetched `origin/main`, one
       `origin/master`, both, none and a failing `for-each-ref`;
     - `TestResolveDefaultBranchHeadPrefersTheFetchedBranch`, with the fake
       runner: the remote ref wins over a different local ref, the local ref
       answers alone, undetermined and absent refs give `found=false`, and a
       Git error is returned.
   - `internal/spec/git_locator_test.go`:
     - `TestResolveGitLocatorReachesTheFetchedDefaultInACIClone`. It builds a
       temporary origin with `main`, an older tag and a newer commit, and a
       clone made with `git init`, `git remote add origin`, `git fetch origin
       +refs/heads/*:refs/remotes/origin/* +refs/tags/*:refs/tags/*` and `git
       checkout --detach <pull request head>`. Without origin/HEAD or a local
       `main`, a locator at the newer commit resolves to it.
     - `TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch`: a clone whose
       local `main` is reset to the older tag resolves the same locator, and a
       removed path's newest commit is read from `origin/main`.
     - `TestResolveGitLocatorReadsTheRemoteDefaultBranch` keeps its name. Its
       last assertion changes to: an unpushed local commit on the default
       branch does not win over the fetched head (declared break, ADR-0295).
   - `internal/speccheck/history_retained_test.go`:
     `TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone`.
     An Archive Record cites a locator whose commit lies between two tags, in
     the CI-like clone, and `SC-HISTORY-RETAINED` reports no finding for it.
     Deleting the remote-tracking ref makes the finding return.

   Every other existing test MUST keep its name and pass.
6. MUST document the change:
   - `docs/user-guide/commands/archive.md` replaces "when the local default
     branch or a tag reaches it" with "when the Default Branch Head or a tag
     reaches it". It then says in one sentence that the Default Branch Head is
     the fetched `origin/<default>`, the local branch only without one, and
     `origin/main` or `origin/master` when origin/HEAD is unset.
   - `.agents/skills/roundfix/references/archive.md` replaces "be reachable
     from the default branch or a tag" with "be reachable from the Default
     Branch Head or a tag".

   No text inside any `### QA settlement` section may change.
7. MUST add **Default Branch Head** to `CONTEXT.md` through domain-modeling,
   with ADR-0295's link. Its definition contains "only when no remote-tracking
   branch exists". MUST revise **Git Locator** so it contains "that the Default
   Branch Head or a tag reaches".
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of each edited owned skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`)
   by one patch over the value the tree holds when this Task starts,
   because the record command refuses changed content under an unchanged
   version. The Agent runs these commands itself before handing back: the
   Daemon only runs the Verification and never regenerates a mirror or a
   record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the `origin-branch` step and the Default Branch Head resolver.
- [ ] Make the Git Locator read the resolver.
- [ ] Add the preflight, locator and history-check tests and change the declared assertion.
- [ ] Document the locator rule in the `archive` guide, the reference and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] In a CI-like clone and in a stale-local clone, a locator `origin/main` reaches is valid and `SC-HISTORY-RETAINED` accepts it.
- [ ] With both `origin/main` and `origin/master` and no origin/HEAD, the default stays undetermined.
- [ ] The guide, the reference and the glossary name the Default Branch Head.

## Consumers

- checkpoint: task_02
- `spec.ResolveGitLocator` — consumers: `./internal/speccheck`, `./internal/spec`

## Context

- instruction: `docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md`
- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `internal/speccheck/history_retained.go`
- creates: `internal/preflight/default_branch_head.go`
- interface: `internal/preflight/default_branch.go`
- interface: `internal/preflight/default_branch_test.go`
- interface: `internal/spec/git_locator.go`
- interface: `internal/spec/git_locator_test.go`
- interface: `internal/speccheck/history_retained_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `grep -q 'DefaultBranchFromOriginBranch' internal/preflight/default_branch.go || { printf 'missing origin-branch source\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestDetectDefaultBranch|TestResolveDefaultBranchHeadPrefersTheFetchedBranch)$' ./internal/preflight 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestResolveDefaultBranchHeadPrefersTheFetchedBranch ' || { printf 'missing pass: TestResolveDefaultBranchHeadPrefersTheFetchedBranch\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestResolveGitLocatorReachesTheFetchedDefaultInACIClone|TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch|TestResolveGitLocatorReadsTheRemoteDefaultBranch)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestResolveGitLocatorReachesTheFetchedDefaultInACIClone TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch TestResolveGitLocatorReadsTheRemoteDefaultBranch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone ' || { printf 'missing pass: TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'when the Default Branch Head or a tag reaches it' || { printf 'missing archive guide locator rule\n' >&2; exit 1; }; ! tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'local default branch' || { printf 'archive guide still names the local default branch\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- 'reachable from the Default Branch Head or a tag' || { printf 'missing archive reference locator rule\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Default Branch Head**:' || { printf 'missing glossary term Default Branch Head\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'only when no remote-tracking branch exists' || { printf 'missing Default Branch Head definition\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'that the Default Branch Head or a tag reaches' || { printf 'missing revised glossary term Git Locator\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Git Locator**:' || { printf 'missing glossary term Git Locator\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/preflight ./internal/spec || exit 1; go test -count=1 -run 'History|Locator|Retained|Receipt|Rollup' ./internal/speccheck || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the `origin-branch` source and the new tests do not exist, so the command fails at its first check. Afterwards these checks pass: the resolver and locator tests in temporary clones, the history check, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Stories 1-2; Core Features 1-2; Core Feature 5; Success Metric 1; Success Metric 2; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 10; Invariants 1-2; Current behavior; Testing Approach; Build Order 1
- ADR-0295; ADR-0278; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285; ADR-0265


## Result

Implemented the Task 01 slice for Daemon Verification; status remains
Daemon-owned. The starting worktree contained only the Daemon's status change
in this Task file. No Task Graph or other Task file was edited.

### Implementation and acceptance evidence

- CI-like and stale-local clones: `DetectDefaultBranch` now uses exactly one
  fetched `origin/main` or `origin/master` after its existing detection steps.
  `ResolveDefaultBranchHead` reads both exact candidate refs in one call,
  prefers the fetched ref, falls back to the local ref only when absent,
  distinguishes absent refs from Git errors, and wraps read errors with `%w`.
  `ResolveGitLocator` uses that head for provenance reachability and removed-path
  history. `internal/speccheck/history_retained.go` is unchanged.
- `TestResolveGitLocatorReachesTheFetchedDefaultInACIClone` and
  `TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch` pass with isolated
  temporary origins and local fetches. They cover a detached checkout without
  origin/HEAD or local main, a stale local main, and newest-held-file recovery
  after deletion. The existing remote-default test now checks that an
  unpushed local commit cannot override the fetched head, both with and
  without a preferred commit.
- `TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone`
  passes: an Archive Record's locator is accepted in the CI-like clone and
  with a stale local main. Removing the local branch and fetched main makes
  the expected invalid-locator finding return. The fixture's older tag and
  later sibling release tag do not independently retain the cited commit.
- Ambiguous defaults: focused preflight cases pass for origin/main alone,
  origin/master alone, both, neither, Git failure, exact-ref filtering, and
  earlier-step precedence. The resolver cases pass for remote preference,
  local fallback, undetermined detection, missing refs, exact-ref filtering,
  and wrapped Git errors.
- Documentation: the archive guide and Roundfix archive reference name the
  Default Branch Head. `CONTEXT.md` adds its ADR-0295-linked definition,
  including "only when no remote-tracking branch exists", and revises Git
  Locator reachability. Roundfix's two canonical version fields rose from
  0.1.106 to 0.1.107; mirrors and records were generated by their owners.
  Byte comparisons confirmed canonical/mirror equality and unchanged QA
  settlement sections. A Python scope/wording audit and `git diff --check`
  both exited 0.

### Focused checks and regeneration

- Red regression proof: a temporary `go test -overlay` substituted HEAD's
  original `internal/preflight/default_branch.go` and
  `internal/spec/git_locator.go` while retaining the new tests. Running
  `go test -overlay <temporary-overlay> ./internal/spec -run
  '^TestResolveGitLocator(ReachesTheFetchedDefaultInACIClone|IgnoresAStaleLocalDefaultBranch)$'
  -count=1` exited 1: the CI clone returned no commit and the stale-local
  clone returned the older commit. Working-tree source was never reverted.
- `go test ./internal/preflight -run
  'Test(DetectDefaultBranch|DefaultBranchIsDefault|ResolveDefaultBranchHead)'
  -count=1`: exit 0.
- `go test ./internal/spec ./internal/speccheck -run
  'Test(ResolveGitLocator|HistoryRetained)' -count=1`: exit 0 after the final
  test edits.
- `GOOS=windows GOARCH=amd64 go test -c -o <temporary-output>
  ./internal/<package>`: exit 0 for preflight, spec and speccheck; spec and
  speccheck were compiled again after their final assertion changes.
- `make skills-sync`: exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: exit 0; recorded Roundfix 0.1.107.
- `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  exit 0; regenerated the archive guide's Behavior Surface Record.
- `make baseline-digests`: exit 0, `changed=false`; no derived Baseline file
  changed.

The symlinked writable root prevented sandbox process startup. Authorized
execution used the resolved Run checkout under `/Volumes/Brain-SSD/cache/`;
this was an environment constraint. All test Git repositories and fetches
were temporary and local; no network access was needed.

The declared Verification command and repository Verification were not run;
they remain Daemon-owned. No commit, push or pull request was made. Consumers
owned by Task 02 and storage work owned by Task 03 remain outside this diff.
