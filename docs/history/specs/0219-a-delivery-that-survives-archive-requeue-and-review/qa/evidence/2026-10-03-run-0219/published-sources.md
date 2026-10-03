# Published evidence obtained on 2026-10-03

Public fetch first attempted with curl-cffi for CommonMark: CONNECT tunnel failed, HTTP 403. No denial loop. Web tool fallback reached the primary published pages.

[CommonMark 0.31.2](https://spec.commonmark.org/0.31.2/), sections 2.4 and 6.1, example 338: backslashes remain literal in code spans and do not escape the closing delimiter. This supports the truncated-span detection and the observed shell command.

[POSIX shell language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html), set -n: direct page open returned 403, but the search tool returned the primary page's indexed set -n passage. It specifies parsing without execution and permits interactive shells to ignore it. The product uses noninteractive sh. Obtained through the primary page's search excerpt rather than full-page fetch; supports parse-before-execute.

[git-grep manual](https://git-scm.com/docs/git-grep), --untracked and --no-exclude-standard: default --untracked includes untracked worktree files; searching ignored files requires the extra no-exclude-standard option. Supports tracked/untracked pin detection and exclusion of ignored files. Tests independently observed those boundaries.

No published source contradicts the design. The two transcript findings are output-contract mismatches, not disagreement with these references.
