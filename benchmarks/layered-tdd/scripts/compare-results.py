#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import statistics
from collections import defaultdict
from pathlib import Path


METRICS = (
    ("input_tokens", "Input tokens", ".0f"),
    ("output_tokens", "Output tokens", ".0f"),
    ("total_cost_usd", "Cost USD", ".4f"),
    ("model_invocations", "Invocations", ".0f"),
    ("input_tokens_per_invocation", "Input / invocation", ".0f"),
    ("cost_per_expected_layer_usd", "Cost / layer", ".4f"),
)


def main() -> None:
    parser = argparse.ArgumentParser(description="Compare medians from collected layered-TDD benchmark runs.")
    parser.add_argument("results", nargs="+", type=Path)
    args = parser.parse_args()

    grouped: dict[str, list[dict]] = defaultdict(list)
    for path in args.results:
        result = json.loads(path.read_text(encoding="utf-8"))
        grouped[result["variant"]].append(result)

    headers = ["Variant", "Runs", *[label for _, label, _ in METRICS], "Oracle"]
    print("| " + " | ".join(headers) + " |")
    print("| " + " | ".join(["---"] * len(headers)) + " |")
    for variant in sorted(grouped):
        runs = grouped[variant]
        cells = [variant, str(len(runs))]
        for key, _, format_spec in METRICS:
            values = [run["metrics"].get(key) for run in runs]
            numeric = [float(value) for value in values if value is not None]
            cells.append(format(statistics.median(numeric), format_spec) if numeric else "n/a")
        evaluated = [run.get("evaluation") for run in runs if run.get("evaluation") is not None]
        passed = sum(1 for item in evaluated if item["oracle_passed"] and item["boundary_passed"])
        cells.append(f"{passed}/{len(evaluated)}" if evaluated else "n/a")
        print("| " + " | ".join(cells) + " |")


if __name__ == "__main__":
    main()

