#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import re
from pathlib import Path


TOKEN_PATTERNS = {
    "input_tokens": re.compile(r"^\s*Input:\s+([\d,]+)\s+tokens\s*$", re.MULTILINE),
    "output_tokens": re.compile(r"^\s*Output:\s+([\d,]+)\s+tokens\s*$", re.MULTILINE),
    "total_tokens": re.compile(r"^\s*Total:\s+([\d,]+)\s+tokens\s*$", re.MULTILINE),
}
COST_LINE = re.compile(r"^\s{2}([A-Za-z0-9_.$-]+):\s+\$([\d.]+)(?:\s+\(\d+%\))?\s*$", re.MULTILINE)
TOTAL_COST = re.compile(r"^\s*Total:\s+\$([\d.]+)\s*$", re.MULTILINE)


def last_number(pattern: re.Pattern[str], text: str, cast):
    matches = pattern.findall(text)
    if not matches:
        return None
    return cast(matches[-1].replace(",", ""))


def main() -> None:
    parser = argparse.ArgumentParser(description="Collect a Conductor benchmark summary into JSON.")
    parser.add_argument("--variant", required=True)
    parser.add_argument("--run-log", required=True, type=Path)
    parser.add_argument("--project", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--notes", default="")
    args = parser.parse_args()

    text = args.run_log.read_text(encoding="utf-8", errors="replace")
    project = args.project.resolve()
    case = json.loads((project / ".benchmark" / "case.json").read_text(encoding="utf-8"))
    evaluation_path = project / ".benchmark" / "evaluation.json"
    evaluation = json.loads(evaluation_path.read_text(encoding="utf-8")) if evaluation_path.exists() else None

    metrics = {name: last_number(pattern, text, int) for name, pattern in TOKEN_PATTERNS.items()}
    metrics["total_cost_usd"] = last_number(TOTAL_COST, text, float)

    agent_costs: dict[str, float] = {}
    invocation_count = 0
    for agent, cost in COST_LINE.findall(text):
        if agent == "Total":
            continue
        invocation_count += 1
        agent_costs[agent] = agent_costs.get(agent, 0.0) + float(cost)

    metrics["model_invocations"] = invocation_count
    metrics["input_tokens_per_invocation"] = (
        metrics["input_tokens"] / invocation_count if metrics["input_tokens"] is not None and invocation_count else None
    )
    metrics["cost_per_invocation_usd"] = (
        metrics["total_cost_usd"] / invocation_count if metrics["total_cost_usd"] is not None and invocation_count else None
    )
    expected_layers = int(case["expected_layers"])
    metrics["input_tokens_per_expected_layer"] = (
        metrics["input_tokens"] / expected_layers if metrics["input_tokens"] is not None else None
    )
    metrics["cost_per_expected_layer_usd"] = (
        metrics["total_cost_usd"] / expected_layers if metrics["total_cost_usd"] is not None else None
    )

    result = {
        "variant": args.variant,
        "case": case["case"],
        "expected_layers": expected_layers,
        "test_command": case["test_command"],
        "graphify": case["graphify"],
        "metrics": metrics,
        "agent_costs_usd": agent_costs,
        "evaluation": evaluation,
        "notes": args.notes,
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
