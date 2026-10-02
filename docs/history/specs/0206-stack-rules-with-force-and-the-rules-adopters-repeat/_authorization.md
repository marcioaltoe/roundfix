---
status: approved
granted: 2026-09-30
action: give every Go, Rust, CLI and TUI Baseline rule a clause with force, state the Rust error policy and the recorded-reason Go dependency preference, and promote the hand-written adopter rules that recur into Baseline clauses with their Source Baseline rows and retention disposition
consuming: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
paths:
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/bun.json
  - internal/baseline/assets/modules/go.json
  - internal/baseline/assets/modules/cli-surface.json
  - internal/baseline/assets/modules/tui-surface.json
  - internal/baseline/assets/modules/rust.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/profiles/go-cli-typescript-monorepo.json
  - internal/baseline/assets/profiles/rust-cli.json
  - internal/baseline/assets/retention/transition.legacy-typescript-bun-to-portable-v3.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/guides/rust.md
  - internal/baseline/assets/source-baselines/index.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/typescript-bun.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/plan_test.go
  - docs/agents/agent-instructions.md
  - docs/agents/skill-dispatch.md
  - docs/agents/spec-routing.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0206

On 2026-09-30 the maintainer approved the program whose ninth wave covers the
Baseline stack rules. The same day the maintainer expressly authorized editing
the Baseline source and its generated guides, answering "Autorizar os dois" to
a structured question that named `internal/baseline/assets/` and the guides
under `docs/agents/`, and said of the skills "considere autorizado a ajustar
todas as skills se necessário". For this Spec the maintainer decided that
"stdlib first" is a preference with a recorded reason, that Cobra and testify
stay banned only in this repository's own rules, the Rust typed-error policy
with force, that a hand-written rule is promoted only when it recurs in two
repositories or is tied to measured rework, that lint warnings block, and that
onioncry is validated by reading only. The governed set was measured with
`GovernedPath` on `5f182757`, through a test overlay that wrote nothing to the
repository, against the files a disposable-clone rehearsal of the four Tasks
changed.

On 2026-10-02 the scope was amended after Spec 0207 merged. A rehearsal of
task_01 at `18ef15eb`, in a disposable clone, showed the composed profile's
reference to the removed rule as the only blocker outside the grant; the same
"Autorizar os dois" answer named `internal/baseline/assets/` as a whole, and
the composed profile path, measured with the same write-free `GovernedPath`
overlay at `18ef15eb`, is governed, so it joins `paths` under that
authorization. No other newly declared path is governed.

## Why each governed path is unavoidable

- The eight modules carry the clauses this Spec adds, splits or removes
  (task_01 `core.json` and `bun.json`; task_02 `go.json`, `cli-surface.json`
  and `tui-surface.json`; task_03 `rust.json`; task_04 `spec-workflow.json`
  and `typescript.json`).
- `profiles/standard-typescript-monorepo.json` stops requiring the removed Bun
  rule (task_01), and its digest pin is rewritten by the sanctioned
  regeneration in task_01 and task_04. `profiles/rust-cli.json` requires the
  new Rust rule (task_03).
- `profiles/go-cli-typescript-monorepo.json`, the composed profile Spec 0207
  shipped after this Spec was authored, also requires the removed Bun rule;
  the catalog refuses to load while it does
  (`catalog.profile.rule.unknown`), so task_01 removes the rule from its
  `requiredRules` exactly as from the Standard TypeScript Monorepo profile.
  No other field of that profile changes, and its setup snapshot is untouched.
- The legacy retention transition names the removed Bun clause as a target;
  it is retargeted to the core clause that replaces it (task_01).
- `templates/guides/rust.md` gains its scope sentence and
  `templates/index.json` raises that template's version (task_03).
- The Source Baseline index, manifest, identity record and four corpus files
  gain one row for each new clause of a Standard TypeScript Monorepo module,
  because catalog validation refuses a required clause without its row and
  the regeneration maintains rows without creating them (task_01, task_04).
  The regeneration fills each row's offsets and digest and the identity
  record's digests.
- The five formatter goldens are rewritten by the sanctioned regeneration and
  never hand-edited (task_01, task_04).
- `internal/baseline/plan_test.go` names the guides whose module grew past the
  frozen parity record; the Go, CLI and TUI guides join that list and no
  other line changes (task_02).
- `docs/agents/agent-instructions.md`, `skill-dispatch.md`, `spec-routing.md`,
  `docs-layout.md` and `setup-context.json` are this repository's managed
  guides and Setup Manifest, rewritten by the public Managed Refresh.

## Sanctioned regeneration

```yaml
command: make baseline-digests
```

This repository's managed guides and Setup Manifest are rendered by
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`,
and a second refresh must report `File changes: 0`.

## Limits

- Exactly one clause is removed,
  `clause.bun.block-warnings-when-profile-treats-them-as-errors`, with its rule
  `rule.bun.warning-free-verification`; its replacement is declared. No other
  clause changes its identifier or enforcement level.
- No clause of the backend or frontend module changes.
- No skill, vendored or owned, is edited, and no setup snapshot changes.
- No adopter repository, onioncry or the Secondbrain is written.
- No Makefile target, lint or formatter configuration, CI workflow or
  dependency changes.
- No top-level test is renamed or removed, and no exported function signature
  changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
