# Suite: generated operational guidance
# Invariant: each operational clause renders through one distinct supporting guide with canonical enforcement.
# Boundary IN: setup asset validation and generated workflow-guide content.
# Boundary OUT: retention accounting, dispatch normalization, formatter execution, and repository delegation scans.

import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    read_json_copy,
    write_json,
)
from test_apply import BASE_DECISIONS, run_apply  # noqa: E402


class OperationalGuidanceTests(unittest.TestCase):
    def test_selected_workflow_guides_render_complete_distinct_contracts(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            decisions = [
                "secondbrain.enabled=true"
                if decision.startswith("secondbrain.enabled=")
                else decision
                for decision in BASE_DECISIONS
            ]

            result = run_apply(repo, "rust-cli", decisions)

            self.assertEqual(result.returncode, 0, result.stderr)
            spec_routing = self.read_guide(repo, "spec-routing.md")
            issue_tracker = self.read_guide(repo, "issue-tracker.md")
            docs_layout = self.read_guide(repo, "docs-layout.md")
            general = self.read_guide(repo, "agent-instructions.md")
            autonomous = self.read_guide(repo, "autonomous-work.md")
            secondbrain = self.read_guide(repo, "secondbrain.md")

            for phrase in (
                "Large or fuzzy product initiative",
                "Standard feature that changes product behavior",
                "Refactor or bug fix without product-behavior change",
                "Trivial one-line fix, typo, or configuration tweak",
            ):
                self.assertIn(phrase, spec_routing)
                self.assertNotIn(phrase, issue_tracker)
            self.assertIn("Dependencies live only in `_tasks.md`", issue_tracker)
            self.assertNotEqual(spec_routing, issue_tracker)

            self.assertIn("status: pending # pending | partial | deferred | done", docs_layout)
            for state in ("`pending`", "`partial`", "`deferred`"):
                self.assertIn(state, docs_layout)
            self.assertIn("`status: done`", docs_layout)
            self.assertIn("Record the reason", docs_layout)
            self.assertIn("append evidence and routing links", docs_layout)
            self.assertIn("Update `updated_at`", docs_layout)

            self.assertIn("evidence for every acceptance criterion", general)
            self.assertIn("Keep follow-up work outside the current slice", general)

            self.assertIn("must not write feature code or tests", autonomous)
            self.assertIn("through a Roundfix Run", autonomous)
            self.assertIn("Default backend work uses `codex gpt-5.5 xhigh`", autonomous)
            self.assertIn("frontend-dominant work uses `claude opus xhigh`", autonomous)

            self.assert_in_order(
                secondbrain,
                "wiki/index.md",
                'qmd query "<question>" --all --files --min-score 0.3',
                "projects/<project>/mirror/",
            )
            for phrase in (
                "Vortex, Tax, Visio, or Gesttione",
                "Cite every Secondbrain file",
                "Do not write to the Secondbrain",
                "Do not edit raw/",
                "Do not edit projects/*/mirror/",
                "Never read, copy, or expose",
            ):
                self.assertIn(phrase, secondbrain)

            for guide in (
                spec_routing,
                issue_tracker,
                docs_layout,
                general,
                autonomous,
                secondbrain,
            ):
                self.assertRegex(
                    guide,
                    r"(?m)^- \*\*(mandatory|prohibited|stop-and-ask)\*\*: ",
                )

    def test_duplicate_effective_clause_lists_are_rejected(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            skill_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, skill_root)
            module_path = skill_root / "assets" / "modules" / "spec-workflow.json"
            module = read_json_copy(module_path)
            guides = module["supportingGuides"]
            routing_rules = next(
                guide["rules"] for guide in guides if guide["id"] == "guide.spec-routing"
            )
            next(
                guide for guide in guides if guide["id"] == "guide.issue-tracker"
            )["rules"] = list(routing_rules)
            write_json(module_path, module)

            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(skill_root)

        self.assertIn(
            "guide.clause-list.duplicate: spec-workflow",
            "\n".join(captured.exception.diagnostics),
        )

    def read_guide(self, repo, name):
        return (repo / "docs" / "agents" / name).read_text(encoding="utf-8")

    def assert_in_order(self, content, *phrases):
        positions = [content.index(phrase) for phrase in phrases]
        self.assertEqual(positions, sorted(positions))


if __name__ == "__main__":
    unittest.main()
