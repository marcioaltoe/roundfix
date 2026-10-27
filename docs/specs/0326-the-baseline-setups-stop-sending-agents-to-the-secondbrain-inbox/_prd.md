---
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
status: active
created: 2026-10-26
surfaces: [backend]
---

# The Baseline setups stop sending agents to the Secondbrain inbox

Roundfix v0.124.0 bundles five Setup Snapshots that pin the
`knowledge-workspace` skill at upstream commit `8d71fade` of
`marcioaltoe/skills`, a tree that still tells agents to create files in the
Secondbrain inbox the maintainer retired on 2026-10-25, while the Baseline's
own `secondbrain` guide forbids any Secondbrain write. Upstream fixed the
skill in `b56fc5ffcad1ab366fbc03839a6ce8c662401178`. After this Spec every
bundled Setup Snapshot names that commit and carries its
`knowledge-workspace` tree, this repository's installed copy and lock entry
follow it through the restore command, the derived Baseline files and the
Setup Manifest follow through their regeneration commands, and the revision
test constant moves with a new test that pins the corrected tree. No skill
joins or leaves a snapshot, and no module, profile, guide or other installed
skill changes.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
from two ai-memory messages (see `references/_index.md`):

1. [`references/2026-10-26-the-knowledge-workspace-skill-still-writes-to-the-secondbrain-inbox.md`](references/2026-10-26-the-knowledge-workspace-skill-still-writes-to-the-secondbrain-inbox.md),
   `type: fix`, from `01a179b6` of `default/gesttione-solutions-pantheon`
   (2026-10-26 16:36 BRT) and item 3 of `01a1795f` of
   `default/marcioaltoe-skills` (2026-10-26 15:01 BRT).

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. This Spec is authored against `f06afdd7` (v0.124.0 plus
the 0321 QA amendment) and delivers after Spec 0322, which regenerates the
same derived catalog files and the Setup Manifest for its module edits; this
Spec regenerates them again from whatever its base holds, so it pins none of
their values. The import fetches commit
`b56fc5ffcad1ab366fbc03839a6ce8c662401178` read-only from the operator's
scratch clone `/Users/marcio/.roundfix-operator/scratch/spec-0326/skills`,
which holds it on 2026-10-26; task_01 stops before any write when it does not.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the snapshots keep their ids and schema. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added, and no Verification or QA row reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0191, "A setup snapshot follows
  its upstream by name": the asset sync is the only writer of a Setup
  Snapshot, so all five move through it; ADR-0204: the composed
  `go-cli-typescript-bun` snapshot stays the union the sync writes from `go`
  and `typescript`, with its `composed` source and no ref; ADR-0221 and
  ADR-0264: a required external skill is held to its Setup Snapshot, and a
  copy that differs is restored only deliberately through
  `roundfix baseline skills restore`, which is how this repository and every
  adopter take the corrected skill; ADR-0241: a repository-owned profile takes
  `knowledge-workspace`'s contract from the embedded snapshots, which agree on
  the new tree; ADR-0322: every upstream snapshot, `python` included, names
  one commit; ADR-0330: the Secondbrain is read-only and ai-memory carries
  cross-project messages, the decision the corrected skill now states;
  ADR-0060: Source Baselines are frozen, project-agnostic inputs, so the
  legacy corpus that still describes the old inbox is left as it is;
  ADR-0224: this repository holds its upstream skills to the snapshot and the
  lock, so a change that moves the `go` snapshot restores this repository's
  skill in the same change. Recorded as not applying: ADR-0206 (no upstream
  list or skill is renamed or removed), ADR-0219 (no profile or profile draft
  changes), ADR-0246 (no Retired Skill changes) and ADR-0338 (the follow-up
  mode's handoff and message handling is untouched; it cites ADR-0330 only
  for the ai-memory channel). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AQ decisions of 2026-10-26 quoted in the authorization record, and
  the Governed Paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox/_authorization.md`;
  bounded files:
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `docs/agents/setup-context.json`,
  `.agents/skills/knowledge-workspace/SKILL.md`.

## Goals

- Every bundled Setup Snapshot, and every upstream skill it pins, names
  `b56fc5ff`; the snapshot check against that commit exits `0`.
- All five snapshots carry the `knowledge-workspace` tree `490af37d…` that
  commit publishes; every other skill tree and every snapshot's membership
  stay as they are.
- This repository's installed `knowledge-workspace` and its lock entry match
  the new snapshot, so the repository follows its own release.
- The derived Baseline files and this repository's Setup Manifest are
  current, and a second managed refresh reports no file change.
- No test pins `8d71fade`; the upstream revision the tests expect stays one
  constant.

## Core Features

1. The snapshot import. The asset sync rewrites the five snapshots from a
   checkout of the published commit, and its `--check` then exits `0`.
2. This repository's skill. `roundfix baseline skills restore --profile
   go-cli-tui --skill knowledge-workspace` with `--source-dir` rewrites
   `.agents/skills/knowledge-workspace/SKILL.md` and its `skills-lock.json`
   entry from the same checkout.
3. The derived files. `make baseline-digests` rewrites the catalog files and
   the plan characterization goldens under `internal/baseline/testdata/`,
   and the managed refresh rewrites the `catalogDigest` of
   `docs/agents/setup-context.json`.
4. The revision tests. The shared constant in
   `internal/baseline/upstream_setup_revision_test.go` moves to `b56fc5ff`,
   and a new test in the same file asserts the `knowledge-workspace` ref and
   tree in all five snapshots.

## User Experience

No command, flag, message or guide changes. After the release, the next
`roundfix baseline update` of every adopter rewrites one managed entry of its
Setup Manifest, the `catalogDigest`, and, because every profile's
`context-workflow` module requires and dispatches `knowledge-workspace`,
lists the installed skill under `skills.differing`; `roundfix doctor` reports
`DR-SKILL-DIFFERS-FROM-SNAPSHOT`. The refresh keeps the installed copy
(ADR-0264). Each adopter then runs
`roundfix baseline skills restore --profile <id> --skill knowledge-workspace`,
reviews the preview (exit `3` with a Plan Digest), and applies it with
`--confirm-plan <digest>`. The `<id>` is the adopter's own profile:
`go-cli-tui`, `go-cli-typescript-monorepo`, `rust-cli`,
`standard-typescript-monorepo`, `python`, or a repository-owned profile,
whose contract comes from the five snapshots that now agree. The restored
skill sends research to `default/marcioaltoe-secondbrain` with
`memory_message_send` and drops the inbox carve-out; its consultation rules
are unchanged.

## Non-Goals / Out of Scope

- Editing the legacy Source Baseline
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/`,
  whose corpus still describes the old inbox: it is a frozen, digest-verified
  input that the managed refresh reads only to account an old adoption's
  clauses (ADR-0060), never a file it renders.
