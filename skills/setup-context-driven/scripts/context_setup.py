#!/usr/bin/env python3
"""Audit setup-context-driven managed agent instructions."""

from __future__ import annotations

import argparse
import base64
from collections.abc import Mapping
import hashlib
import json
import os
import re
import shutil
import stat
import subprocess
import sys
import tempfile
from datetime import date
from dataclasses import dataclass, fields, is_dataclass
from functools import lru_cache
from pathlib import Path
from typing import Iterable

from context_assets import (
    AssetCatalog,
    AssetValidationError,
    ExternalSkillContract,
    PortableTreeError,
    RepositoryOwnedExtension,
    TEMPLATE_TOKEN,
    UpgradeTransition,
    build_standard_profile_plan,
    is_http_contract_decision_value_valid,
    load_asset_catalog,
    portable_file_digest,
    portable_tree_digest,
    render_standard_profile_snapshot,
    setup_snapshot_digest,
)
from context_capabilities import (
    CapabilityEvaluation,
    RequirementStrength,
    UNIVERSAL_CAPABILITIES,
    evaluate_repository_capabilities,
    render_capability_json,
)
from context_baseline import (
    DecisionDocumentDiagnostic,
    DecisionDocumentError,
    IncompatibleSourceBaseline,
    ReadoptionDisposition,
    SourceInventoryError,
    StructuredDecisionDocument,
    inventory_incompatible_source_baseline,
    load_decision_document,
    merge_decision_documents,
    repository_rules_proposed_bytes,
    validate_readoption_decisions,
)


AUDIT_SCHEMA_VERSION = "setup-context-driven/audit-v1"
MANIFEST_SCHEMA_VERSION_0_0_1 = "setup-context-driven/manifest/0.0.1"
OWNED_VERSION_0_0_1 = "0.0.1"
RESTORE_SCHEMA_VERSION = "setup-context-driven/restore-v1"
LOCK_HASH_COMPATIBILITY_SCHEMA_VERSION = (
    "setup-context-driven/external-lock-hash-compatibility-v1"
)
LOCK_HASH_COMPATIBILITY_FIXTURE = Path(
    "assets/lock-hash-compatibility-v1.json"
)
MANIFEST_PATH = Path("docs/agents/setup-context.json")
ROOT_INSTRUCTIONS_PATH = Path("AGENTS.md")
SEVERITY_ORDER = {"error": 0, "decision": 1, "warning": 2, "info": 3}
BEGIN_MARKER = re.compile(
    r"<!--\s*setup-context-driven:begin\s+id=([A-Za-z0-9_.-]+)\s+version=([0-9]+(?:\.[0-9]+\.[0-9]+)?)\s*-->"
)
END_MARKER = re.compile(
    r"<!--\s*setup-context-driven:end\s+id=([A-Za-z0-9_.-]+)\s*-->"
)
MARKER = re.compile(r"<!--\s*setup-context-driven:(begin|end)\b[^>]*-->")
MARKDOWN_LINK = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
DELEGATION_SIGNAL = re.compile(
    r"\b(?:consult|defer|deferred|delegate|delegated|delegation|defined|documented|"
    r"follow|governed|maintained|read|refer|see)\b",
    re.IGNORECASE,
)
DELEGATION_TARGET = re.compile(
    r"(?:^|[^A-Za-z0-9_.-])(?:\.\.?/|/)*(?:AGENTS\.md|CLAUDE\.md|"
    r"docs/agents/[A-Za-z0-9_./-]+\.md)(?=$|[^A-Za-z0-9_-])",
    re.IGNORECASE,
)
DELEGATION_DOCUMENT_NAMES = frozenset({"AGENTS.md", "CLAUDE.md"})
DELEGATION_IGNORED_DIRECTORIES = frozenset(
    {
        ".git",
        ".hg",
        ".svn",
        ".venv",
        "__pycache__",
        "node_modules",
        "vendor",
        "venv",
    }
)
DELEGATION_IGNORED_PREFIXES = (
    (".agents", "skills", "setup-context-driven"),
    ("skills", "setup-context-driven"),
)
NON_ENGLISH_MARKERS = [
    " não ",
    " obrigatório",
    " repositório",
    " arquivo",
    " configuración",
    " repositorio",
]
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
IMMUTABLE_GIT_REF = re.compile(r"^(?:[0-9a-f]{40}|[0-9a-f]{64})$")
GITHUB_REMOTE = re.compile(
    r"^(?:https://github\.com/|git@github\.com:|ssh://git@github\.com/)([A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+?)(?:\.git)?$"
)
REPO_OWNED_SKILLS = {
    "archive-spec",
    "brainstorming",
    "business-analyst",
    "council",
    "evidence-gate",
    "implement-spec",
    "implement-task",
    "qa-gate",
    "roundfix",
    "setup-context-driven",
    "write-idea",
    "write-prd",
    "write-tasks",
    "write-techspec",
}


@dataclass(frozen=True)
class Finding:
    code: str
    severity: str
    path: str
    managed_id: str
    message: str
    action: str
    remediation: dict[str, object] | None = None

    def to_json(self) -> dict[str, object]:
        data: dict[str, object] = {
            "code": self.code,
            "severity": self.severity,
            "path": self.path,
            "managedId": self.managed_id,
            "message": self.message,
            "action": self.action,
        }
        if self.remediation is not None:
            data["remediation"] = self.remediation
        return data


@dataclass(frozen=True)
class ExpectedArtifact:
    managed_id: str
    path: Path
    kind: str
    module_id: str
    template_id: str
    version: int | str
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
    reason: str
    before_digest: str | None = None
    after_digest: str | None = None
    condition: PlanCondition | None = None
    from_path: Path | None = None
    reference_edits: tuple[dict[str, str], ...] = ()

    def to_json(self) -> dict[str, object]:
        data: dict[str, object] = {
            "action": self.action,
            "path": self.path.as_posix(),
            "managedId": self.managed_id,
            "state": self.state,
            "reason": self.reason,
            "beforeDigest": self.before_digest,
            "afterDigest": self.after_digest,
        }
        if self.condition is not None:
            data["condition"] = self.condition.to_json()
        if self.from_path is not None:
            data["fromPath"] = self.from_path.as_posix()
        if self.reference_edits:
            data["referenceEdits"] = list(self.reference_edits)
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
class RetentionEntry:
    from_clause: str
    enforcement: str
    disposition: str
    targets: tuple[str, ...]
    reason: str

    def to_json(self) -> dict[str, object]:
        return {
            "fromClause": self.from_clause,
            "enforcement": self.enforcement,
            "disposition": self.disposition,
            "targets": list(self.targets),
            "reason": self.reason,
        }


@dataclass(frozen=True)
class ManagedBlock:
    managed_id: str
    version: int | str
    body: str


@dataclass(frozen=True)
class ManagedBlockSpan:
    managed_id: str
    version: int | str
    body: str
    start: int
    end: int


@dataclass(frozen=True)
class FileMutation:
    path: Path
    before_digest: str | None
    after_digest: str | None
    content: bytes | None
    operations: tuple[PlannedChange, ...]

    def output_json(self) -> dict[str, object]:
        return {
            "path": self.path.as_posix(),
            "beforeDigest": self.before_digest,
            "afterDigest": self.after_digest,
            "afterBytes": (
                base64.b64encode(self.content).decode("ascii")
                if self.content is not None
                else None
            ),
            "encoding": "base64",
        }


@dataclass(frozen=True)
class ReadoptionPlanContext:
    source_baseline: IncompatibleSourceBaseline
    dispositions: tuple[ReadoptionDisposition, ...]
    decision_document: dict[str, object]
    decision_file_digests: tuple[str, ...]
    capabilities: tuple[dict[str, object], ...]
    setup_snapshot: dict[str, object]
    verification: tuple[dict[str, str], ...]
    repository_rules: bytes


@dataclass(frozen=True)
class ChangePlan:
    kind: str
    mutations: tuple[FileMutation, ...]
    digest: str | None
    manifest: dict
    retention: tuple[RetentionEntry, ...] = ()
    readoption: ReadoptionPlanContext | None = None

    @property
    def changes(self) -> list[FileMutation]:
        """Compatibility view for callers that only need affected paths."""
        return list(self.mutations)


@dataclass(frozen=True)
class InstalledSkill:
    name: str
    path: Path
    root: Path
    locked: bool
    origin: dict
    content_digest: str | None
    unsafe_error: str | None = None


@dataclass(frozen=True)
class GitSourceCheckout:
    root: Path
    repository: str
    revision: str


@dataclass(frozen=True)
class AuditResult:
    findings: list[Finding]
    selection: PreviewSelection | None = None
    planned_changes: tuple[PlannedChange, ...] = ()
    plan_digest: str | None = None
    retention: tuple[RetentionEntry, ...] = ()
    source_baseline: IncompatibleSourceBaseline | None = None
    decision_document: dict[str, object] | None = None
    planned_outputs: tuple[FileMutation, ...] = ()
    capabilities: tuple[dict[str, object], ...] = ()
    setup_snapshot: dict[str, object] | None = None
    verification: tuple[dict[str, str], ...] = ()

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
        if self.plan_digest is not None:
            data["planDigest"] = self.plan_digest
        if self.retention:
            data["retentionAccounting"] = [
                entry.to_json() for entry in self.retention
            ]
        if self.source_baseline is not None:
            data["sourceBaseline"] = self.source_baseline.identity_json()
            data["sourceEntries"] = [
                entry.to_json() for entry in self.source_baseline.entries
            ]
        if self.decision_document is not None:
            data["decisionDocument"] = self.decision_document
        if self.planned_outputs:
            data["plannedOutputs"] = [
                mutation.output_json() for mutation in self.planned_outputs
            ]
        if self.capabilities:
            data["capabilities"] = list(self.capabilities)
        if self.setup_snapshot is not None:
            data["setupSnapshot"] = self.setup_snapshot
        if self.verification:
            data["verification"] = list(self.verification)
        return data


@dataclass(frozen=True)
class RestoreLimits:
    max_files: int = 2_000
    max_bytes: int = 8 * 1024 * 1024


@dataclass(frozen=True)
class DelegationScanLimits:
    max_files: int = 256
    max_bytes: int = 2 * 1024 * 1024


@dataclass(frozen=True)
class RestoreFile:
    path: Path
    content: bytes


@dataclass(frozen=True)
class RestoreFileChange:
    action: str
    path: Path
    skill_name: str
    before_digest: str | None
    after_digest: str | None

    def to_json(self) -> dict[str, object]:
        return {
            "action": self.action,
            "path": self.path.as_posix(),
            "skill": self.skill_name,
            "beforeDigest": self.before_digest,
            "afterDigest": self.after_digest,
        }


@dataclass(frozen=True)
class RestoreLockEdit:
    skill_name: str
    before: dict | None
    after: dict

    def to_json(self) -> dict[str, object]:
        return {
            "action": "update-lock-entry",
            "path": "skills-lock.json",
            "skill": self.skill_name,
            "before": self.before,
            "after": self.after,
        }


@dataclass(frozen=True)
class RestoreSkillPlan:
    contract: ExternalSkillContract
    target: Path
    files: tuple[RestoreFile, ...]
    changes: tuple[RestoreFileChange, ...]
    lock_edit: RestoreLockEdit | None
    observed_digest: str | None

    def to_json(self) -> dict[str, object]:
        source = self.contract.source
        return {
            "skill": self.contract.skill_name,
            "targetPath": self.target.as_posix(),
            "source": {
                "provider": source.provider,
                "repository": source.repository,
                "ref": source.revision,
                "path": source.source_path.as_posix(),
            },
            "expectedDigest": self.contract.tree_digest,
            "observedDigest": self.observed_digest,
            "changes": [change.to_json() for change in self.changes],
            "lockEdit": self.lock_edit.to_json() if self.lock_edit is not None else None,
        }


@dataclass(frozen=True)
class RestorePlan:
    profile_id: str
    setup_id: str
    acquisitions: tuple[dict[str, str], ...]
    skills: tuple[RestoreSkillPlan, ...]
    lock_before: bytes | None
    lock_after: bytes | None
    digest: str

    @property
    def has_changes(self) -> bool:
        return bool(self.skills)

    @property
    def planned_changes(self) -> list[dict[str, object]]:
        changes: list[dict[str, object]] = []
        for skill in self.skills:
            changes.extend(change.to_json() for change in skill.changes)
            if skill.lock_edit is not None:
                changes.append(skill.lock_edit.to_json())
        return changes

    def to_json(
        self,
        *,
        ok: bool,
        applied: bool,
        finding_item: dict[str, object] | None = None,
    ) -> dict[str, object]:
        payload: dict[str, object] = {
            "schemaVersion": RESTORE_SCHEMA_VERSION,
            "ok": ok,
            "applied": applied,
            "profile": self.profile_id,
            "setup": self.setup_id,
            "acquisitions": list(self.acquisitions),
            "skills": [skill.to_json() for skill in self.skills],
            "plannedChanges": self.planned_changes,
            "planDigest": self.digest,
        }
        if finding_item is not None:
            payload["finding"] = finding_item
        return payload


class RestoreError(Exception):
    def __init__(self, code: str, message: str, action: str, exit_code: int = 1):
        self.code = code
        self.message = message
        self.action = action
        self.exit_code = exit_code
        super().__init__(message)

    def to_json(self) -> dict[str, object]:
        return {"code": self.code, "message": self.message, "action": self.action}


class GitSkillSource:
    def __init__(self, source_dir: Path | None = None):
        self.source_dir = source_dir

    def acquire(self, source, target: Path) -> None:
        if source.provider != "github":
            raise RestoreError(
                "source.provider-unsupported",
                f"Skill source provider {source.provider!r} is not supported.",
                "Use a bundled snapshot whose external source provider is github.",
            )
        if self.source_dir is not None:
            origin = str(self.source_dir)
        else:
            origin = f"https://github.com/{source.repository}.git"
        target.parent.mkdir(parents=True, exist_ok=True)
        try:
            run_git_argv("init", "--bare", str(target))
            run_git_argv(
                "--git-dir",
                str(target),
                "fetch",
                "--no-tags",
                "--depth=1",
                origin,
                source.revision,
            )
            resolved = run_git_argv(
                "--git-dir",
                str(target),
                "rev-parse",
                "--verify",
                "FETCH_HEAD^{commit}",
            ).decode("ascii").strip()
        except FileNotFoundError as error:
            raise RestoreError(
                "source.git-missing",
                "Git is required for immutable external skill acquisition.",
                "Install Git and rerun the same restore-skills preview.",
            ) from error
        except (OSError, UnicodeDecodeError, subprocess.CalledProcessError) as error:
            raise RestoreError(
                "source.commit-unavailable",
                f"The exact commit {source.revision} is unavailable: {git_error_message(error)}.",
                "Make the declared commit available in --source-dir or verify the bundled GitHub provenance.",
            ) from error
        if resolved != source.revision:
            raise RestoreError(
                "source.commit-mismatch",
                f"Git resolved {resolved}, not the declared commit {source.revision}.",
                "Reject the source and verify its immutable commit identity before retrying.",
            )


def lock_adapter_incompatible(reason: str) -> RestoreError:
    return RestoreError(
        "lock.adapter-incompatible",
        "The lock compatibility fixture cannot prove agreement with Spec 0036 "
        f"Task 01: {reason}.",
        "Update the isolated lock adapter before restoring skills.",
    )


