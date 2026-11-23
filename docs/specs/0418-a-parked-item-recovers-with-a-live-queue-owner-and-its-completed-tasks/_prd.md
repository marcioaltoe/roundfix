---
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
status: active
created: 2026-11-23
surfaces: [backend, cli, docs]
type: fix
---

# A parked item recovers with a live Queue Owner and its completed Tasks

Reported by quaestor on 2026-11-23: the Implement Run of its Spec 0008 ended Failed with five Tasks completed and only the QA gate's Task failed, and the Delivery Queue item parked as a `delivery-error`. `roundfix deliver retry` answered `Retried ... -> running`. Six hours later the item was still `running`, with no Run, no live Roundfix process and no warning in `roundfix deliver status`. `roundfix deliver resume` then started a Run that began again at the first Task, and the Run Database still named the Failed Run on the item while another Run was active. Reported by the Praetor on 2026-11-23: `deliver status` shows no Task and no Run of a running item, so its periodic sweep reported that a Spec made no progress while it was advancing. After this Spec a Delivery Retry ends with a Queue Owner proven alive or fails and says why, it hands the completed Tasks of a Run that failed only in its QA gate to the item instead of running them again, the item names its latest Run, and status says when nobody is advancing a queue and which Run and Tasks a running item is at.

## Sources

The sources are copied whole in [`references/sources.md`](references/sources.md), which the planning change commits.

1. Sections 1 and 2: the finding and its evidence, ai-memory messages `01a20714-3a6e` and `01a20714-ed79` of 2026-11-23, both measured on Roundfix 0.153.0. Cited below as "reported by quaestor".
2. Section 3: the Praetor's item of 2026-11-23 about the status of a running item. Cited as "reported by the Praetor".
3. Sections 4 and 5: the maintainer's two decisions of 2026-11-23, quoted in `## Decisions`.

The three observations share one context: what happens to a Delivery Queue item after its Run ends, and what the queue's status says meanwhile. They change one component, the Delivery Queue's retry, owner and status, and one completes the other: the retry that leaves no owner is found only because status is silent. One Spec owns them.

This project's inbox held fifteen pending messages when the author listed it. One is the second decision above. The other fourteen belong to other components or only inform, and stay pending for Triage; among them are a report of 2026-11-23 about the pre-PR review, the decisions for Spec 0417, a cancel of a Spec that holds only an Idea and the amendment list of ADRs named by date, which the maintainer's first decision leaves outside this Spec. `docs/findings/` is empty. Seven Backlog Entries are open and none changes the Delivery Queue's retry, owner or status. The entry of 2026-11-11 that names a false status is about the status of a Run, was not reproduced and waits on its reproduction steps, so it stays apart.

## Prerequisites

None in `requires`. Authored against `fe247d4c` (v0.153.0).

