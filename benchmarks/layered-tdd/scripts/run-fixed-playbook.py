#!/usr/bin/env python3
"""Run the fixed layered-TDD benchmark gate playbook without timing races."""

from __future__ import annotations

import argparse
import json
import os
from pathlib import Path

import pexpect


TIMEOUT_SECONDS = 900
UNEXPECTED_GATES = (
    "agent_test_checkpoint_gate",
    "checkpoint_gate",
    "preflight_failure_gate",
    "graphify_refresh_failure_gate",
)


def resolve_gate(child: pexpect.spawn, name: str, option: int, field: str = "", value: str = "") -> None:
    unexpected = "|".join(UNEXPECTED_GATES)
    matched = child.expect(
        [
            rf"Agent: {name} ",
            rf"Agent: (?P<unexpected>{unexpected}) ",
        ]
    )
    if matched == 1:
        raise RuntimeError(
            f"unexpected gate before {name}: {child.match.group('unexpected')}"
        )
    child.expect(r"Select option \[[0-9/]+\]: ")
    child.sendline(str(option))
    if field:
        child.expect(rf"\r?\n  {field}: ")
        child.sendline(value)


def main() -> None:
    parser = argparse.ArgumentParser(description="Run the fixed layered-TDD benchmark gate playbook.")
    parser.add_argument("project", type=Path)
    parser.add_argument("--run-log", required=True, type=Path)
    args = parser.parse_args()

    project = args.project.resolve()
    run_log = args.run_log.resolve()
    try:
        run_log.relative_to(project)
    except ValueError:
        pass
    else:
        parser.error("--run-log must be outside the benchmark project")

    request = (project / ".benchmark" / "request.md").read_text(encoding="utf-8")
    case = json.loads(
        (project / ".benchmark" / "case.json").read_text(encoding="utf-8")
    )
    run_log.parent.mkdir(parents=True, exist_ok=True)
    debug_log = run_log.with_name(f"{run_log.name}.conductor-debug.log")

    command = [
        "run",
        "workflows/conductor/layered-tdd.yaml",
        "--workspace-instructions",
        "--log-file",
        str(debug_log),
        "--input",
        "task_slug=benchmark-idempotent-cancellation",
        "--input",
        f"test_command={case['test_command']}",
        "--input",
        "lint_command=make lint",
        "--input",
        "security_command=make audit",
        "--input",
        f"request={request}",
    ]

    environment = os.environ.copy()
    environment.setdefault("NO_COLOR", "1")
    with run_log.open("w", encoding="utf-8") as log:
        child = pexpect.spawn(
            "conductor",
            command,
            cwd=str(project),
            env=environment,
            encoding="utf-8",
            timeout=TIMEOUT_SECONDS,
        )
        child.logfile_read = log

        resolve_gate(child, "graphify_update_gate", 3)
        resolve_gate(child, "requirements_gate", 1, "feedback")

        layers = ("L10-domain", "L20-service", "L30-http")
        for index, layer in enumerate(layers):
            resolve_gate(child, "layer_selection_gate", 1, "selected_layer", layer)
            resolve_gate(child, "layer_todo_gate", 1, "feedback")
            resolve_gate(child, "red_suite_evidence_gate", 1, "feedback")
            approval_option = 2 if index == len(layers) - 1 else 1
            resolve_gate(child, "layer_approval_gate", approval_option, "feedback")

        resolve_gate(child, "memory_gate", 2, "feedback")
        child.expect(pexpect.EOF)
        child.close()

    if child.exitstatus != 0:
        raise SystemExit(child.exitstatus or 1)


if __name__ == "__main__":
    main()
