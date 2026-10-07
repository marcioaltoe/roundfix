---
schema: roundfix/archive-record/v1
spec: 0200-a-skill-snapshot-that-matches-its-upstream
title: A skill snapshot that matches its upstream
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0200-a-skill-snapshot-that-matches-its-upstream
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0191
sources: []
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "309"
delivery_commit: bfebc0e961f9d40ca185a5da2da857929af2b809
---

# A skill snapshot that matches its upstream

The Baseline ships setup snapshots that name the upstream skills a profile installs, pinned to one commit of `marcioaltoe/skills`. The pin is `14fdf46`, from 2026-07-25; upstream is at `a4e18e4`. Upstream renamed `context7` to `context7-cli`, `feature-systems-pattern` to `app-renderer-systems` and `rust` to `rust-expert`, and moved `bubbletea` and `tui-design` out of `go-cli` into a new `go-tui` list. The catalog, the rendered guides and one universal capability still use the old names, so a new adopter is told to restore a skill no snapshot can supply. The asset sync cannot take the refresh: run with `--check` against the upstream checkout it exits `2` with `catalog.profile.skill.outside-setup` for `bubbletea`, `tui-design`, `context7`, `rust` and `feature-systems-pattern`. The `go` module dispatches three skills no Roundfix setup lists, because the catalog checks only required skills.
