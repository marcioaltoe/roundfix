---
schema: roundfix/archive-record/v1
spec: 0013-codex-runtime-hygiene
title: Codex Runtime Hygiene
status: archived
created: "2026-07-06"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0013-codex-runtime-hygiene
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0032
sources: []
regeneration: []
promoted: []
---

# Codex Runtime Hygiene

macOS Tahoe blocked Roundfix Runs mid-flight with "Malware Blocked: 'codex' was not opened because it contains malware." The codex binary is OpenAI-signed but not Apple-notarized, so XProtect YARA-scans the ~193 MB binary on every exec and intermittently blocks it — worst when the binary carries the Homebrew-Cask `com.apple.quarantine` flag and an agent loop spawns it many times a minute. The environment was repaired by hand (curl install to `~/.local/bin` without quarantine, `CODEX_PATH` pinned in the shell), but no product should depend on a hand-repaired machine. This Spec teaches Roundfix to detect an unsafe codex on PATH and to spawn a verified-clean one, so an agent loop never trips XProtect.