Spec 0417 is reserved for another change and is not in the repository yet. If it changes an owned skill or a user-guide page, the two Specs share the Roundfix skill's version fields, the owned-skill version record and the Behavior Surface Record. Each Task regenerates the shared records on the tree it runs on and writes "the next version", never a number.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The Spec reads and records Run identifiers and process identifiers that already exist, and the new ADR takes the number 0450, the next of the repository's sequence. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. Every command this Spec changes reads the local Run Database, local Git objects and the local process table. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the maintainer's decision of 2026-11-23 changes two recorded rules, so a new accepted ADR, numbered 0450, states the new ones and the two ADRs it changes in part stay accepted and gain their amendment notices in the same change (`docs/agents/docs-layout.md`: "In the change that adds the later ADR, its header gains"). ADR-0279 is the first: it keeps status silent about a queue with no owner (ADR-0279: "A queue with no recorded owner adds no line.") and freezes the rest of the output (ADR-0279: "keeps its exit code, output order and every other line"). ADR-0395 is the second: it rejected the same report (ADR-0395: "Keeping the record of an owner that died names exactly the defect."), and its own rule stays as it is (ADR-0395: "An owner that stops with an item left keeps its record."). ADR-0170 applies unchanged and bounds the retry: the Run an item records never selects what is carried (ADR-0170: "The Run recorded on the item is provenance, not the selector."), and Tasks carried before a refusal stay (ADR-0170: "A later Run's refusal leaves earlier Runs' proved Tasks on the item branch and names them."). ADR-0158 gives the reason a fourth outcome is worth accepting (ADR-0158: "A Run that spent its budget after proving three Tasks has three Tasks worth recovering"), and ADR-0284 keeps deciding which planning amendments a carried Task survives (ADR-0284: "A Task whose named section changed is refused with today's reason."). ADR-0365 and ADR-0304 forbid a schema change (ADR-0365: "a schema change makes older item binaries refuse the shared Run Database"; ADR-0304: "No stage, column or schema is added."), and ADR-0365 keeps the Turn Record for a live owner only (ADR-0365: "Status reads it only for the owner it proves live and only when the PID and identity match, and writes nothing."). ADR-0373 keeps deciding when an owner is gone (ADR-0373: "A proven start-identity mismatch at a recorded owner PID means the owner is gone"). ADR-0028 states what a detached process owes its caller (ADR-0028: "a Run's lifetime must not depend on the caller's"). A parked item still has its way back through a Delivery Retry, and a retry that fails leaves the item parked so the way back stays open (ADR-0398: "A parked item always has a way back through a Delivery Retry"). ADR-0256 makes the skill change part of this Spec (ADR-0256: "A surface whose fingerprint changed, appeared or disappeared is a Lagging Surface unless a covering skill file changed in the same range"). ADR-0030, ADR-0034, ADR-0056, ADR-0098, ADR-0148, ADR-0164, ADR-0311, ADR-0315, ADR-0346, ADR-0349, ADR-0384, ADR-0385, ADR-0406, ADR-0412 and ADR-0435 are listed because each cites an ADR named above; this Spec changes none of them, and the TechSpec relates each one. No accepted ADR is superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's decision of 2026-11-23, "Sim, até release patch (Recomendado)", for the Spec, its delivery and a patch release, and his decision of the same day, "Sim, os 3 arquivos (Recomendado)", for the three files of the Roundfix skill with their regeneration, both quoted with their provenance in the authorization record. The record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. It grants no Makefile, workflow, Go module, lint or formatter change, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Delivery Retry never reports an item as retried while no live Queue Owner holds its queue.
- The completed Tasks of a Run that failed only in its QA gate are not executed again after a Delivery Retry.
- A Delivery Queue item names the latest Run its owner started for it, a Failed Run included.
- A reader of `roundfix deliver status`, a person or a routine, can tell from one read whether someone is advancing the queue.
- The same reader can tell which Run and which Tasks a running item is at.

## User Stories

1. As an operator whose item parked after a Run that failed only in its QA gate, I want the retry to carry the Tasks that Run completed, so that the next Run executes the correction and the gate and not the whole Spec again.
2. As an operator who runs `roundfix deliver retry`, I want the command to end only when a live Queue Owner holds the item, or to fail and say why, so that I never leave believing a queue is running when it is not.
3. As an operator who reads `roundfix deliver status`, I want a line that says nobody is advancing the queue and names the command that fixes it, whatever left the queue without an owner, so that an item does not sit for hours at a stage that looks alive.
4. As a coordinator whose routine reads the status of several queues, I want each running item to show its Run and its Tasks in progress, so that the routine reports a queue as stalled only when it is.
5. As an operator who investigates a parked item, I want the item to name the Run that just failed, so that I open the right Run and not an older one.

## Core Features

