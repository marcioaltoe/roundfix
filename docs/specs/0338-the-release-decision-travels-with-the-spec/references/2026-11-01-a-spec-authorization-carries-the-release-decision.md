---
type: feat
status: promoted
created: 2026-11-01
spec: 0338-the-release-decision-travels-with-the-spec
reason: null
---

# A Spec's authorization carries the release decision

## Symptom

When the maintainer authorizes release while a Spec is planned, the release
still stops for a second approval. `docs/agents/agent-instructions.md` says a
generic release request authorizes only a conclusive patch plan, and the
Roundfix skill's `references/release.md` says approval of setup,
implementation or an earlier plan never answers `approval_required`. Yet
`docs/agents/spec-routing.md` lists `release` in the closed `operations`
vocabulary of `_authorization.md`, so the record looks able to grant release
while no rule lets it answer the release plan's approval question.

## Decision (maintainer, AskUserQuestion, 2026-11-01)

- When: "Próximo ciclo (Recommended)", with the delivery pull request type
  entry, since both change the release plan's classification.
- Shape: "Teto de impacto (Recommended)". Option text: "Campo com o teto
  (patch, minor ou major). O release plan aceita sem perguntar enquanto o plano
  novo ficar no teto ou abaixo; breaking ou acima do teto pergunta de novo. As
  skills de autoria fazem a pergunta no planejamento quando você pede release."

## Expected

- An approved, non-withdrawn `_authorization.md` with `release` in
  `operations` records an impact ceiling (`patch`, `minor` or `major`) and the
  version-only paths it may change (version file, package manifest, CHANGELOG,
  runbook).
- `roundfix release plan` treats matching records of the Specs in the release
  range as the answer to `approval_required` while the fresh plan's impact is at
  or below every ceiling and nothing is breaking; above a ceiling, breaking, or
  manual classification still stops and asks.
- write-prd, write-techspec and write-tasks ask the release question at
  planning time when the maintainer requests release, and record it.
- The Baseline wording (agent-instructions, release reference, spec-routing)
  changes together and reaches adopters through `baseline update`.

## Evidence

Requested by the maintainer through Argus (`default/gesttione-solutions-argus`,
ai-memory message `01a19920`, 2026-11-01): Spec 0006's `_authorization.md`
(`56eb973`) listed no `release` and said "The release request does not by
itself grant a minor or breaking release"; after PR #23 merged, a new session
had to ask for v0.7.0 again and write a separate record (PR #24,
`docs/references/2026-11-01-autorizacao-do-release-0-7-0.md`).
