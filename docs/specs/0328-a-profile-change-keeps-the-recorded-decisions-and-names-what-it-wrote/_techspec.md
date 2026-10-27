---
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
prd: _prd.md
created: 2026-10-26
---

# A profile change keeps the recorded decisions and names what it wrote — Technical Spec

## Executive Summary

Three of the five frictions are output gaps over behavior the ADRs require,
and two are missing behavior. `baseline update` stays on the recorded profile
(ADR-0263); its refusal of `profile` becomes a typed error whose next action
names the plan command. `baseline plan` gains an opt-in `PlanRequest` field
that reads the Setup Manifest and carries the recorded decisions the selected
profile still selects; only the CLI's plan command sets it, so the update,
the interactive workflow and every package test keep their inputs. The plan
refusal joins the missing decisions, their catalog-conditional dependents and
the preservation mode into one message. The apply adds one warning per Root
Backup to its existing `warnings` array, so text and JSON name it with no
schema change. Reconcile's text, and only its text, explains a no-change
result with the Setup Snapshot revision. The restore plans a Claude Skill
Link per selected skill where `.claude/skills` already keeps them and writes
it as a symbolic-link postimage of the existing recoverable transaction.

The trade-off accepted: the transaction learns a second postimage kind, which
is more work than writing links after it commits, because Baseline apply
chose a journaled multi-file transaction over best-effort writes; and a plan's digest now depends on
manifest values a flag did not state, which is what carrying means.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; restore keeps its existing acquisition and every test
  uses an offline source. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0263, ADR-0264, ADR-0241 and
  ADR-0341 govern the design, and ADR-0221, ADR-0224 and ADR-0322 are
  recorded as not applying, as `_prd.md` → Project Constraints states and ADR
  Alternatives below relates. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — six Governed Paths, all test or skill
  files, granted in the authorization record; no production tooling,
  Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote/_authorization.md`;
  bounded files:
  `internal/cli/baseline_plan_test.go`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `6a2630e5`; `_prd.md` → Validation records the runs.

1. `DeclareManifestDecisions` (`internal/baseline/update.go`) refuses
   `preservation.mode` by name and passes any other id to
   `resolveManifestDecisionValues`, whose `normalizePlanDecisions`
   (`internal/baseline/plan.go`) answers `unknown decision "profile"`.
   `runBaselineUpdateCommand` wraps it with category `invalid`, exit `2` and
   the next action "correct the declared decisions and rerun roundfix
   baseline update".
2. `runBaselinePlanCommand` (`internal/cli/baseline_profile.go`) builds the
   `PlanRequest` from `--decision` and `--decision-file` only.
   `buildPlanWithCatalog` returns category `decision` for missing decisions
   before it checks `Preservation.Mode`, and `normalizePlanDecisions`
   computes the required set as the closure over answered values only, so a
   decision required by an unanswered decision's effect (`runtime.backend`
   and `runtime.design` through `autonomous.enabled` and the
   `autonomous-work` module) surfaces in a later round.
3. The interactive workflow already uses the recorded values as defaults
   (`inspectBaselineHumanState`, `promptBaselineDecisions` in
   `internal/cli/baseline_human.go`); the automation path does not.
4. `planRootBackups` (`internal/baseline/preservation.go`) names each backup
   `<base>.<sha256>.md`, and the plan carries a `backup:<carrier>` managed
   entry. `printBaselinePlanText` prints `- create <path> (N managed
   entries)`; `ApplyPlan` (`internal/baseline/apply.go`) copies the plan
   warnings into its `Result`; `writeBaselineApplyResult` prints only
   `Verified postimages: N`. A managed refresh plan carries no backup
   by design, enforced in `apply.go`.
5. `buildSkillsReconcilePlan` (`internal/baseline/skills_reconcile.go`) plans
   only lock-entry removals; `writeBaselineSkillsReconcilePayload`
   (`internal/cli/baseline_skills_reconcile.go`) prints `Baseline skills
   reconcile: no changes` when nothing is planned.
6. `buildSkillsRestorePlan` (`internal/baseline/skills_restore.go`) targets
   `.agents/skills/<name>` and `skills-lock.json`;
   `buildRestoreTransactionDocument` feeds the file transaction
   (`internal/baseline/transaction.go`), which stages, journals, verifies
   and rolls back regular-file postimages only and refuses a destination or
   parent that is a symbolic link (`validateMutationDestination`,
   `validatePathParents`). The update's skills stage restores through
   `baseline.RestoreSkills`.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/baseline/recorded_decisions.go` (new): reads the Setup Manifest
  through `ResolveManifestInput` (whose `Manifest` is available for a
  resolved or incompatible manifest) and computes the Carried Decisions for
  a resolved profile.
