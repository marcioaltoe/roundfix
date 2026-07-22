"""Load and validate setup-context-driven portable assets."""

from __future__ import annotations

import copy
import hashlib
import json
import os
import re
import stat
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable


ASSET_SCHEMA_VERSION = "setup-context-driven/assets-v1"
COVERAGE_SCHEMA_VERSION = "setup-context-driven/coverage-v1"
COVERAGE_SCHEMA_VERSION_V2 = "setup-context-driven/coverage-v2"
DECISIONS_SCHEMA_VERSION = "setup-context-driven/decisions-v1"
MODULE_SCHEMA_VERSION = "setup-context-driven/module-v1"
MODULE_SCHEMA_VERSION_V2 = "setup-context-driven/module-v2"
MODULE_SCHEMA_VERSION_V3 = "setup-context-driven/module-v3"
PROFILE_SCHEMA_VERSION = "setup-context-driven/profile-v1"
PROFILE_SCHEMA_VERSION_V2 = "setup-context-driven/profile-v2"
PROFILE_SCHEMA_VERSION_V3 = "setup-context-driven/profile-v3"
SETUP_SCHEMA_VERSION = "setup-context-driven/setup-snapshot-v1"
SETUP_SCHEMA_VERSION_V2 = "setup-context-driven/setup-snapshot-v2"
TEMPLATES_SCHEMA_VERSION = "setup-context-driven/templates-v1"
UPGRADE_TRANSITION_SCHEMA_VERSION = "setup-context-driven/upgrade-transition-v1"
TEMPLATE_TOKEN = re.compile(r"\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}")
LOWERCASE_SHA256 = re.compile(r"^[0-9a-f]{64}$")
IMMUTABLE_GIT_REF = re.compile(r"^(?:[0-9a-f]{40}|[0-9a-f]{64})$")
GITHUB_REPOSITORY = re.compile(
    r"^[A-Za-z0-9](?:[A-Za-z0-9_.-]{0,38})/[A-Za-z0-9_.-]+$"
)
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
DECISION_TYPES = {"boolean", "enum", "string"}
CLAUSE_ENFORCEMENTS = {"mandatory", "prohibited", "stop-and-ask"}
TRANSITION_DISPOSITIONS = {"retained", "moved", "replaced", "rejected"}
FORMATTER_KINDS = {"none", "selected"}
MAX_DELEGATION_ALIASES = 16
MAX_DELEGATION_ALIAS_LENGTH = 80


class AssetValidationError(Exception):
    """Raised when bundled setup-context-driven assets are invalid."""

    def __init__(self, diagnostics: list[str]):
        self.diagnostics = diagnostics
        super().__init__("\n".join(diagnostics))


class PortableTreeError(ValueError):
    """Raised when a skill tree contains an entry that cannot be restored safely."""


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
class CoverageContract:
    coverage_id: str
    description: str
    delegation_aliases: tuple[str, ...] = ()


@dataclass(frozen=True)
class ClauseContract:
    clause_id: str
    enforcement: str
    guidance: str


@dataclass(frozen=True)
class RuleContract:
    rule_id: str
    coverage: tuple[str, ...]
    guidance: str = ""
    clauses: tuple[ClauseContract, ...] = ()


@dataclass(frozen=True)
class SkillDispatch:
    skill_name: str
    when: str
    trigger_id: str = ""
    owner_module: str = ""


@dataclass(frozen=True)
class RepositoryOwnedExtension:
    extension_id: str
    target_path: Path
    template_id: str
    root_pointer_id: str
    decision_id: str


@dataclass(frozen=True)
class FormatterContract:
    kind: str
    formatter_id: str | None = None
    version: str | None = None
    fixture_paths: tuple[Path, ...] = ()
    golden_digest: str | None = None


@dataclass(frozen=True)
class PriorClauseContract:
    clause_id: str
    enforcement: str
    carrier_id: str
    guidance_digest: str


@dataclass(frozen=True)
class TransitionMapping:
    from_clause: str
    disposition: str
    targets: tuple[str, ...]
    reason: str


@dataclass(frozen=True)
class UpgradeTransition:
    transition_id: str
    version: int
    from_baseline: str
    to_baseline: str
    prior_clauses: tuple[PriorClauseContract, ...]
    mappings: tuple[TransitionMapping, ...]


@dataclass(frozen=True)
class ArtifactReference:
    reference_id: str
    token: str
    ownership: str
    target_managed_id: str | None = None
    repository_path: Path | None = None


@dataclass(frozen=True)
class SkillSourceRef:
    provider: str
    repository: str
    revision: str
    source_path: Path


@dataclass(frozen=True)
class ExternalSkillContract:
    skill_name: str
    target_path: Path
    source: SkillSourceRef
    tree_digest: str


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
    coverage_contracts: dict[str, CoverageContract]
    rule_contracts: dict[str, RuleContract]
    skill_dispatch_by_module: dict[str, tuple[SkillDispatch, ...]]
    references_by_artifact: dict[str, tuple[ArtifactReference, ...]]
    external_sources_by_setup: dict[str, tuple[ExternalSkillContract, ...]]
    repository_extensions: dict[str, RepositoryOwnedExtension]
    formatter_by_profile: dict[str, FormatterContract]
    upgrade_transitions: dict[str, UpgradeTransition]
    skill_dispatch_by_skill: dict[str, tuple[SkillDispatch, ...]]


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
    coverage_path = assets_root / "coverage.json"
    coverage_doc = _read_json(coverage_path, diagnostics) if coverage_path.is_file() else {}
    modules = _read_collection(
        assets_root / "modules",
        (MODULE_SCHEMA_VERSION, MODULE_SCHEMA_VERSION_V2, MODULE_SCHEMA_VERSION_V3),
        "module",
        diagnostics,
    )
    profiles = _read_collection(
        assets_root / "profiles",
        (PROFILE_SCHEMA_VERSION, PROFILE_SCHEMA_VERSION_V2, PROFILE_SCHEMA_VERSION_V3),
        "profile",
        diagnostics,
    )
    setups = _read_collection(
        assets_root / "setups",
        (SETUP_SCHEMA_VERSION, SETUP_SCHEMA_VERSION_V2),
        "setup",
        diagnostics,
    )
    transitions = _read_collection(
        assets_root / "retention",
        UPGRADE_TRANSITION_SCHEMA_VERSION,
        "transition",
        diagnostics,
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
    coverage = _index_assets(
        coverage_doc,
        "coverage",
        (COVERAGE_SCHEMA_VERSION, COVERAGE_SCHEMA_VERSION_V2),
        "coverage",
        diagnostics,
    )

    _validate_versions(coverage, "coverage", diagnostics)
    _validate_versions(decisions, "decision", diagnostics)
    _validate_versions(templates, "template", diagnostics)
    _validate_versions(modules, "module", diagnostics)
    _validate_versions(profiles, "profile", diagnostics)
    _validate_versions(setups, "setup", diagnostics)
    _validate_versions(transitions, "transition", diagnostics)
    coverage_contracts = _validate_coverage_contracts(
        coverage,
        coverage_doc.get("schemaVersion") if isinstance(coverage_doc, dict) else None,
        diagnostics,
    )
    template_tokens, rendered_template_tokens = _validate_templates(
        assets_root / "templates", templates, diagnostics
    )
    _validate_modules(modules, decisions, templates, diagnostics)
    (
        rule_contracts,
        skill_dispatch_by_module,
        references_by_artifact,
        repository_extensions,
    ) = _validate_versioned_module_contracts(
        modules,
        decisions,
        coverage_contracts,
        templates,
        template_tokens,
        rendered_template_tokens,
        diagnostics,
    )
    external_sources_by_setup = _validate_setups(setups, diagnostics)
    profile_entry_decisions = _validate_profile_entry_decisions(
        profiles, decisions, diagnostics
    )
    _validate_decision_contracts(decisions, diagnostics)
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

    _validate_profile_rule_contracts(
        profiles,
        ordered_modules_by_profile,
        modules,
        rule_contracts,
        rendered_template_tokens,
        diagnostics,
    )
    formatter_by_profile = _validate_profile_formatters(profiles, diagnostics)
    skill_dispatch_by_skill = _normalize_profile_dispatch_contracts(
        profiles,
        ordered_modules_by_profile,
        skill_dispatch_by_module,
        diagnostics,
    )
    upgrade_transitions = _validate_upgrade_transitions(
        transitions,
        rule_contracts,
        repository_extensions,
        diagnostics,
    )

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
        coverage_contracts=coverage_contracts,
        rule_contracts=rule_contracts,
        skill_dispatch_by_module=skill_dispatch_by_module,
        references_by_artifact=references_by_artifact,
        external_sources_by_setup=external_sources_by_setup,
        repository_extensions=repository_extensions,
        formatter_by_profile=formatter_by_profile,
        upgrade_transitions=upgrade_transitions,
        skill_dispatch_by_skill=skill_dispatch_by_skill,
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
    schema_versions: str | tuple[str, ...],
    kind: str,
    diagnostics: list[str],
) -> dict[str, dict]:
    accepted_versions = (
        (schema_versions,) if isinstance(schema_versions, str) else schema_versions
    )
    collection: dict[str, dict] = {}
    for path in sorted(directory.glob("*.json")):
        data = _read_json(path, diagnostics)
        if not isinstance(data, dict):
            diagnostics.append(f"{kind}.document.invalid: {path.name}: expected object")
            continue
        if not data:
            continue
        if data.get("schemaVersion") not in accepted_versions:
            expected = ", ".join(accepted_versions)
            diagnostics.append(
                f"{kind}.schemaVersion: {path.name}: expected one of {expected}"
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
    data: object,
    key: str,
    schema_versions: str | tuple[str, ...],
    kind: str,
    diagnostics: list[str],
) -> dict[str, dict]:
    if not data:
        return {}
    if not isinstance(data, dict):
        diagnostics.append(f"{kind}.document.invalid: expected object")
        return {}
    accepted_versions = (
        (schema_versions,) if isinstance(schema_versions, str) else schema_versions
    )
    if data.get("schemaVersion") not in accepted_versions:
        diagnostics.append(
            f"{kind}.schemaVersion: expected one of {', '.join(accepted_versions)}"
        )
    indexed: dict[str, dict] = {}
    raw_items = data.get(key, [])
    if not isinstance(raw_items, list):
        diagnostics.append(f"{kind}.collection.invalid: expected list")
        return {}
    for index, item in enumerate(raw_items):
        if not isinstance(item, dict):
            diagnostics.append(f"{kind}.item.invalid: {index}")
            continue
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


