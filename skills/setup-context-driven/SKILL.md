---
name: setup-context-driven
description: Adopt, update, or verify a repository's Context-Driven Baseline through the public Roundfix Baseline Command. Use when preparing a repository for the write-prd/write-tasks/implement pipeline, changing its Baseline Profile, preserving existing instructions, restoring its Repository Skill Set, or confirming that Baseline-owned content is current.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.0.1
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Context-Driven

Use the public `roundfix baseline` command family. The Roundfix binary is the
only runtime authority for Baseline Profiles, repository inspection, Decision
Plans, Change Plans, confirmation, apply, recovery, and Baseline verification.
This skill supplies recipes and interpretation only. It has no independent
setup engine or behavioral fallback.

Read the durable operating contract before the first adoption or a recovery:
`docs/user-guide/context-driven-development.md#adopt-or-update-the-context-driven-baseline`.

## Boundaries

Baseline owns declared root blocks, setup-owned guides, immutable
root-instruction backups, and `docs/agents/setup-context.json`. It preserves
repository-authored bytes outside managed boundaries.

Baseline never:

- installs application or Agent dependencies;
- connects to a database or live infrastructure;
- runs repository formatter, build, test, lint, migration, or Verification
  commands;
- follows unsafe links or mutates nested instruction carriers;
- infers repository policy from implementation evidence;
- lets ACP proposals authorize mutation;
- combines profiles or loads user-scoped or remote profiles.

Treat profile expectations, repository commands, and recommendations as
different facts. A profile expectation is portable policy, not implementation
proof. A repository command is executable only when the result says
`repositoryExecutable: true` and binds a local declaration. A recommendation
is a command for the maintainer to run; Baseline did not execute it.

## Interpret composed guidance

Read the generated root Instruction Hierarchy from universal instructions
through context and documentation, Spec workflow, enabled autonomous work,
stack guidance, surface guidance, and optional knowledge sources. A narrower
guide may add constraints for its concern but cannot weaken a universal
Normative Clause or confirmed project decision.

Greenfield composition emits only guides and pointers selected by the Baseline
Profile. It creates no generic repository guide and no residual carrier when
there are no Repository-Specific Normative Rules.

For update or Baseline Readoption, interpret the complete retention ledger
before approving it. Every existing rule must retain its exact source bytes in
one active semantic owner, a recognized typed repository document, or
`docs/agents/specific-repository.md`, or have an individual reasoned rejection
or non-governed classification. An empty residual removes the carrier and root
pointer; non-empty accepted residuals retain only their exact rules.

Generated `docs/agents/docs-layout.md` owns the copyable ADR and Findings
contracts. Only `accepted` is active for ADRs; legacy ADRs without lifecycle
frontmatter remain active unless their body marks them inactive. Findings use
`pending`, `partial`, `deferred`, and `done`, keep the original observation,
and append later evidence as dated addenda. Copy those generated templates;
this skill does not render or replace them.

## Human adoption or update

At an interactive terminal, run:

```bash
roundfix baseline --repo . --format text
```

The one workflow detects first adoption or update, collects numbered decisions,
selects exactly one Baseline Profile, and shows one consolidated Change Plan
with file changes first. Do not authorize mutation before the complete plan and
Plan Digest are visible.

Decision and preservation prompts mark one visible default; Enter confirms it.
A valid stored Setup Manifest value wins even when a changed Profile Digest
requires adoption again. Otherwise the CLI uses its embedded catalog
suggestions, including `codex gpt-5.6-sol` for backend work and
`claude opus 5 xhigh` for design work. Existing root instructions default to
Preservation, an empty instruction inventory defaults to Greenfield, and a
recoverable existing profile is preferred. Classification, Plan approval, and
apply still require an explicit non-empty choice.

For first adoption, ask the maintainer to choose one instruction mode:

- Greenfield backs up safe root carriers and imports none of their rules.
- Preservation backs up safe root carriers and requires one reviewed
  disposition per Source Baseline Entry.

Nested instruction carriers remain unchanged and appear only as warnings.
Unsafe root carriers block planning. If sealed ACP classification is
unavailable or discarded, continue through the public command's structured
manual classification review. Never invent Source Baseline identities,
classifications, or destinations.

