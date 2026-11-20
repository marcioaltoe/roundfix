---
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
prd: _prd.md
created: 2026-11-20
---

# The Secondbrain search and the duplicate check go through the Indexus, without qmd — Technical Spec

## Executive Summary

The change is text and its proof. Two clauses of the Baseline's Secondbrain module are reworded under their ids, the research step of the owned skill `write-idea` is reworded the same way, one new accepted ADR records the decision, and ADR-0330 gains its notice. No Go production code changes: the plan code compares a clause's identity and force and never its text, so the repository's own commands regenerate every derived file. Seven new tests in three files pin the wording, the adopter's plan, the decision record and the absence of the retired command.

The primary trade-off is how the absence of `qmd` is held. A test that walks the whole tree against a list of allowed files would match the PRD's sentence word for word, and it would fail at the release, whose changelog names the word, and at every later Archive Record or ADR that recalls it. The design sweeps the classes Roundfix ships instead: every Baseline asset, every guide the Baseline writes here and every owned skill. It gives up catching a new mention in a file Roundfix does not ship, and the QA gate covers that once with a whole-word search of the tree.

A second trade-off is the size of the guide. The search clause grows from 264 to 649 characters, because each added sentence states one fact an Agent needs to read the answer or to act without it (`references/2026-11-20-characterization.md`, section 7).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. Both clauses keep their ids, and the new decision record takes the next free number of the repository's sequence. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. The guide names two tools of an MCP server that the session's own configuration provides, and no code, test or Verification calls it. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands. ADR-0330 is the decision this design amends in part, through the notice of API Contract 6. ADR-0257, ADR-0186, ADR-0222, ADR-0250, ADR-0360, ADR-0438, ADR-0439, ADR-0402 and ADR-0371 are the ones it obeys or uses, and ADR-0372 and ADR-0400 are the earlier rewordings it follows. ADR-0203, ADR-0256, ADR-0208, ADR-0248, ADR-0258, ADR-0338, ADR-0342, ADR-0349, ADR-0364, ADR-0380, ADR-0383, ADR-0390, ADR-0396, ADR-0399, ADR-0401, ADR-0411, ADR-0412, ADR-0415, ADR-0416, ADR-0441 and ADR-0444 are listed only to record why they do not apply; ADR-0412 is added to the PRD's list at this stage because it cites ADR-0438. This stage adds no ADR: the one decision record of this Spec is written by the Task that changes the clauses (`_prd.md` → Decisions). No accepted ADR is superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's two decisions of 2026-11-20, quoted with their provenance in the authorization record, which grants `implement`, `commit`, `push` and `pull_request` and neither `merge` nor `release`. The design changes no Makefile, workflow, Go module, lint or formatter file. It changes one Baseline module, the Standard TypeScript profile's pin and its formatter golden of the Secondbrain guide, this repository's Secondbrain guide and Setup Manifest, and both copies of one owned skill. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd/_authorization.md`; bounded files: `.agents/skills/write-idea/SKILL.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`, `internal/baseline/assets/modules/secondbrain.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/write-idea/SKILL.md`.

## System Architecture

Every component extends something that exists. The new files are three test files and one ADR.

### Baseline text

`internal/baseline/assets/modules/secondbrain.json` holds both clauses. `clause.secondbrain.02-query-order`, in `rule.secondbrain.index-first`, takes the text of API Contract 1. The last sentence of `clause.secondbrain.research-to-the-knowledge-owner`, in `rule.secondbrain.research-digests`, takes the text of API Contract 2. Each keeps its id, its `mandatory` force and its `replaces` value. No rule version and no guide version changes, as in the earlier rewordings; the module's version line and `internal/baseline/module-versions.json` are written by the record command.

### Derived files

`make baseline-digests` rewrites the formatter golden of the Secondbrain guide, the Standard TypeScript profile's pin, the three catalog files and the four plan goldens under `internal/baseline/testdata/`. The Managed Refresh rewrites `docs/agents/secondbrain.md` and `docs/agents/setup-context.json`. The Source Baseline corpus, its manifest and its index do not move.

### Owned skill

The research step of `.agents/skills/write-idea/SKILL.md` takes the text of API Contract 3. `make skills-sync` copies it to `skills/write-idea/SKILL.md`, and the skill record command raises both version fields and records them. This edit moves no derived Baseline file.

### Decision record

One new accepted ADR, described in API Contract 5, and the notice of API Contract 6 in ADR-0330. The documentation part of `roundfix baseline update`, which exists, adds the new path to the Standard Export List of `.secondbrain-export` and the new identifier to ADR-0330's Amendment List in the same run as the Managed Refresh.

### Tests

`internal/baseline/secondbrain_indexus_clauses_test.go`, `skills/indexus_duplicate_check_test.go` and `internal/docscontract/indexus_decision_record_test.go`. Testing Approach names what each holds.

## Implementation Design

### Interfaces

No production interface is added or changed. The tests are the contract a Task implements:

```go
// internal/baseline/secondbrain_indexus_clauses_test.go, package baseline
func TestTheSecondbrainClausesNameTheIndexus(t *testing.T)
func TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten(t *testing.T)
func TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest(t *testing.T)
func TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand(t *testing.T)

