---
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
prd: _prd.md
created: 2026-11-02
---

# A glossary definition written as a list item counts — Technical Spec

## Executive Summary

The Spec Consistency Check reads glossary definitions through one pattern,
`glossaryDefinitionPattern` in `internal/speccheck/glossary.go`, which accepts
only a line starting with `**Term**:`. The fix widens that one pattern to the
Glossary Definition forms the maintainer chose: the bolded term at the start
of a line followed by a colon or by a spaced em dash, bare or after one
top-level `-`, `*` or `+` list marker. Because `glossaryFiles` feeds the same
term set to every glossary code, one change fixes `SC-GLOSSARY-MISSING`, the
changed-term branch of `SC-GLOSSARY-UNPLANNED` and the bold coverage of
`SC-GLOSSARY-UNDECLARED` together.

The same file gains `GlossaryDefinitionForms`, the text the
`SC-GLOSSARY-MISSING` fix and the guides use, and
`GlossaryVerificationCommand`, which renders the line-anchored glossary
Verification that `SC-GLOSSARY-UNPLANNED` prints and `write-tasks` asks for.
A test runs that command against every form the reader is tested on, so the
Task's Verification and the check cannot drift apart.

The trade-off accepted: the command is exact in case and whitespace while the
reader folds both, so a Verification can fail on a definition the check would
accept; it can never pass on one the check rejects.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0244 holds a Spec's declared terms
  to the glossary deterministically; the design keeps one deterministic
  pattern and changes only which lines define a term. ADR-0280 holds this
  repository's own glossary to its Repository Contract Test; that test and
  its `**Term**:` reader stay unchanged, and the new `CONTEXT.md` entry is
  written in that form. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AV answers of 2026-11-02 quoted in the authorization record; the
  Governed Paths are the Roundfix, `write-prd` and `write-tasks` skill files
  the Tasks change; no Makefile, workflow, module or lint configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0341-a-glossary-definition-written-as-a-list-item-counts/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/glossary.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`.

## Current behavior

Read on 2026-11-02 at `900e13f8`; `_prd.md` → Validation records the run.

1. `glossaryDefinitionPattern` is `^\*\*([^*]+)\*\*:`. `glossaryFiles`
   applies it to each line of every glossary file and stores the normalized
   capture (case folded, whitespace collapsed) in the defined term set.
2. `detectGlossary` reads that set three times: a bold candidate whose
   normalized phrase (or its `s`/`es` plural) is in the set is covered; a
   `changes:` term outside the set reports `SC-GLOSSARY-UNPLANNED` ("is not
   defined"); a declared term whose binding Tasks all completed and that is
   outside the set reports `SC-GLOSSARY-MISSING` with the fix "Define the
   declared term in the glossary.".
3. A term without a binding Task reports `SC-GLOSSARY-UNPLANNED` with the fix
   "Add a non-QA Task declaring the glossary under interface: or creates: and
   naming **<term>** in ## Verification." A Task binds when its Verification
   contains the exact bolded term.
4. `write-tasks` (`### Glossary requirement`) asks for a Verification that
   "checks the bolded term followed by a colon in the glossary on
   whitespace-normalized text". Normalizing whitespace joins the lines, so
   that Verification also passes on a `**<term>**:` inside a sentence, which
   the check rejects.
5. The spec command guide says "Definitions begin a line with `**Term**:`."
   The `write-prd` glossary guide and the Roundfix `spec` reference do not
   name a form.
6. This repository's own glossary contract
   (`parseGlossaryContract` in `internal/docscontract/user_guide_contract_test.go`)
   reads `**Term**:` lines with its own pattern; it is not a reader of the
   Spec Consistency Check and stays unchanged.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/speccheck/glossary.go`: the pattern, the two exported names, the
  two fix texts.
- `internal/speccheck/glossary_test.go`: the form table, the agreement test,
  the message test.
- `skills/glossary_definition_test.go` (new): the guides carry the forms and
  the command.
- The spec command guide, `CONTEXT.md`, and the `write-tasks`, `write-prd`
  and Roundfix skills.

## Implementation Design

### Interfaces

```go
// internal/speccheck/glossary.go

