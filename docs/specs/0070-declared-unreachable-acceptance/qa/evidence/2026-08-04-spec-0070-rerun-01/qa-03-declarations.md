# QA-03 — Declaration reader and diagnostics

Status: pass

- Present: the declared-only public journey consumed the author-supplied
  criterion, reason, and `satisfied-by` action, archived, and persisted that
  action under `unproven`.
- Absent: the ordinary pass fixture had no declaration section, archived with
  exit 0, and gained no `unproven` field.
- Malformed: the partial fixture omitted `reason`; the built Archive Command
  exited 2 and named `_prd.md`, line 9, and `missing reason`. A fresh read
  retained the active Spec and found no archive destination.

The full gate also passed the declaration parser's two-entry, absent-section,
and all malformed-shape tests on the current build.
