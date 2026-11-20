---
status: accepted
knowledge: false
created_at: 2026-11-20T14:50:00Z
updated_at: 2026-11-20T14:50:00Z
deprecated_at: null
superseded_by: null
---

# The owner of a Skill Activation requires every skill of its bundle

The module that owns a Skill Activation, or a module that owner depends on,
requires every skill of the activation's bundle, and the catalog refuses to
load otherwise. A profile's Repository Skill Set stays the union of the
`requiredSkills` of its modules, read by one function.

## Context

The skill dispatch guide renders the bundle of every Skill Activation whose
owner the profile selects. The Repository Skill Set that restore, the update
and Doctor read is the union of the selected modules' `requiredSkills`. Nothing
tied the two. ADR-0191 held a module to the skills its own dispatch entries
name, and ADR-0222 held an activation's owner to the Setup Snapshot; no rule
held an activation's bundle to its owner.

On Roundfix 0.152.0 the `typescript` module owned the production-code
activation, whose bundle names `solid`, while only the `frontend` module
required `solid`. Every built-in profile that selects `typescript` also selects
`frontend`, so no built-in profile showed it. On 2026-11-20 a repository
profile with `typescript` and without `frontend` rendered `solid` as mandatory
and never required it: the update did not report it and the restore command
refused it by name. A sweep of the 11 activations and 17 modules found this one
gap.

## Decision

For every Skill Activation, each skill of its bundle is in the
`requiredSkills` of its owner or of a module the owner reaches through
`dependsOn`. Catalog validation reports each exception as
`catalog.skill.activation.skill.unrequired`, naming the activation and the
skill, and a catalog with that diagnostic does not load. A module the owner
does not depend on never satisfies the rule, whatever profiles select beside
it.

A dependency satisfies the rule because every profile, built-in or
repository-owned, must select the dependencies of each module it selects. The
rule therefore holds for every composition the catalog accepts, not only for
the built-in profiles.

The `typescript` module requires `solid`. The `frontend` module keeps it,
because it dispatches `solid` under its own trigger.

`Catalog.RequiredSkills` returns the skills a list of modules requires, and
restore, the update preview, `skills check` and Doctor read the set only
through it.

## Considered options

- **Compute the set from the dispatch.** The readers would add every skill the
  rendered dispatch names to the union. It is correct by construction and
  needs no module edit. It was rejected because `requiredSkills` would stop
  being the whole declaration: every other reader of that field, the user
  guide, the Roundfix skill's runtime reference and ADR-0241 would describe a
  set the code no longer uses, and the catalog would keep a module that names
  a skill it never declares.
- **Move the activation to another module.** ADR-0222 already rejected `core`,
  whose profiles include setups that do not install the bundle; `frontend`
  would drop the production-code bundle from every backend-only repository.
- **Remove `solid` from the bundle.** The TypeScript Setup Snapshot carries
  the same bundle with `solid`, and the maintainer asked for the skill to be
  installed, not dropped.
- **A test instead of a catalog diagnostic.** The sibling rule for dispatch
  entries is a test. A diagnostic was chosen because the catalog already
  validates an activation's owner and bundle at load, and a rule at load also
  refuses a catalog that an asset sync or a local build produces.

## Consequences

A change that adds a skill to an activation's bundle must add it to the
owner's `requiredSkills` in the same change, or the catalog does not load.
A skill may now sit in the `requiredSkills` of two modules. Adopters whose
profile selects `typescript` without `frontend` see `solid` as a Missing
Required Skill on their first update after the release that carries this rule.