def load_lock_hash_compatibility_fixture(
    fixture_path: Path,
) -> tuple[tuple[RestoreFile, ...], str]:
    try:
        document = json.loads(fixture_path.read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        raise lock_adapter_incompatible(
            f"could not load {fixture_path.as_posix()}: {error}"
        ) from error

    required_fields = {
        "schemaVersion",
        "version",
        "files",
        "expectedSha256",
    }
    if not isinstance(document, dict) or set(document) != required_fields:
        raise lock_adapter_incompatible(
            "the fixture must contain only schemaVersion, version, files, and expectedSha256"
        )
    if (
        document["schemaVersion"] != LOCK_HASH_COMPATIBILITY_SCHEMA_VERSION
        or type(document["version"]) is not int
        or document["version"] != 1
    ):
        raise lock_adapter_incompatible(
            "the fixture schemaVersion and version must identify version 1"
        )
    expected = document["expectedSha256"]
    if not isinstance(expected, str) or re.fullmatch(r"[0-9a-f]{64}", expected) is None:
        raise lock_adapter_incompatible(
            "expectedSha256 must be a lowercase SHA-256 digest"
        )
    declared_files = document["files"]
    if not isinstance(declared_files, list) or not declared_files:
        raise lock_adapter_incompatible("files must be a non-empty array")

    files: list[RestoreFile] = []
    seen_paths: set[str] = set()
    for index, item in enumerate(declared_files):
        if not isinstance(item, dict) or set(item) != {"path", "content"}:
            raise lock_adapter_incompatible(
                f"files[{index}] must contain only path and content"
            )
        path_value = item["path"]
        content = item["content"]
        if not isinstance(path_value, str) or not isinstance(content, str):
            raise lock_adapter_incompatible(
                f"files[{index}] path and content must be strings"
            )
        path_parts = path_value.split("/")
        if (
            not path_value
            or "\\" in path_value
            or "\x00" in path_value
            or path_value.startswith("/")
            or any(part in {"", ".", ".."} for part in path_parts)
            or any(part in {".git", "node_modules"} for part in path_parts[:-1])
        ):
            raise lock_adapter_incompatible(
                f"files[{index}] path is not a safe slash-normalized relative path"
            )
        if path_value in seen_paths:
            raise lock_adapter_incompatible(
                f"files[{index}] duplicates path {path_value!r}"
            )
        try:
            path_value.encode("utf-8")
            content_bytes = content.encode("utf-8")
        except UnicodeEncodeError as error:
            raise lock_adapter_incompatible(
                f"files[{index}] path and content must be valid UTF-8"
            ) from error
        seen_paths.add(path_value)
        files.append(RestoreFile(Path(*path_parts), content_bytes))
    return tuple(files), expected


class ExternalSkillLockAdapter:
    def __init__(self, fixture_path: Path | None = None):
        self.fixture_path = fixture_path or (
            Path(__file__).resolve().parents[1] / LOCK_HASH_COMPATIBILITY_FIXTURE
        )

    def assert_compatible(self) -> None:
        files, expected = load_lock_hash_compatibility_fixture(self.fixture_path)
        observed = external_lock_digest(files)
        if observed != expected:
            raise lock_adapter_incompatible(
                f"computed digest {observed} does not match pinned digest {expected}"
            )

    def entry_for(self, source, tree: tuple[RestoreFile, ...]) -> dict:
        self.assert_compatible()
        if not any(item.path == Path("SKILL.md") for item in tree):
            raise RestoreError(
                "source.skill-file-missing",
                "The verified source subtree does not contain SKILL.md.",
                "Fix the immutable source snapshot before restoring this skill.",
            )
        return {
            "source": source.repository,
            "ref": source.revision,
            "sourceType": source.provider,
            "skillPath": (source.source_path / "SKILL.md").as_posix(),
            "computedHash": external_lock_digest(tree),
        }


class RestoreFilesystem:
    def replace(self, source: Path, target: Path) -> None:
        source.replace(target)

    def remove_tree(self, target: Path) -> None:
        if target.is_symlink() or target.is_file():
            target.unlink()
        elif target.exists():
            shutil.rmtree(target)

    def verify_tree(self, target: Path, expected_digest: str) -> None:
        try:
            observed = portable_tree_digest(target)
        except PortableTreeError as error:
            raise OSError(str(error)) from error
        if observed != expected_digest:
            raise OSError(
                f"postwrite tree digest {observed} does not match {expected_digest}"
            )

    def verify_lock(self, target: Path, expected: bytes) -> None:
        if target.read_bytes() != expected:
            raise OSError("postwrite skills-lock.json bytes differ from the authorized plan")


def external_lock_digest(files: Iterable[RestoreFile]) -> str:
    digest = hashlib.sha256()
    for item in sorted(files, key=lambda candidate: candidate.path.as_posix().encode("utf-8")):
        digest.update(item.path.as_posix().encode("utf-8"))
        digest.update(item.content)
    return digest.hexdigest()


def run_git_argv(*args: str) -> bytes:
    environment = os.environ.copy()
    environment.update(
        {
            "GIT_TERMINAL_PROMPT": "0",
            "GCM_INTERACTIVE": "Never",
            "GIT_ASKPASS": "",
            "SSH_ASKPASS": "",
        }
    )
    result = subprocess.run(
        ["git", *args],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
        env=environment,
    )
    if result.returncode != 0:
        raise subprocess.CalledProcessError(
            result.returncode,
            result.args,
            output=result.stdout,
            stderr=result.stderr,
        )
    return result.stdout


def build_restore_plan(
    repo: Path,
    catalog: AssetCatalog,
    profile_id: str,
    selected_skills: list[str] | tuple[str, ...] | None = None,
    source_dir: Path | None = None,
    *,
    skill_source: GitSkillSource | None = None,
    lock_adapter: ExternalSkillLockAdapter | None = None,
    limits: RestoreLimits | None = None,
) -> RestorePlan:
    if profile_id not in catalog.profiles:
        raise RestoreError(
            "restore.profile-unknown",
            f"Unknown bundled profile {profile_id!r}.",
            "Choose a profile id from the bundled setup-context-driven assets.",
            exit_code=2,
        )
    setup_id = catalog.profiles[profile_id]["setup"]
    contracts = {
        contract.skill_name: contract
        for contract in catalog.external_sources_by_setup.get(setup_id, ())
    }
    requested = list(selected_skills or [])
    if requested:
        unknown = sorted({name for name in requested if name not in contracts})
        if unknown:
            raise RestoreError(
                "restore.skill-invalid",
                "Selected skills are not external members of this profile: "
                + ", ".join(unknown)
                + ".",
                "Choose repeatable --skill values from this profile's external Repository Skill Set.",
                exit_code=2,
            )
        selected_names = sorted(set(requested))
    else:
        selected_names = sorted(contracts)

    resolved_source_dir = source_dir.resolve(strict=False) if source_dir is not None else None
    if resolved_source_dir is not None and not resolved_source_dir.is_dir():
        raise RestoreError(
            "restore.source-dir-invalid",
            f"Offline Git object store is not a directory: {resolved_source_dir}.",
            "Pass an existing Git checkout or bare object store to --source-dir.",
            exit_code=2,
        )
    source_adapter = skill_source or GitSkillSource(resolved_source_dir)
    adapter = lock_adapter or ExternalSkillLockAdapter()
    adapter.assert_compatible()
    active_limits = limits or RestoreLimits()
    if active_limits.max_files < 1 or active_limits.max_bytes < 1:
        raise RestoreError(
            "restore.limit-invalid",
            "Restoration limits must be positive.",
            "Use positive file-count and byte limits.",
            exit_code=2,
        )

    lock_path = repo / "skills-lock.json"
    lock_document, lock_before = load_restore_lock(lock_path)
    lock_after_document = copy_json(lock_document)
    lock_entries = lock_after_document["skills"]
    acquisitions: list[dict[str, str]] = []
    acquired_files: dict[str, tuple[RestoreFile, ...]] = {}
    acquired_count = 0
    acquired_bytes = 0

    groups: dict[tuple[str, str, str], list[ExternalSkillContract]] = {}
    for name in selected_names:
        contract = contracts[name]
        source = contract.source
        groups.setdefault(
            (source.provider, source.repository, source.revision), []
        ).append(contract)

    with tempfile.TemporaryDirectory(prefix="setup-context-driven-restore-") as temp_dir:
        temp_root = Path(temp_dir)
        for index, (group_key, group_contracts) in enumerate(sorted(groups.items())):
            provider, repository, revision = group_key
            source = group_contracts[0].source
            object_store = temp_root / f"source-{index}"
            source_adapter.acquire(source, object_store)
            acquisitions.append(
                {"provider": provider, "repository": repository, "ref": revision}
            )
            for contract in sorted(group_contracts, key=lambda item: item.skill_name):
                files = read_verified_git_skill_tree(object_store, contract, active_limits)
                acquired_count += len(files)
                acquired_bytes += sum(len(item.content) for item in files)
                if acquired_count > active_limits.max_files or acquired_bytes > active_limits.max_bytes:
                    raise RestoreError(
                        "source.limit-exceeded",
                        "Acquired skill content exceeds the restoration file-count or byte limit.",
                        "Reduce the declared immutable source subtree before retrying.",
                    )
                acquired_files[contract.skill_name] = files

    skill_plans: list[RestoreSkillPlan] = []
    for name in selected_names:
        contract = contracts[name]
        files = acquired_files[name]
        target = Path(".agents") / "skills" / name
        absolute_target = safe_restore_target(repo, target)
        current_files = inspect_restore_target(absolute_target, active_limits)
        observed_digest = (
            portable_file_digest(
                (item.path.as_posix().encode("utf-8"), item.content)
                for item in current_files
            )
            if absolute_target.exists()
            else None
        )
        changes = restore_file_changes(name, target, current_files, files)
        expected_entry = adapter.entry_for(contract.source, files)
        before_entry = lock_entries.get(name)
        lock_edit = None
        if before_entry != expected_entry:
            lock_edit = RestoreLockEdit(
                skill_name=name,
                before=copy_json(before_entry) if isinstance(before_entry, dict) else None,
                after=expected_entry,
            )
            lock_entries[name] = expected_entry
        if not changes and lock_edit is None:
            continue
        skill_plans.append(
            RestoreSkillPlan(
                contract=contract,
                target=target,
                files=files,
                changes=changes,
                lock_edit=lock_edit,
                observed_digest=observed_digest,
            )
        )

    lock_edits = any(skill.lock_edit is not None for skill in skill_plans)
    lock_after = (
        (json.dumps(lock_after_document, indent=2, sort_keys=False) + "\n").encode("utf-8")
        if lock_edits
        else lock_before
    )
    digest_payload = {
        "kind": "restore-skills",
        "profile": profile_id,
        "setup": setup_id,
        "acquisitions": acquisitions,
        "skills": [skill.to_json() for skill in skill_plans],
        "plannedChanges": [
            item
            for skill in skill_plans
            for item in (
                [change.to_json() for change in skill.changes]
                + ([skill.lock_edit.to_json()] if skill.lock_edit is not None else [])
            )
        ],
        "lockBeforeDigest": bytes_digest(lock_before),
        "lockAfterDigest": bytes_digest(lock_after),
    }
    plan_digest = hashlib.sha256(canonical_json_bytes(digest_payload)).hexdigest()
    return RestorePlan(
        profile_id=profile_id,
        setup_id=setup_id,
        acquisitions=tuple(acquisitions),
        skills=tuple(skill_plans),
        lock_before=lock_before,
        lock_after=lock_after,
        digest=plan_digest,
    )


def copy_json(value):
    return json.loads(json.dumps(value))


def load_restore_lock(path: Path) -> tuple[dict, bytes | None]:
    if path.is_symlink():
        raise RestoreError(
            "lock.unsafe-path",
            "skills-lock.json is a symbolic link.",
            "Replace it with a regular repository file before restoring skills.",
        )
    try:
        raw = path.read_bytes()
    except FileNotFoundError:
        return {"version": 1, "skills": {}}, None
    except OSError as error:
        raise RestoreError(
            "lock.read-failed",
            f"Could not read skills-lock.json: {error}.",
            "Fix repository permissions and rerun restore-skills.",
        ) from error
    try:
        document = json.loads(raw.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as error:
        raise RestoreError(
            "lock.invalid",
            f"skills-lock.json is malformed: {error}.",
            "Fix skills-lock.json before restoring external skills.",
            exit_code=2,
        ) from error
    if (
        not isinstance(document, dict)
        or document.get("version") != 1
        or not isinstance(document.get("skills"), dict)
    ):
        raise RestoreError(
            "lock.invalid",
            "skills-lock.json must use version 1 with a skills object.",
            "Fix skills-lock.json before restoring external skills.",
            exit_code=2,
        )
    return document, raw


def read_verified_git_skill_tree(
    object_store: Path,
    contract: ExternalSkillContract,
    limits: RestoreLimits,
) -> tuple[RestoreFile, ...]:
    source = contract.source
    source_path = source.source_path.as_posix()
    try:
        output = run_git_argv(
            "--git-dir",
            str(object_store),
            "ls-tree",
            "-r",
            "-z",
            "--full-tree",
            source.revision,
            "--",
            source_path,
        )
    except (OSError, subprocess.CalledProcessError) as error:
        raise RestoreError(
            "source.commit-unavailable",
            f"Could not inspect {source.repository}@{source.revision}: {git_error_message(error)}.",
            "Verify the exact commit and source path in the bundled snapshot.",
        ) from error

    files: list[RestoreFile] = []
    total_bytes = 0
    prefix = source_path.encode("utf-8") + b"/"
    for raw_entry in output.split(b"\0"):
        if not raw_entry:
            continue
        try:
            header, path_bytes = raw_entry.split(b"\t", 1)
            mode, object_type, object_id = header.split(b" ", 2)
            decoded_path = path_bytes.decode("utf-8")
        except (ValueError, UnicodeDecodeError) as error:
            raise RestoreError(
                "source.unsafe-tree",
                "The acquired Git tree contains an undecodable or malformed entry.",
                "Remove unsafe entries from the immutable source and publish a new snapshot.",
            ) from error
        if not path_bytes.startswith(prefix):
            raise RestoreError(
                "source.unsafe-tree",
                f"Git tree entry escapes {source_path}: {decoded_path}.",
                "Remove traversal entries from the immutable source.",
            )
        relative_bytes = path_bytes[len(prefix) :]
        try:
            relative_text = relative_bytes.decode("utf-8")
        except UnicodeDecodeError as error:
            raise RestoreError(
                "source.unsafe-tree",
                "The acquired skill subtree contains a non-UTF-8 path.",
                "Rename the source entry to a portable UTF-8 relative path.",
            ) from error
        relative = safe_relative_path(relative_text)
        if relative is None or relative.as_posix() != relative_text:
            raise RestoreError(
                "source.unsafe-tree",
                f"The acquired skill subtree contains an unsafe path: {relative_text!r}.",
                "Remove traversal or non-portable entries from the immutable source.",
            )
        if any(part in {".git", "node_modules"} for part in relative.parts[:-1]):
            continue
        if object_type != b"blob" or mode not in {b"100644", b"100755"}:
            raise RestoreError(
                "source.unsafe-tree",
                f"The acquired skill subtree contains a link, device, or unsupported entry: {relative_text}.",
                "Replace it with regular files before publishing a new immutable snapshot.",
            )
        try:
            content = run_git_argv(
                "--git-dir", str(object_store), "cat-file", "blob", object_id.decode("ascii")
            )
        except (OSError, UnicodeDecodeError, subprocess.CalledProcessError) as error:
            raise RestoreError(
                "source.read-failed",
                f"Could not read acquired file {relative_text}: {git_error_message(error)}.",
                "Verify the offline object store or declared GitHub commit and retry.",
            ) from error
        files.append(RestoreFile(relative, content))
        total_bytes += len(content)
        if len(files) > limits.max_files or total_bytes > limits.max_bytes:
            raise RestoreError(
                "source.limit-exceeded",
                "Acquired skill content exceeds the restoration file-count or byte limit.",
                "Reduce the declared immutable source subtree before retrying.",
            )
    if not files:
        raise RestoreError(
            "source.empty-tree",
            f"The declared source path {source_path} contains no restorable files.",
            "Fix the snapshot source path and immutable revision.",
        )
    files.sort(key=lambda item: item.path.as_posix().encode("utf-8"))
    observed = portable_file_digest(
        (item.path.as_posix().encode("utf-8"), item.content) for item in files
    )
    if observed != contract.tree_digest:
        raise RestoreError(
            "source.digest-mismatch",
            f"Acquired skill {contract.skill_name} digest {observed} does not match {contract.tree_digest}.",
            "Regenerate the bundled snapshot from the exact committed skill subtree.",
        )
    return tuple(files)


def inspect_restore_target(target: Path, limits: RestoreLimits) -> tuple[RestoreFile, ...]:
    try:
        root_stat = target.lstat()
    except FileNotFoundError:
        return ()
    except OSError as error:
        raise RestoreError(
            "target.read-failed",
            f"Could not inspect restoration target {target}: {error}.",
            "Fix repository permissions and rerun restore-skills.",
        ) from error
    if stat.S_ISLNK(root_stat.st_mode) or not stat.S_ISDIR(root_stat.st_mode):
        raise RestoreError(
            "target.unsafe-tree",
            f"Restoration target is not a regular directory: {target}.",
            "Replace the unsafe target manually, then rerun restore-skills.",
        )
    files: list[RestoreFile] = []
    total_bytes = 0
    pending = [target]
    while pending:
        directory = pending.pop()
        try:
            entries = list(os.scandir(directory))
        except OSError as error:
            raise RestoreError(
                "target.read-failed",
                f"Could not read restoration target {directory}: {error}.",
                "Fix repository permissions and rerun restore-skills.",
            ) from error
        for entry in entries:
            path = Path(entry.path)
            relative = path.relative_to(target)
            try:
                entry_stat = entry.stat(follow_symlinks=False)
            except OSError as error:
                raise RestoreError(
                    "target.read-failed",
                    f"Could not inspect target entry {relative.as_posix()}: {error}.",
                    "Fix repository permissions and rerun restore-skills.",
                ) from error
            if stat.S_ISDIR(entry_stat.st_mode):
                pending.append(path)
                continue
            if stat.S_ISLNK(entry_stat.st_mode) or not stat.S_ISREG(entry_stat.st_mode) or entry_stat.st_nlink != 1:
                raise RestoreError(
                    "target.unsafe-tree",
                    f"Restoration target contains a link or special entry: {relative.as_posix()}.",
                    "Remove the unsafe target entry manually, then rerun restore-skills.",
                )
            try:
                content = path.read_bytes()
            except OSError as error:
                raise RestoreError(
                    "target.read-failed",
                    f"Could not read target file {relative.as_posix()}: {error}.",
                    "Fix repository permissions and rerun restore-skills.",
                ) from error
            files.append(RestoreFile(relative, content))
            total_bytes += len(content)
            if len(files) > limits.max_files or total_bytes > limits.max_bytes:
                raise RestoreError(
                    "target.limit-exceeded",
                    "Existing restoration target exceeds the safe inspection limits.",
                    "Reduce the target directory before retrying.",
                )
    files.sort(key=lambda item: item.path.as_posix().encode("utf-8"))
    return tuple(files)


def restore_file_changes(
    skill_name: str,
    target: Path,
    before: tuple[RestoreFile, ...],
    after: tuple[RestoreFile, ...],
) -> tuple[RestoreFileChange, ...]:
    before_by_path = {item.path: item.content for item in before}
    after_by_path = {item.path: item.content for item in after}
    changes: list[RestoreFileChange] = []
    for relative in sorted(
        set(before_by_path) | set(after_by_path),
        key=lambda item: item.as_posix().encode("utf-8"),
    ):
        before_content = before_by_path.get(relative)
        after_content = after_by_path.get(relative)
        if before_content == after_content:
            continue
        action = "refresh"
        if before_content is None:
            action = "create"
        elif after_content is None:
            action = "remove"
        changes.append(
            RestoreFileChange(
                action=action,
                path=target / relative,
                skill_name=skill_name,
                before_digest=bytes_digest(before_content),
                after_digest=bytes_digest(after_content),
            )
        )
    return tuple(changes)


def apply_restore_plan(
    repo: Path,
    plan: RestorePlan,
    *,
    filesystem: RestoreFilesystem | None = None,
) -> None:
    fs = filesystem or RestoreFilesystem()
    assert_restore_preimage(repo, plan)
    stages: dict[Path, Path] = {}
    directory_states: list[dict[str, object]] = []
    lock_state: dict[str, object] | None = None
    cleanup_paths: list[Path] = []
    try:
        for skill in plan.skills:
            if not skill.changes:
                continue
            target = safe_restore_target(repo, skill.target)
            target.parent.mkdir(parents=True, exist_ok=True)
            stage = Path(
                tempfile.mkdtemp(prefix=".restore-stage-", dir=target.parent)
            )
            cleanup_paths.append(stage)
            for item in skill.files:
                destination = stage / item.path
                destination.parent.mkdir(parents=True, exist_ok=True)
                destination.write_bytes(item.content)
            fs.verify_tree(stage, skill.contract.tree_digest)
            stages[target] = stage

        lock_target = repo / "skills-lock.json"
        lock_temp: Path | None = None
        if plan.lock_after is not None and plan.lock_after != plan.lock_before:
            lock_target.parent.mkdir(parents=True, exist_ok=True)
            descriptor, temp_name = tempfile.mkstemp(
                prefix=".restore-lock-", dir=lock_target.parent
            )
            lock_temp = Path(temp_name)
            cleanup_paths.append(lock_temp)
            with os.fdopen(descriptor, "wb") as handle:
                handle.write(plan.lock_after)

        for target, stage in sorted(stages.items(), key=lambda item: item[0].as_posix()):
            backup_container = Path(
                tempfile.mkdtemp(prefix=".restore-backup-", dir=target.parent)
            )
            cleanup_paths.append(backup_container)
            backup = backup_container / "original"
            state: dict[str, object] = {
                "target": target,
                "backup": backup,
                "touched": False,
            }
            directory_states.append(state)
            if target.exists() or target.is_symlink():
                fs.replace(target, backup)
            state["touched"] = True
            fs.replace(stage, target)

        if lock_temp is not None:
            backup_container = Path(
                tempfile.mkdtemp(prefix=".restore-backup-lock-", dir=lock_target.parent)
            )
            cleanup_paths.append(backup_container)
            backup = backup_container / "original"
            lock_state = {
                "target": lock_target,
                "backup": backup,
                "touched": False,
            }
            if lock_target.exists() or lock_target.is_symlink():
                fs.replace(lock_target, backup)
            lock_state["touched"] = True
            fs.replace(lock_temp, lock_target)

        for skill in plan.skills:
            target = safe_restore_target(repo, skill.target)
            fs.verify_tree(target, skill.contract.tree_digest)
        if plan.lock_after is not None:
            fs.verify_lock(lock_target, plan.lock_after)
    except (OSError, PortableTreeError, RestoreError) as error:
        rollback_errors: list[str] = []
        states = ([lock_state] if lock_state is not None else []) + list(
            reversed(directory_states)
        )
        for state in states:
            if state is None or not state["touched"]:
                continue
            target = state["target"]
            backup = state["backup"]
            try:
                if target.exists() or target.is_symlink():
                    fs.remove_tree(target)
                if backup.exists() or backup.is_symlink():
                    fs.replace(backup, target)
            except OSError as rollback_error:
                rollback_errors.append(f"{target}: {rollback_error}")
        message = f"Restoration apply failed and was rolled back: {error}."
        if rollback_errors:
            message = (
                f"Restoration apply failed: {error}; rollback also failed for "
                + "; ".join(rollback_errors)
                + "."
            )
        raise RestoreError(
            "restore.apply-failed",
            message,
            "Fix repository filesystem permissions and rerun the current preview.",
        ) from error
    finally:
        for path in reversed(cleanup_paths):
            try:
                if path.is_dir() and not path.is_symlink():
                    shutil.rmtree(path)
                elif path.exists() or path.is_symlink():
                    path.unlink()
            except OSError:
                pass


def assert_restore_preimage(repo: Path, plan: RestorePlan) -> None:
    limits = RestoreLimits()
    for skill in plan.skills:
        target = safe_restore_target(repo, skill.target)
        current_files = inspect_restore_target(target, limits)
        current_digest = (
            portable_file_digest(
                (item.path.as_posix().encode("utf-8"), item.content)
                for item in current_files
            )
            if target.exists()
            else None
        )
        if current_digest != skill.observed_digest:
            raise RestoreError(
                "plan.confirmation.stale",
                f"Restoration target changed after planning: {skill.target.as_posix()}.",
                "Preview the current Change Plan and confirm its new planDigest.",
                exit_code=3,
            )
    lock_path = repo / "skills-lock.json"
    try:
        current_lock = lock_path.read_bytes()
    except FileNotFoundError:
        current_lock = None
    except OSError as error:
        raise RestoreError(
            "lock.read-failed",
            f"Could not re-read skills-lock.json before apply: {error}.",
            "Fix repository permissions and preview restoration again.",
        ) from error
    if current_lock != plan.lock_before:
        raise RestoreError(
            "plan.confirmation.stale",
            "skills-lock.json changed after restoration planning.",
            "Preview the current Change Plan and confirm its new planDigest.",
            exit_code=3,
        )


def safe_restore_target(repo: Path, relative: Path) -> Path:
    if relative.is_absolute() or ".." in relative.parts:
        raise RestoreError(
            "target.unsafe-path",
            f"Restoration target is not repository-relative: {relative}.",
            "Fix the bundled setup snapshot target before retrying.",
        )
    current = repo
    for part in relative.parts[:-1]:
        current = current / part
        if current.is_symlink():
            raise RestoreError(
                "target.unsafe-path",
                f"Restoration target parent is a symbolic link: {current}.",
                "Replace the linked parent with a repository directory before retrying.",
            )
    target = repo / relative
    try:
        target.parent.resolve(strict=False).relative_to(repo.resolve(strict=False))
    except ValueError as error:
        raise RestoreError(
            "target.unsafe-path",
            f"Restoration target escapes the repository: {relative}.",
            "Fix the repository path before retrying.",
        ) from error
    return target


def run_restore_skills_command(options: argparse.Namespace) -> int:
    repo = resolve_repo(options.repo)
    if not repo.is_dir():
        return render_restore_error(
            RestoreError(
                "restore.repo-invalid",
                f"Repository root is not a directory: {repo}.",
                "Pass an existing repository directory with --repo.",
                exit_code=2,
            ),
            options.format,
        )
    if options.confirm_plan is not None and re.fullmatch(
        r"[0-9a-f]{64}", options.confirm_plan
    ) is None:
        return render_restore_error(
            RestoreError(
                "plan.confirmation.invalid",
                "Plan confirmation must be a lowercase SHA-256 digest.",
                "Pass the exact planDigest returned by restore-skills preview.",
                exit_code=2,
            ),
            options.format,
            profile_id=options.profile,
        )

    skill_root = Path(__file__).resolve().parents[1]
    try:
        catalog = load_asset_catalog(skill_root)
    except AssetValidationError as error:
        return render_restore_error(
            RestoreError(
                "restore.assets-invalid",
                "Bundled setup assets are invalid: " + "; ".join(error.diagnostics) + ".",
                "Fix the canonical setup-context-driven assets before restoring skills.",
                exit_code=2,
            ),
            options.format,
            profile_id=options.profile,
        )

    source_dir = Path(options.source_dir).expanduser() if options.source_dir else None
    if source_dir is not None and not source_dir.is_absolute():
        source_dir = Path.cwd() / source_dir
    try:
        plan = build_restore_plan(
            repo,
            catalog,
            options.profile,
            options.skill,
            source_dir,
        )
    except RestoreError as error:
        return render_restore_error(
            error,
            options.format,
            profile_id=options.profile,
        )

    if not plan.has_changes:
        render_restore_plan(plan, options.format, ok=True, applied=False)
        return 0
    if options.confirm_plan != plan.digest:
        code = (
            "plan.confirmation.required"
            if options.confirm_plan is None
            else "plan.confirmation.stale"
        )
        message = (
            "Restoration requires confirmation of this exact Change Plan."
            if options.confirm_plan is None
            else "The supplied confirmation does not match the current restoration Change Plan."
        )
        finding_item = {
            "code": code,
            "message": message,
            "action": "Review plannedChanges and rerun with --confirm-plan planDigest.",
        }
        render_restore_plan(
            plan,
            options.format,
            ok=False,
            applied=False,
            finding_item=finding_item,
        )
        return 3
    try:
        apply_restore_plan(repo, plan)
    except RestoreError as error:
        render_restore_plan(
            plan,
            options.format,
            ok=False,
            applied=False,
            finding_item=error.to_json(),
        )
        return error.exit_code
    render_restore_plan(
        plan,
        options.format,
        ok=True,
        applied=True,
        finding_item={
            "code": "restore.completed",
            "message": "Selected external Repository Skill Set members match their immutable snapshots.",
            "action": "Run setup-context-driven audit to verify the selected profile.",
        },
    )
    return 0


def render_restore_error(
    error: RestoreError,
    output_format: str,
    *,
    profile_id: str | None = None,
) -> int:
    payload = {
        "schemaVersion": RESTORE_SCHEMA_VERSION,
        "ok": False,
        "applied": False,
        "profile": profile_id,
        "setup": None,
        "acquisitions": [],
        "skills": [],
        "plannedChanges": [],
        "planDigest": None,
        "finding": error.to_json(),
    }
    render_restore_payload(payload, output_format)
    return error.exit_code


def render_restore_plan(
    plan: RestorePlan,
    output_format: str,
    *,
    ok: bool,
    applied: bool,
    finding_item: dict[str, object] | None = None,
) -> None:
    render_restore_payload(
        plan.to_json(ok=ok, applied=applied, finding_item=finding_item),
        output_format,
    )


def render_restore_payload(payload: dict[str, object], output_format: str) -> None:
    if output_format == "json":
        print(json.dumps(payload, indent=2, sort_keys=False))
        return
    finding_item = payload.get("finding")
    if payload.get("applied"):
        heading = "setup-context-driven restore-skills: applied"
    elif payload.get("ok"):
        heading = "setup-context-driven restore-skills: no changes"
    else:
        heading = "setup-context-driven restore-skills: blocked"
    lines = [heading]
    if isinstance(finding_item, dict):
        lines.append(f"{finding_item.get('code')}: {finding_item.get('message')}")
        lines.append(f"action: {finding_item.get('action')}")
    plan_digest = payload.get("planDigest")
    if plan_digest:
        lines.append(f"planDigest: {plan_digest}")
    for change in payload.get("plannedChanges", []):
        if isinstance(change, dict):
            lines.append(
                f"- {change.get('action')} {change.get('path')} [{change.get('skill')}]"
            )
    print("\n".join(lines))


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
    if args and args[0] == "restore-skills":
        parser = restore_skills_parser()
        return run_restore_skills_command(parser.parse_args(args[1:]))
    if args and args[0] == "audit":
        args = args[1:]

    parser = audit_parser()
    options = parser.parse_args(args)
    return run_audit_command(options)


def load_structured_decision_files(
    paths: list[str],
) -> tuple[StructuredDecisionDocument | None, list[Finding]]:
    documents: list[StructuredDecisionDocument] = []
    diagnostics: list[DecisionDocumentDiagnostic] = []
    for path in paths:
        try:
            documents.append(load_decision_document(path))
        except DecisionDocumentError as error:
            diagnostics.extend(error.diagnostics)
    if not diagnostics:
        try:
            merged = merge_decision_documents(tuple(documents))
        except DecisionDocumentError as error:
            diagnostics.extend(error.diagnostics)
        else:
            return merged, []
    return None, [decision_document_finding(item) for item in diagnostics]


def decision_document_finding(diagnostic: DecisionDocumentDiagnostic) -> Finding:
    return finding(
        diagnostic.code,
        "error",
        diagnostic.path.as_posix() or ".",
        diagnostic.item_id,
        diagnostic.message,
        "Correct the named structured decision entry and rerun without changing repository bytes.",
    )


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

    decision_document, decision_file_findings = load_structured_decision_files(
        options.decision_file
    )
    if decision_file_findings:
        render_result(AuditResult(sorted_findings(decision_file_findings)), options.format)
        return 2

    if options.decision or decision_document is not None:
        readoption_result, readoption_invalid = audit_readoption_repository(
            repo,
            catalog,
            decision_document=decision_document,
            decision_args=options.decision,
            profile_override=options.profile,
        )
        if readoption_result is not None:
            render_result(readoption_result, options.format)
            return exit_code_for(readoption_result, readoption_invalid)
        if decision_document is not None and decision_document.readoption is not None:
            unexpected = AuditResult(
                [
                    finding(
                        "readoption.source.unexpected",
                        "error",
                        decision_document.source_paths[0].as_posix(),
                        decision_document.readoption.source_baseline_id,
                        "Decision file contains Readoption data for a compatible repository.",
                        "Remove the stale readoption section and rerun audit.",
                    )
                ],
                decision_document=decision_document.to_json(),
            )
            render_result(unexpected, options.format)
            return 2
        result, invalid_input, _ = plan_apply(
            repo=repo,
            catalog=catalog,
            profile_override=options.profile,
            decision_args=options.decision,
            decision_document=decision_document,
        )
        render_result(result, options.format)
        return exit_code_for(result, invalid_input)

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

    decision_document, decision_file_findings = load_structured_decision_files(
        options.decision_file
    )
    if decision_file_findings:
        render_result(AuditResult(sorted_findings(decision_file_findings)), options.format)
        return 2

    if options.confirm_plan is not None and re.fullmatch(
        r"[0-9a-f]{64}", options.confirm_plan
    ) is None:
        malformed = AuditResult(
            [
                finding(
                    "plan.confirmation.invalid",
                    "error",
                    ".",
                    "apply",
                    "Plan confirmation must be a lowercase SHA-256 digest.",
                    "Pass the exact planDigest returned by audit or apply preview.",
                )
            ]
        )
        render_result(malformed, options.format)
        return 2

    readoption_apply = False
    if decision_document is not None and decision_document.readoption is not None:
        readoption_result, readoption_invalid = audit_readoption_repository(
            repo,
            catalog,
            decision_document=decision_document,
            decision_args=options.decision,
            profile_override=options.profile,
        )
        if readoption_result is not None:
            if (
                readoption_result.summary["errors"]
                or readoption_result.summary["decisions"]
                or readoption_invalid
            ):
                render_result(readoption_result, options.format)
                return exit_code_for(readoption_result, readoption_invalid)
            result = readoption_result
            invalid_input = False
            manifest = manifest_from_mutations(readoption_result.planned_outputs)
            plan = ChangePlan(
                "baseline-readoption",
                readoption_result.planned_outputs,
                readoption_result.plan_digest,
                manifest,
            )
            readoption_apply = True
    if not readoption_apply:
        result, invalid_input, plan = plan_apply(
            repo=repo,
            catalog=catalog,
            profile_override=options.profile,
            decision_args=options.decision,
            decision_document=decision_document,
        )
    if result.summary["errors"] or result.summary["decisions"] or invalid_input:
        render_result(result, options.format)
        return exit_code_for(result, invalid_input)
    nonblocking_findings = [
        item for item in result.findings if item.severity in {"warning", "info"}
    ]

    confirmation = options.confirm_plan
    has_writes = any(
        mutation.before_digest != mutation.after_digest for mutation in plan.mutations
    )
    if has_writes and confirmation != plan.digest:
        code = "plan.confirmation.required" if confirmation is None else "plan.confirmation.stale"
        message = (
            "Apply requires confirmation of this exact Change Plan."
            if confirmation is None
            else "The supplied confirmation does not match the current Change Plan."
        )
        blocked = AuditResult(
            sorted_findings(
                nonblocking_findings
                + [finding(
                    code,
                    "decision",
                    ".",
                    "apply",
                    message,
                    "Review plannedChanges and rerun with --confirm-plan planDigest.",
                )]
            ),
            selection=result.selection,
            planned_changes=result.planned_changes,
            plan_digest=result.plan_digest,
            retention=result.retention,
            source_baseline=result.source_baseline,
            decision_document=result.decision_document,
            planned_outputs=result.planned_outputs,
            capabilities=result.capabilities,
            setup_snapshot=result.setup_snapshot,
            verification=result.verification,
        )
        render_result(blocked, options.format)
        return 3

    if not has_writes:
        empty = AuditResult(
            sorted_findings(
                nonblocking_findings
                + [finding(
                    "managed.apply.empty",
                    "info",
                    ".",
                    "apply",
                    "The repository already matches the selected Change Plan.",
                    "No action needed.",
                )]
            ),
            selection=result.selection,
            planned_changes=(),
            plan_digest=result.plan_digest,
            retention=result.retention,
            source_baseline=result.source_baseline,
            decision_document=result.decision_document,
            capabilities=result.capabilities,
            setup_snapshot=result.setup_snapshot,
            verification=result.verification,
        )
        render_result(empty, options.format)
        return 0

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
            ),
            retention=result.retention,
            source_baseline=result.source_baseline,
            decision_document=result.decision_document,
            capabilities=result.capabilities,
            setup_snapshot=result.setup_snapshot,
            verification=result.verification,
        )
        render_result(failure, options.format)
        return 1

    applied = AuditResult(
        sorted_findings(
            nonblocking_findings
            + [finding(
                "managed.apply.completed",
                "info",
                ".",
                "apply",
                "Managed setup content matches the selected profile.",
                "No action needed.",
            )]
        ),
        selection=result.selection,
        planned_changes=result.planned_changes,
        plan_digest=result.plan_digest,
        retention=result.retention,
        source_baseline=result.source_baseline,
        decision_document=result.decision_document,
        planned_outputs=result.planned_outputs,
        capabilities=result.capabilities,
        setup_snapshot=result.setup_snapshot,
        verification=result.verification,
    )
    render_result(applied, options.format)
    return 0


