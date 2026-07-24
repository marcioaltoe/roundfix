# Evidence — 2026-07-23 baseline reset rerun

Build: `2498252f7921788466e7d48a87acdce930fa518b`

## E-01 — Full repository Verification

Command:

```text
rtk make verify
```

Result: exit `0`.

```text
Go test: 1727 passed in 20 packages
Ran 252 tests in 66.015s
OK
Ran 252 tests in 66.753s
OK
setup-context-driven assets: ok
Roundfix skill check passed: roundfix, write-idea, write-prd,
write-techspec, write-tasks, setup-context-driven, implement-task,
implement-spec, brainstorming, council, business-analyst, archive-spec,
qa-gate, evidence-gate
go build completed for ./cmd/roundfix
```

Both canonical and distributed setup suites, all Go packages, asset loading,
the shipped skill contract check, and the CLI build passed on the same source
state. The QA report itself makes the worktree dirty, so build provenance
correctly carried `2498252-dirty`.

## E-02 — Public clean adoption and Baseline Readoption journeys

Command:

```text
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B \
  docs/specs/0045-context-driven-baseline-0-0-1-reset/qa/evidence/\
2026-07-23-baseline-reset-rerun-2498252/public_setup_probe.py
```

Result: exit `0`. The probe invokes the shipped `context_setup.py` command in
subprocesses against disposable repositories.

Clean Standard TypeScript Monorepo adoption:

- A complete 0.0.1 Decision File without application capability evidence
  exited `1` with 20 `capability.required.missing` findings and preserved
  every repository byte.
- After adding exact package, workspace, PostgreSQL, LogTape, Repository Skill
  Set, design, and HTTP Contract evidence, preview exited `3` with
  `plan.confirmation.required`, 27 capability outcomes, 13 planned outputs,
  and Setup Snapshot schema
  `setup-context-driven/profile-snapshot/0.0.1`. Preview wrote nothing.
- `--confirm-plan not-a-digest` exited `2` with
  `plan.confirmation.invalid` and wrote nothing.
- Exact confirmation exited `0`, preserved preview/apply digest and output
  parity, and wrote Setup Manifest schema
  `setup-context-driven/manifest/0.0.1`, manifest and generator version
  `0.0.1`, and baseline
  `baseline.standard-typescript-monorepo-0.0.1`.
- Generated frontend guidance contained the domain-system and public-boundary
  contract. Generated backend guidance contained domain/application/
  infrastructure layers, thin HTTP handlers, HTTP-independent use cases, and
  infrastructure-owned persistence.
- Two fresh public audits exited `0`. Reapply exited `0`, returned zero
  planned changes, and preserved every repository byte.
- Changing `package.json` after preview made exact confirmation exit `3` with
  `plan.confirmation.stale`; the digest changed from
  `3d6bd80445aec02d41f8ca0232db4192b3e29456a1392ff5de8383fa9baafa9b`
  to
  `2978efbe2bdf4924e315aeace9fca7e2bf651cab2ba30dc745ed970877f31b7d`,
  and the command preserved the changed preimage.

Baseline Readoption:

- Public preview exited `3` with 19 independently identified Source Baseline
  entries and 19 dispositions.
- Exact confirmation exited `0`; the typed document and proposed
  Repository-Specific Normative Rules bytes matched the preview.
- A maintainer edit to Repository-Specific Normative Rules survived reapply
  byte-for-byte with digest
  `e5e0d7e5e66fd8f92fa8988c58f152f00aaec621b1ba531770155469c625978c`.
- Fresh public audit exited `0`; reapply exited `0` with zero planned changes.

## E-03 — Focused criterion and independent-tree checks

Every command exited `0`.

```text
canonical test_readoption_apply.py       13 tests
distributed test_readoption_apply.py     13 tests
canonical test_macro_profiles.py          9 tests
distributed test_macro_profiles.py        9 tests
test_source_baselines.py                  7 tests
test_governed_corpus.py                    6 tests
test_standard_typescript_monorepo.py       9 tests
test_documentation_contract.py             5 tests
test_version_contract.py                   4 tests
test_skill_dispatch.py                    12 tests
test_capabilities.py                       9 tests
test_profile_alignment.py                 10 tests
test_restore_skills.py                    14 tests
test_assets.py                            12 tests
test_decision_plan_contracts.py           10 tests
test_decision_rendering.py                 9 tests
rtk go test ./internal/app ./internal/releaseplan ./internal/cli ./skills
                                            755 tests in 4 packages
rtk make skills-sync-check                pass
rtk diff -qr .agents/skills/setup-context-driven \
  skills/setup-context-driven             no differences
```

