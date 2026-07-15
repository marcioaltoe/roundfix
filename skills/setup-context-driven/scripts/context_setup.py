#!/usr/bin/env python3
"""Audit setup-context-driven managed agent instructions."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
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
class ManagedBlock:
    managed_id: str
    version: int
    body: str


@dataclass(frozen=True)
class AuditResult:
    findings: list[Finding]

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
        return {
            "schemaVersion": AUDIT_SCHEMA_VERSION,
            "ok": self.ok,
            "summary": self.summary,
            "findings": [finding.to_json() for finding in sorted_findings(self.findings)],
        }


def main(argv: list[str] | None = None) -> int:
    args = list(sys.argv[1:] if argv is None else argv)
    if args and args[0] in {"-h", "--help"}:
        print_top_level_help()
        return 0
    if args and args[0] == "apply":
        return parse_unimplemented_command("apply", args[1:])
    if args and args[0] == "sync-setups":
        return parse_unimplemented_command("sync-setups", args[1:])
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

    result, invalid_input = audit_repository(repo, catalog, options.profile)
    render_result(result, options.format)
    return exit_code_for(result, invalid_input)


def audit_repository(
    repo: Path,
    catalog: AssetCatalog,
    profile_override: str | None = None,
) -> tuple[AuditResult, bool]:
    findings: list[Finding] = []
    manifest, invalid_input = load_manifest(repo, findings)
    if manifest is None:
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

    ordered_modules = catalog.ordered_modules_by_profile[profile_id]
    validate_manifest_shape(manifest, profile_id, ordered_modules, catalog, findings)
    expected_artifacts = expected_artifacts_for_profile(catalog, profile_id)
    validate_manifest_artifacts(manifest, expected_artifacts, findings)
    validate_documents(repo, expected_artifacts, findings)

    return AuditResult(sorted_findings(findings)), invalid_input


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


def parse_managed_blocks(
    relative_path: Path,
    content: str,
) -> tuple[dict[str, ManagedBlock], list[Finding]]:
    findings: list[Finding] = []
    blocks: dict[str, ManagedBlock] = {}
    open_marker: tuple[str, int, int] | None = None

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
            open_marker = (managed_id, version, marker.end())
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
            open_id, version, body_start = open_marker
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
            blocks[managed_id] = ManagedBlock(
                managed_id=managed_id,
                version=version,
                body=content[body_start:marker.start()],
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
        managed_id, _, _ = open_marker
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
) -> list[ExpectedArtifact]:
    artifacts: list[ExpectedArtifact] = []
    templates_root = Path(__file__).resolve().parents[1] / "assets" / "templates"
    for module_id in catalog.ordered_modules_by_profile[profile_id]:
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
    if not result.findings:
        return "setup-context-driven audit: ok"

    lines = [
        "setup-context-driven audit: findings",
        (
            f"errors={result.summary['errors']} decisions={result.summary['decisions']} "
            f"warnings={result.summary['warnings']} info={result.summary['info']}"
        ),
    ]
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
    return "\n".join(lines)


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
                "  apply        Planned safe correction command; not implemented in this slice.",
                "  sync-setups  Planned snapshot authoring command; not implemented in this slice.",
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
        help="Reserved for optional skill cleanup reporting.",
    )
    parser.add_argument(
        "--setups-dir",
        help="Reserved for explicit canonical setup drift checks.",
    )
    return parser


def parse_unimplemented_command(command: str, args: list[str]) -> int:
    parser = argparse.ArgumentParser(
        prog=f"context_setup.py {command}",
        description=f"{command} is documented for the setup workflow but not implemented in this slice.",
    )
    if command == "apply":
        parser.add_argument("--repo", default=".")
        parser.add_argument("--format", choices=["text", "json"], default="text")
        parser.add_argument("--profile")
        parser.add_argument("--decision", action="append", default=[])
    else:
        parser.add_argument("--source-dir", required="--help" not in args and "-h" not in args)
        parser.add_argument("--check", action="store_true")
        parser.add_argument("--format", choices=["text", "json"], default="text")
    parser.parse_args(args)
    print(f"{command} is not implemented by this task.", file=sys.stderr)
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
