# General Agent Instructions

Template for Rust CLI projects: clap-based command-line tools,
automation-facing commands, JSON/text output contracts, and multi-manifest
release workflows (crates.io and npm launchers).

## High priority

- **MANDATORY**: Use the relevant local skills before changing code, docs,
  tests, workflows, or agent instructions. Skill activation comes BEFORE any
  planning or code generation for that domain.
- **ALWAYS** prefix shell commands with `rtk` when it is available. In command
  chains, prefix each command.
- **MUST** use `rg` / `rg --files` for local code search. Use `context7` for
  external library/API docs and `exa-web-search` for broader web/source
  research. **NEVER** use web research tools to search local code.
- **MUST** run the repo's full verification gate before claiming completion.
  Any format issue, clippy warning, build failure, or test failure is
  **blocking** — zero tolerance. Note: many Rust CLI repos run tests only
  locally (CI validates conventions), so the local gate is the **only** gate.
- **HARD RULE — validation configuration is a specification**: NEVER modify lint, OnionCry, formatter, typecheck, test-runner, or verification configuration to eliminate errors or warnings. Fix the reported code, tests, paths, or architecture instead. If the required conforming change is very large and a configuration change may be justified, STOP and ask the user explicitly whether they want to change that configuration before editing it.
- **NEVER** use workarounds in production code or tests. Fix the root cause.
- **NEVER** hand-edit `Cargo.toml` dependencies. Use `rtk cargo add` /
  `rtk cargo remove`; respect the pinned `rust-version` (MSRV).
- **ABSOLUTELY FORBIDDEN**: `git reset`, `git checkout --`, `git restore`,
  `git clean`, or any command that discards working-directory changes
  **WITHOUT EXPLICIT USER PERMISSION**. These can permanently lose code.
- Agent-created branches **MUST** use the `ma/` prefix unless the repo
  documents a different human-owned prefix.
- **ALWAYS** use the AskUserQuestion tool for confirmations, clarifying questions, decision points, and any needed user interaction. If this CLI has no such tool, ask as a plain message and stop until the user answers — **NEVER** guess an answer the user can give cheaply.

## Knowledge workspace

<!-- Shared block: keep this section identical across samples/AGENTS.*.md; edit all copies together. -->

Long-lived documentation (`CONTEXT.md`, `docs/` — including `docs/specs/` and `docs/adr/` — and legacy `.compozy`/`.scratch`) lives in the central knowledge repository, mounted as a sparse checkout at `.knowledge/` (gitignored) and exposed through versioned symlinks. **ALWAYS USE** the `knowledge-workspace` skill before bootstrapping the workspace or committing any documentation change.

- **MUST** commit documentation changes inside the workspace — `git -C .knowledge add projects/<project> && git -C .knowledge commit -m "docs: ..." && git -C .knowledge push origin main` — never in the code repository. This includes spec artifacts and `task_NN.md` status flips: one piece of work often produces two commits (code here, docs in `.knowledge`).
- **NEVER** commit `.knowledge/` and **NEVER** replace the symlinks with real files.
- If `.knowledge/` or the symlinks are missing, run `scripts/knowledge-bootstrap.sh` once (`--adopt-local` when local docs already exist).

## Agent docs

Read these only when relevant to the task:

- `docs/specs/<feature-slug>/` — spec artifacts (`_idea.md`, `_prd.md`,
  `_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`); completed specs (all
  tasks done, QA passed) move to `docs/specs/_archived/`. Run
  `setup-context-driven` once if the layout is missing.
- `docs/agents/spec-routing.md` — which pipeline stages a change runs
  through (large initiative / feature / refactor-bugfix / trivial) and what
  marks a spec done
- `docs/agents/issue-tracker.md` — the local `docs/specs/` tracker conventions
- `docs/agents/triage-labels.md` — label mapping for issue triage skills
- `docs/agents/domain.md` — how agents consume `CONTEXT.md` and ADRs
- `CONTEXT.md` — project vocabulary, command concepts, domain rules, and
  product decisions
- `docs/adr/` — architectural decisions; flag conflicts before overriding them
- `docs/release.md` — the release choreography, when the repo publishes

**ALWAYS** use canonical terms from `CONTEXT.md` in command names, help text,
issue titles, test names, and user-facing explanations. If the right term is
missing, call out the gap instead of inventing new language.

## Skill dispatch

Before editing, identify the task domain and **activate every matching skill**:

- **Feature discovery or product idea**: Use `brainstorming`; product-level
  ideas go through `write-idea` (scored by `business-analyst`, debated by
  `council`, challenged by `the-fool`)
- **PRD, tech spec, or task breakdown**: Use `write-prd`, `write-techspec`,
  `write-tasks`. Pick the pipeline entry point per
  `docs/agents/spec-routing.md`: large/fuzzy initiatives start at
  `write-idea`, standard features at `write-prd`, refactors and bug fixes at
  `write-techspec` (it mints the spec folder with a minimal `_prd.md`);
  trivial one-line changes skip the pipeline entirely
