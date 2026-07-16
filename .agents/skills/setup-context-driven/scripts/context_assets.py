"""Load and validate setup-context-driven portable assets."""

from __future__ import annotations

import copy
import hashlib
import json
import re
from dataclasses import dataclass
from pathlib import Path


ASSET_SCHEMA_VERSION = "setup-context-driven/assets-v1"
DECISIONS_SCHEMA_VERSION = "setup-context-driven/decisions-v1"
MODULE_SCHEMA_VERSION = "setup-context-driven/module-v1"
PROFILE_SCHEMA_VERSION = "setup-context-driven/profile-v1"
SETUP_SCHEMA_VERSION = "setup-context-driven/setup-snapshot-v1"
TEMPLATES_SCHEMA_VERSION = "setup-context-driven/templates-v1"
TEMPLATE_TOKEN = re.compile(r"\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}")
EFFECT_FIELDS = {
    "when",
    "activateModules",
    "requireDecisions",
    "includeArtifacts",
    "excludeArtifacts",
    "selectTemplates",
    "renderBindings",
}
CONDITION_OPERATORS = {"equals", "present"}


class AssetValidationError(Exception):
    """Raised when bundled setup-context-driven assets are invalid."""

    def __init__(self, diagnostics: list[str]):
        self.diagnostics = diagnostics
        super().__init__("\n".join(diagnostics))


@dataclass(frozen=True)
class DecisionCondition:
    decision_id: str
    operator: str
    value: object


@dataclass(frozen=True)
class TemplateSelection:
    artifact_id: str
    template_id: str


@dataclass(frozen=True)
class RenderBinding:
    artifact_id: str
    template_id: str
    token: str


@dataclass(frozen=True)
class DecisionEffect:
    decision_id: str
    condition: DecisionCondition
    activate_modules: tuple[str, ...]
    require_decisions: tuple[str, ...]
    include_artifacts: tuple[str, ...]
    exclude_artifacts: tuple[str, ...]
    template_selections: tuple[TemplateSelection, ...]
    render_bindings: tuple[RenderBinding, ...]


@dataclass(frozen=True)
class AssetCatalog:
    decisions: dict[str, dict]
    modules: dict[str, dict]
    profiles: dict[str, dict]
    setups: dict[str, dict]
    templates: dict[str, dict]
    ordered_modules_by_profile: dict[str, list[str]]
    profile_entry_decisions: dict[str, tuple[str, ...]]
    decision_effects: dict[str, tuple[DecisionEffect, ...]]


@dataclass(frozen=True)
class _ArtifactTarget:
    module_id: str
    kind: str


def load_asset_catalog(skill_root: str | Path) -> AssetCatalog:
    """Load and validate the asset catalog under a setup-context-driven skill."""

    root = Path(skill_root)
    assets_root = root / "assets"
    diagnostics: list[str] = []

    contract = _read_json(assets_root / "contract-v1.json", diagnostics)
    if contract and contract.get("schemaVersion") != ASSET_SCHEMA_VERSION:
        diagnostics.append(
            "contract.schemaVersion: expected setup-context-driven/assets-v1"
        )

    decisions_doc = _read_json(assets_root / "decisions.json", diagnostics)
    templates_doc = _read_json(assets_root / "templates" / "index.json", diagnostics)
    modules = _read_collection(
        assets_root / "modules", MODULE_SCHEMA_VERSION, "module", diagnostics
    )
    profiles = _read_collection(
        assets_root / "profiles", PROFILE_SCHEMA_VERSION, "profile", diagnostics
    )
    setups = _read_collection(
        assets_root / "setups", SETUP_SCHEMA_VERSION, "setup", diagnostics
    )

    decisions = _index_assets(
        decisions_doc,
        "decisions",
        DECISIONS_SCHEMA_VERSION,
        "decision",
        diagnostics,
    )
    templates = _index_assets(
        templates_doc,
        "templates",
        TEMPLATES_SCHEMA_VERSION,
        "template",
        diagnostics,
    )

    _validate_versions(decisions, "decision", diagnostics)
    _validate_versions(templates, "template", diagnostics)
    _validate_versions(modules, "module", diagnostics)
    _validate_versions(profiles, "profile", diagnostics)
    _validate_versions(setups, "setup", diagnostics)
    template_tokens = _validate_templates(assets_root / "templates", templates, diagnostics)
    _validate_modules(modules, decisions, templates, diagnostics)
    _validate_setups(setups, diagnostics)
    profile_entry_decisions = _validate_profile_entry_decisions(
        profiles, decisions, diagnostics
    )
    decision_effects = _validate_decision_effects(
        decisions,
        modules,
        templates,
        template_tokens,
        diagnostics,
    )

    ordered_modules_by_profile: dict[str, list[str]] = {}
    for profile_id, profile in sorted(profiles.items()):
        ordered_modules_by_profile[profile_id] = _resolve_profile_modules(
            profile_id, profile, modules, setups, diagnostics
        )

    for profile_id, ordered_modules in ordered_modules_by_profile.items():
        profile = profiles[profile_id]
        setup = setups.get(profile.get("setup"))
        if setup is None:
            continue
        _validate_profile_skills(profile_id, ordered_modules, modules, setup, diagnostics)

    if diagnostics:
        raise AssetValidationError(diagnostics)

    return AssetCatalog(
        decisions=decisions,
        modules=modules,
        profiles=profiles,
        setups=setups,
        templates=templates,
        ordered_modules_by_profile=ordered_modules_by_profile,
        profile_entry_decisions=profile_entry_decisions,
        decision_effects=decision_effects,
    )


