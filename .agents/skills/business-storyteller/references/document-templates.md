# Document templates

Seven document types. Pick by what the reader must do after reading: celebrate and adopt (announcement), approve (proposal, debt case), regain confidence (incident explainer), understand (how-it-works), decide fast (one-pager), or commit capital/partnership (investor memo). Skeletons are in English; write the actual document in the language of the request, translating headings naturally.

Every skeleton follows the same spine: conclusion first → story → proof → one ask. Fill sections only with facts from the source; mark genuinely missing numbers as a question to the user, never as an invention.

## 1. Feature or product announcement

Audience: the whole company. Goal: everyone understands what changed, why it matters, and what to do differently. Inspired by press-release-style internal launches: write it as if announcing to the world, then answer the questions each department will actually ask.

```markdown
# <headline: the outcome, not the feature name>

<One paragraph, written like a news lead: what is now possible, for whom,
and the single strongest number. A colleague should be able to forward
just this paragraph.>

## What changes for you

<2–4 bullets, each "you/your team" framed: what readers can do now,
what they can stop doing. The five-second moment lives here.>

## The story behind it

<Short story-spine narrative: the pain that existed, what it cost
(real numbers), what was built — in business terms — and the after-state.>

## Proof

<The pilot result, the before/after metric, a quote from a real internal
user. A chart earns its place here when the data exists.>

## Questions you may have

<4–6 Q&A pairs, one per affected department: does this change my process?
costs? training? when? who to contact? Anticipate the skeptic's question
and answer it plainly.>

## What to do now

<The single action: try it, enroll, read the guide — one link, one step.>
```

## 2. Approval proposal

Audience: whoever signs off (leadership, finance, a committee). Goal: a confident yes. This is the most persuasion-dense type: cost of inaction, limited options, evidence, one ask with deadline.

```markdown
# <headline: the decision and its payoff in one line>

<Three-line executive summary:
1. The recommendation in one sentence.
2. Why it matters — the strongest real number (gain or avoided loss).
3. The ask: decision, owner, deadline.>

## The problem today

<The pain as the reader experiences it, with its current cost quantified:
hours lost, money spent, risk carried, opportunities missed. This is the
cost-of-inaction section — it does the heavy lifting.>

## What we propose

<The solution in business terms: what changes, for whom, when. Technology
appears only where it earns its place. Paint the after-state in one
concrete paragraph.>

## Options considered

<At most three options including "do nothing", as a small table:
option / cost / what we get / risk. Recommend one and say why in one
sentence. "Do nothing" carries its quantified cost from the section above.>

## Proof it works

<Pilot data, comparable cases, internal measurements, vendor benchmarks
with source. Charts welcome when the data is real.>

## Costs, risks, and limits

<Front-loaded honesty: what it costs, what could go wrong, what this does
NOT solve — each risk paired with its mitigation. This section buys the
credibility the recommendation spends.>

## The plan

<3–5 milestones with dates and owners. Small enough to feel safe,
concrete enough to feel real.>

## Decision needed

<One sentence: what to approve, who decides, by when, and what the
deadline protects ("approval by Friday keeps the rollout inside Q3").>
```

## 3. Fix or incident explainer

Audience: management, support, compliance, affected teams. Goal: restore confidence. Structure is trust-first: lead with what happened and the impact — never minimize — then show control.

```markdown
# <headline: what was fixed and what is now guaranteed>

<One paragraph: what happened in business terms, who was affected and for
how long, and the one-line status now. Bad news first, plainly.>

## What happened

<The event as the business experienced it: "between X and Y, invoices
went out with wrong dates for customers east of UTC". No blame, no
jargon, no drama.>

## The impact

<Honest numbers: how many customers/records/hours. If the impact was
small, the numbers prove it; if it was large, hiding it would cost more
than admitting it.>

## What we did

<The response timeline in plain language: detected, contained, fixed,
verified. Hours and dates show competence better than adjectives.>

## Why it will not happen again

<The structural change, translated: the new safeguard, the test, the
alert. This is the section the reader came for — make it concrete, not
promissory.>

## What we ask of you

<One action if any: inform affected clients with this summary, update a
procedure — or state plainly that no action is needed.>
```

## 4. Technical-debt or refactor case

Audience: decision-makers who see no visible product change. Goal: get invisible work approved. The whole document translates engineering hygiene into risk and speed — the two things invisible work actually buys.

```markdown
# <headline: the business capability this work protects or unlocks>

<Three-line executive summary: recommendation, the risk carried or speed
lost today (with numbers), the ask.>

## What is slowing us down (or what could break)

<The debt in consequence terms: "every change to billing takes 3 weeks
and breaks something every other release" — never "the code is messy".
Quantify: cycle time, incident count, hours of rework, dependency risk.>

## The cost of doing nothing

<Project today's numbers forward honestly: more incidents as volume
grows, slower delivery as the team grows, the risk event waiting to
happen. Loss framing, true numbers only.>

## What we propose

<The work in outcome terms: "after this, billing changes ship in days
and month-end runs without intervention". Duration, team, what pauses
meanwhile — the honest price.>

## How we will prove it worked

<The before/after metrics the reader can check later: cycle time,
incidents, support tickets. Committing to measurable outcomes is what
separates a case from a wish.>

## Decision needed

<One ask, owner, deadline.>
```

## 5. How-it-works explainer

Audience: a department that uses or depends on the software. Goal: understanding and autonomy, not approval. A guided tour by outcomes — organized around what the reader does, never around the architecture.

