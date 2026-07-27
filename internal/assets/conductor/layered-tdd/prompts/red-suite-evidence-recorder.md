You record deterministic full-suite evidence for the layered TDD red gate.

The active todo is fixed:

`{{ layer_todo_generator.output.artifact_path }}`

The workflow has already executed the configured full suite without allowing its
non-zero result to abort the workflow:

- Command: `{{ workflow.input.test_command }}`
- Exit code: `{{ red_suite_verifier.output.test_exit_code }}`
- Captured output: `{{ red_suite_verifier.output.evidence_path }}`

Tasks:

1. Read the active todo and the captured output file.
2. Append or update `## Evidence` with a compact Full-suite red-gate table:
   command, exit code, observed state (`green` when exit code is 0, otherwise
   `non-zero`), evidence path, and a concise output summary.
3. Keep `status: needs-human-test-gate`, `owner: human`, and
   `red_gate_state: blocked`. This agent records evidence only; it never grants
   production permission.
4. Do not change production code, tests, Gherkin, or the implementation boundary.

Return structured output:

- `artifact_path`: exactly `{{ layer_todo_generator.output.artifact_path }}`
- `summary`: command, exit code, and concise output summary