def clone_assets_to(source_root: str | Path, destination_root: str | Path) -> None:
    """Copy asset JSON into a temporary skill root for mutation-based tests."""

    source_assets = Path(source_root) / "assets"
    destination_assets = Path(destination_root) / "assets"
    for source_path in source_assets.rglob("*"):
        if source_path.is_dir():
            continue
        relative_path = source_path.relative_to(source_assets)
        destination_path = destination_assets / relative_path
        destination_path.parent.mkdir(parents=True, exist_ok=True)
        destination_path.write_text(source_path.read_text(encoding="utf-8"), encoding="utf-8")


def read_json_copy(path: str | Path) -> dict:
    return copy.deepcopy(json.loads(Path(path).read_text(encoding="utf-8")))


def write_json(path: str | Path, data: dict) -> None:
    Path(path).write_text(
        json.dumps(data, indent=2, sort_keys=False) + "\n",
        encoding="utf-8",
    )


def _read_json(path: Path, diagnostics: list[str]) -> dict:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError:
        diagnostics.append(f"asset.file.missing: {path}")
    except json.JSONDecodeError as error:
        diagnostics.append(f"asset.json.invalid: {path}: {error.msg}")
    return {}


def _read_collection(
    directory: Path,
    schema_version: str,
    kind: str,
    diagnostics: list[str],
) -> dict[str, dict]:
    collection: dict[str, dict] = {}
    for path in sorted(directory.glob("*.json")):
        data = _read_json(path, diagnostics)
        if not data:
            continue
        if data.get("schemaVersion") != schema_version:
            diagnostics.append(
                f"{kind}.schemaVersion: {path.name}: expected {schema_version}"
            )
        asset_id = data.get("id")
        if not isinstance(asset_id, str) or not asset_id:
            diagnostics.append(f"{kind}.id.missing: {path.name}")
            continue
        expected_name = f"{asset_id}.json"
        if path.name != expected_name:
            diagnostics.append(
                f"{kind}.filename.mismatch: {path.name}: expected {expected_name}"
            )
        if asset_id in collection:
            diagnostics.append(f"{kind}.id.duplicate: {asset_id}")
        collection[asset_id] = data
    return collection


def _index_assets(
    data: dict,
    key: str,
    schema_version: str,
    kind: str,
    diagnostics: list[str],
) -> dict[str, dict]:
    if not data:
        return {}
    if data.get("schemaVersion") != schema_version:
        diagnostics.append(f"{kind}.schemaVersion: expected {schema_version}")
    indexed: dict[str, dict] = {}
    for item in data.get(key, []):
        asset_id = item.get("id")
        if not isinstance(asset_id, str) or not asset_id:
            diagnostics.append(f"{kind}.id.missing")
            continue
        if asset_id in indexed:
            diagnostics.append(f"{kind}.id.duplicate: {asset_id}")
        indexed[asset_id] = item
    return indexed


def _validate_versions(
    assets: dict[str, dict], kind: str, diagnostics: list[str]
) -> None:
    for asset_id, data in sorted(assets.items()):
        version = data.get("version")
        if not isinstance(version, int) or version < 1:
            diagnostics.append(f"{kind}.version.invalid: {asset_id}")


def _validate_templates(
    templates_root: Path,
    templates: dict[str, dict],
    diagnostics: list[str],
) -> dict[str, set[str]]:
    template_tokens: dict[str, set[str]] = {}
    for template_id, template in sorted(templates.items()):
        declared_tokens = _declared_template_tokens(template_id, template, diagnostics)
        template_tokens[template_id] = declared_tokens
        path = template.get("path")
        if not isinstance(path, str) or _is_unsafe_relative_path(path):
            diagnostics.append(f"template.path.invalid: {template_id}")
            continue
        template_path = templates_root / path
        if not template_path.is_file():
            diagnostics.append(f"template.file.missing: {template_id}: {path}")
            continue
        content = template_path.read_text(encoding="utf-8")
        for token in sorted(set(TEMPLATE_TOKEN.findall(content))):
            if token not in declared_tokens:
                diagnostics.append(f"template.token.undeclared: {template_id} -> {token}")
    return template_tokens


def _declared_template_tokens(
    template_id: str,
    template: dict,
    diagnostics: list[str],
) -> set[str]:
    raw_tokens = template.get("tokens", [])
    if not isinstance(raw_tokens, list) or not all(
        isinstance(token, str) and token for token in raw_tokens
    ):
        diagnostics.append(f"template.tokens.invalid: {template_id}")
        return set()

    tokens: set[str] = set()
    for token in raw_tokens:
        if token in tokens:
            diagnostics.append(f"template.token.duplicate: {template_id} -> {token}")
        tokens.add(token)
    return tokens


def _validate_modules(
    modules: dict[str, dict],
    decisions: dict[str, dict],
    templates: dict[str, dict],
    diagnostics: list[str],
) -> None:
    seen_rules: dict[str, str] = {}
    seen_blocks: dict[str, str] = {}
    seen_guides: dict[str, str] = {}

    for module_id, module in sorted(modules.items()):
        for dependency_id in module.get("dependsOn", []):
            if dependency_id not in modules:
                diagnostics.append(
                    f"module.dependency.unknown: {module_id} -> {dependency_id}"
                )
        for conflict_id in module.get("conflictsWith", []):
            if conflict_id not in modules:
                diagnostics.append(
                    f"module.conflict.unknown: {module_id} -> {conflict_id}"
                )
        for decision_id in module.get("requiredDecisions", []):
            if decision_id not in decisions:
                diagnostics.append(
                    f"module.decision.unknown: {module_id} -> {decision_id}"
                )

        for rule in module.get("rules", []):
            rule_id = rule.get("id")
            if not isinstance(rule_id, str) or not rule_id:
                diagnostics.append(f"rule.id.missing: {module_id}")
                continue
            if not isinstance(rule.get("version"), int) or rule["version"] < 1:
                diagnostics.append(f"rule.version.invalid: {rule_id}")
            owner = seen_rules.get(rule_id)
            if owner is not None:
                diagnostics.append(
                    f"rule.id.duplicate: {rule_id}: {owner}, {module_id}"
                )
            seen_rules[rule_id] = module_id

        for block in module.get("rootBlocks", []):
            block_id = block.get("id")
            if not isinstance(block_id, str) or not block_id:
                diagnostics.append(f"managed.block.id.missing: {module_id}")
                continue
            if not isinstance(block.get("version"), int) or block["version"] < 1:
                diagnostics.append(f"managed.block.version.invalid: {block_id}")
            owner = seen_blocks.get(block_id)
            if owner is not None:
                diagnostics.append(
                    f"managed.block.id.duplicate: {block_id}: {owner}, {module_id}"
                )
            seen_blocks[block_id] = module_id
            _validate_template_reference(
                block.get("template"), "root-block", templates, diagnostics, block_id
            )
            _validate_rule_references(
                block.get("rules", []), module, diagnostics, block_id
            )

        for guide in module.get("supportingGuides", []):
            guide_id = guide.get("id")
            if not isinstance(guide_id, str) or not guide_id:
                diagnostics.append(f"guide.id.missing: {module_id}")
                continue
            if not isinstance(guide.get("version"), int) or guide["version"] < 1:
                diagnostics.append(f"guide.version.invalid: {guide_id}")
            owner = seen_guides.get(guide_id)
            if owner is not None:
                diagnostics.append(f"guide.id.duplicate: {guide_id}: {owner}, {module_id}")
            seen_guides[guide_id] = module_id
            target_path = guide.get("path")
            if not isinstance(target_path, str) or _is_unsafe_relative_path(target_path):
                diagnostics.append(f"guide.path.invalid: {guide_id}")
            _validate_template_reference(
                guide.get("template"), "guide", templates, diagnostics, guide_id
            )
            _validate_rule_references(
                guide.get("rules", []), module, diagnostics, guide_id
            )


