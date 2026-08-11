---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Configured fallbacks activate after notification

Roundfix proves the Preferred Selection and every configured Fallback Selection needed by the requested Run through disposable Agent Sessions before creating the Run. If a live selection fails before its first Agent prompt begins, Roundfix emits a structured notification to the user and Supervisor and then automatically starts the next proven Fallback Selection, including a different ACP Runtime when configured; once Agent work begins, Roundfix fails the Work Item instead of switching models over potentially modified state. This supersedes ADR-0041: dynamic catalog probing and human confirmation at failure time are replaced by explicit configuration, preflight proof, notification, and automatic activation.
