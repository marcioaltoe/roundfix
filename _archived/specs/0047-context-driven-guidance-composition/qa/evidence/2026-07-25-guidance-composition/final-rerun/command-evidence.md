# Command evidence — final 2026-07-25 rerun

## Build and source identity

- Roundfix source revision:
  `a2502b70e67e85ac05a9b384f520909ebed1c5db`
- QA binary:
  `/private/tmp/roundfix-qa-0047-final-a2502b7/roundfix`
- `roundfix version`:
  `roundfix 0.0.1 (a2502b70e67e85ac05a9b384f520909ebed1c5db)`
- Binary SHA-256:
  `0b04e9d5b17a1fa65e855ec9d160aa4fa4c578d48aee270ca7841ac9a94db343`
- Fluxus source:
  `1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`
- Oraculum source:
  `ad74f46197500de63dc0d9ff0d3e09f61a6a43ce`
- Disposable root:
  `/private/tmp/roundfix-qa-0047-final-a2502b7`

Both source checkouts were clean at the same revisions before and after the
run. All Baseline writes occurred in `--no-hardlinks` local clones.

## Static and criterion-level gates

The mandatory unmodified gate passed:

```text
rtk make verify
Go test: 2204 passed in 22 packages
Go test: 4 passed in 1 package
Roundfix skill check passed
Roundfix build passed
exit 0
```

The Task-focused commands passed with these totals:

```text
Task 01: 28 passed
Task 02: 7 passed
Task 03: 22 passed
Task 04: 11 passed
Task 05: 18 passed
Task 06: 10 passed
Task 07: 8 passed
Task 08: 28 passed
Task 09 macro journeys: 37 passed
Task 09 Git-lifecycle stress: 300 passed
```

The last command ran 100 repetitions of both
`TestRepositoryInspectionNoMutation` and
`TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout`
without an environment override.

Raw summaries: `focused-tests.txt`.

## LIVE-01 — Fluxus greenfield

Public automation planning exited `0`:

```text
roundfix baseline plan \
  --repo <fluxus-greenfield> \
  --profile standard-typescript-monorepo \
  --decision-file <fluxus-greenfield-decisions.json> \
  --format json
Plan Digest: sha256:c52337dce84d5153752f70bebb671f1907bb1d0b8f12cdd29a8709c73952d121
File changes: 13
Managed entries: 26
Upgrade Retention Contract entries: 0
```

Applying with an all-zero digest exited `3` and left the clone unchanged.
Applying the exact digest exited `0` with `state: verified`.

The first formatter attempt found no installed `oxfmt` in the clean clone.
`rtk bun install --frozen-lockfile` restored the lockfile-defined disposable
environment; the exact `rtk bun run fmt` rerun passed on 589 files. The
disposable Fluxus `rtk make verify` then passed.

Fresh public planning exited `0`:

```text
Plan Digest: sha256:485ab6ff87129660e2ea41f6c4c0b543cba54e7ef7384d4ee4051710311f7434
File changes: 0
Managed entries: 25
Upgrade Retention Contract entries: 0
```

Applying that exact empty Plan exited `0`, returned `state: verified`, and
verified 14 persisted postimages. The generated root contains the ordered,
strengthen-only Instruction Hierarchy. The generated documentation guide
contains all five ADR states and all four Findings states. No generic or
repository-specific carrier path exists.

Artifacts:
`fluxus-greenfield-plan.json`,
`fluxus-greenfield-wrong-digest.json`,
`fluxus-greenfield-apply.json`,
`fluxus-greenfield-format.txt`,
`fluxus-greenfield-bootstrap.txt`,
`fluxus-greenfield-format-rerun.txt`,
`fluxus-greenfield-verify.txt`,
`fluxus-greenfield-replan.json`, and
`fluxus-greenfield-reapply.json`.

## LIVE-02 — Fluxus update

Automation first exited `3` and requested the complete classification
Decision Document without writing a partial Plan. The public PTY flow then:

