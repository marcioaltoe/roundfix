### fetch

```bash
roundfix fetch --source coderabbit --pr <number> [--spec <slug>]
```

Validates local state, creates a Fetch Run, fetches unresolved CodeRabbit
review threads, writes markdown Round artifacts, and stops at the `Fetched`
outcome. It never starts an Agent, commits, pushes, or resolves Review Source
threads. With automatic Round selection it reuses an existing matching Round
when the same HEAD already has the same Review Issue fingerprints and never
overwrites existing Round artifacts.