def manifest_from_mutations(mutations: tuple[FileMutation, ...]) -> dict:
    for mutation in mutations:
        if mutation.path != MANIFEST_PATH or mutation.content is None:
            continue
        try:
            manifest = json.loads(mutation.content)
        except (UnicodeDecodeError, json.JSONDecodeError):
            return {}
        return manifest if isinstance(manifest, dict) else {}
    return {}


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


def delegation_findings(
    repo: Path,
    catalog: AssetCatalog,
    profile_id: str,
    active_modules: Iterable[str] | None = None,
    limits: DelegationScanLimits | None = None,
) -> list[Finding]:
    active_limits = limits or DelegationScanLimits()
    documents, scan_finding = discover_delegation_documents(repo, active_limits)
    if scan_finding is not None:
        return [scan_finding]

    selected_coverage = selected_clause_coverage(
        catalog,
        profile_id,
        active_modules,
    )
    findings: list[Finding] = []
    seen: set[tuple[str, str]] = set()
    for relative_path, content in documents:
        authored = repository_authored_instruction_text(relative_path, content)
        for paragraph in re.split(r"\n\s*\n", authored):
            if not DELEGATION_SIGNAL.search(paragraph):
                continue
            if not DELEGATION_TARGET.search(paragraph):
                continue
            for coverage_id, contract in sorted(catalog.coverage_contracts.items()):
                if coverage_id in selected_coverage:
                    continue
                if not any(
                    contains_delegation_alias(paragraph, alias)
                    for alias in contract.delegation_aliases
                ):
                    continue
                key = (relative_path.as_posix(), coverage_id)
                if key in seen:
                    continue
                seen.add(key)
                findings.append(
                    finding(
                        "delegation.baseline-floor",
                        "info",
                        relative_path.as_posix(),
                        coverage_id,
                        (
                            f"Repository-authored instructions delegate {coverage_id}, "
                            f"which profile {profile_id} does not cover."
                        ),
                        (
                            "Treat the generated baseline as a floor; preserve the "
                            "repository-owned guidance and select or add coverage deliberately."
                        ),
                    )
                )
    return sorted_findings(findings)


def selected_clause_coverage(
    catalog: AssetCatalog,
    profile_id: str,
    active_modules: Iterable[str] | None,
) -> set[str]:
    if active_modules is None:
        rule_ids = catalog.profiles[profile_id].get("requiredRules", [])
    else:
        rule_ids = [
            rule.get("id")
            for module_id in active_modules
            for rule in catalog.modules[module_id].get("rules", [])
            if isinstance(rule, dict)
        ]
    return {
        coverage_id
        for rule_id in rule_ids
        if isinstance(rule_id, str) and rule_id in catalog.rule_contracts
        for coverage_id in catalog.rule_contracts[rule_id].coverage
    }


def discover_delegation_documents(
    repo: Path,
    limits: DelegationScanLimits,
) -> tuple[tuple[tuple[Path, str], ...], Finding | None]:
    if limits.max_files < 1 or limits.max_bytes < 1:
        return (), delegation_scan_limit_finding(
            "Delegation scan limits must be positive."
        )

    documents: list[tuple[Path, str]] = []
    total_bytes = 0
    pending: list[tuple[Path, Path]] = [(repo, Path())]
    while pending:
        absolute_directory, relative_directory = pending.pop()
        try:
            with os.scandir(absolute_directory) as iterator:
                entries = sorted(iterator, key=lambda item: item.name)
        except OSError as error:
            return (), delegation_scan_read_finding(relative_directory, error)

        nested: list[tuple[Path, Path]] = []
        for entry in entries:
            relative_path = relative_directory / entry.name
            try:
                if entry.is_symlink():
                    continue
                if entry.is_dir(follow_symlinks=False):
                    if delegation_directory_is_ignored(relative_path):
                        continue
                    nested.append((Path(entry.path), relative_path))
                    continue
                if entry.name not in DELEGATION_DOCUMENT_NAMES:
                    continue
                if not entry.is_file(follow_symlinks=False):
                    continue
            except OSError as error:
                return (), delegation_scan_read_finding(relative_path, error)

            if len(documents) >= limits.max_files:
                return (), delegation_scan_limit_finding(
                    f"Delegation scan exceeds the {limits.max_files}-file limit."
                )
            remaining_bytes = limits.max_bytes - total_bytes
            try:
                data = read_instruction_bytes(Path(entry.path), remaining_bytes)
            except OSError as error:
                return (), delegation_scan_read_finding(relative_path, error)
            if data is None:
                return (), delegation_scan_limit_finding(
                    f"Delegation scan exceeds the {limits.max_bytes}-byte limit."
                )
            try:
                content = data.decode("utf-8")
            except UnicodeDecodeError:
                return (), delegation_scan_read_finding(
                    relative_path,
                    ValueError("instruction document is not UTF-8"),
                )
            total_bytes += len(data)
            documents.append((relative_path, content))
        pending.extend(reversed(nested))

    return tuple(sorted(documents, key=lambda item: item[0].as_posix())), None


def delegation_directory_is_ignored(relative_path: Path) -> bool:
    if relative_path.name in DELEGATION_IGNORED_DIRECTORIES:
        return True
    parts = relative_path.parts
    return any(parts[: len(prefix)] == prefix for prefix in DELEGATION_IGNORED_PREFIXES)


def read_instruction_bytes(path: Path, max_bytes: int) -> bytes | None:
    flags = os.O_RDONLY | getattr(os, "O_NOFOLLOW", 0)
    descriptor = os.open(path, flags)
    try:
        if not stat.S_ISREG(os.fstat(descriptor).st_mode):
            return b""
        chunks: list[bytes] = []
        total = 0
        while True:
            chunk = os.read(descriptor, min(64 * 1024, max_bytes - total + 1))
            if not chunk:
                return b"".join(chunks)
            chunks.append(chunk)
            total += len(chunk)
            if total > max_bytes:
                return None
    finally:
        os.close(descriptor)


def repository_authored_instruction_text(relative_path: Path, content: str) -> str:
    spans, _ = parse_managed_block_spans(relative_path, content)
    if not spans:
        return content
    characters = list(content)
    for span in spans.values():
        for index in range(span.start, span.end):
            if characters[index] not in {"\n", "\r"}:
                characters[index] = " "
    return "".join(characters)


def contains_delegation_alias(paragraph: str, alias: str) -> bool:
    return re.search(
        rf"(?<![A-Za-z0-9_]){re.escape(alias)}(?![A-Za-z0-9_])",
        paragraph,
        re.IGNORECASE,
    ) is not None


def delegation_scan_limit_finding(message: str) -> Finding:
    return finding(
        "delegation.scan-limit",
        "info",
        ".",
        "delegation.scan",
        message,
        "Review repository instruction delegation directly; the baseline remains a floor.",
    )


def delegation_scan_read_finding(relative_path: Path, error: Exception) -> Finding:
    path = relative_path.as_posix() or "."
    return finding(
        "delegation.scan-unreadable",
        "info",
        path,
        "delegation.scan",
        f"Delegation scan could not read {path}: {error}.",
        "Review that instruction document directly; the baseline remains a floor.",
    )


