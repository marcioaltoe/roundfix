---
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
prd: _prd.md
created: 2026-10-27
---

# The skills lock matches what is installed, and a Run without a QA gate settles green — Technical Spec

## Executive Summary

Two defects, two components. The owned-skill bundle install writes skill
files and never reads `skills-lock.json`, so an adopter's lock entry for an
owned skill keeps the hash of the copy it replaced. A new Baseline function
rewrites the `computedHash` of each existing lock entry that names an
installed owned skill, through the ordered-JSON lock document and the
external skills digest the restore already uses, and both callers of the
install call it. The Daemon appends the repository Verification at
settlement only when the graph has a QA gate Task; the condition loses that
clause, so every non-QA Task settles on it, and the in-process checks keep
their gate scope. Because Settlement Checks see one Task's tree, `TaskCycle`
of a graph without a QA gate runs the repository Verification once more on
the integrated Run Worktree when parallel Tasks could have integrated
unverified, and `implement` ends such a Run `Unresolved` when it fails.

The trade-off accepted: a graph without a QA gate pays one repository
Verification per Task plus, under Task Capacity above one, one at the end,
the same price gated graphs pay since ADR-0182; and the lock keeps a `ref`
that names the upstream copy first installed, because the format offers no
provider that names a binary bundle (ADR-0189 addendum).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; every test uses temporary repositories and fake
  verifiers. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0182, ADR-0342, ADR-0189,
  ADR-0233 and ADR-0256 govern the design, and ADR-0197, ADR-0250, ADR-0272,
  ADR-0279 and ADR-0286 are recorded as not applying, as `_prd.md` → Project Constraints states and
  ADR Alternatives below relates. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — seven Governed Paths, all owned-skill
  files task_01 edits, granted in the authorization record; no production
  tooling, Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green/_authorization.md`;
  bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Read on 2026-10-27 at `5e9475bf`; `_prd.md` → Validation records the runs.

1. `Install` (`skills/skills.go`) writes every bundled file under the target
   directory with `os.WriteFile` and returns; it never opens
   `skills-lock.json`. Its callers are `runBaselineUpdateSkillsStageWith`
   (`internal/cli/baseline_update.go`, through
   `baselineUpdateSkillsDependencies.install`) and `runSkillsInstall`
   (`internal/cli/cli.go`).
2. The lock is read and written with order preserved by `loadSkillsLock`,
   `skillsLockDocument` and `marshalIndent`
   (`internal/baseline/skills_restore_lock.go`); `loadSkillsLock` refuses a
   symbolic link, a non-regular file and a malformed document, and returns
   an empty document for an absent file. The restore computes an entry's
   hash with `externalSkillsLockDigest` over the files
   `inspectRestoreTarget` reads (`internal/baseline/skills_restore.go`);
   `skills.skillFolderHash` computes the same digest for readiness.
3. Owned skills never enter the external set
   (`resolveExternalSkillRequirement` in `internal/cli/doctor.go`), and the
   Setup Snapshots mark them `type: repo`, so no Roundfix reader compares an
   owned skill with its lock entry; an adopter's own test does.
4. `TaskCycle` (`internal/daemon/task_engine.go`) sets
   `taskPlan.settlementChecks = qaTask != nil`. `executeTaskWorker` sets
   `requiredRepositoryVerification` when `enteredOnRedRepository ||
   (plan.settlementChecks && plan.RepositoryVerificationAtSettlement)`, and
   `executeTask` appends it as the attempt's last command. The Spec
   Consistency baseline (in `executeTaskWorker`) and the in-process checks
   (in `verifyTask`) run only when `settlementChecks` is true.
5. After `runTaskScheduler`, `TaskCycle` runs the QA gate when there is one
   and returns. `runImplementCommand` (`internal/cli/implement.go`) chooses Clean
   when no Task failed or was skipped and, for a gated graph, the QA verdict
   was accepted.
6. The built-in repository Verification is `make verify`
   (`defaultVerification`, `internal/config/config.go`). In a scratch
   worktree, removing `plan.settlementChecks &&` from the condition of item 4
   failed one test of `./internal/daemon`,
   `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification`, which
   pins the old scope, and 21 tests of `./internal/cli`, measured in two full
   runs with the same list: 18 in `implement_test.go` and 2 in
   `implement_unix_test.go` fail because their gate-less fixture never sets
   a repository Verification and `make verify` fails in a temporary
   repository without a Makefile, and
   `TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit` hangs
   until the test timeout, because its `onTask` commits user work on every
   Agent turn and the repair turn finds nothing to commit. task_03 names all
   21 and migrates them.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/baseline/owned_skills_lock.go` (new):
  `RecordInstalledOwnedSkills`.
