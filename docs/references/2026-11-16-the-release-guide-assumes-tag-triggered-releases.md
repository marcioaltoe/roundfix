---
type: fix
status: promoted
created: 2026-11-16
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
reason: null
---

# The release guide assumes tag-triggered releases

This entry is the release-guide half of the Backlog Entry `docs/backlog/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md`. The author split it from that entry on 2026-11-16 when Spec 0400 adopted it, following the maintainer's decision below. The Delivery Queue half stays open in that Backlog Entry.

## Symptom

The Quaestor wrote on 2026-11-16 (ai-memory message `01a1e522-a608-7563-8568-289c6138d377`, information only, no reply asked) that its Spec 0115, under its own decision records 0115 and 0116, publishes the tag and the release on the push to `main` when `project.version` changes, with no manual tag. Its second observation about Roundfix, in its words:

2. "`.agents/skills/roundfix/references/release.md` diz \"Preserve the existing tag-triggered workflow\", o que supõe que todo repositório adotante publica release por tag. O Quaestor passa a publicar no merge da versão (ADR 0116 [do Quaestor])". The bracketed words and the dropped hyphen are the author's, so the Quaestor's record number is not read as one of this repository's ADRs.

## Decision

Marcio, 2026-11-16 ~20:55 BRT: "Spec fix depois da 0399 (Recommended)" — one release-plan Spec that adopts the Oraculum finding and this half; the Delivery Queue half stays in the backlog.

## Where

The release guide (`references/release.md`) of the Roundfix skill.

## Expected

The release guide describes a release published on the merge of the version change as well as one triggered by a tag, and tells the session to keep the trigger the repository already uses.

## Evidence

The Quaestor's message above. Roundfix's own release workflow is tag-triggered (`.github/workflows/release.yml`, `on: push: tags: v*`), so it does not reproduce here.
