---
task: task_03
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: pending
type: backend
complexity: high
---

# Task 03: roundfix follow-up lists every pending item in Follow-up Priority order beside the work in progress

## Overview

This Task adds the Follow-up Command's list view: the two `secondbrain` User
Config keys, the inventory and in-progress readers, the call to the judge
task_02 added, the order of ADR-0319, the stop rule and the text and JSON
reports, wired into the root usage and the Skill Coverage Map (`_prd.md` →
Core Features 1-5 and 7; `_techspec.md` → API Contracts 1-7). The status
view and `--status` come in task_04.

## Requirements

1. MUST add the `Secondbrain` section to `internal/config` with the keys,
   the User-Config-only rule for `secondbrain.root`, the `~` expansion, the
   validation messages and the namespace rule of `_techspec.md` → API
   Contract 2, following the `delivery.check_timeout` and
   `jev.monthly_ceiling_usd` pattern (`removeYAMLPath` plus
   `warnIgnoredProjectSetting` for a Project Config `secondbrain.root`), and
   add `internal/config/secondbrain_test.go` with
   `TestSecondbrainRootIsReadFromUserConfig` (an absolute value and a `~/`
   value expanded against the test home),
   `TestProjectConfigCannotSetTheSecondbrainRoot` (the exact warning and an
   unchanged empty root), `TestSecondbrainKeysRefuseInvalidValues` (a
   relative root, `_triaged`, `Fluxus` and `a/b` as namespaces, each with its
   exact message) and `TestTheInboxNamespaceMayComeFromProjectConfig`
   (Project over User).
2. MUST add `roundfix follow-up [--messages <path>] [--format <text|json>]`:
   the root usage line `  roundfix follow-up [--messages <path>] [--format <text|json>]`
   right after the `deliver` line, the Commands row
   `  follow-up  List pending work in priority order, or show status` right
   after the `deliver` row, the dispatch case, the `commandUsage` case, and a
   help text equal to `_techspec.md` → Surface Transcript 1 without the
   `[--status] ` usage token, without the `--status` option line, and with
   the first prose line "List this repository's pending work in Follow-up
   Priority order.". Flags are parsed before any read, so a usage error
   reads nothing (API Contract 1).
3. MUST implement the inventory of `_techspec.md` → API Contract 3 in a new
   `internal/cli/follow_up_inventory.go`, reusing `deliveryPlanIntent` for
   `docs/backlog/`, `docs/findings/` and `docs/_inbox/`,
   `spec.ListActiveDetailed` with `roundconfig.ResolveSpecsRoot`, and
   `parseReadinessRemote` with `config.RepositoryRoot` for the default
   namespace; and the in-progress readers of API Contract 4 reusing
   `openRunsListReader` and `ListRuns`, `DeliveryQueue` with
   `delivery.ClassifyPark`, `preflight.InspectGit`, and a new
   `commandDependencies` field `followUpGHRunner` (default
   `preflight.ExecGHRunner{}`) for `gh pr list`. A missing or unavailable
   source never changes the exit code; an unreadable Run Database exits `1`.
   The inbox reader reads only regular `*.md` files directly in the namespace
   directory, never `_triaged/` and never a symbolic link.
4. MUST call `judge.PrioritizeFollowUp` once per run with the items in
   inventory order, `Project` the namespace, the Judge Key Variables copied
   from the command environment only outside a Run (as `spec judge` does with
   `commandRunID`), the ceiling from `jev.monthly_ceiling_usd`, and the
   `judgeTransport` and `judgeNow` dependencies; and order the items with
   `followUpOrder` exactly as API Contract 5 states, including the type
   fallback and the tie-breaks.
5. MUST print the JSON report of API Contract 6 and the text report of API
   Contract 7, with the stop rule of Core Feature 5, the judge line and the
   ten sources in their fixed order, and never print an item's text, only its
   title.