1. selected Preservation and reused
   `standard-typescript-monorepo`;
2. retained every existing repository decision;
3. reached ready Profile alignment;
4. reviewed 19 exact `AGENTS.md` source segments, each proposed as
   non-governed with reasoned rejection;
5. reviewed 13 file changes, 26 managed entries, and all 19 Upgrade Retention
   Contract entries;
6. confirmed Plan Digest
   `sha256:434be69ff208a5c34da005dd4d18f9094f436baeacba4c8c5847162563eb5e82`.

Apply exited `0` and verified 15 postimages. After lockfile bootstrap,
`rtk bun run fmt` and the disposable Fluxus `rtk make verify` passed.

The fresh public PTY run repeated the complete 19-segment review. Its Change
Plan displayed no file changes and Plan Digest
`sha256:b06920923598178a4dcdca75e4974a8a4e6d13f838249ea8086de548180e27b2`.
Exact-digest apply exited `0` with
`approved Baseline Plan is already applied and verified` and verified 14
postimages.

Artifacts:
`fluxus-update-initial.json`,
`fluxus-update-bootstrap.txt`,
`fluxus-update-format.txt`, and
`fluxus-update-verify.txt`.

## LIVE-03 — Oraculum Profile adaptation

The built-in `standard-typescript-monorepo` automation path exited `3` and
named ten profile-specific blockers. Supplying both `--profile` and
`--profile-file` exited `2` before mutation.

The public PTY flow then:

1. selected Greenfield and reused the built-in Profile;
2. retained the current repository decisions;
3. chose repository-owned Profile adaptation;
4. reviewed removal of `frontend`, `autonomous-work`, and ten
   profile-specific capabilities;
5. accepted the complete removal proposal and supplied
   `oraculum-backend`;
6. re-audited to ready;
7. reviewed the Profile postimage, 14 file changes, 23 managed entries, and
   the empty retention ledger;
8. confirmed Plan Digest
   `sha256:c5e892ce37890bdd76a4e46376a65643d6b9760f48d92df39b3efaebef710911`.

Apply exited `0` and verified 14 postimages.
`roundfix baseline profile validate oraculum-backend --format json` exited
`0`, and the disposable Oraculum `rtk make verify` passed.

A replan with the pre-adaptation Decision Document exited `2` because
`autonomous.enabled` was no longer selected. This out-of-order input changed
nothing. The Decision Document derived from the applied Profile's selected
decisions produced:

```text
Plan Digest: sha256:89aa58e8520967b5c3b981dc10858f075ff9a39c803cd716ac0eb415ff2e1f27
Profile: oraculum-backend
File changes: 0
Managed entries: 22
Upgrade Retention Contract entries: 0
```

Exact-digest reapply exited `0`, returned `state: verified`, and verified 13
persisted postimages.

Artifacts:
`oraculum-built-in.json`,
`oraculum-mutual-exclusion.json`,
`oraculum-profile-validate.json`,
`oraculum-bootstrap.txt`,
`oraculum-verify.txt`,
`oraculum-replan.json`,
`oraculum-replan-selected.json`, and
`oraculum-reapply.json`.

## Documentation, exclusions, and integrity

- Built `baseline --help`, `baseline plan --help`, and
  `baseline assets sync --help` all exited `0`. Help exposes mutually
  exclusive `--profile` / `--profile-file`, digest-confirmed apply, Profile
  validation, restoration, and asset synchronization with stable exit
  categories.
- `rtk make skills-sync-check` passed all four setup-skill contract tests.
- `rtk git diff --check` passed.
- The focused documentation suite passed 28 tests.
- `rtk rg -n 'Fluxus|Oraculum' internal/baseline/assets
  .agents/skills/setup-context-driven skills/setup-context-driven` exited `1`
  with no matches.
- The `main...HEAD` changed-path list contains no tooling configuration,
  dependency manifest, ignore file, plugin declaration, or version pin.
- The upstream ADR-format and skill immutability guards passed.
- No commit, push, tag, release, or pull-request command ran.
