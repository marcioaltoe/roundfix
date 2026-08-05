# QA-15 — Roundfix Skill contract

Status: pass.

The canonical Skill was followed against clean and attention fixtures.
Required slug, text, JSON, `roundfix-specaudit/v1`, four survivor kinds,
evidence, operator-only reclaim, and exits 0, 1, and 2 match built help and
observed behavior.

- `diff -qr .agents/skills/roundfix skills/roundfix` — exit 0.
- `rtk make skills-sync-check` — exit 0; four tests passed.
- `rtk make verify` — exit 0, including Repository Skill Set validation.

The Skill accurately says the audit reports and never reclaims.
