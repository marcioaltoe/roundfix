---
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
prd: _prd.md
created: 2026-11-18
---

# A fresh home keeps the host toolchains, and a renamed Spec shows once — Technical Spec

## Executive Summary

Two independent fixes ship together. The `qa-gate` and `write-tasks` skills gain one section each that makes the QA environment explicit: suites and the repository Verification keep the session's home, and every command under a fresh `HOME` receives the host toolchain locations read before `HOME` changed. No Roundfix code isolates the gate's home, so the fix is text where the gate and the plan author read it, pinned by a contract test. The status view's row reader drops a phase entry's slug when another slug of the same Spec Number holds the Spec. The trade-off accepted is that the view relies on the `write-prd` rule that Spec Numbers are never reused: a repository that broke it sees one of two Specs being authored with the same number (ADR-0434).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the row reader reads the existing Spec Number of a slug with the `specNumber` helper of `internal/cli`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; the row reader reads the Spec Root, the Archive Root, the Delivery Queue it is given and the Phase Record. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0434 (new) and its amendment of ADR-0410; ADR-0080, ADR-0353, ADR-0407, ADR-0418 and ADR-0421 listed only to record why they do not apply. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: "Mesmas do BE (Recommended)" of 2026-11-18, quoted in the record; it bounds the four `SKILL.md` files below and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/write-tasks/SKILL.md`.

## System Architecture

Two components change and nothing new is wired.

1. The Spec queue row reader, `followUpQueueRows` in `internal/cli/follow_up_queue.go`, which the status view of the Follow-up Command calls. Today a candidate slug comes from a folder under the Spec Root that has no PRD, or from the latest counted phase entry of an Authoring Phase that `followUpAuthoringPhases` returns; it is dropped when it is a Delivery Queue item, has a PRD or has an Archive Record. The comparison is by slug, so a renamed Spec's earlier slug keeps a row (references, section 2). The fix adds one rule for a candidate whose folder is not under the Spec Root: it is dropped when its Spec Number is held by another slug. Nothing else in the view, its sources, its forecasts or its output format changes.
2. The owned skills `qa-gate` and `write-tasks`, canonical under `.agents/skills/` and mirrored under `skills/` by `make skills-sync`. Each gains one section; no existing sentence changes, and the `### QA settlement` section of `qa-gate` is untouched. Adopters receive both sections at their next `roundfix baseline update` or skill installation, through the owned-skill bundle.

The QA environment stays an instruction. The Daemon starts the QA Agent session with the host environment and sets no `HOME` (references, section 1), and the gate itself chooses the environment of each command it runs. A code-level isolation would need a new Daemon contract that the maintainer did not ask for.

## Implementation Design

### Interfaces

```go
// internal/cli/follow_up_queue.go

// followUpHeldSpecNumbers maps each Spec Number to the slugs that hold it:
// in.Active, in.Skipped folders, every in.Progress.queue item (merged ones
// included) and every Archive Record name under in.ArchiveRoot.
func followUpHeldSpecNumbers(in followUpQueueInput) map[string]map[string]bool

// followUpRenamedSlug reports whether entry names an earlier slug: another
// slug of its Spec Number is held, or another slug of that number has a
// later counted entry in latest (an equal time goes to the greater slug).
func followUpRenamedSlug(slug string, entry phaserecord.RecordedPhase,
	held map[string]map[string]bool, latest map[string]phaserecord.RecordedPhase) bool
```

In `followUpQueueRows`, a candidate whose folder is not a directory under the Spec Root keeps its row only when `latest` holds it and `followUpRenamedSlug` returns false. A candidate with a folder keeps today's behaviour. An unreadable or missing Archive Root holds nothing and fails nothing.

### Data Models

None. The Phase Record, the Delivery Queue and the Archive Records are read as they are; nothing is written.

### API Contracts