def _validate_coverage_contracts(
    coverage: dict[str, dict],
    schema_version: object,
    diagnostics: list[str],
) -> dict[str, CoverageContract]:
    contracts: dict[str, CoverageContract] = {}
    for coverage_id, item in sorted(coverage.items()):
        description = item.get("description")
        if not isinstance(description, str) or not description.strip():
            diagnostics.append(f"coverage.description.invalid: {coverage_id}")
            continue
        aliases: tuple[str, ...] = ()
        if schema_version == COVERAGE_SCHEMA_VERSION_V2:
            raw_aliases = item.get("delegationAliases")
            if (
                not isinstance(raw_aliases, list)
                or len(raw_aliases) > MAX_DELEGATION_ALIASES
                or not all(_is_valid_delegation_alias(alias) for alias in raw_aliases)
            ):
                diagnostics.append(
                    f"coverage.delegationAliases.invalid: {coverage_id}"
                )
            else:
                normalized_aliases = [alias.casefold() for alias in raw_aliases]
                if len(set(normalized_aliases)) != len(normalized_aliases):
                    diagnostics.append(
                        f"coverage.delegationAliases.duplicate: {coverage_id}"
                    )
                aliases = tuple(sorted(normalized_aliases))
        contracts[coverage_id] = CoverageContract(
            coverage_id=coverage_id,
            description=description,
            delegation_aliases=aliases,
        )
    return contracts


def _is_valid_delegation_alias(alias: object) -> bool:
    return (
        isinstance(alias, str)
        and bool(alias)
        and alias == alias.strip()
        and alias == alias.casefold()
        and len(alias) <= MAX_DELEGATION_ALIAS_LENGTH
        and "\n" not in alias
        and "\r" not in alias
        and "\x00" not in alias
    )


def _validate_templates(
    templates_root: Path,
    templates: dict[str, dict],
    diagnostics: list[str],
) -> tuple[dict[str, set[str]], dict[str, set[str]]]:
    template_tokens: dict[str, set[str]] = {}
    rendered_tokens: dict[str, set[str]] = {}
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
        content_tokens = set(TEMPLATE_TOKEN.findall(content))
        rendered_tokens[template_id] = content_tokens
        for token in sorted(content_tokens):
            if token not in declared_tokens:
                diagnostics.append(f"template.token.undeclared: {template_id} -> {token}")
    return template_tokens, rendered_tokens


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
        if module.get("schemaVersion") in {
            MODULE_SCHEMA_VERSION_V2,
            MODULE_SCHEMA_VERSION_V3,
        }:
            continue
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