- **Executing spec tasks**: Use `implement-task` (one task) or `implement-spec`
  (the whole graph in dependency order)
- **Final QA and QA test execution for a completed spec**: Use `qa-gate`
  exclusively. It owns QA planning, execution, evidence, findings, verdict,
  and the dated report. Keep `testing-boss` and the Rust testing skills limited
  to writing and maintaining implementation tests; they do not replace the QA
  workflow. On QA pass, the spec is complete and `archive-spec` runs
  automatically (merge/release is a separate user-driven step, never the
  archive gate)
- **MANDATORY** for CLI behavior, flags, stdout/stderr, exit codes, JSON
  output, dry-run behavior, non-interactive mode, or introspection:
  `agentic-cli-design`
- **ALWAYS USE** `rust` and `rust-cli` before Rust command behavior, error
  handling, output contracts, integration tests, or Cargo workflows
- **ALWAYS USE** `clap-rust` before touching clap parser code: derive
  patterns, subcommands, flag groups, value parsers, completions
- **Cutting a release**: Use `cut-release` — version sync across every
  manifest (Cargo.toml plus npm launcher packages), publish dry-runs, tag,
  watch the workflow to completion
- **Tests, fixtures, golden files, integration tests**: Use `testing-boss`
  plus the relevant Rust/CLI skill
- **Implementation**: Use `coding-guidelines`
- **Bug fix or failing test**: Use `no-workarounds` plus `systematic-debugging`
- **Docs, PRDs, ADRs, issues, PR descriptions**: Use `tech-writer`
- **Commits or PR titles**: Use `conventional-commits`
- **Completion claim**: Use `evidence-gate`
- **Session handoff**: Use `handoff`
- **OnionCry checks**: Use `onioncry` only when running OnionCry against
  architecture boundaries or modifying the OnionCry project itself. Do NOT
  load OnionCry for ordinary Rust CLI work.

## CLI behavior

- Design commands for humans **and** agents: deterministic output,
  non-interactive flags, stable exit codes, machine-readable modes.
- **MUST** keep stdout for requested command output only. Diagnostics,
  progress, and warnings go to stderr.
- Command names, flags, JSON/SARIF fields, text output relied on by tests,
  and exit-code contracts are **public API** — never change them casually.
  Versioned output formats keep their version markers accurate.
- Help text **MUST** be concise, truthful, and backed by implemented behavior.
- Errors **MUST** name the failed operation and the next useful action when
  one is known.

## Rust conventions

- Prefer explicit domain types over stringly typed command plumbing.
- Typed errors via `thiserror` (plus `miette` when the repo renders terminal
  diagnostics). **NEVER** `unwrap`/`expect`/`panic!` on user-facing paths.
- **MUST** build structured output with `serde`/`serde_json`. Hand-built JSON
  is rejected.
- Keep parsing, domain execution, and output rendering separable enough to
  test without spawning the binary for every case.
- **Integration tests drive the real binary**: `assert_cmd` + `predicates` in
  `tempfile::TempDir` sandboxes, sharing helpers through `tests/support/`.
  Assert stdout/stderr text, JSON via `serde_json::Value`, and exit codes.
  User-facing contracts are tested through the binary; pure parsing and
  domain logic get unit tests in-module.
- Schema-validate machine-readable outputs (JSON Schema fixtures for
  SARIF-style formats) instead of eyeballing shape.

## Verification

Use the repo's declared verification command when present, commonly
`make verify`. If no aggregate command exists, run at least:

```bash
cargo fmt --all -- --check
cargo check --all-targets --all-features
cargo clippy --all-targets --all-features -- -D warnings
cargo test --all-features
```

Also run any documented build, snapshot, integration, or release check that
applies to the touched files. **Skipping any verification check invalidates
the completion claim.**

## Git and delivery

- **MUST** check `git status --short` before staging; keep unrelated user
  changes out of your diff.
- Use `conventional-commits` for commits and PR titles (check `cog.toml` —
  some Rust repos enforce a **closed scope list**).
- **NEVER** rewrite unrelated files or format the whole repo unless asked.
- PR bodies summarize changes, call out risk, and list validation commands run.
- **NEVER** tag a release by hand outside the `cut-release` choreography: a
  version mismatch across Cargo.toml and npm manifests fails the publish
  workflow late.

## Anti-patterns (immediate rejection)

1. `unwrap`/`expect` on user-facing paths, or hand-built JSON strings
2. Marking a spec task `completed` without fresh verification evidence
3. Tracking progress in `_tasks.md` — status lives only in each `task_NN.md`
4. Asking for confirmation before running spec tasks — invocation is the
   authorization
5. Changing exit codes, flags, or output formats without treating it as a
   breaking API change
6. Tagging before the publish dry-runs pass — the tag is the trigger, not the
   test
