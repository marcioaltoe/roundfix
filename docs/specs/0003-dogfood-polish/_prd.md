---
spec: 0003-dogfood-polish
status: active
created: 2026-07-05
surfaces: [cli, docs]
---

# Dogfood Polish

The first two dogfood cycles (Implement Command on spec 0002, watch on a real
pull request) shipped working software and a list of small, confirmed
irritations: contract inconsistencies, misleading header lines, an
environment-sensitive test suite, and opaque infrastructure errors. Each item
is small, already diagnosed in the dogfood findings log, and none needs new
architecture — this Spec batches them so the debris is cleared before the
larger watch and TUI cycles, and so the 0002 QA gate can pass on re-run.

## Goals

- Every generated commit message follows one convention: lowercase Conventional
  Commits subjects, unscoped, for Task commits and the QA Report commit alike.
- The Implement Command's Run header shows only facts that apply to spec Runs.
- The full test suite passes on any machine regardless of the user's global
  git configuration — the environmental QA `fail` on spec 0002 becomes a
  `pass` on re-run.
- Infrastructure errors from the agent layer name enough of the underlying
  tool's output to be actionable without digging into logs.
- Small operability gaps from the findings log are closed: stop by spec
  target, QA opt-in in Interactive Input, spec discovery diagnostics, and one
  agent-log location rule.

## User Stories

1. As a developer reading branch history after a spec Run, I want generated
   commit subjects lowercase and unscoped, so that Roundfix commits read like
   the repository's own style.
2. As a developer starting a spec Run, I want the Run header without the Run
   Budget and Round lines that do not apply, so that the header states only
   real facts.
3. As a developer on a machine with global commit signing, I want the test
   suite and the QA gate to pass without environment overrides, so that
   verification results reflect the code, not my git config.
4. As a developer hitting an agent-layer failure, I want the error to include
   the tool's own stderr tail, so that the next action is obvious.
5. As a developer with a stuck spec Run lock, I want `roundfix stop --spec
   <slug>`, so that releasing a spec target does not require hunting a run id.
6. As a developer using Interactive Input for implement, I want the QA opt-in
   offered as a field, so that `--qa` is not flag-only knowledge.
7. As a developer with a broken spec folder, I want the Spec picker to tell me
   which folder was skipped and why, so that a typo in one `_prd.md` is not
   silently invisible.
8. As a developer configuring `defaults.artifact_dir`, I want spec-Run agent
   logs to honor it like review Runs do, so that one setting governs all logs.

## Core Features

1. **Commit subject normalization.** Task commit subjects lowercase their
   first rune; the QA Report commit drops its scope (`docs: qa report for
   <slug> (<verdict>)`). Trailers are unchanged.
2. **Header truthfulness.** Implement Runs omit the Run Budget and Round lines.
3. **Hermetic git tests.** Test helpers that create git repositories isolate
   them from user/global git config (no inherited signing, deterministic
   identity), matching the daemon's own explicit-config discipline.
4. **Actionable infrastructure errors.** Agent-layer infrastructure errors
   include the failing tool's trailing stderr (bounded) in the message.
5. **Stop by spec target.** `roundfix stop --spec <slug>` resolves the Active
   Run for the current repository's spec target, mirroring `--pr`.
6. **QA field in Interactive Input.** The implement flow offers the QA opt-in
   as a yes/no field defaulting to no.
7. **Spec discovery diagnostics.** The Spec picker lists skipped folders with
   the reason (stderr), while keeping the picker itself robust.
8. **One agent-log rule.** Spec Runs place agent logs under the configured
   Artifact Directory exactly as review Runs do.
9. **Docs.** The repo agent instructions gain the commit-style line (unscoped
   subjects); the shipped Roundfix skill is updated wherever these behaviors
   are documented.

## User Experience

No new commands. Two flag surfaces grow (`stop --spec`, an Interactive Input
field); all other changes make existing output truthful or quieter. Exit
codes, stdout shapes, and journal kinds are untouched except where noted.

## Non-Goals / Out of Scope

- The watch loop changes (poll-first, stdout contract, merge-readiness) — next
  Spec.
- TUI redesign — its own Spec.
- Worktree-per-task, in-repo review artifacts, permission policy, templated
  prompts (work-plan items).
- `roundfix setup`/`upgrade` and graceful-stop semantics (findings 22–24) —
  future command-lifecycle Spec.

## Success Metrics

- Re-running the 0002 QA gate on this machine yields `verdict: pass` with no
  environment overrides.
- A spec Run's commits and QA commit pass `cog verify` under this repo's
  configuration.
- The full suite passes with `commit.gpgsign=true` globally and no overrides.
- Existing stdout/exit-code contract tests pass unchanged except the header
  lines deliberately removed.

## Decisions

- QA Report commit message becomes unscoped — consistency with Task commits
  beats the decorative scope; the product contract text updates accordingly
  (amends the 0001 techspec's commit contract detail; no ADR needed).
- Spec-Run agent logs move under the Artifact Directory (finding 8's "unify"
  option) rather than documenting the split.
- Task titles stay as written in task files; normalization happens only in the
  commit-subject derivation.

## Open Questions

None.
