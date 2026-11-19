---
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
prd: _prd.md
created: 2026-11-19
---

# The Baseline owns the documentation rules and the Secondbrain export — Technical Spec

## Executive Summary

The design has three layers. Text: four clauses of the docs-layout rule and one of the cross-project rule are reworded under their ids, and two new rules carry the Lesson Learned and the export list; the `archive-spec` skill and five references of the Roundfix skill gain their steps. Code: `roundfix baseline update` gains a documentation part, planned beside the history cleanup and outside the guidance plan, that keeps the Standard Export List and fills each Amendment List; `roundfix doctor` gains two warning checks; and the Release Plan counts two more paths as maintenance-only. Content: this repository's export list and decision records follow the rules.

The primary trade-off is where the standard list lives in the update. As a managed artifact of the guidance plan it would get the preservation proof and one transactional apply, and its content follows the repository's ADRs, so the digest the Setup Manifest records would go stale with each marked ADR and the check every settlement makes of the Setup Manifest would stop a Run for a line of documentation. The design therefore plans the documentation part in the command, as the history cleanup is planned, binds it into the Plan Digest and applies it after the guidance. It gives up the one transaction: a failure while the documentation part is written leaves the guidance applied and names the two paths to restore.

A second trade-off decides where the new clauses live. The docs-layout rule is one the Standard TypeScript profile requires, and the catalog refuses a new clause in such a rule until a Source Baseline row is written by hand. Existing clauses are reworded under their ids, which the transition records as retained, and the lesson and the export list are new rules that no profile requires (`references/2026-11-19-prototype-and-governed-paths.md`, section 2). A third is the size of the docs-layout region: the text is kept to 12,948 bytes and the ceiling of the lean Baseline rises from 11,000 to 13,000.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity; the rule ids `rule.context.lessons` and `rule.secondbrain.knowledge-export`, their clause ids and the check codes with the `DR-EXPORT-` and `DR-DOCS-` prefixes follow the forms their owners already use. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added, and no test or Verification runs the sync, reaches another repository or asks a person, a provider or ai-memory anything. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands, with ADR-0438, ADR-0439 and ADR-0440 as the decisions this design implements; ADR-0278, ADR-0248 and ADR-0412 as the ones a Task gives a dated notice; ADR-0282, ADR-0254, ADR-0122, ADR-0220, ADR-0411, ADR-0186, ADR-0222, ADR-0257, ADR-0203, ADR-0256, ADR-0349, ADR-0417 and ADR-0420 as the ones it obeys or uses; and ADR-0145, ADR-0287, ADR-0100, ADR-0330 and ADR-0396 listed to record how they bound it. ADR-0092, ADR-0120, ADR-0247, ADR-0251, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0277, ADR-0281, ADR-0290, ADR-0295, ADR-0299, ADR-0308, ADR-0342, ADR-0352, ADR-0366, ADR-0368, ADR-0391, ADR-0421, ADR-0422, ADR-0423, ADR-0424, ADR-0426 and ADR-0432 are listed only to record why they do not apply. This stage adds no ADR. The standard list is found by its own pair of markers and is not a Managed Region of an instruction carrier (ADR-0145: "A managed region is trusted by its marker"), and the clause texts name no ADR and no Spec of this repository (ADR-0186: "Baseline guidance states what the product does, in adopter-neutral words"). Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18 for this Spec, quoted in the authorization record, which keeps the grant of every owned skill and of the Governed Paths the Spec declares, Baseline modules and records included. The design changes no Makefile, workflow, Go module, lint or formatter file. It changes three Baseline modules, the Standard TypeScript profile's pin and three of its formatter goldens, this repository's three guides and Setup Manifest, two owned skills and the coverage record. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export/_authorization.md`; bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/corrections.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/release.md`, `.agents/skills/roundfix/references/setup.md`, `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/secondbrain.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Every component extends something that exists. The new files are six Go files with their tests, two contract test files and the folder `docs/lessons/`, which appears with its first lesson.

### Baseline text

`internal/baseline/assets/modules/context-workflow.json`: four clauses of `rule.context.docs-layout` take the texts of API Contracts 1 to 4 under their ids, and the module gains `rule.context.lessons`, version 1, with the module's coverage id and the two clauses of API Contract 5. The guide `guide.docs-layout` lists the new rule after `rule.context.docs-layout`, and its guide version rises by one. `core.json`: the clause `clause.core.message-the-owning-project` takes the sentence of API Contract 6. `secondbrain.json`: the module gains `rule.secondbrain.knowledge-export`, version 1, with the coverage id `coverage.secondbrain` and the three clauses of API Contract 7; the guide `guide.secondbrain` lists it last, and its guide version rises by one. No profile lists either new rule under `requiredRules`. Each module version rises through its record command.

The module edit moves the files the prototype measured: the module version record, the Standard TypeScript profile's formatter pin and its goldens of the three guides, three catalog files and four plan goldens, which `make baseline-digests` writes, and this repository's three guides and Setup Manifest, which the Managed Refresh writes. `internal/baseline/clause_characterization_test.go` gains the five new clause ids. Three existing tests pin text that changes and are updated with it (`## Testing Approach`).

### ADR header and the lesson form

`internal/spec` already classifies an ADR (`ClassifyADR`) and splits frontmatter. It gains `adr_header.go`, which reads the Knowledge Mark and the Amendment List, finds the amendment notices of a body, and writes the Amendment List by inserting or replacing one line; and `lesson.go`, which reports how a file departs from the form of a Lesson Learned. Both are pure functions over bytes.

### Export list

`internal/baseline` already owns the checks the update and the Doctor Command share (`MisplacedRepositoryRules`, `ReadClaudeInstructionLink`). It gains `export_list.go`: the standard items, the two markers, the planner that turns the bytes of a file and the repository's facts into the bytes the update writes, and the reader of the repository's own lines. It gains `docs_layout.go`: the folders and root documents the directory clause names. Both take their facts as values, so neither reads the configuration nor imports the consistency check.

### Documentation part of the update

`internal/cli/baseline_update_documentation.go` gathers the facts (the Setup Manifest's modules, the Spec Root and the History Root from the configuration, the glossary files from the consistency check's resolver, the ADRs and the folders from the working tree), plans the export list and the header fills, folds the plan into the Plan Digest, reports it under the `documentation` member and in the text output, and applies it. `runBaselineUpdateCommandWithSkillsStage` calls it where it calls the history plan: after the guidance plan is built, before the digest is final, and after the history apply. `--no-history` and `--no-skills` do not skip it. Its file changes are also appended to `fileChanges`, so the Update Check and the `update` Skill Mode count and show them without a change of shape.

### Doctor checks

`internal/cli/doctor_documentation.go` builds two results from the same facts, `baseline-export` and `baseline-docs`, and `runDoctorCommand` appends them after `baseline-claude-link`. `internal/cli/health.go` gains the two names. Neither result is ever `failed`.

### Release Plan

`internal/releaseplan/classifier.go`: `docs/lessons/` joins the planning-evidence prefixes and `.secondbrain-export` is maintenance-only. The rule holds for every repository and reads no configuration, so it needs nothing from the adopter test of Spec 0414.

### Skills, guides and glossary

