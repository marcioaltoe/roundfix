---
type: fix
status: open
created: 2026-11-13
spec: null
reason: null
---

# A completion notice reads an adopted record as closed

## Symptom

Reported by argus in ai-memory message `01a1d483-25a4-71c3-91d7-86600602a109` of 2026-11-13, from its pre-PR review of the Baseline refresh to Roundfix 0.143.0 (roundfix skill 0.1.228); it asks for no reply. The "A record is closed without a Spec" rule of `skills/roundfix/references/release.md` (`## Completion notices`, lines 145-153 in 0.1.228) treats every path that `git diff --name-only --diff-filter=D <merge commit>^ <merge commit> -- docs/backlog docs/findings` lists as a closed record, without reading why the file left. Under the adoption contract of `docs/agents/docs-layout.md`, adopting a Backlog Entry or Finding moves it to `docs/specs/<slug>/references/` with `status: promoted` (`spec` set). Git reports that move as a deletion whenever rename pairing fails: an edit too large for similarity, rename detection off, or `diff.renameLimit` exceeded. The merge that only plans the Spec would then send the Notice Request's one notice early, as "delivered without a Spec" or "closed", and the Spec's archive would later find another event and send a second one.

The operator confirms the cause: the round-1 fix of Spec 0386 (`3ed725d7`) said that a record a PRD's `## Sources` names was adopted and its notice follows the Spec; the round-2 rewrite (`a7c2b447`) dropped that sentence, and v0.143.0 shipped without it.

## Where

`skills/roundfix/references/release.md`, `## Completion notices`, the closure bullet, and its canonical copy under `.agents/skills/`; `skills/roundfix_modes_test.go` pins the section's phrases.

## Expected

The closure rule reads the deleted record's frontmatter at `<merge commit>^` and sends a notice only for a terminal disposition (`declined`, `deferred`, `done` without a consuming Spec, `closed`, `cancelled`, `deprecated`, `superseded`); a `promoted` record, or one whose path the same merge adds again under the Spec root, was adopted and its notice follows the Spec. A contract test pins the exclusion.

## Evidence

Sender's evidence: gesttione-solutions/argus branch `chore/baseline-roundfix-0-143-0`, `roundfix review` finding F1 of 2026-11-13 on `.agents/skills/roundfix/references/release.md:145-153`. Confirmed by reading `release.md` at `777cd521` (v0.143.0) and the 0386 item commits `3ed725d7` and `a7c2b447`.