def _validate_versioned_module_contracts(
    modules: dict[str, dict],
    decisions: dict[str, dict],
    coverage: dict[str, CoverageContract],
    templates: dict[str, dict],
    template_tokens: dict[str, set[str]],
    rendered_template_tokens: dict[str, set[str]],
    diagnostics: list[str],
) -> tuple[
    dict[str, RuleContract],
    dict[str, tuple[SkillDispatch, ...]],
    dict[str, tuple[ArtifactReference, ...]],
    dict[str, RepositoryOwnedExtension],
]:
    versioned_modules = {
        module_id: module
        for module_id, module in modules.items()
        if module.get("schemaVersion")
        in {MODULE_SCHEMA_VERSION_V2, MODULE_SCHEMA_VERSION_V3}
    }
    if not versioned_modules:
        return {}, {}, {}, {}
    if not coverage:
        diagnostics.append("coverage.catalog.missing: versioned modules require coverage.json")

    rule_contracts: dict[str, RuleContract] = {}
    rule_owners: dict[str, str] = {}
    dispatch_by_module: dict[str, tuple[SkillDispatch, ...]] = {}
    references_by_artifact: dict[str, tuple[ArtifactReference, ...]] = {}
    repository_extensions: dict[str, RepositoryOwnedExtension] = {}
    artifacts: dict[str, tuple[str, str, dict]] = {}
    raw_extensions: list[tuple[str, dict]] = []
    reference_owners: dict[str, str] = {}
    clause_owners: dict[str, str] = {}
    trigger_owners: dict[str, str] = {}
    block_owners: dict[str, str] = {}
    guide_owners: dict[str, str] = {}

    for module_id, module in sorted(modules.items()):
        if module.get("schemaVersion") in {
            MODULE_SCHEMA_VERSION_V2,
            MODULE_SCHEMA_VERSION_V3,
        }:
            continue
        for rule in _dict_items(module.get("rules")):
            rule_id = rule.get("id")
            if isinstance(rule_id, str) and rule_id:
                rule_owners[rule_id] = module_id
        for block in _dict_items(module.get("rootBlocks")):
            block_id = block.get("id")
            if isinstance(block_id, str) and block_id:
                block_owners[block_id] = module_id
        for guide in _dict_items(module.get("supportingGuides")):
            guide_id = guide.get("id")
            if isinstance(guide_id, str) and guide_id:
                guide_owners[guide_id] = module_id

    required_module_fields = {
        "id",
        "version",
        "dependsOn",
        "conflictsWith",
        "rootBlocks",
        "supportingGuides",
        "rules",
        "requiredSkills",
        "skillDispatch",
        "requiredDecisions",
    }
    for module_id, module in sorted(versioned_modules.items()):
        module_fields = set(required_module_fields)
        if module.get("schemaVersion") == MODULE_SCHEMA_VERSION_V3:
            module_fields.add("repositoryExtensions")
        _validate_required_fields(module, module_fields, "module", module_id, diagnostics)
        dependencies = _validated_string_list(
            module.get("dependsOn"), f"module.dependsOn.invalid: {module_id}", diagnostics
        )
        conflicts = _validated_string_list(
            module.get("conflictsWith"),
            f"module.conflictsWith.invalid: {module_id}",
            diagnostics,
        )
        required_decisions = _validated_string_list(
            module.get("requiredDecisions"),
            f"module.requiredDecisions.invalid: {module_id}",
            diagnostics,
        )
        for dependency_id in dependencies:
            if dependency_id not in modules:
                diagnostics.append(
                    f"module.dependency.unknown: {module_id} -> {dependency_id}"
                )
        for conflict_id in conflicts:
            if conflict_id not in modules:
                diagnostics.append(f"module.conflict.unknown: {module_id} -> {conflict_id}")
        for decision_id in required_decisions:
            if decision_id not in decisions:
                diagnostics.append(f"module.decision.unknown: {module_id} -> {decision_id}")

        raw_rules = _validated_dict_list(
            module.get("rules"), f"module.rules.invalid: {module_id}", diagnostics
        )
        for rule in raw_rules:
            rule_id = rule.get("id")
            if not isinstance(rule_id, str) or not rule_id:
                diagnostics.append(f"rule.id.missing: {module_id}")
                continue
            owner = rule_owners.get(rule_id)
            if owner is not None:
                diagnostics.append(f"rule.id.duplicate: {rule_id}: {owner}, {module_id}")
            rule_owners[rule_id] = module_id
            if not isinstance(rule.get("version"), int) or rule["version"] < 1:
                diagnostics.append(f"rule.version.invalid: {rule_id}")
            raw_coverage = _validated_string_list(
                rule.get("coverage"), f"rule.coverage.invalid: {rule_id}", diagnostics
            )
            if len(set(raw_coverage)) != len(raw_coverage):
                diagnostics.append(f"rule.coverage.duplicate: {rule_id}")
            for coverage_id in raw_coverage:
                if coverage_id not in coverage:
                    diagnostics.append(f"rule.coverage.unknown: {rule_id} -> {coverage_id}")
            if module.get("schemaVersion") == MODULE_SCHEMA_VERSION_V3:
                clauses = _validate_clause_contracts(
                    module_id,
                    rule_id,
                    rule.get("clauses"),
                    clause_owners,
                    diagnostics,
                )
                rule_contracts[rule_id] = RuleContract(
                    rule_id=rule_id,
                    coverage=tuple(sorted(raw_coverage)),
                    clauses=clauses,
                )
            else:
                guidance = rule.get("guidance")
                if not isinstance(guidance, str) or not guidance.strip():
                    diagnostics.append(f"rule.guidance.invalid: {rule_id}")
                    continue
                rule_contracts[rule_id] = RuleContract(
                    rule_id=rule_id,
                    coverage=tuple(raw_coverage),
                    guidance=guidance,
                )

        required_skills = _validated_string_list(
            module.get("requiredSkills"),
            f"module.requiredSkills.invalid: {module_id}",
            diagnostics,
        )
        if len(set(required_skills)) != len(required_skills):
            diagnostics.append(f"module.requiredSkills.duplicate: {module_id}")
        dispatch_by_module[module_id] = _validate_skill_dispatch(
            module_id,
            module.get("skillDispatch"),
            required_skills,
            module.get("schemaVersion"),
            trigger_owners,
            diagnostics,
        )

        for field, kind in (("rootBlocks", "root-block"), ("supportingGuides", "guide")):
            raw_artifacts = _validated_dict_list(
                module.get(field), f"module.{field}.invalid: {module_id}", diagnostics
            )
            for artifact in raw_artifacts:
                artifact_id = artifact.get("id")
                if not isinstance(artifact_id, str) or not artifact_id:
                    diagnostics.append(f"managed.artifact.id.missing: {module_id}")
                    continue
                owners = block_owners if kind == "root-block" else guide_owners
                diagnostic_kind = "managed.block" if kind == "root-block" else "guide"
                owner = owners.get(artifact_id)
                if owner is not None:
                    diagnostics.append(
                        f"{diagnostic_kind}.id.duplicate: {artifact_id}: {owner}, {module_id}"
                    )
                owners[artifact_id] = module_id
                if not isinstance(artifact.get("version"), int) or artifact["version"] < 1:
                    diagnostics.append(f"{diagnostic_kind}.version.invalid: {artifact_id}")
                if kind == "guide":
                    target_path = artifact.get("path")
                    if not isinstance(target_path, str) or _is_unsafe_relative_path(target_path):
                        diagnostics.append(f"guide.path.invalid: {artifact_id}")
                _validate_template_reference(
                    artifact.get("template"), kind, templates, diagnostics, artifact_id
                )
                artifact_rules = _validated_string_list(
                    artifact.get("rules"),
                    f"rule.references.invalid: {artifact_id}",
                    diagnostics,
                )
                if len(set(artifact_rules)) != len(artifact_rules):
                    diagnostics.append(f"rule.reference.duplicate: {artifact_id}")
                module_rule_ids = {item.get("id") for item in raw_rules}
                for rule_id in artifact_rules:
                    if rule_id not in module_rule_ids:
                        diagnostics.append(f"rule.reference.unknown: {artifact_id} -> {rule_id}")
                if "references" not in artifact:
                    diagnostics.append(f"reference.collection.missing: {artifact_id}")
                artifacts[artifact_id] = (module_id, kind, artifact)

        if module.get("schemaVersion") == MODULE_SCHEMA_VERSION_V3:
            extensions = _validated_dict_list(
                module.get("repositoryExtensions"),
                f"module.repositoryExtensions.invalid: {module_id}",
                diagnostics,
            )
            raw_extensions.extend((module_id, extension) for extension in extensions)

    repository_extensions = _validate_repository_extensions(
        raw_extensions,
        artifacts,
        templates,
        decisions,
        diagnostics,
    )

    for artifact_id, (_, _, artifact) in sorted(artifacts.items()):
        references_by_artifact[artifact_id] = _validate_artifact_references(
            artifact_id,
            artifact,
            artifacts,
            templates,
            template_tokens,
            rendered_template_tokens,
            reference_owners,
            diagnostics,
        )

    return (
        rule_contracts,
        dispatch_by_module,
        references_by_artifact,
        repository_extensions,
    )


