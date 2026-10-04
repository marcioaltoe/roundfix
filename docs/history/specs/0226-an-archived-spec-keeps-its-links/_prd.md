---
spec: 0226-an-archived-spec-keeps-its-links
status: archived
created: 2026-10-04
surfaces: [backend, cli, docs]
archived: "2026-10-04"
source_slug: 0226-an-archived-spec-keeps-its-links
---


# An archived Spec keeps its links

`roundfix archive` moves a Spec one directory deeper, from the Spec Root to its
archive root. A relative Markdown link in the Spec that leaves it, such as a
PRD's link to a Finding or a Task's link to an ADR, then resolves one level
short and breaks. fluxus measured it on 2026-10-02: Spec 0096 reached `main`
with links that can no longer be fixed in place, because an archived Spec is
immutable, and Specs 0097 and 0098 needed three links fixed by hand in their
archive commits. Roundfix's own history carries the same damage. This Spec
makes the archive keep every outward link reaching the file it reached, and
refuse, naming the link, when one reaches nothing.

## Prerequisites

None. Specs 0222, 0223 and 0224 are authored in the same cycle and may also
raise the Roundfix Skill's version; the operator orders the queue so that each
later Spec raises it from the earlier one's value, and task_01 raises it from
the value on the tree it starts from. The governed test file this Spec changes
was granted by name on 2026-10-04 (`_authorization.md`).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; Specs
  keep their slugs and archive destinations. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the archive reads and rewrites
  local files only; no request, credential or forge read is added. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0230 (this Spec) decides that the
  archive rewrites relative links that leave the Spec and refuses when one
  resolves from neither location. ADR-0223 keeps a delivery's work across
  archive, requeue and review, and the Delivery Queue's resumed archive commit
  accepts the rewrites, ADR-0223: "A delivery keeps its work across archive,
  requeue and review". ADR-0187 splits the Roundfix Skill by command and
  ADR-0189 ties an owned skill's version to its content, so the skill edit
  raises its version. ADR-0184: "A TechSpec now declares numbered Surface
  Transcripts", applied to the archive refusal. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and ADR-0093,
  ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by
  citation and receipt. ADR-0178 authorizes each Task commit by its grant,
  ADR-0182 runs Settlement Checks before it, and ADR-0166 records undeclared
  paths; every Task declares its paths. ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park, which this Spec does not touch; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot, and ADR-0192 cites ADR-0178 but decides how a conflict in declared derived paths is resolved; this Spec changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-10-04: the named grant "Concedo" for `internal/spec/archive.go`, and the
  standing skills authorization "considere autorizado a ajustar todas as
  skills se necessário", and the named grant "Concedo" of 2026-10-04 for the
  governed test `internal/spec/archive_test.go`, limited to the Spec 0058
  replay fixture. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0226-an-archived-spec-keeps-its-links/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `skills/roundfix/SKILL.md`.

## Goals

- Every relative Markdown link that leaves an archived Spec reaches the file it
  reached before the archive.
- No Spec reaches history with an outward link that reaches nothing.
- A Delivery Queue item whose archive rewrote links resumes and merges as it
  did before.

## User Stories

1. As a Supervisor archiving a Spec, I want relative links that leave the Spec
   rewritten to keep reaching their files, so that the archived record stays
   navigable.
2. As a Supervisor, I want the archive to refuse, naming each link, when a link
   that leaves the Spec reaches nothing, so that a broken link never reaches
   history where it cannot be fixed in place.
3. As the operator of the Delivery Queue, I want an item whose archive commit
   carries link rewrites to resume without parking, so that the fix costs no
   intervention.

## Core Features

1. **Outward links are rewritten.** Before it moves a Spec, the archive reads
   every relative link destination in the Spec's Markdown files: inline links,
   images and reference definitions, outside fenced code blocks and code
   spans. A destination that leaves the Spec and reaches an existing path from
   the active location is rewritten to the relative path that reaches the same
   path from the archived location, keeping its fragment, query and
   angle-bracket form. This holds for the built-in Spec Root and any
   configured one, and for a superseded Spec.
