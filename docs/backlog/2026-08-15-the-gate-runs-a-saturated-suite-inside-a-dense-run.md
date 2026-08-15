---
status: pending
created_at: 2026-08-15
updated_at: 2026-08-15
---

# The gate runs a saturated suite inside a dense Run

`test-budget` runs `go test -parallel $(GO_TEST_PARALLEL)`, which is sixteen on a
machine with ten cores. That oversubscription is survivable on a developer's
machine — measured green end to end on 2026-08-15 at build `4482c019`. It is not
survivable inside an Agent Session, where the authoritative gate runs on top of
the Daemon, the adapter, and the Session's own process tree.

Spec 0103 measured the consequence three times. Its gate reported three ACPX
fixture journeys dying under repository-wide load — a fixture killed before its
milestone, an unproven adapter lineage, and a session close returning exit code
`-1` — while the same trio passed 5/5 in isolation and the whole `internal/agent`
package passed at `-parallel 16`. The failures move between runs and between
tests, which is the signature of saturation rather than of any one test.

Spec 0103 removed every cause it could reach: fixtures are compiled and linked
rather than written and executed, waits observe their children, and the stress
case is sized to available parallelism. What remains is the gate's own execution
context, which no fixture change reaches.

Worth settling when this is picked up: whether `GO_TEST_PARALLEL` should follow
the machine's core count rather than a fixed sixteen; whether the authoritative
gate should run outside the Agent Session that requested it, so a Run is not
measuring itself; and whether a Run should reserve capacity for the gate the way
it already separates Task Capacity from Verification Capacity under ADR-0056.

Changing `GO_TEST_PARALLEL` is a protected tooling mutation and needs express
maintainer authorization with the Makefile named. The Makefile is currently
bounded to Spec 0104 in the 2026-08-12 umbrella grant, which does not cover this.
