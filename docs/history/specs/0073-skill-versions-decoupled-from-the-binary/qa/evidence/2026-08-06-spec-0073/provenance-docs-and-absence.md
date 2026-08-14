# Provenance, documentation, and absence checks

The three Go Skills added in commit `23953af7` are
`golang-dependency-management`, `golang-safety`, and
`golang-structs-interfaces`. Their files exist and `go.json` plus the generated
dispatch route to them. Fresh searches found none of the three names in
`skills-lock.json` or any setup snapshot. `roundfix skills list` also omitted
all three from its bundled and recommended lists. Their upstream metadata
versions in installed files are not Roundfix provenance lock entries.

Other checks:

- `rtk make skills-version-check`: exit 0.
- `rtk make skills-sync-check`: exit 0; four Skill contract tests passed.
- `rtk make skills-version-check OWNED_SKILLS='roundfix coding-guidelines'`:
  expected exit 2 naming the injected unversioned member.
- `rtk diff -rq .agents/skills/roundfix skills/roundfix`: exit 0.
- A diff-added-line search across implementation-range Go tests found no
  copied `0.0.2` compatibility version or 64-hex digest assertion.
- A search found no `contentDigest` in setup snapshots, catalog diagnostics,
  or plan-characterization corpora.
- The Roundfix Skill names the minimum comparison, `satisfies`, `below
  minimum`, `unversioned or unresolvable`, all four blocking facts, the
  third-party boundary, and the no-regeneration rule.
- The implementation range changes no `requiredSkills` line, archived Spec,
  or Spec 0067 derived-ownership record.
