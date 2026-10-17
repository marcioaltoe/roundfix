---
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: active
created: 2026-10-17
surfaces: [backend, cli, docs]
---

# The history cleanup reaches every adopter, and the review sees a large Spec

On 2026-10-17 the history cleanup that `roundfix baseline update` plans
refused its whole ADR unit in three of four adopters because one superseded
ADR's successor did not name it, and the guide the same update wrote then
turned the adopter's Spec Consistency Check red over every retired ADR left
behind and over links to files an earlier layout removed
([the first Backlog Entry](references/2026-10-17-one-unnamed-superseded-adr-keeps-an-adopters-whole-history.md)).
The same day the Delivery Queue parked Spec 0272 because its pre-PR review
diff, swollen by deleted history and regenerated goldens, exceeded the review
bound before any reviewer read it
([the second Backlog Entry](references/2026-10-17-the-pre-pr-review-refuses-a-spec-that-deletes-retired-history.md)).
Four smaller gaps of the same cycle each cost a Run or a review: a delivery
convention that asks a QA Task for settlement sections
([third](references/2026-10-16-review-convention-c2-asks-a-qa-task-for-settlement-sections.md)),
a plan that edits an owned skill without its version record and still passes
the strict check
([fourth](references/2026-10-16-spec-check-does-not-flag-an-owned-skill-edit-without-its-version-record.md)),
a QA gate whose archive plan contacts the judge its own Task forbids
([fifth](references/2026-10-16-a-qa-gate-archive-plan-contacts-the-judge-an-offline-qa-forbids.md)),
and an item binary rebuilt in place while a copy still runs, which macOS kills
([sixth](references/2026-10-16-an-item-binary-built-in-place-can-be-killed-by-code-signing.md)).
After this Spec an adopter applies its history cleanup and passes the check,
a large Spec is reviewed instead of refused, and each of the four gaps is
closed where it appears.

## Prerequisites

