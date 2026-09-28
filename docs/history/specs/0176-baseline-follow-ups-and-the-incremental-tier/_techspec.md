---
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Baseline follow-ups and the incremental tier

## Executive Summary

This Spec closes four Baseline gaps:

- a blocked reconcile prints an empty plan and documents its exit;
- skill-lock reconciliation and restoration plan from the lock present after the
  fetch;
- the docs-layout clause routes a QA Archive Override through its command;
- the incremental tier becomes the required `verification.incremental` Baseline
  decision, published through the agent instructions and the Setup Manifest.

## Project Constraints

- Identifier strategy: not applicable — no new persisted identifier; the new
  decision uses the catalog identity `verification.incremental`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the embedded
  catalog only; no credential and no new network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0058, ADR-0066, ADR-0071, ADR-0072,
  ADR-0073, ADR-0080, ADR-0081, ADR-0085, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0099, ADR-0103, ADR-0104, ADR-0117, ADR-0118, ADR-0130,
  ADR-0144, ADR-0149, ADR-0154, ADR-0155 and ADR-0156 hold; ADR-0150 does not
  apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/rust-cli.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`,
  `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`,
  `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go`,
  `internal/cli/baseline_plan_test.go`,
  `internal/cli/baseline_release_gate_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The blocked reconcile result

`buildSkillsReconcilePlan` in `internal/baseline/skills_reconcile.go` builds the
required-removed payload with `PlannedChanges: []RestorePlannedChange{}`, so
the JSON field is `[]`. `PlanDigest` stays `nil` (JSON `null`), every
`LockEdit` stays `nil`, and `reconcileSkillsLock` still returns the
`reconcile.required-removed` error in category `action_required`.
`baselineSkillsRestoreExit` in `internal/cli/baseline_skills_restore.go` maps
that category to exit `3`, unchanged.

The `baseline skills reconcile` help in `internal/cli/cli.go` documents exit
`3` for both cases:

```text
  3  confirmation is required or does not match the current Change Plan, or a
     Profile-required skill is absent at the selected revision (finding
     reconcile.required-removed, no Plan Digest)
```

`docs/user-guide/commands.md` (the `baseline skills reconcile` paragraph) and
the Roundfix skill state the same: the finding code, `plannedChanges: []`,
`planDigest: null`, exit `3`, and nothing written.

## The lock read after the fetch

`buildSkillsReconcilePlan` acquires the source commit through
`withAcquiredRestoreCommit` first, and only then loads `skills-lock.json` with
`loadSkillsLock`. `buildSkillsRestorePlan` in `internal/baseline/skills_restore.go`
loads the lock after the last `acquireRestoreGroup`. Either builder may keep an
earlier read to refuse a malformed or unsafe lock before acquisition, but only
the post-acquisition read feeds the plan digest (`lockBefore`), the lock edits
and the postimage.

`buildSkillsReconcileTransactionDocument` and `buildRestoreTransactionDocument`
take the planned lock bytes (`lock.before`) as a new argument. After capturing
the `skills-lock.json` state, each compares it to those bytes. A missing lock
matches only a missing capture. A present lock matches only a regular file
whose `ContentIdentity` equals `transactionContentIdentity(lock.before)`. On a
mismatch, the builder returns a `SkillsRestoreError` with:

- category `action_required`;
- code `lock.changed-during-plan`;
- message `skills-lock.json changed during planning; the plan no longer
  describes it.`;
- action `Rerun the preview and confirm its new Plan Digest.`

Nothing is written. The restore document checks only when it plans a lock
edit. The existing apply-time revalidation in `applySkillsRestorePlan` is
unchanged, so a rewrite after planning still ends in `plan.confirmation.stale`.
The restore payload shape is unchanged (ADR-0072).

Both command helps add `lock.changed-during-plan` to the exit `3` line.
`docs/user-guide/commands.md` and the Roundfix skill say that the lock is read
after the source is acquired and name the refusal.

## The docs-layout override clause

In `internal/baseline/assets/modules/spec-workflow.json`, the guidance of
`clause.spec.keep-artifacts-in-spec-folder` replaces this sentence:

> Record the approval source, date, covered Spec/revision and actual QA outcome
> or absence; stamp `qa_override: true` and preserve any supplied reason.

with:

> Perform the override only through the runtime's archive override command,
> such as `roundfix archive <slug> --qa-override --approval <source> --reason
> <text>`; that command records the approval source, date, covered Spec revision
> and actual QA outcome or absence, stamps `qa_override: true` and preserves the
> reason. Never hand-edit the override stamp or any other archive front matter,
> because a hand edit skips the command's refusals and provenance.

Every other sentence stays byte-identical, including "A runtime command without
override support must be reported as unsupported rather than given an invented
flag." The versions of `rule.spec.docs-layout`, `guide.spec-docs-layout` and the
module rise by one. The regeneration below renders the change into
`docs/agents/docs-layout.md`, `docs/agents/setup-context.json` and the formatter
golden.

## The incremental verification decision

**Catalog.** `internal/baseline/assets/decisions.json` gains, after
`verification.gate`:

```json
{
  "id": "verification.incremental",
  "version": 1,
  "type": "string",
  "suggestion": "rtk make verify-incremental",
  "summary": "The repository incremental verification command named by generated guidance.",
  "effects": [
    {
      "when": {"present": true},
      "renderBindings": [
        {
          "artifact": "guide.agent-instructions",
          "template": "template.guide.agent-instructions",
          "token": "verification.incremental"
        }
      ]
    }
  ]
}
```

It has no `default`, so `promptBaselineDecision` in
`internal/cli/baseline_human.go` asks for a non-empty value.
`resolveDecisionSuggestions` in `internal/baseline/update.go` offers the
suggestion as a new decision.

`internal/baseline/assets/modules/core.json` lists the decision after
`verification.gate` in `requiredDecisions`, so a Profile Draft must carry it.
The three built-in Profiles list it after `verification.gate` in
`entryDecisions`. The Standard TypeScript Monorepo Profile's `verification`
list stays byte-identical.

`template.guide.agent-instructions` in
`internal/baseline/assets/templates/index.json` gains the token and a version.
Its template gains this line after the gate line:

```text
The selected incremental Verification is {{verification.incremental}}.
```

**Clauses.** In `core.json`, the first three sentences of
`clause.core.verification-two-tiers` become:

> Use the selected incremental Verification named at the top of this guide for
> fast local checks; it answers whether the current change remains valid while
> reusing safe local state. CI must run the selected repository Verification
> from a fresh run; it answers whether the complete tree satisfies the
> repository contract. Baseline planning refuses until the repository selects
> and declares both commands, so neither tier is ever satisfied by omission.

In `spec-workflow.json`, `clause.spec.verification-two-tiers` becomes:

> For each Task, run the selected incremental Verification named in
> `docs/agents/agent-instructions.md` to answer whether the current slice
> remains valid before handoff. CI must run the selected repository
> Verification from a fresh run to answer whether the assembled tree satisfies
> the repository contract. A missing incremental selection is a Baseline
> decision to answer, never a license to skip the local tier or a waiver to
> repeat in each Spec.

The rule, guide and module versions that contain these clauses rise by one.

**Projection.** `resolveVerificationProjection` in
`internal/baseline/profile_alignment.go` projects a `verification.incremental`
decision exactly like `verification.gate`:

- ID `verification.incremental`, role `incremental`, the selected command,
  classification `repository-command`, and `RepositoryExecutable`,
  `DeclarationPath` and `DeclarationDigest` from
  `validateLocalCommandDeclaration`;
- with no local declaration, a blocking `verification.command.undeclared`
  divergence with ID `verification.incremental`, message `selected incremental
  Verification command "<command>" has no matching local declaration`, and the
  gate's next action.

When the decision is present, the built-in Profile loop skips a Profile
declaration with the same ID, so one ID has one projection. When the decision
is absent (a direct `ResolveProfileAlignment` call), the Standard TypeScript
Monorepo's portable `incremental` expectation projects as before. Projections
stay sorted by ID.

**Migration.** No new update code is needed:

- A Setup Manifest without the decision makes `ResolveManifestInput` report it
  in `NewDecisions`. `runBaselineUpdateCommand` then exits `3` with category
  `decision` and writes nothing.
- `--adopt-suggested` adopts `rtk make verify-incremental`, and the plan is
  blocked (`required profile alignment is unresolved: verification.incremental`)
  unless the repository declares it.
- The interactive `roundfix baseline` asks only for that decision.

A greenfield `baseline plan` without the decision exits `3` naming it.

**Fixtures and tests.** Measured on 2026-09-28, the existing test fixtures
pass complete decision lists and write Makefiles with only a `verify:` target.
Each such list gains `verification.incremental`, and each such Makefile gains a
`verify-incremental:` target, in:

- `internal/baseline/plan_test.go`, `internal/baseline/apply_test.go`,
  `internal/baseline/plan_characterization_test.go` and
  `internal/baseline/profile_alignment_test.go`;
- `internal/cli/baseline_plan_test.go`, `internal/cli/baseline_apply_test.go`,
  `internal/cli/baseline_release_gate_test.go`,
  `internal/cli/baseline_update_test.go` and
  `internal/cli/baseline_human_test.go`.

The prompt scripts in `internal/cli/baseline_human_test.go`
(`humanBaselineAdoptionAnswers`, `humanBaselinePreservationAnswers`,
`projectDecisionHumanAnswers` and the inline scripts) gain one answer at the
decision's prompt position.

`TestProfileDeclaresBothVerificationTiers` keeps its name and its independence
assertions. With both decisions, changing either command leaves the other
projection unchanged. Without the incremental decision, the Profile expectation
projects as `profile-expectation`.

`TestBaselineDecisionExamples` and `TestProjectConstraintDocumentation` in
`internal/docscontract/publicdocs_test.go` expect 16 decisions. The frozen
parity corpus stays byte-identical, and no top-level test is removed or
renamed, so `docs/references/coverage-record.json` is untouched.

**Regeneration and docs.**

1. Run `make baseline-digests`.
2. Run `go run -buildvcs=false ./cmd/roundfix baseline update --repo .
   --no-skills --yes --adopt-suggested --format text`. The repository's
   Makefile declares `verify-incremental`, so the decision's projection
   records `Makefile`.
3. Run the refresh again without `--adopt-suggested`. It reports
   `File changes: 0`.

`docs/user-guide/context-driven-development.md` adds the decision to its
published Decision Document and describes the migration. `CONTEXT.md` gains
an **Incremental Verification** entry. The setup-context-driven skill describes
the decision and its migration, and `make skills-sync` regenerates the mirror.

## Regeneration outputs

`make baseline-digests` rewrites these files. They were measured on 2026-09-28
for both catalog changes.

- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, its
  digest pin.
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/`:
  `docs-layout.md` for the override clause, and `agent-instructions.md` and
  `spec-routing.md` for the incremental tier.
- `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest` and `catalog.normalized.json`.
- The plan-characterization goldens `advisory-only-divergences`,
  `clean-adoption`, `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests` and
  `unsatisfied-blocking-capabilities`.

The public Baseline update rewrites `docs/agents/docs-layout.md`,
`docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md` and
`docs/agents/setup-context.json`. Neither step touches the Source Baselines,
`internal/baseline/assets/setups/` or the parity corpus.

## API Contracts

1. `baseline skills reconcile --format json` blocked on a required-removed
   skill: `ok: false`, `plannedChanges: []`, `planDigest: null`,
   `finding.code: reconcile.required-removed`, exit `3`.
2. `baseline skills reconcile` and `baseline skills restore` refuse a lock
   that changed during planning with `finding.code: lock.changed-during-plan`,
   exit `3`, and write nothing.
3. The Baseline catalog declares the string decision `verification.incremental`
   with suggestion `rtk make verify-incremental` and no default, required by
   every built-in Profile and by the core module.
4. A Setup Manifest and a `roundfix/baseline-plan/v1` alignment carry a
   `verification.incremental` Verification projection with role `incremental`.
5. `baseline update` on a manifest without the decision exits `3` with
   category `decision` and `newDecisions` naming it.

## Coverage Map

- Goal 1 → The incremental verification decision; API Contracts 3-4.
- Goal 2 → The incremental verification decision (Migration); API Contract 5.
- Goal 3 → The docs-layout override clause.
- Goal 4 → The lock read after the fetch; API Contract 2.
- Goal 5 → The blocked reconcile result; API Contract 1.
- Core Feature 1 → The blocked reconcile result.
- Core Feature 2 → The lock read after the fetch.
- Core Feature 3 → The docs-layout override clause; Regeneration outputs.
- Core Feature 4 → The incremental verification decision; Regeneration outputs.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- API Contract 1 → The blocked reconcile result.
- API Contract 2 → The lock read after the fetch.
- API Contracts 3-5 → The incremental verification decision.

## Integration Points

- **Spec 0121.** Its Core Feature 6 is delivered here; the operator supersedes
  and archives Spec 0121 at this Spec's archive step.
- **Spec 0148.** Gave the Standard TypeScript Monorepo Profile its
  `incremental` expectation, which stays and projects only without the
  decision.
- **Spec 0163.** Shipped `baseline skills reconcile` and recorded two of its
  limits, closed here.
- **Spec 0169.** Recorded the docs-layout limit closed here.

## Testing Approach

1. **Blocked reconcile.** A package test marshals the blocked payload and finds
   `"plannedChanges":[]` and `"planDigest":null`. A CLI test decodes the JSON,
   asserts an empty array (not `null`), exit `3`, the finding code and an
   unchanged lock, and the help names the exit. The non-blocked preview still
   returns a digest.
2. **Lock read after the fetch.** A test writes a `git` script into a temporary
   directory placed first on `PATH`. The script copies new lock bytes over
   `skills-lock.json` the first time it runs `init --bare`, then execs the real
   `git` resolved beforehand.
   - A preview plans from the new bytes.
   - A confirmed apply of a digest previewed on the old bytes refuses with
     `plan.confirmation.stale` and leaves the new bytes, for reconcile and for
     restore.
   - Each transaction builder refuses a lock whose disk bytes differ from the
     planned bytes with `lock.changed-during-plan`, and accepts equal bytes.
3. **Override clause.** A package test reads the embedded clause and the
   formatter golden: they name the command, contain `Never hand-edit the
   override stamp`, keep the unsupported-runtime sentence, and no longer contain
   `preserve any supplied reason`. A second `baseline update` reports the
   repository current.
4. **Incremental decision.** Package tests cover:
   - the catalog declaration;
   - the projection with a declared and an undeclared command;
   - the planning refusal for a missing decision;
   - the rendered agent-instructions line.

   CLI tests cover `baseline update` on a single-gate manifest: exit `3` naming
   the decision with no writes; `--adopt-suggested` without a declaration
   refused with no writes; with a declaration applied, rendered, recorded and
   converged. The existing Baseline suites stay green with the updated
   fixtures.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The blocked reconcile result (depends on: none).
2. The lock read after the fetch (depends on: 1).
3. The docs-layout override clause (depends on: 2).
4. The incremental verification decision (depends on: 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **One chain.** Tasks 1 and 2 both edit `skills_reconcile.go`, `cli.go`,
  `commands.md` and the Roundfix skill. Tasks 3 and 4 both edit
  `spec-workflow.json`, the regenerated catalog outputs and
  `docs/agents/setup-context.json`. Every Task also shares the implement-task
  instruction, so the graph is a single chain.
- **Fixture breadth.** The incremental decision reaches many Baseline tests.
  Every edit is mechanical: a decision value, a Makefile target or a prompt
  answer. A test that needs any other change is a design signal to report, not
  to paper over.
- **Governed test files.** Three governed test files beyond the two authorized
  on 2026-09-25 (`internal/baseline/plan_test.go`,
  `internal/cli/baseline_human_test.go`,
  `internal/docscontract/publicdocs_test.go`) are bounded in
  [_authorization.md](_authorization.md) with the measurement that requires
  them.