1. **The item names its latest Run.** A Delivery Queue item records the Run its owner started for it as soon as the owner can read that Run's identifier, and never later than the Run's end, whatever the outcome. A Run that ends Failed is recorded before the item parks: today such an item keeps the Run it recorded before, or none, until a retry writes one, and during the next Run the record does not move until that Run ends (reported by quaestor: the item still named the Failed Run while another Run was active). The blocker of the park and its Park Class do not change. The record keeps the role ADR-0170 gives it (ADR-0170: "The Run recorded on the item is provenance, not the selector.").
2. **A Run that failed only in its QA gate hands its Tasks to the retry.** A **QA-Only Failed Run** is a terminal Implement Run whose outcome is Failed and whose graph has the authored QA gate's Task as its only Task not completed on the Run Branch, whether that Task settled failed or the Run broke before it settled. A Delivery Retry accepts a QA-Only Failed Run as a source of Task Carry-Forward exactly as it accepts a Run that ended Stopped, Unresolved or Budget Exceeded: with the same proofs for every Task, the whole remaining set carried or refused, every terminal Run of the item's Spec considered newest first, each carried Task stamped with its Run and commit, and one `Carried forward from Run <run-id>: <task>, <task>` line per Run. The QA gate's Task is never carried. This changes a Glossary Definition (`CONTEXT.md`: "and refuses every other terminal outcome"), and the maintainer decided the change (the maintainer, 2026-11-23: "incluindo a mudança de produto de aproveitar as tarefas completed num Run Failed só por reprovação na QA"), so no Rule Challenge is asked.
3. **Every other Failed Run keeps today's behavior.** A Run that ended Failed with any Task other than the QA gate's not completed is not a source: the retry carries nothing from it, as today. The Reconcile Command's carry-forward switch and the Implement Command's Preflight keep the three outcomes they accept. Widening either is outside the maintainer's decision and needs a new one.
4. **A retry ends with a proven live owner, or fails.** `roundfix deliver retry` exits `0` only when, at its end, one Queue Owner holds the item's queue and is proven alive under its recorded start-time identity: the recorded owner the item was handed to, or the owner the retry started. The successful output keeps today's lines and order. When no live owner can be proven, because the owner could not be started, ended before the proof or cannot be identified, the command exits `2`, prints no `Retried` line, and names the reason on standard error. The item is then parked again with the blocker and the retry count it had before the command, so the failed attempt spends none of the queue's retry allowance. Tasks the retry had already carried stay on the Item Branch, as ADR-0170 rules for a refusal (ADR-0170: "A later Run's refusal leaves earlier Runs' proved Tasks on the item branch and names them."), and the failure names them instead of claiming that nothing changed. A Queue Owner that ends after the command returned is the status report's part, in Core Feature 5.
5. **Status says when nobody is advancing a queue.** When at least one item can still advance, its stage being neither `merged` nor `parked`, and the queue has no recorded Queue Owner, `roundfix deliver status` prints one `Owner:` line that says so and names `roundfix deliver resume`. The line of a recorded owner proven gone keeps its text, an owner whose identity cannot be checked keeps its diagnostic on standard error, a proven live owner adds no line, and a queue whose items are all merged or parked adds none. Status stays read-only. This reverses what two accepted ADRs rejected (ADR-0279: "A queue with no recorded owner adds no line."; ADR-0395: "Keeping the record of an owner that died names exactly the defect."). The maintainer decided it (the maintainer, 2026-11-23: "o status mostra dono ausente, tarefa e Run em curso"), so no Rule Challenge is asked. The two cases those ADRs wanted to keep quiet are now reported on purpose: after `roundfix deliver stop` the line is true and names the next step, and a read that falls inside the moment a new owner takes the queue shows it once.
6. **Status shows the Run and the Tasks of a running item.** For every item at `running`, in queue order, status prints one line of its own that names the item and one of three states: its Active Run with the Tasks in progress, its Active Run with no Task in progress, or no Active Run. The line is printed whether or not a live owner is proven, because it is read from the Run Database and not from the Turn Record. The item row keeps its four tab-separated fields and its place, every line status prints today keeps its text and its order relative to the others, the exit code stays `0`, and status gains no JSON form. This adds lines where ADR-0279 froze the output (ADR-0279: "keeps its exit code, output order and every other line"), by the same decision of the maintainer.
7. **Each defect is characterized before it is fixed.** Before any behavior changes, named tests reproduce on the unchanged tree the four defects of the sources: a Failed Run that leaves an older Run on the item, a retry after a QA-Only Failed Run that carries nothing, a status with a running item and no owner that prints no warning, and a status that shows neither Run nor Task. The same characterization states by which path a retried item was left at `running` with no recorded owner on 2026-11-23, which the sources do not establish, and the fix closes that path. Every existing test, golden or fixture the change moves is named with the contract it moves to.
8. **The rules that change are changed with the work.** One Task of this Spec writes the new accepted ADR, numbered 0450, and no planning change does. The ADR states that status reports a queue no live owner advances and the Run of a running item, that a Delivery Retry ends with a proven live owner or fails, and that it accepts a QA-Only Failed Run. In the same change ADR-0279 and ADR-0395 gain the new ADR in their amendment lists and a dated notice after the text that changed. The glossary gains the term and revises the three entries `## Glossary` names, each in a Task whose Verification names it. The ADR carries its Knowledge Mark, and its path joins the Standard Export List when the mark is true.
9. **The guide and the skill say what the product does.** The `deliver` page of the user guide and both copies of the Roundfix skill's `deliver` reference describe the new owner line, the line of a running item, the retry's failure and the QA-Only Failed Run. The sentence of the reference that lists "no recorded owner" among the states that add no `Owner:` line is replaced. The skill's version rises through the record command and the Behavior Surface Record is re-recorded, as the CLI guide requires of a change to a command (`docs/agents/cli.md`: "a change to command behavior updates that skill or reference in the same pull request").

