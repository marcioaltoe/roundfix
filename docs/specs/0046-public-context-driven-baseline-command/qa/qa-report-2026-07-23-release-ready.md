---
spec: 0046-public-context-driven-baseline-command
date: 2026-07-23
build: 14b30db+task_17-wip.2
status: in-progress
verdict: pending
surfaces: [cli, infra, docs]
---

# QA report — Public Context-Driven Baseline Command

## Scope and environment

This full release gate covers every PRD User Story, Core Feature, Success
Metric, finding correction, maintained Baseline Profile, named product journey,
and applicable Non-Goal. The hermetic rows use a freshly built Roundfix binary
and disposable committed Git repositories. The live Fluxus row requires
separate explicit authorization and cannot inherit authorization from the
hermetic gate.

Actor names:

- Maintainer: a repository owner using the human command.
- Automation: a non-interactive caller using portable JSON plan/apply.
- Agent: an Agent following the shipped setup skill.
- Release verifier: the person running repository and composition gates.

Every CLI row records the command exit, stdout result, stderr diagnostics, a
second public read or audit, and the repository tree before and after. Apply
rows also restart through a new process and finish with an empty reapply.

## Static gate

| Gate | Command | Expected | Evidence | Status |
| --- | --- | --- | --- | --- |
| SG-01 | `rtk go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run 'TestBaselineMacroJourneys\|TestBaselineFindingRegressions\|TestBaselineFormatterComposition\|TestBaselineDocumentationContract'` | Named release suites pass | Exit 0 in all four packages after the WIP.2 repair | pass |
| SG-02 | `rtk grep -R -n 'status: passed' docs/specs/0046-public-context-driven-baseline-command/qa` | This report records the separately authorized live result | Pending live authorization | pending |
| SG-03 | `rtk make verify` | Full repository Verification passes after evidence is final | Pre-live WIP.2 run passed 2,077 tests in 22 packages, the four cutover guards, Repository Skill Set checks, and the binary build; rerun remains required after live evidence | pending |

## Hermetic evidence

- Pre-change discovery:
  `rtk env GOCACHE=/private/tmp/roundfix-task17-go-cache go test
  ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run '^$'
  -list
  'TestBaselineMacroJourneys|TestBaselineFindingRegressions|TestBaselineFormatterComposition|TestBaselineDocumentationContract'`
  listed only `TestBaselineDocumentationContract`. The required macro,
  finding-regression, formatter-composition, and thin-skill gates did not
  exist.
- The first release-suite run failed only
  `TestBaselineFormatterComposition/standard-typescript-monorepo` because the
  frontend guide could not render repository-owned token
  `repository.design-contract`.
- WIP.2 binds safe repository-owned reference paths at the same artifact
  rendering boundary that already binds setup-owned `managedId` references.
  The isolated TypeScript formatter-composition test then passed.
- The exact SG-01 command passed after the repair. Its real-process journeys
  built Roundfix, exchanged portable plans across processes and clones,
  applied exact digests, ran Go, Rust, and pinned Oxfmt formatters externally,
  ran each repository Verification externally, and completed empty reapply
  with unchanged managed bytes.
- The first sandboxed `rtk make verify` attempt could not read the host Go
  build cache. The unchanged command with the required cache access passed:
  2,077 tests in 22 packages, four setup-skill cutover guards, Repository
  Skill Set checks, and the final binary build.

## User story coverage

