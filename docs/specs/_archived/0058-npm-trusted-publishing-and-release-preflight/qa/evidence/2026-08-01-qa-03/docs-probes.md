# Documentation probes

Build: `171f6a378c9e640a8a10c9382e28b501b21ff5a0`

- `rtk ruby .../remediation_scope_probe.rb` — exit 0. It derived and found all
  six coordinates in the setup and shutdown procedures, confirmed the shared
  owner/repository/workflow binding, publish-time validation, empty-fallback
  exit condition, repository removals, per-package registry shutdown, and
  independent confirmation instructions. It also confirmed the two canonical
  glossary terms.
- `rtk ruby .../failure_vocabulary_probe.rb` — exit 0. Both the workflow and
  runbook expose exactly `identity`, `publish`, `registry`, `runtime`, and
  `undetermined`; no emitted prefix lacks a recovery row.
