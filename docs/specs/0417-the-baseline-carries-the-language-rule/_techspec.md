---
spec: 0417-the-baseline-carries-the-language-rule
prd: _prd.md
created: 2026-11-27
---

# The Baseline carries the language rule for code, data and products — Technical Spec

## Executive Summary

The Language Rule enters the Baseline as one new rule of the `core` module with four mandatory clauses, the Reader Language as a required string decision rendered in the opening lines of the agent-instructions guide, and the Name in Code as two sentences inserted into the canonical-terms clause, which keeps its identity. The Code Language Exception is an optional decision that holds a language and a source as copied text. While it is recorded, a clause gate takes three of the four clauses out of the guide and puts one clause in. The Pre-PR Review prompt and the QA gate skill each gain one instruction, and no checker is added. The primary trade-off: the design accepts a required decision that every adopter answers once and that reaches the answer fixtures of ten test files, so that no guide states a language its repository did not give. It also places the new clauses in a rule that no profile lists in `requiredRules`, so that no Source Baseline row is written by hand. A prototype on a scratch clone of `8ac0224b` carried the design through the catalog, the regeneration commands, the Managed Refresh and the tests: it touches the 28 Governed Paths the PRD expected and none of the eleven it held in reserve. Section 8 of `references/2026-11-27-characterization.md` records what it ran and found.

## Project Constraints

- Identifier strategy: not applicable — no Internal Identifier is minted. The rule `rule.core.code-and-reader-language`, its five clause ids, the decision ids `language.reader` and `language.code_exception` and the decision type `code-language-exception` take the forms the catalog already uses, `clause.domain.canonical-language` keeps its id, and the ADR takes the next free number when the Task Graph is authored. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added, and no test or Task Verification reaches the network, another repository or ai-memory. `roundfix baseline update` keeps its flags, its output shape and its exit codes; a new decision id and a new refusal reason are the only changes of its surface. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands. ADR-0118, ADR-0186, ADR-0203, ADR-0222, ADR-0257, ADR-0263 and ADR-0372 shape the catalog change; ADR-0189, ADR-0233, ADR-0250, ADR-0256, ADR-0280 and ADR-0360 bind how modules, owned skills, guides and the glossary change; ADR-0196 binds the review prompt; ADR-0244, ADR-0438 and ADR-0439 bind the terms and the decision record; ADR-0371, ADR-0383, ADR-0431, ADR-0436 and ADR-0437 are obeyed without a component; ADR-0453 is listed to record that no gate time file changes. ADR-0005, ADR-0010, ADR-0139, ADR-0141, ADR-0150, ADR-0182, ADR-0228, ADR-0248, ADR-0252, ADR-0253, ADR-0258, ADR-0259, ADR-0272, ADR-0286, ADR-0317, ADR-0319, ADR-0330, ADR-0341, ADR-0342, ADR-0349, ADR-0359, ADR-0364, ADR-0380, ADR-0386, ADR-0390, ADR-0393, ADR-0395, ADR-0396, ADR-0399, ADR-0411, ADR-0412, ADR-0416, ADR-0442, ADR-0444, ADR-0448 and ADR-0452 are listed only to record why they do not apply. This stage writes no ADR file: Core Feature 12 has one Task write the ADR together with the clauses it describes, and API Contract 8 fixes what it states. No accepted ADR is superseded or amended, and the design takes no alternative an accepted ADR rejected. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — protected tooling mutation is proposed, and express maintainer authorization for its paths is granted. The maintainer's decision of 2026-11-20 authorized the PRD, the TechSpec, the Tasks, the implementation and a Pull Request, and reserved the protected Baseline paths, the merge and the release to a decision of his that names the exact list. He granted that list on 2026-11-27 13:40 BRT, ai-memory message `01a21de1-c3f0-76d0-bc5e-784e4cbffed0`, which the Spec-contained authorization record named below records as `approved` with the grant date 2026-11-27, together with the operations `implement` to `release` and the Release Ceiling `minor`. The list was measured with `speccheck.GovernedPath` in a `go test -overlay` probe over every file the prototype changed and every file Build Order steps 1 to 3 name; it equals the list of the PRD, and the Tasks author measured again every file the Tasks declare and found no governed file outside it. The record also declares the seven regeneration commands of API Contract 7 with their outputs. No Makefile, workflow, Go module, lint, formatter or Project Config file changes. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`; Spec-contained authorization record: `docs/specs/0417-the-baseline-carries-the-language-rule/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/review.md`, `.agents/skills/setup-context-driven/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, `internal/baseline/assets/decisions.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/go-cli-tui.json`, `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`, `internal/baseline/assets/profiles/python.json`, `internal/baseline/assets/profiles/rust-cli.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/index.json`, `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go`, `internal/cli/baseline_plan_test.go`, `internal/cli/baseline_release_gate_test.go`, `internal/docscontract/publicdocs_test.go`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## System Architecture

Nothing new is built beside the Baseline: every part extends a mechanism the catalog, the review or a skill already has.

