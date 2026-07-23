# Suite: normalized skill dispatch contracts
# Invariant: every selected setup skill has one catalog owner and one rendered entry.
# Boundary IN: asset validation, profile normalization, and dispatch rendering.
# Boundary OUT: repository apply, installed-skill discovery, and network access.

import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
FIXTURE_ROOT = Path(__file__).resolve().parent / "fixtures" / "upgrade-contracts"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    read_json_copy,
    setup_snapshot_digest,
    write_json,
)
from context_setup import render_skill_dispatch  # noqa: E402


class SkillDispatchContractTests(unittest.TestCase):
    def test_distinct_triggers_render_under_one_stably_ordered_skill_entry(self):
        catalog = load_asset_catalog(FIXTURE_ROOT)
        modules = catalog.ordered_modules_by_profile["fixture"]

        first = render_skill_dispatch(catalog, modules)
        second = render_skill_dispatch(catalog, modules)

        self.assertEqual(first, second)
        self.assertEqual(
            first.splitlines(),
            [
                "- `coding-guidelines`:",
                "  - `trigger.implementation.change`: Writing or changing implementation code.",
                "  - `trigger.implementation.review`: Reviewing implementation code.",
            ],
        )

    def test_dependent_module_cannot_redeclare_a_skill_owner(self):
        diagnostics = self._load_invalid_fixture(self._add_dependent_duplicate_owner)

        self.assertIn("skill.dispatch.owner.duplicate", diagnostics)

    def test_reused_trigger_id_fails_even_when_wording_differs(self):
        diagnostics = self._load_invalid_fixture(self._reuse_trigger_id_with_new_wording)

        self.assertIn("skill.dispatch.trigger.id.duplicate", diagnostics)

    def test_duplicate_skill_contract_in_one_module_is_rejected(self):
        diagnostics = self._load_invalid_fixture(self._duplicate_skill_contract)

        self.assertIn("skill.dispatch.skill.duplicate", diagnostics)

    def test_required_skill_missing_from_setup_snapshot_is_rejected(self):
        diagnostics = self._load_invalid_fixture(self._add_required_skill_outside_setup)

        self.assertIn("skills.reference.outside-setup", diagnostics)
        self.assertIn("profile.skill.set.mismatch", diagnostics)

    def test_required_skill_missing_from_dispatch_map_is_rejected(self):
        diagnostics = self._load_invalid_fixture(self._remove_dispatch_contract)

        self.assertIn("module.skillDispatch.mismatch", diagnostics)
        self.assertIn("profile.skill.set.mismatch", diagnostics)

    def test_installed_skill_missing_from_required_and_dispatch_sets_is_rejected(self):
        diagnostics = self._load_invalid_fixture(self._add_unowned_setup_skill)

        self.assertIn("profile.skill.set.mismatch", diagnostics)

    def test_every_bundled_profile_has_exact_installed_required_and_dispatch_sets(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        for profile_id, modules in catalog.ordered_modules_by_profile.items():
            with self.subTest(profile=profile_id):
                setup = catalog.setups[catalog.profiles[profile_id]["setup"]]
                installed = {skill["name"] for skill in setup["skills"]}
                required = {
                    skill
                    for module_id in modules
                    for skill in catalog.modules[module_id]["requiredSkills"]
                }
                dispatched = {
                    skill
                    for skill, triggers in catalog.skill_dispatch_by_skill.items()
                    if any(trigger.owner_module in modules for trigger in triggers)
                }

                self.assertEqual(installed, required)
                self.assertEqual(required, dispatched)

    def test_framework_dispatch_only_renders_for_profiles_that_install_it(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        go_dispatch = render_skill_dispatch(
            catalog, catalog.ordered_modules_by_profile["go-cli-tui"]
        )
        rust_dispatch = render_skill_dispatch(
            catalog, catalog.ordered_modules_by_profile["rust-cli"]
        )
        typescript_dispatch = render_skill_dispatch(
            catalog,
            catalog.ordered_modules_by_profile["typescript-bun-monorepo"],
        )

        self.assertIn("- `react`:", typescript_dispatch)
        self.assertNotIn("- `react`:", go_dispatch)
        self.assertNotIn("- `react`:", rust_dispatch)
        self.assertIn("- `golang-cli`:", go_dispatch)
        self.assertNotIn("- `golang-cli`:", typescript_dispatch)
        self.assertNotIn("- `golang-cli`:", rust_dispatch)

    def test_typescript_profile_renders_exact_activation_bundles_deterministically(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        modules = catalog.ordered_modules_by_profile["typescript-bun-monorepo"]

        first = render_skill_dispatch(catalog, modules)
        second = render_skill_dispatch(catalog, modules)

        self.assertEqual(first, second)
        self.assertEqual(
            {
                bundle_id: bundle.skills
                for bundle_id, bundle in catalog.activation_bundles.items()
            },
            {
                "bundle.debugging": (
                    "systematic-debugging",
                    "diagnosing-bugs",
                    "no-workarounds",
                ),
                "bundle.delivery": (
                    "conventional-commits",
                    "github-pr-workflow",
                ),
                "bundle.frontend-react": (
                    "react",
                    "react-best-practices",
                    "react-composition-patterns",
                ),
                "bundle.frontend-ui-quality": (
                    "frontend-design",
                    "interaction-design",
                    "interface-design",
                    "fixing-accessibility",
                    "wcag-audit-patterns",
                    "web-design-guidelines",
                ),
                "bundle.hono-endpoint": (
                    "hono-api-best-practices",
                    "hono",
                    "zod",
                ),
                "bundle.hono-endpoint-persistence": (
                    "hono-api-best-practices",
                    "hono",
                    "zod",
                    "drizzle-orm",
                ),
                "bundle.production-code": (
                    "coding-guidelines",
                    "clean-code",
                    "solid",
                ),
                "bundle.qa": ("qa-gate", "evidence-gate"),
                "bundle.security": (
                    "security-best-practices",
                    "security-threat-model",
                ),
                "bundle.testing": ("testing-boss", "tdd", "vitest"),
            },
        )
        self.assertIn(
            "- `trigger.production-code`: Writing, modifying, or reviewing production code.\n"
            "  - `bundle.production-code`: `coding-guidelines`, `clean-code`, `solid`",
            first,
        )
        self.assertIn(
            "- `trigger.hono.endpoint`: Creating, modifying, or reviewing a Hono endpoint without persistence behavior.\n"
            "  - `bundle.hono-endpoint`: `hono-api-best-practices`, `hono`, `zod`",
            first,
        )
        self.assertIn(
            "- `trigger.hono.endpoint-persistence`: Creating, modifying, or reviewing a Hono endpoint with persistence behavior.\n"
            "  - `bundle.hono-endpoint-persistence`: `hono-api-best-practices`, `hono`, `zod`, `drizzle-orm`",
            first,
        )

    def test_specialist_activation_bundles_remain_distinct_and_snapshot_backed(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        setup = catalog.setups["typescript-bun"]

        self.assertEqual(
            setup["activationBundles"],
            [
                {
                    "id": bundle.bundle_id,
                    "skills": list(bundle.skills),
                }
                for bundle in catalog.activation_bundles_by_setup["typescript-bun"]
            ],
        )
        specialist_triggers = {
            "trigger.frontend.react-feature": "bundle.frontend-react",
            "trigger.frontend.ui-quality": "bundle.frontend-ui-quality",
            "trigger.testing": "bundle.testing",
            "trigger.security": "bundle.security",
            "trigger.delivery": "bundle.delivery",
        }
        self.assertEqual(
            {
                activation.trigger_id: activation.bundle_id
                for activation in catalog.skill_activations
                if activation.trigger_id in specialist_triggers
            },
            specialist_triggers,
        )
        self.assertEqual(
            len(set(specialist_triggers.values())),
            len(specialist_triggers),
        )

    def test_invalid_activation_documents_fail_before_repository_writes(self):
        cases = [
            (
                "missing snapshot member",
                self._remove_snapshot_bundle_member,
                "setup.activationBundle.members.mismatch",
            ),
            (
                "duplicate bundle member",
                self._duplicate_activation_bundle_member,
                "skill.activation.bundle.member.duplicate",
            ),
            (
                "duplicate trigger id",
                self._duplicate_activation_trigger,
                "skill.activation.trigger.id.duplicate",
            ),
            (
                "duplicate bundle owner",
                self._duplicate_activation_bundle_owner,
                "skill.activation.owner.duplicate",
            ),
            (
                "unknown skill",
                self._add_unknown_activation_skill,
                "skill.activation.skill.unknown",
            ),
            (
                "unstable bundle order",
                self._reverse_activation_bundle_order,
                "skill.activation.bundle.order.invalid",
            ),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    temp_root = Path(temp_dir) / "setup-context-driven"
                    clone_assets_to(SKILL_ROOT, temp_root)
                    mutator(temp_root)
                    with mock.patch.object(
                        Path,
                        "write_text",
                        side_effect=AssertionError("repository write attempted"),
                    ):
                        with self.assertRaises(AssetValidationError) as captured:
                            load_asset_catalog(temp_root)
                self.assertIn(expected_code, "\n".join(captured.exception.diagnostics))

    def _load_invalid_fixture(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(FIXTURE_ROOT, temp_root)
            mutator(temp_root)
            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)
        return "\n".join(captured.exception.diagnostics)

    def _core_module(self, temp_root):
        return temp_root / "assets" / "modules" / "core.json"

    def _add_dependent_duplicate_owner(self, temp_root):
        module_path = temp_root / "assets" / "modules" / "dependent.json"
        write_json(
            module_path,
            {
                "schemaVersion": "setup-context-driven/module-v3",
                "id": "dependent",
                "version": 1,
                "dependsOn": ["core"],
                "conflictsWith": [],
                "requiredSkills": ["coding-guidelines"],
                "skillDispatch": [
                    {
                        "skill": "coding-guidelines",
                        "triggers": [
                            {
                                "id": "trigger.dependent.code-change",
                                "when": "Changing code through a dependent workflow.",
                            }
                        ],
                    }
                ],
                "requiredDecisions": [],
                "rootBlocks": [],
                "supportingGuides": [],
                "rules": [],
                "repositoryExtensions": [],
            },
        )
        profile_path = temp_root / "assets" / "profiles" / "fixture.json"
        profile = read_json_copy(profile_path)
        profile["modules"].append("dependent")
        write_json(profile_path, profile)

    def _reuse_trigger_id_with_new_wording(self, temp_root):
        module = read_json_copy(self._core_module(temp_root))
        module["skillDispatch"][0]["triggers"].append(
            {
                "id": "trigger.implementation.change",
                "when": "A differently worded trigger with the same stable intent.",
            }
        )
        write_json(self._core_module(temp_root), module)

    def _duplicate_skill_contract(self, temp_root):
        module = read_json_copy(self._core_module(temp_root))
        module["skillDispatch"].append(
            {
                "skill": "coding-guidelines",
                "triggers": [
                    {
                        "id": "trigger.implementation.duplicate-contract",
                        "when": "Changing code through a duplicate contract.",
                    }
                ],
            }
        )
        write_json(self._core_module(temp_root), module)

    def _add_required_skill_outside_setup(self, temp_root):
        module = read_json_copy(self._core_module(temp_root))
        module["requiredSkills"].append("missing-skill")
        module["skillDispatch"].append(
            {
                "skill": "missing-skill",
                "triggers": [
                    {
                        "id": "trigger.missing-skill.use",
                        "when": "Using a skill absent from the selected setup.",
                    }
                ],
            }
        )
        write_json(self._core_module(temp_root), module)

    def _remove_dispatch_contract(self, temp_root):
        module = read_json_copy(self._core_module(temp_root))
        module["skillDispatch"] = []
        write_json(self._core_module(temp_root), module)

    def _add_unowned_setup_skill(self, temp_root):
        setup_path = temp_root / "assets" / "setups" / "fixture.json"
        setup = read_json_copy(setup_path)
        setup["skills"].append(
            {
                "name": "unowned-skill",
                "path": "skills/unowned-skill",
                "source": {"type": "repo", "name": "unowned-skill"},
                "contentDigest": "d" * 64,
            }
        )
        setup["digest"] = setup_snapshot_digest(
            setup["skills"], setup.get("activationBundles")
        )
        write_json(setup_path, setup)

    def _activation_document(self, temp_root):
        return temp_root / "assets" / "skill-activations.json"

    def _remove_snapshot_bundle_member(self, temp_root):
        setup_path = temp_root / "assets" / "setups" / "typescript-bun.json"
        setup = read_json_copy(setup_path)
        production = next(
            bundle
            for bundle in setup["activationBundles"]
            if bundle["id"] == "bundle.production-code"
        )
        production["skills"].remove("solid")
        write_json(setup_path, setup)

    def _duplicate_activation_bundle_member(self, temp_root):
        path = self._activation_document(temp_root)
        document = read_json_copy(path)
        document["bundles"][0]["skills"].append(
            document["bundles"][0]["skills"][0]
        )
        write_json(path, document)

    def _duplicate_activation_trigger(self, temp_root):
        path = self._activation_document(temp_root)
        document = read_json_copy(path)
        document["activations"].append(dict(document["activations"][0]))
        write_json(path, document)

    def _duplicate_activation_bundle_owner(self, temp_root):
        path = self._activation_document(temp_root)
        document = read_json_copy(path)
        document["activations"].append(
            {
                "id": "trigger.z-duplicate-delivery-owner",
                "owner": "context-workflow",
                "when": "Using the delivery bundle through another owner.",
                "bundle": "bundle.delivery",
            }
        )
        write_json(path, document)

    def _add_unknown_activation_skill(self, temp_root):
        path = self._activation_document(temp_root)
        document = read_json_copy(path)
        document["bundles"][0]["skills"].append("missing-skill")
        write_json(path, document)

    def _reverse_activation_bundle_order(self, temp_root):
        path = self._activation_document(temp_root)
        document = read_json_copy(path)
        document["bundles"].reverse()
        write_json(path, document)


if __name__ == "__main__":
    unittest.main()
