# Published sources observed on 2026-10-03

- [CommonMark 0.31.2](https://spec.commonmark.org/0.31.2/), §§2.4/6.1, example 338. Primary page read through web. Backslashes are literal in code spans; an escaped-looking backtick can end the span. Supports detecting truncated authored command spans; R5 reproduces the command and the detector tests exercise both truncated cases.
- [POSIX.1-2024 shell language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html), set option -n. Direct web open returned Internal Error; a fresh search returned the primary-page indexed passage for -n: syntax is read without execution, with an interactive-shell caveat. Supports parsing with noninteractive sh -n before execution. R7d confirms exit 2; R5 marker assertions independently prove non-execution.
- [Git git-grep manual](https://git-scm.com/docs/git-grep), --untracked and --no-exclude-standard. Primary page read through web. --untracked adds untracked files to the tracked-file search; ignored files require the explicit additional option. Supports the non-ignored tracked/untracked detector boundary, independently exercised by R2 tests.

No contradiction observed. Shell and detector observations agree with each primary rule. Public curl attempt was denied CONNECT 403; no retry loop or escalation. No provider or private project data was sent.
