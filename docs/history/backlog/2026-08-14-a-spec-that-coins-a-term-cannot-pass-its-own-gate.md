---
status: deferred
created_at: 2026-08-14
updated_at: 2026-08-26
---

# A Spec that coins a term cannot pass its own gate

A Spec that declares a Vocabulary Contract promises that its coined term reaches the glossary. The task-authoring contract puts that promise in the closing node: the QA gate checks whether the work introduced, changed, or retired a term and updates the domain context when it found something. The QA gate also refuses before doing anything else when `spec check --strict` reports a finding — and `SC-VOCABULARY-UNDOCUMENTED` is exactly the finding an undocumented coined term produces.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-a-spec-that-coins-a-term-cannot-pass-its-own-gate.md`.
