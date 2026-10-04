---
type: feat
status: promoted
created: 2026-09-30
spec: 0223-adjustments-the-adopters-asked-for
reason: null
---

# The release plan does not report the skill and guide checks

## Problem

Spec 0195 makes the skills-and-guides check a mandatory step of the release runbook. The step is a list of commands a maintainer or an Agent runs. `roundfix release plan` is the first command of every release and is read-only, but it says nothing about that step: it does not report whether the owned skills match the versions the binary carries, or whether the rendered guides match the Baseline modules.

## Expected

`roundfix release plan` prints one `skills:` line and one `baseline:` line, in text and in JSON, computed from local files only. It stays read-only and contacts no service, as its contract requires. A failing line does not change the proposed version; it names the runbook step.

## Why it was not done in Spec 0195

The Release Plan Command has a frozen JSON contract and its own documentation contract tests, which are Governed Paths. Adding two lines is a Spec of its own and did not fit four implementation Tasks.
