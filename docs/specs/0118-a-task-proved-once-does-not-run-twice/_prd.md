---
spec: 0118-a-task-proved-once-does-not-run-twice
status: active
created: 2026-08-27
surfaces: [backend, cli, docs]
---

# A Task proved once does not run twice

An Implement Run commits each settled Task inside its Run Worktree and moves
those commits to the user's branch only when the Run reaches Clean. A Run that
ends with an Unresolved Outcome therefore leaves every Task it completed behind
on its Run Branch, and the checkout still reads `status: pending` for work that
ran, passed its Verification, and was committed. The next Implement Run reads
that checkout and re-executes it — not re-checks it: re-executes it, with a
fresh Agent turn and a fresh Verification cycle. Roundfix already owns the
command that hands proved work back, and that command refuses the one outcome
that produces the problem. This Spec makes the remedy reachable, and makes the
Implement Command stop before spending an Agent turn on work a prior Run
already proved.

## Project Constraints

- Identifier strategy: applicable — Unresolved Outcome, Stopped, Run Branch,
  Run Worktree, Task, Verification, Reconcile Command, Implement Command, and
  Preflight Validation are glossary terms this Spec changes behavior around,
  and a message that invents a synonym for one of them is a defect. Two
  glossary repairs are in scope and are consequences of this work rather than
  additions to it: the glossary carries no term for the carry-forward act that
  Spec 0092 shipped, and its Reconcile Command entry states that `--apply` is
  the command's only mutation switch, which the two disposition flags shipped
  alongside it contradict. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential,
  request, or transport is created or read. The work is local Git evidence, the
  Run Database, and CLI reporting. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0023 decides that a Clean integrated
  Run removes its worktree and branch while every other outcome keeps them as
  the single inspection and settle surface, which is the obligation this Spec is
  most directly accountable to: an Unresolved Run's Run Worktree is that
  surface, and a carry-forward that refuses the state contradicts the reason the
  surface is kept. ADR-0115 decides that disposal of retained Run Branch work is
  a separate named act rather than a widening of the read-only report, and
  rejects a command that silently mutates to unblock itself — which is why this
  Spec keeps carry-forward explicit and has the Implement Command refuse rather
  than carry work forward on its own. ADR-0053 decides that terminal Run
  Worktree reconciliation is explicit and proof-based, leaving ambiguous, dirty,
  and unintegrated work intact, so widening the accepted outcome set may not
  weaken a proof to reach it. ADR-0097 establishes that a carried record names
  the report and head that established it, and this Spec applies the same
  principle to a carried Task. ADR-0024 keeps Run integration porcelain-only,
  which bounds how carried commits may reach the checkout, and ADR-0026 decides
  that a settled Task's commit integrates through a serialized cherry-pick whose
  conflict surfaces a graph defect rather than being auto-resolved — the
  mechanism carry-forward already uses and this Spec must leave unchanged.
  ADR-0104 obliges at least one named acceptance row to rest on evidence from
  outside the Spec's own artifacts and lands that obligation at the gate rather
  than during authoring, which is why this Spec names its outside sources
  without stalling on them. ADR-0141 is not applicable: it governs review Runs,
  which create no Run Worktree and carry no Tasks, while every behavior here
  belongs to spec Runs. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — this Spec edits the Roundfix-owned `roundfix`
  skill, which is protected tooling. Express maintainer authorization: granted
  2026-08-27 in session, recorded at
  `docs/workflow/authorizations/2026-08-27-carry-forward-reaches-an-unresolved-run.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`. Its generated copy under
  `skills/` is rewritten by the declared `make skills-sync` and is sanctioned
  fallout, not a separate target. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. A Task that ran, passed its Verification, and was committed on a Run Branch
   is never executed a second time to reach the same result.
2. The command that hands proved work back accepts the outcome that strands it,
   without relaxing any proof that makes it safe.
3. An Implement Run that would re-execute proved work stops before creating a
   Run, and names the command that recovers it.
4. A refusal never leaves a caller with no next action, so the Preflight refuses
   only when the remedy would succeed.
5. An Agent reading the shipped skill can reach the recovery command, which
   today the skill does not mention.

## User Stories

1. As a Supervisor whose Implement Run ended Unresolved on one failing Task, I
   want the Tasks that already passed to be handed to my checkout, so that the
   next Run works only on what is actually unfinished.
2. As a Supervisor starting the next Implement Run, I want Roundfix to refuse
   before it spends an Agent turn re-running proved work, so that a mistake
   costs me one command rather than an hour of Agent time.
3. As a Supervisor whose prior Run's inputs have since moved, I want implement
   to tell me what it found and continue, so that an unrecoverable Run Branch
   does not block the Spec.
4. As a maintainer reading a refusal, I want it to name the Run, the Tasks, and
   the exact command that clears it, so that I do not have to reconstruct the
   Run's history to act on it.
5. As an Agent following the roundfix skill, I want the recovery command to be
   documented, so that I can reach it without reading the binary's help output.

## Core Features

1. **Carry-forward accepts an Unresolved Run.** The Reconcile Command's
   carry-forward act accepts a Run whose outcome is Unresolved on exactly the
   proofs it already applies to a Stopped one: a passing Verification verdict
   for the Task, exactly one settlement commit for it, declared inputs unmoved
   since settlement, a clean checkout, and a repository-local Specs Root. No
   proof is relaxed, added, reordered, or made bypassable, and a set with one
   refusing Task is still refused whole rather than carried in part.