6. MUST add to `docs/references/skill-coverage.json` the entry
   `command follow-up` whose `skills` is
   `.agents/skills/roundfix/references/follow-up.md` and whose `sources` are
   `docs/user-guide/commands/follow-up.md`, `internal/cli/follow_up.go` and
   `internal/cli/follow_up_inventory.go` (never `internal/cli/cli.go`), and
   run `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
   last, because the new command and the new config keys add or change
   Behavior Surfaces. The record is never hand-edited.
7. MUST add `internal/cli/follow_up_test.go`, driving the command through
   `runCLI` with `setCommandEnvironmentForTest`, a temporary repository and
   home (as `newImplementWorkspace` builds them), a temporary Secondbrain
   Root named in a temporary User Config, the Run Database seeded through
   `store.Open` and `CreateDeliveryQueue`, a fake `followUpGHRunner`, and the
   judge fake of `specJudgeFixture` adapted to the four questions. Every test
   calls `t.Parallel()` first, uses no real-time constant below ten seconds,
   never sets a process environment variable, never runs the real `gh`,
   never reaches the network and never starts the CLI test binary as
   `roundfix`. The tests:
   - `TestFollowUpListsEveryPendingItemOnce` (Success Metric 1): two pending
     Inbox Entries whose `created_at` order is the reverse of their file-name
     order and one entry under `_triaged/`, one open and one promoted Backlog
     Entry, one pending Finding, one active Spec and a messages file with one
     message list exactly six items, the inbox entries in `created_at` order,
     with their kind, reference, type, date, title and `requested_by`;
   - `TestFollowUpOrdersItemsByTheFourAnswers` (Success Metric 2): a table of
     fake answers puts one item in each of the four classes, an item whose
     blocks answer is `0.5` (unsure) and whose fix answer is `0.9` lands in
     `health`, not `unblocking`, two items of one class order small first and
     then oldest, and `blocks` names `requested_by` only for a yes;
   - `TestFollowUpOrdersByTypeWithoutTheJudge` (Success Metric 3): without a
     key, with `ROUNDFIX_RUN_ID` in the command environment and at the
     ceiling, the fake transport receives no request, the order follows the
     type fallback, the judge line names the reason and the exit code is 0;
   - `TestFollowUpSaysNothingToDo` (Success Metric 4): an empty repository
     prints exactly `Nothing to do: no pending item and nothing in
     progress.` after the first line, with the ten sources, and JSON `state`
     `nothing_to_do`; one seeded active Run makes it `work_found`;
   - `TestFollowUpReportsAMissingSecondbrainRoot` (Success Metric 5): no
     `secondbrain.root` makes `secondbrain-inbox` `missing` with its detail
     and still lists the repository items;
   - `TestFollowUpNamespaceDefaultsToTheRepositoryName`: an `origin` of
     `git@github.com:owner/widget.git` gives `widget`; no `origin` gives the
     checkout's directory name; a configured namespace wins;
   - `TestFollowUpReadsNoUnavailableSourceAsAFailure`: a failing fake `gh`
     makes `pull-requests` `unavailable` with its first error line and the
     command still exits 0;
   - `TestFollowUpWritesNothing` (Success Metric 7): the bytes of every file
     under the repository, the Secondbrain Root and the Roundfix Home, and
     `git status --porcelain`, are identical before and after a text and a
     JSON run;
   - `TestFollowUpRefusesAnUnknownFormat`: asserts `_techspec.md` → Surface
     Transcript 2 whole — empty standard output, both standard error lines in
     order, exit 2;
   - `TestFollowUpRefusesAnInvalidMessagesFile`: a missing file, a file over
     1 MiB, a wrong schema, a missing `id` and a duplicate `id` each exit 2
     with a message naming the file.
8. MUST NOT add `--status`, change `deliver`, `runs`, `release` or `spec
   judge`, or change `internal/store`, `internal/delivery` or
   `internal/releaseplan`.

## Subtasks

- [ ] Add the `secondbrain` config keys and their tests.
- [ ] Add the command, its usage and help, and its dispatch.
- [ ] Add the inventory and in-progress readers.
- [ ] Add the judge call, the order, the stop rule and both renderers.
- [ ] Add the map entry, the tests, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Every pending item of every source is listed once with its fields, and no triaged, promoted or terminal record is.
- [ ] The order matches the four answers and falls back to the type with a named reason, and an unsure answer never makes an item `unblocking`.
- [ ] Nothing to do is reported only with no item and nothing in progress, and no run writes a byte.
- [ ] `secondbrain.root` takes effect only from User Config.

## Consumers

- the `roundfix follow-up` command surface and the `config keys` Behavior Surface — consumers: `./internal/docscontract`
- the judge contract task_02 added — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- instruction: `internal/cli/deliver_plan.go`
- instruction: `internal/cli/runs.go`
- instruction: `internal/cli/spec_judge.go`
- instruction: `internal/cli/spec_judge_test.go`
- instruction: `internal/cli/readiness_forge.go`
- instruction: `internal/cli/deliver_limits_test.go`
- instruction: `internal/config/delivery_check_timeout_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/secondbrain_test.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/follow_up.go`
- creates: `internal/cli/follow_up_inventory.go`
- creates: `internal/cli/follow_up_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSecondbrainRootIsReadFromUserConfig|TestProjectConfigCannotSetTheSecondbrainRoot|TestSecondbrainKeysRefuseInvalidValues|TestTheInboxNamespaceMayComeFromProjectConfig)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSecondbrainRootIsReadFromUserConfig TestProjectConfigCannotSetTheSecondbrainRoot TestSecondbrainKeysRefuseInvalidValues TestTheInboxNamespaceMayComeFromProjectConfig; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/config` — expected: exit 0. Before this Task the four config tests do not exist, so the command fails at its first pass check.
- `out="$(go test -count=1 -v -run '^(TestFollowUpListsEveryPendingItemOnce|TestFollowUpOrdersItemsByTheFourAnswers|TestFollowUpOrdersByTypeWithoutTheJudge|TestFollowUpSaysNothingToDo|TestFollowUpReportsAMissingSecondbrainRoot|TestFollowUpNamespaceDefaultsToTheRepositoryName|TestFollowUpReadsNoUnavailableSourceAsAFailure|TestFollowUpWritesNothing|TestFollowUpRefusesAnUnknownFormat|TestFollowUpRefusesAnInvalidMessagesFile|TestEveryCommandPathPrintsItsHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpListsEveryPendingItemOnce TestFollowUpOrdersItemsByTheFourAnswers TestFollowUpOrdersByTypeWithoutTheJudge TestFollowUpSaysNothingToDo TestFollowUpReportsAMissingSecondbrainRoot TestFollowUpNamespaceDefaultsToTheRepositoryName TestFollowUpReadsNoUnavailableSourceAsAFailure TestFollowUpWritesNothing TestFollowUpRefusesAnUnknownFormat TestFollowUpRefusesAnInvalidMessagesFile TestEveryCommandPathPrintsItsHelp; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide)$' ./internal/docscontract || exit 1; grep -qF '"command follow-up"' docs/references/skill-coverage.json || { echo 'missing command map entry' >&2; exit 1; }; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the follow-up tests do not exist, so the command fails at its first pass check. Afterwards the list view, its order, its stop rule, its refusals and its read-only boundary pass with every command's help, the command and its keys are covered and named in the skill and the guide, the map holds the command entry, and the package's tests compile for Windows.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5 and 7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 7
- `_techspec.md` → Interfaces; API Contracts 1-7; Surface Transcript 2; Invariants 1-7; Testing Approach; Build Order 3
- ADR-0095
- ADR-0209
- ADR-0319