## User Experience

The lines below show the intended shape. The TechSpec fixes every byte in its Surface Transcripts.

Status of the queue quaestor reported, as Roundfix 0.153.0 printed it on 2026-11-23:

```text
0008-livros-e-apuracao-do-icms	running	-	<item worktree>
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: 16587262 tokens from 9 of 9 prompt(s) across 1 Run(s); cost not reported
```

The same queue after this Spec:

```text
0008-livros-e-apuracao-do-icms	running	-	<item worktree>
Owner: Delivery Queue has no recorded owner and 1 item(s) can still advance; run roundfix deliver resume
Run: 0008-livros-e-apuracao-do-icms no Active Run
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: 16587262 tokens from 9 of 9 prompt(s) across 1 Run(s); cost not reported
```

A running item under a live owner:

```text
<slug>	running	-	<item worktree>
Turn: <slug> advancing at running
Run: <slug> <run-id>; Tasks in progress: task_03, task_04
```

Between two Tasks the last line reads `Run: <slug> <run-id>; no Task in progress`.

A retry after a QA-Only Failed Run, with an owner proven alive:

```text
Carried forward from Run <run-id>: task_01, task_02, task_03, task_04, task_05
Retried <slug>: delivery-error: run Implement executor: roundfix implement failed with exit code 1 -> running
Delivery Owner: <owner id>
Console Log: <path>
Status: roundfix deliver status
Stop: roundfix deliver stop
```

A retry that cannot prove a live owner prints nothing that reads as success. Standard error names the reason, the item's stage and blocker as they are again, the Tasks that stay carried on the Item Branch, and `roundfix deliver retry <slug>` as the next step. It exits `2`.

## Non-Goals / Out of Scope

- Accepting any other Failed Run as a source of Task Carry-Forward, or a QA-Only Failed Run through the Reconcile Command or the Implement Command's Preflight.
- Carrying the QA gate's Task, a QA row or a Task that is not completed.
- A new blocker, Park Class or stage for an item whose Run ended Failed.
- A record of why a Queue Owner exited, and any restart of an owner that Roundfix makes on its own: `deliver resume` stays an explicit command.
- Any change to what `roundfix deliver start` and `roundfix deliver resume` print, refuse or return.
- A change to the item row of status, a JSON form of status, and any change to the `Turn:`, `Park:`, `Warning:`, `Review:`, `Limits:` and `Usage:` lines.
- A new column, table or stage in the Run Database, and a new command, flag or exit code.
- The counts of completed Tasks of a delivering Spec, which the status view of `roundfix follow-up --status` already shows.
- Spec 0417, the rule about ADRs, the corrections of the Spec Consistency Check, the amendment list of ADRs named by date and the cancel of a Spec that holds only an Idea, which the maintainer's decision leaves for later requests; the priority of a queue and moving a queue between machines.
- Messages to quaestor or to the Praetor: the operator sends them.

## Success Metrics