def _validate_profile_entry_decisions(
    profiles: dict[str, dict],
    decisions: dict[str, dict],
    diagnostics: list[str],
) -> dict[str, tuple[str, ...]]:
    entry_decisions: dict[str, tuple[str, ...]] = {}
    for profile_id, profile in sorted(profiles.items()):
        raw_entries = profile.get("entryDecisions")
        if not isinstance(raw_entries, list) or not all(
            isinstance(decision_id, str) and decision_id for decision_id in raw_entries
        ):
            diagnostics.append(f"profile.entryDecisions.invalid: {profile_id}")
            entry_decisions[profile_id] = ()
            continue

        seen: set[str] = set()
        ordered_entries: list[str] = []
        for decision_id in raw_entries:
            if decision_id in seen:
                diagnostics.append(
                    f"profile.entryDecision.duplicate: {profile_id} -> {decision_id}"
                )
                continue
            seen.add(decision_id)
            if decision_id not in decisions:
                diagnostics.append(
                    f"profile.entryDecision.unknown: {profile_id} -> {decision_id}"
                )
            ordered_entries.append(decision_id)
        entry_decisions[profile_id] = tuple(ordered_entries)
    return entry_decisions


def _validate_decision_effects(
    decisions: dict[str, dict],
    modules: dict[str, dict],
    templates: dict[str, dict],
    template_tokens: dict[str, set[str]],
    diagnostics: list[str],
) -> dict[str, tuple[DecisionEffect, ...]]:
    artifacts = _index_artifact_targets(modules)
    effects_by_decision: dict[str, tuple[DecisionEffect, ...]] = {}
    dependency_graph: dict[str, list[str]] = {decision_id: [] for decision_id in decisions}
    binding_owners: dict[tuple[str, str], str] = {}

    for decision_id, decision in decisions.items():
        raw_effects = decision.get("effects")
        if not isinstance(raw_effects, list) or not raw_effects:
            diagnostics.append(f"decision.effects.invalid: {decision_id}")
            effects_by_decision[decision_id] = ()
            continue

        validated_effects: list[DecisionEffect] = []
        for index, raw_effect in enumerate(raw_effects):
            if not isinstance(raw_effect, dict):
                diagnostics.append(f"decision.effect.invalid: {decision_id}[{index}]")
                continue
            for field in sorted(set(raw_effect) - EFFECT_FIELDS):
                diagnostics.append(
                    f"decision.effect.field.unknown: {decision_id}[{index}] -> {field}"
                )

            condition = _validate_condition(
                decision_id,
                decision,
                raw_effect.get("when"),
                diagnostics,
            )
            activate_modules = _validate_module_targets(
                decision_id,
                raw_effect.get("activateModules", []),
                modules,
                diagnostics,
            )
            require_decisions = _validate_decision_targets(
                decision_id,
                raw_effect.get("requireDecisions", []),
                decisions,
                diagnostics,
            )
            include_artifacts = _validate_artifact_targets(
                decision_id,
                raw_effect.get("includeArtifacts", []),
                artifacts,
                modules,
                diagnostics,
            )
            exclude_artifacts = _validate_artifact_targets(
                decision_id,
                raw_effect.get("excludeArtifacts", []),
                artifacts,
                modules,
                diagnostics,
            )
            template_selections = _validate_template_selections(
                decision_id,
                raw_effect.get("selectTemplates", []),
                artifacts,
                modules,
                templates,
                diagnostics,
            )
            render_bindings = _validate_render_bindings(
                decision_id,
                raw_effect.get("renderBindings", []),
                artifacts,
                modules,
                templates,
                template_tokens,
                binding_owners,
                diagnostics,
            )
            dependency_graph[decision_id].extend(require_decisions)

            if condition is None:
                continue
            validated_effects.append(
                DecisionEffect(
                    decision_id=decision_id,
                    condition=condition,
                    activate_modules=activate_modules,
                    require_decisions=require_decisions,
                    include_artifacts=include_artifacts,
                    exclude_artifacts=exclude_artifacts,
                    template_selections=template_selections,
                    render_bindings=render_bindings,
                )
            )
        effects_by_decision[decision_id] = tuple(validated_effects)

    _validate_decision_dependency_cycles(dependency_graph, diagnostics)
    return effects_by_decision


