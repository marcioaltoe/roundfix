---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-12T12:59:56Z
updated_at: 2026-08-12T12:59:56Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# An archive move is a ledger of identities, never of content

The Baseline Plan already expresses a relocation without a new concept — a
postimage absent at the source and a postimage carrying the bytes at the
destination — but a postimage holds full content and the plan serializes to JSON,
so relocating an archive that only ever grows would put its every byte into the
plan document. Archive relocations are therefore recorded as their own ordered
ledger of source, destination, and content identity, carrying no bytes, and the
apply verifies the destination's identity against the recorded one. The trade-off
accepted is a second mutation shape inside a plan contract that already carries
one, bought against a plan whose size would otherwise grow with the history it
moves rather than with the change it makes.
