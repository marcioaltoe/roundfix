---
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
status: active
created: 2026-11-23
surfaces: [backend, cli, docs]
type: fix
---

# Round 2 of a pre-PR review opens a new session when the configured model changed

Reported by fiscus on 2026-11-23: its Spec 0108 parked in the Delivery Queue as `review-blocked: unclassifiable agent output: no-findings verdict does not account for other content`, although the reviewer's verdict was `No findings.`. Between round 1 and round 2 of the Reviewer Lineage the review profile's model had changed from `gpt-6-luna` to `gpt-6.1-sol`. Round 2 continued the Agent Session round 1 had left open, Codex put a notice about the changed model before the verdict, and the Pre-PR Review Command refuses a no-findings verdict beside other content. The cause is that round 2 prepares the recorded session by its name, on the recorded index of the review profile, and never asks whether that index still names the model the session was opened with: the Pre-PR Review Record does not hold the model. After this Spec each round records the model of its session, and round 2 continues the round-1 session only when the review profile still names that model. Otherwise it ends that session and opens a new one, with the same round-2 prompt.

## Sources

The source is the message of the Fiscus of 2026-11-23, ai-memory message `01a20a55-66b6-7b31-bd03-ed4f3d6a1fbc`, measured on Roundfix 0.153.0 and relayed by the operator's brief. It is copied whole in [`references/sources.md`](references/sources.md), section 1, and cited below as "reported by fiscus". Sections 2 to 4 of that file hold the maintainer's three decisions, sections 5 and 6 the author's characterization and prototype, and section 7 the published sources.

No open Backlog Entry or unresolved Finding shares this context. The seven open Backlog Entries concern a team's status report, the size of a session's context, the name of the guides folder, a `spec check` operand, the semantic analysis deadline, dead probe code and a vacuous Verification; none changes the Pre-PR Review Command or its sessions. `docs/findings/` is empty. The other pending messages of this project's inbox are outside this Spec by the operator's brief and stay pending for Triage.

## Prerequisites

The Task Graph names Spec 0418, `0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks`, in `requires`, so the Delivery Queue starts this Spec only once Spec 0418 has merged. Authored against `50143e74`, which holds v0.153.0 and the plan of Spec 0418.

Spec 0418 is delivered before this Spec, by the maintainer's third decision. It raises the Roundfix skill's version fields and re-records the owned-skill version record and the Behavior Surface Record, and it adds its own decision record, numbered 0450, with that record's line of the Standard Export List. No Verification of this Spec depends on Spec 0418; the `requires` entry holds the order the third decision set, because both Specs change the same files. This Spec's Task regenerates the shared records on the tree it runs on and writes the next version, never a number.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The Spec records an Agent Model name the review profile already holds, and the new ADR takes the number 0451, because 0450 belongs to Spec 0418. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. The change is in which Agent Session name the Pre-PR Review Command prepares and in one field of its local record. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the maintainer's decision of 2026-11-23 changes one recorded rule in part, so a new accepted ADR, numbered 0451, states the new rule, and the ADR it changes stays accepted and gains its amendment notice in the same change (`docs/agents/docs-layout.md`: "In the change that adds the later ADR, its header gains"). That ADR is ADR-0197, whose rule for the session changes (ADR-0197: "round 2 continues it through acpx, on the same selection") and whose other rules stand, the two-round ceiling among them. ADR-0174 applies unchanged: the classifier that refused the answer is not loosened (ADR-0174: "Both verdicts together, or a no-findings verdict beside other content, still block."). ADR-0256 makes the skill change part of this Spec (ADR-0256: "A surface whose fingerprint changed, appeared or disappeared is a Lagging Surface unless a covering skill file changed in the same range"). ADR-0349 is used as it is: the authorization record carries a `patch` Release Ceiling. ADR-0233, ADR-0272, ADR-0286, ADR-0366 and ADR-0391 are listed because each cites the amended decision, and ADR-0196, ADR-0250 and ADR-0412 because each cites a decision named above; this Spec changes none of the eight, and the TechSpec relates each one. No accepted ADR is superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's decision of 2026-11-23, "Sim, antes da 0418 (Recomendado)", for the Spec, its delivery and a patch release; his decision of the same day, "Sim, os 3 arquivos (Recomendado)", for the three files of the Roundfix skill with their regeneration; and his decision of the same day, "Sim, 0418 hoje (Recomendado)", which delivers this Spec after Spec 0418 and names its release `v0.153.2`. The three are quoted with their provenance in the authorization record, which grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. It grants no Makefile, workflow, Go module, lint or formatter change, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`; Spec-contained authorization record: `docs/specs/0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Goals

