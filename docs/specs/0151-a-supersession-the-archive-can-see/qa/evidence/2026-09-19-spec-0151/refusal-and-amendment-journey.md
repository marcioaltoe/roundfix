# Supersede refusals and one-file amendment

The public CLI ran from disposable Git root
`/private/tmp/roundfix-qa0151-refusals.bc2BB4`. Each negative command ran by
itself so its exit code was observed directly; a fresh `git status --porcelain`
after each of the first four refusals was empty.

| Probe | Exit | Named diagnostic | Post-command change |
| --- | ---: | --- | --- |
| unknown superseded Spec `9998-unknown` | 2 | `superseded Spec "9998-unknown" is unknown` | none |
| unknown superseding Spec `9999-unknown` | 2 | `neither active nor archived` | none |
| self-supersession | 2 | `cannot supersede itself` | none |
| unknown flag `--unknown` | 2 | `flag provided but not defined: -unknown` | none |

The accepted command exited `0`. `git status --porcelain` then named only
`docs/specs/0128-release-planning-with-bare-stable-tags/_supersession.md`, and
`git diff --exit-code` over every other path exited `0`. A fresh read confirmed
frontmatter fields `superseded_by`, `date`, and `reason`, plus the same reason
in the body.

A second supersession exited `2` with `already carries a supersession`. The
amendment SHA-256 was
`942d9ca693726820f0ce17b0132c0ef4945a86d1be3c4ff52d54b55eacf90486`
both before and after that refusal, and the fresh diff over all pre-existing
files remained empty.

