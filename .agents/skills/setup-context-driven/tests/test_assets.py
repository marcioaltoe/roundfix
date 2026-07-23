# Suite: setup asset contracts
# Invariant: versioned catalogs accept only complete, safe, deterministic declarations.
# Boundary IN: local asset loading, normalization, and semantic validation.
# Boundary OUT: repository inspection, command execution, network access, and rendering.

import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
V2_FIXTURE_ROOT = Path(__file__).resolve().parent / "fixtures" / "asset-contracts-v2"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_support import setup_skill_roots  # noqa: E402

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    read_json_copy,
    write_json,
)


class AssetContractTests(unittest.TestCase):
    def test_versioned_contract_fixture_loads_normalized_values_deterministically(self):
        first = load_asset_catalog(V2_FIXTURE_ROOT)
        second = load_asset_catalog(V2_FIXTURE_ROOT)

        self.assertEqual(first, second)
        self.assertEqual(
            first.coverage_contracts["coverage.universal-safety"].description,
            "Portable safety behavior required by every profile.",
        )
        self.assertEqual(
            first.rule_contracts["rule.core.root-cause-only"].coverage,
            ("coverage.universal-safety",),
        )
        self.assertEqual(
            first.skill_dispatch_by_module["core"][0].skill_name,
            "coding-guidelines",
        )
        references = first.references_by_artifact["guide.agent-instructions"]
        self.assertEqual(references[0].target_managed_id, "guide.agent-instructions")
        self.assertEqual(references[1].repository_path, Path("DESIGN.md"))
        source = first.external_sources_by_setup["fixture"][0]
        self.assertEqual(source.skill_name, "coding-guidelines")
        self.assertEqual(source.source.revision, "0123456789abcdef0123456789abcdef01234567")

    def test_versioned_contract_mutations_fail_with_stable_diagnostics(self):
        cases = [
            ("missing rule carrier", self._remove_required_rule_carrier, "profile.rule.carrier.missing"),
            ("missing required profile rule", self._remove_required_profile_rule, "profile.rule.required.mismatch"),
            ("missing rule owner module", self._remove_rule_owner_module, "profile.rule.module.missing"),
            ("missing rule render binding", self._remove_rule_render_binding, "profile.rule.binding.missing"),
            ("missing dispatch mapping", self._remove_dispatch_mapping, "module.skillDispatch.mismatch"),
            ("extra dispatch mapping", self._add_extra_dispatch_mapping, "module.skillDispatch.mismatch"),
            ("unknown rule coverage", self._unknown_rule_coverage, "rule.coverage.unknown"),
            ("missing rule coverage", self._missing_rule_coverage, "rule.coverage.invalid"),
            ("missing coverage description", self._missing_coverage_description, "coverage.description.invalid"),
            ("unknown managed reference", self._unknown_managed_reference, "reference.managed.unknown"),
            ("missing reference ownership", self._missing_reference_ownership, "reference.ownership.unknown"),
            ("absolute repository path", self._absolute_repository_reference, "reference.repository.path.invalid"),
            ("mutable external ref", self._mutable_external_ref, "setup.skill.source.ref.mutable"),
            ("mutable setup ref", self._mutable_setup_ref, "setup.source: fixture.ref.mutable"),
            ("unsafe source path", self._unsafe_external_source_path, "setup.skill.source.path.invalid"),
            ("machine-local setup source", self._machine_local_setup_source, "setup.source: fixture.fields.invalid"),
            ("malformed tree digest", self._malformed_tree_digest, "setup.skill.treeDigest.invalid"),
            ("record provenance changes digest", self._change_external_repository, "setup.digest.mismatch"),
            ("missing rule guidance", self._missing_rule_guidance, "rule.guidance.invalid"),
            ("invalid external source shape", self._invalid_external_source_shape, "setup.skill.source.invalid"),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                diagnostics = self._load_invalid_v2_fixture(mutator)
                self.assertIn(expected_code, diagnostics)

    def test_versioned_contract_duplicate_identifiers_are_rejected(self):
        cases = [
            ("rule", self._duplicate_v2_rule, "rule.id.duplicate"),
            ("coverage", self._duplicate_coverage, "coverage.id.duplicate"),
            ("dispatch", self._duplicate_dispatch, "skill.dispatch.id.duplicate"),
            ("reference", self._duplicate_reference, "reference.id.duplicate"),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                diagnostics = self._load_invalid_v2_fixture(mutator)
                self.assertIn(expected_code, diagnostics)

    def test_contract_loading_has_no_write_command_or_network_side_effects(self):
        before = self._asset_bytes(V2_FIXTURE_ROOT)

        with (
            mock.patch.object(Path, "write_text", side_effect=AssertionError("write attempted")),
            mock.patch("subprocess.run", side_effect=AssertionError("command attempted")),
            mock.patch("urllib.request.urlopen", side_effect=AssertionError("network attempted")),
        ):
            load_asset_catalog(V2_FIXTURE_ROOT)

        self.assertEqual(self._asset_bytes(V2_FIXTURE_ROOT), before)

    def test_canonical_and_embedded_catalogs_load_successfully(self):
        canonical_root, embedded_root = setup_skill_roots(SKILL_ROOT)
        canonical = load_asset_catalog(canonical_root)
        embedded = load_asset_catalog(embedded_root)

        self.assertEqual(canonical.ordered_modules_by_profile, embedded.ordered_modules_by_profile)

    def test_every_bundled_external_skill_has_portable_immutable_provenance(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        for setup_id, setup in catalog.setups.items():
            self.assertEqual(setup["schemaVersion"], "setup-context-driven/setup-snapshot-v2")
            self.assertEqual(setup["version"], 2)
            external_names = {
                skill["name"]
                for skill in setup["skills"]
                if skill["source"]["type"] == "github"
            }
            self.assertEqual(
                external_names,
                {contract.skill_name for contract in catalog.external_sources_by_setup[setup_id]},
            )

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
                    "spec-workflow",
                    "external-triage",
                    "secondbrain",
                    "repository-extension",
                ],
                "rust-cli": [
                    "core",
                    "context-workflow",
                    "rust",
                    "cli-surface",
                    "autonomous-work",
                    "spec-workflow",
                    "external-triage",
                    "secondbrain",
                    "repository-extension",
                ],
                "standard-typescript-monorepo": [
                    "core",
                    "context-workflow",
                    "typescript",
                    "bun",
                    "monorepo",
                    "backend",
                    "frontend",
                    "autonomous-work",
                    "spec-workflow",
                    "external-triage",
                    "secondbrain",
                    "repository-extension",
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
                    "spec-workflow",
                    "external-triage",
                    "secondbrain",
                    "repository-extension",
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
                if asset.get("schemaVersion") == "setup-context-driven/profile/0.0.1":
                    self.assertEqual(asset["version"], "0.0.1")
                    continue
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
            self.assertIn("{{reference.", content, template["id"])

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
            ("malformed module document", self._malformed_module_document, "module.document.invalid"),
            ("malformed template item", self._malformed_template_item, "template.item.invalid"),
            ("malformed template collection", self._malformed_template_collection, "template.collection.invalid"),
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

    def _load_invalid_v2_fixture(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(V2_FIXTURE_ROOT, temp_root)
            mutator(temp_root)

            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)

        return "\n".join(captured.exception.diagnostics)

    def _asset_bytes(self, skill_root):
        return {
            path.relative_to(skill_root).as_posix(): path.read_bytes()
            for path in sorted((skill_root / "assets").rglob("*"))
            if path.is_file()
        }

    def _v2_module(self, temp_root):
        return temp_root / "assets" / "modules" / "core.json"

    def _v2_setup(self, temp_root):
        return temp_root / "assets" / "setups" / "fixture.json"

    def _remove_required_rule_carrier(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["supportingGuides"][0]["rules"] = []
        write_json(self._v2_module(temp_root), module)

    def _remove_dispatch_mapping(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["skillDispatch"] = []
        write_json(self._v2_module(temp_root), module)

    def _remove_required_profile_rule(self, temp_root):
        profile_path = temp_root / "assets" / "profiles" / "fixture.json"
        profile = read_json_copy(profile_path)
        profile["requiredRules"] = []
        write_json(profile_path, profile)

    def _remove_rule_owner_module(self, temp_root):
        profile_path = temp_root / "assets" / "profiles" / "fixture.json"
        profile = read_json_copy(profile_path)
        profile["modules"] = []
        write_json(profile_path, profile)

    def _remove_rule_render_binding(self, temp_root):
        template_path = (
            temp_root
            / "assets"
            / "templates"
            / "guides"
            / "agent-instructions.md"
        )
        template_path.write_text(
            template_path.read_text(encoding="utf-8").replace(
                "{{artifact.rules}}", "Portable rules"
            ),
            encoding="utf-8",
        )

    def _add_extra_dispatch_mapping(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["skillDispatch"].append(
            {"id": "extra-skill", "when": "An undeclared trigger runs."}
        )
        write_json(self._v2_module(temp_root), module)

    def _unknown_rule_coverage(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["rules"][0]["coverage"] = ["coverage.missing"]
        write_json(self._v2_module(temp_root), module)

    def _missing_rule_coverage(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        del module["rules"][0]["coverage"]
        write_json(self._v2_module(temp_root), module)

    def _missing_coverage_description(self, temp_root):
        path = temp_root / "assets" / "coverage.json"
        coverage = read_json_copy(path)
        del coverage["coverage"][0]["description"]
        write_json(path, coverage)

    def _unknown_managed_reference(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["supportingGuides"][0]["references"][0]["managedId"] = "guide.missing"
        write_json(self._v2_module(temp_root), module)

    def _missing_reference_ownership(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        del module["supportingGuides"][0]["references"][0]["ownership"]
        write_json(self._v2_module(temp_root), module)

    def _absolute_repository_reference(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["supportingGuides"][0]["references"][1]["path"] = "/tmp/DESIGN.md"
        write_json(self._v2_module(temp_root), module)

    def _mutable_external_ref(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["skills"][0]["source"]["ref"] = "main"
        write_json(self._v2_setup(temp_root), setup)

    def _mutable_setup_ref(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["source"]["ref"] = "main"
        write_json(self._v2_setup(temp_root), setup)

    def _machine_local_setup_source(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["source"]["localPath"] = "/tmp/skills"
        write_json(self._v2_setup(temp_root), setup)

    def _unsafe_external_source_path(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["skills"][0]["source"]["path"] = "../coding-guidelines"
        write_json(self._v2_setup(temp_root), setup)

    def _malformed_tree_digest(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["skills"][0]["treeDigest"] = "not-a-digest"
        write_json(self._v2_setup(temp_root), setup)

    def _change_external_repository(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["skills"][0]["source"]["repository"] = "example/other-skills"
        write_json(self._v2_setup(temp_root), setup)

    def _missing_rule_guidance(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        del module["rules"][0]["guidance"]
        write_json(self._v2_module(temp_root), module)

    def _invalid_external_source_shape(self, temp_root):
        setup = read_json_copy(self._v2_setup(temp_root))
        setup["skills"][0]["source"] = []
        write_json(self._v2_setup(temp_root), setup)

    def _duplicate_v2_rule(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["rules"].append(module["rules"][0])
        write_json(self._v2_module(temp_root), module)

    def _duplicate_coverage(self, temp_root):
        path = temp_root / "assets" / "coverage.json"
        coverage = read_json_copy(path)
        coverage["coverage"].append(coverage["coverage"][0])
        write_json(path, coverage)

    def _duplicate_dispatch(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["skillDispatch"].append(module["skillDispatch"][0])
        write_json(self._v2_module(temp_root), module)

    def _duplicate_reference(self, temp_root):
        module = read_json_copy(self._v2_module(temp_root))
        module["supportingGuides"][0]["references"].append(
            module["supportingGuides"][0]["references"][0]
        )
        write_json(self._v2_module(temp_root), module)

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

    def _malformed_module_document(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "rust.json"
        write_json(module_path, ["not", "an", "object"])

    def _malformed_template_item(self, temp_root):
        templates_path = temp_root / "assets" / "templates" / "index.json"
        templates = read_json_copy(templates_path)
        templates["templates"].append(None)
        write_json(templates_path, templates)

    def _malformed_template_collection(self, temp_root):
        templates_path = temp_root / "assets" / "templates" / "index.json"
        templates = read_json_copy(templates_path)
        templates["templates"] = {"id": "template.not-a-list"}
        write_json(templates_path, templates)


if __name__ == "__main__":
    unittest.main()
