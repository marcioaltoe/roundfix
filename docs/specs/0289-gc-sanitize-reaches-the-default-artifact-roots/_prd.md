---
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# GC sanitize reaches the default Artifact Roots

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix gc sanitize`
groups Runs by the Artifact Directory they recorded and refuses any root that
is not a clean absolute path, and every Implement Run records an empty one, so
sanitation keeps every default root under `~/.roundfix/artifacts/` and reports
nothing to reclaim there, about 605 MiB on the maintainer's machine. After this
Spec, sanitation resolves each recorded root the way the Run wrote to it,
through the one resolver the storage report and Run Retention also use, and
classifies the resolved root with the proofs it already had. The TechSpec holds
the design.

## Sources

One Backlog Entry, adopted by this Spec's planning change into `references/`
(see `references/_index.md`):

1. [`references/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md`](references/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md),
   type `fix`: the maintainer's `roundfix gc sanitize` dry run of 2026-10-20
   listed Artifact Root `""` as `unsafe` and reclaimed nothing under the
   default roots.

No other open Backlog Entry or unresolved Finding shares this context:
`docs/backlog/` held only this entry on 2026-10-20 and `docs/findings/` held
none. The storage report half of the same defect belongs to Spec
`reports-that-count-what-they-say`, which this Spec builds on (Prerequisites).
This planning change also writes one open Backlog Entry for later,
`docs/backlog/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md`:
`roundfix runs causes` reads the same empty recorded root, a read-only report
outside the three consumers the maintainer's decision names. `roundfix spec
judge` suggested grouping that entry with this Spec (P 0.88); it stays apart
because the cycle AL decision scopes this Spec to the storage report,
sanitation and Run Retention, and the run cause report reads diagnostics
through its own containment checks, a surface this Spec does not touch.

## Prerequisites

This plan assumes the code of Spec `reports-that-count-what-they-say` on the
default branch: `roundfix storage report` resolving an empty or relative
recorded Artifact Root in `internal/store/journal.go` (its Task 02). That Spec
was still in the Delivery Queue when this plan was authored against `7a678ddb`
(v0.97.0) and the Spec's final Run branch. The plan merges only after it, so
the Task Graph declares no `requires`. Before this plan enters a Run, the
operator confirms that `internal/store/journal.go` on the default branch holds
`func storageArtifactRoot`, then rebases this plan and runs `roundfix spec
check 0289-gc-sanitize-reaches-the-default-artifact-roots --strict
--run-verification` again. The other cycle AL Spec, on archive citations
rewritten from any worktree, is authored beside this one and delivers first;
it changes neither sanitation, the storage report nor Run Retention, but it
may raise the Roundfix skill version first.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  field is minted; Run IDs, repository IDs and Artifact Root paths keep their
  current derivation. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Sanitation reads the Run Database, Git
  metadata and directories under Roundfix Home, and every test uses temporary
  homes and repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the resolver, the sanitation proofs, the recheck and their documentation follow recorded decisions. ADR-0301: "GC sanitation classifies the Artifact Root each Run wrote to". ADR-0255: "Run Retention removes terminal Runs whole and compacts incrementally". ADR-0296: "The storage report reconciles every page of one Storage Snapshot". ADR-0033: "The Run Event Journal is pruned by a retention window". ADR-0259: "A test runs in parallel unless it names why it cannot". ADR-0271: "A test bounds by real time only what it asserts". ADR-0285: "Every test file compiles for Windows, and a repository contract proves it". ADR-0189: "An owned skill's version names its content, and its minimum is the bundle". ADR-0256: "A release waits for the skills that describe a changed surface". ADR-0280: "A guide Task records its own surface fingerprint, a PRD records who decided, and the glossary is held to its own rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base". ADR-0265: "A Spec proves what it claims, and every active Spec is held to it". ADR-0183: "An attributed claim carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0171 (a retention prune reports only what it reclaimed); ADR-0172 (Doctor reports reclaimable Run storage from cheap reads); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0184 (a TechSpec states a command surface as a transcript); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0266 (a signal goes only to a proven process, and worktrees share one build cache); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset); ADR-0298 (a timed-out profile proof gets a third attempt, and its setup deadline is configurable); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AL approval of 2026-10-20 through AskUserQuestion, option "Mesmas regras do
  AK (Recomendado)", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."), and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0289-gc-sanitize-reaches-the-default-artifact-roots/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/storage.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- `roundfix gc sanitize` reaches the Artifact Root every Run wrote to,
  including Runs that recorded no root or a relative one, and its dry run
  names each resolved root and what it would reclaim there.
- Deletion widens only to what the existing proofs establish on the resolved
  root: live Runs, recent Runs, directories recorded under another root,
  symbolic links, roots outside Roundfix Home, roots that are not every
  resolving Run's default, and everything outside `runs/` stay.
- One resolver serves the storage report, GC sanitation and Run Retention.

## Core Features

1. The shared resolver. The function the storage report already uses becomes
   the one exported resolver; Run Retention calls it instead of its own copy,
   and the storage report's output stays the same.
2. Sanitation over resolved roots. Runs are grouped by their resolved root and
   each group is classified with the existing proofs; a Run whose root cannot
   be resolved keeps its group `unsafe`; a directory proven only by its Run's
   absence is checked again just before removal.
3. A dry run that names what it resolved. Each root lists the recorded values
   that resolved to it and the directories and bytes it would reclaim.
4. Documentation. The `gc` guide, the Roundfix `storage` reference and the
   glossary describe the resolved root and the sanitation proofs.

## Non-Goals / Out of Scope

- Running `roundfix gc sanitize --apply` on any real Roundfix Home: the
  operator does that after the release (Unreachable Acceptance).
- Recording the resolved Artifact Directory on new Implement Runs.
- Changing what the per-repository `roundfix gc` Journal Retention sweep, the
  Worktree Build Cache sweep or the Item Binary sweep removes, or the storage
  report's sections, columns or reconciliation.
- `roundfix runs causes` and every other reader of a recorded root; the open
  Backlog Entry above holds the one found.
- Adding a JSON output or a flag to `gc sanitize`.

## Success Metrics

1. Success Metric: on a temporary Roundfix Home whose terminal Implement Run
   recorded no Artifact Root, `roundfix gc sanitize` prints the repository's
   default root with `Recorded as: "" by 1 Run(s)`, classifies it `orphaned`,
   lists the Run's directory as `would remove` and prints the root's reclaim
   lines exactly as Surface Transcript 1; `--apply` removes that directory and
   an orphaned one, and a second `--apply` removes nothing.
2. Success Metric: on the same kind of home, sanitation keeps the directory of
   an Active Run, of a terminal Run inside Journal Retention, of a Run recorded
   under another root, a relative root resolved outside Roundfix Home, a
   symbolically linked default root, a default root that another repository
   recorded as its own absolute root, a Run whose root cannot be resolved, and
   every file outside `runs/`; a Run whose row appears after discovery and
   before removal keeps its directory.
3. Success Metric: `roundfix storage report` prints the same Artifact Roots
   lines as before this Spec, Run Retention keeps and removes the same Runs
   for absolute and empty recorded roots, and neither carries a resolution of
   its own.
4. Success Metric: on the maintainer's machine after the release, the dry run
   names each default root it resolved with what it would reclaim, and the
   maintainer applies it only if it comes out clean.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The maintainer's `roundfix gc sanitize` dry run of 2026-10-20 with
  `roundfix 0.97.0`, read-only, recorded in Source 1 and quoted in the
  TechSpec: one Artifact Root `""`, `unsafe`, 0 directories and 0 bytes
  reclaimable.
- This repository's sanitation and Run Retention suites, which no Task of
  this Spec authors and which encode the proofs the maintainer approved with
  Spec `run-storage-compaction-and-global-sanitation` and ADR-0255:
  `TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories`,
  `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval`, the five tests in
  `internal/cli/gc_sanitize_pre_key_root_test.go`,
  `TestGCKeepsARunUnderAnUnprovenArtifactRoot` and
  `TestGCRemovesRunsPastRunRetention`. They must pass unchanged at the audited
  head.
- The Rust compiler's incremental-compilation file system notes
  (https://doc.rust-lang.org/nightly/nightly-rustc/rustc_incremental/persist/fs/index.html),
  on collecting session directories: "There is still a small time window
  between the original process creating the lock file and actually locking
  it." A collector that proves a directory unused has to close the window in
  which its owner exists but has not yet claimed it; here the Run row precedes
  its directory, so a second read of the Run Database closes it.
- The Nix store collector (https://github.com/NixOS/nix/blob/ab5a9cf2db31d3840a801385349b9d23deb29ecc/src/libstore/gc.cc):
  "There may be temp directories in the store that are still in use by
  another process. We need to be sure that we can acquire an exclusive lock
  before deleting them." Deletion re-proves ownership at the moment of
  removal rather than trusting the earlier scan.

## Unreachable Acceptance

- criterion: Success Metric 4, on the maintainer's real Roundfix Home after
  the release, `roundfix gc sanitize` names each default Artifact Root it
  resolved with what it would reclaim, and `roundfix gc sanitize --apply` runs
  only if that dry run is clean.
  reason: no Verification or QA row may read or write the real `~/.roundfix`,
  and applying sanitation there deletes the maintainer's data, which the
  maintainer decided to do himself.
  satisfied-by: the maintainer runs `roundfix gc sanitize`, and then
  `roundfix gc sanitize --apply` only if the dry run is clean, after the
  release that carries this Spec.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "roundfix
  gc sanitize artifact root storage retention reclaim" --all --files
  --min-score 0.3`. It returned the mirror of Source 1, the mirrors of this
  repository's `storage` reference, `storage report` and `gc` guides, the
  active PRD of Spec `reports-that-count-what-they-say`, and the records of
  Specs `run-store-retention`, `run-storage-that-says-what-it-holds`,
  `a-run-database-that-keeps-only-recent-runs` and
  `run-storage-compaction-and-global-sanitation`. The last one is where
  sanitation and its six classifications were decided; its TechSpec, read
  at its `source_revision`, states "Ambiguity preserves" and "Sanitation twice
  reclaims something, then zero", which this Spec keeps. No wiki page holds a
  decision on empty recorded roots. No Inbox Entry was written, because the
  briefing forbids writing to the Secondbrain.