| ID | Story | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| US-01 | One guided adoption or update flow | Maintainer / `roundfix baseline` / CLI | Complete preflight, decisions, review, confirmation, apply; observe verified result | New process audits the Setup Manifest; empty reapply is unchanged | `TestBaselineMacroJourneys/human_adoption_and_update` | pass |
| US-02 | Stable non-interactive operation | Automation / `baseline plan` and `baseline apply` / CLI | Emit strict JSON, apply only approved digest, preserve stdout/stderr contract | Parse both documents through public schemas and audit in a new process | `TestBaselineMacroJourneys/automation_plan_apply` | pass |
| US-03 | Explicit greenfield or preservation choice | Maintainer / human and automation adoption / CLI | Exercise both modes; observe immutable backup and selected dispositions | Read backup bytes and Repository-Specific Normative Rules, then reapply | `TestBaselineMacroJourneys/greenfield_and_preservation` | pass |
| US-04 | Select exactly one Baseline Profile | Maintainer / profile selection / CLI | Plan each maintained profile and one explicit profile change | Manifest and audit report exactly one selected profile | `TestBaselineFormatterComposition` | pass |
| US-05 | Consolidated editable classification proposal | Maintainer / preservation review / CLI | Review and edit one consolidated proposal before planning | Decline once and compare complete tree; approved retry persists exact decision | `TestBaselineMacroJourneys/consolidated_review` | pass |
| US-06 | Revisit a rejected plan decision | Maintainer / rejected final review / CLI | Select one revision area, revise, and observe a new digest | Old digest is rejected; new review is deterministic on replay | `TestBaselineMacroJourneys/rejected_plan_revision` | pass |
| US-07 | File projection plus complete ledger | Maintainer / final plan review / CLI | Observe `fileChanges` before the complete managed-entry ledger | Strict parser recomputes projection; audit confirms postimages | `TestBaselineFindingRegressions/file_level_projection` | pass |
| US-08 | Capability and Verification diagnostics stay distinct | Maintainer / TypeScript plan / CLI | Observe implementation evidence, missing contract, roles, and local commands separately | Generated manifest labels only locally validated commands executable | `TestBaselineFindingRegressions/capability_and_verification_diagnostics` | pass |
| US-09 | Setup skill delegates to public CLI | Agent / `.agents/skills/setup-context-driven/SKILL.md` / docs | Follow each recipe through the public parser | Shipped-skill check rejects any executable engine file | `TestBaselineDocumentationContract` and `TestThinSetupSkill` | pass |
| US-10 | Current operating and recovery documentation | Roundfix user / user guide / docs | Follow adoption, update, automation, recovery, and migration examples | Every example dispatches through the shipped parser | `TestBaselineDocumentationContract` | pass |

## Core feature coverage

