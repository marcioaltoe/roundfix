---
name: business-storyteller
description: Transforms technical work into persuasive business-language documents for non-technical stakeholders, including internal announcements, approval proposals, incident explainers, technical-debt cases, product explainers, and investor or partner memos. Use when the user asks to present technical work to the company, explain software to business teams, write a business case, make an approval document, or convert a technical artifact into stakeholder communication in the requested language. Do not use for READMEs, specs, ADRs, external marketing pages, or customer-facing sales copy.
metadata:
  category: writing
  tags: [business-writing, persuasion, storytelling, internal-comms, stakeholders, pdf, charts]
  version: 0.4.0
  author: Marcio Altoé
---

# Business storyteller

Turn technical work into documents that non-technical colleagues read, understand, and approve. The craft is translation plus persuasion: start from the business outcome and work backwards to the technology, prove every claim with a real number, and tell it as a story the reader can repeat in their next meeting. Professional and rich, never banal; persuasive, never manipulative.

## When to use

- Presenting a feature, product, or release to the whole company
- Writing an approval proposal for a feature, refactor, migration, or technical investment
- Explaining a fix or incident to management, compliance, or customer-facing teams
- Explaining how the software works to admin, sales, marketing, HR, finance, or accounting
- Turning any technical artifact (spec, PR, changelog, postmortem) into business communication

## Hard rules

1. **Relevance, not simplification.** The reader does not need to understand how the system works — they need to know what changes for the company: money earned or saved, risk removed, time freed, capability unlocked. Start from that impact and work backwards; mention technology only when it earns its place.
2. **Outcome before capability.** "Reports now load in 2 seconds instead of 40" — never "we optimized the database queries". Every technical fact gets translated through the table in the persuasion playbook.
3. **True numbers only.** Concrete numbers persuade ("87% of users" beats "most users") — but only numbers the source provides. Never invent metrics, savings estimates, percentages, or deadlines. The rule also covers small "harmless" details writers add for texture: a meeting length, a sprint duration, a "30-minute" anything — if the source does not state it, it does not exist. Simple arithmetic on source numbers is fine (400 invoices × 15 min = 100 h/month), but anchor derived time spans to the right endpoints (time-from-detection is not time-from-first-failure). Before delivering, point every number in the draft back to a source fact; a number with no source gets cut or replaced with a qualitative claim. A document that wins approval on invented numbers loses trust forever. If a number would help and does not exist, ask the user for it.
4. **One ask per document.** End with exactly one clear call to action: the decision needed, who decides, and by when. Two asks compete; three asks lose.
5. **Write in the language of the request.** Portuguese request → Portuguese document. English request → English document. The same persuasion and banned-pattern rules apply in every language.
6. **Persuade, never manipulate.** Loss framing, social proof, and urgency are allowed only when factually true. No false scarcity, no invented testimonials, no inflated risk. Front-load bad news — readers forgive problems, not surprises.

## Process