def audit_repository(
    repo: Path,
    catalog: AssetCatalog,
    profile_override: str | None = None,
    show_extra_skills: bool = False,
    setups_dir: str | None = None,
) -> tuple[AuditResult, bool]:
    readoption_result, readoption_invalid = audit_readoption_repository(repo, catalog)
    if readoption_result is not None:
        return readoption_result, readoption_invalid

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

    findings.extend(
        validate_decision_plan_references(
            repo,
            catalog,
            decision_plan,
            existing_manifest=manifest,
        )
    )

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
    validate_manifest_artifacts(repo, manifest, expected_artifacts, findings)
    validate_documents(repo, expected_artifacts, findings)
    validate_secondbrain_documents(repo, ordered_modules, findings)

    plan_result, plan_invalid_input, _ = plan_apply(
        repo=repo,
        catalog=catalog,
        profile_override=profile_id,
        decision_args=[],
    )
    existing_finding_keys = {
        (item.code, item.path, item.managed_id, item.message) for item in findings
    }
    findings.extend(
        item
        for item in plan_result.findings
        if (item.code, item.path, item.managed_id, item.message) not in existing_finding_keys
    )
    return (
        AuditResult(
            sorted_findings(findings),
            selection=plan_result.selection,
            planned_changes=plan_result.planned_changes,
            plan_digest=plan_result.plan_digest,
            retention=plan_result.retention,
        ),
        invalid_input or plan_invalid_input,
    )


def audit_readoption_repository(
    repo: Path,
    catalog: AssetCatalog,
    *,
    decision_document: StructuredDecisionDocument | None = None,
    decision_args: list[str] | None = None,
    profile_override: str | None = None,
) -> tuple[AuditResult | None, bool]:
    manifest_findings: list[Finding] = []
    manifest, invalid_input = load_manifest(repo, manifest_findings)
    if invalid_input:
        return None, invalid_input

    declared_identity: str | None = None
    if manifest is None:
        declared_identity = "manifest.missing"
    elif is_current_strict_manifest(manifest, catalog):
        return None, False
    else:
        _, transition_findings = resolve_upgrade_transition(catalog, manifest)
        if any(
            item.code == "retention.baseline.unknown"
            for item in transition_findings
        ):
            generator = manifest.get("generator")
            baseline = (
                generator.get("baseline") if isinstance(generator, dict) else None
            )
            if isinstance(baseline, str) and baseline:
                declared_identity = baseline
            else:
                fingerprint = legacy_manifest_fingerprint(manifest)
                declared_identity = (
                    f"manifest.fingerprint.{fingerprint}"
                    if fingerprint is not None
                    else "manifest.incompatible"
                )

    if declared_identity is None:
        return None, False

    try:
        inventory = inventory_incompatible_source_baseline(repo, declared_identity)
    except SourceInventoryError as error:
        findings = [
            finding(
                diagnostic.code,
                "error",
                diagnostic.path.as_posix() or ".",
                "source-baseline",
                diagnostic.message,
                "Replace the unsafe carrier with a bounded regular file and rerun audit.",
            )
            for diagnostic in error.diagnostics
        ]
        return AuditResult(sorted_findings(findings)), False

    if manifest is None and not inventory.carriers:
        return None, False

    findings: list[Finding] = [
        finding(
            "readoption.baseline.incompatible",
            "info",
            str(MANIFEST_PATH),
            inventory.baseline_id,
            (
                f"Source Baseline {declared_identity!r} is incompatible; "
                f"audit inventoried {len(inventory.entries)} structural entries."
            ),
            "Review every sourceEntries item before supplying Readoption dispositions.",
        )
    ]
    if decision_document is None or decision_document.readoption is None:
        findings.extend(
            finding(
                "readoption.disposition.required",
                "decision",
                entry.path.as_posix(),
                entry.entry_id,
                "Source Baseline Entry has no Readoption disposition.",
                "Supply one explicit classification and disposition for this entry.",
            )
            for entry in inventory.entries
        )
        return (
            AuditResult(
                sorted_findings(findings),
                source_baseline=inventory,
                decision_document=(
                    decision_document.to_json()
                    if decision_document is not None
                    else None
                ),
            ),
            False,
        )

    cli_decisions, cli_findings = parse_decision_args(decision_args or [], catalog)
    file_decisions, file_findings = structured_decision_values(
        decision_document, catalog
    )
    findings.extend(cli_findings)
    findings.extend(file_findings)
    for decision_id in sorted(set(cli_decisions).intersection(file_decisions)):
        if cli_decisions[decision_id] == file_decisions[decision_id]:
            continue
        findings.append(
            finding(
                "decision-file.decision.conflict",
                "error",
                decision_document.source_paths[0].as_posix(),
                decision_id,
                "--decision conflicts with the structured decision-file value.",
                "Keep the decision in one input or make both values identical.",
            )
        )

    dispositions, disposition_diagnostics = validate_readoption_decisions(
        repo,
        inventory,
        decision_document.readoption,
    )
    findings.extend(
        decision_document_finding(item) for item in disposition_diagnostics
    )
    managed_ids = catalog_managed_entry_ids(catalog)
    for item in dispositions:
        if item.disposition != "managed-entry" or item.destination is None:
            continue
        managed_id = item.destination.get("managedId")
        if managed_id not in managed_ids:
            findings.append(
                finding(
                    "readoption.destination.managed-entry.unknown",
                    "error",
                    decision_document.source_paths[0].as_posix(),
                    item.entry_id,
                    f"Managed destination {managed_id!r} is not in the current catalog.",
                    "Choose a current managed entry or another typed disposition.",
                )
            )

    normalized_document = decision_document.to_json()
    normalized_readoption = dict(normalized_document["readoption"])
    normalized_readoption["dispositions"] = [item.to_json() for item in dispositions]
    normalized_document["readoption"] = normalized_readoption
    if any(item.severity == "error" for item in findings):
        return (
            AuditResult(
                sorted_findings(findings),
                source_baseline=inventory,
                decision_document=normalized_document,
            ),
            True,
        )

    if profile_override is not None:
        readoption_context, context_findings = prepare_readoption_plan_context(
            repo=repo,
            catalog=catalog,
            profile_id=profile_override,
            inventory=inventory,
            dispositions=dispositions,
            decision_document=decision_document,
            normalized_document=normalized_document,
            cli_decisions=cli_decisions,
            file_decisions=file_decisions,
        )
        findings.extend(context_findings)
        if readoption_context is None:
            return (
                AuditResult(
                    sorted_findings(findings),
                    source_baseline=inventory,
                    decision_document=normalized_document,
                ),
                any(item.severity == "error" for item in context_findings),
            )
        if any(item.severity == "error" for item in context_findings):
            return (
                AuditResult(
                    sorted_findings(findings),
                    source_baseline=inventory,
                    decision_document=normalized_document,
                    capabilities=readoption_context.capabilities,
                    setup_snapshot=readoption_context.setup_snapshot,
                    verification=readoption_context.verification,
                ),
                False,
            )
        plan_result, plan_invalid, plan = plan_apply(
            repo=repo,
            catalog=catalog,
            profile_override=profile_override,
            decision_args=decision_args or [],
            decision_document=decision_document,
            readoption=readoption_context,
        )
        return (
            AuditResult(
                sorted_findings(findings + plan_result.findings),
                selection=plan_result.selection,
                planned_changes=plan_result.planned_changes,
                plan_digest=plan_result.plan_digest,
                source_baseline=inventory,
                decision_document=normalized_document,
                planned_outputs=plan.mutations,
                capabilities=readoption_context.capabilities,
                setup_snapshot=readoption_context.setup_snapshot,
                verification=readoption_context.verification,
            ),
            plan_invalid,
        )

    planned_changes: tuple[PlannedChange, ...] = ()
    proposed_rules = repository_rules_proposed_bytes(dispositions)
    repository_rules_path = Path("docs/agents/repository-rules.md")
    if proposed_rules:
        target = repo / repository_rules_path
        if target.exists() or target.is_symlink():
            try:
                target_stat = target.lstat()
            except OSError as error:
                findings.append(
                    finding(
                        "readoption.repository-rules.read",
                        "error",
                        repository_rules_path.as_posix(),
                        "repository-rules",
                        f"Repository-Specific Normative Rules cannot be inspected: {error}.",
                        "Replace the target with a regular repository-owned file and rerun.",
                    )
                )
            else:
                if not stat.S_ISREG(target_stat.st_mode) or stat.S_ISLNK(
                    target_stat.st_mode
                ):
                    findings.append(
                        finding(
                            "readoption.repository-rules.type.invalid",
                            "error",
                            repository_rules_path.as_posix(),
                            "repository-rules",
                            "Repository-Specific Normative Rules must be a regular non-symlink file.",
                            "Replace the target with a regular repository-owned file and rerun.",
                        )
                    )
        else:
            planned_changes = (
                PlannedChange(
                    action="create repository-specific normative rules",
                    path=repository_rules_path,
                    managed_id="repository-rules.readoption",
                    state="definite",
                    reason="Explicit Readoption dispositions propose exact repository-owned bytes.",
                    before_digest=None,
                    after_digest=hashlib.sha256(proposed_rules).hexdigest(),
                ),
            )

    if any(item.severity == "error" for item in findings):
        return (
            AuditResult(
                sorted_findings(findings),
                source_baseline=inventory,
                decision_document=normalized_document,
            ),
            True,
        )

    digest_payload = {
        "kind": "baseline-readoption-decisions",
        "sourceBaseline": inventory.identity_json(),
        "sourceEntries": [entry.to_json() for entry in inventory.entries],
        "decisionDocument": normalized_document,
        "decisionFileDigests": list(decision_document.source_digests),
        "operations": [item.to_json() for item in planned_changes],
    }
    plan_digest = hashlib.sha256(canonical_json_bytes(digest_payload)).hexdigest()
    findings.append(
        finding(
            "readoption.dispositions.resolved",
            "info",
            str(MANIFEST_PATH),
            inventory.baseline_id,
            f"Every one of {len(inventory.entries)} Source Baseline Entries has one explicit disposition.",
            "Review the normalized decision document and deterministic plan digest.",
        )
    )
    if planned_changes:
        findings.append(
            finding(
                "plan.confirmation.required",
                "decision",
                repository_rules_path.as_posix(),
                "repository-rules.readoption",
                "Creating Repository-Specific Normative Rules requires confirmation of this exact plan.",
                "Review the proposed bytes and rerun apply with --confirm-plan planDigest.",
            )
        )
    return (
        AuditResult(
            sorted_findings(findings),
            source_baseline=inventory,
            planned_changes=planned_changes,
            plan_digest=plan_digest,
            decision_document=normalized_document,
        ),
        False,
    )


def prepare_readoption_plan_context(
    repo: Path,
    catalog: AssetCatalog,
    profile_id: str,
    inventory: IncompatibleSourceBaseline,
    dispositions: tuple[ReadoptionDisposition, ...],
    decision_document: StructuredDecisionDocument,
    normalized_document: dict[str, object],
    cli_decisions: dict[str, object],
    file_decisions: dict[str, object],
) -> tuple[ReadoptionPlanContext | None, list[Finding]]:
    findings: list[Finding] = []
    contract = catalog.standard_profiles.get(profile_id)
    if contract is None:
        findings.append(
            finding(
                "profile.readoption.unsupported",
                "error",
                str(MANIFEST_PATH),
                profile_id,
                "Baseline Readoption requires a strict 0.0.1 profile.",
                "Select a maintained 0.0.1 profile for this generation.",
            )
        )
        return None, findings

    decisions = {**file_decisions, **cli_decisions}
    decision_plan = resolve_decision_plan(catalog, profile_id, None, decisions)
    findings.extend(decision_required_findings(decision_plan))
    if decision_plan.unresolved_decisions:
        return None, findings

    http_answer = (
        decision_plan.resolved_decisions.get(contract.http_decision_id)
        if contract.http_decision_id
        else None
    )
    http_value = http_answer.get("value") if isinstance(http_answer, dict) else None
    try:
        profile_plan = build_standard_profile_plan(catalog, http_value, profile_id)
        setup_snapshot = json.loads(render_standard_profile_snapshot(profile_plan))
    except (TypeError, ValueError, json.JSONDecodeError) as error:
        findings.append(
            finding(
                "profile.snapshot.invalid",
                "error",
                str(MANIFEST_PATH),
                profile_id,
                f"Strict Setup Snapshot could not be resolved: {error}.",
                "Supply a valid typed HTTP Contract Decision and rerun preview.",
            )
        )
        return None, findings

    evaluation = evaluate_repository_capabilities(
        repo,
        (*UNIVERSAL_CAPABILITIES, *contract.capabilities),
    )
    capability_document = json.loads(render_capability_json(evaluation))
    capabilities = tuple(capability_document["capabilities"])
    findings.extend(capability_findings(evaluation))

    expected_paths = {
        artifact.path for artifact in expected_artifacts_for_plan(decision_plan)
    }
    for disposition in dispositions:
        if disposition.disposition != "repository-document" or disposition.destination is None:
            continue
        target_path = Path(str(disposition.destination["path"]))
        if target_path in expected_paths:
            findings.append(
                finding(
                    "readoption.destination.managed-collision",
                    "error",
                    target_path.as_posix(),
                    disposition.entry_id,
                    "A typed repository document cannot also be a setup-managed output.",
                    "Choose a distinct typed repository document or a managed-entry disposition.",
                )
            )
    verification = tuple(
        {
            "id": entry.entry_id,
            "kind": entry.kind,
            "tool": entry.tool,
            "command": entry.command,
        }
        for entry in contract.verification
    )
    context = ReadoptionPlanContext(
        source_baseline=inventory,
        dispositions=dispositions,
        decision_document=normalized_document,
        decision_file_digests=decision_document.source_digests,
        capabilities=capabilities,
        setup_snapshot=setup_snapshot,
        verification=verification,
        repository_rules=repository_rules_proposed_bytes(dispositions),
    )
    return context, findings


def capability_findings(evaluation: CapabilityEvaluation) -> list[Finding]:
    findings: list[Finding] = []
    for outcome in evaluation.outcomes:
        if outcome.diagnostic.code == "capability.satisfied":
            continue
        if outcome.blocking:
            severity = "error"
        elif outcome.capability.strength is RequirementStrength.RECOMMENDED:
            severity = "warning"
        else:
            continue
        findings.append(
            finding(
                outcome.diagnostic.code,
                severity,
                (
                    outcome.evidence[0].source_path.as_posix()
                    if outcome.evidence and outcome.evidence[0].source_path is not None
                    else "."
                ),
                outcome.capability.capability_id,
                outcome.diagnostic.message,
                outcome.diagnostic.next_action,
            )
        )
    return findings


def is_current_strict_manifest(manifest: dict, catalog: AssetCatalog) -> bool:
    profile_id = manifest.get("profile")
    generator = manifest.get("generator")
    return (
        manifest.get("schemaVersion") == MANIFEST_SCHEMA_VERSION_0_0_1
        and manifest.get("version") == OWNED_VERSION_0_0_1
        and isinstance(profile_id, str)
        and profile_id in catalog.standard_profiles
        and isinstance(generator, dict)
        and generator.get("version") == OWNED_VERSION_0_0_1
        and generator.get("baseline") == f"baseline.{profile_id}-0.0.1"
    )


