---
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: active
created: 2026-10-13
surfaces: [backend, cli, docs]
---

# History sanitize records keep their summary, links and dates

`roundfix history sanitize` and the history section of `roundfix baseline
update` turn each Legacy Archive Folder into an Archive Record. On 2026-10-12
an adopter's three batches produced records that its pre-PR review flagged in
six ways: a blockquote summary joined with `>` mid-text, a missing summary, an
override with empty approval and reason, summary links and promoted copies
that point at the removed folder, `sources: []` beside a promoted adopted
source, and `created: ""` beside a dated TechSpec
([the Backlog Entry](references/2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md)).
This bug fix makes the generator keep what the legacy folder held. Records
already written stay byte-identical, and only generation from Legacy Archive
Folders changes.

## Prerequisites

None. Specs 0266 to 0270 of cycle AG are authored at the same time and touch
other files; no Verification here depends on another Spec's artifacts. Spec
0270 adds a `spec check` rule for QA byte invariants pinned to a commit; this
Spec's QA compares with the delivery base, so the rule does not apply.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or Archive Record
  field is added; the schema stays `roundfix/archive-record/v1`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; Git is read locally and tests run on temporary repositories
  and homes. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0277 (this Spec) decides the five
  readings. ADR-0251: "The rule applies wherever a Legacy Archive Folder
  becomes a record", so the History Sanitize Command and the folder branch of
  `ReadArchivedSpec` read alike. ADR-0248: "Each folder becomes `<slug>.md`
  through `BuildArchiveRecord`, the builder the archive itself uses"; the
  Archive Command keeps its reading. ADR-0247: "Promoted copies must equal
  removed blobs" stays the Delivery Queue's proof of an archive. ADR-0254:
  "One plan, one digest", so the update's digest binds the changed records.
  ADR-0260: "a promoted file keeps its bytes and its own outward relative
  links are not rebased", which ADR-0277 narrows for conversions. ADR-0261:
  "records already archived stay byte-identical". ADR-0272 reads a
  `source_revision` the way ADR-0261 recovers it and changes no record. ADR-0262: "A cancelled Spec
  leaves an Archive Record and returns its sources"; a cancellation's record is
  not a conversion and keeps its reading. ADR-0269: "Archive and cancel refuse
  ignored Spec files and finish a failed removal forward"; neither command
  changes. ADR-0256: "A release waits
  for the skills that describe a changed surface". ADR-0244: "A Spec declares
  the domain terms it introduces". ADR-0257: "Inside a Run the Daemon is the
  only full gate". ADR-0258: "A Spec is authored against the QA rerun
  classes". ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts"; the new tests
  wait on no clock. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AG grant of 2026-10-13 (the answers to the AskUserQuestion of 2026-10-13
  and "Pode iniciar", under the cycle AE "Concedo" for the Governed Paths each
  Spec declares) and the 2026-09-30 skills grant ("considere autorizado a
  ajustar todas as skills se necessário"); bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, lint, formatter or module
  file change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0271-history-sanitize-records-keep-their-summary-links-and-dates/_authorization.md`.

## Goals

- A record converted from a Legacy Archive Folder carries a readable summary
  whenever the PRD or QA Report holds a paragraph.
- No link the generator writes in a record or a promoted copy points at the
  removed folder.
- A record carries the override, sources and creation date the folder can
  supply, and says so explicitly when the folder held no override approval or
  reason.

## User Stories

1. As an adopter sanitizing history, I want each record's summary to read as
   the legacy author wrote it, so that a reviewer reads the outcome instead of
   flagging stray `>` or an empty body.
2. As an adopter promoting a reference, I want the record's summary and the
   promoted copy to link to files that exist, so that I do not fix links by
   hand in the batch.
3. As a reviewer of a sanitize batch, I want override, sources and creation
   fields that say what the folder recorded, so that a generated record does
   not read as an unauthorized override or a broken trace.

## Core Features

1. Summary. Block quote markers are removed from each line before the
   summary paragraph is joined, and a PRD whose title is followed directly by
   a section yields the first paragraph after its first second-level heading.
2. Fields. A legacy override without approval or reason records
   `legacy: not recorded` there; `sources` falls back to the adopted source
   files under `references/` when the index lists none; `created` falls back
   to the TechSpec date and then to the oldest commit that added the PRD.
3. Links. Legacy Link Rebase rewrites the relative destinations of the
   record's outcome and of each promoted copy, and the folder path text in
   them.
4. Documentation. The `history` user guide, the Roundfix skill's archive and
   baseline references and `CONTEXT.md` describe the new readings.

## Non-Goals / Out of Scope

- Rewriting or backfilling an Archive Record already written, here or in any
  adopter.
- Changing the Archive Command, its exact-archive proof or the outcome of an
  active Spec's record; the sibling defect is the open Backlog Entry
  `docs/backlog/2026-10-13-an-archive-record-outcome-keeps-links-relative-to-the-removed-spec-folder.md`.
- Changing the review conventions, the Archive Record schema or the Plan
  Digest format.
- Repairing a link that did not resolve before the conversion; it keeps
  naming the same repository path.

## Success Metrics

1. Success Metric: a synthetic Legacy Archive Folder for each of the six
   defects converts to a record without the defect, through the built
   `bin/roundfix history sanitize --apply`.
2. Success Metric: a promoted copy without relative links and a record from a
   folder without any of the six shapes are byte-identical to what the
   delivery base's binary writes for the same fixture.
3. Success Metric: every Archive Record present at the delivery base is
   byte-identical at the audited head, and the Archive Command's existing
   tests pass unchanged.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-13:

- CommonMark Spec 0.31.2, section 5.1 Block quotes
  (https://spec.commonmark.org/0.31.2/): "A block quote marker consists of
  0-3 spaces of initial indent, plus (a) the character `>` together with a
  following space, or (b) a single character `>` not followed by a space."
  This is the marker the summary reader removes, and why the joined text keeps
  no `>`.
- RFC 3986, section 4.2 (https://www.rfc-editor.org/rfc/rfc3986.html): "A
  relative reference ... refers to a resource by describing the difference
  within a hierarchical name space between the reference context and the
  target URI." A relative link therefore depends on the file that holds it,
  which is why a summary or a copy that moves must be rebased.

The adopter's records were characterized read-only in a scratch clone and are
not quoted here. The Secondbrain was consulted through `wiki/index.md` and
`qmd query "roundfix history sanitize archive record summary legacy folder"`;
it returned the mirrors of ADR-0251
(`projects/roundfix/mirror/docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`),
ADR-0254 and Spec 0226 (`projects/roundfix/mirror/docs/history/specs/0226-an-archived-spec-keeps-its-links.md`),
which confirmed that Lenient Legacy Reading is the owning rule and that links
were last handled for folders that moved whole. The original note is the
Secondbrain inbox entry
`inbox/roundfix/2026-10-12-history-sanitize-colapsa-blockquote-e-perde-resumo.md`,
already triaged into the adopted Backlog Entry. The CommonMark and RFC passages
are standard references, not research, so no Inbox Entry was captured.

## Glossary

- adds: **Legacy Link Rebase**
- changes: **Lenient Legacy Reading**

## Decisions

- The five readings apply to Legacy Archive Folders only; see ADR-0277.
- Missing override approval or reason records `legacy: not recorded`; see
  ADR-0277.
- Promoted copies of a conversion are rebased, narrowing ADR-0260; see
  ADR-0277.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the characterization, the contracts
and the build order.