def _validate_condition(
    decision_id: str,
    decision: dict,
    condition: object,
    diagnostics: list[str],
) -> DecisionCondition | None:
    if not isinstance(condition, dict) or len(condition) != 1:
        diagnostics.append(f"decision.condition.invalid: {decision_id}")
        return None

    operator, value = next(iter(condition.items()))
    if operator not in CONDITION_OPERATORS:
        diagnostics.append(f"decision.condition.operator.unknown: {decision_id} -> {operator}")
        return None
    if operator == "present":
        if not isinstance(value, bool):
            diagnostics.append(f"decision.condition.type.invalid: {decision_id}: present")
            return None
        return DecisionCondition(decision_id=decision_id, operator=operator, value=value)

    if not _condition_value_matches_decision(decision, value):
        diagnostics.append(f"decision.condition.type.invalid: {decision_id}: equals")
        return None
    return DecisionCondition(decision_id=decision_id, operator=operator, value=value)


def _condition_value_matches_decision(decision: dict, value: object) -> bool:
    decision_type = decision.get("type")
    if decision_type == "boolean":
        return isinstance(value, bool)
    if decision_type == "string":
        return isinstance(value, str)
    if decision_type == "enum":
        return isinstance(value, str) and value in decision.get("values", [])
    return False


def _validate_module_targets(
    decision_id: str,
    raw_targets: object,
    modules: dict[str, dict],
    diagnostics: list[str],
) -> tuple[str, ...]:
    module_ids = _validate_string_list(
        raw_targets, f"decision.effect.modules.invalid: {decision_id}", diagnostics
    )
    for module_id in module_ids:
        module = modules.get(module_id)
        if module is None:
            diagnostics.append(f"decision.effect.module.unknown: {decision_id} -> {module_id}")
            continue
        if decision_id not in module.get("requiredDecisions", []):
            diagnostics.append(f"decision.effect.module.unowned: {decision_id} -> {module_id}")
    return tuple(module_ids)


def _validate_decision_targets(
    decision_id: str,
    raw_targets: object,
    decisions: dict[str, dict],
    diagnostics: list[str],
) -> tuple[str, ...]:
    target_ids = _validate_string_list(
        raw_targets, f"decision.effect.decisions.invalid: {decision_id}", diagnostics
    )
    for target_id in target_ids:
        if target_id not in decisions:
            diagnostics.append(f"decision.effect.decision.unknown: {decision_id} -> {target_id}")
    return tuple(target_ids)


