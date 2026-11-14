---
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
created: 2026-11-14
status: proposed
leverage: compounding
---

# A feature is judged for the Idea stage before its PRD

## Overview

Before the PRD of a feature is written, Roundfix judges whether the request needs the Idea stage, the study that answers "should we build this, and in what shape?", and the agent then starts `write-idea` without waiting to be asked. The judge is the one Roundfix already uses for its Advisory Judgments; when it cannot answer, the agent answers the same criteria itself. Each feature Spec records the judgment and what followed, and a later report re-measures the thresholds from that record.

It is for the person who supervises Spec authoring in any adopting repository. The judge serves a repository with a judge key; the criteria serve every repository.

It is valuable because the Idea stage runs today only when someone remembers to ask for it, and the waste the maintainer wants to avoid is a feature built in the wrong shape.

V1 is ambitious on purpose, in two Specs: the judgment with its record first, then the report that re-measures the thresholds and the judgment shown per item in the follow-up report.

## Problem

The routing rule sends only a "large or fuzzy product initiative" through `write-idea`, and tells the agent to "start with the smaller sufficient route when two routes fit". The agent decides alone, nothing records why the stage was skipped, and `write-prd` never sends a request back to it. In practice the stage runs when the user starts it.

**A feature in the wrong shape.** On 2026-11-14 the request behind Spec 0392 named its own shape: write a model into each skill's header from Project Config. The Idea stage measured that the harness ignores that header whenever the agent calls the skill, which is the main path, and the shape changed to generated subagents in two Specs before any PRD existed. Had the request gone straight to the PRD, as it was about to, the Spec would have been delivered and failed its own purpose.

**A feature that should not exist.** On 2026-11-13 a study of a Homebrew tap, asked for by hand, led the maintainer to decline the feature. Nothing but his request made the study happen.

**Why a rule alone falls short.** An agent asked to judge by itself tends to proceed: the smaller route is the PRD. A judgment that leaves no trace cannot be measured, so nobody learns whether skipping the stage was right.

### Market Data

- Kohavi and colleagues (KDD 2013): "Only one third of the ideas tested at Microsoft improved the metric(s) they were designed to improve." The evidence is strong for web products measured with controlled experiments, after the ideas were built; it shows that ideas often fail, not that an earlier study would have caught them.
- Marty Cagan treats discovery as a judgment of risk and consequence: for minor items "we just add the work to the backlog and move on", and treating every risk as major means "you'll move way too slowly".
- The spec-driven tools surveyed leave the choice to the user or to a size rule: BMAD (tracks by size, the user confirms), Spec Kit (`clarify` is started by the user), Kiro (the user picks Vibe or Spec), OpenSpec and Agent OS (started by the user). Task Master scores complexity with a model, to expand tasks. No tool was found in which a model decides whether a discovery stage runs; this is an absence in the search, not proof.
- Automation bias is measured: people accept an automated recommendation without checking it, experts included (Parasuraman and Manzey, 2010). Showing the confidence and giving information instead of a bare recommendation reduce it (Goddard and colleagues, 2012).
- Not used, because contested: the Standish figure of 64% of features rarely or never used, and Boehm's cost-of-change curve.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The Idea-stage judgment | Critical | Before a feature's PRD, the judge answers a fixed set of yes or no questions about the request (is the shape of the solution still open, does it rest on an assumption with no evidence, does the decision depend on a market question) and Roundfix reports the stage as indicated, not indicated or undecided, with the probabilities. |
| F2  | The rule that acts on it | Critical | Indicated: the agent starts `write-idea` and shows the numbers; the user may dismiss the stage at its first question, and the dismissal is recorded in the PRD with its reason. Not indicated: the PRD follows. Undecided: the user decides, seeing the numbers. No command refuses anything. |
| F3  | The same criteria without the judge | Critical | When the judge cannot answer (no key, inside a Run, the monthly ceiling, a text it does not accept), the agent answers the same criteria itself, says that the judge was unavailable, and follows the same rule; in doubt it asks the user. |
| F4  | The judgment record | High | Each feature Spec records who judged (the judge or the agent) and the answers, whether the Idea stage ran or was dismissed, and what the Idea did to the shape: changed, split, shelved or unchanged. |
| F5  | Requests in Portuguese | High | A request written in Portuguese is judged, and the thresholds are measured on Portuguese cases too, because that is how this fleet's requests arrive. |
| F6  | The calibration report | Medium | A report re-measures the thresholds against the repository's own record (judgments, outcomes and the corrections each Spec needed afterwards) and says what it would change. It changes nothing by itself. Second Spec. |
| F7  | The judgment in the follow-up report | Medium | The follow-up report carries the judgment for each feature item that is not yet a Spec, so a plan shows it before authoring starts. Second Spec. |

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Corrections per feature Spec (plan corrections, QA failures, review rounds with findings) | -30% against the 20 feature Specs delivered before the release | Delivery retries and review rounds recorded for each Spec; the calibration report computes it. |
| Indicated features whose Idea changed the shape (changed, split or shelved) | >= 33% | The judgment record of each feature Spec. |
| Dismissals of an indicated Idea stage | < 20% | Dismissals recorded in the PRDs. |
| Feature Specs with a judgment record | 100% of those authored after the release | The PRDs of feature Specs. |
| False skips: features judged "not indicated" that later changed scope | < 10% | Superseding and corrective Specs that name a feature Spec judged "not indicated". |
| Judge cost per judgment | < US$ 0.01 | The Judge Log. |