def catalog_managed_entry_ids(catalog: AssetCatalog) -> set[str]:
    managed_ids: set[str] = set()
    for module in catalog.modules.values():
        for collection_name in ("rootBlocks", "supportingGuides"):
            for item in module.get(collection_name, []):
                if isinstance(item, dict) and isinstance(item.get("id"), str):
                    managed_ids.add(item["id"])
    return managed_ids


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
    setup_id = catalog.profiles[profile_id]["setup"]
    setup = catalog.setups[setup_id]
    required_skills = [
        skill
        for skill in setup.get("skills", [])
        if isinstance(skill.get("name"), str) and skill["name"]
    ]
    required_names = [skill["name"] for skill in required_skills]
    external_contracts = {
        contract.skill_name: contract
        for contract in catalog.external_sources_by_setup.get(setup_id, ())
    }
    for skill in required_skills:
        skill_name = skill["name"]
        installed_skill = installed.get(skill_name)
        external_contract = external_contracts.get(skill_name)
        if installed_skill is None:
            remediation = (
                external_skill_remediation(profile_id, external_contract)
                if external_contract is not None
                else None
            )
            findings.append(
                finding(
                    "skills.required.missing",
                    "error",
                    f".agents/skills/{skill_name}",
                    f"profile.{profile_id}",
                    f"Required skill {skill_name} is not installed.",
                    (
                        "Use remediation.previewArgv to preview exact immutable restoration; audit made no changes."
                        if remediation is not None
                        else f"Install the {setup_id} canonical skill setup or add .agents/skills/{skill_name}/SKILL.md."
                    ),
                    remediation=remediation,
                )
            )
            continue
        if external_contract is not None:
            drift_reason = installed_skill.unsafe_error
            if drift_reason is None:
                try:
                    actual_digest = portable_tree_digest(repo / installed_skill.root)
                except PortableTreeError as error:
                    drift_reason = str(error)
                else:
                    if actual_digest != external_contract.tree_digest:
                        drift_reason = (
                            f"complete tree digest {actual_digest} does not match "
                            f"{external_contract.tree_digest}"
                        )
            if drift_reason is not None:
                findings.append(
                    finding(
                        "skills.required.drift",
                        "error",
                        installed_skill.root.as_posix(),
                        skill_name,
                        f"Required skill {skill_name} differs from the selected immutable snapshot: {drift_reason}.",
                        "Use remediation.previewArgv to preview exact immutable restoration; audit made no changes.",
                        remediation=external_skill_remediation(
                            profile_id,
                            external_contract,
                        ),
                    )
                )
            continue
        expected_digest = skill.get("contentDigest")
        canonical_skill_file = local_canonical_skill_file(skill_name)
        if (
            canonical_skill_file is not None
            and isinstance(expected_digest, str)
            and (
                installed_skill.unsafe_error is not None
                or installed_skill.content_digest != expected_digest
            )
        ):
            findings.append(
                finding(
                    "skills.required.drift",
                    "error",
                    installed_skill.path.as_posix(),
                    skill_name,
                    f"Required skill {skill_name} content differs from the selected setup snapshot.",
                    f"Refresh .agents/skills/{skill_name}/SKILL.md from the {setup_id} canonical skill setup.",
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


def local_canonical_skill_file(skill_name: str) -> Path | None:
    script_path = Path(__file__).resolve()
    for parent in script_path.parents:
        candidate = parent / ".agents" / "skills" / skill_name / "SKILL.md"
        if candidate.is_file():
            return candidate
    return None


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
    if skills_root.is_symlink() or not skills_root.is_dir():
        return installed
    for skill_root in sorted(skills_root.iterdir(), key=lambda item: item.name.encode("utf-8")):
        skill_name = skill_root.name
        lock_entry = lock_entries.get(skill_name)
        relative_root = skill_root.relative_to(repo)
        if skill_root.is_symlink() or not skill_root.is_dir():
            installed[skill_name] = InstalledSkill(
                name=skill_name,
                path=relative_root,
                root=relative_root,
                locked=lock_entry is not None,
                origin=lock_entry or {},
                content_digest=None,
                unsafe_error="installed skill root is not a regular directory",
            )
            continue
        skill_file = skill_root / "SKILL.md"
        relative_file = skill_file.relative_to(repo)
        if skill_file.is_symlink():
            installed[skill_name] = InstalledSkill(
                name=skill_name,
                path=relative_file,
                root=relative_root,
                locked=lock_entry is not None,
                origin=lock_entry or {},
                content_digest=None,
                unsafe_error="installed SKILL.md is not a regular file",
            )
            continue
        if not skill_file.exists():
            continue
        if not skill_file.is_file():
            installed[skill_name] = InstalledSkill(
                name=skill_name,
                path=relative_file,
                root=relative_root,
                locked=lock_entry is not None,
                origin=lock_entry or {},
                content_digest=None,
                unsafe_error="installed SKILL.md is not a regular file",
            )
            continue
        try:
            content_digest = managed_digest(skill_file.read_text(encoding="utf-8"))
        except (OSError, UnicodeDecodeError) as error:
            content_digest = None
            unsafe_error = f"installed SKILL.md cannot be read as UTF-8: {error}"
        else:
            unsafe_error = None
        installed[skill_name] = InstalledSkill(
            name=skill_name,
            path=relative_file,
            root=relative_root,
            locked=lock_entry is not None,
            origin=lock_entry or {},
            content_digest=content_digest,
            unsafe_error=unsafe_error,
        )
    return installed


def external_skill_remediation(
    profile_id: str,
    contract: ExternalSkillContract,
) -> dict[str, object]:
    source = contract.source
    return {
        "provider": source.provider,
        "skill": contract.skill_name,
        "source": source.repository,
        "ref": source.revision,
        "sourcePath": source.source_path.as_posix(),
        "expectedDigest": contract.tree_digest,
        "previewArgv": [
            "python3",
            ".agents/skills/setup-context-driven/scripts/context_setup.py",
            "restore-skills",
            "--repo",
            ".",
            "--profile",
            profile_id,
            "--skill",
            contract.skill_name,
            "--format",
            "json",
        ],
    }


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
    if not decision_plan.unresolved_decisions:
        findings.extend(
            validate_decision_plan_references(
                repo,
                catalog,
                decision_plan,
                existing_manifest=existing_manifest,
            )
        )
        findings.extend(
            delegation_findings(
                repo,
                catalog,
                profile_id,
                active_modules=decision_plan.active_modules,
            )
        )
    return (
        AuditResult(
            sorted_findings(findings),
            selection=decision_plan.selection,
            planned_changes=planned_changes_for_plan(repo, decision_plan),
        ),
        invalid_input,
    )


def current_baseline_id(catalog: AssetCatalog) -> str | None:
    source_baselines = {
        transition.from_baseline
        for transition in catalog.upgrade_transitions.values()
    }
    target_baselines = {
        transition.to_baseline
        for transition in catalog.upgrade_transitions.values()
    }
    terminal_baselines = target_baselines - source_baselines
    if len(terminal_baselines) != 1:
        return None
    return next(iter(terminal_baselines))


def legacy_manifest_fingerprint(manifest: dict) -> str | None:
    artifacts = manifest.get("managedArtifacts")
    if not isinstance(artifacts, list):
        return None
    normalized: list[dict[str, object]] = []
    for artifact in artifacts:
        if not isinstance(artifact, dict):
            return None
        managed_id = artifact.get("id")
        version = artifact.get("version")
        template = artifact.get("template")
        digest = artifact.get("digest")
        if (
            not isinstance(managed_id, str)
            or not managed_id
            or not isinstance(version, int)
            or isinstance(version, bool)
            or version < 1
            or not isinstance(template, str)
            or not template
            or not isinstance(digest, str)
            or re.fullmatch(r"[0-9a-f]{64}", digest) is None
        ):
            return None
        normalized.append(
            {
                "id": managed_id,
                "version": version,
                "template": template,
                "digest": digest,
            }
        )
    normalized.sort(key=lambda item: str(item["id"]))
    return hashlib.sha256(canonical_json_bytes(normalized)).hexdigest()


def resolve_upgrade_transition(
    catalog: AssetCatalog,
    manifest: dict | None,
) -> tuple[UpgradeTransition | None, list[Finding]]:
    if manifest is None:
        return None, []

    current_baseline = current_baseline_id(catalog)
    if current_baseline is None:
        return None, [
            finding(
                "retention.baseline.catalog-invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "The bundled upgrade catalog has no unique current baseline.",
                "Fix the bundled transition graph before upgrading repositories.",
            )
        ]

    generator = manifest.get("generator")
    declared_baseline = (
        generator.get("baseline") if isinstance(generator, dict) else None
    )
    if declared_baseline is not None:
        if declared_baseline == current_baseline:
            return None, []
        matching = [
            transition
            for transition in catalog.upgrade_transitions.values()
            if transition.from_baseline == declared_baseline
            and transition.to_baseline == current_baseline
        ]
        source_description = f"declared baseline {declared_baseline!r}"
    else:
        fingerprint = legacy_manifest_fingerprint(manifest)
        matching = [
            transition
            for transition in catalog.upgrade_transitions.values()
            if fingerprint in transition.legacy_manifest_fingerprints
        ]
        source_description = (
            f"legacy fingerprint {fingerprint}"
            if fingerprint is not None
            else "an invalid legacy fingerprint"
        )

    if len(matching) == 1:
        return matching[0], []
    return None, [
        finding(
            "retention.baseline.unknown",
            "error",
            str(MANIFEST_PATH),
            "manifest",
            f"The Setup Manifest has {source_description}, which is not a declared source baseline.",
            "Restore an exact supported manifest or add a reviewed transition contract before upgrading.",
        )
    ]


def selected_clause_enforcement(
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
) -> dict[str, str]:
    selected_artifacts = {
        planned.artifact.managed_id
        for planned in decision_plan.artifacts
        if planned.present and planned.state == "definite"
    }
    reachable: dict[str, str] = {}
    for module_id in decision_plan.active_modules:
        module = catalog.modules[module_id]
        for guide in module.get("supportingGuides", []):
            if guide.get("id") not in selected_artifacts:
                continue
            for rule_id in guide.get("rules", []):
                rule = catalog.rule_contracts.get(rule_id)
                if rule is None:
                    continue
                for clause in rule.clauses:
                    reachable[clause.clause_id] = clause.enforcement
    return reachable


def selected_repository_extensions(
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
) -> set[str]:
    selected_artifacts = {
        planned.artifact.managed_id
        for planned in decision_plan.artifacts
        if planned.present and planned.state == "definite"
    }
    selected: set[str] = set()
    for extension_id, extension in catalog.repository_extensions.items():
        decision = decision_plan.resolved_decisions.get(extension.decision_id)
        if (
            extension.root_pointer_id in selected_artifacts
            and isinstance(decision, dict)
            and decision.get("value") is True
        ):
            selected.add(extension_id)
    return selected


def evaluate_retention(
    catalog: AssetCatalog,
    manifest: dict | None,
    decision_plan: DecisionPlan,
) -> tuple[tuple[RetentionEntry, ...], list[Finding]]:
    if isinstance(manifest, dict) and is_current_strict_manifest(manifest, catalog):
        return (), []
    transition, findings = resolve_upgrade_transition(catalog, manifest)
    if transition is None:
        return (), findings

    reachable_clauses = selected_clause_enforcement(catalog, decision_plan)
    reachable_extensions = selected_repository_extensions(catalog, decision_plan)
    mappings = {mapping.from_clause: mapping for mapping in transition.mappings}
    entries: list[RetentionEntry] = []
    for prior in sorted(transition.prior_clauses, key=lambda item: item.clause_id):
        mapping = mappings.get(prior.clause_id)
        if mapping is None:
            findings.append(
                finding(
                    "retention.clause.unaccounted",
                    "error",
                    str(MANIFEST_PATH),
                    prior.clause_id,
                    f"Prior mandatory clause {prior.clause_id} has no transition mapping.",
                    "Add one reviewed retained, moved, replaced, or rejected mapping with a reason.",
                )
            )
            continue

        entries.append(
            RetentionEntry(
                from_clause=prior.clause_id,
                enforcement=prior.enforcement,
                disposition=mapping.disposition,
                targets=mapping.targets,
                reason=mapping.reason,
            )
        )
        if mapping.disposition == "rejected":
            continue
        if not mapping.targets:
            findings.append(
                finding(
                    "retention.clause.unaccounted",
                    "error",
                    str(MANIFEST_PATH),
                    prior.clause_id,
                    f"Accepted prior clause {prior.clause_id} has no target.",
                    "Add at least one reachable current clause or Repository-Owned Extension target.",
                )
            )
            continue
        for target in mapping.targets:
            if target in catalog.repository_extensions:
                if target not in reachable_extensions:
                    findings.append(
                        finding(
                            "retention.target.unreachable",
                            "error",
                            str(MANIFEST_PATH),
                            prior.clause_id,
                            f"Retention target {target} is not selected in the future artifact graph.",
                            "Select its Repository-Owned Extension or revise the transition mapping.",
                        )
                    )
                continue
            target_enforcement = reachable_clauses.get(target)
            if target_enforcement is None:
                findings.append(
                    finding(
                        "retention.target.unreachable",
                        "error",
                        str(MANIFEST_PATH),
                        prior.clause_id,
                        f"Retention target {target} has no selected supporting-guide carrier.",
                        "Select the target carrier or revise the transition mapping.",
                    )
                )
            elif target_enforcement != prior.enforcement:
                findings.append(
                    finding(
                        "retention.target.enforcement-mismatch",
                        "error",
                        str(MANIFEST_PATH),
                        prior.clause_id,
                        f"Retention target {target} uses {target_enforcement} instead of {prior.enforcement} enforcement.",
                        "Restore equivalent enforcement strength before authorizing the upgrade.",
                    )
                )
    return tuple(entries), findings


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
        if decision_id.startswith(("adoption.", "removal.")) and isinstance(existing_decision, dict):
            resolved[decision_id] = existing_decision
    for decision_id, value in sorted(cli_decisions.items()):
        if decision_id.startswith(("adoption.", "removal.")):
            resolved[decision_id] = {"value": value, "confirmedAt": today}
    preserve_compatible_catalog_decisions(
        catalog=catalog,
        existing=existing,
        cli_decisions=cli_decisions,
        today=today,
        resolved=resolved,
    )

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
        resolved_decisions=resolved,
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


def preserve_compatible_catalog_decisions(
    catalog: AssetCatalog,
    existing: dict,
    cli_decisions: dict[str, object],
    today: str,
    resolved: dict[str, dict],
) -> None:
    for decision_id in sorted(catalog.decisions):
        if decision_id in resolved:
            continue
        decision = compatible_decision_answer(
            contract=catalog.decisions[decision_id],
            decision_id=decision_id,
            existing=existing,
            cli_decisions=cli_decisions,
            today=today,
        )
        if decision is not None:
            resolved[decision_id] = decision


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
    resolved_decisions: dict[str, dict],
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
        for artifact in expected_artifacts_for_profile(
            catalog,
            profile_id,
            artifact_modules,
            template_overrides=template_overrides_for_effects(matched_effects),
            render_values=render_values_for_effects(matched_effects, resolved_decisions),
            dispatch_modules=active_modules,
            strict_tokens=False,
        )
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


def validate_decision_plan_references(
    repo: Path,
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
    existing_manifest: dict | None = None,
) -> list[Finding]:
    expected_artifacts = expected_artifacts_for_plan(decision_plan)
    definite_by_id = {
        artifact.managed_id: artifact for artifact in expected_artifacts
    }
    future_repository_paths = {
        extension.target_path
        for extension in repository_extensions_to_create(
            repo,
            catalog,
            decision_plan,
            existing_manifest,
        )
    }
    findings: list[Finding] = []

    for source in expected_artifacts:
        references = catalog.references_by_artifact.get(source.managed_id, ())
        for reference in references:
            if reference.ownership == "setup":
                target = definite_by_id.get(reference.target_managed_id)
                if target is None:
                    findings.append(
                        finding(
                            "reference.managed.missing",
                            "error",
                            source.path.as_posix(),
                            source.managed_id,
                            (
                                f"Declared reference {reference.reference_id} targets "
                                f"{reference.target_managed_id}, which is not present in "
                                "the definite Decision Plan artifact set."
                            ),
                            (
                                "Update the Decision Plan effects or the setup-owned "
                                "reference so the selected artifacts resolve together."
                            ),
                        )
                    )
                continue

            repository_path = reference.repository_path
            if repository_path is None:
                continue
            findings.extend(
                validate_repository_reference(
                    repo=repo,
                    source=source,
                    reference_id=reference.reference_id,
                    repository_path=repository_path,
                    future_repository_paths=future_repository_paths,
                )
            )

    future_paths = {artifact.path for artifact in expected_artifacts}
    future_absent_paths = {
        planned_artifact.artifact.path
        for planned_artifact in decision_plan.artifacts
        if not planned_artifact.present and planned_artifact.state == "definite"
    } - future_paths
    for artifact in expected_artifacts:
        findings.extend(
            validate_internal_references(
                repo=repo,
                relative_path=artifact.path,
                content=artifact.content,
                managed_id=artifact.managed_id,
                future_paths=future_paths,
                future_absent_paths=future_absent_paths,
            )
        )
    return sorted_findings(findings)


def validate_repository_reference(
    repo: Path,
    source: ExpectedArtifact,
    reference_id: str,
    repository_path: Path,
    future_repository_paths: set[Path] | None = None,
) -> list[Finding]:
    if repository_path.is_absolute() or ".." in repository_path.parts:
        return [
            repository_reference_outside_finding(
                source, reference_id, repository_path
            )
        ]

    repo_root = repo.resolve(strict=False)
    try:
        target = (repo_root / repository_path).resolve(strict=False)
        target.relative_to(repo_root)
    except (OSError, RuntimeError, ValueError):
        return [
            repository_reference_outside_finding(
                source, reference_id, repository_path
            )
        ]

    if target.exists() or repository_path in (future_repository_paths or set()):
        return []
    return [
        finding(
            "reference.repository.missing",
            "error",
            repository_path.as_posix(),
            source.managed_id,
            (
                f"Declared reference {reference_id} requires repository-owned path "
                f"{repository_path.as_posix()}, but it does not exist."
            ),
            (
                "Create the repository-authored target or select a profile that does "
                "not require it; setup will not generate repository-owned content."
            ),
        )
    ]


def active_repository_extensions(
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
) -> list[RepositoryOwnedExtension]:
    definite_artifact_ids = {
        artifact.managed_id for artifact in expected_artifacts_for_plan(decision_plan)
    }
    active = []
    for extension_id in sorted(catalog.repository_extensions):
        extension = catalog.repository_extensions[extension_id]
        decision = decision_plan.resolved_decisions.get(extension.decision_id)
        if not isinstance(decision, dict) or decision.get("value") is not True:
            continue
        if extension.root_pointer_id not in definite_artifact_ids:
            continue
        active.append(extension)
    return active


def manifest_repository_extension_ids(manifest: dict | None) -> set[str]:
    if not isinstance(manifest, dict):
        return set()
    records = manifest.get("repositoryExtensions", [])
    if not isinstance(records, list):
        return set()
    return {
        record["id"]
        for record in records
        if isinstance(record, dict) and isinstance(record.get("id"), str)
    }


def repository_extensions_to_create(
    repo: Path,
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
    existing_manifest: dict | None,
) -> list[RepositoryOwnedExtension]:
    recorded_ids = manifest_repository_extension_ids(existing_manifest)
    creations = []
    for extension in active_repository_extensions(catalog, decision_plan):
        target = repo / extension.target_path
        if target.exists() or extension.extension_id in recorded_ids:
            continue
        creations.append(extension)
    return creations


def repository_extension_content(
    catalog: AssetCatalog,
    extension: RepositoryOwnedExtension,
) -> str:
    templates_root = Path(__file__).resolve().parents[1] / "assets" / "templates"
    return template_content(
        templates_root,
        catalog,
        extension.template_id,
        strict_tokens=True,
    )


def repository_reference_outside_finding(
    source: ExpectedArtifact,
    reference_id: str,
    repository_path: Path,
) -> Finding:
    return finding(
        "reference.repository.outside",
        "error",
        repository_path.as_posix(),
        source.managed_id,
        (
            f"Declared reference {reference_id} resolves outside the repository: "
            f"{repository_path.as_posix()}."
        ),
        "Point the repository-owned reference at an existing path inside the repository.",
    )


def template_overrides_for_effects(effects: Iterable) -> dict[str, str]:
    overrides: dict[str, str] = {}
    for effect in effects:
        for selection in effect.template_selections:
            overrides[selection.artifact_id] = selection.template_id
    return overrides


def render_values_for_effects(
    effects: Iterable,
    resolved_decisions: dict[str, dict],
) -> dict[str, dict[str, str]]:
    values: dict[str, dict[str, str]] = {}
    for effect in effects:
        decision = resolved_decisions.get(effect.decision_id)
        if not isinstance(decision, dict) or "value" not in decision:
            continue
        rendered_value = render_inline_code(str(decision["value"]))
        for binding in effect.render_bindings:
            values.setdefault(binding.artifact_id, {})[binding.token] = rendered_value
    return values


def render_inline_code(value: str) -> str:
    runs = re.findall(r"`+", value)
    delimiter = "`" * (max((len(run) for run in runs), default=0) + 1)
    if value.startswith("`") or value.endswith("`"):
        return f"{delimiter} {value} {delimiter}"
    return f"{delimiter}{value}{delimiter}"


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
    existing_manifest: dict | None = None,
) -> tuple[PlannedChange, ...]:
    manifest_before = current_file_bytes(repo, MANIFEST_PATH)
    changes: list[PlannedChange] = [
        PlannedChange(
            action="refresh manifest" if manifest_before is not None else "create manifest",
            path=MANIFEST_PATH,
            managed_id="manifest",
            state="definite",
            reason="Record the selected profile, decisions, and managed inventory.",
            before_digest=bytes_digest(manifest_before),
        )
    ]
    seen: set[tuple[str, str, str, str, object]] = {
        ("create manifest", MANIFEST_PATH.as_posix(), "manifest", "", "")
    }
    for planned_artifact in decision_plan.artifacts:
        change = planned_change_for_artifact(repo, planned_artifact, existing_manifest)
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
    existing_manifest: dict | None,
) -> PlannedChange | None:
    artifact = planned_artifact.artifact
    target = repo / artifact.path
    before = current_file_bytes(repo, artifact.path)
    try:
        current = before.decode("utf-8") if before is not None else ""
    except UnicodeDecodeError:
        return None
    if not planned_artifact.present:
        if not target.exists() or not manifest_owns_artifact(
            existing_manifest, artifact.managed_id, artifact.path
        ):
            return None
        spans, marker_findings = parse_managed_block_spans(artifact.path, current)
        if marker_findings or artifact.managed_id not in spans:
            return None
        span = spans[artifact.managed_id]
        remaining = current[: span.start] + current[span.end :]
        after = None if artifact.kind == "guide" and not remaining.strip() else remaining.encode("utf-8")
        return PlannedChange(
            action="remove managed content",
            path=artifact.path,
            managed_id=artifact.managed_id,
            state=planned_artifact.state,
            reason="The resolved Decision Plan excludes previously managed content.",
            before_digest=bytes_digest(before),
            after_digest=bytes_digest(after),
            condition=planned_artifact.condition,
        )
    action = "create guide" if artifact.kind == "guide" else "create managed block"
    if target.exists():
        action = "refresh managed content"
    after = render_expected_path(current, [artifact]).encode("utf-8")
    return PlannedChange(
        action=action,
        path=artifact.path,
        managed_id=artifact.managed_id,
        state=planned_artifact.state,
        reason="The resolved Decision Plan requires this managed artifact.",
        before_digest=bytes_digest(before),
        after_digest=bytes_digest(after),
        condition=planned_artifact.condition,
    )


def manifest_owns_artifact(
    manifest: dict | None,
    managed_id: str,
    path: Path,
) -> bool:
    if not isinstance(manifest, dict):
        return False
    return any(
        isinstance(item, dict)
        and item.get("id") == managed_id
        and item.get("path") == path.as_posix()
        for item in manifest.get("managedArtifacts", [])
    )


