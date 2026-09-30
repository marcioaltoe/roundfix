---
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
prd: _prd.md
created: 2026-09-30
---

# Claims with receipts and contracts as they ship — Technical Spec

## Executive Summary

Two detectors join the Spec Consistency Check, and one horizon decides which
Specs they bind.

- The receipt detector reads `source: "quote"` pairs from a PRD or TechSpec
  and proves each quote against its source file.
- The transcript detector reads a TechSpec's Surface Transcripts, checks each
  block's shape and traces it to a Task and to the QA gate.
- The horizon is the commit that added the concrete-contract guide to the
  write-techspec skill. A Spec whose PRD was committed before that commit is
  not held to the two declaration gaps.

Both detectors compare text that is already written. The trade-off accepted
is that a receipt proves presence and never support: a quote can exist and
still not establish the claim. That judgment stays with the reader. The
parsed receipts are exported so that a later advisory check can read the exact
quote, and this Spec builds none of that check.

## Project Constraints

- Identifier strategy: applicable — five stable `SC-*` codes, declared in the
  code block of `internal/speccheck/citations.go` where the other citation and
  coverage codes live. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential, no network call and no model call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0183, ADR-0184, ADR-0116, ADR-0093,
  ADR-0094, ADR-0168, ADR-0156, ADR-0117, ADR-0176 and ADR-0130 hold as the PRD
  states. The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0155, ADR-0166 and ADR-0167. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for every Roundfix-owned skill, and the standing grant of
  2026-09-21 for governed Go sources, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-techspec/references/concrete-contracts.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`,
  `skills/write-prd/SKILL.md`,
  `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`,
  `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`, `make baseline-digests`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package and no new command. Three files join `internal/speccheck`, and
the existing sweep calls them:

| Part | New file | Called from |
| --- | --- | --- |
| Contract horizon | `internal/speccheck/authoring_horizon.go` | both detectors |
| Receipts | `internal/speccheck/receipts.go` | `detectCitationCoverageAndReferences` in `citations.go`, and the stage-scoped check in `coherence.go` |
| Surface Transcripts | `internal/speccheck/transcripts.go` | the same two callers |

The full sweep (`Check`) already reaches `detectCitationCoverageAndReferences`,
so `internal/speccheck/constraints.go` does not change. The five codes are
declared in the code block of `citations.go` and join
`citationCoverageDetectorCodes`, so a Spec with no PRD lists them as skipped.
`coherence.go` registers each code with its stage and calls the detectors in
the stage-scoped check.

## Implementation Design

### Interfaces

```go
// internal/speccheck/authoring_horizon.go

// ConcreteContractGuidePath is the guide whose adding commit starts the
// receipt and transcript declaration gaps.
const ConcreteContractGuidePath = ".agents/skills/write-techspec/references/concrete-contracts.md"

type contractHorizon struct {
	held    bool
	missing string // what the skip line names when held is false
}

func newContractHorizon(repoRoot, prdPath string) contractHorizon
```

```go
// internal/speccheck/receipts.go

// Receipt is one source-and-quote pair an artifact wrote.
type Receipt struct {
	Artifact string // display path of the PRD or TechSpec
	Line     int    // line of the source token
	Source   string // "ADR-0116", or a repository-relative path
	Quote    string // the quote, whitespace-normalized
}

// ReceiptedClaim pairs an attribution with the receipts of its paragraph
// that name the same decision record.
type ReceiptedClaim struct {
	Claim    Claim
	Receipts []Receipt
}

func Receipts(artifact string, content []byte) []Receipt
func ReceiptedClaims(artifact string, content []byte) []ReceiptedClaim

// ProveReceipt resolves the receipt's source and reports whether the source
// holds the quote. reason is empty when proven is true.
func ProveReceipt(repoRoot string, receipt Receipt) (sourcePath string, proven bool, reason string, err error)
```

```go
// internal/speccheck/transcripts.go

