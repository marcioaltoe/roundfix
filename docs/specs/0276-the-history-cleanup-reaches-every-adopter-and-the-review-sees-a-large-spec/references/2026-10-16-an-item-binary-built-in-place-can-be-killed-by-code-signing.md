---
type: fix
status: promoted
created: 2026-10-16
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
reason: null
---

# An item binary built in place can be killed by code signing

## Symptom

On macOS, a signed executable rewritten in place while another process still
maps it is killed with `SIGKILL` by the kernel's code-signing check, and the
crash report under `~/Library/Logs/DiagnosticReports/` carries the namespace
`CODESIGNING`. Measured on 2026-10-16 while authoring Spec 0275: 40 kills in
40 tries when the file was overwritten with `cp` while a copy ran; none in
3,400 tries when nothing mapped it; none when the build copied to a new file
and renamed it over the old one.

## Where

The Delivery Queue's item binary build, `delivery.item_binary` in the Project
Config, and any adopter build command that writes its output in place (for
example `cp` over `bin/roundfix`) while an earlier item binary still runs, such
as a Run owner the implement step started.

## Expected

The item binary is replaced atomically (write to a sibling path, then rename
over the target), or the Delivery Queue refuses or warns about an in-place
build, so a running item binary is never rewritten under itself.

## Evidence

Spec 0275 (`docs/history/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced.md`)
fixes the same mechanism in the delivery test fixtures and leaves production
builds unchanged; its authoring probes reproduced F6's notice "was killed by
signal SIGKILL (9), no core dump" five times in five through the fixture.
Sources: Apple, "Updating Mac Software"
(https://developer.apple.com/documentation/security/updating-mac-software),
Open Radar FB8914243 (https://openradar.appspot.com/FB8914243).