- `internal/cli/baseline_update.go`: one dependency, the call after the
  install, `skills.lockUpdated` and its text lines.
- `internal/cli/cli.go`: `runSkillsInstall` calls the function for the
  project target and prints one line per rewritten entry.
- `internal/daemon/task_engine.go`: the settlement condition, the Integrated
  Repository Verification in `TaskCycle`, and one `TaskCycleResult` field.
- `internal/cli/implement.go`: the outcome and the stderr line.

## Implementation Design

### Interfaces

```go
// internal/baseline/owned_skills_lock.go
// RecordInstalledOwnedSkills rewrites computedHash of each existing
// skills-lock.json entry named in names to the digest of
// <root>/.agents/skills/<name>, and returns the rewritten names, sorted.
func RecordInstalledOwnedSkills(ctx context.Context, root string, names []string) ([]string, error)

// internal/cli/baseline_update.go
type baselineUpdateSkillsResult struct {
	// ...existing fields...
	LockUpdated []string `json:"lockUpdated,omitempty"`
}
type baselineUpdateSkillsDependencies struct {
	// ...existing fields...
	recordLock func(context.Context, string, []string) ([]string, error)
}

// internal/daemon/task_engine.go
type TaskCycleResult struct {
	// ...existing fields...
	IntegratedVerificationFailure string // empty unless the check ran and failed
}
```

### Data Models

No persisted schema changes. `skills-lock.json` keeps its document; only
the `computedHash` string of an owned-skill entry can change. The
`roundfix/baseline-update-result/v1` payload gains the optional
`skills.lockUpdated` array, omitted when empty. The Run Database is
unchanged; the integrated check's diagnostics file is
`<artifact-root>/runs/<run-id>/verification/integrated-repository.log`.

### API Contracts

1. API Contract: the lock follows the bundle.
   `RecordInstalledOwnedSkills(ctx, root, names)`:
   - returns `nil, nil` when `<root>/skills-lock.json` does not exist, and
     never creates it;
   - loads the lock through `loadSkillsLock`, so a symbolic link, a
     non-regular file or a malformed document returns an error that names
     `skills-lock.json` and writes nothing;
   - for each name, in sorted order, whose `skills.<name>` value is an
     object and whose `.agents/skills/<name>` folder exists, computes the
     digest with `externalSkillsLockDigest` over the files
     `inspectRestoreTarget` reads, and sets that entry's `computedHash`
     field in place when it differs, keeping every other field, its value
     and the key order; a name without an entry or without a folder is
     skipped;
   - writes the serialized document only when an entry changed, through a
     temporary file in the repository root renamed over the lock, keeping
     the lock's file mode; when nothing changed the lock's bytes and
     modification time are untouched;
   - returns the rewritten names, sorted.
