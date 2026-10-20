---
spec: 0287-reports-that-count-what-they-say
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# Reports that count what they say

Minimal PRD for a bug fix, minted by `write-techspec`. Three reports print a
number or a summary that does not describe what happened: the Linux Lane's
pass line counts the package patterns it handed to the container instead of
the packages that ran, `roundfix storage report` measures the Artifact Root a
Run recorded literally and so finds nothing where every Implement Run wrote
its artifacts, and an Archive Record's summary stops at the colon that
introduces the list it was meant to carry. After this Spec each report states
what it measured, and the TechSpec holds the design.

## Sources

Three Backlog Entries, written and adopted by this Spec's planning change into
`references/` (see `references/_index.md`), one per defect:

1. [`references/2026-10-20-the-linux-lane-counts-package-patterns-as-packages.md`](references/2026-10-20-the-linux-lane-counts-package-patterns-as-packages.md),
   type `fix`: `roundfix deliver lane --all` printed `Linux Lane passed 1
   package(s)` while 43 packages ran.
2. [`references/2026-10-20-the-storage-report-misses-the-default-artifact-roots.md`](references/2026-10-20-the-storage-report-misses-the-default-artifact-roots.md),
   type `fix`: the report listed Artifact Root `.` as `missing` with 508 Run
   rows and 0 bytes while about 605 MiB sat under the default roots.
3. [`references/2026-10-20-an-archive-record-summary-ends-at-the-colon-that-introduces-its-list.md`](references/2026-10-20-an-archive-record-summary-ends-at-the-colon-that-introduces-its-list.md),
   type `fix`: the Archive Record of Spec 0282 ends with "Two environment gaps
   failed it:".

No open Backlog Entry or unresolved Finding shares this context: `docs/backlog/`
and `docs/findings/` held no entry on 2026-10-20, and the two pending
Secondbrain Inbox Entries for Roundfix are the sources of Specs 0285 and 0286.
This planning change also writes one open Backlog Entry for later,
`docs/backlog/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md`:
`roundfix gc sanitize` reads the same empty recorded root, and changing what a
deleting command may remove belongs to a Spec of its own.

## Prerequisites

None. This Spec is authored against `9448aa64` (v0.95.0). Specs 0285 and 0286
are authored beside it and deliver before it; neither changes the Linux Lane,
the storage report or the summary reader. Either may raise the Roundfix skill
version and re-record the Behavior Surface Record first, so this plan is
rebased and checked again before its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the lane's count line and the no-count wording are fixed
  strings, and the Archive Record schema stays `roundfix/archive-record/v1`.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. The lane keeps `--network none` and
  `--pull never`, the storage report reads only Roundfix Home and the
  recorded roots, and every test uses temporary repositories, homes and
  compiled stand-ins. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the three fixes, their tests and
  their documentation follow recorded decisions. ADR-0294: "The Linux Lane
  runs like CI". ADR-0297: "The Linux Lane reads the repository history
  read-only". ADR-0288: "The Delivery Queue gate runs the changed Go packages
  in a Linux Lane". ADR-0296: "The storage report reconciles every page of
  one Storage Snapshot". ADR-0255: "Run Retention removes terminal Runs whole
  and compacts incrementally". ADR-0247: "An archive leaves an Archive
  Record, and the Spec folder stays in Git". ADR-0277: "A Legacy Archive
  Folder conversion keeps the summary, links and dates its folder held".
  ADR-0262: "A cancelled Spec leaves an Archive Record and returns its
  sources". ADR-0125: "A spawned fixture is a compiled binary, never a
  written script". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0271: "A test bounds by real time only what it asserts". ADR-0298: "A timed-out profile proof gets a third attempt, and its setup deadline is configurable".
  ADR-0285: "Every test file compiles for Windows, and a repository contract
  proves it". ADR-0189: "An owned skill's version names its content, and its
  minimum is the bundle". ADR-0256: "A release waits for the skills that
  describe a changed surface". ADR-0280: "A guide Task records its own
  surface fingerprint, a PRD records who decided, and the glossary is held to
  its own rule". ADR-0104: "A Spec accepts on evidence it did not author".
  ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach".
  ADR-0276: "A row the authorization keeps offline qualifies, and a byte
  invariant compares with the Delivery Base". ADR-0265: "A Spec proves what
  it claims, and every active Spec is held to it". ADR-0183: "An attributed
  claim carries a receipt the check proves". Set no obligation here, because
  this Spec changes none of their rules: ADR-0092 (triage intent lives in a
  typed backlog); ADR-0120 (retired documentation lives under one history
  root); ADR-0165 (a blocking review after archive parks publication for a
  corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a
  qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from
  its merge base); ADR-0184 (a TechSpec states a command surface as a
  transcript); ADR-0196 (a pre-PR review finding parks only after
  validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds);
  ADR-0229 (an operator archive resumes any park from the Run start);
  ADR-0233 (a skill version raise is regenerated at merge, and a review-only
  correction returns to review); ADR-0237 (a Delivery Retry records a merge
  made outside the queue); ADR-0248 (existing history is sanitized in batches
  after a history-full tag); ADR-0250 (a module version is chosen when it is
  recorded, and the coverage record lists every platform); ADR-0251 (a Legacy
  Archive Folder is read leniently, and a failed QA keeps its verdict);
  ADR-0254 (baseline update sanitizes pending history in the change it
  plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a
  retirement writes reduced history); ADR-0260 (an archive relinks live
  citations, and a Baseline update lists the citations it breaks); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0261
  (an Archive Record names its pull request, and readers fetch the pull
  request head); ADR-0269 (archive and cancel refuse ignored Spec files and
  finish a failed removal forward); ADR-0272 (a review record survives a fix,
  and a correction inside the archived Spec's declared paths returns to
  review); ADR-0278 (a retirement deletes, and the History Root holds only
  Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a
  dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces
  an executable it rebuilds instead of writing into it); ADR-0282 (the
  history cleanup names a superseded ADR in its successor, and a large
  candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets
  one Correction Review of its fix); ADR-0290 (merged items give back their
  Item Binaries, and Run Worktrees leave only by proof); ADR-0295
  (reachability reads the fetched default branch, and origin names it when
  origin/HEAD is unset). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 through AskUserQuestion, option "Sim, mesmas
  regras do AJ (Recomendado)", which carries the cycle AJ grant "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0287-reports-that-count-what-they-say/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config file change.

