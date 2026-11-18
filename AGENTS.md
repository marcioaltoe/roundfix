<!-- setup-context-driven:begin id=root.core version=0.0.1 -->

### Instruction hierarchy

Apply active guidance in this order. A narrower guide may add constraints for its concern but cannot weaken a universal Normative Clause or confirmed project decision.

1. **Universal instructions**
2. **Context and documentation**
3. **Spec workflow**
4. **Autonomous work**
5. **Stack guidance**
6. **Surface guidance**
7. **Optional knowledge sources**
8. **Repository-Specific Normative Rules**

### Universal instructions

- Follow the mandatory setup-owned baseline in `docs/agents/agent-instructions.md`.
- Activate skills through `docs/agents/skill-dispatch.md` before governed work.

<!-- setup-context-driven:end id=root.core -->
<!-- setup-context-driven:begin id=root.context-workflow version=0.0.1 -->

### CONTEXT-driven workflow

- Domain rules are mandatory: `docs/agents/domain.md`. Read `docs/agents/docs-layout.md`, whose rules then apply, before creating, changing, moving, or retiring `CONTEXT.md` or a document under `docs/`, and before a test or build step reads one.

<!-- setup-context-driven:end id=root.context-workflow -->
<!-- setup-context-driven:begin id=root.spec-workflow version=0.0.1 -->

### Local Spec workflow

- Route and execute local Specs through `docs/agents/spec-routing.md` and `docs/agents/issue-tracker.md`.
- In Claude Code, a skill's `model` and `effort` header applies only when a
  person types the skill (`/skill`); a skill an Agent activates runs on the
  session's model (anthropics/claude-code#98898, measured 2026-11-14). To
  hold an Authoring Phase to its declared model there, hand the phase to its
  Author Subagent, `.claude/agents/roundfix-<phase>-author.md`, which
  `roundfix skills install --target project` generates from
  `skills.authoring.<harness>`; when that model's usage limit ends the
  subagent, run `roundfix phase fallback` and restart it on the model it prints.
  This bullet does not apply to Codex or OpenCode, whose models come from
  the Agent Selection Profiles.

<!-- setup-context-driven:end id=root.spec-workflow -->
<!-- setup-context-driven:begin id=root.autonomous-work version=0.0.1 -->

### Autonomous work

- Supervisor and ACP Runtime rules are mandatory: `docs/agents/autonomous-work.md`.

<!-- setup-context-driven:end id=root.autonomous-work -->
<!-- setup-context-driven:begin id=root.go version=0.0.1 -->

### Go

- Go implementation rules are mandatory: `docs/agents/go.md`.

<!-- setup-context-driven:end id=root.go -->
<!-- setup-context-driven:begin id=root.cli-surface version=0.0.1 -->

### CLI surface

- CLI public-contract rules are mandatory: `docs/agents/cli.md`.

<!-- setup-context-driven:end id=root.cli-surface -->
<!-- setup-context-driven:begin id=root.tui-surface version=0.0.1 -->

### TUI surface

- TUI interaction rules are mandatory: `docs/agents/tui.md`.

<!-- setup-context-driven:end id=root.tui-surface -->
<!-- setup-context-driven:begin id=root.secondbrain version=0.0.1 -->

### Secondbrain

- Read `docs/agents/secondbrain.md` before consulting or writing the Secondbrain or authoring an Idea, PRD, or TechSpec; a Run session, which cannot reach it, skips it.

<!-- setup-context-driven:end id=root.secondbrain -->
<!-- setup-context-driven:begin id=root.repository-extension version=0.0.1 -->

### Repository-owned rules

- Read `docs/agents/specific-repository.md` for project-specific hard rules,
  and write every rule or variation of this repository there, never in this
  file or in a Baseline guide.

<!-- setup-context-driven:end id=root.repository-extension -->