- The historical citation of an inbox entry in a comment of this
  repository's `.roundfixrc.yml`: it names a source, it instructs nothing.
- Restoring any other installed skill or re-pinning other lock entries: no
  other skill tree moves.
- Any Baseline module, profile, guide, ADR or user guide page: the
  `secondbrain` and `core` modules already state the read-only rule.

## Success Metrics

1. Success Metric: each of the four snapshots with a `github` source names
   `b56fc5ffcad1ab366fbc03839a6ce8c662401178` as its `source.ref`; the
   composed `go-cli-typescript-bun.json` keeps its `composed` source
   (`setups: [go, typescript]`, no ref) byte-identical to its value before
   the Task; all five name that commit as the `ref` of every `github` skill,
   keep the same skill names in the same order, and differ from before the
   Task in exactly one skill tree, `knowledge-workspace`'s, now `490af37d…`.
2. Success Metric: `go test ./internal/baseline ./skills` passes,
   `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`,
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` and
   `TestCheckRepositoryMatchesRealRepository` among them, and
   `TestTheSetupManifestIsCurrent` passes.
3. Success Metric: in a clone of the audited head, `make baseline-digests`
   reports `"changed":false` and the managed refresh reports `File changes: 0`.
4. Success Metric: no test file names `8d71fade`, the installed
   `knowledge-workspace` reads `version: 0.6.2` and names neither
   `dev/secondbrain/inbox/` nor `create-only carve-out`, and the Task's
   commits change no file outside its declared set.
5. Success Metric: after release, an adopter's next managed refresh lists
   `knowledge-workspace` under `skills.differing`, and its restore applies
   the `b56fc5ff` tree.

## Acceptance evidence

The outside-evidence row rests on text this Spec did not write: the upstream
diff `8d71fade..b56fc5ff` of
`skills/00-setup/knowledge-workspace/SKILL.md`, whose new lines the adopted
Backlog Entry under `references/` quotes, and the two ai-memory messages it
records (`01a179b6`, `01a1795f`). The QA compares the restored skill in the
audited head with those quoted lines, offline.

## Unreachable Acceptance

- criterion: Success Metric 5, after release, an adopter's next managed
  refresh lists `knowledge-workspace` under `skills.differing`, and its
  restore applies the `b56fc5ff` tree.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the release notice to `default/gesttione-solutions-pantheon`
  and that project's refresh and restore after the release, reported to this
  project through ai-memory.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Prototype in a scratch clone of `f06afdd7` with a clean upstream clone at
  `b56fc5ff` whose origin is `https://github.com/marcioaltoe/skills.git`: the
  sync's `--check` exited `1` before (drift in all five files); the sync
  changed 307 lines per side in the five snapshots (refs, five `digest`
  lines, five `knowledge-workspace` tree lines) and its `--check` then exited
  `0`; `make baseline-digests` rewrote seven files under
  `internal/baseline/testdata/`; the restore preview listed exactly two
  changes (`.agents/skills/knowledge-workspace/SKILL.md` and the
  `skills-lock.json` entry) and applied them; the managed refresh rewrote
  only the `catalogDigest`, and a second refresh reported `File changes: 0`.
  Before the test edit, `./internal/baseline` failed only in the four tests
  that read the revision constant; after it, `./internal/baseline`,
  `./skills`, `./internal/cli`, `./internal/config` and the documentation
  contracts passed.