def plan_apply(
    repo: Path,
    catalog: AssetCatalog,
    profile_override: str | None,
    decision_args: list[str],
    decision_document: StructuredDecisionDocument | None = None,
    readoption: ReadoptionPlanContext | None = None,
) -> tuple[AuditResult, bool, ChangePlan]:
    findings: list[Finding] = []
    existing_manifest, invalid_input = load_manifest_for_apply(repo, findings)
    if invalid_input:
        return empty_apply_result(findings, invalid_input)

    planning_manifest = None if readoption is not None else existing_manifest
    if (
        readoption is None
        and existing_manifest is not None
        and existing_manifest.get("schemaVersion")
        not in {1, MANIFEST_SCHEMA_VERSION_0_0_1}
    ):
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
        planning_manifest.get("profile") if isinstance(planning_manifest, dict) else None
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
    file_decisions, file_findings = structured_decision_values(
        decision_document, catalog
    )
    parse_findings.extend(file_findings)
    for decision_id in sorted(set(cli_decisions).intersection(file_decisions)):
        if cli_decisions[decision_id] == file_decisions[decision_id]:
            continue
        parse_findings.append(
            finding(
                "decision-file.decision.conflict",
                "error",
                decision_document.source_paths[0].as_posix(),
                decision_id,
                "--decision conflicts with the structured decision-file value.",
                "Keep the decision in one input or make both values identical.",
            )
        )
    cli_decisions = {**file_decisions, **cli_decisions}
    findings.extend(parse_findings)
    decision_plan = resolve_decision_plan(
        catalog,
        profile_id,
        planning_manifest,
        cli_decisions,
    )
    preview_changes = planned_changes_for_plan(repo, decision_plan, planning_manifest)
    if parse_findings:
        return empty_apply_result(
            findings,
            True,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
        )

    findings.extend(decision_required_findings(decision_plan))
    if not decision_plan.unresolved_decisions:
        findings.extend(
            validate_decision_plan_references(
                repo,
                catalog,
                decision_plan,
                existing_manifest=planning_manifest,
            )
        )
    retention: tuple[RetentionEntry, ...] = ()
    if not findings and readoption is None:
        retention, retention_findings = evaluate_retention(
            catalog,
            planning_manifest,
            decision_plan,
        )
        findings.extend(retention_findings)
    if findings:
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
            retention=retention,
        )

    decisions = decision_plan.resolved_decisions
    ordered_modules = list(decision_plan.active_modules)
    expected_artifacts = expected_artifacts_for_plan(decision_plan)

    expected_by_id = {artifact.managed_id: artifact for artifact in expected_artifacts}
    current_files = load_current_files(repo, expected_artifacts, planning_manifest, findings)
    if findings:
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
            retention=retention,
        )

    ownership_findings = require_obsolete_artifact_decisions(
        repo=repo,
        existing_manifest=planning_manifest,
        expected_by_id=expected_by_id,
        current_files=current_files,
        decisions=decisions,
    )
    findings.extend(ownership_findings)
    if ownership_findings:
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
            retention=retention,
        )

    adoption_findings = (
        []
        if readoption is not None
        else require_adoption_decisions(
            current_files=current_files,
            expected_artifacts=expected_artifacts,
            decisions=decisions,
        )
    )
    findings.extend(adoption_findings)
    if adoption_findings:
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
            retention=retention,
        )

    changed_contents: dict[Path, str | bytes | None] = {}
    if readoption is not None:
        for relative_path, content in readoption_carrier_outputs(repo, readoption).items():
            if relative_path == MANIFEST_PATH:
                continue
            if relative_path in current_files:
                try:
                    current_files[relative_path] = content.decode("utf-8")
                except UnicodeDecodeError:
                    findings.append(
                        finding(
                            "readoption.carrier.encoding.invalid",
                            "error",
                            relative_path.as_posix(),
                            "source-baseline",
                            "A managed output carrier is not UTF-8 after applying dispositions.",
                            "Move binary evidence to a non-managed carrier before applying.",
                        )
                    )
            elif content != current_file_bytes(repo, relative_path):
                changed_contents[relative_path] = content
        if findings:
            return empty_apply_result(
                findings,
                False,
                selection=decision_plan.selection,
                planned_changes=preview_changes,
            )
    for relative_path, artifacts in artifacts_by_path(expected_artifacts).items():
        current = current_files.get(relative_path, "")
        changed_contents[relative_path] = (
            render_shared_path(current, artifacts)
            if readoption is not None and all(item.kind == "guide" for item in artifacts)
            else render_expected_path(current=current, artifacts=artifacts)
        )

    remove_obsolete_artifacts(
        repo=repo,
        existing_manifest=planning_manifest,
        expected_by_id=expected_by_id,
        current_files=current_files,
        changed_contents=changed_contents,
        decisions=decisions,
    )

    extension_creations = repository_extensions_to_create(
        repo,
        catalog,
        decision_plan,
        planning_manifest,
    )
    for extension in extension_creations:
        changed_contents[extension.target_path] = repository_extension_content(
            catalog,
            extension,
        )

    if (
        readoption is not None
        and readoption.repository_rules
        and current_file_bytes(repo, Path("docs/agents/repository-rules.md")) is None
    ):
        changed_contents[Path("docs/agents/repository-rules.md")] = (
            readoption.repository_rules
        )

    baseline_id = (
        f"baseline.{profile_id}-0.0.1"
        if readoption is not None
        else existing_manifest["generator"]["baseline"]
        if isinstance(existing_manifest, dict)
        and is_current_strict_manifest(existing_manifest, catalog)
        else current_baseline_id(catalog)
    )
    if baseline_id is None:
        findings.append(
            finding(
                "retention.baseline.catalog-invalid",
                "error",
                str(MANIFEST_PATH),
                "manifest",
                "The bundled upgrade catalog has no unique current baseline.",
                "Fix the bundled transition graph before applying managed content.",
            )
        )
        return empty_apply_result(
            findings,
            False,
            selection=decision_plan.selection,
            planned_changes=preview_changes,
            retention=retention,
        )
    repository_extensions = repository_extension_records(
        repo,
        catalog,
        decision_plan,
        planning_manifest,
        extension_creations,
    )
    manifest = (
        build_readoption_manifest(
            profile_id,
            ordered_modules,
            expected_artifacts,
            decisions,
            baseline_id,
            repository_extensions,
            readoption,
            repo,
        )
        if readoption is not None
        else refresh_strict_manifest(
            existing_manifest,
            ordered_modules,
            expected_artifacts,
            decisions,
            repository_extensions,
        )
        if isinstance(existing_manifest, dict)
        and is_current_strict_manifest(existing_manifest, catalog)
        else build_manifest(
            profile_id,
            ordered_modules,
            expected_artifacts,
            decisions,
            planning_manifest,
            baseline_id,
            repository_extensions,
        )
    )
    changed_contents[MANIFEST_PATH] = json.dumps(manifest, indent=2, sort_keys=False) + "\n"
    plan = concrete_change_plan(
        repo=repo,
        catalog=catalog,
        decision_plan=decision_plan,
        existing_manifest=planning_manifest,
        expected_artifacts=expected_artifacts,
        changed_contents=changed_contents,
        current_files=current_files,
        manifest=manifest,
        retention=retention,
        extension_creations=extension_creations,
        decision_document=(
            decision_document.to_json() if decision_document is not None else None
        ),
        decision_file_digests=(
            decision_document.source_digests
            if decision_document is not None
            else ()
        ),
        readoption=readoption,
    )
    validation_findings = validate_change_plan(repo, plan, expected_artifacts)
    findings.extend(validation_findings)
    findings.extend(
        delegation_findings(
            repo,
            catalog,
            profile_id,
            active_modules=decision_plan.active_modules,
        )
    )
    return (
        AuditResult(
            sorted_findings(findings),
            selection=decision_plan.selection,
            planned_changes=plan_operations(plan),
            plan_digest=plan.digest,
            retention=retention,
            decision_document=(
                decision_document.to_json()
                if decision_document is not None
                else None
            ),
            planned_outputs=plan.mutations if readoption is not None else (),
            source_baseline=(
                readoption.source_baseline if readoption is not None else None
            ),
            capabilities=(readoption.capabilities if readoption is not None else ()),
            setup_snapshot=(readoption.setup_snapshot if readoption is not None else None),
            verification=(readoption.verification if readoption is not None else ()),
        ),
        False,
        plan,
    )


def empty_apply_result(
    findings: list[Finding],
    invalid_input: bool,
    selection: PreviewSelection | None = None,
    planned_changes: tuple[PlannedChange, ...] = (),
    retention: tuple[RetentionEntry, ...] = (),
) -> tuple[AuditResult, bool, ChangePlan]:
    return (
        AuditResult(
            sorted_findings(findings),
            selection=selection,
            planned_changes=planned_changes,
            retention=retention,
        ),
        invalid_input,
        ChangePlan("setup", (), None, {}, retention),
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
        value = parse_decision_value(decision_id, raw_value, catalog)
        contract = catalog.decisions.get(decision_id)
        if decision_id.startswith("removal.") and value not in {"preserve", "remove"}:
            findings.append(invalid_decision_finding(decision_id))
            continue
        if (
            contract is not None
            and contract.get("type") == "string"
            and isinstance(value, str)
            and is_unsafe_inline_decision_value(value)
        ):
            findings.append(unsafe_decision_value_finding(decision_id))
            continue
        if contract is not None and not is_decision_value_valid(contract, value):
            findings.append(invalid_decision_finding(decision_id))
            continue
        decisions[decision_id] = value
    return decisions, findings


def structured_decision_values(
    document: StructuredDecisionDocument | None,
    catalog: AssetCatalog,
) -> tuple[dict[str, object], list[Finding]]:
    if document is None:
        return {}, []
    decisions: dict[str, object] = {}
    findings: list[Finding] = []
    source_path = document.source_paths[0].as_posix()
    for decision_id, value in document.decisions:
        if decision_id.startswith("adoption."):
            if not isinstance(value, bool):
                findings.append(invalid_decision_finding(decision_id))
                continue
            decisions[decision_id] = value
            continue
        if decision_id.startswith("removal."):
            if value not in {"preserve", "remove"}:
                findings.append(invalid_decision_finding(decision_id))
                continue
            decisions[decision_id] = value
            continue
        contract = catalog.decisions.get(decision_id)
        if contract is None:
            findings.append(
                finding(
                    "decision-file.decision.unknown",
                    "error",
                    source_path,
                    decision_id,
                    f"Structured decision {decision_id!r} is not in the current catalog.",
                    "Remove the stale decision or select a catalog decision id.",
                )
            )
            continue
        if (
            contract.get("type") == "string"
            and isinstance(value, str)
            and is_unsafe_inline_decision_value(value)
        ):
            findings.append(unsafe_decision_value_finding(decision_id))
            continue
        if not is_decision_value_valid(contract, value):
            findings.append(invalid_decision_finding(decision_id))
            continue
        decisions[decision_id] = value
    return decisions, findings


def parse_decision_value(decision_id: str, raw_value: str, catalog: AssetCatalog) -> object:
    value = raw_value.strip()
    if decision_id.startswith("adoption."):
        return parse_bool_value(value)
    if decision_id.startswith("removal."):
        return value
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


def is_unsafe_inline_decision_value(value: str) -> bool:
    return (
        any(ord(character) < 32 or ord(character) == 127 for character in value)
        or MARKER.search(value) is not None
    )


def unsafe_decision_value_finding(decision_id: str) -> Finding:
    return finding(
        "decision.value.unsafe",
        "error",
        str(MANIFEST_PATH),
        decision_id,
        f"Decision {decision_id} contains unsafe inline Markdown.",
        "Use a single-line value without control characters or setup ownership markers.",
    )


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
        return (
            isinstance(value, str)
            and bool(value.strip())
            and not is_unsafe_inline_decision_value(value)
        )
    if decision_type == "enum":
        return value in decision_contract.get("values", [])
    if decision_type == "http-contract":
        return is_http_contract_decision_value_valid(
            value, decision_contract.get("modes", ())
        )
    return False


def invalid_decision_finding(decision_id: str) -> Finding:
    return finding(
        "decision.value.invalid",
        "error",
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
    expected_paths = {artifact.path for artifact in expected_artifacts}
    paths = set(expected_paths)
    paths.update(manifest_artifact_paths(existing_manifest, repo, findings))
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
        if relative_path in expected_paths:
            findings.extend(marker_findings)
        current[relative_path] = content
    return current


def manifest_artifact_paths(
    existing_manifest: dict | None,
    repo: Path,
    findings: list[Finding],
) -> set[Path]:
    if not isinstance(existing_manifest, dict):
        return set()
    paths: set[Path] = set()
    for artifact in existing_manifest.get("managedArtifacts", []):
        if not isinstance(artifact, dict):
            continue
        managed_id = artifact.get("id")
        if not isinstance(managed_id, str) or not managed_id:
            managed_id = "managedArtifacts"
        path_value = artifact.get("path")
        if not isinstance(path_value, str):
            continue
        relative_path = validate_manifest_artifact_path(path_value, repo, managed_id, findings)
        if relative_path is not None:
            paths.add(relative_path)
    return paths


def validate_manifest_artifact_path(
    path_value: str,
    repo: Path,
    managed_id: str,
    findings: list[Finding],
) -> Path | None:
    relative_path = safe_relative_path(path_value)
    if relative_path is None or not path_is_inside_repo(repo, relative_path):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                managed_id,
                "Managed artifact path must be a safe repository-relative path.",
                "Refresh the managed artifact inventory before applying setup changes.",
            )
        )
        return None
    return relative_path


def safe_relative_path(path_value: str) -> Path | None:
    if not path_value or "\\" in path_value:
        return None
    candidate = Path(path_value)
    if candidate.is_absolute() or ".." in candidate.parts:
        return None
    return candidate


def path_is_inside_repo(repo: Path, relative_path: Path) -> bool:
    repo_root = repo.resolve(strict=False)
    target = (repo_root / relative_path).resolve(strict=False)
    try:
        target.relative_to(repo_root)
    except ValueError:
        return False
    return True


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
        if marker_findings or blocks:
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


