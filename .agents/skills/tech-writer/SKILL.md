---
name: tech-writer
description: Write or rewrite anything technical people write — READMEs, documentation pages, quickstarts, PRDs, tech specs, RFCs, ADRs, design docs, task breakdowns and task files, GitHub issues, bug reports, QA test plans, test cases, exploratory charters, verification reports, PR descriptions, code review comments, and async team updates — so it reads clearly, scans fast, stays faithful to source facts, and sounds human. Use when the user says "write a README", "improve these docs", "write a PRD", "write a tech spec", "draft an ADR", "break this into tasks", "write the task files", "open an issue", "write a test plan", "write test cases", "file a bug report", "write the QA report", "write the PR description", "draft a status update", or wants technical writing that persuades without marketing fluff. Do NOT use for marketing landing pages, sales copy, or ad copywriting.
metadata:
  category: writing
  tags: [readme, docs, technical-writing, prd, techspec, adr, tasks, issues, qa, communication]
  version: 0.4.0
  author: Marcio Altoé
---

# Tech writer

Write technical content that earns the reader's attention: lead with what matters most, prove claims with real facts and runnable examples, and cut everything that does not help the reader decide or act. Covers the full document chain a senior engineer or architect produces — from idea brief and PRD through tech spec, ADRs, task files, and issues, down to QA test plans, bug reports, and verification reports — plus everyday docs, PRs, reviews, and updates.

## When to use

- Writing a README, quickstart, docs page, or guide for a tool, library, or service
- Writing or rewriting a PRD, tech spec, RFC, ADR, design doc, idea brief, GitHub issue, bug report, PR description, code review comment, or async status update
- Decomposing a PRD or spec into a task list and individual task files, or into vertical-slice issues
- Writing QA artifacts: test plans, test cases, exploratory charters, bug reports, verification reports
- Rewriting any technical text that feels bloated, vague, or AI-generated

## Two modes: transcribe or co-author

Decide the mode before writing; it changes what is allowed.

- **Transcribe** — a source exists (code, spec, ticket, diff, conversation, prior doc). Write only what the source supports. Never invent acceptance criteria, metrics, edge cases, alternatives, or behavior. A `TBD — needs <owner>` marker beats an invented detail.
- **Co-author** — the user is creating the document with you. Propose content for empty sections, but mark every proposal (`Proposed:` prefix or an explicit callout) and ask before it becomes fact. Ask one question at a time, prefer multiple-choice when the options are knowable, and park unresolved points under Open questions with a named owner. A proposal the user confirms becomes a fact; one they have not seen stays a proposal.

## Process

1. **Read the source first.** Examine the code, CLI help, API surface, spec, ticket, or diff before writing. Every claim must be backed by something real.
2. **Name the reader and the action.** Decide who reads this (new user, reviewer, PM, implementing engineer, QA, future maintainer) and the one thing they should know or do after reading. Structure everything around that.
3. **Pick the doc type and load its playbook.** Read only the matching reference file from the router below — it carries the structure rules, skeleton, and quality gates for that type.
4. **Draft inverted-pyramid.** The most important sentence comes first. For docs: what it does and the concrete outcome. For a PR: the standalone summary line. For an update: the TL;DR. No throat-clearing ("Welcome to...", "This document describes...").
5. **Cut pass.** Delete needless words, qualifiers, and duplicate ideas. If a sentence survives unchanged after removing a word, the word was noise.
6. **Humanize pass.** Scan for the banned patterns below and rewrite every hit.
7. **Verify.** Run the self-check. Confirm every command runs, every link resolves, every stated fact exists in the source, and the doc passes its playbook's quality gates.

## Doc-type router

| Document                                                                                              | Playbook                                                             |
| ----------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| PRD (full or lightweight), idea/discovery brief, user stories, acceptance criteria                    | [references/prd-playbook.md](references/prd-playbook.md)             |
| Tech spec, design doc, RFC, ADR                                                                       | [references/techspec-playbook.md](references/techspec-playbook.md)   |
| Task list, task files, vertical-slice issues, work decomposition                                      | [references/tasks-playbook.md](references/tasks-playbook.md)         |
| QA test plan, test cases, exploratory charters, QA bug reports, verification report                   | [references/qa-playbook.md](references/qa-playbook.md)               |
| README, docs page, quickstart, tracker issue/bug report, PR description, review comment, async update | [references/doc-type-playbooks.md](references/doc-type-playbooks.md) |

## Senior judgment

These are the habits that separate a senior architect's documents from a note-taker's transcripts. Apply them across every type.