## Feature Assessment

| Criteria        | Question                                            | Score  |
| --------------- | --------------------------------------------------- | ------ |
| Impact          | How much more valuable does this make the product?  | Strong: it catches the wrong shape before Tasks are cut; it did so once on the day it was proposed. |
| Reach           | What % of users would this affect?                  | Strong: the criteria apply in every adopting repository, with or without a judge key. |
| Frequency       | How often would users encounter this value?         | Maybe: once per feature, and most Specs in this fleet are fixes. |
| Differentiation | Does this set us apart or just match competitors?   | Strong: no tool surveyed lets a model decide whether discovery runs. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong: the recorded outcomes re-measure the thresholds, so the judgment improves with each repository's own history. |
| Feasibility     | Can we actually build this?                         | Maybe: the judge needs a door before a PRD exists, a new decision on what text may be sent, Portuguese, and thresholds nobody has measured. |

## Trade-off Insights

- **Recommended approach:** the judgment is information that starts the stage by default and never refuses anything; the criteria live in the rule, so it works without a key; everything is recorded, so the thresholds become measured.
- **Key trade-offs:** starting the stage by default against asking every time (the maintainer chose the default with a recorded dismissal: fewer interruptions, and the dismissal is the check on the judge); the judge against the agent's own criteria (both: the judge is calibrated, the criteria reach everyone); one Spec against two (the report needs recorded cases that only the first Spec produces).
- **Risks identified:** the thresholds come from a four-case pilot and are unmeasured, while every other judge question was measured before use; features are a minority of Specs, so labeled cases accumulate slowly; a raw request may carry a secret and lies outside the texts the maintainer authorized sending to the judge; the judge skips text it reads as not English; a "not indicated" verdict may be accepted without thought; the Idea stage costs an interactive session and cannot run with the user away; the glossary says an Advisory Judgment never gates, so the rule must stay a rule and never become a refusal.
- **Dissenting view:** the operator recommended one Spec with the record and no calibration report yet, because a report over a handful of cases measures noise. A second dissent: the rule alone, with the agent answering the criteria and recording the answer, reaches every adopter with no judge and no new data decision; and on 2026-11-14 the evidence that changed Spec 0392 came from the Idea stage's research, not from the judge, which was undecided.
- **Stretch goal (V2+):** judge a roadmap, comparing several features that are not yet Specs and ordering them by leverage.

## Opportunity Scan

- **Original with the record** — the judge before the PRD, the agent's criteria when there is no judge, and the record in each feature Spec. Medium effort, one Spec.
- **Simpler** — the rule alone: the agent answers three criteria and records the answer; no judge. Small effort; Feasibility Must do, Differentiation Pass.
- **More ambitious** — the original with the record, plus the calibration report and the judgment per item in the follow-up report. Large effort, two Specs; Defensibility rises to Strong.
- **Adjacent** — every feature PRD opens with a short shape check (one smaller and one larger alternative); the full Idea stage runs only when the user picks it. Small to medium effort; no judge.

**Chosen direction:** more ambitious — the maintainer's pick on 2026-11-14.

## Sub-Features

- **0393:** F1 to F5, the judgment, the rule, the criteria without the judge, the record and Portuguese.
- **0395:** F6 and F7, after 0393 is delivered and recorded cases exist.

## Integration with Existing Features

| Integration point | How |
| ----------------- | --- |
| The judge (its questions and thresholds, the Judge Log, the Jev Ceiling, the Judge Key Variables, the Run Marker) | Gains a set of questions answered before a PRD exists; every existing skip still applies and fails open. |
| The Spec routing rule and the `brainstorming`, `write-idea` and `write-prd` skills | Carry the rule, the criteria and the record. |
| The Roundfix skill's follow-up mode | Runs the judgment for each feature before authoring and collects the Idea stage's answers while the user is present. |
| `roundfix follow-up` | Second Spec: the judgment per item that is not yet a Spec. |
| A Spec folder that holds only an Idea | No tool reads it today; the record must not depend on one that ignores it. |

