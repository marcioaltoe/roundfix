---
task: task_03
spec: 0082-the-manifest-already-answered-that
status: completed
type: backend
complexity: medium
---

# Task 03: Read the Setup Manifest back into plan inputs

## Overview

The stored Setup Manifest already holds the profile and every decision an update
needs, and nothing reads it as input today. This task adds the resolver that
projects a stored manifest into the inputs plan assembly already accepts, and
that names any decision the current catalog requires but the manifest does not
carry. It is verifiable on its own against fixture repositories, before any
command consumes it.

## Requirements

1. MUST resolve a stored Setup Manifest into a profile identity and a complete
   set of decision values, without prompting and without writing.
2. MUST distinguish, by returned value rather than by error text, the four
   outcomes: no manifest, an unreadable or incompatible manifest, a manifest
   resolving cleanly, and a manifest resolving except for decisions the current
   catalog newly requires.
3. MUST name every newly required decision individually, each with the catalog's
   suggested value and its summary, so a caller can either report them or adopt
   them without a second lookup.
4. MUST NOT treat a catalog suggestion as an answer; a suggestion is carried as
   a suggestion and only an explicit caller choice adopts it.
5. MUST skip decisions the resolved profile fixes, so a profile-fixed value is
   never presented as an unanswered question.
6. MUST treat a changed profile digest as a catalog move and still resolve to an
   update, provided the stored profile resolves and every stored decision still
   validates against it. Only an unresolvable profile, or decisions that no
   longer validate, resolves to the adoption-required outcome, and the returned
   value MUST distinguish which of the two occurred.
7. MUST cover the `standard-typescript-monorepo` decision set in tests, not only
   this repository's profile, because that profile carries the structured typed
   values most likely to be absent from an older manifest.

## Subtasks

- [ ] Add the manifest resolver and its result type.
- [ ] Detect and name newly required decisions against the current catalog.
- [ ] Carry each new decision's suggested value and summary.
- [ ] Distinguish absent, incompatible, clean, and incomplete outcomes.
- [ ] Cover both profile decision sets in tests, including a structured value.

## Acceptance Criteria

- [ ] A repository with no manifest resolves to the absent outcome and no error
      that a caller must parse text to classify.
- [ ] A repository with a valid current manifest resolves to a profile and a
      decision set equal to the manifest's recorded values.
- [ ] A manifest missing a decision the catalog requires resolves to the
      incomplete outcome and names exactly the missing decision ids.
- [ ] A manifest whose profile digest no longer matches, but whose profile
      resolves and whose decisions all validate, resolves to the update outcome
      and not to adoption.
- [ ] A manifest naming a profile that no longer resolves is distinguishable in
      the returned value from one whose decisions no longer validate.
- [ ] Resolution writes nothing: the fixture repository is byte-identical after
      the call.
- [ ] A profile-fixed decision never appears among newly required decisions.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog.go`
- interface: `docs/agents/setup-context.json`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'ManifestInput' -v 2>&1 | grep -q '^--- PASS: .*ManifestInput'` — expected: exits 0, proving the resolver's cases exist and pass.
- `go test ./internal/baseline/ -run 'ManifestInput' -v 2>&1 | grep -q 'typescript'` — expected: exits 0, proving the second profile's decision set is exercised.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 4; Interfaces: `ManifestInput`, `ResolveManifestInput`; Risks: fleet profiles are not all `go-cli-tui`.
- `_prd.md` → Core Features 1, 6, and 8; User Story 5; Goal 1.
- ADR-0047, ADR-0067.

## Result

### Implementation

- Added the read-only `ResolveManifestInput` boundary with typed absent,
  incompatible, resolved, and incomplete states. Incompatible results separately
  identify unreadable manifests, invalid manifests, unresolved profiles, and
  invalid stored decisions.
- Projected recorded decisions through the existing catalog normalization path,
  retained profile-fixed values without treating them as unanswered, and
  reported profile-digest drift without downgrading a valid manifest to
  adoption.
- Added `DecisionSuggestion` values that carry each missing decision's id,
  catalog suggestion or default, and summary without inserting the suggestion
  into the resolved decision set.
- Added resolver tests over temporary fixture repositories for `go-cli-tui` and
  `standard-typescript-monorepo`, including the structured Better Auth value.

### Focused checks

- Red signal: the focused absent-case test initially failed to compile with
  `undefined: ResolveManifestInput`, `undefined: ErrNoManifest`, and undefined
  result-state symbols.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260807T134150Z_b3d9e9dc5c4bccda.task_03/.gocache go test ./internal/baseline -run '^TestResolveManifestInputAbsent$' -count=1`
  passed after implementation.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260807T134150Z_b3d9e9dc5c4bccda.task_03/.gocache go test ./internal/baseline -run '^TestResolveManifestInput' -count=1`
  passed after the Verification Feedback repair (`ok roundfix/internal/baseline`).
- The Task's authored `## Verification` commands were not run; the Daemon owns
  them.

### Verification Feedback repair

- Attempt 1's diagnostic artifact was inspected. The failed quiet grep produced
  no diagnostic body; the related verbose test output used capitalized
  `Typescript`, so it contained no case-sensitive lowercase `typescript` token.
- Renamed the structured-profile test to
  `TestResolveManifestInput_standard_typescript_monorepoStructuredSuggestion`,
  using the canonical profile slug as the evidence-bearing test name instead of
  adding synthetic log output.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260807T134150Z_b3d9e9dc5c4bccda.task_03/.gocache go test ./internal/baseline -run '^TestResolveManifestInput_standard_typescript_monorepoStructuredSuggestion$' -v -count=1`
  passed and its verbose test identifier contains lowercase `typescript`.

### Acceptance evidence

1. `TestResolveManifestInputAbsent` proves a missing manifest returns
   `ManifestInputAbsent` and the matchable `ErrNoManifest` sentinel, with no
   message parsing.
2. `TestResolveManifestInputCurrent` proves a current manifest returns its
   profile id and decision values equal to the recorded values.
3. `TestResolveManifestInputIncompleteDoesNotAdoptSuggestion` and
   `TestResolveManifestInputNamesEveryNewDecision` prove incomplete resolution
   names exactly every missing decision and keeps suggestions out of resolved
   decisions.
4. `TestResolveManifestInputProfileDigestDriftIsResolved` proves a stale profile
   digest remains a resolved update input when the profile and decisions still
   validate.
5. `TestResolveManifestInputDistinguishesAdoptionReasons` proves an unresolved
   profile and invalid stored decisions return distinct incompatibility values.
6. `TestResolveManifestInputWritesNothing` snapshots every fixture file's bytes
   and mode before and after resolution and proves exact equality.
7. `TestResolveManifestInputProfileFixedDecisionIsNotNew` proves a
   profile-fixed decision is resolved from the profile and never reported as a
   new decision.
   `TestResolveManifestInput_standard_typescript_monorepoStructuredSuggestion`
   additionally proves the second profile's structured suggestion remains a
   suggestion and exposes its canonical profile token to Verification.

### Follow-ups

- None discovered within Task 03's slice.