- A review whose configured model changed between round 1 and round 2 reaches the reviewer's verdict, instead of blocking on a notice the runtime adds when a session continues on another model.
- The Pre-PR Review Record names the model of the session each round used.
- A round 2 whose model did not change continues the round-1 session exactly as it does today.
- The review guide, the Roundfix skill, the glossary and the decision records state the rule.

## Core Features

1. Each review round that sends a prompt records, in the lineage of its Pre-PR Review Record, the Agent Model of the Agent Session that received the prompt.
2. Round 2 continues the Agent Session round 1 left open only when the Agent Selection at the recorded index of the review profile still names the recorded model.
3. When that model changed, or when the round-1 record holds no model because an earlier Roundfix wrote it, round 2 ends the round-1 session and opens a new one through the review profile's normal selection order. The new session receives the round-2 prompt: the delta since round 1 with the round-1 findings and their dispositions. The Reviewer Lineage stays at round 2, the two-round ceiling does not move, and the record reports that the session was not continued.
4. Nothing else changes: the classifier of the reviewer's answer, the new lineage a changed provider starts, the Correction Review, and the fallback to a new session when the recorded one cannot be prepared. No command, flag, output line, exit code or key of Project Config is added.
5. The rule is recorded where the product's rules live. One Task of this Spec writes the new accepted ADR, numbered 0451, and in the same change ADR-0197 gains it in its amendment list with a dated notice. The glossary entry Reviewer Lineage is revised. The review page of the user guide and both copies of the Roundfix skill's review reference state the rule, as the CLI guide requires of a change to a command (`docs/agents/cli.md`: "a change to command behavior updates that skill or reference in the same pull request").

## Non-Goals / Out of Scope

- Removing, ignoring or classifying apart any text of the reviewer's answer. The classifier's rule stays, and the new session covers the reported case (`references/sources.md`, section 6).
- Any other notice an adapter sends as agent text, such as a notice about a compacted conversation. A no-findings verdict beside such a notice still blocks; that is another defect with another cause.
- Comparing the reasoning effort or any other part of the Agent Selection, and continuing a session by setting its recorded model back on it.
- Starting a new Reviewer Lineage at round 1 when the model changed.
- A new output line that says why the session was not continued.
- Any change to a Review Run, to the sessions of a Run, or to the review of a provider that is not an agent.
- Everything of Spec 0418, and every other pending message of this project's inbox.
- The notice to the Fiscus after the release, and its retry of Spec 0108: the operator sends the first, and the Fiscus runs the second.

## Success Metrics

1. Success Metric: through the Pre-PR Review Command over a temporary repository, with a fake runner that puts the notice reported by fiscus before its answer whenever a session name is prepared again with another model, round 2 after the review profile's model changed from `gpt-6-luna` to `gpt-6.1-sol` exits `0` with the outcome `reviewed`, prepares a session name round 1 did not use, ends the round-1 session first, sends the round-2 delta prompt and records `continued: false`. On the unchanged tree the same round 2 exits `2` with the outcome `blocked` and the reason `unclassifiable agent output: no-findings verdict does not account for other content`.
2. Success Metric: a round-1 record whose lineage holds no model makes round 2 open a new session in the same way. On the unchanged tree round 2 prepares the round-1 session name.
3. Success Metric: with the review profile unchanged, round 2 prepares the round-1 session name on the recorded selection, preferred or fallback, and records `continued: true`, as the existing tests of the continued session assert on both trees.
4. Success Metric: the record of every round that sent a prompt holds `lineage.model` with the model of the session that received it, a round 2 on a new session holds the new model, and a `ceiling-closed` record, which sends no prompt, keeps the bytes it has today.
5. Success Metric: the review page of the user guide and both copies of the Roundfix skill's review reference state when round 2 continues the session and when it opens a new one, and name `lineage.model`. ADR-0451 is accepted, ADR-0197 names it in its amendment list with a dated notice, and the glossary entry Reviewer Lineage states the rule. The documentation, skills and glossary contracts pass.
6. Success Metric, outside evidence: the answer the Fiscus measured on Roundfix 0.153.0, quoted whole in `references/sources.md`, section 1, is the answer the fake runner of Success Metric 1 returns for a continued session, with the same two model names; the classifier still refuses that answer, and after this Spec round 2 no longer asks for it.

