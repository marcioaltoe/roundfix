# The queue is empty (2026-09-02)

`docs/specs/` carries no active Spec. Every item the 2026-08-26 triage listed is
delivered, and five releases went out. This records what shipped, what the work
found on the way, and what is left — which is no longer a Spec queue but a set of
decisions about mechanisms the practice has been working around by hand.

## Released

| Version | Spec | What it does |
| --- | --- | --- |
| v0.8.0 | 0117 | The Run Window: a repository-scoped cutoff on when an Implement Run may start |
| v0.9.0 | 0118 | Task Carry-Forward accepts an `Unresolved` Run; `implement` refuses rather than re-executing proved work |
| v0.10.0 | 0116 | A clean verdict states its own coverage; a QA Report names its Auditing Binary |
| v0.11.0 | 0105 | Roundfix owns the QA Task's Verification; a finding blocks the rows it names |
| v0.12.0 | 0097 | A Wave that would collide is reported before it dispatches; bootstrap is serialized across siblings |

The `kickoff` versus `implement-spec` question is settled and recorded at
`docs/workflow/2026-08-31-kickoff-is-settled-implement-spec-remains-the-loop.md`:
`implement-spec` remains the loop, and no `kickoff` implementation shipped.

## Read this first: the mechanisms the practice works around

Three captures in the Secondbrain inbox are the same shape, and they are the most
valuable thing this session found. Each is a rule the repository wrote down and
then quietly stopped being able to enforce, with the practice compensating by
hand every time.

**`reconcile` cannot prove integration under squash merge.** It classifies a Run
Branch by `merge-base --is-ancestor`. A squash merge creates a new commit with
the same tree and a different identity, so a Run Branch's commits never become
ancestors of it. Every terminal Run Branch therefore stays `unintegrated`
forever, `--apply` releases nothing, and each session clears the accumulation
with `git branch -D` — the operation ADR-0115 says the Daemon exists to spare the
maintainer. Thirteen branches and seven worktrees were removed by hand on
2026-09-02; the 2026-08-27 handoff described the same symptom for six earlier
ones without naming the cause. Three alternative proofs that survive squash are
written in the capture.

**`skills/` declares no ownership, so every grant enumerates mirrors.** ADR-0149
says an authorization names only the regeneration command and the audit resolves
its outputs from an `_ownership.yml` declaration — the shape ADR-0081 prefers
over enumerating consequences. `skills/` carries no such declaration, and
`OutputsFor` reads records only under `internal/baseline`, so the mechanism
cannot reach it. Every skill-editing Spec must list the generated copies
explicitly, which both ADRs reject. Note the trap: creating `skills/_ownership.yml`
would produce a file nothing reads. Closing this is a Spec, not a config tweak.

**`go vet` reports 31 `copylocks` the repository gate does not.**
`internal/agent.ACPXRunner` contains a `sync.Mutex` and has value receivers, so
every call copies the mutex. `go vet ./internal/agent` reports 31 findings;
`make verify` passes on the same tree. A stdlib analyzer that `go vet` runs by
default and the gate does not is a difference worth being deliberate about.

## What the delivery itself found

Recorded because each cost real Runs and none is obvious from the code.

**A capability can pass QA while existing only in a template.** Three times in
Spec 0116: the auditor recorded in one of three report producers, then Core
Feature 4 entirely inert — both writers passed empty evidence, so every report
the code wrote said `unknown`, and the `stale:` values in that Spec's own reports
came from the gate's Agent computing them by hand. A gate whose executor fills
the artifact can certify what the code does not do. When reviewing, ask which
production caller supplies the value.

**Ask when a skill runs, not only what it says.** Spec 0116's first cut applied
"name the probing form of the check" to every authoring skill including the QA
gate. The probe asks whether a command already passes before its work exists;
the gate runs after every Task is complete. Applied uniformly it refused the gate
with five of six Tasks vacuous, and would have refused every finished graph in
the repository.

**Documentation Tasks go after every behavior Task.** Spec 0118 ordered them
first, a corrective Task then changed the decision rule they described, and two
shipped surfaces contradicted the code until the gate caught it. Only the
glossary may come early, because it names the act rather than the rule.

