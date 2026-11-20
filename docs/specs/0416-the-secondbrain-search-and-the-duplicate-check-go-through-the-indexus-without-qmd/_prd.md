---
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
status: active
created: 2026-11-20
surfaces: [cli, docs]
type: feat
---

# The Secondbrain search and the duplicate check go through the Indexus, without qmd

Every repository that adopts the Baseline with the Secondbrain reads a guide that orders two uses of the `qmd` command: the search that follows the wiki index, and the advisory duplicate check before a research digest is sent. The owned skill `write-idea` orders the same check. The maintainer decided on 2026-11-20 that the Indexus, the search service of the knowledge base, replaces qmd with no second path once it is shown superior in real use, and that qmd is then removed from the machines. A guide that still orders a removed command would send every Agent to a tool that is gone. After this Spec the guide and the skill name the Indexus tool and no other search tool, every adopter receives the change through its next Baseline update with every clause accounted, and the change waits as an open Pull Request until the result of the Indexus test arrives.

## Sources

1. [`references/2026-11-20-messages-and-decisions.md`](references/2026-11-20-messages-and-decisions.md), recorded by the author: the request of the project that owns the Secondbrain (ai-memory message `01a1f8a4-d7ed`, 2026-11-20) and its update (`01a1f954-fbf4`, 2026-11-20), the report of the skills project about `write-idea` (`01a1f932-8f30`, 2026-11-20), the maintainer's decision about retiring qmd (`01a1f952-79e3`, 2026-11-20), his decision about the Idea stage (`01a1f955-c546`, 2026-11-20) as the operator's brief records it, and the Idea-stage judgment. It is cited below as "the maintainer, 2026-11-20" and "the Secondbrain's request".
2. [`references/2026-11-20-characterization.md`](references/2026-11-20-characterization.md), recorded by the author: where qmd is named, the two clauses and the skill step as they are, a prototype of the change on a scratch clone, the Governed Paths, what the Indexus answered from this repository's session, and the Secondbrain files and published sources consulted. It is cited as "measured by the author".

The three messages share one context: each asks for the same replacement in the same contract, the Baseline's Secondbrain consultation, and the skill step repeats the guide's clause word for word. One Spec owns them. The update `01a1f954-fbf4` asks to be told when the version that carries the change is published; `## Delivery limits` keeps that request.

This project's inbox held seven pending messages when the author listed it. Two belong to other components and stay pending for Triage: a cancel of a Spec that holds only an Idea, and a required skill an adopter lacked after a Baseline update. `docs/findings/` is empty. Seven Backlog Entries are open: two wait on measurements from outside, one proposes renaming folders, and four record defects of the Spec Consistency Check, a Baseline analysis deadline and unused probes. None changes the Secondbrain consultation, and none is adopted.

## Prerequisites

None for Verification. Authored against `851024fa` (v0.152.0).

