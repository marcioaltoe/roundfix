# Maintain setup-context-driven assets

Treat `.agents/skills/setup-context-driven` as the only authorial source for this repo-owned skill. The catalog must prove semantic coverage and source identity before repository inspection or generated-document apply can begin.

## Governed Source Baselines

Author the current governed corpus under `assets/source-baselines/`. Every retained Normative Clause, recommendation, and Operational Contract has one marker-bounded entry, one manifest record with enforcement and carrier, and one independently pinned index identity. Operational Contracts declare their required structure and retain the complete template, lifecycle, decision matrix, ordered procedure, or protocol; a summary is not an acceptable replacement.

Keep root carriers as compact delegation indexes and place durable behavior in `docs/agents/` carriers. Record prior-clause retention or intentional rejection individually in the strict `0.0.1` accounting document. Before refreshing corpus, manifest, or index digests, run the governed-corpus mutation gate and confirm the denylist rejects project names, brands, machine-specific paths, and copied generated markers.

## Catalog contracts

Edit the catalog as one connected contract:

- `assets/coverage.json` declares stable coverage identifiers. Universal categories include safety, selected Verification, Verification configuration integrity, skill dispatch, language, research authority, dependency discipline, Git and delivery, and security/configuration; enabled surfaces add their applicable categories.
- `assets/modules/*.json` declares versioned rules with coverage IDs and portable guidance. Every rule must have a selected supporting-guide carrier whose template renders `{{artifact.rules}}`.
- Each module's `requiredSkills` and `skillDispatch` keys must match exactly. Dispatch text states when an Agent must activate each required skill; no required or extra dispatch entry can hide in prose.
- Root blocks and supporting guides declare typed references. A setup-owned reference binds a template token to a selected managed identity. A repository-owned reference uses a safe repository-relative path and never makes that target setup-owned.
- `assets/profiles/*.json` declares `requiredRules`. Every required rule must belong to a selected module, be reachable through a selected guide, and collectively satisfy all universal and profile-applicable coverage categories.
- `assets/templates/` contains only generated baseline content. Keep root blocks compact, route details to setup-owned guides, and point to repository-owned contracts such as `DESIGN.md` without generating project-specific architecture.

`assets/contract-v1.json` lists the stable loader diagnostics for coverage, rules, dispatch, references, profiles, templates, and setup snapshots. Add mutation tests whenever a new invariant or diagnostic joins that contract.

## Transition-ledger maintenance

The canonical Upgrade Retention Contract assets live under `.agents/skills/setup-context-driven/assets/retention/`. Treat each transition ledger as reviewed migration evidence, not generated output. Its source baseline inventory must list every previously managed mandatory clause, and its mappings must account for each clause exactly once as `retained`, `moved`, `replaced`, or `rejected`, always with a reason. Accepted current-clause targets must be reachable in the destination profile and preserve the exact enforcement enum; a Repository-Owned Extension target records that setup transfers ownership without managing the target bytes.

Review the relevant ledger whenever a baseline ID, legacy fingerprint, clause identity, enforcement value, selected carrier, or destination profile graph changes. Do not infer an unknown source baseline or delete a prior clause to make a transition pass. Update the ledger and its sanitized source fixture together, then run the focused contracts:

```bash
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_legacy_rule_ledger.py'
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_contracts.py'
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_retention.py'
```

The Change Plan's ordered `retentionAccounting` and `planDigest` are the operator-visible proof that the reviewed ledger is active.

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

## Formatter-provenance refresh boundary

Formatter compatibility is profile-specific. The selected formatter contract lives in `assets/profiles/<profile>.json`; the TypeScript/Bun proof record lives at `tests/fixtures/formatter-compatibility/typescript-bun-monorepo/provenance.json`, and its `golden/` tree contains the complete generated managed-Markdown corpus. Profiles without a selected Markdown formatter must declare `kind: none` rather than inherit another profile's proof.

Refresh formatter provenance only when the selected formatter identity or version changes, formatter-sensitive rendering changes, or the generated corpus intentionally changes. Regenerate the entire disposable profile fixture through confirmed apply, run the exact pinned formatter, and update the golden bytes, fixture path list, provenance command and argv, and portable golden digest as one reviewable change. Never hand-normalize selected golden files or substitute a newer formatter version without changing the profile contract.

Ordinary `rtk make setup-context-check` is hermetic: it validates the pinned provenance and golden bytes without installing or executing Oxfmt. Final QA owns the real probe recorded in `realFormatterProbe`; run it in the disposable TypeScript/Bun fixture, then run `fixtureVerification`, a fresh audit, and a second apply and require an empty diff and Change Plan. A real-probe failure blocks provenance refresh instead of being converted into a golden update.

After any transition-ledger, formatter-provenance, workflow, or test change, author under `.agents/skills/setup-context-driven/`, run `rtk make skills-sync`, inspect both tree diffs, and run `rtk make skills-sync-check`. The distributed `skills/setup-context-driven/` tree, including tests and fixtures, must remain byte-identical to the canonical tree.