The canonical `archive-spec` skill gains the section of API Contract 14. The canonical Roundfix references gain the texts of API Contracts 15 to 19. `make skills-sync` writes the mirrors and the version record command raises both skills. The user-guide pages `commands/baseline.md`, `commands/doctor.md`, `context-driven-development.md`, `release-runbook.md` and `usage.md` and the `release plan` help gain the texts of API Contract 20, and the Behavior Surface Record is recorded again after the last guide edit of each Task. `CONTEXT.md` gains the entries of API Contract 21.

### This repository

`.secondbrain-export` becomes the file of API Contract 22. Every file under `docs/adr/` gains the Knowledge Mark, the update fills the Amendment Lists, and ADR-0248, ADR-0278 and ADR-0412 gain the notices of API Contract 23. `internal/docscontract/secondbrain_export_test.go` is rewritten to the new rule and the coverage record is recorded again.

## Implementation Design

### Interfaces

```go
// internal/spec/adr_header.go
type ADRHeader struct {
	Frontmatter bool     // the file opens with a frontmatter block
	Status      string   // lower case, "" when the header has none
	Knowledge   *bool    // nil when the mark is absent or not a boolean
	AmendedBy   []string // identifiers as written, such as ADR-NNNN
}

func ReadADRHeader(content []byte) ADRHeader
func ADRAmendmentNotices(content []byte) []string // identifiers the body names, ascending, each once
func WithADRAmendedBy(content []byte, ids []string) []byte
```

```go
// internal/spec/lesson.go
func LessonLearnedProblems(name string, content []byte) []string // nil when the file follows the form
```

```go
// internal/baseline/export_list.go
const (
	ExportListPath        = ".secondbrain-export"
	ExportListBeginMarker = "# roundfix:standard-export:begin"
	ExportListEndMarker   = "# roundfix:standard-export:end"
)

type ExportListFacts struct {
	GlossaryFiles []string // repository-relative, in the resolver's order
	KnowledgeADRs []string // paths of ADRs marked true and not retired, ascending
	ProcessPaths  []string // a folder ends with a slash
	DocsEntries   []string // the immediate children of docs/, a folder ends with a slash
}

type ExportLoss struct {
	Path string `json:"path"`
	Kind string `json:"kind"` // "process" or "other"
}

type ExportListPlan struct {
	Action         string // "none", "unchanged", "create", "update" or "refused"
	Content        []byte // the whole file after the plan
	Standard       []string
	Added, Removed []string
	Converted      []string
	StopsExporting []ExportLoss
	Refusal        string
}

func PlanExportList(current []byte, exists bool, facts ExportListFacts) ExportListPlan
func ExportVariationDepartures(current []byte, facts ExportListFacts) []DocumentationDeparture
```

```go
// internal/baseline/docs_layout.go
type DocumentationDeparture struct {
	Code, Path string
	Line       int
	Detail     string
}

type DocsLayoutFacts struct {
	RootDocuments      []string // Markdown files at the repository root
	DocsEntries        []string // as in ExportListFacts
	SpecRoot, History  string   // repository-relative, "" when outside the repository
	GlossaryFiles      []string
	ReadmePresent      bool
	SpecificRepository []byte // the repository's own rules guide, nil when absent
}

func DocsLayoutDepartures(facts DocsLayoutFacts) []DocumentationDeparture
```

Invariants:

```text
1. PlanExportList and WithADRAmendedBy are functions of their arguments: the same bytes and facts give the same bytes.
2. Outside the two markers PlanExportList changes a byte only when the file has neither marker, and then only by removing the lines it reports under Converted.
3. With exactly one begin marker followed by exactly one end marker, the plan replaces the bytes from the begin line through the end line and nothing else.
4. With any other count or order of markers, and at least one marker present, the plan is "refused", Content equals the input and Refusal names the marker.
5. The Standard list never holds a process path, a retired ADR or an ADR whose Knowledge Mark is not true.
6. WithADRAmendedBy changes or adds exactly one line, or prepends a header that holds only that line; every other byte, updated_at included, is kept.
7. The Amendment List written is the union of the identifiers the header already lists and the identifiers the notices name, ascending by number; an identifier a person wrote is never removed, and an ADR never lists itself.
8. A notice is read outside fenced code and inline code only.
9. Applying a plan and planning again gives "unchanged" and no header fill.
10. No result of the two Doctor checks has the status failed, and neither check writes a file.
```

### Data Models

The Standard Export List, as the planner writes it. The second line is fixed, the glossary files and the ADR paths follow the repository, and every other line is constant:

```text
# roundfix:standard-export:begin
# Kept by the Baseline: roundfix baseline update rewrites the lines between these two markers.
README.md
CONTEXT.md
docs/user-guide/
docs/references/
docs/design/
docs/lessons/
docs/adr/0007-a-credit-note-never-reopens-a-closed-period.md
# roundfix:standard-export:end
```

- The glossary files are the root `CONTEXT.md`, `GLOSSARY.md`, `CONTEXT-MAP.md` and `GLOSSARY-MAP.md` that exist and each glossary file a map links, as the consistency check already resolves them, in that order.
- The process paths are `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, `docs/backlog/`, `docs/findings/`, `docs/handoffs/`, the Spec Root and the History Root when they lie inside the repository.
- A new file is the block and a final newline. A converted file is the block, one blank line and the kept lines in their order. A line is compared after trimming blanks and a trailing slash.
- Conversion removes an include line that equals a standard line, an include line that is a process path, lies inside one or contains one (`docs/` whole), and an exclusion whose pattern lies inside a process path. For a removed include that contains a process path, each entry of `DocsEntries` that the new file no longer covers is one `ExportLoss`: `process` when it is a process path, `other` otherwise. A removed process line is one `ExportLoss` of kind `process`.
- In the repository's own lines, an entry is a line that is neither blank nor a comment. Its reason is the nearest comment line above it with no blank line between; a marker line is not a reason.

The ADR header. The writer puts `knowledge` after `status`, and `amended_by` after `knowledge` when the header has it and after `status` otherwise; a reader accepts either field anywhere in the header:

```yaml
status: accepted
knowledge: true
amended_by: [ADR-NNNN, ADR-MMMM]
```

- A notice is one of "Amended by", "Narrowed by", "Extended by" and "Superseded in part by", followed by one or more identifiers `ADR-<digits>`, each plain or as the text of a Markdown link, joined by a comma or `and`.
- A header that holds neither `status` nor any other lifecycle field leaves the ADR a legacy ADR: `ClassifyADR` already reads the status from the body when the header has none.

The `documentation` member of the update's JSON result, added after `history`:

```json
{
  "status": "pending",
  "exportList": {
    "path": ".secondbrain-export",
    "action": "update",
    "standard": ["README.md", "CONTEXT.md", "docs/user-guide/", "docs/references/", "docs/design/", "docs/lessons/"],
    "added": ["docs/lessons/"],
    "removed": [],
    "converted": ["docs/", "AGENTS.md"],
    "stopsExporting": [{"path": "AGENTS.md", "kind": "process"}],
    "refusal": ""
  },
  "adrs": {"total": 43, "unmarked": 41},
  "adrHeaders": [{"path": "docs/adr/0003-prices-are-stored-in-cents.md", "amendedBy": ["ADR-NNNN"]}]
}
```

`status` is `current`, `pending` or `applied`. `exportList.action` is `none` where the Secondbrain is turned off. Every array is present and may be empty.

A Lesson Learned follows the form when its name is `YYYY-MM-DD-<kebab-slug>.md` with a valid date, its frontmatter holds `created` with a valid date, it has one `# ` title line, and the eight lines `- Context:`, `- Tried:`, `- Went wrong:`, `- Cause:`, `- Solved by:`, `- Applies when:`, `- Does not apply when:` and `- Links:` appear once each, in that order, each with text after the colon.

