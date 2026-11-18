---
type: fix
status: open
created: 2026-11-18
spec: null
reason: null
waiting_on: "the single release of the 0399-0406 queue (maintainer: Spec after the release)"
since: 2026-11-18
---

# The QA gate's isolated home lacks the host Rust toolchain, and the status view keeps a renamed Spec's phase

## Symptom

1. The QA gates of Specs 0399, 0400, 0401 and 0402 each ended `blocked (environment: …)` because the fresh `HOME` they prescribe leaves rustup without a default toolchain, so `TestGuidanceCompositionJourney/rust-cli` fails in the whole `./internal/cli` package (the same failure on the unchanged Delivery Base). Each Spec needed a QA Archive Override or a plan amendment that passes `RUSTUP_HOME` and `CARGO_HOME` read-only. The `qa-gate` skill and the write-tasks QA guidance should make that the default for every Spec.
2. `roundfix follow-up --status` lists `0401-the-phase-boundary-follows-only-the-release-of-a-spec` as "authoring: tasks" after the plan was amended and the slug renamed (to `0401-a-supervisor-session-ends-only-when-the-maintainer-asks`, delivered): the Phase Record keeps the old slug's last phase.

## Decisions

- 2026-11-18: "Spec após a release"; authorization "Mesmas do BE (Recommended)" (to be recorded in the Spec's `_authorization.md`).
