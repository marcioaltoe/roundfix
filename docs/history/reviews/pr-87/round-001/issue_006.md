---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/spec/spec_test.go
line: 1271
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeOA,comment:PRRC_kwDOS0qyts7c2V0s
review_hash: 872d40b314398b15c7301dc8da5dbfdc1c4f5343d33b83acef021fa2b4b9495c
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 006: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**The byte-identical assertion cannot fail.**

`Load` is called with `tempSpecsRoot` at Line 1252, but Lines 1265-1271 re-read `manifestPath` under `archivedRoot`. `Load` never opens that path, so the comparison always succeeds. The test therefore does not prove the PR's compatibility claim that legacy manifests stay byte-identical.

Compare the copied manifest inside `tempSpecsRoot`, which is the file `Load` actually reads.

<details>
<summary>💚 Proposed fix</summary>

```diff
-		after, err := os.ReadFile(manifestPath)
+		copiedManifest := filepath.Join(tempSpecsRoot, entry.Name(), "_tasks.md")
+		after, err := os.ReadFile(copiedManifest)
 		if err != nil {
-			t.Fatalf("re-read archived manifest %q: %v", entry.Name(), err)
+			t.Fatalf("re-read loaded manifest %q: %v", entry.Name(), err)
 		}
 		if string(after) != string(before) {
 			t.Fatalf("Load changed archived manifest %q", entry.Name())
 		}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		writeSpecDir(t, tempSpecsRoot, entry.Name(), files)

		graph, err := Load(tempSpecsRoot, entry.Name())
		if err != nil {
			t.Fatalf("Load archived Spec %q: %v", entry.Name(), err)
		}
		if graph.QATaskID != "" || graph.QADeclined || graph.QAReason != "" {
			t.Fatalf("archived Spec %q gained QA declaration state: (%q, %t, %q)", entry.Name(), graph.QATaskID, graph.QADeclined, graph.QAReason)
		}
		for _, task := range graph.Tasks {
			if task.Type == TaskTypeQA {
				t.Fatalf("archived Spec %q unexpectedly has QA Task %q", entry.Name(), task.ID)
			}
		}

		copiedManifest := filepath.Join(tempSpecsRoot, entry.Name(), "_tasks.md")
		after, err := os.ReadFile(copiedManifest)
		if err != nil {
			t.Fatalf("re-read loaded manifest %q: %v", entry.Name(), err)
		}
		if string(after) != string(before) {
			t.Fatalf("Load changed archived manifest %q", entry.Name())
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/spec_test.go` around lines 1250 - 1271, Update the
archived-manifest verification in the test around Load to re-read the manifest
copied under tempSpecsRoot, using the same entry-specific path passed to Load
rather than archivedRoot/manifestPath. Keep the existing byte-for-byte
comparison and failure message so the test verifies Load preserves the actual
manifest it reads.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d525fa6e28526f305d145d4f -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The archived-manifest characterization now re-reads the copied `_tasks.md` under `tempSpecsRoot`, which is the file passed to `Load`, instead of re-reading the untouched source artifact. `go test ./internal/spec -count=1` passed with 126 tests.
