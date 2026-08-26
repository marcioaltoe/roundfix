---
spec: mechanical-carrier
status: closed
verdict: partial
rows_blocked_environment: 1
rows_blocked_finding: 1
rows_blocked_declared: 1
---

# QA report — mechanical carrier

## Results

| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| - | --- | --- | --- | --- |
| R01 | Green row | maintainer / backend | pass | [evidence](evidence/pass.txt) |
| R02 | Environment row | maintainer / backend | blocked (environment: fixture boundary) | [evidence](evidence/pass.txt) |
| R03 | Finding row | maintainer / backend | blocked (finding: QA-FIXTURE — waits on fixture) | [evidence](evidence/pass.txt) |
| R04 | Declared row | maintainer / backend | blocked (declared: fixture criterion) | [evidence](evidence/pass.txt) |
| R05 | Explicitly skipped row | maintainer / backend | skipped | [evidence](evidence/pass.txt) |

### Row detail — R01, and the comparison that justifies it

The table below is evidence written in prose. Its cells are names and
observations, not statuses, and it is not a second Results matrix.

| Case | Objection observed | Refused | Recovered |
| --- | --- | --- | --- |
| 82-line function vs 80 | function exceeds the 80-line limit | exit `1`, work staged | exit `0`, byte-identical |
| `sort()` vs `toSorted()` | use toSorted() instead of sort() | exit `1`, work staged | exit `0`, byte-identical |