// A Glossary Definition: the bolded term at the start of a line, bare or
// after one "- ", "* " or "+ ", followed by ":" or by " —" and a space or
// the end of the line.
var glossaryDefinitionPattern = regexp.MustCompile(`^(?:[-*+] )?\*\*([^*]+)\*\*(?::| —(?: |$))`)

// GlossaryDefinitionForms names the accepted forms for messages and guides.
const GlossaryDefinitionForms = "a line that starts with `**Term**:`, `**Term** —`, `- **Term**:` or `- **Term** —`, where the list marker may also be `*` or `+`"

// GlossaryVerificationCommand renders the shell command that passes only
// when glossary holds a Glossary Definition of term, exact in case.
func GlossaryVerificationCommand(term, glossary string) string
```

### Data Models

No schema changes.

### API Contracts

1. API Contract: the reader. `glossaryDefinitionPattern` is the pattern
   above; `glossaryFiles` keeps its loop, its normalization and its file set.
   A line is a Glossary Definition exactly when the pattern matches it.
2. API Contract: the rendered command. `GlossaryVerificationCommand(term,
   glossary)` returns
   `grep -F -- '**<term>**' < <glossary> | grep -qE '^([-*+] )?\*\*<escaped>\*\*(:| —( |$))'`,
   where `<escaped>` is `regexp.QuoteMeta(term)` (whose backslash escapes
   are valid POSIX extended expressions) and each single-quoted argument
   writes an apostrophe as `'\''`. The first `grep` keeps the literal
   bolded term, so the command binds the Task under the unchanged binding
   rule, and reads the glossary on standard input, so `SC-VERIFY-WRAP-FRAGILE`
   (which exempts a `grep` reading standard input) does not report the
   command `write-tasks` asks for; the second anchors the definition at the
   start of the line. For
   `<Term>` and `CONTEXT.md` it renders
   `grep -F -- '**<Term>**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*<Term>\*\*(:| —( |$))'`.
3. API Contract: the messages. The `SC-GLOSSARY-MISSING` fix becomes
   `Define the declared term in the glossary as ` + `GlossaryDefinitionForms`
   + `.`; the `SC-GLOSSARY-UNPLANNED` fix for a term without a binding Task
   appends `, such as: ` + `GlossaryVerificationCommand(term, path)`, where
   `path` is `CONTEXT.md` when it is a glossary file and otherwise the
   lexically first glossary file. Summaries, codes, locations and every other
   fix text keep their bytes.
4. API Contract: the rules in writing.
   - `.agents/skills/write-tasks/SKILL.md`, `### Glossary requirement`: the
     Verification runs the command `roundfix spec check` prints for an
     unbound term, shown rendered for `<Term>` and `CONTEXT.md` as in API
     Contract 2, with the term's metacharacters escaped; it reads lines, not
     whitespace-normalized text, because a definition is a line; it is never
     more permissive than the check. The sentence about "the bolded term
     followed by a colon" is replaced.
   - `.agents/skills/write-prd/references/glossary.md` and
     `.agents/skills/roundfix/references/spec.md`: a Glossary Definition is
     `GlossaryDefinitionForms`, whitespace-normalized; headings, indented or
     ordered list items and bold inside a sentence define nothing.
   - `docs/user-guide/commands/spec.md`: "Definitions begin a line with
     `**Term**:`." becomes the accepted forms, and the `SC-GLOSSARY-MISSING`
     row names them.
   - `CONTEXT.md`: adds **Glossary Definition**.

### Surface Transcripts

None. The changed finding texts are asserted by named tests; no whole-command
transcript is declared.

### Invariants

```text
1. A line the reader accepted at 900e13f8 (a line starting with **Term**:) is still a Glossary Definition.
2. A Glossary Definition is exactly a line matching the pattern of API Contract 1; headings, indented or ordered items, ASCII or en dashes, an unspaced em dash and bold inside a sentence define nothing.
3. The rendered command exits 0 only on a glossary holding a line the reader accepts for that exact term.
4. The binding rule, the Glossary Declaration grammar, the bold candidate rule and the glossary file set keep their behavior.
```

### Declared breaks