```markdown
# <headline: what the software does for this reader>

<One paragraph: the job the system does in this reader's world, and the
one thing they will understand by the end.>

## The big picture

<The system in one analogy-free paragraph: what goes in, what comes out,
who touches it. A simple flow diagram or chart helps here more than
anywhere else.>

## What happens when you <core action 1>

<Step-by-step from the reader's seat: "you click X; the system checks Y;
finance sees Z". Each step is something observable, not a component name.>

## What happens when you <core action 2>

<Same pattern. Cover the 2–4 actions this audience actually performs.>

## When something looks wrong

<The reader's map: what the common situations mean ("pending means the
bank has not confirmed yet — it resolves within 2 hours") and when to
call whom.>

## Words you will hear

<A short glossary of the terms that leak from engineering into meetings,
each in one plain sentence. This section quietly becomes the most
consulted part of the document.>
```

## 6. Executive one-pager

Audience: a decision-maker with two minutes. Goal: a decision. One page, hard limit — it forces the thinking. Often the companion of a fuller proposal; works standalone when trust is high. This is the type that benefits most from the PDF output.

Enforcing the limit: render the HTML with the template's compact mode (`<body class="compact">`) and verify the page count the export script reports. Over one page even in compact mode? Do not trim silently — present the candidate cuts to the user (each with what it would lose) and let the user choose what goes, or accept the extra page. Cutting is the user's decision, never the writer's.

```markdown
# <headline with the outcome and the number>

**Recommendation:** <one sentence>
**Impact:** <the strongest number: gain, saving, or avoided loss>
**Ask:** <decision, owner, deadline>

## Why now

<2–3 sentences: the trigger and the cost of waiting.>

## The picture

<One chart or one comparison table — today vs. after. The single
strongest piece of evidence, visualized.>

## Costs and risks

<2–3 bullets, each risk with its mitigation. Honest and short.>

## Next step

<The one action and the date.>
```

## 7. Investor and strategic-partner product memo

Audience: investors, strategic partners, potential co-founders. Goal: conviction to commit capital, partnership, or a terms conversation. A written memo, not slides: standalone, 2–5 pages, must read without a presenter in the room. The structure condenses what funded pitches share (Sequoia's framework, the YC seed deck and memo, analyzed decks like Airbnb's and Buffer's): the memo is an argument — every section makes one claim and proves it with a number or a named fact.

Ordering rule: when a real traction signal exists (design partner, pilot, build in progress with dates), it goes before the problem — the modern reader decides in the first sections, and proof answers "is this real?" faster than any market analysis. Without traction, open with why-now.

```markdown
# <headline: the thesis carrying its strongest number>

<Three-line summary: (1) what the product is and for whom,
(2) the business model with its strongest economics number,
(3) the ask: what you seek, from whom, in what window.>

## Traction and why now

<Proof before theory: the strongest real signals, with dates.
Then the 2–4 named, falsifiable shifts that make this possible
today and not two years ago. Timing risk is the question this
section exists to answer; it is the most skipped and the most
weighed.>

## The problem

<The pain as the customer lives it: concrete frictions, quantified
where the source allows. Trend statements let the reader nod along
without engaging; specific frictions force a response.>

## The product

<How it works from the buyer's seat, in numbered steps; then the
distribution or positioning angle (white-label, channel strategy).
Behavior change, not feature lists.>

## Market and who we sell to

<Bottom-up math only: count of addressable buyers × price, source
named. "The market is $50B" reads as a googled number; "42,000
property managers spending $3,200/year" reads as homework done.
Name the 2–3 buyer profiles and why each one buys.>

## Competition

<Named competitors plus the two unnamed ones every reader prices in:
"do nothing" and "build it in-house". Say honestly where each wins;
then the moat that is not "we are 6 months ahead".>

## Business model and projected margin

<Price, cost base, gross margin at 2+ volume points. Unknown costs
become explicit variables or labeled hypotheses, never invented
numbers. Include the scenario where the model does NOT close —
honesty here buys credibility everywhere else.>

## Strategic choices that protect the business

<The 1–2 architecture or positioning decisions that de-risk the
company (provider abstraction, regulatory path), translated to
consequences: switching cost, negotiating power, time-to-market.>

## Team

<Founder-market fit, not titles: what in this team's history makes
them the ones to win this market now. Investors consistently rank
team as the first factor.>

## Risks

<Named risks with mitigations, hardest first. A memo that names its
own risks reads as diligence already done.>

## Vision

<Scope evolution in 2–3 steps: today → next → the larger play.
No dates the source does not give.>

## The ask

<One ask: what you want, from whom, by when, and what it buys in
named milestones. "Capital to fuel growth" is a non-ask; "X to reach
Y by Z" is.>
```

## Choosing and combining

| Situation                                       | Type                                   |
| ----------------------------------------------- | -------------------------------------- |
| "We shipped something — tell the company"       | 1. Announcement                        |
| "We need a yes on building/buying something"    | 2. Proposal (+ 6. One-pager as cover)  |
| "Something broke and people heard about it"     | 3. Incident explainer                  |
| "We need time for invisible engineering work"   | 4. Debt case (+ 6. One-pager as cover) |
| "A department keeps asking how this works"      | 5. How-it-works                        |
| "The decision-maker has two minutes"            | 6. One-pager                           |
| "We need capital or a strategic partner to commit" | 7. Investor memo (+ 6. One-pager as cover) |

A proposal or debt case sent to busy leadership travels best as a one-pager on top with the full document attached. Generate both from the same facts; the one-pager is the executive summary promoted to its own artifact.
