# Characterization of the two defects, recorded by the author on 2026-11-18

Measured on the planning branch `docs/plan-0407`, based on `1a8e08b4`, on the
Mac mini that authored the plan. Every probe ran in a temporary directory
under the planning clone and was removed afterwards.

## 1. The fresh home of the QA gate hides the host Rust toolchain

### Where a fresh `HOME` is prescribed

No shipped text and no code gives the QA gate a fresh `HOME`:

- `rg -n -i 'RUSTUP_HOME|CARGO_HOME|rustup' --glob '!docs/history/**' --glob '!CHANGELOG.md'`
  found only the adopted Backlog Entry.
- `rg -n -i 'temporary (home|director)|fresh (temporary|home)|real home' .agents/skills docs/agents docs/user-guide CONTEXT.md internal/baseline/assets`
  found nothing in `qa-gate`, `write-tasks` or any guide about the QA
  environment.
- `rg -n '"HOME="|"HOME"' --type go internal cmd` outside tests found only
  the Delivery Owner's command environment and the ACP runner's reading of
  `HOME`; the Daemon gives the QA Agent session the host environment.

The fresh `HOME` comes from each Spec's own `qa` Task. Its authors wrote that
every row that runs the product sets `HOME` and `TMPDIR` to fresh temporary
directories, and that no row "reads or writes the real `~/.roundfix` or the
real home directory". The gate then also ran the whole `./internal/cli`
package under a fresh `HOME`. Neither the `qa-gate` skill nor the
`write-tasks` skill says what a fresh home must keep, so each author and each
gate decided it again.

### The failures it caused

| Spec | Evidence | What happened |
| --- | --- | --- |
| 0399 | `qa/qa-report-2026-11-16.md` at `a8b702c0`, row Q04 | "Full internal/cli run failed only TestGuidanceCompositionJourney/rust-cli because the installed rustfmt shim cannot select a toolchain in the mandated fresh HOME." The plan was amended (`f7a132d4`) to pass `RUSTUP_HOME` and `CARGO_HOME` read-only. |
| 0400 | Archive Record, `qa_override_reason` | "isolated HOME had no rustup toolchain for TestGuidanceCompositionJourney/rust-cli (same failure on the unchanged Delivery Base)"; archived under a QA Archive Override. |
| 0401 | Archive Record, `qa_override_reason` | The same words; archived under a QA Archive Override. |
| 0402 | `qa/qa-report-2026-11-17.md` at `bf17cb38`, row Q02 | "blocked (environment: isolated HOME has no configured rustfmt toolchain)"; the rerun passed with the host toolchain. |

### The mechanism, measured

`TestGuidanceCompositionJourney/rust-cli` (`internal/cli/baseline_release_gate_test.go`)
resolves `rustfmt` with `testfixture.RequireTool`, which skips only when the
tool is absent. On a host with rustup, `rustfmt` is present as a rustup proxy
in `$HOME/.cargo/bin`, and the proxy reads its default toolchain from
`RUSTUP_HOME`, which defaults to `$HOME/.rustup`.

```text
$ HOME="<fresh temporary directory>" rustfmt --version
error: rustup could not choose a version of rustfmt to run, because one wasn't specified explicitly, and no default is configured.
help: run 'rustup default stable' to download the latest stable release of Rust and set it as your default toolchain.
exit 1

$ env HOME="<fresh temporary directory>" RUSTUP_HOME="<host RUSTUP_HOME or $HOME/.rustup>" CARGO_HOME="<host CARGO_HOME or $HOME/.cargo>" rustfmt --version
rustfmt 1.10.0-stable (58cf8f9c35 2026-10-13)
exit 0
```

On this host `RUSTUP_HOME` and `CARGO_HOME` are unset and `~/.rustup` is a
symbolic link to another volume, so the locations must be read before `HOME`
changes; computing them from the fresh `HOME` finds nothing.

The same pair through the test, with `HOME` and `TMPDIR` fresh temporary
directories outside the checkout (the suite guard refuses them inside it) and
the host's `GOMODCACHE` and `GOCACHE`:

```text
$ go test -count=1 -timeout 15m -v -run '^TestGuidanceCompositionJourney$/^rust-cli$' ./internal/cli
        error: rustup could not choose a version of rustfmt to run, because one wasn't specified explicitly, and no default is configured.
--- FAIL: TestGuidanceCompositionJourney (2.53s)
    --- FAIL: TestGuidanceCompositionJourney/rust-cli (1.34s)
exit 1

$ # the same command with RUSTUP_HOME and CARGO_HOME read before HOME changed
--- PASS: TestGuidanceCompositionJourney (1.87s)
    --- PASS: TestGuidanceCompositionJourney/rust-cli (1.25s)
exit 0
```

### Published sources

