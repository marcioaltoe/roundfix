---
task: task_03
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: completed
type: backend
complexity: high
---

# Task 03: `roundfix deliver plan` shows what is approved to run, and start refuses what is not

## Overview

`runDeliverStart` in `internal/cli/deliver.go` only checks that each slug loads, so an operator learns that a Spec lacks `push`, `pull_request` or `merge` authority only when the finished item parks as `unauthorized`. Nothing lists the open intent that is not a Spec. This Task adds `roundfix deliver plan`, a read-only report built from the Specs Root, each Spec's committed authorization, the strict Spec Consistency Check and the repository's backlog, findings and inbox directories. It also makes `deliver start` refuse a Spec without delivery authority. The plan's stdout is read by the operator and by a Supervisor's automation. It MUST NOT open the Run Database or write any file.

## Requirements

1. MUST add `internal/cli/deliver_plan.go` with the helpers `strictSpecFindings(specsRoot, repoRoot, specSlug string) ([]speccheck.Finding, error)`, which composes `speccheck.Check`, `speccheck.PromoteGaps` and `speccheck.GatePrecondition(...).Findings`, and `productionPremises(graph *spec.Graph) []string`, which returns the sorted, unique `interface:` paths of non-`qa` Tasks that end in `.go` and not in `_test.go`. task_01's revalidation reuses both. The same file MUST implement `roundfix deliver plan [--json] [<slug>...]` as the TechSpec states, and dispatch `plan` from `runDeliverCommand`. With no slug it reports every active Spec from `spec.ListActiveDetailed`, and a slug `spec.Load` refuses exits `2`.
2. MUST compute for each Spec its Task count, its unfinished Task count and its reasons:
   - the authorization reasons from a helper `deliveryAuthorizationReasons(ctx, loaded, specsRoot, slug) []string`, which reads `spec.ReadSpecAuthorization` at the checkout's `HEAD` and yields `authorization <outcome>: <reason code>` or `authorization lacks <op>, <op>` over `implement`, `commit`, `push`, `pull_request` and `merge`;
   - `spec check: <code>, <code>` from `strictSpecFindings`.

   A Spec with no reason is `approved`; any other is `blocked`.
3. MUST report, for each Spec, the `productionPremises` it shares with each earlier Spec in the given order. A `shared` row is information only: it names the files for which the later item will carry a `premise-changed` warning once the earlier one merges, and it never makes a Spec `blocked` or refuses a start.
4. MUST list every `docs/backlog/*.md` and `docs/findings/*.md` with its front-matter `status`, and every regular file under `docs/_inbox/` with status `-`, each sorted by path. A missing directory contributes nothing.
5. MUST print the tab-separated `spec`, `shared`, `backlog`, `finding` and `inbox` rows the TechSpec shows, or with `--json` one `roundfix-deliver-plan/v1` document carrying the same facts. The plan exits `0` when every reported Spec is approved, `1` when any is blocked and `2` on a usage or preflight error.
6. MUST make `runDeliverStart` call `deliveryAuthorizationReasons` for every slug after `spec.Load` and before `store.Open`. When any slug has a reason, it fails through `printDeliverFailure` with exit `2`, names every such slug with its reasons and `roundfix deliver plan`, and records no queue. A strict finding MUST NOT refuse a start.
7. MUST add `roundfix deliver plan [--json] [<slug>...]` and a Commands line to `deliverUsage`, and change the top-level usage line in `internal/cli/cli.go` to `roundfix deliver <plan|start|status|resume|retry|stop> [<slug> ...]`.
8. MUST update only the tests this change invalidates:
   - add the plan line to the `deliver` row of `TestRunCommandHelp` in `internal/cli/cli_test.go`, adding, renaming or removing no top-level test there;
   - make `TestTopLevelUsageNamesDeliverRetry` in `internal/cli/deliver_retry_test.go` expect the new top-level line under the same name;
   - make `TestDeliverStatusPrintsTheItemWorktree` and `TestATerminalQueueIsReplacedByANewStart` in `internal/cli/deliver_test.go` grant all five operations with `setImplementFixtureAuthorizationOperations` before they start a queue.
9. MUST document in the deliver section of `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md`:
   - the command, by the string `roundfix deliver plan`;
   - its rows, JSON document, verdicts and exit codes;
   - that it writes nothing and is never implementation authority;
   - that a `shared` row predicts a warning, never a stop;
   - the start refusal.

   Then MUST regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
10. MUST put the new tests in `internal/cli/deliver_plan_test.go`. They drive the public CLI in disposable repositories, which may start from the `clean` fixture Spec that `newSpecCheckWorkspace` copies, and inject no production-only hook.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A plan over an approved Spec and a Spec whose authorization lacks `merge` exits `1` and reports the first `approved` and the second `blocked` with `authorization lacks merge`. A plan over approved Specs only exits `0`.
- [ ] `--json` carries the same verdicts, reasons, shared premises and intent as the text rows.
- [ ] Two Specs declaring the same production Go file yield a `shared` row for the later one, and both stay `approved`. A shared test or guide yields none.
- [ ] Backlog Entries, Findings and inbox notes appear as intent rows. The plan leaves the checkout's HEAD and status unchanged and creates no Run Database. With no slug, it reports every active Spec, and an unknown slug exits `2`.
- [ ] `deliver start` with a Spec lacking a delivery operation exits `2` and records no queue, while a Spec granting all five operations starts one.
- [ ] `roundfix deliver --help` and `roundfix --help` name the plan command.

## Context

