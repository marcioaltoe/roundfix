---
task: task_05
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: backend
complexity: medium
---

# Task 05: Phrase checks survive wrapping and name what they missed

## Overview

Spec 0173's first Run lost two correct Tasks. Their Verifications ran
`grep -q "<phrase of four to six words>" <guide>`, the Agent wrapped the phrase
across two lines of `docs/user-guide/commands.md` and `CONTEXT.md`, and the
command exited `1` without output, so the diagnostic artifact was empty and
nobody could see which phrase had missed. The Verification text is read from
the Task files by `roundfix spec check` at authoring and by the Daemon, which
runs it verbatim. The author reads the finding and fixes the Task while it is
still open, which is the stage ADR-0117 assigns to an authoring defect. A
detector that flags a structure which cannot wrap, such as an anchored
frontmatter line or a heading, would refuse honest Specs; one that misses the
line-bound form lets the same Run loss happen again.

## Requirements

1. MUST add `CodeVerifyWrapFragile = "SC-VERIFY-WRAP-FRAGILE"`
   (`SeverityError`) and `WrapFragileVerification(task spec.Task) []Finding` to
   `internal/speccheck/verification.go`. For each top-level command of a
   Verification line (split on `&&`, `||`, `;` and `|`, with a leading `!`
   allowed), a `grep` invocation, optionally after `rtk`, is reported when all
   of these hold:
   - its pattern, the value of the first `-e` or else the first operand that is
     not an option, contains whitespace;
   - the pattern does not start with `^`, `#` or `|` and does not end with `$`;
   - at least one file operand ends in `.md`.

   A grep that reads standard input is never reported.
2. MUST make each finding's summary name the phrase and the file, and its fix
   give the wrap-tolerant form with both filled in:
   `tr -s '[:space:]' ' ' < <file> | grep -qF -- "<phrase>" || { printf 'missing phrase in %s: %s\n' <file> "<phrase>" >&2; exit 1; }`
   for presence, and `! { tr -s '[:space:]' ' ' < <file> | grep -qF --
   "<phrase>"; }` for a `!`-prefixed grep.
3. MUST call it from `detectTaskCoverageAndContextReferences` in
   `internal/speccheck/citations.go` for pending non-QA Tasks, locating each
   finding at its Verification line as `InvertedExitVerification` does, and list
   the code as skipped when `_tasks.md` is absent. MUST register it in
   `stagedDetectors` in `internal/speccheck/coherence.go` as a Task-stage
   detector.
4. MUST add the code to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, with count `0` in
   `internal/docscontract/testdata/corpus-golden.json` and in the pin in
   `internal/spec/archive_layout_characterization_test.go`, leaving every other
   count unchanged. The golden's `update` text says the code joined at `0`.
5. MUST teach the wrap-tolerant form in the Verification guidance of
   `.agents/skills/write-tasks/references/task-template.md`, naming
   `SC-VERIFY-WRAP-FRAGILE` and giving the self-reporting presence form
   verbatim. It MUST regenerate `skills/write-tasks/references/task-template.md`
   with `make skills-sync` and pin that guidance with a new test in
   `skills/baseline_skill_contract_test.go`.
6. MUST add the glossary entry **Wrap-Fragile Phrase Check** to `CONTEXT.md`
   beside **Inverted Verification Exit**, naming `SC-VERIFY-WRAP-FRAGILE`.
7. MUST prove by executing the suggested presence form through `sh -c` that it
   exits `0` against a Markdown file whose phrase is wrapped across two lines,
   and exits `1` against a file without the phrase with stderr naming the
   phrase and the file.
8. MUST put the new `speccheck` tests in
   `internal/speccheck/verification_wrap_test.go`, each negative case a test of
   its own, and MUST NOT rename or remove an existing top-level test.

## Subtasks

- [ ] Add the detector, its registration, its call and its skip.
- [ ] Update the corpus code list, the golden and its pin.
- [ ] Teach the form in the Task template, pin it, add the glossary entry and
      run `make skills-sync`.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] `grep -q "records no QA row" docs/user-guide/commands.md` is reported
      once, with the phrase, the file and the filled wrap-tolerant form; a
      `!`-prefixed line-bound grep is reported with the absence form.
- [ ] The wrap-tolerant form, an anchored pattern, a heading pattern, a single
      word, a grep over standard input and a `.json` operand are not reported.
- [ ] A completed Task and the `qa` Task are not checked, and `--strict` keeps
      the finding an error.
- [ ] The suggested form passes on a wrapped phrase and names phrase and file
      when it misses.
- [ ] The active corpus golden records `0` for the code.

## Context

