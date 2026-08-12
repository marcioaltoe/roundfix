---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Agent run logs are opt-in

Per-Batch agent log files are no longer written by default: the Run Event Journal already stores every raw agent payload durably (ADR-0008), making the log files a redundant on-disk copy that accumulated tens of megabytes per day of dogfooding. A config key turns them back on for development or debugging. The one exception is a Detached Run's console log, which stays unconditional because it is the caller-visible record ADR-0028 promises.