// Transcript is one declared Surface Transcript.
type Transcript struct {
	Number    int
	Title     string
	Line      int      // line of the numbered item
	Command   string   // the text after "$ "
	Stdout    []string // lines between "stdout:" and "stderr:"
	Stderr    []string // lines between "stderr:" and "exit:"
	ExitCode  int
	Malformed string // "" or the reason the block is refused
}

func SurfaceTranscripts(content []byte) []Transcript
```

`Receipts`, `ReceiptedClaims`, `ProveReceipt` and `SurfaceTranscripts` are
exported. They are the seam for the guide's own test, for the QA gate and for
a later advisory judgment, which needs the claim sentence, the exact quote and
the source path and nothing else. `CitationClaims` keeps its signature and its
output.

### Data Models

No schema change. The `roundfix-speccheck/v1` JSON document carries the new
codes in its existing `findings` and `skipped` arrays.

### Invariants

1. For an artifact with no receipt and no Surface Transcripts section,
   `CitationClaims` and every earlier detector return what they returned
   before this Spec.
2. A receipt is proven exactly when its quote, after whitespace
   normalization, is a contiguous substring of its source file's text after
   the same normalization. The comparison is case-sensitive and reads no
   markup.
3. A receipt is proved whether or not the Spec is held. Only
   `SC-RECEIPT-MISSING` and `SC-TRANSCRIPT-UNDECLARED` depend on the horizon.
4. A receipt covers a claim only when both sit in the same paragraph and the
   receipt's source is the decision record the claim names.
5. Text inside a fenced code block is never read as a receipt, as a claim or
   as a numbered Surface Transcript item.
6. A Spec is not held only when the guide is absent, when the guide has no
   adding commit, or when Git proves that the guide's adding commit is not an
   ancestor of the PRD's adding commit. Every other case is held.
7. The detectors run no command and open no file outside the repository
   root.

### The contract horizon

`newContractHorizon(repoRoot, prdPath)` answers in this order:

1. The guide at `ConcreteContractGuidePath` is not a regular file under
   `repoRoot`: not held. `missing` is the guide path.
2. History cannot be read: held. This covers a directory that is not a Git
   work tree, a `repoRoot` that is not the work tree's root, a PRD in another
   repository, a shallow repository and any failing Git command. These are
   the conditions `newADRHorizon` already treats as unreadable.
3. The guide has no adding commit
   (`git log --diff-filter=A --format=%H -- <guide>` prints nothing): not
   held. `missing` is `a committed ` followed by the guide path.
4. The PRD has no adding commit: held. This is the Spec being written.
5. `git merge-base --is-ancestor <guide commit> <PRD commit>` exits `0`:
   held. It exits `1`: not held, and `missing` is
   `a PRD committed at or after ` followed by the guide path. Any other
   failure: held.

The guide commit is the oldest commit that added the guide. The PRD commit is
resolved the way `newADRHorizon` resolves it, and the two horizons share that
helper. A commit is its own ancestor, so a PRD added in the guide's commit is
held.

When a Spec is not held, the detector adds one skip for `SC-RECEIPT-MISSING`
and, when a TechSpec is present, one for `SC-TRANSCRIPT-UNDECLARED`. Both name
`missing`. The horizon is computed at most once per checked Spec.

### Claim Receipts

**Form.** A receipt is a source, a colon, whitespace and a quote in straight
double quotes. The whitespace may include one line break, so a receipt wraps
with its paragraph.

```text
ADR-0116: "reads the cited record when an artifact makes a claim"
`internal/speccheck/citations.go`: "CitationClaims parses subject attributions"
```

- A decision-record source is `ADR-` and four digits, at a word boundary.
- A path source is one backticked token with no whitespace that contains `/`
  or ends in a file extension. A backticked word such as a field name is
  therefore not a source.
- The quote is every character up to the next straight double quote. It
  cannot contain one.
- `Receipts` walks the same paragraphs `CitationClaims` walks: blank lines
  and a new `- ` item end a paragraph, and fenced blocks are skipped. The two
  functions share one paragraph walker.

**Proof.** `ProveReceipt` resolves the source, reads it and compares.

- `ADR-NNNN` resolves to the one file named `NNNN-*.md` under `docs/adr/`,
  then under `docs/history/adr/`. The whole file is read, whatever its status.
- A path resolves when it is clean, relative, stays under `repoRoot` and
  names a regular file.
- Normalization replaces every run of whitespace with one space and trims
  both ends. It is applied to the quote and to the source text.
- A quote with fewer than three words after normalization is refused before
  the source is read.

**Findings.**

| Code | Severity | Stage | Raised when |
| --- | --- | --- | --- |
| `SC-RECEIPT-UNPROVEN` | error | prd | a written receipt's source does not resolve, its quote has fewer than three words, or the source does not contain the quote |
| `SC-RECEIPT-MISSING` | gap | prd | the Spec is held and a resolved attribution has no receipt for its record in the same paragraph |

`SC-RECEIPT-MISSING` reads exactly the claims `SC-CITATION-UNSUPPORTED`
resolves: a claim whose decision record is not an accepted record in
`docs/adr/` needs no receipt, as today it gets no support check. Both codes
read the PRD at the `prd` stage and the PRD and TechSpec at later stages.

The three summaries of `SC-RECEIPT-UNPROVEN` end with
`does not contain that text`, `does not resolve to a file` or
`a receipt needs at least three words`. Surface Transcripts 1 and 2 give the
full text of the first kind and of the missing-receipt gap.

### Surface Transcripts in a TechSpec

**Declaration.** A TechSpec declares its transcripts in a section headed
`Surface Transcripts`, at heading level two or three. The section holds
numbered items that start with `Surface Transcript:`, or the single entry
`None.` followed by the reason, read by the rule `API Contracts` already uses.
Numbered lines inside a fenced block are not items.

**Block.** Each item carries exactly one fenced block whose info string is
`transcript`. After the block's indentation is removed, its lines are:

1. `$ ` followed by one command line;
2. the line `stdout:`, then zero or more lines of standard output;
3. the line `stderr:`, then zero or more lines of standard error;
4. the line `exit: ` followed by an integer from `0` to `255`, as the last
   line.

The first `stderr:` line after `stdout:` ends the standard output. Two
conventions keep a transcript exact without pinning what varies, and the
check does not interpret either:

- a line that is exactly `...` stands for any number of lines;
- text in angle brackets, such as `<line>`, stands for a value that differs
  between runs.

**Findings.**

| Code | Severity | Stage | Raised when |
| --- | --- | --- | --- |
| `SC-TRANSCRIPT-UNDECLARED` | gap | techspec | the Spec is held and its TechSpec has no Surface Transcripts declaration |
| `SC-TRANSCRIPT-MALFORMED` | error | techspec | a declared transcript has no block, more than one block, or a block that breaks the four line rules |
| `SC-COVERAGE-UNTASKED` | error | tasks | no Task other than the QA gate names the transcript in its References |
| `SC-TRANSCRIPT-UNGATED` | error | tasks | the graph includes the gate, the QA Task is not completed, and none of its Requirements names the transcript |

A Task names a transcript as `Surface Transcript 2` or
`Surface Transcripts 1-3`, parsed like the other coverage references. For a
transcript, the `SC-COVERAGE-UNTASKED` summary says that no Task other than
the QA gate references it. A
Surface Transcript has no Coverage Map entry, for the reason an API Contract
has none. A declined gate raises no `SC-TRANSCRIPT-UNGATED`. The malformed
reasons are `without a transcript block`, `with more than one transcript
block`, `without a command line`, `without a stdout line`, `without a stderr
line after its stdout line`, `without an exit line` and `with an exit code
outside 0-255`.

### API Contracts

1. API Contract: `roundfix spec check` findings — the command reports
   `SC-RECEIPT-UNPROVEN`, `SC-RECEIPT-MISSING`, `SC-TRANSCRIPT-UNDECLARED`,
   `SC-TRANSCRIPT-MALFORMED` and `SC-TRANSCRIPT-UNGATED` in its text and JSON
   output, with the severities and stages in the two tables above. `--strict`
   promotes the two gaps. Flags, exit codes and the JSON schema do not change.
2. API Contract: `roundfix spec check` skips — a Spec that is not held lists
   `SC-RECEIPT-MISSING` and `SC-TRANSCRIPT-UNDECLARED` under `Skipped:` with
   the reason the horizon gives.
3. API Contract: exported readers — `Receipts`, `ReceiptedClaims`,
   `ProveReceipt` and `SurfaceTranscripts` in `internal/speccheck` have the
   signatures in Interfaces and perform no network or model call.

### Surface Transcripts

Each transcript runs in a repository that holds the committed guide, except
transcript 4. `0200-example` is a Spec committed after the guide.

1. Surface Transcript: a receipt whose quote is not in its source. The PRD
   writes the quote with `records` where the decision record says `record`.

   ```transcript
   $ roundfix spec check 0200-example --stage prd
   stdout:
   Spec 0200-example
   [error] SC-RECEIPT-UNPROVEN: docs/specs/0200-example/_prd.md quotes ADR-0116 as "reads the cited records", but docs/adr/0116-a-citation-is-checked-against-what-it-cites.md does not contain that text
     at docs/specs/0200-example/_prd.md:<line>
     at docs/adr/0116-a-citation-is-checked-against-what-it-cites.md:1
     fix: Copy the passage verbatim from docs/adr/0116-a-citation-is-checked-against-what-it-cites.md into the receipt, or remove the receipt.
   ...
   stderr:
   exit: 1
   ```

2. Surface Transcript: an attribution with no receipt, in a held Spec and
   without `--strict`.

   ```transcript
   $ roundfix spec check 0200-example --stage prd
   stdout:
   Spec 0200-example
   [gap] SC-RECEIPT-MISSING: docs/specs/0200-example/_prd.md attributes "ADR-0116 requires the check to read the cited record." to ADR-0116 without a receipt
     at docs/specs/0200-example/_prd.md:<line>
     at docs/adr/0116-a-citation-is-checked-against-what-it-cites.md:<line>
     fix: Add ADR-0116: "<verbatim passage>" to the same paragraph in docs/specs/0200-example/_prd.md.
   ...
   stderr:
   exit: 0
   ```

3. Surface Transcript: a transcript with no exit line, and a transcript the
   QA Task does not name.

   ```transcript
   $ roundfix spec check 0200-example
   stdout:
   Spec 0200-example
   ...
   [error] SC-TRANSCRIPT-MALFORMED: docs/specs/0200-example/_techspec.md declares Surface Transcript 1 without an exit line
     at docs/specs/0200-example/_techspec.md:<line>
     fix: Write the block as a command line starting with "$ ", then "stdout:", "stderr:" and "exit: <code>", in that order.
   ...
   [error] SC-TRANSCRIPT-UNGATED: docs/specs/0200-example/_techspec.md declares Surface Transcript 2, but QA Task docs/specs/0200-example/task_02.md names it in no Requirement
     at docs/specs/0200-example/_techspec.md:<line>
     at docs/specs/0200-example/task_02.md:<line>
     fix: Name Surface Transcript 2 in a Requirement of docs/specs/0200-example/task_02.md so the gate reproduces it.
   ...
   stderr:
   exit: 1
   ```

4. Surface Transcript: a repository that does not hold the guide. The same
   Spec reports no declaration gap and names the guide in its skips.

   ```transcript
   $ roundfix spec check 0200-example --stage techspec
   stdout:
   Spec 0200-example
   No findings. Authored Verification commands were not executed.
   Skipped:
   ...
     SC-RECEIPT-MISSING: missing .agents/skills/write-techspec/references/concrete-contracts.md
   ...
     SC-TRANSCRIPT-UNDECLARED: missing .agents/skills/write-techspec/references/concrete-contracts.md
   ...
   stderr:
   exit: 0
   ```

### Skills, guides and glossary

Every skill edit is additive and sits under its own heading. No line is added
inside a `### QA settlement` section.