## Out of Scope (V1)

- **A command that refuses a PRD without an Idea** — an Advisory Judgment never gates; the rule and its record are the control.
- **Judging fixes and refactors** — the Idea stage answers a product question; their route stays `write-techspec`.
- **An agent answering the Idea stage's questions for the user** — the stage's product is shared judgment.
- **Judging inside a Run** — the judge sends nothing there, and authoring does not happen in a Run.
- **Changing a threshold automatically** — the calibration report proposes and the maintainer decides.
- **Ordering a roadmap** — kept as the stretch goal.

## Decisions

- The judgment lives in the product, as a new Spec (maintainer, 2026-11-14).
- Waste to avoid first: a feature built in the wrong shape (maintainer, 2026-11-14).
- Authority: an indicated Idea stage starts by default, and the user may dismiss it with the dismissal recorded (maintainer, 2026-11-14).
- Without the judge, the agent applies the same criteria (maintainer, 2026-11-14).
- Success: feature Specs need fewer corrections afterwards (maintainer, 2026-11-14).
- Shape: the more ambitious direction, in two Specs (maintainer, 2026-11-14).

## Open Questions

- Which questions and which thresholds? Default: the pilot's five questions and the documentation's starting thresholds (0.7 and 0.3) are a starting point only; the PRD measures them on a labeled set from this repository's history, Portuguese cases included, before fixing them, and until then the undecided band goes to the user.
- What text may be sent to the judge before a PRD exists? Default: only what the follow-up judge already sends (an item that is not yet a Spec), after the same secret scan; any other text needs the maintainer's authorization, asked when the PRD is written.
- What happens to a request the judge reads as not English if the Portuguese measurement fails? Default: the agent's criteria apply to it.
- Where does the record live when no Idea ran? Default: in the PRD.
- Which work counts as a feature? Default: both feature routes of the routing rule; fixes and refactors are excluded.
- What does an unattended cycle do when the stage is indicated and the user is away? Default: the item waits, recorded with what it waits for.
- Which records hold the corrections a Spec needed? Default: delivery retries and review rounds; the second Spec settles it.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-14-maintainer-request-and-answers.md`](references/2026-11-14-maintainer-request-and-answers.md), the maintainer's request and the answers behind every line of `## Decisions`. No ai-memory message is adopted and nobody asked for a notice.
2. [`references/2026-11-14-idea-stage-judgment-pilot.md`](references/2026-11-14-idea-stage-judgment-pilot.md), the four-case pilot, with its questions, answers and limits.

Secondbrain (read on 2026-11-14; `wiki/index.md`, then `qmd query`):

- `wiki/sources/branas-formacao-ia-processo-de-desenvolvimento.md` — a course's six-stage process starts at the PRD and has no stage that asks whether or in what shape to build. It showed that the Idea stage is not a common step, so a rule that starts it needs its own justification and its own record.

External sources (read through Exa on 2026-11-14):

- Kohavi and colleagues, "Online Controlled Experiments at Large Scale", KDD 2013, https://dl.acm.org/doi/pdf/10.1145/2487575.2488217 — the one-third figure quoted in Market Data.
- Marty Cagan, "Discovery Judgment", https://www.svpg.com/discovery-judgement/, and "The Four Big Risks", https://www.svpg.com/four-big-risks/ — discovery as a judgment of risk, with minor items skipped.
- BMAD, https://docs.bmad-method.org/plan/choose-a-planning-path/; Spec Kit, https://github.github.com/spec-kit/reference/agentic-sdd.html; Agent OS, https://github.com/buildermethods/agent-os/blob/main/commands/agent-os/shape-spec.md; Task Master, https://github.com/eyaltoledano/claude-task-master/blob/main/docs/task-structure.md — who decides whether a discovery or brief stage runs.
- Parasuraman and Manzey, "Complacency and Bias in Human Use of Automation", 2010, https://journals.sagepub.com/doi/10.1177/0018720810376055; Goddard and colleagues, "Automation bias: a systematic review", 2012, https://pmc.ncbi.nlm.nih.gov/articles/PMC3240751/ — automation bias and what reduces it.
- Mike Cohn on the Standish feature-usage figure, https://www.mountaingoatsoftware.com/agile/are-64-of-features-really-rarely-or-never-used — why that number is not used here.
- TypeSafe documentation, https://docs.typesafe.ai/primitives/noul.md and https://docs.typesafe.ai/confidence.md — a yes or no answer is a probability, 0.7 and 0.3 are starting thresholds, and the band between them can go to a person.
