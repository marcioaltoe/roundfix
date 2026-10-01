---
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: archived
created: 2026-09-30
surfaces: [backend, cli, docs]
archived: "2026-10-01"
source_slug: 0190-a-task-settles-on-the-facts-its-gate-will-check
---


# A Task settles on the facts its gate will check

A Task settles `completed` when its own declared Verification passes. The QA
gate then refuses, before it spends an Agent turn, on three machine facts: the
strict Spec Consistency Check, the configured repository Verification and the
authorization audit of each Task commit. A completed Task can carry any of
them, and the Run learns it only at the end, after the Agent Session that
could have fixed it has closed.

- On 2026-09-30, Spec 0187's task_04 settled `completed` after one second of
  declared Verification. Two hours into the Run the gate refused at its
  precondition: `make verify-changed` failed on a contract test that task_04's
  own skill edit had broken. The Run ended Unresolved and needed a corrective
  Task and a second Run.
- On 2026-09-29, Spec 0181's second gate refused the same way, on the same
  contract test, after its task_04.
- On 2026-09-29, Spec 0181's third gate refused task_07's commit in the
  authorization audit.

Each refusal cost a full rerun of 15 to 60 minutes. This Spec makes the Daemon
check those facts when the Task settles, so the failure reaches the session
that caused it.

## Project Constraints

