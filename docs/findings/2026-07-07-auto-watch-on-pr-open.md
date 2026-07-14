# Idea: auto-start watch after opening a pull request

Requested 2026-07-07 (Marcio). A config (e.g. `watch.auto_start: true` or a
`pr-opened` hook) so that opening a pull request automatically starts
`roundfix watch --until-clean` for it — closing the gap where review
feedback accumulates until someone remembers to launch the watch loop.
Open questions for the spec: trigger mechanism (gh webhook? polling? a
`roundfix pr create` wrapper? integration with the github-pr-workflow
skill?), detach semantics, and how it composes with Run Outcome
Notifications. **Not scheduled — parked for triage.**