- **The guide**, `.agents/skills/write-techspec/references/concrete-contracts.md`:
  the receipt form, the transcript block and its two conventions, interfaces
  as code signatures and numbered invariants. Its worked example is written
  as a TechSpec fragment, with one receipt that quotes a file the skill ships
  and one transcript. Every other illustration sits in a fenced block. The
  file's adding commit is the contract horizon.
- **write-techspec**: a `## Concrete contracts` section that points to the
  guide, and a `### Surface Transcripts` section in the template after
  `### API Contracts`. The template names the guide and nests no fenced block.
- **write-prd**: a `## Claim receipts` section, and one comment in the
  template's Project Constraints block.
- **write-tasks**: a `## Surface Transcripts in Tasks` section: the
  implementing Task names each transcript in its References and asserts its
  text in a test.
- **qa-gate**: a `## Surface Transcripts at the gate` section: one row per
  transcript reproduces the command through the built product.
- **Roundfix skill**: the five codes join the identifier list of its Spec
  Consistency Check section, with the horizon rule.
- `CONTEXT.md` gains **Claim Receipt** and **Surface Transcript**, and the
  Spec Consistency Check entry names the contract horizon.
  `docs/user-guide/context-driven-development.md` gains one paragraph for
  each rule.

## Vocabulary Contract

