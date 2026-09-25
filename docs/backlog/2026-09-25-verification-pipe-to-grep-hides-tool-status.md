---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `SC-VERIFY-INVERTED-EXIT` misses `x="$(tool | grep pat)"; test -z "$x"`

## Symptom

A Verification that pipes a tool into `grep` and tests for emptiness passes when the tool itself fails, and the honesty probe approves it.

## Where

`internal/speccheck/verification.go`; task template Verification note.

## Expected

Recognise the pattern and document the status-preserving form in the template.

## Evidence

secondbrain `inbox/roundfix/2026-09-18-verificacao-com-pipe-para-grep-esconde-o-status-da-ferramenta.md`.
