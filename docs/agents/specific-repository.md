# Agent instructions

This repository is Roundfix: a local-first Go CLI and future daemon that picks
up work items (today: PR review issues; next: spec task graphs), resolves them
through the user's selected ACP runtime, and pushes only when nothing
unresolved remains. Stdlib `flag` dispatch and a Bubble Tea v2 TUI.

  Any format failure, test failure, or build failure is **blocking** — zero
  tolerance. CI validates PR titles only; the local gate is the ONLY gate.
- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- Agent-created branches **MUST** use the `ma/` prefix.
- **HARD RULE — roundfix skill sync**: before opening any PR, confirm the
  roundfix skill still matches the shipped CLI behavior; a PR that changes
  CLI behavior ships the skill update too. Contract:
  `docs/agents/skill-dispatch.md`.
- **HARD RULE — skill ownership**: repo-owned authorial workflow skills may
  be adapted locally; every other skill is upstream-managed and **MUST NOT**
  be modified here. Ownership split: `docs/agents/skill-dispatch.md`.
- **HARD RULE — protected-tooling commit choreography**: the express
  authorization record with its exact bounded paths, and any **prerequisite**
  fix — one repairing something already red before the Task — are each their
  own commit landing **before** the authorized Task commit, in either relative
  order. An authorization written into the Spec artifacts normally predates
  every implementation commit, which is the healthy case. A **consequent** fix,
  which only becomes necessary because the authorized change made something
  else stale, is its own commit landing **after** the Task commit; it cannot
  precede the cause that created it. Either kind folded into the Task commit
  fails the tooling-authority gate. Prefer no consequent fix at all: a Task's
  declared scope should include the tests its own change invalidates.
- **HARD RULE — sanctioned digest regeneration**: after an expressly
  authorized Roundfix-owned Skill or Baseline module edit, run
  `make baseline-digests`. Every derived pin rewritten by that command is
  deterministic fallout of the authorized source edit and needs no separate
  express authorization. A hand-edited pin value remains an unauthorized
  mutation.
- **HARD RULE — release planning**: release work starts with the read-only
  `roundfix release plan` before changelog, version, tag, push, package, asset,
  or GitHub Release mutation. A generic release request authorizes only a
  conclusive patch plan; minor, major, version-zero breaking, and manual
  classification outcomes require the decisions in
  `docs/user-guide/release-runbook.md`.
- **HARD RULE — never let a pipe hide a gate's exit status**: a pipeline exits
  with its last command's status, so `make verify | tail` reports the pager's
  success and `&&` proceeds over a red gate. Run the gate on its own, capture
  `$?`, or redirect to a file and read it — this is how a commit landed on a
  failing gate on 2026-07-29, and it is the same defect the Makefile carried
  in `find … | sort`.
- **HARD RULE — an assertion reads the constant it means**: a test that copies
  a pinned version, digest, or identifier as a literal stops testing the day a
  legitimate change moves it, sometimes silently — a fixture that mutates
  `version: 0.0.1` mutates nothing once the value is `0.0.2`, and the test
  still passes. Reference the exported or package constant. When a value must
  be duplicated, change every occurrence in the same commit: `grep` for it
  first, because fixing one of three is the most repeated defect in this
  repository's history.
- **HARD RULE — durable knowledge flows upstream only**: Specs are downstream
  results of the CONTEXT-driven workflow, never sources it depends on — an
  archived Spec may be deleted at any time, so durable knowledge a Spec
  produced moves to its semantic owner (`CONTEXT.md`, an accepted ADR, an
  agent guide, or `docs/references/`) before or at archive. `CONTEXT.md` and
  the agent guides reference accepted ADRs and **NEVER** reference
  `docs/specs/` or `docs/findings/` content; findings are dated reports that
  become Specs, not reference material.
- Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
  `internal/...` (`internal/cli/` owns parsing, output, and exit behavior;
  `internal/app/` holds app metadata)

### Skill dispatch

Skill triggers, ownership, and the Roundfix skill-sync contract live in
`docs/agents/skill-dispatch.md`.

## Verification

The full gate is `make verify` (fmt-check + test + `roundfix skills check` +
build) and it **MUST** pass 100% before any completion claim. For the smallest
relevant gate while iterating:

The Makefile defaults `GOCACHE` to the ignored repository-local `.gocache`
directory only when the environment does not set it; an explicit `GOCACHE`
always wins. After an authorized Roundfix-owned Skill or Baseline module edit,
run `make baseline-digests` before `make verify`. The sanctioned command owns
every derived pin update; never transcribe a digest by hand.

```bash
rtk gofmt -w <changed-go-files>
rtk go test ./...
rtk go run ./cmd/roundfix --help
```

If concurrency changed, also run:

```bash
rtk go test -race ./...
```

The `0.0.x` series is patch-only here: every release on it bumps the patch
component regardless of the conventional-commit mix. `roundfix release plan`
maps `feat` commits to a minor minimum and refuses a manual patch
classification below it, so a `0.0.x` release is a recorded maintainer
decision rather than a plan proposal until the tool learns this policy or the
project moves to `0.1.0`.

- Use `conventional-commits` for commits and PR titles (check `cog.toml`).
- Commit and PR titles are unscoped Conventional Commits subjects here
  (`cog.toml` sets `scopes = []`).
- PR bodies summarize changes, call out risk, and list validation commands run.
- **HARD RULE — a hand-opened code pull request asks for its own review**:
  `.coderabbit.yaml` sets `auto_review.enabled: false`, so no review happens
  unless something asks. Roundfix asks on its own Final Push
  (`review_source.request_review: true`), which covers the watch and resolve
  loops and nothing else. A pull request opened directly — by `gh pr create` or
  any other route — gets no review at all, and its CodeRabbit check reports
  `Review skipped: automatic reviews are disabled`, which reads like a pass.
  Whenever the pull request changes code, post `@coderabbitai review` as a
  comment immediately after opening it, and treat that comment as part of
  opening the pull request rather than a follow-up. Documentation-only pull
  requests may skip it. Three pull requests merged unreviewed on 2026-08-07
  because this rule did not exist: the tool had been taught to ask and the
  agents had not.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers
