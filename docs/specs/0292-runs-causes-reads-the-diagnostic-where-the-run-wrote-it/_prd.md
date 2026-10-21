---
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# Runs causes reads the diagnostic where the Run wrote it

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix runs causes`
classifies each failed Verification attempt partly from the tail of the
diagnostic the Run wrote, but it reads that diagnostic only under the
Artifact Directory the Run recorded, taken literally, and every Implement Run
records none. So the report reads no diagnostic for any Implement Run, and
every signature that reads the diagnostic tail can never match. After this
Spec, the report finds the Run's directory under the Artifact Root the Run
wrote to, through the one resolver that the storage report, GC sanitation and
Run Retention already use, and it keeps its symbolic-link and containment
checks. The TechSpec holds the design.

## Sources

One Backlog Entry, adopted by this Spec's planning change into `references/`
(see `references/_index.md`):

1. [`references/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md`](references/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md),
   type `fix`: `diagnosticTail` returns an empty tail whenever the recorded
   Artifact Directory is empty and otherwise uses the recorded value
   literally; found by code reading while Spec
   `gc-sanitize-reaches-the-default-artifact-roots` was authored, and not yet
   reproduced through the command.

No other open Backlog Entry or unresolved Finding shares this context:
`docs/backlog/` held only this entry on 2026-10-20 and `docs/findings/` held
none. The cycle AM decision gives this Spec the `runs causes` entry alone; the
symbolic-link hardening of `archive` and `cancel` is a separate Spec.

## Prerequisites

None. The resolver this Spec reuses, `store.ResolveRunArtifactRoot`, is on the
default branch since Spec `gc-sanitize-reaches-the-default-artifact-roots`
(v0.100.0, `70c78d15`), the commit this plan was authored against. The other
cycle AM Specs are authored beside this one and deliver before it; none of
them changes `internal/runcause`, `internal/cli/runs_causes.go` or the `runs`
guide, but they may raise the Roundfix skill version first, so task_02 raises
it over whatever value the tree holds when it starts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  field is minted; Run IDs, repository IDs, Artifact Root paths and the
  `roundfix/runs-causes/v1` schema keep their current form. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. The report reads the Run Database, Git
  metadata and files under Roundfix Home, and every test uses temporary homes
  and repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the read path, its checks and their documentation follow recorded decisions. ADR-0301: "GC sanitation classifies the Artifact Root each Run wrote to". ADR-0255: "Run Retention removes terminal Runs whole and compacts incrementally". ADR-0214: "Why a Verification failed is classified by deterministic signatures". ADR-0259: "A test runs in parallel unless it names why it cannot". ADR-0271: "A test bounds by real time only what it asserts". ADR-0285: "Every test file compiles for Windows, and a repository contract proves it". ADR-0189: "An owned skill's version names its content, and its minimum is the bundle". ADR-0256: "A release waits for the skills that describe a changed surface". ADR-0280: "A guide Task records its own surface fingerprint, a PRD records who decided, and the glossary is held to its own rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base". ADR-0265: "A Spec proves what it claims, and every active Spec is held to it". ADR-0183: "An attributed claim carries a receipt the check proves". Set no obligation here, because this Spec changes none of their rules: ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0184 (a TechSpec states a command surface as a transcript); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset); ADR-0296 (the storage report reconciles every page of one Storage Snapshot); ADR-0298 (a timed-out profile proof gets a third attempt, and its setup deadline is configurable); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0306 (archive and cancel refuse a path whose resolved target leaves the repository); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 through AskUserQuestion, authorization "the same
  rules as cycle AL", whose answer "Mesmas regras do AK (Recomendado)" carries
  the cycle AJ grant "Mexer em Governed Paths declarados" ("'Concedo' para os
  caminhos que cada Spec declarar (módulos Baseline, registros, skills,
  testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
  específica)."), and the 2026-09-30 skills grant ("considere autorizado a
  ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- `roundfix runs causes` reads the diagnostic of a failed Verification
  attempt under the Artifact Root the Run wrote to, including Runs that
  recorded no Artifact Directory or a relative one, so diagnostic signatures
  can classify those attempts.
- The report finds that root through the one resolver the storage report, GC
  sanitation and Run Retention use, and carries no resolution of its own.
- Every existing guard on the diagnostic stays: a symbolic link anywhere on
  the path, a path outside the Run's own directory, another Run's directory, a
  directory or a missing file still gives no evidence, and the report still
  writes nothing.

## Core Features

1. The resolved read. The diagnostic tail is read under
   `<resolved root>/runs/<run-id>`, where the resolved root comes from
   `store.ResolveRunArtifactRoot` for the Run's ID, Git root, repository root
   and recorded value, and the user's home that `runs causes` already loads.
2. Kept guards. A Run whose root cannot be resolved, or whose ID is unsafe,
   gives no evidence as today; the containment and component checks run on the
   resolved directory unchanged.
3. Documentation. The `runs` guide, the Roundfix `runs` reference and the
   glossary entry for **Artifact Root** say where the report reads diagnostics
   and that it uses the shared resolver.

## Non-Goals / Out of Scope

- Changing the resolver, the storage report, GC sanitation or Run Retention.
- Recording the resolved Artifact Directory on new Implement Runs.
- Replacing the component walk with `os.Root` or otherwise changing the
  symbolic-link or containment rules; a home reached through a symbolic link
  still yields no diagnostic, as today.
- Changing the signature table, its digest, the flags, help, exit codes or
  JSON fields of `runs causes`.
- `roundfix attach`, which passes a review Run's recorded root to the Review
  Issue selection, and every other reader the TechSpec's sweep lists.
- Re-measuring `docs/references/corrective-causes-measurement.md`.

## Success Metrics

1. Success Metric: on a temporary Roundfix Home whose terminal Implement Run
   recorded no Artifact Directory and wrote its Verification diagnostics under
   the repository's default root, `roundfix runs causes --since 2026-10-01`
   prints the same four items and summary as for a Run that recorded the
   absolute root (`environment database-locked` and `implementation-defect
   go-test-failure` among them), where before this Spec it printed those
   attempts `unclassified`.
2. Success Metric: the same holds for a Run that recorded a relative Artifact
   Directory and wrote under that directory inside its Git root.
3. Success Metric: on the resolved root, a leaf symbolic link, a linked parent
   directory, a path outside the Run's directory, another Run's directory, a
   directory and a missing file each give the attempt no evidence; a Run whose
   root cannot be resolved gives no evidence without failing the report; the
   home's files are byte-identical after the report.
4. Success Metric: the storage report, GC sanitation, Run Retention and
   `runs causes` all call `store.ResolveRunArtifactRoot`, and the production
   code of `internal/runcause` calls neither `config.ResolveArtifactDirectory` nor
   `config.DefaultArtifactDirectoryForPath`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The published measurement `docs/references/corrective-causes-measurement.md`
  (window 2026-09-17 through 2026-10-02, table digest
  `53e5e5a764109b92806c6efe4d88603c13e5de2c6d21816755406deb11b6dd5f`, the
  digest of today's table): "The read-only `roundfix runs causes` report
  contains 177 terminal Runs and 132 failed Verification or corrective items",
  of which `unclassified` is 100 and `environment` is 0, and "There are only
  five observed signatures in this record": `record-or-golden`,
  `review-defect`, `shared-section`, `authorization-or-scope` and
  `help-or-guide`. None of the seven signatures whose only evidence source is
  the diagnostic (`database-locked`, `disk-full`, `network-unreachable`,
  `killed-or-timed-out`, `governed-path`, `lint-or-format`,
  `go-test-failure`) matched any of the 132 items, and its verdict was
  `inconclusive` because unclassified items exceeded one third. That is what
  a report that never reads an Implement Run's diagnostic produces.
- This repository's `runs causes` suites, which no Task of this Spec authors
  and which encode the guards and the transcript that came with the signature
  report of ADR-0214:
  `TestBuildIgnoresALogOutsideTheArtifactDirectory`,
  `TestBuildReportsAttemptsTasksAndCorrectives`,
  `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText`,
  `TestRunsCausesPrintsAClassifiedWindow` and `TestRunsCausesPrintsJSON`.
  They must pass with unchanged bodies at the audited head.
- The Go team's article on traversal-resistant file access
  (https://go.dev/blog/osroot, 2025-03-12), as a challenging source: "If the
  program defends against symlink traversal by first verifying that the
  intended file does not contain any symlinks, it may still be vulnerable to
  time-of-check/time-of-use (TOCTOU) races, where the attacker creates a
  symlink after the program's check". The report's component walk has that
  window today; this Spec keeps the walk as it is and moves only the root it
  walks from (Non-Goals).

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "roundfix
  runs causes diagnostic tail artifact root resolver" --all --files
  --min-score 0.3`. It returned the mirror of Source 1 (score 1.00), the
  mirrors of this repository's `runs` guide and ADR-0301, and older storage
  Archive Records. No wiki page holds a decision on where `runs causes` reads
  diagnostics. No Inbox Entry was written, because the briefing forbids
  writing to the Secondbrain.
- Exa: the Go article quoted above, found with a search for traversal-resistant
  file access in Go. It supports checking every component rather than the
  joined path lexically, and it challenges the check-then-open walk with a
  TOCTOU window, which this Spec records as an accepted limit for a read-only
  report over the user's own Roundfix Home.
- Characterization on 2026-10-20 at `70c78d15`, through the command on a
  temporary Roundfix Home, recorded in the TechSpec.

## Glossary

- changes: **Artifact Root**

## Decisions

- 2026-10-20 (about 18:15–18:35 BRT): the cycle AM scope "R → S → T
  (Recomendado)" (option text: "Três Specs, um minor cada (v0.101–v0.103),
  nesta ordem. R primeiro porque já afeta o oraculum. Mesmas regras do AL.");
  the minor items "Distribuir + U (Recomendado)" (option text: "'runs causes'
  vai em T (mesmo resolvedor e storage). 'Esperar vaga' vai em S (gate e
  preflight da fila). O endurecimento de symlink vira a Spec U, no fim
  (v0.104.0)."), then "Rodar gc; T vira runs causes (Recomendado)"; the
  authorization, the same rules as cycle AL.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AM briefing for Spec 0292 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The report calls `store.ResolveRunArtifactRoot`; no copy of the resolution
  is added.
- Both resolver errors that mean "no root" give no evidence, as a blank root
  and an unsafe Run ID do today; any other resolver error fails the report,
  as an unreadable diagnostic already does.
- No ADR: ADR-0301 already decided that one resolver turns a recorded
  Artifact Directory into the root a Run wrote to; this Spec adds a read-only
  caller and changes no proof. The glossary entry for **Artifact Root** names
  the new caller.

## Open Questions

None.