def _validate_clause_contracts(
    module_id: str,
    rule_id: str,
    raw_clauses: object,
    clause_owners: dict[str, str],
    diagnostics: list[str],
) -> tuple[ClauseContract, ...]:
    clauses = _validated_dict_list(
        raw_clauses, f"rule.clauses.invalid: {rule_id}", diagnostics
    )
    if not clauses:
        diagnostics.append(f"rule.clauses.empty: {rule_id}")
    normalized: list[ClauseContract] = []
    expected_fields = {"id", "enforcement", "guidance"}
    for clause in clauses:
        clause_id = clause.get("id")
        owner_label = clause_id if isinstance(clause_id, str) and clause_id else rule_id
        _validate_exact_fields(
            clause, expected_fields, "clause", owner_label, diagnostics
        )
        if not isinstance(clause_id, str) or not clause_id.startswith("clause."):
            diagnostics.append(f"clause.id.invalid: {rule_id}")
            continue
        owner = clause_owners.get(clause_id)
        if owner is not None:
            diagnostics.append(
                f"clause.id.duplicate: {clause_id}: {owner}, {module_id}"
            )
        clause_owners[clause_id] = module_id
        enforcement = clause.get("enforcement")
        if enforcement not in CLAUSE_ENFORCEMENTS:
            diagnostics.append(
                f"clause.enforcement.invalid: {clause_id} -> {enforcement}"
            )
            continue
        guidance = clause.get("guidance")
        if not isinstance(guidance, str) or not guidance.strip():
            diagnostics.append(f"clause.guidance.invalid: {clause_id}")
            continue
        normalized.append(
            ClauseContract(
                clause_id=clause_id,
                enforcement=enforcement,
                guidance=guidance,
            )
        )
    return tuple(sorted(normalized, key=lambda clause: clause.clause_id))


def _validate_skill_dispatch(
    module_id: str,
    raw_dispatch: object,
    required_skills: list[str],
    schema_version: object,
    trigger_owners: dict[str, str],
    diagnostics: list[str],
) -> tuple[SkillDispatch, ...]:
    if schema_version == MODULE_SCHEMA_VERSION_V3:
        return _validate_triggered_skill_dispatch(
            module_id,
            raw_dispatch,
            required_skills,
            trigger_owners,
            diagnostics,
        )

    items = _validated_dict_list(
        raw_dispatch, f"module.skillDispatch.invalid: {module_id}", diagnostics
    )
    dispatch: list[SkillDispatch] = []
    seen: set[str] = set()
    for item in items:
        skill_name = item.get("id")
        when = item.get("when")
        if not isinstance(skill_name, str) or not skill_name:
            diagnostics.append(f"skill.dispatch.id.missing: {module_id}")
            continue
        if skill_name in seen:
            diagnostics.append(f"skill.dispatch.id.duplicate: {module_id} -> {skill_name}")
        seen.add(skill_name)
        if not isinstance(when, str) or not when.strip():
            diagnostics.append(f"skill.dispatch.when.invalid: {module_id} -> {skill_name}")
            continue
        dispatch.append(
            SkillDispatch(
                skill_name=skill_name,
                when=when,
                trigger_id=skill_name,
                owner_module=module_id,
            )
        )

    required = set(required_skills)
    declared = set(seen)
    if required != declared:
        missing = ",".join(sorted(required - declared)) or "-"
        extra = ",".join(sorted(declared - required)) or "-"
        diagnostics.append(
            f"module.skillDispatch.mismatch: {module_id}: missing={missing}: extra={extra}"
        )
    return tuple(dispatch)


def _validate_triggered_skill_dispatch(
    module_id: str,
    raw_dispatch: object,
    required_skills: list[str],
    trigger_owners: dict[str, str],
    diagnostics: list[str],
) -> tuple[SkillDispatch, ...]:
    items = _validated_dict_list(
        raw_dispatch, f"module.skillDispatch.invalid: {module_id}", diagnostics
    )
    dispatch: list[SkillDispatch] = []
    seen_skills: set[str] = set()
    for item in items:
        skill_name = item.get("skill")
        if set(item) != {"skill", "triggers"}:
            diagnostics.append(f"skill.dispatch.fields.invalid: {module_id}")
        if not isinstance(skill_name, str) or not skill_name:
            diagnostics.append(f"skill.dispatch.skill.invalid: {module_id}")
            continue
        if skill_name in seen_skills:
            diagnostics.append(
                f"skill.dispatch.skill.duplicate: {module_id} -> {skill_name}"
            )
        seen_skills.add(skill_name)
        triggers = _validated_dict_list(
            item.get("triggers"),
            f"skill.dispatch.triggers.invalid: {module_id} -> {skill_name}",
            diagnostics,
        )
        if not triggers:
            diagnostics.append(
                f"skill.dispatch.triggers.empty: {module_id} -> {skill_name}"
            )
        for trigger in triggers:
            trigger_id = trigger.get("id")
            when = trigger.get("when")
            trigger_label = (
                trigger_id if isinstance(trigger_id, str) and trigger_id else skill_name
            )
            _validate_exact_fields(
                trigger,
                {"id", "when"},
                "skill.dispatch.trigger",
                trigger_label,
                diagnostics,
            )
            if not isinstance(trigger_id, str) or not trigger_id.startswith("trigger."):
                diagnostics.append(
                    f"skill.dispatch.trigger.id.invalid: {module_id} -> {skill_name}"
                )
                continue
            owner = trigger_owners.get(trigger_id)
            if owner is not None:
                diagnostics.append(
                    f"skill.dispatch.trigger.id.duplicate: {trigger_id}: {owner}, {module_id}"
                )
            trigger_owners[trigger_id] = module_id
            if not isinstance(when, str) or not when.strip():
                diagnostics.append(
                    f"skill.dispatch.trigger.when.invalid: {trigger_id}"
                )
                continue
            dispatch.append(
                SkillDispatch(
                    skill_name=skill_name,
                    when=when,
                    trigger_id=trigger_id,
                    owner_module=module_id,
                )
            )

    required = set(required_skills)
    declared = set(seen_skills)
    if required != declared:
        missing = ",".join(sorted(required - declared)) or "-"
        extra = ",".join(sorted(declared - required)) or "-"
        diagnostics.append(
            f"module.skillDispatch.mismatch: {module_id}: missing={missing}: extra={extra}"
        )
    return tuple(
        sorted(dispatch, key=lambda entry: (entry.skill_name, entry.trigger_id))
    )


def _validate_repository_extensions(
    raw_extensions: list[tuple[str, dict]],
    artifacts: dict[str, tuple[str, str, dict]],
    templates: dict[str, dict],
    decisions: dict[str, dict],
    diagnostics: list[str],
) -> dict[str, RepositoryOwnedExtension]:
    extensions: dict[str, RepositoryOwnedExtension] = {}
    expected_fields = {"id", "path", "template", "rootPointer", "decision"}
    for module_id, item in sorted(
        raw_extensions, key=lambda pair: str(pair[1].get("id", ""))
    ):
        extension_id = item.get("id")
        owner_label = (
            extension_id if isinstance(extension_id, str) and extension_id else module_id
        )
        _validate_exact_fields(
            item, expected_fields, "extension", owner_label, diagnostics
        )
        if not isinstance(extension_id, str) or not extension_id.startswith("extension."):
            diagnostics.append(f"extension.id.invalid: {module_id}")
            continue
        if extension_id in extensions:
            diagnostics.append(f"extension.id.duplicate: {extension_id}")
        target_path = item.get("path")
        if not isinstance(target_path, str) or _is_unsafe_relative_path(target_path):
            diagnostics.append(f"extension.path.invalid: {extension_id}")
            continue
        template_id = item.get("template")
        template = templates.get(template_id) if isinstance(template_id, str) else None
        if template is None:
            diagnostics.append(f"extension.template.unknown: {extension_id} -> {template_id}")
            continue
        if template.get("kind") != "repository-extension":
            diagnostics.append(f"extension.template.kind.invalid: {extension_id}")
        root_pointer_id = item.get("rootPointer")
        root_pointer = artifacts.get(root_pointer_id) if isinstance(root_pointer_id, str) else None
        if root_pointer is None or root_pointer[1] != "root-block":
            diagnostics.append(
                f"extension.rootPointer.invalid: {extension_id} -> {root_pointer_id}"
            )
            continue
        decision_id = item.get("decision")
        if not isinstance(decision_id, str) or decision_id not in decisions:
            diagnostics.append(
                f"extension.decision.unknown: {extension_id} -> {decision_id}"
            )
            continue
        extensions[extension_id] = RepositoryOwnedExtension(
            extension_id=extension_id,
            target_path=Path(target_path),
            template_id=template_id,
            root_pointer_id=root_pointer_id,
            decision_id=decision_id,
        )
    return extensions


