# How an open authoring question is handled today

Measured by the author of the Idea on 2026-11-19, in this repository at `348962d2da564ab47970e4bd41c16dff0255ebfc`. It answers three questions: what the instructions tell a session to do with a question an author cannot settle, how often a PRD is written with nobody present, and how often the maintainer takes the option a question recommends. A fourth section holds what published tools do with an unanswered question.

## 1. What the instructions say

Line numbers are of this repository at the commit above.

| Subject | Where | What it says |
| --- | --- | --- |
| An Author Subagent with a person behind the session | `.agents/skills/write-idea/SKILL.md:70`, `.agents/skills/write-prd/SKILL.md:89`, `.agents/skills/write-techspec/SKILL.md:70`, `.agents/skills/write-tasks/SKILL.md:456` | "When a decision this skill needs is missing from your brief and the brief does not say that no person is present, write no artifact and return the open questions, each with its options and the one you recommend." |
| The session that started it | `.agents/skills/write-prd/SKILL.md:91` and the same sentence in the other three skills | "When it returns questions, ask the user and continue it with the answers." Nothing says what the session does when the user is away. |
| The PRD with nobody present | `.agents/skills/write-prd/SKILL.md:148-149`, `.agents/skills/write-prd/references/clarification-record.md:14-17` | "When no person is present, settle questions from the available sources and name that route in the Clarification Record", with the line `Autonomous route: no person present; decisions taken from <source>.` |
| The Task Graph | `docs/agents/autonomous-work.md:22` | "Decide the Task Graph autonomously", and "Stop and ask only when the next decision would change authority, architecture, or blast radius". |
| The Idea stage | `.agents/skills/write-idea/SKILL.md:26`, `:44`; `.agents/skills/roundfix/references/modes.md:42` | "Never answer this stage's questions for the user." "Do not write the file until the user approves the draft". "The Idea stage is the user's". |
| What counts as an answer | `.agents/skills/roundfix/references/modes.md:44` | "A recommendation, silence, a dismissal or a timeout is never an answer." |
| Where an answer is recorded | `.agents/skills/roundfix/references/modes.md:38` | "Answers and authorizations about an item that becomes a Spec go into that Spec's `_authorization.md`, and its decisions also go into the PRD's `## Decisions`." |
| The decisions an Author Subagent needs | `.agents/skills/roundfix/references/modes.md:34` | The `unattended` Skill Mode gathers them before the phase starts, because "that subagent cannot ask the user, and it returns its questions instead of guessing one". |
| The record of who decided | `CONTEXT.md:1691` | The Clarification Record holds "one dated line for each question a person answered, or the autonomous-route line naming the source of the author's decisions". Its Avoid line names "assumed approval". |
| The stop-and-ask clause | `docs/agents/agent-instructions.md:90` | "Never guess a decision the user can answer cheaply." |
| A wait | `CONTEXT.md:213` | The Waiting Declaration is "the optional `waiting_on` and `since` pair on an open Backlog Entry or an unresolved Finding". No Spec artifact carries one. |

The gap: a question that no source settles has no route when the user is away. The PRD author may settle only what "the available sources" settle. The session may not take the recommended option, because a recommendation "is never an answer". The Spec therefore stops, and nothing records that it stopped or why, unless the session writes a Backlog Entry by its own judgment.

The plan of Spec 0409, read in its planning clone on 2026-11-19, decides the same case for one kind of question. Its PRD says: "With no person present, no artifact of that phase is written and nobody answers for the user. Only that Spec waits: the feature stays or becomes a Backlog Entry that holds the whole question, names the Spec folder when one exists and carries a Waiting Declaration that names the Rule Challenge, and the other work goes on." The maintainer's answer behind it is "A Spec espera (Recommended)".

## 2. What Project Config allows

| Subject | Where | What it says |
| --- | --- | --- |
| An unknown key | `internal/config/config.go:1910` | The decoder runs with `KnownFields(true)`, so a binary that does not know a key refuses the file. |
| An older binary | `docs/user-guide/configuration.md:215` | "An older binary refuses every command that loads configuration with `field skills not found in type config.configOverlay`". |
| The guard for that | ADR-0392 (`docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`) | A Roundfix older than the one that last applied the repository's Baseline refuses before it reads Project Config, "naming both versions and `roundfix upgrade`". |
| A key that only a project may set | `internal/config/authoring.go:44-62` | `skills` and `delivery.require_defect_class` in User Config are ignored with a warning that names Project Config. |
| Who may write the file | `internal/speccheck/governed.go:50`, `.agents/skills/roundfix/references/modes.md:60`, `CONTEXT.md:770` | `.roundfixrc.yml` is a Governed Path. The `update` Skill Mode leaves it "to the user; never write it". A Task commit includes it "only when the frozen Spec authorization bounds `.roundfixrc.yml`". |

