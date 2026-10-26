---
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
status: active
created: 2026-10-26
surfaces: [backend]
---

# The Baseline setups follow the consolidated skills catalog

Roundfix v0.123.0 bundles five Setup Snapshots that name upstream commit
`65d5f365` of `marcioaltoe/skills`, while upstream has since consolidated its
catalog and published `8d71fade`; the snapshot check against the published
commit reports drift in all five. After this Spec every bundled Setup
Snapshot names `8d71fade68c54693bedbd917371f0514311af157` and carries the
`onioncry` tree that commit publishes, the derived Baseline files and this
repository's Setup Manifest follow through their regeneration commands, and
the tests that pin the upstream revision read one shared constant. No skill
joins or leaves a snapshot, and no module, profile, guide or installed skill
of this repository changes.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-26-the-baseline-setups-trail-the-consolidated-skills-catalog.md`](references/2026-10-26-the-baseline-setups-trail-the-consolidated-skills-catalog.md),
   `type: fix`, from two ai-memory messages of `default/marcioaltoe-skills`
   (`01a17807`, 2026-10-26 08:45 BRT, and `01a1781b`, 09:07 BRT) and the
   operator's check of 2026-10-26.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. This Spec is authored against `86cb6b73` (v0.123.0) and
delivers after Spec 0319, which touches no file this Spec declares. The
import fetches commit `8d71fade68c54693bedbd917371f0514311af157` read-only
from the local checkout `~/dev/skills`, which holds it on 2026-10-26; task_01
stops before any write when it does not.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the snapshots keep their ids and schema. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added, and no Verification or QA row reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0191, "A setup snapshot follows
  its upstream by name": a Setup Snapshot mirrors one upstream list and the
  asset sync is its only writer, so this Spec moves the snapshots only
  through the sync; ADR-0204: the composed `go-cli-typescript-bun` snapshot
  stays the union the sync writes from `go` and `typescript`; ADR-0241: a
  repository-owned profile takes `onioncry`'s contract from the embedded
  snapshots, which agree on the new tree; ADR-0322: every upstream snapshot,
  `python` included, names one commit. Recorded as not applying: ADR-0206
  (no upstream list or skill is renamed or removed in this import), ADR-0246
  (no Retired Skill changes; `skill-architect` and `writing-for-agents` left
  upstream without ever reaching a snapshot) and ADR-0219 (no profile or
  profile draft changes). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AQ decision of 2026-10-26 quoted in the authorization record, and
  the Governed Paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0320-the-baseline-setups-follow-the-consolidated-skills-catalog/_authorization.md`;
  bounded files:
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `docs/agents/setup-context.json`.

## Goals

- Every bundled Setup Snapshot, and every upstream skill it pins, names
  `8d71fade`; the snapshot check against that commit exits `0`.
- The `typescript` and `go-cli-typescript-bun` snapshots carry the `onioncry`
  tree `be502b02…` the commit publishes; every other skill tree and every
  snapshot's membership stay as they are.
- The derived Baseline files and this repository's Setup Manifest are
  current, and a second managed refresh reports no file change.
- No test pins `65d5f365`; the upstream revision the tests expect is one
  constant.
- This repository's installed skills, lock, modules, profiles and rendered
  guides stay byte-identical.

## Core Features

1. The snapshot import. The asset sync rewrites the five snapshots from a
   checkout of the published commit, and its `--check` then exits `0`.
2. The derived files. `make baseline-digests` rewrites the catalog files and
   the plan characterization goldens under `internal/baseline/testdata/`,
   and the managed refresh rewrites the `catalogDigest` of
   `docs/agents/setup-context.json`.
3. One revision constant. A new test file declares the upstream revision and
   the published `onioncry` tree and asserts the tree in both TypeScript
   snapshots; the two tests that pinned `65d5f365` read the constant.

## User Experience

No command, flag, message or guide changes. After the release, the next
`roundfix baseline update` of every adopter rewrites one managed entry of
its Setup Manifest, the `catalogDigest`. An adopter of the
`standard-typescript-monorepo` or `go-cli-typescript-monorepo` profile, whose
`typescript` module requires and dispatches `onioncry`, and that holds the
skill at its earlier tree sees
the skill under `skills.differing` in the refresh and
`DR-SKILL-DIFFERS-FROM-SNAPSHOT` in `roundfix doctor`; the refresh keeps the
installed copy, and `roundfix baseline skills restore --profile <id> --skill
onioncry` takes the new one. The only upstream change in that skill is its
`metadata.source` line, which now names `https://github.com/marcioaltoe/onioncry`;
its instructions are unchanged. Go, Rust and Python adopters see only the
`catalogDigest` line.

## Non-Goals / Out of Scope

- Restoring or re-pinning this repository's installed skills or
  `skills-lock.json`: the `go` snapshot's skill trees are unchanged, and the
  lock records each skill's install revision.
- Any change for `skill-creator`, `skill-architect` or `writing-for-agents`:
  no upstream preset at `65d5f365` or `8d71fade` lists them, and no module,
  dispatch rule, Repository Skill Set entry or file of this repository names
  them.
- The upstream mirrors of Roundfix-owned skills (`roundfix`, `write-idea`,
  `write-prd`, `write-techspec`, `setup-context-driven`) that changed
  upstream: the snapshots take those skills from source `repo`, so Roundfix's
  own copies stay authoritative.
- Editing any ADR, user guide page, skill or Baseline module.

## Success Metrics

