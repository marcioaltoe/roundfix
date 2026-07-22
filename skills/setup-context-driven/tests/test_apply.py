import json
import io
import hashlib
import subprocess
import sys
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

import context_setup  # noqa: E402
from context_setup import managed_block, parse_managed_blocks  # noqa: E402
from context_assets import clone_assets_to, load_asset_catalog, read_json_copy, write_json  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_context_setup as run_fixture_context_setup,
    snapshot_files,
    write_compliant_repository,
)


BASE_DECISIONS = [
    "spec.scaffold=true",
    "domain.layout=single-context",
    "triage.external=false",
    "autonomous.enabled=true",
    "runtime.backend=codex gpt-5.5 xhigh",
    "runtime.design=claude opus xhigh",
    "verification.gate=make verify",
    "language.generated=English",
    "secondbrain.enabled=false",
]


class ApplyCliTests(unittest.TestCase):
    def test_non_empty_apply_requires_matching_plan_digest_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)

            missing = run_apply_preview(repo, "rust-cli", BASE_DECISIONS)
            stale = run_apply_with_confirmation(repo, "rust-cli", BASE_DECISIONS, "0" * 64)
            malformed = run_apply_with_confirmation(repo, "rust-cli", BASE_DECISIONS, "not-a-digest")

            self.assertEqual(missing.returncode, 3, missing.stderr)
            self.assertFinding(missing, "plan.confirmation.required", "decision")
            self.assertRegex(json.loads(missing.stdout)["planDigest"], r"^[0-9a-f]{64}$")
            self.assertEqual(stale.returncode, 3, stale.stderr)
            self.assertFinding(stale, "plan.confirmation.stale", "decision")
            self.assertEqual(malformed.returncode, 2, malformed.stderr)
            self.assertFinding(malformed, "plan.confirmation.invalid", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_confirmed_plan_matches_complete_tree_delta_and_second_apply_is_empty(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)
            preview = run_apply_preview(repo, "rust-cli", BASE_DECISIONS)
            payload = json.loads(preview.stdout)

            applied = run_apply_with_confirmation(
                repo,
                "rust-cli",
                BASE_DECISIONS,
                payload["planDigest"],
            )
            after = snapshot_files(repo)

            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual(applied.stderr, "")
            self.assertEqual(json.loads(applied.stdout)["planDigest"], payload["planDigest"])
            planned_delta = {
                (change["path"], change["beforeDigest"], change["afterDigest"])
                for change in payload["plannedChanges"]
                if change["beforeDigest"] != change["afterDigest"]
            }
            observed_delta = tree_delta(before, after)
            self.assertEqual(planned_delta, observed_delta)

            second = run_apply_preview(repo, "rust-cli", BASE_DECISIONS)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(json.loads(second.stdout)["plannedChanges"], [])
            self.assertEqual(snapshot_files(repo), after)

    def test_catalog_path_move_plans_rename_and_reference_edit_without_delta_gaps(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            write_compliant_repository(repo, "rust-cli")
            before = snapshot_files(repo)
            temp_skill = root / "setup-context-driven"
            clone_assets_to(SKILL_ROOT, temp_skill)
            module_path = temp_skill / "assets" / "modules" / "rust.json"
            module = read_json_copy(module_path)
            module["supportingGuides"][0]["path"] = "docs/agents/rust-renamed.md"
            write_json(module_path, module)
            catalog = load_asset_catalog(temp_skill)

            result, invalid_input, plan = context_setup.plan_apply(
                repo=repo,
                catalog=catalog,
                profile_override="rust-cli",
                decision_args=[],
            )
            payload = result.to_json()

            self.assertFalse(invalid_input)
            self.assertEqual(result.findings, [])
            rename = next(
                change
                for change in payload["plannedChanges"]
                if change["action"] == "rename managed content"
            )
            self.assertEqual(rename["path"], "docs/agents/rust-renamed.md")
            self.assertEqual(rename["fromPath"], "docs/agents/rust.md")
            reference_edit = next(
                change
                for change in payload["plannedChanges"]
                if change["action"] == "edit managed references"
                and change["managedId"] == "root.rust"
            )
            self.assertTrue(reference_edit["referenceEdits"])

            context_setup.apply_change_plan(repo, plan)
            after = snapshot_files(repo)
            planned_delta = {
                (change["path"], change["beforeDigest"], change["afterDigest"])
                for change in payload["plannedChanges"]
                if change["beforeDigest"] != change["afterDigest"]
            }
            self.assertEqual(planned_delta, tree_delta(before, after))
            self.assertFalse((repo / "docs" / "agents" / "rust.md").exists())
            self.assertTrue((repo / "docs" / "agents" / "rust-renamed.md").is_file())

    def test_apply_creates_manifest_root_blocks_and_guides(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 0, result.stderr)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(manifest["profile"], "rust-cli")
            self.assertEqual(
                manifest["modules"],
                [
                    "core",
                    "context-workflow",
                    "rust",
                    "cli-surface",
                    "autonomous-work",
                    "spec-workflow",
                ],
            )
            self.assertIn("runtime.design", manifest["decisions"])
            self.assertTrue(manifest["managedArtifacts"])
            self.assertIn("<!-- setup-context-driven:begin id=root.rust version=2 -->", (repo / "AGENTS.md").read_text(encoding="utf-8"))
            self.assertTrue((repo / "docs" / "agents" / "rust.md").is_file())
            install_profile_skills(repo, "rust-cli")
            self.assertEqual(run_audit(repo).returncode, 0)

    def test_apply_preserves_custom_content_and_is_idempotent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "AGENTS.md").write_text(
                "custom before\n\ncustom after\n",
                encoding="utf-8",
            )

            first = run_apply(repo, "rust-cli", BASE_DECISIONS)
            after_first = snapshot_files(repo)
            second = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(snapshot_files(repo), after_first)
            content = (repo / "AGENTS.md").read_text(encoding="utf-8")
            self.assertIn("custom before\n\ncustom after\n", content)

    def test_missing_decision_prevents_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            before = snapshot_files(repo)
            decisions = [item for item in BASE_DECISIONS if not item.startswith("runtime.design=")]

            result = run_apply(repo, "rust-cli", decisions)

            self.assertEqual(result.returncode, 3)
            self.assertFinding(result, "decision.required", "decision")
            self.assertEqual(snapshot_files(repo), before)

    def test_invalid_duplicate_markers_prevent_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            duplicated = (
                "<!-- setup-context-driven:begin id=root.core version=1 -->\n"
                "first\n"
                "<!-- setup-context-driven:end id=root.core -->\n"
                "<!-- setup-context-driven:begin id=root.core version=1 -->\n"
                "second\n"
                "<!-- setup-context-driven:end id=root.core -->\n"
            )
            (repo / "AGENTS.md").write_text(duplicated, encoding="utf-8")
            before = snapshot_files(repo)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.block.duplicate", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_failure_before_commit_preserves_target_files(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            docs_agents = repo / "docs" / "agents"
            docs_agents.mkdir(parents=True)
            (repo / "AGENTS.md").write_text("custom root\n", encoding="utf-8")
            before = snapshot_files(repo)

            with mock.patch.object(context_setup.Path, "replace", side_effect=OSError("injected replace failure")):
                result = run_apply_in_process(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.apply.failed", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_postwrite_delta_mismatch_restores_every_original_byte(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "AGENTS.md").write_text("custom root\n", encoding="utf-8")
            before = snapshot_files(repo)
            preview = run_apply_in_process(
                repo,
                "rust-cli",
                BASE_DECISIONS,
                auto_confirm=False,
            )
            digest = json.loads(preview.stdout)["planDigest"]
            original_replace = context_setup.Path.replace

            def replace_then_corrupt(source, target):
                result = original_replace(source, target)
                if Path(target).name == "setup-context.json":
                    (repo / "AGENTS.md").write_text("forced mismatch\n", encoding="utf-8")
                return result

            with mock.patch.object(context_setup.Path, "replace", autospec=True, side_effect=replace_then_corrupt):
                result = run_apply_in_process(
                    repo,
                    "rust-cli",
                    BASE_DECISIONS,
                    confirm_plan=digest,
                    auto_confirm=False,
                )

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "managed.apply.failed", "error")
            self.assertEqual(snapshot_files(repo), before)

    def test_apply_preserves_existing_predictable_temp_named_user_file(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            existing_temp = repo / ".AGENTS.md.setup-context.tmp"
            existing_temp.write_text("user-owned temp file\n", encoding="utf-8")

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual(existing_temp.read_text(encoding="utf-8"), "user-owned temp file\n")

    def test_manifest_artifact_paths_must_stay_inside_repo_before_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir) / "repo"
            outside = Path(temp_dir) / "outside.md"
            write_compliant_repository(repo, "rust-cli")
            outside.write_text(
                managed_block("stale.outside", 1, "setup-owned stale outside file\n"),
                encoding="utf-8",
            )
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["managedArtifacts"].append(
                {
                    "id": "stale.outside",
                    "path": "../outside.md",
                    "kind": "guide",
                    "module": "stale",
                    "template": "template.stale",
                    "version": 1,
                    "digest": "0" * 64,
                }
            )
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
            repo_before = snapshot_files(repo)
            outside_before = outside.read_text(encoding="utf-8")

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "manifest.invalid", "error")
            self.assertEqual(snapshot_files(repo), repo_before)
            self.assertEqual(outside.read_text(encoding="utf-8"), outside_before)

    def test_invalid_explicit_decisions_are_invalid_input_without_writes(self):
        cases = [
            ("domain.layout=dual-context", "enum"),
            ("autonomous.enabled=maybe", "boolean"),
        ]

        for decision, name in cases:
            with self.subTest(name=name):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    before = snapshot_files(repo)
                    decisions = [
                        decision if item.split("=", 1)[0] == decision.split("=", 1)[0] else item
                        for item in BASE_DECISIONS
                    ]

                    result = run_apply(repo, "rust-cli", decisions)

                    self.assertEqual(result.returncode, 2)
                    self.assertFinding(result, "decision.value.invalid", "error")
                    self.assertEqual(snapshot_files(repo), before)

    def test_obsolete_managed_artifacts_are_removed_without_unowned_deletions(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "go-cli-tui")
            go_guide = repo / "docs" / "agents" / "go.md"
            go_guide.write_text(
                "custom go intro\n" + go_guide.read_text(encoding="utf-8") + "custom go tail\n",
                encoding="utf-8",
            )
            (repo / "docs" / "agents" / "local.md").write_text("local guide\n", encoding="utf-8")

            result = run_apply(repo, "rust-cli", [])

            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertFalse((repo / "docs" / "agents" / "tui.md").exists())
            self.assertTrue((repo / "docs" / "agents" / "local.md").exists())
            self.assertEqual(go_guide.read_text(encoding="utf-8"), "custom go intro\ncustom go tail\n")
            root_blocks, _ = parse_managed_blocks(
                Path("AGENTS.md"),
                (repo / "AGENTS.md").read_text(encoding="utf-8"),
            )
            self.assertIn("root.rust", root_blocks)
            self.assertNotIn("root.go", root_blocks)

    def test_unmarked_conditional_guide_outside_manifest_is_never_removed(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            unmarked = repo / "docs" / "agents" / "tui.md"
            unmarked.write_text("repository-authored conditional guide\n", encoding="utf-8")

            preview = run_apply_preview(repo, "rust-cli", [])
            payload = json.loads(preview.stdout)

            self.assertFalse(
                any(
                    change["path"] == "docs/agents/tui.md"
                    and change["action"].startswith("remove")
                    for change in payload["plannedChanges"]
                ),
                payload,
            )
            applied = run_apply(repo, "rust-cli", [])
            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual(unmarked.read_text(encoding="utf-8"), "repository-authored conditional guide\n")

    def test_ambiguous_old_inventory_requires_explicit_preserve_or_remove(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "go-cli-tui")
            stale = repo / "docs" / "agents" / "tui.md"
            stale.write_text("repository-authored replacement\n", encoding="utf-8")
            before = snapshot_files(repo)

            blocked = run_apply_preview(repo, "rust-cli", [])
            self.assertEqual(blocked.returncode, 3, blocked.stderr)
            self.assertFinding(blocked, "decision.required", "decision")
            self.assertEqual(snapshot_files(repo), before)

            preserved = run_apply(repo, "rust-cli", ["removal.guide.tui-surface=preserve"])
            self.assertEqual(preserved.returncode, 0, preserved.stderr)
            self.assertEqual(stale.read_text(encoding="utf-8"), "repository-authored replacement\n")

            write_compliant_repository(repo, "go-cli-tui")
            stale.write_text("repository-authored replacement\n", encoding="utf-8")
            removed = run_apply(repo, "rust-cli", ["removal.guide.tui-surface=remove"])
            self.assertEqual(removed.returncode, 0, removed.stderr)
            self.assertFalse(stale.exists())

    def test_adoption_decision_required_before_owning_unmarked_guide(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            rust_guide = repo / "docs" / "agents" / "rust.md"
            rust_guide.parent.mkdir(parents=True)
            rust_guide.write_text("repo-authored rust notes\n", encoding="utf-8")
            before = snapshot_files(repo)

            blocked = run_apply(repo, "rust-cli", BASE_DECISIONS)
            self.assertEqual(blocked.returncode, 3)
            self.assertFinding(blocked, "decision.required", "decision")
            self.assertEqual(snapshot_files(repo), before)

            adopted = run_apply(
                repo,
                "rust-cli",
                BASE_DECISIONS + ["adoption.guide.rust=true"],
            )

            self.assertEqual(adopted.returncode, 0, adopted.stderr)
            self.assertIn(
                "<!-- setup-context-driven:begin id=guide.rust version=2 -->",
                rust_guide.read_text(encoding="utf-8"),
            )
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertEqual(manifest["decisions"]["adoption.guide.rust"]["value"], True)

    def test_unsupported_manifest_version_blocks_apply_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["schemaVersion"] = 0
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
            before = snapshot_files(repo)

            result = run_apply(repo, "rust-cli", BASE_DECISIONS)

            self.assertEqual(result.returncode, 1)
            self.assertFinding(result, "manifest.migration-required", "error")
            self.assertEqual(snapshot_files(repo), before)

    def assertFinding(self, result, code, severity):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        self.assertEqual(matches[0]["severity"], severity)


def run_apply(repo, profile, decisions):
    preview = run_apply_preview(repo, profile, decisions)
    if preview.returncode != 3:
        return preview
    payload = json.loads(preview.stdout)
    if not any(
        finding["code"] == "plan.confirmation.required"
        for finding in payload["findings"]
    ):
        return preview
    return run_apply_with_confirmation(repo, profile, decisions, payload["planDigest"])


def run_apply_preview(repo, profile, decisions):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", profile]
    for decision in decisions:
        args.extend(["--decision", decision])
    return run_context_setup(*args)


def run_apply_with_confirmation(repo, profile, decisions, confirmation):
    args = [
        "apply", "--repo", str(repo), "--format", "json", "--profile", profile,
        "--confirm-plan", confirmation,
    ]
    for decision in decisions:
        args.extend(["--decision", decision])
    return run_context_setup(*args)


def run_apply_in_process(repo, profile, decisions, confirm_plan=None, auto_confirm=True):
    args = ["apply", "--repo", str(repo), "--format", "json", "--profile", profile]
    for decision in decisions:
        args.extend(["--decision", decision])
    if confirm_plan is not None:
        args.extend(["--confirm-plan", confirm_plan])
    result = run_context_setup_in_process(*args)
    if not auto_confirm or result.returncode != 3:
        return result
    payload = json.loads(result.stdout)
    if not any(finding["code"] == "plan.confirmation.required" for finding in payload["findings"]):
        return result
    args.extend(["--confirm-plan", payload["planDigest"]])
    return run_context_setup_in_process(*args)


def run_audit(repo):
    return run_context_setup("audit", "--repo", str(repo), "--format", "json")


def run_context_setup(*args):
    return run_fixture_context_setup(*args)


def run_context_setup_in_process(*args):
    stdout = io.StringIO()
    stderr = io.StringIO()
    with redirect_stdout(stdout), redirect_stderr(stderr):
        returncode = context_setup.main(list(args))
    return subprocess.CompletedProcess(args, returncode, stdout.getvalue(), stderr.getvalue())


def tree_delta(before, after):
    paths = set(before) | set(after)
    return {
        (path, digest(before.get(path)), digest(after.get(path)))
        for path in paths
        if before.get(path) != after.get(path)
    }


def digest(content):
    return hashlib.sha256(content).hexdigest() if content is not None else None


if __name__ == "__main__":
    unittest.main()