def require_obsolete_artifact_decisions(
    repo: Path,
    existing_manifest: dict | None,
    expected_by_id: dict[str, ExpectedArtifact],
    current_files: dict[Path, str],
    decisions: dict[str, dict],
) -> list[Finding]:
    findings: list[Finding] = []
    if not isinstance(existing_manifest, dict):
        return findings
    for artifact in existing_manifest.get("managedArtifacts", []):
        if not isinstance(artifact, dict):
            continue
        managed_id = artifact.get("id")
        path_value = artifact.get("path")
        if not isinstance(managed_id, str) or not isinstance(path_value, str):
            continue
        relative_path = safe_relative_path(path_value)
        if relative_path is None or not path_is_inside_repo(repo, relative_path):
            continue
        expected = expected_by_id.get(managed_id)
        if expected is not None and expected.path == relative_path:
            continue
        content = current_files.get(relative_path)
        if content is None:
            continue
        spans, marker_findings = parse_managed_block_spans(relative_path, content)
        if not marker_findings and managed_id in spans:
            continue
        decision_id = f"removal.{managed_id}"
        decision = decisions.get(decision_id)
        if isinstance(decision, dict) and decision.get("value") in {"preserve", "remove"}:
            continue
        findings.append(
            finding(
                "decision.required",
                "decision",
                relative_path.as_posix(),
                decision_id,
                "Legacy managed artifact is stale but no ownership marker proves setup owns the content.",
                f"Pass --decision {decision_id}=preserve|remove after reviewing the file.",
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
    repo: Path,
    existing_manifest: dict | None,
    expected_by_id: dict[str, ExpectedArtifact],
    current_files: dict[Path, str],
    changed_contents: dict[Path, str | None],
    decisions: dict[str, dict],
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
        relative_path = safe_relative_path(path_value)
        if relative_path is None or not path_is_inside_repo(repo, relative_path):
            continue
        expected = expected_by_id.get(managed_id)
        if expected is not None and expected.path == relative_path:
            continue
        content = changed_contents.get(relative_path, current_files.get(relative_path))
        if content is None:
            continue
        spans, marker_findings = parse_managed_block_spans(relative_path, content)
        if marker_findings:
            removal = decisions.get(f"removal.{managed_id}")
            if isinstance(removal, dict) and removal.get("value") == "remove":
                changed_contents[relative_path] = None
            continue
        span = spans.get(managed_id)
        if span is None:
            removal = decisions.get(f"removal.{managed_id}")
            if isinstance(removal, dict) and removal.get("value") == "remove":
                changed_contents[relative_path] = None
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
    baseline_id: str,
    repository_extensions: list[dict[str, str]],
) -> dict:
    local_skills = []
    if isinstance(existing_manifest, dict) and isinstance(existing_manifest.get("localSkills"), list):
        local_skills = existing_manifest["localSkills"]
    generator: dict[str, object] = {}
    if isinstance(existing_manifest, dict) and isinstance(
        existing_manifest.get("generator"), dict
    ):
        generator = dict(existing_manifest["generator"])
    generator.update(
        {
            "skill": "setup-context-driven",
            "version": 1,
            "baseline": baseline_id,
        }
    )
    ordered_decisions: dict[str, dict] = {}
    existing_decisions = (
        existing_manifest.get("decisions", {}) if isinstance(existing_manifest, dict) else {}
    )
    if isinstance(existing_decisions, dict):
        for decision_id in existing_decisions:
            if decision_id in decisions:
                ordered_decisions[decision_id] = decisions[decision_id]
    for decision_id, decision in decisions.items():
        if decision_id not in ordered_decisions:
            ordered_decisions[decision_id] = decision
    manifest = {
        "schemaVersion": 1,
        "generator": generator,
        "profile": profile_id,
        "modules": ordered_modules,
        "decisions": ordered_decisions,
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
    if repository_extensions:
        manifest["repositoryExtensions"] = repository_extensions
    return manifest


def build_readoption_manifest(
    profile_id: str,
    ordered_modules: list[str],
    expected_artifacts: list[ExpectedArtifact],
    decisions: dict[str, dict],
    baseline_id: str,
    repository_extensions: list[dict[str, str]],
    readoption: ReadoptionPlanContext,
    repo: Path,
) -> dict:
    manifest: dict[str, object] = {
        "schemaVersion": MANIFEST_SCHEMA_VERSION_0_0_1,
        "version": OWNED_VERSION_0_0_1,
        "generator": {
            "skill": "setup-context-driven",
            "version": OWNED_VERSION_0_0_1,
            "baseline": baseline_id,
        },
        "profile": profile_id,
        "modules": ordered_modules,
        "decisions": {
            decision_id: {"value": decision["value"]}
            for decision_id, decision in sorted(decisions.items())
            if isinstance(decision, dict) and "value" in decision
        },
        "sourceBaseline": readoption.source_baseline.identity_json(),
        "sourceEntries": [
            entry.to_json() for entry in readoption.source_baseline.entries
        ],
        "readoptionDispositions": [
            disposition.to_json() for disposition in readoption.dispositions
        ],
        "capabilities": list(readoption.capabilities),
        "setupSnapshot": readoption.setup_snapshot,
        "verification": list(readoption.verification),
        "managedArtifacts": [
            {
                "id": artifact.managed_id,
                "path": artifact.path.as_posix(),
                "kind": artifact.kind,
                "module": artifact.module_id,
                "template": artifact.template_id,
                "version": OWNED_VERSION_0_0_1,
                "digest": artifact.digest,
            }
            for artifact in expected_artifacts
        ],
        "localSkills": [],
    }
    if repository_extensions:
        manifest["repositoryExtensions"] = repository_extensions
    repository_rules_path = Path("docs/agents/repository-rules.md")
    if readoption.repository_rules or (repo / repository_rules_path).is_file():
        manifest["repositoryRules"] = {
            "path": repository_rules_path.as_posix(),
            "ownership": "repository",
        }
    return manifest


def refresh_strict_manifest(
    existing_manifest: dict,
    ordered_modules: list[str],
    expected_artifacts: list[ExpectedArtifact],
    decisions: dict[str, dict],
    repository_extensions: list[dict[str, str]],
) -> dict:
    manifest = dict(existing_manifest)
    manifest["modules"] = ordered_modules
    manifest["decisions"] = {
        decision_id: {"value": decision["value"]}
        for decision_id, decision in sorted(decisions.items())
        if isinstance(decision, dict) and "value" in decision
    }
    manifest["managedArtifacts"] = [
        {
            "id": artifact.managed_id,
            "path": artifact.path.as_posix(),
            "kind": artifact.kind,
            "module": artifact.module_id,
            "template": artifact.template_id,
            "version": OWNED_VERSION_0_0_1,
            "digest": artifact.digest,
        }
        for artifact in expected_artifacts
    ]
    if repository_extensions:
        manifest["repositoryExtensions"] = repository_extensions
    else:
        manifest.pop("repositoryExtensions", None)
    return manifest


def readoption_carrier_outputs(
    repo: Path,
    readoption: ReadoptionPlanContext,
) -> dict[Path, bytes]:
    outputs: dict[Path, bytes] = {}
    entries_by_path: dict[Path, list] = {}
    for entry in readoption.source_baseline.entries:
        entries_by_path.setdefault(entry.path, []).append(entry)
    for relative_path, entries in entries_by_path.items():
        current = current_file_bytes(repo, relative_path)
        if current is None:
            continue
        output = b"".join(
            entry.source_bytes
            for entry in entries
            if entry.kind != "managed-block"
        )
        if output != current:
            outputs[relative_path] = output
    return outputs


def repository_extension_records(
    repo: Path,
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
    existing_manifest: dict | None,
    extension_creations: list[RepositoryOwnedExtension],
) -> list[dict[str, str]]:
    records: dict[str, dict[str, str]] = {}
    if isinstance(existing_manifest, dict):
        existing_records = existing_manifest.get("repositoryExtensions", [])
        if isinstance(existing_records, list):
            for record in existing_records:
                if (
                    isinstance(record, dict)
                    and isinstance(record.get("id"), str)
                    and isinstance(record.get("path"), str)
                ):
                    records[record["id"]] = {
                        "id": record["id"],
                        "path": record["path"],
                    }

    creation_ids = {extension.extension_id for extension in extension_creations}
    for extension in active_repository_extensions(catalog, decision_plan):
        if not (repo / extension.target_path).exists() and extension.extension_id not in creation_ids:
            continue
        records[extension.extension_id] = {
            "id": extension.extension_id,
            "path": extension.target_path.as_posix(),
        }
    return [records[extension_id] for extension_id in sorted(records)]


def concrete_change_plan(
    repo: Path,
    catalog: AssetCatalog,
    decision_plan: DecisionPlan,
    existing_manifest: dict | None,
    expected_artifacts: list[ExpectedArtifact],
    changed_contents: dict[Path, str | bytes | None],
    current_files: dict[Path, str],
    manifest: dict,
    retention: tuple[RetentionEntry, ...],
    extension_creations: list[RepositoryOwnedExtension],
    decision_document: dict[str, object] | None,
    decision_file_digests: tuple[str, ...],
    readoption: ReadoptionPlanContext | None = None,
) -> ChangePlan:
    expected_by_path = artifacts_by_path(expected_artifacts)
    old_by_id = manifest_artifacts_by_id(existing_manifest)
    expected_by_id = {artifact.managed_id: artifact for artifact in expected_artifacts}
    operations_by_path: dict[Path, list[PlannedChange]] = {}

    def add_operation(path: Path, operation: PlannedChange) -> None:
        operations_by_path.setdefault(path, []).append(operation)

    if readoption is not None:
        carrier_outputs = readoption_carrier_outputs(repo, readoption)
        for entry in readoption.source_baseline.entries:
            if entry.kind != "managed-block" or entry.path not in carrier_outputs:
                continue
            add_operation(
                entry.path,
                PlannedChange(
                    action="remove source managed content",
                    path=entry.path,
                    managed_id=entry.entry_id,
                    state="definite",
                    reason="The confirmed Readoption disposition retires prior setup-owned bytes.",
                    before_digest=bytes_digest(current_file_bytes(repo, entry.path)),
                    after_digest=bytes_digest(
                        content_bytes(changed_contents.get(entry.path, carrier_outputs[entry.path]))
                    ),
                ),
            )
        repository_rules_path = Path("docs/agents/repository-rules.md")
        if (
            readoption.repository_rules
            and current_file_bytes(repo, repository_rules_path) is None
        ):
            add_operation(
                repository_rules_path,
                PlannedChange(
                    action="create repository-specific normative rules",
                    path=repository_rules_path,
                    managed_id="repository-rules.readoption",
                    state="definite",
                    reason="Explicit Readoption dispositions authorize these exact repository-owned bytes.",
                    before_digest=None,
                    after_digest=bytes_digest(readoption.repository_rules),
                ),
            )

    for extension in extension_creations:
        path = extension.target_path
        before = current_file_bytes(repo, path)
        after = content_bytes(changed_contents[path])
        add_operation(
            path,
            PlannedChange(
                action="create repository extension",
                path=path,
                managed_id=extension.extension_id,
                state="definite",
                reason=(
                    "The resolved Decision Plan authorizes the one-time "
                    "Repository-Owned Extension scaffold."
                ),
                before_digest=bytes_digest(before),
                after_digest=bytes_digest(after),
            ),
        )

    for path, artifacts in expected_by_path.items():
        before = current_file_bytes(repo, path)
        after = content_bytes(changed_contents.get(path, current_files.get(path)))
        if before == after:
            continue
        before_content = current_files.get(path, "")
        before_spans, _ = parse_managed_block_spans(path, before_content)
        for artifact in artifacts:
            old = old_by_id.get(artifact.managed_id)
            old_path = safe_relative_path(old.get("path", "")) if isinstance(old, dict) else None
            from_path = old_path if old_path is not None and old_path != path else None
            before_span = before_spans.get(artifact.managed_id)
            expected_block = managed_block(
                artifact.managed_id,
                artifact.version,
                artifact.content,
            )
            if (
                from_path is None
                and before_span is not None
                and before_content[before_span.start : before_span.end] == expected_block
            ):
                continue
            if from_path is not None:
                action = "rename managed content"
                reason = "The selected artifact moved to its catalog path."
            elif artifact.managed_id not in before_spans:
                action = "create guide" if artifact.kind == "guide" else "create managed block"
                reason = "The selected profile requires this managed artifact."
            elif catalog.references_by_artifact.get(artifact.managed_id):
                action = "edit managed references"
                reason = "The selected Decision Plan changes this artifact's managed references."
            else:
                action = "refresh managed content"
                reason = "The selected catalog content differs from the repository bytes."
            add_operation(
                path,
                PlannedChange(
                    action=action,
                    path=path,
                    managed_id=artifact.managed_id,
                    state="definite",
                    reason=reason,
                    before_digest=bytes_digest(before),
                    after_digest=bytes_digest(after),
                    from_path=from_path,
                    reference_edits=reference_edit_details(catalog, artifact.managed_id)
                    if action == "edit managed references"
                    else (),
                ),
            )

    for managed_id, old in sorted(old_by_id.items()):
        path = safe_relative_path(old.get("path", ""))
        if path is None:
            continue
        expected = expected_by_id.get(managed_id)
        if expected is not None and expected.path == path:
            continue
        before = current_file_bytes(repo, path)
        after = content_bytes(changed_contents.get(path, current_files.get(path)))
        removal_id = f"removal.{managed_id}"
        removal = decision_plan.resolved_decisions.get(removal_id)
        if before == after and isinstance(removal, dict) and removal.get("value") == "preserve":
            add_operation(
                path,
                PlannedChange(
                    action="preserve unmarked content",
                    path=path,
                    managed_id=managed_id,
                    state="definite",
                    reason="The explicit removal decision preserves repository-authored bytes.",
                    before_digest=bytes_digest(before),
                    after_digest=bytes_digest(after),
                ),
            )
        elif before != after:
            add_operation(
                path,
                PlannedChange(
                    action="remove managed content",
                    path=path,
                    managed_id=managed_id,
                    state="definite",
                    reason="The selected profile no longer includes previously owned content.",
                    before_digest=bytes_digest(before),
                    after_digest=bytes_digest(after),
                ),
            )

    manifest_before = current_file_bytes(repo, MANIFEST_PATH)
    manifest_after = content_bytes(changed_contents[MANIFEST_PATH])
    if manifest_before != manifest_after:
        add_operation(
            MANIFEST_PATH,
            PlannedChange(
                action="refresh manifest" if manifest_before is not None else "create manifest",
                path=MANIFEST_PATH,
                managed_id="manifest",
                state="definite",
                reason="Record the selected profile, decisions, and managed inventory.",
                before_digest=bytes_digest(manifest_before),
                after_digest=bytes_digest(manifest_after),
            ),
        )

    mutations: list[FileMutation] = []
    all_paths = set(operations_by_path)
    all_paths.update(
        path
        for path, content in changed_contents.items()
        if current_file_bytes(repo, path) != content_bytes(content)
    )
    for path in sorted(all_paths, key=lambda item: item.as_posix()):
        before = current_file_bytes(repo, path)
        after = content_bytes(changed_contents.get(path, current_files.get(path)))
        mutations.append(
            FileMutation(
                path=path,
                before_digest=bytes_digest(before),
                after_digest=bytes_digest(after),
                content=after,
                operations=tuple(operations_by_path.get(path, ())),
            )
        )

    digest_payload = {
        "kind": "baseline-readoption" if readoption is not None else "setup",
        "selection": decision_plan.selection.to_json(),
        "decisions": {
            key: value.get("value")
            for key, value in sorted(decision_plan.resolved_decisions.items())
        },
        "catalogDigest": catalog_digest(catalog),
        "decisionDocument": decision_document,
        "decisionFileDigests": list(decision_file_digests),
        "retentionAccounting": [entry.to_json() for entry in retention],
        "operations": [
            operation.to_json()
            for mutation in mutations
            for operation in mutation.operations
        ],
        "paths": [
            {
                "path": mutation.path.as_posix(),
                "beforeDigest": mutation.before_digest,
                "afterDigest": mutation.after_digest,
            }
            for mutation in mutations
        ],
    }
    if readoption is not None:
        digest_payload.update(
            {
                "sourceBaseline": readoption.source_baseline.identity_json(),
                "sourceEntries": [
                    entry.to_json() for entry in readoption.source_baseline.entries
                ],
                "dispositions": [
                    disposition.to_json() for disposition in readoption.dispositions
                ],
                "capabilities": list(readoption.capabilities),
                "setupSnapshot": readoption.setup_snapshot,
                "verification": list(readoption.verification),
                "outputs": [mutation.output_json() for mutation in mutations],
            }
        )
    digest = hashlib.sha256(canonical_json_bytes(digest_payload)).hexdigest()
    return ChangePlan(
        "baseline-readoption" if readoption is not None else "setup",
        tuple(mutations),
        digest,
        manifest,
        retention,
        readoption,
    )


def manifest_artifacts_by_id(manifest: dict | None) -> dict[str, dict]:
    if not isinstance(manifest, dict):
        return {}
    return {
        item["id"]: item
        for item in manifest.get("managedArtifacts", [])
        if isinstance(item, dict) and isinstance(item.get("id"), str)
    }


def reference_edit_details(
    catalog: AssetCatalog,
    managed_id: str,
) -> tuple[dict[str, str], ...]:
    details: list[dict[str, str]] = []
    for reference in catalog.references_by_artifact.get(managed_id, ()):
        detail = {"id": reference.reference_id, "ownership": reference.ownership}
        if reference.target_managed_id is not None:
            detail["targetManagedId"] = reference.target_managed_id
        if reference.repository_path is not None:
            detail["repositoryPath"] = reference.repository_path.as_posix()
        details.append(detail)
    return tuple(details)


def plan_operations(plan: ChangePlan) -> tuple[PlannedChange, ...]:
    return tuple(
        operation
        for mutation in plan.mutations
        for operation in mutation.operations
    )


def bytes_digest(content: bytes | None) -> str | None:
    if content is None:
        return None
    return hashlib.sha256(content).hexdigest()


def canonical_json_bytes(value: object) -> bytes:
    return json.dumps(
        normalize_digest_value(value),
        sort_keys=True,
        separators=(",", ":"),
        ensure_ascii=False,
    ).encode("utf-8")


def normalize_digest_value(value: object) -> object:
    if is_dataclass(value):
        return {
            field.name: normalize_digest_value(getattr(value, field.name))
            for field in fields(value)
        }
    if isinstance(value, Path):
        return value.as_posix()
    if isinstance(value, Mapping):
        return {str(key): normalize_digest_value(item) for key, item in value.items()}
    if isinstance(value, (list, tuple)):
        return [normalize_digest_value(item) for item in value]
    return value


def catalog_digest(catalog: AssetCatalog) -> str:
    return hashlib.sha256(canonical_json_bytes(catalog)).hexdigest()


def validate_change_plan(
    repo: Path,
    plan: ChangePlan,
    expected_artifacts: list[ExpectedArtifact],
) -> list[Finding]:
    findings: list[Finding] = []
    planned = {mutation.path: mutation.content for mutation in plan.mutations}
    for artifact in expected_artifacts:
        content = planned.get(artifact.path)
        if content is None:
            path = repo / artifact.path
            if path.exists() and path.is_file():
                content = path.read_bytes()
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
        try:
            text = content.decode("utf-8")
        except UnicodeDecodeError:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    artifact.path.as_posix(),
                    artifact.managed_id,
                    "Change plan output is not UTF-8 text.",
                    "Rebuild the apply plan from bundled assets.",
                )
            )
            continue
        blocks, marker_findings = parse_managed_blocks(artifact.path, text)
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
                    "error",
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
    temp_by_target: dict[Path, Path] = {}
    mutated_targets: set[Path] = set()
    try:
        for mutation in plan.mutations:
            target = safe_repo_target(repo, mutation.path)
            originals[target] = target.read_bytes() if target.exists() and target.is_file() else None
            if bytes_digest(originals[target]) != mutation.before_digest:
                raise OSError(f"preimage changed for {mutation.path}")
            if mutation.before_digest == mutation.after_digest or mutation.content is None:
                continue
            ensure_parent_dir(target.parent, created_dirs)
            temp_by_target[target] = write_unique_temp_bytes(target, mutation.content, temp_paths)

        for mutation in plan.mutations:
            target = safe_repo_target(repo, mutation.path)
            observed = target.read_bytes() if target.exists() and target.is_file() else None
            if bytes_digest(observed) != mutation.before_digest:
                raise OSError(f"preimage changed for {mutation.path}")

        for mutation in plan.mutations:
            target = safe_repo_target(repo, mutation.path)
            if mutation.before_digest == mutation.after_digest:
                continue
            mutated_targets.add(target)
            if mutation.content is None:
                if target.exists() and target.is_file():
                    target.unlink()
                continue
            temp_path = temp_by_target[target]
            temp_path.replace(target)

        mismatches = [
            mutation.path.as_posix()
            for mutation in plan.mutations
            if bytes_digest(current_file_bytes(repo, mutation.path)) != mutation.after_digest
        ]
        if mismatches:
            raise OSError(f"postwrite delta mismatch for {', '.join(mismatches)}")
    except OSError:
        for target in mutated_targets:
            original = originals[target]
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


def safe_repo_target(repo: Path, relative_path: Path) -> Path:
    if relative_path.is_absolute() or ".." in relative_path.parts:
        raise OSError(f"refusing to write outside repository: {relative_path}")
    repo_root = repo.resolve(strict=False)
    target = (repo_root / relative_path).resolve(strict=False)
    try:
        target.relative_to(repo_root)
    except ValueError as error:
        raise OSError(f"refusing to write outside repository: {relative_path}") from error
    return target


def write_unique_temp_bytes(target: Path, content: bytes, temp_paths: list[Path]) -> Path:
    with tempfile.NamedTemporaryFile(
        "wb",
        dir=target.parent,
        prefix=f".{target.name}.setup-context.",
        suffix=".tmp",
        delete=False,
    ) as handle:
        temp_path = Path(handle.name)
        temp_paths.append(temp_path)
        handle.write(content)
    return temp_path


def write_unique_temp_text(target: Path, content: str, temp_paths: list[Path]) -> Path:
    return write_unique_temp_bytes(target, content.encode("utf-8"), temp_paths)


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


def content_bytes(content: str | bytes | None) -> bytes | None:
    if content is None:
        return None
    if isinstance(content, bytes):
        return content
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
    source_dir = source_dir.resolve(strict=False)
    source_doc, source_path, invalid_input = load_source_setup_doc(setup_id, source_dir, findings)
    if invalid_input or source_doc is None or source_path is None:
        return None, findings, invalid_input
    checkout, checkout_findings = inspect_git_source_checkout(
        source_dir,
        source_path,
        setup_id,
    )
    findings.extend(checkout_findings)
    if checkout is None:
        return None, findings, True
    declared_setup_source = source_doc.get("source")
    if declared_setup_source is not None:
        source_relative = source_path.relative_to(checkout.root).as_posix()
        declared_error = validate_declared_external_source(
            declared_setup_source,
            checkout,
            source_relative,
        )
        if declared_error is not None:
            findings.append(
                setup_snapshot_invalid_finding(source_path, setup_id, declared_error)
            )
            return None, findings, True

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
    raw_skills = source_doc.get("skills")
    if not isinstance(raw_skills, list) or not raw_skills:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup must contain a non-empty skills list.",
            )
        )
        return None, findings, True

    for raw_skill in raw_skills:
        skill, skill_findings = normalize_source_skill(
            raw_skill=raw_skill,
            source_dir=source_dir,
            current_by_path=current_by_path,
            current_by_name=current_by_name,
            source_path=source_path,
            checkout=checkout,
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

    if not skills:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup must contain a non-empty skills list.",
            )
        )
        invalid_input = True

    source_relative = source_path.relative_to(checkout.root).as_posix()
    source_metadata = {
        "type": "github",
        "repository": checkout.repository,
        "ref": checkout.revision,
        "path": source_relative,
    }

    activation_bundles = current_snapshot.get("activationBundles")
    snapshot = {
        "schemaVersion": "setup-context-driven/setup-snapshot/0.0.1",
        "id": setup_id,
        "version": "0.0.1",
        "source": source_metadata,
        "digest": setup_snapshot_digest(skills, activation_bundles),
        "skills": skills,
    }
    if isinstance(activation_bundles, list):
        snapshot["activationBundles"] = activation_bundles
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
    checkout: GitSourceCheckout,
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
    current_source = current_skill.get("source")
    repo_owned = str(name).strip() in REPO_OWNED_SKILLS or (
        isinstance(current_source, dict) and current_source.get("type") == "repo"
    )
    if repo_owned:
        digest = current_skill.get("contentDigest")
        if not isinstance(digest, str) or not re.fullmatch(r"[0-9a-f]{64}", digest):
            digest = raw_data.get("contentDigest")
        if not isinstance(digest, str) or not re.fullmatch(r"[0-9a-f]{64}", digest):
            findings.append(
                setup_snapshot_invalid_finding(
                    source_path,
                    str(name),
                    f"Repository-owned skill {name} is missing a valid content digest.",
                )
            )
            return None, findings
        return {
            "name": str(name).strip(),
            "path": normalized_path,
            "source": {"type": "repo", "name": "roundfix"},
            "contentDigest": digest,
        }, findings

    declared_source = raw_data.get("source")
    effective_revision = checkout.revision
    if declared_source is not None:
        declared_error = validate_declared_external_source(
            declared_source,
            checkout,
            normalized_path,
        )
        if declared_error is not None:
            findings.append(
                setup_snapshot_invalid_finding(source_path, str(name), declared_error)
            )
            return None, findings
        effective_revision = declared_source["ref"]

    source_skill_path = resolve_canonical_skill_tree(source_dir, raw_path)
    if source_skill_path is None or not source_skill_path.exists():
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} source directory is missing.",
            )
        )
        return None, findings
    try:
        source_relative = source_skill_path.relative_to(checkout.root)
    except ValueError:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} source directory escapes the Git checkout.",
            )
        )
        return None, findings
    if source_relative.as_posix() != normalized_path:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} source path does not match {normalized_path}.",
            )
        )
        return None, findings

    try:
        tree_digest = portable_tree_digest(source_skill_path)
        committed_digest = committed_skill_tree_digest(
            checkout.root.as_posix(),
            effective_revision,
            normalized_path,
        )
    except (PortableTreeError, OSError, ValueError) as error:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} cannot be proven: {error}.",
            )
        )
        return None, findings
    if tree_digest != committed_digest:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} bytes do not match commit {effective_revision}.",
            )
        )
        return None, findings
    declared_digest = raw_data.get("treeDigest")
    if declared_digest is not None and declared_digest != tree_digest:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                str(name),
                f"Canonical setup skill {name} tree digest does not match its source bytes.",
            )
        )
        return None, findings

    return {
        "name": str(name).strip(),
        "path": normalized_path,
        "source": {
            "type": "github",
            "repository": checkout.repository,
            "ref": effective_revision,
            "path": normalized_path,
        },
        "treeDigest": tree_digest,
    }, findings