Repository-Specific Normative Rules have one canonical carrier:
`docs/agents/specific-repository.md`. The enabling decision permits this
destination; it does not create an empty scaffold. Baseline creates the file
and its managed root pointer only for approved non-empty rules. It migrates
either legacy `docs/agents/repository.md` or
`docs/agents/repository-rules.md` byte-for-byte, removes the known empty legacy
scaffold, and blocks divergent non-empty legacy carriers for manual
reconciliation. Existing `0.0.1` Decision Documents that name either legacy
destination are normalized to the canonical path.

For updates, the current Baseline Profile and compatible decisions are offered
first. If only the Profile Digest changed, the command enters adoption but
retains still-valid profile and decision values as defaults. A profile change
produces a new full plan. If the maintainer rejects the plan, use the command's
decision-area revision flow. Every accepted correction must produce a new Plan
Digest and another complete review.

Profile alignment runs before instruction classification. For
profile-specific blockers, the interactive command offers **Change Baseline
Profile**, a reviewed **repository-owned Profile adaptation**, or **Decline
without writing**. An adaptation lists every module and profile-specific
capability removal, validates a repository-owned Profile ID, re-audits the
in-memory draft, and proceeds only when alignment is ready. The Profile file is
part of the final Change Plan and is not written before Plan Digest approval.

Universal required capabilities cannot be removed or waived. Run the exact
restoration preview named by the result, review its Plan Digest, confirm the
same current preview, then rerun Profile alignment:

```bash
roundfix baseline skills restore --repo . --profile <built-in-id> --skill <skill-name> --format json
roundfix baseline skills restore --repo . --profile <built-in-id> --skill <skill-name> --confirm-plan <digest> --format json
```

## Project decisions and Spec constraints

The CLI owns project-decision collection and rendering. This skill does not
collect, derive, validate, or render decisions; it explains the public result
and sends every correction back through `roundfix baseline`.

For the Standard TypeScript Monorepo Profile, UUID version 7 is a visible
suggestion for `identifier.strategy`. The human must explicitly keep it or
provide one non-empty repository-defined rule. When Better Auth remains
selected, `auth.provider` proposes `GET` and `POST` under `/api/auth/*`, owner
`Better Auth`, with the Session, OAuth redirect, callback, and related provider
protocol reason. The human must keep or change that complete typed proposal.

Automation supplies both structured values through the strict Decision
Document passed to `--decision-file`. If either value is unresolved, planning
exits `3`, names every missing decision in `roundfix/baseline-result/v1`, emits
no partial Plan, and writes nothing. Never treat a catalog suggestion,
repository evidence, setup approval, or an empty answer as authorization.

Every new PRD and TechSpec must contain `Project Constraints` rows for:

- Identifier strategy from `docs/agents/domain.md`;
- Authentication and HTTP from `docs/agents/backend.md`;
- active ADR obligations from `docs/agents/spec-routing.md`; and
- tooling authority from `docs/agents/agent-instructions.md`.

Each row must say applicable or not applicable with a reason. Protected
tooling work remains blocked until the Spec records express maintainer
authorization and the exact bounded repository-relative files. Task assignment
and generic implementation requests do not grant that authority. Keep
completed and archived legacy Specs byte-identical.

## Non-interactive planning

The interactive root command refuses redirected or absent terminal input.
Automation and Agents use:

```bash
roundfix baseline plan --repo . --profile <profile-id> --decision-file <decision-file> --format json
```

Automation can pass the equivalent reviewed strict draft:

```bash
roundfix baseline plan --repo . --profile-file <draft.json> --decision-file <decision-file> --format json
```

`--profile-file` and `--profile` are mutually exclusive. The CLI validates the
draft against the embedded catalog, resolves it in memory, and includes its
canonical repository path and exact bytes in the portable Plan. Do not write,
normalize, classify, or render the draft through this skill.

Scalar answers may use repeatable `--decision <id=value>` flags. Structured
answers use repeatable strict Decision Documents with schema
`setup-context-driven/decisions/0.0.1`.

Interpret the result by exit category:

- `0`: stdout is one complete `roundfix/baseline-plan/v1` document.
- `2`: repair invalid input, schema, Git state, or an unsafe carrier.
- `3`: follow the `roundfix/baseline-result/v1` `nextAction`; no partial plan
  exists.

