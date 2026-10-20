---
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# The history conversion rewrites every live citation

The history conversion of `roundfix baseline update` and `roundfix history
sanitize` now rewrites every citation it lists, in every tracked text file
outside the skill trees and `testdata`, into the Git Locator of the removed
file. A date-named ADR is named by its file name and dated with the day it was
superseded. A link to a file a pending Task will create is left alone, and a
link that was already broken is resolved by its file name. `roundfix spec
check` reports a citation of a removed history file that is still not a Git
Locator.

Minimal PRD for a bug fix, minted by `write-techspec`. It carries the contract
the downstream skills parse; the TechSpec holds the design.

On 2026-10-19 an adopter ran `roundfix baseline update` 0.91.0 with the
history conversion. It removed 27 files from the retired history families,
listed 45 History Citations and rewrote 30 citations in 7 files. Four defects
appeared, and the check called the result clean.

1. Citations in ADRs, agent guides outside the setup markers, reference
   pages, infrastructure READMEs and a Makefile comment were listed or not
   found at all, and stayed pointing at removed files: 18 of them. They were
   paths between backticks, relative links and relative paths between
   backticks.
2. A date-named ADR's successor received a `Supersedes` line that named it
   `ADR-` and the year `2026` of its file name, dated with the day of the
   update instead of the day the predecessor was superseded.
3. Four links of an active Spec to a file its pending Task declares in
   `creates:` became text marked not reachable.
4. Two relative links that were already wrong became text marked not
   reachable instead of the Git Locator of the removed file with that name.

## Sources

One Secondbrain Inbox Entry, captured manually from the Pantheon repository on
2026-10-19, is adopted verbatim into this Spec's `references/` (see
`references/_index.md`). The lead moves the Inbox Entry to `_triaged/` after
this plan merges; this Spec never writes to the Secondbrain.

1. [`references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt`](references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt),
   from `inbox/roundfix/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.md`.
   It records the four defects, the plan digest, the tool and hand-fix
   commits on the adopter's branch, the grep that found the 18 citations, and
   the `--no-history` observation. The copy keeps the Inbox Entry's bytes under a
   `.txt` name, because `spec check` reads every Markdown file of a Spec for
   decision-record citations and would take the wrong identifier the entry
   quotes for a record this Spec must list.

No open Backlog Entry or unresolved Finding shares this context:
`docs/backlog/` and `docs/findings/` hold none at `9448aa64`. Specs 0285 and
0287 of the same cycle own the preflight profile proof and three report
counts; neither touches the history conversion or `SC-HISTORY-RETAINED`.

## Prerequisites

None. This Spec is authored against `9448aa64` (v0.95.0) and requires no
other Spec. Specs 0285 and 0287, authored beside it, may merge first and
raise the Roundfix skill version or re-record the Behavior Surface Record,
so this plan is rebased and checked again before its Run.

## Project Constraints

