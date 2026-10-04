---
task: task_03
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: completed
type: backend
complexity: low
---

# Task 03: Project Config declares the repository's Roundfix item build

## Overview

The queue owner can build an item's own Roundfix binary only when the
repository says how. This Task adds the `delivery.item_binary` declaration,
with a `build` command and the repository-relative `path` it writes, to the
configuration model, with the validation and errors of "The declaration".
Nothing reads it yet; task_04 does.

## Requirements

1. MUST add `ItemBinaryDeclaration` and the `Declared` method with the shape
   of the TechSpec's Interfaces, and the `ItemBinary` field to `Delivery`.
2. MUST decode `delivery.item_binary` through a custom `UnmarshalYAML` that
   refuses any key other than `build` and `path` with
   `delivery.item_binary.<key> is not a supported config key`.
3. MUST make `applyOverlay` replace the whole declaration when an overlay
   carries it, so Project Config replaces User Config, and leave the zero
   value when no overlay declares it.
4. MUST validate a declared value beside `validateDerivedPaths` with exactly
   `delivery.item_binary requires build and path` for an empty field and
   `delivery.item_binary has unsafe path "<path>"` for a path that is
   absolute, holds a backslash, has a `..` segment, is `.`, or differs from
   its `path.Clean` form, on every load path that validates `derived_paths`,
   including `ResolveConfigProposal`.
5. MUST add the tests named in Verification to the new file
   `internal/config/delivery_item_binary_test.go`: a Project Config
   declaration is read; a Project declaration replaces a User one; an
   undeclared configuration yields a value whose `Declared` is false; each
   error of API Contract 4 is returned for its input.
6. MUST NOT change `delivery.derived_paths` behavior, any other key, or
   `.roundfixrc.yml`.

## Subtasks

- [ ] Add the declaration type and field.
- [ ] Decode it with the unknown-key refusal.
- [ ] Layer and validate it.
- [ ] Add the declaration tests.

## Acceptance Criteria

- [ ] `build: make build` with `path: bin/roundfix` in Project Config loads
      as a declared value with those two fields.
- [ ] A User declaration is replaced by a Project one; no declaration loads
      as undeclared.
- [ ] An empty field, `/abs/roundfix`, `../roundfix`, `bin\\roundfix`, `.` and
      an unknown sub-key each fail with their exact message.

## Context

- interface: `internal/config/config.go`
- interface: `internal/config/delivery.go`
- creates: `internal/config/delivery_item_binary_test.go`
- instruction: `internal/config/delivery_derived_paths_test.go`
- instruction: `docs/user-guide/configuration.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliveryItemBinaryIsReadFromProjectConfig|TestDeliveryItemBinaryProjectReplacesUser|TestDeliveryItemBinaryIsUndeclaredByDefault|TestDeliveryItemBinaryRefusesAnIncompleteDeclaration|TestDeliveryItemBinaryRefusesAnUnsafePath|TestDeliveryItemBinaryRefusesAnUnknownKey)$" ./internal/config 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliveryItemBinaryIsReadFromProjectConfig TestDeliveryItemBinaryProjectReplacesUser TestDeliveryItemBinaryIsUndeclaredByDefault TestDeliveryItemBinaryRefusesAnIncompleteDeclaration TestDeliveryItemBinaryRefusesAnUnsafePath TestDeliveryItemBinaryRefusesAnUnknownKey; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six tests do not exist, so the command fails.

## References

- `_prd.md` → User Story 4; Core Feature 1
- `_techspec.md` → The declaration; Interfaces; API Contract 4; Testing Approach 2; Build Order 3
- ADR-0225; ADR-0192; ADR-0027


## Result

Implemented the Task 03 configuration slice. `Delivery.ItemBinary` holds the
specified two-field `ItemBinaryDeclaration`; its zero value is undeclared.
The YAML decoder refuses unknown nested keys and explicitly empty mappings.
Project Config replaces the whole User Config declaration. Shared `Validate`
checks completeness and repository-relative path safety beside derived-path
validation, covering both `Load` and `ResolveConfigProposal`.

Acceptance evidence from `internal/config/delivery_item_binary_test.go`:

| Acceptance criterion | Focused-check evidence |
| --- | --- |
| Project Config reads `make build` and `bin/roundfix` as declared | `TestDeliveryItemBinaryIsReadFromProjectConfig` asserts both fields and `Declared()` through `Load` and `ResolveConfigProposal`. |
| Project replaces User; absence is undeclared | `TestDeliveryItemBinaryProjectReplacesUser` proves replacement, inheritance when Project omits the declaration, and refusal of a partial Project declaration rather than field merging. `TestDeliveryItemBinaryIsUndeclaredByDefault` proves the zero value and absent configuration. |
| Incomplete declarations, unsafe paths, and unknown keys fail with the contract messages | `TestDeliveryItemBinaryRefusesAnIncompleteDeclaration`, `TestDeliveryItemBinaryRefusesAnUnsafePath`, and `TestDeliveryItemBinaryRefusesAnUnknownKey` exercise both entry points. Cases include missing, empty and whitespace fields, an empty mapping, `/abs/roundfix`, `../roundfix`, a nested parent segment, a backslash, `.`, non-clean paths, and an unknown sub-key in both scopes. Parsing retains its existing scope wrapper; tests compare the underlying contract error exactly. |

Focused checks:

- Before implementation, `GOCACHE=/tmp/roundfix-task03-gocache rtk go test
  ./internal/config -run '^TestDeliveryItemBinaryIsReadFromProjectConfig$'
  -count=1` exited 1: the new type and field did not exist.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/config
  -run 'TestDeliveryItemBinary|TestProjectConfigReadsDerivedPathDeclarations|TestDerivedPathDeclarationsRefuseUnsafeEntries'
  -count=1` exited 0.
- After the final code edit, `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy
  go test -count=1 ./internal/config` exited 0 (`ok`, 0.881s), exercising the
  six new named tests and the existing config tests, including the unedited
  derived-path tests.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The only pre-existing worktree change was the Daemon's `status: in_progress`
in this Task file, which was preserved. No declared Verification command,
repository-wide gate, commit, push, or Pull Request operation was run.
`delivery.derived_paths`, `.roundfixrc.yml`, the Task Graph, and other Task
files were not changed. Task 04 owns consuming the declaration; this diff
adds no binary build or selection behavior. Implementation is handed back
for Daemon Verification and settlement.

## Carry-forward provenance

- Source Run: `run_20261003T234640Z_3e854d86e9475f64`
- Source commit: `7bdbc8838a2afd281e485e90ed33fafa5d40f6d2`
