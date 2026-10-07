---
status: done
created_at: 2026-07-29
updated_at: 2026-09-08
absorbed_by: 0121-baseline-decisions-and-complete-regeneration
---

# Repository Skill Set — Doctor demands Roundfix's own development skills in every repository (2026-07-29)

Refreshing four consumer repositories after the `0.0.2` release, three of them failed `roundfix doctor` because the required external skill set includes `golang-cli`, `golang-concurrency`, `golang-context`, `golang-error-handling`, `golang-lint`, `golang-testing`, `bubbletea`, and `tui-design`. None of the three writes Go or renders a terminal UI. The requirement is Roundfix's own development stack leaking into every repository that installs its skills.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-29-doctor-requires-roundfix-own-development-skills.md`.