### Baseline text

The `core` module gains `rule.core.code-and-reader-language`, rendered in the agent-instructions guide after `rule.core.english-guidance`. A guide renders the clauses of a rule sorted by id, so the ids are chosen for the reading order: the code-names clauses first, the reader clause last. The rule is listed in no profile's `requiredRules`, as `rule.core.cross-project-continuity` is not, and it reaches every profile because every profile composes `core`. The `context-workflow` module rewords `clause.domain.canonical-language` by insertion: every word the Source Baseline recorded stays, in order, and two sentences are added between its two sentences.

### Reader Language decision

`language.reader` is a required decision of type `string` with the suggested value `English`. The `core` module requires it, the five built-in profiles list it after `language.generated`, and the agent-instructions template gains one line that the decision's token fills through the renderer every plain string decision uses. Nothing in the planner, the update or the interactive flow changes for it: the stop for a new required decision, `--decision`, `--adopt-suggested` and the profile decision gap already do what Core Features 2 and 10 ask.

### Code Language Exception

`language.code_exception` is an optional decision of a new type, `code-language-exception`, whose value is an object with a language and a source. It has no default and no suggestion: absence is its only unrecorded state. Two mechanisms carry it. Its token, on the Reader Language line of the template, renders nothing while the decision is absent and one sentence with the language and the source while it is recorded. A clause gate `appliesWhen` gains a second form, `present`, beside `equals`: three clauses apply only while the decision is absent and one only while it is recorded. The interactive `roundfix baseline` flow never prompts for it and carries a recorded value unchanged; like `auth.base_path`, it is declared deliberately, with `roundfix baseline update --decision`.

### Review and QA

`buildReviewPrompt` gains one instruction sentence, held in a constant, between the absence-claim instruction and the finding grammar. It is no Delivery Convention: the conventions list what a delivery writes by design, and this sentence asks for a finding. The `qa-gate` skill gains one section under `## 4. Exercise real flows`.

### This repository

This repository answers `language.reader=English` and declares no exception. Its two guides and its Setup Manifest change through the Managed Refresh, the derived Baseline files through `make baseline-digests`, and the module versions through the record command.

### Record of the decision

One ADR, five Glossary Definitions, the Baseline reference and the review reference of the Roundfix skill, the `setup-context-driven` skill, three user-guide pages and the published Decision Document.

## Implementation Design

### Interfaces

```go
// internal/baseline/code_language_exception.go (new)
const codeLanguageExceptionDecisionID = "language.code_exception"
const codeLanguageExceptionDecisionType = "code-language-exception"

type CodeLanguageException struct {
	Language string `json:"language"`
	Source   string `json:"source"`
}

// Accepts an object with exactly language and source; refusals are API Contract 9.
func normalizeCodeLanguageException(value any) (CodeLanguageException, error)
```

```go
// internal/baseline/plan.go (changed): a gate holds decision and exactly one of equals, present.
func clauseApplies(catalog *Catalog, clause document, recorded map[string]any) bool

// internal/cli/review.go (changed)
const reviewLanguageRuleInstruction = "..." // API Contract 13
```

```text
1. No guide states a Reader Language its repository did not answer or adopt by flag.
2. A recorded Code Language Exception always holds a non-empty language other than English and a non-empty source.
3. The guide of a repository holds clause.core.code-names-follow-the-existing-codebase or the three clauses gated on absence, never both and never neither.
4. Every word of clause.domain.canonical-language that the Source Baseline recorded stays in the clause, in order.
5. A gate with present never reads a decision's default: an unrecorded decision is absent.
6. No profile lists rule.core.code-and-reader-language in requiredRules, and no Source Baseline file changes.
7. The review prompt keeps its finding grammar, its Delivery Conventions and their version.
```

### Data Models

The two decisions, added to `internal/baseline/assets/decisions.json` (effects elided: each has one effect, `when {"present": true}`, with one `renderBindings` entry for `guide.agent-instructions`, `template.guide.agent-instructions` and the token named after the decision):

```json
{"id": "language.reader", "version": 1, "type": "string", "suggestion": "English",
 "summary": "The Reader Language: the language of what the product's users read."}
{"id": "language.code_exception", "version": 1, "optional": true, "type": "code-language-exception",
 "summary": "The Code Language Exception: the language an existing codebase keeps for code and infrastructure names, with the source of its owner's decision."}
```

A recorded exception in the Setup Manifest, whose decision record still carries only `value`:

```json
"language.code_exception": {"value": {"language": "Portuguese", "source": "docs/adr/0141-example.md"}}
```

The gate of a clause, in a module:

```json
"appliesWhen": {"decision": "language.code_exception", "present": false}
```

### API Contracts

