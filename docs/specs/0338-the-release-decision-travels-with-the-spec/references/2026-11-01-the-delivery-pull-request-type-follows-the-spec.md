---
type: fix
status: promoted
created: 2026-11-01
spec: 0338-the-release-decision-travels-with-the-spec
reason: null
---

# The delivery pull request type follows the Spec

## Symptom

The Delivery Queue titles every pull request `feat: deliver <slug>`
(`internal/cli/deliver_workflow.go`, the `title := "feat: deliver " + specSlug`
assignment and the numbered-prefix variant below it). The squash merge keeps
that subject, so `roundfix release plan` classifies every queue-delivered Spec
as a feature and sets the automatic minimum impact to minor. A fix or refactor
Spec can never ship as a patch.

## Where

`internal/cli/deliver_workflow.go` (pull request title) and the release plan's
change classification.

## Expected

The pull request and squash type come from the Spec, not a constant: `fix:`
or `refactor:` for a Spec authored through the refactor/bug-fix route or whose
adopted sources are fixes, `feat:` for a feature Spec. The type must be
recorded somewhere the release plan can read without guessing (a declared
field in the plan, or the Archive Record's source type). The decision on the
mechanism belongs to the Spec's TechSpec.

## Evidence

Reported by Argus (`default/gesttione-solutions-argus`, ai-memory message
`01a193d4`, 2026-10-31, roundfix 0.126.0): Spec
`0002-instaladores-no-caminho-certo-do-spaces`, a pure fix with a minimal PRD
and one `type: infra` Task, merged as PR #10 with squash `eaa337a`
`feat: deliver 0002-...`; `roundfix release plan --impact patch` then refused
with `manual impact "patch" is invalid for automatic minimum "minor"`. The
maintainer had approved 0.3.1 and had to accept 0.4.0. The operator confirmed
the constant on `340c9004`.
