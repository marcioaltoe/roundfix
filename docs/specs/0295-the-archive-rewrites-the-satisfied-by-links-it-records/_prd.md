---
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
status: active
created: 2026-10-21
surfaces: [backend, cli, docs]
---

# The archive rewrites the satisfied-by links it records as unproven

Minimal PRD for a bug fix, minted by `write-techspec`. When a Spec archives
with a qualifying `partial` QA verdict, `roundfix archive` copies the
`satisfied-by` text of each `## Unreachable Acceptance` declaration into the
Archive Record's `unproven` list without change. A relative link in that text
was written against the Spec folder, so in `docs/history/specs/<slug>.md` it
names a file under `docs/history/specs/references/` that never exists, and
`roundfix spec check --strict` fails with `SC-HISTORY-RETAINED` right after
the archive. `roundfix history sanitize` cannot repair the record either: it
points a file the archive promoted to `docs/references/` at the same missing
path, and applying needs a `history-full` tag. After this Spec the archive
rewrites each `unproven` entry exactly as it rewrites the record's outcome,
the sanitize repair points a promoted file at its `docs/references/` copy,
and the strict check suggests that same copy. The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md`](references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md),
   origin pantheon, type hint `fix`. On 2026-10-20, at about 16:40 BRT,
   roundfix 0.97.0 archived pantheon's Spec 0049 with a qualifying `partial`
   and 12 Unreachable Acceptance declarations whose `satisfied-by` linked
   `references/<file>.md`. `roundfix spec check --strict` then reported 15
   `SC-HISTORY-RETAINED` findings, each a link to
   `docs/history/specs/references/<file>.md`. `roundfix history sanitize`
   planned Git Locators for the removed files and
   `docs/history/specs/references/medicao-de-sete-dias.md` "(pull request
   #91)" for the file the archive had promoted to
   `docs/references/medicao-de-sete-dias.md`, and `--apply` needs a
   `history-full` tag the repository does not keep. The operator fixed the
   record by hand.

No open Backlog Entry or unresolved Finding shares this context. The only open
Backlog Entry at the authoring base,
`docs/backlog/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md`,
is owned by Spec 0292 and is about `roundfix runs causes`. `docs/findings/`
holds no entry.

## Prerequisites

None. This Spec was characterized at `956066a7` (v0.101.0) and is authored
against `8fb2da1f`, which adds Spec 0291's delivery and a CI timeout change
that touch none of its paths. The cycle AM Specs 0292 and 0293 are active on
that commit and deliver first. Spec 0293 changes
`internal/spec/repository_path.go`, `internal/spec/archive.go`,
`internal/spec/cancel.go`, `internal/cli/archive.go` and
`internal/cli/cancel.go`; this Spec changes none of them and keeps away from
their containment logic. Spec 0293 also edits the `archive` user guide, the
Roundfix skill's `archive` reference, `CONTEXT.md` and the Behavior Surface
Record, which this Spec edits in other paragraphs after it. The plan is
rebased and checked again before its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Archive Record keeps `roundfix/archive-record/v1` and its
  fields, and only the text of `unproven` entries changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the commands read and write only the
  checkout they run in and its local Git history, and every test uses
  temporary repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the fix records its decision and
  keeps every other archive and history decision. ADR-0308: "An archive
  rebases the unproven links it records, and the Record Link Repair points a
  promoted target at its copy". ADR-0247: "An archive leaves an Archive
  Record, and the Spec folder stays in Git". ADR-0277: "A Legacy Archive
  Folder conversion keeps the summary, links and dates its folder held".
  ADR-0278: "A retirement deletes, and the History Root holds only Archive
  Records". ADR-0230: "An archived Spec keeps its relative links". ADR-0260:
  "An archive relinks live citations, and a Baseline update lists the
  citations it breaks". ADR-0261: "An Archive Record names its pull request,
  and readers fetch the pull request head". ADR-0248: "Existing history is
  sanitized in batches after a history-full tag". ADR-0254: "Baseline update
  sanitizes pending history in the change it plans". ADR-0299: "The history
  conversion rewrites every citation it lists, and the check reports a
  removed history file cited without a Git Locator". ADR-0251: "A Legacy
  Archive Folder is read leniently, and a failed QA keeps its verdict".
  ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach".
  ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts". ADR-0285:
  "Every test file compiles for Windows, and a repository contract proves
  it". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "A Spec proves what it claims,
  and every active Spec is held to it". ADR-0183: "An attributed claim
  carries a receipt the check proves". ADR-0184: "A TechSpec states a command
  surface as a transcript". ADR-0280: "A guide Task records its own surface
  fingerprint, a PRD records who decided, and the glossary is held to its own
  rule". ADR-0256: "A release waits for the skills that describe a changed
  surface". Set no obligation here, because this Spec changes none of their
  rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120
  (retired documentation lives under one history root); ADR-0167 (the pre-PR
  Pull Request row never decides a qualifying partial); ADR-0229 (an operator
  archive resumes any park from the Run start); ADR-0250 (a module version is
  chosen when it is recorded, and the coverage record lists every platform);
  ADR-0258 (a Spec is authored against the QA rerun classes); ADR-0262 (a
  cancelled Spec leaves an Archive Record and returns its sources); ADR-0269
  (archive and cancel refuse ignored Spec files and finish a failed removal
  forward); ADR-0272 (a review record survives a fix); ADR-0281 (a test
  replaces an executable it rebuilds); ADR-0282 (the history cleanup names a
  superseded ADR in its successor); ADR-0286 (a fixed round-2 finding gets
  one Correction Review); ADR-0295 (reachability reads the fetched default
  branch); ADR-0298 (a timed-out profile proof gets a third attempt);
  ADR-0165 (a blocking review after archive parks publication for a
  corrective Spec); ADR-0197 (a pre-PR reviewer lineage spans at most two
  rounds); ADR-0237 (a Delivery Retry records a merge made outside the
  queue); ADR-0279 (status reports a Queue Owner proven gone); ADR-0290
  (merged items give back their Item Binaries); ADR-0306 (archive and cancel
  refuse a path whose resolved target leaves the repository, whose
  containment logic this Spec does not touch); ADR-0169 (the pre-PR review
  diffs the candidate from its merge base); ADR-0233 (a skill version raise
  is regenerated at merge); ADR-0196 (a pre-PR review finding parks only
  after validation). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21 through AskUserQuestion, authorization "the same
  rules as cycle AM", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).") and the 2026-09-30
  skills grant "considere autorizado a ajustar todas as skills se
  necessário". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0295-the-archive-rewrites-the-satisfied-by-links-it-records/_authorization.md`;
  bounded files: `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`. No
  Makefile, workflow, Go module, lint, formatter, Baseline module or Project
  Config file change.