1. Success Metric: each of the five snapshots names
   `8d71fade68c54693bedbd917371f0514311af157` for itself and for every
   `github` skill, holds the same skill names in the same order as at the
   Delivery Base, and differs from the Delivery Base in exactly one skill
   tree, `onioncry`, in `typescript.json` and `go-cli-typescript-bun.json`.
2. Success Metric: `go test ./internal/baseline ./skills` passes,
   `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree` among them, and
   `TestTheSetupManifestIsCurrent` passes.
3. Success Metric: in a clone of the audited head, `make baseline-digests`
   reports `"changed":false` and the managed refresh reports `File changes: 0`.
4. Success Metric: no test file names `65d5f365`, and every Go file outside
   `internal/baseline/*_test.go`, every module, profile and template, every
   installed skill, `skills-lock.json` and every rendered guide other than
   `docs/agents/setup-context.json` is byte-identical to the Delivery Base.
5. Success Metric: after release, a TypeScript adopter's next managed refresh
   rewrites its `catalogDigest` and lists `onioncry` under `skills.differing`.

## Acceptance evidence

The outside-evidence row rests on statements this Spec did not write: the two
ai-memory messages of `default/marcioaltoe-skills` quoted in the adopted
Backlog Entry under `references/`, received on 2026-10-26.

- `01a17807`: "skill-architect and writing-for-agents removed from the
  catalog (canonical folders, links, registry, lock); no setup preset named
  either; skill-creator keeps its slug/path but is now sourced from
  anthropics/claude-plugins-official".
- `01a1781b`: "onioncry updated and set to auto".

The Claude plugins repository publishes the skill at
`plugins/skill-creator/skills/skill-creator/SKILL.md` of
`anthropics/claude-plugins-official` (Exa search, 2026-10-26); no preset
lists it, so it reaches no snapshot.

## Unreachable Acceptance

- criterion: Success Metric 5, after release, a TypeScript adopter's next
  managed refresh rewrites its `catalogDigest` and lists `onioncry` under
  `skills.differing`.
  reason: it needs the released binary in another repository that holds
  `onioncry` at its earlier tree, which no Verification or QA row may read.
  satisfied-by: the operator's next managed refresh of a TypeScript adopter
  after the release, reported to this project through ai-memory.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Prototype in a scratch worktree of `86cb6b73` with a clean upstream clone
  at `8d71fade` whose origin is `https://github.com/marcioaltoe/skills.git`:
  the sync changed 304 lines in the five snapshots (293 skill refs, five
  snapshot refs, five `digest` lines, two `onioncry` tree lines) and its
  `--check` then exited `0`; `make baseline-digests` rewrote seven files under
  `internal/baseline/testdata/`; the managed refresh rewrote only the
  `catalogDigest` of `docs/agents/setup-context.json`, and a second refresh
  reported `File changes: 0`. Before the test edit, `go test
  ./internal/baseline` failed only in
  `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
  `TestEveryUpstreamSetupFollowsThePythonPresetRevision` and
  `TestThePythonSetupIsThePublishedUpstreamPreset`, on the revision literal;
  with the shared constant and the new test, `./internal/baseline`,
  `./skills`, `./internal/config` and the documentation contracts passed.
  `./internal/cli` passed apart from `TestGuidanceCompositionJourney/rust-cli`,
  which failed because the prototype ran with a fresh `HOME` that has no
  rustup default toolchain, an environment effect unrelated to the change.
- Upstream diff `65d5f365..8d71fade`: `setups/` is unchanged; under the
  skills the presets list, only `skills/06-review-repair/onioncry/SKILL.md`
  changed, one line, `source: https://github.com/marcioaltoe/skills` to
  `source: https://github.com/marcioaltoe/onioncry`. The other changed files
  are upstream mirrors of Roundfix-owned skills, the removed `skill-architect`
  and `writing-for-agents`, `skill-creator` and repository metadata.
- `rg -n "skill-creator|skill-architect|writing-for-agents"` over this
  repository outside `docs/history/` finds nothing; no Go file reads a
  skill's `metadata.source` or names `onioncry`.
- Secondbrain, read only: read `wiki/index.md`; `qmd query "baseline setup
  snapshot skills catalog onioncry skill-creator" --all --files --min-score
  0.3` returned the upstream study
  `projects/skills/mirror/docs/research/2026-10-26-agent-skills-and-tools-study.md`,
  the upstream `onioncry` skill, this repository's Baseline guide and the
  Archive Record of Spec 0208; none records a decision against the import.
- Exa: the published `skill-creator` location quoted under Acceptance
  evidence.

## Glossary

None.

## Decisions

- 2026-10-26 (about 15:00 BRT): the maintainer wrote "Sim, pode fazer a
  higiene." and "Vamos seguir a sua sugestão de sequencia para a spec e os
  demais itens. Talvez seja importante usar o jev para validar.", approving
  the hygiene step that checks the Baseline snapshots against
  `marcioaltoe/skills` `d12d7452` or later.

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 1 briefing for Spec 0320 of 2026-10-26 and the characterization above.

Taken from those sources, with no further interview:

- Target `8d71fade`, the newest published upstream `main` at authoring time,
  rather than `d12d7452`: the snapshots and the skills they pin must name one
  commit, and `8d71fade` contains `d12d7452`.
- Keep this repository's installed skills and lock as they are: its `go`
  snapshot changes no skill tree.
- Replace the two revision literals with one test constant, as the
  agent-instructions clause on pinned values asks, and pin the `onioncry`
  tree in a new test so a sync that skipped the tree is caught.
- No ADR: the import follows ADR-0191 and decides nothing new.

## Open Questions

None.