1. API Contract: the Reader Language decision. `language.reader` as Data Models gives it, placed after `language.generated` in the catalog. `core` adds it to `requiredDecisions`, and each of the five built-in profiles lists it in `entryDecisions` directly after `language.generated`, which is also where the interactive flow asks it. Any non-empty string is a valid answer; the user guide and the skills give `English` and `pt-BR` as examples.
2. API Contract: the template. `internal/baseline/assets/templates/guides/agent-instructions.md` gains, between the incremental Verification line and the `{{branch.prefix}}` line, the line `The Reader Language is {{language.reader}}.{{language.code_exception}}`, and the template's entry in `templates/index.json` lists both tokens. The line sits before the branch-prefix sentence because `internal/baseline/branch_prefix_optional_test.go` holds that sentence to be the last of the opening lines. Build Order step 1 adds the line with the first token only; step 2 appends the second.
3. API Contract: the rule. `rule.core.code-and-reader-language`, version 1, coverage `coverage.language`, listed in the `rules` of `guide.agent-instructions` after `rule.core.english-guidance`, with four clauses, all `mandatory`. Each line below is a clause id, a tab and its guidance:

   ```text
   clause.core.code-names-in-english	Write code and infrastructure names in English: variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs.
   clause.core.code-names-keep-untranslated-terms	A business term with no faithful English translation keeps its original form inside an English name: an official acronym as it is (`icms_amount`), another term without accents (`apuracao_summary`). Brazilian fiscal and legal terms are the usual case.
   clause.core.code-names-new-only	A new name follows this rule. An existing name changes only when its code is touched, a large rename takes a Spec of its own, and nothing is renamed in bulk without a request. The Pre-PR Review reports a new name that breaks the rule, and the QA gate reports text a user reads outside the Reader Language.
   clause.core.reader-language	Write what the product's users read in the Reader Language named above: screens, user messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Specs, ADRs, agent guides, commits and pull requests keep the language they already have.
   ```

4. API Contract: the canonical-terms clause. `clause.domain.canonical-language` keeps its id and its force, and its guidance becomes the text below: the Source Baseline's two sentences with two sentences inserted between them. The line is written after the definition and before `_Avoid_`, as `_In code_: ` followed by the name in backticks. `TestStandardTypeScriptStructuralClauseRetention` keeps the clause in its list and compares the guidance, with the two inserted sentences removed, to the Source Baseline text.

   ```text
   Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. In a code name, a term takes its Name in Code: the term itself, or the name its Glossary Definition records on a line `_In code_:` when the term was translated to English or kept without accents. A Spec that introduces or changes such a term writes that line. Call out a missing term instead of inventing a competing synonym.
   ```