## Acceptance evidence

The outside-evidence row of Success Metric 6 rests on a repository this Spec did not build and a measurement it did not design. `references/sources.md`, section 1, holds the answer the reviewer returned in the Fiscus's Delivery Queue on 2026-11-23 and the reason Roundfix 0.153.0 parked the item with. The planning change commits that file, so the QA gate reads the answer from it and reads no other repository and no network.

Published literature states the rule the fix follows. A contributor from OpenAI, on the Codex issue 9054 (https://github.com/openai/codex/issues/9054), on 2026-01-12: "You should stick to one model per thread. If you want to use a different model, create a new thread."

## Unreachable Acceptance

None. Every criterion is proved with temporary repositories, a fake runner and the committed reference; none needs the network, an open Pull Request, a real Agent session or another repository.

## Validation

The decisions were checked on 2026-11-23 against the code of the Pre-PR Review Command on `fe247d4c` and `50143e74`, the glossary, the ADRs `## Project Constraints` names, the review page of the user guide and the Roundfix skill's review reference.

- The code confirms the source. Round 2 prepares the recorded session name on the recorded index and compares no model, and the lineage holds no model. A probe on the unchanged tree reproduces the blocked round 2 the Fiscus reported, and a prototype of three changed lines turns it into a reviewed round 2 on a new session (`references/sources.md`, sections 5 and 6).
- The provider half of the report already holds: a changed provider starts a new lineage and ends the open session, and every selection of the review profile is held to the provider's runtime.
- Secondbrain, read only: `wiki/index.md` was read in part by the Research Subagent, and the Indexus tool `pesquisar_acervos` answered two queries, on continuing a reviewer session with another model and on the review lineage. The results were the captures of 2026-09-08 about the reviewer's selection and this repository's own mirrored user guide. Nothing covers a Codex session continued with another model, and no Secondbrain page informed the design.
- Exa: the two published sources of `references/sources.md`, section 7, read by the author. The Codex issue supports opening a new session over continuing one on another model. The Codex commit shows that the notice compares only the model, which is why round 2 compares the model and not the whole Agent Selection. One source challenges the design: the Agent Client Protocol lets a client set a session's model, so the recorded model could be set back on the continued session. That was rejected, because the review would then run on a model the operator no longer configures.

The Research Subagent ran the document, Secondbrain-index and web sweeps in the first authoring round; the author ran the Indexus search and every command. The second round reused that research and repeated none of it.

## Glossary

- changes: **Reviewer Lineage**

## Decisions

- 2026-11-23: Does a short corrective Spec go ahead for the review that continues the reviewer's old session after the configured model changed, and when? — "Sim, antes da 0418 (Recomendado)" (ai-memory message `01a20a6c-7dce-7bb3-9a50-b09b27461492`, from the Praetor, 18:48 BRT).
- 2026-11-23: May the Spec change the three protected files of the Roundfix skill, with their regeneration? — "Sim, os 3 arquivos (Recomendado)" (ai-memory message `01a20ac2-3453-71c1-8227-b02e7cb2c905`, from the Praetor, 20:33 BRT).
- 2026-11-23: Is the order reversed, with Spec 0418 started today and this Spec delivered after it? — "Sim, 0418 hoje (Recomendado)" (ai-memory message `01a20ace-16f5-7c32-b81d-7e73d0340f02`, from the Praetor, 20:46 BRT).

The first decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a5a-ebf0-7ba2-8ac0-95a330f45eeb
pergunta: Roundfix: a fila do Fiscus está parada na Spec 0108 porque, na segunda rodada da review antes do PR, o Roundfix retoma a sessão antiga do revisor (Codex) mesmo com o modelo configurado diferente (gpt-6-luna → gpt-6.1-sol); o Codex avisa a troca e o classificador recusa a resposta, mesmo "No findings". Abrir uma Spec curta de correção? Opções: antes da 0418; depois da 0418; não, o Fiscus contorna.
resposta: "Sim, antes da 0418 (Recomendado)"
alcance: repositório roundfix, Spec nova 0419: PRD, TechSpec, Tasks, implementação, PR, merge e release patch v0.153.1 (a 0418 passa a v0.153.2). Código em internal/cli/review*.go e testes; sem Makefile, CI, go.mod, lint nem skill protegida — se precisar de caminho protegido, novo PEDIDO DE DECISÃO. Avise o Fiscus (default/gesttione-solutions-fiscus) quando o release sair, para ele dar o retry da 0108.
quando: 23/11/2026 18:48 BRT
```

The second decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a8f-f044-7402-ab98-56f7e59b3420
pergunta: Roundfix, Spec 0419 (a revisão abre sessão nova quando o modelo muda): a skill descreve a retomada da sessão e a regra do repositório exige mudá-la no mesmo PR. Três arquivos são protegidos (references/review.md e as duas cópias do SKILL.md). Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0419, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/references/review.md, .agents/skills/roundfix/SKILL.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 20:33 BRT
```

The third decision, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a93-6180-76a2-b59f-20ee9895af31
pergunta: O plano da 0418 está pronto e pode rodar esta noite (6 a 10 h); a 0419 acabou de receber a autorização da skill (01a20ac2-3453) mas ainda precisa ser planejada, e o Fiscus sai com um segundo retry. Inverter a ordem e começar a 0418 hoje?
resposta: "Sim, 0418 hoje (Recomendado)"
alcance: repositório roundfix: merge do PR #946 (plano da 0418) e roundfix deliver start da 0418 hoje; a 0418 sai como v0.153.1 e a 0419 como v0.153.2, entrando na fila logo que planejada. Nada mais muda nas autorizações das duas Specs (0418: 01a20a59-a318; 0419: 01a20a6c-7dce e 01a20ac2-3453).
quando: 23/11/2026 20:46 BRT
```

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0419 of 2026-11-23 and its second-round addendum, the maintainer's three decisions above, the report of the Fiscus, and the code of the Pre-PR Review Command on `50143e74`.

Settled by a source, with no question left:

- The route. The Spec is a bug fix, so it takes `write-techspec` and `write-tasks` with this minimal PRD. A fix is never judged for the Idea stage, so this PRD carries no Idea-stage judgment. Source: the Spec routing guide.
- The correction: round 2 opens a new session when the configured model changed, and the classifier is not loosened. Source: the first decision, whose question names the defect, and the operator's brief, which asks for the new session and allows a change to the classifier only if the new session cannot cover the case. It covers it (`references/sources.md`, section 6).
- The Release Ceiling, `patch`, the three Governed Paths with their regeneration, the delivery after Spec 0418 and the release `v0.153.2`. Source: the three decisions.
- No Rule Challenge. One accepted rule forbids a part of the goal (ADR-0197: "round 2 continues it through acpx, on the same selection"), and the first decision already decides the change. Core Feature 5 names it where the Spec records the change.
- When ADR-0451 is written: in the Task that changes the behavior, not in the planning change. Source: the docs-layout guide, which writes the amendment notice in the change that adds the later ADR.
- The size: one context, one implementation Task and the QA gate. Source: the operator's brief.

The Delegated Decisions below are the four questions no source settles. The author returned each with its options in the first round; the Supervisor session read none as reserved and adopted the recommended option under the project's Authoring Question Preference.

- delegated: 2026-11-23, techspec
  - question: What does round 2 compare before it continues the round-1 session?
  - options: the model of the Agent Selection only; the model and the reasoning effort; the whole Agent Selection, recorded as an object in the lineage
  - adopted: the model of the Agent Selection only — the measured defect is the model, the Codex notice compares only the model, a changed provider already starts a new lineage, and no defect is recorded for a changed reasoning effort
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: A round-1 record written before this fix holds no model. Does round 2 continue its session or open a new one?
  - options: open a new session; continue the session as today
  - adopted: open a new session — the model cannot be proven the same, and the cost is at most one lost continuation for each lineage that is open when Roundfix is updated
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: Which prompt does the new session of round 2 receive?
  - options: the round-2 delta prompt with the round-1 findings and their dispositions; the full round-1 prompt
  - adopted: the round-2 delta prompt with the round-1 findings and their dispositions — it is what round 2 already sends when the recorded session cannot be prepared, and a full prompt would start the lineage again and reset the two-round ceiling
  - under: Project Config
  - answer: none yet

- delegated: 2026-11-23, techspec
  - question: Does the review print why it did not continue the round-1 session?
  - options: no new output line, because the record's continued and lineage.model fields tell; one line on standard error that names both models
  - adopted: no new output line, because the record's continued and lineage.model fields tell — the fallback after a failed preparation prints nothing either, and a fix released as a patch adds no output line
  - under: Project Config
  - answer: none yet

## Open Questions

None. What the real Codex adapter answers on a new session is not exercised by this Spec's tests, which use a fake runner; `_techspec.md` → Risks & Considerations says what supports it and what confirms it after the release.