A project that sets no key is read by every binary. Only a project that declares the new preference needs the release that knows it.

## 3. How PRDs were authored, Specs 0364 to 0408

Method:

1. List every `_prd.md` that a commit added under `docs/specs/` (`git log --diff-filter=A --name-only`), keep the Specs numbered 0364 to 0408 and one PRD per Spec number.
2. Read each file as it was in the commit that added it.
3. Count the files that hold a line starting `Autonomous route: `, and the files whose `## Decisions` section holds at least one line starting `- YYYY-MM-DD: `.

Result, over 43 PRDs:

| PRDs | Count |
| --- | --- |
| With the autonomous-route line | 41 |
| With at least one dated answer of a person | 14 |
| With both | 12 |
| With the autonomous-route line and no dated answer | 29 |
| With dated answers and no autonomous-route line | 2 |

Authoring with nobody present is the usual case here: 41 of 43 PRDs name the autonomous route, and 29 of 43 hold no answer of a person. The count says nothing about how many of those phases met a question that no source settled, because no file records one.

## 4. How often the maintainer takes the recommended option

The three answer records of this cycle hold 33 answers of the maintainer: the record of Spec 0409 (13), the record of Spec 0410 (11 that the 0409 record does not also hold) and [the record of this Spec](2026-11-19-maintainer-answers.md) (9, of which one is the free-text answer that is this Spec's request). The first two were read in their planning clones on 2026-11-19.

| Answer | Count |
| --- | --- |
| The option labeled `(Recommended)` | 29 |
| Another option ("Só as marcadas", Spec 0410) | 1 |
| Free text in place of the options (the request of this Spec) | 1 |
| An option the record quotes without a label ("Spec após a release", "Spec 0412 com Ideia agora") | 2 |

The sample is small, comes from one person and three days, and is made mostly of Idea-stage, routing and authorization questions, which stay the user's under this Idea. It is the only count that exists. No record holds a question returned by a PRD, TechSpec or Tasks author, because no form for one exists.

## 5. What published tools do with an unanswered question

Found through Exa by the Research Subagent on 2026-11-19. The author re-read each passage quoted here on the same day.

- LangGraph, "Interrupts", https://docs.langchain.com/oss/python/langgraph/interrupts — "When an interrupt is triggered, LangGraph saves the graph state using its persistence layer and waits indefinitely until you resume execution."
- AG2, "Approval Required", https://docs.ag2.ai/docs/user-guide/tools/approval_required/ — "By default approval_required() waits as long as the human takes." "When the timeout expires the turn ends with HumanInputTimeoutError and the gated tool does not run — an unanswered approval is not an approval."
- Jurgen Appelo, "Leading Agile Developers: The Seven Levels of Authority (Part 2)", InformIT, 2011-02-03, https://www.informit.com/articles/article.aspx?p=1675546&seqNum=3 — level 6, Inquire: "You let the team decide. And afterwards you inquire about their motivations, or you ask that they actively keep you informed." Level 7, Delegate: "you don't even need to know which decisions they make."
- Feng, McDonald and Zhang, "Levels of Autonomy for AI Agents", https://arxiv.org/abs/2506.12469 — "we define five levels of escalating agent autonomy, characterized by the roles a user can take when interacting with an agent: operator, collaborator, consultant, approver, and observer."
- Cuadron and Gupta, "Built on Broken Assumptions (BOBA): Why Coding Agents Cannot Be Left Alone", https://arxiv.org/abs/2611.02207 — "they complete 95% of requests taken one at a time, but when requests are chained together, they complete 51%", and "A single external ``that's wrong'' after a failed request and one repair round lift chain completion from 51% to 59%". The study measures chains of code requests, not planning documents.

Not re-read, and so not quoted: the Research Subagent's search highlights also named the OpenAI Agents SDK, the Claude Agent SDK and two smaller tools as waiting for the person or treating an unanswered request as a refusal. The sweep found no published rate at which people accept a recommended option, and no tool that records a decision as taken on someone's behalf with a fixed form.
