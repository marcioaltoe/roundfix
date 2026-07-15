import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    read_json_copy,
    write_json,
)


class AssetContractTests(unittest.TestCase):
    def test_supported_profiles_resolve_to_deterministic_module_order(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        self.assertEqual(
            catalog.ordered_modules_by_profile,
            {
                "go-cli-tui": [
                    "core",
                    "context-workflow",
                    "go",
                    "cli-surface",
                    "tui-surface",
                    "autonomous-work",
                ],
                "rust-cli": [
                    "core",
                    "context-workflow",
                    "rust",
                    "cli-surface",
                    "autonomous-work",
                ],
                "typescript-bun-monorepo": [
                    "core",
                    "context-workflow",
                    "typescript",
                    "bun",
                    "monorepo",
                    "backend",
                    "frontend",
                    "autonomous-work",
                ],
            },
        )

    def test_every_managed_asset_has_stable_id_and_version(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        for collection in [
            catalog.decisions,
            catalog.modules,
            catalog.profiles,
            catalog.setups,
            catalog.templates,
        ]:
            for asset_id, asset in collection.items():
                self.assertEqual(asset["id"], asset_id)
                self.assertIsInstance(asset["version"], int)
                self.assertGreaterEqual(asset["version"], 1)

        for module in catalog.modules.values():
            for rule in module["rules"]:
                self.assertRegex(rule["id"], r"^rule\.")
                self.assertIsInstance(rule["version"], int)
            for block in module["rootBlocks"]:
                self.assertRegex(block["id"], r"^root\.")
                self.assertIsInstance(block["version"], int)
            for guide in module["supportingGuides"]:
                self.assertRegex(guide["id"], r"^guide\.")
                self.assertIsInstance(guide["version"], int)

    def test_profile_referenced_skills_exist_in_bundled_setup_snapshot(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        for profile_id, module_ids in catalog.ordered_modules_by_profile.items():
            setup = catalog.setups[catalog.profiles[profile_id]["setup"]]
            setup_skills = {skill["name"] for skill in setup["skills"]}
            required_skills = {
                skill
                for module_id in module_ids
                for skill in catalog.modules[module_id]["requiredSkills"]
            }
            self.assertLessEqual(required_skills, setup_skills)

    def test_root_templates_are_compact_pointers(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        templates_root = SKILL_ROOT / "assets" / "templates"

        root_templates = [
            template
            for template in catalog.templates.values()
            if template["kind"] == "root-block"
        ]
        self.assertGreater(len(root_templates), 0)

        for template in root_templates:
            content = (templates_root / template["path"]).read_text(encoding="utf-8")
            self.assertLessEqual(len(content.split()), 45, template["id"])
            self.assertIn("docs/agents/", content, template["id"])

    def test_assets_are_portable(self):
        load_asset_catalog(SKILL_ROOT)

        assets_root = SKILL_ROOT / "assets"
        for path in assets_root.rglob("*"):
            if not path.is_file():
                continue
            content = path.read_text(encoding="utf-8")
            self.assertNotIn("~/dev/skills", content)
            self.assertNotIn("https://", content)
            self.assertNotIn("http://", content)

    def test_invalid_contract_fixtures_fail_with_precise_diagnostics(self):
        cases = [
            ("unknown profile module", self._unknown_profile_module, "profile.module.unknown"),
            ("unknown setup", self._unknown_setup, "profile.setup.unknown"),
            (
                "missing profile dependency",
                self._missing_profile_dependency,
                "profile.dependency.missing",
            ),
            ("dependency cycle", self._dependency_cycle, "module.dependency.cycle"),
            ("conflicting modules", self._conflicting_modules, "module.conflict"),
            ("duplicate rule id", self._duplicate_rule_id, "rule.id.duplicate"),
            (
                "duplicate root block id",
                self._duplicate_root_block_id,
                "managed.block.id.duplicate",
            ),
            ("unknown decision", self._unknown_decision, "module.decision.unknown"),
            ("unknown template", self._unknown_template, "template.reference.unknown"),
            (
                "skill outside snapshot",
                self._missing_setup_skill,
                "skills.reference.outside-setup",
            ),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                diagnostics = self._load_invalid_fixture(mutator)
                self.assertIn(expected_code, diagnostics)

    def _load_invalid_fixture(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_root)
            mutator(temp_root)

            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)

        return "\n".join(captured.exception.diagnostics)

    def _unknown_profile_module(self, temp_root):
        profile_path = temp_root / "assets" / "profiles" / "rust-cli.json"
        profile = read_json_copy(profile_path)
        profile["modules"].append("missing-module")
        write_json(profile_path, profile)

    def _unknown_setup(self, temp_root):
        profile_path = temp_root / "assets" / "profiles" / "rust-cli.json"
        profile = read_json_copy(profile_path)
        profile["setup"] = "missing-setup"
        write_json(profile_path, profile)

    def _missing_profile_dependency(self, temp_root):
        profile_path = temp_root / "assets" / "profiles" / "go-cli-tui.json"
        profile = read_json_copy(profile_path)
        profile["modules"].remove("core")
        write_json(profile_path, profile)

    def _dependency_cycle(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "core.json"
        module = read_json_copy(module_path)
        module["dependsOn"] = ["context-workflow"]
        write_json(module_path, module)

    def _conflicting_modules(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "backend.json"
        module = read_json_copy(module_path)
        module["conflictsWith"] = ["frontend"]
        write_json(module_path, module)

    def _duplicate_rule_id(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "bun.json"
        module = read_json_copy(module_path)
        module["rules"][0]["id"] = "rule.typescript.current-docs"
        write_json(module_path, module)

    def _duplicate_root_block_id(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "bun.json"
        module = read_json_copy(module_path)
        module["rootBlocks"][0]["id"] = "root.typescript"
        write_json(module_path, module)

    def _unknown_decision(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "rust.json"
        module = read_json_copy(module_path)
        module["requiredDecisions"] = ["decision.missing"]
        write_json(module_path, module)

    def _unknown_template(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "rust.json"
        module = read_json_copy(module_path)
        module["rootBlocks"][0]["template"] = "template.missing"
        write_json(module_path, module)

    def _missing_setup_skill(self, temp_root):
        setup_path = temp_root / "assets" / "setups" / "go-cli.json"
        setup = read_json_copy(setup_path)
        setup["skills"] = [
            skill for skill in setup["skills"] if skill["name"] != "golang-cli"
        ]
        write_json(setup_path, setup)


if __name__ == "__main__":
    unittest.main()
