---
status: accepted
created_at: 2026-10-20T22:00:00Z
updated_at: 2026-10-20T22:00:00Z
deprecated_at: null
superseded_by: null
---

# A Linux Lane skipped for its environment is shown where the delivery is read

On 2026-10-20 the `golang:1.26` image had been pruned, and the Delivery Queue
gate of three Specs skipped the Linux Lane with one console notice and merged
them; nobody noticed for three releases. A skip still passes the gate and never
parks the item, as ADR-0288 and ADR-0294 decided, but a skip of an enabled lane
for an environment reason is now recorded on the item and shown by `roundfix
deliver status`, the queue's Pull Request body and `roundfix doctor`.

## Decision

**Two kinds of skip.** A skip because the lane has nothing to test, `no go.mod
at the repository root` or `no changed Go package`, is an expected skip and
stays as quiet as before: the console notice and the lane log. Every other
skip reason of an enabled lane is a Lane Environment Skip: Docker or Go not on
PATH, `docker info` failing, the image absent, Docker refusing to start the
container, an object store the lane cannot mount whole, a `go.mod` without a
`go` version when no image is configured, and a fetched default whose
configuration cannot be read, because then nobody can prove the lane was off.

**The item carries it.** The gate returns the reason with its passing result,
and the Delivery Engine records it on the item as the warning
`linux-lane-skipped: <reason>`, always the last warning segment and replaced
by the next gate. `roundfix deliver status` prints it on the item's existing
`Warning:` line, merged items included, so the record outlives the queue run.

**The Pull Request says so.** When the item carries a Lane Environment Skip,
the body the queue writes gains a paragraph naming the reason and saying that
the candidate's Go packages were not tested on Linux by the queue gate.

**Doctor checks the lane it will need.** When the checkout's effective
`delivery.linux_lane` is `true`, `roundfix doctor` prints a `linux-lane:` line
that reports, without pulling or starting anything, whether Docker is on PATH
and answers `docker info`, and whether the image the lane would use is present
locally. Its findings are warnings: the lane would skip, not fail. With the
lane off the line is absent, so no other output changes.

## Considered Options

- Parking the item when an enabled lane is skipped was rejected by the
  maintainer's decision: a missing image is an operator's repair, and the
  repository gate already passed.
- Making every skip visible was rejected: an item that changes no Go package
  would carry a warning on every delivery of documentation, and the warning
  would stop meaning anything.
- A new item field or Run Database column was rejected: the schema version
  would rise and older item binaries sharing the database would refuse it; the
  existing item warning is durable and already printed by status.
- Reading the opt-in from the fetched default in `doctor` was rejected: doctor
  reports this checkout's configuration and reads no remote-tracking state;
  the user guide says that the queue reads the fetched default instead.

## Consequences

- An operator sees an unrun lane on the item, on the Pull Request and in
  doctor, and repairs it with `docker pull <image>` before the next delivery.
- The console and lane log notices keep their exact text, so existing readers
  of the log are unchanged.