- `internal/baseline/plan.go`: `PlanRequest.CarryRecordedDecisions`,
  `PlanOutcome.CarriedDecisions`, and the one-refusal message.
- `internal/baseline/update.go`: the typed profile refusal.
- `internal/cli/baseline_profile.go`, `internal/cli/baseline_update.go`:
  wiring, next actions and the text lines.
- `internal/baseline/apply.go`: one warning per verified backup entry.
- `internal/baseline/skills_reconcile.go`, `internal/cli/baseline_skills_reconcile.go`:
  the Setup Snapshot revisions and the no-change lines.
- `internal/baseline/skills_restore.go`, `internal/baseline/transaction.go`:
  the Claude Skill Link plan and the symbolic-link postimage.

## Implementation Design

### Interfaces

```go
// internal/baseline/plan.go
type PlanRequest struct {
	// ...existing fields...
	CarryRecordedDecisions bool // only `baseline plan` sets it
}
type PlanOutcome struct {
	Plan             *PlanDocument
	Result           Result
	CarriedDecisions []string // sorted ids taken from the Setup Manifest; not serialized
}

// internal/baseline/recorded_decisions.go
func carryRecordedDecisions(root string, profile ResolvedProfile, explicit []DecisionValue,
	catalog *Catalog) (decisions []DecisionValue, carried []string, notCarried []Finding, err error)

// internal/baseline/update.go
type ProfileChangeRefusal struct{ Recorded, Requested string } // implements error

// internal/baseline/skills_reconcile.go
func SetupSnapshotRevisions(repoRoot, profileID, sourceRepository string) ([]string, error)
```

A Claude Skill Link travels through the transaction as a `Postimage` with
`Kind: PreimageSymlink`, `Mode: 0`, and `Content` holding the link target
bytes; its preimage is `PreimageMissing`.

### Data Models

No persisted schema changes. `roundfix/baseline-plan/v1`,
`roundfix/baseline-result/v1`, the update result and the reconcile payload
keep their fields; the restore payload's `plannedChanges` gains the action
value `create-link`, marshaled as `{"action","path","skill","target"}`. The
transaction journal records a symbolic-link postimage with the fields it
already has.

### API Contracts

1. API Contract: the update refuses a profile change. `roundfix baseline
   update --decision profile=<id>` (also with other flags) exits `2`,
   category `invalid`, writes nothing, with message `"profile" is not a
   Baseline decision; baseline update keeps the recorded Baseline Profile
   <recorded>` and next action `change the profile with roundfix baseline
   plan --repo <repo> --profile <id> --decision
   preservation.mode=preservation --format json, which carries the recorded
   decisions, then apply the reviewed plan with roundfix baseline apply
   --repo <repo> --plan <plan-file> --confirm-plan <planDigest>`, where
   `<repo>` is the `--repo` value as given. Any other unknown id keeps
   today's bytes.
2. API Contract: the plan carries. With a Setup Manifest present, `roundfix
   baseline plan` (with `--profile` or `--profile-file`) uses, for each
   decision the selected profile selects given the final answers and does
   not fix, the flag value when one is given, else the recorded value when
   the catalog accepts it. The text plan prints `Carried decisions: <ids>`
   (comma-separated, sorted) after `Baseline Profile:` when any was carried.
   Each recorded decision not carried yields one plan warning, code
   `baseline.decision.not-carried`, path `docs/agents/setup-context.json`,
   message `recorded decision <id> is not carried: <reason>`, where reason is
   `profile <p> does not select it`, `profile <p> fixes it`, or `the catalog
   refuses its recorded value`. A manifest that cannot be parsed contributes
   nothing and yields one warning, code `baseline.decision.manifest-unreadable`,
   on the same path. Without a manifest the plan is unchanged.
3. API Contract: one refusal. When answers are missing, the result keeps
   state `action_required`, category `decision`, exit `3`, and its message
   joins, with `; `, the existing `required Baseline decisions are missing:
   <ids>` (when any), one `required when <id> is <value>: <ids>` per
   unanswered decision whose catalog `equals` effect requires non-optional
   decisions not yet required, and `instruction-preservation mode is
   required` (when absent). When more than the preservation mode is missing,
   the next action is `answer every named decision with --decision or
   --decision-file, including --decision preservation.mode=greenfield or
   preservation, and rerun`; when only one part is missing, both texts keep
   today's bytes.