1. Success Metric: through the production Delivery Engine over a temporary repository and Run Database, with a fake runner, an item whose Implement Run ends Failed parks with that Run recorded on the item. On the unchanged tree the item keeps the Run it recorded before, or none.
2. Success Metric: on a temporary repository whose item parked after a QA-Only Failed Run with five completed Tasks, `roundfix deliver retry` prints one `Carried forward from Run` line naming the five Tasks, the Item Branch holds them completed and stamped with that Run, the QA gate's Task stays pending and the item re-enters at `running`. On the unchanged tree the retry prints no such line and every Task stays pending. Two controls hold on both trees: a Failed Run with another Task not completed carries nothing, and a Run that ended Unresolved carries as before.
3. Success Metric: when the owner a retry starts cannot be proven alive, by a fake process controller, the command exits `2`, standard output holds no `Retried` line, standard error names the reason, and the Run Database holds the item parked with its earlier blocker and retry count and no recorded owner. On the unchanged tree the same retry prints `Retried <slug>: <blocker> -> running` before it starts the owner and leaves the item at `running`. With a live owner, the successful transcript is byte-identical to the one the unchanged tree prints.
4. Success Metric: with an item at `running` and no recorded owner, `roundfix deliver status` prints the new `Owner:` line with `roundfix deliver resume` and exits `0`. On the unchanged tree it prints no `Owner:` line. A queue whose items are all merged or parked prints none, and a recorded owner proven gone prints the line it prints today, byte for byte.
5. Success Metric: with an item at `running`, status prints its `Run:` line in each of the three states of Core Feature 6, naming the Active Run and the Tasks in progress when there are any, and the item row is byte-identical to the one the unchanged tree prints. On the unchanged tree no line names the Run or a Task.
6. Success Metric: the `deliver` guide and both copies of the Roundfix `deliver` reference describe the two new status lines, the retry's failure and the QA-Only Failed Run, and no copy still lists a queue with no recorded owner among the states that add no `Owner:` line. ADR-0450 is accepted, ADR-0279 and ADR-0395 name it in their amendment lists with a dated notice, and the glossary defines the term and the three revised entries. The documentation, skills and glossary contracts pass.
7. Success Metric, outside evidence: on the state quaestor measured on 2026-11-23, rebuilt in a temporary repository and Run Database from the committed reference, the built binary's `roundfix deliver status` prints the owner line and the `Run:` line of an item with no Active Run, where the output quaestor captured holds neither.

## Acceptance evidence

The outside-evidence row of Success Metric 7 rests on a repository this Spec did not build and a measurement it did not design. `references/sources.md`, section 2, holds the state quaestor measured on Roundfix 0.153.0: the three lines `roundfix deliver status` printed at 03:12 BRT of 2026-11-23 for an item at `running`, with no live Roundfix process and no Active Run at that time, and the item as the Run Database held it after the resume, with `retry_count` 1 and a recorded Run that was still the Failed one. The planning change commits that file, so the QA gate rebuilds the state in a temporary repository and reads no other repository and no network.

Published literature states both rules. On a start that reports success before the process is known to run, the manual page systemd.service(5), under `Type=`: "Note that this means systemctl start command lines for `simple` services will report success even if the service's binary cannot be invoked successfully". On keeping what a failed run already completed, GitHub Docs, "Re-running workflows and jobs": "When you re-run failed jobs or a specific job, GitHub also re-runs the jobs that depend on them. Jobs that aren't included keep their outputs from the previous attempt, which the newly run jobs can use."

## Unreachable Acceptance

None. Every criterion is proved with temporary repositories, a temporary Run Database, a fake runner and a fake process controller; none needs the network, an open Pull Request or another repository.

## Risks