- emits: `internal/speccheck/citations.go`
  pattern: `SC-(RECEIPT|TRANSCRIPT)-[A-Z]+`
  documented-in: `.agents/skills/roundfix/SKILL.md`

## Coverage Map

- Goal 1 → Claim Receipts; API Contract 1.
- Goal 2 → Surface Transcripts in a TechSpec; API Contract 1.
- Goal 3 → The contract horizon; Invariants 1, 3 and 6; API Contract 2.
- Goal 4 → Skills, guides and glossary.
- Core Feature 1 → Claim Receipts (proof).
- Core Feature 2 → Claim Receipts (findings); The contract horizon.
- Core Feature 3 → Surface Transcripts in a TechSpec (declaration and block).
- Core Feature 4 → Surface Transcripts in a TechSpec (findings).
- Core Feature 5 → The contract horizon.
- Core Feature 6 → Skills, guides and glossary.
- Success Metric 1 → Testing Approach 2, Testing Approach 3.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 1, Testing Approach 5.
- Success Metric 4 → Testing Approach 5, Testing Approach 7.

## Integration Points

- **Git.** The horizon runs read-only Git commands with the environment and
  the root checks `newADRHorizon` uses. No fetch and no write.
- **Stage-scoped check.** `coherence.go` registers the codes and calls the
  detectors, so `--stage prd` and `--stage techspec` report them.
