"""Load and validate setup-context-driven portable assets."""

from __future__ import annotations

import copy
import hashlib
import json
from dataclasses import dataclass
from pathlib import Path


ASSET_SCHEMA_VERSION = "setup-context-driven/assets-v1"
DECISIONS_SCHEMA_VERSION = "setup-context-driven/decisions-v1"
MODULE_SCHEMA_VERSION = "setup-context-driven/module-v1"
PROFILE_SCHEMA_VERSION = "setup-context-driven/profile-v1"
SETUP_SCHEMA_VERSION = "setup-context-driven/setup-snapshot-v1"
TEMPLATES_SCHEMA_VERSION = "setup-context-driven/templates-v1"


class AssetValidationError(Exception):
    """Raised when bundled setup-context-driven assets are invalid."""

    def __init__(self, diagnostics: list[str]):
        self.diagnostics = diagnostics
        super().__init__("\n".join(diagnostics))


@dataclass(frozen=True)
class AssetCatalog:
    decisions: dict[str, dict]
    modules: dict[str, dict]
    profiles: dict[str, dict]
    setups: dict[str, dict]
    templates: dict[str, dict]
    ordered_modules_by_profile: dict[str, list[str]]


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
    _validate_templates(assets_root / "templates", templates, diagnostics)
    _validate_modules(modules, decisions, templates, diagnostics)
    _validate_setups(setups, diagnostics)

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
) -> None:
    for template_id, template in sorted(templates.items()):
        path = template.get("path")
        if not isinstance(path, str) or _is_unsafe_relative_path(path):
            diagnostics.append(f"template.path.invalid: {template_id}")
            continue
        if not (templates_root / path).is_file():
            diagnostics.append(f"template.file.missing: {template_id}: {path}")


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
