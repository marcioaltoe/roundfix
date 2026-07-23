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
    "http.contract",
    "spec.scaffold",
    "domain.layout",
    "triage.external",
    "autonomous.enabled",
    "secondbrain.enabled",
    "repository.extension.enabled",
}


class PreviewCliTests(unittest.TestCase):
    def test_top_level_help_names_every_exit_three_condition(self):
        result = run_context_setup("--help")

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stderr, "")
        self.assertIn(
            "Exit codes: 0 ok, 1 blocking findings, 2 invalid input, "
            "3 decisions required or plan confirmation required/stale.",
            result.stdout.splitlines(),
        )

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
                "repository.extension.enabled=false",
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

    def test_structured_scalar_decisions_are_normalized_and_digest_bound(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir) / "repo"
            repo.mkdir()
            decision_path = Path(temp_dir) / "decisions.json"
            decisions = [
                {"id": "spec.scaffold", "value": True},
                {"id": "domain.layout", "value": "single-context"},
                {"id": "triage.external", "value": False},
                {"id": "autonomous.enabled", "value": False},
                {"id": "verification.gate", "value": "make verify"},
                {"id": "language.generated", "value": "English"},
                {"id": "secondbrain.enabled", "value": False},
                {"id": "repository.extension.enabled", "value": False},
            ]
            write_decision_document(decision_path, decisions)

            first = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "rust-cli",
                "--decision-file",
                str(decision_path),
            )
            decisions[2]["value"] = True
            write_decision_document(decision_path, decisions)
            changed = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "rust-cli",
                "--decision-file",
                str(decision_path),
            )

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(changed.returncode, 0, changed.stderr)
            first_payload = json.loads(first.stdout)
            changed_payload = json.loads(changed.stdout)
            self.assertEqual(
                first_payload["decisionDocument"]["schemaVersion"],
                "setup-context-driven/decisions/0.0.1",
            )
            self.assertNotEqual(first_payload["planDigest"], changed_payload["planDigest"])

    def test_conflicting_scalar_flag_and_decision_file_is_invalid_input(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir) / "repo"
            repo.mkdir()
            decision_path = Path(temp_dir) / "decisions.json"
            write_decision_document(
                decision_path,
                [{"id": "domain.layout", "value": "single-context"}],
            )

            result = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "rust-cli",
                "--decision-file",
                str(decision_path),
                "--decision",
                "domain.layout=multi-context",
            )

            self.assertEqual(result.returncode, 2, result.stderr)
            payload = json.loads(result.stdout)
            self.assertIn(
                "decision-file.decision.conflict",
                {item["code"] for item in payload["findings"]},
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
                "standard-typescript-monorepo",
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
            self.assertEqual(payload["selection"]["profile"], "standard-typescript-monorepo")
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
                "standard-typescript-monorepo",
            )
            after_audit = snapshot_files(repo)
            apply = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "standard-typescript-monorepo",
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
                "standard-typescript-monorepo",
            )
            second = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "standard-typescript-monorepo",
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


def write_decision_document(path, decisions):
    path.write_text(
        json.dumps(
            {
                "schemaVersion": "setup-context-driven/decisions/0.0.1",
                "version": "0.0.1",
                "decisions": decisions,
            },
            indent=2,
        )
        + "\n",
        encoding="utf-8",
    )


if __name__ == "__main__":
    unittest.main()
