---
task: task_03
spec: 0056-profiles-configure-merge-semantics
status: pending
type: backend
complexity: high
---

# Task 03: Merge by category instead of replacing the map

## Overview

This is the slice that fixes the reported defect. The writer builds a fresh
profiles mapping from the fragment and swaps the whole node, so every category
the fragment did not name loses its Fallback Chain along with its comments and
formatting. This Task merges into the existing mapping instead, editing only
the categories the change set names and leaving every other entry's nodes
untouched — which is what preserves both the data and the formatting.

## Requirements

1. MUST merge into the existing profiles mapping rather than replacing it: a
   replaced category swaps only its own value, an added category appends a new
   entry, and a removed category drops its entry.
2. MUST leave every category absent from the change set byte-identical in the
   written file, including its comments, key order, and indentation.
3. MUST keep a replaced category atomic — the whole profile object is replaced,
   never merged field by field.
4. MUST keep working for a config with no profiles section and for an empty
   file, where the merge is an add rather than an error.
5. MUST fail without writing when the operation would produce a document that
   no longer parses, rather than writing it — specifically when an alias
   elsewhere refers to an anchor defined inside a replaced category.
6. MUST NOT change the confirmation, the summary, exit codes, or the machine
   output shape in this Task.

## Subtasks

- [ ] Merge the change set into the existing mapping node.
- [ ] Preserve untouched entries' nodes rather than rebuilding them.
- [ ] Keep the absent-section and empty-file paths working as adds.
- [ ] Fail closed on an alias whose anchor a replacement would remove.

## Acceptance Criteria

- [ ] Configuring one category in a five-category file leaves the other four
      categories byte-identical, including their Fallback Chains and comments.
- [ ] A single-value change produces a diff touching only that value's lines.
- [ ] A replaced category is replaced as one object, with no field-level merge.
- [ ] A removed category's entry is gone and no other entry moved.
- [ ] Merging into a file with no profiles section, and into an empty file, both
      succeed and produce a valid document.
- [ ] An alias referring to an anchor inside a replaced category fails the
      operation with nothing written.
- [ ] Every corpus config from task 01 still round-trips without a parse
      failure; only the intended values differ.
- [ ] `git status --porcelain` shows no path outside `internal/config/` and this
      task file.

## Context

- interface: `internal/config/profile_config.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/config -run TestProfilesConfigureMergePreservesOtherCategories -count=1`
  — expected: exit 0; untouched categories are byte-identical.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0; every corpus config still writes without corruption, with
  the goldens re-recorded only for the intended merge difference.
- `go test ./internal/config -run TestEffectiveChangeSet -count=1` — expected:
  exit 0; classification still holds.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Stories 1 and 3; Core Features 1 and 4; Success Metrics
  (byte-identical against a five-category config).
- `_techspec.md` → Implementation Design: Interfaces; Build Order 3; Risks
  (anchors and aliases).
- ADR-0049.
