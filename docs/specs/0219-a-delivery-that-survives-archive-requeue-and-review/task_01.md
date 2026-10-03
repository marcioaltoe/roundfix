---
task: task_01
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: pending
type: backend
complexity: high
---

# Task 01: A file that pins an active Spec's directory is refused, and Task Context follows the archive

## Overview

On 2026-10-01 Spec 0205's test read its own active TechSpec and failed the
repository gate after the archive commit, and Spec 0218's Task Context named
Spec 0205's active path and failed the corpus check on the delivery's Pull
Request. This Task makes the Spec Consistency Check refuse a non-Markdown file
that names the checked Spec's active directory, so the Daemon's Settlement
Check returns it to the Task that wrote it; makes the Archive Command refuse
such a Spec before moving anything; and resolves another Spec's Task Context
path through the archive root once that Spec is archived.

## Requirements

1. MUST answer items 1 and 2 of the Backlog Entry of 2026-10-01, "Archiving or
   requeueing a Spec loses or breaks its delivery", as `_techspec.md` → The pin
   detector, The archive refusal and Task Context through the archive state.
2. MUST add `CodeSpecPathPinned`, `SpecPathPin` and `ActiveSpecPathPins` in the
   new file `internal/speccheck/active_spec_paths.go` with the signatures of
   `_techspec.md` → Interfaces and invariants 1 and 2: in a Git work tree it
   searches tracked and untracked, non-ignored files with Git, excluding
   Markdown, the Spec Root, its archive root and `docs/history`; outside a
   work tree it walks the directory with the same exclusions; a Spec Root
   outside the repository root yields no pins.
3. MUST call the detector from `internal/speccheck/citations.go` for the
   checked Spec and add one `SC-SPEC-PATH-PINNED` error per pin, located at the
   pin's path and line, with the fix sentence of `_techspec.md` → The pin
   detector; it MUST run whatever the Task statuses are.
4. MUST make `detectTaskContextReferences` resolve a missing path under the
   Spec Root through the archive root as `_techspec.md` → Task Context through
   the archive states (invariant 3), and MUST keep reporting
   `SC-REF-UNRESOLVED` when the path is missing in both places.
5. MUST make `roundfix archive <slug>` refuse, before `spec.Archive` runs and
   with exit `2`, with the reason of API Contract 2 and Surface Transcript 1,
   leaving the Spec directory and every file in place.
6. MUST NOT change `internal/speccheck/constraints.go`,
   `internal/speccheck/coherence.go` or `internal/spec/archive.go`, and MUST
   NOT rewrite any file other than the ones this Task declares.
7. MUST add the tests named in Verification to the new files
   `internal/speccheck/active_spec_paths_test.go`,
   `internal/daemon/settlement_active_spec_path_test.go` and
   `internal/cli/archive_active_spec_path_test.go`, each building its own
   temporary repository; no test reads a file under this repository's Spec
   Root or history root.
8. MUST describe the refusal in `docs/user-guide/commands/archive.md` and the
   Roundfix Skill's `archive` reference, and the `SC-SPEC-PATH-PINNED` finding
   and the archive resolution of Task Context in
   `docs/user-guide/commands/spec.md` and the Skill's `spec` reference; MUST
   raise the Roundfix Skill's version by one patch level above the version on
   this Task's starting tree in both front-matter fields, run
   `make skills-sync`, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.

## Subtasks

- [ ] Add the pin detector and its finding.
- [ ] Resolve Task Context through the archive root.
- [ ] Refuse the archive on a pin.
- [ ] Add the tests.
- [ ] Describe it in the guides and the Skill, and raise its version.

## Acceptance Criteria

- [ ] A Go test in a temporary repository that names `docs/specs/<slug>/_techspec.md`
      yields one `SC-SPEC-PATH-PINNED` error at its line; a Markdown file, a
      file under the archive root, a file under `docs/history` and an ignored
      file naming the same path yield none.
- [ ] Through `SpecCheckSettlementChecker.RefusingFindings`, the pin is a
      refusing finding.
- [ ] A Task Context path of an archived Spec resolves when the archived file
      exists, and stays `SC-REF-UNRESOLVED` when it does not.
- [ ] `roundfix archive` prints Surface Transcript 1's reason, exits `2`, and
      leaves the Spec directory in place; with the pin removed it archives.
- [ ] The guides and the Skill describe the refusal and the finding, the
      mirrors equal their canonical files, and the raised version is recorded.

## Context

- creates: `internal/speccheck/active_spec_paths.go`
- creates: `internal/speccheck/active_spec_paths_test.go`
- creates: `internal/daemon/settlement_active_spec_path_test.go`
- creates: `internal/cli/archive_active_spec_path_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/cli/archive.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/daemon/settlement_checks.go`
- instruction: `internal/spec/archive.go`
- instruction: `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/references/2026-10-01-archiving-or-requeueing-a-spec-loses-or-breaks-its-delivery.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestActiveSpecPathPinIsAnError|TestActiveSpecPathPinSkipsMarkdownSpecRootsAndIgnoredFiles|TestTaskContextResolvesThroughTheArchive|TestTaskContextMissingInBothPlacesIsUnresolved)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestActiveSpecPathPinIsAnError TestActiveSpecPathPinSkipsMarkdownSpecRootsAndIgnoredFiles TestTaskContextResolvesThroughTheArchive TestTaskContextMissingInBothPlacesIsUnresolved; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestSettlementRefusesATaskThatPinsItsSpecPath$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestSettlementRefusesATaskThatPinsItsSpecPath" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestArchiveRefusesASpecAFileStillPins|TestArchiveIgnoresMarkdownThatNamesTheSpec)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveRefusesASpecAFileStillPins TestArchiveIgnoresMarkdownThatNamesTheSpec; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.
- `for file in docs/user-guide/commands/archive.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "cannot archive while another file names its active directory" || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "SC-SPEC-PATH-PINNED" || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md && cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && out="$(go test -count=1 -v -run "^TestEveryOwnedSkillVersionIsRecorded$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; before this Task neither guide nor reference names the refusal or the finding, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1-3, 9; Success Metrics 1-2; Acceptance evidence
- `_techspec.md` → The pin detector; The archive refusal; Task Context through the archive; Interfaces; API Contract 1; API Contract 2; Surface Transcript 1; Testing Approach 1; Build Order 1
- ADR-0223; ADR-0182; ADR-0187; ADR-0189
