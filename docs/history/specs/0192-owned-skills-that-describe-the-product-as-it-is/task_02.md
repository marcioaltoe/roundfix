---
task: task_02
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: completed
type: test
complexity: medium
---

# Task 02: The user guide states what ships, and its links resolve

## Overview

The user guide and the README are what a person reads before running Roundfix. They document a `--qa` flag that `implement` refuses, state an archive destination the command no longer uses, print a Doctor skills count this repository does not produce, leave two commands unnamed, and carry fifteen links that resolve to nothing. This Task adds a repository test file that fails for each of those classes and corrects the guide until it passes. The test reuses `commandPaths` and `undocumentedCommands`, which task_01 added in `internal/docscontract/command_documentation_test.go`.

## Requirements

1. MUST add `internal/docscontract/user_guide_contract_test.go`, in the package and under the `docscontract` build tag the directory's other tests use, with `brokenLinks(repoRoot, file, content string) []string` as the TechSpec's Interfaces section names it. It resolves each relative Markdown link destination from the directory of `file`. It skips lines inside fenced code blocks, destinations with a URL scheme, and fragment-only destinations. It strips a `#fragment`, a `?query`, a title and `<…>` wrapping before resolving.
2. MUST add these tests to that file:
   - `TestEveryCommandIsNamedInTheUserGuide` reads the root help through `cli.Run`, calls `commandPaths` and `undocumentedCommands` from task_01's file over the concatenated text of every `docs/user-guide/*.md`, and fails with the list of unnamed paths.
   - `TestUserGuideLinksResolve` fails with every entry `brokenLinks` returns for each `docs/user-guide/*.md` file and for `README.md`.
   - `TestABrokenUserGuideLinkIsReported` feeds a document holding one link to a missing path, one link to an existing path, one URL and one link inside a code fence, and expects only the missing one.
   - `TestUserGuideNamesNoRefusedQAFlag` fails when any `docs/user-guide/*.md` file or `README.md` holds `--qa` not followed by `-override`.
3. MUST remove `--qa` from the `implement` synopsis in `docs/user-guide/commands.md` and from the `implement` examples in `docs/user-guide/usage.md` and `README.md`. Each of the three files MUST say that the gate is the Spec's authored terminal `qa` Task and contain the words `no flag requests it`. `usage.md` keeps the statement that the report gains a `qa <verdict> — <report path>` line, which the command still prints.
4. MUST state the archive destination the command uses. `docs/user-guide/commands.md` names `docs/history/specs/<slug>/` for the built-in Spec Root and `<specs.root>/_archived/<slug>/` for any other root. The `roundfix archive` row of `docs/user-guide/context-driven-development.md` names `docs/history/specs/` and no longer names `docs/specs/_archived/`.
5. MUST retarget every link `TestUserGuideLinksResolve` reports to the document's current path under `docs/history/`. On the starting main there are fifteen, in `commands.md`, `configuration.md`, `run-database-lifecycle.md` and `usage.md`. They point at `../specs/_archived/`, `../findings/_archived/`, `../specs/0036-doctor-skill-readiness/` and two ADRs now under `docs/history/adr/`. When a target exists nowhere in the repository, remove the link and keep the sentence.
6. MUST replace the three literal lines `skills: ok (39 required: 14 Roundfix-owned, 25 external)` in `commands.md`, `usage.md` and `README.md` with the shape `skills: ok (<required> required: <owned> Roundfix-owned, <external> external)`, and say next to each that the numbers come from the repository's Repository Skill Set.
7. MUST add entries to `docs/user-guide/commands.md` for `roundfix spec audit` and `roundfix storage report`, each written from the command's own `--help`: synopsis, what it reads, what it never changes, and exit codes where the help states them.
8. MUST NOT edit any existing test, any skill, or any file under `docs/agents/`. The adapter versions in the Doctor example are left as they are; another Spec owns them.
9. MUST keep every phrase the existing documentation contract tests require from the guide files. `TestBaselineDocumentationContract`, `TestProfilesDocumentationContractMatchesPublicGuidance` and `TestReleasePlanDocumentationContract` name them.

## Subtasks

- [ ] Add the link helper and the four tests.
- [ ] Remove the refused flag from the guide and the README.
- [ ] Correct the archive destination in both guide files.
- [ ] Retarget or remove every broken link.
- [ ] Show the Doctor skills line as a shape.
- [ ] Document `spec audit` and `storage report` in the commands guide.

## Acceptance Criteria

- [ ] Every command path the root help lists is named as `roundfix <path>` somewhere in `docs/user-guide/`.
- [ ] No relative link in `docs/user-guide/*.md` or `README.md` points at a missing path.
- [ ] A document with one dead link, one live link, one URL and one fenced link yields exactly the dead link.
- [ ] Neither the guide nor the README names `--qa` except as `--qa-override`.
- [ ] The guide states the archive destination for the built-in and for a configured Spec Root.
- [ ] The Doctor skills line is a shape, with no literal count.

## Context

- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/run-database-lifecycle.md`
- interface: `README.md`
- creates: `internal/docscontract/user_guide_contract_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestEveryCommandIsNamedInTheUserGuide|TestUserGuideLinksResolve|TestABrokenUserGuideLinkIsReported|TestUserGuideNamesNoRefusedQAFlag|TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract)$" ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryCommandIsNamedInTheUserGuide TestUserGuideLinksResolve TestABrokenUserGuideLinkIsReported TestUserGuideNamesNoRefusedQAFlag TestBaselineDocumentationContract TestProfilesDocumentationContractMatchesPublicGuidance TestReleasePlanDocumentationContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in "docs/user-guide/commands.md|no flag requests it" "docs/user-guide/usage.md|no flag requests it" "README.md|no flag requests it" "docs/user-guide/commands.md|docs/history/specs/<slug>/" "docs/user-guide/context-driven-development.md|docs/history/specs/" "docs/user-guide/commands.md|<required> required: <owned> Roundfix-owned, <external> external" "docs/user-guide/usage.md|<required> required: <owned> Roundfix-owned, <external> external" "README.md|<required> required: <owned> Roundfix-owned, <external> external" "docs/user-guide/commands.md|roundfix spec audit" "docs/user-guide/commands.md|roundfix storage report"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; for pair in "docs/user-guide/commands.md|39 required" "docs/user-guide/usage.md|39 required" "README.md|39 required" "docs/user-guide/context-driven-development.md|docs/specs/_archived/"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done` — expected: exit 0; before this Task none of the four new tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Interfaces; The user guide
- `_prd.md` → Goal 2; Goal 3; Goal 4; Core Feature 2; Core Feature 6; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Testing Approach 2; Testing Approach 3; Build Order 2

## Result

Implemented the Task 02 documentation contract slice. Task status and the
authored Verification section remain Daemon-owned and unchanged by this Agent.
No existing test, skill, agent guide, Task Graph, or other Task file was edited.

Acceptance evidence:

- Command coverage: added `TestEveryCommandIsNamedInTheUserGuide`, using the
  real root help through `cli.Run` and task_01's helpers. Its first run named
  `spec audit` and `storage report` as missing; after documenting both from
  their own help, the focused test passes.
- Link resolution: added `brokenLinks` and `TestUserGuideLinksResolve`. The
  first run reported all fifteen broken links. All fifteen targets exist
  under `docs/history/` and were retargeted; the focused test now passes for
  every guide file and README.
- Dead-link detection: `TestABrokenUserGuideLinkIsReported` passes with exactly
  the missing destination from a fixture containing a dead link, a live link,
  a URL, and a fenced link. An additional focused fixture proves angle-wrapped
  and plain destinations with query, fragment and title, paths with spaces,
  fragment-only links, URL schemes and tilde fences containing backticks.
- QA flag: the initial check reported four refused-flag mentions. Removed
  them, retained the archive override, and stated in commands, usage and README
  that the gate is the Spec's authored terminal `qa` Task and no flag requests
  it. `TestUserGuideNamesNoRefusedQAFlag` now passes. The usage report retains
  `qa <verdict> — <report path>`.
- Archive destinations: the command reference names `docs/history/specs/<slug>/` for the
  built-in Spec Root and `<specs.root>/_archived/<slug>/` otherwise. The
  context-driven pipeline row names both destinations and no longer names
  `docs/specs/_archived/`. Focused Python assertions confirmed these strings.
- Doctor skills: all three examples use
  `skills: ok (<required> required: <owned> Roundfix-owned, <external> external)`
  and explain that the counts come from the repository's Repository Skill
  Set. Focused Python assertions confirmed the shape and removal of the old
  `39 required` line. Adapter versions were preserved.

Focused checks (not the authored Daemon Verification):

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -count=1 -tags docscontract -run 'TestEveryCommandIsNamedInTheUserGuide|TestUserGuideLinksResolve|TestABrokenUserGuideLinkIsReported|TestBrokenLinksHandles|TestUserGuideNamesNoRefusedQAFlag' ./internal/docscontract`
  — exit 0 after the final test edit. The initial run reproduced the two
  missing commands, fifteen broken links and four refused flag mentions.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -count=1 -tags docscontract -run 'TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract' ./internal/docscontract`
  — exit 0; all three existing documentation contracts were preserved.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go run ./cmd/roundfix spec audit --help`
  and the corresponding `storage report --help` call — both exit 0; synopsis,
  reads, mutation boundaries and stated exit codes were checked against them.
- Focused Python assertions over commands, usage, README and the pipeline row
  — exit 0 for authored QA wording, report line, archive destinations, skills
  shapes and count provenance.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk make verify-incremental`
  — exit 0 with host process-table access on the rerun; formatting, vet,
  repository unit/integration tests, skill sync/readiness and build passed.
  The first sandboxed attempt exited 2 because force-stop tests could not
  inspect the host process table. Its suite guard also detected this Agent's
  final fixture edit during execution. The rerun held source files unchanged.

The first focused Go invocation hit the sandbox's denied default build cache;
the task-scoped cache above allowed it to run. No follow-up scope was added.
The complete authored Verification is left for the Daemon.

## Carry-forward provenance

- Source Run: `run_20260930T144240Z_33c18e8a522f7217`
- Source commit: `563ca5aa82b33a35878b8c0f673b06a62615ff03`