def _validate_artifact_references(
    artifact_id: str,
    artifact: dict,
    artifacts: dict[str, tuple[str, str, dict]],
    templates: dict[str, dict],
    template_tokens: dict[str, set[str]],
    rendered_template_tokens: dict[str, set[str]],
    reference_owners: dict[str, str],
    diagnostics: list[str],
) -> tuple[ArtifactReference, ...]:
    items = _validated_dict_list(
        artifact.get("references"),
        f"reference.collection.invalid: {artifact_id}",
        diagnostics,
    )
    template_id = artifact.get("template")
    declared_tokens = template_tokens.get(template_id, set())
    rendered_tokens = rendered_template_tokens.get(template_id, set())
    references: list[ArtifactReference] = []
    bound_tokens: set[str] = set()
    for item in items:
        reference_id = item.get("id")
        token = item.get("token")
        ownership = item.get("ownership")
        if not isinstance(reference_id, str) or not reference_id:
            diagnostics.append(f"reference.id.missing: {artifact_id}")
            continue
        owner = reference_owners.get(reference_id)
        if owner is not None:
            diagnostics.append(f"reference.id.duplicate: {reference_id}: {owner}, {artifact_id}")
        reference_owners[reference_id] = artifact_id
        if not isinstance(token, str) or not token:
            diagnostics.append(f"reference.token.invalid: {reference_id}")
            continue
        if token in bound_tokens:
            diagnostics.append(f"reference.token.duplicate: {artifact_id} -> {token}")
        bound_tokens.add(token)
        if token not in declared_tokens or token not in rendered_tokens:
            diagnostics.append(f"reference.token.unknown: {reference_id} -> {token}")

        if ownership == "setup":
            managed_id = item.get("managedId")
            if not isinstance(managed_id, str) or not managed_id:
                diagnostics.append(f"reference.managed.invalid: {reference_id}")
                continue
            if managed_id not in artifacts:
                diagnostics.append(f"reference.managed.unknown: {reference_id} -> {managed_id}")
            if "path" in item:
                diagnostics.append(f"reference.shape.invalid: {reference_id}")
            references.append(
                ArtifactReference(
                    reference_id=reference_id,
                    token=token,
                    ownership=ownership,
                    target_managed_id=managed_id,
                )
            )
            continue

        if ownership == "repository":
            repository_path = item.get("path")
            if not isinstance(repository_path, str) or _is_unsafe_relative_path(repository_path):
                diagnostics.append(f"reference.repository.path.invalid: {reference_id}")
                continue
            if "managedId" in item:
                diagnostics.append(f"reference.shape.invalid: {reference_id}")
            references.append(
                ArtifactReference(
                    reference_id=reference_id,
                    token=token,
                    ownership=ownership,
                    repository_path=Path(repository_path),
                )
            )
            continue

        diagnostics.append(f"reference.ownership.unknown: {reference_id} -> {ownership}")
    return tuple(references)


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


def _validate_decision_contracts(
    decisions: dict[str, dict],
    diagnostics: list[str],
) -> None:
    for decision_id, decision in sorted(decisions.items()):
        decision_type = decision.get("type")
        if decision_type not in DECISION_TYPES:
            diagnostics.append(f"decision.type.invalid: {decision_id}: {decision_type}")
            continue

        values = decision.get("values")
        if decision_type != "enum":
            if "values" in decision:
                diagnostics.append(f"decision.values.invalid: {decision_id}")
            continue
        if (
            not isinstance(values, list)
            or not values
            or not all(isinstance(value, str) and value for value in values)
        ):
            diagnostics.append(f"decision.values.invalid: {decision_id}")
            continue
        if len(set(values)) != len(values):
            diagnostics.append(f"decision.values.duplicate: {decision_id}")


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


def _validate_required_fields(
    data: dict,
    fields: set[str],
    kind: str,
    asset_id: str,
    diagnostics: list[str],
) -> None:
    for field in sorted(fields - set(data)):
        diagnostics.append(f"{kind}.field.missing: {asset_id} -> {field}")


def _validate_exact_fields(
    data: dict,
    fields: set[str],
    kind: str,
    asset_id: str,
    diagnostics: list[str],
) -> None:
    _validate_required_fields(data, fields, kind, asset_id, diagnostics)
    for field in sorted(set(data) - fields):
        diagnostics.append(f"{kind}.field.unknown: {asset_id} -> {field}")


def _validated_string_list(
    raw_items: object,
    diagnostic: str,
    diagnostics: list[str],
) -> list[str]:
    if not isinstance(raw_items, list) or not all(
        isinstance(item, str) and item for item in raw_items
    ):
        diagnostics.append(diagnostic)
        return []
    return list(raw_items)


def _validated_dict_list(
    raw_items: object,
    diagnostic: str,
    diagnostics: list[str],
) -> list[dict]:
    if not isinstance(raw_items, list) or not all(
        isinstance(item, dict) for item in raw_items
    ):
        diagnostics.append(diagnostic)
        return []
    return list(raw_items)


def _dict_items(raw_items: object) -> list[dict]:
    if not isinstance(raw_items, list):
        return []
    return [item for item in raw_items if isinstance(item, dict)]


def _string_items(raw_items: object) -> list[str]:
    if not isinstance(raw_items, list):
        return []
    return [item for item in raw_items if isinstance(item, str)]


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


def _validate_setups(
    setups: dict[str, dict], diagnostics: list[str]
) -> dict[str, tuple[ExternalSkillContract, ...]]:
    external_by_setup: dict[str, tuple[ExternalSkillContract, ...]] = {}
    for setup_id, setup in sorted(setups.items()):
        if setup.get("schemaVersion") == SETUP_SCHEMA_VERSION_V2:
            external_by_setup[setup_id] = _validate_versioned_setup(
                setup_id, setup, diagnostics
            )
            continue
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
    return external_by_setup