### API Contracts

Each text is prescribed word for word. A clause keeps the id, the force and every sentence this section does not quote.

1. API Contract: the ADR clause. `clause.context.adr-01-template` keeps its text and its code block and gains, after the block, this paragraph: "An accepted ADR that a later accepted ADR changes in part, without replacing it, stays `accepted`. In the change that adds the later ADR, its header gains `amended_by: [ADR-NNNN]`, and its body gains, after the text that changed, a dated notice that begins `Amended by ADR-NNNN (YYYY-MM-DD)` and says what changed. `roundfix baseline update` fills the list from the notices an ADR already carries."
2. API Contract: the lifecycle exception. `clause.context.adr-03-legacy-compatibility` reads: "Treat a legacy ADR without lifecycle frontmatter as active unless its body explicitly marks it inactive. Do not rewrite existing ADRs solely to adopt lifecycle metadata. The one exception is a header field that a guide of this repository defines for existing ADRs, such as `amended_by`: it is written into them and changes no other line, and a header that holds no lifecycle field leaves its ADR a legacy ADR."
3. API Contract: the directory clause. In `clause.context.docs-one-job-per-directory`, everything before "Preserve repository-authored bytes outside setup markers" is replaced by: "Give each documentation directory one job: `docs/_inbox/` for raw notes only when a repository has no Secondbrain, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts of the product, such as an interface, API or data design, and never a decision or a plan, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/lessons/` for Lessons Learned, `docs/references/` for external pointers and durable project reference documents, including one that would otherwise sit at the repository root, the `DESIGN.md` design contract excepted, and `docs/user-guide/` for human documentation. A folder under `docs/` or a reference document at the root that this clause does not name is a variation: name its path and its job in `docs/agents/specific-repository.md`, or move its content to the directory whose job it is. `roundfix doctor` reports one that is neither." The rest of the clause is unchanged.
4. API Contract: the retirement clause. In `clause.context.findings-09-archive`, "a rejected or deprecated ADR, and a superseded ADR once an accepted successor names it with its date are deleted" becomes "a rejected or deprecated ADR, a superseded ADR once an accepted successor names it with its date, and a Lesson Learned once an accepted ADR, a guide or a skill states what it teaches are deleted". The clause still opens with the retirement sentinel and keeps its `replaces` list.
5. API Contract: the lesson rule, `rule.context.lessons`, two mandatory clauses.
   - `clause.context.lesson-01-record`, with the text "A Lesson Learned records something this repository tried, what went wrong, why, and what solved it. It is one file, `docs/lessons/YYYY-MM-DD-<kebab-slug>.md`, in this complete copyable form; keep each label as written and write the rest in the language of the repository's documentation:", then a `markdown` block with the frontmatter `created: YYYY-MM-DD`, the title `# <The lesson in one line>` and these eight lines: `- Context: <the work it happened in>`, `- Tried: <what was tried>`, `- Went wrong: <what was observed>`, `- Cause: <the proven cause>`, `- Solved by: <what solved it>`, `- Applies when: <when a reader should use this lesson>`, `- Does not apply when: <what looks alike and is different>`, `- Links: <the Finding, Spec, Pull Request or commit that holds the evidence>`.
   - `clause.context.lesson-02-test`, with the text "Write a Lesson Learned only when every part of this test holds: something was tried, it went wrong, the cause is proven, something solved it, and no accepted ADR, reference document or glossary entry already states it. Otherwise write nothing. Apply the test at two moments: when a Finding closes, in the change that answers it or retires it, and when a Spec that needed a correction, a failed QA or a retry is archived. Write the file in that change, so the Pull Request that carries it shows the lesson to a person, who may delete it."
6. API Contract: the message clause. In `clause.core.message-the-owning-project`, the two sentences from "When a Finding closes, state whether" through "pinned as `<repository>@<commit>:<path>`." become: "When a Finding closes, send no message about what it taught: the docs-layout guide says when that becomes a Lesson Learned in this repository."
7. API Contract: the export rule, `rule.secondbrain.knowledge-export`, three mandatory clauses.
   - `clause.secondbrain.export-01-standard-list`, with the text "This repository publishes to the Secondbrain only what `.secondbrain-export` lists, and that file has two parts. The first is the Standard Export List, between the lines `# roundfix:standard-export:begin` and `# roundfix:standard-export:end`: the README, the glossary files, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and each ADR whose header says `knowledge: true`, one path per line. The Baseline owns those lines: `roundfix baseline update` rewrites them, and the only line added by hand is the path of an ADR, in order, when the ADR is written with that mark. Never publish process: `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, the Spec folders, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and `docs/history/` stay out of both parts."
   - `clause.secondbrain.export-02-variation`, with the text "The second part is the repository's own lines, everything outside those two markers, which no Baseline operation changes. There the repository may add a file or a folder of knowledge that the standard list lacks, or drop a standard item with an exclusion line: `!` followed by the item as the list writes it, or `!docs/adr/*` for every ADR. Give each such line its reason in a comment line directly above it; one comment stands for the lines that follow it without a blank line. `roundfix doctor` reports a line without a reason, a process path, and a standard list that differs from what the update would write."
   - `clause.secondbrain.export-03-knowledge-mark`, with the text "Mark every ADR you write with whether the Secondbrain wants it: add `knowledge: true` or `knowledge: false` to its header, on the line after `status`. It is `true` when the decision states a rule of the business or of the product that someone outside this repository's development can use, and `false` when it is about how this repository is built, tested, organized or delivered, or about a tool it uses. An ADR without the mark is not published. The ADRs a repository already has are marked once, by the same test, by the Agent of the session that applies the Baseline update that brings this rule: in one commit that a person reviews in a Pull Request, followed by a second update so the list names them."
8. API Contract: the size of the docs-layout region. After API Contracts 1 to 5 the managed region `guide.docs-layout` of the Standard TypeScript golden is at most 13,000 bytes, and `TestTheDocsLayoutRegionStaysLean` holds that number.
9. API Contract: the closed list of clause forces gains `clause.context.lesson-01-record`, `clause.context.lesson-02-test`, `clause.secondbrain.export-01-standard-list`, `clause.secondbrain.export-02-variation` and `clause.secondbrain.export-03-knowledge-mark`, each `mandatory`.
10. API Contract: the documentation part of `roundfix baseline update`.
    - Plan. `status` is `pending` when the export list's action is `create` or `update` or `adrHeaders` is not empty, and `current` otherwise. A retired ADR is neither listed nor given a header. Where the Setup Manifest's modules lack `secondbrain`, the export list's action is `none` and the file is neither read nor written.
    - Digest. With a `pending` status the Plan Digest is the SHA-256, written `sha256:<hex>`, of the digest so far and the JSON of the export list's action, content digest and refusal and of each header fill. With `current` the digest is unchanged, so a repository with nothing to do keeps the digest it has today.
    - State. With no guidance change, no history unit and a `pending` documentation part, the result is `plan_ready` with the message `guidance matches the current Baseline catalog; <n> documentation change(s) pending` and the next action the other pending plans carry. `<n>` counts files.
    - File changes. Each file the part writes is one more `fileChanges` entry: `.secondbrain-export` with the action `create` or `update` and the managed entry `export.standard-list`, and each ADR with the action `update` and the managed entry `adr.amendment-list`.
    - Apply. After the guidance and the history are applied, each header fill is `WithADRAmendedBy` over the file's bytes at that moment, skipped when the file is gone, and the export list is planned again from the tree and written. A failure reports the category `documentation` and the next action `restore .secondbrain-export and docs/adr with git restore, then rerun roundfix baseline update`, with exit 1.
    - Text. The section is printed after the history section when the status is `pending` or `applied`, when the export list is refused, or when an ADR lacks the mark. Its lines, in this order, each printed only when it has something to say:

```text
Documentation: <current|pending|applied>
Export list: <create|update> .secondbrain-export (<n> line(s) in the Standard Export List)
- adds <line>
- removes <line>
- stops exporting <path> (process)
- stops exporting <path> (not in the standard list)
Export list refused: <refusal>; nothing is written to .secondbrain-export
ADRs without a Knowledge Mark: <n> of <m>; they are not exported
ADR headers: <n>
- <path> gains amended_by: [<identifiers>]
```

    - Refusals, word for word: `# roundfix:standard-export:begin is missing`, `# roundfix:standard-export:end is missing`, `# roundfix:standard-export:begin appears <n> times`, `# roundfix:standard-export:end appears <n> times`, and `# roundfix:standard-export:end comes before # roundfix:standard-export:begin`.
11. API Contract: the check `baseline-export` of `roundfix doctor`.
    - `skipped`, with the detail `no Setup Manifest at docs/agents/setup-context.json`, or `the Secondbrain is turned off in the Setup Manifest`.
    - `ok`, with the detail `the Standard Export List is current; <n> line(s) of the repository's own`.
    - `warn`, with the detail `<n> departure(s) from the export rules`, the next action `run roundfix baseline update, then edit the repository's own lines`, and one line per departure:
      - `DR-EXPORT-LIST-MISSING: .secondbrain-export has no Standard Export List`
      - `DR-EXPORT-LIST-STALE: the Standard Export List differs from what roundfix baseline update would write: <a> line(s) to add, <r> to remove`
      - `DR-EXPORT-MARKER: <refusal>; restore the pair of markers or delete both`
      - `DR-EXPORT-PROCESS-PATH: .secondbrain-export:<line>: <entry> publishes process; remove the line`
      - `DR-EXPORT-REASON-MISSING: .secondbrain-export:<line>: <entry> has no reason; add a comment line directly above it`
      - `DR-EXPORT-ADR-UNMARKED: <n> of <m> ADR(s) carry no Knowledge Mark and are not exported`
    - A read error is `warn` with the detail `DR-EXPORT-UNREAD: <error>` and the next action `re-run roundfix doctor after fixing the Setup Manifest`.
12. API Contract: the check `baseline-docs` of `roundfix doctor`.
    - `skipped`, with the detail `no Setup Manifest at docs/agents/setup-context.json`.
    - `ok`, with the detail `the documentation follows the docs-layout guide`.
    - `warn`, with the detail `<n> departure(s) from the documentation rules`, the next action `see each line; roundfix baseline update fills an Amendment List and removes a retired ADR`, and one line per departure:
      - `DR-DOCS-REQUIRED-MISSING: README.md is missing` and `DR-DOCS-REQUIRED-MISSING: the glossary is missing`
      - `DR-DOCS-OUTSIDE-LAYOUT: <path>: a folder the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md`, and the same line with `a file directly under docs/` or `a reference document at the root` in place of `a folder`
      - `DR-DOCS-ADR-RETIRED: <path>: a <status> ADR is still in the tree`
      - `DR-DOCS-ADR-AMENDMENT-UNLISTED: <path>: the body names <identifiers> and amended_by does not`
      - `DR-DOCS-LESSON-FORM: <path>: <problem>`, where the problem is `the name is not YYYY-MM-DD-<kebab-slug>.md`, `the header has no created date`, `the title line is missing`, `the line "<label>:" is missing`, `the line "<label>:" is empty` or `the line "<label>:" is out of order`
    - A code prints at most 20 lines, and then `<code>: <n> more`.
    - The folders the layout names under `docs/` are `_inbox`, `adr`, `agents`, `backlog`, `design`, `findings`, `handoffs`, `history`, `lessons`, `references`, `user-guide` and `workflow`, and the first folder of the Spec Root and of the History Root. `docs/README.md` is not reported. A root Markdown file is read as a reference document when its name before `.md` is in capitals, digits, hyphens and underscores, and these names are named by the layout: `README`, `AGENTS`, `CLAUDE`, `CONTEXT`, `CONTEXT-MAP`, `GLOSSARY`, `GLOSSARY-MAP`, `CHANGELOG`, `CONTRIBUTING`, `LICENSE`, `LICENCE`, `SECURITY`, `CODE_OF_CONDUCT`, `SUPPORT` and `DESIGN`. A path that `docs/agents/specific-repository.md` contains is not reported.
13. API Contract: the Release Plan. `MaintenanceOnly` is true for a path under `docs/lessons/` and for `.secondbrain-export`. No other classification changes.
14. API Contract: the `archive-spec` skill gains this section before `## Steps`:

```markdown
## Lesson Learned

Before the archive, find out whether the Spec stumbled or answered a Finding. Read five signs:

- `references/_index.md` adopts a Finding;
- `qa/` holds a QA Report whose verdict is not `pass`, or more than one report;
- the Task Graph holds a Corrective Task;
- `roundfix corrections report --spec <slug> --format json` lists a declaration other than `not-a-correction`, or an undeclared sign;
- a Task of the Run failed its Verification and ran again, or the Delivery Queue item was retried.

With no sign, write nothing. With any, apply the Lesson Learned test of `docs/agents/docs-layout.md` to what went wrong. When every part of the test holds, write the file under `docs/lessons/` in the form that guide gives and commit it with the archive, so the same Pull Request shows it to a person, who may delete it. When a part fails, write nothing and say in the report that the test was applied and which part failed. Never write a lesson to fill the folder.

When the Delivery Queue writes the archive, no Agent runs this step there: the Supervisor session applies the test as the Roundfix skill's `corrections` reference says. A repository whose docs-layout guide has no Lesson Learned rule skips this section.
```

15. API Contract: the Roundfix skill's `corrections` reference gains this section at its end:

```markdown
## Lesson Learned

A Correction is where a delivery's lesson is cheapest to state, because its declarer has just named the cause. After `roundfix corrections declare` records a Correction other than `not-a-correction`, and after a `roundfix deliver retry` that follows a repair, apply the Lesson Learned test of `docs/agents/docs-layout.md` to that cause. When every part of the test holds, write the file under `docs/lessons/` in the commit that makes the correction, on the Item Branch, so the delivery Pull Request shows it. When the item has already merged, or the queue made the correction with no session, write the file in the session's next change to the repository. When a part fails, write nothing. A repository whose docs-layout guide has no Lesson Learned rule skips this section.
```

16. API Contract: the Roundfix skill's `modes` reference, section `## --update`. The step count stays five, and every sentence the section holds at the delivery base is kept.
    - Step 3 gains, after the sentence that ends "and its `planDigest`.": "Show its `documentation` member too: the lines the Standard Export List gains and loses, each path under `stopsExporting` with its kind, each ADR header it fills and the count of ADRs without a Knowledge Mark."
    - Step 4 gains, after the sentence that ends "ask the same single question before applying either.": "Ask that same question, naming the paths, when `documentation.exportList.stopsExporting` is not empty."
    - After step 5, a paragraph: "Marking the existing ADRs: when the applied result's `documentation.adrs.unmarked` is not zero and `docs/agents/secondbrain.md` holds the rule of the Knowledge Mark, mark each ADR that lacks the mark by that rule's test, changing no other line. Then preview and apply the update once more, so the Standard Export List names the ADRs marked true. Report how many ADRs were marked each way and the path of each one marked true, and say that the marks are a judgment for a person to review in the Pull Request, in a commit that holds only the marks and the list. When the mode runs at session start because `update.mode` is `auto`, mark nothing and report the count."
17. API Contract: the Roundfix skill's `baseline` reference gains, in the part that describes `roundfix baseline update`: "`roundfix baseline update` also keeps a repository's documentation to the Baseline's rules, in a `documentation` part of its plan that the Plan Digest covers and the same confirmation applies. Where the Setup Manifest has the Secondbrain turned on, it writes the Standard Export List between the two markers of `.secondbrain-export` and no other line; the first time, it converts the list the file had and names every path that stops being exported. In every repository it fills an ADR's `amended_by` from the amendment notices in its body. It reports the count of ADRs without a Knowledge Mark and marks none. A missing or repeated marker is refused: nothing is written to that file, and the rest of the plan applies. `--no-history` does not skip this part."
18. API Contract: the Roundfix skill's `setup` reference gains, in the part that describes `roundfix doctor`: "`roundfix doctor` reports two documentation checks as warnings that never change its exit. `baseline-export`, where the Secondbrain is turned on, compares `.secondbrain-export` with the Standard Export List and reads the repository's own lines; its codes start with `DR-EXPORT-`. `baseline-docs` reads the documentation layout, the ADR headers and `docs/lessons/`; its codes start with `DR-DOCS-`. `roundfix baseline update` repairs the list, an Amendment List and a retired ADR, and the rest is an edit by hand."
19. API Contract: the release sentence. The Roundfix skill's `release` reference, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md` gain, right after the sentence that ends "commits changing only those paths are maintenance-only.", as that sentence reads on the default branch when the Task runs: "Lessons Learned under `docs/lessons/` are planning evidence too, and the export list `.secondbrain-export` is maintenance-only." The `release plan` help gains, right after its sentence about planning evidence: "Neither does docs/lessons or the export list .secondbrain-export."
20. API Contract: the user guide.
    - `docs/user-guide/commands/baseline.md` gains a subsection `### Documentation part` that states every rule of API Contract 10 and holds these sentences: "The Plan Digest covers the documentation part, and the one confirmation applies it.", "The first update converts the list the file had and names every path that stops being exported.", "A missing or repeated marker is refused, and nothing is written to `.secondbrain-export`." and "`--no-history` does not skip the documentation part."
    - `docs/user-guide/commands/doctor.md` gains the bullets `baseline-export` and `baseline-docs` after `baseline-claude-link`, the eleven codes in its code table, each with its meaning and what repairs it, and a subsection `#### Documentation rules` with the line forms of API Contracts 11 and 12 and the sentence "Both checks warn and neither changes the exit code."
    - `docs/user-guide/context-driven-development.md` gains a subsection `### Knowledge a repository publishes`, before `### Spec Event Announcements`, that describes the Standard Export List, the Export Variation, the Knowledge Mark, the Amendment List and the Lesson Learned, with the file of the PRD's first illustration, the lesson form, and these sentences: "A Lesson Learned is written only when there is a lesson, and it leaves when a rule takes it in.", "An ADR without the Knowledge Mark is not published." and "The marks of the ADRs a repository already has are one change a person reviews."
21. API Contract: the glossary. `CONTEXT.md` gains these five entries after the entry of Managed Refresh, each followed by a blank line, and the entry of Managed Refresh gains, after the sentence that ends "converts Pending History on approval.", the sentence: "It also keeps the Standard Export List and fills each Amendment List, in a documentation part of its plan that the Plan Digest covers."

```markdown
**Standard Export List**:
The lines of `.secondbrain-export` between its two marker lines, which the Baseline keeps in a repository that has the Secondbrain turned on: the README, the glossary files, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and each ADR whose Knowledge Mark is true. A Managed Refresh rewrites them and no other line of the file ([ADR-0438](docs/adr/0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md)).
_Avoid_: Export template, mirror filter, managed block

**Export Variation**:
The lines of `.secondbrain-export` outside the Standard Export List, which belong to the repository: a file or a folder of knowledge it adds, or an exclusion that drops a standard item, each with its reason in a comment above it. No Baseline operation changes them ([ADR-0438](docs/adr/0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md)).
_Avoid_: Override file, custom list, local exclusion

**Knowledge Mark**:
The header field `knowledge` of an ADR in a repository that has the Secondbrain turned on: true when the decision states a rule of the business or of the product that someone outside the repository's development can use, false when it is about how the repository is built or about a tool. An ADR without the mark is not published ([ADR-0439](docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md)).
_Avoid_: Export flag, public ADR, tag

**Amendment List**:
The header field `amended_by` of an accepted ADR that later accepted ADRs change in part, naming those ADRs in ascending order. A Managed Refresh fills it from the notices in the ADR's body ([ADR-0439](docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md)).
_Avoid_: Amended status, change log, supersession

**Lesson Learned**:
One file under `docs/lessons/` that records something a repository tried, what went wrong, the proven cause and what solved it, with when it applies and when it does not. It is written only when the Baseline's test holds, when a Finding closes or a Spec that needed a correction is archived, and it is deleted once an ADR, a guide or a skill states what it teaches ([ADR-0440](docs/adr/0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md)).
_Avoid_: Postmortem, retrospective, incident report
```

22. API Contract: this repository's export list. The file is the Standard Export List the update writes, a blank line, and these lines:

```text
# This repository publishes its README, its user guide and its lessons (maintainer, 2026-11-18).
# The glossary, the references and the design notes describe how the tool is built.
!CONTEXT.md
!docs/references/
!docs/design/
# Its decision records are about the tool; the user guide carries what other projects use.
!docs/adr/*
```

23. API Contract: the amendment notices, each written after the text it changes, with the header's `amended_by` and without a change of `updated_at`. ADR-0248 takes the first paragraph below after its paragraph that ends "A documentation contract holds both directions."; ADR-0278 takes the second after its paragraph that begins "**Lessons leave before the text.**"; ADR-0412 takes the third at the end of its decision text.

```markdown
**Amended by [ADR-0438](0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md) (2026-11-19).** The export no longer mirrors `docs/` whole: the Baseline keeps a Standard Export List that names knowledge and never the History Root, and the documentation contract holds that this repository exports no process path.

**Amended by [ADR-0440](0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md) (2026-11-19).** A lesson is no longer sent as a message: it is written as a Lesson Learned under `docs/lessons/` when its test holds, and a Lesson Learned joins what a retirement deletes once an ADR, a guide or a skill states it.

**Extended by [ADR-0438](0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md) and [ADR-0440](0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md) (2026-11-19).** `.secondbrain-export` is maintenance-only and `docs/lessons/` is planning evidence.
```