def normalize_skill_path(raw_path: str) -> str | None:
    value = raw_path.strip().replace("\\", "/")
    if not value:
        return None
    candidate = Path(value)
    if (
        candidate.is_absolute()
        or ".." in candidate.parts
        or candidate.as_posix() == "."
        or ("/" in value and candidate.as_posix() != value)
    ):
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


def resolve_canonical_skill_tree(source_dir: Path, raw_path: str) -> Path | None:
    source_path = resolve_canonical_skill_file(source_dir, raw_path)
    if source_path is None:
        return None
    if source_path.name == "SKILL.md":
        return source_path.parent
    return source_path


def inspect_git_source_checkout(
    source_dir: Path,
    source_path: Path,
    setup_id: str,
) -> tuple[GitSourceCheckout | None, list[Finding]]:
    findings: list[Finding] = []
    try:
        root_output = run_git(source_dir, "rev-parse", "--show-toplevel")
        root = Path(root_output.decode("utf-8").strip()).resolve(strict=False)
        revision = run_git(root, "rev-parse", "--verify", "HEAD^{commit}").decode(
            "ascii"
        ).strip()
        status = run_git(root, "status", "--porcelain=v1", "--untracked-files=all")
        remote = run_git(root, "remote", "get-url", "origin").decode("utf-8").strip()
    except (OSError, UnicodeDecodeError, subprocess.CalledProcessError) as error:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                f"Canonical setup source is not a readable Git checkout: {git_error_message(error)}.",
            )
        )
        return None, findings

    if not IMMUTABLE_GIT_REF.fullmatch(revision):
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup source did not resolve to a full immutable commit.",
            )
        )
    if status:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup Git checkout contains dirty or untracked source bytes.",
            )
        )
    repository = github_repository_from_remote(remote)
    if repository is None:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup Git origin is not a portable GitHub repository identity.",
            )
        )
    try:
        relative_source = source_path.relative_to(root).as_posix()
    except ValueError:
        relative_source = ""
    if not relative_source or safe_relative_path(relative_source) is None:
        findings.append(
            setup_snapshot_invalid_finding(
                source_path,
                setup_id,
                "Canonical setup source file is outside the Git checkout.",
            )
        )
    else:
        try:
            committed_source = run_git(root, "show", f"{revision}:{relative_source}")
            if committed_source != source_path.read_bytes():
                findings.append(
                    setup_snapshot_invalid_finding(
                        source_path,
                        setup_id,
                        "Canonical setup source file bytes do not match the declared commit.",
                    )
                )
        except (OSError, subprocess.CalledProcessError) as error:
            findings.append(
                setup_snapshot_invalid_finding(
                    source_path,
                    setup_id,
                    f"Canonical setup source file is not committed: {git_error_message(error)}.",
                )
            )
    if findings or repository is None:
        return None, findings
    return GitSourceCheckout(root=root, repository=repository, revision=revision), findings


def run_git(directory: Path, *args: str) -> bytes:
    environment = os.environ.copy()
    environment["GIT_TERMINAL_PROMPT"] = "0"
    result = subprocess.run(
        ["git", "-C", str(directory), *args],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        check=False,
        env=environment,
    )
    if result.returncode != 0:
        raise subprocess.CalledProcessError(
            result.returncode,
            result.args,
            output=result.stdout,
            stderr=result.stderr,
        )
    return result.stdout


def git_error_message(error: BaseException) -> str:
    if isinstance(error, subprocess.CalledProcessError):
        stderr = error.stderr.decode("utf-8", errors="replace").strip()
        if stderr:
            return stderr
    return str(error)


def github_repository_from_remote(remote: str) -> str | None:
    match = GITHUB_REMOTE.fullmatch(remote)
    if match is None:
        return None
    return match.group(1)


def validate_declared_external_source(
    raw_source: object,
    checkout: GitSourceCheckout,
    source_path: str,
) -> str | None:
    if not isinstance(raw_source, dict) or not raw_source:
        return "Canonical setup skill has empty external provenance."
    if set(raw_source) != {"type", "repository", "ref", "path"}:
        return "Canonical setup skill provenance fields are incomplete or machine-local."
    if raw_source.get("type") != "github":
        return "Canonical setup skill provider must be github."
    if raw_source.get("repository") != checkout.repository:
        return "Canonical setup skill repository does not match the Git origin."
    revision = raw_source.get("ref")
    if not isinstance(revision, str) or not IMMUTABLE_GIT_REF.fullmatch(revision):
        return "Canonical setup skill ref must be a full immutable commit."
    try:
        resolved = run_git(
            checkout.root,
            "rev-parse",
            "--verify",
            f"{revision}^{{commit}}",
        ).decode("ascii").strip()
    except (OSError, UnicodeDecodeError, subprocess.CalledProcessError):
        return "Canonical setup skill ref is not available in the Git checkout."
    if resolved != revision:
        return "Canonical setup skill ref did not resolve to the declared full commit."
    declared_path = raw_source.get("path")
    if not isinstance(declared_path, str) or safe_relative_path(declared_path) is None:
        return "Canonical setup skill source path is not portable."
    if declared_path != source_path:
        return "Canonical setup skill source path does not match its setup path."
    return None


@lru_cache(maxsize=None)
def committed_skill_tree_digest(
    checkout_root: str,
    revision: str,
    source_path: str,
) -> str:
    root = Path(checkout_root)
    output = run_git(
        root,
        "ls-tree",
        "-r",
        "-z",
        "--full-tree",
        revision,
        "--",
        source_path,
    )
    records: list[tuple[bytes, bytes]] = []
    prefix = source_path.encode("utf-8") + b"/"
    for raw_entry in output.split(b"\0"):
        if not raw_entry:
            continue
        header, path_bytes = raw_entry.split(b"\t", 1)
        mode, object_type, object_id = header.split(b" ", 2)
        if not path_bytes.startswith(prefix):
            raise ValueError(f"Git tree entry escapes source path {source_path}")
        relative = path_bytes[len(prefix) :]
        if any(
            part in {b".git", b"node_modules"}
            for part in relative.split(b"/")[:-1]
        ):
            continue
        if object_type != b"blob" or mode not in {b"100644", b"100755"}:
            raise ValueError(
                f"Git tree entry is not a regular file: {os.fsdecode(relative)}"
            )
        content = run_git(root, "cat-file", "blob", object_id.decode("ascii"))
        records.append((relative, content))
    if not records:
        raise ValueError(f"Git commit has no regular files under {source_path}")
    return portable_file_digest(records)


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
    temp_by_target: dict[Path, Path] = {}
    try:
        for target, content in sorted(changes.items(), key=lambda item: item[0].as_posix()):
            originals[target] = target.read_bytes() if target.exists() and target.is_file() else None
            temp_by_target[target] = write_unique_temp_text(target, content, temp_paths)
        for target in sorted(changes, key=lambda item: item.as_posix()):
            temp_path = temp_by_target[target]
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
    strict_manifest = is_current_strict_manifest(manifest, catalog)
    if manifest.get("schemaVersion") != 1 and not strict_manifest:
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
    if strict_manifest:
        snapshot = manifest.get("setupSnapshot")
        if (
            not isinstance(snapshot, dict)
            or snapshot.get("schemaVersion")
            != "setup-context-driven/profile-snapshot/0.0.1"
            or snapshot.get("version") != OWNED_VERSION_0_0_1
        ):
            findings.append(
                finding(
                    "manifest.snapshot.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "setupSnapshot",
                    "Strict Setup Manifest has no valid 0.0.1 Setup Snapshot.",
                    "Reapply the confirmed Change Plan to restore the strict snapshot.",
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

    validate_manifest_repository_extensions(manifest, catalog, findings)


def validate_manifest_repository_extensions(
    manifest: dict,
    catalog: AssetCatalog,
    findings: list[Finding],
) -> None:
    records = manifest.get("repositoryExtensions", [])
    if not isinstance(records, list):
        findings.append(
            finding(
                "manifest.invalid",
                "error",
                str(MANIFEST_PATH),
                "repositoryExtensions",
                "Manifest repositoryExtensions must be a list.",
                "Refresh the Repository-Owned Extension records.",
            )
        )
        return

    seen: set[str] = set()
    for record in records:
        if not isinstance(record, dict) or set(record) != {"id", "path"}:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    "repositoryExtensions",
                    "Each Repository-Owned Extension record must contain only id and path.",
                    "Refresh the Repository-Owned Extension records.",
                )
            )
            continue
        extension_id = record.get("id")
        extension = catalog.repository_extensions.get(extension_id)
        if extension is None or record.get("path") != extension.target_path.as_posix():
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    str(extension_id),
                    "Manifest Repository-Owned Extension record does not match the catalog.",
                    "Refresh the Repository-Owned Extension records.",
                )
            )
            continue
        if extension_id in seen:
            findings.append(
                finding(
                    "manifest.invalid",
                    "error",
                    str(MANIFEST_PATH),
                    extension_id,
                    "Repository-Owned Extension appears more than once in the manifest.",
                    "Keep one record per Repository-Owned Extension.",
                )
            )
        seen.add(extension_id)


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
    elif decision_type == "http-contract":
        valid = is_http_contract_decision_value_valid(
            value, decision_contract.get("modes", ())
        )

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
    repo: Path,
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

        path_value = artifact.get("path")
        if not isinstance(path_value, str) or validate_manifest_artifact_path(
            path_value,
            repo,
            managed_id,
            findings,
        ) is None:
            continue

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
                        "error",
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
    open_marker: tuple[str, int | str, int, int] | None = None

    for marker in MARKER.finditer(content):
        text = marker.group(0)
        begin = BEGIN_MARKER.fullmatch(text)
        end = END_MARKER.fullmatch(text)
        if begin:
            managed_id = begin.group(1)
            version_text = begin.group(2)
            version: int | str = (
                int(version_text) if version_text.isdigit() else version_text
            )
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
    future_paths: set[Path] | None = None,
    future_absent_paths: set[Path] | None = None,
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
            candidate_relative = candidate.relative_to(repo.resolve(strict=False))
        except ValueError:
            findings.append(
                broken_reference_finding(relative_path, managed_id, target)
            )
            continue
        if future_paths is not None and candidate_relative in future_paths:
            continue
        if (
            future_absent_paths is not None
            and candidate_relative in future_absent_paths
        ):
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
    template_overrides: dict[str, str] | None = None,
    render_values: dict[str, dict[str, str]] | None = None,
    dispatch_modules: Iterable[str] | None = None,
    strict_tokens: bool = True,
) -> list[ExpectedArtifact]:
    artifacts: list[ExpectedArtifact] = []
    templates_root = Path(__file__).resolve().parents[1] / "assets" / "templates"
    modules = ordered_modules or catalog.ordered_modules_by_profile[profile_id]
    dispatch_modules = list(dispatch_modules) if dispatch_modules is not None else modules
    template_overrides = template_overrides or {}
    render_values = render_values or {}
    marker_version: int | str | None = None
    strict_profile = catalog.standard_profiles.get(profile_id)
    if strict_profile is not None:
        marker_version = strict_profile.marker_version
    for module_id in modules:
        module = catalog.modules[module_id]
        for block in module.get("rootBlocks", []):
            template_id = template_overrides.get(block["id"], block["template"])
            values = computed_render_values(
                catalog,
                dispatch_modules,
                block,
            )
            values.update(render_values.get(block["id"], {}))
            content = template_content(
                templates_root,
                catalog,
                template_id,
                values,
                strict_tokens,
            )
            artifacts.append(
                ExpectedArtifact(
                    managed_id=block["id"],
                    path=ROOT_INSTRUCTIONS_PATH,
                    kind="root-block",
                    module_id=module_id,
                    template_id=template_id,
                    version=marker_version or block["version"],
                    content=content,
                    digest=managed_digest(content),
                )
            )
        for guide in module.get("supportingGuides", []):
            template_id = template_overrides.get(guide["id"], guide["template"])
            values = computed_render_values(
                catalog,
                dispatch_modules,
                guide,
            )
            values.update(render_values.get(guide["id"], {}))
            content = template_content(
                templates_root,
                catalog,
                template_id,
                values,
                strict_tokens,
            )
            artifacts.append(
                ExpectedArtifact(
                    managed_id=guide["id"],
                    path=Path(guide["path"]),
                    kind="guide",
                    module_id=module_id,
                    template_id=template_id,
                    version=marker_version or guide["version"],
                    content=content,
                    digest=managed_digest(content),
                )
            )
    return artifacts


def computed_render_values(
    catalog: AssetCatalog,
    active_modules: Iterable[str],
    artifact: dict,
) -> dict[str, str]:
    values: dict[str, str] = {}
    rule_lines: list[str] = []
    for rule_id in artifact.get("rules", []):
        rule = catalog.rule_contracts.get(rule_id)
        if rule is None:
            continue
        if rule.clauses:
            rule_lines.extend(
                f"- **{clause.enforcement}**: {clause.guidance.strip()}"
                for clause in rule.clauses
            )
        elif rule.guidance.strip():
            rule_lines.append(f"- {rule.guidance.strip()}")
    if rule_lines:
        values["artifact.rules"] = "\n\n".join(rule_lines)

    if artifact.get("id") == "guide.skill-dispatch":
        values["active-modules.skill-dispatch"] = render_skill_dispatch(
            catalog,
            active_modules,
        )

    artifact_paths = managed_artifact_paths(catalog)
    for reference in catalog.references_by_artifact.get(artifact.get("id"), ()):
        if reference.ownership == "setup":
            path = artifact_paths.get(reference.target_managed_id)
        else:
            path = reference.repository_path
        if path is not None:
            values[reference.token] = render_inline_code(path.as_posix())
    return values


def render_skill_dispatch(
    catalog: AssetCatalog,
    active_modules: Iterable[str],
) -> str:
    active_module_ids = set(active_modules)
    lines: list[str] = []
    active_activations = [
        activation
        for activation in catalog.skill_activations
        if activation.owner_module in active_module_ids
    ]
    if active_activations:
        lines.extend(["Exact activation bundles:", ""])
        for activation in active_activations:
            description = activation.when
            if activation.capability_condition is not None:
                description = (
                    f"When `{activation.capability_condition}` is selected: "
                    f"{description}"
                )
            bundle = catalog.activation_bundles[activation.bundle_id]
            rendered_skills = ", ".join(f"`{skill}`" for skill in bundle.skills)
            lines.append(f"- `{activation.trigger_id}`: {description}")
            lines.append(f"  - `{bundle.bundle_id}`: {rendered_skills}")
        lines.extend(["", "Individual skill triggers:", ""])
    for skill_name, triggers in catalog.skill_dispatch_by_skill.items():
        active_triggers = [
            trigger for trigger in triggers if trigger.owner_module in active_module_ids
        ]
        if not active_triggers:
            continue
        lines.append(f"- `{skill_name}`:")
        lines.extend(
            f"  - `{trigger.trigger_id}`: {trigger.when}"
            for trigger in active_triggers
        )
    return "\n".join(lines)


def managed_artifact_paths(catalog: AssetCatalog) -> dict[str, Path]:
    paths: dict[str, Path] = {}
    for module in catalog.modules.values():
        for block in module.get("rootBlocks", []):
            paths[block["id"]] = ROOT_INSTRUCTIONS_PATH
        for guide in module.get("supportingGuides", []):
            paths[guide["id"]] = Path(guide["path"])
    return paths


def template_content(
    templates_root: Path,
    catalog: AssetCatalog,
    template_id: str,
    render_values: dict[str, str] | None = None,
    strict_tokens: bool = True,
) -> str:
    template = catalog.templates[template_id]
    content = (templates_root / template["path"]).read_text(encoding="utf-8")
    render_values = render_values or {}

    def replace_token(match: re.Match) -> str:
        token = match.group(1)
        rendered = render_values.get(token)
        if rendered is not None:
            return rendered
        if strict_tokens:
            raise ValueError(f"missing render value for {template_id}:{token}")
        return match.group(0)

    return TEMPLATE_TOKEN.sub(replace_token, content)


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
    if not result.findings and not result.planned_changes:
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
    if result.source_baseline is not None:
        source = result.source_baseline
        lines.append("source baseline:")
        lines.append(
            f"- {source.baseline_id} declared={source.declared_identity} "
            f"carriers={len(source.carriers)} entries={len(source.entries)} "
            f"bytes={source.byte_count} digest={source.digest}"
        )
        lines.append("source entries:")
        for entry in source.entries:
            provenance = json.dumps(
                dict(entry.provenance), sort_keys=True, separators=(",", ":")
            )
            lines.append(
                f"- {entry.entry_id} {entry.path.as_posix()} kind={entry.kind} "
                f"range={entry.start}:{entry.end} digest={entry.digest} "
                f"provenance={provenance}"
            )
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
    if result.retention:
        lines.append("retention accounting:")
        for entry in result.retention:
            targets = ", ".join(entry.targets) if entry.targets else "-"
            lines.append(
                f"- {entry.from_clause} enforcement={entry.enforcement} "
                f"disposition={entry.disposition} targets={targets}"
            )
            lines.append(f"  reason: {entry.reason}")
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
    remediation: dict[str, object] | None = None,
) -> Finding:
    return Finding(
        code=code,
        severity=severity,
        path=path,
        managed_id=managed_id,
        message=message,
        action=action,
        remediation=remediation,
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


def managed_block(managed_id: str, version: int | str, content: str) -> str:
    body = content.strip() + "\n"
    return (
        f"<!-- setup-context-driven:begin id={managed_id} version={version} -->\n"
        f"\n{body}\n"
        f"<!-- setup-context-driven:end id={managed_id} -->\n"
    )


def print_top_level_help() -> None:
    print(
        "\n".join(
            [
                "usage: context_setup.py [audit] [--repo PATH] [--format text|json]",
                "       context_setup.py apply --repo PATH [--format text|json]",
                "       context_setup.py sync-setups --source-dir PATH [--check] [--format text|json]",
                "       context_setup.py restore-skills --repo PATH --profile ID [--skill NAME ...] [--format text|json]",
                "",
                "Audit is the read-only default when no subcommand is supplied.",
                "Output formats: text, json. Results go to stdout; diagnostics go to stderr.",
                "Exit codes: 0 ok, 1 blocking findings, 2 invalid input, 3 decisions required or plan confirmation required/stale.",
                "",
                "Subcommands:",
                "  audit        Read bundled assets and repository state without writes.",
                "  apply        Write confirmed managed content through an atomic change plan.",
                "  restore-skills  Restore selected external skills from immutable provenance.",
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
        "--decision",
        action="append",
        default=[],
        help="Decision in ID=VALUE form. Repeat to resolve a concrete read-only plan.",
    )
    parser.add_argument(
        "--decision-file",
        action="append",
        default=[],
        help="Structured setup-context-driven/decisions/0.0.1 document. Repeatable.",
    )
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
    parser.add_argument(
        "--decision-file",
        action="append",
        default=[],
        help="Structured setup-context-driven/decisions/0.0.1 document. Repeatable.",
    )
    parser.add_argument(
        "--confirm-plan",
        help="Exact lowercase SHA-256 planDigest authorizing a non-empty apply.",
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


def restore_skills_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="context_setup.py restore-skills",
        description=(
            "Preview or apply immutable external Repository Skill Set restoration.\n\n"
            "Exit codes: 0 applied or already current; 1 source, proof, or write failure; "
            "2 invalid input; 3 confirmation required or stale."
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--repo", default=".", help="Repository root. Defaults to cwd.")
    parser.add_argument("--profile", required=True, help="Bundled profile id to restore.")
    parser.add_argument(
        "--skill",
        action="append",
        default=[],
        help="External required skill name. Repeat to select multiple skills; omit for all drift.",
    )
    parser.add_argument(
        "--source-dir",
        help="Offline Git checkout or bare object store containing the declared exact commit.",
    )
    parser.add_argument(
        "--confirm-plan",
        help="Exact lowercase SHA-256 planDigest authorizing a non-empty restoration.",
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
