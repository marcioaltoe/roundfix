---
task: task_01
spec: 0056-profiles-configure-merge-semantics
status: completed
type: test
complexity: medium
---

# Task 01: Record how the writer behaves today

## Overview

Every later slice changes how the config file is written, and the Spec promises
that a file the current writer handles can never fail or be corrupted by the
new one. This Task records that behavior first, over the config shapes that
actually break naive YAML rewriting. It changes no product behavior — it is the
regression gate the rest of the Spec is measured against, and it only works if
it lands before anything else moves.

## Requirements

1. MUST assemble a corpus of real config shapes and record, for each, the exact
   bytes today's writer produces for a representative configure operation.
2. MUST cover at minimum: a five-category config with Fallback Chains, comments
   above and beside entries, non-default indentation, key order that is not
   alphabetical, a multiline scalar, a YAML anchor and an alias to it, non-ASCII
   values, an empty file, and a file with no `profiles` section.
3. MUST fail with a readable diff naming the affected config and the changed
   region when a later change alters any recorded output.
4. MUST be regenerable through an explicit flag so an intended change is
   re-recorded deliberately, never silently.
5. MUST NOT change any production behavior, exported API, command surface, or
   output in this Task.

## Subtasks

- [ ] Assemble the config corpus covering every listed shape.
- [ ] Record today's writer output for each as a golden.
- [ ] Add the comparison with a readable failure diff.
- [ ] Add the explicit regeneration flag.

## Acceptance Criteria

- [ ] The corpus contains a case for each shape named in Requirement 2.
- [ ] A test writes each corpus config through the current writer and compares
      against its golden, passing on the unmodified tree.
- [ ] Deliberately altering one recorded golden makes the test fail and name the
      affected config.
- [ ] Running the comparison twice in a row produces the same result and does
      not rewrite its own goldens.
- [ ] The goldens can be re-recorded only through the explicit flag.
- [ ] `git status --porcelain` shows no path outside `internal/config/` and this
      task file.

## Context

- interface: `internal/config/profile_config.go`
- interface: `internal/config/config_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0; the corpus matches the current writer.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0 on a second consecutive run, proving the comparison is
  stable and self-recording is gated.
- `grep -rq "TestProfilesConfigWriterCharacterization" internal/config` —
  expected: exit 0; the test exists at the declared name.
- `go test ./internal/config -count=1` — expected: exit 0; the existing suite is
  unaffected.

## References

- `_prd.md` → User Story 6; Core Features 7; Success Metrics (round-trips with
  only the intended change).
- `_techspec.md` → Testing Approach: characterization corpus; Build Order 1;
  Risks (anchors and aliases).

## Result

### Implementation

- Added nine input/golden pairs under
  `internal/config/testdata/profiles_config_writer/`, with one named case for
  each required config shape.
- Added `TestProfilesConfigWriterCharacterization` to drive every input through
  `WriteProfilesConfig`, reload the written config through `Load`, and compare
  the persisted bytes with its golden.
- Added a focused line diff with the config fixture name, hunk line numbers,
  context, and `-`/`+` lines.
- Added the test-only `-update-profiles-config-goldens` flag. Normal comparison
  runs never write golden files.
- Changed no production implementation, exported API, CLI command, or output.

### Focused checks

- Initial focused execution could not use the host Go cache because the
  sandbox denied access. All subsequent Go checks used the repository-local
  writable `.gocache`.
- `go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$' -update-profiles-config-goldens`
  with local `GOCACHE`: exit 0; 10 tests passed and recorded the nine goldens
  only after the explicit flag was supplied.
- Two consecutive `go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$'`
  runs with local `GOCACHE`: exit 0 each; 10 tests passed on each run. SHA-1
  output for every golden was byte-identical before and after both runs.
- After deliberately changing `non_ascii_values.golden.yml`,
  `go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$/non-ASCII'`
  with local `GOCACHE`: expected exit 1. The failure named
  `non_ascii_values` and printed `@@ -1,4 +1,4 @@` with the changed
  `artifact_dir` as `-` and `+` lines. The explicit update flag restored the
  recorded writer output; the same focused subtest then passed.
- Final `go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$'`
  with local `GOCACHE`: exit 0; 10 tests passed after the last test edit,
  including a successful `Load` of every written config.
- `go test ./internal/config -run '^(TestProfileConfig|TestWriteProfilesConfig|TestNormalizeProfilesFragment)'`
  with local `GOCACHE`: exit 0; 7 neighboring profile-config tests passed.
- `git diff --check` and `gofmt -d internal/config/config_test.go`: exit 0 with
  no output.
- Declared Task Verification commands were not run; the Daemon owns them.

### Acceptance criterion evidence

- Corpus coverage: explicit fixtures cover five categories with Fallback
  Chains, comments above and beside entries, non-default indentation,
  non-alphabetical key order, a multiline scalar, a YAML anchor and alias,
  non-ASCII values, an empty file, and a file without `profiles`.
- Current-writer comparison: all nine cases execute the real writer, reload
  successfully, and match exact golden bytes in the final focused run.
- Drift diagnostics: the deliberate `non_ascii_values` mutation failed with
  the fixture name and a readable changed-region hunk.
- Stable comparison: two consecutive normal runs passed and left every golden
  hash unchanged.
- Explicit regeneration: goldens were first created and the deliberate
  mutation was restored only with `-update-profiles-config-goldens`; normal
  runs performed comparison only.
- Changed-path scope: final porcelain status lists only this Task file,
  `internal/config/config_test.go`, and the corpus beneath
  `internal/config/testdata/profiles_config_writer/`.

### Follow-ups

None discovered within this Task slice.
