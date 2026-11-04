---
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: active
created: 2026-11-04
surfaces: [backend, docs]
type: fix
---

# The vendored catalog skills follow the catalog's main through a pinned bump

This repository vendors 33 skills from the skills catalog `marcioaltoe/skills`,
and its `skills-lock.json` pins them at seven different refs, six of them the
branch `main`, while `systematic-debugging` still carries five files the
catalog removed. After this Spec one Makefile target, `make
skills-third-party-bump`, moves the 33 vendored folders, their lock entries and
the five upstream Setup Snapshots to one full commit of the catalog's `main`
through the existing asset sync and confirmed restore, and a second run at the
same commit changes nothing. The Spec runs it once at `c8b2e85`, so every
vendored folder equals the catalog's tree and `systematic-debugging` loses the
five files. The setup tests read the expected revision from the lock instead
of a hand-pinned constant, and a documentation contract fails on any vendored
folder that differs from the tree its Catalog Revision publishes.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
from an ai-memory message (see `references/_index.md`):

1. [`references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md`](references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md),
   `type: fix`, reported by marcioaltoe-skills, message `01a1a608` of
   `default/marcioaltoe-skills` (2026-11-04 07:08 BRT), reproduced verbatim.

The same sender's message `01a1a5f7` (an adopter's `baseline update` keeps an
old `ref` on the Roundfix-owned entries) is the source of Spec 0349, which
owns that defect; this Spec does not adopt it. At `6a0ed904` the only open
Backlog Entry, about the waiting declaration of the backlog contract, belongs
to Spec 0346, which has since adopted it, and `docs/findings/` holds none;
nothing open shares this context.

## Prerequisites

None in `requires`. The plan was measured on `6a0ed904` (v0.130.1),
rebased on `364d1250` (the 0345-0348 plans and the 0348 delivery), and
delivers last in the Cycle AX queue, after Specs 0345, 0348, 0347, 0346 and
0349. Files it may share with them:

- `docs/agents/setup-context.json`, which 0346 and 0349 may also refresh: this
  Spec regenerates it with the managed refresh from whatever its base holds
  and pins no value.
- The derived Baseline files under `internal/baseline/testdata/`, which 0346
  regenerates for its module edit: this Spec regenerates them with `make
  baseline-digests` and pins no value.
- `CONTEXT.md`, where 0346 to 0349 may add terms: this Spec adds two entries
  after **Vendored Skill** and touches no other line.
- Spec 0349 owns `internal/baseline/owned_skills_lock.go`,
  `internal/baseline/skills_restore*.go` and the Doctor skills check; this
  Spec edits none of them and only runs `roundfix baseline skills restore`.
  It relies on that command writing each restored catalog entry's `ref` as
  the snapshot's commit, as it does at `6a0ed904`; if 0349 changes that for
  catalog entries, task_02 stops and reports.

task_02 reads the catalog commit `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`
from the operator's clone
`/Users/marcio/.roundfix-operator/scratch/spec-0350/skills`, which holds it
on 2026-11-04; the Task stops before any write when it does not.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the lock, the snapshots and the manifest keep their schemas.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Only the new Makefile target may reach the network,
  to resolve and fetch one catalog commit, and every Verification and QA row
  runs it against a local catalog. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0360, "The vendored catalog skills
  follow the catalog through one pinned commit", the decision this Spec
  implements; ADR-0191: the asset sync stays the only writer of a Setup
  Snapshot, so the bump moves all five through it; ADR-0204: the composed
  `go-cli-typescript-bun` snapshot stays the union the sync writes; ADR-0221
  and ADR-0264: a required external skill is held to its Setup Snapshot and
  replaced only through the confirmed restore, which the bump runs; ADR-0224,
  which ADR-0360 amends: this repository holds its skills to the snapshot and
  the lock, so the bump moves both; ADR-0241: a repository-owned profile takes
  its contracts from snapshots that agree; ADR-0322: every upstream snapshot,
  `python` included, names one commit; ADR-0326: a test that reads
  `.agents/skills` is a `docscontract` file declared `//verify:always`.
  Recorded as not applying: ADR-0206 (no skill is renamed or leaves a
  preset), ADR-0219 (no profile draft changes; it cites ADR-0204 only for
  the composed snapshot), ADR-0350 (the Retired Skill set is unchanged). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  Cycle AX answers of 2026-11-04 quoted in the authorization record, whose
  Rules answer releases the Makefile when the design needs it; bounded files:
  `.agents/skills/systematic-debugging/CREATION-LOG.md`,
  `.agents/skills/systematic-debugging/test-academic.md`,
  `.agents/skills/systematic-debugging/test-pressure-1.md`,
  `.agents/skills/systematic-debugging/test-pressure-2.md`,
  `.agents/skills/systematic-debugging/test-pressure-3.md`, `Makefile`,
  `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump/_authorization.md`.

