---
task: task_02
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
status: pending
type: test
complexity: medium
---

# Task 02: Documentation contract tests read the conventions and the event categories from the code

## Overview

Nothing compared the Roundfix skill's references with the code they describe,
so the reviewer's conventions and the Run Event Stream categories moved and
the references stayed behind. This Task adds documentation contract tests in
`internal/docscontract` that read the conventions and their version by parsing
`internal/cli/review_conventions.go`, and the categories from
`runevent.AllStreamCategories()`, and compare them with the documents task_01
corrected; it also adds the glossary's Run Event Stream entry to the
closed-list contract (`_prd.md` → Core Feature 4). It answers the Inbox Entry
of 2026-10-20 adopted at
[`references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md`](references/2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md).

## Requirements

1. MUST add to `internal/docscontract/command_documentation_test.go`
   `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`, as
   `_techspec.md` → API Contract 4 specifies. It parses
   `internal/cli/review_conventions.go` with `go/parser` from the root that
   `baselineDocumentationRepoRoot()` returns, takes the string value of the
   `deliveryConventionsVersion` constant and, from the body of the function
   `deliveryConventions`, every composite literal of exactly two string
   literals whose first matches `^C[0-9]+$`, in source order, and fails when
   the version is empty or no convention is found. It then checks the
   canonical `.agents/skills/roundfix/references/review.md`, the mirror
   `skills/roundfix/references/review.md` and
   `docs/user-guide/commands/review.md`: on each document's text with every
   whitespace run collapsed to one space, the text contains
   "Delivery Conventions version " followed by the version in code font,
   contains "- <ID>. <prompt>" for each convention with the prompt's
   whitespace collapsed the same way, and holds exactly as many list items
   matching `- C[0-9]+\. ` as there are conventions. A failure names the
   document's repository-relative path and the convention ID, or the version.
   It never copies a convention or the version as a literal into the test.
2. MUST add `TestTheDeliveryConventionContractReportsAStaleConvention`, which
   takes the canonical reference's text, replaces the parsed C3 prompt with
   "The archive commit moves a completed Spec's directory to the archive root
   and stamps its archive front matter.", runs the same check on it, and
   requires an error that names
   `.agents/skills/roundfix/references/review.md` and `C3`.
3. MUST add `TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem`, as
   `_techspec.md` → API Contract 5 specifies. It takes the keys of
   `runevent.AllStreamCategories()`, sorted. In the canonical
   `.agents/skills/roundfix/references/events.md` and its mirror, after the
   anchor "Default replay emits these public categories in journal cursor
   order:", and in `docs/user-guide/commands/events.md`, after the anchor "The
   public Supervisor categories, in journal cursor order, are", on
   whitespace-collapsed text, it collects every code span up to the first `.`
   or `;` and requires them, sorted, to equal the categories; a missing anchor
   fails. In both reference copies it also requires, for each category, a
   line starting with "| " followed by the category in code font and " |". A
   failure names the document's path and each missing or extra category.
4. MUST add `TestTheEventCategoryContractReportsAMissingCategory`, which
   removes `usage` from the canonical reference's category sentence (the
   code span and its preceding ", and " joiner, leaving the sentence's
   ending), runs the same check, and requires an error that names
   `.agents/skills/roundfix/references/events.md` and `usage`.
5. MUST add to `glossaryClosedLists` in
   `internal/docscontract/user_guide_contract_test.go` the term
   `Run Event Stream` with the values of `runevent.AllStreamCategories()`,
   built in the entry as the existing Task Type entry builds its values, so
   `TestGlossaryClosedListsEqualTheirConstants` holds the glossary entry to
   the six categories and
   `TestGlossaryContractReportsAPlantedViolationOfEachRule` keeps passing.
6. MUST make each new top-level test call `t.Parallel()` as its first
   statement, read only repository files, and compile for Windows; add no new
   file to `internal/docscontract`, because new contract tests go into the
   existing `//verify:always` files.
7. MUST NOT change any production Go source, any document, any other test,
   or the Contract Relevance declarations.

## Subtasks

- [ ] Parse the reviewer's conventions and version and check the three `review` documents, with its planted-violation test.
- [ ] Check the category sentence and table rows of the `events` documents against the stream, with its planted-violation test.
- [ ] Add the Run Event Stream closed list.

## Acceptance Criteria

- [ ] The convention guard passes on the corrected documents and fails, naming the path and `C3`, on the old C3 sentence.
- [ ] The category guard passes on the corrected documents and fails, naming the path and `usage`, when `usage` is removed.
- [ ] The glossary closed-list contract requires the six categories in the Run Event Stream entry.

## Consumers

None. This Task adds tests only; no Shared Contract changes and no other suite consumes these helpers.

## Context

- instruction: `internal/cli/review_conventions.go`
- instruction: `internal/runevent/stream.go`
- instruction: `.agents/skills/roundfix/references/review.md`
- instruction: `.agents/skills/roundfix/references/events.md`
- instruction: `docs/user-guide/commands/review.md`
- instruction: `docs/user-guide/commands/events.md`
- instruction: `CONTEXT.md`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `internal/docscontract/user_guide_contract_test.go`

## Verification

- `for name in TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem TestTheDeliveryConventionContractReportsAStaleConvention TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem TestTheEventCategoryContractReportsAMissingCategory; do grep -q "^func $name(t .testing.T) {" internal/docscontract/command_documentation_test.go || { printf 'missing %s\n' "$name" >&2; exit 1; }; done; grep -q 'term: "Run Event Stream"' internal/docscontract/user_guide_contract_test.go || { echo 'missing Run Event Stream closed list' >&2; exit 1; }; GOOS=windows go vet -tags docscontract ./internal/docscontract || exit 1; go test -count=1 -tags docscontract -run '^(TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem|TestTheDeliveryConventionContractReportsAStaleConvention|TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem|TestTheEventCategoryContractReportsAMissingCategory|TestGlossaryClosedListsEqualTheirConstants|TestGlossaryContractReportsAPlantedViolationOfEachRule)$' ./internal/docscontract` — expected: exit 0. Before this Task none of the four tests exists and the glossary has no Run Event Stream closed list, so the command fails at its first name check. Afterwards the four guards and the glossary closed-list contract pass against the documents task_01 corrected, and the package compiles for Windows.

## References

- `_prd.md` → Core Feature 4; Goals; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contracts 4-7; Invariants 1-4; Testing Approach; Build Order 2
