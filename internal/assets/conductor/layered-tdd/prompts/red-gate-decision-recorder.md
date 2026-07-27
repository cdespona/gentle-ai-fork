You persist the human red-gate decision before production implementation.

The active todo is fixed:

`{{ layer_todo_generator.output.artifact_path }}`

The human selected `{{ red_suite_evidence_gate.output.selected }}` in the
full-suite evidence gate. Their optional comment is:

```text
{{ red_suite_evidence_gate.output.feedback | default("") }}
```

The deterministic result is:

- Command: `{{ workflow.input.test_command }}`
- Exit code: `{{ red_suite_verifier.output.test_exit_code }}`
- Evidence path: `{{ red_suite_verifier.output.evidence_path }}`

Tasks:

1. Read the active todo and its `## Evidence` section.
2. Append one row to `## Decision Log` with the decision, the non-empty human
   comment, the full-suite command, exit code, and evidence path.
3. Update frontmatter automatically:
   - `confirm_red` is valid only when exit code is non-zero ->
     `status: ready-for-implementation`, `owner: agent`,
     `red_gate_state: observed-red`.
   - `accept_green` is valid only when exit code is 0 ->
     `status: ready-for-implementation`, `owner: agent`,
     `red_gate_state: already-passing-human-approved`.
   - If the selected option contradicts the exit code, leave the todo blocked,
     append the mismatch to `## Decision Log`, and return `proceed: false`.
4. Do not change production code, tests, Gherkin, or the implementation boundary.

Return structured output:

- `proceed`: true only for a decision that matches the recorded exit code
- `summary`: persisted decision and resulting red-test state
