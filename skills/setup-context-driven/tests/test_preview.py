import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_audit import snapshot_files  # noqa: E402


ENTRY_DECISIONS = {
    "language.generated",
    "verification.gate",
    "spec.scaffold",
    "domain.layout",
    "triage.external",
    "autonomous.enabled",
    "secondbrain.enabled",
}


class PreviewCliTests(unittest.TestCase):
    def test_audit_repeated_decisions_exposes_authorizable_concrete_digest(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = [
                "spec.scaffold=true",
                "domain.layout=single-context",
                "triage.external=false",
                "autonomous.enabled=false",
                "verification.gate=make verify",
                "language.generated=English",
                "secondbrain.enabled=false",
            ]
            args = [
                "audit", "--repo", str(repo), "--format", "json", "--profile", "rust-cli"
            ]
            for decision in decisions:
                args.extend(["--decision", decision])

            first = run_context_setup(*args)
            second = run_context_setup(*args)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            first_payload = json.loads(first.stdout)
            self.assertRegex(first_payload["planDigest"], r"^[0-9a-f]{64}$")
            self.assertEqual(first_payload, json.loads(second.stdout))
            self.assertTrue(first_payload["plannedChanges"])
            for change in first_payload["plannedChanges"]:
                self.assertEqual(
                    {
                        "path", "managedId", "state", "reason",
                        "beforeDigest", "afterDigest",
                    }.difference(change),
                    set(),
                )

    def test_typescript_first_apply_returns_entry_decisions_and_preview(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            result = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "typescript-bun-monorepo",
            )

            self.assertEqual(result.returncode, 3, result.stderr)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(result.stdout)
            self.assertEqual(payload["schemaVersion"], "setup-context-driven/audit-v1")
            self.assertEqual(payload["summary"]["decisions"], len(ENTRY_DECISIONS))
            self.assertEqual(
                {finding["managedId"] for finding in payload["findings"]},
                ENTRY_DECISIONS,
            )
            self.assertEqual(
                set(payload["findings"][0]),
                {"code", "severity", "path", "managedId", "message", "action"},
            )
            self.assertEqual(payload["selection"]["profile"], "typescript-bun-monorepo")
            self.assertEqual(payload["selection"]["setup"], "typescript-bun")
            self.assertIn(
                {"id": "core", "state": "active"},
                payload["selection"]["modules"],
            )

            planned = payload["plannedChanges"]
            self.assertGreater(len(planned), 0)
            self.assertTrue(any(change["state"] == "definite" for change in planned))
            definite = next(change for change in planned if change["state"] == "definite")
            self.assertNotIn("condition", definite)
            conditional = self._planned_change(
                planned,
                "guide.autonomous-work",
                "conditional",
            )
            self.assertEqual(
                conditional["condition"],
                {"decisionId": "autonomous.enabled", "equals": True},
            )

    def test_audit_and_blocked_apply_share_preview_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            audit = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "typescript-bun-monorepo",
            )
            after_audit = snapshot_files(repo)
            apply = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "typescript-bun-monorepo",
            )

            self.assertEqual(audit.returncode, 3, audit.stderr)
            self.assertEqual(apply.returncode, 3, apply.stderr)
            self.assertEqual(after_audit, before)
            self.assertEqual(snapshot_files(repo), before)
            audit_payload = json.loads(audit.stdout)
            apply_payload = json.loads(apply.stdout)
            self.assertEqual(audit_payload["selection"], apply_payload["selection"])
            self.assertEqual(audit_payload["plannedChanges"], apply_payload["plannedChanges"])

    def test_preview_output_is_deterministic(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            first = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "typescript-bun-monorepo",
            )
            second = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "typescript-bun-monorepo",
            )

            self.assertEqual(first.returncode, 3, first.stderr)
            self.assertEqual(second.returncode, 3, second.stderr)
            self.assertEqual(json.loads(first.stdout), json.loads(second.stdout))

    def test_text_preview_names_selection_and_planned_changes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            result = run_context_setup(
                "audit",
                "--repo",
                temp_dir,
                "--format",
                "text",
                "--profile",
                "rust-cli",
            )

            self.assertEqual(result.returncode, 3, result.stderr)
            self.assertIn("selection:", result.stdout)
            self.assertIn("- profile rust-cli setup rust-cli", result.stdout)
            self.assertIn("planned changes:", result.stdout)
            self.assertIn("state=conditional", result.stdout)

    def _planned_change(self, planned, managed_id, state):
        for change in planned:
            if change["managedId"] == managed_id and change["state"] == state:
                return change
        self.fail(f"missing {state} planned change for {managed_id}: {planned}")


def run_context_setup(*args):
    env = os.environ.copy()
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        text=True,
        capture_output=True,
        check=False,
        env=env,
    )


if __name__ == "__main__":
    unittest.main()