- **QA gate precondition and Delivery Queue revalidation.** Both run the
  strict check, so a gap in a held Spec refuses there. A Spec that is not
  held is unaffected.
- **Later advisory judgment.** A later wave plans an advisory check of
  whether the quoted passage supports its claim. The supervising session
  reported, on 2026-09-30, a shadow benchmark in which such a check did better
  than the word comparison. That check will read `ReceiptedClaims` and
  `ProveReceipt`, which give it the claim sentence, the exact quote and the
  source path. It is not part of this Spec, and nothing here calls a model.

## Testing Approach

1. **Characterization first.** Before any parser changes,
   `internal/speccheck/receipt_characterization_test.go` records the claims
   `CitationClaims` returns for fixture copies of four archived artifacts
   (the PRD and TechSpec of Specs 0181 and 0182) in
   `internal/speccheck/testdata/receipt-characterization/`. The test compares
   every claim's artifact, line, target and subject with a recorded golden
   file. It passes on the unchanged parser and after the change.
2. **Horizon.** `internal/speccheck/authoring_horizon_test.go` builds
   temporary Git repositories with `gittest`:
   - a PRD committed before the guide is not held;
   - a PRD committed in the guide's commit, or after it, is held;
   - an uncommitted PRD is held once the guide is committed;
   - a repository without the guide, and one whose guide is uncommitted, hold
     nothing and name the reason;
   - a directory with the guide and no Git holds every Spec.