def _validate_artifact_targets(
    decision_id: str,
    raw_targets: object,
    artifacts: dict[str, _ArtifactTarget],
    modules: dict[str, dict],
    diagnostics: list[str],
) -> tuple[str, ...]:
    artifact_ids = _validate_string_list(
        raw_targets, f"decision.effect.artifacts.invalid: {decision_id}", diagnostics
    )
    for artifact_id in artifact_ids:
        artifact = artifacts.get(artifact_id)
        if artifact is None:
            diagnostics.append(f"decision.effect.artifact.unknown: {decision_id} -> {artifact_id}")
            continue
        _validate_artifact_owner(decision_id, artifact_id, artifact, modules, diagnostics)
    return tuple(artifact_ids)


def _validate_template_selections(
    decision_id: str,
    raw_selections: object,
    artifacts: dict[str, _ArtifactTarget],
    modules: dict[str, dict],
    templates: dict[str, dict],
    diagnostics: list[str],
) -> tuple[TemplateSelection, ...]:
    if not isinstance(raw_selections, list):
        diagnostics.append(f"decision.effect.templateSelections.invalid: {decision_id}")
        return ()

    selections: list[TemplateSelection] = []
    for index, raw_selection in enumerate(raw_selections):
        if not isinstance(raw_selection, dict):
            diagnostics.append(
                f"decision.effect.templateSelection.invalid: {decision_id}[{index}]"
            )
            continue
        artifact_id = raw_selection.get("artifact")
        template_id = raw_selection.get("template")
        if not isinstance(artifact_id, str) or not isinstance(template_id, str):
            diagnostics.append(
                f"decision.effect.templateSelection.invalid: {decision_id}[{index}]"
            )
            continue
        artifact = artifacts.get(artifact_id)
        if artifact is None:
            diagnostics.append(f"decision.effect.artifact.unknown: {decision_id} -> {artifact_id}")
            continue
        _validate_artifact_owner(decision_id, artifact_id, artifact, modules, diagnostics)
        if _validate_template_target(decision_id, artifact, template_id, templates, diagnostics):
            selections.append(
                TemplateSelection(artifact_id=artifact_id, template_id=template_id)
            )
    return tuple(selections)


def _validate_render_bindings(
    decision_id: str,
    raw_bindings: object,
    artifacts: dict[str, _ArtifactTarget],
    modules: dict[str, dict],
    templates: dict[str, dict],
    template_tokens: dict[str, set[str]],
    binding_owners: dict[tuple[str, str], str],
    diagnostics: list[str],
) -> tuple[RenderBinding, ...]:
    if not isinstance(raw_bindings, list):
        diagnostics.append(f"decision.effect.renderBindings.invalid: {decision_id}")
        return ()

    bindings: list[RenderBinding] = []
    for index, raw_binding in enumerate(raw_bindings):
        if not isinstance(raw_binding, dict):
            diagnostics.append(f"decision.effect.renderBinding.invalid: {decision_id}[{index}]")
            continue
        artifact_id = raw_binding.get("artifact")
        template_id = raw_binding.get("template")
        token = raw_binding.get("token")
        if not all(isinstance(value, str) and value for value in [artifact_id, template_id, token]):
            diagnostics.append(f"decision.effect.renderBinding.invalid: {decision_id}[{index}]")
            continue
        artifact = artifacts.get(artifact_id)
        if artifact is None:
            diagnostics.append(f"decision.effect.artifact.unknown: {decision_id} -> {artifact_id}")
            continue
        _validate_artifact_owner(decision_id, artifact_id, artifact, modules, diagnostics)
        if not _validate_template_target(decision_id, artifact, template_id, templates, diagnostics):
            continue
        if token not in template_tokens.get(template_id, set()):
            diagnostics.append(
                f"decision.effect.binding.token.unknown: {decision_id} -> {template_id}: {token}"
            )
            continue
        binding_key = (template_id, token)
        owner = binding_owners.get(binding_key)
        if owner is not None:
            diagnostics.append(
                f"decision.effect.binding.duplicate: {template_id}:{token}: {owner}, {decision_id}"
            )
            continue
        binding_owners[binding_key] = decision_id
        bindings.append(
            RenderBinding(
                artifact_id=artifact_id,
                template_id=template_id,
                token=token,
            )
        )
    return tuple(bindings)


