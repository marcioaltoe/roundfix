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
- `golang-error-handling`:
  - `trigger.go.golang-error-handling`: Creating, wrapping, matching, logging, or exposing Go errors.
- `golang-lint`:
  - `trigger.go.golang-lint`: Changing Go lint configuration or adding a lint suppression.
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
  - `trigger.context-workflow.setup-context-driven`: Auditing, applying, or refreshing the Context-Driven Baseline.
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

<!-- roundfix:repository-rule:begin id=rule.fa2871b863891fa9aa1200e64da20a493ff91b274d1e49f1d5cec763b060b54f -->
## Repository dispatch extensions

- **Broader web or source research**: Use `exa-web-search` when current
  authoritative library documentation is insufficient. Search from multiple
  angles and validate conclusions against primary sources. Never use it to
  discover or infer local repository code.
- **Skill optimization**: Use `autoresearch` only to improve a repo-owned
  skill through its evaluation loop. Never use it to mutate an
  upstream-managed skill locally.

### Skill ownership

Repo-owned workflow skills are plain files under `.agents/skills/` and are not
listed in `skills-lock.json`; they may be adapted to Roundfix. Skills declared
in `skills-lock.json` are upstream-managed and must not be edited locally.
Propose needed changes to their upstream source.

### Roundfix skill synchronization

Before opening a PR, confirm `.agents/skills/roundfix/SKILL.md` matches the
shipped CLI commands, flags, output formats, exit codes, and Batch semantics.
When CLI behavior changes, update that canonical skill in the same PR and run
`make skills-sync` so the embedded `skills/roundfix/` copy matches it.

The skill's `metadata.version` follows the released CLI version. Do not keep a
detached machine-global copy that can shadow this repository's canonical
contract; use a symlink to `skills/roundfix/` when a global install is needed.

<!-- roundfix:repository-rule:end id=rule.fa2871b863891fa9aa1200e64da20a493ff91b274d1e49f1d5cec763b060b54f -->
<!-- roundfix:repository-rule:begin id=rule.b7ee6fb14a4e620c351e803f52d4dc4752636e9319c1445ae8d01a9939b6ba62 -->
- **Feature discovery or product idea**: Use `brainstorming`; product-level
  ideas go through `write-idea` (scored by `business-analyst`, debated by
  `council`, challenged by `the-fool`)

<!-- roundfix:repository-rule:end id=rule.b7ee6fb14a4e620c351e803f52d4dc4752636e9319c1445ae8d01a9939b6ba62 -->

<!-- roundfix:repository-rule:begin id=rule.44abd68d3af68ac77a9668ca6df1a4fd2fd03807b0c18570614051e14bd4d2b0 -->
- **PRD, tech spec, or task breakdown**: Use `write-prd`, `write-techspec`,
  `write-tasks`

<!-- roundfix:repository-rule:end id=rule.44abd68d3af68ac77a9668ca6df1a4fd2fd03807b0c18570614051e14bd4d2b0 -->

<!-- roundfix:repository-rule:begin id=rule.91383959f8b5a1f4bc15edc36a1ac18bacea397adb5a3848e77d6e3ca751a8fe -->
- **Executing spec tasks**: Use `implement-task` (one task) or `implement-spec`
  (the whole graph in dependency order)

<!-- roundfix:repository-rule:end id=rule.91383959f8b5a1f4bc15edc36a1ac18bacea397adb5a3848e77d6e3ca751a8fe -->

<!-- roundfix:repository-rule:begin id=rule.a709b50ca54a0862333e30497aa73778928b0a3f8133c87e5c812fa878895ece -->
- **Final QA of a completed spec**: Use `qa-gate`; archive after release with
  `archive-spec`

<!-- roundfix:repository-rule:end id=rule.a709b50ca54a0862333e30497aa73778928b0a3f8133c87e5c812fa878895ece -->

<!-- roundfix:repository-rule:begin id=rule.b62e0e1064e62aedee90356df08c6ca6517c05a95cf404a6dce68377fa83fa02 -->
- **MANDATORY** for CLI behavior, flags, stdout/stderr, exit codes, JSON
  output, dry-run behavior, non-interactive mode, or introspection:
  `agentic-cli-design`

<!-- roundfix:repository-rule:end id=rule.b62e0e1064e62aedee90356df08c6ca6517c05a95cf404a6dce68377fa83fa02 -->

<!-- roundfix:repository-rule:begin id=rule.257cded33420ff9e2c0a1233b04f3a9bb4d3e69a927ceaa9fb2bc82d132dd42e -->
- **ALWAYS USE** `golang-cli` before writing Go command behavior, package
  layout, version output, or command tests. CLI style is stdlib `flag.FlagSet`
  dispatch with a `Run() int` exit-code contract — **no Cobra**.

<!-- roundfix:repository-rule:end id=rule.257cded33420ff9e2c0a1233b04f3a9bb4d3e69a927ceaa9fb2bc82d132dd42e -->

<!-- roundfix:repository-rule:begin id=rule.7223a2001fea50552284adbccc0fcddf59c01bbf5bb943b7f94c1c78bfca68aa -->
- **ALWAYS USE** `golang-error-handling` for error paths: `%w` wrapping,
  `errors.Is`/`As`, sentinels

