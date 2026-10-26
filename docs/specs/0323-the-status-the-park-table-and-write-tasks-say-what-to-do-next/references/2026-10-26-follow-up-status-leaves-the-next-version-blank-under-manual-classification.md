---
type: fix
status: promoted
created: 2026-10-26
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
reason: null
---

# Follow-up status leaves the next version blank under manual classification

## Symptom

With the release Pull Request open and an unclassified `docs:` commit on the
default branch, `roundfix follow-up --status` printed

```text
Next release:  from v0.117.0 (manual_classification_required)
```

with two spaces where the version belongs, and listed
`release : manual classification required` under `Waiting on the
maintainer`, again with an empty version before the colon. Neither line says
why there is no version or what the maintainer runs to settle it. The JSON
report has `release.proposed_version: ""`, `release.reason: null` and a
`release_classification` maintainer item whose `detail` is only
`manual classification required`.

## Where

`internal/cli/follow_up_status.go`: `runFollowUpStatus` copies the Release
Plan into `followUpRelease` and adds the maintainer item;
`renderFollowUpStatusText` prints `Next release: %s from %s (%s)` for every
state other than `no_unreleased_changes`, `no_release` and `unavailable`,
and `release %s: %s` for every maintainer item that is not a pending
question. `releaseplan.CalculateProposal` proposes no version in
`manual_classification_required`, so both version fields print empty.

## Expected

The manual classification state renders a clear sentence and no blank field,
in text and JSON: how many commits since the base tag need a manual release
classification, and the `roundfix release plan` command that classifies
them. The other release states keep their bytes.

## Evidence

Operator observation on the Mac mini delivery machine, 2026-10-24 ~11:22 BRT
(Roundfix 0.117.0), recorded as item 11 of the cycle AO triage inputs:
"follow-up --status: com o release PR aberto e um commit docs: sem
classificação, mostra 'Next release:  from v0.117.0
(manual_classification_required)' (versão em branco) e põe 'release: manual
classification required' em Waiting on the maintainer. Rever a renderização
do caso manual." Reproduced at `aa8b3db7` by
`TestFollowUpStatusReportsTheNextRelease/manual`, whose fixture has the same
shape and whose assertions accept the blank version.
