---
type: fix
status: promoted
created: 2026-10-12
spec: 0264-baseline-group-usage-errors-name-their-own-help
reason: null
---

# A bare baseline group's usage error doubles "baseline" and points at the wrong help

## Symptom

`roundfix baseline capabilities`, `roundfix baseline skills` and
`roundfix baseline assets`, run without their command, exit 2 with
`roundfix: baseline baseline capabilities failed: baseline capabilities requires the check command`
(and the matching `skills` and `assets` text), followed by
`Run 'roundfix baseline profile --help' for usage.` The word `baseline`
appears twice, and the pointer names `baseline profile`, which is not the
group the user typed.

## Where

`internal/cli/baseline_profile.go`: `runBaselineCommand` passes an operation
already prefixed with `baseline` to `printBaselineProfileFailure`, which adds
the prefix again and always prints the `baseline profile --help` pointer.

## Expected

The error names the group once, and the pointer names a help that lists the
group's commands. After Spec 0260, `roundfix baseline <group> --help` prints
the `baseline` usage and exits 0, so the pointer can name it.

## Evidence

Reproduced on 2026-10-12 with `bin/roundfix` built at `7a8b7f70`, during the
help sweep of Spec 0260 (`0260-help-on-every-command-path-exits-zero`), which
fixes the help exit codes and leaves this error text unchanged.