- creates: `internal/cli/deliver_plan.go`
- creates: `internal/cli/deliver_plan_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/deliver_retry_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliverPlanReportsAnApprovedAndABlockedSpec|TestDeliverPlanExitsZeroWhenEverySpecIsApproved|TestDeliverPlanJSONCarriesTheSameFacts|TestDeliverPlanNamesSharedProductionPremises|TestDeliverPlanIgnoresSharedTestsAndGuides|TestDeliverPlanListsIntentThatIsNotApprovedToRun|TestDeliverPlanWritesNothing|TestDeliverPlanDefaultsToEveryActiveSpec|TestDeliverPlanRefusesAnUnknownSlug|TestDeliverStartRefusesASpecWithoutDeliveryAuthority|TestDeliverStartAcceptsASpecWithEveryDeliveryOperation|TestDeliverHelpNamesThePlanCommand|TestRunCommandHelp|TestTopLevelUsageNamesDeliverRetry|TestDeliverStatusPrintsTheItemWorktree|TestATerminalQueueIsReplacedByANewStart|TestDeliverCommandRejectsUnknownSlugBeforeRecordingQueue|TestDeliverCommandRefusesUnknownFlags)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverPlanReportsAnApprovedAndABlockedSpec TestDeliverPlanExitsZeroWhenEverySpecIsApproved TestDeliverPlanJSONCarriesTheSameFacts TestDeliverPlanNamesSharedProductionPremises TestDeliverPlanIgnoresSharedTestsAndGuides TestDeliverPlanListsIntentThatIsNotApprovedToRun TestDeliverPlanWritesNothing TestDeliverPlanDefaultsToEveryActiveSpec TestDeliverPlanRefusesAnUnknownSlug TestDeliverStartRefusesASpecWithoutDeliveryAuthority TestDeliverStartAcceptsASpecWithEveryDeliveryOperation TestDeliverHelpNamesThePlanCommand TestRunCommandHelp TestTopLevelUsageNamesDeliverRetry TestDeliverStatusPrintsTheItemWorktree TestATerminalQueueIsReplacedByANewStart TestDeliverCommandRejectsUnknownSlugBeforeRecordingQueue TestDeliverCommandRefusesUnknownFlags; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < internal/cli/cli_test.go | grep -qF -- "roundfix deliver plan [--json] [<slug>...]" && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "roundfix deliver plan" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "roundfix deliver plan" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the twelve new named tests exists and the command is documented nowhere, so the command fails.

## References

- [_techspec.md](_techspec.md) — The Delivery Plan and the start refusal
- `_prd.md` → Goals 1-2; Core Features 1-2; Success Metric 1
- `_techspec.md` → API Contracts 1-2; Testing Approach 3

## Result

Implemented the read-only Delivery Plan and the queue authorization preflight.
The plan now reads committed Spec authorization, strict Spec Consistency
findings, Task state, production premises, and repository intent; renders the
documented tab-separated or `roundfix-deliver-plan/v1` output; and never opens
the Run Database. `deliver start` now loads every Spec and refuses every
missing delivery operation before `store.Open`. Help, the user guide, and the
canonical/generated Roundfix skills carry the new contract.

Focused evidence by acceptance criterion:

- Approved/blocked verdicts and exit codes: the initial focused test was red
  because `plan` was unknown. After implementation,
  `TestDeliverPlanReportsAnApprovedAndABlockedSpec` and
  `TestDeliverPlanExitsZeroWhenEverySpecIsApproved` pass and assert
  `authorization lacks merge`, exit `1` for a mixed plan, and exit `0` for an
  approved-only plan.
- JSON parity: `TestDeliverPlanJSONCarriesTheSameFacts` passes and decodes the
  schema, verdicts, reasons, Task counts, shared premises, and intent.
- Shared premises: `TestDeliverPlanNamesSharedProductionPremises` and
  `TestDeliverPlanIgnoresSharedTestsAndGuides` pass; the shared production Go
  path produces an informational row while both Specs remain approved, and
  shared tests/guides produce no row.
- Intent and read-only behavior:
  `TestDeliverPlanListsIntentThatIsNotApprovedToRun`,
  `TestDeliverPlanWritesNothing`, `TestDeliverPlanDefaultsToEveryActiveSpec`,
  `TestDeliverPlanRefusesAnUnknownSlug`, and
  `TestDeliverPlanRefusesAnUnknownFlag` pass. They cover Backlog Entry,
  Finding, and inbox rows; unchanged HEAD/status; no Run Database; active-Spec
  discovery; and exit `2` preflight refusals.
- Start authorization: `TestDeliverStartRefusesASpecWithoutDeliveryAuthority`
  and `TestDeliverStartAcceptsASpecWithEveryDeliveryOperation` pass. The
  refusal names the Spec, missing operation, and `roundfix deliver plan`,
  starts no owner, and creates no Run Database; the full five-operation grant
  records the queue and starts its owner.
- Help: `TestDeliverHelpNamesThePlanCommand`, `TestRunCommandHelp`, and
  `TestTopLevelUsageNamesDeliverRetry` pass with the new command and top-level
  usage.

Focused checks:

- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestDeliver' ./internal/cli` — passed.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^(TestRunCommandHelp|TestTopLevelUsageNamesDeliverRetry|TestATerminalQueueIsReplacedByANewStart)$' ./internal/cli` — passed.
- `make skills-sync-check` — passed after `make skills-sync` regenerated the
  distributed Roundfix skill.
- `make verify-incremental` — the sandboxed run reached the suite and failed
  only because two force-stop integration tests could not read the host
  process table; the required-permission rerun passed, including `go vet`, all
  Go packages, skill checks, and the build.
- `git diff --check` — passed.

The Task's authored `## Verification` command was not run; Daemon Verification
remains the settlement authority.
