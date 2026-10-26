---
type: fix
status: promoted
created: 2026-10-26
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
reason: null
---

# The Baseline setups trail the consolidated skills catalog

## Symptom

`roundfix baseline assets sync --source-dir <checkout>/setups --check`, run
on 2026-10-26 by the operator against a clean checkout of
`marcioaltoe/skills` at `8d71fade68c54693bedbd917371f0514311af157`, exits `1`
and reports `skills.setup-snapshot.drift` for all five bundled Setup
Snapshots (`go`, `go-cli-typescript-bun`, `python`, `rust`, `typescript`).
Against `65d5f365cace27c4b468db1723213723b10c870b`, the commit the snapshots
name today, the same check exits `0`. Roundfix v0.123.0 therefore ships
snapshots one upstream catalog behind.

The upstream project `default/marcioaltoe-skills` announced the change in two
ai-memory messages to this project. The operator recorded them as follows.

Message `01a17807-93df-7a91-be10-195aad68863a`, 2026-10-26 08:45 BRT:

> skill-architect and writing-for-agents removed from the catalog (canonical
> folders, links, registry, lock); no setup preset named either; skill-creator
> keeps its slug/path but is now sourced from anthropics/claude-plugins-official
> `plugins/skill-creator/skills/skill-creator`; AGENTS requires only
> skill-creator for skill authoring. "Validate Baseline references/snapshots
> before consuming."

Message `01a1781b-6c32-72a1-b192-b3ee00265a4d`, 2026-10-26 09:07 BRT:

> published by PR marcioaltoe/skills#134, squash
> `d12d7452fd740c7152ed1937f40126b00c75c70e`; onioncry updated and set to
> auto; catalog 145 entries, 127 lock entries, 317 preset references.

Upstream main then moved to `8d71fade68c54693bedbd917371f0514311af157`
(pull request #135, "update upstream skills and ai-memory release pins"),
still the newest `origin/main` on 2026-10-26 at 15:10 BRT.

## Where

- `internal/baseline/assets/setups/*.json`, written only by
  `roundfix baseline assets sync`.
- The derived files `make baseline-digests` rewrites under
  `internal/baseline/testdata/`, and this repository's Setup Manifest
  `docs/agents/setup-context.json`, whose `catalogDigest` follows the catalog.
- `internal/baseline/ai_memory_setup_test.go` and
  `internal/baseline/python_setup_test.go`, which each pin `65d5f365` as a
  separate literal.

## Expected

Every bundled Setup Snapshot names the published upstream commit and carries
the skill trees that commit publishes; the derived files and this
repository's Setup Manifest follow through their regeneration commands; and
the tests that pin the upstream revision read one constant, so the next sync
moves one line.

## Evidence

Operator measurement of 2026-10-26: a sync in a disposable worktree changed
304 lines in the five snapshots, every `ref` from `65d5f365` to `8d71fade`,
each file's `digest`, and one skill tree, `onioncry`
(`f918508d…` to `be502b02…`), present only in `typescript.json` and
`go-cli-typescript-bun.json`; no skill was added or removed. No file of this
repository names `skill-architect` or `writing-for-agents`.
