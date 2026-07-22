import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_setup import parse_managed_blocks  # noqa: E402
from test_audit import install_profile_skills, snapshot_files  # noqa: E402


class AutonomousSecondbrainDecisionTests(unittest.TestCase):
    def test_autonomous_disabled_omits_guidance_and_dependent_questions(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(
                repo,
                decisions_for(
                    autonomous=False,
                    secondbrain=False,
                    include_runtime=False,
                ),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = read_manifest(repo)
            self.assertNotIn("autonomous-work", manifest["modules"])
            self.assertNotIn("runtime.backend", manifest["decisions"])
            self.assertNotIn("runtime.design", manifest["decisions"])
            self.assertEqual(
                manifest["decisions"]["verification.gate"]["value"],
                "make verify",
            )
            self.assertNotIn(
                "root.autonomous-work",
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertFalse((repo / "docs" / "agents" / "autonomous-work.md").exists())

    def test_autonomous_enabled_requires_only_runtime_and_verification_values(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            blocked = run_apply(
                repo,
                decisions_for(
                    autonomous=True,
                    secondbrain=False,
                    include_runtime=False,
                ),
            )

            self.assertEqual(blocked.returncode, 3)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(blocked.stdout)
            missing = {
                finding["managedId"]
                for finding in payload["findings"]
                if finding["code"] == "decision.required"
            }
            self.assertEqual(missing, {"runtime.backend", "runtime.design"})

            applied = run_apply(repo, decisions_for(autonomous=True, secondbrain=False))

            self.assertEqual(applied.returncode, 0, applied.stderr)
            manifest = read_manifest(repo)
            self.assertIn("autonomous-work", manifest["modules"])
            root_blocks, _ = parse_managed_blocks(
                Path("AGENTS.md"),
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertIn("root.autonomous-work", root_blocks)
            self.assertTrue((repo / "docs" / "agents" / "autonomous-work.md").is_file())

    def test_unanswered_optional_decisions_preview_conditional_modules_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            result = run_context_setup(
                "audit",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "rust-cli",
            )

            self.assertEqual(result.returncode, 3, result.stderr)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(result.stdout)
            modules = {module["id"]: module for module in payload["selection"]["modules"]}
            self.assertEqual(modules["autonomous-work"]["state"], "conditional")
            self.assertEqual(modules["secondbrain"]["state"], "conditional")
            self.assert_planned_condition(
                payload["plannedChanges"],
                "guide.autonomous-work",
                "autonomous.enabled",
            )
            self.assert_planned_condition(
                payload["plannedChanges"],
                "guide.secondbrain",
                "secondbrain.enabled",
            )

    def test_secondbrain_enabled_preserves_complete_safety_guidance(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(repo, decisions_for(autonomous=False, secondbrain=True))
            install_profile_skills(repo, "rust-cli")

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = read_manifest(repo)
            self.assertIn("secondbrain", manifest["modules"])
            root_blocks, _ = parse_managed_blocks(
                Path("AGENTS.md"),
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertIn("root.secondbrain", root_blocks)
            self.assertLessEqual(len(root_blocks["root.secondbrain"].body.split()), 45)
            guide = (repo / "docs" / "agents" / "secondbrain.md").read_text(
                encoding="utf-8"
            )
            for phrase in SECONDBRAIN_REQUIRED_GUIDE_PHRASES:
                self.assertIn(phrase, guide)
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_disabling_optional_modules_preserves_surrounding_owner_bytes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            enabled = run_apply(repo, decisions_for(autonomous=True, secondbrain=True))
            self.assertEqual(enabled.returncode, 0, enabled.stderr)

            root_path = repo / "AGENTS.md"
            autonomous_guide = repo / "docs" / "agents" / "autonomous-work.md"
            secondbrain_guide = repo / "docs" / "agents" / "secondbrain.md"
            root_path.write_text(
                "repo root note\n"
                + root_path.read_text(encoding="utf-8")
                + "repo root tail\n",
                encoding="utf-8",
            )
            autonomous_guide.write_text(
                "repo autonomous note\n"
                + autonomous_guide.read_text(encoding="utf-8")
                + "repo autonomous tail\n",
                encoding="utf-8",
            )
            secondbrain_guide.write_text(
                "repo secondbrain note\n"
                + secondbrain_guide.read_text(encoding="utf-8")
                + "repo secondbrain tail\n",
                encoding="utf-8",
            )

            disabled = run_apply(repo, decisions_for(autonomous=False, secondbrain=False))

            self.assertEqual(disabled.returncode, 0, disabled.stderr)
            manifest = read_manifest(repo)
            self.assertNotIn("autonomous-work", manifest["modules"])
            self.assertNotIn("secondbrain", manifest["modules"])
            root_content = root_path.read_text(encoding="utf-8")
            self.assertIn("repo root note\n", root_content)
            self.assertIn("repo root tail\n", root_content)
            self.assertNotIn("id=root.autonomous-work", root_content)
            self.assertNotIn("id=root.secondbrain", root_content)
            self.assertEqual(
                autonomous_guide.read_text(encoding="utf-8"),
                "repo autonomous note\nrepo autonomous tail\n",
            )
            self.assertEqual(
                secondbrain_guide.read_text(encoding="utf-8"),
                "repo secondbrain note\nrepo secondbrain tail\n",
            )

    def test_apply_and_audit_are_idempotent_for_enabled_and_disabled_branches(self):
        cases = [
            (False, False),
            (False, True),
            (True, False),
            (True, True),
        ]
        for autonomous_enabled, secondbrain_enabled in cases:
            with self.subTest(
                autonomous=autonomous_enabled,
                secondbrain=secondbrain_enabled,
            ):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    first = run_apply(
                        repo,
                        decisions_for(
                            autonomous=autonomous_enabled,
                            secondbrain=secondbrain_enabled,
                        ),
                    )
                    self.assertEqual(first.returncode, 0, first.stderr)
                    install_profile_skills(repo, "rust-cli")
                    after_first = snapshot_files(repo)

                    second = run_apply(
                        repo,
                        decisions_for(
                            autonomous=autonomous_enabled,
                            secondbrain=secondbrain_enabled,
                        ),
                    )
                    first_audit = run_audit(repo)
                    second_audit = run_audit(repo)

                    self.assertEqual(second.returncode, 0, second.stderr)
                    self.assertEqual(snapshot_files(repo), after_first)
                    self.assertEqual(first_audit.returncode, 0, first_audit.stderr)
                    self.assertEqual(second_audit.returncode, 0, second_audit.stderr)
                    self.assertEqual(
                        json.loads(first_audit.stdout),
                        json.loads(second_audit.stdout),
                    )

    def assert_planned_condition(self, planned_changes, managed_id, decision_id):
        matches = [
            change
            for change in planned_changes
            if change["managedId"] == managed_id and change["state"] == "conditional"
        ]
        self.assertEqual(len(matches), 1, planned_changes)
        self.assertEqual(
            matches[0]["condition"],
            {"decisionId": decision_id, "equals": True},
        )


SECONDBRAIN_REQUIRED_GUIDE_PHRASES = [
    "wiki/index.md",
    "qmd query",
    "projects/<project>/mirror",
    "Cite every Secondbrain file",
    "Do not write to the Secondbrain",
    "Do not edit raw/",
    "Do not edit projects/*/mirror/",
    "Hermes",
    "Never read, copy, or expose",
]


def decisions_for(autonomous, secondbrain, include_runtime=True):
    decisions = [
        "spec.scaffold=true",
        "domain.layout=single-context",
        "triage.external=false",
        f"autonomous.enabled={str(autonomous).lower()}",
        "verification.gate=make verify",
        "language.generated=English",
        f"secondbrain.enabled={str(secondbrain).lower()}",
    ]
    if include_runtime:
        decisions.extend(
            [
                "runtime.backend=codex gpt-5.5 xhigh",
                "runtime.design=claude opus xhigh",
            ]
        )
    return decisions


def run_apply(repo, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", "rust-cli"]
    for decision in decisions:
        args.extend(["--decision", decision])
    result = run_context_setup(*args)
    payload = json.loads(result.stdout)
    if result.returncode == 3 and any(
        finding["code"] == "plan.confirmation.required" for finding in payload["findings"]
    ):
        args.extend(["--confirm-plan", payload["planDigest"]])
        return run_context_setup(*args)
    return result


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def read_manifest(repo):
    return json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
    )


def run_context_setup(*args):
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()
