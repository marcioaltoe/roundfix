---
date: 2026-08-04
surface: internal/cli, .agents/skills, docs/workflow
status: done
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# What still needs a Supervisor between a PRD and a merge

The stated destination is a process where the maintainer writes the PRD and
everything after it is implemented autonomously. This records, from one
measured session driving Spec 0064 end to end, every point where a Supervisor
was structurally required — not where one merely helped.

## Measured session shape

One Spec, 2026-08-03 into 2026-08-04. Eight implementation Tasks settled
`completed`, each passing its declared Verification on attempt 1, in 1 h 28 min
of Agent time. Around that: four Supervisor-authored artifacts, three Pull
Requests, one failed QA gate, two corrective Tasks, one relaunch.

The implementation was never the constraint. Everything below is.

## The five structural stops

### 1. The Spec arrives as a PRD and nothing else

Nine of ten active Specs carried `_prd.md` alone. That is the intended state —
a TechSpec is authored when the Spec reaches the front of the queue — but it
means the destination process starts exactly where the current process has the
most Supervisor work: `write-techspec`, then `write-tasks`, then the graph's
own validation.

That authoring is not clerical. This session's TechSpec settled a detection
model by probe: only 1 of 89 accepted ADRs names a repository path, so an
ADR-to-code index cannot work, while the ADR citation graph can — the
depth-one closure over Spec 0056's cited set is exactly two ADRs, one of them
the ADR its F-001 named. That is architecture, and it currently has no
autonomous home between "PRD exists" and "Tasks exist".

### 2. The Task Graph must be committed before `implement` can read it

`roundfix implement` reads the graph from the commit, not the working tree. So
every Spec costs a **planning Pull Request** that must merge before any Task
can run, then a second branch for the work, then a second Pull Request. Two
human-visible merges per Spec, one of which exists purely to make a file
readable to a process running on the same machine.

### 3. A non-`pass` gate stops the whole workflow

The contract is correct — merge only on `pass` — but the recovery is entirely
manual: read the report, diagnose the finding, decide the repair shape, author
the corrective Tasks, amend the manifest, reset the gate Task to `pending`,
relaunch. This session did exactly that, and the diagnosis required knowing
that `StaleGateError` refuses to load a `failed` gate whose dependency closure
grew.

There is no autonomous path from "gate returned a finding" to "corrective Task
exists", even for findings whose required repair the report states explicitly.

### 4. The tooling boundary splits repairs that are one repair

F-001's fix was one change with two halves: split a test assertion, and add the
gate step that runs it. Because an authorized tooling Task may mutate only its
bounded files, that became **two** Tasks — task_10 for the test, task_11 for
`Makefile`. The split is correct under the current rule and the audit that
enforces it, but it consumed both corrective Tasks that
`docs/agents/autonomous-work.md` allows before decomposition must be
re-examined. One finding exhausted the corrective budget.

### 5. Progress is not readable without writing a parser

`roundfix runs list` gives a Run's state; `roundfix attach` gives a human the
Live Run View. Between them there is no compact, non-interactive answer to
"which Task is running and which have settled". Every status check this session
required piping `roundfix events --filter task-status` through a hand-written
Python one-liner to reduce JSONL to eight lines.

For a Supervisor that is an agent rather than a person, that is the difference
between checking progress and building a tool to check progress. It also makes
long Runs opaque: `roundfix events --follow` emits nothing to a pipe until it
exits, so a 1 h 28 min Run reported once, at the end.

## Two blind spots worth naming

- **The consistency gate does not read the documents that route the work.**
  Spec 0064 checks Spec folders. This session began by trusting
  `docs/workflow/spec-queue.md` and a handoff, both of which stated Spec 0064
  had a TechSpec. It did not. The gate being built to catch artifact
  contradictions has no reach into the artifacts that told the Supervisor what
  to do.
- **An approval is not a review.** CodeRabbit returned `APPROVED` on PR #98
  with an empty body and zero comments; the walkthrough comment says *"Review
  skipped due to path filters"* for all fourteen files. `reviewDecision` alone
  would have recorded that as reviewed. ADR-0054 already requires evidence
  rather than presence; this is a live instance of why.

## Against the reference material

The second brain's `wiki/concepts/agent-workflows-e-loop-engineering.md`
defines the minimum loop: goal, context builder, maker, capture, checker, state
file, stop conditions. Roundfix has every element, and its checker is genuinely
independent — the Daemon runs Verification, the gate owns a separate Agent
Session. That is the hard part, and it is done.

What the same page calls **closed looping** — "objetivo claro, passos
definidos, eval em cada etapa e ponto de parada/handoff" — is exactly what
exists per Spec. The gap is one level up: the *queue* is open-looped and
human-gated at every Spec boundary, so the loop that works cannot chain.

`wiki/concepts/agentic-coding-e-unknowns.md` names the tension the destination
process has to resolve. Its **interview** technique — "agente pergunta uma
coisa por vez, priorizando respostas que mudariam arquitetura" — is how
`write-techspec` surfaces *unknown knowns*, the criteria a maintainer
recognizes on sight but never wrote down. Autonomous mode instructs the
Supervisor not to ask. Both are right; a PRD-only input makes the collision
routine rather than occasional, and something has to decide which unknowns are
worth one question and which are derivations.

## Corroboration from adopting repositories

Eight findings dated the same day, written from fluxus and vortex sessions,
were sitting untracked in this repository's `docs/findings/` when this one was
written. They observe the same ceiling from the consumer side, and three of
them sharpen a stop recorded above:

- *An accepted gap has no terminal state, so the autonomous loop cannot close*
  — a vortex session needed six maintainer interventions to close one Pull
  Request, and **four carried no judgement at all**. Its proximate cause is one
  missing Review Issue status meaning "valid finding, deliberately not fixed,
  accepted". Every available terminal status either lies about the finding or
  blocks the loop forever. That is stop 3 above, one level deeper: the loop
  cannot converge by construction, not merely by racing.
- *Fail-fast Verification spends the single repair turn on the first of N
  defects* — a sequential gate stops at the first failure, so the Agent's one
  Verification Feedback turn repairs defect 1 while defect 2 is still unknown.
  A Task with two independent defects cannot settle, whatever the Agent does.
- *Review Runs halt autonomous delivery on unrelated dirty files* — a
  four-Spec delivery stopped at the first review step because eight files
  unrelated to the Pull Request were dirty in the checkout.

Together they say the same thing this session measured from the Supervisor
side: the implementation half is close to autonomous, and the **closing** half
hands control back for reasons that mostly carry no judgement.

## Evidence

- Spec 0064 artifacts and manifest,
  `docs/specs/0064-spec-artifact-consistency-gate/`.
- QA report `qa/qa-report-2026-08-03.md`: verdict `fail`, F-001, 12 of 15 rows
  blocked.
- Pull Requests #98 (planning), #99 (finding), #100 (skill), all 2026-08-03/04.
- Run `run_20260803T233822Z_3ffcad0ced4ba246`: 8 Tasks completed, gate failed.
- `internal/spec/spec.go` `StaleGateError` — the load-time refusal that makes
  the gate reset necessary.
- Second brain: `wiki/concepts/agent-workflows-e-loop-engineering.md`,
  `wiki/concepts/agentic-coding-e-unknowns.md`,
  `wiki/concepts/claude-code-operacao-produtiva.md`.