3. **Receipts.** `internal/speccheck/receipts_test.go` covers:
   - a verbatim quote, and a quote that wraps across lines in the artifact
     and in the source, are proven;
   - one changed word, an unresolved source and a two-word quote are
     unproven;
   - a file source is proved against the file;
   - a receipt inside a fenced block is not read;
   - a held attribution with no receipt is a gap, and a receipt for another
     record does not cover it;
   - Surface Transcripts 1 and 2: the rendered text holds each transcript's
     finding lines.
4. **Transcripts.** `internal/speccheck/transcripts_test.go` covers:
   - the declaration gap in a held Spec, and `None.` with a reason;
   - a well-formed block, and each malformed reason;
   - numbered output inside a block is not an item;
   - a transcript named only by the QA Task is untasked;
   - a transcript the QA Task does not name is ungated, and a declined gate
     is not;
   - Surface Transcripts 3 and 4: the rendered text holds each transcript's
     finding and skip lines.
5. **Corpus.** `internal/docscontract/corpus_test.go` lists the five codes,
   and the corpus golden holds each at `0` with every earlier count
   unchanged. The active corpus includes this Spec, so the sweep proves this
   Spec's own receipts and transcripts.
6. **Sabotage.** Every code has one test in which the defect is planted and
   the finding is reported, and one in which the defect is absent and nothing
   is reported. A detector that returns nothing fails the first.
7. **Guide and skills.** `skills/concrete_contract_guide_test.go` proves that
   the embedded bundle holds the guide at `ConcreteContractGuidePath`. It
   also proves that the guide's worked example is read as exactly one
   receipt, proved against a file the skill ships, and exactly one
   well-formed transcript.
   Phrase checks cover each skill heading, the glossary and the user guide,
   with `make skills-sync-check`.

## Build Order

1. Characterization, the contract horizon and Claim Receipts, with the two
   receipt codes in the stage table, the corpus golden and the Roundfix
   skill, task_01 (depends on: none).
2. Surface Transcripts, with the three transcript codes in the same places,
   task_02 (depends on: 1). Both edit `citations.go`, `coherence.go`, the
   corpus golden and the Roundfix skill.
3. The guide, the authoring skills and their templates, the glossary and the
   user guide, task_03 (depends on: 2). It describes what 1 and 2 ship, and
   its commit starts the horizon.
4. Terminal QA, task_04 (depends on: 1, 2, 3).

## Risks & Considerations

- **A receipt can be true and beside the point.** The check proves presence
  only. The summary quotes the claim and names the source, so a reader
  settles support in seconds, the property ADR-0116 already relies on.
- **The horizon binds planning branches that merge late.** A Spec authored
  before the guide and squash-merged after it is held once merged. The PRD
  records this limit. The authoring skills ask for both forms from the day
  the guide ships, which is the same commit.
- **A backticked path followed by a quoted value.** Prose such as a
  configuration key with a slash, a colon and a quoted value would be read as
  a file receipt. Measured on 2026-09-30, no PRD or TechSpec in the archive
  (343 files) or in `docs/specs/` holds the form.
- **Transcripts can rot.** The check proves shape and references. The gate
  reproduces each transcript, and a Task test asserts its text.
- **Cost.** One horizon per Spec costs at most four Git calls. Proving a
  receipt reads one file.
- **Skill edits from other Specs.** Other Wave 7 Specs edit the same skills.
  Every edit here is a new section under its own heading.

## Decisions

- The horizon is the guide's adding commit. See ADR-0183.
- A receipt is proved by contiguous, whitespace-normalized presence. See
  ADR-0183.
- A command surface is declared as a transcript, and the check never runs it.
  See ADR-0184.
- A transcript needs a Task other than the gate, because the QA Task's
  References list every promise and would hide a transcript nobody
  implements.
- The template points to the guide instead of nesting a transcript block
  inside its own fenced block.