// skills/indexus_duplicate_check_test.go, package skills
func TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus(t *testing.T)
func TestNoOwnedSkillNamesTheRetiredSearchCommand(t *testing.T)

// internal/docscontract/indexus_decision_record_test.go, build tag docscontract, verify:always
func TestTheIndexusDecisionAmendsTheConsultationRule(t *testing.T)
```

```text
1. Both clauses keep their id, their mandatory force and their replaces value.
2. No Baseline asset outside the Source Baseline, no guide under docs/agents and no owned skill holds the whole word qmd, in any letter case.
3. A phrase of API Contract 4 occurs exactly once in the text it belongs to, after whitespace is collapsed.
4. No test reaches the network, calls an MCP server, runs qmd, or reads a file under the Spec root or outside the repository.
5. The decision record is found by its title, never by its number.
```

### Data Models

No schema changes. Two documents change shape by one field each:

- The new ADR's header is the lifecycle overlay of the docs-layout guide with `status: accepted` and, on the next line, `knowledge: true`. Its file is `docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md` when the next free number is still 0449; the uncommitted plan of Spec 0415 holds 0448 (`_prd.md` → Open Questions, item 4). The Tasks author confirms the number on the tree the Tasks are authored on.
- ADR-0330's header gains the new identifier as the last entry of `amended_by`, written by the update, and its `updated_at` takes the new ADR's `created_at`, written by the Task.

### API Contracts

1. API Contract: the search clause. The `guidance` of `clause.secondbrain.02-query-order` becomes this one line, with the JSON escapes the file needs:

   ```text
   Read `wiki/index.md` first. Then search with the `pesquisar_acervos` tool of the `indexus` MCP server, which returns one group per acervo the session's key reaches; `abrir_conteudo` opens a returned page by its `caminho` as returned, which for a `wiki` result is relative to `wiki/`. An answer without a `wiki` result is a limitation to record. When the tool is not in the session or does not answer, the Secondbrain search is unavailable: substitute no other search tool. Inspect `projects/<project>/mirror/` only when the index and the search point there, and open only the files required for the task; treat mirrors as references, not workspaces.
   ```

   Failure modes it states: an acervo the key does not reach, an answer without a wiki page, and a tool that is absent or silent. `clause.secondbrain.05-research-limitations` is unchanged and says what an Agent then does: record the attempted consultation and its limitation, and continue.

2. API Contract: the duplicate check. In `clause.secondbrain.research-to-the-knowledge-owner`, the sentence that begins "Run the advisory" becomes the sentence below. Every sentence before it is unchanged.

   ```text
   Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`, verify that each returned page opens, and name a strong verified match in the message so the owner extends existing knowledge instead of duplicating it; a relevance label or a rank alone never decides.
   ```

