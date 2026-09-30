---
task: task_03
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: pending
type: docs
complexity: medium
---

# Task 03: The authoring skills teach receipts and concrete contracts

## Overview

task_01 and task_02 made the Spec Consistency Check prove Claim Receipts and read Surface Transcripts. This Task teaches the two forms where Specs are written. It adds the concrete-contract guide to the write-techspec skill, whose adding commit is the contract horizon, and small sections in the write-techspec, write-prd, write-tasks and qa-gate skills. It also adds the two terms to the glossary and the user guide. Every edit is additive and sits under its own heading.

## Requirements

1. MUST create `.agents/skills/write-techspec/references/concrete-contracts.md`, the guide, with:
   - the Claim Receipt form and what the check proves;
   - the Surface Transcript block, its four line rules and its two conventions and how the gate matches them (a line `...` matches zero or more consecutive lines, text in angle brackets matches one or more characters within the same line, everything else and the exit code match exactly);
   - interfaces written as code signatures and numbered invariants;
   - a worked example written exactly as a TechSpec fragment: one paragraph that carries one receipt, and a `Surface Transcripts` section with one item and its block. The receipt's source MUST be a file that ships with the write-techspec skill, so that the example is provable wherever the skill is installed. Every other illustration of either form in the guide MUST sit inside a fenced block, where the check does not read it.
2. MUST add a `## Concrete contracts` section to `.agents/skills/write-techspec/SKILL.md` that points to the guide, and a `### Surface Transcripts` section to `.agents/skills/write-techspec/references/techspec-template.md` directly after `### API Contracts`. The template section MUST name the guide and MUST NOT nest a fenced block.
3. MUST add a `## Claim receipts` section to `.agents/skills/write-prd/SKILL.md`, and one comment that names the Claim Receipt in the Project Constraints block of `.agents/skills/write-prd/references/prd-template.md`.
4. MUST add a `## Surface Transcripts in Tasks` section to `.agents/skills/write-tasks/SKILL.md`: the implementing Task names each transcript in its References and asserts its text in a test, and the QA Task names it in a Requirement.
5. MUST add a `## Surface Transcripts at the gate` section to `.agents/skills/qa-gate/SKILL.md`: one row per transcript reproduces the command through the built product and compares standard output, standard error and the exit code by the matching rules of `_techspec.md` → Surface Transcripts in a TechSpec: a line `...` matches zero or more consecutive lines, text in angle brackets matches one or more characters within the same line, and everything else matches exactly. It MUST NOT add or change any line inside the `### QA settlement` section.
6. MUST add **Claim Receipt** and **Surface Transcript** to `CONTEXT.md` in the glossary's entry form, name the contract horizon in the Spec Consistency Check entry, and add one paragraph for each rule to `docs/user-guide/context-driven-development.md`.
7. MUST run `make skills-sync` and then `make baseline-digests`, and name in the Result every file either command rewrote. If a repository contract requires a version change with a content change, it MUST raise the patch version of each edited skill in both version fields; otherwise it MUST leave the versions unchanged.
8. MUST add `skills/concrete_contract_guide_test.go`, which proves that:
   - the embedded bundle holds the guide at `speccheck.ConcreteContractGuidePath`, read from the constant and never copied as a literal;
   - `speccheck.Receipts` reads exactly one receipt from the guide, and `speccheck.ProveReceipt` proves it against the shipped file it names;
   - `speccheck.SurfaceTranscripts` reads exactly one transcript from the guide, and it is well-formed.
9. MUST NOT edit any skill other than the four named here, any Baseline asset by hand, or any file outside Context and the outputs of the two sanctioned commands.

## Subtasks

- [ ] Write the guide with its two worked examples.
- [ ] Add the section and the template section to write-techspec.
- [ ] Add the sections to write-prd, write-tasks and qa-gate, and the template comment.
- [ ] Add the glossary terms and the user guide paragraphs.
- [ ] Regenerate the mirrors and the digests with the sanctioned commands.
- [ ] Add the guide test.

## Acceptance Criteria

- [ ] The embedded bundle holds the guide at the path the horizon constant names. The check reads the guide's worked example as exactly one receipt, proved against a file the skill ships, and exactly one well-formed transcript.
- [ ] The four skills, the two templates, the glossary and the user guide carry the headings and terms this Task names, and each mirror is byte-identical to its canonical skill.
- [ ] `TestSettlementGuidanceIsOneTable` passes, so the `### QA settlement` section is unchanged.
- [ ] After this Task's commit, a Spec whose PRD was committed earlier reports no `SC-RECEIPT-MISSING` and no `SC-TRANSCRIPT-UNDECLARED`.

## Context

- instruction: `docs/adr/0183-an-attributed-claim-carries-a-receipt-the-check-proves.md`
- instruction: `docs/adr/0184-a-techspec-states-a-command-surface-as-a-transcript.md`
- creates: `.agents/skills/write-techspec/references/concrete-contracts.md`
- creates: `skills/write-techspec/references/concrete-contracts.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `skills/concrete_contract_guide_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheConcreteContractGuideShipsAtTheHorizonPath|TestTheGuideExamplesAreReadByTheCheck|TestSettlementGuidanceIsOneTable)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheConcreteContractGuideShipsAtTheHorizonPath TestTheGuideExamplesAreReadByTheCheck TestSettlementGuidanceIsOneTable; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new named tests do not exist, so the command fails.
- `for pair in ".agents/skills/write-techspec/references/concrete-contracts.md|Surface Transcript" ".agents/skills/qa-gate/SKILL.md|zero or more consecutive lines" ".agents/skills/write-techspec/references/concrete-contracts.md|Claim Receipt" ".agents/skills/write-techspec/SKILL.md|## Concrete contracts" ".agents/skills/write-techspec/references/techspec-template.md|### Surface Transcripts" ".agents/skills/write-prd/SKILL.md|## Claim receipts" ".agents/skills/write-prd/references/prd-template.md|Claim Receipt" ".agents/skills/write-tasks/SKILL.md|## Surface Transcripts in Tasks" ".agents/skills/qa-gate/SKILL.md|## Surface Transcripts at the gate" "CONTEXT.md|**Claim Receipt**" "CONTEXT.md|**Surface Transcript**" "docs/user-guide/context-driven-development.md|Claim Receipt" "docs/user-guide/context-driven-development.md|Surface Transcript"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && for skill in write-techspec write-prd write-tasks qa-gate; do diff -r ".agents/skills/$skill" "skills/$skill" >/dev/null || { printf 'mirror drift: %s\n' "$skill" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task the guide does not exist and no skill carries the new headings, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; Core Features 5 and 6; Success Metric 4
- [_techspec.md](_techspec.md) — Skills, guides and glossary; The contract horizon; API Contract 3; Testing Approach 7; Build Order 3
- ADR-0183; ADR-0184

## Result
