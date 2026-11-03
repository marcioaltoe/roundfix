---
task: task_04
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: completed
type: docs
complexity: medium
---

# Task 04: The guides, the Roundfix references and the glossary describe omit paths, Claude segments and the refusal

## Overview

The behavior task_01 to task_03 build reaches adopters only through the
guides and the Roundfix skill. This Task writes the Review Omit Path, the
Claude Review Segments, the actionable refusal and the park reason into the
review and deliver guides, the Roundfix skill's review and deliver
references, `CONTEXT.md`, and dated amendment lines on ADR-0227 and
ADR-0282. This answers the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md)),
reported by vortex, which will adopt the key from these guides.

## Requirements

1. MUST, in `docs/user-guide/commands/review.md` and
   `.agents/skills/roundfix/references/review.md`:
   - after the `generated` omission sentences, add that a file matched by an
     entry of `pre_pr_review.omit_paths` in the same base-commit Project
     Config is omitted as `omit-path`, that entries take the grammar of
     `delivery.derived_paths[].paths` and a directory entry ending in `/`
     omits every file beneath it, and that the key never comes from the head;
   - add `omit-path` last to the omission precedence, and state that the
     summary sentence gains `, <o> of pre_pr_review.omit_paths.` only when
     `<o>` is not zero;
   - replace "Claude also retains its token budget below and is never
     segmented" with the sentence "A Claude round-1 prompt above its budget
     is reviewed in at most four Review Segments, each prompt as sent
     estimated within the budget", and describe the per-segment
     `estimatedPromptTokens`, the Claude segment-limit refusal and the
     unchanged single-commit and round-2 refusals (`_techspec.md` → API
     Contracts 4-5);
   - state the suffix of API Contract 6, containing the phrase
     "largest files in the diff", and how files are ranked.
2. MUST, in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`, state next to the Park
   Class table that a `review-blocked` park carries the first line of the
   reason in the Pre-PR Review Record after a colon, and that its Park Class
   stays `review` (API Contract 7).
3. MUST add `**Review Omit Path**` to `CONTEXT.md` through domain-modeling,
   right after **Review Segment**, defined as a Project Config path in
   `pre_pr_review.omit_paths` at the review base whose matching changed files
   the pre-PR review leaves out with reason `omit-path`, distinct from a
   Derived Path Declaration, citing ADR-0352, with an `_Avoid_` line; and
   revise **Review Segment** so it covers a candidate whose diff exceeds the
   review bound or, for Claude, whose prompt exceeds its token budget, citing
   ADR-0352 beside ADR-0282.
4. MUST append to ADR-0227 and ADR-0282 one dated line each,
   `**Amended by [ADR-0352](0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md) (2026-11-03).**`
   followed by one sentence: for ADR-0227, the budget bounds each Claude
   Review Segment's prompt as sent; for ADR-0282, Claude round 1 is segmented
   too and `pre_pr_review.omit_paths` joins the base-commit omissions.
5. MUST copy the canonical Roundfix references with `make skills-sync`, then
   raise both version fields of the Roundfix skill and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   Never hand-edit a mirror or a version field.
6. MUST re-record the Behavior Surface Record as the LAST step, after the
   last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and commit `docs/references/behavior-surfaces.json`.
7. MUST NOT edit any Go file, the configuration guide task_01 owns, the `### QA settlement` section, any other
   Roundfix reference (`modes.md`, `follow-up.md` and `implement.md` belong to
   the parallel Specs), any other skill or any Baseline file.

## Subtasks

- [ ] Write the review guide and reference text.
- [ ] Write the deliver text.
- [ ] Add and revise the glossary entries and the ADR amendment lines.
- [ ] Sync the mirrors, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] The review guide and both reference copies describe the key, the Claude segments and the refusal suffix, and no longer say Claude is never segmented.
- [ ] The deliver guide and both deliver reference copies describe the park reason.
- [ ] `CONTEXT.md` defines **Review Omit Path** and the revised **Review Segment**; the mirrors equal their canonical copies; the skill version and Behavior Surface records pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`
- interface: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill references, the guides and the owned-skill version record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Claude round-1 prompt above its budget is reviewed in at most four Review Segments' || { printf 'missing Claude segments in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'largest files in the diff' || { printf 'missing refusal suffix in %s\n' "$f" >&2; exit 1; }; grep -qF -- 'omit-path' "$f" || { printf 'missing omit-path in %s\n' "$f" >&2; exit 1; }; ! tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'and is never segmented' || { printf 'stale Claude sentence in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'carries the first line of the reason in the Pre-PR Review Record' || { printf 'missing park reason in %s\n' "$f" >&2; exit 1; }; done; grep -F -- '**Review Omit Path**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Review Omit Path\*\*(:| —( |$))' || { echo 'glossary lacks Review Omit Path' >&2; exit 1; }; grep -F -- '**Review Segment**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Review Segment\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'for Claude, whose prompt exceeds its token budget' || { echo 'glossary lacks the revised Review Segment' >&2; exit 1; }; for f in docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md; do grep -qF -- 'Amended by [ADR-0352]' "$f" || { printf 'missing amendment in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides lack the Claude segment sentence, so the first check fails; afterwards every guide, reference copy, glossary entry and ADR line is present, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 7; Glossary; Decisions
- `_techspec.md` → API Contracts 1-7; Build Order 4; Glossary
- ADR-0352; ADR-0227; ADR-0282

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implementation evidence:
- Updated the review guide and canonical Roundfix review reference with pre_pr_review.omit_paths, omit-path precedence and conditional summary counts, Claude round-1 Review Segments and per-segment estimatedPromptTokens, unchanged single-commit/round-2 refusals, and the API Contract 6 suffix ranked by file-section bytes.
- Updated the deliver guide and canonical deliver reference with the review-blocked reason suffix and unchanged review Park Class.
- Revised the Review Segment glossary entry, added Review Omit Path immediately after it with its _Avoid_ line and ADR-0352 citation, and appended the required ADR-0352 amendments to ADR-0227 and ADR-0282.
- make skills-sync: passed; the canonical Roundfix references were copied to their mirrors.
- go test ./skills -run TestEveryOwnedSkillVersionIsRecorded -record-skill-versions -count=1: passed, 1 package; both Roundfix skill version fields and the owned-skill record now contain 0.1.173.
- Mirror comparison with diff -u for the review and deliver references: passed with no differences.
- go test -count=1 -tags docscontract ./internal/docscontract -run TestTheSkillCoverageMapIsCurrent -record-skill-coverage: passed, 1 package; the Behavior Surface Record was re-recorded last after the guide edits.
- git diff --check: passed.

Acceptance evidence:
- Review surfaces: all three review documents contain the omit-path rule, Claude segment sentence, refusal details and “largest files in the diff”; the obsolete “and is never segmented” sentence was removed.
- Deliver surfaces: all three deliver documents contain the Pre-PR Review Record reason after the colon and state that the Park Class remains review.
- Glossary, ADR amendments, mirrors, skill versions and Behavior Surface recording are present in the changed tree; no Go file, other Task, _tasks.md, QA settlement section or unrelated Roundfix reference was changed.

Not run:
- The commands in this Task's ## Verification section and the repository's full Verification, because this Daemon-assigned turn leaves those checks to the Daemon.
