---
task: task_04
spec: 0239-a-glossary-every-spec-keeps-current
status: completed
type: docs
complexity: medium
---

# Task 04: The glossary gains the terms dropped since 2026-10-02 and this Spec's own

## Overview

`CONTEXT.md` last changed on 2026-10-02 while the Specs after it introduced
terms their reports named and nobody wrote down
([the adopted Backlog Entry of 2026-10-06](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)).
On 2026-10-06 the maintainer chose to catch up by editing `CONTEXT.md` directly
through `domain-modeling` ("Sim"). This is the Spec's docs Task: it writes every
term the Glossary Declaration in `_prd.md` adds or changes, from
`_techspec.md` → Catch-up glossary entries, each checked against the ADR that
decided it (ADR-0244). It is verifiable on its own: the glossary defines every
declared term.

## Requirements

1. MUST activate `domain-modeling` and write to `CONTEXT.md`, the repository's
   selected domain context, even though the skill names its glossary file
   `GLOSSARY.md`; MUST NOT create a `GLOSSARY.md`.
2. MUST append, at the end of `CONTEXT.md` and in the given order, every entry
   of `_techspec.md` → Catch-up glossary entries, each as a bold heading
   followed by a colon, its definition and its `_Avoid_` line, with one blank
   line between entries; and MUST add the sentence that section gives to the
   end of the first paragraph of the existing Spec Consistency Check entry.
   Every other line of `CONTEXT.md` MUST stay as it is.
3. MUST read, before writing each catch-up entry, the ADR named for it in that
   section, and correct the entry to the ADR's meaning where they disagree. The
   Result MUST name each ADR read and each correction, or state that none was
   needed.
4. MUST keep each definition free of Spec numbers, Spec paths, file paths
   other than the Item Branch name pattern, and implementation detail, per the
   glossary format, and MUST NOT add the five candidates that section leaves
   out.
5. MUST NOT edit any other file.

## Subtasks

- [ ] Read each source ADR.
- [ ] Append the entries and the Spec Consistency Check sentence.
- [ ] Record the ADRs read and any correction.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines the twelve added terms and carries the Spec
      Consistency Check sentence.
- [ ] Each catch-up definition agrees with its ADR.
- [ ] No `GLOSSARY.md` exists and no other file changed.

## Context

- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `docs/adr/0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`
- instruction: `docs/adr/0231-the-jev-ceiling-is-a-user-config-value.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`
- instruction: `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- instruction: `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- interface: `CONTEXT.md`

## Verification

- `for phrase in '**Glossary Declaration**:' '**Glossary Gap**:' '**Light Tier**:' '**Light Spend Log**:' '**Jev Ceiling**:' '**Stage Key**:' '**Tested Base**:' '**Merge Evidence**:' '**Delivery Commit**:' '**Item Branch**:' '**Network-Denied Row**:' '**Pre-PR Review Command**:' '**Spec Consistency Check**:' 'Its Glossary Declaration gap begins at the glossary horizon'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$phrase" || { printf 'missing in CONTEXT.md: %s\n' "$phrase" >&2; exit 1; }; done; test ! -e GLOSSARY.md || { printf 'GLOSSARY.md must not exist\n' >&2; exit 1; }` — expected: exit 0; before this Task the glossary defines none of the twelve added terms, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 1; Core Feature 7; Success Metric 5; Glossary
- [_techspec.md](_techspec.md) — Catch-up glossary entries; API Contract 8; Testing Approach 6; Build Order 4
- ADR-0244; ADR-0174; ADR-0196; ADR-0231; ADR-0232; ADR-0236; ADR-0238; ADR-0239; ADR-0240

## Result

Implemented the twelve catch-up glossary entries in the specified order at the
end of `CONTEXT.md` and added the Glossary Declaration horizon sentence to the
first paragraph of the existing Spec Consistency Check entry. The definitions
contain no Spec identifiers or paths beyond the required Item Branch pattern,
and the five excluded candidates were not added.

Read ADR-0244, ADR-0174, ADR-0196, ADR-0231, ADR-0232, ADR-0236, ADR-0238,
ADR-0239, and ADR-0240 before writing the entries. No correction to the
TechSpec wording was needed: each entry agrees with its named ADR or ADRs.

Focused checks after the edit:

- `rg -n` confirmed all twelve bold headings, the Spec Consistency Check
  sentence, and none of the five excluded candidate names in the changed
  glossary section.
- `git diff --check` exited 0.
- `git status --short` showed only the pre-existing Task file change and this
  Task's permitted `CONTEXT.md` change; no `GLOSSARY.md` was created.

The declared Verification and repository-wide verification remain for the
Daemon, per the daemon-assigned execution contract.