1. **The path of 2026-11-23 may not be reproduced.** The sources show a retried item at `running` with no recorded owner six hours later, and a Turn Record that an owner wrote about a minute before the time the retry is reported at. They do not show how that owner ended or why its record was gone, since an owner that dies keeps its record. If the characterization cannot reproduce the path, the Spec says so, and the retry's proof and the status line stand as the answer.
2. **The owner line appears where it was once unwanted.** After `roundfix deliver stop` and inside the moment a new owner takes a queue, status now prints it. A routine that reads status must treat the line as a fact about that read and not as an alarm by itself.
3. **A carried Task meets an amended plan.** An operator often adds a corrective Task before the retry. Carry-forward already refuses a Task whose Binding Text an amendment changed, and a refusal is whole. The retry then fails as it does today for the other outcomes, and the operator sees why.
4. **The proof lengthens the command.** A retry that waits for a new owner to take the queue returns later than today. The TechSpec bounds the wait with the timeouts the detached start already has.
5. **A routine parses status by position.** New lines appear between the lines it knows. The item row and every existing line keep their bytes, and the release notes name the two new prefixes.

## Validation

The decisions were checked on 2026-11-23 against the code of the Delivery Engine, the `deliver` commands and the carry-forward on `fe247d4c`, the glossary, the ADRs `## Project Constraints` names, the `deliver` guide and the Roundfix skill's `deliver` reference.

- The code confirms the sources and bounds them. An Implement Run that returns an error parks the item before its Run is recorded. A retry already writes the newest Run it carried from, so the record is stale between the park and the retry, and while a Run is active. The retry considers every terminal Run of the Spec and skips a Failed one in silence. It stores the stage and prints `Retried` before it starts the owner, and exits `0` once the owner takes the queue. Status prints no owner line whenever no owner is recorded, and a named test pins that silence.
- Two facts the sources leave open are recorded in `## Open Questions`: which path ends a Run as Failed with its QA gate's Task settled failed, since a failing QA verdict ends a Run Unresolved, and how the item of 2026-11-23 lost its owner record.
- Secondbrain, read only: `wiki/index.md`, then the Indexus tool `pesquisar_acervos`, which answered with wiki results. `concepts/agendas-e-reexecucao-em-motores-de-workflow-durable.md` records, from a guide on scheduled jobs, that a monitor alerts on "a run started but did not finish within the deadline" and judges by the deadline and not by the run in progress; it supports giving a routine the Run and the Tasks instead of leaving it to infer progress from a stage. The search excerpt of `entities/roundfix-historico.md` and of a capture of 2026-09-17 records an earlier Run that ended Failed after its Tasks had committed, where a new Run would repeat them; it shows that a Failed Run with proved Tasks is not new.
- Exa: the two published sources `## Acceptance evidence` quotes, systemd.service(5) at freedesktop.org and "Re-running workflows and jobs" at docs.github.com. The first changed the retry's rule from "the owner was started" to "the owner is proven alive"; the second supports carrying what a failed run completed, and its limit, that dependents run again, matches the QA gate running again.

The author ran the code reads and the two consultations in the session; no sweep went to the Research Subagent in this round, because the research of the first round was reused.

## Rule challenges

None.

## Glossary

- adds: **QA-Only Failed Run**
- changes: **Task Carry-Forward**
- changes: **Delivery Retry**
- changes: **Queue Owner**

## Skills

- unchanged: command reconcile — the Reconcile Command's source gains one field of its carry-forward selection that only a Delivery Retry sets, so its help, flags, output and exit codes do not change and its skill reference stays true (Core Feature 3); task_01 records the Coverage Review

## Decisions

- 2026-11-23: Does one Spec of delivery recovery go ahead, with which scope, delivery and release? — "Sim, até release patch (Recomendado)" (ai-memory message `01a20996-258c-7fa1-af9d-a053fbce2441`, from the Praetor, 15:25 BRT).
- 2026-11-23: May the Spec change the three protected files of the Roundfix skill, with their regeneration? — "Sim, os 3 arquivos (Recomendado)" (ai-memory message `01a20a59-a318-7df1-89e7-a3c9baca4a62`, from the Praetor, 18:42 BRT).

