# Concrete contracts

Concrete contracts make the parts of a Spec that other people and tools must
agree on readable and checkable. Use the forms below when a claim depends on a
decision record or when a TechSpec changes a command surface.

## Claim Receipts

A Claim Receipt puts the source and a verbatim quote in the same paragraph as
the attribution. The source is an `ADR-NNNN` identifier or a backticked
repository path. The quote is in straight double quotes and must contain at
least three words. The Spec Consistency Check finds the source, reads it, and
proves that the quote occurs there after collapsing whitespace. It proves
presence only; a reader still decides whether the quote supports the claim.

The check reads receipts in authored paragraphs and skips fenced blocks. A
receipt whose source cannot be resolved, whose quote is too short, or whose
quote is absent from the source is an error. A recognized attribution without a
receipt is a gap when the Spec is held at the contract horizon.

For example, this is illustrative syntax, so it is fenced and is not read by
the check:

```text
ADR-0183: "receipts prove presence and never support"
`docs/agents/domain.md`: "the source owns this vocabulary"
```

## Surface transcript blocks

A Surface Transcript records one command surface in the TechSpec. Number each
item and give it one fenced block with the `transcript` info string. The block
has exactly four lines and regions, in this order:

1. The first line starts with `$ ` and contains the command.
2. The next line is `stdout:`, followed by zero or more standard-output lines.
3. The next line is `stderr:`, followed by zero or more standard-error lines.
4. The final line is `exit: ` followed by an integer from 0 through 255.

The check validates the shape and does not run the command. The QA gate runs
the command through the built product and compares standard output, standard
error, and exit code. Two conventions allow controlled variation: a line that
is exactly `...` matches zero or more consecutive lines, and text in angle
brackets matches one or more characters within the same line. Everything else
matches exactly, including the exit code; standard output and standard error
are compared separately.

Copy each line from a run of the real command, not from memory. Keep the lines
a reader might skip: a final digest or summary line, a `Usage` block, leading
indentation and, for a `go run` command that exits non-zero, the `exit status
<n>` line that `go run` adds to standard error. The implementing Task's test
asserts every line, so a line left out here becomes a QA rerun later.

This is an illustrative block, so it is fenced and is not read by the check:

```transcript
$ roundfix spec check <slug>
stdout:
Spec <slug>
...
stderr:
exit: 0
```

## Interfaces and invariants

Write interfaces as code signatures so their shape is unambiguous. Number
invariants so a Task or review can refer to one rule without restating it.

```go
type Receipt struct {
	Source string
	Quote  string
}

func Receipts(artifact string, content []byte) []Receipt
func ProveReceipt(repoRoot string, receipt Receipt) (sourcePath string, proven bool, reason string, err error)

type Transcript struct {
	Command  string
	Stdout   []string
	Stderr   []string
	ExitCode int
}

func SurfaceTranscripts(content []byte) []Transcript
```

```text
1. A receipt is proven only when its normalized quote occurs in its source.
2. A transcript is compared through the built product at the QA gate.
3. Fenced examples are not declarations.
```

## Worked TechSpec fragment

The authoring skill carries the technical approach: `.agents/skills/write-techspec/SKILL.md`: "The PRD said _what_ and _why_; this document decides _how_, _where_, and _with which_". That receipt lets the check prove the attribution against a file shipped with this skill.

### Surface Transcripts

1. Surface Transcript: the guide example is discoverable through the consistency check.

   ```transcript
   $ roundfix spec check <slug>
   stdout:
   Spec <slug>
   No findings. Authored Verification commands were not executed.
   stderr:
   exit: 0
   ```