- A glossary line in one of the new forms now defines its term, so a Spec
  that relied on such a term being undefined (a `changes:` declaration now
  defined, a bold candidate now covered) reports fewer findings. In this
  repository's `CONTEXT.md` and the local adopter glossaries no line changes
  classification except Pantheon's 16 em dash items (`_prd.md` →
  Validation). The `SC-GLOSSARY-MISSING` and unbound `SC-GLOSSARY-UNPLANNED`
  fix texts change; no test or golden pins them
  (`grep -rn 'Define the declared term\|in ## Verification' internal skills --include=*_test.go`
  finds none).

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `glossaryDefinitionPattern`
  (API Contract 1), Invariant 2.
- Goal 2, Core Feature 1, Success Metric 2 → the form table of
  `TestGlossaryDefinitionFormsAreRead`, Invariants 1, 2 and 4.
- Goal 3, Core Feature 2, Success Metric 3 → the fix texts (API Contracts 2
  and 3).
- Goal 4, Core Features 3 and 4, Success Metrics 4 and 5 →
  `GlossaryVerificationCommand`, the agreement test and the guide texts (API
  Contracts 2 and 4), Invariant 3.

Deviations from the PRD: None.

## Absence Evidence

1. `GlossaryDefinitionForms` and `GlossaryVerificationCommand` — command: `grep -rn 'GlossaryDefinitionForms\|GlossaryVerificationCommand\|DefinitionForms' internal skills`; result: no match.
2. A second definition reader in the check — command: `grep -rn 'glossaryDefinitionPattern' internal`; result: only its declaration and its use in `glossaryFiles`.
3. Form guidance in the authoring guides — command: `grep -n -- '- \*\*Term\*\*\|em dash' .agents/skills/write-prd/references/glossary.md .agents/skills/write-tasks/SKILL.md .agents/skills/roundfix/references/spec.md`; result: no match.
4. `skills/glossary_definition_test.go` — command: `ls skills/glossary_definition_test.go`; result: no such file.

## ADR Alternatives

- ADR-0244: kept as recorded; it rejected a model judging "is this a domain
  term" because the check must answer the same way on every run, and the
  widened reader is still one deterministic pattern.
- ADR-0280: kept as recorded; it rejected holding this repository's glossary
  through `roundfix spec check`, and the design leaves that contract to its
  own test, untouched.

## Integration Points

None external.

## Testing Approach

All tests are hermetic: temporary Git repositories through the existing
`glossaryRepo`, `glossaryTask` and `glossaryCheck` helpers, and no network.

- `internal/speccheck/glossary_test.go`:
  - a package-level table `glossaryDefinitionFormCases` (name, line,
    defined) holding, for a two-word term written as a fixture: accepted
    `**<t>**: d`, `- **<t>**: d`, `* **<t>**: d`, `+ **<t>**: d`,
    `**<t>** — d`, `- **<t>** — d`, `- **<t>** —`; rejected `### <t>`,
    `  - **<t>**: d`, a tab-indented `- **<t>** — d`, `1. **<t>**: d`,
    `- **<t>** - d`, `- **<t>** -- d`, `- **<t>** – d`, `- **<t>**—d`,
    `-**<t>**: d`, `- **<t>** is prose`, `- **<t>.** Prose`,
    `See **<t>**: d`.
  - `TestGlossaryDefinitionFormsAreRead`: for each row, a Spec adding the
    term with a completed binding Task reports `SC-GLOSSARY-MISSING` at the
    declaration line exactly when the row is rejected, and a PRD with
    `None.` that bolds the term reports `SC-GLOSSARY-UNDECLARED` exactly when
    the row is rejected.
  - `TestGlossaryVerificationCommandAgreesWithTheDefinitionReader`: resolves
    `grep` with `testfixture.RequireTool`; for each row, plus a term with a
    `.` whose line holds a different term in the definition position, a term
    with an apostrophe in an accepted form, and a line defining another term
    that cites the term later, it writes a temporary `CONTEXT.md`, checks
    the reader's acceptance against the row and runs
    `sh -c GlossaryVerificationCommand(term, "CONTEXT.md")` there; exit 0
    must equal the row's `defined`.
  - `TestGlossaryFindingsNameTheDefinitionForms`: the `SC-GLOSSARY-MISSING`
    fix contains `GlossaryDefinitionForms`; `WrapFragileVerification` reports
    nothing for a Task whose Verification is the rendered command; the unbound
    `SC-GLOSSARY-UNPLANNED` fix ends with
    `GlossaryVerificationCommand(term, "CONTEXT.md")`, and with
    `GlossaryVerificationCommand(term, "GLOSSARY.md")` when only `GLOSSARY.md`
    exists.
