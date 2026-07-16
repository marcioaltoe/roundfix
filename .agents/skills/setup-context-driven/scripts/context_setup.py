#!/usr/bin/env python3
"""Audit setup-context-driven managed agent instructions."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from datetime import date
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable

from context_assets import AssetCatalog, AssetValidationError, load_asset_catalog


AUDIT_SCHEMA_VERSION = "setup-context-driven/audit-v1"
MANIFEST_PATH = Path("docs/agents/setup-context.json")
ROOT_INSTRUCTIONS_PATH = Path("AGENTS.md")
SEVERITY_ORDER = {"error": 0, "decision": 1, "warning": 2, "info": 3}
BEGIN_MARKER = re.compile(
    r"<!--\s*setup-context-driven:begin\s+id=([A-Za-z0-9_.-]+)\s+version=([0-9]+)\s*-->"
)
END_MARKER = re.compile(
    r"<!--\s*setup-context-driven:end\s+id=([A-Za-z0-9_.-]+)\s*-->"
)
MARKER = re.compile(r"<!--\s*setup-context-driven:(begin|end)\b[^>]*-->")
MARKDOWN_LINK = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
NON_ENGLISH_MARKERS = [
    " não ",
    " obrigatório",
    " repositório",
    " arquivo",
    " configuración",
    " repositorio",
]
OPTIONAL_MODULE_DECISIONS = {
    "secondbrain": ("secondbrain.enabled", True),
}
SECONDBRAIN_REQUIRED_GUIDE_PHRASES = [
    "wiki/index.md",
    "qmd query",
    "projects/<project>/mirror",
    "Cite every Secondbrain file",
    "Do not write to the Secondbrain",
    "Do not edit raw/",
    "Do not edit projects/*/mirror/",
    "Hermes",
    "Never read, copy, or expose",
]


@dataclass(frozen=True)
class Finding:
    code: str
    severity: str
    path: str
    managed_id: str
    message: str
    action: str

    def to_json(self) -> dict[str, str]:
        return {
            "code": self.code,
            "severity": self.severity,
            "path": self.path,
            "managedId": self.managed_id,
            "message": self.message,
            "action": self.action,
        }


@dataclass(frozen=True)
class ExpectedArtifact:
    managed_id: str
    path: Path
    kind: str
    module_id: str
    template_id: str
    version: int
    content: str
    digest: str


@dataclass(frozen=True)
class PlanCondition:
    decision_id: str
    operator: str
    value: object

    def to_json(self) -> dict[str, object]:
        return {"decisionId": self.decision_id, self.operator: self.value}


@dataclass(frozen=True)
class PlannedArtifact:
    artifact: ExpectedArtifact
    present: bool
    state: str
    condition: PlanCondition | None = None


@dataclass(frozen=True)
class PlannedChange:
    action: str
    path: Path
    managed_id: str
    state: str
    condition: PlanCondition | None = None

    def to_json(self) -> dict[str, object]:
        data: dict[str, object] = {
            "action": self.action,
            "path": self.path.as_posix(),
            "managedId": self.managed_id,
            "state": self.state,
        }
        if self.condition is not None:
            data["condition"] = self.condition.to_json()
        return data


@dataclass(frozen=True)
class SelectionModule:
    module_id: str
    state: str
    condition: PlanCondition | None = None

    def to_json(self) -> dict[str, object]:
        data: dict[str, object] = {"id": self.module_id, "state": self.state}
        if self.condition is not None:
            data["condition"] = self.condition.to_json()
        return data


@dataclass(frozen=True)
class PreviewSelection:
    profile_id: str
    setup_id: str
    modules: tuple[SelectionModule, ...]

    def to_json(self) -> dict[str, object]:
        return {
            "profile": self.profile_id,
            "setup": self.setup_id,
            "modules": [module.to_json() for module in self.modules],
        }


@dataclass(frozen=True)
class DecisionPlan:
    profile_id: str
    setup_id: str
    active_modules: tuple[str, ...]
    resolved_decisions: dict[str, dict]
    unresolved_decisions: tuple[str, ...]
    selection: PreviewSelection
    artifacts: tuple[PlannedArtifact, ...]


@dataclass(frozen=True)
class ManagedBlock:
    managed_id: str
    version: int
    body: str


@dataclass(frozen=True)
class ManagedBlockSpan:
    managed_id: str
    version: int
    body: str
    start: int
    end: int


@dataclass(frozen=True)
class FileChange:
    path: Path
    content: str | None


@dataclass(frozen=True)
class ChangePlan:
    changes: list[FileChange]
    manifest: dict


@dataclass(frozen=True)
class InstalledSkill:
    name: str
    path: Path
    locked: bool
    origin: dict


@dataclass(frozen=True)
class AuditResult:
    findings: list[Finding]
    selection: PreviewSelection | None = None
    planned_changes: tuple[PlannedChange, ...] = ()

    @property
    def ok(self) -> bool:
        return not any(finding.severity in {"error", "decision"} for finding in self.findings)

    @property
    def summary(self) -> dict[str, int]:
        return {
            "errors": sum(1 for finding in self.findings if finding.severity == "error"),
            "decisions": sum(1 for finding in self.findings if finding.severity == "decision"),
            "warnings": sum(1 for finding in self.findings if finding.severity == "warning"),
            "info": sum(1 for finding in self.findings if finding.severity == "info"),
        }

    def to_json(self) -> dict:
        data = {
            "schemaVersion": AUDIT_SCHEMA_VERSION,
            "ok": self.ok,
            "summary": self.summary,
            "findings": [finding.to_json() for finding in sorted_findings(self.findings)],
            "plannedChanges": planned_changes_for_result(self),
        }
        if self.selection is not None:
            data["selection"] = self.selection.to_json()
        return data


def main(argv: list[str] | None = None) -> int:
    args = list(sys.argv[1:] if argv is None else argv)
    if args and args[0] in {"-h", "--help"}:
        print_top_level_help()
        return 0
    if args and args[0] == "apply":
        parser = apply_parser()
        return run_apply_command(parser.parse_args(args[1:]))
    if args and args[0] == "sync-setups":
        parser = sync_setups_parser()
        return run_sync_setups_command(parser.parse_args(args[1:]))
    if args and args[0] == "audit":
        args = args[1:]

    parser = audit_parser()
    options = parser.parse_args(args)
    return run_audit_command(options)


def run_audit_command(options: argparse.Namespace) -> int:
    repo = Path(options.repo).expanduser()
    if not repo.is_absolute():
        repo = Path.cwd() / repo
    repo = repo.resolve(strict=False)
    if not repo.is_dir():
        print(f"Repository root is not a directory: {repo}", file=sys.stderr)
        return 2

    skill_root = Path(__file__).resolve().parents[1]
    try:
        catalog = load_asset_catalog(skill_root)
    except AssetValidationError as error:
        for diagnostic in error.diagnostics:
            print(diagnostic, file=sys.stderr)
        return 2

    result, invalid_input = audit_repository(
        repo=repo,
        catalog=catalog,
        profile_override=options.profile,
        show_extra_skills=options.show_extra_skills,
        setups_dir=options.setups_dir,
    )
    render_result(result, options.format)
    return exit_code_for(result, invalid_input)


def run_apply_command(options: argparse.Namespace) -> int:
    repo = resolve_repo(options.repo)
    if not repo.is_dir():
        print(f"Repository root is not a directory: {repo}", file=sys.stderr)
        return 2

    skill_root = Path(__file__).resolve().parents[1]
    try:
        catalog = load_asset_catalog(skill_root)
    except AssetValidationError as error:
        for diagnostic in error.diagnostics:
            print(diagnostic, file=sys.stderr)
        return 2

    result, invalid_input, plan = plan_apply(
        repo=repo,
        catalog=catalog,
        profile_override=options.profile,
        decision_args=options.decision,
    )
    if result.summary["errors"] or result.summary["decisions"] or invalid_input:
        render_result(result, options.format)
        return exit_code_for(result, invalid_input)

    try:
        apply_change_plan(repo, plan)
    except OSError as error:
        failure = AuditResult(
            sorted_findings(
                [
                    finding(
                        "managed.apply.failed",
                        "error",
                        ".",
                        "apply",
                        f"Atomic apply failed: {error}.",
                        "Fix filesystem permissions and rerun apply.",
                    )
                ]
            )
        )
        render_result(failure, options.format)
        return 1

    applied = AuditResult(
        [
            finding(
                "managed.apply.completed",
                "info",
                ".",
                "apply",
                "Managed setup content matches the selected profile.",
                "No action needed.",
            )
        ]
    )
    render_result(applied, options.format)
    return 0


def run_sync_setups_command(options: argparse.Namespace) -> int:
    source_dir = Path(options.source_dir).expanduser()
    if not source_dir.is_absolute():
        source_dir = Path.cwd() / source_dir
    source_dir = source_dir.resolve(strict=False)
    if not source_dir.is_dir():
        print(f"Canonical setups directory is not a directory: {source_dir}", file=sys.stderr)
        return 2

    skill_root = Path(__file__).resolve().parents[1]
    result, invalid_input = sync_setup_snapshots(
        skill_root=skill_root,
        source_dir=source_dir,
        check=options.check,
    )
    render_result(result, options.format)
    return exit_code_for(result, invalid_input)


def resolve_repo(repo_arg: str) -> Path:
    repo = Path(repo_arg).expanduser()
    if not repo.is_absolute():
        repo = Path.cwd() / repo
    return repo.resolve(strict=False)


def audit_repository(
    repo: Path,
    catalog: AssetCatalog,
    profile_override: str | None = None,
    show_extra_skills: bool = False,
    setups_dir: str | None = None,
) -> tuple[AuditResult, bool]:
    findings: list[Finding] = []
    manifest, invalid_input = load_manifest(repo, findings)
    if manifest is None:
        if profile_override is not None and not invalid_input:
            findings.clear()
            return preview_result(
                repo=repo,
                catalog=catalog,
                profile_id=profile_override,
                existing_manifest=None,
                cli_decisions={},
                findings=findings,
                invalid_input=invalid_input,
            )
        return AuditResult(sorted_findings(findings)), invalid_input

    profile_id = profile_override or manifest.get("profile")
    if not isinstance(profile_id, str) or profile_id not in catalog.profiles:
        findings.append(
            finding(
                "profile.unknown",
                "error",
                str(MANIFEST_PATH),
                str(profile_id),
                f"Profile {profile_id!r} is not bundled.",
                "Select one bundled profile or update the manifest.",
            )
        )
        return AuditResult(sorted_findings(findings)), invalid_input

    decision_plan = resolve_decision_plan(
        catalog,
        profile_id,
        manifest,
        {},
    )
    decision_findings = decision_required_findings(decision_plan)
    if decision_findings:
        findings.extend(decision_findings)
        return AuditResult(
            sorted_findings(findings),
            selection=decision_plan.selection,
            planned_changes=planned_changes_for_plan(repo, decision_plan),
        ), invalid_input

    ordered_modules = list(decision_plan.active_modules)
    validate_manifest_shape(manifest, profile_id, ordered_modules, catalog, findings)
    validate_profile_skill_references(catalog, profile_id, ordered_modules, findings)
    skills_invalid_input = validate_installed_skills(
        repo=repo,
        catalog=catalog,
        manifest=manifest,
        profile_id=profile_id,
        show_extra_skills=show_extra_skills,
        findings=findings,
    )
    invalid_input = invalid_input or skills_invalid_input
    if setups_dir is not None:
        setups_path = Path(setups_dir).expanduser()
        if not setups_path.is_absolute():
            setups_path = Path.cwd() / setups_path
        setups_invalid_input = validate_selected_setup_snapshot(
            catalog=catalog,
            profile_id=profile_id,
            source_dir=setups_path.resolve(strict=False),
            findings=findings,
        )
        invalid_input = invalid_input or setups_invalid_input
    expected_artifacts = expected_artifacts_for_plan(decision_plan)
    validate_manifest_artifacts(manifest, expected_artifacts, findings)
    validate_documents(repo, expected_artifacts, findings)
    validate_secondbrain_documents(repo, ordered_modules, findings)

    return AuditResult(sorted_findings(findings)), invalid_input


def validate_profile_skill_references(
    catalog: AssetCatalog,
    profile_id: str,
    ordered_modules: list[str],
    findings: list[Finding],
) -> None:
    setup = catalog.setups[catalog.profiles[profile_id]["setup"]]
    setup_skills = {skill["name"] for skill in setup.get("skills", []) if isinstance(skill.get("name"), str)}
    for module_id in ordered_modules:
        for skill_name in catalog.modules[module_id].get("requiredSkills", []):
            if skill_name in setup_skills:
                continue
            findings.append(
                finding(
                    "skills.reference.outside-setup",
                    "error",
                    f"assets/modules/{module_id}.json",
                    skill_name,
                    f"Module {module_id} requires skill {skill_name}, but the selected setup snapshot does not include it.",
                    "Update the module or selected setup snapshot so generated rules reference available capabilities.",
                )
            )


def validate_installed_skills(
    repo: Path,
    catalog: AssetCatalog,
    manifest: dict,
    profile_id: str,
    show_extra_skills: bool,
    findings: list[Finding],
) -> bool:
    lock_entries, invalid_input = load_skills_lock(repo, findings)
    installed = discover_installed_skills(repo, lock_entries)
    installed_names = set(installed)
    setup_id = catalog.profiles[profile_id]["setup"]
    setup = catalog.setups[setup_id]
    required_names = [
        skill["name"]
        for skill in setup.get("skills", [])
        if isinstance(skill.get("name"), str) and skill["name"]
    ]
    for skill_name in required_names:
        if skill_name in installed_names:
            continue
        findings.append(
            finding(
                "skills.required.missing",
                "error",
                f".agents/skills/{skill_name}",
                f"profile.{profile_id}",
                f"Required skill {skill_name} is not installed.",
                f"Install the {setup_id} canonical skill setup or add .agents/skills/{skill_name}/SKILL.md.",
            )
        )

    if not show_extra_skills:
        return invalid_input

    local_skills = manifest_local_skills(manifest)
    setup_names = set(required_names)
    for skill_name, installed_skill in sorted(installed.items()):
        if skill_name in setup_names or skill_name in local_skills:
            continue
        if installed_skill.locked:
            findings.append(
                finding(
                    "skills.extra.installed",
                    "info",
                    installed_skill.path.as_posix(),
                    skill_name,
                    "Installed locked skill is outside the selected setup.",
                    "Review whether this locked skill is still useful; no file changes are performed.",
                )
            )
        else:
            findings.append(
                finding(
                    "skills.local.untracked",
                    "info",
                    installed_skill.path.as_posix(),
                    skill_name,
                    "Installed skill is not locked, setup-required, or manifest-declared repository-owned.",
                    "Review this untracked skill manually; the workflow will not modify it.",
                )
            )
    return invalid_input


def load_skills_lock(repo: Path, findings: list[Finding]) -> tuple[dict[str, dict], bool]:
    lock_path = repo / "skills-lock.json"
    try:
        lock = json.loads(lock_path.read_text(encoding="utf-8"))
    except FileNotFoundError:
        return {}, False
    except json.JSONDecodeError as error:
        findings.append(
            finding(
                "skills.lockfile.invalid",
                "error",
                "skills-lock.json",
                "skills-lock",
                f"skills-lock.json is not valid JSON: {error.msg}.",
                "Fix skills-lock.json before auditing installed skills.",
            )
        )
        return {}, True
    if not isinstance(lock, dict) or not isinstance(lock.get("skills"), dict):
        findings.append(
            finding(
                "skills.lockfile.invalid",
                "error",
                "skills-lock.json",
                "skills-lock",
                "skills-lock.json must be an object with a skills object.",
                "Fix skills-lock.json before auditing installed skills.",
            )
        )
        return {}, True
    entries: dict[str, dict] = {}
    for skill_name, entry in lock["skills"].items():
        if isinstance(skill_name, str) and isinstance(entry, dict):
            entries[skill_name] = entry
    return entries, False


def discover_installed_skills(
    repo: Path,
    lock_entries: dict[str, dict],
) -> dict[str, InstalledSkill]:
    skills_root = repo / ".agents" / "skills"
    installed: dict[str, InstalledSkill] = {}
    if not skills_root.is_dir():
        return installed
    for skill_file in sorted(skills_root.glob("*/SKILL.md")):
        skill_name = skill_file.parent.name
        lock_entry = lock_entries.get(skill_name)
        installed[skill_name] = InstalledSkill(
            name=skill_name,
            path=skill_file.relative_to(repo),
            locked=lock_entry is not None,
            origin=lock_entry or {},
        )
    return installed


def manifest_local_skills(manifest: dict) -> set[str]:
    local_skills: set[str] = set()
    for item in manifest.get("localSkills", []):
        if isinstance(item, str) and item:
            local_skills.add(item)
        elif isinstance(item, dict) and isinstance(item.get("name"), str) and item["name"]:
            local_skills.add(item["name"])
    return local_skills


def validate_selected_setup_snapshot(
    catalog: AssetCatalog,
    profile_id: str,
    source_dir: Path,
    findings: list[Finding],
) -> bool:
    if not source_dir.is_dir():
        findings.append(
            finding(
                "skills.setup-snapshot.drift",
                "error",
                str(source_dir),
                f"profile.{profile_id}",
                "Canonical setups directory is not available.",
                "Pass an existing --setups-dir or omit the optional drift check.",
            )
        )
        return True
    setup_id = catalog.profiles[profile_id]["setup"]
    snapshot, snapshot_findings, invalid_input = build_source_setup_snapshot(
        setup_id=setup_id,
        source_dir=source_dir,
        current_snapshot=catalog.setups[setup_id],
    )
    findings.extend(snapshot_findings)
    if invalid_input or snapshot is None:
        return invalid_input
    if snapshot != catalog.setups[setup_id]:
        findings.append(setup_snapshot_drift_finding(setup_id))
    return False


def ordered_modules_for_decisions(
    catalog: AssetCatalog,
    profile_id: str,
    decisions: dict,
) -> list[str]:
    ordered_modules = list(catalog.ordered_modules_by_profile[profile_id])
    for module_id, (decision_id, enabled_value) in OPTIONAL_MODULE_DECISIONS.items():
        if module_id not in catalog.modules:
            continue
        if decision_value(decisions, decision_id) == enabled_value and module_id not in ordered_modules:
            ordered_modules.append(module_id)
    return ordered_modules


def decision_value(decisions: dict, decision_id: str) -> object:
    if not isinstance(decisions, dict):
        return None
    decision = decisions.get(decision_id)
    if isinstance(decision, dict):
        return decision.get("value")
    return None


def preview_result(
    repo: Path,
    catalog: AssetCatalog,
    profile_id: str,
    existing_manifest: dict | None,
    cli_decisions: dict[str, object],
    findings: list[Finding],
    invalid_input: bool,
) -> tuple[AuditResult, bool]:
    if profile_id not in catalog.profiles:
        findings.append(
            finding(
                "profile.unknown",
                "error",
                str(MANIFEST_PATH),
                str(profile_id),
                f"Profile {profile_id!r} is not bundled.",
                "Select one bundled profile or update the manifest.",
            )
        )
        return AuditResult(sorted_findings(findings)), invalid_input

    decision_plan = resolve_decision_plan(
        catalog=catalog,
        profile_id=profile_id,
        existing_manifest=existing_manifest,
        cli_decisions=cli_decisions,
    )
    findings.extend(decision_required_findings(decision_plan))
    return (
        AuditResult(
            sorted_findings(findings),
            selection=decision_plan.selection,
            planned_changes=planned_changes_for_plan(repo, decision_plan),
        ),
        invalid_input,
    )


def resolve_decision_plan(
    catalog: AssetCatalog,
    profile_id: str,
    existing_manifest: dict | None,
    cli_decisions: dict[str, object],
) -> DecisionPlan:
    today = date.today().isoformat()
    existing = existing_manifest.get("decisions", {}) if isinstance(existing_manifest, dict) else {}
    if not isinstance(existing, dict):
        existing = {}

    required = list(catalog.profile_entry_decisions.get(profile_id, ()))
    if not required:
        required = required_decisions(catalog, catalog.ordered_modules_by_profile[profile_id])
    required_set = set(required)
    resolved: dict[str, dict] = {}
    unresolved: list[str] = []
    matched_effects = []
    unresolved_effects = []
    controlled_modules = controlled_module_ids(catalog)
    active_modules = {
        module_id
        for module_id in catalog.ordered_modules_by_profile[profile_id]
        if module_id not in controlled_modules
    }

    index = 0
    while index < len(required):
        decision_id = required[index]
        index += 1
        contract = catalog.decisions.get(decision_id)
        if contract is None:
            continue
        decision = compatible_decision_answer(
            contract=contract,
            decision_id=decision_id,
            existing=existing,
            cli_decisions=cli_decisions,
            today=today,
        )
        if decision is None:
            if decision_id not in unresolved:
                unresolved.append(decision_id)
            unresolved_effects.extend(catalog.decision_effects.get(decision_id, ()))
            continue

        resolved[decision_id] = decision
        for effect in catalog.decision_effects.get(decision_id, ()):
            if not plan_condition_matches(effect.condition, decision["value"]):
                continue
            matched_effects.append(effect)
            for module_id in effect.activate_modules:
                active_modules.add(module_id)
                for module_decision_id in catalog.modules[module_id].get("requiredDecisions", []):
                    if module_decision_id not in required_set:
                        required.append(module_decision_id)
                        required_set.add(module_decision_id)
            for dependent_decision_id in effect.require_decisions:
                if dependent_decision_id not in required_set:
                    required.append(dependent_decision_id)
                    required_set.add(dependent_decision_id)

    for decision_id, existing_decision in sorted(existing.items()):
        if decision_id.startswith("adoption.") and isinstance(existing_decision, dict):
            resolved[decision_id] = existing_decision
    for decision_id, value in sorted(cli_decisions.items()):
        if decision_id.startswith("adoption."):
            resolved[decision_id] = {"value": value, "confirmedAt": today}

    conditional_modules = conditional_module_conditions(
        catalog=catalog,
        unresolved_effects=unresolved_effects,
        active_modules=active_modules,
    )
    active_order = ordered_module_ids(catalog, profile_id, active_modules)
    selection = PreviewSelection(
        profile_id=profile_id,
        setup_id=catalog.profiles[profile_id]["setup"],
        modules=selection_modules(
            catalog=catalog,
            profile_id=profile_id,
            active_modules=set(active_order),
            conditional_modules=conditional_modules,
        ),
    )
    artifacts = planned_artifacts_for_decision_plan(
        catalog=catalog,
        profile_id=profile_id,
        active_modules=active_order,
        matched_effects=matched_effects,
        unresolved_effects=unresolved_effects,
    )
    return DecisionPlan(
        profile_id=profile_id,
        setup_id=catalog.profiles[profile_id]["setup"],
        active_modules=tuple(active_order),
        resolved_decisions=resolved,
        unresolved_decisions=tuple(unresolved),
        selection=selection,
        artifacts=artifacts,
    )


def compatible_decision_answer(
    contract: dict,
    decision_id: str,
    existing: dict,
    cli_decisions: dict[str, object],
    today: str,
) -> dict | None:
    if decision_id in cli_decisions:
        value = cli_decisions[decision_id]
        if is_decision_value_valid(contract, value):
            return {"value": value, "confirmedAt": today}
        return None

    existing_decision = existing.get(decision_id)
    if isinstance(existing_decision, dict) and "value" in existing_decision:
        value = existing_decision["value"]
        if is_decision_value_valid(contract, value):
            confirmed_at = existing_decision.get("confirmedAt")
            if not isinstance(confirmed_at, str) or not confirmed_at:
                confirmed_at = today
            return {"value": value, "confirmedAt": confirmed_at}
    return None


def plan_condition_matches(condition, value: object) -> bool:
    if condition.operator == "equals":
        return value == condition.value
    if condition.operator == "present":
        return (value is not None) == condition.value
    return False


def to_plan_condition(condition) -> PlanCondition:
    return PlanCondition(
        decision_id=condition.decision_id,
        operator=condition.operator,
        value=condition.value,
    )


def controlled_module_ids(catalog: AssetCatalog) -> set[str]:
    return {
        module_id
        for effects in catalog.decision_effects.values()
        for effect in effects
        for module_id in effect.activate_modules
    }


def controlled_artifact_ids(catalog: AssetCatalog) -> set[str]:
    controlled: set[str] = set()
    for effects in catalog.decision_effects.values():
        for effect in effects:
            controlled.update(effect.include_artifacts)
            controlled.update(effect.exclude_artifacts)
            controlled.update(selection.artifact_id for selection in effect.template_selections)
    return controlled


def conditional_module_conditions(
    catalog: AssetCatalog,
    unresolved_effects: list,
    active_modules: set[str],
) -> dict[str, PlanCondition]:
    conditions: dict[str, PlanCondition] = {}
    for effect in unresolved_effects:
        condition = to_plan_condition(effect.condition)
        for module_id in effect.activate_modules:
            if module_id in active_modules or module_id in conditions:
                continue
            conditions[module_id] = condition
    return conditions


def ordered_module_ids(
    catalog: AssetCatalog,
    profile_id: str,
    module_ids: set[str],
) -> list[str]:
    ordered: list[str] = []
    for module_id in catalog.ordered_modules_by_profile[profile_id]:
        if module_id in module_ids:
            ordered.append(module_id)
    for module_id in sorted(module_ids):
        if module_id not in ordered:
            ordered.append(module_id)
    return ordered


def selection_modules(
    catalog: AssetCatalog,
    profile_id: str,
    active_modules: set[str],
    conditional_modules: dict[str, PlanCondition],
) -> tuple[SelectionModule, ...]:
    module_order = list(catalog.ordered_modules_by_profile[profile_id])
    module_order.extend(
        module_id
        for module_id in sorted(conditional_modules)
        if module_id not in module_order
    )
    modules: list[SelectionModule] = []
    for module_id in module_order:
        if module_id in active_modules:
            modules.append(SelectionModule(module_id=module_id, state="active"))
        elif module_id in conditional_modules:
            modules.append(
                SelectionModule(
                    module_id=module_id,
                    state="conditional",
                    condition=conditional_modules[module_id],
                )
            )
    return tuple(modules)


def planned_artifacts_for_decision_plan(
    catalog: AssetCatalog,
    profile_id: str,
    active_modules: list[str],
    matched_effects: list,
    unresolved_effects: list,
) -> tuple[PlannedArtifact, ...]:
    artifact_modules = list(active_modules)
    for effect in [*matched_effects, *unresolved_effects]:
        for module_id in effect.activate_modules:
            if module_id not in artifact_modules:
                artifact_modules.append(module_id)
        effect_artifacts = list(effect.include_artifacts)
        effect_artifacts.extend(effect.exclude_artifacts)
        effect_artifacts.extend(selection.artifact_id for selection in effect.template_selections)
        for module_id in artifact_owner_modules(catalog, effect_artifacts):
            if module_id not in artifact_modules:
                artifact_modules.append(module_id)

    artifacts = {
        artifact.managed_id: artifact
        for artifact in expected_artifacts_for_profile(catalog, profile_id, artifact_modules)
    }
    module_artifacts = artifacts_by_module(artifacts.values())
    controlled_artifacts = controlled_artifact_ids(catalog)
    planned: list[PlannedArtifact] = []
    seen: set[tuple[str, bool, str, str, object]] = set()

    def add_artifact(
        managed_id: str,
        present: bool,
        state: str,
        condition: PlanCondition | None = None,
    ) -> None:
        artifact = artifacts.get(managed_id)
        if artifact is None:
            return
        condition_key = (
            condition.decision_id if condition else "",
            condition.operator if condition else "",
            condition.value if condition else "",
        )
        key = (managed_id, present, state, *condition_key)
        if key in seen:
            return
        seen.add(key)
        planned.append(
            PlannedArtifact(
                artifact=artifact,
                present=present,
                state=state,
                condition=condition,
            )
        )

    for module_id in active_modules:
        for artifact in module_artifacts.get(module_id, []):
            if artifact.managed_id not in controlled_artifacts:
                add_artifact(artifact.managed_id, True, "definite")

    for effect in matched_effects:
        for module_id in effect.activate_modules:
            for artifact in module_artifacts.get(module_id, []):
                add_artifact(artifact.managed_id, True, "definite")
        for managed_id in effect.include_artifacts:
            add_artifact(managed_id, True, "definite")
        for selection in effect.template_selections:
            add_artifact(selection.artifact_id, True, "definite")
        for managed_id in effect.exclude_artifacts:
            add_artifact(managed_id, False, "definite")

    for effect in unresolved_effects:
        condition = to_plan_condition(effect.condition)
        for module_id in effect.activate_modules:
            for artifact in module_artifacts.get(module_id, []):
                add_artifact(artifact.managed_id, True, "conditional", condition)
        for managed_id in effect.include_artifacts:
            add_artifact(managed_id, True, "conditional", condition)
        for selection in effect.template_selections:
            add_artifact(selection.artifact_id, True, "conditional", condition)
        for managed_id in effect.exclude_artifacts:
            add_artifact(managed_id, False, "conditional", condition)

    return tuple(planned)


def expected_artifacts_for_plan(decision_plan: DecisionPlan) -> list[ExpectedArtifact]:
    return [
        planned_artifact.artifact
        for planned_artifact in decision_plan.artifacts
        if planned_artifact.present and planned_artifact.state == "definite"
    ]


def artifact_owner_modules(catalog: AssetCatalog, managed_ids: Iterable[str]) -> list[str]:
    wanted = set(managed_ids)
    owner_modules: list[str] = []
    if not wanted:
        return owner_modules
    for module_id, module in sorted(catalog.modules.items()):
        module_artifact_ids = {
            block.get("id")
            for block in module.get("rootBlocks", [])
            if isinstance(block.get("id"), str)
        }
        module_artifact_ids.update(
            guide.get("id")
            for guide in module.get("supportingGuides", [])
            if isinstance(guide.get("id"), str)
        )
        if wanted.intersection(module_artifact_ids):
            owner_modules.append(module_id)
    return owner_modules


def artifacts_by_module(artifacts: Iterable[ExpectedArtifact]) -> dict[str, list[ExpectedArtifact]]:
    grouped: dict[str, list[ExpectedArtifact]] = {}
    for artifact in artifacts:
        grouped.setdefault(artifact.module_id, []).append(artifact)
    return grouped


def decision_required_findings(decision_plan: DecisionPlan) -> list[Finding]:
    return [
        finding(
            "decision.required",
            "decision",
            str(MANIFEST_PATH),
            decision_id,
            f"Decision {decision_id} has no compatible durable answer.",
            "Pass --decision with a valid value before applying.",
        )
        for decision_id in decision_plan.unresolved_decisions
    ]


def planned_changes_for_plan(
    repo: Path,
    decision_plan: DecisionPlan,
) -> tuple[PlannedChange, ...]:
    changes: list[PlannedChange] = [
        PlannedChange(
            action="create manifest",
            path=MANIFEST_PATH,
            managed_id="manifest",
            state="definite",
        )
    ]
    seen: set[tuple[str, str, str, str, object]] = {
        ("create manifest", MANIFEST_PATH.as_posix(), "manifest", "", "")
    }
    for planned_artifact in decision_plan.artifacts:
        change = planned_change_for_artifact(repo, planned_artifact)
        if change is None:
            continue
        condition = change.condition
        condition_key = (
            condition.decision_id if condition else "",
            condition.operator if condition else "",
            condition.value if condition else "",
        )
        key = (
            change.action,
            change.path.as_posix(),
            change.managed_id,
            change.state,
            *condition_key,
        )
        if key in seen:
            continue
        seen.add(key)
        changes.append(change)
    return tuple(changes)


def planned_change_for_artifact(
    repo: Path,
    planned_artifact: PlannedArtifact,
) -> PlannedChange | None:
    artifact = planned_artifact.artifact
    target = repo / artifact.path
    if not planned_artifact.present:
        if not target.exists():
            return None
        return PlannedChange(
            action="remove managed content",
            path=artifact.path,
            managed_id=artifact.managed_id,
            state=planned_artifact.state,
            condition=planned_artifact.condition,
        )
    action = "create guide" if artifact.kind == "guide" else "create managed block"
    if target.exists():
        action = "refresh managed content"
    return PlannedChange(
        action=action,
        path=artifact.path,
        managed_id=artifact.managed_id,
        state=planned_artifact.state,
        condition=planned_artifact.condition,
    )


def plan_apply(
    repo: Path,
    catalog: AssetCatalog,
    profile_override: str | None,
    decision_args: list[str],
) -> tuple[AuditResult, bool, ChangePlan]:
    findings: list[Finding] = []
    existing_manifest, invalid_input = load_manifest_for_apply(repo, findings)
    if invalid_input:
        return empty_apply_result(findings, invalid_input)

    if existing_manifest is not None and existing_manifest.get("schemaVersion") != 1:
        findings.append(
            finding(
                "manifest.migration-required",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "Setup Manifest schemaVersion is not supported by apply.",
                "Run a manifest migration before applying managed content.",
            )
        )
        return empty_apply_result(findings, False)

    profile_id = profile_override or (
        existing_manifest.get("profile") if isinstance(existing_manifest, dict) else None
    )
    if not isinstance(profile_id, str) or profile_id not in catalog.profiles:
        findings.append(
            finding(
                "profile.unknown",
                "error",
                str(MANIFEST_PATH),
                str(profile_id),
                "Apply requires one bundled profile.",
                "Pass --profile with a bundled profile id.",
            )
        )
        return empty_apply_result(findings, False)

    cli_decisions, parse_findings = parse_decision_args(decision_args, catalog)
    findings.extend(parse_findings)
    decision_plan = resolve_decision_plan(
        catalog,
        profile_id,
        existing_manifest,
        cli_decisions,
    )
    preview_changes = planned_changes_for_plan(repo, decision_plan)
    if parse_findings:
        return empty_apply_result(
            findings,
            True,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
        )

    findings.extend(decision_required_findings(decision_plan))
    if findings:
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
        )

    decisions = decision_plan.resolved_decisions
    ordered_modules = list(decision_plan.active_modules)
    expected_artifacts = expected_artifacts_for_plan(decision_plan)

    expected_by_id = {artifact.managed_id: artifact for artifact in expected_artifacts}
    current_files = load_current_files(repo, expected_artifacts, existing_manifest, findings)
    if findings:
        return empty_apply_result(findings, False)

    adoption_findings = require_adoption_decisions(
        current_files=current_files,
        expected_artifacts=expected_artifacts,
        decisions=decisions,
    )
    findings.extend(adoption_findings)
    if adoption_findings:
        return empty_apply_result(findings, False)

    changed_contents: dict[Path, str | None] = {}
    for relative_path, artifacts in artifacts_by_path(expected_artifacts).items():
        current = current_files.get(relative_path, "")
        changed_contents[relative_path] = render_expected_path(
            current=current,
            artifacts=artifacts,
        )

    remove_obsolete_artifacts(
        existing_manifest=existing_manifest,
        expected_by_id=expected_by_id,
        current_files=current_files,
        changed_contents=changed_contents,
    )

    manifest = build_manifest(profile_id, ordered_modules, expected_artifacts, decisions, existing_manifest)
    changed_contents[MANIFEST_PATH] = json.dumps(manifest, indent=2, sort_keys=False) + "\n"
    changes = [
        FileChange(path=path, content=content)
        for path, content in sorted(changed_contents.items(), key=lambda item: item[0].as_posix())
        if current_file_bytes(repo, path) != content_bytes(content)
    ]
    plan = ChangePlan(changes=changes, manifest=manifest)
    validation_findings = validate_change_plan(repo, plan, expected_artifacts)
    findings.extend(validation_findings)
    return AuditResult(sorted_findings(findings)), False, plan


def empty_apply_result(
    findings: list[Finding],
    invalid_input: bool,
    selection: PreviewSelection | None = None,
    planned_changes: tuple[PlannedChange, ...] = (),
) -> tuple[AuditResult, bool, ChangePlan]:
    return (
        AuditResult(
            sorted_findings(findings),
            selection=selection,
            planned_changes=planned_changes,
        ),
        invalid_input,
        ChangePlan([], {}),
    )


def load_manifest_for_apply(
    repo: Path,
    findings: list[Finding],
) -> tuple[dict | None, bool]:
    manifest_path = repo / MANIFEST_PATH
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except FileNotFoundError:
        return None, False
    except json.JSONDecodeError as error:
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                f"Setup Manifest is not valid JSON: {error.msg}.",
                "Fix the JSON syntax before applying managed content.",
            )
        )
        return None, True
    if not isinstance(manifest, dict):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "Setup Manifest must be a JSON object.",
                "Replace it with the versioned manifest object.",
            )
        )
        return None, True
    return manifest, False


def parse_decision_args(
    decision_args: list[str],
    catalog: AssetCatalog,
) -> tuple[dict[str, object], list[Finding]]:
    decisions: dict[str, object] = {}
    findings: list[Finding] = []
    for raw in decision_args:
        if "=" not in raw:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "decision",
                    f"Decision argument {raw!r} must use ID=VALUE.",
                    "Pass decisions as --decision decision.id=value.",
                )
            )
            continue
        decision_id, raw_value = raw.split("=", 1)
        decision_id = decision_id.strip()
        if not decision_id:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "decision",
                    "Decision argument is missing an id.",
                    "Pass decisions as --decision decision.id=value.",
                )
            )
            continue
        decisions[decision_id] = parse_decision_value(decision_id, raw_value, catalog)
    return decisions, findings


def parse_decision_value(decision_id: str, raw_value: str, catalog: AssetCatalog) -> object:
    value = raw_value.strip()
    if decision_id.startswith("adoption."):
        return parse_bool_value(value)
    contract = catalog.decisions.get(decision_id)
    if contract and contract.get("type") == "boolean":
        return parse_bool_value(value)
    return value


def parse_bool_value(value: str) -> object:
    normalized = value.lower()
    if normalized in {"1", "true", "yes", "y"}:
        return True
    if normalized in {"0", "false", "no", "n"}:
        return False
    return value


def resolve_decisions(
    catalog: AssetCatalog,
    ordered_modules: list[str],
    existing_manifest: dict | None,
    cli_decisions: dict[str, object],
    findings: list[Finding],
) -> dict[str, dict]:
    today = date.today().isoformat()
    existing = existing_manifest.get("decisions", {}) if isinstance(existing_manifest, dict) else {}
    if not isinstance(existing, dict):
        existing = {}
    resolved: dict[str, dict] = {}
    for decision_id in required_decisions(catalog, ordered_modules):
        contract = catalog.decisions[decision_id]
        value = cli_decisions.get(decision_id)
        if value is not None:
            if not is_decision_value_valid(contract, value):
                findings.append(invalid_decision_finding(decision_id))
                continue
            resolved[decision_id] = {"value": value, "confirmedAt": today}
            continue
        existing_decision = existing.get(decision_id)
        if isinstance(existing_decision, dict) and "value" in existing_decision:
            if is_decision_value_valid(contract, existing_decision["value"]):
                confirmed_at = existing_decision.get("confirmedAt")
                if not isinstance(confirmed_at, str) or not confirmed_at:
                    confirmed_at = today
                resolved[decision_id] = {
                    "value": existing_decision["value"],
                    "confirmedAt": confirmed_at,
                }
                continue
        findings.append(
            finding(
                "decision.required",
                "decision",
                str(MANIFEST_PATH),
                decision_id,
                f"Decision {decision_id} has no compatible durable answer.",
                "Pass --decision with a valid value before applying.",
            )
        )

    for decision_id, existing_decision in sorted(existing.items()):
        if decision_id.startswith("adoption.") and isinstance(existing_decision, dict):
            resolved[decision_id] = existing_decision
    for decision_id, value in sorted(cli_decisions.items()):
        if decision_id.startswith("adoption."):
            resolved[decision_id] = {"value": value, "confirmedAt": today}
    return resolved


def is_decision_value_valid(decision_contract: dict, value: object) -> bool:
    decision_type = decision_contract.get("type")
    if decision_type == "boolean":
        return isinstance(value, bool)
    if decision_type == "string":
        return isinstance(value, str) and bool(value.strip())
    if decision_type == "enum":
        return value in decision_contract.get("values", [])
    return True


def invalid_decision_finding(decision_id: str) -> Finding:
    return finding(
        "decision.required",
        "decision",
        str(MANIFEST_PATH),
        decision_id,
        f"Decision {decision_id} has an incompatible value.",
        "Pass --decision with a value accepted by the decision contract.",
    )


def load_current_files(
    repo: Path,
    expected_artifacts: list[ExpectedArtifact],
    existing_manifest: dict | None,
    findings: list[Finding],
) -> dict[Path, str]:
    paths = {artifact.path for artifact in expected_artifacts}
    paths.update(manifest_artifact_paths(existing_manifest))
    current: dict[Path, str] = {}
    for relative_path in sorted(paths, key=lambda item: item.as_posix()):
        path = repo / relative_path
        if not path.exists():
            continue
        if path.is_dir():
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(relative_path),
                    "document",
                    "Managed destination is a directory.",
                    "Move the directory before applying managed content.",
                )
            )
            continue
        try:
            content = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(relative_path),
                    "document",
                    "Managed destination is not UTF-8 text.",
                    "Rewrite the file as UTF-8 Markdown before applying.",
                )
            )
            continue
        _, marker_findings = parse_managed_blocks(relative_path, content)
        findings.extend(marker_findings)
        current[relative_path] = content
    return current


def manifest_artifact_paths(existing_manifest: dict | None) -> set[Path]:
    if not isinstance(existing_manifest, dict):
        return set()
    paths: set[Path] = set()
    for artifact in existing_manifest.get("managedArtifacts", []):
        if isinstance(artifact, dict) and isinstance(artifact.get("path"), str):
            paths.add(Path(artifact["path"]))
    return paths


def require_adoption_decisions(
    current_files: dict[Path, str],
    expected_artifacts: list[ExpectedArtifact],
    decisions: dict[str, dict],
) -> list[Finding]:
    findings: list[Finding] = []
    for artifact in expected_artifacts:
        if artifact.kind != "guide":
            continue
        content = current_files.get(artifact.path)
        if content is None or not content.strip():
            continue
        blocks, marker_findings = parse_managed_blocks(artifact.path, content)
        if marker_findings or artifact.managed_id in blocks:
            continue
        adoption_id = f"adoption.{artifact.managed_id}"
        decision = decisions.get(adoption_id)
        if not isinstance(decision, dict) or decision.get("value") is not True:
            findings.append(
                finding(
                    "decision.required",
                    "decision",
                    artifact.path.as_posix(),
                    adoption_id,
                    "Existing unmarked content needs an adoption decision before setup can own it.",
                    f"Pass --decision {adoption_id}=true after reviewing the file.",
                )
            )
    return findings


def artifacts_by_path(
    expected_artifacts: list[ExpectedArtifact],
) -> dict[Path, list[ExpectedArtifact]]:
    grouped: dict[Path, list[ExpectedArtifact]] = {}
    for artifact in expected_artifacts:
        grouped.setdefault(artifact.path, []).append(artifact)
    return grouped


def render_expected_path(current: str, artifacts: list[ExpectedArtifact]) -> str:
    if all(artifact.kind == "guide" for artifact in artifacts):
        return render_guide_path(current, artifacts)
    return render_shared_path(current, artifacts)


def render_guide_path(current: str, artifacts: list[ExpectedArtifact]) -> str:
    content = current
    for artifact in artifacts:
        replacement = managed_block(artifact.managed_id, artifact.version, artifact.content)
        spans, _ = parse_managed_block_spans(artifact.path, content)
        span = spans.get(artifact.managed_id)
        if span is None:
            if content.strip() and not spans:
                content = replacement
            else:
                content = append_block(content, replacement)
        else:
            content = content[: span.start] + replacement + content[span.end :]
    return content


def render_shared_path(current: str, artifacts: list[ExpectedArtifact]) -> str:
    content = current
    for artifact in artifacts:
        replacement = managed_block(artifact.managed_id, artifact.version, artifact.content)
        spans, _ = parse_managed_block_spans(artifact.path, content)
        span = spans.get(artifact.managed_id)
        if span is None:
            content = append_block(content, replacement)
        else:
            content = content[: span.start] + replacement + content[span.end :]
    return content


def append_block(content: str, block: str) -> str:
    if not content:
        return block
    separator = "" if content.endswith("\n") else "\n"
    return content + separator + "\n" + block


def remove_obsolete_artifacts(
    existing_manifest: dict | None,
    expected_by_id: dict[str, ExpectedArtifact],
    current_files: dict[Path, str],
    changed_contents: dict[Path, str | None],
) -> None:
    if not isinstance(existing_manifest, dict):
        return
    for artifact in existing_manifest.get("managedArtifacts", []):
        if not isinstance(artifact, dict):
            continue
        managed_id = artifact.get("id")
        path_value = artifact.get("path")
        if not isinstance(managed_id, str) or not isinstance(path_value, str):
            continue
        if managed_id in expected_by_id:
            continue
        relative_path = Path(path_value)
        content = changed_contents.get(relative_path, current_files.get(relative_path))
        if content is None:
            continue
        spans, marker_findings = parse_managed_block_spans(relative_path, content)
        if marker_findings:
            continue
        span = spans.get(managed_id)
        if span is None:
            continue
        remaining = content[: span.start] + content[span.end :]
        if artifact.get("kind") == "guide" and not remaining.strip():
            changed_contents[relative_path] = None
        else:
            changed_contents[relative_path] = remaining


def build_manifest(
    profile_id: str,
    ordered_modules: list[str],
    expected_artifacts: list[ExpectedArtifact],
    decisions: dict[str, dict],
    existing_manifest: dict | None,
) -> dict:
    local_skills = []
    if isinstance(existing_manifest, dict) and isinstance(existing_manifest.get("localSkills"), list):
        local_skills = existing_manifest["localSkills"]
    return {
        "schemaVersion": 1,
        "generator": {"skill": "setup-context-driven", "version": 1},
        "profile": profile_id,
        "modules": ordered_modules,
        "decisions": decisions,
        "managedArtifacts": [
            {
                "id": artifact.managed_id,
                "path": artifact.path.as_posix(),
                "kind": artifact.kind,
                "module": artifact.module_id,
                "template": artifact.template_id,
                "version": artifact.version,
                "digest": artifact.digest,
            }
            for artifact in expected_artifacts
        ],
        "localSkills": local_skills,
    }


def validate_change_plan(
    repo: Path,
    plan: ChangePlan,
    expected_artifacts: list[ExpectedArtifact],
) -> list[Finding]:
    findings: list[Finding] = []
    planned = {change.path: change.content for change in plan.changes}
    for artifact in expected_artifacts:
        content = planned.get(artifact.path)
        if content is None:
            path = repo / artifact.path
            if path.exists() and path.is_file():
                content = path.read_text(encoding="utf-8")
        if content is None:
            findings.append(
                finding(
                    "managed.block.missing",
                    "error",
                    artifact.path.as_posix(),
                    artifact.managed_id,
                    "Change plan does not produce the expected managed block.",
                    "Rebuild the apply plan before writing.",
                )
            )
            continue
        blocks, marker_findings = parse_managed_blocks(artifact.path, content)
        findings.extend(marker_findings)
        block = blocks.get(artifact.managed_id)
        if block is None:
            findings.append(
                finding(
                    "managed.block.missing",
                    "error",
                    artifact.path.as_posix(),
                    artifact.managed_id,
                    "Change plan omits an expected managed block.",
                    "Rebuild the apply plan before writing.",
                )
            )
        elif managed_digest(block.body) != artifact.digest:
            findings.append(
                finding(
                    "managed.content.modified",
                    "warning",
                    artifact.path.as_posix(),
                    artifact.managed_id,
                    "Change plan does not match the bundled template digest.",
                    "Rebuild the apply plan from bundled assets.",
                )
            )
    return findings


def apply_change_plan(repo: Path, plan: ChangePlan) -> None:
    temp_paths: list[Path] = []
    originals: dict[Path, bytes | None] = {}
    created_dirs: list[Path] = []
    try:
        for change in plan.changes:
            target = repo / change.path
            originals[target] = target.read_bytes() if target.exists() and target.is_file() else None
            if change.content is None:
                continue
            ensure_parent_dir(target.parent, created_dirs)
            temp_path = target.with_name(f".{target.name}.setup-context.tmp")
            temp_path.write_text(change.content, encoding="utf-8")
            temp_paths.append(temp_path)

        for change in plan.changes:
            target = repo / change.path
            if change.content is None:
                if target.exists() and target.is_file():
                    target.unlink()
                continue
            temp_path = target.with_name(f".{target.name}.setup-context.tmp")
            temp_path.replace(target)
    except OSError:
        for target, original in originals.items():
            if original is None:
                if target.exists() and target.is_file():
                    target.unlink()
            else:
                ensure_parent_dir(target.parent, created_dirs)
                target.write_bytes(original)
        raise
    finally:
        for temp_path in temp_paths:
            if temp_path.exists():
                temp_path.unlink()
        for directory in reversed(created_dirs):
            try:
                directory.rmdir()
            except OSError:
                pass


def ensure_parent_dir(directory: Path, created_dirs: list[Path]) -> None:
    missing: list[Path] = []
    current = directory
    while not current.exists():
        missing.append(current)
        current = current.parent
    for item in reversed(missing):
        item.mkdir()
        created_dirs.append(item)


def current_file_bytes(repo: Path, path: Path) -> bytes | None:
    target = repo / path
    if not target.exists() or not target.is_file():
        return None
    return target.read_bytes()


def content_bytes(content: str | None) -> bytes | None:
    if content is None:
        return None
    return content.encode("utf-8")


def sync_setup_snapshots(
    skill_root: Path,
    source_dir: Path,
    check: bool,
) -> tuple[AuditResult, bool]:
    findings: list[Finding] = []
    current_snapshots, invalid_input = load_current_setup_snapshots(skill_root, findings)
    if invalid_input:
        return AuditResult(sorted_findings(findings)), True

    planned: dict[Path, str] = {}
    changed_setup_ids: list[str] = []
    for setup_id, current_snapshot in sorted(current_snapshots.items()):
        snapshot, snapshot_findings, snapshot_invalid = build_source_setup_snapshot(
            setup_id=setup_id,
            source_dir=source_dir,
            current_snapshot=current_snapshot,
        )
        findings.extend(snapshot_findings)
        invalid_input = invalid_input or snapshot_invalid
        if snapshot is None:
            continue
        if snapshot == current_snapshot:
            continue
        changed_setup_ids.append(setup_id)
        if check:
            findings.append(setup_snapshot_drift_finding(setup_id))
        else:
            target = skill_root / "assets" / "setups" / f"{setup_id}.json"
            planned[target] = snapshot_json(snapshot)

    if invalid_input or check:
        return AuditResult(sorted_findings(findings)), invalid_input

    try:
        write_atomic_text_changes(planned)
    except OSError as error:
        findings.append(
            finding(
                "skills.setup-snapshot.drift",
                "error",
                "assets/setups",
                "setups",
                f"Could not update setup snapshots atomically: {error}.",
                "Fix filesystem permissions and rerun sync-setups.",
            )
        )
        return AuditResult(sorted_findings(findings)), False

    for setup_id in changed_setup_ids:
        findings.append(
            finding(
                "skills.setup-snapshot.updated",
                "info",
                f"assets/setups/{setup_id}.json",
                f"setup.{setup_id}",
                "Setup snapshot was synchronized from the canonical source.",
                "Review the snapshot diff and run asset validation.",
            )
        )
    return AuditResult(sorted_findings(findings)), False


def load_current_setup_snapshots(
    skill_root: Path,
    findings: list[Finding],
) -> tuple[dict[str, dict], bool]:
    setups_root = skill_root / "assets" / "setups"
    snapshots: dict[str, dict] = {}
    invalid_input = False
    for path in sorted(setups_root.glob("*.json")):
        try:
            snapshot = json.loads(path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as error:
            findings.append(
                finding(
                    "skills.setup-snapshot.drift",
                    "error",
                    path.relative_to(skill_root).as_posix(),
                    path.stem,
                    f"Bundled setup snapshot is not valid JSON: {error.msg}.",
                    "Fix the bundled setup snapshot before synchronizing.",
                )
            )
            invalid_input = True
            continue
        if not isinstance(snapshot, dict) or snapshot.get("id") != path.stem:
            findings.append(
                finding(
                    "skills.setup-snapshot.drift",
                    "error",
                    path.relative_to(skill_root).as_posix(),
                    path.stem,
                    "Bundled setup snapshot id does not match its filename.",
                    "Fix the bundled setup snapshot before synchronizing.",
                )
            )
            invalid_input = True
            continue
        snapshots[path.stem] = snapshot
    return snapshots, invalid_input


def build_source_setup_snapshot(
    setup_id: str,
    source_dir: Path,
    current_snapshot: dict,
) -> tuple[dict | None, list[Finding], bool]:
    findings: list[Finding] = []
    source_doc, source_path, invalid_input = load_source_setup_doc(setup_id, source_dir, findings)
    if invalid_input or source_doc is None or source_path is None:
        return None, findings, invalid_input

    current_by_path = {
        skill.get("path"): skill
        for skill in current_snapshot.get("skills", [])
        if isinstance(skill, dict) and isinstance(skill.get("path"), str)
    }
    current_by_name = {
        skill.get("name"): skill
        for skill in current_snapshot.get("skills", [])
        if isinstance(skill, dict) and isinstance(skill.get("name"), str)
    }
    skills: list[dict] = []
    seen_names: set[str] = set()
    seen_paths: set[str] = set()
    for raw_skill in source_doc.get("skills", []):
        skill, skill_findings = normalize_source_skill(
            raw_skill=raw_skill,
            source_dir=source_dir,
            current_by_path=current_by_path,
            current_by_name=current_by_name,
            source_path=source_path,
        )
        findings.extend(skill_findings)
        if skill is None:
            invalid_input = True
            continue
        if skill["name"] in seen_names:
            findings.append(
                setup_snapshot_invalid_finding(
                    source_path,
                    setup_id,
                    f"Canonical setup contains duplicate skill name {skill['name']}.",
                )
            )
            invalid_input = True
        if skill["path"] in seen_paths:
            findings.append(
                setup_snapshot_invalid_finding(
                    source_path,
                    setup_id,
                    f"Canonical setup contains duplicate skill path {skill['path']}.",
                )
            )
            invalid_input = True
        seen_names.add(skill["name"])
        seen_paths.add(skill["path"])
        skills.append(skill)

    if not isinstance(source_doc.get("skills"), list) or not skills:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup must contain a non-empty skills list.",
            )
        )
        invalid_input = True

    source_metadata = source_doc.get("source")
    if not isinstance(source_metadata, dict):
        source_metadata = current_snapshot.get("source")
    if not isinstance(source_metadata, dict):
        source_metadata = {
            "name": setup_id,
            "type": "canonical-skill-setup",
            "revision": "unknown",
        }

    snapshot = {
        "schemaVersion": "setup-context-driven/setup-snapshot-v1",
        "id": setup_id,
        "version": current_snapshot.get("version", 1),
        "source": source_metadata,
        "digest": setup_paths_digest([skill["path"] for skill in skills]),
        "skills": skills,
    }
    return snapshot, findings, invalid_input


def load_source_setup_doc(
    setup_id: str,
    source_dir: Path,
    findings: list[Finding],
) -> tuple[dict | None, Path | None, bool]:
    json_path = source_dir / f"{setup_id}.json"
    if json_path.is_file():
        try:
            doc = json.loads(json_path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as error:
            findings.append(
                setup_snapshot_invalid_finding(
                    json_path,
                    setup_id,
                    f"Canonical setup JSON is invalid: {error.msg}.",
                )
            )
            return None, json_path, True
        if not isinstance(doc, dict):
            findings.append(
                setup_snapshot_invalid_finding(
                    json_path,
                    setup_id,
                    "Canonical setup JSON must be an object.",
                )
            )
            return None, json_path, True
        return doc, json_path, False

    for suffix in (".txt", ".md"):
        path = source_dir / f"{setup_id}{suffix}"
        if path.is_file():
            skills = [
                {"path": line.strip()}
                for line in path.read_text(encoding="utf-8").splitlines()
                if line.strip() and not line.lstrip().startswith("#")
            ]
            return {"skills": skills}, path, False

    findings.append(
        setup_snapshot_invalid_finding(
            source_dir / f"{setup_id}.json",
            setup_id,
            "Canonical setup source file is missing.",
        )
    )
    return None, None, True


def normalize_source_skill(
    raw_skill: object,
    source_dir: Path,
    current_by_path: dict[str, dict],
    current_by_name: dict[str, dict],
    source_path: Path,
) -> tuple[dict | None, list[Finding]]:
    findings: list[Finding] = []
    if isinstance(raw_skill, str):
        raw_data: dict[str, object] = {"path": raw_skill}
    elif isinstance(raw_skill, dict):
        raw_data = raw_skill
    else:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                "setup",
                "Each canonical setup skill must be a string or object.",
            )
        )
        return None, findings

    raw_path = raw_data.get("path") or raw_data.get("skillPath") or raw_data.get("name")
    if not isinstance(raw_path, str) or not raw_path.strip():
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                "setup",
                "Canonical setup skill is missing a path.",
            )
        )
        return None, findings
    normalized_path = normalize_skill_path(raw_path)
    if normalized_path is None:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                "setup",
                f"Canonical setup skill path {raw_path!r} is not portable.",
            )
        )
        return None, findings

    name = raw_data.get("name")
    if not isinstance(name, str) or not name.strip():
        name = infer_skill_name(normalized_path)
    current_skill = current_by_path.get(normalized_path) or current_by_name.get(str(name), {})
    source = raw_data.get("source")
    if not isinstance(source, dict):
        source = lock_style_source(raw_data)
    if not source and isinstance(current_skill.get("source"), dict):
        source = current_skill["source"]
    digest = raw_data.get("contentDigest") or raw_data.get("computedHash")
    if not isinstance(digest, str) or not digest:
        source_skill_path = resolve_canonical_skill_file(source_dir, raw_path)
        if source_skill_path is not None and source_skill_path.is_file():
            digest = managed_digest(source_skill_path.read_text(encoding="utf-8"))
        elif isinstance(current_skill.get("contentDigest"), str):
            digest = current_skill["contentDigest"]
    if not isinstance(digest, str) or not digest:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} is missing a content digest.",
            )
        )
        return None, findings

    return {
        "name": str(name).strip(),
        "path": normalized_path,
        "source": source,
        "contentDigest": digest,
    }, findings


def normalize_skill_path(raw_path: str) -> str | None:
    value = raw_path.strip().replace("\\", "/")
    if not value:
        return None
    candidate = Path(value)
    if candidate.is_absolute() or ".." in candidate.parts:
        return None
    if "/" not in value:
        return f".agents/skills/{value}/SKILL.md"
    return candidate.as_posix()


def resolve_canonical_skill_file(source_dir: Path, raw_path: str) -> Path | None:
    value = raw_path.strip().replace("\\", "/")
    candidate = Path(value)
    if not value or candidate.is_absolute() or ".." in candidate.parts:
        return None
    if candidate.parts and candidate.parts[0] == "skills":
        source_path = source_dir.parent / candidate
    else:
        source_path = source_dir / candidate
    if source_path.is_dir():
        return source_path / "SKILL.md"
    return source_path


def infer_skill_name(path: str) -> str:
    candidate = Path(path)
    if candidate.name == "SKILL.md" and candidate.parent.name:
        return candidate.parent.name
    return candidate.stem


def lock_style_source(raw_data: dict[str, object]) -> dict:
    source_name = raw_data.get("source")
    source_type = raw_data.get("sourceType")
    if not isinstance(source_name, str) and not isinstance(source_type, str):
        return {}
    source: dict[str, str] = {}
    if isinstance(source_type, str):
        source["type"] = source_type
    if isinstance(source_name, str):
        source["name"] = source_name
    if isinstance(raw_data.get("skillPath"), str):
        source["path"] = raw_data["skillPath"]
    if isinstance(raw_data.get("ref"), str):
        source["ref"] = raw_data["ref"]
    return source


def setup_paths_digest(paths: list[str]) -> str:
    payload = "\n".join(paths) + "\n"
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def snapshot_json(snapshot: dict) -> str:
    return json.dumps(snapshot, indent=2, sort_keys=False) + "\n"


def write_atomic_text_changes(changes: dict[Path, str]) -> None:
    temp_paths: list[Path] = []
    originals: dict[Path, bytes | None] = {}
    try:
        for target, content in sorted(changes.items(), key=lambda item: item[0].as_posix()):
            originals[target] = target.read_bytes() if target.exists() and target.is_file() else None
            temp_path = target.with_name(f".{target.name}.setup-context.tmp")
            temp_path.write_text(content, encoding="utf-8")
            temp_paths.append(temp_path)
        for target in sorted(changes, key=lambda item: item.as_posix()):
            temp_path = target.with_name(f".{target.name}.setup-context.tmp")
            temp_path.replace(target)
    except OSError:
        for target, original in originals.items():
            if original is not None:
                target.write_bytes(original)
            elif target.exists() and target.is_file():
                target.unlink()
        raise
    finally:
        for temp_path in temp_paths:
            if temp_path.exists():
                temp_path.unlink()


def setup_snapshot_invalid_finding(path: Path, setup_id: str, message: str) -> Finding:
    return finding(
        "skills.setup-snapshot.drift",
        "error",
        path.as_posix(),
        f"setup.{setup_id}",
        message,
        "Fix the canonical setup source before synchronizing snapshots.",
    )


def setup_snapshot_drift_finding(setup_id: str) -> Finding:
    return finding(
        "skills.setup-snapshot.drift",
        "error",
        f"assets/setups/{setup_id}.json",
        f"setup.{setup_id}",
        "Bundled setup snapshot differs from the canonical source.",
        "Run sync-setups without --check to refresh the bundled snapshot.",
    )


def load_manifest(repo: Path, findings: list[Finding]) -> tuple[dict | None, bool]:
    manifest_path = repo / MANIFEST_PATH
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except FileNotFoundError:
        findings.append(
            finding(
                "manifest.missing",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "Setup Manifest is missing.",
                "Run apply after selecting a profile.",
            )
        )
        return None, False
    except json.JSONDecodeError as error:
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                f"Setup Manifest is not valid JSON: {error.msg}.",
                "Fix the JSON syntax before running audit again.",
            )
        )
        return None, True

    if not isinstance(manifest, dict):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "Setup Manifest must be a JSON object.",
                "Replace it with the versioned manifest object.",
            )
        )
        return None, True
    return manifest, False


def validate_manifest_shape(
    manifest: dict,
    profile_id: str,
    ordered_modules: list[str],
    catalog: AssetCatalog,
    findings: list[Finding],
) -> None:
    if manifest.get("schemaVersion") != 1:
        findings.append(
            finding(
                "manifest.migration-required",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "Setup Manifest schemaVersion is not 1.",
                "Migrate the manifest before auditing generated guidance.",
            )
        )

    manifest_modules = manifest.get("modules")
    if manifest_modules != ordered_modules:
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                f"profile.{profile_id}",
                "Manifest modules do not match the selected profile order.",
                "Refresh the manifest from the selected profile.",
            )
        )

    decisions = manifest.get("decisions", {})
    if not isinstance(decisions, dict):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "decisions",
                "Manifest decisions must be a JSON object.",
                "Store decisions by stable decision identifier.",
            )
        )
        decisions = {}

    for decision_id in required_decisions(catalog, ordered_modules):
        decision = decisions.get(decision_id)
        if not isinstance(decision, dict) or "value" not in decision:
            findings.append(
                finding(
                    "decision.required",
                    "decision",
                    str(MANIFEST_PATH),
                    decision_id,
                    f"Decision {decision_id} has no durable answer.",
                    "Confirm and store the decision in the Setup Manifest.",
                )
            )
            continue
        validate_decision_value(catalog.decisions[decision_id], decision["value"], findings)


def validate_decision_value(
    decision_contract: dict,
    value: object,
    findings: list[Finding],
) -> None:
    decision_id = decision_contract["id"]
    decision_type = decision_contract.get("type")
    valid = True
    if decision_type == "boolean":
        valid = isinstance(value, bool)
    elif decision_type == "string":
        valid = isinstance(value, str) and bool(value.strip())
    elif decision_type == "enum":
        valid = value in decision_contract.get("values", [])

    if not valid:
        findings.append(
            finding(
                "decision.required",
                "decision",
                str(MANIFEST_PATH),
                decision_id,
                f"Decision {decision_id} has an invalid value.",
                "Confirm a valid value and update the Setup Manifest.",
            )
        )


def validate_manifest_artifacts(
    manifest: dict,
    expected_artifacts: list[ExpectedArtifact],
    findings: list[Finding],
) -> None:
    expected_by_id = {artifact.managed_id: artifact for artifact in expected_artifacts}
    seen: set[str] = set()
    managed_artifacts = manifest.get("managedArtifacts", [])
    if not isinstance(managed_artifacts, list):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "managedArtifacts",
                "Manifest managedArtifacts must be a list.",
                "Refresh the managed artifact inventory.",
            )
        )
        return

    for artifact in managed_artifacts:
        if not isinstance(artifact, dict):
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "managedArtifacts",
                    "Each managed artifact entry must be an object.",
                    "Refresh the managed artifact inventory.",
                )
            )
            continue
        managed_id = artifact.get("id")
        if not isinstance(managed_id, str) or not managed_id:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "managedArtifacts",
                    "Managed artifact entry is missing id.",
                    "Refresh the managed artifact inventory.",
                )
            )
            continue
        if managed_id in seen:
            findings.append(
                finding(
                    "managed.block.duplicate",
                    "error",
                    str(MANIFEST_PATH),
                    managed_id,
                    "Managed artifact appears more than once in the manifest.",
                    "Keep one inventory entry per managed identifier.",
                )
            )
        seen.add(managed_id)

        expected = expected_by_id.get(managed_id)
        if expected is None:
            findings.append(
                finding(
                    "docs.reference.broken",
                    "error",
                    str(MANIFEST_PATH),
                    managed_id,
                    "Managed artifact references an unknown generated asset.",
                    "Remove the stale inventory entry or update the selected profile.",
                )
            )
            continue
        if artifact.get("template") != expected.template_id or artifact.get("version") != expected.version:
            findings.append(stale_template_finding(str(MANIFEST_PATH), managed_id))


def validate_documents(
    repo: Path,
    expected_artifacts: list[ExpectedArtifact],
    findings: list[Finding],
) -> None:
    artifacts_by_path: dict[Path, list[ExpectedArtifact]] = {}
    for artifact in expected_artifacts:
        artifacts_by_path.setdefault(artifact.path, []).append(artifact)

    for relative_path, artifacts in sorted(artifacts_by_path.items(), key=lambda item: str(item[0])):
        path = repo / relative_path
        if not path.exists():
            for artifact in artifacts:
                if artifact.kind == "guide":
                    findings.append(
                        finding(
                            "docs.guide.missing",
                            "error",
                            str(relative_path),
                            artifact.managed_id,
                            "Supporting guide is missing.",
                            "Restore the setup-owned guide from its template.",
                        )
                    )
                else:
                    findings.append(
                        finding(
                            "managed.block.missing",
                            "error",
                            str(relative_path),
                            artifact.managed_id,
                            "Managed root block is missing.",
                            "Restore the setup-owned root block.",
                        )
                    )
            continue

        try:
            content = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(relative_path),
                    "document",
                    "Managed document is not UTF-8 text.",
                    "Rewrite the document as UTF-8 Markdown.",
                )
            )
            continue

        blocks, marker_findings = parse_managed_blocks(relative_path, content)
        findings.extend(marker_findings)
        for artifact in artifacts:
            block = blocks.get(artifact.managed_id)
            if block is None:
                code = "docs.guide.missing" if artifact.kind == "guide" else "managed.block.missing"
                message = (
                    "Supporting guide managed block is missing."
                    if artifact.kind == "guide"
                    else "Managed root block is missing."
                )
                findings.append(
                    finding(
                        code,
                        "error",
                        str(relative_path),
                        artifact.managed_id,
                        message,
                        "Restore the setup-owned block from its template.",
                    )
                )
                continue
            if block.version != artifact.version:
                findings.append(stale_template_finding(str(relative_path), artifact.managed_id))
            if managed_digest(block.body) != artifact.digest:
                findings.append(
                    finding(
                        "managed.content.modified",
                        "warning",
                        str(relative_path),
                        artifact.managed_id,
                        "Managed content digest differs from the bundled template.",
                        "Review the local edit or refresh the setup-owned block.",
                    )
                )
            if contains_non_english_marker(block.body):
                findings.append(
                    finding(
                        "docs.language.non-english",
                        "warning",
                        str(relative_path),
                        artifact.managed_id,
                        "Generated content appears to use a non-English phrase.",
                        "Rewrite setup-generated content in English.",
                    )
                )
            findings.extend(validate_internal_references(repo, relative_path, block.body, artifact.managed_id))


def validate_secondbrain_documents(
    repo: Path,
    ordered_modules: list[str],
    findings: list[Finding],
) -> None:
    if "secondbrain" not in ordered_modules:
        return

    root_blocks, root_findings = read_document_blocks(repo, ROOT_INSTRUCTIONS_PATH)
    guide_blocks, guide_findings = read_document_blocks(repo, Path("docs/agents/secondbrain.md"))
    if root_findings or guide_findings:
        return

    pointer = root_blocks.get("root.secondbrain")
    if pointer is None:
        findings.append(
            finding(
                "secondbrain.pointer.missing",
                "error",
                str(ROOT_INSTRUCTIONS_PATH),
                "root.secondbrain",
                "Secondbrain is enabled but the root pointer is missing.",
                "Restore the compact Secondbrain root pointer.",
            )
        )
    elif "docs/agents/secondbrain.md" not in pointer.body:
        findings.append(
            finding(
                "secondbrain.pointer.missing",
                "error",
                str(ROOT_INSTRUCTIONS_PATH),
                "root.secondbrain",
                "Secondbrain root pointer does not reference its supporting guide.",
                "Point the root block to docs/agents/secondbrain.md.",
            )
        )

    guide = guide_blocks.get("guide.secondbrain")
    if guide is None:
        findings.append(
            finding(
                "secondbrain.guide.missing",
                "error",
                "docs/agents/secondbrain.md",
                "guide.secondbrain",
                "Secondbrain is enabled but the read-only supporting guide is missing.",
                "Restore docs/agents/secondbrain.md from the bundled template.",
            )
        )
        return
    for phrase in SECONDBRAIN_REQUIRED_GUIDE_PHRASES:
        if phrase not in guide.body:
            findings.append(
                finding(
                    "secondbrain.safety-rule.missing",
                    "error",
                    "docs/agents/secondbrain.md",
                    "guide.secondbrain",
                    f"Secondbrain guide is missing required guidance: {phrase}.",
                    "Restore the complete read-only Secondbrain guide from the bundled template.",
                )
            )


def read_document_blocks(
    repo: Path,
    relative_path: Path,
) -> tuple[dict[str, ManagedBlock], list[Finding]]:
    try:
        content = (repo / relative_path).read_text(encoding="utf-8")
    except FileNotFoundError:
        return {}, []
    except UnicodeDecodeError:
        return {}, [
            finding(
                "manifest.invalid",
                "error",
                str(relative_path),
                "document",
                "Managed document is not UTF-8 text.",
                "Rewrite the document as UTF-8 Markdown.",
            )
        ]
    return parse_managed_blocks(relative_path, content)


def parse_managed_blocks(
    relative_path: Path,
    content: str,
) -> tuple[dict[str, ManagedBlock], list[Finding]]:
    spans, findings = parse_managed_block_spans(relative_path, content)
    return {
        managed_id: ManagedBlock(
            managed_id=span.managed_id,
            version=span.version,
            body=span.body,
        )
        for managed_id, span in spans.items()
    }, findings


def parse_managed_block_spans(
    relative_path: Path,
    content: str,
) -> tuple[dict[str, ManagedBlockSpan], list[Finding]]:
    findings: list[Finding] = []
    blocks: dict[str, ManagedBlockSpan] = {}
    open_marker: tuple[str, int, int, int] | None = None

    for marker in MARKER.finditer(content):
        text = marker.group(0)
        begin = BEGIN_MARKER.fullmatch(text)
        end = END_MARKER.fullmatch(text)
        if begin:
            managed_id = begin.group(1)
            version = int(begin.group(2))
            if open_marker is not None:
                findings.append(
                    marker_invalid_finding(
                        relative_path,
                        managed_id,
                        "Managed block markers are nested.",
                    )
                )
                continue
            open_marker = (managed_id, version, marker.start(), marker.end())
            continue
        if end:
            managed_id = end.group(1)
            if open_marker is None:
                findings.append(
                    marker_invalid_finding(
                        relative_path,
                        managed_id,
                        "Managed end marker has no matching begin marker.",
                    )
                )
                continue
            open_id, version, marker_start, body_start = open_marker
            if managed_id != open_id:
                findings.append(
                    marker_invalid_finding(
                        relative_path,
                        managed_id,
                        f"Managed end marker closes {managed_id}, expected {open_id}.",
                    )
                )
                open_marker = None
                continue
            if managed_id in blocks:
                findings.append(
                    finding(
                        "managed.block.duplicate",
                        "error",
                        str(relative_path),
                        managed_id,
                        "Managed block appears more than once in the document.",
                        "Keep one block for each managed identifier.",
                    )
                )
            block_end = marker.end()
            if content.startswith("\r\n", block_end):
                block_end += 2
            elif content.startswith("\n", block_end):
                block_end += 1
            blocks[managed_id] = ManagedBlockSpan(
                managed_id=managed_id,
                version=version,
                body=content[body_start:marker.start()],
                start=marker_start,
                end=block_end,
            )
            open_marker = None
            continue
        findings.append(
            marker_invalid_finding(
                relative_path,
                "unknown",
                "Managed marker is malformed.",
            )
        )

    if open_marker is not None:
        managed_id, _, _, _ = open_marker
        findings.append(
            marker_invalid_finding(
                relative_path,
                managed_id,
                "Managed begin marker has no matching end marker.",
            )
        )

    return blocks, findings


def validate_internal_references(
    repo: Path,
    relative_path: Path,
    content: str,
    managed_id: str,
) -> list[Finding]:
    findings: list[Finding] = []
    document_dir = (repo / relative_path).parent
    for match in MARKDOWN_LINK.finditer(content):
        target = match.group(1).strip()
        if not target or "://" in target or target.startswith(("mailto:", "#")):
            continue
        target_path = target.split("#", 1)[0]
        if not target_path:
            continue
        candidate = (document_dir / target_path).resolve(strict=False)
        try:
            candidate.relative_to(repo)
        except ValueError:
            findings.append(
                broken_reference_finding(relative_path, managed_id, target)
            )
            continue
        if not candidate.exists():
            findings.append(
                broken_reference_finding(relative_path, managed_id, target)
            )
    return findings


def expected_artifacts_for_profile(
    catalog: AssetCatalog,
    profile_id: str,
    ordered_modules: list[str] | None = None,
) -> list[ExpectedArtifact]:
    artifacts: list[ExpectedArtifact] = []
    templates_root = Path(__file__).resolve().parents[1] / "assets" / "templates"
    modules = ordered_modules or catalog.ordered_modules_by_profile[profile_id]
    for module_id in modules:
        module = catalog.modules[module_id]
        for block in module.get("rootBlocks", []):
            template_id = block["template"]
            content = template_content(templates_root, catalog, template_id)
            artifacts.append(
                ExpectedArtifact(
                    managed_id=block["id"],
                    path=ROOT_INSTRUCTIONS_PATH,
                    kind="root-block",
                    module_id=module_id,
                    template_id=template_id,
                    version=block["version"],
                    content=content,
                    digest=managed_digest(content),
                )
            )
        for guide in module.get("supportingGuides", []):
            template_id = guide["template"]
            content = template_content(templates_root, catalog, template_id)
            artifacts.append(
                ExpectedArtifact(
                    managed_id=guide["id"],
                    path=Path(guide["path"]),
                    kind="guide",
                    module_id=module_id,
                    template_id=template_id,
                    version=guide["version"],
                    content=content,
                    digest=managed_digest(content),
                )
            )
    return artifacts


def template_content(
    templates_root: Path,
    catalog: AssetCatalog,
    template_id: str,
) -> str:
    template = catalog.templates[template_id]
    return (templates_root / template["path"]).read_text(encoding="utf-8")


def required_decisions(catalog: AssetCatalog, ordered_modules: Iterable[str]) -> list[str]:
    seen: set[str] = set()
    decisions: list[str] = []
    for module_id in ordered_modules:
        for decision_id in catalog.modules[module_id].get("requiredDecisions", []):
            if decision_id not in seen:
                decisions.append(decision_id)
                seen.add(decision_id)
    return decisions


def render_result(result: AuditResult, output_format: str) -> None:
    if output_format == "json":
        print(json.dumps(result.to_json(), indent=2, sort_keys=False))
        return
    print(render_text(result))


def render_text(result: AuditResult) -> str:
    if not result.findings and result.selection is None and not result.planned_changes:
        return "setup-context-driven audit: ok"

    lines = [
        "setup-context-driven audit: findings",
        (
            f"errors={result.summary['errors']} decisions={result.summary['decisions']} "
            f"warnings={result.summary['warnings']} info={result.summary['info']}"
        ),
    ]
    if result.selection is not None:
        lines.append("selection:")
        lines.append(
            f"- profile {result.selection.profile_id} setup {result.selection.setup_id}"
        )
        module_summary = ", ".join(
            f"{module.module_id}({module.state})" for module in result.selection.modules
        )
        lines.append(f"- modules {module_summary}")
    grouped: dict[str, list[Finding]] = {severity: [] for severity in SEVERITY_ORDER}
    for finding_item in sorted_findings(result.findings):
        grouped[finding_item.severity].append(finding_item)
    for severity in ["error", "decision", "warning", "info"]:
        if not grouped[severity]:
            continue
        lines.append(f"{severity}:")
        for finding_item in grouped[severity]:
            location = finding_item.path
            if finding_item.managed_id:
                location = f"{location} [{finding_item.managed_id}]"
            lines.append(f"- {finding_item.code} {location}: {finding_item.message}")
            lines.append(f"  action: {finding_item.action}")
    planned_changes = planned_changes_for_result(result)
    if planned_changes:
        lines.append("planned changes:")
        for change in planned_changes:
            condition = change.get("condition")
            suffix = ""
            if isinstance(condition, dict):
                condition_items = [
                    f"{key}={value}"
                    for key, value in condition.items()
                    if key != "decisionId"
                ]
                suffix = f" if {condition.get('decisionId')} {' '.join(condition_items)}"
            lines.append(
                f"- {change['action']} {change['path']} [{change['managedId']}] "
                f"state={change.get('state', 'unknown')}{suffix}"
            )
    return "\n".join(lines)


def planned_changes_for_result(result: AuditResult) -> list[dict[str, object]]:
    planned = [change.to_json() for change in result.planned_changes]
    planned.extend(planned_changes_for_findings(result.findings))
    return planned


def planned_changes_for_findings(findings: Iterable[Finding]) -> list[dict[str, object]]:
    actions = {
        "manifest.missing": "create manifest",
        "managed.block.missing": "create managed block",
        "docs.guide.missing": "create guide",
        "managed.template.stale": "refresh managed content",
        "managed.content.modified": "refresh managed content",
        "docs.language.non-english": "refresh managed content",
        "docs.reference.broken": "refresh managed content",
    }
    planned: list[dict[str, object]] = []
    seen: set[tuple[str, str, str]] = set()
    for finding_item in sorted_findings(findings):
        action = actions.get(finding_item.code)
        if action is None:
            continue
        key = (action, finding_item.path, finding_item.managed_id)
        if key in seen:
            continue
        seen.add(key)
        planned.append(
            {
                "action": action,
                "path": finding_item.path,
                "managedId": finding_item.managed_id,
            }
        )
    return planned


def exit_code_for(result: AuditResult, invalid_input: bool) -> int:
    if invalid_input:
        return 2
    summary = result.summary
    if summary["decisions"]:
        return 3
    if summary["errors"]:
        return 1
    return 0


def sorted_findings(findings: Iterable[Finding]) -> list[Finding]:
    return sorted(
        findings,
        key=lambda item: (
            SEVERITY_ORDER[item.severity],
            item.code,
            item.path,
            item.managed_id,
            item.message,
        ),
    )


def finding(
    code: str,
    severity: str,
    path: str,
    managed_id: str,
    message: str,
    action: str,
) -> Finding:
    return Finding(
        code=code,
        severity=severity,
        path=path,
        managed_id=managed_id,
        message=message,
        action=action,
    )


def marker_invalid_finding(relative_path: Path, managed_id: str, message: str) -> Finding:
    return finding(
        "managed.marker.invalid",
        "error",
        str(relative_path),
        managed_id,
        message,
        "Fix setup-context-driven ownership marker pairing.",
    )


def stale_template_finding(path: str, managed_id: str) -> Finding:
    return finding(
        "managed.template.stale",
        "warning",
        path,
        managed_id,
        "Managed artifact version or template identity is stale.",
        "Refresh the setup-owned content from the bundled template.",
    )


def broken_reference_finding(relative_path: Path, managed_id: str, target: str) -> Finding:
    return finding(
        "docs.reference.broken",
        "error",
        str(relative_path),
        managed_id,
        f"Generated content references missing path {target}.",
        "Update or restore the referenced setup-owned guide.",
    )


def contains_non_english_marker(content: str) -> bool:
    normalized = f" {content.lower()} "
    return any(marker in normalized for marker in NON_ENGLISH_MARKERS)


def managed_digest(content: str) -> str:
    normalized = content.strip() + "\n"
    return hashlib.sha256(normalized.encode("utf-8")).hexdigest()


def managed_block(managed_id: str, version: int, content: str) -> str:
    body = content.strip() + "\n"
    return (
        f"<!-- setup-context-driven:begin id={managed_id} version={version} -->\n"
        f"{body}"
        f"<!-- setup-context-driven:end id={managed_id} -->\n"
    )


def print_top_level_help() -> None:
    print(
        "\n".join(
            [
                "usage: context_setup.py [audit] [--repo PATH] [--format text|json]",
                "       context_setup.py apply --repo PATH [--format text|json]",
                "       context_setup.py sync-setups --source-dir PATH [--check] [--format text|json]",
                "",
                "Audit is the read-only default when no subcommand is supplied.",
                "Output formats: text, json. Results go to stdout; diagnostics go to stderr.",
                "Exit codes: 0 ok, 1 blocking findings, 2 invalid input, 3 decisions required.",
                "",
                "Subcommands:",
                "  audit        Read bundled assets and repository state without writes.",
                "  apply        Write confirmed managed content through an atomic change plan.",
                "  sync-setups  Check or refresh bundled canonical skill setup snapshots.",
            ]
        )
    )


def audit_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="context_setup.py audit",
        description="Audit setup-context-driven managed agent instructions without writes.",
    )
    parser.add_argument("--repo", default=".", help="Repository root. Defaults to cwd.")
    parser.add_argument(
        "--format",
        choices=["text", "json"],
        default="text",
        help="Result format written to stdout.",
    )
    parser.add_argument("--profile", help="Override the manifest profile for audit.")
    parser.add_argument(
        "--show-extra-skills",
        action="store_true",
        help="Show informational findings for installed skills outside the selected setup.",
    )
    parser.add_argument(
        "--setups-dir",
        help="Compare the bundled setup snapshot with a canonical setups directory.",
    )
    return parser


def apply_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="context_setup.py apply",
        description="Apply setup-context-driven managed agent instructions after decisions are confirmed.",
    )
    parser.add_argument("--repo", default=".", help="Repository root. Defaults to cwd.")
    parser.add_argument(
        "--format",
        choices=["text", "json"],
        default="text",
        help="Result format written to stdout.",
    )
    parser.add_argument("--profile", help="Bundled profile id to apply.")
    parser.add_argument(
        "--decision",
        action="append",
        default=[],
        help="Confirmed decision in ID=VALUE form. Repeat for multiple decisions.",
    )
    return parser


def sync_setups_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="context_setup.py sync-setups",
        description="Check or refresh bundled setup snapshots from an explicit canonical setups directory.",
    )
    parser.add_argument("--source-dir", required=True, help="Canonical setups directory.")
    parser.add_argument(
        "--check",
        action="store_true",
        help="Report snapshot drift without writing bundled assets.",
    )
    parser.add_argument(
        "--format",
        choices=["text", "json"],
        default="text",
        help="Result format written to stdout.",
    )
    return parser


def parse_unimplemented_command(command: str, args: list[str]) -> int:
    parser = argparse.ArgumentParser(
        prog=f"context_setup.py {command}",
        description=f"{command} is documented for the setup workflow but not implemented in this slice.",
    )
    parser.add_argument("--source-dir", required="--help" not in args and "-h" not in args)
    parser.add_argument("--check", action="store_true")
    parser.add_argument("--format", choices=["text", "json"], default="text")
    parser.parse_args(args)
    print(f"{command} is not implemented by this task.", file=sys.stderr)
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