| ID | Core Feature | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| CF-01 | Public CLI is the sole Baseline authority | Maintainer / `roundfix baseline` / CLI | Complete plan and apply without an installed skill or Python | Run embedded profile show with helper tools absent | `TestBaselineMacroJourneys/public_authority` | pass |
| CF-02 | One deterministic state machine | Maintainer / root command / CLI | Traverse preflight through Baseline verification without skipped states | Automation-equivalent input produces the same Plan Digest | `TestBaselineMacroJourneys/state_machine_parity` | pass |
| CF-03 | Preflight selects adoption or update | Maintainer / root command / CLI | Run against unconfigured, configured, and profile-change repositories | New process reports expected current profile and workflow | `TestBaselineMacroJourneys/human_adoption_and_update` | pass |
| CF-04 | Built-in and repository-owned profiles | Maintainer / `baseline profile` / CLI | Show built-ins; init, validate, and plan a repository-owned profile | Invalid composition and remote/custom executable content are rejected | `TestBaselineMacroJourneys/profile_selection` | pass |
| CF-05 | Two instruction-preservation modes | Maintainer / adoption / CLI | Greenfield backs up without import; preservation requires dispositions | Exact backup and retained rule bytes survive restart | `TestBaselineMacroJourneys/greenfield_and_preservation` | pass |
| CF-06 | Bounded audit; root-only preservation | Maintainer / plan / CLI | Include root and nested carriers; mutate only root carriers | Nested identities match before/after and warning persists | `TestBaselineMacroJourneys/nested_carrier_boundary` | pass |
| CF-07 | Safe aliases retain evidence; unsafe aliases block | Maintainer / plan / CLI | Exercise safe, external, escaping, cyclic, unreadable, and special targets | Safe target backup is exact; each unsafe case has zero mutation | `TestBaselineMacroJourneys/unsafe_carriers` | pass |
| CF-08 | Preferred/fallback/manual classification | Maintainer / preservation / CLI | Accept preferred, discard preferred and accept fallback, then use manual fallback | Snapshot bytes match; no ACP output authorizes mutation | `TestBaselineMacroJourneys/classification_proposals` | pass |
| CF-09 | ACP proposals remain review-only | Maintainer / consolidated review / CLI | Decline proposal and observe zero writes; approve deterministic disposition | Applied Plan Digest excludes proposal origin and raw ACP data | `TestBaselineMacroJourneys/classification_proposals` | pass |
| CF-10 | Required divergence blocks; advisory stays visible | Maintainer / TypeScript plan / CLI | Omit required evidence, add it, and leave advisory evidence absent | Only required resolution changes readiness | `TestBaselineFindingRegressions/capability_and_verification_diagnostics` | pass |
| CF-11 | Rejected plan returns to structured revision | Maintainer / root command / CLI | Revise profile, rules, divergence, or files; optionally propose scoped change | Each accepted revision creates a fresh approval digest | `TestBaselineMacroJourneys/rejected_plan_revision` | pass |
| CF-12 | File projection derives from canonical ledger | Automation / JSON plan / CLI | Inspect one row per path and ordered managed IDs | Strict parser rejects a mismatched projection | `TestBaselineFindingRegressions/file_level_projection` | pass |
| CF-13 | Apply is exact, stale-safe, rollback-capable, idempotent | Automation / apply / CLI | Apply correct digest, stale plan, injected failure, and empty reapply | Repository tree is exact after rollback/recovery and unchanged on reapply | `TestBaselineMacroJourneys/apply_safety` | pass |
| CF-14 | Baseline verification owns only Baseline outcomes | Release verifier / apply and external commands / CLI/infra | Apply, then run formatter and Repository Verification externally | Trap commands prove Baseline never executed them | `TestBaselineFormatterComposition` | pass |
| CF-15 | Automation uses complete plan then exact apply | Automation / plan/apply / CLI | Plan once, apply same file/digest, then mutate bounded preimage | Stale apply exits actionable and writes nothing | `TestBaselineMacroJourneys/automation_plan_apply` | pass |
| CF-16 | Stable schemas, ordering, streams, and exits | Automation / JSON commands / CLI | Exercise success, usage, unsafe, action-required, execution, cancellation contracts | Strict parser accepts results; stderr contains diagnostics only | `TestBaselineMacroJourneys/machine_contract` | pass |
| CF-17 | Audit is local, read-only, and network-free | Maintainer / TypeScript plan / CLI | Discover routes, PostgreSQL evidence, and commands with process/network traps | Trap counters remain zero and repository tree is unchanged | `TestBaselineFindingRegressions` | pass |
| CF-18 | Docs and parser describe one contract | Agent / docs and skill recipes / docs/CLI | Dispatch every example and parse Decision Documents | Unknown fields remain rejected by the same parser | `TestBaselineDocumentationContract` | pass |
| CF-19 | Spec 0045 safety properties remain binding | Release verifier / parity corpus and macro flows / infra | Validate Readoption, dispositions, retention, formatter stability, ownership, digest, rollback | Compatibility corpus and macro post-state both pass | `TestBaselineMacroJourneys/spec_0045_safety` | pass |
| CF-20 | Every maintained Python contract has a Go destination | Release verifier / compatibility corpus / infra | Validate all 243 rows and 26 Go destinations after Python removal | No retired or unassigned maintained row exists | `TestBaselineCompatibilityCorpus` | pass |
| CF-21 | Setup skill has no executable engine | Agent / shipped setup skill / docs | Inspect canonical, distributed, and embedded skill trees | Negative guard rejects a reintroduced executable artifact | `TestBaselineDocumentationContract` and `TestCheckRejectsExecutableSetupEngineArtifacts` | pass |
| CF-22 | User docs cover every operating flow | Roundfix user / user guide / docs | Follow adoption, update, preservation, profile, revision, automation, recovery, troubleshooting, migration | Documentation contract dispatches all public examples | `TestBaselineDocumentationContract` | pass |

## Success metric coverage