def _validate_versioned_setup(
    setup_id: str, setup: dict, diagnostics: list[str]
) -> tuple[ExternalSkillContract, ...]:
    _validate_required_fields(
        setup,
        {"id", "version", "source", "digest", "skills"},
        "setup",
        setup_id,
        diagnostics,
    )
    if setup.get("version") != 2:
        diagnostics.append(f"setup.version.invalid: {setup_id}")
    source_metadata = setup.get("source")
    if not isinstance(source_metadata, dict):
        diagnostics.append(f"setup.source.invalid: {setup_id}")
    else:
        _validate_github_source(
            source_metadata,
            f"setup.source: {setup_id}",
            diagnostics,
        )
    skills = _validated_dict_list(
        setup.get("skills"), f"setup.skills.invalid: {setup_id}", diagnostics
    )
    external: list[ExternalSkillContract] = []
    normalized_records: list[dict] = []
    seen_names: set[str] = set()
    for skill in skills:
        name = skill.get("name")
        target_path = skill.get("path")
        if not isinstance(name, str) or not name:
            diagnostics.append(f"setup.skill.name.missing: {setup_id}")
            continue
        if name in seen_names:
            diagnostics.append(f"setup.skill.name.duplicate: {setup_id}: {name}")
        seen_names.add(name)
        if not isinstance(target_path, str) or _is_unsafe_relative_path(target_path):
            diagnostics.append(f"setup.skill.path.invalid: {setup_id}: {name}")
            continue

        source = skill.get("source")
        if not isinstance(source, dict):
            diagnostics.append(f"setup.skill.source.invalid: {setup_id}: {name}")
            continue
        source_type = source.get("type")
        if source_type == "github":
            repository = source.get("repository")
            revision = source.get("ref")
            source_path = source.get("path")
            tree_digest = skill.get("treeDigest")
            valid = True
            if set(source) != {"type", "repository", "ref", "path"}:
                diagnostics.append(f"setup.skill.source.fields.invalid: {setup_id}: {name}")
                valid = False
            if set(skill) != {"name", "path", "source", "treeDigest"}:
                diagnostics.append(f"setup.skill.fields.invalid: {setup_id}: {name}")
                valid = False
            if not isinstance(repository, str) or not _is_github_repository(repository):
                diagnostics.append(f"setup.skill.source.repository.invalid: {setup_id}: {name}")
                valid = False
            if not isinstance(revision, str) or not IMMUTABLE_GIT_REF.fullmatch(revision):
                diagnostics.append(f"setup.skill.source.ref.mutable: {setup_id}: {name}")
                valid = False
            if not isinstance(source_path, str) or _is_unsafe_relative_path(source_path):
                diagnostics.append(f"setup.skill.source.path.invalid: {setup_id}: {name}")
                valid = False
            if not isinstance(tree_digest, str) or not LOWERCASE_SHA256.fullmatch(tree_digest):
                diagnostics.append(f"setup.skill.treeDigest.invalid: {setup_id}: {name}")
                valid = False
            normalized_records.append(
                {
                    "name": name,
                    "path": target_path,
                    "source": {
                        "type": source_type,
                        "repository": repository,
                        "ref": revision,
                        "path": source_path,
                    },
                    "treeDigest": tree_digest,
                }
            )
            if valid:
                source_ref = SkillSourceRef(
                    provider="github",
                    repository=repository,
                    revision=revision,
                    source_path=Path(source_path),
                )
                external.append(
                    ExternalSkillContract(
                        skill_name=name,
                        target_path=Path(target_path),
                        source=source_ref,
                        tree_digest=tree_digest,
                    )
                )
            continue

        if source_type == "repo":
            source_name = source.get("name")
            content_digest = skill.get("contentDigest")
            if set(source) != {"type", "name"}:
                diagnostics.append(f"setup.skill.source.fields.invalid: {setup_id}: {name}")
            if set(skill) != {"name", "path", "source", "contentDigest"}:
                diagnostics.append(f"setup.skill.fields.invalid: {setup_id}: {name}")
            if not isinstance(source_name, str) or not source_name:
                diagnostics.append(f"setup.skill.source.repo.invalid: {setup_id}: {name}")
            if not isinstance(content_digest, str) or not LOWERCASE_SHA256.fullmatch(
                content_digest
            ):
                diagnostics.append(f"setup.skill.contentDigest.invalid: {setup_id}: {name}")
            normalized_records.append(
                {
                    "name": name,
                    "path": target_path,
                    "source": {"type": source_type, "name": source_name},
                    "contentDigest": content_digest,
                }
            )
            continue

        diagnostics.append(f"setup.skill.source.type.unknown: {setup_id}: {name} -> {source_type}")

    declared_digest = setup.get("digest")
    expected_digest = _normalized_records_digest(normalized_records)
    if not isinstance(declared_digest, str) or not LOWERCASE_SHA256.fullmatch(declared_digest):
        diagnostics.append(f"setup.digest.invalid: {setup_id}")
    elif declared_digest != expected_digest:
        diagnostics.append(f"setup.digest.mismatch: {setup_id}: expected {expected_digest}")
    return tuple(external)


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


def _validate_profile_rule_contracts(
    profiles: dict[str, dict],
    ordered_modules_by_profile: dict[str, list[str]],
    modules: dict[str, dict],
    rule_contracts: dict[str, RuleContract],
    rendered_template_tokens: dict[str, set[str]],
    diagnostics: list[str],
) -> None:
    rule_owners: dict[str, str] = {}
    for module_id, module in sorted(modules.items()):
        if module.get("schemaVersion") not in {
            MODULE_SCHEMA_VERSION_V2,
            MODULE_SCHEMA_VERSION_V3,
        }:
            continue
        for rule in _dict_items(module.get("rules")):
            rule_id = rule.get("id")
            if isinstance(rule_id, str) and rule_id:
                rule_owners[rule_id] = module_id

    for profile_id, profile in sorted(profiles.items()):
        if profile.get("schemaVersion") not in {
            PROFILE_SCHEMA_VERSION_V2,
            PROFILE_SCHEMA_VERSION_V3,
        }:
            continue
        required_profile_fields = {
            "id",
            "version",
            "setup",
            "entryDecisions",
            "modules",
            "requiredRules",
        }
        if profile.get("schemaVersion") == PROFILE_SCHEMA_VERSION_V3:
            required_profile_fields.add("formatter")
        _validate_required_fields(
            profile,
            required_profile_fields,
            "profile",
            profile_id,
            diagnostics,
        )
        required_rules = _validated_string_list(
            profile.get("requiredRules"),
            f"profile.requiredRules.invalid: {profile_id}",
            diagnostics,
        )
        if len(set(required_rules)) != len(required_rules):
            diagnostics.append(f"profile.rule.required.duplicate: {profile_id}")
        selected_modules = set(ordered_modules_by_profile.get(profile_id, []))
        selected_rules = {
            rule.get("id")
            for module_id in selected_modules
            for rule in _dict_items(modules[module_id].get("rules"))
            if isinstance(rule.get("id"), str) and rule.get("id")
        }
        missing_rules = selected_rules - set(required_rules)
        extra_rules = set(required_rules) - selected_rules
        if missing_rules or extra_rules:
            diagnostics.append(
                f"profile.rule.required.mismatch: {profile_id}: "
                f"missing={','.join(sorted(missing_rules)) or '-'}: "
                f"extra={','.join(sorted(extra_rules)) or '-'}"
            )
        for rule_id in required_rules:
            if rule_id not in rule_contracts:
                diagnostics.append(f"profile.rule.required.unknown: {profile_id} -> {rule_id}")
                continue
            owner = rule_owners.get(rule_id)
            if owner not in selected_modules:
                diagnostics.append(
                    f"profile.rule.module.missing: {profile_id}: {rule_id} requires {owner}"
                )
                continue
            carriers = [
                guide
                for guide in _dict_items(modules[owner].get("supportingGuides"))
                if rule_id in _string_items(guide.get("rules"))
            ]
            if not carriers:
                diagnostics.append(f"profile.rule.carrier.missing: {profile_id} -> {rule_id}")
                continue
            if not any(
                "artifact.rules"
                in rendered_template_tokens.get(guide.get("template"), set())
                for guide in carriers
            ):
                diagnostics.append(f"profile.rule.binding.missing: {profile_id} -> {rule_id}")