3. API Contract: the research step of `write-idea`. The paragraph that begins "When the session performed substantive external research" keeps its first sentence, and the rest becomes the two sentences below, wrapped as the file wraps its lines.

   ```text
   Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`, verify that each returned page opens, and name a strong verified match in the message; a relevance label or a rank alone never decides. When that tool is not in the session or does not answer, say so in the message and substitute no other search tool.
   ```

4. API Contract: the pinned phrases, one per guarantee. A Task names one subtest per phrase.

   ```text
   search clause
     search with the `pesquisar_acervos` tool of the `indexus` MCP server
     one group per acervo the session's key reaches
     `abrir_conteudo` opens a returned page by its `caminho` as returned
     for a `wiki` result is relative to `wiki/`
     An answer without a `wiki` result is a limitation to record
     the Secondbrain search is unavailable: substitute no other search tool
     Read `wiki/index.md` first
     treat mirrors as references, not workspaces
   duplicate-check clause
     Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`
     verify that each returned page opens
     so the owner extends existing knowledge instead of duplicating it
     a relevance label or a rank alone never decides
   write-idea
     Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`
     verify that each returned page opens
     a relevance label or a rank alone never decides
     say so in the message and substitute no other search tool
   ```

5. API Contract: the decision record. Title: "The Secondbrain search and the duplicate check go through the Indexus, with no second path". It states, in this order:
   - the context: the maintainer's decision of 2026-11-20 in his words, and what the earlier command cost;
   - the decision: the search that follows the wiki index and the advisory duplicate check use the Indexus tool, and no clause, guide or owned skill names a second search tool;
   - what an Agent does without the tool, without its answer, or with an answer that holds no wiki page: it records the limitation, continues with the evidence it has and substitutes nothing;
   - that both clauses keep their identity and force, and that Roundfix's code never calls the Indexus and adds no Baseline decision and no Doctor check for it;
   - that the decision takes effect with the release that carries it, and that the project that owns the Secondbrain removes the earlier command from the machines after that release;
   - the options considered and rejected: the Indexus first with the earlier command as a reserve, a Baseline decision that says whether the Indexus is available, new clause ids, describing the capability without naming the tool, and naming one harness's joined tool identifier;
   - that it amends ADR-0330.

   It names no client and no acervo, and it cites no Spec folder path.

6. API Contract: the notice in ADR-0330. One paragraph directly after the first paragraph of `## Decision`, which holds the sentence that names the earlier command. It begins with the bold lead `Amended by`, the new ADR's identifier as a link to its file, and the date of the new ADR's `created_at` in parentheses, in the form the three notices below it use. It says that the consultation rule after the wiki index is a search with the Indexus tool, that the duplicate check uses the same tool, and that no second search tool is named. The sentence it follows is not edited.

7. API Contract: the plan an adopter sees. For a repository that holds the wording of the Delivery Base with the Secondbrain turned on, `roundfix baseline update` plans two file changes, the Secondbrain guide and the Setup Manifest. For an adopter whose Setup Manifest names the Source Baseline, the clause accounting is the one it already had: `clause.secondbrain.02-query-order` is `retained` with `mandatory` force, `clause.secondbrain.research-capture` is `replaced` with the target `clause.secondbrain.research-to-the-knowledge-owner` and `mandatory` force, and no clause is `unaccounted`. A second update after the apply plans no change.

### Surface Transcripts

None. No command, flag, exit code, configuration key or output line changes. `roundfix baseline update` prints digests of the tree it runs on and the paths of API Contract 7, which the adopter test asserts through the plan and the QA gate reproduces through the built binary.

## Coverage Map

