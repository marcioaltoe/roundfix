# QA-15 — Roundfix Skill contract

Status: pass.

The canonical Roundfix Skill documents the required slug, text and JSON
formats, `roundfix-specaudit/v1`, all four survivor kinds, classification
evidence, undelivered artifacts, operator-only reclaim commands, and exits
0/1/2. Fresh built help and the clean, attention, invalid-input, and real-Spec
journeys matched that contract.

- `diff -qr .agents/skills/roundfix skills/roundfix` exited 0.
- `rtk make skills-sync-check` passed four tests.
- `rtk make verify` passed Repository Skill Set validation.

The docs correctly state that the audit reports and never reclaims.
