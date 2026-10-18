---
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: active
created: 2026-10-18
surfaces: [backend, cli, docs]
---

# A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers

During cycle AH, on 2026-10-17 and 2026-10-18, the operator observed four
gaps. Each one cost a hand merge, a Run, a QA row or a release fix.

- When round 2 of the pre-PR review found a real defect, no sanctioned path
  remained. Spec 0276's delivery was merged by hand.
- A Verification named a file that no section of its Task declared, and the
  Run failed twice.
- Two QA gates asked for work no Task did. The proposed judge question for
  it was measured and rejected before the Run (Validation, ADR-0287).
- The release found the Setup Manifest stale after a Makefile edit, for the
  fourth time.

After this Spec:

- A fix of round-2 findings gets exactly one Correction Review of its delta and
  then proceeds or parks.
- The Spec Consistency Check refuses a Verification path the plan does not
  provide.
- Every Task settlement reports a stale Setup Manifest.

## Sources

There are no open Backlog Entries or Findings. The maintainer named these four
sources for cycle AH item G on 2026-10-18. Each is an operator observation,
cited here with its evidence. This Spec creates no Backlog Entry or Finding.

1. No sanctioned path after a real round-2 finding (2026-10-17, Spec 0276,
   pull request #568).
   - The pre-PR review round 2 found two real defects in the operator's CI fix
     `a4551833`: a retry with a sleep on `ETXTBSY`, a timing allowance that
     ADR-0271 rejects, and a cancelled wait that returned a stale error.
   - The Delivery Queue parked the item `corrective-spec-required`.
   - The operator replaced the fix structurally (`c7148f13`, the fork lock).
     They reviewed the delta in a scratch clone with
     `roundfix review --base f497cce8`. It found one more defect, fixed in
     `938edc44`, and was then clean.
   - The operator ran `make verify-changed` and merged pull request #568 by
     hand with the maintainer's approval.
   - Those commits are on the delivery branch of Spec 0276, not on `main`. The
     squash commit `23ab0904` carries their result.
2. A Verification named a file that no section declared (2026-10-17, Spec
   0274).
   - Pull request #559 moved two contract tests into existing files.
   - It left task_03's Verification grepping
     `internal/docscontract/glossary_contract_test.go`. The Task no longer
     declared that file, and the file did not exist.
   - The Run failed twice with `grep: … No such file or directory`. Pull
     request #560 fixed the plan.
3. QA criteria that no Task delivers (2026-10-17). Measured and not adopted;
   see Validation and ADR-0287.
   - Spec 0272's QA asked for "exactly 243" files under the History Root, a
     count that Spec 0273's merged record moved to 244. Pull request #554
     rewrote the criterion.
   - Spec 0274's QA asked the glossary to name two check codes that no Task
     added. Pull request #561 rewrote the criterion.
   - Each failed one QA row because the plan's QA requirement exceeded its
     Tasks.
4. The release found the Setup Manifest stale (2026-10-18, v0.88.0).
   - After Spec 0277 merged (`a5d5a7e1`), the release checks showed `baseline:
     plan_ready: 1 file change`.
   - The Managed Refresh rewrote the two `declarationDigest` values of
     `docs/agents/setup-context.json`, and the release pull request #573
     carried the rewrite.
   - Both digests hash the whole Makefile. Spec 0277 changed one Makefile line
     (`REPO_CONTRACT_TESTS`).
   - The same drift followed Specs 0247, 0248 and 0259 and was fixed by
     releases v0.58.0, v0.59.0 and v0.72.0.

## Prerequisites

None. Spec 0277 merged as pull request #572, and v0.88.0 is released from the
same tree. This Spec releases as v0.89.0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; the review record's lineage gains one boolean and
  the Spec Consistency Check gains one code. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; a Correction Review calls the configured
  review provider exactly as round 2 does, and every test runs on temporary
  repositories with a fake reviewer and sends nothing to a provider. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the review lineage, the retry proof,
  the Spec Consistency Check and the documentation contracts each
  follow recorded decisions. ADR-0286 (this Spec): "A fixed round-2 finding
  gets one Correction Review." ADR-0287 (this Spec): "The Spec Consistency
  Check reports a Verification path the plan does not provide." ADR-0197,
  which ADR-0286 amends: "A pre-PR reviewer lineage spans at most two rounds".
  ADR-0272, which ADR-0286 amends: "A review record survives a fix, and a
  correction inside the archived Spec's declared paths returns to review".
  ADR-0103, which ADR-0287 amends: "An applied refresh republishes the Setup
  Manifest so the update converges". ADR-0165: "A blocking review after
  archive parks publication for a corrective Spec". ADR-0196: "A pre-PR review
  finding parks only after validation". ADR-0169: "The pre-PR review diffs the
  candidate from its merge base". ADR-0282: "a large candidate is reviewed in
  segments". ADR-0233: "a review-only correction returns to review". ADR-0279:
  "a dismissal survives a fix that leaves its anchor". ADR-0237: "A Delivery
  Retry records a merge made outside the queue". ADR-0117: "A defect is checked by the
  stage that can produce it". ADR-0252: "The selective gate runs the
  repository contracts a change makes relevant". ADR-0253: "Every Repository
  Contract Test runs on main and before a release". ADR-0271: "A test bounds
  by real time only what it asserts". ADR-0189: "A version names one
  content." ADR-0280: "A guide Task records its own fingerprint." ADR-0256:
  "A release waits for the skills that describe a changed surface". ADR-0259:
  "A test runs in parallel unless it names why it cannot". ADR-0285: "A
  repository contract compiles every test for Windows." ADR-0104: "A Spec
  accepts on evidence it did not author". ADR-0276: "A row the authorization
  keeps offline qualifies, and a byte invariant compares with the Delivery
  Base". ADR-0265: "Each claim carries its proof as a section." Set no obligation here, because this Spec changes none of their rules: ADR-0173 (History Citations of a Baseline Plan); ADR-0182 (settling on the gate's facts); ADR-0209 (the source-grouping question); ADR-0231 (the Jev Ceiling in User Config); ADR-0239 (each OpenRouter stage's key); ADR-0240 (the QA partial policy); ADR-0250 (module versions and the coverage record's platforms); ADR-0273 (no secret leaves a Run, which ADR-0283 amended); ADR-0278 (retirement deletes); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0177 (the Relocation Citation scan reads through a repository root); ADR-0208 (work items that share a context share a Spec, and an open item is extended); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0235 (OpenAI and Anthropic models run only through the codex and claude subscriptions); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0238 (a light tier runs low-complexity Tasks on an open model); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0257 (inside a Run the Daemon is the only full gate, and the Baseline states each rule once); ADR-0200, ADR-0201 and ADR-0283 (the advisory judge, its data boundary and its Run Marker, which this Spec leaves unchanged after measuring and rejecting a new question); ADR-0281 (rebuilt executables in tests). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH approval of 2026-10-16 through AskUserQuestion, option "Governed Paths
  declarados" ("\"Concedo\" como em AE/AF/AG para os Governed Paths que cada
  Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."), and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile,
  workflow, Go module, lint, formatter or Baseline module change.

## Goals

- After round 2 finds a real defect and every finding is disposed as fixed,
  the fix gets exactly one bounded review of its delta through
  `roundfix deliver retry`. It then proceeds as after a clean round 2, or parks
  for a corrective Spec as today, with no review loop.
- The Spec Consistency Check refuses a Verification command of a pending Task
  that names a repository path the plan neither declares, creates nor finds,
  before any Run starts.
- Every Task settlement fails when the Setup Manifest no longer matches a
  fresh Managed Refresh, so the release never finds it stale.

## User Stories

1. As the operator of a delivery whose round-2 review found a real defect, I
   want `roundfix deliver retry` to review my fix once and carry on, so that I
   never review the delta by hand and merge outside the queue.
2. As a Spec author, I want the check to refuse a Verification that names a
   file my plan does not provide, so that a Run never fails on
   `No such file or directory`.
3. As the maintainer, I want a Makefile edit to fail its own Task's settlement
   until the Setup Manifest is regenerated, so that a release does not carry
   that fix.

## Core Features

1. Correction Review in the review lineage.
   - Applies when round 2 of a Reviewer Lineage ended with findings, every
     standing round-2 finding has exactly one valid `fixed` disposition whose
     commit the current head contains, and the head descends from the round-2
     head.
   - The next `roundfix review` runs one **Correction Review** of the delta
     from the round-2 head, instead of closing at the ceiling.
   - It records `correction: true` in its lineage and never repeats.
2. Retry admits the fix.
   - The trigger is `roundfix deliver retry` on an item parked
     `corrective-spec-required` by a round-2 review, where every standing
     finding is disposed as fixed between the parked candidate and the item
     head.
   - The retry returns the item to review whatever paths the fix touched.
   - The Correction Review then decides.
3. Undeclared Verification path.
   - `SC-VERIFY-PATH-UNDECLARED` is an error from the Tasks stage.
   - It names the Task, the command's line and the path.
4. Setup Manifest contract. A test in an existing always-relevant
   documentation contract file fails when a Managed Refresh plan of the
   repository proposes any change, and names the regeneration command.
5. Documentation.
   - The `review`, `deliver` and `spec` user guides, the Roundfix skill
     references that cover them, and the `write-tasks` skill describe each
     change.
   - The glossary adds **Correction Review** and revises **Reviewer Lineage**.

## Non-Goals / Out of Scope

- A second Correction Review, a review after a dismissal-only disposition, or
  any change to round 1, round 2, Review Segments or the review bound.
- Changing the corrective-Spec park itself (ADR-0165), validation of findings
  (ADR-0196), the Delivery Queue's repository gate, required checks or merge
  order.
- Any judge question, threshold, transport or key, or the Jev Ceiling: the
  QA-delivery question was measured and not adopted (Validation, ADR-0287).
- Editing the Makefile, a CI workflow, `REPO_CONTRACT_TESTS`, or adding a
  Repository Contract Test file.
- Changing the release plan, the Managed Refresh or the Setup Manifest format.
- Checking Verification paths of completed, archived or QA Tasks.

## Success Metrics

1. Success Metric: on a temporary repository whose round-2 review left two
   findings disposed as fixed, a fake reviewer that counts its calls sees
   exactly one Correction Review of the round-2-head-to-head delta. A clean
   answer gives outcome `reviewed` with `correction: true`, and the Delivery
   Queue advances. A findings answer parks as round-2 findings do. A further
   review of the same lineage never calls the reviewer.
2. Success Metric: a retry of an item parked `corrective-spec-required` by a
   round-2 review, whose fix changes a file no Task of the archived Spec
   declared and whose every finding is disposed as fixed, returns the item to
   `reviewing`. The same fix with one finding dismissed is refused with
   today's reason.
3. Success Metric: the check raises `SC-VERIFY-PATH-UNDECLARED` for the 0274
   shape of pull request #559 and no finding for the repaired plan. On the
   active corpus it counts 0, as the corpus golden records.
4. Success Metric: on a copy of the repository with one byte appended to the
   Makefile, the Setup Manifest contract fails naming
   `docs/agents/setup-context.json`, and on the unchanged tree it passes in
   well under 10 s.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-18.

- This repository's history, which no Task of this Spec authors: commits
  `ae56aba0`, `26a22f6d`, `db281987` and `a5d5a7e1` each changed the Makefile
  (14, 11, 8 and 1 added lines). The next release commits `70536eb9` (v0.58.0),
  `2a80ee3c` (v0.59.0), `8e88e9b0` (v0.72.0) and `9471f698` (v0.88.0) each
  changed only the two `declarationDigest` lines of
  `docs/agents/setup-context.json` (2 added, 2 removed). The QA gate reads this
  with `git show --numstat`, offline.
- GitHub Docs, "About protected branches"
  (https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches):
  "the \"last reviewer\" can check that the latest set of changes incorporates
  feedback from other reviews, and does not add new, unreviewed content." A
  fix that follows a review is itself reviewed, which is the Correction
  Review's premise.
- Google Engineering Practices, "Speed of Code Reviews"
  (https://google.github.io/eng-practices/review/reviewer/speed.html): "a
  typical CL should get multiple rounds of review (if needed) within a single
  day." Supporting: a bounded further round is ordinary review practice.
  Challenging: the same guide favors approving "LGTM With Comments" when "the
  reviewer is confident that the developer will appropriately address all the
  reviewer's remaining comments", which is today's `ceiling-closed`; the
  Correction Review leaves that path to dismissals.

## Validation

The decisions were checked on 2026-10-18 against these sources.

- Secondbrain.
  - Read `wiki/index.md`, then ran three `qmd query` runs from the Secondbrain
    root:
    - "pre-PR review round ceiling correction of a real finding after round 2"
    - "stale setup manifest declaration digest after Makefile change release"
    - "advisory model judges whether QA criteria are delivered by
      implementation tasks"
  - They returned mirrors of ADR-0197, ADR-0103, the release runbook, the
    `write-tasks` skill and adopters' `autonomous-work` guides. All agree with
    the repository and hold no decision on a review after the ceiling, a
    manifest check per settlement or a QA-delivery question.
  - The pending inbox `inbox/roundfix/` holds no entry.
- Exa.
  - Supporting: GitHub Docs and Google's review guide, quoted above.
  - Challenging: Google's "LGTM With Comments", quoted above.
  - No research Inbox Entry was captured; the sources are standard
    references.
- Measurement, source 2. A throwaway local probe applied the path rule to all
  247 archived Specs with Verification (3,187 commands). It read each Spec's
  tree at the parent of the commit that removed its live Task Graph, or at its
  Archive Record's `source_revision`.
  - It raised 34 findings in 11 Specs:
    - 19 named a path the Task created without declaring it, which is the
      rule's own subject.
    - 15 were false positives (0.47% of commands, 1.6% of Specs).
  - It raised none in Specs 0131 to 0277.
  - It raised exactly the 0274 file on the tree after pull request #559 and
    nothing before or after.
- Measurement, source 3. The candidate was an advisory Noul question per QA
  requirement, "Does some requirement in `implementation_tasks` deliver what
  `qa_requirement` expects the delivered work to show?". On 2026-10-18 the
  lead, with the maintainer's authorization, asked it of 74 QA requirements
  through OpenRouter, outside any Run: the two real cases, 0272 requirement 2
  at `123f6524` and 0274 requirement 10 at `f363ad6d`, and every QA
  requirement of Specs 0272 to 0277 at their planning commits. Every answer
  came from `typesafe/jev-1.13-20260917`, with no error.
  - The real cases scored 0.88 and 0.51, towards "delivered". 0272's task_04
    stated the 243 files itself, so that defect was a wrong count, not an
    undelivered requirement.
  - Below 0.3, 14 of the 72 negatives were raised, all requirements the gate
    fulfils by itself (outside evidence, the Pull Request row, the glossary
    check, "MUST NOT accept a row whose only evidence …").
  - Recall 0 of 2 and 19% false alarms at 0.3, with no threshold that
    separates them, so the question was dropped before the Run (ADR-0287).
- Measurement, source 4. The Managed Refresh plan of the tree at `a5d5a7e1`
  proposes one change to `docs/agents/setup-context.json`, and at `9471f698`
  none. The in-process plan takes about 0.1 s, and the documentation contract
  package takes about 1.2 s with it against 1.0 s without.

## Glossary

- adds: **Correction Review**
- changes: **Reviewer Lineage**

## Decisions

- 2026-10-16: which failure rules govern cycle AH — "Regras de AE–AG
  (Recomendado)".
- 2026-10-16: which authority covers the Governed Paths — "Governed Paths
  declarados", option text "\"Concedo\" como em AE/AF/AG para os Governed
  Paths que cada Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."
- 2026-09-30: the skills grant — "considere autorizado a ajustar todas as
  skills se necessário".
- 2026-10-18 (about 09:20 BRT): what to do after v0.88.0 — "Triagem e Spec G
  (Recomendado)", option text "Faço a triagem e, se houver itens, escrevo,
  entrego e lanço uma Spec G (v0.89.0) com até 4 tasks e as mesmas regras.
  Candidatos: um round de correção quando a revisão do round 2 acha algo real
  e a correção é estrutural, e um detector de critério da QA sem task que o
  entregue. Paro na proposta seguinte."
- 2026-10-18: what to do with a real round-2 finding — "Corrigir e mergear
  (Recomendado)".
- 2026-10-18, relayed by the lead: the maintainer authorized the source-3
  validation calls and, before the Run, amending the Spec on their result; on
  the measured result the lead decided to remove the QA-delivery judgment.

Autonomous route: no person present. Decisions were taken from the lead's
briefing of 2026-10-18 for Spec 0278 and from the code and evidence
characterized in the TechSpec.

Taken from those sources, with no further interview:

- A fixed round-2 finding opens one Correction Review, recorded in the
  lineage, and never a second. See ADR-0286.
- The retry of a round-2 park whose findings are all fixed skips ADR-0272's
  declared-path bound, because the Correction Review reads the fix. See
  ADR-0286.
- `SC-VERIFY-PATH-UNDECLARED` is an error with a conservative path rule. See
  ADR-0287.
- The QA-delivery question is not adopted: measured on 2026-10-18 it found
  neither real case and raised 19% of the negatives. See ADR-0287.
- The Setup Manifest check is a test in an existing always-relevant
  documentation contract file. See ADR-0287.

## Open Questions

None.
