# Static gate evidence

Audited head: `90ce84c244f40967a8a74e061c681ef914c66fdb`.

- `make build`: exit 0.
- `./bin/roundfix --version`: exit 0; `roundfix 0.78.0 (90ce84c2-dirty, built 2026-10-14 11:47:55 -0300)`.
- `./bin/roundfix spec check 0268-a-turn-that-did-not-finish-is-not-accepted --strict`: exit 0; no findings; authored Verification commands not executed.
- `make verify-changed`: Daemon-provided pass, exit 0; successful diagnostics removed. This Agent did not rerun it.

The auditor identity and staleness fields in the seeded report remain unchanged.
