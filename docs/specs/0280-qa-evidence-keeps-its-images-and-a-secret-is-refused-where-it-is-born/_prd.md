---
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: active
created: 2026-10-18
surfaces: [backend, cli, docs]
---

# QA evidence keeps its images, and a secret is refused where it is born

On 2026-10-18 the Conexus repository delivered its first Spec with Roundfix
0.87.0 and lost two full cycles to Roundfix itself.

- The QA gate saved screenshots as the `qa-gate` skill asks, the Daemon
  committed them, and the next gate's mechanical stage refused every PNG as a
  built binary. One QA run, about 11 minutes and 5.7 million tokens, found no
  behavior defect and still failed.
- A test fixture that every Task committed matched a Secret Pattern only at
  the final push. The printed next action asked for a history rewrite, and
  the retry of the archived item then refused the rewritten head.

After this Spec:

- PNG, JPEG and SVG evidence is readable evidence. The QA Report commit
  refuses any other binary under `qa/evidence/`, so the gate that creates the
  defect reports it.
- The Secret Scan reads every Task commit, so a fixture value is refused when
  it is born.
- A Delivery Retry accepts a **Secret Rewrite** of a `secret-found` item, and
  the next action names the rewrite that works.

## Sources

Two Secondbrain Inbox Entries from the Conexus repository, captured manually
on 2026-10-18, were promoted to Findings and adopted into this Spec's
`references/` (see `references/_index.md`). The lead moves the Inbox Entries
to `_triaged/` after this plan merges.

1. Finding `2026-10-18-the-qa-gate-commits-screenshots-the-mechanical-audit-then-refuses.md`,
   from `inbox/roundfix/2026-10-18-o-qa-gate-versiona-screenshot-que-a-auditoria-mecanica-recusa.md`.
   - Conexus Spec `0001-usuario-verificado-com-segundo-fator`: the third QA run
     raised 67 `QA-EVIDENCE-PATH` findings, one per PNG the first two runs had
     committed, each "is a built binary".
   - The Conexus workaround, commit `d132944`, removed the PNGs and added a
     text-only evidence requirement to the QA Task. The fourth run, with that
     requirement, still committed 26 new PNGs in the Daemon's report commit
     `afe23e9`.
2. Finding `2026-10-18-secret-found-asks-for-a-rewrite-the-retry-of-an-archived-item-refuses.md`,
   from `inbox/roundfix/2026-10-18-secret-found-manda-reescrever-mas-o-retry-recusa-o-head-arquivado.md`.
   - The same Conexus Spec's archived Delivery Queue item parked
     `secret-found` with five `credential-assignment` matches of one test
     fixture.
   - A `git filter-branch` over `6bcddf3..HEAD` was refused by the retry
     because the archived item head differed from the candidate head. The
     operator recovered by resetting to the commit before the archive, which
     lost the Reviewer Lineage.

No open Backlog Entry or unresolved Finding shares this context:
`docs/backlog/` holds none and `docs/findings/` held only these two.

## Prerequisites

