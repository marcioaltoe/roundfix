---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Roundfix detects and avoids an unsafe codex

macOS XProtect intermittently blocks the OpenAI-signed but un-notarized codex
binary when an agent loop execs it frequently, killing Runs mid-flight —
especially when the binary carries the Homebrew-Cask `com.apple.quarantine`
attribute. Roundfix therefore treats codex hygiene as a Run prerequisite: the
Doctor Command inspects the codex resolved on PATH (and any configured codex
path) for quarantine and notarization acceptance and, on failure, names the
curl-to-`~/.local/bin` reinstall fix, while codex spawning through `codex-acp`
prefers a verified-clean binary so repeated execs stop tripping XProtect.
Detection is macOS-only — other platforms report the check not-applicable — and
Roundfix only diagnoses and instructs; it never reinstalls or de-quarantines
codex automatically.
