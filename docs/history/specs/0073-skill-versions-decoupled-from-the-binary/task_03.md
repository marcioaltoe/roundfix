---
task: task_03
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
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

## Result

### Implementation

- Setup catalog identity now canonicalizes the compatibility projection: it
  excludes repository-owned `treeDigest` and legacy `contentDigest` fields,
  while external GitHub `treeDigest` values remain identity-bearing.
- Setup validation and Assets Sync no longer require, preserve, or write a
  content pin for repository-owned skills. Their compatibility contract keeps
  the repository source and `minimumVersion`.
- Catalog, parity, diagnostic, and plan-characterization fixtures were
  re-recorded for the structural contract change. Generated-guide golden
  digests and ADR-0085's regeneration-only relaxation were left intact.
- Integration-tagged regression coverage now edits the canonical and mirrored
  `roundfix` skill in an isolated tracked repository, runs the exact real
  `make verify` target without a regeneration target, and proves derived and
  archived artifacts remain byte-identical. Keeping this macro assertion out
  of ordinary package tests prevents a test suite from recursively invoking
  itself.

### Focused checks

- Before the catalog identity change,
  `TestCatalogDigestExcludesOwnedSkillContent` failed because a legacy owned
  content-pin edit moved the catalog digest. The same focused test passed after
  the implementation; its external-source control still moves the digest.
- The focused Baseline suite covering catalog identity, compatibility,
  diagnostics, both characterization paths, Assets Sync, and regeneration
  semantics passed: 19 tests in `internal/baseline`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-repair-go-cache rtk go test
  ./skills -run
  '^(TestAuthorialSkillSync|TestAuthorialSkillSyncUpdateModeRoundTrip|TestCharacterizationCorporaDoNotRecordOwnedSkillDigests)$'
  -count=1` passed: 20 tests.
- Daemon Verification attempt 1 exposed a test-harness recursion: the ordinary
  `skills` suite ran `TestOwnedSkillEditLeavesMakeVerifyGreen`, which launched
  another repository-wide suite while the outer suite was still active. The
  nested run then hit an unrelated temporary Git worktree read failure under
  duplicated load. A two-package focused reproduction passed before the
  repair, confirming the failure was load-sensitive rather than a catalog
  behavior regression.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-repair-go-cache go test
  -tags=integration ./skills -run
  '^TestOwnedSkillEditLeavesMakeVerifyGreen$' -count=1 -v` passed in 135.96
  seconds. Its nested exact `make verify RTK=` run succeeded, did not mention
  or invoke `baseline-digests`, and left the asserted derived and archived path
  digests unchanged.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-repair-go-cache rtk go
  test ./internal/cli ./skills -run
  '^(TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestOwnedSkillEditLeavesMakeVerifyGreen)$'
  -parallel 16 -count=1 -v` passed: one ordinary test across two packages. This
  proves the integration macro is absent from the default suite while the CLI
  worktree regression named by Daemon diagnostics remains green.
- `git diff --check` and `git diff --quiet -- docs/specs/_archived` exited 0.
  A focused search found no `contentDigest` in setup snapshots, catalog
  diagnostics, plan characterization goldens, or the Assets Sync parity
  fixture.
- The Task's declared `## Verification` commands were not run; terminal
  verification remains Daemon-owned.

### Acceptance evidence

- **Owned skill edit:** the integration-tagged
  `TestOwnedSkillEditLeavesMakeVerifyGreen` exercised the exact real
  verification target after editing owned skill bytes, with no regeneration
  step or artifact rewrite. Ordinary package tests do not recursively launch
  that repository-wide target.
- **Stable catalog digest:** `TestCatalogDigestExcludesOwnedSkillContent`
  asserts identical catalog digest and normalized identity after an owned
  content edit, plus the negative external-tree control.
- **Digest-free characterization:**
  `TestCharacterizationCorporaDoNotRecordOwnedSkillDigests` dynamically hashes
  every owned skill and rejects those values from the catalog-diagnostic and
  plan-characterization corpora; the structural absence audit also passed.
- **Generated-guide protection:** `TestCatalogRegenerationMode` and
  `TestRegenerationBreaksGoldenDigestCycle` passed with generated-guide golden
  digests still strict on ordinary loads and deferred only during regeneration.
- **Pre-Spec Baseline compatibility:** `TestCatalogCompatibility` passed its
  legacy Baseline corpus, and the same-baseline characterization remained
  valid after removing only volatile owned-skill pins.
- **Archived Specs:** the end-to-end path digest comparison passed, and the
  repository diff contains no archived Spec change.
- **ADR-0085:** the focused regeneration tests passed both halves of the
  contract: ordinary loads remain strict, while regeneration is not gated on
  the pins it rewrites.