Spec 0275 (release v0.86.0), authored beside this one, changes the delivery
item binary tests this Spec extends and decides ADR-0281, which this Spec
carries from tests to the Delivery Queue. This Spec runs after 0275 merges and
releases as v0.87.0.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one Spec Consistency Check code joins the existing
  code set, and the Pre-PR Review Record gains optional fields that older
  records omit. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; the change removes a judge call from the QA path, and every
  test runs on temporary repositories and a fake reviewer. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0282 (this Spec) decides every
  change: "The cleanup writes the lineage a successor lacks.", "Live Documents
  get their Git Locators.", "A QA Task settles on its QA Report." and "Spec
  check holds an owned skill edit to its version.". ADR-0278: "A retirement
  deletes." and "A link to a deleted file becomes a Git Locator.", which the
  cleanup now completes. ADR-0254: "Refused Units never block.", so a single
  kept ADR is listed and the rest applies. ADR-0260: "The update lists.", the
  citation list the lineage line joins. ADR-0281: "A test that rebuilds an
  executable replaces it.", the rule the Delivery Queue now follows. ADR-0225:
  "A queue item runs the Roundfix binary its branch builds", whose built
  binary is now run from its copy. ADR-0196: "A pre-PR review finding parks
  only after validation", whose Delivery Conventions C2 belongs to. ADR-0197:
  "A pre-PR reviewer lineage spans at most two rounds", so one round covers
  every segment. ADR-0227: "The Claude review prompt is bounded in tokens, not
  bytes", so segments apply only where the byte bound applies. ADR-0273: "The
  judge's keys never enter a Run.", which the offline archive plan no longer
  depends on. ADR-0189: "A version names one content.", the rule the new gap
  holds plans to. ADR-0280: "A guide Task records its own fingerprint.", whose
  gap the new one copies. ADR-0256: "A release waits for the skills that
  describe a changed surface", so each Task records its surfaces. ADR-0104: "A
  Spec accepts on evidence it did not author". ADR-0276: "A row the
  authorization keeps offline qualifies, and a byte invariant compares with the
  Delivery Base". ADR-0169: "The pre-PR review diffs the candidate from its merge base" is kept, since every Review Segment is a range of the candidate measured from that merge base. Set no obligation here, because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0251, ADR-0261, ADR-0262 and ADR-0269 (history families, Archive Records, cancellation and the History Full Tag); ADR-0233 and ADR-0250 (version regeneration at merge and module versions); ADR-0258 (reduced history); ADR-0277 (legacy folder conversion); ADR-0165, ADR-0272 and ADR-0279 (review parks, records and dismissals); ADR-0167, ADR-0229, ADR-0237 and ADR-0240 (QA rows, operator archive and delivery retry). Source:
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
  `docs/specs/0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `docs/references/coverage-record.json`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter or Baseline module change.

## Goals

- An adopter whose ADR log names a successor that does not name back, or
  whose Live Documents link files an earlier layout removed, applies its
  history cleanup and then passes `roundfix spec check --strict`, except for
  each superseded ADR that has no accepted successor, which is named with its
  fix.
- A Spec that deletes retired history or rewrites regenerated files is
  reviewed before its Pull Request instead of refused for size.
- A pre-PR review never asks a QA Task for settlement sections.
- A plan whose Task edits an owned skill without its version record fails
  `roundfix spec check --strict` before a Run starts.
- A QA gate lists its archive candidates without contacting the judge, and an
  item binary build never rewrites a running item binary.

## User Stories

1. As an adopter's maintainer, I want the history cleanup to remove every
   retired ADR it can and to write the lineage line a successor lacks, so that
   one unnamed ADR does not keep the whole history and turn my check red.
2. As an adopter's maintainer, I want links from my Live Documents and user
   guide to files an earlier layout removed rewritten to Git Locators by the
   same cleanup, so that I need no hand edit to pass the check.
3. As an operator of the Delivery Queue, I want a large candidate reviewed in
   bounded segments, with deleted history and generated files left out, so
   that a history Spec is not parked `review-blocked`.
4. As an operator, I want the reviewer told that a QA Task settles on its QA
   Report, so that a completed QA Task without settlement sections is not
   raised as a finding.
5. As a Spec author, I want `spec check` to name a Task that edits an owned
   skill without its mirror, version fields, version record and record step,
   so that the Run does not fail at settlement.
6. As a QA Agent under a Task that forbids provider requests, I want an
   offline archive plan, so that I can list the archive candidates and still
   pass my own gate.
7. As an operator on macOS, I want the Delivery Queue to run a fresh copy of
   each item binary build, so that a build writing its output in place never
   gets a running item binary killed by code signing.

## Core Features

1. Lineage completion. When a superseded ADR names an accepted successor that
   does not name it back, the history plan removes it and adds to the
   successor one dated line that names it and the Git Locator of its text,
   listed with the plan's citations and bound by the Plan Digest. A superseded
   ADR without an accepted successor stays and is refused by its own path with
   the fix; the unit removes every other retired ADR.
2. Live link repair. A new history unit, after the Archive Record repair,
   rewrites every relative link from a Live Document or a user guide page to a
   missing file into its Git Locator, exactly as the Spec Consistency Check's
   history gate prints it. Both the History Sanitize Command and
   `roundfix baseline update` plan it, and the JSON plan lists its pages.
3. Review omissions. The pre-PR review diff also omits, with their reasons in
   the prompt and the record, each file deleted under the History Root and
   each file a Project Config derived-path declaration of the base commit
   names.
4. Review segments. A `codex` round-1 candidate whose diff still exceeds the
   review bound is reviewed in up to four segments, each a **Review Segment**,
   in one review session, and their findings form one Pre-PR Review Record; a single commit
   above the bound, a fifth segment and a round-2 delta above the bound still
   block, each with its reason.
5. Delivery convention C2. The reviewer is told that the Daemon writes
   settlement sections for non-QA Tasks only, and that a QA Task settles on
   its QA Report.
6. Owned skill version gap. The Spec Consistency Check reports a pending
   non-QA Task that declares a file of an owned skill without its mirror, both
   `SKILL.md` files, the version record and the record step, from the Tasks
   stage, with no horizon.
7. Offline archive plan. `roundfix archive <slug> --plan --offline` prints
   the same cut without Archive Advice, says the judge was skipped, and sends
   nothing; the `qa-gate` skill uses it.
8. Item binary copy. After each item binary build, the Delivery Queue copies
   the built file to a new file under the item's artifacts by writing a
   sibling temporary file and renaming it, then probes and runs that copy.
9. Documentation. The `history`, `baseline`, `review`, `spec`, `archive` and
   `deliver` user guides, the Roundfix skill references, the `write-tasks` and
   `qa-gate` skills, and the glossary describe each change.

## Non-Goals / Out of Scope

- Editing any adopter repository. The maintainer decided that the operator
  sends each adopter an Inbox Entry with its work instead; no Task touches
  vortex, fluxus, oraculum or conexus.
- Changing `SC-HISTORY-RETAINED`, its sentinel or its fix text.
- Inventing a successor for a superseded ADR that has none, or rewriting a
  successor's existing text.
- Segmenting a `claude` review, a round-2 delta, or a single commit; raising
  the review bound.
- Omitting files under a directory entry of a derived-path declaration, or
  reading derived declarations from the candidate.
- Holding a Task that edits only an owned skill's mirror, or an upstream
  skill, to the version record.
- Changing any other judge call, the Judge Key Variables or how a Run
  receives its environment.
- Reclaiming item binary copies, or changing how the item binary is built.

## Success Metrics

1. Success Metric: on a temporary repository holding a superseded ADR whose
   accepted successor does not name it, a superseded ADR with no accepted
   successor, a retired ADR, and a Live Document and a user guide page each
   linking a missing file, the built binary's history plan removes the first
   and third ADRs, lists the lineage line and both link rewrites, refuses the
   second by its path; after applying, `roundfix spec check --strict` reports
   `SC-HISTORY-RETAINED` only for the kept ADR, and a second plan reports
   nothing pending but that refusal.
2. Success Metric: a review fixture whose candidate deletes history files and
   rewrites a declared generated file above the bound lists both omissions and
   is reviewed; one whose remaining diff spans commits above the bound is
   reviewed in segments with one merged record, and a single commit above the
   bound, five segments and an oversized round-2 delta each block with their
   reason.
3. Success Metric: the review prompt's C2 text tells the reviewer that a QA
   Task carries no settlement sections, and a test fails when it does not.
4. Success Metric: a plan with a pending Task that edits an owned skill
   without its record yields `SC-SKILL-VERSION-UNTASKED` and fails `--strict`;
   the same plan with the declarations and step yields none, and every earlier
   corpus count at the Delivery Base is unchanged.
5. Success Metric: `roundfix archive <slug> --plan --offline` prints
   `judge: skipped (offline)` and makes no judge request even with a judge key
   in its environment.
6. Success Metric: a delivery test that rebuilds the item binary in place
   while the previous copy runs shows that the executed path differs from the
   build output, that the running copy's bytes and file identity are
   unchanged, and that no temporary file remains.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-17:

- PR-Agent documentation, "Review"
  (https://docs.pr-agent.ai/tools/review/): "the diff is split into up to
  `max_number_of_calls` chunks, each chunk is reviewed on its own, and the
  answers are merged into a single review that says how many chunks it was
  built from", with `max_number_of_calls` defaulting to 3. A published review
  tool bounds the number of calls and merges their answers, which is the shape
  of the Review Segments.
- GitHub Docs, "Customizing how changed files appear on GitHub"
  (https://docs.github.com/en/repositories/working-with-files/managing-files/customizing-how-changed-files-appear-on-github):
  "To keep certain files from displaying in diffs by default, or counting
  toward the repository language, you can mark them with the
  linguist-generated attribute". The forge readers use hides declared
  generated files from the reviewed diff, as the review omission does.
- ADR Tools (https://sourceforge.net/projects/adr-tools.mirror/, a mirror of
  https://github.com/npryce/adr-tools): "Create a new ADR file that is flagged
  as superceding ADR 9, and changes the status of ADR 9 to indicate that it is
  superceded by the new ADR." The common ADR tool writes the supersession into
  the successor, which is the line the cleanup adds when it is missing.
- Apple Developer Documentation, "Updating Mac Software"
  (https://developer.apple.com/documentation/security/updating-mac-software):
  "write the updated code to a temporary file and replace the existing file
  with that temporary one", as quoted by ADR-0281, whose authors read it
  through two copies because the page renders only with JavaScript.

## Validation

The decisions were checked on 2026-10-17 against two sources.

- Secondbrain. `wiki/index.md`, then two `qmd query` runs from the
  Secondbrain root. "superseded ADR successor names it baseline update history
  refused adopter" returned the mirrors of three Baseline and retirement ADRs and
  the `baseline` user guide under `projects/roundfix/mirror/`; they agree with
  the repository and hold no decision on a successor without its lineage
  line. "pre-PR review diff too large segment bound codex" returned the mirror
  of the ADR that makes the pre-PR review an explicit provider policy, which
  keeps an enabled review mandatory, the reason a large candidate must be
  reviewed rather than skipped. The pending inbox
  `inbox/roundfix/` holds no entry; the triaged entries
  `inbox/roundfix/_triaged/2026-10-14-baseline-text-keeps-a-history-of-retired-adrs-and-names-two-capture-doors.md`
  and
  `inbox/roundfix/_triaged/2026-10-09-o-baseline-update-sanitiza-o-historico-sem-relatar-os-links-que-quebra.md`
  were resolved by earlier Specs and share this
  context without adding to it.
- Exa. Supporting: the PR-Agent, GitHub Docs and ADR Tools passages above.
  Challenging: `codejury` (https://github.com/aiseclabs/codejury) and Qwen
  Code's pull request 6612 (https://github.com/QwenLM/qwen-code/pull/6612)
  split a large diff by file and require every chunk to tile the diff exactly.
  This Spec splits by commit range instead, so a finding's anchor names a line
  of a tree the author committed, and coverage is the whole commit range
  rather than an exact tiling of the final diff; ADR-0282 records the trade.
  These are standard references, not new research, so no Inbox Entry was
  captured.

## Glossary

- adds: **Review Segment**

## Decisions

- 2026-10-16: which failure rules govern cycle AH — "Regras de AE–AG
  (Recomendado)".
- 2026-10-16: where the leftover sources of cycle AH go — "Spec E no fim
  (Recomendado)".
- 2026-10-17: what to do with Spec E — "Escrever e entregar (Recomendado)".
- 2026-10-17: whether to edit the adopters — "Não faça a edição, só envie por
  inbox o que temos a trabalhar em cada um dos repos para que no melhor
  momento seja atualizado. Alguns estão com runs ativas e abrir PR enquanto
  isso pode ser um problema."
- 2026-10-17: whether a candidate still above the review bound is reviewed in
  bounded segments — "Sim, como na 0272 (Recomendado)".
- 2026-10-17: whether a failed queue item is retried — "Sim, uma vez por item
  (Recomendado)".
- 2026-10-16: which authority covers the Governed Paths — "Governed Paths
  declarados".

Autonomous route: no person present; decisions taken from the lead's briefing of 2026-10-17 and the six adopted Backlog Entries.

Taken from those sources, with no further interview:

- The cleanup writes the successor's lineage line and refuses only an ADR
  without an accepted successor. See ADR-0282.
- Live Documents and user guide pages get their Git Locators in a unit of
  their own; the history gate is unchanged. See ADR-0282.
- Derived declarations are read from the base commit, and directory entries
  are not used for omission. See ADR-0282.
- Segments apply to `codex` round 1 only, at most four. See ADR-0282.
- The item binary copy lives under the item's artifact directory. See
  ADR-0282.

## Open Questions

None.
