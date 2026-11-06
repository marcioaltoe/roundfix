---
status: accepted
created_at: 2026-11-06T15:00:00Z
updated_at: 2026-11-06T15:00:00Z
deprecated_at: null
superseded_by: null
---

# The Baseline names the archive destination, a profile's missing decision and a guide it does not own

On 2026-11-06 two adopters reported three Baseline defects after
`roundfix baseline update` 0.136.0. OnionCry found that the docs layout
guide's archive clause still sends a file a later reader needs to "the
Secondbrain inbox", which ADR-0330 retired on 2026-10-25. sql-library found
that a repository-owned Baseline Profile derived from `go-cli-tui` before
`verification.incremental` existed made `baseline update` fail with
`render managed entry "guide.agent-instructions": missing render value for
token "verification.incremental"` while `baseline profile validate` called the
profile valid. It also found `baseline.inventory.nested-carrier-conflict` on
`docs/agents/sql.md`, a guide the repository wrote, which only names the
marker `setup-context-driven` in prose.

Characterization on 2026-11-06 at `9a285d77` showed two facts the reports did
not state. The template `template.guide.agent-instructions` declares the token
`verification.incremental`, which only the decision of the same name binds,
and the core module lists it among its required decisions, but nothing
checks a repository profile's `decisions` against either. The inventory warns
for every guide under `docs/agents/` that carries no managed block, so
`docs/agents/sql.md` warns with or without the marker name; the name matters
only in the colon form `setup-context-driven:`, which makes a carrier
unclassifiable wherever it appears, prose included.

On 2026-11-06 the maintainer chose "Uma Spec de correção (Recommended)"
(option text: "Uma Spec fix: a regra passa a mandar uma mensagem ai-memory ao
default/marcioaltoe-secondbrain. O profile validate e o baseline update nomeiam
a decisão que falta, a linha a acrescentar e o comando --decision. O detector
deixa de casar o nome do marcador citado em prosa.").

## Decision

**The archive clause names an ai-memory message to the project that owns the
Secondbrain; a repository profile that lacks a decision a template it renders
requires is refused with the decision, the line to add and the `--decision`
command; and a guide under `docs/agents/` that the Baseline does not own is
not a carrier conflict.**

- The archive clause (`clause.spec.keep-artifacts-in-spec-folder`) ends: a
  file a later reader needs is copied to `docs/references/`, an accepted ADR,
  the glossary or an agent guide before or at archive, or its knowledge is
  sent in an ai-memory message to the project that owns the Secondbrain, for
  its ingestion. The clause names the destination by its role, as the
  Secondbrain guide does, because every adopter renders it; in the
  maintainer's repositories that project is `default/marcioaltoe-secondbrain`.
  The clause keeps its identity and force, so a Source Baseline transition
  records it as retained.
- A profile decision gap is a decision that binds a token of a template the
  profile renders, that the catalog does not mark optional, and that neither
  the profile's `decisions` nor the decisions they require reach. The
  template must be the only one the catalog can select for its artifact, so a
  template another decision may replace is never counted. `baseline plan`,
  `baseline update` and `baseline profile validate` refuse a profile with a
  gap before they render anything, naming each decision, its template, the
  quoted entry to add to the profile's `decisions` and a `--decision` flag
  with the catalog's suggested value. The refusal never edits the profile.
  Once the profile lists the decision, the update's existing new-decision
  path takes over: it exits 3 and names the decision, as a built-in profile
  does.
- A carrier under `docs/agents/`, other than an `AGENTS.md` or `CLAUDE.md`,
  that has no marker line and whose path is not a path any catalog artifact
  or the recorded Setup Manifest manages is a repository-authored guide. It
  raises no `baseline.inventory.nested-carrier-conflict`. A marker line is a
  `setup-context-driven` begin or end comment that starts a line after at most
  three spaces, the position in which CommonMark 0.31.2, section 4.6, starts
  an HTML block; the marker name inside a code span or a sentence is prose.
  Every other nested carrier keeps its warning, and a nested carrier with an
  unknown real marker stays unclassifiable.

## Considered Options

- Adding the missing decision to the repository profile automatically.
  Rejected: the maintainer chose to name the gap, and a repository profile is
  the repository's own reviewed file, which no command rewrites unasked.
- Checking a repository profile against the module `requiredDecisions` only.
  Rejected as the rule: it names a module but not the template whose render
  fails, which is what the reporter could not find. On every built-in profile
  and on each built-in profile without `verification.incremental` the two
  rules found the same gaps on 2026-11-06.
- Matching the marker name only outside backticks. Rejected: the reported
  guide warns without the name at all, so a narrower text match would leave
  the false warning in place.
- Silencing `nested-carrier-conflict` for every guide under `docs/agents/`.
  Rejected: a guide at a path the Baseline manages, without its markers, is a
  real conflict that the next refresh would meet.

## Consequences

- `baseline profile validate` exits 2 for a repository profile that it called
  valid before, and `baseline plan` and `baseline update` exit 3 with category
  `decision` where they failed with the render error.
- The docs layout guide changes for every adopter at the next
  `baseline update`, as a Managed Refresh.
- A repository-authored guide under `docs/agents/` stops warning in plans and
  updates; nested `AGENTS.md` and `CLAUDE.md` files keep warning, and the
  Baseline still changes no nested carrier.
- The archive decision of 2026-10-06 still names the Secondbrain inbox for a
  cross-project lesson, a sentence ADR-0330 overtook; this decision leaves
  that record's text as it is and states the current destination here.