- The rustup book, Environment variables
  (<https://rust-lang.github.io/rustup/environment-variables.html>):
  "`RUSTUP_HOME` (default: `~/.rustup` or `%USERPROFILE%/.rustup`). Sets the
  root `rustup` folder, used for storing installed toolchains and
  configuration options."
- The rustup book, Installation
  (<https://rust-lang.github.io/rustup/installation/>): "Note that you will
  need to ensure these environment variables are always set and that
  `CARGO_HOME/bin` is in the `$PATH` environment variable when using the
  toolchain."
- The Cargo Book, Environment Variables
  (<https://doc.rust-lang.org/cargo/reference/environment-variables.html>):
  "By default these are stored under `$HOME/.cargo` (`%USERPROFILE%\.cargo`
  on Windows), but this variable overrides the location of this directory."
- Go, Managing Go installations (<https://go.dev/doc/manage-install>): "Go
  stores intermediate build artifacts in the directory returned by `go env
  GOCACHE`" and "Go stores downloaded dependencies in the directory returned
  by `go env GOMODCACHE`".

Each stack keeps its installed toolchain or caches under the home directory
unless a variable moves them, so a fresh `HOME` hides them for every stack
that does this, not only Rust.

## 2. The status view keeps a renamed Spec's authoring row

### How the rows are chosen today

`followUpQueueRows` (`internal/cli/follow_up_queue.go`) makes a slug a Spec
being authored when it is not a Delivery Queue item, has no `_prd.md` under
the Spec Root and no Archive Record, and either has a folder under the Spec
Root or a phase entry of an Authoring Phase of the last seven days that
`followUpAuthoringPhases` counts for the repository. The comparison is by
slug only. When an amendment renames a Spec's folder, the phase entries of
the old slug stay in the Phase Record, and nothing holds that slug any more,
so it is shown as a separate Spec being authored for up to seven days.

### Reproduction

A temporary Git repository with `origin` set, one planned Spec under its new
slug and one phase entry under its old slug, written with the built binary
in a fresh `HOME`:

```text
$ roundfix phase start tasks --spec 0401-the-old-slug
Recorded tasks of 0401-the-old-slug for claude session <session>: model unknown, effort high (environment), run by the session
$ roundfix follow-up --status
...
Spec queue: 2
  #  Spec  Status                Start (-03)               End (-03)                     Note
  1  0401  not queued 0/0 Tasks  after it is queued        not enough history (0 of 11)  -
  2  0401  authoring: tasks      after its plan is merged  not enough history (0 of 11)  tasks since 11-18 20:05
...
```

The same Spec Number is two rows. An Archive Record or a Delivery Queue item
under the new slug gives the same second row, which is what the maintainer
saw for `0401-the-phase-boundary-follows-only-the-release-of-a-spec` after
`0401-a-supervisor-session-ends-only-when-the-maintainer-asks` was delivered.

## 3. Prototype of the design

On a scratch worktree of the same base, removed afterwards:

- The rule of `_techspec.md` → API Contract 1 in `followUpQueueRows`, and a
  test `TestARenamedSpecLeavesNoAuthoringRow` with eight subtests. Before the
  change seven subtests failed and the control `another_number_stays`
  passed; after it all eight passed, and `TestASpecIsBeingAuthored` and
  `TestTheQueueRowsComeInTableOrder` still passed unchanged.
- The two skill sections of `_techspec.md` → API Contract 2. `make
  skills-sync` and the owned-skill record command raised `qa-gate` and
  `write-tasks` by one patch each, and `TestGlossaryPhrasesAreDefinedTerms`
  passed. `TestTheSkillCoverageMapIsCurrent` stayed current after the skill
  edits and went stale only after the `follow-up` guide edit, so only the
  status Task re-records the Behavior Surface Record.
- `speccheck.GovernedPath`, probed with `go test -overlay` and nothing
  written to the tree: the four `SKILL.md` files of `qa-gate` and
  `write-tasks` are governed; the Go sources and tests, the `follow-up`
  guide, `CONTEXT.md`, the ADRs, `docs/references/behavior-surfaces.json`,
  `skills/testdata/owned-skill-versions.json` and the planning files are not.
  Measured again after rebasing on `d1c43f9e`, with
  `docs/references/skill-coverage.json` added to the probe: same result, and
  the Skill Coverage Map is not governed either.
- Each sabotage of `task_03.md` requirement 6, applied alone to the prototype
  and reverted afterwards, failed its named subtest: the Archive Root not read
  (`renamed_to_an_archive_record`), a merged queue item skipped
  (`renamed_to_a_merged_queue_item`), a later entry of the same Spec Number
  ignored, with the code still compiling (`renamed_while_authoring`), a slug of
  another Spec Number treated as the same Spec (`another_number_stays`), the
  rule applied to a candidate with a folder (`the_old_folder_stays`, and also
  `renamed_to_an_idea_folder`), the Rust sentence removed from the canonical
  `qa-gate` section (`canonical/qa-gate/rust_locations`) and the last sentence
  removed from the mirror
  (`embedded/qa-gate/a_hidden_toolchain_is_no_environment_block`). Without a
  sabotage all eight status subtests and the 34 skill guarantees passed, and
  the prototype's skill sections equal `_techspec.md` → API Contract 2 after
  whitespace normalization.