2. **Every other terminal outcome is refused by name.** A Run whose outcome is
   not Stopped or Unresolved is refused with a reason that names the outcome it
   has and why that outcome is not carriable, rather than with a sentence that
   claims carry-forward accepts one stopped Run. A reader learns which outcomes
   are carriable from the refusal itself.

3. **The Implement Command refuses to re-execute proved work.** Before creating
   a Run, implement inspects the Spec's prior terminal Runs in the current
   repository. When at least one Task would be carried forward — meaning it
   satisfies every proof in Core Feature 1 — implement refuses during Preflight
   Validation, creating no Run, opening no Agent Session, and leaving no Git or
   Run Database side effect. The refusal names the Run, the Tasks it would
   recover, and the exact carry-forward command.

4. **A refusal happens only when the remedy exists.** When prior Runs hold
   completed Tasks but no Task satisfies the proofs — declared inputs moved, the
   Run Worktree is no longer present, or a Task carries other than exactly one
   settlement commit — implement reports what it found and why nothing is
   carriable, then proceeds to create the Run. The Preflight never refuses a Run
   whose only offered alternative would refuse as well.

5. **The refusal names one Run.** Where several prior Runs hold carriable work,
   implement names the Run whose carry-forward would recover the most Tasks, and
   breaks a tie by preferring the most recently created. Carry-forward remains
   one Run per invocation, so the caller is given one command rather than a
   sequence to compose.

6. **The glossary carries the act and stops misdescribing the command.** The
   domain context gains a term for handing a settled Task back, and its
   Reconcile Command entry describes the command's actual mutation switches
   rather than naming `--apply` as the only one.

7. **The shipped skill documents the recovery path.** The roundfix skill's
   reconcile guidance names the disposition flags it omits today and records
   which Run outcomes carry-forward accepts, and its implement guidance
   describes the Preflight refusal and the command that clears it.

## User Experience

Not applicable as a browser surface; the surfaces are the terminal reports of
two commands.

A refused Implement Run prints, before any Run exists, the Run that holds the
proved work, each Task it would recover, and one copyable command. A carry-
forward that refuses prints, per Task, which proof failed and against what
evidence. Both are diagnostics and go to stderr, keeping the requested report on
stdout, as the two commands already do.

An implement that finds stranded but unrecoverable work prints why it is
unrecoverable and then continues, so the Run's own output follows the notice
rather than replacing it.

## Non-Goals / Out of Scope

- Carrying work forward automatically. Carry-forward stays an explicit act; the
  Implement Command refuses and names it rather than performing it.
- Changing which outcome integrates. Integration stays gated on Clean alone.
- Making a Failed Run carriable. A Failed Run means the Run itself broke, so its
  own records are not trustworthy evidence about the Tasks it holds.
- Carrying from more than one Run in a single invocation.
- Relaxing, bypassing, or adding a force switch to any existing carry-forward
  proof.
- Changing what `--discard-superseded` does; this Spec only documents it.
- Changing the Settle Command, which cannot reach a Task that is already
  completed and committed on an abandoned Run Branch.
- Re-creating the removed Run Database, cutting a release, or any other item on
  the current queue.

## Success Metrics

- Across the Run Branches of one Spec delivered in several Runs, no Task
  identifier appears in more than one Run's settlement set. The three preserved
  Run Branches of Spec 0098 record the opposite today: seven Task settlements
  that repeat a Task an earlier Run had already completed and verified.
- An Implement Run started against a Spec with carriable stranded work creates
  no Run, and its exit reports Preflight Validation rather than a Run outcome.
- The number of Agent turns spent delivering a Spec that needs more than one Run
  falls by the number of Tasks its earlier Runs had already proved.
- A reader of the roundfix skill can name every mutation switch the Reconcile
  Command accepts.

## Outside Evidence

Two sources originate outside this Spec's own artifacts and are named here so
the Task Graph can rest an acceptance row on them.

- A `fiscus` session on 2026-08-07/08 delivered one Spec through five terminal
  Runs, all Unresolved, totalling 5h21m — a repository this Spec did not build,
  recorded in the finding adopted under this Spec's references.
- The three Run Branches this repository still carries from the Spec 0098
  delivery on 2026-08-25 settle `task_01, task_06`, then
  `task_01, task_02, task_03, task_04, task_06`, then `task_01` through
  `task_07`. Those Runs were created to deliver Spec 0098, not to measure this
  defect, so the measurement is one this Spec did not design.

## Decisions

- Carry-forward stays an explicit act rather than becoming automatic, and the
  Implement Command refuses rather than mutating to unblock itself. See
  ADR-0115.
- The Preflight refuses instead of warning, because the caller that pays the
  cost of a warning is an autonomous loop that does not read one.
- The refusal is conditional on a carriable Task existing, so that a Run Branch
  whose work cannot be proved never produces a refusal with no next action.
- Widening the accepted outcome set is preferred to a new command, because the
  proofs, refusals, and reporting the act needs already exist and duplicating
  them would create a second surface that can drift.
- No new ADR is minted. Every decision above either applies an accepted record
  or completes the scope ADR-0023 and Spec 0092 already stated, and none of them
  reverses a prior decision.

## Open Questions

- Whether a Failed Run should later become carriable on the same per-Task
  proofs. Until answered, a Failed Run is not carriable, which is the current
  behavior.
- Whether carry-forward should accept several Runs in one invocation once
  Unresolved Runs are carriable. Until answered, it stays one Run per
  invocation, and the Implement Command names the single best one.
