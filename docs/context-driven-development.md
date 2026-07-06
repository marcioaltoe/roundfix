# CONTEXT-driven development

CONTEXT-driven development is how this repository plans and builds: every change
flows through a pipeline of planning artifacts — idea, PRD, tech spec, tasks —
that live as local markdown under `docs/specs/<slug>/`, grounded in a shared
vocabulary (`CONTEXT.md`) and recorded decisions (`docs/adr/`). Downstream stages
read the artifacts, not the conversation, so any fresh agent session continues
from the files alone.

The method is adapted from Matt Pocock's **skills** work; see
[Source and attribution](#source-and-attribution).

## Why "CONTEXT-driven"

The name points at the artifact at the center of the method: `CONTEXT.md`, the
project glossary. Code, docs, prompts, task names, and TUI copy all draw their
vocabulary from it, so an agent and a human describe the same thing the same way.
Shared, written-down language is what keeps artifacts consistent across sessions
and models — it is the context that drives the work, rather than a single long
conversation an agent cannot reload.

Three artifacts carry that context:

- **`CONTEXT.md`** — the domain vocabulary. New terms are added as Specs
  introduce them.
- **`docs/adr/`** — Architecture Decision Records. Product and technical
  decisions are recorded here in the same authoring pass, never only in chat.
- **`docs/specs/<slug>/`** — the planning artifacts for one change, the only home
  of that work's plan.

## The pipeline

```text
write-idea -> write-prd -> write-techspec -> write-tasks -> implement -> qa-gate -> roundfix archive
```

Each stage reads and writes `docs/specs/<slug>/` and produces one artifact:

| Stage | Produces | Altitude |
| --- | --- | --- |
| `write-idea` | `_idea.md` — the opportunity, scored and debated | Why |
| `write-prd` | `_prd.md` — user stories and acceptance criteria | What and why |
| `write-techspec` | `_techspec.md` — contracts, data models, failure modes | How |
| `write-tasks` | `_tasks.md` + `task_NN.md` — a dependency-ordered Task Graph | Work units |
| `implement` | Completed Tasks, each with verification evidence | Execution |
| `qa-gate` | `qa/` report validating the feature against the PRD | Verdict |
| `roundfix archive` | The Spec stamped and moved to `docs/specs/_archived/` | Record |

Not every change runs the full pipeline. The entry point depends on the change —
a large initiative starts at `write-idea`, a standard feature at `write-prd`, a
refactor or bug fix at `write-techspec`, and a trivial fix skips the spec folder
entirely. The routing rules live in
[`docs/agents/spec-routing.md`](agents/spec-routing.md); every route converges on
`write-tasks`, so implementation always executes from a Task Graph rather than an
ad-hoc plan.

## How Roundfix executes it

Roundfix is the runtime for the implementation half of the pipeline. `roundfix
implement --spec <slug>` executes the Task Graph as one Run — Tasks in dependency
order, each gated by its own Verification commands — and `roundfix archive`
closes the loop after `qa-gate` passes. The operational flow is documented in the
[usage guide](usage.md); the autonomous role split (an orchestrator authors
Specs, an ACP Runtime implements them) is in
[`docs/agents/autonomous-work.md`](agents/autonomous-work.md).

The result is that planning and execution share one artifact contract: the same
`docs/specs/<slug>/` files a human reads are what the agent implements from and
what QA checks against.

## Source and attribution

This method is adapted from **Matt Pocock's skills** —
<https://github.com/mattpocock/skills> — described there as "Skills for Real
Engineers. Straight from my .claude directory." Pocock's approach organizes
agent work as composable skills built on software-engineering fundamentals:

- **Alignment before execution** — resolve gaps between developer and agent
  before writing code.
- **Domain-driven language** — build shared terminology through `CONTEXT.md` and
  ADRs to keep agents aligned. This is the origin of the "CONTEXT" in
  CONTEXT-driven development.
- **Tight feedback loops** — test-driven development and fast verification to
  validate quality.

Roundfix adapts those ideas into this repository's pipeline (`write-idea`,
`write-prd`, `write-techspec`, `write-tasks`, `implement-spec`/`implement-task`,
`qa-gate`, `archive-spec`) and ships the authorial skill bundle in the binary.
The skills are maintained as an adaptation at
[`marcioaltoe/skills`](https://github.com/marcioaltoe/skills) and pinned through
`skills-lock.json`; the wording, task contract, and Roundfix integration are this
repository's own. Credit for the underlying method belongs to Matt Pocock.