The first decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20994-3bf2-7ed3-af99-543c9f49b0fc (itens a e b)
pergunta: Spec única de recuperação da entrega: o retry sobe dono vivo ou falha; o status mostra dono ausente, tarefa e Run em curso; um Run reprovado só na QA aproveita as tarefas concluídas no retry, e o run_id da Base é atualizado. Até 4 tasks mais a QA. Autoriza? Opções: sim, até release patch; sim, sem aproveitar tarefas; só o plano agora.
resposta: Sim, até release patch (Recomendado)
alcance: no repositório roundfix: uma Spec nova com PRD, TechSpec, tasks, implementação, PR, merge e release patch, incluindo a mudança de produto de aproveitar as tarefas completed num Run Failed só por reprovação na QA. Atende os achados 01a20714-3a6e e 01a20714-ed79 (Quaestor) e o item do Praetor (status sem a tarefa em curso, SEM AVANÇO falso no fluxus/0160). Fora do alcance: 0417, a regra de ADR, as correções do spec check, o amended_by e o cancel. Esses vêm em pedidos separados depois desta Spec.
quando: 23/11/2026 15:25 BRT
```

The second decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a58-06b0-7f02-8adb-995c3ad6ebdf
pergunta: Roundfix, Spec 0418 (a recuperação da entrega autorizada às 15:25): ela muda o que o status e o retry da fila mostram. A skill do Roundfix descreve o comportamento antigo, e o próprio Roundfix recusa o release se a skill não mudar no mesmo PR. Três arquivos dela são protegidos e precisam de autorização. Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código; encolher a Spec.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/deliver.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 18:42 BRT
```

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0418 of 2026-11-23 and its addendum, the maintainer's two decisions above, the reports of quaestor and of the Praetor, and the code of the Delivery Queue on `fe247d4c`.

Settled by a source, with no question left:

- The route. The Spec is a `type: fix` that also changes product behavior in one point, the carry-forward after a Failed Run, so it takes `write-prd`, `write-techspec` and `write-tasks`. A fix is never judged for the Idea stage, so this PRD carries no Idea-stage judgment. Source: the Spec routing guide.
- The carry-forward after a Run that failed only in its QA gate, the retry that ends with a live owner or fails, the status that shows an absent owner, the Task and the Run, and the Run the item records. Source: the first decision, whose question lists all four.
- The Release Ceiling, `patch`, and the three Governed Paths with their regeneration. Source: the two decisions.
- No Rule Challenge. ADR-0279, ADR-0395 and the Glossary Definition of Task Carry-Forward each forbid a part of the goal, and the first decision already decides each change. Core Features 2, 5 and 6 name it where they record the change.
- When ADR-0450 is written: in the Task that changes status, not in the planning change. Source: the Spec routing guide, which has a Task of the same Spec change a rule together with the work that depended on it, and the docs-layout guide, which writes the amendment notice in the change that adds the later ADR.
- The size: one context, and up to four implementation Tasks plus the QA gate. Source: the first decision.

The Delegated Decisions below are the five questions no source settles. The author returned each with its options in the first round; the Supervisor session read none as reserved and adopted the recommended option under the project's Authoring Question Preference.

