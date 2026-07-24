import json
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_support import repository_root  # noqa: E402

REPO_ROOT = repository_root(SKILL_ROOT)
USER_GUIDE = REPO_ROOT / "docs" / "user-guide" / "context-driven-development.md"
ASSET_GUIDE = SKILL_ROOT / "references" / "asset-maintenance.md"

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
                    "repository.extension.enabled=false",
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
        normalized = " ".join(skill.split())

        plan_index = skill.index("roundfix baseline plan")
        preview_index = skill.index("Review at least:")
        confirm_index = skill.index("Ask the maintainer to approve")
        apply_index = skill.index("roundfix baseline apply")

        self.assertLess(plan_index, preview_index)
        self.assertLess(preview_index, confirm_index)
        self.assertLess(confirm_index, apply_index)
        self.assertIn("Baseline owns declared root blocks", skill)
        self.assertIn(
            "preserves repository-authored bytes outside managed boundaries",
            normalized,
        )

    def test_skill_names_canonical_setup_and_optional_extra_skill_report(self):
        skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")
        normalized = " ".join(skill.split())

        self.assertIn("## Baseline Profiles", skill)
        self.assertIn("profile show", skill)
        self.assertIn("profile validate", skill)
        self.assertIn("## Repository Skill Set restoration", skill)
        self.assertIn("never restores them as a side effect", normalized)
        self.assertIn("## Canonical asset synchronization", skill)

    def test_operator_docs_publish_schema_exit_and_confirmation_contract(self):
        skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")
        user_guide = USER_GUIDE.read_text(encoding="utf-8")

        for document in (skill, user_guide):
            normalized = " ".join(document.split())
            self.assertIn("roundfix/baseline-plan/v1", document)
            self.assertIn("roundfix/baseline-result/v1", document)
            self.assertIn("roundfix baseline plan", document)
            self.assertIn("roundfix baseline apply", document)
            self.assertIn("--confirm-plan", document)
            self.assertIn("baseline skills restore", document)
            for exit_code in ("`0`", "`1`", "`2`", "`3`"):
                self.assertIn(exit_code, normalized)
            for field in (
                "fileChanges",
                "managedEntries",
                "preimage",
                "postimage",
                "planDigest",
                "recommendations",
            ):
                self.assertIn(field, document)
            self.assertNotIn("context_setup.py", document)
            self.assertNotIn("python3", document.lower())

        plan_index = skill.index("roundfix baseline plan")
        review_index = skill.index("Review at least:")
        apply_index = skill.index("roundfix baseline apply")
        restore_index = skill.index("roundfix baseline skills restore")
        completion_index = skill.index("## Completion")
        self.assertLess(plan_index, review_index)
        self.assertLess(review_index, apply_index)
        self.assertLess(apply_index, restore_index)
        self.assertLess(restore_index, completion_index)
        normalized_skill = " ".join(skill.split())
        self.assertIn("never restores them as a side effect", normalized_skill)
        self.assertIn(
            "no independent setup engine or behavioral fallback",
            normalized_skill,
        )
        self.assertIn("Never substitute a generic skill refresh", normalized_skill)

    def test_asset_maintenance_doc_publishes_catalog_and_source_boundaries(self):
        guide = ASSET_GUIDE.read_text(encoding="utf-8")

        for token in (
            "assets/coverage.json",
            "requiredRules",
            "skillDispatch",
            "typed references",
            "setup-context-driven/setup-snapshot-v2",
            "treeDigest",
            "sync-setups",
            "make skills-sync",
            ".agents/skills/setup-context-driven",
            "skills/setup-context-driven",
        ):
            self.assertIn(token, guide)
        self.assertIn("Do not edit upstream-managed skill content", guide)
        self.assertIn("Spec 0036", guide)
        self.assertIn("Doctor", guide)


def run_context_setup(*args):
    return run_fixture_context_setup(*args)


if __name__ == "__main__":
    unittest.main()