1. API Contract: the authoring rows of the Spec queue. A phase entry for a slug whose folder is not under the Spec Root does not make that slug a Spec being authored when another slug of the same Spec Number has a folder under the Spec Root, an Archive Record, or a Delivery Queue item of any stage, or when another slug of that number, held only by phase entries, has a later counted entry. A slug with a folder under the Spec Root keeps its row. The text and JSON shape of the status view does not change. The `follow-up` guide appends this sentence, word for word, to its paragraph that begins "A Spec is being authored when": "A Spec keeps its Spec Number when its slug is renamed, so an entry for a slug whose folder is not under the Spec Root does not count when another slug of the same Spec Number has a folder under the Spec Root, an Archive Record or a Delivery Queue item, a merged one included, or a counted entry recorded later: that entry names the Spec's earlier slug, and the Spec is one row." The glossary entry **Spec Queue Table** gains, before its ADR links, "A renamed Spec keeps its Spec Number, so an entry for a slug without a folder there does not count when another slug of the same number has a folder, an Archive Record, a Delivery Queue item or a later counted entry", and an ADR-0434 link. ADR-0410 gains, after its paragraph "Specs being authored", a paragraph that starts with the bold notice `**Amended by [ADR-0434](0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md) (2026-11-18).**` and goes on, word for word: "A phase entry for a slug whose folder is not under the Spec Root no longer counts when another slug of the same Spec Number has a folder there, an Archive Record, a Delivery Queue item or a later counted entry: that slug is a renamed Spec's earlier slug. The seven days, the repository match and the row a Spec folder makes still hold." and its `updated_at` moves to the time of the change. The Skill Coverage Map records the Coverage Review: the entry `command follow-up` gains the `review` "2026-11-18 — Spec 0407 task_01: a phase entry for a renamed Spec's earlier slug no longer makes a row of the Spec queue; no command, flag, field, output line or exit code changes, and the Roundfix follow-up reference relays the rows as the command gives them", and the entry `guide docs/user-guide/commands/follow-up.md` gains the `review` "2026-11-18 — Spec 0407 task_01: the guide states that a renamed Spec keeps its Spec Number and is one row; the Roundfix follow-up reference relays the rows as the command gives them and needs no change".
2. API Contract: the QA environment sections. The `qa-gate` skill gains, right before `### Evidence links resolve in the committed tree`, this section, word for word:

   > `### A fresh home keeps the host toolchains`
   >
   > A test-suite run, the repository Verification and every other command that does not run the built product keep the session's own `HOME` and `TMPDIR`, unless the `qa` Task requires otherwise. A row that runs the built product gives it a fresh temporary `HOME` when the Task asks for one, so that the product never reads the user's own configuration.
   >
   > Whenever a command runs under a fresh `HOME`, read the locations of the toolchains and caches the repository's stacks use before `HOME` changes, and pass them to that command unchanged: the fresh home hides the user's configuration, never the installed toolchains. For Go, pass `GOMODCACHE` and `GOCACHE` as `go env` prints them. For Rust, pass `RUSTUP_HOME` and `CARGO_HOME`, their values or `$HOME/.rustup` and `$HOME/.cargo` when they are unset. For another stack, pass the variables that locate its installed toolchain. The command uses those toolchains and never installs, updates or removes one. A requirement that forbids reading the real home directory still lets the gate pass these locations, unless it names them. Create the fresh home outside the checkout and remove it afterwards.
   >
   > A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: run it again with the toolchain locations and record that result.

   The `write-tasks` skill gains, right before `## Verification paths`, this section, word for word:

   > `## QA environment`
   >
   > Write a `qa` Task so that the gate keeps the toolchains the repository's stacks install under the home directory. Test-suite runs and the repository Verification keep the session's own `HOME` and `TMPDIR`; only a row that runs the built product needs a fresh temporary `HOME`, so that the product never reads the user's own configuration. Every command under a fresh `HOME` receives the locations of those toolchains and caches, read before `HOME` changes and passed unchanged, as the qa-gate skill's section "A fresh home keeps the host toolchains" lists them: for Go, `GOMODCACHE` and `GOCACHE`; for Rust, `RUSTUP_HOME` and `CARGO_HOME`. A requirement that forbids reading the real home directory names these locations as its exception. A suite runs under a fresh `HOME` only when the Task states why, and then it keeps the same locations.

   In both, the backticked heading is written as a plain Markdown heading.

### Surface Transcripts

None. No command line, field, message or exit code changes; API Contract 1 changes which rows the status view lists, and the QA gate reproduces it through the built binary (`_prd.md` → Success Metric 3).

## Coverage Map

- Goal 1 → the `qa-gate` section of API Contract 2.
- Goal 2 → the `write-tasks` section of API Contract 2.
- Goal 3 → `followUpHeldSpecNumbers`, `followUpRenamedSlug` in `followUpQueueRows` (API Contract 1).
- Core Feature 1 → the `qa-gate` section; Core Feature 2 → the `write-tasks` section.
- Core Feature 3 → `followUpHeldSpecNumbers`, `followUpRenamedSlug`.
- Core Feature 4 → the `follow-up` guide sentence, the **Spec Queue Table** entry, the ADR-0410 amendment, and the Coverage Review entries of the Skill Coverage Map for `command follow-up` and `guide docs/user-guide/commands/follow-up.md` (API Contract 1).
- Success Metric 1 → `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome`.
- Success Metric 2 → `TestARenamedSpecLeavesNoAuthoringRow`.
- Success Metric 3 → the QA reproduction through `bin/roundfix follow-up --status`.
- Success Metric 4 → the QA outside-evidence row, which applies the `qa-gate` section.

Deviations from the PRD: None.

## Absence Evidence

1. `followUpHeldSpecNumbers` and `followUpRenamedSlug` — command: `rg -n 'specNumber|Spec Number|renamed' internal/cli/follow_up_queue.go internal/cli/follow_up_status.go`; result: no match; the row reader compares slugs only.
2. The `qa-gate` section on a fresh home — command: `rg -n -i 'RUSTUP_HOME|CARGO_HOME|fresh home|fresh .HOME|TMPDIR' .agents/skills/qa-gate .agents/skills/write-tasks`; result: no match; neither skill names the QA environment.
3. The `write-tasks` QA environment section — command: `rg -n '^## QA environment' .agents/skills/write-tasks/SKILL.md`; result: no match.
4. `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome` and `TestARenamedSpecLeavesNoAuthoringRow` — command: `rg -n 'TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome|TestARenamedSpecLeavesNoAuthoringRow' skills internal`; result: no match.