- delegated: 2026-11-23, prd
  - question: Does every route to Task Carry-Forward accept a Run that failed only in its QA gate, or only a Delivery Retry?
  - options: only a Delivery Retry accepts it and the Reconcile Command and the Implement Command's Preflight keep their three outcomes; every route to Task Carry-Forward accepts it
  - adopted: only a Delivery Retry accepts it and the Reconcile Command and the Implement Command's Preflight keep their three outcomes — the maintainer's decision names the retry of the queue, and the other routes would change a second command and a skill reference the grant does not bound
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, prd
  - question: What makes a Failed Run one that failed only in its QA gate?
  - options: its outcome is Failed and the authored QA gate's Task is the only Task of its graph not completed on the Run Branch, whether that Task settled failed or the Run broke before it settled; only a Run whose QA gate's Task itself settled failed; a Failed Run whose recorded reason names the QA gate
  - adopted: its outcome is Failed and the authored QA gate's Task is the only Task of its graph not completed on the Run Branch, whether that Task settled failed or the Run broke before it settled — it reads only the outcome and the Task statuses the Daemon wrote, parses no reason text, and covers a Run that broke while its gate was running, whose implementation Tasks are as proved as the other's
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, prd
  - question: What does a Delivery Retry do when it cannot prove a live Queue Owner at its end?
  - options: it exits 2 with the reason and no Retried line, parks the item again with the blocker and retry count it had, and names the Tasks it carried, which stay on the Item Branch; it exits 2 and leaves the item at its re-entry stage for the status report to flag; it keeps today's exit code and adds a warning
  - adopted: it exits 2 with the reason and no Retried line, parks the item again with the blocker and retry count it had, and names the Tasks it carried, which stay on the Item Branch — a parked item is the state the operator already knows how to retry, a failed start spends none of the retry allowance, and a refused retry already leaves its carried Tasks on the Item Branch
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, prd
  - question: How does the delivery status show the Run and the Tasks of an item?
  - options: one line of its own for every item at running, in one of three states, while the item row keeps its four tab-separated fields; two more fields at the end of the item row; a JSON form of status
  - adopted: one line of its own for every item at running, in one of three states, while the item row keeps its four tab-separated fields — routines read the four-field row today, every other fact of status already has its own prefixed line, and status has no JSON form to extend
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, prd
  - question: Does a Queue Owner record why it exited?
  - options: no, the status line and the retry's proof answer the defect without a cause; yes, an exit record beside the console log that status prints
  - adopted: no, the status line and the retry's proof answer the defect without a cause — the line names the next step whatever the cause, and an exit record is a new file with its own reader and lifetime that no source asks for
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: When does the Queue Owner record on the item the Run it started?
  - options: when the Implement executor returns, whatever it returned, while status reads the Active Run from the Run Database; at the Run's start, by the owner reading the Run Database on a timer beside its child process; at the Run's start, written by the Implement Command when it creates a Run in an item worktree
  - adopted: when the Implement executor returns, whatever it returned, while status reads the Active Run from the Run Database — the owner reads its child's output only when the child exits, a timer adds a goroutine and a real-time wait to the owner and to its tests, the Implement Command cannot name the queue from an item worktree without a second identity rule, and neither the retry nor status reads the item's record while a Run is active
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: What does the Run line of status show while a Run is in its QA gate?
  - options: the Active Run with no Task in progress, read from the Task events the Run Database holds; a fourth state that names the QA gate, inferred from the gate's own events; a start event the Daemon would journal for the gate's Task
  - adopted: the Active Run with no Task in progress, read from the Task events the Run Database holds — the PRD fixes three states, the Run Database holds no start event for the gate's Task and no event that names it before it settles, and a new event is a change to the Daemon, which is another component
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: What does a Delivery Retry do when the owner it started has already left every item merged or parked?
  - options: it exits 0 with the lines of a successful retry, because nothing can advance and the owner's handshake shows it claimed the queue; it exits 2 and parks the item again with its earlier blocker; it exits 2 and leaves the item as the owner left it
  - adopted: it exits 0 with the lines of a successful retry, because nothing can advance and the owner's handshake shows it claimed the queue — the owner took the item, which is what the retry promises, and parking it again would overwrite the blocker that owner just recorded
  - under: Project Config
  - answer: none yet

## Open Questions

1. Which path ends a Run as Failed with its QA gate's Task settled failed? The glossary says a failing QA verdict ends a Run Unresolved, which a retry already carries from, and quaestor's Run ended Failed. Default until answered: the definition of Core Feature 2 depends on no path, and the TechSpec's characterization names the one it reproduces.
2. How did the item of 2026-11-23 come to stand at `running` with no recorded owner? Default: Core Feature 7 has the TechSpec reproduce it; when it cannot, `## Risks`, item 1, applies.
3. Does the item's record move when its Run starts or when it ends? Default: at the start when the owner can read the Run's identifier without a schema change, and otherwise at the end, with status naming the Active Run in both cases.
4. Where do the Tasks in progress of an Active Run come from? Default: from what the Run Database already holds for the Run, read without writing and without reading another working tree; if it holds no such fact, the line names the Active Run and the TechSpec says what it shows of the Tasks.
5. What are the exact words and places of the two new lines and of the retry's failure? Default: the shapes of `## User Experience`, fixed byte for byte by the TechSpec's Surface Transcripts.