- `skills/glossary_definition_test.go`,
  `TestTheGlossaryGuidesStateTheDefinitionForms`: both copies of the
  `write-tasks` skill contain `GlossaryVerificationCommand("<Term>",
  "CONTEXT.md")`; both copies of the `write-tasks` skill, the `write-prd`
  glossary guide and the Roundfix `spec` reference contain
  `GlossaryDefinitionForms` on whitespace-normalized text.
- Existing suites: `./internal/speccheck`, `./skills`, and
  `-tags docscontract ./internal/docscontract` (the corpus, the glossary
  contract over the new `CONTEXT.md` entry and the Behavior Surface Record).
- Sabotages, each run by the author on a scratch copy at `900e13f8` with the
  prototype, with the case whose input makes it observable:
  - list marker dropped from the reader: the agreement test fails on
    `dash list colon`, `star list colon`, `plus list colon`, and the form
    test on the same rows;
  - leading whitespace allowed in the reader: both tests fail on
    `indented list item` and `tab indented list item`;
  - em dash dropped from the reader: both fail on `em dash` and
    `dash list em dash`;
  - the command also accepts ` - `: the agreement test fails on
    `ascii hyphen`;
  - `regexp.QuoteMeta` dropped from the command: the agreement test fails on
    `metacharacter term`;
  - list marker dropped from the command only: the agreement test fails on
    `dash list colon`;
  - the MISSING fix reverted: `TestGlossaryFindingsNameTheDefinitionForms`
    fails;
  - the command reads the glossary as a file operand instead of standard
    input: `TestGlossaryFindingsNameTheDefinitionForms` fails on the
    wrap-fragile assertion;
  - the `write-tasks` command line removed:
    `TestTheGlossaryGuidesStateTheDefinitionForms` fails.

## Build Order

1. The reader, the exported form text and command, the fix texts, their
   tests, the spec command guide and `CONTEXT.md`, with the Behavior Surface
   Record re-recorded last (task_01).
2. The `write-tasks`, `write-prd` and Roundfix skill texts, the skills test,
   with the mirrors synced and the owned-skill versions recorded last
   (task_02) (depends on: 1).
3. QA gate (task_03) (depends on: 2).

task_02 follows task_01 because its test reads the two exported names and the
spec command guide that task_01 changes is a source of the Roundfix `spec`
reference's Behavior Surface. No other Spec is active; no Verification depends
on another Spec.

## Risks & Considerations

- Pantheon's nine `- **Term** (_alias_) —` definitions stay unread; a Spec
  that declares one of those terms still needs the plain form. The guides
  name the accepted forms, so the author learns it at authoring time.
- The command is exact in case and whitespace, so a definition written with
  another case passes the check but fails the Task's Verification; the
  Verification is then fixed by writing the declared spelling.
- A glossary line inside a fenced block is read as before; the reader never
  skipped fences.

## Glossary

- adds: **Glossary Definition** — the glossary line shapes that define a term

## Decisions

- One reader pattern for every glossary code; no second reader.
- Accepted markers are the three CommonMark bullet markers; leading
  whitespace is refused because an indented item is a sub-point of another
  entry (a field or a note), not a vocabulary entry; ordered markers are
  refused because no glossary in the corpus numbers its terms.
- Separators are the colon (any text after it, as today) and the em dash
  with a space or the line end after it; ASCII hyphen, double hyphen and en
  dash are refused because they appear in prose and no glossary in the
  corpus uses them, and a later Spec can widen the set without breaking
  anyone.
- The glossary Verification is the rendered line-anchored command, printed
  by `SC-GLOSSARY-UNPLANNED` and taught by `write-tasks`, never more
  permissive than the reader.
- No ADR: the forms are a fact for `CONTEXT.md` and the guides.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding a form. Exa: the CommonMark list item
  section.
- Two implementation Tasks in series and a QA gate.