Spec 0415 is authored in parallel on another branch, and the author of this PRD did not read it. If it changes a Baseline module or an owned skill, the two Specs share the module version record, the catalog files, the plan goldens, the Standard TypeScript profile's pin, this repository's Setup Manifest and the owned-skill version record. Each Task regenerates the shared records on the tree it runs on and writes "the next version", never a number. The same holds for every change that merges while this Spec's Pull Request waits (`## Risks`, item 3).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The two clauses keep the ids they have, and the new ADR takes the next number of the repository's sequence. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Roundfix adds no endpoint, credential or network call. The guide names a tool of an MCP server that the session's own configuration provides, with a key Roundfix never reads, and no code, test or Task Verification calls it. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0330 states the consultation rule this Spec changes, so it gains an amendment notice and stays accepted (ADR-0330: "The Secondbrain stays the business knowledge base, and it is read-only."); the notice follows the lifecycle of the docs-layout guide (`docs/agents/docs-layout.md`: "An accepted ADR that a later accepted ADR changes in part, without replacing it, stays"). Both clauses keep their identity and their force, because a new identity costs what a rewording does not (ADR-0257: "a new identity needs a Source Baseline row that the regenerator"; ADR-0186: "Every clause keeps its enforcement level when its"), as two earlier rewordings did (ADR-0372: "The clause keeps its identity and force, so a Source Baseline transition"; ADR-0400: "the clause id, force and Source Baseline row"). The new words state only what the shipped product and the named tool do, and cite no Spec and no ADR (ADR-0186: "A Baseline clause now states only behavior the shipped product has."). The guide still reaches only repositories that have the Secondbrain turned on, where the named tool is the knowledge base's own (ADR-0222, whose title is the rule: "A Baseline guide says only what holds for the repository that reads it"). The module's version is chosen by the record step (ADR-0250: "A Baseline module's version names one content, and the record step"). The Vendored Skill that also names qmd is not edited here (ADR-0360: "The Catalog Skill Bump is the only writer."). The new ADR carries its Knowledge Mark (ADR-0439: "knowledge the Secondbrain wants, and only a marked ADR is published."), and its path joins the Standard Export List as the update would write it (ADR-0438: "added in order, which is what the update would write."). A Run's Agent still does not consult the Secondbrain (ADR-0257: "Session cannot reach the Secondbrain."), and the Research Subagent keeps its tool list (ADR-0402: "Sweeps that need a command stay with the author."). Adopters learn of the change at session start (ADR-0371: "It compares the installed Roundfix with the latest"). ADR-0203 is listed to record why it does not apply: no hand-written rule of a repository is promoted, and the evidence from outside Roundfix is in `## Acceptance evidence` (ADR-0203: "The Baseline grows only with evidence from outside Roundfix."). ADR-0256 is listed for the same purpose: no command help, configuration key, exit code or user-guide page changes (ADR-0256: "A Behavior Surface is one user-visible"). ADR-0208, ADR-0248, ADR-0258, ADR-0338, ADR-0342, ADR-0349, ADR-0364, ADR-0380, ADR-0383, ADR-0390, ADR-0396, ADR-0399, ADR-0401, ADR-0411, ADR-0412, ADR-0415, ADR-0416, ADR-0441 and ADR-0444 are listed only to record that they change nothing this Spec builds; the TechSpec relates each one. No accepted ADR is superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's decision of 2026-11-20, "já podemos usar somente o indexus e eliminar o qmd", whose recorded scope names this repository's part of the replacement, "trocar toda referência ao qmd por Indexus, sem reserva: guia docs/agents/secondbrain.md e write-idea do baseline (Roundfix)", and his decision of the same day that lets this Spec reach an open Pull Request before the test's result, both quoted with their provenance in the authorization record. The record grants `implement`, `commit`, `push` and `pull_request`, and neither `merge` nor `release`. It grants no Makefile, workflow, Go module, lint or formatter change, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd/_authorization.md`; bounded files: `.agents/skills/write-idea/SKILL.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`, `internal/baseline/assets/modules/secondbrain.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/write-idea/SKILL.md`.

## Goals

- An Agent that follows the Secondbrain guide searches the knowledge base through one path, the Indexus, in every repository that adopts the Baseline with the Secondbrain.
- Removing qmd from the machines leaves no instruction that Roundfix ships ordering it.
- An adopter's Baseline update accounts every clause the adopter held, and changes only its Secondbrain guide and its Setup Manifest.
- An Agent whose session has no Indexus, or whose Indexus does not answer, knows what to do from the guide alone and reaches for no other search tool.
- The change is ready and reviewed as an open Pull Request, and nothing of it ships before the Indexus is shown superior.

## User Stories

1. As an Agent about to propose a decision or to author an Idea, a PRD or a TechSpec, I want the guide to name the search tool the knowledge base answers by, so that I consult it in one call instead of running a command that is slow, locked or removed.
2. As an Agent that did substantive external research, I want the duplicate check named with the same tool, in the guide and in `write-idea`, so that my digest tells the knowledge base's owner what it already holds.
3. As an Agent whose session does not expose the Indexus, or whose call gets no answer, I want one rule for that case, so that I record the limitation and go on with what I have.
4. As the maintainer of an adopter, I want the change to arrive with my next Baseline update, with every clause I held accounted and nothing else rewritten, so that I can apply it without reading a diff of my whole instruction set.
5. As the maintainer, I want the change prepared now and held at an open Pull Request, so that qmd keeps answering until the Indexus is proven and the switch then costs one merge and one release.
6. As the project that owns the Secondbrain, I want one line when the version that carries the change is published, so that I remove qmd from the machines only after the guides stopped ordering it.