## ADR Alternatives

- ADR-0434: the design takes its decision; the rejected options (a marked row, folder-only counting, rewriting the Phase Record, reading renames from Git) stay rejected.
- ADR-0410: amended, not superseded; its rejected options stay rejected, matching by slug alone, which it refused for repositories, is now also refused within one repository, and matching by repository path alone stays refused because a folder-less entry from another clone still counts.
- ADR-0080: not applicable — no QA verdict rule changes; an environment-blocked row still needs a proved cause, and the `qa-gate` section only states that a toolchain the gate itself hid is not one.
- ADR-0353: not applicable — the follow-up list and its Waiting Items are unchanged; ADR-0410 cites it.
- ADR-0407: not applicable — no phase entry is written or changed and the Phase Record stays append-only; ADR-0410 cites it.
- ADR-0418: not applicable — the model and effort a phase entry holds are neither read nor changed; it cites ADR-0410.
- ADR-0421: not applicable — the Repository Identity of a phase entry is neither read nor changed, and the repository match of ADR-0410 stays; it cites ADR-0410.

## Integration Points

None external. The row reader reads local files and the Delivery Queue store the status view already loads. The skills reach adopters through the owned-skill bundle and are read by the Agent that runs the gate or authors the Tasks.

## Testing Approach

- `TestARenamedSpecLeavesNoAuthoringRow` in `internal/cli/follow_up_queue_test.go`, on the existing seams `newFollowUpReadInput`, `followUpQueuePhase`, `followUpArchiveFixture`, `followUpWrite` and `followUpQueueDiscover`, with a phase entry for `0401-old-slug` and the subtests `renamed_to_a_planned_spec`, `renamed_to_an_idea_folder`, `renamed_to_a_queue_item`, `renamed_to_a_merged_queue_item`, `renamed_to_an_archive_record`, `renamed_while_authoring`, `another_number_stays` and `the_old_folder_stays`. Each asserts the exact list of row slugs. The prototype measured seven failing before the change and the control passing (references, section 3). `TestASpecIsBeingAuthored` and `TestTheQueueRowsComeInTableOrder` stay unchanged and keep passing.
- `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome` in a new `skills/qa_fresh_home_test.go`, reading both copies through `roundfixLayoutSources`, with one subtest per guarantee, named `<copy>/<skill>/<check>`, and a check that the `qa-gate` section precedes `### QA settlement`. `TestSettlementGuidanceIsOneTable` keeps the settlement section byte-identical.
- The sabotages of the QA gate mutate `followUpHeldSpecNumbers`, `followUpRenamedSlug` and one skill phrase, each observed by the subtest that reads that function or phrase.
- No test reads the clock, the network, the real home or the real Phase Record.

## Build Order

1. The status view drops a renamed Spec's earlier slug, with its test, the `follow-up` guide sentence, the **Spec Queue Table** entry, the ADR-0410 amendment and the re-recorded Behavior Surface Record (API Contract 1).
2. The `qa-gate` and `write-tasks` sections, their mirrors, the owned-skill versions and the contract test (API Contract 2). Independent of step 1: the two share no file.
3. The QA gate (depends on: 1, 2).

Files shared with Spec 0408, authored beside this one and delivered after it: `CONTEXT.md` (step 1 edits only the **Spec Queue Table** entry), `docs/references/behavior-surfaces.json` (re-recorded by step 1), `skills/testdata/owned-skill-versions.json` (step 2 appends `qa-gate` and `write-tasks` versions). No Verification of this Spec depends on another Spec.

## Risks & Considerations

- A repository that reused a Spec Number hides one of the two Specs being authored; ADR-0434 accepts it because `write-prd` forbids the reuse.
- The skill sections are instructions, so a gate that ignores them still hides a toolchain. The `write-tasks` section puts the same default in every new `qa` Task, which is where the defect started.
- A QA host without rustup cannot reproduce the Rust half of Success Metric 4; the PRD says what the row records then.
- Listing the Archive Root adds one directory read per status call, bounded by the number of Archive Records.
- The suite guard refuses a fresh `HOME` or `TMPDIR` inside the checkout, which is why the `qa-gate` section says to create the fresh home outside it.

## Glossary

- changes: **Spec Queue Table**

## Decisions

- One Spec for two defects that share no component: the maintainer answered "Uma Spec só (Recommended)" on 2026-11-18 after the measurement (`_prd.md` → Sources, Decisions).
- The QA environment rule is skill text, not code: no Roundfix code sets the gate's `HOME`.
- Go and Rust are named as examples of one stack-neutral rule; a repository without Rust passes nothing for it.
- A renamed Spec's earlier slug is dropped, identified by the Spec Number. See ADR-0434.
- A folder under the Spec Root keeps its row even when a phase entry names a later slug of its number, so the checkout stays the source of truth until the renamed plan reaches it. See ADR-0434.