- interface: `internal/speccheck/verification.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- creates: `internal/speccheck/verification_wrap_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestWrapFragileVerificationReportsALineBoundPhraseAgainstMarkdown|TestWrapFragileVerificationNamesThePhraseTheFileAndTheFix|TestWrapFragileVerificationReportsANegatedLineBoundPhrase|TestWrapFragileVerificationAcceptsTheWrapTolerantForm|TestWrapFragileVerificationAcceptsAnAnchoredPattern|TestWrapFragileVerificationAcceptsAHeadingPattern|TestWrapFragileVerificationAcceptsASingleWord|TestWrapFragileVerificationAcceptsStandardInput|TestWrapFragileVerificationAcceptsANonMarkdownFile|TestWrapFragileCheckSkipsACompletedTask|TestWrapFragileCheckSkipsTheQATask|TestStrictKeepsAWrapFragilePhraseAnError|TestMissingTaskGraphListsTheWrapFragileSkip|TestWrapTolerantFormPassesOnAWrappedPhrase|TestWrapTolerantFormNamesTheMissingPhraseAndFile|TestTaskTemplateStatesTheWrapTolerantPhraseForm|TestTaskTemplateStatesTheStatusPreservingVerificationForm|TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/speccheck ./internal/spec ./internal/docscontract ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestWrapFragileVerificationReportsALineBoundPhraseAgainstMarkdown TestWrapFragileVerificationNamesThePhraseTheFileAndTheFix TestWrapFragileVerificationReportsANegatedLineBoundPhrase TestWrapFragileVerificationAcceptsTheWrapTolerantForm TestWrapFragileVerificationAcceptsAnAnchoredPattern TestWrapFragileVerificationAcceptsAHeadingPattern TestWrapFragileVerificationAcceptsASingleWord TestWrapFragileVerificationAcceptsStandardInput TestWrapFragileVerificationAcceptsANonMarkdownFile TestWrapFragileCheckSkipsACompletedTask TestWrapFragileCheckSkipsTheQATask TestStrictKeepsAWrapFragilePhraseAnError TestMissingTaskGraphListsTheWrapFragileSkip TestWrapTolerantFormPassesOnAWrappedPhrase TestWrapTolerantFormNamesTheMissingPhraseAndFile TestTaskTemplateStatesTheWrapTolerantPhraseForm TestTaskTemplateStatesTheStatusPreservingVerificationForm TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095 TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q '"SC-VERIFY-WRAP-FRAGILE": 0' internal/docscontract/testdata/corpus-golden.json && grep -q "SC-VERIFY-WRAP-FRAGILE" .agents/skills/write-tasks/references/task-template.md && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "Wrap-Fragile Phrase Check" && diff -r .agents/skills/write-tasks skills/write-tasks >/dev/null` — expected: exit 0; before this Task none of the new named tests exists, the golden carries no `SC-VERIFY-WRAP-FRAGILE` entry and `CONTEXT.md` has no Wrap-Fragile Phrase Check entry, so the command fails.

## References

- `_prd.md` → Goal 5; Core Feature 5; Success Metric 5; Decisions.
- `_techspec.md` → The wrap-fragile phrase check; API Contract 7; Testing
  Approach 5; ADR-0093; ADR-0117; ADR-0133; ADR-0148.

## Result

Implemented the Task-stage `SC-VERIFY-WRAP-FRAGILE` detector. It classifies
top-level `grep` commands (including `rtk` and leading `!` forms), derives the
first `-e` or operand pattern, refuses only unanchored multi-word patterns with
a Markdown file operand, and supplies filled presence or absence repairs. The
Spec checker now runs it only for pending non-QA Tasks, locates findings at the
declaring Verification line, and names its skip when `_tasks.md` is absent.

Updated the active-corpus code registry and pinned zero count, the canonical
write-tasks template and generated mirror, the skill contract, and the project
glossary. `make skills-sync` regenerated the distributed template from the
authorized canonical source.

Acceptance evidence:

- The positive, diagnostic, negated, first-`-e`, and command-chain tests pass;
  findings carry `SeverityError`, name `records no QA row` and
  `docs/user-guide/commands.md`, and render the filled presence or absence
  form.
- Separate negative tests pass for the wrap-tolerant pipeline, start and end
  anchors, headings, table rows, a single word, piped/explicit/redirected
  standard input, and a non-Markdown operand.
- Spec-check fixtures pass for completed and QA Task exclusion, strict error
  severity, Task-stage registration, missing-graph skip reporting, and the
  Verification source line.
- The suggested presence repair was executed through `sh -c`: it exits zero
  for `records no\nQA row` and exits one for a missing phrase while stderr
  includes both the phrase and the temporary Markdown path.
- The docscontract corpus checks pass with `SC-VERIFY-WRAP-FRAGILE: 0`; the
  archive-layout characterization and both Task-template contract tests pass.

Focused checks:

- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 ./internal/speccheck`
  — passed.
- Targeted `go test` runs for the write-tasks template contracts,
  archive-layout characterization, and docscontract corpus checks — passed.
- `make skills-sync` followed by a canonical/distributed template diff —
  passed with no difference.
- `GOCACHE=/tmp/roundfix-task05-gocache make verify-incremental` — the sandboxed
  run reached two existing force-stop integration tests but could not read the
  host process table; the permitted rerun passed vet, all Go packages, skill
  checks, and build.

The Daemon-owned `## Verification` command was not run in this Agent turn.
