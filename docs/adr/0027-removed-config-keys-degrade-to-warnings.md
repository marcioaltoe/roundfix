---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Removed config keys degrade to warnings

Removing or renaming a config key never hard-fails an existing configuration: recognized deprecated keys are ignored with one stderr warning naming the replacement, while truly unknown keys keep failing strict validation. The 0009 cycle proved the alternative hostile — dropping `resolve.concurrent` with a Preflight rejection broke every Run on the machine whose config still carried it, which is precisely the population a migration must protect. A deprecated-keys table in the config package is the single place future removals register.
