---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# History sanitize writes Archive Records that lose their summary, links and dates

## Symptom

`roundfix history sanitize --apply --batch 26 --advise` (0.69.0) in fluxus wrote
Archive Records that the pre-PR review then flagged, across three batches:

1. A multi-line blockquote in the legacy `_prd.md` is joined into one summary
   line with `>` left mid-text (`0009-deployed-revision-and-publication-signal`).
2. A legacy `_prd.md` that starts straight at `## Problem statement` yields a
   record with the title and no summary (`0001`, `0004`).
3. A legacy `qa_override: true` without approval or reason yields empty
   `qa_override_approval` and `qa_override_reason`, which a review reads as an
   unauthorized override (0003, 0010, 0011, 0013–0017, 0019).
4. Summary links `references/<file>.md` are copied verbatim, so from
   `docs/history/specs/` they point at a folder that no longer exists, although
   the same command promoted the file to `docs/references/` (0040, 0041); links
   inside promoted references also keep pointing at the old Spec folder.
5. A record promotes an adopted reference and still declares `sources: []`
   (0068–0072).
6. `created: ""` when `_prd.md` has no date, although `_techspec.md` in the same
   folder has one (0063, 0066–0072).

## Where

The Archive Record generator used by `roundfix history sanitize` and
`roundfix baseline update` for Legacy Archive Folders (Lenient Legacy Reading).

## Expected

The summary keeps blockquotes readable and falls back to the first paragraph of
the first section; a legacy override records an explicit value such as
`legacy: not recorded`; summary links to a promoted file are rewritten to its
promoted copy (or to `source_revision`); `sources` lists what the record
promotes; `created` falls back to the TechSpec date or the first commit.

## Evidence

Secondbrain inbox note from fluxus, 2026-10-12
(`inbox/roundfix/2026-10-12-history-sanitize-colapsa-blockquote-e-perde-resumo.md`);
fluxus dismissed the review findings with evidence because the record is
generated.