- Identifier strategy: applicable — the lineage line identifies a retired
  ADR: `ADR-NNNN` for a numbered ADR and the file name without `.md`
  otherwise, the identifier the repository's own `superseded_by` uses. No
  new identifier, slug or record schema is minted. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the conversion and the check read only
  the repository and its Git history, and every test uses temporary
  repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the conversion, the check and their
  documentation follow recorded decisions. ADR-0299 (this Spec): "The history
  conversion rewrites every citation it lists, and the check reports a
  removed history file cited without a Git Locator". ADR-0278, which ADR-0299
  extends: "A retirement deletes, and the History Root holds only Archive
  Records". ADR-0282, which ADR-0299 amends: "The history cleanup names a
  superseded ADR in its successor, and a large candidate is reviewed in
  segments". ADR-0260: "An archive relinks live citations, and a Baseline
  update lists the citations it breaks". ADR-0254: "Baseline update
  sanitizes pending history in the change it plans". ADR-0295: "Reachability
  reads the fetched default branch, and origin names it when origin/HEAD is
  unset". ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts". ADR-0285:
  "Every test file compiles for Windows, and a repository contract proves
  it". ADR-0189: "An owned skill's version names its content, and its minimum
  is the bundle". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0280: "A guide Task records its own surface
  fingerprint, a PRD records who decided, and the glossary is held to its own
  rule". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240:
  "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A
  row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "A Spec proves what it claims,
  and every active Spec is held to it". ADR-0183: "An attributed claim
  carries a receipt the check proves". ADR-0184: "A TechSpec states a command
  surface as a transcript". Set no obligation here, because this Spec changes none of
  their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120
  (retired documentation lives under one history root); ADR-0165 (a blocking
  review after archive parks publication for a corrective Spec); ADR-0167
  (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169
  (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a
  pre-PR review finding parks only after validation); ADR-0197 (a pre-PR
  reviewer lineage spans at most two rounds); ADR-0229 (an operator archive
  resumes any park from the Run start); ADR-0233 (a skill version raise is
  regenerated at merge, and a review-only correction returns to review);
  ADR-0237 (a Delivery Retry records a merge made outside the queue);
  ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays
  in Git); ADR-0248 (existing history is sanitized in batches after a
  history-full tag); ADR-0250 (a module version is chosen when it is
  recorded, and the coverage record lists every platform); ADR-0251 (a
  Legacy Archive Folder is read leniently, and a failed QA keeps its
  verdict); ADR-0258 (a Spec is authored against the QA rerun classes, and a
  retirement writes reduced history); ADR-0261 (an Archive Record names its
  pull request, and readers fetch the pull request head); ADR-0262 (a
  cancelled Spec leaves an Archive Record and returns its sources); ADR-0269
  (archive and cancel refuse ignored Spec files and finish a failed removal
  forward); ADR-0272 (a review record survives a fix, and a correction
  inside the archived Spec's declared paths returns to review); ADR-0277 (a
  Legacy Archive Folder conversion keeps the summary, links and dates its
  folder held); ADR-0279 (status reports a Queue Owner proven gone, and a
  dismissal survives a fix that leaves its anchor); ADR-0281 (a test
  replaces an executable it rebuilds instead of writing into it); ADR-0286
  (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290
  (merged items give back their Item Binaries, and Run Worktrees leave only
  by proof). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 through AskUserQuestion, option "Sim, mesmas
  regras do AJ (Recomendado)", which carries the cycle AJ grant "Mexer em
  Governed Paths declarados" ("'Concedo' para os caminhos que cada Spec
  declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0286-the-history-conversion-rewrites-every-live-citation/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/spec.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, Go module, lint, formatter, Baseline module or
  Project Config file change.

## Goals

- No citation the conversion lists still points at a file it removed.
- A citation outside Markdown, and a relative path between backticks, is
  found and rewritten like a Markdown link.
- A lineage line names a retired ADR by the identifier the repository uses
  and dates it with the day it was superseded.
- A link to a file a pending Task will create survives the conversion and the
  check.
- A link that was already broken gets the Git Locator of the removed file
  with that name when exactly one exists.
- `roundfix spec check --strict` fails while a citation of a removed history
  file is not a Git Locator, and its diagnostic prints the locator to use.

## Core Features

1. Citation Sources. The conversion looks for citations in every tracked
   text file outside the skill trees and `testdata`, including prose paths
   that start with `./` or `../`, and rewrites each one it lists; outside
   Markdown it writes a plain `<path>@<commit>`.
2. Lineage identity and date. The successor's line names a date-named ADR by
   its file name without `.md` and dates it with the day the ADR was
   superseded.
3. Pending Task targets and broken links. `live-links` and the check leave a
   link to a file a pending Task declares in `creates:`, and resolve a missing
   link by its file name among the removed history files.
4. The check. `SC-HISTORY-RETAINED` reports each citation of a Removed
   History Path that is not a Git Locator, with exemptions for files under
   `docs/history/` and pending `creates:` targets.
5. Documentation. The `baseline`, `history` and `spec` guides, the Roundfix
   skill references and `CONTEXT.md` describe the four corrections, the check
   and the `--no-history` interaction.

## Non-Goals / Out of Scope

- Repairing citations an earlier update already left broken in an adopter:
  the check reports them with the locator to use, and the Stale History
  Citation report already offers the restore command.
- Rewriting citations in skill trees, `testdata` or binary files.
- Rewriting inside the Archive Relink of `roundfix archive`, or changing its
  warnings.
- Making `--no-history` keep `SC-HISTORY-RETAINED` green while history is
  pending: the guide's retirement rule is right; the documentation names the
  path to green.
- Changing a Baseline module, the Plan Digest's fields, or the JSON field set.

## Success Metrics

1. Success Metric: on a temporary repository of the adopter's shape, the plan
   lists and rewrites every citation of a removed file in an ADR, an agent
   guide, a reference page, a README and a Makefile comment, including
   relative paths between backticks. After apply, no citation of a removed
   path outside `docs/history/` lacks a Git Locator, and no file under a skill
   tree or `testdata` changes.
2. Success Metric: a date-named superseded ADR whose `updated_at` is
   2026-10-01 gives its successor `Supersedes <file name without .md>
   (2026-10-01)`, whatever day the plan runs.
3. Success Metric: a link to a pending Task's `creates:` target is neither
   rewritten nor reported, and an already broken link whose file name matches
   exactly one removed history file gets that file's Git Locator.
4. Success Metric: on a tree where an earlier conversion left such
   citations, `roundfix spec check --strict` exits 1 with one
   `SC-HISTORY-RETAINED` error per citation, each naming the file, the line
   and the locator to use; Archive Records and pending `creates:` targets
   are not reported.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The adopter's run, a measurement this Spec did not design, recorded in the
  adopted Inbox Entry: 45 citations listed, 30 rewritten, 18 left in files
  outside `docs/history`, a lineage line naming `ADR-` and the year `2026`,
  dated 2026-10-16 against an `updated_at` of 2026-10-01, and `spec check
  --strict` exiting 0. The QA gate reads the committed copy under
  `references/`.
- MADR, "Use Dashes in Filenames"
  (https://adr.github.io/madr/decisions/0005-use-dashes-in-filenames.html):
  "`NNNN` provides a unique number, which can be used for referencing in the
  forms", the first of which is the `ADR-` prefix and the number in plain
  text; the other option it considered is
  "`YYYY-MM-DD Title`".
- Log4brains example ADR, "Use the ADR slug as its unique ID"
  (https://ruzickap.github.io/myteam-adr/log4brains/adr/20201016-use-the-adr-slug-as-its-unique-id/):
  "An ADR will be uniquely identified by its slug (ie. its filename without
  the extension)".
- Git's binary test, as quoted from `xdiff-interface.c` in "How to determine
  if Git handles a file as binary or as text?"
  (https://stackoverflow.com/questions/6119956/how-to-determine-if-git-handles-a-file-as-binary-or-as-text):
  "checks for any occurrence of a zero byte (NUL “character”) in the first
  8000 bytes (or the entire length if shorter)".

## Unreachable Acceptance

None. Every criterion is reached on temporary repositories with the built
binary; the published passages are quoted above, and the QA gate records
them as blocked by the Spec's own network prohibition, which never decides a
qualifying `partial`.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Secondbrain, read only. Read `wiki/index.md`, then ran `qmd query
  "baseline update history conversion broken citations git locator" --all
  --files --min-score 0.3` and `qmd query "date-named ADR identifier
  superseded_by lineage" --all --files --min-score 0.3`. They returned this
  repository's mirrored `baseline` guide, its live-citations characterization
  reference, ADR-0260, ADR-0278 and ADR-0282, and unrelated ADRs of other
  projects. None decides how a date-named ADR is identified or which files
  hold citations, so the four corrections stand on the adopter's evidence.
  No Inbox Entry was written, because the briefing forbids writing there.
- Exa: the MADR, Log4brains and Git passages quoted above. Supporting: MADR
  ties `ADR-NNNN` to numbered names only, and Log4brains identifies a
  date-named ADR by its file name. Challenging: Log4brains' file names are
  `YYYYMMDD-<slug>`, not `YYYY-MM-DD-<slug>`, so the identifier rule is
  stated for any name that is not `NNNN-<slug>`, not for one date shape.
- Characterization on 2026-10-20 at `9448aa64`, recorded in the TechSpec: a
  temporary repository of the adopter's shape reproduced all four defects
  through `bin/roundfix history sanitize`, and `bin/roundfix spec check
  --strict` reported none of the citations left behind.

## Glossary

- adds: **Citation Source**
- adds: **Removed History Path**
- changes: **History Citation**
- changes: **Git Locator**

## Decisions

- 2026-10-20 (about 06:30 BRT): the cycle AK plan — "Duas Specs + um lote
  (Recomendado)"; the scope of this Spec — "Corrigir os 4 + o check acusa
  (Recomendado)" (option text: "Reescrever como Git Locator toda citação
  listada, inclusive caminho entre crases e relativo. Não tocar link para
  arquivo que uma Task pendente declara em creates:. Resolver link relativo
  quebrado pelo nome do arquivo. Derivar o ID da ADR com nome por data do
  nome do arquivo ou de superseded_by, com a data da substituição. O spec
  check passa a acusar citação viva quebrada para docs/history. A nota sobre
  --no-history vai para a documentação."); the authorization — "Sim, mesmas
  regras do AJ (Recomendado)".
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".

Autonomous route: no person present; decisions taken from the operator's cycle AK briefing of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- Citation Sources are every tracked text file outside the skill trees and
  `testdata`; every listed citation is rewritten. See ADR-0299.
- The check is an extension of `SC-HISTORY-RETAINED`, not a new rule ID. See
  ADR-0299.
- The lineage date is the retired ADR's `updated_at`, with commit dates as
  fallbacks. See ADR-0299.

## Open Questions

None.