def _validate_profile_formatters(
    profiles: dict[str, dict], diagnostics: list[str]
) -> dict[str, FormatterContract]:
    formatters: dict[str, FormatterContract] = {}
    selected_fields = {"kind", "id", "version", "fixturePaths", "goldenDigest"}
    for profile_id, profile in sorted(profiles.items()):
        if profile.get("schemaVersion") != PROFILE_SCHEMA_VERSION_V3:
            continue
        raw_formatter = profile.get("formatter")
        if not isinstance(raw_formatter, dict):
            diagnostics.append(f"profile.formatter.invalid: {profile_id}")
            continue
        kind = raw_formatter.get("kind")
        if kind not in FORMATTER_KINDS:
            diagnostics.append(f"profile.formatter.kind.invalid: {profile_id} -> {kind}")
            continue
        if kind == "none":
            _validate_exact_fields(
                raw_formatter, {"kind"}, "profile.formatter", profile_id, diagnostics
            )
            formatters[profile_id] = FormatterContract(kind=kind)
            continue

        _validate_exact_fields(
            raw_formatter,
            selected_fields,
            "profile.formatter",
            profile_id,
            diagnostics,
        )
        formatter_id = raw_formatter.get("id")
        if not isinstance(formatter_id, str) or not formatter_id.startswith("formatter."):
            diagnostics.append(f"profile.formatter.id.invalid: {profile_id}")
            continue
        version = raw_formatter.get("version")
        if not isinstance(version, str) or not version.strip():
            diagnostics.append(f"profile.formatter.version.invalid: {profile_id}")
            continue
        raw_paths = raw_formatter.get("fixturePaths")
        if not isinstance(raw_paths, list) or not raw_paths:
            diagnostics.append(f"profile.formatter.fixturePaths.invalid: {profile_id}")
            continue
        fixture_paths: list[Path] = []
        seen_paths: set[str] = set()
        paths_valid = True
        for path in raw_paths:
            if not isinstance(path, str) or _is_unsafe_relative_path(path):
                diagnostics.append(
                    f"profile.formatter.fixturePath.invalid: {profile_id} -> {path}"
                )
                paths_valid = False
                continue
            if path in seen_paths:
                diagnostics.append(
                    f"profile.formatter.fixturePath.duplicate: {profile_id} -> {path}"
                )
            seen_paths.add(path)
            fixture_paths.append(Path(path))
        golden_digest = raw_formatter.get("goldenDigest")
        if not isinstance(golden_digest, str) or not LOWERCASE_SHA256.fullmatch(
            golden_digest
        ):
            diagnostics.append(f"profile.formatter.goldenDigest.invalid: {profile_id}")
            continue
        if not paths_valid:
            continue
        formatters[profile_id] = FormatterContract(
            kind=kind,
            formatter_id=formatter_id,
            version=version,
            fixture_paths=tuple(sorted(fixture_paths, key=lambda path: path.as_posix())),
            golden_digest=golden_digest,
        )
    return formatters


def _normalize_profile_dispatch_contracts(
    profiles: dict[str, dict],
    ordered_modules_by_profile: dict[str, list[str]],
    dispatch_by_module: dict[str, tuple[SkillDispatch, ...]],
    diagnostics: list[str],
) -> dict[str, tuple[SkillDispatch, ...]]:
    normalized: dict[str, dict[tuple[str, str], SkillDispatch]] = {}
    for profile_id, profile in sorted(profiles.items()):
        if profile.get("schemaVersion") != PROFILE_SCHEMA_VERSION_V3:
            continue
        owners: dict[str, str] = {}
        for module_id in ordered_modules_by_profile.get(profile_id, []):
            for entry in dispatch_by_module.get(module_id, ()):
                owner = owners.get(entry.skill_name)
                if owner is not None and owner != module_id:
                    diagnostics.append(
                        f"skill.dispatch.owner.duplicate: {profile_id}: "
                        f"{entry.skill_name}: {owner}, {module_id}"
                    )
                owners[entry.skill_name] = module_id
                normalized.setdefault(entry.skill_name, {})[
                    (entry.trigger_id, entry.owner_module)
                ] = entry
    return {
        skill_name: tuple(
            sorted(entries.values(), key=lambda entry: entry.trigger_id)
        )
        for skill_name, entries in sorted(normalized.items())
    }


def _validate_upgrade_transitions(
    transitions: dict[str, dict],
    rule_contracts: dict[str, RuleContract],
    repository_extensions: dict[str, RepositoryOwnedExtension],
    diagnostics: list[str],
) -> dict[str, UpgradeTransition]:
    current_clauses = {
        clause.clause_id: clause
        for rule in rule_contracts.values()
        for clause in rule.clauses
    }
    validated: dict[str, UpgradeTransition] = {}
    transition_fields = {
        "schemaVersion",
        "id",
        "version",
        "fromBaseline",
        "toBaseline",
        "priorClauses",
        "mappings",
    }
    prior_fields = {"id", "enforcement", "carrier", "guidanceDigest"}
    mapping_fields = {"fromClause", "disposition", "targets", "reason"}

    for transition_id, transition in sorted(transitions.items()):
        _validate_exact_fields(
            transition, transition_fields, "transition", transition_id, diagnostics
        )
        version = transition.get("version")
        if version != 1:
            diagnostics.append(f"transition.version.invalid: {transition_id}")
        from_baseline = transition.get("fromBaseline")
        to_baseline = transition.get("toBaseline")
        if not isinstance(from_baseline, str) or not from_baseline.startswith("baseline."):
            diagnostics.append(f"transition.fromBaseline.invalid: {transition_id}")
            continue
        if not isinstance(to_baseline, str) or not to_baseline.startswith("baseline."):
            diagnostics.append(f"transition.toBaseline.invalid: {transition_id}")
            continue
        if from_baseline == to_baseline:
            diagnostics.append(f"transition.baseline.same: {transition_id}")

        raw_prior_clauses = _validated_dict_list(
            transition.get("priorClauses"),
            f"transition.priorClauses.invalid: {transition_id}",
            diagnostics,
        )
        if not raw_prior_clauses:
            diagnostics.append(f"transition.priorClauses.empty: {transition_id}")
        prior_clauses: list[PriorClauseContract] = []
        prior_enforcement: dict[str, str] = {}
        for prior in raw_prior_clauses:
            clause_id = prior.get("id")
            clause_label = (
                clause_id if isinstance(clause_id, str) and clause_id else transition_id
            )
            _validate_exact_fields(
                prior, prior_fields, "transition.priorClause", clause_label, diagnostics
            )
            if not isinstance(clause_id, str) or not clause_id.startswith("clause."):
                diagnostics.append(f"transition.priorClause.id.invalid: {transition_id}")
                continue
            if clause_id in prior_enforcement:
                diagnostics.append(
                    f"transition.priorClause.duplicate: {transition_id} -> {clause_id}"
                )
            enforcement = prior.get("enforcement")
            if enforcement not in CLAUSE_ENFORCEMENTS:
                diagnostics.append(
                    f"transition.priorClause.enforcement.invalid: {clause_id}"
                )
                continue
            carrier = prior.get("carrier")
            if not isinstance(carrier, str) or not carrier:
                diagnostics.append(f"transition.priorClause.carrier.invalid: {clause_id}")
                continue
            guidance_digest = prior.get("guidanceDigest")
            if not isinstance(guidance_digest, str) or not LOWERCASE_SHA256.fullmatch(
                guidance_digest
            ):
                diagnostics.append(f"transition.priorClause.digest.invalid: {clause_id}")
                continue
            prior_enforcement[clause_id] = enforcement
            prior_clauses.append(
                PriorClauseContract(
                    clause_id=clause_id,
                    enforcement=enforcement,
                    carrier_id=carrier,
                    guidance_digest=guidance_digest,
                )
            )

        raw_mappings = _validated_dict_list(
            transition.get("mappings"),
            f"transition.mappings.invalid: {transition_id}",
            diagnostics,
        )
        mappings: list[TransitionMapping] = []
        mapped_clauses: set[str] = set()
        for mapping in raw_mappings:
            from_clause = mapping.get("fromClause")
            mapping_label = (
                from_clause
                if isinstance(from_clause, str) and from_clause
                else transition_id
            )
            _validate_exact_fields(
                mapping, mapping_fields, "transition.mapping", mapping_label, diagnostics
            )
            if not isinstance(from_clause, str) or not from_clause:
                diagnostics.append(f"transition.mapping.fromClause.invalid: {transition_id}")
                continue
            if from_clause in mapped_clauses:
                diagnostics.append(
                    f"transition.mapping.duplicate: {transition_id} -> {from_clause}"
                )
            mapped_clauses.add(from_clause)
            if from_clause not in prior_enforcement:
                diagnostics.append(
                    f"transition.mapping.source.unknown: {transition_id} -> {from_clause}"
                )
            disposition = mapping.get("disposition")
            if disposition not in TRANSITION_DISPOSITIONS:
                diagnostics.append(
                    f"transition.disposition.invalid: {from_clause} -> {disposition}"
                )
                continue
            targets = _validated_string_list(
                mapping.get("targets"),
                f"transition.targets.invalid: {from_clause}",
                diagnostics,
            )
            if len(set(targets)) != len(targets):
                diagnostics.append(f"transition.target.duplicate: {from_clause}")
            if disposition == "rejected" and targets:
                diagnostics.append(f"transition.rejected.targets.invalid: {from_clause}")
            if disposition != "rejected" and not targets:
                diagnostics.append(f"transition.accepted.targets.missing: {from_clause}")
            for target in targets:
                current_clause = current_clauses.get(target)
                if current_clause is None and target not in repository_extensions:
                    diagnostics.append(f"transition.target.unknown: {from_clause} -> {target}")
                    continue
                if (
                    current_clause is not None
                    and prior_enforcement.get(from_clause) != current_clause.enforcement
                ):
                    diagnostics.append(
                        f"transition.target.enforcement.mismatch: {from_clause} -> {target}"
                    )
            reason = mapping.get("reason")
            if not isinstance(reason, str) or not reason.strip():
                diagnostics.append(f"transition.reason.invalid: {from_clause}")
                continue
            mappings.append(
                TransitionMapping(
                    from_clause=from_clause,
                    disposition=disposition,
                    targets=tuple(sorted(targets)),
                    reason=reason,
                )
            )

        prior_ids = set(prior_enforcement)
        if prior_ids != mapped_clauses:
            missing = ",".join(sorted(prior_ids - mapped_clauses)) or "-"
            extra = ",".join(sorted(mapped_clauses - prior_ids)) or "-"
            diagnostics.append(
                f"transition.mapping.incomplete: {transition_id}: "
                f"missing={missing}: extra={extra}"
            )
        validated[transition_id] = UpgradeTransition(
            transition_id=transition_id,
            version=version if isinstance(version, int) else 0,
            from_baseline=from_baseline,
            to_baseline=to_baseline,
            prior_clauses=tuple(
                sorted(prior_clauses, key=lambda clause: clause.clause_id)
            ),
            mappings=tuple(sorted(mappings, key=lambda mapping: mapping.from_clause)),
        )
    return validated


