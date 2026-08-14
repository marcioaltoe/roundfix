# Documentation and scope evidence

Build: `75161e9c3a5f7554cd1e0b9290bce6c61820b5c7`.

## Discoverability

`docs/agents/specific-repository.md` independently states:

- authorization records and prerequisite fixes are separate commits before
  the authorized Task, in either relative order;
- consequent fixes are separate commits after their cause;
- folding either kind into the Task commit fails the tooling-authority gate;
- `make baseline-digests` owns deterministic fallout after an expressly
  authorized Roundfix-owned Skill or Baseline module edit;
- hand-edited pins remain unauthorized;
- `GOCACHE` defaults to the ignored repository-local cache only when the
  environment supplies none.

ADR-0081 agrees with the standing regeneration policy. `rtk make help` exposes
`baseline-digests` with the expected description. The canonical and embedded
Roundfix Skill copies compare byte-identical and document the shipped entry
precondition, executable refusal, and regeneration command.

## Managed-block and scope audit

The base-to-build guide diff changes only:

- one repository-authored docs-layout table row before the first setup-owned
  block; and
- repository-authored `specific-repository.md`, which contains no
  `setup-context-driven:begin` marker.

No setup-owned managed block changed. No production Run Worktree maintenance
implementation path appears in the Spec delta.

The Makefile still declares the same repository Verification graph:
`fmt-check`, full `test`, Skill sync/check, and `build`. `rtk make -n verify`
printed all of those commands. The Spec adds the cache default and separate
regeneration target without removing or bypassing a verified check.