1. **Name the audience and the decision.** Who reads this (a department, the leadership, the whole company)? What should they think, feel, or approve after reading? Weigh the message with the CRG model: finance reads cost, operations reads risk, leadership reads growth — same facts, different lead.
2. **Gather the source facts.** Read the spec, PR, changelog, metrics, or conversation. List the facts and real numbers available. These are the only raw materials.
3. **Pick the document type and load both references.** Read [references/persuasion-playbook.md](references/persuasion-playbook.md) for the techniques and [references/document-templates.md](references/document-templates.md) for the matching skeleton: feature announcement, approval proposal, fix/incident explainer, technical-debt case, how-it-works explainer, executive one-pager, or investor/strategic-partner product memo.
4. **Translate.** Run every technical fact through capability → outcome. Replace each jargon term with its business meaning or cut it. If a term must stay (compliance, audit), define it in one plain sentence on first use. Jargon includes the words engineers stop noticing: module, refactor, sprint, deploy, rollback, staging, pipeline, environment, dependency, endpoint, backend, schema, migration, API, latency, p95 (módulo, refatoração, sprint, implantação, migração). The test: would a person in HR or accounting know this word from their own job? If not, translate it ("the billing module" → "the part of the system that calculates invoices"; "sprint time" → "the team's working time") or define it in passing. Translate the term, not into a guess: "sprint" does not become "two weeks" unless the source states the sprint length.
5. **Draft answer-first.** The first sentence carries the conclusion: what this is and why the reader should care. Then the story, then the proof, then the ask. No throat-clearing, no "this document describes".
6. **Persuasion pass.** Apply the playbook: quantified cost of inaction, the reader as hero of the story, concrete numbers, one memorable phrase the reader will repeat, confident language without hedging.
7. **Humanize pass.** Remove every banned pattern (playbook list). Vary sentence rhythm. The document must read like a sharp colleague wrote it, not a machine or a press release.
8. **Produce the output.** Markdown is the canonical deliverable. For PDF or HTML, follow the Output pipeline below.
9. **Self-check** (bottom of this file) before delivering.

## Output pipeline

Markdown first, always — it is the reviewable source of truth.

**Charts** (only when data exists in the source): generate self-contained SVG with the bundled script (bootstrap helper — it creates files, modifies nothing else). Run it from the skill directory:

```sh
python3 scripts/make_chart.py --type bar --title "Hours saved per week" \
  --data "Support:14,Finance:9,Operations:6" --out chart-hours.svg
```

Types: `bar`, `line`, `donut`. Data values must come from the source — a chart of invented numbers is fabrication with extra polish.

**HTML**: copy [templates/document.html](templates/document.html), replace the `{{...}}` placeholders, and write each section as simple HTML (`<h2>`, `<p>`, `<ul>`, `<table>`, `<blockquote>`). Inline any SVG charts inside `<figure>` blocks. The template is self-contained (no CDN, no network) and print-ready.

**PDF**: render the HTML with the bundled export script (mutating helper — it writes the PDF file). It finds Chrome/Chromium automatically and reports the page count:

```sh
sh scripts/export_pdf.sh document.html document.pdf
# OK: document.pdf (pages: 1)
```

If no Chrome/Chromium is installed, deliver the HTML and tell the user it prints to PDF from any browser (File → Print → Save as PDF).

**One-pager page limit.** The executive one-pager must fit a single A4 page. Use the template's compact mode (`<body class="compact">`) and check the `pages:` count the export script prints. If it still exceeds one page, do NOT cut content on your own — every section was built from facts the user cares about. List the candidate cuts (with what each would lose) and ask the user to decide what to cut, or whether a second page is acceptable. The user decides; you propose.

## Self-check

Every answer must be yes before delivering:

1. Does the first sentence state the conclusion and why this reader should care?
2. Could a person from HR or accounting read the whole document without stumbling on one unexplained technical term?
3. Is every number real — present in the source or confirmed by the user?
4. Is the cost of doing nothing stated (when the document asks for a decision)?
5. Is there exactly one call to action, with owner and deadline?
6. Does it pass the banned-patterns scan (AI vocabulary, hedging, decoration, sycophancy)?
7. Is there one phrase memorable enough that a reader would repeat it in a meeting?
8. Would the target reader forward this to their boss as-is?

## Anti-patterns

- Explaining the implementation ("we migrated to microservices") instead of the consequence ("new features now ship in days, not months").
- Burying the ask on the last page or splitting it into several requests.
- Inventing savings, percentages, or adoption numbers to strengthen the case.
- Dumbing down instead of translating — business readers are smart; they lack context, not intelligence.
- Marketing fluff: superlatives, exclamation points, "revolutionary", emoji decoration.
- Hedging the recommendation ("we believe it might be beneficial...") — recommend with confidence or do not recommend.
- Hiding bad news in the middle of a paragraph. Lead with it and pair it with the plan.
