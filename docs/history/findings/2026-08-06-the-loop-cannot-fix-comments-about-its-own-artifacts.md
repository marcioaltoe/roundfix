---
status: done
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
created_at: 2026-08-06
updated_at: 2026-08-26
---

# The loop cannot fix comments about its own artifacts

**Date:** 2026-08-06
**Found by:** a `watch --until-clean` Run on PR #136 settling `Unresolved`
with nine Review Issues failed for one identical reason.

## The deadlock

Roundfix writes one Markdown file per Review Issue under
`docs/specs/_reviews/pr-<n>/round-<nnn>/`, commits them, and pushes. CodeRabbit
then reviews that push — including those artifacts — and files new Review
Issues against them. Roundfix fetches those as the next round, assigns them to
a Batch, and the Batch refuses every one:

```text
issue 007 failed — reason: The finding is valid, but its target is an
unassigned prior-round Review Issue artifact that this Batch is forbidden to
edit.
```

The refusal is correct. A Batch may edit only the files assigned to it, and a
prior round's artifacts are never assigned. So the loop reaches a state it
cannot leave under its own rules: the comments are valid, the fix is forbidden,
and every subsequent round regenerates the situation because pushing the new
round's artifacts invites review of those too.

Final tally for the Run: 2 resolved, 9 failed, all nine on the same cause.
Cumulative for the Pull Request: 10 resolved, 1 invalid, 9 failed.

## Why the artifacts attract comments at all

They are not neutral text. `internal/rounds/rounds.go` composes each file as
`fmt.Sprintf("# Issue %03d: %s\n\n## Review Comment\n\n%s\n\n…", n, title, body)`,
where `title` is the review comment's first line and `body` is the comment
verbatim. Two consequences follow mechanically:

- CodeRabbit titles look like `_🗄️ Data Integrity & Integration_ | _🟠 Major_`.
  With the emoji removed the emphasis becomes `_ Data Integrity & Integration_`
  — a space inside the emphasis markers, which is MD037, and it renders as
  literal underscores rather than italics.
- The embedded body carries CodeRabbit's own fenced blocks, which have no
  language tag — MD040 — and sometimes no surrounding blank line — MD031.

So Roundfix emits Markdown that fails the lint it asks other repositories to
pass, and it emits it on every round, in every repository that runs the loop.

## Two independent fixes, and they are not substitutes

**Stop reviewing the artifacts.** A path filter excluding
`docs/**/_reviews/**/*.md` breaks the cycle at its source, and on this Pull
Request that exclusion was the only exit from the deadlock. This is the
maintainer's change and it is correct: a review record is evidence of a review,
not code, and reviewing it is self-reference.

**Stop emitting malformed Markdown.** The exclusion hides the defect from
review; it does not repair the committed bytes. Normalizing the derived title
(strip inline emphasis when composing a heading) and labelling embedded fences
would cost a few lines in the composer and a test, and would improve every
artifact the fleet has already committed going forward.

## The wider shape

Any tool that commits its own operational records into the repository it works
on will meet this: the record becomes input to the next cycle. The rule worth
carrying is that machine-written operational records are excluded from review
by default, and that the tool is still responsible for the quality of what it
writes — because a repository that never reviews those files is exactly the one
that will never notice they are wrong.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