## Core Features

1. **The search step names the Indexus.** The clause that orders the search still starts with reading the wiki index. Its second step becomes a search with the Indexus MCP tool `pesquisar_acervos` of the `indexus` server. The clause says that the answer holds one group for each acervo the session's key reaches, that `abrir_conteudo` opens a result, and that the path of a wiki result is relative to the wiki folder. The rest of the clause keeps its meaning: an Agent inspects a project mirror only when the index and the search point there, opens only the files the task requires, and treats mirrors as references.
2. **The duplicate check names the Indexus.** The clause that sends substantive external research to the knowledge base's owner still orders the advisory duplicate check before the message. The check is run with the same tool. The rest keeps its meaning: the Agent verifies that the returned paths exist, names a strong verified match in the message so the owner extends what exists, and no relevance label or position decides alone.
3. **No second search path.** No clause of the guide names qmd or any other search tool as a reserve (the maintainer, 2026-11-20: "já podemos usar somente o indexus"; the scope recorded with his decision reads "sem reserva"). When the session does not expose the Indexus tool, or the tool does not answer, the guide's existing rule about an unavailable source applies unchanged: the Agent records the attempted consultation and its limitation, and continues with the evidence it has. An Agent never substitutes another search tool and never claims a consultation that did not happen.
4. **Both clauses keep their identity and their force.** Each stays mandatory under the id it has. An adopter's next Baseline update rewrites the two sentences in its Secondbrain guide and accounts its clauses as it did before: `clause.secondbrain.02-query-order` is recorded `retained`, and `clause.secondbrain.research-to-the-knowledge-owner`, which the Source Baseline does not hold, has no disposition of its own and stays the target of the `replaced` record of `clause.secondbrain.research-capture`. Both keep their mandatory force, and no clause is left unaccounted. The Source Baseline keeps the wording it recorded.
5. **`write-idea` orders the same check.** The research step of `write-idea` names the Indexus tool for the advisory duplicate check, in the canonical copy and in its mirror, with the same conditions as the guide. The skill's version rises and is recorded, and the catalog that vendors the skill receives it by its next sync.
6. **This repository reads what it ships.** This repository's own Secondbrain guide and Setup Manifest change through the Managed Refresh, the Standard TypeScript golden and pin and the other derived Baseline files through the sanctioned regeneration, and the module's version through the record step. None is edited by hand.
7. **The decision is recorded as a change of a rule.** ADR-0330 lists the consultation rules the Baseline keeps, and the second of them is qmd (ADR-0330: "The Baseline keeps its consultation rules: read"). The maintainer decided the change (the maintainer, 2026-11-20: "já podemos usar somente o indexus e eliminar o qmd"), so no Rule Challenge is asked. A new accepted ADR states that the search and the duplicate check go through the Indexus with no second path, why, what an Agent does without it, and that the decision takes effect with the release that carries it. Its Knowledge Mark is true, because the Agents of every adopter act on it, so its path joins the Standard Export List. In the same change, ADR-0330 gains the new ADR in its Amendment List and a dated notice after the sentence that changed. One Task makes this change together with the clause edit it describes, so the record and the guide never disagree on the default branch.
8. **What still names qmd is declared.** After the change the word appears in the tree only where it records the past or belongs to another owner: the Source Baseline corpus, two fixtures of the receipt check, one dated measurement under the references, the sentence of ADR-0330 that its notice now follows, the Archive Records, and the Vendored Skill `knowledge-workspace` until a Catalog Skill Bump brings its next text. No Baseline module, no guide the Baseline writes, no formatter golden and no Roundfix-owned skill names it.

## User Experience

An adopter's maintainer sees the change as every Baseline change is seen. At session start the Update Check reports the new release by its title and a Baseline plan for the repository. The update previews two files, the Secondbrain guide and the Setup Manifest, and applies them on confirmation. A repository whose Baseline has the Secondbrain turned off receives no guide text from this change.