## Goals

- The Linux Lane's pass line, in the Delivery Queue gate and in `roundfix
  deliver lane`, states the number of packages that ran, and never a number
  it did not count.
- `roundfix storage report` measures every Run's artifacts where the Run
  wrote them, including Runs that recorded no Artifact Root or a relative
  one, and stays read-only.
- An Archive Record's summary carries the list or paragraph its colon
  introduces, whoever writes the record: the Archive Command, the Cancel
  Command or the conversion of a Legacy Archive Folder.
- Records already written stay byte-identical.

## Core Features

1. The lane's package count. The lane script reports how many packages `go
   list` selected before it runs `go test`, and the host prints that number;
   when the script reported none, the pass line says so without a number.
2. The storage report's resolved roots. An empty recorded root resolves to
   the Run's repository default root, a relative one to the Run's Git root,
   and an absolute one stays as recorded; the repository, state and Artifact
   Root groups all measure the resolved root.
3. The summary that carries its list. A summary paragraph that ends with a
   colon continues with the list or paragraph that follows it, flattened
   into the same line.
4. Documentation. The `deliver`, `storage report`, `archive` and `history`
   guides, the Roundfix `deliver`, `storage` and `archive` references and
   the glossary's Archive Record entry describe each change.

## Non-Goals / Out of Scope

- Changing what `roundfix gc sanitize` may remove: it still preserves an
  empty recorded root; the open Backlog Entry above holds that work.
- Recording the resolved Artifact Directory on new Implement Runs: the
  report resolves the empty value it finds, which also covers every Run
  already recorded.
- Rewriting any Archive Record already written, including the 0282 record.
- Counting packages on the Linux Lane's failure lines, or changing the
  lane's outcome classification, image, mounts or package selection.
- Changing the storage report's sections, columns or reconciliation.

## Success Metrics

1. Success Metric: with compiled Docker and Go stand-ins whose Docker
   reports `roundfix-linux-lane: packages: 43`, `roundfix deliver lane
   --all` and the Delivery Queue gate print `Linux Lane passed 43
   package(s)`; a stand-in that reports no count gives the no-count line;
   the lane script prints the count line before `go test`.
2. Success Metric: on a temporary Roundfix Home whose Implement Run recorded
   no Artifact Root, `roundfix storage report` lists the repository's
   default root as `present` with that Run's bytes, the Run's repository and
   state groups carry the same bytes, a relative recorded root is measured
   under the Run's Git root, and no directory is created.
3. Success Metric: a PRD whose title paragraph ends with a colon and is
   followed by a list gives an Archive Record whose summary holds every item;
   the same holds for a QA Report Outcome, a cancelled Spec and a Legacy
   Archive Folder; a record that already ends with a colon reads back and
   renders byte-identically.
4. Success Metric: on the maintainer's machine after merge, the storage
   report lists the default roots as `present` with their bytes, and a real
   `roundfix deliver lane --all` on `main` prints the number of packages
   that ran.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- This repository's Archive Records, which no Task of this Spec writes: on
  2026-10-20, 38 of the 256 records under `docs/history/specs/` ended with a
  colon, among them
  `docs/history/specs/0282-the-linux-lane-runs-like-ci.md`. The QA gate
  counts them at the Delivery Base and proves none changed.
- The operator's on-demand lane run on `main`, reported in the cycle AK
  briefing of 2026-10-20 and recorded in Source 1: 43 packages ran and the
  pass line said 1.
- The maintainer's `roundfix storage report` and `roundfix gc sanitize` dry
  run of 2026-10-20, read-only, recorded in Source 2: Artifact Root `.`
  `missing` with 508 Run rows and 0 bytes, and a recorded root `""`.
- The Go command documentation, "Package lists and patterns"
  (https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns): "Usually,
  [packages] is a list of package patterns, which can take several forms:"
  and, of the `...` wildcard, "elements in filesystem or import paths expand
  to match 0 or more path elements."; and `go list` (https://pkg.go.dev/cmd/go#hdr-List_packages_or_modules):
  "List lists the named packages, one per line."

## Unreachable Acceptance

- criterion: Success Metric 4, on the maintainer's machine, `roundfix
  storage report` lists the default Artifact Roots as `present` with their
  bytes, and a real `roundfix deliver lane --all` on `main` prints the number
  of packages that ran.
  reason: no Verification or QA row may read the real `~/.roundfix`, and a
  Run sandbox may not reach the Docker daemon; the QA gate runs the real lane
  only when it can.
  satisfied-by: the operator runs `roundfix storage report` and `roundfix
  deliver lane --all` on `main` after this Spec merges.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query "storage
  report artifact root empty artifact_dir missing" --all --files --min-score
  0.3`, `qmd query "linux lane package count go list patterns" --all --files
  --min-score 0.3` and `qmd query "archive record summary first paragraph
  colon list" --all --files --min-score 0.3`. They returned mirrors of this
  repository's `storage` and `deliver` guides and references, ADR-0294,
  ADR-0247, Spec 0271's record and an unrelated writing guide; none holds a
  decision on these defects. No Inbox Entry was written, because the briefing
  forbids writing there.
