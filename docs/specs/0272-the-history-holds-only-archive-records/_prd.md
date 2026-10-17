---
spec: 0272-the-history-holds-only-archive-records
status: active
created: 2026-10-16
surfaces: [backend, cli, docs]
---

# The history holds only Archive Records

The History Root holds every delivered Spec's Archive Record, which the
Delivery Queue, Merge Evidence and the suite guard read, and beside them the
Reduced History Entries of 98 Findings and 32 Backlog Entries that no code
reads. A Jev pass on 2026-10-14 judged none of those 383 files reusable
knowledge; the seven lessons that no ADR states were captured in the
Secondbrain on 2026-10-15 before anything was removed. The maintainer decided
that retired material is deleted, not kept as a reference: "Não quero nada no
histórico que não seja relevante para o secondbrain" (2026-10-08) and "backlogs
e findings resolvidos e adrs deprecadas não acrescentam nada" (2026-10-15).
After this Spec a retirement deletes the file, every link to it names the file
in Git, existing history in this repository and in each adopter leaves through
the product's own commands, and the Baseline states the rule once, briefly,
with force only where the product enforces it.

## Prerequisites

The glossary and ADR congruence change of 2026-10-14 (pull request #545,
branch `docs/glossary-congruence`) merges first: this Spec quotes ADRs and
glossary terms in the state it leaves, and is authored on a branch stacked on
it. No Verification here depends on another unmerged Spec.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  Archive Record field is added; the record schema stays
  `roundfix/archive-record/v1` and a Git Locator reuses the path and commit
  already in use. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; reachability of a commit is proved against local refs, and
  tests run on temporary repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0278 (this Spec) decides the
  deletion: "A retirement deletes, and the History Root holds only Archive
  Records". It narrows ADR-0258, whose title says "a retirement writes reduced
  history"; ADR-0120, whose root holds "`specs/`, `findings/`, `adr/`,
  `backlog/`, and `reviews/`"; ADR-0248, whose sanitize stays a "Dry run by
  default."; ADR-0254, which keeps "One plan, one digest."; ADR-0260, whose
  rule "The archive relinks." extends to Git Locators; ADR-0261, where
  "records already archived stay byte-identical" now excepts only link
  destinations to deleted files; ADR-0277 for the destination of an
  unpromoted file; and ADR-0092, where "promotion moves the entry into the
  adopting Spec's `references/`" is unchanged. ADR-0247 still decides "An
  archive leaves an Archive Record, and the Spec folder stays in Git", and
  ADR-0232 still finds a delivery where "the default-branch head holds its
  archived `_prd.md`". ADR-0183: "A claim that attributes behavior to a
  decision record now carries a Claim Receipt". ADR-0186: "A Baseline clause
  now states only behavior the shipped product has." ADR-0202: "A preference
  is written as an obligation about its exception." ADR-0257: "the Baseline
  states each rule once". ADR-0244: "A Spec declares the domain terms it
  introduces". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0271: "A test bounds by real time only what it asserts".
  ADR-0095: "The secondbrain inbox is the single capture door". ADR-0203:
  "A repository rule becomes a Baseline clause only when it recurs"; the
  retirement clause is product behavior every adopter runs, not a
  repository rule. ADR-0184 asks a TechSpec for its Surface Transcripts,
  which the TechSpec gives for the History Sanitize Command. ADR-0251 reads a
  Legacy Archive Folder leniently and still applies to the folders the
  sanitize converts. ADR-0262 and ADR-0269 decide the Cancel Command's record
  and the archive's refusals, which stay as they are. ADR-0165, ADR-0169,
  ADR-0196, ADR-0233, ADR-0237, ADR-0272 and ADR-0279 decide review, delivery
  and queue rules this Spec does not change. ADR-0250: "A module version is
  chosen when it is recorded", which the module record step follows. ADR-0280:
  "A guide Task records its own fingerprint", which every Task here that
  edits a user guide page does. ADR-0281 does not apply: no test here rebuilds
  an executable in place. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on
  2026-10-14 the maintainer asked for a Spec "para ajustar o baseline de
  roundfix e skills", and on 2026-10-16, answering the plan that named the
  module edits, the version records and `make baseline-digests` "com a sua
  autorização expressa", wrote "Ok. vamos seguir o caminho da SPEC" and asked
  for "um baseline enxuto, claro, rígido onde deve ser e aberto onde pode ser";
  the 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
  se necessário") covers the skill references. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0272-the-history-holds-only-archive-records/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive.go`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/archive-spec/SKILL.md`, `skills/baseline_skill_contract_test.go`,
  `skills/roundfix/SKILL.md`, and, by the amendments of 2026-10-17,
  `docs/references/coverage-record.json`.