def _validate_template_target(
    decision_id: str,
    artifact: _ArtifactTarget,
    template_id: str,
    templates: dict[str, dict],
    diagnostics: list[str],
) -> bool:
    template = templates.get(template_id)
    if template is None:
        diagnostics.append(f"decision.effect.template.unknown: {decision_id} -> {template_id}")
        return False
    if template.get("kind") != artifact.kind:
        diagnostics.append(
            f"decision.effect.template.kind: {decision_id} -> {template_id}: expected {artifact.kind}"
        )
        return False
    return True


def _validate_artifact_owner(
    decision_id: str,
    artifact_id: str,
    artifact: _ArtifactTarget,
    modules: dict[str, dict],
    diagnostics: list[str],
) -> None:
    module = modules.get(artifact.module_id, {})
    if decision_id not in module.get("requiredDecisions", []):
        diagnostics.append(f"decision.effect.artifact.unowned: {decision_id} -> {artifact_id}")


def _validate_string_list(
    raw_items: object,
    diagnostic: str,
    diagnostics: list[str],
) -> list[str]:
    if raw_items is None:
        return []
    if not isinstance(raw_items, list) or not all(
        isinstance(item, str) and item for item in raw_items
    ):
        diagnostics.append(diagnostic)
        return []
    return list(raw_items)


def _index_artifact_targets(modules: dict[str, dict]) -> dict[str, _ArtifactTarget]:
    artifacts: dict[str, _ArtifactTarget] = {}
    for module_id, module in sorted(modules.items()):
        for block in module.get("rootBlocks", []):
            artifact_id = block.get("id")
            if isinstance(artifact_id, str) and artifact_id:
                artifacts[artifact_id] = _ArtifactTarget(
                    module_id=module_id,
                    kind="root-block",
                )
        for guide in module.get("supportingGuides", []):
            artifact_id = guide.get("id")
            if isinstance(artifact_id, str) and artifact_id:
                artifacts[artifact_id] = _ArtifactTarget(
                    module_id=module_id,
                    kind="guide",
                )
    return artifacts


def _validate_decision_dependency_cycles(
    dependency_graph: dict[str, list[str]],
    diagnostics: list[str],
) -> None:
    visiting: list[str] = []
    visited: set[str] = set()

    def visit(decision_id: str) -> None:
        if decision_id in visited:
            return
        if decision_id in visiting:
            cycle = visiting[visiting.index(decision_id) :] + [decision_id]
            diagnostics.append(f"decision.dependency.cycle: {' -> '.join(cycle)}")
            return
        visiting.append(decision_id)
        for dependent_id in dependency_graph.get(decision_id, []):
            if dependent_id in dependency_graph:
                visit(dependent_id)
        visiting.pop()
        visited.add(decision_id)

    for decision_id in dependency_graph:
        visit(decision_id)


def _validate_template_reference(
    template_id: object,
    expected_kind: str,
    templates: dict[str, dict],
    diagnostics: list[str],
    owner_id: str,
) -> None:
    if not isinstance(template_id, str) or not template_id:
        diagnostics.append(f"template.reference.missing: {owner_id}")
        return
    template = templates.get(template_id)
    if template is None:
        diagnostics.append(f"template.reference.unknown: {owner_id} -> {template_id}")
        return
    if template.get("kind") != expected_kind:
        diagnostics.append(
            f"template.kind.mismatch: {owner_id} -> {template_id}: "
            f"expected {expected_kind}"
        )


def _validate_rule_references(
    rule_ids: list[str],
    module: dict,
    diagnostics: list[str],
    owner_id: str,
) -> None:
    module_rules = {rule.get("id") for rule in module.get("rules", [])}
    for rule_id in rule_ids:
        if rule_id not in module_rules:
            diagnostics.append(f"rule.reference.unknown: {owner_id} -> {rule_id}")