For preservation, use only a current, strict Decision Document derived from the
public human classification review. It must bind the current Source Baseline
and every Source Baseline Entry exactly once. Do not reuse identities or
digests after repository bytes change.

Review at least:

- `repository`, `catalog`, and the one resolved `profile`;
- normalized `decisions`;
- `fileChanges` and the complete canonical `managedEntries` ledger;
- `retention`, `warnings`, and the proposed `setupManifest`;
- every exact `preimage` and `postimage`;
- `planDigest`.

The plan can contain repository policy and exact generated bytes. Store and
transport it as a sensitive review artifact.

## Apply the approved plan

Ask the maintainer to approve the exact current `planDigest`. Then run:

```bash
roundfix baseline apply --repo . --plan <plan-file> --confirm-plan <plan-digest> --format json
```

Apply strictly parses the supplied `roundfix/baseline-plan/v1` document and
applies only its postimages. It never recalculates or substitutes a plan.

On exit:

- `0`: apply and Baseline verification passed, or the exact postimages were
  already verified.
- `1`: inspect execution, verification, recovery, output, or rollback failure.
- `2`: correct invalid input, plan schema, or an unsafe repository.
- `3`: confirmation, preimage, catalog, profile, or Git lineage is no longer
  current. Generate and review a new plan.
- `130`: the operation was canceled.

A portable plan can apply in another clone only when clone-stable Git lineage,
object format, catalog identity, profile, and all bounded preimages match.

## Baseline Profiles

Inspect and validate profiles through the public CLI:

```bash
roundfix baseline profile show <profile-id> --format json
roundfix baseline profile validate <profile-id> --format text
```

Create one repository-owned profile from one built-in profile:

```bash
roundfix baseline profile init --id <repository-profile-id> --from <built-in-id>
roundfix baseline profile validate <repository-profile-id> --format text
```

Repository-owned profiles live only under
`.roundfix/baseline/profiles/<id>.json`, compose only embedded catalog entry
IDs, and must be versioned with the repository. They cannot combine profiles,
load user-scoped catalogs, or declare executable or remote content.

## Repository Skill Set restoration

Adoption reports missing or drifted external skills but never restores them as
a side effect. Preview an explicitly requested restoration:

```bash
roundfix baseline skills restore --repo . --profile <built-in-id> --skill <skill-name> --format json
```

A non-empty preview exits `3` with a Plan Digest and writes nothing. After the
maintainer approves that exact preview, rerun:

```bash
roundfix baseline skills restore --repo . --profile <built-in-id> --skill <skill-name> --confirm-plan <plan-digest> --format json
```

Use `--source-dir <path>` only for a declared offline Git checkout or bare
object store containing the exact immutable source commit. Never substitute a
generic skill refresh.

After restoration, generate the same Baseline Plan again. Restoration success
proves selected Repository Skill Set bytes, not repository-authored policy.

## Canonical asset synchronization

This is a maintainer-only operation and is never part of ordinary adoption:

```bash
roundfix baseline assets sync --source-dir <canonical-setups> --check --format json
```

Run without `--check` only when the maintainer explicitly authorizes refreshing
the Go-owned canonical setup snapshots. The command validates source
provenance and the generated catalog before changing only canonical setup
assets.

## Recovery

- Missing decisions or manual classification: follow `nextAction`, complete
  current structured input, and rerun planning.
- Confirmation mismatch: copy the full current Plan Digest from the supplied
  plan after review.
- Stale plan: generate and approve a replacement; never edit or force the old
  plan.
- Unsafe carrier: repair the exact reported path; never follow or replace it
  through an ad-hoc workaround.
- Interrupted transaction: rerun the same public apply command so Roundfix can
  lock and recover its Git-private journal. Never delete recovery state.
- Incomplete rollback: stop writes, preserve checkout and recovery state,
  inspect every reported path, restore trustworthy bytes, and generate a new
  plan.

## Completion

Do not claim adoption or update complete from planning alone.

1. Require apply exit `0` and inspect `verifiedPostimages`.
2. Run each reported formatter or Verification recommendation outside
   Baseline.
3. Repair any repository failure outside managed Baseline output.
4. Generate a fresh plan with the same current Profile input and decisions.
5. Require no file changes, or perform an exact idempotent reapply at exit `0`.

Report remaining warnings, recommendations not run, the approved Plan Digest,
and the final Baseline result state.
