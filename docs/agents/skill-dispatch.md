<!-- setup-context-driven:begin id=guide.skill-dispatch version=0.0.1 -->

# Skill dispatch

This setup-owned guide is generated from the active modules. Activate every
matching skill before governed work; repository-authored dispatch extensions
may add stricter triggers.

- **mandatory**: Activate every matching required skill before governed work. When one skill has distinct active-module triggers, retain and follow each trigger.

- **mandatory**: Use the governing `conventional-commits` skill before staging changes, writing commit messages, or preparing pull request titles.

- **mandatory**: Use the governing `github-pr-workflow` skill before preparing, opening, updating, or handing off a pull request.

Exact activation bundles:

- `trigger.delivery`: Preparing commits or pull request delivery.
  - `bundle.delivery`: `conventional-commits`, `github-pr-workflow`
- `trigger.qa`: Running final Spec QA.
  - `bundle.qa`: `qa-gate`, `evidence-gate`

Individual skill triggers:

- `agentic-cli-design`:
  - `trigger.cli-surface.agentic-cli-design`: Changing CLI flags, streams, exit codes, JSON, dry-run, non-interactive, or introspection behavior.
- `archive-spec`:
  - `trigger.context-workflow.archive-spec`: Archiving a completed and QA-passed Spec.
- `brainstorming`:
  - `trigger.context-workflow.brainstorming`: Starting creative feature or behavior design before implementation.
- `bubbletea`:
  - `trigger.tui-surface.bubbletea`: Building or changing Bubble Tea models, commands, messages, or component composition.
- `business-analyst`:
  - `trigger.context-workflow.business-analyst`: Evaluating product viability, options, KPIs, or business trade-offs.
- `coding-guidelines`:
  - `trigger.core.coding-guidelines`: Writing, modifying, refactoring, or reviewing implementation code.
- `context7`:
  - `trigger.core.context7`: Consulting authoritative current documentation for a library, framework, runtime, or toolchain.
- `conventional-commits`:
  - `trigger.core.conventional-commits`: Staging changes, writing commit messages, or preparing pull request titles.
- `council`:
  - `trigger.context-workflow.council`: Debating a high-impact product or architecture decision through multiple advisors.
- `domain-modeling`:
  - `trigger.context-workflow.domain-modeling`: Defining or changing domain vocabulary, ownership, or bounded-context relationships.
- `evidence-gate`:
  - `trigger.core.evidence-gate`: Making any completion, readiness, or handoff claim.
- `github-pr-workflow`:
  - `trigger.core.github-pr-workflow`: Preparing, opening, updating, or handing off a pull request.
- `golang-cli`:
  - `trigger.go.golang-cli`: Changing Go command structure, flags, package layout, version output, or command tests.
- `golang-concurrency`:
  - `trigger.go.golang-concurrency`: Using goroutines, channels, locks, worker pools, or other concurrent Go behavior.
- `golang-context`:
  - `trigger.go.golang-context`: Designing Go IO boundaries, cancellation, timeouts, deadlines, or request-scoped values.
- `golang-dependency-management`:
  - `trigger.go.golang-dependency-management`: Adding, upgrading, or removing Go dependencies, auditing vulnerabilities, resolving version conflicts, or configuring automated updates.
- `golang-error-handling`:
  - `trigger.go.golang-error-handling`: Creating, wrapping, matching, logging, or exposing Go errors.
- `golang-lint`:
  - `trigger.go.golang-lint`: Changing Go lint configuration or adding a lint suppression.
- `golang-safety`:
  - `trigger.go.golang-safety`: Writing or reviewing Go for nil safety, append aliasing, concurrent map access, numeric conversion overflow, defer in loops, or defensive copying.
- `golang-structs-interfaces`:
  - `trigger.go.golang-structs-interfaces`: Designing Go types and interfaces — composition, embedding, type assertions and switches, struct field tags, or pointer versus value receivers.
- `golang-testing`:
  - `trigger.go.golang-testing`: Writing or reviewing Go tests, fixtures, fuzz tests, or integration tests.
- `grill-with-docs`:
  - `trigger.context-workflow.grill-with-docs`: Stress-testing a plan or decision against repository documents.
- `grilling`:
  - `trigger.context-workflow.grilling`: Stress-testing an idea or plan through direct questions.
- `handoff`:
  - `trigger.core.handoff`: Recording resumable session state for another Agent or maintainer.
- `implement-spec`:
  - `trigger.autonomous-work.implement-spec`: Delegating a complete autonomous Spec Task Graph.
  - `trigger.context-workflow.implement-spec`: Executing every pending Task in a Spec Task Graph.
- `implement-task`:
  - `trigger.autonomous-work.implement-task`: Delegating one autonomous Task to an ACP Runtime.
  - `trigger.context-workflow.implement-task`: Executing exactly one assigned Spec Task.
- `knowledge-workspace`:
  - `trigger.context-workflow.knowledge-workspace`: Consulting the configured knowledge workspace for relevant prior context.
- `no-workarounds`:
  - `trigger.core.no-workarounds`: Fixing defects, regressions, or verification failures at their root cause.
- `qa-gate`:
  - `trigger.context-workflow.qa-gate`: Running final Spec QA after implementation Tasks complete.
- `review`:
  - `trigger.core.review`: Reviewing a change against repository standards and its originating contract.
- `setup-context-driven`:
  - `trigger.context-workflow.setup-context-driven`: Auditing or applying the Context-Driven Baseline for first adoption, or refreshing an adopted repository with `roundfix baseline update`.
- `systematic-debugging`:
  - `trigger.core.systematic-debugging`: Diagnosing a bug, failing check, or unexpected behavior before proposing a fix.
- `tech-writer`:
  - `trigger.core.tech-writer`: Writing or revising technical documentation, Specs, findings, or delivery notes.
- `testing-boss`:
  - `trigger.core.testing-boss`: Writing, changing, or reviewing tests and fixtures.
- `the-fool`:
  - `trigger.context-workflow.the-fool`: Challenging a proposal through a pre-mortem, red-team, or evidence audit.
- `tui-design`:
  - `trigger.tui-surface.tui-design`: Designing terminal layout, navigation, keyboard, mouse, color, or accessibility behavior.
- `write-idea`:
  - `trigger.context-workflow.write-idea`: Expanding a product-level idea into a researched opportunity artifact.
- `write-prd`:
  - `trigger.context-workflow.write-prd`: Defining product requirements and user outcomes.
- `write-tasks`:
  - `trigger.context-workflow.write-tasks`: Decomposing an approved Spec into a dependency-ordered Task Graph.
- `write-techspec`:
  - `trigger.context-workflow.write-techspec`: Settling implementation architecture and technical contracts.

<!-- setup-context-driven:end id=guide.skill-dispatch -->