| ID | Metric | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| SM-01 | Seven finding categories have passing evidence | Release verifier / macro suite / CLI | Execute each named regression | Each row points to an observable public result | `TestBaselineFindingRegressions` | pass |
| SM-02 | All maintained Python contracts have Go parity | Release verifier / parity corpus / infra | Validate matrix, fixtures, and destinations | Skill engine remains removed | `TestBaselineCompatibilityCorpus` | pass |
| SM-03 | Profiles and states preserve normalized outcomes and safety | Release verifier / all profiles / CLI | Plan/apply maintained profiles and refusal states | Final audit, rollback, and empty reapply match | `TestBaselineFormatterComposition` and `TestBaselineMacroJourneys` | pass |
| SM-04 | Equivalent interactive and automation inputs have identical digests | Maintainer and Automation / root and plan commands / CLI | Provide equivalent answers | Compare exact Plan bytes and digest | `TestBaselineMacroJourneys/state_machine_parity` | pass |
| SM-05 | Every modified root carrier has an exact immutable backup | Maintainer / preservation and greenfield / CLI | Apply root-carrier changes | Backup name digest and bytes match after restart | `TestBaselineMacroJourneys/greenfield_and_preservation` | pass |
| SM-06 | Nested carriers receive zero automatic mutation | Maintainer / plan/apply / CLI | Include nested carrier and conflict | Byte identity stays fixed and warning remains | `TestBaselineMacroJourneys/nested_carrier_boundary` | pass |
| SM-07 | ACP proposals never mutate or bypass review | Maintainer / proposal review / CLI | Exercise preferred, fallback, manual, decline, and approval | Pre-review trees match; approved plan is deterministic | `TestBaselineMacroJourneys/classification_proposals` | pass |
| SM-08 | Unsafe targets are never followed as trusted source | Maintainer / plan / CLI | Exercise every unsafe target category | No target bytes enter trusted evidence and tree is unchanged | `TestBaselineMacroJourneys/unsafe_carriers` | pass |
| SM-09 | Plans expose file and managed-entry views | Maintainer and Automation / plan / CLI | Read text and JSON views | Strict parser recomputes and validates projection | `TestBaselineFindingRegressions/file_level_projection` | pass |
| SM-10 | Only locally validated commands are repository-executable | Maintainer / TypeScript plan / CLI | Compare absent and declared scripts/targets | Manifest and recommendations preserve the distinction | `TestBaselineFindingRegressions/verification_commands` | pass |
| SM-11 | Baseline executes no dependency, DB, formatter, or Verification command | Release verifier / plan/apply/audit / infra | Install execution traps and run all Baseline phases | Trap counters remain zero | `TestBaselineFormatterComposition/non_execution` | pass |
| SM-12 | Generated output composes with formatter, Verification, audit, and empty reapply | Release verifier / all profiles / CLI/infra | Apply, run external tools, audit, reapply | Managed-file identity set remains unchanged | `TestBaselineFormatterComposition` | pass |
| SM-13 | Authorized Fluxus journey needs no four manual corrections | Maintainer / live Fluxus repository / CLI | Adopt or update using current public workflow | Record identity, approved digest, audit, recovery, and absence of seven friction outcomes | Live evidence in this report | pending |
| SM-14 | Setup skill contains zero executable setup scripts | Agent / shipped skill / docs | Inspect all shipped copies | Negative validator catches a synthetic engine file | `TestBaselineDocumentationContract` | pass |
| SM-15 | All documented examples pass contract checks | Roundfix user / docs / docs/CLI | Parse and dispatch every fenced example | Public help and strict Decision parser agree | `TestBaselineDocumentationContract` | pass |
| SM-16 | Full repository Verification passes | Release verifier / `make verify` / infra | Run after final QA evidence is recorded | Exit 0 with no skipped repository gate | SG-03 | pending |

## Fluxus finding regressions

| ID | Finding category | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| FX-01 | Safe alias evidence survives remediation | Maintainer / plan / CLI | Audit a safe alias before any mutation | Alias, target, relationship, and target digest remain in plan | `TestBaselineFindingRegressions/safe_alias_evidence` | pass |
| FX-02 | Decisions are collected without repeated full-plan noise | Maintainer / root command / CLI | Answer all required decisions once | One consolidated final plan appears after decision collection | `TestBaselineFindingRegressions/consolidated_decisions` | pass |
| FX-03 | Emitted Decision Documents include the required schema | Automation / emitted skeleton / CLI | Parse the emitted document without repair | Strict parser accepts it and rejects a missing schema field | `TestBaselineFindingRegressions/decision_document_schema` | pass |
| FX-04 | HTTP candidates include bounded facts and source digests only | Maintainer / TypeScript audit / CLI | Inspect route candidates | Facts include path/method/scope/digest and no policy inference | `TestBaselineFindingRegressions/http_candidates` | pass |
| FX-05 | PostgreSQL implementation and contract evidence are distinct | Maintainer / TypeScript audit / CLI | Supply implementation evidence without contract, then add contract | Diagnostic names accepted contract paths and preserves implementation evidence | `TestBaselineFindingRegressions/postgresql_diagnostics` | pass |
| FX-06 | Plan review leads with one file row per path | Maintainer / final review / CLI | Render plan with multiple managed entries per path | `fileChanges` precedes complete canonical ledger | `TestBaselineFindingRegressions/file_level_projection` | pass |
| FX-07 | Persisted commands are locally validated | Maintainer / TypeScript audit / CLI | Compare profile expectations with repository scripts and Make targets | Only exact local declarations become repository-executable | `TestBaselineFindingRegressions/verification_commands` | pass |

## Product journeys