## Goals

- One command moves the vendored catalog skills, their lock entries and the
  upstream Setup Snapshots to one full commit of the catalog's `main`, and
  repeats with no change at the same commit.
- After this Spec, the 33 catalog entries name
  `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, each vendored folder equals the
  catalog's tree there, and `systematic-debugging` has none of the five
  removed files.
- No Roundfix-owned skill enters the lock, and the bump refuses one that does.
- No test hand-pins the catalog revision or a tree that the next bump moves.
- A local edit of a vendored folder fails a contract even when its lock entry
  is re-hashed.

## Core Features

1. The Catalog Skill Bump. `make skills-third-party-bump` resolves the
   catalog's `main` (or takes `SKILLS_REF`), fetches that commit from
   `SKILLS_CATALOG`, refuses what it must not pin, runs the asset sync, the
   confirmed restore of every catalog entry, `make baseline-digests` and the
   managed refresh, and checks the result. `make skills-update` is retired.
2. The bump at `c8b2e85`. The Spec runs the target once from the operator's
   clone and commits what it writes.
3. Tests that follow the bump. The setup revision tests read the Catalog
   Revision from `skills-lock.json`; the three tree tests refuse the tree each
   earlier fix replaced; a new documentation contract compares every vendored
   folder with the tree its snapshots pin.
4. The rule in words. `CONTEXT.md` defines **Catalog Revision** and **Catalog
   Skill Bump**, and the repository rules name the bump as the only writer of
   the vendored skills.

## User Experience

The maintainer runs `make skills-third-party-bump` at each cycle, reviews the
diff and commits it; with no argument it follows the catalog's `main`. The
target prints one final line naming how many catalog skills and which commit
it pinned, and every refusal names its cause. Adopters see the moved Setup
Snapshots after the release: their next `roundfix baseline update` rewrites
the Setup Manifest's `catalogDigest` and lists each moved skill they hold
(`systematic-debugging` in every profile, `ui-ux-pro-max` in the TypeScript
profiles, `developing-with-streamlit` in Python) under `skills.differing`,
and `roundfix doctor` names it with `DR-SKILL-DIFFERS-FROM-SNAPSHOT` without a
warning; each adopter restores it deliberately with
`roundfix baseline skills restore --profile <id> --skill <name>`.

## Non-Goals / Out of Scope

- The ref an adopter's `baseline update` writes for Roundfix-owned entries
  and the Doctor warning about incoherent entries: Spec 0349.
- Adding, removing or renaming a vendored skill or a setup member: a bump
  that meets one stops, and that move is a Spec (ADR-0360).
- A CI job or scheduled routine that runs the bump: the maintainer runs it
  each cycle, by decision.
- The file mode of `find-polluter.sh`: neither the lock's `computedHash` nor
  the snapshots' tree digest reads modes, and a Task commit refuses an
  executable it was not granted.

## Success Metrics

1. Success Metric: in a temporary repository fixture with a local fake
   catalog, the target pins the fake `main` commit, restores exactly the
   lock's catalog entries, builds twice, runs the sync, its check, the
   confirmed restore, `make baseline-digests` and two managed refreshes, and
   refuses, before any command, a lock naming an owned skill, a branch name, a
   commit the catalog lacks and a vendored folder without a lock entry; after
   the restore it refuses a folder that differs from the catalog and an entry
   left at another commit.
2. Success Metric: at the audited head the 33 catalog entries of
   `skills-lock.json` name `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, no
   entry names a Roundfix-owned skill, and each upstream snapshot with a
   `github` source and every `github` skill of the five name the same commit.
3. Success Metric: each of the 33 folders equals, byte for byte, the
   catalog's folder at that commit, and `.agents/skills/systematic-debugging`
   holds none of the five removed files; the 12 owned skills keep their bytes.
4. Success Metric: `go test ./internal/baseline ./skills` and the documentation
   contracts pass, `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`
   among them, and in a clone of the audited head `make baseline-digests`
   reports `"changed":false` and the managed refresh reports `File changes: 0`.
5. Success Metric: a second run of the target at the same commit against a
   clone of the audited head changes no file.
6. Success Metric: no test file names `86df8608726979064a8fe46003dfaf403ffe0087`
   or a tree digest the bump writes; a hand edit of a vendored folder with a
   re-hashed lock entry fails the new contract.

## Acceptance evidence