An Agent that follows the guide reads the wiki index, calls the search tool with its question, and receives groups. It opens the results it needs, cites each Secondbrain file it used by path, and records in the artifact how the sources informed the decision, as today. An acervo the session's key does not reach has no group in the answer; the guide says so in the same sentence that names the tool.

When the tool is absent or silent, the Agent writes one line in the artifact or the response: what it attempted and why it got nothing. It then continues. This is the state the guide already describes for either required source.

## Non-Goals / Out of Scope

- Removing qmd from any machine, stopping its services or deleting its indexes. The project that owns the Secondbrain coordinates that, after this Spec's release.
- The side-by-side test of the Indexus against qmd, its questions and its numbers. The Indexus project owns it.
- The Vendored Skill `knowledge-workspace`, which also names qmd. Its text belongs to the skills catalog, and a Catalog Skill Bump brings it; this Spec runs none.
- The global instructions, the Secondbrain's own contract and skills, and the jobs of other projects that name qmd.
- Any other clause of the Secondbrain module, the first one included, which calls the Secondbrain "local". The read-only rule, the citation rules, the export rules and the Exa requirement keep their words.
- Giving the Research Subagent an Indexus tool, or changing any Project Config key. The Indexus search stays with the session that runs the phase, as the qmd search did.
- Any Roundfix code that calls the Indexus, checks that its MCP server is configured, or reports its absence in the Doctor Command.
- Naming the acervos, their levels, the client mark or the search filters in the guide. The guide names the tool and the two facts an Agent needs to read its answer.
- Rewriting the Source Baseline, the receipt fixtures, the dated measurement or any Archive Record.
- The merge and the release of this change, and the changelog section of the release (`## Delivery limits`).

## Delivery limits

The maintainer's decision makes the replacement conditional: "Se for, já podemos usar somente o indexus". The request's update repeats it: "Preparem já, mas publiquem só depois do resultado do teste do Indexus".

- The authorization record grants `implement`, `commit`, `push` and `pull_request`. It grants no `merge` and no `release`, so the Delivery Queue does not deliver this Spec, and its implementation stops at an open Pull Request.
- The Pull Request is not merged before the notice, so that no release of another Spec ships this change early.
- When the Secondbrain's project sends the result of the test and the Indexus is superior, the operator records the notice, widens the authorization record in a change of its own, merges and releases. When the Indexus loses in some use, the Pull Request waits until the Indexus project corrects it, and nothing of qmd is removed.
- The section of the changelog that the release publishes tells an adopter, in its title, that the Secondbrain search needs the Indexus MCP server in the session and that the guide no longer names qmd. The operator writes it at release.
- Once the version is published, the session that follows that release sends one line to the project that owns the Secondbrain, naming the version. That project asked for it and waits for it before it removes qmd. No Task sends a message.

## Risks

1. **A project session may not reach the wiki through the Indexus.** Measured by the author on 2026-11-20 from this repository's session: no answer held the personal acervo, and two searches restricted to wiki pages returned nothing, one of them the single word "Indexus". The Indexus's own documents say the whole wiki sits in the personal acervo and that a key below the personal level does not reach it. With qmd an Agent found wiki pages; with the Indexus alone, an Agent of such a session finds mirrors and raw sources and relies on the wiki index for the rest. Whether that is a use in which the Indexus loses is a question for the test that holds the release (`## Open Questions`, item 1).
2. **A session without the Indexus consults less than it did.** qmd answered from local files with no server. A session whose harness does not hold the Indexus server, or whose client holds a stale tool description, now records a limitation where it used to search. The maintainer chose one path, and the published sources recorded in the characterization say why a reserve path that is rarely exercised is itself a cost. The limitation is visible in every artifact that records it.
3. **The waiting Pull Request goes stale.** Every Baseline module or owned-skill change that merges meanwhile moves the shared derived records. Before the merge the operator brings the branch to the new base and runs the sanctioned regeneration again; no derived file is merged by hand.
4. **The guide and a Vendored Skill may disagree for a while.** Until a Catalog Skill Bump brings the catalog's next text, an adopter can read a guide that names only the Indexus and a `knowledge-workspace` skill that still shows qmd commands. The guide is the mandatory rule. The operator decides at release whether a bump goes first (`## Open Questions`, item 2).
5. **The guide names another project's tool.** If the Indexus renames the tool or changes the shape of its answer, the guide goes stale and no test of this repository can see it, because no test reaches the network. The contract is versioned by its owner, which tells this project of a change by message.
6. **The guide grows.** The trial wording added about 170 characters to a guide every authoring session reads. The TechSpec keeps the two sentences as short as the facts allow.