24. API Contract: the marks of this repository's ADRs. Every file under `docs/adr/` in the tree when the Task runs gains `knowledge: true` or `knowledge: false` on the line after `status`, by the test of API Contract 7, and no other line of it changes for the mark. For this repository the product is Roundfix: a decision about what a command, the Daemon, the Delivery Queue, the Baseline or a skill does for the people and repositories that use Roundfix is true, and a decision about this repository's tests, build, CI, release procedure or internal structure is false.
25. API Contract: the contract of this repository's export file. `internal/docscontract/secondbrain_export_test.go` holds two tests in place of the two it has: `TestThisRepositoryExportsNoProcessPath`, which fails when the file's own lines publish a process path or lack a reason, and `TestTheExportRuleRefusesAProcessPath`, which proves on a temporary directory that the same rule refuses `docs/` whole and `AGENTS.md`. Neither asserts that every ADR carries the mark or that the Standard Export List is current.

### Surface Transcripts

None. The new lines of `roundfix baseline update` and of `roundfix doctor` are prescribed word for word in API Contracts 10 to 12, and each is asserted through the command's entry point on a temporary repository by a named test. A whole transcript of either command cannot be fixed here: the update prints digests of the tree it runs on, the Doctor Command prints the state of the machine and takes its exit from checks this Spec does not change, and neither surface exists to copy a run from. The QA gate reproduces the lines through the built binary (`_prd.md` → Success Metrics 2 to 4).

## Coverage Map

- Goal 1 → Export list, Documentation part of the update, Baseline text (API Contracts 7 and 10)
- Goal 2 → Baseline text, Doctor checks (the export rule and the `baseline-export` check, API Contracts 7 and 11)
- Goal 3 → ADR header and the lesson form, Documentation part of the update, Doctor checks (the ADR clause, `WithADRAmendedBy` and the `baseline-docs` check, API Contracts 1, 10 and 12)
- Goal 4 → Baseline text, Skills, guides and glossary (the lesson rule and the lesson steps of the skills, API Contracts 5, 14 and 15)
- Goal 5 → Baseline text, Doctor checks (the directory clause, `DocsLayoutDepartures` and the `baseline-docs` check, API Contracts 3 and 12)
- Story 1 → export list planner, export rule
- Story 2 → documentation part of the update, Managed Refresh of the Secondbrain guide
- Story 3 → `ExportVariationDepartures`, `clause.secondbrain.export-02-variation`
- Story 4 → `ExportListPlan.StopsExporting`, text and JSON of the documentation part, step 4 of the `update` Skill Mode (API Contracts 10 and 16)
- Story 5 → `clause.secondbrain.export-03-knowledge-mark`, `ReadADRHeader`
- Story 6 → the marking paragraph of the `update` Skill Mode (API Contract 16)
- Story 7 → `ADRAmendmentNotices`, `WithADRAmendedBy`, ADR clause
- Story 8 → lesson rule, `docs/lessons/` in the Standard Export List
- Story 9 → `clause.context.lesson-02-test`, lesson steps of the skills
- Story 10 → directory clause
- Story 11 → `baseline-export` and `baseline-docs` checks
- Story 12 → Release Plan rule (API Contract 13)
- Story 13 → the `## Decisions` line of the PRD on the notice; no component, the operator's session sends it
- Core Feature 1 → export list planner, `clause.secondbrain.export-01-standard-list`, documentation part of the update
- Core Feature 2 → process paths of `ExportListFacts`, invariant 5, `DR-EXPORT-PROCESS-PATH`
- Core Feature 3 → `ExportVariationDepartures`, `clause.secondbrain.export-02-variation`
- Core Feature 4 → conversion in `PlanExportList`, `Converted` and `StopsExporting`
- Core Feature 5 → invariant 4, the refusals of API Contract 10, `DR-EXPORT-MARKER`
- Core Feature 6 → `clause.secondbrain.export-03-knowledge-mark`, `ReadADRHeader`, `DR-EXPORT-ADR-UNMARKED`
- Core Feature 7 → the marking paragraph of the `update` Skill Mode, `documentation.adrs`
- Core Feature 8 → ADR clause, `ADRAmendmentNotices`, `WithADRAmendedBy`, header fills of the documentation part
- Core Feature 9 → lifecycle exception (API Contract 2)
- Core Feature 10 → `DR-DOCS-ADR-RETIRED` through `ClassifyADR`; the history cleanup is unchanged
- Core Feature 11 → `clause.context.lesson-01-record`, `LessonLearnedProblems`
- Core Feature 12 → `clause.context.lesson-02-test`
- Core Feature 13 → `clause.context.lesson-02-test`, the sections of API Contracts 14 and 15
- Core Feature 14 → message clause, retirement clause (API Contracts 4 and 6)
- Core Feature 15 → directory clause, `DocsLayoutDepartures`
- Core Feature 16 → documentation part of the update, its `fileChanges` entries, steps 3 and 4 of the `update` Skill Mode
- Core Feature 17 → `baseline-export` and `baseline-docs` checks
- Core Feature 18 → Release Plan rule, release sentence (API Contracts 13 and 19)
- Core Feature 19 → API Contracts 22 to 25
- Core Feature 20 → API Contracts 14 to 21; ADR-0438, ADR-0439 and ADR-0440 arrive with the planning change
- Success Metric 1 → Baseline text, Managed Refresh, `internal/baseline/documentation_rules_clauses_test.go`
- Success Metric 2 → `PlanExportList`, documentation part of the update and its tests
- Success Metric 3 → `ADRAmendmentNotices`, `WithADRAmendedBy` and their tests
- Success Metric 4 → the two Doctor checks and `internal/cli/doctor_documentation_test.go`
- Success Metric 5 → Release Plan rule and its test
- Success Metric 6 → API Contracts 22 to 24, read through the Doctor Command on this repository
- Success Metric 7 → `skills/documentation_rules_steps_test.go`
- Success Metric 8 → API Contracts 20 and 21, the glossary and documentation contracts
- Success Metric 9 → API Contracts 8 and 9, the three updated tests
- Success Metric 10 → declared under Unreachable Acceptance; the export list planner and API Contract 22 are what the measurement reads
- Success Metric 11 → declared under Unreachable Acceptance; the lesson rule and the lesson steps

Deviations from the PRD: None.

## Absence Evidence

1. `rule.context.lessons` and its two clauses — command: `rg -n -i 'lesson' internal/baseline/assets/modules`; result: one match, the sentence of `clause.core.message-the-owning-project` that API Contract 6 replaces.
2. `rule.secondbrain.knowledge-export` and its three clauses — command: `rg -n 'secondbrain-export|knowledge:|standard-export' internal/baseline/assets/modules internal/baseline/assets/templates`; result: no match.
3. `ReadADRHeader`, `ADRAmendmentNotices` and `WithADRAmendedBy` — command: `rg -n -g '!*_test.go' 'amended_by|AmendedBy|ADRHeader|knowledge:' internal cmd`; result: three matches, all the unexported `legacyADRHeader` of `internal/spec/retirement.go`, which reads a status from a body.
4. The export list planner — command: `rg -n -g '*.go' 'secondbrain-export' internal cmd`; result: five matches, all in `internal/docscontract/secondbrain_export_test.go`, the test API Contract 25 rewrites.
5. The documentation part of the update — command: `rg -n '"documentation"|Documentation:' internal/cli`; result: no match.
6. The two Doctor checks and their codes — command: `rg -n 'baseline-export|baseline-docs|DR-EXPORT|DR-DOCS' internal docs/user-guide .agents/skills`; result: no match.
7. `LessonLearnedProblems` and the lessons folder — command: `rg -n 'docs/lessons|LessonLearned' internal .agents/skills docs/agents docs/user-guide CONTEXT.md`; result: no match.
8. The lesson sections of the skills and the user-guide subsection — command: `rg -n 'Lesson Learned' .agents/skills docs/user-guide CONTEXT.md`; result: no match.
9. The five glossary entries — command: `rg -n '^\*\*(Standard Export List|Export Variation|Knowledge Mark|Amendment List|Lesson Learned)\*\*' CONTEXT.md`; result: no match.
10. The Release Plan rule and the release sentence — command: `rg -n 'docs/lessons|secondbrain-export' internal/releaseplan internal/cli/cli.go docs/user-guide/release-runbook.md docs/user-guide/usage.md`; result: no match.

