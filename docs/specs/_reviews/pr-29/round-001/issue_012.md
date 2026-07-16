---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/tests/test_sync_setups.py
line: 1
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYf3,comment:PRRC_kwDOS0qyts7WgS5k
review_hash: 1a5d56a0c3b1f76895fd1fea82774c0af97d848d994c2e2cb89644bbf99f0922
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:21Z"
---

# Issue 012: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the exact skill-file digest.**

A wrong constant digest currently passes. Compare `contentDigest` with the expected SHA-256 value.

<details>
<summary>Proposed fix</summary>

```diff
 import json
+from hashlib import sha256
...
             added = next(skill for skill in snapshot["skills"] if skill["name"] == "canonical-example")
             self.assertEqual(added["path"], canonical_path)
-            self.assertTrue(added["contentDigest"])
+            self.assertEqual(
+                added["contentDigest"],
+                sha256(skill_file.read_bytes()).hexdigest(),
+            )
```
</details>






Also applies to: 131-133

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/tests/test_sync_setups.py` at line 1, Update the
test assertions in test_sync_setups.py to compare contentDigest against the
expected SHA-256 digest value, rather than only validating its presence or
format. Apply this to both the primary assertion and the additional setup
assertion around the referenced lines, preserving the existing test structure.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:9bb038587b35b8808522a971 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The sync setup test only asserted digest presence. It now compares `contentDigest` to the exact SHA-256 digest of the generated skill file content.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
