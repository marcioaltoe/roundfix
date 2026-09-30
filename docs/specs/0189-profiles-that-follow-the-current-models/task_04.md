---
task: task_04
spec: 0189-profiles-that-follow-the-current-models
status: pending
type: docs
complexity: medium
---

# Task 04: The model reference states the shipped snapshot

## Overview

`docs/references/model-selection.md` was last updated on 2026-09-08 and explains the 2026-08-07 ranking. After task_02 the binary ships the Recommended Profile of 2026-09-30, and nothing ties the two together: the catalog and the reference have drifted apart twice. This Task refreshes the reference from `_techspec.md` → Reference data, adds a documentation contract test that fails when the reference and the binary disagree, and corrects the comment in `.roundfixrc.yml` that says the reference is still the August snapshot.

## Requirements

1. MUST open `docs/references/model-selection.md` with `Updated 2026-09-30.` and add current sections that state, from `_techspec.md` → Reference data only:
   - the Recommended Profile of every category, each selection in `runtime / model / effort` form, with its rationale;
   - the Model Catalog and the picker reasoning efforts, with the adapter versions and the date they were read;
   - the published evidence table, each figure with its source link, the note that the figures were read on 2026-09-30 and not re-measured, and `not found` wherever the TechSpec says so;
   - the retirement dates with their source, and the local session measurement.
2. MUST keep every earlier section as history under a heading that says it is not current selection advice. It MUST invent no figure, date or source that the TechSpec does not carry.
3. MUST add `internal/docscontract/model_selection_test.go`, under the `docscontract` build tag, with a test that reads the reference and requires `Updated ` followed by `ModelRecommendationSnapshotVersion`, and, for every Agent Work Category, every selection of `RecommendedProfile` in `runtime / model / effort` form. It MUST read the constant and the function, never a copied literal.
4. MUST prove the new gate can fail. The Result MUST record that the test failed when the reference's date was changed, and that the document was restored.
5. MUST replace, in `.roundfixrc.yml`, only the comment sentence that says `docs/references/model-selection.md` still holds the 2026-08-07 snapshot, with one that says the reference states the same snapshot as the binary. It MUST change no key, no value and no other comment line.
6. MUST NOT edit any guide, skill or Go source other than the new test file.

## Subtasks

- [ ] Refresh the reference from the TechSpec's Reference data, keeping the earlier sections as history.
- [ ] Add the documentation contract test and prove it fails on a wrong date.
- [ ] Correct the one comment in `.roundfixrc.yml`.

## Acceptance Criteria

- [ ] The reference opens with `Updated 2026-09-30.` and lists the Recommended Profile of all ten categories.
- [ ] Every published figure carries a source link, and a missing figure is written as `not found`.
- [ ] Changing the reference's date, or removing a recommended selection from it, fails the new test.
- [ ] `.roundfixrc.yml` differs from its previous content only in that comment sentence.

## Context

- interface: `docs/references/model-selection.md`
- interface: `.roundfixrc.yml`
- creates: `internal/docscontract/model_selection_test.go`
- instruction: `internal/config/recommendations.go`
- instruction: `internal/docscontract/publicdocs_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestModelSelectionReferenceStatesTheShippedSnapshot)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestModelSelectionReferenceStatesTheShippedSnapshot; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/references/model-selection.md|Updated 2026-09-30" "docs/references/model-selection.md|codex / gpt-6.1-sol / high" "docs/references/model-selection.md|https://learn.chatgpt.com/docs/models"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && ! grep -q 'still holds the 2026-08-07 snapshot' .roundfixrc.yml && grep -q 'model-selection.md' .roundfixrc.yml` — expected: exit 0; before this Task the named test does not exist and the reference says `Updated 2026-09-08`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 6; Success Metric 5; Acceptance evidence
- [_techspec.md](_techspec.md) — Reference data; Testing Approach 6; Build Order 4
- ADR-0180

## Result
