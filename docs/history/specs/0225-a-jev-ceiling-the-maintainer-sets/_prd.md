---
spec: 0225-a-jev-ceiling-the-maintainer-sets
status: archived
created: 2026-10-04
surfaces: [backend, cli, docs]
archived: "2026-10-04"
source_slug: 0225-a-jev-ceiling-the-maintainer-sets
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: row 09 needs the operator''s key-status record and installation; the operator read GET /api/v1/key on 2026-10-04 (limit 50, limit_reset monthly) and the maintainer''s decision is quoted in _authorization.md. Row 12 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 94761ac03fe623575143ccf90703580620e1999c
---


# A Jev ceiling the maintainer sets

Roundfix stops spending on Jev once the month's spend reaches a ceiling that
is compiled into the binary at US$5. The judge skips at that point, and a
routed Jev Router prompt is refused, together with any OpenRouter key whose
monthly credit limit is above US$5. On 2026-10-04 the maintainer raised the
limit of Roundfix's OpenRouter key and asked for the work to run without
that bound: "A chave do openrouter está com limite de $50 mensal e quero que
seja tudo implementado e testado sem limitações. Depois vemos o custo e como
reduzir o custo." Today only a release can change the number, and the key's
new limit makes every routed prompt fail as an unbounded key. This Spec makes
the ceiling a User Config value that the judge, the router gate and its
key-limit check all read, with US$5 as the default, so the maintainer sets
US$50 on this machine and an adopter who sets nothing sees no change.

## Prerequisites

None. Specs 0222 and 0223 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  key `jev.monthly_ceiling_usd` follows the existing dotted configuration
  names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, endpoint, credential
  or key handling changes; the judge and the router gate keep their
  recipients and keys, and no test or Verification reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0231 (this Spec) decides that the
  ceiling is a User Config value with the US$5 default. It supersedes in part
  ADR-0201, which sets "a spending ceiling of US$5 per calendar month", and
  ADR-0218, whose routed prompts run "under the US$5 monthly ceiling that
  ADR-0201 sets"; every other part of both stands, including ADR-0201's sum
  "across every repository and both recipients" and ADR-0218's refusal of an
  unbounded key. ADR-0218 also keeps the router selectable only from Project
  Config and calls User Config the configuration "which spans every
  repository on the machine", the reason the ceiling lives there. ADR-0002:
  "Roundfix uses YAML for User Config at", which names the file the operator
  edits. ADR-0027: "truly unknown keys keep failing strict validation", so an
  older binary refuses the new key and the operator sets it only after
  upgrading. ADR-0114: "A Fallback Selection may switch ACP Runtime
  automatically only while Agent work has not begun", so a refusal before
  Agent work stays a fallback and after it a failed Work Item, unchanged.
  ADR-0187 splits the
  Roundfix Skill by command and ADR-0189 ties an owned skill's version to its
  content, so the skill edit raises the version. ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts", applied to `roundfix spec judge`.
  The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155,
  ADR-0156 and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and
  ADR-0183 check this Spec's consistency by citation and receipt. ADR-0178
  authorizes each Task commit by its grant, ADR-0182 runs Settlement Checks
  before it, and ADR-0166 records undeclared paths; every Task declares its
  paths. ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park, which this Spec does not touch; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0192 cites ADR-0178 but decides derived-path conflict regeneration, and ADR-0200 and ADR-0209 cite ADR-0201 but decide that the judge never gates and that it only suggests which sources share a Spec, ADR-0208 cites ADR-0209 but decides how sources are grouped, and ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; this Spec changes none of them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, and the maintainer authorized skill
  edits ("considere autorizado a ajustar todas as skills se necessário") and
  this change on 2026-10-04 (quoted above). No other Governed Path changes;
  the repository's own Project Config is not edited, because Project Config
  cannot set this value. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0225-a-jev-ceiling-the-maintainer-sets/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- The maintainer raises the Jev ceiling to US$50 on this machine by editing
  User Config, with no release.
- With the ceiling at US$50, a routed prompt on a key whose monthly limit is
  US$50 runs, and the judge keeps asking until the month's spend reaches
  US$50.
- A machine whose User Config does not set the ceiling behaves exactly as
  before: US$5, the same messages, the same refusals.
- A repository cannot change how much of a machine's Jev budget is spent.

## User Stories

1. As the maintainer, I want to set the monthly Jev ceiling in my User
   Config, so that the work runs under the US$50 my OpenRouter key allows.
2. As the maintainer, I want the judge and the Jev Router gate to read the
   same ceiling, so that the month's shared spend stops at one number.
3. As the maintainer, I want the router gate to keep refusing a key whose
   monthly limit is above my configured ceiling, so that OpenRouter still
   stops a running prompt at my bound.
