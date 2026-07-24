# ADR-0072: Baseline Go cutover preserves Python contracts

Status: Accepted

Before the Python setup runtime is removed, every input schema, supported
state, transition, behavior, action, refusal, planned byte sequence, digest,
and rollback outcome that it currently maintains must have characterization
evidence and a Go destination. New product behavior is recorded explicitly as
a designed delta rather than being mistaken for parity. The former
`restore-skills` and `sync-setups` responsibilities remain Go operations under
`roundfix baseline skills restore` and `roundfix baseline assets sync`;
neither expands the interactive adoption state machine.