The public-boundary Readoption suite covered clean adoption, all four
structural entry kinds, REST and Post-only contracts, required and recommended
capabilities, incomplete decisions, stale confirmation, changed preimages,
rollback, postwrite tampering, audit, and byte-empty reapply. The macro suites
ran apply/formatter/Verification/audit/reapply for all three maintained
profiles from both shipped trees.

## E-04 — Live read-only release reset plan

A clean local clone at build
`2498252f7921788466e7d48a87acdce930fa518b` used the real Roundfix GitHub
origin. The managed sandbox initially denied DNS resolution. The exact command
then ran with read-only network access.

Commands:

```text
rtk roundfix release plan --reset-to v0.0.1 --format json
rtk roundfix release plan --reset-to v0.0.1
```

Both exited `3` with state `approval_required`, target `v0.0.1`, target commit
`2498252f7921788466e7d48a87acdce930fa518b`, and plan digest:

```text
sha256:0f26c0e3e4b2991d38f06c6512f07026c19af589a7330a021ffc241b66f8a724
```

JSON used schema `roundfix.release-plan/0.0.1`. Text and JSON each inventoried
six immutable tag records: local and remote identities for `v0.1.0`, `v0.2.0`,
and `v0.3.0`. They also inventoried three paginated GitHub Releases, including
database IDs, node IDs, tags, target commitishes, and resolved commits.

`--reset-to v0.0.1 --from v0.3.0` exited `2` with a useful incompatible-flag
diagnostic and no partial plan.

Independent post-plan reads passed:

```text
rtk git status --short                   clean
rtk git rev-parse HEAD                   2498252f7921788466e7d48a87acdce930fa518b
rtk git tag --list                       v0.1.0, v0.2.0, v0.3.0
rtk git ls-remote --tags origin          all three remote tags present
rtk gh release list --repo marcioaltoe/roundfix
                                           v0.1.0, v0.2.0, v0.3.0 present
```

No local file, ref, tag, remote tag, GitHub Release, package, release, or
configuration mutation occurred.

## E-05 — Public help, version, documentation, and ownership boundary

The shipped setup help exited `0` and documented read-only audit, atomic apply,
skill restoration, setup synchronization, stdout/stderr separation, and exit
categories `0`, `1`, `2`, and `3`.

The built CLI reported:

```text
roundfix 0.0.1 (2498252-dirty, built 2026-07-23 03:15:26 -0300)
```

`roundfix release plan --help` exited `0` and documented reset mode as
read-only, complete local/remote tag and paginated GitHub Release inventory,
exit `3`, and a separately authorized deletion boundary.

The focused documentation contract passed five tests. Direct inspection found
the complete Standard TypeScript Monorepo setup and Baseline Readoption
commands in `docs/user-guide/context-driven-development.md`, the owned and
protected 0.0.1 boundary in the same guide and shipped setup skill, and the
post-QA read-only plan boundary in `docs/user-guide/release-runbook.md` and the
shipped Roundfix skill.

## E-06 — QA report integrity and mutation boundary

`git diff --check` exited `0`. `git rev-parse HEAD` remained
`2498252f7921788466e7d48a87acdce930fa518b` on the existing Roundfix Run
Branch. Worktree changes were limited to:

```text
docs/specs/0045-context-driven-baseline-0-0-1-reset/qa/qa-report-2026-07-23.md
docs/specs/0045-context-driven-baseline-0-0-1-reset/qa/qa-report-2026-07-23-aa84d5a-fail.md
docs/specs/0045-context-driven-baseline-0-0-1-reset/qa/evidence/2026-07-23-baseline-reset-rerun-2498252/
```

The current report contains 26 terminal `pass` rows and zero `pending`,
`fail`, `blocked`, or `skipped` rows. Both referenced evidence files resolve.
No commit, push, branch creation, pull request, tag deletion, GitHub Release
deletion, or package publication occurred.