## ADR Alternatives

- ADR-0438: implements it. Of the options it rejected, the design takes none: the list is at the source, in a block of its own, in the one file the sync reads, outside the Setup Manifest, with every standard folder listed.
- ADR-0439: implements it. No command decides a mark, no status is added, and the notice stays in the body beside the header field.
- ADR-0440: implements it. No judge, no approval before the file exists, no field in the Archive Record, and the QA gate writes no lesson.
- ADR-0278: obeys the retirement rule and changes its paragraph on lessons through the notice of API Contract 23. It rejected keeping retired files; a Lesson Learned is deleted the same way once a rule states it.
- ADR-0248: changes its sentence on the export through the notice of API Contract 23. Its rejected alternatives concern the history conversion, which this design does not touch.
- ADR-0412: extends its list of paths that set no release impact with two more, through the notice of API Contract 23.
- ADR-0282: obeys it; the history unit still removes a superseded ADR and names it in its successor, and the header fill skips a file that unit removed.
- ADR-0254: obeys it and repeats its shape: a part of the update planned in the command, bound into the Plan Digest and applied on the one approval. It rejected a separate operator procedure, and so does this design.
- ADR-0122: obeys it; a proposed ADR is not retired, so `DR-DOCS-ADR-RETIRED` reports only what `ClassifyADR` calls retired.
- ADR-0220: obeys it; no new check reads a definite fact that stops delivery, so each one warns and none fails.
- ADR-0411: obeys it and repeats its shape: a departure is named with where it belongs, and nothing is moved or rewritten.
- ADR-0186: obeys it; the clause texts name commands and paths every adopter has and no record of this repository.
- ADR-0222: obeys it; the rule of the Knowledge Mark lives in the guide of the module it depends on, so a repository without the Secondbrain never reads it.
- ADR-0257: obeys it; the test of a lesson is stated once, in the docs-layout guide, and the message clause and the skills point at it.
- ADR-0203: obeys it; the evidence is another project's request and the count of fifteen mirrors.
- ADR-0256: obeys it; each changed command changes the reference the Skill Coverage Map gives it, in the same delivery.
- ADR-0349: obeys it; the authorization record carries `release` and the ceiling `minor`.
- ADR-0417: uses the declaration of a Correction as the moment of a lesson and changes nothing in it.
- ADR-0420: uses the Corrections Report as one of the five signs and changes nothing in the record.
- ADR-0145: not applicable to the export list, which is no Managed Region of an instruction carrier; the design borrows its principle that a region is found by its markers and refuses a broken pair instead of refreshing it.
- ADR-0287: not applicable to the documentation part, which stays out of the guidance plan so that the settlement's check of the Setup Manifest never fails for a line of documentation.
- ADR-0100: not applicable to the documentation part; the preservation proof it describes still covers the instruction carriers, and the export file and the ADRs are not carriers.
- ADR-0330: not applicable beyond the one sentence API Contract 6 replaces; the Secondbrain stays read-only and messages still go through ai-memory.
- ADR-0396: not applicable to any component; the Notice Request of the adopted message is answered by the operator's session after the release.
- ADR-0092: not applicable — it keeps backlog and findings apart as intent and evidence, and this design renames and merges neither.
- ADR-0120: not applicable — the History Root is used as it is; a Lesson Learned is deleted and never moved there.
- ADR-0247: not applicable — the archive and its Archive Record are unchanged; the lesson step runs before the archive and adds no field.
- ADR-0251: not applicable — no Legacy Archive Folder is read or written.
- ADR-0258: not applicable — no Reduced History Entry is written, and the QA rerun classes bind the Task Graph, which a later stage writes.
- ADR-0260: not applicable — the Archive Relink and the list of broken citations are unchanged, and the documentation part removes no file a citation names.
- ADR-0261: not applicable — the pull request field of the Archive Record is untouched.
- ADR-0262: not applicable — the Cancel Command and the Source Return are unchanged.
- ADR-0277: not applicable — no Legacy Archive Folder is converted.
- ADR-0281: not applicable — no test rebuilds an executable.
- ADR-0290: not applicable — Item Binaries and Run Worktrees are untouched.
- ADR-0295: not applicable — no reachability read is added.
- ADR-0299: not applicable — the history conversion and its citation rewrite are unchanged.
- ADR-0308: not applicable — the archive's link rebase is untouched.
- ADR-0342: not applicable — settlement is unchanged, and the documentation part stays out of what settlement checks.
- ADR-0352: not applicable — the pre-PR review is unchanged; the size of the last step's diff is a risk this document names.
- ADR-0366: not applicable — no rule of the review delta changes.
- ADR-0368: not applicable — no authored data file is added.
- ADR-0391: not applicable — no rule of the review delta changes.
- ADR-0421: not applicable — the identity of the corrections record is untouched.
- ADR-0422: not applicable — a lesson is written in the commit that corrects, so the design adds no commit without a mark to an Item Branch.
- ADR-0423: not applicable — the park for an undeclared Correction Sign is unchanged.
- ADR-0424: not applicable — the corrections commands are used as they are.
- ADR-0426: not applicable — the Corrections Report is read and not changed.
- ADR-0432: not applicable — the Claude Instruction Link is untouched.

## Integration Points

- The Secondbrain's sync, in another project: it reads `.secondbrain-export` after each push to the default branch. The boundary is the file. The design relies on three facts read in its script: one path per line, a `!` line is an rsync exclusion, and a listed path that does not exist is skipped with a warning.
- ai-memory: no component calls it. The operator's session sends the completion notice and tells the project that owns the Secondbrain.
- The Vendored Skills `domain-modeling` and `knowledge-workspace`: not edited. The header fields are the Baseline's overlay on the ADR format, and `knowledge-workspace` already tells a session to add a new exported artifact to the file, which the Export Variation keeps true.
- The Update Check: it previews the update's JSON and counts `fileChanges`, so it reports a pending documentation part with no change of its own.

## Testing Approach

No new seam. Every test uses temporary repositories, the fake runners the packages have, and no network.

