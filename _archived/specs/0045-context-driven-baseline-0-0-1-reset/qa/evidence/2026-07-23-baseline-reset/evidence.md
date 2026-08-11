# Evidence — 2026-07-23 baseline reset

Build: `aa84d5a2919fcb454bc30f930b9a09ea1aede1b0`

## E-01 — Full repository Verification

Command:

```text
rtk make verify
```

Final result: exit `0`.

```text
Go test: 1727 passed in 20 packages
Ran 251 tests in 64.278s
OK
Ran 251 tests in 65.070s
OK
setup-context-driven assets: ok
Roundfix skill check passed: roundfix, write-idea, write-prd,
write-techspec, write-tasks, setup-context-driven, implement-task,
implement-spec, brainstorming, council, business-analyst, archive-spec,
qa-gate, evidence-gate
go build completed for ./cmd/roundfix
```

The first unchanged invocation reached `1726 passed, 1 failed, 1 skipped` in
the Go suite because macOS returned an `unlinkat` error while
`testing.TempDir` cleaned up
`TestRunArchiveUsesConfiguredExternalSpecRoot`. The exact test then passed once
with `-v` and 20 consecutive times with `-count=20`. The unchanged full
`rtk make verify` rerun above passed. No source or test code changed between
the failed cleanup, focused reproductions, and passing full gate.

## E-02 — Public Readoption inventory

Command:

```text
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B \
  .agents/skills/setup-context-driven/scripts/context_setup.py audit \
  --repo <disposable-readoption-repo> \
  --profile standard-typescript-monorepo --format json
```

Result: exit `3`, stdout parsed as
`setup-context-driven/audit-v1`, and stderr was empty. The incompatible Source
Baseline contained 3 carriers, 19 entries, and 763 bytes. The result reported
19 `readoption.disposition.required` findings plus
`readoption.baseline.incompatible`, with zero planned changes.

Two fresh runs returned the same Source Baseline identity and the same 19
ordered entry IDs, kinds, byte ranges, and digests. SHA-256 snapshots of all
110 disposable-repository files matched before and after both audits.

## E-03 — Required and recommended Repository Capabilities

The public apply command ran against independent disposable Readoption
repositories.

- With `.agents/skills/context7/SKILL.md` absent, apply exited `1` with
  `capability.required.missing`, named `capability.context7`, instructed the
  maintainer to add Context7 to the Repository Skill Set, planned no changes,
  and preserved every repository byte.
- With only Firecrawl evidence absent, apply exited `3` for plan confirmation
  and also emitted one non-blocking `capability.recommended.missing` warning
  for `capability.firecrawl`. The warning explained that Firecrawl is useful
  for structured web extraction. The preview preserved every repository byte.
- A fully evidenced Readoption preview reported 27 ordered capability outcomes:
  Context7 and Exa were satisfied; Firecrawl, `rtk`, and `rg` were satisfied;
  all required Standard TypeScript Monorepo stack capabilities were satisfied;
  and absent Inngest remained optional.

The clean-adoption capability failure is recorded separately in E-06.

## E-04 — Generated instructions and Skill Activation

The confirmed public setup generated `AGENTS.md` plus 15 Markdown guides.
`docs/agents/skill-dispatch.md` contained these exact bundles:

```text
bundle.production-code: coding-guidelines, clean-code, solid
bundle.hono-endpoint: hono-api-best-practices, hono, zod
bundle.hono-endpoint-persistence: hono-api-best-practices, hono, zod, drizzle-orm
bundle.qa: qa-gate, evidence-gate
bundle.delivery: conventional-commits, github-pr-workflow
```

Frontend, UI quality, testing, debugging, security, and delivery remained
distinct triggers. Generated guides also retained Spec routing, the findings
frontmatter/status lifecycle, Supervisor-to-ACP-Runtime delegation, and the
complete read-only Secondbrain protocol. Canonical and distributed
setup-context-driven trees had no recursive diff.

Oxfmt `0.59.0` checked all 16 generated Markdown files successfully. The first
`bunx` invocation added a fixture-local lockfile because the scratch repository
had no pre-existing Bun lock; a second `--no-install` check, the selected
Verification, two audits, and reapply preserved the resulting repository
snapshot byte-for-byte. The missing Standard TypeScript Monorepo architecture
contract is part of E-06.

## E-05 — Digest-bound Readoption plan and rejection probes

The complete public Readoption preview exited `3` with
`plan.confirmation.required` and plan digest
`ce949866a3ace5bd07aea06748649e7f52eb1620a9e2d719ba8cd5d2332f5537`.
It bound all 19 dispositions, 25 planned changes, 14 exact planned outputs,
27 capability outcomes, and five persisted Verification entries. The
disposable decision file used individual `repository-rules` and `rejected`
dispositions; the current 12-test `test_readoption_apply.py` suite covered all
four destinations and all structural entry kinds.

Negative public probes:

- `--confirm-plan not-a-digest` exited `2` with
  `plan.confirmation.invalid` and no write.
- After changing `package.json` by one byte, confirmation of the prior digest
  exited `3` with `plan.confirmation.stale`. The recomputed digest changed to
  `b25c2a33e064a0a866b2218bc8dbfc0e5e6b8ad79e9bd9415a84fb54c301e190`,
  and the repository snapshot did not change during the rejected apply.
- Missing decisions and all 19 unresolved dispositions were already visible
  in E-02 and produced no write.

## E-06 — Finding F-01: clean adoption bypasses the 0.0.1 profile contract

Two independent disposable repositories reproduced the same public behavior.
Each had `DESIGN.md`, a typed REST HTTP contract, a complete strict `0.0.1`
decision file, and the selected Verification fixture, but deliberately had no
`package.json`, `packages/frontend`, `packages/backend`, PostgreSQL evidence,
LogTape evidence, or other required application-stack evidence.