<!-- roundfix:repository-rule:end id=rule.7223a2001fea50552284adbccc0fcddf59c01bbf5bb943b7f94c1c78bfca68aa -->

<!-- roundfix:repository-rule:begin id=rule.5ecd1c28260d6faad79e50f91d615b54ad54bbdcf6f533666df04facc90206f6 -->
- **ALWAYS USE** `golang-concurrency` before goroutines, channels, worker
  pools, or anything with leak/race exposure

<!-- roundfix:repository-rule:end id=rule.5ecd1c28260d6faad79e50f91d615b54ad54bbdcf6f533666df04facc90206f6 -->

<!-- roundfix:repository-rule:begin id=rule.6bb7ff6acfe19f3b76058b07f09510505b96cd30c590f1c4850508b27e5b90e0 -->
- **ALWAYS USE** `golang-context` for context propagation, cancellation, and
  timeouts

<!-- roundfix:repository-rule:end id=rule.6bb7ff6acfe19f3b76058b07f09510505b96cd30c590f1c4850508b27e5b90e0 -->

<!-- roundfix:repository-rule:begin id=rule.ad9c13f96bd7ec59d54989a8cd065338c06ead025129dadfbefef3d418b29bd9 -->
- **Lint config or nolint**: Use `golang-lint` (golangci-lint, vet,
  staticcheck discipline)

<!-- roundfix:repository-rule:end id=rule.ad9c13f96bd7ec59d54989a8cd065338c06ead025129dadfbefef3d418b29bd9 -->

<!-- roundfix:repository-rule:begin id=rule.742b018445608f7db9454c22779b77c2a00b659964d8bf900e390bcd63b6733b -->
- **Tests, fixtures, golden files, integration tests**: Use `golang-testing`
  plus `testing-boss`

<!-- roundfix:repository-rule:end id=rule.742b018445608f7db9454c22779b77c2a00b659964d8bf900e390bcd63b6733b -->

<!-- roundfix:repository-rule:begin id=rule.65cb06de682d619743c6549e45005b020c46272b8c614e8ce90c18afbfb02760 -->
- **Bubble Tea or Lip Gloss TUI work**: Use `bubbletea` and `tui-design`

<!-- roundfix:repository-rule:end id=rule.65cb06de682d619743c6549e45005b020c46272b8c614e8ce90c18afbfb02760 -->

<!-- roundfix:repository-rule:begin id=rule.e6998fd3ffda4286cad26995d5c58da80bb8c93ce665fa2fd1def9748fc7fd18 -->
- **Implementation**: Use `coding-guidelines`

<!-- roundfix:repository-rule:end id=rule.e6998fd3ffda4286cad26995d5c58da80bb8c93ce665fa2fd1def9748fc7fd18 -->

<!-- roundfix:repository-rule:begin id=rule.eb416a2e487d856a0649451b1fc9580e42abfa08df38ab7bf23d37a1642f7b1b -->
- **Bug fix or failing test**: Use `no-workarounds` plus `systematic-debugging`

<!-- roundfix:repository-rule:end id=rule.eb416a2e487d856a0649451b1fc9580e42abfa08df38ab7bf23d37a1642f7b1b -->

<!-- roundfix:repository-rule:begin id=rule.d0ab371252128a34bf2b6000031026be1019930f88d714bed6f7147db803ebdb -->
- **Docs, PRDs, ADRs, issues, PR descriptions**: Use `tech-writer`

<!-- roundfix:repository-rule:end id=rule.d0ab371252128a34bf2b6000031026be1019930f88d714bed6f7147db803ebdb -->

<!-- roundfix:repository-rule:begin id=rule.ec539cc69cd3996d889c1a1b493c8b77ce4f96a71402c09f568f84ec6571925e -->
- **Commits or PR titles**: Use `conventional-commits`

<!-- roundfix:repository-rule:end id=rule.ec539cc69cd3996d889c1a1b493c8b77ce4f96a71402c09f568f84ec6571925e -->

<!-- roundfix:repository-rule:begin id=rule.b33d4d55a40b58baea327824fde5611396391e677119020377419b6dadb1ad4e -->
- **Completion claim**: Use `evidence-gate`

<!-- roundfix:repository-rule:end id=rule.b33d4d55a40b58baea327824fde5611396391e677119020377419b6dadb1ad4e -->

<!-- roundfix:repository-rule:begin id=rule.0815cdacadf1822cebd71fd4c406d891bfd1a70c682cbf21f01e5a7f11955375 -->
- **Session handoff**: Use `handoff`

<!-- roundfix:repository-rule:end id=rule.0815cdacadf1822cebd71fd4c406d891bfd1a70c682cbf21f01e5a7f11955375 -->

<!-- roundfix:repository-rule:begin id=rule.5f91916a7ee99b20b60f63c5855c3af4bf775950c037b12c9dfa3d7e0e0bb094 -->
- **Roundfix dogfooding or assigned-Batch contract checks**: Use `roundfix`
  when driving Roundfix against an Open Pull Request or validating the Batch
  resolution contract


<!-- roundfix:repository-rule:end id=rule.5f91916a7ee99b20b60f63c5855c3af4bf775950c037b12c9dfa3d7e0e0bb094 -->