| ID | Journey | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| J-01 | Greenfield adoption | Maintainer / root and plan/apply / CLI | Back up root carriers, import no rules, apply verified plan | Audit and empty reapply in fresh processes | `TestBaselineMacroJourneys/greenfield` | pass |
| J-02 | Instruction preservation | Maintainer / root command / CLI | Classify every root entry and apply retained rules | Backup and repository-rule bytes match source | `TestBaselineMacroJourneys/preservation` | pass |
| J-03 | Update | Maintainer / root command / CLI | Detect current profile and retain current decisions | New process reports update state and clean audit | `TestBaselineMacroJourneys/update` | pass |
| J-04 | Profile change | Maintainer / root command / CLI | Change from Go to Rust profile and approve new digest | Manifest reports Rust; old digest cannot apply | `TestBaselineMacroJourneys/profile_change` | pass |
| J-05 | Manual ACP fallback | Maintainer / preservation / CLI | Make both selections unavailable or invalid and classify manually | Complete deterministic proposal names backup and repository-rules destinations | `TestBaselineMacroJourneys/manual_fallback` | pass |
| J-06 | Preferred/fallback proposal | Maintainer / preservation / CLI | Accept preferred once; reject preferred and accept fallback once | Both attempts receive byte-identical snapshot; no tool or checkout mutation | `TestBaselineMacroJourneys/preferred_fallback` | pass |
| J-07 | Rejected-plan revision | Maintainer / root command / CLI | Reject, revise one scoped decision, and review again | New digest differs; replay matches; original digest is stale | `TestBaselineMacroJourneys/rejected_plan_revision` | pass |
| J-08 | Automation plan/apply | Automation / plan and apply / CLI | Exchange one portable JSON plan and exact approval | Public audit and empty reapply confirm persisted state | `TestBaselineMacroJourneys/automation_plan_apply` | pass |
| J-09 | Stale plan | Automation / apply / CLI | Change a consulted and a mutation-target preimage before apply | Exit is actionable and complete tree remains unchanged | `TestBaselineMacroJourneys/stale_plan` | pass |
| J-10 | Cross-clone apply | Automation / plan in clone A, apply in clone B / CLI | Use matching lineage and preimages; also try unrelated history | Matching clone verifies; unrelated clone writes nothing | `TestBaselineMacroJourneys/cross_clone` | pass |
| J-11 | Unsafe carrier | Maintainer / plan / CLI | Exercise external, escape, cycle, unreadable, and special targets | Each blocks before trusted read or mutation | `TestBaselineMacroJourneys/unsafe_carriers` | pass |
| J-12 | Rollback | Automation / apply / CLI | Inject post-write verification failure | Exact preimage returns and result names failure | `TestBaselineMacroJourneys/rollback` | pass |
| J-13 | Recovery | Automation / apply / CLI | Abandon an interrupted journal, then start another apply | Next process recovers exact tree before proceeding | `TestBaselineMacroJourneys/recovery` | pass |
| J-14 | Empty reapply | Automation / apply / CLI | Reapply already verified exact plan | Exit 0, verified noop, zero managed-file delta | `TestBaselineMacroJourneys/empty_reapply` | pass |

## Maintained profile composition

| ID | Profile | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| PF-01 | `go-cli-tui` | Release verifier / public plan/apply plus external commands / CLI/infra | Plan, apply, run external formatter and `make verify`, audit, reapply | Managed identities remain unchanged and Baseline trap counters stay zero | `TestBaselineFormatterComposition/go-cli-tui` | pass |
| PF-02 | `rust-cli` | Release verifier / public plan/apply plus external commands / CLI/infra | Plan, apply, run external formatter and repository Verification, audit, reapply | Managed identities remain unchanged and Baseline trap counters stay zero | `TestBaselineFormatterComposition/rust-cli` | pass |
| PF-03 | `standard-typescript-monorepo` | Release verifier / public plan/apply plus external commands / CLI/infra | Plan, apply, run pinned Oxfmt and repository Verification, audit, reapply | Managed identities remain unchanged and Baseline trap counters stay zero | `TestBaselineFormatterComposition/standard-typescript-monorepo` | pass |

## Non-Goal checks

