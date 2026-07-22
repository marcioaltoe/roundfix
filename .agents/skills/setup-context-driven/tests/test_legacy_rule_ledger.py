# Suite: reviewed legacy hard-rule ledgers
# Invariant: every reviewed clause has one strength-preserving mapping or one recorded exclusion.
# Boundary IN: canonical rule assets, transition ledgers, and generated guidance.
# Boundary OUT: upgrade planning, manifest migration, formatter execution, and network access.

import json
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
from test_decision_rendering import decisions_for, run_apply_for_profile  # noqa: E402


LEDGER_IDS = {
    "transition.managed-v2-to-portable-v3",
    "transition.legacy-typescript-bun-to-portable-v3",
}

RETIRED_LEGACY_CLAUSES = {
    "clause.legacy.coding-reference-trio",
    "clause.legacy.knowledge-workspace-sparse-checkout",
    "clause.legacy.separate-documentation-commit",
    "clause.legacy.fable-runtime-defaults",
    "clause.legacy.template-product-stack",
    "clause.legacy.standard-rest-contract",
    "clause.legacy.feature-systems-layout",
    "clause.legacy.unconfirmed-domain-skill-matrix",
    "clause.legacy.code-language-policy",
    "clause.legacy.environment-example-parity",
}


class LegacyRuleLedgerTests(unittest.TestCase):
    def test_every_reviewed_clause_has_one_mapping_with_equivalent_enforcement(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        self.assertEqual(set(catalog.upgrade_transitions), LEDGER_IDS)
        current_clauses = {
            clause.clause_id: clause
            for rule in catalog.rule_contracts.values()
            for clause in rule.clauses
        }
        for ledger in catalog.upgrade_transitions.values():
            prior = {clause.clause_id: clause for clause in ledger.prior_clauses}
            mappings = {mapping.from_clause: mapping for mapping in ledger.mappings}
            self.assertEqual(set(mappings), set(prior), ledger.transition_id)
            self.assertEqual(len(ledger.mappings), len(set(mappings)), ledger.transition_id)
            for mapping in ledger.mappings:
                self.assertTrue(mapping.reason.strip(), mapping.from_clause)
                if mapping.disposition == "rejected":
                    self.assertEqual(mapping.targets, (), mapping.from_clause)
                    continue
                self.assertTrue(mapping.targets, mapping.from_clause)
                for target in mapping.targets:
                    if target.startswith("clause."):
                        self.assertEqual(
                            current_clauses[target].enforcement,
                            prior[mapping.from_clause].enforcement,
                            mapping.from_clause,
                        )

    def test_missing_mapping_and_weakened_accepted_mapping_fail_validation(self):
        cases = [
            (self._remove_mapping, "transition.mapping.incomplete"),
            (self._weaken_target, "transition.target.enforcement.mismatch"),
        ]

        for mutator, diagnostic in cases:
            with self.subTest(diagnostic=diagnostic):
                with tempfile.TemporaryDirectory() as temp_dir:
                    root = Path(temp_dir) / "setup-context-driven"
                    clone_assets_to(SKILL_ROOT, root)
                    mutator(root)
                    with self.assertRaises(AssetValidationError) as captured:
                        load_asset_catalog(root)
                self.assertIn(diagnostic, "\n".join(captured.exception.diagnostics))

    def test_retired_sample_behavior_is_excluded_and_repository_policy_survives(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        legacy = catalog.upgrade_transitions[
            "transition.legacy-typescript-bun-to-portable-v3"
        ]
        rejected = {
            mapping.from_clause
            for mapping in legacy.mappings
            if mapping.disposition == "rejected"
        }
        self.assertEqual(rejected, RETIRED_LEGACY_CLAUSES)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            authored = "# Repository policy\n\nKeep this repository-authored rule.\n"
            (repo / "AGENTS.md").write_text(authored, encoding="utf-8")
            (repo / "DESIGN.md").write_text("# Design contract\n", encoding="utf-8")

            result = run_apply_for_profile(
                repo,
                "typescript-bun-monorepo",
                decisions_for(autonomous=False),
            )

            self.assertEqual(result.returncode, 0, result.stderr)
            root_guidance = (repo / "AGENTS.md").read_text(encoding="utf-8")
            self.assertIn(authored, root_guidance)
            generated = root_guidance + "\n" + "\n".join(
                path.read_text(encoding="utf-8")
                for path in sorted((repo / "docs" / "agents").glob("*.md"))
            )
            for retired_text in (
                ".knowledge/",
                "Fable",
                "gpt-5.5",
                "Opus 4.8",
                "<Project name>",
                "PostgreSQL 18",
                "Better Auth",
                "standard REST",
                "systems/<domain>",
            ):
                self.assertNotIn(retired_text, generated)

    def _remove_mapping(self, root):
        path = (
            root
            / "assets"
            / "retention"
            / "transition.managed-v2-to-portable-v3.json"
        )
        data = read_json_copy(path)
        data["mappings"].pop()
        write_json(path, data)

    def _weaken_target(self, root):
        path = root / "assets" / "modules" / "core.json"
        data = read_json_copy(path)
        for rule in data["rules"]:
            for clause in rule["clauses"]:
                if clause["id"] == "clause.core.require-fresh-evidence":
                    clause["enforcement"] = "stop-and-ask"
                    write_json(path, data)
                    return
        self.fail("missing clause.core.require-fresh-evidence")


if __name__ == "__main__":
    unittest.main()
