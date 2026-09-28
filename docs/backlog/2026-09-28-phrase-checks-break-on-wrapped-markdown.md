---
type: fix
status: open
created: 2026-09-28
spec: null
reason: null
---

# A Verification phrase check fails when the author wraps the Markdown line

## Symptom

Spec 0173's first Run lost task_01 and task_02 after two Verification attempts each although both implementations were correct: their Verifications ran `grep -q "<phrase of four to six words>" <guide>`, and the Agent wrote the phrase wrapped across two lines of `docs/user-guide/commands.md` and `CONTEXT.md`. The diagnostic artifacts were empty, because the honest pattern exits 1 silently on a missing phrase, so neither the Agent nor the operator could see the cause.

## Where

Task authoring (`write-tasks` task template, Verification guidance) and `internal/speccheck/verification.go`.

## Expected

`spec check` flags a single-line `grep -q` of a multi-word phrase against a Markdown file and suggests the wrap-tolerant form `tr -s '[:space:]' ' ' < <file> | grep -qF -- "<phrase>"`; the template teaches that form; a failing phrase check prints which phrase and file it missed.

## Evidence

Run `run_20260928T111838Z_19370d95d9647ce9`, 2026-09-28.
