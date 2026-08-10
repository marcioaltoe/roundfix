---
type: refactor # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The roundfix skill carries every command in one file

## Opportunity

`.agents/skills/roundfix/SKILL.md` is **2,172 lines across 29 sections** — six
times the next largest owned skill, and nine times the median.

| Skill | Lines |
| --- | ---: |
| `roundfix` | **2,172** |
| `setup-context-driven` | 362 |
| `qa-gate` | 359 |
| `write-tasks` | 240 |
| `archive-spec` | 240 |
| `write-prd` | 215 |
| `implement-task` | 166 |
| `write-techspec` | 124 |

It covers release planning, CodeRabbit resolution, runtime readiness, Task Graph
execution, Run monitoring, storage reclamation, Spec archiving, the bounded
Batch resolution contract, and the configuration reference — in one file that an
Agent loads whole to answer any one of them.

## Value

Every dispatch pays for all nine subjects. An Agent resolving a review comment
loads the release-planning reference, the retention policy and the whole
configuration schema alongside what it needs. That is the cost the instruction
architecture concept calls out: the global layer should carry invariants, and
specific knowledge should be discovered by links rather than loaded up front.

It also makes the file hard to change safely. The review-request rule added on
2026-08-09 had to be placed by finding an anchor in a 2,000-line document, and
the contract test that guards it asserts seven literal strings whose only
location is that same file.

## Shape

Non-binding. The straightforward move is a `references/` directory beside
`SKILL.md`, with the skill keeping what routes and each reference keeping one
subject — the shape `write-techspec` and `write-tasks` already use for their
templates.

Worth settling in the same work: which subjects are genuinely separate. Release
planning, storage and retention, and the configuration reference look
independent. The Batch resolution contract and Run monitoring may not be, since
an Agent inside a Batch needs both.

Worth measuring rather than assuming: how much of the 2,172 lines a typical
dispatch actually uses. The answer decides whether this is a real cost or a
tidiness preference.
