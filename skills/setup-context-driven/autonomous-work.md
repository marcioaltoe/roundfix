# Autonomous work: Supervisor and implementation runtimes

Who does what when this repo works autonomously. One Supervisor session orchestrates;
implementation is delegated to an ACP Runtime, normally through a Roundfix Run. The split binds
every autonomous session, interactive or unattended: Supervisor capacity is reserved for
orchestration and judgment, and operational work goes to implementation runtimes. The repo's
agent instructions (`AGENTS.md`/`CLAUDE.md`) enforce this as a hard rule. Command contracts
(flags, outcomes, monitoring) live in the `roundfix` skill; this doc fixes the role split,
runtime routing, and how the Supervisor authors Specs.

## Roles

| Role | Model | Job |
| --- | --- | --- |
| Supervisor | Supervising Claude Code session | Author Specs, launch and monitor Runs, integrate outcomes, run `qa-gate`, archive, boundary commits, route work to a runtime |
| Default implementer | Codex `gpt-5.5` with `xhigh` Default Reasoning Effort | Every Spec Task and Review Issue Batch not routed below |
| Design implementer | Claude with `opus` at `high` or `xhigh` Default Reasoning Effort | Tasks dominated by design, UI, UX, or frontend work |

## The Supervisor does not implement

Delegation is mandatory, not a preference: Supervisor capacity is limited and reserved for the
work that needs its judgment. The Supervisor session:

- routes changes through the spec pipeline (`docs/agents/spec-routing.md`) and authors the
  planning artifacts
- launches Runs detached, monitors them to a terminal outcome, integrates, runs `qa-gate`, and
  archives on pass
- picks the implementation runtime per Run (rules below)
- keeps its own edits boundary-scoped: spec and doc fixes, Run recovery, boundary commits — it
  does not implement Spec Tasks itself while an ACP Runtime can do the work

Writing feature code, tests, or any other operational implementation directly in the Supervisor
session violates the hard rule — delegate it, even when doing it inline looks faster.

Fable is an Agent Model label in the Claude catalog, not the supervising role. Never point
`--agent` at a Supervisor-backed runtime.

## Authoring Specs

The Supervisor authors every Spec through the pipeline skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`); which stages a change needs lives in
`docs/agents/spec-routing.md`. Automated authoring keeps these behaviors:

- **Ask before minting a Spec.** When new work surfaces — findings triage, an idea, a bug —
  propose the Spec (slug, scope, pipeline route) and ask via AskUserQuestion whether to create
  it with `_prd.md`, `_techspec.md`, and the Task Graph before writing any folder. The
  authorization boundary sits here: running an approved Spec's Tasks needs no confirmation
  (invocation is the authorization), but creating a new Spec always does.
- **Record decisions as ADRs.** Product and technical decisions settled during `write-prd` and
  `write-techspec` land in `docs/adr/` in the same authoring pass — never only in conversation.
- **Enrich CONTEXT.md as terms appear.** When a Spec introduces a domain concept the glossary
  lacks, add the canonical term with its avoid-list instead of inventing synonyms; ask when the
  term or its meaning needs the user's call.
- **Hold shared-doc edits to Run boundaries.** CONTEXT.md, ADRs, AGENTS.md, and spec files are
  edited only while no Run is Active — uncommitted edits to files an Active Run also touches
  block its integration and leave the Run IntegrationPending.
- **Author Tasks for the runtime that executes them.** Task Verification commands run verbatim
  in a bare worktree: they must be fully self-contained — real paths, no `<placeholders>`, no
  assumed build, link, or install state, and the repo's real build flags. Verifications that
  fit the repo's standard test and verify commands work; install-heavy verifications do not.
- **Sequence the queue.** Keep one explicit, dependency-and-risk-ordered implementation
  sequence across pending Specs. New Specs slot into that queue when approved; the Supervisor
  advances it one Run at a time and re-plans the order at each Run boundary.

## Default implementer: Codex gpt-5.5 xhigh

```bash
roundfix implement --spec <slug> --agent codex --qa --detach
```

Review-surface Runs (`resolve`, `watch`) use the same default runtime. Roundfix owns the
effective Agent Model and Default Reasoning Effort:

1. The repo's Roundfix Project Config pins `defaults.agent: codex`.
2. The repo's Roundfix Project Config pins `runtimes.codex.model: gpt-5.5` and
   `runtimes.codex.reasoning_effort: xhigh`.
3. acpx (`~/.acpx/config.json`) maps the `codex` agent to the local `codex-acp` adapter.

Do not set `defaults.model`; Roundfix ignores it and prints a deprecation warning. Use
`--model` and `--reasoning-effort` together only for a one-Run exception.

## Design implementer: Claude Code with Opus 4.8

Claude Code outperforms Codex on design-heavy work. Route a Run to it when the Task Graph is
dominated by:

- design — visual and interaction design decisions
- UI and UX — layout, navigation, information hierarchy, feedback states, user-facing copy
- frontend — <this repo's design surface: its TUI, web frontend, or both>

```bash
roundfix implement --spec <slug> --agent claude --model opus --reasoning-effort high --qa --detach
```

The `claude` runtime launches through `claude-agent-acp`. Roundfix forwards both the Agent Model
and Default Reasoning Effort. Use the one-Run flags above for a design-heavy exception, or pin
`runtimes.claude.model` and `runtimes.claude.reasoning_effort` in Project Config when a whole
repository should use that route.

## Routing rules

- The routing unit is the Run: one Run drives one Agent. A Spec that mixes backend and
  design/UI work either runs on the default Codex runtime or is sliced at `write-tasks` time so
  the design surface gets its own Spec.
- Work delegated outside a Run (a one-off fix handed to a coding CLI) follows the same routing:
  Codex `gpt-5.5 xhigh` by default, Claude Code Opus 4.8 for design/UI/UX/frontend.
- When in doubt, use Codex. The Claude route exists for work where design quality is the point,
  not as a general alternative.

## What this does not change

- Verification is runtime-independent: every Task passes its Verification, the repo's
  verification gate binds every completion claim, and `qa-gate` runs after the last Task
  regardless of runtime.
- The spec workflow, issue-tracker conventions, and skill dispatch in the repo's agent
  instructions bind every runtime equally.