## Success Metrics

1. Success Metric: the Secondbrain guide this repository reads and the Standard TypeScript golden of that guide each name `pesquisar_acervos`, `abrir_conteudo` and the `indexus` server in the search clause and `pesquisar_acervos` in the duplicate-check clause, and hold no mention of qmd; a named test pins those phrases against the module, fails on the unchanged tree and fails again when any of them is removed.
2. Success Metric: on a temporary repository that holds the Baseline of the Delivery Base with the Secondbrain turned on, the built binary's `roundfix baseline update` previews a plan that changes the Secondbrain guide and the Setup Manifest and no other file, records `clause.secondbrain.02-query-order` as `retained` and `clause.secondbrain.research-capture` as `replaced` with the target `clause.secondbrain.research-to-the-knowledge-owner`, each with mandatory force, and leaves no clause unaccounted; after it is applied, a second update plans no change.
3. Success Metric: both copies of `write-idea` name the Indexus tool in the research step, hold no mention of qmd and are byte-identical, and the skill's version is one the owned-skill version record holds for those bytes; a named test fails on the unchanged tree.
4. Success Metric, outside evidence: the tool names and the shape of the answer that the guide states are the ones the Indexus publishes in its contract `indexus-acervos/1`, as the characterization records them from the contract's revision of 2026-11-18 and from the author's live calls of 2026-11-20.
5. Success Metric: an accepted ADR whose Knowledge Mark is true states the decision, ADR-0330 names it in its Amendment List and carries the dated notice after the sentence that names qmd, and the Standard Export List is the one the Baseline update would write.
6. Success Metric: a whole-word search of the tree finds qmd only in the places Core Feature 8 declares, and a named test that holds that list fails when a Baseline module, a guide the Baseline writes, a formatter golden or a Roundfix-owned skill names it.

## Acceptance evidence

The outside-evidence row of Success Metric 4 rests on a repository this Spec did not build and on a measurement it did not design. The Indexus's contract, read in the Secondbrain's mirror of that project, says: "Há um grupo por acervo alcançado e pedido." The author's calls from this repository's session returned the contract name `indexus-acervos/1` and one group per reached acervo, and `abrir_conteudo` returned a page. Both are copied in `references/2026-11-20-characterization.md`, section 5, which the planning change commits, so the QA gate compares the guide's words with that record and reaches no network and no other repository.

Published literature states the reason for a single path: Jacob Gabrielson, "Avoiding fallback in distributed systems", Amazon Builders' Library: "One of the worst things about fallback is that it isn't exercised regularly and is likely to fail or increase the scope of impact when it triggers during an outage."

## Unreachable Acceptance

- criterion: The Indexus is superior to qmd in what the projects and the maintainer's assistant use.
  reason: The Indexus project measures it in a side-by-side test outside this repository, and the QA gate reads no other repository and reaches no network.
  satisfied-by: The notice of the test's result that the project that owns the Secondbrain sends to this project's inbox, which the operator records before the merge.
- criterion: An Agent session that exposes the Indexus follows the reworded clauses and receives groups from the live service.
  reason: The QA gate starts no Agent session and calls no MCP server.
  satisfied-by: The author's live calls of 2026-11-20, recorded in the characterization, section 5.

## Validation

The decisions were checked on 2026-11-20 against the Secondbrain module, the three guides and the skill that render or repeat its two clauses, the plan code that classifies a clause transition, the tests that pin the two clauses' force, ADR-0330 and the ADRs `## Project Constraints` names, and a prototype of the change on a scratch clone whose sixteen changed files and seven Governed Paths `references/2026-11-20-characterization.md` records in sections 3 and 4. The Secondbrain was consulted through its wiki index, the Indexus page and three documents of the Indexus mirror, and through the Indexus itself; the qmd query the guide still requires was attempted and ended on a locked database. Exa gave two reliability sources on fallback paths, the MCP specification and client documentation on tools that vary by session, and the qmd README. Section 6 of the characterization lists each file and URL and what it changed here. The Research Subagent swept the codebase and the web; the author ran the calls that need a command or an MCP tool.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: no (0.13)
- unverified_assumption: unsure (0.34)
- market_question: no (0.10)
- verdict: undecided
- stage: dismissed — the maintainer chose "Pular a Ideia" on 2026-11-20; the one criterion left unsure, that the Indexus is superior to qmd, is answered by the side-by-side test of the Indexus project, which holds the release and not the authoring
- shape: not applicable