Spec 0279 (`0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage`)
is planned and delivers first, as v0.90.0. It changes the `deliver` user guide,
the Roundfix skill references and their version record, which this Spec also
changes, so this Spec names it in the Task Graph's `requires` and releases as
v0.91.0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; the QA Task gains one refusal reason, the Task
  settlement one refusal reason, and the Delivery Retry one acceptance rule.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the Secret Scan reads local files and Git
  objects only, and every test runs on temporary repositories and sends
  nothing to a provider. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the QA evidence rule, the Task commit
  scan and the retry follow recorded decisions. ADR-0291 (this Spec): "QA
  evidence holds text and images, and the QA Report commit refuses any other
  binary". ADR-0292 (this Spec): "A Task commit is scanned for secrets, and a
  retry accepts a Secret Rewrite". ADR-0273, which both amend: "No secret and
  no Spec content leaves the machine from a Run". ADR-0283: "A Run Marker
  keeps the judge silent inside a Run". ADR-0237: "A Delivery Retry records a
  merge made outside the queue". ADR-0272: "A review record survives a fix,
  and a correction inside the archived Spec's declared paths returns to
  review". ADR-0197: "A pre-PR reviewer lineage spans at most two rounds".
  ADR-0117: "A defect is checked by the stage that can produce it". ADR-0189:
  "A version names one content." ADR-0280: "A guide Task records its own
  fingerprint." ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0285: "A repository contract compiles every test for
  Windows." ADR-0125: "A spawned fixture is a compiled binary, never a written
  script". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "Each claim carries its proof as
  a section." ADR-0193: "A Spec names its prerequisites, and the owner waits for
  them". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0182 (a Task settles on the facts its gate will check); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0271 (a test bounds by real time only what it asserts); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  approval of 2026-10-16 through AskUserQuestion, option "Governed Paths
  declarados" ("\"Concedo\" como em AE/AF/AG para os Governed Paths que cada
  Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."), and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born/_authorization.md`;
  bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `skills/implement-task/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module,
  lint, formatter or Baseline module change.

## Goals

- A screenshot the QA gate saves under `qa/evidence/` as PNG, JPEG or SVG is
  committed and accepted by every later gate of the same Spec.
- A binary under `qa/evidence/` that is not an image fails the gate whose
  report commit would stage it, naming the path, and never reaches a commit.
- A Task whose commit would add a line that matches a Secret Pattern settles
  `failed` with the location, before QA, review or archive.
- After a `secret-found` park, the operator can follow the printed next
  action, rewrite the item branch, and have `roundfix deliver retry` return
  the item to review without undoing the archive.

## User Stories

1. As the QA gate of a frontend Spec, I want my screenshots to count as
   evidence, so that the next gate does not fail on files the flow committed.
2. As the maintainer, I want a non-image binary refused at the report commit
   that would stage it, so that the defect costs one gate, not two.
3. As a Spec author whose Task adds a test fixture, I want the fixture refused
   by the Task that adds it, so that I fix one value instead of rewriting a
   delivered branch.
4. As the operator of a `secret-found` item, I want the printed next action to
   lead to a retry that succeeds, so that I never reset the branch to before
   its archive.

## Core Features

1. Readable QA evidence.
   - One classifier reads the first 8000 bytes of a file: no NUL byte is text,
     an SVG included; the PNG signature is a PNG image; the JPEG Start Of
     Image marker followed by a marker is a JPEG image; anything else is a
     binary that is not an image. Content decides, never the name.
   - The mechanical stage raises `QA-EVIDENCE-PATH` for a tracked evidence
     file only when it is a gitlink or a non-image binary.
2. The QA Report commit refuses a non-image binary.
   - Applies to every path under the report's `evidence/` directory that the
     QA Report commit would stage, after the Secret Scan finds nothing.
   - The QA Task settles `failed` with the paths named, and no QA Report
     commit is made.
3. The Task commit Secret Scan.
   - Applies to the lines a Task commit would add relative to `HEAD`, after
     the Task's Verification passes and before it settles.
   - A Secret Match settles the Task `failed` with each location and makes no
     Task commit.
4. The Secret Rewrite.
   - Applies to `roundfix deliver retry` of an archived item parked
     `secret-found` whose head moved without descending from the parked
     candidate.
   - It accepts the head when the rewrite starts at or after the item's Run
     start, and the trees differ only at lines of the files the blocker names
     that held a Secret Match at the candidate and hold none at the head.
   - The item returns to `reviewing`; any other rewrite is refused with
     today's reason and the rule that failed.
   - `secret-found` prints a next action that describes that rewrite.
5. Documentation.
   - The `implement` and `deliver` user guides, the Roundfix `implement` and
     `deliver` references, the `qa-gate` skill and the `implement-task` skill
     describe each change, including fixture values the
     `credential-assignment` pattern does not match.
   - The glossary adds **Secret Rewrite** and revises **Secret Scan** and
     **Mechanical Refusal Code**.

## Non-Goals / Out of Scope

- Accepting a PDF, a video, a trace archive or any other non-image binary as
  QA evidence, or bounding an image's size.
- Changing any Secret Pattern, exempting test files or fixture-looking
  values, or retaining any part of a matched value.
- Giving the Agent a repair turn after a Task Secret Match, or changing the
  Verification retry policy.
- Carrying the Reviewer Lineage across a Secret Rewrite: the rewritten head is
  reviewed from round 1.
- Accepting a rewrite that deletes or adds a line, changes another file, or
  rewrites commits before the item's Run start.
- Changing the Final Push or the Delivery Queue push scans, the retry of any
  other blocker, or the mechanical stage's other codes.

## Success Metrics

1. Success Metric: on a temporary repository, tracked PNG, JPEG and SVG files
   under a QA Report's `evidence/` directory raise no `QA-EVIDENCE-PATH`,
   while an executable, a SQLite database, a PDF and a file named `.png` that
   holds an executable each raise one finding naming the path.
2. Success Metric: a QA run whose evidence holds a PNG commits it with the
   report, and a QA run whose evidence holds a non-image binary settles the QA
   Task `failed` with `QA evidence carries a non-image binary: <path>` and
   creates no QA Report commit.
3. Success Metric: a Task whose change adds a `credential-assignment` line
   settles `failed` with `Task commit carries a secret: <path>:<line>:
   credential-assignment` and creates no Task commit. A Task that edits
   another line of a file whose committed line already matches, and a Task
   that adds a letters-only fixture value, each commit.
4. Success Metric: a retry of an archived `secret-found` item whose rewrite
   replaced only the matched lines returns the item to `reviewing` with the
   rewritten head as its newest candidate. Rewrites that change one more
   line, change an unnamed file, or start before the Run start are each
   refused with the archived-head reason followed by the failed rule.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-18.

- This repository's history, which no Task of this Spec authors. At
  `513b22eb`, `docs/history/specs/0005-tui-cockpit/design/roundfix-01.png` is
  a PNG and
  `docs/history/specs/0035-agent-selection-profiles/references/artificial-analysis-coding-agent.jpeg`
  is a JPEG, while the QA evidence of Spec 0090 held SQLite Run Databases
  (`docs/history/specs/0090-a-gate-that-could-have-failed/qa/evidence/2026-08-09-spec-0090/positive-home/.roundfix/roundfix.db`)
  and the QA evidence of Spec 0209 held a PDF
  (`docs/history/specs/0209-sources-that-share-a-context-share-a-spec/qa/evidence/2026-10-01-gate/zhang-author.pdf`).
  The QA gate reads them offline with `git show` and records the classifier's
  answer for each.
- W3C, "Portable Network Graphics (PNG) Specification (Third Edition)",
  section 5.2 (https://www.w3.org/TR/png-3/): "The first eight bytes of a PNG
  datastream always contain the following hexadecimal values: 89 50 4E 47 0D
  0A 1A 0A".
- Library of Congress, "JFIF, JPEG File Interchange Format, Version 1.02"
  (https://www.loc.gov/preservation/digital/formats/fdd/fdd000018.shtml): "The
  first two bytes of every JPEG stream are the Start Of Image (SOI) marker
  values FFh D8h."
- GitHub Docs, "Working with push protection from the command line"
  (https://docs.github.com/en/code-security/how-tos/secure-your-secrets/work-with-leak-prevention/push-protection-on-the-command-line):
  "To resolve a blocked push, you must remove the secret from all of the
  commits it appears in." For a secret in an earlier commit the page
  prescribes an interactive rebase from that commit, which is the rewrite a
  Secret Rewrite accepts.

## Validation

The decisions were checked on 2026-10-18 against these sources.

- Secondbrain.
  - Read `wiki/index.md`, then ran two `qmd query` runs from the Secondbrain
    root:
    - "QA evidence screenshot binary refused mechanical audit"
    - "secret found push history rewrite retry delivery queue fixture
      credential"
  - The first returned
    `projects/fiscus/mirror/docs/history/backlog/2026-08-28-o-gate-de-qa-commita-binarios-que-ele-mesmo-recusa-depois.md`:
    on 2026-08-28 the Fiscus QA gate committed PNG and webm captures that the
    next mechanical stage refused, and the entry closed because the loop had
    not recurred since 2026-09-13. Conexus shows it recurring, which supports
    fixing the classifier rather than the gate's habits.
  - Both queries returned the mirrors of ADR-0273, the `deliver` guide and
    the Roundfix `deliver` reference, which agree with this repository and
    hold no decision on images or on a rewritten head.
- Exa.
  - Supporting: the W3C PNG specification, the Library of Congress JFIF
    description and GitHub's push protection page, quoted above.
  - Challenging: GitHub lets a person bypass a block, "If the secret is only
    used in tests and poses no threat, click It's used in tests." Roundfix
    keeps no such human bypass inside a Run, so it documents fixture values
    the pattern does not match instead of exempting tests.
  - No research Inbox Entry was captured; the sources are standard
    references.
- Measurement. A throwaway `go test -overlay` probe ran the Secret Scan over
  every tracked file outside `docs/history/` at `2ef4f080` and found no Secret
  Match, so scanning this repository's Task commits refuses nothing that is
  already committed.

## Glossary

- adds: **Secret Rewrite**
- changes: **Secret Scan**
- changes: **Mechanical Refusal Code**

## Decisions

- 2026-10-16: which failure rules govern the cycle — "Regras de AE–AG
  (Recomendado)".
- 2026-10-16: which authority covers the Governed Paths — "Governed Paths
  declarados", option text "\"Concedo\" como em AE/AF/AG para os Governed
  Paths que cada Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".
- 2026-10-18 (about 14:00 BRT): where this work goes — "Spec I depois da 0279
  (Recomendado)".
- 2026-10-18 (about 14:00 BRT): what to do with PNG evidence — "Aceitar
  imagens (Recomendado)", option text "O detector aceita PNG, JPEG e SVG como
  evidência legível e recusa só binário de build ou executável. O Daemon
  recusa no commit do relatório qualquer binário que não seja imagem."
- 2026-10-18 (about 14:00 BRT): what to do with `secret-found` — "Retry aceita
  head reescrito", option text "Para secret-found, o retry aceita um head
  reescrito que preserva a árvore final, exceto as linhas apontadas, e a
  próxima ação impressa explica o caminho real."; "Scan em cada Task", option
  text "O Secret Scan roda no commit de cada Task, como já roda no do
  relatório de QA, e recusa a fixture quando ela nasce."; and "Documentar
  fixture", option text "Documentar no guia e na skill como escrever um valor
  de fixture que não case a regra, sem afrouxar o detector."

Autonomous route: no person present. Decisions were taken from the lead's
briefing of 2026-10-18 for Spec 0280 and from the code and evidence
characterized in the TechSpec.

Taken from those sources, with no further interview:

- The classifier decides by content: text, the PNG signature, or the JPEG
  Start Of Image marker; an SVG is text. See ADR-0291.
- The QA Report commit refuses a non-image binary the way it refuses a Secret
  Match, and no image size bound is added because the code has no natural
  one. See ADR-0291.
- The Task commit scan reads the lines the commit adds, which is what the
  Final Push reads. See ADR-0292.
- A Secret Rewrite is compared line by line at the lines that match at the
  parked candidate, the strictest comparison a Secret Match allows, and it
  returns the item to `reviewing`. See ADR-0292.
- The fixture guidance also goes into the `implement-task` skill, because the
  Agent that writes a fixture reads that skill. See ADR-0292.

## Open Questions

None.