4. As an adopter who sets nothing, I want the US$5 default, so that an
   upgrade does not raise my spend.
5. As a user of a cloned repository, I want a Project Config value for the
   ceiling to be ignored with a warning, so that a repository cannot raise
   what my machine spends.

## Core Features

1. **A configured ceiling.** User Config accepts `jev.monthly_ceiling_usd`, a
   number of US dollars. A value that is not a finite number greater than
   zero is a configuration error naming the key.
2. **A default that changes nothing.** When User Config does not set the
   key, the ceiling is the built-in US$5 (ADR-0201).
3. **User Config only.** A Project Config value is removed before it is
   read and Roundfix prints the warning it already gives for a User
   Config-only setting; the effective ceiling is unchanged by it.
4. **The judge reads it.** `roundfix spec judge` skips and stops at the
   configured ceiling, and its text summary and JSON report show it.
5. **The router gate reads it.** A routed prompt is refused with
   `jev_ceiling_reached` when the month's Jev spend is at or above the
   configured ceiling, and with `jev_router_key_unbounded` when the key's
   monthly limit is above it; both messages name the configured value.
6. **The guides and the skill say so.** The configuration guide, the `spec`
   command guide and the Roundfix Skill's `spec` and `runtime` references
   name the key, its default, its scope and the order in which the operator
   sets it, and no longer state US$5 as the fixed ceiling. ADR-0201 and
   ADR-0218 carry a note that ADR-0231 supersedes their fixed value.

## User Experience

A machine without the key sees nothing new. With
`jev.monthly_ceiling_usd: 50` in User Config, the judge's summary reads
"month US$… of US$50.00" and the router gate accepts a key with a US$50
monthly limit. A Project Config that sets the key prints one warning on
standard error and keeps the User Config value. An invalid value fails the
command that loads configuration, as any invalid configuration value does.

## Non-Goals / Out of Scope

- Reducing Jev spend or choosing a different ceiling for adopters; the
  maintainer will look at cost afterwards.
- A per-repository or per-recipient ceiling, or letting Project Config lower
  the value (ADR-0231).
- Writing the operator's User Config: the Task writes no file outside the
  repository; the operator sets the value after merge.
- Changing the Judge Log, the spend sum, the recipients, the keys, the
  router's selection rules, or the measured judgment thresholds.
- Adding the key to the files `roundfix config init` writes.
- Relaxing strict validation of unknown keys for older binaries.

## Success Metrics

1. Success Metric: with `jev.monthly_ceiling_usd: 50` in a disposable User
   Config and a Judge Log month at US$50.25, `roundfix spec judge` skips with
   "monthly ceiling reached (US$50.2500 of US$50.00)"; with the key unset and
   the month at US$5.25 it skips with "of US$5.00", as before.
2. Success Metric: the router gate built with a configured ceiling of US$50
   accepts a key reporting a US$50 monthly limit, and the gate built without
   one refuses that key with `jev_router_key_unbounded` naming US$5.
3. Success Metric: a Project Config that sets the key leaves the effective
   ceiling at the User Config value or the default and prints the User
   Config-only warning; a value of zero, a negative value or an infinite
   value fails with the key named.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The maintainer's decision of 2026-10-04, quoted above, and a read-only
  query of the key's status on 2026-10-04 that reported `limit: 50` and
  `limit_reset: monthly` for `ROUNDFIX_OPENROUTER_API_KEY`.
- OpenRouter's published limits reference
  (<https://openrouter.ai/docs/api_reference/limits>): the key's `limit` is
  the "Credit limit for the key, or null if unlimited", and `usage_monthly`
  counts the "current UTC month", the quantities the router gate compares.
- A reproduction during authoring on 2026-10-04: the binary built from
  `main` at `6ea9e1e9`, given a disposable User Config holding
  `jev.monthly_ceiling_usd: 50`, refused `roundfix spec judge` with
  `field jev not found in type config.configOverlay` and exit `2`; with no
  User Config and a Judge Log month at US$5.25 it printed
  "monthly ceiling reached (US$5.2500 of US$5.00)".
- The Secondbrain holds no entry on the Jev ceiling beyond this repository's
  own mirrored Specs, searched on 2026-10-04.

## Decisions

- The ceiling is a User Config value with the US$5 default; Project Config
  cannot set it. See ADR-0231.
- The router gate keeps requiring a key limit no greater than the ceiling,
  now the configured one.
- This repository sets US$50 through the operator's User Config, not its
  Project Config.

## Open Questions

- When does the operator set US$50? Default: after this Spec merges and the
  Roundfix binary on the machine's `PATH`, the Delivery Queue owner and any
  running Run include it, the operator adds `jev:` with
  `monthly_ceiling_usd: 50` to `~/.roundfix/config.yml`; set earlier, every
  older binary on the machine refuses that User Config (ADR-0027).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
