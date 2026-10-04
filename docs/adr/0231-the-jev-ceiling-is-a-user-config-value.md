---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Jev ceiling is a User Config value

ADR-0201 fixed the judge's monthly spending ceiling at US$5 and ADR-0218 made
every routed Jev Router prompt share it, refusing an OpenRouter key whose
monthly credit limit is above it. On 2026-10-04 the maintainer raised the
key's limit and asked for the work to run unconstrained: "A chave do
openrouter está com limite de $50 mensal e quero que seja tudo implementado e
testado sem limitações. Depois vemos o custo e como reduzir o custo." A value
compiled into the binary can only be changed by a release, so the ceiling
becomes configuration.

The ceiling is `jev.monthly_ceiling_usd` in User Config. It must be a finite
number greater than zero. When it is unset, the built-in US$5 of ADR-0201
applies, so a machine that sets nothing sees no change. Project Config cannot
set it: a Project Config value is ignored with the warning Roundfix already
gives for a User Config-only setting. The judge, the Jev Router gate and its
key-limit check all read the same effective value, and the key's monthly
credit limit must still be at most that value.

User Config is the scope because the ceiling is compared with a sum that is
already per machine: the Judge Log in Roundfix Home adds up every
repository's and both recipients' spend, so a value chosen by one repository
would decide whether another repository may spend. A repository also must not
be able to raise how much of the maintainer's money a clone spends. Letting
Project Config lower the value was considered and left out, because no one
asked for a per-repository budget and a lower value against a shared sum
would only make that repository stop first.

## Consequences

- ADR-0201's US$5 becomes the default rather than the rule, and ADR-0218's
  "US$5 monthly ceiling" reads as the configured ceiling; both are superseded
  in that part only.
- An older Roundfix binary refuses a User Config that carries the new key, so
  the operator sets it only after every Roundfix binary that reads that User
  Config includes this decision.
- Raising the ceiling raises real spend: the OpenRouter key's own monthly
  limit stays the hard stop for the router, and TypeSafe-direct calls are
  bounded only by the Judge Log sum.