5. API Contract: the fixtures that gain the answer. Each helper or literal that lists the answers of the required decisions gains `language.reader` with the value `English`: in `internal/baseline`, `plan_test.go` (five places, one a repository profile's `Decisions`) and `profile_alignment_test.go`; in `internal/cli`, `baseline_apply_test.go`, `baseline_human_test.go`, `baseline_incremental_verification_test.go`, `baseline_plan_test.go`, `baseline_release_gate_test.go`, `baseline_update_documentation_test.go` and `baseline_update_repository_profile_test.go`. `humanBaselineAdoptionAnswers` and the prompt sequences of `baseline_human_test.go` gain one answer after the first decision. `baseline_profile_change_test.go` gains the id in the message it expects. `clause_characterization_test.go` lists the new clauses. `internal/cli/baseline_update_test.go` and `internal/config/regeneration_declared_test.go` take the new text of the canonical-terms clause. The parity-corpus fixtures do not change: every test that reads them passed unchanged on the prototype.
6. API Contract: the published Decision Document. In `docs/user-guide/context-driven-development.md` the Decision Document gains `{"id": "language.reader", "value": "English"}` after `language.generated`, the table of suggested values gains the row `| Reader Language | `English`; answer with the language the product's users read |`, and `internal/docscontract/publicdocs_test.go` expects 17 decisions in its two assertions. The step re-records `docs/references/behavior-surfaces.json` with the command of API Contract 7.
7. API Contract: regeneration in this repository. No derived file is edited by hand. The block below lists every command a step runs, with the files it writes; the authorization record declares each command once, with these outputs. Step 1 runs the first three after its module edits, then the last two. Step 2 runs the first three and the last one, with `core.json` as the only module it changed: repeating a recorded answer is no change, so its refresh rewrites only what the changed catalog moves, expected to be the Setup Manifest alone. Step 3 runs the last four. On the prototype `make baseline-digests` rewrote the Standard TypeScript profile's `goldenDigest`, its two formatter goldens for the agent-instructions and domain guides, and, under `internal/baseline/testdata/`, `catalog.digest`, `catalog.normalized.json`, `catalog.diagnostics.golden.json` and the five plan goldens. Each Task writes "the next version" of a module or a skill, never a number.

   ```text
   go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
     internal/baseline/module-versions.json
     internal/baseline/assets/modules/core.json
     internal/baseline/assets/modules/context-workflow.json
   make baseline-digests
   go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text --decision language.reader=English
     docs/agents/agent-instructions.md
     docs/agents/domain.md
     docs/agents/setup-context.json
   make skills-sync
   go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
     skills/testdata/owned-skill-versions.json
     .agents/skills/qa-gate/SKILL.md
     skills/qa-gate/SKILL.md
     .agents/skills/roundfix/SKILL.md
     skills/roundfix/SKILL.md
     .agents/skills/setup-context-driven/SKILL.md
     skills/setup-context-driven/SKILL.md
   go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
     docs/references/behavior-surfaces.json
   go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
     docs/references/coverage-record.json
   ```

8. API Contract: the ADR. One accepted ADR with `knowledge: true`, written by Build Order step 1 under the next free number, before the Managed Refresh runs. Its path joins the Standard Export List of `.secondbrain-export` in order, as the update would write it. It states: the Language Rule and where its clauses live; the Reader Language as a required decision with the suggested value `English`, and why a suggestion is kept; the Name in Code and its line; the Code Language Exception as an optional decision of language and source, its two conditions, the source as copied text and why no path is checked; the clause gate on presence; why the new rule is outside `requiredRules`; and why the check is the Pre-PR Review and the QA gate. Its considered options are the ones `## Decisions` names.
9. API Contract: the exception's value. `normalizeCodeLanguageException` refuses a value with one of the reasons below, in this order of checks, and also refuses a language or a source that `validateStructuredRenderText` refuses as a single line. `validateDecisionValue` calls it for the type, so `baseline plan`, `baseline update --decision` and a Decision Document all refuse the same way and write nothing. The catalog accepts the type in its type list, counts it among the types with exactly one render binding, and for an optional decision of this type refuses a `default` or a `suggestion` instead of requiring a default.

   ```text
   not an object                  must be an object with language and source
   a field other than the two     accepts only language and source
   source absent or blank         requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception
   language absent or blank       requires a language
   language English, any case     names a language other than English, which is the rule and needs no exception
   ```

10. API Contract: the clause gate. `appliesWhen` holds `decision` and exactly one of `equals` and `present`. With `present`, the value is a boolean and the decision is optional, of any type; `clauseApplies` returns whether the decision is recorded, compared with the gate. With `equals`, validation and evaluation are unchanged. The Source Baseline transition names presence in the reason of a clause a `present` gate takes out.
11. API Contract: the exception in the guide. `clause.core.code-names-in-english`, `clause.core.code-names-keep-untranslated-terms` and `clause.core.code-names-new-only` gain a gate on `language.code_exception` with `present` false. The rule gains `clause.core.code-names-follow-the-existing-codebase`, `mandatory`, gated with `present` true, whose guidance is the first block below. The token renders nothing while the decision is absent; while it is recorded it renders a line break and the sentence of the second block, the language escaped as Markdown text and the source in a code span. The five built-in profiles list the decision in `entryDecisions` after `language.reader`. `clause_characterization_test.go` lists the clause, and the helper `manifestInputDecisionsForProfile` of `update_test.go` skips the decision, which has no suggested value to record.

   ```text
   Write new code and infrastructure names in the language of the existing code, as the Code Language Exception named above declares. The exception covers the codebase that existed when its owner took the cited decision and never starts a new codebase in another language; that owner answers for it. The Pre-PR Review reports a new name in another language, and the QA gate reports text a user reads outside the Reader Language.
   ```

   ```text
   The Code Language Exception: code and infrastructure names are written in <language>, as the existing code is. Source of the owner's decision: `<source>`.
   ```

12. API Contract: the interactive flow. `promptBaselineDecisions` in `internal/cli/baseline_human.go` never prompts for `language.code_exception`: it carries a recorded value into the answers unchanged and otherwise skips the decision.
13. API Contract: the review instruction. `buildReviewPrompt` writes `reviewLanguageRuleInstruction` and a line break directly after the absence-claim instruction; the constant holds the text below. The Roundfix skill's review reference, its mirror and `docs/user-guide/commands/review.md` quote it, and a test of `internal/docscontract/command_documentation_test.go` reads the constant from the Go source and holds the three files to it.

   ```text
   Read `docs/agents/agent-instructions.md` when the repository has it. When it states a language for code and infrastructure names, report as a finding each name the candidate adds that breaks that rule: an identifier, a file or folder, a database name or another kind the rule lists. Do not report a name the candidate did not add or a business term the guide lets keep its original form. When the guide states no such rule, report nothing for it.
   ```

14. API Contract: the QA instruction. `.agents/skills/qa-gate/SKILL.md` gains, after `### Frontend sweep`, the section `### Reader Language` with the paragraph below. `make skills-sync` copies it to `skills/qa-gate/SKILL.md`, and the owned-skill record command raises the version. The `### QA settlement` section is untouched.

   ```text
   When `docs/agents/agent-instructions.md` names a Reader Language, read the text a user reads on every surface the Spec declares: screen copy, command output, messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Record each text written in another language as a finding with the surface, the text and the expected language, and classify it by user impact as any other finding. A business term the guide lets keep its original form is not a finding. When the guide names no Reader Language, report nothing for it.
   ```

15. API Contract: the skills and the guides. `.agents/skills/roundfix/references/baseline.md` and `docs/user-guide/commands/baseline.md` gain, beside the `--decision` paragraph, what `language.reader` is, that an update without it stops with exit `3`, that `--adopt-suggested` records `English` and so is right only for a product whose users read English, the `--decision 'language.code_exception={"language":"<language>","source":"<source>"}'` form, its refusal with exit `2`, and that a repository-owned profile lists `language.reader`, and lists `language.code_exception` before it can declare one. `.agents/skills/setup-context-driven/SKILL.md` gains two sentences in `## Project decisions and Spec constraints` that name both decisions. `docs/user-guide/context-driven-development.md` gains a short subsection on the Language Rule, the Name in Code line and the exception. The `roundfix` and `setup-context-driven` skills are mirrored and versioned as API Contract 14 says for `qa-gate`.
16. API Contract: the glossary. `CONTEXT.md` gains the five Glossary Definitions the PRD's `## Glossary` declares, written through `domain-modeling` from the definitions of Core Features 1, 2, 3, 4 and 7, each citing the ADR. None takes an `_In code_` line: in this repository each term's name in code is the term itself.

### Surface Transcripts

Both blocks were copied from the prototype binary, run on a clone of this repository at `8ac0224b`. The QA gate runs them on a temporary clone checked out at the Delivery Base, whose Setup Manifest lacks the answer.

1. Surface Transcript: an update without the Reader Language stops and names the decision.

   ```transcript
   $ roundfix baseline update --repo <repository> --no-skills --no-history --format text
   stdout:
   Baseline update: action required
   Category: decision
   Result: the current Baseline catalog requires new decisions: language.reader
   Prior catalog: <digest>
   Current catalog: <digest>
   File changes: 0
   Misplaced repository rules: 0
   Retention evidence: 0
   Declared decisions: 0
   New decision: language.reader suggested="English" (The Reader Language: the language of what the product's users read.)
   Skills: not run
   Skills installed: 0
   Skills restored: 0
   Skills drifted: 0
   Next action: answer the named decisions through roundfix baseline, or rerun roundfix baseline update with --adopt-suggested
   stderr:
   exit: 3
   ```

2. Surface Transcript: a Code Language Exception without a source is refused and nothing is written.

   ```transcript
   $ roundfix baseline update --repo <repository> --no-skills --no-history --format text --decision language.reader=pt-BR --decision 'language.code_exception={"language":"Portuguese"}'
   stdout:
   Baseline update: failed
   Category: invalid
   Result: declare Setup Manifest decisions: resolve Setup Manifest decisions: normalize Baseline decision "language.code_exception": requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception
   Prior catalog: <digest>
   Current catalog: <digest>
   File changes: 0
   Misplaced repository rules: 0
   Retention evidence: 0
   Declared decisions: 0
   Skills: not run
   Skills installed: 0
   Skills restored: 0
   Skills drifted: 0
   Next action: correct the declared decisions and rerun roundfix baseline update
   stderr:
   roundfix: baseline update failed: declare Setup Manifest decisions: resolve Setup Manifest decisions: normalize Baseline decision "language.code_exception": requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception
   exit: 2
   ```

## Coverage Map

- Goal 1 → Baseline text, Reader Language decision (API Contracts 2 and 3)
- Goal 2 → Reader Language decision, invariant 1 (API Contract 1)
- Goal 3 → Baseline text (the canonical-terms clause, API Contract 4)
- Goal 4 → Review and QA (API Contracts 13 and 14)
- Goal 5 → Baseline text, Code Language Exception (`clause.core.code-names-new-only`, API Contracts 3, 9 and 11)
- Story 1 → `clause.core.code-names-in-english`
- Story 2 → `clause.core.reader-language`, the template line
- Story 3 → `clause.core.code-names-keep-untranslated-terms`, the `_In code_` line of API Contract 4
- Story 4 → `language.reader`, Surface Transcript 1
- Story 5 → `language.code_exception`, the clause gate, API Contract 11
- Story 6 → `reviewLanguageRuleInstruction`
- Story 7 → the `### Reader Language` section of the `qa-gate` skill
- Story 8 → the `## Decisions` line of this document on the adopters' messages; no component, the session that follows the release sends them
- Core Feature 1 → `clause.core.code-names-in-english`
- Core Feature 2 → API Contracts 1 and 2, `clause.core.reader-language`, Surface Transcript 1
- Core Feature 3 → `clause.core.code-names-keep-untranslated-terms`
- Core Feature 4 → API Contract 4
- Core Feature 5 → API Contract 4, invariant 4
- Core Feature 6 → `clause.core.code-names-new-only`
- Core Feature 7 → API Contracts 9, 10, 11 and 12, invariants 2, 3 and 5, Surface Transcript 2
- Core Feature 8 → API Contract 13, invariant 7
- Core Feature 9 → API Contract 14
- Core Feature 10 → API Contracts 1, 5 and 7; the profile decision gap as it is
- Core Feature 11 → API Contract 4, invariants 4 and 6, the Managed Refresh of API Contract 7
- Core Feature 12 → API Contracts 6, 8, 15 and 16
- Core Feature 13 → no component; the session that follows the release sends the messages
- Success Metric 1 → `TestEveryBuiltInProfileRendersTheLanguageRule`
- Success Metric 2 → `TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles`, `TestBaselineUpdateFleetSweep`, `TestStandardTypeScriptStructuralClauseRetention`, Surface Transcript 1
- Success Metric 3 → `TestACodeLanguageExceptionSwapsTheCodeClauses`, Surface Transcript 2
- Success Metric 4 → `TestARepositoryProfileWithoutTheReaderLanguageIsRefused`
- Success Metric 5 → `TestTheReviewPromptCarriesTheLanguageRuleInstruction` and the documentation contract of API Contract 13
- Success Metric 6 → `TestTheQAGateReportsTextOutsideTheReaderLanguage`, `TestAuthorialSkillSync`, `TestEveryOwnedSkillVersionIsRecorded`, `TestSettlementGuidanceIsOneTable`
- Success Metric 7 → API Contracts 8 and 16, the glossary command of `roundfix spec check`
- Success Metric 8 → the clause texts of API Contracts 3 and 4, measured at 1,433 bytes on the prototype
- Success Metric 9 → the clause texts of API Contracts 3 and 4, compared with the passages of the PRD's `## Acceptance evidence`
- Success Metric 10 → Build Order, which renames no file and names no file of the Gate Time Rule
- Success Metric 11 → declared under Unreachable Acceptance; the rule and the adopters' messages are what the sweep reads
- Success Metric 12 → declared under Unreachable Acceptance; the review instruction and the `_In code_` line are what the sweep reads
- Success Metric 13 → declared under Unreachable Acceptance; no component

Deviations from the PRD: None.

## Absence Evidence

1. `rule.core.code-and-reader-language` and its five clauses — command: `rg -n 'code-names|reader-language|code-and-reader-language' internal/baseline/assets`; result: no match.
2. The decisions `language.reader` and `language.code_exception` — command: `rg -n 'language\.reader|language\.code_exception' internal docs/user-guide docs/agents .agents/skills skills`; result: no match.
3. The `present` form of the clause gate and the decision type — command: `rg -n -g '!*_test.go' 'gate\["present"\]|code-language-exception|CodeLanguageException' internal`; result: no match. `rg -n '"present"' internal/baseline/assets/modules` also finds nothing: the three gates the catalog holds use `equals`.
4. The review instruction — command: `rg -n -i -g '!*_test.go' 'language' internal/cli/review.go internal/cli/review_conventions.go`; result: no match.
5. The QA instruction and the texts of API Contract 15 — command: `rg -n -i 'reader language|another language|wrong language' .agents/skills/qa-gate/SKILL.md .agents/skills/roundfix docs/user-guide`; result: no match.
6. The `_In code_` line — command: `rg -n '_In code_' CONTEXT.md .agents/skills docs/agents internal/baseline/assets`; result: no match.
7. The five Glossary Definitions — command: `rg -n '^\*\*(Language Rule|Reader Language|Untranslated Term|Name in Code|Code Language Exception)\*\*' CONTEXT.md`; result: no match.
8. The ADR — command: `rg -l -i 'reader language|untranslated term|name in code' docs/adr CONTEXT.md docs/agents`; result: no match.

## ADR Alternatives

- ADR-0118: implements its shape: the obligation is a clause and the value a decision with a suggested value, which a repository takes only by `--adopt-suggested`. It rejected rendering a value from a repository-local file that the catalog reads; the exception's source is copied into the Setup Manifest and no file is read for it.
- ADR-0186: obeys it. The clauses name the Pre-PR Review, the QA gate and the owner of the cited decision as who holds what the product does not enforce, and cite no Spec and no ADR of this repository.
- ADR-0203: obeys it. The rule is promoted on evidence from six repositories, and it supersedes no clause, so no clause declares `replaces`.
- ADR-0222: obeys it. The clause gate keeps a guide to what holds for its repository. Its rule that a reworded clause takes a new identity was narrowed by ADR-0257, which this design follows.
- ADR-0257: obeys it. The reworded clause keeps its identity, and each rule is stated once: the canonical-terms clause carries the Name in Code and no second clause repeats it.
- ADR-0263: repeats its shape for the exception: an optional decision that an update records through `--decision`. It rejected a required decision where the answer is already recorded; the Reader Language is required because no repository has recorded it.
- ADR-0372: uses its refusal as it is. A repository profile that does not list `language.reader` is refused with the decision, the template and the entry to add.
- ADR-0189: obeys it. Each changed owned skill takes a version that names its new content.
- ADR-0233: obeys it. The owned-skill record command chooses the version.
- ADR-0250: obeys it. The record step chooses each module's version, and the coverage record is re-recorded by its command.
- ADR-0256: obeys it. The skills that describe the changed commands change in the same delivery.
- ADR-0280: obeys it. The Task that edits a user-guide page records its own surface fingerprint, and the glossary entries follow the glossary's own test.
- ADR-0360: obeys it. The Vendored Skill that owns the glossary's format is not edited; the line that records a name in code is defined by a Baseline clause.
- ADR-0196: obeys it. The finding grammar, the Delivery Conventions and their version are unchanged; the instruction is a sentence of the prompt and no convention.
- ADR-0244: obeys it. The PRD declares the five terms, and one Task writes them.
- ADR-0438: obeys it. The ADR's path is the one line added by hand to the Standard Export List, in order.
- ADR-0439: obeys it. The ADR carries `knowledge: true`.
- ADR-0371: obeys it without a component. Adopters learn of the release through the Update Check.
- ADR-0383: obeys it without a component. The adopters' messages expect no reply and no Task sends one.
- ADR-0431: obeys it without a component. The rule leaves a user-level file for the Baseline.
- ADR-0436: obeys it. No rule blocks the goal, and the design contradicts no accepted ADR, so no Rule Challenge is asked.
- ADR-0437: obeys it without a component. No rule was kept by a Rule Challenge, so no Kept Line is written.
- ADR-0453: obeys it. No file of the Gate Time Rule, no Makefile and no workflow changes, and the new tests run in-process.
- ADR-0005: not applicable — the number names another repository's decision record in the PRD's sources; this repository has no such record.
- ADR-0010: not applicable — it settles failed Batches of the resolve cycle; the PRD cites another repository's record of that number.
- ADR-0141: not applicable — it settles where Review Runs execute; the PRD cites another repository's record of that number as the exception's first source.
- ADR-0139: not applicable — Run keys are untouched.
- ADR-0150: not applicable — branch naming is untouched.
- ADR-0182: not applicable — settlement facts are untouched.
- ADR-0228: not applicable — the branch-prefix sentence keeps its bytes and its place as the last opening line.
- ADR-0248: not applicable — history is not sanitized or converted.
- ADR-0252: not applicable — the selective gate and Contract Relevance are unchanged; the new contract tests join existing files.
- ADR-0253: not applicable — the Full Contract Run is unchanged.
- ADR-0258: not applicable — the QA rerun classes bind the Task Graph, which a later stage writes.
- ADR-0259: not applicable — no test changes process-wide state; each new test runs in parallel.
- ADR-0272: not applicable — the review record is unchanged.
- ADR-0286: not applicable — the Correction Review is unchanged beyond carrying the prompt's new sentence in its header.
- ADR-0317: not applicable — Result Status Notes are untouched.
- ADR-0319: not applicable — the follow-up order is untouched.
- ADR-0330: not applicable — the Secondbrain stays read-only and nothing is written to it.
- ADR-0341: not applicable — a profile change carries both new decisions as it carries any recorded decision.
- ADR-0342: not applicable — settlement is unchanged.
- ADR-0349: not applicable — this stage asks no release question; the authorization record carries the maintainer's answer.
- ADR-0359: not applicable — lock entries of installed skills are untouched.
- ADR-0364: not applicable — required external skills are untouched.
- ADR-0380: not applicable — no skill is replaced or retired.
- ADR-0386: not applicable — user-level skill copies are untouched.
- ADR-0390: not applicable — message precedence is untouched.
- ADR-0393: not applicable — `skills-lock.json` is untouched.
- ADR-0395: not applicable — queue stops are untouched.
- ADR-0396: not applicable — the adopted source asked for no completion notice.
- ADR-0399: not applicable — carrier inventory is untouched.
- ADR-0411: not applicable — no misplaced rule is named or moved; an adopter's own line that repeats the rule is its owner's to remove.
- ADR-0412: not applicable — release impact rules are untouched.
- ADR-0416: not applicable — the Update Check's report is untouched.
- ADR-0442: not applicable — this stage leaves no Open Authoring Question for the session to adopt.
- ADR-0444: not applicable — no harness permission file is written.
- ADR-0448: not applicable — no Skill Activation changes owner.
- ADR-0452: not applicable — the platforms of Verification are untouched.

## Integration Points

None at run time. After the release, the session that follows it sends the adopters' messages of Core Feature 13 and one message to the skills catalog about the sentence of its glossary skill, through ai-memory; no Task sends a message. The reviewer and the QA gate read `docs/agents/agent-instructions.md` of the repository they work in, as a file.

## Testing Approach

Every seam exists; no test seam is new.

- Catalog and rendering, in-process, package `internal/baseline`, in a new file: `TestEveryBuiltInProfileRendersTheLanguageRule` renders each built-in profile and asserts the Reader Language line, the four clause lines with their force and the reworded domain clause, and that the same assertion fails on a catalog copy with one clause removed. `TestACodeLanguageExceptionSwapsTheCodeClauses` covers the three cases of Success Metric 3 and invariant 3. Table tests cover each refusal of API Contract 9, each form of the gate and the catalog's refusal of a default on the new type.
- The update, package `internal/cli`, in a new file, on the existing update fixtures: `TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles` (stop, answer, three files, second update plans nothing), the Standard TypeScript copies of `TestBaselineUpdateFleetSweep` for the clause accounting, `TestTheUpdateNamesTheReaderLanguageDecision` and `TestAnExceptionWithoutASourceIsRefused`, which assert Surface Transcripts 1 and 2 line by line, `TestARepositoryProfileWithoutTheReaderLanguageIsRefused`, and a test that the interactive flow asks nothing for the exception.
- The review prompt, package `internal/cli`, in a new file: `TestTheReviewPromptCarriesTheLanguageRuleInstruction` asserts the constant's place after the absence-claim instruction and the unchanged grammar and conventions.
- Skills, package `skills`: `TestTheQAGateReportsTextOutsideTheReaderLanguage` pins the section in the embedded skill.
- Documentation contracts join existing files, because `TestRepositoryContractTestsDeclareTheirRelevance` pins the number of `docscontract` files: the review quotation in `internal/docscontract/command_documentation_test.go`, the Decision Document in `publicdocs_test.go`.
- Every added test changes `docs/references/coverage-record.json`, re-recorded by its command in the step that adds it.

## Build Order

1. The Language Rule clauses, the Reader Language decision and the decision record: API Contracts 1 to 8, with the first token of API Contract 2 only. It edits the catalog sources, the fixtures, the Decision Document and its test, writes the ADR and its export line, runs the commands API Contract 7 gives it and adds the tests of Success Metrics 1, 2 and 4.
2. The Code Language Exception: API Contracts 9 to 12 and the second token of API Contract 2, the commands API Contract 7 gives it, and the tests of Success Metric 3 (depends on: 1).
3. The review instruction, the QA instruction, the skills, the user-guide pages and the glossary: API Contracts 13 to 16, the commands API Contract 7 gives it, and the tests of Success Metrics 5 and 6 (depends on: 2).
4. The QA gate (depends on: 1, 2, 3).

Steps 1 to 3 are one chain because each re-records `docs/references/coverage-record.json`, steps 1 and 2 edit the same module and catalog files, and step 3 describes the refusals step 2 builds. No step depends on another Spec.

## Risks & Considerations

- The suggested value is wrong for most adopters. Nine of the eleven read Portuguese, and `English` appears in the stop, in `--adopt-suggested` and in the repair line of the profile decision gap. Each is explicit, the adopters' message names the answer, the skill says when adoption is right, and the `update` Skill Mode already never passes `--adopt-suggested` unasked.
- Step 1 is wide: fourteen test files beside the catalog. The edits are mechanical and listed in API Contract 5. The prototype reached a green `internal/baseline` and a green Full Contract Run, and left only the scripted answers of the interactive tests in `internal/cli` to extend.
- The guides grow by 1,433 bytes against a ceiling of 1,500. A Task that changes a clause text measures again before it settles.
- The update records or changes a decision and has no flag that removes one, so this design builds no withdrawal of a recorded exception. The first repository that needs it brings the evidence for that work.
- The source is text. A path that does not exist or a decision that does not say what the declaration claims is the owner's error, which the guide's own sentence exposes to every reader.
- `--decision` takes the exception as JSON, which a shell quotes with single quotes; the skill and the guide print the whole command.
- The reviewer and the QA gate are models. Both instructions are bounded to what the candidate adds and to what a user reads, and both report nothing where the guide states no rule.

## Glossary

None.

## Decisions

- The Reader Language carries the suggested value `English`. A decision without one would end an update as an invalid manifest instead of the stop Core Feature 2 asks for, and would depart from the shape ADR-0118 gives a decision. The ADR of API Contract 8 records it.
- The code-language clauses and the reader clause sit in a new rule outside `requiredRules`. Inside `rule.core.english-guidance` each would need a hand-written Source Baseline row in five governed files; extending the one existing clause would make a text no gate can split.
- The canonical-terms clause is reworded by insertion and carries the Name in Code; it gains no neighbour, which would need a Source Baseline row and state the rule twice.
- The Name in Code is recorded on a line `_In code_:`, parallel to `_Avoid_`, so that every session writes it the same way and a later check can read it.
- The exception is one optional structured decision with no default, and the clause gate gains `present`. Two string decisions would need a default for a source, and an optional enum cannot hold a language and a source.
- The exception's source is copied text. A checked path would make a decision's validity depend on the repository's files at every update, would refuse a decision recorded outside a file, and would still accept any file.
- `English` is refused as an exception's language, so a guide never drops the three clauses for the language they already state.
- The update's next-action sentence is unchanged; the `--decision` form is taught by the skill, the user guide and the adopters' message.
- The plan is three implementation steps and the QA gate. Splitting step 1 would leave a clause that names the Reader Language without its line, or an ADR apart from the clauses it describes.