- Exa: the Go command documentation quoted above. Challenging: the pattern
  `./...` excludes `testdata` and nested modules, so even the number of
  directories under the module would not be the number of packages that ran;
  only the container's own `go list` knows it.
- Probes on 2026-10-20 at `9448aa64`, each a `go test -overlay` file that
  wrote nothing to the repository: a store whose Implement Run recorded no
  root reported `root="." missing=true rows=1 bytes=0` although the Run's
  default root held 10 bytes, and a relative root reported `missing` too; a
  PRD whose title paragraph ended "Two gaps failed it:" before a list gave
  the summary "Two gaps failed it:", and a QA Outcome ending with a colon
  gave "The gate held three facts:". The TechSpec records each.

## Glossary

- changes: **Archive Record**

## Decisions

- 2026-10-20 (about 06:30 BRT): the cycle AK decomposition "Duas Specs + um
  lote (Recomendado)", M (0285) → N (0286) → O (0287), one minor each; the
  authorization "Sim, mesmas regras do AJ (Recomendado)".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AK briefing for Spec 0287 of 2026-10-20 and the probes characterized in the TechSpec.

Taken from those sources, with no further interview:

- The lane counts the packages `go list` selected in the container, reported
  by the lane script, rather than renaming the number "patterns".
- The storage report resolves an empty or relative recorded root the way the
  Run resolved it when it wrote its artifacts.
- A summary paragraph that ends with a colon carries the list or paragraph
  that follows it.
- No new ADR: each fix restores what its report already promised.

## Open Questions

None.
