---
name: brainstorming
description: "You MUST use this before any creative work - creating features, building components, adding functionality, or modifying behavior. Explores user intent, requirements and design through one question at a time, then routes the outcome to write-idea, write-prd, write-techspec (refactors/bug fixes), or a direct task."
metadata:
  category: discovery
  tags: [product, requirements, research, workflow]
  version: 0.0.5
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.0.5
---

# Brainstorming

Turn a raw intention into a shared understanding through collaborative dialogue — before any spec is written or any code is touched. This skill is the questioning discipline the discovery stage runs on; `write-idea` and `write-prd` both assume a conversation shaped like this has happened.

<HARD-GATE>
Do NOT invoke any implementation skill, write any code, or scaffold anything until the exploration below has run and the user has confirmed the direction. This applies to EVERY request regardless of perceived simplicity.
</HARD-GATE>

## Anti-pattern: "this is too simple to need exploration"

A todo list, a single-function utility, a config change — all of them go through this. "Simple" requests are where unexamined assumptions cause the most wasted work. The exploration can be three questions and two minutes for genuinely simple work, but it must happen.

## The questioning discipline

- **One question per message.** A topic that needs depth becomes a sequence of questions, never a questionnaire.
- **Multiple choice preferred** when the options are enumerable — use two or three options, put the recommended one first and label it `(Recommended)`, with no `Other` option. Open-ended is fine when options genuinely can't be predetermined. The structured question tool supplies the custom answer itself; without such a tool, show the same single question in chat and accept a custom answer.
- **Explore before asking.** Check project state first — files, docs, `CONTEXT.md`, the accepted ADRs under `docs/adr/`, the guides under `docs/agents/`, active Specs under `docs/specs/`, retired Specs under `docs/history/specs/`, recent commits. A question the codebase can answer is wasted user time.
- **WHAT, WHY, WHO — not HOW.** Databases, APIs, frameworks, and architecture belong to the tech-spec stage; raising them here derails the product conversation.
- **YAGNI ruthlessly.** Challenge every feature against the smallest version that delivers the value.

## Question phases

Walk these in order, skipping what the conversation or the codebase already answered. Minimum 3 questions, maximum ~6 — stop early once problem, user, and scope are clear.

1. **Problem / pain point** — what concrete problem does this solve? What prompted it? Offer 2–3 interpretations if ambiguous.
2. **Target user and context** — who is it for, where in their workflow, what's their current workaround?
3. **Scope** — ideal V1 size: minimal MVP / complete feature / platform. Recommend one of the first two.
4. **Ambition** — quick win / strategic bet / compounding feature; and "what would make this 10x more valuable instead of incremental?"
5. **Dependencies** — does it touch or extend existing features? Name candidates from the codebase.
6. **Success criteria** — how will we know it worked?

## Converging

When understanding feels complete:

1. **Propose 2–3 approaches** with trade-offs, leading with your recommendation and why.
2. **Present the direction in sections** scaled to their complexity (a few sentences when straightforward, 200–300 words when nuanced), validating each with the user before moving on.
3. **Capture language as it resolves** — a term the conversation sharpened goes into `CONTEXT.md` right then; a hard-to-reverse decision that survives the three-part gate becomes an ADR.

## Where the outcome goes

The conversation itself is the artifact — route it, don't duplicate it:

- **Product-level idea** (market questions, open solution shape) → `write-idea` (produces `docs/specs/<slug>/_idea.md`).
- **Feature ready to specify** → `write-prd` (synthesizes this conversation; it will not re-interview), once the Idea-stage judgment below does not send it to `write-idea`.
- **Refactor or bug fix** (no product behavior change) → `write-techspec` directly — it mints the spec folder with a minimal `_prd.md`, no product interview.
- **Trivial, well-understood change** (one-line fix, typo, config tweak) → direct implementation, no spec folder; say so explicitly.

When a repo carries `docs/agents/spec-routing.md`, follow its routing table — it is the canonical version of this list for that project.

