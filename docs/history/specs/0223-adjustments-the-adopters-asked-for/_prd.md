---
spec: 0223-adjustments-the-adopters-asked-for
status: archived
created: 2026-10-04
surfaces: [backend, cli, docs]
archived: "2026-10-04"
source_slug: 0223-adjustments-the-adopters-asked-for
---


# Adjustments the adopters asked for

Two repositories that adopted Roundfix, and Roundfix's own release step, hit
three small walls in the same week. The oraculum maintainer could not leave the
branch prefix to the commit-type rule, because the Baseline made the decision
mandatory and always rendered a prefix sentence. fluxus found that a QA gate
reopened over a report written before row inputs existed had no valid move
under the `qa-gate` skill. And `roundfix release plan`, the first command of
every release, says nothing about the skills and guides check that the release
runbook makes mandatory. This Spec removes each wall without changing what any
adopter has already recorded. The fourth request of the same week, relative
links that the archive breaks, needs its own grant and is owned by Spec 0226.

## Prerequisites

None. Specs 0222, 0224 and 0226 are authored in the same cycle and may also
raise the Roundfix Skill's version; the operator orders the queue so that each
later Spec raises it from the earlier one's value, and task_01 raises it from
the value on the tree it starts from.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  decision keeps its identifier `branch.prefix`, and Specs, Tasks and QA rows
  keep theirs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change is a local file
  read or local guidance; no request, credential or forge read is added, and
  `roundfix release plan` keeps its offline contract. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0228 (this Spec) decides that an
  unrecorded branch prefix follows the commit-type rule. ADR-0205 is the
  optional-decision precedent, ADR-0205: "The decision is optional: while a
  repository records none, the Baseline states the suggestion". ADR-0118 made
  the prefix a decision and ADR-0150 made its value purpose-based, ADR-0150:
  "Keep the compatible `branch.prefix` string key with `<type>/` as its pattern
  default"; both hold, and ADR-0228 narrows ADR-0118 for this one decision.
  ADR-0222 holds that a guide says only what holds for the repository that
  reads it, ADR-0222: "four sentences the Baseline renders still did not hold
  for every repository that reads them", and the unrecorded wording obeys it.
  ADR-0187 splits the Roundfix Skill by command and ADR-0189 ties an owned
  skill's version to its content, so each skill edit raises its version.
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts", applied to
  the release plan's new lines. The gate is bound by ADR-0080, ADR-0088,
  ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by citation and
  receipt. ADR-0178 authorizes each Task commit by its grant, ADR-0182 runs
  Settlement Checks before it, and ADR-0166 records undeclared paths; every
  Task declares its paths. ADR-0229 cites ADR-0167 but decides how an operator
  archive resumes a park, which this Spec does not touch. ADR-0096 and ADR-0097 decide the gate's machine
  stage and row carry, and ADR-0194, ADR-0195 and ADR-0210 decide what a QA row
  records, when it is observed again and its evidence snapshot; the `qa-gate`
  wording this Spec adds restates that carry rule for a report without inputs
  and changes none of them. The adopted branch-prefix entry cites ADR-0191 for clause retention, but ADR-0191 decides how a setup snapshot follows its upstream by name, ADR-0204 and ADR-0206 cite it but decide composed setups and upstream setup names, ADR-0219 cites ADR-0204 but decides which built-in profile a draft adapts (a draft may now omit the branch prefix only because no module requires it), and ADR-0192 cites ADR-0178 but decides how a conflict in declared derived paths is resolved; this Spec changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-10-04, "Autorizar os dois", for the Baseline source under
  `internal/baseline/assets/`, the agent guides and the tests that the source
  change intentionally invalidates; and the standing skills authorization
  "considere autorizado a ajustar todas as skills se necessário". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0223-adjustments-the-adopters-asked-for/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/release.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/templates/guides/agent-instructions.md`,
  `internal/cli/baseline_plan_test.go`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A repository can leave the branch prefix unrecorded, and its guide then
  states the commit-type rule; every repository that recorded a prefix keeps
  its sentence and is asked nothing.
- A QA gate reopened over a report without row inputs has one valid move:
  re-execute the row and declare its inputs for the new pass.
- The release plan reports, without writing anything, whether the skills and
  the Baseline guidance are current, and never changes its decision for it.

## User Stories

1. As an adopting maintainer, I want to leave the branch prefix unrecorded, so
   that my agent guide names branches by the Conventional Commit type rule
   instead of repeating a prefix.
2. As a maintainer who already recorded a prefix, I want my next Baseline
   update to ask nothing and keep my value and its sentence, so that the
   change costs me nothing.
3. As an Agent running a reopened QA gate over an old report, I want the skill
   to tell me to re-execute a row without inputs and declare its inputs for the
   new pass, so that I am not forced to dismiss a contradiction.
4. As a release maintainer, I want `roundfix release plan` to print one
   `skills:` line and one `baseline:` line, in text and JSON, so that I see the
   state of the mandatory skills and guides check before any release mutation.

## Core Features

1. **An optional branch prefix.** The Baseline treats `branch.prefix` as an
   optional decision. The core module no longer requires it; the built-in
   profiles still offer it, and first adoption still asks it. Planning,
   update and profile alignment never report it missing.
2. **Wording that follows the record.** While the decision is unrecorded, the
   agent-instructions guide says that no branch prefix is recorded and that
   new work branches are named `<type>/<description>` from the work's
   Conventional Commit type. A recorded value renders today's sentence byte for
   byte, whatever the value.
3. **Nothing asked of adopters.** An adopter's next Baseline update reports no
   new decision, keeps its recorded value, and leaves its agent-instructions
   guide unchanged. No Baseline clause is removed or renamed.
4. **A reopened gate re-executes rows without inputs.** The `qa-gate` skill
   states that a row from a report without `inputs:` is never carriable and is
   re-executed, and that a row executed again in a new pass declares its inputs
   for that pass, which is not adding inputs after execution. The skill's
   version is raised.
5. **The release plan reports two checks.** A range release plan prints one
   `skills:` line, the comparison the Doctor's `skills:` line makes, and one
   `baseline:` line, what `roundfix baseline update --no-skills` would report.
   The JSON output carries both under one new `checks` field. Each line that is
   not passing names the release runbook's skills and guides check. The checks
   read local files only, write nothing, and never change the plan's state,
   proposed version or exit code. A reset plan and a failed plan are
   unchanged.
6. **The skills and guides say so.** The Roundfix Skill's release reference,
   the release runbook and the Context-Driven Development guide describe the
   behavior above, and each owned skill whose text changes raises its version.

## User Experience

An adopter who never answered the branch prefix sees, after the update, one
paragraph in the agent-instructions guide that names the commit-type rule. An
adopter who answered sees nothing new. `roundfix release plan` gains two lines
after its `Next action:` line, and its JSON gains a `checks` object.

## Non-Goals / Out of Scope

- Clearing a recorded branch prefix through a command; a recorded value is
  kept, as for the frontend layout, ADR-0205: "A recorded value is never
  overwritten."
- Changing the human first-adoption questions or the decision's prompt text.
- Making the release plan run the runbook's test commands, contact any
  service, or refuse a release on a failing check.
- Changing the `### QA settlement` section of any skill, the Daemon's carry
  logic or the QA Report format.