2. **Already archived targets keep their link.** A destination that reaches
   nothing from the active location but reaches an existing path unchanged from
   the archived location keeps its bytes.
3. **Everything else keeps its bytes.** Destinations inside the Spec, absolute
   destinations, fragments, URLs, HTML anchors and non-Markdown files are not
   touched.
4. **A broken outward link refuses the archive.** When a destination that
   leaves the Spec reaches nothing from either location, the archive exits `2`
   before stamping or moving anything and names every such link with its
   Spec-relative file, line and destination.
5. **The command says what it rewrote.** A successful archive that rewrote
   links appends the count to its confirmation line.
6. **A resumed delivery accepts the rewrites.** The Delivery Queue treats an
   archive commit whose content changes are the archive stamp and these
   rewrites as the exact archive it expects.
7. **The skill and the guide say so.** The Roundfix Skill's archive reference
   and the archive command guide describe the rewrite, the refusal, the count
   and the limits.

## User Experience

`roundfix archive <slug>` prints `archived <slug> -> <path>` as before, with
`; rewrote <n> relative link(s)` appended when it rewrote any. When it refuses,
standard error lists each broken outward link as `<file>:<line> "<destination>"`
and nothing in the Spec changes.

## Non-Goals / Out of Scope

- Evidence scripts or other non-Markdown files that climb a fixed number of
  directories, HTML anchors, wiki links and absolute destinations; the archive
  guide names this limit.
- Repairing links in Specs archived before this Spec; archived Specs stay
  byte-identical.
- Rewriting links in other files when their target is archived later, and
  links from outside the Spec that point into it.
- Changing archive eligibility, the QA override, the archive stamp or the
  archive destinations.

## Success Metrics

1. Success Metric: archiving a Spec whose files link a Finding, an ADR, a
   sibling Spec, a reference definition and an angle-bracket destination leaves
   every such link reaching the same file from the archived location and
   changes no other byte.
2. Success Metric: archiving a Spec with one outward link that reaches nothing
   exits `2`, names the file, line and destination, and leaves every file of
   the Spec and its location unchanged.
3. Success Metric: a Delivery Queue item whose archive commit rewrote links is
   accepted as an exact archive on resume, and one with any other content change
   is still refused.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- fluxus's measurement of 2026-10-02 in the Secondbrain Inbox Entry
  `inbox/roundfix/_triaged/2026-10-02-o-archive-nao-reescreve-links-relativos.md`
  (Specs 0096, 0097 and 0098), and fluxus's own Spec 0067, which on 2026-09-03
  repaired forty broken documentation link destinations by resolving each
  against the file that holds it.
- A measurement during authoring on 2026-10-04 over this repository's
  archived Specs: 43 of 204 contain 134 outward relative links; 52 are broken
  by the move but resolve from the active location, 78 resolve only because
  their target was archived too, and 4 resolve from neither.
- The CommonMark specification (<https://spec.commonmark.org/0.31.2/>),
  section 4.7: "A link reference definition consists of a link label, ...
  a link destination, ... and an optional link title", which fixes the link
  forms the archive reads.
- The published `docmv` tool (<https://pypi.org/project/docmv/>), which moves
  Markdown files by recomputing "the relative path from its post-move
  container to its post-move target" and whose `apply` "refuses and changes
  nothing" when "the move would introduce any broken link".

## Decisions

- Rewrite and refuse, including links already broken before the archive. See
  ADR-0230.
- Keep a link whose target was already archived into the mirrored history
  layout, rather than point it back at an active path that no longer exists.
- Split from Spec 0223, which was authored from the same week's adopter
  requests, because this change needs a grant beyond the one named for the
  archive source file and adds a Delivery Queue change.

## Open Questions

- None. The named grant for `internal/spec/archive_test.go`, which the replay
  test `TestSpec0058ReplayArchivesDeclaredUnreachableRelease` needs because its
  three outward links reach nothing in an empty temporary repository, was
  given on 2026-10-04 ("Concedo"); the change creates the three link targets in
  that temporary repository and leaves the assertions unchanged.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