- **Control the altitude.** A PRD argues what and why and never prescribes architecture. A spec defines exactly what the system does — contracts, data models, failure modes — at the precision the decision needs. A task says what to accomplish, not how. Mixing altitudes is the most common structural failure.
- **Non-goals are load-bearing.** An explicit out-of-scope list prevents more rework than any other section. Every planning doc gets one.
- **Name the trade-offs.** A spec with no downsides is not finished. State what the design gives up, which alternatives were rejected and why, and the risks that remain. Reviewers trust documents that name their own weaknesses.
- **Numbers over adjectives.** "Fast" is a wish; "p95 under 200ms" is a requirement. Push every non-functional claim toward a measurable target — and in transcribe mode, only the targets the source states.
- **Acceptance criteria are observables.** Each criterion must be checkable by someone who did not write the code: a behavior, an output, a passing command. Given-When-Then where it helps; include negative and edge cases.
- **Keep the traceability chain.** Every spec component traces to a PRD goal; every task to a spec section; every test to a task or journey; every bug to a reproduction. A reader must be able to walk the chain in both directions through links, not duplication.
- **Reference, never duplicate.** Downstream docs cite upstream sections by name ("see Tech spec → Data Models") instead of restating them. Duplicated content drifts; links don't.
- **Decompose vertically.** Slice work into thin end-to-end paths (schema + API + UI + tests) that are independently implementable and demoable — never horizontal layers. Map dependencies explicitly; front-load the riskiest slice.
- **Plan the rollout and the rollback.** A design that ships behind a flag with a tested rollback path is a different (better) design than the same code without one. Specs state both when the source decides them.
- **Write for the executor — human or agent.** Tasks and issues are prompts: a clear problem, complete acceptance criteria, explicit boundaries, and the minimum context needed. Vague input produces vague work.

## Writing rules

- **Active voice.** "The CLI writes a report" — not "a report is written".
- **Concrete over abstract — but only as concrete as the source.** Name the commands, versions, file paths, and numbers the source gives. "Retries 3 times with exponential backoff" beats "handles failures gracefully". But if the source says only "exponential backoff", stop there: never compute example delays, multipliers, or formulas the source does not state. Source-stated concreteness is credibility; invented concreteness is a bug.
- **Examples and glosses obey the same rule.** A usage example may only demonstrate flag values, formats, and inputs the source documents. When accepted values are undocumented, keep the placeholder (`--level <level>`) or point the reader to `--help` — never demonstrate with a guessed value (`--level warn`). Likewise, describe what a field or option does only as far as its comment or docs state: "default true" documents the default, not what `false` does or why anyone would set it.
- **Why over what.** Code, diffs, and schemas already show what. Your text adds why: intent, trade-offs, rejected alternatives.
- **One topic per artifact.** One decision per ADR, one idea per issue, one purpose per page.
- **Must / can / might — not "should".** "Must" for required, "can" for optional, "might" for possible outcomes. "Should" leaves the reader guessing.
- **Benefits with proof.** State what the reader gets, then show the code, command, or evidence that delivers it.
- **One idea per paragraph, four sentences max.** Long paragraphs hide the point.
- **Sentence-case headings.** "Getting started", not "Getting Started".
- **Runnable examples.** Snippets must be copy-pasteable, complete, and shorter than one screen. Prefer one real example over three abstract ones.
- **Plain words.** "Use" not "utilize", "help" not "facilitate", "fast" not "performant".

## Banned patterns

Rewrite any occurrence of these. They mark text as machine-generated filler.

- **AI vocabulary:** delve, leverage, seamless, robust, cutting-edge, crucial, pivotal, comprehensive, testament, landscape (abstract), showcase, foster, tapestry, intricate, vibrant, streamline, empower, elevate, unlock, supercharge.
- **Filler phrases:** "in order to" → "to"; "due to the fact that" → "because"; "it is important to note that" → delete; "at this point in time" → "now"; "please note that" → delete.
- **Minimizers:** "simply", "easily", "just", "it's that simple". If it were simple, the reader would not need the doc.
- **Inflated significance:** "plays a vital role", "marks a pivotal step", "represents a paradigm shift".
- **Superficial -ing tails:** "...ensuring reliability", "...showcasing flexibility", "...highlighting its power". Cut the tail or replace it with a fact.
- **Negative parallelism:** "It's not just X, it's Y."
- **Rule-of-three padding:** forced triads like "fast, flexible, and powerful". Use the one adjective that is true and provable.
- **Bold-header bullet lists:** bullets shaped like "**Performance:** Performance has been improved...". Write prose or plain bullets instead.
- **Decoration:** emojis in headings or bullets, exclamation points, Title Case headings, bold scattered for emphasis.
- **Vague attribution:** "many developers say", "users want", "it is widely considered". Cite the ticket, survey, or source — or cut the claim.
- **Generic endings:** "The future looks bright", "Happy coding!". End with the next concrete step instead.

## Self-check

Before delivering, every answer must be yes:

1. Does the first sentence carry the single most important point for this reader?
2. Is every code block complete and runnable as written?
3. Did a scan for banned patterns come back clean?
4. Is every claim specific — backed by a command, number, option name, citation, or example?
5. Does the doc follow its type's structure and pass the quality gates in its playbook?
6. In transcribe mode: does every fact exist in the source? In co-author mode: is every unconfirmed proposal still marked as one?
7. Could the target reader (or an executing agent) act using only this text plus its links?

## Anti-patterns

- Writing before reading the code, spec, ticket, or diff.
- Opening with history, philosophy, or "In today's fast-paced world...".
- Prescribing implementation in a PRD, or re-arguing product strategy in a spec.
- Duplicating upstream content downstream instead of referencing the section.
- Tasks sliced by layer ("build all the models, then all the APIs") instead of by outcome.
- Acceptance criteria nobody can check ("works correctly", "handles errors gracefully").
- Padding length to look thorough. A short doc that works beats a long one that impresses.
- Inventing benchmarks, testimonials, metrics, or behavior not present in the source.
- Claiming QA verdicts without fresh evidence, or omitting what was not tested.
