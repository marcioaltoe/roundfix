---
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# Claims with receipts and contracts as they ship

On 2026-09-29 and 2026-09-30 the planning changes of Waves 4, 5 and 6 (#276,
#280 and #286) each needed two to seven pre-PR review rounds. Most findings
were of two kinds:

- A PRD, TechSpec or Task contradicted another about what a command prints or
  returns.
- A claim about what a decision record or a file says could be settled only by
  opening that record.

Every Spec delivered in those waves then needed one to three corrective Tasks.
Several had the same cause: the Task author, the QA gate and the reviewer had
each rebuilt the expected output from prose.

This Spec gives both kinds of statement a form the Spec Consistency Check can
prove:

- A claim that attributes behavior to a decision record carries a **Claim
  Receipt**: a verbatim quote that the check finds in the record.
- A TechSpec states each command-line surface it changes as a **Surface
  Transcript**: the command, its output and its exit code, written as they
  will ship, and named by Tasks and by the QA gate.

## Project Constraints

- Identifier strategy: applicable — five `SC-*` diagnostic codes are added, and
  they are stable once shipped. A Surface Transcript is addressed by its name
  and number, like the other declared promises. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential is read, no network call is added and no model is called.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0183 and ADR-0184 (this Spec)
  record the two rules.
  ADR-0116 says what the citation check does today (ADR-0116: "reads the cited
  record when an artifact makes a claim about what that record establishes"),
  and its word comparison stays unchanged.
  ADR-0093 keeps the check on written text (ADR-0093: "it never judges whether
  a decision is correct"), so a receipt is proved by presence and never by
  meaning.
  ADR-0094 requires a detector to skip an absent artifact (ADR-0094: "Each
  detector declares the artifacts it reads and is skipped"), which is how a
  repository without the guide is read.
  ADR-0168 keeps a revised PRD on the horizon of its first commit
  (ADR-0168: "A PRD revised after a later ADR landed keeps the horizon of its
  first commit"), and the guide horizon reuses that rule.
  ADR-0156 makes a declared promise a coverage unit (ADR-0156: "Each numbered
  item is a coverage unit that some Task names in its References"), and a
  Surface Transcript is one more such promise.
  ADR-0117 places each check with its authoring stage (ADR-0117: "places each
  check with the stage that can produce the defect it catches").
  ADR-0176 reads citations only from authored text, and ADR-0130 holds the
  governed set to the record; neither changes.
  This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0155, ADR-0166 and ADR-0167. All hold.
  ADR-0097 cites ADR-0080 but carries a QA row forward on unmoved evidence.
  ADR-0179 cites ADR-0130 but decides what an empty grant authorizes.
  ADR-0182 (Spec 0190) cites ADR-0096 and ADR-0117 but moves mechanical facts
  to Task settlement. This Spec changes none of the three, so they do not
  apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for every Roundfix-owned skill ("considere autorizado a ajustar
  todas as skills se necessário"), and the standing grant of 2026-09-21 for
  the governed Go sources a slice needs, recorded in
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

## Goals

1. A false or unverifiable attribution to a decision record is found when the
   artifact is written, by a file read and not by a reviewer.
2. The expected output of a changed command exists once, in the TechSpec, and
   the Task tests and the QA gate name it.
3. A Spec written before these rules existed keeps passing every check it
   passed.
4. Authors learn both forms from the authoring skills and their templates, at
   a cost of a few lines per Spec.

## Core Features

1. **A written Claim Receipt is proved.** A receipt is a source, a colon and a
   quote in straight double quotes. The source is a decision record written as
   its `ADR-` identifier, or a backticked repository path. The check finds the
   quote in the source after it collapses whitespace in both. A quote that is
   not there, a source that does not resolve, and a quote shorter than three
   words are reported as the error `SC-RECEIPT-UNPROVEN`. Every receipt in a
   PRD or TechSpec is proved, whatever the Spec's age.
2. **An attribution without a receipt is a gap.** For a Spec held to the
   contract, each claim the citation check already reads as an attribution to
   a decision record needs a receipt for that record in the same paragraph.
   A claim without one is reported as the gap `SC-RECEIPT-MISSING`.
3. **A TechSpec declares its Surface Transcripts.** A TechSpec held to the
   contract declares numbered Surface Transcripts, or `None.` with the reason.
   A missing declaration is the gap `SC-TRANSCRIPT-UNDECLARED`. Each declared
   transcript holds one block with the command line, the standard output, the
   standard error and the exit code. A block with another shape is the error
   `SC-TRANSCRIPT-MALFORMED`, whatever the Spec's age.
4. **A Surface Transcript is traced.** Each declared transcript is named in
   the References of a Task that is not the QA gate, or `SC-COVERAGE-UNTASKED`
   reports it. When the Task Graph includes the gate, the QA Task names the
   transcript in a Requirement, or the error `SC-TRANSCRIPT-UNGATED` reports
   it. A QA Task that is already completed is historical evidence and is not
   read, and a declined gate raises nothing.
5. **The contract starts at its guide.** The guide is the concrete-contract
   guide of the repository's write-techspec skill. A repository where that
   guide is absent, or present and not yet committed, holds no Spec. Where the
   guide is committed, a Spec is held when its PRD was committed in or after
   the commit that added the guide, when its PRD is not committed yet, and
   when history cannot be read. A Spec whose PRD was committed before the
   guide is not held. For every Spec that is not held, the report names each
   skipped detector and the reason.
6. **The authoring skills teach both forms.** The write-techspec skill ships
   the concrete-contract guide and a template with the Surface Transcripts
   section. The write-prd skill and its template ask for a receipt on every
   attribution. The write-tasks skill asks each transcript to be named by its
   implementing Task and by the gate. The qa-gate skill asks the gate to
   reproduce each transcript through the built product. The guide also asks
   for interfaces written as code signatures and for numbered invariants.

## Non-Goals / Out of Scope

- Judging whether a quote supports its claim. No model is called. The check
  only exposes the parsed receipts so that a later advisory judgment can read
  them.
- Changing which sentences count as an attribution, or changing
  `SC-CITATION-UNSUPPORTED`.
- Requiring a receipt for a claim about a file. A file receipt is proved when
  written and never demanded.
- Running a transcript's command inside the Spec Consistency Check.
- Checking code signatures or numbered invariants mechanically.
- Rewriting an archived Spec, or a Spec that predates the guide.
- A new command, a new flag, or a change to the `roundfix-speccheck/v1`
  schema.

## Success Metrics

1. In a repository whose guide is committed, a PRD committed after it with an
   attribution that has no receipt reports `SC-RECEIPT-MISSING` as a gap, and
   as an error under `--strict`. Adding the verbatim receipt clears it.
   Changing one word of the quote reports `SC-RECEIPT-UNPROVEN`.
2. In a held Spec, a TechSpec with no Surface Transcripts declaration
   reports `SC-TRANSCRIPT-UNDECLARED`. A declared transcript with
   no exit line reports `SC-TRANSCRIPT-MALFORMED`. A transcript that no non-QA
   Task references reports `SC-COVERAGE-UNTASKED`, and one the QA Task does
   not name reports `SC-TRANSCRIPT-UNGATED`.
3. Every artifact with no receipt and no Surface Transcripts section is read
   exactly as before: the claims recorded by the characterization are
   unchanged, the corpus golden holds the five new codes at `0`, and every
   earlier count in it is unchanged.
4. This Spec's own PRD and TechSpec satisfy the two rules they introduce:
   every attribution carries a proved receipt, and every Surface Transcript is
   well-formed, referenced by a non-QA Task and named by the gate.

## Recorded limits

- A receipt proves that the quote exists in the source. It does not prove that
  the quote supports the claim.
- A receipt is required only for the attribution verbs the citation check
  already recognizes. A claim phrased with another verb is not read, as today.
- A planning branch cut before the guide landed and squash-merged after it is
  held once merged, because the merge commit adds its PRD after the guide. Its
  author rebases and adds the receipts and transcripts before merging.
- The check verifies a transcript's shape and its references. Only the QA gate
  runs the command and compares the output.
- A quote cannot contain a straight double quote, and a path source cannot
  contain a backtick or a space.

## Decisions

- **Receipts are proved, never judged.** See ADR-0183.
- **A command surface is a transcript.** See ADR-0184.
- **The horizon is the guide's commit, not a marker in the Spec.** A marker
  that a Spec declares can be left out, and the omission would switch the rule
  off silently. A date in the frontmatter is typed by the author and ties on
  the day the guide lands. The commit that added the guide is a fact the
  repository already holds, and it is absent exactly where the skill that
  teaches the rule is absent.
- **A written receipt is always proved.** Writing one asserts that the quote
  exists, so the assertion is checked even in a Spec that is not held.
- **File receipts are optional.** A path in prose is usually a pointer, not a
  claim, and no grammar separates the two reliably.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The outside-evidence rows rest on sources this Spec did not
produce:

- Published guidance on grounding a claim in a quote: the "Reduce
  hallucinations" page of the Claude documentation
  (<https://docs.anthropic.com/en/docs/test-and-evaluate/strengthen-guardrails/reduce-hallucinations>)
  tells the reader to verify each claim by finding a supporting quote and to
  retract a claim that has none.
- Published prior art for transcripts: the Cram documentation
  (<https://pypi.org/project/cram/>) describes tests of command-line
  applications written as shell sessions, with the command after `$`, its
  output on the following lines and a non-zero exit code in brackets.
- A replay on text this Spec did not write: the archived PRD and TechSpec of
  Spec 0182, copied into a disposable repository as a held Spec. On 2026-09-30
  the attribution pattern matched 6 sentences in those two files, and matched
  536 across the archived PRDs and TechSpecs, none of which carries a receipt.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the queries
`qmd query "citação literal verificada por código afirmação sobre ADR recibo"`
and
`qmd query "spec com contratos concretos assinaturas invariantes transcrição de CLI antes de implementar"`.
Two results informed the design:

- `inbox/secondbrain/2026-09-30-orquestradores-e-ecossistema-jev-o-que-se-transfere-ao-roundfix.md`
  records both mechanisms as plain-code ideas: a quote that code proves
  present, and concrete contracts in the Spec.
- `projects/fiscus/mirror/docs/adr/0065-a-guarda-de-citacao-recusa-o-rascunho-inteiro.md`
  refuses a whole draft on one unverifiable citation. It supports failing
  closed on an unproven receipt.

Exa located the Claude documentation page and the Cram documentation above,
and two independent open-source citation gates that first require a verbatim,
contiguous quote and only then ask whether it supports the claim. They confirm
the order this Spec takes: prove presence mechanically and leave support to a
reader or a later judgment.

## Technical candidate

The [_techspec.md](_techspec.md) records the forms, the invariants, the
transcripts, the coverage and the build order.