**An ambiguous sentence in a TechSpec is expensive.** Spec 0097's TechSpec said
the collision rule "runs no command", meaning it executes no authored
Verification. It was read as "no subprocess", and the result was 695 lines of
hand-written Git — loose objects, pack indexes, delta decoding — with a test that
removed `git` from `PATH` to prove it. The repository invokes `git` for exactly
those questions in four packages. Deleted in review; the sentence is corrected in
the archived TechSpec.

**One defect, four times in one file, in both directions.** Spec 0097's
prior-Run source treats a failure to read as an absence of evidence. Review found
it in `git log`, in the `HEAD` probe added to fix the first, in `diff-tree`, and
in `repositoryFileExact` — the middle one introduced by the fix for the first.
Under-reporting is this rule's whole failure mode, so each instance silently
permitted the Wave it exists to refuse. The stable answer was to stop inferring
from the presence of an error and read the exit status: `rev-parse --verify
--quiet HEAD` gives 1 and silence for a repository with no commits, 128 with a
fatal for no repository, and everything else propagates. Related: a path filter
built for shell tokens was reused for paths Git reports verbatim, which trimmed
them and discarded any name containing a glob character — `fixture[1].go` is an
ordinary filename.

**Three families of gate that hides its exit status**, all found in authored
Verification commands this session: a pipe into `grep -q "^ok"` passes when any
package prints `ok`; the same idiom over several packages hides a red one; and a
`;` before the final command breaks the `&&` chain so the exit status comes from
the tests alone. Use `|| exit 1` explicitly and let `go test` report its own
status. Spec 0117's QA Task still carries the first shape.

## Process that worked

- **The review marker belongs in the pull request description at opening.**
  `.coderabbit.yaml` sets `auto_review.enabled: false` with
  `description_keyword: "coderabbit:review"`. Without it the check reports
  automatic review disabled, which reads as a pass. PR #168 merged that way and
  never had a review. With it, PRs #169 through #173 produced roughly thirty
  findings, including a shell injection through a Spec slug, a checker that read
  `granted 2026` as citing ADR-2026 and suppressed its own finding, and a
  production `ESRCH` defect invisible on macOS.
- **An absent review verdict is a block, not a pass.** Twice the checks were
  green with no formal verdict at the head; waiting produced
  `CHANGES_REQUESTED` both times.
- **Prove a red CI is not a regression before re-running.** Two red runs this
  session were pre-existing races; one was a real production defect. The proof
  is cheap: does the same tree pass the same job, and does the intervening
  commit touch that path.
- **`GOOS=linux go vet ./...` before trusting a green local gate.** A
  `//go:build linux` file does not compile on macOS, so no local `make verify`
  can reach it. That is how the `ESRCH` defect stayed hidden.

## Machine state left behind

- No Run Branches, no Run Worktrees, no retained artifacts. `~/.roundfix/worktrees`
  is empty.
- The Run Database is at schema 13 and the installed binary is 0.12.0, the
  released one, so the machine-wide database and the fleet agree. Do not run a
  `bin/roundfix` built from a branch carrying an unreleased migration.
- A release bump touches two files, and `dist/npm/roundfix/package.json` is the
  hinge between two different checks. `make verify` fails when
  `internal/app.Version` and that manifest disagree; the release workflow's tag
  validation fails when the pushed tag and that manifest disagree. So bumping
  only `internal/app.Version` is caught locally, which is how v0.12.0 was
  caught, and bumping only the manifest is caught at the tag.
- Nine pending Inbox Entries under `inbox/roundfix/`, five of them from this
  session. Two captures were triaged to `_triaged/` when their defects shipped.

## If you pick this up

The queue is not a list of Specs any more. The next decisions are the three
mechanisms above — which proof `reconcile` should use, whether `skills/` gets an
ownership declaration and the resolver that reads it, and whether the gate should
run the analyzer `go vet` already runs. None is urgent; each is a place where the
written rule and the practice have diverged without anyone deciding they should.