| ID | Exclusion | Actor / entry / surface | Steps and expected observable | Independent confirmation and persistence | Evidence | Status |
| --- | --- | --- | --- | --- | --- | --- |
| NG-01 | No user-scoped custom profile catalog | Maintainer / profile commands / CLI | Place a user-scoped profile and plan by ID | Profile remains undiscoverable | Existing profile rejection suite plus macro profile journey | pass |
| NG-02 | No profile composition | Maintainer / profile validate / CLI | Declare composed profiles | Strict validation rejects it | Existing custom-profile rejection suite | pass |
| NG-03 | No external profile registry or plugins | Maintainer / profile validate / CLI | Declare remote or plugin fields | Strict validation rejects them | Existing custom-profile rejection suite | pass |
| NG-04 | No Python fallback | Agent / setup skill and binary / docs/CLI | Run with Python unavailable and inspect shipped skill | Public Go command works; no engine file exists | `TestBaselineDocumentationContract` | pass |
| NG-05 | Setup Command remains unchanged | Roundfix user / `roundfix setup` / CLI | Run compatibility fixture | Healthy-machine and profile-proof contract remains unchanged | `TestSetupCommandCompatibility` | pass |
| NG-06 | No bypassing update path | Maintainer / root and automation commands / CLI | Attempt update without preflight/audit/approval | Workflow returns the required next action | `TestBaselineMacroJourneys/update` | pass |
| NG-07 | No nested carrier mutation | Maintainer / plan/apply / CLI | Include nested carriers | Nested paths remain byte-identical | `TestBaselineMacroJourneys/nested_carrier_boundary` | pass |
| NG-08 | ACP cannot authorize mutation | Maintainer / proposal path / CLI | Supply proposal then decline review | Complete tree remains unchanged | `TestBaselineMacroJourneys/classification_proposals` | pass |
| NG-09 | No policy inference from source | Maintainer / TypeScript audit / CLI | Inspect HTTP and capability candidates | No owner, reason, mode, or Normative Clause is inferred | `TestBaselineFindingRegressions/http_candidates` | pass |
| NG-10 | No unsafe-link following | Maintainer / plan / CLI | Exercise unsafe carrier matrix | No opaque target bytes become trusted evidence | `TestBaselineMacroJourneys/unsafe_carriers` | pass |
| NG-11 | No dependency, database, or repository-command execution | Release verifier / plan/apply/audit / infra | Install traps and run Baseline | All trap counters stay zero | `TestBaselineFormatterComposition/non_execution` | pass |
| NG-12 | Repository policy stays repository-owned | Maintainer / preservation / CLI | Retain a repository rule and apply | Managed ledger excludes retained rule ownership | `TestBaselineMacroJourneys/preservation` | pass |

## Acceptance criteria

| ID | Criterion | Evidence | Status |
| --- | --- | --- | --- |
| AC-01 | All 10 User Stories, 22 Core Features, and Success Metrics have fresh passing evidence | US, CF, and SM tables | pending |
| AC-02 | All seven Fluxus finding categories have named regression evidence | FX table | pass |
| AC-03 | Every maintained profile completes plan, apply, formatter/Verification composition, final audit, and empty reapply | PF table | pass |
| AC-04 | Failure journeys prove zero unauthorized or unrecoverable mutation | J-09 through J-13 and tree evidence | pass |
| AC-05 | Interactive and non-interactive equivalent inputs produce identical Plan Digests | US-02, CF-02, SM-04 | pass |
| AC-06 | Live Fluxus completes only after separate authorization and without the four manual corrections | SM-13 and live evidence | pending |
| AC-07 | Full repository Verification passes after all QA evidence is recorded | SG-03 | pending |

## Findings

### F-01 — Repository-owned guide references did not render

Impact: Blocks-Completion. A maintainer selecting
`standard-typescript-monorepo` could not produce a Change Plan because
`guide.frontend` failed on the missing `repository.design-contract` render
value. The profile declares that reference as repository-owned path
`DESIGN.md`; it is not a repository decision.

Reproduction:

1. Run
   `TestBaselineFormatterComposition/standard-typescript-monorepo` against
   WIP.1.
2. Observe exit 2 and
   `render managed entry "guide.frontend": missing render value for token
   "repository.design-contract"`.

WIP.2 resolves safe repository-owned `path` references alongside setup-owned
`managedId` references. The isolated reproduction and the complete SG-01
command both pass. Affected rows: US-04, CF-04, SM-03, SM-12, PF-03.

## Blocked and skipped

No row is blocked or skipped. SG-02, SG-03, SM-13, SM-16, AC-01, AC-06, and
AC-07 remain pending until the live Fluxus journey receives separate explicit
authorization and the final repository gate is rerun afterward.

## Coverage

Planned: 10 User Stories, 22 Core Features, 16 Success Metrics, seven Fluxus
finding categories, 14 product journeys, three maintained profiles, 12
applicable Non-Goal checks, and seven Task acceptance criteria.

## Final verdict

Pending. The report remains open until every row is settled and the live
Fluxus journey has separate explicit authorization.