def _is_unsafe_relative_path(path: str) -> bool:
    candidate = Path(path)
    return (
        not path
        or candidate.as_posix() == "."
        or path != candidate.as_posix()
        or candidate.is_absolute()
        or ".." in candidate.parts
        or "\\" in path
        or re.match(r"^[A-Za-z]:", path) is not None
    )


def _is_github_repository(repository: str) -> bool:
    if not GITHUB_REPOSITORY.fullmatch(repository):
        return False
    owner, name = repository.split("/", 1)
    return owner not in {".", ".."} and name not in {".", ".."}


def _validate_github_source(
    source: dict,
    owner: str,
    diagnostics: list[str],
) -> None:
    if set(source) != {"type", "repository", "ref", "path"}:
        diagnostics.append(f"{owner}.fields.invalid")
    if source.get("type") != "github":
        diagnostics.append(f"{owner}.type.invalid")
    repository = source.get("repository")
    if not isinstance(repository, str) or not _is_github_repository(repository):
        diagnostics.append(f"{owner}.repository.invalid")
    revision = source.get("ref")
    if not isinstance(revision, str) or not IMMUTABLE_GIT_REF.fullmatch(revision):
        diagnostics.append(f"{owner}.ref.mutable")
    source_path = source.get("path")
    if not isinstance(source_path, str) or _is_unsafe_relative_path(source_path):
        diagnostics.append(f"{owner}.path.invalid")


def _paths_digest(paths: list[str]) -> str:
    payload = "\n".join(paths) + "\n"
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def _normalized_records_digest(records: list[dict]) -> str:
    payload = json.dumps(
        records,
        sort_keys=True,
        separators=(",", ":"),
        ensure_ascii=True,
    )
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def setup_records_digest(records: list[dict]) -> str:
    """Return the canonical digest for complete normalized setup records."""

    return _normalized_records_digest(records)


def portable_tree_digest(root: str | Path) -> str:
    """Hash a complete safe skill tree with portable, unambiguous framing."""

    tree_root = Path(root)
    try:
        root_stat = tree_root.lstat()
    except OSError as error:
        raise PortableTreeError(f"cannot inspect skill tree {tree_root}: {error}") from error
    if stat.S_ISLNK(root_stat.st_mode) or not stat.S_ISDIR(root_stat.st_mode):
        raise PortableTreeError(f"skill tree root is not a regular directory: {tree_root}")

    files: list[tuple[bytes, Path]] = []
    pending = [tree_root]
    while pending:
        directory = pending.pop()
        try:
            entries = list(os.scandir(directory))
        except OSError as error:
            raise PortableTreeError(f"cannot read skill tree directory {directory}: {error}") from error
        for entry in entries:
            path = Path(entry.path)
            relative = path.relative_to(tree_root)
            relative_bytes = os.fsencode(relative.as_posix())
            try:
                entry_stat = entry.stat(follow_symlinks=False)
            except OSError as error:
                raise PortableTreeError(f"cannot inspect skill tree entry {relative.as_posix()}: {error}") from error
            mode = entry_stat.st_mode
            if stat.S_ISLNK(mode):
                raise PortableTreeError(f"skill tree entry is a symbolic link: {relative.as_posix()}")
            if stat.S_ISDIR(mode):
                if entry.name not in {".git", "node_modules"}:
                    pending.append(path)
                continue
            if not stat.S_ISREG(mode):
                raise PortableTreeError(f"skill tree entry is not a regular file: {relative.as_posix()}")
            if entry_stat.st_nlink != 1:
                raise PortableTreeError(f"skill tree entry is hard linked: {relative.as_posix()}")
            files.append((relative_bytes, path))

    framed_files: list[tuple[bytes, bytes]] = []
    for relative_bytes, path in files:
        try:
            content = path.read_bytes()
        except OSError as error:
            relative = os.fsdecode(relative_bytes)
            raise PortableTreeError(f"cannot read skill tree file {relative}: {error}") from error
        framed_files.append((relative_bytes, content))
    return portable_file_digest(framed_files)


def portable_file_digest(files: Iterable[tuple[bytes, bytes]]) -> str:
    """Hash relative path/content pairs using the portable tree framing."""

    digest = hashlib.sha256()
    for relative_bytes, content in sorted(files, key=lambda item: item[0]):
        digest.update(len(relative_bytes).to_bytes(8, "big"))
        digest.update(relative_bytes)
        digest.update(len(content).to_bytes(8, "big"))
        digest.update(content)
    return digest.hexdigest()
