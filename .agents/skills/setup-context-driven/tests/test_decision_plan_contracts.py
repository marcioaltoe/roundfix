import sys
import tempfile
import unittest
from dataclasses import FrozenInstanceError, replace
from itertools import product
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
import context_setup  # noqa: E402


EXPECTED_DECISIONS = (
    "spec.scaffold",
    "domain.layout",
    "triage.external",
    "autonomous.enabled",
    "verification.gate",
    "runtime.backend",
    "runtime.design",
    "language.generated",
    "secondbrain.enabled",
)
EXPECTED_ENTRY_DECISIONS = (
    "language.generated",
    "verification.gate",
    "spec.scaffold",
    "domain.layout",
    "triage.external",
    "autonomous.enabled",
    "secondbrain.enabled",
)


class DecisionPlanContractTests(unittest.TestCase):
    def test_canonical_catalog_declares_entry_decisions_and_effects(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        self.assertEqual(tuple(catalog.decisions), EXPECTED_DECISIONS)
        self.assertEqual(set(catalog.decision_effects), set(EXPECTED_DECISIONS))
        for decision_id in EXPECTED_DECISIONS:
            self.assertGreater(len(catalog.decision_effects[decision_id]), 0)

        self.assertEqual(
            catalog.profile_entry_decisions,
            {
                "go-cli-tui": EXPECTED_ENTRY_DECISIONS,
                "rust-cli": EXPECTED_ENTRY_DECISIONS,
                "typescript-bun-monorepo": EXPECTED_ENTRY_DECISIONS,
            },
        )

        autonomous_effect = self._effect_for(
            catalog, "autonomous.enabled", "equals", True
        )
        self.assertEqual(autonomous_effect.activate_modules, ("autonomous-work",))
        self.assertEqual(
            autonomous_effect.require_decisions,
            ("runtime.backend", "runtime.design"),
        )
        self.assertEqual(
            autonomous_effect.include_artifacts,
            ("root.autonomous-work", "guide.autonomous-work"),
        )

        render_tokens = {
            decision_id: catalog.decision_effects[decision_id][0].render_bindings[0].token
            for decision_id in [
                "runtime.backend",
                "runtime.design",
                "verification.gate",
            ]
        }
        self.assertEqual(
            render_tokens,
            {
                "runtime.backend": "runtime.backend",
                "runtime.design": "runtime.design",
                "verification.gate": "verification.gate",
            },
        )

    def test_validated_effect_models_are_immutable(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        effect = catalog.decision_effects["autonomous.enabled"][0]

        with self.assertRaises(FrozenInstanceError):
            effect.decision_id = "changed"

        with self.assertRaises(FrozenInstanceError):
            effect.condition.operator = "changed"

    def test_invalid_effect_contracts_fail_with_stable_diagnostics(self):
        cases = [
            ("unsupported decision type", self._unsupported_decision_type, "decision.type.invalid"),
            ("enum values must be strings", self._non_string_enum_value, "decision.values.invalid"),
            ("unknown module target", self._unknown_module_target, "decision.effect.module.unknown"),
            ("unknown artifact target", self._unknown_artifact_target, "decision.effect.artifact.unknown"),
            ("unknown template target", self._unknown_template_target, "decision.effect.template.unknown"),
            ("unknown dependent decision", self._unknown_dependent_decision, "decision.effect.decision.unknown"),
            ("incompatible condition", self._incompatible_condition, "decision.condition.type.invalid"),
            ("duplicate binding", self._duplicate_binding, "decision.effect.binding.duplicate"),
            ("undeclared token", self._undeclared_template_token, "template.token.undeclared"),
            ("dependency cycle", self._decision_dependency_cycle, "decision.dependency.cycle"),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                diagnostics = self._load_invalid_fixture(mutator)
                self.assertIn(expected_code, diagnostics)

    def test_loading_same_assets_twice_is_deterministic(self):
        first = load_asset_catalog(SKILL_ROOT)
        second = load_asset_catalog(SKILL_ROOT)

        self.assertEqual(self._catalog_snapshot(first), self._catalog_snapshot(second))

    def test_canonical_and_embedded_portable_assets_load_together(self):
        repo_root = self._repo_root(SKILL_ROOT)

        canonical = load_asset_catalog(repo_root / ".agents" / "skills" / "setup-context-driven")
        embedded = load_asset_catalog(repo_root / "skills" / "setup-context-driven")

        self.assertEqual(self._catalog_snapshot(canonical), self._catalog_snapshot(embedded))

    def test_every_finite_artifact_branch_resolves_declared_references(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "DESIGN.md").write_text("# Repository design\n", encoding="utf-8")
            for profile_id in catalog.profiles:
                finite_decisions = [
                    decision_id
                    for decision_id in catalog.profile_entry_decisions[profile_id]
                    if catalog.decisions[decision_id]["type"] in {"boolean", "enum"}
                ]
                finite_values = [
                    self._finite_values(catalog.decisions[decision_id])
                    for decision_id in finite_decisions
                ]
                for values in product(*finite_values):
                    decisions = dict(zip(finite_decisions, values, strict=True))
                    decisions["verification.gate"] = "make verify"
                    if decisions["autonomous.enabled"]:
                        decisions["runtime.backend"] = "codex gpt-5.5 xhigh"
                        decisions["runtime.design"] = "claude opus xhigh"
                    plan = context_setup.resolve_decision_plan(
                        catalog,
                        profile_id,
                        {"decisions": self._manifest_decisions(decisions)},
                        {},
                    )

                    self.assertEqual(plan.unresolved_decisions, (), (profile_id, decisions))
                    self.assertEqual(
                        context_setup.validate_decision_plan_references(repo, catalog, plan),
                        [],
                        (profile_id, decisions),
                    )

    def test_single_context_monorepo_selects_its_referenced_guide(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        decisions = self._complete_decisions(**{"domain.layout": "single-context"})
        plan = context_setup.resolve_decision_plan(
            catalog,
            "typescript-bun-monorepo",
            {"decisions": self._manifest_decisions(decisions)},
            {},
        )

        definite_ids = {
            artifact.artifact.managed_id
            for artifact in plan.artifacts
            if artifact.present and artifact.state == "definite"
        }
        self.assertIn("guide.monorepo", definite_ids)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "DESIGN.md").write_text("# Repository design\n", encoding="utf-8")
            self.assertEqual(
                context_setup.validate_decision_plan_references(repo, catalog, plan),
                [],
            )

    def test_definite_source_rejects_excluded_managed_reference_target(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_root)
            self._exclude_domain_reference_target(temp_root)
            catalog = load_asset_catalog(temp_root)
            plan = context_setup.resolve_decision_plan(
                catalog,
                "rust-cli",
                {"decisions": self._manifest_decisions(self._complete_decisions())},
                {},
            )
            stale_target = Path(temp_dir) / "docs" / "agents" / "domain.md"
            stale_target.parent.mkdir(parents=True)
            stale_target.write_text("# Stale domain guide\n", encoding="utf-8")
            plan = replace(
                plan,
                artifacts=tuple(
                    replace(
                        planned_artifact,
                        artifact=replace(
                            planned_artifact.artifact,
                            content=(
                                planned_artifact.artifact.content
                                + "\n[Domain guide](docs/agents/domain.md)\n"
                            ),
                        ),
                    )
                    if planned_artifact.artifact.managed_id == "root.context-workflow"
                    else planned_artifact
                    for planned_artifact in plan.artifacts
                ),
            )

            findings = context_setup.validate_decision_plan_references(
                Path(temp_dir), catalog, plan
            )

        managed = [
            finding for finding in findings if finding.code == "reference.managed.missing"
        ]
        markdown = [
            finding for finding in findings if finding.code == "docs.reference.broken"
        ]
        self.assertEqual(len(managed), 1, findings)
        self.assertEqual(managed[0].path, "AGENTS.md")
        self.assertEqual(managed[0].managed_id, "root.context-workflow")
        self.assertIn("guide.domain", managed[0].message)
        self.assertEqual(len(markdown), 1, findings)
        self.assertEqual(markdown[0].path, "AGENTS.md")

    def test_repository_reference_paths_reject_absolute_and_escaping_values(self):
        for name, path_value in (
            ("absolute", "/tmp/DESIGN.md"),
            ("escaping", "../DESIGN.md"),
        ):
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    temp_root = Path(temp_dir) / "setup-context-driven"
                    clone_assets_to(SKILL_ROOT, temp_root)
                    module_path = temp_root / "assets" / "modules" / "frontend.json"
                    module = read_json_copy(module_path)
                    module["supportingGuides"][0]["references"][0]["path"] = path_value
                    write_json(module_path, module)

                    with self.assertRaises(AssetValidationError) as captured:
                        load_asset_catalog(temp_root)

                self.assertIn(
                    "reference.repository.path.invalid",
                    "\n".join(captured.exception.diagnostics),
                )

    def test_each_declared_reference_token_has_one_typed_path_binding(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        managed_paths = context_setup.managed_artifact_paths(catalog)

        for artifact_id, references in catalog.references_by_artifact.items():
            tokens = [reference.token for reference in references]
            self.assertEqual(len(tokens), len(set(tokens)), artifact_id)
            for reference in references:
                if reference.ownership == "setup":
                    self.assertIn(reference.target_managed_id, managed_paths)
                    self.assertIsNone(reference.repository_path)
                else:
                    self.assertIsNotNone(reference.repository_path)
                    self.assertIsNone(reference.target_managed_id)

        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_root)
            module_path = temp_root / "assets" / "modules" / "frontend.json"
            module = read_json_copy(module_path)
            module["supportingGuides"][0]["references"].append(
                {
                    "id": "reference.frontend.duplicate-token",
                    "token": "repository.design-contract",
                    "ownership": "setup",
                    "managedId": "guide.frontend",
                }
            )
            write_json(module_path, module)

            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)

        self.assertIn(
            "reference.token.duplicate",
            "\n".join(captured.exception.diagnostics),
        )

    def _effect_for(self, catalog, decision_id, operator, value):
        for effect in catalog.decision_effects[decision_id]:
            if effect.condition.operator == operator and effect.condition.value == value:
                return effect
        self.fail(f"missing effect for {decision_id} when {operator}={value!r}")

    def _finite_values(self, decision):
        if decision["type"] == "boolean":
            return (False, True)
        return tuple(decision["values"])

    def _complete_decisions(self, **overrides):
        decisions = {
            "spec.scaffold": True,
            "domain.layout": "single-context",
            "triage.external": False,
            "autonomous.enabled": True,
            "runtime.backend": "codex gpt-5.5 xhigh",
            "runtime.design": "claude opus xhigh",
            "verification.gate": "make verify",
            "language.generated": "English",
            "secondbrain.enabled": False,
        }
        decisions.update(overrides)
        return decisions

    def _manifest_decisions(self, decisions):
        return {
            decision_id: {"value": value, "confirmedAt": "2026-07-21"}
            for decision_id, value in decisions.items()
        }

    def _exclude_domain_reference_target(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        effect = self._decision(decisions, "domain.layout")["effects"][0]
        effect["includeArtifacts"] = []
        effect["excludeArtifacts"] = ["guide.domain"]
        effect["selectTemplates"] = []
        write_json(decisions_path, decisions)

    def _load_invalid_fixture(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_root)
            mutator(temp_root)

            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)

        return "\n".join(captured.exception.diagnostics)

    def _unknown_module_target(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "autonomous.enabled")["effects"][0][
            "activateModules"
        ].append("module.missing")
        write_json(decisions_path, decisions)

    def _unsupported_decision_type(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "runtime.backend")["type"] = "bool"
        write_json(decisions_path, decisions)

    def _non_string_enum_value(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "domain.layout")["values"].append(True)
        write_json(decisions_path, decisions)

    def _unknown_artifact_target(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "spec.scaffold")["effects"][0][
            "includeArtifacts"
        ].append("root.missing")
        write_json(decisions_path, decisions)

    def _unknown_template_target(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "domain.layout")["effects"][0]["selectTemplates"][0][
            "template"
        ] = "template.missing"
        write_json(decisions_path, decisions)

    def _unknown_dependent_decision(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "autonomous.enabled")["effects"][0][
            "requireDecisions"
        ].append("decision.missing")
        write_json(decisions_path, decisions)

    def _incompatible_condition(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "spec.scaffold")["effects"][0]["when"] = {
            "equals": "true"
        }
        write_json(decisions_path, decisions)

    def _duplicate_binding(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "runtime.design")["effects"][0]["renderBindings"][0][
            "token"
        ] = "runtime.backend"
        write_json(decisions_path, decisions)

    def _undeclared_template_token(self, temp_root):
        template_path = temp_root / "assets" / "templates" / "guides" / "autonomous-work.md"
        template_path.write_text(
            template_path.read_text(encoding="utf-8") + "\n{{runtime.missing}}\n",
            encoding="utf-8",
        )

    def _decision_dependency_cycle(self, temp_root):
        decisions_path = temp_root / "assets" / "decisions.json"
        decisions = read_json_copy(decisions_path)
        self._decision(decisions, "spec.scaffold")["effects"][0].setdefault(
            "requireDecisions", []
        ).append("secondbrain.enabled")
        self._decision(decisions, "secondbrain.enabled")["effects"][0].setdefault(
            "requireDecisions", []
        ).append("spec.scaffold")
        write_json(decisions_path, decisions)

    def _decision(self, decisions, decision_id):
        for decision in decisions["decisions"]:
            if decision["id"] == decision_id:
                return decision
        self.fail(f"missing decision fixture {decision_id}")

    def _catalog_snapshot(self, catalog):
        return (
            tuple(catalog.profiles),
            tuple((key, tuple(value)) for key, value in catalog.ordered_modules_by_profile.items()),
            tuple((key, tuple(value)) for key, value in catalog.profile_entry_decisions.items()),
            tuple(catalog.decisions),
            tuple(
                (decision_id, tuple(self._effect_snapshot(effect) for effect in effects))
                for decision_id, effects in catalog.decision_effects.items()
            ),
        )

    def _effect_snapshot(self, effect):
        return (
            effect.decision_id,
            (effect.condition.decision_id, effect.condition.operator, effect.condition.value),
            effect.activate_modules,
            effect.require_decisions,
            effect.include_artifacts,
            effect.exclude_artifacts,
            tuple(
                (selection.artifact_id, selection.template_id)
                for selection in effect.template_selections
            ),
            tuple(
                (binding.artifact_id, binding.template_id, binding.token)
                for binding in effect.render_bindings
            ),
        )

    def _repo_root(self, skill_root):
        for parent in [skill_root, *skill_root.parents]:
            if (
                parent / ".agents" / "skills" / "setup-context-driven"
            ).is_dir() and (parent / "skills" / "setup-context-driven").is_dir():
                return parent
        self.fail("could not locate repository root")


if __name__ == "__main__":
    unittest.main()