4. API Contract: Root Backups are named. The text plan prints a backup file
   change as `- create <path> (N managed entries; root backup of
   <carrier>)`. The apply result carries, after the plan's warnings, one
   warning per `backup:` managed entry, code
   `baseline.preservation.backup.written`, path the backup file, message
   `exact copy of <carrier> before this apply; commit it with this change,
   Roundfix never deletes it`; text prints it on the existing `Warning:`
   line and JSON in `warnings`.
5. API Contract: reconcile explains no change. When the payload is OK, not
   applied, with no planned change and no finding, the text prints after the
   status line: `Reconcile removes only skills-lock.json entries absent at
   <revision>; it never moves an installed skill to that revision.`, then
   `Installed skills from <source> follow the Setup Snapshot this binary
   carries: <refs>.` (or `No skill of profile <profile> comes from
   <source>.`), then `A newer revision reaches them with the Roundfix release
   whose Setup Snapshot names it; restore one with roundfix baseline skills
   restore --profile <profile> --skill <name>.` JSON is unchanged.
6. API Contract: Claude Skill Links. When `.claude/skills` is a real
   directory (neither it nor `.claude` a symbolic link) holding at least one
   entry that is a symbolic link whose target is `../../.agents/skills/<its
   own name>`, the restore plans `create-link .claude/skills/<name>` with
   target `../../.agents/skills/<name>` for each selected skill without a
   `.claude/skills/<name>` entry, includes it in the Plan Digest, writes it
   in the same transaction as the tree and lock, verifies it with `Lstat`
   and `Readlink`, and removes it on rollback. An existing entry of any kind
   is left untouched; otherwise nothing is planned. The text preview lists
   `- create-link .claude/skills/<name> [<name>]`.

### Surface Transcripts

None. No command gains a flag, an argument or an exit code; the new lines
are fixed in the API Contracts above and asserted, line for line, by the
named tests of Testing Approach.

### Invariants