The outside-evidence row rests on text this Spec did not write: the adopted
message from marcioaltoe-skills, reproduced in the Backlog Entry under
`references/`, and the catalog's published commit `c8b2e85`, whose
`skills/06-review-repair/systematic-debugging/` holds six files and none of
the five the message names. The QA checks the message's expected result
against the audited head offline and runs the target against the operator's
clone, read-only, with network protocols disabled.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Prototype in a worktree of `6a0ed904` with the catalog fetched read-only at
  `c8b2e85`: the asset sync rewrote the five snapshots (only
  `systematic-debugging`'s tree moved in `go` and `rust`, plus `ui-ux-pro-max`
  in the two TypeScript snapshots and `developing-with-streamlit` in
  `python`; no skill joined or left one). A restore with the old binary still
  named `86df860`, because it reads the embedded snapshots, so the target
  builds again after the sync. The restore of the 33 then planned 33 lock
  updates and five removals; `make baseline-digests` rewrote seven files under
  `internal/baseline/testdata/`; the managed refresh rewrote one file and a
  second reported `File changes: 0`. The prototyped target produced the same
  20 changed paths and, run again, changed nothing.
- Before the test edits `./internal/baseline` failed in the five tests that
  read the revision constant and in the `ui-ux-pro-max` tree test; after them
  it passed. A hand edit of `ai-memory-retrieval` with a re-hashed lock entry
  passed every existing test (the skill is not required by `go-cli-tui`, so
  ADR-0224's repository test skips it) and failed the new contract.
- Secondbrain, read only: `wiki/index.md` and `qmd query "skills catalog
  marcioaltoe/skills vendored skills lock pinned commit bump" --all --files
  --min-score 0.3` returned the catalog's README and AGENTS mirrors, this
  repository's archived Spec 0208 (the Baseline following the reshaped
  catalog) and three adopters' tooling-authorization notes; none records a
  decision against a pinned bump.
- Exa: published guidance on Git-sourced dependencies resolves a branch with
  `git ls-remote`, writes the full commit and moves a pin only by a deliberate
  update ("Pinning Git Dependencies to Commit SHAs", git-automation.com), and
  a vendoring tool exposes a check that fails on stale vendoring
  (`git-third-party update --check`).

## Glossary

- adds: **Catalog Revision**
- adds: **Catalog Skill Bump**

## Decisions

Autonomous route: no person present; decisions taken from the operator's Cycle AX addendum briefing for Spec 0350 of 2026-11-04 and the characterization above.

The maintainer's answers (AskUserQuestion, 2026-11-04, about 07:15 BRT),
verbatim:

- 0350: "Bump fixado e repetível (Recommended)" — option text: "Um alvo de
  Makefile move as 33 entradas para o commit atual da main do catálogo (ref =
  SHA, não 'main'), reinstala e atualiza o hash; mantém o lock reprodutível e
  os digests da Baseline determinísticos. Rodo o alvo agora (c8b2e85) e a cada
  ciclo."
- Placement: "No ciclo AX, mesma v0.131.0 (Recommended)" — option text:
  "Ficam depois da 0346 na fila (as três mexem em Baseline/skills); uma
  v0.131.0 no fim com as seis; mesmas regras do AX (merge pela fila,
  Makefile/CI liberados). Nota ao catálogo e aos adotantes ao final."
- Rules: "Mesmas do AW (Recommended)" — option text: "Sem supervisão: push e
  PR dos planos, entrega e merge pela fila, consertar e seguir em falha
  técnica; Makefile, .github/workflows, go.mod/go.sum e lint liberados quando
  o desenho pedir; mudança de Baseline que chega aos adotantes liberada."
- Release: "Uma v0.131.0 no fim (Recommended)" — option text: "Uma minor (a
  mudança da Baseline e o runs list fora de repo são feat); notas aos 8
  adotantes, resposta ao Pantheon com o defeito corrigido e o contorno que
  podem largar."

Taken from those sources and the characterization, with no further interview:

- `type: fix`: the vendored skills trail the catalog and the lock does not
  reproduce them; the target is the fix's mechanism.
- The bump also moves the five upstream Setup Snapshots, because ADR-0224
  holds this repository to the `go` snapshot and Doctor would otherwise name
  `systematic-debugging` a Setup Snapshot Difference (ADR-0360).
- Target `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, the catalog's `main` at
  authoring (2026-11-04 04:04 BRT); the next cycle's bump follows a later one.
- The Run gets the commit from the operator's local clone, read-only, as
  Specs 0326 and 0333 did; only the target itself may use the network.
- No local edit of a vendored folder, enforced by the new documentation
  contract and the existing lock-hash tests.

## Open Questions

None.