- The archive's relative links, owned by Spec 0226.
- The Grok fallback for `docs` and `chore` Tasks, which the maintainer
  declined on 2026-10-04.

## Success Metrics

1. Success Metric: planning a built-in profile with no `branch.prefix` answer
   reports no missing decision and renders the commit-type paragraph; planning
   with `<type>/` or `ma/` renders the same bytes as before this Spec.
2. Success Metric: `roundfix baseline update` on an adopted repository whose
   Setup Manifest lacks `branch.prefix` does not stop for a decision, and on one
   that records it reports no new decision and leaves its agent-instructions
   guide byte-identical.
3. Success Metric: `roundfix release plan` prints exactly one `skills:` and
   one `baseline:` line and a JSON `checks` object, with the same state,
   proposed version and exit code a plan without them would give, and leaves
   the repository's bytes and Git status unchanged.
4. Success Metric: the `qa-gate` skill states the re-execution rule for a
   report without inputs, its version is raised and recorded, and its
   `### QA settlement` section is byte-identical.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The oraculum maintainer's request of 2026-10-01, captured in the Secondbrain
  Inbox Entry
  `inbox/roundfix/_triaged/2026-10-01-branch-prefix-e-decisao-obrigatoria-e-nao-deixa-o-repo-usar-a-regra-normal.md`,
  with its reproduction on Roundfix 0.23.0 and oraculum's commit `9ff79081`,
  and oraculum's rendered agent-instructions guide in the Secondbrain mirror,
  which carries the prefix sentence this Spec makes optional.
- fluxus's pre-PR review finding of 2026-10-02 on the `qa-gate` skill,
  captured in the Inbox Entry
  `inbox/roundfix/_triaged/2026-10-02-qa-gate-0-0-6-relatorio-antigo-sem-inputs.md`.
- The release runbook's mandatory skills and guides check, added by Spec 0195,
  whose steps the two release plan lines summarize, measured against this
  repository's release history at the tag `v0.33.0`.

## Decisions

- The branch prefix follows the optional-decision precedent rather than a
  `none` value or an empty rendering. See ADR-0228.
- The release plan's checks are reported, never enforced: a failing check
  leaves the decision and exit code alone, because the runbook step, not the
  plan, owns the release's go or no-go.
- The `qa-gate` fix is wording only; the Daemon already refuses to carry a row
  without inputs.
- The archive source moved to Spec 0226: it needs a grant beyond the one named
  for the archive source file, and its delivery-side change would push this
  Spec past four implementation Tasks.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
