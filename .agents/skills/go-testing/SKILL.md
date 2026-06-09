---
name: go-testing
description: Use when writing, modifying, or reviewing Go tests in this repository, including CLI tests, config validation tests, SQLite persistence tests, GitHub/API boundary tests, ACP process tests, daemon-loop tests, and concurrency tests.
---

# Go testing

Write tests that prove behavior, not implementation shape.

## Workflow

1. State the invariant in one sentence before adding a test.
2. Put the test at the lowest package layer that can catch the broken invariant.
3. Extend an existing test file when it already owns the behavior.
4. Add a negative case for validation, failure, or edge behavior when possible.
5. Run the focused test first when useful, then `rtk go test ./...`.

## Patterns

- Use table-driven tests when only inputs and expected outputs vary.
- Name subtests as behavior: `rejects missing pull request number`, not
  `case 1`.
- Capture loop variables before `t.Parallel()` in subtests.
- Use `t.Helper()` for helper functions that can fail the test.
- Use `t.TempDir()` for filesystem tests and `t.Setenv()` for environment
  tests.
- Keep assertions tied to observable behavior: exit code, stdout, stderr,
  file content, stored state, or returned error.
- Prefer real owned logic with fake external boundaries. Do not test that a
  mock received the value the same test wrote into it.

## Avoid

- No `time.Sleep` for synchronization. Use channels, contexts, test hooks, or
  observable conditions.
- No external network calls in unit tests.
- No tests that depend on local user config, current branch, wall clock, or
  global mutable state unless the dependency is injected or isolated.
- No test-only branches or methods in production code.
- No snapshot or golden file unless the artifact is an intentional product
  contract.

## Roundfix-specific checks

- CLI tests must check stdout, stderr, and exit code separately.
- Preflight tests must assert fail-first behavior before remote fetch or agent
  startup.
- Loop tests must prove termination on context cancellation, max run duration,
  budget exhaustion, or clean review state as separate cases.
- Duplicate review issue tests must prove the newest duplicate is selected and
  the older local issue is marked duplicated.

## Verification

Use these checks before finishing test work:

```bash
rtk go test ./...
```

If concurrency changed or goroutines are under test, also run:

```bash
rtk go test -race ./...
```