## Goals

- A retirement leaves no copy in the tree: the History Root holds only Archive
  Records, and the check refuses anything else under it.
- Every link from a Live Document or an Archive Record to a deleted file names
  that file in Git at a commit that stays reachable.
- Existing history leaves this repository and each adopter through the
  product's own commands, in a plan a person reviews before anything is
  deleted.
- The Baseline states the retirement rule once, each clause one obligation,
  with mechanical force where the product enforces it and room for the
  adopter where nothing is at stake.

## User Stories

1. As a maintainer, I want a closed Finding or Backlog Entry to disappear from
   the tree in the change that closes it, so that no Agent reads retired
   material as current.
2. As a reader of a user guide page or an Archive Record, I want a link to a
   retired file to open that file as it was, so that provenance survives the
   deletion.
3. As an adopter, I want `roundfix baseline update` to show me every file it
   will delete and every citation it breaks before I approve, so that I review
   one pull request instead of cleaning history by hand.
4. As an Agent closing a Finding, I want one short rule that tells me whether
   to capture its lesson, so that a lesson is not lost with the file and the
   Secondbrain is not flooded with duplicates.
5. As an Agent reading the Baseline, I want each clause to state one obligation
   and name the command that enforces it, so that I follow fewer, clearer
   rules.

## Core Features

1. **Deletion on retirement.** Recording a terminal status for a Finding,
   Rollup or Backlog Entry, `deferred` included, removes the file in the same
   change; handoffs and Review Artifacts are removed the same way; a rejected
   or deprecated ADR, and a superseded ADR whose accepted successor names it
   with its date, are removed. No Reduced History Entry is written. A Spec that
   absorbs a source adopts it, so its Archive Record names it; any other
   retirement names its reason in the commit that deletes the file.
2. **The History Root holds only Archive Records.** Once a repository's
   Baseline carries the deletion rule, the Spec Consistency Check reports an
   error for any file under the Findings, Backlog or ADR families of the
   History Root, for a terminal Finding left among unresolved Findings, and
   for a relative link from an Archive Record or a Live Document to a missing
   file or from a record to itself, and its diagnostic prints the Git Locator
   to use. Before that update the binary refuses none of them. Claim Receipts
   resolve an ADR in the active ADR directory and otherwise in Git; a Rollup
   member resolves as an unresolved Finding or as a file Git holds.
3. **Git Locator.** A Live Document or an Archive Record that cites a file the
   change deletes is rewritten in the same change: a link becomes its text
   followed by the Git Locator `<path>@<commit>` written as code, and a prose
   path becomes the Git Locator. The commit is the provenance commit the file
   already names when the default branch or a tag reaches it, and otherwise
   the newest commit the default branch reaches that holds the file; a file no
   reachable commit holds is named by its path and the record's pull request.
   Only the citation changes, and the rewrite lists every touched file. The
   Archive Command applies the same rule to the outcome it copies from a
   removed Spec folder and to each promoted copy: a destination inside the
   folder that is not promoted becomes a Git Locator, never a link to the
   record itself, whether the Archive Command or a Legacy Archive Folder
   conversion wrote it.
4. **History Sanitize and Baseline update remove the old families.** The
   History Sanitize Command and the history section of `roundfix baseline
   update` plan the removal of every Reduced History Entry, every retired ADR
   and every `superseded` Archive Record that nothing cites and that holds no
   regeneration grant, rewrite each citation the removal breaks to its Git
   Locator, give every Archive Record whose relative link points at a missing
   file its Git Locator, keep the dry run default, and refuse to apply without
   the History Full Tag. A removed path must be held by the tag or by a commit
   the default branch reaches. Retired ADRs and terminal Backlog Entries that
   the Baseline used to relocate into the old families are removed in the
   same plan, so a second run reports the history current.
5. **This repository's history.** The delivery removes the 98 reduced
   Findings, the 32 reduced Backlog Entries and the ten `superseded` Archive
   Records of Specs 0120 to 0129 through the built History Sanitize Command,
   rewrites the five user guide and reference links that named them, and
   gives the 40 broken relative links already in 32 Archive Records their
   Git Locators.