## Rule challenges

None.

## Glossary

None.

## Decisions

- 2026-11-20: What is the criterion for retiring qmd, and what happens when it is met? — "Agora, os testes devem ser com o objetivo inicial de ser superiores ao que o qmd/secondbrain entrega no que os projetos/jarvis usam. Se for, já podemos usar somente o indexus e eliminar o qmd no mini e macbook." The maintainer's assistant recorded the scope with it: the replacement of every reference, "sem reserva", in the Baseline's Secondbrain guide and in `write-idea`, authorized with no new request once the Indexus is superior (ai-memory message `01a1f952-79e3`).
- 2026-11-20: Does the Idea stage run for this Spec, whose judgment was undecided? — "Pular a Ideia": the Spec goes to the PRD, the TechSpec, the Tasks and an implementation in a Pull Request, and the release waits for the notice of the Secondbrain's project (ai-memory message `01a1f955-c546`, as the operator's brief records it).

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0416 of 2026-11-20, the maintainer's two decisions above, the Secondbrain's request and its update, the report of the skills project and the characterization.

No Open Authoring Question was left: each question below is settled by a source, named with it, and none is a Delegated Decision.

- Does the guide keep a reserve search tool? No. Source: the maintainer's decision, whose recorded scope reads "sem reserva", and the request's update, which drops "qmd as the fallback".
- What does an Agent do without the Indexus? What the guide already orders for an unavailable source. Source: the unchanged limitation clause of the same rule, copied in the characterization, section 2.
- Does the search still start with the wiki index? Yes. Source: the Secondbrain's request, whose proposed sentence keeps it.
- What does the guide say about the session's key? One fact, that the answer holds one group per acervo the key reaches, and no new duty to list the acervos in each record. Source: the Secondbrain's request, whose proposed sentence says exactly that; the limitation clause already covers a source that yields nothing.
- Do the clauses take new ids? No. Source: ADR-0257 and the repository rule that a reworded clause keeps its identity.
- When is the decision record written? In the Task that changes the clauses. Source: the Spec routing guide, which has a Task of the same Spec change a rule together with the work that depended on it, and the docs-layout guide, which writes the amendment notice in the change that adds the later ADR. A record that arrived with the plan would state a rule the shipped guide does not yet follow, for as long as the Pull Request waits.
- Does this Spec bring the Vendored Skill's new text? No. Source: the maintainer's decision, which gives `knowledge-workspace` to the skills project, and ADR-0360.
- Which operations does the authorization grant? Everything up to an open Pull Request. Source: the operator's brief, which forbids `release` and holds the merge.

## Open Questions

1. Does a project session reach the wiki through the Indexus? The author's session did not (`## Risks`, item 1). The Indexus project and the maintainer decide whether keys, the placement of the wiki or the guide's reliance on the wiki index answer it. Default until answered: this Spec goes on unchanged, the session that started this phase sends the measurement to the Indexus project as input to its test, and the test's result holds the merge.
2. Does a Catalog Skill Bump that brings `knowledge-workspace` without qmd go before the release? Default: this Spec depends on none; the operator decides at release, when the catalog's text is known.
3. Is the path of a wiki result relative to the wiki folder in the contract an adopter's session meets? The Secondbrain's request states it, the test of 2026-11-18 shows it, and the author's session received no wiki result to confirm it. Default: the clause states it as the request does, and the TechSpec author confirms it against the contract before fixing the words.
4. Which number does the new ADR take? Default: the next free number when the Task Graph is authored. The authoring base ends at 0447, and the uncommitted plan of Spec 0415 already holds a file numbered 0448, so 0449 is the likely one.