## Goals

- An Archive Record written for a qualifying `partial` holds no relative link
  that leaves the record's directory without resolving, so `roundfix spec
  check --strict` passes right after the archive.
- Each `unproven` entry names a removed folder file by its Git Locator and a
  promoted file by a relative link to its `docs/references/` copy, exactly as
  the record's outcome does.
- `roundfix history sanitize` and the history section of `roundfix baseline
  update` repair a record already written that way: a promoted target points
  at its `docs/references/` copy; every other target keeps today's repair.
- The strict check's fix line for such a link names the same promoted copy.

## Core Features

1. The Archive Command, and every archive inside a Run or a Delivery Queue
   item, rewrite each `unproven` entry with the Legacy Link Rebase that
   already rewrites the outcome, as text written in the Spec folder's
   `_prd.md`. The conversion of a Legacy Archive Folder does the same.
2. The Record Link Repair resolves a broken record link against the record's
   `source` folder first for promotion: when the target's file name names a
   `docs/references/` copy the record lists under `promoted` and the tree
   holds, it rewrites the link to that copy and names the copy in the plan.
3. `SC-HISTORY-RETAINED` prints that copy in the fix line of such a link.
4. Documentation. The `archive`, `history`, `baseline` and `spec` user guides,
   the Roundfix skill's `archive`, `baseline` and `spec` references and the
   glossary state the new rewrites.

## Behavior changes

An `unproven` entry with a relative link to a file inside the Spec folder is
now written rewritten, where it used to be copied verbatim: a promoted target
becomes `../../references/<file>` (for the built-in Spec Root) and any other
folder file becomes its label followed by its Git Locator. The plan line of
`roundfix history sanitize` and of `roundfix baseline update` for a promoted
target now reads `rewrites to docs/references/<file>`, and the Plan Digest of
the update binds that text. An entry without such a link, and every record
already written, keep their bytes until an operator applies the repair. No
existing test, golden or fixture changes.

## Non-Goals / Out of Scope

- The `history-full` tag requirement of `--apply`: unchanged, because it
  anchors every Git Locator.
- The repair of a record link whose target is neither promoted nor held by a
  reachable commit: it keeps today's record-relative path with the pull
  request.
- Any change to the containment checks Spec 0293 adds, to the Archive
  Relink, to the exact-archive or Pull Request Record Commit proofs, or to
  the record schema.
- Rewriting records already written by any command other than an applied
  Record Link Repair.

## Success Metrics

1. Success Metric: in a temporary repository, a Spec with a qualifying
   `partial` and two Unreachable Acceptance declarations whose `satisfied-by`
   link `references/removed.md` and `references/promoted.md`, archived with
   `--promote references/promoted.md`, writes a record whose entries read
   `` `references/removed.md` (`docs/specs/<slug>/references/removed.md@<commit>`) ``
   (or with `(pull request #<n>)` when no reachable commit holds it) and
   `` [`references/promoted.md`](../../references/promoted.md) ``, and
   `roundfix spec check --strict` reports no `SC-HISTORY-RETAINED` finding.
2. Success Metric: on a record written by the Delivery Base binary for the
   same Spec, `roundfix history sanitize` names `docs/references/promoted.md`
   as the rewrite of the promoted link, where the Delivery Base names
   `docs/history/specs/references/promoted.md`, and keeps the removed link's
   Git Locator; the strict check's fix line names the same copy.
3. Success Metric: an `unproven` entry without a relative link, and every
   record already written, keep their bytes through the archive and through
   a plan of the repair.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-21.

- The Go standard library documentation of `net/url.URL.ResolveReference`,
  readable offline with `go doc net/url.URL.ResolveReference`:
  "ResolveReference resolves a URI reference to an absolute URI from an
  absolute base URI u, per RFC 3986 Section 5.2." A relative link's target
  depends on the document that holds it, which is why text moved to another
  directory must be rewritten.
- GitHub Docs, "Basic writing and formatting syntax", section Relative links
  (https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax):
  "A relative link is a link that is relative to the current file."
- The symptom the adopted Inbox Entry reports, reproduced by a binary built
  from the Delivery Base, which this Spec did not build.

## Unreachable Acceptance

None. Every criterion runs in temporary repositories; the GitHub Docs page is
recorded under its own network-prohibited row and never decides a qualifying
`partial`.

## Validation

The decisions were checked on 2026-10-21 against these sources.

- Characterization at `956066a7` with `bin/roundfix` built by `make build`
  from that commit, in temporary repositories under a temporary `HOME` and
  `GIT_CONFIG_GLOBAL`. The archive copied both `satisfied-by` links
  verbatim, `spec check --strict` reported two `SC-HISTORY-RETAINED`
  findings, and `history sanitize` planned the promoted file's
  `docs/history/specs/references/promoted.md` "(pull request #91)" when the
  file exists only on the delivery branch, as in the Source. The TechSpec
  records each line.
- A prototype of Core Features 1-3 on a scratch copy of `956066a7` wrote the
  entries Success Metric 1 names, made the strict check pass, made the
  sanitize plan name `docs/references/promoted.md`, applied it to the same
  bytes the fixed archive writes, and left every test in `internal/spec` and
  `internal/speccheck` and every archive, history, cancel and record test in
  `internal/cli` passing. The prototype was discarded.
- Secondbrain, read only. Read `wiki/index.md` and
  `wiki/concepts/ciclo-de-vida-de-decisoes-e-raiz-docs-history.md`, which
  records the retirement direction: retired material leaves the tree and Git
  keeps its text, so a record must name removed files by locator. Ran `qmd
  query "archive record relative links unproven satisfied-by git locator
  promoted references" --all --files --min-score 0.3`, which returned the
  mirrors of ADR-0247, ADR-0230, ADR-0299, the `history`, `spec` and
  `deliver` guides and the records of Specs 0226, 0255 and 0272, and `qmd
  query "link relativo quebrado no registro de arquivamento promovido para
  docs/references" --all --files --min-score 0.3`, which returned the
  2026-10-12 live-citations characterization, ADR-0277, ADR-0230, ADR-0260
  and the `archive-spec` skill. They confirmed that the outcome already gets
  the rewrite and that promotion keeps one owner in `docs/references/`; no
  page described the `unproven` gap. No Inbox Entry was written, because the
  briefing forbids writing there.
- Exa: the GitHub Docs page above; the CommonMark discussion "What is
  CommonMark's semantic definition for relative links?"
  (https://talk.commonmark.org/t/what-is-commonmarks-semantic-definition-for-relative-links/4023),
  which says the specification leaves relative destinations to the
  publishing platform and that GitHub resolves them against the containing
  file; and the pandoc manual's `rebase_relative_paths` extension
  (https://pandoc.org/demo/example33/8.21-non-default-extensions.html):
  "For each link or image, pandoc will compute the directory of the
  containing file, relative to the working directory, and prepend the
  resulting path to the link or image path." Challenging: pandoc rebases
  links mechanically and never checks the target exists, which would leave a
  removed file's link broken; the decision rebases inside the folder to the
  promoted copy or the Git Locator instead.

## Glossary

- changes: **Legacy Link Rebase**
- changes: **Lenient Legacy Reading**
- adds: **Record Link Repair**

## Decisions

- 2026-10-21 (about 01:00 BRT): the cycle AN decisions through
  AskUserQuestion: scope "V → W → X (Recomendado)" (option text: "Três
  Specs: v0.105 a v0.107. Y fica no backlog como entrada aberta; são só 3
  linhas de ruído no reconcile e não custam nada."); for this Spec "Archive +
  alvo promovido (Recomendado)" (option text: "O archive reescreve os
  `satisfied-by` (Git Locator para arquivo removido, `docs/references/` para
  arquivo promovido), e o sanitize passa a apontar arquivos promovidos para
  `docs/references/`. A exigência da tag fica como está; ela é a âncora do
  Git Locator."); the authorization, the same rules as cycle AM.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AN briefing for Spec 0295 of 2026-10-21 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- Rewrite the `unproven` entries with the Legacy Link Rebase the outcome
  already uses, in `BuildArchiveRecord`, so the Archive Command, the in-Run
  and queue archives and the Legacy Archive Folder conversion share one path.
- Read the entries as written in `_prd.md`, the file that holds the
  declarations.
- Match a promoted copy by file name, because the record keeps only the
  copy's path.
- Let the strict check's fix line name the same copy as the repair, so the
  two never disagree.
- Record the decision as ADR-0308, because it changes what every new record
  holds and what the repair writes.

## Open Questions

None.
