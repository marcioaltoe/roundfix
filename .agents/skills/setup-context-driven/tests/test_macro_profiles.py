"""Macro profile flows for setup-context-driven.

Suite: supported profile repository flows
Invariant: every bundled profile applies, audits cleanly, and re-applies without file changes.
Boundary IN: context_setup.py CLI, bundled assets, temporary repository files.
Boundary OUT: Makefile orchestration and embedded skill synchronization checks.
"""

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_apply import BASE_DECISIONS, run_apply, run_audit  # noqa: E402
from test_audit import install_profile_skills, snapshot_files, write_skill  # noqa: E402
from test_skills import write_lockfile  # noqa: E402


SUPPORTED_PROFILES = [
    "typescript-bun-monorepo",
    "go-cli-tui",
    "rust-cli",
]


class ProfileMacroFlowTests(unittest.TestCase):
    def test_supported_profiles_apply_audit_clean_and_reapply_without_changes(self):
        for profile_id in SUPPORTED_PROFILES:
            with self.subTest(profile=profile_id):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)

                    first_apply = run_apply(repo, profile_id, BASE_DECISIONS)
                    install_profile_skills(repo, profile_id)
                    clean_audit = run_audit(repo)
                    after_audit = snapshot_files(repo)
                    second_apply = run_apply(repo, profile_id, [])

                    self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
                    self.assertEqual(clean_audit.returncode, 0, clean_audit.stderr)
                    self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
                    self.assertEqual(snapshot_files(repo), after_audit)

    def test_required_skill_failure_and_extra_reporting_keep_exit_semantics(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            first_apply = run_apply(repo, "rust-cli", BASE_DECISIONS)
            install_profile_skills(repo, "rust-cli", omit={"agentic-cli-design"})
            missing_required = run_audit(repo)
            write_skill(repo, "agentic-cli-design")
            write_skill(repo, "the-fool")
            write_lockfile(repo, ["the-fool"])
            compliant_without_extras = run_audit(repo)
            extras_visible = run_audit_cli(repo, "--show-extra-skills")

            self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
            self.assertEqual(missing_required.returncode, 1)
            self.assertFinding(missing_required, "skills.required.missing", "error")
            self.assertEqual(compliant_without_extras.returncode, 0, compliant_without_extras.stderr)
            self.assertNoFinding(compliant_without_extras, "skills.extra.installed")
            self.assertEqual(extras_visible.returncode, 0, extras_visible.stderr)
            extra = self.finding(extras_visible, "skills.extra.installed")
            self.assertEqual(extra["severity"], "info")

    def test_secondbrain_opt_in_profile_applies_audits_and_reapplies_without_changes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            first_apply = run_apply(repo, "go-cli-tui", secondbrain_decisions())
            install_profile_skills(repo, "go-cli-tui")
            clean_audit = run_audit(repo)
            after_audit = snapshot_files(repo)
            second_apply = run_apply(repo, "go-cli-tui", [])

            self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
            self.assertEqual(clean_audit.returncode, 0, clean_audit.stderr)
            self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
            self.assertEqual(snapshot_files(repo), after_audit)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertIn("secondbrain", manifest["modules"])
            self.assertIn("root.secondbrain", (repo / "AGENTS.md").read_text(encoding="utf-8"))
            self.assertTrue((repo / "docs" / "agents" / "secondbrain.md").is_file())

    def finding(self, result, code):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        return matches[0]

    def assertFinding(self, result, code, severity):
        match = self.finding(result, code)
        self.assertEqual(match["severity"], severity)

    def assertNoFinding(self, result, code):
        payload = json.loads(result.stdout)
        self.assertEqual(
            [finding for finding in payload["findings"] if finding["code"] == code],
            [],
            payload,
        )


def secondbrain_decisions():
    return [
        "secondbrain.enabled=true"
        if decision.startswith("secondbrain.enabled=")
        else decision
        for decision in BASE_DECISIONS
    ]


def run_audit_cli(repo, *extra_args):
    return subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            "audit",
            "--repo",
            str(repo),
            "--format",
            "json",
            *extra_args,
        ],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()
