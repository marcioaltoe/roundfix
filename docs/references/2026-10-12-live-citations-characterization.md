# Live citations broken by a history sanitize and by an archive (2026-10-12)

A read-only measurement of an adopter repository (the vortex) taken on
2026-10-12 while authoring Spec 0255. Nothing was written to that repository
and none of its content was sent to an external service. Roundfix fixtures stay
synthetic; this document keeps only counts, path classes and commit
identifiers.

## History sanitize through `baseline update`

- History Full Tag: annotated `history-full` at
  `4b5d2cf2b5a2fc637e76b7e2eceb430a4d756b88` (2026-10-02).
- State measured: `544003f1a06597b73505c411a0a03d26e41b04ac^`, the commit
  before the adopter repaired the links by hand (2026-10-09), after
  `roundfix baseline update` 0.65.0 applied 79 Pending History units.
- Method: for every tracked Markdown file outside `docs/history/`, resolve
  each relative Markdown link destination against its file and collect each
  repository-relative `docs/history/...` path; count the targets absent at the
  measured commit and present in the tag tree.
- Result: 46 citations in 22 files. Three name directories rather than
  files (`docs/history/handoffs` and `docs/history/reviews` in a Baseline
  guide, and one converted folder's `qa` directory), so 43 name files. The
  adopter's own scan reported 46 links in 23 files.
- Citing files by class: active Spec folders 5, open Backlog Entries 4,
  unresolved Findings 4, files under `docs/references/` 14, handoffs 17, and
  the Baseline guide 2. The first four classes, 27 citations, are Live
  Documents in the sense of ADR-0260.
- Cited targets by class: a file under a converted folder's `references/` 16;
  its `_prd.md`, `_tasks.md` or a Task file 15; a QA report 4; another file
  under `qa/` or `evidence/` 6; a removed handoff 2; a directory 3.
- The update's JSON plan for that apply carried `warnings: []` and no citation
  field (adopter report, 2026-10-09).

## Secondary Specs linking an owner's `references/`

Measured at `544003f1a06597b73505c411a0a03d26e41b04ac`:

- Two open Specs each hold one relative link from `_prd.md` into
  `../<owner 0078>/references/<file>.md`.
- Two open Specs each hold four relative links (`_prd.md` and three Task
  files) into `../<owner 0095>/references/<file>.md`.
- `roundfix archive` 0.65.0 rewrites only links that leave the archived Spec,
  so archiving either owner first leaves those ten links unresolved.

## Sources

- Adopter reports triaged in the Secondbrain:
  `inbox/roundfix/_triaged/2026-10-09-o-baseline-update-sanitiza-o-historico-sem-relatar-os-links-que-quebra.md`
  and
  `inbox/roundfix/_triaged/2026-10-10-spec-secundaria-linka-a-references-de-outra-spec-e-quebra-no-archive.md`.
- Published practice for moving Markdown files: MD Links for VS Code
  (marketplace.visualstudio.com/items?itemName=TarasErmolenko.fix-refs-vscode)
  "rewrites only the URL token of broken links. Everything else in the document
  stays byte-for-byte unchanged", and offers to update every Markdown link that
  pointed to a moved file; its `fixRefs.warnOnDelete` setting warns "if any
  Markdown links still point to" a deleted file. Zensical Studio
  (zensical.org/docs/studio/workflows/link-intelligence/) updates incoming
  links on a move and also rebases the moved file's own relative links.
