---
task: task_02
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: pending
type: backend
complexity: low
---

# Task 02: The Project Config reads the verification.format command

## Overview

Adds the optional Format Command to the configuration that task_03 runs. This
Task answers the Backlog Entry "An unformatted QA report breaks the next Run's
precondition" of 2026-10-06, whose first expected outcome asks how the
formatter is discovered: by an explicit Project Config command (ADR-0249). It
is verifiable on its own through the config package's tests.

## Requirements

1. MUST add `Format string` to the `Verification` configuration and accept the
   key `format` in the `verification` mapping of User Config and Project
   Config. A value that is not a YAML string scalar MUST be refused with an
   error containing `verification.format must be a string`
   (`_techspec.md` → Invariant 1).
2. MUST default the value to empty, treat a whitespace-only value as empty, and
   let a Project Config value replace a User Config value, as the other
   `verification` keys do.
3. MUST render `format: ""` under `verification:` in the default configuration
   that `roundfix init` writes, with a one-line comment saying that the command
   formats the files a QA step commits under the Spec's `qa/` directory and
   that empty runs nothing (`_techspec.md` → API Contract 1).
4. MUST add `internal/config/verification_format_test.go` with the tests
   `TestVerificationFormatDefaultsEmpty`,
   `TestVerificationFormatProjectReplacesUser`,
   `TestVerificationFormatRejectsANonString` and
   `TestVerificationFormatRendersInTheDefaultConfig`, each driving the
   package's public load or render entry point.
5. MUST NOT change any other configuration key, its default or its error text,
   and MUST NOT wire the value into a Run; task_03 does that.

## Subtasks

- [ ] Add the field, the key and its validation.
- [ ] Merge Project over User and trim whitespace.
- [ ] Render the default with its comment.
- [ ] Write the four tests.

## Acceptance Criteria

- [ ] An unset key loads as empty, and a Project value replaces a User value.
- [ ] A list or mapping value is refused with the named error.
- [ ] The rendered default configuration contains the key under
      `verification:`.
- [ ] Every existing config test passes unchanged.

## Context

- interface: `internal/config/config.go`
- creates: `internal/config/verification_format_test.go`
- instruction: `internal/config/verification_tools_test.go`
- instruction: `docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`

## Verification

- `out="$(go test -count=1 -v -run '^TestVerificationFormat' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestVerificationFormatDefaultsEmpty TestVerificationFormatProjectReplacesUser TestVerificationFormatRejectsANonString TestVerificationFormatRendersInTheDefaultConfig; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0; before this Task no test of that name exists, so the first name is missing, and after it all four pass.

## References

- `_prd.md` → Core Feature 1; Goals
- `_techspec.md` → API Contract 1; Invariant 1; Build Order 2
- ADR-0249
