---
task: task_07
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: pending
type: backend
complexity: low
---

# Task 07: An ADR whose header already lists itself loses only its own identifier

## Overview

Corrective Task 1 of 2. The pre-PR review of 2026-11-19 (provider codex, artifact `81eee2e07bb96026`, reviewed head `d8506932`, finding F1, anchored at `internal/spec/adr_header.go:184`) found that an existing `amended_by` list holding the ADR's own identifier survives the union `WithADRAmendedBy` writes, so `roundfix baseline update` leaves an ADR listing itself, contrary to `_techspec.md` → Invariant 7. The caller `documentationNotices` in `internal/cli/baseline_update_documentation.go` drops the ADR's own identifier only from the notices; both header fills, in the plan and in the apply, still union the header's own identifier back.

Invariant 7 says the list written is the union of the identifiers the header already lists and those the notices name, that "an identifier a person wrote is never removed", and that "an ADR never lists itself". The two clauses meet only when a person wrote the ADR's own identifier, and "never lists itself" is the narrower rule, so it wins there: the list written drops the ADR's own identifier and keeps every other identifier the header lists. The ADR's identifier is the number its file name begins with, compared by number, so `ADR-55` and `ADR-055` are the same identifier. When dropping it leaves no identifier, the line is written as `amended_by: []`, so the fill still changes exactly one line (Invariant 6).

## Requirements

1. MUST answer finding F1 of the pre-PR review of 2026-11-19: no header fill of `roundfix baseline update`, in the plan or in the apply, writes an ADR's own identifier into its Amendment List, and every other identifier its header lists stays.
2. MUST keep `WithADRAmendedBy(content []byte, ids []string) []byte` with the signature `_techspec.md` → Interfaces declares and its nine existing subtests unchanged and passing, and add to `internal/spec/adr_header.go` the function `WithADRAmendedByExcept(content []byte, ids []string, own string) []byte`, which writes the union of the header's identifiers and `ids` without any identifier whose number equals the number of `own`; `WithADRAmendedBy(content, ids)` returns the same bytes as `WithADRAmendedByExcept(content, ids, "")`.
3. MUST, when the header lists the ADR's own identifier and nothing else remains, write `amended_by: []` in place of that line; a header with no Amendment List and no identifier to write stays byte for byte as it was.
4. MUST make both callers in `internal/cli/baseline_update_documentation.go`, the plan's fill and the apply's fill, call `WithADRAmendedByExcept` with the ADR's own identifier taken from its file name, the same number `documentationNotices` already compares, through one helper both use.
5. MUST add `TestWithADRAmendedByExcept` to `internal/spec/adr_header_test.go` with the subtests:
   - `drops_its_own_identifier`: the header `amended_by: [ADR-42, ADR-55]`, the ids `ADR-43` and the own identifier `ADR-55` give `amended_by: [ADR-42, ADR-43]`;
   - `own_by_number`: the header `amended_by: [ADR-55, ADR-42]` and the own identifier `ADR-055` give `amended_by: [ADR-42]`;
   - `only_itself_leaves_an_empty_list`: the header `amended_by: [ADR-55]`, no ids and the own identifier `ADR-55` give `amended_by: []`, with every other byte kept;
   - `empty_own_is_the_union`: with the own identifier empty the bytes equal those of `WithADRAmendedBy` for the same input;
   - each subtest also proves the input is not mutated and that a second write gives the same bytes.
6. MUST add the subtest `drops_its_own_identifier` to `TestBaselineUpdateDocumentationPart` in `internal/cli/baseline_update_documentation_test.go`: an adopted fixture holds `docs/adr/0055-lists-itself.md` whose header is `status: accepted`, `knowledge: false` and `amended_by: [ADR-42, ADR-55]` and whose body carries `Amended by ADR-43.`; after `roundfix baseline update --yes` the file's header reads `amended_by: [ADR-42, ADR-43]`, every other byte is kept, and a second update reports the documentation part current.
7. MUST keep every existing subtest of `TestBaselineUpdateDocumentationPart`, `TestReadADRHeader`, `TestADRAmendmentNotices` and `TestWithADRAmendedBy` unchanged and passing, and change no other test, guide, glossary, record, skill, Baseline module, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Add `WithADRAmendedByExcept` and make `WithADRAmendedBy` call it with an empty own identifier.
- [ ] Make the plan's and the apply's header fills pass the ADR's own identifier through one helper.
- [ ] Add the four writer subtests and the command subtest.

## Acceptance Criteria

- [ ] An ADR whose header lists its own identifier and another identifier is filled without its own identifier and keeps the other one, through the pure writer and through `roundfix baseline update`, and a second update is current.
- [ ] The nine subtests of `TestWithADRAmendedBy` and the fourteen existing subtests of `TestBaselineUpdateDocumentationPart` still pass.

## Context

- interface: `internal/spec/adr_header.go`
- interface: `internal/spec/adr_header_test.go`
- interface: `internal/cli/baseline_update_documentation.go`
- interface: `internal/cli/baseline_update_documentation_test.go`

## Consumers

- `WithADRAmendedBy` and `WithADRAmendedByExcept` — consumers: `./internal/spec`, `./internal/cli`
- The header fills of the documentation part of `roundfix baseline update` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestWithADRAmendedBy|TestWithADRAmendedByExcept)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWithADRAmendedBy/after_knowledge TestWithADRAmendedBy/after_status TestWithADRAmendedBy/replaces_the_list TestWithADRAmendedBy/keeps_written_identifiers TestWithADRAmendedBy/ascending_by_number TestWithADRAmendedBy/keeps_updated_at TestWithADRAmendedBy/crlf TestWithADRAmendedBy/no_frontmatter_gains_header TestWithADRAmendedBy/same_bytes_twice TestWithADRAmendedByExcept/drops_its_own_identifier TestWithADRAmendedByExcept/own_by_number TestWithADRAmendedByExcept/only_itself_leaves_an_empty_list TestWithADRAmendedByExcept/empty_own_is_the_union; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec` — expected: exit 0. Before this Task `WithADRAmendedByExcept` and its test do not exist, so the package does not compile and the command exits 1; afterwards the nine writer cases and the four own-identifier cases pass, and the package passes whole.
- `out="$(go test -count=1 -timeout 40m -v -run '^TestBaselineUpdateDocumentationPart$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in plan_names_the_list_and_the_losses text_section json_member digest_moves_only_when_pending plan_ready_message file_changes confirmed_apply second_run_current secondbrain_off marker_refused_rest_applies fills_amendment_lists fill_beside_history_unit no_history_does_not_skip apply_failure drops_its_own_identifier; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestBaselineUpdateDocumentationPart/$name " || { echo "missing pass: TestBaselineUpdateDocumentationPart/$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the subtest `drops_its_own_identifier` does not exist, and with only the subtest added the update writes `amended_by: [ADR-42, ADR-43, ADR-55]`, so the command exits 1; afterwards the fifteen subtests pass.

## References

- Pre-PR review of 2026-11-19 (artifact `81eee2e07bb96026`, head `d8506932`) → F1
- `_techspec.md` → Interfaces; Invariants 6 and 7
- task_02 → Requirement 1