- Upstream diff `8d71fade..b56fc5ff`: one file,
  `skills/00-setup/knowledge-workspace/SKILL.md`, 4 insertions and 4
  deletions: the version line and the "Feeding business knowledge",
  "Project coordination and continuity" and "Territories" bullets. The
  skill's description, triggers and consultation rules are unchanged, so
  the two dispatch triggers of `context-workflow`
  (`trigger.context-workflow.knowledge-workspace`,
  `trigger.context-workflow.knowledge-workspace-secondbrain-authoring`) still
  describe it, and no module clause, rendered guide or verification step
  reads the changed text.
- The Baseline's rendered guidance already agrees: `secondbrain.json`
  forbids every Secondbrain write, `inbox/` included, and `core.json` routes
  cross-project observations through `memory_message_send` (ADR-0330).
- Secondbrain, read only: read `wiki/index.md`; `qmd query
  "knowledge-workspace inbox retired ai-memory messages secondbrain" --all
  --files --min-score 0.3` returned the upstream skill's mirror, this
  repository's `CONTEXT.md`, ADR-0330 and the `roundfix` follow-up
  reference; none records a decision against the import.
- Exa: published vendored-skill tools pin each skill to a commit plus a
  content hash and move a pin only on an explicit update ("The only command
  that moves a pin", `luisalima/skills-lock`), the same deliberate model as
  ADR-0264's restore.

## Glossary

None.

## Decisions

- 2026-10-26 (about 15:00 BRT): the maintainer wrote "Sim, pode fazer a
  higiene." and "Vamos seguir a sua sugestão de sequencia para a spec e os
  demais itens. Talvez seja importante usar o jev para validar.", the cycle
  AQ approval this batch runs under.

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 4 briefing for Spec 0326 of 2026-10-26 and the characterization above.

Taken from those sources, with no further interview:

- Target `b56fc5ff`, upstream `main` at authoring time (re-measured
  2026-10-26 ~20:30 BRT).
- Restore this repository's `knowledge-workspace` in the same Task: its
  `go` snapshot tree changes, and
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` holds the
  installed copy to it.
- Extend Spec 0320's test file rather than add a parallel one.
- Leave the legacy Source Baseline untouched (ADR-0060).
- No ADR: the import follows ADR-0191 and ADR-0264 and decides nothing new.

## Open Questions

None.