Do not write a separate design document — the spec pipeline's artifacts (`_idea.md`, `_prd.md`) are where the outcome lives.

## Rule Challenge

A rule guides the solution and never bounds the goal in silence. A rule is a decision the repository recorded: an accepted ADR, a Normative Clause of a guide under `docs/agents/`, a definition in `CONTEXT.md`, or a document under `docs/references/` that the repository treats as binding. A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem. Extra work, taste and friction are not a block.

When a rule blocks the goal, neither drop or narrow the solution nor contradict the rule: put a Rule Challenge to the user, one question per rule, through the structured question tool the session exposes. Name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why the rule exists and whether it looks out of date or badly written. Offer keep, change and retire, each with its consequence, the recommended one first and labeled `(Recommended)`. The rule binds until the user answers: asking changes nothing, and a recommendation, silence, a dismissal or a timeout is never an answer.

- Ask no Rule Challenge about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records. Name that request or answer where the artifact records the change.
- Ask none about a rule that carries a Kept Line, a line `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.` after its text, unless the question names a fact the earlier one did not have. Record that fact as `new fact`.
- Ask at most three in one Authoring Phase or brainstorming session. When more rules block, ask about the three that prevent the most and list the others, which keep binding, in the artifact's open questions.
- Never answer for the user. An Author Subagent returns the question, with its options and the one it recommends, to the session that started it. When no person is present, write no artifact for this phase: the Spec waits as a Backlog Entry that holds the whole question in the record form with `answer: waiting`, names the Spec folder when one exists, and carries `waiting_on`, naming the Rule Challenge, and `since` in its frontmatter.

What follows the answer:

- Change or retire: the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it. An ADR changes through the lifecycle of `docs/agents/docs-layout.md`.
- Keep: the solution stays inside the rule, and the PRD requires a Task to add the Kept Line after the rule's text. The line names no Spec.
- A rule inside `setup-context-driven` markers changes only in the Baseline, and its Kept Line goes in `docs/agents/specific-repository.md`, under a heading of its own, naming the rule. On change or retire, a repository that does not own the Baseline sends the user's decision as one ai-memory message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out.

In this conversation, when the direction the user wants meets a rule that blocks it, say so and ask the Rule Challenge instead of steering the user away from that direction. Hand the question and the user's answer to the skill the outcome routes to, which records them in the `## Rule challenges` section of its artifact; when no Spec follows a keep answer, add the Kept Line yourself.

## Idea-stage judgment

Before a feature goes to `write-prd`, judge whether it needs the Idea stage, unless the user asked for `write-idea` or the feature already has an `_idea.md`. The criteria are three: the right shape of the solution is still open; the request rests on an assumption about its users, the business or the market that it gives no evidence for; deciding to build it depends on a market or competitive question.

When the repository runs Roundfix, ask the judge. For a request recorded as an open Backlog Entry, a note in `docs/_inbox/` or a message handed in as a file, run `roundfix spec judge --stage idea --item <reference>`. For a request that exists only in the session, write it to a file in a temporary directory outside the repository, run `roundfix spec judge --stage idea --request <path>` and remove the file; only the project's name, the request's first line and its first 4,000 characters reach the judge, after the secret scan. When the command says the judge did not answer, or the repository does not run Roundfix, answer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer.

- Any criterion yes: the Idea stage is indicated. Start `write-idea` and show the answers; the user may dismiss the stage at its first question.
- Every criterion no: go to `write-prd`.
- Anything else, or any doubt: ask the user whether the Idea stage runs, showing the answers.

Never present your own answers as the judge's, and never answer the Idea stage's questions for the user. Hand the answers to the skill that runs next, which records them. A fix or a refactor is never judged.

## Key principles

One question at a time · multiple choice preferred · explore before asking · YAGNI ruthlessly · always propose alternatives before settling · validate incrementally · be ready to loop back when an answer breaks an earlier assumption.