def _validate_setups(setups: dict[str, dict], diagnostics: list[str]) -> None:
    for setup_id, setup in sorted(setups.items()):
        seen_names: set[str] = set()
        normalized_paths: list[str] = []
        for skill in setup.get("skills", []):
            name = skill.get("name")
            path = skill.get("path")
            if not isinstance(name, str) or not name:
                diagnostics.append(f"setup.skill.name.missing: {setup_id}")
                continue
            if name in seen_names:
                diagnostics.append(f"setup.skill.name.duplicate: {setup_id}: {name}")
            seen_names.add(name)
            if not isinstance(path, str) or _is_unsafe_relative_path(path):
                diagnostics.append(f"setup.skill.path.invalid: {setup_id}: {name}")
                continue
            normalized_paths.append(path)
            if not isinstance(skill.get("contentDigest"), str) or not skill[
                "contentDigest"
            ]:
                diagnostics.append(f"setup.skill.digest.missing: {setup_id}: {name}")

        expected_digest = _paths_digest(normalized_paths)
        if setup.get("digest") != expected_digest:
            diagnostics.append(
                f"setup.digest.mismatch: {setup_id}: expected {expected_digest}"
            )


def _resolve_profile_modules(
    profile_id: str,
    profile: dict,
    modules: dict[str, dict],
    setups: dict[str, dict],
    diagnostics: list[str],
) -> list[str]:
    selected = profile.get("modules", [])
    if not isinstance(selected, list) or not all(isinstance(item, str) for item in selected):
        diagnostics.append(f"profile.modules.invalid: {profile_id}")
        return []

    if profile.get("setup") not in setups:
        diagnostics.append(f"profile.setup.unknown: {profile_id} -> {profile.get('setup')}")

    selected_set = set(selected)
    for module_id in selected:
        if module_id not in modules:
            diagnostics.append(f"profile.module.unknown: {profile_id} -> {module_id}")

    for module_id in selected:
        module = modules.get(module_id)
        if module is None:
            continue
        for dependency_id in module.get("dependsOn", []):
            if dependency_id not in selected_set:
                diagnostics.append(
                    f"profile.dependency.missing: {profile_id}: "
                    f"{module_id} requires {dependency_id}"
                )

    visiting: list[str] = []
    visited: set[str] = set()
    ordered: list[str] = []

    def visit(module_id: str) -> None:
        if module_id in visited or module_id not in modules:
            return
        if module_id in visiting:
            cycle = visiting[visiting.index(module_id) :] + [module_id]
            diagnostics.append(
                f"module.dependency.cycle: {profile_id}: {' -> '.join(cycle)}"
            )
            return
        visiting.append(module_id)
        for dependency_id in modules[module_id].get("dependsOn", []):
            if dependency_id in selected_set:
                visit(dependency_id)
        visiting.pop()
        visited.add(module_id)
        ordered.append(module_id)

    for module_id in selected:
        visit(module_id)

    order_index = {module_id: index for index, module_id in enumerate(ordered)}
    for module_id in selected:
        module = modules.get(module_id)
        if module is None:
            continue
        for conflict_id in module.get("conflictsWith", []):
            if conflict_id in selected_set:
                diagnostics.append(
                    f"module.conflict: {profile_id}: {module_id} conflicts with {conflict_id}"
                )
        for dependency_id in module.get("dependsOn", []):
            if dependency_id in selected_set and order_index.get(dependency_id, -1) > order_index.get(module_id, -1):
                diagnostics.append(
                    f"profile.order.invalid: {profile_id}: {dependency_id} after {module_id}"
                )

    return ordered


def _validate_profile_skills(
    profile_id: str,
    ordered_modules: list[str],
    modules: dict[str, dict],
    setup: dict,
    diagnostics: list[str],
) -> None:
    setup_skills = {skill.get("name") for skill in setup.get("skills", [])}
    for module_id in ordered_modules:
        for skill_name in modules[module_id].get("requiredSkills", []):
            if skill_name not in setup_skills:
                diagnostics.append(
                    f"skills.reference.outside-setup: {profile_id}: "
                    f"{module_id} requires {skill_name}"
                )


def _is_unsafe_relative_path(path: str) -> bool:
    candidate = Path(path)
    return candidate.is_absolute() or ".." in candidate.parts or "\\" in path


def _paths_digest(paths: list[str]) -> str:
    payload = "\n".join(paths) + "\n"
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()