6. **A lean Baseline.** The docs-layout rule of the context-workflow module
   states retirement once: one clause for what each directory holds, one for
   deletion, one for provenance, one for using the runtime's commands. The
   clauses that described reduced entries, history families and "retired ADRs
   stay whole" are replaced; the reading hint that an empty findings
   directory is health is removed; `docs/_inbox/` is listed only for a
   repository without a Secondbrain, and handoffs are no longer called a
   capture door. The Secondbrain module gains one clause: the closing change
   states whether a Finding teaches a lesson that its ADR, Archive Record or
   reference does not already state, and only then captures it as one pending
   Inbox Entry with its source pinned. Each changed module carries a new
   version, and the derived files are regenerated.
7. **Documentation.** The glossary, the history, archive and baseline pages of
   the user guide, the context-driven development guide and the Roundfix
   skill's archive and baseline references describe deletion and Git
   Locators; the narrowed ADRs carry a dated note naming ADR-0278.

## Non-Goals / Out of Scope

- Removing Archive Records of delivered or cancelled Specs, or changing what
  the Delivery Queue, Merge Evidence or the suite guard read.
- Changing the Archive Record schema, the Plan Digest format or the Backlog
  and Findings Operational Contracts beyond the retirement paragraph.
- The Secondbrain side of retirement: Triage deleting an entry instead of
  moving it to `_triaged/`, the lint pinning a vanished mirror source, and
  rewriting the wiki's decision lifecycle page. That is a separate Spec in the
  Secondbrain.
- An authoring gate that records the grilling round in a PRD and checks the
  glossary against the code. That is a separate Spec.
- Opening the adopters' pull requests; each is an operator action after the
  release.
- A new enforcement level in the Baseline schema. Openness is expressed with
  the existing pattern of ADR-0202 and with recorded adopter decisions.

## Success Metrics

1. Success Metric: on the delivered tree the History Root holds exactly 243
   files, all Archive Records, and none under the Findings, Backlog or ADR
   families.
2. Success Metric: a planted file under any of those families, a planted
   Reduced History Entry and a planted relative link to a missing file each
   make the Spec Consistency Check fail, and the unplanted tree passes.
3. Success Metric: every Git Locator in a Live Document or an Archive Record
   names a commit that holds the file and is reachable from the default branch
   or a tag, or is replaced by the path and the record's pull request when no
   reachable commit holds it; no Archive Record links to itself.
4. Success Metric: the rendered docs-layout managed region is at most 11,000
   bytes, down from 14,494, and no two clauses of the changed modules state
   the retirement rule.
5. Success Metric: `roundfix baseline update` run in a scratch clone of an
   adopter lists every removal and every broken citation before applying, and
   a second run reports the history current.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-16:

- quereus documentation conventions
  (https://github.com/gotchoices/quereus/blob/main/docs/doc-conventions.md):
  "History is deleted, not archived." The repository keeps Git and one written
  account per change, the same shape as Git plus Archive Records.
- RFC 2119, section 6 (https://www.rfc-editor.org/rfc/rfc2119.html):
  imperatives "MUST only be used where it is actually required for
  interoperation or to limit behavior which has potential for causing harm".
  This is the test the Baseline applies before giving a clause mandatory
  force.
- GitHub Docs, "Getting permanent links to files"
  (https://docs.github.com/en/repositories/working-with-files/using-files/getting-permanent-links-to-files):
  a commit identifier in the URL "will permanently link to the exact version
  of the file in that commit", which is the property a Git Locator relies on.

## Validation

The decisions were checked against three sources on 2026-10-16. Where they
challenged the proposal, the change is named.

- **Secondbrain.** `wiki/index.md`, then `qmd query "política de retirada:
  apagar findings e backlog resolvidos do history, git guarda o texto"`, which
  returned `wiki/concepts/politica-de-retirada-e-poda.md`,
  `wiki/concepts/ciclo-de-vida-de-decisoes-e-raiz-docs-history.md` and
  `wiki/sources/roundfix-politica-retirada-e-poda-2026-10-15.md`. They confirm
  that the seven lessons were ingested before deletion and that the brain's
  decision lifecycle page still describes the history as a permanent ledger,
  which is the Secondbrain Spec's concern. A second query, "baseline normative
  clause mandatory rigid versus open adopter decision roundfix", returned the
  mirrors of ADR-0186, ADR-0202, ADR-0203 and ADR-0257, the existing rules for
  force and wording that the Baseline change follows.
- **Exa.** Supporting: quereus, the "Delete over deprecate" post, an agent
  documentation skill, Procida's lifecycle, the AGENTS.md evaluation (arXiv
  2602.11988, context files lowered success and raised cost by over 20%), the
  instruction-density benchmark IFScale (arXiv 2507.11538), the configuration
  smells study (arXiv 2606.15828, Lint Leakage in 62% of files) and
  Anthropic's guidance that a rule that must always hold belongs in a
  deterministic check. Challenging: ContextQB ("the agent reads the document,
  not its git history"), the ADR practice of keeping superseded ADRs, OpenSpec
  keeping completed changes, and the vLLM report of permalinks "correct as
  history, wrong as documentation". They changed the proposal three ways: the
  Archive Record and the successor ADR keep the lineage in the tree, a Git
  Locator must name a reachable commit, and the lesson capture is conditional.
  The digest is the Secondbrain Inbox Entry
  `inbox/secondbrain/2026-10-16-validacao-externa-do-history-so-com-archive-records.md`.
- **Jev.** jev-1.13.0, two passes. Two negative controls scored low
  ("delete every Archive Record too": support 0.14, breaks a reader 0.96;
  "delete without tag, plan or pull request": 0.10 and 0.89), so the judge
  discriminates. The decisions scored between 0.50 and 0.60 after their
  safeguards were added: contestable, not rejected. It asked for review before
  deletion, a check that fails closed and a reachable commit, which Core
  Features 2 to 4 carry. On the drafted clauses it found every one bundling
  several obligations and the command-enforced ones mostly restating the tool,
  which is why Core Feature 6 splits them and names the command. It classed
  the "empty findings is health" hint as an adopter's choice, and it is
  removed.

## Skills

- unchanged: command reconcile — the Managed Refresh of task_04 rewrites one link to a removed Finding in the reconcile user guide to its Git Locator; no flag, output or exit code changes, so the Roundfix skill's reconcile reference keeps its text.
- unchanged: command stop — the same rewrite of one link in the stop user guide; the skill's stop reference keeps its text.
- unchanged: guide docs/user-guide/commands/detached-runs.md — the same rewrite of one link; the covering skill references describe no history path.
- unchanged: guide docs/user-guide/commands/reconcile.md — the same rewrite as for command reconcile.
- unchanged: guide docs/user-guide/commands/stop.md — the same rewrite as for command stop.

## Glossary

- adds: **Git Locator**
- changes: **History Root**
- changes: **History Sanitize Command**
- changes: **Reduced History Entry**
- changes: **Rollup**
- changes: **Backlog Entry**
- changes: **Legacy Link Rebase**
- not a term: **Deletion on retirement** — feature heading
- not a term: **The History Root holds only Archive Records** — feature heading
- not a term: **History Sanitize and Baseline update remove the old families** — feature heading
- not a term: **This repository's history** — feature heading
- not a term: **A lean Baseline** — feature heading
- not a term: **Secondbrain** — product name already in use
- not a term: **Exa** — research tool name
- not a term: **Jev** — judge model name
- not a term: **Validation** — section label

## Decisions

The grilling rounds with the maintainer, recorded as answered:

- 2026-10-14: retired material is deleted, not kept as a reference; a Spec is
  required for the Baseline and the skills.
- 2026-10-15: the history keeps only Archive Records; resolved Findings and
  Backlog Entries and deprecated ADRs add nothing. See ADR-0278.
- 2026-10-15: the ten `superseded` Archive Records of Specs 0120 to 0129 are
  deleted ("Apagar os 10").
- 2026-10-15: a Finding's lesson goes to the Secondbrain inbox when it closes
  ("Inbox do brain ao fechar").
- 2026-10-15: adopters receive the cleanup through `roundfix baseline update`
  ("baseline update apaga o legado").
- 2026-10-15: a Live Document citing a deleted file is rewritten to a Git
  Locator ("Reescrever para localizador Git").
- 2026-10-15: Graphify does not replace the history; it stays out.
- 2026-10-16: `deferred` entries are deleted like resolved ones ("Apagar como
  os resolvidos").
- 2026-10-16: adopters' `superseded` records that nothing cites and that hold
  no grant are removed too ("Incluir nos adotantes").
- 2026-10-16: three Backlog Entries are adopted because they share the
  removed-file citation: the Archive Command's outcome links, the Baseline text
  about retired ADRs and capture doors, and prose paths in promoted copies,
  which `roundfix spec judge` suggested (same Spec 0.75). The Git Locator is
  the one rule the last entry asked to decide.
- 2026-10-16: the Baseline is lean and clear, rigid where it must be and open
  where it can be; openness uses the pattern of ADR-0202.

## Open Questions

None.
