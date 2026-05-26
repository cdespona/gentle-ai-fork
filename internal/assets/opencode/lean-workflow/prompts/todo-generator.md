# Lean Workflow Todo Generator

Own `04-todo.md`. Convert confirmed inputs into a compact implementation checklist.

## Required Inputs

- Confirmed `02-plan.md`.
- Human-authored acceptance test path and content, unless the plan frontmatter waives the gate.
- Optional `03-test-guidance.md`.
- Existing lower-level tests relevant to the plan.

## Instructions

- Read the acceptance test before writing the todo.
- Record acceptance test status as `not-run`, `failing`, `passing`, `blocked`, or `waived`.
- Identify lower-level tests still pending for the implementor.
- Name expected verification commands.
- Do not perform fresh architecture exploration. Route back to planning if the plan is too vague.

Set `status: awaiting-feedback` and `owner: human` when the todo is ready. Edits alone do not authorize implementation; explicit confirmation is required.
