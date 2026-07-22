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

    def test_operator_docs_publish_schema_exit_and_confirmation_contract(self):
        skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")
        user_guide = USER_GUIDE.read_text(encoding="utf-8")

        for document in (skill, user_guide):
            self.assertIn("setup-context-driven/audit-v1", document)
            self.assertIn("setup-context-driven/restore-v1", document)
            self.assertIn("context_setup.py audit", document)
            self.assertIn("context_setup.py apply", document)
            self.assertIn("--confirm-plan", document)
            self.assertIn("restore-skills", document)
            self.assertIn("exit `0`", document)
            self.assertIn("exit `1`", document)
            self.assertIn("exit `2`", document)
            self.assertIn("exit `3`", document)
            for field in (
                "action",
                "path",
                "managedId",
                "state",
                "reason",
                "beforeDigest",
                "afterDigest",
                "condition",
                "fromPath",
                "referenceEdits",
            ):
                self.assertIn(f"`{field}`", document)
            self.assertNotIn("bunx skills update", document)
            self.assertNotIn("bunx skills experimental_install", document)

        audit_index = skill.index("context_setup.py audit")
        resolve_index = skill.index("fully resolved Decision Plan")
        apply_index = skill.index("context_setup.py apply")
        final_audit_index = skill.index("After apply, rerun the same resolved audit")
        restore_index = skill.index("context_setup.py restore-skills")
        self.assertLess(audit_index, resolve_index)
        self.assertLess(resolve_index, apply_index)
        self.assertLess(apply_index, final_audit_index)
        self.assertLess(final_audit_index, restore_index)
        self.assertIn("never removes skills", skill)
        self.assertIn("never generates project-specific architecture", skill)
        self.assertIn("never execute that argv automatically", skill)
        self.assertIn("There is no branch, default-revision, or generic skill-refresh fallback", skill)

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