```text
1. A flag answer always wins over a recorded decision, and a plan without a Setup Manifest is byte-identical to today's.
2. baseline update never changes the recorded profile, its exit codes or its JSON fields.
3. A managed refresh writes no Root Backup and adds no backup warning.
4. Reconcile plans exactly what it planned before; only its no-change text grows.
5. Restore never replaces, removes or follows an existing .claude/skills entry, and creates no link where the repository keeps none.
6. Every restore write, a link included, is staged, verified and rolled back by the one transaction.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `ProfileChangeRefusal`,
  `runBaselineUpdateCommand` (API Contract 1).
- Goal 2, Core Features 2–3, Success Metric 2 → `carryRecordedDecisions`,
  `buildPlanWithCatalog`, `runBaselinePlanCommand` (API Contracts 2–3).
- Goal 3, Core Feature 4, Success Metric 3 → `ApplyPlan`,
  `printBaselinePlanText` (API Contract 4).
- Goal 4, Core Feature 5, Success Metric 4 → `SetupSnapshotRevisions`,
  `writeBaselineSkillsReconcilePayload` (API Contract 5).
- Goal 5, Core Feature 6, Success Metric 5 → `buildSkillsRestorePlan`,
  `buildRestoreTransactionDocument`, the file transaction (API Contract 6).
- Success Metric 6 → declared unreachable; the guides and skills task_01
  writes tell the adopter the new path.

Deviations from the PRD: None.

## Absence Evidence

1. Carried Decisions in the plan command — command: `grep -n "ResolveManifestInput\|setup-context\|Manifest" internal/cli/baseline_profile.go`; result: one hit, the `capabilities check` next action; `runBaselinePlanCommand` never reads the manifest.
2. Root Backup warning — command: `grep -rn "backup.written\|root backup of" internal`; result: no match.
3. Reconcile snapshot explanation — command: `grep -n "Setup Snapshot\|snapshot" internal/cli/baseline_skills_reconcile.go`; result: no match.
4. Claude Skill Link planning — command: `grep -rn "create-link\|claude/skills" internal/baseline`; result: only the retired-skill listing in `retired_skills_installed.go` and its test.
5. Symbolic-link postimage — command: `grep -n "os.Symlink\|Symlink(" internal/baseline/transaction.go`; result: no match.
6. Profile refusal — command: `grep -n '"profile"' internal/baseline/update.go`; result: no match.

## ADR Alternatives

- ADR-0221: untouched; installed skills stay held to the Setup Snapshot, which the reconcile text states instead of moving them.
- ADR-0224: untouched; this repository's skills and lock are not restored by this Spec.
- ADR-0241: followed; the revisions reconcile names come from the same snapshot contracts a repository-owned profile takes.
- ADR-0263: followed; the update stays on the recorded profile, and the plan answers the alternative it rejected (restating every recorded decision) by carrying them.
- ADR-0264: followed; reconcile and restore never move a vendored skill to an unproven revision, so friction (d) changes text only.
- ADR-0322: untouched; the `python` profile and its preset are the target of the adopter's change, not changed by it.
- ADR-0341: implemented as recorded, including its rejection of a post-commit link write and of linking every `.claude/skills` directory.

## Integration Points

None external. The restore keeps its existing acquisition (`--source-dir` or
the declared GitHub commit); tests use an offline source.

## Testing Approach

All tests are hermetic: temporary repositories, `t.TempDir()` homes, the
offline restore source the existing restore tests build, and in-process
`RunContext` for the CLI. Test names below are the ones the Tasks' Verification
runs.

- task_02, `internal/baseline/recorded_decisions_test.go`:
  `TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects` (flag wins,
  fixed and unselected not carried, refused value not carried, closure through
  `autonomous.enabled`), `TestCarryRecordedDecisionsWithoutAManifestChangesNothing`.
  `internal/cli/baseline_profile_change_test.go`:
  `TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange` (adopt
  `go-cli-tui` with plan and apply, then plan `python` with only
  `preservation.mode`: exit `0`, decisions equal the recorded values, text
  lists `Carried decisions:`),
  `TestBaselinePlanNamesEveryMissingAnswerInOneRefusal` (no manifest, no
  decisions: one message naming every base decision, `required when
  autonomous.enabled is true: runtime.backend, runtime.design` and the
  preservation mode),
  `TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand` (text and
  JSON, exit `2`, tree unchanged).
- task_03, `internal/cli/baseline_root_backup_test.go`:
  `TestBaselineApplyNamesTheRootBackupItWrote` (text and JSON),
  `TestBaselinePlanTextLabelsTheRootBackup`,
  `TestBaselineUpdateWritesNoRootBackupWarning`;
  `internal/cli/baseline_skills_reconcile_test.go`:
  `TestBaselineSkillsReconcileNoChangesNamesTheSetupSnapshotRevision`, whose
  expected ref is read from the embedded catalog, never a literal.
- task_04, `internal/baseline/skills_restore_claude_link_test.go`:
  `TestSkillsRestoreLinksARestoredSkillWhereClaudeLinksAreKept`,
  `TestSkillsRestoreAddsNoLinkWithoutClaudeLinks` (absent, symbolic-link and
  copies-only `.claude/skills`), `TestSkillsRestoreLeavesAnExistingClaudeEntry`,
  `TestSkillsRestoreRollsBackTheClaudeLinkWithTheTree`,
  `TestSkillsRestoreLinksAnInstalledSkillThatLacksItsLink`;
  `internal/baseline/transaction_test.go`:
  `TestTransactionWritesVerifiesAndRollsBackASymbolicLinkPostimage`;
  `internal/cli/baseline_skills_restore_test.go`:
  `TestBaselineSkillsRestoreTextListsTheClaudeLink`.
- Existing suites that consume the changed contracts run in the same
  Verification commands: `./internal/baseline` and `./internal/cli` tests of
  plan, update, apply, reconcile and restore.

## Build Order

1. Carried Decisions, one refusal and the update refusal (task_02).
2. Root Backup naming and the reconcile text (task_03) (depends on: 1 — both
   edit `internal/cli/baseline_profile.go`).
3. Claude Skill Links in the restore transaction (task_04).
4. Guides, skills, glossary and Behavior Surface record (task_01) (depends
   on: 1, 2, 3).
5. QA gate (task_05) (depends on: 4).

## Risks & Considerations

- A plan's digest now depends on manifest values: an adopter who relied on a
  refusal to force restating decisions now gets a plan; the text lists the
  Carried Decisions so the review shows them.
- Symbolic links on Windows can need a privilege; the restore plans links
  only where the repository already holds one, so a Windows checkout of
  copies gets none, and a failed link write rolls the transaction back with
  its path named.
- `internal/baseline/transaction.go` is shared by every Baseline apply; the
  symlink kind must be unreachable from Baseline plans, which never build
  one (Invariant 6 and the existing apply tests guard it).

## Glossary

- adds: **Carried Decision**
- adds: **Root Backup**
- adds: **Claude Skill Link**

## Decisions

- Typed refusal over supporting a profile change in the update. See ADR-0341.
- Opt-in `PlanRequest` field over changing `BuildPlan` for every caller, so
  the update and the interactive workflow keep their inputs.
- Warnings in the existing array over a new result field, so no schema moves.
- Reconcile text only; JSON consumers already read `skills[].disposition`.
- Links inside the transaction, only where links already exist. See ADR-0341.