- Goal 1 → Baseline text (API Contract 1)
- Goal 2 → Baseline text, Owned skill, the two sweep tests (invariant 2)
- Goal 3 → Baseline text, Derived files (API Contract 7, the adopter test)
- Goal 4 → Baseline text (API Contract 1, the sentences about an answer without a wiki page and about an absent tool, with the unchanged limitation clause)
- Goal 5 → Build Order, Risks & Considerations (the authorization record grants no `merge`)
- Story 1 → Baseline text (API Contract 1)
- Story 2 → Baseline text, Owned skill (API Contracts 2 and 3)
- Story 3 → Baseline text, Owned skill (API Contracts 1 and 3)
- Story 4 → Derived files (API Contract 7)
- Story 5 → Build Order (the delivery stops at an open Pull Request)
- Story 6 → Integration Points (the operator's notice; no Task sends it)
- Core Feature 1 → Baseline text (API Contract 1)
- Core Feature 2 → Baseline text (API Contract 2)
- Core Feature 3 → Baseline text, the two sweep tests (API Contract 1, invariant 2)
- Core Feature 4 → Baseline text, the clause test, the adopter test (invariant 1, API Contract 7)
- Core Feature 5 → Owned skill (API Contract 3)
- Core Feature 6 → Derived files, Build Order 1 and 2
- Core Feature 7 → Decision record (API Contracts 5 and 6)
- Core Feature 8 → the two sweep tests, the QA gate's search of the tree
- Success Metric 1 → `TestTheSecondbrainClausesNameTheIndexus`, `TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten`
- Success Metric 2 → `TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest`, the QA gate through the built binary
- Success Metric 3 → `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus`, with the existing `TestAuthorialSkillSync` and `TestEveryOwnedSkillVersionIsRecorded`
- Success Metric 4 → API Contracts 1 to 3 against `references/2026-11-20-characterization.md`, sections 5 and 7, and `references/2026-11-20-messages-and-decisions.md`, section 9
- Success Metric 5 → `TestTheIndexusDecisionAmendsTheConsultationRule`, the QA gate's run of the update
- Success Metric 6 → `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand`, `TestNoOwnedSkillNamesTheRetiredSearchCommand`, the QA gate's search of the tree

Deviations from the PRD:

- Core Feature 1 gains one sentence the PRD does not list, that an answer without a wiki page is a limitation to record. The Indexus project confirmed after the PRD was written that a project session's key does not reach the acervo that holds the wiki, and the session that started this phase asked for the guide to make such a session record it (`references/2026-11-20-messages-and-decisions.md`, section 9). The PRD's decision that the guide adds no duty to list acervos stands: the sentence names no acervo.
- Core Feature 2 says the Agent verifies that the returned paths exist. The clause says "verify that each returned page opens". It is the same guarantee, stated by the one method that needs no knowledge of what a path is relative to.
- Core Feature 8 lists where the word still appears. It also appears in the new decision record and in ADR-0330's notice when they name what was replaced, in the three new test files, which hold it as the forbidden word and as the earlier sentence, and in this Spec's folder and later its Archive Record.
- Success Metric 6 asks for one named test that holds the list. Two tests hold it, one in the package that embeds the Baseline assets and one in the package that embeds the owned skills, and each sweeps a class instead of the tree, for the reason the Executive Summary gives.

## Absence Evidence

1. The Indexus wording of the two clauses, the guide and the skill step — command: `rg -n -i 'pesquisar_acervos|abrir_conteudo|indexus' internal/baseline/assets docs/agents .agents/skills/write-idea skills/write-idea`; result: no match.
2. The seven tests and their three files — command: `rg -n -i 'indexus|pesquisar_acervos' -g '*_test.go' internal skills cmd`; result: no match.
3. A check that forbids the retired command — command: `rg -n -i -w 'qmd' -g '*.go' internal skills cmd`; result: no match, so no Go file names the word and no test sweeps for it.
4. The decision record — command: `rg -l -i 'indexus' docs/adr`; result: four files, the ADRs numbered 0401, 0402, 0406 and 0416, each naming the Indexus as an adopter of Roundfix and none deciding how the Secondbrain is searched.
5. The notice in ADR-0330 — command: `rg -n 'Amended by' docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`; result: three matches, the notices of ADR-0383, ADR-0390 and ADR-0396.

## ADR Alternatives

- ADR-0330: amends it in part and takes none of the options it rejected. It rejected a second door for one job, a Baseline decision for a service whose reach is a property of the session, and a call from Roundfix's code to that service. The design keeps one search path, adds no decision, and adds no call.
- ADR-0257: obeys it. It rejected renaming every reworded clause; both clauses keep their ids.
- ADR-0186: obeys it. It records no rejected option. The new words state what the named tools do and cite no Spec and no ADR, and each clause keeps its enforcement level.
- ADR-0222: obeys its title rule; the guide reaches only repositories that turned the Secondbrain on. Its two rejected options concern skill triggers, which this design does not touch, and its rename rule is the one ADR-0257 narrowed.
- ADR-0250: uses it. It rejected a record step inside `make baseline-digests`; the module record command runs first, as its own command.
- ADR-0360: obeys it. It rejected a bump that moves only part of the Vendored Skills; this design runs no bump and edits no Vendored Skill.
- ADR-0372: follows its precedent of a clause reworded under its id. Its rejected options concern profile decisions and guide warnings.
- ADR-0400: follows its precedent of a clause reworded under its id with its Source Baseline row. Its rejected options concern how a Spec is sized.
- ADR-0402: obeys it. It rejected a memory server or a search provider built into the Research Subagent; the Indexus search stays with the session that runs the phase.
- ADR-0438: uses it. It rejected recording the list as a managed artifact of the Setup Manifest; the update adds the new ADR's path through its documentation part.
- ADR-0439: uses it. It rejected an `amended` status and a notice kept in the body alone; ADR-0330 stays accepted and gains both the notice and the header entry.
- ADR-0371: uses it. It rejected a session hook and a full release note at session start; adopters learn of the release from the Update Check as it is.
- ADR-0203: not applicable — no hand-written rule of a repository is promoted; two shipped clauses are reworded.
- ADR-0256: not applicable — no command help, configuration key, exit code or user-guide page changes, so no Behavior Surface moves.
- ADR-0208: not applicable — it governs how sources are grouped into a Spec, which the PRD settled.
- ADR-0248: not applicable — no history is sanitized or converted.
- ADR-0258: not applicable — no retirement is written; no file is deleted.
- ADR-0338: not applicable — the follow-up mode and its handoff are not changed.
- ADR-0342: not applicable — the settlement Verification is used as it is.
- ADR-0349: not applicable — the authorization record grants no `release` and carries no Release Ceiling.
- ADR-0364: not applicable — no required skill is added or checked.
- ADR-0380: not applicable — no skill is replaced or retired.
- ADR-0383: not applicable — the message rules it added to ADR-0330 are not changed.
- ADR-0390: not applicable — the message rules it added to ADR-0330 are not changed.
- ADR-0396: not applicable — the completion notice it defines is sent by the operator, outside every Task.
- ADR-0399: not applicable — no nested carrier is inventoried differently.
- ADR-0401: not applicable — no Authoring Model Declaration changes.
- ADR-0411: not applicable — no rule of this repository moves between guides.
- ADR-0412: not applicable — the Release Plan's classification of a commit is not changed, and this Spec is granted no release.
- ADR-0415: not applicable — the Doctor Command gains no check.
- ADR-0416: not applicable — the Update Check is not changed.
- ADR-0441: not applicable — the Research Subagent's tool list is not changed.
- ADR-0444: not applicable — no harness permission file is written.

## Integration Points

- The Indexus MCP server, in another project: named and never called. The boundary is the guide's text. The names and the shape of the answer are the ones of its contract `indexus-acervos/1`, and its owner tells this project of a change by message.
- The skills catalog, which vendors `write-idea`: it receives the new step by its next sync. No Catalog Skill Bump runs here.
- The Secondbrain's sync: it reads `.secondbrain-export` after a push to the default branch. This repository's own lines exclude every ADR, so the new line of the Standard Export List publishes nothing from here.
- ai-memory: no component calls it. The operator's session sends the line the Secondbrain's project asked for, after the release.
- The Update Check: unchanged. It tells an adopter of the release and of the pending Baseline plan.

## Testing Approach

No new seam. Every test reads the embedded catalog, the embedded skills or files of this repository, and plans on a temporary repository.

- `TestTheSecondbrainClausesNameTheIndexus`: finds both clauses with `embeddedBaselineClauses`, and asserts invariant 1, each phrase of API Contract 4 by invariant 3, and that neither holds the retired word.
- `TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten`: asserts the same phrases in three texts: the postimage of `docs/agents/secondbrain.md` in a plan of the Standard TypeScript profile with `secondbrain.enabled` true, the formatter golden read with `catalog.Asset`, and this repository's `docs/agents/secondbrain.md`.
- `TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest`: clones the catalog, gives both clauses the sentences of the Delivery Base with `mutateCatalogClause`, adopts the Standard TypeScript profile with the Secondbrain turned on, commits, and plans a Managed Refresh with the embedded catalog. It asserts API Contract 7, applies the plan and asserts that a second plan changes nothing. On the unchanged tree the two catalogs are equal and the plan changes no file, so the test fails.
- `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand`: walks `assets/` and `../../docs/agents`, as `TestAuthorizationClausesDropTheRemovedPhrases` walks them, for the whole word in any case. It holds one exempt path, the Source Baseline corpus copy of the Secondbrain guide, with its reason, fails when that path no longer needs the exemption, and fails when it scanned no module or no golden.
- `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus`: asserts the phrases of API Contract 4 in the canonical copy and in the mirror with `assertIdeaStageSkillPhrases`.
- `TestNoOwnedSkillNamesTheRetiredSearchCommand`: sweeps every file `skills.Files()` returns. `TestAuthorialSkillSync` already holds the canonical copies equal to them.
- `TestTheIndexusDecisionAmendsTheConsultationRule`: finds the ADR by its title and asserts, with `spec.ReadADRHeader`, that it is accepted and its Knowledge Mark is true; that ADR-0330's Amendment List holds its identifier; that the notice of API Contract 6 comes after the sentence that names the earlier command and carries the date of the new ADR's `created_at`; and that the Standard Export List holds the ADR's path once, between the two markers, with its ADR lines in ascending order.
- Existing tests that keep passing unchanged: `TestBaselineClauseForceIsCharacterized`, `TestEveryInboxClauseAnAdopterHoldsIsReplaced`, `TestFormatterComposition`, `TestTheSetupManifestIsCurrent`, `TestEveryBaselineModuleVersionIsRecorded`. No existing test pins the words this Spec changes, so none is edited.
- The QA gate, through the built binary and without the network: the plan of API Contract 7 on a temporary repository adopted with a binary built from the Delivery Base; `roundfix baseline update --no-skills --no-history` on the candidate, which plans no change; a whole-word search of the tree for the retired word, compared with Core Feature 8 and the deviation above; and the guide's words against the outside evidence of Success Metric 4.

On the prototype tree, with the wording of API Contracts 1 to 3, `go test ./internal/baseline ./skills` passed and `make verify-docs` exited 0, so the coverage record, the skill coverage record and the Behavior Surface Record do not move.

## Build Order

1. The guide and its decision record: API Contracts 1, 2, 5, 6 and 7, with the four tests of `internal/baseline` and the test of `internal/docscontract`. The order that worked on the prototype: edit the module; run the module record command; run `make baseline-digests`; write the new ADR and the notice with ADR-0330's `updated_at`; run the Managed Refresh once, which writes the guide, the Setup Manifest, the Standard Export List and ADR-0330's Amendment List. The clause edit and the decision record are one step because the PRD requires that the record and the guide never disagree.
2. The `write-idea` step: API Contract 3, with the two tests of `skills`. Edit the canonical copy; run `make skills-sync`; run the skill record command (depends on: none; steps 1 and 2 share no file, and the skill edit moves no derived Baseline file, measured on a second scratch clone).
3. The QA gate (depends on: 1 and 2).

Step 2 is five lines of one skill, three files and one test file. Under the sizing rule of the Spec routing guide it may join step 1 in one Task. A joined Task runs the order of the PRD's prototype: both edits, `make skills-sync`, the skill record command, the module record command, `make baseline-digests`, the ADR and the notice, then the Managed Refresh.

Files shared with other Specs. Every Baseline module or owned-skill change that merges while this Spec's Pull Request waits moves the module version record, the catalog files, the plan goldens, the Standard TypeScript profile's pin, this repository's Setup Manifest and the owned-skill version record. Each step regenerates them on the tree it runs on and names the next version, never a number. No Verification of this Spec depends on another Spec, and the Task Graph needs no `requires`.

## Risks & Considerations

- A condition of the merge: the wiki must be reachable through the Indexus from a project session. The Indexus project confirmed that team keys do not reach the acervo that holds the wiki, counts it as a loss in its test, and recommends that this change not be published without a second path before the maintainer decides between widening those keys and moving the non-personal part of the wiki. The maintainer excluded a second path, and the authorization record already holds the merge. Until then a project session finds mirrors and raw sources, reads the wiki index, and records the limitation the search clause names.
- The sentence about a wiki result's `caminho` rests on the Indexus's answer and on the Secondbrain's request. The contract's one example of `abrir_conteudo` still shows a path that begins with `wiki/`, and its owner said it will be corrected. Before the merge the operator reads the contract's then-current revision. If the wiki moves to another acervo with another path form, API Contract 1 changes before the merge.
- The Legal acervo leaves the reach of every project key, and legal questions go to the Fiscus. Neither the guide nor `write-idea` names law, a legal collection or the Fiscus today, so this design adds no word about it, and the reworded clauses tell no session to search law.
- The Secondbrain guide grows by 413 characters, about seven percent. No ceiling binds this guide.
- The ADR number can be taken while the Pull Request waits. The operator then renames the file and its link in the notice, and runs the update again; invariant 5 keeps the test independent of the number.
- The phrase tests pin one phrase per guarantee and not the whole text, so a rewording that keeps every phrase passes. The formatter golden holds the whole text.
- A client that holds a stale tool schema refuses the Indexus's answer. That is the case "does not answer" of API Contract 1.
- No test of this repository can see the Indexus rename a tool or change its answer. The comparison of Success Metric 4 is the only one, and it reads a record of 2026-11-20.

## Glossary

None.

## Decisions

- The guide names the server and each tool separately and no harness's joined identifier, because that identifier differs by harness (`references/2026-11-20-characterization.md`, section 7).
- The search clause states the rule for an absent or silent tool itself and leaves what to record to the unchanged limitation clause, so one clause says when the search is unavailable and one says what an Agent then does.
- The skill step carries its own sentence for an absent or silent tool, because a session can read the skill without the guide, and a Vendored Skill may still show the earlier command until its next bump.
- The absence of the retired word is held by class, in two tests, and by the QA gate's search of the tree once.
- The adopter test builds the earlier wording by mutating a cloned catalog, so no fixture file of the earlier guide is added.
- The wording follows the Indexus's answer of 2026-11-20 on the path of a wiki result, which settles the PRD's third open question, and adds nothing about law.
- No Open Authoring Question was left by this phase, and no Delegated Decision was taken.