2. API Contract: `baseline update` records the lock. In
   `runBaselineUpdateSkillsStageWith`, right after a successful install,
   `recordLock(ctx, root, ownedNames())` runs; its names become
   `result.LockUpdated`. An error sets the stage `failed` and returns
   `record installed Roundfix skills in skills-lock.json: <error>`, which
   the update reports as today's skills-stage failure. Text prints, after
   the `- installed <name>` lines and only when the list is not empty,
   `Skills lock updated: <count>` and one `- lock updated <name>` line per
   name. The preview without `--yes` or `--confirm-plan` and the Plan Digest
   are unchanged. `runBaselineUpdateSkillsStage` wires
   `baseline.RecordInstalledOwnedSkills`.
3. API Contract: `skills install` records the lock. `runSkillsInstall`
   calls `baseline.RecordInstalledOwnedSkills(ctx, projectRoot,
   roundskills.Names())` only for `--target project` without `--dir`,
   after the `Installed Roundfix skill for project: …` line, and prints
   `Updated skills-lock.json computedHash: <name>` per rewritten name on
   stdout. An error prints `roundfix: skills install failed: <error>` on
   stderr and exits `2`, as an install failure does today. Other targets
   and `--dir` never read a lock.
4. API Contract: every Task settles on the repository Verification.
   `executeTaskWorker` sets `requiredRepositoryVerification` to the trimmed
   `plan.RepositoryVerification` when `enteredOnRedRepository ||
   plan.RepositoryVerificationAtSettlement`. An empty command appends
   nothing. `settlementChecks` keeps `qaTask != nil` and still gates the
   Spec Consistency baseline and the two in-process checks. Order, repair
   turn, final-failure reason and the precondition-repair rule are
   unchanged.
5. API Contract: the Integrated Repository Verification. After
   `runTaskScheduler` returns without error, `TaskCycle` runs the check when
   all of these hold: the plan has no QA gate Task,
   `plan.RepositoryVerificationAtSettlement` is true, the trimmed
   `plan.RepositoryVerification` is not empty, `plan.Concurrency > 1`,
   `result.Completed >= 2`, and `result.Failed == 0 && result.Skipped == 0`.
   It runs the command once, in `plan.WorkDir` (the Run Worktree), through
   the Verifier under the Verification Capacity gate, with the diagnostics
   file of Data Models, a single exclusive retry after a temporary failure
   (exit `75`) as for any Verification command, and a `daemon_verification`
   event whose summary is `Integrated Repository Verification passed:
   <command>.` or `Integrated Repository Verification failed: <command>.`.
   A failure sets `result.IntegratedVerificationFailure` to
   `Integrated Repository Verification failed: <command>; diagnostics:
   <path>`, changes no Task status and sends no Agent turn. A Stop Request
   or an exhausted Run Budget during the check ends the cycle as either
   does today.
6. API Contract: the Run outcome. In `runImplementCommand`, a non-empty
   `IntegratedVerificationFailure` makes the outcome `Unresolved` before
   integration with the user's branch; the reason is printed on stderr
   before `Implement Run <id> reached Unresolved.` and carried as the
   terminal completion reason; the Run Worktree is kept and the exit code
   is `1`, as for any `Unresolved` Run.

### Surface Transcripts

None. No command gains a flag, an argument or an exit code; the new lines
are fixed in the API Contracts above and asserted by the named tests of
Testing Approach.

### Invariants