Command:

```text
rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B \
  .agents/skills/setup-context-driven/scripts/context_setup.py apply \
  --repo <empty-disposable-repo> \
  --profile standard-typescript-monorepo --format json \
  --decision-file <strict-0.0.1-decisions.json>
```

Observed:

1. Preview exited `3` only for `plan.confirmation.required`; it emitted no
   required-capability finding.
2. Confirmation of digest
   `290e1772fcb64fb5827a8c0a38449526d0c2417344280812f1189baeaf55a7f1`
   exited `0` with `managed.apply.completed`.
3. The written Setup Manifest had `schemaVersion: 1`, no top-level `version`,
   generator version `1`, and baseline `baseline.portable-v3`.
4. Generated frontend/backend guides contained the older generic portable
   rules. They did not contain the required frontend systems/public-boundary
   contract or backend domain/application/infrastructure, thin-handler,
   HTTP-independent-use-case, and Drizzle-owned-persistence contract.
5. The repository still had no `package.json` after the successful apply.

Expected: required stack and workspace gaps block without writes; a successful
apply emits only strict `0.0.1` Setup Manifest/Snapshot state and the exact
Standard TypeScript Monorepo architecture contract.

Impact: `Blocks-Completion`. The maintainer receives a successful result but
cannot adopt the promised 0.0.1 baseline through the normal clean-repository
entry point.

## E-07 — Confirmed Readoption apply and repository-owned preservation

Confirmation of the exact E-05 digest exited `0`; apply returned the same
digest and exact planned outputs. The written manifest used
`setup-context-driven/manifest/0.0.1`, top-level version `0.0.1`, and generator
version `0.0.1`. The exact unmarked Repository-Specific Normative Rules bytes
were created.

After adding `Maintainer-owned follow-up.` to that repository-owned file, a
fresh public reapply exited `0` with `managed.apply.empty`, planned zero
changes, preserved the edit, and left the full repository snapshot unchanged.

The final canonical audit could not close because the disposable repository
contained fixture skill bytes rather than all 86 immutable external skill
trees pinned by the shipped setup snapshot. It exited `1` with 86
`skills.required.drift` findings and no mutation. Exact immutable Repository
Skill Set restoration is the unblocking action. The same-build focused
Readoption suite completed its digest-rewritten fixture audit at exit `0`.

## E-08 — Formatter, selected Verification, audit, and reapply

For the clean Standard TypeScript Monorepo setup:

```text
rtk bunx --no-install oxfmt@0.59.0 --check AGENTS.md docs/agents
```

exited `0` and reported all 16 files formatted. The persisted
`python3 -B .formatter-fixture-verify.py` command exited `0`. Two public audits
exited `0` with no findings, and public reapply exited `0` with
`managed.apply.empty` and zero planned changes. The repository snapshot stayed
byte-identical across the second formatter/audit/reapply cycle.

This proves composition and idempotency of the behavior that was emitted; it
does not cure the wrong clean-adoption generation identified in E-06.

## E-09 — Distribution identity and protected versions

- `rtk ./bin/roundfix version` reported
  `roundfix 0.0.1 (aa84d5a-dirty, built 2026-07-23 02:20:46 -0300)`.
- `rtk go test ./internal/app ./internal/releaseplan ./skills` passed 149 tests.
- `test_version_contract.py` passed 4 tests.
- `CHANGELOG.md` starts and ends its release history at `0.0.1`.
- The checked-in launcher, five platform packages, profiles, setups, Source
  Baseline documents, owned skills, and Release Plan schema identify `0.0.1`;
  protected operational/upstream fixtures remained accepted.

E-06 proves that the emitted clean-adoption Setup Manifest is still outside
this otherwise coherent checked-in identity.

## E-10 — Live read-only release reset plan

The built CLI ran from a clean local clone at the QA build with its origin set
to the real Roundfix GitHub repository. Network inventory was read-only.

Text and JSON commands both exited `3`, returned
`approval_required`, targeted `v0.0.1` at
`aa84d5a2919fcb454bc30f930b9a09ea1aede1b0`, and produced the same digest:

```text
sha256:ea3563e9742d21bef24739da6c1a77ce2c9e04ee95483427b5d27cfefb3c7bdf
```

The JSON result used `roundfix.release-plan/0.0.1` and contained six immutable
tag records (local and remote identities for `v0.1.0`, `v0.2.0`, and
`v0.3.0`) plus three paginated GitHub Releases with database IDs, node IDs,
tags, and resolved target commits. `--reset-to` combined with `--from` exited
`2` with a useful diagnostic and no partial plan.

After both complete plans, the clone remained clean at the same commit with
all three local tags and the same remote. No tag, Release, ref, file, package,
or configuration mutation occurred.

## E-11 — Focused current-build criterion checks

All commands exited `0`:

```text
test_source_baselines.py                 7 tests
test_governed_corpus.py                  6 tests
test_standard_typescript_monorepo.py     9 tests
test_readoption_apply.py                12 tests
test_macro_profiles.py                   9 tests
test_documentation_contract.py           5 tests
make skills-sync-check                    pass
```

These checks support the isolated contracts they name. E-06 is an assembled
public-flow gap that the passing focused suites did not detect.

## E-12 — QA mutation boundary

The source worktree started clean. After QA, `git status --short` listed only
the new Spec-local `qa/` directory. `git diff --check` and
`make skills-sync-check` passed. All setup mutations occurred under
`/private/tmp/roundfix-qa0045.*`; the live release commands were read-only.
No commit, push, branch change, tag deletion, GitHub Release deletion, package
publication, or pull request occurred.