- `internal/spec/adr_header_test.go`: table tests of the reader, of each notice form with one and two identifiers, plain and linked, of a notice in a fenced block and in inline code, and of the writer's byte invariants 6 and 7, CRLF files included.
- `internal/spec/lesson_test.go`: one valid lesson and one case per problem of API Contract 12.
- `internal/baseline/export_list_test.go`: creation, conversion of the template's list with its losses, replacement between the markers, each refusal, invariants 1 to 5 and 9.
- `internal/baseline/docs_layout_test.go`: each named and unnamed folder and root document, and the path the repository's own rules guide names.
- `internal/baseline/documentation_rules_clauses_test.go`: one phrase per guarantee of API Contracts 1 to 7, in the module and in the rendered golden, and that no profile requires the two new rules.
- `internal/cli/baseline_update_documentation_test.go`: the plan, the digest that moves only when the part is pending, the confirmed apply, the second run, the Secondbrain turned off, the refusal that lets the rest apply, the header fill beside a history unit that appends to the same successor, and every line of the text section of API Contract 10. It drives `runBaselineUpdateCommand`.
- `internal/cli/doctor_documentation_test.go`: each of the eleven codes, both `ok` details, both `skipped` details, the cap of 20 lines, and that the exit code is the one the other checks give.
- `internal/releaseplan/classifier_test.go`: the two paths, alone and with another path.
- `skills/documentation_rules_steps_test.go`: one phrase per guarantee of API Contracts 14 to 19, in both copies.
- Existing tests that pin text this Spec changes, each updated in the Task that changes the text: `TestTheCaptureTriggerAsksForAConditionalLesson` takes the phrases "send no message about what it taught" and "Lesson Learned" in place of the three about the message; `docsOneJobPerDirectoryGuidance` in `internal/baseline/repository_rule_placement_test.go` takes the text of API Contract 3; `TestTheDocsLayoutRegionStaysLean` takes 13,000. `TestADRLifecycleContract` in `internal/baseline/plan_test.go` is not changed: API Contract 2 keeps the sentence it pins. The two tests of `internal/docscontract/secondbrain_export_test.go` are replaced, so the coverage record is recorded again. A test that pins the release sentence or a line of the update's or the Doctor Command's output on the default branch when a Task runs is updated by that Task.
- The glossary contract: the prescribed texts were applied to the prototype, and `TestGlossaryPhrasesAreDefinedTerms` passed with the five entries of API Contract 21.

## Build Order

1. The Baseline text: API Contracts 1 to 9, the three updated tests, the new clause test, the module record command, `make baseline-digests` and the Managed Refresh; and the glossary entries of API Contract 21, because the rendered guides and the later steps use the five terms.
2. The ADR header functions, the export list planner, the documentation part of the update, the Release Plan rule, and their descriptions: API Contracts 10, 13, 17 and 19, the `baseline.md` subsection of API Contract 20 and the release sentence in the two guides, with the Roundfix skill's version and the Behavior Surface Record (depends on: 1, because the part reads the module the Setup Manifest names and its texts use the glossary terms, and because both steps move the derived Baseline records and the skill version record). Its last step runs this repository's Managed Refresh once the part exists, which converts this repository's export file and fills the Amendment Lists its ADRs' notices call for, so that this repository's own update stays current and the two tests that expect it pass (`_prd.md` → Decisions, Tasks stage).
3. The lesson form check, the layout departures, the two Doctor checks, and their descriptions: API Contracts 11, 12 and 18 and the `doctor.md` part of API Contract 20, with the Roundfix skill's version and the Behavior Surface Record (depends on: 2, because the export check calls the planner and both steps edit the Roundfix skill's version fields and the Behavior Surface Record, and Tasks that share a file run in series).
4. The steps of the skills and the user guide: API Contracts 14 to 16, the `context-driven-development.md` subsection of API Contract 20, the skill contract test, with both skills' versions and the Behavior Surface Record (depends on: 3, for the shared skill version record and Behavior Surface Record; the `modes` text names the JSON member of step 2).
5. This repository: API Contracts 22 to 25, by marking every ADR, running the Managed Refresh again so the update writes the Standard Export List with the ADRs marked true and the Amendment Lists of the three notices, writing the three notices and the own lines of the export file, and recording the coverage record again (depends on: 2 and 3, because the update writes the list and the fills and the Doctor Command is how the step proves itself; and on 4 only for the order of the shared records).
6. The QA gate (depends on: 1, 2, 3, 4 and 5).

Files shared with other Specs. Specs 0411, 0414 and 0409 are delivered before this one and change `CONTEXT.md`, `docs/references/behavior-surfaces.json`, the Roundfix skill's `SKILL.md` in both copies, its `modes`, `setup` and `release` references, `skills/testdata/owned-skill-versions.json`, the module version record, the derived Baseline files, the closed list of clause forces, `docs/user-guide/context-driven-development.md`, `docs/user-guide/release-runbook.md`, `docs/user-guide/usage.md` and the `release plan` help. Each step of this Spec raises a version with its record command on the tree it runs on, inserts its sentences by the anchors this document names, and records the Behavior Surface Record after its last guide edit. No Verification of this Spec depends on another Spec, and the Task Graph needs no `requires`.

## Risks & Considerations

- Step 5 changes about 370 files by a few lines each. The diff stays under the bound of the pre-PR review because each hunk is one or two lines, and the Task must not reflow or re-save a file: `WithADRAmendedBy` and a line insertion for the mark, never a YAML round trip.
- The marks of step 5 are a judgment the Verification cannot prove. It proves that every ADR has a mark and that no other line changed; the Pull Request shows the marks to the maintainer.
- A plan merged after step 5 runs brings an ADR without a mark, and the Doctor Command then warns in this repository. That is the designed signal, and no contract fails.
- The exclusion `!docs/adr/*` was measured with the rsync of the author's machine, not with the one the sync's workflow runs. If it differs there, this repository's mirror keeps its ADRs until its own lines are corrected; no adopter is affected, and the PRD declares the row unreachable.
- The update applies the guidance, then the history, then the documentation part. A crash between them leaves a repository that a second run completes, because each part plans from the tree.
- The docs-layout region sits 52 bytes under its new ceiling. The next Spec that adds a sentence to that guide meets the ceiling, which is what a ceiling is for.
- A repository that upgrades Roundfix and does not apply the Baseline update sees the two checks warn about rules its guides do not state yet. Each warning names `roundfix baseline update`.
- The Secondbrain's sync prints one warning for each standard folder a repository lacks. It is another project's log line, and the operator tells that project with the plan.
- The conversion of the first update leaves a comment that introduced only removed lines where it was. It is harmless to the sync and a person may delete it.

## Glossary

None.

## Decisions

- The documentation part is planned in the command, beside the history plan, and stays out of the guidance plan and the Setup Manifest. See ADR-0438.
- Existing clauses of the docs-layout rule are reworded under their ids, and the lesson and the export list are new rules that no profile requires, so no Source Baseline row is written by hand. The prototype measured both.
- The header is edited by lines and never through a YAML encoder: no Go YAML library promises the same bytes back, and about 370 files must change by one line each.
- The planner and the checks are pure functions that take the repository's facts as values. The command gathers the facts once and gives them to the update and to the Doctor Command, so the two cannot disagree about what the list should hold.
- The apply fills a header over the bytes the file has at that moment and plans the list again from the tree, so the history unit that appends a line to a successor and the header fill of the same file do not overwrite each other.
- The documentation part's files are also `fileChanges` entries, so the Update Check and the `update` Skill Mode need no new field to count and show them.
- `--no-history` does not skip the documentation part, and no flag is added. This repository's Managed Refresh command therefore keeps its export list and its Amendment Lists current whenever a later Spec runs it.
- The marks and the lists move no `updated_at`.
- The Release Plan rule is general and reads no configuration, so it needs nothing from Spec 0414 and holds in this repository too.
- A code prints at most 20 lines, so a repository with hundreds of departures still gets a report a person reads.
- No Surface Transcript is fixed, for the reasons its section gives; every new line is prescribed in an API Contract and asserted by a named test.
- This stage adds no ADR.
