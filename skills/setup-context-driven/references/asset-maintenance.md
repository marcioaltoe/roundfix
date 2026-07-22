# Maintain setup-context-driven assets

Treat `.agents/skills/setup-context-driven` as the only authorial source for this repo-owned skill. The catalog must prove semantic coverage and source identity before repository inspection or generated-document apply can begin.

## Catalog contracts

Edit the catalog as one connected contract:

- `assets/coverage.json` declares stable coverage identifiers. Universal categories include safety, selected Verification, Verification configuration integrity, skill dispatch, language, research authority, dependency discipline, Git and delivery, and security/configuration; enabled surfaces add their applicable categories.
- `assets/modules/*.json` declares versioned rules with coverage IDs and portable guidance. Every rule must have a selected supporting-guide carrier whose template renders `{{artifact.rules}}`.
- Each module's `requiredSkills` and `skillDispatch` keys must match exactly. Dispatch text states when an Agent must activate each required skill; no required or extra dispatch entry can hide in prose.
- Root blocks and supporting guides declare typed references. A setup-owned reference binds a template token to a selected managed identity. A repository-owned reference uses a safe repository-relative path and never makes that target setup-owned.
- `assets/profiles/*.json` declares `requiredRules`. Every required rule must belong to a selected module, be reachable through a selected guide, and collectively satisfy all universal and profile-applicable coverage categories.
- `assets/templates/` contains only generated baseline content. Keep root blocks compact, route details to setup-owned guides, and point to repository-owned contracts such as `DESIGN.md` without generating project-specific architecture.

`assets/contract-v1.json` lists the stable loader diagnostics for coverage, rules, dispatch, references, profiles, templates, and setup snapshots. Add mutation tests whenever a new invariant or diagnostic joins that contract.

## Snapshot-v2 validation

Every `assets/setups/*.json` file uses `setup-context-driven/setup-snapshot-v2`.

- An external skill declares `source.type: github`, a normalized `owner/repository`, a full immutable commit in `ref`, a safe source-relative `path`, and a lowercase complete-directory `treeDigest`.
- A Roundfix-owned skill declares `source.type: repo` and a `contentDigest`. Keep this source class separate from external provenance.
- The snapshot `digest` covers canonical serialization of the complete normalized skill records, not only their paths.
- Complete-tree hashing includes every regular file in bytewise POSIX-path order with path/content length framing, excludes `.git` and `node_modules`, and rejects symlinks and special files.

Validate or refresh external snapshot records only from an explicit canonical setups directory:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --source-dir <canonical-setups> --check --format json
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --source-dir <canonical-setups> --format json
```

`sync-setups` validates the declared immutable Git revision and complete source tree. It does not install skills, rewrite upstream content, select a convenient branch, or make the external checkout a runtime dependency.

## Source synchronization boundary

Author changes under `.agents/skills/setup-context-driven` first. Regenerate `skills/setup-context-driven` only through the repository workflow:

```bash
rtk make skills-sync
rtk make skills-sync-check
```

`make skills-sync` copies every repo-owned authorial skill and refreshes the recommended external-skill name list. Review the resulting diff and keep this Task limited to setup-context-driven files. Do not edit upstream-managed skill content; external skills remain represented only by immutable provenance and `skills-lock.json` compatibility metadata.

Before restoration writes ship, the isolated lock adapter must agree with the Spec 0036 lock compatibility fixture. Spec 0036 and the Doctor Command retain ownership of Repository Skill Set readiness and lock-hash compatibility. Setup documentation defines no new Doctor behavior.

Run the setup asset suite, `rtk make skills-sync-check`, and the repository Verification after catalog or workflow changes.
