---
task: task_03
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: backend
complexity: high
---

# Task 03: Stop gating compatibility on content

## Overview

`treeDigest` lives in four places: declared per skill in the setup snapshot,
folded into the catalog digest, validated in `catalog_validate.go`, and
rewritten by `assets_sync.go`. That coupling is why editing one owned skill
broke `make verify` twice in one session — once through each characterization
corpus — while `make baseline-digests` reported no changes, because neither
corpus is a member of its steps.

This slice removes content from the compatibility gate. It does not remove
digests that protect artifacts Roundfix genuinely owns.

## Requirements

1. MUST stop the catalog digest from folding skill content, so an owned skill
   edit no longer moves it.
2. MUST stop the characterization corpora from embedding volatile skill digests
   in recorded diagnostics, so a skill edit cannot invalidate them.
3. MUST keep every digest that protects an artifact Roundfix generates. The
   change is narrow: content pins stop gating *compatibility*, nothing else.
4. MUST leave a Baseline applied before this Spec validating unchanged.
5. MUST leave archived Spec artifacts byte-identical.
6. MUST assert that editing an owned skill leaves `make verify` green with no
   regeneration step — the Spec's first Success Metric and the reason it
   exists.
7. MUST preserve ADR-0085: a regeneration run stays ungated by the pins it
   rewrites while every other load stays strict.

## Subtasks

- [ ] Remove skill content from the catalog digest.
- [ ] Remove volatile skill digests from both corpora's recorded diagnostics.
- [ ] Assert the edit test and the back-compatibility corpus.

## Acceptance Criteria

- [ ] Editing an owned skill leaves `make verify` green with no regeneration
      step, asserted end to end.
- [ ] The catalog digest is unchanged by an owned skill edit.
- [ ] Neither characterization corpus records a skill digest.
- [ ] Digests protecting generated guides are unchanged, asserted.
- [ ] A Baseline applied before this Spec still validates.
- [ ] Archived Spec artifacts are byte-identical.
- [ ] A regeneration run stays ungated by the pins it rewrites.

## Context

- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/catalog_validate.go`
- instruction: `docs/adr/0085-a-regeneration-run-is-not-gated-on-the-pins-it-rewrites.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/baseline -count=1 -run 'Digest|Catalog|Characterization|Compatibility' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the digest and corpus tests ran and passed.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 7, 8 and 9; Success Metrics 1 and 8.
- `_techspec.md` → Build Order 3; Risks & Considerations.
- ADR-0081, ADR-0085.