```text
1. A lock entry that names no installed owned skill, and every field other than computedHash, keep their bytes; a lock without such an entry is never written.
2. A repository without skills-lock.json never gets one from the bundle install.
3. With verification.repository_at_settlement false, no settlement command and no Integrated Repository Verification run in any graph.
4. A gated graph runs exactly the settlement commands, in-process checks and QA precondition it ran before.
5. The Spec Consistency and authorization Settlement Checks never run in a graph without a QA gate Task.
6. The Integrated Repository Verification never changes a Task status, never sends an Agent turn, and never runs when a Task failed or was skipped.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1–3 → `RecordInstalledOwnedSkills`,
  `runBaselineUpdateSkillsStageWith`, `runSkillsInstall` (API Contracts 1–3).
- Goal 2, Core Feature 2, Success Metrics 4 and 6 → `executeTaskWorker`
  (API Contract 4).
- Goal 3, Core Feature 3, Success Metric 5 → `TaskCycle`, `runImplementCommand`
  (API Contracts 5–6).
- Success Metric 7 → declared unreachable; the guides and skills task_01
  writes tell the adopter the new behavior.

Deviations from the PRD: None.

## Absence Evidence

1. `RecordInstalledOwnedSkills` — command: `grep -rn "skills-lock\|computedHash" skills/skills.go`; result: no match; `Install` never reads the lock.
2. `skills.lockUpdated` and its text — command: `grep -rn "lockUpdated\|Skills lock updated" internal .agents docs/user-guide`; result: no match.
3. Integrated Repository Verification — command: `grep -rn "integrated-repository\|IntegratedVerificationFailure" internal`; result: no match.
4. Gate-less settlement command — command: `grep -n "settlementChecks &&" internal/daemon/task_engine.go`; result: one match, the condition API Contract 4 changes.

## ADR Alternatives

- ADR-0182: narrowed by ADR-0342 for the repository Verification only; its choice of keeping these facts to the gate is answered by the Baseline clause that the Daemon runs the repository Verification at settlement, and its in-process checks keep their scope.
- ADR-0342: implemented as recorded, including its rejection of refusing Clean without a check and of an end-only check.
- ADR-0189: implemented through its 2026-10-27 addendum; the bundle stays the owned skill's source of truth, and the lock now describes it.
- ADR-0233: followed; task_01 raises the owned-skill versions with the record command and never by hand.
- ADR-0256: followed; the skills that describe `implement`, `baseline update` and `skills install` change in the same Spec.
- ADR-0250: untouched; no Baseline module, module version or coverage record changes.
- ADR-0272: untouched; no review record or review correction path changes.
- ADR-0197: untouched; no pre-PR reviewer lineage changes.
- ADR-0279: untouched; no Queue Owner status or review dismissal changes.
- ADR-0286: untouched; no Correction Review changes.

## Integration Points

None external. The lock format is the one the upstream skills CLI writes
(vercel-labs/skills `src/local-lock.ts`); no command calls that tool.

## Testing Approach

All tests are hermetic: temporary repositories and homes, the fake runners,
verifiers and committers of the existing daemon and implement suites, and
in-process `RunContext` for the CLI. The test names below are the ones the
Tasks' Verification runs.

- task_02, `internal/baseline/owned_skills_lock_test.go`:
  `TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill`
  (a stale owned entry is rewritten to the folder digest; `source`, `ref`,
  `sourceType`, `skillPath`, an extra field, the key order and an external
  entry keep their bytes),
  `TestRecordInstalledOwnedSkillsLeavesAMatchingLockUntouched` (second call
  returns no name and the file's bytes and modification time are equal),
  `TestRecordInstalledOwnedSkillsAddsNoEntryAndNoLock` (no lock: none
  created; a lock without the owned name: unchanged),
  `TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock` (symbolic link and
  malformed document: error naming `skills-lock.json`, bytes unchanged).
  `internal/cli/owned_skills_lock_test.go`:
  `TestSkillsInstallRewritesTheOwnedSkillLockHash` (real install, the
  `Updated skills-lock.json computedHash: qa-gate` line, the hash equal to
  the installed folder digest, a second run prints no such line),
  `TestBaselineUpdateSkillsStageReportsTheRewrittenLockEntries` (real
  install through `runBaselineUpdateSkillsStageWith`, `LockUpdated` and the
  text lines, JSON `skills.lockUpdated`, omitted when empty),
  `TestBaselineUpdateSkillsStageFailsOnAnUnsafeLock`.
- task_03, `internal/daemon/settlement_repository_verification_test.go`:
  `TestUngatedGraphRunsRepositoryVerificationAtSettlement` replaces
  `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification`
  (declared break: commands become `declared check`, `repository check`),
  `TestUngatedGraphRepositoryFailureAtSettlementReturnsFeedback`,
  `TestUngatedGraphRepositoryAtSettlementOffRunsOnlyDeclaredCommands`,
  `TestUngatedGraphRunsNoInProcessSettlementCheck` (a recording
  `SettlementChecker` is never called). `internal/cli/implement_settlement_checks_test.go`:
  `TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement`.
  The existing gated tests of both files run unchanged. The 21 implement
  tests of Current behavior 6 state their repository Verification
  explicitly, a passing command or the switch off where they count
  Verification calls for another purpose, and keep every other assertion.
- task_04, `internal/daemon/integrated_repository_verification_test.go`:
  `TestIntegratedRepositoryVerificationRunsOnceAfterAParallelUngatedGraph`
  (Task Capacity 2, two completed Tasks: one extra command in the Run
  Worktree, its event and diagnostics path),
  `TestIntegratedRepositoryVerificationFailureLeavesTaskStatuses`,
  `TestIntegratedRepositoryVerificationDoesNotRun` with subtests
  `capacity-one`, `one-completed`, `failed-task`, `switch-off`,
  `empty-command` and `gated-graph`.
  `internal/cli/implement_settlement_checks_test.go`:
  `TestImplementUngatedRunEndsUnresolvedWhenTheIntegratedTreeIsRed`
  (exit `1`, the stderr line, the Run record `Unresolved`, the Run
  Worktree kept) and `TestImplementUngatedRunEndsCleanWhenTheIntegratedTreeIsGreen`.
  The 21 migrated implement tests keep passing; a count that now includes
  the integrated call is the declared contract.
- Existing suites that consume the changed contracts run in the same
  Verification commands: `./internal/daemon` settlement tests,
  `./internal/cli` implement and baseline update tests, `./skills`.

## Build Order

1. Guides, skills, glossary and Behavior Surface record (task_01).
2. The lock follows the bundle (task_02) (depends on: 1).
3. Every Task settles on the repository Verification (task_03) (depends on:
   1).
4. The Integrated Repository Verification and the Run outcome (task_04)
   (depends on: 3 — both edit `internal/daemon/task_engine.go` and
   `internal/cli/implement_settlement_checks_test.go`).
5. QA gate (task_05) (depends on: 2, 4).

## Risks & Considerations

- A repository whose Verification is slow pays it once more per Task of a
  graph without a QA gate; `verification.repository_at_settlement: false`
  remains the operator's switch, documented beside the change.
- A repository already red on entry now fails every Task of a gate-less
  graph, as gated graphs do; the precondition-repair rule is the documented
  path.
- An adopter that never configured `defaults.verification` gets the
  built-in `make verify` at settlement in a gate-less graph, as in a gated
  one; a repository without that target must configure its command or turn
  the switch off, which the configuration guide says.
- The lock rewrite reserializes the document with two-space indentation and
  a trailing newline, the layout the upstream tool and the restore already
  write; a lock written in another layout is reformatted the first time an
  entry changes, never when none does (Invariant 1).
- The Integrated Repository Verification runs after the last Task
  integrated, so a failure has no Agent Session to repair it; the Run stays
  recoverable through its kept Run Worktree, a corrective Task and Task
  Carry-Forward.

## Glossary

- adds: **Integrated Repository Verification**
- changes: **Settlement Check**
- changes: **Clean**

## Decisions

- Rewrite only `computedHash` and keep the upstream-recorded source. See the
  ADR-0189 addendum.
- One Baseline function shared by both install callers, not a lock writer in
  `skills`, because the ordered lock document and the digest live in
  `internal/baseline` and `skills` cannot import it.
- Repository Verification at settlement in every graph, plus one integrated
  check for a parallel gate-less Run. See ADR-0342.
- The integrated check reports through the existing Run outcome, not a new
  state, because `Unresolved` already means "work remains" and keeps the
  Run Worktree.