- Exa: the two collectors quoted above support re-proving a directory unused
  at removal. Challenging: both use file locks the owner holds, which
  Roundfix Runs do not take on `runs/<run-id>`; the ordering of the Run row
  before the directory is what makes a database read a sufficient proof here,
  and the TechSpec records where that ordering comes from.
- Characterization on 2026-10-20 at `7a678ddb` and at the final Run branch
  of Spec `reports-that-count-what-they-say`, on a temporary Roundfix Home,
  recorded in the TechSpec.

## Glossary

- adds: **Artifact Root**

## Decisions

- 2026-10-20 (about 13:10 BRT): the cycle AL sequence "Aprovar como está
  (Recomendado)" (option text: "Duas Specs separadas, um minor cada, nesta
  ordem; P primeiro porque trava o caminho de recuperação documentado."); the
  scope of this Spec "Resolver + eu aplico (Recomendado)" (option text: "O
  sanitize resolve a raiz vazia como o Run Retention e usa as provas
  existentes; depois da release eu rodo o dry-run no seu ~/.roundfix e aplico
  só se ele sair limpo."); the authorization "Mesmas regras do AK
  (Recomendado)".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AL briefing for Spec 0289 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The storage report's resolver becomes the one resolver; Run Retention and
  sanitation call it, and no copy stays.
- Sanitation groups Runs by their resolved root and keeps every existing proof
  in its existing order; an unresolvable Run keeps its whole group.
- A directory proven only by absence is checked against the Run Database again
  just before removal, because the default roots are where every new
  Implement Run writes.
- ADR-0301 records the decision, because it widens what a deleting command may
  remove.

## Open Questions

None.
