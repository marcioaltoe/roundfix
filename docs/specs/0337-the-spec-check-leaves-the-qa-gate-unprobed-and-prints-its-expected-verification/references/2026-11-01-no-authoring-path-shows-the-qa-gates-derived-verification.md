---
type: fix
status: promoted
created: 2026-11-01
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
reason: null
---

# No authoring path shows the QA gate's derived Verification

## Symptom

A QA Task written from the installed `write-tasks` template fails
`SC-QA-VERIFICATION-AUTHORED`, and nothing the CLI prints tells the author
what to write instead. The template teaches an authored `## Verification`
for every Task and has no instruction for the QA gate's derived one. The
finding's fix says only "Regenerate the qa Task so its ## Verification
renders Roundfix's derived command unchanged.", and no command renders that
text.

## Where

- `.agents/skills/write-tasks/references/task-template.md` (installed in
  adopters as `.agents/skills/write-tasks/references/task-template.md`) and
  `.agents/skills/write-tasks/SKILL.md`.
- `internal/speccheck/coherence.go`, `AuthoredQAVerification`, which builds
  the finding with code, summary, location and fix only.
- `internal/spec/task.go`, `DerivedQAVerification` and `taskVerification`,
  which hold the only rendering of the expected command.

## Expected

The `SC-QA-VERIFICATION-AUTHORED` finding prints the exact expected
Verification block, in text and in JSON, and the `write-tasks` template and
`SKILL.md` teach copying it from there. No new command.

## Evidence

Reported by the Quaestor project (`default/gesttione-solutions-quaestor`)
through ai-memory message `01a19861` on 2026-11-01: their QA Task written from
the template failed `SC-QA-VERIFICATION-AUTHORED`, and they had to read the
`DerivedQAVerification` and `taskVerification` source to render the derived
text. The maintainer chose the Expected shape on 2026-11-01 ("spec check
mostra (Recommended)"). The operator reproduced the finding on `340c9004`
with that commit's `bin/roundfix`: a pending QA Task whose Verification is
`test -f docs/specs/9001-fixture/qa/report.md` gives one
`SC-QA-VERIFICATION-AUTHORED` whose text and JSON carry no expected text.