- Identifier strategy: not applicable — no new entity identifier. The new
  glossary term, the config key and the check labels follow the existing
  naming style. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes, Git, the Spec
  tree and the Run Database only; no credential is read and no network call is
  added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0182 (this Spec) makes a Task of a
  gated Task Graph settle on the facts its gate will check. ADR-0014 keeps the
  Daemon as the Verification authority, and the new stage is the Daemon's own.
  ADR-0038 keeps the single Verification repair, which now also serves a
  failed Settlement Check. ADR-0096 keeps the gate proving its machine facts
  before an Agent turn, and ADR-0117 places a check with the stage that can
  produce the defect. ADR-0160 keeps the frozen authorization as the only
  authority that opens a red repository gate, and ADR-0159 keeps Verification
  stopping at the first failure unless the Task declares `independent`; the
  appended repository command follows both. ADR-0148 keeps one Verification prober
  that refuses vacuous Tasks before the Agent, and it still probes only the
  declared commands. ADR-0056 keeps Verification Capacity
  separate from Task Capacity, and the appended command uses it. ADR-0111
  keeps an unobserved Verification `unknown`, and ADR-0135 reports an absent
  diagnostic as a state; a Settlement Check always writes its diagnostics.
  ADR-0130 and ADR-0178 keep the authorization audit's rules, and ADR-0179
  lets an explicitly empty paths list grant operations only; the settlement
  audit applies all three unchanged.
  ADR-0166 records a Task's undeclared paths, and ADR-0057 keeps Task status
  Daemon-written. ADR-0093 and ADR-0094 check Spec consistency by citation and
  artifact presence, and ADR-0176 reads citations only from authored text.
  This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167. All of these hold.
  ADR-0020 cites ADR-0014 but ranks a parsed prompt result over the acpx exit
  code. ADR-0127 cites ADR-0014 but reports process residue as a readiness
  fact. ADR-0097 cites ADR-0080 but carries a QA row forward. ADR-0168 cites
  ADR-0093 but narrows the related-ADR check. ADR-0170 cites ADR-0057 but
  treats a Task already completed on its target as nothing to carry. ADR-0183
  cites ADR-0093 but requires a receipt on an attributed claim, and ADR-0184
  cites ADR-0156 but states a command surface as a transcript; both bind a
  Spec authored from their guide commit on, which follows this Spec. This Spec
  changes none of these seven, so they do not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized adjusting every
  skill on 2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and approved the three-wave program ("Três ondas"). The
  authorization is recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`, `make baseline-digests`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A gate fact that a Task's own change causes reaches the Agent Session that
  caused it, as Verification Feedback, before the Task settles.
- A Task that cannot clear such a fact in its one repair settles `failed`
  there, naming the check, instead of letting the Run continue to a gate that
  will refuse.
- The added cost is bounded and visible: one more run of the repository
  Verification per Task settlement, with a switch to turn it off.
- A Task Graph without a QA gate, and the QA gate itself, behave as before.

## User Stories

1. As a Supervisor, I want a Task that breaks the repository Verification to
   be repaired or failed in its own session, so that the Run does not spend
   the rest of the graph and end Unresolved at the gate.
2. As a Supervisor, I want a Task that changes a Governed Path outside its
   grant to hear it before it settles, so that the Agent can undo the change
   while it still knows why it made it.
3. As an operator whose repository Verification is slow, or red for a reason
   of the environment, I want to turn the settlement run of it off, so that
   my Runs keep today's behavior.
4. As a Task author, I want the authoring guidance to say that every Task
   leaves the repository Verification green, so that I slice graphs the
   Daemon can settle.

## Core Features

1. **Settlement Checks run for every Task of a gated graph.** In a Task Graph
   that has a QA gate Task, the Daemon runs Settlement Checks for every other
   Task, after the Task's declared Verification and before the Task settles. A
   Task Graph without a QA gate Task settles on declared Verification alone,
   as today.
2. **The repository Verification runs at settlement.** The configured
   repository Verification runs as the last Verification command of each
   attempt. A Task that already declares that exact command runs it once, as
   declared. The config key `verification.repository_at_settlement`, `true`
   by default, turns this one check off.
3. **Spec Consistency findings the Task introduced are refused.** The Daemon
   takes the refusing findings of the gate's own precondition check when the
   Task starts and again after the Agent's work. A finding that is new since
   the start fails the Settlement Check. A finding that was already present
   does not.
4. **The prospective Task commit is audited.** When the Task changed a
   Governed Path, the Daemon applies the gate's authorization audit to the
   commit it is about to create, with the same rules and the same grant the
   gate would read for that commit.
5. **A failed Settlement Check is a Verification failure.** Its diagnostics
   return to the same Agent Session as Verification Feedback, under the
   existing single repair. A failure on the final attempt settles the Task
   `failed` with a reason that names the check, and no Task commit is created.
   A check the Daemon cannot evaluate fails the same way and never passes
   silently.
6. **The stage is visible and documented.** Each Settlement Check appears in
   the Run Event Stream as Verification events. The user guide, the Roundfix
   skill and the glossary describe the stage and its switch, and the
   `write-tasks` skill states that every Task leaves the repository
   Verification green.

## User Experience

An operator following a Run sees, for each Task of a gated graph, the declared
Verification commands, then the repository Verification command, then two
check lines named `settlement check: spec consistency` and
`settlement check: authorization`. A failure returns the Task to the Agent
once, with the existing `Verification Feedback` line. A Task that still fails
prints `Task <id> failed:` with the check's name and its diagnostics path.

## Non-Goals / Out of Scope

- Changing the QA gate's precondition, its mechanical stage or any verdict
  rule. The gate keeps every check it runs today.
- A separate check of the integrated tree after a parallel Wave.
- An entry observation of the repository Verification for every Task. The
  rule that opens a red repository gate only for a named repair Task is
  unchanged.
- More repair attempts, or a second Feedback turn for a Settlement Check.
- Settlement Checks in the Settle Command, for the QA Task, or in a graph
  without a QA gate Task.
- The gate's consequent-fix order, report shape and evidence path detectors.
  One Task commit does not cause those facts.
- Editing the baseline guides or the Task prompt. The guides belong to the
  wave's guide Spec.

## Success Metrics

1. A Task whose declared Verification passes, and whose change makes the
   configured repository Verification fail, receives Verification Feedback
   naming that command. After the repair it settles `completed`. Without a
   successful repair it settles `failed`, its reason names the command, and no
   Task commit exists.
2. A Task that introduces a refusing Spec Consistency finding receives
   Feedback naming the finding's code. A refusing finding that was present
   when the Task started is not attributed to it.
3. A Task that changes a Governed Path outside its grant receives Feedback
   naming the path. For the same parent and the same changed paths, the
   settlement audit and the gate's audit of the created commit report the
   same findings.
4. A Task Graph without a QA gate Task, and a gated graph run with
   `verification.repository_at_settlement: false`, run the declared
   Verification commands and no repository Verification at settlement. The
   existing test that pins this for a gate-less graph passes unedited.
5. The pre-work prober never runs the repository Verification that settlement
   appends.

## Recorded limits

- The repository Verification runs only after the declared commands pass. A
  Task whose declared Verification fails on the first attempt reaches the
  repository Verification on the final attempt, and a failure there settles
  the Task `failed` with no further repair.
- Settlement Checks see one Task's tree. Two Tasks of one Wave can each pass
  and still integrate into a red tree. The next settlement, or the gate,
  finds it.
- The Daemon takes no entry observation of the repository Verification. A
  Task that starts on a red repository fails at settlement for a failure it
  did not cause, and its reason names the command.
- A finding is matched between start and settlement by its code and its
  sentence. A pre-existing finding whose sentence changes reads as new.
- The Settle Command creates a Task commit without Settlement Checks. The QA
  gate still checks that commit.

## Decisions

- **Only a gated graph.** The stage moves the gate's own facts earlier, so it
  applies where a gate will run. Measured on a prototype: running the
  repository Verification for every Task broke 20 existing tests, because
  their gate-less fixtures have no such command; conditioning on the gate
  broke none.
- **Append, do not add a loop.** The repository Verification joins the Task's
  Verification as its last command, the path a named repair Task already
  uses. The repair bound, Temporary Verification Failure and Repeated Failure
  apply to it unchanged.
- **One switch, for the costly check.** The repository Verification took 206
  to 223 seconds at the gate in four Runs of 2026-09-29 and 2026-09-30. Fleet
  repositories configure their complete gate as that command, and a command
  that is red for a reason of the environment would otherwise fail every
  Task. The two in-process checks take under a second and have no switch.
- **New findings only.** A graph can be inconsistent until a later Task runs.
  Comparing with the Task's start attributes to a Task only what it
  introduced.
- **One audit.** The settlement audit and the gate audit are one function, so
  they cannot disagree about a commit.
- See ADR-0182.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- Run `run_20260930T102010Z_231eaa2e314ed8d8` of Spec 0187. Its event stream
  shows task_04 settled `completed` at 12:24:52Z and the gate refused at
  12:28:20Z. Its verification log `batch-005-attempt-1.log` shows
  `TestSettlementGuidanceIsOneTable` failing under `make verify-changed`.
- Chromium's presubmit documentation
  (<https://www.chromium.org/developers/how-tos/depottools/presubmit-scripts/>).
  It states that skipping slow tests until commit "can result in problems only
  being caught at the last minute". It also states that presubmit checks "do
  not guarantee invariants", because two changes can each pass and fail
  together. The first supports this Spec's premise and the second its recorded
  limit for a parallel Wave.

## Research basis

No pending Inbox Entry addresses this repository. The Secondbrain was
consulted through `wiki/index.md` and the query
`qmd query "checar fatos do gate quando a task assenta descoberta tardia rerun QA"`,
which returned only mirrored Task files of other projects. The research entry
`inbox/secondbrain/2026-09-30-orquestradores-e-ecossistema-jev-o-que-se-transfere-ao-roundfix.md`
records the mechanism this Spec adopts: gates that travel with the slice
instead of a gate only at the end. Its cited case is a graph with a tail-only
gate that produced 9 feature commits and 52 repair commits. Exa located the
Chromium presubmit documentation and a Google Testing Blog article that calls
presubmit "the most valuable real estate" for stopping obvious errors. Both
support checking at the author's change, and the Chromium page also names the
limit this Spec records.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
