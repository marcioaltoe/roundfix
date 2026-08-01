---
task: task_03
spec: 0056-profiles-configure-merge-semantics
status: completed
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

## Result

### Implementation

- `PrepareProfilesConfig` now derives one Effective Change Set from the
  fragment, declared removals, and categories already present in the target
  file, then passes that value to the writer.
- The writer edits the existing `profiles` mapping by key/value pair: a
  replacement keeps the category key and swaps its complete profile value, an
  addition appends one pair, and a removal drops only its pair. Unnamed pairs
  retain their original YAML nodes and order.
- The profile write path detects the existing document indentation before
  encoding, so a two-space or four-space file keeps its indentation instead of
  inheriting the YAML encoder's default.
- Candidate bytes are reparsed before a proposal is returned. A replacement
  that removes an anchor still referenced by a surviving alias now fails before
  persistence.
- Added writer-level regressions for five-category byte preservation,
  single-value diffs, atomic replacement, append order, removal order,
  empty/sectionless adds, and dangling-alias rejection. Re-recorded the
  characterization goldens through their explicit update flag for the intended
  merge and indentation-preservation changes.
- No CLI file changed; confirmation, summaries, exit codes, and machine output
  remain outside this Task slice.

### Focused checks

- Red signal: the focused five-category subtest initially exited `1` and showed
  that configuring `backend` deleted the other four category blocks. The
  dangling-alias subtest also exited `1` because the current writer wrote
  instead of rejecting the invalid candidate.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/config -run
  '^TestProfilesConfigureMergePreservesOtherCategories/' -count=1`: exit `0`;
  all Task 03 writer scenarios passed after the final implementation edit.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/config -run '^TestProfilesConfigWriterCharacterization/'
  -count=1`: exit `0`; every corpus case wrote and reloaded successfully against
  the re-recorded intended output.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/config -run
  '^Test(ProfileConfigAtomic(WritesUserAndProjectProfiles|DryRunAndFailuresLeaveBytesUnchanged)|EffectiveChangeSet)/'
  -count=1`: exit `0`; neighboring proposal, failure, and classification paths
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go vet
  ./internal/config`, `rtk gofmt -d internal/config/profile_config.go
  internal/config/config_test.go`, and `rtk git diff --check`: exit `0` with no
  diagnostics.
- `rtk git diff --exit-code -- internal/cli`: exit `0`; no CLI behavior file
  changed.
- Declared Task Verification commands were not run; the Daemon owns them.

### Acceptance criterion evidence

- Five-category preservation: the focused regression begins with five
  categories at two-space indentation and comments above untouched entries;
  configuring only `backend` produces exact expected bytes with all other
  Fallback Chains, comments, order, and indentation unchanged.
- Single-value diff: that same expected document differs from its input only at
  `backend.preferred.model`; exact whole-file comparison passed.
- Atomic replacement: replacing `backend` removes its prior two-entry Fallback
  Chain and writes only the new complete profile object; no old field survives.
- Removal stability: removing `qa` drops only the `qa` pair while `general`,
  `backend`, and `review` retain their bytes and relative positions.
- Add paths: an added category appends after existing entries, and both an empty
  file and a file without `profiles` produce documents accepted by `Load` while
  preserving unrelated config.
- Alias failure: replacing a category that owns `&backend_model` while
  `artifact_dir` retains `*backend_model` returns an `unknown anchor` validation
  error and leaves the original file byte-identical.
- Corpus compatibility: all nine task-01 corpus shapes passed the focused
  characterization run and reloaded after writing; unchanged nodes survive and
  only the configured profile plus document-indent preservation differ from
  the old replacement-writer goldens.
- Changed-path scope: `git status --porcelain` lists only this Task file,
  `internal/config/profile_config.go`, `internal/config/config_test.go`, and
  characterization goldens under `internal/config/testdata/`.

### Follow-ups

None discovered within this Task slice.
