import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_apply import BASE_DECISIONS, run_apply  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_context_setup as run_fixture_context_setup,
    snapshot_files,
    write_compliant_repository,
)


class SetupWorkflowTests(unittest.TestCase):
    def test_unchanged_rerun_reuses_stored_compatible_decisions(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            first = run_apply(repo, "rust-cli", BASE_DECISIONS)
            install_profile_skills(repo, "rust-cli")
            after_first = snapshot_files(repo)

            rerun = run_apply(repo, "rust-cli", [])

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(rerun.returncode, 0, rerun.stderr)
            self.assertEqual(snapshot_files(repo), after_first)
            audit = run_context_setup("audit", "--repo", str(repo), "--format", "json")
            self.assertEqual(audit.returncode, 0, audit.stderr)

    def test_new_required_decision_routes_one_question_and_persists_answer(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli", omit_decision="secondbrain.enabled")

            blocked = run_context_setup("audit", "--repo", str(repo), "--format", "json")
            applied = run_apply(repo, "rust-cli", ["secondbrain.enabled=false"])

            self.assertEqual(blocked.returncode, 3)
            payload = json.loads(blocked.stdout)
            decision_findings = [
                finding
                for finding in payload["findings"]
                if finding["code"] == "decision.required"
            ]
            self.assertEqual(len(decision_findings), 1, payload)
            self.assertEqual(decision_findings[0]["managedId"], "secondbrain.enabled")
            self.assertEqual(applied.returncode, 0, applied.stderr)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(
                manifest["decisions"]["secondbrain.enabled"]["value"],
                False,
            )

    def test_enabled_autonomous_work_routes_one_missing_dependent_question(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            disabled = run_apply(
                repo,
                "rust-cli",
                [
                    "spec.scaffold=true",
                    "domain.layout=single-context",
                    "triage.external=false",
                    "autonomous.enabled=false",
                    "verification.gate=make workflow-verify",
                    "language.generated=English",
                    "secondbrain.enabled=false",
                ],
            )
            before = snapshot_files(repo)

            blocked = run_apply(
                repo,
                "rust-cli",
                [
                    "autonomous.enabled=true",
                    "runtime.backend=codex workflow-backend xhigh",
                ],
            )
            blocked_payload = json.loads(blocked.stdout)
            required = [
                finding["managedId"]
                for finding in blocked_payload["findings"]
                if finding["code"] == "decision.required"
            ]

            self.assertEqual(disabled.returncode, 0, disabled.stderr)
            self.assertEqual(blocked.returncode, 3)
            self.assertEqual(required, ["runtime.design"], blocked_payload)
            self.assertEqual(snapshot_files(repo), before)

            answered = run_apply(
                repo,
                "rust-cli",
                [
                    "autonomous.enabled=true",
                    "runtime.backend=codex workflow-backend xhigh",
                    "runtime.design=claude workflow-design xhigh",
                    "verification.gate=make workflow-verify",
                ],
            )

            self.assertEqual(answered.returncode, 0, answered.stderr)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(manifest["decisions"]["verification.gate"]["value"], "make workflow-verify")

    def test_skill_workflow_requires_preview_and_confirmation_before_apply(self):
        skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")

        audit_index = skill.index("Run audit before asking setup questions")
        preview_index = skill.index("plannedChanges")
        confirm_index = skill.index("Ask for confirmation")
        apply_index = skill.index("context_setup.py apply")

        self.assertLess(audit_index, preview_index)
        self.assertLess(preview_index, confirm_index)
        self.assertLess(confirm_index, apply_index)
        self.assertIn("only `docs/agents/setup-context.json` and declared setup-owned Markdown boundaries can change", skill)
        self.assertIn("repository-authored bytes outside managed markers remain untouched", skill)

    def test_skill_names_canonical_setup_and_optional_extra_skill_report(self):
        skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")

        self.assertIn("selected canonical skill setup", skill)
        self.assertIn("--show-extra-skills", skill)
        self.assertIn("Never suggest a removal command", skill)
        self.assertIn("Do not dump generated Markdown by default", skill)


def run_context_setup(*args):
    return run_fixture_context_setup(*args)


if __name__ == "__main__":
    unittest.main()
